package reflect

import (
	"reflect"
)

// Create returns a new value of type T, where T is any regular value or reference type. Panics when T is a multi-pointer (e.g. **string).
func Create[T any]() T {
	t := reflect.TypeOf(new(T)).Elem()

	ptrCount := 0

	for t.Kind() == reflect.Ptr {
		ptrCount++
		t = t.Elem()
	}

	if ptrCount > 1 {
		panic("only single pointers are allowed")
	}

	value := reflect.New(t)

	if ptrCount > 0 {
		return value.Interface().(T)
	}

	return value.Elem().Interface().(T)
}
