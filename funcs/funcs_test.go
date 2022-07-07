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

func TestMustNoErr(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNoErr(nil)
	})
	assert.Panics(t, func() { MustNoErr(errors.New("some error")) })
}

func TestNot(t *testing.T) {
	isLessThanOrEqualToZero := Not(func(i int64) bool { return i > 0 })
	isPositive := Not(isLessThanOrEqualToZero)

	assert.False(t, isLessThanOrEqualToZero(1))
	assert.True(t, isLessThanOrEqualToZero(0))
	assert.True(t, isLessThanOrEqualToZero(-1))

	assert.True(t, isPositive(1))
	assert.False(t, isPositive(0))
	assert.False(t, isPositive(-1))
}

func TestAnd(t *testing.T) {
	isPositive := func(i int64) bool { return i > 0 }
	isEven := func(i int64) bool { return i%2 == 0 }
	isLessThanFive := func(i int64) bool { return i < 5 }

	and := And(isPositive, isEven, isLessThanFive)

	assert.False(t, and(0))
	assert.False(t, and(1))
	assert.True(t, and(2))
	assert.False(t, and(3))
	assert.True(t, and(4))
	assert.False(t, and(5))
}

func TestOr(t *testing.T) {
	isNegative := func(i int64) bool { return i < 0 }
	isEven := func(i int64) bool { return i%2 == 0 }
	isLessThanFive := func(i int64) bool { return i < 5 }

	or := Or(isNegative, isEven, isLessThanFive)

	assert.True(t, or(0))
	assert.True(t, or(1))
	assert.True(t, or(2))
	assert.True(t, or(3))
	assert.True(t, or(4))
	assert.False(t, or(5))
	assert.True(t, or(6))
	assert.True(t, or(-1))
	assert.True(t, or(-2))
}
