package logger_test

import (
	"net/http/httptest"
	"testing"

	"github.com/toomanybyt3s/sab_monitarr/internal/logger"
)

func TestGetClientIP(t *testing.T) {
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.RemoteAddr = "192.168.1.100:12345"
	if ip := logger.GetClientIP(req1); ip != "192.168.1.100" {
		t.Errorf("expected 192.168.1.100, got %s", ip)
	}

	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "10.0.0.1:12345"
	req2.Header.Set("X-Forwarded-For", "1.2.3.4, 10.0.0.1")
	if ip := logger.GetClientIP(req2); ip != "10.0.0.1" {
		t.Errorf("expected 10.0.0.1 (last XFF entry), got %s", ip)
	}

	req3 := httptest.NewRequest("GET", "/", nil)
	req3.Header.Set("X-Forwarded-For", "192.168.1.100")
	if ip := logger.GetClientIP(req3); ip != "192.168.1.100" {
		t.Errorf("expected 192.168.1.100 from single XFF, got %s", ip)
	}

	req4 := httptest.NewRequest("GET", "/", nil)
	req4.RemoteAddr = "[::1]:54321"
	if ip := logger.GetClientIP(req4); ip != "::1" {
		t.Errorf("expected ::1 from IPv6 RemoteAddr, got %s", ip)
	}
}
