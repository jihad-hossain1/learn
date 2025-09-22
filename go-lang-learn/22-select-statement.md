# 22. Select Statement

The `select` statement is one of Go's most powerful concurrency features. It allows a goroutine to wait on multiple communication operations, making it essential for coordinating concurrent operations and implementing complex communication patterns.

## Introduction to Select

The `select` statement is similar to a `switch` statement, but it operates on channel operations. It blocks until one of its cases can proceed, then executes that case. If multiple cases are ready, it chooses one at random.

### Basic Select Syntax

```go
// basic_select.go
package main

import (
    "fmt"
    "time"
)

func basicSelect() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // Send data to channels in separate goroutines
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Message from channel 1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message from channel 2"
    }()
    
    // Select waits for the first available channel
    select {
    case msg1 := <-ch1:
        fmt.Println("Received:", msg1)
    case msg2 := <-ch2:
        fmt.Println("Received:", msg2)
    }
    
    fmt.Println("Select completed")
}

func main() {
    fmt.Println("=== Basic Select Example ===")
    basicSelect()
}
```

## Select with Default Case

The `default` case in a `select` statement executes immediately if no other case is ready, making the select non-blocking.

```go
// select_default.go
package main

import (
    "fmt"
    "time"
)

func selectWithDefault() {
    ch := make(chan string)
    
    // Try to receive from channel
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    default:
        fmt.Println("No message available")
    }
    
    // Send a message and try again
    go func() {
        ch <- "Hello from goroutine"
    }()
    
    // Give goroutine time to send
    time.Sleep(100 * time.Millisecond)
    
    select {
    case msg := <-ch:
        fmt.Println("Received:", msg)
    default:
        fmt.Println("No message available")
    }
}

// Non-blocking send
func nonBlockingSend() {
    ch := make(chan string, 1) // Buffered channel with capacity 1
    
    // First send will succeed
    select {
    case ch <- "First message":
        fmt.Println("Sent first message")
    default:
        fmt.Println("Could not send first message")
    }
    
    // Second send will fail (buffer full)
    select {
    case ch <- "Second message":
        fmt.Println("Sent second message")
    default:
        fmt.Println("Could not send second message - buffer full")
    }
    
    // Receive the message
    msg := <-ch
    fmt.Println("Received:", msg)
    
    // Now we can send again
    select {
    case ch <- "Third message":
        fmt.Println("Sent third message")
    default:
        fmt.Println("Could not send third message")
    }
}

func main() {
    fmt.Println("=== Select with Default ===")
    selectWithDefault()
    
    fmt.Println("\n=== Non-blocking Send ===")
    nonBlockingSend()
}
```

## Select with Timeout

Using `time.After()` with select provides a timeout mechanism for channel operations.

```go
// select_timeout.go
package main

import (
    "fmt"
    "time"
)

func selectWithTimeout() {
    ch := make(chan string)
    
    // Simulate a slow operation
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "Data received"
    }()
    
    // Wait for data or timeout
    select {
    case data := <-ch:
        fmt.Println("Received:", data)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout: operation took too long")
    }
}

// Multiple timeouts
func multipleTimeouts() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(500 * time.Millisecond)
        ch1 <- "Quick response"
    }()
    
    go func() {
        time.Sleep(1500 * time.Millisecond)
        ch2 <- "Slow response"
    }()
    
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("From ch1:", msg1)
        case msg2 := <-ch2:
            fmt.Println("From ch2:", msg2)
        case <-time.After(1 * time.Second):
            fmt.Println("Timeout waiting for response")
        }
    }
}

// Periodic operations with timeout
func periodicWithTimeout() {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    timeout := time.After(3 * time.Second)
    counter := 0
    
    for {
        select {
        case <-ticker.C:
            counter++
            fmt.Printf("Tick %d\n", counter)
        case <-timeout:
            fmt.Println("Timeout reached, stopping")
            return
        }
    }
}

func main() {
    fmt.Println("=== Select with Timeout ===")
    selectWithTimeout()
    
    fmt.Println("\n=== Multiple Timeouts ===")
    multipleTimeouts()
    
    fmt.Println("\n=== Periodic with Timeout ===")
    periodicWithTimeout()
}
```

## Advanced Select Patterns

### Fan-in Pattern

```go
// fan_in_select.go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

// Fan-in: multiple channels into one
func fanIn(input1, input2 <-chan string) <-chan string {
    output := make(chan string)
    
    go func() {
        defer close(output)
        
        for {
            select {
            case msg, ok := <-input1:
                if !ok {
                    input1 = nil // Disable this case
                } else {
                    output <- msg
                }
            case msg, ok := <-input2:
                if !ok {
                    input2 = nil // Disable this case
                } else {
                    output <- msg
                }
            }
            
            // Exit when both channels are closed
            if input1 == nil && input2 == nil {
                return
            }
        }
    }()
    
    return output
}

// Producer function
func producer(name string, count int) <-chan string {
    ch := make(chan string)
    
    go func() {
        defer close(ch)
        
        for i := 0; i < count; i++ {
            // Random delay
            time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
            ch <- fmt.Sprintf("%s: message %d", name, i+1)
        }
    }()
    
    return ch
}

func demonstrateFanIn() {
    // Create two producers
    ch1 := producer("Producer1", 3)
    ch2 := producer("Producer2", 3)
    
    // Fan-in the channels
    merged := fanIn(ch1, ch2)
    
    // Consume merged messages
    for msg := range merged {
        fmt.Println("Received:", msg)
    }
    
    fmt.Println("All messages received")
}

func main() {
    fmt.Println("=== Fan-in Pattern ===")
    demonstrateFanIn()
}
```

### Select for Graceful Shutdown

```go
// graceful_shutdown.go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Worker struct {
    id   int
    jobs <-chan Job
    quit <-chan struct{}
}

type Job struct {
    ID   int
    Data string
}

func (w *Worker) Start() {
    fmt.Printf("Worker %d started\n", w.id)
    
    for {
        select {
        case job, ok := <-w.jobs:
            if !ok {
                fmt.Printf("Worker %d: jobs channel closed\n", w.id)
                return
            }
            
            // Process job
            fmt.Printf("Worker %d processing job %d: %s\n", w.id, job.ID, job.Data)
            time.Sleep(500 * time.Millisecond) // Simulate work
            
        case <-w.quit:
            fmt.Printf("Worker %d received quit signal\n", w.id)
            return
        }
    }
}

func gracefulShutdownExample() {
    // Create channels
    jobs := make(chan Job, 10)
    quit := make(chan struct{})
    
    // Start workers
    numWorkers := 3
    for i := 0; i < numWorkers; i++ {
        worker := &Worker{
            id:   i + 1,
            jobs: jobs,
            quit: quit,
        }
        go worker.Start()
    }
    
    // Send some jobs
    go func() {
        for i := 0; i < 10; i++ {
            jobs <- Job{
                ID:   i + 1,
                Data: fmt.Sprintf("Task %d", i+1),
            }
            time.Sleep(200 * time.Millisecond)
        }
        close(jobs)
    }()
    
    // Wait for interrupt signal
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    
    select {
    case sig := <-sigChan:
        fmt.Printf("\nReceived signal: %v\n", sig)
        fmt.Println("Initiating graceful shutdown...")
        
        // Signal all workers to quit
        close(quit)
        
        // Give workers time to finish
        time.Sleep(1 * time.Second)
        fmt.Println("Shutdown complete")
        
    case <-time.After(10 * time.Second):
        fmt.Println("\nTimeout reached, shutting down")
        close(quit)
    }
}

func main() {
    fmt.Println("=== Graceful Shutdown Example ===")
    fmt.Println("Press Ctrl+C to trigger graceful shutdown")
    gracefulShutdownExample()
}
```

### Load Balancer with Select

```go
// load_balancer.go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Request struct {
    ID       int
    Data     string
    Response chan string
}

type Server struct {
    ID       int
    Requests chan Request
    Load     int
}

func (s *Server) Start() {
    for req := range s.Requests {
        // Simulate processing time
        processingTime := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(processingTime)
        
        response := fmt.Sprintf("Server %d processed request %d", s.ID, req.ID)
        req.Response <- response
        
        fmt.Printf("Server %d completed request %d (took %v)\n", 
            s.ID, req.ID, processingTime)
    }
}

// Load balancer using select
func loadBalancer(servers []*Server, requests <-chan Request) {
    for req := range requests {
        // Try to send to any available server
        sent := false
        
        // Create cases for select
        cases := make([]chan Request, len(servers))
        for i, server := range servers {
            cases[i] = server.Requests
        }
        
        // Try to send to first available server
        for !sent {
            switch len(servers) {
            case 1:
                select {
                case cases[0] <- req:
                    sent = true
                    fmt.Printf("Request %d sent to server %d\n", req.ID, servers[0].ID)
                default:
                    time.Sleep(10 * time.Millisecond)
                }
            case 2:
                select {
                case cases[0] <- req:
                    sent = true
                    fmt.Printf("Request %d sent to server %d\n", req.ID, servers[0].ID)
                case cases[1] <- req:
                    sent = true
                    fmt.Printf("Request %d sent to server %d\n", req.ID, servers[1].ID)
                default:
                    time.Sleep(10 * time.Millisecond)
                }
            case 3:
                select {
                case cases[0] <- req:
                    sent = true
                    fmt.Printf("Request %d sent to server %d\n", req.ID, servers[0].ID)
                case cases[1] <- req:
                    sent = true
                    fmt.Printf("Request %d sent to server %d\n", req.ID, servers[1].ID)
                case cases[2] <- req:
                    sent = true
                    fmt.Printf("Request %d sent to server %d\n", req.ID, servers[2].ID)
                default:
                    time.Sleep(10 * time.Millisecond)
                }
            }
        }
    }
}

func demonstrateLoadBalancer() {
    // Create servers
    numServers := 3
    servers := make([]*Server, numServers)
    
    for i := 0; i < numServers; i++ {
        servers[i] = &Server{
            ID:       i + 1,
            Requests: make(chan Request, 2), // Small buffer
        }
        go servers[i].Start()
    }
    
    // Create request channel
    requests := make(chan Request, 10)
    
    // Start load balancer
    go loadBalancer(servers, requests)
    
    // Send requests
    numRequests := 10
    responses := make([]chan string, numRequests)
    
    for i := 0; i < numRequests; i++ {
        responses[i] = make(chan string, 1)
        
        requests <- Request{
            ID:       i + 1,
            Data:     fmt.Sprintf("Request data %d", i+1),
            Response: responses[i],
        }
    }
    
    // Collect responses
    for i := 0; i < numRequests; i++ {
        response := <-responses[i]
        fmt.Printf("Response %d: %s\n", i+1, response)
    }
    
    // Cleanup
    close(requests)
    for _, server := range servers {
        close(server.Requests)
    }
}

func main() {
    fmt.Println("=== Load Balancer Example ===")
    demonstrateLoadBalancer()
}
```

## Select Best Practices

### 1. Handling Closed Channels

```go
// handling_closed_channels.go
package main

import (
    "fmt"
    "time"
)

func handleClosedChannels() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // Close ch1 after a delay
    go func() {
        time.Sleep(1 * time.Second)
        close(ch1)
    }()
    
    // Send to ch2 after a delay
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message from ch2"
        close(ch2)
    }()
    
    // Handle both channels until both are closed
    for {
        select {
        case msg, ok := <-ch1:
            if !ok {
                fmt.Println("ch1 is closed")
                ch1 = nil // Disable this case
            } else {
                fmt.Println("From ch1:", msg)
            }
        case msg, ok := <-ch2:
            if !ok {
                fmt.Println("ch2 is closed")
                ch2 = nil // Disable this case
            } else {
                fmt.Println("From ch2:", msg)
            }
        }
        
        // Exit when both channels are nil (closed)
        if ch1 == nil && ch2 == nil {
            break
        }
    }
    
    fmt.Println("All channels processed")
}

func main() {
    fmt.Println("=== Handling Closed Channels ===")
    handleClosedChannels()
}
```

### 2. Avoiding Goroutine Leaks

```go
// avoiding_leaks.go
package main

import (
    "context"
    "fmt"
    "time"
)

// Good: Using context for cancellation
func workerWithContext(ctx context.Context, id int, jobs <-chan int, results chan<- int) {
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                fmt.Printf("Worker %d: jobs channel closed\n", id)
                return
            }
            
            // Simulate work with context checking
            select {
            case <-time.After(500 * time.Millisecond):
                results <- job * 2
                fmt.Printf("Worker %d completed job %d\n", id, job)
            case <-ctx.Done():
                fmt.Printf("Worker %d cancelled during job %d\n", id, job)
                return
            }
            
        case <-ctx.Done():
            fmt.Printf("Worker %d cancelled\n", id)
            return
        }
    }
}

func demonstrateContextCancellation() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    
    jobs := make(chan int, 10)
    results := make(chan int, 10)
    
    // Start workers
    numWorkers := 3
    for i := 0; i < numWorkers; i++ {
        go workerWithContext(ctx, i+1, jobs, results)
    }
    
    // Send jobs
    go func() {
        defer close(jobs)
        for i := 1; i <= 10; i++ {
            select {
            case jobs <- i:
                fmt.Printf("Sent job %d\n", i)
            case <-ctx.Done():
                fmt.Println("Job sender cancelled")
                return
            }
            time.Sleep(200 * time.Millisecond)
        }
    }()
    
    // Collect results
    go func() {
        defer close(results)
        for {
            select {
            case result, ok := <-results:
                if !ok {
                    return
                }
                fmt.Printf("Result: %d\n", result)
            case <-ctx.Done():
                fmt.Println("Result collector cancelled")
                return
            }
        }
    }()
    
    // Wait for context to be done
    <-ctx.Done()
    fmt.Println("Context cancelled, cleaning up...")
    
    // Give goroutines time to clean up
    time.Sleep(500 * time.Millisecond)
}

func main() {
    fmt.Println("=== Context Cancellation Example ===")
    demonstrateContextCancellation()
}
```

## Common Pitfalls and Solutions

### 1. Select in Loops

```go
// select_pitfalls.go
package main

import (
    "fmt"
    "time"
)

// Bad: Infinite loop with select
func badSelectLoop() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch <- "message"
        close(ch)
    }()
    
    // This will loop forever after channel is closed
    // for {
    //     select {
    //     case msg := <-ch:
    //         fmt.Println("Received:", msg)
    //     }
    // }
}

// Good: Proper handling of closed channels
func goodSelectLoop() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch <- "message 1"
        time.Sleep(1 * time.Second)
        ch <- "message 2"
        close(ch)
    }()
    
    for {
        select {
        case msg, ok := <-ch:
            if !ok {
                fmt.Println("Channel closed, exiting loop")
                return
            }
            fmt.Println("Received:", msg)
        }
    }
}

// Alternative: Using range
func usingRange() {
    ch := make(chan string)
    
    go func() {
        defer close(ch)
        ch <- "message 1"
        time.Sleep(1 * time.Second)
        ch <- "message 2"
        time.Sleep(1 * time.Second)
        ch <- "message 3"
    }()
    
    // Range automatically handles channel closure
    for msg := range ch {
        fmt.Println("Received:", msg)
    }
    
    fmt.Println("Range loop completed")
}

func main() {
    fmt.Println("=== Good Select Loop ===")
    goodSelectLoop()
    
    fmt.Println("\n=== Using Range ===")
    usingRange()
}
```

## Exercises

### Exercise 1: Message Router
Implement a message router that:
- Receives messages from multiple input channels
- Routes messages to different output channels based on message type
- Handles timeouts and graceful shutdown

### Exercise 2: Rate-Limited Worker Pool
Create a worker pool that:
- Processes jobs from a queue
- Limits the rate of job processing
- Uses select for coordination
- Supports dynamic scaling

### Exercise 3: Circuit Breaker
Implement a circuit breaker pattern using select that:
- Monitors service health
- Opens circuit on failures
- Attempts recovery after timeout
- Provides fallback responses

### Exercise 4: Event Aggregator
Build an event aggregator that:
- Collects events from multiple sources
- Batches events by time or count
- Handles backpressure
- Supports multiple subscribers

## Key Takeaways

- `select` enables non-blocking communication on multiple channels
- Use `default` case for non-blocking operations
- Combine `select` with `time.After()` for timeouts
- Always handle closed channels properly in select loops
- Use context for cancellation and timeouts
- `select` chooses randomly when multiple cases are ready
- Nil channels are ignored in select statements
- `select` is essential for implementing complex concurrency patterns
- Always consider goroutine lifecycle and cleanup
- Use `select` for graceful shutdown patterns

## Next Steps

Next, we'll explore [Context Package](23-context-package.md) to learn about managing request-scoped data, cancellation signals, and timeouts across API boundaries.