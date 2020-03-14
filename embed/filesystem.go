package embed

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"

	"time"
)

// File system inspired by vfsgen

type FS map[string]interface{}

func (fs FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {
	case *CompressedFileInfo:
		gr, err := gzip.NewReader(bytes.NewReader(f.compressedContent))
		if err != nil {
			// This should never happen because we generate the gzip bytes such that they are always valid.
			panic("unexpected error reading own gzip compressed bytes: " + err.Error())
		}
		return &CompressedFile{
			CompressedFileInfo: f,
			gr:                 gr,
		}, nil
	case *FileInfo:
		return &File{
			FileInfo: f,
			Reader:   bytes.NewReader(f.content),
		}, nil
	case *DirInfo:
		return &Dir{
			DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// CompressedFileInfo is a static definition of a gzip compressed file.
type CompressedFileInfo struct {
	name              string
	modTime           time.Time
	compressedContent []byte
	uncompressedSize  int64
}

func (f *CompressedFileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *CompressedFileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *CompressedFileInfo) GzipBytes() []byte {
	return f.compressedContent
}

func (f *CompressedFileInfo) Name() string       { return f.name }
func (f *CompressedFileInfo) Size() int64        { return f.uncompressedSize }
func (f *CompressedFileInfo) Mode() os.FileMode  { return 0444 }
func (f *CompressedFileInfo) ModTime() time.Time { return f.modTime }
func (f *CompressedFileInfo) IsDir() bool        { return false }
func (f *CompressedFileInfo) Sys() interface{}   { return nil }

// CompressedFile is an opened compressedFile instance.
type CompressedFile struct {
	*CompressedFileInfo
	gr      *gzip.Reader
	grPos   int64 // Actual gr uncompressed position.
	seekPos int64 // Seek uncompressed position.
}

func (f *CompressedFile) Read(p []byte) (n int, err error) {
	if f.grPos > f.seekPos {
		// Rewind to beginning.
		err = f.gr.Reset(bytes.NewReader(f.compressedContent))
		if err != nil {
			return 0, err
		}
		f.grPos = 0
	}
	if f.grPos < f.seekPos {
		// Fast-forward.
		_, err = io.CopyN(ioutil.Discard, f.gr, f.seekPos-f.grPos)
		if err != nil {
			return 0, err
		}
		f.grPos = f.seekPos
	}
	n, err = f.gr.Read(p)
	f.grPos += int64(n)
	f.seekPos = f.grPos
	return n, err
}
func (f *CompressedFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		f.seekPos = 0 + offset
	case io.SeekCurrent:
		f.seekPos += offset
	case io.SeekEnd:
		f.seekPos = f.uncompressedSize + offset
	default:
		panic(fmt.Errorf("invalid whence value: %v", whence))
	}
	return f.seekPos, nil
}
func (f *CompressedFile) Close() error {
	return f.gr.Close()
}

// FileInfo is a static definition of an uncompressed file (because it's not worth gzip compressing).
type FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *FileInfo) NotWorthGzipCompressing() {}

func (f *FileInfo) Name() string       { return f.name }
func (f *FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *FileInfo) Mode() os.FileMode  { return 0444 }
func (f *FileInfo) ModTime() time.Time { return f.modTime }
func (f *FileInfo) IsDir() bool        { return false }
func (f *FileInfo) Sys() interface{}   { return nil }

// File is an opened file instance.
type File struct {
	*FileInfo
	*bytes.Reader
}

func (f *File) Close() error {
	return nil
}

// DirInfo is a static definition of a directory.
type DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *DirInfo) Close() error               { return nil }
func (d *DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *DirInfo) Name() string       { return d.name }
func (d *DirInfo) Size() int64        { return 0 }
func (d *DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *DirInfo) ModTime() time.Time { return d.modTime }
func (d *DirInfo) IsDir() bool        { return true }
func (d *DirInfo) Sys() interface{}   { return nil }

// Dir is an opened dir instance.
type Dir struct {
	*DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
