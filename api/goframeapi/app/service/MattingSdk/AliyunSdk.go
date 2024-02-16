package MattingSdk

import (
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	viapiutil "github.com/alibabacloud-go/viapi-utils/client"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/gogf/gf/frame/g"
)

func GetMattingImg(imgUrl string) string {
	accessKeyId := g.Cfg("config").GetString("aliyunoss.AccessKeyId")
	accessKeySecret := g.Cfg("config").GetString("aliyunoss.AccessKeySecret")
	client, err := sdk.NewClientWithAccessKey("cn-qingdao", accessKeyId, accessKeySecret)
	if err != nil {
		panic(err)
	}
	accessKeyId1 := tea.String(accessKeyId)
	// 你的 access_key_secret
	accessKeySecret1 := tea.String(accessKeySecret)
	// 要上传的文件路径，url 或 filePath
	fileUrl := tea.String(imgUrl)
	// 上传成功后，返回上传后的文件地址
	fileLoadAddress, _err := viapiutil.Upload(accessKeyId1, accessKeySecret1, fileUrl)
	if _err != nil {
		fmt.Println(_err)
	}
	fmt.Println(*fileLoadAddress)
	fmt.Println(22)
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "imageseg.cn-shanghai.aliyuncs.com"
	request.Version = "2019-12-30"
	request.ApiName = "SegmentCommonImage"
	request.QueryParams["RegionId"] = "cn-shanghai"
	request.QueryParams["ImageURL"] = *fileLoadAddress
	//request.QueryParams["ImageURL"] = "https://cykj-composite-video.oss-cn-qingdao.aliyuncs.com/5dbbf19a86a2d_162.jpg?Expires=1607414893&OSSAccessKeyId=TMP.3Kg3rQujwDjnmjNDwLm5cwtNSDgvKtze5WK8o11WQpNHmszT1iQ3VnYMwhhuktha9hDAh2uWPzBrr3FxryAdZndiYTc92F&Signature=ZQqHGPXtHMwf5GBx7%2BeHV3uNkXM%3D"

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
	return ""
}