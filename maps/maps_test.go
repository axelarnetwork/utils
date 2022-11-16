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

func TestFilter(t *testing.T) {
	m := map[int]string{
		1: "foo",
		2: "bar",
		3: "test"}

	assert.EqualValues(t, map[int]string{2:"bar"}, maps.Filter(m, func(k int, v string) bool { return k !=1 && v != "test" }))
}
