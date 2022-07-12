package wrapper

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

// NewResult wraps a (value, error) tuple ina Result struct
func NewResult[T any](ok T, err error) Result[T] {
	res := Result[T]{
		ok:  ok,
		err: err,
	}
	if err != nil {
		res.ok = *new(T)
	}
	return res
}

// ContinueWithResult only executes f if res.Error() is nil, returns the original error otherwise
func ContinueWithResult[T1, T2 any](res Result[T1], f func(T1) Result[T2]) Result[T2] {
	if res.Error() != nil {
		return NewResult(*new(T2), res.Error())
	}
	return f(res.Ok())
}
