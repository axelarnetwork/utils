package maps

// Has returns true if the given key is included in the map
func Has[T1 comparable, T2 any](m map[T1]T2, key T1) bool {
	_, ok := m[key]
	return ok
}

// Filter returns a new map that only contains elements that match the predicate
func Filter[T1 comparable, T2 any](m map[T1]T2, predicate func(key T1, value T2) bool) map[T1]T2 {
	filtered := make(map[T1]T2, len(m))

	for k, v := range m {
		if predicate(k, v) {
			filtered[k] = v
		}
	}
	return filtered
}
