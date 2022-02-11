package function

import (
	"errors"
	"math"
	"strings"
	"net/url"
)

const (
	//total length is 62 (0~9 + a~z + A~Z)
	charTable = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length = uint64(62)
)

func Encode(number uint64) string {
	var encoder strings.Builder
	encoder.Grow(11)

	for ; number > 0; number = number / length {
		encoder.WriteByte(charTable[(number % length)])
	}

	return encoder.String()
}

func Decode(encoded string) (uint64, error) {
	var number uint64

	for i, char := range encoded {
	   charPosition := strings.IndexRune(charTable, char)

	   if charPosition == -1 {
		  return uint64(charPosition), errors.New("Invalid character: " + string(char))
	   }
	   number += uint64(charPosition) * uint64(math.Pow(float64(length), float64(i)))
	}
	return number, nil
}

//Valid the URL, if URL not correct then return false
func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}