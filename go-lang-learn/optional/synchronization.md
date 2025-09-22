# 17. Synchronization Primitives

Go provides several synchronization primitives in the `sync` package to coordinate access to shared resources and manage goroutine execution. While channels are the preferred way to communicate between goroutines, synchronization primitives are essential for protecting shared state and coordinating goroutine lifecycles.

## Mutex (Mutual Exclusion)

### Basic Mutex Usage

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Counter with mutex protection
type SafeCounter struct {
    mu    sync.Mutex
    value int
}

// Increment safely increments the counter
func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// Decrement safely decrements the counter
func (c *SafeCounter) Decrement() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value--
}

// Value safely returns the current value
func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

// Add safely adds a value to the counter
func (c *SafeCounter) Add(delta int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value += delta
}

func demonstrateBasicMutex() {
    fmt.Println("=== Basic Mutex Usage ===")
    
    counter := &SafeCounter{}
    var wg sync.WaitGroup
    
    // Start multiple goroutines that increment
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 100; j++ {
                counter.Increment()
                if j%25 == 0 {
                    fmt.Printf("Goroutine %d: Counter at %d\n", id, counter.Value())
                }
                time.Sleep(time.Microsecond)
            }
        }(i)
    }
    
    // Start multiple goroutines that decrement
    for i := 0; i < 2; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 50; j++ {
                counter.Decrement()
                time.Sleep(time.Microsecond)
            }
        }(i + 5)
    }
    
    wg.Wait()
    fmt.Printf("Final counter value: %d\n", counter.Value())
    fmt.Printf("Expected value: %d\n", 5*100-2*50)
}

// Demonstrate race condition without mutex
type UnsafeCounter struct {
    value int
}

func (c *UnsafeCounter) Increment() {
    c.value++ // Race condition!
}

func (c *UnsafeCounter) Value() int {
    return c.value // Race condition!
}

func demonstrateRaceCondition() {
    fmt.Println("\n=== Race Condition Demo ===")
    
    counter := &UnsafeCounter{}
    var wg sync.WaitGroup
    
    // Multiple goroutines incrementing without synchronization
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := 0; j < 1000; j++ {
                counter.Increment()
            }
        }()
    }
    
    wg.Wait()
    fmt.Printf("Unsafe counter final value: %d\n", counter.Value())
    fmt.Printf("Expected value: %d\n", 10*1000)
    fmt.Println("Note: Unsafe counter likely shows incorrect result due to race condition")
}

func main() {
    demonstrateBasicMutex()
    demonstrateRaceCondition()
}
```

### RWMutex (Read-Write Mutex)

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Cache with read-write mutex
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

// NewCache creates a new cache
func NewCache() *Cache {
    return &Cache{
        data: make(map[string]string),
    }
}

// Set adds or updates a key-value pair
func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    fmt.Printf("Setting %s = %s\n", key, value)
    c.data[key] = value
    time.Sleep(10 * time.Millisecond) // Simulate work
}

// Get retrieves a value by key
func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    value, exists := c.data[key]
    fmt.Printf("Getting %s = %s (exists: %t)\n", key, value, exists)
    time.Sleep(5 * time.Millisecond) // Simulate work
    return value, exists
}

// GetAll returns all key-value pairs
func (c *Cache) GetAll() map[string]string {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    // Create a copy to avoid exposing internal map
    result := make(map[string]string)
    for k, v := range c.data {
        result[k] = v
    }
    
    fmt.Printf("Getting all data: %d items\n", len(result))
    time.Sleep(15 * time.Millisecond) // Simulate work
    return result
}

// Delete removes a key-value pair
func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    fmt.Printf("Deleting %s\n", key)
    delete(c.data, key)
    time.Sleep(8 * time.Millisecond) // Simulate work
}

// Size returns the number of items in cache
func (c *Cache) Size() int {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    return len(c.data)
}

func demonstrateRWMutex() {
    fmt.Println("=== RWMutex Demonstration ===")
    
    cache := NewCache()
    var wg sync.WaitGroup
    
    // Writers (fewer, slower)
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 5; j++ {
                key := fmt.Sprintf("key%d_%d", id, j)
                value := fmt.Sprintf("value%d_%d", id, j)
                cache.Set(key, value)
                time.Sleep(20 * time.Millisecond)
            }
        }(i)
    }
    
    // Readers (more, faster)
    for i := 0; i < 8; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 10; j++ {
                if j%3 == 0 {
                    // Occasionally get all data
                    cache.GetAll()
                } else {
                    // Usually get specific keys
                    key := fmt.Sprintf("key%d_%d", id%3, j%5)
                    cache.Get(key)
                }
                time.Sleep(10 * time.Millisecond)
            }
        }(i)
    }
    
    // Deleters (occasional)
    wg.Add(1)
    go func() {
        defer wg.Done()
        
        time.Sleep(100 * time.Millisecond) // Let some data accumulate
        
        for i := 0; i < 3; i++ {
            key := fmt.Sprintf("key0_%d", i)
            cache.Delete(key)
            time.Sleep(50 * time.Millisecond)
        }
    }()
    
    wg.Wait()
    
    fmt.Printf("\nFinal cache size: %d\n", cache.Size())
    fmt.Println("Final cache contents:")
    for k, v := range cache.GetAll() {
        fmt.Printf("  %s: %s\n", k, v)
    }
}

// Performance comparison: Mutex vs RWMutex
func comparePerformance() {
    fmt.Println("\n=== Performance Comparison ===")
    
    const numReaders = 10
    const numReads = 1000
    
    // Test with regular mutex
    fmt.Println("Testing with regular Mutex...")
    start := time.Now()
    
    type MutexCache struct {
        mu   sync.Mutex
        data map[string]string
    }
    
    mutexCache := &MutexCache{data: map[string]string{"key": "value"}}
    var wg1 sync.WaitGroup
    
    for i := 0; i < numReaders; i++ {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            for j := 0; j < numReads; j++ {
                mutexCache.mu.Lock()
                _ = mutexCache.data["key"]
                mutexCache.mu.Unlock()
            }
        }()
    }
    
    wg1.Wait()
    mutexTime := time.Since(start)
    
    // Test with RWMutex
    fmt.Println("Testing with RWMutex...")
    start = time.Now()
    
    type RWMutexCache struct {
        mu   sync.RWMutex
        data map[string]string
    }
    
    rwMutexCache := &RWMutexCache{data: map[string]string{"key": "value"}}
    var wg2 sync.WaitGroup
    
    for i := 0; i < numReaders; i++ {
        wg2.Add(1)
        go func() {
            defer wg2.Done()
            for j := 0; j < numReads; j++ {
                rwMutexCache.mu.RLock()
                _ = rwMutexCache.data["key"]
                rwMutexCache.mu.RUnlock()
            }
        }()
    }
    
    wg2.Wait()
    rwMutexTime := time.Since(start)
    
    fmt.Printf("Mutex time: %v\n", mutexTime)
    fmt.Printf("RWMutex time: %v\n", rwMutexTime)
    fmt.Printf("RWMutex is %.2fx faster for read-heavy workloads\n", 
        float64(mutexTime)/float64(rwMutexTime))
}

func main() {
    demonstrateRWMutex()
    comparePerformance()
}
```

## WaitGroup

### Basic WaitGroup Usage

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Task represents a unit of work
type Task struct {
    ID       int
    Name     string
    Duration time.Duration
}

// Worker processes tasks
func worker(id int, tasks []Task, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Printf("Worker %d started\n", id)
    
    for _, task := range tasks {
        fmt.Printf("Worker %d processing task %d: %s\n", id, task.ID, task.Name)
        time.Sleep(task.Duration)
        fmt.Printf("Worker %d completed task %d\n", id, task.ID)
    }
    
    fmt.Printf("Worker %d finished\n", id)
}

func demonstrateBasicWaitGroup() {
    fmt.Println("=== Basic WaitGroup Usage ===")
    
    var wg sync.WaitGroup
    
    // Create tasks
    allTasks := []Task{
        {1, "Process data", 200 * time.Millisecond},
        {2, "Generate report", 300 * time.Millisecond},
        {3, "Send email", 150 * time.Millisecond},
        {4, "Update database", 400 * time.Millisecond},
        {5, "Backup files", 250 * time.Millisecond},
        {6, "Clean logs", 100 * time.Millisecond},
    }
    
    // Distribute tasks among workers
    numWorkers := 3
    tasksPerWorker := len(allTasks) / numWorkers
    
    for i := 0; i < numWorkers; i++ {
        start := i * tasksPerWorker
        end := start + tasksPerWorker
        if i == numWorkers-1 {
            end = len(allTasks) // Last worker gets remaining tasks
        }
        
        wg.Add(1)
        go worker(i+1, allTasks[start:end], &wg)
    }
    
    fmt.Println("Waiting for all workers to complete...")
    wg.Wait()
    fmt.Println("All workers completed!")
}

// Demonstrate WaitGroup with dynamic task creation
func demonstrateDynamicWaitGroup() {
    fmt.Println("\n=== Dynamic WaitGroup Usage ===")
    
    var wg sync.WaitGroup
    taskChan := make(chan Task, 10)
    
    // Start workers
    numWorkers := 3
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            fmt.Printf("Dynamic Worker %d started\n", workerID)
            
            for task := range taskChan {
                fmt.Printf("Dynamic Worker %d processing task %d: %s\n", 
                    workerID, task.ID, task.Name)
                time.Sleep(task.Duration)
                fmt.Printf("Dynamic Worker %d completed task %d\n", workerID, task.ID)
            }
            
            fmt.Printf("Dynamic Worker %d finished\n", workerID)
        }(i)
    }
    
    // Generate tasks dynamically
    go func() {
        defer close(taskChan)
        
        for i := 1; i <= 8; i++ {
            task := Task{
                ID:       i,
                Name:     fmt.Sprintf("Dynamic task %d", i),
                Duration: time.Duration(100+i*50) * time.Millisecond,
            }
            
            fmt.Printf("Generating task %d\n", i)
            taskChan <- task
            time.Sleep(100 * time.Millisecond)
        }
        
        fmt.Println("All tasks generated")
    }()
    
    wg.Wait()
    fmt.Println("All dynamic workers completed!")
}

// Nested WaitGroups example
func demonstrateNestedWaitGroups() {
    fmt.Println("\n=== Nested WaitGroups ===")
    
    var mainWG sync.WaitGroup
    
    // Main phases
    phases := []string{"Preparation", "Execution", "Cleanup"}
    
    for phaseIdx, phase := range phases {
        mainWG.Add(1)
        
        go func(phaseIndex int, phaseName string) {
            defer mainWG.Done()
            
            fmt.Printf("Phase %d: %s started\n", phaseIndex+1, phaseName)
            
            var phaseWG sync.WaitGroup
            
            // Sub-tasks within each phase
            numSubTasks := 3
            for i := 1; i <= numSubTasks; i++ {
                phaseWG.Add(1)
                
                go func(taskID int) {
                    defer phaseWG.Done()
                    
                    taskName := fmt.Sprintf("%s-SubTask-%d", phaseName, taskID)
                    duration := time.Duration(100+taskID*50) * time.Millisecond
                    
                    fmt.Printf("  %s started\n", taskName)
                    time.Sleep(duration)
                    fmt.Printf("  %s completed\n", taskName)
                }(i)
            }
            
            phaseWG.Wait()
            fmt.Printf("Phase %d: %s completed\n", phaseIndex+1, phaseName)
        }(phaseIdx, phase)
    }
    
    mainWG.Wait()
    fmt.Println("All phases completed!")
}

// Error handling with WaitGroup
func demonstrateErrorHandling() {
    fmt.Println("\n=== Error Handling with WaitGroup ===")
    
    var wg sync.WaitGroup
    errorChan := make(chan error, 5)
    
    // Tasks that might fail
    tasks := []struct {
        id      int
        name    string
        willFail bool
    }{
        {1, "Task 1", false},
        {2, "Task 2", true},
        {3, "Task 3", false},
        {4, "Task 4", true},
        {5, "Task 5", false},
    }
    
    for _, task := range tasks {
        wg.Add(1)
        
        go func(t struct {
            id      int
            name    string
            willFail bool
        }) {
            defer wg.Done()
            
            fmt.Printf("Starting %s\n", t.name)
            time.Sleep(200 * time.Millisecond)
            
            if t.willFail {
                err := fmt.Errorf("%s failed", t.name)
                fmt.Printf("Error in %s: %v\n", t.name, err)
                errorChan <- err
                return
            }
            
            fmt.Printf("%s completed successfully\n", t.name)
        }(task)
    }
    
    // Close error channel when all goroutines are done
    go func() {
        wg.Wait()
        close(errorChan)
    }()
    
    // Collect errors
    var errors []error
    for err := range errorChan {
        errors = append(errors, err)
    }
    
    fmt.Printf("\nCompleted with %d errors:\n", len(errors))
    for _, err := range errors {
        fmt.Printf("  - %v\n", err)
    }
}

func main() {
    demonstrateBasicWaitGroup()
    demonstrateDynamicWaitGroup()
    demonstrateNestedWaitGroups()
    demonstrateErrorHandling()
}
```

## Once

### Singleton Pattern with sync.Once

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Singleton using sync.Once
type Database struct {
    connection string
    connected  bool
}

var (
    dbInstance *Database
    dbOnce     sync.Once
)

// GetDatabase returns the singleton database instance
func GetDatabase() *Database {
    dbOnce.Do(func() {
        fmt.Println("Creating database instance...")
        time.Sleep(100 * time.Millisecond) // Simulate expensive initialization
        
        dbInstance = &Database{
            connection: "database://localhost:5432/mydb",
            connected:  true,
        }
        
        fmt.Println("Database instance created")
    })
    
    return dbInstance
}

// Configuration singleton
type Config struct {
    AppName    string
    Version    string
    Debug      bool
    MaxWorkers int
}

var (
    configInstance *Config
    configOnce     sync.Once
)

// LoadConfig loads configuration once
func LoadConfig() *Config {
    configOnce.Do(func() {
        fmt.Println("Loading configuration...")
        time.Sleep(50 * time.Millisecond) // Simulate file reading
        
        configInstance = &Config{
            AppName:    "MyApp",
            Version:    "1.0.0",
            Debug:      true,
            MaxWorkers: 10,
        }
        
        fmt.Println("Configuration loaded")
    })
    
    return configInstance
}

// Logger singleton
type Logger struct {
    level string
    file  string
}

var (
    loggerInstance *Logger
    loggerOnce     sync.Once
)

func GetLogger() *Logger {
    loggerOnce.Do(func() {
        fmt.Println("Initializing logger...")
        time.Sleep(30 * time.Millisecond)
        
        loggerInstance = &Logger{
            level: "INFO",
            file:  "/var/log/app.log",
        }
        
        fmt.Println("Logger initialized")
    })
    
    return loggerInstance
}

func (l *Logger) Log(message string) {
    fmt.Printf("[%s] %s\n", l.level, message)
}

func demonstrateSingletons() {
    fmt.Println("=== Singleton Pattern with sync.Once ===")
    
    var wg sync.WaitGroup
    
    // Multiple goroutines trying to get database instance
    for i := 1; i <= 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            fmt.Printf("Goroutine %d requesting database...\n", id)
            db := GetDatabase()
            fmt.Printf("Goroutine %d got database: %s\n", id, db.connection)
        }(i)
    }
    
    // Multiple goroutines trying to load config
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            fmt.Printf("Goroutine %d loading config...\n", id)
            config := LoadConfig()
            fmt.Printf("Goroutine %d got config: %s v%s\n", id, config.AppName, config.Version)
        }(i)
    }
    
    // Multiple goroutines trying to get logger
    for i := 1; i <= 4; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            fmt.Printf("Goroutine %d getting logger...\n", id)
            logger := GetLogger()
            logger.Log(fmt.Sprintf("Message from goroutine %d", id))
        }(i)
    }
    
    wg.Wait()
    
    // Verify all instances are the same
    fmt.Println("\nVerifying singleton instances:")
    db1 := GetDatabase()
    db2 := GetDatabase()
    fmt.Printf("Database instances are same: %t\n", db1 == db2)
    
    config1 := LoadConfig()
    config2 := LoadConfig()
    fmt.Printf("Config instances are same: %t\n", config1 == config2)
    
    logger1 := GetLogger()
    logger2 := GetLogger()
    fmt.Printf("Logger instances are same: %t\n", logger1 == logger2)
}

// One-time initialization example
func demonstrateOneTimeInit() {
    fmt.Println("\n=== One-time Initialization ===")
    
    var (
        expensiveResource string
        initOnce          sync.Once
    )
    
    initResource := func() {
        initOnce.Do(func() {
            fmt.Println("Performing expensive initialization...")
            time.Sleep(200 * time.Millisecond)
            expensiveResource = "Initialized expensive resource"
            fmt.Println("Expensive initialization completed")
        })
    }
    
    var wg sync.WaitGroup
    
    // Multiple goroutines that need the resource
    for i := 1; i <= 6; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            fmt.Printf("Worker %d needs resource...\n", id)
            initResource()
            fmt.Printf("Worker %d using resource: %s\n", id, expensiveResource)
            
            // Simulate work with the resource
            time.Sleep(50 * time.Millisecond)
            fmt.Printf("Worker %d finished\n", id)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All workers completed")
}

// Multiple sync.Once for different initialization phases
func demonstrateMultiPhaseInit() {
    fmt.Println("\n=== Multi-phase Initialization ===")
    
    var (
        phase1Done sync.Once
        phase2Done sync.Once
        phase3Done sync.Once
        
        phase1Data string
        phase2Data string
        phase3Data string
    )
    
    initPhase1 := func() {
        phase1Done.Do(func() {
            fmt.Println("Phase 1: Basic initialization")
            time.Sleep(100 * time.Millisecond)
            phase1Data = "Phase 1 complete"
        })
    }
    
    initPhase2 := func() {
        initPhase1() // Ensure phase 1 is done first
        phase2Done.Do(func() {
            fmt.Println("Phase 2: Advanced initialization")
            time.Sleep(150 * time.Millisecond)
            phase2Data = "Phase 2 complete"
        })
    }
    
    initPhase3 := func() {
        initPhase2() // Ensure phase 2 is done first
        phase3Done.Do(func() {
            fmt.Println("Phase 3: Final initialization")
            time.Sleep(100 * time.Millisecond)
            phase3Data = "Phase 3 complete"
        })
    }
    
    var wg sync.WaitGroup
    
    // Workers that need different phases
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            switch id {
            case 1:
                fmt.Printf("Worker %d needs phase 1\n", id)
                initPhase1()
                fmt.Printf("Worker %d: %s\n", id, phase1Data)
            case 2:
                fmt.Printf("Worker %d needs phase 2\n", id)
                initPhase2()
                fmt.Printf("Worker %d: %s, %s\n", id, phase1Data, phase2Data)
            case 3:
                fmt.Printf("Worker %d needs phase 3\n", id)
                initPhase3()
                fmt.Printf("Worker %d: %s, %s, %s\n", id, phase1Data, phase2Data, phase3Data)
            }
        }(i)
    }
    
    wg.Wait()
    fmt.Println("Multi-phase initialization completed")
}

func main() {
    demonstrateSingletons()
    demonstrateOneTimeInit()
    demonstrateMultiPhaseInit()
}
```

## Cond (Condition Variables)

### Basic Condition Variable Usage

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Queue with condition variable
type Queue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []interface{}
    maxSize int
}

// NewQueue creates a new queue with maximum size
func NewQueue(maxSize int) *Queue {
    q := &Queue{
        items:   make([]interface{}, 0),
        maxSize: maxSize,
    }
    q.cond = sync.NewCond(&q.mu)
    return q
}

// Enqueue adds an item to the queue (blocks if full)
func (q *Queue) Enqueue(item interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    // Wait while queue is full
    for len(q.items) >= q.maxSize {
        fmt.Printf("Queue full, waiting to enqueue %v\n", item)
        q.cond.Wait()
    }
    
    q.items = append(q.items, item)
    fmt.Printf("Enqueued %v (size: %d)\n", item, len(q.items))
    
    // Signal waiting dequeuers
    q.cond.Signal()
}

// Dequeue removes an item from the queue (blocks if empty)
func (q *Queue) Dequeue() interface{} {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    // Wait while queue is empty
    for len(q.items) == 0 {
        fmt.Println("Queue empty, waiting to dequeue")
        q.cond.Wait()
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    fmt.Printf("Dequeued %v (size: %d)\n", item, len(q.items))
    
    // Signal waiting enqueuers
    q.cond.Signal()
    
    return item
}

// TryEnqueue attempts to enqueue without blocking
func (q *Queue) TryEnqueue(item interface{}) bool {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if len(q.items) >= q.maxSize {
        return false
    }
    
    q.items = append(q.items, item)
    fmt.Printf("Try-enqueued %v (size: %d)\n", item, len(q.items))
    q.cond.Signal()
    return true
}

// TryDequeue attempts to dequeue without blocking
func (q *Queue) TryDequeue() (interface{}, bool) {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    if len(q.items) == 0 {
        return nil, false
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    fmt.Printf("Try-dequeued %v (size: %d)\n", item, len(q.items))
    q.cond.Signal()
    return item, true
}

// Size returns the current queue size
func (q *Queue) Size() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    return len(q.items)
}

// Broadcast wakes up all waiting goroutines
func (q *Queue) Broadcast() {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.cond.Broadcast()
}

func demonstrateBasicCond() {
    fmt.Println("=== Basic Condition Variable Usage ===")
    
    queue := NewQueue(3)
    var wg sync.WaitGroup
    
    // Producers
    for i := 1; i <= 2; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 1; j <= 5; j++ {
                item := fmt.Sprintf("P%d-Item%d", id, j)
                queue.Enqueue(item)
                time.Sleep(200 * time.Millisecond)
            }
            fmt.Printf("Producer %d finished\n", id)
        }(i)
    }
    
    // Consumers
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 1; j <= 3; j++ {
                item := queue.Dequeue()
                fmt.Printf("Consumer %d got: %v\n", id, item)
                time.Sleep(300 * time.Millisecond)
            }
            fmt.Printf("Consumer %d finished\n", id)
        }(i)
    }
    
    wg.Wait()
    fmt.Printf("Final queue size: %d\n", queue.Size())
}

// Worker pool with condition variable
type WorkerPool struct {
    mu       sync.Mutex
    cond     *sync.Cond
    tasks    []func()
    workers  int
    shutdown bool
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int) *WorkerPool {
    wp := &WorkerPool{
        tasks:   make([]func(), 0),
        workers: numWorkers,
    }
    wp.cond = sync.NewCond(&wp.mu)
    return wp
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
    for i := 1; i <= wp.workers; i++ {
        go wp.worker(i)
    }
}

// worker processes tasks
func (wp *WorkerPool) worker(id int) {
    fmt.Printf("Worker %d started\n", id)
    
    for {
        wp.mu.Lock()
        
        // Wait for tasks or shutdown
        for len(wp.tasks) == 0 && !wp.shutdown {
            fmt.Printf("Worker %d waiting for tasks\n", id)
            wp.cond.Wait()
        }
        
        // Check if shutting down
        if wp.shutdown {
            wp.mu.Unlock()
            fmt.Printf("Worker %d shutting down\n", id)
            return
        }
        
        // Get a task
        task := wp.tasks[0]
        wp.tasks = wp.tasks[1:]
        wp.mu.Unlock()
        
        fmt.Printf("Worker %d executing task\n", id)
        task()
        fmt.Printf("Worker %d completed task\n", id)
    }
}

// AddTask adds a task to the pool
func (wp *WorkerPool) AddTask(task func()) {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    if wp.shutdown {
        return
    }
    
    wp.tasks = append(wp.tasks, task)
    wp.cond.Signal() // Wake up one worker
}

// Shutdown shuts down the worker pool
func (wp *WorkerPool) Shutdown() {
    wp.mu.Lock()
    defer wp.mu.Unlock()
    
    wp.shutdown = true
    wp.cond.Broadcast() // Wake up all workers
}

func demonstrateWorkerPool() {
    fmt.Println("\n=== Worker Pool with Condition Variable ===")
    
    pool := NewWorkerPool(3)
    pool.Start()
    
    // Add tasks
    for i := 1; i <= 8; i++ {
        taskID := i
        pool.AddTask(func() {
            fmt.Printf("  Executing task %d\n", taskID)
            time.Sleep(200 * time.Millisecond)
        })
        time.Sleep(100 * time.Millisecond)
    }
    
    // Let workers process tasks
    time.Sleep(2 * time.Second)
    
    // Shutdown
    fmt.Println("Shutting down worker pool...")
    pool.Shutdown()
    
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Worker pool demonstration completed")
}

func main() {
    demonstrateBasicCond()
    demonstrateWorkerPool()
}
```

## Atomic Operations

### Basic Atomic Operations

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
    "time"
)

// Atomic counter
type AtomicCounter struct {
    value int64
}

// Increment atomically increments the counter
func (c *AtomicCounter) Increment() {
    atomic.AddInt64(&c.value, 1)
}

// Decrement atomically decrements the counter
func (c *AtomicCounter) Decrement() {
    atomic.AddInt64(&c.value, -1)
}

// Add atomically adds a value to the counter
func (c *AtomicCounter) Add(delta int64) {
    atomic.AddInt64(&c.value, delta)
}

// Value atomically returns the current value
func (c *AtomicCounter) Value() int64 {
    return atomic.LoadInt64(&c.value)
}

// Set atomically sets the value
func (c *AtomicCounter) Set(value int64) {
    atomic.StoreInt64(&c.value, value)
}

// CompareAndSwap atomically compares and swaps
func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
    return atomic.CompareAndSwapInt64(&c.value, old, new)
}

func demonstrateAtomicCounter() {
    fmt.Println("=== Atomic Counter ===")
    
    counter := &AtomicCounter{}
    var wg sync.WaitGroup
    
    // Multiple goroutines incrementing
    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 1000; j++ {
                counter.Increment()
                if j%250 == 0 {
                    fmt.Printf("Goroutine %d: Counter at %d\n", id, counter.Value())
                }
            }
        }(i)
    }
    
    // Multiple goroutines decrementing
    for i := 0; i < 2; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 500; j++ {
                counter.Decrement()
            }
        }(i + 5)
    }
    
    wg.Wait()
    fmt.Printf("Final atomic counter value: %d\n", counter.Value())
    fmt.Printf("Expected value: %d\n", 5*1000-2*500)
}

// Atomic flags and states
type AtomicFlags struct {
    running   int32
    connected int32
    ready     int32
}

// SetRunning atomically sets the running flag
func (f *AtomicFlags) SetRunning(running bool) {
    var value int32
    if running {
        value = 1
    }
    atomic.StoreInt32(&f.running, value)
}

// IsRunning atomically checks if running
func (f *AtomicFlags) IsRunning() bool {
    return atomic.LoadInt32(&f.running) == 1
}

// SetConnected atomically sets the connected flag
func (f *AtomicFlags) SetConnected(connected bool) {
    var value int32
    if connected {
        value = 1
    }
    atomic.StoreInt32(&f.connected, value)
}

// IsConnected atomically checks if connected
func (f *AtomicFlags) IsConnected() bool {
    return atomic.LoadInt32(&f.connected) == 1
}

// SetReady atomically sets the ready flag
func (f *AtomicFlags) SetReady(ready bool) {
    var value int32
    if ready {
        value = 1
    }
    atomic.StoreInt32(&f.ready, value)
}

// IsReady atomically checks if ready
func (f *AtomicFlags) IsReady() bool {
    return atomic.LoadInt32(&f.ready) == 1
}

// IsFullyReady checks if all flags are set
func (f *AtomicFlags) IsFullyReady() bool {
    return f.IsRunning() && f.IsConnected() && f.IsReady()
}

func demonstrateAtomicFlags() {
    fmt.Println("\n=== Atomic Flags ===")
    
    flags := &AtomicFlags{}
    var wg sync.WaitGroup
    
    // Goroutine that sets flags
    wg.Add(1)
    go func() {
        defer wg.Done()
        
        fmt.Println("Setting running flag...")
        flags.SetRunning(true)
        time.Sleep(200 * time.Millisecond)
        
        fmt.Println("Setting connected flag...")
        flags.SetConnected(true)
        time.Sleep(300 * time.Millisecond)
        
        fmt.Println("Setting ready flag...")
        flags.SetReady(true)
    }()
    
    // Goroutines that check flags
    for i := 1; i <= 3; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            for j := 0; j < 10; j++ {
                fmt.Printf("Checker %d: Running=%t, Connected=%t, Ready=%t, FullyReady=%t\n",
                    id, flags.IsRunning(), flags.IsConnected(), flags.IsReady(), flags.IsFullyReady())
                
                if flags.IsFullyReady() {
                    fmt.Printf("Checker %d: System is fully ready!\n", id)
                    return
                }
                
                time.Sleep(100 * time.Millisecond)
            }
        }(i)
    }
    
    wg.Wait()
}

// Atomic pointer operations
type Node struct {
    value int
    next  *Node
}

type AtomicStack struct {
    head unsafe.Pointer
}

func NewAtomicStack() *AtomicStack {
    return &AtomicStack{}
}

func (s *AtomicStack) Push(value int) {
    newNode := &Node{value: value}
    
    for {
        head := (*Node)(atomic.LoadPointer(&s.head))
        newNode.next = head
        
        if atomic.CompareAndSwapPointer(&s.head, unsafe.Pointer(head), unsafe.Pointer(newNode)) {
            break
        }
    }
}

func (s *AtomicStack) Pop() (int, bool) {
    for {
        head := (*Node)(atomic.LoadPointer(&s.head))
        if head == nil {
            return 0, false
        }
        
        if atomic.CompareAndSwapPointer(&s.head, unsafe.Pointer(head), unsafe.Pointer(head.next)) {
            return head.value, true
        }
    }
}

// Note: This example requires unsafe package
// import "unsafe"

// Performance comparison: Atomic vs Mutex
func compareAtomicVsMutex() {
    fmt.Println("\n=== Performance: Atomic vs Mutex ===")
    
    const numGoroutines = 10
    const numOperations = 100000
    
    // Test with atomic operations
    fmt.Println("Testing atomic operations...")
    var atomicCounter int64
    start := time.Now()
    
    var wg1 sync.WaitGroup
    for i := 0; i < numGoroutines; i++ {
        wg1.Add(1)
        go func() {
            defer wg1.Done()
            for j := 0; j < numOperations; j++ {
                atomic.AddInt64(&atomicCounter, 1)
            }
        }()
    }
    wg1.Wait()
    atomicTime := time.Since(start)
    
    // Test with mutex
    fmt.Println("Testing mutex operations...")
    var mutexCounter int64
    var mu sync.Mutex
    start = time.Now()
    
    var wg2 sync.WaitGroup
    for i := 0; i < numGoroutines; i++ {
        wg2.Add(1)
        go func() {
            defer wg2.Done()
            for j := 0; j < numOperations; j++ {
                mu.Lock()
                mutexCounter++
                mu.Unlock()
            }
        }()
    }
    wg2.Wait()
    mutexTime := time.Since(start)
    
    fmt.Printf("Atomic result: %d, time: %v\n", atomicCounter, atomicTime)
    fmt.Printf("Mutex result: %d, time: %v\n", mutexCounter, mutexTime)
    fmt.Printf("Atomic is %.2fx faster\n", float64(mutexTime)/float64(atomicTime))
}

func main() {
    demonstrateAtomicCounter()
    demonstrateAtomicFlags()
    compareAtomicVsMutex()
}
```

## Advanced Synchronization Patterns

### Barrier Synchronization

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

// Barrier allows multiple goroutines to wait for each other
type Barrier struct {
    mu       sync.Mutex
    cond     *sync.Cond
    count    int
    waiting  int
    generation int
}

// NewBarrier creates a new barrier for n goroutines
func NewBarrier(n int) *Barrier {
    b := &Barrier{
        count: n,
    }
    b.cond = sync.NewCond(&b.mu)
    return b
}

// Wait waits for all goroutines to reach the barrier
func (b *Barrier) Wait() {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    generation := b.generation
    b.waiting++
    
    if b.waiting == b.count {
        // Last goroutine to arrive
        b.waiting = 0
        b.generation++
        b.cond.Broadcast()
    } else {
        // Wait for others
        for generation == b.generation {
            b.cond.Wait()
        }
    }
}

func demonstrateBarrier() {
    fmt.Println("=== Barrier Synchronization ===")
    
    const numWorkers = 5
    barrier := NewBarrier(numWorkers)
    var wg sync.WaitGroup
    
    for i := 1; i <= numWorkers; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // Phase 1
            workTime := time.Duration(id*100) * time.Millisecond
            fmt.Printf("Worker %d: Phase 1 starting (will take %v)\n", id, workTime)
            time.Sleep(workTime)
            fmt.Printf("Worker %d: Phase 1 completed, waiting at barrier\n", id)
            
            barrier.Wait()
            
            // Phase 2
            fmt.Printf("Worker %d: Phase 2 starting\n", id)
            time.Sleep(200 * time.Millisecond)
            fmt.Printf("Worker %d: Phase 2 completed, waiting at barrier\n", id)
            
            barrier.Wait()
            
            // Phase 3
            fmt.Printf("Worker %d: Phase 3 starting\n", id)
            time.Sleep(150 * time.Millisecond)
            fmt.Printf("Worker %d: All phases completed\n", id)
        }(i)
    }
    
    wg.Wait()
    fmt.Println("All workers completed all phases")
}

// Semaphore implementation
type Semaphore struct {
    ch chan struct{}
}

// NewSemaphore creates a semaphore with given capacity
func NewSemaphore(capacity int) *Semaphore {
    return &Semaphore{
        ch: make(chan struct{}, capacity),
    }
}

// Acquire acquires a semaphore permit
func (s *Semaphore) Acquire() {
    s.ch <- struct{}{}
}

// Release releases a semaphore permit
func (s *Semaphore) Release() {
    <-s.ch
}

// TryAcquire tries to acquire a permit without blocking
func (s *Semaphore) TryAcquire() bool {
    select {
    case s.ch <- struct{}{}:
        return true
    default:
        return false
    }
}

// AvailablePermits returns the number of available permits
func (s *Semaphore) AvailablePermits() int {
    return cap(s.ch) - len(s.ch)
}

func demonstrateSemaphore() {
    fmt.Println("\n=== Semaphore Pattern ===")
    
    const maxConcurrent = 3
    const numTasks = 8
    
    semaphore := NewSemaphore(maxConcurrent)
    var wg sync.WaitGroup
    
    for i := 1; i <= numTasks; i++ {
        wg.Add(1)
        go func(taskID int) {
            defer wg.Done()
            
            fmt.Printf("Task %d: Waiting for permit (available: %d)\n", 
                taskID, semaphore.AvailablePermits())
            
            semaphore.Acquire()
            defer semaphore.Release()
            
            fmt.Printf("Task %d: Acquired permit, starting work (concurrent: %d)\n", 
                taskID, maxConcurrent-semaphore.AvailablePermits())
            
            // Simulate work
            time.Sleep(500 * time.Millisecond)
            
            fmt.Printf("Task %d: Work completed\n", taskID)
        }(i)
        
        time.Sleep(50 * time.Millisecond)
    }
    
    wg.Wait()
    fmt.Println("All tasks completed")
}

func main() {
    demonstrateBarrier()
    demonstrateSemaphore()
}
```

## Exercises

### Exercise 1: Thread-Safe Cache
Implement a thread-safe cache with TTL (Time To Live).

```go
// Create a thread-safe cache that:
// - Supports Get, Set, Delete operations
// - Has TTL (time-to-live) for entries
// - Automatically cleans up expired entries
// - Supports maximum size with LRU eviction
// - Provides cache statistics

package main

type CacheEntry struct {
    Value     interface{}
    ExpiresAt time.Time
    AccessCount int64
    LastAccess  time.Time
}

type TTLCache struct {
    // Add fields for thread-safe operations
}

// Implement these methods:
// NewTTLCache(maxSize int, defaultTTL time.Duration) *TTLCache
// Get(key string) (interface{}, bool)
// Set(key string, value interface{}, ttl time.Duration)
// Delete(key string) bool
// Clear()
// Size() int
// Stats() CacheStats
// StartCleanup(interval time.Duration)
// StopCleanup()
```

### Exercise 2: Rate Limiter
Implement different types of rate limiters.

```go
// Create rate limiters that:
// - Token bucket rate limiter
// - Sliding window rate limiter
// - Fixed window rate limiter
// - Support burst capacity
// - Provide rate limiting statistics

package main

type RateLimiter interface {
    Allow() bool
    Wait() error
    Reserve() Reservation
    Limit() Rate
    Burst() int
}

type TokenBucketLimiter struct {
    // Add fields for token bucket implementation
}

type SlidingWindowLimiter struct {
    // Add fields for sliding window implementation
}

// Implement the RateLimiter interface for both types
```

### Exercise 3: Producer-Consumer with Backpressure
Implement a producer-consumer system with backpressure handling.

```go
// Create a system that:
// - Handles multiple producers and consumers
// - Implements backpressure when consumers are slow
// - Supports different backpressure strategies (drop, block, buffer)
// - Provides monitoring and metrics
// - Supports graceful shutdown

package main

type BackpressureStrategy int

const (
    DropOldest BackpressureStrategy = iota
    DropNewest
    Block
    Buffer
)

type ProducerConsumerSystem struct {
    // Add fields for managing producers, consumers, and backpressure
}

// Implement these methods:
// NewProducerConsumerSystem(config Config) *ProducerConsumerSystem
// AddProducer(producer Producer) error
// AddConsumer(consumer Consumer) error
// Start() error
// Stop() error
// GetMetrics() SystemMetrics
```

## Key Takeaways

1. **Mutex protects shared state** from concurrent access
2. **RWMutex allows multiple readers** but exclusive writers
3. **WaitGroup coordinates goroutine completion** 
4. **sync.Once ensures one-time initialization** across goroutines
5. **Condition variables coordinate complex waiting** scenarios
6. **Atomic operations provide lock-free synchronization** for simple operations
7. **Choose the right primitive** based on your use case
8. **Avoid deadlocks** by consistent lock ordering
9. **Prefer channels for communication**, sync primitives for protection
10. **Always use defer** to ensure locks are released

## Next Steps

Now that you understand synchronization primitives, let's learn about [Context Package](18-context.md) for managing request lifecycles and cancellation!

---

**Previous**: [← Channels](16-channels.md) | **Next**: [Context Package →](18-context.md)