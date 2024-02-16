package user

import (
	"fmt"
	"gf-app/app/service/user"
	"gf-app/library/response"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

type userReq struct {
	Openid string `p:"openid" v:"required#用户openid不能为空"`
}
type userVideoRequest struct {
	Openid string `p:"openid" v:"required#用户openid不能为空"`
	Limit string `p:"limit" v:"required#条数必填"`
	Page  string `p:"page" v:"required#页数不能为空"`
}
//获取当前用户
func GetInfo(r *ghttp.Request) {
	var data *userReq
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}

	info := user.GetUserInfo(data.Openid)

	response.JsonExit(r, 200, "获取成功", info)
}
func GetUserVideo(r *ghttp.Request){
	var data *userVideoRequest
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	All := user.GetUserVideo(data.Openid,gconv.Int(data.Limit), gconv.Int(data.Page))
	DataCount := user.GetUserVideoCount(data.Openid)
	if All == nil {
		response.JsonExit(r, 200, "暂无数据", g.Map{})
	}
	response.JsonExit(r, 200, "获取成功", g.Map{"Count": DataCount, "Data": All})
}
func Ceshi(r *ghttp.Request) {
	id:= "22222"
	fmt.Println(len(id)>2)
	//conn := g.Redis("cache").Conn()
	//fmt.Println(conn)
	//defer conn.Close()
	//bloom := bloom.NewBloom(conn,"Bloom") //创建过滤器
	//b := bloom.Exist("newClien") //判断是否存在这个值
	//fmt.Println(b)

}
func getName(params ...interface{}) string {
	var stringSlice []string
	for _, param := range params {
		stringSlice = append(stringSlice, param.(string))
	}
	return strings.Join(stringSlice, "_")
}
