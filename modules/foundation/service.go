package main

import (
	"crypto/rand"
	"fmt"
	"sync"
)

// InMemoryUserService implements UserService using in-memory storage
type InMemoryUserService struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewInMemoryUserService creates a new instance of InMemoryUserService
func NewInMemoryUserService() *InMemoryUserService {
	service := &InMemoryUserService{
		users: make(map[string]*User),
	}

	// Seed with some initial data
	service.seedData()

	return service
}

// seedData adds some initial users for demonstration
func (s *InMemoryUserService) seedData() {
	users := []*User{
		NewUser("John Doe", "john.doe@example.com"),
		NewUser("Jane Smith", "jane.smith@example.com"),
		NewUser("Bob Johnson", "bob.johnson@example.com"),
	}

	for _, user := range users {
		s.users[user.ID] = user
	}
}

// GetUsers returns all users
func (s *InMemoryUserService) GetUsers() ([]User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, *user)
	}

	return users, nil
}

// GetUserByID returns a user by their ID
func (s *InMemoryUserService) GetUserByID(id string) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, NewNotFoundError("user", id)
	}

	// Return a copy to prevent external modification
	userCopy := *user
	return &userCopy, nil
}

// CreateUser creates a new user
func (s *InMemoryUserService) CreateUser(name, email string) (*User, error) {
	// Check if email already exists
	if err := s.checkEmailExists(email); err != nil {
		return nil, err
	}

	user := NewUser(name, email)

	// Validate the user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.users[user.ID] = user

	// Return a copy
	userCopy := *user
	return &userCopy, nil
}

// UpdateUser updates an existing user
func (s *InMemoryUserService) UpdateUser(id, name, email string) (*User, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user, exists := s.users[id]
	if !exists {
		return nil, NewNotFoundError("user", id)
	}

	// Check if email already exists for another user
	if email != "" && email != user.Email {
		for _, existingUser := range s.users {
			if existingUser.ID != id && existingUser.Email == email {
				return nil, NewConflictError("email already exists")
			}
		}
	}

	// Update the user
	user.Update(name, email)

	// Validate the updated user
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Return a copy
	userCopy := *user
	return &userCopy, nil
}

// DeleteUser deletes a user by ID
func (s *InMemoryUserService) DeleteUser(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.users[id]; !exists {
		return NewNotFoundError("user", id)
	}

	delete(s.users, id)
	return nil
}

// checkEmailExists checks if an email already exists
func (s *InMemoryUserService) checkEmailExists(email string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return NewConflictError("email already exists")
		}
	}
	return nil
}

// generateID generates a simple random ID for demonstration
func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}
