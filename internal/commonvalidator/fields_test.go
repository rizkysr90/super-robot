package commonvalidator

import "testing"

func TestIsRequired(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		// String tests
		{
			name:     "non-empty string",
			value:    "hello",
			expected: true,
		},
		{
			name:     "empty string",
			value:    "",
			expected: false,
		},
		{
			name:     "whitespace string",
			value:    "   ",
			expected: false,
		},

		// Integer tests
		{
			name:     "positive int",
			value:    42,
			expected: true,
		},
		{
			name:     "zero int",
			value:    0,
			expected: false,
		},
		{
			name:     "negative int",
			value:    -1,
			expected: true,
		},
		{
			name:     "int64 value",
			value:    int64(100),
			expected: true,
		},

		// Slice tests
		{
			name:     "non-empty interface slice",
			value:    []interface{}{"hello", 123},
			expected: true,
		},
		{
			name:     "empty interface slice",
			value:    []interface{}{},
			expected: false,
		},
		{
			name:     "non-empty string slice",
			value:    []string{"hello", "world"},
			expected: true,
		},
		{
			name:     "empty string slice",
			value:    []string{},
			expected: false,
		},
		{
			name:     "non-empty int slice",
			value:    []int{1, 2, 3},
			expected: true,
		},
		{
			name:     "empty int slice",
			value:    []int{},
			expected: false,
		},

		// Map tests
		{
			name:     "non-empty map",
			value:    map[string]interface{}{"key": "value"},
			expected: true,
		},
		{
			name:     "empty map",
			value:    map[string]interface{}{},
			expected: false,
		},

		// Nil tests
		{
			name:     "nil value",
			value:    nil,
			expected: false,
		},

		// Custom struct tests - structs are always considered valid if not nil
		{
			name:     "empty struct",
			value:    struct{}{},
			expected: true,
		},
		{
			name:     "struct with empty fields",
			value:    struct{ name string }{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRequired(tt.value)
			if result != tt.expected {
				t.Errorf("IsRequired(%v) = %v; want %v", tt.value, result, tt.expected)
			}
		})
	}
}

// Test for different numeric types that fall into the default case
func TestIsRequiredDefaultNumericTypes(t *testing.T) {
	numericTests := []struct {
		name     string
		value    interface{}
		expected bool
	}{
		{"uint non-zero", uint(42), true},
		{"uint zero", uint(0), true}, // Falls to default case
		{"float64 non-zero", 3.14, true},
		{"float64 zero", 0.0, true}, // Falls to default case
		{"float32 non-zero", float32(3.14), true},
		{"float32 zero", float32(0.0), true}, // Falls to default case
	}

	for _, tt := range numericTests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRequired(tt.value)
			if result != tt.expected {
				t.Errorf("IsRequired(%v) = %v; want %v", tt.value, result, tt.expected)
			}
		})
	}
}
func TestMaxLen(t *testing.T) {
	tests := []struct {
		name      string
		str       string
		maxLength int
		expected  bool
	}{
		{
			name:      "empty string",
			str:       "",
			maxLength: 5,
			expected:  true,
		},
		{
			name:      "exact length",
			str:       "hello",
			maxLength: 5,
			expected:  true,
		},
		{
			name:      "shorter than max",
			str:       "hi",
			maxLength: 5,
			expected:  true,
		},
		{
			name:      "longer than max",
			str:       "hello world",
			maxLength: 5,
			expected:  false,
		},
		{
			name:      "zero max length with empty string",
			str:       "",
			maxLength: 0,
			expected:  true,
		},
		{
			name:      "zero max length with non-empty string",
			str:       "a",
			maxLength: 0,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaxLen(tt.str, tt.maxLength)
			if result != tt.expected {
				t.Errorf("MaxLen(%q, %d) = %v; want %v",
					tt.str, tt.maxLength, result, tt.expected)
			}
		})
	}
}
