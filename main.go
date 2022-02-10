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
	/*
		//var a uint64
		var b uint64
		//a = rand.Uint64()
		for i := 1; i <= 5; i++ {
			b = rand.Uint64()
		}
		fmt.Println(b)
		c := function.Encode(b)
		fmt.Println(c)
		s := "SSASFASDFFAW"
		d, err := function.Decode(s)
		fmt.Println(d)
		fmt.Println(err)
	*/
	layout := "2006-01-02T15:04:05Z"
	str := "2021-02-15T09:20:41Z"
	t, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	xx, _ := store.Save(RC, "152165464", t)
	fmt.Println(xx)
	xx, _ = store.Save(RC, "15216541264", t)
	fmt.Println(xx)
	xx, _ = store.Save(RC, "15216546FQW4", t)
	fmt.Println(xx)
	dec, _ := function.Decode("EJlQv3wKYD6")
	fmt.Println(dec)
	fmt.Println(RC.HMSet(ctx, "key1", "val"))
	fmt.Println(RC.Keys(ctx, "1"))
	fmt.Println(RC.Do(ctx, "hset", "KEY", "VAL").Result())
	fmt.Println(RC.Get(ctx, "KEY"))
}
