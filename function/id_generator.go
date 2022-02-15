package function

import (
	"math/big"

	"github.com/google/uuid"
)

func Id() string {
	uuid, _ := uuid.NewRandom()
	return toBase62(uuid)
}

func toBase62(uuid uuid.UUID) string {
	var i big.Int
	i.SetBytes(uuid[:])
	return i.Text(62)[:11]
}
