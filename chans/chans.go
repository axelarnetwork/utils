package chans

import (
	"context"
	"golang.org/x/exp/constraints"
	"math/big"
)

var oneBig = big.NewInt(1)

// Concat produces a channel containing all items from all channels concatenated in the order they are passed.
func Concat[T any](channels ...<-chan T) <-chan T {
	newCh := make(chan T)

	go func() {
		defer close(newCh)

		for _, ch := range channels {
			for v := range ch {
				newCh <- v
			}
		}
	}()

	return newCh
}

// Empty returns an empty chanel
func Empty[T any]() <-chan T {
	newCh := make(chan T)
	close(newCh)

	return newCh
}

// Filter returns a new channel that only contains elements that match the predicate. Runs until source channel is closed
func Filter[T any](source <-chan T, predicate func(T) bool) <-chan T {
	out := make(chan T, cap(source))

	go func() {
		defer close(out)
		for x := range source {
			if predicate(x) {
				out <- x
			}
		}
	}()

	return out
}

// Flatten flattens a chan of chans into a chan of elements
func Flatten[T any](source <-chan <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		for c := range source {
			for element := range c {
				out <- element
			}
		}
	}()

	return out
}

// ForEach performs the given function on every element in the channel.  Runs until source channel is closed
func ForEach[T any](source <-chan T, f func(T)) {
	go func() {
		for x := range source {
			f(x)
		}
	}()
}

// FromValues creates a ComposableChannel from a list of values
func FromValues[T any](values ...T) <-chan T {
	newCh := make(chan T, len(values))
	defer close(newCh)

	for _, v := range values {
		newCh <- v
	}

	return newCh
}

// Map maps a channel of T to a channel of S. Runs until source channel is closed.
// By default, output channel has the same capacity as the source channel.
// Desired capacity of the output channel can be specified via an optional argument.
func Map[T, S any](source <-chan T, f func(T) S, capacity ...int) <-chan S {
	var out chan S
	if len(capacity) > 0 {
		out = make(chan S, capacity[0])
	} else {
		out = make(chan S)
	}

	go func() {
		defer close(out)
		for x := range source {
			out <- f(x)
		}
	}()

	return out
}

// Pop gets the next element from the given channel. Returns false if the context expires before an element is available.
func Pop[T any](ctx context.Context, c <-chan T) (T, bool) {
	if ctx.Err() != nil {
		return *new(T), false
	}

	select {
	case <-ctx.Done():
		return *new(T), false
	case x := <-c:
		return x, true
	}
}

// Push adds the given element to the given channel. Returns false if the context expires before the channel unblocks.
func Push[T any](ctx context.Context, c chan<- T, x T) bool {
	if ctx.Err() != nil {
		return false
	}
	select {
	case <-ctx.Done():
		return false
	case c <- x:
		return true
	}
}

// Range creates a channel that inclusively contains all values from `from` to `to`.
func Range[T constraints.Integer](from, to T) <-chan T {
	if from > to {
		return Empty[T]()
	}

	newCh := make(chan T)

	go func() {
		defer close(newCh)

		for i := from; i <= to; i++ {
			newCh <- i
		}
	}()

	return newCh
}

// RangeBig creates a channel that inclusively contains all values from `from` to `to`.
func RangeBig(from, to *big.Int) <-chan *big.Int {
	if from.Cmp(to) > 0 {
		return Empty[*big.Int]()
	}

	newCh := make(chan *big.Int)

	go func() {
		defer close(newCh)

		for i := (&big.Int{}).Set(from); i.Cmp(to) <= 0; i.Add(i, oneBig) {
			newCh <- (&big.Int{}).Set(i)
		}
	}()

	return newCh
}

// Drain enumerates all items from the channel and discards them.
// Returns number of items drained.
func Drain[T any](channel <-chan T) int {
	drained := 0
	for range channel {
		drained++
	}
	return drained
}
