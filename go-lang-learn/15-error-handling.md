# 15. Error Handling

Error handling is a fundamental aspect of Go programming. Unlike many other languages that use exceptions, Go uses explicit error values to handle errors. This approach makes error handling more predictable and forces developers to think about error cases.

## The Error Interface

In Go, errors are represented by the built-in `error` interface:

```go
type error interface {
    Error() string
}
```

### Basic Error Handling

```go
package main

import (
    "errors"
    "fmt"
    "strconv"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func parseAndDivide(aStr, bStr string) (float64, error) {
    a, err := strconv.ParseFloat(aStr, 64)
    if err != nil {
        return 0, fmt.Errorf("failed to parse first number: %v", err)
    }
    
    b, err := strconv.ParseFloat(bStr, 64)
    if err != nil {
        return 0, fmt.Errorf("failed to parse second number: %v", err)
    }
    
    result, err := divide(a, b)
    if err != nil {
        return 0, fmt.Errorf("division failed: %v", err)
    }
    
    return result, nil
}

func main() {
    // Successful case
    result, err := parseAndDivide("10", "2")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %.2f\n", result)
    }
    
    // Error cases
    _, err = parseAndDivide("10", "0")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    _, err = parseAndDivide("abc", "2")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Creating Custom Errors

### Using errors.New()

```go
package main

import (
    "errors"
    "fmt"
)

var (
    ErrInvalidAge    = errors.New("age must be between 0 and 150")
    ErrEmptyName     = errors.New("name cannot be empty")
    ErrInvalidEmail  = errors.New("invalid email format")
)

type Person struct {
    Name  string
    Age   int
    Email string
}

func NewPerson(name string, age int, email string) (*Person, error) {
    if name == "" {
        return nil, ErrEmptyName
    }
    
    if age < 0 || age > 150 {
        return nil, ErrInvalidAge
    }
    
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    
    return &Person{
        Name:  name,
        Age:   age,
        Email: email,
    }, nil
}

func isValidEmail(email string) bool {
    // Simple email validation
    return len(email) > 0 && 
           len(email) > 3 && 
           email[len(email)-4:] == ".com" || 
           email[len(email)-4:] == ".org"
}

func main() {
    // Valid person
    person, err := NewPerson("Alice", 30, "alice@example.com")
    if err != nil {
        fmt.Printf("Error creating person: %v\n", err)
    } else {
        fmt.Printf("Created person: %+v\n", person)
    }
    
    // Invalid cases
    _, err = NewPerson("", 30, "alice@example.com")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    _, err = NewPerson("Bob", 200, "bob@example.com")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    _, err = NewPerson("Charlie", 25, "invalid-email")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Custom Error Types

```go
package main

import (
    "fmt"
    "time"
)

// Custom error type
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
    Time    time.Time
}

func (ve ValidationError) Error() string {
    return fmt.Sprintf("validation error in field '%s': %s (value: %v, time: %s)",
        ve.Field, ve.Message, ve.Value, ve.Time.Format(time.RFC3339))
}

// Another custom error type
type NetworkError struct {
    Operation string
    URL       string
    Err       error
}

func (ne NetworkError) Error() string {
    return fmt.Sprintf("network error during %s to %s: %v", ne.Operation, ne.URL, ne.Err)
}

func (ne NetworkError) Unwrap() error {
    return ne.Err
}

// User validation function
func validateUser(name string, age int, email string) error {
    if name == "" {
        return ValidationError{
            Field:   "name",
            Value:   name,
            Message: "cannot be empty",
            Time:    time.Now(),
        }
    }
    
    if age < 0 || age > 150 {
        return ValidationError{
            Field:   "age",
            Value:   age,
            Message: "must be between 0 and 150",
            Time:    time.Now(),
        }
    }
    
    if len(email) < 5 {
        return ValidationError{
            Field:   "email",
            Value:   email,
            Message: "must be at least 5 characters",
            Time:    time.Now(),
        }
    }
    
    return nil
}

// Simulated network operation
func fetchUserData(url string) error {
    // Simulate network error
    if url == "" {
        return NetworkError{
            Operation: "GET",
            URL:       url,
            Err:       fmt.Errorf("empty URL provided"),
        }
    }
    
    if url == "http://invalid.com" {
        return NetworkError{
            Operation: "GET",
            URL:       url,
            Err:       fmt.Errorf("connection timeout"),
        }
    }
    
    return nil
}

func main() {
    // Validation errors
    err := validateUser("", 25, "test@example.com")
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
        
        // Type assertion to access custom fields
        if ve, ok := err.(ValidationError); ok {
            fmt.Printf("Field: %s, Value: %v\n", ve.Field, ve.Value)
        }
    }
    
    err = validateUser("Alice", 200, "alice@example.com")
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
    }
    
    // Network errors
    err = fetchUserData("")
    if err != nil {
        fmt.Printf("Network error: %v\n", err)
        
        if ne, ok := err.(NetworkError); ok {
            fmt.Printf("Operation: %s, URL: %s\n", ne.Operation, ne.URL)
        }
    }
    
    err = fetchUserData("http://invalid.com")
    if err != nil {
        fmt.Printf("Network error: %v\n", err)
    }
}
```

## Error Wrapping and Unwrapping

### Using fmt.Errorf with %w

```go
package main

import (
    "errors"
    "fmt"
    "strconv"
)

var (
    ErrInvalidInput = errors.New("invalid input")
    ErrOutOfRange   = errors.New("value out of range")
)

func parseNumber(s string) (int, error) {
    num, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("failed to parse number '%s': %w", s, err)
    }
    return num, nil
}

func validateRange(num int, min, max int) error {
    if num < min || num > max {
        return fmt.Errorf("number %d is out of range [%d, %d]: %w", num, min, max, ErrOutOfRange)
    }
    return nil
}

func processInput(input string, min, max int) (int, error) {
    if input == "" {
        return 0, fmt.Errorf("empty input: %w", ErrInvalidInput)
    }
    
    num, err := parseNumber(input)
    if err != nil {
        return 0, fmt.Errorf("input processing failed: %w", err)
    }
    
    if err := validateRange(num, min, max); err != nil {
        return 0, fmt.Errorf("validation failed: %w", err)
    }
    
    return num, nil
}

func main() {
    inputs := []string{"42", "abc", "200", ""}
    
    for _, input := range inputs {
        result, err := processInput(input, 1, 100)
        if err != nil {
            fmt.Printf("Error processing '%s': %v\n", input, err)
            
            // Check for specific errors using errors.Is
            if errors.Is(err, ErrInvalidInput) {
                fmt.Println("  -> This is an invalid input error")
            }
            if errors.Is(err, ErrOutOfRange) {
                fmt.Println("  -> This is an out of range error")
            }
            if errors.Is(err, strconv.ErrSyntax) {
                fmt.Println("  -> This is a syntax error")
            }
            
            fmt.Println()
        } else {
            fmt.Printf("Successfully processed '%s': %d\n\n", input, result)
        }
    }
}
```

### Using errors.Is and errors.As

```go
package main

import (
    "errors"
    "fmt"
    "net"
    "os"
    "syscall"
)

type CustomError struct {
    Code    int
    Message string
    Cause   error
}

func (ce CustomError) Error() string {
    return fmt.Sprintf("custom error [%d]: %s", ce.Code, ce.Message)
}

func (ce CustomError) Unwrap() error {
    return ce.Cause
}

func simulateFileOperation(filename string) error {
    if filename == "" {
        return CustomError{
            Code:    400,
            Message: "filename cannot be empty",
            Cause:   os.ErrInvalid,
        }
    }
    
    if filename == "nonexistent.txt" {
        return CustomError{
            Code:    404,
            Message: "file not found",
            Cause:   os.ErrNotExist,
        }
    }
    
    if filename == "permission.txt" {
        return CustomError{
            Code:    403,
            Message: "permission denied",
            Cause:   os.ErrPermission,
        }
    }
    
    return nil
}

func simulateNetworkOperation(host string) error {
    if host == "timeout.com" {
        return CustomError{
            Code:    408,
            Message: "request timeout",
            Cause:   &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ETIMEDOUT},
        }
    }
    
    if host == "refused.com" {
        return CustomError{
            Code:    503,
            Message: "connection refused",
            Cause:   &net.OpError{Op: "dial", Net: "tcp", Err: syscall.ECONNREFUSED},
        }
    }
    
    return nil
}

func handleError(err error, operation string) {
    if err == nil {
        fmt.Printf("%s completed successfully\n", operation)
        return
    }
    
    fmt.Printf("%s failed: %v\n", operation, err)
    
    // Check for specific error types using errors.Is
    if errors.Is(err, os.ErrNotExist) {
        fmt.Println("  -> File does not exist")
    } else if errors.Is(err, os.ErrPermission) {
        fmt.Println("  -> Permission denied")
    } else if errors.Is(err, os.ErrInvalid) {
        fmt.Println("  -> Invalid argument")
    }
    
    // Check for custom error type using errors.As
    var customErr CustomError
    if errors.As(err, &customErr) {
        fmt.Printf("  -> Custom error code: %d\n", customErr.Code)
        
        switch customErr.Code {
        case 400:
            fmt.Println("  -> Bad request")
        case 403:
            fmt.Println("  -> Forbidden")
        case 404:
            fmt.Println("  -> Not found")
        case 408:
            fmt.Println("  -> Timeout")
        case 503:
            fmt.Println("  -> Service unavailable")
        }
    }
    
    // Check for network error
    var netErr *net.OpError
    if errors.As(err, &netErr) {
        fmt.Printf("  -> Network operation: %s on %s\n", netErr.Op, netErr.Net)
        
        if errors.Is(netErr.Err, syscall.ETIMEDOUT) {
            fmt.Println("  -> Connection timed out")
        } else if errors.Is(netErr.Err, syscall.ECONNREFUSED) {
            fmt.Println("  -> Connection refused")
        }
    }
    
    fmt.Println()
}

func main() {
    // File operations
    fmt.Println("=== File Operations ===")
    filenames := []string{"valid.txt", "", "nonexistent.txt", "permission.txt"}
    
    for _, filename := range filenames {
        err := simulateFileOperation(filename)
        handleError(err, fmt.Sprintf("File operation on '%s'", filename))
    }
    
    // Network operations
    fmt.Println("=== Network Operations ===")
    hosts := []string{"valid.com", "timeout.com", "refused.com"}
    
    for _, host := range hosts {
        err := simulateNetworkOperation(host)
        handleError(err, fmt.Sprintf("Network operation to '%s'", host))
    }
}
```

## Error Handling Patterns

### Sentinel Errors

```go
package main

import (
    "errors"
    "fmt"
)

// Sentinel errors
var (
    ErrUserNotFound     = errors.New("user not found")
    ErrInvalidPassword  = errors.New("invalid password")
    ErrAccountLocked    = errors.New("account locked")
    ErrSessionExpired   = errors.New("session expired")
)

type User struct {
    ID       int
    Username string
    Password string
    Locked   bool
}

type AuthService struct {
    users map[string]User
}

func NewAuthService() *AuthService {
    return &AuthService{
        users: map[string]User{
            "alice": {ID: 1, Username: "alice", Password: "secret123", Locked: false},
            "bob":   {ID: 2, Username: "bob", Password: "password456", Locked: true},
            "charlie": {ID: 3, Username: "charlie", Password: "mypass789", Locked: false},
        },
    }
}

func (as *AuthService) Authenticate(username, password string) (*User, error) {
    user, exists := as.users[username]
    if !exists {
        return nil, ErrUserNotFound
    }
    
    if user.Locked {
        return nil, ErrAccountLocked
    }
    
    if user.Password != password {
        return nil, ErrInvalidPassword
    }
    
    return &user, nil
}

func handleAuthError(err error) string {
    switch {
    case errors.Is(err, ErrUserNotFound):
        return "User does not exist. Please check your username."
    case errors.Is(err, ErrInvalidPassword):
        return "Incorrect password. Please try again."
    case errors.Is(err, ErrAccountLocked):
        return "Your account has been locked. Please contact support."
    case errors.Is(err, ErrSessionExpired):
        return "Your session has expired. Please log in again."
    default:
        return "An unexpected error occurred. Please try again later."
    }
}

func main() {
    authService := NewAuthService()
    
    testCases := []struct {
        username string
        password string
    }{
        {"alice", "secret123"},     // Valid
        {"alice", "wrongpass"},     // Invalid password
        {"bob", "password456"},     // Account locked
        {"nonexistent", "pass"},   // User not found
    }
    
    for _, tc := range testCases {
        user, err := authService.Authenticate(tc.username, tc.password)
        if err != nil {
            fmt.Printf("Login failed for %s: %s\n", tc.username, handleAuthError(err))
        } else {
            fmt.Printf("Login successful for %s (ID: %d)\n", user.Username, user.ID)
        }
    }
}
```

### Retry Pattern

```go
package main

import (
    "errors"
    "fmt"
    "math/rand"
    "time"
)

type RetryableError struct {
    Err error
}

func (re RetryableError) Error() string {
    return fmt.Sprintf("retryable error: %v", re.Err)
}

func (re RetryableError) Unwrap() error {
    return re.Err
}

func IsRetryable(err error) bool {
    var retryableErr RetryableError
    return errors.As(err, &retryableErr)
}

// Simulated unreliable operation
func unreliableOperation(id string) error {
    rand.Seed(time.Now().UnixNano())
    
    switch rand.Intn(4) {
    case 0:
        return nil // Success
    case 1:
        return RetryableError{Err: errors.New("temporary network error")}
    case 2:
        return RetryableError{Err: errors.New("service temporarily unavailable")}
    default:
        return errors.New("permanent error: invalid request") // Non-retryable
    }
}

func retry(operation func() error, maxAttempts int, delay time.Duration) error {
    var lastErr error
    
    for attempt := 1; attempt <= maxAttempts; attempt++ {
        err := operation()
        if err == nil {
            return nil // Success
        }
        
        lastErr = err
        
        if !IsRetryable(err) {
            return fmt.Errorf("non-retryable error on attempt %d: %w", attempt, err)
        }
        
        if attempt < maxAttempts {
            fmt.Printf("Attempt %d failed: %v. Retrying in %v...\n", attempt, err, delay)
            time.Sleep(delay)
            delay *= 2 // Exponential backoff
        }
    }
    
    return fmt.Errorf("operation failed after %d attempts: %w", maxAttempts, lastErr)
}

func main() {
    fmt.Println("=== Retry Pattern Example ===")
    
    operations := []string{"op1", "op2", "op3"}
    
    for _, opID := range operations {
        fmt.Printf("\nExecuting operation %s:\n", opID)
        
        err := retry(func() error {
            return unreliableOperation(opID)
        }, 3, 100*time.Millisecond)
        
        if err != nil {
            fmt.Printf("Operation %s failed: %v\n", opID, err)
        } else {
            fmt.Printf("Operation %s succeeded!\n", opID)
        }
    }
}
```

### Circuit Breaker Pattern

```go
package main

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

type CircuitState int

const (
    Closed CircuitState = iota
    Open
    HalfOpen
)

func (cs CircuitState) String() string {
    switch cs {
    case Closed:
        return "Closed"
    case Open:
        return "Open"
    case HalfOpen:
        return "Half-Open"
    default:
        return "Unknown"
    }
}

type CircuitBreaker struct {
    mu              sync.Mutex
    state           CircuitState
    failureCount    int
    successCount    int
    lastFailureTime time.Time
    
    maxFailures     int
    resetTimeout    time.Duration
    halfOpenMaxCalls int
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration, halfOpenMaxCalls int) *CircuitBreaker {
    return &CircuitBreaker{
        state:            Closed,
        maxFailures:      maxFailures,
        resetTimeout:     resetTimeout,
        halfOpenMaxCalls: halfOpenMaxCalls,
    }
}

func (cb *CircuitBreaker) Call(operation func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if cb.state == Open {
        if time.Since(cb.lastFailureTime) > cb.resetTimeout {
            cb.state = HalfOpen
            cb.successCount = 0
            fmt.Println("Circuit breaker: Transitioning to Half-Open")
        } else {
            return errors.New("circuit breaker is open")
        }
    }
    
    if cb.state == HalfOpen && cb.successCount >= cb.halfOpenMaxCalls {
        cb.state = Closed
        cb.failureCount = 0
        cb.successCount = 0
        fmt.Println("Circuit breaker: Transitioning to Closed")
    }
    
    err := operation()
    
    if err != nil {
        cb.failureCount++
        cb.lastFailureTime = time.Now()
        
        if cb.state == HalfOpen {
            cb.state = Open
            fmt.Println("Circuit breaker: Transitioning to Open (failed in half-open)")
        } else if cb.failureCount >= cb.maxFailures {
            cb.state = Open
            fmt.Printf("Circuit breaker: Transitioning to Open (max failures: %d)\n", cb.maxFailures)
        }
        
        return err
    }
    
    // Success
    if cb.state == HalfOpen {
        cb.successCount++
    } else {
        cb.failureCount = 0
    }
    
    return nil
}

func (cb *CircuitBreaker) GetState() CircuitState {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}

// Simulated service that fails sometimes
var serviceCallCount int

func unreliableService() error {
    serviceCallCount++
    
    // Fail for the first 5 calls, then succeed
    if serviceCallCount <= 5 {
        return errors.New("service is down")
    }
    
    return nil
}

func main() {
    fmt.Println("=== Circuit Breaker Pattern Example ===")
    
    cb := NewCircuitBreaker(
        3,                // Max failures before opening
        2*time.Second,    // Reset timeout
        2,                // Half-open max calls
    )
    
    // Simulate multiple service calls
    for i := 1; i <= 15; i++ {
        fmt.Printf("\nCall %d (Circuit: %s):\n", i, cb.GetState())
        
        err := cb.Call(unreliableService)
        if err != nil {
            fmt.Printf("  Error: %v\n", err)
        } else {
            fmt.Printf("  Success!\n")
        }
        
        // Wait a bit between calls
        time.Sleep(300 * time.Millisecond)
        
        // After call 8, wait for reset timeout
        if i == 8 {
            fmt.Println("\nWaiting for circuit breaker reset timeout...")
            time.Sleep(2500 * time.Millisecond)
        }
    }
}
```

### Error Aggregation

```go
package main

import (
    "errors"
    "fmt"
    "strings"
)

// MultiError aggregates multiple errors
type MultiError struct {
    errors []error
}

func (me *MultiError) Add(err error) {
    if err != nil {
        me.errors = append(me.errors, err)
    }
}

func (me *MultiError) Error() string {
    if len(me.errors) == 0 {
        return "no errors"
    }
    
    if len(me.errors) == 1 {
        return me.errors[0].Error()
    }
    
    var messages []string
    for i, err := range me.errors {
        messages = append(messages, fmt.Sprintf("%d: %v", i+1, err))
    }
    
    return fmt.Sprintf("multiple errors occurred:\n%s", strings.Join(messages, "\n"))
}

func (me *MultiError) HasErrors() bool {
    return len(me.errors) > 0
}

func (me *MultiError) Errors() []error {
    return me.errors
}

func (me *MultiError) Count() int {
    return len(me.errors)
}

// Validation functions
func validateName(name string) error {
    if name == "" {
        return errors.New("name cannot be empty")
    }
    if len(name) < 2 {
        return errors.New("name must be at least 2 characters")
    }
    return nil
}

func validateAge(age int) error {
    if age < 0 {
        return errors.New("age cannot be negative")
    }
    if age > 150 {
        return errors.New("age cannot be greater than 150")
    }
    return nil
}

func validateEmail(email string) error {
    if email == "" {
        return errors.New("email cannot be empty")
    }
    if !strings.Contains(email, "@") {
        return errors.New("email must contain @ symbol")
    }
    if !strings.Contains(email, ".") {
        return errors.New("email must contain a domain")
    }
    return nil
}

func validatePhone(phone string) error {
    if phone == "" {
        return errors.New("phone cannot be empty")
    }
    if len(phone) < 10 {
        return errors.New("phone must be at least 10 digits")
    }
    return nil
}

// User validation with error aggregation
func validateUser(name string, age int, email, phone string) error {
    var multiErr MultiError
    
    multiErr.Add(validateName(name))
    multiErr.Add(validateAge(age))
    multiErr.Add(validateEmail(email))
    multiErr.Add(validatePhone(phone))
    
    if multiErr.HasErrors() {
        return &multiErr
    }
    
    return nil
}

// Batch processing with error aggregation
func processUsers(users []map[string]interface{}) error {
    var multiErr MultiError
    
    for i, user := range users {
        name, _ := user["name"].(string)
        age, _ := user["age"].(int)
        email, _ := user["email"].(string)
        phone, _ := user["phone"].(string)
        
        if err := validateUser(name, age, email, phone); err != nil {
            multiErr.Add(fmt.Errorf("user %d validation failed: %w", i+1, err))
        }
    }
    
    if multiErr.HasErrors() {
        return &multiErr
    }
    
    return nil
}

func main() {
    fmt.Println("=== Error Aggregation Example ===")
    
    // Single user validation
    fmt.Println("\n--- Single User Validation ---")
    err := validateUser("", -5, "invalid-email", "123")
    if err != nil {
        fmt.Printf("Validation failed:\n%v\n", err)
        
        // Check if it's a MultiError
        if multiErr, ok := err.(*MultiError); ok {
            fmt.Printf("\nTotal errors: %d\n", multiErr.Count())
            for i, e := range multiErr.Errors() {
                fmt.Printf("Error %d: %v\n", i+1, e)
            }
        }
    }
    
    // Batch processing
    fmt.Println("\n--- Batch Processing ---")
    users := []map[string]interface{}{
        {"name": "Alice", "age": 30, "email": "alice@example.com", "phone": "1234567890"},
        {"name": "", "age": 25, "email": "bob@example.com", "phone": "0987654321"},
        {"name": "Charlie", "age": -5, "email": "invalid", "phone": "123"},
        {"name": "Diana", "age": 28, "email": "diana@example.com", "phone": "5555555555"},
    }
    
    err = processUsers(users)
    if err != nil {
        fmt.Printf("Batch processing failed:\n%v\n", err)
    } else {
        fmt.Println("All users processed successfully!")
    }
}
```

## Practical Examples

### HTTP Client with Error Handling

```go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
    "time"
)

type HTTPError struct {
    StatusCode int
    Status     string
    URL        string
    Method     string
}

func (he HTTPError) Error() string {
    return fmt.Sprintf("HTTP %s %s failed: %d %s", he.Method, he.URL, he.StatusCode, he.Status)
}

type APIClient struct {
    baseURL    string
    httpClient *http.Client
}

func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        httpClient: &http.Client{
            Timeout: timeout,
        },
    }
}

func (c *APIClient) get(endpoint string) ([]byte, error) {
    url := c.baseURL + endpoint
    
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to make GET request to %s: %w", url, err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, HTTPError{
            StatusCode: resp.StatusCode,
            Status:     resp.Status,
            URL:        url,
            Method:     "GET",
        }
    }
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body from %s: %w", url, err)
    }
    
    return body, nil
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func (c *APIClient) GetUser(id int) (*User, error) {
    endpoint := fmt.Sprintf("/users/%d", id)
    
    data, err := c.get(endpoint)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    
    var user User
    if err := json.Unmarshal(data, &user); err != nil {
        return nil, fmt.Errorf("failed to parse user data: %w", err)
    }
    
    return &user, nil
}

func (c *APIClient) GetUsers() ([]User, error) {
    data, err := c.get("/users")
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }
    
    var users []User
    if err := json.Unmarshal(data, &users); err != nil {
        return nil, fmt.Errorf("failed to parse users data: %w", err)
    }
    
    return users, nil
}

func handleAPIError(err error, operation string) {
    fmt.Printf("%s failed: %v\n", operation, err)
    
    var httpErr HTTPError
    if errors.As(err, &httpErr) {
        switch httpErr.StatusCode {
        case 404:
            fmt.Println("  -> Resource not found")
        case 401:
            fmt.Println("  -> Unauthorized - check your credentials")
        case 403:
            fmt.Println("  -> Forbidden - insufficient permissions")
        case 429:
            fmt.Println("  -> Rate limited - try again later")
        case 500:
            fmt.Println("  -> Server error - try again later")
        default:
            fmt.Printf("  -> HTTP error: %d\n", httpErr.StatusCode)
        }
    }
    
    // Check for network errors
    if errors.Is(err, io.EOF) {
        fmt.Println("  -> Connection closed unexpectedly")
    }
}

func main() {
    fmt.Println("=== HTTP Client Error Handling Example ===")
    
    // Note: This uses a real API endpoint for demonstration
    client := NewAPIClient("https://jsonplaceholder.typicode.com", 10*time.Second)
    
    // Successful request
    fmt.Println("\n--- Getting User 1 ---")
    user, err := client.GetUser(1)
    if err != nil {
        handleAPIError(err, "Get User 1")
    } else {
        fmt.Printf("Success: %+v\n", user)
    }
    
    // Non-existent user (404 error)
    fmt.Println("\n--- Getting Non-existent User ---")
    _, err = client.GetUser(999999)
    if err != nil {
        handleAPIError(err, "Get User 999999")
    }
    
    // Invalid endpoint (using wrong client for demonstration)
    fmt.Println("\n--- Invalid Endpoint ---")
    invalidClient := NewAPIClient("https://invalid-domain-that-does-not-exist.com", 5*time.Second)
    _, err = invalidClient.GetUser(1)
    if err != nil {
        handleAPIError(err, "Get User from Invalid Domain")
    }
}
```

### File Processing with Error Handling

```go
package main

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

type FileProcessingError struct {
    Filename string
    Line     int
    Column   int
    Err      error
}

func (fpe FileProcessingError) Error() string {
    if fpe.Line > 0 {
        return fmt.Sprintf("error in file %s at line %d, column %d: %v", 
            fpe.Filename, fpe.Line, fpe.Column, fpe.Err)
    }
    return fmt.Sprintf("error in file %s: %v", fpe.Filename, fpe.Err)
}

func (fpe FileProcessingError) Unwrap() error {
    return fpe.Err
}

type CSVProcessor struct {
    filename string
    file     *os.File
    scanner  *bufio.Scanner
    lineNum  int
}

func NewCSVProcessor(filename string) (*CSVProcessor, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, FileProcessingError{
            Filename: filename,
            Err:      fmt.Errorf("failed to open file: %w", err),
        }
    }
    
    return &CSVProcessor{
        filename: filename,
        file:     file,
        scanner:  bufio.NewScanner(file),
    }, nil
}

func (cp *CSVProcessor) Close() error {
    if cp.file != nil {
        return cp.file.Close()
    }
    return nil
}

func (cp *CSVProcessor) ProcessLine() ([]string, error) {
    if !cp.scanner.Scan() {
        if err := cp.scanner.Err(); err != nil {
            return nil, FileProcessingError{
                Filename: cp.filename,
                Line:     cp.lineNum,
                Err:      fmt.Errorf("scanner error: %w", err),
            }
        }
        return nil, io.EOF
    }
    
    cp.lineNum++
    line := cp.scanner.Text()
    
    if strings.TrimSpace(line) == "" {
        return nil, FileProcessingError{
            Filename: cp.filename,
            Line:     cp.lineNum,
            Err:      errors.New("empty line"),
        }
    }
    
    fields := strings.Split(line, ",")
    
    // Trim whitespace from fields
    for i, field := range fields {
        fields[i] = strings.TrimSpace(field)
    }
    
    return fields, nil
}

type Person struct {
    Name string
    Age  int
    City string
}

func parsePersonFromFields(fields []string, filename string, lineNum int) (*Person, error) {
    if len(fields) != 3 {
        return nil, FileProcessingError{
            Filename: filename,
            Line:     lineNum,
            Err:      fmt.Errorf("expected 3 fields, got %d", len(fields)),
        }
    }
    
    name := fields[0]
    if name == "" {
        return nil, FileProcessingError{
            Filename: filename,
            Line:     lineNum,
            Column:   1,
            Err:      errors.New("name cannot be empty"),
        }
    }
    
    age, err := strconv.Atoi(fields[1])
    if err != nil {
        return nil, FileProcessingError{
            Filename: filename,
            Line:     lineNum,
            Column:   2,
            Err:      fmt.Errorf("invalid age: %w", err),
        }
    }
    
    if age < 0 || age > 150 {
        return nil, FileProcessingError{
            Filename: filename,
            Line:     lineNum,
            Column:   2,
            Err:      fmt.Errorf("age %d is out of valid range (0-150)", age),
        }
    }
    
    city := fields[2]
    if city == "" {
        return nil, FileProcessingError{
            Filename: filename,
            Line:     lineNum,
            Column:   3,
            Err:      errors.New("city cannot be empty"),
        }
    }
    
    return &Person{
        Name: name,
        Age:  age,
        City: city,
    }, nil
}

func processCSVFile(filename string) ([]Person, error) {
    processor, err := NewCSVProcessor(filename)
    if err != nil {
        return nil, err
    }
    defer processor.Close()
    
    var people []Person
    var errors []error
    
    for {
        fields, err := processor.ProcessLine()
        if err != nil {
            if errors.Is(err, io.EOF) {
                break
            }
            
            // Collect error but continue processing
            errors = append(errors, err)
            continue
        }
        
        person, err := parsePersonFromFields(fields, filename, processor.lineNum)
        if err != nil {
            errors = append(errors, err)
            continue
        }
        
        people = append(people, *person)
    }
    
    if len(errors) > 0 {
        // Return partial results with aggregated errors
        var multiErr MultiError
        for _, err := range errors {
            multiErr.Add(err)
        }
        return people, &multiErr
    }
    
    return people, nil
}

func createSampleCSV(filename string) error {
    content := `Alice, 30, New York
Bob, 25, Los Angeles
, 35, Chicago
Diana, invalid_age, Miami
Eve, 28, Seattle
Frank, 200, Boston
Grace, 22, `
    
    return os.WriteFile(filename, []byte(content), 0644)
}

func main() {
    fmt.Println("=== File Processing Error Handling Example ===")
    
    // Create sample CSV file
    filename := "sample.csv"
    if err := createSampleCSV(filename); err != nil {
        fmt.Printf("Failed to create sample file: %v\n", err)
        return
    }
    defer os.Remove(filename) // Clean up
    
    fmt.Printf("\nProcessing file: %s\n", filename)
    
    people, err := processCSVFile(filename)
    
    fmt.Printf("\nSuccessfully processed %d records:\n", len(people))
    for i, person := range people {
        fmt.Printf("  %d: %s, %d, %s\n", i+1, person.Name, person.Age, person.City)
    }
    
    if err != nil {
        fmt.Printf("\nErrors encountered during processing:\n%v\n", err)
        
        // Handle specific error types
        if multiErr, ok := err.(*MultiError); ok {
            fmt.Printf("\nDetailed error breakdown (%d errors):\n", multiErr.Count())
            for i, e := range multiErr.Errors() {
                fmt.Printf("  %d: %v\n", i+1, e)
                
                var fpe FileProcessingError
                if errors.As(e, &fpe) {
                    if fpe.Line > 0 && fpe.Column > 0 {
                        fmt.Printf("     -> Location: Line %d, Column %d\n", fpe.Line, fpe.Column)
                    }
                }
            }
        }
    }
    
    // Test with non-existent file
    fmt.Println("\n--- Testing with non-existent file ---")
    _, err = processCSVFile("nonexistent.csv")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        
        var fpe FileProcessingError
        if errors.As(err, &fpe) {
            fmt.Printf("File processing error detected for: %s\n", fpe.Filename)
        }
        
        if errors.Is(err, os.ErrNotExist) {
            fmt.Println("File does not exist")
        }
    }
}
```

## Best Practices

### 1. Always Handle Errors

```go
// Good
result, err := someFunction()
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Bad - ignoring errors
result, _ := someFunction()
```

### 2. Provide Context in Error Messages

```go
// Good
func processFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read file %s: %w", filename, err)
    }
    // ...
}

// Bad - no context
func processFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    // ...
}
```

### 3. Use Sentinel Errors for Expected Conditions

```go
var ErrNotFound = errors.New("item not found")

func findItem(id string) (*Item, error) {
    // ...
    if !found {
        return nil, ErrNotFound
    }
    // ...
}
```

### 4. Wrap Errors to Preserve Context

```go
func processData(data []byte) error {
    if err := validate(data); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // ...
}
```

### 5. Use errors.Is and errors.As for Error Checking

```go
// Check for specific error
if errors.Is(err, ErrNotFound) {
    // handle not found
}

// Extract error type
var netErr *net.OpError
if errors.As(err, &netErr) {
    // handle network error
}
```

## Exercises

### Exercise 1: Configuration Validator

Create a configuration validator that checks multiple fields and returns all validation errors.

```go
type Config struct {
    Host     string
    Port     int
    Username string
    Password string
    Timeout  time.Duration
}

// Implement:
// - ValidateConfig(config Config) error
// - Custom error types for different validation failures
// - Error aggregation for multiple validation errors
```

### Exercise 2: Retry Mechanism

Implement a generic retry mechanism with exponential backoff.

```go
// Implement:
// - Retry(operation func() error, maxAttempts int, initialDelay time.Duration) error
// - Support for retryable vs non-retryable errors
// - Exponential backoff with jitter
```

### Exercise 3: Error Recovery

Create a system that can recover from certain types of errors.

```go
// Implement:
// - A service that can fail and recover
// - Different error types (temporary, permanent, recoverable)
// - Automatic recovery mechanisms
```

## Key Takeaways

1. **Explicit error handling**: Go uses explicit error values instead of exceptions.

2. **Error interface**: Errors implement the simple `error` interface with an `Error() string` method.

3. **Error wrapping**: Use `fmt.Errorf` with `%w` to wrap errors and preserve context.

4. **Error checking**: Use `errors.Is` and `errors.As` to check for specific errors and types.

5. **Custom error types**: Create custom error types for domain-specific error information.

6. **Sentinel errors**: Use predefined error variables for expected error conditions.

7. **Error aggregation**: Collect multiple errors when processing batches or validating multiple fields.

8. **Context in errors**: Always provide meaningful context in error messages.

9. **Fail fast vs. collect errors**: Choose between stopping at first error or collecting all errors based on use case.

10. **Error handling patterns**: Use patterns like retry, circuit breaker, and graceful degradation for robust systems.

## Next Steps

Now that you understand error handling, let's learn about [Packages and Modules](16-packages-modules.md) to see how to organize and distribute Go code!

---

**Previous**: [← Interfaces](14-interfaces.md) | **Next**: [Packages and Modules →](16-packages-modules.md)