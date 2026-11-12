# Phase 1 — Go Fundamentals

## Lesson 9: Pointers in Go

Pointers are an important Go concept. They help you work with memory addresses, modify variables through functions, and avoid unnecessary copying of large data structures.

By the end of this lesson, you will understand:

- What pointers are.
- Address-of (`&`) and dereference (`*`) operators.
- How to declare and use pointers.
- How to modify variables through pointers.
- Nil pointers and pointer safety.
- How pointers work with functions and structs.

## 1. What is a pointer?

A normal variable stores a value.

A pointer stores the memory address of another variable.

For example:

```
package main

import "fmt"

func main() {
    age := 30

    fmt.Println("Value:", age)
    fmt.Println("Address:", &age)
}
```

Possible output:

```
Value: 30
Address: 0xc000012080
```

The memory address will be different when you run the program.

Here:

- `age` stores `30`.
- `&age` gives the memory address of `age`.

## 2. Declare a pointer

A pointer's type is written using `*`.

```
var age int = 30

var ptr *int

ptr = &age

fmt.Println(ptr)
```

Explanation:

```
var ptr *int
```

This declares a pointer that can hold the address of an integer.

```
ptr = &age
```

This assigns the address of `age` to `ptr`.

We can also declare and initialize a pointer in one statement:

```
age := 30
ptr := &age
```

## 3. Dereferencing a pointer

Dereferencing means accessing the value stored at the memory address held by a pointer.

Use the `*` operator.

```
package main

import "fmt"

func main() {
    age := 30

    ptr := &age

    fmt.Println("Age:", age)
    fmt.Println("Pointer:", ptr)
    fmt.Println("Value through pointer:", *ptr)
}
```

Possible output:

```
Age: 30
Pointer: 0xc000012080
Value through pointer: 30
```

Understand these three expressions:

| Expression | Meaning                        |
| ---------- | ------------------------------ |
| `age`      | Variable value                 |
| `&age`     | Memory address                 |
| `*ptr`     | Value at the pointer's address |

Important: `*` has two uses:

```
var ptr *int // Pointer type
value := *ptr // Dereferencing
```

## 4. Modify a variable using a pointer

You can change the original variable through its pointer.

```
package main

import "fmt"

func main() {
    price := 500.0

    ptr := &price

    *ptr = 750.0

    fmt.Println("Updated price:", price)
}
```

Output:

```
Updated price: 750
```

Why did `price` change?

Because `ptr` contains the address of `price`. Writing to `*ptr` updates the value at that address.

How the pointer works

ptr

Stores the address of price

Address of price

price

# 750.0

Original variable updated through pointer

## 5. Pointers with functions

This is one of the most useful applications of pointers.

First, consider a normal function:

```
package main

import "fmt"

func updatePrice(price float64) {
    price = 1000
}

func main() {
    price := 500.0

    updatePrice(price)

    fmt.Println(price)
}
```

Output:

```
500
```

Why? Because Go passes the value of `price` to the function. The function receives a copy.

Now let's use a pointer.

```
package main

import "fmt"

func updatePrice(price *float64) {
    *price = 1000
}

func main() {
    price := 500.0

    updatePrice(&price)

    fmt.Println(price)
}
```

Output:

```
1000
```

Explanation:

1. `&price` passes the address of the original variable.
2. `price *float64` receives that address.
3. `*price = 1000` updates the original variable.

Important Go rule: Go always passes arguments by value. When you pass a pointer, Go copies the pointer value (the address), so both pointers refer to the same underlying variable. Go does not have traditional pass-by-reference parameters.

## 6. Nil pointers

A pointer that does not point to a valid value has the value `nil`.

```
package main

import "fmt"

func main() {
    var ptr *int

    fmt.Println(ptr)
}
```

Output:

```
<nil>
```

You must not dereference a nil pointer.

This will panic:

```
var ptr *int

fmt.Println(*ptr) // Runtime panic
```

Instead, check it first:

```
package main

import "fmt"

func main() {
    var ptr *int

    if ptr == nil {
        fmt.Println("Pointer is nil")
        return
    }

    fmt.Println(*ptr)
}
```

Output:

```
Pointer is nil
```

## 7. Using `new()` to create pointers

Go provides the built-in `new()` function.

```
package main

import "fmt"

func main() {
    ptr := new(int)

    fmt.Println(*ptr)

    *ptr = 100

    fmt.Println(*ptr)
}
```

Output:

```
0
100
```

`new(int)` allocates a zero-valued integer and returns a pointer to it.

Other examples:

```
intPtr := new(int)
floatPtr := new(float64)
boolPtr := new(bool)
```

You will usually use `&variable` when you already have a variable. `new()` is useful when you need a pointer to a zero-valued allocation.

## 8. Pointers with structs

Pointers are frequently used with structs in backend applications.

```
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
}

func updateProduct(product *Product) {
    product.Price = 1500
}

func main() {
    product := Product{
        Name:  "Laptop",
        Price: 1200,
    }

    updateProduct(&product)

    fmt.Println(product)
}
```

Output:

```
{Laptop 1500}
```

Notice:

```
product.Price = 1500
```

Go automatically dereferences the pointer when accessing struct fields. This is equivalent to:

```
(*product).Price = 1500
```

You don't need to write `(*product).Price` in normal code.

## 9. Pointer best practices

| Situation                           | Recommendation                             |
| ----------------------------------- | ------------------------------------------ |
| Need to modify an original variable | Use a pointer                              |
| Working with large structs          | Consider a pointer to avoid copying        |
| Need an optional value              | A pointer can represent absence with `nil` |
| Working with small values           | Passing by value is often simpler          |
| Pointer may be nil                  | Check before dereferencing                 |
| Slice or map modification           | A pointer is usually unnecessary           |

A few important points:

- Go has no pointer arithmetic in ordinary safe Go code.
- You don't need to manually allocate or free memory.
- Go's garbage collector manages memory.
- Go can safely return pointers to local variables; the compiler manages their lifetime.

## 10. Practical example: Inventory management

Let's update a product's stock and price using pointers.

```
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
    Stock int
}

func updateInventory(product *Product, newPrice float64, quantity int) {
    product.Price = newPrice
    product.Stock += quantity
}

func main() {
    product := Product{
        Name:  "Keyboard",
        Price: 800,
        Stock: 10,
    }

    updateInventory(&product, 950, 5)

    fmt.Println("Product:", product.Name)
    fmt.Println("Price:", product.Price)
    fmt.Println("Stock:", product.Stock)
}
```

Output:

```
Product: Keyboard
Price: 950
Stock: 15
```

This pattern is useful when building inventory, invoice, purchase, and other backend modules.

## 11. Practice exercises

Try solving these without looking at the solutions.

Exercise checklist

0/6 completed

Declare an integer variable and print its memory address.

Create a pointer to a variable and print the value through the pointer.

Change a variable's value using a pointer.

Write a function named swap that swaps two integer values using pointers.

Create a Product struct and write a function that updates its price and stock using a pointer.

Write a function that safely handles a nil pointer.

### Challenge: Swap two numbers

Complete this function:

```
func swap(a *int, b *int) {
    // Your code here
}
```

Expected behavior:

```
x := 10
y := 20

swap(&x, &y)

fmt.Println(x, y)
```

Output:

```
20 10
```
