//go:build windows

package capture

import (
	"fmt"
	"image"
	"image/color"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/kbinani/screenshot"
)

var (
	user32                     = syscall.NewLazyDLL("user32.dll")
	gdi32                      = syscall.NewLazyDLL("gdi32.dll")
	procEnumWindows            = user32.NewProc("EnumWindows")
	procIsWindowVisible        = user32.NewProc("IsWindowVisible")
	procGetWindowTextW         = user32.NewProc("GetWindowTextW")
	procGetWindowRect          = user32.NewProc("GetWindowRect")
	procGetClientRect          = user32.NewProc("GetClientRect")
	procPrintWindow            = user32.NewProc("PrintWindow")
	procGetWindowDC            = user32.NewProc("GetWindowDC")
	procReleaseDC              = user32.NewProc("ReleaseDC")
	procCreateCompatibleDC     = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject           = gdi32.NewProc("SelectObject")
	procDeleteObject           = gdi32.NewProc("DeleteObject")
	procDeleteDC               = gdi32.NewProc("DeleteDC")
	procGetDIBits              = gdi32.NewProc("GetDIBits")
)

const (
	pwClientOnly          = 0x00000001
	pwRenderFullContent   = 0x00000002
	maxWindowScreenRatio  = 0.92 // skip listing near-fullscreen shells in app picker
)

type rect struct {
	Left, Top, Right, Bottom int32
}

type bitmapInfoHeader struct {
	Size, Width, Height    int32
	Planes, BitCount       uint16
	Compression, SizeImage uint32
	XPelsPerMeter          int32
	YPelsPerMeter          int32
	ClrUsed, ClrImportant  uint32
}

type bitmapInfo struct {
	Header bitmapInfoHeader
}

func ListWindows() ([]Source, error) {
	primary := screenshot.GetDisplayBounds(0)
	primaryArea := primary.Dx() * primary.Dy()

	var hwnds []syscall.Handle
	cb := syscall.NewCallback(func(hwnd syscall.Handle, _ uintptr) uintptr {
		if isVisible(hwnd) {
			hwnds = append(hwnds, hwnd)
		}
		return 1
	})
	_, _, _ = procEnumWindows.Call(cb, 0)

	out := make([]Source, 0, len(hwnds))
	for _, hwnd := range hwnds {
		title := windowTitle(hwnd)
		if title == "" {
			continue
		}
		r, ok := windowRect(hwnd)
		if !ok || r.width() < 80 || r.height() < 80 {
			continue
		}
		winArea := int(r.width()) * int(r.height())
		if primaryArea > 0 && float64(winArea)/float64(primaryArea) > maxWindowScreenRatio {
			// Skip desktop-sized surfaces (often the whole monitor, not a normal app window).
			continue
		}
		id := fmt.Sprintf("%d", uintptr(hwnd))
		src := Source{ID: id, Name: title, Type: "app"}
		if img, err := captureWindow(id); err == nil {
			src.ThumbPNG = thumb(img, 160)
		}
		out = append(out, src)
	}
	return out, nil
}

// captureWindow renders the HWND via PrintWindow (WM_PRINT), not screen BitBlt.
// This keeps sharing the window's own pixels when it is in the background or covered.
func captureWindow(id string) (image.Image, error) {
	hwnd := syscall.Handle(uintptr(parseUint(id)))
	if hwnd == 0 {
		return nil, fmt.Errorf("invalid window handle %q", id)
	}

	var lastErr error
	if wr, ok := windowRect(hwnd); ok && wr.width() > 0 && wr.height() > 0 {
		if img, err := captureHWNDPrint(hwnd, int(wr.width()), int(wr.height()), uint32(pwRenderFullContent)); err == nil {
			return img, nil
		} else {
			lastErr = err
		}
	}
	if cr, ok := clientRect(hwnd); ok && cr.width() > 0 && cr.height() > 0 {
		flags := uint32(pwRenderFullContent | pwClientOnly)
		if img, err := captureHWNDPrint(hwnd, int(cr.width()), int(cr.height()), flags); err == nil {
			return img, nil
		} else {
			lastErr = err
		}
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no bounds")
	}
	return nil, fmt.Errorf("window %s: %w", id, lastErr)
}

func captureHWNDPrint(hwnd syscall.Handle, width, height int, printFlags uint32) (image.Image, error) {
	hdcWnd, _, _ := procGetWindowDC.Call(uintptr(hwnd))
	if hdcWnd == 0 {
		return nil, fmt.Errorf("GetWindowDC failed")
	}
	defer procReleaseDC.Call(uintptr(hwnd), hdcWnd)

	hdcMem, _, _ := procCreateCompatibleDC.Call(hdcWnd)
	if hdcMem == 0 {
		return nil, fmt.Errorf("CreateCompatibleDC failed")
	}
	defer procDeleteDC.Call(hdcMem)

	hbmp, _, _ := procCreateCompatibleBitmap.Call(hdcWnd, uintptr(width), uintptr(height))
	if hbmp == 0 {
		return nil, fmt.Errorf("CreateCompatibleBitmap failed")
	}
	defer procDeleteObject.Call(hbmp)

	_, _, _ = procSelectObject.Call(hdcMem, hbmp)

	ok, _, _ := procPrintWindow.Call(uintptr(hwnd), hdcMem, uintptr(printFlags))
	if ok == 0 && printFlags&pwRenderFullContent != 0 {
		ok, _, _ = procPrintWindow.Call(uintptr(hwnd), hdcMem, uintptr(printFlags&^pwRenderFullContent))
	}
	if ok == 0 {
		return nil, fmt.Errorf("PrintWindow failed")
	}
	return dibToImage(hdcMem, hbmp, width, height)
}

func dibToImage(hdcMem, hbmp uintptr, width, height int) (image.Image, error) {
	var bi bitmapInfo
	bi.Header.Size = int32(unsafe.Sizeof(bi.Header))
	bi.Header.Width = int32(width)
	bi.Header.Height = -int32(height)
	bi.Header.Planes = 1
	bi.Header.BitCount = 32
	bi.Header.Compression = 0

	stride := width * 4
	buf := make([]byte, stride*height)
	ret, _, _ := procGetDIBits.Call(
		hdcMem,
		hbmp,
		0,
		uintptr(height),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bi)),
		0,
	)
	if ret == 0 {
		return nil, fmt.Errorf("GetDIBits failed")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			i := y*stride + x*4
			img.SetRGBA(x, y, color.RGBA{R: buf[i+2], G: buf[i+1], B: buf[i], A: 255})
		}
	}
	return img, nil
}

func parseUint(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func isVisible(hwnd syscall.Handle) bool {
	r, _, _ := procIsWindowVisible.Call(uintptr(hwnd))
	return r != 0
}

func windowTitle(hwnd syscall.Handle) string {
	buf := make([]uint16, 256)
	_, _, _ = procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf)
}

func windowRect(hwnd syscall.Handle) (rect, bool) {
	var r rect
	ok, _, _ := procGetWindowRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r, ok != 0
}

func clientRect(hwnd syscall.Handle) (rect, bool) {
	var r rect
	ok, _, _ := procGetClientRect.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&r)))
	return r, ok != 0
}

func (r rect) width() int32  { return r.Right - r.Left }
func (r rect) height() int32 { return r.Bottom - r.Top }
