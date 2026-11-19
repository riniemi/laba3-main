package ds

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func put(t *testing.T, h *HashTable, k, v string) bool {
	t.Helper()
	isNew, err := h.Put(k, v)
	require.NoError(t, err)
	return isNew
}

func TestHashTablePutGetEraseForBothMethods(t *testing.T) {
	for _, m := range []Method{Linear, Double} {
		h := NewHashTable(m)
		assert.True(t, put(t, h, "apple", "10"))
		assert.True(t, put(t, h, "banana", "20"))
		assert.False(t, put(t, h, "apple", "15"), "обновление, а не новый ключ")

		v, ok := h.Get("apple")
		assert.True(t, ok)
		assert.Equal(t, "15", v)
		assert.True(t, h.Contains("banana"))
		assert.False(t, h.Contains("cherry"))
		assert.Equal(t, 2, h.Size())
		assert.Equal(t, m, h.MethodOf())
		assert.Equal(t, byte('H'), h.Type())

		assert.True(t, h.Erase("apple"))
		assert.False(t, h.Erase("apple"))
		assert.False(t, h.Contains("apple"))
		assert.Equal(t, 1, h.Size())
	}
}

func TestHashTableGrowsAndKeepsEverything(t *testing.T) {
	for _, m := range []Method{Linear, Double} {
		h := NewHashTable(m)
		for i := 0; i < 200; i++ {
			put(t, h, "k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
		}
		assert.Equal(t, 200, h.Size())
		assert.Greater(t, h.Capacity(), 200)
		for i := 0; i < 200; i += 2 {
			assert.True(t, h.Erase("k"+strconv.Itoa(i)))
		}
		for i := 0; i < 200; i++ {
			v, ok := h.Get("k" + strconv.Itoa(i))
			if i%2 == 0 {
				assert.False(t, ok, "k%d удален", i)
			} else {
				assert.True(t, ok)
				assert.Equal(t, "v"+strconv.Itoa(i), v)
			}
		}
	}
}

func TestHashTableCapacityGrowthMatchesCpp(t *testing.T) {
	h := NewHashTable(Double)
	for i := 0; i < 30; i++ {
		put(t, h, "k"+strconv.Itoa(i), "v")
	}
	assert.Equal(t, 47, h.Capacity())
	assert.Equal(t, 30, h.Size())
}

func TestHashTableTombstonesAreCleanedWithoutGrowth(t *testing.T) {
	h := NewHashTable(Linear)
	for i := 0; i < 50; i++ { // живых почти нет, надгробия копятся и должны вычищаться
		k := "k" + strconv.Itoa(i)
		put(t, h, k, "v")
		assert.True(t, h.Erase(k))
	}
	assert.Equal(t, 11, h.Capacity())
	assert.Equal(t, 0, h.Size())
}

func TestHashTableReusesTombstone(t *testing.T) {
	h := NewHashTable(Linear)
	put(t, h, "a", "1")
	assert.True(t, h.Erase("a"))
	assert.True(t, put(t, h, "a", "2"))
	v, ok := h.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "2", v)
	assert.Equal(t, 1, h.Size())
}

func TestHashTableKeyWithSeparatorRejected(t *testing.T) {
	h := NewHashTable(Double)
	_, err := h.Put("bad\x1Fkey", "v")
	assert.ErrorIs(t, err, ErrBadKey)
	assert.Equal(t, 0, h.Size())
}

func TestHashTableFullWithoutAutoGrow(t *testing.T) {
	h := NewHashTableWith(Linear, 3, false)
	assert.Equal(t, 3, h.Capacity())
	for i := 0; i < 3; i++ {
		put(t, h, "k"+strconv.Itoa(i), "v")
	}
	_, err := h.Put("extra", "v")
	assert.ErrorIs(t, err, ErrFull)
	assert.False(t, put(t, h, "k0", "new"), "существующий ключ обновляется даже в полной таблице")
}

func TestHashTableMinimumCapacityIsPrime(t *testing.T) {
	assert.Equal(t, 3, NewHashTableWith(Double, 0, true).Capacity())
	assert.Equal(t, 11, NewHashTableWith(Double, 8, true).Capacity())
}

func TestHashTableDumpMatchesCpp(t *testing.T) {
	for _, m := range []Method{Linear, Double} {
		name := map[Method]string{Linear: "linear", Double: "double"}[m]
		h := NewHashTable(m)
		put(t, h, "apple", "10")
		put(t, h, "banana", "20")
		put(t, h, "cherry", "hello world")
		assert.Equal(t, []string{name, "11", "3", "0", "0", "cherry", "hello world", "2", "apple", "10", "9", "banana", "20"}, h.Dump())

		assert.True(t, h.Erase("banana"))
		assert.Equal(t, []string{name, "11", "2", "1", "0", "cherry", "hello world", "2", "apple", "10", "9"}, h.Dump())
	}
}

func TestHashTableShowMatchesCpp(t *testing.T) {
	h := NewHashTable(Double)
	put(t, h, "apple", "10")
	put(t, h, "banana", "20")
	put(t, h, "cherry", "hello world")
	want := "method: double hashing, capacity: 11, size: 3, deleted: 0, load: 27%\n" +
		"[0] cherry => \"hello world\"\n[2] apple => 10\n[9] banana => 20"
	assert.Equal(t, want, h.Show())

	h.Erase("banana")
	assert.Contains(t, h.Show(), "[9] (deleted)")
}

func TestHashTableStatsOfEmptyTable(t *testing.T) {
	want := "method: linear probing, capacity: 11, size: 0, deleted: 0, load: 0.00%"
	assert.Equal(t, want, NewHashTable(Linear).Stats())
}

func TestHashTableDumpLoadRoundTripWithTombstones(t *testing.T) {
	src := NewHashTable(Linear)
	for i := 0; i < 8; i++ {
		put(t, src, "k"+strconv.Itoa(i), "v"+strconv.Itoa(i))
	}
	src.Erase("k3")
	src.Erase("k5")

	dst := NewHashTable(Double) // метод и размер берутся из данных
	require.NoError(t, dst.Load(src.Dump()))
	assert.Equal(t, src.Dump(), dst.Dump())
	assert.Equal(t, src.Show(), dst.Show())
	assert.Equal(t, Linear, dst.MethodOf())
	for _, k := range []string{"k0", "k1", "k2", "k4", "k6", "k7"} {
		assert.True(t, dst.Contains(k), k)
	}
	assert.False(t, dst.Contains("k3"))
}

func TestHashTableLoadRejectsCorruptedData(t *testing.T) {
	bad := [][]string{
		{"x"},                                       // слишком мало элементов
		{"triple", "11", "0", "0"},                  // неизвестный метод
		{"double", "12", "0", "0"},                  // размер не простое число
		{"double", "abc", "0", "0"},                 // не число
		{"double", "11", "-1", "0"},                 // отрицательное число элементов
		{"double", "11", "1", "0"},                  // заявлен 1 элемент, данных нет
		{"double", "11", "1", "0", "99", "k", "v"},  // индекс живой ячейки вне таблицы
		{"double", "11", "0", "1", "99"},            // индекс надгробия вне таблицы
		{"double", "11", "1", "0", "x", "k", "v"},   // индекс не число
		{"double", "11", "0", "1", "x"},             // индекс надгробия не число
		{"double", "4194319", "0", "0"},             // размер больше допустимого
		{"double", "11", "99999999", "0"},           // нелепое число элементов
	}
	for _, tokens := range bad {
		h := NewHashTable(Double)
		assert.ErrorIs(t, h.Load(tokens), ErrCorrupted, "%v", tokens)
	}
}
