package tp

import (
	"net"
	"net/http"
	"strings"
)

func GetClientIp(req *http.Request) string {
	// get ip from header X-Forwarded-For
	ip := req.Header.Get("X-Forwarded-For")
	if ip != "" {
		ip = strings.Split(ip, ",")[0]
		return strings.TrimSpace(ip)
	}

	// get ip from X-Real-IP
	ip = req.Header.Get("X-Real-IP")
	if ip != "" {
		return strings.TrimSpace(ip)
	}

	// get ip from RemoteAddr
	ip, _, _ = net.SplitHostPort(strings.TrimSpace(req.RemoteAddr))
	return ip
}
