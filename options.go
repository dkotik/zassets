package goresminpack

import (
	"fmt"
	"os"
	"regexp"
)

type option func(*Compiler) error

// OptIgnore eliminates all matching files from processing.
// Ignored files will still be accessible as includes.
func OptIgnore(patterns ...string) option {
	return func(c *Compiler) error {
		for _, p := range patterns {
			r, err := regexp.Compile(p)
			if err != nil {
				return err
			}
			c.ignore = append(c.ignore, r)
		}
		return nil
	}
}

// OptInclude adds an additional path to look for includes, when neccessary.
func OptInclude(paths ...string) option {
	return func(c *Compiler) error {
		for _, p := range paths {
			s, err := os.Stat(p)
			if err != nil {
				return err
			}
			if !s.IsDir() {
				return fmt.Errorf("%s is not a directory", p)
			}
			c.include = append(c.include, p)
		}
		return nil
	}
}

// OptDebug presents the directory files in readable format.
func OptDebug() option {
	return func(c *Compiler) error {
		c.debug = true
		return nil
	}
}
