# TOON for Go

[![CI](https://github.com/soy4rias/toongo/actions/workflows/ci.yml/badge.svg)](https://github.com/soy4rias/toongo/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/soy4rias/toongo)](https://goreportcard.com/report/github.com/soy4rias/toongo)
[![SPEC v1.4](https://img.shields.io/badge/spec-v1.4-lightgray)](https://github.com/toon-format/spec)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

> **Go port of [toon-format/toon](https://github.com/toon-format/toon)** - A compact, human-readable serialization format designed for passing structured data to Large Language Models.

**Token-Oriented Object Notation (TOON)** is a serialization format that significantly reduces token usage when sending structured data to LLMs. This is the official Go implementation, ported from the reference TypeScript implementation.

TOON's sweet spot is **uniform arrays of objects** – multiple fields per row, same structure across items. It borrows YAML's indentation-based structure for nested objects and CSV's tabular format for uniform data rows, then optimizes both for token efficiency in LLM contexts.

## Table of Contents

- [Why TOON?](#why-toon)
- [Key Features](#key-features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Format Overview](#format-overview)
- [API Reference](#api-reference)
- [Testing](#testing)
- [Conformance](#conformance)
- [Project Status](#project-status)
- [Other Implementations](#other-implementations)
- [Contributing](#contributing)
- [License](#license)

## Why TOON?

LLM tokens still cost money – and standard JSON is verbose and token-expensive:

```json
{
  "users": [
    { "id": 1, "name": "Alice", "role": "admin" },
    { "id": 2, "name": "Bob", "role": "user" }
  ]
}
```

TOON conveys the same information with **fewer tokens**:

```
users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user
```

**Token savings:** TOON typically achieves **30–60% fewer tokens** than formatted JSON for uniform tabular data.

## Key Features

- 💸 **Token-efficient:** 30–60% fewer tokens than JSON for tabular data
- 🤿 **LLM-friendly guardrails:** Explicit lengths and fields enable validation
- 🍱 **Minimal syntax:** Removes redundant punctuation (braces, brackets, most quotes)
- 📐 **Indentation-based structure:** Like YAML, uses whitespace instead of braces
- 🧺 **Tabular arrays:** Declare keys once, stream data as rows
- ✅ **100% conformance:** Passes all official TOON v1.4 specification tests

## Installation

```bash
go get github.com/soy4rias/toongo
```

**Requirements:**
- Go 1.21 or later

## Quick Start

### Basic Encoding

```go
package main

import (
    "fmt"
    "github.com/soy4rias/toongo"
)

func main() {
    data := map[string]interface{}{
        "users": []interface{}{
            map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
            map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
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
}
```

### Basic Decoding

```go
package main

import (
    "fmt"
    "github.com/soy4rias/toongo"
)

func main() {
    input := `users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user`

    decoded, err := toon.Decode(input, nil)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", decoded)
    // Output: map[users:[map[id:1 name:Alice role:admin] map[id:2 name:Bob role:user]]]
}
```

### Custom Options

```go
// Encode with tab delimiter and length markers
encoded, err := toon.Encode(data, &toon.EncodeOptions{
    Delimiter:    "\t",
    LengthMarker: "#",
    Indent:       4,
})

// Decode with custom indentation
decoded, err := toon.Decode(input, &toon.DecodeOptions{
    Indent: 4,
    Strict: true,
})
```

## Format Overview

For complete format details, see the [TOON v1.4 Specification](https://github.com/toon-format/spec/blob/main/SPEC.md).

### Objects

Simple objects with primitive values:

```go
data := map[string]interface{}{
    "id":     123,
    "name":   "Ada",
    "active": true,
}

// Encodes to:
// id: 123
// name: Ada
// active: true
```

### Arrays

#### Primitive Arrays (Inline)

```go
data := map[string]interface{}{
    "tags": []interface{}{"admin", "ops", "dev"},
}

// Encodes to:
// tags[3]: admin,ops,dev
```

#### Arrays of Objects (Tabular)

When all objects share the same primitive fields, TOON uses an efficient tabular format:

```go
data := map[string]interface{}{
    "items": []interface{}{
        map[string]interface{}{"sku": "A1", "qty": 2, "price": 9.99},
        map[string]interface{}{"sku": "B2", "qty": 1, "price": 14.5},
    },
}

// Encodes to:
// items[2]{sku,qty,price}:
//   A1,2,9.99
//   B2,1,14.5
```

### Alternative Delimiters

TOON supports comma (default), tab, and pipe delimiters:

```go
// Tab-separated (often more token-efficient)
encoded, _ := toon.Encode(data, &toon.EncodeOptions{
    Delimiter: "\t",
})
// items[2	]{sku	qty	price}:
//   A1	2	9.99
//   B2	1	14.5

// Pipe-separated
encoded, _ := toon.Encode(data, &toon.EncodeOptions{
    Delimiter: "|",
})
// items[2|]{sku|qty|price}:
//   A1|2|9.99
//   B2|1|14.5
```

## API Reference

### Encode

```go
func Encode(value interface{}, options *EncodeOptions) (string, error)
```

Converts any Go value to TOON format.

**Parameters:**
- `value` - Any JSON-serializable Go value (maps, slices, primitives, structs)
- `options` - Optional encoding configuration:
  - `Indent` (int) - Number of spaces per indentation level (default: 2)
  - `Delimiter` (string) - Array delimiter: `","`, `"\t"`, or `"|"` (default: `","`)
  - `LengthMarker` (string) - Optional marker to prefix array lengths: `"#"` or `""` (default: `""`)

**Returns:**
- TOON-formatted string with no trailing newline
- Error if encoding fails

**Example:**

```go
data := map[string]interface{}{
    "items": []interface{}{
        map[string]interface{}{"id": 1, "name": "Widget"},
        map[string]interface{}{"id": 2, "name": "Gadget"},
    },
}

encoded, err := toon.Encode(data, &toon.EncodeOptions{
    Indent:       2,
    Delimiter:    ",",
    LengthMarker: "#",
})
```

### Decode

```go
func Decode(input string, options *DecodeOptions) (interface{}, error)
```

Converts a TOON-formatted string back to Go values.

**Parameters:**
- `input` - A TOON-formatted string to parse
- `options` - Optional decoding configuration:
  - `Indent` (int) - Expected number of spaces per indentation level (default: 2)
  - `Strict` (bool) - Enable strict validation (default: true)

**Returns:**
- Go value (map, slice, or primitive) representing the parsed data
- Error if parsing fails

**Example:**

```go
input := `items[2]{id,name}:
  1,Widget
  2,Gadget`

decoded, err := toon.Decode(input, &toon.DecodeOptions{
    Indent: 2,
    Strict: true,
})
// Result: map[items:[map[id:1 name:Widget] map[id:2 name:Gadget]]]
```

### Type Conversions

Some non-JSON types are automatically normalized:

| Input Type | Output |
|---|---|
| Finite numbers | Decimal form (no scientific notation) |
| `NaN`, `±Infinity` | `null` |
| `time.Time` | ISO 8601 string in quotes |
| Structs | Converted to objects (map) |
| Pointers | Dereferenced |

## Testing

Run the test suite:

```bash
# Run all tests
go test -v ./...

# Run only conformance tests
go test -v -run TestConformance

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

**Test Coverage:**
- ✅ Basic unit tests: 100% passing
- ✅ Encode/decode round-trip tests: 100% passing
- ✅ Official conformance tests: 100% passing (323/323)

## Conformance

This implementation achieves **100% conformance** with the [TOON v1.4 Specification](https://github.com/toon-format/spec):

- **Decoder tests:** 185/185 passing (100%)
- **Encoder tests:** 138/138 passing (100%)
- **Overall:** 323/323 passing (100%)

See [CONFORMANCE.md](./CONFORMANCE.md) for detailed test results.

## Project Status

This is a **production-ready** Go port of the [toon-format/toon](https://github.com/toon-format/toon) reference implementation.

**What's Complete:**
- ✅ Full encoder implementation (tabular, inline, list formats)
- ✅ Full decoder implementation with strict validation
- ✅ All delimiter options (comma, tab, pipe)
- ✅ Length markers
- ✅ Proper quoting and escaping
- ✅ 100% conformance with TOON v1.4 spec
- ✅ Comprehensive test suite
- ✅ CI/CD pipeline

**Not Yet Implemented:**
- ⏳ CLI tool (planned)
- ⏳ Benchmarks vs. JSON encoding
- ⏳ Examples directory

## Other Implementations

This is the **official Go implementation**. For the reference TypeScript implementation and other language ports, see:

### Official Implementations

- **TypeScript/JavaScript:** [toon-format/toon](https://github.com/toon-format/toon) *(reference implementation)*
- **Python:** [toon-format/toon-python](https://github.com/toon-format/toon-python) *(in development)*
- **Rust:** [toon-format/toon-rust](https://github.com/toon-format/toon-rust) *(in development)*

### Community Implementations

- **.NET:** [ToonSharp](https://github.com/0xZunia/ToonSharp)
- **C++:** [ctoon](https://github.com/mohammadraziei/ctoon)
- **Clojure:** [toon](https://github.com/vadelabs/toon)
- **Crystal:** [toon-crystal](https://github.com/mamantoha/toon-crystal)
- **Dart:** [toon](https://github.com/wisamidris77/toon)
- **Elixir:** [toon_ex](https://github.com/kentaro/toon_ex)
- **Gleam:** [toon_codec](https://github.com/axelbellec/toon_codec)
- **Go (alternative):** [gotoon](https://github.com/alpkeskin/gotoon)
- **Java:** [JToon](https://github.com/felipestanzani/JToon)
- **Lua/Neovim:** [toon.nvim](https://github.com/thalesgelinger/toon.nvim)
- **OCaml:** [ocaml-toon](https://github.com/davesnx/ocaml-toon)
- **PHP:** [toon-php](https://github.com/HelgeSverre/toon-php)
- **Python (alternative):** [python-toon](https://github.com/xaviviro/python-toon)
- **Ruby:** [toon-ruby](https://github.com/andrepcg/toon-ruby)
- **Swift:** [TOONEncoder](https://github.com/mattt/TOONEncoder)

## Contributing

Contributions are welcome! This is a port of the official TOON specification, so please ensure changes align with the [TOON v1.4 spec](https://github.com/toon-format/spec).

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`go test -v ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Related Resources

- **[TOON Specification v1.4](https://github.com/toon-format/spec/blob/main/SPEC.md)** - Complete technical specification
- **[Conformance Tests](https://github.com/toon-format/spec/tree/main/tests)** - Language-agnostic test fixtures
- **[Original TypeScript Implementation](https://github.com/toon-format/toon)** - Reference implementation
- **[Format Tokenization Playground](https://www.curiouslychase.com/playground/format-tokenization-exploration)** - Interactive comparison tool

## License

[MIT](./LICENSE) License © 2025-PRESENT [Johann Schopplich](https://github.com/johannschopplich) (original implementation)

This Go port is maintained by [SOY4RIAS](https://github.com/soy4rias).

---

**Note:** This is a **Go port** of the original [toon-format/toon](https://github.com/toon-format/toon) TypeScript implementation. All credit for the TOON format design and specification goes to the original authors.
