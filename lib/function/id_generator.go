package function

import (
	"math/big"

	"github.com/google/uuid"
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
