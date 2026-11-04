# Phase 1 — Go Fundamentals

## Lesson 2: Variables, Constants & Data Types

In this lesson, you'll learn how to store and manage data in Go. These concepts are essential for building backend APIs and applications.

## 1. What is a variable?

A variable is a named storage location used to hold data.

In Go, we can declare variables using the `var` keyword.

```
package main

import "fmt"

func main() {
    var name string = "Jihad"
    var age int = 30

    fmt.Println(name)
    fmt.Println(age)
}
```

Output:

```
Jihad
30
```

### Syntax

```
var variableName dataType = value
```

Examples:

```
var username string = "jihad"
var age int = 30
var price float64 = 99.99
var isActive bool = true
```

## 2. Short variable declaration `:=`

Go provides a shorter way to declare variables inside functions.

```
package main

import "fmt"

func main() {
    name := "Jihad"
    age := 30
    salary := 45000.50
    isDeveloper := true

    fmt.Println(name, age, salary, isDeveloper)
}
```

Important difference:

| Declaration                 | Usage                |
| --------------------------- | -------------------- |
| `var name string = "Jihad"` | Explicit declaration |
| `var name = "Jihad"`        | Type inference       |
| `name := "Jihad"`           | Short declaration    |

Remember: `:=` is used inside functions. You cannot use it for package-level variable declarations.

## 3. Basic data types in Go

[Data Types in Go An Introduction (Golang)](https://images.openai.com/static-rsc-4/h1d2KjyOYp_8CJ706K1XNgoXTcjodCIeMoJ9vRYVjMiv0T_Iq6VjpJIK6ZwBRqGJgyf1dtWlsFi8DFA1NGBnBfuKxSQhzxolmW9QLtFIgwh6JqdEuZNX4vf3QeMze5V2ReXTAO-DShilUKScNOH_DrYJ_nrWvj6Ig8W6GkWh5Xo?purpose=inline)

[Supported Go types for gomobile bind](https://images.openai.com/static-rsc-4/WW96RSSqYY4YEqUA85Nfcar3hyO3S24Q3Tk9ThFWBGa1s6Ni6uWjgWwIC3EzCPn-ozOmMrez1i4bUivL3Df1-yY005SUHy2A6_VRicR_yu5QOaFky-cGE2YJnGpf6laqWTl0iF_b-bLQN4W4lgmbZLniWA0TO4LVWH-druIPg7k?purpose=inline)

[Zero Values - KodeKloud](https://images.openai.com/static-rsc-4/egqPjTbRx3kLcqlJEIzlHqovhUqdo9JS4JxmG5R39aql5hwfCagiXoI4L85-8NSKovPXS9P68qRK3b5CJNV4NVA9cciD9V0E5mmf0y7MiigA-yEg8SjX2a4YAJK1gLy11iNYgWY4trtyOvp4thep4MrPweofR9lIb2ZWYf4KU9Y?purpose=inline)

6

| Data type | Example          | Description       |
| --------- | ---------------- | ----------------- |
| `int`     | `age := 30`      | Integer numbers   |
| `float64` | `price := 99.99` | Decimal numbers   |
| `string`  | `name := "Go"`   | Text              |
| `bool`    | `active := true` | `true` or `false` |

Example:

```
package main

import "fmt"

func main() {
    var age int = 25
    var height float64 = 5.9
    var name string = "Jihad"
    var isMarried bool = false

    fmt.Println(age)
    fmt.Println(height)
    fmt.Println(name)
    fmt.Println(isMarried)
}
```

## 4. Type inference

Go can automatically determine a variable's type from its assigned value.

```
name := "Jihad"       // string
age := 30             // int
price := 250.50        // float64
isActive := true       // bool
```

You don't always need to specify the data type explicitly.

## 5. Constants in Go

A constant is a value that cannot be changed after declaration.

Use the `const` keyword.

```
package main

import "fmt"

func main() {
    const appName = "Olo Invoicer"
    const taxRate = 15.0

    fmt.Println(appName)
    fmt.Println(taxRate)
}
```

You cannot reassign a constant:

```
const pi = 3.14159

// pi = 3.14  // Compilation error
```

### Variable vs Constant

| Variable (`var`)    | Constant (`const`)      |
| ------------------- | ----------------------- |
| Can be reassigned   | Cannot be reassigned    |
| Stores mutable data | Represents fixed values |
| `var age = 30`      | `const taxRate = 15`    |

## 6. Zero values in Go

If you declare a variable without assigning a value, Go automatically gives it a zero value.

```
package main

import "fmt"

func main() {
    var age int
    var price float64
    var name string
    var isActive bool

    fmt.Println(age)
    fmt.Println(price)
    fmt.Println(name)
    fmt.Println(isActive)
}
```

Output:

```
0
0

false
```

| Type      | Zero value          |
| --------- | ------------------- |
| `int`     | `0`                 |
| `float64` | `0.0`               |
| `string`  | `""` (empty string) |
| `bool`    | `false`             |

This is an important Go feature: variables are initialized automatically.

## 7. Practical example: Developer profile

```

package main

import "fmt"

func main() {
    name := "Jihad"
    profession := "Full Stack Developer"
    experience := 3
    salary := 50000.0
    isEmployed := true

    fmt.Println("Name:", name)
    fmt.Println("Profession:", profession)
    fmt.Println("Experience:", experience, "years")
    fmt.Println("Salary:", salary, "TK")
    fmt.Println("Employed:", isEmployed)
}

```

## 8. Practice exercises

Try writing these programs yourself.

1. Declare variables for your name, age, city, and profession. Print them.
2. Declare two integers and print their sum, subtraction, multiplication, and division.
3. Declare a constant `taxRate = 15` and print it.
4. Declare variables without values and print their zero values.
5. Create a product with `name`, `price`, and `isAvailable` variables, then print its information.

## Next lesson — Lesson 3

Operators and Expressions in Go

- Arithmetic operators
- Comparison operators
- Logical operators
- Assignment operators
- Operator precedence
