package encode

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/toon-format/toon/go/shared"
)

const (
	nullLiteral  = "null"
	trueLiteral  = "true"
	falseLiteral = "false"
	doubleQuote  = '"'
	comma        = ','
)

// EncodePrimitive encodes a JSON primitive value (null, bool, number, string)
func EncodePrimitive(value interface{}, delimiter rune) string {
	if value == nil {
		return nullLiteral
	}

	switch v := value.(type) {
	case bool:
		if v {
			return trueLiteral
		}
		return falseLiteral

	case float64:
		return formatNumber(v)

	case string:
		return EncodeStringLiteral(v, delimiter)

	default:
		// Shouldn't happen if value is already normalized, but handle it
		return nullLiteral
	}
}

// EncodeStringLiteral encodes a string, quoting it if necessary
func EncodeStringLiteral(value string, delimiter rune) string {
	if IsSafeUnquoted(value, delimiter) {
		return value
	}

	return string(doubleQuote) + shared.EscapeString(value) + string(doubleQuote)
}

// EncodeKey encodes an object key, quoting it if necessary
func EncodeKey(key string) string {
	if IsValidUnquotedKey(key) {
		return key
	}

	return string(doubleQuote) + shared.EscapeString(key) + string(doubleQuote)
}

// EncodeAndJoinPrimitives encodes an array of primitives and joins them with the delimiter
func EncodeAndJoinPrimitives(values []interface{}, delimiter rune) string {
	encoded := make([]string, len(values))
	for i, v := range values {
		encoded[i] = EncodePrimitive(v, delimiter)
	}
	return strings.Join(encoded, string(delimiter))
}

// FormatHeaderOptions contains options for formatting array headers
type FormatHeaderOptions struct {
	Key          string   // Optional key name
	Fields       []string // Optional field names for tabular arrays
	Delimiter    rune     // Delimiter character
	LengthMarker rune     // Optional length marker ('#' or 0)
}

// FormatHeader formats an array header line
// Examples:
//   - items[3]:
//   - items[#3]:
//   - items[3|]:
//   - items[3]{id,name,price}:
func FormatHeader(length int, opts FormatHeaderOptions) string {
	var header strings.Builder

	// Add key if present
	if opts.Key != "" {
		header.WriteString(EncodeKey(opts.Key))
	}

	// Add array length with optional marker and delimiter
	header.WriteRune('[')
	if opts.LengthMarker != 0 {
		header.WriteRune(opts.LengthMarker)
	}
	header.WriteString(strconv.Itoa(length))

	// Only include delimiter if it's not the default (comma)
	if opts.Delimiter != 0 && opts.Delimiter != comma {
		header.WriteRune(opts.Delimiter)
	}

	header.WriteRune(']')

	// Add fields if present (for tabular arrays)
	if len(opts.Fields) > 0 {
		header.WriteRune('{')
		encodedFields := make([]string, len(opts.Fields))
		for i, field := range opts.Fields {
			encodedFields[i] = EncodeKey(field)
		}
		header.WriteString(strings.Join(encodedFields, string(opts.Delimiter)))
		header.WriteRune('}')
	}

	header.WriteRune(':')

	return header.String()
}

// formatNumber formats a float64 as a string without scientific notation
func formatNumber(f float64) string {
	// For whole numbers, format without decimal point
	if f == float64(int64(f)) {
		return fmt.Sprintf("%.0f", f)
	}
	// For decimal numbers, use %g but ensure we don't get scientific notation
	// Use %f with enough precision, then trim trailing zeros
	s := strconv.FormatFloat(f, 'f', -1, 64)
	return s
}
