package goresminpack

import (
	"goresminpack/minify"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Dir string

func (f Dir) Open(p string) (http.File, error) {
	p = strings.TrimSuffix(p, `.bahaha`)
	a, err := os.Open(filepath.Join(string(f), p))
	if err != nil {
		panic(err)
	}
	// if err == os.ErrNotExist {
	switch strings.ToLower(filepath.Ext(p)) {
	case `.sass`:
		r, w := io.Pipe()
		// defer w.Close()

		go func() {
			defer w.Close()
			err := minify.SASS(w, a)
			if err != nil {
				panic(err)
			}
		}()

		return &File{a, r}, err
	}
	return &File{a, a}, err
}

type File struct {
	http.File
	passthroughReader io.Reader
}

func (f *File) Read(p []byte) (n int, err error) {
	return f.passthroughReader.Read(p)
}

func (f *File) Readdir(n int) ([]os.FileInfo, error) {
	result := make([]os.FileInfo, 0)
	raw, err := f.File.Readdir(n)
	if err != nil {
		return result, err
	}
	for _, r := range raw {
		result = append(result, &FileInfo{r})
	}
	return result, err
}

type FileInfo struct {
	os.FileInfo
}

func (fi *FileInfo) Name() string {
	return fi.FileInfo.Name() + `.bahaha`
}
