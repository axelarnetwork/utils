package wrapper

type Result[T any] struct {
	ok  T
	err error
}

func (res Result[T]) Ok() T {
	return res.ok
}

func (res Result[T]) Error() error {
	return res.err
}

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
