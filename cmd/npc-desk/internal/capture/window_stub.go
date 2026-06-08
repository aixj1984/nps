//go:build !windows

package capture

import "errors"

func ListWindows() ([]Source, error) {
	return nil, errors.New("application window capture is only supported on Windows")
}

func captureWindow(_ string) (image.Image, error) {
	return nil, errors.New("application window capture is only supported on Windows")
}
