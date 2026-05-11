package main

import (
	"context"
	"net/http"

	"visitor-ip-api-go/internal/clientinfo"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

type serverConfig struct {
}

type ipOutput struct {
	Body ipResponse
}

type ipResponse struct {
	IP string `json:"ip" example:"2606:4700:4700::1111" doc:"Best-known client public IP address."`
}

func newRouter() http.Handler {
	detector := clientinfo.NewDetector()
	router := chi.NewRouter()
	router.Use(corsMiddleware)
	router.Use(detector.Middleware)

	apiConfig := huma.DefaultConfig("visitor-ip-api-go", "1.0.0")
	apiConfig.Info.Description = "Return the best-known client IP address."
	apiConfig.DocsRenderer = huma.DocsRendererSwaggerUI
	apiConfig.CreateHooks = nil

	api := humachi.New(router, apiConfig)

	huma.Get(api, "/ip", func(ctx context.Context, input *struct{}) (*ipOutput, error) {
		info, _ := clientinfo.FromContext(ctx)
		return &ipOutput{
			Body: ipResponse{
				IP: info.IP,
			},
		}, nil
	}, huma.OperationTags("IP"))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusTemporaryRedirect)
	})

	return router
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
