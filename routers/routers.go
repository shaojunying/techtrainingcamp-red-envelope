package routers

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"red_envelope/api/redenvelope"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	//测试阶段，先将令牌桶注释
	//加入限制器，限制能通过的最大流量，多余流量将被舍弃
	//limitRate := int64(viper.GetInt("limitRate"))
	//limitCapacity := int64(viper.GetInt("limitCapacity"))
	//router.Use(middleware.RateLimitMiddleware(time.Second, limitRate, limitCapacity))
	setUpRouter(router)
	return router
}

// 设置路由
func setUpRouter(router *gin.Engine) {
	pprof.Register(router) // 注册pprof路由
	api := router.Group("/")
	redenvelope.RegisterRedEnvelopeRouter(api.Group("/redenvelope"))
}
