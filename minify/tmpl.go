package minify

import (
	"bytes"
	"io"
)

func Tmpl(w io.Writer, r io.Reader) error {
	// TODO: not very efficient here!
	var b bytes.Buffer
	_, err := io.Copy(&b, r)
	if err != nil {
		return nil
	}
	_, err = io.WriteString(w, reNewLineWithWhiteSpace.ReplaceAllString(b.String(), ""))
	return err
}
