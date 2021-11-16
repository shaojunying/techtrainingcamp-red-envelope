package routers

import (
	"red_envelope/api/redenvelope"
	"red_envelope/middleware"

	//"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	//测试阶段，先将令牌桶注释
	//加入限制器，限制能通过的最大流量，多余流量将被舍弃
	//limitRate := int64(viper.GetInt("limitRate"))
	//limitCapacity := int64(viper.GetInt("limitCapacity"))
	//router.Use(middleware.RateLimitMiddleware(time.Second, limitRate, limitCapacity))

	milliseconds := int64(viper.GetInt("cheat-preventing.milliseconds"))
	if milliseconds == 0 {
		milliseconds = 1000
	}
	setUpRouter(router, milliseconds)
	return router
}

// 设置路由
func setUpRouter(router *gin.Engine, milliseconds int64) {
	//pprof.Register(router) // 注册pprof路由
	otherApi := router.Group("/")

	// 这组路由不走防作弊检查
	redenvelope.RegisterOtherRouter(otherApi.Group("/redenvelope"))

	router.Use(middleware.CheatPreventingMiddleware(milliseconds))
	router.Use(middleware.ConfigLoadingMiddleware())
	coreApi := router.Group("/")
	redenvelope.RegisterRedEnvelopeRouter(coreApi.Group("/redenvelope"))
}
