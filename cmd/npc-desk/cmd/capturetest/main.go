// Capture test tool: go run ./cmd/capturetest -mode app -id <hwnd> -out shot.png
package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"

	"mochi-deskreen/internal/capture"
)

func main() {
	mode := flag.String("mode", "app", "app or screen")
	id := flag.String("id", "0", "window HWND (decimal) or display index")
	out := flag.String("out", "capture-test.png", "output PNG")
	flag.Parse()

	src := capture.Source{ID: *id, Name: "test", Type: *mode}
	img, err := capture.CaptureWithMode(src, *mode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "capture failed: %v\n", err)
		os.Exit(1)
	}
	b := img.Bounds()
	fmt.Printf("ok: %dx%d -> %s\n", b.Dx(), b.Dy(), *out)
	f, err := os.Create(*out)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_ = png.Encode(f, img)
}
