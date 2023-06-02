package log

// NOPLogger is a logger that doesn't do anything
type NOPLogger struct{}

// Interface assertions
var _ BaseLogger = NOPLogger{}

// NewNOPLogger returns a logger that doesn't do anything.
func NewNOPLogger() BaseLogger { return NOPLogger{} }

// Info does nothing
func (NOPLogger) Info(string, ...any) {}

// Debug does nothing
func (NOPLogger) Debug(string, ...any) {}

// Error does nothing
func (NOPLogger) Error(string, ...any) {}

// With does nothing
func (l NOPLogger) With(...any) BaseLogger {
	return l
}
