package toon

import (
	"fmt"
	"strings"
)

// EscapeString escapes special characters in a string for encoding.
// Handles backslashes, quotes, newlines, carriage returns, and tabs.
func EscapeString(value string) string {
	result := value
	result = strings.ReplaceAll(result, "\\", "\\\\")
	result = strings.ReplaceAll(result, "\"", "\\\"")
	result = strings.ReplaceAll(result, "\n", "\\n")
	result = strings.ReplaceAll(result, "\r", "\\r")
	result = strings.ReplaceAll(result, "\t", "\\t")
	return result
}

// UnescapeString unescapes a string by processing escape sequences.
// Handles \n, \t, \r, \\, and \" escape sequences.
func UnescapeString(value string) (string, error) {
	var result strings.Builder
	i := 0

	for i < len(value) {
		if value[i] == Backslash {
			if i+1 >= len(value) {
				return "", fmt.Errorf("invalid escape sequence: backslash at end of string")
			}

			next := value[i+1]
			switch next {
			case 'n':
				result.WriteByte('\n')
				i += 2
				continue
			case 't':
				result.WriteByte('\t')
				i += 2
				continue
			case 'r':
				result.WriteByte('\r')
				i += 2
				continue
			case Backslash:
				result.WriteByte(Backslash)
				i += 2
				continue
			case DoubleQuote:
				result.WriteByte(DoubleQuote)
				i += 2
				continue
			default:
				return "", fmt.Errorf("invalid escape sequence: \\%c", next)
			}
		}

		result.WriteByte(value[i])
		i++
	}

	return result.String(), nil
}

// FindClosingQuote finds the index of the closing double quote in a string,
// accounting for escape sequences.
// Returns the index of the closing quote, or -1 if not found.
func FindClosingQuote(content string, start int) int {
	i := start + 1
	for i < len(content) {
		if content[i] == Backslash && i+1 < len(content) {
			// Skip escaped character
			i += 2
			continue
		}
		if content[i] == DoubleQuote {
			return i
		}
		i++
	}
	return -1 // Not found
}

// FindUnquotedChar finds the index of a specific character outside of quoted sections.
// Returns the index of the character, or -1 if not found outside quotes.
func FindUnquotedChar(content string, char byte, start int) int {
	inQuotes := false
	i := start

	for i < len(content) {
		if content[i] == Backslash && i+1 < len(content) && inQuotes {
			// Skip escaped character
			i += 2
			continue
		}

		if content[i] == DoubleQuote {
			inQuotes = !inQuotes
			i++
			continue
		}

		if content[i] == char && !inQuotes {
			return i
		}

		i++
	}

	return -1
}
