package zassets

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/OneOfOne/xxhash"
)

// Public provides access to assets distributed among different FileSystems
// by hash name. This is useful for serving static assets from a single
// http.Handler sitting behing a caching layer or from a CDN.
// To find out the hash name for a given asset, call the PublicName function.
var Public http.FileSystem = &publicSystem{}
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
func PublicName(namespace, path string) (string, bool) {
	publicMutex.Lock()
	defer publicMutex.Unlock()
	result, ok := assetLookUpMap[namespace+`:`+path]
	return result, ok
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
			line := scanner.Text()
			if i := strings.Index(line, ` `); i > 5 {
				add(line[i+1:]+path.Ext(path.Base(line)), line[:i])
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
