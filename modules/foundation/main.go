package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort = "8080"
	defaultHost = "localhost"
)

func main() {
	// Get configuration from environment variables
	port := getEnv("PORT", defaultPort)
	host := getEnv("HOST", defaultHost)

	// Create user service
	userService := NewInMemoryUserService()

	// Create handlers
	userHandler := NewUserHandler(userService)

	// Setup routes
	mux := http.NewServeMux()

	// API routes
	mux.Handle("/users", userHandler)
	mux.Handle("/users/", userHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", rootHandler)

	// Create server
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s:%s", host, port)
		log.Printf("API endpoints:")
		log.Printf("  GET    /              - API information")
		log.Printf("  GET    /health        - Health check")
		log.Printf("  GET    /users         - Get all users")
		log.Printf("  POST   /users         - Create user")
		log.Printf("  GET    /users/{id}    - Get user by ID")
		log.Printf("  PUT    /users/{id}    - Update user")
		log.Printf("  DELETE /users/{id}    - Delete user")
		log.Printf("")
		log.Printf("Example requests:")
		log.Printf("  curl http://%s:%s/users", host, port)
		log.Printf("  curl -X POST http://%s:%s/users -H 'Content-Type: application/json' -d '{\"name\":\"Alice\",\"email\":\"alice@example.com\"}'", host, port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// getEnv gets an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Log the request
		duration := time.Since(start)
		log.Printf("%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			wrapper.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
