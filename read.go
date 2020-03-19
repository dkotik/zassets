package zassets

import (
	"bytes"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

// Must panics if there is an error associated with loading assets.
func Must(s *Store, err error) *Store {
	if err != nil {
		panic(err)
	}
	return s
}

// FromBytes serves assets from bytes encoding a zip archive.
// Used for accessing assets embedded into a Go binary.
func FromBytes(b []byte) (*Store, error) {
	return &Store{bytes.NewReader(b), int64(len(b))}, nil
}

// FromArchive serves assets from a zip archive.
func FromArchive(p string) (*Store, error) {
	r, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	s, err := r.Stat()
	if err != nil {
		return nil, err
	}
	return &Store{r, s.Size()}, nil
}

// // FromDirectory serves assets from disk.
// func FromDirectory(p string) (*Store, error) {
// 	return &Store{http.Dir(p)}, nil
// }

// Walk emulates filepath.Walk()
func Walk(s http.FileSystem, p string, f filepath.WalkFunc) error {
	root, err := s.Open(p)
	if err != nil {
		return err
	}
	files, err := root.Readdir(-1)
	root.Close()
	if err != nil {
		return err
	}

	for _, info := range files {
		err = f(path.Join(p, info.Name()), info, err)
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = Walk(s, path.Join(p, info.Name()), f)
		}
		if err != nil {
			return err
		}
	}
	return err
}
