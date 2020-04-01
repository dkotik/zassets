package compile

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"time"
)

// Refiner converts one type of asset into another or optimizes the content.
type Refiner interface {
	// Match will activate the refiner on matching paths.
	Match(path string) (ok bool)
	// Some refiners transform one asset into another, like SASS => CSS.
	Rename(oldPath string) (newPath string)
	// Change the content from source to destination.
	Refine(destination, source string) error
	// Same as Refine, but keep the changed files as readable as possible.
	Debug(destination, source string) error
}

// NewCompiler creates a configured compiler.
func NewCompiler(opts ...func(*Compiler) error) (*Compiler, error) {
	c := &Compiler{
		false,
		make([]string, 0),
		make([]*regexp.Regexp, 0),
		make([]Refiner, 0),
		log.New(os.Stdout, `📁 `, log.Ltime|log.Lmsgprefix),
		50,
	}
	var err error
	for _, o := range opts {
		if err = o(c); err != nil {
			return c, err
		}
	}
	return c, nil
}

// Compiler packs resources for distribution.
type Compiler struct {
	debug    bool
	include  []string
	ignore   []*regexp.Regexp
	refiners []Refiner
	logger   *log.Logger
	maxTasks int
}

// Run gathers and compiles assets from files and folders in given paths.
func (c *Compiler) Run(destination string, i *Iterator) (err error) {
	errs := make(chan error, c.maxTasks)
	tasks := make(chan string, c.maxTasks)
	i.Walk(func(source, relative string, info os.FileInfo) error {
		tasks <- source
		if len(errs) > 0 { // there is at least one error on stack
			<-tasks // end the task
			return errors.New(`iteration interrupted`)
		}
		go func() {
			// for i := 0; i < 1+rand.Intn(3); i++ { // for testing
			// 	time.Sleep(time.Second)
			// }
			if err := c.each(filepath.Join(destination, relative), source); err != nil {
				if c.debug {
					c.logger.Printf("Error in %s: %s", source, err)
				}
				errs <- err
			}
			<-tasks
		}()
		return nil
	})

	var remaining int
	for { // wait on workers
		time.Sleep(time.Second)
		if remaining = len(tasks); remaining > 0 {
			if c.debug {
				c.logger.Printf("Waiting for %d tasks...", remaining)
			}
		} else {
			break
		}
	}
	close(errs)
	close(tasks)
	return <-errs // return only the first error that occured
}

// Each searches for the first matching Refiner and runs it.
// If no refiner is discovered, the source file is copied.
func (c *Compiler) each(destination, source string) (err error) {
	err = os.MkdirAll(filepath.Dir(destination), 0700)
	if err != nil {
		return err
	}
	for _, r := range c.refiners {
		if r.Match(source) {
			if c.debug {
				c.logger.Printf("Refining %s to %s using %s.", source, r.Rename(destination), reflect.TypeOf(r))
				return r.Debug(r.Rename(destination), source)
			}
			err = r.Refine(r.Rename(destination), source)
			if err != nil {
				return fmt.Errorf(`could not refine %s: %w`, source, err)
			}
			return nil // stop processing
		}
	}
	if c.debug {
		c.logger.Printf("Copying %s to %s.", source, destination)
	}
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, r)
	return err
}
