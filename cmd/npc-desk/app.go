package main

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"npc-deskreen/internal/capture"
	"npc-deskreen/internal/deskreen"
	"npc-deskreen/internal/netutil"
	"npc-deskreen/internal/server"
	"npc-deskreen/internal/session"
	"npc-deskreen/internal/stream"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// HostApp is exposed to the Wails Vue frontend (Deskreen-style host UI).
type HostApp struct {
	ctx context.Context

	mu     sync.Mutex
	room   *deskreen.Room
	ctrl   *deskreen.Controller
	srv    *server.Server
	port   string
	bindIP string
}

// HostStatus is returned to the Vue host UI.
type HostStatus struct {
	ViewerURL    string             `json:"viewerUrl"`
	LanAvailable bool               `json:"lanAvailable"`
	ViewerOnline bool               `json:"viewerOnline"`
	Room         deskreen.RoomState `json:"room"`
}

// CaptureSourceDTO is sent to the Vue UI (includes base64 thumbnail).
type CaptureSourceDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Thumb string `json:"thumb,omitempty"`
}

// SelectSourceRequest picks screen or app capture source.
// Use sharingMode (not "type") — Wails/JSON may drop a field named "type".
type SelectSourceRequest struct {
	SharingMode string           `json:"sharingMode"`
	Source      CaptureSourceDTO `json:"source"`
}

func NewHostApp() *HostApp {
	bindIP := flagLocalIP()
	port, err := server.ListenPort("0.0.0.0")
	if err != nil {
		port = server.PortPrimary
	}
	room := deskreen.NewRoom(port)
	viewers := session.NewManager()
	streamHost, err := stream.NewHost()
	if err != nil {
		log.Fatalf("stream: %v", err)
	}
	srv := server.New(server.Config{Addr: ":" + port, BindIP: bindIP}, room, viewers, streamHost)
	return &HostApp{
		room: room,
		ctrl: srv.Controller(),
		srv:  srv,
		port: port,
		bindIP: bindIP,
	}
}

func (a *HostApp) startup(ctx context.Context) {
	a.ctx = ctx
	a.room.SetOnChange(a.emitState)
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("http: %v", err)
		}
	}()
	a.emitState()
}

func (a *HostApp) shutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = a.srv.Shutdown(shutdownCtx)
}

func (a *HostApp) emitState() {
	if a.ctx == nil {
		return
	}
	runtime.EventsEmit(a.ctx, "room:update", a.statusMap())
}

// GetStatus returns current host/room state for the Vue UI (plain map for reliable JSON).
func (a *HostApp) GetStatus() map[string]any {
	return a.statusMap()
}

func (a *HostApp) Allow() error {
	a.ctrl.Allow()
	a.emitState()
	return nil
}

func (a *HostApp) Deny() error {
	a.ctrl.Deny()
	a.emitState()
	return nil
}

func (a *HostApp) ListSources(typ string) ([]CaptureSourceDTO, error) {
	// Lock share mode as soon as the user picks screen vs app (before thumbnail selection).
	a.room.SetSharingType(deskreen.SharingType(typ))
	var list []capture.Source
	var err error
	switch deskreen.SharingType(typ) {
	case deskreen.ShareScreen:
		list, err = capture.ListScreens()
	case deskreen.ShareApp:
		list, err = capture.ListWindows()
	default:
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	out := make([]CaptureSourceDTO, len(list))
	for i, s := range list {
		out[i] = CaptureSourceDTO{ID: s.ID, Name: s.Name, Type: s.Type}
		if len(s.ThumbPNG) > 0 {
			out[i].Thumb = base64.StdEncoding.EncodeToString(s.ThumbPNG)
		}
	}
	return out, nil
}

// SetSharingMode records whether the host will share a full display or one app window.
func (a *HostApp) SetSharingMode(mode string) error {
	a.room.SetSharingType(deskreen.SharingType(mode))
	a.emitState()
	return nil
}

func (a *HostApp) SelectSource(req SelectSourceRequest) error {
	mode := req.SharingMode
	if mode == "" {
		mode = string(a.room.SharingTypeLocked())
	}
	if mode == "" {
		mode = req.Source.Type
	}
	if mode == "app" {
		return a.SelectAppWindow(req.Source.ID, req.Source.Name)
	}
	if mode == "screen" {
		return a.SelectScreenDisplay(req.Source.ID, req.Source.Name)
	}
	typ := deskreen.SharingType(mode)
	a.room.SetSharingType(typ)
	a.room.SelectSource(capture.Source{
		ID:   req.Source.ID,
		Name: req.Source.Name,
		Type: string(typ),
	})
	a.emitState()
	return nil
}

// SelectAppWindow selects one application window (HWND as decimal string). Prefer over SelectSource for Wails.
func (a *HostApp) SelectAppWindow(windowID, windowName string) error {
	log.Printf("host: SelectAppWindow id=%s name=%q", windowID, windowName)
	a.room.SetSharingType(deskreen.ShareApp)
	a.room.SelectSource(capture.Source{
		ID:   windowID,
		Name: windowName,
		Type: "app",
	})
	a.emitState()
	return nil
}

// SelectScreenDisplay selects one monitor by index string ("0", "1", …).
func (a *HostApp) SelectScreenDisplay(displayID, displayName string) error {
	log.Printf("host: SelectScreenDisplay id=%s name=%q", displayID, displayName)
	a.room.SetSharingType(deskreen.ShareScreen)
	a.room.SelectSource(capture.Source{
		ID:   displayID,
		Name: displayName,
		Type: "screen",
	})
	a.emitState()
	return nil
}

func (a *HostApp) ConfirmShare() error {
	if !a.ctrl.ConfirmAndShare() {
		return errCannotShare
	}
	a.emitState()
	return nil
}

func (a *HostApp) Disconnect() error {
	a.ctrl.DisconnectViewer()
	a.emitState()
	return nil
}

func (a *HostApp) CopyViewerURL() string {
	url := a.srv.ViewerURL()
	runtime.ClipboardSetText(a.ctx, url)
	return url
}

func (a *HostApp) OpenViewerURL() {
	runtime.BrowserOpenURL(a.ctx, a.srv.ViewerURL())
}

// GetBindIP returns advertised LAN IP for the UI header.
func (a *HostApp) GetBindIP() string { return a.bindIP }

// GetPort returns HTTP server port (3131 or 3132).
func (a *HostApp) GetPort() string { return a.port }

var errCannotShare = &wailsError{msg: "cannot start sharing: pick a source first"}

type wailsError struct{ msg string }

func (e *wailsError) Error() string { return e.msg }

func flagLocalIP() string {
	args := os.Args[1:]
	for i, arg := range args {
		if (arg == "--ip" || arg == "--local-ip") && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, "--ip=") {
			return strings.TrimPrefix(arg, "--ip=")
		}
	}
	return netutil.LocalIPv4()
}
