package chans

import (
	"sync"
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

func TestFilter(t *testing.T) {
	source := make(chan int, 100)

	for i := 0; i < 100; i++ {
		source <- i
	}
	close(source)

	even := Filter(source, func(i int) bool { return i%2 == 0 })

	for x := range even {
		assert.Equal(t, 0, x%2)
	}
}

func TestForEach(t *testing.T) {
	source := make(chan int, 100)
	total := 0

	wg := &sync.WaitGroup{}
	wg.Add(100)
	ForEach(source, func(n int) {
		total += n
		wg.Done()
	})

	for i := 0; i < 100; i++ {
		source <- 1
	}
	close(source)
	wg.Wait()

	assert.Equal(t, 100, total)
}

func TestFlatten(t *testing.T) {
	source := make(chan []int, 1)

	out := Flatten(source)

	total := 0
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for x := range out {
			total += x
		}
	}()

	for i := 0; i < 99; i += 3 {
		source <- []int{i, i + 1, i + 2}
	}
	close(source)
	wg.Wait()
	assert.Equal(t, 4851, total)
}
