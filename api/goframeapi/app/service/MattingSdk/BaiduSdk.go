package MattingSdk

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"gf-app/app/service/oss"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"io/ioutil"
	"net/http"
)

func GetBaiduMattingImg(imgUrl string,name string) string {
	CLIENTID := g.Cfg("config").GetString("Baidu.BaiduKey")
	CLIENTSECTET := g.Cfg("config").GetString("Baidu.BaiduSectet")
	//imgUrl := "https://img.ipcfun.com/uploads/ishoulu/pic/2013/05/9215193abd60b5ff099795216.jpg"
	//获取远端图片
	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred!")
	}
	defer res.Body.Close()

	// 读取获取的[]byte数据
	data, _ := ioutil.ReadAll(res.Body)

	imageBase64 := base64.StdEncoding.EncodeToString(data)

	request_url := "https://aip.baidubce.com/rest/2.0/image-classify/v1/body_seg?access_token=%s"
	host := "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=%s&client_secret=%s"
	getTokenUrl := fmt.Sprintf(host,CLIENTID,CLIENTSECTET)

	r1, err := g.Client().Post(getTokenUrl)
	if err != nil {
		fmt.Println(err)
	}
	var accessTokenArr map[string]interface{}
	json.Unmarshal([]byte(r1.ReadAllString()),&accessTokenArr)
	fmt.Println(accessTokenArr["access_token"])
	getImgUrl := fmt.Sprintf(request_url,accessTokenArr["access_token"])

	r2, err := g.Client().Post(getImgUrl,g.Map{"image":imageBase64,"type":"foreground"})
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(r2.ReadAllString())
	var R2Arr map[string]interface{}
	json.Unmarshal([]byte(r2.ReadAllString()),&R2Arr)
	//fmt.Println(gconv.String(R2Arr["foreground"]))
	ddd, _ := base64.StdEncoding.DecodeString(gconv.String(R2Arr["foreground"])) //成图片文件并把文件写入到buffer
	//err2 := ioutil.WriteFile("./output.jpg", ddd, 0666)   //buffer输出到jpg文件中（不做处理，直接写到文件）
	//if err2 != nil {
	//	fmt.Println(err2)
	//}
	baiduImgName := "baidu"+name
	resultfile,err := oss.UploadFileByFileStream(gconv.Bytes(ddd), baiduImgName)

	glog.Line().Debug(resultfile)
	glog.Line().Debug(err)
	return ""
}