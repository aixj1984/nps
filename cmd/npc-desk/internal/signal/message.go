package signal

import "github.com/pion/webrtc/v4"

// Envelope mirrors Deskreen / darkwire app messages over WebSocket.
type Envelope struct {
	Type string `json:"type"`

	SDP       string                   `json:"sdp,omitempty"`
	SDPType   string                   `json:"sdpType,omitempty"`
	Candidate *webrtc.ICECandidateInit `json:"candidate,omitempty"`

	Device  map[string]any `json:"device,omitempty"`
	Quality int            `json:"quality,omitempty"`
	Data    string         `json:"data,omitempty"`
}

const (
	TypeUserEnter        = "USER_ENTER"
	TypeUserExit         = "USER_EXIT"
	TypeAllowedToConnect = "ALLOWED_TO_CONNECT"
	TypeDenyToConnect    = "DENY_TO_CONNECT"
	TypeWaitingForSource = "WAITING_FOR_SOURCE"
	TypeSharingStarted   = "SHARING_STARTED"
	TypeSharingFailed    = "SHARING_FAILED"
	TypeDisconnectByHost = "DISCONNECT_BY_HOST"
	TypeNotAllowed       = "NOT_ALLOWED"
	TypeRoomLocked       = "ROOM_LOCKED"

	TypeOffer  = "offer"
	TypeAnswer = "answer"
	TypeICE    = "ice"

	TypeScreenFrame = "screen-frame"
	TypeQuality     = "QUALITY"
)
