package middleware

import (
	"log"
	"red_envelope/api/redenvelope"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
)

func CheatPreventing(c *gin.Context, milliseconds int64) {
	//匹配参数
	r := redenvelope.CheatPreventing{}
	err := c.ShouldBindBodyWith(&r, binding.JSON)
	if err != nil {
		// 匹配不到说明没有uid，可能是别的接口
		return
	}
	if r.UID == nil {
		redenvelope.HandleERR(c, 102, errors.New("user.UID is nil"))
		return
	}
	// 获取该uid上一次请求的时间
	lastTime, err := redenvelope.Mapper.GetLastRequestTime(c, *r.UID)
	if err != nil {
		log.Println("获取上一次请求时间失败, err: ", err)
		c.Next()
		return
	}
	// time.Now().UnixNano()/1_000_000 获取的是毫秒时间戳（13位）
	now := time.Now().UnixNano() / 1_000_000
	if lastTime != 0 && lastTime+milliseconds > now {
		redenvelope.HandleERR(c, 103, errors.New("请求过于频繁"))
		c.Abort()
		return
	}
	// 更新该uid的请求时间
	err = redenvelope.Mapper.UpdateLastRequestTime(c, *r.UID, now, milliseconds)
	if err != nil {
		log.Println("更新uid的请求时间失败, err: ", err)
	}
}

//CheatPreventingMiddleware 防作弊中间件
func CheatPreventingMiddleware(milliseconds int64) gin.HandlerFunc {

	return func(c *gin.Context) {
		CheatPreventing(c, milliseconds)
		c.Next()
	}
}
