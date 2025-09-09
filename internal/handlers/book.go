package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/shahinzaman102/Go_JumpStart/internal/data"
	"github.com/shahinzaman102/Go_JumpStart/internal/models"

	"github.com/go-chi/chi/v5"
)

// GetBooks returns all books in JSON format.
func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.GetAllBooks())
}

// GetBookByID returns a single book by ID.
func GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"message": "invalid book ID"}`, http.StatusBadRequest)
		return
	}

	book, err := data.GetBookByID(id)
	if err != nil {
		http.Error(w, `{"message": "book not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// PostBook creates a new book.
func PostBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, `{"message": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Trim whitespace
	newBook.Title = strings.TrimSpace(newBook.Title)
	newBook.Author = strings.TrimSpace(newBook.Author)

	// Validation
	if newBook.Title == "" || newBook.Author == "" {
		http.Error(w, `{"message": "Title and Author are required"}`, http.StatusBadRequest)
		return
	}
	if len(newBook.Title) > 200 {
		http.Error(w, `{"message": "Title must be 1-200 characters"}`, http.StatusBadRequest)
		return
	}
	if len(newBook.Author) > 100 {
		http.Error(w, `{"message": "Author must be 1-100 characters"}`, http.StatusBadRequest)
		return
	}

	book := data.AddBook(newBook)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook updates an existing book by ID.
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"message": "invalid book ID"}`, http.StatusBadRequest)
		return
	}

	var updatedData models.Book
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, `{"message": "invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Trim whitespace
	updatedData.Title = strings.TrimSpace(updatedData.Title)
	updatedData.Author = strings.TrimSpace(updatedData.Author)

	// Validate at least one field
	if updatedData.Title == "" && updatedData.Author == "" {
		http.Error(w, `{"message": "No fields to update"}`, http.StatusBadRequest)
		return
	}
	// Validate lengths if provided
	if updatedData.Title != "" && len(updatedData.Title) > 200 {
		http.Error(w, `{"message": "Title must be 1-200 characters"}`, http.StatusBadRequest)
		return
	}
	if updatedData.Author != "" && len(updatedData.Author) > 100 {
		http.Error(w, `{"message": "Author must be 1-100 characters"}`, http.StatusBadRequest)
		return
	}

	updatedBook, err := data.UpdateBook(id, updatedData)
	if err != nil {
		http.Error(w, `{"message": "book not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"status":  "success",
		"message": "book updated successfully",
		"book":    updatedBook,
	})
}

// DeleteBook removes a book by ID.
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"message": "invalid book ID"}`, http.StatusBadRequest)
		return
	}

	if err := data.DeleteBook(id); err != nil {
		http.Error(w, `{"message": "book not found"}`, http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "book deleted successfully",
	})
}
