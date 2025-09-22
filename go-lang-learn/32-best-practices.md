# 32. Best Practices

This final chapter consolidates the most important best practices for writing clean, efficient, and maintainable Go code. These practices are derived from the Go community, official guidelines, and real-world experience.

## Code Organization and Structure

### Package Design

```go
// Good package design example

// Package user provides user management functionality
package user

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// User represents a user in the system
type User struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Service defines the user service interface
type Service interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int64) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

// Repository defines the user repository interface
type Repository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id int64) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int64) error
}

// Validator defines the user validation interface
type Validator interface {
    Validate(user *User) error
}

// Common errors
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidUser      = errors.New("invalid user")
    ErrUserAlreadyExists = errors.New("user already exists")
)

// service implements the Service interface
type service struct {
    repo      Repository
    validator Validator
}

// NewService creates a new user service
func NewService(repo Repository, validator Validator) Service {
    return &service{
        repo:      repo,
        validator: validator,
    }
}

// Create creates a new user
func (s *service) Create(ctx context.Context, user *User) error {
    if user == nil {
        return ErrInvalidUser
    }
    
    if err := s.validator.Validate(user); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()
    
    if err := s.repo.Save(ctx, user); err != nil {
        return fmt.Errorf("failed to save user: %w", err)
    }
    
    return nil
}

// GetByID retrieves a user by ID
func (s *service) GetByID(ctx context.Context, id int64) (*User, error) {
    if id <= 0 {
        return nil, ErrInvalidUser
    }
    
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    
    return user, nil
}

// Update updates an existing user
func (s *service) Update(ctx context.Context, user *User) error {
    if user == nil || user.ID <= 0 {
        return ErrInvalidUser
    }
    
    if err := s.validator.Validate(user); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    user.UpdatedAt = time.Now()
    
    if err := s.repo.Update(ctx, user); err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}

// Delete deletes a user by ID
func (s *service) Delete(ctx context.Context, id int64) error {
    if id <= 0 {
        return ErrInvalidUser
    }
    
    if err := s.repo.Delete(ctx, id); err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }
    
    return nil
}
```

### Directory Structure

```
project/
├── cmd/                    # Main applications
│   ├── server/
│   │   └── main.go
│   └── cli/
│       └── main.go
├── internal/               # Private application code
│   ├── user/              # User domain
│   │   ├── service.go
│   │   ├── repository.go
│   │   └── handler.go
│   ├── auth/              # Authentication domain
│   └── config/            # Configuration
├── pkg/                   # Public library code
│   ├── logger/
│   ├── database/
│   └── middleware/
├── api/                   # API definitions
│   └── openapi.yaml
├── web/                   # Web assets
├── scripts/               # Build and deployment scripts
├── docs/                  # Documentation
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

## Error Handling

### Comprehensive Error Handling

```go
// error_handling.go
package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

// Custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type NotFoundError struct {
    Resource string
    ID       string
}

func (e NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID '%s' not found", e.Resource, e.ID)
}

// Error wrapping and unwrapping
func processUser(userID string) error {
    // Validate input
    if userID == "" {
        return ValidationError{
            Field:   "userID",
            Message: "cannot be empty",
        }
    }
    
    // Convert to integer
    id, err := strconv.Atoi(userID)
    if err != nil {
        return fmt.Errorf("invalid user ID format: %w", err)
    }
    
    // Simulate database operation
    if id == 404 {
        return NotFoundError{
            Resource: "User",
            ID:       userID,
        }
    }
    
    if id < 0 {
        return fmt.Errorf("user ID must be positive, got %d", id)
    }
    
    return nil
}

// Error handling in HTTP handlers
func userHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    
    err := processUser(userID)
    if err != nil {
        handleError(w, err)
        return
    }
    
    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "User %s processed successfully", userID)
}

func handleError(w http.ResponseWriter, err error) {
    // Log the error
    log.Printf("Error: %v", err)
    
    // Handle different error types
    var validationErr ValidationError
    var notFoundErr NotFoundError
    
    switch {
    case errors.As(err, &validationErr):
        http.Error(w, validationErr.Error(), http.StatusBadRequest)
    case errors.As(err, &notFoundErr):
        http.Error(w, notFoundErr.Error(), http.StatusNotFound)
    default:
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}

// Sentinel errors
var (
    ErrInvalidInput = errors.New("invalid input")
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

// Error checking with errors.Is
func checkError(err error) {
    if errors.Is(err, ErrNotFound) {
        log.Println("Resource not found")
    } else if errors.Is(err, ErrUnauthorized) {
        log.Println("Access denied")
    }
}
```

## Concurrency Best Practices

### Safe Concurrency Patterns

```go
// concurrency_patterns.go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Worker pool with graceful shutdown
type WorkerPool struct {
    workers    int
    jobs       chan Job
    results    chan Result
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID int
    Data  string
    Error error
}

func NewWorkerPool(workers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &WorkerPool{
        workers: workers,
        jobs:    make(chan Job, workers*2), // Buffered channel
        results: make(chan Result, workers*2),
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for {
        select {
        case job, ok := <-wp.jobs:
            if !ok {
                return // Channel closed
            }
            
            // Process job
            result := wp.processJob(job)
            
            select {
            case wp.results <- result:
            case <-wp.ctx.Done():
                return
            }
            
        case <-wp.ctx.Done():
            return
        }
    }
}

func (wp *WorkerPool) processJob(job Job) Result {
    // Simulate work
    time.Sleep(100 * time.Millisecond)
    
    return Result{
        JobID: job.ID,
        Data:  fmt.Sprintf("Processed: %s", job.Data),
        Error: nil,
    }
}

func (wp *WorkerPool) Submit(job Job) error {
    select {
    case wp.jobs <- job:
        return nil
    case <-wp.ctx.Done():
        return wp.ctx.Err()
    }
}

func (wp *WorkerPool) Results() <-chan Result {
    return wp.results
}

func (wp *WorkerPool) Shutdown() {
    close(wp.jobs)  // Signal no more jobs
    wp.wg.Wait()    // Wait for workers to finish
    wp.cancel()     // Cancel context
    close(wp.results) // Close results channel
}

// Rate limiter
type RateLimiter struct {
    tokens chan struct{}
    ticker *time.Ticker
    done   chan struct{}
}

func NewRateLimiter(rate int, burst int) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, burst),
        ticker: time.NewTicker(time.Second / time.Duration(rate)),
        done:   make(chan struct{}),
    }
    
    // Fill initial tokens
    for i := 0; i < burst; i++ {
        rl.tokens <- struct{}{}
    }
    
    go rl.refill()
    return rl
}

func (rl *RateLimiter) refill() {
    for {
        select {
        case <-rl.ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default: // Bucket is full
            }
        case <-rl.done:
            rl.ticker.Stop()
            return
        }
    }
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (rl *RateLimiter) Stop() {
    close(rl.done)
}

// Circuit breaker
type CircuitBreaker struct {
    mu           sync.RWMutex
    state        State
    failures     int
    lastFailTime time.Time
    threshold    int
    timeout      time.Duration
}

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:     StateClosed,
        threshold: threshold,
        timeout:   timeout,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.canExecute() {
        return fmt.Errorf("circuit breaker is open")
    }
    
    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) canExecute() bool {
    cb.mu.RLock()
    defer cb.mu.RUnlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        return time.Since(cb.lastFailTime) > cb.timeout
    case StateHalfOpen:
        return true
    default:
        return false
    }
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailTime = time.Now()
        
        if cb.failures >= cb.threshold {
            cb.state = StateOpen
        }
    } else {
        cb.failures = 0
        cb.state = StateClosed
    }
}
```

## Testing Best Practices

### Comprehensive Testing Strategy

```go
// testing_best_practices_test.go
package main

import (
    "context"
    "errors"
    "testing"
    "time"
)

// Table-driven tests
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        want    bool
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "user@example.com",
            want:    true,
            wantErr: false,
        },
        {
            name:    "invalid email - no @",
            email:   "userexample.com",
            want:    false,
            wantErr: true,
        },
        {
            name:    "invalid email - no domain",
            email:   "user@",
            want:    false,
            wantErr: true,
        },
        {
            name:    "empty email",
            email:   "",
            want:    false,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := validateEmail(tt.email)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if got != tt.want {
                t.Errorf("validateEmail() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Mock interfaces for testing
type MockUserRepository struct {
    users map[int64]*User
    err   error
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[int64]*User),
    }
}

func (m *MockUserRepository) Save(ctx context.Context, user *User) error {
    if m.err != nil {
        return m.err
    }
    
    if user.ID == 0 {
        user.ID = int64(len(m.users) + 1)
    }
    
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int64) (*User, error) {
    if m.err != nil {
        return nil, m.err
    }
    
    user, exists := m.users[id]
    if !exists {
        return nil, ErrUserNotFound
    }
    
    return user, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *User) error {
    if m.err != nil {
        return m.err
    }
    
    if _, exists := m.users[user.ID]; !exists {
        return ErrUserNotFound
    }
    
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
    if m.err != nil {
        return m.err
    }
    
    if _, exists := m.users[id]; !exists {
        return ErrUserNotFound
    }
    
    delete(m.users, id)
    return nil
}

func (m *MockUserRepository) SetError(err error) {
    m.err = err
}

// Test with mocks
func TestUserService_Create(t *testing.T) {
    mockRepo := NewMockUserRepository()
    mockValidator := &MockValidator{}
    service := NewService(mockRepo, mockValidator)
    
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    err := service.Create(context.Background(), user)
    if err != nil {
        t.Errorf("Create() error = %v, want nil", err)
    }
    
    if user.ID == 0 {
        t.Error("Create() should set user ID")
    }
    
    if user.CreatedAt.IsZero() {
        t.Error("Create() should set CreatedAt")
    }
}

// Benchmark tests
func BenchmarkUserService_Create(b *testing.B) {
    mockRepo := NewMockUserRepository()
    mockValidator := &MockValidator{}
    service := NewService(mockRepo, mockValidator)
    
    user := &User{
        Username: "testuser",
        Email:    "test@example.com",
    }
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        user.ID = 0 // Reset ID for each iteration
        err := service.Create(context.Background(), user)
        if err != nil {
            b.Fatalf("Create() error = %v", err)
        }
    }
}

// Test helpers
func setupTestUser(t *testing.T) *User {
    t.Helper()
    
    return &User{
        ID:        1,
        Username:  "testuser",
        Email:     "test@example.com",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func assertUserEqual(t *testing.T, got, want *User) {
    t.Helper()
    
    if got.ID != want.ID {
        t.Errorf("ID = %v, want %v", got.ID, want.ID)
    }
    
    if got.Username != want.Username {
        t.Errorf("Username = %v, want %v", got.Username, want.Username)
    }
    
    if got.Email != want.Email {
        t.Errorf("Email = %v, want %v", got.Email, want.Email)
    }
}

// Mock validator
type MockValidator struct {
    err error
}

func (m *MockValidator) Validate(user *User) error {
    return m.err
}

func (m *MockValidator) SetError(err error) {
    m.err = err
}

// Helper functions for testing
func validateEmail(email string) (bool, error) {
    if email == "" {
        return false, errors.New("email cannot be empty")
    }
    
    // Simple validation
    if !contains(email, "@") || !contains(email, ".") {
        return false, errors.New("invalid email format")
    }
    
    return true, nil
}

func contains(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
```

## Performance Best Practices

### Memory and CPU Optimization

```go
// performance_optimization.go
package main

import (
    "bytes"
    "fmt"
    "strings"
    "sync"
)

// String building optimization
func buildStringEfficient(parts []string) string {
    // Use strings.Builder for efficient string concatenation
    var builder strings.Builder
    
    // Pre-allocate capacity if known
    totalLen := 0
    for _, part := range parts {
        totalLen += len(part)
    }
    builder.Grow(totalLen)
    
    for _, part := range parts {
        builder.WriteString(part)
    }
    
    return builder.String()
}

// Slice optimization
func processSliceEfficient(data []int) []int {
    // Pre-allocate slice with known capacity
    result := make([]int, 0, len(data))
    
    for _, item := range data {
        if item > 0 {
            result = append(result, item*2)
        }
    }
    
    return result
}

// Object pooling for frequent allocations
var bufferPool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func processWithPool(data []byte) []byte {
    // Get buffer from pool
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    // Use buffer
    buf.Write(data)
    buf.WriteString(" processed")
    
    // Return copy since buffer will be reused
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    
    return result
}

// Map optimization
func createMapEfficient(size int) map[string]int {
    // Pre-allocate map with estimated size
    m := make(map[string]int, size)
    
    for i := 0; i < size; i++ {
        key := fmt.Sprintf("key_%d", i)
        m[key] = i
    }
    
    return m
}

// Interface optimization - avoid boxing when possible
type Processor interface {
    Process(data []byte) []byte
}

type ConcreteProcessor struct{}

func (cp ConcreteProcessor) Process(data []byte) []byte {
    // Process data
    return data
}

// Use concrete type when interface is not needed
func processDirectly(cp ConcreteProcessor, data []byte) []byte {
    return cp.Process(data)
}

// Avoid allocations in hot paths
func hotPathFunction(data []int) int {
    sum := 0
    
    // Avoid creating slices or maps in hot paths
    for _, v := range data {
        sum += v
    }
    
    return sum
}
```

## Security Best Practices

### Secure Coding Guidelines

```go
// security_practices.go
package main

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/hex"
    "fmt"
    "html/template"
    "net/http"
    "regexp"
    "strings"
    "time"
)

// Input validation
func validateInput(input string) error {
    // Length validation
    if len(input) > 1000 {
        return fmt.Errorf("input too long")
    }
    
    // Character validation
    validPattern := regexp.MustCompile(`^[a-zA-Z0-9\s\-_\.]+$`)
    if !validPattern.MatchString(input) {
        return fmt.Errorf("input contains invalid characters")
    }
    
    return nil
}

// SQL injection prevention
func safeQuery(db Database, userID string) (*User, error) {
    // Always use parameterized queries
    query := "SELECT id, username, email FROM users WHERE id = ?"
    
    var user User
    err := db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// XSS prevention
func renderTemplate(w http.ResponseWriter, data interface{}) error {
    // Use html/template for automatic escaping
    tmpl := template.Must(template.New("page").Parse(`
        <html>
        <body>
            <h1>Welcome {{.Username}}</h1>
            <p>Email: {{.Email}}</p>
        </body>
        </html>
    `))
    
    return tmpl.Execute(w, data)
}

// Secure random generation
func generateSecureToken(length int) (string, error) {
    bytes := make([]byte, length)
    
    _, err := rand.Read(bytes)
    if err != nil {
        return "", err
    }
    
    return hex.EncodeToString(bytes), nil
}

// Constant-time comparison
func secureCompare(a, b string) bool {
    return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// Rate limiting middleware
func rateLimitMiddleware(limiter *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if !limiter.Allow() {
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Security headers middleware
func securityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set security headers
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        
        next.ServeHTTP(w, r)
    })
}

// Timeout middleware
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.TimeoutHandler(next, timeout, "Request timeout")
    }
}

// Database interface for example
type Database interface {
    QueryRow(query string, args ...interface{}) Row
}

type Row interface {
    Scan(dest ...interface{}) error
}
```

## Documentation and Code Quality

### Documentation Best Practices

```go
// Package documentation should describe the package purpose
// and provide usage examples.
//
// Package calculator provides basic mathematical operations
// with support for floating-point arithmetic.
//
// Example usage:
//
//     calc := calculator.New()
//     result := calc.Add(2.5, 3.7)
//     fmt.Printf("Result: %.2f\n", result)
//
package calculator

import (
    "errors"
    "math"
)

// Calculator provides mathematical operations.
// All operations are safe for concurrent use.
type Calculator struct {
    precision int
}

// New creates a new Calculator with default precision.
func New() *Calculator {
    return &Calculator{
        precision: 2,
    }
}

// NewWithPrecision creates a new Calculator with specified precision.
// Precision must be between 0 and 10.
func NewWithPrecision(precision int) (*Calculator, error) {
    if precision < 0 || precision > 10 {
        return nil, errors.New("precision must be between 0 and 10")
    }
    
    return &Calculator{
        precision: precision,
    }, nil
}

// Add returns the sum of a and b.
func (c *Calculator) Add(a, b float64) float64 {
    result := a + b
    return c.round(result)
}

// Subtract returns the difference of a and b.
func (c *Calculator) Subtract(a, b float64) float64 {
    result := a - b
    return c.round(result)
}

// Multiply returns the product of a and b.
func (c *Calculator) Multiply(a, b float64) float64 {
    result := a * b
    return c.round(result)
}

// Divide returns the quotient of a and b.
// Returns an error if b is zero.
func (c *Calculator) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    
    result := a / b
    return c.round(result), nil
}

// SetPrecision sets the precision for calculations.
// Precision must be between 0 and 10.
func (c *Calculator) SetPrecision(precision int) error {
    if precision < 0 || precision > 10 {
        return errors.New("precision must be between 0 and 10")
    }
    
    c.precision = precision
    return nil
}

// round rounds the value to the calculator's precision.
func (c *Calculator) round(value float64) float64 {
    multiplier := math.Pow(10, float64(c.precision))
    return math.Round(value*multiplier) / multiplier
}
```

## Configuration Management

### Environment-based Configuration

```go
// config.go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

// Config holds application configuration
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Logger   LoggerConfig
}

type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type DatabaseConfig struct {
    Host         string
    Port         int
    Username     string
    Password     string
    Database     string
    MaxOpenConns int
    MaxIdleConns int
    MaxLifetime  time.Duration
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    Database int
}

type LoggerConfig struct {
    Level  string
    Format string
    Output string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
    config := &Config{}
    
    // Server configuration
    config.Server.Host = getEnvString("SERVER_HOST", "localhost")
    
    port, err := getEnvInt("SERVER_PORT", 8080)
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
    }
    config.Server.Port = port
    
    readTimeout, err := getEnvDuration("SERVER_READ_TIMEOUT", "30s")
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_READ_TIMEOUT: %w", err)
    }
    config.Server.ReadTimeout = readTimeout
    
    writeTimeout, err := getEnvDuration("SERVER_WRITE_TIMEOUT", "30s")
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_WRITE_TIMEOUT: %w", err)
    }
    config.Server.WriteTimeout = writeTimeout
    
    idleTimeout, err := getEnvDuration("SERVER_IDLE_TIMEOUT", "120s")
    if err != nil {
        return nil, fmt.Errorf("invalid SERVER_IDLE_TIMEOUT: %w", err)
    }
    config.Server.IdleTimeout = idleTimeout
    
    // Database configuration
    config.Database.Host = getEnvString("DB_HOST", "localhost")
    
    dbPort, err := getEnvInt("DB_PORT", 5432)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_PORT: %w", err)
    }
    config.Database.Port = dbPort
    
    config.Database.Username = getEnvString("DB_USERNAME", "")
    config.Database.Password = getEnvString("DB_PASSWORD", "")
    config.Database.Database = getEnvString("DB_DATABASE", "")
    
    maxOpenConns, err := getEnvInt("DB_MAX_OPEN_CONNS", 25)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_MAX_OPEN_CONNS: %w", err)
    }
    config.Database.MaxOpenConns = maxOpenConns
    
    maxIdleConns, err := getEnvInt("DB_MAX_IDLE_CONNS", 5)
    if err != nil {
        return nil, fmt.Errorf("invalid DB_MAX_IDLE_CONNS: %w", err)
    }
    config.Database.MaxIdleConns = maxIdleConns
    
    maxLifetime, err := getEnvDuration("DB_MAX_LIFETIME", "5m")
    if err != nil {
        return nil, fmt.Errorf("invalid DB_MAX_LIFETIME: %w", err)
    }
    config.Database.MaxLifetime = maxLifetime
    
    // Redis configuration
    config.Redis.Host = getEnvString("REDIS_HOST", "localhost")
    
    redisPort, err := getEnvInt("REDIS_PORT", 6379)
    if err != nil {
        return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
    }
    config.Redis.Port = redisPort
    
    config.Redis.Password = getEnvString("REDIS_PASSWORD", "")
    
    redisDB, err := getEnvInt("REDIS_DATABASE", 0)
    if err != nil {
        return nil, fmt.Errorf("invalid REDIS_DATABASE: %w", err)
    }
    config.Redis.Database = redisDB
    
    // Logger configuration
    config.Logger.Level = getEnvString("LOG_LEVEL", "info")
    config.Logger.Format = getEnvString("LOG_FORMAT", "json")
    config.Logger.Output = getEnvString("LOG_OUTPUT", "stdout")
    
    return config, nil
}

// Helper functions for environment variable parsing
func getEnvString(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvInt(key string, defaultValue int) (int, error) {
    if value := os.Getenv(key); value != "" {
        return strconv.Atoi(value)
    }
    return defaultValue, nil
}

func getEnvDuration(key, defaultValue string) (time.Duration, error) {
    value := getEnvString(key, defaultValue)
    return time.ParseDuration(value)
}

// Validate validates the configuration
func (c *Config) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    
    if c.Database.Username == "" {
        return fmt.Errorf("database username is required")
    }
    
    if c.Database.Database == "" {
        return fmt.Errorf("database name is required")
    }
    
    validLogLevels := map[string]bool{
        "debug": true,
        "info":  true,
        "warn":  true,
        "error": true,
    }
    
    if !validLogLevels[strings.ToLower(c.Logger.Level)] {
        return fmt.Errorf("invalid log level: %s", c.Logger.Level)
    }
    
    return nil
}
```

## Key Principles Summary

### 1. Code Organization
- Use clear package structure
- Separate concerns properly
- Follow Go naming conventions
- Keep interfaces small and focused

### 2. Error Handling
- Always handle errors explicitly
- Use custom error types when appropriate
- Wrap errors with context
- Don't ignore errors

### 3. Concurrency
- Use channels for communication
- Avoid shared mutable state
- Always handle goroutine lifecycle
- Use context for cancellation

### 4. Performance
- Profile before optimizing
- Minimize allocations
- Use appropriate data structures
- Consider memory pooling for hot paths

### 5. Security
- Validate all inputs
- Use parameterized queries
- Implement proper authentication/authorization
- Set security headers

### 6. Testing
- Write tests first (TDD)
- Use table-driven tests
- Mock external dependencies
- Aim for high test coverage

### 7. Documentation
- Document public APIs
- Provide usage examples
- Keep documentation up to date
- Use meaningful names

## Final Recommendations

1. **Follow Go idioms**: Write Go code that looks like Go code
2. **Keep it simple**: Prefer simple solutions over complex ones
3. **Be explicit**: Make your code's intent clear
4. **Handle errors**: Never ignore errors
5. **Test thoroughly**: Write comprehensive tests
6. **Profile and optimize**: Measure before optimizing
7. **Document well**: Write clear documentation
8. **Review regularly**: Conduct code reviews
9. **Stay updated**: Keep up with Go best practices
10. **Practice consistently**: Apply these practices in all projects

## Conclusion

This comprehensive guide covers the essential best practices for Go development. By following these guidelines, you'll write more maintainable, efficient, and secure Go applications. Remember that best practices evolve, so stay engaged with the Go community and continue learning.

Congratulations on completing the Go learning path! You now have a solid foundation in Go programming and the knowledge to build robust, scalable applications.