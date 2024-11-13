package tp

import (
	"net/http"
	"strings"
)

func GetUserIp(req *http.Request) string {
	// get ip from header X-Forwarded-For
	ip := req.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For 可能包含多个 IP 地址，取第一个
		ip = strings.Split(ip, ",")[0]
		return strings.TrimSpace(ip)
	}

	// get ip from X-Real-IP
	ip = req.Header.Get("X-Real-IP")
	if ip != "" {
		return strings.TrimSpace(ip)
	}

	// get ip from RemoteAddr
	ip = req.RemoteAddr
	if ip != "" {
		// filter out the port
		if colon := strings.LastIndex(ip, ":"); colon != -1 {
			ip = ip[:colon]
		}
		return ip
	}

	return ""
}
