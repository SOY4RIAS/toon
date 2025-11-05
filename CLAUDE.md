# CLAUDE.md - AI Context for TOON Go Implementation

> **Purpose**: This file provides essential context for AI assistants working on this codebase. It's designed to minimize context window usage while maximizing understanding.

## 🎯 Project Overview

**What**: Port of TOON (Token-Oriented Object Notation) format from TypeScript to Go
**TOON**: A token-efficient JSON alternative for LLM prompts (30-60% fewer tokens than JSON)
**Original**: TypeScript implementation at `packages/toon/` (~2089 lines)
**Target**: 100% compatible Go implementation
**Status**: Encoder ✅ Complete | Decoder 🚧 20% Complete

## ⚠️ Critical Constraints

1. **100% Output Compatibility**: Go encoder MUST produce byte-identical output to TypeScript version
2. **Everything in English**: Code, comments, documentation, commit messages - NO Spanish
3. **Spec Compliance**: Must pass [TOON Spec v1.4](https://github.com/toon-format/spec) conformance tests
4. **No Breaking Changes**: Port behavior exactly, don't "improve" the format

## 📁 Architecture Overview

```
/home/user/toon/
├── packages/toon/src/          # Original TypeScript (READ-ONLY reference)
│   ├── encode/                 # 📚 Reference for encoding logic
│   ├── decode/                 # 📚 Reference for decoding logic (960 lines)
│   └── shared/                 # 📚 Reference for utilities
│
├── [Go Implementation ROOT]    # New Go code at repository root
│   ├── encode/                 # ✅ COMPLETE - Encoding implementation
│   │   ├── normalize.go        # Value normalization (Date, BigInt, NaN→null)
│   │   ├── primitives.go       # Primitive encoding + header formatting
│   │   ├── encoders.go         # Main encoding logic (objects, arrays)
│   │   └── writer.go           # Line writer with indentation
│   │
│   ├── decode/                 # 🚧 IN PROGRESS
│   │   ├── scanner.go          # ✅ DONE - Line scanning/tokenization
│   │   ├── decoders.go         # ⚠️  STUB - Needs full implementation
│   │   ├── parser.go           # ❌ TODO - Parse headers, primitives, keys
│   │   └── validation.go       # ❌ TODO - Strict mode validation
│   │
│   ├── shared/                 # ✅ COMPLETE - Shared utilities
│   │   ├── strings.go          # Escape/unescape, quote finding
│   │   ├── validation.go       # Quote/key validation
│   │   └── literals.go         # Boolean/null/numeric detection
│   │
│   ├── internal/types/         # ✅ Internal types (avoids import cycles)
│   │   └── types.go            # ParsedLine, ScanResult, BlankLineInfo
│   │
│   ├── toon.go                 # ✅ Public API (Encode/Decode)
│   ├── types.go                # ✅ Public types (EncodeOptions, etc.)
│   ├── constants.go            # ✅ Constants (delimiters, literals)
│   └── errors.go               # ✅ Error types
```

## 🔑 Key Concepts

### TOON Format Basics
```
// JSON (verbose)
{"users": [
  {"id": 1, "name": "Alice", "role": "admin"},
  {"id": 2, "name": "Bob", "role": "user"}
]}

// TOON (compact)
users[2]{id,name,role}:
  1,Alice,admin
  2,Bob,user
```

**Three Array Formats**:
1. **Inline**: `tags[3]: admin,ops,dev` - for primitive arrays
2. **Tabular**: `users[2]{id,name}:` - for uniform objects (same keys, all primitive values)
3. **List**: `items[2]:` with `- ` prefixes - for mixed/non-uniform arrays

### Go-Specific Design Decisions

**Import Cycle Prevention**:
- `internal/types/` holds `ParsedLine`, `ScanResult`, `BlankLineInfo`
- Prevents cycles: `toon` → `decode` → would import `toon` (❌)
- Solution: `decode` uses `internal/types` (✅)

**Type Mapping**:
| TypeScript | Go | Notes |
|------------|-----|-------|
| `unknown` | `interface{}` | Use type switches/assertions |
| `JsonValue` | `interface{}` | Runtime type checking |
| `number` | `float64` | All numbers as float64 |
| `BigInt` | `string` or `float64` | Convert if in safe range |
| `-0` | `0` | Normalize negative zero |
| `NaN`/`Infinity` | `null` | Non-finite → null |

## ✅ What's Working

### Encoder (100% Complete)
```bash
go test -v ./examples/
# All tests pass ✅
```

**Features**:
- Primitives (string, number, boolean, null)
- Objects (nested, empty)
- Arrays (inline, tabular, list)
- Delimiters (`,` `\t` `|`)
- Length markers (`[#3]`)
- Quoting rules (minimal quotes)

**Test Coverage**:
- ✅ Basic encode (tabular arrays)
- ✅ Primitive objects
- ✅ Inline arrays

### Scanner (100% Complete)
- Line parsing with indentation tracking
- Blank line detection
- Strict mode validation (tab rejection, indent multiples)
- Depth calculation

## 🚧 What Needs Work

### 1. Decoder Parser (`decode/parser.go`) - ~300 lines
**Reference**: `packages/toon/src/decode/parser.ts`

**Must Implement**:
- `parseArrayHeaderLine()` - Parse `key[N]{fields}:` or `[N]:` headers
- `parseBracketSegment()` - Extract length, delimiter, lengthMarker from `[#3|]`
- `parseDelimitedValues()` - Split by delimiter, respect quotes
- `parsePrimitiveToken()` - Parse string/number/boolean/null
- `parseStringLiteral()` - Handle quoted strings with escapes
- `parseKeyToken()` - Extract keys from `key:` or `"quoted key":`

**Key Functions**:
```go
// Parse array header: users[2]{id,name}: value,value
type ArrayHeaderInfo struct {
    Key             string    // "users" or empty for root
    Length          int       // 2
    Delimiter       rune      // ',' '\t' or '|'
    Fields          []string  // ["id", "name"] or nil
    HasLengthMarker bool      // true if [#2]
}

func parseArrayHeaderLine(content string, defaultDelim rune) (*ArrayHeaderInfo, string, error)
func parsePrimitiveToken(token string) (interface{}, error)
func parseDelimitedValues(input string, delim rune) ([]string, error)
```

### 2. Decoder Logic (`decode/decoders.go`) - ~500 lines
**Reference**: `packages/toon/src/decode/decoders.ts`

**Current State**: Stub that returns `map[string]interface{}` with placeholder data

**Must Implement**:
```go
// Entry point - determine root type (object, array, primitive)
func decodeValueFromLines(cursor *LineCursor, opts *ResolvedDecodeOptions) (interface{}, error)

// Object decoding - parse key-value pairs at same depth
func decodeObject(cursor *LineCursor, baseDepth int, opts) (map[string]interface{}, error)
func decodeKeyValuePair(line *ParsedLine, cursor, depth, opts) (string, interface{})

// Array decoding - handle inline, tabular, and list formats
func decodeArrayFromHeader(header *ArrayHeaderInfo, inlineVals string, cursor, depth, opts) ([]interface{}, error)
func decodeInlinePrimitiveArray(header, values, opts) ([]interface{}, error)
func decodeTabularArray(header, cursor, depth, opts) ([]map[string]interface{}, error)
func decodeListArray(header, cursor, depth, opts) ([]interface{}, error)

// List item decoding
func decodeListItemValue(cursor, depth, opts) (interface{}, error)
func decodeObjectAsListItem(cursor, depth, opts) (map[string]interface{}, error)
```

**Critical Logic**:
- Detect root type (object, array, single primitive)
- Track depth correctly for nested structures
- Handle `- ` prefixed list items
- Validate array lengths in strict mode
- Parse tabular rows (no `- ` prefix, just comma-separated values)

### 3. Validation (`decode/validation.go`) - ~160 lines
**Reference**: `packages/toon/src/decode/validation.ts`

```go
func assertExpectedCount(actual, expected int, itemType string, opts) error
func validateNoExtraListItems(cursor *LineCursor, itemDepth, expectedCount int) error
func validateNoExtraTabularRows(cursor *LineCursor, rowDepth int, header *ArrayHeaderInfo) error
func validateNoBlankLinesInRange(startLine, endLine int, blankLines []BlankLineInfo, strict bool, context string) error
```

## 🧪 Testing Strategy

### Phase 1: Unit Tests (TODO)
Port from `packages/toon/test/*.test.ts`:
- `encode.test.ts` → `encode_test.go`
- `decode.test.ts` → `decode_test.go`
- `normalization.test.ts` → `normalization_test.go`

### Phase 2: Conformance Tests (TODO)
Download and run: https://github.com/toon-format/spec/tree/main/tests
```bash
# Clone spec tests
git clone https://github.com/toon-format/spec.git /tmp/toon-spec

# Run conformance suite (after implementing decoder)
go test ./conformance/ -v
```

### Phase 3: Cross-Validation (TODO)
For each test case:
1. Encode in Go → compare with TypeScript output (byte-exact)
2. Decode TOON → compare structures (deep equal)

## 🚨 Common Pitfalls

### 1. Import Cycles
❌ **Don't**: Import `github.com/SOY4RIAS/toon` from `decode/` or `encode/`
✅ **Do**: Use `internal/types` for shared types

### 2. Map Iteration Order
Go maps have random iteration order. When encoding objects:
```go
// ❌ Wrong - random key order
for key, val := range obj {
    encodeKeyValuePair(key, val, ...)
}

// ✅ For consistent output, sort keys first (if needed for tests)
keys := make([]string, 0, len(obj))
for k := range obj {
    keys = append(keys, k)
}
sort.Strings(keys)  // Only if spec requires it
```

**Note**: Current implementation doesn't sort - verify if TypeScript does!

### 3. Numeric Precision
```go
// TypeScript: Number.parseFloat("123.45")
// Go: strconv.ParseFloat("123.45", 64)

// TypeScript: -0 === 0 (but Object.is(-0, 0) = false)
// Go: math.Signbit(f) to detect -0
```

### 4. Indentation Strictness
Strict mode requires:
- No tabs in indentation
- Indent must be exact multiple of `indentSize` (default 2)
- No blank lines inside arrays/tables

## 📋 Next Steps Checklist

### Immediate (Decoder Completion)
- [ ] Create `decode/parser.go` with all parsing functions
- [ ] Implement full `decode/decoders.go` logic
- [ ] Create `decode/validation.go` for strict mode
- [ ] Test decoder with simple cases
- [ ] Verify decode(encode(data)) === data

### Testing
- [ ] Port encode tests from TypeScript
- [ ] Port decode tests from TypeScript
- [ ] Port normalization tests
- [ ] Download and run spec conformance tests
- [ ] Create cross-validation tests (Go ↔ TypeScript)

### Polish
- [ ] Create CLI tool (`cmd/toon/main.go`)
- [ ] Write comprehensive README.md
- [ ] Add godoc comments to all exports
- [ ] Setup GitHub Actions CI
- [ ] Add performance benchmarks

## 🔍 Quick Reference

### Run Tests
```bash
go test -v ./examples/        # Current encoder tests
go test -v ./...              # All tests (when decoder done)
go test -bench=. ./...        # Benchmarks (create these)
```

### Build
```bash
go build ./...                # Build all packages
go build ./cmd/toon           # Build CLI (create this)
```

### Key Files to Reference
- TypeScript encoder: `packages/toon/src/encode/encoders.ts`
- TypeScript decoder: `packages/toon/src/decode/decoders.ts`
- TypeScript parser: `packages/toon/src/decode/parser.ts`
- Spec: https://github.com/toon-format/spec/blob/main/SPEC.md

## 📞 Getting Help

### Understanding TOON Format
1. Read `README.md` in repository root
2. Check `SPEC.md` for precise rules
3. Look at examples in `benchmarks/` directory

### Debugging Encoder/Decoder Mismatch
1. Enable verbose logging (add debug prints)
2. Compare outputs character-by-character
3. Check delimiter detection (comma vs tab vs pipe)
4. Verify indentation (2 spaces default)

### Import Errors
If you see "import cycle not allowed":
- Check if `decode/` or `encode/` imports `toon`
- Move types to `internal/types` if needed
- Use interfaces to break dependency

## 📊 Progress Tracking

Last updated: 2025-11-05

| Component | Status | Lines | Notes |
|-----------|--------|-------|-------|
| Types & Constants | ✅ | ~150 | Complete |
| Shared Utils | ✅ | ~250 | Complete |
| Encoder | ✅ | ~600 | **All tests pass** |
| Scanner | ✅ | ~160 | Complete |
| Parser | ❌ | 0/~300 | **Next priority** |
| Decoders | ⚠️ | ~40/~500 | Stub only |
| Validation | ❌ | 0/~160 | Not started |
| Tests | ⚠️ | ~80/~500 | Only encoder examples |
| CLI | ❌ | 0/~200 | Not started |

**Total Progress**: ~40% complete

---

## 💡 Pro Tips for AI Assistants

1. **Reference First**: Always check TypeScript implementation before implementing
2. **Test Early**: Create test after each function, not at the end
3. **Byte-Exact**: Use `go test` to verify output matches TypeScript
4. **Context Management**: This file should be updated as work progresses
5. **Commit Often**: Small, focused commits with clear messages
6. **English Only**: Double-check all text is in English before committing

## 🎓 Learning Resources

- **Go Reflection**: Used in `encode/normalize.go` for struct handling
- **String Building**: `strings.Builder` for efficient concatenation
- **Regex in Go**: `regexp.Compile()` patterns in `shared/validation.go`
- **Error Handling**: Return `error` types, use `fmt.Errorf()` for context

---

**End of Context Document** - This file should be your first read when resuming work on this project.
