package compile

import (
	"fmt"
	"os"
	"regexp"
)

// OptRefiners links refiners to the compiler.
// They will run in the order added for each file.
func OptRefiners(refiners ...Refiner) func(c *Compiler) error {
	return func(c *Compiler) error {
		for _, r := range refiners {
			c.refiners = append(c.refiners, r)
		}
		return nil
	}
}

// OptIgnore eliminates all matching files from processing.
// Ignored files will still be accessible as includes.
func OptIgnore(patterns ...string) func(c *Compiler) error {
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
func OptInclude(paths ...string) func(c *Compiler) error {
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
func OptDebug() func(c *Compiler) error {
	return func(c *Compiler) error {
		c.debug = true
		return nil
	}
}
