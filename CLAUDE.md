# Go Migration Progress

This document tracks the progress of migrating the TOON TypeScript implementation to Go.

## Project Overview

TOON (Token-Oriented Object Notation) is a compact, human-readable serialization format designed for passing structured data to Large Language Models with significantly reduced token usage.

**TypeScript Source:** `packages/toon/src/`
**Go Target:** `go/`

## Architecture

The TypeScript implementation consists of:

### Core Modules

| Module | TypeScript Path | Go Path | Status | Notes |
|--------|----------------|---------|--------|-------|
| **Types & Constants** | `src/types.ts`, `src/constants.ts` | `go/types.go`, `go/constants.go` | ✅ Complete | Core type definitions and constants |
| **Main API** | `src/index.ts` | `go/toon.go` | ⬜ Not Started | Public encode/decode API |

### Decode (Parser)

| Component | TypeScript Path | Go Path | Status | Notes |
|-----------|----------------|---------|--------|-------|
| **Scanner** | `src/decode/scanner.ts` | `go/decode/scanner.go` | ✅ Complete | Line scanning and cursor |
| **Parser** | `src/decode/parser.ts` | `go/decode/parser.go` | ✅ Complete | Array headers, delimited values, primitives |
| **Decoders** | `src/decode/decoders.ts` | `go/decode/decoders.go` | ✅ Complete | Value decoding logic |
| **Validation** | `src/decode/validation.ts` | `go/decode/validation.go` | ✅ Complete | Strict mode validation |

### Encode (Serializer)

| Component | TypeScript Path | Go Path | Status | Notes |
|-----------|----------------|---------|--------|-------|
| **Normalize** | `src/encode/normalize.ts` | `go/encode/normalize.go` | ⬜ Not Started | Value normalization |
| **Encoders** | `src/encode/encoders.ts` | `go/encode/encoders.go` | ⬜ Not Started | Value encoding logic |
| **Writer** | `src/encode/writer.ts` | `go/encode/writer.go` | ⬜ Not Started | Output formatting |
| **Primitives** | `src/encode/primitives.ts` | `go/encode/primitives.go` | ⬜ Not Started | Primitive value handling |

### Shared Utilities

| Component | TypeScript Path | Go Path | Status | Notes |
|-----------|----------------|---------|--------|-------|
| **String Utils** | `src/shared/string-utils.ts` | `go/shared/string_utils.go` | ✅ Complete | Quote handling, escaping |
| **Literal Utils** | `src/shared/literal-utils.ts` | `go/shared/literal_utils.go` | ✅ Complete | Boolean, null, number parsing |
| **Validation** | `src/shared/validation.ts` | `go/shared/validation.go` | ⬜ Not Started | Common validation |

### Testing

| Component | TypeScript Path | Go Path | Status | Notes |
|-----------|----------------|---------|--------|-------|
| **Encode Tests** | `test/encode.test.ts` | `go/encode_test.go` | ⬜ Not Started | Encoding tests |
| **Decode Tests** | `test/decode.test.ts` | `go/decode_test.go` | ⬜ Not Started | Decoding tests |
| **Normalization Tests** | `test/normalization.test.ts` | `go/normalization_test.go` | ⬜ Not Started | Normalization tests |

### CLI

| Component | TypeScript Path | Go Path | Status | Notes |
|-----------|----------------|---------|--------|-------|
| **CLI** | `packages/cli/` | `go/cmd/toon/` | ⬜ Not Started | Command-line interface |

## Next Steps Checklist

### Phase 1: Foundation ✅ Complete
- [x] Create Go module structure (`go.mod`)
- [x] Implement core types and constants
- [x] Implement shared utilities (string-utils, literal-utils)
- [x] Set up basic testing infrastructure

### Phase 2: Parser (Decode) ✅ Complete
- [x] Implement scanner (line parsing, cursor)
- [x] Implement parser (array headers, delimited values, primitives)
- [x] Implement decoders (value decoding logic)
- [x] Implement validation
- [x] Add comprehensive decode tests (all passing)

### Phase 3: Serializer (Encode)
- [ ] Implement value normalization
- [ ] Implement encoders
- [ ] Implement writer
- [ ] Implement primitive handling
- [ ] Port encode tests

### Phase 4: API & CLI
- [ ] Implement public encode/decode API
- [ ] Implement CLI
- [ ] Integration testing
- [ ] Documentation

## Current Focus

**Decoder Implementation** - Implementing the decode/decoders module to complete the parsing functionality

## Recent Progress (2025-11-05)

- ✅ Created Go module structure
- ✅ Implemented core types and constants
- ✅ Implemented shared string utilities (escaping, unescaping, quote finding)
- ✅ Implemented shared literal utilities (boolean, null, number detection)
- ✅ Implemented scanner module (line parsing, cursor navigation)
- ✅ Implemented parser module (array headers, delimited values, primitives)
- ✅ Added comprehensive tests for parser functionality (all passing)

## Key Differences: TypeScript vs Go

### Type System
- TypeScript uses union types (`JsonValue = JsonPrimitive | JsonObject | JsonArray`)
- Go will use interfaces and type switches
- Consider using `interface{}` or `any` for JsonValue

### String Handling
- TypeScript: strings are immutable, use string methods
- Go: strings are immutable, but use bytes for manipulation
- Go rune iteration for unicode handling

### Error Handling
- TypeScript: throw exceptions
- Go: return errors as values

### Options Pattern
- TypeScript: optional parameters with defaults
- Go: functional options or config structs

## References

- [TOON Specification v1.4](https://github.com/toon-format/spec)
- [TypeScript Implementation](packages/toon/)
- [Conformance Tests](https://github.com/toon-format/spec/tree/main/tests)

## Status Legend

- ✅ Complete
- 🚧 In Progress
- ⬜ Not Started
- ❌ Blocked

## Files Created

### Core
- `go/types.go` - Type definitions
- `go/constants.go` - Constants and delimiters

### Shared Utilities
- `go/shared/string_utils.go` - String escaping and quote handling
- `go/shared/literal_utils.go` - Literal value parsing

### Decode (Parser)
- `go/decode/scanner.go` - Line scanning and cursor
- `go/decode/parser.go` - Array header and value parsing
- `go/decode/parser_test.go` - Parser tests (all passing)

---

Last Updated: 2025-11-05
