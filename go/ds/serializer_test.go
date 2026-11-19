package ds

import (
	"bytes"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// sample создает по одной заполненной структуре каждого типа.
func sample(t *testing.T) []Storable {
	t.Helper()
	all := make([]Storable, 0)
	for _, ty := range []byte("MFLSQTH") {
		s, err := Create(ty)
		require.NoError(t, err)
		all = append(all, s)
	}
	require.NoError(t, all[0].Load([]string{"10", "hello world", "-5", `q"uote`, ""}))
	require.NoError(t, all[1].Load([]string{"1", "2", "3"}))
	require.NoError(t, all[2].Load([]string{"a", "b", "c"}))
	require.NoError(t, all[3].Load([]string{"1", "2", "3"}))
	require.NoError(t, all[4].Load([]string{"first", "second"}))
	tr, ok := all[5].(*CompleteTree)
	require.True(t, ok)
	for _, v := range []string{"10", "20", "30", "15", "25", "5", "1"} {
		_, err := tr.Insert(v)
		require.NoError(t, err)
	}
	h, ok := all[6].(*HashTable)
	require.True(t, ok)
	_, err := h.Put("apple", "10")
	require.NoError(t, err)
	_, err = h.Put("banana", "hello world")
	require.NoError(t, err)
	h.Erase("apple") // в файле будет надгробие
	return all
}

func TestTextRoundTripForAllStructures(t *testing.T) {
	for _, src := range sample(t) {
		var buf bytes.Buffer
		require.NoError(t, SaveText(src, &buf))
		dst, err := Create(src.Type())
		require.NoError(t, err)
		require.NoError(t, LoadText(dst, &buf))
		assert.Equal(t, src.Dump(), dst.Dump(), "тип %c", src.Type())
	}
}

func TestBinaryRoundTripForAllStructures(t *testing.T) {
	for _, src := range sample(t) {
		var buf bytes.Buffer
		require.NoError(t, SaveBinary(src, &buf))
		dst, err := Create(src.Type())
		require.NoError(t, err)
		require.NoError(t, LoadBinary(dst, &buf))
		assert.Equal(t, src.Dump(), dst.Dump(), "тип %c", src.Type())
	}
}

func TestTextFormatExactBytes(t *testing.T) {
	s := NewStack()
	require.NoError(t, s.Load([]string{"1", "2", "3"}))
	var buf bytes.Buffer
	require.NoError(t, SaveText(s, &buf))
	assert.Equal(t, "S 1 2 3\n", buf.String())

	a := NewArray()
	require.NoError(t, a.Load([]string{"hello world", `a"b`, ""}))
	buf.Reset()
	require.NoError(t, SaveText(a, &buf))
	assert.Equal(t, "M \"hello world\" \"a\\\"b\" \"\"\n", buf.String())

	h := NewHashTable(Double)
	_, err := h.Put("a", "1")
	require.NoError(t, err)
	buf.Reset()
	require.NoError(t, SaveText(h, &buf))
	assert.Equal(t, "H double 11 1 0 9 a 1\n", buf.String())
}

func TestBinaryFormatExactBytes(t *testing.T) {
	s := NewStack()
	require.NoError(t, s.Load([]string{"1", "2", "3"}))
	var buf bytes.Buffer
	require.NoError(t, SaveBinary(s, &buf))
	want := []byte{'D', 'S', 'B', '1', 'S', 3, 0, 0, 0,
		1, 0, 0, 0, '1', 1, 0, 0, 0, '2', 1, 0, 0, 0, '3'}
	assert.Equal(t, want, buf.Bytes())
	assert.Len(t, buf.Bytes(), 24)
}

func TestLoadTextWithoutTrailingNewline(t *testing.T) {
	s := NewStack()
	require.NoError(t, LoadText(s, strings.NewReader("S 1 2")))
	assert.Equal(t, []string{"1", "2"}, s.Dump())
}

func TestLoadTextErrors(t *testing.T) {
	assert.ErrorIs(t, LoadText(NewStack(), strings.NewReader("")), ErrCorrupted)
	assert.ErrorIs(t, LoadText(NewStack(), strings.NewReader("Q 1 2\n")), ErrTypeMatch)
	assert.ErrorIs(t, LoadText(NewStack(), strings.NewReader("\n")), ErrTypeMatch)
	assert.ErrorIs(t, LoadText(NewStack(), strings.NewReader("S \"open\n")), ErrCorrupted)
	assert.ErrorIs(t, LoadText(NewCompleteTree(), strings.NewReader("T x\n")), ErrNotInteger)

	busy := NewStack()
	busy.Push("x")
	assert.ErrorIs(t, LoadText(busy, strings.NewReader("S 1\n")), ErrNotEmpty)
}

// bin собирает бинарные данные для проверки повреждений.
func bin(parts ...[]byte) *bytes.Reader {
	return bytes.NewReader(bytes.Join(parts, nil))
}

func TestLoadBinaryErrors(t *testing.T) {
	magicS := []byte{'D', 'S', 'B', '1', 'S'}
	one := []byte{1, 0, 0, 0}

	assert.ErrorIs(t, LoadBinary(NewStack(), bin()), ErrCorrupted)                               // пусто
	assert.ErrorIs(t, LoadBinary(NewStack(), bin([]byte("XXXXS"))), ErrCorrupted)                // не та сигнатура
	assert.ErrorIs(t, LoadBinary(NewStack(), bin([]byte{'D', 'S', 'B', '1', 'Q'}, one)), ErrTypeMatch)
	assert.ErrorIs(t, LoadBinary(NewStack(), bin(magicS)), ErrCorrupted)                         // нет числа элементов
	assert.ErrorIs(t, LoadBinary(NewStack(), bin(magicS, one)), ErrCorrupted)                    // нет длины элемента
	assert.ErrorIs(t, LoadBinary(NewStack(), bin(magicS, one, one)), ErrCorrupted)               // нет данных элемента
	assert.ErrorIs(t, LoadBinary(NewStack(), bin(magicS, one, []byte{0xFF, 0xFF, 0xFF, 0xFF})), ErrCorrupted) // слишком длинный

	busy := NewStack()
	busy.Push("x")
	assert.ErrorIs(t, LoadBinary(busy, bin(magicS, []byte{0, 0, 0, 0})), ErrNotEmpty)

	ok := NewStack()
	require.NoError(t, LoadBinary(ok, bin(magicS, one, one, []byte{'7'})))
	assert.Equal(t, []string{"7"}, ok.Dump())
}

// failingWriter отказывает на заданной по счету записи.
type failingWriter struct {
	calls, failOn int
}

var errDisk = errors.New("disk full")

func (w *failingWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.calls == w.failOn {
		return 0, errDisk
	}
	return len(p), nil
}

func TestSaveReportsWriteErrors(t *testing.T) {
	s := NewStack()
	require.NoError(t, s.Load([]string{"1"}))
	for call := 1; call <= 5; call++ { // подпись, тип, число элементов, длина, данные
		w := &failingWriter{failOn: call}
		assert.ErrorIs(t, SaveBinary(s, w), errDisk, "отказ на записи №%d", call)
	}
	assert.ErrorIs(t, SaveText(s, &failingWriter{failOn: 1}), errDisk)
}

func TestFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	for _, src := range sample(t) {
		txt := filepath.Join(dir, "d_"+string(rune(src.Type()))+".txt")
		bn := filepath.Join(dir, "d_"+string(rune(src.Type()))+".bin")
		require.NoError(t, SaveTextFile(src, txt))
		require.NoError(t, SaveBinaryFile(src, bn))

		fromText, err := Create(src.Type())
		require.NoError(t, err)
		fromBin, err := Create(src.Type())
		require.NoError(t, err)
		require.NoError(t, LoadTextFile(fromText, txt))
		require.NoError(t, LoadBinaryFile(fromBin, bn))
		assert.Equal(t, src.Dump(), fromText.Dump())
		assert.Equal(t, src.Dump(), fromBin.Dump())
	}
}

func TestFileErrors(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "no_such_dir", "f")
	s := NewStack()
	assert.Error(t, SaveTextFile(s, missing))
	assert.Error(t, SaveBinaryFile(s, missing))
	assert.Error(t, LoadTextFile(NewStack(), missing))
	assert.Error(t, LoadBinaryFile(NewStack(), missing))
}
