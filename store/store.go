package store

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Data struct {
	Id          uint64 `json:"id" redis:"id"`
	OriginalURL string `json:"OriginalURL" redis:"OriginalURL"`
	ExpTime     string `json:"exp" redis:"exp"`
}

func NewClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connect!")

	return rdb
}

func CheckId(rdb *redis.Client) bool {
	n, err := rdb.Exists(ctx, "key1").Result()
	if err != nil {
		panic(err)
	}

	if n > 0 {
		fmt.Println("存在")
		return true
	} else {
		fmt.Println("不存在")
		return false
	}

}
