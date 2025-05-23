package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/melsincostan/dotenv/helpers"
)

const (
	quoteChar   = '"'
	commentChar = '#'
	lineBreak   = "\\n"
	separator   = '='
)

var (
	ErrTooManyQuotes        = errors.New("too many quotes in a single value")
	ErrContentOutsideQuotes = errors.New("content after quoted string")
	ErrEmptyLine            = errors.New("no content on this line")
	ErrMalformedLine        = errors.New("line isn't in the key=value format and isn't a comment")
	ErrQuoteInKey           = fmt.Errorf("key contains '%c'", quoteChar)
)

func ParseFile(name string) (res map[string]string, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	return Parse(scanner)
}

func Parse(scanner *bufio.Scanner) (res map[string]string, err error) {
	inquotes := false
	multiline := ""
	rawLines := []string{}
	for scanner.Scan() {
		pl := scanner.Text()
		switch helpers.CountRuneN(pl, quoteChar, 2) {
		case 0:
			if inquotes {
				multiline += lineBreak + pl
			} else {
				rawLines = append(rawLines, pl)
			}
		case 1:
			if inquotes {
				rawLines = append(rawLines, multiline+lineBreak+pl)
				inquotes = false
			} else {
				// check that the quote isn't in the key part of the string, which would be problematic
				if strings.IndexRune(pl, quoteChar) < strings.IndexRune(pl, separator) {
					return nil, ErrQuoteInKey
				}
				multiline = pl
				inquotes = true
			}
		case 2:
			if inquotes {
				return nil, ErrTooManyQuotes
			} else {
				// check that the quote isn't in the key part of the string, which would be problematic
				if strings.IndexRune(pl, quoteChar) < strings.IndexRune(pl, separator) {
					return nil, ErrQuoteInKey
				}
				rawLines = append(rawLines, pl)
			}
		}
	}

	res = map[string]string{}

	for _, line := range rawLines {
		key, value, err := parseLine(line)
		if err != nil {
			if errors.Is(err, ErrEmptyLine) {
				// doesn't interest us, continue
				continue
			}
			return nil, err
		}
		res[key] = value
	}
	return
}

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

	spl := strings.SplitN(cl, string(separator), 2)
	if len(spl) < 2 {
		return "", "", ErrMalformedLine
	}

	key = spl[0]
	// remove quotes from the string.
	// It was checked before that there should only be a valid amount, so doing it this way shouldn't result in trimming the beginning of a string.
	noquotes := strings.TrimSuffix(strings.TrimPrefix(spl[1], string(quoteChar)), string(quoteChar))
	if len(noquotes) < len(spl[1]) { // if this is a quoted string
		noquotes = strings.ReplaceAll(noquotes, lineBreak, "\n") // replace all "line breaks" with actual line breaks
	}
	value = noquotes
	return
}
