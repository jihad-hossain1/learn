# 8. Functions

Functions are reusable blocks of code that perform specific tasks. They help organize code, reduce repetition, and make programs more modular and maintainable.

## Function Declaration

### Basic Syntax

```go
func functionName(parameters) returnType {
    // function body
    return value
}
```

### Simple Function Examples

```go
package main

import "fmt"

// Function with no parameters and no return value
func sayHello() {
    fmt.Println("Hello, World!")
}

// Function with parameters and return value
func add(a int, b int) int {
    return a + b
}

// Function with same type parameters (shorthand)
func multiply(a, b int) int {
    return a * b
}

// Function with multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("division by zero")
    }
    return a / b, nil
}

func main() {
    sayHello()
    
    sum := add(5, 3)
    fmt.Printf("5 + 3 = %d\n", sum)
    
    product := multiply(4, 6)
    fmt.Printf("4 * 6 = %d\n", product)
    
    result, err := divide(10, 2)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("10 / 2 = %.2f\n", result)
    }
}
```

## Function Parameters

### Pass by Value

Go passes arguments by value (copies the value).

```go
package main

import "fmt"

func modifyValue(x int) {
    x = 100
    fmt.Printf("Inside function: x = %d\n", x)
}

func main() {
    num := 42
    fmt.Printf("Before function call: num = %d\n", num)
    
    modifyValue(num)
    
    fmt.Printf("After function call: num = %d\n", num)
    // num is still 42 because Go passes by value
}
```

### Pass by Reference (Pointers)

To modify the original value, use pointers.

```go
package main

import "fmt"

func modifyValueByPointer(x *int) {
    *x = 100
    fmt.Printf("Inside function: *x = %d\n", *x)
}

func main() {
    num := 42
    fmt.Printf("Before function call: num = %d\n", num)
    
    modifyValueByPointer(&num)
    
    fmt.Printf("After function call: num = %d\n", num)
    // num is now 100 because we passed a pointer
}
```

### Variadic Functions

Functions that accept a variable number of arguments.

```go
package main

import "fmt"

// Variadic function
func sum(numbers ...int) int {
    total := 0
    for _, num := range numbers {
        total += num
    }
    return total
}

// Mixed parameters (fixed + variadic)
func greetPeople(greeting string, names ...string) {
    for _, name := range names {
        fmt.Printf("%s, %s!\n", greeting, name)
    }
}

func main() {
    // Call with different number of arguments
    fmt.Printf("Sum of no numbers: %d\n", sum())
    fmt.Printf("Sum of 1, 2, 3: %d\n", sum(1, 2, 3))
    fmt.Printf("Sum of 1, 2, 3, 4, 5: %d\n", sum(1, 2, 3, 4, 5))
    
    // Pass slice to variadic function
    numbers := []int{10, 20, 30, 40}
    fmt.Printf("Sum of slice: %d\n", sum(numbers...))
    
    // Mixed parameters
    greetPeople("Hello", "Alice", "Bob", "Carol")
}
```

## Return Values

### Single Return Value

```go
func square(x int) int {
    return x * x
}
```

### Multiple Return Values

```go
func divmod(a, b int) (int, int) {
    quotient := a / b
    remainder := a % b
    return quotient, remainder
}
```

### Named Return Values

```go
package main

import "fmt"

// Named return values
func rectangle(length, width float64) (area, perimeter float64) {
    area = length * width
    perimeter = 2 * (length + width)
    return // naked return
}

// Named returns with explicit return
func circle(radius float64) (area, circumference float64) {
    const pi = 3.14159
    area = pi * radius * radius
    circumference = 2 * pi * radius
    return area, circumference // explicit return
}

func main() {
    area, perimeter := rectangle(5, 3)
    fmt.Printf("Rectangle: Area = %.2f, Perimeter = %.2f\n", area, perimeter)
    
    area, circumference := circle(5)
    fmt.Printf("Circle: Area = %.2f, Circumference = %.2f\n", area, circumference)
}
```

### Ignoring Return Values

```go
package main

import "fmt"

func getNameAndAge() (string, int) {
    return "Alice", 25
}

func main() {
    // Use both return values
    name, age := getNameAndAge()
    fmt.Printf("Name: %s, Age: %d\n", name, age)
    
    // Ignore age using blank identifier
    name, _ = getNameAndAge()
    fmt.Printf("Name only: %s\n", name)
    
    // Ignore name
    _, age = getNameAndAge()
    fmt.Printf("Age only: %d\n", age)
}
```

## Function Types and Variables

### Function as a Type

```go
package main

import "fmt"

// Define a function type
type MathOperation func(int, int) int

// Functions that match the type
func add(a, b int) int {
    return a + b
}

func subtract(a, b int) int {
    return a - b
}

func multiply(a, b int) int {
    return a * b
}

// Function that takes another function as parameter
func calculate(a, b int, operation MathOperation) int {
    return operation(a, b)
}

func main() {
    // Assign functions to variables
    var op MathOperation
    
    op = add
    fmt.Printf("10 + 5 = %d\n", op(10, 5))
    
    op = subtract
    fmt.Printf("10 - 5 = %d\n", op(10, 5))
    
    // Pass function as argument
    result := calculate(10, 5, multiply)
    fmt.Printf("10 * 5 = %d\n", result)
}
```

### Anonymous Functions (Lambdas)

```go
package main

import "fmt"

func main() {
    // Anonymous function assigned to variable
    square := func(x int) int {
        return x * x
    }
    
    fmt.Printf("Square of 5: %d\n", square(5))
    
    // Anonymous function called immediately
    result := func(a, b int) int {
        return a + b
    }(10, 20)
    
    fmt.Printf("Immediate function result: %d\n", result)
    
    // Anonymous function in slice
    operations := []func(int, int) int{
        func(a, b int) int { return a + b },
        func(a, b int) int { return a - b },
        func(a, b int) int { return a * b },
    }
    
    operationNames := []string{"Addition", "Subtraction", "Multiplication"}
    
    for i, op := range operations {
        result := op(8, 3)
        fmt.Printf("%s: 8 and 3 = %d\n", operationNames[i], result)
    }
}
```

## Closures

Closures are functions that capture variables from their surrounding scope.

```go
package main

import "fmt"

// Function that returns a closure
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

// Function that returns a closure with parameter
func multiplier(factor int) func(int) int {
    return func(x int) int {
        return x * factor
    }
}

func main() {
    // Create counter closures
    counter1 := counter()
    counter2 := counter()
    
    fmt.Printf("Counter1: %d\n", counter1()) // 1
    fmt.Printf("Counter1: %d\n", counter1()) // 2
    fmt.Printf("Counter2: %d\n", counter2()) // 1
    fmt.Printf("Counter1: %d\n", counter1()) // 3
    
    // Create multiplier closures
    double := multiplier(2)
    triple := multiplier(3)
    
    fmt.Printf("Double 5: %d\n", double(5))   // 10
    fmt.Printf("Triple 5: %d\n", triple(5))   // 15
    fmt.Printf("Double 10: %d\n", double(10)) // 20
}
```

## Recursive Functions

Functions that call themselves.

```go
package main

import "fmt"

// Factorial using recursion
func factorial(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorial(n-1)
}

// Fibonacci using recursion
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}

// Binary search using recursion
func binarySearch(arr []int, target, left, right int) int {
    if left > right {
        return -1 // Not found
    }
    
    mid := left + (right-left)/2
    
    if arr[mid] == target {
        return mid
    } else if arr[mid] > target {
        return binarySearch(arr, target, left, mid-1)
    } else {
        return binarySearch(arr, target, mid+1, right)
    }
}

func main() {
    // Factorial examples
    for i := 0; i <= 5; i++ {
        fmt.Printf("Factorial of %d: %d\n", i, factorial(i))
    }
    
    // Fibonacci examples
    fmt.Println("\nFibonacci sequence:")
    for i := 0; i < 10; i++ {
        fmt.Printf("%d ", fibonacci(i))
    }
    fmt.Println()
    
    // Binary search example
    arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
    target := 7
    index := binarySearch(arr, target, 0, len(arr)-1)
    
    if index != -1 {
        fmt.Printf("\nFound %d at index %d\n", target, index)
    } else {
        fmt.Printf("\n%d not found\n", target)
    }
}
```

## Defer Statement

The `defer` statement postpones function execution until the surrounding function returns.

```go
package main

import "fmt"

func deferExample() {
    fmt.Println("Start of function")
    
    defer fmt.Println("Deferred statement 1")
    defer fmt.Println("Deferred statement 2")
    defer fmt.Println("Deferred statement 3")
    
    fmt.Println("Middle of function")
    fmt.Println("End of function")
    // Deferred statements execute in LIFO order
}

// Practical use of defer
func processFile(filename string) error {
    fmt.Printf("Opening file: %s\n", filename)
    
    // Simulate file opening
    file := "file_handle"
    
    // Defer closing the file
    defer func() {
        fmt.Printf("Closing file: %s\n", filename)
        // Close file here
    }()
    
    fmt.Printf("Processing file: %s\n", filename)
    
    // File processing logic here
    
    return nil
}

// Defer with function parameters
func deferWithParams() {
    x := 10
    
    defer func(val int) {
        fmt.Printf("Deferred: x was %d\n", val)
    }(x) // x is evaluated now (10)
    
    x = 20
    fmt.Printf("Current: x is %d\n", x)
}

func main() {
    fmt.Println("=== Defer Example ===")
    deferExample()
    
    fmt.Println("\n=== File Processing ===")
    processFile("data.txt")
    
    fmt.Println("\n=== Defer with Parameters ===")
    deferWithParams()
}
```

## Error Handling with Functions

```go
package main

import (
    "errors"
    "fmt"
)

// Function that returns an error
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Function that validates input
func validateAge(age int) error {
    if age < 0 {
        return fmt.Errorf("age cannot be negative: %d", age)
    }
    if age > 150 {
        return fmt.Errorf("age seems unrealistic: %d", age)
    }
    return nil
}

// Function with multiple error conditions
func processUser(name string, age int) error {
    if name == "" {
        return errors.New("name cannot be empty")
    }
    
    if err := validateAge(age); err != nil {
        return fmt.Errorf("age validation failed: %w", err)
    }
    
    fmt.Printf("Processing user: %s, age %d\n", name, age)
    return nil
}

func main() {
    // Division examples
    result, err := divide(10, 2)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("10 / 2 = %.2f\n", result)
    }
    
    result, err = divide(10, 0)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("10 / 0 = %.2f\n", result)
    }
    
    // User processing examples
    users := []struct {
        name string
        age  int
    }{
        {"Alice", 25},
        {"", 30},
        {"Bob", -5},
        {"Carol", 200},
    }
    
    for _, user := range users {
        if err := processUser(user.name, user.age); err != nil {
            fmt.Printf("Failed to process user: %v\n", err)
        }
    }
}
```

## Practical Examples

### Example 1: String Utilities

```go
package main

import (
    "fmt"
    "strings"
)

// Check if string is palindrome
func isPalindrome(s string) bool {
    s = strings.ToLower(s)
    left, right := 0, len(s)-1
    
    for left < right {
        if s[left] != s[right] {
            return false
        }
        left++
        right--
    }
    return true
}

// Count words in a string
func countWords(s string) int {
    words := strings.Fields(s)
    return len(words)
}

// Reverse a string
func reverseString(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

// Capitalize first letter of each word
func titleCase(s string) string {
    words := strings.Fields(s)
    for i, word := range words {
        if len(word) > 0 {
            words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
        }
    }
    return strings.Join(words, " ")
}

func main() {
    testStrings := []string{
        "racecar",
        "hello world",
        "A man a plan a canal Panama",
        "go programming language",
    }
    
    for _, s := range testStrings {
        fmt.Printf("String: \"%s\"\n", s)
        fmt.Printf("  Is palindrome: %t\n", isPalindrome(s))
        fmt.Printf("  Word count: %d\n", countWords(s))
        fmt.Printf("  Reversed: \"%s\"\n", reverseString(s))
        fmt.Printf("  Title case: \"%s\"\n", titleCase(s))
        fmt.Println()
    }
}
```

### Example 2: Mathematical Functions

```go
package main

import (
    "fmt"
    "math"
)

// Calculate greatest common divisor
func gcd(a, b int) int {
    for b != 0 {
        a, b = b, a%b
    }
    return a
}

// Calculate least common multiple
func lcm(a, b int) int {
    return (a * b) / gcd(a, b)
}

// Check if number is prime
func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i <= int(math.Sqrt(float64(n))); i++ {
        if n%i == 0 {
            return false
        }
    }
    return true
}

// Generate prime numbers up to n
func generatePrimes(n int) []int {
    var primes []int
    for i := 2; i <= n; i++ {
        if isPrime(i) {
            primes = append(primes, i)
        }
    }
    return primes
}

// Calculate power using recursion
func power(base, exp int) int {
    if exp == 0 {
        return 1
    }
    if exp == 1 {
        return base
    }
    if exp%2 == 0 {
        half := power(base, exp/2)
        return half * half
    }
    return base * power(base, exp-1)
}

func main() {
    // GCD and LCM examples
    a, b := 48, 18
    fmt.Printf("GCD of %d and %d: %d\n", a, b, gcd(a, b))
    fmt.Printf("LCM of %d and %d: %d\n", a, b, lcm(a, b))
    
    // Prime number examples
    fmt.Printf("\nPrime numbers up to 30: %v\n", generatePrimes(30))
    
    // Power examples
    fmt.Printf("\nPower calculations:\n")
    fmt.Printf("2^10 = %d\n", power(2, 10))
    fmt.Printf("3^5 = %d\n", power(3, 5))
    fmt.Printf("5^0 = %d\n", power(5, 0))
}
```

## Exercises

### Exercise 1: Temperature Converter
Create functions to convert between different temperature scales.

```go
package main

import "fmt"

// Implement these functions:
// celsiusToFahrenheit(c float64) float64
// fahrenheitToCelsius(f float64) float64
// celsiusToKelvin(c float64) float64
// kelvinToCelsius(k float64) float64
// fahrenheitToKelvin(f float64) float64
// kelvinToFahrenheit(k float64) float64

func main() {
    // Test all conversion functions
    // Display a temperature conversion table
}
```

### Exercise 2: Array/Slice Utilities
Create utility functions for working with slices.

```go
package main

import "fmt"

// Implement these functions:
// findMax(numbers []int) int
// findMin(numbers []int) int
// average(numbers []float64) float64
// contains(slice []int, value int) bool
// reverse(slice []int) []int
// removeDuplicates(slice []int) []int

func main() {
    // Test all utility functions
}
```

### Exercise 3: Calculator with Functions
Create a calculator using functions for different operations.

```go
package main

import "fmt"

// Implement calculator functions and a menu system
// Include: add, subtract, multiply, divide, power, sqrt, factorial

func main() {
    // Implement menu-driven calculator
}
```

## Key Takeaways

1. **Functions** organize code into reusable blocks
2. **Parameters** pass data into functions
3. **Return values** pass data out of functions
4. **Multiple returns** are common in Go for error handling
5. **Named returns** can improve readability
6. **Variadic functions** accept variable arguments
7. **Function types** enable higher-order functions
8. **Closures** capture variables from surrounding scope
9. **Recursion** enables elegant solutions for certain problems
10. **Defer** ensures cleanup code runs
11. **Error handling** is explicit and important

## Next Steps

Now that you understand functions, let's learn about [Arrays and Slices](09-arrays-slices.md) to work with collections of data!

---

**Previous**: [← Control Structures](07-control-structures.md) | **Next**: [Arrays and Slices →](09-arrays-slices.md)