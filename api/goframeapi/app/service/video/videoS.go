package video

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"gf-app/app/model/makeVideo"
	"gf-app/app/model/music"
	"gf-app/app/model/video"
	"gf-app/app/service/goRedis"
	"gf-app/app/service/oss"
	"gf-app/app/service/redisKey"
	"gf-app/app/service/until"
	"gf-app/app/service/user"
	"gf-app/app/service/wechat"
	"github.com/gogf/gf/container/glist"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
	"github.com/gomodule/redigo/redis"
	"github.com/nahid/gohttp"
	"github.com/noaway/dateparse"
	redisType "github.com/go-redis/redis"
	"math"
	"strconv"
	"strings"
	"time"
)

const OPENID_TEMPLATEID_VIDEOID = "cyapi.%s.%v.%v"
const MERGE_VIDEO_QUEUE = "video.queue.prepare"

//新用户刷新次数
const REDIS_API_KEY_NEW_FLOW_OPEN_NUM = "%s.api.list.%s.flowOpenNum.%s"
const REDIS_API_RANDOM2_NEAR_NUMBER = "%s.api.random.near.number.%s"

//第一段数据(无序)
const REDIS_API_RANDOM2_SHOW_LIST = "%s.apigo.random.show.list.%s"

//原创是空的
const REDIS_API_KEY_ORIGINAL_DATA_ISNULL = "%s.api.list.original.isNull"

//原创列表无序
const REDIS_API_RANDOM2_NEW_ORIGIN_LIST = "%s.api.random.neworigin.list.%s"

//新文章列表KEy(有序)
const REDIS_API_RANDOM2_NEW_LIST = "%s.api.random.new.list.%s"

//新文章读取数量
const REDIS_API_RANDOM2_NEW_NUMBER = "%s.api.random.new.number.%s"

//随机列表2-邻核心数据-列表的RedisKey
const REDIS_API_RANDOM2_NEAR_LIST = "%s.api.random.near.list.%s"

//详情key
const REDIS_API_DETAILS = "cygo_api.details.id.%s"
const REDIS_API_DETAILS_USER = "cygo_api.user.details.id.%s"
const APP_ID = "wxca9f101ee94bd958"

//详情中间key
const REDIS_API_DETAILS_OAUTHORITY = "cygo_api.details.oauthority.id.%s"
const REDIS_API_DETAILS_OAUTHORITY_USER = "cygo_api.user.details.oauthority.id.%s"

//头像key
const REDIS_API_AVATAR_USERNAME = "cygo_api.avatar.username.administrators.id"

//热门,7天内
const CACHE_HOST = 1

//邻近,7-30天
const CACHE_NEAR = 2
const AVATARIMG = "http://zcom.xazzp.com/default.jpg"
const REDIS_API_USER_AVATARIMG = "cygo_api.user.info"
//模板key
const UREDIS_API_TEPLATELIST = "cygo_api.template.list"
//作者的条数
const REDIS_API_VIDEO_VIDEOSCOUNT_WHITE = "cykj_api.list.videoscount.%s"

//制作视频redis
//上传图片未制作
const REDIS_API_OPENID_RND_TYPE = "cygo_api.%s.0"
const REDIS_API_OPENID_RND_IMG = "cygo_api.%s.img"

const REDIS_API_CYGO_API_OPENID_MAKED_VIDEO = "cygo_api.openid.maked.video"
//队列存放到表
const REDIS_API_CYGO_API_QUEUE_TABLE = "cygo_api.queue.table"
//百分比
const REDIS_API_CODE_VIDEO_PLAN = "cygo_api.%s.video.plan"
var db *sql.DB

func GetTemplate(page int, limit int) interface{} {
	redisKey := UREDIS_API_TEPLATELIST
	exists, err := g.Redis().DoVar("EXISTS", redisKey)
	if err != nil {
		return err
	}
	if exists.Bool() == false {
		Allpage, _ := video.Model.Order("id desc").Where("is_show",1).Page(page, limit).All()
		for _,value := range Allpage{
			json,err := json.Marshal(value)
			if err != nil {
				fmt.Println(err)
			}
			g.Redis().Do("ZADD",redisKey,value.Id,json)
			g.Redis().Do("expire",redisKey,7200)
		}
	}
	errRedisConn := goRedis.InitClient()
	if errRedisConn != nil {
		//redis连接错误
		fmt.Println(errRedisConn)
	}
	offset := (page - 1) * limit
	op := redisType.ZRangeBy{
		Max:    "+inf", // 最大分数
		Min:    "-inf", // 最小分数
		Offset: gconv.Int64(offset), // 类似sql的limit, 表示开始偏移量
		Count:  gconv.Int64(limit),  // 一次返回多少数据
	}
	list, err := goRedis.RedisDb.ZRevRangeByScore(redisKey, op).Result()
	if err != nil {
		panic(err)
	}
	l := glist.New()
	for _, v := range list {
		time := gtime.Now().Add(+time.Hour * 2).Timestamp()
		var newV map[string]interface{}
		newV = gconv.Map(v)
		if newV != nil {
			if strings.Index(gconv.String(newV["thumbnail"]), "http") < 0 {
				fmt.Println(gconv.String(newV["thumbnail"]))
				newV["thumbnail"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(newV["thumbnail"]))
			}
			if strings.Index(gconv.String(newV["video_template_url"]), "http") < 0 {
				newV["video_template_url"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(newV["video_template_url"]))
			}
			newV["fixed_pic_num"] = gconv.Int(newV["fixed_pic_num"])
			newV["is_head"] = gconv.Int(newV["is_head"])
		}
		l.PushBack(newV)
	}
	return l
}
func GetTemCount() int {
	r, err := video.Model.Count()
	if err != nil {
		return 0
	}
	return r
}

func GetMusicList(page int, limit int) []*music.Entity {
	Allpage, _ := music.Model.Order("id desc").Page(page, limit).All()
	return Allpage
}
func GetMusicCount() int {
	r, err := music.Model.Count()
	if err != nil {
		return 0
	}
	return r
}
func InserVideoImgToRedis(openid string, templateId int, url string, sort int) (string, error) {
	userInfo := user.GetUserInfo(openid)
	redisKeyType := GetRedisKeyOpenIdRandType(openid)
	redisKeyImg := GetRedisKeyOpenIdRandImg(openid)
	sort +=1
	conn := g.Redis().Conn()
	defer conn.Close()
	exists, err := conn.DoVar("EXISTS", redisKeyImg)
	existsType, err := conn.DoVar("EXISTS", redisKeyType)
	if err != nil {
		return "", err
	}
	if exists.Bool(){
		if existsType.Bool() {
			var ttl int
			typeTTl,_ := conn.DoVar("TTl",redisKeyType)
			fmt.Println(typeTTl)
			conn.Do("ZREMRANGEBYSCORE",redisKeyImg,sort,sort)
			conn.Do("ZADD", redisKeyImg, sort, url)
			if typeTTl.Int() < 10 {
				ttl = typeTTl.Int()
			}else{
				ttl = typeTTl.Int()-1
			}

			conn.Do("expire", redisKeyImg, ttl)
			code, err := conn.DoVar("HGET", redisKeyType, "code")
			if err != nil {
				return "", err
			}
			return code.String(), nil
		}else{
			conn.Do("DEL", redisKeyImg)
		}
	}
	code := guid.S()
	conn.Send("MULTI")
	conn.Send("HMSET", redis.Args{}.Add(redisKeyType).AddFlat(g.Map{
		"template_id": templateId,
		"code":        code,
		"user_id":     gconv.Map(userInfo)["id"],
		"state":       0,
		"created_at":  gtime.Datetime(),
	})...)
	conn.Send("expire", redisKeyType, 7200)
	conn.Send("ZADD", redisKeyImg, sort, url)
	conn.Send("expire", redisKeyImg, 7200)
	_, err2 := conn.Do("EXEC")
	if err2 != nil {
		return "", err2
	}
	return code, nil
}
func InserVideoImg(openid string, templateId int, url string, sort int) (string, error) {
	uid := user.GetUserValue(openid, "id")
	if tx, err := g.DB().Begin(); err == nil {
		// 方法退出时检验返回值，
		// 如果结果成功则执行tx.Commit()提交,
		// 否则执行tx.Rollback()回滚操作。
		defer func() {
			if err != nil {
				tx.Rollback()
			} else {
				tx.Commit()
			}
		}()
		code := guid.S()
		// 写入合成视频表
		r, err := tx.Table("merge_video").Insert(g.Map{
			"template_id": templateId,
			"code":        code,
			"user_id":     uid,
			"state":       0,
			"created_at":  gtime.Date(),
		})
		if err != nil {
			return "", err
		}
		videoId, err := r.LastInsertId()
		// 写入图片
		r, err = tx.Table("merge_video_imgs").Insert(g.Map{
			"video_id":   videoId,
			"img_url":    url,
			"sort":       sort,
			"created_at": gtime.Date(),
		})
		if err != nil {
			return "", err
		}
		return code, nil
	}
	return "", nil
}
func MakeVideoToAdmin(appid string,Openid string, Title string, Content string, TemplateId int, Code string, UserState int,templateFunc string,author string) error {
	//userInfo := user.GetUserInfo(Openid)
	redisKeyType := GetRedisKeyOpenIdRandType(Openid)
	//redisKeyImg := GetRedisKeyOpenIdRandImg(Openid)
	conn := g.Redis().Conn()
	defer conn.Close()
	exists, err := conn.DoVar("EXISTS", redisKeyType)
	if err != nil {
		return err
	}
	if exists.Bool() {
		req := gohttp.NewRequest()
		ch := make(chan *gohttp.AsyncResponse)
		//makeVideoUrl := "http://47.105.140.133:8199/videoTemplate/index"
		makeVideoUrl := g.Cfg("config").GetString("MakeAdmin.url")
		req.
			FormData(map[string]string{"id":Code,
			"openid":Openid,
			"appid":appid,
			"title":Title,
			"funcName":templateFunc,
			"content":Content,
			"author":author,
			"templateId":gconv.String(TemplateId),
			"userState":gconv.String(UserState),}).
			AsyncPost(makeVideoUrl, ch)
		glog.Line().Info("制作视频中。。。。")
		return nil
	}
	return errors.New("重新上传图片")
}
func MakeVideoToRedis(appid string,Openid string, Title string, Content string, TemplateId int, Code string, UserState int) error {
	//userInfo := user.GetUserInfo(Openid)
	redisKeyType := GetRedisKeyOpenIdRandType(Openid)
	redisKeyImg := GetRedisKeyOpenIdRandImg(Openid)
	redisKeyTypeMake := redisKey.GetRedisKeyOpenIdRandTypeMake(Openid)
	redisKeyImgMake := redisKey.GetRedisKeyOpenIdRandImgMake(Openid,Code)
	conn := g.Redis().Conn()
	defer conn.Close()
	exists, err := conn.DoVar("EXISTS", redisKeyType)
	if err != nil {
		return err
	}
	if exists.Bool() {
		redisKeyTypeContent, err := conn.DoVar("HGetAll", redisKeyType)
		if err != nil {
			return err
		}
		redisKeyContent := redisKeyTypeContent.Map()
		redisKeyContent["title"] = Title
		redisKeyContent["content"] = Content
		redisKeyContent["user_state"] = UserState
		redisKeyContent["openid"] = Openid
		redisKeyContent["updated_at"] = gtime.Datetime()

		redisKeyImgContent, err := conn.DoVar("ZRANGE", redisKeyImg, 0, -1, "WITHSCORES")
		redisKeyImgContent1, err := conn.DoVar("ZRangeByScore", redisKeyImg, 1, 10)
		if err != nil {
			return err
		}
		videofileName := guid.S()
		var Data = g.Map{
			"jobId":      Code,
			"imgList":    redisKeyImgContent1.Array(),
			"templateId": TemplateId,
			"openid":     Openid,
			"fileName":   fmt.Sprintf("%s.mp4",videofileName),
			"videoPath":  fmt.Sprintf("video/%s/%s/%s.mp4", gtime.Now().Format("Ym"), gtime.Now().Format("Ymd"), videofileName),
		}
		data,err := json.Marshal(Data)
		redisKeyContentJson,err := json.Marshal(redisKeyContent)
		conn.Send("MULTI")
		conn.Send("HSET", redisKeyTypeMake,redisKeyContent["code"],redisKeyContentJson)
		for i, v := range redisKeyImgContent.Map() {
			conn.Send("ZADD", redisKeyImgMake, v, i)
		}
		conn.Send("RPUSH", MERGE_VIDEO_QUEUE, data)
		conn.Send("DEL", redisKeyType)
		conn.Send("DEL", redisKeyImg)
		_, err1 := conn.Do("EXEC")
		if err1 != nil {
			return err1
		}
		return nil
	}
	return errors.New("重新上传图片")
}
func MakeVideo(Openid string, Title string, Content string, TemplateId int, Code string, UserState int) error {
	uid := user.GetUserValue(Openid, "id")
	VideoId, err := makeVideo.Model.Where("user_id=?", uid).Where("template_id=?", TemplateId).Where("code=?", Code).Fields("id").Value()
	if err != nil {
		return err
	}
	_, err = makeVideo.Model.Where("id=?", VideoId).Save(g.Map{
		"id":         VideoId,
		"title":      Title,
		"content":    Content,
		"user_state": UserState,
		//"music_id":MusicId,
		"updated_at": gtime.Date(),
	})
	if err != nil {
		return err
	}
	imgList, err := g.DB().Table("merge_video_imgs").Order("sort asc").Fields("img_url").Where("video_id=?", VideoId).Array()
	if err != nil {
		return err
	}
	videofileName := guid.S()
	var Data = g.Map{
		"jobId":      Code,
		"imgList":    imgList,
		"templateId": TemplateId,
		"openid":     Openid,
		"fileName":   fmt.Sprintf("%s.mp4", videofileName),
		"videoPath":  fmt.Sprintf("video/%s/%s/%s.mp4", gtime.Now().Format("Ym"), gtime.Now().Format("Ymd"), videofileName),
	}
	_, err = g.Redis().Do("RPUSH", MERGE_VIDEO_QUEUE, Data)
	if err != nil {
		return err
	}
	return nil
}

func RedisKeyOpenidTemplateId(openid string, TemplateId int, VideoId int) string {
	return fmt.Sprintf(OPENID_TEMPLATEID_VIDEOID, openid, TemplateId, VideoId)
}

func GetOpenIdNumByOpenId(openid string) interface{} {
	uid := openid[len(openid)-1:]
	var newNum = 0
	var rnum = 0
	openNumKey := GetRedisKeyFlowOpenNum(gtime.Now().Format("Ymd"), uid, APP_ID)
	num, _ := g.Redis().Do("HGET", openNumKey, openid)
	if num == nil {
		rnum = 0
	} else {
		rnum = gconv.Int(gconv.String(num.([]byte)))
	}

	if rnum > 0 {
		newNum = rnum
	} else {
		newNum = 1
	}
	g.Redis().Do("HSET", openNumKey, openid, newNum+1)
	return newNum
}
func GetTwoNum() string {
	RedisKey := GetRedisKeyrandom2NearNumber()
	return RedisKey
}
func AllListGetDataByReloadNum(reloadNum int, twoNum int) interface{} {
	var newNum = 0
	num := twoNum / 10
	if num >= 3 {
		num = 3
	}
	switch {
	case num == 1:
		if reloadNum%4 == 0 {
			newNum = 2
		} else {
			newNum = 1
		}
	case num == 2:
		if math.Abs(gconv.Float64(reloadNum%5-2)) == 2 {
			newNum = 2
		} else {
			newNum = 1
		}
	case num == 3:
		if gconv.Int(math.Ceil(gconv.Float64(reloadNum/3)))%2 == 0 {
			newNum = 2
		} else {
			newNum = 1
		}
	default:
		newNum = 1
	}
	if newNum == 1 {
		return GetFirstFlowList(reloadNum)
	} else {
		return GetTwoFlowList()
	}
	return nil
}

//获取第一段数据
func GetFirstFlowList(reloadNum int) interface{} {
	var todayNum int
	firstFlowListRedisKey := GetRedisKeyrandom2ShowList()

	data, _ := g.Redis().DoVar("srandmember", firstFlowListRedisKey, 10)
	if len(data.Array()) > 0 {
		test1 := glist.New()
		for _, value := range data.Array() {
			var res map[string]interface{}
			err := json.Unmarshal(gconv.Bytes(value), &res)

			if err != nil {
				fmt.Println(err)
			}
			test1.PushBack(res)
		}
		return test1
	}
	//if len(gconv.Map(data)) > 0 {
	//	return GetAllListOriginal(data,reloadNum)
	//}
	redisKey := GetRedisKeyrandom2NewList()
	todayNum = GetTodayNumFromRedis()
	if todayNum == 0 {
		todayNum = 100
	}
	newData, err := g.Redis().DoVar("ZREVRANGEBYSCORE", redisKey, gtime.Now().Unix(), "-inf")
	if err != nil {
		return err
	}
	//return gconv.String(newData)
	//if newData == nil {
	//	inserData, _ := g.DB().Table("video").Fields("code as id', 'title', 'img0 as img', 'audittime").Where(g.Map{
	//		"is_del =" : 0,
	//		"status" : 1,
	//		"is_look" : 1,
	//		"blacklist" : 0,
	//		"audittime <" : gtime.Now().Unix(),
	//	}).Limit(100).All()
	//	time := gtime.Now().Add(+time.Hour*24).Unix()
	//	for _,v:=range inserData {
	//		d,_ := json.Marshal(v)
	//		g.Redis().Do("ZADD",redisKey, d)
	//		g.Redis().Do("expireat",redisKey, time)
	//	}
	//	dataRedis,_ := g.Redis().DoVar("zRange",redisKey,0,99)
	//	if dataRedis != nil {
	//		return dataRedis
	//	}
	//}
	newData1 := DbProcessData(newData)
	//获取新核心数据
	length := newData1.Len()
	if length > 0 {
		for i, e := 0, newData1.Front(); i < length; i, e = i+1, e.Next() {
			//return e.Value
			g.Redis().Do("SADD", firstFlowListRedisKey, e.Value)
		}
	}
	//用户视频
	redisKeyUserVideo := fmt.Sprintf(user.USER_VIDEO, gtime.Now().Format("Y-m-d"))
	userVideo, userVideoErr := g.Redis().DoVar("SRANDMEMBER", redisKeyUserVideo, 40)
	if userVideoErr != nil {
		fmt.Println("用户视频错误")
	} else {
		fmt.Println("1111")
		if len(userVideo.Array()) > 0 {
			for _, v := range userVideo.Array() {
				var res map[string]interface{}
				err := json.Unmarshal(gconv.Bytes(v), &res)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("222")
				if res != nil {
					fmt.Println(res)
					time := gtime.Now().Add(+time.Hour * 2).Timestamp()
					resTime, _ := gtime.StrToTime(gconv.String(res["audittime"]))
					res["audittime"] = resTime.Unix()
					res["type"] = 1
					res["img"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(res["img"]))
					json, err := json.Marshal(res)
					if err != nil {
						fmt.Println(err)
					} else {
						g.Redis().Do("SADD", firstFlowListRedisKey, json)
					}

				}
			}
		}
	}
	g.Redis().Do("expire", firstFlowListRedisKey, 3600)
	newData2, _ := g.Redis().DoVar("srandmember", firstFlowListRedisKey, 10)
	if len(newData2.Array()) > 0 {
		test2 := glist.New()
		for _, value := range newData2.Array() {
			var res2 map[string]interface{}
			err := json.Unmarshal(gconv.Bytes(value), &res2)

			if err != nil {
				fmt.Println(err)
			}
			test2.PushBack(res2)
		}
		return test2
	}
	return newData2
}
func GetTwoFlowList() *glist.List {
	RedisKey := fmt.Sprintf(REDIS_API_RANDOM2_NEAR_LIST, APP_ID, gtime.Now().Format("Y-m-d"))
	data, _ := g.Redis().DoVar("srandmember", RedisKey, 20)
	test1 := glist.New()
	if len(gconv.Map(data)) > 0 {

		for _, value := range gconv.Map(data) {
			var res map[string]interface{}
			err := json.Unmarshal(gconv.Bytes(value), &res)

			if err != nil {
				fmt.Println(err)
			}
			test1.PushBack(res)
		}
		return test1
	}
	return test1
}

func GetVideo(id string, openId string) interface{} {
	RedisKey := GetRedisCyapiDetails(id)
	redisVal, _ := g.Redis().DoVar("HGETALL", RedisKey)
	if len(gconv.Map(redisVal)) > 0 {
		var newRedisval = make(map[string]interface{})
		for i, v := range gconv.Map(redisVal) {
			newRedisval[i] = v
		}
		return newRedisval
	}
	rediskeyoauthority := GetRedisCyapiDetailsOaut(id)
	redisValRawdata, _ := g.Redis().DoVar("HGETALL", rediskeyoauthority)
	var newDetails map[string]interface{}
	fmt.Println(len(gconv.Map(redisValRawdata)))
	if len(gconv.Map(redisValRawdata)) > 0 {
		newDetails = gconv.Map(redisValRawdata)
		vidPathObj := 0
		for i, _ := range gconv.Map(redisValRawdata) {
			if i == "vidPathObj" {
				vidPathObj = 1
			}
		}
		if vidPathObj == 0 {
			return redisValRawdata
		}
		fmt.Println(111)
	} else {
		fmt.Println("查视频详情")
		details := Videodetails(id)
		fmt.Println(details)
		_, err := g.Redis().Do("HMSET", redis.Args{}.Add(rediskeyoauthority).AddFlat(details)...)
		if err != nil {
			fmt.Println(err)
		}
		g.Redis().Do("expire", rediskeyoauthority, 86400*7)
		newDetails = details
	}
	fmt.Println(333)
	var detailsRawdata interface{}
	if len(id) > 20 {
		detailsRawdata = UserInfoOauthority(newDetails, openId)
	} else {
		detailsRawdata = Oauthority(newDetails)
	}
	g.Redis().Do("HMSET", redis.Args{}.Add(RedisKey).AddFlat(detailsRawdata)...)
	g.Redis().Do("expire", RedisKey, 1800)
	return detailsRawdata
}
func UserInfoOauthority(details gdb.Map, openId string) interface{} {
	//userInfo := GetMiniUserInfo(openId)
	//return userInfo
	time := gtime.Now().Add(+time.Hour * 2).Timestamp()
	dataValue := g.Map{}
	count, _ := g.DB().Table("merge_video").Where(g.Map{"code": details["id"]}).Count()
	dataValue["operatorImg"] = AVATARIMG
	dataValue["operatorVideos"] = gconv.Int(count)
	dataValue["operatorId"] = details["user_mini_id"]
	dataValue["operator"] = details["author"]
	dataValue["title"] = details["title"]
	dataValue["introduce"] = dataValue["title"]
	dataValue["water_imgsrc"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(details["img0"]))
	dataValue["imgSrc"] = dataValue["water_imgsrc"]
	dataValue["time"] = gconv.String(details["audittime"])
	dataValue["leftId"] = 0
	dataValue["rightId"] = 0
	dataValue["ON"] = 0
	dataValue["videoProportion"] = ""
	dataValue["operatorClick"] = ""
	dataValue["videoType"] = ""
	dataValue["vidPathObj"] = ""
	dataValue["pripath"] = ""
	dataValue["othervideo_path"] = ""
	dataValue["is_othervideo"] = ""
	dataValue["excitationTime"] = ""
	dataValue["excitationState"] = ""
	dataValue["slideFlag"] = ""
	dataValue["template_id"] = ""
	newPath := fmt.Sprintf("/%s", details["path"])
	url := oss.GetPrivateUrl(strconv.FormatInt(time, 10), newPath)
	dataValue["path"] = url

	return dataValue
}
func Oauthority(details gdb.Map) interface{} {
	cropUrl := oss.GetOss("crop")
	dataValue := g.Map{}
	userInfo := GetusereInfo(gconv.Int(details["user_mini_id"]))
	if userInfo["avatar"] == nil {
		dataValue["operatorImg"] = AVATARIMG
	} else {
		dataValue["operatorImg"] = gconv.String(userInfo["avatar"])
	}
	dataValue["operatorVideos"] = VideosCount(gconv.Int(details["user_mini_id"]))
	dataValue["operatorId"] = details["user_mini_id"]
	dataValue["operator"] = userInfo["login_name"]
	dataValue["title"] = until.DynamicTitle(gconv.String(details["title"]))
	dataValue["introduce"] = dataValue["title"]
	dataValue["water_imgsrc"] = fmt.Sprintf("%s%s", cropUrl, details["img0"])
	dataValue["imgSrc"] = fmt.Sprintf("%s%s%s", cropUrl, details["img0"], "?x-oss-process=image/watermark,image_bG9vazUucG5nP3gtb3NzLXByb2Nlc3M9aW1hZ2UvcmVzaXplLFBfMzA,t_100,g_sw,t_80/watermark,image_Ym9mYW5nLnBuZz94LW9zcy1wcm9jZXNzPWltYWdlL3Jlc2l6ZSxQXzMw,t_90,g_center,t_100")
	t, _ := dateparse.ParseAny(gconv.String(details["audittime"]))
	dataValue["time"] = t.Format("2006-01-02")
	dataValue["leftId"] = 0
	dataValue["rightId"] = 0
	dataValue["ON"] = 0
	dataValue["videoProportion"] = ""
	dataValue["operatorClick"] = ""
	dataValue["videoType"] = ""
	dataValue["vidPathObj"] = ""
	dataValue["pripath"] = ""
	dataValue["othervideo_path"] = ""
	dataValue["is_othervideo"] = ""
	dataValue["excitationTime"] = ""
	dataValue["excitationState"] = ""
	dataValue["slideFlag"] = ""
	dataValue["template_id"] = ""
	time := gtime.Now().Add(+time.Hour * 2).Timestamp()
	newPath := fmt.Sprintf("/%s", details["path"])
	url := oss.TencentAmethod(strconv.FormatInt(time, 10), newPath)
	dataValue["path"] = url

	return dataValue
}

func VideosCount(userMiniId int) int {
	rediskey := GetRedisKeyVideoCount(gconv.String(userMiniId))
	redisval, err := g.Redis().DoVar("GET", rediskey)
	if err != nil {
		fmt.Println(err)
	}
	if gconv.String(redisval) == "" {
		count, err := g.DB().Table("video").Where(g.Map{"user_mini_id": userMiniId, "is_del": 0, "status": 1}).Count()
		if err != nil {
			fmt.Println(err)
		}
		g.Redis().Do("SET", rediskey, count)
		g.Redis().Do("expire", rediskey, 1800)
		return count
	}
	return gconv.Int(redisval)
}
func GetusereInfo(userMiniId int) map[string]interface{} {
	avatarkey := REDIS_API_AVATAR_USERNAME
	userChildInfo, err := g.Redis().DoVar("HGET", avatarkey, userMiniId)
	if err != nil {
		fmt.Println(err)
	}
	if len(gconv.Map(userChildInfo)) > 0 {
		var arr map[string]interface{}
		err = json.Unmarshal(gconv.Bytes(userChildInfo), &arr)
		return arr
	}

	userInfo, err := g.DB("admin").Table("child_user").Fields("id,avatar,login_name").Where("id=?", userMiniId).One()
	if err != nil {
		fmt.Println(err)
	}
	jsonInfo, _ := json.Marshal(userInfo)
	g.Redis().Do("HSET", avatarkey, userMiniId, jsonInfo)
	g.Redis().Do("expire", avatarkey, 7200)
	return gconv.Map(userInfo)
}
func GetMiniUserInfo(openId string) interface{} {
	avatarkey := REDIS_API_USER_AVATARIMG
	userInfo, err := g.Redis().DoVar("HGET", avatarkey, openId)
	if err != nil {
		fmt.Println(err)
	}
	if len(userInfo.Map()) > 0 {
		var arr map[string]interface{}
		err = json.Unmarshal(gconv.Bytes(userInfo), &arr)
		return arr
	}

	userNewInfo, err1 := g.DB().Table("merge_video_user").Fields("id,avatar,name as login_name").Where("openid=?", openId).One()
	if err1 != nil {
		fmt.Println(err1)
	}
	//if gconv.Map(userNewInfo)["avatar"] == "" {
	appid := g.Cfg("config").GetString("AppInfo.AppId")
	appSecret := g.Cfg("config").GetString("AppInfo.AppSecret")
	token := wechat.GetAccessToken(appid, appSecret)
	getTokenUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN", token, openId)
	if r, err := g.Client().Get(getTokenUrl); err != nil {
		fmt.Println(err)
		fmt.Println("获取用户信息")
	} else {
		defer r.Close()
		fmt.Println(r.ReadAllString())
		return r.ReadAllString()
	}
	//}
	return "1111"
	jsonInfo, _ := json.Marshal(userNewInfo)
	g.Redis().Do("HSET", avatarkey, openId, jsonInfo)
	g.Redis().Do("expire", avatarkey, 7200)
	return userNewInfo.Map()
}
func Videodetails(id string) gdb.Map {
	if len(id) > 20 {
		fields := "code as id,title,created_at as audittime,img_url as img0,user_state,state,video_url as path,video_url as pripath,user_id as operator,user_id,author"
		queryVideo, _ := g.DB().Table("merge_video").Where(g.Map{"code": id, "state": 2}).Fields(fields).One()
		newMap := queryVideo.Map()
		newMap["vidPathObj"] = ""
		newMap["videoType"] = "1"
		newMap["is_othervideo"] = 0
		newMap["excitationTime"] = ""
		newMap["excitationState"] = "1"
		return newMap
	}
	fields := "audittime,title,img0,videoType,is_loangvideo,othervideo_path,is_othervideo,path,user_mini_id as operator,user_mini_id,pripath,custom_ratio,second,excitationTime,template_id,excitation_state as excitationState"
	queryVideo, _ := g.DB().Table("video").Where(g.Map{"code": id, "is_del": 0, "status": 1}).Fields(fields).One()
	newMap := queryVideo.Map()
	newMap["vidPathObj"] = ""
	return newMap
}

//func GetAllListCoreDataFromRedis() interface{} {
//	redisKey := GetRedisKeyrandom2ShowList()
//	r,_ := g.Redis().DoVar("sMembers",redisKey)
//	return r
//}
/**
 * 根据用户刷新次数处理原创集合到第一段数据中
 * @param firstData //第一段数据
 * @param reloadNum //刷新次数
 */
func GetAllListOriginal(firstData interface{}, reloadNum int) interface{} {
	if reloadNum > 3 {
		return firstData
	}
	originalData := GetOriginalData()
	if originalData == nil {
		return firstData
	}
	//if reloadNum == 1 {
	//	newFirstData := until.RemoveDuplicatesAndEmpty(append(gconv.Map(originalData)[:0],firstData))
	//	a := append(newFirstData[:1], newFirstData[2:]...)
	//}else{
	//
	//}
	return firstData
}
func DbProcessData(data interface{}) *glist.List {

	//if len(gconv.Map(data)) == 0 {
	//	return data
	//}
	cropUrl := oss.GetOss("crop")
	test := glist.New()
	for _, value := range gconv.Map(data) {
		var res map[string]interface{}
		err := json.Unmarshal(gconv.Bytes(value), &res)

		if err != nil {
			fmt.Println(err)
		}
		if res != nil {
			if strings.Index(gconv.String(res["img"]), "http") < 0 {
				res["img"] = fmt.Sprintf("%s%s", cropUrl, gconv.String(res["img"]))
			}
			res["title"] = until.DynamicTitle(gconv.String(res["title"]))
			res["type"] = 0
			//return until.DynamicTitle(gconv.String(res["title"]))
			test.PushBack(res)
		}

	}
	return test
}
func GetOriginalData() interface{} {
	originalNumDataRedisKey := GetOriginalNullDataRedisKey()
	originalNumData, _ := g.Redis().Do("get", originalNumDataRedisKey)
	if originalNumData != nil {
		return g.Map{}
	}
	return GetOriginalDataFromRedis()
}
func GetOriginalDataFromRedis() interface{} {
	originalRedisKey := getRedisKeyrandom2NewOriginList()
	r, _ := g.Redis().Do("sRandMember", originalRedisKey, 2)
	return r
}
func SendUserComplate(code string, videoPath string, openid string) interface{} {
	conn := g.Redis().Conn()
	defer conn.Close()

	imgKey := fmt.Sprintf(redisKey.REDIS_API_OPENID_RND_IMG_MAKE, openid,code)
	typeMakeKey := fmt.Sprintf(redisKey.REDIS_API_OPENID_RND_TYPE_MAKE, openid)
	redisKeyMyself := fmt.Sprintf(user.USER_MYSELF_VIDEO_LIST,openid)
	exists, err := conn.DoVar("HEXISTS", typeMakeKey,code)
	if err != nil {
		return err
	}

	var imgUrl string
	if exists.Bool() {
		err := goRedis.InitClient()
		if err != nil {
			fmt.Println(err)
		}
		imgValues, err := goRedis.RedisDb.ZRange(imgKey, 0, 0).Result()
		if err != nil {
			fmt.Println(err)
		}
		imgUrl = imgValues[0]
	}else {
		videoId, _ := g.DB().Table("merge_video").Fields("id").Where("code=?", code).Value()
		imgValue, _ := g.DB().Table("merge_video_imgs").Fields("img_url").Where("video_id=?", videoId).Value()
		imgUrl = imgValue.String()
	}
	//var data map[string]interface{}
	//if imgUrl == "" {
	//	data = g.Map{"video_url": videoPath, "state": 2}
	//} else {
	//	data = g.Map{"video_url": videoPath, "state": 2, "img_url": imgUrl}
	//}
	if exists.Bool() {
		jsonS,err := g.Redis().DoVar("HGET",typeMakeKey,code)
		if err != nil {
			fmt.Println(err)
		}
		var arr map[string]interface{}
		json.Unmarshal(jsonS.Bytes(),&arr)
		arr["video_url"] = videoPath
		arr["state"] = 2
		if imgUrl == "" {
			arr["img_url"] = imgUrl
		}
		newJson,err := json.Marshal(arr)
		if err != nil {
			return err
		}
		conn.Send("MULTI")
		conn.Send("HSET",typeMakeKey,code,newJson)
		conn.Send("RPUSH", REDIS_API_CYGO_API_QUEUE_TABLE, newJson)
		conn.Send("DEL",redisKeyMyself)
		_, err2 := conn.Do("EXEC")
		if err2 != nil {
			return err2
		}
	}
	//else{
	//	r, err := g.DB().Table("merge_video").Data(data).Where("code=?", code).Update()
	//	if err != nil {
	//		return err
	//	}
	//	rows,err := r.RowsAffected()
	//	if err != nil {
	//		return err
	//	}
	//	if rows == 1 {
	//		conn.Do("DEL",redisKeyMyself)
	//	}else{
	//		log.Fatalf("修改数据库信息合成video表错误%d", rows)
	//		return err
	//	}
	//
	//}
	return wechat.SendUserMessage(openid)
}
func DoRedisVideo() error {
	conn := g.Redis().Conn()
	defer conn.Close()
	redisKeyQue := REDIS_API_CYGO_API_QUEUE_TABLE
	exists, err := conn.DoVar("EXISTS", redisKeyQue)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if exists.Bool() {
		var arr map[string]interface{}
		r,err := conn.DoVar("LRANGE",redisKeyQue,0,0)
		if err != nil {
			return err
		}
		arr1 := r.Array()[0]
		err1 := json.Unmarshal(gconv.Bytes(arr1),&arr)
		if err1 != nil {
			fmt.Println(err1)
			return err1
		}

		redisKeyTypeMake := redisKey.GetRedisKeyOpenIdRandTypeMake(gconv.String(arr["openid"]))
		redisKeyImgMake := redisKey.GetRedisKeyOpenIdRandImgMake(gconv.String(arr["openid"]),gconv.String(arr["code"]))
		redisKeyImgContent, err := conn.DoVar("ZRANGE", redisKeyImgMake, 0, 0)
		//conn.DoVar("")
		redisKeyImgAll, err := conn.DoVar("ZRANGE", redisKeyImgMake, 0, -1, "WITHSCORES")
		imgUrlOne := redisKeyImgContent.Array()[0]
		if tx, err := g.DB().Begin(); err == nil {
			// 方法退出时检验返回值，
			// 如果结果成功则执行tx.Commit()提交,
			// 否则执行tx.Rollback()回滚操作。
			defer func() {
				if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()
			// 写入合成视频表
			r, err := tx.Table("merge_video").Insert(g.Map{
				"template_id": arr["template_id"],
				"code":        arr["code"],
				"user_id":     arr["user_id"],
				"state":       arr["state"],
				"created_at":  arr["created_at"],
				"title":      arr["title"],
				"img_url":      imgUrlOne,
				"content":    arr["content"],
				"user_state": arr["user_state"],
				"video_url": arr["video_url"],
				"updated_at": arr["updated_at"],
			})
			if err != nil {
				return err
			}
			videoId, err := r.LastInsertId()
			for i,v := range redisKeyImgAll.Map(){
				// 写入图片
				r, err = tx.Table("merge_video_imgs").Insert(g.Map{
					"video_id":   videoId,
					"img_url":    i,
					"sort":       v,
					"created_at": gtime.Date(),
				})
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			if videoId>0 {
				fmt.Println("删除redis")
				conn.Send("MULTI")
				conn.Send("HDEL",redisKeyTypeMake,arr["code"])
				conn.Send("DEL",redisKeyImgMake)
				conn.Send("Lpop",redisKeyQue)
				_, err := conn.Do("EXEC")
				if err != nil {
					return err
				}
			}
			return nil
		}
		return nil
	}
	fmt.Println("没有redis队列")
	return nil
}
func QueryVideoState(videoCode string,openid string) map[string]interface{} {
	redisKey := REDIS_API_CYGO_API_OPENID_MAKED_VIDEO
	redisKeyType := GetRedisKeyOpenIdRandType(openid)
	redisKeyImg := GetRedisKeyOpenIdRandImg(openid)
	redisKeyMyself := fmt.Sprintf(user.USER_MYSELF_VIDEO_LIST,openid)
	redisKeyPer := fmt.Sprintf(REDIS_API_CODE_VIDEO_PLAN,videoCode)
	videoInfoJson, err := g.Redis().DoVar("HGET", redisKey, videoCode)
	if err != nil {
		fmt.Println(err)
	}
	var videoInfoMap map[string]interface{}
	if videoInfoJson.Bool() {
		time := gtime.Now().Add(+time.Hour*2).Timestamp()
		err = json.Unmarshal(gconv.Bytes(videoInfoJson),&videoInfoMap)
		if err != nil {
			fmt.Println(err)
		}
		videoInfoMap["path"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(fmt.Sprintf("/%s",videoInfoMap["path"])))
		videoInfoMap["img"] = oss.GetPrivateUrl(strconv.FormatInt(time, 10), gconv.String(videoInfoMap["img"]))
		g.Redis().Do("HDEL",redisKey,videoCode)
		g.Redis().Do("Del",redisKeyType)
		g.Redis().Do("Del",redisKeyImg)
		g.Redis().Do("Del",redisKeyMyself)
		//制作百分比percent
		g.Redis().Do("Del",redisKeyPer)

		return videoInfoMap
	}
	return videoInfoMap
}
func QueryVideoPercent(videoCode string) int {

	redisKeyPer := fmt.Sprintf(REDIS_API_CODE_VIDEO_PLAN,videoCode)
	//_,err1 := g.Redis().Do("ZADD",redisKeyPer,30,1)
	//if err1 != nil{
	//	glog.Line().Println(err1)
	//}
	highestScore,err := g.Redis().DoVar("ZRANGE",redisKeyPer,-1,-1,"WITHSCORES")
	if err != nil{
		glog.Line().Println(err)
	}
	if highestScore.Bool() {
		return gconv.Int(highestScore.Array()[1])
	}
	return 10
}
/**
 * 从redis中获取今天最新数据取值条数
 */
func GetTodayNumFromRedis() int {
	toDayNumRedisKey := GetRedisKeyrandom2NewNumber()
	toDayData, _ := g.Redis().Do("get", toDayNumRedisKey)
	return gconv.Int(gconv.String(toDayData))
}
func GetRedisKeyrandom2NewNumber() string {
	return fmt.Sprintf(REDIS_API_RANDOM2_NEW_NUMBER, APP_ID, gtime.Now().Format("Y-m-d"))
}
func GetRedisKeyrandom2NewList() string {
	return fmt.Sprintf(REDIS_API_RANDOM2_NEW_LIST, APP_ID, gtime.Now().Format("Y-m-d"))
}
func getRedisKeyrandom2NewOriginList() string {
	return fmt.Sprintf(REDIS_API_RANDOM2_NEW_ORIGIN_LIST, APP_ID, gtime.Now().Format("Y-m-d"))
}
func GetOriginalNullDataRedisKey() string {
	return fmt.Sprintf(REDIS_API_KEY_ORIGINAL_DATA_ISNULL, APP_ID)
}
func GetRedisKeyFlowOpenNum(date string, uid string, appId string) string {
	return fmt.Sprintf(REDIS_API_KEY_NEW_FLOW_OPEN_NUM, date, uid, appId)
}
func GetRedisKeyrandom2NearNumber() string {
	return fmt.Sprintf(REDIS_API_RANDOM2_NEAR_NUMBER, APP_ID, gtime.Now().Format("Y-m-d"))
}
func GetRedisKeyrandom2ShowList() string {
	return fmt.Sprintf(REDIS_API_RANDOM2_SHOW_LIST, APP_ID, gtime.Now().Format("Y-m-d"))
}
func GetRedisCyapiDetails(id string) string {
	if len(id) > 20 {
		return fmt.Sprintf(REDIS_API_DETAILS_USER, id)
	}
	return fmt.Sprintf(REDIS_API_DETAILS, id)
}
func GetRedisCyapiDetailsOaut(id string) string {
	if len(id) > 20 {
		return fmt.Sprintf(REDIS_API_DETAILS_OAUTHORITY_USER, id)
	}
	return fmt.Sprintf(REDIS_API_DETAILS_OAUTHORITY, id)
}
func GetRedisKeyVideoCount(id string) string {
	return fmt.Sprintf(REDIS_API_VIDEO_VIDEOSCOUNT_WHITE, id)
}
func GetRedisKeyOpenIdRandType(openid string) string {
	return fmt.Sprintf(REDIS_API_OPENID_RND_TYPE, openid)
}
func GetRedisKeyOpenIdRandImg(openid string) string {
	return fmt.Sprintf(REDIS_API_OPENID_RND_IMG, openid)
}
