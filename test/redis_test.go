package test

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var count int = 1
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "192.168.150.132:30103",
	Password: "test", // no password set
	DB:       0,      // use default DB
})

func TestRedisSet(t *testing.T) {
	rdb.Set(ctx, "name", "mmc", time.Second*100)
}

func TestRedisGet(t *testing.T) {
	v, err := rdb.Get(ctx, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)
}
