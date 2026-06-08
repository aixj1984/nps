package stream

import (
	"bytes"
	"image"
	"image/jpeg"
)

type bytesBuffer struct {
	bytes.Buffer
}

func encodeJPEG(buf *bytesBuffer, img image.Image, quality int) error {
	if quality < 30 {
		quality = 30
	}
	if quality > 95 {
		quality = 95
	}
	return jpeg.Encode(buf, img, &jpeg.Options{Quality: quality})
}
