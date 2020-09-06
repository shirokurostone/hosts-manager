package hosts

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	testCases := []struct {
		input    string
		expected *ParseResult
	}{
		{
			input: "192.0.2.0 host001.example\thost002.example# comment\n",
			expected: &ParseResult{
				Lines: []Line{
					Line{Tokens: []Token{
						Token{Type: Text, Text: "192.0.2.0"},
						Token{Type: Separator, Text: " "},
						Token{Type: Text, Text: "host001.example"},
						Token{Type: Separator, Text: "\t"},
						Token{Type: Text, Text: "host002.example"},
						Token{Type: Comment, Text: "# comment"},
					}},
				},
			},
		},
	}

	for i, testCase := range testCases {
		reader := strings.NewReader(testCase.input)
		actual, err := Parse(reader)
		if err != nil {
			t.Errorf("case:%d, expected:%#v, actual:%#v", i, nil, err)
		}
		if !reflect.DeepEqual(testCase.expected, actual) {
			t.Errorf("case:%d, expected:%#v, actual:%#v", i, testCase.expected, actual)
		}
	}

}

func TestCheckSyntax(t *testing.T) {
	testCases := []struct {
		input    []Token
		expected bool
	}{
		{
			input: []Token{
				Token{Type: Text, Text: "192.0.2.0"},
				Token{Type: Separator, Text: " "},
				Token{Type: Text, Text: "host001.example"},
				Token{Type: Separator, Text: "\t"},
				Token{Type: Text, Text: "host002.example"},
				Token{Type: Comment, Text: "# comment"},
			},
			expected: true,
		},
		{
			input: []Token{
				Token{Type: Text, Text: "2001:db8::1"},
				Token{Type: Separator, Text: " "},
				Token{Type: Text, Text: "host003.example"},
				Token{Type: Separator, Text: "\t"},
				Token{Type: Text, Text: "host004.example"},
				Token{Type: Comment, Text: "# comment"},
			},
			expected: true,
		},
		{
			input: []Token{
				Token{Type: Separator, Text: " "},
			},
			expected: true,
		},
		{
			input: []Token{
				Token{Type: Comment, Text: "# comment"},
			},
			expected: true,
		},
		{
			input: []Token{
				Token{Type: Text, Text: "host001.example"},
			},
			expected: false,
		},
		{
			input: []Token{
				Token{Type: Text, Text: "192.0.2.0"},
			},
			expected: false,
		},
		{
			input: []Token{
				Token{Type: Text, Text: "2001:db8::1"},
			},
			expected: false,
		},
	}

	for i, testCase := range testCases {
		parseResult := ParseResult{
			Lines: []Line{
				Line{
					Tokens: testCase.input,
				},
			},
		}

		actual := parseResult.CheckSyntax()
		if testCase.expected != actual {
			t.Errorf("case:%d, expected:%#v, actual:%#v", i, testCase.expected, actual)
		}
	}

}
