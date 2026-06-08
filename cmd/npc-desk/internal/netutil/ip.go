package netutil

import (
	"net"
	"strings"
)

// LocalIPv4 returns the first non-loopback IPv4 address, or 127.0.0.1.
func LocalIPv4() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "127.0.0.1"
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
			if !ok || ipNet.IP == nil || ipNet.IP.To4() == nil {
				continue
			}
			ip := ipNet.IP.String()
			if strings.HasPrefix(ip, "169.254.") {
				continue
			}
			return ip
		}
	}
	return "127.0.0.1"
}
