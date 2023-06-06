package log

import (
	"context"
	"fmt"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// Logger is a simple interface to log at three log levels with additional formatting methods for convenience
type Logger interface {
	Debug(msg string)
	Debugf(format string, a ...any)
	Info(msg string)
	Infof(format string, a ...any)
	Error(msg string)
	Errorf(format string, a ...any)
}

var (
	defaultLogger = logWrapper{tmlog.NewNopLogger()}
	frozen        bool
)

// Setup sets the logger that the application should use. The default is a nop logger, i.e. all logs are discarded.
// Panics if called more than once without calling Reset first.
func Setup(logger tmlog.Logger) {
	if frozen {
		panic("logger was already set")
	}

	defaultLogger = logWrapper{logger}
	frozen = true
}

// Reset returns the logger state to the default nop logger and enables Setup to be called again.
func Reset() {
	defaultLogger = logWrapper{tmlog.NewNopLogger()}
	frozen = false
}

// Debug logs the given msg at log level DEBUG
func Debug(msg string) {
	defaultLogger.Debug(msg)
}

// Debugf logs the given formatted msg at log level DEBUG
func Debugf(format string, a ...any) {
	defaultLogger.Debugf(format, a...)
}

// Info logs the given msg at log level INFO
func Info(msg string) {
	defaultLogger.Info(msg)
}

// Infof logs the given formatted msg at log level INFO
func Infof(format string, a ...any) {
	defaultLogger.Infof(format, a...)
}

// Error logs the given msg at log level ERROR
func Error(msg string) {
	defaultLogger.Error(msg)
}

// Errorf logs the given formatted msg at log level ERROR
func Errorf(format string, a ...any) {
	defaultLogger.Errorf(format, a...)
}

type key int

var logKey key

// AppendKeyVals adds the given keyvals to the given context. If the context already stores keyvals, the new ones get appended.
// This should be used to track the logging context across the application.
func AppendKeyVals(ctx context.Context, keyvals ...any) context.Context {
	if len(keyvals)%2 != 0 {
		return ctx
	}

	existingKeyvals, _ := ctx.Value(logKey).([]any)

	return context.WithValue(ctx, logKey, append(existingKeyvals, keyvals...))
}

// FromCtx reads logging keyvals from the given context if any and adds them to any logs the returned Logger puts out
func FromCtx(ctx context.Context) Logger {
	keyVals, ok := ctx.Value(logKey).([]any)

	if !ok || len(keyVals) == 0 {
		return defaultLogger
	}

	return logWrapper{defaultLogger.With(keyVals...)}
}

// WithKeyVals returns a logger that adds the given keyvals to any logs it puts out.
// This should be used for immediate log enrichment, for tracking of a logging context across the application use AppendKeyVals
func WithKeyVals(keyvals ...any) Logger {
	if len(keyvals)%2 != 0 {
		return defaultLogger
	}

	return logWrapper{defaultLogger.With(keyvals...)}
}

type logWrapper struct {
	tmlog.Logger
}

func (l logWrapper) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l logWrapper) Debugf(format string, a ...any) {
	l.Logger.Debug(fmt.Sprintf(format, a...))
}

func (l logWrapper) Info(msg string) {
	l.Logger.Info(msg)
}

func (l logWrapper) Infof(format string, a ...any) {
	l.Logger.Info(fmt.Sprintf(format, a...))
}

func (l logWrapper) Error(msg string) {
	l.Logger.Error(msg)
}

func (l logWrapper) Errorf(format string, a ...any) {
	l.Logger.Error(fmt.Sprintf(format, a...))
}
