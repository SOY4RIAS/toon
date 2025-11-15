# TOON Go Package - Status

## Package Information

**Module**: `github.com/soy4rias/toongo`
**Status**: Complete and Production Ready
**Conformance**: 100% (323/323 tests passing)
**Spec Version**: v1.4

> **Fork Notice**: This is a community fork of the original [toon-format/toon](https://github.com/toon-format/toon) TypeScript implementation. This Go package is independently maintained and not affiliated with the original creators.

## Feature Status

### Core Features (Complete)

- **Encoder**: Full TOON encoding support
  - Objects (flat and nested)
  - Arrays (inline, list, tabular formats)
  - Automatic tabular detection for uniform object arrays
  - Deterministic output (alphabetically sorted keys)
  - Smart quoting (only when necessary)
  - Custom delimiters (comma, tab, pipe)
  - Optional length markers
  - All primitive types

- **Decoder**: Full TOON decoding support
  - All TOON formats supported
  - Strict mode validation
  - Error recovery and handling
  - Unicode support
  - All primitive types

- **Conformance**: 100% compliance with official spec
  - Encode: 138/138 passing
  - Decode: 185/185 passing
  - Total: 323/323 passing

### Optional Features (Not Implemented)

- CLI tool for command-line usage
- Performance benchmarking vs. other implementations
- Publishing to pkg.go.dev

## Test Results

```
=== Unit Tests ===
Decode tests: 5/5 passing
Encode tests: 8/8 passing
Round-trip tests: passing

=== Conformance Tests ===
Decode: 185/185 (100%)
Encode: 138/138 (100%)
Total:  323/323 (100%)
```

## Known Limitations

1. **Field Ordering**: Object keys are sorted alphabetically (Go map semantics)
   - Trade-off for deterministic output
   - Semantically correct (order doesn't affect meaning)

2. **Numeric Types**: Numbers decode as `float64` (standard Go behavior)

3. **No CLI**: Command-line tool not implemented

## Project Structure

```
/
├── *.go                     # Go implementation (100% conformance)
├── encode/                  # Encoder subdirectory
├── shared/                  # Shared utilities
├── go.mod                   # Go module definition
├── *_test.go               # Test files
├── spec-tests/              # Official spec test fixtures
├── ts-version/              # Original TypeScript (reference only)
└── Documentation/
    ├── README.md            # Package documentation
    ├── CLAUDE.md            # Development guide
    ├── CONFORMANCE.md       # Test results
    └── SPEC.md              # Spec reference
```

## Usage

```bash
# Install
go get github.com/soy4rias/toongo

# Test
go test ./...

# Conformance tests
go test -run Conformance
```

## Next Steps (Optional)

1. CLI tool for command-line usage
2. Performance benchmarking
3. Publish to pkg.go.dev
4. Additional documentation

---

Last updated: 2025-11-15
