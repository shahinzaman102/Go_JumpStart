package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shahinzaman102/Go_JumpStart/internal/data"
	"github.com/shahinzaman102/Go_JumpStart/internal/models"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

// local test DB helper (self-contained for handlers tests)
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:") // in-memory DB
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE album (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT, artist TEXT, price REAL, quantity INTEGER
	)`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`INSERT INTO album (title, artist, price, quantity)
		VALUES ("Go Beats", "Gopher", 9.99, 5)`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func setupAlbumHandlerDB(t *testing.T) {
	db := setupTestDB(t)
	data.InitDBConnection(db)
}

func TestGetAlbumByID(t *testing.T) {
	setupAlbumHandlerDB(t)

	// create a chi router and mount the handler
	r := chi.NewRouter()
	r.Get("/albums/{id}", GetAlbumByID)

	// request for /albums/1
	req := httptest.NewRequest(http.MethodGet, "/albums/1", nil)
	w := httptest.NewRecorder()

	// pass through chi router (so URLParam works)
	r.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var alb models.Album
	if err := json.NewDecoder(resp.Body).Decode(&alb); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if strings.ToLower(alb.Artist) != "gopher" {
		t.Errorf("expected artist Gopher, got %s", alb.Artist)
	}
}
