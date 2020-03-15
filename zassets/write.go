package zassets

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path"
)

func write(w io.Writer, s http.FileSystem) error {
	a := zip.NewWriter(w)
	a.SetComment(`Resource pack generated by github.com/dkotik/zassets.`)
	defer a.Close()
	return Walk(s, `/`, func(p string, i os.FileInfo, err error) error {
		if i.IsDir() || err != nil {
			return err
		}
		h, err := zip.FileInfoHeader(i)
		if err != nil {
			return err
		}
		h.Name = path.Clean(p)
		h.Method = zip.Deflate
		w, err := a.CreateHeader(h)
		if err != nil {
			return err
		}
		r, err := s.Open(h.Name)
		if err != nil {
			return err
		}
		defer r.Close()
		_, err = io.Copy(w, r)
		return err
	})
}
