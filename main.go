package main

import (
	"context"
	"fmt"

	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
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
	xx, _ := store.Save(RC, "152165464", t)
	fmt.Println(xx)
	xx, _ = store.Save(RC, "15216541264", t)
	fmt.Println(xx)
	xx, _ = store.Save(RC, "15216546FQW4", t)
	fmt.Println(xx)
	dec, _ := function.Decode("EJlQv3wKYD6")
	fmt.Println(dec)
	val, err := RC.Set(ctx, "5577006791947779410", "GOOGLE", 0).Result()
	fmt.Printf("val = %s \n", val)

}
