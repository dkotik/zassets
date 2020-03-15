package compile

import (
	"os"
	"regexp"

	"github.com/tdewolff/minify"
)

// Interfacing with github.com/tdewolff/minify.

var _ Refiner = &RefineMinify{}
var reMinPass = regexp.MustCompile(`(?i)\.min\.[^\.]+$`)

type RefineMinify struct {
	passthrough

	MatchPath *regexp.Regexp
	Minifier  minify.MinifierFunc
}

func (rf *RefineMinify) Match(p string) bool {
	if reMinPass.MatchString(p) {
		return false // skip already minified assets
	}
	return rf.MatchPath.MatchString(p)
}

func (rf *RefineMinify) Refine(destination, source string) error {
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
	return rf.Minifier(nil, w, r, nil)
}
