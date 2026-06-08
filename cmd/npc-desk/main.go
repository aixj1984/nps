package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewNpcDeskApp()
	webAssets, err := fs.Sub(assets, "frontend/dist")
	if err != nil {
		log.Fatal(err)
	}

	err = wails.Run(&options.App{
		Title:            "NPC Desk",
		Width:            1080,
		Height:           760,
		MinWidth:         900,
		MinHeight:        640,
		DisableResize:    false,
		AssetServer:      &assetserver.Options{Assets: webAssets},
		BackgroundColour: &options.RGBA{R: 245, G: 248, B: 250, A: 255},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind:             []any{app},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			Theme:                windows.SystemDefault,
		},
	})
	if err != nil {
		writeStartupLog(err)
		log.Fatal(err)
	}
}

func writeStartupLog(err error) {
	exe, e := os.Executable()
	if e != nil {
		exe = "."
	}
	path := filepath.Join(filepath.Dir(exe), "npc-desk-startup.log")
	msg := fmt.Sprintf("[%s] startup failed: %v\n", time.Now().Format(time.RFC3339), err)
	_ = os.WriteFile(path, []byte(msg), 0o644)
}
