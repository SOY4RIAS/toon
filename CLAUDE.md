# TOON Go Implementation

This document provides context for working with the TOON Go implementation.

## Project Overview

**Goal**: Provide a production-ready Go implementation of the TOON (Token-Oriented Object Notation) format.

**Repository**: `github.com/soy4rias/toongo`
**Spec Version**: v1.4 ([spec repository](https://github.com/toon-format/spec))
**Status**: ✅ **Production Ready** - 100% Conformance Achieved
**Module**: `github.com/soy4rias/toongo`

## What is TOON?

TOON is a compact, human-readable serialization format designed for passing structured data to Large Language Models with significantly reduced token usage. This is a **Go port** of the original [toon-format/toon](https://github.com/toon-format/toon) TypeScript implementation.

## Project Structure

```
/
├── *.go                     # Go source files
├── encode/                  # Encoder implementation
│   ├── normalize.go         # Value normalization
│   ├── primitives.go        # Primitive encoding
│   ├── writer.go            # Output writer
│   ├── encoders.go          # Main encoding logic
│   └── validation.go        # Quoting and validation
├── shared/                  # Shared utilities
│   ├── string_utils.go      # String manipulation
│   └── literal_utils.go     # Literal parsing
├── decode_scanner.go        # Scanner - line tokenization
├── decode_parser.go         # Parser - header parsing
├── decode_decoders.go       # Decoders - value decoding
├── decode_validation.go     # Validation - strict mode
├── toon.go                  # Main API (Encode/Decode)
├── types.go                 # Type definitions
├── constants.go             # Constants and delimiters
├── *_test.go                # Test files
├── conformance_test.go      # Conformance tests (100% passing)
├── go.mod                   # Go module definition
├── spec-tests/              # Official TOON spec test fixtures
├── README.md                # User-facing documentation
├── CONFORMANCE.md           # Conformance test results
├── SPEC.md                  # Spec reference
└── CLAUDE.md                # This file (project context)
```

## Implementation Status

### ✅ Complete (100%)

**Core Implementation:**
- ✅ Constants and types (constants.go, types.go)
- ✅ Shared utilities (shared/string_utils.go, shared/literal_utils.go)
- ✅ Complete Decoder (scanner → parser → decoders → validation)
- ✅ Complete Encoder (normalize → primitives → writer → encoders → validation)
- ✅ Main API (toon.go - Encode/Decode functions)

**Testing:**
- ✅ Basic unit tests (13/13 passing - 100%)
- ✅ Round-trip tests (encode → decode → verify)
- ✅ Conformance tests (323/323 passing - **100% compliance**)
  - Decoder: 185/185 passing (100%)
  - Encoder: 138/138 passing (100%)

**Documentation:**
- ✅ Comprehensive README.md with Go-specific examples
- ✅ CONFORMANCE.md with detailed test results
- ✅ SPEC.md reference
- ✅ Inline code documentation

**Project Infrastructure:**
- ✅ CI/CD pipeline (.github/workflows/ci.yml)
- ✅ Go module configuration (go.mod)
- ✅ Dev container setup (.devcontainer/devcontainer.json)
- ✅ VS Code settings (.vscode/settings.json)

### ⏳ Planned Future Work

- ⏳ CLI tool (optional - TypeScript CLI exists in original repo)
- ⏳ Performance benchmarks (vs JSON, vs TypeScript)
- ⏳ Examples directory
- ⏳ pkg.go.dev publication

## Key Features

- **100% Conformance**: Passes all 323 official TOON v1.4 spec tests
- **Production Ready**: Full encoder and decoder with comprehensive error handling
- **Strict Validation**: Optional strict mode for validating TOON structure
- **Multiple Delimiters**: Support for comma, tab, and pipe delimiters
- **Length Markers**: Optional `#` prefix for array lengths
- **Type Safe**: Proper Go type handling with struct support
- **Well Tested**: Comprehensive test suite with high coverage

## Technical Decisions

### Field Ordering (Alphabetical Sorting)

**Problem**: Go's `map` type has non-deterministic iteration order.

**Solution**: All object keys are sorted alphabetically before encoding.
- Uses `sort.Strings()` from stdlib (zero dependencies)
- Provides deterministic, reproducible output
- Passes all conformance tests with semantic comparison

**Trade-offs**:
- ✅ Deterministic output
- ✅ Language-agnostic (alphabetical order)
- ✅ Simple implementation
- ⚠️ Field order differs from insertion order (but TOON semantics preserved)

### Tab Validation Fix

**Issue**: Leading tabs/spaces were being stripped during scanning, breaking strict mode validation.

**Solution**: Modified `decode_scanner.go` to preserve leading whitespace for validation while still correctly parsing indentation.

## Usage Examples

### Basic Encoding

```go
import "github.com/soy4rias/toongo"

data := map[string]interface{}{
    "users": []interface{}{
        map[string]interface{}{"id": 1, "name": "Alice", "role": "admin"},
        map[string]interface{}{"id": 2, "name": "Bob", "role": "user"},
    },
}

encoded, err := toon.Encode(data, nil)
// Output:
// users[2]{id,name,role}:
//   1,Alice,admin
//   2,Bob,user
```

### Basic Decoding

```go
input := `users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user`

decoded, err := toon.Decode(input, nil)
// Result: map[users:[map[id:1 name:Alice role:admin] map[id:2 name:Bob role:user]]]
```

### Custom Options

```go
// Tab-separated with length markers
encoded, err := toon.Encode(data, &toon.EncodeOptions{
    Delimiter:    "\t",
    LengthMarker: "#",
    Indent:       4,
})

// Lenient decoding
decoded, err := toon.Decode(input, &toon.DecodeOptions{
    Strict: false,
})
```

## Testing

```bash
# Run all tests
go test -v ./...

# Run only conformance tests
go test -v -run TestConformance

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

## Conformance Results

**Overall**: 323/323 tests passing (**100% compliance**)

**Breakdown**:
- Decoder tests: 185/185 (100%)
- Encoder tests: 138/138 (100%)

See [CONFORMANCE.md](./CONFORMANCE.md) for detailed results.

## Migration from TypeScript

This is a **complete port** of the [toon-format/toon](https://github.com/toon-format/toon) TypeScript implementation. Key differences:

| Aspect | TypeScript | Go |
|--------|-----------|-----|
| Type system | Union types, `unknown` | `interface{}`, type assertions |
| Error handling | Exceptions | Explicit error returns |
| Field order | Insertion order preserved | Alphabetically sorted (deterministic) |
| JSON handling | Native JSON types | `encoding/json` package |
| String building | Template literals | `strings.Builder` |
| Unicode | Native | Rune-based iteration |

## References

- **[TOON Specification v1.4](https://github.com/toon-format/spec/blob/main/SPEC.md)** - Complete technical specification
- **[Original TypeScript Implementation](https://github.com/toon-format/toon)** - Reference implementation
- **[Conformance Tests](https://github.com/toon-format/spec/tree/main/tests)** - Official test fixtures
- **[Other Implementations](https://github.com/toon-format/toon#other-implementations)** - Community ports

## Contributing

This is a port of the official TOON specification. When contributing:

1. **Follow the spec**: Ensure changes align with [TOON v1.4](https://github.com/toon-format/spec)
2. **Run tests**: All conformance tests must pass (`go test -v -run TestConformance`)
3. **Maintain compatibility**: Output should match the TypeScript implementation semantically
4. **Document changes**: Update README.md and CONFORMANCE.md as needed

## Version History

- **v1.0.0** (2025-11-18) - Production release with 100% conformance
  - Complete encoder and decoder implementation
  - All 323 conformance tests passing
  - Comprehensive documentation
  - CI/CD pipeline
  - Removed TypeScript reference implementation (now standalone Go project)

- **v0.9.0** (2025-11-07) - Beta release
  - 100% conformance achieved (was 97%)
  - Field ordering implemented (alphabetical sorting)
  - Tab validation fixed
  - Project reorganized (TypeScript moved to ts-version/)

- **v0.1.0** (2025-11-06) - Initial implementation
  - Core encoder and decoder
  - Basic test coverage

## Acknowledgments

All credit for the TOON format design and specification goes to:
- **[Johann Schopplich](https://github.com/johannschopplich)** - Original TOON creator and TypeScript implementation
- **[toon-format/spec](https://github.com/toon-format/spec)** contributors - Specification development

This Go port is maintained by [SOY4RIAS](https://github.com/soy4rias).

---

**Note:** This is a **Go port** of the original [toon-format/toon](https://github.com/toon-format/toon) TypeScript implementation.

Last updated: 2025-11-18
