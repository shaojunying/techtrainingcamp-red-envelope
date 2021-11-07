package database

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

var Rdx *redis.Client

func InitRedis() {
	addr := viper.GetString("redis.addr")
	username := viper.GetString("redis.username")
	password := viper.GetString("redis.password")
	dbNumber := viper.GetInt("redis.dbNumber")

	rdx := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       dbNumber,
	})
	// 测试redis是否可以正常连接
	ctx := context.Background()

	if err := rdx.Ping(ctx).Err(); err != nil {
        panic("failed to connect to redis server, err: " + err.Error())
    }

	Rdx = rdx
}

func GetRdx() *redis.Client {
	return Rdx
}
