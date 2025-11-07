package toon

import (
	"fmt"
	"strings"

	"github.com/soy4rias/toongo/shared"
)

// DecodeValueFromLines decodes a value from parsed lines
func DecodeValueFromLines(cursor *LineCursor, options ResolvedDecodeOptions) (JsonValue, error) {
	first := cursor.Peek()
	if first == nil {
		return nil, fmt.Errorf("no content to decode")
	}

	// Check for root array
	if IsArrayHeaderAfterHyphen(first.Content) {
		header, inlineValues, err := ParseArrayHeaderLine(first.Content, DefaultDelimiter)
		if err == nil && header != nil {
			cursor.Advance() // Move past the header line
			return decodeArrayFromHeader(header, inlineValues, cursor, 0, options)
		}
	}

	// Check for single primitive value
	if cursor.Length() == 1 && !isKeyValueLine(first) {
		return parsePrimitiveToken(strings.TrimSpace(first.Content)), nil
	}

	// Default to object
	return decodeObject(cursor, 0, options)
}

func isKeyValueLine(line *ParsedLine) bool {
	content := line.Content
	// Look for unquoted colon or quoted key followed by colon
	if strings.HasPrefix(content, "\"") {
		// Quoted key - find the closing quote
		closingQuoteIndex := shared.FindClosingQuote(content, 0)
		if closingQuoteIndex == -1 {
			return false
		}
		// Check if colon exists after quoted key
		return strings.Contains(content[closingQuoteIndex+1:], ":")
	}
	// Unquoted key - look for first colon
	return strings.Contains(content, ":")
}

// decodeObject decodes an object from lines
func decodeObject(cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (JsonObject, error) {
	obj := make(JsonObject)

	// Detect the actual depth of the first field
	computedDepth := -1

	for !cursor.AtEnd() {
		line := cursor.Peek()
		if line == nil || line.Depth < baseDepth {
			break
		}

		if computedDepth == -1 && line.Depth >= baseDepth {
			computedDepth = line.Depth
		}

		if line.Depth == computedDepth {
			key, value, err := decodeKeyValuePair(line, cursor, computedDepth, options)
			if err != nil {
				return nil, err
			}
			obj[key] = value
		} else {
			// Different depth - stop object parsing
			break
		}
	}

	return obj, nil
}

func decodeKeyValue(content string, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (string, JsonValue, int, error) {
	// Check for array header first
	arrayHeader, inlineValues, err := ParseArrayHeaderLine(content, DefaultDelimiter)
	if err == nil && arrayHeader != nil && arrayHeader.Key != "" {
		value, err := decodeArrayFromHeader(arrayHeader, inlineValues, cursor, baseDepth, options)
		if err != nil {
			return "", nil, 0, err
		}
		return arrayHeader.Key, value, baseDepth + 1, nil
	}

	// Regular key-value pair
	key, end, err := parseKeyToken(content, 0)
	if err != nil {
		return "", nil, 0, err
	}
	rest := strings.TrimSpace(content[end:])

	// No value after colon - expect nested object or empty
	if rest == "" {
		nextLine := cursor.Peek()
		if nextLine != nil && nextLine.Depth > baseDepth {
			nested, err := decodeObject(cursor, baseDepth+1, options)
			if err != nil {
				return "", nil, 0, err
			}
			return key, nested, baseDepth + 1, nil
		}
		// Empty object
		return key, make(JsonObject), baseDepth + 1, nil
	}

	// Inline primitive value
	value := parsePrimitiveToken(rest)
	return key, value, baseDepth + 1, nil
}

func decodeKeyValuePair(line *ParsedLine, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (string, JsonValue, error) {
	cursor.Advance()
	key, value, _, err := decodeKeyValue(line.Content, cursor, baseDepth, options)
	return key, value, err
}

// decodeArrayFromHeader decodes an array based on header info
func decodeArrayFromHeader(header *ArrayHeaderInfo, inlineValues string, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (JsonValue, error) {
	// Inline primitive array
	if inlineValues != "" {
		return decodeInlinePrimitiveArray(header, inlineValues, options)
	}

	// Tabular array
	if len(header.Fields) > 0 {
		return decodeTabularArray(header, cursor, baseDepth, options)
	}

	// List array
	return decodeListArray(header, cursor, baseDepth, options)
}

func decodeInlinePrimitiveArray(header *ArrayHeaderInfo, inlineValues string, options ResolvedDecodeOptions) (JsonArray, error) {
	if strings.TrimSpace(inlineValues) == "" {
		if err := AssertExpectedCount(0, header.Length, "inline array items", options); err != nil {
			return nil, err
		}
		return JsonArray{}, nil
	}

	values := parseDelimitedValues(inlineValues, header.Delimiter)
	primitives := mapRowValuesToPrimitives(values)

	if err := AssertExpectedCount(len(primitives), header.Length, "inline array items", options); err != nil {
		return nil, err
	}

	result := make(JsonArray, len(primitives))
	for i, p := range primitives {
		result[i] = p
	}
	return result, nil
}

func decodeListArray(header *ArrayHeaderInfo, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (JsonArray, error) {
	items := JsonArray{}
	itemDepth := baseDepth + 1

	startLine, endLine := -1, -1

	for !cursor.AtEnd() && len(items) < header.Length {
		line := cursor.Peek()
		if line == nil || line.Depth < itemDepth {
			break
		}

		// Check for list item
		isListItem := strings.HasPrefix(line.Content, ListItemPrefix) || line.Content == "-"

		if line.Depth == itemDepth && isListItem {
			// Track first and last item line numbers
			if startLine == -1 {
				startLine = line.LineNumber
			}
			endLine = line.LineNumber

			item, err := decodeListItem(cursor, itemDepth, options)
			if err != nil {
				return nil, err
			}
			items = append(items, item)

			// Update endLine to current cursor position
			currentLine := cursor.Current()
			if currentLine != nil {
				endLine = currentLine.LineNumber
			}
		} else {
			break
		}
	}

	if err := AssertExpectedCount(len(items), header.Length, "list array items", options); err != nil {
		return nil, err
	}

	// Validate blank lines in strict mode
	if options.Strict && startLine != -1 && endLine != -1 {
		if err := ValidateNoBlankLinesInRange(startLine, endLine, cursor.GetBlankLines(), options.Strict, "list array"); err != nil {
			return nil, err
		}
	}

	// Check for extra items in strict mode
	if options.Strict {
		if err := ValidateNoExtraListItems(cursor, itemDepth, header.Length); err != nil {
			return nil, err
		}
	}

	return items, nil
}

func decodeTabularArray(header *ArrayHeaderInfo, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (JsonArray, error) {
	objects := []JsonObject{}
	rowDepth := baseDepth + 1

	startLine, endLine := -1, -1

	for !cursor.AtEnd() && len(objects) < header.Length {
		line := cursor.Peek()
		if line == nil || line.Depth < rowDepth {
			break
		}

		if line.Depth == rowDepth {
			// Track first and last row line numbers
			if startLine == -1 {
				startLine = line.LineNumber
			}
			endLine = line.LineNumber

			cursor.Advance()
			values := parseDelimitedValues(line.Content, header.Delimiter)
			if err := AssertExpectedCount(len(values), len(header.Fields), "tabular row values", options); err != nil {
				return nil, err
			}

			primitives := mapRowValuesToPrimitives(values)
			obj := make(JsonObject)

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

	// Validate blank lines in strict mode
	if options.Strict && startLine != -1 && endLine != -1 {
		if err := ValidateNoBlankLinesInRange(startLine, endLine, cursor.GetBlankLines(), options.Strict, "tabular array"); err != nil {
			return nil, err
		}
	}

	// Check for extra rows in strict mode
	if options.Strict {
		if err := ValidateNoExtraTabularRows(cursor, rowDepth, header); err != nil {
			return nil, err
		}
	}

	result := make(JsonArray, len(objects))
	for i, obj := range objects {
		result[i] = obj
	}
	return result, nil
}

// decodeListItem decodes a single list item
func decodeListItem(cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (JsonValue, error) {
	line := cursor.Next()
	if line == nil {
		return nil, fmt.Errorf("expected list item")
	}

	// Check for list item
	var afterHyphen string

	// Empty list item should be an empty object
	if line.Content == "-" {
		return make(JsonObject), nil
	} else if strings.HasPrefix(line.Content, ListItemPrefix) {
		afterHyphen = line.Content[len(ListItemPrefix):]
	} else {
		return nil, fmt.Errorf("expected list item to start with \"%s\"", ListItemPrefix)
	}

	// Empty content after list item should also be an empty object
	if strings.TrimSpace(afterHyphen) == "" {
		return make(JsonObject), nil
	}

	// Check for array header after hyphen
	if IsArrayHeaderAfterHyphen(afterHyphen) {
		arrayHeader, inlineValues, err := ParseArrayHeaderLine(afterHyphen, DefaultDelimiter)
		if err == nil && arrayHeader != nil {
			return decodeArrayFromHeader(arrayHeader, inlineValues, cursor, baseDepth, options)
		}
	}

	// Check for object first field after hyphen
	if IsObjectFirstFieldAfterHyphen(afterHyphen) {
		return decodeObjectFromListItem(line, afterHyphen, cursor, baseDepth, options)
	}

	// Primitive value
	return parsePrimitiveToken(afterHyphen), nil
}

func decodeObjectFromListItem(firstLine *ParsedLine, afterHyphen string, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (JsonObject, error) {
	key, value, followDepth, err := decodeKeyValue(afterHyphen, cursor, baseDepth, options)
	if err != nil {
		return nil, err
	}

	obj := make(JsonObject)
	obj[key] = value

	// Read subsequent fields
	for !cursor.AtEnd() {
		line := cursor.Peek()
		if line == nil || line.Depth < followDepth {
			break
		}

		if line.Depth == followDepth && !strings.HasPrefix(line.Content, ListItemPrefix) {
			k, v, err := decodeKeyValuePair(line, cursor, followDepth, options)
			if err != nil {
				return nil, err
			}
			obj[k] = v
		} else {
			break
		}
	}

	return obj, nil
}
