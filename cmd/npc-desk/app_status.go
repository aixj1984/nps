package main

import (
	"encoding/json"

	"npc-deskreen/internal/deskreen"
	"npc-deskreen/internal/netutil"
)

// statusMap builds a JSON-friendly map for Wails EventsEmit and Vue (avoids nested struct IPC issues).
func (a *HostApp) statusMap() map[string]any {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.statusMapLocked()
}

func (a *HostApp) statusMapLocked() map[string]any {
	room := a.room.PublicState()
	roomMap := map[string]any{
		"roomId":       room.RoomID,
		"phase":        room.Phase,
		"sharingType":  room.SharingType,
		"viewerOnline": room.ViewerOnline,
		"connectedIp":  room.ConnectedIP,
		"selectedSource": map[string]any{
			"id":   room.SelectedSource.ID,
			"name": room.SelectedSource.Name,
			"type": room.SelectedSource.Type,
		},
	}
	if room.PendingDevice != nil {
		roomMap["pendingDevice"] = deviceMap(*room.PendingDevice)
	}
	return map[string]any{
		"viewerUrl":    a.srv.ViewerURL(),
		"lanAvailable": netutil.LANAvailable(),
		"viewerOnline": room.ViewerOnline || a.srv.ViewersOnline(),
		"room":         roomMap,
	}
}

func deviceMap(d deskreen.DeviceInfo) map[string]any {
	return map[string]any{
		"ip":           d.IP,
		"userAgent":    d.UserAgent,
		"os":           d.OS,
		"deviceType":   d.DeviceType,
		"browser":      d.Browser,
		"screenWidth":  d.ScreenWidth,
		"screenHeight": d.ScreenHeight,
	}
}

// hostStatusJSON is used when a typed struct must be marshaled consistently.
func hostStatusJSON(st HostStatus) map[string]any {
	b, _ := json.Marshal(st)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	return m
}
