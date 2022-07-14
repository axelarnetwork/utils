package funcs

import (
	"fmt"

	"github.com/pkg/errors"
)

// Compose composes two compatible functions
func Compose[T1, T2, T3 any](f func(T1) T2, g func(T2) T3) func(T1) T3 {
	return func(x T1) T3 {
		return g(f(x))
	}
}

// Identity returns the given element
func Identity[T any](x T) T { return x }

// Must returns the result if err is nil, panics otherwise
func Must[T any](result T, err error) T {
	if err != nil {
		panic(errors.Wrap(err, "call should not have failed"))
	}

	return result
}

// MustNoErr panics if err is not nil
func MustNoErr(err error) {
	if err != nil {
		panic(errors.Wrap(err, "call should not have failed"))
	}
}

// MustOk returns the result if ok, panics otherwise
func MustOk[T any](result T, ok bool) T {
	if !ok {
		panic(fmt.Errorf("result is not found"))
	}

	return result
}

// Not returns a new predicate function that would return true if the given
// predicate function returns false; false otherwise
func Not[T any](predicateFunc func(T) bool) func(T) bool {
	return func(t T) bool {
		return !predicateFunc(t)
	}
}

// And returns a new predicate function that would return true if all of the
// given predicate functions return true for the given value; false otherwise
func And[T any](predicateFuncs ...func(T) bool) func(T) bool {
	return func(t T) bool {
		for _, prepredicateFunc := range predicateFuncs {
			if !prepredicateFunc(t) {
				return false
			}
		}

		return true
	}
}

// Or returns a new predicate function that would return true if any of the
// given predicate functions returns true for the given value; false otherwise
func Or[T any](predicateFuncs ...func(T) bool) func(T) bool {
	return func(t T) bool {
		for _, prepredicateFunc := range predicateFuncs {
			if prepredicateFunc(t) {
				return true
			}
		}

		return false
	}
}
