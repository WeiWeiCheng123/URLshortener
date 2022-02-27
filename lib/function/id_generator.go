package function

import (
	"math/big"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	//total length is 62 (0~9 + a~z + A~Z)
	charTable = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length    = uint64(62)
)

//Generate a unique id and convert it into base62 type
func Id() string {
	uuid, _ := uuid.NewRandom()
	return toBase62(uuid)
}

//Convert id into base62 type
func toBase62(uuid uuid.UUID) string {
	var i big.Int
	i.SetBytes(uuid[:])
	return i.Text(62)[:7]
}

func Generator() string {
	t := uint64(time.Now().Unix())
	rand.Seed(time.Now().UnixNano())
	r_1 := rand.Uint64()
	rand.Seed(time.Now().Unix())
	r_2 := uint64(rand.Int())
	t = t + r_1 + r_2

	return encode(t)
}

func encode(num uint64) string {
	var encoder strings.Builder
	encoder.Grow(7)

	for ; num > 0; num = num / length {
		encoder.WriteByte(charTable[(num % length)])
	}

	return encoder.String()[:7]
}
