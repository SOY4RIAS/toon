package toon

import (
	"math"
	"strconv"
	"strings"
)

// IsBooleanOrNullLiteral checks if a token is a boolean or null literal (true, false, null).
func IsBooleanOrNullLiteral(token string) bool {
	return token == TrueLiteral || token == FalseLiteral || token == NullLiteral
}

// IsNumericLiteral checks if a token represents a valid numeric literal.
// Rejects numbers with leading zeros (except "0" itself or decimals like "0.5").
func IsNumericLiteral(token string) bool {
	if token == "" {
		return false
	}

	// Must not have leading zeros (except for "0" itself or decimals like "0.5")
	if len(token) > 1 && token[0] == '0' && token[1] != '.' {
		return false
	}

	// Check if it's a valid number
	numericValue, err := strconv.ParseFloat(token, 64)
	if err != nil {
		return false
	}

	// Check if it's finite (not NaN or Inf)
	return !math.IsNaN(numericValue) && !math.IsInf(numericValue, 0)
}

// IsValidUnquotedKey checks if a key can be used without quotes.
// Valid unquoted keys must start with a letter or underscore,
// followed by letters, digits, underscores, or dots.
func IsValidUnquotedKey(key string) bool {
	if key == "" {
		return false
	}

	// First character must be letter or underscore
	first := rune(key[0])
	if !isLetter(first) && first != '_' {
		return false
	}

	// Remaining characters must be letters, digits, underscores, or dots
	for _, ch := range key[1:] {
		if !isLetter(ch) && !isDigit(ch) && ch != '_' && ch != '.' {
			return false
		}
	}

	return true
}

// IsSafeUnquoted determines if a string value can be safely encoded without quotes.
// A string needs quoting if it:
// - Is empty
// - Has leading or trailing whitespace
// - Could be confused with a literal (boolean, null, number)
// - Contains structural characters (colons, brackets, braces)
// - Contains quotes or backslashes (need escaping)
// - Contains control characters (newlines, tabs, etc.)
// - Contains the active delimiter
// - Starts with a list marker (hyphen)
func IsSafeUnquoted(value string, delimiter Delimiter) bool {
	if value == "" {
		return false
	}

	if value != strings.TrimSpace(value) {
		return false
	}

	// Check if it looks like any literal value (boolean, null, or numeric)
	if IsBooleanOrNullLiteral(value) || isNumericLike(value) {
		return false
	}

	// Check for colon (always structural)
	if strings.Contains(value, ":") {
		return false
	}

	// Check for quotes and backslash (always need escaping)
	if strings.Contains(value, "\"") || strings.Contains(value, "\\") {
		return false
	}

	// Check for brackets and braces (always structural)
	if strings.ContainsAny(value, "[]{}") {
		return false
	}

	// Check for control characters (newline, carriage return, tab - always need quoting/escaping)
	if strings.ContainsAny(value, "\n\r\t") {
		return false
	}

	// Check for the active delimiter
	if strings.ContainsRune(value, rune(delimiter)) {
		return false
	}

	// Check for hyphen at start (list marker)
	if strings.HasPrefix(value, string(ListItemMarker)) {
		return false
	}

	return true
}

// isNumericLike checks if a string looks like a number.
// Match numbers like 42, -3.14, 1e-6, 05, etc.
func isNumericLike(value string) bool {
	// Try parsing as float - if it succeeds, it looks numeric
	_, err := strconv.ParseFloat(value, 64)
	if err == nil {
		return true
	}

	// Check for numbers with leading zeros like "05"
	if len(value) > 1 && value[0] == '0' && isDigit(rune(value[1])) {
		return true
	}

	return false
}

// isLetter checks if a rune is a letter (a-z, A-Z)
func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isDigit checks if a rune is a digit (0-9)
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
