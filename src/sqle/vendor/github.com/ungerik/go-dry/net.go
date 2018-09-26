package dry

import (
	"net"
	"os"
)

// NetIP returns the primary IP address of the system or an empty string.
func NetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		ip := addr.String()
		if ip != "127.0.0.1" {
			return ip
		}
	}
	return ""
}

func NetHostname() string {
	name, _ := os.Hostname()
	return name
}
