package routers

import (
	"github.com/gin-gonic/gin"
	"red_envelope/api/snatch"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	setUpRouter(router)
	return router
}

// 设置路由
func setUpRouter(router *gin.Engine) {
	api := router.Group("/")
	snatch.RegisterRedEnvelopeRouter(api.Group("/snatch"))
}
