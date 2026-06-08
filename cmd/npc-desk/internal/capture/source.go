package capture

import (
	"image"
	"log"
)

// Source is a capturable screen or application window.
type Source struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"` // screen | app
	ThumbPNG []byte `json:"thumb,omitempty"`
}

// Capture grabs a frame for the given source.
func Capture(src Source) (image.Image, error) {
	return CaptureWithMode(src, src.Type)
}

// CaptureWithMode uses explicit share mode (screen/app); infers HWND when mode is missing.
func CaptureWithMode(src Source, mode string) (image.Image, error) {
	if UseWindowCapture(src, mode) {
		log.Printf("capture: window hwnd=%s mode=%q src.type=%q name=%q", src.ID, mode, src.Type, src.Name)
		return captureWindow(src.ID)
	}
	log.Printf("capture: display index=%s mode=%q src.type=%q name=%q", src.ID, mode, src.Type, src.Name)
	return captureDisplay(src.ID)
}
