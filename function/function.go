package function

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"
)

const (
	//total length is 62 (0~9 + a~z + A~Z)
	charTable = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length    = uint64(62)
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
func IsUrl(OriginalURL string) bool {
	u, err := url.Parse(OriginalURL)
	return err == nil && u.Scheme != "" && u.Host != ""
}

//Time formater
func TimeFormat(expTime string) (time.Time, error) {
	var localLocation *time.Location
	localLocation, _ = time.LoadLocation("Asia/Shanghai")
	layout := "2006-01-02T15:04:05Z"
	expireTime, err := time.Parse(layout, expTime)
	//Time format error
	if err != nil {
		return time.Time{}, err
	}

	expireTime = expireTime.In(localLocation)
	fmt.Println("EXP TIME = ", expTime)
	return expireTime, nil
}
