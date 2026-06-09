package stream

import (
	"encoding/json"
	"image"
	"log"
	"time"

	"npc-deskreen/internal/capture"
	"npc-deskreen/internal/session"

	"github.com/nfnt/resize"
	"github.com/pion/webrtc/v4"
)

const screenFPS = 10

func (h *Host) runCapture(dc *webrtc.DataChannel, c *session.Conn, stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second / screenFPS)
	defer ticker.Stop()

	var captureErrLogged bool
	var transportLogged string
	var modeLogged bool
	var lastGood image.Image
	var lastGoodAt time.Time

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
		}

		h.mu.Lock()
		src := h.source
		mode := h.shareMode
		jpegQ := h.jpegQuality
		maxW := h.maxWidth
		h.mu.Unlock()

		if src.ID == "" {
			continue
		}
		if !modeLogged {
			modeLogged = true
			log.Printf("stream: capturing mode=%s type=%s id=%s name=%q", mode, src.Type, src.ID, src.Name)
		}

		frame, err := capture.CaptureWithMode(src, mode)
		if err != nil {
			if lastGood != nil && time.Since(lastGoodAt) < 2*time.Second {
				frame = lastGood
			} else {
				if !captureErrLogged {
					captureErrLogged = true
					log.Printf("capture mode=%s type=%s id=%s: %v", mode, src.Type, src.ID, err)
				}
				continue
			}
		} else {
			lastGood = frame
			lastGoodAt = time.Now()
		}

		bounds := frame.Bounds()
		fw := bounds.Dx()
		if maxW > 0 && fw > maxW {
			// Bilinear is faster than Lanczos and reduces resize flicker on app windows.
			frame = resize.Resize(uint(maxW), 0, frame, resize.Bilinear)
		}

		var buf bytesBuffer
		if err := encodeJPEG(&buf, frame, jpegQ); err != nil {
			continue
		}
		payload, err := json.Marshal(map[string]string{
			"type": typeScreenFrame,
			"data": bytesToBase64(buf.Bytes()),
		})
		if err != nil {
			continue
		}

		dcOpen := dc != nil && dc.ReadyState() == webrtc.DataChannelStateOpen
		if dcOpen {
			if err := dc.Send(payload); err == nil {
				if transportLogged != "dc" {
					transportLogged = "dc"
					log.Printf("stream: frames via webrtc datachannel")
				}
				continue
			}
		}
		// Only use WebSocket when DataChannel is not ready (avoid dual-path flicker).
		if c != nil {
			h.send(c, payload)
			if transportLogged != "ws" {
				transportLogged = "ws"
				log.Printf("stream: frames via websocket (datachannel not ready)")
			}
		}
	}
}
