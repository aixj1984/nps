package deskreen

import "mochi-deskreen/internal/capture"

// RoomState is JSON returned to the Fyne wizard (/api/host).
type RoomState struct {
	RoomID         string          `json:"roomId"`
	Phase          string          `json:"phase"`
	SharingType    string          `json:"sharingType"`
	SelectedSource capture.Source  `json:"selectedSource"`
	PendingDevice  *DeviceInfo     `json:"pendingDevice"`
	ViewerOnline   bool            `json:"viewerOnline"`
	ConnectedIP    string          `json:"connectedIp"`
}

func (r *Room) PublicState() RoomState {
	r.mu.RLock()
	defer r.mu.RUnlock()
	st := RoomState{
		RoomID:       r.ID,
		Phase:        r.Phase.String(),
		SharingType:  string(r.SharingType),
		SelectedSource: r.Selected,
		ViewerOnline: r.ViewerOnline,
		ConnectedIP:  r.ConnectedIP,
	}
	if r.Pending != nil {
		p := *r.Pending
		st.PendingDevice = &p
	}
	return st
}
