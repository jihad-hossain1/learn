# 19. Testing in Go

Go has a built-in testing framework in the standard library's `testing` package. This framework provides a simple way to write unit tests, benchmarks, and examples for your code. In this chapter, we'll explore Go's testing capabilities and best practices.

## Basic Unit Testing

### Writing Your First Test

Go tests are functions that:
- Are in files with names ending in `_test.go`
- Are in the same package as the code they test
- Are functions named `TestXxx` (where `Xxx` is the function being tested)
- Take a single parameter of type `*testing.T`

Let's start with a simple example:

```go
// math.go
package math

// Add returns the sum of a and b
func Add(a, b int) int {
    return a + b
}

// Subtract returns the difference of a and b
func Subtract(a, b int) int {
    return a - b
}

// Multiply returns the product of a and b
func Multiply(a, b int) int {
    return a * b
}

// Divide returns the quotient of a and b
func Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}
```

```go
// math_test.go
package math

import (
    "testing"
)

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    want := 5
    
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}

func TestSubtract(t *testing.T) {
    got := Subtract(5, 3)
    want := 2
    
    if got != want {
        t.Errorf("Subtract(5, 3) = %d; want %d", got, want)
    }
}

func TestMultiply(t *testing.T) {
    got := Multiply(2, 3)
    want := 6
    
    if got != want {
        t.Errorf("Multiply(2, 3) = %d; want %d", got, want)
    }
}

func TestDivide(t *testing.T) {
    got, err := Divide(6, 3)
    want := 2
    
    if err != nil {
        t.Errorf("Divide(6, 3) returned unexpected error: %v", err)
    }
    
    if got != want {
        t.Errorf("Divide(6, 3) = %d; want %d", got, want)
    }
}

func TestDivideByZero(t *testing.T) {
    _, err := Divide(6, 0)
    
    if err == nil {
        t.Error("Divide(6, 0) expected an error, got nil")
    }
}
```

### Running Tests

To run tests, use the `go test` command:

```bash
# Run all tests in the current package
go test

# Run tests with verbose output
go test -v

# Run a specific test
go test -run TestAdd

# Run tests matching a pattern
go test -run "TestDivide.*"

# Run tests in all subdirectories
go test ./...
```

### Table-Driven Tests

Table-driven tests are a common pattern in Go for testing multiple inputs and expected outputs:

```go
func TestAdd_TableDriven(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed signs", -2, 3, 1},
        {"zeros", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", 
                    tt.a, tt.b, got, tt.expected)
            }
        })
    }
}

func TestDivide_TableDriven(t *testing.T) {
    tests := []struct {
        name        string
        a, b        int
        expected    int
        expectError bool
    }{
        {"simple division", 6, 3, 2, false},
        {"division with remainder", 7, 3, 2, false},
        {"divide by zero", 6, 0, 0, true},
        {"zero divided by number", 0, 5, 0, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Divide(tt.a, tt.b)
            
            // Check error expectation
            if (err != nil) != tt.expectError {
                t.Errorf("Divide(%d, %d) error = %v; expectError = %v", 
                    tt.a, tt.b, err, tt.expectError)
                return
            }
            
            // If we expect an error, don't check the result
            if tt.expectError {
                return
            }
            
            if got != tt.expected {
                t.Errorf("Divide(%d, %d) = %d; want %d", 
                    tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

## Test Organization

### Test Helpers

Test helpers are functions that help set up or tear down test environments:

```go
func setupTestCase(t *testing.T) func(t *testing.T) {
    t.Log("Setting up test case")
    
    // Return a function to be called at the end of the test
    return func(t *testing.T) {
        t.Log("Tearing down test case")
    }
}

func TestWithHelper(t *testing.T) {
    // Call the setup function, which returns a teardown function
    teardown := setupTestCase(t)
    // Ensure teardown is called at the end of the test
    defer teardown(t)
    
    // Test code here
    t.Log("Running test")
}
```

### Subtests

Subtests allow you to group related tests together:

```go
func TestMathOperations(t *testing.T) {
    t.Run("Addition", func(t *testing.T) {
        got := Add(2, 3)
        want := 5
        if got != want {
            t.Errorf("Add(2, 3) = %d; want %d", got, want)
        }
    })
    
    t.Run("Subtraction", func(t *testing.T) {
        got := Subtract(5, 3)
        want := 2
        if got != want {
            t.Errorf("Subtract(5, 3) = %d; want %d", got, want)
        }
    })
    
    t.Run("Multiplication", func(t *testing.T) {
        got := Multiply(2, 3)
        want := 6
        if got != want {
            t.Errorf("Multiply(2, 3) = %d; want %d", got, want)
        }
    })
    
    t.Run("Division", func(t *testing.T) {
        t.Run("Valid", func(t *testing.T) {
            got, err := Divide(6, 3)
            want := 2
            if err != nil {
                t.Errorf("Divide(6, 3) returned unexpected error: %v", err)
            }
            if got != want {
                t.Errorf("Divide(6, 3) = %d; want %d", got, want)
            }
        })
        
        t.Run("DivideByZero", func(t *testing.T) {
            _, err := Divide(6, 0)
            if err == nil {
                t.Error("Divide(6, 0) expected an error, got nil")
            }
        })
    })
}
```

### Parallel Tests

Parallel tests can speed up test execution:

```go
func TestAddParallel(t *testing.T) {
    t.Parallel() // Mark this test as parallel
    
    got := Add(2, 3)
    want := 5
    
    if got != want {
        t.Errorf("Add(2, 3) = %d; want %d", got, want)
    }
}

func TestSubtractParallel(t *testing.T) {
    t.Parallel() // Mark this test as parallel
    
    got := Subtract(5, 3)
    want := 2
    
    if got != want {
        t.Errorf("Subtract(5, 3) = %d; want %d", got, want)
    }
}
```

## Testing HTTP Handlers

### Basic HTTP Handler Testing

```go
// handlers.go
package handlers

import (
    "encoding/json"
    "net/http"
)

type Response struct {
    Message string `json:"message"`
    Status  int    `json:"status"`
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Message: "Hello, World!",
        Status:  http.StatusOK,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func GreetHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "Guest"
    }
    
    response := Response{
        Message: "Hello, " + name + "!",
        Status:  http.StatusOK,
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}
```

```go
// handlers_test.go
package handlers

import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHelloHandler(t *testing.T) {
    // Create a request
    req, err := http.NewRequest("GET", "/hello", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Create the handler
    handler := http.HandlerFunc(HelloHandler)
    
    // Serve the request
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    // Check the content type
    expectedContentType := "application/json"
    if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
        t.Errorf("handler returned wrong content type: got %v want %v",
            contentType, expectedContentType)
    }
    
    // Check the response body
    var response Response
    if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
        t.Errorf("error decoding response: %v", err)
    }
    
    expectedMessage := "Hello, World!"
    if response.Message != expectedMessage {
        t.Errorf("handler returned unexpected message: got %v want %v",
            response.Message, expectedMessage)
    }
}

func TestGreetHandler(t *testing.T) {
    tests := []struct {
        name           string
        queryParam     string
        expectedStatus int
        expectedMsg    string
    }{
        {"with name", "John", http.StatusOK, "Hello, John!"},
        {"without name", "", http.StatusOK, "Hello, Guest!"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a request with query parameter
            req, err := http.NewRequest("GET", "/greet", nil)
            if err != nil {
                t.Fatal(err)
            }
            
            // Add query parameter if provided
            if tt.queryParam != "" {
                q := req.URL.Query()
                q.Add("name", tt.queryParam)
                req.URL.RawQuery = q.Encode()
            }
            
            // Create a response recorder
            rr := httptest.NewRecorder()
            
            // Create the handler
            handler := http.HandlerFunc(GreetHandler)
            
            // Serve the request
            handler.ServeHTTP(rr, req)
            
            // Check the status code
            if status := rr.Code; status != tt.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatus)
            }
            
            // Check the response body
            var response Response
            if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
                t.Errorf("error decoding response: %v", err)
            }
            
            if response.Message != tt.expectedMsg {
                t.Errorf("handler returned unexpected message: got %v want %v",
                    response.Message, tt.expectedMsg)
            }
        })
    }
}
```

### Testing HTTP Middleware

```go
// middleware.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

// Logger is a middleware that logs request details
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call the next handler
        next.ServeHTTP(w, r)
        
        // Log request details
        log.Printf(
            "%s %s %s %s",
            r.Method,
            r.RequestURI,
            r.RemoteAddr,
            time.Since(start),
        )
    })
}

// Auth is a middleware that checks for an API key
func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        
        if apiKey != "valid-key" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Call the next handler
        next.ServeHTTP(w, r)
    })
}
```

```go
// middleware_test.go
package middleware

import (
    "bytes"
    "log"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestLogger(t *testing.T) {
    // Create a buffer to capture log output
    var buf bytes.Buffer
    log.SetOutput(&buf)
    
    // Create a test handler
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })
    
    // Wrap the test handler with the logger middleware
    handler := Logger(testHandler)
    
    // Create a test request
    req, err := http.NewRequest("GET", "/test", nil)
    if err != nil {
        t.Fatal(err)
    }
    
    // Create a response recorder
    rr := httptest.NewRecorder()
    
    // Serve the request
    handler.ServeHTTP(rr, req)
    
    // Check the status code
    if status := rr.Code; status != http.StatusOK {
        t.Errorf("handler returned wrong status code: got %v want %v",
            status, http.StatusOK)
    }
    
    // Check the response body
    expected := "OK"
    if rr.Body.String() != expected {
        t.Errorf("handler returned unexpected body: got %v want %v",
            rr.Body.String(), expected)
    }
    
    // Check that the log contains expected information
    logOutput := buf.String()
    if !strings.Contains(logOutput, "GET") || !strings.Contains(logOutput, "/test") {
        t.Errorf("log output missing expected information: %s", logOutput)
    }
}

func TestAuth(t *testing.T) {
    // Create a test handler
    testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Protected Resource"))
    })
    
    // Wrap the test handler with the auth middleware
    handler := Auth(testHandler)
    
    tests := []struct {
        name           string
        apiKey         string
        expectedStatus int
        expectedBody   string
    }{
        {"valid key", "valid-key", http.StatusOK, "Protected Resource"},
        {"invalid key", "invalid-key", http.StatusUnauthorized, "Unauthorized\n"},
        {"missing key", "", http.StatusUnauthorized, "Unauthorized\n"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create a test request
            req, err := http.NewRequest("GET", "/protected", nil)
            if err != nil {
                t.Fatal(err)
            }
            
            // Add API key header if provided
            if tt.apiKey != "" {
                req.Header.Set("X-API-Key", tt.apiKey)
            }
            
            // Create a response recorder
            rr := httptest.NewRecorder()
            
            // Serve the request
            handler.ServeHTTP(rr, req)
            
            // Check the status code
            if status := rr.Code; status != tt.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v",
                    status, tt.expectedStatus)
            }
            
            // Check the response body
            if rr.Body.String() != tt.expectedBody {
                t.Errorf("handler returned unexpected body: got %v want %v",
                    rr.Body.String(), tt.expectedBody)
            }
        })
    }
}
```

## Mocking and Dependency Injection

### Interface-Based Mocking

```go
// user.go
package user

import (
    "errors"
    "time"
)

type User struct {
    ID        int
    Name      string
    Email     string
    CreatedAt time.Time
}

// UserRepository defines the interface for user data access
type UserRepository interface {
    GetByID(id int) (*User, error)
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
}

// UserService provides user-related operations
type UserService struct {
    repo UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user ID")
    }
    
    return s.repo.GetByID(id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name, email string) (*User, error) {
    if name == "" || email == "" {
        return nil, errors.New("name and email are required")
    }
    
    user := &User{
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
    
    if err := s.repo.Create(user); err != nil {
        return nil, err
    }
    
    return user, nil
}
```

```go
// user_test.go
package user

import (
    "errors"
    "testing"
    "time"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
    users map[int]*User
    // Track method calls for verification
    getByIDCalled bool
    createCalled  bool
    updateCalled  bool
    deleteCalled  bool
}

// NewMockUserRepository creates a new mock repository
func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[int]*User),
    }
}

// GetByID retrieves a user by ID
func (m *MockUserRepository) GetByID(id int) (*User, error) {
    m.getByIDCalled = true
    
    user, ok := m.users[id]
    if !ok {
        return nil, errors.New("user not found")
    }
    
    return user, nil
}

// Create adds a new user
func (m *MockUserRepository) Create(user *User) error {
    m.createCalled = true
    
    // Simulate ID generation
    user.ID = len(m.users) + 1
    m.users[user.ID] = user
    
    return nil
}

// Update updates an existing user
func (m *MockUserRepository) Update(user *User) error {
    m.updateCalled = true
    
    if _, ok := m.users[user.ID]; !ok {
        return errors.New("user not found")
    }
    
    m.users[user.ID] = user
    return nil
}

// Delete removes a user
func (m *MockUserRepository) Delete(id int) error {
    m.deleteCalled = true
    
    if _, ok := m.users[id]; !ok {
        return errors.New("user not found")
    }
    
    delete(m.users, id)
    return nil
}

// AddTestUser adds a test user to the mock repository
func (m *MockUserRepository) AddTestUser(id int, name, email string) {
    m.users[id] = &User{
        ID:        id,
        Name:      name,
        Email:     email,
        CreatedAt: time.Now(),
    }
}

func TestGetUser(t *testing.T) {
    // Create a mock repository
    mockRepo := NewMockUserRepository()
    
    // Add a test user
    mockRepo.AddTestUser(1, "John Doe", "john@example.com")
    
    // Create the service with the mock repository
    service := NewUserService(mockRepo)
    
    // Test valid user ID
    t.Run("Valid ID", func(t *testing.T) {
        user, err := service.GetUser(1)
        
        if err != nil {
            t.Errorf("Expected no error, got %v", err)
        }
        
        if user == nil {
            t.Fatal("Expected user, got nil")
        }
        
        if user.Name != "John Doe" {
            t.Errorf("Expected name 'John Doe', got '%s'", user.Name)
        }
        
        if !mockRepo.getByIDCalled {
            t.Error("Expected GetByID to be called")
        }
    })
    
    // Test invalid user ID
    t.Run("Invalid ID", func(t *testing.T) {
        user, err := service.GetUser(0)
        
        if err == nil {
            t.Error("Expected an error, got nil")
        }
        
        if user != nil {
            t.Errorf("Expected nil user, got %v", user)
        }
    })
    
    // Test non-existent user ID
    t.Run("Non-existent ID", func(t *testing.T) {
        user, err := service.GetUser(999)
        
        if err == nil {
            t.Error("Expected an error, got nil")
        }
        
        if user != nil {
            t.Errorf("Expected nil user, got %v", user)
        }
    })
}

func TestCreateUser(t *testing.T) {
    // Create a mock repository
    mockRepo := NewMockUserRepository()
    
    // Create the service with the mock repository
    service := NewUserService(mockRepo)
    
    // Test valid user creation
    t.Run("Valid User", func(t *testing.T) {
        user, err := service.CreateUser("Jane Doe", "jane@example.com")
        
        if err != nil {
            t.Errorf("Expected no error, got %v", err)
        }
        
        if user == nil {
            t.Fatal("Expected user, got nil")
        }
        
        if user.Name != "Jane Doe" {
            t.Errorf("Expected name 'Jane Doe', got '%s'", user.Name)
        }
        
        if user.Email != "jane@example.com" {
            t.Errorf("Expected email 'jane@example.com', got '%s'", user.Email)
        }
        
        if user.ID <= 0 {
            t.Errorf("Expected positive ID, got %d", user.ID)
        }
        
        if !mockRepo.createCalled {
            t.Error("Expected Create to be called")
        }
    })
    
    // Test invalid user creation (missing name)
    t.Run("Missing Name", func(t *testing.T) {
        user, err := service.CreateUser("", "missing@example.com")
        
        if err == nil {
            t.Error("Expected an error, got nil")
        }
        
        if user != nil {
            t.Errorf("Expected nil user, got %v", user)
        }
    })
    
    // Test invalid user creation (missing email)
    t.Run("Missing Email", func(t *testing.T) {
        user, err := service.CreateUser("Missing Email", "")
        
        if err == nil {
            t.Error("Expected an error, got nil")
        }
        
        if user != nil {
            t.Errorf("Expected nil user, got %v", user)
        }
    })
}
```

## Benchmarking

### Writing Benchmarks

Benchmarks in Go are functions that:
- Are in files with names ending in `_test.go`
- Are functions named `BenchmarkXxx` (where `Xxx` is the function being benchmarked)
- Take a single parameter of type `*testing.B`

```go
// math_bench_test.go
package math

import (
    "testing"
)

func BenchmarkAdd(b *testing.B) {
    // Run the Add function b.N times
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}

func BenchmarkSubtract(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Subtract(5, 3)
    }
}

func BenchmarkMultiply(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Multiply(2, 3)
    }
}

func BenchmarkDivide(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Divide(6, 3)
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkAdd

# Run benchmarks with more iterations
go test -bench=. -benchtime=5s

# Run benchmarks and see memory allocations
go test -bench=. -benchmem
```

### Comparing Implementations

```go
// string_concat.go
package concat

import (
    "bytes"
    "strings"
)

// ConcatWithPlus concatenates strings using +
func ConcatWithPlus(strs []string) string {
    var result string
    for _, s := range strs {
        result += s
    }
    return result
}

// ConcatWithJoin concatenates strings using strings.Join
func ConcatWithJoin(strs []string) string {
    return strings.Join(strs, "")
}

// ConcatWithBuilder concatenates strings using strings.Builder
func ConcatWithBuilder(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

// ConcatWithBuffer concatenates strings using bytes.Buffer
func ConcatWithBuffer(strs []string) string {
    var buffer bytes.Buffer
    for _, s := range strs {
        buffer.WriteString(s)
    }
    return buffer.String()
}
```

```go
// string_concat_test.go
package concat

import (
    "testing"
)

// Generate test data
func generateStrings(n int) []string {
    result := make([]string, n)
    for i := 0; i < n; i++ {
        result[i] = "a"
    }
    return result
}

func BenchmarkConcatWithPlus(b *testing.B) {
    strs := generateStrings(100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithPlus(strs)
    }
}

func BenchmarkConcatWithJoin(b *testing.B) {
    strs := generateStrings(100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithJoin(strs)
    }
}

func BenchmarkConcatWithBuilder(b *testing.B) {
    strs := generateStrings(100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithBuilder(strs)
    }
}

func BenchmarkConcatWithBuffer(b *testing.B) {
    strs := generateStrings(100)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithBuffer(strs)
    }
}

// Benchmark with different input sizes
func BenchmarkConcatWithPlus10(b *testing.B) {
    strs := generateStrings(10)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithPlus(strs)
    }
}

func BenchmarkConcatWithPlus1000(b *testing.B) {
    strs := generateStrings(1000)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithPlus(strs)
    }
}

func BenchmarkConcatWithBuilder10(b *testing.B) {
    strs := generateStrings(10)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithBuilder(strs)
    }
}

func BenchmarkConcatWithBuilder1000(b *testing.B) {
    strs := generateStrings(1000)
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        ConcatWithBuilder(strs)
    }
}
```

## Test Coverage

### Measuring Test Coverage

```bash
# Run tests with coverage
go test -cover

# Generate coverage profile
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# View coverage in terminal
go tool cover -func=coverage.out
```

### Improving Test Coverage

```go
// calculator.go
package calculator

import "errors"

func Calculate(a, b float64, op string) (float64, error) {
    switch op {
    case "+":
        return a + b, nil
    case "-":
        return a - b, nil
    case "*":
        return a * b, nil
    case "/":
        if b == 0 {
            return 0, errors.New("division by zero")
        }
        return a / b, nil
    default:
        return 0, errors.New("invalid operation")
    }
}
```

```go
// calculator_test.go
package calculator

import (
    "testing"
)

func TestCalculate(t *testing.T) {
    tests := []struct {
        name        string
        a, b        float64
        op          string
        expected    float64
        expectError bool
    }{
        {"addition", 2, 3, "+", 5, false},
        {"subtraction", 5, 3, "-", 2, false},
        {"multiplication", 2, 3, "*", 6, false},
        {"division", 6, 3, "/", 2, false},
        {"division by zero", 6, 0, "/", 0, true},
        {"invalid operation", 2, 3, "^", 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Calculate(tt.a, tt.b, tt.op)
            
            // Check error expectation
            if (err != nil) != tt.expectError {
                t.Errorf("Calculate(%f, %f, %s) error = %v; expectError = %v", 
                    tt.a, tt.b, tt.op, err, tt.expectError)
                return
            }
            
            // If we expect an error, don't check the result
            if tt.expectError {
                return
            }
            
            if result != tt.expected {
                t.Errorf("Calculate(%f, %f, %s) = %f; want %f", 
                    tt.a, tt.b, tt.op, result, tt.expected)
            }
        })
    }
}
```

## Examples as Tests

### Writing Examples

Examples in Go are functions that:
- Are in files with names ending in `_test.go`
- Are functions named `ExampleXxx` (where `Xxx` is the function being exemplified)
- Include output comments that match the expected output

```go
// math_example_test.go
package math

import (
    "fmt"
)

func ExampleAdd() {
    sum := Add(2, 3)
    fmt.Println(sum)
    // Output: 5
}

func ExampleSubtract() {
    difference := Subtract(5, 3)
    fmt.Println(difference)
    // Output: 2
}

func ExampleMultiply() {
    product := Multiply(2, 3)
    fmt.Println(product)
    // Output: 6
}

func ExampleDivide() {
    quotient, err := Divide(6, 3)
    fmt.Println(quotient, err)
    // Output: 2 <nil>
}

func ExampleDivide_byZero() {
    quotient, err := Divide(6, 0)
    fmt.Printf("quotient: %d, error: %v\n", quotient, err != nil)
    // Output: quotient: 0, error: true
}
```

## Testing Best Practices

### Test Structure

1. **Arrange**: Set up the test data and environment
2. **Act**: Call the function or method being tested
3. **Assert**: Verify the results

```go
func TestSomething(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := FunctionUnderTest(input)
    
    // Assert
    if result != expected {
        t.Errorf("FunctionUnderTest(%q) = %q; want %q", input, result, expected)
    }
}
```

### Test Naming

Use descriptive names for tests:

```go
func TestAdd_PositiveNumbers(t *testing.T) { /* ... */ }
func TestAdd_NegativeNumbers(t *testing.T) { /* ... */ }
func TestAdd_MixedSigns(t *testing.T) { /* ... */ }
```

### Test Independence

Ensure tests are independent and can run in any order:

```go
func TestIndependent1(t *testing.T) {
    // This test should not depend on any other test
    // or global state
}

func TestIndependent2(t *testing.T) {
    // This test should not depend on any other test
    // or global state
}
```

### Test Readability

Make tests easy to understand:

```go
func TestReadable(t *testing.T) {
    // Use descriptive variable names
    username := "john_doe"
    password := "secret123"
    
    // Use clear assertions
    isValid := ValidateCredentials(username, password)
    if !isValid {
        t.Errorf("ValidateCredentials(%q, %q) = false; want true", username, password)
    }
}
```

## Advanced Testing Topics

### Testing with Race Detection

```bash
# Run tests with race detection
go test -race
```

### Fuzzing

Fuzzing is a testing technique that provides random inputs to a function to find bugs:

```go
// fuzz_test.go
package calculator

import (
    "testing"
)

func FuzzCalculate(f *testing.F) {
    // Seed corpus
    f.Add(float64(2), float64(3), "+")
    f.Add(float64(5), float64(3), "-")
    f.Add(float64(2), float64(3), "*")
    f.Add(float64(6), float64(3), "/")
    
    f.Fuzz(func(t *testing.T, a, b float64, op string) {
        // Skip division by zero and invalid operations for simplicity
        if op != "+" && op != "-" && op != "*" && op != "/" {
            return
        }
        if op == "/" && b == 0 {
            return
        }
        
        result, err := Calculate(a, b, op)
        if err != nil {
            t.Errorf("Calculate(%f, %f, %s) returned unexpected error: %v", a, b, op, err)
        }
        
        // Verify basic properties
        switch op {
        case "+":
            if result != a+b {
                t.Errorf("Calculate(%f, %f, %s) = %f; want %f", a, b, op, result, a+b)
            }
        case "-":
            if result != a-b {
                t.Errorf("Calculate(%f, %f, %s) = %f; want %f", a, b, op, result, a-b)
            }
        case "*":
            if result != a*b {
                t.Errorf("Calculate(%f, %f, %s) = %f; want %f", a, b, op, result, a*b)
            }
        case "/":
            if result != a/b {
                t.Errorf("Calculate(%f, %f, %s) = %f; want %f", a, b, op, result, a/b)
            }
        }
    })
}
```

### Testing with Build Tags

Use build tags to control which tests are run:

```go
// integration_test.go
//go:build integration
// +build integration

package mypackage

import (
    "testing"
)

func TestIntegration(t *testing.T) {
    // This test will only run when the "integration" tag is specified
}
```

Run with:

```bash
go test -tags=integration
```

## Exercises

### Exercise 1: Write Unit Tests for a String Utility Package

Create a string utility package with functions for common string operations, then write comprehensive tests for it.

```go
// stringutil/stringutil.go
package stringutil

import (
    "strings"
    "unicode"
)

// Reverse returns the reverse of a string
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// Capitalize capitalizes the first letter of each word
func Capitalize(s string) string {
    words := strings.Fields(s)
    for i, word := range words {
        if len(word) > 0 {
            runes := []rune(word)
            runes[0] = unicode.ToUpper(runes[0])
            words[i] = string(runes)
        }
    }
    return strings.Join(words, " ")
}

// CountWords counts the number of words in a string
func CountWords(s string) int {
    if len(strings.TrimSpace(s)) == 0 {
        return 0
    }
    return len(strings.Fields(s))
}

// Truncate truncates a string to the specified length
func Truncate(s string, maxLength int) string {
    if maxLength <= 0 {
        return ""
    }
    if len(s) <= maxLength {
        return s
    }
    return s[:maxLength]
}
```

Write tests for this package.

### Exercise 2: Test an HTTP API Client

Create an HTTP API client and write tests for it using mocks.

```go
// apiclient/client.go
package apiclient

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Client struct {
    BaseURL    string
    HTTPClient *http.Client
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func NewClient(baseURL string) *Client {
    return &Client{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *Client) GetUser(id int) (*User, error) {
    url := fmt.Sprintf("%s/users/%d", c.BaseURL, id)
    resp, err := c.HTTPClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

func (c *Client) CreateUser(name, email string) (*User, error) {
    // Implementation left as an exercise
    return nil, nil
}

func (c *Client) UpdateUser(user *User) error {
    // Implementation left as an exercise
    return nil
}

func (c *Client) DeleteUser(id int) error {
    // Implementation left as an exercise
    return nil
}
```

Write tests for this client.

### Exercise 3: Benchmark Different Sorting Algorithms

Implement different sorting algorithms and benchmark their performance.

```go
// sorting/sorting.go
package sorting

// BubbleSort implements the bubble sort algorithm
func BubbleSort(arr []int) []int {
    n := len(arr)
    result := make([]int, n)
    copy(result, arr)
    
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    
    return result
}

// InsertionSort implements the insertion sort algorithm
func InsertionSort(arr []int) []int {
    n := len(arr)
    result := make([]int, n)
    copy(result, arr)
    
    for i := 1; i < n; i++ {
        key := result[i]
        j := i - 1
        
        for j >= 0 && result[j] > key {
            result[j+1] = result[j]
            j--
        }
        
        result[j+1] = key
    }
    
    return result
}

// QuickSort implements the quick sort algorithm
func QuickSort(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    quickSortHelper(result, 0, len(result)-1)
    return result
}

func quickSortHelper(arr []int, low, high int) {
    if low < high {
        pivot := partition(arr, low, high)
        quickSortHelper(arr, low, pivot-1)
        quickSortHelper(arr, pivot+1, high)
    }
}

func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1
    
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}

// MergeSort implements the merge sort algorithm
func MergeSort(arr []int) []int {
    result := make([]int, len(arr))
    copy(result, arr)
    
    if len(result) <= 1 {
        return result
    }
    
    mid := len(result) / 2
    left := MergeSort(result[:mid])
    right := MergeSort(result[mid:])
    
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, len(left)+len(right))
    i, j, k := 0, 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result[k] = left[i]
            i++
        } else {
            result[k] = right[j]
            j++
        }
        k++
    }
    
    for i < len(left) {
        result[k] = left[i]
        i++
        k++
    }
    
    for j < len(right) {
        result[k] = right[j]
        j++
        k++
    }
    
    return result
}
```

Write benchmarks for these sorting algorithms.

## Key Takeaways

1. **Go's testing package is simple but powerful**: The standard library's `testing` package provides all the tools you need for unit testing, benchmarking, and examples.

2. **Table-driven tests are idiomatic**: They make it easy to test multiple inputs and expected outputs with minimal code duplication.

3. **Subtests and parallel tests improve organization and performance**: Use `t.Run()` for subtests and `t.Parallel()` for parallel tests.

4. **Mocking is done through interfaces**: Go's interface system makes it easy to create mocks for testing.

5. **Benchmarks help identify performance bottlenecks**: Use benchmarks to compare different implementations and optimize your code.

6. **Examples serve as both documentation and tests**: They show how to use your code and verify that it works as expected.

7. **Test coverage helps identify untested code**: Use the `-cover` flag to see how much of your code is covered by tests.

8. **Tests should be independent and repeatable**: Each test should be able to run on its own without depending on other tests or global state.

9. **Follow the Arrange-Act-Assert pattern**: Structure your tests with clear setup, execution, and verification phases.

10. **Use descriptive test names and error messages**: Make it easy to understand what each test is checking and why it failed.

## Next Steps

Now that you understand testing in Go, let's learn about [Reflection and Code Generation](20-reflection-code-generation.md) to explore Go's metaprogramming capabilities!

---

**Previous**: [← Context Package](18-context.md) | **Next**: [Reflection and Code Generation →](20-reflection-code-generation.md)