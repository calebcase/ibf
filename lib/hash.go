package ibf

import (
	"github.com/dchest/siphash"
)

// Hash maintains the state for a siphash hasher.
type Hash struct {
	Key [2]uint64 `json:"key"`
}

// NewHash returns a new siphash hasher.
func NewHash(key0, key1 uint64) *Hash {
	return &Hash{
		Key: [2]uint64{key0, key1},
	}
}

// Hash retuns the digest of the value.
func (h *Hash) Hash(value []byte) (digest uint64) {
	return siphash.Hash(h.Key[0], h.Key[1], value)
}
