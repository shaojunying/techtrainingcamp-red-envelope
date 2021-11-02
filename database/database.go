package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("database.driverName")
	host := viper.GetString("database.host")
	userName := viper.GetString("database.userName")
	password := viper.GetString("database.password")
	database := viper.GetString("database.database")
	// port
	args := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", userName, password, host, database)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database, err: " + err.Error())
	}

	MaxIdleConns := viper.GetInt("database.MaxIdleConns")            //空闲时最大连接数
	MaxOpenConns := viper.GetInt("database.MaxOpenConns")            //数据库最大连接数
	ConnMaxLifeTime := viper.GetDuration("database.ConnMaxLifeTime") //单连接最长生命周期
	db.DB().SetMaxIdleConns(MaxIdleConns)
	db.DB().SetMaxOpenConns(MaxOpenConns)
	db.DB().SetConnMaxLifetime(ConnMaxLifeTime)
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
