# 16. Channels

Channels are the pipes that connect concurrent goroutines. They allow you to pass values between goroutines, enabling safe communication and synchronization. Channels are a core feature of Go's concurrency model and follow the principle: "Don't communicate by sharing memory; share memory by communicating."

## Channel Basics

### Creating and Using Channels

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== Basic Channel Operations ===")
    
    // Creating channels
    var ch1 chan int                    // nil channel
    ch2 := make(chan int)              // unbuffered channel
    ch3 := make(chan int, 5)           // buffered channel with capacity 5
    ch4 := make(chan string, 2)        // buffered string channel
    
    fmt.Printf("ch1 (nil): %v\n", ch1)
    fmt.Printf("ch2 (unbuffered): %v\n", ch2)
    fmt.Printf("ch3 (buffered): %v, capacity: %d\n", ch3, cap(ch3))
    fmt.Printf("ch4 (string): %v, capacity: %d\n", ch4, cap(ch4))
    
    // Basic send and receive operations
    fmt.Println("\n=== Send and Receive ===")
    
    // Using buffered channel for demonstration
    numbers := make(chan int, 3)
    
    // Sending values
    numbers <- 10
    numbers <- 20
    numbers <- 30
    
    fmt.Printf("Channel length: %d, capacity: %d\n", len(numbers), cap(numbers))
    
    // Receiving values
    first := <-numbers
    second := <-numbers
    third := <-numbers
    
    fmt.Printf("Received: %d, %d, %d\n", first, second, third)
    fmt.Printf("Channel length after receiving: %d\n", len(numbers))
    
    // Channel with goroutines
    fmt.Println("\n=== Channel with Goroutines ===")
    
    messages := make(chan string)
    
    // Sender goroutine
    go func() {
        messages <- "Hello"
        messages <- "World"
        messages <- "from"
        messages <- "Goroutine"
        close(messages) // Close channel when done sending
    }()
    
    // Receiver (main goroutine)
    for msg := range messages {
        fmt.Printf("Received: %s\n", msg)
        time.Sleep(200 * time.Millisecond)
    }
    
    fmt.Println("Channel closed and all messages received")
}
```

### Channel Directions

```go
package main

import (
    "fmt"
    "time"
)

// Send-only channel parameter
func sender(ch chan<- string, messages []string) {
    fmt.Println("Sender started")
    for i, msg := range messages {
        fmt.Printf("Sending: %s\n", msg)
        ch <- fmt.Sprintf("%d: %s", i+1, msg)
        time.Sleep(100 * time.Millisecond)
    }
    close(ch)
    fmt.Println("Sender finished")
}

// Receive-only channel parameter
func receiver(ch <-chan string, name string) {
    fmt.Printf("%s receiver started\n", name)
    for msg := range ch {
        fmt.Printf("%s received: %s\n", name, msg)
        time.Sleep(150 * time.Millisecond)
    }
    fmt.Printf("%s receiver finished\n", name)
}

// Bidirectional channel that gets converted
func processor(input <-chan int, output chan<- int) {
    fmt.Println("Processor started")
    for num := range input {
        processed := num * num
        fmt.Printf("Processing %d -> %d\n", num, processed)
        output <- processed
        time.Sleep(100 * time.Millisecond)
    }
    close(output)
    fmt.Println("Processor finished")
}

// Function that demonstrates channel direction conversion
func channelDirections() {
    fmt.Println("\n=== Channel Directions ===")
    
    // String channel example
    messages := make(chan string, 3)
    
    go sender(messages, []string{"Hello", "World", "Go", "Channels"})
    receiver(messages, "Main")
    
    // Processing pipeline example
    fmt.Println("\n=== Processing Pipeline ===")
    
    input := make(chan int, 5)
    output := make(chan int, 5)
    
    // Start processor
    go processor(input, output)
    
    // Send input data
    go func() {
        for i := 1; i <= 5; i++ {
            input <- i
        }
        close(input)
    }()
    
    // Receive processed data
    for result := range output {
        fmt.Printf("Final result: %d\n", result)
    }
}

func main() {
    channelDirections()
}
```

## Buffered vs Unbuffered Channels

### Understanding the Differences

```go
package main

import (
    "fmt"
    "time"
)

func demonstrateUnbuffered() {
    fmt.Println("=== Unbuffered Channels ===")
    
    ch := make(chan string)
    
    // This would block forever if not in a goroutine
    go func() {
        fmt.Println("Goroutine: About to send")
        ch <- "Hello from unbuffered channel"
        fmt.Println("Goroutine: Sent successfully")
    }()
    
    // Small delay to show the goroutine starts
    time.Sleep(100 * time.Millisecond)
    
    fmt.Println("Main: About to receive")
    msg := <-ch
    fmt.Printf("Main: Received: %s\n", msg)
    
    // Demonstrate synchronous nature
    fmt.Println("\n--- Synchronous Communication ---")
    
    sync := make(chan bool)
    
    go func() {
        fmt.Println("Worker: Starting work")
        time.Sleep(500 * time.Millisecond)
        fmt.Println("Worker: Work completed")
        sync <- true // Signal completion
    }()
    
    fmt.Println("Main: Waiting for worker")
    <-sync // Wait for signal
    fmt.Println("Main: Worker finished, continuing")
}

func demonstrateBuffered() {
    fmt.Println("\n=== Buffered Channels ===")
    
    ch := make(chan string, 3)
    
    // Can send without blocking (up to capacity)
    fmt.Println("Sending to buffered channel...")
    ch <- "Message 1"
    fmt.Printf("Sent 1, length: %d\n", len(ch))
    
    ch <- "Message 2"
    fmt.Printf("Sent 2, length: %d\n", len(ch))
    
    ch <- "Message 3"
    fmt.Printf("Sent 3, length: %d\n", len(ch))
    
    // Channel is now full
    fmt.Printf("Channel full: length=%d, capacity=%d\n", len(ch), cap(ch))
    
    // Receive messages
    for i := 0; i < 3; i++ {
        msg := <-ch
        fmt.Printf("Received: %s, remaining: %d\n", msg, len(ch))
    }
    
    // Demonstrate non-blocking behavior
    fmt.Println("\n--- Non-blocking Sends ---")
    
    buffer := make(chan int, 2)
    
    // These won't block
    buffer <- 1
    buffer <- 2
    
    fmt.Println("Sent 2 values without blocking")
    
    // This would block, so we do it in a goroutine
    go func() {
        fmt.Println("Goroutine: Trying to send third value...")
        buffer <- 3
        fmt.Println("Goroutine: Third value sent")
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    // Receive to make space
    val := <-buffer
    fmt.Printf("Received: %d, goroutine can now send\n", val)
    
    time.Sleep(100 * time.Millisecond)
    
    // Receive remaining
    fmt.Printf("Received: %d\n", <-buffer)
    fmt.Printf("Received: %d\n", <-buffer)
}

func comparePerformance() {
    fmt.Println("\n=== Performance Comparison ===")
    
    const numMessages = 1000
    
    // Unbuffered channel test
    start := time.Now()
    unbuffered := make(chan int)
    
    go func() {
        for i := 0; i < numMessages; i++ {
            unbuffered <- i
        }
        close(unbuffered)
    }()
    
    for range unbuffered {
        // Just consume
    }
    
    unbufferedTime := time.Since(start)
    
    // Buffered channel test
    start = time.Now()
    buffered := make(chan int, 100)
    
    go func() {
        for i := 0; i < numMessages; i++ {
            buffered <- i
        }
        close(buffered)
    }()
    
    for range buffered {
        // Just consume
    }
    
    bufferedTime := time.Since(start)
    
    fmt.Printf("Unbuffered channel time: %v\n", unbufferedTime)
    fmt.Printf("Buffered channel time: %v\n", bufferedTime)
    fmt.Printf("Buffered is %.2fx faster\n", float64(unbufferedTime)/float64(bufferedTime))
}

func main() {
    demonstrateUnbuffered()
    demonstrateBuffered()
    comparePerformance()
}
```

## Channel Operations

### Select Statement

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func basicSelect() {
    fmt.Println("=== Basic Select Statement ===")
    
    ch1 := make(chan string)
    ch2 := make(chan string)
    
    // Goroutine 1
    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "Message from channel 1"
    }()
    
    // Goroutine 2
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "Message from channel 2"
    }()
    
    // Select waits for the first available channel
    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Printf("Received from ch1: %s\n", msg1)
        case msg2 := <-ch2:
            fmt.Printf("Received from ch2: %s\n", msg2)
        }
    }
}

func selectWithDefault() {
    fmt.Println("\n=== Select with Default Case ===")
    
    ch := make(chan string, 1)
    
    // Non-blocking send
    select {
    case ch <- "Hello":
        fmt.Println("Sent 'Hello' to channel")
    default:
        fmt.Println("Channel is full, couldn't send")
    }
    
    // Non-blocking receive
    select {
    case msg := <-ch:
        fmt.Printf("Received: %s\n", msg)
    default:
        fmt.Println("No message available")
    }
    
    // Another non-blocking receive (channel is now empty)
    select {
    case msg := <-ch:
        fmt.Printf("Received: %s\n", msg)
    default:
        fmt.Println("Channel is empty")
    }
}

func selectWithTimeout() {
    fmt.Println("\n=== Select with Timeout ===")
    
    ch := make(chan string)
    
    // Goroutine that sends after 2 seconds
    go func() {
        time.Sleep(2 * time.Second)
        ch <- "Delayed message"
    }()
    
    // Wait for message or timeout
    select {
    case msg := <-ch:
        fmt.Printf("Received: %s\n", msg)
    case <-time.After(1 * time.Second):
        fmt.Println("Timeout: No message received within 1 second")
    }
    
    // Wait for the delayed message
    select {
    case msg := <-ch:
        fmt.Printf("Finally received: %s\n", msg)
    case <-time.After(2 * time.Second):
        fmt.Println("Timeout again")
    }
}

func randomSelect() {
    fmt.Println("\n=== Random Select Behavior ===")
    
    ch1 := make(chan string, 1)
    ch2 := make(chan string, 1)
    
    // Fill both channels
    ch1 <- "Channel 1"
    ch2 <- "Channel 2"
    
    // Select randomly chooses when multiple cases are ready
    for i := 0; i < 10; i++ {
        // Refill channels
        select {
        case ch1 <- "Channel 1":
        default:
        }
        select {
        case ch2 <- "Channel 2":
        default:
        }
        
        // Random selection
        select {
        case msg := <-ch1:
            fmt.Printf("Iteration %d: %s\n", i+1, msg)
        case msg := <-ch2:
            fmt.Printf("Iteration %d: %s\n", i+1, msg)
        }
    }
}

func multiplexer() {
    fmt.Println("\n=== Channel Multiplexer ===")
    
    // Create multiple input channels
    input1 := make(chan int)
    input2 := make(chan int)
    input3 := make(chan int)
    output := make(chan int)
    
    // Multiplexer goroutine
    go func() {
        defer close(output)
        
        for {
            select {
            case val, ok := <-input1:
                if !ok {
                    input1 = nil
                } else {
                    output <- val * 10 // Transform from input1
                }
            case val, ok := <-input2:
                if !ok {
                    input2 = nil
                } else {
                    output <- val * 100 // Transform from input2
                }
            case val, ok := <-input3:
                if !ok {
                    input3 = nil
                } else {
                    output <- val * 1000 // Transform from input3
                }
            }
            
            // Exit when all inputs are closed
            if input1 == nil && input2 == nil && input3 == nil {
                break
            }
        }
    }()
    
    // Send data to different inputs
    go func() {
        for i := 1; i <= 3; i++ {
            input1 <- i
            time.Sleep(100 * time.Millisecond)
        }
        close(input1)
    }()
    
    go func() {
        for i := 1; i <= 2; i++ {
            input2 <- i
            time.Sleep(150 * time.Millisecond)
        }
        close(input2)
    }()
    
    go func() {
        for i := 1; i <= 4; i++ {
            input3 <- i
            time.Sleep(80 * time.Millisecond)
        }
        close(input3)
    }()
    
    // Receive multiplexed output
    for val := range output {
        fmt.Printf("Multiplexed output: %d\n", val)
    }
}

func main() {
    rand.Seed(time.Now().UnixNano())
    
    basicSelect()
    selectWithDefault()
    selectWithTimeout()
    randomSelect()
    multiplexer()
}
```

### Channel Closing and Range

```go
package main

import (
    "fmt"
    "time"
)

func basicClosing() {
    fmt.Println("=== Basic Channel Closing ===")
    
    ch := make(chan int, 3)
    
    // Send some values
    ch <- 1
    ch <- 2
    ch <- 3
    
    // Close the channel
    close(ch)
    
    // Can still receive from closed channel
    for {
        val, ok := <-ch
        if !ok {
            fmt.Println("Channel is closed and empty")
            break
        }
        fmt.Printf("Received: %d\n", val)
    }
    
    // Receiving from closed empty channel returns zero value
    val, ok := <-ch
    fmt.Printf("After close: value=%d, ok=%t\n", val, ok)
}

func rangeOverChannel() {
    fmt.Println("\n=== Range Over Channel ===")
    
    numbers := make(chan int)
    
    // Producer goroutine
    go func() {
        defer close(numbers) // Important: close when done
        
        for i := 1; i <= 5; i++ {
            fmt.Printf("Sending: %d\n", i)
            numbers <- i
            time.Sleep(200 * time.Millisecond)
        }
        fmt.Println("Producer finished")
    }()
    
    // Consumer using range
    fmt.Println("Consumer starting...")
    for num := range numbers {
        fmt.Printf("Received: %d\n", num)
    }
    fmt.Println("Consumer finished")
}

func multipleProducers() {
    fmt.Println("\n=== Multiple Producers ===")
    
    results := make(chan string, 10)
    
    // Producer 1
    go func() {
        for i := 1; i <= 3; i++ {
            results <- fmt.Sprintf("Producer1: %d", i)
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // Producer 2
    go func() {
        for i := 1; i <= 3; i++ {
            results <- fmt.Sprintf("Producer2: %d", i)
            time.Sleep(150 * time.Millisecond)
        }
    }()
    
    // Producer 3
    go func() {
        for i := 1; i <= 3; i++ {
            results <- fmt.Sprintf("Producer3: %d", i)
            time.Sleep(80 * time.Millisecond)
        }
    }()
    
    // Wait for all producers to finish
    go func() {
        time.Sleep(1 * time.Second)
        close(results)
    }()
    
    // Consume all results
    for result := range results {
        fmt.Printf("Result: %s\n", result)
    }
}

func gracefulShutdown() {
    fmt.Println("\n=== Graceful Shutdown Pattern ===")
    
    data := make(chan int)
    done := make(chan bool)
    
    // Worker goroutine
    go func() {
        defer func() {
            done <- true
        }()
        
        for {
            select {
            case val, ok := <-data:
                if !ok {
                    fmt.Println("Worker: Data channel closed, shutting down")
                    return
                }
                fmt.Printf("Worker: Processing %d\n", val)
                time.Sleep(100 * time.Millisecond)
            }
        }
    }()
    
    // Send some data
    for i := 1; i <= 5; i++ {
        data <- i
        time.Sleep(50 * time.Millisecond)
    }
    
    // Signal shutdown by closing channel
    fmt.Println("Main: Closing data channel")
    close(data)
    
    // Wait for worker to finish
    <-done
    fmt.Println("Main: Worker finished gracefully")
}

func channelOfChannels() {
    fmt.Println("\n=== Channel of Channels ===")
    
    // Channel that carries channels
    chOfCh := make(chan chan string, 3)
    
    // Create and send channels
    for i := 1; i <= 3; i++ {
        ch := make(chan string, 1)
        ch <- fmt.Sprintf("Message from channel %d", i)
        close(ch)
        chOfCh <- ch
    }
    close(chOfCh)
    
    // Receive and process channels
    for ch := range chOfCh {
        for msg := range ch {
            fmt.Printf("Received: %s\n", msg)
        }
    }
}

func main() {
    basicClosing()
    rangeOverChannel()
    multipleProducers()
    gracefulShutdown()
    channelOfChannels()
}
```

## Advanced Channel Patterns

### Channel Synchronization Patterns

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Barrier synchronization using channels
func barrierPattern() {
    fmt.Println("=== Barrier Synchronization ===")
    
    const numWorkers = 5
    barrier := make(chan bool, numWorkers)
    
    // Start workers
    for i := 1; i <= numWorkers; i++ {
        go func(id int) {
            // Simulate work
            workTime := time.Duration(id*100) * time.Millisecond
            fmt.Printf("Worker %d: Starting work (%v)\n", id, workTime)
            time.Sleep(workTime)
            
            fmt.Printf("Worker %d: Finished work, waiting at barrier\n", id)
            barrier <- true
            
            // Wait for all workers to reach barrier
            <-barrier
            
            fmt.Printf("Worker %d: Proceeding after barrier\n", id)
        }(i)
    }
    
    // Collect all workers at barrier
    for i := 0; i < numWorkers; i++ {
        <-barrier
    }
    
    fmt.Println("All workers reached barrier, releasing them")
    
    // Release all workers
    for i := 0; i < numWorkers; i++ {
        barrier <- true
    }
    
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Barrier synchronization complete")
}

// Semaphore pattern using channels
func semaphorePattern() {
    fmt.Println("\n=== Semaphore Pattern ===")
    
    const maxConcurrent = 3
    const numTasks = 8
    
    semaphore := make(chan struct{}, maxConcurrent)
    var wg sync.WaitGroup
    
    for i := 1; i <= numTasks; i++ {
        wg.Add(1)
        
        go func(taskID int) {
            defer wg.Done()
            
            // Acquire semaphore
            semaphore <- struct{}{}
            defer func() { <-semaphore }() // Release semaphore
            
            fmt.Printf("Task %d: Started (concurrent tasks: %d)\n", taskID, len(semaphore))
            
            // Simulate work
            time.Sleep(500 * time.Millisecond)
            
            fmt.Printf("Task %d: Completed\n", taskID)
        }(i)
        
        time.Sleep(50 * time.Millisecond) // Stagger task starts
    }
    
    wg.Wait()
    fmt.Println("All tasks completed")
}

// Broadcast pattern
func broadcastPattern() {
    fmt.Println("\n=== Broadcast Pattern ===")
    
    const numListeners = 4
    
    // Create channels for each listener
    listeners := make([]chan string, numListeners)
    for i := range listeners {
        listeners[i] = make(chan string, 5)
    }
    
    // Broadcaster goroutine
    go func() {
        messages := []string{"Hello", "World", "Broadcast", "Message"}
        
        for _, msg := range messages {
            fmt.Printf("Broadcasting: %s\n", msg)
            
            // Send to all listeners
            for i, ch := range listeners {
                select {
                case ch <- fmt.Sprintf("[Listener%d] %s", i+1, msg):
                default:
                    fmt.Printf("Listener %d channel full, dropping message\n", i+1)
                }
            }
            
            time.Sleep(200 * time.Millisecond)
        }
        
        // Close all channels
        for _, ch := range listeners {
            close(ch)
        }
    }()
    
    // Start listeners
    var wg sync.WaitGroup
    for i, ch := range listeners {
        wg.Add(1)
        
        go func(id int, ch <-chan string) {
            defer wg.Done()
            
            for msg := range ch {
                fmt.Printf("Listener %d received: %s\n", id, msg)
                time.Sleep(100 * time.Millisecond)
            }
        }(i+1, ch)
    }
    
    wg.Wait()
    fmt.Println("Broadcast complete")
}

// Request-Response pattern
func requestResponsePattern() {
    fmt.Println("\n=== Request-Response Pattern ===")
    
    type Request struct {
        ID       int
        Data     string
        Response chan string
    }
    
    requests := make(chan Request, 5)
    
    // Server goroutine
    go func() {
        for req := range requests {
            fmt.Printf("Server: Processing request %d: %s\n", req.ID, req.Data)
            
            // Simulate processing
            time.Sleep(200 * time.Millisecond)
            
            response := fmt.Sprintf("Processed: %s (ID: %d)", req.Data, req.ID)
            req.Response <- response
            close(req.Response)
        }
    }()
    
    // Client requests
    var wg sync.WaitGroup
    
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        
        go func(id int) {
            defer wg.Done()
            
            // Create request
            req := Request{
                ID:       id,
                Data:     fmt.Sprintf("Task-%d", id),
                Response: make(chan string, 1),
            }
            
            fmt.Printf("Client %d: Sending request\n", id)
            requests <- req
            
            // Wait for response
            response := <-req.Response
            fmt.Printf("Client %d: Received response: %s\n", id, response)
        }(i)
    }
    
    wg.Wait()
    close(requests)
    
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Request-Response pattern complete")
}

func main() {
    barrierPattern()
    semaphorePattern()
    broadcastPattern()
    requestResponsePattern()
}
```

### Channel-based Data Structures

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Channel-based Queue
type ChannelQueue struct {
    items chan interface{}
    size  int
}

func NewChannelQueue(capacity int) *ChannelQueue {
    return &ChannelQueue{
        items: make(chan interface{}, capacity),
        size:  capacity,
    }
}

func (q *ChannelQueue) Enqueue(item interface{}) bool {
    select {
    case q.items <- item:
        return true
    default:
        return false // Queue is full
    }
}

func (q *ChannelQueue) Dequeue() (interface{}, bool) {
    select {
    case item := <-q.items:
        return item, true
    default:
        return nil, false // Queue is empty
    }
}

func (q *ChannelQueue) DequeueBlocking() interface{} {
    return <-q.items
}

func (q *ChannelQueue) Size() int {
    return len(q.items)
}

func (q *ChannelQueue) Capacity() int {
    return q.size
}

func (q *ChannelQueue) Close() {
    close(q.items)
}

// Channel-based Stack
type ChannelStack struct {
    items chan interface{}
    stack []interface{}
    mutex sync.Mutex
}

func NewChannelStack(capacity int) *ChannelStack {
    cs := &ChannelStack{
        items: make(chan interface{}, capacity),
        stack: make([]interface{}, 0, capacity),
    }
    
    // Start the stack manager goroutine
    go cs.manager()
    
    return cs
}

func (s *ChannelStack) manager() {
    for item := range s.items {
        s.mutex.Lock()
        s.stack = append(s.stack, item)
        s.mutex.Unlock()
    }
}

func (s *ChannelStack) Push(item interface{}) bool {
    select {
    case s.items <- item:
        return true
    default:
        return false
    }
}

func (s *ChannelStack) Pop() (interface{}, bool) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if len(s.stack) == 0 {
        return nil, false
    }
    
    item := s.stack[len(s.stack)-1]
    s.stack = s.stack[:len(s.stack)-1]
    return item, true
}

func (s *ChannelStack) Size() int {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    return len(s.stack)
}

// Channel-based Priority Queue
type PriorityItem struct {
    Value    interface{}
    Priority int
}

type ChannelPriorityQueue struct {
    items chan PriorityItem
    queue []PriorityItem
    mutex sync.Mutex
    cond  *sync.Cond
}

func NewChannelPriorityQueue(capacity int) *ChannelPriorityQueue {
    pq := &ChannelPriorityQueue{
        items: make(chan PriorityItem, capacity),
        queue: make([]PriorityItem, 0),
    }
    pq.cond = sync.NewCond(&pq.mutex)
    
    go pq.manager()
    
    return pq
}

func (pq *ChannelPriorityQueue) manager() {
    for item := range pq.items {
        pq.mutex.Lock()
        
        // Insert in priority order (higher priority first)
        inserted := false
        for i, existing := range pq.queue {
            if item.Priority > existing.Priority {
                // Insert at position i
                pq.queue = append(pq.queue[:i], append([]PriorityItem{item}, pq.queue[i:]...)...)
                inserted = true
                break
            }
        }
        
        if !inserted {
            pq.queue = append(pq.queue, item)
        }
        
        pq.cond.Signal()
        pq.mutex.Unlock()
    }
}

func (pq *ChannelPriorityQueue) Enqueue(value interface{}, priority int) bool {
    item := PriorityItem{Value: value, Priority: priority}
    select {
    case pq.items <- item:
        return true
    default:
        return false
    }
}

func (pq *ChannelPriorityQueue) Dequeue() (interface{}, bool) {
    pq.mutex.Lock()
    defer pq.mutex.Unlock()
    
    if len(pq.queue) == 0 {
        return nil, false
    }
    
    item := pq.queue[0]
    pq.queue = pq.queue[1:]
    return item.Value, true
}

func (pq *ChannelPriorityQueue) DequeueBlocking() interface{} {
    pq.mutex.Lock()
    defer pq.mutex.Unlock()
    
    for len(pq.queue) == 0 {
        pq.cond.Wait()
    }
    
    item := pq.queue[0]
    pq.queue = pq.queue[1:]
    return item.Value
}

func (pq *ChannelPriorityQueue) Size() int {
    pq.mutex.Lock()
    defer pq.mutex.Unlock()
    return len(pq.queue)
}

func demonstrateDataStructures() {
    fmt.Println("=== Channel-based Data Structures ===")
    
    // Queue demonstration
    fmt.Println("\n--- Channel Queue ---")
    queue := NewChannelQueue(5)
    
    // Producer
    go func() {
        for i := 1; i <= 7; i++ {
            if queue.Enqueue(fmt.Sprintf("Item-%d", i)) {
                fmt.Printf("Enqueued: Item-%d (size: %d)\n", i, queue.Size())
            } else {
                fmt.Printf("Queue full, couldn't enqueue Item-%d\n", i)
            }
            time.Sleep(100 * time.Millisecond)
        }
    }()
    
    // Consumer
    time.Sleep(300 * time.Millisecond)
    for i := 0; i < 5; i++ {
        if item, ok := queue.Dequeue(); ok {
            fmt.Printf("Dequeued: %v (size: %d)\n", item, queue.Size())
        }
        time.Sleep(200 * time.Millisecond)
    }
    
    // Stack demonstration
    fmt.Println("\n--- Channel Stack ---")
    stack := NewChannelStack(5)
    
    // Push items
    for i := 1; i <= 5; i++ {
        if stack.Push(fmt.Sprintf("Item-%d", i)) {
            fmt.Printf("Pushed: Item-%d\n", i)
        }
        time.Sleep(50 * time.Millisecond)
    }
    
    time.Sleep(100 * time.Millisecond)
    
    // Pop items
    for i := 0; i < 5; i++ {
        if item, ok := stack.Pop(); ok {
            fmt.Printf("Popped: %v (size: %d)\n", item, stack.Size())
        }
        time.Sleep(50 * time.Millisecond)
    }
    
    // Priority Queue demonstration
    fmt.Println("\n--- Channel Priority Queue ---")
    pq := NewChannelPriorityQueue(10)
    
    // Enqueue with different priorities
    items := []struct {
        value    string
        priority int
    }{
        {"Low priority task", 1},
        {"High priority task", 10},
        {"Medium priority task", 5},
        {"Critical task", 15},
        {"Normal task", 3},
    }
    
    for _, item := range items {
        if pq.Enqueue(item.value, item.priority) {
            fmt.Printf("Enqueued: %s (priority: %d)\n", item.value, item.priority)
        }
        time.Sleep(50 * time.Millisecond)
    }
    
    time.Sleep(200 * time.Millisecond)
    
    // Dequeue in priority order
    fmt.Println("\nDequeuing in priority order:")
    for i := 0; i < len(items); i++ {
        if item, ok := pq.Dequeue(); ok {
            fmt.Printf("Dequeued: %v\n", item)
        }
    }
}

func main() {
    demonstrateDataStructures()
}
```

## Channel Anti-patterns and Best Practices

### Common Mistakes and Solutions

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Anti-pattern 1: Sending on closed channel
func antipatternClosedChannel() {
    fmt.Println("=== Anti-pattern: Sending on Closed Channel ===")
    
    ch := make(chan int, 1)
    close(ch)
    
    // This will panic!
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
        }
    }()
    
    ch <- 1 // This causes panic
}

// Better pattern: Check if channel is closed
func betterClosedChannelHandling() {
    fmt.Println("\n=== Better: Safe Channel Operations ===")
    
    ch := make(chan int, 1)
    var closed bool
    var mutex sync.Mutex
    
    safeSend := func(value int) bool {
        mutex.Lock()
        defer mutex.Unlock()
        
        if closed {
            return false
        }
        
        select {
        case ch <- value:
            return true
        default:
            return false
        }
    }
    
    safeClose := func() {
        mutex.Lock()
        defer mutex.Unlock()
        
        if !closed {
            close(ch)
            closed = true
        }
    }
    
    // Use safe operations
    if safeSend(1) {
        fmt.Println("Successfully sent 1")
    }
    
    safeClose()
    
    if !safeSend(2) {
        fmt.Println("Failed to send 2 (channel closed)")
    }
    
    // Receive remaining values
    for {
        val, ok := <-ch
        if !ok {
            break
        }
        fmt.Printf("Received: %d\n", val)
    }
}

// Anti-pattern 2: Goroutine leaks
func antipatternGoroutineLeak() {
    fmt.Println("\n=== Anti-pattern: Goroutine Leak ===")
    
    ch := make(chan int)
    
    // This goroutine will leak because no one receives from ch
    go func() {
        fmt.Println("Goroutine: Trying to send...")
        ch <- 1 // This will block forever
        fmt.Println("Goroutine: Sent successfully") // Never reached
    }()
    
    time.Sleep(100 * time.Millisecond)
    fmt.Println("Main: Exiting without receiving")
    // Goroutine is still blocked on send
}

// Better pattern: Use context or timeout
func betterGoroutineManagement() {
    fmt.Println("\n=== Better: Proper Goroutine Management ===")
    
    ch := make(chan int)
    done := make(chan bool)
    
    go func() {
        defer func() { done <- true }()
        
        fmt.Println("Goroutine: Trying to send...")
        
        select {
        case ch <- 1:
            fmt.Println("Goroutine: Sent successfully")
        case <-time.After(200 * time.Millisecond):
            fmt.Println("Goroutine: Send timeout, exiting")
        }
    }()
    
    time.Sleep(100 * time.Millisecond)
    
    // Wait for goroutine to finish
    <-done
    fmt.Println("Main: Goroutine finished properly")
}

// Anti-pattern 3: Unbuffered channel in same goroutine
func antipatternSameGoroutine() {
    fmt.Println("\n=== Anti-pattern: Unbuffered Channel in Same Goroutine ===")
    
    ch := make(chan int)
    
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from deadlock\n")
        }
    }()
    
    // This will deadlock!
    go func() {
        ch <- 1 // Send
        val := <-ch // Receive in same goroutine - deadlock!
        fmt.Printf("Received: %d\n", val)
    }()
    
    time.Sleep(100 * time.Millisecond)
    fmt.Println("This won't be reached due to deadlock")
}

// Better pattern: Use buffered channel or separate goroutines
func betterSameGoroutineHandling() {
    fmt.Println("\n=== Better: Avoid Same Goroutine Deadlock ===")
    
    // Solution 1: Buffered channel
    fmt.Println("Solution 1: Buffered channel")
    ch1 := make(chan int, 1)
    ch1 <- 1
    val := <-ch1
    fmt.Printf("Received from buffered: %d\n", val)
    
    // Solution 2: Separate goroutines
    fmt.Println("Solution 2: Separate goroutines")
    ch2 := make(chan int)
    
    go func() {
        ch2 <- 2
    }()
    
    val2 := <-ch2
    fmt.Printf("Received from separate goroutine: %d\n", val2)
}

// Anti-pattern 4: Not closing channels
func antipatternNotClosing() {
    fmt.Println("\n=== Anti-pattern: Not Closing Channels ===")
    
    ch := make(chan int, 3)
    
    // Send some values
    ch <- 1
    ch <- 2
    ch <- 3
    
    // Try to range over channel without closing
    go func() {
        fmt.Println("Goroutine: Starting range (will block forever)")
        for val := range ch {
            fmt.Printf("Received: %d\n", val)
        }
        fmt.Println("Goroutine: Range finished") // Never reached
    }()
    
    time.Sleep(200 * time.Millisecond)
    fmt.Println("Main: Exiting, goroutine still blocked")
}

// Better pattern: Always close channels when done
func betterChannelClosing() {
    fmt.Println("\n=== Better: Proper Channel Closing ===")
    
    ch := make(chan int, 3)
    
    // Producer
    go func() {
        defer close(ch) // Always close when done producing
        
        for i := 1; i <= 3; i++ {
            ch <- i
            time.Sleep(50 * time.Millisecond)
        }
        fmt.Println("Producer: Finished, channel closed")
    }()
    
    // Consumer
    for val := range ch {
        fmt.Printf("Received: %d\n", val)
    }
    
    fmt.Println("Consumer: Range finished properly")
}

// Best practices demonstration
func bestPractices() {
    fmt.Println("\n=== Channel Best Practices ===")
    
    // 1. Use buffered channels for known capacity
    fmt.Println("\n1. Buffered channels for known capacity:")
    results := make(chan string, 3) // Buffer size matches number of workers
    
    for i := 1; i <= 3; i++ {
        go func(id int) {
            time.Sleep(time.Duration(id*50) * time.Millisecond)
            results <- fmt.Sprintf("Result from worker %d", id)
        }(i)
    }
    
    for i := 0; i < 3; i++ {
        fmt.Printf("Collected: %s\n", <-results)
    }
    
    // 2. Use select with default for non-blocking operations
    fmt.Println("\n2. Non-blocking operations with select:")
    ch := make(chan int, 1)
    
    // Non-blocking send
    select {
    case ch <- 42:
        fmt.Println("Sent 42")
    default:
        fmt.Println("Channel full")
    }
    
    // Non-blocking receive
    select {
    case val := <-ch:
        fmt.Printf("Received: %d\n", val)
    default:
        fmt.Println("Channel empty")
    }
    
    // 3. Use context for cancellation
    fmt.Println("\n3. Context for cancellation:")
    // This would typically use context.Context
    cancel := make(chan bool)
    data := make(chan int)
    
    go func() {
        defer close(data)
        
        for i := 1; i <= 5; i++ {
            select {
            case data <- i:
                time.Sleep(100 * time.Millisecond)
            case <-cancel:
                fmt.Println("Producer: Cancelled")
                return
            }
        }
    }()
    
    // Cancel after receiving 2 items
    count := 0
    for val := range data {
        fmt.Printf("Received: %d\n", val)
        count++
        if count == 2 {
            close(cancel)
            break
        }
    }
    
    // Drain remaining values
    for range data {
        // Consume remaining
    }
    
    fmt.Println("Cancellation handled properly")
}

func main() {
    antipatternClosedChannel()
    betterClosedChannelHandling()
    antipatternGoroutineLeak()
    betterGoroutineManagement()
    antipatternSameGoroutine()
    betterSameGoroutineHandling()
    antipatternNotClosing()
    betterChannelClosing()
    bestPractices()
}
```

## Practical Examples

### Web Scraper with Channels

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// URL represents a URL to scrape
type URL struct {
    Address string
    Depth   int
}

// Result represents scraping result
type Result struct {
    URL     string
    Title   string
    Links   []string
    Error   error
    Duration time.Duration
}

// WebScraper manages concurrent web scraping
type WebScraper struct {
    maxWorkers   int
    maxDepth     int
    urlQueue     chan URL
    results      chan Result
    visited      map[string]bool
    visitedMutex sync.Mutex
    wg           sync.WaitGroup
}

// NewWebScraper creates a new web scraper
func NewWebScraper(maxWorkers, maxDepth int) *WebScraper {
    return &WebScraper{
        maxWorkers: maxWorkers,
        maxDepth:   maxDepth,
        urlQueue:   make(chan URL, maxWorkers*2),
        results:    make(chan Result, maxWorkers*2),
        visited:    make(map[string]bool),
    }
}

// isVisited checks if URL has been visited
func (ws *WebScraper) isVisited(url string) bool {
    ws.visitedMutex.Lock()
    defer ws.visitedMutex.Unlock()
    
    if ws.visited[url] {
        return true
    }
    ws.visited[url] = true
    return false
}

// scrapeURL simulates scraping a URL
func (ws *WebScraper) scrapeURL(url URL) Result {
    start := time.Now()
    
    // Simulate network delay
    delay := time.Duration(rand.Intn(500)+100) * time.Millisecond
    time.Sleep(delay)
    
    // Simulate occasional errors
    if rand.Float32() < 0.1 {
        return Result{
            URL:      url.Address,
            Error:    fmt.Errorf("failed to scrape %s", url.Address),
            Duration: time.Since(start),
        }
    }
    
    // Simulate successful scraping
    title := fmt.Sprintf("Title of %s", url.Address)
    links := []string{}
    
    // Generate some child links if not at max depth
    if url.Depth < ws.maxDepth {
        numLinks := rand.Intn(3) + 1
        for i := 0; i < numLinks; i++ {
            links = append(links, fmt.Sprintf("%s/page%d", url.Address, i+1))
        }
    }
    
    return Result{
        URL:      url.Address,
        Title:    title,
        Links:    links,
        Duration: time.Since(start),
    }
}

// worker processes URLs from the queue
func (ws *WebScraper) worker(id int) {
    defer ws.wg.Done()
    
    fmt.Printf("Worker %d started\n", id)
    
    for url := range ws.urlQueue {
        if ws.isVisited(url.Address) {
            continue
        }
        
        fmt.Printf("Worker %d scraping: %s (depth: %d)\n", id, url.Address, url.Depth)
        
        result := ws.scrapeURL(url)
        ws.results <- result
        
        // Add discovered links to queue
        if result.Error == nil && url.Depth < ws.maxDepth {
            for _, link := range result.Links {
                select {
                case ws.urlQueue <- URL{Address: link, Depth: url.Depth + 1}:
                default:
                    // Queue full, skip this link
                }
            }
        }
    }
    
    fmt.Printf("Worker %d finished\n", id)
}

// Scrape starts the scraping process
func (ws *WebScraper) Scrape(startURLs []string) <-chan Result {
    // Start workers
    for i := 1; i <= ws.maxWorkers; i++ {
        ws.wg.Add(1)
        go ws.worker(i)
    }
    
    // Add initial URLs
    go func() {
        for _, url := range startURLs {
            ws.urlQueue <- URL{Address: url, Depth: 0}
        }
    }()
    
    // Close results channel when all workers are done
    go func() {
        ws.wg.Wait()
        close(ws.results)
    }()
    
    return ws.results
}

// Stop stops the scraper
func (ws *WebScraper) Stop() {
    close(ws.urlQueue)
}

func demonstrateWebScraper() {
    fmt.Println("=== Web Scraper with Channels ===")
    
    rand.Seed(time.Now().UnixNano())
    
    scraper := NewWebScraper(3, 2)
    
    startURLs := []string{
        "https://example.com",
        "https://test.com",
    }
    
    results := scraper.Scrape(startURLs)
    
    // Collect results
    var successCount, errorCount int
    var totalDuration time.Duration
    
    // Stop scraper after some time
    go func() {
        time.Sleep(3 * time.Second)
        scraper.Stop()
    }()
    
    for result := range results {
        if result.Error != nil {
            fmt.Printf("❌ Error scraping %s: %v (took %v)\n", 
                result.URL, result.Error, result.Duration)
            errorCount++
        } else {
            fmt.Printf("✅ Scraped %s: %s (found %d links, took %v)\n", 
                result.URL, result.Title, len(result.Links), result.Duration)
            successCount++
        }
        totalDuration += result.Duration
    }
    
    fmt.Printf("\nScraping Summary:\n")
    fmt.Printf("Successful: %d\n", successCount)
    fmt.Printf("Errors: %d\n", errorCount)
    fmt.Printf("Total time: %v\n", totalDuration)
    if successCount+errorCount > 0 {
        fmt.Printf("Average time per URL: %v\n", totalDuration/time.Duration(successCount+errorCount))
    }
}

func main() {
    demonstrateWebScraper()
}
```

## Exercises

### Exercise 1: Chat Server
Create a chat server using channels for message distribution.

```go
// Create a chat server that:
// - Manages multiple chat rooms
// - Handles user join/leave events
// - Broadcasts messages to room members
// - Supports private messages
// - Maintains user presence

package main

type Message struct {
    Type     string // "join", "leave", "message", "private"
    From     string
    To       string // For private messages
    Room     string
    Content  string
    Timestamp time.Time
}

type User struct {
    Name     string
    Messages chan Message
    Rooms    map[string]bool
}

type ChatServer struct {
    // Add fields for managing users, rooms, etc.
}

// Implement these methods:
// NewChatServer() *ChatServer
// AddUser(name string) *User
// RemoveUser(name string)
// JoinRoom(user, room string) error
// LeaveRoom(user, room string) error
// SendMessage(msg Message) error
// BroadcastToRoom(room string, msg Message)
// Start() error
// Stop() error
```

### Exercise 2: Load Balancer
Implement a load balancer using channels.

```go
// Create a load balancer that:
// - Distributes requests across multiple servers
// - Implements different balancing algorithms (round-robin, least-connections)
// - Handles server health checks
// - Supports server addition/removal
// - Provides request/response metrics

package main

type Request struct {
    ID       string
    Data     interface{}
    Response chan Response
}

type Response struct {
    Data  interface{}
    Error error
}

type Server struct {
    ID          string
    Address     string
    Healthy     bool
    Connections int
    Requests    chan Request
}

type LoadBalancer struct {
    // Add fields for servers, algorithms, metrics, etc.
}

// Implement these methods:
// NewLoadBalancer(algorithm string) *LoadBalancer
// AddServer(server *Server) error
// RemoveServer(serverID string) error
// HandleRequest(req Request) error
// GetMetrics() Metrics
// Start() error
// Stop() error
```

### Exercise 3: Pipeline Processor
Create a configurable data processing pipeline.

```go
// Create a pipeline processor that:
// - Supports multiple processing stages
// - Allows dynamic stage addition/removal
// - Handles different data types
// - Provides processing metrics
// - Supports error handling and retries

package main

type ProcessingStage interface {
    Process(data interface{}) (interface{}, error)
    Name() string
}

type PipelineProcessor struct {
    // Add fields for stages, metrics, error handling, etc.
}

// Implement these methods:
// NewPipelineProcessor() *PipelineProcessor
// AddStage(stage ProcessingStage) error
// RemoveStage(name string) error
// Process(data interface{}) (interface{}, error)
// ProcessBatch(data []interface{}) ([]interface{}, error)
// GetMetrics() ProcessingMetrics
// Start() error
// Stop() error
```

## Key Takeaways

1. **Channels are the primary communication mechanism** between goroutines
2. **Unbuffered channels are synchronous** - sender blocks until receiver is ready
3. **Buffered channels are asynchronous** up to their capacity
4. **Always close channels** when done sending to enable proper range loops
5. **Use select for non-blocking operations** and multiplexing
6. **Channel directions** (send-only, receive-only) provide compile-time safety
7. **Avoid sending on closed channels** - it causes panics
8. **Use context or timeouts** to prevent goroutine leaks
9. **Channel-based patterns** (fan-in, fan-out, pipeline) solve common concurrency problems
10. **Proper channel management** prevents deadlocks and resource leaks

## Next Steps

Now that you understand channels, let's learn about [Synchronization Primitives](17-synchronization.md) for more advanced concurrency control!

---

**Previous**: [← Goroutines](15-goroutines.md) | **Next**: [Synchronization Primitives →](17-synchronization.md)