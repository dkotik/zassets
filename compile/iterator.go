package compile

import (
	"os"
	"path/filepath"
	"regexp"
)

// IteratorFunc is run on every Iterator cycle.
// First returned error stops the Iterator.
type IteratorFunc func(target, relative string, info os.FileInfo) error

// NewIterator sets up an Iterator.
// If the Ignore list is empty, fills it with default values.
func NewIterator(entries, ignore []string) (*Iterator, error) {
	i := &Iterator{entries, make([]*regexp.Regexp, 0)}
	if len(ignore) == 0 {
		ignore = []string{
			`(\A|\\|\/)[\.\_][^\\\/]+$`,
			`node_modules$`,
		}
	}
	for _, p := range ignore {
		r, err := regexp.Compile(p)
		if err != nil {
			return i, err
		}
		i.Ignore = append(i.Ignore, r)
	}
	return i, nil
}

// Iterator returns a walking function through paths pointing to files and folders.
type Iterator struct {
	Entries []string // paths to files and folders
	Ignore  []*regexp.Regexp
}

// AllowPath returns true, if the path does not match any of the Ignore patterns.
func (i *Iterator) AllowPath(p string) bool {
	for _, r := range i.Ignore {
		if r.MatchString(p) {
			return false
		}
	}
	return true
}

// Walk calls the IteratorFunc for every discovered object.
func (i *Iterator) Walk(f IteratorFunc) error {
	for _, p := range i.Entries {
		if !i.AllowPath(p) {
			continue
		}
		info, err := os.Stat(p)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = filepath.Walk(p, func(s string, info os.FileInfo, err error) error {
				if info.IsDir() || err != nil {
					if !i.AllowPath(s) {
						return filepath.SkipDir
					}
					return err
				}
				if !i.AllowPath(s) {
					return nil
				}
				// chop off walk root and filepath separator
				return f(s, s[len(p)+1:], info)
			})
			if err != nil {
				return err
			}
			continue
		}
		err = f(p, filepath.Base(p), info)
		if err != nil {
			return err
		}
	}
	return nil
}
