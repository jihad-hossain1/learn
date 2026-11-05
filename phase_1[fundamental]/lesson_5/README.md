# Phase 1 — Go Fundamentals

## Lesson 5: Control Flow in Go

Control flow determines how a program makes decisions and executes instructions.

Today we'll learn:

1. `if`, `else if`, `else`
2. `switch`
3. `for` loops
4. `break` and `continue`
5. Nested conditions and loops

These concepts are essential for backend business logic, authentication, validation, and API development.

## 1. `if` Statement

The `if` statement executes a block of code when a condition is `true`.

```
package main

import "fmt"

func main() {
    age := 25

    if age >= 18 {
        fmt.Println("You are an adult.")
    }
}
```

Output:

```
You are an adult.
```

Syntax:

```
if condition {
    // Execute when condition is true
}
```

Notice that Go does not require parentheses around the condition.

## 2. `if...else`

When the condition is false, `else` executes.

```
package main

import "fmt"

func main() {
    age := 16

    if age >= 18 {
        fmt.Println("Eligible to vote.")
    } else {
        fmt.Println("Not eligible to vote.")
    }
}
```

Output:

```
Not eligible to vote.
```

## 3. `if...else if...else`

Use this when you need to check multiple conditions.

```

package main

import "fmt"

func main() {
    marks := 85

    if marks >= 80 {
        fmt.Println("Grade: A+")
    } else if marks >= 70 {
        fmt.Println("Grade: A")
    } else if marks >= 60 {
        fmt.Println("Grade: A-")
    } else if marks >= 50 {
        fmt.Println("Grade: B")
    } else {
        fmt.Println("Grade: Fail")
    }
}

```

Output:

```
Grade: A+
```

Important: Go evaluates conditions from top to bottom and executes only the first matching branch.

### Practical example: Login validation

```
username := "admin"
password := "123456"

if username == "admin" && password == "123456" {
    fmt.Println("Login successful")
} else {
    fmt.Println("Invalid credentials")
}
```

## 4. `if` with a short statement

Go allows you to declare a variable directly inside an `if` statement.

```
package main

import "fmt"

func main() {
    if age := 20; age >= 18 {
        fmt.Println("Adult")
    } else {
        fmt.Println("Minor")
    }
}
```

Here, `age` exists only within the `if` and `else` blocks.

This is useful when checking the result of a function.

## 5. `switch` Statement

A `switch` statement is a cleaner alternative to multiple `if...else` conditions when comparing one value against several possible cases.

```
package main

import "fmt"

func main() {
    day := 3

    switch day {
    case 1:
        fmt.Println("Sunday")
    case 2:
        fmt.Println("Monday")
    case 3:
        fmt.Println("Tuesday")
    case 4:
        fmt.Println("Wednesday")
    default:
        fmt.Println("Invalid day")
    }
}
```

Output:

```
Tuesday
```

### Multiple values in one case

```
day := "Saturday"

switch day {
case "Friday", "Saturday":
    fmt.Println("Weekend!")
default:
    fmt.Println("Working day")
}
```

### Switch without an expression

Go also supports condition-based switches:

```
age := 25

switch {
case age < 13:
    fmt.Println("Child")
case age < 20:
    fmt.Println("Teenager")
case age < 60:
    fmt.Println("Adult")
default:
    fmt.Println("Senior")
}
```

Important: Go automatically stops after a matching case. Unlike C, Go does not fall through by default.

## 6. `for` Loop

Go has only one looping keyword: `for`.

### A. Traditional `for` loop

```
package main

import "fmt"

func main() {
    for i := 1; i <= 5; i++ {
        fmt.Println(i)
    }
}
```

Output:

```
1
2
3
4
5
```

Syntax:

```
for initialization; condition; post {
    // Loop body
}
```

### B. `for` as a while loop

Go doesn't have a `while` keyword. Instead, use `for` with a condition.

```
i := 1

for i <= 5 {
    fmt.Println(i)
    i++
}
```

### C. Infinite loop

```
for {
    fmt.Println("Running...")
    break
}
```

An infinite loop continues until you explicitly stop it or return from the function.

## 7. `break` and `continue`

### `break`

Terminates the loop immediately.

```
for i := 1; i <= 10; i++ {
    if i == 5 {
        break
    }

    fmt.Println(i)
}
```

Output:

```
1
2
3
4
```

### `continue`

Skips the current iteration and moves to the next one.

```
for i := 1; i <= 5; i++ {
    if i == 3 {
        continue
    }

    fmt.Println(i)
}
```

Output:

```
1
2
4
5
```

## 8. Nested loops

A loop inside another loop is called a nested loop.

```
package main

import "fmt"

func main() {
    for i := 1; i <= 3; i++ {
        for j := 1; j <= 2; j++ {
            fmt.Printf("i=%d, j=%d\n", i, j)
        }
    }
}
```

Output:

```
i=1, j=1
i=1, j=2
i=2, j=1
i=2, j=2
i=3, j=1
i=3, j=2
```

Nested loops are useful for processing multidimensional data, but be careful with performance when working with large datasets.

## 9. Practical project: Invoice grading system

Let's combine input, conditions, and operators.

```

package main

import "fmt"

func main() {
    var total float64

    fmt.Print("Enter invoice total (TK): ")
    fmt.Scan(&total)

    switch {
    case total <= 0:
        fmt.Println("Invalid invoice amount")
    case total < 1000:
        fmt.Println("Small invoice")
    case total < 10000:
        fmt.Println("Medium invoice")
    default:
        fmt.Println("Large invoice")
    }

    for i := 1; i <= 3; i++ {
        fmt.Printf("Invoice processing step %d\n", i)
    }
}

```

## 10. Practice exercises

Try these without looking at the examples.

Exercise 1: Take an integer and check whether it is positive, negative, or zero.

Exercise 2: Take a user's age and categorize them as:

- Child (0–12)
- Teenager (13–19)
- Adult (20–59)
- Senior (60+)

Exercise 3: Create a calculator using `switch` supporting `+`, `-`, `*`, `/`.

Exercise 4: Print numbers from 1 to 100 using a `for` loop.

Exercise 5: Print all even numbers between 1 and 50.

Exercise 6: Print a multiplication table for a number entered by the user.

Example:

```
Enter number: 5

5 x 1 = 5
5 x 2 = 10
...
5 x 10 = 50
```

Exercise 7: Calculate the sum of numbers from 1 to N, where N is entered by the user.
