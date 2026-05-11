package clientinfo

import (
	"context"
	"net"
	"net/http"
	"net/netip"
	"strings"
)

type contextKey struct{}

type Detector struct {
	pseudoIPv4CIDR netip.Prefix
}

type Info struct {
	IP string `json:"ip"`
}

func NewDetector() *Detector {
	return &Detector{
		pseudoIPv4CIDR: netip.MustParsePrefix("240.0.0.0/4"),
	}
}

func (d *Detector) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		info := d.Detect(r)
		ctx := context.WithValue(r.Context(), contextKey{}, info)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func FromContext(ctx context.Context) (Info, bool) {
	info, ok := ctx.Value(contextKey{}).(Info)
	return info, ok
}

func (d *Detector) Detect(r *http.Request) Info {
	info := Info{}

	remoteIP, _ := parseRemoteIP(r.RemoteAddr)

	if clientIP, ok := d.bestEffortClientIP(r.Header); ok {
		info.IP = clientIP.String()
		return info
	}

	if remoteIP.IsValid() {
		info.IP = remoteIP.String()
	}

	return info
}

func (d *Detector) bestEffortClientIP(headers http.Header) (netip.Addr, bool) {
	connectingIP := headerValue(headers, "cf-connecting-ip")
	connectingIPv6 := headerValue(headers, "cf-connecting-ipv6")
	if clientIP, ok := d.cloudflareClientIP(connectingIP, connectingIPv6); ok {
		return clientIP, true
	}

	if clientIP, ok := parseXForwardedFor(headerValue(headers, "x-forwarded-for")); ok {
		return clientIP, true
	}

	if clientIP, ok := parseIP(headerValue(headers, "x-real-ip")); ok {
		return clientIP, true
	}

	return netip.Addr{}, false
}

func (d *Detector) cloudflareClientIP(connectingIP string, connectingIPv6 string) (netip.Addr, bool) {
	primary, ok := parseIP(connectingIP)
	if !ok {
		return netip.Addr{}, false
	}

	if primary.Is4() && d.pseudoIPv4CIDR.Contains(primary) {
		if realIPv6, ok := parseIP(connectingIPv6); ok && realIPv6.Is6() {
			return realIPv6, true
		}
	}

	return primary, true
}

func parseRemoteIP(remoteAddr string) (netip.Addr, string) {
	if remoteAddr == "" {
		return netip.Addr{}, ""
	}

	if host, _, err := net.SplitHostPort(remoteAddr); err == nil {
		if ip, ok := parseIP(host); ok {
			return ip, ip.String()
		}
	}

	if ip, ok := parseIP(remoteAddr); ok {
		return ip, ip.String()
	}

	return netip.Addr{}, remoteAddr
}

func parseXForwardedFor(value string) (netip.Addr, bool) {
	for _, part := range strings.Split(value, ",") {
		if ip, ok := parseIP(part); ok {
			return ip, true
		}
	}

	return netip.Addr{}, false
}

func parseIP(value string) (netip.Addr, bool) {
	if value == "" {
		return netip.Addr{}, false
	}

	addr, err := netip.ParseAddr(strings.TrimSpace(value))
	if err != nil {
		return netip.Addr{}, false
	}

	return addr.Unmap(), true
}

func headerValue(headers http.Header, name string) string {
	return strings.TrimSpace(headers.Get(name))
}
