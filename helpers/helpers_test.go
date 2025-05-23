package helpers

import "testing"

func TestRuneIndexN(t *testing.T) {
	cases := []struct {
		In     string
		Needle rune
		N      int
		Want   int
	}{
		{"aba", 'a', 2, 2},
		{"bab", 'c', 2, -1},
		{"bba", 'a', 2, -1},
		{"", 'a', 2, -1},
		{"abb", 'a', 0, -1},
	}

	for _, c := range cases {
		c := c
		t.Run(c.In, func(t *testing.T) {
			t.Parallel()
			res := RuneIndexN(c.In, c.Needle, c.N)
			if res != c.Want {
				t.Errorf("wanted %d, got %d (line: '%s')", c.Want, res, c.In)
			}
		})
	}
}

func TestCountRuneN(t *testing.T) {
	cases := []struct {
		In     string
		Needle rune
		N      int
		Want   int
	}{
		{"aba", 'a', 2, 2},
		{"abaa", 'a', 2, 2},
		{"aaba", 'a', -1, 3},
		{"aaab", 'c', 2, 0},
		{"baaa", 'a', 0, 0},
	}

	for _, c := range cases {
		c := c
		t.Run(c.In, func(t *testing.T) {
			t.Parallel()
			res := CountRuneN(c.In, c.Needle, c.N)

			if res != c.Want {
				t.Errorf("wanted %d, got %d (line: '%s')", c.Want, res, c.In)
			}
		})
	}
}
