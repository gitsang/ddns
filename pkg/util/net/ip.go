package net

import (
	"errors"
	"net"
	"strings"
)

var (
	InterfaceNotFoundErr = errors.New("interface not found")
	IPNotFoundErr        = errors.New("ip not found")
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

		return ip.String(), nil
	}

	return "", IPNotFoundErr
}

func GetIpsWithPrefix(ifacename, prefix string) ([]string, error) {
	var ips []string = make([]string, 0)

	iface, err := GetInterface(ifacename)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
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

		ips = append(ips, ip.String())
	}

	if len(ips) == 0 {
		return nil, IPNotFoundErr
	}

	return ips, nil
}
