package compile

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch sources and re-compile them to serve through the asset object.
// github.com/fsnotify/fsnotify

type debugger struct {
	d string
	c *Compiler
	i *Iterator
	// v *EmbedValues
	w *fsnotify.Watcher
	l *log.Logger
}

func (d *debugger) watch() {
	t := time.Tick(time.Second)
	var somethingChanged bool
	for {
		select {
		case <-t: // time to update
			if somethingChanged {
				somethingChanged = false
				err := d.c.Run(d.d, d.i)
				if err != nil {
					d.l.Println("refining error:", err.Error())
				}
			}
		case event, ok := <-d.w.Events:
			if !ok {
				return
			}
			// log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				d.l.Println("modified file:", event.Name)
				somethingChanged = true
			}
		case err, ok := <-d.w.Errors:
			if !ok {
				return
			}
			d.l.Println("error:", err)
		}
	}
}

func (d *debugger) Watch(p string) error {
	i, _ := NewIterator([]string{p}, []string{})
	return i.Walk(func(target, relative string, info os.FileInfo) error {
		return d.w.Add(target)
	})
}

func Debug(entries, ignore []string, refine bool) (err error) {
	d := new(debugger)
	d.w, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	_, gofile, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New(`cannot determine asset origin file`)
	}
	for _, e := range entries {
		if filepath.IsAbs(e) {
			d.Watch(e)
		} else {
			d.Watch(filepath.Join(filepath.Dir(gofile), e))
		}
	}

	d.d, err = ioutil.TempDir(os.TempDir(), `zassets-debug-*`)
	if err != nil {
		return err
	}
	d.i, err = NewIterator(entries, ignore)
	if err != nil {
		return err
	}
	d.c, err = NewCompiler(
		WithDefaultOptions(), WithDebug())
	if err != nil {
		return err
	}
	go d.watch()
	return err
}
