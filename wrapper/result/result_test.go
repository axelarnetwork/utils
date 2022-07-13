package result_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/axelarnetwork/utils/wrapper/result"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/strings"
)

func TestResult(t *testing.T) {
	t.Run("constructors", func(t *testing.T) {
		res1 := result.New(5, nil)
		assert.Equal(t, 5, res1.Ok())
		assert.NoError(t, res1.Error())

		res2 := result.New("", errors.New("some error"))
		assert.Error(t, res2.Error())

		res3 := result.New("should not be a value", errors.New("some error"))
		assert.Error(t, res3.Error())
		assert.Equal(t, "", res3.Ok())

		assert.Equal(t, result.New(6, nil), result.FromOk(6))
		assert.Equal(t, result.New(0, errors.New("some error")), result.FromErr[int](errors.New("some error")))
	})

	t.Run("Pipe", func(t *testing.T) {
		assert.Error(t, result.Pipe(successfulFunc(3), unsuccessfulFunc).Error())
		assert.Error(t, result.Pipe(unsuccessfulFunc("fail"), successfulFunc).Error())

		assert.NoError(t, result.Pipe(successfulFunc(7), successfulFunc2).Error())
		assert.Equal(t, '7', result.Pipe(successfulFunc(7), successfulFunc2).Ok())
	})

	t.Run("Try", func(t *testing.T) {
		assert.Error(t, result.Try(unsuccessfulFunc("fail"), strconv.Itoa).Error())
		assert.NoError(t, result.Try(successfulFunc(20), strings.IsASCIIText).Error())
	})
}

func successfulFunc(i int) result.Result[string] {
	return result.New(strconv.Itoa(i), nil)
}

func successfulFunc2(s string) result.Result[rune] {
	return result.New([]rune(s)[0], nil)
}

func unsuccessfulFunc(string) result.Result[int] {
	return result.New(0, errors.New("some error"))
}
