package cached_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/monads/cached"
	"github.com/axelarnetwork/utils/test/rand"
)

func TestCached(t *testing.T) {
	val := cached.New(func() int { return 5 })

	assert.Equal(t, 5, val.Value())

	valFromRand := cached.New(func() string { return rand.Str(10) })

	cachedValue := valFromRand.Value()
	assert.Equal(t, cachedValue, valFromRand.Value())

	valFromRand.Clear()
	assert.NotEqual(t, cachedValue, valFromRand.Value())
}
