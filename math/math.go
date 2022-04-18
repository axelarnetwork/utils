package math

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](first, second T) T {
	if first > second {
		return first
	}
	return second
}

func Min[T constraints.Ordered](first, second T) T {
	if first > second {
		return second
	}
	return first
}
