package ibf

//package main

import (
	"errors"
	"math/big"

	"github.com/dchest/siphash"
)

// Interfaces

type Celler interface {
	Insert(key *big.Int, hash uint64)
	Remove(key *big.Int, hash uint64)
	Subtract(Celler)

	Clone() Celler

	GetId() *big.Int
	GetHash() uint64
	GetCount() int64
}

type IBFer interface {
	Insert(key *big.Int)
	Remove(key *big.Int)
	Pop() (*big.Int, error)
	Subtract(IBFer)

	Clone() IBFer

	GetSize() uint64
	GetCells() []Celler
	GetCardinality() uint64

	IsEmpty() bool
}

// Hasher

type Hasher struct {
	Key [2]uint64 `json:"key"`
}

func NewHasher(key0, key1 uint64) *Hasher {
	return &Hasher{[2]uint64{key0, key1}}
}

func (self *Hasher) Hash(key *big.Int) uint64 {
	return siphash.Hash(self.Key[0], self.Key[1], key.Bytes())
}

// Cell

type cell struct {
	Id    *big.Int `json:"id"`
	Hash  uint64   `json:"hash"`
	Count int64    `json:"count"`
}

func NewCell() *cell {
	return &cell{
		Id:    big.NewInt(0),
		Hash:  0,
		Count: 0,
	}
}

func (self *cell) Insert(key *big.Int, hash uint64) {
	self.Id.Xor(self.Id, key)
	self.Hash = self.Hash ^ hash
	self.Count += 1
}

func (self *cell) Remove(key *big.Int, hash uint64) {
	self.Id.Xor(self.Id, key)
	self.Hash = self.Hash ^ hash
	self.Count -= 1
}

func (self *cell) Subtract(cell Celler) {
	self.Id.Xor(self.Id, cell.GetId())
	self.Hash = self.Hash ^ cell.GetHash()
	self.Count -= cell.GetCount()
}

func (self *cell) Clone() Celler {
	return &cell{
		Id:    big.NewInt(0).Set(self.Id),
		Hash:  self.Hash,
		Count: self.Count,
	}
}

func (self *cell) GetId() *big.Int {
	return big.NewInt(0).Set(self.Id)
}

func (self *cell) GetHash() uint64 {
	return self.Hash
}

func (self *cell) GetCount() int64 {
	return self.Count
}

// IBF

type IBF struct {
	Positioners []*Hasher `json:"positioners"`
	Hasher      *Hasher   `json:"hasher"`

	Size  uint64  `json:"size"`
	Cells []*cell `json:"cells"`

	Cardinality uint64 `json:"cardinality"`
}

func NewIBF(size uint64, positioners []*Hasher, hasher *Hasher) *IBF {
	cells := make([]*cell, size)
	for i, _ := range cells {
		cells[i] = NewCell()
	}

	return &IBF{
		Positioners: positioners,
		Hasher:      hasher,

		Size:  size,
		Cells: cells,
	}
}

func (self *IBF) Insert(key *big.Int) {
	hash := self.Hasher.Hash(key)
	used := map[uint64]bool{}

	for _, positioner := range self.Positioners {
		index := positioner.Hash(key) % self.Size
		for used[index] {
			index = (index + 1) % self.Size
		}
		used[index] = true

		self.Cells[index].Insert(key, hash)
	}

	self.Cardinality += 1
}

func (self *IBF) Remove(key *big.Int) {
	total := len(self.Positioners)
	cells := make([]*cell, total)
	hash := self.Hasher.Hash(key)
	used := map[uint64]bool{}

	// Find all the positions.
	for i, positioner := range self.Positioners {
		index := positioner.Hash(key) % self.Size
		for used[index] {
			index = (index + 1) % self.Size
		}
		used[index] = true

		cells[i] = self.Cells[index]
	}

	// Determine if all cells are filled.
	all_filled := true
	for _, cell := range cells {
		if empty(cell) {
			all_filled = false
			break
		}
	}
	if !all_filled {
		// It can't be in the set if all cells aren't filled.
		return
	}

	for _, cell := range cells {
		cell.Remove(key, hash)
	}

	self.Cardinality -= 1
}

var ZERO = big.NewInt(0)

func empty(cell *cell) bool {
	if cell.Count == 0 && cell.Id.Cmp(ZERO) == 0 && cell.Hash == 0 {
		return true
	}
	return false
}

func pure(cell *cell, hasher *Hasher) bool {
	if cell.Count == 1 {
		hash := hasher.Hash(cell.Id)
		if cell.Hash == hash {
			return true
		}
	}
	return false
}

func (self *IBF) pure() (*cell, error) {
	impure := false
	for _, cell := range self.Cells {
		if cell.Count == 1 {
			if self.Hasher.Hash(cell.Id) == cell.Hash {
				return cell, nil
			} // Otherwise this cell is actually irretrievable.
		}

		if !impure {
			if cell.Count < 0 || cell.Count > 1 {
				impure = true
			}

			if cell.Count == 0 {
				if cell.Id.Cmp(big.NewInt(0)) != 0 {
					impure = true
				}
				if cell.Hash != 0 {
					impure = true
				}
			}
		}
	}

	if impure {
		return nil, errors.New("Cannot get another pure value and values remain.")
	} else {
		return nil, errors.New("Nothing left...")
	}
}

func (self *IBF) Pop() (*big.Int, error) {
	all_empty := true

	// Look for a pure cell.
	for _, cell := range self.Cells {
		if pure(cell, self.Hasher) {
			result := big.NewInt(0).Set(cell.Id)
			self.Remove(result)
			return result, nil
		}

		if all_empty && !empty(cell) {
			all_empty = false
		}
	}

	// Are there non-empty cells?
	if !all_empty {
		return nil, errors.New("More elements in the set, but unable to retrieve.")
	}

	// Empty set, nothing to pop.
	return nil, errors.New("Empty set.")
}

func (self *IBF) Subtract(ibf IBFer) {
	cells := ibf.GetCells()

	for i := 0; i < len(self.Cells); i++ {
		self.Cells[i].Subtract(cells[i])
	}

	self.Cardinality -= ibf.GetCardinality()
}

func (self *IBF) Clone() IBFer {
	clone := NewIBF(self.Size, self.Positioners, self.Hasher)
	for i, c := range self.Cells {
		clone.Cells[i] = c.Clone().(*cell)
	}

	return clone
}

func (self *IBF) GetSize() uint64 {
	return self.Size
}

func (self *IBF) GetCells() []Celler {
	cells := make([]Celler, len(self.Cells))
	for i, c := range self.Cells {
		cells[i] = c
	}
	return cells
}

func (self *IBF) GetCardinality() uint64 {
	return self.Cardinality
}

func (self *IBF) IsEmpty() bool {
	all_empty := true

	for _, cell := range self.Cells {
		if !empty(cell) {
			all_empty = false
			break
		}
	}

	return all_empty
}
