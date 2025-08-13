//go:build !dev

package routes

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// spaHandler implements the http.Handler interface, so we can use it
// to serve a single-page application from disk.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP for production static file serving
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Prepend the static directory to the requested path
	path := filepath.Join(h.staticPath, r.URL.Path)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// File does not exist, serve index.html
		log.Printf("Go Backend: File not found: %s. Serving index.html for %s", path, r.URL.Path)
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		log.Printf("Go Backend: Error checking file existence for %s: %v", path, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// File exists, use a standard FileServer to serve it
	log.Printf("Go Backend: Serving static file: %s", r.URL.Path)
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// RegisterFrontendHandlers configures the HTTP server for production mode.
// It serves static files from the specified frontendDistPath.
func RegisterFrontendHandlers(mux *http.ServeMux) {
	log.Println("Go Backend: Running in PRODUCTION mode. Serving static files from ./frontend/dist.")

	// The path inside the Docker container where frontend/dist will be mounted/copied
	const frontendDistPath = "./frontend/dist"
	spa := spaHandler{staticPath: frontendDistPath, indexPath: "index.html"}
	mux.Handle("/", spa)
}
