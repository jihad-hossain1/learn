# 31. Advanced Topics

This chapter covers advanced Go programming concepts and techniques that are essential for building sophisticated applications and understanding Go's deeper capabilities.

## Reflection

Reflection in Go allows programs to examine their own structure at runtime. The `reflect` package provides this capability.

### Basic Reflection

```go
// reflection_basics.go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name    string `json:"name" validate:"required"`
    Age     int    `json:"age" validate:"min=0,max=150"`
    Email   string `json:"email" validate:"email"`
    Address *Address
}

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
}

func examineType(i interface{}) {
    // Get the type and value
    t := reflect.TypeOf(i)
    v := reflect.ValueOf(i)
    
    fmt.Printf("Type: %s\n", t.Name())
    fmt.Printf("Kind: %s\n", t.Kind())
    fmt.Printf("Value: %v\n", v)
    fmt.Printf("Is valid: %t\n", v.IsValid())
    
    // Handle pointer types
    if t.Kind() == reflect.Ptr {
        fmt.Printf("Pointer to: %s\n", t.Elem().Name())
        if !v.IsNil() {
            v = v.Elem()
            t = t.Elem()
        }
    }
    
    // Examine struct fields
    if t.Kind() == reflect.Struct {
        fmt.Printf("Number of fields: %d\n", t.NumField())
        
        for i := 0; i < t.NumField(); i++ {
            field := t.Field(i)
            value := v.Field(i)
            
            fmt.Printf("  Field %d: %s (type: %s, value: %v)\n", 
                i, field.Name, field.Type, value)
            
            // Examine struct tags
            if tag := field.Tag; tag != "" {
                fmt.Printf("    JSON tag: %s\n", tag.Get("json"))
                fmt.Printf("    Validate tag: %s\n", tag.Get("validate"))
            }
        }
    }
    
    fmt.Println()
}

func main() {
    // Basic types
    examineType(42)
    examineType("hello")
    examineType([]int{1, 2, 3})
    
    // Struct
    person := Person{
        Name:  "John Doe",
        Age:   30,
        Email: "john@example.com",
        Address: &Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
        },
    }
    
    examineType(person)
    examineType(&person)
}
```

### Dynamic Function Calls

```go
// dynamic_calls.go
package main

import (
    "fmt"
    "reflect"
    "strings"
)

type Calculator struct{}

func (c Calculator) Add(a, b int) int {
    return a + b
}

func (c Calculator) Subtract(a, b int) int {
    return a - b
}

func (c Calculator) Multiply(a, b int) int {
    return a * b
}

func (c Calculator) Divide(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

func (c Calculator) Concat(strs ...string) string {
    return strings.Join(strs, " ")
}

// Call method by name using reflection
func callMethod(obj interface{}, methodName string, args ...interface{}) ([]reflect.Value, error) {
    v := reflect.ValueOf(obj)
    method := v.MethodByName(methodName)
    
    if !method.IsValid() {
        return nil, fmt.Errorf("method %s not found", methodName)
    }
    
    // Check number of arguments
    methodType := method.Type()
    if methodType.NumIn() != len(args) {
        return nil, fmt.Errorf("method %s expects %d arguments, got %d", 
            methodName, methodType.NumIn(), len(args))
    }
    
    // Convert arguments to reflect.Value
    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }
    
    // Call the method
    return method.Call(in), nil
}

// List all methods of a type
func listMethods(obj interface{}) {
    t := reflect.TypeOf(obj)
    v := reflect.ValueOf(obj)
    
    fmt.Printf("Methods of %s:\n", t.Name())
    
    for i := 0; i < t.NumMethod(); i++ {
        method := t.Method(i)
        fmt.Printf("  %s: %s\n", method.Name, method.Type)
        
        // Show method signature details
        methodType := method.Type
        fmt.Printf("    Inputs: %d, Outputs: %d\n", 
            methodType.NumIn(), methodType.NumOut())
        
        for j := 1; j < methodType.NumIn(); j++ { // Skip receiver
            fmt.Printf("      In[%d]: %s\n", j-1, methodType.In(j))
        }
        
        for j := 0; j < methodType.NumOut(); j++ {
            fmt.Printf("      Out[%d]: %s\n", j, methodType.Out(j))
        }
    }
    fmt.Println()
}

func main() {
    calc := Calculator{}
    
    // List all methods
    listMethods(calc)
    
    // Call methods dynamically
    result, err := callMethod(calc, "Add", 10, 5)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Add(10, 5) = %v\n", result[0].Int())
    }
    
    result, err = callMethod(calc, "Divide", 10, 2)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Divide(10, 2) = %v, error = %v\n", 
            result[0].Int(), result[1].Interface())
    }
    
    // Call with error
    result, err = callMethod(calc, "Divide", 10, 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Divide(10, 0) = %v, error = %v\n", 
            result[0].Int(), result[1].Interface())
    }
    
    // Variadic function
    result, err = callMethod(calc, "Concat", "Hello", "World", "from", "Go")
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Concat result = %v\n", result[0].String())
    }
}
```

### Struct Tag Processing

```go
// struct_tags.go
package main

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

type User struct {
    ID       int    `db:"id" json:"id" validate:"required"`
    Username string `db:"username" json:"username" validate:"required,min=3,max=20"`
    Email    string `db:"email" json:"email" validate:"required,email"`
    Age      int    `db:"age" json:"age" validate:"min=0,max=150"`
    IsActive bool   `db:"is_active" json:"is_active"`
}

// Validation rules
type ValidationRule struct {
    Field string
    Rules []string
}

// Extract validation rules from struct tags
func extractValidationRules(obj interface{}) []ValidationRule {
    var rules []ValidationRule
    
    t := reflect.TypeOf(obj)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        validateTag := field.Tag.Get("validate")
        
        if validateTag != "" {
            fieldRules := strings.Split(validateTag, ",")
            rules = append(rules, ValidationRule{
                Field: field.Name,
                Rules: fieldRules,
            })
        }
    }
    
    return rules
}

// Simple validator
func validateStruct(obj interface{}) []string {
    var errors []string
    
    v := reflect.ValueOf(obj)
    t := reflect.TypeOf(obj)
    
    if t.Kind() == reflect.Ptr {
        v = v.Elem()
        t = t.Elem()
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        validateTag := field.Tag.Get("validate")
        
        if validateTag == "" {
            continue
        }
        
        rules := strings.Split(validateTag, ",")
        
        for _, rule := range rules {
            rule = strings.TrimSpace(rule)
            
            switch {
            case rule == "required":
                if isZeroValue(value) {
                    errors = append(errors, fmt.Sprintf("%s is required", field.Name))
                }
                
            case strings.HasPrefix(rule, "min="):
                minStr := strings.TrimPrefix(rule, "min=")
                min, err := strconv.Atoi(minStr)
                if err != nil {
                    continue
                }
                
                switch value.Kind() {
                case reflect.String:
                    if len(value.String()) < min {
                        errors = append(errors, fmt.Sprintf("%s must be at least %d characters", field.Name, min))
                    }
                case reflect.Int:
                    if int(value.Int()) < min {
                        errors = append(errors, fmt.Sprintf("%s must be at least %d", field.Name, min))
                    }
                }
                
            case strings.HasPrefix(rule, "max="):
                maxStr := strings.TrimPrefix(rule, "max=")
                max, err := strconv.Atoi(maxStr)
                if err != nil {
                    continue
                }
                
                switch value.Kind() {
                case reflect.String:
                    if len(value.String()) > max {
                        errors = append(errors, fmt.Sprintf("%s must be at most %d characters", field.Name, max))
                    }
                case reflect.Int:
                    if int(value.Int()) > max {
                        errors = append(errors, fmt.Sprintf("%s must be at most %d", field.Name, max))
                    }
                }
                
            case rule == "email":
                if value.Kind() == reflect.String {
                    email := value.String()
                    if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
                        errors = append(errors, fmt.Sprintf("%s must be a valid email", field.Name))
                    }
                }
            }
        }
    }
    
    return errors
}

func isZeroValue(v reflect.Value) bool {
    switch v.Kind() {
    case reflect.String:
        return v.String() == ""
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return v.Int() == 0
    case reflect.Bool:
        return !v.Bool()
    case reflect.Ptr, reflect.Interface:
        return v.IsNil()
    default:
        return false
    }
}

// Generate SQL from struct tags
func generateSQL(obj interface{}, tableName string) string {
    t := reflect.TypeOf(obj)
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
    }
    
    var columns []string
    var placeholders []string
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        dbTag := field.Tag.Get("db")
        
        if dbTag != "" && dbTag != "-" {
            columns = append(columns, dbTag)
            placeholders = append(placeholders, "?")
        }
    }
    
    return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
        tableName,
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "))
}

func main() {
    user := User{
        ID:       1,
        Username: "john_doe",
        Email:    "john@example.com",
        Age:      25,
        IsActive: true,
    }
    
    // Extract validation rules
    rules := extractValidationRules(user)
    fmt.Println("Validation Rules:")
    for _, rule := range rules {
        fmt.Printf("  %s: %v\n", rule.Field, rule.Rules)
    }
    fmt.Println()
    
    // Validate valid struct
    errors := validateStruct(user)
    if len(errors) == 0 {
        fmt.Println("User is valid")
    } else {
        fmt.Println("Validation errors:")
        for _, err := range errors {
            fmt.Printf("  - %s\n", err)
        }
    }
    
    // Validate invalid struct
    invalidUser := User{
        Username: "jo", // Too short
        Email:    "invalid-email", // Invalid format
        Age:      200, // Too high
    }
    
    errors = validateStruct(invalidUser)
    fmt.Println("\nValidation errors for invalid user:")
    for _, err := range errors {
        fmt.Printf("  - %s\n", err)
    }
    
    // Generate SQL
    sql := generateSQL(user, "users")
    fmt.Printf("\nGenerated SQL: %s\n", sql)
}
```

## Unsafe Package

The `unsafe` package provides low-level memory operations. Use with extreme caution.

### Basic Unsafe Operations

```go
// unsafe_operations.go
package main

import (
    "fmt"
    "unsafe"
)

type Person struct {
    Name string
    Age  int32
    City string
}

func demonstrateUnsafe() {
    // Basic pointer arithmetic
    numbers := []int{1, 2, 3, 4, 5}
    
    // Get pointer to first element
    ptr := unsafe.Pointer(&numbers[0])
    
    // Access elements using pointer arithmetic
    for i := 0; i < len(numbers); i++ {
        // Calculate offset
        offset := uintptr(i) * unsafe.Sizeof(numbers[0])
        elementPtr := unsafe.Pointer(uintptr(ptr) + offset)
        
        // Dereference pointer
        value := *(*int)(elementPtr)
        fmt.Printf("Element %d: %d\n", i, value)
    }
    
    fmt.Println()
}

func structFieldAccess() {
    person := Person{
        Name: "John Doe",
        Age:  30,
        City: "New York",
    }
    
    // Get pointer to struct
    ptr := unsafe.Pointer(&person)
    
    // Access Name field (offset 0)
    namePtr := (*string)(ptr)
    fmt.Printf("Name: %s\n", *namePtr)
    
    // Access Age field (after Name field)
    nameSize := unsafe.Sizeof(person.Name)
    agePtr := (*int32)(unsafe.Pointer(uintptr(ptr) + nameSize))
    fmt.Printf("Age: %d\n", *agePtr)
    
    // Access City field (after Name and Age fields)
    ageSize := unsafe.Sizeof(person.Age)
    cityPtr := (*string)(unsafe.Pointer(uintptr(ptr) + nameSize + ageSize))
    fmt.Printf("City: %s\n", *cityPtr)
    
    fmt.Println()
}

// Convert string to []byte without allocation
func stringToBytes(s string) []byte {
    return *(*[]byte)(unsafe.Pointer(&s))
}

// Convert []byte to string without allocation
func bytesToString(b []byte) string {
    return *(*string)(unsafe.Pointer(&b))
}

// Zero-copy string operations
func zeroCopyOperations() {
    original := "Hello, World!"
    
    // Convert to bytes without copying
    bytes := stringToBytes(original)
    fmt.Printf("Original: %s\n", original)
    fmt.Printf("As bytes: %v\n", bytes)
    
    // Modify bytes (dangerous!)
    // bytes[0] = 'h' // This would cause a runtime panic
    
    // Convert back to string
    converted := bytesToString(bytes)
    fmt.Printf("Converted back: %s\n", converted)
    
    // Create a new byte slice for safe modification
    safeCopy := make([]byte, len(bytes))
    copy(safeCopy, bytes)
    safeCopy[0] = 'h'
    
    modified := bytesToString(safeCopy)
    fmt.Printf("Modified: %s\n", modified)
    
    fmt.Println()
}

// Memory layout inspection
func inspectMemoryLayout() {
    type TestStruct struct {
        A int8   // 1 byte
        B int64  // 8 bytes
        C int16  // 2 bytes
        D int32  // 4 bytes
    }
    
    var ts TestStruct
    
    fmt.Printf("TestStruct size: %d bytes\n", unsafe.Sizeof(ts))
    fmt.Printf("TestStruct alignment: %d bytes\n", unsafe.Alignof(ts))
    
    fmt.Printf("Field A offset: %d, size: %d, align: %d\n", 
        unsafe.Offsetof(ts.A), unsafe.Sizeof(ts.A), unsafe.Alignof(ts.A))
    fmt.Printf("Field B offset: %d, size: %d, align: %d\n", 
        unsafe.Offsetof(ts.B), unsafe.Sizeof(ts.B), unsafe.Alignof(ts.B))
    fmt.Printf("Field C offset: %d, size: %d, align: %d\n", 
        unsafe.Offsetof(ts.C), unsafe.Sizeof(ts.C), unsafe.Alignof(ts.C))
    fmt.Printf("Field D offset: %d, size: %d, align: %d\n", 
        unsafe.Offsetof(ts.D), unsafe.Sizeof(ts.D), unsafe.Alignof(ts.D))
    
    fmt.Println()
}

func main() {
    fmt.Println("=== Unsafe Package Demonstrations ===")
    fmt.Println("WARNING: These operations are unsafe and should be used with extreme caution!\n")
    
    demonstrateUnsafe()
    structFieldAccess()
    zeroCopyOperations()
    inspectMemoryLayout()
}
```

## Code Generation

Go supports code generation through the `go generate` command and build tags.

### Using go:generate

```go
// generator.go
//go:generate go run generator.go

package main

import (
    "fmt"
    "os"
    "text/template"
)

type EnumValue struct {
    Name  string
    Value int
}

type EnumData struct {
    PackageName string
    EnumName    string
    Values      []EnumValue
}

const enumTemplate = `// Code generated by go generate; DO NOT EDIT.

package {{.PackageName}}

import "fmt"

// {{.EnumName}} represents the enum type
type {{.EnumName}} int

const (
{{range $i, $v := .Values}}    {{$v.Name}} {{if eq $i 0}}{{$.EnumName}} = iota{{end}}
{{end}})

// String returns the string representation of {{.EnumName}}
func (e {{.EnumName}}) String() string {
    switch e {
{{range .Values}}    case {{.Name}}:
        return "{{.Name}}"
{{end}}    default:
        return fmt.Sprintf("{{.EnumName}}(%d)", int(e))
    }
}

// IsValid checks if the {{.EnumName}} value is valid
func (e {{.EnumName}}) IsValid() bool {
    switch e {
{{range .Values}}    case {{.Name}}:
        return true
{{end}}    default:
        return false
    }
}

// All{{.EnumName}}s returns all valid {{.EnumName}} values
func All{{.EnumName}}s() []{{.EnumName}} {
    return []{{.EnumName}}{
{{range .Values}}        {{.Name}},
{{end}}    }
}
`

func main() {
    // Define enum data
    data := EnumData{
        PackageName: "main",
        EnumName:    "Status",
        Values: []EnumValue{
            {"StatusPending", 0},
            {"StatusInProgress", 1},
            {"StatusCompleted", 2},
            {"StatusFailed", 3},
            {"StatusCancelled", 4},
        },
    }
    
    // Parse template
    tmpl, err := template.New("enum").Parse(enumTemplate)
    if err != nil {
        fmt.Printf("Error parsing template: %v\n", err)
        os.Exit(1)
    }
    
    // Create output file
    file, err := os.Create("status_enum.go")
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        os.Exit(1)
    }
    defer file.Close()
    
    // Execute template
    err = tmpl.Execute(file, data)
    if err != nil {
        fmt.Printf("Error executing template: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("Generated status_enum.go")
}
```

### Build Tags

```go
// debug_build.go
//go:build debug
// +build debug

package main

import "fmt"

func init() {
    fmt.Println("Debug build enabled")
}

func debugLog(msg string) {
    fmt.Printf("[DEBUG] %s\n", msg)
}
```

```go
// release_build.go
//go:build !debug
// +build !debug

package main

func debugLog(msg string) {
    // No-op in release builds
}
```

```go
// platform_specific.go
//go:build windows
// +build windows

package main

import "fmt"

func platformSpecificFunction() {
    fmt.Println("Running on Windows")
}
```

```go
// platform_unix.go
//go:build unix
// +build unix

package main

import "fmt"

func platformSpecificFunction() {
    fmt.Println("Running on Unix-like system")
}
```

## Assembly Integration

Go allows integration with assembly code for performance-critical operations.

### Assembly Function Declaration

```go
// math_asm.go
package main

// Add two integers using assembly
func addAsm(a, b int64) int64

// Multiply two integers using assembly
func mulAsm(a, b int64) int64
```

```assembly
// math_amd64.s
#include "textflag.h"

// func addAsm(a, b int64) int64
TEXT ·addAsm(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

// func mulAsm(a, b int64) int64
TEXT ·mulAsm(SB), NOSPLIT, $0-24
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    IMULQ BX, AX
    MOVQ AX, ret+16(FP)
    RET
```

## Memory Management

### Custom Memory Allocators

```go
// memory_pool.go
package main

import (
    "fmt"
    "sync"
    "unsafe"
)

// MemoryPool represents a simple memory pool
type MemoryPool struct {
    mu       sync.Mutex
    blocks   [][]byte
    blockSize int
    available [][]byte
}

// NewMemoryPool creates a new memory pool
func NewMemoryPool(blockSize, initialBlocks int) *MemoryPool {
    pool := &MemoryPool{
        blockSize: blockSize,
        blocks:    make([][]byte, 0, initialBlocks),
        available: make([][]byte, 0, initialBlocks),
    }
    
    // Pre-allocate blocks
    for i := 0; i < initialBlocks; i++ {
        block := make([]byte, blockSize)
        pool.blocks = append(pool.blocks, block)
        pool.available = append(pool.available, block)
    }
    
    return pool
}

// Get retrieves a block from the pool
func (p *MemoryPool) Get() []byte {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    if len(p.available) == 0 {
        // Allocate new block if pool is empty
        block := make([]byte, p.blockSize)
        p.blocks = append(p.blocks, block)
        return block
    }
    
    // Get block from available list
    block := p.available[len(p.available)-1]
    p.available = p.available[:len(p.available)-1]
    
    // Clear the block
    for i := range block {
        block[i] = 0
    }
    
    return block
}

// Put returns a block to the pool
func (p *MemoryPool) Put(block []byte) {
    if len(block) != p.blockSize {
        return // Wrong size, don't return to pool
    }
    
    p.mu.Lock()
    defer p.mu.Unlock()
    
    p.available = append(p.available, block)
}

// Stats returns pool statistics
func (p *MemoryPool) Stats() (total, available int) {
    p.mu.Lock()
    defer p.mu.Unlock()
    
    return len(p.blocks), len(p.available)
}

// Object pool for reusing objects
type ObjectPool struct {
    pool sync.Pool
    new  func() interface{}
}

// NewObjectPool creates a new object pool
func NewObjectPool(newFunc func() interface{}) *ObjectPool {
    return &ObjectPool{
        pool: sync.Pool{
            New: newFunc,
        },
        new: newFunc,
    }
}

// Get retrieves an object from the pool
func (p *ObjectPool) Get() interface{} {
    return p.pool.Get()
}

// Put returns an object to the pool
func (p *ObjectPool) Put(obj interface{}) {
    p.pool.Put(obj)
}

// Example usage
type Buffer struct {
    data []byte
}

func (b *Buffer) Reset() {
    b.data = b.data[:0]
}

func (b *Buffer) Write(data []byte) {
    b.data = append(b.data, data...)
}

func (b *Buffer) Bytes() []byte {
    return b.data
}

func demonstrateMemoryManagement() {
    // Memory pool example
    fmt.Println("=== Memory Pool Example ===")
    pool := NewMemoryPool(1024, 5)
    
    // Get some blocks
    block1 := pool.Get()
    block2 := pool.Get()
    
    total, available := pool.Stats()
    fmt.Printf("After getting 2 blocks: total=%d, available=%d\n", total, available)
    
    // Use blocks
    copy(block1, []byte("Hello, World!"))
    copy(block2, []byte("Memory Pool"))
    
    fmt.Printf("Block 1: %s\n", string(block1[:13]))
    fmt.Printf("Block 2: %s\n", string(block2[:11]))
    
    // Return blocks
    pool.Put(block1)
    pool.Put(block2)
    
    total, available = pool.Stats()
    fmt.Printf("After returning blocks: total=%d, available=%d\n", total, available)
    
    // Object pool example
    fmt.Println("\n=== Object Pool Example ===")
    bufferPool := NewObjectPool(func() interface{} {
        return &Buffer{
            data: make([]byte, 0, 1024),
        }
    })
    
    // Get buffer from pool
    buf := bufferPool.Get().(*Buffer)
    buf.Write([]byte("Hello from object pool"))
    
    fmt.Printf("Buffer content: %s\n", string(buf.Bytes()))
    
    // Reset and return to pool
    buf.Reset()
    bufferPool.Put(buf)
    
    fmt.Println("Buffer returned to pool")
}

func main() {
    demonstrateMemoryManagement()
}
```

## Performance Optimization Techniques

### CPU Profiling Integration

```go
// profiling.go
package main

import (
    "fmt"
    "os"
    "runtime"
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

// Memory-intensive function
func allocateMemory() {
    var slices [][]int
    
    for i := 0; i < 1000; i++ {
        slice := make([]int, 1000)
        for j := range slice {
            slice[j] = i * j
        }
        slices = append(slices, slice)
    }
    
    // Force garbage collection
    runtime.GC()
}

// Profiling wrapper
func profileFunction(name string, fn func()) {
    // CPU profiling
    cpuFile, err := os.Create(fmt.Sprintf("%s_cpu.prof", name))
    if err != nil {
        fmt.Printf("Error creating CPU profile: %v\n", err)
        return
    }
    defer cpuFile.Close()
    
    if err := pprof.StartCPUProfile(cpuFile); err != nil {
        fmt.Printf("Error starting CPU profile: %v\n", err)
        return
    }
    defer pprof.StopCPUProfile()
    
    // Memory profiling
    memFile, err := os.Create(fmt.Sprintf("%s_mem.prof", name))
    if err != nil {
        fmt.Printf("Error creating memory profile: %v\n", err)
        return
    }
    defer memFile.Close()
    
    // Run function
    start := time.Now()
    fn()
    duration := time.Since(start)
    
    // Write memory profile
    runtime.GC() // Force GC before memory profile
    if err := pprof.WriteHeapProfile(memFile); err != nil {
        fmt.Printf("Error writing memory profile: %v\n", err)
        return
    }
    
    fmt.Printf("%s completed in %v\n", name, duration)
    fmt.Printf("Profiles written: %s_cpu.prof, %s_mem.prof\n", name, name)
}

func main() {
    fmt.Println("=== Performance Profiling Example ===")
    
    // Profile fibonacci calculation
    profileFunction("fibonacci", func() {
        result := fibonacci(35)
        fmt.Printf("Fibonacci(35) = %d\n", result)
    })
    
    // Profile memory allocation
    profileFunction("memory", func() {
        allocateMemory()
        fmt.Println("Memory allocation completed")
    })
    
    fmt.Println("\nTo analyze profiles, use:")
    fmt.Println("go tool pprof fibonacci_cpu.prof")
    fmt.Println("go tool pprof fibonacci_mem.prof")
    fmt.Println("go tool pprof memory_cpu.prof")
    fmt.Println("go tool pprof memory_mem.prof")
}
```

## Best Practices

### 1. Reflection Best Practices

- Use reflection sparingly - it's slower than direct code
- Cache reflection results when possible
- Always check for nil values and validity
- Prefer type assertions over reflection when possible

### 2. Unsafe Package Guidelines

- Only use when absolutely necessary
- Thoroughly test unsafe code
- Document unsafe operations clearly
- Consider alternatives before using unsafe
- Be aware of garbage collector implications

### 3. Code Generation

- Use `go:generate` for repetitive code
- Keep generated code separate from hand-written code
- Include generation commands in documentation
- Version control generation scripts

### 4. Memory Management

- Use object pools for frequently allocated objects
- Monitor memory usage with profiling
- Understand garbage collector behavior
- Minimize allocations in hot paths

## Exercises

### Exercise 1: Generic Validator
Build a generic validation framework using reflection that can:
- Validate struct fields based on tags
- Support custom validation functions
- Provide detailed error messages
- Handle nested structs

### Exercise 2: ORM Framework
Create a simple ORM using reflection that can:
- Map structs to database tables
- Generate SQL queries from struct tags
- Handle relationships between entities
- Support basic CRUD operations

### Exercise 3: Performance Optimizer
Build a performance analysis tool that:
- Profiles CPU and memory usage
- Identifies bottlenecks in code
- Suggests optimization strategies
- Compares before/after performance

### Exercise 4: Code Generator
Create a code generator that:
- Generates REST API handlers from struct definitions
- Creates mock implementations for interfaces
- Generates test cases from function signatures
- Supports custom templates

## Key Takeaways

- Reflection provides runtime introspection but comes with performance costs
- The unsafe package allows low-level memory operations but should be used carefully
- Code generation can eliminate repetitive coding tasks
- Build tags enable conditional compilation
- Assembly integration allows performance optimization for critical code
- Memory management techniques can significantly improve performance
- Profiling is essential for identifying performance bottlenecks
- Advanced features should be used judiciously and with proper testing

## Next Steps

Next, we'll explore [Best Practices](32-best-practices.md) covering comprehensive guidelines for writing maintainable, efficient, and idiomatic Go code.