# 3. Basic Syntax and Hello World

## Go Program Structure

Every Go program consists of the following elements:

1. **Package declaration**
2. **Import statements**
3. **Functions, variables, constants, and types**

### Basic Program Template

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## Package Declaration

- Every Go file starts with a `package` declaration
- `package main` indicates this is an executable program
- Other packages are libraries

```go
package main  // Executable program
package utils // Library package
```

## Import Statements

Import statements bring other packages into your program:

```go
// Single import
import "fmt"

// Multiple imports
import (
    "fmt"
    "os"
    "strings"
)

// Import with alias
import f "fmt"

// Import for side effects only
import _ "image/png"
```

## The main Function

- Entry point of executable programs
- Must be in `package main`
- Takes no parameters and returns nothing

```go
func main() {
    // Program starts here
    fmt.Println("Hello, World!")
}
```

## Comments

### Single-line Comments
```go
// This is a single-line comment
fmt.Println("Hello") // Comment at end of line
```

### Multi-line Comments
```go
/*
This is a multi-line comment
It can span multiple lines
*/
```

### Documentation Comments
```go
// Package main demonstrates basic Go syntax
package main

// main is the entry point of the program
func main() {
    fmt.Println("Hello, World!")
}
```

## Identifiers and Keywords

### Identifiers
- Names for variables, functions, types, etc.
- Must start with letter or underscore
- Can contain letters, digits, and underscores
- Case-sensitive

```go
// Valid identifiers
var name string
var _count int
var userName string
var HTTP_PORT int

// Invalid identifiers
// var 2name string    // starts with digit
// var user-name string // contains hyphen
```

### Keywords (Reserved Words)

Go has 25 keywords that cannot be used as identifiers:

```
break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
```

## Naming Conventions

### Visibility Rules
- **Exported** (public): Start with uppercase letter
- **Unexported** (private): Start with lowercase letter

```go
// Exported (can be accessed from other packages)
var PublicVariable int
func PublicFunction() {}

// Unexported (only accessible within same package)
var privateVariable int
func privateFunction() {}
```

### Naming Styles
```go
// Variables and functions: camelCase
var userName string
func getUserName() string

// Constants: camelCase or ALL_CAPS
const maxRetries = 3
const MAX_BUFFER_SIZE = 1024

// Types: PascalCase
type UserAccount struct{}
type HTTPClient interface{}
```

## Semicolons

- Go automatically inserts semicolons
- You rarely need to write them explicitly
- Lexer adds semicolons at end of lines that could end a statement

```go
// These are equivalent
fmt.Println("Hello")
fmt.Println("World")

// Same as
fmt.Println("Hello"); fmt.Println("World")
```

## Code Formatting

### gofmt Tool

Go has a standard formatting tool:

```bash
# Format a file
gofmt -w main.go

# Format all Go files in current directory
gofmt -w .

# Show differences without modifying
gofmt -d main.go
```

### Standard Formatting Rules

1. **Indentation**: Use tabs, not spaces
2. **Braces**: Opening brace on same line
3. **Line length**: No strict limit, but keep reasonable
4. **Imports**: Group standard library, third-party, and local

```go
// Correct formatting
if condition {
    doSomething()
}

// Incorrect (will be auto-corrected)
if condition 
{
    doSomething()
}
```

## Hello World Examples

### Basic Hello World
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Hello World with Variables
```go
package main

import "fmt"

func main() {
    message := "Hello, World!"
    fmt.Println(message)
}
```

### Hello World with User Input
```go
package main

import "fmt"

func main() {
    var name string
    fmt.Print("Enter your name: ")
    fmt.Scanln(&name)
    fmt.Printf("Hello, %s!\n", name)
}
```

### Hello World with Function
```go
package main

import "fmt"

func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    greet("World")
    greet("Go")
}
```

## Common fmt Package Functions

```go
package main

import "fmt"

func main() {
    name := "Alice"
    age := 30
    
    // Print with newline
    fmt.Println("Hello, World!")
    
    // Print without newline
    fmt.Print("Hello ")
    fmt.Print("World\n")
    
    // Formatted print
    fmt.Printf("Name: %s, Age: %d\n", name, age)
    
    // Sprint functions (return string)
    message := fmt.Sprintf("Hello, %s!", name)
    fmt.Println(message)
}
```

## Format Verbs (Printf)

| Verb | Description | Example |
|------|-------------|----------|
| `%s` | String | `fmt.Printf("%s", "hello")` |
| `%d` | Decimal integer | `fmt.Printf("%d", 42)` |
| `%f` | Floating point | `fmt.Printf("%f", 3.14)` |
| `%t` | Boolean | `fmt.Printf("%t", true)` |
| `%v` | Default format | `fmt.Printf("%v", anything)` |
| `%+v` | Struct with field names | `fmt.Printf("%+v", struct)` |
| `%#v` | Go syntax representation | `fmt.Printf("%#v", value)` |
| `%T` | Type of value | `fmt.Printf("%T", value)` |
| `%%` | Literal percent sign | `fmt.Printf("100%%")` |

## Program Execution Flow

1. **Package initialization**: Import packages and initialize variables
2. **main function**: Program execution starts here
3. **Sequential execution**: Statements execute line by line
4. **Program termination**: When main function ends

```go
package main

import "fmt"

// Package-level variable (initialized before main)
var globalMessage = "Global variable initialized"

func main() {
    fmt.Println("1. Program starts")
    fmt.Println("2.", globalMessage)
    
    localMessage := "Local variable"
    fmt.Println("3.", localMessage)
    
    fmt.Println("4. Program ends")
}
```

## Exercises

### Exercise 1: Basic Hello World
Create a program that prints your name and favorite programming language.

```go
// Your solution here
package main

import "fmt"

func main() {
    // Print your name and favorite language
}
```

### Exercise 2: Formatted Output
Create a program that displays information about yourself using Printf.

```go
// Your solution here
package main

import "fmt"

func main() {
    name := "Your Name"
    age := 25
    city := "Your City"
    
    // Use Printf to display formatted information
}
```

### Exercise 3: Multiple Functions
Create a program with separate functions for greeting and farewell.

```go
// Your solution here
package main

import "fmt"

func greet() {
    // Greeting function
}

func farewell() {
    // Farewell function
}

func main() {
    // Call both functions
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
    fmt.Println("My name is John Doe")
    fmt.Println("My favorite programming language is Go")
}
```

### Solution 2:
```go
package main

import "fmt"

func main() {
    name := "John Doe"
    age := 25
    city := "New York"
    
    fmt.Printf("Name: %s\n", name)
    fmt.Printf("Age: %d years old\n", age)
    fmt.Printf("City: %s\n", city)
    fmt.Printf("Summary: %s is %d years old and lives in %s\n", name, age, city)
}
```

### Solution 3:
```go
package main

import "fmt"

func greet() {
    fmt.Println("Welcome to Go programming!")
}

func farewell() {
    fmt.Println("Thanks for learning Go!")
}

func main() {
    greet()
    fmt.Println("This is the main program.")
    farewell()
}
```

</details>

## Key Takeaways

1. Every Go program starts with a package declaration
2. `package main` creates an executable program
3. The `main()` function is the entry point
4. Use `fmt` package for input/output operations
5. Go automatically formats code with `gofmt`
6. Comments start with `//` or `/* */`
7. Exported names start with uppercase letters

## Next Steps

Now that you understand basic Go syntax, let's learn about [Variables and Data Types](04-variables-datatypes.md)!

---

**Previous**: [← Installation and Setup](02-installation-setup.md) | **Next**: [Variables and Data Types →](04-variables-datatypes.md)