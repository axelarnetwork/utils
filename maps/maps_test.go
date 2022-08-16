package maps_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/maps"
)

func TestHas(t *testing.T) {
	m := map[int]string{3: "test"}

	assert.True(t, maps.Has(m, 3))
	assert.False(t, maps.Has(m, 10))
}
