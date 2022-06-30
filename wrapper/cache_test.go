package wrapper_test

import (
	"testing"

	"github.com/axelarnetwork/utils/test/rand"
	"github.com/axelarnetwork/utils/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestCached(t *testing.T) {
	val := wrapper.NewCached(func() int { return 5 })

	assert.Equal(t, 5, val.Value())

	valFromRand := wrapper.NewCached(func() string { return rand.Str(10) })

	assert.Equal(t, valFromRand.Value(), valFromRand.Value())
}
