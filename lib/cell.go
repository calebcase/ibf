package ibf

// Cell contains the state of an individual position in an IBF.
type Cell struct {
	Key    *block `json:"key"`
	Digest uint64 `json:"digest"`
	Count  int64  `json:"count"`
}

// NewCell returns a new empty cell.
func NewCell() *Cell {
	return &Cell{
		Key:    newBlock([]byte{}),
		Digest: 0,
		Count:  0,
	}
}

// Insert adds the key with the given digest to this cell.
//
// NOTE: This assumes the key does not already exist in the cell. If it does
// this effectively removes it and the count will be incorrect.
func (c *Cell) Insert(key []byte, digest uint64) {
	c.Key.Xor(newBlock(key))
	c.Digest = c.Digest ^ digest
	c.Count++
}

// Remove deletes the key with the given digest from this cell.
//
// NOTE: This assumes the key already exists in the cell. If it does not this
// effectively adds it and the count will be incorrect.
func (c *Cell) Remove(key []byte, digest uint64) {
	c.Key.Xor(newBlock(key))
	c.Digest = c.Digest ^ digest
	c.Count--
}

// Union adds all keys from the given cell to this one.
//
// NOTE: This assumes cells were disjoint sets. If there weren't this effectly
// performs a symmetric difference and the count will be incorrect.
func (c *Cell) Union(cell *Cell) {
	c.Key.Xor(cell.Key)
	c.Digest = c.Digest ^ cell.GetDigest()
	c.Count += cell.GetCount()
}

// Subtract removes all keys in the given cell from this one.
//
// NOTE: This assumes given cell is a subset of this one. If it wasn't this
// effectly performs a symmetric difference and the count will be incorrect.
func (c *Cell) Subtract(cell *Cell) {
	c.Key.Xor(cell.Key)
	c.Digest = c.Digest ^ cell.GetDigest()
	c.Count -= cell.GetCount()
}

// Invert negates the count.
func (c *Cell) Invert() {
	c.Count *= -1
}

// Clone returns a deep copy of this cell.
func (c *Cell) Clone() *Cell {
	return &Cell{
		Key:    c.Key.Clone(),
		Digest: c.Digest,
		Count:  c.Count,
	}
}

// GetKey returns a copy of the key.
func (c *Cell) GetKey() []byte {
	raw := c.Key.Value()

	key := make([]byte, len(raw))
	copy(key, raw)

	return key
}

// GetDigest returns the cell's digest.
func (c *Cell) GetDigest() uint64 {
	return c.Digest
}

// GetCount returns the cell's count.
func (c *Cell) GetCount() int64 {
	return c.Count
}

// IsEmpty returns true if the cell's count is zero, key is empty, and digest
// is zero.
func (c *Cell) IsEmpty() bool {
	return c.Count == 0 && len(c.Key.Value()) == 0 && c.Digest == 0
}

// IsPure returns true if the cell contains exactly one value and the hash is
// valid.
func (c *Cell) IsPure(h *Hash) bool {
	if c.Count == 1 {
		return c.Digest == h.Hash(c.Key.Value())
	}

	return false
}
