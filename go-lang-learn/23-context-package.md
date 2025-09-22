# 23. Context Package

The `context` package is fundamental to Go's approach to handling cancellation, timeouts, and request-scoped values across API boundaries. It provides a standardized way to carry deadlines, cancellation signals, and other request-scoped values across goroutines and API calls.

## Introduction to Context

Context is an interface that carries deadlines, cancellation signals, and other request-scoped values across API boundaries and between processes.

### The Context Interface

```go
// context_interface.go
package main

import (
    "context"
    "fmt"
    "time"
)

// The Context interface (from standard library)
// type Context interface {
//     Deadline() (deadline time.Time, ok bool)
//     Done() <-chan struct{}
//     Err() error
//     Value(key interface{}) interface{}
// }

func exploreContextInterface() {
    // Create a context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // Check deadline
    if deadline, ok := ctx.Deadline(); ok {
        fmt.Printf("Context deadline: %v\n", deadline)
        fmt.Printf("Time until deadline: %v\n", time.Until(deadline))
    }
    
    // Check if context is done
    select {
    case <-ctx.Done():
        fmt.Printf("Context is done: %v\n", ctx.Err())
    default:
        fmt.Println("Context is not done yet")
    }
    
    // Wait for context to be done
    <-ctx.Done()
    fmt.Printf("Context finished with error: %v\n", ctx.Err())
}

func main() {
    fmt.Println("=== Context Interface Exploration ===")
    exploreContextInterface()
}
```

## Creating Contexts

### Background and TODO Contexts

```go
// context_creation.go
package main

import (
    "context"
    "fmt"
    "time"
)

func backgroundContext() {
    // Background context - typically used at the top level
    ctx := context.Background()
    
    fmt.Printf("Background context: %v\n", ctx)
    fmt.Printf("Background deadline: %v\n", ctx.Deadline())
    fmt.Printf("Background done channel: %v\n", ctx.Done())
    fmt.Printf("Background error: %v\n", ctx.Err())
}

func todoContext() {
    // TODO context - used when you're not sure which context to use
    ctx := context.TODO()
    
    fmt.Printf("TODO context: %v\n", ctx)
    fmt.Printf("TODO deadline: %v\n", ctx.Deadline())
    fmt.Printf("TODO done channel: %v\n", ctx.Done())
    fmt.Printf("TODO error: %v\n", ctx.Err())
}

func main() {
    fmt.Println("=== Background Context ===")
    backgroundContext()
    
    fmt.Println("\n=== TODO Context ===")
    todoContext()
}
```

### WithCancel

```go
// with_cancel.go
package main

import (
    "context"
    "fmt"
    "time"
)

func withCancelExample() {
    // Create a cancellable context
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start a goroutine that will be cancelled
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Goroutine cancelled:", ctx.Err())
                return
            default:
                fmt.Println("Goroutine working...")
                time.Sleep(500 * time.Millisecond)
            }
        }
    }()
    
    // Let it run for a while
    time.Sleep(2 * time.Second)
    
    // Cancel the context
    fmt.Println("Cancelling context...")
    cancel()
    
    // Give goroutine time to handle cancellation
    time.Sleep(1 * time.Second)
}

// Multiple goroutines with shared cancellation
func multipleGoroutinesCancel() {
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start multiple workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }
    
    // Let workers run
    time.Sleep(3 * time.Second)
    
    // Cancel all workers
    fmt.Println("Cancelling all workers...")
    cancel()
    
    // Wait for cleanup
    time.Sleep(1 * time.Second)
}

func worker(ctx context.Context, id int) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d stopped: %v\n", id, ctx.Err())
            return
        default:
            fmt.Printf("Worker %d is working\n", id)
            time.Sleep(800 * time.Millisecond)
        }
    }
}

func main() {
    fmt.Println("=== WithCancel Example ===")
    withCancelExample()
    
    fmt.Println("\n=== Multiple Goroutines Cancel ===")
    multipleGoroutinesCancel()
}
```

### WithTimeout and WithDeadline

```go
// with_timeout.go
package main

import (
    "context"
    "fmt"
    "time"
)

func withTimeoutExample() {
    // Create context with 2-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel() // Always call cancel to release resources
    
    // Simulate work that might take too long
    go func() {
        select {
        case <-time.After(3 * time.Second):
            fmt.Println("Work completed")
        case <-ctx.Done():
            fmt.Println("Work cancelled due to timeout:", ctx.Err())
        }
    }()
    
    // Wait for context to be done
    <-ctx.Done()
    fmt.Println("Main function finished")
}

func withDeadlineExample() {
    // Create context with specific deadline
    deadline := time.Now().Add(1500 * time.Millisecond)
    ctx, cancel := context.WithDeadline(context.Background(), deadline)
    defer cancel()
    
    fmt.Printf("Deadline set for: %v\n", deadline)
    
    // Check remaining time periodically
    ticker := time.NewTicker(300 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            if deadline, ok := ctx.Deadline(); ok {
                remaining := time.Until(deadline)
                fmt.Printf("Time remaining: %v\n", remaining)
            }
        case <-ctx.Done():
            fmt.Println("Context deadline exceeded:", ctx.Err())
            return
        }
    }
}

// Timeout with early completion
func timeoutWithEarlyCompletion() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    // Simulate work that completes early
    done := make(chan bool)
    
    go func() {
        // Simulate work
        time.Sleep(1 * time.Second)
        fmt.Println("Work completed successfully")
        done <- true
    }()
    
    select {
    case <-done:
        fmt.Println("Work finished before timeout")
    case <-ctx.Done():
        fmt.Println("Work cancelled due to timeout:", ctx.Err())
    }
}

func main() {
    fmt.Println("=== WithTimeout Example ===")
    withTimeoutExample()
    
    fmt.Println("\n=== WithDeadline Example ===")
    withDeadlineExample()
    
    fmt.Println("\n=== Timeout with Early Completion ===")
    timeoutWithEarlyCompletion()
}
```

### WithValue

```go
// with_value.go
package main

import (
    "context"
    "fmt"
)

// Define custom key types to avoid collisions
type userIDKey struct{}
type requestIDKey struct{}
type authTokenKey struct{}

func withValueExample() {
    // Create context with values
    ctx := context.Background()
    
    // Add user ID
    ctx = context.WithValue(ctx, userIDKey{}, "user123")
    
    // Add request ID
    ctx = context.WithValue(ctx, requestIDKey{}, "req-456")
    
    // Add auth token
    ctx = context.WithValue(ctx, authTokenKey{}, "token-789")
    
    // Pass context to functions
    processRequest(ctx)
}

func processRequest(ctx context.Context) {
    // Extract values from context
    userID := getUserID(ctx)
    requestID := getRequestID(ctx)
    authToken := getAuthToken(ctx)
    
    fmt.Printf("Processing request:\n")
    fmt.Printf("  User ID: %s\n", userID)
    fmt.Printf("  Request ID: %s\n", requestID)
    fmt.Printf("  Auth Token: %s\n", authToken)
    
    // Pass context to other functions
    authenticateUser(ctx)
    logRequest(ctx)
}

func getUserID(ctx context.Context) string {
    if userID, ok := ctx.Value(userIDKey{}).(string); ok {
        return userID
    }
    return "unknown"
}

func getRequestID(ctx context.Context) string {
    if requestID, ok := ctx.Value(requestIDKey{}).(string); ok {
        return requestID
    }
    return "unknown"
}

func getAuthToken(ctx context.Context) string {
    if token, ok := ctx.Value(authTokenKey{}).(string); ok {
        return token
    }
    return "unknown"
}

func authenticateUser(ctx context.Context) {
    userID := getUserID(ctx)
    authToken := getAuthToken(ctx)
    
    fmt.Printf("Authenticating user %s with token %s\n", userID, authToken)
}

func logRequest(ctx context.Context) {
    requestID := getRequestID(ctx)
    userID := getUserID(ctx)
    
    fmt.Printf("Logging request %s for user %s\n", requestID, userID)
}

// Helper functions for context values
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, userIDKey{}, userID)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, requestIDKey{}, requestID)
}

func WithAuthToken(ctx context.Context, token string) context.Context {
    return context.WithValue(ctx, authTokenKey{}, token)
}

func helperFunctionsExample() {
    ctx := context.Background()
    
    // Use helper functions
    ctx = WithUserID(ctx, "user456")
    ctx = WithRequestID(ctx, "req-789")
    ctx = WithAuthToken(ctx, "token-abc")
    
    processRequest(ctx)
}

func main() {
    fmt.Println("=== WithValue Example ===")
    withValueExample()
    
    fmt.Println("\n=== Helper Functions Example ===")
    helperFunctionsExample()
}
```

## Practical Context Patterns

### HTTP Server with Context

```go
// http_server_context.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type UserService struct {
    users map[int]User
}

func NewUserService() *UserService {
    return &UserService{
        users: map[int]User{
            1: {ID: 1, Name: "Alice"},
            2: {ID: 2, Name: "Bob"},
            3: {ID: 3, Name: "Charlie"},
        },
    }
}

func (s *UserService) GetUser(ctx context.Context, id int) (*User, error) {
    // Simulate database operation with context
    select {
    case <-time.After(500 * time.Millisecond): // Simulate DB delay
        if user, exists := s.users[id]; exists {
            return &user, nil
        }
        return nil, fmt.Errorf("user %d not found", id)
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (s *UserService) CreateUser(ctx context.Context, name string) (*User, error) {
    // Simulate database operation
    select {
    case <-time.After(1 * time.Second): // Simulate DB delay
        id := len(s.users) + 1
        user := User{ID: id, Name: name}
        s.users[id] = user
        return &user, nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// Middleware to add request timeout
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            
            // Replace request context with timeout context
            r = r.WithContext(ctx)
            next.ServeHTTP(w, r)
        })
    }
}

// Middleware to add request ID
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())
        ctx := context.WithValue(r.Context(), requestIDKey{}, requestID)
        
        // Add request ID to response header
        w.Header().Set("X-Request-ID", requestID)
        
        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}

func (s *UserService) handleGetUser(w http.ResponseWriter, r *http.Request) {
    // Extract request ID from context
    requestID := r.Context().Value(requestIDKey{})
    log.Printf("[%v] Handling GET user request", requestID)
    
    // Parse user ID
    idStr := r.URL.Query().Get("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    // Get user with context
    user, err := s.GetUser(r.Context(), id)
    if err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusRequestTimeout)
            return
        }
        if err == context.Canceled {
            http.Error(w, "Request cancelled", http.StatusRequestTimeout)
            return
        }
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
    
    log.Printf("[%v] Successfully returned user %d", requestID, user.ID)
}

func (s *UserService) handleCreateUser(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value(requestIDKey{})
    log.Printf("[%v] Handling POST user request", requestID)
    
    var req struct {
        Name string `json:"name"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user, err := s.CreateUser(r.Context(), req.Name)
    if err != nil {
        if err == context.DeadlineExceeded {
            http.Error(w, "Request timeout", http.StatusRequestTimeout)
            return
        }
        if err == context.Canceled {
            http.Error(w, "Request cancelled", http.StatusRequestTimeout)
            return
        }
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
    
    log.Printf("[%v] Successfully created user %d", requestID, user.ID)
}

func startHTTPServer() {
    userService := NewUserService()
    
    mux := http.NewServeMux()
    mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            userService.handleGetUser(w, r)
        case http.MethodPost:
            userService.handleCreateUser(w, r)
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    })
    
    // Apply middleware
    handler := timeoutMiddleware(2 * time.Second)(mux)
    handler = requestIDMiddleware(handler)
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: curl 'http://localhost:8080/users?id=1'")
    fmt.Println("Try: curl -X POST -H 'Content-Type: application/json' -d '{\"name\":\"David\"}' http://localhost:8080/users")
    
    log.Fatal(http.ListenAndServe(":8080", handler))
}

func main() {
    startHTTPServer()
}
```

### Database Operations with Context

```go
// database_context.go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"
    
    _ "github.com/mattn/go-sqlite3"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}

type ProductRepository struct {
    db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
    return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProduct(ctx context.Context, id int) (*Product, error) {
    query := "SELECT id, name, price FROM products WHERE id = ?"
    
    // Use QueryRowContext to respect context cancellation
    row := r.db.QueryRowContext(ctx, query, id)
    
    var product Product
    err := row.Scan(&product.ID, &product.Name, &product.Price)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("product %d not found", id)
        }
        return nil, fmt.Errorf("failed to get product: %w", err)
    }
    
    return &product, nil
}

func (r *ProductRepository) GetProducts(ctx context.Context, limit int) ([]Product, error) {
    query := "SELECT id, name, price FROM products LIMIT ?"
    
    // Use QueryContext to respect context cancellation
    rows, err := r.db.QueryContext(ctx, query, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to query products: %w", err)
    }
    defer rows.Close()
    
    var products []Product
    for rows.Next() {
        // Check if context is cancelled during iteration
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
        
        var product Product
        if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
            return nil, fmt.Errorf("failed to scan product: %w", err)
        }
        products = append(products, product)
    }
    
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error during rows iteration: %w", err)
    }
    
    return products, nil
}

func (r *ProductRepository) CreateProduct(ctx context.Context, name string, price float64) (*Product, error) {
    query := "INSERT INTO products (name, price) VALUES (?, ?)"
    
    // Use ExecContext to respect context cancellation
    result, err := r.db.ExecContext(ctx, query, name, price)
    if err != nil {
        return nil, fmt.Errorf("failed to create product: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("failed to get last insert id: %w", err)
    }
    
    return &Product{
        ID:    int(id),
        Name:  name,
        Price: price,
    }, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, id int, name string, price float64) error {
    query := "UPDATE products SET name = ?, price = ? WHERE id = ?"
    
    result, err := r.db.ExecContext(ctx, query, name, price, id)
    if err != nil {
        return fmt.Errorf("failed to update product: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product %d not found", id)
    }
    
    return nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, id int) error {
    query := "DELETE FROM products WHERE id = ?"
    
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete product: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("product %d not found", id)
    }
    
    return nil
}

// Transaction with context
func (r *ProductRepository) TransferProducts(ctx context.Context, fromID, toID int, quantity int) error {
    // Begin transaction with context
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback() // Will be ignored if tx.Commit() succeeds
    
    // Check source product
    var fromQuantity int
    err = tx.QueryRowContext(ctx, "SELECT quantity FROM inventory WHERE product_id = ?", fromID).Scan(&fromQuantity)
    if err != nil {
        return fmt.Errorf("failed to get source quantity: %w", err)
    }
    
    if fromQuantity < quantity {
        return fmt.Errorf("insufficient quantity in source product")
    }
    
    // Update source
    _, err = tx.ExecContext(ctx, "UPDATE inventory SET quantity = quantity - ? WHERE product_id = ?", quantity, fromID)
    if err != nil {
        return fmt.Errorf("failed to update source: %w", err)
    }
    
    // Update destination
    _, err = tx.ExecContext(ctx, "UPDATE inventory SET quantity = quantity + ? WHERE product_id = ?", quantity, toID)
    if err != nil {
        return fmt.Errorf("failed to update destination: %w", err)
    }
    
    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    return nil
}

func setupDatabase() (*sql.DB, error) {
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        return nil, err
    }
    
    // Create tables
    _, err = db.Exec(`
        CREATE TABLE products (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            price REAL NOT NULL
        );
        
        CREATE TABLE inventory (
            product_id INTEGER PRIMARY KEY,
            quantity INTEGER NOT NULL,
            FOREIGN KEY (product_id) REFERENCES products(id)
        );
    `)
    if err != nil {
        return nil, err
    }
    
    // Insert sample data
    _, err = db.Exec(`
        INSERT INTO products (name, price) VALUES 
        ('Laptop', 999.99),
        ('Mouse', 29.99),
        ('Keyboard', 79.99);
        
        INSERT INTO inventory (product_id, quantity) VALUES
        (1, 10),
        (2, 50),
        (3, 25);
    `)
    if err != nil {
        return nil, err
    }
    
    return db, nil
}

func demonstrateDatabaseContext() {
    db, err := setupDatabase()
    if err != nil {
        log.Fatal("Failed to setup database:", err)
    }
    defer db.Close()
    
    repo := NewProductRepository(db)
    
    // Test with timeout context
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    // Get all products
    fmt.Println("=== Getting Products ===")
    products, err := repo.GetProducts(ctx, 10)
    if err != nil {
        log.Printf("Failed to get products: %v", err)
        return
    }
    
    for _, product := range products {
        fmt.Printf("Product: %+v\n", product)
    }
    
    // Get specific product
    fmt.Println("\n=== Getting Specific Product ===")
    product, err := repo.GetProduct(ctx, 1)
    if err != nil {
        log.Printf("Failed to get product: %v", err)
        return
    }
    fmt.Printf("Product 1: %+v\n", product)
    
    // Create new product
    fmt.Println("\n=== Creating Product ===")
    newProduct, err := repo.CreateProduct(ctx, "Monitor", 299.99)
    if err != nil {
        log.Printf("Failed to create product: %v", err)
        return
    }
    fmt.Printf("Created product: %+v\n", newProduct)
    
    // Test with cancelled context
    fmt.Println("\n=== Testing Cancelled Context ===")
    cancelledCtx, cancel := context.WithCancel(context.Background())
    cancel() // Cancel immediately
    
    _, err = repo.GetProduct(cancelledCtx, 1)
    if err != nil {
        fmt.Printf("Expected error with cancelled context: %v\n", err)
    }
}

func main() {
    demonstrateDatabaseContext()
}
```

### Concurrent Operations with Context

```go
// concurrent_context.go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Worker pool with context
type WorkerPool struct {
    workerCount int
    jobs        chan Job
    results     chan Result
    wg          sync.WaitGroup
}

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID  int
    Output string
    Error  error
}

func NewWorkerPool(workerCount int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobs:        make(chan Job, workerCount*2),
        results:     make(chan Result, workerCount*2),
    }
}

func (wp *WorkerPool) Start(ctx context.Context) {
    // Start workers
    for i := 0; i < wp.workerCount; i++ {
        wp.wg.Add(1)
        go wp.worker(ctx, i+1)
    }
    
    // Start result collector
    go wp.collectResults(ctx)
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
    defer wp.wg.Done()
    
    fmt.Printf("Worker %d started\n", id)
    
    for {
        select {
        case job, ok := <-wp.jobs:
            if !ok {
                fmt.Printf("Worker %d: jobs channel closed\n", id)
                return
            }
            
            // Process job with context
            result := wp.processJob(ctx, job, id)
            
            select {
            case wp.results <- result:
            case <-ctx.Done():
                fmt.Printf("Worker %d cancelled while sending result\n", id)
                return
            }
            
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled\n", id)
            return
        }
    }
}

func (wp *WorkerPool) processJob(ctx context.Context, job Job, workerID int) Result {
    // Simulate work with random duration
    workDuration := time.Duration(rand.Intn(2000)) * time.Millisecond
    
    select {
    case <-time.After(workDuration):
        output := fmt.Sprintf("Worker %d processed job %d: %s", workerID, job.ID, job.Data)
        return Result{
            JobID:  job.ID,
            Output: output,
            Error:  nil,
        }
    case <-ctx.Done():
        return Result{
            JobID:  job.ID,
            Output: "",
            Error:  fmt.Errorf("job %d cancelled: %w", job.ID, ctx.Err()),
        }
    }
}

func (wp *WorkerPool) collectResults(ctx context.Context) {
    for {
        select {
        case result, ok := <-wp.results:
            if !ok {
                fmt.Println("Results channel closed")
                return
            }
            
            if result.Error != nil {
                fmt.Printf("Job %d failed: %v\n", result.JobID, result.Error)
            } else {
                fmt.Printf("Job %d completed: %s\n", result.JobID, result.Output)
            }
            
        case <-ctx.Done():
            fmt.Println("Result collector cancelled")
            return
        }
    }
}

func (wp *WorkerPool) AddJob(job Job) {
    wp.jobs <- job
}

func (wp *WorkerPool) Close() {
    close(wp.jobs)
    wp.wg.Wait()
    close(wp.results)
}

// Fan-out pattern with context
func fanOutPattern(ctx context.Context, input <-chan int) (<-chan int, <-chan int, <-chan int) {
    out1 := make(chan int)
    out2 := make(chan int)
    out3 := make(chan int)
    
    go func() {
        defer close(out1)
        defer close(out2)
        defer close(out3)
        
        for {
            select {
            case value, ok := <-input:
                if !ok {
                    return
                }
                
                // Send to all outputs
                select {
                case out1 <- value:
                case <-ctx.Done():
                    return
                }
                
                select {
                case out2 <- value * 2:
                case <-ctx.Done():
                    return
                }
                
                select {
                case out3 <- value * 3:
                case <-ctx.Done():
                    return
                }
                
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return out1, out2, out3
}

// Pipeline pattern with context
func pipeline(ctx context.Context, input <-chan int) <-chan string {
    // Stage 1: Square numbers
    squared := make(chan int)
    go func() {
        defer close(squared)
        for {
            select {
            case num, ok := <-input:
                if !ok {
                    return
                }
                select {
                case squared <- num * num:
                case <-ctx.Done():
                    return
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    
    // Stage 2: Convert to string
    output := make(chan string)
    go func() {
        defer close(output)
        for {
            select {
            case num, ok := <-squared:
                if !ok {
                    return
                }
                result := fmt.Sprintf("squared: %d", num)
                select {
                case output <- result:
                case <-ctx.Done():
                    return
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    
    return output
}

func demonstrateWorkerPool() {
    fmt.Println("=== Worker Pool with Context ===")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    wp := NewWorkerPool(3)
    wp.Start(ctx)
    
    // Add jobs
    for i := 1; i <= 10; i++ {
        wp.AddJob(Job{
            ID:   i,
            Data: fmt.Sprintf("task-%d", i),
        })
    }
    
    // Let workers process for a while
    time.Sleep(3 * time.Second)
    
    // Close worker pool
    wp.Close()
    
    fmt.Println("Worker pool demonstration completed")
}

func demonstrateFanOut() {
    fmt.Println("\n=== Fan-out Pattern ===")
    
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    // Create input channel
    input := make(chan int)
    
    // Start fan-out
    out1, out2, out3 := fanOutPattern(ctx, input)
    
    // Start consumers
    var wg sync.WaitGroup
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for value := range out1 {
            fmt.Printf("Consumer 1 received: %d\n", value)
        }
    }()
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for value := range out2 {
            fmt.Printf("Consumer 2 received: %d\n", value)
        }
    }()
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        for value := range out3 {
            fmt.Printf("Consumer 3 received: %d\n", value)
        }
    }()
    
    // Send data
    go func() {
        defer close(input)
        for i := 1; i <= 5; i++ {
            select {
            case input <- i:
                time.Sleep(200 * time.Millisecond)
            case <-ctx.Done():
                return
            }
        }
    }()
    
    wg.Wait()
    fmt.Println("Fan-out demonstration completed")
}

func demonstratePipeline() {
    fmt.Println("\n=== Pipeline Pattern ===")
    
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    // Create input channel
    input := make(chan int)
    
    // Start pipeline
    output := pipeline(ctx, input)
    
    // Start consumer
    go func() {
        for result := range output {
            fmt.Printf("Pipeline result: %s\n", result)
        }
    }()
    
    // Send data
    go func() {
        defer close(input)
        for i := 1; i <= 5; i++ {
            select {
            case input <- i:
                time.Sleep(300 * time.Millisecond)
            case <-ctx.Done():
                return
            }
        }
    }()
    
    // Wait for context to finish
    <-ctx.Done()
    fmt.Println("Pipeline demonstration completed")
}

func main() {
    demonstrateWorkerPool()
    demonstrateFanOut()
    demonstratePipeline()
}
```

## Context Best Practices

### 1. Context Propagation

```go
// context_propagation.go
package main

import (
    "context"
    "fmt"
    "time"
)

// Always pass context as the first parameter
func goodFunction(ctx context.Context, userID string) error {
    // Use context in function
    select {
    case <-time.After(1 * time.Second):
        fmt.Printf("Processed user %s\n", userID)
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Don't store context in structs (anti-pattern)
type BadService struct {
    ctx context.Context // Don't do this
}

// Good: Pass context to methods
type GoodService struct {
    name string
}

func (s *GoodService) ProcessUser(ctx context.Context, userID string) error {
    return goodFunction(ctx, userID)
}

// Chain contexts properly
func chainContexts() {
    // Root context
    rootCtx := context.Background()
    
    // Add timeout
    timeoutCtx, cancel1 := context.WithTimeout(rootCtx, 5*time.Second)
    defer cancel1()
    
    // Add cancellation
    cancelCtx, cancel2 := context.WithCancel(timeoutCtx)
    defer cancel2()
    
    // Add values
    valueCtx := context.WithValue(cancelCtx, userIDKey{}, "user123")
    
    // Use the final context
    service := &GoodService{name: "UserService"}
    if err := service.ProcessUser(valueCtx, "user123"); err != nil {
        fmt.Printf("Error: %v\n", err)
    }
}

func main() {
    fmt.Println("=== Context Propagation ===")
    chainContexts()
}
```

### 2. Context Values Guidelines

```go
// context_values.go
package main

import (
    "context"
    "fmt"
)

// Good: Use custom types for keys
type contextKey string

const (
    UserIDKey    contextKey = "userID"
    RequestIDKey contextKey = "requestID"
    TraceIDKey   contextKey = "traceID"
)

// Good: Type-safe value extraction
func GetUserID(ctx context.Context) (string, bool) {
    userID, ok := ctx.Value(UserIDKey).(string)
    return userID, ok
}

func GetRequestID(ctx context.Context) (string, bool) {
    requestID, ok := ctx.Value(RequestIDKey).(string)
    return requestID, ok
}

// Good: Helper functions for setting values
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, UserIDKey, userID)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, requestID)
}

// Bad: Using string keys directly (anti-pattern)
func badValueUsage(ctx context.Context) {
    // Don't do this - string keys can collide
    ctx = context.WithValue(ctx, "userID", "user123")
    userID := ctx.Value("userID").(string) // Unsafe type assertion
    fmt.Println(userID)
}

// Good: Proper value usage
func goodValueUsage() {
    ctx := context.Background()
    
    // Set values using helper functions
    ctx = WithUserID(ctx, "user123")
    ctx = WithRequestID(ctx, "req-456")
    
    // Extract values safely
    if userID, ok := GetUserID(ctx); ok {
        fmt.Printf("User ID: %s\n", userID)
    }
    
    if requestID, ok := GetRequestID(ctx); ok {
        fmt.Printf("Request ID: %s\n", requestID)
    }
}

func main() {
    fmt.Println("=== Context Values Guidelines ===")
    goodValueUsage()
}
```

### 3. Error Handling with Context

```go
// context_errors.go
package main

import (
    "context"
    "errors"
    "fmt"
    "time"
)

// Custom errors
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)

func processWithTimeout(ctx context.Context, userID string) error {
    // Simulate processing
    select {
    case <-time.After(2 * time.Second):
        if userID == "invalid" {
            return ErrInvalidInput
        }
        if userID == "notfound" {
            return ErrUserNotFound
        }
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func handleContextErrors() {
    // Test with timeout
    fmt.Println("=== Testing with timeout ===")
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    err := processWithTimeout(ctx, "user123")
    if err != nil {
        switch {
        case errors.Is(err, context.DeadlineExceeded):
            fmt.Println("Operation timed out")
        case errors.Is(err, context.Canceled):
            fmt.Println("Operation was cancelled")
        case errors.Is(err, ErrUserNotFound):
            fmt.Println("User not found")
        case errors.Is(err, ErrInvalidInput):
            fmt.Println("Invalid input provided")
        default:
            fmt.Printf("Unknown error: %v\n", err)
        }
    } else {
        fmt.Println("Operation completed successfully")
    }
    
    // Test with cancellation
    fmt.Println("\n=== Testing with cancellation ===")
    ctx2, cancel2 := context.WithCancel(context.Background())
    
    go func() {
        time.Sleep(500 * time.Millisecond)
        cancel2()
    }()
    
    err = processWithTimeout(ctx2, "user456")
    if err != nil {
        switch {
        case errors.Is(err, context.DeadlineExceeded):
            fmt.Println("Operation timed out")
        case errors.Is(err, context.Canceled):
            fmt.Println("Operation was cancelled")
        default:
            fmt.Printf("Other error: %v\n", err)
        }
    }
    
    // Test with business logic error
    fmt.Println("\n=== Testing with business logic error ===")
    ctx3, cancel3 := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel3()
    
    err = processWithTimeout(ctx3, "notfound")
    if err != nil {
        switch {
        case errors.Is(err, context.DeadlineExceeded):
            fmt.Println("Operation timed out")
        case errors.Is(err, context.Canceled):
            fmt.Println("Operation was cancelled")
        case errors.Is(err, ErrUserNotFound):
            fmt.Println("User not found")
        case errors.Is(err, ErrInvalidInput):
            fmt.Println("Invalid input provided")
        default:
            fmt.Printf("Unknown error: %v\n", err)
        }
    }
}

func main() {
    handleContextErrors()
}
```

## Common Anti-patterns

### 1. What NOT to do with Context

```go
// context_antipatterns.go
package main

import (
    "context"
    "fmt"
)

// Anti-pattern 1: Storing context in struct
type BadService struct {
    ctx context.Context // Don't do this!
}

// Anti-pattern 2: Passing nil context
func badFunction() {
    // Don't pass nil context
    // someFunction(nil) // This will panic
    
    // Use context.Background() or context.TODO() instead
    someFunction(context.Background())
}

func someFunction(ctx context.Context) {
    fmt.Println("Function called with context")
}

// Anti-pattern 3: Not checking context cancellation
func badLongRunningOperation(ctx context.Context) {
    for i := 0; i < 1000000; i++ {
        // Bad: Not checking context
        // Do some work...
        
        // This loop will continue even if context is cancelled
    }
}

// Good: Check context regularly
func goodLongRunningOperation(ctx context.Context) error {
    for i := 0; i < 1000000; i++ {
        // Check context every iteration (or every N iterations)
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        
        // Do some work...
    }
    return nil
}

// Anti-pattern 4: Using context for optional parameters
func badAPICall(ctx context.Context) {
    // Don't use context for optional parameters
    // timeout := ctx.Value("timeout").(time.Duration) // Bad!
    
    // Use function parameters instead
}

// Good: Use explicit parameters
func goodAPICall(ctx context.Context, timeout time.Duration, retries int) {
    // Use explicit parameters for configuration
}

// Anti-pattern 5: Not calling cancel function
func badContextUsage() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    // Forgot to call cancel() - this leaks resources!
    
    someFunction(ctx)
    // cancel() should be called here
}

// Good: Always call cancel
func goodContextUsage() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel() // Always call cancel
    
    someFunction(ctx)
}

func main() {
    fmt.Println("=== Context Anti-patterns ===")
    fmt.Println("See comments in code for what NOT to do")
    
    goodContextUsage()
}
```

## Exercises

### Exercise 1: HTTP Client with Context
Implement an HTTP client that:
- Uses context for request timeouts
- Supports request cancellation
- Includes retry logic with exponential backoff
- Propagates context through all operations

### Exercise 2: Background Job Processor
Create a background job processor that:
- Processes jobs from a queue
- Uses context for graceful shutdown
- Supports job timeouts
- Handles context cancellation properly

### Exercise 3: Distributed System Simulator
Build a distributed system simulator that:
- Simulates multiple services
- Uses context for request tracing
- Implements circuit breaker pattern
- Handles cascading failures

### Exercise 4: Real-time Data Pipeline
Implement a real-time data pipeline that:
- Processes streaming data
- Uses context for backpressure handling
- Supports dynamic scaling
- Implements proper cleanup on shutdown

## Key Takeaways

- Context carries deadlines, cancellation signals, and request-scoped values
- Always pass context as the first parameter to functions
- Use `context.Background()` for top-level contexts
- Use `context.TODO()` when you're unsure which context to use
- Always call the cancel function to prevent resource leaks
- Don't store context in structs
- Use custom types for context keys to avoid collisions
- Check context cancellation in long-running operations
- Handle context errors appropriately
- Context is safe for concurrent use
- Use context for cancellation, not for passing optional parameters
- Context values should be request-scoped data, not configuration

## Next Steps

Next, we'll explore [Testing](24-testing.md) to learn about writing comprehensive tests for Go applications, including unit tests, integration tests, and benchmarks.