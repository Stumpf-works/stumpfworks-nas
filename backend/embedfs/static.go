// Package embedfs provides embedded frontend static files
package embedfs

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/Stumpf-works/stumpfworks-nas/pkg/logger"
	"go.uber.org/zap"
)

// Embed the entire dist directory
// The Makefile copies frontend/dist to backend/embedfs/dist before building
//
//go:embed dist
var frontendFS embed.FS

// GetFileSystem returns the embedded filesystem
// This strips the "dist" prefix so files are served from root
func GetFileSystem() (http.FileSystem, error) {
	// Get subdirectory to strip the "dist" prefix
	fsys, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(fsys), nil
}

// SPAHandler wraps the http.FileServer to handle Single Page Application routing
// Falls back to index.html for client-side routes
type SPAHandler struct {
	fileServer http.Handler
	indexFile  []byte
}

// NewSPAHandler creates a new SPA handler
func NewSPAHandler() (*SPAHandler, error) {
	fsys, err := GetFileSystem()
	if err != nil {
		logger.Error("Failed to get embedded filesystem", zap.Error(err))
		return nil, err
	}

	// Read index.html for fallback
	indexFile, err := fs.ReadFile(frontendFS, "dist/index.html")
	if err != nil {
		logger.Warn("index.html not found in embedded files, SPA routing may not work", zap.Error(err))
		indexFile = []byte("<!DOCTYPE html><html><body><h1>StumpfWorks NAS</h1><p>Frontend not built. Run 'make build' to include the web interface.</p></body></html>")
	}

	return &SPAHandler{
		fileServer: http.FileServer(fsys),
		indexFile:  indexFile,
	}, nil
}

// ServeHTTP implements http.Handler
func (h *SPAHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the cleaned path
	urlPath := r.URL.Path

	// Don't serve static files for API routes or WebSocket
	if strings.HasPrefix(urlPath, "/api/") ||
	   strings.HasPrefix(urlPath, "/ws") ||
	   strings.HasPrefix(urlPath, "/health") ||
	   strings.HasPrefix(urlPath, "/metrics") {
		http.NotFound(w, r)
		return
	}

	// Clean the path
	urlPath = path.Clean(urlPath)

	// Check if the file exists
	_, err := frontendFS.Open("dist" + urlPath)

	// If file exists, serve it
	if err == nil {
		h.fileServer.ServeHTTP(w, r)
		return
	}

	// If it's a file request (has extension) but doesn't exist, 404
	if strings.Contains(path.Base(urlPath), ".") {
		http.NotFound(w, r)
		return
	}

	// Otherwise, serve index.html for client-side routing
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(h.indexFile)
}
