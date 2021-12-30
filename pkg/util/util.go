package util

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetRealIPAddress(r *http.Request) (string, error) {
	ip := r.Header.Get("X-REAL-IP")
	parseIP := net.ParseIP(ip)
	if parseIP != nil {
		return ip, nil
	}

	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		parseIP = net.ParseIP(ip)
		if parseIP != nil {
			return ip, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	parseIP = net.ParseIP(ip)
	if parseIP != nil {
		return ip, nil
	}

	return "", fmt.Errorf("can not get real IP address")
}

func ValidRequest(r *http.Request) bool {
	if r.Header.Get("Req-Source") == "AutoK3s" {
		return true
	}
	return false
}
