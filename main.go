package main

import (
	"github.com/WeiWeiCheng123/URLshortener/function"
	"fmt"
	"math/rand"
)

func main() {
	//var a uint64
	var b uint64
	//a = rand.Uint64()
	for i:=1; i<=5; i++{
		b = rand.Uint64()
	}
	fmt.Println(b)
	c := function.Encode(b)
	fmt.Println(c)

}