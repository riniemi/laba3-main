package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStackLifoOrder(t *testing.T) {
	s := NewStack()
	s.Push("1")
	s.Push("2")
	s.Push("3")
	assert.Equal(t, 3, s.Size())
	for _, want := range []string{"3", "2", "1"} {
		v, err := s.Pop()
		require.NoError(t, err)
		assert.Equal(t, want, v)
	}
	assert.Equal(t, 0, s.Size())
}

func TestStackItemsAndDumpAreMirrored(t *testing.T) {
	s := NewStack()
	s.Push("1")
	s.Push("2")
	assert.Equal(t, []string{"2", "1"}, s.Items()) // от вершины к низу
	assert.Equal(t, []string{"1", "2"}, s.Dump())  // от низа к вершине (так записывается в файл)
	assert.Equal(t, byte('S'), s.Type())
}

func TestStackShow(t *testing.T) {
	s := NewStack()
	assert.Equal(t, "(empty stack)", s.Show())
	s.Push("1")
	s.Push("2")
	assert.Equal(t, "| 2 |  <- top\n| 1 |\n+---+", s.Show())
}

func TestStackLoadRestoresTopElement(t *testing.T) {
	s := NewStack()
	require.NoError(t, s.Load([]string{"1", "2", "3"}))
	v, err := s.Pop()
	require.NoError(t, err)
	assert.Equal(t, "3", v)
}

func TestStackPopFromEmpty(t *testing.T) {
	_, err := NewStack().Pop()
	assert.ErrorIs(t, err, ErrEmpty)
}

func TestQueueFifoOrder(t *testing.T) {
	q := NewQueue()
	q.Push("a")
	q.Push("b")
	q.Push("c")
	assert.Equal(t, 3, q.Size())
	v, err := q.Pop()
	require.NoError(t, err)
	assert.Equal(t, "a", v)
	v, err = q.Pop()
	require.NoError(t, err)
	assert.Equal(t, "b", v)
	assert.Equal(t, []string{"c"}, q.Items())
	assert.Equal(t, byte('Q'), q.Type())
}

func TestQueueShowDumpLoad(t *testing.T) {
	q := NewQueue()
	assert.Equal(t, "(empty queue)", q.Show())
	require.NoError(t, q.Load([]string{"1", "2"}))
	assert.Equal(t, "front -> 1, 2 <- back", q.Show())
	assert.Equal(t, []string{"1", "2"}, q.Dump())
}

func TestQueuePopFromEmpty(t *testing.T) {
	_, err := NewQueue().Pop()
	assert.ErrorIs(t, err, ErrEmpty)
}
