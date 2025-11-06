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

// decodeObject decodes a TOON object from the cursor at the specified depth.
func decodeObject(cursor *LineCursor, baseDepth int, options *toon.DecodeOptions) (toon.JsonObject, error) {
	obj := make(toon.JsonObject)

	// Detect the actual depth of the first field (may differ from baseDepth in nested structures)
	var computedDepth *int

	for !cursor.AtEnd() {
		line := cursor.Peek()
		if line == nil || line.Depth < baseDepth {
			break
		}

		if computedDepth == nil && line.Depth >= baseDepth {
			depth := line.Depth
			computedDepth = &depth
		}

		if line.Depth == *computedDepth {
			key, value, err := decodeKeyValuePair(line, cursor, *computedDepth, options)
			if err != nil {
				return nil, err
			}
			obj[key] = value
		} else {
			// Different depth (shallower or deeper) - stop object parsing
			break
		}
	}

	return obj, nil
}

// decodeKeyValue parses a key-value line and returns the key, value, and the depth for following fields.
func decodeKeyValue(
	content string,
	cursor *LineCursor,
	baseDepth int,
	options *toon.DecodeOptions,
) (key string, value toon.JsonValue, followDepth int, err error) {
	// Check for array header first (before parsing key)
	arrayHeader, inlineValues, parseErr := ParseArrayHeaderLine(content, toon.DefaultDelimiter)
	if parseErr == nil && arrayHeader != nil && arrayHeader.Key != "" {
		val, arrErr := decodeArrayFromHeader(arrayHeader, inlineValues, cursor, baseDepth, options)
		if arrErr != nil {
			return "", nil, 0, arrErr
		}
		// After an array, subsequent fields are at baseDepth + 1 (where array content is)
		return arrayHeader.Key, val, baseDepth + 1, nil
	}

	// Regular key-value pair
	parsedKey, end, keyErr := ParseKeyToken(content, 0)
	if keyErr != nil {
		return "", nil, 0, keyErr
	}

	runes := []rune(content)
	rest := strings.TrimSpace(string(runes[end:]))

	// No value after colon - expect nested object or empty
	if rest == "" {
		nextLine := cursor.Peek()
		if nextLine != nil && nextLine.Depth > baseDepth {
			nested, nestedErr := decodeObject(cursor, baseDepth+1, options)
			if nestedErr != nil {
				return "", nil, 0, nestedErr
			}
			return parsedKey, nested, baseDepth + 1, nil
		}
		// Empty object
		return parsedKey, make(toon.JsonObject), baseDepth + 1, nil
	}

	// Inline primitive value
	primValue, primErr := ParsePrimitiveToken(rest)
	if primErr != nil {
		return "", nil, 0, primErr
	}
	return parsedKey, primValue, baseDepth + 1, nil
}

// decodeKeyValuePair decodes a single key-value pair from a line.
func decodeKeyValuePair(
	line *toon.ParsedLine,
	cursor *LineCursor,
	baseDepth int,
	options *toon.DecodeOptions,
) (key string, value toon.JsonValue, err error) {
	cursor.Advance()
	key, value, _, err = decodeKeyValue(line.Content, cursor, baseDepth, options)
	return key, value, err
}

// decodeArrayFromHeader dispatches to the appropriate array decoder based on header info.
func decodeArrayFromHeader(
	header *toon.ArrayHeaderInfo,
	inlineValues string,
	cursor *LineCursor,
	baseDepth int,
	options *toon.DecodeOptions,
) (toon.JsonArray, error) {
	// Inline primitive array
	if inlineValues != "" {
		return decodeInlinePrimitiveArray(header, inlineValues, options)
	}

	// Tabular array
	if header.Fields != nil && len(header.Fields) > 0 {
		return decodeTabularArray(header, cursor, baseDepth, options)
	}

	// List array
	return decodeListArray(header, cursor, baseDepth, options)
}

// decodeInlinePrimitiveArray decodes an inline array like [3]: a,b,c
func decodeInlinePrimitiveArray(
	header *toon.ArrayHeaderInfo,
	inlineValues string,
	options *toon.DecodeOptions,
) (toon.JsonArray, error) {
	if strings.TrimSpace(inlineValues) == "" {
		if err := AssertExpectedCount(0, header.Length, "inline array items", options); err != nil {
			return nil, err
		}
		return toon.JsonArray{}, nil
	}

	values, err := ParseDelimitedValues(inlineValues, header.Delimiter)
	if err != nil {
		return nil, err
	}

	primitives, err := MapRowValuesToPrimitives(values)
	if err != nil {
		return nil, err
	}

	if err := AssertExpectedCount(len(primitives), header.Length, "inline array items", options); err != nil {
		return nil, err
	}

	// Convert []JsonPrimitive to JsonArray
	result := make(toon.JsonArray, len(primitives))
	for i, p := range primitives {
		result[i] = p
	}

	return result, nil
}

// decodeTabularArray decodes a tabular array with fields like items[2]{id,name,price}:
func decodeTabularArray(
	header *toon.ArrayHeaderInfo,
	cursor *LineCursor,
	baseDepth int,
	options *toon.DecodeOptions,
) (toon.JsonArray, error) {
	objects := make([]toon.JsonObject, 0, header.Length)
	rowDepth := baseDepth + 1

	// Track line range for blank line validation
	var startLine, endLine int

	for !cursor.AtEnd() && len(objects) < header.Length {
		line := cursor.Peek()
		if line == nil || line.Depth < rowDepth {
			break
		}

		if line.Depth == rowDepth {
			// Track first and last row line numbers
			if len(objects) == 0 {
				startLine = line.LineNumber
			}
			endLine = line.LineNumber

			cursor.Advance()
			values, err := ParseDelimitedValues(line.Content, header.Delimiter)
			if err != nil {
				return nil, err
			}

			if err := AssertExpectedCount(len(values), len(header.Fields), "tabular row values", options); err != nil {
				return nil, err
			}

			primitives, err := MapRowValuesToPrimitives(values)
			if err != nil {
				return nil, err
			}

			obj := make(toon.JsonObject)
			for i := 0; i < len(header.Fields); i++ {
				obj[header.Fields[i]] = primitives[i]
			}

			objects = append(objects, obj)
		} else {
			break
		}
	}

	if err := AssertExpectedCount(len(objects), header.Length, "tabular rows", options); err != nil {
		return nil, err
	}

	// In strict mode, check for blank lines inside the array
	if options.Strict && len(objects) > 0 {
		if err := ValidateNoBlankLinesInRange(
			startLine,
			endLine,
			cursor.GetBlankLines(),
			options.Strict,
			"tabular array",
		); err != nil {
			return nil, err
		}
	}

	// In strict mode, check for extra rows
	if options.Strict {
		if err := ValidateNoExtraTabularRows(cursor, rowDepth, header); err != nil {
			return nil, err
		}
	}

	// Convert []JsonObject to JsonArray
	result := make(toon.JsonArray, len(objects))
	for i, obj := range objects {
		result[i] = obj
	}

	return result, nil
}

func decodeListArray(
	header *toon.ArrayHeaderInfo,
	cursor *LineCursor,
	baseDepth int,
	options *toon.DecodeOptions,
) (toon.JsonArray, error) {
	return nil, fmt.Errorf("decodeListArray not yet implemented")
}
