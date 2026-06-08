package capture

import (
	"strconv"

	"github.com/kbinani/screenshot"
)

// IsDisplayIndexID reports whether id is a screen index (0..n-1), not a window HWND.
func IsDisplayIndexID(id string) bool {
	i, err := strconv.Atoi(id)
	if err != nil || i < 0 {
		return false
	}
	return i < screenshot.NumActiveDisplays()
}

// UseWindowCapture decides app-window vs full-display capture from mode and source id.
func UseWindowCapture(src Source, mode string) bool {
	if mode == "app" || src.Type == "app" {
		return true
	}
	if mode == "screen" || src.Type == "screen" {
		return false
	}
	// Fallback: large numeric ids are HWNDs; small ids are display indices.
	return !IsDisplayIndexID(src.ID)
}
