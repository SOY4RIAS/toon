package shared

import (
	"strconv"
	"strings"
)

// IsBooleanOrNullLiteral checks if a token is a boolean or null literal
// (true, false, null).
func IsBooleanOrNullLiteral(token string) bool {
	return token == "true" || token == "false" || token == "null"
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
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// IsValidUnquotedKey checks if a string can be used as an unquoted key.
// Keys must start with a letter or underscore, followed by letters, digits,
// underscores, or dots.
func IsValidUnquotedKey(key string) bool {
	if key == "" {
		return false
	}

	runes := []rune(key)
	first := runes[0]

	// First character must be letter or underscore
	if !isLetter(first) && first != '_' {
		return false
	}

	// Remaining characters must be letter, digit, underscore, or dot
	for i := 1; i < len(runes); i++ {
		ch := runes[i]
		if !isLetter(ch) && !isDigit(ch) && ch != '_' && ch != '.' {
			return false
		}
	}

	return true
}

// HasLeadingOrTrailingSpace checks if a string has leading or trailing whitespace
func HasLeadingOrTrailingSpace(s string) bool {
	if s == "" {
		return false
	}
	return strings.TrimSpace(s) != s
}

// ContainsAny checks if a string contains any of the specified runes
func ContainsAny(s string, chars string) bool {
	return strings.ContainsAny(s, chars)
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}
