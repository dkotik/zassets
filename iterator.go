package zassets

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Iterator returns a walking function through paths pointing to files and folders.
func Iterator(entries, ignore []string) (func(func(target, relative string) error) error, error) {
	ir := make([]*regexp.Regexp, 0)
	for _, p := range ignore {
		r, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		ir = append(ir, r)
	}
	allowPath := func(p string) bool {
		for _, r := range ir {
			if r.MatchString(p) {
				return false
			}
		}
		return true
	}

	return func(f func(target, relative string) error) error {
		for _, p := range entries {
			if !allowPath(p) {
				continue
			}
			s, err := os.Stat(p)
			if err != nil {
				return err
			}
			if s.IsDir() {
				err = filepath.Walk(p, func(s string, i os.FileInfo, err error) error {
					if i.IsDir() || err != nil {
						if !allowPath(s) {
							return filepath.SkipDir
						}
						return err
					}
					if !allowPath(s) {
						return nil
					}
					return f(s, strings.TrimPrefix(s, p))
				})
				if err != nil {
					return err
				}
				continue
			}
			err = f(p, filepath.Base(p))
			if err != nil {
				return err
			}
		}
		return nil
	}, nil
}
