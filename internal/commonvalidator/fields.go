package commonvalidator

import "strings"

// IsRequired checks if a value is empty based on its type
func IsRequired(value interface{}) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case string:
		return len(strings.TrimSpace(v)) > 0
	case int, int8, int16, int32, int64:
		return v != 0
	case uint, uint8, uint16, uint32, uint64:
		return v != 0
	case float32, float64:
		return v != 0
	case bool:
		return v
	case []interface{}:
		return len(v) > 0
	case map[string]interface{}:
		return len(v) > 0
	case []string:
		return len(v) > 0
	case []int:
		return len(v) > 0
	default:
		// For custom types or structs, just check if not nil
		return true
	}
}

// MaxLen checks if a string's length is less than or equal to the maximum allowed length
func MaxLen(str string, maxLength int) bool {
	return len(str) <= maxLength
}
