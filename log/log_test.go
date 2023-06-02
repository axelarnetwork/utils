package log_test

import (
	"context"
	"github.com/axelarnetwork/utils/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoSetup(t *testing.T) {
	t.Cleanup(log.Reset)

	assert.NotPanics(t, func() {
		log.Error("output")
	})
}

func TestMultipleSetups(t *testing.T) {
	t.Cleanup(log.Reset)

	assert.Panics(t, func() {
		log.Setup(log.NewNOPLogger())
		log.Setup(log.NewNOPLogger())
	})
}

func TestDebug(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	log.Debug("debug")
	assert.Equal(t, "debug", <-output)
	assert.Nil(t, <-keyvals)

	log.Debugf("debug%s", "f")
	assert.Equal(t, "debugf", <-output)
	assert.Nil(t, <-keyvals)
}

func TestInfo(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	log.Info("info")
	assert.Equal(t, "info", <-output)
	assert.Nil(t, <-keyvals)

	log.Infof("info%s", "f")
	assert.Equal(t, "infof", <-output)
	assert.Nil(t, <-keyvals)
}

func TestError(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	log.Error("error")
	assert.Equal(t, "error", <-output)
	assert.Nil(t, <-keyvals)

	log.Errorf("error%s", "f")
	assert.Equal(t, "errorf", <-output)
	assert.Nil(t, <-keyvals)
}

func TestDebugWithCtx(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	ctx := context.Background()

	log.FromCtx(ctx).Debug("debug")
	assert.Equal(t, "debug", <-output)
	assert.Nil(t, <-keyvals)

	ctx = log.AppendKeyVals(ctx, "key1", "val1", "key2", 2)

	log.FromCtx(ctx).Debug("debug2")
	assert.Equal(t, "debug2", <-output)
	assert.Equal(t, []any{"key1", "val1", "key2", 2}, <-keyvals)

	ctx = log.AppendKeyVals(ctx, "key3", true)

	log.FromCtx(ctx).Debugf("debug%d", 3)
	assert.Equal(t, "debug3", <-output)
	assert.Equal(t, []any{"key1", "val1", "key2", 2, "key3", true}, <-keyvals)
}

func TestInfoWithCtx(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	ctx := context.Background()

	log.FromCtx(ctx).Info("info")
	assert.Equal(t, "info", <-output)
	assert.Nil(t, <-keyvals)

	ctx = log.AppendKeyVals(ctx, "key1", "val1", "key2", 2)

	log.FromCtx(ctx).Info("info2")
	assert.Equal(t, "info2", <-output)
	assert.Equal(t, []any{"key1", "val1", "key2", 2}, <-keyvals)

	ctx = log.AppendKeyVals(ctx, "key3", true)

	log.FromCtx(ctx).Infof("info%d", 3)
	assert.Equal(t, "info3", <-output)
	assert.Equal(t, []any{"key1", "val1", "key2", 2, "key3", true}, <-keyvals)
}

func TestErrorWithCtx(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	ctx := context.Background()

	log.FromCtx(ctx).Error("error")
	assert.Equal(t, "error", <-output)
	assert.Nil(t, <-keyvals)

	ctx = log.AppendKeyVals(ctx, "key1", "val1", "key2", 2)

	log.FromCtx(ctx).Error("error2")
	assert.Equal(t, "error2", <-output)
	assert.Equal(t, []any{"key1", "val1", "key2", 2}, <-keyvals)

	ctx = log.AppendKeyVals(ctx, "key3", true)

	log.FromCtx(ctx).Errorf("error%d", 3)
	assert.Equal(t, "error3", <-output)
	assert.Equal(t, []any{"key1", "val1", "key2", 2, "key3", true}, <-keyvals)
}

func TestWrongKeyVals(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	ctx := log.AppendKeyVals(context.Background(), "key1", "val1", "key2", 2, "key3")

	log.FromCtx(ctx).Debug("test")
	assert.Equal(t, "test", <-output)
	assert.Nil(t, <-keyvals)
}

func TestWithKeyVals(t *testing.T) {
	t.Cleanup(log.Reset)

	output := make(chan string, 1000)
	keyvals := make(chan []any, 1000)

	log.Setup(&testLogger{
		Output:  output,
		Keyvals: keyvals,
	})

	log.WithKeyVals("key1", false).Debug("debug")
	assert.Equal(t, "debug", <-output)
	assert.Equal(t, []any{"key1", false}, <-keyvals)

	log.WithKeyVals("key2", 5).Info("info")
	assert.Equal(t, "info", <-output)
	assert.Equal(t, []any{"key2", 5}, <-keyvals)

	log.WithKeyVals("key3", "17", "invalid").Error("error")
	assert.Equal(t, "error", <-output)
	assert.Nil(t, <-keyvals)
}

type testLogger struct {
	Output  chan<- string
	Keyvals chan<- []any
	keyvals []any
}

func (t *testLogger) Debug(msg string, keyvals ...interface{}) {
	t.Output <- msg
	t.Keyvals <- append(t.keyvals, keyvals...)
}

func (t *testLogger) Info(msg string, keyvals ...interface{}) {
	t.Output <- msg
	t.Keyvals <- append(t.keyvals, keyvals...)
}

func (t *testLogger) Error(msg string, keyvals ...interface{}) {
	t.Output <- msg
	t.Keyvals <- append(t.keyvals, keyvals...)
}

func (t *testLogger) With(keyvals ...interface{}) log.BaseLogger {
	return &testLogger{
		Output:  t.Output,
		Keyvals: t.Keyvals,
		keyvals: append(t.keyvals, keyvals...),
	}
}
