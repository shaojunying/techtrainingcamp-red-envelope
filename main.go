package main

import (
	"github.com/spf13/viper"
	"red_envelope/config"
	"red_envelope/database"
	"red_envelope/routers"
	"github.com/gin-contrib/pprof" // 性能分析使用，请在正式版本移除
)

func main() {

	//读取配置
	config.InitConf()

	//启动数据库
	db := database.InitDB()
	defer db.Close()

	r := routers.InitRouter()
	pprof.Register(r)

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
