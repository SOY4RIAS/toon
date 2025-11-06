package encode

import (
	"reflect"
)

const listItemMarker = '-'

// EncodeValue encodes a normalized JsonValue to TOON format.
func EncodeValue(value interface{}, indent int, delimiter rune, lengthMarker bool) string {
	if IsJsonPrimitive(value) {
		return EncodePrimitive(value, delimiter)
	}

	writer := NewLineWriter(indent)

	if IsJsonArray(value) {
		arr := toSlice(value)
		encodeArray("", arr, writer, 0, delimiter, lengthMarker)
	} else if IsJsonObject(value) {
		obj := toMap(value)
		encodeObject(obj, writer, 0, delimiter, lengthMarker)
	}

	return writer.String()
}

// encodeObject encodes a JSON object.
func encodeObject(value map[string]interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	for key, val := range value {
		encodeKeyValuePair(key, val, writer, depth, delimiter, lengthMarker)
	}
}

// encodeKeyValuePair encodes a single key-value pair.
func encodeKeyValuePair(key string, value interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	encodedKey := EncodeKey(key)

	if IsJsonPrimitive(value) {
		writer.Push(depth, encodedKey+": "+EncodePrimitive(value, delimiter))
	} else if IsJsonArray(value) {
		arr := toSlice(value)
		encodeArray(key, arr, writer, depth, delimiter, lengthMarker)
	} else if IsJsonObject(value) {
		obj := toMap(value)
		if len(obj) == 0 {
			// Empty object
			writer.Push(depth, encodedKey+":")
		} else {
			writer.Push(depth, encodedKey+":")
			encodeObject(obj, writer, depth+1, delimiter, lengthMarker)
		}
	}
}

// encodeArray encodes a JSON array.
func encodeArray(key string, value []interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	if len(value) == 0 {
		opts := &HeaderOptions{Key: key, Delimiter: delimiter, LengthMarker: lengthMarker}
		header := FormatHeader(0, opts)
		writer.Push(depth, header)
		return
	}

	// Primitive array
	if IsArrayOfPrimitives(value) {
		formatted := encodeInlineArrayLine(value, delimiter, key, lengthMarker)
		writer.Push(depth, formatted)
		return
	}

	// Array of arrays (all primitives)
	if IsArrayOfArrays(value) {
		allPrimitiveArrays := true
		for _, item := range value {
			arr := toSlice(item)
			if !IsArrayOfPrimitives(arr) {
				allPrimitiveArrays = false
				break
			}
		}
		if allPrimitiveArrays {
			encodeArrayOfArraysAsListItems(key, value, writer, depth, delimiter, lengthMarker)
			return
		}
	}

	// Array of objects
	if IsArrayOfObjects(value) {
		objects := make([]map[string]interface{}, len(value))
		for i, item := range value {
			objects[i] = toMap(item)
		}

		header := extractTabularHeader(objects)
		if header != nil {
			encodeArrayOfObjectsAsTabular(key, objects, header, writer, depth, delimiter, lengthMarker)
		} else {
			encodeMixedArrayAsListItems(key, value, writer, depth, delimiter, lengthMarker)
		}
		return
	}

	// Mixed array: fallback to expanded format
	encodeMixedArrayAsListItems(key, value, writer, depth, delimiter, lengthMarker)
}

// encodeArrayOfArraysAsListItems encodes an array of primitive arrays as list items.
func encodeArrayOfArraysAsListItems(prefix string, values []interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	opts := &HeaderOptions{Key: prefix, Delimiter: delimiter, LengthMarker: lengthMarker}
	header := FormatHeader(len(values), opts)
	writer.Push(depth, header)

	for _, item := range values {
		arr := toSlice(item)
		if IsArrayOfPrimitives(arr) {
			inline := encodeInlineArrayLine(arr, delimiter, "", lengthMarker)
			writer.PushListItem(depth+1, inline)
		}
	}
}

// encodeInlineArrayLine encodes a primitive array on a single line.
func encodeInlineArrayLine(values []interface{}, delimiter rune, prefix string, lengthMarker bool) string {
	opts := &HeaderOptions{Key: prefix, Delimiter: delimiter, LengthMarker: lengthMarker}
	header := FormatHeader(len(values), opts)
	joinedValue := EncodeAndJoinPrimitives(values, delimiter)

	// Only add space if there are values
	if len(values) == 0 {
		return header
	}
	return header + " " + joinedValue
}

// encodeArrayOfObjectsAsTabular encodes an array of objects in tabular format.
func encodeArrayOfObjectsAsTabular(prefix string, rows []map[string]interface{}, header []string, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	opts := &HeaderOptions{
		Key:          prefix,
		Fields:       header,
		Delimiter:    delimiter,
		LengthMarker: lengthMarker,
	}
	formattedHeader := FormatHeader(len(rows), opts)
	writer.Push(depth, formattedHeader)

	writeTabularRows(rows, header, writer, depth+1, delimiter)
}

// extractTabularHeader extracts field names if the array can be encoded in tabular format.
func extractTabularHeader(rows []map[string]interface{}) []string {
	if len(rows) == 0 {
		return nil
	}

	firstRow := rows[0]
	if len(firstRow) == 0 {
		return nil
	}

	// Get keys from first row
	firstKeys := make([]string, 0, len(firstRow))
	for key := range firstRow {
		firstKeys = append(firstKeys, key)
	}

	if isTabularArray(rows, firstKeys) {
		return firstKeys
	}

	return nil
}

// isTabularArray checks if an array of objects can be encoded in tabular format.
func isTabularArray(rows []map[string]interface{}, header []string) bool {
	for _, row := range rows {
		// All objects must have the same keys (but order can differ)
		if len(row) != len(header) {
			return false
		}

		// Check that all header keys exist in the row and all values are primitives
		for _, key := range header {
			val, exists := row[key]
			if !exists {
				return false
			}
			if !IsJsonPrimitive(val) {
				return false
			}
		}
	}

	return true
}

// writeTabularRows writes the data rows for a tabular array.
func writeTabularRows(rows []map[string]interface{}, header []string, writer *LineWriter, depth int, delimiter rune) {
	for _, row := range rows {
		values := make([]interface{}, len(header))
		for i, key := range header {
			values[i] = row[key]
		}
		joinedValue := EncodeAndJoinPrimitives(values, delimiter)
		writer.Push(depth, joinedValue)
	}
}

// encodeMixedArrayAsListItems encodes a mixed array as list items.
func encodeMixedArrayAsListItems(prefix string, items []interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	opts := &HeaderOptions{Key: prefix, Delimiter: delimiter, LengthMarker: lengthMarker}
	header := FormatHeader(len(items), opts)
	writer.Push(depth, header)

	for _, item := range items {
		encodeListItemValue(item, writer, depth+1, delimiter, lengthMarker)
	}
}

// encodeObjectAsListItem encodes an object as a list item.
func encodeObjectAsListItem(obj map[string]interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	if len(obj) == 0 {
		writer.Push(depth, string(listItemMarker))
		return
	}

	// Get keys (order matters for consistency)
	keys := make([]string, 0, len(obj))
	for key := range obj {
		keys = append(keys, key)
	}

	// First key-value on the same line as "- "
	firstKey := keys[0]
	encodedKey := EncodeKey(firstKey)
	firstValue := obj[firstKey]

	if IsJsonPrimitive(firstValue) {
		writer.PushListItem(depth, encodedKey+": "+EncodePrimitive(firstValue, delimiter))
	} else if IsJsonArray(firstValue) {
		arr := toSlice(firstValue)
		if IsArrayOfPrimitives(arr) {
			// Inline format for primitive arrays
			formatted := encodeInlineArrayLine(arr, delimiter, firstKey, lengthMarker)
			writer.PushListItem(depth, formatted)
		} else if IsArrayOfObjects(arr) {
			objects := make([]map[string]interface{}, len(arr))
			for i, item := range arr {
				objects[i] = toMap(item)
			}

			// Check if array of objects can use tabular format
			header := extractTabularHeader(objects)
			if header != nil {
				// Tabular format for uniform arrays of objects
				opts := &HeaderOptions{
					Key:          firstKey,
					Fields:       header,
					Delimiter:    delimiter,
					LengthMarker: lengthMarker,
				}
				formattedHeader := FormatHeader(len(arr), opts)
				writer.PushListItem(depth, formattedHeader)
				writeTabularRows(objects, header, writer, depth+1, delimiter)
			} else {
				// Fall back to list format for non-uniform arrays of objects
				writer.PushListItem(depth, encodedKey+"["+string(rune(len(arr)))+"]:")
				for _, item := range arr {
					encodeObjectAsListItem(toMap(item), writer, depth+1, delimiter, lengthMarker)
				}
			}
		} else {
			// Complex arrays on separate lines (array of arrays, etc.)
			writer.PushListItem(depth, encodedKey+"["+string(rune(len(arr)))+"]:")

			// Encode array contents at depth + 1
			for _, item := range arr {
				encodeListItemValue(item, writer, depth+1, delimiter, lengthMarker)
			}
		}
	} else if IsJsonObject(firstValue) {
		nestedObj := toMap(firstValue)
		if len(nestedObj) == 0 {
			writer.PushListItem(depth, encodedKey+":")
		} else {
			writer.PushListItem(depth, encodedKey+":")
			encodeObject(nestedObj, writer, depth+2, delimiter, lengthMarker)
		}
	}

	// Remaining keys on indented lines
	for i := 1; i < len(keys); i++ {
		key := keys[i]
		encodeKeyValuePair(key, obj[key], writer, depth+1, delimiter, lengthMarker)
	}
}

// encodeListItemValue encodes a value as a list item.
func encodeListItemValue(value interface{}, writer *LineWriter, depth int, delimiter rune, lengthMarker bool) {
	if IsJsonPrimitive(value) {
		writer.PushListItem(depth, EncodePrimitive(value, delimiter))
	} else if IsJsonArray(value) {
		arr := toSlice(value)
		if IsArrayOfPrimitives(arr) {
			inline := encodeInlineArrayLine(arr, delimiter, "", lengthMarker)
			writer.PushListItem(depth, inline)
		}
	} else if IsJsonObject(value) {
		obj := toMap(value)
		encodeObjectAsListItem(obj, writer, depth, delimiter, lengthMarker)
	}
}

// Helper functions to convert interface{} to concrete types

func toSlice(value interface{}) []interface{} {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil
	}

	result := make([]interface{}, v.Len())
	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}
	return result
}

func toMap(value interface{}) map[string]interface{} {
	if value == nil {
		return nil
	}

	if m, ok := value.(map[string]interface{}); ok {
		return m
	}

	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Map {
		return nil
	}

	result := make(map[string]interface{})
	iter := v.MapRange()
	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		keyStr := ""
		if key.Kind() == reflect.String {
			keyStr = key.String()
		} else {
			keyStr = key.String() // Fallback
		}
		result[keyStr] = val.Interface()
	}
	return result
}
