package deskreen

import "strconv"

func mapToDevice(m map[string]any) DeviceInfo {
	d := DeviceInfo{}
	if v, ok := m["ip"].(string); ok {
		d.IP = v
	}
	if v, ok := m["userAgent"].(string); ok {
		d.UserAgent = v
	}
	if v, ok := m["os"].(string); ok {
		d.OS = v
	}
	if v, ok := m["deviceType"].(string); ok {
		d.DeviceType = v
	}
	if v, ok := m["browser"].(string); ok {
		d.Browser = v
	}
	d.ScreenWidth = intFromAny(m["screenWidth"])
	d.ScreenHeight = intFromAny(m["screenHeight"])
	return d
}

func intFromAny(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case float32:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)
		return i
	default:
		return 0
	}
}
