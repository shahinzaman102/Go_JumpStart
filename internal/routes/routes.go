package routes

import (
	"net/http"

	"github.com/shahinzaman102/Go_JumpStart/internal/handlers"
	"github.com/shahinzaman102/Go_JumpStart/internal/middleware"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Register() http.Handler {
	r := chi.NewRouter()
	// Creates a new Chi router (Chi = lightweight HTTP router for Go).
	// This router decides what function runs when a request hits a certain URL.

	// --- Middleware ---
	r.Use(chimiddleware.RequestID) // RequestID → trace requests.
	r.Use(chimiddleware.RealIP)    // RealIP → know who’s really calling.
	r.Use(chimiddleware.Logger)    // Logger → track activity.
	r.Use(chimiddleware.Recoverer) // Recoverer → server never dies on error.
	r.Use(middleware.Tracing)

	// --- CORS ---
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // front-end URLs
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by browsers
	}))

	// HTTP routers -->

	// --- App Home ---
	r.Get("/", handlers.TestUI)

	// --- Form ---
	r.Route("/form", func(r chi.Router) {
		r.Get("/", handlers.Form)
		r.Post("/", handlers.Form)
	})

	// --- Static Files ---
	handlers.ServeStatic(r)

	// --- Authentication Flow ---
	r.Group(func(r chi.Router) {
		r.Get("/login", handlers.LoginForm)
		r.Post("/login", handlers.Login)
		r.Get("/logout", handlers.Logout)
	})

	// --- Dashboard ---
	r.Get("/dashboard", handlers.Dashboard)

	// --- Users API ---
	r.Route("/users", func(r chi.Router) {
		r.Get("/", handlers.GetUsers)
		r.Post("/", handlers.CreateUser)
		r.Get("/{id}", handlers.GetUserByID)
		r.Put("/{id}", handlers.UpdateUser)
		r.Delete("/{id}", handlers.DeleteUser)
	})

	// --- Books API ---
	r.Route("/books", func(r chi.Router) {
		r.Get("/", handlers.GetBooks)
		r.Post("/", handlers.PostBook)
		r.Get("/total", handlers.GetTotalBookPrice)
		r.Get("/{id}", handlers.GetBookByID)
		r.Put("/{id}", handlers.UpdateBook)
		r.Delete("/{id}", handlers.DeleteBook)
	})

	// --- Albums API ---
	r.Route("/albums", func(r chi.Router) {
		r.Get("/", handlers.GetAllAlbums)
		r.Post("/", handlers.CreateAlbum)
		r.Get("/artist/{name}", handlers.GetAlbumsByArtist)
		r.Get("/timeout", handlers.QueryWithTimeout)
		r.Get("/{id}/can-purchase", handlers.CanPurchaseAlbum)
		r.Get("/{id}", handlers.GetAlbumByID)
	})

	// --- Orders API ---
	r.Route("/orders", func(r chi.Router) {
		r.Get("/", handlers.GetOrdersByUser)
		r.Post("/", handlers.CreateOrderByUser)
	})

	// --- Misc Handlers ---
	r.Get("/customer-name", handlers.GetCustomerName)
	r.Get("/admin/multi-query", handlers.HandleMultipleResultSets)

	// --- Wiki Pages ---
	r.Get("/view", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
	})
	r.Get("/view/{title}", handlers.ViewWiki)
	r.Get("/edit/{title}", handlers.EditWiki)
	r.Post("/save/{title}", handlers.SaveWiki)

	// --- JSON Utilities ---
	r.Post("/json/encode", handlers.JsonEncode)
	r.Post("/json/decode", handlers.JsonDecode)

	// --- WebSocket ---
	r.Get("/ws", handlers.Echo)                  // special endpoint → upgrades HTTP to a WebSocket connection.
	r.Get("/websockets", handlers.WebsocketPage) // opens the HTML page in browser.

	// --- Concurrency ---
	r.Get("/concurrency/goroutines_waitgroup", handlers.GoroutinesWaitGroupHandler)
	r.Get("/concurrency/channels_unbuffered", handlers.ChannelsUnbufferedHandler)
	r.Get("/concurrency/buffered_channels", handlers.BufferedChannelsHandler)
	r.Get("/concurrency/mutex", handlers.MutexHandler)
	r.Get("/concurrency/rwmutex", handlers.RWMutexHandler)
	r.Get("/concurrency/worker_pool", handlers.WorkerPoolHandler)
	r.Get("/concurrency/atomic_counters", handlers.AtomicCountersHandler)
	r.Get("/concurrency/cond_synccond", handlers.CondSyncCondHandler)
	r.Get("/concurrency/pool_once_map", handlers.PoolOnceMapHandler)
	r.Get("/concurrency/context_cancellation", handlers.ContextCancellationHandler)

	r.Get("/go-basics", handlers.GoBasics)

	r.Get("/pathfinder", handlers.Pathfinder)

	r.Get("/runtime-errors", handlers.RuntimeErrorsHandler)

	return r
}
