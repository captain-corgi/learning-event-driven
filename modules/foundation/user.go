package main

import (
	"time"
)

// User represents a user entity in our system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserService defines the interface for user operations
type UserService interface {
	// GetUsers returns all users
	GetUsers() ([]User, error)

	// GetUserByID returns a user by their ID
	GetUserByID(id string) (*User, error)

	// CreateUser creates a new user
	CreateUser(name, email string) (*User, error)

	// UpdateUser updates an existing user
	UpdateUser(id, name, email string) (*User, error)

	// DeleteUser deletes a user by ID
	DeleteUser(id string) error
}

// NewUser creates a new User instance with generated ID and timestamps
func NewUser(name, email string) *User {
	now := time.Now()
	return &User{
		ID:        generateID(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update updates the user's fields and timestamp
func (u *User) Update(name, email string) {
	// Create a temporary user to validate new values
	temp := &User{Name: name, Email: email}
	if err := temp.Validate(); err != nil {
		return // or return the error
	}
	if name != "" {
		u.Name = name
	}
	if email != "" {
		u.Email = email
	}
	u.UpdatedAt = time.Now()
}

// Validate checks if the user has valid data
func (u *User) Validate() error {
	if u.Name == "" {
		return NewValidationError("name", "name cannot be empty")
	}
	if u.Email == "" {
		return NewValidationError("email", "email cannot be empty")
	}
	// Simple email validation
	if !isValidEmail(u.Email) {
		return NewValidationError("email", "email format is invalid")
	}
	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Simple validation - contains @ and at least one dot after @
	atIndex := -1
	for i, char := range email {
		if char == '@' {
			if atIndex != -1 {
				return false // Multiple @ symbols
			}
			atIndex = i
		}
	}

	if atIndex == -1 || atIndex == 0 || atIndex == len(email)-1 {
		return false
	}

	// Check for dot after @
	for i := atIndex + 1; i < len(email); i++ {
		if email[i] == '.' && i < len(email)-1 {
			return true
		}
	}

	return false
}
