//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package ghw

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	disableWarningsEnv = "GHW_DISABLE_WARNINGS"
)

type closer interface {
	Close() error
}

func safeClose(c closer) {
	err := c.Close()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to close: %s", err)
	}
}

func warn(msg string, args ...interface{}) {
	if _, ok := os.LookupEnv(disableWarningsEnv); ok {
		return
	}
	_, _ = fmt.Fprint(os.Stderr, "WARNING: ")
	_, _ = fmt.Fprintf(os.Stderr, msg, args...)
}

// Reads a supplied filepath and converts the contents to an integer. Returns
// -1 if there were file permissions or existence errors or if the contents
// could not be successfully converted to an integer. In any error, a warning
// message is printed to STDERR and -1 is returned.
func safeIntFromFile(path string) int {
	msg := "failed to read int from file: %s\n"
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		warn(msg, err)
		return -1
	}
	contents := strings.TrimSpace(string(buf))
	res, err := strconv.Atoi(contents)
	if err != nil {
		warn(msg, err)
		return -1
	}
	return res
}
