package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
	"idcos.io/cloudboot/utils"
	http3 "idcos.io/cloudboot/utils/http"
)

//RoutesConf 路由配置
type RoutesConf struct {
	HTTPMethod               string `json:"http_method"`
	URL                      string `json:"url"`
	CategoryCode             string `json:"category_code"`
	CategoryName             string `json:"category_name"`
	ReqParamDesensitization  string `json:"req_param_desensitization"`
	RespParamDesensitization string `json:"resp_param_desensitization"`
	GetDBMethod              string `json:"method"`
}

var once sync.Once
var routesConf []*RoutesConf

//GetRoutesConf 获取路由的配置信息
func GetRoutesConf(repo model.Repo, logger logger.Logger) {
	//TODO 前期需要大量的测试，不断地修改数据库，所以此处的单例先取消
	if len(routesConf) != 0 {
		return
	}

	once.Do(func() {
		sys, err := repo.GetSystemSetting("route_conf")
		if err != nil {
			logger.Errorf("GetSystemSetting route_conf err,%s", err.Error())
			return
		}

		if err = json.Unmarshal([]byte(sys.Value), &routesConf); err != nil {
			logger.Errorf("json unmarshal route_conf value err,%s", err.Error())
		}
	})
}

//Responser 手工构造response，用于获取http response信息
type Responser struct {
	w   http.ResponseWriter
	buf *bytes.Buffer
}

//Header  Header
func (r *Responser) Header() http.Header {
	return r.w.Header()
}

//Write Write
func (r *Responser) Write(buf []byte) (int, error) {
	r.buf.Write(buf)
	return r.w.Write(buf)
}

//WriteHeader WriteHeader
func (r *Responser) WriteHeader(statusCode int) {
	r.w.WriteHeader(statusCode)
}

// OperateInterceptor 操作记录拦截器
func OperateInterceptor(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// 声明
		// 源和目标数据
		var source, destination interface{}
		url := strings.TrimSpace(r.URL.String())
		method := r.Method
		log, _ := LoggerFromContext(r.Context())
		repo, _ := RepoFromContext(r.Context())
		t1 := time.Now()

		//获取http RequestBody
		reqBody, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))

		if http.MethodGet == r.Method {
			next.ServeHTTP(w, r)
			return
		}

		// 根据url获取配置信息
		conf := isURLExist(url, method, repo, log)

		if conf == nil {
			log.Debugf("[OperateInterceptor]--> config is null")
			next.ServeHTTP(w, r)
			return
		}

		ww := &Responser{
			w:   w,
			buf: bytes.NewBufferString(""),
		}

		//catch 异常，进行log，不影响http请求的结果
		defer func() {
			if rvr := recover(); rvr != nil {
				next.ServeHTTP(w, r)
				log.Errorf("[OperateInterceptor]--> operate interceptor fail：%s\n  stack: %s", rvr, debug.Stack())
			}
		}()

		//处理API调用记录信息
		defer func() {
			if r.ContentLength >= 10000 {
				log.Debugf("[OperateInterceptor]--> HttpRequest ContentLength is too large, do not SaveAPILog ------> url: %s, method: %s,time: %s",
					r.URL.String(), r.Method, fmt.Sprintf("%#v", time.Now().Sub(t1)))
				return
			}

			if len(ww.buf.String()) >= 10000 {
				log.Debugf("[OperateInterceptor]--> HttpReponse Body is too large, ignore it  ------> url: %s, method: %s,time: %s",
					r.URL.String(), r.Method, fmt.Sprintf("%#v", time.Now().Sub(t1)))

				SaveAPILog(r, conf, repo, log, `{"status":"success","content":{},"message":"处理成功"}`, "", time.Now().Sub(t1).Seconds())
				return

			}

			log.Debugf("[OperateInterceptor]--> start to  SaveAPILog ------> url: %s, method: %s, form: %s, response: %s, time: %s",
				r.URL.String(), r.Method, fmt.Sprintf("%#v", r.Form), ww.buf.String(), fmt.Sprintf("%#v", time.Now().Sub(t1)))
			SaveAPILog(r, conf, repo, log, ww.buf.String(), string(reqBody), time.Now().Sub(t1).Seconds())
		}()

		// 若conf当中的getmethod没有配置
		if conf.GetDBMethod == "" {
			log.Warnf("[OperateInterceptor]--> route config GetDBMethod is null, config details: %s", utils.ToJsonString(conf))
			next.ServeHTTP(ww, r)
			return
		}

		//从URL当中拿到ID
		ids := GetIDsFromReq(r, conf)

		if len(ids) > 0 {
			source = GetRespVal(repo, conf.GetDBMethod, ids)
		}

		//若url在配置当中存在，保存操作记录信息
		defer func() {
			if strings.Contains(ww.buf.String(), "操作成功") {
				// 若为新增，需要去拿id
				switch method {
				case http.MethodPost:
					values := GetIDsFromResp(r, ww)
					if len(values) > 0 {
						destination = GetRespVal(repo, conf.GetDBMethod, values)
					}
				case http.MethodPut:
					if len(ids) > 0 {
						destination = GetRespVal(repo, conf.GetDBMethod, ids)
					} else {
						values := GetIDsFromResp(r, ww)
						if len(values) > 0 {
							destination = GetRespVal(repo, conf.GetDBMethod, values)
						}
					}
				case http.MethodDelete:
					destination = struct{}{}
				}

				saveOperate(r, conf, utils.ToJsonString(source), utils.ToJsonString(destination))
			} else {
				log.Debugf("[OperateInterceptor]--> http response not success ------> url: %s, method: %s,form: %s,  time: %s",
					r.URL.String(), r.Method, fmt.Sprintf("%#v", r.Form), fmt.Sprintf("%#v", time.Now().Sub(t1)))
			}

		}()
		next.ServeHTTP(ww, r)
	}

	return http.HandlerFunc(fn)
}

//SaveAPILog 保存API log信息
func SaveAPILog(req *http.Request, conf *RoutesConf, repo model.Repo, log logger.Logger, resp string, reqBody string, cost float64) {

	user, _ := LoginUserFromContext(req.Context())

	userName := "unknown user"
	if user != nil {
		userName = user.LoginName
	}

	var resbody http3.RespBody
	if resp != "" {
		log.Debugf("SaveAPILog ------> begin to unmarshal resp: %v", resp)
		if err := json.Unmarshal([]byte(resp), &resbody); err != nil {
			log.Error(err)
			return
		}
	} else {
		log.Debugf("SaveAPILog ------> resp is %v, skip it", resp)	
		return
	}

	description := ""
	reqParamDesensitization := ""
	respParamDesensitization := ""
	if conf != nil {
		description = conf.CategoryName
		reqParamDesensitization = conf.ReqParamDesensitization
		respParamDesensitization = conf.RespParamDesensitization
	}

	al := &model.APILog{
		CreatedAt:   time.Now(),
		API:         req.URL.String(),
		Operator:    userName,
		Description: description,
		Method:      req.Method,
		ReqBody:     desensitizationParam(reqBody, reqParamDesensitization),
		RemoteAddr:  req.Host,
		Status:      resbody.Status,
		Msg:         resbody.Message,
		Result:      desensitizationParam(utils.ToJsonString(resbody.Content), respParamDesensitization),
		Time:        cost,
	}

	if _, err := repo.SaveAPILog(al); err != nil {
		log.Error(err)
	}
}

//desensitizationParam 脱敏参数信息 desValues是需要脱敏的字段列表
func desensitizationParam(param, desValues string) string {
	res := param
	if desValues == "" || !strings.Contains(param, ",") {
		return param
	}

	for _, item := range strings.Split(param, ",") {
		if !strings.Contains(item, ":") {
			continue
		}

		k := strings.Split(item, ":")[0]
		v := strings.Split(item, ":")[1]

		if k == "" || v == "" {
			continue
		}

		k = strings.Trim(strings.Trim(strings.Trim(k, "["), "{"), "\n\t")
		v = strings.Trim(strings.Trim(v, "]"), "}")

		for _, des := range strings.Split(desValues, ",") {
			if strings.Trim(k, "\"") == des {
				res = strings.Replace(res, v, "\"*****\"", -1)
			}
		}

	}
	return res
}

//GetIDsFromResp 从resp当中拿到id
func GetIDsFromResp(r *http.Request, ww *Responser) (values []reflect.Value) {
	var respBody http3.RespBody
	log, _ := LoggerFromContext(r.Context())
	if err := json.Unmarshal([]byte(ww.buf.String()), &respBody); err != nil {
		log.Warnf("[OperateInterceptor]--> son unmarshal (%s) error:%s", ww.buf.String(), err.Error())
		return
	}
	// 这里拿到的id，竟然是float64?
	id := fmt.Sprintf("%g", respBody.Content["id"])
	if id == "" {
		log.Errorf("[OperateInterceptor]--> post method not return id, url is : %s", strings.TrimSpace(r.URL.String()))
		return
	}
	if strToUint(id) != nil {
		values = append(values, reflect.ValueOf(*strToUint(id)))
	}
	return
}

//GetIDsFromReq 从request当中拿到http的id
func GetIDsFromReq(r *http.Request, conf *RoutesConf) (values []reflect.Value) {
	//从URLParam当中拿不到id，为空
	//id := chi.URLParam(r, "id")

	//折中取法，若url当中不为id，则取不到
	id := strings.Replace(strings.TrimSpace(r.URL.String()), strings.Replace(conf.URL, "{id}", "", -1), "", -1)

	if id == "" {
		//若id为空，从queryParam当中再取一次
		id = r.URL.Query().Get("id")
	}
	uintID := strToUint(id)
	if uintID != nil {
		values = append(values, reflect.ValueOf(*uintID))
		return
	}

	//TODO 若Id还为空，则从body当中取
	if id == "" {

	}

	return
}

func strToUint(str string) *uint {
	// string to uint
	uintID, err := strconv.ParseUint(str, 10, 0)
	if err != nil {
		return nil
	}

	uid := uint(uintID)

	return &uid
}

//saveOperate 保存操作记录信息
func saveOperate(r *http.Request, conf *RoutesConf, source, destination string) {
	userName := "unknow user"
	repo, _ := RepoFromContext(r.Context())
	log, _ := LoggerFromContext(r.Context())
	user, ok := LoginUserFromContext(r.Context())

	if ok || user != nil {
		userName = user.LoginName
	}

	operate := &model.OperateLog{
		CreatedAt:    time.Now(),
		Operator:     userName,
		URL:          r.URL.String(),
		HTTPMethod:   r.Method,
		CategoryName: conf.CategoryName,
		CategoryCode: conf.CategoryCode,
		Source:       source,
		Destination:  destination,
	}

	if _, err := repo.SaveOperateLog(operate); err != nil {
		log.Errorf("[OperateInterceptor]--> save operate error, %s", err.Error())
	}
}

// isURLExist 判断URL是否存在
func isURLExist(url, method string, repo model.Repo, logger logger.Logger) *RoutesConf {
	//if len(routesConf) <= 0 {
	//	GetRoutesConf(repo, logger)
	//}

	GetRoutesConf(repo, logger)

	for _, conf := range routesConf {
		if strings.ToUpper(conf.HTTPMethod) != method {
			continue
		}

		if urlCompare(conf.URL, url) {
			return conf
		}
	}
	return nil
}

//urlCompare url对比
func urlCompare(url1, url2 string) bool {
	url1 = strings.Replace(url1, "/api/cloudboot/v1/", "", -1)
	url2 = strings.Replace(url2, "/api/cloudboot/v1/", "", -1)

	//去掉url2当中的?
	if strings.Contains(url2, "?") {
		url2 = url2[0:strings.Index(url2, "?")]
	}

	res := true
	if url1 == url2 {
		return true
	}

	u1 := strings.Split(url1, "/")
	u2 := strings.Split(url2, "/")

	if len(u1) != len(u2) {
		return false
	}

	for i, item := range u1 {
		if strings.Contains(item, "{") || strings.Contains(item, "}") {
			if item != u2[i] {
				//res = false
				//break
				//比如实际的SN:123456,而URL配置占位符{sn}，这个值不比较
				continue
			}
		}

		if strings.Contains(u2[i], "{") || strings.Contains(u2[i], "}") {
			continue
		}

		if item != u2[i] {
			res = false
			break
		}
	}

	return res
}

//GetRespVal 反射获取结果
func GetRespVal(repo model.Repo, method string, args []reflect.Value) interface{} {
	//getType := reflect.TypeOf(repo)
	getVal := reflect.ValueOf(repo)

	//fmt.Printf("type: %s\n", getType)
	// 获取interface value当中的method
	methodVal := getVal.MethodByName(method)

	respMap := map[string]interface{}{}

	// 反射调用func，并返回结果
	for _, resp := range methodVal.Call(args) {

		//fmt.Printf("type string %s\n", resp.Type().String())

		if resp.IsNil() && resp.Type().String() == "error" {
			continue
		}

		//fmt.Println(utils.ToJsonString(resp))

		respV := resp.Elem()

		//fmt.Println(respV)

		respT := respV.Type()

		for i := 0; i < respT.NumField(); i++ {
			if !resp.CanInterface() {
				continue
			}

			field := respT.Field(i)

			key := strings.Replace(string(respT.Field(i).Tag), "gorm:\"column:", "", -1)
			key = strings.Replace(key, "\"", "", -1)

			if key == "" {
				continue
			}

			respMap[key] = GetValue(field.Name, field.Type.Kind().String(), respV)

			//fmt.Printf("key:%s, value: %s\n", key, respV.FieldByName(respT.Field(i).Name))
		}

		//如何通过反射去获取这四个字段呢？
		respMap["updated_at"] = GetValue("UpdatedAt", "time", respV)
		respMap["created_at"] = GetValue("CreatedAt", "time", respV)
		respMap["deleted_at"] = GetValue("DeletedAt", "time", respV)
		respMap["id"] = GetValue("ID", reflect.Uint.String(), respV)
	}
	return respMap
}

//GetValue 转换field信息
func GetValue(key, typ string, val reflect.Value) (v interface{}) {
	str := fmt.Sprintf("%v", val.FieldByName(key))
	var err error

	switch typ {
	case reflect.Int.String():
		v, err = strconv.Atoi(str)
	case reflect.Uint.String():
		v = strToUint(str)
	case reflect.Bool.String():
		v, err = strconv.ParseBool(str)
	case reflect.Float32.String():
		v, err = strconv.ParseFloat(str, 32)
	case reflect.Float64.String():
		v, err = strconv.ParseFloat(str, 64)
	case "time":
		if strings.Contains(str, "CST") {
			v, err = time.Parse("2006-01-02 15:04:05 +0800 CST", str)
		}
	default:
		v = str
	}

	if err != nil {
		fmt.Printf("fail to convert str to %s, err: %s\n", typ, err.Error())
	}

	return
}
