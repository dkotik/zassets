package goresminpack

import (
	"fmt"
	"goresminpack/minify"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/svg"
)

// Compiler packs resources for distribution.
type Compiler struct {
	debug   bool
	include []string
	ignore  []*regexp.Regexp
}

// WithOptions applies options that will alter the minifiers.
func (c *Compiler) WithOptions(opts ...option) (err error) {
	for _, o := range opts {
		if err = o(c); err != nil {
			return err
		}
	}
	return nil
}

// Process compiles and minifies assets from the source directory to the destination directory.
func (c *Compiler) Process(source, destination string) (err error) {
	d, err := os.Stat(destination)
	if err != nil {
		return err
	}
	if !d.IsDir() {
		return fmt.Errorf("path %s is not a directory", destination)
	}

	err = filepath.Walk(source, filepath.WalkFunc(func(path string, info os.FileInfo, err error) error {
		if !c.allowPath(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return err
		} else if info.IsDir() {
			return nil
		}

		// fullPath := filepath.Join(source, path)
		ext := filepath.Ext(path)
		lext := strings.ToLower(ext)
		target := filepath.Join(destination, strings.TrimPrefix(path, source))
		// log.Fatalln(target)
		err = os.MkdirAll(filepath.Dir(target), 0700)
		if err != nil {
			return err
		}
		switch lext {
		case `.js`, `.jsx`, `.json`:
			target, err = filepath.Abs(target)
			if err != nil {
				return err
			}
			return minify.CompiledJS(path, target, c.debug)
		case `.sass`, `.scss`:
			target = strings.TrimSuffix(target, ext) + `.css`
		}

		r, err := os.Open(path)
		if err != nil {
			return err
		}
		defer r.Close()
		w, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
		defer w.Close()

		switch lext {
		case `.html`:
			err = html.Minify(nil, w, r, nil)
		case `.css`:
			err = css.Minify(nil, w, r, nil)
		case `.scss`:
			err = minify.SCSS(w, r, c.includesWithCurrent(filepath.Dir(path))...)
		case `.sass`:
			err = minify.SASS(w, r, c.includesWithCurrent(filepath.Dir(path))...)
		case `.svg`:
			err = svg.Minify(nil, w, r, nil)
		case `.tmpl`, `.sql`:
			err = minify.EatLineWhiteSpace(w, r)
		case `.jpg`, `.jpeg`:
			err = minify.ResizeJPG(w, r)
		case `.png`:
			err = minify.ResizePNG(w, r)
		default:
			_, err = io.Copy(w, r)
		}
		return err
	}))
	return err
}

func (c *Compiler) includesWithCurrent(p string) []string {
	return append([]string{p}, c.include...)
}

func (c *Compiler) allowPath(p string) bool {
	for _, r := range c.ignore {
		if r.MatchString(p) {
			return false
		}
	}
	return true
}
