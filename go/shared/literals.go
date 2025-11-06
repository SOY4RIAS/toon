package shared

import (
	"strconv"
	"strings"
)

// IsBooleanOrNullLiteral checks if a token is a boolean or null literal
// (true, false, null)
func IsBooleanOrNullLiteral(token string) bool {
	return token == "true" || token == "false" || token == "null"
}

// IsNumericLiteral checks if a token represents a valid numeric literal
// Rejects numbers with leading zeros (except "0" itself or decimals like "0.5")
func IsNumericLiteral(token string) bool {
	if token == "" {
		return false
	}

	// Must not have leading zeros (except for "0" itself or decimals like "0.5")
	if len(token) > 1 && token[0] == '0' && token[1] != '.' {
		// Also check for negative numbers like "-0.5"
		if !strings.HasPrefix(token, "-") {
			return false
		}
	}

	// Check if it's a valid number
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}
