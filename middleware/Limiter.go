package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

// RateLimitMiddleware 限流器 超出桶的流量会被舍弃
func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		//if bucket.TakeAvailable(1) < 1 {
		//	c.String(http.StatusForbidden, "rate limit...")
		//	c.Abort()
		//	return
		//}
		count := 0 // 用来对等待次数做记录
		//桶中没有令牌时程序会休眠0.2~0.6s 若还没有则return
		//后期可以根据情况调整为上面的if语句
		for bucket.TakeAvailable(1) < 1 {
			time.Sleep(time.Duration(200) * time.Millisecond)
			count++
			if count >= 3 {
				c.String(http.StatusForbidden, "rate limit...")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
