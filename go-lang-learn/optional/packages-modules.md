# 14. Packages and Modules

Go's package system is fundamental to organizing code and managing dependencies. Go modules, introduced in Go 1.11, provide a modern dependency management system that makes it easy to version and distribute Go code.

## Understanding Packages

### Package Basics

```go
// Package declaration - must be the first non-comment line
package main

import (
    "fmt"     // Standard library package
    "strings" // Standard library package
    "os"      // Standard library package
)

func main() {
    fmt.Println("Hello from main package")
    fmt.Println(strings.ToUpper("hello world"))
    fmt.Println("Args:", os.Args)
}
```

### Package Naming Rules

1. Package names should be lowercase
2. Package names should be short and descriptive
3. Package names should not use underscores or camelCase
4. The package name should match the directory name

### Visibility Rules

```go
// In package "mypackage"
package mypackage

// Exported (public) - starts with uppercase
var PublicVariable = "I am public"
const PublicConstant = 42

type PublicStruct struct {
    PublicField    string // Exported field
    privateField   string // Unexported field
}

func PublicFunction() string {
    return "I am a public function"
}

// Unexported (private) - starts with lowercase
var privateVariable = "I am private"
const privateConstant = 24

type privateStruct struct {
    field string
}

func privateFunction() string {
    return "I am a private function"
}

// Public method on public type
func (p *PublicStruct) PublicMethod() string {
    return p.PublicField + " " + p.privateField
}

// Private method on public type
func (p *PublicStruct) privateMethod() string {
    return "private method"
}
```

## Creating Custom Packages

### Example: Math Utilities Package

```go
// File: mathutils/basic.go
package mathutils

import "math"

// Add returns the sum of two integers
func Add(a, b int) int {
    return a + b
}

// Subtract returns the difference of two integers
func Subtract(a, b int) int {
    return a - b
}

// Multiply returns the product of two integers
func Multiply(a, b int) int {
    return a * b
}

// Divide returns the quotient of two floats and an error if division by zero
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Power returns base raised to the power of exponent
func Power(base, exponent float64) float64 {
    return math.Pow(base, exponent)
}

// internal helper function (unexported)
func abs(x float64) float64 {
    if x < 0 {
        return -x
    }
    return x
}
```

```go
// File: mathutils/advanced.go
package mathutils

import "math"

// Factorial calculates the factorial of n
func Factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * Factorial(n-1)
}

// IsPrime checks if a number is prime
func IsPrime(n int) bool {
    if n <= 1 {
        return false
    }
    if n <= 3 {
        return true
    }
    if n%2 == 0 || n%3 == 0 {
        return false
    }
    
    for i := 5; i*i <= n; i += 6 {
        if n%i == 0 || n%(i+2) == 0 {
            return false
        }
    }
    return true
}

// GCD calculates the greatest common divisor
func GCD(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return abs(a)
}

// LCM calculates the least common multiple
func LCM(a, b int) int {
    return abs(a*b) / GCD(a, b)
}

// SquareRoot calculates square root using Newton's method
func SquareRoot(x float64) float64 {
    if x < 0 {
        return math.NaN()
    }
    if x == 0 {
        return 0
    }
    
    guess := x / 2
    for i := 0; i < 10; i++ {
        guess = (guess + x/guess) / 2
    }
    return guess
}
```

```go
// File: mathutils/constants.go
package mathutils

import "math"

// Mathematical constants
const (
    Pi      = math.Pi
    E       = math.E
    Phi     = 1.618033988749 // Golden ratio
    Sqrt2   = math.Sqrt2
    Sqrt3   = 1.732050807568
    Log2E   = math.Log2E
    Log10E  = math.Log10E
    Ln2     = math.Ln2
    Ln10    = math.Ln10
)

// Conversion factors
const (
    DegreesToRadians = Pi / 180
    RadiansToDegrees = 180 / Pi
    FeetToMeters     = 0.3048
    MetersToFeet     = 1 / FeetToMeters
    MilesToKm        = 1.609344
    KmToMiles        = 1 / MilesToKm
)
```

### Example: String Utilities Package

```go
// File: stringutils/manipulation.go
package stringutils

import (
    "strings"
    "unicode"
)

// Reverse returns the reversed string
func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// IsPalindrome checks if a string is a palindrome
func IsPalindrome(s string) bool {
    s = strings.ToLower(s)
    cleaned := ""
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            cleaned += string(r)
        }
    }
    return cleaned == Reverse(cleaned)
}

// WordCount returns the number of words in a string
func WordCount(s string) int {
    return len(strings.Fields(s))
}

// Capitalize capitalizes the first letter of each word
func Capitalize(s string) string {
    return strings.Title(strings.ToLower(s))
}

// TruncateWithEllipsis truncates string and adds ellipsis if needed
func TruncateWithEllipsis(s string, maxLength int) string {
    if len(s) <= maxLength {
        return s
    }
    if maxLength <= 3 {
        return s[:maxLength]
    }
    return s[:maxLength-3] + "..."
}

// RemoveWhitespace removes all whitespace characters
func RemoveWhitespace(s string) string {
    return strings.ReplaceAll(strings.ReplaceAll(s, " ", ""), "\t", "")
}

// CountOccurrences counts occurrences of substring in string
func CountOccurrences(s, substr string) int {
    return strings.Count(s, substr)
}
```

```go
// File: stringutils/validation.go
package stringutils

import (
    "regexp"
    "strings"
    "unicode"
)

// IsEmail validates email format
func IsEmail(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRegex.MatchString(email)
}

// IsPhoneNumber validates phone number format (simple validation)
func IsPhoneNumber(phone string) bool {
    phoneRegex := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
    cleaned := strings.ReplaceAll(strings.ReplaceAll(phone, "-", ""), " ", "")
    return phoneRegex.MatchString(cleaned)
}

// IsURL validates URL format
func IsURL(url string) bool {
    urlRegex := regexp.MustCompile(`^https?://[\w\-]+(\.[\w\-]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`)
    return urlRegex.MatchString(url)
}

// IsAlphanumeric checks if string contains only alphanumeric characters
func IsAlphanumeric(s string) bool {
    for _, r := range s {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
            return false
        }
    }
    return len(s) > 0
}

// IsNumeric checks if string contains only numeric characters
func IsNumeric(s string) bool {
    for _, r := range s {
        if !unicode.IsDigit(r) {
            return false
        }
    }
    return len(s) > 0
}

// HasMinLength checks if string has minimum length
func HasMinLength(s string, minLength int) bool {
    return len(s) >= minLength
}

// HasMaxLength checks if string has maximum length
func HasMaxLength(s string, maxLength int) bool {
    return len(s) <= maxLength
}

// ContainsUppercase checks if string contains uppercase letters
func ContainsUppercase(s string) bool {
    for _, r := range s {
        if unicode.IsUpper(r) {
            return true
        }
    }
    return false
}

// ContainsLowercase checks if string contains lowercase letters
func ContainsLowercase(s string) bool {
    for _, r := range s {
        if unicode.IsLower(r) {
            return true
        }
    }
    return false
}

// ContainsDigit checks if string contains digits
func ContainsDigit(s string) bool {
    for _, r := range s {
        if unicode.IsDigit(r) {
            return true
        }
    }
    return false
}

// ContainsSpecialChar checks if string contains special characters
func ContainsSpecialChar(s string) bool {
    for _, r := range s {
        if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
            return true
        }
    }
    return false
}
```

## Go Modules

### Creating a Module

```bash
# Initialize a new module
go mod init example.com/myproject

# This creates a go.mod file
```

### go.mod File Structure

```go
// go.mod
module example.com/myproject

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
)

require (
    github.com/gorilla/context v1.1.1 // indirect
)

replace github.com/old/package => github.com/new/package v1.2.3

exclude github.com/broken/package v1.0.0
```

### Module Commands

```bash
# Add a dependency
go get github.com/gorilla/mux

# Add a specific version
go get github.com/gorilla/mux@v1.8.0

# Update dependencies
go get -u
go get -u github.com/gorilla/mux

# Remove unused dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Show module graph
go mod graph

# Show why a dependency is needed
go mod why github.com/gorilla/mux

# Create vendor directory
go mod vendor
```

## Package Import Patterns

### Standard Import

```go
package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

func main() {
    fmt.Println("Standard import")
    http.HandleFunc("/", handler)
    json.Marshal(map[string]string{"key": "value"})
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}
```

### Aliased Import

```go
package main

import (
    "fmt"
    mux "github.com/gorilla/mux"  // Alias
    pq "github.com/lib/pq"        // Alias
)

func main() {
    router := mux.NewRouter()
    fmt.Println("Using aliased imports")
    
    // Use pq alias
    _ = pq.Driver{}
}
```

### Dot Import (Use Sparingly)

```go
package main

import (
    . "fmt"  // Dot import - imports into current namespace
    . "math" // Not recommended for most cases
)

func main() {
    // Can use functions without package prefix
    Println("Hello from dot import")
    result := Sqrt(16)
    Printf("Square root of 16 is %.2f\n", result)
}
```

### Blank Import

```go
package main

import (
    "database/sql"
    _ "github.com/lib/pq" // Blank import for side effects (driver registration)
    "fmt"
)

func main() {
    // The pq driver is registered but we don't use it directly
    db, err := sql.Open("postgres", "connection_string")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    defer db.Close()
}
```

## Package Organization Best Practices

### Project Structure Example

```
myproject/
├── go.mod
├── go.sum
├── main.go
├── README.md
├── cmd/
│   ├── server/
│   │   └── main.go
│   └── client/
│       └── main.go
├── internal/
│   ├── auth/
│   │   ├── auth.go
│   │   └── auth_test.go
│   ├── database/
│   │   ├── connection.go
│   │   ├── migrations/
│   │   └── models/
│   └── handlers/
│       ├── user.go
│       └── product.go
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── utils/
│       ├── string.go
│       └── math.go
├── api/
│   └── openapi.yaml
├── web/
│   ├── static/
│   └── templates/
├── scripts/
│   └── deploy.sh
└── docs/
    └── api.md
```

### Package Documentation

```go
// Package mathutils provides mathematical utility functions and constants.
//
// This package includes basic arithmetic operations, advanced mathematical
// functions, and commonly used mathematical constants.
//
// Example usage:
//
//     result := mathutils.Add(5, 3)
//     fmt.Println(result) // Output: 8
//
//     isPrime := mathutils.IsPrime(17)
//     fmt.Println(isPrime) // Output: true
//
package mathutils

import "errors"

// ErrDivisionByZero is returned when attempting to divide by zero.
var ErrDivisionByZero = errors.New("division by zero")

// Calculator represents a mathematical calculator with memory.
type Calculator struct {
    memory float64
}

// NewCalculator creates a new Calculator instance.
//
// The calculator starts with zero in memory.
func NewCalculator() *Calculator {
    return &Calculator{memory: 0}
}

// Add performs addition and returns the result.
//
// Parameters:
//   - a: first operand
//   - b: second operand
//
// Returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}

// Divide performs division with error handling.
//
// Parameters:
//   - dividend: the number to be divided
//   - divisor: the number to divide by
//
// Returns:
//   - float64: the quotient
//   - error: ErrDivisionByZero if divisor is zero
//
// Example:
//
//     result, err := mathutils.Divide(10, 2)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println(result) // Output: 5
//
func Divide(dividend, divisor float64) (float64, error) {
    if divisor == 0 {
        return 0, ErrDivisionByZero
    }
    return dividend / divisor, nil
}

// StoreInMemory stores a value in the calculator's memory.
func (c *Calculator) StoreInMemory(value float64) {
    c.memory = value
}

// RecallFromMemory returns the value stored in memory.
func (c *Calculator) RecallFromMemory() float64 {
    return c.memory
}

// ClearMemory resets the memory to zero.
func (c *Calculator) ClearMemory() {
    c.memory = 0
}
```

## Practical Examples

### Example 1: Complete Web Server Package

```go
// File: webserver/server.go
package webserver

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

// Server represents an HTTP server with graceful shutdown
type Server struct {
    httpServer *http.Server
    logger     *log.Logger
}

// Config holds server configuration
type Config struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

// DefaultConfig returns a default server configuration
func DefaultConfig() Config {
    return Config{
        Host:         "localhost",
        Port:         8080,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
}

// NewServer creates a new server instance
func NewServer(config Config, handler http.Handler) *Server {
    addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
    
    httpServer := &http.Server{
        Addr:         addr,
        Handler:      handler,
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
        IdleTimeout:  config.IdleTimeout,
    }
    
    logger := log.New(os.Stdout, "[SERVER] ", log.LstdFlags)
    
    return &Server{
        httpServer: httpServer,
        logger:     logger,
    }
}

// Start starts the server with graceful shutdown
func (s *Server) Start() error {
    // Start server in a goroutine
    go func() {
        s.logger.Printf("Starting server on %s", s.httpServer.Addr)
        if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            s.logger.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Wait for interrupt signal to gracefully shutdown
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    s.logger.Println("Shutting down server...")
    
    // Create a deadline for shutdown
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Attempt graceful shutdown
    if err := s.httpServer.Shutdown(ctx); err != nil {
        s.logger.Printf("Server forced to shutdown: %v", err)
        return err
    }
    
    s.logger.Println("Server exited")
    return nil
}

// Stop stops the server immediately
func (s *Server) Stop() error {
    return s.httpServer.Close()
}
```

```go
// File: webserver/middleware.go
package webserver

import (
    "log"
    "net/http"
    "time"
)

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Create a response writer wrapper to capture status code
        wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapper, r)
        
        duration := time.Since(start)
        log.Printf("%s %s %d %v", r.Method, r.URL.Path, wrapper.statusCode, duration)
    })
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// CORSMiddleware adds CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("Panic recovered: %v", err)
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// ChainMiddleware chains multiple middleware functions
func ChainMiddleware(handler http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
```

```go
// File: webserver/handlers.go
package webserver

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Response represents a JSON response
type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

// HealthHandler returns server health status
func HealthHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Success: true,
        Message: "Server is healthy",
        Data: map[string]interface{}{
            "timestamp": time.Now().UTC(),
            "status":    "ok",
        },
    }
    
    writeJSONResponse(w, http.StatusOK, response)
}

// HomeHandler returns a welcome message
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Success: true,
        Message: "Welcome to the API",
        Data: map[string]interface{}{
            "version": "1.0.0",
            "endpoints": []string{
                "/health",
                "/api/users",
                "/api/products",
            },
        },
    }
    
    writeJSONResponse(w, http.StatusOK, response)
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    response := Response{
        Success: false,
        Error:   "Endpoint not found",
    }
    
    writeJSONResponse(w, http.StatusNotFound, response)
}

// writeJSONResponse writes a JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(data); err != nil {
        http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
    }
}

// CreateRouter creates and configures the main router
func CreateRouter() http.Handler {
    mux := http.NewServeMux()
    
    // Register routes
    mux.HandleFunc("/", HomeHandler)
    mux.HandleFunc("/health", HealthHandler)
    
    // Apply middleware
    handler := ChainMiddleware(
        mux,
        LoggingMiddleware,
        CORSMiddleware,
        RecoveryMiddleware,
    )
    
    return handler
}
```

### Example 2: Using the Web Server Package

```go
// File: main.go
package main

import (
    "log"
    "time"
    
    "./webserver" // Local package import
)

func main() {
    // Create server configuration
    config := webserver.Config{
        Host:         "localhost",
        Port:         8080,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  30 * time.Second,
    }
    
    // Create router with handlers
    router := webserver.CreateRouter()
    
    // Create and start server
    server := webserver.NewServer(config, router)
    
    log.Println("Starting web server...")
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

### Example 3: Database Package

```go
// File: database/connection.go
package database

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database configuration
type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// DB wraps sql.DB with additional functionality
type DB struct {
    *sql.DB
    config Config
}

// Connect establishes a database connection
func Connect(config Config) (*DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
    
    sqlDB, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(5)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{
        DB:     sqlDB,
        config: config,
    }, nil
}

// HealthCheck checks database connectivity
func (db *DB) HealthCheck() error {
    return db.Ping()
}

// GetStats returns database statistics
func (db *DB) GetStats() sql.DBStats {
    return db.Stats()
}
```

```go
// File: database/models.go
package database

import (
    "database/sql"
    "time"
)

// User represents a user in the database
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository handles user database operations
type UserRepository struct {
    db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
    return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *User) error {
    query := `
        INSERT INTO users (username, email, created_at, updated_at)
        VALUES ($1, $2, $3, $4)
        RETURNING id`
    
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now
    
    err := r.db.QueryRow(query, user.Username, user.Email, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*User, error) {
    query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`
    
    user := &User{}
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err == sql.ErrNoRows {
        return nil, fmt.Errorf("user with id %d not found", id)
    }
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

// GetAll retrieves all users
func (r *UserRepository) GetAll() ([]*User, error) {
    query := `SELECT id, username, email, created_at, updated_at FROM users ORDER BY id`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query users: %w", err)
    }
    defer rows.Close()
    
    var users []*User
    for rows.Next() {
        user := &User{}
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, user)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating users: %w", err)
    }
    
    return users, nil
}

// Update updates a user
func (r *UserRepository) Update(user *User) error {
    query := `
        UPDATE users 
        SET username = $1, email = $2, updated_at = $3
        WHERE id = $4`
    
    user.UpdatedAt = time.Now()
    
    result, err := r.db.Exec(query, user.Username, user.Email, user.UpdatedAt, user.ID)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user with id %d not found", user.ID)
    }
    
    return nil
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
    query := `DELETE FROM users WHERE id = $1`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("user with id %d not found", id)
    }
    
    return nil
}
```

## Testing Packages

### Package Testing Example

```go
// File: mathutils/mathutils_test.go
package mathutils

import (
    "math"
    "testing"
)

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed numbers", -2, 3, 1},
        {"zero", 0, 5, 5},
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

func TestDivide(t *testing.T) {
    // Test successful division
    result, err := Divide(10, 2)
    if err != nil {
        t.Errorf("Divide(10, 2) returned error: %v", err)
    }
    if result != 5 {
        t.Errorf("Divide(10, 2) = %f; want 5", result)
    }
    
    // Test division by zero
    _, err = Divide(10, 0)
    if err == nil {
        t.Error("Divide(10, 0) should return error")
    }
}

func TestIsPrime(t *testing.T) {
    tests := []struct {
        n        int
        expected bool
    }{
        {2, true},
        {3, true},
        {4, false},
        {17, true},
        {25, false},
        {1, false},
        {0, false},
        {-5, false},
    }
    
    for _, tt := range tests {
        t.Run(fmt.Sprintf("IsPrime(%d)", tt.n), func(t *testing.T) {
            result := IsPrime(tt.n)
            if result != tt.expected {
                t.Errorf("IsPrime(%d) = %t; want %t", tt.n, result, tt.expected)
            }
        })
    }
}

func BenchmarkFactorial(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Factorial(10)
    }
}

func BenchmarkIsPrime(b *testing.B) {
    for i := 0; i < b.N; i++ {
        IsPrime(97)
    }
}

func ExampleAdd() {
    result := Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

func ExampleDivide() {
    result, err := Divide(10, 2)
    if err != nil {
        fmt.Printf("Error: %v", err)
        return
    }
    fmt.Printf("%.1f", result)
    // Output: 5.0
}
```

## Exercises

### Exercise 1: Logger Package
Create a comprehensive logging package.

```go
// Create a logger package with the following features:
// - Different log levels (DEBUG, INFO, WARN, ERROR, FATAL)
// - Multiple output destinations (console, file, network)
// - Structured logging (JSON format)
// - Log rotation
// - Context-aware logging

package logger

// Define interfaces and types:
// Logger interface
// Level type
// Output interface
// Config struct

// Implement:
// NewLogger(config Config) Logger
// SetLevel(level Level)
// Debug(msg string, fields ...Field)
// Info(msg string, fields ...Field)
// Warn(msg string, fields ...Field)
// Error(msg string, fields ...Field)
// Fatal(msg string, fields ...Field)
// WithContext(ctx context.Context) Logger
// WithFields(fields ...Field) Logger
```

### Exercise 2: Configuration Package
Create a configuration management package.

```go
// Create a config package that:
// - Loads configuration from multiple sources (files, environment, flags)
// - Supports multiple formats (JSON, YAML, TOML)
// - Provides validation
// - Supports hot reloading
// - Has default values

package config

// Define:
// Config interface
// Source interface
// Validator interface
// Loader struct

// Implement:
// NewLoader() *Loader
// AddSource(source Source) *Loader
// AddValidator(validator Validator) *Loader
// Load(target interface{}) error
// Watch(callback func(interface{})) error
// Get(key string) interface{}
// Set(key string, value interface{}) error
```

### Exercise 3: Cache Package
Create a caching package with multiple backends.

```go
// Create a cache package that:
// - Supports multiple backends (memory, Redis, file)
// - Has TTL support
// - Provides cache statistics
// - Supports cache warming
// - Has eviction policies

package cache

// Define:
// Cache interface
// Backend interface
// Stats struct
// Config struct

// Implement:
// NewCache(backend Backend, config Config) Cache
// Get(key string) (interface{}, bool)
// Set(key string, value interface{}, ttl time.Duration) error
// Delete(key string) error
// Clear() error
// Stats() Stats
// Warm(data map[string]interface{}) error
```

## Key Takeaways

1. **Package naming** should be lowercase, short, and descriptive
2. **Visibility** is controlled by capitalization (exported vs unexported)
3. **Go modules** provide modern dependency management
4. **Package documentation** should explain purpose and provide examples
5. **Import patterns** include standard, aliased, dot, and blank imports
6. **Project structure** should separate concerns logically
7. **Internal packages** are only accessible within the same module
8. **Package initialization** happens through `init()` functions
9. **Testing** should be in the same package or `_test` package
10. **Vendor directory** can be used for dependency vendoring

## Next Steps

Now that you understand packages and modules, let's explore [Concurrency with Goroutines](15-goroutines.md) to learn about Go's powerful concurrency model!

---

**Previous**: [← Error Handling](13-error-handling.md) | **Next**: [Goroutines →](15-goroutines.md)