package minify

import (
	"bytes"
	"io"
	"regexp"
)

var reNewLineWithWhiteSpace = regexp.MustCompile(`\s*?\n\s*`)

// EatLineWhiteSpace eliminates new lines and surrounding white space.
func EatLineWhiteSpace(w io.Writer, r io.Reader) error {
	// TODO: not very efficient here!
	var b bytes.Buffer
	_, err := io.Copy(&b, r)
	if err != nil {
		return nil
	}
	_, err = io.WriteString(w, reNewLineWithWhiteSpace.ReplaceAllString(b.String(), ""))
	return err
}
