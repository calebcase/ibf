package ibf

import (
	"encoding/json"
	"math/big"

	"github.com/dchest/siphash"
)

// Interface

type Hasher interface {
	Hash(key *big.Int) uint64

	json.Marshaler
	json.Unmarshaler
}

// Hash

type hash struct {
	Key [2]uint64 `json:"key"`
}

type Hash struct {
	p hash
}

var _ Hasher = (*Hash)(nil)

func NewHash(key0, key1 uint64) *Hash {
	return &Hash{hash{[2]uint64{key0, key1}}}
}

func (self *Hash) Hash(key *big.Int) uint64 {
	return siphash.Hash(self.p.Key[0], self.p.Key[1], key.Bytes())
}

func (self *Hash) MarshalJSON() ([]byte, error) {
	return json.Marshal(&self.p)
}

func (self *Hash) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &self.p)
}
