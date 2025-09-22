# 13. Error Handling

Go takes a unique approach to error handling using explicit error values rather than exceptions. This makes error handling visible and forces developers to deal with errors explicitly, leading to more robust code.

## Error Basics

### The Error Interface

```go
// Built-in error interface
type error interface {
    Error() string
}
```

### Basic Error Handling

```go
package main

import (
    "fmt"
    "strconv"
)

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
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
    // Successful operation
    result, err := divide(10, 2)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %.2f\n", result)
    }
    
    // Error case
    result, err = divide(10, 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %.2f\n", result)
    }
    
    // Chained operations with error handling
    result, err = parseAndDivide("10", "2")
    if err != nil {
        fmt.Printf("Parse and divide error: %v\n", err)
    } else {
        fmt.Printf("Parse and divide result: %.2f\n", result)
    }
    
    // Error in parsing
    result, err = parseAndDivide("abc", "2")
    if err != nil {
        fmt.Printf("Parse and divide error: %v\n", err)
    } else {
        fmt.Printf("Parse and divide result: %.2f\n", result)
    }
}
```

### Creating Custom Errors

```go
package main

import (
    "fmt"
    "time"
)

// Simple custom error type
type ValidationError struct {
    Field   string
    Value   interface{}
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error for field '%s' with value '%v': %s", 
        e.Field, e.Value, e.Message)
}

// More complex custom error type
type DatabaseError struct {
    Operation string
    Table     string
    Err       error
    Timestamp time.Time
    Code      int
}

func (e DatabaseError) Error() string {
    return fmt.Sprintf("database error [%d] during %s on table '%s' at %s: %v",
        e.Code, e.Operation, e.Table, e.Timestamp.Format(time.RFC3339), e.Err)
}

// Error that wraps another error
func (e DatabaseError) Unwrap() error {
    return e.Err
}

// User validation function
func validateUser(name string, age int, email string) error {
    if name == "" {
        return ValidationError{
            Field:   "name",
            Value:   name,
            Message: "name cannot be empty",
        }
    }
    
    if age < 0 || age > 150 {
        return ValidationError{
            Field:   "age",
            Value:   age,
            Message: "age must be between 0 and 150",
        }
    }
    
    if email == "" {
        return ValidationError{
            Field:   "email",
            Value:   email,
            Message: "email cannot be empty",
        }
    }
    
    return nil
}

// Database operation simulation
func saveUser(name string, age int, email string) error {
    // Validate first
    if err := validateUser(name, age, email); err != nil {
        return fmt.Errorf("user validation failed: %w", err)
    }
    
    // Simulate database error
    if name == "error" {
        dbErr := fmt.Errorf("connection timeout")
        return DatabaseError{
            Operation: "INSERT",
            Table:     "users",
            Err:       dbErr,
            Timestamp: time.Now(),
            Code:      1001,
        }
    }
    
    fmt.Printf("User saved successfully: %s, %d, %s\n", name, age, email)
    return nil
}

func main() {
    // Test cases
    testCases := []struct {
        name  string
        age   int
        email string
    }{
        {"John Doe", 30, "john@example.com"},
        {"", 25, "empty@example.com"},
        {"Jane", -5, "jane@example.com"},
        {"Bob", 25, ""},
        {"error", 30, "error@example.com"},
    }
    
    for i, tc := range testCases {
        fmt.Printf("\nTest case %d:\n", i+1)
        err := saveUser(tc.name, tc.age, tc.email)
        
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            
            // Type assertion to check error type
            var validationErr ValidationError
            var dbErr DatabaseError
            
            if fmt.Errorf("%w", err); validationErr != (ValidationError{}) {
                fmt.Printf("This is a validation error\n")
            } else if fmt.Errorf("%w", err); dbErr != (DatabaseError{}) {
                fmt.Printf("This is a database error\n")
            }
        }
    }
}
```

### Error Wrapping and Unwrapping

```go
package main

import (
    "errors"
    "fmt"
    "io"
    "os"
)

// Custom error types
type FileError struct {
    Filename string
    Op       string
    Err      error
}

func (e FileError) Error() string {
    return fmt.Sprintf("file error: %s %s: %v", e.Op, e.Filename, e.Err)
}

func (e FileError) Unwrap() error {
    return e.Err
}

type ProcessingError struct {
    Stage string
    Err   error
}

func (e ProcessingError) Error() string {
    return fmt.Sprintf("processing error at stage '%s': %v", e.Stage, e.Err)
}

func (e ProcessingError) Unwrap() error {
    return e.Err
}

// Functions that wrap errors
func readFile(filename string) ([]byte, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, FileError{
            Filename: filename,
            Op:       "read",
            Err:      err,
        }
    }
    return data, nil
}

func processData(data []byte) (string, error) {
    if len(data) == 0 {
        return "", ProcessingError{
            Stage: "validation",
            Err:   errors.New("empty data"),
        }
    }
    
    // Simulate processing error
    if string(data[:1]) == "#" {
        return "", ProcessingError{
            Stage: "parsing",
            Err:   errors.New("invalid format: starts with #"),
        }
    }
    
    return string(data), nil
}

func processFile(filename string) (string, error) {
    data, err := readFile(filename)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }
    
    result, err := processData(data)
    if err != nil {
        return "", fmt.Errorf("failed to process data: %w", err)
    }
    
    return result, nil
}

func main() {
    // Create test files
    testFiles := map[string]string{
        "valid.txt":   "Hello, World!",
        "empty.txt":   "",
        "invalid.txt": "#invalid format",
    }
    
    for filename, content := range testFiles {
        err := os.WriteFile(filename, []byte(content), 0644)
        if err != nil {
            fmt.Printf("Failed to create test file %s: %v\n", filename, err)
            continue
        }
        defer os.Remove(filename) // Clean up
    }
    
    // Test processing files
    filenames := []string{"valid.txt", "empty.txt", "invalid.txt", "nonexistent.txt"}
    
    for _, filename := range filenames {
        fmt.Printf("\nProcessing %s:\n", filename)
        result, err := processFile(filename)
        
        if err != nil {
            fmt.Printf("Error: %v\n", err)
            
            // Check for specific error types using errors.Is
            if errors.Is(err, os.ErrNotExist) {
                fmt.Println("  -> File does not exist")
            }
            
            // Check for custom error types using errors.As
            var fileErr FileError
            if errors.As(err, &fileErr) {
                fmt.Printf("  -> File error: operation=%s, filename=%s\n", 
                    fileErr.Op, fileErr.Filename)
            }
            
            var procErr ProcessingError
            if errors.As(err, &procErr) {
                fmt.Printf("  -> Processing error: stage=%s\n", procErr.Stage)
            }
            
            // Unwrap errors manually
            fmt.Println("  -> Error chain:")
            currentErr := err
            level := 1
            for currentErr != nil {
                fmt.Printf("    %d. %v\n", level, currentErr)
                currentErr = errors.Unwrap(currentErr)
                level++
            }
        } else {
            fmt.Printf("Success: %s\n", result)
        }
    }
    
    // Demonstrate error comparison
    fmt.Println("\n=== Error Comparison Examples ===")
    
    // Create some errors
    err1 := errors.New("base error")
    err2 := fmt.Errorf("wrapped: %w", err1)
    err3 := fmt.Errorf("double wrapped: %w", err2)
    
    fmt.Printf("err1: %v\n", err1)
    fmt.Printf("err2: %v\n", err2)
    fmt.Printf("err3: %v\n", err3)
    
    // Check if err3 contains err1
    fmt.Printf("errors.Is(err3, err1): %t\n", errors.Is(err3, err1))
    fmt.Printf("errors.Is(err2, err1): %t\n", errors.Is(err2, err1))
    
    // Unwrap chain
    fmt.Println("\nUnwrapping err3:")
    current := err3
    for current != nil {
        fmt.Printf("  %v\n", current)
        current = errors.Unwrap(current)
    }
}
```

### Error Handling Patterns

```go
package main

import (
    "errors"
    "fmt"
    "time"
)

// Sentinel errors
var (
    ErrNotFound     = errors.New("item not found")
    ErrInvalidInput = errors.New("invalid input")
    ErrTimeout      = errors.New("operation timeout")
    ErrUnauthorized = errors.New("unauthorized access")
)

// Result type for operations that might fail
type Result[T any] struct {
    Value T
    Error error
}

func (r Result[T]) IsSuccess() bool {
    return r.Error == nil
}

func (r Result[T]) IsError() bool {
    return r.Error != nil
}

// Retry pattern
func retryOperation(operation func() error, maxRetries int, delay time.Duration) error {
    var lastErr error
    
    for attempt := 0; attempt <= maxRetries; attempt++ {
        if attempt > 0 {
            fmt.Printf("Retry attempt %d/%d\n", attempt, maxRetries)
            time.Sleep(delay)
        }
        
        err := operation()
        if err == nil {
            return nil // Success
        }
        
        lastErr = err
        
        // Don't retry for certain types of errors
        if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrInvalidInput) {
            return err
        }
    }
    
    return fmt.Errorf("operation failed after %d retries: %w", maxRetries, lastErr)
}

// Circuit breaker pattern (simplified)
type CircuitBreaker struct {
    maxFailures int
    failures    int
    lastFailure time.Time
    timeout     time.Duration
    state       string // "closed", "open", "half-open"
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       "closed",
    }
}

func (cb *CircuitBreaker) Call(operation func() error) error {
    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
            fmt.Println("Circuit breaker: half-open")
        } else {
            return fmt.Errorf("circuit breaker is open")
        }
    }
    
    err := operation()
    
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
            fmt.Println("Circuit breaker: opened")
        }
        
        return err
    }
    
    // Success - reset circuit breaker
    if cb.state == "half-open" {
        cb.state = "closed"
        fmt.Println("Circuit breaker: closed")
    }
    cb.failures = 0
    
    return nil
}

// Error aggregation
type MultiError struct {
    Errors []error
}

func (me MultiError) Error() string {
    if len(me.Errors) == 0 {
        return "no errors"
    }
    
    if len(me.Errors) == 1 {
        return me.Errors[0].Error()
    }
    
    result := fmt.Sprintf("%d errors occurred:", len(me.Errors))
    for i, err := range me.Errors {
        result += fmt.Sprintf("\n  %d. %v", i+1, err)
    }
    return result
}

func (me MultiError) HasErrors() bool {
    return len(me.Errors) > 0
}

func (me *MultiError) Add(err error) {
    if err != nil {
        me.Errors = append(me.Errors, err)
    }
}

// Batch operation with error collection
func processBatch(items []string) error {
    var multiErr MultiError
    
    for i, item := range items {
        err := processItem(item)
        if err != nil {
            multiErr.Add(fmt.Errorf("item %d (%s): %w", i, item, err))
        }
    }
    
    if multiErr.HasErrors() {
        return multiErr
    }
    
    return nil
}

func processItem(item string) error {
    if item == "" {
        return ErrInvalidInput
    }
    if item == "notfound" {
        return ErrNotFound
    }
    if item == "timeout" {
        return ErrTimeout
    }
    if item == "unauthorized" {
        return ErrUnauthorized
    }
    
    fmt.Printf("Processed item: %s\n", item)
    return nil
}

// Simulated unreliable operation
var operationCount int

func unreliableOperation() error {
    operationCount++
    
    switch operationCount % 4 {
    case 1:
        return ErrTimeout
    case 2:
        return ErrTimeout
    case 3:
        return nil // Success
    default:
        return ErrTimeout
    }
}

func main() {
    fmt.Println("=== Error Handling Patterns ===")
    
    // 1. Retry pattern
    fmt.Println("\n1. Retry Pattern:")
    err := retryOperation(unreliableOperation, 3, 100*time.Millisecond)
    if err != nil {
        fmt.Printf("Retry failed: %v\n", err)
    } else {
        fmt.Println("Retry succeeded!")
    }
    
    // 2. Circuit breaker pattern
    fmt.Println("\n2. Circuit Breaker Pattern:")
    cb := NewCircuitBreaker(2, 1*time.Second)
    
    // Simulate multiple failures
    for i := 0; i < 5; i++ {
        err := cb.Call(func() error {
            return ErrTimeout
        })
        fmt.Printf("Call %d: %v\n", i+1, err)
    }
    
    // Wait for circuit breaker to reset
    fmt.Println("Waiting for circuit breaker timeout...")
    time.Sleep(1100 * time.Millisecond)
    
    // Try again
    err = cb.Call(func() error {
        return nil // Success
    })
    fmt.Printf("After timeout: %v\n", err)
    
    // 3. Error aggregation
    fmt.Println("\n3. Error Aggregation:")
    items := []string{"valid1", "", "valid2", "notfound", "valid3", "timeout", "unauthorized"}
    
    err = processBatch(items)
    if err != nil {
        fmt.Printf("Batch processing errors:\n%v\n", err)
        
        // Check if it's a MultiError
        var multiErr MultiError
        if errors.As(err, &multiErr) {
            fmt.Printf("Total errors: %d\n", len(multiErr.Errors))
        }
    } else {
        fmt.Println("Batch processing completed successfully")
    }
    
    // 4. Result type pattern
    fmt.Println("\n4. Result Type Pattern:")
    results := []Result[string]{
        {Value: "success", Error: nil},
        {Value: "", Error: ErrNotFound},
        {Value: "another success", Error: nil},
        {Value: "", Error: ErrInvalidInput},
    }
    
    for i, result := range results {
        if result.IsSuccess() {
            fmt.Printf("Result %d: Success - %s\n", i+1, result.Value)
        } else {
            fmt.Printf("Result %d: Error - %v\n", i+1, result.Error)
        }
    }
    
    // 5. Error type checking
    fmt.Println("\n5. Error Type Checking:")
    testErrors := []error{
        ErrNotFound,
        fmt.Errorf("wrapped not found: %w", ErrNotFound),
        ErrTimeout,
        errors.New("unknown error"),
    }
    
    for i, err := range testErrors {
        fmt.Printf("Error %d: %v\n", i+1, err)
        
        switch {
        case errors.Is(err, ErrNotFound):
            fmt.Println("  -> Handle not found error")
        case errors.Is(err, ErrTimeout):
            fmt.Println("  -> Handle timeout error")
        case errors.Is(err, ErrInvalidInput):
            fmt.Println("  -> Handle invalid input error")
        case errors.Is(err, ErrUnauthorized):
            fmt.Println("  -> Handle unauthorized error")
        default:
            fmt.Println("  -> Handle unknown error")
        }
    }
}
```

## Practical Examples

### Example 1: HTTP Client with Error Handling

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

// Custom error types
type HTTPError struct {
    StatusCode int
    Status     string
    URL        string
    Body       string
}

func (e HTTPError) Error() string {
    return fmt.Sprintf("HTTP %d %s for URL %s: %s", 
        e.StatusCode, e.Status, e.URL, e.Body)
}

type NetworkError struct {
    URL string
    Err error
}

func (e NetworkError) Error() string {
    return fmt.Sprintf("network error for URL %s: %v", e.URL, e.Err)
}

func (e NetworkError) Unwrap() error {
    return e.Err
}

type JSONError struct {
    Data string
    Err  error
}

func (e JSONError) Error() string {
    return fmt.Sprintf("JSON parsing error: %v", e.Err)
}

func (e JSONError) Unwrap() error {
    return e.Err
}

// HTTP client with error handling
type APIClient struct {
    BaseURL    string
    HTTPClient *http.Client
    MaxRetries int
    RetryDelay time.Duration
}

func NewAPIClient(baseURL string) *APIClient {
    return &APIClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        MaxRetries: 3,
        RetryDelay: 1 * time.Second,
    }
}

func (c *APIClient) Get(endpoint string) ([]byte, error) {
    url := c.BaseURL + endpoint
    
    var lastErr error
    for attempt := 0; attempt <= c.MaxRetries; attempt++ {
        if attempt > 0 {
            fmt.Printf("Retrying request to %s (attempt %d/%d)\n", 
                url, attempt, c.MaxRetries)
            time.Sleep(c.RetryDelay)
        }
        
        resp, err := c.HTTPClient.Get(url)
        if err != nil {
            lastErr = NetworkError{URL: url, Err: err}
            continue
        }
        defer resp.Body.Close()
        
        body, err := io.ReadAll(resp.Body)
        if err != nil {
            lastErr = NetworkError{URL: url, Err: err}
            continue
        }
        
        if resp.StatusCode >= 400 {
            httpErr := HTTPError{
                StatusCode: resp.StatusCode,
                Status:     resp.Status,
                URL:        url,
                Body:       string(body),
            }
            
            // Don't retry client errors (4xx)
            if resp.StatusCode >= 400 && resp.StatusCode < 500 {
                return nil, httpErr
            }
            
            lastErr = httpErr
            continue
        }
        
        return body, nil
    }
    
    return nil, fmt.Errorf("request failed after %d retries: %w", 
        c.MaxRetries, lastErr)
}

func (c *APIClient) GetJSON(endpoint string, target interface{}) error {
    data, err := c.Get(endpoint)
    if err != nil {
        return err
    }
    
    err = json.Unmarshal(data, target)
    if err != nil {
        return JSONError{
            Data: string(data),
            Err:  err,
        }
    }
    
    return nil
}

// Example API response structures
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type APIResponse struct {
    Success bool        `json:"success"`
    Data    interface{} `json:"data"`
    Error   string      `json:"error,omitempty"`
}

// Service layer with error handling
type UserService struct {
    client *APIClient
}

func NewUserService(baseURL string) *UserService {
    return &UserService{
        client: NewAPIClient(baseURL),
    }
}

func (s *UserService) GetUser(id int) (*User, error) {
    endpoint := fmt.Sprintf("/users/%d", id)
    
    var response APIResponse
    err := s.client.GetJSON(endpoint, &response)
    if err != nil {
        return nil, fmt.Errorf("failed to get user %d: %w", id, err)
    }
    
    if !response.Success {
        return nil, fmt.Errorf("API error: %s", response.Error)
    }
    
    // Convert response data to User
    userData, err := json.Marshal(response.Data)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal user data: %w", err)
    }
    
    var user User
    err = json.Unmarshal(userData, &user)
    if err != nil {
        return nil, JSONError{Data: string(userData), Err: err}
    }
    
    return &user, nil
}

func (s *UserService) GetUsers() ([]User, error) {
    var users []User
    err := s.client.GetJSON("/users", &users)
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }
    
    return users, nil
}

// Error handling middleware
func handleUserServiceError(err error) {
    if err == nil {
        return
    }
    
    fmt.Printf("Error occurred: %v\n", err)
    
    // Handle different error types
    var httpErr HTTPError
    var networkErr NetworkError
    var jsonErr JSONError
    
    switch {
    case errors.As(err, &httpErr):
        switch httpErr.StatusCode {
        case 404:
            fmt.Println("  -> Resource not found")
        case 401:
            fmt.Println("  -> Authentication required")
        case 403:
            fmt.Println("  -> Access forbidden")
        case 429:
            fmt.Println("  -> Rate limit exceeded")
        case 500:
            fmt.Println("  -> Server error - try again later")
        default:
            fmt.Printf("  -> HTTP error: %d\n", httpErr.StatusCode)
        }
        
    case errors.As(err, &networkErr):
        fmt.Println("  -> Network connectivity issue")
        fmt.Println("  -> Check internet connection")
        
    case errors.As(err, &jsonErr):
        fmt.Println("  -> Data format error")
        fmt.Println("  -> Server returned invalid JSON")
        
    default:
        fmt.Println("  -> Unknown error type")
    }
}

func main() {
    // Note: This example uses a mock API that doesn't exist
    // In a real scenario, you would use an actual API endpoint
    
    userService := NewUserService("https://jsonplaceholder.typicode.com")
    
    fmt.Println("=== HTTP Client Error Handling Example ===")
    
    // Test successful request
    fmt.Println("\n1. Testing successful request:")
    user, err := userService.GetUser(1)
    if err != nil {
        handleUserServiceError(err)
    } else {
        fmt.Printf("Success: %+v\n", user)
    }
    
    // Test not found error
    fmt.Println("\n2. Testing not found error:")
    user, err = userService.GetUser(999)
    if err != nil {
        handleUserServiceError(err)
    } else {
        fmt.Printf("Success: %+v\n", user)
    }
    
    // Test network error (invalid URL)
    fmt.Println("\n3. Testing network error:")
    invalidService := NewUserService("https://invalid-domain-that-does-not-exist.com")
    user, err = invalidService.GetUser(1)
    if err != nil {
        handleUserServiceError(err)
    } else {
        fmt.Printf("Success: %+v\n", user)
    }
    
    // Test getting multiple users
    fmt.Println("\n4. Testing get multiple users:")
    users, err := userService.GetUsers()
    if err != nil {
        handleUserServiceError(err)
    } else {
        fmt.Printf("Success: Retrieved %d users\n", len(users))
        for i, u := range users {
            if i < 3 { // Show first 3 users
                fmt.Printf("  User %d: %s (%s)\n", u.ID, u.Name, u.Email)
            }
        }
        if len(users) > 3 {
            fmt.Printf("  ... and %d more users\n", len(users)-3)
        }
    }
    
    // Demonstrate error chain analysis
    fmt.Println("\n5. Error chain analysis:")
    _, err = invalidService.GetUser(1)
    if err != nil {
        fmt.Println("Error chain:")
        current := err
        level := 1
        for current != nil {
            fmt.Printf("  %d. %v\n", level, current)
            current = errors.Unwrap(current)
            level++
        }
    }
}
```

### Example 2: File Processing with Comprehensive Error Handling

```go
package main

import (
    "bufio"
    "encoding/csv"
    "errors"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
)

// Custom error types for file processing
type FileProcessingError struct {
    Filename string
    Line     int
    Column   int
    Field    string
    Value    string
    Err      error
}

func (e FileProcessingError) Error() string {
    if e.Line > 0 {
        return fmt.Sprintf("file '%s' line %d: %v", e.Filename, e.Line, e.Err)
    }
    return fmt.Sprintf("file '%s': %v", e.Filename, e.Err)
}

func (e FileProcessingError) Unwrap() error {
    return e.Err
}

type ValidationError struct {
    Field   string
    Value   string
    Rule    string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field '%s' with value '%s': %s (rule: %s)",
        e.Field, e.Value, e.Message, e.Rule)
}

// Data structures
type Employee struct {
    ID       int
    Name     string
    Email    string
    Salary   float64
    Department string
}

type ProcessingResult struct {
    Employees []Employee
    Errors    []error
    Warnings  []string
    Stats     ProcessingStats
}

type ProcessingStats struct {
    TotalLines    int
    ProcessedLines int
    ErrorLines    int
    WarningLines  int
}

// File processor
type FileProcessor struct {
    validators map[string]func(string) error
}

func NewFileProcessor() *FileProcessor {
    fp := &FileProcessor{
        validators: make(map[string]func(string) error),
    }
    
    // Add default validators
    fp.AddValidator("email", validateEmail)
    fp.AddValidator("salary", validateSalary)
    fp.AddValidator("name", validateName)
    
    return fp
}

func (fp *FileProcessor) AddValidator(field string, validator func(string) error) {
    fp.validators[field] = validator
}

func validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return ValidationError{
            Field:   "email",
            Value:   email,
            Rule:    "must_contain_at",
            Message: "email must contain @ symbol",
        }
    }
    if len(email) < 5 {
        return ValidationError{
            Field:   "email",
            Value:   email,
            Rule:    "min_length",
            Message: "email must be at least 5 characters",
        }
    }
    return nil
}

func validateSalary(salaryStr string) error {
    salary, err := strconv.ParseFloat(salaryStr, 64)
    if err != nil {
        return ValidationError{
            Field:   "salary",
            Value:   salaryStr,
            Rule:    "numeric",
            Message: "salary must be a valid number",
        }
    }
    if salary < 0 {
        return ValidationError{
            Field:   "salary",
            Value:   salaryStr,
            Rule:    "positive",
            Message: "salary must be positive",
        }
    }
    if salary > 1000000 {
        return ValidationError{
            Field:   "salary",
            Value:   salaryStr,
            Rule:    "max_value",
            Message: "salary seems unreasonably high",
        }
    }
    return nil
}

func validateName(name string) error {
    if strings.TrimSpace(name) == "" {
        return ValidationError{
            Field:   "name",
            Value:   name,
            Rule:    "required",
            Message: "name cannot be empty",
        }
    }
    if len(name) < 2 {
        return ValidationError{
            Field:   "name",
            Value:   name,
            Rule:    "min_length",
            Message: "name must be at least 2 characters",
        }
    }
    return nil
}

func (fp *FileProcessor) ProcessCSVFile(filename string) (*ProcessingResult, error) {
    // Check if file exists
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        return nil, FileProcessingError{
            Filename: filename,
            Err:      errors.New("file does not exist"),
        }
    }
    
    file, err := os.Open(filename)
    if err != nil {
        return nil, FileProcessingError{
            Filename: filename,
            Err:      fmt.Errorf("failed to open file: %w", err),
        }
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    reader.FieldsPerRecord = -1 // Allow variable number of fields
    
    result := &ProcessingResult{
        Employees: make([]Employee, 0),
        Errors:    make([]error, 0),
        Warnings:  make([]string, 0),
    }
    
    lineNumber := 0
    headerProcessed := false
    
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        
        lineNumber++
        result.Stats.TotalLines++
        
        if err != nil {
            procErr := FileProcessingError{
                Filename: filename,
                Line:     lineNumber,
                Err:      fmt.Errorf("CSV parsing error: %w", err),
            }
            result.Errors = append(result.Errors, procErr)
            result.Stats.ErrorLines++
            continue
        }
        
        // Skip header
        if !headerProcessed {
            headerProcessed = true
            continue
        }
        
        employee, warnings, err := fp.processRecord(record, filename, lineNumber)
        if err != nil {
            result.Errors = append(result.Errors, err)
            result.Stats.ErrorLines++
            continue
        }
        
        if len(warnings) > 0 {
            for _, warning := range warnings {
                result.Warnings = append(result.Warnings, 
                    fmt.Sprintf("Line %d: %s", lineNumber, warning))
            }
            result.Stats.WarningLines++
        }
        
        result.Employees = append(result.Employees, *employee)
        result.Stats.ProcessedLines++
    }
    
    return result, nil
}

func (fp *FileProcessor) processRecord(record []string, filename string, lineNumber int) (*Employee, []string, error) {
    if len(record) < 5 {
        return nil, nil, FileProcessingError{
            Filename: filename,
            Line:     lineNumber,
            Err:      fmt.Errorf("insufficient fields: expected 5, got %d", len(record)),
        }
    }
    
    var warnings []string
    
    // Parse ID
    id, err := strconv.Atoi(strings.TrimSpace(record[0]))
    if err != nil {
        return nil, nil, FileProcessingError{
            Filename: filename,
            Line:     lineNumber,
            Field:    "id",
            Value:    record[0],
            Err:      fmt.Errorf("invalid ID: %w", err),
        }
    }
    
    // Validate and process other fields
    name := strings.TrimSpace(record[1])
    if validator, exists := fp.validators["name"]; exists {
        if err := validator(name); err != nil {
            return nil, nil, FileProcessingError{
                Filename: filename,
                Line:     lineNumber,
                Field:    "name",
                Value:    name,
                Err:      err,
            }
        }
    }
    
    email := strings.TrimSpace(record[2])
    if validator, exists := fp.validators["email"]; exists {
        if err := validator(email); err != nil {
            return nil, nil, FileProcessingError{
                Filename: filename,
                Line:     lineNumber,
                Field:    "email",
                Value:    email,
                Err:      err,
            }
        }
    }
    
    // Parse salary
    salaryStr := strings.TrimSpace(record[3])
    if validator, exists := fp.validators["salary"]; exists {
        if err := validator(salaryStr); err != nil {
            return nil, nil, FileProcessingError{
                Filename: filename,
                Line:     lineNumber,
                Field:    "salary",
                Value:    salaryStr,
                Err:      err,
            }
        }
    }
    
    salary, err := strconv.ParseFloat(salaryStr, 64)
    if err != nil {
        return nil, nil, FileProcessingError{
            Filename: filename,
            Line:     lineNumber,
            Field:    "salary",
            Value:    salaryStr,
            Err:      fmt.Errorf("invalid salary: %w", err),
        }
    }
    
    department := strings.TrimSpace(record[4])
    if department == "" {
        warnings = append(warnings, "department is empty")
        department = "Unknown"
    }
    
    // Check for duplicate ID (simplified check)
    if id <= 0 {
        warnings = append(warnings, "ID should be positive")
    }
    
    employee := &Employee{
        ID:         id,
        Name:       name,
        Email:      email,
        Salary:     salary,
        Department: department,
    }
    
    return employee, warnings, nil
}

func (fp *FileProcessor) ProcessDirectory(dirPath string) (map[string]*ProcessingResult, error) {
    files, err := filepath.Glob(filepath.Join(dirPath, "*.csv"))
    if err != nil {
        return nil, fmt.Errorf("failed to list CSV files in directory %s: %w", dirPath, err)
    }
    
    if len(files) == 0 {
        return nil, fmt.Errorf("no CSV files found in directory %s", dirPath)
    }
    
    results := make(map[string]*ProcessingResult)
    
    for _, file := range files {
        fmt.Printf("Processing file: %s\n", file)
        result, err := fp.ProcessCSVFile(file)
        if err != nil {
            fmt.Printf("Failed to process file %s: %v\n", file, err)
            continue
        }
        results[file] = result
    }
    
    return results, nil
}

func createTestFiles() error {
    // Create test directory
    err := os.MkdirAll("test_data", 0755)
    if err != nil {
        return err
    }
    
    // Valid CSV file
    validCSV := `ID,Name,Email,Salary,Department
1,John Doe,john@example.com,50000,Engineering
2,Jane Smith,jane@example.com,60000,Marketing
3,Bob Johnson,bob@example.com,55000,Engineering
4,Alice Brown,alice@example.com,65000,Sales`
    
    err = os.WriteFile("test_data/valid_employees.csv", []byte(validCSV), 0644)
    if err != nil {
        return err
    }
    
    // CSV with errors
    errorCSV := `ID,Name,Email,Salary,Department
1,John Doe,john@example.com,50000,Engineering
invalid_id,Jane Smith,jane@example.com,60000,Marketing
3,,bob@example.com,55000,Engineering
4,Alice Brown,invalid_email,65000,Sales
5,Bob Wilson,bob@example.com,-1000,Finance
6,Carol Davis,carol@example.com,abc,HR`
    
    err = os.WriteFile("test_data/error_employees.csv", []byte(errorCSV), 0644)
    if err != nil {
        return err
    }
    
    return nil
}

func printProcessingResult(filename string, result *ProcessingResult) {
    fmt.Printf("\n=== Processing Result for %s ===\n", filename)
    fmt.Printf("Statistics:\n")
    fmt.Printf("  Total lines: %d\n", result.Stats.TotalLines)
    fmt.Printf("  Processed successfully: %d\n", result.Stats.ProcessedLines)
    fmt.Printf("  Lines with errors: %d\n", result.Stats.ErrorLines)
    fmt.Printf("  Lines with warnings: %d\n", result.Stats.WarningLines)
    
    if len(result.Employees) > 0 {
        fmt.Printf("\nSuccessfully processed employees:\n")
        for _, emp := range result.Employees {
            fmt.Printf("  ID: %d, Name: %s, Email: %s, Salary: %.2f, Dept: %s\n",
                emp.ID, emp.Name, emp.Email, emp.Salary, emp.Department)
        }
    }
    
    if len(result.Errors) > 0 {
        fmt.Printf("\nErrors encountered:\n")
        for i, err := range result.Errors {
            fmt.Printf("  %d. %v\n", i+1, err)
            
            // Show error type information
            var fileErr FileProcessingError
            var validErr ValidationError
            
            if errors.As(err, &fileErr) {
                if fileErr.Field != "" {
                    fmt.Printf("     Field: %s, Value: %s\n", fileErr.Field, fileErr.Value)
                }
            }
            
            if errors.As(err, &validErr) {
                fmt.Printf("     Validation rule: %s\n", validErr.Rule)
            }
        }
    }
    
    if len(result.Warnings) > 0 {
        fmt.Printf("\nWarnings:\n")
        for i, warning := range result.Warnings {
            fmt.Printf("  %d. %s\n", i+1, warning)
        }
    }
}

func main() {
    fmt.Println("=== File Processing Error Handling Example ===")
    
    // Create test files
    err := createTestFiles()
    if err != nil {
        fmt.Printf("Failed to create test files: %v\n", err)
        return
    }
    defer os.RemoveAll("test_data") // Clean up
    
    processor := NewFileProcessor()
    
    // Process valid file
    fmt.Println("\n1. Processing valid CSV file:")
    result, err := processor.ProcessCSVFile("test_data/valid_employees.csv")
    if err != nil {
        fmt.Printf("Failed to process file: %v\n", err)
    } else {
        printProcessingResult("valid_employees.csv", result)
    }
    
    // Process file with errors
    fmt.Println("\n2. Processing CSV file with errors:")
    result, err = processor.ProcessCSVFile("test_data/error_employees.csv")
    if err != nil {
        fmt.Printf("Failed to process file: %v\n", err)
    } else {
        printProcessingResult("error_employees.csv", result)
    }
    
    // Process non-existent file
    fmt.Println("\n3. Processing non-existent file:")
    result, err = processor.ProcessCSVFile("test_data/nonexistent.csv")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        
        var fileErr FileProcessingError
        if errors.As(err, &fileErr) {
            fmt.Printf("File error for: %s\n", fileErr.Filename)
        }
    }
    
    // Process entire directory
    fmt.Println("\n4. Processing entire directory:")
    results, err := processor.ProcessDirectory("test_data")
    if err != nil {
        fmt.Printf("Failed to process directory: %v\n", err)
    } else {
        fmt.Printf("Processed %d files\n", len(results))
        for filename, result := range results {
            fmt.Printf("\nFile: %s - Processed: %d, Errors: %d, Warnings: %d\n",
                filepath.Base(filename), result.Stats.ProcessedLines, 
                result.Stats.ErrorLines, result.Stats.WarningLines)
        }
    }
}
```

## Exercises

### Exercise 1: Database Connection Manager
Create a database connection manager with comprehensive error handling.

```go
package main

import (
    "errors"
    "time"
)

// Define custom error types:
// ConnectionError, QueryError, TransactionError, TimeoutError

type DatabaseManager struct {
    // Add fields for connection pool, timeouts, retry logic
}

// Implement these methods with proper error handling:
// Connect() error
// Disconnect() error
// Query(sql string, args ...interface{}) ([]map[string]interface{}, error)
// Execute(sql string, args ...interface{}) (int64, error)
// BeginTransaction() (Transaction, error)
// WithRetry(operation func() error, maxRetries int) error

type Transaction struct {
    // Add transaction fields
}

// Implement transaction methods:
// Commit() error
// Rollback() error
// Query(sql string, args ...interface{}) ([]map[string]interface{}, error)
// Execute(sql string, args ...interface{}) (int64, error)

func main() {
    // Test the database manager with various error scenarios
}
```

### Exercise 2: Configuration Loader
Create a configuration loader with validation and error reporting.

```go
package main

import "time"

// Define error types:
// ConfigError, ValidationError, ParseError

type Config struct {
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Logging  LoggingConfig  `json:"logging"`
}

type DatabaseConfig struct {
    Host     string `json:"host" validate:"required"`
    Port     int    `json:"port" validate:"range:1-65535"`
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
    Database string `json:"database" validate:"required"`
}

type ServerConfig struct {
    Host         string        `json:"host" validate:"required"`
    Port         int           `json:"port" validate:"range:1-65535"`
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
}

type LoggingConfig struct {
    Level  string `json:"level" validate:"oneof:debug info warn error"`
    Format string `json:"format" validate:"oneof:json text"`
    Output string `json:"output" validate:"required"`
}

type ConfigLoader struct {
    // Add fields for validation rules, file paths, etc.
}

// Implement these methods:
// LoadFromFile(filename string) (*Config, error)
// LoadFromJSON(data []byte) (*Config, error)
// Validate(config *Config) error
// LoadWithDefaults(filename string) (*Config, error)

func main() {
    // Test configuration loading with various error scenarios
}
```

### Exercise 3: Batch Job Processor
Create a batch job processor with error recovery and reporting.

```go
package main

import "time"

// Define error types:
// JobError, BatchError, ProcessingError

type Job struct {
    ID       string
    Type     string
    Data     interface{}
    Priority int
    Retries  int
    MaxRetries int
}

type JobResult struct {
    Job       Job
    Success   bool
    Error     error
    Duration  time.Duration
    Timestamp time.Time
}

type BatchProcessor struct {
    // Add fields for job queue, workers, error handling
}

// Implement these methods:
// AddJob(job Job) error
// ProcessBatch(jobs []Job) ([]JobResult, error)
// ProcessWithRecovery(job Job) JobResult
// GetFailedJobs() []Job
// RetryFailedJobs() error
// GenerateErrorReport() string

func main() {
    // Test batch processing with various error scenarios
}
```

## Key Takeaways

1. **Explicit error handling** makes code more robust and predictable
2. **Error interface** is simple but powerful: `Error() string`
3. **Custom error types** provide rich context and type safety
4. **Error wrapping** with `fmt.Errorf("%w", err)` preserves error chains
5. **errors.Is()** checks for specific errors in the chain
6. **errors.As()** extracts specific error types from the chain
7. **Sentinel errors** are predefined error values for common cases
8. **Error aggregation** collects multiple errors for batch operations
9. **Retry patterns** handle transient failures gracefully
10. **Circuit breakers** prevent cascading failures

## Next Steps

Now that you understand error handling, let's learn about [Packages and Modules](14-packages-modules.md) for organizing and sharing Go code!

---

**Previous**: [← Interfaces](12-interfaces.md) | **Next**: [Packages and Modules →](14-packages-modules.md)