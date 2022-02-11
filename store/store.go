package store

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Data struct {
	Id          uint64    `json:"id" redis:"id"`
	OriginalURL string    `json:"OriginalURL" redis:"OriginalURL"`
	ExpTime     time.Time `json:"exp" redis:"exp"`
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
func Save(rdb *redis.Client, url string, expireTime time.Time) (string, error) {
	var id uint64

	for exist := true; exist; exist = CheckId(rdb, id) {
		id = rand.Uint64()
	}

	err := rdb.Set(ctx, strconv.FormatUint(id, 10), url, 0).Err()
	//Error to save data
	if err != nil {
		return "", err
	}

	res, err := rdb.ExpireAt(ctx, strconv.FormatUint(id, 10), expireTime).Result()
	//Error to set expire time
	if err != nil {
		return "", err
	}
	if res {
		return function.Encode(id), nil
	} else {
		return "", errors.New("Fail to Set")
	}
}

//Give shortURL return original URL
func Load(rdb *redis.Client, shortURL string) (string, error) {
	id, err := function.Decode(shortURL)

	if err != nil {
		return "", err
	}

	url, err := rdb.Get(ctx, strconv.FormatUint(id, 10)).Result()

	if err == redis.Nil {
		return "", err
	}
	return url, nil
}
