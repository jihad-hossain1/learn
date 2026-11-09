# Phase 1 — Go Fundamentals

## Lesson 6: Functions in Go

Functions are one of the most important concepts in programming. They allow you to organize code into reusable, independent blocks.

As we move toward backend development with Go, you'll use functions to build:

- Business logic
- Authentication services
- Database operations
- API handlers
- Validation and utility functions

## 1. What is a function?

A function is a reusable block of code designed to perform a specific task.

Instead of writing the same code repeatedly, you can define a function once and call it whenever needed.

For example, without a function:

```
package main

import "fmt"

func main() {
    fmt.Println("Hello, Jihad!")
    fmt.Println("Hello, Rahim!")
    fmt.Println("Hello, Karim!")
}
```

Using a function:

```
package main

import "fmt"

func sayHello() {
    fmt.Println("Hello!")
}

func main() {
    sayHello()
    sayHello()
    sayHello()
}
```

Output:

```
Hello!
Hello!
Hello!
```

## 2. Function syntax

The general syntax of a Go function is:

```
func functionName(parameters) returnType {
    // Function body
    return value
}
```

Let's understand each part:

| Component      | Meaning                            |
| -------------- | ---------------------------------- |
| `func`         | Keyword used to declare a function |
| `functionName` | Name of the function               |
| `parameters`   | Input values                       |
| `returnType`   | Type of returned value             |
| `return`       | Sends a result back to the caller  |

## 3. Functions without parameters and return values

A function doesn't always need input or output.

```
package main

import "fmt"

func greet() {
    fmt.Println("Welcome to Go programming!")
}

func main() {
    greet()
}
```

Here:

- `greet()` is the function declaration.
- `main()` calls the `greet()` function.
- The function performs its task without receiving or returning anything.

## 4. Functions with parameters

Parameters allow you to pass data into a function.

```
package main

import "fmt"

func greet(name string) {
    fmt.Println("Hello,", name)
}

func main() {
    greet("Jihad")
    greet("Rahim")
    greet("Karim")
}
```

Output:

```
Hello, Jihad
Hello, Rahim
Hello, Karim
```

### Multiple parameters

You can pass multiple parameters to a function.

```
package main

import "fmt"

func displayUser(name string, age int) {
    fmt.Println("Name:", name)
    fmt.Println("Age:", age)
}

func main() {
    displayUser("Jihad", 30)
}
```

Go syntax shortcut:

If consecutive parameters have the same type, you can declare them together.

```
func add(a, b int) {
    fmt.Println(a + b)
}
```

This is equivalent to:

```
func add(a int, b int) {
    fmt.Println(a + b)
}
```

## 5. Functions with return values

A function can return a value to its caller.

```
package main

import "fmt"

func add(a int, b int) int {
    result := a + b
    return result
}

func main() {
    total := add(10, 20)

    fmt.Println("Total:", total)
}
```

Output:

```
Total: 30
```

Notice:

```
func add(a int, b int) int
```

The final `int` specifies that this function returns an integer.

Important: The return type comes after the parameter list.

## 6. Multiple return values

One of Go's useful features is that a function can return multiple values.

```
package main

import "fmt"

func calculate(a, b int) (int, int, int, int) {
    sum := a + b
    sub := a - b
    mul := a * b
    div := a / b

    return sum, sub, mul, div
}

func main() {
    sum, sub, mul, div := calculate(20, 5)

    fmt.Println("Addition:", sum)
    fmt.Println("Subtraction:", sub)
    fmt.Println("Multiplication:", mul)
    fmt.Println("Division:", div)
}
```

Output:

```
Addition: 25
Subtraction: 15
Multiplication: 100
Division: 4
```

You can also ignore returned values using `_`:

```
sum, _, _, _ := calculate(20, 5)
fmt.Println(sum)
```

The underscore is called the blank identifier.

## 7. Named return values

Go allows you to name the return values in the function signature.

```
package main

import "fmt"

func calculate(a, b int) (sum int, difference int) {
    sum = a + b
    difference = a - b

    return
}

func main() {
    total, diff := calculate(20, 10)

    fmt.Println("Sum:", total)
    fmt.Println("Difference:", diff)
}
```

Output:

```
Sum: 30
Difference: 10
```

Here:

- `sum` and `difference` are named return variables.
- They are automatically initialized to their zero values.
- A bare `return` returns their current values.

Best practice: Named returns are useful when they make the returned values clearer. Avoid using them unnecessarily in long functions.

## 8. Variadic functions

A variadic function accepts a variable number of arguments.

For example, you might want to calculate the total of 2, 3, 5, or more numbers without creating separate functions.

```
package main

import "fmt"

func sum(numbers ...int) int {
    total := 0

    for _, number := range numbers {
        total += number
    }

    return total
}

func main() {
    fmt.Println(sum(10, 20))
    fmt.Println(sum(10, 20, 30))
    fmt.Println(sum(10, 20, 30, 40, 50))
}
```

Output:

```
30
60
150
```

### Understanding `...int`

```
func sum(numbers ...int) int
```

This means:

- Accept zero or more integer arguments.
- Inside the function, `numbers` behaves like a slice of integers.
- You can iterate over it using `range`.

You can also pass an existing slice using `...`:

```
values := []int{10, 20, 30}

result := sum(values...)
fmt.Println(result)
```

Output:

```
60
```

## 9. Function scope in Go

Scope determines where a variable or function can be accessed.

Go has different scopes, including:

A. Local scope

Variables declared inside a function are accessible only within that function.

```
func greet() {
    message := "Hello Go!"
    fmt.Println(message)
}
```

You cannot access `message` directly from another function.

B. Package-level scope

Variables declared outside functions can be accessed by other functions in the same package.

```
package main

import "fmt"

var appName = "Go Learning"

func showApp() {
    fmt.Println(appName)
}

func main() {
    showApp()
    fmt.Println(appName)
}
```

C. Block scope

Variables declared inside an `if`, `for`, or other block are generally limited to that block.

```
package main

import "fmt"

func main() {
    age := 25

    if age >= 18 {
        message := "Adult"
        fmt.Println(message)
    }

    // fmt.Println(message) // Compilation error
}
```

### Important: Exported and unexported names

Go uses capitalization to determine whether an identifier is exported from a package.

```
func CalculateTotal() {} // Exported

func calculateTax() {}   // Unexported
```

- Uppercase first letter: accessible from other packages.
- Lowercase first letter: accessible only within the declaring package.

This is particularly important when designing reusable Go packages.

## 10. Functions as values

In Go, functions are first-class values. This means you can assign a function to a variable and call it through that variable.

```
package main

import "fmt"

func add(a, b int) int {
    return a + b
}

func main() {
    operation := add

    result := operation(10, 20)

    fmt.Println(result)
}
```

Output:

```
30
```

You can also pass functions as arguments to other functions.

```
package main

import "fmt"

func calculate(a, b int, operation func(int, int) int) int {
    return operation(a, b)
}

func add(a, b int) int {
    return a + b
}

func multiply(a, b int) int {
    return a * b
}

func main() {
    fmt.Println(calculate(10, 5, add))
    fmt.Println(calculate(10, 5, multiply))
}
```

Output:

```
15
50
```

This concept is useful for callbacks, reusable business logic, and middleware.

## 11. Anonymous functions

An anonymous function is a function without a declared name.

```
package main

import "fmt"

func main() {
    greet := func(name string) {
        fmt.Println("Hello,", name)
    }

    greet("Jihad")
}
```

You can also execute an anonymous function immediately:

```
func() {
    fmt.Println("Anonymous function executed!")
}()
```

The final `()` calls the function immediately after defining it.

## 12. Practical project: Invoice calculation system

Now let's combine what we've learned into a small backend-style example.

```

package main

import "fmt"

func calculateSubtotal(price float64, quantity int) float64 {
    return price * float64(quantity)
}

func calculateDiscount(subtotal, discountRate float64) float64 {
    return subtotal * discountRate / 100
}

func calculateTotal(subtotal, discount float64) float64 {
    return subtotal - discount
}

func generateInvoice(
    productName string,
    price float64,
    quantity int,
    discountRate float64,
) {
    subtotal := calculateSubtotal(price, quantity)
    discount := calculateDiscount(subtotal, discountRate)
    total := calculateTotal(subtotal, discount)

    fmt.Println("----- INVOICE -----")
    fmt.Println("Product:", productName)
    fmt.Printf("Price: %.2f TK\n", price)
    fmt.Println("Quantity:", quantity)
    fmt.Printf("Subtotal: %.2f TK\n", subtotal)
    fmt.Printf("Discount: %.2f TK\n", discount)
    fmt.Printf("Total: %.2f TK\n", total)
}

func main() {
    generateInvoice("Keyboard", 1500, 2, 10)
}

```

Output:

```
----- INVOICE -----
Product: Keyboard
Price: 1500.00 TK
Quantity: 2
Subtotal: 3000.00 TK
Discount: 300.00 TK
Total: 2700.00 TK
```

Notice how we divided the invoice calculation into smaller functions.

This makes code:

- Easier to maintain.
- Easier to test.
- Reusable across multiple modules.
- Easier to understand.

This is the foundation of clean backend architecture.

## 13. Practice exercises

Now it's your turn. Try these exercises without copying the examples.

Exercise 1 — Basic function

Easy

Create a function called `greet()` that prints your name and profession.

Exercise 2 — Addition

Easy

Create an `add(a, b int) int` function that returns the sum of two integers.

Exercise 3 — Multiple returns

Medium

Create a function that accepts two integers and returns their sum, subtraction, and multiplication.

Exercise 4 — Variadic function

Medium

Create a function that accepts any number of integers and calculates their average. Handle the case where no numbers are provided.

Exercise 5 — Calculator

Challenge

Create separate functions for:

- `add()`
- `subtract()`
- `multiply()`
- `divide()`

Build a calculator that accepts user input and calls the appropriate function. Handle division by zero.

Exercise 6 — Invoice system

Challenge

Create reusable functions to calculate:

- Subtotal
- Discount
- VAT
- Final payable amount

## 14. Important Go function concepts — Quick revision

| Concept           | Syntax                        |
| ----------------- | ----------------------------- |
| Function          | `func greet() {}`             |
| Parameter         | `func greet(name string) {}`  |
| Return value      | `func add(a, b int) int`      |
| Multiple returns  | `func calculate() (int, int)` |
| Named returns     | `func calculate() (sum int)`  |
| Variadic          | `func sum(nums ...int) int`   |
| Anonymous         | `func() {}`                   |
| Function as value | `operation := add`            |

Professional recommendation: In real Go projects, keep functions focused on one responsibility. Prefer small, readable functions rather than one large function handling everything.
