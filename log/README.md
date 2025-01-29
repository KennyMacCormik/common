# README

## Overview
The `log` package provides flexible, thread-safe global logging utility built on top of `log/slog`. It supports multiple output formats, dynamic log levels, and configurable outputs for seamless integration into your Go projects. The package is designed for simplicity and configurability, ensuring developers can adjust logging behavior to suit various environments.

---

## Features
- **Dynamic Log Levels**: Adjust log levels at runtime (`debug`, `info`, `warn`, `error`).
- **Multiple Output Formats**: Supports both JSON and text-based logging.
- **Customizable Output**: Redirect logs to `os.Stdout`, files, or any `io.Writer` implementation.
- **Thread-Safe**: Ensures global logger operations are safe for concurrent use.
- **Default Configuration**: Pre-configured for JSON format, `warn` level, and `os.Stdout` output.
- **Composable Options**: Use the `Configure` function to chain multiple options (e.g., log level, format, output).

---

## Installation

### Dependencies
Ensure you have Go 1.23 or later installed. This package relies on the standard library and the following third-party dependencies:
- `github.com/stretchr/testify` (for unit testing)

### Package Installation
To install the `log` package, run the following command:

```bash
go get github.com/KennyMacCormik/common/log
```

---

## Usage

### Basic Example
```go
package main

import (
	"github.com/KennyMacCormik/common/log"
)

func main() {
	log.Configure(
		log.WithJSONFormat(),
		log.WithLogLevel("info"),
	)

	log.Info("Starting http server")
	log.Debug("This won't appear due to log level")
	log.Error("An error occurred")
}
```

**Output:**
```json
{"time":"2025-01-29T01:37:34.931171+03:00","level":"INFO","msg":"Starting http server"}
{"time":"2025-01-29T01:37:34.931171+03:00","level":"ERROR","msg":"An error occurred"}
```

### Example Without `Configure`
```go
package main

import (
	"github.com/KennyMacCormik/common/log"
)

func main() {
	log.Info("Application started with default configuration")
	log.Warn("This is a warning")
	log.Error("An error occurred")
}
```

**Output:**
```json
{"time":"2025-01-29T01:44:46.887+03:00","level":"INFO","msg":"Application started with default configuration"}
{"time":"2025-01-29T01:44:46.887+03:00","level":"WARN","msg":"This is a warning"}
{"time":"2025-01-29T01:44:46.887+03:00","level":"ERROR","msg":"An error occurred"}
```

### Advanced Example
```go
package main

import (
	"github.com/KennyMacCormik/common/log"
	"os"
)

func main() {
	// Redirect logs to a file
	file, err := os.Create("app.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.Configure(
		log.WithTextFormat(),
		log.WithLogLevel("debug"),
		log.WithOutput(file),
	)

	log.Debug("Debugging application state")
	log.Info("Informational message")
}
```

**Output in `app.log`:**
```
time=2025-01-29T01:44:46.887+03:00 level=DEBUG msg="Debugging application state"
time=2025-01-29T01:44:46.887+03:00 level=INFO msg="Informational message"
```

---

## API Documentation

### Exported Functions

#### `func CopyLogger(msg string, args ...any)`
CopyLogger copies the global logger and returns it.

#### `func Configure(options ...LoggingOptions)`
Configures the global logger with the specified options. Options can include log level, format, and output.

#### `func WithLogLevel(level string) LoggingOptions`
Sets the log level. Accepted values: `debug`, `info`, `warn`, `error`. Defaults to `warn` for invalid values.

#### `func WithJSONFormat() LoggingOptions`
Configures the logger to use JSON output format.

#### `func WithTextFormat() LoggingOptions`
Configures the logger to use text output format.

#### `func WithOutput(out io.Writer) LoggingOptions`
Redirects the logger output to the specified `io.Writer`. Defaults to `os.Stdout` if `nil` or invalid values are provided.

#### `func Debug(msg string, args ...any)`
Logs a message at the `DEBUG` level.

#### `func Info(msg string, args ...any)`
Logs a message at the `INFO` level.

#### `func Warn(msg string, args ...any)`
Logs a message at the `WARN` level.

#### `func Error(msg string, args ...any)`
Logs a message at the `ERROR` level.

---

## Type Descriptions

### Interfaces

#### `type LoggingOptions func()`
Represents a functional option for configuring the global logger.

---

## Variable Descriptions

---

## Package Behavior

The `log` package initializes the logger with the following default configuration:
- JSON format
- Warn log level
- `os.Stdout` as the output

#### Attribute Processing
The `Debug`, `Info`, `Warn`, and `Error` functions emit a log record with the current time, level, and message. Attributes are processed as follows:
- If an argument is an `slog.Attr`, it is used as is.
- If an argument is a string and not the last argument, the next argument is treated as its value, forming a key-value pair.
- Otherwise, the argument is treated as a value with the key `"!BADKEY"`.

---

## License
This project is licensed under the MIT License. See the [LICENSE](https://opensource.org/licenses/MIT) file for details.

---

## Thanks
Special thanks to the contributors and maintainers of the following packages:
- [log/slog](https://pkg.go.dev/log/slog)
- [stretchr/testify](https://pkg.go.dev/github.com/stretchr/testify)

