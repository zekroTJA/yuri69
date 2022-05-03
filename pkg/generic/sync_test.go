package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadStore(t *testing.T) {
	m := SyncMap[string, int]{}

	m.Store("test", 2)

	v, ok := m.Load("not existent")
	assert.False(t, ok)
	assert.Equal(t, 0, v)

	v, ok = m.Load("test")
	assert.True(t, ok)
	assert.Equal(t, 2, v)
}
