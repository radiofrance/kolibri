// log package is based on zap logger, but only with minimal required methods by the controller.
// Other loggers can be wrapped through log.Logger interface.
package log

import (
	"go.uber.org/zap"
)

// Field is an alias for zap.Field. A Field is a marshaling operation used to add
// a key-value pair to a logger's context. Most fields are lazily marshaled, so it's inexpensive to add fields
// to disabled debug-level log statements.
type Field zap.Field

// A Logger provides fast, leveled, structured logging. All methods are safe
// for concurrent use.
type Logger interface {
	// With creates a child logger and adds structured context to it. Fields added to the child don't affect the parent,
	// and vice versa.
	With(fields ...Field) Logger
	// Named adds a new path segment to the logger's name. Segments are joined by periods. By default, Loggers are
	// unnamed.
	Named(s string) Logger

	// Debug logs a message at DebugLevel. The message includes any fields passed at the log site, as well as any fields
	// accumulated on the logger.
	Debug(message string)
	// Debugf logs a message according to a format specifier at DebugLevel. The message includes any fields passed at
	// the log site, as well as any fields accumulated on the logger.
	Debugf(format string, a ...interface{})
	// Info logs a message at InfoLevel. The message includes any fields passed at the log site, as well as any fields
	// accumulated on the logger.
	Info(message string)
	// Infof logs a message according to a format specifier at InfoLevel. The message includes any fields passed at the
	// log site, as well as any fields accumulated on the logger.
	Infof(format string, a ...interface{})
	// Warn logs a message at WarnLevel. The message includes any fields passed at the log site, as well as any fields
	// accumulated on the logger.
	Warn(message string)
	// Warnf logs a message according to a format specifier at WarnLevel. The message includes any fields passed at the
	// log site, as well as any fields accumulated on the logger.
	Warnf(format string, a ...interface{})
	// Error logs a message at ErrorLevel. The message includes any fields passed at the log site, as well as any fields
	// accumulated on the logger.
	Error(message string)
	// Errorf logs a message according to a format specifier at ErrorLevel. The message includes any fields passed at
	// the log site, as well as any fields accumulated on the logger.
	Errorf(format string, a ...interface{})
	// Fatal logs a message at FatalLevel. The message includes any fields passed at the log site, as well as any fields
	// accumulated on the logger.
	Fatal(message string)
	// Fatalf logs a message according to a format specifier at FatalLevel. The message includes any fields passed at
	// the log site, as well as any fields accumulated on the logger.
	Fatalf(format string, a ...interface{})
	// Panic logs a message at PanicLevel. The message includes any fields passed at the log site, as well as any fields
	// accumulated on the logger.
	Panic(message string)
	// Panicf logs a message according to a format specifier at PanicLevel. The message includes any fields passed at
	// the log site, as well as any fields accumulated on the logger.
	Panicf(format string, a ...interface{})
}
