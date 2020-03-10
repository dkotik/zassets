package minify

import (
	"github.com/nfnt/resize"

	"image/jpeg"
	"image/png"
	"io"
)

// Raster provide basic image compression.

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
