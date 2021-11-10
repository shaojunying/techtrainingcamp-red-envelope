package main

import (
	"log"
	"net/http"
	"red_envelope/config"
	"red_envelope/database"
	"red_envelope/routers"

	"github.com/spf13/viper"
)


func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()
	database.InitRedis()

	err := database.InitMQ()
	if err != nil {
		log.Println("Mq init error.")
		return
	}
	defer database.CloseMQ()

	r := routers.InitRouter()

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
