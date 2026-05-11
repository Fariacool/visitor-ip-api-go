package clientinfo

import (
	"net/http/httptest"
	"testing"
)

func TestDetectDirectIPv4(t *testing.T) {
	detector := NewDetector()
	req := httptest.NewRequest("GET", "http://example.com/ip", nil)
	req.RemoteAddr = "203.0.113.10:54321"

	info := detector.Detect(req)

	if info.IP != "203.0.113.10" {
		t.Fatalf("info.IP = %q, want %q", info.IP, "203.0.113.10")
	}
}

func TestDetectCloudflareIPv4(t *testing.T) {
	detector := NewDetector()
	req := httptest.NewRequest("GET", "http://example.com/ip", nil)
	req.RemoteAddr = "172.64.1.10:443"
	req.Header.Set("CF-Connecting-IP", "198.51.100.23")

	info := detector.Detect(req)

	if info.IP != "198.51.100.23" {
		t.Fatalf("info.IP = %q, want %q", info.IP, "198.51.100.23")
	}
}

func TestDetectCloudflarePseudoIPv4Overwrite(t *testing.T) {
	detector := NewDetector()
	req := httptest.NewRequest("GET", "http://example.com/ip", nil)
	req.RemoteAddr = "172.64.1.10:443"
	req.Header.Set("CF-Connecting-IP", "240.16.0.1")
	req.Header.Set("CF-Connecting-IPv6", "2606:4700:4700::1111")

	info := detector.Detect(req)

	if info.IP != "2606:4700:4700::1111" {
		t.Fatalf("info.IP = %q, want %q", info.IP, "2606:4700:4700::1111")
	}
}

func TestDetectXForwardedFor(t *testing.T) {
	detector := NewDetector()

	req := httptest.NewRequest("GET", "http://example.com/ip", nil)
	req.RemoteAddr = "10.0.0.5:54321"
	req.Header.Set("X-Forwarded-For", "198.51.100.23, 10.0.0.5")

	info := detector.Detect(req)

	if info.IP != "198.51.100.23" {
		t.Fatalf("info.IP = %q, want %q", info.IP, "198.51.100.23")
	}
}

func TestDetectXRealIP(t *testing.T) {
	detector := NewDetector()

	req := httptest.NewRequest("GET", "http://example.com/ip", nil)
	req.RemoteAddr = "10.0.0.5:54321"
	req.Header.Set("X-Real-IP", "198.51.100.99")

	info := detector.Detect(req)

	if info.IP != "198.51.100.99" {
		t.Fatalf("info.IP = %q, want %q", info.IP, "198.51.100.99")
	}
}

func TestDetectFallsBackToRemoteAddr(t *testing.T) {
	detector := NewDetector()

	req := httptest.NewRequest("GET", "http://example.com/ip", nil)
	req.RemoteAddr = "203.0.113.10:54321"
	req.Header.Set("X-Forwarded-For", "unknown")
	req.Header.Set("X-Real-IP", "also-bad")

	info := detector.Detect(req)

	if info.IP != "203.0.113.10" {
		t.Fatalf("info.IP = %q, want %q", info.IP, "203.0.113.10")
	}
}
