package goresminpack

// Public creates two look-up maps of content-based hashes
// for <resourceDir>/public.
// The maps can be used to serve static assets from one or more
// resource packs

import (
	"fmt"
	"io"
	"net/http"

	"github.com/OneOfOne/xxhash"
)

var PublicHashSeed string = `goresminpack`

func contentHash(r io.Reader) (string, error) {
	h := xxhash.NewS32(0)
	_, err := io.Copy(h, r)
	if err != nil {
		return "", err
	}
	h.WriteString(PublicHashSeed)
	return fmt.Sprintf(`%x`, h.Sum32()), nil
}

func PublicName(p string) (string, error) {
	return "", nil
}

// PublicHTTPHandler serves public assets by their hash names.
func PublicHTTPHandler(w http.ResponseWriter, r *http.Request) error {

	return nil
}
