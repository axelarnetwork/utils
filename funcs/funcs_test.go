package funcs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {
	f := func(s string) int { return len(s) }
	g := func(i int) bool { return i%2 == 0 }
	h := Compose(f, g)

	assert.False(t, h("hello"))
}
