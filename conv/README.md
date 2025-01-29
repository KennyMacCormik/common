
# Package `conv`

## Overview

The `conv` package provides efficient conversion utilities between Go `string` and `[]byte` types using the `unsafe` package. These conversions avoid memory copying by sharing underlying data, making them suitable for performance-critical scenarios. However, these utilities come with inherent risks due to Go's immutability guarantees for `string` types.

---

## Features

- **Zero-Copy Conversion**: Avoids memory duplication by sharing underlying data between `string` and `[]byte`.
- **Performance-Critical Use Cases**: Designed for scenarios where conversion overhead must be minimized.
- **Extensive Documentation**: Includes warnings and usage recommendations to mitigate potential misuse.

---

## Installation

### Dependencies

The `conv` package relies on Go's standard library and requires Go version 1.23 or newer for compatibility with the `unsafe` utilities used.

### Installation Command

To install the package, run:

```bash
go get github.com/KennyMacCormik/HerdMaster/pkg/conv
```

---

## Usage

#### `StrToBytes(s string) []byte`

Converts a `string` to a `[]byte` without copying data.

Example:

```go
package main

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/pkg/conv"
)

func main() {
	str := "example"
	bytes := conv.StrToBytes(str)
	fmt.Printf("%v\n", bytes)
}
```

#### `BytesToStr(b []byte) string`

Converts a `[]byte` to a `string` without copying data.

Example:

```go
package main

import (
	"fmt"
	"github.com/KennyMacCormik/HerdMaster/pkg/conv"
)

func main() {
	bytes := []byte("example")
	str := conv.BytesToStr(bytes)
	fmt.Printf("%s\n", str)
}
```

---

## API Documentation

### Functions

#### `func StrToBytes(s string) []byte`

- Converts a `string` to a `[]byte`.
- **Warning**: Modifying the resulting `[]byte` can lead to undefined behavior.
- Suitable for scenarios where immutability of the string can be guaranteed.

#### `func BytesToStr(b []byte) string`

- Converts a `[]byte` to a `string`.
- **Warning**: The input `[]byte` must not be modified after conversion.
- Use in contexts where the byte slice's immutability is ensured.

---

## License

This package is licensed under the [MIT License](https://opensource.org/licenses/MIT).

---

## Thanks

Special thanks to the contributors and maintainers of the Go programming language and the authors of the `unsafe` package, whose tools made this package possible.
