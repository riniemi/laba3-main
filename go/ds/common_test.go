package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuote(t *testing.T) {
	assert.Equal(t, "abc", Quote("abc"))
	assert.Equal(t, `"a b"`, Quote("a b"))
	assert.Equal(t, `""`, Quote(""))
	assert.Equal(t, `"a\"b"`, Quote(`a"b`))
	assert.Equal(t, `"a\\b"`, Quote(`a\b`))
	assert.Equal(t, "\"tab\there\"", Quote("tab\there"))
}

func TestTokenize(t *testing.T) {
	tok, err := Tokenize(`MPUSH arr "hello world"`)
	require.NoError(t, err)
	assert.Equal(t, []string{"MPUSH", "arr", "hello world"}, tok)

	tok, err = Tokenize(`"a\"b" c "d\\e"`)
	require.NoError(t, err)
	assert.Equal(t, []string{`a"b`, "c", `d\e`}, tok)

	tok, err = Tokenize(`x "" y`)
	require.NoError(t, err)
	assert.Equal(t, []string{"x", "", "y"}, tok)

	tok, err = Tokenize("   \t ")
	require.NoError(t, err)
	assert.Empty(t, tok)
}

func TestTokenizeUnterminatedQuote(t *testing.T) {
	_, err := Tokenize(`a "unterminated`)
	assert.ErrorIs(t, err, ErrCorrupted)
	_, err = Tokenize(`"ends with backslash\`)
	assert.ErrorIs(t, err, ErrCorrupted)
}

func TestQuoteTokenizeRoundTrip(t *testing.T) {
	for _, s := range []string{"plain", "with space", "", `q"uote`, `back\slash`, "multi\nline"} {
		tok, err := Tokenize(Quote(s))
		require.NoError(t, err)
		require.Len(t, tok, 1)
		assert.Equal(t, s, tok[0])
	}
}

func TestIsInteger(t *testing.T) {
	for _, s := range []string{"0", "-5", "123", "999999999999999999", "-999999999999999999", "1234567890123456789"} {
		assert.True(t, IsInteger(s), s)
	}
	for _, s := range []string{"", "-", "1a", "--1", "1.5", "+3", " 1"} {
		assert.False(t, IsInteger(s), s)
	}
}

func TestNormInt(t *testing.T) {
	v, err := NormInt("-007")
	require.NoError(t, err)
	assert.Equal(t, "-7", v)

	v, err = NormInt("-0")
	require.NoError(t, err)
	assert.Equal(t, "0", v)

	_, err = NormInt("abc")
	assert.ErrorIs(t, err, ErrNotInteger)
}
