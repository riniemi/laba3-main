package ds

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactoryCreatesEveryType(t *testing.T) {
	for _, ty := range []byte("MFLSQTH") {
		s, err := Create(ty)
		require.NoError(t, err)
		assert.Equal(t, ty, s.Type())
		assert.Equal(t, 0, s.Size())
	}
}

func TestFactoryRejectsUnknownType(t *testing.T) {
	s, err := Create('X')
	assert.Nil(t, s)
	assert.ErrorIs(t, err, ErrUnknown)
}

func TestPolymorphicShowThroughInterface(t *testing.T) {
	// один и тот же вызов Show() выполняет код реального класса
	want := map[byte]string{
		'M': "[] (len=0)",
		'F': "NULL",
		'L': "NULL",
		'S': "(empty stack)",
		'Q': "(empty queue)",
		'T': "(empty tree)",
	}
	for ty, show := range want {
		s, err := Create(ty)
		require.NoError(t, err)
		assert.Equal(t, show, s.Show(), "тип %c", ty)
	}
}
