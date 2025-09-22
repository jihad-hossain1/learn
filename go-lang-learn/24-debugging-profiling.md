# 24. Debugging and Profiling

Go provides excellent tools for debugging and profiling applications. This chapter covers debugging techniques, performance profiling, memory analysis, and optimization strategies.

## Debugging Basics

### Print Debugging

```go
package main

import (
    "fmt"
    "log"
    "os"
)

// Simple print debugging
func calculateSum(numbers []int) int {
    fmt.Printf("DEBUG: Input numbers: %v\n", numbers)
    
    sum := 0
    for i, num := range numbers {
        sum += num
        fmt.Printf("DEBUG: Step %d: num=%d, sum=%d\n", i, num, sum)
    }
    
    fmt.Printf("DEBUG: Final sum: %d\n", sum)
    return sum
}

// Conditional debugging
var debugMode = os.Getenv("DEBUG") == "true"

func debugPrint(format string, args ...interface{}) {
    if debugMode {
        log.Printf("DEBUG: "+format, args...)
    }
}

func processData(data []string) []string {
    debugPrint("Processing %d items", len(data))
    
    var result []string
    for i, item := range data {
        debugPrint("Processing item %d: %s", i, item)
        
        if len(item) > 0 {
            processed := fmt.Sprintf("processed_%s", item)
            result = append(result, processed)
            debugPrint("Added to result: %s", processed)
        } else {
            debugPrint("Skipping empty item at index %d", i)
        }
    }
    
    debugPrint("Processing complete. Result count: %d", len(result))
    return result
}

func main() {
    // Example 1: Simple debugging
    numbers := []int{1, 2, 3, 4, 5}
    sum := calculateSum(numbers)
    fmt.Printf("Sum: %d\n", sum)
    
    fmt.Println("\n--- Conditional Debugging ---")
    
    // Example 2: Conditional debugging
    data := []string{"hello", "", "world", "test"}
    result := processData(data)
    fmt.Printf("Result: %v\n", result)
    
    fmt.Println("\nSet DEBUG=true environment variable to see debug output")
}
```

### Logging for Debugging

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"
)

// Custom logger with different levels
type Logger struct {
    debugLogger *log.Logger
    infoLogger  *log.Logger
    errorLogger *log.Logger
}

func NewLogger() *Logger {
    return &Logger{
        debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
        infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
        errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
    }
}

func (l *Logger) Debug(format string, args ...interface{}) {
    if os.Getenv("LOG_LEVEL") == "DEBUG" {
        l.debugLogger.Printf(format, args...)
    }
}

func (l *Logger) Info(format string, args ...interface{}) {
    l.infoLogger.Printf(format, args...)
}

func (l *Logger) Error(format string, args ...interface{}) {
    l.errorLogger.Printf(format, args...)
}

// Example service with logging
type UserService struct {
    logger *Logger
    users  map[int]string
}

func NewUserService() *UserService {
    return &UserService{
        logger: NewLogger(),
        users:  make(map[int]string),
    }
}

func (s *UserService) CreateUser(id int, name string) error {
    s.logger.Debug("CreateUser called with id=%d, name=%s", id, name)
    
    if id <= 0 {
        s.logger.Error("Invalid user ID: %d", id)
        return fmt.Errorf("user ID must be positive")
    }
    
    if name == "" {
        s.logger.Error("Empty user name for ID: %d", id)
        return fmt.Errorf("user name cannot be empty")
    }
    
    if _, exists := s.users[id]; exists {
        s.logger.Error("User with ID %d already exists", id)
        return fmt.Errorf("user with ID %d already exists", id)
    }
    
    s.users[id] = name
    s.logger.Info("User created successfully: id=%d, name=%s", id, name)
    s.logger.Debug("Current user count: %d", len(s.users))
    
    return nil
}

func (s *UserService) GetUser(id int) (string, error) {
    s.logger.Debug("GetUser called with id=%d", id)
    
    name, exists := s.users[id]
    if !exists {
        s.logger.Error("User not found: id=%d", id)
        return "", fmt.Errorf("user with ID %d not found", id)
    }
    
    s.logger.Debug("User found: id=%d, name=%s", id, name)
    return name, nil
}

func main() {
    service := NewUserService()
    
    // Test the service
    fmt.Println("=== Testing UserService ===")
    
    // Create users
    err := service.CreateUser(1, "John Doe")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    err = service.CreateUser(2, "Jane Smith")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // Try to create duplicate
    err = service.CreateUser(1, "Duplicate User")
    if err != nil {
        fmt.Printf("Expected error: %v\n", err)
    }
    
    // Get users
    name, err := service.GetUser(1)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("User 1: %s\n", name)
    }
    
    // Try to get non-existent user
    _, err = service.GetUser(999)
    if err != nil {
        fmt.Printf("Expected error: %v\n", err)
    }
    
    fmt.Println("\nSet LOG_LEVEL=DEBUG to see debug output")
}
```

## Using Delve Debugger

### Installing and Basic Usage

```bash
# Install Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a program
dlv debug main.go

# Debug with arguments
dlv debug main.go -- arg1 arg2

# Attach to running process
dlv attach <pid>

# Debug test
dlv test
```

### Delve Commands Example

```go
// debug_example.go
package main

import (
    "fmt"
    "time"
)

type Person struct {
    Name string
    Age  int
}

func (p *Person) Greet() string {
    return fmt.Sprintf("Hello, I'm %s and I'm %d years old", p.Name, p.Age)
}

func processNumbers(numbers []int) []int {
    var result []int
    
    for i, num := range numbers {
        // Set breakpoint here: (dlv) break main.processNumbers:25
        processed := num * 2
        
        if processed > 10 {
            result = append(result, processed)
        }
        
        fmt.Printf("Step %d: num=%d, processed=%d\n", i, num, processed)
    }
    
    return result
}

func main() {
    fmt.Println("Starting debug example...")
    
    // Set breakpoint here: (dlv) break main.main:40
    person := &Person{
        Name: "Alice",
        Age:  30,
    }
    
    greeting := person.Greet()
    fmt.Println(greeting)
    
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8}
    result := processNumbers(numbers)
    
    fmt.Printf("Result: %v\n", result)
    
    time.Sleep(1 * time.Second)
    fmt.Println("Program finished")
}
```

### Common Delve Commands

```bash
# Basic commands
(dlv) help                    # Show help
(dlv) break main.main         # Set breakpoint at function
(dlv) break main.go:25        # Set breakpoint at line
(dlv) breakpoints             # List breakpoints
(dlv) clear 1                 # Clear breakpoint by ID

# Execution control
(dlv) continue                # Continue execution
(dlv) next                    # Step over
(dlv) step                    # Step into
(dlv) stepout                 # Step out
(dlv) restart                 # Restart program

# Inspection
(dlv) print variable          # Print variable value
(dlv) locals                  # Show local variables
(dlv) args                    # Show function arguments
(dlv) vars                    # Show package variables
(dlv) stack                   # Show stack trace
(dlv) goroutines              # List goroutines
(dlv) goroutine 1             # Switch to goroutine

# Advanced
(dlv) condition 1 i > 5       # Set conditional breakpoint
(dlv) watch variable          # Watch variable changes
(dlv) disassemble             # Show assembly code
```

## Performance Profiling

### CPU Profiling

```go
package main

import (
    "fmt"
    "math/rand"
    "os"
    "runtime/pprof"
    "time"
)

// CPU-intensive function for profiling
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// Another CPU-intensive function
func bubbleSort(arr []int) []int {
    n := len(arr)
    result := make([]int, n)
    copy(result, arr)
    
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if result[j] > result[j+1] {
                result[j], result[j+1] = result[j+1], result[j]
            }
        }
    }
    
    return result
}

// Generate random data
func generateRandomData(size int) []int {
    data := make([]int, size)
    for i := range data {
        data[i] = rand.Intn(1000)
    }
    return data
}

func main() {
    // Start CPU profiling
    cpuFile, err := os.Create("cpu.prof")
    if err != nil {
        panic(err)
    }
    defer cpuFile.Close()
    
    if err := pprof.StartCPUProfile(cpuFile); err != nil {
        panic(err)
    }
    defer pprof.StopCPUProfile()
    
    fmt.Println("Starting CPU-intensive operations...")
    
    // CPU-intensive work
    start := time.Now()
    
    // Calculate fibonacci numbers
    for i := 0; i < 5; i++ {
        result := fibonacci(35)
        fmt.Printf("Fibonacci(35) = %d\n", result)
    }
    
    // Sort random data
    for i := 0; i < 10; i++ {
        data := generateRandomData(1000)
        sorted := bubbleSort(data)
        fmt.Printf("Sorted %d elements\n", len(sorted))
    }
    
    elapsed := time.Since(start)
    fmt.Printf("Total time: %v\n", elapsed)
    
    fmt.Println("CPU profile saved to cpu.prof")
    fmt.Println("Analyze with: go tool pprof cpu.prof")
}
```

### Memory Profiling

```go
package main

import (
    "fmt"
    "os"
    "runtime"
    "runtime/pprof"
    "time"
)

// Memory-intensive operations
type DataProcessor struct {
    data [][]byte
}

func NewDataProcessor() *DataProcessor {
    return &DataProcessor{
        data: make([][]byte, 0),
    }
}

func (dp *DataProcessor) ProcessLargeData() {
    // Allocate large chunks of memory
    for i := 0; i < 1000; i++ {
        // Allocate 1MB chunks
        chunk := make([]byte, 1024*1024)
        
        // Fill with some data
        for j := range chunk {
            chunk[j] = byte(i % 256)
        }
        
        dp.data = append(dp.data, chunk)
        
        if i%100 == 0 {
            fmt.Printf("Allocated %d chunks\n", i+1)
        }
    }
}

func (dp *DataProcessor) ProcessAndRelease() {
    // Process data and release some memory
    for i := 0; i < len(dp.data); i += 2 {
        // Process every other chunk
        chunk := dp.data[i]
        
        // Simulate processing
        sum := 0
        for _, b := range chunk {
            sum += int(b)
        }
        
        // Release the chunk
        dp.data[i] = nil
    }
    
    // Force garbage collection
    runtime.GC()
}

func createMemoryProfile() {
    memFile, err := os.Create("mem.prof")
    if err != nil {
        panic(err)
    }
    defer memFile.Close()
    
    runtime.GC() // Force GC before profiling
    if err := pprof.WriteHeapProfile(memFile); err != nil {
        panic(err)
    }
    
    fmt.Println("Memory profile saved to mem.prof")
}

func printMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Memory Stats:\n")
    fmt.Printf("  Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("  TotalAlloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("  Sys: %d KB\n", m.Sys/1024)
    fmt.Printf("  NumGC: %d\n", m.NumGC)
    fmt.Printf("  HeapObjects: %d\n", m.HeapObjects)
}

func main() {
    fmt.Println("Starting memory profiling example...")
    
    fmt.Println("\n=== Initial Memory Stats ===")
    printMemStats()
    
    processor := NewDataProcessor()
    
    fmt.Println("\n=== Processing Large Data ===")
    start := time.Now()
    processor.ProcessLargeData()
    elapsed := time.Since(start)
    
    fmt.Printf("Processing completed in %v\n", elapsed)
    
    fmt.Println("\n=== Memory Stats After Allocation ===")
    printMemStats()
    
    // Create memory profile
    createMemoryProfile()
    
    fmt.Println("\n=== Processing and Releasing Memory ===")
    processor.ProcessAndRelease()
    
    fmt.Println("\n=== Final Memory Stats ===")
    printMemStats()
    
    fmt.Println("\nAnalyze with: go tool pprof mem.prof")
}
```

### HTTP Server Profiling

```go
package main

import (
    "fmt"
    "log"
    "math/rand"
    "net/http"
    _ "net/http/pprof" // Import for side effects
    "runtime"
    "time"
)

// Simulate CPU-intensive work
func heavyComputation(n int) int {
    result := 0
    for i := 0; i < n*1000000; i++ {
        result += i
    }
    return result
}

// Simulate memory allocation
func allocateMemory(size int) []byte {
    data := make([]byte, size*1024) // size in KB
    for i := range data {
        data[i] = byte(rand.Intn(256))
    }
    return data
}

// HTTP handlers
func cpuIntensiveHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // Simulate CPU-intensive work
    result := heavyComputation(100)
    
    elapsed := time.Since(start)
    
    fmt.Fprintf(w, "CPU-intensive operation completed\n")
    fmt.Fprintf(w, "Result: %d\n", result)
    fmt.Fprintf(w, "Time taken: %v\n", elapsed)
}

func memoryIntensiveHandler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // Allocate memory
    var data [][]byte
    for i := 0; i < 100; i++ {
        chunk := allocateMemory(1024) // 1MB chunks
        data = append(data, chunk)
    }
    
    elapsed := time.Since(start)
    
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Fprintf(w, "Memory allocation completed\n")
    fmt.Fprintf(w, "Allocated %d chunks\n", len(data))
    fmt.Fprintf(w, "Time taken: %v\n", elapsed)
    fmt.Fprintf(w, "Current memory usage: %d KB\n", m.Alloc/1024)
    
    // Keep reference to prevent GC
    _ = data
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Fprintf(w, "Server Status\n")
    fmt.Fprintf(w, "=============\n")
    fmt.Fprintf(w, "Goroutines: %d\n", runtime.NumGoroutine())
    fmt.Fprintf(w, "Memory Alloc: %d KB\n", m.Alloc/1024)
    fmt.Fprintf(w, "Total Alloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Fprintf(w, "Sys: %d KB\n", m.Sys/1024)
    fmt.Fprintf(w, "GC Cycles: %d\n", m.NumGC)
}

func main() {
    // Register handlers
    http.HandleFunc("/cpu", cpuIntensiveHandler)
    http.HandleFunc("/memory", memoryIntensiveHandler)
    http.HandleFunc("/status", statusHandler)
    
    // Default handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Profiling Server\n")
        fmt.Fprintf(w, "===============\n")
        fmt.Fprintf(w, "Available endpoints:\n")
        fmt.Fprintf(w, "  /cpu     - CPU-intensive operation\n")
        fmt.Fprintf(w, "  /memory  - Memory-intensive operation\n")
        fmt.Fprintf(w, "  /status  - Server status\n")
        fmt.Fprintf(w, "\nProfiling endpoints:\n")
        fmt.Fprintf(w, "  /debug/pprof/        - Profile index\n")
        fmt.Fprintf(w, "  /debug/pprof/profile - CPU profile\n")
        fmt.Fprintf(w, "  /debug/pprof/heap    - Heap profile\n")
        fmt.Fprintf(w, "  /debug/pprof/goroutine - Goroutine profile\n")
    })
    
    fmt.Println("Starting profiling server on :8080")
    fmt.Println("Visit http://localhost:8080 for available endpoints")
    fmt.Println("\nProfiling commands:")
    fmt.Println("  go tool pprof http://localhost:8080/debug/pprof/profile")
    fmt.Println("  go tool pprof http://localhost:8080/debug/pprof/heap")
    fmt.Println("  go tool pprof http://localhost:8080/debug/pprof/goroutine")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## Analyzing Profiles

### Using go tool pprof

```bash
# Analyze CPU profile
go tool pprof cpu.prof

# Analyze memory profile
go tool pprof mem.prof

# Analyze live server
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Common pprof commands
(pprof) top           # Show top functions
(pprof) top10         # Show top 10 functions
(pprof) list main     # Show source code for main
(pprof) web           # Open web interface
(pprof) png           # Generate PNG graph
(pprof) pdf           # Generate PDF report
(pprof) help          # Show help
```

### Benchmark Profiling

```go
package main

import (
    "strings"
    "testing"
)

// Functions to benchmark
func ConcatenateLoop(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s
    }
    return result
}

func ConcatenateBuilder(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

func ConcatenateJoin(strs []string) string {
    return strings.Join(strs, "")
}

// Benchmark tests
func BenchmarkConcatenateLoop(b *testing.B) {
    strs := []string{"hello", "world", "this", "is", "a", "test", "of", "string", "concatenation"}
    
    for i := 0; i < b.N; i++ {
        ConcatenateLoop(strs)
    }
}

func BenchmarkConcatenateBuilder(b *testing.B) {
    strs := []string{"hello", "world", "this", "is", "a", "test", "of", "string", "concatenation"}
    
    for i := 0; i < b.N; i++ {
        ConcatenateBuilder(strs)
    }
}

func BenchmarkConcatenateJoin(b *testing.B) {
    strs := []string{"hello", "world", "this", "is", "a", "test", "of", "string", "concatenation"}
    
    for i := 0; i < b.N; i++ {
        ConcatenateJoin(strs)
    }
}

// Benchmark with different sizes
func BenchmarkConcatenateSizes(b *testing.B) {
    sizes := []int{10, 100, 1000}
    
    for _, size := range sizes {
        strs := make([]string, size)
        for i := range strs {
            strs[i] = "test"
        }
        
        b.Run(fmt.Sprintf("Loop-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ConcatenateLoop(strs)
            }
        })
        
        b.Run(fmt.Sprintf("Builder-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ConcatenateBuilder(strs)
            }
        })
        
        b.Run(fmt.Sprintf("Join-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                ConcatenateJoin(strs)
            }
        })
    }
}
```

### Running Benchmark Profiles

```bash
# Run benchmarks with CPU profiling
go test -bench=. -cpuprofile=cpu.prof

# Run benchmarks with memory profiling
go test -bench=. -memprofile=mem.prof

# Run benchmarks with both
go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof

# Analyze the profiles
go tool pprof cpu.prof
go tool pprof mem.prof
```

## Memory Analysis

### Memory Leak Detection

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

// Example of potential memory leak
type LeakyService struct {
    data map[string][]byte
}

func NewLeakyService() *LeakyService {
    return &LeakyService{
        data: make(map[string][]byte),
    }
}

// This method has a memory leak - data is never cleaned up
func (s *LeakyService) ProcessData(key string, size int) {
    // Allocate memory but never clean it up
    s.data[key] = make([]byte, size*1024) // size in KB
    
    // Fill with some data
    for i := range s.data[key] {
        s.data[key][i] = byte(i % 256)
    }
}

// Fixed version with cleanup
type FixedService struct {
    data    map[string][]byte
    maxSize int
}

func NewFixedService(maxSize int) *FixedService {
    return &FixedService{
        data:    make(map[string][]byte),
        maxSize: maxSize,
    }
}

func (s *FixedService) ProcessData(key string, size int) {
    // Clean up old data if we're at capacity
    if len(s.data) >= s.maxSize {
        // Remove oldest entries (simple FIFO)
        count := 0
        for k := range s.data {
            delete(s.data, k)
            count++
            if count >= s.maxSize/2 {
                break
            }
        }
        runtime.GC() // Force garbage collection
    }
    
    s.data[key] = make([]byte, size*1024)
    for i := range s.data[key] {
        s.data[key][i] = byte(i % 256)
    }
}

func printMemoryStats(label string) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("%s:\n", label)
    fmt.Printf("  Alloc: %d KB\n", m.Alloc/1024)
    fmt.Printf("  TotalAlloc: %d KB\n", m.TotalAlloc/1024)
    fmt.Printf("  Sys: %d KB\n", m.Sys/1024)
    fmt.Printf("  NumGC: %d\n", m.NumGC)
    fmt.Printf("  HeapObjects: %d\n\n", m.HeapObjects)
}

func demonstrateMemoryLeak() {
    fmt.Println("=== Demonstrating Memory Leak ===")
    
    service := NewLeakyService()
    
    printMemoryStats("Initial")
    
    // Simulate continuous data processing
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("data_%d", i)
        service.ProcessData(key, 1024) // 1MB each
        
        if i%100 == 0 {
            printMemoryStats(fmt.Sprintf("After %d iterations", i+1))
        }
    }
    
    fmt.Println("Notice how memory keeps growing!")
}

func demonstrateFixedVersion() {
    fmt.Println("\n=== Demonstrating Fixed Version ===")
    
    service := NewFixedService(50) // Limit to 50 entries
    
    printMemoryStats("Initial")
    
    // Simulate continuous data processing
    for i := 0; i < 1000; i++ {
        key := fmt.Sprintf("data_%d", i)
        service.ProcessData(key, 1024) // 1MB each
        
        if i%100 == 0 {
            printMemoryStats(fmt.Sprintf("After %d iterations", i+1))
        }
    }
    
    fmt.Println("Memory usage is now controlled!")
}

func main() {
    demonstrateMemoryLeak()
    
    // Force GC to clean up
    runtime.GC()
    time.Sleep(100 * time.Millisecond)
    
    demonstrateFixedVersion()
}
```

### Goroutine Leak Detection

```go
package main

import (
    "context"
    "fmt"
    "runtime"
    "sync"
    "time"
)

// Example of goroutine leak
func leakyWorker(data <-chan string) {
    for {
        select {
        case msg := <-data:
            fmt.Printf("Processing: %s\n", msg)
            time.Sleep(100 * time.Millisecond)
        }
        // No case for context cancellation - this will leak!
    }
}

// Fixed version with proper cancellation
func fixedWorker(ctx context.Context, data <-chan string) {
    for {
        select {
        case msg := <-data:
            fmt.Printf("Processing: %s\n", msg)
            time.Sleep(100 * time.Millisecond)
        case <-ctx.Done():
            fmt.Println("Worker shutting down")
            return
        }
    }
}

func printGoroutineCount(label string) {
    count := runtime.NumGoroutine()
    fmt.Printf("%s: %d goroutines\n", label, count)
}

func demonstrateGoroutineLeak() {
    fmt.Println("=== Demonstrating Goroutine Leak ===")
    
    printGoroutineCount("Initial")
    
    // Start leaky workers
    dataChan := make(chan string, 10)
    
    for i := 0; i < 10; i++ {
        go leakyWorker(dataChan)
    }
    
    printGoroutineCount("After starting workers")
    
    // Send some data
    for i := 0; i < 5; i++ {
        dataChan <- fmt.Sprintf("message_%d", i)
    }
    
    time.Sleep(1 * time.Second)
    
    // Close channel and try to clean up
    close(dataChan)
    time.Sleep(1 * time.Second)
    
    printGoroutineCount("After closing channel")
    fmt.Println("Goroutines are still running - they leaked!")
}

func demonstrateFixedVersion() {
    fmt.Println("\n=== Demonstrating Fixed Version ===")
    
    printGoroutineCount("Initial")
    
    ctx, cancel := context.WithCancel(context.Background())
    dataChan := make(chan string, 10)
    
    var wg sync.WaitGroup
    
    // Start fixed workers
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fixedWorker(ctx, dataChan)
        }()
    }
    
    printGoroutineCount("After starting workers")
    
    // Send some data
    for i := 0; i < 5; i++ {
        dataChan <- fmt.Sprintf("message_%d", i)
    }
    
    time.Sleep(1 * time.Second)
    
    // Properly shut down
    cancel()
    close(dataChan)
    
    // Wait for all workers to finish
    wg.Wait()
    
    time.Sleep(100 * time.Millisecond)
    printGoroutineCount("After proper shutdown")
    fmt.Println("All goroutines cleaned up properly!")
}

func main() {
    demonstrateGoroutineLeak()
    
    // Give some time for cleanup attempts
    time.Sleep(2 * time.Second)
    
    demonstrateFixedVersion()
}
```

## Performance Optimization

### String Operations Optimization

```go
package main

import (
    "fmt"
    "strings"
    "testing"
)

// Inefficient string concatenation
func inefficientConcat(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s // Creates new string each time
    }
    return result
}

// Efficient string concatenation
func efficientConcat(strs []string) string {
    var builder strings.Builder
    builder.Grow(len(strs) * 10) // Pre-allocate capacity
    
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

// Most efficient for simple joining
func mostEfficientConcat(strs []string) string {
    return strings.Join(strs, "")
}

// Benchmark functions
func BenchmarkInefficientConcat(b *testing.B) {
    strs := make([]string, 100)
    for i := range strs {
        strs[i] = fmt.Sprintf("string%d", i)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        inefficientConcat(strs)
    }
}

func BenchmarkEfficientConcat(b *testing.B) {
    strs := make([]string, 100)
    for i := range strs {
        strs[i] = fmt.Sprintf("string%d", i)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        efficientConcat(strs)
    }
}

func BenchmarkMostEfficientConcat(b *testing.B) {
    strs := make([]string, 100)
    for i := range strs {
        strs[i] = fmt.Sprintf("string%d", i)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        mostEfficientConcat(strs)
    }
}
```

### Slice Operations Optimization

```go
package main

import (
    "fmt"
    "testing"
)

// Inefficient slice operations
func inefficientSliceOps(data []int) []int {
    var result []int
    
    for _, v := range data {
        if v%2 == 0 {
            result = append(result, v*2) // Frequent reallocations
        }
    }
    
    return result
}

// Efficient slice operations
func efficientSliceOps(data []int) []int {
    // Pre-allocate with estimated capacity
    result := make([]int, 0, len(data)/2)
    
    for _, v := range data {
        if v%2 == 0 {
            result = append(result, v*2)
        }
    }
    
    return result
}

// Most efficient - count first, then allocate
func mostEfficientSliceOps(data []int) []int {
    // Count even numbers first
    count := 0
    for _, v := range data {
        if v%2 == 0 {
            count++
        }
    }
    
    // Allocate exact size
    result := make([]int, 0, count)
    
    for _, v := range data {
        if v%2 == 0 {
            result = append(result, v*2)
        }
    }
    
    return result
}

// Benchmark functions
func BenchmarkInefficientSliceOps(b *testing.B) {
    data := make([]int, 10000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        inefficientSliceOps(data)
    }
}

func BenchmarkEfficientSliceOps(b *testing.B) {
    data := make([]int, 10000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        efficientSliceOps(data)
    }
}

func BenchmarkMostEfficientSliceOps(b *testing.B) {
    data := make([]int, 10000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        mostEfficientSliceOps(data)
    }
}
```

## Debugging Best Practices

### Structured Logging

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "time"
)

// Structured log entry
type LogEntry struct {
    Timestamp string                 `json:"timestamp"`
    Level     string                 `json:"level"`
    Message   string                 `json:"message"`
    Fields    map[string]interface{} `json:"fields,omitempty"`
    Error     string                 `json:"error,omitempty"`
}

// Structured logger
type StructuredLogger struct {
    logger *log.Logger
}

func NewStructuredLogger() *StructuredLogger {
    return &StructuredLogger{
        logger: log.New(os.Stdout, "", 0),
    }
}

func (sl *StructuredLogger) log(level, message string, fields map[string]interface{}, err error) {
    entry := LogEntry{
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Level:     level,
        Message:   message,
        Fields:    fields,
    }
    
    if err != nil {
        entry.Error = err.Error()
    }
    
    jsonData, _ := json.Marshal(entry)
    sl.logger.Println(string(jsonData))
}

func (sl *StructuredLogger) Info(message string, fields map[string]interface{}) {
    sl.log("INFO", message, fields, nil)
}

func (sl *StructuredLogger) Error(message string, err error, fields map[string]interface{}) {
    sl.log("ERROR", message, fields, err)
}

func (sl *StructuredLogger) Debug(message string, fields map[string]interface{}) {
    if os.Getenv("LOG_LEVEL") == "DEBUG" {
        sl.log("DEBUG", message, fields, nil)
    }
}

// Example usage
func processOrder(orderID string, userID int, amount float64) error {
    logger := NewStructuredLogger()
    
    logger.Info("Processing order", map[string]interface{}{
        "order_id": orderID,
        "user_id":  userID,
        "amount":   amount,
    })
    
    // Simulate processing
    if amount <= 0 {
        err := fmt.Errorf("invalid amount")
        logger.Error("Order processing failed", err, map[string]interface{}{
            "order_id": orderID,
            "amount":   amount,
            "reason":   "invalid_amount",
        })
        return err
    }
    
    // Simulate some processing steps
    logger.Debug("Validating payment method", map[string]interface{}{
        "order_id": orderID,
        "step":     "payment_validation",
    })
    
    time.Sleep(100 * time.Millisecond)
    
    logger.Debug("Updating inventory", map[string]interface{}{
        "order_id": orderID,
        "step":     "inventory_update",
    })
    
    time.Sleep(50 * time.Millisecond)
    
    logger.Info("Order processed successfully", map[string]interface{}{
        "order_id":      orderID,
        "user_id":       userID,
        "amount":        amount,
        "processing_ms": 150,
    })
    
    return nil
}

func main() {
    fmt.Println("=== Structured Logging Example ===")
    
    // Process some orders
    processOrder("ORD-001", 123, 99.99)
    processOrder("ORD-002", 456, -10.00) // This will fail
    processOrder("ORD-003", 789, 149.99)
    
    fmt.Println("\nSet LOG_LEVEL=DEBUG to see debug messages")
}
```

## Exercises

### Exercise 1: Debug a Concurrent Program
Create a program with multiple goroutines that has race conditions and deadlocks. Use debugging tools to identify and fix the issues.

### Exercise 2: Profile and Optimize
Write a program that processes large amounts of data inefficiently. Use profiling tools to identify bottlenecks and optimize the code.

### Exercise 3: Memory Leak Detection
Create a web server that has memory leaks. Use memory profiling to identify the leaks and fix them.

## Key Takeaways

- Use print debugging for simple issues, structured logging for complex applications
- Delve debugger provides powerful debugging capabilities for Go programs
- CPU profiling helps identify performance bottlenecks
- Memory profiling helps detect memory leaks and optimize memory usage
- The `net/http/pprof` package enables profiling of running web servers
- Regular profiling and benchmarking help maintain application performance
- Proper resource management prevents memory and goroutine leaks
- Structured logging improves debugging in production environments
- Use appropriate data structures and algorithms for better performance

## Next Steps

Next, we'll explore [Deployment](25-deployment.md) to learn about building, packaging, and deploying Go applications.