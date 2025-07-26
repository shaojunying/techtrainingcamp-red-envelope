package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	"red_envelope/internal/app"
	"red_envelope/internal/infrastructure/database"
	"red_envelope/internal/interface/http/router"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//读取配置
	app.InitConf()

	//启动数据库
	//db := database.InitDB()
	//defer db.Close()
	database.InitRedis()

	err := database.InitMQ()
	if err != nil {
		log.Println("Mq init error.")
		return
	}
	defer database.CloseMQ()

	r := router.InitRouter()

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
