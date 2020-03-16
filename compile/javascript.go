package compile

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var _ Refiner = &RefineJavascript{}
var reJavascriptMatch = regexp.MustCompile(`(?i)\.(js|json|jsx)$`)

type RefineJavascript struct {
	passthrough
}

func (rf *RefineJavascript) Match(p string) bool {
	if reMinPass.MatchString(p) {
		return false // skip already minified assets
	}
	return reJavascriptMatch.MatchString(p)
}

func (rf *RefineJavascript) Debug(destination, source string) error {
	return rf.compile(source, destination, true)
}

func (rf *RefineJavascript) Refine(destination, source string) error {
	return rf.compile(source, destination, false)
}

func (rf *RefineJavascript) compile(file, output string, debug bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel() // interrupts the process running
	args := []string{
		filepath.Base(file),
		`--outfile=` + output,
		`--bundle`,
	}
	if !debug {
		args = append(args, `--minify`)
	}
	p := exec.CommandContext(ctx, `esbuild`, args...)

	if debug {
		p.Stderr = os.Stderr
	}

	p.Dir = filepath.Dir(file)
	return p.Run()
}
