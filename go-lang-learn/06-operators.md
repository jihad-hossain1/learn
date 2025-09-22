# 6. Operators

Operators are symbols that perform operations on operands (variables and values). Go supports various types of operators for different operations.

## Arithmetic Operators

Arithmetic operators perform mathematical operations on numeric values.

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `+` | Addition | `5 + 3` | `8` |
| `-` | Subtraction | `5 - 3` | `2` |
| `*` | Multiplication | `5 * 3` | `15` |
| `/` | Division | `10 / 3` | `3` (integer division) |
| `%` | Modulus (remainder) | `10 % 3` | `1` |

### Examples

```go
package main

import "fmt"

func main() {
    a := 10
    b := 3
    
    fmt.Printf("a = %d, b = %d\n", a, b)
    fmt.Printf("Addition: %d + %d = %d\n", a, b, a+b)
    fmt.Printf("Subtraction: %d - %d = %d\n", a, b, a-b)
    fmt.Printf("Multiplication: %d * %d = %d\n", a, b, a*b)
    fmt.Printf("Division: %d / %d = %d\n", a, b, a/b)
    fmt.Printf("Modulus: %d %% %d = %d\n", a, b, a%b)
    
    // Floating point division
    x := 10.0
    y := 3.0
    fmt.Printf("Float Division: %.2f / %.2f = %.2f\n", x, y, x/y)
}
```

### Integer vs Floating Point Division

```go
package main

import "fmt"

func main() {
    // Integer division
    intResult := 7 / 2  // Result: 3 (truncated)
    
    // Floating point division
    floatResult := 7.0 / 2.0  // Result: 3.5
    
    // Mixed division (convert to float)
    mixedResult := float64(7) / float64(2)  // Result: 3.5
    
    fmt.Printf("Integer division: 7 / 2 = %d\n", intResult)
    fmt.Printf("Float division: 7.0 / 2.0 = %.1f\n", floatResult)
    fmt.Printf("Mixed division: float64(7) / float64(2) = %.1f\n", mixedResult)
}
```

## Assignment Operators

Assignment operators assign values to variables.

| Operator | Description | Example | Equivalent |
|----------|-------------|---------|------------|
| `=` | Simple assignment | `x = 5` | - |
| `+=` | Add and assign | `x += 3` | `x = x + 3` |
| `-=` | Subtract and assign | `x -= 3` | `x = x - 3` |
| `*=` | Multiply and assign | `x *= 3` | `x = x * 3` |
| `/=` | Divide and assign | `x /= 3` | `x = x / 3` |
| `%=` | Modulus and assign | `x %= 3` | `x = x % 3` |

### Examples

```go
package main

import "fmt"

func main() {
    x := 10
    fmt.Printf("Initial value: x = %d\n", x)
    
    x += 5  // x = x + 5
    fmt.Printf("After x += 5: x = %d\n", x)
    
    x -= 3  // x = x - 3
    fmt.Printf("After x -= 3: x = %d\n", x)
    
    x *= 2  // x = x * 2
    fmt.Printf("After x *= 2: x = %d\n", x)
    
    x /= 4  // x = x / 4
    fmt.Printf("After x /= 4: x = %d\n", x)
    
    x %= 3  // x = x % 3
    fmt.Printf("After x %%= 3: x = %d\n", x)
}
```

## Comparison Operators

Comparison operators compare two values and return a boolean result.

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `==` | Equal to | `5 == 5` | `true` |
| `!=` | Not equal to | `5 != 3` | `true` |
| `<` | Less than | `3 < 5` | `true` |
| `<=` | Less than or equal | `5 <= 5` | `true` |
| `>` | Greater than | `5 > 3` | `true` |
| `>=` | Greater than or equal | `5 >= 5` | `true` |

### Examples

```go
package main

import "fmt"

func main() {
    a := 10
    b := 5
    c := 10
    
    fmt.Printf("a = %d, b = %d, c = %d\n", a, b, c)
    fmt.Printf("a == c: %t\n", a == c)
    fmt.Printf("a != b: %t\n", a != b)
    fmt.Printf("a > b: %t\n", a > b)
    fmt.Printf("a >= c: %t\n", a >= c)
    fmt.Printf("b < a: %t\n", b < a)
    fmt.Printf("b <= a: %t\n", b <= a)
    
    // String comparison
    str1 := "apple"
    str2 := "banana"
    fmt.Printf("\nString comparison:\n")
    fmt.Printf("\"%s\" < \"%s\": %t\n", str1, str2, str1 < str2)
}
```

## Logical Operators

Logical operators perform logical operations on boolean values.

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `&&` | Logical AND | `true && false` | `false` |
| `\|\|` | Logical OR | `true \|\| false` | `true` |
| `!` | Logical NOT | `!true` | `false` |

### Truth Tables

#### AND (&&)
| A | B | A && B |
|---|---|--------|
| true | true | true |
| true | false | false |
| false | true | false |
| false | false | false |

#### OR (||)
| A | B | A \|\| B |
|---|---|--------|
| true | true | true |
| true | false | true |
| false | true | true |
| false | false | false |

#### NOT (!)
| A | !A |
|---|----|
| true | false |
| false | true |

### Examples

```go
package main

import "fmt"

func main() {
    a := true
    b := false
    
    fmt.Printf("a = %t, b = %t\n", a, b)
    fmt.Printf("a && b: %t\n", a && b)
    fmt.Printf("a || b: %t\n", a || b)
    fmt.Printf("!a: %t\n", !a)
    fmt.Printf("!b: %t\n", !b)
    
    // Practical example
    age := 25
    hasLicense := true
    
    canDrive := age >= 18 && hasLicense
    fmt.Printf("\nCan drive (age >= 18 && hasLicense): %t\n", canDrive)
    
    // Short-circuit evaluation
    fmt.Printf("\nShort-circuit examples:\n")
    fmt.Printf("false && (expensive operation): %t\n", false && expensiveOperation())
    fmt.Printf("true || (expensive operation): %t\n", true || expensiveOperation())
}

func expensiveOperation() bool {
    fmt.Println("  Expensive operation called!")
    return true
}
```

## Bitwise Operators

Bitwise operators perform operations on individual bits.

| Operator | Description | Example | Result |
|----------|-------------|---------|--------|
| `&` | Bitwise AND | `5 & 3` | `1` |
| `\|` | Bitwise OR | `5 \| 3` | `7` |
| `^` | Bitwise XOR | `5 ^ 3` | `6` |
| `&^` | Bit clear (AND NOT) | `5 &^ 3` | `4` |
| `<<` | Left shift | `5 << 1` | `10` |
| `>>` | Right shift | `5 >> 1` | `2` |

### Examples

```go
package main

import "fmt"

func main() {
    a := 5  // Binary: 101
    b := 3  // Binary: 011
    
    fmt.Printf("a = %d (binary: %b)\n", a, a)
    fmt.Printf("b = %d (binary: %b)\n", b, b)
    
    fmt.Printf("a & b = %d (binary: %b)\n", a&b, a&b)
    fmt.Printf("a | b = %d (binary: %b)\n", a|b, a|b)
    fmt.Printf("a ^ b = %d (binary: %b)\n", a^b, a^b)
    fmt.Printf("a &^ b = %d (binary: %b)\n", a&^b, a&^b)
    
    // Bit shifting
    fmt.Printf("\nBit shifting:\n")
    fmt.Printf("a << 1 = %d (binary: %b)\n", a<<1, a<<1)
    fmt.Printf("a >> 1 = %d (binary: %b)\n", a>>1, a>>1)
    
    // Practical use: Powers of 2
    fmt.Printf("\nPowers of 2 using bit shifting:\n")
    for i := 0; i < 5; i++ {
        power := 1 << i
        fmt.Printf("2^%d = %d\n", i, power)
    }
}
```

## Increment and Decrement Operators

Go provides increment (`++`) and decrement (`--`) operators.

```go
package main

import "fmt"

func main() {
    x := 5
    fmt.Printf("Initial value: x = %d\n", x)
    
    x++  // Increment by 1
    fmt.Printf("After x++: x = %d\n", x)
    
    x--  // Decrement by 1
    fmt.Printf("After x--: x = %d\n", x)
    
    // Note: Go only supports postfix increment/decrement
    // ++x and --x are not valid in Go
    
    // Also, increment/decrement are statements, not expressions
    // y := x++  // This is invalid
    
    // Correct way:
    y := x
    x++
    fmt.Printf("y = %d, x = %d\n", y, x)
}
```

## Operator Precedence

Operator precedence determines the order of operations in expressions.

| Precedence | Operators | Associativity |
|------------|-----------|---------------|
| 5 (highest) | `*`, `/`, `%`, `<<`, `>>`, `&`, `&^` | Left to right |
| 4 | `+`, `-`, `\|`, `^` | Left to right |
| 3 | `==`, `!=`, `<`, `<=`, `>`, `>=` | Left to right |
| 2 | `&&` | Left to right |
| 1 (lowest) | `\|\|` | Left to right |

### Examples

```go
package main

import "fmt"

func main() {
    // Precedence examples
    result1 := 2 + 3 * 4      // 2 + (3 * 4) = 14
    result2 := (2 + 3) * 4    // (2 + 3) * 4 = 20
    result3 := 10 / 2 * 3     // (10 / 2) * 3 = 15
    result4 := 10 / (2 * 3)   // 10 / (2 * 3) = 1
    
    fmt.Printf("2 + 3 * 4 = %d\n", result1)
    fmt.Printf("(2 + 3) * 4 = %d\n", result2)
    fmt.Printf("10 / 2 * 3 = %d\n", result3)
    fmt.Printf("10 / (2 * 3) = %d\n", result4)
    
    // Logical operator precedence
    a, b, c := true, false, true
    result5 := a || b && c    // a || (b && c) = true
    result6 := (a || b) && c  // (a || b) && c = true
    
    fmt.Printf("true || false && true = %t\n", result5)
    fmt.Printf("(true || false) && true = %t\n", result6)
}
```

## Practical Examples

### Example 1: Calculator

```go
package main

import "fmt"

func main() {
    var num1, num2 float64
    var operator string
    
    fmt.Print("Enter first number: ")
    fmt.Scanln(&num1)
    
    fmt.Print("Enter operator (+, -, *, /): ")
    fmt.Scanln(&operator)
    
    fmt.Print("Enter second number: ")
    fmt.Scanln(&num2)
    
    var result float64
    var valid bool = true
    
    switch operator {
    case "+":
        result = num1 + num2
    case "-":
        result = num1 - num2
    case "*":
        result = num1 * num2
    case "/":
        if num2 != 0 {
            result = num1 / num2
        } else {
            fmt.Println("Error: Division by zero!")
            valid = false
        }
    default:
        fmt.Println("Error: Invalid operator!")
        valid = false
    }
    
    if valid {
        fmt.Printf("%.2f %s %.2f = %.2f\n", num1, operator, num2, result)
    }
}
```

### Example 2: Grade Calculator

```go
package main

import "fmt"

func main() {
    var score int
    fmt.Print("Enter your score (0-100): ")
    fmt.Scanln(&score)
    
    // Validate input
    if score < 0 || score > 100 {
        fmt.Println("Invalid score! Please enter a value between 0 and 100.")
        return
    }
    
    // Determine grade using comparison operators
    var grade string
    var passed bool
    
    if score >= 90 {
        grade = "A"
    } else if score >= 80 {
        grade = "B"
    } else if score >= 70 {
        grade = "C"
    } else if score >= 60 {
        grade = "D"
    } else {
        grade = "F"
    }
    
    passed = score >= 60
    
    fmt.Printf("Score: %d\n", score)
    fmt.Printf("Grade: %s\n", grade)
    fmt.Printf("Passed: %t\n", passed)
    
    // Additional feedback
    if passed && score >= 90 {
        fmt.Println("Excellent work!")
    } else if passed {
        fmt.Println("Good job!")
    } else {
        fmt.Println("Better luck next time!")
    }
}
```

### Example 3: Bitwise Flags

```go
package main

import "fmt"

const (
    ReadFlag    = 1 << iota // 1 (binary: 001)
    WriteFlag               // 2 (binary: 010)
    ExecuteFlag             // 4 (binary: 100)
)

func main() {
    // Set permissions
    var permissions int
    
    // Grant read and write permissions
    permissions |= ReadFlag
    permissions |= WriteFlag
    
    fmt.Printf("Current permissions: %d (binary: %b)\n", permissions, permissions)
    
    // Check permissions
    canRead := (permissions & ReadFlag) != 0
    canWrite := (permissions & WriteFlag) != 0
    canExecute := (permissions & ExecuteFlag) != 0
    
    fmt.Printf("Can read: %t\n", canRead)
    fmt.Printf("Can write: %t\n", canWrite)
    fmt.Printf("Can execute: %t\n", canExecute)
    
    // Grant execute permission
    permissions |= ExecuteFlag
    fmt.Printf("\nAfter granting execute permission: %d (binary: %b)\n", permissions, permissions)
    
    // Remove write permission
    permissions &^= WriteFlag
    fmt.Printf("After removing write permission: %d (binary: %b)\n", permissions, permissions)
    
    // Check if has all permissions
    allPermissions := ReadFlag | WriteFlag | ExecuteFlag
    hasAllPermissions := (permissions & allPermissions) == allPermissions
    fmt.Printf("Has all permissions: %t\n", hasAllPermissions)
}
```

## Exercises

### Exercise 1: BMI Calculator
Create a BMI calculator using arithmetic and comparison operators.

```go
package main

import "fmt"

func main() {
    // Get weight (kg) and height (m) from user
    // Calculate BMI = weight / (height * height)
    // Determine BMI category:
    // - Underweight: BMI < 18.5
    // - Normal: 18.5 <= BMI < 25
    // - Overweight: 25 <= BMI < 30
    // - Obese: BMI >= 30
}
```

### Exercise 2: Number Properties
Check various properties of a number using different operators.

```go
package main

import "fmt"

func main() {
    var num int
    fmt.Print("Enter a number: ")
    fmt.Scanln(&num)
    
    // Check and display:
    // - Is even or odd
    // - Is positive, negative, or zero
    // - Is divisible by 3
    // - Is divisible by 5
    // - Is a power of 2
}
```

### Exercise 3: Logical Operations
Create a program that demonstrates logical operators with real-world scenarios.

```go
package main

import "fmt"

func main() {
    // Define variables for:
    // - age
    // - hasDriversLicense
    // - hasInsurance
    // - hasValidRegistration
    
    // Determine:
    // - Can drive (age >= 16 && hasDriversLicense)
    // - Can drive legally (canDrive && hasInsurance && hasValidRegistration)
    // - Needs adult supervision (age < 18)
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
    var weight, height float64
    
    fmt.Print("Enter your weight (kg): ")
    fmt.Scanln(&weight)
    
    fmt.Print("Enter your height (m): ")
    fmt.Scanln(&height)
    
    bmi := weight / (height * height)
    
    fmt.Printf("Your BMI: %.2f\n", bmi)
    
    var category string
    if bmi < 18.5 {
        category = "Underweight"
    } else if bmi < 25 {
        category = "Normal"
    } else if bmi < 30 {
        category = "Overweight"
    } else {
        category = "Obese"
    }
    
    fmt.Printf("Category: %s\n", category)
}
```

### Solution 2:
```go
package main

import "fmt"

func main() {
    var num int
    fmt.Print("Enter a number: ")
    fmt.Scanln(&num)
    
    fmt.Printf("Number: %d\n", num)
    
    // Even or odd
    if num%2 == 0 {
        fmt.Println("Even")
    } else {
        fmt.Println("Odd")
    }
    
    // Positive, negative, or zero
    if num > 0 {
        fmt.Println("Positive")
    } else if num < 0 {
        fmt.Println("Negative")
    } else {
        fmt.Println("Zero")
    }
    
    // Divisible by 3
    fmt.Printf("Divisible by 3: %t\n", num%3 == 0)
    
    // Divisible by 5
    fmt.Printf("Divisible by 5: %t\n", num%5 == 0)
    
    // Power of 2 (positive numbers only)
    isPowerOf2 := num > 0 && (num&(num-1)) == 0
    fmt.Printf("Power of 2: %t\n", isPowerOf2)
}
```

### Solution 3:
```go
package main

import "fmt"

func main() {
    age := 17
    hasDriversLicense := true
    hasInsurance := false
    hasValidRegistration := true
    
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Has driver's license: %t\n", hasDriversLicense)
    fmt.Printf("Has insurance: %t\n", hasInsurance)
    fmt.Printf("Has valid registration: %t\n", hasValidRegistration)
    
    canDrive := age >= 16 && hasDriversLicense
    canDriveLegally := canDrive && hasInsurance && hasValidRegistration
    needsSupervision := age < 18
    
    fmt.Printf("\nCan drive: %t\n", canDrive)
    fmt.Printf("Can drive legally: %t\n", canDriveLegally)
    fmt.Printf("Needs adult supervision: %t\n", needsSupervision)
    
    if !canDriveLegally {
        fmt.Println("\nReasons cannot drive legally:")
        if !canDrive {
            if age < 16 {
                fmt.Println("- Too young to drive")
            }
            if !hasDriversLicense {
                fmt.Println("- No driver's license")
            }
        }
        if !hasInsurance {
            fmt.Println("- No insurance")
        }
        if !hasValidRegistration {
            fmt.Println("- No valid registration")
        }
    }
}
```

</details>

## Key Takeaways

1. **Arithmetic operators** perform mathematical calculations
2. **Assignment operators** provide shortcuts for common operations
3. **Comparison operators** return boolean values
4. **Logical operators** work with boolean values and support short-circuit evaluation
5. **Bitwise operators** manipulate individual bits
6. **Operator precedence** determines evaluation order
7. Use parentheses to make precedence explicit
8. Go only supports postfix increment/decrement
9. Division behavior depends on operand types (integer vs floating-point)

## Next Steps

Now that you understand operators, let's learn about [Control Structures](07-control-structures.md) to control program flow!

---

**Previous**: [← Constants](05-constants.md) | **Next**: [Control Structures →](07-control-structures.md)