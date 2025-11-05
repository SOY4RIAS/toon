package decode

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/toon-format/toon-go"
	"github.com/toon-format/toon-go/shared"
)

// ArrayHeaderParseResult contains the result of parsing an array header line
type ArrayHeaderParseResult struct {
	Header       toon.ArrayHeaderInfo
	InlineValues string
}

// ParseArrayHeaderLine parses an array header line like "items[2]{id,name}: val1,val2"
// Returns nil if the line doesn't contain an array header
func ParseArrayHeaderLine(content string, defaultDelimiter toon.Delimiter) *ArrayHeaderParseResult {
	trimmed := strings.TrimLeft(content, " ")

	// Find the bracket segment, accounting for quoted keys that may contain brackets
	bracketStart := -1

	// For quoted keys, find bracket after closing quote (not inside the quoted string)
	if len(trimmed) > 0 && trimmed[0] == toon.DoubleQuote {
		closingQuoteIndex := shared.FindClosingQuote(trimmed, 0)
		if closingQuoteIndex == -1 {
			return nil
		}

		afterQuote := trimmed[closingQuoteIndex+1:]
		if len(afterQuote) == 0 || afterQuote[0] != byte(toon.OpenBracket) {
			return nil
		}

		// Calculate position in original content and find bracket after the quoted key
		leadingWhitespace := len(content) - len(trimmed)
		keyEndIndex := leadingWhitespace + closingQuoteIndex + 1
		bracketStart = strings.IndexRune(content[keyEndIndex:], toon.OpenBracket)
		if bracketStart != -1 {
			bracketStart += keyEndIndex
		}
	} else {
		// Unquoted key - find first bracket
		bracketStart = strings.IndexRune(content, toon.OpenBracket)
	}

	if bracketStart == -1 {
		return nil
	}

	bracketEnd := strings.IndexRune(content[bracketStart:], toon.CloseBracket)
	if bracketEnd == -1 {
		return nil
	}
	bracketEnd += bracketStart

	// Find the colon that comes after all brackets and braces
	colonIndex := bracketEnd + 1
	braceEnd := colonIndex

	// Check for fields segment (braces come after bracket)
	braceStart := strings.IndexRune(content[bracketEnd:], toon.OpenBrace)
	if braceStart != -1 {
		braceStart += bracketEnd
		colonPos := strings.IndexRune(content[bracketEnd:], toon.Colon)
		if colonPos != -1 {
			colonPos += bracketEnd
			if braceStart < colonPos {
				foundBraceEnd := strings.IndexRune(content[braceStart:], toon.CloseBrace)
				if foundBraceEnd != -1 {
					braceEnd = braceStart + foundBraceEnd + 1
				}
			}
		}
	}

	// Now find colon after brackets and braces
	maxPos := braceEnd
	if bracketEnd > maxPos {
		maxPos = bracketEnd
	}
	colonPos := strings.IndexRune(content[maxPos:], toon.Colon)
	if colonPos == -1 {
		return nil
	}
	colonIndex = maxPos + colonPos

	// Extract and parse the key (might be quoted)
	var key *string
	if bracketStart > 0 {
		rawKey := strings.TrimSpace(content[:bracketStart])
		if len(rawKey) > 0 {
			parsedKey := ParseStringLiteral(rawKey)
			key = &parsedKey
		}
	}

	afterColon := strings.TrimSpace(content[colonIndex+1:])

	bracketContent := content[bracketStart+1 : bracketEnd]

	// Try to parse bracket segment
	parsedBracket, err := ParseBracketSegment(bracketContent, defaultDelimiter)
	if err != nil {
		return nil
	}

	// Check for fields segment
	var fields []string
	if braceStart != -1 && braceStart < colonIndex {
		foundBraceEnd := strings.IndexRune(content[braceStart:], toon.CloseBrace)
		if foundBraceEnd != -1 {
			foundBraceEnd += braceStart
			if foundBraceEnd < colonIndex {
				fieldsContent := content[braceStart+1 : foundBraceEnd]
				fieldValues := ParseDelimitedValues(fieldsContent, parsedBracket.Delimiter)
				fields = make([]string, len(fieldValues))
				for i, field := range fieldValues {
					fields[i] = ParseStringLiteral(strings.TrimSpace(field))
				}
			}
		}
	}

	inlineValues := ""
	if afterColon != "" {
		inlineValues = afterColon
	}

	return &ArrayHeaderParseResult{
		Header: toon.ArrayHeaderInfo{
			Key:             key,
			Length:          parsedBracket.Length,
			Delimiter:       parsedBracket.Delimiter,
			Fields:          fields,
			HasLengthMarker: parsedBracket.HasLengthMarker,
		},
		InlineValues: inlineValues,
	}
}

// BracketSegmentResult contains the result of parsing a bracket segment
type BracketSegmentResult struct {
	Length          int
	Delimiter       toon.Delimiter
	HasLengthMarker bool
}

// ParseBracketSegment parses a bracket segment like "2" or "#3|" or "5	"
func ParseBracketSegment(seg string, defaultDelimiter toon.Delimiter) (*BracketSegmentResult, error) {
	hasLengthMarker := false
	content := seg

	// Check for length marker
	if len(content) > 0 && content[0] == byte(toon.Hash) {
		hasLengthMarker = true
		content = content[1:]
	}

	// Check for delimiter suffix
	delimiter := defaultDelimiter
	if len(content) > 0 {
		lastChar := rune(content[len(content)-1])
		if lastChar == toon.Tab {
			delimiter = toon.DelimiterTab
			content = content[:len(content)-1]
		} else if lastChar == toon.Pipe {
			delimiter = toon.DelimiterPipe
			content = content[:len(content)-1]
		}
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

// ParseDelimitedValues splits a string by delimiter, accounting for quoted strings
func ParseDelimitedValues(input string, delimiter toon.Delimiter) []string {
	values := []string{}
	current := ""
	inQuotes := false
	runes := []rune(input)
	i := 0

	for i < len(runes) {
		ch := runes[i]

		if ch == toon.Backslash && i+1 < len(runes) && inQuotes {
			// Escape sequence in quoted string
			current += string(ch) + string(runes[i+1])
			i += 2
			continue
		}

		if ch == toon.DoubleQuote {
			inQuotes = !inQuotes
			current += string(ch)
			i++
			continue
		}

		if ch == rune(delimiter) && !inQuotes {
			values = append(values, strings.TrimSpace(current))
			current = ""
			i++
			continue
		}

		current += string(ch)
		i++
	}

	// Add last value
	if current != "" || len(values) > 0 {
		values = append(values, strings.TrimSpace(current))
	}

	return values
}

// MapRowValuesToPrimitives converts string values to primitive types
func MapRowValuesToPrimitives(values []string) []toon.JsonPrimitive {
	result := make([]toon.JsonPrimitive, len(values))
	for i, v := range values {
		result[i] = ParsePrimitiveToken(v)
	}
	return result
}

// ParsePrimitiveToken parses a token into a primitive value
func ParsePrimitiveToken(token string) toon.JsonPrimitive {
	trimmed := strings.TrimSpace(token)

	// Empty token
	if trimmed == "" {
		return ""
	}

	// Quoted string (if starts with quote, it MUST be properly quoted)
	if len(trimmed) > 0 && trimmed[0] == byte(toon.DoubleQuote) {
		return ParseStringLiteral(trimmed)
	}

	// Boolean or null literals
	if shared.IsBooleanOrNullLiteral(trimmed) {
		if trimmed == toon.TrueLiteral {
			return true
		}
		if trimmed == toon.FalseLiteral {
			return false
		}
		if trimmed == toon.NullLiteral {
			return nil
		}
	}

	// Numeric literal
	if shared.IsNumericLiteral(trimmed) {
		parsedNumber, _ := strconv.ParseFloat(trimmed, 64)
		// Normalize negative zero to positive zero
		if parsedNumber == 0 && strings.HasPrefix(trimmed, "-") {
			return 0.0
		}
		return parsedNumber
	}

	// Unquoted string
	return trimmed
}

// ParseStringLiteral parses a string literal, handling quotes and escape sequences
func ParseStringLiteral(token string) string {
	trimmedToken := strings.TrimSpace(token)

	if len(trimmedToken) > 0 && trimmedToken[0] == byte(toon.DoubleQuote) {
		// Find the closing quote, accounting for escaped quotes
		closingQuoteIndex := shared.FindClosingQuote(trimmedToken, 0)

		if closingQuoteIndex == -1 {
			// No closing quote was found
			panic("unterminated string: missing closing quote")
		}

		if closingQuoteIndex != len([]rune(trimmedToken))-1 {
			panic("unexpected characters after closing quote")
		}

		content := string([]rune(trimmedToken)[1:closingQuoteIndex])
		unescaped, err := shared.UnescapeString(content)
		if err != nil {
			panic(err)
		}
		return unescaped
	}

	return trimmedToken
}

// KeyParseResult contains the result of parsing a key
type KeyParseResult struct {
	Key string
	End int
}

// ParseUnquotedKey parses an unquoted key from content starting at start
func ParseUnquotedKey(content string, start int) (*KeyParseResult, error) {
	runes := []rune(content)
	end := start
	for end < len(runes) && runes[end] != toon.Colon {
		end++
	}

	// Validate that a colon was found
	if end >= len(runes) || runes[end] != toon.Colon {
		return nil, fmt.Errorf("missing colon after key")
	}

	key := strings.TrimSpace(string(runes[start:end]))

	// Skip the colon
	end++

	return &KeyParseResult{Key: key, End: end}, nil
}

// ParseQuotedKey parses a quoted key from content starting at start
func ParseQuotedKey(content string, start int) (*KeyParseResult, error) {
	// Find the closing quote, accounting for escaped quotes
	closingQuoteIndex := shared.FindClosingQuote(content, start)

	if closingQuoteIndex == -1 {
		return nil, fmt.Errorf("unterminated quoted key")
	}

	// Extract and unescape the key content
	runes := []rune(content)
	keyContent := string(runes[start+1 : closingQuoteIndex])
	key, err := shared.UnescapeString(keyContent)
	if err != nil {
		return nil, err
	}
	end := closingQuoteIndex + 1

	// Validate and skip colon after quoted key
	if end >= len(runes) || runes[end] != toon.Colon {
		return nil, fmt.Errorf("missing colon after key")
	}
	end++

	return &KeyParseResult{Key: key, End: end}, nil
}

// ParseKeyToken parses a key (quoted or unquoted) from content starting at start
func ParseKeyToken(content string, start int) (*KeyParseResult, error) {
	runes := []rune(content)
	if start < len(runes) && runes[start] == toon.DoubleQuote {
		return ParseQuotedKey(content, start)
	}
	return ParseUnquotedKey(content, start)
}

// IsArrayHeaderAfterHyphen checks if content after hyphen is an array header
func IsArrayHeaderAfterHyphen(content string) bool {
	trimmed := strings.TrimSpace(content)
	if len(trimmed) == 0 || trimmed[0] != byte(toon.OpenBracket) {
		return false
	}
	return shared.FindUnquotedChar(trimmed, toon.Colon, 0) != -1
}

// IsObjectFirstFieldAfterHyphen checks if content after hyphen is an object field
func IsObjectFirstFieldAfterHyphen(content string) bool {
	return shared.FindUnquotedChar(content, toon.Colon, 0) != -1
}
