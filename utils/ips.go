package utils

import "net"

func GetLocalIps() ([]string, error) {
	ips := make([]string, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ips = append(ips, ip.String())
		}
	}
	return ips, nil
}

func IsLocalIp(ip string) (bool, error) {
	isLocal := false
	ips, err := GetLocalIps()
	if err != nil {
		return isLocal, err
	}
	for _, item := range ips {
		if item == ip {
			return true, nil
		}
	}
	return isLocal, nil
}
