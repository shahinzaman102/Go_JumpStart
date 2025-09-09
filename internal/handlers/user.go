package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shahinzaman102/Go_JumpStart/internal/data"
	"github.com/shahinzaman102/Go_JumpStart/internal/models"

	"github.com/go-chi/chi/v5"
)

// UserResponse defines the JSON output for API clients (hides password)
type UserResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// mapUser converts internal User model to UserResponse
func mapUser(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

// GetUsers returns a list of all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := data.GetAllUsers()
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}

	resp := make([]*UserResponse, len(users))
	for i, u := range users {
		userResp := mapUser(u)
		resp[i] = &userResp
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetUserByID returns a single user by ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := data.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapUser(*user))
}

// CreateUser adds a new user to the database
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Trim whitespace
	input.Username = strings.TrimSpace(input.Username) // removes extra spaces
	input.Password = strings.TrimSpace(input.Password)

	// Validation
	if input.Username == "" || input.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}
	if len(input.Username) < 3 || len(input.Username) > 50 {
		http.Error(w, "Username must be between 3 and 50 characters", http.StatusBadRequest)
		return
	}
	if len(input.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	id, err := data.CreateUser(input.Username, input.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user, _ := data.GetUserByID(int(id))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "User created successfully",
		"user":    mapUser(*user),
	})
}

// UpdateUser updates username and/or password for a given user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Trim whitespace
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	input.Username = strings.TrimSpace(input.Username)
	input.Password = strings.TrimSpace(input.Password)

	if input.Username == "" && input.Password == "" {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}
	if input.Username != "" && (len(input.Username) < 3 || len(input.Username) > 50) {
		http.Error(w, "Username must be between 3 and 50 characters", http.StatusBadRequest)
		return
	}
	if input.Password != "" && len(input.Password) < 6 {
		http.Error(w, "Password must be at least 6 characters", http.StatusBadRequest)
		return
	}

	if err := data.UpdateUserByID(id, input.Username, input.Password); err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	updatedUser, err := data.GetUserByID(id)
	if err != nil {
		http.Error(w, "Error fetching updated user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"user":   mapUser(*updatedUser),
	})
}

// DeleteUser removes a user by ID
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if err := data.DeleteUserByID(id); err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": "deleted",
		"id":     id,
	})
}
