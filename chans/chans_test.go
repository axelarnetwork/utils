package chans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	source := make(chan int, 1)
	defer close(source)

	even := Map(source, func(i int) bool { return i%2 == 0 })

	for i := 0; i < 100; i++ {
		source <- i
		assert.Equal(t, i%2 == 0, <-even)
	}
}
