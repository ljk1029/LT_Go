package statistics

import (
	"gf-app/app/service/statistics"
	"gf-app/library/response"
	"github.com/gogf/gf/net/ghttp"
)
type VideoClickReq struct {
	Id string `p:"id" v:"required#id不能为空"`
	OpenId string `p:"openID" v:"required#openID不能为空"`
}
func Videoclick(r *ghttp.Request)  {
	var data *VideoClickReq
	if err := r.Parse(&data); err != nil {
		response.JsonExit(r, 1, err.Error())
	}
	err := statistics.VideoClickPvUv(data.Id,data.OpenId)
	if err != nil {
		response.JsonExit(r, 1, "失败")
	}
	response.JsonExit(r, 200, "成功")
}