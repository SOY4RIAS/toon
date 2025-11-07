package encode

import (
	"fmt"
)

const listItemMarker = '-'

// ResolvedEncodeOptions contains resolved encoding options
type ResolvedEncodeOptions struct {
	Indent       int
	Delimiter    rune
	LengthMarker rune
}

// EncodeValue encodes a normalized JSON value to TOON format
func EncodeValue(value interface{}, options ResolvedEncodeOptions) string {
	if IsJsonPrimitive(value) {
		return EncodePrimitive(value, options.Delimiter)
	}

	writer := NewLineWriter(options.Indent)

	if IsJsonArray(value) {
		EncodeArray("", value.([]interface{}), writer, 0, options)
	} else if IsJsonObject(value) {
		EncodeObject(value.(map[string]interface{}), writer, 0, options)
	}

	return writer.String()
}

// EncodeObject encodes a JSON object
func EncodeObject(obj map[string]interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	for key, value := range obj {
		EncodeKeyValuePair(key, value, writer, depth, options)
	}
}

// EncodeKeyValuePair encodes a key-value pair
func EncodeKeyValuePair(key string, value interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	encodedKey := EncodeKey(key)

	if IsJsonPrimitive(value) {
		line := fmt.Sprintf("%s: %s", encodedKey, EncodePrimitive(value, options.Delimiter))
		writer.Push(depth, line)
	} else if IsJsonArray(value) {
		EncodeArray(key, value.([]interface{}), writer, depth, options)
	} else if IsJsonObject(value) {
		obj := value.(map[string]interface{})
		if len(obj) == 0 {
			// Empty object
			writer.Push(depth, encodedKey+":")
		} else {
			writer.Push(depth, encodedKey+":")
			EncodeObject(obj, writer, depth+1, options)
		}
	}
}

// EncodeArray encodes a JSON array
func EncodeArray(key string, arr []interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	if len(arr) == 0 {
		header := FormatHeader(0, FormatHeaderOptions{
			Key:          key,
			Delimiter:    options.Delimiter,
			LengthMarker: options.LengthMarker,
		})
		writer.Push(depth, header)
		return
	}

	// Primitive array - inline format
	if IsArrayOfPrimitives(arr) {
		formatted := EncodeInlineArrayLine(arr, options.Delimiter, key, options.LengthMarker)
		writer.Push(depth, formatted)
		return
	}

	// Array of arrays (all primitives) - list format
	if IsArrayOfArrays(arr) {
		allPrimitiveArrays := true
		for _, item := range arr {
			if !IsArrayOfPrimitives(item) {
				allPrimitiveArrays = false
				break
			}
		}
		if allPrimitiveArrays {
			EncodeArrayOfArraysAsListItems(key, arr, writer, depth, options)
			return
		}
	}

	// Array of objects - try tabular format
	if IsArrayOfObjects(arr) {
		header := ExtractTabularHeader(arr)
		if header != nil {
			EncodeArrayOfObjectsAsTabular(key, arr, header, writer, depth, options)
		} else {
			EncodeMixedArrayAsListItems(key, arr, writer, depth, options)
		}
		return
	}

	// Mixed array - list format
	EncodeMixedArrayAsListItems(key, arr, writer, depth, options)
}

// EncodeInlineArrayLine encodes a primitive array as an inline format
// Example: tags[3]: a,b,c
func EncodeInlineArrayLine(values []interface{}, delimiter rune, key string, lengthMarker rune) string {
	header := FormatHeader(len(values), FormatHeaderOptions{
		Key:          key,
		Delimiter:    delimiter,
		LengthMarker: lengthMarker,
	})

	if len(values) == 0 {
		return header
	}

	joinedValue := EncodeAndJoinPrimitives(values, delimiter)
	return header + " " + joinedValue
}

// EncodeArrayOfArraysAsListItems encodes an array of primitive arrays as list items
func EncodeArrayOfArraysAsListItems(key string, arrays []interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	header := FormatHeader(len(arrays), FormatHeaderOptions{
		Key:          key,
		Delimiter:    options.Delimiter,
		LengthMarker: options.LengthMarker,
	})
	writer.Push(depth, header)

	for _, item := range arrays {
		if arr, ok := item.([]interface{}); ok && IsArrayOfPrimitives(arr) {
			inline := EncodeInlineArrayLine(arr, options.Delimiter, "", options.LengthMarker)
			writer.PushListItem(depth+1, inline)
		}
	}
}

// ExtractTabularHeader checks if an array of objects can use tabular format
// Returns the field names if all objects have the same primitive fields, nil otherwise
func ExtractTabularHeader(arr []interface{}) []string {
	if len(arr) == 0 {
		return nil
	}

	firstObj, ok := arr[0].(map[string]interface{})
	if !ok || len(firstObj) == 0 {
		return nil
	}

	// Extract keys from first object
	var header []string
	for key := range firstObj {
		header = append(header, key)
	}

	// Check if this is a valid tabular array
	if IsTabularArray(arr, header) {
		return header
	}

	return nil
}

// IsTabularArray checks if all objects in the array have the same primitive fields
func IsTabularArray(arr []interface{}, header []string) bool {
	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return false
		}

		// Check that object has exactly the same keys
		if len(obj) != len(header) {
			return false
		}

		// Check that all header keys exist and all values are primitives
		for _, key := range header {
			value, exists := obj[key]
			if !exists {
				return false
			}
			if !IsJsonPrimitive(value) {
				return false
			}
		}
	}

	return true
}

// EncodeArrayOfObjectsAsTabular encodes an array of objects in tabular format
func EncodeArrayOfObjectsAsTabular(key string, rows []interface{}, header []string, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	formattedHeader := FormatHeader(len(rows), FormatHeaderOptions{
		Key:          key,
		Fields:       header,
		Delimiter:    options.Delimiter,
		LengthMarker: options.LengthMarker,
	})
	writer.Push(depth, formattedHeader)

	WriteTabularRows(rows, header, writer, depth+1, options)
}

// WriteTabularRows writes the data rows for a tabular array
func WriteTabularRows(rows []interface{}, header []string, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	for _, row := range rows {
		obj := row.(map[string]interface{})
		values := make([]interface{}, len(header))
		for i, key := range header {
			values[i] = obj[key]
		}
		joinedValue := EncodeAndJoinPrimitives(values, options.Delimiter)
		writer.Push(depth, joinedValue)
	}
}

// EncodeMixedArrayAsListItems encodes a mixed array as list items
func EncodeMixedArrayAsListItems(key string, items []interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	header := FormatHeader(len(items), FormatHeaderOptions{
		Key:          key,
		Delimiter:    options.Delimiter,
		LengthMarker: options.LengthMarker,
	})
	writer.Push(depth, header)

	for _, item := range items {
		EncodeListItemValue(item, writer, depth+1, options)
	}
}

// EncodeListItemValue encodes a single list item value
func EncodeListItemValue(value interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	if IsJsonPrimitive(value) {
		writer.PushListItem(depth, EncodePrimitive(value, options.Delimiter))
	} else if arr, ok := value.([]interface{}); ok && IsArrayOfPrimitives(arr) {
		inline := EncodeInlineArrayLine(arr, options.Delimiter, "", options.LengthMarker)
		writer.PushListItem(depth, inline)
	} else if IsJsonObject(value) {
		EncodeObjectAsListItem(value.(map[string]interface{}), writer, depth, options)
	}
}

// EncodeObjectAsListItem encodes an object as a list item
func EncodeObjectAsListItem(obj map[string]interface{}, writer *LineWriter, depth int, options ResolvedEncodeOptions) {
	if len(obj) == 0 {
		writer.Push(depth, string(listItemMarker))
		return
	}

	// Get keys (order may vary, but we need to process them consistently)
	var keys []string
	for key := range obj {
		keys = append(keys, key)
	}

	// First key-value on the same line as "- "
	firstKey := keys[0]
	encodedKey := EncodeKey(firstKey)
	firstValue := obj[firstKey]

	if IsJsonPrimitive(firstValue) {
		line := fmt.Sprintf("%s: %s", encodedKey, EncodePrimitive(firstValue, options.Delimiter))
		writer.PushListItem(depth, line)
	} else if arr, ok := firstValue.([]interface{}); ok {
		if IsArrayOfPrimitives(arr) {
			// Inline format for primitive arrays
			formatted := EncodeInlineArrayLine(arr, options.Delimiter, firstKey, options.LengthMarker)
			writer.PushListItem(depth, formatted)
		} else if IsArrayOfObjects(arr) {
			// Check if array of objects can use tabular format
			header := ExtractTabularHeader(arr)
			if header != nil {
				// Tabular format for uniform arrays of objects
				formattedHeader := FormatHeader(len(arr), FormatHeaderOptions{
					Key:          firstKey,
					Fields:       header,
					Delimiter:    options.Delimiter,
					LengthMarker: options.LengthMarker,
				})
				writer.PushListItem(depth, formattedHeader)
				WriteTabularRows(arr, header, writer, depth+1, options)
			} else {
				// Fall back to list format for non-uniform arrays of objects
				header := fmt.Sprintf("%s[%d]:", encodedKey, len(arr))
				writer.PushListItem(depth, header)
				for _, item := range arr {
					if IsJsonObject(item) {
						EncodeObjectAsListItem(item.(map[string]interface{}), writer, depth+1, options)
					}
				}
			}
		} else {
			// Complex arrays on separate lines (array of arrays, etc.)
			header := fmt.Sprintf("%s[%d]:", encodedKey, len(arr))
			writer.PushListItem(depth, header)

			// Encode array contents at depth + 1
			for _, item := range arr {
				EncodeListItemValue(item, writer, depth+1, options)
			}
		}
	} else if nestedObj, ok := firstValue.(map[string]interface{}); ok {
		if len(nestedObj) == 0 {
			writer.PushListItem(depth, encodedKey+":")
		} else {
			writer.PushListItem(depth, encodedKey+":")
			EncodeObject(nestedObj, writer, depth+2, options)
		}
	}

	// Remaining keys on indented lines
	for i := 1; i < len(keys); i++ {
		key := keys[i]
		EncodeKeyValuePair(key, obj[key], writer, depth+1, options)
	}
}
