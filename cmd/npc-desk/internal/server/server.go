package server

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"mochi-deskreen/internal/capture"
	"mochi-deskreen/internal/deskreen"
	"mochi-deskreen/internal/netutil"
	"mochi-deskreen/internal/session"
	"mochi-deskreen/internal/signal"
	"mochi-deskreen/internal/stream"
	appweb "mochi-deskreen/web"

	"github.com/gorilla/websocket"
)

const (
	PortPrimary = "3131"
	PortBackup  = "3132"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// Config for embedded HTTP server.
type Config struct {
	Addr   string
	BindIP string
}

// Server serves Vue viewer + Deskreen CE APIs.
type Server struct {
	cfg        Config
	room       *deskreen.Room
	ctrl       *deskreen.Controller
	viewers    *session.Manager
	stream     *stream.Host
	srv        *http.Server
	onRoomChange func()
}

func New(cfg Config, room *deskreen.Room, viewers *session.Manager, streamHost *stream.Host) *Server {
	ctrl := &deskreen.Controller{Room: room, Viewers: viewers, Stream: streamHost, BindIP: cfg.BindIP}
	s := &Server{cfg: cfg, room: room, ctrl: ctrl, viewers: viewers, stream: streamHost}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/host", s.handleHost)
	mux.HandleFunc("/api/host/allow", s.handleAllow)
	mux.HandleFunc("/api/host/deny", s.handleDeny)
	mux.HandleFunc("/api/host/sources", s.handleSources)
	mux.HandleFunc("/api/host/source", s.handleSelectSource)
	mux.HandleFunc("/api/host/confirm", s.handleConfirm)
	mux.HandleFunc("/api/host/disconnect", s.handleDisconnect)
	mux.HandleFunc("/api/ws", s.handleWS)
	mux.HandleFunc("/", s.handleRoot)
	s.srv = &http.Server{
		Addr:              cfg.Addr,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	return s
}

func (s *Server) SetOnRoomChange(fn func()) {
	s.onRoomChange = fn
	s.room.SetOnChange(fn)
}

func ListenPort(bind string) (string, error) {
	for _, p := range []string{PortPrimary, PortBackup} {
		ln, err := net.Listen("tcp", bind+":"+p)
		if err == nil {
			_ = ln.Close()
			return p, nil
		}
	}
	return PortPrimary, nil
}

func (s *Server) ListenAndServe() error {
	log.Printf("deskreen http %s room=%s url=%s", s.cfg.Addr, s.room.ID, s.ViewerURL())
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) Controller() *deskreen.Controller { return s.ctrl }

func (s *Server) ViewersOnline() bool { return s.viewers.Online() }

func (s *Server) ViewerURL() string {
	host := s.cfg.BindIP
	if host == "" {
		host = netutil.LocalIPv4()
	}
	return s.room.ViewerURL(host)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"ok": true, "lan": netutil.LANAvailable()})
}

func (s *Server) handleHost(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{
		"viewerUrl":    s.ViewerURL(),
		"lanAvailable": netutil.LANAvailable(),
		"viewerOnline": s.viewers.Online(),
		"room":         s.room.PublicState(),
	})
}

func (s *Server) handleAllow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	s.ctrl.Allow()
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleDeny(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	s.ctrl.Deny()
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleSources(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("type")
	var list []capture.Source
	var err error
	switch deskreen.SharingType(t) {
	case deskreen.ShareScreen:
		list, err = capture.ListScreens()
	case deskreen.ShareApp:
		list, err = capture.ListWindows()
	default:
		http.Error(w, "type=screen|app", http.StatusBadRequest)
		return
	}
	if err != nil {
		writeJSON(w, map[string]any{"error": err.Error(), "sources": []capture.Source{}})
		return
	}
	writeJSON(w, map[string]any{"sources": list})
}

func (s *Server) handleSelectSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		SharingMode string         `json:"sharingMode"`
		Type        string         `json:"type"`
		Source      capture.Source `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	typ := body.SharingMode
	if typ == "" {
		typ = body.Type
	}
	if typ == "" {
		typ = body.Source.Type
	}
	body.Source.Type = typ
	s.room.SetSharingType(deskreen.SharingType(typ))
	s.room.SelectSource(body.Source)
	writeJSON(w, map[string]any{"ok": true, "room": s.room.PublicState()})
}

func (s *Server) handleConfirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	if !s.ctrl.ConfirmAndShare() {
		http.Error(w, "cannot start sharing", http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", http.StatusMethodNotAllowed)
		return
	}
	s.ctrl.DisconnectViewer()
	writeJSON(w, map[string]any{"ok": true})
}

func (s *Server) handleWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	if roomID == "" {
		roomID = strings.TrimPrefix(r.URL.Path, "/api/ws/")
	}
	if roomID != s.room.ID {
		log.Printf("ws: room mismatch got=%q want=%q", roomID, s.room.ID)
		writeJSONError(w, signalTypeNotAllowed())
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := &session.Conn{RoomID: roomID, Send: make(chan []byte, 64), PCMu: &sync.Mutex{}}
	remoteIP := r.RemoteAddr
	if host, _, err := net.SplitHostPort(remoteIP); err == nil {
		remoteIP = host
	}

	registered := false
	defer func() {
		log.Printf("ws: viewer disconnected room=%s", roomID)
		if registered {
			s.viewers.Clear(conn)
			if s.stream != nil {
				s.stream.CloseViewer(conn)
			}
		}
		s.room.ViewerDisconnected()
		_ = ws.Close()
	}()

	startWriter := func() {
		go func() {
			for msg := range conn.Send {
				if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
					return
				}
			}
		}()
	}

	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			return
		}
		var peek struct {
			Type string `json:"type"`
		}
		_ = json.Unmarshal(data, &peek)

		if peek.Type == signal.TypeUserEnter {
			var full struct {
				Type   string         `json:"type"`
				Device map[string]any `json:"device"`
			}
			_ = json.Unmarshal(data, &full)
			if full.Device == nil {
				full.Device = map[string]any{}
			}
			full.Device["ip"] = remoteIP
			data, _ = json.Marshal(map[string]any{"type": signal.TypeUserEnter, "device": full.Device})

			if !registered {
				if !s.viewers.Set(conn) {
					log.Printf("viewer rejected: slot busy room=%s", roomID)
					_ = ws.WriteMessage(websocket.TextMessage, signalTypeRoomLocked())
					return
				}
				registered = true
				startWriter()
			}
			s.ctrl.HandleWSMessage(data, conn)
			if s.room.HasPendingViewer() {
				log.Printf("viewer USER_ENTER pending approval ip=%s room=%s", remoteIP, roomID)
			}
			continue
		}

		if !registered {
			continue
		}
		s.ctrl.HandleWSMessage(data, conn)
	}
}

func signalTypeNotAllowed() []byte {
	b, _ := json.Marshal(map[string]string{"type": "NOT_ALLOWED"})
	return b
}

func signalTypeRoomLocked() []byte {
	b, _ := json.Marshal(map[string]string{"type": signal.TypeRoomLocked})
	return b
}

func writeJSONError(w http.ResponseWriter, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write(body)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" || path == "index.html" {
		s.serveSPA(w)
		return
	}
	if !strings.Contains(path, ".") {
		s.serveSPA(w)
		return
	}
	s.serveAsset(w, path)
}

func (s *Server) serveSPA(w http.ResponseWriter) {
	dist, err := fs.Sub(appweb.Embedded, "dist")
	if err != nil {
		http.Error(w, "embed missing", 500)
		return
	}
	b, err := fs.ReadFile(dist, "index.html")
	if err != nil {
		http.Error(w, "build webui first", 404)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(b)
}

func (s *Server) serveAsset(w http.ResponseWriter, path string) {
	dist, _ := fs.Sub(appweb.Embedded, "dist")
	if b, err := fs.ReadFile(dist, path); err == nil {
		w.Header().Set("Content-Type", contentType(path))
		_, _ = w.Write(b)
		return
	}
	s.serveSPA(w)
}

func contentType(path string) string {
	switch {
	case strings.HasSuffix(path, ".html"):
		return "text/html; charset=utf-8"
	case strings.HasSuffix(path, ".js"):
		return "application/javascript"
	case strings.HasSuffix(path, ".css"):
		return "text/css"
	case strings.HasSuffix(path, ".svg"):
		return "image/svg+xml"
	case strings.HasSuffix(path, ".ico"):
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
