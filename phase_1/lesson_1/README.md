# Phase 1 — Go Fundamentals

## Lesson 1: Introduction to Go (Golang)

Welcome to your Go programming journey! We'll learn Go from basic to advanced, step by step, with practical examples and exercises.

## 1. What is Go?

Go (Golang) is an open-source, statically typed, compiled programming language developed by Google.

It is designed to be:

- Simple and readable
- Fast and efficient
- Excellent for backend development
- Powerful for concurrent programming
- Suitable for APIs, microservices, CLI tools, and cloud applications

[Writing Better Modules in Go: Best Practices for Design and Testing | by Rado Salov | Medium](https://images.openai.com/static-rsc-4/wJWrO09B_KDEjdb1YOsqVHFzjsAnP626dS5QZO8HIVYyI__bk6nyB9zmwiDhcZYJTytQlbymz2lobHImd1mcVlpyNt6yC00YhxsM262H9ad4uvqvi-doH1qrtJCfJt5FxFp_F7w8JU3ldLV1sDjRsyp6VUZuaKeQdFGZKLWAhvA?purpose=inline)

[Введение в Golang: Основы программирования на Go — Stepik](https://images.openai.com/static-rsc-4/S7GxkhQJTTWXnzTDXOd4DuVvhho6-w_UNCaLaND_Nc0JifFIjBS8FJ5NqTTeT0p8pa0kzUr-mMGrXPj28CuAfLCumgz5P_EbMgbwUQg2nA2df3f_Ud7KSk3ww9nD_umFCn7hJqGT5iVhutdzWPo89RQrsRYXFRy8BGkeqh6UErg?purpose=inline)

## 2. Why learn Go?

As a backend developer, Go is especially useful for building:

| Technology         | Examples                     |
| ------------------ | ---------------------------- |
| REST APIs          | Go + Gin / Fiber             |
| Microservices      | Distributed backend services |
| CLI tools          | Developer automation         |
| Cloud applications | Docker, Kubernetes           |
| Concurrent systems | Goroutines, channels         |

## 3. Install Go

Download Go from the official website:

[Download Go — Official Website](https://go.dev/dl/?utm_source=chatgpt.com)

For Fedora Linux, you can also install it using:

```
sudo dnf install golang
```

Verify installation:

```
go version
```

Expected output (version may differ):

```
go version go1.x.x linux/amd64
```

## 4. Your first Go program

Create a project directory:

```
mkdir hello-go
cd hello-go
```

Initialize a Go module:

```
go mod init hello-go
```

Create a file named `main.go`.

```

package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}

```

Run the program:

```
go run main.go
```

Output:

```
Hello, World!
```

## 5. Understanding the code

| Code            | Meaning                                      |
| --------------- | -------------------------------------------- |
| `package main`  | Defines the executable program's package     |
| `import "fmt"`  | Imports Go's formatting and printing package |
| `func main()`   | Entry point of the program                   |
| `fmt.Println()` | Prints text to the terminal                  |

Important: Go executes your program starting from the `main()` function.

## 6. Go syntax rules

```
// Single-line comment

/*
   Multi-line comment
*/

package main

import "fmt"

func main() {
    fmt.Println("Learning Go")
    fmt.Println(100)
    fmt.Println(true)
}
```

Output:

```
Learning Go
100
true
```

Notice:

- Go statements usually don't require semicolons.
- Curly braces `{}` define blocks.
- Go is case-sensitive.
- Unused imported packages and declared variables cause compilation errors.
- Go automatically formats code using `gofmt`.

Format your file:

```
gofmt -w main.go
```

## 7. Practice exercises

Try these yourself before looking for solutions.

Exercise 1: Print your name.

Exercise 2: Print your name, age, and profession on separate lines.

Exercise 3: Print the following output:

```
Hello, Go!
I am learning backend development.
Go is awesome!
```

Exercise 4: Create a program that prints the result of:

```
25 + 15
100 - 40
12 * 5
100 / 4
```

## What comes next?

Lesson 2

Variables and Constants in Go

- `var` keyword
- Short declaration `:=`
- Constants using `const`
- Data types: `int`, `float64`, `string`, `bool`
- Zero values and type inference
