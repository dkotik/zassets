package compile

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Debugger watches sources and re-compile them to serve through the asset object.
type Debugger struct {
	d string
	c *Compiler
	i *Iterator
	w *fsnotify.Watcher
	l *log.Logger
}

// SetLogger changes the logger for all the messages that the Debugger produces.
func (d *Debugger) SetLogger(l *log.Logger) {
	if l == nil {
		l = log.New(os.Stdout, `📁 `, log.Ltime|log.Lmsgprefix)
	}
	d.l = l
	// if d.c != nil {
	WithLogger(d.l)(d.c)
	// }
}

// Watch observes source objects for changes.
// Call multiple times, if additional objects need to be added.
func (d *Debugger) Watch(p ...string) (err error) {
	for _, dir := range p {
		err = filepath.Walk(dir, func(p string, i os.FileInfo, err error) error {
			d.w.Add(p)
			if i.IsDir() {
			}
			return err
		})
		if err != nil {
			return err
		}
		d.l.Printf("👁 Watching %s for changes.", dir)
	}
	return nil
}

func (d *Debugger) watch() { // the boring watcher logic
	t := time.Tick(time.Second)
	var somethingChanged bool
	for {
		select {
		case <-t: // time to update
			if somethingChanged {
				somethingChanged = false
				err := d.c.Run(d.d, d.i)
				if err != nil {
					d.l.Println("Refining error:", err.Error())
				}
			}
		case event, ok := <-d.w.Events:
			if !ok {
				d.l.Println("error: could not gather the event")
				return
			}
			somethingChanged = true
			d.l.Println("Detected a changed file:", event.Name)
			if event.Op&fsnotify.Create == fsnotify.Create {
				d.w.Add(event.Name)
			}
		case err, ok := <-d.w.Errors:
			d.l.Println("Error:", err)
			if !ok {
				return
			}
		}
	}
}

// Open fulfills the http.FileSystem interface.
// If there is no associated compiler, point to a file in source directory.
// Otherwise, point to a temporary file that was built using the compiler.
func (d *Debugger) Open(p string) (http.File, error) {
	// if d.c == nil {
	// 	result := ""
	// 	if strings.HasPrefix(p, `/`) {
	// 		p = p[1:]
	// 	}
	// 	d.i.Walk(func(target, relative string, info os.FileInfo) error {
	// 		if relative == p {
	// 			result = target
	// 		}
	// 		return nil
	// 	})
	// 	return os.Open(result)
	// }
	return os.Open(filepath.Join(d.d, p))
}

// NewDebugger watches entry directories and files and copies them to
// a temporary directory. If the Compiler is provided, the files will
// also be refined by the compiler while being copied. If any of the watched
// files are changed, the changes are reflected within 1-2 seconds.
// The returned Debugger is meant to replace a zassets.Store object
// for live-editing.
func NewDebugger(entries, ignore []string, c *Compiler) *Debugger {
	var err error
	panicOnError := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	d := new(Debugger)
	if c == nil {
		c, err = NewCompiler(WithDebug()) // pass through compiler
		panicOnError(err)
	}
	d.c = c
	d.SetLogger(nil)
	_, gofile, _, ok := runtime.Caller(1)
	if !ok {
		panic(errors.New(`cannot determine asset origin file`))
	}
	fromDirectory := filepath.Dir(gofile)
	adjusted := make([]string, 0)
	for _, e := range entries {
		if filepath.IsAbs(e) {
			adjusted = append(adjusted, e)
		} else {
			adjusted = append(adjusted, filepath.Join(fromDirectory, e))
		}
	}
	d.i, err = NewIterator(adjusted, ignore)
	panicOnError(err)
	// spew.Dump(adjusted)

	d.w, err = fsnotify.NewWatcher()
	panicOnError(err)
	d.Watch(adjusted...)
	d.d, err = ioutil.TempDir(os.TempDir(), `zassets-debug-*`)
	panicOnError(err)
	panicOnError(d.c.Run(d.d, d.i)) // initial compilation
	go d.watch()
	d.l.Printf(`Compiled output served from <%s>.`, d.d)
	return d
}
