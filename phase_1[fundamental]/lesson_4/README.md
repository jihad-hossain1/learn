# Phase 1 — Go Fundamentals

## Lesson 4: Input and Output in Go

In the previous lessons, we printed fixed values. Now we'll learn how to take input from users and display formatted output.

This is an important step toward building interactive CLI applications and backend programs.

## 1. Output using `fmt`

Go provides the `fmt` package for formatted input and output.

### `fmt.Print()`

Prints text without automatically adding a newline.

```
fmt.Print("Hello ")
fmt.Print("Go!")
```

Output:

```
Hello Go!
```

### `fmt.Println()`

Prints text and automatically adds a newline.

```
fmt.Println("Hello")
fmt.Println("World")
```

Output:

```
Hello
World
```

### `fmt.Printf()`

Prints formatted output using format specifiers.

```
name := "Jihad"
age := 30

fmt.Printf("My name is %s and I am %d years old.\n", name, age)
```

Output:

```
My name is Jihad and I am 30 years old.
```

## 2. Format specifiers

| Specifier | Data type                   | Example         |
| --------- | --------------------------- | --------------- |
| `%s`      | String                      | `"Jihad"`       |
| `%d`      | Integer                     | `30`            |
| `%f`      | Float                       | `99.99`         |
| `%t`      | Boolean                     | `true`          |
| `%v`      | Any value (default format)  | `30`, `"Go"`    |
| `%T`      | Data type                   | `int`, `string` |
| `%.2f`    | Float with 2 decimal places | `99.99`         |
| `%%`      | Literal percent sign        | `%`             |

Example:

```
package main

import "fmt"

func main() {
    name := "Jihad"
    age := 30
    salary := 45000.567
    isDeveloper := true

    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Salary: %.2f TK\n", salary)
    fmt.Printf("Developer: %t\n", isDeveloper)

    fmt.Printf("Type of name: %T\n", name)
}
```

## 3. Taking user input with `fmt.Scan()`

The `fmt.Scan()` function reads input from the keyboard.

```
package main

import "fmt"

func main() {
    var name string

    fmt.Print("Enter your name: ")
    fmt.Scan(&name)

    fmt.Println("Hello,", name)
}
```

Example output:

```
Enter your name: Jihad
Hello, Jihad
```

### What does `&name` mean?

The `&` operator gives the memory address of a variable.

`fmt.Scan(&name)` passes the address so that `Scan()` can store the user's input inside the variable.

We'll explore pointers in a later lesson.

## 4. Taking multiple inputs

```
package main

import "fmt"

func main() {
    var name string
    var age int
    var salary float64

    fmt.Print("Enter your name: ")
    fmt.Scan(&name)

    fmt.Print("Enter your age: ")
    fmt.Scan(&age)

    fmt.Print("Enter your salary: ")
    fmt.Scan(&salary)

    fmt.Println("\n--- User Information ---")
    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d\n", age)
    fmt.Printf("Salary: %.2f TK\n", salary)
}
```

You can enter a string, integer, and decimal number separately.

Note: `fmt.Scan()` reads space-separated values. For example, if you enter `Jihad Hossain`, it reads only `Jihad` into a single string variable.

## 5. Reading a full line with `bufio`

To read a full name containing spaces, use `bufio` with `os.Stdin`.

```
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Enter your full name: ")

    name, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading input:", err)
        return
    }

    name = strings.TrimSpace(name)

    fmt.Println("Hello,", name)
}
```

Input:

```
Jihad Hossain
```

Output:

```
Hello, Jihad Hossain
```

### New concepts

| Code                  | Meaning                                    |
| --------------------- | ------------------------------------------ |
| `bufio.NewReader()`   | Creates a buffered reader                  |
| `os.Stdin`            | Standard input (keyboard)                  |
| `ReadString('\n')`    | Reads until newline                        |
| `err`                 | Stores information about an error          |
| `strings.TrimSpace()` | Removes surrounding whitespace and newline |

We'll study error handling in detail in Phase 3.

## 6. Practical project: Simple calculator

Let's combine variables, operators, and user input.

```

package main

import "fmt"

func main() {
    var a, b float64
    var operator string

    fmt.Print("Enter first number: ")
    fmt.Scan(&a)

    fmt.Print("Enter operator (+, -, *, /): ")
    fmt.Scan(&operator)

    fmt.Print("Enter second number: ")
    fmt.Scan(&b)

    switch operator {
    case "+":
        fmt.Printf("Result: %.2f\n", a+b)
    case "-":
        fmt.Printf("Result: %.2f\n", a-b)
    case "*":
        fmt.Printf("Result: %.2f\n", a*b)
    case "/":
        if b == 0 {
            fmt.Println("Error: Division by zero!")
        } else {
            fmt.Printf("Result: %.2f\n", a/b)
        }
    default:
        fmt.Println("Invalid operator!")
    }
}

```

Example output:

```
Enter first number: 20
Enter operator (+, -, *, /): *
Enter second number: 5
Result: 100.00
```

Notice that we introduced `switch` and `if` statements. We'll study them deeply in the next lesson.

## 7. Practice exercises

Try these yourself:

1. Take a user's name and age, then print a formatted greeting.
2. Take two numbers and display their sum.
3. Take a product name, price, and quantity, then calculate the total price.
4. Build a calculator supporting `+`, `-`, `*`, `/`.
5. Take a user's full name using `bufio.Reader`.
6. Print a float value with 2, 3, and 4 decimal places using `fmt.Printf()`.
