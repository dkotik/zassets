package zassets

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

// TODO: https://github.com/gnue/httpfs/blob/master/zipfs/dir.go uses caching and might be more efficient

// NewStore indexes a zipped archive.
func NewStore(r io.ReaderAt, size int64) (*Store, error) {
	z, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	s := &Store{z, make(map[string]*zip.File)}
	for _, f := range z.File {
		s.m[f.Name] = f
	}
	return s, nil
}

// Store presents as zipped archive as http.FileSystem.
type Store struct {
	r *zip.Reader
	m map[string]*zip.File
}

// String reduces the asset and returns it as a string.
func (s *Store) String(p string) (string, error) {
	r, err := s.Open(p)
	if err != nil {
		return "", err
	}
	defer r.Close()
	var b bytes.Buffer
	_, err = io.Copy(&b, r)
	return b.String(), err
}

// Bytes reduces the asset and returns it as a byte array.
func (s *Store) Bytes(p string) ([]byte, error) {
	r, err := s.Open(p)
	if err != nil {
		return []byte{}, err
	}
	defer r.Close()
	var b bytes.Buffer
	_, err = io.Copy(&b, r)
	return b.Bytes(), err
}

// Open returns a handle to the underlying file.
func (s *Store) Open(p string) (http.File, error) {
	f, ok := s.m[p]
	if ok {
		handle, err := f.Open()
		if err != nil {
			return nil, err
		}
		return &zipFile{f, handle, 0}, nil
	}
	// Contruct a directory interface.
	dir := s.makeDir(p)
	if len(dir.i) == 0 { //
		return nil, &os.PathError{Op: `open`, Path: p, Err: os.ErrNotExist}
	}
	return dir, nil
}

func (s *Store) makeDir(p string) *zipDir {
	dir := &zipDir{i: make([]os.FileInfo, 0), n: path.Base(p)}
	prefix := p
	if prefix == "" {
		// root folder
	} else if prefix == "/" || prefix == `\` || prefix == "." {
		prefix = "" // also root folder, but by another expression
	} else if !strings.HasSuffix(prefix, `/`) {
		prefix = prefix + `/`
	}
	for k, f := range s.m {
		if strings.HasPrefix(k, prefix) {
			remainder := strings.TrimPrefix(k, prefix)
			if index := strings.Index(remainder, `/`); index > 0 {
				subDir := s.makeDir(path.Join(p, remainder[:index+1]))
				if len(subDir.i) > 0 {
					dir.i = append(dir.i, subDir)
				}
			} else {
				i := f.FileInfo()
				dir.i = append(dir.i, i)
				if t := i.ModTime(); t.After(dir.m) {
					dir.m = t
				}
			}
		}
	}
	// TODO: this is inefficient! I should cache dir file info here
	// and check cache first line of the function
	return dir
}

type zipDir struct {
	i []os.FileInfo
	n string
	m time.Time
	p int
}

func (d *zipDir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.p = 0
		return 0, nil
	}
	return 0, os.ErrInvalid
}

func (d *zipDir) Readdir(count int) ([]os.FileInfo, error) {
	if d.p >= len(d.i) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.i)-d.p {
		count = len(d.i) - d.p
	}
	e := d.i[d.p : d.p+count]
	d.p += count
	return e, nil
}

func (d *zipDir) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.n)
}
func (d *zipDir) Close() error               { return nil }
func (d *zipDir) Stat() (os.FileInfo, error) { return d, nil }
func (d *zipDir) Name() string               { return d.n }
func (d *zipDir) Size() int64                { return 0 }
func (d *zipDir) Mode() os.FileMode          { return 0755 | os.ModeDir }
func (d *zipDir) ModTime() time.Time         { return d.m }
func (d *zipDir) IsDir() bool                { return true }
func (d *zipDir) Sys() interface{}           { return nil }

type zipFile struct {
	*zip.File
	f    io.ReadCloser
	seek int64
}

func (z *zipFile) Open() (err error) {
	if z.f != nil {
		return z.f.Close()
	}
	z.f, err = z.File.Open()
	if err != nil {
		return err
	}
	if z.seek > 0 {
		_, err = io.CopyN(ioutil.Discard, z.f, z.seek)
	}
	return err
}

func (z *zipFile) Read(p []byte) (n int, err error) {
	n, err = z.f.Read(p)
	z.seek += int64(n)
	return n, err
}

func (z *zipFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", z.Name)
}

func (z *zipFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		z.seek = 0 + offset
	case io.SeekCurrent:
		z.seek += offset
	case io.SeekEnd:
		z.seek = z.FileInfo().Size() + offset
	default: // should never happen
		panic(os.ErrInvalid)
	}
	return z.seek, z.Open() // open also rewinds to seek position
}

func (z *zipFile) Stat() (os.FileInfo, error) { return z.FileInfo(), nil }

func (z *zipFile) Close() error {
	if z.f == nil {
		return os.ErrClosed
	}
	z.f.Close()
	z.f = nil
	z.seek = 0
	return nil
}
