package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"red_envelope/api/redenvelope"
	"red_envelope/middleware"
	"time"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	//加入限制器，限制能通过的最大流量，多余流量将被舍弃
	limitRate := int64(viper.GetInt("limitRate"))
	limitCapacity := int64(viper.GetInt("limitCapacity"))
	router.Use(middleware.RateLimitMiddleware(time.Second, limitRate, limitCapacity))
	router.Use(middleware.ConfigLoadingMiddleware())
	setUpRouter(router)
	return router
}

// 设置路由
func setUpRouter(router *gin.Engine) {
	api := router.Group("/")
	redenvelope.RegisterRedEnvelopeRouter(api.Group("/redenvelope"))
}
