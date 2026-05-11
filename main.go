package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	router := newRouter()
	addr := envString("LISTEN_ADDR", ":8466")
	log.Printf("listening on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("serve: %v", err)
	}
}

func envString(name string, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
