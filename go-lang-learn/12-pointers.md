# 12. Pointers

Pointers are one of the most important concepts in Go. A pointer holds the memory address of a value. Go pointers are similar to pointers in C, but they are safer and more restricted. Understanding pointers is crucial for efficient memory management and for working with references to data.

## What are Pointers?

A pointer is a variable that stores the memory address of another variable. Instead of holding a value directly, a pointer "points to" the location where the value is stored in memory.

### Basic Pointer Syntax

```go
package main

import "fmt"

func main() {
    var x int = 42
    var p *int = &x  // p is a pointer to int, pointing to x
    
    fmt.Println("Value of x:", x)        // 42
    fmt.Println("Address of x:", &x)     // memory address
    fmt.Println("Value of p:", p)        // same memory address
    fmt.Println("Value at p:", *p)       // 42 (dereferencing)
}
```

### Pointer Operators

- `&` (address-of operator): Gets the memory address of a variable
- `*` (dereference operator): Gets the value stored at a memory address

```go
package main

import "fmt"

func main() {
    var num int = 100
    
    // Create a pointer
    var ptr *int = &num
    
    fmt.Printf("num = %d\n", num)           // 100
    fmt.Printf("&num = %p\n", &num)         // memory address
    fmt.Printf("ptr = %p\n", ptr)           // same memory address
    fmt.Printf("*ptr = %d\n", *ptr)         // 100
    
    // Modify value through pointer
    *ptr = 200
    fmt.Printf("After *ptr = 200:\n")
    fmt.Printf("num = %d\n", num)           // 200
    fmt.Printf("*ptr = %d\n", *ptr)         // 200
}
```

## Declaring and Initializing Pointers

### Zero Value of Pointers

The zero value of a pointer is `nil`.

```go
package main

import "fmt"

func main() {
    var p *int
    fmt.Println("Zero value of pointer:", p) // <nil>
    
    // Check if pointer is nil
    if p == nil {
        fmt.Println("Pointer is nil")
    }
    
    // Initialize the pointer
    var x int = 42
    p = &x
    fmt.Println("After initialization:", p)  // memory address
    fmt.Println("Value at pointer:", *p)     // 42
}
```

### Different Ways to Create Pointers

```go
package main

import "fmt"

func main() {
    // Method 1: Declare variable first, then get its address
    var x int = 10
    var p1 *int = &x
    
    // Method 2: Use new() function
    p2 := new(int)  // Creates a pointer to a zero value int
    *p2 = 20
    
    // Method 3: Short declaration with address operator
    y := 30
    p3 := &y
    
    fmt.Printf("p1 points to: %d\n", *p1)  // 10
    fmt.Printf("p2 points to: %d\n", *p2)  // 20
    fmt.Printf("p3 points to: %d\n", *p3)  // 30
}
```

## Pointers and Functions

### Pass by Value vs Pass by Reference

By default, Go passes arguments by value. To modify the original variable, you need to pass a pointer.

```go
package main

import "fmt"

// Pass by value - original variable is not modified
func incrementByValue(x int) {
    x++
    fmt.Printf("Inside incrementByValue: %d\n", x)
}

// Pass by reference - original variable is modified
func incrementByReference(x *int) {
    *x++
    fmt.Printf("Inside incrementByReference: %d\n", *x)
}

func main() {
    num := 10
    
    fmt.Printf("Original value: %d\n", num)
    
    incrementByValue(num)
    fmt.Printf("After incrementByValue: %d\n", num)  // Still 10
    
    incrementByReference(&num)
    fmt.Printf("After incrementByReference: %d\n", num)  // Now 11
}
```

### Returning Pointers from Functions

```go
package main

import "fmt"

// Function that returns a pointer
func createInt(value int) *int {
    x := value  // Local variable
    return &x   // Return address of local variable (safe in Go)
}

// Function that creates and returns a pointer using new
func createIntWithNew(value int) *int {
    p := new(int)
    *p = value
    return p
}

func main() {
    p1 := createInt(42)
    p2 := createIntWithNew(84)
    
    fmt.Printf("p1 points to: %d\n", *p1)  // 42
    fmt.Printf("p2 points to: %d\n", *p2)  // 84
}
```

### Swapping Values Using Pointers

```go
package main

import "fmt"

func swap(a, b *int) {
    *a, *b = *b, *a
}

func main() {
    x, y := 10, 20
    
    fmt.Printf("Before swap: x=%d, y=%d\n", x, y)
    swap(&x, &y)
    fmt.Printf("After swap: x=%d, y=%d\n", x, y)
}
```

## Pointers to Structs

### Basic Struct Pointers

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    // Create a struct
    p1 := Person{Name: "Alice", Age: 30}
    
    // Create a pointer to the struct
    ptr := &p1
    
    // Access fields through pointer (automatic dereferencing)
    fmt.Printf("Name: %s\n", ptr.Name)  // Same as (*ptr).Name
    fmt.Printf("Age: %d\n", ptr.Age)    // Same as (*ptr).Age
    
    // Modify through pointer
    ptr.Age = 31
    fmt.Printf("Modified age: %d\n", p1.Age)  // 31
}
```

### Creating Struct Pointers with new

```go
package main

import "fmt"

type Rectangle struct {
    Width  float64
    Height float64
}

func (r *Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func main() {
    // Create pointer to struct using new
    rect := new(Rectangle)
    rect.Width = 10
    rect.Height = 5
    
    fmt.Printf("Area: %.2f\n", rect.Area())
    
    rect.Scale(2)
    fmt.Printf("After scaling - Width: %.2f, Height: %.2f\n", rect.Width, rect.Height)
    fmt.Printf("New area: %.2f\n", rect.Area())
}
```

### Struct Literal with Address Operator

```go
package main

import "fmt"

type Point struct {
    X, Y int
}

func main() {
    // Create pointer to struct using struct literal
    p := &Point{X: 10, Y: 20}
    
    fmt.Printf("Point: (%d, %d)\n", p.X, p.Y)
    
    // Modify through pointer
    p.X = 30
    p.Y = 40
    
    fmt.Printf("Modified point: (%d, %d)\n", p.X, p.Y)
}
```

## Pointers to Arrays and Slices

### Pointers to Arrays

```go
package main

import "fmt"

func modifyArray(arr *[5]int) {
    for i := range arr {
        arr[i] *= 2
    }
}

func main() {
    numbers := [5]int{1, 2, 3, 4, 5}
    
    fmt.Printf("Original array: %v\n", numbers)
    
    modifyArray(&numbers)
    
    fmt.Printf("Modified array: %v\n", numbers)
}
```

### Slices are Reference Types

Slices are already reference types, so you don't usually need pointers to slices.

```go
package main

import "fmt"

func modifySlice(slice []int) {
    for i := range slice {
        slice[i] *= 2
    }
}

func appendToSlice(slice *[]int, value int) {
    *slice = append(*slice, value)
}

func main() {
    numbers := []int{1, 2, 3, 4, 5}
    
    fmt.Printf("Original slice: %v\n", numbers)
    
    // Modify existing elements (no pointer needed)
    modifySlice(numbers)
    fmt.Printf("After modification: %v\n", numbers)
    
    // Append new elements (pointer needed to modify slice header)
    appendToSlice(&numbers, 12)
    fmt.Printf("After append: %v\n", numbers)
}
```

## Pointer Arithmetic (Limited)

Unlike C, Go has very limited pointer arithmetic for safety reasons.

```go
package main

import (
    "fmt"
    "unsafe"
)

func main() {
    arr := [3]int{10, 20, 30}
    
    // Get pointer to first element
    ptr := &arr[0]
    
    fmt.Printf("arr[0] = %d, address = %p\n", *ptr, ptr)
    
    // Move to next element using unsafe package (not recommended)
    ptr = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) + unsafe.Sizeof(int(0))))
    
    fmt.Printf("arr[1] = %d, address = %p\n", *ptr, ptr)
    
    // Better approach: use array indexing
    for i := 0; i < len(arr); i++ {
        fmt.Printf("arr[%d] = %d, address = %p\n", i, arr[i], &arr[i])
    }
}
```

## Practical Examples

### Linked List Implementation

```go
package main

import "fmt"

type Node struct {
    Data int
    Next *Node
}

type LinkedList struct {
    Head *Node
}

func (ll *LinkedList) Insert(data int) {
    newNode := &Node{Data: data}
    
    if ll.Head == nil {
        ll.Head = newNode
        return
    }
    
    current := ll.Head
    for current.Next != nil {
        current = current.Next
    }
    current.Next = newNode
}

func (ll *LinkedList) Display() {
    current := ll.Head
    for current != nil {
        fmt.Printf("%d -> ", current.Data)
        current = current.Next
    }
    fmt.Println("nil")
}

func (ll *LinkedList) Delete(data int) {
    if ll.Head == nil {
        return
    }
    
    if ll.Head.Data == data {
        ll.Head = ll.Head.Next
        return
    }
    
    current := ll.Head
    for current.Next != nil && current.Next.Data != data {
        current = current.Next
    }
    
    if current.Next != nil {
        current.Next = current.Next.Next
    }
}

func main() {
    ll := &LinkedList{}
    
    // Insert elements
    ll.Insert(10)
    ll.Insert(20)
    ll.Insert(30)
    ll.Insert(40)
    
    fmt.Println("Linked List:")
    ll.Display()
    
    // Delete an element
    ll.Delete(20)
    fmt.Println("After deleting 20:")
    ll.Display()
}
```

### Binary Tree Implementation

```go
package main

import "fmt"

type TreeNode struct {
    Value int
    Left  *TreeNode
    Right *TreeNode
}

type BinaryTree struct {
    Root *TreeNode
}

func (bt *BinaryTree) Insert(value int) {
    bt.Root = insertNode(bt.Root, value)
}

func insertNode(node *TreeNode, value int) *TreeNode {
    if node == nil {
        return &TreeNode{Value: value}
    }
    
    if value < node.Value {
        node.Left = insertNode(node.Left, value)
    } else if value > node.Value {
        node.Right = insertNode(node.Right, value)
    }
    
    return node
}

func (bt *BinaryTree) InorderTraversal() {
    inorder(bt.Root)
    fmt.Println()
}

func inorder(node *TreeNode) {
    if node != nil {
        inorder(node.Left)
        fmt.Printf("%d ", node.Value)
        inorder(node.Right)
    }
}

func (bt *BinaryTree) Search(value int) bool {
    return search(bt.Root, value)
}

func search(node *TreeNode, value int) bool {
    if node == nil {
        return false
    }
    
    if value == node.Value {
        return true
    } else if value < node.Value {
        return search(node.Left, value)
    } else {
        return search(node.Right, value)
    }
}

func main() {
    bt := &BinaryTree{}
    
    // Insert values
    values := []int{50, 30, 70, 20, 40, 60, 80}
    for _, v := range values {
        bt.Insert(v)
    }
    
    fmt.Println("Inorder traversal:")
    bt.InorderTraversal()
    
    // Search for values
    searchValues := []int{40, 90, 20}
    for _, v := range searchValues {
        if bt.Search(v) {
            fmt.Printf("%d found in tree\n", v)
        } else {
            fmt.Printf("%d not found in tree\n", v)
        }
    }
}
```

### Memory Pool Example

```go
package main

import "fmt"

type Object struct {
    ID   int
    Data string
    next *Object // For free list
}

type ObjectPool struct {
    free *Object
    size int
}

func NewObjectPool(size int) *ObjectPool {
    pool := &ObjectPool{size: size}
    
    // Pre-allocate objects and link them
    for i := 0; i < size; i++ {
        obj := &Object{next: pool.free}
        pool.free = obj
    }
    
    return pool
}

func (p *ObjectPool) Get() *Object {
    if p.free == nil {
        return nil // Pool exhausted
    }
    
    obj := p.free
    p.free = obj.next
    obj.next = nil
    
    return obj
}

func (p *ObjectPool) Put(obj *Object) {
    // Reset object
    obj.ID = 0
    obj.Data = ""
    
    // Add to free list
    obj.next = p.free
    p.free = obj
}

func main() {
    pool := NewObjectPool(3)
    
    // Get objects from pool
    obj1 := pool.Get()
    obj2 := pool.Get()
    obj3 := pool.Get()
    obj4 := pool.Get() // Should be nil
    
    if obj1 != nil {
        obj1.ID = 1
        obj1.Data = "Object 1"
        fmt.Printf("Got object: %+v\n", obj1)
    }
    
    if obj2 != nil {
        obj2.ID = 2
        obj2.Data = "Object 2"
        fmt.Printf("Got object: %+v\n", obj2)
    }
    
    if obj3 != nil {
        obj3.ID = 3
        obj3.Data = "Object 3"
        fmt.Printf("Got object: %+v\n", obj3)
    }
    
    if obj4 == nil {
        fmt.Println("Pool exhausted, obj4 is nil")
    }
    
    // Return objects to pool
    pool.Put(obj1)
    pool.Put(obj2)
    
    // Get objects again
    obj5 := pool.Get()
    obj6 := pool.Get()
    
    if obj5 != nil {
        obj5.ID = 5
        obj5.Data = "Object 5"
        fmt.Printf("Reused object: %+v\n", obj5)
    }
    
    if obj6 != nil {
        obj6.ID = 6
        obj6.Data = "Object 6"
        fmt.Printf("Reused object: %+v\n", obj6)
    }
}
```

## Common Pointer Pitfalls

### Nil Pointer Dereference

```go
package main

import "fmt"

func safeDereference(p *int) {
    if p != nil {
        fmt.Printf("Value: %d\n", *p)
    } else {
        fmt.Println("Pointer is nil")
    }
}

func main() {
    var p *int
    
    // This would cause a panic
    // fmt.Println(*p) // panic: runtime error: invalid memory address
    
    // Safe way
    safeDereference(p)
    
    // Initialize pointer
    x := 42
    p = &x
    safeDereference(p)
}
```

### Pointer to Loop Variable

```go
package main

import "fmt"

func main() {
    // Wrong way - all pointers point to the same variable
    var ptrs []*int
    for i := 0; i < 3; i++ {
        ptrs = append(ptrs, &i) // All point to the same 'i'
    }
    
    fmt.Println("Wrong way:")
    for _, p := range ptrs {
        fmt.Printf("%d ", *p) // Prints: 3 3 3
    }
    fmt.Println()
    
    // Correct way - create a copy of the loop variable
    var correctPtrs []*int
    for i := 0; i < 3; i++ {
        temp := i // Create a copy
        correctPtrs = append(correctPtrs, &temp)
    }
    
    fmt.Println("Correct way:")
    for _, p := range correctPtrs {
        fmt.Printf("%d ", *p) // Prints: 0 1 2
    }
    fmt.Println()
}
```

### Memory Leaks with Pointers

```go
package main

import (
    "fmt"
    "runtime"
)

type LargeStruct struct {
    Data [1000000]int
    Next *LargeStruct
}

func createChain(length int) *LargeStruct {
    var head *LargeStruct
    var current *LargeStruct
    
    for i := 0; i < length; i++ {
        node := &LargeStruct{}
        if head == nil {
            head = node
            current = node
        } else {
            current.Next = node
            current = node
        }
    }
    
    return head
}

func printMemStats(label string) {
    var m runtime.MemStats
    runtime.GC()
    runtime.ReadMemStats(&m)
    fmt.Printf("%s: Alloc = %d KB", label, m.Alloc/1024)
    fmt.Printf(", TotalAlloc = %d KB", m.TotalAlloc/1024)
    fmt.Printf(", Sys = %d KB\n", m.Sys/1024)
}

func main() {
    printMemStats("Initial")
    
    // Create a chain of large structs
    chain := createChain(10)
    printMemStats("After creating chain")
    
    // Break the chain to allow garbage collection
    current := chain
    for current != nil {
        next := current.Next
        current.Next = nil // Break the reference
        current = next
    }
    chain = nil
    
    printMemStats("After breaking chain")
    
    runtime.GC()
    printMemStats("After GC")
}
```

## Best Practices

### 1. Always Check for Nil

```go
func processPointer(p *int) {
    if p == nil {
        return // or handle the nil case appropriately
    }
    // Safe to dereference
    *p *= 2
}
```

### 2. Use Pointers for Large Structs

```go
type LargeStruct struct {
    // Many fields...
    Data [1000]int
}

// Efficient - passes pointer
func processLargeStruct(ls *LargeStruct) {
    // Process the struct
}

// Inefficient - copies entire struct
func processLargeStructBad(ls LargeStruct) {
    // Process the struct
}
```

### 3. Use Pointers When You Need to Modify

```go
func increment(x *int) {
    *x++
}

func main() {
    num := 5
    increment(&num)
    fmt.Println(num) // 6
}
```

### 4. Prefer Value Types When Possible

```go
// Good for small, immutable data
type Point struct {
    X, Y int
}

func distance(p1, p2 Point) float64 {
    // Calculate distance
    return 0.0 // placeholder
}
```

## Exercises

### Exercise 1: Pointer Basics

Write a program that:
1. Creates an integer variable with value 10
2. Creates a pointer to that variable
3. Modifies the value through the pointer
4. Prints both the original variable and the pointer value

### Exercise 2: Function with Pointers

Write a function that takes two integer pointers and swaps their values. Test it with different values.

### Exercise 3: Struct Pointers

Create a `Student` struct with fields for name, age, and grade. Write functions to:
1. Create a new student (return pointer)
2. Update student information (using pointer parameter)
3. Display student information

### Exercise 4: Linked List Operations

Extend the linked list example to include:
1. A function to find a node with specific data
2. A function to insert a node at a specific position
3. A function to reverse the linked list

### Exercise 5: Pointer Array

Create an array of pointers to integers. Initialize each pointer to point to a different integer value. Write a function that sorts the array based on the values pointed to by the pointers.

## Key Takeaways

1. **Pointers store memory addresses**: They don't hold values directly but point to where values are stored.

2. **Use `&` to get address and `*` to dereference**: These are the fundamental pointer operators.

3. **Nil is the zero value**: Always check for nil before dereferencing.

4. **Pointers enable efficient parameter passing**: Avoid copying large structs by passing pointers.

5. **Pointers allow modification**: Functions can modify original variables through pointers.

6. **Go pointers are safe**: No pointer arithmetic (except with unsafe package).

7. **Automatic dereferencing for structs**: `ptr.field` is equivalent to `(*ptr).field`.

8. **Memory management**: Go's garbage collector handles memory, but be aware of reference cycles.

9. **Use pointers judiciously**: Not everything needs to be a pointer; value types are often simpler and safer.

10. **Slices and maps are reference types**: They already contain pointers internally.

## Next Steps

Now that you understand pointers, let's learn about [Methods](13-methods.md) to see how pointers work with method receivers!

---

**Previous**: [← Structs](11-structs.md) | **Next**: [Methods →](13-methods.md)