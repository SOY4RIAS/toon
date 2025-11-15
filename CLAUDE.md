# TOON Go Package - Development Guide

This document provides context for developing and maintaining the TOON Go package.

## Project Overview

**Repository**: `github.com/soy4rias/toongo`
**Type**: Community fork of [toon-format/toon](https://github.com/toon-format/toon)
**Purpose**: Go implementation of TOON (Token-Oriented Object Notation) format
**Spec Version**: v1.4 ([specification](https://github.com/toon-format/spec))
**Conformance**: 100% (323/323 tests passing)

> **Important**: This is a community fork and is not affiliated with the original TOON creators. The original TypeScript implementation was created by Johann Schopplich. This repository provides a Go implementation as a workaround for Go-based applications.

## Repository Structure

```
/ (root) - Go Package
├── go.mod                    # Go module: github.com/soy4rias/toongo
├── toon.go                   # Main API (Encode/Decode functions)
├── types.go                  # Type definitions (EncodeOptions, DecodeOptions, etc.)
├── constants.go              # Constants and delimiters
├── shared/                   # Shared utilities
│   ├── string_utils.go       # String manipulation (escape, unescape, quoting)
│   └── literal_utils.go      # Literal parsing (bool, number, null)
├── encode/                   # Encoder implementation
│   ├── encoders.go           # Main encoding logic (objects, arrays, tabular)
│   ├── normalize.go          # Value normalization
│   ├── primitives.go         # Primitive encoding & header formatting
│   ├── writer.go             # Output writer with indentation
│   └── validation.go         # Quoting and key validation
├── decode_scanner.go         # Scanner - line tokenization, indentation
├── decode_parser.go          # Parser - header parsing, structural analysis
├── decode_decoders.go        # Decoders - value decoding logic
├── decode_validation.go      # Validation - strict mode validation
├── *_test.go                 # Test files
├── conformance_test.go       # Conformance tests (100% passing)
├── spec-tests/               # Official spec test fixtures
├── ts-version/               # Original TypeScript implementation (reference only)
└── Documentation
    ├── README.md             # Go package documentation
    ├── CONFORMANCE.md        # Conformance test results
    ├── SPEC.md               # TOON specification reference
    └── STATUS.md             # Development status
```

## Development Guidelines

### Building and Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run conformance tests only
go test -run Conformance

# Check for race conditions
go test -race ./...
```

### Code Standards

1. **Error Handling**: Return errors explicitly, no panics except in internal recovery
2. **Type Safety**: Use `interface{}` for JSON-compatible values
3. **String Building**: Use `strings.Builder` for efficient concatenation
4. **Unicode**: Handle runes properly for character operations
5. **Determinism**: Object keys are sorted alphabetically for consistent output

### Key Implementation Details

#### Encoder (`encode/`)
- **Tabular Detection**: Automatically identifies uniform object arrays
- **Field Ordering**: Alphabetical sorting via `sort.Strings()` for deterministic output
- **Smart Quoting**: Only quotes strings when necessary (see `validation.go`)
- **Delimiter Support**: Comma, tab, pipe delimiters
- **Normalization**: Handles Go-specific types (BigInt as int64, time.Time, etc.)

#### Decoder (`decode_*.go`)
- **Scanner**: Tokenizes lines, preserves indentation for validation
- **Parser**: Extracts headers, handles delimited values
- **Decoders**: Recursive decoding for nested structures
- **Validation**: Strict mode checks for spec compliance

### Technical Decisions

1. **Alphabetical Field Order**
   - Go maps have non-deterministic iteration order
   - Solution: Sort keys before encoding
   - Trade-off: Consistent output vs. insertion order preservation

2. **Semantic Test Comparison**
   - Conformance tests compare decoded values, not raw strings
   - Handles field order differences gracefully

3. **Module Path**
   - `github.com/soy4rias/toongo` - clearly identifies this as a separate fork
   - Not using `github.com/toon-format/*` to avoid confusion with official repos

## What This Package IS

- A Go implementation of the TOON format specification
- A tool for encoding/decoding TOON in Go applications
- 100% spec-compliant according to official conformance tests
- A community contribution for Go developers

## What This Package IS NOT

- An official TOON implementation
- Affiliated with or endorsed by the original creators
- A replacement for the TypeScript implementation
- The authoritative TOON reference (that's [toon-format/spec](https://github.com/toon-format/spec))

## References

- [TOON Specification v1.4](https://github.com/toon-format/spec/blob/main/SPEC.md) - The authoritative specification
- [Official Conformance Tests](https://github.com/toon-format/spec/tree/main/tests) - Language-agnostic test fixtures
- [Original TypeScript Implementation](https://github.com/toon-format/toon) - Reference implementation by original creators
- [TypeScript Reference in this repo](./ts-version/) - Preserved copy for reference

## Maintenance Notes

- Follow the official TOON spec v1.4 for any feature additions
- Run conformance tests after any changes: `go test -run Conformance`
- Keep README.md updated with any API changes
- Document breaking changes clearly
- Preserve the TypeScript reference in `ts-version/` but don't modify it

## Contributing

See [README.md](./README.md#contributing) for contribution guidelines.

---

Last updated: 2025-11-15
