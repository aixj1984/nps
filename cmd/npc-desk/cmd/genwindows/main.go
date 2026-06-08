// Generates build/windows/icon.ico and mochi-deskreen-res.syso
// See https://wails.io/docs/guides/manual-builds/
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/leaanthony/winicon"
	"github.com/tc-hib/winres"
	"github.com/tc-hib/winres/version"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	buildDir := filepath.Join(root, "build")
	winDir := filepath.Join(buildDir, "windows")
	if err := os.MkdirAll(winDir, 0o755); err != nil {
		fatal(err)
	}

	appicon := filepath.Join(buildDir, "appicon.png")
	if _, err := os.Stat(appicon); err != nil {
		fatal(fmt.Errorf("missing %s", appicon))
	}

	icoPath := filepath.Join(winDir, "icon.ico")
	_ = os.Remove(icoPath)
	png, err := os.ReadFile(appicon)
	if err != nil {
		fatal(err)
	}
	icoFile, err := os.OpenFile(icoPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		fatal(err)
	}
	if err := winicon.GenerateIcon(bytes.NewReader(png), icoFile, []int{256, 128, 64, 48, 32, 16}); err != nil {
		_ = icoFile.Close()
		fatal(err)
	}
	_ = icoFile.Close()
	fmt.Println("icon:", icoPath)

	manifest, err := os.ReadFile(filepath.Join(winDir, "wails.exe.manifest"))
	if err != nil {
		fatal(err)
	}

	f, err := os.Open(icoPath)
	if err != nil {
		fatal(err)
	}
	ico, err := winres.LoadICO(f)
	_ = f.Close()
	if err != nil {
		fatal(err)
	}

	rs := winres.ResourceSet{}
	if err := rs.SetIcon(winres.RT_ICON, ico); err != nil {
		fatal(err)
	}
	xmlData, err := winres.AppManifestFromXML(manifest)
	if err != nil {
		fatal(err)
	}
	rs.SetManifest(xmlData)

	if raw, err := os.ReadFile(filepath.Join(winDir, "info.json")); err == nil && len(bytes.TrimSpace(raw)) > 0 && !bytes.Contains(raw, []byte("{{")) {
		var v version.Info
		if err := v.UnmarshalJSON(raw); err != nil {
			fatal(err)
		}
		rs.SetVersionInfo(v)
	}

	sysoPath := filepath.Join(root, strings.ReplaceAll("mochi-deskreen", " ", "_")+"-res.syso")
	out, err := os.Create(sysoPath)
	if err != nil {
		fatal(err)
	}
	if err := rs.WriteObject(out, winres.ArchAMD64); err != nil {
		_ = out.Close()
		fatal(err)
	}
	_ = out.Close()
	fmt.Println("syso:", sysoPath)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "genwindows:", err)
	os.Exit(1)
}
