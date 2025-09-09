package handlers

import (
	"html/template"
	"log"
	"net/http"

	assets "github.com/shahinzaman102/Go_JumpStart"
	"github.com/shahinzaman102/Go_JumpStart/internal/models"
)

func Form(w http.ResponseWriter, r *http.Request) {
	// Load form.html from embedded templates
	tmpl := template.Must(template.ParseFS(assets.Templates, "templates/form.html"))

	if r.Method == http.MethodPost {
		// Parse submitted form fields
		r.ParseForm()
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		log.Println("Form received:", email, subject, message)

		// Render template with submitted data
		data := models.FormResponse{
			Success: true,
			Email:   email,
			Subject: subject,
			Message: message,
		}
		tmpl.Execute(w, data)
		return
	}

	// Render empty form on GET
	tmpl.Execute(w, nil)
}
