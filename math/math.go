package math

import "golang.org/x/exp/constraints"

// Max returns the bigger of two numbers
func Max[T constraints.Ordered](first, second T) T {
	if first > second {
		return first
	}
	return second
}

// Min returns the smaller of two numbers
func Min[T constraints.Ordered](first, second T) T {
	if first > second {
		return second
	}
	return first
}
