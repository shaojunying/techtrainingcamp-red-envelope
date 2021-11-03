package routers

import (
	"github.com/gin-gonic/gin"
	"red_envelope/api/redenvelope"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	setUpRouter(router)
	return router
}

// 设置路由
func setUpRouter(router *gin.Engine) {
	api := router.Group("/")
	redenvelope.RegisterRedEnvelopeRouter(api.Group("/redenvelope"))
}
