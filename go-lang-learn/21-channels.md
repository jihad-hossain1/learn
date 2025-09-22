# 21. Channels

Channels are Go's way of implementing communication between goroutines. They provide a type-safe way to pass data between concurrent processes, following the principle "Don't communicate by sharing memory; share memory by communicating."

## Channel Basics

### Creating Channels

```go
package main

import "fmt"

func main() {
    // Create an unbuffered channel
    ch := make(chan int)
    
    // Create a buffered channel
    bufferedCh := make(chan string, 3)
    
    // Create a receive-only channel
    var recvOnly <-chan int = ch
    
    // Create a send-only channel
    var sendOnly chan<- int = ch
    
    fmt.Printf("Channel: %T\n", ch)
    fmt.Printf("Buffered channel: %T\n", bufferedCh)
    fmt.Printf("Receive-only: %T\n", recvOnly)
    fmt.Printf("Send-only: %T\n", sendOnly)
}
```

### Basic Channel Operations

```go
package main

import (
    "fmt"
    "time"
)

func sender(ch chan<- string) {
    messages := []string{"Hello", "World", "From", "Goroutine"}
    
    for _, msg := range messages {
        fmt.Printf("Sending: %s\n", msg)
        ch <- msg // Send to channel
        time.Sleep(500 * time.Millisecond)
    }
    
    close(ch) // Close the channel
}

func receiver(ch <-chan string) {
    for {
        msg, ok := <-ch // Receive from channel
        if !ok {
            fmt.Println("Channel closed")
            break
        }
        fmt.Printf("Received: %s\n", msg)
    }
}

func main() {
    ch := make(chan string)
    
    go sender(ch)
    go receiver(ch)
    
    // Wait for goroutines to complete
    time.Sleep(3 * time.Second)
}
```

### Range Over Channels

```go
package main

import (
    "fmt"
    "time"
)

func fibonacci(n int, ch chan<- int) {
    x, y := 0, 1
    for i := 0; i < n; i++ {
        ch <- x
        x, y = y, x+y
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
}

func main() {
    ch := make(chan int, 10)
    
    go fibonacci(10, ch)
    
    // Range automatically breaks when channel is closed
    for num := range ch {
        fmt.Printf("Fibonacci: %d\n", num)
    }
    
    fmt.Println("Done")
}
```

## Buffered vs Unbuffered Channels

### Unbuffered Channels (Synchronous)

```go
package main

import (
    "fmt"
    "time"
)

func unbufferedExample() {
    ch := make(chan string) // Unbuffered channel
    
    go func() {
        fmt.Println("Goroutine: About to send")
        ch <- "Hello" // Blocks until someone receives
        fmt.Println("Goroutine: Sent message")
    }()
    
    time.Sleep(1 * time.Second) // Simulate delay
    
    fmt.Println("Main: About to receive")
    msg := <-ch // Blocks until someone sends
    fmt.Printf("Main: Received '%s'\n", msg)
}

func main() {
    fmt.Println("=== Unbuffered Channel Example ===")
    unbufferedExample()
}
```

### Buffered Channels (Asynchronous)

```go
package main

import (
    "fmt"
    "time"
)

func bufferedExample() {
    ch := make(chan string, 2) // Buffered channel with capacity 2
    
    go func() {
        fmt.Println("Goroutine: Sending first message")
        ch <- "First"  // Doesn't block (buffer has space)
        fmt.Println("Goroutine: Sending second message")
        ch <- "Second" // Doesn't block (buffer has space)
        fmt.Println("Goroutine: Sending third message")
        ch <- "Third"  // Blocks (buffer is full)
        fmt.Println("Goroutine: Sent third message")
    }()
    
    time.Sleep(1 * time.Second)
    
    fmt.Printf("Main: Received '%s'\n", <-ch)
    fmt.Printf("Main: Received '%s'\n", <-ch)
    
    time.Sleep(500 * time.Millisecond)
    
    fmt.Printf("Main: Received '%s'\n", <-ch)
}

func main() {
    fmt.Println("=== Buffered Channel Example ===")
    bufferedExample()
}
```

### Channel Capacity and Length

```go
package main

import "fmt"

func main() {
    ch := make(chan int, 5)
    
    fmt.Printf("Initial - Capacity: %d, Length: %d\n", cap(ch), len(ch))
    
    // Send some values
    ch <- 1
    ch <- 2
    ch <- 3
    
    fmt.Printf("After sending 3 - Capacity: %d, Length: %d\n", cap(ch), len(ch))
    
    // Receive one value
    <-ch
    
    fmt.Printf("After receiving 1 - Capacity: %d, Length: %d\n", cap(ch), len(ch))
    
    // Fill the buffer
    ch <- 4
    ch <- 5
    ch <- 6
    
    fmt.Printf("Buffer full - Capacity: %d, Length: %d\n", cap(ch), len(ch))
}
```

## Select Statement

### Basic Select

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Message from channel 1"
    }()
    
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message from channel 2"
    }()
    
    // Select waits for one of the channels to be ready
    select {
    case msg1 := <-ch1:
        fmt.Printf("Received: %s\n", msg1)
    case msg2 := <-ch2:
        fmt.Printf("Received: %s\n", msg2)
    }
    
    fmt.Println("Done")
}
```

### Select with Default Case

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch := make(chan string)
    
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "Hello after 2 seconds"
    }()
    
    for i := 0; i < 5; i++ {
        select {
        case msg := <-ch:
            fmt.Printf("Received: %s\n", msg)
            return
        default:
            fmt.Printf("No message yet (attempt %d)\n", i+1)
            time.Sleep(500 * time.Millisecond)
        }
    }
}
```

### Select with Timeout

```go
package main

import (
    "fmt"
    "time"
)

func fetchData(ch chan<- string, delay time.Duration) {
    time.Sleep(delay)
    ch <- "Data fetched"
}

func main() {
    ch := make(chan string)
    
    // Start a goroutine that takes 3 seconds
    go fetchData(ch, 3*time.Second)
    
    select {
    case data := <-ch:
        fmt.Printf("Success: %s\n", data)
    case <-time.After(2 * time.Second):
        fmt.Println("Timeout: Operation took too long")
    }
}
```

### Multiple Channel Operations

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func producer(name string, ch chan<- string) {
    for i := 0; i < 3; i++ {
        delay := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(delay)
        
        message := fmt.Sprintf("%s: message %d", name, i+1)
        ch <- message
    }
    close(ch)
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    ch1 := make(chan string)
    ch2 := make(chan string)
    ch3 := make(chan string)
    
    go producer("Producer1", ch1)
    go producer("Producer2", ch2)
    go producer("Producer3", ch3)
    
    openChannels := 3
    
    for openChannels > 0 {
        select {
        case msg, ok := <-ch1:
            if ok {
                fmt.Printf("From ch1: %s\n", msg)
            } else {
                fmt.Println("ch1 closed")
                ch1 = nil // Disable this case
                openChannels--
            }
        case msg, ok := <-ch2:
            if ok {
                fmt.Printf("From ch2: %s\n", msg)
            } else {
                fmt.Println("ch2 closed")
                ch2 = nil // Disable this case
                openChannels--
            }
        case msg, ok := <-ch3:
            if ok {
                fmt.Printf("From ch3: %s\n", msg)
            } else {
                fmt.Println("ch3 closed")
                ch3 = nil // Disable this case
                openChannels--
            }
        }
    }
    
    fmt.Println("All channels closed")
}
```

## Channel Patterns

### Worker Pool with Channels

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
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for job := range jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        
        // Simulate work
        processingTime := time.Duration(rand.Intn(1000)) * time.Millisecond
        time.Sleep(processingTime)
        
        result := Result{
            JobID:  job.ID,
            Output: fmt.Sprintf("Processed %s by worker %d", job.Data, id),
        }
        
        results <- result
    }
    
    fmt.Printf("Worker %d finished\n", id)
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    const numWorkers = 3
    const numJobs = 10
    
    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)
    
    var wg sync.WaitGroup
    
    // Start workers
    for w := 1; w <= numWorkers; w++ {
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }
    
    // Send jobs
    for j := 1; j <= numJobs; j++ {
        job := Job{
            ID:   j,
            Data: fmt.Sprintf("task-%d", j),
        }
        jobs <- job
    }
    close(jobs)
    
    // Close results channel when all workers are done
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    fmt.Println("\nResults:")
    for result := range results {
        fmt.Printf("Job %d: %s\n", result.JobID, result.Output)
    }
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
            time.Sleep(100 * time.Millisecond)
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
            time.Sleep(200 * time.Millisecond)
        }
    }()
    return out
}

// Stage 3: Filter even numbers
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

// Stage 4: Format as string
func format(in <-chan int) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for n := range in {
            out <- fmt.Sprintf("Number: %d", n)
        }
    }()
    return out
}

func main() {
    // Create pipeline
    numbers := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squared := square(numbers)
    filtered := filterEven(squared)
    formatted := format(filtered)
    
    // Consume results
    fmt.Println("Pipeline results:")
    for result := range formatted {
        fmt.Println(result)
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

// Fan-out: distribute work to multiple workers
func fanOut(input <-chan int, workers int) []<-chan int {
    outputs := make([]<-chan int, workers)
    
    for i := 0; i < workers; i++ {
        output := make(chan int)
        outputs[i] = output
        
        go func(workerID int, out chan<- int) {
            defer close(out)
            for n := range input {
                // Simulate variable processing time
                processingTime := time.Duration(rand.Intn(500)) * time.Millisecond
                time.Sleep(processingTime)
                
                result := n * n
                fmt.Printf("Worker %d: %d -> %d\n", workerID, n, result)
                out <- result
            }
            fmt.Printf("Worker %d finished\n", workerID)
        }(i+1, output)
    }
    
    return outputs
}

// Fan-in: merge multiple channels into one
func fanIn(inputs ...<-chan int) <-chan int {
    output := make(chan int)
    var wg sync.WaitGroup
    
    // Start a goroutine for each input channel
    for i, input := range inputs {
        wg.Add(1)
        go func(workerID int, in <-chan int) {
            defer wg.Done()
            for value := range in {
                fmt.Printf("Fan-in from worker %d: %d\n", workerID, value)
                output <- value
            }
        }(i+1, input)
    }
    
    // Close output channel when all inputs are done
    go func() {
        wg.Wait()
        close(output)
        fmt.Println("Fan-in completed")
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
        for i := 1; i <= 8; i++ {
            fmt.Printf("Sending: %d\n", i)
            input <- i
            time.Sleep(200 * time.Millisecond)
        }
        fmt.Println("Input generation completed")
    }()
    
    // Fan-out to 3 workers
    workerOutputs := fanOut(input, 3)
    
    // Fan-in results
    results := fanIn(workerOutputs...)
    
    // Collect all results
    fmt.Println("\nFinal results:")
    var finalResults []int
    for result := range results {
        finalResults = append(finalResults, result)
    }
    
    fmt.Printf("Collected %d results: %v\n", len(finalResults), finalResults)
}
```

### Request-Response Pattern

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

type Request struct {
    ID       int
    Data     string
    Response chan Response
}

type Response struct {
    ID     int
    Result string
    Error  error
}

type Server struct {
    requests chan Request
}

func NewServer() *Server {
    s := &Server{
        requests: make(chan Request),
    }
    
    // Start server goroutine
    go s.run()
    
    return s
}

func (s *Server) run() {
    for req := range s.requests {
        go s.handleRequest(req)
    }
}

func (s *Server) handleRequest(req Request) {
    fmt.Printf("Processing request %d: %s\n", req.ID, req.Data)
    
    // Simulate processing time
    processingTime := time.Duration(rand.Intn(1000)) * time.Millisecond
    time.Sleep(processingTime)
    
    // Simulate random errors
    var resp Response
    if rand.Float32() < 0.2 { // 20% chance of error
        resp = Response{
            ID:    req.ID,
            Error: fmt.Errorf("processing failed for request %d", req.ID),
        }
    } else {
        resp = Response{
            ID:     req.ID,
            Result: fmt.Sprintf("Processed: %s", req.Data),
        }
    }
    
    // Send response back
    req.Response <- resp
    close(req.Response)
}

func (s *Server) ProcessRequest(id int, data string) Response {
    responseCh := make(chan Response, 1)
    
    req := Request{
        ID:       id,
        Data:     data,
        Response: responseCh,
    }
    
    s.requests <- req
    
    // Wait for response with timeout
    select {
    case resp := <-responseCh:
        return resp
    case <-time.After(2 * time.Second):
        return Response{
            ID:    id,
            Error: fmt.Errorf("request %d timed out", id),
        }
    }
}

func (s *Server) Close() {
    close(s.requests)
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    server := NewServer()
    defer server.Close()
    
    // Make multiple concurrent requests
    var responses []Response
    responseCh := make(chan Response, 5)
    
    for i := 1; i <= 5; i++ {
        go func(id int) {
            resp := server.ProcessRequest(id, fmt.Sprintf("data-%d", id))
            responseCh <- resp
        }(i)
    }
    
    // Collect responses
    for i := 0; i < 5; i++ {
        resp := <-responseCh
        responses = append(responses, resp)
    }
    
    // Print results
    fmt.Println("\nResults:")
    for _, resp := range responses {
        if resp.Error != nil {
            fmt.Printf("Request %d failed: %v\n", resp.ID, resp.Error)
        } else {
            fmt.Printf("Request %d succeeded: %s\n", resp.ID, resp.Result)
        }
    }
}
```

## Channel Synchronization Patterns

### Done Channel Pattern

```go
package main

import (
    "fmt"
    "time"
)

func worker(id int, done <-chan bool) {
    for {
        select {
        case <-done:
            fmt.Printf("Worker %d: Received done signal\n", id)
            return
        default:
            fmt.Printf("Worker %d: Working...\n", id)
            time.Sleep(500 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan bool)
    
    // Start workers
    for i := 1; i <= 3; i++ {
        go worker(i, done)
    }
    
    // Let workers run for 2 seconds
    time.Sleep(2 * time.Second)
    
    // Signal all workers to stop
    fmt.Println("Signaling workers to stop...")
    close(done) // Closing broadcasts to all receivers
    
    // Give time for cleanup
    time.Sleep(1 * time.Second)
    fmt.Println("Main function ending")
}
```

### Semaphore Pattern

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Semaphore chan struct{}

func NewSemaphore(capacity int) Semaphore {
    return make(Semaphore, capacity)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}

func worker(id int, sem Semaphore, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d: Waiting for semaphore...\n", id)
    sem.Acquire()
    defer sem.Release()
    
    fmt.Printf("Worker %d: Acquired semaphore, working...\n", id)
    time.Sleep(2 * time.Second) // Simulate work
    fmt.Printf("Worker %d: Work completed\n", id)
}

func main() {
    // Create semaphore with capacity 2 (max 2 concurrent workers)
    sem := NewSemaphore(2)
    var wg sync.WaitGroup
    
    // Start 5 workers
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go worker(i, sem, &wg)
    }
    
    wg.Wait()
    fmt.Println("All workers completed")
}
```

### Barrier Pattern

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Barrier struct {
    n       int
    count   int
    mutex   sync.Mutex
    cond    *sync.Cond
    waiting chan struct{}
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{
        n:       n,
        waiting: make(chan struct{}),
    }
    b.cond = sync.NewCond(&b.mutex)
    return b
}

func (b *Barrier) Wait() {
    b.mutex.Lock()
    defer b.mutex.Unlock()
    
    b.count++
    if b.count == b.n {
        // Last goroutine to arrive
        close(b.waiting)
        b.cond.Broadcast()
    } else {
        // Wait for all others
        for b.count < b.n {
            b.cond.Wait()
        }
    }
}

func worker(id int, barrier *Barrier) {
    // Phase 1: Do some work
    workTime := time.Duration(rand.Intn(2000)) * time.Millisecond
    fmt.Printf("Worker %d: Starting phase 1 (will take %v)\n", id, workTime)
    time.Sleep(workTime)
    fmt.Printf("Worker %d: Completed phase 1\n", id)
    
    // Wait at barrier
    fmt.Printf("Worker %d: Waiting at barrier\n", id)
    barrier.Wait()
    
    // Phase 2: Continue after all workers reach barrier
    fmt.Printf("Worker %d: Starting phase 2\n", id)
    time.Sleep(500 * time.Millisecond)
    fmt.Printf("Worker %d: Completed phase 2\n", id)
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    const numWorkers = 4
    barrier := NewBarrier(numWorkers)
    
    var wg sync.WaitGroup
    
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(id, barrier)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All workers completed both phases")
}
```

## Error Handling with Channels

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

type Result struct {
    TaskID int
    Value  string
    Error  error
}

func processTask(task Task) Result {
    // Simulate random processing time
    processingTime := time.Duration(rand.Intn(1000)) * time.Millisecond
    time.Sleep(processingTime)
    
    // Simulate random errors
    if rand.Float32() < 0.3 {
        return Result{
            TaskID: task.ID,
            Error:  errors.New("random processing error"),
        }
    }
    
    return Result{
        TaskID: task.ID,
        Value:  fmt.Sprintf("Processed: %s", task.Data),
    }
}

func worker(tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for task := range tasks {
        result := processTask(task)
        results <- result
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    tasks := make(chan Task, 10)
    results := make(chan Result, 10)
    
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
    
    // Collect results and separate successes from errors
    var successes []Result
    var errors []Result
    
    for result := range results {
        if result.Error != nil {
            errors = append(errors, result)
        } else {
            successes = append(successes, result)
        }
    }
    
    // Report results
    fmt.Printf("\nSuccessful tasks (%d):\n", len(successes))
    for _, result := range successes {
        fmt.Printf("  Task %d: %s\n", result.TaskID, result.Value)
    }
    
    fmt.Printf("\nFailed tasks (%d):\n", len(errors))
    for _, result := range errors {
        fmt.Printf("  Task %d: %v\n", result.TaskID, result.Error)
    }
}
```

## Channel Best Practices

### 1. Channel Ownership

```go
// Good: Clear ownership - sender closes the channel
func producer(ch chan<- int) {
    defer close(ch) // Producer owns and closes the channel
    
    for i := 0; i < 5; i++ {
        ch <- i
    }
}

func consumer(ch <-chan int) {
    for value := range ch { // Consumer just reads
        fmt.Println(value)
    }
}
```

### 2. Avoid Channel Leaks

```go
// Bad: Potential goroutine leak
func badExample() <-chan int {
    ch := make(chan int)
    go func() {
        ch <- 42 // If no one reads, this goroutine leaks
    }()
    return ch
}

// Good: Use buffered channel or ensure consumption
func goodExample() <-chan int {
    ch := make(chan int, 1) // Buffered
    go func() {
        ch <- 42
        close(ch)
    }()
    return ch
}
```

### 3. Use Directional Channels

```go
// Good: Use directional channels for clarity
func sender(ch chan<- string) { // Send-only
    ch <- "message"
}

func receiver(ch <-chan string) { // Receive-only
    msg := <-ch
    fmt.Println(msg)
}
```

### 4. Handle Channel Closure

```go
// Good: Always check if channel is closed
func safeReceive(ch <-chan int) {
    for {
        value, ok := <-ch
        if !ok {
            fmt.Println("Channel closed")
            break
        }
        fmt.Printf("Received: %d\n", value)
    }
}

// Better: Use range
func betterReceive(ch <-chan int) {
    for value := range ch {
        fmt.Printf("Received: %d\n", value)
    }
    fmt.Println("Channel closed")
}
```

## Exercises

### Exercise 1: Rate Limiter
Implement a rate limiter using channels that:
- Allows a maximum number of operations per time window
- Blocks when the limit is exceeded
- Supports different rate limiting strategies

### Exercise 2: Message Broker
Build a simple message broker that:
- Supports multiple topics
- Allows publishers and subscribers
- Handles subscriber disconnections gracefully
- Implements message persistence

### Exercise 3: Parallel Web Scraper
Create a web scraper that:
- Fetches multiple URLs concurrently
- Limits concurrent requests
- Handles errors and retries
- Aggregates results from all workers

## Key Takeaways

- Channels are Go's primary mechanism for goroutine communication
- Unbuffered channels provide synchronous communication
- Buffered channels allow asynchronous communication up to capacity
- Use `select` for non-blocking operations and timeouts
- Always close channels when done sending (by the sender)
- Use directional channels for better API design
- Range over channels for clean consumption patterns
- Implement proper error handling with channels
- Avoid channel leaks by ensuring proper cleanup
- Use channels to implement various synchronization patterns

## Next Steps

Next, we'll explore [Context](22-context.md) to learn about Go's powerful mechanism for handling cancellation, timeouts, and request-scoped values.