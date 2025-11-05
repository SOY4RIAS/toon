package shared

import (
	"fmt"
	"strings"
)

// EscapeString escapes special characters in a string for encoding.
// Handles backslashes, quotes, newlines, carriage returns, and tabs.
func EscapeString(value string) string {
	var result strings.Builder
	result.Grow(len(value))

	for _, ch := range value {
		switch ch {
		case '\\':
			result.WriteString("\\\\")
		case '"':
			result.WriteString("\\\"")
		case '\n':
			result.WriteString("\\n")
		case '\r':
			result.WriteString("\\r")
		case '\t':
			result.WriteString("\\t")
		default:
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// UnescapeString unescapes a string by processing escape sequences.
// Handles \n, \t, \r, \\, and \" escape sequences.
func UnescapeString(value string) (string, error) {
	var result strings.Builder
	result.Grow(len(value))

	runes := []rune(value)
	i := 0

	for i < len(runes) {
		if runes[i] == '\\' {
			if i+1 >= len(runes) {
				return "", fmt.Errorf("invalid escape sequence: backslash at end of string")
			}

			next := runes[i+1]
			switch next {
			case 'n':
				result.WriteRune('\n')
				i += 2
			case 't':
				result.WriteRune('\t')
				i += 2
			case 'r':
				result.WriteRune('\r')
				i += 2
			case '\\':
				result.WriteRune('\\')
				i += 2
			case '"':
				result.WriteRune('"')
				i += 2
			default:
				return "", fmt.Errorf("invalid escape sequence: \\%c", next)
			}
			continue
		}

		result.WriteRune(runes[i])
		i++
	}

	return result.String(), nil
}

// FindClosingQuote finds the index of the closing double quote in a string,
// accounting for escape sequences.
//
// Parameters:
//   - content: The string to search in
//   - start: The index of the opening quote
//
// Returns the index of the closing quote, or -1 if not found.
func FindClosingQuote(content string, start int) int {
	runes := []rune(content)
	i := start + 1

	for i < len(runes) {
		if runes[i] == '\\' && i+1 < len(runes) {
			// Skip escaped character
			i += 2
			continue
		}
		if runes[i] == '"' {
			return i
		}
		i++
	}

	return -1 // Not found
}

// FindUnquotedChar finds the index of a specific character outside of quoted sections.
//
// Parameters:
//   - content: The string to search in
//   - char: The character to look for
//   - start: Starting index (defaults to 0)
//
// Returns the index of the character, or -1 if not found outside quotes.
func FindUnquotedChar(content string, char rune, start int) int {
	runes := []rune(content)
	inQuotes := false
	i := start

	for i < len(runes) {
		if runes[i] == '\\' && i+1 < len(runes) && inQuotes {
			// Skip escaped character
			i += 2
			continue
		}

		if runes[i] == '"' {
			inQuotes = !inQuotes
			i++
			continue
		}

		if runes[i] == char && !inQuotes {
			return i
		}

		i++
	}

	return -1
}
