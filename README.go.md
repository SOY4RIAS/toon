# TOON - Go Implementation

[![Go Reference](https://pkg.go.dev/badge/github.com/SOY4RIAS/toon.svg)](https://pkg.go.dev/github.com/SOY4RIAS/toon)
[![Go Report Card](https://goreportcard.com/badge/github.com/SOY4RIAS/toon)](https://goreportcard.com/report/github.com/SOY4RIAS/toon)
[![SPEC v1.4](https://img.shields.io/badge/spec-v1.4-lightgray)](https://github.com/toon-format/spec)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

**Token-Oriented Object Notation (TOON)** is a compact, human-readable serialization format designed for passing structured data to Large Language Models with significantly reduced token usage.

This is the **official Go implementation** of the TOON format, providing 100% compatibility with the [TOON specification v1.4](https://github.com/toon-format/spec).

## Features

- 💸 **Token-efficient:** typically 30–60% fewer tokens than JSON
- 🤿 **LLM-friendly guardrails:** explicit lengths and fields enable validation
- 🍱 **Minimal syntax:** removes redundant punctuation (braces, brackets, most quotes)
- 📐 **Indentation-based structure:** like YAML, uses whitespace instead of braces
- 🧺 **Tabular arrays:** declare keys once, stream data as rows
- ⚡ **High performance:** native Go implementation with zero dependencies
- 🔒 **Type safe:** strong typing with comprehensive error handling

## Installation

```bash
go get github.com/SOY4RIAS/toon
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/SOY4RIAS/toon"
)

func main() {
    // Encode Go values to TOON
    data := map[string]interface{}{
        "users": []map[string]interface{}{
            {"id": 1, "name": "Alice", "role": "admin"},
            {"id": 2, "name": "Bob", "role": "user"},
        },
    }

    encoded, err := toon.Encode(data, nil)
    if err != nil {
        panic(err)
    }

    fmt.Println(encoded)
    // Output:
    // users[2]{id,name,role}:
    //   1,Alice,admin
    //   2,Bob,user

    // Decode TOON back to Go values
    decoded, err := toon.Decode(encoded, nil)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", decoded)
}
```

## API Reference

### Encoding

```go
func Encode(value interface{}, opts *EncodeOptions) (string, error)
```

Converts Go values to TOON format.

**Options:**
```go
type EncodeOptions struct {
    Indent       int        // Indentation spaces (default: 2)
    Delimiter    Delimiter  // ',', '\t', '|' (default: ',')
    LengthMarker bool       // Use '#' prefix for array lengths (default: false)
}
```

**Example:**
```go
data := map[string]interface{}{
    "items": []map[string]interface{}{
        {"sku": "A1", "qty": 2, "price": 9.99},
        {"sku": "B2", "qty": 1, "price": 14.5},
    },
}

// Default options (comma delimiter, 2-space indent)
toon.Encode(data, nil)

// Tab-delimited with length markers
toon.Encode(data, &toon.EncodeOptions{
    Delimiter:    toon.DelimiterTab,
    LengthMarker: true,
})
```

### Decoding

```go
func Decode(input string, opts *DecodeOptions) (interface{}, error)
```

Parses TOON format back to Go values.

**Options:**
```go
type DecodeOptions struct {
    Indent int  // Expected indentation spaces (default: 2)
    Strict bool // Enable strict validation (default: true)
}
```

**Example:**
```go
input := `users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user`

decoded, err := toon.Decode(input, nil)
if err != nil {
    panic(err)
}
```

## CLI Tool

Install the CLI:
```bash
go install github.com/SOY4RIAS/toon/cmd/toon@latest
```

**Usage:**
```bash
# Encode JSON to TOON
toon encode input.json -o output.toon

# Decode TOON to JSON
toon decode data.toon -o output.json

# From stdin
echo '{"name": "Alice", "age": 30}' | toon encode

# With options
toon encode data.json --delimiter tab --length-marker
```

## Format Overview

### Objects

```go
data := map[string]interface{}{
    "id":     123,
    "name":   "Alice",
    "active": true,
}
```

**TOON output:**
```
id: 123
name: Alice
active: true
```

### Arrays of Objects (Tabular Format)

```go
data := map[string]interface{}{
    "users": []map[string]interface{}{
        {"id": 1, "name": "Alice", "role": "admin"},
        {"id": 2, "name": "Bob", "role": "user"},
    },
}
```

**TOON output:**
```
users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user
```

### Primitive Arrays

```go
data := map[string]interface{}{
    "tags": []string{"admin", "ops", "dev"},
}
```

**TOON output:**
```
tags[3]: admin,ops,dev
```

## Specification Compliance

This implementation is fully compliant with [TOON Specification v1.4](https://github.com/toon-format/spec/blob/main/SPEC.md) and passes all official conformance tests.

## Performance

The Go implementation is optimized for performance:
- Zero allocations for encoding primitives
- Efficient string building with pre-allocation
- Streaming parser for decoding
- Minimal memory overhead

Benchmarks coming soon.

## Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run benchmarks
go test -bench=. ./...
```

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code follows Go conventions
- Changes maintain spec compliance
- Documentation is updated

## License

[MIT](./LICENSE) License © 2025-PRESENT [Johann Schopplich](https://github.com/johannschopplich)

Go implementation by [SOY4RIAS](https://github.com/SOY4RIAS)

## Related Projects

- [TOON Specification](https://github.com/toon-format/spec) - Official format specification
- [TOON TypeScript/JavaScript](https://github.com/toon-format/toon) - Reference implementation
- [Other implementations](https://github.com/toon-format/toon#other-implementations) - Community ports

## Links

- [Documentation](https://pkg.go.dev/github.com/SOY4RIAS/toon)
- [TOON Format Website](https://toonformat.dev)
- [Specification](https://github.com/toon-format/spec)
- [Issue Tracker](https://github.com/SOY4RIAS/toon/issues)
