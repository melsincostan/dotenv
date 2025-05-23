// This package provides a way to parse a .env file (or any [io.Reader] returning content akin to it) into a map with string keys and string values.
// It supports comments, quoted values and multiline strings.
package dotenv

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/melsincostan/dotenv/helpers"
)

const (
	quoteChar   = '"'
	commentChar = '#'
	lineBreak   = "\\n"
	separator   = '='
)

var (
	// ErrTooManyQuotes means a line had more than two " characters.
	ErrTooManyQuotes = errors.New("too many quotes in a single value")
	// ErrContentOutsideQuotes means a line had a value with quotes, but there is extra content before or after the quote.
	ErrContentOutsideQuotes = errors.New("content after quoted string")
	// ErrEmptyLine means a line is empty or is entirely composed of a comment.
	// This is used to differentiate between an empty line and a line with an empty key and value but an equal sign, and should only be used internally.
	ErrEmptyLine = errors.New("no content on this line")
	// ErrMalformedLine means the format of the line couldn't be parsed to a key and a value, due to the lack of an equals sign.
	ErrMalformedLine = errors.New("line isn't in the key=value format and isn't a comment")
	// ErrQuoteInKey means that there is a quote sign in the key part of the line.
	ErrQuoteInKey = fmt.Errorf("key contains '%c'", quoteChar)
)

// Parse takes in an [io.Reader] representation of a .env file and parses through it, returning a map keys and values.
// No parsing is done on the values beyond setting "\n" to an actual newline in quoted values and removing quotes.
// It returns no map and an error in case a line has an unexpected format.
// It supportes unquoted values, quoted values, quoted multiline values by using "\n" on a single line, multiline values and comments, either on their own line or at the end of a line.
// Double quotes shouldn't be used as part of keys or values, as this will return an error.
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
				multiline = pl
				inquotes = true
			}
		case 2:
			if inquotes {
				return nil, ErrTooManyQuotes
			} else {
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
