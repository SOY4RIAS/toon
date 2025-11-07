package shared

import (
	"fmt"
	"strings"
)

const (
	Backslash      = '\\'
	DoubleQuote    = '"'
	Newline        = '\n'
	Tab            = '\t'
	CarriageReturn = '\r'
)

// EscapeString escapes special characters in a string for encoding.
// Handles backslashes, quotes, newlines, carriage returns, and tabs.
func EscapeString(value string) string {
	var result strings.Builder
	result.Grow(len(value) + 10) // Pre-allocate with some extra space

	for _, ch := range value {
		switch ch {
		case Backslash:
			result.WriteRune(Backslash)
			result.WriteRune(Backslash)
		case DoubleQuote:
			result.WriteRune(Backslash)
			result.WriteRune(DoubleQuote)
		case Newline:
			result.WriteRune(Backslash)
			result.WriteRune('n')
		case CarriageReturn:
			result.WriteRune(Backslash)
			result.WriteRune('r')
		case Tab:
			result.WriteRune(Backslash)
			result.WriteRune('t')
		default:
			result.WriteRune(ch)
		}
	}

	return result.String()
}

// UnescapeString processes escape sequences in a string.
// Handles \n, \t, \r, \\, and \" escape sequences.
func UnescapeString(value string) (string, error) {
	var result strings.Builder
	result.Grow(len(value))

	runes := []rune(value)
	i := 0

	for i < len(runes) {
		if runes[i] == Backslash {
			if i+1 >= len(runes) {
				return "", fmt.Errorf("invalid escape sequence: backslash at end of string")
			}

			next := runes[i+1]
			switch next {
			case 'n':
				result.WriteRune(Newline)
				i += 2
			case 't':
				result.WriteRune(Tab)
				i += 2
			case 'r':
				result.WriteRune(CarriageReturn)
				i += 2
			case Backslash:
				result.WriteRune(Backslash)
				i += 2
			case DoubleQuote:
				result.WriteRune(DoubleQuote)
				i += 2
			default:
				return "", fmt.Errorf("invalid escape sequence: \\%c", next)
			}
		} else {
			result.WriteRune(runes[i])
			i++
		}
	}

	return result.String(), nil
}

// FindClosingQuote finds the index of the closing double quote in a string,
// accounting for escape sequences.
// Returns the index of the closing quote, or -1 if not found.
func FindClosingQuote(content string, start int) int {
	runes := []rune(content)
	i := start + 1

	for i < len(runes) {
		if runes[i] == Backslash && i+1 < len(runes) {
			// Skip escaped character
			i += 2
			continue
		}
		if runes[i] == DoubleQuote {
			return i
		}
		i++
	}

	return -1 // Not found
}

// FindUnquotedChar finds the index of a specific character outside of quoted sections.
// Returns the index of the character, or -1 if not found outside quotes.
func FindUnquotedChar(content string, char rune, start int) int {
	runes := []rune(content)
	inQuotes := false
	i := start

	for i < len(runes) {
		if runes[i] == Backslash && i+1 < len(runes) && inQuotes {
			// Skip escaped character
			i += 2
			continue
		}

		if runes[i] == DoubleQuote {
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
