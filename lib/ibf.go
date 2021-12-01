package ibf

import "math/rand"

// IBF holds the state of an invertable bloom filter.
type IBF struct {
	Positioners []*Hash `json:"positioners"`
	Hasher      *Hash   `json:"hasher"`

	Size  uint64  `json:"size"`
	Cells []*Cell `json:"cells"`

	Cardinality int64 `json:"cardinality"`
}

// NewIBF creates a new IBF of the given size. An IBF can accurately handle
// differences of approximately 2/3rds the configured size (e.g. a size of 100
// would allow for ~66 differences to be accurately retrieved). 3 positioners
// and a hasher are created using the output from a random number generator
// initialized with the seed.
func NewIBF(size uint64, seed int64) *IBF {
	rng := rand.New(rand.NewSource(seed))

	positioners := []*Hash{
		NewHash(uint64(rng.Int63()), uint64(rng.Int63())),
		NewHash(uint64(rng.Int63()), uint64(rng.Int63())),
		NewHash(uint64(rng.Int63()), uint64(rng.Int63())),
	}
	hasher := NewHash(uint64(rng.Int63()), uint64(rng.Int63()))

	return NewIBFWithHash(size, positioners, hasher)
}

// NewIBFWithHash creates a new IBF with the provided positioners and hasher.
// It will use the given hashers for positioning and computing the key hashes.
// The positioners must all be initialized with different seeds to ensure they
// do not produce the same positions for the same key.
func NewIBFWithHash(size uint64, positioners []*Hash, hasher *Hash) *IBF {
	cells := make([]*Cell, size)
	for i := range cells {
		cells[i] = NewCell()
	}

	return &IBF{
		Positioners: positioners,
		Hasher:      hasher,

		Size:  size,
		Cells: cells,

		Cardinality: 0,
	}
}

// getPositions returns the cells that the key would occupy. It always returns
// len(positioners) many cells ensuring that no key is under represented.
func (i *IBF) getPositions(key []byte, digest uint64) (cells []*Cell) {
	cells = make([]*Cell, len(i.Positioners))
	used := map[uint64]bool{}

	for j, positioner := range i.Positioners {
		index := positioner.Hash(key) % i.Size

		// NOTE: We need to keep looking if we have found a collision
		// with an already used position.
		for used[index] {
			index = (index + 1) % i.Size
		}

		used[index] = true
		cells[j] = i.Cells[index]
	}

	return cells
}

// Insert adds the key to the set.
//
// NOTE: This does not know if the key already exists and will add it
// unconditionally. If the key did already exist in the set, then that
// effectively would remove it!
func (i *IBF) Insert(key []byte) {
	digest := i.Hasher.Hash(key)
	cells := i.getPositions(key, digest)

	for _, c := range cells {
		c.Insert(key, digest)
	}

	i.Cardinality++
}

// Remove deletes the key from the set.
//
// NOTE: This does not know if the key already exists and will add it
// unconditionally. If the key did already exist in the set, then that
// effectively would add it!
func (i *IBF) Remove(key []byte) {
	digest := i.Hasher.Hash(key)
	cells := i.getPositions(key, digest)

	for _, c := range cells {
		c.Remove(key, digest)
	}

	i.Cardinality--
}

// Invert flips the cardinality of the set and the cells. As if all elements
// has instead been removed from the set instead of added.
func (i *IBF) Invert() {
	for _, cell := range i.Cells {
		cell.Invert()
	}

	i.Cardinality *= -1
}

// Pop finds a key in a pure cell, removes it from the set, and returns it. If
// no pure cell can be found it returns ErrNoPureCell indicating that there are
// more elements in the set, but they cannot be popped. If the set is empty it
// returns ErrEmptySet.
func (i *IBF) Pop() ([]byte, error) {
	allEmpty := true

	// Look for a pure cell.
	for _, cell := range i.Cells {
		if cell.IsPure(i.Hasher) {
			key := cell.GetKey()
			i.Remove(key)

			return key, nil
		}

		if allEmpty && !cell.IsEmpty() {
			allEmpty = false
		}
	}

	// Are there non-empty cells?
	if !allEmpty {
		return nil, ErrNoPureCell
	}

	// Empty set, nothing to pop.
	return nil, ErrEmptySet
}

// Union inserts all the elements from the provided set to this set.
//
// NOTE: This assumes the two sets are disjoint and configured the same. If the
// two sets are not disjoint this will actually perform a symmetric difference
// and the cardinality will be incorrect! If the two sets are not configured
// the same then the behavior is undefined and could potentially panic.
func (i *IBF) Union(other *IBF) {
	cells := other.GetCells()

	for j := 0; j < len(i.Cells); j++ {
		i.Cells[j].Union(cells[j])
	}

	i.Cardinality += other.GetCardinality()
}

// Subtract removes all the elements from the provided set from this set.
//
// NOTE: This assumes the other set is a subset of this one. If that isn't true
// then this will actually perform a symmetric difference and the cardinality
// will be incorrect! If the two sets are not configured the same then the
// behavior is undefined and could potentially panic.
func (i *IBF) Subtract(other *IBF) {
	cells := other.GetCells()

	for j := 0; j < len(i.Cells); j++ {
		i.Cells[j].Subtract(cells[j])
	}

	i.Cardinality -= other.GetCardinality()
}

// Clone returns a copy of this set.
func (i *IBF) Clone() (clone *IBF) {
	clone = NewIBFWithHash(i.Size, i.Positioners, i.Hasher)

	for j, c := range i.Cells {
		clone.Cells[j] = c.Clone()
	}

	clone.Cardinality = i.Cardinality

	return clone
}

// GetSize returns the IBF's size.
func (i *IBF) GetSize() uint64 {
	return i.Size
}

// GetCells returns the IBF's cells.
func (i *IBF) GetCells() []*Cell {
	return i.Cells
}

// GetCardinality returns the IBF's cardinality.
func (i *IBF) GetCardinality() int64 {
	return i.Cardinality
}

// IsEmpty returns true if all the cells are empty and the cardinality is zero.
func (i *IBF) IsEmpty() bool {
	if i.Cardinality != 0 {
		return false
	}

	for _, cell := range i.Cells {
		if !cell.IsEmpty() {
			return false
		}
	}

	return true
}
