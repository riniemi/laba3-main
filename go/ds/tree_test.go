package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTree(t *testing.T, keys ...string) *CompleteTree {
	t.Helper()
	tr := NewCompleteTree()
	for _, k := range keys {
		_, err := tr.Insert(k)
		require.NoError(t, err)
	}
	return tr
}

func TestTreeInsertNormalizesAndRejectsNonIntegers(t *testing.T) {
	tr := NewCompleteTree()
	v, err := tr.Insert("007")
	require.NoError(t, err)
	assert.Equal(t, "7", v)

	_, err = tr.Insert("abc")
	assert.ErrorIs(t, err, ErrNotInteger)
	assert.Equal(t, 1, tr.Size(), "после ошибки дерево не изменилось")
	assert.Equal(t, byte('T'), tr.Type())
}

func TestTreeContains(t *testing.T) {
	tr := buildTree(t, "10", "-5", "7")
	ok, err := tr.Contains("-5")
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = tr.Contains("0007")
	require.NoError(t, err)
	assert.True(t, ok)
	ok, err = tr.Contains("99")
	require.NoError(t, err)
	assert.False(t, ok)
	_, err = tr.Contains("x")
	assert.ErrorIs(t, err, ErrNotInteger)
}

func TestTreeStaysComplete(t *testing.T) {
	tr := NewCompleteTree()
	assert.True(t, tr.IsComplete(), "пустое дерево считается complete")
	for i := 0; i < 20; i++ {
		_, err := tr.Insert(string(rune('0' + i%10)))
		require.NoError(t, err)
		assert.True(t, tr.IsComplete())
	}
	assert.Equal(t, 20, tr.Size())
}

func TestTreeTraversals(t *testing.T) {
	tr := buildTree(t, "10", "20", "30", "15", "25")
	cases := map[string][]string{
		"PREORDER":  {"10", "20", "15", "25", "30"},
		"INORDER":   {"15", "20", "25", "10", "30"},
		"POSTORDER": {"15", "25", "20", "30", "10"},
		"BFS":       {"10", "20", "30", "15", "25"},
		"SORTED":    {"10", "15", "20", "25", "30"},
	}
	for order, want := range cases {
		got, err := tr.Traverse(order)
		require.NoError(t, err, order)
		assert.Equal(t, want, got, order)
	}
	// SORTED не меняет само дерево
	assert.Equal(t, []string{"10", "20", "15", "#", "#", "25", "#", "#", "30", "#", "#"}, tr.Dump())
}

func TestTreeSortedWithNegatives(t *testing.T) {
	tr := buildTree(t, "3", "-1", "2")
	got, err := tr.Traverse("SORTED")
	require.NoError(t, err)
	assert.Equal(t, []string{"-1", "2", "3"}, got)
}

func TestTreeTraverseUnknownOrderAndEmpty(t *testing.T) {
	tr := NewCompleteTree()
	_, err := tr.Traverse("ZIGZAG")
	assert.ErrorIs(t, err, ErrUnknown)
	got, err := tr.Traverse("BFS")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestTreeDumpLoadRoundTrip(t *testing.T) {
	src := buildTree(t, "10", "20", "30", "15", "25")
	dst := NewCompleteTree()
	require.NoError(t, dst.Load(src.Dump()))
	assert.Equal(t, src.Dump(), dst.Dump())
	assert.Equal(t, src.Show(), dst.Show())
}

func TestTreeLoadIncompleteTree(t *testing.T) {
	tr := NewCompleteTree()
	// корень 1, у него только правый ребенок 2
	require.NoError(t, tr.Load([]string{"1", "#", "2", "#", "#"}))
	assert.Equal(t, 2, tr.Size())
	assert.False(t, tr.IsComplete())
}

func TestTreeLoadErrors(t *testing.T) {
	tr := NewCompleteTree()
	err := tr.Load([]string{"1", "x", "#", "#"})
	assert.ErrorIs(t, err, ErrNotInteger)
	assert.Equal(t, 0, tr.Size(), "при ошибке дерево остается пустым")

	err = tr.Load([]string{"1", "2", "bad"})
	assert.ErrorIs(t, err, ErrNotInteger)

	require.NoError(t, tr.Load(nil))
	assert.Equal(t, 0, tr.Size())
	require.NoError(t, tr.Load([]string{"#"}))
	assert.Equal(t, 0, tr.Size())
}

func TestTreeShowMatchesCppVersion(t *testing.T) {
	assert.Equal(t, "(empty tree)", NewCompleteTree().Show())
	assert.Equal(t, "7", buildTree(t, "7").Show())
	assert.Equal(t, "  1\n┌─┘\n2", buildTree(t, "1", "2").Show())
	assert.Equal(t, "  1\n┌─┴─┐\n2   3", buildTree(t, "1", "2", "3").Show())
	assert.Equal(t,
		"       10\n  ┌────┴────┐\n  20        30\n┌─┴──┐    ┌─┴─┐\n15   25   5   1",
		buildTree(t, "10", "20", "30", "15", "25", "5", "1").Show())
	assert.Equal(t, "    -5\n  ┌─┴──┐\n 100   7\n┌─┘\n8", buildTree(t, "-5", "100", "7", "8").Show())
}

func TestTreeShowRightChildOnly(t *testing.T) {
	tr := NewCompleteTree()
	require.NoError(t, tr.Load([]string{"1", "#", "2", "#", "#"}))
	got := tr.Show()
	assert.Contains(t, got, "└")
	assert.Contains(t, got, "2")
}
