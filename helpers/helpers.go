// This package contains a few helper functions to assist in parsing strings.
package helpers

// RuneIndexN returns the index of the n-th instance of the rune r in the string s.
// If n < 1, the function always returns -1.
func RuneIndexN(s string, r rune, n int) (idx int) {
	if n < 1 {
		return -1
	}
	count := 0
	for i, sr := range s {
		if sr == r {
			count++
			if count == n {
				return i
			}
		}
	}
	return -1
}

// CountRuneN goes through the string s looking for rune r, either up until it has seen r n times or it reaches the end of the string.
// It returns the amount of times it has seen r.
// If n < 0, it will count all instances of r until it reaches the end of the string.
// if n == 0, it will immediately return 0.
func CountRuneN(s string, r rune, n int) (res int) {
	if n == 0 {
		return
	}
	for _, sr := range s {
		if sr == r {
			res++
			if n > 0 && res >= n {
				return
			}
		}

	}
	return
}
