package utils

import (
	"net"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func GetInterface(name string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Name != name {
			continue
		}
		return &iface, nil
	}

	return nil, InterfaceNotFoundErr
}

func GetIpWithPrefix(ifacename, prefix string) (string, error) {
	iface, err := GetInterface(ifacename)
	if err != nil {
		return "", err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, isIpNet := addr.(*net.IPNet)
		if !isIpNet {
			continue
		}

		ip := ipNet.IP
		if !strings.HasPrefix(ip.String(), prefix) {
			continue
		}

		log.Info("ip get", zap.String("iface", iface.Name), zap.String("ip", ip.String()))
		return ip.String(), nil
	}

	return "", IPNotFoundErr
}
