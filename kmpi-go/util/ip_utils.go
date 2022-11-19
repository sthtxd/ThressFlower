package util

import (
	"kmpi-go/log"
	"net"
)

func GetIp() (*string, error) {
	var ipResult string
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Error("Interfaces error", err.Error())
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			log.Error("Addrs error", err.Error())
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ipResult = ip.String()
			//	fmt.Println("ip: ", ip.String(), "mac: ", iface.HardwareAddr.String())
		}
	}
	return &ipResult, nil
}
