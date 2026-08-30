package main

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/dporkka/go-sqlite-htmx/internal/db"
	"github.com/dporkka/go-sqlite-htmx/internal/handlers"
)

//go:embed ui/html/*.tmpl
var templatesFS embed.FS

//go:embed ui/static
var staticFS embed.FS

func main() {
	// Configuration
	const (
		dbPath  = "./app.db"
		port    = "8080"
		timeout = 30 * time.Second
	)

	// Initialize logger
	logger := log.New(os.Stdout, "[perma-app] ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Starting Permacomputing Web Application")

	// Initialize database
	logger.Println("Initializing SQLite database...")
	database, err := db.Open(dbPath)
	if err != nil {
		logger.Fatalf("Failed to open database: %v", err)
	}
	defer database.Close()
	logger.Println("Database initialized successfully")

	// Parse templates from embedded filesystem
	logger.Println("Loading embedded templates...")
	templates, err := template.ParseFS(templatesFS, "ui/html/*.tmpl")
	if err != nil {
		logger.Fatalf("Failed to parse templates: %v", err)
	}
	logger.Println("Templates loaded successfully")

	// Create application instance
	app := &handlers.Application{
		DB:        database,
		Templates: templates,
	}

	// Create router using Go 1.22+ ServeMux with method and path patterns
	mux := http.NewServeMux()

	// Serve static files from embedded filesystem
	staticSubFS, err := fs.Sub(staticFS, "ui/static")
	if err != nil {
		logger.Fatalf("Failed to create static sub-filesystem: %v", err)
	}
	fileServer := http.FileServer(http.FS(staticSubFS))
	mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	// Application routes
	mux.HandleFunc("GET /{$}", app.Home)                     // Exact match for root
	mux.HandleFunc("GET /user/{id}/view", app.UserView)      // View a user
	mux.HandleFunc("GET /user/{id}/edit", app.UserEdit)      // Edit form for a user
	mux.HandleFunc("POST /user/{id}/update", app.UserUpdate) // Update a user

	// Apply middleware chain
	handler := loggingMiddleware(logger)(panicRecoveryMiddleware(logger)(mux))

	// Configure server
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		IdleTimeout:  timeout,
	}

	// Start server
	logger.Printf("Server starting on http://localhost:%s", port)
	logger.Println("Press Ctrl+C to stop")

	if err := server.ListenAndServe(); err != nil {
		logger.Fatalf("Server failed to start: %v", err)
	}
}

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Printf(
				"%s %s %s %d %v",
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				wrapped.statusCode,
				duration,
			)
		})
	}
}

// panicRecoveryMiddleware recovers from panics and logs the error
func panicRecoveryMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.Printf("PANIC: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.ResponseWriter.Write(b)
}

func (rw *responseWriter) Flush() {
	if flusher, ok := rw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Compile-time check that responseWriter implements http.ResponseWriter
var _ http.ResponseWriter = (*responseWriter)(nil)

// Optional: Implement http.Hijacker if needed
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter does not implement http.Hijacker")
	}
	return hijacker.Hijack()
}
