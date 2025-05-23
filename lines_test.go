package dotenv

import "testing"

func TestParseLine(t *testing.T) {
	cases := []struct {
		In    string
		Key   string
		Value string
		Error error
	}{
		{"", "", "", ErrEmptyLine},
		{"# fully a comment line", "", "", ErrEmptyLine}, // should already be covered by the tests in comments_test.go but one is never too sure.
		{"line with to many \"\"quotes\"", "", "", ErrTooManyQuotes},
		{"\"good\" bad", "", "", ErrContentOutsideQuotes},
		{"k\"ey=value", "", "", ErrQuoteInKey},
		{"bad", "", "", ErrMalformedLine},
		{"key=value", "key", "value", nil},
		{"key=\"value\"", "key", "value", nil},
		{"key=value # comment", "key", "value", nil},
		{"key=\"value\" # comment", "key", "value", nil},
		{"key=\"value\\nline2\"", "key", "value\nline2", nil},
		{"key=\"value\nline2\"", "key", "value\nline2", nil},
		{"key=\"value\\nline2\" # comment", "key", "value\nline2", nil},
		{"key=\"value\nline2\" # comment", "key", "value\nline2", nil},
		{"key=ab\"value\"", "", "", ErrContentOutsideQuotes},
	}

	for _, c := range cases {
		c := c
		t.Run(c.In, func(t *testing.T) {
			t.Parallel()

			key, value, err := parseLine(c.In)
			if key != c.Key {
				t.Errorf("wanted key '%s', got '%s' (line: '%s')", c.Key, key, c.In)
			}

			if value != c.Value {
				t.Errorf("wanted value '%s', got '%s' (line: '%s')", c.Value, value, c.In)
			}

			if err != c.Error {
				t.Errorf("wanted error '%s', got '%s' (line: '%s')", c.Error, err, c.In)
			}
		})
	}
}
