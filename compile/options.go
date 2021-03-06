package compile

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/chai2010/webp"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/svg"
)

// WithDefaultOptions configures standard compiler behavior.
func WithDefaultOptions() func(c *Compiler) error {
	return func(c *Compiler) (err error) {
		err = WithRefiners(
			&RefineMinify{ // Clean up HTML files.
				MatchPath: regexp.MustCompile(`(?i)\.html?$`),
				Minifier:  html.Minify,
			},
			&RefineJavascript{},
			// &RefineText{ // Strip extra white space from SQL and Tmpl files.
			// 	MatchPath: regexp.MustCompile(`(?i)\.(sql|tmpl)$`),
			// 	Search:    regexp.MustCompile(`\s*?\n\s*`),
			// 	Replace:   ``,
			// },
		)(c)
		if err != nil {
			return err
		}
		err = WithDefaultCSSRefiners()(c)
		if err != nil {
			return err
		}
		return WithDefaultImageRefiners()(c)
	}
}

// WithDefaultCSSRefiners sets up SASS and SCSS compilers and a CSS minifier.
func WithDefaultCSSRefiners() func(c *Compiler) error {
	return WithRefiners(
		&RefineMinify{ // Clean up CSS files.
			MatchPath: regexp.MustCompile(`(?i)\.css$`),
			Minifier:  css.Minify,
		},
		&RefineSASS{}, &RefineSCSS{},
	)
}

// WithDefaultImageRefiners sets up refiners for SVG, JPEG, PNG, and Webp images.
func WithDefaultImageRefiners() func(c *Compiler) error {
	return WithRefiners(
		&RefineMinify{ // Compress SVG files.
			MatchPath: regexp.MustCompile(`(?i)\.svg$`),
			Minifier:  svg.Minify,
		},
		&RefineRaster{ // Compress JPG files.
			MatchPath: regexp.MustCompile(`(?i)\.jpe?g$`),
			Encode: func(w io.Writer, img image.Image) error {
				return jpeg.Encode(w, img, &jpeg.Options{Quality: 60})
			},
			Decode: jpeg.Decode,
		},
		&RefineRaster{ // Compress PNG files.
			MatchPath: regexp.MustCompile(`(?i)\.png$`),
			Encode:    png.Encode,
			Decode:    png.Decode,
		},
		&RefineRaster{ // Compress Webp files.
			MatchPath: regexp.MustCompile(`(?i)\.webp$`),
			Encode: func(w io.Writer, img image.Image) error {
				return webp.Encode(w, img, &webp.Options{Quality: 60})
			},
			Decode: webp.Decode,
		},
	)
}

// WithRefiners links refiners to the compiler.
// They will run in the order added for each file.
func WithRefiners(refiners ...Refiner) func(c *Compiler) error {
	return func(c *Compiler) error {
		for _, r := range refiners {
			c.refiners = append(c.refiners, r)
		}
		return nil
	}
}

// WithParallelTasks sets the maximum number of simultaneous refiners.
// Higher number speeds up the compilation but consumes more resources.
// Default value is 50.
func WithParallelTasks(max int) func(c *Compiler) error {
	return func(c *Compiler) error {
		if max < 1 {
			return errors.New("cannot run less than 1 parallel tasks")
		}
		if max > 9999 {
			return errors.New(string(max) + " parallel tasks is too many")
		}
		c.maxTasks = max
		return nil
	}
}

// WithInclude adds an additional path to look for includes, when neccessary.
func WithInclude(paths ...string) func(c *Compiler) error {
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

// WithDebug presents the directory files in readable format.
func WithDebug() func(c *Compiler) error {
	return func(c *Compiler) error {
		c.debug = true
		return nil
	}
}

// WithLogger overrides the default Compiler logger.
func WithLogger(l *log.Logger) func(c *Compiler) error {
	return func(c *Compiler) error {
		c.logger = l
		return nil
	}
}
