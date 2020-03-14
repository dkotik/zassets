package goresminpack

// Public creates two look-up maps of content-based hashes
// for <resourceDir>/public.
// The maps can be used to serve static assets from one or more
// resource packs

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/OneOfOne/xxhash"
)

var PublicHashSeed string = `goresminpack`
var assetMap = make(map[string]*publicAsset)
var assetLookUpMap = make(map[string]string)

// ContentType returns content type based on file extension.
func ContentType(file, defaultValue string) string {
	switch ext := strings.ToLower(strings.TrimPrefix(path.Ext(file), `.`)); ext {
	case `html`:
		return `text/html; charset=utf-8`
	case `css`:
		return `text/css; charset=utf-8`
	case `jpg`, `jpeg`:
		return `image/jpeg`
	case `png`:
		return `image/png`
	case `gif`:
		return `image/gif`
	case `svg`:
		return `image/svg+xml`
	case `ico`:
		// return `image/vnd.microsoft.icon`
		return `image/x-icon`
	case `js`:
		// used to be `application/javascript` This is in accordance with an IETF draft that treats application/javascript as obsolete.
		return `text/javascript`
	case `txt`:
		return `text/plain; charset=utf-8`
	case `mp3`:
		return `audio/mpeg`
	case `pdf`, `zip`, `xml`:
		return fmt.Sprintf("application/%s", ext)
	}
	return defaultValue
}

type publicAsset struct {
	Name, ContentType string
	Storage           http.FileSystem
}

func PublicName(p string) (string, bool) {
	result, ok := assetLookUpMap[p]
	return result, ok
}

func PublicRegister(p string, d http.FileSystem) error {
	f, err := d.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()

	h := xxhash.NewS32(0)
	_, err = io.Copy(h, f)
	if err != nil {
		return err
	}
	h.WriteString(PublicHashSeed)
	tag := fmt.Sprintf(`%x%s`, h.Sum32(), strings.ToLower(path.Ext(p)))

	assetLookUpMap[p] = tag
	assetMap[tag] = &publicAsset{p, ContentType(p, `text/plain; charset=utf-8`), d}
	return nil
}

// PublicHTTPHandler serves public assets by their hash names.
func PublicHTTPHandler(w http.ResponseWriter, r *http.Request) error {
	asset, ok := assetMap[path.Base(r.URL.Path)]
	if !ok {
		http.Error(w, "file not found", http.StatusNotFound)
		return nil
	}
	f, err := asset.Storage.Open(asset.Name)
	if !ok {
		http.Error(w, "file not found", http.StatusNotFound)
		return nil
	}
	defer f.Close()
	w.Header().Add("content-type", asset.ContentType)
	// TODO: add eternal expiration headers
	_, err = io.Copy(w, f)
	return err
}
