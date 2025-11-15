# TOON Go - Token-Oriented Object Notation

[![Go Reference](https://pkg.go.dev/badge/github.com/soy4rias/toongo.svg)](https://pkg.go.dev/github.com/soy4rias/toongo)
[![SPEC v1.4](https://img.shields.io/badge/spec-v1.4-lightgray)](https://github.com/toon-format/spec)
[![Conformance](https://img.shields.io/badge/conformance-100%25-brightgreen)](./CONFORMANCE.md)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

**Go implementation of the TOON (Token-Oriented Object Notation) format** - a compact, human-readable serialization format designed for passing structured data to Large Language Models with reduced token usage.

> **Fork Notice**: This is a community fork of the original [toon-format/toon](https://github.com/toon-format/toon) TypeScript implementation. This repository focuses on providing a Go implementation of the TOON specification and is not affiliated with or maintained by the original creators.

## Overview

TOON (Token-Oriented Object Notation) is a serialization format optimized for LLM contexts. This Go package provides a complete encoder and decoder implementation that achieves **100% conformance** with the [official TOON specification v1.4](https://github.com/toon-format/spec).

### Why TOON?

TOON reduces token usage by 30-60% compared to JSON for uniform tabular data, making it cost-effective for LLM applications:

**JSON** (verbose):
```json
{
  "users": [
    { "id": 1, "name": "Alice", "role": "admin" },
    { "id": 2, "name": "Bob", "role": "user" }
  ]
}
```

**TOON** (compact):
```
users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user
```

## Installation

```bash
go get github.com/soy4rias/toongo
```

## Quick Start

```go
package main

import (
    "fmt"
    toon "github.com/soy4rias/toongo"
)

func main() {
    // Encode Go data to TOON format
    data := map[string]interface{}{
        "users": []interface{}{
            map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
            map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
        },
    }

    encoded := toon.Encode(data, toon.EncodeOptions{})
    fmt.Println(encoded)
    // Output:
    // users[2]{id,name,role}:
    //   1,Alice,admin
    //   2,Bob,user

    // Decode TOON back to Go data
    toonString := `users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user`

    decoded, err := toon.Decode(toonString, toon.DecodeOptions{})
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", decoded)
}
```

## API Reference

### Encode

```go
func Encode(value interface{}, options EncodeOptions) string
```

Converts any Go value to TOON format.

**EncodeOptions:**
- `Indent` (int): Number of spaces per indentation level (default: 2)
- `Delimiter` (string): Array delimiter - `","`, `"\t"`, or `"|"` (default: `","`)
- `LengthMarker` (string): Optional marker to prefix array lengths, e.g., `"#"` (default: `""`)

**Examples:**

```go
// Basic object
data := map[string]interface{}{
    "id":     123,
    "name":   "Ada",
    "active": true,
}
toon.Encode(data, toon.EncodeOptions{})
// Output:
// active: true
// id: 123
// name: Ada

// Tabular array (uniform objects)
items := map[string]interface{}{
    "items": []interface{}{
        map[string]interface{}{"sku": "A1", "qty": 2, "price": 9.99},
        map[string]interface{}{"sku": "B2", "qty": 1, "price": 14.5},
    },
}
toon.Encode(items, toon.EncodeOptions{})
// Output:
// items[2]{price,qty,sku}:
//   9.99,2,A1
//   14.5,1,B2

// With tab delimiter
toon.Encode(items, toon.EncodeOptions{Delimiter: "\t"})
// Output (tabs shown as spaces):
// items[2	]{price	qty	sku}:
//   9.99	2	A1
//   14.5	1	B2

// With length marker
toon.Encode(items, toon.EncodeOptions{LengthMarker: "#"})
// Output:
// items[#2]{price,qty,sku}:
//   9.99,2,A1
//   14.5,1,B2
```

### Decode

```go
func Decode(input string, options DecodeOptions) (interface{}, error)
```

Converts a TOON-formatted string back to Go values.

**DecodeOptions:**
- `Indent` (int): Expected number of spaces per indentation level (default: 2)
- `Strict` (bool): Enable strict validation (default: true)

**Examples:**

```go
toonString := `items[2]{sku,qty,price}:
  A1,2,9.99
  B2,1,14.5`

data, err := toon.Decode(toonString, toon.DecodeOptions{})
if err != nil {
    log.Fatal(err)
}
// Result: map[string]interface{}{
//   "items": []interface{}{
//     map[string]interface{}{"sku": "A1", "qty": 2.0, "price": 9.99},
//     map[string]interface{}{"sku": "B2", "qty": 1.0, "price": 14.5},
//   },
// }

// Lenient decoding (skip strict validation)
data, err := toon.Decode(toonString, toon.DecodeOptions{Strict: false})
```

**Strict Mode Validation:**
- Invalid escape sequences are rejected
- Syntax errors throw errors
- Array length mismatches are caught
- Delimiter mismatches are detected
- Indentation must be consistent

## Features

### Supported Data Types

- **Primitives**: strings, numbers (int, float), booleans, null
- **Objects**: nested objects, quoted keys, special characters in keys
- **Arrays**:
  - Inline format for primitive arrays
  - List format for mixed/non-uniform arrays
  - Tabular format for uniform object arrays
- **Special handling**: empty arrays, empty objects, unicode, emoji

### Encoding Features

- **Automatic tabular detection**: Uniform object arrays are encoded efficiently
- **Deterministic output**: Object keys are alphabetically sorted
- **Smart quoting**: Strings are only quoted when necessary
- **Delimiter options**: Comma, tab, or pipe delimiters
- **Optional length markers**: Add `#` prefix for emphasis

### Decoding Features

- **Strict validation**: Comprehensive error checking
- **All formats supported**: Inline, list, and tabular arrays
- **Error recovery**: Panics are converted to errors
- **Unicode support**: Full UTF-8 compatibility

## Conformance

This implementation achieves **100% conformance** (323/323 tests passing) with the official TOON specification v1.4.

See [CONFORMANCE.md](./CONFORMANCE.md) for detailed test results.

### Test Results

- **Decode conformance**: 185/185 passing (100%)
- **Encode conformance**: 138/138 passing (100%)
- **Overall conformance**: 323/323 passing (100%)

## Project Structure

```
/
├── toon.go                  # Main API (Encode/Decode functions)
├── types.go                 # Type definitions
├── constants.go             # Constants and delimiters
├── shared/                  # Shared utilities
│   ├── string_utils.go      # String manipulation
│   └── literal_utils.go     # Literal parsing
├── encode/                  # Encoder implementation
│   ├── encoders.go          # Main encoding logic
│   ├── normalize.go         # Value normalization
│   ├── primitives.go        # Primitive encoding
│   ├── writer.go            # Output writer
│   └── validation.go        # Quoting validation
├── decode_*.go              # Decoder implementation
│   ├── decode_scanner.go    # Line tokenization
│   ├── decode_parser.go     # Header parsing
│   ├── decode_decoders.go   # Value decoding
│   └── decode_validation.go # Strict validation
├── *_test.go                # Test files
└── ts-version/              # Original TypeScript implementation (reference)
```

## Limitations

- **Field ordering**: Object keys are sorted alphabetically (not insertion order) due to Go's map semantics. This is semantically correct but may differ from JSON source order.
- **Numeric types**: All numbers decode as `float64` (standard Go JSON behavior).
- **Date handling**: Dates should be passed as ISO 8601 strings.

## Related Projects

- **Original TypeScript Implementation**: [toon-format/toon](https://github.com/toon-format/toon) - The reference implementation this fork is based on
- **TOON Specification**: [toon-format/spec](https://github.com/toon-format/spec) - Official TOON format specification v1.4
- **Conformance Tests**: [toon-format/spec/tests](https://github.com/toon-format/spec/tree/main/tests) - Language-agnostic test fixtures

### Other Go Implementations

- [gotoon](https://github.com/alpkeskin/gotoon) - Community Go implementation

### Other Language Implementations

See the [original repository](https://github.com/toon-format/toon#other-implementations) for a complete list of community implementations in various languages.

## Contributing

This is a fork maintained independently from the original TOON project. Contributions are welcome:

1. Fork this repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

## License

[MIT](./LICENSE) License

---

**Disclaimer**: This Go implementation is a community fork and is not affiliated with, endorsed by, or maintained by the original TOON format creators. The original TOON format and TypeScript implementation are created by [Johann Schopplich](https://github.com/johannschopplich). This fork provides only a Go workaround for using TOON in Go-based applications.
