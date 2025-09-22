# 5. Constants

## What are Constants?

Constants are immutable values that are known at compile time and do not change during program execution. In Go, constants are declared using the `const` keyword.

## Declaring Constants

### Basic Syntax

```go
const constantName = value
const constantName dataType = value
```

### Examples

```go
package main

import "fmt"

func main() {
    // Basic constants
    const name = "John Doe"
    const age = 25
    const pi = 3.14159
    const isActive = true
    
    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Pi: %.5f\n", pi)
    fmt.Printf("Active: %t\n", isActive)
}
```

## Multiple Constants Declaration

### Method 1: Individual Declaration

```go
const name = "Alice"
const age = 30
const city = "New York"
```

### Method 2: Grouped Declaration

```go
const (
    name = "Alice"
    age  = 30
    city = "New York"
)
```

### Method 3: Multiple Constants in One Line

```go
const name, age, city = "Alice", 30, "New York"
```

## Typed vs Untyped Constants

### Untyped Constants (Preferred)

```go
const (
    name = "John"        // untyped string
    age = 25             // untyped integer
    pi = 3.14159         // untyped floating-point
    isValid = true       // untyped boolean
)
```

### Typed Constants

```go
const (
    name string = "John"
    age int = 25
    pi float64 = 3.14159
    isValid bool = true
)
```

### Why Untyped Constants are Better

```go
package main

import "fmt"

func main() {
    const untypedNum = 42
    const typedNum int = 42
    
    var int32Var int32 = 10
    var int64Var int64 = 20
    
    // Untyped constant can be used with different types
    result1 := int32Var + untypedNum  // OK
    result2 := int64Var + untypedNum  // OK
    
    // Typed constant requires conversion
    // result3 := int32Var + typedNum     // Error
    result3 := int32Var + int32(typedNum) // OK with conversion
    
    fmt.Printf("Result1: %d\n", result1)
    fmt.Printf("Result2: %d\n", result2)
    fmt.Printf("Result3: %d\n", result3)
}
```

## Constant Rules and Limitations

### 1. Constants Must Be Known at Compile Time

```go
// Valid constants
const validConst1 = 42
const validConst2 = "hello"
const validConst3 = 3.14 * 2
const validConst4 = len("hello")  // len() is allowed

// Invalid constants (runtime values)
// const invalidConst1 = time.Now()     // Error
// const invalidConst2 = math.Sin(1.0)  // Error
// const invalidConst3 = rand.Intn(100) // Error
```

### 2. Constants Cannot Be Changed

```go
const maxRetries = 3
// maxRetries = 5  // Error: cannot assign to constant
```

### 3. Limited to Basic Types

```go
// Valid constant types
const str = "hello"           // string
const num = 42               // numeric
const flag = true            // boolean

// Invalid constant types
// const arr = [3]int{1, 2, 3}    // Error: arrays not allowed
// const slice = []int{1, 2, 3}   // Error: slices not allowed
// const m = map[string]int{}     // Error: maps not allowed
```

## iota - Constant Generator

`iota` is a predeclared identifier that represents successive untyped integer constants.

### Basic iota Usage

```go
package main

import "fmt"

const (
    a = iota  // 0
    b = iota  // 1
    c = iota  // 2
)

// Simplified version
const (
    x = iota  // 0
    y         // 1 (implicitly y = iota)
    z         // 2 (implicitly z = iota)
)

func main() {
    fmt.Printf("a=%d, b=%d, c=%d\n", a, b, c)
    fmt.Printf("x=%d, y=%d, z=%d\n", x, y, z)
}
```

### iota with Expressions

```go
package main

import "fmt"

const (
    _  = iota             // 0 (ignored)
    KB = 1 << (10 * iota) // 1 << (10 * 1) = 1024
    MB = 1 << (10 * iota) // 1 << (10 * 2) = 1048576
    GB = 1 << (10 * iota) // 1 << (10 * 3) = 1073741824
)

func main() {
    fmt.Printf("KB = %d\n", KB)
    fmt.Printf("MB = %d\n", MB)
    fmt.Printf("GB = %d\n", GB)
}
```

### iota Resets in Each const Block

```go
package main

import "fmt"

const (
    first = iota   // 0
    second         // 1
)

const (
    third = iota   // 0 (resets)
    fourth         // 1
)

func main() {
    fmt.Printf("first=%d, second=%d\n", first, second)
    fmt.Printf("third=%d, fourth=%d\n", third, fourth)
}
```

## Practical Examples

### Example 1: Configuration Constants

```go
package main

import "fmt"

const (
    // Application configuration
    AppName    = "MyApp"
    AppVersion = "1.0.0"
    
    // Server configuration
    DefaultPort = 8080
    MaxRetries  = 3
    Timeout     = 30 // seconds
    
    // Database configuration
    DBHost = "localhost"
    DBPort = 5432
    DBName = "myapp_db"
)

func main() {
    fmt.Printf("%s v%s\n", AppName, AppVersion)
    fmt.Printf("Server: localhost:%d\n", DefaultPort)
    fmt.Printf("Database: %s:%d/%s\n", DBHost, DBPort, DBName)
    fmt.Printf("Max retries: %d, Timeout: %ds\n", MaxRetries, Timeout)
}
```

### Example 2: Status Codes with iota

```go
package main

import "fmt"

const (
    StatusPending = iota
    StatusInProgress
    StatusCompleted
    StatusFailed
    StatusCancelled
)

func getStatusName(status int) string {
    switch status {
    case StatusPending:
        return "Pending"
    case StatusInProgress:
        return "In Progress"
    case StatusCompleted:
        return "Completed"
    case StatusFailed:
        return "Failed"
    case StatusCancelled:
        return "Cancelled"
    default:
        return "Unknown"
    }
}

func main() {
    fmt.Printf("Status codes:\n")
    fmt.Printf("Pending: %d (%s)\n", StatusPending, getStatusName(StatusPending))
    fmt.Printf("In Progress: %d (%s)\n", StatusInProgress, getStatusName(StatusInProgress))
    fmt.Printf("Completed: %d (%s)\n", StatusCompleted, getStatusName(StatusCompleted))
    fmt.Printf("Failed: %d (%s)\n", StatusFailed, getStatusName(StatusFailed))
    fmt.Printf("Cancelled: %d (%s)\n", StatusCancelled, getStatusName(StatusCancelled))
}
```

### Example 3: File Permissions with iota

```go
package main

import "fmt"

const (
    ReadPermission = 1 << iota  // 1 (binary: 001)
    WritePermission             // 2 (binary: 010)
    ExecutePermission           // 4 (binary: 100)
)

func hasPermission(userPermissions, requiredPermission int) bool {
    return userPermissions&requiredPermission != 0
}

func main() {
    // User has read and write permissions
    userPerms := ReadPermission | WritePermission
    
    fmt.Printf("User permissions: %d (binary: %b)\n", userPerms, userPerms)
    fmt.Printf("Can read: %t\n", hasPermission(userPerms, ReadPermission))
    fmt.Printf("Can write: %t\n", hasPermission(userPerms, WritePermission))
    fmt.Printf("Can execute: %t\n", hasPermission(userPerms, ExecutePermission))
    
    // Grant execute permission
    userPerms |= ExecutePermission
    fmt.Printf("\nAfter granting execute permission: %d (binary: %b)\n", userPerms, userPerms)
    fmt.Printf("Can execute: %t\n", hasPermission(userPerms, ExecutePermission))
}
```

### Example 4: Mathematical Constants

```go
package main

import "fmt"

const (
    // Mathematical constants
    Pi    = 3.14159265359
    E     = 2.71828182846
    Phi   = 1.61803398875  // Golden ratio
    
    // Physical constants
    SpeedOfLight = 299792458 // m/s
    Gravity      = 9.80665   // m/s²
)

func main() {
    radius := 5.0
    
    // Circle calculations
    area := Pi * radius * radius
    circumference := 2 * Pi * radius
    
    fmt.Printf("Circle with radius %.1f:\n", radius)
    fmt.Printf("Area: %.2f\n", area)
    fmt.Printf("Circumference: %.2f\n", circumference)
    
    // Physics calculation
    mass := 10.0 // kg
    weight := mass * Gravity
    fmt.Printf("\nWeight of %.1f kg: %.2f N\n", mass, weight)
}
```

## Best Practices

### 1. Use Meaningful Names

```go
// Good
const (
    MaxRetries = 3
    TimeoutSeconds = 30
    DefaultPort = 8080
)

// Bad
const (
    X = 3
    Y = 30
    Z = 8080
)
```

### 2. Group Related Constants

```go
// Good - grouped by functionality
const (
    // HTTP status codes
    StatusOK       = 200
    StatusNotFound = 404
    StatusError    = 500
)

const (
    // Database configuration
    DBHost = "localhost"
    DBPort = 5432
    DBName = "myapp"
)
```

### 3. Use iota for Enumerations

```go
// Good - using iota for related constants
const (
    Monday = iota
    Tuesday
    Wednesday
    Thursday
    Friday
    Saturday
    Sunday
)

// Bad - manual numbering
const (
    Monday    = 0
    Tuesday   = 1
    Wednesday = 2
    // ... error-prone
)
```

### 4. Document Complex Constants

```go
const (
    // MaxFileSize represents the maximum allowed file size in bytes (10MB)
    MaxFileSize = 10 * 1024 * 1024
    
    // TokenExpiryHours defines how long authentication tokens remain valid
    TokenExpiryHours = 24
)
```

## Exercises

### Exercise 1: Basic Constants
Create constants for a simple e-commerce application.

```go
package main

import "fmt"

func main() {
    // Define constants for:
    // - Store name
    // - Tax rate (as percentage)
    // - Free shipping threshold
    // - Maximum items per order
    
    // Calculate total for an order
    itemPrice := 25.99
    quantity := 3
    
    // Use your constants to calculate:
    // - Subtotal
    // - Tax amount
    // - Total
    // - Whether free shipping applies
}
```

### Exercise 2: Status Codes with iota
Create an order status system using iota.

```go
package main

import "fmt"

func main() {
    // Define order status constants using iota:
    // OrderPlaced, OrderConfirmed, OrderShipped, OrderDelivered, OrderCancelled
    
    // Create a function that returns status description
    
    // Test with different status values
}
```

### Exercise 3: File Size Constants
Create file size constants using iota and bit shifting.

```go
package main

import "fmt"

func main() {
    // Define file size constants:
    // Byte, KB, MB, GB, TB using iota and bit shifting
    
    // Convert and display file sizes
    fileSize := 2.5 * 1024 * 1024 * 1024 // 2.5 GB in bytes
    
    // Display in different units
}
```

## Solutions

<details>
<summary>Click to see solutions</summary>

### Solution 1:
```go
package main

import "fmt"

const (
    StoreName = "TechStore"
    TaxRate = 8.5 // percentage
    FreeShippingThreshold = 50.0
    MaxItemsPerOrder = 10
)

func main() {
    itemPrice := 25.99
    quantity := 3
    
    subtotal := itemPrice * float64(quantity)
    taxAmount := subtotal * TaxRate / 100
    total := subtotal + taxAmount
    freeShipping := subtotal >= FreeShippingThreshold
    
    fmt.Printf("Welcome to %s!\n\n", StoreName)
    fmt.Printf("Item Price: $%.2f\n", itemPrice)
    fmt.Printf("Quantity: %d\n", quantity)
    fmt.Printf("Subtotal: $%.2f\n", subtotal)
    fmt.Printf("Tax (%.1f%%): $%.2f\n", TaxRate, taxAmount)
    fmt.Printf("Total: $%.2f\n", total)
    fmt.Printf("Free Shipping: %t\n", freeShipping)
}
```

### Solution 2:
```go
package main

import "fmt"

const (
    OrderPlaced = iota
    OrderConfirmed
    OrderShipped
    OrderDelivered
    OrderCancelled
)

func getOrderStatus(status int) string {
    switch status {
    case OrderPlaced:
        return "Order Placed"
    case OrderConfirmed:
        return "Order Confirmed"
    case OrderShipped:
        return "Order Shipped"
    case OrderDelivered:
        return "Order Delivered"
    case OrderCancelled:
        return "Order Cancelled"
    default:
        return "Unknown Status"
    }
}

func main() {
    fmt.Println("Order Status System:")
    
    for status := OrderPlaced; status <= OrderCancelled; status++ {
        fmt.Printf("Status %d: %s\n", status, getOrderStatus(status))
    }
}
```

### Solution 3:
```go
package main

import "fmt"

const (
    Byte = 1 << (10 * iota)
    KB
    MB
    GB
    TB
)

func main() {
    fmt.Println("File Size Constants:")
    fmt.Printf("Byte: %d\n", Byte)
    fmt.Printf("KB: %d\n", KB)
    fmt.Printf("MB: %d\n", MB)
    fmt.Printf("GB: %d\n", GB)
    fmt.Printf("TB: %d\n", TB)
    
    fileSize := 2.5 * float64(GB) // 2.5 GB in bytes
    
    fmt.Printf("\nFile size: %.0f bytes\n", fileSize)
    fmt.Printf("In KB: %.2f\n", fileSize/float64(KB))
    fmt.Printf("In MB: %.2f\n", fileSize/float64(MB))
    fmt.Printf("In GB: %.2f\n", fileSize/float64(GB))
}
```

</details>

## Key Takeaways

1. Constants are immutable and known at compile time
2. Use untyped constants for flexibility
3. Group related constants together
4. Use `iota` for enumerations and sequential values
5. Constants improve code readability and maintainability
6. Choose meaningful names for constants
7. Document complex constant calculations
8. Constants can only be basic types (string, numeric, boolean)

## Next Steps

Now that you understand constants, let's learn about [Operators](06-operators.md) to perform operations on variables and constants!

---

**Previous**: [← Variables and Data Types](04-variables-datatypes.md) | **Next**: [Operators →](06-operators.md)