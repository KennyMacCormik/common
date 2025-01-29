# Common Utilities Repository

## Overview

This repository contains general implementations of commonly used features for application development. Each package is designed to simplify and optimize specific tasks such as logging, validation, data conversion, and HTTP router management.

---

## Packages

### [Log](https://github.com/KennyMacCormik/common/tree/main/log)
- **Purpose:** Provides a thread-safe implementation of a global logger using the `slog` library.
- **Key Features:**
    - Default configuration: JSON format, `warn` log level, and output to `os.Stdout`.
    - Logging methods for different levels: `Debug`, `Info`, `Warn`, and `Error`.
    - Configurable options:
        - `WithJSONFormat` or `WithTextFormat` to select the log format.
        - `WithOutput` to customize the output destination.
        - `WithLogLevel` to dynamically adjust log levels (`debug`, `info`, `warn`, `error`).
    - Thread-safe design utilizing `sync.Mutex` and `atomic` primitives.
- **Usage:** The logger is preconfigured and ready to use out-of-the-box.

---

### [Conv](https://github.com/KennyMacCormik/common/tree/main/conv)
- **Purpose:** Provides utilities for efficient and unsafe conversions between strings and byte slices.
- **Key Features:**
    - `StrToBytes`: Converts a string to a byte slice without copying data (uses shared memory).
    - `BytesToStr`: Converts a byte slice to a string without copying data.
    - Designed for performance-critical scenarios where immutability can be guaranteed.
- **Warnings:**
    - Modifying the underlying byte slice or string can lead to undefined behavior.
    - Use only when immutability of the underlying data is assured.

---

### [Val](https://github.com/KennyMacCormik/common/tree/main/val)
- **Purpose:** Provides a thread-safe validation mechanism using the `go-playground/validator/v10` library.
- **Key Features:**
    - Built as a singleton to ensure a single validation instance across the application.
    - Supports:
        - Struct validation using `ValidateStruct`.
        - Tag-based validation using `ValidateWithTag`.
    - Custom validation rules can be added using `RegisterValidation`.
    - Preconfigured with custom validators such as `urlprefix` (validates URLs starting with `http://` or `https://`).
- **Usage:** The validator is initialized internally and requires no setup.

---

### [Gin Factory](https://github.com/KennyMacCormik/common/tree/main/http/gin_factory)
- **Purpose:** Provides a `GinFactory` for managing middleware and handlers in a modular way.
- **Key Features:**
    - Simplifies the creation of `gin.Engine` with preconfigured middleware and route handlers.
    - Supports adding middleware with `AddMiddleware`.
    - Supports adding route handlers with `AddHandlers`.
    - Creates a new router instance with `CreateRouter` that applies all configured middleware and handlers.
    - Preconfigured with `gin.Recovery` middleware for handling panics gracefully.
- **Usage:** Ideal for creating modular and reusable HTTP routing logic in `Gin` applications.

---

## License

This repository is licensed under the [MIT License](https://opensource.org/licenses/MIT).

---

## Acknowledgments

Special thanks to the contributors and maintainers of the following libraries that this repository is based on:
- [slog](https://pkg.go.dev/log/slog) for logging.
- [go-playground/validator](https://pkg.go.dev/github.com/go-playground/validator/v10) for validation.
- [gin](https://pkg.go.dev/github.com/gin-gonic/gin) for HTTP routing.

