package reflect_test

import (
	"github.com/axelarnetwork/utils/reflect"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreate(t *testing.T) {
	assert.NotPanics(t, func() {
		val := reflect.Create[struct{}]()
		assert.NotNil(t, val)
	})

	assert.NotPanics(t, func() {
		val := reflect.Create[*struct{}]()
		assert.NotNil(t, val)
	})

	assert.NotPanics(t, func() {
		val := reflect.Create[*[]byte]()
		assert.NotNil(t, val)
	})

	assert.NotPanics(t, func() {
		val := reflect.Create[int]()
		assert.NotNil(t, val)
	})

	assert.Panics(t, func() {
		_ = reflect.Create[**struct{}]()
	})

	assert.Panics(t, func() {
		_ = reflect.Create[any]()
	})

	assert.Panics(t, func() {
		_ = reflect.Create[*any]()
	})

	assert.Panics(t, func() {
		_ = reflect.Create[testInterface]()
	})
}

func TestIsInterface(t *testing.T) {
	assert.False(t, reflect.IsInterface[int]())
	assert.False(t, reflect.IsInterface[**struct{}]())

	assert.True(t, reflect.IsInterface[any]())
	assert.True(t, reflect.IsInterface[*any]())
	assert.True(t, reflect.IsInterface[testInterface]())
}

type testInterface interface {
	doWork()
}
