package deskreen

import (
	"encoding/json"
	"log"
	"strings"

	"mochi-deskreen/internal/session"
	"mochi-deskreen/internal/signal"
	"mochi-deskreen/internal/stream"
)

// Controller ties Room state, viewer socket, and streaming together.
type Controller struct {
	Room    *Room
	Viewers *session.Manager
	Stream  *stream.Host
	BindIP  string
}

func (c *Controller) SendViewer(msg signal.Envelope) {
	b, _ := json.Marshal(msg)
	c.Viewers.Send(b)
}

func (c *Controller) HandleWSMessage(raw []byte, conn *session.Conn) {
	var env signal.Envelope
	if err := json.Unmarshal(raw, &env); err != nil {
		return
	}

	switch strings.ToLower(env.Type) {
	case strings.ToLower(signal.TypeUserEnter):
		if env.Device == nil {
			env.Device = map[string]any{}
		}
		dev := mapToDevice(env.Device)
		if !c.Room.CanAcceptViewer() {
			c.SendViewer(signal.Envelope{Type: signal.TypeRoomLocked})
			return
		}
		if !c.Room.SetPending(dev) {
			c.SendViewer(signal.Envelope{Type: signal.TypeRoomLocked})
			return
		}
	case strings.ToLower(signal.TypeUserExit):
		c.DisconnectViewer()
	case strings.ToLower(signal.TypeAnswer), strings.ToLower(signal.TypeICE):
		if c.Stream != nil {
			c.Stream.HandleViewerMessage(conn, raw)
		}
	case strings.ToLower(signal.TypeQuality):
		if c.Stream != nil {
			c.Stream.SetQuality(env.Quality)
		}
	default:
		if c.Stream != nil {
			c.Stream.HandleViewerMessage(conn, raw)
		}
	}
}

func (c *Controller) Allow() {
	c.Room.Allow()
	c.SendViewer(signal.Envelope{Type: signal.TypeAllowedToConnect})
	c.NotifyWaitingSource()
}

func (c *Controller) Deny() {
	c.SendViewer(signal.Envelope{Type: signal.TypeDenyToConnect})
	if c.Stream != nil {
		c.Stream.StopSharing()
	}
	if conn := c.Viewers.CloseActive(); conn != nil && c.Stream != nil {
		c.Stream.CloseViewer(conn)
	}
	c.Room.Deny()
}

func (c *Controller) NotifyWaitingSource() {
	c.SendViewer(signal.Envelope{Type: signal.TypeWaitingForSource})
}

func (c *Controller) ConfirmAndShare() bool {
	conn := c.Viewers.Get()
	if conn == nil {
		log.Printf("deskreen: ConfirmAndShare: no viewer websocket")
		return false
	}
	if c.Stream == nil {
		return false
	}
	src, ok := c.Room.StartSharing()
	if !ok {
		return false
	}
	mode := string(c.Room.SharingTypeLocked())
	if mode == "" {
		mode = src.Type
	}
	log.Printf("deskreen: ConfirmAndShare mode=%q src.type=%q id=%s name=%q", mode, src.Type, src.ID, src.Name)
	c.SendViewer(signal.Envelope{Type: signal.TypeSharingStarted})
	c.Stream.StartSharing(conn, src, mode)
	return true
}

func (c *Controller) DisconnectViewer() {
	c.SendViewer(signal.Envelope{Type: signal.TypeDisconnectByHost})
	if c.Stream != nil {
		c.Stream.StopSharing()
	}
	conn := c.Viewers.CloseActive()
	if conn != nil && c.Stream != nil {
		c.Stream.CloseViewer(conn)
	}
	c.Room.DisconnectViewer()
}
