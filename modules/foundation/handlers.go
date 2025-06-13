package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	service UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// ServeHTTP implements http.Handler interface for routing
func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set common headers
	w.Header().Set("Content-Type", "application/json")

	// Parse the path
	path := strings.TrimPrefix(r.URL.Path, "/users")

	switch {
	case path == "" || path == "/":
		switch r.Method {
		case http.MethodGet:
			h.handleGetUsers(w, r)
		case http.MethodPost:
			h.handleCreateUser(w, r)
		default:
			h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	case strings.HasPrefix(path, "/"):
		userID := strings.TrimPrefix(path, "/")
		switch r.Method {
		case http.MethodGet:
			h.handleGetUser(w, r, userID)
		case http.MethodPut:
			h.handleUpdateUser(w, r, userID)
		case http.MethodDelete:
			h.handleDeleteUser(w, r, userID)
		default:
			h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	default:
		h.writeErrorResponse(w, http.StatusNotFound, "endpoint not found")
	}
}

// handleGetUsers handles GET /users
func (h *UserHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetUsers()
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, users)
}

// handleGetUser handles GET /users/{id}
func (h *UserHandler) handleGetUser(w http.ResponseWriter, r *http.Request, userID string) {
	user, err := h.service.GetUserByID(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

// CreateUserRequest represents the request body for creating a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// handleCreateUser handles POST /users
func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := h.service.CreateUser(req.Name, req.Email)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, user)
}

// UpdateUserRequest represents the request body for updating a user
type UpdateUserRequest struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// handleUpdateUser handles PUT /users/{id}
func (h *UserHandler) handleUpdateUser(w http.ResponseWriter, r *http.Request, userID string) {
	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	user, err := h.service.UpdateUser(userID, req.Name, req.Email)
	if err != nil {
		h.handleError(w, err)
		return
	}

	h.writeJSONResponse(w, http.StatusOK, user)
}

// handleDeleteUser handles DELETE /users/{id}
func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, userID string) {
	err := h.service.DeleteUser(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError handles application errors and writes appropriate HTTP responses
func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	if appErr, ok := IsAppError(err); ok {
		h.writeJSONResponse(w, appErr.HTTPStatusCode(), map[string]interface{}{
			"error": map[string]interface{}{
				"type":    appErr.Type,
				"message": appErr.Message,
				"field":   appErr.Field,
			},
		})
		return
	}

	// Log unexpected errors
	log.Printf("Unexpected error: %v", err)
	h.writeErrorResponse(w, http.StatusInternalServerError, "internal server error")
}

// writeJSONResponse writes a JSON response
func (h *UserHandler) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

// writeErrorResponse writes a simple error response
func (h *UserHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	h.writeJSONResponse(w, statusCode, map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
		},
	})
}

// healthHandler handles health check requests
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"status":  "healthy",
		"service": "user-service",
		"version": "1.0.0",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// rootHandler handles requests to the root path
func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message": "Welcome to User Service API",
		"version": "1.0.0",
		"endpoints": map[string]interface{}{
			"users": map[string]interface{}{
				"GET /users":         "Get all users",
				"POST /users":        "Create a new user",
				"GET /users/{id}":    "Get user by ID",
				"PUT /users/{id}":    "Update user by ID",
				"DELETE /users/{id}": "Delete user by ID",
			},
			"health": "GET /health - Health check",
		},
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding root response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
