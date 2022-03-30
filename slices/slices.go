package slices

func Map[T, S any](source []T, f func(T) S) []S {
	out := make([]S, 0, cap(source))

	// avoid allocating a copy of the slice element
	for i := range source {
		out = append(out, f(source[i]))
	}

	return out
}
