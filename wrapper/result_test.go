package wrapper_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/axelarnetwork/utils/wrapper"
	"github.com/stretchr/testify/assert"
)

func TestResult(t *testing.T) {
	res1 := wrapper.NewResult(5, nil)
	assert.Equal(t, 5, res1.Ok())
	assert.NoError(t, res1.Error())

	res2 := wrapper.NewResult("", errors.New("some error"))
	assert.Error(t, res2.Error())

	res3 := wrapper.NewResult("should not be a value", errors.New("some error"))
	assert.Error(t, res3.Error())
	assert.Equal(t, "", res3.Ok())

	assert.Error(t, wrapper.ContinueWithResult(successfulFunc(3), unsuccessfulFunc).Error())
	assert.Error(t, wrapper.ContinueWithResult(unsuccessfulFunc("fail"), successfulFunc).Error())

	assert.NoError(t, wrapper.ContinueWithResult(successfulFunc(7), successfulFunc2).Error())
	assert.Equal(t, '7', wrapper.ContinueWithResult(successfulFunc(7), successfulFunc2).Ok())
}

func successfulFunc(i int) wrapper.Result[string] {
	return wrapper.NewResult(strconv.Itoa(i), nil)
}

func successfulFunc2(s string) wrapper.Result[rune] {
	return wrapper.NewResult([]rune(s)[0], nil)
}

func unsuccessfulFunc(string) wrapper.Result[int] {
	return wrapper.NewResult(0, errors.New("some error"))
}
