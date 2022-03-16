package function

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	//total length is 62 (0~9 + a~z + A~Z)
	charTable = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length    = uint64(62)
)

//generate a unique id and convert it into base62 type
func Generator() string {
	t := uint64(time.Now().Unix())
	rand.Seed(time.Now().UnixNano())
	r_1 := rand.Uint64()
	rand.Seed(time.Now().UnixNano() + int64(r_1))
	r_2 := uint64(rand.Int())
	t = t + r_1 + r_2
	if t % 2 == 1 {
		t = t + 1
	}
	return encode(t)
}

// give an integer return based62 string
func encode(num uint64) string {
	var encoder strings.Builder
	encoder.Grow(7)

	for ; num > 0; num = num / length {
		encoder.WriteByte(charTable[(num % length)])
	}

	return encoder.String()[:7]
}

// give base62 string return an integre
func decode(encoded string) uint64 {
	var number uint64

	for i, char := range encoded {
		charPosition := strings.IndexRune(charTable, char)
		if charPosition == -1 {
			return uint64(charPosition)
		}
		number += uint64(charPosition) * uint64(math.Pow(float64(length), float64(i)))
	}

	return number
}
