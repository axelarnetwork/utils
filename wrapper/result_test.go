package wrapper_test

import (
	"errors"
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
}
