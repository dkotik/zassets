package compile

import (
	"image"
	"os"
	"regexp"

	"github.com/nfnt/resize"

	"io"
)

var _ Refiner = &RefineRaster{}

// RefineRaster provides basic image compression.
type RefineRaster struct {
	passthrough

	MatchPath *regexp.Regexp
	Encode    func(io.Writer, image.Image) error
	Decode    func(io.Reader) (image.Image, error)
}

func (rf *RefineRaster) Match(p string) bool {
	if reMinPass.MatchString(p) {
		return false // skip already minified assets
	}
	return rf.MatchPath.MatchString(p)
}

func (rf *RefineRaster) Refine(destination, source string) error {
	r, err := os.Open(source)
	if err != nil {
		return err
	}
	defer r.Close()
	img, err := rf.Decode(r)
	if err != nil {
		return err
	}
	w, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer w.Close()
	return rf.Encode(w, resize.Thumbnail(1920, 1080, img, resize.Lanczos3))
}
