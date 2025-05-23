package dotenv

import "testing"

func TestStripComment(t *testing.T) {
	cases := []struct {
		In  string
		Out string
	}{
		{"aabbcc", "aabbcc"},
		{"", ""},
		{"# this line is all comments no data", ""},
		{"aabbcc # a comment", "aabbcc "}, // the space at the end is important! We're stripping comments, not extra space
		{"\"quoted string\"", "\"quoted string\""},
		{"\"# not a comment\"", "\"# not a comment\""},
		{"\"#not a comment\" # yes a comment", "\"#not a comment\" "}, // the space is also important here!
	}

	for _, c := range cases {
		c := c
		t.Run(c.In, func(t *testing.T) {
			t.Parallel()
			res := stripComments(c.In)
			if res != c.Out {
				t.Errorf("wanted '%s', got '%s' (input: '%s')", c.Out, res, c.In)
			}
		})
	}
}
