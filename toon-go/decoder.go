package toon

import (
	"fmt"
	"strings"
)

// DecodeValueFromLines decodes a value from the parsed lines
func DecodeValueFromLines(cursor *LineCursor, options ResolvedDecodeOptions) (interface{}, error) {
	first := cursor.Peek()
	if first == nil {
		return nil, fmt.Errorf("no content to decode")
	}

	// Check for root array
	if IsArrayHeaderAfterHyphen(first.Content) {
		headerInfo := ParseArrayHeaderLine(first.Content, DefaultDelimiter)
		if headerInfo != nil {
			cursor.Advance() // Move past the header line
			inlineVals := ""
			if headerInfo.InlineValues != "" {
				inlineVals = headerInfo.InlineValues
			}
			return DecodeArrayFromHeader(headerInfo.Header, inlineVals, cursor, 0, options)
		}
	}

	// Check for single primitive value
	if cursor.Length() == 1 && !isKeyValueLine(first) {
		return ParsePrimitiveToken(strings.TrimSpace(first.Content))
	}

	// Default to object
	return DecodeObject(cursor, 0, options)
}

// isKeyValueLine checks if a line contains a key-value pair
func isKeyValueLine(line *ParsedLine) bool {
	content := line.Content
	// Look for unquoted colon or quoted key followed by colon
	if strings.HasPrefix(content, string(DoubleQuote)) {
		// Quoted key - find the closing quote
		closingQuoteIndex := FindClosingQuote(content, 0)
		if closingQuoteIndex == -1 {
			return false
		}
		// Check if colon exists after quoted key (may have array/brace syntax between)
		return strings.Contains(content[closingQuoteIndex+1:], string(Colon))
	}
	// Unquoted key - look for first colon not inside quotes
	return strings.ContainsRune(content, Colon)
}

// DecodeObject decodes an object from the cursor
func DecodeObject(cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (map[string]interface{}, error) {
	obj := make(map[string]interface{})

	// Detect the actual depth of the first field (may differ from baseDepth in nested structures)
	var computedDepth *int

	for !cursor.AtEnd() {
		line := cursor.Peek()
		if line == nil || line.Depth < baseDepth {
			break
		}

		if computedDepth == nil && line.Depth >= baseDepth {
			d := line.Depth
			computedDepth = &d
		}

		if line.Depth == *computedDepth {
			key, value, err := DecodeKeyValuePair(line, cursor, *computedDepth, options)
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

// KeyValueResult contains the result of decoding a key-value pair
type KeyValueResult struct {
	Key         string
	Value       interface{}
	FollowDepth int
}

// DecodeKeyValue decodes a key-value pair from a line
func DecodeKeyValue(content string, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (*KeyValueResult, error) {
	// Check for array header first (before parsing key)
	arrayHeader := ParseArrayHeaderLine(content, DefaultDelimiter)
	if arrayHeader != nil && arrayHeader.Header.Key != nil {
		inlineVals := ""
		if arrayHeader.InlineValues != "" {
			inlineVals = arrayHeader.InlineValues
		}
		value, err := DecodeArrayFromHeader(arrayHeader.Header, inlineVals, cursor, baseDepth, options)
		if err != nil {
			return nil, err
		}
		// After an array, subsequent fields are at baseDepth + 1 (where array content is)
		return &KeyValueResult{
			Key:         *arrayHeader.Header.Key,
			Value:       value,
			FollowDepth: baseDepth + 1,
		}, nil
	}

	// Regular key-value pair
	keyResult, err := ParseKeyToken(content, 0)
	if err != nil {
		return nil, err
	}

	rest := strings.TrimSpace(content[keyResult.End:])

	// No value after colon - expect nested object or empty
	if rest == "" {
		nextLine := cursor.Peek()
		if nextLine != nil && nextLine.Depth > baseDepth {
			nested, err := DecodeObject(cursor, baseDepth+1, options)
			if err != nil {
				return nil, err
			}
			return &KeyValueResult{
				Key:         keyResult.Key,
				Value:       nested,
				FollowDepth: baseDepth + 1,
			}, nil
		}
		// Empty object
		return &KeyValueResult{
			Key:         keyResult.Key,
			Value:       make(map[string]interface{}),
			FollowDepth: baseDepth + 1,
		}, nil
	}

	// Inline primitive value
	value, err := ParsePrimitiveToken(rest)
	if err != nil {
		return nil, err
	}
	return &KeyValueResult{
		Key:         keyResult.Key,
		Value:       value,
		FollowDepth: baseDepth + 1,
	}, nil
}

// DecodeKeyValuePair decodes a key-value pair from a parsed line
func DecodeKeyValuePair(line *ParsedLine, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (string, interface{}, error) {
	cursor.Advance()
	result, err := DecodeKeyValue(line.Content, cursor, baseDepth, options)
	if err != nil {
		return "", nil, err
	}
	return result.Key, result.Value, nil
}

// DecodeArrayFromHeader decodes an array from a header
func DecodeArrayFromHeader(header ArrayHeaderInfo, inlineValues string, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (interface{}, error) {
	// Inline primitive array
	if inlineValues != "" {
		return DecodeInlinePrimitiveArray(header, inlineValues, options)
	}

	// Tabular array
	if header.Fields != nil && len(header.Fields) > 0 {
		return DecodeTabularArray(header, cursor, baseDepth, options)
	}

	// List array
	return DecodeListArray(header, cursor, baseDepth, options)
}

// DecodeInlinePrimitiveArray decodes an inline primitive array
func DecodeInlinePrimitiveArray(header ArrayHeaderInfo, inlineValues string, options ResolvedDecodeOptions) ([]interface{}, error) {
	trimmed := strings.TrimSpace(inlineValues)
	if trimmed == "" {
		if err := AssertExpectedCount(0, header.Length, "inline array items", options); err != nil {
			return nil, err
		}
		return []interface{}{}, nil
	}

	values := ParseDelimitedValues(inlineValues, header.Delimiter)
	primitives, err := MapRowValuesToPrimitives(values)
	if err != nil {
		return nil, err
	}

	if err := AssertExpectedCount(len(primitives), header.Length, "inline array items", options); err != nil {
		return nil, err
	}

	return primitives, nil
}

// DecodeListArray decodes a list-style array
func DecodeListArray(header ArrayHeaderInfo, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) ([]interface{}, error) {
	items := []interface{}{}
	itemDepth := baseDepth + 1

	// Track line range for blank line validation
	var startLine, endLine int

	for !cursor.AtEnd() && len(items) < header.Length {
		line := cursor.Peek()
		if line == nil || line.Depth < itemDepth {
			break
		}

		// Check for list item (with or without space after hyphen)
		isListItem := strings.HasPrefix(line.Content, ListItemPrefix) || line.Content == "-"

		if line.Depth == itemDepth && isListItem {
			// Track first and last item line numbers
			if len(items) == 0 {
				startLine = line.LineNumber
			}
			endLine = line.LineNumber

			item, err := DecodeListItem(cursor, itemDepth, options)
			if err != nil {
				return nil, err
			}
			items = append(items, item)

			// Update endLine to the current cursor position (after item was decoded)
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

	// In strict mode, check for blank lines inside the array
	if options.Strict && len(items) > 0 {
		if err := ValidateNoBlankLinesInRange(
			startLine,
			endLine,
			cursor.GetBlankLines(),
			options.Strict,
			"list array",
		); err != nil {
			return nil, err
		}
	}

	// In strict mode, check for extra items
	if options.Strict {
		if err := ValidateNoExtraListItems(cursor, itemDepth, header.Length); err != nil {
			return nil, err
		}
	}

	return items, nil
}

// DecodeTabularArray decodes a tabular-style array
func DecodeTabularArray(header ArrayHeaderInfo, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) ([]map[string]interface{}, error) {
	objects := []map[string]interface{}{}
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
			values := ParseDelimitedValues(line.Content, header.Delimiter)
			if err := AssertExpectedCount(len(values), len(header.Fields), "tabular row values", options); err != nil {
				return nil, err
			}

			primitives, err := MapRowValuesToPrimitives(values)
			if err != nil {
				return nil, err
			}

			obj := make(map[string]interface{})
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

	return objects, nil
}

// DecodeListItem decodes a single list item
func DecodeListItem(cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (interface{}, error) {
	line := cursor.Next()
	if line == nil {
		return nil, fmt.Errorf("expected list item")
	}

	// Empty list item should be an empty object
	if line.Content == "-" {
		return make(map[string]interface{}), nil
	}

	var afterHyphen string
	if strings.HasPrefix(line.Content, ListItemPrefix) {
		afterHyphen = line.Content[len(ListItemPrefix):]
	} else {
		return nil, fmt.Errorf("expected list item to start with \"%s\"", ListItemPrefix)
	}

	// Empty content after list item should also be an empty object
	if strings.TrimSpace(afterHyphen) == "" {
		return make(map[string]interface{}), nil
	}

	// Check for array header after hyphen
	if IsArrayHeaderAfterHyphen(afterHyphen) {
		arrayHeader := ParseArrayHeaderLine(afterHyphen, DefaultDelimiter)
		if arrayHeader != nil {
			inlineVals := ""
			if arrayHeader.InlineValues != "" {
				inlineVals = arrayHeader.InlineValues
			}
			return DecodeArrayFromHeader(arrayHeader.Header, inlineVals, cursor, baseDepth, options)
		}
	}

	// Check for object first field after hyphen
	if IsObjectFirstFieldAfterHyphen(afterHyphen) {
		return DecodeObjectFromListItem(line, cursor, baseDepth, options)
	}

	// Primitive value
	return ParsePrimitiveToken(afterHyphen)
}

// DecodeObjectFromListItem decodes an object that starts on a list item line
func DecodeObjectFromListItem(firstLine *ParsedLine, cursor *LineCursor, baseDepth int, options ResolvedDecodeOptions) (map[string]interface{}, error) {
	afterHyphen := firstLine.Content[len(ListItemPrefix):]
	result, err := DecodeKeyValue(afterHyphen, cursor, baseDepth, options)
	if err != nil {
		return nil, err
	}

	obj := make(map[string]interface{})
	obj[result.Key] = result.Value

	// Read subsequent fields
	for !cursor.AtEnd() {
		line := cursor.Peek()
		if line == nil || line.Depth < result.FollowDepth {
			break
		}

		if line.Depth == result.FollowDepth && !strings.HasPrefix(line.Content, ListItemPrefix) {
			key, value, err := DecodeKeyValuePair(line, cursor, result.FollowDepth, options)
			if err != nil {
				return nil, err
			}
			obj[key] = value
		} else {
			break
		}
	}

	return obj, nil
}
