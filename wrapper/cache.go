package wrapper

// Cached wraps a lazy value getter and stores the result for quick access
type Cached[T any] struct {
	value T
	isSet bool
	get   func() T
}

// NewCached returns a new cached value
func NewCached[T any](getValue func() T) Cached[T] {
	return Cached[T]{
		get: getValue,
	}
}

func (c *Cached[T]) Value() T {
	if !c.isSet {
		c.value = c.get()
		c.isSet = true
	}

	return c.value
}
