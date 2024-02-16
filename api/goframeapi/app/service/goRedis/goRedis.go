package goRedis

import (
	"github.com/go-redis/redis"
	"github.com/gogf/gf/frame/g"
)

// 声明一个全局的redisDb变量
var RedisDb *redis.Client

// 根据redis配置初始化一个客户端
func InitClient() (err error) {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:g.Cfg("config").GetString("redis.Addr"), // redis地址
		Password:g.Cfg("config").GetString("redis.Password"),               // redis密码，没有则留空
		DB:g.Cfg("config").GetInt("redis.DB"),                // 默认数据库，默认是0
	})

	//通过 *redis.Client.Ping() 来检查是否成功连接到了redis服务器
	_, err = RedisDb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

