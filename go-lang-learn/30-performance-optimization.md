# 27. Performance Optimization

This chapter covers techniques for optimizing Go applications, including profiling, memory management, CPU optimization, and best practices for writing high-performance Go code.

## Profiling and Benchmarking

### CPU Profiling

```go
// profile/cpu.go
package main

import (
    "fmt"
    "log"
    "os"
    "runtime/pprof"
    "time"
)

func expensiveOperation() {
    // Simulate CPU-intensive work
    for i := 0; i < 1000000; i++ {
        for j := 0; j < 1000; j++ {
            _ = i * j
        }
    }
}

func main() {
    // Create CPU profile file
    f, err := os.Create("cpu.prof")
    if err != nil {
        log.Fatal("could not create CPU profile: ", err)
    }
    defer f.Close()
    
    // Start CPU profiling
    if err := pprof.StartCPUProfile(f); err != nil {
        log.Fatal("could not start CPU profile: ", err)
    }
    defer pprof.StopCPUProfile()
    
    // Run the code to profile
    start := time.Now()
    expensiveOperation()
    duration := time.Since(start)
    
    fmt.Printf("Operation took: %v\n", duration)
}

// To analyze the profile:
// go run profile/cpu.go
// go tool pprof cpu.prof
// (pprof) top
// (pprof) list expensiveOperation
// (pprof) web
```

### Memory Profiling

```go
// profile/memory.go
package main

import (
    "fmt"
    "log"
    "os"
    "runtime"
    "runtime/pprof"
)

func allocateMemory() {
    // Simulate memory allocation
    var slices [][]int
    for i := 0; i < 1000; i++ {
        slice := make([]int, 1000)
        for j := range slice {
            slice[j] = j
        }
        slices = append(slices, slice)
    }
    
    // Keep reference to prevent GC
    _ = slices
}

func main() {
    // Run the code that allocates memory
    allocateMemory()
    
    // Force garbage collection
    runtime.GC()
    
    // Create memory profile file
    f, err := os.Create("mem.prof")
    if err != nil {
        log.Fatal("could not create memory profile: ", err)
    }
    defer f.Close()
    
    // Write heap profile
    if err := pprof.WriteHeapProfile(f); err != nil {
        log.Fatal("could not write memory profile: ", err)
    }
    
    fmt.Println("Memory profile written to mem.prof")
}

// To analyze the profile:
// go run profile/memory.go
// go tool pprof mem.prof
// (pprof) top
// (pprof) list allocateMemory
```

### HTTP Server Profiling

```go
// profile/server.go
package main

import (
    "fmt"
    "log"
    "net/http"
    _ "net/http/pprof" // Import for side effects
    "time"
)

func heavyHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate heavy computation
    start := time.Now()
    
    var result int
    for i := 0; i < 1000000; i++ {
        result += i * i
    }
    
    duration := time.Since(start)
    fmt.Fprintf(w, "Result: %d, Time: %v\n", result, duration)
}

func memoryHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate memory allocation
    data := make([][]byte, 1000)
    for i := range data {
        data[i] = make([]byte, 1024)
    }
    
    fmt.Fprintf(w, "Allocated %d KB\n", len(data))
}

func main() {
    http.HandleFunc("/heavy", heavyHandler)
    http.HandleFunc("/memory", memoryHandler)
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Profiling endpoints:")
    fmt.Println("  http://localhost:8080/debug/pprof/")
    fmt.Println("  http://localhost:8080/debug/pprof/profile")
    fmt.Println("  http://localhost:8080/debug/pprof/heap")
    fmt.Println("  http://localhost:8080/debug/pprof/goroutine")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Usage:
// go run profile/server.go
// curl http://localhost:8080/heavy
// go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30
// go tool pprof http://localhost:8080/debug/pprof/heap
```

### Benchmarking

```go
// benchmark/string_test.go
package benchmark

import (
    "fmt"
    "strings"
    "testing"
)

// Different string concatenation methods
func BenchmarkStringConcatPlus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var result string
        for j := 0; j < 100; j++ {
            result += "hello"
        }
    }
}

func BenchmarkStringConcatBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var builder strings.Builder
        for j := 0; j < 100; j++ {
            builder.WriteString("hello")
        }
        _ = builder.String()
    }
}

func BenchmarkStringConcatJoin(b *testing.B) {
    for i := 0; i < b.N; i++ {
        parts := make([]string, 100)
        for j := 0; j < 100; j++ {
            parts[j] = "hello"
        }
        _ = strings.Join(parts, "")
    }
}

func BenchmarkStringConcatSprintf(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var result string
        for j := 0; j < 100; j++ {
            result = fmt.Sprintf("%s%s", result, "hello")
        }
    }
}

// Memory allocation benchmarks
func BenchmarkSliceAppend(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var slice []int
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}

func BenchmarkSlicePrealloc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        slice := make([]int, 0, 1000)
        for j := 0; j < 1000; j++ {
            slice = append(slice, j)
        }
    }
}

func BenchmarkSliceIndex(b *testing.B) {
    for i := 0; i < b.N; i++ {
        slice := make([]int, 1000)
        for j := 0; j < 1000; j++ {
            slice[j] = j
        }
    }
}

// Run benchmarks:
// go test -bench=. -benchmem
// go test -bench=BenchmarkStringConcat -benchmem
// go test -bench=. -benchmem -cpuprofile=cpu.prof
// go test -bench=. -benchmem -memprofile=mem.prof
```

### Custom Benchmarking

```go
// benchmark/custom_test.go
package benchmark

import (
    "testing"
    "time"
)

// Benchmark with setup and teardown
func BenchmarkWithSetup(b *testing.B) {
    // Setup (not measured)
    data := make([]int, 10000)
    for i := range data {
        data[i] = i
    }
    
    b.ResetTimer() // Reset timer after setup
    
    for i := 0; i < b.N; i++ {
        // Code to benchmark
        sum := 0
        for _, v := range data {
            sum += v
        }
    }
}

// Benchmark with sub-benchmarks
func BenchmarkSorting(b *testing.B) {
    sizes := []int{100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
            for i := 0; i < b.N; i++ {
                b.StopTimer()
                data := generateRandomSlice(size)
                b.StartTimer()
                
                sort.Ints(data)
            }
        })
    }
}

// Benchmark with parallel execution
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            // Code to benchmark in parallel
            expensiveFunction()
        }
    })
}

// Custom timing
func BenchmarkCustomTiming(b *testing.B) {
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        // Setup for each iteration
        data := setupData()
        b.StartTimer()
        
        // Measured code
        processData(data)
        
        b.StopTimer()
        // Cleanup (not measured)
        cleanup(data)
        b.StartTimer()
    }
}

func generateRandomSlice(size int) []int {
    slice := make([]int, size)
    for i := range slice {
        slice[i] = rand.Intn(1000)
    }
    return slice
}

func expensiveFunction() {
    time.Sleep(time.Microsecond)
}

func setupData() interface{} {
    return make([]int, 1000)
}

func processData(data interface{}) {
    // Process data
}

func cleanup(data interface{}) {
    // Cleanup
}
```

## Memory Optimization

### Memory Pool Pattern

```go
// pool/buffer.go
package pool

import (
    "bytes"
    "sync"
)

// Buffer pool for reusing byte buffers
var bufferPool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func GetBuffer() *bytes.Buffer {
    return bufferPool.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
    buf.Reset() // Clear the buffer
    bufferPool.Put(buf)
}

// Example usage
func ProcessData(data []byte) []byte {
    buf := GetBuffer()
    defer PutBuffer(buf)
    
    // Use the buffer
    buf.Write(data)
    buf.WriteString(" processed")
    
    // Return a copy since we're returning the buffer to the pool
    result := make([]byte, buf.Len())
    copy(result, buf.Bytes())
    return result
}

// Custom object pool
type Worker struct {
    ID   int
    Data []byte
}

func (w *Worker) Reset() {
    w.ID = 0
    w.Data = w.Data[:0] // Keep capacity, reset length
}

type WorkerPool struct {
    pool sync.Pool
}

func NewWorkerPool() *WorkerPool {
    return &WorkerPool{
        pool: sync.Pool{
            New: func() interface{} {
                return &Worker{
                    Data: make([]byte, 0, 1024), // Pre-allocate capacity
                }
            },
        },
    }
}

func (wp *WorkerPool) Get() *Worker {
    return wp.pool.Get().(*Worker)
}

func (wp *WorkerPool) Put(w *Worker) {
    w.Reset()
    wp.pool.Put(w)
}

// Example usage
func (wp *WorkerPool) ProcessJob(jobID int, data []byte) {
    worker := wp.Get()
    defer wp.Put(worker)
    
    worker.ID = jobID
    worker.Data = append(worker.Data, data...)
    
    // Process the job
    // ...
}
```

### Slice Optimization

```go
// optimization/slice.go
package optimization

import "fmt"

// Inefficient: multiple allocations
func AppendInefficient(data []int) []int {
    var result []int
    for _, v := range data {
        result = append(result, v*2)
    }
    return result
}

// Efficient: pre-allocate capacity
func AppendEfficient(data []int) []int {
    result := make([]int, 0, len(data)) // Pre-allocate capacity
    for _, v := range data {
        result = append(result, v*2)
    }
    return result
}

// Most efficient: pre-allocate length
func AppendMostEfficient(data []int) []int {
    result := make([]int, len(data)) // Pre-allocate length
    for i, v := range data {
        result[i] = v * 2
    }
    return result
}

// Slice filtering without allocations
func FilterInPlace(slice []int, predicate func(int) bool) []int {
    n := 0
    for _, v := range slice {
        if predicate(v) {
            slice[n] = v
            n++
        }
    }
    return slice[:n]
}

// Slice copying optimization
func CopySliceEfficient(src []int) []int {
    dst := make([]int, len(src))
    copy(dst, src) // More efficient than manual loop
    return dst
}

// Avoiding slice header copies
func ProcessSliceByPointer(slice *[]int) {
    for i := range *slice {
        (*slice)[i] *= 2
    }
}

// String to byte slice optimization
func StringToBytesUnsafe(s string) []byte {
    // WARNING: This is unsafe and should only be used when
    // you're sure the byte slice won't be modified
    return *(*[]byte)(unsafe.Pointer(&s))
}

// Safe version
func StringToBytesSafe(s string) []byte {
    return []byte(s)
}

// Reusing slices
type SliceReuser struct {
    buffer []int
}

func NewSliceReuser(capacity int) *SliceReuser {
    return &SliceReuser{
        buffer: make([]int, 0, capacity),
    }
}

func (sr *SliceReuser) Process(data []int) []int {
    // Reset length but keep capacity
    sr.buffer = sr.buffer[:0]
    
    // Ensure capacity
    if cap(sr.buffer) < len(data) {
        sr.buffer = make([]int, 0, len(data))
    }
    
    for _, v := range data {
        sr.buffer = append(sr.buffer, v*2)
    }
    
    // Return a copy to avoid sharing the underlying array
    result := make([]int, len(sr.buffer))
    copy(result, sr.buffer)
    return result
}
```

### String Optimization

```go
// optimization/string.go
package optimization

import (
    "strings"
    "unsafe"
)

// Efficient string building
func BuildStringEfficient(parts []string) string {
    var builder strings.Builder
    
    // Pre-calculate total size
    totalSize := 0
    for _, part := range parts {
        totalSize += len(part)
    }
    
    builder.Grow(totalSize) // Pre-allocate
    
    for _, part := range parts {
        builder.WriteString(part)
    }
    
    return builder.String()
}

// Byte slice to string without allocation (unsafe)
func BytesToStringUnsafe(b []byte) string {
    // WARNING: This is unsafe and should only be used when
    // you're sure the byte slice won't be modified
    return *(*string)(unsafe.Pointer(&b))
}

// Safe version
func BytesToStringSafe(b []byte) string {
    return string(b)
}

// String concatenation comparison
func ConcatWithPlus(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s // Inefficient: creates new string each time
    }
    return result
}

func ConcatWithBuilder(strs []string) string {
    var builder strings.Builder
    for _, s := range strs {
        builder.WriteString(s)
    }
    return builder.String()
}

func ConcatWithJoin(strs []string) string {
    return strings.Join(strs, "")
}

// String interning for memory savings
type StringInterner struct {
    cache map[string]string
}

func NewStringInterner() *StringInterner {
    return &StringInterner{
        cache: make(map[string]string),
    }
}

func (si *StringInterner) Intern(s string) string {
    if interned, exists := si.cache[s]; exists {
        return interned
    }
    
    // Make a copy to ensure we own the string
    interned := string([]byte(s))
    si.cache[s] = interned
    return interned
}

// Efficient string searching
func ContainsAny(s string, chars string) bool {
    for _, c := range chars {
        if strings.ContainsRune(s, c) {
            return true
        }
    }
    return false
}

// More efficient version using strings.ContainsAny
func ContainsAnyEfficient(s string, chars string) bool {
    return strings.ContainsAny(s, chars)
}
```

## CPU Optimization

### Algorithm Optimization

```go
// optimization/algorithm.go
package optimization

import (
    "sort"
)

// Inefficient: O(n²) search
func LinearSearch(slice []int, target int) int {
    for i, v := range slice {
        if v == target {
            return i
        }
    }
    return -1
}

// Efficient: O(log n) search (requires sorted slice)
func BinarySearch(slice []int, target int) int {
    return sort.SearchInts(slice, target)
}

// Map-based lookup: O(1) average case
func MapLookup(m map[int]bool, target int) bool {
    return m[target]
}

// Efficient sorting for small slices
func InsertionSort(slice []int) {
    for i := 1; i < len(slice); i++ {
        key := slice[i]
        j := i - 1
        
        for j >= 0 && slice[j] > key {
            slice[j+1] = slice[j]
            j--
        }
        slice[j+1] = key
    }
}

// Hybrid sort: insertion sort for small slices, quicksort for large
func HybridSort(slice []int) {
    if len(slice) < 10 {
        InsertionSort(slice)
    } else {
        sort.Ints(slice)
    }
}

// Efficient duplicate removal
func RemoveDuplicates(slice []int) []int {
    if len(slice) == 0 {
        return slice
    }
    
    seen := make(map[int]bool)
    result := make([]int, 0, len(slice))
    
    for _, v := range slice {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    
    return result
}

// In-place duplicate removal for sorted slices
func RemoveDuplicatesSorted(slice []int) []int {
    if len(slice) == 0 {
        return slice
    }
    
    j := 0
    for i := 1; i < len(slice); i++ {
        if slice[i] != slice[j] {
            j++
            slice[j] = slice[i]
        }
    }
    
    return slice[:j+1]
}

// Efficient set operations
type IntSet map[int]struct{}

func NewIntSet() IntSet {
    return make(IntSet)
}

func (s IntSet) Add(v int) {
    s[v] = struct{}{}
}

func (s IntSet) Contains(v int) bool {
    _, exists := s[v]
    return exists
}

func (s IntSet) Remove(v int) {
    delete(s, v)
}

func (s IntSet) Union(other IntSet) IntSet {
    result := NewIntSet()
    
    for v := range s {
        result.Add(v)
    }
    
    for v := range other {
        result.Add(v)
    }
    
    return result
}

func (s IntSet) Intersection(other IntSet) IntSet {
    result := NewIntSet()
    
    // Iterate over the smaller set
    smaller, larger := s, other
    if len(other) < len(s) {
        smaller, larger = other, s
    }
    
    for v := range smaller {
        if larger.Contains(v) {
            result.Add(v)
        }
    }
    
    return result
}
```

### Concurrency Optimization

```go
// optimization/concurrency.go
package optimization

import (
    "context"
    "runtime"
    "sync"
)

// Worker pool for CPU-bound tasks
type WorkerPool struct {
    workers    int
    jobs       chan func()
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewWorkerPool(workers int) *WorkerPool {
    if workers <= 0 {
        workers = runtime.NumCPU()
    }
    
    ctx, cancel := context.WithCancel(context.Background())
    
    wp := &WorkerPool{
        workers: workers,
        jobs:    make(chan func(), workers*2), // Buffered channel
        ctx:     ctx,
        cancel:  cancel,
    }
    
    wp.start()
    return wp
}

func (wp *WorkerPool) start() {
    for i := 0; i < wp.workers; i++ {
        wp.wg.Add(1)
        go wp.worker()
    }
}

func (wp *WorkerPool) worker() {
    defer wp.wg.Done()
    
    for {
        select {
        case job := <-wp.jobs:
            job()
        case <-wp.ctx.Done():
            return
        }
    }
}

func (wp *WorkerPool) Submit(job func()) {
    select {
    case wp.jobs <- job:
    case <-wp.ctx.Done():
    }
}

func (wp *WorkerPool) Close() {
    wp.cancel()
    wp.wg.Wait()
}

// Parallel processing with optimal worker count
func ProcessParallel(data []int, processor func(int) int) []int {
    numWorkers := runtime.NumCPU()
    chunkSize := len(data) / numWorkers
    if chunkSize == 0 {
        chunkSize = 1
        numWorkers = len(data)
    }
    
    result := make([]int, len(data))
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        start := i * chunkSize
        end := start + chunkSize
        if i == numWorkers-1 {
            end = len(data) // Handle remainder
        }
        
        wg.Add(1)
        go func(start, end int) {
            defer wg.Done()
            for j := start; j < end; j++ {
                result[j] = processor(data[j])
            }
        }(start, end)
    }
    
    wg.Wait()
    return result
}

// Pipeline pattern for stream processing
type Pipeline struct {
    stages []func(<-chan interface{}) <-chan interface{}
}

func NewPipeline() *Pipeline {
    return &Pipeline{}
}

func (p *Pipeline) AddStage(stage func(<-chan interface{}) <-chan interface{}) {
    p.stages = append(p.stages, stage)
}

func (p *Pipeline) Execute(input <-chan interface{}) <-chan interface{} {
    current := input
    for _, stage := range p.stages {
        current = stage(current)
    }
    return current
}

// Example pipeline stages
func FilterStage(predicate func(interface{}) bool) func(<-chan interface{}) <-chan interface{} {
    return func(input <-chan interface{}) <-chan interface{} {
        output := make(chan interface{})
        go func() {
            defer close(output)
            for item := range input {
                if predicate(item) {
                    output <- item
                }
            }
        }()
        return output
    }
}

func MapStage(mapper func(interface{}) interface{}) func(<-chan interface{}) <-chan interface{} {
    return func(input <-chan interface{}) <-chan interface{} {
        output := make(chan interface{})
        go func() {
            defer close(output)
            for item := range input {
                output <- mapper(item)
            }
        }()
        return output
    }
}

// Batch processing for better throughput
func BatchProcessor(batchSize int, processor func([]interface{}) []interface{}) func(<-chan interface{}) <-chan interface{} {
    return func(input <-chan interface{}) <-chan interface{} {
        output := make(chan interface{})
        
        go func() {
            defer close(output)
            
            batch := make([]interface{}, 0, batchSize)
            
            for item := range input {
                batch = append(batch, item)
                
                if len(batch) == batchSize {
                    results := processor(batch)
                    for _, result := range results {
                        output <- result
                    }
                    batch = batch[:0] // Reset batch
                }
            }
            
            // Process remaining items
            if len(batch) > 0 {
                results := processor(batch)
                for _, result := range results {
                    output <- result
                }
            }
        }()
        
        return output
    }
}
```

## Garbage Collection Optimization

### GC Tuning

```go
// gc/tuning.go
package gc

import (
    "runtime"
    "runtime/debug"
    "time"
)

// GC configuration
type GCConfig struct {
    TargetPercent int
    MaxHeap       int64
    GCPercent     int
}

func OptimizeGC(config GCConfig) {
    // Set GC target percentage
    if config.GCPercent > 0 {
        debug.SetGCPercent(config.GCPercent)
    }
    
    // Set memory limit (Go 1.19+)
    if config.MaxHeap > 0 {
        debug.SetMemoryLimit(config.MaxHeap)
    }
    
    // Force initial GC
    runtime.GC()
}

// Monitor GC performance
type GCStats struct {
    NumGC        uint32
    PauseTotal   time.Duration
    LastGC       time.Time
    HeapSize     uint64
    HeapObjects  uint64
}

func GetGCStats() GCStats {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    return GCStats{
        NumGC:       m.NumGC,
        PauseTotal:  time.Duration(m.PauseTotalNs),
        LastGC:      time.Unix(0, int64(m.LastGC)),
        HeapSize:    m.HeapAlloc,
        HeapObjects: m.HeapObjects,
    }
}

// Reduce GC pressure by reusing objects
type ObjectPool struct {
    pool sync.Pool
    new  func() interface{}
}

func NewObjectPool(newFunc func() interface{}) *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{New: newFunc},
        new:  newFunc,
    }
}

func (op *ObjectPool) Get() interface{} {
    return op.pool.Get()
}

func (op *ObjectPool) Put(obj interface{}) {
    op.pool.Put(obj)
}

// Example: Reusable buffer pool
type Buffer struct {
    data []byte
}

func (b *Buffer) Reset() {
    b.data = b.data[:0]
}

func (b *Buffer) Write(p []byte) {
    b.data = append(b.data, p...)
}

func (b *Buffer) Bytes() []byte {
    return b.data
}

var bufferPool = NewObjectPool(func() interface{} {
    return &Buffer{
        data: make([]byte, 0, 1024),
    }
})

func GetBuffer() *Buffer {
    return bufferPool.Get().(*Buffer)
}

func PutBuffer(buf *Buffer) {
    buf.Reset()
    bufferPool.Put(buf)
}

// Minimize allocations in hot paths
func ProcessDataEfficient(data [][]byte) int {
    buf := GetBuffer()
    defer PutBuffer(buf)
    
    totalSize := 0
    for _, chunk := range data {
        buf.Write(chunk)
        totalSize += len(chunk)
    }
    
    return totalSize
}

// Use finalizers sparingly
type Resource struct {
    handle uintptr
    closed bool
}

func NewResource() *Resource {
    r := &Resource{
        handle: allocateResource(), // Hypothetical C function
    }
    
    // Set finalizer as safety net
    runtime.SetFinalizer(r, (*Resource).finalize)
    return r
}

func (r *Resource) Close() {
    if !r.closed {
        freeResource(r.handle) // Hypothetical C function
        r.closed = true
        runtime.SetFinalizer(r, nil) // Remove finalizer
    }
}

func (r *Resource) finalize() {
    if !r.closed {
        // Log warning about resource leak
        freeResource(r.handle)
    }
}

// Hypothetical C functions
func allocateResource() uintptr { return 0 }
func freeResource(handle uintptr) {}
```

### Memory-Efficient Data Structures

```go
// structures/efficient.go
package structures

import (
    "unsafe"
)

// Compact struct layout
type CompactStruct struct {
    // Group fields by size to minimize padding
    flag1    bool   // 1 byte
    flag2    bool   // 1 byte
    smallInt uint16 // 2 bytes
    id       uint32 // 4 bytes
    value    uint64 // 8 bytes
    name     string // 16 bytes (on 64-bit)
}

// Inefficient struct layout (more padding)
type InefficientStruct struct {
    flag1    bool   // 1 byte + 7 bytes padding
    value    uint64 // 8 bytes
    flag2    bool   // 1 byte + 1 byte padding
    smallInt uint16 // 2 bytes + 4 bytes padding
    id       uint32 // 4 bytes + 4 bytes padding
    name     string // 16 bytes
}

// Bit fields for flags
type Flags uint32

const (
    FlagA Flags = 1 << iota
    FlagB
    FlagC
    FlagD
)

func (f Flags) Has(flag Flags) bool {
    return f&flag != 0
}

func (f *Flags) Set(flag Flags) {
    *f |= flag
}

func (f *Flags) Clear(flag Flags) {
    *f &^= flag
}

// Slice of structs vs slice of pointers
type Point struct {
    X, Y float64
}

// More cache-friendly: data is contiguous
func ProcessPointsValue(points []Point) float64 {
    sum := 0.0
    for _, p := range points {
        sum += p.X + p.Y
    }
    return sum
}

// Less cache-friendly: data is scattered
func ProcessPointsPointer(points []*Point) float64 {
    sum := 0.0
    for _, p := range points {
        sum += p.X + p.Y
    }
    return sum
}

// String interning to save memory
type StringInterner struct {
    strings map[string]string
}

func NewStringInterner() *StringInterner {
    return &StringInterner{
        strings: make(map[string]string),
    }
}

func (si *StringInterner) Intern(s string) string {
    if interned, exists := si.strings[s]; exists {
        return interned
    }
    
    // Create a copy to ensure we own the memory
    interned := string([]byte(s))
    si.strings[interned] = interned
    return interned
}

// Memory-efficient map for small keys
type SmallKeyMap struct {
    keys   []uint32
    values []interface{}
}

func NewSmallKeyMap() *SmallKeyMap {
    return &SmallKeyMap{
        keys:   make([]uint32, 0, 8),
        values: make([]interface{}, 0, 8),
    }
}

func (m *SmallKeyMap) Get(key uint32) (interface{}, bool) {
    for i, k := range m.keys {
        if k == key {
            return m.values[i], true
        }
    }
    return nil, false
}

func (m *SmallKeyMap) Set(key uint32, value interface{}) {
    for i, k := range m.keys {
        if k == key {
            m.values[i] = value
            return
        }
    }
    
    m.keys = append(m.keys, key)
    m.values = append(m.values, value)
}

// Zero-allocation string operations
func HasPrefixBytes(s, prefix []byte) bool {
    return len(s) >= len(prefix) && 
           string(s[:len(prefix)]) == string(prefix)
}

// Unsafe string/byte conversions (use with caution)
func StringToBytes(s string) []byte {
    return *(*[]byte)(unsafe.Pointer(
        &struct {
            string
            Cap int
        }{s, len(s)},
    ))
}

func BytesToString(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}
```

## Performance Monitoring

### Runtime Metrics

```go
// monitoring/metrics.go
package monitoring

import (
    "context"
    "runtime"
    "time"
)

type Metrics struct {
    // Memory metrics
    HeapAlloc     uint64
    HeapSys       uint64
    HeapIdle      uint64
    HeapInuse     uint64
    HeapReleased  uint64
    HeapObjects   uint64
    
    // GC metrics
    GCCycles      uint32
    GCPauseTotal  time.Duration
    GCPauseAvg    time.Duration
    
    // Goroutine metrics
    NumGoroutines int
    NumCPU        int
    
    // Custom metrics
    Timestamp time.Time
}

func CollectMetrics() Metrics {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    metrics := Metrics{
        HeapAlloc:     m.HeapAlloc,
        HeapSys:       m.HeapSys,
        HeapIdle:      m.HeapIdle,
        HeapInuse:     m.HeapInuse,
        HeapReleased:  m.HeapReleased,
        HeapObjects:   m.HeapObjects,
        GCCycles:      m.NumGC,
        GCPauseTotal:  time.Duration(m.PauseTotalNs),
        NumGoroutines: runtime.NumGoroutine(),
        NumCPU:        runtime.NumCPU(),
        Timestamp:     time.Now(),
    }
    
    if m.NumGC > 0 {
        metrics.GCPauseAvg = time.Duration(m.PauseTotalNs / uint64(m.NumGC))
    }
    
    return metrics
}

type MetricsCollector struct {
    metrics chan Metrics
    done    chan struct{}
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        metrics: make(chan Metrics, 100),
        done:    make(chan struct{}),
    }
}

func (mc *MetricsCollector) Start(ctx context.Context, interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            metrics := CollectMetrics()
            select {
            case mc.metrics <- metrics:
            default:
                // Channel full, drop metric
            }
        case <-ctx.Done():
            close(mc.done)
            return
        }
    }
}

func (mc *MetricsCollector) GetMetrics() <-chan Metrics {
    return mc.metrics
}

func (mc *MetricsCollector) Stop() {
    <-mc.done
}

// Performance profiler
type Profiler struct {
    startTime time.Time
    samples   []time.Duration
}

func NewProfiler() *Profiler {
    return &Profiler{
        samples: make([]time.Duration, 0, 1000),
    }
}

func (p *Profiler) Start() {
    p.startTime = time.Now()
}

func (p *Profiler) Stop() time.Duration {
    duration := time.Since(p.startTime)
    p.samples = append(p.samples, duration)
    return duration
}

func (p *Profiler) Stats() (min, max, avg time.Duration) {
    if len(p.samples) == 0 {
        return 0, 0, 0
    }
    
    min = p.samples[0]
    max = p.samples[0]
    total := time.Duration(0)
    
    for _, sample := range p.samples {
        if sample < min {
            min = sample
        }
        if sample > max {
            max = sample
        }
        total += sample
    }
    
    avg = total / time.Duration(len(p.samples))
    return min, max, avg
}

func (p *Profiler) Reset() {
    p.samples = p.samples[:0]
}

// Function timing decorator
func TimeFunction(name string, fn func()) time.Duration {
    start := time.Now()
    fn()
    duration := time.Since(start)
    
    // Log or store timing information
    // log.Printf("%s took %v", name, duration)
    
    return duration
}

// Method timing with defer
func (p *Profiler) TimeMethod(name string) func() {
    start := time.Now()
    return func() {
        duration := time.Since(start)
        p.samples = append(p.samples, duration)
        // log.Printf("%s took %v", name, duration)
    }
}

// Usage example:
// defer profiler.TimeMethod("MyMethod")()
```

## Exercises

### Exercise 1: Optimize String Processing
Optimize a function that processes large text files:
- Use efficient string operations
- Minimize memory allocations
- Implement proper buffering
- Add benchmarks to measure improvements

### Exercise 2: Concurrent Data Processing
Implement a concurrent data processing pipeline:
- Process data in parallel using worker pools
- Optimize for CPU and memory usage
- Add proper error handling and cancellation
- Measure throughput and latency

### Exercise 3: Memory Pool Implementation
Create a memory pool for frequently allocated objects:
- Implement object reuse to reduce GC pressure
- Add metrics to track pool efficiency
- Compare performance with and without pooling
- Handle different object sizes

## Key Takeaways

- Profile before optimizing - measure to identify bottlenecks
- Use `go test -bench` and `go tool pprof` for performance analysis
- Pre-allocate slices and maps when size is known
- Reuse objects with sync.Pool to reduce GC pressure
- Choose appropriate data structures for your use case
- Optimize for cache locality with contiguous data
- Use goroutines and channels efficiently
- Minimize allocations in hot code paths
- Consider unsafe operations only when necessary and safe
- Monitor production performance with runtime metrics
- Optimize algorithms before micro-optimizations
- Balance readability with performance

## Next Steps

Next, we'll explore [Microservices](28-microservices.md) to learn about building distributed systems and microservice architectures with Go.