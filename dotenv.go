package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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

func Parse(reader io.Reader) (res map[string]string, err error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines) // explicitely set the function we want to use to split.
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
