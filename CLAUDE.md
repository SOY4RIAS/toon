# TOON Go Migration - Progress Tracker

This document tracks the progress of migrating the TOON (Token-Oriented Object Notation) library from TypeScript to Go.

## Project Overview

**Goal**: Create a Go implementation of the TOON format parser/encoder that matches the TypeScript reference implementation.

**Repository**: `toon-format/toon`
**Branch**: `claude/go-migration-parser-011CUqhpS7DatVmPCe2tsdQX`
**Spec Version**: v1.4 ([spec repository](https://github.com/toon-format/spec))

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
├── go.mod
├── go.sum
├── toon.go               # Main entry point (Encode/Decode)
├── types.go              # Type definitions
├── constants.go          # Constants and delimiters
├── shared/               # Shared utilities
│   ├── strings.go        # String manipulation
│   ├── literals.go       # Literal parsing
│   └── validation.go     # Validation helpers
├── encode/               # Encoder implementation
│   ├── primitives.go     # Primitive encoding
│   ├── writer.go         # Output writer
│   ├── normalize.go      # Value normalization
│   └── encoders.go       # Main encoding logic
├── decode/               # Decoder implementation
│   ├── decoders.go       # Main decoding logic
│   ├── scanner.go        # Line scanning/tokenization
│   ├── parser.go         # Parsing logic
│   └── validation.go     # Decode validation
└── toon_test.go          # Main tests
```

## Migration Progress

### Phase 1: Foundation ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Project setup | ✅ | go.mod | Initialize Go module |
| Constants | ✅ | constants.go | Port delimiters, markers, literals |
| Types | ✅ | types.go | Define Go types for options, parsing |
| Main API | ✅ | toon.go | Decode function implemented |

### Phase 2: Shared Utilities ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| String utils | ✅ | shared/strings.go | Quoting, escaping, unquoting |
| Literal utils | ✅ | shared/literals.go | Parse bool, number, null |

### Phase 3: Decoder ✅

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Scanner | ✅ | decode_scanner.go | Line tokenization, indentation |
| Parser | ✅ | decode_parser.go | Header parsing, structural analysis |
| Decoders | ✅ | decode_decoders.go | Value decoding logic |
| Validation | ✅ | decode_validation.go | Strict mode validation |

### Phase 4: Encoder ⏳

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Normalize | ⏳ | encode/normalize.go | Value normalization (Date, BigInt, etc) |
| Primitives | ⏳ | encode/primitives.go | Primitive encoding |
| Writer | ⏳ | encode/writer.go | Output writing with indentation |
| Encoders | ⏳ | encode/encoders.go | Main encoding logic |

### Phase 5: Testing ⏳

| Component | Status | File | Notes |
|-----------|--------|------|-------|
| Basic tests | ✅ | toon_test.go | Basic decode tests passing |
| Unit tests | ⏳ | *_test.go | More comprehensive tests needed |
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
- [ ] Implement encoder (normalize → primitives → writer → encoders)
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

## Questions / Issues

*(Track any questions or blockers here)*

---

Last updated: 2025-11-06
