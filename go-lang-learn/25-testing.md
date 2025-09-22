# 23. Testing

Go has excellent built-in support for testing with the `testing` package. This chapter covers unit testing, benchmarking, table-driven tests, test coverage, and testing best practices.

## Basic Testing

### Writing Your First Test

```go
// math_utils.go
package main

import "math"

func Add(a, b int) int {
    return a + b
}

func Multiply(a, b int) int {
    return a * b
}

func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

func IsPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}
```

```go
// math_utils_test.go
package main

import (
    "testing"
)

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    expected := 5
    
    if result != expected {
        t.Errorf("Add(2, 3) = %d; want %d", result, expected)
    }
}

func TestMultiply(t *testing.T) {
    result := Multiply(4, 5)
    expected := 20
    
    if result != expected {
        t.Errorf("Multiply(4, 5) = %d; want %d", result, expected)
    }
}

func TestDivide(t *testing.T) {
    // Test normal division
    result, err := Divide(10, 2)
    if err != nil {
        t.Errorf("Divide(10, 2) returned error: %v", err)
    }
    
    expected := 5.0
    if result != expected {
        t.Errorf("Divide(10, 2) = %f; want %f", result, expected)
    }
}

func TestDivideByZero(t *testing.T) {
    _, err := Divide(10, 0)
    if err == nil {
        t.Error("Divide(10, 0) should return an error")
    }
}

func TestIsPrime(t *testing.T) {
    // Test prime numbers
    primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23}
    for _, prime := range primes {
        if !IsPrime(prime) {
            t.Errorf("IsPrime(%d) = false; want true", prime)
        }
    }
    
    // Test non-prime numbers
    nonPrimes := []int{0, 1, 4, 6, 8, 9, 10, 12, 14, 15}
    for _, nonPrime := range nonPrimes {
        if IsPrime(nonPrime) {
            t.Errorf("IsPrime(%d) = true; want false", nonPrime)
        }
    }
}
```

### Running Tests

```bash
# Run all tests in current package
go test

# Run tests with verbose output
go test -v

# Run specific test
go test -run TestAdd

# Run tests matching pattern
go test -run "Test.*Divide"

# Run tests in all subdirectories
go test ./...
```

## Table-Driven Tests

### Basic Table-Driven Test

```go
package main

import (
    "testing"
)

func TestAddTableDriven(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed numbers", -2, 3, 1},
        {"zero values", 0, 0, 0},
        {"large numbers", 1000000, 2000000, 3000000},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}

func TestDivideTableDriven(t *testing.T) {
    tests := []struct {
        name        string
        a, b        float64
        expected    float64
        expectError bool
    }{
        {"normal division", 10, 2, 5, false},
        {"division by zero", 10, 0, 0, true},
        {"negative numbers", -10, 2, -5, false},
        {"decimal result", 7, 2, 3.5, false},
        {"zero dividend", 0, 5, 0, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.a, tt.b)
            
            if tt.expectError {
                if err == nil {
                    t.Errorf("Divide(%f, %f) should return error", tt.a, tt.b)
                }
                return
            }
            
            if err != nil {
                t.Errorf("Divide(%f, %f) returned unexpected error: %v", tt.a, tt.b, err)
                return
            }
            
            if result != tt.expected {
                t.Errorf("Divide(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Advanced Table-Driven Tests

```go
package main

import (
    "fmt"
    "strings"
    "testing"
)

// String utilities to test
func ReverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func CountWords(s string) int {
    if strings.TrimSpace(s) == "" {
        return 0
    }
    return len(strings.Fields(s))
}

func IsPalindrome(s string) bool {
    s = strings.ToLower(strings.ReplaceAll(s, " ", ""))
    return s == ReverseString(s)
}

func TestStringUtilities(t *testing.T) {
    tests := []struct {
        name     string
        function string
        input    string
        expected interface{}
    }{
        // ReverseString tests
        {"reverse empty string", "reverse", "", ""},
        {"reverse single char", "reverse", "a", "a"},
        {"reverse simple word", "reverse", "hello", "olleh"},
        {"reverse with spaces", "reverse", "hello world", "dlrow olleh"},
        {"reverse unicode", "reverse", "🙂🙃", "🙃🙂"},
        
        // CountWords tests
        {"count empty string", "count", "", 0},
        {"count single word", "count", "hello", 1},
        {"count multiple words", "count", "hello world", 2},
        {"count with extra spaces", "count", "  hello   world  ", 2},
        {"count only spaces", "count", "   ", 0},
        
        // IsPalindrome tests
        {"palindrome empty", "palindrome", "", true},
        {"palindrome single char", "palindrome", "a", true},
        {"palindrome simple", "palindrome", "racecar", true},
        {"palindrome with spaces", "palindrome", "race car", true},
        {"not palindrome", "palindrome", "hello", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            switch tt.function {
            case "reverse":
                result := ReverseString(tt.input)
                if result != tt.expected.(string) {
                    t.Errorf("ReverseString(%q) = %q; want %q", tt.input, result, tt.expected)
                }
                
            case "count":
                result := CountWords(tt.input)
                if result != tt.expected.(int) {
                    t.Errorf("CountWords(%q) = %d; want %d", tt.input, result, tt.expected)
                }
                
            case "palindrome":
                result := IsPalindrome(tt.input)
                if result != tt.expected.(bool) {
                    t.Errorf("IsPalindrome(%q) = %t; want %t", tt.input, result, tt.expected)
                }
            }
        })
    }
}
```

## Test Helpers and Utilities

### Custom Test Helpers

```go
package main

import (
    "reflect"
    "testing"
)

// Helper function to check equality
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper() // This marks the function as a test helper
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v, want %v", got, want)
    }
}

// Helper function to check errors
func assertError(t *testing.T, err error, wantError bool) {
    t.Helper()
    if wantError && err == nil {
        t.Error("expected error but got none")
    }
    if !wantError && err != nil {
        t.Errorf("unexpected error: %v", err)
    }
}

// Helper function to check if slice contains element
func assertContains(t *testing.T, slice []string, element string) {
    t.Helper()
    for _, item := range slice {
        if item == element {
            return
        }
    }
    t.Errorf("slice %v does not contain %q", slice, element)
}

// Example usage of helpers
func TestWithHelpers(t *testing.T) {
    // Test Add function
    result := Add(2, 3)
    assertEqual(t, result, 5)
    
    // Test Divide function
    quotient, err := Divide(10, 2)
    assertError(t, err, false)
    assertEqual(t, quotient, 5.0)
    
    // Test error case
    _, err = Divide(10, 0)
    assertError(t, err, true)
    
    // Test slice contains
    fruits := []string{"apple", "banana", "orange"}
    assertContains(t, fruits, "banana")
}
```

### Setup and Teardown

```go
package main

import (
    "fmt"
    "os"
    "testing"
    "time"
)

// Global test setup
func TestMain(m *testing.M) {
    fmt.Println("Setting up tests...")
    
    // Setup code here
    setupDatabase()
    setupTempFiles()
    
    // Run tests
    code := m.Run()
    
    // Cleanup code here
    fmt.Println("Cleaning up tests...")
    cleanupDatabase()
    cleanupTempFiles()
    
    os.Exit(code)
}

func setupDatabase() {
    fmt.Println("Setting up test database...")
    // Database setup logic
}

func cleanupDatabase() {
    fmt.Println("Cleaning up test database...")
    // Database cleanup logic
}

func setupTempFiles() {
    fmt.Println("Creating temporary files...")
    // File setup logic
}

func cleanupTempFiles() {
    fmt.Println("Removing temporary files...")
    // File cleanup logic
}

// Test with setup and teardown
func TestWithSetupTeardown(t *testing.T) {
    // Setup for this specific test
    tempData := "test data"
    defer func() {
        // Teardown for this specific test
        fmt.Println("Test-specific cleanup")
    }()
    
    // Test logic
    if tempData == "" {
        t.Error("tempData should not be empty")
    }
}
```

## Benchmarking

### Basic Benchmarks

```go
package main

import (
    "fmt"
    "strings"
    "testing"
)

// Functions to benchmark
func ConcatenateStrings(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s
    }
    return result
}

func ConcatenateStringsBuilder(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

func ConcatenateStringsJoin(strs []string) string {
    return strings.Join(strs, "")
}

// Benchmark tests
func BenchmarkConcatenateStrings(b *testing.B) {
    strs := []string{"hello", "world", "this", "is", "a", "test"}
    
    for i := 0; i < b.N; i++ {
        ConcatenateStrings(strs)
    }
}

func BenchmarkConcatenateStringsBuilder(b *testing.B) {
    strs := []string{"hello", "world", "this", "is", "a", "test"}
    
    for i := 0; i < b.N; i++ {
        ConcatenateStringsBuilder(strs)
    }
}

func BenchmarkConcatenateStringsJoin(b *testing.B) {
    strs := []string{"hello", "world", "this", "is", "a", "test"}
    
    for i := 0; i < b.N; i++ {
        ConcatenateStringsJoin(strs)
    }
}

// Benchmark with different input sizes
func BenchmarkConcatenateStringsSizes(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        strs := make([]string, size)
        for i := range strs {
            strs[i] = fmt.Sprintf("string%d", i)
        }
        
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ConcatenateStringsBuilder(strs)
            }
        })
    }
}

// Benchmark with memory allocation tracking
func BenchmarkWithMemory(b *testing.B) {
    b.ReportAllocs() // Report memory allocations
    
    for i := 0; i < b.N; i++ {
        strs := []string{"hello", "world", "test"}
        ConcatenateStringsBuilder(strs)
    }
}
```

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=.

# Run specific benchmark
go test -bench=BenchmarkConcatenateStrings

# Run benchmarks with memory stats
go test -bench=. -benchmem

# Run benchmarks multiple times for accuracy
go test -bench=. -count=5

# Set benchmark time
go test -bench=. -benchtime=10s
```

## Test Coverage

### Measuring Coverage

```bash
# Run tests with coverage
go test -cover

# Generate detailed coverage report
go test -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Show coverage by function
go tool cover -func=coverage.out
```

### Coverage Example

```go
// calculator.go
package main

import "errors"

type Calculator struct {
    memory float64
}

func (c *Calculator) Add(a, b float64) float64 {
    result := a + b
    c.memory = result
    return result
}

func (c *Calculator) Subtract(a, b float64) float64 {
    result := a - b
    c.memory = result
    return result
}

func (c *Calculator) Multiply(a, b float64) float64 {
    result := a * b
    c.memory = result
    return result
}

func (c *Calculator) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    result := a / b
    c.memory = result
    return result, nil
}

func (c *Calculator) GetMemory() float64 {
    return c.memory
}

func (c *Calculator) ClearMemory() {
    c.memory = 0
}
```

```go
// calculator_test.go
package main

import "testing"

func TestCalculator(t *testing.T) {
    calc := &Calculator{}
    
    // Test Add
    result := calc.Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %f; want 5", result)
    }
    
    // Test memory
    if calc.GetMemory() != 5 {
        t.Errorf("Memory = %f; want 5", calc.GetMemory())
    }
    
    // Test Subtract
    result = calc.Subtract(10, 3)
    if result != 7 {
        t.Errorf("Subtract(10, 3) = %f; want 7", result)
    }
    
    // Test Multiply
    result = calc.Multiply(4, 5)
    if result != 20 {
        t.Errorf("Multiply(4, 5) = %f; want 20", result)
    }
    
    // Test Divide
    result, err := calc.Divide(20, 4)
    if err != nil {
        t.Errorf("Divide(20, 4) returned error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(20, 4) = %f; want 5", result)
    }
    
    // Test divide by zero
    _, err = calc.Divide(10, 0)
    if err == nil {
        t.Error("Divide(10, 0) should return error")
    }
    
    // Test clear memory
    calc.ClearMemory()
    if calc.GetMemory() != 0 {
        t.Errorf("Memory after clear = %f; want 0", calc.GetMemory())
    }
}
```

## Mocking and Test Doubles

### Interface-Based Mocking

```go
package main

import (
    "fmt"
    "testing"
)

// Interface for external dependency
type EmailSender interface {
    SendEmail(to, subject, body string) error
}

// Real implementation
type SMTPSender struct {
    host string
    port int
}

func (s *SMTPSender) SendEmail(to, subject, body string) error {
    // Real email sending logic
    fmt.Printf("Sending email to %s via SMTP\n", to)
    return nil
}

// Service that uses the dependency
type UserService struct {
    emailSender EmailSender
}

func (u *UserService) RegisterUser(email, name string) error {
    // Registration logic
    fmt.Printf("Registering user: %s\n", name)
    
    // Send welcome email
    subject := "Welcome!"
    body := fmt.Sprintf("Hello %s, welcome to our service!", name)
    
    return u.emailSender.SendEmail(email, subject, body)
}

// Mock implementation for testing
type MockEmailSender struct {
    sentEmails []Email
    shouldFail bool
}

type Email struct {
    To      string
    Subject string
    Body    string
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
    if m.shouldFail {
        return fmt.Errorf("mock email sending failed")
    }
    
    m.sentEmails = append(m.sentEmails, Email{
        To:      to,
        Subject: subject,
        Body:    body,
    })
    
    return nil
}

// Test with mock
func TestUserService_RegisterUser(t *testing.T) {
    tests := []struct {
        name        string
        email       string
        userName    string
        shouldFail  bool
        expectError bool
    }{
        {"successful registration", "test@example.com", "John Doe", false, false},
        {"email sending fails", "test@example.com", "Jane Doe", true, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock
            mockSender := &MockEmailSender{
                shouldFail: tt.shouldFail,
            }
            
            // Create service with mock
            service := &UserService{
                emailSender: mockSender,
            }
            
            // Test the service
            err := service.RegisterUser(tt.email, tt.userName)
            
            // Check error expectation
            if tt.expectError && err == nil {
                t.Error("expected error but got none")
            }
            if !tt.expectError && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            
            // Check email was sent (if no error expected)
            if !tt.expectError {
                if len(mockSender.sentEmails) != 1 {
                    t.Errorf("expected 1 email, got %d", len(mockSender.sentEmails))
                }
                
                email := mockSender.sentEmails[0]
                if email.To != tt.email {
                    t.Errorf("email to = %s; want %s", email.To, tt.email)
                }
                if email.Subject != "Welcome!" {
                    t.Errorf("email subject = %s; want Welcome!", email.Subject)
                }
            }
        })
    }
}
```

### HTTP Testing

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

type UserHandler struct {
    users []User
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        h.getUsers(w, r)
    case http.MethodPost:
        h.createUser(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func (h *UserHandler) getUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(h.users)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user.ID = len(h.users) + 1
    h.users = append(h.users, user)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func TestUserHandler(t *testing.T) {
    handler := &UserHandler{
        users: []User{
            {ID: 1, Name: "John Doe", Email: "john@example.com"},
        },
    }
    
    t.Run("GET users", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodGet, "/users", nil)
        w := httptest.NewRecorder()
        
        handler.ServeHTTP(w, req)
        
        if w.Code != http.StatusOK {
            t.Errorf("status code = %d; want %d", w.Code, http.StatusOK)
        }
        
        var users []User
        if err := json.NewDecoder(w.Body).Decode(&users); err != nil {
            t.Errorf("failed to decode response: %v", err)
        }
        
        if len(users) != 1 {
            t.Errorf("got %d users; want 1", len(users))
        }
    })
    
    t.Run("POST user", func(t *testing.T) {
        userJSON := `{"name":"Jane Doe","email":"jane@example.com"}`
        req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        
        handler.ServeHTTP(w, req)
        
        if w.Code != http.StatusCreated {
            t.Errorf("status code = %d; want %d", w.Code, http.StatusCreated)
        }
        
        var user User
        if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
            t.Errorf("failed to decode response: %v", err)
        }
        
        if user.Name != "Jane Doe" {
            t.Errorf("user name = %s; want Jane Doe", user.Name)
        }
        
        if user.ID != 2 {
            t.Errorf("user ID = %d; want 2", user.ID)
        }
    })
    
    t.Run("Invalid method", func(t *testing.T) {
        req := httptest.NewRequest(http.MethodDelete, "/users", nil)
        w := httptest.NewRecorder()
        
        handler.ServeHTTP(w, req)
        
        if w.Code != http.StatusMethodNotAllowed {
            t.Errorf("status code = %d; want %d", w.Code, http.StatusMethodNotAllowed)
        }
    })
}
```

## Testing Best Practices

### Test Organization

```go
package main

import (
    "testing"
    "time"
)

// Good: Clear test names that describe what is being tested
func TestUserService_CreateUser_WithValidData_ReturnsUser(t *testing.T) {
    // Test implementation
}

func TestUserService_CreateUser_WithInvalidEmail_ReturnsError(t *testing.T) {
    // Test implementation
}

// Good: Group related tests using subtests
func TestUserValidation(t *testing.T) {
    t.Run("valid email", func(t *testing.T) {
        // Test valid email
    })
    
    t.Run("invalid email format", func(t *testing.T) {
        // Test invalid email
    })
    
    t.Run("empty email", func(t *testing.T) {
        // Test empty email
    })
}

// Good: Test edge cases and error conditions
func TestDivision(t *testing.T) {
    tests := []struct {
        name        string
        dividend    float64
        divisor     float64
        expected    float64
        expectError bool
    }{
        {"normal division", 10, 2, 5, false},
        {"division by zero", 10, 0, 0, true},
        {"negative numbers", -10, 2, -5, false},
        {"very small numbers", 0.000001, 0.000001, 1, false},
        {"very large numbers", 1e10, 1e5, 1e5, false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := Divide(tt.dividend, tt.divisor)
            
            if tt.expectError {
                if err == nil {
                    t.Error("expected error but got none")
                }
                return
            }
            
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            
            if result != tt.expected {
                t.Errorf("got %f, want %f", result, tt.expected)
            }
        })
    }
}

// Good: Test timeouts and cancellation
func TestWithTimeout(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping timeout test in short mode")
    }
    
    start := time.Now()
    
    // Test that should complete quickly
    done := make(chan bool)
    go func() {
        time.Sleep(100 * time.Millisecond)
        done <- true
    }()
    
    select {
    case <-done:
        elapsed := time.Since(start)
        if elapsed > 200*time.Millisecond {
            t.Errorf("operation took too long: %v", elapsed)
        }
    case <-time.After(1 * time.Second):
        t.Error("operation timed out")
    }
}
```

### Test Data Management

```go
package main

import (
    "testing"
)

// Good: Use test fixtures for complex data
func getTestUsers() []User {
    return []User{
        {ID: 1, Name: "John Doe", Email: "john@example.com"},
        {ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
        {ID: 3, Name: "Bob Johnson", Email: "bob@example.com"},
    }
}

// Good: Builder pattern for test data
type UserBuilder struct {
    user User
}

func NewUserBuilder() *UserBuilder {
    return &UserBuilder{
        user: User{
            ID:    1,
            Name:  "Test User",
            Email: "test@example.com",
        },
    }
}

func (b *UserBuilder) WithID(id int) *UserBuilder {
    b.user.ID = id
    return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
    b.user.Name = name
    return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
    b.user.Email = email
    return b
}

func (b *UserBuilder) Build() User {
    return b.user
}

func TestUserBuilder(t *testing.T) {
    user := NewUserBuilder().
        WithID(42).
        WithName("Custom User").
        WithEmail("custom@example.com").
        Build()
    
    if user.ID != 42 {
        t.Errorf("user ID = %d; want 42", user.ID)
    }
    
    if user.Name != "Custom User" {
        t.Errorf("user name = %s; want Custom User", user.Name)
    }
}
```

## Integration Testing

### Database Integration Tests

```go
package main

import (
    "database/sql"
    "fmt"
    "testing"
    
    _ "github.com/mattn/go-sqlite3"
)

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Create(user User) error {
    query := `INSERT INTO users (name, email) VALUES (?, ?)`
    _, err := r.db.Exec(query, user.Name, user.Email)
    return err
}

func (r *UserRepository) GetByID(id int) (*User, error) {
    query := `SELECT id, name, email FROM users WHERE id = ?`
    row := r.db.QueryRow(query, id)
    
    var user User
    err := row.Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("failed to open database: %v", err)
    }
    
    // Create table
    query := `
        CREATE TABLE users (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            email TEXT NOT NULL UNIQUE
        )
    `
    
    if _, err := db.Exec(query); err != nil {
        t.Fatalf("failed to create table: %v", err)
    }
    
    return db
}

func TestUserRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test in short mode")
    }
    
    db := setupTestDB(t)
    defer db.Close()
    
    repo := &UserRepository{db: db}
    
    t.Run("create and retrieve user", func(t *testing.T) {
        user := User{
            Name:  "Integration Test User",
            Email: "integration@example.com",
        }
        
        // Create user
        err := repo.Create(user)
        if err != nil {
            t.Fatalf("failed to create user: %v", err)
        }
        
        // Retrieve user
        retrieved, err := repo.GetByID(1)
        if err != nil {
            t.Fatalf("failed to retrieve user: %v", err)
        }
        
        if retrieved.Name != user.Name {
            t.Errorf("name = %s; want %s", retrieved.Name, user.Name)
        }
        
        if retrieved.Email != user.Email {
            t.Errorf("email = %s; want %s", retrieved.Email, user.Email)
        }
    })
}
```

## Exercises

### Exercise 1: Testing a Calculator
Create a comprehensive test suite for a calculator that supports:
- Basic arithmetic operations
- Memory functions (store, recall, clear)
- History of operations
- Error handling for invalid operations

### Exercise 2: HTTP API Testing
Build tests for a REST API that manages a todo list:
- Test all CRUD operations
- Test error cases (invalid data, not found, etc.)
- Test authentication and authorization
- Include integration tests with a real database

### Exercise 3: Concurrent Code Testing
Test a concurrent worker pool system:
- Test that workers process jobs correctly
- Test graceful shutdown
- Test error handling and recovery
- Benchmark different pool sizes

## Key Takeaways

- Go's testing package provides excellent built-in testing capabilities
- Table-driven tests make it easy to test multiple scenarios
- Use `t.Helper()` in test helper functions for better error reporting
- Benchmarks help identify performance bottlenecks
- Test coverage helps identify untested code paths
- Mocking enables testing of code with external dependencies
- Integration tests verify that components work together correctly
- Good test organization and naming improve maintainability
- Test edge cases and error conditions, not just happy paths

## Next Steps

Next, we'll explore [Debugging and Profiling](24-debugging-profiling.md) to learn about Go's debugging tools and performance profiling techniques.