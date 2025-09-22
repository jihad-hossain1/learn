# 22. Context

The `context` package in Go provides a way to carry deadlines, cancellation signals, and request-scoped values across API boundaries and between processes. It's essential for building robust, cancellable operations in concurrent programs.

## Context Basics

### What is Context?

Context is an interface that carries:
- **Deadlines**: When an operation should be cancelled
- **Cancellation signals**: Manual cancellation triggers
- **Values**: Request-scoped data (use sparingly)

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func main() {
    // Background context - never cancelled, no deadline, no values
    ctx := context.Background()
    
    fmt.Printf("Context type: %T\n", ctx)
    fmt.Printf("Deadline: %v, %v\n", ctx.Deadline())
    fmt.Printf("Done channel: %v\n", ctx.Done())
    fmt.Printf("Error: %v\n", ctx.Err())
    
    // TODO context - placeholder for when you don't know what context to use
    todoCtx := context.TODO()
    fmt.Printf("\nTODO Context type: %T\n", todoCtx)
}
```

### Context Interface

```go
package main

import (
    "context"
    "fmt"
)

// The context.Context interface
type Context interface {
    // Deadline returns the time when work done on behalf of this context
    // should be canceled.
    Deadline() (deadline time.Time, ok bool)
    
    // Done returns a channel that's closed when work done on behalf of this
    // context should be canceled.
    Done() <-chan struct{}
    
    // Err returns a non-nil error value after Done is closed.
    Err() error
    
    // Value returns the value associated with this context for key.
    Value(key interface{}) interface{}
}

func demonstrateContextInterface() {
    ctx := context.Background()
    
    // Check if context has a deadline
    if deadline, ok := ctx.Deadline(); ok {
        fmt.Printf("Context has deadline: %v\n", deadline)
    } else {
        fmt.Println("Context has no deadline")
    }
    
    // Check if context is done
    select {
    case <-ctx.Done():
        fmt.Printf("Context is done: %v\n", ctx.Err())
    default:
        fmt.Println("Context is not done")
    }
}

func main() {
    demonstrateContextInterface()
}
```

## Context Creation Functions

### WithCancel

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: Cancelled (%v)\n", id, ctx.Err())
            return
        default:
            fmt.Printf("Worker %d: Working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // Create cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }
    
    // Let workers run for 2 seconds
    time.Sleep(2 * time.Second)
    
    // Cancel all workers
    fmt.Println("\nCancelling workers...")
    cancel()
    
    // Give time for cleanup
    time.Sleep(1 * time.Second)
    fmt.Println("Main function ending")
}
```

### WithTimeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func longRunningTask(ctx context.Context, taskName string) error {
    fmt.Printf("%s: Starting\n", taskName)
    
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("%s: Cancelled at step %d (%v)\n", taskName, i, ctx.Err())
            return ctx.Err()
        default:
            fmt.Printf("%s: Step %d\n", taskName, i+1)
            time.Sleep(300 * time.Millisecond)
        }
    }
    
    fmt.Printf("%s: Completed successfully\n", taskName)
    return nil
}

func main() {
    // Context with 2-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel() // Always call cancel to release resources
    
    // Start multiple tasks
    for i := 1; i <= 3; i++ {
        taskName := fmt.Sprintf("Task-%d", i)
        go func(name string) {
            err := longRunningTask(ctx, name)
            if err != nil {
                fmt.Printf("%s failed: %v\n", name, err)
            }
        }(taskName)
    }
    
    // Wait for context to expire
    <-ctx.Done()
    fmt.Printf("\nContext expired: %v\n", ctx.Err())
    
    // Give time for cleanup
    time.Sleep(500 * time.Millisecond)
}
```

### WithDeadline

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func processWithDeadline(ctx context.Context, data string) {
    deadline, ok := ctx.Deadline()
    if ok {
        fmt.Printf("Processing '%s' with deadline: %v\n", data, deadline.Format("15:04:05"))
        fmt.Printf("Time remaining: %v\n", time.Until(deadline))
    }
    
    // Simulate work
    for i := 0; i < 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("Processing cancelled: %v\n", ctx.Err())
            return
        default:
            fmt.Printf("Processing step %d\n", i+1)
            time.Sleep(400 * time.Millisecond)
        }
    }
    
    fmt.Printf("Processing '%s' completed\n", data)
}

func main() {
    // Set deadline to 1.5 seconds from now
    deadline := time.Now().Add(1500 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    fmt.Printf("Current time: %v\n", time.Now().Format("15:04:05"))
    fmt.Printf("Deadline set for: %v\n\n", deadline.Format("15:04:05"))
    
    go processWithDeadline(ctx, "important-data")
    
    // Wait for completion or deadline
    <-ctx.Done()
    
    fmt.Printf("\nMain function ending at: %v\n", time.Now().Format("15:04:05"))
    time.Sleep(100 * time.Millisecond)
}
```

### WithValue

```go
package main

import (
    "context"
    "fmt"
)

// Define custom key types to avoid collisions
type contextKey string

const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
    traceIDKey   contextKey = "traceID"
)

func processRequest(ctx context.Context) {
    // Extract values from context
    userID := ctx.Value(userIDKey)
    requestID := ctx.Value(requestIDKey)
    traceID := ctx.Value(traceIDKey)
    
    fmt.Printf("Processing request:\n")
    fmt.Printf("  User ID: %v\n", userID)
    fmt.Printf("  Request ID: %v\n", requestID)
    fmt.Printf("  Trace ID: %v\n", traceID)
    
    // Call another function that also needs these values
    authenticateUser(ctx)
    logActivity(ctx)
}

func authenticateUser(ctx context.Context) {
    userID := ctx.Value(userIDKey)
    if userID == nil {
        fmt.Println("Authentication: No user ID found")
        return
    }
    
    fmt.Printf("Authentication: Validating user %v\n", userID)
}

func logActivity(ctx context.Context) {
    userID := ctx.Value(userIDKey)
    requestID := ctx.Value(requestIDKey)
    traceID := ctx.Value(traceIDKey)
    
    fmt.Printf("Logging: [%v] User %v performed action in request %v\n", 
        traceID, userID, requestID)
}

func main() {
    // Create context with values
    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey, "user123")
    ctx = context.WithValue(ctx, requestIDKey, "req456")
    ctx = context.WithValue(ctx, traceIDKey, "trace789")
    
    processRequest(ctx)
    
    fmt.Println("\n--- Testing with missing values ---")
    
    // Create context with only some values
    ctx2 := context.Background()
    ctx2 = context.WithValue(ctx2, requestIDKey, "req999")
    
    processRequest(ctx2)
}
```

## Practical Context Patterns

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

type contextKey string

const requestIDKey contextKey = "requestID"

// Middleware to add request ID to context
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := fmt.Sprintf("%d", time.Now().UnixNano())
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)
        
        // Add request ID to response header
        w.Header().Set("X-Request-ID", requestID)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Middleware to add timeout to requests
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    requestID := ctx.Value(requestIDKey)
    
    log.Printf("[%v] Starting slow operation", requestID)
    
    // Simulate slow operation
    select {
    case <-time.After(3 * time.Second):
        fmt.Fprintf(w, "Operation completed for request %v", requestID)
        log.Printf("[%v] Operation completed", requestID)
    case <-ctx.Done():
        log.Printf("[%v] Operation cancelled: %v", requestID, ctx.Err())
        http.Error(w, "Request timeout", http.StatusRequestTimeout)
    }
}

func fastHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    requestID := ctx.Value(requestIDKey)
    
    log.Printf("[%v] Fast operation", requestID)
    fmt.Fprintf(w, "Fast response for request %v", requestID)
}

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/slow", slowHandler)
    mux.HandleFunc("/fast", fastHandler)
    
    // Apply middlewares
    handler := requestIDMiddleware(
        timeoutMiddleware(2 * time.Second)(mux),
    )
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: curl http://localhost:8080/fast")
    fmt.Println("Try: curl http://localhost:8080/slow (will timeout)")
    
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

### Database Operations with Context

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "time"
)

// Simulated database
type Database struct {
    name string
}

type User struct {
    ID   int
    Name string
    Email string
}

func (db *Database) GetUser(ctx context.Context, userID int) (*User, error) {
    // Simulate database query time
    queryTime := time.Duration(rand.Intn(2000)) * time.Millisecond
    
    fmt.Printf("DB: Starting query for user %d (estimated time: %v)\n", userID, queryTime)
    
    select {
    case <-time.After(queryTime):
        // Query completed
        user := &User{
            ID:    userID,
            Name:  fmt.Sprintf("User%d", userID),
            Email: fmt.Sprintf("user%d@example.com", userID),
        }
        fmt.Printf("DB: Query completed for user %d\n", userID)
        return user, nil
        
    case <-ctx.Done():
        // Query cancelled
        fmt.Printf("DB: Query cancelled for user %d: %v\n", userID, ctx.Err())
        return nil, ctx.Err()
    }
}

func (db *Database) UpdateUser(ctx context.Context, user *User) error {
    updateTime := time.Duration(rand.Intn(1500)) * time.Millisecond
    
    fmt.Printf("DB: Starting update for user %d (estimated time: %v)\n", user.ID, updateTime)
    
    select {
    case <-time.After(updateTime):
        fmt.Printf("DB: Update completed for user %d\n", user.ID)
        return nil
        
    case <-ctx.Done():
        fmt.Printf("DB: Update cancelled for user %d: %v\n", user.ID, ctx.Err())
        return ctx.Err()
    }
}

type UserService struct {
    db *Database
}

func (s *UserService) GetAndUpdateUser(ctx context.Context, userID int, newName string) error {
    // Get user with context
    user, err := s.db.GetUser(ctx, userID)
    if err != nil {
        return fmt.Errorf("failed to get user: %w", err)
    }
    
    // Update user data
    user.Name = newName
    
    // Update user with context
    err = s.db.UpdateUser(ctx, user)
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    fmt.Printf("Successfully updated user %d to name '%s'\n", userID, newName)
    return nil
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    db := &Database{name: "UserDB"}
    service := &UserService{db: db}
    
    // Test with timeout
    fmt.Println("=== Test with 3-second timeout ===")
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    err := service.GetAndUpdateUser(ctx, 123, "John Doe")
    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
    }
    
    time.Sleep(500 * time.Millisecond)
    
    // Test with short timeout
    fmt.Println("\n=== Test with 500ms timeout ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
    defer cancel2()
    
    err = service.GetAndUpdateUser(ctx2, 456, "Jane Smith")
    if err != nil {
        fmt.Printf("Operation failed: %v\n", err)
    }
    
    time.Sleep(1 * time.Second)
}
```

### Concurrent Operations with Context

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Task struct {
    ID   int
    Name string
}

type Result struct {
    TaskID int
    Value  string
    Error  error
}

func processTask(ctx context.Context, task Task) Result {
    processingTime := time.Duration(rand.Intn(2000)) * time.Millisecond
    
    fmt.Printf("Task %d (%s): Starting (estimated time: %v)\n", 
        task.ID, task.Name, processingTime)
    
    select {
    case <-time.After(processingTime):
        result := Result{
            TaskID: task.ID,
            Value:  fmt.Sprintf("Processed %s", task.Name),
        }
        fmt.Printf("Task %d: Completed successfully\n", task.ID)
        return result
        
    case <-ctx.Done():
        result := Result{
            TaskID: task.ID,
            Error:  ctx.Err(),
        }
        fmt.Printf("Task %d: Cancelled (%v)\n", task.ID, ctx.Err())
        return result
    }
}

func processConcurrentTasks(ctx context.Context, tasks []Task) []Result {
    var wg sync.WaitGroup
    results := make(chan Result, len(tasks))
    
    // Start all tasks concurrently
    for _, task := range tasks {
        wg.Add(1)
        go func(t Task) {
            defer wg.Done()
            result := processTask(ctx, t)
            results <- result
        }(task)
    }
    
    // Close results channel when all tasks are done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var allResults []Result
    for result := range results {
        allResults = append(allResults, result)
    }
    
    return allResults
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    tasks := []Task{
        {ID: 1, Name: "Process Data A"},
        {ID: 2, Name: "Process Data B"},
        {ID: 3, Name: "Process Data C"},
        {ID: 4, Name: "Process Data D"},
        {ID: 5, Name: "Process Data E"},
    }
    
    // Test with generous timeout
    fmt.Println("=== Test with 5-second timeout ===")
    ctx1, cancel1 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel1()
    
    results1 := processConcurrentTasks(ctx1, tasks)
    
    fmt.Println("\nResults:")
    for _, result := range results1 {
        if result.Error != nil {
            fmt.Printf("  Task %d: FAILED - %v\n", result.TaskID, result.Error)
        } else {
            fmt.Printf("  Task %d: SUCCESS - %s\n", result.TaskID, result.Value)
        }
    }
    
    time.Sleep(1 * time.Second)
    
    // Test with short timeout
    fmt.Println("\n=== Test with 1-second timeout ===")
    ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel2()
    
    results2 := processConcurrentTasks(ctx2, tasks)
    
    fmt.Println("\nResults:")
    for _, result := range results2 {
        if result.Error != nil {
            fmt.Printf("  Task %d: FAILED - %v\n", result.TaskID, result.Error)
        } else {
            fmt.Printf("  Task %d: SUCCESS - %s\n", result.TaskID, result.Value)
        }
    }
}
```

## Context Propagation

### Context Chain

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type contextKey string

const (
    userIDKey   contextKey = "userID"
    sessionKey  contextKey = "sessionID"
    operationKey contextKey = "operation"
)

func serviceA(ctx context.Context) error {
    // Add operation info to context
    ctx = context.WithValue(ctx, operationKey, "ServiceA")
    
    fmt.Printf("ServiceA: User=%v, Session=%v, Operation=%v\n",
        ctx.Value(userIDKey), ctx.Value(sessionKey), ctx.Value(operationKey))
    
    // Create timeout for this service
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    return serviceB(ctx)
}

func serviceB(ctx context.Context) error {
    // Update operation info
    ctx = context.WithValue(ctx, operationKey, "ServiceB")
    
    fmt.Printf("ServiceB: User=%v, Session=%v, Operation=%v\n",
        ctx.Value(userIDKey), ctx.Value(sessionKey), ctx.Value(operationKey))
    
    // Add cancellation for this service
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    return serviceC(ctx)
}

func serviceC(ctx context.Context) error {
    // Update operation info
    ctx = context.WithValue(ctx, operationKey, "ServiceC")
    
    fmt.Printf("ServiceC: User=%v, Session=%v, Operation=%v\n",
        ctx.Value(userIDKey), ctx.Value(sessionKey), ctx.Value(operationKey))
    
    // Simulate work
    select {
    case <-time.After(1 * time.Second):
        fmt.Println("ServiceC: Work completed")
        return nil
    case <-ctx.Done():
        fmt.Printf("ServiceC: Cancelled - %v\n", ctx.Err())
        return ctx.Err()
    }
}

func main() {
    // Create root context with user info
    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey, "user123")
    ctx = context.WithValue(ctx, sessionKey, "session456")
    
    fmt.Println("Starting service chain...")
    
    err := serviceA(ctx)
    if err != nil {
        fmt.Printf("Service chain failed: %v\n", err)
    } else {
        fmt.Println("Service chain completed successfully")
    }
}
```

### Context Best Practices

```go
package main

import (
    "context"
    "fmt"
    "time"
)

// Good: Use custom types for context keys
type contextKey string

const userKey contextKey = "user"

type User struct {
    ID   string
    Name string
}

// Good: Helper functions for context values
func WithUser(ctx context.Context, user *User) context.Context {
    return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(userKey).(*User)
    return user, ok
}

// Good: Always accept context as first parameter
func processUserData(ctx context.Context, data string) error {
    user, ok := UserFromContext(ctx)
    if !ok {
        return fmt.Errorf("no user in context")
    }
    
    fmt.Printf("Processing data for user %s: %s\n", user.Name, data)
    
    // Check for cancellation during work
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(500 * time.Millisecond):
        // Work completed
    }
    
    return nil
}

// Good: Return context from functions that modify it
func withTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
    return context.WithTimeout(ctx, timeout)
}

// Bad: Don't store context in structs (with exceptions)
type BadService struct {
    ctx context.Context // Generally avoid this
}

// Good: Pass context to methods
type GoodService struct {
    name string
}

func (s *GoodService) ProcessData(ctx context.Context, data string) error {
    fmt.Printf("Service %s processing data\n", s.name)
    return processUserData(ctx, data)
}

func main() {
    // Create context with user
    user := &User{ID: "123", Name: "John Doe"}
    ctx := WithUser(context.Background(), user)
    
    // Add timeout
    ctx, cancel := withTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // Use service
    service := &GoodService{name: "DataProcessor"}
    
    err := service.ProcessData(ctx, "important data")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Common Context Patterns

### Graceful Shutdown

```go
package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"
)

type Worker struct {
    id   int
    name string
}

func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %s starting\n", w.name)
    
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %s shutting down: %v\n", w.name, ctx.Err())
            return
        case <-ticker.C:
            fmt.Printf("Worker %s: tick\n", w.name)
        }
    }
}

type Server struct {
    workers []Worker
}

func (s *Server) Start(ctx context.Context) {
    var wg sync.WaitGroup
    
    // Start all workers
    for _, worker := range s.workers {
        wg.Add(1)
        go worker.Start(ctx, &wg)
    }
    
    // Wait for context cancellation
    <-ctx.Done()
    fmt.Println("Server: Shutdown signal received")
    
    // Wait for all workers to finish
    done := make(chan struct{})
    go func() {
        wg.Wait()
        close(done)
    }()
    
    // Wait for graceful shutdown or timeout
    select {
    case <-done:
        fmt.Println("Server: All workers stopped gracefully")
    case <-time.After(5 * time.Second):
        fmt.Println("Server: Shutdown timeout, forcing exit")
    }
}

func main() {
    // Create server with workers
    server := &Server{
        workers: []Worker{
            {id: 1, name: "worker-1"},
            {id: 2, name: "worker-2"},
            {id: 3, name: "worker-3"},
        },
    }
    
    // Create context that cancels on interrupt signal
    ctx, cancel := context.WithCancel(context.Background())
    
    // Handle interrupt signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        sig := <-sigChan
        fmt.Printf("\nReceived signal: %v\n", sig)
        cancel()
    }()
    
    fmt.Println("Server starting... Press Ctrl+C to stop")
    server.Start(ctx)
    fmt.Println("Server stopped")
}
```

## Context Anti-Patterns

### What NOT to Do

```go
package main

import (
    "context"
    "fmt"
)

// BAD: Don't use context for optional parameters
func badFunction(ctx context.Context) {
    // Don't do this
    if timeout := ctx.Value("timeout"); timeout != nil {
        // Use timeout
    }
}

// GOOD: Use explicit parameters
func goodFunction(ctx context.Context, timeout int) {
    // Explicit parameter is better
}

// BAD: Don't store context in structs (usually)
type BadStruct struct {
    ctx context.Context
}

// GOOD: Pass context to methods
type GoodStruct struct {
    name string
}

func (g *GoodStruct) DoWork(ctx context.Context) {
    // Pass context as parameter
}

// BAD: Don't use nil context
func badNilContext() {
    // Don't do this
    // someFunction(nil)
}

// GOOD: Use context.Background() or context.TODO()
func goodContext() {
    ctx := context.Background()
    // or
    ctx = context.TODO()
    _ = ctx
}

// BAD: Don't ignore context cancellation
func badIgnoreContext(ctx context.Context) {
    // This ignores cancellation
    for i := 0; i < 1000000; i++ {
        // Do work without checking ctx.Done()
        fmt.Printf("Working... %d\n", i)
    }
}

// GOOD: Check context regularly
func goodCheckContext(ctx context.Context) {
    for i := 0; i < 1000000; i++ {
        select {
        case <-ctx.Done():
            fmt.Println("Work cancelled")
            return
        default:
            // Do a small amount of work
            if i%1000 == 0 {
                fmt.Printf("Working... %d\n", i)
            }
        }
    }
}

func main() {
    fmt.Println("Context anti-patterns demonstrated")
}
```

## Best Practices Summary

### Do's

1. **Always pass context as the first parameter**
2. **Use context.Background() for root contexts**
3. **Use context.TODO() when you're unsure**
4. **Always call cancel() to release resources**
5. **Check ctx.Done() in long-running operations**
6. **Use custom types for context keys**
7. **Create helper functions for context values**

### Don'ts

1. **Don't store context in structs (usually)**
2. **Don't pass nil context**
3. **Don't use context for optional parameters**
4. **Don't ignore context cancellation**
5. **Don't use string keys for context values**

## Exercises

### Exercise 1: HTTP Client with Context
Create an HTTP client that:
- Supports request timeouts via context
- Can be cancelled mid-request
- Propagates request IDs through context
- Handles retries with exponential backoff

### Exercise 2: Worker Pool with Graceful Shutdown
Implement a worker pool that:
- Processes jobs from a queue
- Supports graceful shutdown via context
- Handles in-flight jobs during shutdown
- Reports progress and statistics

### Exercise 3: Distributed Tracing
Build a simple distributed tracing system that:
- Propagates trace IDs through context
- Measures operation durations
- Supports nested spans
- Handles context cancellation properly

## Key Takeaways

- Context provides cancellation, timeouts, and request-scoped values
- Always pass context as the first parameter to functions
- Use context.Background() for root contexts
- Always call cancel() to release resources
- Check ctx.Done() in long-running operations
- Use custom types for context keys to avoid collisions
- Don't store context in structs (with rare exceptions)
- Context is essential for building robust, cancellable operations

## Next Steps

Next, we'll explore [Testing](23-testing.md) to learn about Go's built-in testing framework and best practices for writing reliable tests.