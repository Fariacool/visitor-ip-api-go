package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIPEndpoint(t *testing.T) {
	router := newRouter()
	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.RemoteAddr = "172.64.1.10:443"
	req.Header.Set("CF-Connecting-IP", "198.51.100.23")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	if recorder.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatalf("missing CORS header")
	}

	var raw map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &raw); err != nil {
		t.Fatalf("json.Unmarshal(raw) error = %v", err)
	}

	if _, ok := raw["$schema"]; ok {
		t.Fatalf("unexpected $schema in response")
	}

	if _, ok := raw["checked_at"]; ok {
		t.Fatalf("unexpected checked_at in response")
	}

	var payload struct {
		IP string `json:"ip"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if payload.IP != "198.51.100.23" {
		t.Fatalf("payload.IP = %q, want %q", payload.IP, "198.51.100.23")
	}

	if len(raw) != 1 {
		t.Fatalf("response field count = %d, want %d; body = %s", len(raw), 1, recorder.Body.String())
	}
}
