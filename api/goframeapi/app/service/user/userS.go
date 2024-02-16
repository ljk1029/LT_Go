package user

import (
	"encoding/json"
	"fmt"
	"gf-app/app/model/user"
	"gf-app/app/service/goRedis"
	"gf-app/app/service/oss"
	"github.com/go-redis/redis"
	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/issue9/unique"
	"strconv"
	"time"
)

//用户制作完成，公开的视频
const USER_VIDEO = "cygo.api.user.video.%s"
const USER_MYSELF_VIDEO_LIST = "cygo.api.user.myself.video.%s"
const USER_MYSELF_VIDEO_LIST_COUT = "cygo.api.user.myself.count"
const USER_OPENID = "cygo.api.user.openid"
type Entity struct {
	Name      string `orm:"name"   json:"name"`             // 用户名称
	Openid    string `orm:"openid"   json:"openid"`         // openid
	Avatar    string `orm:"avatar"   json:"avatar"`         // avatar
	Code    string `orm:"code"   json:"code"`         // code
	CreatedAt time.Time `orm:"created_at"   json:"created_at"` // 添加时间
	UpdatedAt time.Time `orm:"updated_at"   json:"updated_at"` // 修改时间
}
func GetUserInfo(openid string) interface{} {
	var tem interface{}
	tem = make([]interface{}, 0)
	userInfo,err := g.Redis().DoVar("HGET",USER_OPENID,openid)
	if err != nil {
		fmt.Println(err)
	}
	if len(userInfo.Map()) > 0 {
		var arr map[string]interface{}
		err = json.Unmarshal(gconv.Bytes(userInfo),&arr)
		return arr
	}
	info, _ := g.DB("").Table("merge_video_user").Where("openid=?",openid).One()
	if info == nil{
		uniqueId := unique.Number().String()
		var userdata = &Entity{
			Name:      "john",
			Openid:    openid,
			Avatar:    "111",
			Code:    uniqueId,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := user.Model.Data(userdata).Insert()
		if err != nil {
			return tem
		}
		info, _ = g.DB("").Table("merge_video_user").Where("openid=?",openid).One()
	}
	jsonInfo,_ := json.Marshal(info)
	g.Redis().Do("HSET",USER_OPENID,openid,jsonInfo)
	g.Redis().Do("expire",USER_OPENID,7200)
	return info
}
func GetUserVideo(openid string, limit int, page int) interface{} {
	redisKeyMyself := fmt.Sprintf(USER_MYSELF_VIDEO_LIST,openid)
	exists, err := g.Redis().DoVar("EXISTS", redisKeyMyself)
	if err != nil {
		return err
	}
	if exists.Bool() == false {
		userId,_ := g.DB().Table("merge_video_user").Where("openid=?",openid).Fields("id").Value()
		All, _ := g.DB().Table("merge_video").Fields("code as id,title,created_at as audittime,img_url as img,user_state,state").Where("user_id=?",userId).Order("id desc").Page(page, limit).All()
		for _,v := range All{
			var newV map[string]interface{}
			newV = gconv.Map(v)
			if newV != nil {
				res,_:= gtime.StrToTime(gconv.String(newV["audittime"]))
				//放到有序集合start
				zjson,err := json.Marshal(newV)
				if err != nil{
					fmt.Println(err)
				}
				g.Redis().Do("ZADD",redisKeyMyself,res.Unix(),zjson)
				g.Redis().Do("expire",redisKeyMyself,7200)
			}
		}
	}
	errRedisConn := goRedis.InitClient()
	if errRedisConn != nil {
		//redis连接错误
		fmt.Println(errRedisConn)
	}
	//第一版，没有使用
	//正在redis制作的加入有序集合
	//makeRedisKey := fmt.Sprintf(redisKey.REDIS_API_OPENID_RND_TYPE_MAKE, openid)
	//data, err := goRedis.RedisDb.HGetAll(makeRedisKey).Result()
	//if err != nil {
	//	panic(err)
	//}
	//if len(data) >0 {
	//	for key,value := range data{
	//		imgKey := fmt.Sprintf(redisKey.REDIS_API_OPENID_RND_IMG_MAKE, openid,key)
	//		imgValue, err := goRedis.RedisDb.ZRange(imgKey, 0, 0).Result()
	//		if err != nil {
	//			fmt.Println(err)
	//		}
	//		var arr map[string]interface{}
	//		json.Unmarshal(gconv.Bytes(value),&arr)
	//		if len(arr) >0 {
	//			res,_:= gtime.StrToTime(gconv.String(arr["created_at"]))
	//			var zarr = g.Map{
	//				"audittime" :arr["created_at"],
	//				"id" :arr["code"],
	//				"img" :imgValue[0],
	//				"state" :arr["state"],
	//				"title" :arr["title"],
	//				"user_state" :arr["user_state"],
	//			}
	//			zjson,err := json.Marshal(zarr)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//			g.Redis().Do("ZADD",redisKeyMyself,res.Unix(),zjson)
	//			g.Redis().Do("expire",redisKeyMyself,7200)
	//		}
	//
	//	}
	//
	//}
	offset := (page - 1) * limit
	op := redis.ZRangeBy{
		Max:    "+inf", // 最大分数
		Min:    "-inf", // 最小分数
		Offset: gconv.Int64(offset), // 类似sql的limit, 表示开始偏移量
		Count:  gconv.Int64(limit),  // 一次返回多少数据
	}
	list, err := goRedis.RedisDb.ZRevRangeByScore(redisKeyMyself, op).Result()
	if err != nil {
		panic(err)
	}
	time := gtime.Now().Add(+time.Hour*2).Timestamp()
	l := glist.New()
	redisKey := fmt.Sprintf(USER_VIDEO,gtime.Now().Format("Y-m-d"))
	for _,v := range list{
		var newV map[string]interface{}
		newV = gconv.Map(v)
		if newV != nil {
			res,_:= gtime.StrToTime(gconv.String(newV["audittime"]))
			newV["audittime"] = res.Format("Y.m.d")
			newV["type"] = 1
			if gconv.Int(newV["state"]) == 1 {
				newV["makimg_time"] = res.Add(+gtime.M*5).Format("Y年m月d日 H:i")
			}
			if gconv.Int(newV["user_state"]) ==1 {
				json,err := json.Marshal(newV)
				if err != nil{
					fmt.Println(err)
				}else {
					g.Redis().Do("SADD",redisKey,json)
					g.Redis().Do("expire",redisKey,7200)
				}
			}
			if gconv.String(newV["img"]) == "" {
				fmt.Sprintf("缩略图为空")
			}else{
				newV["img"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(newV["img"]))
			}
			l.PushBack(newV)
		}
	}
	return l
}
func GetUserVideoCount(openid string) int {
	key := USER_MYSELF_VIDEO_LIST_COUT
	count, err := g.Redis().DoVar("HGET", key,openid)
	if err != nil {
		fmt.Println(err)
	}
	if count.Int() > 0 {
		return count.Int()
	}
	userId,_ := g.DB().Table("merge_video_user").Where("openid=?",openid).Fields("id").Value()
	r, err := g.DB().Table("merge_video").Where("user_id=?",userId).Count()
	if err != nil {
		return 0
	}
	g.Redis().DoVar("HSET", key,openid,r)
	return r
}
func GetUserValue(openid string,filed string) interface{} {
	value, _ := user.Model.Where("openid=?",openid).Fields(filed).Value()
	return value
}
