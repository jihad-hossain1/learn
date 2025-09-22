# 20. Goroutines

Goroutines are lightweight threads managed by the Go runtime. They are one of Go's most powerful features, enabling concurrent programming with minimal overhead.

## What are Goroutines?

Goroutines are functions that run concurrently with other functions. They are:
- Lightweight (start with ~2KB stack)
- Multiplexed onto OS threads by the Go runtime
- Communicate via channels
- Managed by the Go scheduler

## Basic Goroutine Usage

### Creating Goroutines

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello, %s! (%d)\n", name, i+1)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Regular function call (synchronous)
    fmt.Println("=== Synchronous execution ===")
    sayHello("Alice")
    sayHello("Bob")
    
    fmt.Println("\n=== Asynchronous execution ===")
    // Goroutine (asynchronous)
    go sayHello("Charlie")
    go sayHello("Diana")
    
    // Wait for goroutines to complete
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Main function ending")
}
```

### Anonymous Goroutines

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // Anonymous goroutine
    go func() {
        for i := 0; i < 3; i++ {
            fmt.Printf("Anonymous goroutine: %d\n", i)
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // Anonymous goroutine with parameters
    message := "Hello from goroutine"
    go func(msg string) {
        fmt.Println(msg)
    }(message)
    
    // Wait for goroutines
    time.Sleep(500 * time.Millisecond)
}
```

## Goroutine Communication

### Using Channels

```go
package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job)
        time.Sleep(100 * time.Millisecond)
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 5)
    results := make(chan int, 5)
    
    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }
    
    // Send 5 jobs
    for j := 1; j <= 5; j++ {
        jobs <- j
    }
    close(jobs)
    
    // Collect results
    for r := 1; r <= 5; r++ {
        result := <-results
        fmt.Printf("Result: %d\n", result)
    }
}
```

### WaitGroup for Synchronization

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // Decrement counter when done
    
    fmt.Printf("Worker %d starting\n", id)
    time.Sleep(time.Duration(id) * 100 * time.Millisecond)
    fmt.Printf("Worker %d done\n", id)
}

func main() {
    var wg sync.WaitGroup
    
    // Start 5 workers
    for i := 1; i <= 5; i++ {
        wg.Add(1) // Increment counter
        go worker(i, &wg)
    }
    
    fmt.Println("Waiting for workers to complete...")
    wg.Wait() // Wait for all workers to complete
    fmt.Println("All workers completed")
}
```

## Advanced Goroutine Patterns

### Worker Pool Pattern

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Result struct {
    JobID  int
    Output string
    Error  error
}

type WorkerPool struct {
    workerCount int
    jobs        chan Job
    results     chan Result
    wg          sync.WaitGroup
}

func NewWorkerPool(workerCount, jobQueueSize int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        jobs:        make(chan Job, jobQueueSize),
        results:     make(chan Result, jobQueueSize),
    }
}

func (wp *WorkerPool) Start() {
    for i := 0; i < wp.workerCount; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for job := range wp.jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        
        // Simulate work
        processingTime := time.Duration(rand.Intn(500)) * time.Millisecond
        time.Sleep(processingTime)
        
        result := Result{
            JobID:  job.ID,
            Output: fmt.Sprintf("Processed: %s", job.Data),
        }
        
        wp.results <- result
    }
}

func (wp *WorkerPool) AddJob(job Job) {
    wp.jobs <- job
}

func (wp *WorkerPool) GetResult() Result {
    return <-wp.results
}

func (wp *WorkerPool) Close() {
    close(wp.jobs)
    wp.wg.Wait()
    close(wp.results)
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Create worker pool with 3 workers
    pool := NewWorkerPool(3, 10)
    pool.Start()
    
    // Add jobs
    jobCount := 10
    for i := 1; i <= jobCount; i++ {
        job := Job{
            ID:   i,
            Data: fmt.Sprintf("Task-%d", i),
        }
        pool.AddJob(job)
    }
    
    // Collect results
    go func() {
        for i := 0; i < jobCount; i++ {
            result := pool.GetResult()
            fmt.Printf("Result for job %d: %s\n", result.JobID, result.Output)
        }
    }()
    
    // Close pool and wait
    pool.Close()
    
    fmt.Println("All jobs completed")
}
```

### Pipeline Pattern

```go
package main

import (
    "fmt"
    "time"
)

// Stage 1: Generate numbers
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

// Stage 2: Square numbers
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
            time.Sleep(100 * time.Millisecond) // Simulate work
        }
    }()
    return out
}

// Stage 3: Add constant
func addConstant(in <-chan int, constant int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n + constant
            time.Sleep(50 * time.Millisecond) // Simulate work
        }
    }()
    return out
}

// Stage 4: Filter even numbers
func filterEven(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
    }()
    return out
}

func main() {
    // Create pipeline
    numbers := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squared := square(numbers)
    added := addConstant(squared, 10)
    filtered := filterEven(added)
    
    // Consume results
    fmt.Println("Pipeline results:")
    for result := range filtered {
        fmt.Printf("Result: %d\n", result)
    }
}
```

### Fan-Out/Fan-In Pattern

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Fan-out: distribute work to multiple goroutines
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output
        
        go func(out chan<- int) {
            defer close(out)
            for n := range input {
                // Simulate processing
                processingTime := time.Duration(rand.Intn(500)) * time.Millisecond
                time.Sleep(processingTime)
                
                result := n * n
                fmt.Printf("Worker processed %d -> %d\n", n, result)
                out <- result
            }
        }(output)
    }
    
    return outputs
}

// Fan-in: merge multiple channels into one
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    // Start a goroutine for each input channel
    for _, input := range inputs {
        wg.Add(1)
        go func(in <-chan int) {
            defer wg.Done()
            for value := range in {
                output <- value
            }
        }(input)
    }
    
    // Close output channel when all inputs are done
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    // Create input channel
    input := make(chan int)
    
    // Start input generator
    go func() {
        defer close(input)
        for i := 1; i <= 10; i++ {
            input <- i
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // Fan-out to 3 workers
    workerOutputs := fanOut(input, 3)
    
    // Fan-in results
    results := fanIn(workerOutputs...)
    
    // Collect all results
    fmt.Println("\nFinal results:")
    for result := range results {
        fmt.Printf("Final result: %d\n", result)
    }
}
```

## Context for Cancellation

### Basic Context Usage

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
            fmt.Printf("Worker %d: Received cancellation signal\n", id)
            return
        default:
            fmt.Printf("Worker %d: Working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    // Create context with cancellation
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start workers
    for i := 1; i <= 3; i++ {
        go worker(ctx, i)
    }
    
    // Let workers run for 2 seconds
    time.Sleep(2 * time.Second)
    
    // Cancel all workers
    fmt.Println("Cancelling workers...")
    cancel()
    
    // Give time for cleanup
    time.Sleep(1 * time.Second)
    fmt.Println("Main function ending")
}
```

### Context with Timeout

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func longRunningTask(ctx context.Context, taskID int) error {
    fmt.Printf("Task %d: Starting\n", taskID)
    
    for i := 0; i < 10; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("Task %d: Cancelled due to %v\n", taskID, ctx.Err())
            return ctx.Err()
        default:
            fmt.Printf("Task %d: Step %d\n", taskID, i+1)
            time.Sleep(300 * time.Millisecond)
        }
    }
    
    fmt.Printf("Task %d: Completed successfully\n", taskID)
    return nil
}

func main() {
    // Context with 2-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    // Start multiple tasks
    for i := 1; i <= 3; i++ {
        go func(id int) {
            err := longRunningTask(ctx, id)
            if err != nil {
                fmt.Printf("Task %d failed: %v\n", id, err)
            }
        }(i)
    }
    
    // Wait for context to expire or tasks to complete
    <-ctx.Done()
    fmt.Printf("Context cancelled: %v\n", ctx.Err())
    
    // Give time for cleanup
    time.Sleep(500 * time.Millisecond)
}
```

### Context with Values

```go
package main

import (
    "context"
    "fmt"
    "time"
)

type contextKey string

const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

func processRequest(ctx context.Context) {
    userID := ctx.Value(userIDKey)
    requestID := ctx.Value(requestIDKey)
    
    fmt.Printf("Processing request %v for user %v\n", requestID, userID)
    
    // Simulate work
    select {
    case <-time.After(1 * time.Second):
        fmt.Printf("Request %v completed\n", requestID)
    case <-ctx.Done():
        fmt.Printf("Request %v cancelled: %v\n", requestID, ctx.Err())
    }
}

func main() {
    // Create context with values
    ctx := context.Background()
    ctx = context.WithValue(ctx, userIDKey, "user123")
    ctx = context.WithValue(ctx, requestIDKey, "req456")
    
    // Add timeout
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel()
    
    // Process request
    go processRequest(ctx)
    
    // Wait for completion or timeout
    <-ctx.Done()
    time.Sleep(100 * time.Millisecond)
}
```

## Error Handling in Goroutines

### Error Channel Pattern

```go
package main

import (
    "errors"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Task struct {
    ID   int
    Data string
}

type TaskResult struct {
    TaskID int
    Result string
    Error  error
}

func processTask(task Task) TaskResult {
    // Simulate random failure
    if rand.Float32() < 0.3 {
        return TaskResult{
            TaskID: task.ID,
            Error:  errors.New("random processing error"),
        }
    }
    
    // Simulate work
    time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    
    return TaskResult{
        TaskID: task.ID,
        Result: fmt.Sprintf("Processed: %s", task.Data),
    }
}

func worker(tasks <-chan Task, results chan<- TaskResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for task := range tasks {
        result := processTask(task)
        results <- result
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    tasks := make(chan Task, 10)
    results := make(chan TaskResult, 10)
    
    var wg sync.WaitGroup
    
    // Start 3 workers
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go worker(tasks, results, &wg)
    }
    
    // Send tasks
    taskCount := 10
    go func() {
        defer close(tasks)
        for i := 1; i <= taskCount; i++ {
            tasks <- Task{
                ID:   i,
                Data: fmt.Sprintf("task-%d", i),
            }
        }
    }()
    
    // Close results channel when all workers are done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results and handle errors
    successCount := 0
    errorCount := 0
    
    for result := range results {
        if result.Error != nil {
            fmt.Printf("Task %d failed: %v\n", result.TaskID, result.Error)
            errorCount++
        } else {
            fmt.Printf("Task %d succeeded: %s\n", result.TaskID, result.Result)
            successCount++
        }
    }
    
    fmt.Printf("\nSummary: %d successful, %d failed\n", successCount, errorCount)
}
```

## Goroutine Leaks and Best Practices

### Avoiding Goroutine Leaks

```go
package main

import (
    "context"
    "fmt"
    "runtime"
    "time"
)

// BAD: Goroutine leak
func badExample() {
    ch := make(chan int)
    
    go func() {
        // This goroutine will block forever
        // because no one is reading from the channel
        ch <- 42
    }()
    
    // Channel is never read, goroutine leaks
}

// GOOD: Proper cleanup
func goodExample() {
    ch := make(chan int, 1) // Buffered channel
    
    go func() {
        ch <- 42
    }()
    
    // Read from channel or use select with timeout
    select {
    case value := <-ch:
        fmt.Printf("Received: %d\n", value)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout")
    }
}

// BETTER: Use context for cancellation
func bestExample() {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    ch := make(chan int)
    
    go func() {
        select {
        case ch <- 42:
            fmt.Println("Sent value")
        case <-ctx.Done():
            fmt.Println("Goroutine cancelled")
            return
        }
    }()
    
    select {
    case value := <-ch:
        fmt.Printf("Received: %d\n", value)
    case <-ctx.Done():
        fmt.Printf("Context cancelled: %v\n", ctx.Err())
    }
}

func printGoroutineCount(label string) {
    fmt.Printf("%s: %d goroutines\n", label, runtime.NumGoroutine())
}

func main() {
    printGoroutineCount("Start")
    
    // This will create a goroutine leak
    badExample()
    time.Sleep(100 * time.Millisecond)
    printGoroutineCount("After bad example")
    
    // This won't leak
    goodExample()
    time.Sleep(100 * time.Millisecond)
    printGoroutineCount("After good example")
    
    // This is the best approach
    bestExample()
    time.Sleep(100 * time.Millisecond)
    printGoroutineCount("After best example")
    
    // Force garbage collection
    runtime.GC()
    time.Sleep(100 * time.Millisecond)
    printGoroutineCount("After GC")
}
```

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

type Server struct {
    workers []Worker
    wg      sync.WaitGroup
    ctx     context.Context
    cancel  context.CancelFunc
}

type Worker struct {
    id   int
    name string
}

func (w Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %s starting\n", w.name)
    
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %s shutting down gracefully\n", w.name)
            return
        case <-ticker.C:
            fmt.Printf("Worker %s: tick\n", w.name)
        }
    }
}

func NewServer() *Server {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &Server{
        workers: []Worker{
            {id: 1, name: "worker-1"},
            {id: 2, name: "worker-2"},
            {id: 3, name: "worker-3"},
        },
        ctx:    ctx,
        cancel: cancel,
    }
}

func (s *Server) Start() {
    fmt.Println("Server starting...")
    
    for _, worker := range s.workers {
        s.wg.Add(1)
        go worker.Start(s.ctx, &s.wg)
    }
    
    fmt.Println("All workers started")
}

func (s *Server) Shutdown() {
    fmt.Println("Server shutting down...")
    
    // Cancel context to signal all workers
    s.cancel()
    
    // Wait for all workers to finish
    done := make(chan struct{})
    go func() {
        s.wg.Wait()
        close(done)
    }()
    
    // Wait for graceful shutdown or timeout
    select {
    case <-done:
        fmt.Println("All workers stopped gracefully")
    case <-time.After(5 * time.Second):
        fmt.Println("Shutdown timeout, forcing exit")
    }
}

func main() {
    server := NewServer()
    server.Start()
    
    // Set up signal handling
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    // Wait for signal
    sig := <-sigChan
    fmt.Printf("\nReceived signal: %v\n", sig)
    
    // Graceful shutdown
    server.Shutdown()
    
    fmt.Println("Server stopped")
}
```

## Performance Considerations

### Goroutine Pool vs Unlimited Goroutines

```go
package main

import (
    "fmt"
    "runtime"
    "sync"
    "time"
)

// Unlimited goroutines (can cause resource exhaustion)
func unlimitedGoroutines(tasks []int) {
    var wg sync.WaitGroup
    
    start := time.Now()
    
    for _, task := range tasks {
        wg.Add(1)
        go func(t int) {
            defer wg.Done()
            // Simulate work
            time.Sleep(10 * time.Millisecond)
        }(task)
    }
    
    wg.Wait()
    
    fmt.Printf("Unlimited goroutines: %v, Peak goroutines: %d\n", 
        time.Since(start), runtime.NumGoroutine())
}

// Limited goroutine pool
func limitedGoroutinePool(tasks []int, poolSize int) {
    var wg sync.WaitGroup
    taskChan := make(chan int, len(tasks))
    
    start := time.Now()
    
    // Start worker pool
    for i := 0; i < poolSize; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for task := range taskChan {
                // Simulate work
                time.Sleep(10 * time.Millisecond)
                _ = task
            }
        }()
    }
    
    // Send tasks
    for _, task := range tasks {
        taskChan <- task
    }
    close(taskChan)
    
    wg.Wait()
    
    fmt.Printf("Limited pool (%d workers): %v, Peak goroutines: %d\n", 
        poolSize, time.Since(start), runtime.NumGoroutine())
}

func main() {
    // Create 1000 tasks
    tasks := make([]int, 1000)
    for i := range tasks {
        tasks[i] = i
    }
    
    fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())
    
    // Test unlimited goroutines
    unlimitedGoroutines(tasks)
    runtime.GC()
    time.Sleep(100 * time.Millisecond)
    
    // Test limited pool
    limitedGoroutinePool(tasks, 10)
    runtime.GC()
    time.Sleep(100 * time.Millisecond)
    
    fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())
}
```

## Best Practices

### 1. Always Handle Goroutine Lifecycle

```go
// Good: Use WaitGroup or channels to wait for completion
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // work here
}()
wg.Wait()
```

### 2. Use Context for Cancellation

```go
// Good: Use context for cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    select {
    case <-ctx.Done():
        return
    case <-time.After(1 * time.Second):
        // work here
    }
}()
```

### 3. Limit Goroutine Creation

```go
// Good: Use worker pools for high-volume tasks
const maxWorkers = 10
semaphore := make(chan struct{}, maxWorkers)

for _, task := range tasks {
    semaphore <- struct{}{} // Acquire
    go func(t Task) {
        defer func() { <-semaphore }() // Release
        processTask(t)
    }(task)
}
```

### 4. Avoid Shared State

```go
// Good: Use channels for communication
type Counter struct {
    ch chan int
    value int
}

func (c *Counter) Increment() {
    c.ch <- 1
}

func (c *Counter) run() {
    for delta := range c.ch {
        c.value += delta
    }
}
```

## Exercises

### Exercise 1: Parallel File Processing
Create a program that:
- Reads multiple files concurrently
- Processes each file's content in parallel
- Aggregates results safely
- Handles errors gracefully

### Exercise 2: Rate-Limited API Client
Build an API client that:
- Makes concurrent requests
- Respects rate limits using goroutines
- Implements retry logic with exponential backoff
- Supports cancellation via context

### Exercise 3: Real-time Data Pipeline
Implement a data pipeline that:
- Ingests data from multiple sources concurrently
- Processes data through multiple stages
- Handles backpressure
- Provides monitoring and metrics

## Key Takeaways

- Goroutines are lightweight and efficient for concurrent programming
- Use channels for communication between goroutines
- Always handle goroutine lifecycle with WaitGroups or channels
- Use context for cancellation and timeouts
- Avoid goroutine leaks by ensuring proper cleanup
- Limit goroutine creation for resource-intensive tasks
- Prefer communication over shared memory
- Implement graceful shutdown patterns
- Monitor goroutine count in production applications

## Next Steps

Next, we'll explore [Channels](21-channels.md) in detail to learn about Go's powerful communication mechanism for goroutines.