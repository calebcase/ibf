package ibf

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlock(t *testing.T) {
	type TC struct {
		name string

		b *block // accumulated block
		i []byte // input value
		d []byte // block data
		v []byte // value (optional)
	}

	tcs := []TC{
		{
			name: "same length 0",
			b:    newBlock([]byte{}),
			i:    []byte{},
			d:    []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			v:    []byte{},
		},
		{
			name: "same length 1",
			b:    newBlock([]byte{0x0F}),
			i:    []byte{0xF0},
			d:    []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 1 ^ 1, 0xFF},
			v:    nil,
		},
		{
			name: "diff lengths 1,2",
			b:    newBlock([]byte{0x0F}),
			i:    []byte{0x00, 0xF0},
			d:    []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 1 ^ 2, 0x0F, 0xF0},
			v:    nil,
		},
		{
			name: "remove",
			b: &block{
				Data: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 1 ^ 2, 0x0F, 0xF0},
			},
			i: []byte{0x00, 0xF0},
			d: []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 1, 0x0F, 0x00},
			v: []byte{0x0F},
		},
	}

	for i, tc := range tcs {
		t.Run(fmt.Sprintf("[%d] %s", i, tc.name), func(t *testing.T) {
			tc.b.Xor(newBlock(tc.i))

			require.Equal(t, tc.d, tc.b.Data)

			if tc.v != nil {
				require.Equal(t, tc.v, tc.b.Value())
			}
		})
	}
}
