package app

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConf() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wordDir + "/configs")  // 更新配置文件路径
	viper.AddConfigPath(wordDir)               // 兼容旧路径
	err := viper.ReadInConfig()
	if err != nil {
		panic("load config failed: " + err.Error())
	}
	fmt.Println("load config success")
}
