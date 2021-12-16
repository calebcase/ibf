package ibf

import (
	"encoding/binary"

	xor "github.com/go-faster/xor"
)

// block is a byte array where the first 8 bytes are the big endian uint64
// length of the value. When blocks are combined via xor they are extended to
// the right with zero bytes. When a block is "pure" then it can be losslessly
// recovered by truncating to the stored length. Blocks rely on external logic
// to know when it is "pure" (e.g. see Cell) and attempts to get the value of
// an "unpure" block will likely panic (and if it doesn't will return
// nonsense).
type block struct {
	Data []byte `json:"data"`
}

func newBlock(value []byte) (b *block) {
	data := make([]byte, 8+len(value))
	binary.BigEndian.PutUint64(data, uint64(len(value)))
	copy(data[8:], value)

	return &block{
		Data: data,
	}
}

func (b *block) Xor(other *block) {
	if len(other.Data) > len(b.Data) {
		b.Data = append(b.Data, make([]byte, len(other.Data)-len(b.Data))...)
	}

	// In theory we could extend others.Data harmlessly here, but for now
	// we will avoid modifying other unnecessarily.
	data := make([]byte, len(b.Data))
	copy(data, other.Data)

	xor.Bytes(b.Data, b.Data, data)
}

func (b *block) Value() []byte {
	size := binary.BigEndian.Uint64(b.Data[:8])

	// This truncates the value. The user of value should be comparing it
	// to the value hash (done elsewhere) and that will (usually) catch
	// this error.
	if 8+size > uint64(len(b.Data)) {
		size = uint64(len(b.Data)) - 8
	}

	return b.Data[8 : 8+size]
}

func (b *block) Clone() *block {
	data := make([]byte, len(b.Data))
	copy(data, b.Data)

	return &block{
		Data: data,
	}
}
