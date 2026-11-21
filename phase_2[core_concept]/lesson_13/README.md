# **Phase 2 — Go Core Concepts**

## **Lesson 13: Go Standard Library**

Phase 2

Lesson 13 of 18

Intermediate Beginner

In this lesson, you'll learn how to use Go's **Standard Library** to build practical applications without installing unnecessary third-party packages.

Since you already have backend development experience with Node.js, Express and NestJS, I'll connect some Go packages to familiar backend concepts.

## **1. What is the Go Standard Library?**

The Go Standard Library is a collection of built-in packages that come with Go.

For example:

| Go package      | Purpose                       | Similar concept in Node.js |
| --------------- | ----------------------------- | -------------------------- |
| `fmt`           | Formatting and printing       | `console.log()`            |
| `strings`       | String manipulation           | String methods             |
| `strconv`       | String conversion             | `Number()`, `String()`     |
| `math`          | Mathematical operations       | `Math`                     |
| `time`          | Date and time                 | `Date`                     |
| `os`            | Environment and OS operations | `process`, `fs`            |
| `path/filepath` | File paths                    | `path`                     |
| `encoding/json` | JSON encoding/decoding        | `JSON.stringify()`         |
| `net/http`      | HTTP server/client            | Node HTTP                  |
| `log`           | Logging                       | Console/logger             |
| `errors`        | Error creation                | Error objects              |

**Important:** You don't need to install these packages. They are included with Go.

## **2. The&#x20;**`fmt`**&#x20;package — Formatting**

You've already used `fmt.Println()`. Now let's explore some additional functions.

```
package main

import "fmt"

func main() {
    name := "Laptop"
    price := 45000.50
    quantity := 3

    fmt.Println("Product:", name)

    fmt.Printf("Price: %.2f\n", price)
    fmt.Printf("Quantity: %d\n", quantity)

    fmt.Printf("Product: %s | Price: %.2f | Quantity: %d\n",
        name, price, quantity)
}
```

**Output:**

```
Product: Laptop
Price: 45000.50
Quantity: 3
Product: Laptop | Price: 45000.50 | Quantity: 3
```

### **Important formatting verbs**

| Verb   | Meaning              | Example   |
| ------ | -------------------- | --------- |
| `%s`   | String               | `"Go"`    |
| `%d`   | Integer              | `100`     |
| `%f`   | Floating-point       | `10.50`   |
| `%.2f` | Two decimal places   | `10.50`   |
| `%t`   | Boolean              | `true`    |
| `%v`   | Default value format | Any value |
| `%T`   | Variable type        | `int`     |
| `%q`   | Quoted string        | `"Hello"` |

### `Sprintf()`**&#x20;vs&#x20;**`Printf()`

`Printf()` prints formatted output.

`Sprintf()` returns a formatted string.

```
price := 1250.50

message := fmt.Sprintf("Product price: %.2f BDT", price)

fmt.Println(message)
```

This is particularly useful when generating API response messages, logs, and formatted reports.

## **3. The&#x20;**`strings`**&#x20;package — String operations**

You've already learned about strings and runes. Now let's explore some useful functions from the standard library.

```
package main

import (
    "fmt"
    "strings"
)

func main() {
    name := "  Go Programming Language  "

    fmt.Println(strings.TrimSpace(name))
    fmt.Println(strings.ToUpper(name))
    fmt.Println(strings.ToLower(name))
    fmt.Println(strings.Contains(name, "Programming"))
    fmt.Println(strings.HasPrefix(name, "  Go"))
    fmt.Println(strings.ReplaceAll(name, "Go", "Golang"))

    parts := strings.Split("apple,banana,mango", ",")

    fmt.Println(parts)
}
```

**Output:**

```
Go Programming Language
  GO PROGRAMMING LANGUAGE
  go programming language
true
true
  Golang Programming Language
[apple banana mango]
```

### **Common string functions**

| Function               | Purpose                           |
| ---------------------- | --------------------------------- |
| `strings.TrimSpace()`  | Remove surrounding whitespace     |
| `strings.ToUpper()`    | Convert to uppercase              |
| `strings.ToLower()`    | Convert to lowercase              |
| `strings.Contains()`   | Check substring                   |
| `strings.HasPrefix()`  | Check beginning                   |
| `strings.HasSuffix()`  | Check ending                      |
| `strings.ReplaceAll()` | Replace occurrences               |
| `strings.Split()`      | Split into a slice                |
| `strings.Join()`       | Join a slice                      |
| `strings.Count()`      | Count non-overlapping occurrences |

**Backend example:** When validating an email address, normalizing a search query, or cleaning user input, these functions are useful.

## **4. The&#x20;**`strconv`**&#x20;package — Type conversion**

Go is a statically typed language. You often need to convert values between strings and numeric types.

```
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // String to integer
    quantity, err := strconv.Atoi("25")
    if err != nil {
        fmt.Println("Invalid quantity:", err)
        return
    }

    fmt.Println(quantity + 5)

    // String to float
    price, err := strconv.ParseFloat("1250.75", 64)
    if err != nil {
        fmt.Println("Invalid price:", err)
        return
    }

    fmt.Println(price)

    // Integer to string
    id := strconv.Itoa(1001)
    fmt.Println(id)

    // Float to string
    formattedPrice := strconv.FormatFloat(price, 'f', 2, 64)
    fmt.Println(formattedPrice)
}
```

**Output:**

```
30
1250.75
1001
1250.75
```

### **Important functions**

| Function                | Conversion       |
| ----------------------- | ---------------- |
| `strconv.Atoi()`        | String → int     |
| `strconv.Itoa()`        | int → string     |
| `strconv.ParseInt()`    | String → int64   |
| `strconv.ParseFloat()`  | String → float64 |
| `strconv.FormatInt()`   | int64 → string   |
| `strconv.FormatFloat()` | float64 → string |
| `strconv.ParseBool()`   | String → bool    |
| `strconv.FormatBool()`  | bool → string    |

**Remember:** Conversion functions such as `Atoi()` and `ParseFloat()` return an error. Always check it when working with user input.

## **5. The&#x20;**`math`**&#x20;package — Mathematical operations**

Go provides built-in mathematical functions.

```
package main

import (
    "fmt"
    "math"
)

func main() {
    fmt.Println(math.Abs(-25.5))
    fmt.Println(math.Ceil(10.2))
    fmt.Println(math.Floor(10.9))
    fmt.Println(math.Round(10.6))
    fmt.Println(math.Max(100, 200))
    fmt.Println(math.Min(100, 200))
    fmt.Println(math.Sqrt(144))
}
```

**Output:**

```
25.5
11
10
11
200
100
12
```

### **Useful functions**

| Function           | Result |
| ------------------ | ------ |
| `math.Abs(-10)`    | `10`   |
| `math.Ceil(10.2)`  | `11`   |
| `math.Floor(10.9)` | `10`   |
| `math.Round(10.6)` | `11`   |
| `math.Max(10, 20)` | `20`   |
| `math.Min(10, 20)` | `10`   |
| `math.Sqrt(144)`   | `12`   |
| `math.Pow(2, 3)`   | `8`    |

**Practical note:** Financial applications should generally use integer minor units (such as paisa) or a suitable decimal representation for monetary calculations. `float64` can introduce precision errors.

## **6. The&#x20;**`time`**&#x20;package — Date and time**

This is one of the most important packages for backend development.

You will frequently need to work with:

- Invoice creation dates
- Order timestamps
- Payment dates
- Subscription expiration
- Report date ranges
- API timestamps

### **A. Get current time**

```
package main

import (
    "fmt"
    "time"
)

func main() {
    now := time.Now()

    fmt.Println("Current time:", now)
    fmt.Println("Year:", now.Year())
    fmt.Println("Month:", now.Month())
    fmt.Println("Day:", now.Day())
}
```

### **B. Format dates**

Go uses a special reference time:

`Mon Jan 2 15:04:05 MST 2006`

You must use this reference layout when formatting dates.

```
package main

import (
    "fmt"
    "time"
)

func main() {
    now := time.Now()

    fmt.Println(now.Format("2006-01-02"))
    fmt.Println(now.Format("02-01-2006"))
    fmt.Println(now.Format("2006-01-02 15:04:05"))
}
```

Example output:

```
2026-11-21
21-11-2026
2026-11-21 20:23:00
```

The actual time depends on when you run the program.

### **C. Parse a date string**

```
package main

import (
    "fmt"
    "time"
)

func main() {
    dateString := "2026-11-21"

    parsedDate, err := time.Parse("2006-01-02", dateString)

    if err != nil {
        fmt.Println("Invalid date:", err)
        return
    }

    fmt.Println(parsedDate)
}
```

### **D. Add and subtract time**

```
package main

import (
    "fmt"
    "time"
)

func main() {
    now := time.Now()

    tomorrow := now.Add(24 * time.Hour)
    yesterday := now.Add(-24 * time.Hour)

    fmt.Println("Today:", now)
    fmt.Println("Tomorrow:", tomorrow)
    fmt.Println("Yesterday:", yesterday)
}
```

### **E. Calculate duration**

```
start := time.Now()

time.Sleep(2 * time.Second)

end := time.Now()

duration := end.Sub(start)

fmt.Println("Duration:", duration)
```

Output:

```
Duration: 2s
```

**Backend tip:** Use `time.Time` for timestamps and `time.Duration` for elapsed time or timeouts. For database applications, be deliberate about UTC storage and local timezone presentation.

## **7. The&#x20;**`os`**&#x20;package — Operating system operations**

The `os` package helps you interact with the operating system.

In backend development, one important use is reading environment variables.

### **A. Read environment variables**

```
package main

import (
    "fmt"
    "os"
)

func main() {
    dbHost := os.Getenv("DB_HOST")

    if dbHost == "" {
        dbHost = "localhost"
    }

    fmt.Println("Database Host:", dbHost)
}
```

Run:

```
DB_HOST=127.0.0.1 go run .
```

Output:

```
Database Host: 127.0.0.1
```

### **B. Set environment variables**

```
os.Setenv("APP_ENV", "development")

environment := os.Getenv("APP_ENV")

fmt.Println(environment)
```

### **C. Read command-line arguments**

```
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Arguments:", os.Args)
}
```

Run:

```
go run . hello world
```

The arguments are available in `os.Args`. Index `0` is the program name; subsequent elements are the supplied arguments.

## **8. The&#x20;**`path/filepath`**&#x20;package — File paths**

This package handles filesystem paths in a platform-independent way.

```
package main

import (
    "fmt"
    "path/filepath"
)

func main() {
    path := filepath.Join("uploads", "products", "image.png")

    fmt.Println("Path:", path)
    fmt.Println("Directory:", filepath.Dir(path))
    fmt.Println("Filename:", filepath.Base(path))
    fmt.Println("Extension:", filepath.Ext(path))
}
```

Output on Linux:

```
Path: uploads/products/image.png
Directory: uploads/products
Filename: image.png
Extension: .png
```

Use `filepath.Join()` instead of manually joining path strings.

## **9. Other important standard library packages**

You don't need to master all these packages today. Just understand what they do.

| Package         | What it does                                            |
| --------------- | ------------------------------------------------------- |
| `encoding/json` | Convert Go data to/from JSON                            |
| `net/http`      | Build HTTP clients and servers                          |
| `log`           | Application logging                                     |
| `errors`        | Create and inspect errors                               |
| `context`       | Carry cancellation, deadlines and request-scoped values |
| `sort`          | Sort collections                                        |
| `regexp`        | Regular expression matching                             |
| `bufio`         | Buffered input and output                               |
| `io`            | General input/output interfaces                         |

We'll study JSON, errors, HTTP and concurrency-related concepts in their dedicated lessons and phases.

## **10. Practical project — Product Inventory Summary**

Let's combine several standard library packages in one small program.

**Requirements:**

- Store product information.
- Clean product names.
- Calculate total inventory value.
- Format prices.
- Display the current timestamp.

```

package main

import (
    "fmt"
    "strings"
    "time"
)

type Product struct {
    Name     string
    Price    float64
    Quantity int
}

func main() {
    products := []Product{
        {"  Laptop  ", 75000.50, 2},
        {" Keyboard ", 2500.00, 5},
        {" Mouse ", 800.00, 10},
    }

    var totalValue float64

    fmt.Println("INVENTORY REPORT")
    fmt.Println(strings.Repeat("-", 40))

    for _, product := range products {
        name := strings.TrimSpace(product.Name)

        value := product.Price * float64(product.Quantity)
        totalValue += value

        fmt.Printf(
            "%-12s | Price: %10.2f | Qty: %d | Value: %.2f\n",
            name,
            product.Price,
            product.Quantity,
            value,
        )
    }

    fmt.Println(strings.Repeat("-", 40))

    fmt.Printf("Total Inventory Value: %.2f BDT\n", totalValue)

    fmt.Println(
        "Generated At:",
        time.Now().Format("2006-01-02 15:04:05"),
    )
}

```

**Expected output:**

```
INVENTORY REPORT
----------------------------------------
Laptop       | Price:   75000.50 | Qty: 2 | Value: 150001.00
Keyboard     | Price:    2500.00 | Qty: 5 | Value: 12500.00
Mouse        | Price:     800.00 | Qty: 10 | Value: 8000.00
----------------------------------------
Total Inventory Value: 170501.00 BDT
Generated At: 2026-11-21 20:23:00
```

This example combines:

- `fmt` — formatted output
- `strings` — cleaning product names
- `time` — timestamps
- Structs and slices — storing product data
- Loops — calculating inventory values

## **11. How to explore Go documentation**

A Go developer should become comfortable reading official documentation.

Run these commands in your terminal:

```
# Check Go version
go version

# Open documentation for fmt
go doc fmt

# Open documentation for time
go doc time

# Read a specific function
go doc strings.Contains

# List all standard library packages
go list std
```

Official reference:

[Go Standard Library Documentation](https://pkg.go.dev/std?utm_source=chatgpt.com)

You don't need to memorize every function. Learn how to find the right package and read its documentation.

## **12. Common mistakes and best practices**

| Mistake                                          | Better approach                               |
| ------------------------------------------------ | --------------------------------------------- |
| Ignoring conversion errors                       | Always check `err`                            |
| Using `float64` for precise money                | Use integer minor units or decimal arithmetic |
| Hardcoding database configuration                | Read environment variables                    |
| Manually joining file paths                      | Use `filepath.Join()`                         |
| Formatting dates like JavaScript                 | Use Go's reference time layout                |
| Installing packages unnecessarily                | Check the standard library first              |
| Using `fmt.Println()` for all production logging | Learn structured logging with `log/slog`      |

## **13. Practice exercises**

Try these yourself before looking for a solution.

**Exercise 1 — String utilities**

Create a function that:

1. Accepts a full name.
2. Removes surrounding whitespace.
3. Converts it to uppercase.
4. Counts the number of words.

Example:

```
Input:  "  jihad hossain  "
Output:
Name: JIHAD HOSSAIN
Words: 2
```

**Exercise 2 — Price calculator**

Create a program that:

1. Accepts a product price as a string.
2. Converts it into `float64`.
3. Accepts a quantity.
4. Calculates the total.
5. Displays the result with two decimal places.

Example:

```
Price: "1250.50"
Quantity: 4
Total: 5002.00
```

**Exercise 3 — Invoice due date**

Create a program that:

1. Gets the current date.
2. Adds 30 days.
3. Displays the invoice date.
4. Displays the due date in `YYYY-MM-DD` format.

## **Phase 2 progress**

**Go Core Concepts**

# **3/8**

- Lesson 11 — Strings and Runes
- Lesson 12 — Packages and Modules
- Lesson 13 — Go Standard Library
- Lesson 14 — Working with JSON
- Lesson 15 — File Handling
- Lesson 16 — Interfaces and Types (Deeper Understanding)
- Lesson 17 — Generics
- Lesson 18 — Defer, Panic and Recover

**Next lesson:** Phase 2 — Lesson 14: Working with JSON in Go (`encoding/json`), including structs, JSON tags, marshaling, unmarshaling and building API-ready JSON responses.

I can review your solutions to all three exercises and explain any mistakes before we move to Lesson 14.
