package chans

// Map maps a channel of T to a channel of S. Runs until source channel is closed
func Map[T, S any](source <-chan T, f func(T) S) <-chan S {
	ch := make(chan S, cap(source))

	go func() {
		for x := range source {
			ch <- f(x)
		}
	}()

	return ch
}
