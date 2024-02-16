package redisKey

import "fmt"

//上传图片点击制作
const REDIS_API_OPENID_RND_TYPE_MAKE = "cygo_api.%s.1"
const REDIS_API_OPENID_RND_IMG_MAKE = "cygo_api.%s.%s.make"
func GetRedisKeyOpenIdRandTypeMake(openid string) string {
	return fmt.Sprintf(REDIS_API_OPENID_RND_TYPE_MAKE, openid)
}
func GetRedisKeyOpenIdRandImgMake(openid string, code string) string {
	return fmt.Sprintf(REDIS_API_OPENID_RND_IMG_MAKE, openid, code)
}