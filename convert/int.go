package convert

import (
	"encoding/binary"

	"golang.org/x/exp/constraints"
)

// IntToBytes converts a signed or unsigned integer into 8 bytes
func IntToBytes[T constraints.Integer](i T) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(i))
	return bz
}
