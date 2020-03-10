package minify

import (
	"io"

	libsass "github.com/wellington/go-libsass"
)

// SCSS compiles and minifies style sheets.
func SCSS(w io.Writer, r io.Reader, includes ...string) error {
	comp, err := libsass.New(w, r)
	if err != nil {
		return err
	}
	injectConf(comp, includes)
	comp.Option(libsass.WithSyntax(libsass.SCSSSyntax))
	return comp.Run()
}
