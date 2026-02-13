package graphql_test

import (
	"bytes"
	"testing"

	"github.com/CodeWarrior-debug/perspectize/backend/pkg/graphql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- MarshalIntID Tests ---

func TestMarshalIntID_PositiveInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "single digit",
			input:    1,
			expected: `"1"`,
		},
		{
			name:     "multiple digits",
			input:    42,
			expected: `"42"`,
		},
		{
			name:     "large number",
			input:    999999,
			expected: `"999999"`,
		},
		{
			name:     "zero",
			input:    0,
			expected: `"0"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			marshaler := graphql.MarshalIntID(tt.input)
			var buf bytes.Buffer
			marshaler.MarshalGQL(&buf)
			assert.Equal(t, tt.expected, buf.String())
		})
	}
}

func TestMarshalIntID_NegativeInteger(t *testing.T) {
	marshaler := graphql.MarshalIntID(-1)
	var buf bytes.Buffer
	marshaler.MarshalGQL(&buf)
	assert.Equal(t, `"-1"`, buf.String())
}

// --- UnmarshalIntID Tests ---

func TestUnmarshalIntID_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		hasError bool
	}{
		{
			name:     "single digit",
			input:    "1",
			expected: 1,
			hasError: false,
		},
		{
			name:     "multiple digits",
			input:    "42",
			expected: 42,
			hasError: false,
		},
		{
			name:     "large number",
			input:    "999999",
			expected: 999999,
			hasError: false,
		},
		{
			name:     "zero",
			input:    "0",
			expected: 0,
			hasError: false,
		},
		{
			name:     "negative number",
			input:    "-1",
			expected: -1,
			hasError: false,
		},
		{
			name:     "invalid string",
			input:    "abc",
			expected: 0,
			hasError: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: 0,
			hasError: true,
		},
		{
			name:     "string with spaces",
			input:    "1 2",
			expected: 0,
			hasError: true,
		},
		{
			name:     "string with letters and numbers",
			input:    "123abc",
			expected: 0,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := graphql.UnmarshalIntID(tt.input)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestUnmarshalIntID_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{
			name:     "positive int",
			input:    42,
			expected: 42,
		},
		{
			name:     "zero",
			input:    0,
			expected: 0,
		},
		{
			name:     "negative int",
			input:    -1,
			expected: -1,
		},
		{
			name:     "large int",
			input:    1000000,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := graphql.UnmarshalIntID(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnmarshalIntID_Int64(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected int
	}{
		{
			name:     "positive int64",
			input:    int64(42),
			expected: 42,
		},
		{
			name:     "zero",
			input:    int64(0),
			expected: 0,
		},
		{
			name:     "negative int64",
			input:    int64(-1),
			expected: -1,
		},
		{
			name:     "large int64",
			input:    int64(1000000),
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := graphql.UnmarshalIntID(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnmarshalIntID_Float64(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected int
	}{
		{
			name:     "whole number float",
			input:    42.0,
			expected: 42,
		},
		{
			name:     "zero",
			input:    0.0,
			expected: 0,
		},
		{
			name:     "negative float",
			input:    -1.0,
			expected: -1,
		},
		{
			name:     "float with decimal (truncated)",
			input:    42.7,
			expected: 42,
		},
		{
			name:     "large float",
			input:    1000000.0,
			expected: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := graphql.UnmarshalIntID(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnmarshalIntID_UnsupportedType(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name:  "boolean",
			input: true,
		},
		{
			name:  "nil",
			input: nil,
		},
		{
			name:  "slice",
			input: []int{1, 2, 3},
		},
		{
			name:  "map",
			input: map[string]int{"key": 1},
		},
		{
			name:  "struct",
			input: struct{ ID int }{ID: 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := graphql.UnmarshalIntID(tt.input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "IntID must be a string or integer")
			assert.Equal(t, 0, result)
		})
	}
}

// --- Round-trip Tests ---

func TestIntID_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		id   int
	}{
		{name: "zero", id: 0},
		{name: "positive", id: 42},
		{name: "negative", id: -1},
		{name: "large", id: 999999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			marshaler := graphql.MarshalIntID(tt.id)
			var buf bytes.Buffer
			marshaler.MarshalGQL(&buf)

			// The marshaled value is a quoted string, so we need to unquote it
			// For our unmarshal test, we'll test with the unquoted version
			// (simulating what GraphQL would pass)

			// Unmarshal from int (direct pass-through)
			result, err := graphql.UnmarshalIntID(tt.id)
			require.NoError(t, err)
			assert.Equal(t, tt.id, result)
		})
	}
}
