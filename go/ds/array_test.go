package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArrayPushGetSetGrowth(t *testing.T) {
	a := NewArray()
	assert.Equal(t, byte('M'), a.Type())
	for i := 0; i < 10; i++ { // больше начальной емкости 4: проверяем удвоение
		a.Push(string(rune('a' + i)))
	}
	assert.Equal(t, 10, a.Size())
	v, err := a.Get(9)
	require.NoError(t, err)
	assert.Equal(t, "j", v)

	require.NoError(t, a.Set(0, "Z"))
	v, err = a.Get(0)
	require.NoError(t, err)
	assert.Equal(t, "Z", v)
}

func TestArrayInsertAndErase(t *testing.T) {
	a := NewArray()
	require.NoError(t, a.Load([]string{"b", "c"}))
	require.NoError(t, a.Insert(0, "a")) // в начало
	require.NoError(t, a.Insert(3, "d")) // в конец (idx == len)
	require.NoError(t, a.Insert(2, "x")) // в середину
	assert.Equal(t, []string{"a", "b", "x", "c", "d"}, a.Dump())

	v, err := a.Erase(2)
	require.NoError(t, err)
	assert.Equal(t, "x", v)
	v, err = a.Erase(3) // последний
	require.NoError(t, err)
	assert.Equal(t, "d", v)
	assert.Equal(t, []string{"a", "b", "c"}, a.Dump())
}

func TestArrayIndexErrors(t *testing.T) {
	a := NewArray()
	a.Push("x")

	_, err := a.Get(1)
	assert.ErrorIs(t, err, ErrIndex)
	_, err = a.Get(-1)
	assert.ErrorIs(t, err, ErrIndex)
	assert.ErrorIs(t, a.Set(5, "y"), ErrIndex)
	assert.ErrorIs(t, a.Insert(2, "y"), ErrIndex)
	assert.ErrorIs(t, a.Insert(-1, "y"), ErrIndex)
	_, err = a.Erase(1)
	assert.ErrorIs(t, err, ErrIndex)
	assert.Equal(t, 1, a.Size(), "после ошибок массив не изменился")
}

func TestArrayDumpIsCopy(t *testing.T) {
	a := NewArray()
	a.Push("1")
	d := a.Dump()
	d[0] = "changed"
	v, err := a.Get(0)
	require.NoError(t, err)
	assert.Equal(t, "1", v)
}

func TestArrayShow(t *testing.T) {
	a := NewArray()
	assert.Equal(t, "[] (len=0)", a.Show())
	assert.Equal(t, []string{}, a.Dump())
	a.Push("1")
	a.Push("hello world")
	assert.Equal(t, "[1, hello world] (len=2)", a.Show())
}
