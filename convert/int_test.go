package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntToBytes(t *testing.T) {
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 1, 0}, IntToBytes(256))
	assert.Equal(t, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, IntToBytes(-1))
}
