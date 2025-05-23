package dotenv

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		In  string
		Out map[string]string
		Err error
	}{
		{"key=value", map[string]string{"key": "value"}, nil},
		{"key=\"multi\nline\"", map[string]string{"key": "multi\nline"}, nil},
		{"key=multi\nline", nil, ErrMalformedLine},
		{"key=\"val\"ue\"", nil, ErrTooManyQuotes},
		{"key=\"long\nmulti\nline\"", map[string]string{"key": "long\nmulti\nline"}, nil},
		{"key=\"multi\nkey=\"value\"", nil, ErrTooManyQuotes},
		{"# comment only line\nkey=value", map[string]string{"key": "value"}, nil},
		{"key=\"multi\n#inthere\"", map[string]string{"key": "multi\n#inthere"}, nil},
		{"", map[string]string{}, nil},
		{"key=\"multi\nline\" # comment", map[string]string{"key": "multi\nline"}, nil},
		{"ke\"y=test\"", nil, ErrQuoteInKey},                // technically already covered in the test for the parseLine() function.
		{"key=bwaa\"value\"", nil, ErrContentOutsideQuotes}, // ditto
		{"key=\"value\"bwaa", nil, ErrContentOutsideQuotes}, // ditto
		{"# top level comment\nkey=\"value\" # line level comment\nkey2=\"multi\\nline\"\nkey3=\"multi\nline\"\nkey4=test", map[string]string{
			"key":  "value",
			"key2": "multi\nline",
			"key3": "multi\nline",
			"key4": "test",
		}, nil},
		{"key=unquoted\\ntest", map[string]string{"key": "unquoted\\ntest"}, nil}, // shouldn't get converted to a newline since it isn't quoted
	}

	for _, c := range cases {
		c := c
		t.Run(c.In, func(t *testing.T) {
			t.Parallel()
			res, err := Parse(bytes.NewReader([]byte(c.In)))
			if err != c.Err {
				t.Errorf("wanted error '%s', got '%s' (input: '%s')", c.Err, err, c.In)
			}

			if !reflect.DeepEqual(res, c.Out) {
				t.Errorf("wanted %#v, got %#v\n(input: '%s')", c.Out, res, c.In)
			}
		})
	}
}
