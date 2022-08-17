package maps

// Has returns true if the given key is included in the map
func Has[T1 comparable, T2 any](m map[T1]T2, key T1) bool {
	_, ok := m[key]
	return ok
}
