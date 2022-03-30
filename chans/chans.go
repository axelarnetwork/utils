package chans

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
