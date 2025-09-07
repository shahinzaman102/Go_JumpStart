package handlers

import (
	"html/template"
	"net/http"

	assets "github.com/shahinzaman102/Go_JumpStart"
	"github.com/shahinzaman102/Go_JumpStart/internal/data"

	"github.com/gorilla/sessions"
)

var (
	store    *sessions.CookieStore
	authRepo *data.AuthRepo
)

// Init initializes the session store and auth repository.
func Init(s *sessions.CookieStore, repo *data.AuthRepo) {
	store = s
	authRepo = repo
}

// Login handles user login: verifies credentials and sets session values.
func Login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	username := r.FormValue("username")
	password := r.FormValue("password")

	ok, err := authRepo.VerifyUser(username, password)
	if err != nil || !ok {
		tmpl := template.Must(template.ParseFS(assets.Templates, "templates/unauthorized.html"))
		w.WriteHeader(http.StatusUnauthorized)
		tmpl.Execute(w, nil)
		return
	}

	userID, err := authRepo.GetUserID(username)
	if err != nil {
		http.Error(w, "Failed to fetch user ID", http.StatusInternalServerError)
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = userID

	// Determine redirect path
	redirectPath := r.FormValue("redirect")
	if redirectPath == "" {
		redirectPath = r.URL.Query().Get("redirect")
	}
	if redirectPath == "" {
		if v, ok := session.Values["redirect_after_login"].(string); ok && v != "" {
			redirectPath = v
			delete(session.Values, "redirect_after_login")
		}
	}
	if redirectPath == "" {
		redirectPath = "/dashboard"
	}

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

// LoginForm renders the login page with an optional redirect.
func LoginForm(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	tmpl := template.Must(template.ParseFS(assets.Templates, "templates/login.html"))
	tmpl.Execute(w, map[string]string{
		"Redirect": redirect,
	})
}

// Logout clears the session and redirects to the home page.
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	// Clear session completely
	session.Options.MaxAge = -1
	session.Values = make(map[any]any)

	if err := session.Save(r, w); err != nil {
		http.Error(w, "Failed to clear session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
