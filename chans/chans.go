package chans

import "context"

// Map maps a channel of T to a channel of S. Runs until source channel is closed
func Map[T, S any](source <-chan T, f func(T) S) <-chan S {
	out := make(chan S, cap(source))

	go func() {
		defer close(out)
		for x := range source {
			out <- f(x)
		}
	}()

	return out
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

// ForEach performs the given function on every element in the channel.  Runs until source channel is closed
func ForEach[T any](source <-chan T, f func(T)) {
	go func() {
		for x := range source {
			f(x)
		}
	}()
}

// Flatten flattens a chan of chans into a chan of elements
func Flatten[T any](source <-chan <-chan T) <-chan T {
	out := make(chan T, cap(source))

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
