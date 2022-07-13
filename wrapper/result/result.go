package result

// Result wraps the idiomatic tuple of (value, error)
type Result[T any] struct {
	ok  T
	err error
}

// Ok returns a value of Error is nil
func (res Result[T]) Ok() T {
	return res.ok
}

// Error returns an error, Ok is invalid in that case
func (res Result[T]) Error() error {
	return res.err
}

// New wraps a (value, error) tuple in a Result
func New[T any](ok T, err error) Result[T] {
	if err != nil {
		return FromErr[T](err)
	}
	return FromOk[T](ok)
}

// FromOk returns a Result without error
func FromOk[T any](ok T) Result[T] {
	res := Result[T]{
		ok:  ok,
		err: nil,
	}
	return res
}

// FromErr returns a result with error
func FromErr[T any](err error) Result[T] {
	res := Result[T]{
		ok:  *new(T),
		err: err,
	}
	return res
}

// Pipe only executes f if res.Error() is nil, returns the original error otherwise
func Pipe[T1, T2 any](res Result[T1], f func(T1) Result[T2]) Result[T2] {
	if res.Error() != nil {
		return FromErr[T2](res.Error())
	}
	return f(res.Ok())
}

// Try transforms the value of Result to the new type if OK, returns the original error otherwise
func Try[T1, T2 any](res Result[T1], f func(T1) T2) Result[T2] {
	if res.Error() != nil {
		return FromErr[T2](res.Error())
	}
	return FromOk(f(res.Ok()))
}
