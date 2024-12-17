package middleware

import (
	"errors"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"strings"
	"strconv"
	"net/http"

	"idcos.io/cloudboot/config"
	"idcos.io/cloudboot/logger"
	"idcos.io/cloudboot/model"
)

// ESB - Enterprise Service Bus 企业服务总线功能
//URL请求：
//
//http://xxx.xxx.xxx.xxx/service/[serviceCode]/[apiCode]?[加密参数]
//
//serviceCode：服务代码；apiCode：api代码
//服务和api的关系是一个服务包含多个api，如统一消息服务，代码为message；里面包含多个api，mailsend（邮件发送），wechatsend（邮件发送）等。
//请求方式：支持Get、Post、Put、Delete
//响应格式支持：Xml、Json
//加密参数是api访问鉴权

//关于message/mailsend API的参数说明
//参数名	必填	说明
//From	Y	发件人邮箱
//To	N	收件人邮箱，多人用分号隔开。（to，cc，bcc不能全部同时为空）
//CC	N	抄送人邮箱，多人用分号隔开。（to，cc，bcc不能全部同时为空）
//BCC	N	暗送人邮箱，多人用分号隔开。（to，cc，bcc不能全部同时为空）
//Title	Y	邮件标题
//Content	Y	邮件内容,如果带图片的话，要配合下面的Attachments参数使用。图片：<img src='cid:下面的ContentId' />
//EmailType	N	邮件类型，可选值0（普通邮件），1（会议邮件）。默认为0。
//Priority	N	邮件优先级，0普通，1高优先级。默认值为0。
//BodyFormat	N	邮件格式，0 文本、1 Html。默认值为0。
//Appt_Location	N	当邮件为约会邮件时，约会地点
//Appt_Organizer	N	当邮件为约会邮件时，约会组织者
//Appt_StartTime	N	当邮件为约会邮件时，约会开始时间；时间格式为yyyyMMddHHmmss，如时间2016/6/2 11:17:09表现为20160602111709
//Appt_EndTime	N	当邮件为约会邮件时，约会结束时间；时间格式跟Appt_StartTime相同。
//Attachments	N	附件，为Json格式的字符串（请传入字符串）。"[{Name:文件名1,Data:文件内容1},{Name:文件名2,Data:文件内容2}]" 文件名：就是显示在邮件的文件名称；文件内容是把文件内容读成byte使用base64编码的字符串。
//MailToInternet	N	判断邮件是否需要发送到外网, 0否（默认）、1是。这个设置需要开启ESB后台App管理字段：是否允许发送外网下拉选项，来配合使用。【无法重发】
//AuthUser	N	默认是smtp无密码验证发送，如果有需要指定用户名可以设置AuthUser。但必须AuthPassword不为空才有效。如果AuthUser为空有需要密码校验的话，AuthUser默认为发件人From。 【无法重发】
//AuthPassword	N	默认是smtp无密码验证发送，设置了的话就相当于需要smtp密码发送。【无法重发】

// SendMailReq 调用邮件发送API的通用请求结构体
type SendMailReq struct {
	From			string				`json:"from"`
	To				string				`json:"to"`
	CC				string				`json:"cc"`
	Title			string				`json:"title"`
	Content 		string				`json:"content"`
	BodyFormat		string				`json:"bodyformat"`
	Priority		string				`json:"priority"`
	Attachments		[]MailAttachment	`json:"attachments"`
}

// 邮件附件结构体
type MailAttachment struct {
	Name	string		`json:"name"`
	Data 	interface{}	`json:"data"`
}

// SendMailResp 调用邮件发送API的通用请求的返回结构体
type SendMailResp struct {
	//返回码
	Code uint 
	//返回信息
	Message string
	//返回结果内容，也是一个key-value  Result:{"MessageId":"XXX"}
	Result map[string]string
}

// 生成32位MD5
func MD5(text string) string{
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

// 触发邮件发送
func SendMail(log logger.Logger, repo model.Repo, conf *config.Config, sendmailreq *SendMailReq) (err error) {
	//参数	说明
	//ESB为应用生成的appid apptoken
	appid := conf.ExternalService.ESBAppID
	apptoken := conf.ExternalService.ESBAppToken
	//nonce	5位随机数字
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	nonce := strconv.FormatInt(r.Int63(),10)[:5]
	//timestamp	Unix的时间戳，即当前时间的utc距离1970年1月1日0点整的秒数。如2016/8/31 0:00:00，转化后值为 1472572800
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	//signature	加密签名。算法为 MD5(MD5(AppId + Nonce + Timestamp) + appToken); appToken为指派的appToken。 
	//注意:AppId + Nonce + Timestamp是字符串相加，MD5中算法的编码要用UTF-8，全部都是小写，比如：MD5(MD5(AppId + Nonce + Timestamp).ToLower() + appToken)）.ToLower()
	signature := strings.ToLower(MD5(strings.ToLower(MD5(appid + nonce + timestamp)) + apptoken))

	netSendMailURL := fmt.Sprintf("%s/service/message/mailsend?appid=%s&nonce=%s&timestamp=%s&signature=%s", conf.ExternalService.ESBBaseURL, 
	appid, nonce, timestamp, signature)
	// 参数校验
	if sendmailreq.From == "" {
		log.Error("发件人邮箱不允许为空")
		return errors.New("发件人邮箱不允许为空")
	}
	if sendmailreq.Title == "" {
		log.Error("邮件标题不允许为空")
		return errors.New("邮件标题不允许为空")
	}
	if sendmailreq.Content == "" {
		log.Error("邮件正文不允许为空")
		return errors.New("邮件标题不允许为空")
	}	
	if sendmailreq.To == "" && sendmailreq.CC == "" {
		log.Error("收件人邮箱与抄送人邮箱不允许同时为空")
		return errors.New("收件人邮箱与抄送人邮箱不允许同时为空")		
	}
	reqBody, err := json.Marshal(sendmailreq)
	if err != nil {
		log.Error(err)
		return err
	}
	//log.Debugf("POST %s, request body: %s", netSendMailURL, reqBody)
	req, err := http.NewRequest("POST", netSendMailURL, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Error(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	defer res.Body.Close()
	
	respBytes, err := ioutil.ReadAll(res.Body)
	sendmailresp := new(SendMailResp)
	if err = json.Unmarshal(respBytes, &sendmailresp); err != nil {
		log.Errorf("unmarshal /service/message/mailsend resp err: %s", err.Error())
		return fmt.Errorf("解析json编码数据错误:%s", err.Error())
	}
	if sendmailresp.Code != 0 {
		log.Error("mailsend error: Code != 0")
		log.Errorf("mailsend error detail: %s", sendmailresp.Message)
		return fmt.Errorf("调用邮件发送接口错误:%s", "Code != 0")
	}
	log.Infof("mailsend success, Message: %s", sendmailresp.Message)
	return nil
}