package toon

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ArrayHeaderResult contains the result of parsing an array header line
type ArrayHeaderResult struct {
	Header       ArrayHeaderInfo
	InlineValues string
}

// ParseArrayHeaderLine parses an array header line and returns header info and inline values.
// Returns nil if the line is not a valid array header.
func ParseArrayHeaderLine(content string, defaultDelimiter Delimiter) *ArrayHeaderResult {
	trimmed := strings.TrimLeft(content, " ")

	// Find the bracket segment, accounting for quoted keys that may contain brackets
	bracketStart := -1

	// For quoted keys, find bracket after closing quote (not inside the quoted string)
	if strings.HasPrefix(trimmed, string(DoubleQuote)) {
		closingQuoteIndex := FindClosingQuote(trimmed, 0)
		if closingQuoteIndex == -1 {
			return nil
		}

		afterQuote := trimmed[closingQuoteIndex+1:]
		if !strings.HasPrefix(afterQuote, string(OpenBracket)) {
			return nil
		}

		// Calculate position in original content and find bracket after the quoted key
		leadingWhitespace := len(content) - len(trimmed)
		keyEndIndex := leadingWhitespace + closingQuoteIndex + 1
		bracketStart = strings.IndexByte(content[keyEndIndex:], OpenBracket)
		if bracketStart != -1 {
			bracketStart += keyEndIndex
		}
	} else {
		// Unquoted key - find first bracket
		bracketStart = strings.IndexByte(content, OpenBracket)
	}

	if bracketStart == -1 {
		return nil
	}

	bracketEnd := strings.IndexByte(content[bracketStart:], CloseBracket)
	if bracketEnd == -1 {
		return nil
	}
	bracketEnd += bracketStart

	// Find the colon that comes after all brackets and braces
	colonIndex := bracketEnd + 1
	braceEnd := colonIndex

	// Check for fields segment (braces come after bracket)
	braceStart := strings.IndexByte(content[bracketEnd:], OpenBrace)
	if braceStart != -1 {
		braceStart += bracketEnd
		colonPos := strings.IndexByte(content[bracketEnd:], Colon)
		if colonPos != -1 {
			colonPos += bracketEnd
			if braceStart < colonPos {
				foundBraceEnd := strings.IndexByte(content[braceStart:], CloseBrace)
				if foundBraceEnd != -1 {
					braceEnd = braceStart + foundBraceEnd + 1
				}
			}
		}
	}

	// Now find colon after brackets and braces
	maxEnd := bracketEnd
	if braceEnd > maxEnd {
		maxEnd = braceEnd
	}
	colonIndex = strings.IndexByte(content[maxEnd:], Colon)
	if colonIndex == -1 {
		return nil
	}
	colonIndex += maxEnd

	// Extract and parse the key (might be quoted)
	var key *string
	if bracketStart > 0 {
		rawKey := strings.TrimSpace(content[:bracketStart])
		parsedKey, err := ParseStringLiteral(rawKey)
		if err != nil {
			return nil
		}
		key = &parsedKey
	}

	afterColon := strings.TrimSpace(content[colonIndex+1:])

	bracketContent := content[bracketStart+1 : bracketEnd]

	// Try to parse bracket segment
	parsedBracket, err := ParseBracketSegment(bracketContent, defaultDelimiter)
	if err != nil {
		return nil
	}

	length := parsedBracket.Length
	delimiter := parsedBracket.Delimiter
	hasLengthMarker := parsedBracket.HasLengthMarker

	// Check for fields segment
	var fields []string
	if braceStart != -1 && braceStart < colonIndex {
		foundBraceEnd := strings.IndexByte(content[braceStart:], CloseBrace)
		if foundBraceEnd != -1 {
			foundBraceEnd += braceStart
			if foundBraceEnd < colonIndex {
				fieldsContent := content[braceStart+1 : foundBraceEnd]
				values := ParseDelimitedValues(fieldsContent, delimiter)
				fields = make([]string, len(values))
				for i, field := range values {
					parsed, err := ParseStringLiteral(strings.TrimSpace(field))
					if err != nil {
						return nil
					}
					fields[i] = parsed
				}
			}
		}
	}

	result := &ArrayHeaderResult{
		Header: ArrayHeaderInfo{
			Key:             key,
			Length:          length,
			Delimiter:       delimiter,
			Fields:          fields,
			HasLengthMarker: hasLengthMarker,
		},
	}

	if afterColon != "" {
		result.InlineValues = afterColon
	}

	return result
}

// BracketSegmentResult contains the result of parsing a bracket segment
type BracketSegmentResult struct {
	Length          int
	Delimiter       Delimiter
	HasLengthMarker bool
}

// ParseBracketSegment parses the content inside brackets [...]
func ParseBracketSegment(seg string, defaultDelimiter Delimiter) (*BracketSegmentResult, error) {
	hasLengthMarker := false
	content := seg

	// Check for length marker
	if strings.HasPrefix(content, string(Hash)) {
		hasLengthMarker = true
		content = content[1:]
	}

	// Check for delimiter suffix
	delimiter := defaultDelimiter
	if strings.HasSuffix(content, string(Tab)) {
		delimiter = DelimiterTab
		content = content[:len(content)-1]
	} else if strings.HasSuffix(content, string(Pipe)) {
		delimiter = DelimiterPipe
		content = content[:len(content)-1]
	}

	length, err := strconv.Atoi(content)
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %s", seg)
	}

	return &BracketSegmentResult{
		Length:          length,
		Delimiter:       delimiter,
		HasLengthMarker: hasLengthMarker,
	}, nil
}

// ParseDelimitedValues splits a string by delimiter, respecting quoted strings
func ParseDelimitedValues(input string, delimiter Delimiter) []string {
	values := []string{}
	current := ""
	inQuotes := false
	i := 0

	for i < len(input) {
		char := input[i]

		if char == Backslash && i+1 < len(input) && inQuotes {
			// Escape sequence in quoted string
			current += string(char) + string(input[i+1])
			i += 2
			continue
		}

		if char == DoubleQuote {
			inQuotes = !inQuotes
			current += string(char)
			i++
			continue
		}

		if char == byte(delimiter) && !inQuotes {
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

// MapRowValuesToPrimitives converts string values to primitives
func MapRowValuesToPrimitives(values []string) ([]interface{}, error) {
	result := make([]interface{}, len(values))
	for i, v := range values {
		primitive, err := ParsePrimitiveToken(v)
		if err != nil {
			return nil, err
		}
		result[i] = primitive
	}
	return result, nil
}

// ParsePrimitiveToken parses a token into a primitive value (string, number, boolean, or null)
func ParsePrimitiveToken(token string) (interface{}, error) {
	trimmed := strings.TrimSpace(token)

	// Empty token
	if trimmed == "" {
		return "", nil
	}

	// Quoted string (if starts with quote, it MUST be properly quoted)
	if strings.HasPrefix(trimmed, string(DoubleQuote)) {
		return ParseStringLiteral(trimmed)
	}

	// Boolean or null literals
	if IsBooleanOrNullLiteral(trimmed) {
		if trimmed == TrueLiteral {
			return true, nil
		}
		if trimmed == FalseLiteral {
			return false, nil
		}
		if trimmed == NullLiteral {
			return nil, nil
		}
	}

	// Numeric literal
	if IsNumericLiteral(trimmed) {
		parsedNumber, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return nil, err
		}
		// Normalize negative zero to positive zero
		if parsedNumber == 0 && math.Signbit(parsedNumber) {
			parsedNumber = 0
		}
		return parsedNumber, nil
	}

	// Unquoted string
	return trimmed, nil
}

// ParseStringLiteral parses a string literal, handling quoted and unquoted strings
func ParseStringLiteral(token string) (string, error) {
	trimmedToken := strings.TrimSpace(token)

	if strings.HasPrefix(trimmedToken, string(DoubleQuote)) {
		// Find the closing quote, accounting for escaped quotes
		closingQuoteIndex := FindClosingQuote(trimmedToken, 0)

		if closingQuoteIndex == -1 {
			return "", fmt.Errorf("unterminated string: missing closing quote")
		}

		if closingQuoteIndex != len(trimmedToken)-1 {
			return "", fmt.Errorf("unexpected characters after closing quote")
		}

		content := trimmedToken[1:closingQuoteIndex]
		return UnescapeString(content)
	}

	return trimmedToken, nil
}

// KeyTokenResult contains the result of parsing a key token
type KeyTokenResult struct {
	Key string
	End int
}

// ParseUnquotedKey parses an unquoted key from content starting at the given position
func ParseUnquotedKey(content string, start int) (*KeyTokenResult, error) {
	end := start
	for end < len(content) && content[end] != Colon {
		end++
	}

	// Validate that a colon was found
	if end >= len(content) || content[end] != Colon {
		return nil, fmt.Errorf("missing colon after key")
	}

	key := strings.TrimSpace(content[start:end])

	// Skip the colon
	end++

	return &KeyTokenResult{
		Key: key,
		End: end,
	}, nil
}

// ParseQuotedKey parses a quoted key from content starting at the given position
func ParseQuotedKey(content string, start int) (*KeyTokenResult, error) {
	// Find the closing quote, accounting for escaped quotes
	closingQuoteIndex := FindClosingQuote(content, start)

	if closingQuoteIndex == -1 {
		return nil, fmt.Errorf("unterminated quoted key")
	}

	// Extract and unescape the key content
	keyContent := content[start+1 : closingQuoteIndex]
	key, err := UnescapeString(keyContent)
	if err != nil {
		return nil, err
	}

	end := closingQuoteIndex + 1

	// Validate and skip colon after quoted key
	if end >= len(content) || content[end] != Colon {
		return nil, fmt.Errorf("missing colon after key")
	}
	end++

	return &KeyTokenResult{
		Key: key,
		End: end,
	}, nil
}

// ParseKeyToken parses a key token (quoted or unquoted) from content
func ParseKeyToken(content string, start int) (*KeyTokenResult, error) {
	if start < len(content) && content[start] == DoubleQuote {
		return ParseQuotedKey(content, start)
	}
	return ParseUnquotedKey(content, start)
}

// IsArrayHeaderAfterHyphen checks if content after hyphen is an array header
func IsArrayHeaderAfterHyphen(content string) bool {
	trimmed := strings.TrimSpace(content)
	return strings.HasPrefix(trimmed, string(OpenBracket)) && FindUnquotedChar(trimmed, Colon, 0) != -1
}

// IsObjectFirstFieldAfterHyphen checks if content after hyphen is an object field
func IsObjectFirstFieldAfterHyphen(content string) bool {
	return FindUnquotedChar(content, Colon, 0) != -1
}
