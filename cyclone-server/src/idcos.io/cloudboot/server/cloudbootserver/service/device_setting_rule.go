package service

import (
	"fmt"
	"strings"
	"reflect"
	"net/http"
	"encoding/json"

	"github.com/voidint/binding"
	"github.com/voidint/page"
	"idcos.io/cloudboot/middleware"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// 事实对象的属性（与规则库所需的属性匹配）
type DeviceFact struct {
	SN                 string    // 设备SN
	Category           string    // 设备类型
	PhysicalArea       string    // 物理区域， 如： ServerFarm bonding区1
	Arch               string    // cpu架构：x86\arm
	Vendor             string    // 设备厂商
	IsFITIEcoProduct   string      // 是否金融信创生态产品
}

// 推理规则前件是否包含逻辑运算符号AND OR
type ConditionHasLogicalOperator struct {
	HasAND    bool
	HasOR     bool
	HasBoth   bool
}

// 规则推理返回结构体
type ActionResult struct {
	RuleType        string    // os raid network
	OSResult        string
	RaidResult      string
	NetworkResult   string
}

// 元规则推理，仅针对单个 RuleP 结构体进行推理，返回 bool
func MetaRuleInduction(log logger.Logger, fact *DeviceFact, rulep *model.RuleP) (bool) {
	if rulep.Attribute != "" && rulep.Attribute != model.AttributeLogicalOperator {
		var factvalue string
		// 获取fact对应的属性值
		switch rulep.Attribute {
		case model.AttributeCategory:
			factvalue = fact.Category
		case model.AttributePhysicalArea:
			factvalue = fact.PhysicalArea
		case model.AttributeArch:
			factvalue = fact.Arch
		case model.AttributeVendor:
			factvalue = fact.Vendor
		case model.AttributeIsFITIEcoProduct:
			factvalue = fact.IsFITIEcoProduct			
		default:
			log.Errorf("Attribute of rule is not defined.")
			return false
		}
		
		// 逻辑处理 equal contains in
		switch rulep.Operator {
		case model.OperatorEqual:
			if factvalue == rulep.Value[0] {
				log.Infof("Device %s  Meta_Rule_Hit: %s equal %s ", fact.SN, factvalue, rulep.Value)
				return true
			}
		case model.OperatorContains:
			if strings.Contains(factvalue, rulep.Value[0]) {
				log.Infof("Device %s  Meta_Rule_Hit: %s contains %s ", fact.SN, factvalue, rulep.Value)
				return true
			}
		case model.OperatorIN:
			for _, each := range rulep.Value {
				if factvalue == each {
					log.Infof("Device %s  Meta_Rule_Hit: %s in %s ", fact.SN, factvalue, rulep.Value)
					return true
				}
			}
		default:
			log.Errorf("Operator of rule is not defined.")
			return false
		}
	}
    return false
}

// 前件解析器
func ConditionParser(log logger.Logger, repo model.Repo, fact *DeviceFact, rule *model.DeviceSettingRule) (bool, error) {
	if rule.Condition != "" {
		var ruleconditions []model.RuleP
		chlo := ConditionHasLogicalOperator{
			HasAND:  false,
			HasOR:   false,
			HasBoth: false,
		}

		// 解析规则前件
		err := json.Unmarshal([]byte(rule.Condition), &ruleconditions)
		if err != nil {
			log.Error(err)
			return false,err
		}
        // fact 与 fact 之间必须通过逻辑运算符连接，故index=奇数可增加校验是否为AND OR
		log.Infof("Condition of rule %d : %s", rule.ID, ruleconditions)

		// 推理是否包含逻辑运算符 AND OR
		if len(ruleconditions) > 1 {
		    for _, rc := range ruleconditions {
		    	if rc.Attribute == model.AttributeLogicalOperator && rc.Operator == model.OperatorEqual {
		    		switch rc.Value[0] {
		    		case model.OperatorOR:
		    			chlo.HasOR = true
		    		case model.OperatorAND:
		    			chlo.HasAND = true
		    		}
		    	}
		    }
		    if chlo.HasOR && chlo.HasAND {
		    	chlo.HasBoth = true
		    }
	    } else if len(ruleconditions) == 1 {
			for _, rc := range ruleconditions {
				return MetaRuleInduction(log, fact, &rc),nil
			}			
		}
		
		// 前件包含AND + OR 优先处理 OR
		if chlo.HasBoth {
			// 针对包含AND + OR 临时存储结果 [true,and,false,or,true]
        	var logicresult []string			
			for _, rc := range ruleconditions {
				if rc.Attribute == model.AttributeLogicalOperator {
		    		switch rc.Value[0] {
		    		case model.OperatorOR:
		    			logicresult = append(logicresult, "or")
		    		case model.OperatorAND:
						logicresult = append(logicresult, "and")
		    		}
		    	} else {
					if MetaRuleInduction(log, fact, &rc) {
						logicresult = append(logicresult, "true")
					} else {
						logicresult = append(logicresult, "false")
					}
				}
			}
			
			// 逻辑结果数组按and拆分 [[true],[false,or,true]]
			var logicresult_split_by_and [][]string
			i := 0
			for _, v := range logicresult {
				if v == "and" {
					i++
				} else {
					logicresult_split_by_and[i] = append(logicresult_split_by_and[i], v)
				}
			}

			// 处理返回最终bool
			contain_true := false
			for _, splitted_list := range logicresult_split_by_and {
				if len(splitted_list) == 1 {
					if splitted_list[0] == "false"{
						return false,nil
					} else {
						continue
					}
				} else {
					for _, va := range splitted_list {
						if va == "true" {
							contain_true = true
						} else {
							continue
						}
					}
					if contain_true == false {
						return false,nil
					}
				}
			}
			return true, nil
		} else if chlo.HasAND {
			for _, rc := range ruleconditions {
				if rc.Attribute == model.AttributeLogicalOperator {
		    		switch rc.Value[0] {
		    		case model.OperatorOR:
						log.Errorf("Not expected OR of rule %d ", rule.ID)
						return false, fmt.Errorf("Not expected OR of rule %d ", rule.ID)
		    		case model.OperatorAND:
						continue
		    		}
		    	} else {
					if MetaRuleInduction(log, fact, &rc) {
						continue
					} else {
						return false, nil
					}
				}
			}
			return true, nil
		} else if chlo.HasOR {
			for _, rc := range ruleconditions {
				if rc.Attribute == model.AttributeLogicalOperator {
		    		switch rc.Value[0] {
					case model.OperatorOR:
						continue
		    		case model.OperatorAND:
						log.Errorf("Not expected AND of rule %d ", rule.ID)
						return false, fmt.Errorf("Not expected AND of rule %d ", rule.ID)
		    		}
		    	} else {
					if MetaRuleInduction(log, fact, &rc) {
						return true, nil
					} else {
						continue
					}
				}
			}
			return false, nil
		}
	} 
	log.Errorf("Empty of rule %d ", rule.ID)
	return false, fmt.Errorf("Empty of rule %d ", rule.ID)
}


// 规则推理机（仅推理某一种类）
func RuleTypeInduction(log logger.Logger, repo model.Repo, fact *DeviceFact, ruletype string) (*ActionResult, error) {
	// 根据规则分类提取规则进行处理
	var result = &ActionResult{}

	switch ruletype {
	case model.DeviceSettingRuleOS:
		ruleItems, err := repo.GetDeviceSettingRulesByType(ruletype)
		if err != nil {
			log.Errorf("Fact(device sn: %s ) failed to get os rules.", fact.SN)
			return nil, err
		}
		for _, rule := range ruleItems {
			b, err := ConditionParser(log, repo, fact, rule)
			if err != nil {
				return nil, err
			}
			if b {
				log.Infof("Fact(device sn: %s ) successed to get os result.", fact.SN)
				result.OSResult = rule.Action
				return result, nil
			}
		}

	case model.DeviceSettingRuleRaid:
		ruleItems, err := repo.GetDeviceSettingRulesByType(ruletype)
		if err != nil {
			log.Errorf("Fact(device sn: %s ) failed to get raid rules.", fact.SN)
			return nil, err
		}
		for _, rule := range ruleItems {
			b, err := ConditionParser(log, repo, fact, rule)
			if err != nil {
				return nil, err
			}
			if b {
				log.Infof("Fact(device sn: %s ) successed to get raid result.", fact.SN)
				result.RaidResult = rule.Action
				return result, nil
			}
		}
	
	case model.DeviceSettingRuleNetwork:
		ruleItems, err := repo.GetDeviceSettingRulesByType(ruletype)
		if err != nil {
			log.Errorf("Fact(device sn: %s ) failed to get network rules.", fact.SN)
			return nil, err
		}
		for _, rule := range ruleItems {
			b, err := ConditionParser(log, repo, fact, rule)
			if err != nil {
				return nil, err
			}
			if b {
				log.Infof("Fact(device sn: %s ) successed to get network result.", fact.SN)
				result.NetworkResult = rule.Action
				return result, nil
			}
		}
	default:
		log.Errorf("Not supported type(%s) of rules.", ruletype)
		return nil, fmt.Errorf("Not supported type(%s) of rules.", ruletype)
	}
	log.Errorf("Fact(device sn: %s ) failed to get type(%s) of rules.", fact.SN, ruletype)
	return nil, fmt.Errorf("Fact(device sn: %s ) failed to get type(%s) of rules.", fact.SN, ruletype)
}

// 规则推理机
func RuleInduction(log logger.Logger, repo model.Repo, fact *DeviceFact) (*ActionResult, error) {
	// 根据规则分类提取规则进行处理
	var result = &ActionResult{}
	var ruletypes = []string {model.DeviceSettingRuleOS, model.DeviceSettingRuleRaid, model.DeviceSettingRuleNetwork}
	for _, ruletype := range ruletypes {
		switch ruletype {
		case model.DeviceSettingRuleOS:
			ruleItems, err := repo.GetDeviceSettingRulesByType(ruletype)
			if err != nil {
				log.Errorf("Fact(device sn: %s ) failed to get os rules.", fact.SN)
				return nil, err
			}
			for _, rule := range ruleItems {
				b, err := ConditionParser(log, repo, fact, rule)
				if err != nil {
					return nil, err
				}
				if b {
					log.Infof("Fact(device sn: %s ) successed to get os result.", fact.SN)
					result.OSResult = rule.Action
				}
			}
	
		case model.DeviceSettingRuleRaid:
			ruleItems, err := repo.GetDeviceSettingRulesByType(ruletype)
			if err != nil {
				log.Errorf("Fact(device sn: %s ) failed to get raid rules.", fact.SN)
				return nil, err
			}
			for _, rule := range ruleItems {
				b, err := ConditionParser(log, repo, fact, rule)
				if err != nil {
					return nil, err
				}
				if b {
					log.Infof("Fact(device sn: %s ) successed to get raid result.", fact.SN)
					result.RaidResult = rule.Action
				}
			}
		
		case model.DeviceSettingRuleNetwork:
			ruleItems, err := repo.GetDeviceSettingRulesByType(ruletype)
			if err != nil {
				log.Errorf("Fact(device sn: %s ) failed to get network rules.", fact.SN)
				return nil, err
			}
			for _, rule := range ruleItems {
				b, err := ConditionParser(log, repo, fact, rule)
				if err != nil {
					return nil, err
				}
				if b {
					log.Infof("Fact(device sn: %s ) successed to get network result.", fact.SN)
					result.NetworkResult = rule.Action
				}
			}
		}
	}
	// 仅返回完整的参数推理结果
	if result.OSResult == "" || result.RaidResult == "" || result.NetworkResult == "" {
		log.Errorf("Fact(device sn: %s ) do not hit all types of rules.", fact.SN)
		return nil, fmt.Errorf("Fact(device sn: %s ) do not hit all types of rules.", fact.SN)
	}
	return result, nil
}


// GetDeviceSettingRulePageReq 获取规则记录表分页请求参数
type GetDeviceSettingRulePageReq struct {
	RuleCategory        string      `json:"rule_category"`
	Page    	 		int64    	`json:"page"`
	PageSize 			int64    	`json:"page_size"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *GetDeviceSettingRulePageReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.RuleCategory: "rule_category",
		&reqData.Page:         "page",
		&reqData.PageSize:     "page_size",
	}
}

// GetDeviceSettingRulePage 获取规则记录分页
func GetDeviceSettingRulePage(log logger.Logger, repo model.Repo, reqData *GetDeviceSettingRulePageReq) (*page.Page, error) {
	if reqData.PageSize <= 0 || reqData.PageSize > 100 {
		reqData.PageSize = 20
	}
	if reqData.Page < 0 {
		reqData.Page = 0
	}

	cond := model.DeviceSettingRule{
		RuleCategory: reqData.RuleCategory,
	}

	totalRecords, err := repo.CountDeviceSettingRules(&cond)
	if err != nil {
		return nil, err
	}

	pager := page.NewPager(reflect.TypeOf(&DeviceSettingRuleResp{}), reqData.Page, reqData.PageSize, totalRecords)
	items, err := repo.GetDeviceSettingRules(&cond, model.OneOrderBy("id", model.DESC), pager.BuildLimiter())
	if err != nil {
		return nil, err
	}
	for i := range items {
		item, err := convert2DeviceSettingRulesResult(log, repo, items[i])
		if err != nil {
			return nil, err
		}
		if item != nil {
			pager.AddRecords(item)
		}
	}
	return pager.BuildPage(), nil

}

// 规则表分页查询信息
type DeviceSettingRuleResp struct {
	ID        				 uint      			`json:"id"`
	//Condition                []model.RuleP      `json:"condition"`      // 前件
	Condition                string		        `json:"condition"`      // 前件
	Action                   string    			`json:"action"`         // 结论
	RuleCategory             string    			`json:"rule_category"`  // 规则分类 enum('os','network','raid')	
}

// 返回分页数据结构
func convert2DeviceSettingRulesResult(log logger.Logger, repo model.Repo, mod *model.DeviceSettingRule) (*DeviceSettingRuleResp, error) {
	if mod == nil {
		return nil, nil
	}
	// 解析规则前件
	//var ruleconditions []model.RuleP
	//err := json.Unmarshal([]byte(mod.Condition), &ruleconditions)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	result := DeviceSettingRuleResp{
		ID:        		mod.ID,
		Condition:		mod.Condition,
		Action:			mod.Action,
		RuleCategory:   mod.RuleCategory,
	}
	return &result, nil
}

//SaveDeviceSettingRuleReq 保存规则请求参数
type SaveDeviceSettingRuleReq struct {
	//DeviceSettingRuleBase
	ID uint `json:"id"`
	// 用户登录名
	LoginName string `json:"-"`
	Action                   string    			`json:"action"`         // 结论
	RuleCategory             string    			`json:"rule_category"`  // 规则分类 enum('os','network','raid')	
	Condition                string      		`json:"condition"`      // 前件
}

// 基本字段
//type DeviceSettingRuleBase struct {
//	Condition                []model.RuleP      `json:"condition"`      // 前件
//	Action                   string    			`json:"action"`         // 结论
//	RuleCategory             string    			`json:"rule_category"`  // 规则分类 enum('os','network','raid')	
//}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *SaveDeviceSettingRuleReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.ID:       			"id",
		&reqData.Condition: 		"condition",
		&reqData.Action: 			"action",
		&reqData.RuleCategory:    	"rule_category",
	}
}

// Validate 结构体数据校验
func (reqData *SaveDeviceSettingRuleReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	//必须的基本数据不能为空
	if errs = reqData.baseValidate(req, errs); errs != nil {
		return errs
	}

	//更新规则时，校验指定ID的规则是否存在
	if reqData.ID > 0 {
		if _, err := repo.GetDeviceSettingRuleByID(reqData.ID); errs != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("查询规则id(%d)出现错误: %s", reqData.ID, err.Error()))
			return errs
		}
	} 
	return nil
}

//baseValidate 必要参数不能为空
func (reqData *SaveDeviceSettingRuleReq) baseValidate(req *http.Request, errs binding.Errors) binding.Errors {
	if reqData.Condition == "" {
		errs.Add([]string{"condition"}, binding.RequiredError, "规则前件不能为空")
		return errs
	}
	if reqData.Action == "" {
		errs.Add([]string{"action"}, binding.RequiredError, "规则推论不能为空")
		return errs
	}
	if reqData.RuleCategory == "" {
		errs.Add([]string{"rule_category"}, binding.RequiredError, "规则类别不能为空")
		return errs
	}	
	return errs
}

//SaveDeviceSettingRule 保存规则记录
func SaveDeviceSettingRule(log logger.Logger, repo model.Repo, reqData *SaveDeviceSettingRuleReq) error {
	
	sr := model.DeviceSettingRule{
		Condition: 			reqData.Condition,
		Action: 			reqData.Action,
		RuleCategory:    	reqData.RuleCategory,
	}
	sr.Model.ID = reqData.ID

	_, err := repo.SaveDeviceSettingRule(&sr)
	if err != nil {
		return err
	}

	reqData.ID = sr.Model.ID
	return err
}

// 根据ID删除规则记录
type DelDeviceSettingRuleReq struct {
	IDs []uint `json:"ids"`
}

// FieldMap 请求参数与结构体字段建立映射
func (reqData *DelDeviceSettingRuleReq) FieldMap(req *http.Request) binding.FieldMap {
	// 此处只做字段映射，不要包含校验逻辑。
	return binding.FieldMap{
		&reqData.IDs: "ids",
	}
}

// Validate 结构体数据校验
func (reqData *DelDeviceSettingRuleReq) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	repo, _ := middleware.RepoFromContext(req.Context())
	for _, id := range reqData.IDs {
		if _, err := repo.GetDeviceSettingRuleByID(id); err != nil {
			errs.Add([]string{"id"}, binding.RequiredError, fmt.Sprintf("规则id(%d)不存在", id))
			return errs
		}
	}
	return nil
}

//RemoveDeviceSettingRules 删除指定ID的规则记录
func RemoveDeviceSettingRules(log logger.Logger, repo model.Repo, reqData *DelDeviceSettingRuleReq) (affected int64, err error) {
	for _, id := range reqData.IDs {
		_, err := repo.RemoveDeviceSettingRuleByID(id)
		if err != nil {
			log.Errorf("delete device setting rule (id=%d) fail,err:%v", id, err)
			return affected, err
		}
		affected++
	}
	return affected, err
}