package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shahinzaman102/Go_JumpStart/internal/data"
	"github.com/shahinzaman102/Go_JumpStart/internal/models"

	"github.com/go-chi/chi/v5"
)

// GetAllAlbums responds with all albums in JSON format.
func GetAllAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := data.AllAlbums()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(albums)
}

// GetAlbumsByArtist responds with albums filtered by artist name.
func GetAlbumsByArtist(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimSpace(chi.URLParam(r, "name"))
	if name == "" {
		http.Error(w, "Artist name is required", http.StatusBadRequest)
		return
	}

	albums, err := data.AlbumsByArtist(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(albums)
}

// GetAlbumByID responds with a single album by its ID.
func GetAlbumByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	album, err := data.AlbumByID(id)
	if err != nil {
		http.Error(w, "Album not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(album)
}

// CreateAlbum handles adding a new album to the database.
func CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var album models.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Trim whitespace
	album.Title = strings.TrimSpace(album.Title)
	album.Artist = strings.TrimSpace(album.Artist)

	// Validation
	if album.Title == "" || album.Artist == "" {
		http.Error(w, "Title and Artist are required", http.StatusBadRequest)
		return
	}
	if len(album.Title) > 200 {
		http.Error(w, "Title must be 1-200 characters", http.StatusBadRequest)
		return
	}
	if len(album.Artist) > 100 {
		http.Error(w, "Artist must be 1-100 characters", http.StatusBadRequest)
		return
	}

	id, err := data.AddAlbum(album)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	album.ID = id
	json.NewEncoder(w).Encode(album)
}

// CanPurchaseAlbum checks if the requested quantity can be purchased.
func CanPurchaseAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	qtyStr := r.URL.Query().Get("qty")
	if qtyStr == "" {
		http.Error(w, "Quantity parameter 'qty' is required", http.StatusBadRequest)
		return
	}

	qty, err := strconv.ParseInt(qtyStr, 10, 64)
	if err != nil || qty <= 0 {
		http.Error(w, "Quantity must be a positive integer", http.StatusBadRequest)
		return
	}

	ok, err := data.CanPurchase(id, qty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"canPurchase": ok})
}

// GetOrdersByUser serves the last 10 orders for a logged-in user (HTML page).
func GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	auth, authOk := session.Values["authenticated"].(bool)
	userID, idOk := session.Values["user_id"].(int64)

	if !authOk || !auth || !idOk {
		session.Values["redirect_after_login"] = "/orders"
		session.Save(r, w)
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/html")
		tmpl := template.Must(template.ParseFiles("templates/login_required.html"))
		tmpl.Execute(w, map[string]string{"Resource": "Orders API"})
		return
	}

	cacheKey := fmt.Sprintf("orders:user:%d:last10", userID)
	var orders []models.GetOrder

	if err := data.GetOrdersCache(cacheKey, &orders); err != nil {
		log.Printf("Cache MISS for user: %d", userID)
		log.Printf("[SIMULATION] Sleeping 2s to simulate slow DB query for user %d...", userID)

		orders, err = data.GetOrdersByUser(userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("Returning data from DB...")
		_ = data.SetOrdersCache(cacheKey, orders)

	} else {
		log.Printf("Cache HIT for user: %d", userID)
		log.Println("Returning data from cache...")
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("templates/orders.html"))
	tmpl.Execute(w, map[string]any{
		"Orders": orders,
	})
}

// CreateOrderByUser handles creating a new order for the logged-in user.
func CreateOrderByUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userID, ok := session.Values["user_id"].(int64)
	auth, authOk := session.Values["authenticated"].(bool)

	if !ok || !authOk || !auth {
		session.Values["redirect_after_login"] = "/orders"
		session.Save(r, w)
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Unauthorized: You must log in first to create an order.\n" +
			"Visit /login to log in and then retry your request."))
		return
	}

	var order models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validation
	if order.AlbumID <= 0 {
		http.Error(w, "AlbumID must be a positive integer", http.StatusBadRequest)
		return
	}
	if order.Quantity <= 0 {
		http.Error(w, "Quantity must be a positive integer", http.StatusBadRequest)
		return
	}

	// Always use session user
	order.Customer = userID

	// Insert into DB
	id, err := data.CreateOrderByUser(r.Context(), order.AlbumID, order.Quantity, order.Customer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update per-user cache
	cacheKey := fmt.Sprintf("orders:user:%d:last10", order.Customer)
	var cached []models.GetOrder
	if err := data.GetOrdersCache(cacheKey, &cached); err == nil {
		newOrder := models.GetOrder{
			ID:       id,
			AlbumID:  order.AlbumID,
			Customer: order.Customer,
			Quantity: order.Quantity,
			Date:     time.Now(),
		}
		cached = append([]models.GetOrder{newOrder}, cached...)
		if len(cached) > 10 {
			cached = cached[:10]
		}
		_ = data.SetOrdersCache(cacheKey, cached)
		log.Printf("[CACHE UPDATE] User %d cache updated with new order %d", newOrder.Customer, newOrder.ID)

	} else {
		_ = data.OrderCache.Delete(cacheKey)
		log.Printf("[CACHE INVALIDATE] User %d cache cleared (will refresh on next GetOrders)", order.Customer)
	}

	// Respond with JSON
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"order_id": id,
		"message":  "Order created successfully",
	})
}

// GetCustomerName returns the full name of a customer by ID.
func GetCustomerName(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	name, err := data.GetCustomerName(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"name": name})
}

// HandleMultipleResultSets demonstrates fetching multiple result sets (albums + customers).
func HandleMultipleResultSets(w http.ResponseWriter, r *http.Request) {
	result, err := data.GetAlbumsAndCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

// QueryWithTimeout executes a DB query with a timeout context.
func QueryWithTimeout(w http.ResponseWriter, r *http.Request) {
	albums, err := data.QueryAlbumsWithTimeout(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
		return
	}
	json.NewEncoder(w).Encode(albums)
}
