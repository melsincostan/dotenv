# .env format parser written in go

This is likely not the fastest or most efficient version out there, but it works.
Since the `Parse()` function takes an `io.Reader`-compatible value, it can be used directly with files, network connections (why would one do that), and of course simple strings (going through `bytes.NewReader()`, for example).

Features:

- [x] basic `key=value` lines
- [x] comments (both on their own line and inline)
- [x] quoted values
- [x] multiline values (on a single line and as actual multiline values)
- [ ] support for `'` as quote in addition to `"` (might require slight changes to how `Parse()` internally prepares lines)
- [ ] support for quoted keys

Example of using a file and a string:

```golang
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/melsincostan/dotenv"
)

func main() {
	file, err := os.Open(".env")
	if err != nil {
		log.Fatalf("Error opening the file: %s", err.Error())
	}

	parsed, err := dotenv.Parse(file)
	if err != nil {
		log.Fatalf("Error parsing the file: %s", err.Error())
	}

	fmt.Printf("%#v\n", parsed)

	raw := `test="value"
other_key="multi
line
value"`
	reader := bytes.NewReader([]byte(raw))
	parsed, err = dotenv.Parse(reader)
	if err != nil {
		log.Fatalf("Error parsing the value: %s", err.Error())
	}
	fmt.Printf("%#v\n", parsed)
}
```