# 16. Packages and Modules

Go's package system is fundamental to organizing and sharing code. Packages provide modularity and encapsulation, while Go modules (introduced in Go 1.11) provide dependency management and versioning.

## Understanding Packages

### What is a Package?

A package is a collection of Go source files in the same directory that are compiled together. Every Go program is made up of packages, and programs start running in package `main`.

### Package Declaration

```go
// Every Go file starts with a package declaration
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Package Naming Rules

1. Package names should be lowercase
2. Use short, concise names
3. Avoid underscores or mixed caps
4. Package name should match the directory name

```go
// Good package names
package http
package json
package strings
package math

// Avoid
package HTTP
package json_parser
package stringUtils
```

## Visibility and Exports

### Exported vs Unexported Identifiers

In Go, visibility is determined by the first letter of an identifier:
- **Exported**: Starts with uppercase letter (public)
- **Unexported**: Starts with lowercase letter (private)

```go
package calculator

// Exported function (can be used by other packages)
func Add(a, b int) int {
    return a + b
}

// Unexported function (only available within this package)
func multiply(a, b int) int {
    return a * b
}

// Exported variable
var MaxValue = 1000

// Unexported variable
var defaultTimeout = 30

// Exported type
type Calculator struct {
    Name string // Exported field
    version int // Unexported field
}

// Exported method
func (c Calculator) GetName() string {
    return c.Name
}

// Unexported method
func (c Calculator) getVersion() int {
    return c.version
}
```

## Creating Custom Packages

### Example: Math Utilities Package

Let's create a custom math utilities package:

**Directory structure:**
```
myproject/
├── main.go
└── mathutils/
    ├── basic.go
    ├── advanced.go
    └── constants.go
```

**mathutils/constants.go:**
```go
package mathutils

import "math"

// Exported constants
const (
    Pi     = math.Pi
    E      = math.E
    Golden = 1.618033988749
)

// Unexported constants
const (
    precision = 1e-9
    maxIterations = 1000
)
```

**mathutils/basic.go:**
```go
package mathutils

import (
    "errors"
    "math"
)

// Add returns the sum of two numbers
func Add(a, b float64) float64 {
    return a + b
}

// Subtract returns the difference of two numbers
func Subtract(a, b float64) float64 {
    return a - b
}

// Multiply returns the product of two numbers
func Multiply(a, b float64) float64 {
    return a * b
}

// Divide returns the quotient of two numbers
func Divide(a, b float64) (float64, error) {
    if isZero(b) {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Power returns a raised to the power of b
func Power(a, b float64) float64 {
    return math.Pow(a, b)
}

// Sqrt returns the square root of a number
func Sqrt(a float64) (float64, error) {
    if a < 0 {
        return 0, errors.New("square root of negative number")
    }
    return math.Sqrt(a), nil
}

// isZero checks if a number is effectively zero
func isZero(a float64) bool {
    return math.Abs(a) < precision
}

// Abs returns the absolute value of a number
func Abs(a float64) float64 {
    return math.Abs(a)
}

// Min returns the minimum of two numbers
func Min(a, b float64) float64 {
    return math.Min(a, b)
}

// Max returns the maximum of two numbers
func Max(a, b float64) float64 {
    return math.Max(a, b)
}
```

**mathutils/advanced.go:**
```go
package mathutils

import (
    "errors"
    "math"
)

// Factorial calculates the factorial of n
func Factorial(n int) (int64, error) {
    if n < 0 {
        return 0, errors.New("factorial of negative number")
    }
    if n == 0 || n == 1 {
        return 1, nil
    }
    
    result := int64(1)
    for i := 2; i <= n; i++ {
        result *= int64(i)
    }
    return result, nil
}

// Fibonacci returns the nth Fibonacci number
func Fibonacci(n int) (int64, error) {
    if n < 0 {
        return 0, errors.New("fibonacci of negative number")
    }
    if n == 0 {
        return 0, nil
    }
    if n == 1 {
        return 1, nil
    }
    
    a, b := int64(0), int64(1)
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b, nil
}

// IsPrime checks if a number is prime
func IsPrime(n int) bool {
    if n < 2 {
        return false
    }
    if n == 2 {
        return true
    }
    if n%2 == 0 {
        return false
    }
    
    sqrt := int(math.Sqrt(float64(n)))
    for i := 3; i <= sqrt; i += 2 {
        if n%i == 0 {
            return false
        }
    }
    return true
}

// GCD calculates the Greatest Common Divisor
func GCD(a, b int) int {
    a, b = abs(a), abs(b)
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// LCM calculates the Least Common Multiple
func LCM(a, b int) int {
    if a == 0 || b == 0 {
        return 0
    }
    return abs(a*b) / GCD(a, b)
}

// abs returns the absolute value of an integer
func abs(n int) int {
    if n < 0 {
        return -n
    }
    return n
}

// Round rounds a number to the specified decimal places
func Round(num float64, places int) float64 {
    shift := math.Pow(10, float64(places))
    return math.Round(num*shift) / shift
}

// Clamp constrains a value between min and max
func Clamp(value, min, max float64) float64 {
    if value < min {
        return min
    }
    if value > max {
        return max
    }
    return value
}
```

**main.go:**
```go
package main

import (
    "fmt"
    "./mathutils" // Local import
)

func main() {
    fmt.Println("=== Math Utilities Demo ===")
    
    // Basic operations
    fmt.Println("\n--- Basic Operations ---")
    fmt.Printf("Add(10, 5): %.2f\n", mathutils.Add(10, 5))
    fmt.Printf("Subtract(10, 5): %.2f\n", mathutils.Subtract(10, 5))
    fmt.Printf("Multiply(10, 5): %.2f\n", mathutils.Multiply(10, 5))
    
    if result, err := mathutils.Divide(10, 5); err != nil {
        fmt.Printf("Divide error: %v\n", err)
    } else {
        fmt.Printf("Divide(10, 5): %.2f\n", result)
    }
    
    if result, err := mathutils.Divide(10, 0); err != nil {
        fmt.Printf("Divide(10, 0) error: %v\n", err)
    }
    
    // Advanced operations
    fmt.Println("\n--- Advanced Operations ---")
    if fact, err := mathutils.Factorial(5); err != nil {
        fmt.Printf("Factorial error: %v\n", err)
    } else {
        fmt.Printf("Factorial(5): %d\n", fact)
    }
    
    if fib, err := mathutils.Fibonacci(10); err != nil {
        fmt.Printf("Fibonacci error: %v\n", err)
    } else {
        fmt.Printf("Fibonacci(10): %d\n", fib)
    }
    
    fmt.Printf("IsPrime(17): %t\n", mathutils.IsPrime(17))
    fmt.Printf("IsPrime(18): %t\n", mathutils.IsPrime(18))
    
    fmt.Printf("GCD(48, 18): %d\n", mathutils.GCD(48, 18))
    fmt.Printf("LCM(12, 8): %d\n", mathutils.LCM(12, 8))
    
    // Constants
    fmt.Println("\n--- Constants ---")
    fmt.Printf("Pi: %.6f\n", mathutils.Pi)
    fmt.Printf("E: %.6f\n", mathutils.E)
    fmt.Printf("Golden Ratio: %.6f\n", mathutils.Golden)
    
    // Utility functions
    fmt.Println("\n--- Utility Functions ---")
    fmt.Printf("Round(3.14159, 2): %.2f\n", mathutils.Round(3.14159, 2))
    fmt.Printf("Clamp(15, 0, 10): %.2f\n", mathutils.Clamp(15, 0, 10))
    fmt.Printf("Min(5, 3): %.2f\n", mathutils.Min(5, 3))
    fmt.Printf("Max(5, 3): %.2f\n", mathutils.Max(5, 3))
}
```

### Example: String Utilities Package

**stringutils/stringutils.go:**
```go
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
    cleaned := removeNonAlphanumeric(s)
    return cleaned == Reverse(cleaned)
}

// removeNonAlphanumeric removes non-alphanumeric characters
func removeNonAlphanumeric(s string) string {
    var result strings.Builder
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            result.WriteRune(r)
        }
    }
    return result.String()
}

// WordCount returns the number of words in a string
func WordCount(s string) int {
    fields := strings.Fields(s)
    return len(fields)
}

// Capitalize capitalizes the first letter of each word
func Capitalize(s string) string {
    return strings.Title(strings.ToLower(s))
}

// TruncateWithEllipsis truncates a string and adds ellipsis if needed
func TruncateWithEllipsis(s string, maxLength int) string {
    if len(s) <= maxLength {
        return s
    }
    if maxLength <= 3 {
        return s[:maxLength]
    }
    return s[:maxLength-3] + "..."
}

// ContainsAny checks if string contains any of the substrings
func ContainsAny(s string, substrings []string) bool {
    for _, substr := range substrings {
        if strings.Contains(s, substr) {
            return true
        }
    }
    return false
}

// RemoveWhitespace removes all whitespace characters
func RemoveWhitespace(s string) string {
    var result strings.Builder
    for _, r := range s {
        if !unicode.IsSpace(r) {
            result.WriteRune(r)
        }
    }
    return result.String()
}

// IsEmpty checks if string is empty or contains only whitespace
func IsEmpty(s string) bool {
    return strings.TrimSpace(s) == ""
}

// Repeat repeats a string n times with a separator
func Repeat(s string, n int, separator string) string {
    if n <= 0 {
        return ""
    }
    if n == 1 {
        return s
    }
    
    parts := make([]string, n)
    for i := 0; i < n; i++ {
        parts[i] = s
    }
    return strings.Join(parts, separator)
}
```

## Go Modules

### What are Go Modules?

Go modules are collections of related Go packages that are versioned together as a single unit. Modules provide:
- Dependency management
- Version control
- Reproducible builds
- Semantic versioning

### Creating a Go Module

```bash
# Initialize a new module
go mod init example.com/myproject

# This creates a go.mod file
```

**go.mod file structure:**
```go
module example.com/myproject

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/lib/pq v1.10.9
)

require (
    github.com/gorilla/context v1.1.1 // indirect
)
```

### Module Commands

```bash
# Add a dependency
go get github.com/gorilla/mux

# Add a specific version
go get github.com/gorilla/mux@v1.8.0

# Update dependencies
go get -u

# Update a specific dependency
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
```

### Example: Creating a Web Server Module

**go.mod:**
```go
module example.com/webserver

go 1.21

require (
    github.com/gorilla/mux v1.8.0
    github.com/rs/cors v1.10.1
)
```

**main.go:**
```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "time"
    
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type Server struct {
    users []User
    nextID int
}

func NewServer() *Server {
    return &Server{
        users: []User{
            {ID: 1, Name: "Alice", Email: "alice@example.com"},
            {ID: 2, Name: "Bob", Email: "bob@example.com"},
        },
        nextID: 3,
    }
}

func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(s.users)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    for _, user := range s.users {
        if user.ID == id {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(user)
            return
        }
    }
    
    http.Error(w, "User not found", http.StatusNotFound)
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user.ID = s.nextID
    s.nextID++
    s.users = append(s.users, user)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (s *Server) setupRoutes() http.Handler {
    r := mux.NewRouter()
    
    r.HandleFunc("/users", s.getUsers).Methods("GET")
    r.HandleFunc("/users/{id}", s.getUser).Methods("GET")
    r.HandleFunc("/users", s.createUser).Methods("POST")
    
    // Add CORS support
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
        AllowedHeaders: []string{"*"},
    })
    
    return c.Handler(r)
}

func main() {
    server := NewServer()
    handler := server.setupRoutes()
    
    srv := &http.Server{
        Addr:         ":8080",
        Handler:      handler,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
    }
    
    log.Println("Server starting on :8080")
    log.Fatal(srv.ListenAndServe())
}
```

## Package Import Patterns

### Standard Import

```go
import "fmt"
import "net/http"
import "encoding/json"

// Or grouped
import (
    "fmt"
    "net/http"
    "encoding/json"
)
```

### Aliased Import

```go
import (
    "fmt"
    mux "github.com/gorilla/mux"
    pq "github.com/lib/pq"
)

func main() {
    router := mux.NewRouter()
    // ...
}
```

### Dot Import (Use Sparingly)

```go
import (
    . "fmt"
    . "math"
)

func main() {
    Println(Sqrt(16)) // No package prefix needed
}
```

### Blank Import

```go
import (
    "database/sql"
    _ "github.com/lib/pq" // Import for side effects only
)

func main() {
    db, err := sql.Open("postgres", "connection_string")
    // ...
}
```

## Package Organization Best Practices

### 1. Flat Package Structure

```
myproject/
├── main.go
├── user.go
├── auth.go
├── database.go
└── config.go
```

### 2. Domain-Driven Structure

```
myproject/
├── main.go
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── user/
│   │   ├── user.go
│   │   ├── repository.go
│   │   └── service.go
│   ├── auth/
│   │   ├── auth.go
│   │   └── middleware.go
│   └── database/
│       └── connection.go
└── pkg/
    └── utils/
        └── helpers.go
```

### 3. Layered Architecture

```
myproject/
├── main.go
├── handlers/
│   ├── user.go
│   └── auth.go
├── services/
│   ├── user.go
│   └── auth.go
├── repositories/
│   ├── user.go
│   └── auth.go
├── models/
│   ├── user.go
│   └── auth.go
└── config/
    └── config.go
```

## Package Documentation

### Package-Level Documentation

```go
// Package mathutils provides mathematical utility functions and constants.
//
// This package includes basic arithmetic operations, advanced mathematical
// functions, and commonly used mathematical constants.
//
// Example usage:
//
//     result := mathutils.Add(10, 5)
//     factorial, err := mathutils.Factorial(5)
//     if err != nil {
//         log.Fatal(err)
//     }
//
package mathutils

import "math"

// Pi represents the mathematical constant π
const Pi = math.Pi
```

### Function Documentation

```go
// Divide performs division of two floating-point numbers.
//
// It returns an error if the divisor is zero.
//
// Example:
//     result, err := Divide(10, 2)
//     if err != nil {
//         log.Fatal(err)
//     }
//     fmt.Println(result) // Output: 5
func Divide(dividend, divisor float64) (float64, error) {
    if divisor == 0 {
        return 0, errors.New("division by zero")
    }
    return dividend / divisor, nil
}
```

### Generating Documentation

```bash
# Generate documentation for current package
go doc

# Generate documentation for specific function
go doc Divide

# Generate documentation for external package
go doc fmt.Println

# Start documentation server
godoc -http=:6060
```

## Testing Packages

### Package Testing Structure

```
mathutils/
├── basic.go
├── basic_test.go
├── advanced.go
├── advanced_test.go
└── example_test.go
```

**basic_test.go:**
```go
package mathutils

import (
    "testing"
)

func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     float64
        expected float64
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
                t.Errorf("Add(%f, %f) = %f; want %f", tt.a, tt.b, result, tt.expected)
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

func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(10, 5)
    }
}
```

**example_test.go:**
```go
package mathutils_test

import (
    "fmt"
    "./mathutils"
)

func ExampleAdd() {
    result := mathutils.Add(2, 3)
    fmt.Println(result)
    // Output: 5
}

func ExampleDivide() {
    result, err := mathutils.Divide(10, 2)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result)
    // Output: 5
}

func ExampleFactorial() {
    result, err := mathutils.Factorial(5)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println(result)
    // Output: 120
}
```

## Practical Examples

### Database Package

**database/database.go:**
```go
package database

import (
    "database/sql"
    "fmt"
    "time"
    
    _ "github.com/lib/pq"
)

type Config struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

type DB struct {
    conn *sql.DB
}

func New(config Config) (*DB, error) {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
    
    conn, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    conn.SetMaxOpenConns(25)
    conn.SetMaxIdleConns(5)
    conn.SetConnMaxLifetime(5 * time.Minute)
    
    // Test connection
    if err := conn.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    return &DB{conn: conn}, nil
}

func (db *DB) Close() error {
    return db.conn.Close()
}

func (db *DB) Ping() error {
    return db.conn.Ping()
}

func (db *DB) GetConnection() *sql.DB {
    return db.conn
}

// Health check
func (db *DB) HealthCheck() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return db.conn.PingContext(ctx)
}
```

### Configuration Package

**config/config.go:**
```go
package config

import (
    "encoding/json"
    "fmt"
    "os"
    "strconv"
    "time"
)

type Config struct {
    Server   ServerConfig   `json:"server"`
    Database DatabaseConfig `json:"database"`
    Redis    RedisConfig    `json:"redis"`
    Logging  LoggingConfig  `json:"logging"`
}

type ServerConfig struct {
    Host         string        `json:"host"`
    Port         int           `json:"port"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    User     string `json:"user"`
    Password string `json:"password"`
    DBName   string `json:"db_name"`
    SSLMode  string `json:"ssl_mode"`
}

type RedisConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Password string `json:"password"`
    DB       int    `json:"db"`
}

type LoggingConfig struct {
    Level  string `json:"level"`
    Format string `json:"format"`
    Output string `json:"output"`
}

// Load loads configuration from file and environment variables
func Load(filename string) (*Config, error) {
    config := &Config{}
    
    // Load from file if it exists
    if filename != "" {
        if err := loadFromFile(config, filename); err != nil {
            return nil, fmt.Errorf("failed to load config from file: %w", err)
        }
    }
    
    // Override with environment variables
    loadFromEnv(config)
    
    // Set defaults
    setDefaults(config)
    
    return config, nil
}

func loadFromFile(config *Config, filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    decoder := json.NewDecoder(file)
    return decoder.Decode(config)
}

func loadFromEnv(config *Config) {
    // Server configuration
    if host := os.Getenv("SERVER_HOST"); host != "" {
        config.Server.Host = host
    }
    if port := os.Getenv("SERVER_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Server.Port = p
        }
    }
    
    // Database configuration
    if host := os.Getenv("DB_HOST"); host != "" {
        config.Database.Host = host
    }
    if port := os.Getenv("DB_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Database.Port = p
        }
    }
    if user := os.Getenv("DB_USER"); user != "" {
        config.Database.User = user
    }
    if password := os.Getenv("DB_PASSWORD"); password != "" {
        config.Database.Password = password
    }
    if dbname := os.Getenv("DB_NAME"); dbname != "" {
        config.Database.DBName = dbname
    }
    
    // Redis configuration
    if host := os.Getenv("REDIS_HOST"); host != "" {
        config.Redis.Host = host
    }
    if port := os.Getenv("REDIS_PORT"); port != "" {
        if p, err := strconv.Atoi(port); err == nil {
            config.Redis.Port = p
        }
    }
    if password := os.Getenv("REDIS_PASSWORD"); password != "" {
        config.Redis.Password = password
    }
    
    // Logging configuration
    if level := os.Getenv("LOG_LEVEL"); level != "" {
        config.Logging.Level = level
    }
}

func setDefaults(config *Config) {
    // Server defaults
    if config.Server.Host == "" {
        config.Server.Host = "localhost"
    }
    if config.Server.Port == 0 {
        config.Server.Port = 8080
    }
    if config.Server.ReadTimeout == 0 {
        config.Server.ReadTimeout = 15 * time.Second
    }
    if config.Server.WriteTimeout == 0 {
        config.Server.WriteTimeout = 15 * time.Second
    }
    
    // Database defaults
    if config.Database.Host == "" {
        config.Database.Host = "localhost"
    }
    if config.Database.Port == 0 {
        config.Database.Port = 5432
    }
    if config.Database.SSLMode == "" {
        config.Database.SSLMode = "disable"
    }
    
    // Redis defaults
    if config.Redis.Host == "" {
        config.Redis.Host = "localhost"
    }
    if config.Redis.Port == 0 {
        config.Redis.Port = 6379
    }
    
    // Logging defaults
    if config.Logging.Level == "" {
        config.Logging.Level = "info"
    }
    if config.Logging.Format == "" {
        config.Logging.Format = "json"
    }
    if config.Logging.Output == "" {
        config.Logging.Output = "stdout"
    }
}

// Validate validates the configuration
func (c *Config) Validate() error {
    if c.Server.Port < 1 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    
    if c.Database.User == "" {
        return fmt.Errorf("database user is required")
    }
    
    if c.Database.DBName == "" {
        return fmt.Errorf("database name is required")
    }
    
    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    
    if !validLogLevels[c.Logging.Level] {
        return fmt.Errorf("invalid log level: %s", c.Logging.Level)
    }
    
    return nil
}
```

## Best Practices

### 1. Package Naming

```go
// Good
package http
package json
package user

// Avoid
package httputils
package jsonparser
package userservice
```

### 2. Package Size

- Keep packages focused and cohesive
- Avoid overly large packages
- Split packages when they serve different purposes

### 3. Dependency Management

```go
// Use specific versions in production
go get github.com/gorilla/mux@v1.8.0

// Keep dependencies up to date
go get -u
go mod tidy
```

### 4. Internal Packages

```
myproject/
├── cmd/
├── internal/  // Cannot be imported by external projects
│   ├── auth/
│   └── database/
└── pkg/       // Can be imported by external projects
    └── utils/
```

### 5. Avoid Circular Dependencies

```go
// Bad: Package A imports B, Package B imports A

// Good: Extract common functionality to a third package
package common

package a
import "./common"

package b
import "./common"
```

## Exercises

### Exercise 1: Logger Package

Create a logging package with different log levels and outputs.

```go
// Implement:
// - Logger interface
// - Console logger
// - File logger
// - JSON formatter
// - Different log levels (DEBUG, INFO, WARN, ERROR)
```

### Exercise 2: HTTP Client Package

Create an HTTP client package with retry logic and timeout handling.

```go
// Implement:
// - HTTP client with configurable timeouts
// - Retry mechanism with exponential backoff
// - Request/response logging
// - Error handling
```

### Exercise 3: Cache Package

Create a caching package with multiple backends.

```go
// Implement:
// - Cache interface
// - In-memory cache
// - Redis cache (simulated)
// - TTL support
// - Cache statistics
```

## Key Takeaways

1. **Packages provide modularity**: Organize related functionality together.

2. **Visibility through naming**: Uppercase = exported, lowercase = unexported.

3. **Go modules manage dependencies**: Use semantic versioning and reproducible builds.

4. **Keep packages focused**: Each package should have a single, well-defined purpose.

5. **Document your packages**: Provide clear documentation for exported functions and types.

6. **Test your packages**: Write comprehensive tests for all exported functionality.

7. **Use internal packages**: Keep implementation details private with internal packages.

8. **Avoid circular dependencies**: Design packages to have clear dependency hierarchies.

9. **Follow naming conventions**: Use clear, concise package names.

10. **Version your modules**: Use semantic versioning for public modules.

## Next Steps

Now that you understand packages and modules, let's learn about [File I/O](17-file-io.md) to see how to work with files and directories in Go!

---

**Previous**: [← Error Handling](15-error-handling.md) | **Next**: [File I/O →](17-file-io.md)