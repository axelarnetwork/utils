package checks_test

import (
	"github.com/axelarnetwork/utils/checks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFieldDeclarations(t *testing.T) {
	cmd := checks.FieldDeclarations()
	assert.NoError(t, cmd.Flags().Set(checks.ExcludeFieldsFlag, "XXX_"))
	assert.NoError(t, cmd.Flags().Set(checks.ExcludeTypesFlag, "abi.Argument,cobra.Command,proto/tendermint/types.Header,zerolog.ConsoleWriter,packages.Config"))

	out := &testWriter{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"./..."})
	assert.NoError(t, cmd.Execute())
	assert.Equal(t, 0, out.Called, out.String())

	out = &testWriter{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"./testdata"})
	assert.NoError(t, cmd.Execute())
	assert.Equal(t, 11, out.Called, out.String())

	out = &testWriter{}
	cmd.SetOut(out)
	cmd.SetArgs([]string{"github.com/axelarnetwork/utils/checks/testdata"})
	assert.NoError(t, cmd.Execute())
	assert.Equal(t, 11, out.Called, out.String())
}

type testWriter struct {
	Called int
	out    []byte
}

func (t *testWriter) Write(bz []byte) (int, error) {
	t.Called++
	t.out = append(t.out, bz...)

	return len(bz), nil
}

func (t *testWriter) String() string {
	return string(t.out)
}
