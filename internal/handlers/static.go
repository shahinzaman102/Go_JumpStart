package handlers

import (
	"io/fs"    // Provides filesystem abstraction, used to wrap embed.FS for http.FileServer
	"net/http" // HTTP server and FileServer

	assets "github.com/shahinzaman102/Go_JumpStart"

	"github.com/go-chi/chi/v5"
)

// ServeStatic registers a handler to serve embedded static files at /static/*.
// Example: /static/style.css serves embedded static/style.css
func ServeStatic(r chi.Router) {
	// Wrap embedded static folder as a normal http.FileSystem
	staticFiles, _ := fs.Sub(assets.StaticFS, "static")
	fileServer := http.FileServer(http.FS(staticFiles))

	// Register route and strip /static prefix before serving files
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))
}
