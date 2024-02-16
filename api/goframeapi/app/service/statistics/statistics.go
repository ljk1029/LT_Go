package statistics

import (
	"fmt"
	"gf-app/app/service/bloom"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)
//编辑的
//每分钟pv
const VIDEO_HOUR_PV = "cygo_api.hour.video.pv.%s.%s"
//每分钟uv
const REDIS_HOUR_UV = "cygo_api.hour.video.uv.%s.%s"
//每天Cache
const REDIS_HOUR_BlOOM = "cygo_api.hour.video.bloom.%s"

//用户的
//每分钟pv
const VIDEO_HOUR_PV_TYPE = "cygo_api.user.hour.video.pv.%s.%s"
//每分钟uv
const REDIS_HOUR_UV_TYPE = "cygo_api.user.hour.video.uv.%s.%s"
//每天Cache
const REDIS_HOUR_BlOOM_TYPE = "cygo_api.user.hour.video.bloom.%s"

func VideoClickPvUv(id string,openId string) error {
	date := gtime.Now().Format("Ymd")
	min := gtime.Now().Format("H:i")
	if len(openId) > 0 {
		conn := g.Redis("cache").Conn()
		fmt.Println(conn)
		defer conn.Close()
		hourPvKey := GetRedisKeyVideoHourPv(date, id)
		hourUvKey := GetRedisKeyVideoHourUv(date, id)
		hourBloomKey := GetRedisKeyVideoHourBloom(date,id)
		bloom := bloom.NewBloom(conn,hourBloomKey) //创建过滤器
		if bloom.Exist(id+openId) == false {
			bloom.Add(id+openId)
			g.Redis("cache").Do("hincrby",hourUvKey,min,1)
		}
		g.Redis("cache").Do("hincrby",hourPvKey,min,1)
	}
	return nil
}
func GetRedisKeyVideoHourPv(date string, code string) string {
	if len(code) > 20 {
		return fmt.Sprintf(VIDEO_HOUR_PV_TYPE, date, code)
	}
	return fmt.Sprintf(VIDEO_HOUR_PV, date, code)
}
func GetRedisKeyVideoHourUv(date string, code string) string {
	if len(code) > 20 {
		return fmt.Sprintf(REDIS_HOUR_UV_TYPE, date, code)
	}
	return fmt.Sprintf(REDIS_HOUR_UV, date, code)
}
func GetRedisKeyVideoHourBloom(date string, code string) string {
	if len(code) > 20 {
		return fmt.Sprintf(REDIS_HOUR_BlOOM_TYPE, date)
	}
	return fmt.Sprintf(REDIS_HOUR_BlOOM, date)
}
