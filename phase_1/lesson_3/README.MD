# Phase 1 — Go Fundamentals

## Lesson 3: Operators and Expressions in Go

In this lesson, we'll learn how to perform calculations, compare values, and make logical decisions using Go operators.

## 1. Arithmetic Operators

Arithmetic operators are used to perform mathematical calculations.

```
package main

import "fmt"

func main() {
    a := 20
    b := 6

    fmt.Println(a + b) // Addition:       26
    fmt.Println(a - b) // Subtraction:   14
    fmt.Println(a * b) // Multiplication: 120
    fmt.Println(a / b) // Division:       3
    fmt.Println(a % b) // Remainder:     2
}
```

| Operator | Meaning            | Example       |
| -------- | ------------------ | ------------- |
| `+`      | Addition           | `10 + 5 = 15` |
| `-`      | Subtraction        | `10 - 5 = 5`  |
| `*`      | Multiplication     | `10 * 5 = 50` |
| `/`      | Division           | `10 / 5 = 2`  |
| `%`      | Remainder (modulo) | `10 % 3 = 1`  |

Important: Integer division discards the fractional part.

```
fmt.Println(10 / 3)    // 3
fmt.Println(10.0 / 3.0) // 3.3333333333333335
```

## 2. Comparison Operators

Comparison operators compare two values and return a `bool` (`true` or `false`).

```
package main

import "fmt"

func main() {
    a := 10
    b := 20

    fmt.Println(a == b) // false
    fmt.Println(a != b) // true
    fmt.Println(a > b)  // false
    fmt.Println(a < b)  // true
    fmt.Println(a >= b) // false
    fmt.Println(a <= b) // true
}
```

| Operator | Meaning               |
| -------- | --------------------- |
| `==`     | Equal to              |
| `!=`     | Not equal to          |
| `>`      | Greater than          |
| `<`      | Less than             |
| `>=`     | Greater than or equal |
| `<=`     | Less than or equal    |

These operators are very useful for validation, authentication, and business logic in backend development.

## 3. Logical Operators

Logical operators combine multiple Boolean conditions.

```
package main

import "fmt"

func main() {
    age := 25
    hasID := true

    fmt.Println(age >= 18 && hasID) // true
    fmt.Println(age < 18 || hasID) // true
    fmt.Println(!hasID)            // false
}
```

\| Operator | Meaning | Example |\
\|---|---|---|\
\| `&&` | AND — both conditions must be true | `a && b` |\
\| `||` | OR — at least one condition is true | `a || b` |\
\| `!` | NOT — reverses Boolean value | `!a` |

### Truth table

\| A | B | `A && B` | `A || B` |\
\|---|---|---|---|\
\| true | true | true | true |\
\| true | false | false | true |\
\| false | true | false | true |\
\| false | false | false | false |

## 4. Assignment Operators

Assignment operators assign or update variable values.

```
package main

import "fmt"

func main() {
    x := 10

    x += 5
    fmt.Println(x) // 15

    x -= 3
    fmt.Println(x) // 12

    x *= 2
    fmt.Println(x) // 24

    x /= 4
    fmt.Println(x) // 6

    x %= 4
    fmt.Println(x) // 2
}
```

| Operator | Equivalent  |
| -------- | ----------- |
| `x += 5` | `x = x + 5` |
| `x -= 5` | `x = x - 5` |
| `x *= 5` | `x = x * 5` |
| `x /= 5` | `x = x / 5` |
| `x %= 5` | `x = x % 5` |

## 5. Increment and Decrement

Go supports increment (`++`) and decrement (`--`).

```
package main

import "fmt"

func main() {
    count := 10

    count++
    fmt.Println(count) // 11

    count--
    fmt.Println(count) // 10
}
```

Important Go rule: Unlike C, C++, and JavaScript, Go does not support `++count` or using `count++` inside an expression.

```
count++       // Valid
// ++count    // Invalid
// result := count++ // Invalid
```

## 6. Operator Precedence

When an expression contains multiple operators, Go follows precedence rules.

```
package main

import "fmt"

func main() {
    result := 10 + 5*2
    fmt.Println(result) // 20

    result = (10 + 5) * 2
    fmt.Println(result) // 30
}
```

Multiplication happens before addition. Parentheses can change the order.

\| Priority | Operators |\
\|---|---|\
\| Highest | `*` `/` `%` `<<` `>>` `&` `&^` |\
\| Medium | `+` `-` `|` `^` |\
\| Lower | `==` `!=` `<` `<=` `>` `>=` |\
\| Lower | `&&` |\
\| Lowest | `||` |

## 7. Practical backend example: Invoice discount

Let's apply operators to a simple invoicing calculation.

```

package main

import "fmt"

func main() {
    productPrice := 5000.0
    quantity := 3
    discount := 10.0

    subtotal := productPrice * float64(quantity)
    discountAmount := subtotal * discount / 100
    total := subtotal - discountAmount

    fmt.Println("Subtotal:", subtotal, "TK")
    fmt.Println("Discount:", discountAmount, "TK")
    fmt.Println("Total:", total, "TK")

    isHighValue := total >= 10000
    fmt.Println("High-value invoice:", isHighValue)
}

```

Output:

```
Subtotal: 15000 TK
Discount: 1500 TK
Total: 13500 TK
High-value invoice: true
```

Notice `float64(quantity)`: Go requires explicit conversion when multiplying a `float64` by an `int`.

## 8. Practice exercises

Try solving these without looking at the examples.

1. Write a program that takes two variables and prints all five arithmetic results.
2. Check whether a number is even or odd using `%`.
3. Check whether a person is eligible to vote (age >= 18).
4. Check whether a user is eligible for a discount: age >= 60 OR total purchase >= 10,000.
5. Calculate a product's price after a 15% discount.
6. Predict the output:

```
x := 10
x += 5
x *= 2
x -= 4

fmt.Println(x)
```

## Next lesson — Lesson 4: Input and Output in Go

We'll learn:

- `fmt.Print()`, `fmt.Println()`, `fmt.Printf()`
- Taking user input with `fmt.Scan()`
- Formatting strings and numbers
- Building interactive command-line programs
