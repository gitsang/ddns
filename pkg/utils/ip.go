package utils

import (
	"errors"
	"net"
	"strings"

	log "github.com/gitsang/golog"
	"go.uber.org/zap"
)

func GetIp(ipv6, w bool, ifaceName string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		if iface.Name != ifaceName {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, isIpNet := addr.(*net.IPNet)
			if !isIpNet {
				continue
			}

			ip := ipNet.IP
			if ipv6 && ip.To4() != nil {
				continue
			}

			if w && (strings.HasPrefix(ip.String(), "fe80:") ||
				strings.HasPrefix(ip.String(), "192.168") ||
				strings.HasPrefix(ip.String(), "10.")) {
				continue
			}

			log.Info("ip get", zap.Reflect("iface", iface.Name), zap.String("ip", ip.String()))
			return ip.String(), nil
		}
	}

	return "", errors.New("ip not found")
}
