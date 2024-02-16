package until

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/noaway/dateparse"
	"math"
	"regexp"
	"strings"
	"time"
)
//TimeContentId集中营key
const CY_API_NEW_TIMEPERIOD = "cykj_api.timePeriod.note"

func GetAppid(referer string) string {
	appid := strings.Split(referer,"/")[3]
	if len(appid) >0 {
		return appid
	}
	glog.Debug(appid)
	return ""
}
func RemoveDuplicatesAndEmpty(a []string) (ret []string){
	aLen := len(a)
	for i:=0; i < aLen; i++{
		if (i > 0 && a[i-1] == a[i]) || len(a[i])==0{
			continue
		}
		ret = append(ret, a[i])
	}
	return
}
func DynamicTitle(title string) string {
	pattern := "\\#.*?\\#"
	rm := regexp.MustCompile(pattern) //正则规则匹配
	match := rm.FindAllString(title,5)
	match1 := rm.MatchString(title)
	//fmt.Println(match)
	if match1 == true {
		data := RemoveDuplicatesAndEmpty(match)

		for _,v := range data{
			titleKey := gconv.String(strings.Split(v,"#")[1])
			var titleV string
			if (strings.Contains(v, "-")){
				titleKeyh := strings.Split(gconv.String(strings.Split(v,"#")[1]),"-")
				titleKey = titleKeyh[0]
				titleV = titleKeyh[1]
			}

			//return gconv.String(GetTimeNoteFromRedis())
			//titleValue := titleKey[1]
			var value string
			//return gconv.String(titleKey)
			switch (gconv.String(titleKey)) {
				case "time":
					value = GetTimeNoteFromRedis()
					break
				case "day":
					t, _ := dateparse.ParseAny(fmt.Sprintf("%s%s",gconv.String(gtime.Now().Format("Y")),titleV))
					day := gconv.Int(math.Ceil(float64(t.Add(-time.Hour*8).Unix()-gtime.Now().Unix())/86400))
					dayName := map[int]string{0: "前天", 1: "昨天",2: "今天",3: "明天",4: "后天"}
					if day <= 2 && day >= -2 {
						value = dayName[int(day) + 2]
					} else if day > 2 {
						value = fmt.Sprintf("%d天后",day)
					} else if day < -60 {
						day += 365
						day += IsLeapYear(gconv.Int(gtime.Now().Format("Y"))+1)
						value = fmt.Sprintf("%d天后",day)
					} else if day < 1 {
						value = fmt.Sprintf("%v天前",math.Abs(float64(day)))
					}
					break
				case "week":
					weeks := map[int]string{0: "日", 1: "一",2: "二",3: "三",4: "四",5: "五",6: "六"}
					value = fmt.Sprintf("周%s",weeks[int(gtime.Now().Weekday())])
					break
				case "date":
					value = fmt.Sprintf("%d月%d号",gtime.Now().Month(),gtime.Now().Day())
					break
				case "distance":
					t, _ := dateparse.ParseAny(titleV)
					tend,_:= dateparse.ParseAny(gconv.String(gtime.Now().Format("Ymd")))
					value = gconv.String(math.Abs(float64((tend.Add(-time.Hour*8).Unix()-t.Add(-time.Hour*8).Unix())/86400)))
					break
					default:
					value = v

			}

			title = strings.Replace(title, v, value, -1 )

		}
		return title

	}
	return title
}

func IsLeapYear(year int) int {
	if (year%4 == 0 && year%100 != 0) || (year%400 == 0) ||(year%1000==0){
		return 1
	}
	return 0
}

func GetTimeNoteFromRedis() string {
	redisKey := CY_API_NEW_TIMEPERIOD
	hour := gconv.Int(time.Now().Hour())
	var r string
	if  hour == 24 {
		r = gconv.String(fmt.Sprintf("24%s",gconv.String(gtime.Now().Format("i"))))
	}else{
		r = gconv.String(gtime.Now().Format("Hi"))
	}
	rr,_ := g.Redis().DoVar("ZREVRANGEBYSCORE",redisKey,r,0,"WITHSCORES")
	return gconv.String(rr.Slice()[0])
}
func getName(params ...interface{}) string {
	var stringSlice []string
	for _, param := range params {
		stringSlice = append(stringSlice, param.(string))
	}
	return strings.Join(stringSlice, "_")
}