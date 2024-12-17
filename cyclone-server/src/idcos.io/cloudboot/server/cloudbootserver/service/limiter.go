package service

import (
	"net/http"

	"github.com/voidint/binding"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/limiter"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/model"
)

// ReturnLimiterTokenBySNReq 归还限流器令牌请求结构体
type ReturnLimiterTokenBySNReq struct {
	SN    string        `json:"sn"`
	Token limiter.Token `json:"-"`
}

// FieldMap 请求字段映射
func (reqData *ReturnLimiterTokenBySNReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.SN: "sn",
	}
}

// Validate 结构体数据校验
func (reqData *ReturnLimiterTokenBySNReq) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(r.Context())
	token, _ := repo.GetTokenBySN(reqData.SN)
	//if err == gorm.ErrRecordNotFound {
	//	errs.Add([]string{"sn"}, binding.BusinessError, "令牌已归还或未曾申请令牌")
	//	return errs
	//}
	//if err != nil {
	//	errs.Add([]string{"sn"}, binding.SystemError, "系统内部错误")
	//	return errs
	//}
	reqData.Token = limiter.Token(token)
	return errs
}

// ReturnLimiterTokenBySN 归还目标设备的限流器令牌
func ReturnLimiterTokenBySN(log logger.Logger, repo model.Repo, lim limiter.Limiter, conf *config.Config, reqData *ReturnLimiterTokenBySNReq) (err error) {
	log.Infof("The device(%s) attempts to return token", reqData.SN)

	//把机器关了
	o, err := BatchOperateOOBPower(log, repo, PowerOff, conf, false, []string{reqData.SN})
	if err != nil {
		log.Errorf("power off device:%s fail,err:%v,stderr:%s", reqData.SN, err, string(o))
		//continue ?
	}
	if lim != nil {
		bucket, err := lim.Route(reqData.SN)
		if err != nil {
			return err
		}
		return bucket.Return(reqData.SN, reqData.Token)
	}
	return err
}

type ReturnLimiterTokensReq struct {
	Tokens []*ReturnLimiterTokenBySNReq `json:"tokens"`
}

// FieldMap 请求字段映射
func (reqData *ReturnLimiterTokensReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.Tokens: "tokens",
	}
}

// Validate 结构体数据校验
func (reqData *ReturnLimiterTokensReq) Validate(r *http.Request, errs binding.Errors) binding.Errors {
	for _, t := range reqData.Tokens {
		errs = t.Validate(r, errs)
		if len(errs) != 0 {
			return errs
		}
	}
	return nil
}

// ReturnLimiterToken 归还目标设备的限流器令牌
func ReturnLimiterTokens(log logger.Logger, repo model.Repo, lim limiter.Limiter, conf *config.Config, reqData *ReturnLimiterTokensReq) (succeed int, err error) {
	for _, t := range reqData.Tokens {
		//log.Infof("The device(%s) attempts to return token", t.SN)
		//isFailure := false //记录各操作是否有错误
		//if lim != nil {
		//	bucket, err := lim.Route(t.SN)
		//	if err != nil {
		//		log.Errorf("The device(%s) get token bucket fail,%v", t.SN, err)
		//		isFailure = true
		//		//continue
		//	} else if bucket != nil {
		//		if err = bucket.Return(t.SN, t.Token); err != nil {
		//			log.Errorf("The device(%s) return token fail,%v", t.SN, err)
		//			isFailure = true
		//			//continue
		//		}
		//	}
		//}
		//
		////把机器关了
		//o, errPwr := BatchOperateOOBPower(log, repo, PowerOff, conf, false, []string{t.SN})
		//if errPwr != nil {
		//	log.Errorf("power off device:%s fail,err:%v,stderr:%s", t.SN, errPwr, string(o))
		//	isFailure = true
		//	continue
		//}
		//if isFailure == false {
		//	succeed++
		//}

		err := ReturnLimiterTokenBySN(log, repo, lim, conf, t)
		if err != nil && err != limiter.ErrInvalidOrReturnedToken {
			log.Error(err)
			continue
		}
		succeed++
	}
	return succeed, nil
}
