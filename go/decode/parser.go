package decode

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/toon-format/toon-go"
	"github.com/toon-format/toon-go/shared"
)

// ParseArrayHeaderLine parses an array header line and returns the header info and inline values
func ParseArrayHeaderLine(content string, defaultDelimiter toon.Delimiter) (*toon.ArrayHeaderInfo, string, error) {
	trimmed := strings.TrimLeft(content, " ")

	// Find the bracket segment, accounting for quoted keys that may contain brackets
	bracketStart := -1

	// For quoted keys, find bracket after closing quote (not inside the quoted string)
	if strings.HasPrefix(trimmed, "\"") {
		closingQuoteIndex := shared.FindClosingQuote(trimmed, 0)
		if closingQuoteIndex == -1 {
			return nil, "", nil
		}

		afterQuote := trimmed[closingQuoteIndex+1:]
		if !strings.HasPrefix(afterQuote, "[") {
			return nil, "", nil
		}

		// Calculate position in original content and find bracket after the quoted key
		leadingWhitespace := len(content) - len(trimmed)
		keyEndIndex := leadingWhitespace + closingQuoteIndex + 1
		bracketStart = strings.Index(content[keyEndIndex:], "[")
		if bracketStart != -1 {
			bracketStart += keyEndIndex
		}
	} else {
		// Unquoted key - find first bracket
		bracketStart = strings.Index(content, "[")
	}

	if bracketStart == -1 {
		return nil, "", nil
	}

	bracketEnd := strings.Index(content[bracketStart:], "]")
	if bracketEnd == -1 {
		return nil, "", nil
	}
	bracketEnd += bracketStart

	// Find the colon that comes after all brackets and braces
	colonIndex := bracketEnd + 1
	braceEnd := colonIndex

	// Check for fields segment (braces come after bracket)
	braceStart := strings.Index(content[bracketEnd:], "{")
	if braceStart != -1 {
		braceStart += bracketEnd
		colonPos := strings.Index(content[bracketEnd:], ":")
		if colonPos != -1 {
			colonPos += bracketEnd
			if braceStart < colonPos {
				foundBraceEnd := strings.Index(content[braceStart:], "}")
				if foundBraceEnd != -1 {
					braceEnd = braceStart + foundBraceEnd + 1
				}
			}
		}
	}

	// Now find colon after brackets and braces
	maxPos := bracketEnd
	if braceEnd > maxPos {
		maxPos = braceEnd
	}
	colonPos := strings.Index(content[maxPos:], ":")
	if colonPos == -1 {
		return nil, "", nil
	}
	colonIndex = maxPos + colonPos

	// Extract and parse the key (might be quoted)
	var key string
	if bracketStart > 0 {
		rawKey := strings.TrimSpace(content[:bracketStart])
		if strings.HasPrefix(rawKey, "\"") {
			parsed, err := ParseStringLiteral(rawKey)
			if err != nil {
				return nil, "", err
			}
			key = parsed
		} else {
			key = rawKey
		}
	}

	afterColon := strings.TrimSpace(content[colonIndex+1:])

	bracketContent := content[bracketStart+1 : bracketEnd]

	// Try to parse bracket segment
	length, delimiter, hasLengthMarker, err := ParseBracketSegment(bracketContent, defaultDelimiter)
	if err != nil {
		return nil, "", nil
	}

	// Check for fields segment
	var fields []string
	braceStart = strings.Index(content[bracketEnd:], "{")
	if braceStart != -1 {
		braceStart += bracketEnd
		if braceStart < colonIndex {
			foundBraceEnd := strings.Index(content[braceStart:], "}")
			if foundBraceEnd != -1 {
				foundBraceEnd += braceStart
				if foundBraceEnd < colonIndex {
					fieldsContent := content[braceStart+1 : foundBraceEnd]
					parsedFields, err := ParseDelimitedValues(fieldsContent, delimiter)
					if err != nil {
						return nil, "", err
					}
					fields = make([]string, len(parsedFields))
					for i, field := range parsedFields {
						parsed, err := ParseStringLiteral(strings.TrimSpace(field))
						if err != nil {
							return nil, "", err
						}
						fields[i] = parsed
					}
				}
			}
		}
	}

	header := &toon.ArrayHeaderInfo{
		Key:             key,
		Length:          length,
		Delimiter:       delimiter,
		Fields:          fields,
		HasLengthMarker: hasLengthMarker,
	}

	return header, afterColon, nil
}

// ParseBracketSegment parses a bracket segment like "5", "#5", "5\t", "#5|"
func ParseBracketSegment(seg string, defaultDelimiter toon.Delimiter) (length int, delimiter toon.Delimiter, hasLengthMarker bool, err error) {
	hasLengthMarker = false
	content := seg
	delimiter = defaultDelimiter

	// Check for length marker
	if strings.HasPrefix(content, "#") {
		hasLengthMarker = true
		content = content[1:]
	}

	// Check for delimiter suffix
	if strings.HasSuffix(content, "\t") {
		delimiter = toon.DelimiterTab
		content = content[:len(content)-1]
	} else if strings.HasSuffix(content, "|") {
		delimiter = toon.DelimiterPipe
		content = content[:len(content)-1]
	}

	length, err = strconv.Atoi(content)
	if err != nil {
		return 0, defaultDelimiter, false, fmt.Errorf("invalid array length: %s", seg)
	}

	return length, delimiter, hasLengthMarker, nil
}

// ParseDelimitedValues parses a delimited string into values
func ParseDelimitedValues(input string, delimiter toon.Delimiter) ([]string, error) {
	values := []string{}
	current := strings.Builder{}
	inQuotes := false

	runes := []rune(input)
	i := 0

	for i < len(runes) {
		char := runes[i]

		if char == '\\' && i+1 < len(runes) && inQuotes {
			// Escape sequence in quoted string
			current.WriteRune(char)
			current.WriteRune(runes[i+1])
			i += 2
			continue
		}

		if char == '"' {
			inQuotes = !inQuotes
			current.WriteRune(char)
			i++
			continue
		}

		if char == rune(delimiter) && !inQuotes {
			values = append(values, strings.TrimSpace(current.String()))
			current.Reset()
			i++
			continue
		}

		current.WriteRune(char)
		i++
	}

	// Add last value
	if current.Len() > 0 || len(values) > 0 {
		values = append(values, strings.TrimSpace(current.String()))
	}

	return values, nil
}

// MapRowValuesToPrimitives converts string values to primitives
func MapRowValuesToPrimitives(values []string) ([]toon.JsonPrimitive, error) {
	primitives := make([]toon.JsonPrimitive, len(values))
	for i, v := range values {
		prim, err := ParsePrimitiveToken(v)
		if err != nil {
			return nil, err
		}
		primitives[i] = prim
	}
	return primitives, nil
}

// ParsePrimitiveToken parses a token into a primitive value
func ParsePrimitiveToken(token string) (toon.JsonPrimitive, error) {
	trimmed := strings.TrimSpace(token)

	// Empty token
	if trimmed == "" {
		return "", nil
	}

	// Quoted string (if starts with quote, it MUST be properly quoted)
	if strings.HasPrefix(trimmed, "\"") {
		return ParseStringLiteral(trimmed)
	}

	// Boolean or null literals
	if shared.IsBooleanOrNullLiteral(trimmed) {
		if trimmed == shared.TrueLiteral {
			return true, nil
		}
		if trimmed == shared.FalseLiteral {
			return false, nil
		}
		if trimmed == shared.NullLiteral {
			return nil, nil
		}
	}

	// Numeric literal
	if shared.IsNumericLiteral(trimmed) {
		parsedNumber, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return nil, err
		}
		// Normalize negative zero to positive zero
		if parsedNumber == 0 && math.Signbit(parsedNumber) {
			return 0.0, nil
		}
		return parsedNumber, nil
	}

	// Unquoted string
	return trimmed, nil
}

// ParseStringLiteral parses a string literal (quoted or unquoted)
func ParseStringLiteral(token string) (string, error) {
	trimmedToken := strings.TrimSpace(token)

	if strings.HasPrefix(trimmedToken, "\"") {
		// Find the closing quote, accounting for escaped quotes
		closingQuoteIndex := shared.FindClosingQuote(trimmedToken, 0)

		if closingQuoteIndex == -1 {
			return "", fmt.Errorf("unterminated string: missing closing quote")
		}

		if closingQuoteIndex != len([]rune(trimmedToken))-1 {
			return "", fmt.Errorf("unexpected characters after closing quote")
		}

		content := trimmedToken[1:closingQuoteIndex]
		return shared.UnescapeString(content)
	}

	return trimmedToken, nil
}

// ParseUnquotedKey parses an unquoted key from content starting at the given position
func ParseUnquotedKey(content string, start int) (key string, end int, err error) {
	runes := []rune(content)
	end = start

	for end < len(runes) && runes[end] != ':' {
		end++
	}

	// Validate that a colon was found
	if end >= len(runes) || runes[end] != ':' {
		return "", 0, fmt.Errorf("missing colon after key")
	}

	key = strings.TrimSpace(string(runes[start:end]))

	// Skip the colon
	end++

	return key, end, nil
}

// ParseQuotedKey parses a quoted key from content starting at the given position
func ParseQuotedKey(content string, start int) (key string, end int, err error) {
	runes := []rune(content)

	// Find the closing quote, accounting for escaped quotes
	closingQuoteIndex := shared.FindClosingQuote(content, start)

	if closingQuoteIndex == -1 {
		return "", 0, fmt.Errorf("unterminated quoted key")
	}

	// Extract and unescape the key content
	keyContent := string(runes[start+1 : closingQuoteIndex])
	key, err = shared.UnescapeString(keyContent)
	if err != nil {
		return "", 0, err
	}

	end = closingQuoteIndex + 1

	// Validate and skip colon after quoted key
	if end >= len(runes) || runes[end] != ':' {
		return "", 0, fmt.Errorf("missing colon after key")
	}
	end++

	return key, end, nil
}

// ParseKeyToken parses a key token (quoted or unquoted)
func ParseKeyToken(content string, start int) (key string, end int, err error) {
	runes := []rune(content)
	if start >= len(runes) {
		return "", 0, fmt.Errorf("invalid start position")
	}

	if runes[start] == '"' {
		return ParseQuotedKey(content, start)
	}
	return ParseUnquotedKey(content, start)
}

// IsArrayHeaderAfterHyphen checks if content after hyphen is an array header
func IsArrayHeaderAfterHyphen(content string) bool {
	trimmed := strings.TrimSpace(content)
	return strings.HasPrefix(trimmed, "[") && shared.FindUnquotedChar(content, ':', 0) != -1
}

// IsObjectFirstFieldAfterHyphen checks if content after hyphen is an object field
func IsObjectFirstFieldAfterHyphen(content string) bool {
	return shared.FindUnquotedChar(content, ':', 0) != -1
}
