package function

import (
	"errors"
	"math"
	"strings"
)

const (
	//total length is 62 (a~z + A~Z + 0~9)
	charTable = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length = uint64(62)
)

func Encode(number uint64) string {
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(11)

	for ; number > 0; number = number / length {
	   encodedBuilder.WriteByte(charTable[(number % length)])
	}
  
	return encodedBuilder.String()
}
  
func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, char := range encoded {
	   alphabeticPosition := strings.IndexRune(charTable, char)

	   if alphabeticPosition == -1 {
		  return uint64(alphabeticPosition), errors.New("invalid character: " + string(char))
	   }
	   number += uint64(alphabeticPosition) * uint64(math.Pow(float64(length), float64(i)))
	}
	return number, nil
}