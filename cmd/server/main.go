package main

import (
	"log"
	"net/http"
	"os"
	"runtime/trace"

	"github.com/shahinzaman102/Go_JumpStart/internal/config"
	"github.com/shahinzaman102/Go_JumpStart/internal/data"
	"github.com/shahinzaman102/Go_JumpStart/internal/db"
	"github.com/shahinzaman102/Go_JumpStart/internal/handlers"
	"github.com/shahinzaman102/Go_JumpStart/internal/routes"

	_ "net/http/pprof"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Start runtime tracing
	traceFile, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("failed to create trace output file: %v", err)
	}
	defer func() {
		traceFile.Close()
		log.Println("trace file closed")
	}()

	if err := trace.Start(traceFile); err != nil {
		log.Fatalf("failed to start trace: %v", err)
	}
	defer func() {
		trace.Stop()
		log.Println("trace stopped")
	}()
	log.Println("Go trace started")

	// Ensure data directory exists
	config.EnsureDataDir()

	// Initialize database, session store & cache
	conn := config.InitDB()
	defer func() {
		conn.Close()
		log.Println("database connection closed")
	}()
	config.InitSession()
	data.InitDBConnection(conn)
	data.InitCache()

	authRepo := data.NewAuthRepo(conn)
	handlers.Init(config.Store, authRepo)

	// Preload wiki templates
	if err := handlers.LoadWikiTemplates(); err != nil {
		log.Fatalf("failed to load wiki templates: %v", err)
	}

	// Execute database schema
	if err := db.ExecuteSchema(conn, "schema.sql"); err != nil {
		log.Fatalf("schema execution failed: %v", err)
	}
	log.Println("database schema executed successfully")

	// Register routes
	router := routes.Register()

	// Start pprof server in background
	go func() {
		log.Println("pprof server listening on :6060")
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof server error: %v", err)
		}
	}()

	// Start main HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local development
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	log.Printf("server running at http://localhost:%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("main server error: %v", err)
	}

	log.Println("server stopped")
}
