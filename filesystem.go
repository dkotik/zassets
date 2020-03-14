package goresminpack

// Doodled this trying to emulate http.FileSystem - was not worth it.
// Decided to use a temporary directory instead.

//
// import (
// 	"goresminpack/minify"
// 	"io"
// 	"net/http"
// 	"os"
// 	"path"
// 	"path/filepath"
// 	"regexp"
// 	"strings"
//
// 	"github.com/tdewolff/minify/css"
// 	"github.com/tdewolff/minify/html"
// 	"github.com/tdewolff/minify/svg"
// )
//
// // Dir presents processed assets as http.FileSystem interface.
// type Dir struct {
// 	Path    string
// 	Debug   bool
// 	Include []string
// 	Ignore  []*regexp.Regexp
//
// 	// sources points renamed files to their source files
// 	// useful for .sass => .css transformations and the like
// 	sources map[string]string
// }
//
// // Open satisfies http.FileSystem interface.
// func (f *Dir) Open(p string) (http.File, error) {
// 	if renamed, ok := f.sources[p]; ok {
// 		p = renamed
// 	}
// 	a, err := os.Open(filepath.Join(f.Path, p))
// 	if err != nil {
// 		panic(err)
// 	}
// 	// if err == os.ErrNotExist {
// 	switch strings.ToLower(filepath.Ext(p)) {
// 	case `.js`:
// 		r, err := CompiledJS(filepath.Join(f.Path, filepath.FromSlash(p)), f.Debug)
// 		// log.Println("reading javascript", p)
// 		// io.Copy(os.Stdout, r)
// 		return &wrappedFile{a, r}, err
// 	case `.sass`:
// 		r, w := io.Pipe()
// 		// defer w.Close()
// 		go func() {
// 			defer w.Close()
// 			err := minify.SASS(w, a,
// 				// Look for additional files from the same directory and in includes.
// 				append([]string{filepath.Join(f.Path, filepath.FromSlash(path.Dir(p)))}, f.Include...)...)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		return &wrappedFile{a, r}, err
// 	case `.css`:
// 		r, w := io.Pipe()
// 		// defer w.Close()
// 		go func() {
// 			defer w.Close()
// 			err := css.Minify(nil, w, a, nil)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		return &wrappedFile{a, r}, err
// 	case `.html`:
// 		r, w := io.Pipe()
// 		// defer w.Close()
// 		go func() {
// 			defer w.Close()
// 			err := html.Minify(nil, w, a, nil)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		return &wrappedFile{a, r}, err
// 	case `.svg`:
// 		r, w := io.Pipe()
// 		// defer w.Close()
// 		go func() {
// 			defer w.Close()
// 			err := svg.Minify(nil, w, a, nil)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		return &wrappedFile{a, r}, err
// 	case `.tmpl`, `.sql`:
// 		r, w := io.Pipe()
// 		// defer w.Close()
// 		go func() {
// 			defer w.Close()
// 			err := minify.EatLineWhiteSpace(w, a)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		return &wrappedFile{a, r}, err
// 	case "": // possibly a directory
// 		files, err := a.Readdir(-1)
// 		if err == nil { // found a directory
// 			d := &wrappedFilteredDirectory{a, make([]os.FileInfo, 0)}
// 			for _, o := range files {
// 				checkPath := path.Join(p, o.Name())
// 				allow := func() bool {
// 					for _, r := range f.Ignore {
// 						if r.MatchString(checkPath) {
// 							return false
// 						}
// 					}
// 					return true
// 				}()
// 				if allow {
// 					switch ext := path.Ext(checkPath); strings.ToLower(ext) {
// 					case `.sass`, `.scss`: // File must be renamed.
// 						if f.sources == nil { // Initialize map if empty.
// 							f.sources = make(map[string]string)
// 						}
// 						newName := strings.TrimSuffix(o.Name(), ext) + `.css`
// 						f.sources[path.Join(p, newName)] = checkPath
// 						d.files = append(d.files, &renamed{o, newName})
// 					default:
// 						d.files = append(d.files, o)
// 					}
// 				}
// 			}
// 			return d, err
// 		}
// 	}
// 	return a, err
// }
//
// type wrappedFile struct {
// 	http.File
// 	passthroughReader io.ReadCloser
// }
//
// func (f *wrappedFile) Read(p []byte) (n int, err error) {
// 	return f.passthroughReader.Read(p)
// }
//
// func (f *wrappedFile) Close() error {
// 	err1 := f.File.Close()
// 	err2 := f.passthroughReader.Close()
// 	if err1 != nil {
// 		return err1
// 	}
// 	return err2
// }
//
// type wrappedFilteredDirectory struct {
// 	http.File
// 	files []os.FileInfo
// }
//
// func (w *wrappedFilteredDirectory) Readdir(n int) ([]os.FileInfo, error) {
// 	return w.files, nil
// }
//
// type renamed struct {
// 	os.FileInfo
// 	name string
// }
//
// func (fi *renamed) Name() string {
// 	return fi.name
// }
