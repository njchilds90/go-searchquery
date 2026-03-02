// Package searchquery provides a deterministic lexer and parser for
// GitHub-style search queries, turning strings like `is:open "system crash"` 
// into structured Abstract Syntax Trees.
package searchquery

import (
	"strings"
	"unicode"
)

// Query represents the parsed search query.
type Query struct {
	Original string
	Terms    []*Term
}

// Term represents a single parsed element of the search query.
type Term struct {
	Value    string // The text or the value of the key-value pair
	Key      string // The key if it is a key-value pair, otherwise empty
	Exclude  bool   // True if the term was prefixed with a minus (-)
	IsPhrase bool   // True if the value was explicitly quoted
}

// Parse takes a search string and returns a structured Query.
// It is fully deterministic, stateless, and contains zero external dependencies.
func Parse(input string) *Query {
	q := &Query{
		Original: strings.TrimSpace(input),
		Terms:    make([]*Term, 0),
	}

	if q.Original == "" {
		return q
	}

	runes := []rune(q.Original)
	length := len(runes)
	index := 0

	for index < length {
		// Skip whitespace
		for index < length && unicode.IsSpace(runes[index]) {
			index++
		}
		if index >= length {
			break
		}

		exclude := false
		if runes[index] == '-' {
			// Check if the next character is not a space, meaning it is an exclusion modifier
			if index+1 < length && !unicode.IsSpace(runes[index+1]) {
				exclude = true
				index++
			}
		}

		key := ""
		// Look ahead to check if the sequence is a key:value pair
		// Keys should consist of alphanumeric characters, hyphens, or underscores
		lookahead := index
		for lookahead < length && (unicode.IsLetter(runes[lookahead]) || unicode.IsDigit(runes[lookahead]) || runes[lookahead] == '_' || runes[lookahead] == '-') {
			lookahead++
		}
		
		if lookahead > index && lookahead < length && runes[lookahead] == ':' {
			key = string(runes[index:lookahead])
			index = lookahead + 1 // Skip the colon character
		}

		value := ""
		isPhrase := false

		// Parse explicitly quoted phrases
		if index < length && runes[index] == '"' {
			isPhrase = true
			index++ // Skip opening quote
			start := index
			for index < length {
				if runes[index] == '\\' && index+1 < length && runes[index+1] == '"' {
					// Handle escaped quotes
					index += 2
					continue
				}
				if runes[index] == '"' {
					break
				}
				index++
			}
			// Unescape the captured value
			rawVal := string(runes[start:index])
			value = strings.ReplaceAll(rawVal, `\"`, `"`)
			if index < length && runes[index] == '"' {
				index++ // Skip closing quote
			}
		} else {
			// Read standard text until the next whitespace
			start := index
			for index < length && !unicode.IsSpace(runes[index]) {
				index++
			}
			value = string(runes[start:index])
		}

		// Ensure we do not append empty trailing artifacts
		if value != "" || isPhrase {
			q.Terms = append(q.Terms, &Term{
				Key:      key,
				Value:    value,
				Exclude:  exclude,
				IsPhrase: isPhrase,
			})
		}
	}

	return q
}

// Texts returns all positive free-text search terms.
func (q *Query) Texts() []string {
	var results []string
	for _, term := range q.Terms {
		if term.Key == "" && !term.Exclude {
			results = append(results, term.Value)
		}
	}
	return results
}

// Get returns all positive values associated with a specific key.
func (q *Query) Get(key string) []string {
	var results []string
	for _, term := range q.Terms {
		if strings.EqualFold(term.Key, key) && !term.Exclude {
			results = append(results, term.Value)
		}
	}
	return results
}

// Excludes returns all terms that have been marked for exclusion.
func (q *Query) Excludes() []*Term {
	var results []*Term
	for _, term := range q.Terms {
		if term.Exclude {
			results = append(results, term)
		}
	}
	return results
}

// String reconstructs the term back into a formatted string.
func (t *Term) String() string {
	result := ""
	if t.Exclude {
		result += "-"
	}
	if t.Key != "" {
		result += t.Key + ":"
	}
	value := t.Value
	if t.IsPhrase || strings.ContainsAny(value, " \t\n\r\"") {
		value = `"` + strings.ReplaceAll(value, `"`, `\"`) + `"`
	}
	result += value
	return result
}

// String reconstructs the entire AST back into a deterministic search string.
func (q *Query) String() string {
	var parts []string
	for _, term := range q.Terms {
		parts = append(parts, term.String())
	}
	return strings.Join(parts, " ")
}
