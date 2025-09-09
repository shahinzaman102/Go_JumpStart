package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/shahinzaman102/Go_JumpStart/internal/models"
)

// JsonEncode: Take JSON -> struct -> return JSON
func JsonEncode(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if u.Username == "" || u.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	u.CreatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// JsonDecode: Take JSON -> struct -> return text summary
func JsonDecode(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if u.Username == "" || u.Password == "" {
		http.Error(w, "Username and Password are required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Received user: %+v", u)
}
