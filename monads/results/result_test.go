package results_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/axelarnetwork/utils/monads/results"
)

func TestResult(t *testing.T) {
	t.Run("constructors", func(t *testing.T) {
		res1 := results.New(5, nil)
		assert.Equal(t, 5, res1.Ok())
		assert.NoError(t, res1.Err())

		res2 := results.New("", errors.New("some error"))
		assert.Error(t, res2.Err())

		res3 := results.New("should not be a value", errors.New("some error"))
		assert.Error(t, res3.Err())
		assert.Equal(t, "", res3.Ok())

		assert.Equal(t, results.New(6, nil), results.FromOk(6))
		assert.Equal(t, results.New(0, errors.New("some error")), results.FromErr[int](errors.New("some error")))
	})

	t.Run("Pipe", func(t *testing.T) {
		assert.Error(t, results.Pipe(successfulFunc(3), unsuccessfulFunc).Err())
		assert.Error(t, results.Pipe(unsuccessfulFunc("fail"), successfulFunc).Err())

		assert.NoError(t, results.Pipe(successfulFunc(7), successfulFunc2).Err())
		assert.Equal(t, '7', results.Pipe(successfulFunc(7), successfulFunc2).Ok())
	})

	t.Run("Try", func(t *testing.T) {
		assert.Error(t, results.Try(unsuccessfulFunc("fail"), strconv.Itoa).Err())
		assert.NoError(t, results.Try(successfulFunc(20), func(s string) bool { return len(s) == 2 }).Err())
	})
}

func successfulFunc(i int) results.Result[string] {
	return results.New(strconv.Itoa(i), nil)
}

func successfulFunc2(s string) results.Result[rune] {
	return results.New([]rune(s)[0], nil)
}

func unsuccessfulFunc(string) results.Result[int] {
	return results.New(0, errors.New("some error"))
}
