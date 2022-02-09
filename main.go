package main

import (
	"fmt"
	"math/rand"

	
	"github.com/go-redis/redis/v8"
	"github.com/WeiWeiCheng123/URLshortener/function"
	"github.com/WeiWeiCheng123/URLshortener/store"
)

var RC *redis.Client

func main() {

	RC = store.NewClient()
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
}
