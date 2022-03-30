package slices

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/test/rand"
)

func TestMap(t *testing.T) {
	source := make([]string, 0, 20)

	for i := 0; i < 20; i++ {
		source = append(source, rand.StrBetween(5, 20))
	}

	out := Map(source, func(s string) int { return len(s) })

	for i := range out {
		assert.Equal(t, len(source[i]), out[i])
	}
}
