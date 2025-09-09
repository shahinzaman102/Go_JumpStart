package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// TestUI serves the test_ui.html template from the filesystem.
func TestUI(w http.ResponseWriter, r *http.Request) {
	// Parse template from the templates/ folder on disk
	tmpl, err := template.ParseFiles(filepath.Join("templates", "test_ui.html"))
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Render the template and write HTML to the response
	tmpl.ExecuteTemplate(w, "test_ui.html", nil)
}
