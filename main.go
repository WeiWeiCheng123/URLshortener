package main

import (
	"context"
	"fmt"

	"time"

	//"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/go-redis/redis/v8"
)

var RC *redis.Client
var ctx = context.Background()

func main() {

	RC = store.NewClient()
	var localLocation *time.Location
	localLocation, _ = time.LoadLocation("Asia/Shanghai")

	layout := "2006-01-02T15:04:05Z"
	str := "2022-02-15T09:20:41Z"
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
	}
	t = t.In(localLocation)
	fmt.Println(t)
	xx, _ := store.Save(RC, "GOOGLE", t)
	fmt.Println(xx)
	xx, _ = store.Save(RC, "DCARD", t)
	fmt.Println(xx)
	xx, _ = store.Save(RC, "YAHOO", t)
	fmt.Println(xx)
	fmt.Println(store.Load(RC, "v2dBmrUHmW8"))
}
