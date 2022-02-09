package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/WeiWeiCheng123/URLshortener/store"
	"github.com/go-redis/redis/v8"
)

var RC *redis.Client

func main() {

	RC = store.NewClient()
	fmt.Println(store.CheckId(RC, 1))
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
	date := "2021-02-08T09:20:41Z"
	layout := "2006-01-02T15:04:05.000Z"
	t, err := time.Parse(layout, date)
	xx, _ := store.Save(RC, "152165464", t)
	fmt.Println(xx)
}
