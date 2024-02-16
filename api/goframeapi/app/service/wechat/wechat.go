package wechat

import (
	"encoding/json"
	"fmt"
	"gf-app/app/service/oss"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

type CustomServiceMsg struct {
	ToUser  string         `json:"touser"`
	TemplateId string         `json:"template_id"`
	Data    TextMsgContent `json:"data"`
}
type TextMsgContent struct {
	Thing5 Thing5Type `json:"thing5"`
	Time3 Time3Type `json:"time3"`
	Thing4 Thing4Type `json:"thing4"`
}
type Thing5Type struct {
	Value string `json:"value"`
}
type Time3Type struct {
	Value string `json:"value"`
}
type Thing4Type struct {
	Value string `json:"value"`
}
const ACCESS_TOKEN = "easywechat.kernel.access_token.goNew"
const SUBSCRIBEMESSAGE = "Fh8JO0jNTUIM90CUr6QF6u6wpW-t5VY-eJ_EXhm1TjU"
const VIDEO_QUEUE_SUCCESS = "send.message.success"
const VIDEO_QUEUE_FAIL = "send.message.fail"
func SendUserMessage(openid string) interface{} {
	var token string
	var tokenMap map[string]interface{}
	appid := g.Cfg("config").GetString("AppInfo.AppId")
	appSecret := g.Cfg("config").GetString("AppInfo.AppSecret")

	r,err := g.Redis().DoVar("GET",ACCESS_TOKEN)
	if err != nil {
		fmt.Println(err)
		fmt.Println("24")
	}
	rerr := json.Unmarshal(gconv.Bytes(r),&tokenMap)
	if rerr != nil {
		fmt.Println(rerr)
		fmt.Println("29")
	}
	token = gconv.String(tokenMap["access_token"])

	if gconv.String(r) == "" {
		token = GetAccessToken(appid,appSecret)
		fmt.Println(token)
		fmt.Println("重新获取token")
	}
	fmt.Println("发送")
	fmt.Println(token)
	wxSendUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s",token)
	msg := "\U0001F449" + "请点击查看" + "\U0001F448"
	json1 := &CustomServiceMsg{
		ToUser:  openid,
		TemplateId: SUBSCRIBEMESSAGE,
		Data:    TextMsgContent{
			Thing5: Thing5Type{Value:"您的视频制作已完成"},
			Time3: Time3Type{Value:gtime.Now().Format("Y-m-d H:i:s")},
			Thing4: Thing4Type{Value:msg},
		},
	}
	body, err := json.MarshalIndent(json1, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	r1, err1 := g.Client().Post(wxSendUrl,body)
	fmt.Println(r1.ReadAllString())
	if err1 != nil {
		panic(err1)
		fmt.Println("发送失败")
		g.Redis().Do("RPUSH",VIDEO_QUEUE_FAIL,fmt.Sprintf("{\"发送失败\":\"%s\"}",err1))
		g.Redis().Do("expire",VIDEO_QUEUE_FAIL,86400*7)
	} else {
		defer r1.Close()
		suMap := gconv.Map(r1.ReadAllString())
		fmt.Println("123")
		if gconv.Int(suMap["errcode"]) == 0 && gconv.String(suMap["errmsg"])=="ok" {
			fmt.Println("发送成功")
			json,_ := json.Marshal(suMap)
			g.Redis().Do("RPUSH",VIDEO_QUEUE_SUCCESS,json)
			g.Redis().Do("expire",VIDEO_QUEUE_SUCCESS,86400*7)
		}else{
			fmt.Println("发送失败")
			json,_ := json.Marshal(suMap)
			g.Redis().Do("RPUSH",VIDEO_QUEUE_FAIL,fmt.Sprintf("{\"发送失败\":\"%s\"}",json))
			g.Redis().Do("expire",VIDEO_QUEUE_FAIL,86400*7)
		}

	}
	fmt.Println("发送完成")
	return nil
}
func GetAccessToken(appid string, appSecret string) string {
	var token string
	getTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",appid,appSecret)
	if r, err := g.Client().Get(getTokenUrl); err != nil {
		panic(err)
		return ""
	} else {
		defer r.Close()
		var tokenArr = make(map[string]interface{})
		tokenArr = gconv.Map(r.ReadAllString())
		if len(tokenArr) == 0 {
			fmt.Println("token空的")
			return ""
		}
		token = gconv.String(tokenArr["access_token"])
		jsonToken,_ := json.Marshal(tokenArr)
		g.Redis().Do("SET",ACCESS_TOKEN,jsonToken)
		g.Redis().Do("expire",ACCESS_TOKEN,7000)
	}
	return token
}
func GetOpenidByCode(code string) string {
	appid := g.Cfg("config").GetString("AppInfo.AppId")
	appSecret := g.Cfg("config").GetString("AppInfo.AppSecret")
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	if r, err := g.Client().Get(fmt.Sprintf(url,appid,appSecret,code)); err != nil {
		panic(err)
	} else {
		defer r.Close()
		suMap := gconv.Map(r.ReadAllString())
		return gconv.String(suMap["openid"])
	}
	return ""
}
func GetRediesKeyAccessToken(grant string,appid string,appSecret string) string {
	msg := fmt.Sprintf("{\"grant_type\":\"client_credential\",\"appid\":\"%s\",\"secret\":\"%s\"}",g.Cfg("config").GetString("AppInfo.AppId"),g.Cfg("config").GetString("AppInfo.AppSecret"))
	has,_ := oss.MD5(msg)
	return fmt.Sprintf(ACCESS_TOKEN,has)
}


