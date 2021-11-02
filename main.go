package main

import (
	"github.com/spf13/viper"
	"red_envelope/config"
	"red_envelope/database"
	"red_envelope/routers"
)

func main() {

	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()

	r := routers.InitRouter()

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
