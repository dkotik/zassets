package debug

import (
	"log"

	"github.com/dkotik/zassets/compile"
)

// // WithDefaultOptions configures standard Debugger behavior.
// func WithDefaultOptions() func(d *Debugger) error {
// 	return func(d *Debugger) (err error) {
// 		err = WithLogger()(d)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	}
// }

// WithDefaultCompiler connects a default Compiler to Debugger.
func WithDefaultCompiler() func(d *Debugger) error {
	return func(d *Debugger) error {
		c, err := compile.NewCompiler(compile.WithDefaultOptions(), compile.WithDebug())
		d.c = c
		return err
	}
}

// WithCompiler connects a Compiler to Debugger.
func WithCompiler(c *compile.Compiler) func(d *Debugger) error {
	return func(d *Debugger) error {
		d.c = c
		return nil
	}
}

// WithLogger overrides the default Debugger logger.
func WithLogger(l *log.Logger) func(d *Debugger) error {
	return func(d *Debugger) error {
		d.l = l
		if d.c != nil {
			compile.WithLogger(l)(d.c)
		}
		return nil
	}
}
