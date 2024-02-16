package app

import (
	"gf-app/app/service/video"
	"github.com/gogf/gf/os/gtimer"
	"time"
)

func StartSchedule() {
	timer := gtimer.New(10, 10*time.Millisecond)

	interval := time.Second
	gtimer.AddSingleton(interval, func() {
		video.DoRedisVideo()
	})
	timer.Start()

}

