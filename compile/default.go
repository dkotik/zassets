package compile

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"regexp"

	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/svg"
)

// DefaultRefiners provides a queue of most common asset refiners.
var DefaultRefiners = []Refiner{
	&RefineMinify{ // Clean up HTML files.
		MatchPath: regexp.MustCompile(`(?i)\.html?$`),
		Minifier:  html.Minify,
	},
	&RefineJavascript{},
	&RefineMinify{ // Clean up CSS files.
		MatchPath: regexp.MustCompile(`(?i)\.css$`),
		Minifier:  css.Minify,
	},
	&RefineSASS{}, &RefineSCSS{},
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
	&RefineText{ // Strip extra white space from SQL and Tmpl files.
		MatchPath: regexp.MustCompile(`(?i)\.(sql|tmpl)$`),
		Search:    regexp.MustCompile(`\s*?\n\s*`),
		Replace:   ``,
	},
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
