package store

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
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

func CheckId(rdb *redis.Client, id uint64) bool {
	n, err := rdb.Exists(ctx, strconv.FormatUint(id, 10)).Result()
	if err != nil {
		panic(err)
	}

	if n > 0 {
		return true
	} else {
		return false
	}
}

func Save(rdb *redis.Client, url string, exp time.Time) (string, error) {
	var id uint64
	for exist := true; exist; exist = CheckId(rdb, id) {
		id = rand.Uint64()
	}
	fmt.Println(id)
	shortURL := Data{id, url, exp.Format("2006-01-02 15:04:05")}
	fmt.Print("shortURL")
	fmt.Println(shortURL)
	rdb.HMSet(ctx, strconv.FormatUint(id, 10))
	rdb.ExpireAt(ctx, strconv.FormatUint(id, 10), exp)

	return function.Encode(id), nil
}
