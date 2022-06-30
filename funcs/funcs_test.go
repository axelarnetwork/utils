package funcs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompose(t *testing.T) {
	f := func(s string) int { return len(s) }
	g := func(i int) bool { return i%2 == 0 }
	h := Compose(f, g)

	assert.False(t, h("hello"))
}

func TestIdentity(t *testing.T) {
	assert.Equal(t, "a", Identity("a"))
}

func TestMust(t *testing.T) {
	withResult := func() (int, error) { return 9, nil }
	withErr := func() (string, error) { return "", errors.New("some error") }

	assert.Equal(t, 9, Must(withResult()))
	assert.Panics(t, func() { Must(withErr()) })
}
