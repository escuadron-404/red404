//go:build dev

package routes

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

// Proxy Vite in dev
func RegisterFrontendHandlers(mux *http.ServeMux) {
	log.Println("Go Backend: Running in DEVELOPMENT mode. Proxying frontend requests to Vite dev server.")

	viteDevServerHost := os.Getenv("DEV_SERVER")
	if viteDevServerHost == "" {
		log.Fatal("Go Backend: VITE_DEV_SERVER_HOST not set in development mode. Cannot proxy to Vite.")
	}

	viteURL, err := url.Parse("http://" + viteDevServerHost)
	if err != nil {
		log.Fatalf("Go Backend: Failed to parse VITE_DEV_SERVER_HOST URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(viteURL)

	// Optional: Add a custom director for debugging/logging proxy requests
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		log.Printf("Go Backend: Proxying request: %s %s to Vite dev server at %s", req.Method, req.URL.Path, viteURL.Host)
	}

	mux.Handle("/", proxy) // This catches everything not caught by previous Handlers
}
