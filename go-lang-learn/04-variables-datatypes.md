# 4. Variables and Data Types

## Variable Declaration

Go provides several ways to declare variables:

### 1. var Keyword

```go
// Basic syntax
var variableName dataType

// Examples
var name string
var age int
var isActive bool
```

### 2. var with Initialization

```go
// With explicit type
var name string = "John"
var age int = 25

// Type inference
var name = "John"    // Go infers string type
var age = 25         // Go infers int type
```

### 3. Short Variable Declaration (:=)

```go
// Only inside functions
name := "John"
age := 25
isActive := true
```

### 4. Multiple Variable Declaration

```go
// Multiple variables of same type
var x, y, z int

// Multiple variables with initialization
var name, city string = "John", "New York"

// Multiple variables with type inference
var name, age, isActive = "John", 25, true

// Short declaration
name, age := "John", 25
```

## Basic Data Types

### Numeric Types

#### Integers

| Type | Size | Range |
|------|------|-------|
| `int8` | 8 bits | -128 to 127 |
| `int16` | 16 bits | -32,768 to 32,767 |
| `int32` | 32 bits | -2,147,483,648 to 2,147,483,647 |
| `int64` | 64 bits | -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807 |
| `int` | Platform dependent | 32 or 64 bits |

#### Unsigned Integers

| Type | Size | Range |
|------|------|-------|
| `uint8` | 8 bits | 0 to 255 |
| `uint16` | 16 bits | 0 to 65,535 |
| `uint32` | 32 bits | 0 to 4,294,967,295 |
| `uint64` | 64 bits | 0 to 18,446,744,073,709,551,615 |
| `uint` | Platform dependent | 32 or 64 bits |

#### Special Integer Types

| Type | Description |
|------|-------------|
| `byte` | Alias for `uint8` |
| `rune` | Alias for `int32` (Unicode code point) |
| `uintptr` | Integer type to hold pointer |

#### Floating Point

| Type | Size | Precision |
|------|------|----------|
| `float32` | 32 bits | ~7 decimal digits |
| `float64` | 64 bits | ~15 decimal digits |

#### Complex Numbers

| Type | Description |
|------|-------------|
| `complex64` | Complex number with float32 real and imaginary parts |
| `complex128` | Complex number with float64 real and imaginary parts |

### Examples of Numeric Types

```go
package main

import "fmt"

func main() {
    // Integers
    var age int = 25
    var population int64 = 7800000000
    
    // Unsigned integers
    var count uint = 100
    var data byte = 255
    
    // Floating point
    var price float32 = 19.99
    var pi float64 = 3.14159265359
    
    // Complex numbers
    var c1 complex64 = 1 + 2i
    var c2 complex128 = complex(3, 4) // 3 + 4i
    
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Population: %d\n", population)
    fmt.Printf("Count: %d\n", count)
    fmt.Printf("Data: %d\n", data)
    fmt.Printf("Price: %.2f\n", price)
    fmt.Printf("Pi: %.10f\n", pi)
    fmt.Printf("Complex1: %v\n", c1)
    fmt.Printf("Complex2: %v\n", c2)
}
```

### Boolean Type

```go
var isActive bool = true
var isComplete bool = false

// Short declaration
isValid := true
isEmpty := false

// Boolean operations
result := true && false  // false
result = true || false   // true
result = !true          // false
```

### String Type

```go
// String declaration
var name string = "John Doe"
var message string = `This is a
multi-line string`

// Short declaration
greeting := "Hello, World!"

// String operations
fullName := "John" + " " + "Doe"  // Concatenation
length := len("Hello")            // Length: 5

// String indexing (read-only)
firstChar := name[0]  // 'J' (byte value)

// Raw strings (backticks)
path := `C:\Users\John\Documents`
regex := `\d+\.\d+`
```

## Zero Values

Go initializes variables with zero values if not explicitly initialized:

```go
package main

import "fmt"

func main() {
    var i int        // 0
    var f float64    // 0.0
    var b bool       // false
    var s string     // "" (empty string)
    var p *int       // nil
    
    fmt.Printf("int: %d\n", i)
    fmt.Printf("float64: %f\n", f)
    fmt.Printf("bool: %t\n", b)
    fmt.Printf("string: '%s'\n", s)
    fmt.Printf("pointer: %v\n", p)
}
```

## Type Conversion

Go requires explicit type conversion:

```go
package main

import "fmt"

func main() {
    var i int = 42
    var f float64 = float64(i)  // Convert int to float64
    var u uint = uint(f)        // Convert float64 to uint
    
    // String conversions
    var s string = string(65)   // "A" (ASCII)
    
    fmt.Printf("i: %d, f: %f, u: %d, s: %s\n", i, f, u, s)
    
    // Common conversions
    x := 3.7
    y := int(x)  // y = 3 (truncated, not rounded)
    
    fmt.Printf("x: %f, y: %d\n", x, y)
}
```

## Variable Scope

### Package Level Variables

```go
package main

import "fmt"

// Package level variables
var globalVar = "I'm global"
var (
    name = "John"
    age  = 25
)

func main() {
    fmt.Println(globalVar)
    fmt.Printf("%s is %d years old\n", name, age)
}
```

### Function Level Variables

```go
func main() {
    // Function level variables
    var localVar = "I'm local"
    
    if true {
        // Block level variables
        blockVar := "I'm in a block"
        fmt.Println(localVar)  // Accessible
        fmt.Println(blockVar)  // Accessible
    }
    
    // fmt.Println(blockVar)  // Error: undefined
}
```

## Constants vs Variables

```go
package main

import "fmt"

func main() {
    // Variables (can be changed)
    var count int = 10
    count = 20  // OK
    
    // Constants (cannot be changed)
    const maxRetries = 3
    // maxRetries = 5  // Error: cannot assign to constant
    
    fmt.Printf("Count: %d, Max Retries: %d\n", count, maxRetries)
}
```

## Type Inference Examples

```go
package main

import "fmt"

func main() {
    // Go infers types
    name := "Alice"           // string
    age := 30                 // int
    height := 5.6             // float64
    isStudent := false        // bool
    grade := 'A'              // rune (int32)
    
    // Check types
    fmt.Printf("name: %T\n", name)
    fmt.Printf("age: %T\n", age)
    fmt.Printf("height: %T\n", height)
    fmt.Printf("isStudent: %T\n", isStudent)
    fmt.Printf("grade: %T\n", grade)
}
```

## Practical Examples

### Example 1: User Information

```go
package main

import "fmt"

func main() {
    // User information
    var firstName string = "John"
    var lastName string = "Doe"
    age := 28
    salary := 75000.50
    isEmployed := true
    
    // Display information
    fmt.Printf("Name: %s %s\n", firstName, lastName)
    fmt.Printf("Age: %d years\n", age)
    fmt.Printf("Salary: $%.2f\n", salary)
    fmt.Printf("Employed: %t\n", isEmployed)
    
    // Calculate annual bonus (10% of salary)
    bonus := salary * 0.10
    fmt.Printf("Annual Bonus: $%.2f\n", bonus)
}
```

### Example 2: Mathematical Calculations

```go
package main

import "fmt"

func main() {
    // Circle calculations
    radius := 5.0
    pi := 3.14159
    
    area := pi * radius * radius
    circumference := 2 * pi * radius
    
    fmt.Printf("Circle with radius %.1f:\n", radius)
    fmt.Printf("Area: %.2f\n", area)
    fmt.Printf("Circumference: %.2f\n", circumference)
    
    // Temperature conversion
    celsius := 25.0
    fahrenheit := (celsius * 9.0 / 5.0) + 32.0
    
    fmt.Printf("\nTemperature:\n")
    fmt.Printf("%.1f°C = %.1f°F\n", celsius, fahrenheit)
}
```

### Example 3: String Operations

```go
package main

import "fmt"

func main() {
    firstName := "John"
    lastName := "Doe"
    
    // String concatenation
    fullName := firstName + " " + lastName
    
    // String length
    nameLength := len(fullName)
    
    // String formatting
    greeting := fmt.Sprintf("Hello, %s!", fullName)
    
    fmt.Printf("First Name: %s\n", firstName)
    fmt.Printf("Last Name: %s\n", lastName)
    fmt.Printf("Full Name: %s\n", fullName)
    fmt.Printf("Name Length: %d characters\n", nameLength)
    fmt.Printf("Greeting: %s\n", greeting)
}
```

## Exercises

### Exercise 1: Variable Declaration
Declare variables for a student's information and display them.

```go
package main

import "fmt"

func main() {
    // Declare variables for:
    // - Student name (string)
    // - Student ID (int)
    // - GPA (float64)
    // - Is graduated (bool)
    
    // Display all information
}
```

### Exercise 2: Type Conversion
Convert between different numeric types and display the results.

```go
package main

import "fmt"

func main() {
    // Start with a float64 value
    originalValue := 42.7
    
    // Convert to different types and display
    // int, int32, int64, float32, uint
}
```

### Exercise 3: Calculator
Create a simple calculator that performs basic operations.

```go
package main

import "fmt"

func main() {
    // Declare two numbers
    num1 := 15.5
    num2 := 4.2
    
    // Perform and display:
    // Addition, Subtraction, Multiplication, Division
}
```

## Solutions

<details>
<summary>Click to see solutions</summary>

### Solution 1:
```go
package main

import "fmt"

func main() {
    // Student information
    var studentName string = "Alice Johnson"
    var studentID int = 12345
    var gpa float64 = 3.85
    var isGraduated bool = false
    
    // Display information
    fmt.Printf("Student Name: %s\n", studentName)
    fmt.Printf("Student ID: %d\n", studentID)
    fmt.Printf("GPA: %.2f\n", gpa)
    fmt.Printf("Graduated: %t\n", isGraduated)
}
```

### Solution 2:
```go
package main

import "fmt"

func main() {
    originalValue := 42.7
    
    // Type conversions
    intValue := int(originalValue)
    int32Value := int32(originalValue)
    int64Value := int64(originalValue)
    float32Value := float32(originalValue)
    uintValue := uint(originalValue)
    
    fmt.Printf("Original (float64): %.2f\n", originalValue)
    fmt.Printf("int: %d\n", intValue)
    fmt.Printf("int32: %d\n", int32Value)
    fmt.Printf("int64: %d\n", int64Value)
    fmt.Printf("float32: %.2f\n", float32Value)
    fmt.Printf("uint: %d\n", uintValue)
}
```

### Solution 3:
```go
package main

import "fmt"

func main() {
    num1 := 15.5
    num2 := 4.2
    
    addition := num1 + num2
    subtraction := num1 - num2
    multiplication := num1 * num2
    division := num1 / num2
    
    fmt.Printf("Number 1: %.2f\n", num1)
    fmt.Printf("Number 2: %.2f\n", num2)
    fmt.Printf("Addition: %.2f + %.2f = %.2f\n", num1, num2, addition)
    fmt.Printf("Subtraction: %.2f - %.2f = %.2f\n", num1, num2, subtraction)
    fmt.Printf("Multiplication: %.2f * %.2f = %.2f\n", num1, num2, multiplication)
    fmt.Printf("Division: %.2f / %.2f = %.2f\n", num1, num2, division)
}
```

</details>

## Key Takeaways

1. Go has strong static typing
2. Variables must be declared before use
3. Go provides type inference with `:=`
4. Zero values are automatically assigned
5. Type conversion must be explicit
6. Variable scope determines accessibility
7. Use meaningful variable names
8. Choose appropriate data types for your needs

## Next Steps

Now that you understand variables and data types, let's learn about [Constants](05-constants.md)!

---

**Previous**: [← Basic Syntax](03-basic-syntax.md) | **Next**: [Constants →](05-constants.md)