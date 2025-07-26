package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func InitConf() {
	wordDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(wordDir)
	err := viper.ReadInConfig()
	if err != nil {
		panic("load config failed: " + err.Error())
	}
	fmt.Println("load config success")
}
