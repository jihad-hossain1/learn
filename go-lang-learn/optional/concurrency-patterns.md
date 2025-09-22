# 18. Concurrency Patterns

Go's concurrency model, built around goroutines and channels, enables powerful concurrent programming patterns. This chapter explores advanced concurrency patterns and best practices for building robust concurrent applications.

## Worker Pool Pattern

The worker pool pattern is used to limit the number of goroutines processing tasks concurrently.

### Basic Worker Pool

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID   int
    Data string
}

type Result struct {
    Job    Job
    Output string
    Error  error
}

type WorkerPool struct {
    numWorkers int
    jobs       chan Job
    results    chan Result
    wg         sync.WaitGroup
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    return &WorkerPool{
        numWorkers: numWorkers,
        jobs:       make(chan Job, 100),
        results:    make(chan Result, 100),
    }
}

func (wp *WorkerPool) Start() {
    // Start workers
    for i := 0; i < wp.numWorkers; i++ {
        wp.wg.Add(1)
        go wp.worker(i)
    }
}

func (wp *WorkerPool) worker(id int) {
    defer wp.wg.Done()
    
    for job := range wp.jobs {
        fmt.Printf("Worker %d processing job %d\n", id, job.ID)
        
        // Simulate work
        time.Sleep(time.Millisecond * 500)
        
        result := Result{
            Job:    job,
            Output: fmt.Sprintf("Processed: %s", job.Data),
        }
        
        wp.results <- result
    }
    
    fmt.Printf("Worker %d finished\n", id)
}

func (wp *WorkerPool) AddJob(job Job) {
    wp.jobs <- job
}

func (wp *WorkerPool) Close() {
    close(wp.jobs)
    wp.wg.Wait()
    close(wp.results)
}

func (wp *WorkerPool) Results() <-chan Result {
    return wp.results
}

func main() {
    // Create worker pool with 3 workers
    pool := NewWorkerPool(3)
    pool.Start()
    
    // Add jobs
    go func() {
        for i := 1; i <= 10; i++ {
            job := Job{
                ID:   i,
                Data: fmt.Sprintf("task-%d", i),
            }
            pool.AddJob(job)
        }
        pool.Close()
    }()
    
    // Collect results
    for result := range pool.Results() {
        if result.Error != nil {
            fmt.Printf("Job %d failed: %v\n", result.Job.ID, result.Error)
        } else {
            fmt.Printf("Job %d completed: %s\n", result.Job.ID, result.Output)
        }
    }
    
    fmt.Println("All jobs completed")
}
```

### Advanced Worker Pool with Context

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type Task func(ctx context.Context) (interface{}, error)

type WorkItem struct {
    ID   int
    Task Task
}

type WorkResult struct {
    ID     int
    Result interface{}
    Error  error
}

type AdvancedWorkerPool struct {
    numWorkers int
    workChan   chan WorkItem
    resultChan chan WorkResult
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
}

func NewAdvancedWorkerPool(ctx context.Context, numWorkers int) *AdvancedWorkerPool {
    ctx, cancel := context.WithCancel(ctx)
    
    return &AdvancedWorkerPool{
        numWorkers: numWorkers,
        workChan:   make(chan WorkItem, 100),
        resultChan: make(chan WorkResult, 100),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (awp *AdvancedWorkerPool) Start() {
    for i := 0; i < awp.numWorkers; i++ {
        awp.wg.Add(1)
        go awp.worker(i)
    }
    
    // Start result collector
    go awp.resultCollector()
}

func (awp *AdvancedWorkerPool) worker(id int) {
    defer awp.wg.Done()
    
    for {
        select {
        case <-awp.ctx.Done():
            fmt.Printf("Worker %d stopping due to context cancellation\n", id)
            return
        case work, ok := <-awp.workChan:
            if !ok {
                fmt.Printf("Worker %d stopping - work channel closed\n", id)
                return
            }
            
            fmt.Printf("Worker %d processing work item %d\n", id, work.ID)
            
            result, err := work.Task(awp.ctx)
            
            select {
            case awp.resultChan <- WorkResult{
                ID:     work.ID,
                Result: result,
                Error:  err,
            }:
            case <-awp.ctx.Done():
                return
            }
        }
    }
}

func (awp *AdvancedWorkerPool) resultCollector() {
    defer close(awp.resultChan)
    
    awp.wg.Wait()
}

func (awp *AdvancedWorkerPool) Submit(id int, task Task) error {
    select {
    case awp.workChan <- WorkItem{ID: id, Task: task}:
        return nil
    case <-awp.ctx.Done():
        return awp.ctx.Err()
    }
}

func (awp *AdvancedWorkerPool) Results() <-chan WorkResult {
    return awp.resultChan
}

func (awp *AdvancedWorkerPool) Stop() {
    close(awp.workChan)
    awp.cancel()
}

// Example tasks
func slowTask(ctx context.Context) (interface{}, error) {
    select {
    case <-time.After(2 * time.Second):
        return "Slow task completed", nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func fastTask(ctx context.Context) (interface{}, error) {
    select {
    case <-time.After(100 * time.Millisecond):
        return "Fast task completed", nil
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func errorTask(ctx context.Context) (interface{}, error) {
    return nil, fmt.Errorf("task failed")
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    pool := NewAdvancedWorkerPool(ctx, 3)
    pool.Start()
    
    // Submit tasks
    tasks := []Task{slowTask, fastTask, errorTask, fastTask, slowTask}
    
    go func() {
        for i, task := range tasks {
            if err := pool.Submit(i+1, task); err != nil {
                fmt.Printf("Failed to submit task %d: %v\n", i+1, err)
                break
            }
        }
        pool.Stop()
    }()
    
    // Collect results
    for result := range pool.Results() {
        if result.Error != nil {
            fmt.Printf("Task %d failed: %v\n", result.ID, result.Error)
        } else {
            fmt.Printf("Task %d completed: %v\n", result.ID, result.Result)
        }
    }
    
    fmt.Println("All tasks processed")
}
```

## Pipeline Pattern

The pipeline pattern chains multiple processing stages together using channels.

### Basic Pipeline

```go
package main

import (
    "fmt"
    "strconv"
    "strings"
    "time"
)

// Stage 1: Generate numbers
func generateNumbers(nums ...int) <-chan int {
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
func squareNumbers(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

// Stage 3: Convert to string
func numbersToStrings(in <-chan int) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for n := range in {
            out <- strconv.Itoa(n)
        }
    }()
    return out
}

// Stage 4: Add prefix
func addPrefix(in <-chan string, prefix string) <-chan string {
    out := make(chan string)
    go func() {
        defer close(out)
        for s := range in {
            out <- prefix + s
        }
    }()
    return out
}

func main() {
    // Create pipeline
    numbers := generateNumbers(1, 2, 3, 4, 5)
    squared := squareNumbers(numbers)
    strings := numbersToStrings(squared)
    prefixed := addPrefix(strings, "Result: ")
    
    // Consume results
    for result := range prefixed {
        fmt.Println(result)
    }
}
```

### Advanced Pipeline with Error Handling

```go
package main

import (
    "context"
    "fmt"
    "strconv"
    "sync"
    "time"
)

type PipelineData struct {
    Value interface{}
    Error error
}

type PipelineStage func(context.Context, <-chan PipelineData) <-chan PipelineData

// Pipeline orchestrator
type Pipeline struct {
    stages []PipelineStage
    ctx    context.Context
    cancel context.CancelFunc
}

func NewPipeline(ctx context.Context) *Pipeline {
    ctx, cancel := context.WithCancel(ctx)
    return &Pipeline{
        ctx:    ctx,
        cancel: cancel,
    }
}

func (p *Pipeline) AddStage(stage PipelineStage) {
    p.stages = append(p.stages, stage)
}

func (p *Pipeline) Run(input <-chan PipelineData) <-chan PipelineData {
    current := input
    
    for _, stage := range p.stages {
        current = stage(p.ctx, current)
    }
    
    return current
}

func (p *Pipeline) Stop() {
    p.cancel()
}

// Example stages
func parseIntStage(ctx context.Context, in <-chan PipelineData) <-chan PipelineData {
    out := make(chan PipelineData)
    
    go func() {
        defer close(out)
        
        for {
            select {
            case <-ctx.Done():
                return
            case data, ok := <-in:
                if !ok {
                    return
                }
                
                if data.Error != nil {
                    out <- data
                    continue
                }
                
                str, ok := data.Value.(string)
                if !ok {
                    out <- PipelineData{
                        Error: fmt.Errorf("expected string, got %T", data.Value),
                    }
                    continue
                }
                
                num, err := strconv.Atoi(str)
                if err != nil {
                    out <- PipelineData{
                        Error: fmt.Errorf("failed to parse int: %w", err),
                    }
                    continue
                }
                
                out <- PipelineData{Value: num}
            }
        }
    }()
    
    return out
}

func multiplyStage(multiplier int) PipelineStage {
    return func(ctx context.Context, in <-chan PipelineData) <-chan PipelineData {
        out := make(chan PipelineData)
        
        go func() {
            defer close(out)
            
            for {
                select {
                case <-ctx.Done():
                    return
                case data, ok := <-in:
                    if !ok {
                        return
                    }
                    
                    if data.Error != nil {
                        out <- data
                        continue
                    }
                    
                    num, ok := data.Value.(int)
                    if !ok {
                        out <- PipelineData{
                            Error: fmt.Errorf("expected int, got %T", data.Value),
                        }
                        continue
                    }
                    
                    out <- PipelineData{Value: num * multiplier}
                }
            }
        }()
        
        return out
    }
}

func formatStage(ctx context.Context, in <-chan PipelineData) <-chan PipelineData {
    out := make(chan PipelineData)
    
    go func() {
        defer close(out)
        
        for {
            select {
            case <-ctx.Done():
                return
            case data, ok := <-in:
                if !ok {
                    return
                }
                
                if data.Error != nil {
                    out <- data
                    continue
                }
                
                num, ok := data.Value.(int)
                if !ok {
                    out <- PipelineData{
                        Error: fmt.Errorf("expected int, got %T", data.Value),
                    }
                    continue
                }
                
                formatted := fmt.Sprintf("Result: %d", num)
                out <- PipelineData{Value: formatted}
            }
        }
    }()
    
    return out
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Create pipeline
    pipeline := NewPipeline(ctx)
    pipeline.AddStage(parseIntStage)
    pipeline.AddStage(multiplyStage(2))
    pipeline.AddStage(formatStage)
    
    // Create input
    input := make(chan PipelineData)
    go func() {
        defer close(input)
        
        inputs := []string{"1", "2", "invalid", "4", "5"}
        for _, inp := range inputs {
            select {
            case input <- PipelineData{Value: inp}:
            case <-ctx.Done():
                return
            }
        }
    }()
    
    // Run pipeline and collect results
    output := pipeline.Run(input)
    
    for result := range output {
        if result.Error != nil {
            fmt.Printf("Error: %v\n", result.Error)
        } else {
            fmt.Printf("Success: %v\n", result.Value)
        }
    }
}
```

## Fan-Out/Fan-In Pattern

This pattern distributes work across multiple goroutines (fan-out) and then combines the results (fan-in).

### Basic Fan-Out/Fan-In

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

// Fan-out: distribute work to multiple workers
func fanOut(input <-chan int, numWorkers int) []<-chan int {
    workers := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        worker := make(chan int)
        workers[i] = worker
        
        go func(w chan<- int) {
            defer close(w)
            
            for n := range input {
                // Simulate work
                time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
                w <- n * n
            }
        }(worker)
    }
    
    return workers
}

// Fan-in: combine results from multiple workers
func fanIn(workers ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    output := make(chan int)
    
    // Start a goroutine for each worker
    for _, worker := range workers {
        wg.Add(1)
        go func(w <-chan int) {
            defer wg.Done()
            for n := range w {
                output <- n
            }
        }(worker)
    }
    
    // Close output when all workers are done
    go func() {
        wg.Wait()
        close(output)
    }()
    
    return output
}

func main() {
    // Create input
    input := make(chan int)
    go func() {
        defer close(input)
        for i := 1; i <= 10; i++ {
            input <- i
        }
    }()
    
    // Fan-out to 3 workers
    workers := fanOut(input, 3)
    
    // Fan-in results
    results := fanIn(workers...)
    
    // Collect results
    var allResults []int
    for result := range results {
        allResults = append(allResults, result)
    }
    
    fmt.Printf("Results: %v\n", allResults)
}
```

### Advanced Fan-Out/Fan-In with Load Balancing

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type WorkItem struct {
    ID   int
    Data interface{}
}

type WorkResult struct {
    WorkerID int
    Item     WorkItem
    Result   interface{}
    Error    error
    Duration time.Duration
}

type LoadBalancer struct {
    workers    []chan WorkItem
    results    chan WorkResult
    ctx        context.Context
    cancel     context.CancelFunc
    wg         sync.WaitGroup
    numWorkers int
}

func NewLoadBalancer(ctx context.Context, numWorkers int) *LoadBalancer {
    ctx, cancel := context.WithCancel(ctx)
    
    lb := &LoadBalancer{
        workers:    make([]chan WorkItem, numWorkers),
        results:    make(chan WorkResult, numWorkers*2),
        ctx:        ctx,
        cancel:     cancel,
        numWorkers: numWorkers,
    }
    
    // Initialize worker channels
    for i := 0; i < numWorkers; i++ {
        lb.workers[i] = make(chan WorkItem, 10)
    }
    
    return lb
}

func (lb *LoadBalancer) Start(processor func(WorkItem) (interface{}, error)) {
    // Start workers
    for i := 0; i < lb.numWorkers; i++ {
        lb.wg.Add(1)
        go lb.worker(i, processor)
    }
    
    // Start result collector
    go func() {
        lb.wg.Wait()
        close(lb.results)
    }()
}

func (lb *LoadBalancer) worker(id int, processor func(WorkItem) (interface{}, error)) {
    defer lb.wg.Done()
    
    for {
        select {
        case <-lb.ctx.Done():
            return
        case item, ok := <-lb.workers[id]:
            if !ok {
                return
            }
            
            start := time.Now()
            result, err := processor(item)
            duration := time.Since(start)
            
            select {
            case lb.results <- WorkResult{
                WorkerID: id,
                Item:     item,
                Result:   result,
                Error:    err,
                Duration: duration,
            }:
            case <-lb.ctx.Done():
                return
            }
        }
    }
}

func (lb *LoadBalancer) Submit(item WorkItem) error {
    // Simple round-robin load balancing
    workerID := item.ID % lb.numWorkers
    
    select {
    case lb.workers[workerID] <- item:
        return nil
    case <-lb.ctx.Done():
        return lb.ctx.Err()
    }
}

func (lb *LoadBalancer) Results() <-chan WorkResult {
    return lb.results
}

func (lb *LoadBalancer) Stop() {
    for _, worker := range lb.workers {
        close(worker)
    }
    lb.cancel()
}

// Example processor function
func processItem(item WorkItem) (interface{}, error) {
    // Simulate variable processing time
    processingTime := time.Duration(rand.Intn(500)+100) * time.Millisecond
    time.Sleep(processingTime)
    
    // Simulate occasional errors
    if rand.Float32() < 0.1 {
        return nil, fmt.Errorf("processing failed for item %d", item.ID)
    }
    
    return fmt.Sprintf("Processed: %v", item.Data), nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Create load balancer with 4 workers
    lb := NewLoadBalancer(ctx, 4)
    lb.Start(processItem)
    
    // Submit work items
    go func() {
        for i := 1; i <= 20; i++ {
            item := WorkItem{
                ID:   i,
                Data: fmt.Sprintf("task-%d", i),
            }
            
            if err := lb.Submit(item); err != nil {
                fmt.Printf("Failed to submit item %d: %v\n", i, err)
                break
            }
        }
        lb.Stop()
    }()
    
    // Collect and display results
    workerStats := make(map[int]int)
    totalDuration := time.Duration(0)
    successCount := 0
    errorCount := 0
    
    for result := range lb.Results() {
        workerStats[result.WorkerID]++
        totalDuration += result.Duration
        
        if result.Error != nil {
            fmt.Printf("Worker %d - Item %d failed: %v (took %v)\n", 
                result.WorkerID, result.Item.ID, result.Error, result.Duration)
            errorCount++
        } else {
            fmt.Printf("Worker %d - Item %d: %v (took %v)\n", 
                result.WorkerID, result.Item.ID, result.Result, result.Duration)
            successCount++
        }
    }
    
    // Display statistics
    fmt.Println("\n=== Statistics ===")
    fmt.Printf("Total items processed: %d\n", successCount+errorCount)
    fmt.Printf("Successful: %d\n", successCount)
    fmt.Printf("Failed: %d\n", errorCount)
    fmt.Printf("Average processing time: %v\n", totalDuration/time.Duration(successCount+errorCount))
    
    fmt.Println("\nWork distribution:")
    for workerID, count := range workerStats {
        fmt.Printf("Worker %d: %d items\n", workerID, count)
    }
}
```

## Producer-Consumer Pattern

This pattern separates data production from consumption, allowing them to operate at different rates.

### Basic Producer-Consumer

```go
package main

import (
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type Message struct {
    ID        int
    Content   string
    Timestamp time.Time
}

type Producer struct {
    id       int
    output   chan<- Message
    stopChan chan struct{}
    wg       *sync.WaitGroup
}

func NewProducer(id int, output chan<- Message, wg *sync.WaitGroup) *Producer {
    return &Producer{
        id:       id,
        output:   output,
        stopChan: make(chan struct{}),
        wg:       wg,
    }
}

func (p *Producer) Start() {
    p.wg.Add(1)
    go p.produce()
}

func (p *Producer) Stop() {
    close(p.stopChan)
}

func (p *Producer) produce() {
    defer p.wg.Done()
    
    messageID := 1
    ticker := time.NewTicker(time.Duration(rand.Intn(500)+200) * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-p.stopChan:
            fmt.Printf("Producer %d stopping\n", p.id)
            return
        case <-ticker.C:
            message := Message{
                ID:        messageID,
                Content:   fmt.Sprintf("Message from producer %d", p.id),
                Timestamp: time.Now(),
            }
            
            select {
            case p.output <- message:
                fmt.Printf("Producer %d sent message %d\n", p.id, messageID)
                messageID++
            case <-p.stopChan:
                return
            }
        }
    }
}

type Consumer struct {
    id       int
    input    <-chan Message
    stopChan chan struct{}
    wg       *sync.WaitGroup
}

func NewConsumer(id int, input <-chan Message, wg *sync.WaitGroup) *Consumer {
    return &Consumer{
        id:       id,
        input:    input,
        stopChan: make(chan struct{}),
        wg:       wg,
    }
}

func (c *Consumer) Start() {
    c.wg.Add(1)
    go c.consume()
}

func (c *Consumer) Stop() {
    close(c.stopChan)
}

func (c *Consumer) consume() {
    defer c.wg.Done()
    
    for {
        select {
        case <-c.stopChan:
            fmt.Printf("Consumer %d stopping\n", c.id)
            return
        case message, ok := <-c.input:
            if !ok {
                fmt.Printf("Consumer %d: input channel closed\n", c.id)
                return
            }
            
            // Simulate processing time
            processingTime := time.Duration(rand.Intn(300)+100) * time.Millisecond
            time.Sleep(processingTime)
            
            fmt.Printf("Consumer %d processed message %d: %s (took %v)\n", 
                c.id, message.ID, message.Content, processingTime)
        }
    }
}

func main() {
    // Create message channel with buffer
    messageChannel := make(chan Message, 10)
    
    var wg sync.WaitGroup
    
    // Create and start producers
    producers := make([]*Producer, 2)
    for i := 0; i < 2; i++ {
        producers[i] = NewProducer(i+1, messageChannel, &wg)
        producers[i].Start()
    }
    
    // Create and start consumers
    consumers := make([]*Consumer, 3)
    for i := 0; i < 3; i++ {
        consumers[i] = NewConsumer(i+1, messageChannel, &wg)
        consumers[i].Start()
    }
    
    // Run for 5 seconds
    time.Sleep(5 * time.Second)
    
    // Stop producers
    fmt.Println("Stopping producers...")
    for _, producer := range producers {
        producer.Stop()
    }
    
    // Wait for producers to finish
    wg.Wait()
    
    // Close message channel
    close(messageChannel)
    
    // Wait for consumers to finish processing remaining messages
    wg.Wait()
    
    fmt.Println("All producers and consumers stopped")
}
```

### Advanced Producer-Consumer with Backpressure

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

type QueueItem struct {
    ID       int64
    Data     interface{}
    Priority int
    Created  time.Time
}

type QueueStats struct {
    Produced  int64
    Consumed  int64
    Dropped   int64
    QueueSize int64
}

type BackpressureQueue struct {
    items       chan QueueItem
    maxSize     int
    stats       QueueStats
    ctx         context.Context
    cancel      context.CancelFunc
    wg          sync.WaitGroup
    statsLock   sync.RWMutex
}

func NewBackpressureQueue(ctx context.Context, maxSize int) *BackpressureQueue {
    ctx, cancel := context.WithCancel(ctx)
    
    return &BackpressureQueue{
        items:   make(chan QueueItem, maxSize),
        maxSize: maxSize,
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (bq *BackpressureQueue) Produce(item QueueItem) error {
    select {
    case bq.items <- item:
        atomic.AddInt64(&bq.stats.Produced, 1)
        atomic.AddInt64(&bq.stats.QueueSize, 1)
        return nil
    case <-bq.ctx.Done():
        return bq.ctx.Err()
    default:
        // Queue is full, apply backpressure
        atomic.AddInt64(&bq.stats.Dropped, 1)
        return fmt.Errorf("queue is full, item dropped")
    }
}

func (bq *BackpressureQueue) ProduceWithTimeout(item QueueItem, timeout time.Duration) error {
    ctx, cancel := context.WithTimeout(bq.ctx, timeout)
    defer cancel()
    
    select {
    case bq.items <- item:
        atomic.AddInt64(&bq.stats.Produced, 1)
        atomic.AddInt64(&bq.stats.QueueSize, 1)
        return nil
    case <-ctx.Done():
        if ctx.Err() == context.DeadlineExceeded {
            atomic.AddInt64(&bq.stats.Dropped, 1)
            return fmt.Errorf("timeout waiting to produce item")
        }
        return ctx.Err()
    }
}

func (bq *BackpressureQueue) Consume() <-chan QueueItem {
    return bq.items
}

func (bq *BackpressureQueue) StartConsumer(id int, processor func(QueueItem) error) {
    bq.wg.Add(1)
    go func() {
        defer bq.wg.Done()
        
        for {
            select {
            case <-bq.ctx.Done():
                fmt.Printf("Consumer %d stopping\n", id)
                return
            case item, ok := <-bq.items:
                if !ok {
                    fmt.Printf("Consumer %d: queue closed\n", id)
                    return
                }
                
                atomic.AddInt64(&bq.stats.QueueSize, -1)
                
                if err := processor(item); err != nil {
                    fmt.Printf("Consumer %d failed to process item %d: %v\n", 
                        id, item.ID, err)
                } else {
                    atomic.AddInt64(&bq.stats.Consumed, 1)
                }
            }
        }
    }()
}

func (bq *BackpressureQueue) GetStats() QueueStats {
    return QueueStats{
        Produced:  atomic.LoadInt64(&bq.stats.Produced),
        Consumed:  atomic.LoadInt64(&bq.stats.Consumed),
        Dropped:   atomic.LoadInt64(&bq.stats.Dropped),
        QueueSize: atomic.LoadInt64(&bq.stats.QueueSize),
    }
}

func (bq *BackpressureQueue) Close() {
    bq.cancel()
    close(bq.items)
    bq.wg.Wait()
}

// Example processor
func processItem(item QueueItem) error {
    // Simulate processing time
    time.Sleep(100 * time.Millisecond)
    
    fmt.Printf("Processed item %d with data: %v\n", item.ID, item.Data)
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Create queue with small buffer to demonstrate backpressure
    queue := NewBackpressureQueue(ctx, 5)
    
    // Start consumers
    numConsumers := 2
    for i := 1; i <= numConsumers; i++ {
        queue.StartConsumer(i, processItem)
    }
    
    // Start producer
    go func() {
        var itemID int64 = 1
        ticker := time.NewTicker(50 * time.Millisecond)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                item := QueueItem{
                    ID:      itemID,
                    Data:    fmt.Sprintf("data-%d", itemID),
                    Created: time.Now(),
                }
                
                // Try to produce with timeout
                if err := queue.ProduceWithTimeout(item, 100*time.Millisecond); err != nil {
                    fmt.Printf("Failed to produce item %d: %v\n", itemID, err)
                } else {
                    fmt.Printf("Produced item %d\n", itemID)
                }
                
                itemID++
            }
        }
    }()
    
    // Monitor stats
    statsTicker := time.NewTicker(1 * time.Second)
    defer statsTicker.Stop()
    
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-statsTicker.C:
                stats := queue.GetStats()
                fmt.Printf("\n=== Stats ===\n")
                fmt.Printf("Produced: %d\n", stats.Produced)
                fmt.Printf("Consumed: %d\n", stats.Consumed)
                fmt.Printf("Dropped: %d\n", stats.Dropped)
                fmt.Printf("Queue Size: %d\n", stats.QueueSize)
                fmt.Printf("=============\n\n")
            }
        }
    }()
    
    // Wait for context to expire
    <-ctx.Done()
    
    // Close queue and wait for cleanup
    queue.Close()
    
    // Final stats
    finalStats := queue.GetStats()
    fmt.Printf("\n=== Final Stats ===\n")
    fmt.Printf("Total Produced: %d\n", finalStats.Produced)
    fmt.Printf("Total Consumed: %d\n", finalStats.Consumed)
    fmt.Printf("Total Dropped: %d\n", finalStats.Dropped)
    fmt.Printf("Remaining in Queue: %d\n", finalStats.QueueSize)
}
```

## Rate Limiting Pattern

Rate limiting controls the rate at which operations are performed.

### Token Bucket Rate Limiter

```go
package main

import (
    "context"
    "fmt"
    "sync"
    "time"
)

type TokenBucket struct {
    capacity     int
    tokens       int
    refillRate   int           // tokens per second
    lastRefill   time.Time
    mutex        sync.Mutex
    refillTicker *time.Ticker
    stopChan     chan struct{}
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
    tb := &TokenBucket{
        capacity:   capacity,
        tokens:     capacity,
        refillRate: refillRate,
        lastRefill: time.Now(),
        stopChan:   make(chan struct{}),
    }
    
    // Start refill goroutine
    tb.refillTicker = time.NewTicker(time.Second / time.Duration(refillRate))
    go tb.refillLoop()
    
    return tb
}

func (tb *TokenBucket) refillLoop() {
    for {
        select {
        case <-tb.refillTicker.C:
            tb.mutex.Lock()
            if tb.tokens < tb.capacity {
                tb.tokens++
            }
            tb.mutex.Unlock()
        case <-tb.stopChan:
            tb.refillTicker.Stop()
            return
        }
    }
}

func (tb *TokenBucket) TryConsume(tokens int) bool {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()
    
    if tb.tokens >= tokens {
        tb.tokens -= tokens
        return true
    }
    return false
}

func (tb *TokenBucket) ConsumeWithTimeout(tokens int, timeout time.Duration) error {
    deadline := time.Now().Add(timeout)
    
    for time.Now().Before(deadline) {
        if tb.TryConsume(tokens) {
            return nil
        }
        time.Sleep(10 * time.Millisecond)
    }
    
    return fmt.Errorf("timeout waiting for tokens")
}

func (tb *TokenBucket) AvailableTokens() int {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()
    return tb.tokens
}

func (tb *TokenBucket) Close() {
    close(tb.stopChan)
}

// Rate-limited worker
type RateLimitedWorker struct {
    id          int
    rateLimiter *TokenBucket
    workChan    chan func()
    ctx         context.Context
    wg          *sync.WaitGroup
}

func NewRateLimitedWorker(id int, rateLimiter *TokenBucket, ctx context.Context, wg *sync.WaitGroup) *RateLimitedWorker {
    return &RateLimitedWorker{
        id:          id,
        rateLimiter: rateLimiter,
        workChan:    make(chan func(), 10),
        ctx:         ctx,
        wg:          wg,
    }
}

func (w *RateLimitedWorker) Start() {
    w.wg.Add(1)
    go w.work()
}

func (w *RateLimitedWorker) Submit(work func()) error {
    select {
    case w.workChan <- work:
        return nil
    case <-w.ctx.Done():
        return w.ctx.Err()
    default:
        return fmt.Errorf("worker %d queue is full", w.id)
    }
}

func (w *RateLimitedWorker) work() {
    defer w.wg.Done()
    
    for {
        select {
        case <-w.ctx.Done():
            fmt.Printf("Worker %d stopping\n", w.id)
            return
        case work, ok := <-w.workChan:
            if !ok {
                return
            }
            
            // Wait for rate limiter
            for !w.rateLimiter.TryConsume(1) {
                select {
                case <-w.ctx.Done():
                    return
                case <-time.After(10 * time.Millisecond):
                    // Continue waiting
                }
            }
            
            // Execute work
            start := time.Now()
            work()
            duration := time.Since(start)
            
            fmt.Printf("Worker %d completed work in %v (tokens: %d)\n", 
                w.id, duration, w.rateLimiter.AvailableTokens())
        }
    }
}

func (w *RateLimitedWorker) Close() {
    close(w.workChan)
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    
    // Create rate limiter: 5 tokens capacity, refill 2 tokens per second
    rateLimiter := NewTokenBucket(5, 2)
    defer rateLimiter.Close()
    
    var wg sync.WaitGroup
    
    // Create workers
    numWorkers := 3
    workers := make([]*RateLimitedWorker, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        workers[i] = NewRateLimitedWorker(i+1, rateLimiter, ctx, &wg)
        workers[i].Start()
    }
    
    // Submit work
    go func() {
        workID := 1
        ticker := time.NewTicker(200 * time.Millisecond)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                currentWorkID := workID
                work := func() {
                    // Simulate work
                    time.Sleep(100 * time.Millisecond)
                    fmt.Printf("Completed work item %d\n", currentWorkID)
                }
                
                // Round-robin work distribution
                workerIndex := (workID - 1) % numWorkers
                if err := workers[workerIndex].Submit(work); err != nil {
                    fmt.Printf("Failed to submit work %d to worker %d: %v\n", 
                        workID, workerIndex+1, err)
                }
                
                workID++
            }
        }
    }()
    
    // Monitor rate limiter
    go func() {
        ticker := time.NewTicker(1 * time.Second)
        defer ticker.Stop()
        
        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                fmt.Printf("Rate limiter tokens: %d/%d\n", 
                    rateLimiter.AvailableTokens(), rateLimiter.capacity)
            }
        }
    }()
    
    // Wait for context to expire
    <-ctx.Done()
    
    // Close workers
    for _, worker := range workers {
        worker.Close()
    }
    
    // Wait for all workers to finish
    wg.Wait()
    
    fmt.Println("All workers finished")
}
```