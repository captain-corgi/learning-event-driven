package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errType ErrorType
	}{
		{
			name: "valid user",
			user: User{
				ID:    "123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			user: User{
				ID:    "123",
				Name:  "",
				Email: "john@example.com",
			},
			wantErr: true,
			errType: ErrorTypeValidation,
		},
		{
			name: "empty email",
			user: User{
				ID:    "123",
				Name:  "John Doe",
				Email: "",
			},
			wantErr: true,
			errType: ErrorTypeValidation,
		},
		{
			name: "invalid email format",
			user: User{
				ID:    "123",
				Name:  "John Doe",
				Email: "invalid-email",
			},
			wantErr: true,
			errType: ErrorTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if appErr, ok := IsAppError(err); ok {
					if appErr.Type != tt.errType {
						t.Errorf("User.Validate() error type = %v, want %v", appErr.Type, tt.errType)
					}
				} else {
					t.Errorf("User.Validate() expected AppError, got %T", err)
				}
			}
		})
	}
}

func TestInMemoryUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name    string
		svcName string
		email   string
		wantErr bool
		errType ErrorType
	}{
		{
			name:    "valid user creation",
			svcName: "Alice Johnson",
			email:   "alice@example.com",
			wantErr: false,
		},
		{
			name:    "empty name",
			svcName: "",
			email:   "test@example.com",
			wantErr: true,
			errType: ErrorTypeValidation,
		},
		{
			name:    "invalid email",
			svcName: "Test User",
			email:   "invalid-email",
			wantErr: true,
			errType: ErrorTypeValidation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewInMemoryUserService()

			user, err := service.CreateUser(tt.svcName, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("CreateUser() returned nil user")
					return
				}
				if user.Name != tt.svcName {
					t.Errorf("CreateUser() name = %v, want %v", user.Name, tt.svcName)
				}
				if user.Email != tt.email {
					t.Errorf("CreateUser() email = %v, want %v", user.Email, tt.email)
				}
				if user.ID == "" {
					t.Error("CreateUser() ID should not be empty")
				}
			}
		})
	}
}

func TestInMemoryUserService_GetUserByID(t *testing.T) {
	service := NewInMemoryUserService()

	// Create a test user
	createdUser, err := service.CreateUser("Test User", "test@example.com")
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	tests := []struct {
		name    string
		userID  string
		wantErr bool
		errType ErrorType
	}{
		{
			name:    "existing user",
			userID:  createdUser.ID,
			wantErr: false,
		},
		{
			name:    "non-existing user",
			userID:  "non-existing-id",
			wantErr: true,
			errType: ErrorTypeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := service.GetUserByID(tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if user == nil {
					t.Error("GetUserByID() returned nil user")
					return
				}
				if user.ID != tt.userID {
					t.Errorf("GetUserByID() ID = %v, want %v", user.ID, tt.userID)
				}
			}
		})
	}
}

func TestUserHandler_GetUsers(t *testing.T) {
	service := NewInMemoryUserService()
	handler := NewUserHandler(service)

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var users []User
	if err := json.Unmarshal(rr.Body.Bytes(), &users); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Should have seeded users
	if len(users) == 0 {
		t.Error("Expected seeded users, got empty list")
	}
}

func TestUserHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
	}{
		{
			name:           "valid user creation",
			requestBody:    `{"name":"Test User","email":"test@example.com"}`,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid JSON",
			requestBody:    `{"name":"Test User","email":}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing name",
			requestBody:    `{"email":"test@example.com"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid email",
			requestBody:    `{"name":"Test User","email":"invalid-email"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewInMemoryUserService()
			handler := NewUserHandler(service)

			req, err := http.NewRequest("POST", "/users", strings.NewReader(tt.requestBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "user@mail.example.com", true},
		{"no @ symbol", "testexample.com", false},
		{"multiple @ symbols", "test@@example.com", false},
		{"no domain", "test@", false},
		{"no local part", "@example.com", false},
		{"no dot in domain", "test@example", false},
		{"dot at end", "test@example.", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidEmail(tt.email); got != tt.want {
				t.Errorf("isValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUser(t *testing.T) {
	name := "Test User"
	email := "test@example.com"

	user := NewUser(name, email)

	if user.Name != name {
		t.Errorf("NewUser() name = %v, want %v", user.Name, name)
	}
	if user.Email != email {
		t.Errorf("NewUser() email = %v, want %v", user.Email, email)
	}
	if user.ID == "" {
		t.Error("NewUser() ID should not be empty")
	}
	if user.CreatedAt.IsZero() {
		t.Error("NewUser() CreatedAt should not be zero")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("NewUser() UpdatedAt should not be zero")
	}
}

func TestUser_Update(t *testing.T) {
	user := NewUser("Original Name", "original@example.com")
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure timestamp difference
	time.Sleep(time.Millisecond)

	newName := "Updated Name"
	newEmail := "updated@example.com"
	user.Update(newName, newEmail)

	if user.Name != newName {
		t.Errorf("Update() name = %v, want %v", user.Name, newName)
	}
	if user.Email != newEmail {
		t.Errorf("Update() email = %v, want %v", user.Email, newEmail)
	}
	if !user.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Update() should update the UpdatedAt timestamp")
	}
}
