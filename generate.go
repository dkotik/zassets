package goresminpack

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	pathpkg "path"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/shurcooL/httpfs/vfsutil"
	"github.com/shurcooL/vfsgen"
)

const hotSwapCode = `
package {{ .package }}
+build debug,dev

// Locate current directory? Do I even need this? Yes
_, gofile, _, ok := runtime.Caller(1)
if ok {
	directory, _ = filepath.Abs(filepath.Join(filepath.Dir(gofile), filepath.FromSlash(directory)))
}
`

// Generate packs the resources into a Go file.
func Generate(assetDirectory, outputFilePath, outputObjectName string) (err error) {
	opt := vfsgen.Options{
		Filename:        outputFilePath,
		PackageName:     strings.ToLower(path.Base(path.Dir(outputFilePath))),
		VariableName:    outputObjectName,
		VariableComment: fmt.Sprintf("%s statically implements the virtual filesystem provided to vfsgen.", outputObjectName),
		BuildTags:       "!debug,!dev",
	}

	// Use an in-memory buffer to generate the entire output.
	buf := new(bytes.Buffer)

	err = t.ExecuteTemplate(buf, "Header", opt)
	if err != nil {
		return err
	}

	var toc toc
	err = findAndWriteFiles(buf, http.Dir(assetDirectory), &toc)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(buf, "DirEntries", toc.dirs)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(buf, "Trailer", toc)
	if err != nil {
		return err
	}

	// Write output file (all at once).
	fmt.Println("writing", opt.Filename)
	err = ioutil.WriteFile(opt.Filename, buf.Bytes(), 0644)
	return err
}

type toc struct {
	dirs []*dirInfo

	HasCompressedFile bool // There's at least one compressedFile.
	HasFile           bool // There's at least one uncompressed file.
}

// fileInfo is a definition of a file.
type fileInfo struct {
	Path             string
	Name             string
	ModTime          time.Time
	UncompressedSize int64
}

// dirInfo is a definition of a directory.
type dirInfo struct {
	Path    string
	Name    string
	ModTime time.Time
	Entries []string
}

// findAndWriteFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findAndWriteFiles(buf *bytes.Buffer, fs http.FileSystem, toc *toc) error {
	walkFn := func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			// Consider all errors reading the input filesystem as fatal.
			return err
		}

		switch fi.IsDir() {
		case false:
			file := &fileInfo{
				Path:             path,
				Name:             pathpkg.Base(path),
				ModTime:          fi.ModTime().UTC(),
				UncompressedSize: fi.Size(),
			}

			marker := buf.Len()

			// Write CompressedFileInfo.
			err = writeCompressedFileInfo(buf, file, r)
			switch err {
			default:
				return err
			case nil:
				toc.HasCompressedFile = true
			// If compressed file is not smaller than original, revert and write original file.
			case errCompressedNotSmaller:
				_, err = r.Seek(0, io.SeekStart)
				if err != nil {
					return err
				}

				buf.Truncate(marker)

				// Write FileInfo.
				err = writeFileInfo(buf, file, r)
				if err != nil {
					return err
				}
				toc.HasFile = true
			}
		case true:
			entries, err := readDirPaths(fs, path)
			if err != nil {
				return err
			}

			dir := &dirInfo{
				Path:    path,
				Name:    pathpkg.Base(path),
				ModTime: fi.ModTime().UTC(),
				Entries: entries,
			}

			toc.dirs = append(toc.dirs, dir)

			// Write DirInfo.
			err = t.ExecuteTemplate(buf, "DirInfo", dir)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := vfsutil.WalkFiles(fs, "/", walkFn)
	return err
}

// readDirPaths reads the directory named by dirname and returns
// a sorted list of directory paths.
func readDirPaths(fs http.FileSystem, dirname string) ([]string, error) {
	fis, err := vfsutil.ReadDir(fs, dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(fis))
	for i := range fis {
		paths[i] = pathpkg.Join(dirname, fis[i].Name())
	}
	sort.Strings(paths)
	return paths, nil
}

// writeCompressedFileInfo writes CompressedFileInfo.
// It returns errCompressedNotSmaller if compressed file is not smaller than original.
func writeCompressedFileInfo(w io.Writer, file *fileInfo, r io.Reader) error {
	err := t.ExecuteTemplate(w, "CompressedFileInfo-Before", file)
	if err != nil {
		return err
	}
	sw := &stringWriter{Writer: w}
	gw := gzip.NewWriter(sw)
	_, err = io.Copy(gw, r)
	if err != nil {
		return err
	}
	err = gw.Close()
	if err != nil {
		return err
	}
	if sw.N >= file.UncompressedSize {
		return errCompressedNotSmaller
	}
	err = t.ExecuteTemplate(w, "CompressedFileInfo-After", file)
	return err
}

var errCompressedNotSmaller = errors.New("compressed file is not smaller than original")

// Write FileInfo.
func writeFileInfo(w io.Writer, file *fileInfo, r io.Reader) error {
	err := t.ExecuteTemplate(w, "FileInfo-Before", file)
	if err != nil {
		return err
	}
	sw := &stringWriter{Writer: w}
	_, err = io.Copy(sw, r)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "FileInfo-After", file)
	return err
}

var t = template.Must(template.New("").Funcs(template.FuncMap{
	"quote": strconv.Quote,
}).Parse(`{{define "Header"}}// Code generated by vfsgen; DO NOT EDIT.

{{with .BuildTags}}// +build {{.}}

{{end}}package {{.PackageName}}

import (
	"net/http"
	"os"
	"time"

	vfsgen "goresminpack/embed"
)

var {{.VariableName}} = func() http.FileSystem {
	fs := vfsgen.FS{
{{end}}



{{define "CompressedFileInfo-Before"}}		{{quote .Path}}: &vfsgen.CompressedFileInfo{
			{{quote .Name}},
			{{template "Time" .ModTime}},
			{{.UncompressedSize}},
{{/* This blank line separating compressedContent is neccessary to prevent potential gofmt issues. See issue #19. */}}
			[]byte("{{end}}{{define "CompressedFileInfo-After"}}"),
		},
{{end}}



{{define "FileInfo-Before"}}		{{quote .Path}}: &vfsgen.FileInfo{
			{{quote .Name}},
			{{template "Time" .ModTime}},
			[]byte("{{end}}{{define "FileInfo-After"}}"),
		},
{{end}}



{{define "DirInfo"}}		{{quote .Path}}: &vfsgen.DirInfo{
			{{quote .Name}},
			{{template "Time" .ModTime}},
		},
{{end}}



{{define "DirEntries"}}	}
{{range .}}{{if .Entries}}	fs[{{quote .Path}}].(*vfsgen.DirInfo).entries = []os.FileInfo{{"{"}}{{range .Entries}}
		fs[{{quote .}}].(os.FileInfo),{{end}}
	}
{{end}}{{end}}
	return fs
}()
{{end}}

{{define "Time"}}
{{- if .IsZero -}}
	time.Time{}
{{- else -}}
	time.Date({{.Year}}, {{printf "%d" .Month}}, {{.Day}}, {{.Hour}}, {{.Minute}}, {{.Second}}, {{.Nanosecond}}, time.UTC)
{{- end -}}
{{end}}

{{define "Trailer"}}

{{end}}
`))
