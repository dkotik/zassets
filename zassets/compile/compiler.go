package compile

import (
	"io"
	"os"
)

// Refiner converts one type of asset into another or optimizes the content.
type Refiner interface {
	Match(path string) (ok bool)
	Rename(oldPath string) (newPath string)
	Refine(destination, source string) error
	Debug(destination, source string) error
}

type passthrough struct{}

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
