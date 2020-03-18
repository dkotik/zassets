package compile

import (
	"bytes"
	"io"
	"os"
	"regexp"
)

var _ Refiner = &RefineText{}

type RefineText struct {
	MatchPath *regexp.Regexp
	Search    *regexp.Regexp
	Replace   string
	passthrough
}

func (rf *RefineText) Match(p string) bool { return rf.MatchPath.MatchString(p) }

func (rf *RefineText) Refine(destination, source string) error {
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
	// TODO: not very efficient here!
	var b bytes.Buffer
	_, err = io.Copy(&b, r)
	if err != nil {
		return nil
	}
	_, err = io.WriteString(w, rf.Search.ReplaceAllString(b.String(), rf.Replace))
	return err
}

type passthrough struct{}

func (p *passthrough) Match(s string) bool { return true }

func (p *passthrough) Debug(destination, source string) error {
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
	_, err = io.Copy(w, r)
	return err
}

func (p *passthrough) Rename(path string) string { return path }
