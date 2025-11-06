package shared

import (
	"errors"
	"fmt"
	"strings"
)

// EscapeString escapes special characters in a string for encoding
// Handles backslashes, quotes, newlines, carriage returns, and tabs
func EscapeString(value string) string {
	value = strings.ReplaceAll(value, "\\", "\\\\")
	value = strings.ReplaceAll(value, "\"", "\\\"")
	value = strings.ReplaceAll(value, "\n", "\\n")
	value = strings.ReplaceAll(value, "\r", "\\r")
	value = strings.ReplaceAll(value, "\t", "\\t")
	return value
}

// UnescapeString unescapes a string by processing escape sequences
// Handles \n, \t, \r, \\, and \" escape sequences
func UnescapeString(value string) (string, error) {
	var result strings.Builder
	i := 0

	for i < len(value) {
		if value[i] == '\\' {
			if i+1 >= len(value) {
				return "", errors.New("invalid escape sequence: backslash at end of string")
			}

			next := value[i+1]
			switch next {
			case 'n':
				result.WriteByte('\n')
				i += 2
			case 't':
				result.WriteByte('\t')
				i += 2
			case 'r':
				result.WriteByte('\r')
				i += 2
			case '\\':
				result.WriteByte('\\')
				i += 2
			case '"':
				result.WriteByte('"')
				i += 2
			default:
				return "", fmt.Errorf("invalid escape sequence: \\%c", next)
			}
		} else {
			result.WriteByte(value[i])
			i++
		}
	}

	return result.String(), nil
}

// FindClosingQuote finds the index of the closing double quote in a string,
// accounting for escape sequences.
// Returns the index of the closing quote, or -1 if not found.
func FindClosingQuote(content string, start int) int {
	i := start + 1
	for i < len(content) {
		if content[i] == '\\' && i+1 < len(content) {
			// Skip escaped character
			i += 2
			continue
		}
		if content[i] == '"' {
			return i
		}
		i++
	}
	return -1 // Not found
}

// FindUnquotedChar finds the index of a specific character outside of quoted sections
// Returns the index of the character, or -1 if not found outside quotes
func FindUnquotedChar(content string, char rune, start int) int {
	inQuotes := false
	i := start

	for i < len(content) {
		if content[i] == '\\' && i+1 < len(content) && inQuotes {
			// Skip escaped character
			i += 2
			continue
		}

		if content[i] == '"' {
			inQuotes = !inQuotes
			i++
			continue
		}

		if rune(content[i]) == char && !inQuotes {
			return i
		}

		i++
	}

	return -1
}
