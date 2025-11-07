package toon

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/soy4rias/toongo/shared"
)

// ParseArrayHeaderLine parses an array header line
// Returns header info and optional inline values after colon
func ParseArrayHeaderLine(content string, defaultDelimiter Delimiter) (*ArrayHeaderInfo, string, error) {
	trimmed := strings.TrimLeft(content, " ")

	// Find the bracket segment, accounting for quoted keys
	bracketStart := -1

	if strings.HasPrefix(trimmed, "\"") {
		closingQuoteIndex := shared.FindClosingQuote(trimmed, 0)
		if closingQuoteIndex == -1 {
			return nil, "", nil
		}

		afterQuote := trimmed[closingQuoteIndex+1:]
		if !strings.HasPrefix(afterQuote, "[") {
			return nil, "", nil
		}

		leadingWhitespace := len(content) - len(trimmed)
		keyEndIndex := leadingWhitespace + closingQuoteIndex + 1
		bracketStart = strings.Index(content[keyEndIndex:], "[")
		if bracketStart != -1 {
			bracketStart += keyEndIndex
		}
	} else {
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

	// Find the colon that comes after brackets and braces
	colonIndex := bracketEnd + 1
	braceEnd := colonIndex

	// Check for fields segment (braces come after bracket)
	braceStart := strings.Index(content[bracketEnd:], "{")
	if braceStart != -1 {
		braceStart += bracketEnd
		colonAfterBracket := strings.Index(content[bracketEnd:], ":")
		if colonAfterBracket != -1 && braceStart < bracketEnd+colonAfterBracket {
			foundBraceEnd := strings.Index(content[braceStart:], "}")
			if foundBraceEnd != -1 {
				braceEnd = braceStart + foundBraceEnd + 1
			}
		}
	}

	// Find colon after brackets and braces
	maxPos := bracketEnd
	if braceEnd > maxPos {
		maxPos = braceEnd
	}
	colonIndex = strings.Index(content[maxPos:], ":")
	if colonIndex == -1 {
		return nil, "", nil
	}
	colonIndex += maxPos

	// Extract and parse the key (might be quoted)
	var key string
	if bracketStart > 0 {
		rawKey := strings.TrimSpace(content[:bracketStart])
		if strings.HasPrefix(rawKey, "\"") {
			key = parseStringLiteral(rawKey)
		} else {
			key = rawKey
		}
	}

	afterColon := strings.TrimSpace(content[colonIndex+1:])
	bracketContent := content[bracketStart+1 : bracketEnd]

	// Parse bracket segment
	length, delimiter, hasLengthMarker, err := parseBracketSegment(bracketContent, defaultDelimiter)
	if err != nil {
		return nil, "", err
	}

	// Check for fields segment
	var fields []string
	if braceStart != -1 && braceStart < colonIndex {
		foundBraceEnd := strings.Index(content[braceStart:], "}")
		if foundBraceEnd != -1 {
			foundBraceEnd += braceStart
			if foundBraceEnd < colonIndex {
				fieldsContent := content[braceStart+1 : foundBraceEnd]
				fieldValues := parseDelimitedValues(fieldsContent, delimiter)
				for _, field := range fieldValues {
					fields = append(fields, parseStringLiteral(strings.TrimSpace(field)))
				}
			}
		}
	}

	header := &ArrayHeaderInfo{
		Key:             key,
		Length:          length,
		Delimiter:       delimiter,
		Fields:          fields,
		HasLengthMarker: hasLengthMarker,
	}

	return header, afterColon, nil
}

func parseBracketSegment(seg string, defaultDelimiter Delimiter) (int, Delimiter, bool, error) {
	hasLengthMarker := false
	content := seg

	// Check for length marker
	if strings.HasPrefix(content, "#") {
		hasLengthMarker = true
		content = content[1:]
	}

	// Check for delimiter suffix
	delimiter := defaultDelimiter
	if strings.HasSuffix(content, "\t") {
		delimiter = DelimiterTab
		content = content[:len(content)-1]
	} else if strings.HasSuffix(content, "|") {
		delimiter = DelimiterPipe
		content = content[:len(content)-1]
	}

	length, err := strconv.Atoi(content)
	if err != nil {
		return 0, 0, false, fmt.Errorf("invalid array length: %s", seg)
	}

	return length, delimiter, hasLengthMarker, nil
}

// parseDelimitedValues parses a delimited string into values
func parseDelimitedValues(input string, delimiter Delimiter) []string {
	values := []string{}
	current := ""
	inQuotes := false
	i := 0

	delimChar := rune(delimiter)

	for i < len(input) {
		char := rune(input[i])

		if char == '\\' && i+1 < len(input) && inQuotes {
			// Escape sequence in quoted string
			current += string(char) + string(input[i+1])
			i += 2
			continue
		}

		if char == '"' {
			inQuotes = !inQuotes
			current += string(char)
			i++
			continue
		}

		if char == delimChar && !inQuotes {
			values = append(values, strings.TrimSpace(current))
			current = ""
			i++
			continue
		}

		current += string(char)
		i++
	}

	// Add last value
	if current != "" || len(values) > 0 {
		values = append(values, strings.TrimSpace(current))
	}

	return values
}

// mapRowValuesToPrimitives converts string values to primitives
func mapRowValuesToPrimitives(values []string) []JsonPrimitive {
	primitives := make([]JsonPrimitive, len(values))
	for i, v := range values {
		primitives[i] = parsePrimitiveToken(v)
	}
	return primitives
}

// parsePrimitiveToken parses a token into a primitive value
func parsePrimitiveToken(token string) JsonPrimitive {
	trimmed := strings.TrimSpace(token)

	// Empty token
	if trimmed == "" {
		return ""
	}

	// Quoted string
	if strings.HasPrefix(trimmed, "\"") {
		return parseStringLiteral(trimmed)
	}

	// Boolean or null literals
	if shared.IsBooleanOrNullLiteral(trimmed) {
		if trimmed == "true" {
			return true
		}
		if trimmed == "false" {
			return false
		}
		if trimmed == "null" {
			return nil
		}
	}

	// Numeric literal
	if shared.IsNumericLiteral(trimmed) {
		num, _ := strconv.ParseFloat(trimmed, 64)
		// Normalize negative zero to positive zero
		if num == 0 && strings.HasPrefix(trimmed, "-") {
			return 0.0
		}
		return num
	}

	// Unquoted string
	return trimmed
}

func parseStringLiteral(token string) string {
	trimmedToken := strings.TrimSpace(token)

	if strings.HasPrefix(trimmedToken, "\"") {
		closingQuoteIndex := shared.FindClosingQuote(trimmedToken, 0)

		if closingQuoteIndex == -1 {
			// This should probably be an error, but for compatibility return as-is
			panic("unterminated string: missing closing quote")
		}

		if closingQuoteIndex != len(trimmedToken)-1 {
			panic("unexpected characters after closing quote")
		}

		content := trimmedToken[1:closingQuoteIndex]
		unescaped, err := shared.UnescapeString(content)
		if err != nil {
			panic(err)
		}
		return unescaped
	}

	return trimmedToken
}

// parseKeyToken parses a key from content starting at the given position
func parseKeyToken(content string, start int) (string, int, error) {
	if start < len(content) && content[start] == '"' {
		return parseQuotedKey(content, start)
	}
	return parseUnquotedKey(content, start)
}

func parseQuotedKey(content string, start int) (string, int, error) {
	closingQuoteIndex := shared.FindClosingQuote(content, start)

	if closingQuoteIndex == -1 {
		return "", 0, fmt.Errorf("unterminated quoted key")
	}

	keyContent := content[start+1 : closingQuoteIndex]
	key, err := shared.UnescapeString(keyContent)
	if err != nil {
		return "", 0, err
	}
	end := closingQuoteIndex + 1

	// Validate and skip colon after quoted key
	if end >= len(content) || content[end] != ':' {
		return "", 0, fmt.Errorf("missing colon after key")
	}
	end++

	return key, end, nil
}

func parseUnquotedKey(content string, start int) (string, int, error) {
	end := start
	for end < len(content) && content[end] != ':' {
		end++
	}

	// Validate that a colon was found
	if end >= len(content) || content[end] != ':' {
		return "", 0, fmt.Errorf("missing colon after key")
	}

	key := strings.TrimSpace(content[start:end])

	// Skip the colon
	end++

	return key, end, nil
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
