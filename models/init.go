package models

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//    golang中程序首先加载依赖，加载依赖的顺序为率先加载不依赖别的依赖的依赖，之后加载全局变量，然后执行init()函数，
//   最后是main函数，在函数中加载变量的顺序同加载依赖的方法相同，率先加载不依赖别的变量的变量。

var DB = Init()
var RDB = InitRedisDB()

func Init() *gorm.DB {
	dsn := "root:123456@tcp(192.168.150.132:31234)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// log.Println("gorm Init Error: ", err)
		panic("gorm Init Error: " + err.Error())
	}
	return db
}

func InitRedisDB() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "192.168.150.132:30103",
		Password: "test", // no password set
		DB:       0,      // use default DB
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic("Failed to connect to Redis server:" + err.Error())
	}

	return client
}
