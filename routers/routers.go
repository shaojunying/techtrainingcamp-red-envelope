package routers

import (
	"log"
	"net/http"
	"red_envelope/api/redenvelope"
	"red_envelope/middleware"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(Cors()) //开启中间件 允许使用跨域请求
	//测试阶段，先将令牌桶注释
	//加入限制器，限制能通过的最大流量，多余流量将被舍弃
	limitRate := int64(viper.GetInt("limitRate"))
	limitCapacity := int64(viper.GetInt("limitCapacity"))
	router.Use(middleware.RateLimitMiddleware(time.Second, limitRate, limitCapacity))

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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持POST方法
			c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
