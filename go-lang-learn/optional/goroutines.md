# 15. Goroutines

Goroutines are lightweight threads managed by the Go runtime. They are one of Go's most powerful features, enabling easy concurrent programming. Goroutines are much more efficient than traditional threads, with thousands of goroutines able to run concurrently.

## Goroutine Basics

### Creating Goroutines

```go
package main

import (
    "fmt"
    "time"
)

// Simple function to run as goroutine
func sayHello(name string) {
    for i := 0; i < 5; i++ {
        fmt.Printf("Hello %s! (%d)\n", name, i+1)
        time.Sleep(100 * time.Millisecond)
    }
}

// Function with parameters
func countDown(name string, start int) {
    for i := start; i >= 0; i-- {
        fmt.Printf("%s: %d\n", name, i)
        time.Sleep(200 * time.Millisecond)
    }
    fmt.Printf("%s: Done!\n", name)
}

func main() {
    fmt.Println("=== Basic Goroutines ===")
    
    // Run function normally (synchronous)
    fmt.Println("Running synchronously:")
    sayHello("Alice")
    
    fmt.Println("\nRunning with goroutines:")
    
    // Run function as goroutine (asynchronous)
    go sayHello("Bob")
    go sayHello("Charlie")
    
    // Run multiple goroutines with different functions
    go countDown("Timer1", 3)
    go countDown("Timer2", 5)
    
    // Anonymous function as goroutine
    go func() {
        for i := 0; i < 3; i++ {
            fmt.Printf("Anonymous goroutine: %d\n", i)
            time.Sleep(150 * time.Millisecond)
        }
    }()
    
    // Anonymous function with parameters
    go func(message string, count int) {
        for i := 0; i < count; i++ {
            fmt.Printf("%s: %d\n", message, i)
            time.Sleep(300 * time.Millisecond)
        }
    }("Parameterized", 4)
    
    // Wait for goroutines to complete
    // Note: This is not the best way - we'll learn better methods
    time.Sleep(2 * time.Second)
    
    fmt.Println("Main function ending")
}
```

### Goroutine Lifecycle

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func worker(id int, work chan int, results chan int) {
    for w := range work {
        fmt.Printf("Worker %d processing job %d\n", id, w)
        time.Sleep(time.Second) // Simulate work
        results <- w * 2
    }
    fmt.Printf("Worker %d finished\n", id)
}

func monitorGoroutines() {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for i := 0; i < 10; i++ {
        select {
        case <-ticker.C:
            fmt.Printf("Active goroutines: %d\n", runtime.NumGoroutine())
        }
    }
}

func main() {
    fmt.Println("=== Goroutine Lifecycle ===")
    
    fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())
    
    // Start monitoring goroutines
    go monitorGoroutines()
    
    // Create channels
    work := make(chan int, 5)
    results := make(chan int, 5)
    
    // Start workers
    for i := 1; i <= 3; i++ {
        go worker(i, work, results)
    }
    
    // Send work
    for i := 1; i <= 5; i++ {
        work <- i
    }
    close(work)
    
    // Collect results
    for i := 1; i <= 5; i++ {
        result := <-results
        fmt.Printf("Result: %d\n", result)
    }
    
    time.Sleep(2 * time.Second)
    fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())
}
```

## Goroutine Communication Patterns

### Fan-Out Pattern

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// Producer generates work
func producer(jobs chan<- int, numJobs int) {
    defer close(jobs)
    
    for i := 1; i <= numJobs; i++ {
        fmt.Printf("Producing job %d\n", i)
        jobs <- i
        time.Sleep(100 * time.Millisecond)
    }
    fmt.Println("Producer finished")
}

// Worker processes jobs
func worker(id int, jobs <-chan int, results chan<- string) {
    for job := range jobs {
        // Simulate variable processing time
        processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
        time.Sleep(processingTime)
        
        result := fmt.Sprintf("Worker %d processed job %d in %v", id, job, processingTime)
        fmt.Println(result)
        results <- result
    }
    fmt.Printf("Worker %d finished\n", id)
}

func main() {
    fmt.Println("=== Fan-Out Pattern ===")
    
    rand.Seed(time.Now().UnixNano())
    
    const numJobs = 10
    const numWorkers = 3
    
    jobs := make(chan int, numJobs)
    results := make(chan string, numJobs)
    
    // Start producer
    go producer(jobs, numJobs)
    
    // Start workers (fan-out)
    for i := 1; i <= numWorkers; i++ {
        go worker(i, jobs, results)
    }
    
    // Collect results
    for i := 1; i <= numJobs; i++ {
        result := <-results
        fmt.Printf("Collected: %s\n", result)
    }
    
    fmt.Println("All jobs completed")
}
```

### Fan-In Pattern

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// Data source that generates values
func dataSource(name string, output chan<- string, duration time.Duration) {
    defer close(output)
    
    for i := 1; i <= 5; i++ {
        data := fmt.Sprintf("%s-data-%d", name, i)
        fmt.Printf("%s producing: %s\n", name, data)
        output <- data
        time.Sleep(duration)
    }
    fmt.Printf("%s finished\n", name)
}

// Fan-in function that merges multiple channels
func fanIn(inputs ...<-chan string) <-chan string {
    output := make(chan string)
    
    // Start a goroutine for each input channel
    for i, input := range inputs {
        go func(id int, ch <-chan string) {
            for data := range ch {
                output <- fmt.Sprintf("[Source-%d] %s", id+1, data)
            }
        }(i, input)
    }
    
    // Start a goroutine to close output when all inputs are done
    go func() {
        // Wait for all input channels to close
        // This is a simplified approach - in production, use sync.WaitGroup
        time.Sleep(3 * time.Second)
        close(output)
    }()
    
    return output
}

// Alternative fan-in using select
func fanInSelect(input1, input2 <-chan string) <-chan string {
    output := make(chan string)
    
    go func() {
        defer close(output)
        
        for {
            select {
            case data, ok := <-input1:
                if !ok {
                    input1 = nil
                } else {
                    output <- fmt.Sprintf("[Input1] %s", data)
                }
            case data, ok := <-input2:
                if !ok {
                    input2 = nil
                } else {
                    output <- fmt.Sprintf("[Input2] %s", data)
                }
            }
            
            // Exit when both channels are closed
            if input1 == nil && input2 == nil {
                break
            }
        }
    }()
    
    return output
}

func main() {
    fmt.Println("=== Fan-In Pattern ===")
    
    rand.Seed(time.Now().UnixNano())
    
    // Create multiple data sources
    source1 := make(chan string)
    source2 := make(chan string)
    source3 := make(chan string)
    
    // Start data sources
    go dataSource("Source1", source1, 200*time.Millisecond)
    go dataSource("Source2", source2, 300*time.Millisecond)
    go dataSource("Source3", source3, 150*time.Millisecond)
    
    // Fan-in all sources
    merged := fanIn(source1, source2, source3)
    
    // Process merged data
    fmt.Println("\nProcessing merged data:")
    for data := range merged {
        fmt.Printf("Received: %s\n", data)
    }
    
    fmt.Println("\n=== Fan-In with Select ===")
    
    // Demonstrate select-based fan-in
    sourceA := make(chan string)
    sourceB := make(chan string)
    
    go dataSource("SourceA", sourceA, 250*time.Millisecond)
    go dataSource("SourceB", sourceB, 180*time.Millisecond)
    
    mergedSelect := fanInSelect(sourceA, sourceB)
    
    for data := range mergedSelect {
        fmt.Printf("Select received: %s\n", data)
    }
    
    fmt.Println("All data processed")
}
```

### Pipeline Pattern

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
    "time"
)

// Stage 1: Generate numbers
func generateNumbers(count int) <-chan int {
    output := make(chan int)
    
    go func() {
        defer close(output)
        for i := 1; i <= count; i++ {
            fmt.Printf("Generating: %d\n", i)
            output <- i
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    return output
}

// Stage 2: Square numbers
func squareNumbers(input <-chan int) <-chan int {
    output := make(chan int)
    
    go func() {
        defer close(output)
        for num := range input {
            squared := num * num
            fmt.Printf("Squaring %d = %d\n", num, squared)
            output <- squared
            time.Sleep(50 * time.Millisecond)
        }
    }()
    
    return output
}

// Stage 3: Convert to string
func convertToString(input <-chan int) <-chan string {
    output := make(chan string)
    
    go func() {
        defer close(output)
        for num := range input {
            str := strconv.Itoa(num)
            fmt.Printf("Converting %d to string: %s\n", num, str)
            output <- str
            time.Sleep(75 * time.Millisecond)
        }
    }()
    
    return output
}

// Stage 4: Add prefix
func addPrefix(input <-chan string, prefix string) <-chan string {
    output := make(chan string)
    
    go func() {
        defer close(output)
        for str := range input {
            prefixed := prefix + str
            fmt.Printf("Adding prefix to %s: %s\n", str, prefixed)
            output <- prefixed
            time.Sleep(25 * time.Millisecond)
        }
    }()
    
    return output
}

// Complex pipeline with multiple transformations
func textProcessingPipeline(texts []string) <-chan string {
    // Stage 1: Input
    input := make(chan string)
    go func() {
        defer close(input)
        for _, text := range texts {
            input <- text
            time.Sleep(50 * time.Millisecond)
        }
    }()
    
    // Stage 2: Convert to uppercase
    uppercase := make(chan string)
    go func() {
        defer close(uppercase)
        for text := range input {
            upper := strings.ToUpper(text)
            fmt.Printf("Uppercase: %s -> %s\n", text, upper)
            uppercase <- upper
        }
    }()
    
    // Stage 3: Add word count
    withCount := make(chan string)
    go func() {
        defer close(withCount)
        for text := range uppercase {
            words := strings.Fields(text)
            result := fmt.Sprintf("%s [%d words]", text, len(words))
            fmt.Printf("Word count: %s\n", result)
            withCount <- result
        }
    }()
    
    // Stage 4: Add timestamp
    final := make(chan string)
    go func() {
        defer close(final)
        for text := range withCount {
            timestamp := time.Now().Format("15:04:05")
            result := fmt.Sprintf("[%s] %s", timestamp, text)
            fmt.Printf("Timestamped: %s\n", result)
            final <- result
        }
    }()
    
    return final
}

func main() {
    fmt.Println("=== Pipeline Pattern ===")
    
    // Simple number processing pipeline
    fmt.Println("Number processing pipeline:")
    numbers := generateNumbers(5)
    squared := squareNumbers(numbers)
    strings := convertToString(squared)
    prefixed := addPrefix(strings, "NUM-")
    
    // Consume final results
    for result := range prefixed {
        fmt.Printf("Final result: %s\n", result)
    }
    
    fmt.Println("\n=== Text Processing Pipeline ===")
    
    // Text processing pipeline
    texts := []string{
        "hello world",
        "go programming",
        "concurrent processing",
        "pipeline pattern",
    }
    
    processed := textProcessingPipeline(texts)
    
    fmt.Println("\nFinal processed texts:")
    for result := range processed {
        fmt.Printf("FINAL: %s\n", result)
    }
    
    fmt.Println("Pipeline processing completed")
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

// Job represents work to be done
type Job struct {
    ID       int
    Data     string
    Priority int
}

// Result represents the result of processing a job
type Result struct {
    Job       Job
    Output    string
    Duration  time.Duration
    WorkerID  int
    Error     error
}

// Worker represents a worker in the pool
type Worker struct {
    ID       int
    JobChan  <-chan Job
    ResultChan chan<- Result
    QuitChan <-chan bool
    wg       *sync.WaitGroup
}

// NewWorker creates a new worker
func NewWorker(id int, jobChan <-chan Job, resultChan chan<- Result, quitChan <-chan bool, wg *sync.WaitGroup) *Worker {
    return &Worker{
        ID:         id,
        JobChan:    jobChan,
        ResultChan: resultChan,
        QuitChan:   quitChan,
        wg:         wg,
    }
}

// Start begins the worker's processing loop
func (w *Worker) Start() {
    go func() {
        defer w.wg.Done()
        
        for {
            select {
            case job := <-w.JobChan:
                fmt.Printf("Worker %d started job %d\n", w.ID, job.ID)
                result := w.processJob(job)
                w.ResultChan <- result
                
            case <-w.QuitChan:
                fmt.Printf("Worker %d stopping\n", w.ID)
                return
            }
        }
    }()
}

// processJob simulates job processing
func (w *Worker) processJob(job Job) Result {
    start := time.Now()
    
    // Simulate variable processing time based on priority
    processingTime := time.Duration(rand.Intn(1000)+500) * time.Millisecond
    if job.Priority > 5 {
        processingTime /= 2 // High priority jobs process faster
    }
    
    time.Sleep(processingTime)
    
    // Simulate occasional errors
    var err error
    if rand.Float32() < 0.1 { // 10% chance of error
        err = fmt.Errorf("processing error for job %d", job.ID)
    }
    
    output := fmt.Sprintf("Processed: %s (Priority: %d)", job.Data, job.Priority)
    
    return Result{
        Job:      job,
        Output:   output,
        Duration: time.Since(start),
        WorkerID: w.ID,
        Error:    err,
    }
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
    NumWorkers int
    JobChan    chan Job
    ResultChan chan Result
    QuitChan   chan bool
    workers    []*Worker
    wg         sync.WaitGroup
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers, jobBufferSize, resultBufferSize int) *WorkerPool {
    return &WorkerPool{
        NumWorkers: numWorkers,
        JobChan:    make(chan Job, jobBufferSize),
        ResultChan: make(chan Result, resultBufferSize),
        QuitChan:   make(chan bool),
        workers:    make([]*Worker, numWorkers),
    }
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start() {
    fmt.Printf("Starting worker pool with %d workers\n", wp.NumWorkers)
    
    for i := 0; i < wp.NumWorkers; i++ {
        wp.wg.Add(1)
        worker := NewWorker(i+1, wp.JobChan, wp.ResultChan, wp.QuitChan, &wp.wg)
        wp.workers[i] = worker
        worker.Start()
    }
}

// Stop gracefully stops all workers
func (wp *WorkerPool) Stop() {
    fmt.Println("Stopping worker pool...")
    
    close(wp.QuitChan)
    wp.wg.Wait()
    close(wp.ResultChan)
    
    fmt.Println("Worker pool stopped")
}

// AddJob adds a job to the pool
func (wp *WorkerPool) AddJob(job Job) {
    wp.JobChan <- job
}

// FinishJobs closes the job channel to indicate no more jobs
func (wp *WorkerPool) FinishJobs() {
    close(wp.JobChan)
}

func main() {
    fmt.Println("=== Worker Pool Pattern ===")
    
    rand.Seed(time.Now().UnixNano())
    
    // Create worker pool
    pool := NewWorkerPool(3, 10, 20)
    pool.Start()
    
    // Start result collector
    var resultWg sync.WaitGroup
    resultWg.Add(1)
    
    go func() {
        defer resultWg.Done()
        
        successCount := 0
        errorCount := 0
        totalDuration := time.Duration(0)
        
        for result := range pool.ResultChan {
            if result.Error != nil {
                fmt.Printf("❌ Job %d failed (Worker %d): %v\n", 
                    result.Job.ID, result.WorkerID, result.Error)
                errorCount++
            } else {
                fmt.Printf("✅ Job %d completed (Worker %d) in %v: %s\n", 
                    result.Job.ID, result.WorkerID, result.Duration, result.Output)
                successCount++
            }
            totalDuration += result.Duration
        }
        
        fmt.Printf("\n=== Results Summary ===\n")
        fmt.Printf("Successful jobs: %d\n", successCount)
        fmt.Printf("Failed jobs: %d\n", errorCount)
        fmt.Printf("Average processing time: %v\n", totalDuration/time.Duration(successCount+errorCount))
    }()
    
    // Generate and submit jobs
    jobs := []Job{
        {ID: 1, Data: "Process user data", Priority: 8},
        {ID: 2, Data: "Generate report", Priority: 3},
        {ID: 3, Data: "Send email", Priority: 9},
        {ID: 4, Data: "Update database", Priority: 7},
        {ID: 5, Data: "Backup files", Priority: 2},
        {ID: 6, Data: "Analyze logs", Priority: 5},
        {ID: 7, Data: "Compress images", Priority: 4},
        {ID: 8, Data: "Sync data", Priority: 6},
        {ID: 9, Data: "Clean cache", Priority: 1},
        {ID: 10, Data: "Validate input", Priority: 10},
    }
    
    fmt.Printf("Submitting %d jobs...\n", len(jobs))
    for _, job := range jobs {
        pool.AddJob(job)
        time.Sleep(50 * time.Millisecond) // Simulate job arrival rate
    }
    
    // Signal no more jobs
    pool.FinishJobs()
    
    // Wait a bit for processing to complete
    time.Sleep(3 * time.Second)
    
    // Stop the pool
    pool.Stop()
    
    // Wait for result collector to finish
    resultWg.Wait()
    
    fmt.Println("All jobs processed")
}
```

### Rate Limiting with Goroutines

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// RateLimiter controls the rate of operations
type RateLimiter struct {
    ticker   *time.Ticker
    tokens   chan struct{}
    capacity int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate time.Duration, capacity int) *RateLimiter {
    rl := &RateLimiter{
        ticker:   time.NewTicker(rate),
        tokens:   make(chan struct{}, capacity),
        capacity: capacity,
    }
    
    // Fill initial tokens
    for i := 0; i < capacity; i++ {
        rl.tokens <- struct{}{}
    }
    
    // Start token refill goroutine
    go rl.refillTokens()
    
    return rl
}

// refillTokens adds tokens at the specified rate
func (rl *RateLimiter) refillTokens() {
    for range rl.ticker.C {
        select {
        case rl.tokens <- struct{}{}:
            // Token added
        default:
            // Token bucket is full
        }
    }
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait() {
    <-rl.tokens
}

// TryWait attempts to get a token without blocking
func (rl *RateLimiter) TryWait() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
    rl.ticker.Stop()
    close(rl.tokens)
}

// APIClient simulates an API client with rate limiting
type APIClient struct {
    rateLimiter *RateLimiter
    name        string
}

// NewAPIClient creates a new API client
func NewAPIClient(name string, rateLimiter *RateLimiter) *APIClient {
    return &APIClient{
        rateLimiter: rateLimiter,
        name:        name,
    }
}

// MakeRequest simulates making an API request
func (c *APIClient) MakeRequest(requestID int) {
    fmt.Printf("%s: Waiting for rate limit (Request %d)\n", c.name, requestID)
    
    start := time.Now()
    c.rateLimiter.Wait()
    waitTime := time.Since(start)
    
    fmt.Printf("%s: Making request %d (waited %v)\n", c.name, requestID, waitTime)
    
    // Simulate API call
    time.Sleep(100 * time.Millisecond)
    
    fmt.Printf("%s: Request %d completed\n", c.name, requestID)
}

// TryMakeRequest attempts to make a request without waiting
func (c *APIClient) TryMakeRequest(requestID int) bool {
    if c.rateLimiter.TryWait() {
        fmt.Printf("%s: Making immediate request %d\n", c.name, requestID)
        time.Sleep(100 * time.Millisecond)
        fmt.Printf("%s: Immediate request %d completed\n", c.name, requestID)
        return true
    }
    
    fmt.Printf("%s: Request %d skipped (rate limited)\n", c.name, requestID)
    return false
}

// Semaphore for limiting concurrent operations
type Semaphore struct {
    ch chan struct{}
}

// NewSemaphore creates a new semaphore with the given capacity
func NewSemaphore(capacity int) *Semaphore {
    return &Semaphore{
        ch: make(chan struct{}, capacity),
    }
}

// Acquire acquires a semaphore slot
func (s *Semaphore) Acquire() {
    s.ch <- struct{}{}
}

// Release releases a semaphore slot
func (s *Semaphore) Release() {
    <-s.ch
}

// TryAcquire attempts to acquire a slot without blocking
func (s *Semaphore) TryAcquire() bool {
    select {
    case s.ch <- struct{}{}:
        return true
    default:
        return false
    }
}

// ConcurrentProcessor demonstrates semaphore usage
type ConcurrentProcessor struct {
    semaphore *Semaphore
    name      string
}

// NewConcurrentProcessor creates a new processor
func NewConcurrentProcessor(name string, maxConcurrent int) *ConcurrentProcessor {
    return &ConcurrentProcessor{
        semaphore: NewSemaphore(maxConcurrent),
        name:      name,
    }
}

// ProcessTask processes a task with concurrency control
func (cp *ConcurrentProcessor) ProcessTask(taskID int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("%s: Task %d waiting for slot\n", cp.name, taskID)
    cp.semaphore.Acquire()
    defer cp.semaphore.Release()
    
    fmt.Printf("%s: Task %d started processing\n", cp.name, taskID)
    
    // Simulate processing time
    time.Sleep(time.Duration(500+taskID*100) * time.Millisecond)
    
    fmt.Printf("%s: Task %d completed\n", cp.name, taskID)
}

func main() {
    fmt.Println("=== Rate Limiting and Concurrency Control ===")
    
    // Demonstrate rate limiting
    fmt.Println("\n1. Rate Limiting Example:")
    
    // Create rate limiter: 1 token every 500ms, capacity of 3
    rateLimiter := NewRateLimiter(500*time.Millisecond, 3)
    defer rateLimiter.Stop()
    
    // Create API clients
    client1 := NewAPIClient("Client1", rateLimiter)
    client2 := NewAPIClient("Client2", rateLimiter)
    
    var wg sync.WaitGroup
    
    // Client 1 makes requests
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 1; i <= 5; i++ {
            client1.MakeRequest(i)
        }
    }()
    
    // Client 2 makes requests
    wg.Add(1)
    go func() {
        defer wg.Done()
        for i := 1; i <= 5; i++ {
            client2.MakeRequest(i)
        }
    }()
    
    wg.Wait()
    
    fmt.Println("\n2. Try-based Rate Limiting:")
    
    // Demonstrate try-based requests
    for i := 1; i <= 10; i++ {
        client1.TryMakeRequest(i)
        time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("\n3. Semaphore-based Concurrency Control:")
    
    // Demonstrate semaphore usage
    processor := NewConcurrentProcessor("Processor", 2) // Max 2 concurrent tasks
    
    var processingWg sync.WaitGroup
    
    // Start multiple tasks
    for i := 1; i <= 8; i++ {
        processingWg.Add(1)
        go processor.ProcessTask(i, &processingWg)
        time.Sleep(50 * time.Millisecond) // Stagger task starts
    }
    
    processingWg.Wait()
    
    fmt.Println("\n4. Combined Rate Limiting and Concurrency Control:")
    
    // Create a new rate limiter for this example
    combinedRateLimiter := NewRateLimiter(200*time.Millisecond, 2)
    defer combinedRateLimiter.Stop()
    
    combinedSemaphore := NewSemaphore(3)
    
    var combinedWg sync.WaitGroup
    
    // Function that uses both rate limiting and semaphore
    processWithLimits := func(id int) {
        defer combinedWg.Done()
        
        // Wait for rate limit
        fmt.Printf("Task %d: Waiting for rate limit\n", id)
        combinedRateLimiter.Wait()
        
        // Acquire semaphore
        fmt.Printf("Task %d: Waiting for semaphore\n", id)
        combinedSemaphore.Acquire()
        defer combinedSemaphore.Release()
        
        fmt.Printf("Task %d: Processing\n", id)
        time.Sleep(300 * time.Millisecond)
        fmt.Printf("Task %d: Completed\n", id)
    }
    
    // Start tasks
    for i := 1; i <= 6; i++ {
        combinedWg.Add(1)
        go processWithLimits(i)
    }
    
    combinedWg.Wait()
    
    fmt.Println("\nAll rate limiting and concurrency control examples completed")
}
```

## Goroutine Best Practices

### Proper Goroutine Management

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

// Example 1: Using context for cancellation
func workerWithContext(ctx context.Context, id int, work <-chan string, results chan<- string) {
    for {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: Context cancelled, stopping\n", id)
            return
        case job, ok := <-work:
            if !ok {
                fmt.Printf("Worker %d: Work channel closed, stopping\n", id)
                return
            }
            
            // Simulate work with context checking
            select {
            case <-ctx.Done():
                fmt.Printf("Worker %d: Context cancelled during work\n", id)
                return
            case <-time.After(500 * time.Millisecond):
                result := fmt.Sprintf("Worker %d processed: %s", id, job)
                
                select {
                case results <- result:
                    fmt.Printf("Worker %d: Result sent\n", id)
                case <-ctx.Done():
                    fmt.Printf("Worker %d: Context cancelled while sending result\n", id)
                    return
                }
            }
        }
    }
}

// Example 2: Proper error handling in goroutines
type WorkerError struct {
    WorkerID int
    Job      string
    Err      error
}

func (we WorkerError) Error() string {
    return fmt.Sprintf("Worker %d error processing '%s': %v", we.WorkerID, we.Job, we.Err)
}

func workerWithErrorHandling(id int, jobs <-chan string, results chan<- string, errors chan<- error, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        // Simulate potential error
        if job == "error" {
            errors <- WorkerError{
                WorkerID: id,
                Job:      job,
                Err:      fmt.Errorf("simulated processing error"),
            }
            continue
        }
        
        // Simulate work
        time.Sleep(200 * time.Millisecond)
        result := fmt.Sprintf("Worker %d: %s -> processed", id, job)
        results <- result
    }
    
    fmt.Printf("Worker %d: Finished processing all jobs\n", id)
}

// Example 3: Resource cleanup
type Resource struct {
    ID   int
    Name string
    file *mockFile
}

type mockFile struct {
    name   string
    closed bool
}

func (mf *mockFile) Close() error {
    if mf.closed {
        return fmt.Errorf("file already closed")
    }
    mf.closed = true
    fmt.Printf("File %s closed\n", mf.name)
    return nil
}

func NewResource(id int, name string) *Resource {
    return &Resource{
        ID:   id,
        Name: name,
        file: &mockFile{name: fmt.Sprintf("file_%d.txt", id)},
    }
}

func (r *Resource) Close() error {
    return r.file.Close()
}

func resourceWorker(ctx context.Context, id int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // Create resource
    resource := NewResource(id, fmt.Sprintf("Resource-%d", id))
    defer func() {
        if err := resource.Close(); err != nil {
            fmt.Printf("Error closing resource %d: %v\n", id, err)
        }
    }()
    
    fmt.Printf("Worker %d: Resource created and working\n", id)
    
    // Simulate work with context checking
    ticker := time.NewTicker(300 * time.Millisecond)
    defer ticker.Stop()
    
    for i := 0; i < 5; i++ {
        select {
        case <-ctx.Done():
            fmt.Printf("Worker %d: Context cancelled, cleaning up\n", id)
            return
        case <-ticker.C:
            fmt.Printf("Worker %d: Working with %s (iteration %d)\n", id, resource.Name, i+1)
        }
    }
    
    fmt.Printf("Worker %d: Work completed\n", id)
}

// Example 4: Graceful shutdown
type Server struct {
    workers []chan bool
    wg      sync.WaitGroup
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewServer() *Server {
    ctx, cancel := context.WithCancel(context.Background())
    return &Server{
        workers: make([]chan bool, 0),
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (s *Server) Start(numWorkers int) {
    fmt.Printf("Starting server with %d workers\n", numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        quit := make(chan bool)
        s.workers = append(s.workers, quit)
        
        s.wg.Add(1)
        go s.worker(i+1, quit)
    }
}

func (s *Server) worker(id int, quit <-chan bool) {
    defer s.wg.Done()
    
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    fmt.Printf("Worker %d started\n", id)
    
    for {
        select {
        case <-s.ctx.Done():
            fmt.Printf("Worker %d: Server context cancelled\n", id)
            return
        case <-quit:
            fmt.Printf("Worker %d: Received quit signal\n", id)
            return
        case <-ticker.C:
            fmt.Printf("Worker %d: Heartbeat\n", id)
        }
    }
}

func (s *Server) Shutdown() {
    fmt.Println("Initiating server shutdown...")
    
    // Cancel context
    s.cancel()
    
    // Send quit signals to all workers
    for i, quit := range s.workers {
        fmt.Printf("Sending quit signal to worker %d\n", i+1)
        close(quit)
    }
    
    // Wait for all workers to finish
    s.wg.Wait()
    fmt.Println("Server shutdown complete")
}

func main() {
    fmt.Println("=== Goroutine Best Practices ===")
    
    // Example 1: Context-based cancellation
    fmt.Println("\n1. Context-based Cancellation:")
    
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()
    
    work := make(chan string, 5)
    results := make(chan string, 5)
    
    // Start workers
    for i := 1; i <= 2; i++ {
        go workerWithContext(ctx, i, work, results)
    }
    
    // Send work
    jobs := []string{"job1", "job2", "job3", "job4"}
    go func() {
        for _, job := range jobs {
            work <- job
            time.Sleep(300 * time.Millisecond)
        }
        close(work)
    }()
    
    // Collect results until context is done
    go func() {
        for result := range results {
            fmt.Printf("Result: %s\n", result)
        }
    }()
    
    <-ctx.Done()
    fmt.Println("Context timeout reached")
    
    // Example 2: Error handling
    fmt.Println("\n2. Error Handling in Goroutines:")
    
    jobsWithErrors := make(chan string, 10)
    resultsWithErrors := make(chan string, 10)
    errorsChannel := make(chan error, 10)
    
    var errorWg sync.WaitGroup
    
    // Start workers
    for i := 1; i <= 2; i++ {
        errorWg.Add(1)
        go workerWithErrorHandling(i, jobsWithErrors, resultsWithErrors, errorsChannel, &errorWg)
    }
    
    // Send jobs (including error-inducing ones)
    errorJobs := []string{"task1", "error", "task2", "task3", "error", "task4"}
    go func() {
        for _, job := range errorJobs {
            jobsWithErrors <- job
        }
        close(jobsWithErrors)
    }()
    
    // Collect results and errors
    go func() {
        errorWg.Wait()
        close(resultsWithErrors)
        close(errorsChannel)
    }()
    
    // Process results
    for result := range resultsWithErrors {
        fmt.Printf("Success: %s\n", result)
    }
    
    // Process errors
    for err := range errorsChannel {
        fmt.Printf("Error: %v\n", err)
    }
    
    // Example 3: Resource cleanup
    fmt.Println("\n3. Resource Cleanup:")
    
    resourceCtx, resourceCancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer resourceCancel()
    
    var resourceWg sync.WaitGroup
    
    for i := 1; i <= 3; i++ {
        resourceWg.Add(1)
        go resourceWorker(resourceCtx, i, &resourceWg)
    }
    
    resourceWg.Wait()
    
    // Example 4: Graceful shutdown
    fmt.Println("\n4. Graceful Shutdown:")
    
    server := NewServer()
    server.Start(3)
    
    // Let server run for a bit
    time.Sleep(2 * time.Second)
    
    // Shutdown server
    server.Shutdown()
    
    fmt.Println("\nAll best practices examples completed")
}
```

## Exercises

### Exercise 1: Download Manager
Create a concurrent download manager using goroutines.

```go
// Create a download manager that:
// - Downloads multiple files concurrently
// - Limits the number of concurrent downloads
// - Provides progress reporting
// - Handles errors gracefully
// - Supports cancellation

package main

type DownloadManager struct {
    // Add fields for concurrency control, progress tracking, etc.
}

type DownloadTask struct {
    URL      string
    Filename string
    Size     int64
}

type DownloadProgress struct {
    Task       DownloadTask
    Downloaded int64
    Total      int64
    Speed      float64 // bytes per second
    Status     string  // "downloading", "completed", "failed"
    Error      error
}

// Implement these methods:
// NewDownloadManager(maxConcurrent int) *DownloadManager
// Download(tasks []DownloadTask) <-chan DownloadProgress
// Cancel()
// GetStats() DownloadStats
```

### Exercise 2: Log Aggregator
Create a log aggregator that collects logs from multiple sources.

```go
// Create a log aggregator that:
// - Collects logs from multiple sources concurrently
// - Buffers logs and flushes periodically
// - Filters logs by level
// - Supports multiple output destinations
// - Handles backpressure

package main

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Source    string
    Message   string
    Fields    map[string]interface{}
}

type LogAggregator struct {
    // Add fields for buffering, filtering, etc.
}

// Implement these methods:
// NewLogAggregator(config Config) *LogAggregator
// AddSource(source LogSource) error
// AddDestination(dest LogDestination) error
// Start() error
// Stop() error
// GetMetrics() AggregatorMetrics
```

### Exercise 3: Task Scheduler
Create a task scheduler with priority queues and retry logic.

```go
// Create a task scheduler that:
// - Schedules tasks with different priorities
// - Supports delayed execution
// - Implements retry logic with exponential backoff
// - Provides task status tracking
// - Supports task cancellation

package main

type Task struct {
    ID          string
    Priority    int
    ScheduledAt time.Time
    MaxRetries  int
    Timeout     time.Duration
    Payload     interface{}
    Handler     TaskHandler
}

type TaskScheduler struct {
    // Add fields for priority queues, workers, etc.
}

// Implement these methods:
// NewTaskScheduler(config SchedulerConfig) *TaskScheduler
// ScheduleTask(task Task) error
// ScheduleDelayed(task Task, delay time.Duration) error
// CancelTask(taskID string) error
// GetTaskStatus(taskID string) TaskStatus
// Start() error
// Stop() error
```

## Key Takeaways

1. **Goroutines are lightweight** - you can create thousands of them
2. **Use channels for communication** between goroutines
3. **Always handle goroutine lifecycle** properly
4. **Context is essential** for cancellation and timeouts
5. **Resource cleanup** must be handled with defer statements
6. **Error handling** in goroutines requires careful design
7. **Rate limiting and semaphores** control resource usage
8. **Worker pools** are efficient for processing many tasks
9. **Pipeline patterns** enable stream processing
10. **Graceful shutdown** ensures clean application termination

## Next Steps

Now that you understand goroutines, let's learn about [Channels](16-channels.md) for communication between goroutines!

---

**Previous**: [← Packages and Modules](14-packages-modules.md) | **Next**: [Channels →](16-channels.md)