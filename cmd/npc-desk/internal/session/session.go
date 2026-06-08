package session

import (
	"log"
	"sync"
	"time"

	"github.com/pion/webrtc/v4"
)

// Conn is a single viewer WebSocket + optional WebRTC peer.
type Conn struct {
	RoomID string
	Send   chan []byte
	PC     *webrtc.PeerConnection
	PCMu   *sync.Mutex
	// PendingICE holds viewer trickle candidates until the SDP answer is applied.
	PendingICE []webrtc.ICECandidateInit

	sendMu     sync.Mutex
	sendClosed bool
}

// Manager holds the single active viewer connection (Deskreen CE: one slot).
type Manager struct {
	mu   sync.RWMutex
	conn *Conn
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Set(c *Conn) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.conn != nil {
		return false
	}
	m.conn = c
	return true
}

func (m *Manager) Get() *Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conn
}

func (m *Manager) Clear(c *Conn) {
	m.mu.Lock()
	if m.conn == c {
		m.conn = nil
	}
	m.mu.Unlock()
}

func (m *Manager) CloseActive() *Conn {
	m.mu.Lock()
	c := m.conn
	m.conn = nil
	m.mu.Unlock()
	if c != nil {
		c.CloseSend()
	}
	return c
}

func (m *Manager) Online() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conn != nil
}

// CloseSend marks the outbound queue closed (stops capture/signaling from sending; ends the WS writer loop).
func (c *Conn) CloseSend() {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	if c.sendClosed {
		return
	}
	c.sendClosed = true
	close(c.Send)
}

// Enqueue delivers a message to the viewer send queue (blocks up to timeout; signaling must not be dropped).
func (c *Conn) Enqueue(msg []byte, timeout time.Duration) bool {
	if timeout <= 0 {
		timeout = 15 * time.Second
	}
	c.sendMu.Lock()
	if c.sendClosed {
		c.sendMu.Unlock()
		return false
	}
	c.sendMu.Unlock()
	select {
	case c.Send <- msg:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (m *Manager) Send(msg []byte) {
	m.mu.RLock()
	c := m.conn
	m.mu.RUnlock()
	if c == nil {
		return
	}
	if !c.Enqueue(msg, 15*time.Second) {
		log.Printf("session: viewer send queue full, message dropped (len=%d)", len(msg))
	}
}
