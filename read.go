package zassets

import (
	"bytes"
	"fmt"
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

// MustOpen panics if there is an error opening an asset.
// Use it as a convenience method when loading required assets.
// Example: asset := zassets.MustOpen(Assets.Open(`schema.sql`))
func MustOpen(f http.File, err error) http.File {
	if err != nil {
		panic(fmt.Errorf(`require resource failed to open: %w`, err))
	}
	return f
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
