# 7. Control Structures

Control structures determine the flow of execution in a program. Go provides several control structures: conditional statements, loops, and switch statements.

## If Statements

The `if` statement executes code based on a boolean condition.

### Basic If Statement

```go
if condition {
    // code to execute if condition is true
}
```

### If-Else Statement

```go
if condition {
    // code if condition is true
} else {
    // code if condition is false
}
```

### If-Else If-Else Statement

```go
if condition1 {
    // code if condition1 is true
} else if condition2 {
    // code if condition2 is true
} else {
    // code if all conditions are false
}
```

### Examples

```go
package main

import "fmt"

func main() {
    age := 18
    
    // Basic if
    if age >= 18 {
        fmt.Println("You are an adult")
    }
    
    // If-else
    if age >= 21 {
        fmt.Println("You can drink alcohol in the US")
    } else {
        fmt.Println("You cannot drink alcohol in the US")
    }
    
    // If-else if-else
    score := 85
    if score >= 90 {
        fmt.Println("Grade: A")
    } else if score >= 80 {
        fmt.Println("Grade: B")
    } else if score >= 70 {
        fmt.Println("Grade: C")
    } else if score >= 60 {
        fmt.Println("Grade: D")
    } else {
        fmt.Println("Grade: F")
    }
}
```

### If with Short Statement

Go allows you to execute a short statement before the condition.

```go
package main

import "fmt"

func main() {
    // Short statement before condition
    if num := 42; num%2 == 0 {
        fmt.Printf("%d is even\n", num)
    } else {
        fmt.Printf("%d is odd\n", num)
    }
    
    // num is not accessible here
    // fmt.Println(num) // Error: undefined
    
    // Practical example: error handling
    if err := someFunction(); err != nil {
        fmt.Printf("Error occurred: %v\n", err)
        return
    }
    fmt.Println("Function executed successfully")
}

func someFunction() error {
    // Simulate a function that might return an error
    return nil // No error
}
```

## For Loops

Go has only one looping construct: the `for` loop. However, it can be used in different ways.

### Traditional For Loop

```go
for initialization; condition; post {
    // loop body
}
```

### While-style Loop

```go
for condition {
    // loop body
}
```

### Infinite Loop

```go
for {
    // infinite loop
    // use break to exit
}
```

### Examples

```go
package main

import "fmt"

func main() {
    // Traditional for loop
    fmt.Println("Counting from 1 to 5:")
    for i := 1; i <= 5; i++ {
        fmt.Printf("%d ", i)
    }
    fmt.Println()
    
    // While-style loop
    fmt.Println("\nCountdown:")
    count := 5
    for count > 0 {
        fmt.Printf("%d ", count)
        count--
    }
    fmt.Println("Blast off!")
    
    // Infinite loop with break
    fmt.Println("\nInfinite loop with break:")
    i := 0
    for {
        if i >= 3 {
            break
        }
        fmt.Printf("Iteration %d\n", i)
        i++
    }
    
    // Loop with continue
    fmt.Println("\nSkipping even numbers:")
    for i := 1; i <= 10; i++ {
        if i%2 == 0 {
            continue // Skip even numbers
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()
}
```

### Range Loop

The `range` keyword is used to iterate over arrays, slices, maps, strings, and channels.

```go
package main

import "fmt"

func main() {
    // Range over slice
    numbers := []int{10, 20, 30, 40, 50}
    
    fmt.Println("Index and value:")
    for index, value := range numbers {
        fmt.Printf("Index: %d, Value: %d\n", index, value)
    }
    
    // Range with only value
    fmt.Println("\nOnly values:")
    for _, value := range numbers {
        fmt.Printf("%d ", value)
    }
    fmt.Println()
    
    // Range with only index
    fmt.Println("\nOnly indices:")
    for index := range numbers {
        fmt.Printf("%d ", index)
    }
    fmt.Println()
    
    // Range over string
    fmt.Println("\nIterating over string:")
    text := "Hello"
    for index, char := range text {
        fmt.Printf("Index: %d, Character: %c\n", index, char)
    }
    
    // Range over map
    fmt.Println("\nIterating over map:")
    ages := map[string]int{
        "Alice": 25,
        "Bob":   30,
        "Carol": 35,
    }
    
    for name, age := range ages {
        fmt.Printf("%s is %d years old\n", name, age)
    }
}
```

## Switch Statements

The `switch` statement is a cleaner way to write multiple if-else statements.

### Basic Switch

```go
switch variable {
case value1:
    // code for value1
case value2:
    // code for value2
default:
    // code if no case matches
}
```

### Examples

```go
package main

import "fmt"

func main() {
    // Basic switch
    day := "Monday"
    
    switch day {
    case "Monday":
        fmt.Println("Start of the work week")
    case "Tuesday", "Wednesday", "Thursday":
        fmt.Println("Midweek")
    case "Friday":
        fmt.Println("TGIF!")
    case "Saturday", "Sunday":
        fmt.Println("Weekend!")
    default:
        fmt.Println("Invalid day")
    }
    
    // Switch with expressions
    score := 85
    
    switch {
    case score >= 90:
        fmt.Println("Grade: A")
    case score >= 80:
        fmt.Println("Grade: B")
    case score >= 70:
        fmt.Println("Grade: C")
    case score >= 60:
        fmt.Println("Grade: D")
    default:
        fmt.Println("Grade: F")
    }
    
    // Switch with short statement
    switch num := 42; {
    case num < 0:
        fmt.Println("Negative")
    case num == 0:
        fmt.Println("Zero")
    case num > 0:
        fmt.Println("Positive")
    }
}
```

### Switch with Fallthrough

By default, Go's switch cases don't fall through. Use `fallthrough` to continue to the next case.

```go
package main

import "fmt"

func main() {
    grade := "B"
    
    fmt.Printf("Grade %s means: ", grade)
    switch grade {
    case "A":
        fmt.Print("Excellent ")
        fallthrough
    case "B":
        fmt.Print("Good ")
        fallthrough
    case "C":
        fmt.Print("Average ")
        fallthrough
    case "D":
        fmt.Print("Below Average ")
        fallthrough
    default:
        fmt.Println("performance")
    }
}
```

### Type Switch

Type switch is used to determine the type of an interface variable.

```go
package main

import "fmt"

func main() {
    var value interface{} = 42
    
    switch v := value.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s\n", v)
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
    
    // Test with different types
    testTypeSwitch(42)
    testTypeSwitch("hello")
    testTypeSwitch(true)
    testTypeSwitch(3.14)
}

func testTypeSwitch(value interface{}) {
    switch v := value.(type) {
    case int:
        fmt.Printf("Integer: %d (doubled: %d)\n", v, v*2)
    case string:
        fmt.Printf("String: %s (length: %d)\n", v, len(v))
    case bool:
        fmt.Printf("Boolean: %t (negated: %t)\n", v, !v)
    default:
        fmt.Printf("Unknown type: %T, value: %v\n", v, v)
    }
}
```

## Break and Continue

### Break Statement

The `break` statement exits the current loop or switch.

```go
package main

import "fmt"

func main() {
    // Break in for loop
    fmt.Println("Finding first even number:")
    numbers := []int{1, 3, 5, 8, 9, 12, 15}
    
    for _, num := range numbers {
        if num%2 == 0 {
            fmt.Printf("Found first even number: %d\n", num)
            break
        }
        fmt.Printf("Checking %d...\n", num)
    }
    
    // Break in nested loops with labels
    fmt.Println("\nBreaking from nested loop:")
Outer:
    for i := 1; i <= 3; i++ {
        for j := 1; j <= 3; j++ {
            if i*j > 4 {
                fmt.Printf("Breaking at i=%d, j=%d\n", i, j)
                break Outer
            }
            fmt.Printf("i=%d, j=%d, product=%d\n", i, j, i*j)
        }
    }
}
```

### Continue Statement

The `continue` statement skips the current iteration and continues with the next.

```go
package main

import "fmt"

func main() {
    // Continue in for loop
    fmt.Println("Printing odd numbers from 1 to 10:")
    for i := 1; i <= 10; i++ {
        if i%2 == 0 {
            continue // Skip even numbers
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()
    
    // Continue with labels
    fmt.Println("\nSkipping multiples of 3:")
Loop:
    for i := 1; i <= 15; i++ {
        if i%3 == 0 {
            continue Loop
        }
        fmt.Printf("%d ", i)
    }
    fmt.Println()
}
```

## Practical Examples

### Example 1: Number Guessing Game

```go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    // Seed random number generator
    rand.Seed(time.Now().UnixNano())
    
    // Generate random number between 1 and 100
    target := rand.Intn(100) + 1
    attempts := 0
    maxAttempts := 7
    
    fmt.Println("Welcome to the Number Guessing Game!")
    fmt.Printf("I'm thinking of a number between 1 and 100.\n")
    fmt.Printf("You have %d attempts to guess it.\n\n", maxAttempts)
    
    for attempts < maxAttempts {
        var guess int
        fmt.Printf("Attempt %d: Enter your guess: ", attempts+1)
        fmt.Scanln(&guess)
        
        attempts++
        
        if guess == target {
            fmt.Printf("Congratulations! You guessed it in %d attempts!\n", attempts)
            break
        } else if guess < target {
            fmt.Println("Too low!")
        } else {
            fmt.Println("Too high!")
        }
        
        if attempts == maxAttempts {
            fmt.Printf("Game over! The number was %d.\n", target)
        } else {
            fmt.Printf("You have %d attempts left.\n\n", maxAttempts-attempts)
        }
    }
}
```

### Example 2: Menu-Driven Program

```go
package main

import "fmt"

func main() {
    var choice int
    
    for {
        // Display menu
        fmt.Println("\n=== Calculator Menu ===")
        fmt.Println("1. Addition")
        fmt.Println("2. Subtraction")
        fmt.Println("3. Multiplication")
        fmt.Println("4. Division")
        fmt.Println("5. Exit")
        fmt.Print("Enter your choice (1-5): ")
        
        fmt.Scanln(&choice)
        
        switch choice {
        case 1:
            performOperation("addition")
        case 2:
            performOperation("subtraction")
        case 3:
            performOperation("multiplication")
        case 4:
            performOperation("division")
        case 5:
            fmt.Println("Thank you for using the calculator!")
            return
        default:
            fmt.Println("Invalid choice! Please enter a number between 1 and 5.")
        }
    }
}

func performOperation(operation string) {
    var num1, num2 float64
    
    fmt.Printf("\nPerforming %s:\n", operation)
    fmt.Print("Enter first number: ")
    fmt.Scanln(&num1)
    fmt.Print("Enter second number: ")
    fmt.Scanln(&num2)
    
    var result float64
    var valid bool = true
    
    switch operation {
    case "addition":
        result = num1 + num2
    case "subtraction":
        result = num1 - num2
    case "multiplication":
        result = num1 * num2
    case "division":
        if num2 != 0 {
            result = num1 / num2
        } else {
            fmt.Println("Error: Division by zero!")
            valid = false
        }
    }
    
    if valid {
        fmt.Printf("Result: %.2f\n", result)
    }
}
```

### Example 3: Pattern Printing

```go
package main

import "fmt"

func main() {
    // Right triangle
    fmt.Println("Right Triangle:")
    for i := 1; i <= 5; i++ {
        for j := 1; j <= i; j++ {
            fmt.Print("* ")
        }
        fmt.Println()
    }
    
    // Inverted triangle
    fmt.Println("\nInverted Triangle:")
    for i := 5; i >= 1; i-- {
        for j := 1; j <= i; j++ {
            fmt.Print("* ")
        }
        fmt.Println()
    }
    
    // Number pyramid
    fmt.Println("\nNumber Pyramid:")
    for i := 1; i <= 5; i++ {
        // Print spaces
        for j := 1; j <= 5-i; j++ {
            fmt.Print(" ")
        }
        // Print numbers
        for j := 1; j <= i; j++ {
            fmt.Printf("%d ", j)
        }
        fmt.Println()
    }
    
    // Multiplication table
    fmt.Println("\nMultiplication Table (1-5):")
    for i := 1; i <= 5; i++ {
        for j := 1; j <= 5; j++ {
            fmt.Printf("%2d ", i*j)
        }
        fmt.Println()
    }
}
```

## Exercises

### Exercise 1: Grade Statistics
Create a program that calculates grade statistics.

```go
package main

import "fmt"

func main() {
    // Read number of students
    // Read grades for each student
    // Calculate and display:
    // - Average grade
    // - Highest grade
    // - Lowest grade
    // - Number of students in each grade category (A, B, C, D, F)
    // - Pass/fail statistics
}
```

### Exercise 2: Prime Number Checker
Create a program that finds all prime numbers up to a given number.

```go
package main

import "fmt"

func main() {
    // Read upper limit from user
    // Find and display all prime numbers up to that limit
    // Display count of prime numbers found
}
```

### Exercise 3: Simple ATM System
Create a simple ATM system with basic operations.

```go
package main

import "fmt"

func main() {
    // Initial balance
    balance := 1000.0
    
    // Menu options:
    // 1. Check Balance
    // 2. Deposit
    // 3. Withdraw
    // 4. Exit
    
    // Implement proper validation and error handling
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
    var numStudents int
    fmt.Print("Enter number of students: ")
    fmt.Scanln(&numStudents)
    
    if numStudents <= 0 {
        fmt.Println("Invalid number of students!")
        return
    }
    
    var total float64
    var highest, lowest float64
    var gradeCount [5]int // A, B, C, D, F
    var passCount int
    
    for i := 1; i <= numStudents; i++ {
        var grade float64
        fmt.Printf("Enter grade for student %d: ", i)
        fmt.Scanln(&grade)
        
        if grade < 0 || grade > 100 {
            fmt.Println("Invalid grade! Please enter a value between 0 and 100.")
            i-- // Retry this student
            continue
        }
        
        total += grade
        
        // Track highest and lowest
        if i == 1 {
            highest = grade
            lowest = grade
        } else {
            if grade > highest {
                highest = grade
            }
            if grade < lowest {
                lowest = grade
            }
        }
        
        // Count grade categories
        switch {
        case grade >= 90:
            gradeCount[0]++ // A
        case grade >= 80:
            gradeCount[1]++ // B
        case grade >= 70:
            gradeCount[2]++ // C
        case grade >= 60:
            gradeCount[3]++ // D
            passCount++
        default:
            gradeCount[4]++ // F
        }
        
        if grade >= 60 {
            passCount++
        }
    }
    
    average := total / float64(numStudents)
    failCount := numStudents - passCount
    
    fmt.Printf("\n=== Grade Statistics ===\n")
    fmt.Printf("Average Grade: %.2f\n", average)
    fmt.Printf("Highest Grade: %.2f\n", highest)
    fmt.Printf("Lowest Grade: %.2f\n", lowest)
    fmt.Printf("\nGrade Distribution:\n")
    fmt.Printf("A (90-100): %d students\n", gradeCount[0])
    fmt.Printf("B (80-89):  %d students\n", gradeCount[1])
    fmt.Printf("C (70-79):  %d students\n", gradeCount[2])
    fmt.Printf("D (60-69):  %d students\n", gradeCount[3])
    fmt.Printf("F (0-59):   %d students\n", gradeCount[4])
    fmt.Printf("\nPass/Fail Statistics:\n")
    fmt.Printf("Passed: %d students (%.1f%%)\n", passCount, float64(passCount)/float64(numStudents)*100)
    fmt.Printf("Failed: %d students (%.1f%%)\n", failCount, float64(failCount)/float64(numStudents)*100)
}
```

### Solution 2:
```go
package main

import "fmt"

func main() {
    var limit int
    fmt.Print("Enter the upper limit: ")
    fmt.Scanln(&limit)
    
    if limit < 2 {
        fmt.Println("No prime numbers found.")
        return
    }
    
    fmt.Printf("Prime numbers up to %d:\n", limit)
    
    primeCount := 0
    
    for num := 2; num <= limit; num++ {
        isPrime := true
        
        // Check if num is prime
        for i := 2; i*i <= num; i++ {
            if num%i == 0 {
                isPrime = false
                break
            }
        }
        
        if isPrime {
            fmt.Printf("%d ", num)
            primeCount++
            
            // Print newline every 10 numbers for better formatting
            if primeCount%10 == 0 {
                fmt.Println()
            }
        }
    }
    
    fmt.Printf("\n\nTotal prime numbers found: %d\n", primeCount)
}
```

### Solution 3:
```go
package main

import "fmt"

func main() {
    balance := 1000.0
    var choice int
    
    fmt.Println("Welcome to Simple ATM System!")
    
    for {
        fmt.Println("\n=== ATM Menu ===")
        fmt.Println("1. Check Balance")
        fmt.Println("2. Deposit")
        fmt.Println("3. Withdraw")
        fmt.Println("4. Exit")
        fmt.Print("Enter your choice (1-4): ")
        
        fmt.Scanln(&choice)
        
        switch choice {
        case 1:
            fmt.Printf("Your current balance: $%.2f\n", balance)
            
        case 2:
            var amount float64
            fmt.Print("Enter deposit amount: $")
            fmt.Scanln(&amount)
            
            if amount <= 0 {
                fmt.Println("Invalid amount! Please enter a positive value.")
            } else {
                balance += amount
                fmt.Printf("Successfully deposited $%.2f\n", amount)
                fmt.Printf("New balance: $%.2f\n", balance)
            }
            
        case 3:
            var amount float64
            fmt.Print("Enter withdrawal amount: $")
            fmt.Scanln(&amount)
            
            if amount <= 0 {
                fmt.Println("Invalid amount! Please enter a positive value.")
            } else if amount > balance {
                fmt.Printf("Insufficient funds! Your balance is $%.2f\n", balance)
            } else {
                balance -= amount
                fmt.Printf("Successfully withdrew $%.2f\n", amount)
                fmt.Printf("Remaining balance: $%.2f\n", balance)
            }
            
        case 4:
            fmt.Println("Thank you for using our ATM service!")
            return
            
        default:
            fmt.Println("Invalid choice! Please enter a number between 1 and 4.")
        }
    }
}
```

</details>

## Key Takeaways

1. **If statements** control conditional execution
2. **For loops** are the only looping construct in Go
3. **Range** provides easy iteration over collections
4. **Switch statements** offer clean multi-way branching
5. **Break** exits loops or switches
6. **Continue** skips to next iteration
7. **Labels** allow breaking/continuing outer loops
8. Go's switch doesn't fall through by default
9. Short statements in if/switch provide local scope
10. Type switches help with interface type checking

## Next Steps

Now that you understand control structures, let's learn about [Functions](08-functions.md) to organize and reuse code!

---

**Previous**: [← Operators](06-operators.md) | **Next**: [Functions →](08-functions.md)