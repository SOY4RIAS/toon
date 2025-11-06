package encode

import (
	"math"
	"math/big"
	"reflect"
	"time"
)

// NormalizeValue converts a Go value to a JSON-compatible value.
// It handles various Go types and normalizes them to JSON primitives, objects, or arrays.
//
// Type conversions:
// - nil → null
// - bool, string → unchanged
// - int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64 → float64
// - float32, float64 → float64 (with special handling for -0, NaN, Infinity)
// - *big.Int → number (if safe) or string
// - time.Time → ISO 8601 string
// - []T → array
// - map[string]T → object
// - struct → object (exported fields only)
// - pointer → dereference and normalize
// - function, channel, complex → null
func NormalizeValue(value interface{}) interface{} {
	if value == nil {
		return nil
	}

	v := reflect.ValueOf(value)

	// Handle pointers by dereferencing
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return NormalizeValue(v.Elem().Interface())
	}

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()

	case reflect.String:
		return v.String()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint())

	case reflect.Float32, reflect.Float64:
		f := v.Float()
		// Canonicalize -0 to 0
		if f == 0 && math.Signbit(f) {
			return float64(0)
		}
		// Convert NaN and Infinity to null
		if math.IsNaN(f) || math.IsInf(f, 0) {
			return nil
		}
		return f

	case reflect.Slice, reflect.Array:
		// Handle []byte specially (don't convert to array of numbers)
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return string(v.Bytes())
		}

		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = NormalizeValue(v.Index(i).Interface())
		}
		return result

	case reflect.Map:
		// Only support map[string]T
		if v.Type().Key().Kind() != reflect.String {
			return nil
		}

		result := make(map[string]interface{})
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			result[key] = NormalizeValue(iter.Value().Interface())
		}
		return result

	case reflect.Struct:
		// Special handling for time.Time
		if t, ok := value.(time.Time); ok {
			return t.Format(time.RFC3339Nano)
		}

		// Special handling for *big.Int
		if bigInt, ok := value.(*big.Int); ok {
			// Try to convert to float64 if within safe integer range
			if bigInt.IsInt64() {
				i64 := bigInt.Int64()
				if i64 >= -9007199254740991 && i64 <= 9007199254740991 { // Number.MIN_SAFE_INTEGER to MAX_SAFE_INTEGER
					return float64(i64)
				}
			}
			// Otherwise convert to string
			return bigInt.String()
		}

		// Convert struct to map (only exported fields)
		result := make(map[string]interface{})
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			// Only include exported fields
			if field.PkgPath == "" {
				fieldValue := v.Field(i)
				// Use json tag if present, otherwise use field name
				name := field.Name
				if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
					// Simple tag parsing (doesn't handle all json tag options)
					for idx := 0; idx < len(jsonTag); idx++ {
						if jsonTag[idx] == ',' {
							name = jsonTag[:idx]
							break
						}
					}
					if name == "" {
						name = jsonTag
					}
				}
				result[name] = NormalizeValue(fieldValue.Interface())
			}
		}
		return result

	case reflect.Func, reflect.Chan, reflect.Complex64, reflect.Complex128, reflect.UnsafePointer:
		return nil

	default:
		return nil
	}
}

// IsJsonPrimitive checks if a value is a JSON primitive (null, bool, number, string)
func IsJsonPrimitive(value interface{}) bool {
	if value == nil {
		return true
	}

	switch value.(type) {
	case bool, float64, string:
		return true
	default:
		return false
	}
}

// IsJsonArray checks if a value is a JSON array (slice)
func IsJsonArray(value interface{}) bool {
	if value == nil {
		return false
	}
	v := reflect.ValueOf(value)
	return v.Kind() == reflect.Slice
}

// IsJsonObject checks if a value is a JSON object (map[string]interface{})
func IsJsonObject(value interface{}) bool {
	if value == nil {
		return false
	}
	_, ok := value.(map[string]interface{})
	return ok
}

// IsArrayOfPrimitives checks if a value is an array of primitives
func IsArrayOfPrimitives(value interface{}) bool {
	if !IsJsonArray(value) {
		return false
	}

	arr := value.([]interface{})
	for _, item := range arr {
		if !IsJsonPrimitive(item) {
			return false
		}
	}
	return true
}

// IsArrayOfArrays checks if a value is an array of arrays
func IsArrayOfArrays(value interface{}) bool {
	if !IsJsonArray(value) {
		return false
	}

	arr := value.([]interface{})
	for _, item := range arr {
		if !IsJsonArray(item) {
			return false
		}
	}
	return true
}

// IsArrayOfObjects checks if a value is an array of objects
func IsArrayOfObjects(value interface{}) bool {
	if !IsJsonArray(value) {
		return false
	}

	arr := value.([]interface{})
	for _, item := range arr {
		if !IsJsonObject(item) {
			return false
		}
	}
	return true
}
