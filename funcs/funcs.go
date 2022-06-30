package funcs

import "github.com/pkg/errors"

// Compose composes two compatible functions
func Compose[T1, T2, T3 any](f func(T1) T2, g func(T2) T3) func(T1) T3 {
	return func(x T1) T3 {
		return g(f(x))
	}
}

// Identity returns the given element
func Identity[T any](x T) T { return x }

func Must[T any](result T, err error) T {
	if err != nil {
		panic(errors.Wrap(err, "call should not have failed"))
	}

	return result
}
