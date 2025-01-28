package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

const randomStrLength = 16

func getRandomString() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := make([]byte, randomStrLength)
	for i := range result {
		result[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(result)
}

func resetLoggerConf() {
	output = os.Stdout
	handler.Store(0)
	logLevel.Set(slog.LevelWarn)
	globalLogger = slog.New(
		slog.NewJSONHandler(
			output,
			&slog.HandlerOptions{Level: logLevel},
		),
	)
}

func changeStdout() (*os.File, *os.File, func()) {
	oldStdout := os.Stdout

	r, w, _ := os.Pipe()
	closer := func() {
		os.Stdout = oldStdout
		_ = r.Close()
		_ = w.Close()
	}

	os.Stdout = w

	return r, w, closer
}

func TestLog_ValidateDefaults(t *testing.T) {
	require.Equal(t, slog.LevelWarn, logLevel.Level())
	assert.Equal(t, os.Stdout, output)
}

func TestLog_WithLogLevel(t *testing.T) {
	defer resetLoggerConf()

	t.Run("debug", func(t *testing.T) {
		WithLogLevel("debug")()
		require.Equal(t, slog.LevelDebug, logLevel.Level())
	})

	t.Run("info", func(t *testing.T) {
		WithLogLevel("info")()
		require.Equal(t, slog.LevelInfo, logLevel.Level())
	})

	t.Run("warn", func(t *testing.T) {
		WithLogLevel("warn")()
		require.Equal(t, slog.LevelWarn, logLevel.Level())
	})

	t.Run("error", func(t *testing.T) {
		WithLogLevel("error")()
		require.Equal(t, slog.LevelError, logLevel.Level())
	})

	t.Run("empty string", func(t *testing.T) {
		WithLogLevel("")()
		require.Equal(t, slog.LevelWarn, logLevel.Level())
	})

	t.Run("random string", func(t *testing.T) {
		val := getRandomString()
		WithLogLevel(val)()
		require.Equal(t, slog.LevelWarn, logLevel.Level())
	})
}

func TestLog_WithOutput(t *testing.T) {
	defer resetLoggerConf()

	r, w, closer := changeStdout()
	defer closer()

	WithOutput(w)()
	val := getRandomString()
	Error(val)

	_ = w.Sync()
	_ = w.Close()
	out := &bytes.Buffer{}
	_, _ = io.Copy(out, r)

	require.Contains(t, out.String(), val)
}

func TestLog_WithXXXFormat(t *testing.T) {
	defer resetLoggerConf()

	t.Run("WithJSONFormat", func(t *testing.T) {
		WithJSONFormat()()

		handler := reflect.ValueOf(globalLogger.Handler())
		if handler.Kind() == reflect.Ptr {
			handler = handler.Elem()
		}

		require.Equal(t, "JSONHandler", handler.Type().Name())
	})

	t.Run("WithTextFormat", func(t *testing.T) {
		WithTextFormat()()

		handler := reflect.ValueOf(globalLogger.Handler())
		if handler.Kind() == reflect.Ptr {
			handler = handler.Elem()
		}

		require.Equal(t, "TextHandler", handler.Type().Name())
	})
}

func TestLog_Configure(t *testing.T) {
	defer resetLoggerConf()

	t.Run("WithJSONFormat", func(t *testing.T) {
		Configure(WithJSONFormat())

		handler := reflect.ValueOf(globalLogger.Handler())
		if handler.Kind() == reflect.Ptr {
			handler = handler.Elem()
		}

		require.Equal(t, "JSONHandler", handler.Type().Name())
	})

	t.Run("WithTextFormat", func(t *testing.T) {
		Configure(WithTextFormat())

		handler := reflect.ValueOf(globalLogger.Handler())
		if handler.Kind() == reflect.Ptr {
			handler = handler.Elem()
		}

		require.Equal(t, "TextHandler", handler.Type().Name())
	})

	t.Run("WithOutput", func(t *testing.T) {
		t.Run("correct io.Writer", func(t *testing.T) {
			r, w, closer := changeStdout()
			defer closer()

			Configure(WithOutput(w))
			val := getRandomString()
			Error(val)

			_ = w.Close()
			out := &bytes.Buffer{}
			_, _ = io.Copy(out, r)

			require.Contains(t, out.String(), val)
		})

		t.Run("nil value", func(t *testing.T) {
			Configure(WithOutput(nil))

			handler := reflect.ValueOf(globalLogger.Handler())
			if handler.Kind() == reflect.Ptr {
				handler = handler.Elem()
			}
			writerField := handler.FieldByName("w")
			require.True(t, writerField.IsValid())

			// Use reflection to access the value of the unexported field
			writer := reflect.NewAt(writerField.Type(), unsafe.Pointer(writerField.UnsafeAddr())).Elem().Interface()
			require.Implements(t, (*io.Writer)(nil), writer)
			assert.Equal(t, os.Stdout, writer)
		})

		t.Run("nil pointer", func(t *testing.T) {
			Configure(WithOutput((*os.File)(nil)))

			handler := reflect.ValueOf(globalLogger.Handler())
			if handler.Kind() == reflect.Ptr {
				handler = handler.Elem()
			}
			writerField := handler.FieldByName("w")
			require.True(t, writerField.IsValid())

			// Use reflection to access the value of the unexported field
			writer := reflect.NewAt(writerField.Type(), unsafe.Pointer(writerField.UnsafeAddr())).Elem().Interface()
			require.Implements(t, (*io.Writer)(nil), writer)
			assert.Equal(t, os.Stdout, writer)
		})
	})

	t.Run("WithLogLevel", func(t *testing.T) {
		Configure(WithLogLevel("debug"))
		require.Equal(t, logLevel.Level(), slog.LevelDebug)
	})

	t.Run("several", func(t *testing.T) {
		r, w, closer := changeStdout()
		defer closer()
		Configure(WithTextFormat(), WithLogLevel("info"), WithOutput(w))

		require.Equal(t, logLevel.Level(), slog.LevelInfo)

		val := getRandomString()
		Error(val)
		_ = w.Close()
		out := &bytes.Buffer{}
		_, _ = io.Copy(out, r)
		require.Contains(t, out.String(), val)

		handler := reflect.ValueOf(globalLogger.Handler())
		if handler.Kind() == reflect.Ptr {
			handler = handler.Elem()
		}

		require.Equal(t, "TextHandler", handler.Type().Name())
	})
}

func TestLog_Emitters(t *testing.T) {
	defer resetLoggerConf()

	t.Run("debug", func(t *testing.T) {
		r, w, closer := changeStdout()
		defer closer()

		Configure(WithLogLevel("debug"), WithOutput(w))
		val := getRandomString()
		Debug(val)

		_ = w.Close()
		out := &bytes.Buffer{}
		_, _ = io.Copy(out, r)

		require.Contains(t, out.String(), val)
		assert.Contains(t, out.String(), "\"level\":\"DEBUG\"")
	})

	t.Run("info", func(t *testing.T) {
		r, w, closer := changeStdout()
		defer closer()

		WithLogLevel("info")()
		WithOutput(w)()
		val := getRandomString()
		Info(val)

		_ = w.Sync()
		_ = w.Close()
		out := &bytes.Buffer{}
		_, _ = io.Copy(out, r)

		require.Contains(t, out.String(), val)
		assert.Contains(t, out.String(), "\"level\":\"INFO\"")
	})

	t.Run("warn", func(t *testing.T) {
		r, w, closer := changeStdout()
		defer closer()

		WithLogLevel("warn")()
		WithOutput(w)()
		val := getRandomString()
		Warn(val)

		_ = w.Sync()
		_ = w.Close()
		out := &bytes.Buffer{}
		_, _ = io.Copy(out, r)

		require.Contains(t, out.String(), val)
		assert.Contains(t, out.String(), "\"level\":\"WARN\"")
	})

	t.Run("error", func(t *testing.T) {
		r, w, closer := changeStdout()
		defer closer()

		WithLogLevel("error")()
		WithOutput(w)()
		val := getRandomString()
		Error(val)

		_ = w.Sync()
		_ = w.Close()
		out := &bytes.Buffer{}
		_, _ = io.Copy(out, r)

		require.Contains(t, out.String(), val)
		assert.Contains(t, out.String(), "\"level\":\"ERROR\"")
	})
}
