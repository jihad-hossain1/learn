# **Phase 2 — Go Core Concepts**

## **Lesson 12: Packages and Modules in Go**

In the previous lesson, you learned about strings, runes, Unicode, and string manipulation.

Now we'll learn how to organize Go applications into multiple files and packages.

This is particularly important when building production applications such as ERP systems, REST APIs, inventory management, and microservices.

### **What you'll learn**

1. What are packages in Go?
2. The `package` declaration.
3. Exported vs unexported identifiers.
4. Importing standard library packages.
5. Creating your own packages.
6. Understanding Go modules.
7. `go.mod` and `go.sum`.
8. Managing dependencies.
9. Organizing a backend project.
10. Practical multi-package application.

## **1. What is a package?**

A package is a way to organize related Go source files and functionality.

For example, Go provides standard library packages:

| Package         | Purpose                        |
| --------------- | ------------------------------ |
| `fmt`           | Input and output               |
| `strings`       | String manipulation            |
| `strconv`       | Type conversion                |
| `time`          | Date and time                  |
| `math`          | Mathematical operations        |
| `net/http`      | HTTP server and client         |
| `encoding/json` | JSON encoding and decoding     |
| `os`            | Operating system functionality |

You use packages through the `import` statement.

```
package main

import "fmt"

func main() {
    fmt.Println("Hello Go!")
}
```

Explanation:

- `package main` declares the current package.
- `import "fmt"` imports the formatting package.
- `fmt.Println()` calls an exported function from that package.
- `main()` is the entry point of an executable program.

## **2. Understanding package declarations**

Every Go source file must start with a package declaration.

Example:

```
package main
```

Or:

```
package user
```

There are two common categories of packages:

**A. Executable package**

```
package main
```

A program intended to run as an executable must have a `main` package and a `main()` function.

**B. Reusable package**

```
package calculator
```

This package provides functionality that other packages can import.

For example:

```
package calculator

func Add(a int, b int) int {
    return a + b
}
```

This package can be imported and reused by another package.

## **3. Exported vs unexported identifiers**

This is one of the most important rules in Go.

Go controls visibility using the first letter of an identifier.

- Uppercase first letter = exported.
- Lowercase first letter = unexported.

Example:

```
package product

type Product struct {
    Name  string
    Price float64

    stock int
}

func CreateProduct() {
    // Exported function
}

func validatePrice() {
    // Unexported function
}
```

| Identifier        | Visibility |
| ----------------- | ---------- |
| `Product`         | Exported   |
| `Name`            | Exported   |
| `Price`           | Exported   |
| `stock`           | Unexported |
| `CreateProduct()` | Exported   |
| `validatePrice()` | Unexported |

An exported identifier can be accessed from another package.

An unexported identifier can only be accessed within its own package.

For example, another package can do:

```
product := product.Product{
    Name:  "Laptop",
    Price: 75000,
}
```

But it cannot directly access:

```
product.stock
```

This helps you hide internal implementation details.

## **4. Importing multiple packages**

You can import multiple packages in one import block.

```
package main

import (
    "fmt"
    "strings"
    "strconv"
)

func main() {
    name := "  Go Programming  "

    cleanName := strings.TrimSpace(name)

    age, err := strconv.Atoi("30")
    if err != nil {
        fmt.Println("Invalid age")
        return
    }

    fmt.Println(cleanName)
    fmt.Println(age)
}
```

Go's standard formatter automatically organizes imports with `gofmt`.

Run:

```
gofmt -w main.go
```

## **5. What is a Go module?**

A Go module is a collection of related Go packages managed together.

A module typically contains:

- A `go.mod` file.
- Go source files.
- Packages.
- External dependencies, if needed.

Think of a module as the root of a Go project.

Let's create your first Go module.

### **Step 1 — Create a project folder**

```
mkdir go-learning
cd go-learning
```

### **Step 2 — Initialize a module**

```
go mod init example.com/go-learning
```

Output:

```
go: creating new go.mod: module example.com/go-learning
```

This creates a `go.mod` file.

### **Step 3 — Inspect go.mod**

```
module example.com/go-learning

go 1.24.0
```

Your Go version directive may differ depending on your installed Go version.

Let's understand these lines:

| Directive | Meaning                                                                     |
| --------- | --------------------------------------------------------------------------- |
| `module`  | Module path                                                                 |
| `go`      | Minimum Go language/toolchain compatibility version specified by the module |

The module path identifies your module and is used in import paths.

For example:

```
import "example.com/go-learning/product"
```

If you later publish the module at a Git repository, the module path usually corresponds to its repository location.

## **6. Create your own package**

Let's create a project with two packages.

Project structure:

```
go-learning/
│
├── go.mod
├── main.go
│
└── calculator/
    └── calculator.go
```

### **Step 1 — Create calculator package**

Create `calculator/calculator.go`:

```
package calculator

func Add(a int, b int) int {
    return a + b
}

func Subtract(a int, b int) int {
    return a - b
}

func Multiply(a int, b int) int {
    return a * b
}
```

Notice:

```
package calculator
```

This declares the package name.

### **Step 2 — Import your package**

Create `main.go`:

```
package main

import (
    "fmt"

    "example.com/go-learning/calculator"
)

func main() {
    result := calculator.Add(10, 20)

    fmt.Println("Result:", result)
}
```

Output:

```
Result: 30
```

Here:

```
"example.com/go-learning/calculator"
```

This is the import path.

And:

```
calculator.Add(10, 20)
```

This calls the exported `Add` function.

**Important:** Replace `example.com/go-learning` with the actual module path from your `go.mod` file.

## **7. Multiple files in the same package**

Go allows multiple `.go` files to belong to the same package.

Example:

```
calculator/
├── calculator.go
└── validation.go
```

`calculator.go`:

```
package calculator

func Add(a int, b int) int {
    return a + b
}
```

`validation.go`:

```
package calculator

func IsPositive(number int) bool {
    return number > 0
}
```

Both files belong to the same package.

You can use `IsPositive()` inside `calculator.go` without importing anything.

**Remember:** Files in the same directory normally belong to the same package. You cannot declare different package names in ordinary `.go` files within that directory.

## **8. Understanding go.sum**

When you install external dependencies, Go may create a `go.sum` file.

For example, suppose your project imports an external package.

```
import "github.com/google/uuid"
```

You can add the dependency using:

```
go get github.com/google/uuid
```

Go downloads the dependency and updates the module files as necessary.

Your project might look like:

```
go-learning/
│
├── go.mod
├── go.sum
├── main.go
└── calculator/
    └── calculator.go
```

### **Difference between go.mod and go.sum**

| File     | Responsibility                                       |
| -------- | ---------------------------------------------------- |
| `go.mod` | Module identity, Go version, dependency requirements |
| `go.sum` | Checksums used to verify downloaded module content   |

Do not manually delete or edit these files without understanding the consequences.

Both files should normally be committed to Git for an application.

## **9. Useful Go module commands**

These are important commands you'll use regularly.

| Command           | Purpose                                    |
| ----------------- | ------------------------------------------ |
| `go mod init`     | Initialize a module                        |
| `go mod tidy`     | Add missing and remove unused dependencies |
| `go get`          | Add or update dependencies                 |
| `go mod download` | Download dependencies                      |
| `go mod verify`   | Verify downloaded dependencies             |
| `go list ./...`   | List packages in the module                |
| `go test ./...`   | Run tests across packages                  |
| `go build ./...`  | Build packages                             |
| `go run .`        | Run the current main package               |

For example:

```
go mod tidy
```

This is particularly useful after adding or removing imports.

## **10. Organizing a backend Go project**

Now let's think about a real backend application.

Suppose you're developing an ERP system with:

- Products
- Customers
- Invoices
- Inventory
- Authentication
- Database access

A possible project structure:

```
erp-api/
│
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── product/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   ├── customer/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   ├── invoice/
│   │   ├── model.go
│   │   ├── service.go
│   │   └── handler.go
│   │
│   └── auth/
│       ├── service.go
│       └── middleware.go
│
├── pkg/
│   └── response/
│       └── response.go
│
├── configs/
│
├── go.mod
├── go.sum
└── README.md
```

Let's understand the folders.

| Folder      | Responsibility                                      |
| ----------- | --------------------------------------------------- |
| `cmd/`      | Application entry points                            |
| `internal/` | Private application packages                        |
| `pkg/`      | Packages intended for reuse outside the application |
| `configs/`  | Configuration files                                 |
| `go.mod`    | Module definition                                   |
| `go.sum`    | Dependency checksums                                |

### **What is special about internal?**

Go has a special rule for directories named `internal`.

Packages inside an `internal` directory can only be imported by packages within the permitted parent directory tree.

For example:

```
erp-api/
├── internal/
│   └── product/
└── cmd/
    └── api/
```

Your API application can import `erp-api/internal/product`.

An unrelated external module cannot import that internal package.

This is useful for keeping application-specific implementation private.

**Architecture advice:** Don't create too many packages too early. Start with clear domain boundaries and split functionality when it makes the application easier to maintain.

## **11. Practical project: Product package**

Let's build a small multi-package product application.

Project structure:

```
product-app/
│
├── go.mod
├── main.go
│
└── product/
    ├── model.go
    └── service.go
```

### `product/model.go`

```
package product

type Product struct {
    ID    int
    Name  string
    Price float64
    Stock int
}
```

### `product/service.go`

```
package product

import "errors"

func CreateProduct(name string, price float64, stock int) (Product, error) {
    if name == "" {
        return Product{}, errors.New("product name is required")
    }

    if price <= 0 {
        return Product{}, errors.New("price must be greater than zero")
    }

    if stock < 0 {
        return Product{}, errors.New("stock cannot be negative")
    }

    return Product{
        ID:    1,
        Name:  name,
        Price: price,
        Stock: stock,
    }, nil
}
```

### `main.go`

```
package main

import (
    "fmt"

    "example.com/product-app/product"
)

func main() {
    item, err := product.CreateProduct("Laptop", 75000, 10)

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Product created successfully")
    fmt.Println("ID:", item.ID)
    fmt.Println("Name:", item.Name)
    fmt.Println("Price:", item.Price)
    fmt.Println("Stock:", item.Stock)
}
```

Output:

```
Product created successfully
ID: 1
Name: Laptop
Price: 75000
Stock: 10
```

### **Run the project**

```
go mod init example.com/product-app

go run .
```

Make sure the import path in `main.go` matches the module path in `go.mod`.

This example demonstrates how you can separate product models and business logic into a reusable package.

## **12. Common mistakes**

| Mistake                                             | Correct approach                                                              |
| --------------------------------------------------- | ----------------------------------------------------------------------------- |
| Importing a package with the wrong path             | Match the module path and directory                                           |
| Calling an unexported function from another package | Capitalize its name if it should be exported                                  |
| Using different package names in one directory      | Keep ordinary source files in that directory in the same package              |
| Forgetting to initialize a module                   | Run `go mod init`                                                             |
| Not updating dependency metadata                    | Run `go mod tidy`                                                             |
| Creating unnecessary packages                       | Organize around clear responsibilities                                        |
| Confusing package name with import path             | Import path identifies location; package name identifies the declared package |

## **13. Practice exercises**

Try these exercises in your Fedora Go environment.

**Practice checklist**

0/7

Create a Go module named example.com/student-app.

Create a package named student with a Student struct.

Add an exported function CreateStudent() that returns a Student.

Add an unexported function validateAge() inside the student package.

Import your student package into main.go and print a student.

Create two separate files inside the same package and call a function from one file in the other.

Run go mod tidy, go list ./..., and go build ./....

### **Challenge: Build a small package-based application**

Create this structure:

```
student-app/
├── go.mod
├── main.go
└── student/
    ├── model.go
    └── service.go
```

Requirements:

1. Create a `Student` struct with `ID`, `Name`, and `Age`.
2. Implement `CreateStudent()`.
3. Validate that the name isn't empty.
4. Validate that age is greater than zero.
5. Return an error if validation fails.
6. Import and use the package from `main.go`.

Try to implement this without copying the Product example.
