package dotenv

import (
	"strings"

	"github.com/melsincostan/dotenv/helpers"
)

// parseLine takes in a line from the Parse function (a multiline value would be a single entry)! and parses it to a key and value.
// If there is an issue with the line, a corresponding error is returned and the key and value are left empty.
// An [ErrEmptyLine] error is returned if the line is empty to distinguish from an empty key and value.
// TODO: this should maybe also be an error in the future.
func parseLine(line string) (key, value string, err error) {
	cl := strings.TrimSpace(stripComments(line))

	if len(cl) < 1 { // no need to parse an empty line, break
		return "", "", ErrEmptyLine // not really an error, it's just that there needs to be a way to indicate stuff to callers and this feels better than nil?
	}

	if helpers.CountRuneN(cl, quoteChar, 3) > 2 { // check that there isn't a stray extra quote.
		// limit to three, doesn't matter if there are more since it is already an error to have three.
		// maybe this can help be a tad bit faster if its a very long line and there is an error in the beginning already.
		// isn't constant time though but would there be a timing attack on this??
		return "", "", ErrTooManyQuotes
	}

	if lqidx := helpers.RuneIndexN(cl, quoteChar, 2); lqidx != -1 && lqidx < (len(cl)-1) { // check for stray, non comment content after the last quote
		return "", "", ErrContentOutsideQuotes
	}

	// check that the first quote isn't in the key part of the string, which would be problematic
	if fqidx := strings.IndexRune(cl, quoteChar); fqidx != -1 && fqidx < strings.IndexRune(cl, separator) {
		return "", "", ErrQuoteInKey
	}

	spl := strings.SplitN(cl, string(separator), 2)
	if len(spl) < 2 {
		return "", "", ErrMalformedLine
	}

	if strings.IndexRune(spl[1], quoteChar) != -1 {
		// the string is quoted
		key = strings.TrimSpace(spl[0])

		nspval := strings.TrimSpace(spl[1])

		if nspval[0] != quoteChar { // check that there is no content before the first quote
			return "", "", ErrContentOutsideQuotes
		}
		// remove quotes from the string and replace "\n" by actual line breaks, since this is a multiline string.
		// It was checked before that there should only be a valid amount, so doing it this way shouldn't result in trimming the beginning of a string.
		value = strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(nspval, string(quoteChar)), string(quoteChar)), lineBreak, "\n")
	} else {
		// unquoted string
		key = spl[0]
		value = spl[1]
	}

	return
}
