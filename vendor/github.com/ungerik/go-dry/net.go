package dry

import (
	"fmt"
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

// RealNetIP returns the real local IP of the system or an empty string.
func RealNetIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	// get real local IP
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return ""
}

func NetHostname() string {
	name, _ := os.Hostname()
	return name
}
