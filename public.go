package zassets

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/OneOfOne/xxhash"
)

// Handler serves files only from Public.
// URL.Path does not have to be stripped.
var Handler http.Handler = new(handler)

// Public provides access to assets distributed among different FileSystems
// by hash name. This is useful for serving static assets from a single
// http.Handler sitting behing a caching layer or from a CDN.
// To find out the hash name for a given asset, call the PublicName function.
var Public http.FileSystem = &publicSystem{}

// PublicPrefix is added to all the outputs of PublicName.
var PublicPrefix = `/asset/`

// There must be two spaces or a space and an asterisk between each sum value and filename to be compared (the second space indicates text mode, the asterisk binary mode).
var reParseSum = regexp.MustCompile(`^([^\s]{5,}) [ \*](.*)$`)
var assetMap = make(map[string]*publicAsset)
var assetLookUpMap = make(map[string]string)
var publicMutex = &sync.Mutex{} // protects from async asset map errors

type publicAsset struct {
	Name    string
	Storage http.FileSystem
}

type publicSystem struct{}

func (s *publicSystem) Open(p string) (http.File, error) {
	publicMutex.Lock()
	defer publicMutex.Unlock()
	a, ok := assetMap[p]
	if !ok {
		return nil, &os.PathError{Op: `open`, Path: p, Err: os.ErrNotExist}
	}
	return a.Storage.Open(a.Name)
}

// PublicName returns the hash name associated with an asset namespace and path.
// This function is useful for pointing asset URLs in any templating engine
// contained in the Public virtual file system.
func PublicName(namespace, path string) (string, error) {
	publicMutex.Lock()
	defer publicMutex.Unlock()
	result, ok := assetLookUpMap[namespace+`:`+path]
	if !ok {
		// TODO: make this error standard ErrNoAsset?
		return "", &os.PathError{Op: `open`, Path: fmt.Sprintf(`%s:%s`, namespace, path), Err: errors.New(`asset does not exist or was not registered`)}
	}
	return PublicPrefix + result, nil
}

// PublicRegister connects a FileSystem to Public, differentiated by name space.
// If the FileSystem includes "sum.xxhash", "sum.md5", or "sum.sha256" file,
// the hash map is constructed by using values from this file. Otherwise,
// the hash map is constructed by hashing name space and asset path, which is
// useful for debugging Public.
func PublicRegister(namespace string, d http.FileSystem) (err error) {
	publicMutex.Lock()
	defer publicMutex.Unlock()
	add := func(p, h string) {
		// log.Println(`public`, namespace+`:`+p, h)
		assetLookUpMap[namespace+`:`+p] = h
		assetMap[h] = &publicAsset{p, d}
	}
	var sumFile http.File
	sumFile, err = d.Open(`sum.xxh64`)
	if err != nil {
		sumFile, err = d.Open(`sum.md5`)
		if err != nil {
			sumFile, err = d.Open(`sum.sha256`)
		}
	}

	if err == nil { // sumFile located
		scanner := bufio.NewScanner(sumFile)
		for scanner.Scan() {
			// If line parses <hash> *<path>, add it with proper extension.
			if m := reParseSum.FindStringSubmatch(scanner.Text()); m[0] != `` {
				add(m[2], m[1]+path.Ext(m[2]))
			}
		}
		return scanner.Err()
	}

	h := xxhash.New64()
	err = Walk(d, ``, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		h.Reset()
		h.WriteString(namespace)
		h.WriteString(p)
		add(p, fmt.Sprintf(`%x%s`, h.Sum([]byte{}), path.Ext(info.Name())))
		return nil
	})
	return err
}

type handler struct{}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := path.Base(r.URL.Path)
	f, err := Public.Open(name)
	if err != nil {
		http.Error(w, fmt.Errorf("could not locate asset: %w", err).Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	io.Copy(w, f)
	f.Close()
}
