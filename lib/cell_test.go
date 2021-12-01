package ibf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCell(t *testing.T) {
	h := NewHash(0, 0)

	a := []byte{0x00, 0x00, 0x01}
	b := []byte{0x01}

	cell := NewCell()
	require.Equal(t, int64(0), cell.Count)
	require.True(t, cell.IsEmpty())

	cell.Insert(a, h.Hash(a))
	cell.Insert(b, h.Hash(b))
	require.Equal(t, int64(2), cell.Count)
	require.False(t, cell.IsEmpty())
	require.False(t, cell.IsPure(h))

	cell.Remove(a, h.Hash(a))
	require.Equal(t, int64(1), cell.Count)
	require.False(t, cell.IsEmpty())
	require.True(t, cell.IsPure(h))
	require.Equal(t, b, cell.GetKey())

	cell.Insert(a, h.Hash(a))
	cell.Remove(b, h.Hash(b))
	require.Equal(t, int64(1), cell.Count)
	require.False(t, cell.IsEmpty())
	require.True(t, cell.IsPure(h))
	require.Equal(t, a, cell.GetKey())
}
