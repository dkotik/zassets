package compile

import (
	"os"
	"path/filepath"
	"regexp"

	libsass "github.com/wellington/go-libsass"
)

var _ Refiner = &RefineSASS{}
var _ Refiner = &RefineSCSS{}

var reMatchSCSS = regexp.MustCompile(`(?i)\.scss$`)
var reMatchSASS = regexp.MustCompile(`(?i)\.sass$`)

// RefineSCSS compiles SCSS files to minified CSS.
type RefineSCSS struct {
	Paths []string
}

func (rf *RefineSCSS) Rename(p string) string {
	return p[:len(p)-4] + `css`
}

func (rf *RefineSCSS) Match(p string) bool {
	return reMatchSCSS.MatchString(p)
}

func (rf *RefineSCSS) prepare(comp libsass.Compiler, debug bool) {
	if debug {
		comp.Option(libsass.OutputStyle(libsass.EXPANDED_STYLE))
		comp.Option(libsass.Comments(true))
		comp.Option(libsass.LineComments(true))
	} else {
		comp.Option(libsass.OutputStyle(libsass.COMPRESSED_STYLE))
		comp.Option(libsass.Comments(false))
		comp.Option(libsass.LineComments(false))
	}
	comp.Option(libsass.IncludePaths(rf.Paths))
}

func (rf *RefineSCSS) Refine(destination, source string) error {
	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	comp, err := libsass.New(w, r)
	if err != nil {
		return err
	}
	rf.prepare(comp, false)
	comp.Option(libsass.IncludePaths([]string{filepath.Dir(source)}))
	return comp.Run()
}

func (rf *RefineSCSS) Debug(destination, source string) error {
	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	comp, err := libsass.New(w, r)
	if err != nil {
		return err
	}
	rf.prepare(comp, true)
	return comp.Run()
}

// RefineSASS compiles SASS files to minified CSS.
type RefineSASS struct {
	RefineSCSS
}

func (rf *RefineSASS) Match(p string) bool {
	return reMatchSASS.MatchString(p)
}

func (rf *RefineSASS) Refine(destination, source string) error {
	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	comp, err := libsass.New(w, r)
	if err != nil {
		return err
	}
	rf.prepare(comp, false)
	comp.Option(libsass.IncludePaths([]string{filepath.Dir(source)}))
	comp.Option(libsass.WithSyntax(libsass.SassSyntax))
	return comp.Run()
}

func (rf *RefineSASS) Debug(destination, source string) error {
	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer w.Close()
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	comp, err := libsass.New(w, r)
	if err != nil {
		return err
	}
	rf.prepare(comp, true)
	comp.Option(libsass.WithSyntax(libsass.SassSyntax))
	return comp.Run()
}
