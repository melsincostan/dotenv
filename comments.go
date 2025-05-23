package dotenv

// stripComments takes in a line (as returned by [Parse] - single entry for a multiline value), and removes any potential comments from it.
func stripComments(line string) (cleaned string) {
	cutoff := len(line)
	idx := len(line)
out:
	for range len(line) {
		idx--
		switch line[idx] {
		case commentChar:
			cutoff = idx
		case quoteChar:
			break out // break breaks out of select as well, break out of for instead using a label
		}
	}

	if cutoff < 1 {
		return ""
	}

	return line[:cutoff]
}
