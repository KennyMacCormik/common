// Package log provides a thread-safe implementation of a global logger.
// It uses slog.Logger under the hood.
//
// Package log initializes the logger with the following default configuration:
//   - JSON format
//   - Warn log level
//   - os.Stdout as the output
//
// Debug, Info, Warn, Error emits a log record with the current time and the given level and message.
// Their attributes are processed as follows:
//   - If an argument is a slog.Attr, it is used as is.
//   - If an argument is a string and not the last argument, the next argument is treated as its value, forming a key-value pair.
//   - Otherwise, the argument is treated as a value with the key "!BADKEY".
package log

import (
	"io"
	"log/slog"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
)

func init() {
	logLevel = new(slog.LevelVar)
	output = os.Stdout
	handler.Store(0)
	logLevel.Set(slog.LevelWarn)
	globalLogger = slog.New(
		slog.NewJSONHandler(
			output,
			&slog.HandlerOptions{Level: logLevel},
		),
	)
	globalLogger.Debug("logger init success", "log level", "warn", "output", "os.Stdout", "format", "json")
}

// LoggingOptions represents a configuration option for the logger.
type LoggingOptions func()

var (
	globalLogger *slog.Logger
	logLevel     *slog.LevelVar
	output       io.Writer
	handler      atomic.Int64 // 0 = JSON, 1 = Text
	mtx          sync.Mutex
)

// WithJSONFormat configures the logger to use JSON output format.
// If provided alongside WithTextFormat latest provided wins
func WithJSONFormat() LoggingOptions {
	return func() {
		handler.Store(0)
		storeLogger(output, logLevel.Level())
	}
}

// WithTextFormat configures the logger to use text output format.
// If provided alongside WithJSONFormat latest provided wins
func WithTextFormat() LoggingOptions {
	return func() {
		handler.Store(1)
		storeLogger(output, logLevel.Level())
	}
}

// WithOutput sets the output for the logger. The default output is os.Stdout.
// If the provided value is nil or invalid, os.Stdout will be used instead.
func WithOutput(out io.Writer) LoggingOptions {
	return func() {
		mtx.Lock()
		defer mtx.Unlock()

		if isNotNilOrNilPointer(out) {
			output = out
		} else {
			output = os.Stdout
		}
		storeLogger(output, logLevel.Level())
	}
}

// WithLogLevel sets the log level of the logger. If an invalid value is provided, the log level defaults to "warn".
//
// Accepted values are:
//   - debug: equivalent to slog.LevelDebug
//   - info:  equivalent to slog.LevelInfo
//   - warn:  equivalent to slog.LevelWarn
//   - error: equivalent to slog.LevelError
func WithLogLevel(level string) LoggingOptions {
	return func() {
		logLevelMap := map[string]slog.Level{
			"debug": slog.LevelDebug,
			"info":  slog.LevelInfo,
			"warn":  slog.LevelWarn,
			"error": slog.LevelError,
		}

		if level != "debug" && level != "info" && level != "warn" && level != "error" {
			level = "warn"
		}

		logLevel.Set(logLevelMap[level])
	}
}

// Configure applies the provided LoggingOptions to configure the global logger.
func Configure(options ...LoggingOptions) {
	for _, option := range options {
		option()
	}
}

// Debug logs a message at the slog.LevelDebug level.
func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

// Info logs a message at the slog.LevelInfo level.
func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

// Warn logs a message at the slog.LevelWarn level.
func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

// Error logs a message at the slog.LevelError level.
func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}

// isNotNilOrNilPointer checks if the provided io.Writer is not nil, a nil pointer, or a nil interface.
func isNotNilOrNilPointer(out io.Writer) bool {
	if out == nil {
		return false
	}

	val := reflect.ValueOf(out)
	kind := val.Kind()

	if kind == reflect.Ptr || kind == reflect.Interface {
		if val.IsNil() {
			return false
		}
	}

	return true
}

// storeLogger generates new *slog.Logger with supplied values and stores it as global logger
func storeLogger(out io.Writer, level slog.Level) {
	if mtx.TryLock() {
		defer mtx.Unlock()
	}

	if handler.Load() == 0 {
		globalLogger = slog.New(
			slog.NewJSONHandler(
				out,
				&slog.HandlerOptions{Level: level},
			),
		)
	} else {
		globalLogger = slog.New(
			slog.NewTextHandler(
				output,
				&slog.HandlerOptions{Level: level},
			),
		)
	}
}
