package reflect

import (
	"fmt"
	"reflect"
)

// Create returns a new value of type T, where T is any regular value or reference type. Panics when T is a multi-pointer (e.g. **string) or interface type.
func Create[T any]() T {
	t := reflect.TypeOf(new(T)).Elem()

	ptrCount := 0

	for t.Kind() == reflect.Pointer {
		ptrCount++
		t = t.Elem()
	}

	if ptrCount > 1 {
		panic(fmt.Sprintf("only single pointers are allowed, got %d instead", ptrCount))
	}

	if t.Kind() == reflect.Interface {
		panic(fmt.Sprintf("only concrete types can be instantiated, got %s instead", t.String()))
	}

	value := reflect.New(t)

	if ptrCount > 0 {
		return value.Interface().(T)
	}

	return value.Elem().Interface().(T)
}

// IsInterface returns true if the generic type is an interface and not a concrete type
func IsInterface[T any]() bool {
	t := reflect.TypeOf(new(T)).Elem()

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t.Kind() == reflect.Interface
}
