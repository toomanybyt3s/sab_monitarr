package logger

import (
	"log"
	"net"
	"net/http"
	"strings"
)

// Log writes a levelled log line. DEBUG lines are suppressed unless debug
// mode is enabled. When logClientInfo is true and a request is provided,
// the client IP, User-Agent and URI are appended on a second line.
func Log(debug bool, level, message string, r *http.Request, logClientInfo bool) {
	if !debug && level == "DEBUG" {
		return
	}

	log.Printf("[%s] %s", level, message)

	if logClientInfo && r != nil {
		clientIP := GetClientIP(r)
		log.Printf("  Client: %s, User-Agent: %s, URI: %s",
			clientIP, r.UserAgent(), r.RequestURI)
	}
}

// GetClientIP returns the client's IP address from the request.
// It uses the last (outermost) entry in X-Forwarded-For, which is set by
// the nearest trusted proxy and cannot be spoofed by the client.
// net.SplitHostPort is used so IPv6 addresses are handled correctly.
func GetClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[len(parts)-1])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// Middleware logs every incoming request before passing it to the next handler.
func Middleware(next http.Handler, debug bool, logClientInfo bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log(debug, "INFO", r.Method+" "+r.URL.Path, r, logClientInfo)
		next.ServeHTTP(w, r)
	})
}
