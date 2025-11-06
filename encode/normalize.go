package encode

import (
	"math"
	"reflect"
	"time"
)

// NormalizeValue converts any Go value to a JSON-compatible value.
// This handles special cases like time.Time, NaN, Infinity, etc.
func NormalizeValue(value interface{}) interface{} {
	// Handle nil
	if value == nil {
		return nil
	}

	// Use reflection to handle different types
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String:
		return value

	case reflect.Bool:
		return value

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())

	case reflect.Float32, reflect.Float64:
		f := v.Float()
		// Handle -0
		if f == 0 && math.Signbit(f) {
			return 0.0
		}
		// Handle NaN and Infinity
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil
		}
		return f

	case reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		// Dereference pointer
		return NormalizeValue(v.Elem().Interface())

	case reflect.Slice, reflect.Array:
		// Handle []byte specially (convert to string)
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes())
		}

		// Convert to []interface{}
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = NormalizeValue(v.Index(i).Interface())
		}
		return result

	case reflect.Map:
		// Convert to map[string]interface{}
		result := make(map[string]interface{})
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			val := iter.Value()
			// Convert key to string
			keyStr := ""
			switch key.Kind() {
			case reflect.String:
				keyStr = key.String()
			default:
				keyStr = key.String() // Use String() method for other types
			}
			result[keyStr] = NormalizeValue(val.Interface())
		}
		return result

	case reflect.Struct:
		// Special case: time.Time
		if t, ok := value.(time.Time); ok {
			return t.Format(time.RFC3339Nano)
		}

		// Convert struct to map
		result := make(map[string]interface{})
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			// Skip unexported fields
			if !field.IsExported() {
				continue
			}

			// Get field name (check for json tag first)
			fieldName := field.Name
			if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
				// Simple tag parsing (just get the name part before comma)
				for commaIdx := 0; commaIdx < len(tag); commaIdx++ {
					if tag[commaIdx] == ',' {
						fieldName = tag[:commaIdx]
						break
					}
				}
				if fieldName == field.Name {
					// No comma found, use entire tag
					fieldName = tag
				}
			}

			fieldValue := v.Field(i)
			result[fieldName] = NormalizeValue(fieldValue.Interface())
		}
		return result

	default:
		// For unsupported types (func, chan, etc.), return null
		return nil
	}
}

// IsJsonPrimitive checks if a value is a JSON primitive (string, number, boolean, or null).
func IsJsonPrimitive(value interface{}) bool {
	if value == nil {
		return true
	}

	switch value.(type) {
	case string, bool, float64, float32, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	default:
		return false
	}
}

// IsJsonArray checks if a value is a JSON array (slice).
func IsJsonArray(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array
}

// IsJsonObject checks if a value is a JSON object (map).
func IsJsonObject(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Map
}

// IsArrayOfPrimitives checks if a slice contains only primitive values.
func IsArrayOfPrimitives(value []interface{}) bool {
	for _, item := range value {
		if !IsJsonPrimitive(item) {
			return false
		}
	}
	return true
}

// IsArrayOfArrays checks if a slice contains only arrays.
func IsArrayOfArrays(value []interface{}) bool {
	for _, item := range value {
		if !IsJsonArray(item) {
			return false
		}
	}
	return true
}

// IsArrayOfObjects checks if a slice contains only objects (maps).
func IsArrayOfObjects(value []interface{}) bool {
	for _, item := range value {
		if !IsJsonObject(item) {
			return false
		}
	}
	return true
}
