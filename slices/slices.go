package slices

// Map maps a slice of T to a slice of S
func Map[T, S any](source []T, f func(T) S) []S {
	out := make([]S, len(source))

	// avoid allocating a copy of the slice element
	for i := range source {
		out[i] = f(source[i])
	}

	return out
}

// Reduce performs a reduction to a single value of the source slice according to the given function
func Reduce[T, S any](source []T, initial S, f func(current S, element T) S) S {
	v := initial

	for i := range source {
		v = f(v, source[i])
	}

	return v
}

// Filter returns a new slice that only contains elements that match the predicate
func Filter[T any](source []T, predicate func(T) bool) []T {
	out := make([]T, 0, cap(source))

	for i := range source {
		if predicate(source[i]) {
			out = append(out, source[i])
		}
	}

	return out
}

// ForEach performs the given function on every element of the slice
func ForEach[T any](source []T, f func(T)) {
	for i := range source {
		f(source[i])
	}
}

// While executes the given function on every element of the slice until the function returns false
func While[T any](source []T, f func(T) bool) {
	for i := range source {
		if ok := f(source[i]); !ok {
			return
		}
	}
}

// Any tests if any of the elements of the slice match the predicate
func Any[T any](source []T, predicate func(T) bool) bool {
	for i := range source {
		if predicate(source[i]) {
			return true
		}
	}
	return false
}

// All tests if all of the elements of the slice match the predicate
func All[T any](source []T, predicate func(T) bool) bool {
	for i := range source {
		if !predicate(source[i]) {
			return false
		}
	}
	return true
}

// Flatten flattens a slice of slices into a single slice
func Flatten[T any](source [][]T) []T {
	out := make([]T, 0)

	for i := range source {
		for j := range source[i] {
			out = append(out, source[i][j])
		}
	}

	return out
}

// Expand creates a slice by executing the generator function count times
func Expand[T any](generator func(idx int) T, count int) []T {
	out := make([]T, count)

	for i := 0; i < count; i++ {
		out[i] = generator(i)
	}
	return out
}

// Distinct returns a new slice where duplicate entries are removed
func Distinct[T comparable](source []T) []T {
	seen := make(map[T]struct{})
	out := make([]T, 0, len(source))
	for i := range source {
		item := source[i]
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		out = append(out, item)
	}

	return out
}
