package compile

import (
	"os"
	"regexp"

	"github.com/tdewolff/minify"
)

// Interfacing with github.com/tdewolff/minify.

var _ Refiner = &RefineMinify{}
var reMinPass = regexp.MustCompile(`(?i)\.min\.[^\.]+$`)

// RefineMinify interfaces with the popular Minifier library.
type RefineMinify struct {
	passthrough

	MatchPath *regexp.Regexp
	Minifier  minify.MinifierFunc
}

// Match returns true if pattern fits the file path.
func (rf *RefineMinify) Match(p string) bool {
	if reMinPass.MatchString(p) {
		return false // skip already minified assets
	}
	return rf.MatchPath.MatchString(p)
}

// Refine applies a minifier to source and writes the result to destination.
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
