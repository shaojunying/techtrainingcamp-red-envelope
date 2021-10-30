package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"red_envelope/config"
)

func main() {

	config.InitConf()

	port := viper.GetString("server.port")
	engine := gin.Default()
	engine.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg": "test",
		})
	})
	engine.Run(":" + port)
}
