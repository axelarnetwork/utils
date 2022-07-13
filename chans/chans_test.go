package chans_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/axelarnetwork/utils/chans"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	source := make(chan int, 1)
	defer close(source)

	even := chans.Map(source, func(i int) bool { return i%2 == 0 })

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

	even := chans.Filter(source, func(i int) bool { return i%2 == 0 })

	for x := range even {
		assert.Equal(t, 0, x%2)
	}
}

func TestForEach(t *testing.T) {
	source := make(chan int, 100)
	total := 0

	wg := &sync.WaitGroup{}
	wg.Add(100)
	chans.ForEach(source, func(n int) {
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
	source := make(chan (<-chan int), 1)

	out := chans.Flatten(source)

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
		c := make(chan int, 3)
		source <- c
		c <- i
		c <- i + 1
		c <- i + 2
	}
	close(source)
	wg.Wait()
	assert.Equal(t, 4851, total)
}

func TestPushPop(t *testing.T) {
	t.Run("valid ctx", func(t *testing.T) {
		c := make(chan int, 1)
		assert.True(t, chans.Push(context.Background(), c, 7))
		i, ok := chans.Pop(context.Background(), c)
		assert.True(t, ok)
		assert.Equal(t, 7, i)
	})

	t.Run("cancelled ctx", func(t *testing.T) {
		c := make(chan int, 1)
		cancelledCtx, cancel := context.WithCancel(context.Background())
		cancel()
		assert.False(t, chans.Push(cancelledCtx, c, 7))
		_, ok := chans.Pop(cancelledCtx, c)
		assert.False(t, ok)
	})

	t.Run("cancel ctx while blocked", func(t *testing.T) {
		c := make(chan int)
		cancelledCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		assert.False(t, chans.Push(cancelledCtx, c, 7))

		cancelledCtx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel2()
		_, ok := chans.Pop(cancelledCtx2, c)
		assert.False(t, ok)
	})
}
