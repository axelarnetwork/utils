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

func TestAll(t *testing.T) {
	even := make([]int, 0, 20)

	for i := 0; i < 20; i += 2 {
		even = append(even, i)
	}

	assert.True(t, All(even, func(i int) bool { return i%2 == 0 }))
	assert.False(t, All(append(even, 5), func(i int) bool { return i%2 == 0 }))
}

func TestAny(t *testing.T) {
	even := make([]int, 0, 20)

	for i := 0; i < 20; i += 2 {
		even = append(even, i)
	}

	assert.False(t, Any(even, func(i int) bool { return i%2 != 0 }))
	assert.True(t, Any(append(even, 5), func(i int) bool { return i%2 != 0 }))
}

func TestFilter(t *testing.T) {
	source := make([]int, 0, 100)

	for i := 0; i < 100; i++ {
		source = append(source, i)
	}

	even := Filter(source, func(i int) bool { return i%2 == 0 })

	for _, x := range even {
		assert.Equal(t, 0, x%2)
	}
}

func TestForEach(t *testing.T) {
	source := make([]int, 0, 100)
	total := 0

	for i := 0; i < 100; i++ {
		source = append(source, 1)
	}

	ForEach(source, func(n int) {
		total += n
	})

	assert.Equal(t, 100, total)
}

func TestReduce(t *testing.T) {
	source := make([]string, 0, 100)
	for i := 0; i < 100; i++ {
		source = append(source, rand.Str(i))
	}

	assert.Equal(t, 4950, Reduce(source, 0, func(v int, i string) int { return v + len(i) }))
}
