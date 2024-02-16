package GreenSdk

import (
	"encoding/json"
	"fmt"
	"gf-app/app/service/pigoModel"
	gohash "github.com/corona10/goimagehash"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"image"
	"net/http"
)
type Newresult struct{
	Code int
	Data []Results1
	Msg string		`json:"msg"`
	RequestId string		`json:"requestId"`
}
type Results1 struct {
	Code int		`json:"code"`
	DataId string		`json:"dataId"`
	Msg string		`json:"msg"`
	TaskId string		`json:"taskId"`
	Url string		`json:"url"`
	Results []Sug		`json:"results"`
}
type Sug struct{
	Label string		`json:"label"`
	Scene string		`json:"scene"`
	Suggestion string		`json:"suggestion"`
}
func ImgSyncScan(imgUrl string) string {
	accessKeyId := g.Cfg("config").GetString("aliyunoss.AccessKeyId")
	accessKeySecret := g.Cfg("config").GetString("aliyunoss.AccessKeySecret")
	profile := Profile{AccessKeyId:accessKeyId, AccessKeySecret:accessKeySecret}

	path := "/green/image/scan"

	clientInfo := ClinetInfo{Ip:"127.0.0.1"}

	// 构造请求数据
	bizType := "Green"
	scenes := []string{"porn","terrorism"}

	task := Task{DataId:Rand().Hex(), Url:imgUrl}
	tasks := []Task{task}

	bizData := BizData{ bizType, scenes, tasks}

	var client IAliYunClient = DefaultClient{Profile:profile}

	// your biz code
	response := client.GetResponse(path, clientInfo, bizData)

	var revMsg Newresult
	err := json.Unmarshal([]byte(response),&revMsg)
	if err != nil {
		return err.Error()
	}
	if revMsg.Code == 200 {
		results := revMsg.Data[0].Results
		for _,v := range results{
			if v.Suggestion != "pass" {
				return g.Cfg("config").GetString("Tips.ImgSugTip")
			}
		}
	}
	return ""
}

func ImgIsHead(imgUrl string) string {
	res, err := http.Get(imgUrl)
	if err != nil {
		glog.Line().Debug("A error occurred!")
	}
	defer res.Body.Close()
	img1,gs, err := image.Decode(res.Body)
	fmt.Println(gs)
	if err != nil {
		glog.Line().Debug(err)
	}
	if varr := pigoModel.DetectFace(img1, pigoModel.GetArr); len(varr) > 0 {
		if _, err := gohash.AverageHash(varr[0]); err == nil {
			glog.Line().Debug("是人头像")
			return ""
		}
	} else {
		glog.Line().Debug("未从图像检测到人脸,请重新上传")
		return g.Cfg("config").GetString("Tips.IsHeadTip")
	}
	return ""
}