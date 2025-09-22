# 18. Context Package

The `context` package in Go provides a way to carry deadlines, cancellation signals, and request-scoped values across API boundaries and between processes. It's particularly useful for controlling the execution of goroutines, managing timeouts, and passing request-scoped values through call chains.

## Context Basics

### What is Context?

A `context.Context` is an interface that provides four main methods:

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

- `Deadline()`: Returns the time when the context will be canceled (if any)
- `Done()`: Returns a channel that's closed when the context is canceled
- `Err()`: Returns the error why the context was canceled (if any)
- `Value()`: Returns the value associated with a key (if any)

### Creating Contexts

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func demonstrateContextCreation() {
    fmt.Println("=== Context Creation ===")
    
    // Background context - root of all contexts
    bgCtx := context.Background()
    fmt.Printf("Background context: %v\n", bgCtx)
    
    // TODO context - placeholder for when you're not sure which context to use
    todoCtx := context.TODO()
    fmt.Printf("TODO context: %v\n", todoCtx)
    
    // WithCancel context
    cancelCtx, cancel := context.WithCancel(bgCtx)
    fmt.Printf("WithCancel context: %v\n", cancelCtx)
    defer cancel() // Always call cancel to release resources
    
    // WithTimeout context
    timeoutCtx, cancel := context.WithTimeout(bgCtx, 5*time.Second)
    fmt.Printf("WithTimeout context: %v\n", timeoutCtx)
    defer cancel()
    
    // WithDeadline context
    deadline := time.Now().Add(10 * time.Second)
    deadlineCtx, cancel := context.WithDeadline(bgCtx, deadline)
    fmt.Printf("WithDeadline context: %v\n", deadlineCtx)
    defer cancel()
    
    // WithValue context
    valueCtx := context.WithValue(bgCtx, "key", "value")
    fmt.Printf("WithValue context: %v\n", valueCtx)
    fmt.Printf("Value from context: %v\n", valueCtx.Value("key"))
}

func main() {
    demonstrateContextCreation()
}
```

## Context Cancellation

### Manual Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    fmt.Printf("Worker %d: Started\n", id)
    
    select {
    case <-time.After(5 * time.Second):
        fmt.Printf("Worker %d: Completed work\n", id)
    case <-ctx.Done():
        fmt.Printf("Worker %d: Canceled with error: %v\n", id, ctx.Err())
    }
}

func demonstrateManualCancellation() {
    fmt.Println("=== Manual Cancellation ===")
    
    // Create a context with cancel function
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }
    
    // Wait a bit, then cancel all workers
    fmt.Println("Main: Starting workers...")
    time.Sleep(2 * time.Second)
    fmt.Println("Main: Canceling workers...")
    cancel()
    
    // Give workers time to handle cancellation
    time.Sleep(1 * time.Second)
    fmt.Println("Main: Done")
}

func main() {
    demonstrateManualCancellation()
}
```

### Timeout Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func slowOperation(ctx context.Context) (string, error) {
    fmt.Println("SlowOperation: Starting...")
    
    select {
    case <-time.After(3 * time.Second):
        fmt.Println("SlowOperation: Completed successfully")
        return "Operation result", nil
    case <-ctx.Done():
        fmt.Printf("SlowOperation: Canceled with error: %v\n", ctx.Err())
        return "", ctx.Err()
    }
}

func demonstrateTimeoutCancellation() {
    fmt.Println("=== Timeout Cancellation ===")
    
    // Case 1: Timeout is longer than operation
    fmt.Println("\nCase 1: Timeout longer than operation")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel1()
    
    result1, err1 := slowOperation(ctx1)
    if err1 != nil {
        fmt.Printf("Error: %v\n", err1)
    } else {
        fmt.Printf("Result: %s\n", result1)
    }
    
    // Case 2: Timeout is shorter than operation
    fmt.Println("\nCase 2: Timeout shorter than operation")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel2()
    
    result2, err2 := slowOperation(ctx2)
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
    } else {
        fmt.Printf("Result: %s\n", result2)
    }
}

func main() {
    demonstrateTimeoutCancellation()
}
```

### Deadline Cancellation

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func processWithDeadline(ctx context.Context, taskID string) error {
    // Check if we already have a deadline
    deadline, ok := ctx.Deadline()
    if ok {
        fmt.Printf("Task %s: Deadline set to %v (in %v)\n", 
            taskID, deadline, time.Until(deadline))
    } else {
        fmt.Printf("Task %s: No deadline set\n", taskID)
    }
    
    // Simulate work
    select {
    case <-time.After(2 * time.Second):
        fmt.Printf("Task %s: Completed successfully\n", taskID)
        return nil
    case <-ctx.Done():
        fmt.Printf("Task %s: Canceled due to: %v\n", taskID, ctx.Err())
        return ctx.Err()
    }
}

func demonstrateDeadlineCancellation() {
    fmt.Println("=== Deadline Cancellation ===")
    
    // Case 1: Deadline in the future
    fmt.Println("\nCase 1: Deadline in the future")
    deadline1 := time.Now().Add(3 * time.Second)
    ctx1, cancel1 := context.WithDeadline(context.Background(), deadline1)
    defer cancel1()
    
    err1 := processWithDeadline(ctx1, "future-deadline")
    if err1 != nil {
        fmt.Printf("Error: %v\n", err1)
    }
    
    // Case 2: Deadline in the past
    fmt.Println("\nCase 2: Deadline in the past")
    deadline2 := time.Now().Add(-1 * time.Second) // 1 second in the past
    ctx2, cancel2 := context.WithDeadline(context.Background(), deadline2)
    defer cancel2()
    
    err2 := processWithDeadline(ctx2, "past-deadline")
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
    }
    
    // Case 3: Derived context with earlier deadline
    fmt.Println("\nCase 3: Derived context with earlier deadline")
    parentDeadline := time.Now().Add(10 * time.Second)
    parentCtx, parentCancel := context.WithDeadline(context.Background(), parentDeadline)
    defer parentCancel()
    
    childDeadline := time.Now().Add(1 * time.Second)
    childCtx, childCancel := context.WithDeadline(parentCtx, childDeadline)
    defer childCancel()
    
    err3 := processWithDeadline(childCtx, "child-deadline")
    if err3 != nil {
        fmt.Printf("Error: %v\n", err3)
    }
}

func main() {
    demonstrateDeadlineCancellation()
}
```

## Context Values

### Passing Values Through Context

```go
package main

import (
    "context"
    "fmt"
)

// Custom key types to avoid collisions
type contextKey string

const (
    userIDKey   contextKey = "userID"
    authTokenKey contextKey = "authToken"
    requestIDKey contextKey = "requestID"
)

func processRequest(ctx context.Context) {
    // Extract values from context
    userID, ok := ctx.Value(userIDKey).(string)
    if !ok {
        userID = "anonymous"
    }
    
    authToken, _ := ctx.Value(authTokenKey).(string)
    requestID, _ := ctx.Value(requestIDKey).(string)
    
    fmt.Printf("Processing request:\n")
    fmt.Printf("  RequestID: %s\n", requestID)
    fmt.Printf("  UserID: %s\n", userID)
    fmt.Printf("  AuthToken: %s\n", authToken)
    
    // Call a sub-function with the same context
    validatePermissions(ctx)
}

func validatePermissions(ctx context.Context) {
    // Extract values from context
    userID, _ := ctx.Value(userIDKey).(string)
    requestID, _ := ctx.Value(requestIDKey).(string)
    
    fmt.Printf("Validating permissions:\n")
    fmt.Printf("  RequestID: %s\n", requestID)
    fmt.Printf("  UserID: %s\n", userID)
    
    // Call database function with the same context
    queryDatabase(ctx)
}

func queryDatabase(ctx context.Context) {
    // Extract values from context
    userID, _ := ctx.Value(userIDKey).(string)
    requestID, _ := ctx.Value(requestIDKey).(string)
    
    fmt.Printf("Querying database:\n")
    fmt.Printf("  RequestID: %s\n", requestID)
    fmt.Printf("  UserID: %s\n", userID)
}

func demonstrateContextValues() {
    fmt.Println("=== Context Values ===")
    
    // Create a base context
    ctx := context.Background()
    
    // Add values to the context chain
    ctx = context.WithValue(ctx, requestIDKey, "req-123")
    ctx = context.WithValue(ctx, userIDKey, "user-456")
    ctx = context.WithValue(ctx, authTokenKey, "token-789")
    
    // Process the request with the context
    processRequest(ctx)
    
    // Create a different context with different values
    fmt.Println("\nProcessing another request:")
    ctx2 := context.Background()
    ctx2 = context.WithValue(ctx2, requestIDKey, "req-999")
    ctx2 = context.WithValue(ctx2, userIDKey, "admin-user")
    
    processRequest(ctx2)
}

func main() {
    demonstrateContextValues()
}
```

### Best Practices for Context Values

```go
package main

import (
    "context"
    "fmt"
    "log"
)

// Good: Use custom types for context keys
type contextKey int

const (
    userKey contextKey = iota
    traceIDKey
    loggerKey
)

// Good: Define accessor functions
func GetUserFromContext(ctx context.Context) (string, bool) {
    user, ok := ctx.Value(userKey).(string)
    return user, ok
}

func GetTraceIDFromContext(ctx context.Context) (string, bool) {
    traceID, ok := ctx.Value(traceIDKey).(string)
    return traceID, ok
}

func GetLoggerFromContext(ctx context.Context) (*log.Logger, bool) {
    logger, ok := ctx.Value(loggerKey).(*log.Logger)
    return logger, ok
}

// Good: Define setter functions
func WithUser(ctx context.Context, user string) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, traceIDKey, traceID)
}

func WithLogger(ctx context.Context, logger *log.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}

// Bad: Using string keys (collision risk)
func badContextUsage() {
    ctx := context.Background()
    
    // Bad: Using string as key
    ctx = context.WithValue(ctx, "user", "john")
    
    // Bad: Using primitive type as value without type checking
    user := ctx.Value("user") // No type safety
    fmt.Printf("Bad usage - User: %v\n", user)
}

// Good: Using custom types and accessors
func goodContextUsage() {
    ctx := context.Background()
    
    // Create a logger
    logger := log.New(log.Writer(), "[REQUEST] ", log.LstdFlags)
    
    // Good: Using setter functions
    ctx = WithUser(ctx, "john")
    ctx = WithTraceID(ctx, "trace-123")
    ctx = WithLogger(ctx, logger)
    
    // Good: Using accessor functions with type safety
    if user, ok := GetUserFromContext(ctx); ok {
        fmt.Printf("Good usage - User: %s\n", user)
    }
    
    if traceID, ok := GetTraceIDFromContext(ctx); ok {
        fmt.Printf("Good usage - TraceID: %s\n", traceID)
    }
    
    if logger, ok := GetLoggerFromContext(ctx); ok {
        logger.Println("This is a log message with context")
    }
}

func demonstrateContextBestPractices() {
    fmt.Println("=== Context Values Best Practices ===")
    
    badContextUsage()
    fmt.Println()
    goodContextUsage()
}

func main() {
    demonstrateContextBestPractices()
}
```

## Context in HTTP Servers

### HTTP Server with Context

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
)

// Middleware to add request ID to context
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Generate a request ID
        requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
        
        // Add it to the context
        ctx := context.WithValue(r.Context(), "requestID", requestID)
        
        // Add it to response headers
        w.Header().Set("X-Request-ID", requestID)
        
        // Call the next handler with the updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware to add timeout to context
func timeoutMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Create a context with timeout
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()
        
        // Create a channel to signal when the handler is done
        done := make(chan struct{})
        
        // Create a response writer that can detect if headers were written
        var responseWritten bool
        
        go func() {
            // Call the next handler with the timeout context
            next.ServeHTTP(w, r.WithContext(ctx))
            responseWritten = true
            close(done)
        }()
        
        select {
        case <-done:
            // Handler finished before timeout
            return
        case <-ctx.Done():
            // Timeout occurred
            if !responseWritten {
                w.WriteHeader(http.StatusGatewayTimeout)
                w.Write([]byte("Request timed out"))
            }
            return
        }
    })
}

// Handler that uses context values
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Get request ID from context
    requestID, ok := r.Context().Value("requestID").(string)
    if !ok {
        requestID = "unknown"
    }
    
    log.Printf("[%s] Request started", requestID)
    
    // Simulate work with context awareness
    select {
    case <-time.After(2 * time.Second):
        log.Printf("[%s] Work completed", requestID)
        fmt.Fprintf(w, "Request processed successfully. Request ID: %s", requestID)
    case <-r.Context().Done():
        // This will happen if the client disconnects or the timeout is reached
        log.Printf("[%s] Request canceled: %v", requestID, r.Context().Err())
        return
    }
}

// Slow handler that will trigger timeout
func handleSlowRequest(w http.ResponseWriter, r *http.Request) {
    requestID, _ := r.Context().Value("requestID").(string)
    log.Printf("[%s] Slow request started", requestID)
    
    select {
    case <-time.After(10 * time.Second): // This is longer than our timeout
        log.Printf("[%s] Slow work completed", requestID)
        fmt.Fprintf(w, "Slow request processed. Request ID: %s", requestID)
    case <-r.Context().Done():
        log.Printf("[%s] Slow request canceled: %v", requestID, r.Context().Err())
        return
    }
}

func setupHTTPServer() {
    // Create a mux for routing
    mux := http.NewServeMux()
    
    // Register handlers
    mux.HandleFunc("/", handleRequest)
    mux.HandleFunc("/slow", handleSlowRequest)
    
    // Apply middleware (in order)
    handler := requestIDMiddleware(timeoutMiddleware(mux))
    
    // Start the server
    server := &http.Server{
        Addr:    ":8080",
        Handler: handler,
    }
    
    log.Println("Starting server on :8080")
    log.Println("Try these endpoints:")
    log.Println("  - http://localhost:8080/ (normal request)")
    log.Println("  - http://localhost:8080/slow (will timeout)")
    
    if err := server.ListenAndServe(); err != nil {
        log.Fatalf("Server error: %v", err)
    }
}

// Note: This example won't actually run in the playground
// To run this, save it as a file and execute it locally
func main() {
    setupHTTPServer()
}
```

## Context in HTTP Clients

### HTTP Client with Context

```go
package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

// Function to make HTTP request with context
func fetchURL(ctx context.Context, url string) (string, error) {
    // Create a new request
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return "", fmt.Errorf("error creating request: %w", err)
    }
    
    // Execute the request
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("error executing request: %w", err)
    }
    defer resp.Body.Close()
    
    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("error reading response: %w", err)
    }
    
    return string(body), nil
}

func demonstrateHTTPClientWithContext() {
    fmt.Println("=== HTTP Client with Context ===")
    
    // Case 1: Request with sufficient timeout
    fmt.Println("\nCase 1: Request with sufficient timeout")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel1()
    
    result1, err1 := fetchURL(ctx1, "https://httpbin.org/delay/1")
    if err1 != nil {
        fmt.Printf("Error: %v\n", err1)
    } else {
        fmt.Printf("Response length: %d bytes\n", len(result1))
    }
    
    // Case 2: Request with insufficient timeout
    fmt.Println("\nCase 2: Request with insufficient timeout")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel2()
    
    result2, err2 := fetchURL(ctx2, "https://httpbin.org/delay/2")
    if err2 != nil {
        fmt.Printf("Error: %v\n", err2)
    } else {
        fmt.Printf("Response length: %d bytes\n", len(result2))
    }
    
    // Case 3: Request with manual cancellation
    fmt.Println("\nCase 3: Request with manual cancellation")
    ctx3, cancel3 := context.WithCancel(context.Background())
    
    // Cancel the request after 300ms
    go func() {
        time.Sleep(300 * time.Millisecond)
        fmt.Println("Manually canceling request...")
        cancel3()
    }()
    
    result3, err3 := fetchURL(ctx3, "https://httpbin.org/delay/3")
    if err3 != nil {
        fmt.Printf("Error: %v\n", err3)
    } else {
        fmt.Printf("Response length: %d bytes\n", len(result3))
    }
}

// Note: This example won't work in the playground due to network restrictions
// To run this, save it as a file and execute it locally
func main() {
    demonstrateHTTPClientWithContext()
}
```

## Context Propagation

### Propagating Context Through Function Calls

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Trace represents a hierarchical trace of function calls
type Trace struct {
    Name     string
    Start    time.Time
    End      time.Time
    Children []*Trace
}

// TraceKey is the context key for the current trace
type traceKey struct{}

// StartTrace begins a new trace and adds it to the context
func StartTrace(ctx context.Context, name string) (context.Context, *Trace) {
    trace := &Trace{
        Name:     name,
        Start:    time.Now(),
        Children: make([]*Trace, 0),
    }
    
    return context.WithValue(ctx, traceKey{}, trace), trace
}

// GetTraceFromContext retrieves the current trace from context
func GetTraceFromContext(ctx context.Context) (*Trace, bool) {
    trace, ok := ctx.Value(traceKey{}).(*Trace)
    return trace, ok
}

// AddChild adds a child trace to the parent trace
func AddChild(parent *Trace, child *Trace) {
    parent.Children = append(parent.Children, child)
}

// EndTrace completes a trace by setting its end time
func EndTrace(trace *Trace) {
    trace.End = time.Now()
}

// PrintTrace prints a trace hierarchy with indentation
func PrintTrace(trace *Trace, indent string) {
    duration := trace.End.Sub(trace.Start)
    fmt.Printf("%s- %s: %v\n", indent, trace.Name, duration)
    
    for _, child := range trace.Children {
        PrintTrace(child, indent+"  ")
    }
}

// Function A calls B and C
func functionA(ctx context.Context) error {
    // Start a trace for function A
    ctx, trace := StartTrace(ctx, "FunctionA")
    defer EndTrace(trace)
    
    // Simulate some work
    time.Sleep(100 * time.Millisecond)
    
    // Call function B
    if err := functionB(ctx); err != nil {
        return err
    }
    
    // Call function C
    if err := functionC(ctx); err != nil {
        return err
    }
    
    return nil
}

// Function B calls D
func functionB(ctx context.Context) error {
    // Get parent trace and create a child trace
    parentTrace, ok := GetTraceFromContext(ctx)
    if !ok {
        return fmt.Errorf("no trace in context")
    }
    
    // Start a trace for function B
    childTrace := &Trace{
        Name:     "FunctionB",
        Start:    time.Now(),
        Children: make([]*Trace, 0),
    }
    AddChild(parentTrace, childTrace)
    
    // Create a new context with this trace
    ctx = context.WithValue(ctx, traceKey{}, childTrace)
    
    // Simulate some work
    time.Sleep(200 * time.Millisecond)
    
    // Call function D
    if err := functionD(ctx); err != nil {
        EndTrace(childTrace)
        return err
    }
    
    EndTrace(childTrace)
    return nil
}

// Function C is a leaf function
func functionC(ctx context.Context) error {
    // Get parent trace and create a child trace
    parentTrace, ok := GetTraceFromContext(ctx)
    if !ok {
        return fmt.Errorf("no trace in context")
    }
    
    // Start a trace for function C
    childTrace := &Trace{
        Name:     "FunctionC",
        Start:    time.Now(),
        Children: make([]*Trace, 0),
    }
    AddChild(parentTrace, childTrace)
    
    // Simulate some work
    time.Sleep(150 * time.Millisecond)
    
    // Check for cancellation
    select {
    case <-ctx.Done():
        childTrace.Name += " (canceled)"
        EndTrace(childTrace)
        return ctx.Err()
    default:
        // Continue processing
    }
    
    EndTrace(childTrace)
    return nil
}

// Function D is a leaf function
func functionD(ctx context.Context) error {
    // Get parent trace and create a child trace
    parentTrace, ok := GetTraceFromContext(ctx)
    if !ok {
        return fmt.Errorf("no trace in context")
    }
    
    // Start a trace for function D
    childTrace := &Trace{
        Name:     "FunctionD",
        Start:    time.Now(),
        Children: make([]*Trace, 0),
    }
    AddChild(parentTrace, childTrace)
    
    // Simulate some work
    time.Sleep(300 * time.Millisecond)
    
    EndTrace(childTrace)
    return nil
}

func demonstrateContextPropagation() {
    fmt.Println("=== Context Propagation ===")
    
    // Create a root context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // Start the function chain
    rootTrace := &Trace{
        Name:     "Main",
        Start:    time.Now(),
        Children: make([]*Trace, 0),
    }
    
    ctx = context.WithValue(ctx, traceKey{}, rootTrace)
    
    err := functionA(ctx)
    EndTrace(rootTrace)
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // Print the trace hierarchy
    fmt.Println("\nTrace Hierarchy:")
    PrintTrace(rootTrace, "")
}

func main() {
    demonstrateContextPropagation()
}
```

## Advanced Context Patterns

### Context with Retry Logic

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "time"
)

// RetryOptions configures the retry behavior
type RetryOptions struct {
    MaxRetries  int
    InitialWait time.Duration
    MaxWait     time.Duration
    Multiplier  float64
}

// DefaultRetryOptions provides sensible defaults
func DefaultRetryOptions() RetryOptions {
    return RetryOptions{
        MaxRetries:  3,
        InitialWait: 100 * time.Millisecond,
        MaxWait:     2 * time.Second,
        Multiplier:  2.0,
    }
}

// RetryWithContext retries a function with exponential backoff
func RetryWithContext(ctx context.Context, fn func(context.Context) (interface{}, error), opts RetryOptions) (interface{}, error) {
    var lastErr error
    wait := opts.InitialWait
    
    for attempt := 0; attempt <= opts.MaxRetries; attempt++ {
        // Check if context is canceled before making an attempt
        if ctx.Err() != nil {
            return nil, ctx.Err()
        }
        
        // If this isn't the first attempt, log the retry
        if attempt > 0 {
            fmt.Printf("Retry attempt %d after %v\n", attempt, wait)
            
            // Wait with exponential backoff, but respect context cancellation
            select {
            case <-time.After(wait):
                // Continue with retry
            case <-ctx.Done():
                return nil, ctx.Err()
            }
            
            // Increase wait time for next attempt
            wait = time.Duration(float64(wait) * opts.Multiplier)
            if wait > opts.MaxWait {
                wait = opts.MaxWait
            }
        }
        
        // Attempt the function
        result, err := fn(ctx)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        fmt.Printf("Attempt %d failed: %v\n", attempt, err)
    }
    
    return nil, fmt.Errorf("all %d attempts failed, last error: %w", opts.MaxRetries+1, lastErr)
}

// Simulate an unreliable function that sometimes fails
func unreliableFunction(ctx context.Context) (interface{}, error) {
    // Simulate random failures
    if rand.Float64() < 0.7 {
        return nil, fmt.Errorf("random failure")
    }
    
    // Simulate work
    select {
    case <-time.After(200 * time.Millisecond):
        return "success result", nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func demonstrateRetryPattern() {
    fmt.Println("=== Context with Retry Pattern ===")
    
    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())
    
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Set up retry options
    opts := DefaultRetryOptions()
    opts.MaxRetries = 5
    
    // Try to execute the unreliable function with retries
    result, err := RetryWithContext(ctx, unreliableFunction, opts)
    if err != nil {
        fmt.Printf("Final error: %v\n", err)
    } else {
        fmt.Printf("Success: %v\n", result)
    }
}

func main() {
    demonstrateRetryPattern()
}
```

### Context with Circuit Breaker

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
    mu                sync.Mutex
    name              string
    state             string
    failureThreshold  int
    resetTimeout      time.Duration
    failureCount      int
    lastFailure       time.Time
    halfOpenSuccess   int
    halfOpenSuccessThreshold int
}

// CircuitBreakerState represents the state of a circuit breaker
const (
    StateClosed    = "CLOSED"    // Normal operation, requests pass through
    StateOpen      = "OPEN"      // Circuit is open, requests fail fast
    StateHalfOpen  = "HALF-OPEN" // Testing if service is healthy again
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        name:             name,
        state:            StateClosed,
        failureThreshold: failureThreshold,
        resetTimeout:     resetTimeout,
        halfOpenSuccessThreshold: 2,
    }
}

// Execute runs a function with circuit breaker protection
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) (interface{}, error)) (interface{}, error) {
    cb.mu.Lock()
    state := cb.state
    
    // Check if circuit is open
    if state == StateOpen {
        // Check if we should try half-open state
        if time.Since(cb.lastFailure) > cb.resetTimeout {
            cb.state = StateHalfOpen
            cb.halfOpenSuccess = 0
            state = StateHalfOpen
            fmt.Printf("[%s] Circuit changed from OPEN to HALF-OPEN\n", cb.name)
        } else {
            cb.mu.Unlock()
            return nil, fmt.Errorf("circuit breaker '%s' is OPEN", cb.name)
        }
    }
    cb.mu.Unlock()
    
    // Execute the function
    result, err := fn(ctx)
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    // Handle the result based on current state
    if err != nil {
        // Handle failure
        switch state {
        case StateClosed:
            cb.failureCount++
            if cb.failureCount >= cb.failureThreshold {
                cb.state = StateOpen
                cb.lastFailure = time.Now()
                fmt.Printf("[%s] Circuit changed from CLOSED to OPEN\n", cb.name)
            }
        case StateHalfOpen:
            cb.state = StateOpen
            cb.lastFailure = time.Now()
            fmt.Printf("[%s] Circuit changed from HALF-OPEN to OPEN\n", cb.name)
        }
        return nil, err
    }
    
    // Handle success
    switch state {
    case StateClosed:
        cb.failureCount = 0
    case StateHalfOpen:
        cb.halfOpenSuccess++
        if cb.halfOpenSuccess >= cb.halfOpenSuccessThreshold {
            cb.state = StateClosed
            cb.failureCount = 0
            fmt.Printf("[%s] Circuit changed from HALF-OPEN to CLOSED\n", cb.name)
        }
    }
    
    return result, nil
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() string {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    return cb.state
}

// Simulate an unreliable service
func unreliableService(ctx context.Context, shouldFail bool) (interface{}, error) {
    // Check for context cancellation
    if ctx.Err() != nil {
        return nil, ctx.Err()
    }
    
    // Simulate work
    time.Sleep(100 * time.Millisecond)
    
    if shouldFail {
        return nil, fmt.Errorf("service error")
    }
    
    return "service result", nil
}

func demonstrateCircuitBreaker() {
    fmt.Println("=== Context with Circuit Breaker ===")
    
    // Create a circuit breaker
    cb := NewCircuitBreaker("service-cb", 3, 2*time.Second)
    
    // Create a base context
    baseCtx := context.Background()
    
    // Phase 1: Service is working
    fmt.Println("\nPhase 1: Service is working")
    for i := 1; i <= 3; i++ {
        ctx, cancel := context.WithTimeout(baseCtx, 1*time.Second)
        
        result, err := cb.Execute(ctx, func(ctx context.Context) (interface{}, error) {
            return unreliableService(ctx, false)
        })
        
        if err != nil {
            fmt.Printf("Request %d failed: %v (Circuit: %s)\n", i, err, cb.GetState())
        } else {
            fmt.Printf("Request %d succeeded: %v (Circuit: %s)\n", i, result, cb.GetState())
        }
        
        cancel()
        time.Sleep(200 * time.Millisecond)
    }
    
    // Phase 2: Service starts failing
    fmt.Println("\nPhase 2: Service starts failing")
    for i := 1; i <= 5; i++ {
        ctx, cancel := context.WithTimeout(baseCtx, 1*time.Second)
        
        result, err := cb.Execute(ctx, func(ctx context.Context) (interface{}, error) {
            return unreliableService(ctx, true)
        })
        
        if err != nil {
            fmt.Printf("Request %d failed: %v (Circuit: %s)\n", i, err, cb.GetState())
        } else {
            fmt.Printf("Request %d succeeded: %v (Circuit: %s)\n", i, result, cb.GetState())
        }
        
        cancel()
        time.Sleep(200 * time.Millisecond)
    }
    
    // Phase 3: Circuit is open, requests fail fast
    fmt.Println("\nPhase 3: Circuit is open, requests fail fast")
    for i := 1; i <= 3; i++ {
        ctx, cancel := context.WithTimeout(baseCtx, 1*time.Second)
        
        start := time.Now()
        result, err := cb.Execute(ctx, func(ctx context.Context) (interface{}, error) {
            return unreliableService(ctx, true)
        })
        duration := time.Since(start)
        
        if err != nil {
            fmt.Printf("Request %d failed: %v (Circuit: %s, Duration: %v)\n", 
                i, err, cb.GetState(), duration)
        } else {
            fmt.Printf("Request %d succeeded: %v (Circuit: %s, Duration: %v)\n", 
                i, result, cb.GetState(), duration)
        }
        
        cancel()
        time.Sleep(200 * time.Millisecond)
    }
    
    // Phase 4: Wait for reset timeout
    fmt.Println("\nPhase 4: Waiting for reset timeout...")
    time.Sleep(2 * time.Second)
    
    // Phase 5: Circuit goes to half-open and service recovers
    fmt.Println("\nPhase 5: Service recovers")
    for i := 1; i <= 3; i++ {
        ctx, cancel := context.WithTimeout(baseCtx, 1*time.Second)
        
        result, err := cb.Execute(ctx, func(ctx context.Context) (interface{}, error) {
            return unreliableService(ctx, false)
        })
        
        if err != nil {
            fmt.Printf("Request %d failed: %v (Circuit: %s)\n", i, err, cb.GetState())
        } else {
            fmt.Printf("Request %d succeeded: %v (Circuit: %s)\n", i, result, cb.GetState())
        }
        
        cancel()
        time.Sleep(200 * time.Millisecond)
    }
}

func main() {
    demonstrateCircuitBreaker()
}
```

## Exercises

### Exercise 1: Implement a Context-Aware Cache

Create a cache that respects context cancellation and timeouts.

```go
// Implement a cache with the following features:
// - Get/Set operations that respect context cancellation
// - Automatic expiration of entries
// - Background cleanup of expired entries
// - Ability to cancel all operations when a parent context is canceled

package main

import (
    "context"
    "sync"
    "time"
)

type CacheEntry struct {
    Value     interface{}
    ExpiresAt time.Time
}

type Cache struct {
    mu      sync.RWMutex
    data    map[string]CacheEntry
    cleanup context.CancelFunc
}

// Implement these methods:
// NewCache(ctx context.Context, cleanupInterval time.Duration) *Cache
// Get(ctx context.Context, key string) (interface{}, bool)
// Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
// Delete(ctx context.Context, key string) error
// Clear(ctx context.Context) error
// Close() error
```

### Exercise 2: Implement a Context-Aware Worker Pool

Create a worker pool that can be controlled and monitored using context.

```go
// Implement a worker pool with the following features:
// - Submit tasks that respect context cancellation
// - Control the number of workers
// - Graceful shutdown when context is canceled
// - Monitor worker status and task completion

package main

import (
    "context"
    "sync"
)

type Task func(ctx context.Context) (interface{}, error)

type Result struct {
    Value interface{}
    Err   error
    Task  Task
}

type WorkerPool struct {
    workers int
    tasks   chan Task
    results chan Result
    wg      sync.WaitGroup
    ctx     context.Context
    cancel  context.CancelFunc
}

// Implement these methods:
// NewWorkerPool(ctx context.Context, numWorkers int) *WorkerPool
// Start()
// Submit(ctx context.Context, task Task) error
// Results() <-chan Result
// Shutdown(ctx context.Context) error
```

### Exercise 3: Implement a Context-Aware Rate Limiter

Create a rate limiter that respects context cancellation and timeouts.

```go
// Implement a rate limiter with the following features:
// - Limit requests per second
// - Wait for permits with context cancellation
// - Support for bursts
// - Track statistics

package main

import (
    "context"
    "sync"
    "time"
)

type RateLimiter struct {
    rate      float64 // permits per second
    burst     int     // maximum burst size
    tokens    float64 // current token count
    lastRefill time.Time
    mu        sync.Mutex
}

// Implement these methods:
// NewRateLimiter(rate float64, burst int) *RateLimiter
// Allow() bool
// Wait(ctx context.Context) error
// Reserve() (Reservation, bool)
// SetRate(rate float64)
// SetBurst(burst int)
```

## Key Takeaways

1. **Context is for request-scoped data**: Use context to carry request-scoped values, deadlines, and cancellation signals across API boundaries.

2. **Context is immutable**: Never modify an existing context. Instead, create a new context derived from the original.

3. **Context should flow through your program**: Pass context as the first parameter to functions that may block or need to be canceled.

4. **Always check for context cancellation**: Regularly check `ctx.Done()` in long-running operations to allow early cancellation.

5. **Use context values sparingly**: Only use context values for request-scoped data that truly needs to transit process boundaries, not for passing optional parameters.

6. **Use custom types for context keys**: To avoid key collisions, use custom types for context keys, not strings or basic types.

7. **Always call cancel functions**: When creating a context with a cancel function, always call the cancel function, typically with `defer`.

8. **Context is not for dependency injection**: Don't use context to pass dependencies like database connections or configuration.

9. **Context should have a lifetime**: The lifetime of a context should match the lifetime of the request or operation it's associated with.

10. **Context is a Go-specific pattern**: The context pattern is specific to Go and may not translate directly to other languages.

## Next Steps

Now that you understand the context package, let's learn about [Testing in Go](19-testing.md) to ensure your code works correctly!

---

**Previous**: [← Synchronization Primitives](17-synchronization.md) | **Next**: [Testing in Go →](19-testing.md)