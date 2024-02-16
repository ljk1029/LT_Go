package router

import (
	"gf-app/app/api/hello"
	"gf-app/app/api/statistics"
	"gf-app/app/api/user"
	"gf-app/app/api/video"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/", hello.Hello)
		group.Group("/video", func(group *ghttp.RouterGroup) {
			group.POST("/templateList", video.Template)
			group.POST("/musicList", video.Music)
			group.POST("/makeVideo", video.MakeVideoByImg)
			group.POST("/uploadImg", video.Uploadimg)
			group.POST("/getFlowList", video.GetFlowList)
			group.POST("/getVideoDetail", video.GetVideoDetail)
			group.GET("/makeComplete", video.MakeComplete)
			group.POST("/queryVideoState", video.QueryVideoState)
			group.POST("/code", video.GetOpenidByCode)
		})
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.POST("/getUserInfo", user.GetInfo)
			group.POST("/getUserVideo", user.GetUserVideo)
			group.POST("/ceshi", user.Ceshi)
		})
		group.Group("/statistics", func(group *ghttp.RouterGroup) {
			group.POST("/videoclick", statistics.Videoclick)
		})
	})
}
