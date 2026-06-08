package netutil

import "net"

// LANAvailable reports whether a non-loopback interface is up (Deskreen WiFi/LAN guard).
func LANAvailable() bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipNet, ok := a.(*net.IPNet)
			if ok && ipNet.IP != nil && ipNet.IP.To4() != nil && !ipNet.IP.IsLoopback() {
				return true
			}
		}
	}
	return false
}
