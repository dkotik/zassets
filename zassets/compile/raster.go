package compile

import (
	"image"
	"os"
	"regexp"

	"github.com/nfnt/resize"

	"image/jpeg"
	"image/png"
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
	return rf.Encode(w, resize.Thumbnail(800, 600, img, resize.Lanczos3))
}

// ResizeJPG reduces a JPG image to web-dimentions.
func ResizeJPG(w io.Writer, r io.Reader) error {
	img, err := jpeg.Decode(r)
	if err != nil {
		return err
	}
	return jpeg.Encode(w, resize.Thumbnail(800, 600, img, resize.Lanczos3), nil)
}

// ResizePNG reduces a PNG image to web-dimentions.
func ResizePNG(w io.Writer, r io.Reader) error {
	img, err := png.Decode(r)
	if err != nil {
		return err
	}
	return png.Encode(w, resize.Thumbnail(800, 600, img, resize.Lanczos3))
}
