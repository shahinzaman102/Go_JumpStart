package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/shahinzaman102/Go_JumpStart/internal/models"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	// Dashboard handler → shows tasks if user logged in, else 401 + login page.
	if !isAuthenticated(r) {
		session, _ := store.Get(r, "session")
		session.Values["redirect_after_login"] = "/dashboard"
		session.Save(r, w)

		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/login_required.html"))
		tmpl.Execute(w, map[string]string{"Resource": "dashboard"})
		return
	}

	// Example due dates
	dueBuildApp := time.Date(2025, 8, 20, 0, 0, 0, 0, time.UTC)
	dueLearnGo := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	todos := []models.Todo{
		{Title: "Learn Go", Done: true, Progress: 100, Due: &dueLearnGo},
		{Title: "Build the 2nd App", Done: false, Progress: 40, Due: &dueBuildApp},
	}

	tmpl := template.Must(template.ParseFiles("templates/dashboard.html"))
	tmpl.Execute(w, todos)
}

func isAuthenticated(r *http.Request) bool {
	session, _ := store.Get(r, "session")
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}
