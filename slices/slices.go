package slices

import "fmt"

// Map maps a slice of T to a slice of S
func Map[T, S any](source []T, f func(T) S) []S {
	out := make([]S, len(source))

	// avoid allocating a copy of the slice element
	for i := range source {
		out[i] = f(source[i])
	}

	return out
}

// FlatMap composes Map and Flatten
func FlatMap[T, S any](source []T, f func(T) []S) []S {
	return Flatten(Map(source, f))
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

// ToMap returns a map from the given slice with keys associated by the lookup function. Panics if strictUniqueness is set, otherwise overrides values for colliding keys.
func ToMap[T1 any, T2 comparable](source []T1, lookup func(T1) T2, strictUniqueness ...bool) map[T2]T1 {
	m := make(map[T2]T1, len(source))

	strict := false
	if len(strictUniqueness) > 0 {
		strict = strictUniqueness[0]
	}

	for i := range source {
		key := lookup(source[i])
		if strict {
			// separate the conditions so the lookup is only done if strict flag is set
			if value, ok := m[key]; ok {
				panic(fmt.Sprintf("key %v is not unique, points to %v and %v", key, value, source[i]))
			}
		}

		m[key] = source[i]
	}

	return m
}

// TryCast tries to cast each element of the slice from type T1 to type T2. Elements get filtered out of the cast is unsuccessful.
func TryCast[T1 any, T2 any](source []T1) []T2 {
	out := make([]T2, 0, cap(source))
	for i := range source {
		// T1 could be a concrete type, so the compiler does not allow a type assertion without boxing T1 first
		if x, ok := (interface{}(source[i])).(T2); ok {
			out = append(out, x)
		}
	}
	return out
}

// Reverse returns a new slice in which the elements aree listed in reverse order
func Reverse[T1 any](source []T1) []T1 {
	out := make([]T1, 0, cap(source))
	for i := len(source) - 1; i >= 0; i-- {
		out = append(out, source[i])
	}
	return out
}

// GroupBy returns a map with given items each group into a slice
func GroupBy[T any, K comparable](source []T, fn func(T) K) map[K][]T {
	results := make(map[K][]T)

	for _, s := range source {
		k := fn(s)
		results[k] = append(results[k], s)
	}

	return results
}
