package main

import (
	"function/encoder"
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
	c := Encode(b)
	fmt.Println(c)

}