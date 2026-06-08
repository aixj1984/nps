package deskreen

import (
	"sync"

	"mochi-deskreen/internal/capture"

	"github.com/google/uuid"
)

// Room is a Deskreen-style waiting/sharing room (single viewer in CE).
type Room struct {
	mu sync.RWMutex

	ID   string
	Port string

	Phase        Phase
	SharingType  SharingType
	Selected     capture.Source
	Pending      *DeviceInfo
	ConnectedIP  string
	ViewerOnline bool
	Lang         string

	onChange func()
}

func NewRoom(port string) *Room {
	return &Room{
		ID:    uuid.New().String(),
		Port:  port,
		Phase: PhaseConnect,
		Lang:  "zh_CN",
	}
}

func (r *Room) SetOnChange(fn func()) {
	r.mu.Lock()
	r.onChange = fn
	r.mu.Unlock()
}

func (r *Room) notify() {
	r.mu.RLock()
	fn := r.onChange
	r.mu.RUnlock()
	if fn != nil {
		fn()
	}
}

func (r *Room) Snapshot() map[string]any {
	r.mu.RLock()
	defer r.mu.RUnlock()
	pending := any(nil)
	if r.Pending != nil {
		pending = r.Pending
	}
	return map[string]any{
		"roomId":        r.ID,
		"phase":         r.Phase.String(),
		"sharingType":   string(r.SharingType),
		"selectedSource": r.Selected,
		"pendingDevice":  pending,
		"viewerOnline":  r.ViewerOnline,
		"connectedIp":   r.ConnectedIP,
		"lang":          r.Lang,
	}
}

func (r *Room) ResetWaiting() {
	r.mu.Lock()
	r.Phase = PhaseConnect
	r.SharingType = ShareNone
	r.Selected = capture.Source{}
	r.Pending = nil
	r.ViewerOnline = false
	r.ConnectedIP = ""
	r.ID = uuid.New().String()
	r.mu.Unlock()
	r.notify()
}

func (r *Room) SetPending(d DeviceInfo) bool {
	r.mu.Lock()
	if r.Phase != PhaseConnect {
		r.mu.Unlock()
		return false
	}
	r.Pending = &d
	r.ViewerOnline = true
	r.ConnectedIP = d.IP
	r.mu.Unlock()
	r.notify()
	return true
}

func (r *Room) Allow() {
	r.mu.Lock()
	if r.Pending == nil {
		r.mu.Unlock()
		return
	}
	r.Phase = PhaseSelect
	r.mu.Unlock()
	r.notify()
}

func (r *Room) Deny() {
	r.mu.Lock()
	r.Pending = nil
	r.ViewerOnline = false
	r.ConnectedIP = ""
	r.Phase = PhaseConnect
	r.mu.Unlock()
	r.notify()
}

func (r *Room) SetSharingType(t SharingType) {
	r.mu.Lock()
	r.SharingType = t
	r.Selected = capture.Source{}
	r.mu.Unlock()
	r.notify()
}

func (r *Room) SelectSource(src capture.Source) {
	r.mu.Lock()
	r.Selected = src
	r.Phase = PhaseConfirm
	r.mu.Unlock()
	r.notify()
}

func (r *Room) StartSharing() (capture.Source, bool) {
	r.mu.Lock()
	if r.Selected.ID == "" || r.Phase != PhaseConfirm {
		r.mu.Unlock()
		return capture.Source{}, false
	}
	r.Phase = PhaseSharing
	src := r.Selected
	switch r.SharingType {
	case ShareApp:
		src.Type = "app"
	case ShareScreen:
		src.Type = "screen"
	default:
		if src.Type == "" && r.SharingType != ShareNone {
			src.Type = string(r.SharingType)
		}
	}
	r.mu.Unlock()
	r.notify()
	return src, true
}

func (r *Room) DisconnectViewer() {
	r.mu.Lock()
	r.Pending = nil
	r.ViewerOnline = false
	r.ConnectedIP = ""
	r.SharingType = ShareNone
	r.Selected = capture.Source{}
	r.Phase = PhaseConnect
	r.ID = uuid.New().String()
	r.mu.Unlock()
	r.notify()
}

// ViewerDisconnected clears viewer state when the browser socket closes (keep room URL).
func (r *Room) ViewerDisconnected() {
	r.mu.Lock()
	r.Pending = nil
	r.ViewerOnline = false
	r.ConnectedIP = ""
	if r.Phase == PhaseSharing || r.Phase == PhaseSelect || r.Phase == PhaseConfirm {
		r.Phase = PhaseConnect
		r.SharingType = ShareNone
		r.Selected = capture.Source{}
	}
	r.mu.Unlock()
	r.notify()
}

func (r *Room) CanAcceptViewer() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Phase == PhaseConnect && r.Pending == nil
}

// HasPendingViewer reports whether a device is waiting for host approval.
// SharingTypeLocked returns the current share mode (screen / app).
func (r *Room) SharingTypeLocked() SharingType {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.SharingType
}

func (r *Room) HasPendingViewer() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Pending != nil && r.Phase == PhaseConnect
}

func (r *Room) ViewerURL(hostIP string) string {
	return "http://" + hostIP + ":" + r.Port + "/" + r.ID
}
