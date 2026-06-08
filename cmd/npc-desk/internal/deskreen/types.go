package deskreen

// Phase matches Deskreen CE stepper: Connect → Select → Confirm → Sharing.
type Phase int

const (
	PhaseConnect Phase = iota
	PhaseSelect
	PhaseConfirm
	PhaseSharing
)

func (p Phase) String() string {
	switch p {
	case PhaseConnect:
		return "connect"
	case PhaseSelect:
		return "select"
	case PhaseConfirm:
		return "confirm"
	case PhaseSharing:
		return "sharing"
	default:
		return "unknown"
	}
}

// SharingType is screen (entire display) or app (single window).
type SharingType string

const (
	ShareScreen SharingType = "screen"
	ShareApp    SharingType = "app"
	ShareNone   SharingType = ""
)

// DeviceInfo is sent by the browser viewer on USER_ENTER (Deskreen parity).
type DeviceInfo struct {
	IP          string `json:"ip"`
	UserAgent   string `json:"userAgent"`
	OS          string `json:"os"`
	DeviceType  string `json:"deviceType"`
	Browser     string `json:"browser"`
	ScreenWidth int    `json:"screenWidth"`
	ScreenHeight int   `json:"screenHeight"`
}
