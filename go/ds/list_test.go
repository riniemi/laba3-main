package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// listAPI - общие операции односвязного и двусвязного списка: один набор тестов на оба класса.
type listAPI interface {
	Storable
	PushHead(v string)
	PushTail(v string)
	PushAfter(target, v string) error
	PushBefore(target, v string) error
	PopHead() (string, error)
	PopTail() (string, error)
	DelAfter(target string) (string, error)
	DelBefore(target string) (string, error)
	DelValue(v string) bool
	Find(v string) int
	At(i int) (string, error)
	Items() []string
}

var listKinds = map[string]func() listAPI{
	"singly": func() listAPI { return NewSinglyList() },
	"doubly": func() listAPI { return NewDoublyList() },
}

func eachList(t *testing.T, body func(t *testing.T, l listAPI)) {
	for name, mk := range listKinds {
		mk := mk
		t.Run(name, func(t *testing.T) { body(t, mk()) })
	}
}

func TestListPushOrder(t *testing.T) {
	eachList(t, func(t *testing.T, l listAPI) {
		l.PushTail("a")
		l.PushTail("b")
		l.PushHead("z")
		assert.Equal(t, []string{"z", "a", "b"}, l.Items())
		assert.Equal(t, 3, l.Size())
	})
}

func TestListInsertAroundTarget(t *testing.T) {
	eachList(t, func(t *testing.T, l listAPI) {
		require.NoError(t, l.Load([]string{"a", "b", "c"}))
		require.NoError(t, l.PushAfter("c", "d")) // после хвоста: хвост должен обновиться
		l.PushTail("e")
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, l.Items())
		require.NoError(t, l.PushBefore("a", "s")) // перед головой
		require.NoError(t, l.PushBefore("c", "x")) // в середину
		assert.Equal(t, []string{"s", "a", "b", "x", "c", "d", "e"}, l.Items())

		assert.ErrorIs(t, l.PushAfter("zz", "q"), ErrNotFound)
		assert.ErrorIs(t, l.PushBefore("zz", "q"), ErrNotFound)
		assert.Equal(t, 7, l.Size())
	})
}

func TestListPops(t *testing.T) {
	eachList(t, func(t *testing.T, l listAPI) {
		_, err := l.PopHead()
		assert.ErrorIs(t, err, ErrEmpty)
		_, err = l.PopTail()
		assert.ErrorIs(t, err, ErrEmpty)

		require.NoError(t, l.Load([]string{"1", "2", "3"}))
		v, err := l.PopTail()
		require.NoError(t, err)
		assert.Equal(t, "3", v)
		v, err = l.PopHead()
		require.NoError(t, err)
		assert.Equal(t, "1", v)
		v, err = l.PopTail() // единственный элемент
		require.NoError(t, err)
		assert.Equal(t, "2", v)
		assert.Equal(t, 0, l.Size())

		l.PushTail("x") // после опустошения голова и хвост должны быть согласованы
		assert.Equal(t, []string{"x"}, l.Items())
	})
}

func TestListDelAfterBefore(t *testing.T) {
	eachList(t, func(t *testing.T, l listAPI) {
		require.NoError(t, l.Load([]string{"a", "b", "c", "d"}))
		v, err := l.DelAfter("a")
		require.NoError(t, err)
		assert.Equal(t, "b", v)
		v, err = l.DelAfter("c") // удаляем хвост
		require.NoError(t, err)
		assert.Equal(t, "d", v)
		l.PushTail("e") // хвост обновился: [a c e]
		assert.Equal(t, []string{"a", "c", "e"}, l.Items())

		v, err = l.DelBefore("e")
		require.NoError(t, err)
		assert.Equal(t, "c", v)
		assert.Equal(t, []string{"a", "e"}, l.Items())

		v, err = l.DelBefore("e") // перед e стоит голова
		require.NoError(t, err)
		assert.Equal(t, "a", v)
		assert.Equal(t, []string{"e"}, l.Items())

		_, err = l.DelBefore("e")
		assert.ErrorIs(t, err, ErrNoNeighbor)
		_, err = l.DelAfter("e")
		assert.ErrorIs(t, err, ErrNoNeighbor)
		_, err = l.DelBefore("zz")
		assert.ErrorIs(t, err, ErrNotFound)
		_, err = l.DelAfter("zz")
		assert.ErrorIs(t, err, ErrNotFound)
	})
}

func TestListDelValue(t *testing.T) {
	eachList(t, func(t *testing.T, l listAPI) {
		require.NoError(t, l.Load([]string{"x", "y", "z"}))
		assert.True(t, l.DelValue("z")) // хвост
		l.PushTail("w")
		assert.Equal(t, []string{"x", "y", "w"}, l.Items())
		assert.False(t, l.DelValue("nope"))
		assert.True(t, l.DelValue("y")) // середина
		assert.True(t, l.DelValue("x")) // голова
		assert.Equal(t, []string{"w"}, l.Items())
		assert.Equal(t, 1, l.Size())
	})
}

func TestListFindAt(t *testing.T) {
	eachList(t, func(t *testing.T, l listAPI) {
		require.NoError(t, l.Load([]string{"a", "b", "c"}))
		assert.Equal(t, 2, l.Find("c"))
		assert.Equal(t, -1, l.Find("zz"))
		v, err := l.At(1)
		require.NoError(t, err)
		assert.Equal(t, "b", v)
		_, err = l.At(3)
		assert.ErrorIs(t, err, ErrIndex)
		_, err = l.At(-1)
		assert.ErrorIs(t, err, ErrIndex)
	})
}

func TestListDumpLoadType(t *testing.T) {
	s := NewSinglyList()
	d := NewDoublyList()
	assert.Equal(t, byte('F'), s.Type())
	assert.Equal(t, byte('L'), d.Type())
	eachList(t, func(t *testing.T, l listAPI) {
		assert.Equal(t, []string{}, l.Dump())
		require.NoError(t, l.Load([]string{"1", "2", "3"}))
		assert.Equal(t, []string{"1", "2", "3"}, l.Dump())
	})
}

func TestSinglyListShow(t *testing.T) {
	l := NewSinglyList()
	assert.Equal(t, "NULL", l.Show())
	require.NoError(t, l.Load([]string{"a", "b"}))
	assert.Equal(t, "a -> b -> NULL", l.Show())
}

func TestDoublyListShowAndReversed(t *testing.T) {
	l := NewDoublyList()
	assert.Equal(t, "NULL", l.Show())
	require.NoError(t, l.Load([]string{"a", "b", "c"}))
	assert.Equal(t, "NULL <-> a <-> b <-> c <-> NULL", l.Show())
	assert.Equal(t, []string{"c", "b", "a"}, l.ItemsReversed())

	// после смешанных операций прямой и обратный обход должны совпадать
	require.NoError(t, l.PushBefore("b", "x"))
	require.NoError(t, l.PushAfter("c", "y"))
	_, err := l.DelAfter("a")
	require.NoError(t, err)
	l.PushHead("h")
	fwd := l.Items()
	back := l.ItemsReversed()
	require.Len(t, back, len(fwd))
	for i := range fwd {
		assert.Equal(t, fwd[i], back[len(fwd)-1-i])
	}
}
