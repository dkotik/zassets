package compile

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Refiner converts one type of asset into another or optimizes the content.
type Refiner interface {
	Match(path string) (ok bool)
	Rename(oldPath string) (newPath string)
	Refine(destination, source string) error
	Debug(destination, source string) error
}

// NewCompiler creates a configured compiler.
func NewCompiler(opts ...func(*Compiler) error) (*Compiler, error) {
	c := &Compiler{
		false,
		make([]string, 0),
		make([]*regexp.Regexp, 0),
		make([]Refiner, 0),
		make(chan string, 50),
		make(chan error),
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

	tasks  chan string
	errors chan error
}

// Run gathers and compiles assets from files and folders in given paths.
func (c *Compiler) Run(destination string, paths ...string) (err error) {
	for _, p := range paths {
		// fmt.Println(`compiling assets in`, p, paths)
		if !c.allowPath(p) {
			continue
		}
		s, err := os.Stat(p)
		if err != nil {
			return err
		}
		if s.IsDir() {
			err = filepath.Walk(p, func(s string, i os.FileInfo, err error) error {
				if i.IsDir() || err != nil {
					if !c.allowPath(s) {
						return filepath.SkipDir
					}
					return err
				}
				if !c.allowPath(s) {
					return nil
				}
				relative := strings.TrimPrefix(s, p)
				return c.task(
					filepath.Join(destination, relative),
					filepath.Join(p, relative))
			})
			if err != nil {
				return err
			}
			continue
		}
		err = c.task(filepath.Join(destination, filepath.Base(p)), p)
		if err != nil {
			return err
		}
	}
	time.Sleep(time.Millisecond * 100)
	for { // wait on workers
		select {
		case err = <-c.errors:
			return err
		default:
		}
		if !c.working() {
			break
		}
	}
	return nil
}

func (c *Compiler) working() bool { // remaining tasks
	remaining := len(c.tasks)
	if remaining == 0 {
		return false
	}
	if c.debug {
		log.Printf("Waiting for %d tasks...", remaining)
	}
	time.Sleep(time.Second)
	return true
}

func (c *Compiler) task(destination, source string) (err error) {
	select {
	case err = <-c.errors:
		for c.working() {
		}
		return err
	default:
	}
	c.tasks <- source
	go func() {
		if err := c.each(destination, source); err != nil {
			c.errors <- err
		}
		// for i := 0; i < 1+rand.Intn(3); i++ {
		// 	time.Sleep(time.Second)
		// }
		<-c.tasks
	}()
	return nil
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
				log.Printf("Refining %s to %s using %s.", source, r.Rename(destination), reflect.TypeOf(r))
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
		log.Printf("Copying %s to %s.", source, destination)
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

func (c *Compiler) allowPath(p string) bool {
	for _, r := range c.ignore {
		if r.MatchString(p) {
			// log.Printf("Ignoring %s.\n", p)
			return false
		}
	}
	return true
}
