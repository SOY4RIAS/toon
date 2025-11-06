package encode

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/SOY4RIAS/toon/shared"
)

const (
	nullLiteral   = "null"
	doubleQuote   = '"'
	comma         = ','
	defaultDelim  = comma
	hash          = '#'
	openBracket   = '['
	closeBracket  = ']'
	openBrace     = '{'
	closeBrace    = '}'
	colon         = ':'
)

// EncodePrimitive encodes a primitive JSON value to TOON format.
func EncodePrimitive(value interface{}, delimiter rune) string {
	if value == nil {
		return nullLiteral
	}

	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)

	case float64:
		return formatNumber(v)

	case float32:
		return formatNumber(float64(v))

	case int:
		return strconv.Itoa(v)

	case int64:
		return strconv.FormatInt(v, 10)

	case int32:
		return strconv.FormatInt(int64(v), 10)

	case int16:
		return strconv.FormatInt(int64(v), 10)

	case int8:
		return strconv.FormatInt(int64(v), 10)

	case uint:
		return strconv.FormatUint(uint64(v), 10)

	case uint64:
		return strconv.FormatUint(v, 10)

	case uint32:
		return strconv.FormatUint(uint64(v), 10)

	case uint16:
		return strconv.FormatUint(uint64(v), 10)

	case uint8:
		return strconv.FormatUint(uint64(v), 10)

	case string:
		return EncodeStringLiteral(v, delimiter)

	default:
		// Fallback: try to encode as string
		return EncodeStringLiteral(fmt.Sprint(v), delimiter)
	}
}

// formatNumber formats a number without scientific notation.
func formatNumber(f float64) string {
	// Handle special cases
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return nullLiteral
	}

	// Handle -0
	if f == 0 && math.Signbit(f) {
		return "0"
	}

	// Format without scientific notation
	s := strconv.FormatFloat(f, 'f', -1, 64)

	// Remove trailing zeros after decimal point
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	return s
}

// EncodeStringLiteral encodes a string value, adding quotes if necessary.
func EncodeStringLiteral(value string, delimiter rune) string {
	if shared.IsSafeUnquoted(value, delimiter) {
		return value
	}

	return fmt.Sprintf("%c%s%c", doubleQuote, shared.EscapeString(value), doubleQuote)
}

// EncodeKey encodes an object key, adding quotes if necessary.
func EncodeKey(key string) string {
	if shared.IsValidUnquotedKey(key) {
		return key
	}

	return fmt.Sprintf("%c%s%c", doubleQuote, shared.EscapeString(key), doubleQuote)
}

// EncodeAndJoinPrimitives encodes and joins multiple primitive values with a delimiter.
func EncodeAndJoinPrimitives(values []interface{}, delimiter rune) string {
	if len(values) == 0 {
		return ""
	}

	var parts []string
	for _, v := range values {
		parts = append(parts, EncodePrimitive(v, delimiter))
	}

	return strings.Join(parts, string(delimiter))
}

// HeaderOptions configures the format of an array header.
type HeaderOptions struct {
	Key          string
	Fields       []string
	Delimiter    rune
	LengthMarker bool
}

// FormatHeader formats an array header line.
func FormatHeader(length int, opts *HeaderOptions) string {
	var sb strings.Builder

	// Add key if present
	if opts != nil && opts.Key != "" {
		sb.WriteString(EncodeKey(opts.Key))
	}

	// Determine delimiter
	delim := comma
	if opts != nil && opts.Delimiter != 0 {
		delim = opts.Delimiter
	}

	// Build array length part: [N] or [#N] with optional delimiter
	sb.WriteRune(openBracket)

	// Add length marker if requested
	if opts != nil && opts.LengthMarker {
		sb.WriteRune(hash)
	}

	// Add length
	sb.WriteString(strconv.Itoa(length))

	// Add delimiter to header if not comma (default)
	if delim != defaultDelim {
		sb.WriteRune(delim)
	}

	sb.WriteRune(closeBracket)

	// Add fields if present (for tabular arrays)
	if opts != nil && len(opts.Fields) > 0 {
		sb.WriteRune(openBrace)

		encodedFields := make([]string, len(opts.Fields))
		for i, field := range opts.Fields {
			encodedFields[i] = EncodeKey(field)
		}

		sb.WriteString(strings.Join(encodedFields, string(delim)))
		sb.WriteRune(closeBrace)
	}

	// Add colon
	sb.WriteRune(colon)

	return sb.String()
}
