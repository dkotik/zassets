package minify

// All collects all minification tools into one object.

import (
	"os"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/svg"
)

var (
	Minifier      = minify.New()
	debug    bool = os.Getenv(`DEBUG`) != ``
)

func init() {
	// TODO: remove legacy application/javascript minifier?
	Minifier.AddFunc("application/javascript", js.Minify) // archaic
	Minifier.AddFunc("text/javascript", js.Minify)
	Minifier.AddFunc("text/html", html.Minify)
	Minifier.AddFunc("text/svg+xml", svg.Minify)
	Minifier.AddFunc("text/css", css.Minify)
	// Minifier.AddFunc("application/jpeg", )
	// Minifier.AddFunc("application/png", )
	// Minifier.AddFunc("application/webp", )
}
