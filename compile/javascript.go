package compile

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var _ Refiner = &RefineJavascript{}
var reJavascriptMatch = regexp.MustCompile(`(?i)\.(js|json|jsx)$`)

// RefineJavascript compiles a Javascript file to ESNext.
type RefineJavascript struct {
	passthrough
}

// Match returns true if pattern fits the file path.
func (rf *RefineJavascript) Match(p string) bool {
	if reMinPass.MatchString(p) {
		return false // skip already minified assets
	}
	return reJavascriptMatch.MatchString(p)
}

// Debug preserves the comments and keeps Javascript readable.
func (rf *RefineJavascript) Debug(destination, source string) error {
	return rf.compile(source, destination, true)
}

// Refine runs the compilation and minification.
func (rf *RefineJavascript) Refine(destination, source string) error {
	return rf.compile(source, destination, false)
}

func (rf *RefineJavascript) compile(file, output string, debug bool) (err error) {
	_, err = exec.LookPath(`esbuild`)
	if err != nil {
		return errors.New(`"esbuild" javascript compiler is not installed`)
	}

	// spew.Dump(exec.LookPath(`esbuild`))
	// spew.Dump(exec.LookPath(`go`))

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

	var b bytes.Buffer
	p.Stderr = &b
	p.Stdout = &b
	p.Dir = filepath.Dir(file)
	err = p.Run()
	if err != nil {
		return fmt.Errorf(`%s: %s`, filepath.Dir(file), b.String())
	}
	return nil
}
