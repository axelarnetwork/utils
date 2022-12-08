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
		val := reflect.Create[**struct{}]()
		assert.NotNil(t, val)
	})

}
