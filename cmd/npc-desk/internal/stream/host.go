package stream

import (
	"encoding/json"
	"log"
	"strings"
	"sync"
	"time"

	"npc-deskreen/internal/capture"
	"npc-deskreen/internal/session"
	"npc-deskreen/internal/signal"

	"github.com/pion/webrtc/v4"
)

// Host manages WebRTC screen sharing (Deskreen CE flow).
type Host struct {
	mu        sync.Mutex
	jpegQuality int
	maxWidth    int
	source    capture.Source
	shareMode string
	stopCh    chan struct{}
}

func NewHost() (*Host, error) {
	return &Host{jpegQuality: 72, maxWidth: 1280}, nil
}

// SetQuality applies viewer quality preset (0 = auto). Adjusts JPEG quality and max encode width.
func (h *Host) SetQuality(pct int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	switch {
	case pct <= 0:
		h.jpegQuality = 72
		h.maxWidth = 1280
	case pct <= 25:
		h.jpegQuality = 38
		h.maxWidth = 640
	case pct <= 40:
		h.jpegQuality = 50
		h.maxWidth = 800
	case pct <= 60:
		h.jpegQuality = 62
		h.maxWidth = 1024
	case pct <= 80:
		h.jpegQuality = 78
		h.maxWidth = 1280
	default:
		h.jpegQuality = 90
		h.maxWidth = 1600
	}
	log.Printf("stream: quality preset %d%% -> jpeg=%d maxW=%d", pct, h.jpegQuality, h.maxWidth)
}

func (h *Host) StartSharing(c *session.Conn, src capture.Source, mode string) {
	h.mu.Lock()
	h.source = src
	h.shareMode = mode
	if h.stopCh != nil {
		close(h.stopCh)
	}
	h.stopCh = make(chan struct{})
	stop := h.stopCh
	h.mu.Unlock()
	log.Printf("stream: StartSharing mode=%s src.type=%s id=%s name=%q", mode, src.Type, src.ID, src.Name)
	go h.startSession(c, stop)
}

// StopSharing ends the screen capture loop (call before closing the viewer send channel).
func (h *Host) StopSharing() {
	h.mu.Lock()
	if h.stopCh != nil {
		close(h.stopCh)
		h.stopCh = nil
	}
	h.source = capture.Source{}
	h.mu.Unlock()
}

func (h *Host) HandleViewerMessage(c *session.Conn, raw []byte) {
	var env signal.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return
	}
	switch strings.ToLower(env.Type) {
	case signal.TypeAnswer:
		h.applyAnswer(c, env)
	case signal.TypeICE:
		h.applyICE(c, env)
	case strings.ToLower(signal.TypeQuality):
		h.SetQuality(env.Quality)
	}
}

func (h *Host) startSession(c *session.Conn, stop <-chan struct{}) {
	c.PCMu.Lock()
	if c.PC != nil {
		_ = c.PC.Close()
		c.PC = nil
	}
	c.PendingICE = nil
	c.PCMu.Unlock()

	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
			{URLs: []string{"stun:stun1.l.google.com:19302"}},
		},
	})
	if err != nil {
		log.Printf("webrtc: %v", err)
		h.sendSignal(c, signal.Envelope{Type: signal.TypeSharingFailed})
		return
	}

	c.PCMu.Lock()
	c.PC = pc
	c.PCMu.Unlock()

	dc, err := pc.CreateDataChannel("screen", nil)
	if err != nil {
		log.Printf("webrtc dc: %v", err)
		h.closePC(c)
		h.sendSignal(c, signal.Envelope{Type: signal.TypeSharingFailed})
		return
	}

	// Start capture immediately; use WebSocket until the data channel opens (LAN-friendly fallback).
	go h.runCapture(dc, c, stop)

	dc.OnOpen(func() {
		log.Printf("webrtc: datachannel open")
	})

	pc.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Printf("webrtc: connection state %s", s.String())
		if s == webrtc.PeerConnectionStateFailed {
			h.sendSignal(c, signal.Envelope{Type: signal.TypeSharingFailed})
		}
	})

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate == nil {
			return
		}
		init := candidate.ToJSON()
		h.sendSignal(c, signal.Envelope{Type: signal.TypeICE, Candidate: &init})
	})

	offer, err := pc.CreateOffer(nil)
	if err != nil {
		log.Printf("webrtc offer: %v", err)
		h.sendSignal(c, signal.Envelope{Type: signal.TypeSharingFailed})
		return
	}
	if err := pc.SetLocalDescription(offer); err != nil {
		log.Printf("webrtc set local: %v", err)
		h.sendSignal(c, signal.Envelope{Type: signal.TypeSharingFailed})
		return
	}
	h.sendSignal(c, signal.Envelope{Type: signal.TypeOffer, SDP: offer.SDP, SDPType: offer.Type.String()})
	log.Printf("webrtc: offer sent to viewer")
}

func (h *Host) sendSignal(c *session.Conn, env signal.Envelope) {
	b, err := json.Marshal(env)
	if err != nil {
		return
	}
	h.send(c, b)
}

func (h *Host) applyAnswer(c *session.Conn, env signal.Envelope) {
	c.PCMu.Lock()
	defer c.PCMu.Unlock()
	pc := c.PC
	if pc == nil {
		log.Printf("webrtc: answer but no peer connection")
		return
	}
	if err := pc.SetRemoteDescription(webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  env.SDP,
	}); err != nil {
		log.Printf("webrtc answer: %v", err)
		go h.sendSignal(c, signal.Envelope{Type: signal.TypeSharingFailed})
		return
	}
	pending := c.PendingICE
	c.PendingICE = nil
	for _, cand := range pending {
		if err := pc.AddICECandidate(cand); err != nil {
			log.Printf("webrtc flush ice: %v", err)
		}
	}
	log.Printf("webrtc: answer applied (pending ICE=%d)", len(pending))
}

func (h *Host) applyICE(c *session.Conn, env signal.Envelope) {
	if env.Candidate == nil {
		return
	}
	c.PCMu.Lock()
	defer c.PCMu.Unlock()
	pc := c.PC
	if pc == nil {
		return
	}
	if pc.RemoteDescription() == nil {
		c.PendingICE = append(c.PendingICE, *env.Candidate)
		return
	}
	if err := pc.AddICECandidate(*env.Candidate); err != nil {
		log.Printf("webrtc ice: %v", err)
	}
}

func (h *Host) send(c *session.Conn, msg []byte) {
	if !c.Enqueue(msg, 15*time.Second) {
		log.Printf("stream: viewer send queue full (len=%d)", len(msg))
	}
}

func (h *Host) closePC(c *session.Conn) {
	c.PCMu.Lock()
	if c.PC != nil {
		_ = c.PC.Close()
		c.PC = nil
	}
	c.PendingICE = nil
	c.PCMu.Unlock()
}

func (h *Host) CloseViewer(c *session.Conn) {
	h.StopSharing()
	h.closePC(c)
}
