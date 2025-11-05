# toon-go

Go implementation of the TOON (Terse Object Oriented Notation) decoder.

## Overview

TOON is a human-readable data serialization format designed for clarity and compactness. This package provides a complete decoder implementation for parsing TOON-formatted strings into Go data structures.

## Installation

```bash
go get github.com/toon-format/toon-go
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/toon-format/toon-go"
)

func main() {
    input := `
name: John Doe
age: 30
active: true
tags[3]: web, backend, api
`

    result, err := toon.Decode(input, nil)
    if err != nil {
        log.Fatal(err)
    }

    obj := result.(map[string]interface{})
    fmt.Printf("Name: %s\n", obj["name"])
    fmt.Printf("Age: %.0f\n", obj["age"])
    fmt.Printf("Active: %t\n", obj["active"])
    fmt.Printf("Tags: %v\n", obj["tags"])
}
```

### With Options

```go
indent := 4
strict := false
options := &toon.DecodeOptions{
    Indent: &indent,
    Strict: &strict,
}

result, err := toon.Decode(input, options)
```

## Options

- **Indent** (`*int`): Number of spaces per indentation level. Default: `2`
- **Strict** (`*bool`): Enforce strict validation of array lengths and tabular row counts. Default: `true`

## Package Structure

The implementation is organized into the following files:

- **toon.go** - Main entry point with the `Decode()` function
- **types.go** - Type definitions, constants, and options
- **string.go** - String utility functions (escape, unescape, find closing quote, etc.)
- **literal.go** - Literal checking functions (isNumeric, isBoolean, etc.)
- **scanner.go** - `ToParsedLines` function and `LineCursor` type
- **parser.go** - Parsing functions (ParseArrayHeaderLine, ParsePrimitiveToken, etc.)
- **decoder.go** - Main decoder logic (DecodeValueFromLines and all decode functions)
- **validation.go** - Validation functions

## Features

- ✅ Object decoding
- ✅ Array decoding (inline, list, tabular)
- ✅ Primitive types (string, number, boolean, null)
- ✅ Quoted and unquoted keys
- ✅ Escape sequences
- ✅ Multiple delimiters (comma, tab, pipe)
- ✅ Strict mode validation
- ✅ Configurable indentation

## Testing

Run the tests:

```bash
go test -v
```

## License

This implementation follows the TOON specification and is part of the TOON format ecosystem.

## Related

- [TOON Specification](https://github.com/toon-format/toon)
- [TypeScript Implementation](https://github.com/toon-format/toon/tree/main/packages/toon)
