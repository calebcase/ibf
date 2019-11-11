package ibf

import (
	"encoding/json"
	"errors"
	"math/big"
)

// Interface

type IBFer interface {
	Insert(key *big.Int)
	Remove(key *big.Int)

	Pop() (*big.Int, error)

	Union(IBFer)
	Subtract(IBFer)

	Invert()

	Clone() IBFer

	GetSize() uint64
	GetCells() []Celler
	GetCardinality() *big.Int

	IsEmpty() bool

	json.Marshaler
	json.Unmarshaler
}

// IBF

type ibf struct {
	Positioners []*Hash `json:"positioners"`
	Hasher      *Hash   `json:"hasher"`

	Size  uint64  `json:"size"`
	Cells []*Cell `json:"cells"`

	Cardinality *big.Int `json:"cardinality"`
}

type IBF struct {
	p ibf
}

var _ IBFer = (*IBF)(nil)

func NewIBF(size uint64, positioners []*Hash, hasher *Hash) *IBF {
	cells := make([]*Cell, size)
	for i, _ := range cells {
		cells[i] = NewCell()
	}

	return &IBF{ibf{
		Positioners: positioners,
		Hasher:      hasher,

		Size:  size,
		Cells: cells,

		Cardinality: big.NewInt(0),
	}}
}

func NewEmptyIBF() *IBF {
	return &IBF{}
}

func (self *IBF) Insert(key *big.Int) {
	hash := self.p.Hasher.Hash(key)
	used := map[uint64]bool{}

	for _, positioner := range self.p.Positioners {
		index := positioner.Hash(key) % self.p.Size
		for used[index] {
			index = (index + 1) % self.p.Size
		}
		used[index] = true

		self.p.Cells[index].Insert(key, hash)
	}

	self.p.Cardinality.Add(self.p.Cardinality, ONE)
}

func (self *IBF) Remove(key *big.Int) {
	total := len(self.p.Positioners)
	cells := make([]*Cell, total)
	hash := self.p.Hasher.Hash(key)
	used := map[uint64]bool{}

	// Find all the positions.
	for i, positioner := range self.p.Positioners {
		index := positioner.Hash(key) % self.p.Size
		for used[index] {
			index = (index + 1) % self.p.Size
		}
		used[index] = true

		cells[i] = self.p.Cells[index]
	}

	// Determine if all cells are filled.
	all_filled := true
	for _, cell := range cells {
		if cell.IsEmpty() {
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

	self.p.Cardinality.Sub(self.p.Cardinality, ONE)
}

func (self *IBF) Invert() {
	for _, cell := range self.p.Cells {
		cell.Invert()
	}

	self.p.Cardinality.Mul(self.p.Cardinality, NEG)
}

func (self *IBF) Pop() (*big.Int, error) {
	all_empty := true

	// Look for a pure cell.
	for _, cell := range self.p.Cells {
		if cell.IsPure(self.p.Hasher) {
			result := big.NewInt(0).Set(cell.GetId())
			self.Remove(result)
			return result, nil
		}

		if all_empty && !cell.IsEmpty() {
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

func (self *IBF) Union(ibf IBFer) {
	cells := ibf.GetCells()

	for i := 0; i < len(self.p.Cells); i++ {
		self.p.Cells[i].Union(cells[i])
	}

	self.p.Cardinality.Add(self.p.Cardinality, ibf.GetCardinality())
}

func (self *IBF) Subtract(ibf IBFer) {
	cells := ibf.GetCells()

	for i := 0; i < len(self.p.Cells); i++ {
		self.p.Cells[i].Subtract(cells[i])
	}

	self.p.Cardinality.Sub(self.p.Cardinality, ibf.GetCardinality())
}

func (self *IBF) Clone() IBFer {
	clone := NewIBF(self.p.Size, self.p.Positioners, self.p.Hasher)
	for i, c := range self.p.Cells {
		clone.p.Cells[i] = c.Clone().(*Cell)
	}

	return clone
}

func (self *IBF) GetSize() uint64 {
	return self.p.Size
}

func (self *IBF) GetCells() []Celler {
	cells := make([]Celler, len(self.p.Cells))
	for i, c := range self.p.Cells {
		cells[i] = c
	}
	return cells
}

func (self *IBF) GetCardinality() *big.Int {
	return big.NewInt(0).Set(self.p.Cardinality)
}

func (self *IBF) IsEmpty() bool {
	all_empty := true

	for _, cell := range self.p.Cells {
		if !cell.IsEmpty() && cell.GetCount() > 0 {
			all_empty = false
			break
		}
	}

	return all_empty
}

func (self *IBF) MarshalJSON() ([]byte, error) {
	return json.Marshal(&self.p)
}

func (self *IBF) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &self.p)
}
