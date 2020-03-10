package minify

import (
	"io"

	libsass "github.com/wellington/go-libsass"
)

func injectConf(comp libsass.Compiler, includes []string) {
	if debug {
		comp.Option(libsass.OutputStyle(libsass.EXPANDED_STYLE))
		comp.Option(libsass.Comments(true))
		comp.Option(libsass.LineComments(true))
	} else {
		comp.Option(libsass.OutputStyle(libsass.COMPRESSED_STYLE))
		comp.Option(libsass.Comments(false))
		comp.Option(libsass.LineComments(false))
	}
	comp.Option(libsass.IncludePaths(includes))
}

// SASS compiles and minifies style sheets written in SASS.
func SASS(w io.Writer, r io.Reader, includes ...string) error {
	comp, err := libsass.New(w, r)
	if err != nil {
		return err
	}
	injectConf(comp, includes)
	comp.Option(libsass.WithSyntax(libsass.SassSyntax))
	return comp.Run()
}
