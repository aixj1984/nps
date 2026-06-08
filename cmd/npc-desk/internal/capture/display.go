package capture

import (
	"image"
	"strconv"

	"github.com/kbinani/screenshot"
)

func ListScreens() ([]Source, error) {
	n := screenshot.NumActiveDisplays()
	out := make([]Source, 0, n)
	for i := 0; i < n; i++ {
		b := screenshot.GetDisplayBounds(i)
		name := "Display " + strconv.Itoa(i+1)
		src := Source{
			ID:   strconv.Itoa(i),
			Name: name + " (" + strconv.Itoa(b.Dx()) + "×" + strconv.Itoa(b.Dy()) + ")",
			Type: "screen",
		}
		if img, err := screenshot.CaptureRect(b); err == nil {
			src.ThumbPNG = thumb(img, 160)
		}
		out = append(out, src)
	}
	return out, nil
}

func captureDisplay(id string) (image.Image, error) {
	i, err := strconv.Atoi(id)
	n := screenshot.NumActiveDisplays()
	if err != nil || i < 0 || i >= n {
		i = 0
	}
	return screenshot.CaptureRect(screenshot.GetDisplayBounds(i))
}
