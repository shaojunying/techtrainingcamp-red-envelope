package main

import (
	"github.com/spf13/viper"
	"log"
	"net/http"
	_ "net/http/pprof"
	"red_envelope/config"
	"red_envelope/database"
	"red_envelope/routers"
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

	r := routers.InitRouter()

	port := viper.GetString("server.port")
	r.Run(":" + port)
}
