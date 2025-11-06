package decode

import (
	"fmt"
	"strings"

	"github.com/toon-format/toon-go"
	"github.com/toon-format/toon-go/shared"
)

// DecodeValueFromLines is the entry point for decoding TOON format.
// It determines whether the input is a root array, single primitive, or object.
func DecodeValueFromLines(cursor *LineCursor, options *toon.DecodeOptions) (toon.JsonValue, error) {
	first := cursor.Peek()
	if first == nil {
		return nil, fmt.Errorf("no content to decode")
	}

	// Check for root array
	if IsArrayHeaderAfterHyphen(first.Content) {
		headerInfo, inlineValues, err := ParseArrayHeaderLine(first.Content, toon.DefaultDelimiter)
		if err != nil {
			return nil, err
		}
		if headerInfo != nil {
			cursor.Advance() // Move past the header line
			return decodeArrayFromHeader(headerInfo, inlineValues, cursor, 0, options)
		}
	}

	// Check for single primitive value
	if cursor.Length() == 1 && !isKeyValueLine(first) {
		return ParsePrimitiveToken(first.Content)
	}

	// Default to object
	return decodeObject(cursor, 0, options)
}

// isKeyValueLine checks if a line contains a key-value pair (has a colon).
func isKeyValueLine(line *toon.ParsedLine) bool {
	content := line.Content

	// Look for unquoted colon or quoted key followed by colon
	if strings.HasPrefix(content, "\"") {
		// Quoted key - find the closing quote
		closingQuoteIndex := shared.FindClosingQuote(content, 0)
		if closingQuoteIndex == -1 {
			return false
		}
		// Check if colon exists after quoted key (may have array/brace syntax between)
		return strings.Contains(content[closingQuoteIndex+1:], string(toon.Colon))
	}

	// Unquoted key - look for first colon not inside quotes
	return strings.ContainsRune(content, toon.Colon)
}

// Placeholder functions - will be implemented in subsequent commits
func decodeObject(cursor *LineCursor, baseDepth int, options *toon.DecodeOptions) (toon.JsonObject, error) {
	return nil, fmt.Errorf("decodeObject not yet implemented")
}

func decodeArrayFromHeader(
	header *toon.ArrayHeaderInfo,
	inlineValues string,
	cursor *LineCursor,
	baseDepth int,
	options *toon.DecodeOptions,
) (toon.JsonArray, error) {
	return nil, fmt.Errorf("decodeArrayFromHeader not yet implemented")
}
