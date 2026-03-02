package searchquery_test

import (
	"testing"

	"github.com/njchilds90/go-searchquery"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []searchquery.Term
	}{
		{
			name:  "basic text",
			input: "hello world",
			expected: []searchquery.Term{
				{Value: "hello"},
				{Value: "world"},
			},
		},
		{
			name:  "key values",
			input: "is:open label:bug",
			expected: []searchquery.Term{
				{Key: "is", Value: "open"},
				{Key: "label", Value: "bug"},
			},
		},
		{
			name:  "exclusion",
			input: "-label:stale -\"ignored phrase\"",
			expected: []searchquery.Term{
				{Key: "label", Value: "stale", Exclude: true},
				{Value: "ignored phrase", Exclude: true, IsPhrase: true},
			},
		},
		{
			name:  "complex query",
			input: `author:alice "panic error" -status:closed priority:-1`,
			expected: []searchquery.Term{
				{Key: "author", Value: "alice"},
				{Value: "panic error", IsPhrase: true},
				{Key: "status", Value: "closed", Exclude: true},
				{Key: "priority", Value: "-1"},
			},
		},
		{
			name:  "escaped quotes",
			input: `message:"invalid \"token\" error"`,
			expected: []searchquery.Term{
				{Key: "message", Value: `invalid "token" error`, IsPhrase: true},
			},
		},
		{
			name:  "empty string",
			input: "   ",
			expected: []searchquery.Term{},
		},
		{
			name:  "standalone dash",
			input: "hello - world",
			expected: []searchquery.Term{
				{Value: "hello"},
				{Value: "-"},
				{Value: "world"},
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			query := searchquery.Parse(testCase.input)
			if len(query.Terms) != len(testCase.expected) {
				t.Fatalf("expected %d terms, got %d", len(testCase.expected), len(query.Terms))
			}
			for index, term := range query.Terms {
				expected := testCase.expected[index]
				if term.Key != expected.Key || term.Value != expected.Value || term.Exclude != expected.Exclude || term.IsPhrase != expected.IsPhrase {
					t.Errorf("term %d mismatch: expected %+v, got %+v", index, expected, term)
				}
			}
		})
	}
}

func TestStringReconstruction(t *testing.T) {
	input := `author:alice "panic error" -status:closed priority:-1`
	query := searchquery.Parse(input)
	reconstructed := query.String()
	
	if reconstructed != input {
		t.Errorf("expected %q, got %q", input, reconstructed)
	}
}
