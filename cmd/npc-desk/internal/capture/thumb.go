package capture

import (
	"bytes"
	"image"
	"image/jpeg"

	"github.com/nfnt/resize"
)

func thumb(img image.Image, width uint) []byte {
	if img.Bounds().Dx() > int(width) {
		img = resize.Resize(width, 0, img, resize.Lanczos3)
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 60})
	return buf.Bytes()
}
