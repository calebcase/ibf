package ibf

import (
	"encoding/json"
	"math/big"
)

// Interface

type Celler interface {
	Insert(key *big.Int, hash uint64)
	Remove(key *big.Int, hash uint64)
	Subtract(cell Celler)
	Invert()

	Clone() Celler

	GetId() *big.Int
	GetHash() uint64
	GetCount() int64

	IsEmpty() bool
	IsPure(hasher Hasher) bool

	json.Marshaler
	json.Unmarshaler
}

// Cell

type cell struct {
	Id    *big.Int `json:"id"`
	Hash  uint64   `json:"hash"`
	Count int64    `json:"count"`
}

type Cell struct {
	p cell
}

func NewCell() *Cell {
	return &Cell{cell{
		Id:    big.NewInt(0),
		Hash:  0,
		Count: 0,
	}}
}

func (self *Cell) Insert(key *big.Int, hash uint64) {
	self.p.Id.Xor(self.p.Id, key)
	self.p.Hash = self.p.Hash ^ hash
	self.p.Count += 1
}

func (self *Cell) Remove(key *big.Int, hash uint64) {
	self.p.Id.Xor(self.p.Id, key)
	self.p.Hash = self.p.Hash ^ hash
	self.p.Count -= 1
}

func (self *Cell) Subtract(cell Celler) {
	self.p.Id.Xor(self.p.Id, cell.GetId())
	self.p.Hash = self.p.Hash ^ cell.GetHash()
	self.p.Count -= cell.GetCount()
}

func (self *Cell) Invert() {
	self.p.Count *= -1
}

func (self *Cell) Clone() Celler {
	return &Cell{cell{
		Id:    big.NewInt(0).Set(self.p.Id),
		Hash:  self.p.Hash,
		Count: self.p.Count,
	}}
}

func (self *Cell) GetId() *big.Int {
	return big.NewInt(0).Set(self.p.Id)
}

func (self *Cell) GetHash() uint64 {
	return self.p.Hash
}

func (self *Cell) GetCount() int64 {
	return self.p.Count
}

func (self *Cell) IsEmpty() bool {
	if self.p.Count == 0 && self.p.Id.Cmp(ZERO) == 0 && self.p.Hash == 0 {
		return true
	}
	return false
}

func (self *Cell) IsPure(hasher Hasher) bool {
	if self.p.Count == 1 {
		hash := hasher.Hash(self.p.Id)
		if self.p.Hash == hash {
			return true
		}
	}
	return false
}

func (self *Cell) MarshalJSON() ([]byte, error) {
	return json.Marshal(&self.p)
}

func (self *Cell) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &self.p)
}
