# Module 1: Foundations of Go & Microservices

This module demonstrates the fundamentals of Go programming and building basic HTTP microservices. It covers Go types, functions, error handling, interfaces, and implements a simple REST API for user management.

## Learning Objectives

- ✅ Install Go, set up workspace, understand modules & packages
- ✅ Write idiomatic Go: types, functions, error handling, interfaces
- ✅ Basic HTTP microservice with `net/http`
- ✅ Implement REST endpoints with proper error handling
- ✅ Use in-memory storage for data persistence

## Project Structure

```shell
modules/foundation/
├── go.mod              # Go module definition
├── main.go             # HTTP server and application entry point
├── user.go             # User entity and domain logic
├── service.go          # User service implementation (in-memory)
├── handlers.go         # HTTP handlers for REST API
├── errors.go           # Custom error types and error handling
├── main_test.go        # Unit tests (table-driven testing)
└── README.md           # This documentation
```

## Key Concepts Demonstrated

### 1. Go Fundamentals

- **Structs and Methods**: `User` struct with validation methods
- **Interfaces**: `UserService` interface for dependency injection
- **Error Handling**: Custom error types with proper wrapping
- **Packages**: Organized code structure with clear separation of concerns

### 2. HTTP Microservice

- **REST API**: Full CRUD operations for user management
- **Middleware**: Logging middleware for request tracking
- **Graceful Shutdown**: Proper server lifecycle management
- **Configuration**: Environment variable support

### 3. Error Patterns

- **Custom Error Types**: `AppError` with different error categories
- **Error Wrapping**: Using `github.com/pkg/errors` for context
- **HTTP Error Mapping**: Converting domain errors to HTTP status codes

## API Endpoints

| Method | Endpoint | Description | Request Body | Response |
|--------|----------|-------------|--------------|----------|
| GET | `/` | API information | - | API metadata |
| GET | `/health` | Health check | - | Service status |
| GET | `/users` | Get all users | - | Array of users |
| POST | `/users` | Create user | `{"name":"string","email":"string"}` | Created user |
| GET | `/users/{id}` | Get user by ID | - | User object |
| PUT | `/users/{id}` | Update user | `{"name":"string","email":"string"}` | Updated user |
| DELETE | `/users/{id}` | Delete user | - | 204 No Content |

## Running the Application

### Prerequisites

- Go 1.24.0 or later
- Git (for dependency management)

### Installation and Setup

1. **Navigate to the module directory:**

   ```bash
   cd modules/foundation
   ```

2. **Download dependencies:**

   ```bash
   go mod tidy
   ```

3. **Run the application:**

   ```bash
   go run .
   ```

4. **The server will start on `localhost:8080`**

### Environment Variables

- `PORT`: Server port (default: 8080)
- `HOST`: Server host (default: localhost)

### Example Usage

1. **Get all users:**

   ```bash
   curl http://localhost:8080/users
   ```

2. **Create a new user:**

   ```bash
   curl -X POST http://localhost:8080/users \
     -H "Content-Type: application/json" \
     -d '{"name":"Alice Johnson","email":"alice@example.com"}'
   ```

3. **Get user by ID:**

   ```bash
   curl http://localhost:8080/users/{user-id}
   ```

4. **Update user:**

   ```bash
   curl -X PUT http://localhost:8080/users/{user-id} \
     -H "Content-Type: application/json" \
     -d '{"name":"Alice Smith","email":"alice.smith@example.com"}'
   ```

5. **Delete user:**

   ```bash
   curl -X DELETE http://localhost:8080/users/{user-id}
   ```

## Testing

### Running Tests

```bash
# Run all tests
go test -v

# Run tests with coverage
go test -v -cover

# Run specific test
go test -v -run TestUser_Validate
```

### Test Coverage

The test suite includes:

- **Unit Tests**: Testing individual functions and methods
- **Integration Tests**: Testing HTTP handlers with mock requests
- **Table-Driven Tests**: Comprehensive test cases using Go's testing patterns
- **Error Handling Tests**: Validating error scenarios and edge cases

## Architecture Patterns

### 1. Dependency Injection

```go
type UserService interface {
    GetUsers() ([]User, error)
    GetUserByID(id string) (*User, error)
    CreateUser(name, email string) (*User, error)
    UpdateUser(id, name, email string) (*User, error)
    DeleteUser(id string) error
}
```

### 2. Error Handling

```go
type AppError struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Field   string    `json:"field,omitempty"`
    Cause   error     `json:"-"`
}
```

### 3. HTTP Handler Pattern

```go
type UserHandler struct {
    service UserService
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Route handling logic
}
```

## Key Go Concepts Learned

1. **Structs and Methods**: Defining data structures and behavior
2. **Interfaces**: Defining contracts for dependency injection
3. **Error Handling**: Custom error types and proper error propagation
4. **HTTP Server**: Building REST APIs with `net/http`
5. **JSON Marshaling**: Converting between Go structs and JSON
6. **Middleware**: Request/response processing pipeline
7. **Testing**: Table-driven tests and HTTP testing utilities
8. **Modules**: Dependency management with `go.mod`

## Next Steps

This foundation module prepares you for:

- **Module 2**: Clean Architecture patterns
- **Module 3**: Domain-Driven Design principles
- **Module 4**: Event-driven architecture concepts
- **Module 5**: Message brokers and event publishing

## Additional Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)
- [Go Testing](https://golang.org/pkg/testing/)
