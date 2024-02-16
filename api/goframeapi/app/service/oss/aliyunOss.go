package oss

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"path"
	"strings"
)
//OSS缓存
const REDIS_API_VIDEO_LIST_OSSCONFIG = "cykj_api.list.oss_config1"

//初始化oss服务
func Initserver()(client *oss.Client,err error){
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := g.Cfg("config").GetString("aliyunoss.Endpoint")
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := g.Cfg("config").GetString("aliyunoss.AccessKeyId")
	accessKeySecret := g.Cfg("config").GetString("aliyunoss.AccessKeySecret")
	// 创建OSSClient实例。
	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil{
		return
	}
	return
}

//上传文件
func UploadFile(localfile string,uploadfile string,appid string)(resultfile string,err error) {
	resultfile = ""
	// 创建OSSClient实例。
	client, err := Initserver()

	bucketName := g.Cfg("config").GetString("aliyunoss.Bucket")
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	var uploaddir string
	if len(appid)>0 {
		uploaddir = fmt.Sprintf("%s/imgs/%s/%s",appid,gtime.Now().Format("Ym"),gtime.Now().Format("Ymd"))
	}else{
		uploaddir = fmt.Sprintf("video/imgs/%s",gtime.Date())
	}
	uploadfile = strings.Trim(uploadfile, "/")
	objectName := fmt.Sprintf("%s/%s", uploaddir, uploadfile) //完整的oss路径
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := localfile
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return
	}
	resultfile = objectName
	return
}
//上传文件通过文件流
func UploadFileByFileStream(localfile []byte,uploadfile string)(resultfile string,err error) {
	resultfile = ""
	// 创建OSSClient实例。
	client, err := Initserver()

	bucketName := g.Cfg("config").GetString("aliyunoss.Bucket")
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	uploaddir := fmt.Sprintf("video/imgs/%s",gtime.Date())
	uploadfile = strings.Trim(uploadfile, "/")
	objectName := fmt.Sprintf("%s/%s", uploaddir, uploadfile) //完整的oss路径
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := localfile
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return
	}
	// 上传文件流。
	err = bucket.PutObject(objectName, bytes.NewReader([]byte(localFileName)))
	if err != nil {
		return
	}
	resultfile = objectName
	return
}
func GetPrivateUrl(time string,filename string) string {
	cdnurl := g.Cfg("config").GetString("aliyunoss.OssCdn")
	key := g.Cfg("config").GetString("aliyunoss.Key")
	sstr := fmt.Sprintf("%s-%s-0-0-%s",filename,time,key)
	has,_ := MD5(sstr)
	authKey := fmt.Sprintf("auth_key=%s-0-0-%s",time, has)
	url := fmt.Sprintf("%s%s?%s",cdnurl, filename,authKey)
	return url
}
func TencentAmethod(time string,filename string) string {
	cdnurl := g.Cfg("config").GetString("TencentKey.XinyuanCdn")
	key := g.Cfg("config").GetString("TencentKey.TencentKey")
	fileArdd,_ := path.Split(filename)
	sstr := fmt.Sprintf("%s-%s-0-0-%s",fileArdd,time,key)
	has,_ := MD5(sstr)
	authKey := fmt.Sprintf("auth_key=%s-0-0-%s",time, has)
	url := fmt.Sprintf("%s%s?%s",cdnurl, filename,authKey)
	return url
}
func MD5(s string) (md5String string, err error) {
	hasher := md5.New() // #nosec
	_, err = hasher.Write([]byte(s))
	if err != nil {
		return
	}
	cipherStr := hasher.Sum(nil)
	md5String = hex.EncodeToString(cipherStr)
	return
}
func GetOss(ossType string) interface{} {
	Redis := REDIS_API_VIDEO_LIST_OSSCONFIG
	key := fmt.Sprintf("%s_cdn",ossType)
	value,_ := g.Redis().Do("HGET",Redis,key)
	if value != nil {
		return value
	}
	data,_ := g.DB().Table("oss_config").Fields("endpoint,video_bucket,head_bucket,video_cdn,head_cdn,content_bucket,content_cdn,crop_cdn").Where("is_direct=?",1).One()
	var cropCdn string
	for k,v := range data {
		g.Redis().Do("HSET",Redis,k,gconv.String(v))
		g.Redis().Do("expire",Redis,1800)
		if k == key {
			cropCdn = gconv.String(v)
		}
	}
	return cropCdn
}