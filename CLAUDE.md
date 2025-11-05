# TOON Go Migration Guide

This document tracks the migration of the TOON (Token-Oriented Object Notation) library from TypeScript to Go.

## Overview

TOON is a compact, human-readable serialization format designed for passing structured data to Large Language Models with significantly reduced token usage. This Go implementation aims to be a complete port of the reference TypeScript implementation.

## Project Structure

### TypeScript Source (Reference)
```
packages/toon/src/
├── index.ts              # Main entry point (encode/decode)
├── types.ts              # Type definitions
├── constants.ts          # Constants
├── shared/
│   ├── literal-utils.ts  # Boolean/null/number literal parsing
│   ├── string-utils.ts   # String escaping/unescaping
│   └── validation.ts     # Validation utilities
├── decode/
│   ├── parser.ts         # Array header & primitive parsing
│   ├── scanner.ts        # Line scanning & cursor
│   ├── decoders.ts       # Main decode logic
│   └── validation.ts     # Decode-specific validation
└── encode/
    ├── encoders.ts       # Main encode logic
    ├── normalize.ts      # Value normalization
    ├── writer.ts         # Line writing utilities
    └── primitives.ts     # Primitive encoding
```

### Go Target Structure
```
toon-go/
├── go.mod
├── go.sum
├── toon.go               # Main entry point (Encode/Decode)
├── types.go              # Type definitions
├── constants.go          # Constants
├── shared/
│   ├── literal.go        # Literal parsing utilities
│   ├── string.go         # String utilities
│   └── validation.go     # Validation utilities
├── decode/
│   ├── parser.go         # Array header & primitive parsing
│   ├── scanner.go        # Line scanning & cursor
│   ├── decoder.go        # Main decode logic
│   └── validation.go     # Decode validation
├── encode/
│   ├── encoder.go        # Main encode logic
│   ├── normalize.go      # Value normalization
│   ├── writer.go         # Line writing utilities
│   └── primitives.go     # Primitive encoding
└── tests/
    ├── encode_test.go
    ├── decode_test.go
    └── integration_test.go
```

## Migration Progress

| Module | TypeScript File | Go File | Status | Notes |
|--------|----------------|---------|--------|-------|
| **Core** | | | | |
| Main API | `index.ts` | `toon.go` | ⬜ Not Started | Encode/Decode entry points |
| Types | `types.ts` | `types.go` | ⬜ Not Started | Type definitions |
| Constants | `constants.ts` | `constants.go` | ⬜ Not Started | Character & delimiter constants |
| **Shared Utilities** | | | | |
| Literal Utils | `shared/literal-utils.ts` | `shared/literal.go` | ⬜ Not Started | Boolean/null/number parsing |
| String Utils | `shared/string-utils.ts` | `shared/string.go` | ⬜ Not Started | Escape/unescape, quoting |
| Validation | `shared/validation.ts` | `shared/validation.go` | ⬜ Not Started | Common validation |
| **Decoder** | | | | |
| Parser | `decode/parser.ts` | `decode/parser.go` | ⬜ Not Started | Array header parsing, primitives |
| Scanner | `decode/scanner.ts` | `decode/scanner.go` | ⬜ Not Started | Line scanning, cursor |
| Decoders | `decode/decoders.ts` | `decode/decoder.go` | ⬜ Not Started | Main decode logic |
| Validation | `decode/validation.ts` | `decode/validation.go` | ⬜ Not Started | Decode validation |
| **Encoder** | | | | |
| Encoders | `encode/encoders.ts` | `encode/encoder.go` | ⬜ Not Started | Main encode logic |
| Normalize | `encode/normalize.ts` | `encode/normalize.go` | ⬜ Not Started | Value normalization |
| Writer | `encode/writer.ts` | `encode/writer.go` | ⬜ Not Started | Line writing utilities |
| Primitives | `encode/primitives.ts` | `encode/primitives.go` | ⬜ Not Started | Primitive encoding |
| **Tests** | | | | |
| Encode Tests | `test/encode.test.ts` | `tests/encode_test.go` | ⬜ Not Started | Encoding tests |
| Decode Tests | `test/decode.test.ts` | `tests/decode_test.go` | ⬜ Not Started | Decoding tests |
| Integration | - | `tests/integration_test.go` | ⬜ Not Started | Round-trip tests |

## Legend
- ✅ Completed
- 🚧 In Progress
- ⬜ Not Started
- ❌ Blocked

## Next Steps Checklist

### Phase 1: Foundation (Current)
- [ ] Set up Go module (`go.mod`)
- [ ] Create directory structure
- [ ] Implement `types.go` (type definitions)
- [ ] Implement `constants.go` (constants)

### Phase 2: Shared Utilities
- [ ] Implement `shared/literal.go` (literal parsing)
- [ ] Implement `shared/string.go` (string utilities)
- [ ] Implement `shared/validation.go` (validation)

### Phase 3: Parser Implementation
- [ ] Implement `decode/parser.go` (array header parsing, primitive parsing)
- [ ] Implement `decode/scanner.go` (line scanning, cursor)

### Phase 4: Decoder
- [ ] Implement `decode/decoder.go` (main decode logic)
- [ ] Implement `decode/validation.go` (decode validation)
- [ ] Create decoder tests

### Phase 5: Encoder
- [ ] Implement `encode/normalize.go` (value normalization)
- [ ] Implement `encode/primitives.go` (primitive encoding)
- [ ] Implement `encode/writer.go` (line writing)
- [ ] Implement `encode/encoder.go` (main encode logic)
- [ ] Create encoder tests

### Phase 6: Integration & Testing
- [ ] Create integration tests
- [ ] Run conformance tests from spec repo
- [ ] Performance benchmarks
- [ ] Documentation

## Implementation Notes

### Key Differences: TypeScript → Go

1. **Type System**
   - TypeScript has union types (`string | number | boolean | null`)
   - Go will use `interface{}` or type assertions with custom types

2. **JSON Values**
   - TypeScript: `JsonValue` type alias
   - Go: Custom type or `interface{}`

3. **Options Pattern**
   - TypeScript: Optional parameters with defaults
   - Go: Struct with options (functional options pattern or config struct)

4. **Error Handling**
   - TypeScript: Throw exceptions
   - Go: Return errors explicitly

5. **String Handling**
   - TypeScript: Direct string manipulation
   - Go: Use `strings`, `bytes` packages; runes for Unicode

6. **Arrays & Objects**
   - TypeScript: Native arrays and objects
   - Go: Slices and maps

### Critical Functions to Port

1. **Parser** (`decode/parser.ts`)
   - `parseArrayHeaderLine()` - Parse array headers like `items[2]{id,name}:`
   - `parseBracketSegment()` - Parse `[N]` or `[N|]` segments
   - `parseDelimitedValues()` - Split values by delimiter (accounting for quotes)
   - `parsePrimitiveToken()` - Parse primitives (string, number, boolean, null)

2. **Scanner** (`decode/scanner.ts`)
   - `toParsedLines()` - Convert input to parsed lines
   - `LineCursor` class - Track current position in parsed lines

3. **Decoders** (`decode/decoders.ts`)
   - `decodeValueFromLines()` - Main decode entry point
   - Array/object/primitive decoding logic

4. **Encoders** (`encode/encoders.ts`)
   - `encodeValue()` - Main encode entry point
   - Array/object/primitive encoding logic

## Testing Strategy

1. **Unit Tests**: Test each module independently
2. **Integration Tests**: Test encode/decode round-trips
3. **Conformance Tests**: Use test fixtures from https://github.com/toon-format/spec
4. **Benchmarks**: Compare performance with TypeScript implementation

## Resources

- **TOON Spec**: https://github.com/toon-format/spec/blob/main/SPEC.md (v1.4)
- **Conformance Tests**: https://github.com/toon-format/spec/tree/main/tests
- **TypeScript Reference**: `packages/toon/src/`
- **Other Go Implementation**: https://github.com/alpkeskin/gotoon (for reference)

## Development Workflow

1. Read the TypeScript implementation
2. Implement the Go equivalent
3. Write tests (referencing TypeScript tests)
4. Update this document's progress table
5. Commit with clear messages
6. Repeat for next module

## Questions / Decisions

- [ ] Use functional options pattern or config struct for options?
- [ ] Use custom types or `interface{}` for JSON values?
- [ ] Package structure: single package or multiple packages?
