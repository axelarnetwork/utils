package math

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	testutils "github.com/axelarnetwork/utils/test"
)

func TestMax(t *testing.T) {
	assert.Equal(t, 10, Max(5, 10))
	assert.EqualValues(t, 5, Max(5.0, -5.0))

	t.Run("checking max values", testutils.Func(func(t *testing.T) {
		first := rand.Intn(10000)
		second := rand.Intn(10000)
		assert.Equal(t, Max(first, second), Max(first, second))
		assert.Equal(t, first+second, Max(first+second, first-second))
		assert.Equal(t, -first+second, Max(-first-second, -first+second))
	}).Repeat(20))
}

func TestMin(t *testing.T) {
	assert.Equal(t, 5, Min(5, 10))
	assert.EqualValues(t, -5, Min(5.0, -5.0))

	t.Run("checking min values", testutils.Func(func(t *testing.T) {
		first := rand.Intn(10000)
		second := rand.Intn(10000)
		assert.Equal(t, Min(first, second), Min(first, second))
		assert.Equal(t, first-second, Min(first+second, first-second))
		assert.Equal(t, -first-second, Min(-first-second, -first+second))
	}).Repeat(20))
}
