package video

import (
	"encoding/json"
	"fmt"
	"gf-app/app/service/GreenSdk"
	"gf-app/app/service/oss"
	"gf-app/app/service/until"
	"gf-app/app/service/video"
	"gf-app/app/service/wechat"
	"gf-app/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	gconv "github.com/gogf/gf/util/gconv"
	"strconv"
	"strings"
	"time"
)

type TemRequest struct {
	Limit string `p:"limit" v:"required#条数必填"`
	Page  string `p:"page" v:"required#页数不能为空"`
}
type UploadValidity struct {
	TemplateId string `p:"template_id" v:"required#模板id不能为空"`
	Openid     string `p:"openid" v:"required#openid不能为空"`
	Sort     string `p:"sort" v:"required#sort不能为空"`
	Img        string `p:"img" `
	IsHead     string `p:"is_head" v:"required#is_head不能为空"`
}
type VideoValidity struct {
	TemplateId string `p:"template_id" v:"required#模板id不能为空"`
	Openid     string `p:"openId" v:"required#openid不能为空"`
	Code    string `p:"code" v:"required#code不能为空"`
	Title      string `p:"title" v:"required#标题不能为空"`
	Content    string `p:"content" v:"required#详情不能为空"`
	UserState     string `p:"user_state"`
	TemplateFunc     string `p:"template_func"`
	DefaultTitle     string `p:"default_title"`
	Author     string `p:"author"`

	//MusicId    string `p:"music_id" v:"required#音乐不能为空"`
}
type ReArr struct {
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}
type FlowValidity struct {
	//Openid     string `p:"openId" v:"required#openid不能为空"`
	Data  int         `json:"data"`
}
type DetailValidity struct {
	Id     string `p:"id" v:"required#id不能为空"`
}
type NoticeValidity struct {
	Openid     string `p:"openid" v:"required#openid不能为空"`
	VideoPath     string `p:"videoPath" v:"required#videoPath不能为空"`
	JobId     string `p:"jobId" v:"required#jobId不能为空"`
}
type QueryValidity struct {
	Openid     string `p:"openId" v:"required#openid不能为空"`
	VideoCode     string `p:"code" v:"required#code不能为空"`
}
type CodeValidity struct {
	Code     string `p:"code" v:"required#code不能为空"`
}
//获取模板列表
func Template(r *ghttp.Request) {
	var tem interface{}
	var data *TemRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	glog.Line().Info("小程序appid是"+until.GetAppid(r.Header["Referer"][0]))
	limit, _ := strconv.Atoi(data.Limit)
	page, _ := strconv.Atoi(data.Page)
	All := video.GetTemplate(page, limit)
	DataCount := video.GetTemCount()
	if All == nil {
		tem = make([]interface{}, 0)
		response.JsonExit(r, 200, "暂无数据", tem)
	}
	response.JsonExit(r, 200, "获取成功", ReArr{Count: DataCount, Data: All})
}

//获取音乐列表
func Music(r *ghttp.Request) {
	var tem interface{}
	var data *TemRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	limit, _ := strconv.Atoi(data.Limit)
	page, _ := strconv.Atoi(data.Page)
	All := video.GetMusicList(page, limit)
	DataCount := video.GetMusicCount()
	if All == nil {
		tem = make([]interface{}, 0)
		response.JsonExit(r, 200, "暂无数据", tem)
	}
	response.JsonExit(r, 200, "获取成功", ReArr{Count: DataCount, Data: All})
}

//上传图片
func Uploadimg(r *ghttp.Request) {
	var data *UploadValidity
	filetmp := "/tmp/"

	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	file := r.GetUploadFile("img")
	names, err := file.Save(filetmp, true)
	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	//
	localFileName := fmt.Sprintf("%s%s", filetmp, names)
	resultfile, err := oss.UploadFile(localFileName, names,until.GetAppid(r.Header["Referer"][0]))

	if err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	time := gtime.Now().Add(+time.Hour*2).Timestamp()
	newPath := fmt.Sprintf("/%s", resultfile)
	url := oss.GetPrivateUrl(strconv.FormatInt(time, 10), newPath)
	//鉴黄鉴暴
	reImg := GreenSdk.ImgSyncScan(url)
	if len(reImg) > 0 {
		response.JsonExit(r, 1, reImg)
	}
	if gconv.Int(data.IsHead) == 1 {
		isHeadString := GreenSdk.ImgIsHead(url)
		if len(isHeadString) > 0 {
			response.JsonExit(r, 1, isHeadString)
		}
	}

	//抠图api阿里云
	//reAvarImg := MattingSdk.GetMattingImg(url)
	//if len(reAvarImg) > 0 {
	//	response.JsonExit(r, 1, reAvarImg)
	//}
	//抠图api百度
	//_ = MattingSdk.GetBaiduMattingImg(url,names)

	//操作表
	//code, _ := video.InserVideoImg(data.Openid, gconv.Int(data.TemplateId), newPath,gconv.Int(data.Sort))
	code, _ := video.InserVideoImgToRedis(data.Openid, gconv.Int(data.TemplateId), newPath,gconv.Int(data.Sort))
	if code == "" {
		response.JsonExit(r, 1, "上传失败")
	}
	response.JsonExit(r, 200, "上传成功", g.Map{"url": url, "code": code})
}

//制作视频
func MakeVideoByImg(r *ghttp.Request) {
	var data *VideoValidity
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	defaultTitle := "我的影集标题"
	if  strings.Index(gconv.String(data.Title), defaultTitle) == 0 && len(data.Title) == len(defaultTitle){
		data.Title = data.DefaultTitle
	}
	if gconv.Int(g.Cfg("config").GetString("Tips.wordsSugSwitch")) == 1 {
		titleCheck := GreenSdk.TencenCloudWordsScan(data.Title,g.Cfg("config").GetString("Tips.titleSugTip"))
		if len(titleCheck) > 0 {
			response.JsonExit(r, 1, titleCheck)
		}
		contentCheck := GreenSdk.TencenCloudWordsScan(data.Content,g.Cfg("config").GetString("Tips.contentSugTip"))
		if len(contentCheck) > 0 {
			response.JsonExit(r, 1, contentCheck)
		}
		authorCheck := GreenSdk.TencenCloudWordsScan(data.Author,g.Cfg("config").GetString("Tips.authorSugTip"))
		if len(authorCheck) > 0 {
			response.JsonExit(r, 1, authorCheck)
		}
	}
	fmt.Println("标题2")
	//rerr := video.MakeVideo(data.Openid, data.Title, data.Content, gconv.Int(data.TemplateId), data.Code,gconv.Int(data.UserState))
	rerr := video.MakeVideoToAdmin(
		until.GetAppid(r.Header["Referer"][0]),
		data.Openid,
		data.Title,
		data.Content,
		gconv.Int(data.TemplateId),
		data.Code,
		gconv.Int(data.UserState),
		data.TemplateFunc,
		data.Author,
	)
	if rerr != nil {
		response.JsonExit(r, 1, gconv.String(rerr))
	}
	response.JsonExit(r, 200, "正在制作请等待")
}

func GetFlowList(r *ghttp.Request){
	//var data *FlowValidity
	//if err := r.Parse(&data); err != nil {
	//	response.JsonExit(r, 1, err.Error())
	//}
	openId := gconv.Map(r.Get("data"))["openId"]
	if openId == nil {
		response.JsonExit(r, 1, "openid不能为空")
	}
	reloadNum := video.GetOpenIdNumByOpenId(gconv.String(openId))
	twoNum := video.GetTwoNum()
	info := video.AllListGetDataByReloadNum(gconv.Int(reloadNum),gconv.Int(twoNum))
	cacheType := video.CACHE_HOST
	if gconv.Int(reloadNum) > 4 {
		// 取 邻核心数据
		cacheType = video.CACHE_NEAR
	}
	freq := [2]string{"4","4"}
	freq3 := [2]string{"3","3"}
	arr := [...]string{}

	res := g.Map{"data":info,"AD":g.Map{
		"AD" : g.Map{"hasAD1":"1|0|1|adunit-49aabe2c32d06bf3|||||首页流量主||0|0",
			"hasADcontinue":"1|0|9|adunit-0a8551401571fded|||||首页循环banner广告||0|0",
			"hasADv_1":"1|0|4|adunit-946d0b8516b86ba4|||||首页第四位banner广告||0|0",
		},
		"adPropBackOpenFlag": "0",
		"adPropState": g.Map{
			"broad" : "all",
			"freq" : freq,
		},
		"adPropStateBack":g.Map{
			"broad":"noclass",
			"freq" : freq,
		},
		"excitationShowState":g.Map{
			"broad":"noclass",
			"freq" : freq3,
		},
		"excitationTime":0,
		"money":arr,
		"payInfo":"0||||",
	},
	"firstNewId":cacheType,
	"navList" :g.Map{"0": "首页", "1": "祝福", "2": "正能量", "3": "搞笑", "4": "生活", "6": "其他"},
	"shareRecom":arr,
	"topVideo":arr,
	"_key":"0",
	"adSubscribeMessage":false,
	"Subscribe_tmplIds":arr,
	}

	newJson, _ := json.Marshal(res)

	r.Response.WriteJson(newJson)
	//response.JsonExit(r, 200, "获取成功",11)
}

func GetVideoDetail(r *ghttp.Request)  {
	//var data *DetailValidity
	//if err := r.Parse(&data); err != nil {
	//	response.JsonExit(r, 1, err.Error())
	//}
	id := r.Get("id")
	//userType := r.Get("type")
	openId := r.Get("openId")
	if id == nil {
		response.JsonExit(r, 1, "id不能为空")
	}
	//if userType == "" {
	//	userType = 0
	//}
	list := video.GetVideo(gconv.String(id),gconv.String(openId))
	freq := [2]string{"4","4"}
	freq3 := [2]string{"3","3"}
	arr := [...]string{}
	newArr := gconv.Map(list)
	newArr["ADlist"] = g.Map{
		"AD" : g.Map{"hasAD1":"1|0|1|adunit-49aabe2c32d06bf3|||||首页流量主||0|0",
			"hasADGroup":"1||3|||http://zcom.xazzp.com/adapp/wx87fdc5cbebc4f2e2/addGroup.png|||首页进群卡片|https://mp.weixin.qq.com/s/hL2WsWFt4RBQM9ZN1_HO8A|0|0",
			"hasADcontinue":"1|0|9|adunit-0a8551401571fded|||||首页循环banner广告||0|0",
			"hasADv_1":"1|0|4|adunit-946d0b8516b86ba4|||||首页第四位banner广告||0|0",
		},
		"adPropBackOpenFlag": "0",
		"adPropState": g.Map{
			"broad" : "all",
			"freq" : freq,
		},
		"adPropStateBack":g.Map{
			"broad":"noclass",
			"freq" : freq,
		},
		"excitationShowState":g.Map{
			"broad":"noclass",
			"freq" : freq3,
		},
		"excitationTime":0,
		"money":arr,
		"payInfo":"0||||",
	}
	newArr["recommendList"] = arr
	newArr["maskRemVideoOpenFlag"] = arr
	newArr["maskList"] = arr
	newArr["btnShareText"] = "点这里分享，让更多的朋友看看"
	newArr["videoType"] = 1
	newArr["explrecommList"] = arr


	response.JsonExit(r, 200, "获取成功",newArr)
}
func MakeComplete(r *ghttp.Request)  {
	var data *NoticeValidity
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	hash := make(map[string]interface{})
	hash["openid"] = data.Openid
	hash["jobId"] = data.JobId
	hash["videoPath"] = data.VideoPath
	//_,err := g.Redis().Do("HMSET",redis.Args{}.Add("video.ceshi").AddFlat(hash)...)
	//if err != nil {
	//	fmt.Println("报错了")
	//}
	err := video.SendUserComplate(data.JobId,data.VideoPath,data.Openid)
	if err != nil {
		response.JsonExit(r, 1, "发送失败",err)
	}
	response.JsonExit(r, 200, "发送成功")
}
func QueryVideoState(r *ghttp.Request)  {
	var data *QueryValidity
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	videoInfoMap := video.QueryVideoState(data.VideoCode,data.Openid)
	if len(videoInfoMap) >0 {
		response.JsonExit(r, 200, "获取成功",videoInfoMap)
	}
	percent := video.QueryVideoPercent(data.VideoCode)
	response.JsonExit(r, 200, "正在合成中",g.Map{"state":1,"percent":percent})
}
func GetOpenidByCode(r *ghttp.Request)  {
	//var data *CodeValidity
	//if err := r.Parse(&data); err != nil {
	//	response.JsonExit(r, 1, err.Error())
	//}
	data := gconv.Map(r.Get("data"))
	openId := wechat.GetOpenidByCode(gconv.String(data["code"]))
	r.Response.Write(openId)
}