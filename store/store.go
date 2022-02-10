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


//Check the id is exist or not
func CheckId(rdb *redis.Client, id uint64) bool {
	err := rdb.Get(ctx, strconv.FormatUint(id, 10)).Err()

	if err == redis.Nil {
		return false
	} else if err != nil {
		panic(err)
	} else {
		return true
	}
}

//Give original URL and expire time, save to Redis
func Save(rdb *redis.Client, url string, exp string) (string, error) {
	var id uint64
	var localLocation *time.Location
	localLocation, _ = time.LoadLocation("Asia/Shanghai")

	for exist := true; exist; exist = CheckId(rdb, id) {
		id = rand.Uint64()
	}

	layout := "2006-01-02T15:04:05Z"
	expireTime, err := time.Parse(layout, exp)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	expireTime = expireTime.In(localLocation)

	fmt.Println(id)
	shortURL := Data{id, url, expireTime.Format("2006-01-02 15:04:05")}
	fmt.Println("shortURL")
	fmt.Println(shortURL)
	err = rdb.Set(ctx, strconv.FormatUint(id, 10), url, 0).Err()

	if err != nil {
		return "", err
	}

	res, err := rdb.ExpireAt(ctx, strconv.FormatUint(id, 10), expireTime).Result()

	if err != nil {
		return "", err
	}
	if res {
		fmt.Println("Set")
	} else {
		fmt.Println("Fall to set")
	}

	return function.Encode(id), nil
}

//Give shortURL return original URL
func Load(rdb *redis.Client, shortURL string) (string, error) {
	id, err := function.Decode(shortURL)

	if err != nil {
		return "", err
	}

	url, err := rdb.Get(ctx, strconv.FormatUint(id, 10)).Result()
	if err != nil {
		return "", err
	}
	return url, nil
}
