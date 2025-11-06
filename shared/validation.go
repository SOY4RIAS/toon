package shared

import (
	"regexp"
	"strings"
)

var (
	// validUnquotedKeyRegex matches valid unquoted keys.
	// Must start with a letter or underscore, followed by letters, digits, underscores, or dots.
	validUnquotedKeyRegex = regexp.MustCompile(`^[A-Za-z_][\w.]*$`)

	// numericLikeRegex matches strings that look like numbers.
	numericLikeRegex = regexp.MustCompile(`^-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?$`)

	// leadingZeroRegex matches numbers with leading zeros (like "05").
	leadingZeroRegex = regexp.MustCompile(`^0\d+$`)

	// structuralCharsRegex matches structural characters (brackets and braces).
	structuralCharsRegex = regexp.MustCompile(`[[\]{}]`)

	// controlCharsRegex matches control characters (newline, carriage return, tab).
	controlCharsRegex = regexp.MustCompile(`[\n\r\t]`)
)

// IsValidUnquotedKey checks if a key can be used without quotes.
// Valid unquoted keys must start with a letter or underscore,
// followed by letters, digits, underscores, or dots.
func IsValidUnquotedKey(key string) bool {
	return validUnquotedKeyRegex.MatchString(key)
}

// IsSafeUnquoted determines if a string value can be safely encoded without quotes.
// A string needs quoting if it:
//   - Is empty
//   - Has leading or trailing whitespace
//   - Could be confused with a literal (boolean, null, number)
//   - Contains structural characters (colons, brackets, braces)
//   - Contains quotes or backslashes (need escaping)
//   - Contains control characters (newlines, tabs, etc.)
//   - Contains the active delimiter
//   - Starts with a list marker (hyphen + space)
func IsSafeUnquoted(value string, delimiter rune) bool {
	if value == "" {
		return false
	}

	// Check for leading or trailing whitespace
	if value != strings.TrimSpace(value) {
		return false
	}

	// Check if it looks like any literal value (boolean, null, or numeric)
	if IsBooleanOrNullLiteral(value) || isNumericLike(value) {
		return false
	}

	// Check for colon (always structural)
	if strings.ContainsRune(value, ':') {
		return false
	}

	// Check for quotes and backslash (always need escaping)
	if strings.ContainsRune(value, '"') || strings.ContainsRune(value, '\\') {
		return false
	}

	// Check for brackets and braces (always structural)
	if structuralCharsRegex.MatchString(value) {
		return false
	}

	// Check for control characters (newline, carriage return, tab - always need quoting/escaping)
	if controlCharsRegex.MatchString(value) {
		return false
	}

	// Check for the active delimiter
	if strings.ContainsRune(value, delimiter) {
		return false
	}

	// Check for hyphen at start (list marker "- ")
	if strings.HasPrefix(value, "-") {
		return false
	}

	return true
}

// isNumericLike checks if a string looks like a number.
// Matches numbers like 42, -3.14, 1e-6, 05, etc.
func isNumericLike(value string) bool {
	return numericLikeRegex.MatchString(value) || leadingZeroRegex.MatchString(value)
}
