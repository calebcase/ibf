package ibf

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
)

func TestIBF(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		vs := []string{
			"a",
			"b",
			"c",
			"d",
		}

		i0 := NewIBF(3, 1)

		for _, v := range vs {
			i0.Insert([]byte(v))
		}

		require.Equal(t, int64(len(vs)), i0.Cardinality)

		i1 := i0.Clone()
		i1.Remove([]byte(vs[0]))
		require.Equal(t, int64(len(vs)), i0.Cardinality)
		require.Equal(t, int64(len(vs)-1), i1.Cardinality)

		i2 := i0.Clone()
		i2.Subtract(i1)
		require.Equal(t, int64(1), i2.Cardinality)

		value, err := i2.Pop()
		require.NoError(t, err)
		require.Equal(t, vs[0], string(value))
	})

	t.Run("leading zeros", func(t *testing.T) {
		vs := [][]byte{
			[]byte{0x00, 0x00, 0x00, 0x01},
			[]byte{0x00, 0x00, 0x01},
			[]byte{0x00, 0x01},
			[]byte{0x01},
		}

		i0 := NewIBF(3, 2)

		for _, v := range vs {
			i0.Insert(v)
		}

		i1 := i0.Clone()
		i1.Remove(vs[0])

		i2 := i0.Clone()
		i2.Subtract(i1)

		value, err := i2.Pop()
		require.NoError(t, err)
		require.Equal(t, vs[0], value)
	})

	t.Run("fuzz", func(t *testing.T) {
		f := fuzz.New().NilChance(0).NumElements(0, 1024)

		vs := make([][]byte, 10)

		for i := 0; i < len(vs); i++ {
			var data []byte
			f.Fuzz(&data)

			vs[i] = data
		}

		t.Log("vs:", spew.Sdump(vs))

		var seed int64
		f.Fuzz(&seed)

		i0 := NewIBF(3, seed)

		for _, v := range vs {
			i0.Insert(v)
		}

		i1 := i0.Clone()
		i1.Remove(vs[0])

		i2 := i0.Clone()
		i2.Subtract(i1)

		value, err := i2.Pop()
		t.Log("value:", spew.Sdump(value))
		require.NoError(t, err)
		require.Equal(t, vs[0], value)
	})
}
