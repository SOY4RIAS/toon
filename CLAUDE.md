# TOON Go Migration - Progress Tracker

This document tracks the progress of migrating the TOON (Token-Oriented Object Notation) library from TypeScript to Go.

## Project Overview

**Goal**: Create a Go implementation of the TOON format parser/encoder that matches the TypeScript reference implementation.

**Repository**: `toon-format/toon`
**Branch**: `feature/go-port`
**Spec Version**: v1.4 ([spec repository](https://github.com/toon-format/spec))
**Status**: ✅ Core Implementation Complete (Encoder + Decoder)

## TypeScript Source Structure

```
packages/toon/src/
├── index.ts              # Main entry point (encode/decode)
├── types.ts              # Type definitions
├── constants.ts          # Constants and delimiters
├── shared/               # Shared utilities
│   ├── string-utils.ts   # String manipulation
│   ├── literal-utils.ts  # Literal parsing
│   └── validation.ts     # Validation helpers
├── encode/               # Encoder implementation
│   ├── primitives.ts     # Primitive encoding
│   ├── writer.ts         # Output writer
│   ├── normalize.ts      # Value normalization
│   └── encoders.ts       # Main encoding logic
└── decode/               # Decoder implementation
    ├── decoders.ts       # Main decoding logic
    ├── scanner.ts        # Line scanning/tokenization
    ├── parser.ts         # Parsing logic
    └── validation.ts     # Decode validation
```

## Go Project Structure

```
go/
├── go.mod                    # Go module definition
├── toon.go                   # Main entry point (Encode/Decode)
├── types.go                  # Type definitions
├── constants.go              # Constants and delimiters
├── shared/                   # Shared utilities
│   ├── string_utils.go       # String manipulation (escape, unescape, quote finding)
│   └── literal_utils.go      # Literal parsing (bool, number, null)
├── encode/                   # Encoder implementation
│   ├── normalize.go          # Value normalization (Date, BigInt, structs, etc)
│   ├── primitives.go         # Primitive encoding & header formatting
│   ├── writer.go             # Output writer with indentation
│   ├── encoders.go           # Main encoding logic (objects, arrays, tabular)
│   └── validation.go         # Quoting and key validation
├── decode/                   # Decoder implementation (in root go/ for now)
│   ├── decoders.go           # Main decoding logic
│   ├── scanner.go            # Line scanning/tokenization
│   ├── parser.go             # Parsing logic
│   └── validation.go         # Decode validation
├── decode_scanner.go         # Scanner (also in root for legacy compatibility)
├── decode_parser.go          # Parser (also in root for legacy compatibility)
├── decode_decoders.go        # Decoders (also in root for legacy compatibility)
├── decode_validation.go      # Validation (also in root for legacy compatibility)
├── toon_test.go              # Decoder tests
└── encode_test.go            # Encoder tests
```

## Migration Progress

### Phase 1: Foundation ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Project setup | ✅ | go/go.mod | Initialize Go module |
| Constants | ✅ | go/constants.go | Port delimiters, markers, literals |
| Types | ✅ | go/types.go | Define Go types for options, parsing |
| Main API | ✅ | go/toon.go | Encode + Decode functions implemented |

### Phase 2: Shared Utilities ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| String utils | ✅ | go/shared/string_utils.go | Quoting, escaping, unquoting, quote finding |
| Literal utils | ✅ | go/shared/literal_utils.go | Parse bool, number, null |

### Phase 3: Decoder ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Scanner | ✅ | go/decode_scanner.go | Line tokenization, indentation |
| Parser | ✅ | go/decode_parser.go | Header parsing, structural analysis |
| Decoders | ✅ | go/decode_decoders.go | Value decoding logic |
| Validation | ✅ | go/decode_validation.go | Strict mode validation |

### Phase 4: Encoder ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Normalize | ✅ | go/encode/normalize.go | Value normalization (Date, BigInt, etc) |
| Primitives | ✅ | go/encode/primitives.go | Primitive encoding |
| Writer | ✅ | go/encode/writer.go | Output writing with indentation |
| Encoders | ✅ | go/encode/encoders.go | Main encoding logic |
| Validation | ✅ | go/encode/validation.go | Quoting and key validation |

### Phase 5: Testing ⏳

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Decode tests | ✅ | go/toon_test.go | Basic decode tests passing |
| Encode tests | ✅ | go/encode_test.go | Basic encode tests passing |
| Round-trip tests | ✅ | go/encode_test.go | Encode/decode round-trip tests passing |
| Conformance tests | ⏳ | - | Use spec conformance tests |

### Phase 6: CLI (Optional) ⏳

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| CLI tool | ⏳ | cmd/toon/main.go | Command-line interface |

## Next Steps Checklist

- [x] Initialize Go module and project structure
- [x] Implement constants.go and types.go
- [x] Implement shared utilities (strings, literals)
- [x] Implement decoder (scanner → parser → decoders → validation)
- [x] Write basic decoder tests (all passing)
- [x] Implement encoder (normalize → primitives → writer → encoders)
- [x] Write basic encoder tests (all passing)
- [ ] Port comprehensive unit tests from TypeScript
- [ ] Run conformance tests from spec repository
- [ ] Add CLI tool (optional)
- [ ] Documentation and README for Go package
- [ ] Benchmark against TypeScript implementation

## Key Implementation Notes

### Go-Specific Considerations

1. **Type System**: Go doesn't have union types like TypeScript. Use `interface{}` or `any` (Go 1.18+) for JSON values.
2. **Error Handling**: Return errors explicitly instead of throwing exceptions.
3. **JSON Compatibility**: Use `encoding/json` standard library for JSON value handling.
4. **String Builder**: Use `strings.Builder` for efficient string concatenation.
5. **Runes vs Bytes**: Be careful with Unicode handling (use runes for character iteration).

### Testing Strategy

1. Port existing TypeScript test cases to Go
2. Use table-driven tests (idiomatic Go pattern)
3. Validate against conformance tests from [toon-format/spec](https://github.com/toon-format/spec/tree/main/tests)
4. Compare output with TypeScript implementation for the same inputs

## References

- [TOON Specification v1.4](https://github.com/toon-format/spec/blob/main/SPEC.md)
- [TypeScript Reference Implementation](https://github.com/toon-format/toon)
- [Conformance Tests](https://github.com/toon-format/spec/tree/main/tests)
- [Other Go Implementation](https://github.com/alpkeskin/gotoon) (community, for reference)

## Current Implementation Summary

### ✅ What's Working

**Decoder (go/):**
- ✅ Parses all TOON formats: objects, arrays (tabular, inline, list), primitives
- ✅ Handles nested structures and mixed types
- ✅ Strict mode validation (array lengths, structure)
- ✅ Proper unescaping and quote handling
- ✅ All decode tests passing (5/5)

**Encoder (go/encode/):**
- ✅ Encodes all Go types to TOON format
- ✅ Automatic tabular format detection for uniform object arrays
- ✅ Value normalization (Date → ISO, BigInt → number/string, structs → objects)
- ✅ Custom delimiters (comma, tab, pipe)
- ✅ Optional length markers
- ✅ Proper quoting and escaping
- ✅ All encode tests passing (8/8)

**Round-trip:**
- ✅ Encode → Decode → works correctly
- ✅ Data integrity maintained

### ⏳ What's Next

1. **Comprehensive Testing**: Port full TypeScript test suite
2. **Conformance Tests**: Run against official spec tests
3. **CLI Tool**: Optional command-line interface (like TypeScript version)
4. **Documentation**: README.md for Go package
5. **Performance**: Benchmarking vs TypeScript implementation

### 📝 Implementation Notes

- **File Organization**: Decoder files are currently in `go/` root (decode_*.go) but could be moved to `go/decode/` for consistency with encoder
- **Map Iteration**: Go maps have non-deterministic iteration order, so field order in encoded objects may vary
- **Type Handling**: Go's type system requires explicit type assertions; using `interface{}` for JSON values
- **Error Handling**: All errors returned explicitly (no panics)

## Questions / Issues

None currently - core implementation complete!

---

Last updated: 2025-11-06
