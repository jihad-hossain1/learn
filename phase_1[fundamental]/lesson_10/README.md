# **Phase 1 — Go Fundamentals**

## **Lesson 10: Structs & Methods in Go**

Final lesson of Phase 1

So far, you've learned variables, functions, arrays, slices, maps, and pointers.

Now it's time to learn one of the most important concepts in Go: **Structs and Methods**.

Structs allow you to combine multiple related values into a single custom data type. Methods allow you to associate behavior with that type.

These concepts are essential when building backend applications, REST APIs, inventory systems, and ERP software.

### **What you'll learn**

1. What is a struct?
2. How to declare and initialize structs.
3. Accessing and updating struct fields.
4. Structs with functions.
5. Nested structs.
6. Struct pointers.
7. Methods in Go.
8. Value receivers vs pointer receivers.
9. Struct composition.
10. A practical product management example.

## **1. What is a struct?**

A struct is a collection of fields grouped together under one custom type.

Imagine you want to store product information.

Without a struct:

```
productName := "Laptop"
productPrice := 75000.0
productStock := 10
```

You have three separate variables.

With a struct:

```
type Product struct {
    Name  string
    Price float64
    Stock int
}
```

Now you have created a custom type named `Product`.

You can create a product using:

```
product := Product{
    Name:  "Laptop",
    Price: 75000,
    Stock: 10,
}
```

Print it:

```
fmt.Println(product)
```

Output:

```
{Laptop 75000 10}
```

**Understanding the syntax:**

| Syntax    | Meaning               |
| --------- | --------------------- |
| `type`    | Declares a new type   |
| `Product` | Struct name           |
| `struct`  | Defines a struct type |
| `Name`    | Field name            |
| `string`  | Field data type       |

## **2. Different ways to initialize structs**

Go provides several ways to initialize a struct.

**Method 1 — Named fields (recommended)**

```
product := Product{
    Name:  "Keyboard",
    Price: 1200,
    Stock: 25,
}
```

This is the most readable approach.

**Method 2 — Positional initialization**

```
product := Product{
    "Mouse",
    800,
    15,
}
```

The values must follow the exact field order.

For maintainability, prefer named fields.

**Method 3 — Declare first, assign later**

```
var product Product

product.Name = "Monitor"
product.Price = 15000
product.Stock = 5
```

**Method 4 — Using a pointer**

```
product := &Product{
    Name:  "Laptop",
    Price: 75000,
    Stock: 10,
}
```

Here, `product` is a pointer to a `Product`.

## **3. Accessing and updating struct fields**

Use the dot (`.`) operator.

```
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
    Stock int
}

func main() {
    product := Product{
        Name:  "Laptop",
        Price: 75000,
        Stock: 10,
    }

    // Read fields
    fmt.Println("Name:", product.Name)
    fmt.Println("Price:", product.Price)

    // Update fields
    product.Price = 80000
    product.Stock = 15

    fmt.Println("Updated:", product)
}
```

Output:

```
Name: Laptop
Price: 75000
Updated: {Laptop 80000 15}
```

You can also compare struct values when all their fields are comparable:

```
p1 := Product{Name: "Mouse", Price: 800, Stock: 10}
p2 := Product{Name: "Mouse", Price: 800, Stock: 10}

fmt.Println(p1 == p2) // true
```

## **4. Struct zero values**

When you declare a struct without initializing its fields, Go assigns each field its zero value.

```
type Product struct {
    Name  string
    Price float64
    Stock int
}

func main() {
    var product Product

    fmt.Println(product)
}
```

Output:

```
{ 0 0}
```

| Data type | Zero value |
| --------- | ---------- |
| `string`  | `""`       |
| `int`     | `0`        |
| `float64` | `0.0`      |
| `bool`    | `false`    |
| Pointer   | `nil`      |

This is important because Go doesn't require you to initialize every field manually.

## **5. Structs with functions**

You can pass structs to functions just like other values.

```
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
    Stock int
}

func printProduct(product Product) {
    fmt.Println("Product:", product.Name)
    fmt.Println("Price:", product.Price)
    fmt.Println("Stock:", product.Stock)
}

func main() {
    product := Product{
        Name:  "Keyboard",
        Price: 1200,
        Stock: 20,
    }

    printProduct(product)
}
```

Remember: Go passes arguments by value.

This means `printProduct` receives a copy of the struct.

If you want to modify the original struct, use a pointer:

```
func updatePrice(product *Product, price float64) {
    product.Price = price
}
```

Usage:

```
updatePrice(&product, 1500)
```

## **6. Nested structs**

A struct can contain another struct as a field.

This is useful when modeling real-world relationships.

Example: A customer has an address.

```
package main

import "fmt"

type Address struct {
    City    string
    Country string
}

type Customer struct {
    Name    string
    Email   string
    Address Address
}

func main() {
    customer := Customer{
        Name:  "Rahim",
        Email: "rahim@example.com",
        Address: Address{
            City:    "Dhaka",
            Country: "Bangladesh",
        },
    }

    fmt.Println(customer.Name)
    fmt.Println(customer.Address.City)
    fmt.Println(customer.Address.Country)
}
```

Output:

```
Rahim
Dhaka
Bangladesh
```

Notice:

```
Address Address
```

The first `Address` is the field name, and the second `Address` is the struct type.

You can build more complex structures using this approach.

**Struct relationship**

**Customer**

Name · Email

**Address**

City · Country

## **7. Struct pointers**

You learned pointers in Lesson 9. Now let's combine pointers with structs.

```
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
}

func main() {
    product := Product{
        Name:  "Monitor",
        Price: 15000,
    }

    ptr := &product

    fmt.Println(ptr.Name)

    ptr.Price = 18000

    fmt.Println(product.Price)
}
```

Output:

```
Monitor
18000
```

Go automatically dereferences struct pointers when accessing fields.

These two expressions are equivalent:

```
ptr.Price = 18000
```

```
(*ptr).Price = 18000
```

## **8. What are methods in Go?**

A method is a function associated with a particular type.

In Go, you can define methods on your own types.

Let's start with a simple example.

```
package main

import "fmt"

type Product struct {
    Name  string
    Price float64
}

func (p Product) display() {
    fmt.Println("Product:", p.Name)
    fmt.Println("Price:", p.Price)
}

func main() {
    product := Product{
        Name:  "Laptop",
        Price: 75000,
    }

    product.display()
}
```

Output:

```
Product: Laptop
Price: 75000
```

Look carefully at this declaration:

```
func (p Product) display()
```

The `(p Product)` part is called the **receiver**.

- `p` is the receiver variable.
- `Product` is the receiver type.
- `display()` is the method name.

Compare a normal function and a method:

```
// Normal function
func display(p Product) {}

// Method
func (p Product) display() {}
```

A method belongs to a type, while a normal function does not have a receiver.

## **9. Value receiver vs pointer receiver**

This is an important Go concept.

There are two common receiver types.

### **A. Value receiver**

A value receiver receives a copy of the struct.

```
type Product struct {
    Name  string
    Price float64
}

func (p Product) updatePrice(price float64) {
    p.Price = price
}
```

Usage:

```
product := Product{
    Name:  "Laptop",
    Price: 75000,
}

product.updatePrice(80000)

fmt.Println(product.Price)
```

Output:

```
75000
```

Why?

Because the method modifies only its local copy.

### **B. Pointer receiver**

A pointer receiver can modify the original struct.

```
func (p *Product) updatePrice(price float64) {
    p.Price = price
}
```

Usage:

```
product := Product{
    Name:  "Laptop",
    Price: 75000,
}

product.updatePrice(80000)

fmt.Println(product.Price)
```

Output:

```
80000
```

Go automatically takes the address of an addressable variable when calling a pointer-receiver method.

For example:

```
product.updatePrice(80000)
```

is effectively:

```
(&product).updatePrice(80000)
```

### **Comparison**

| Feature               | Value receiver         | Pointer receiver |
| --------------------- | ---------------------- | ---------------- |
| Syntax                | `(p Product)`          | `(p *Product)`   |
| Receives              | Copy                   | Pointer          |
| Modifies original     | No                     | Yes              |
| Useful for            | Small immutable values | Updating structs |
| Can avoid struct copy | No                     | Yes              |

**Best practice:** For mutable structs, generally use pointer receivers consistently. Use value receivers for small, immutable value types when copying is appropriate.

## **10. Struct composition**

Go does not use traditional class inheritance.

Instead, Go encourages composition: combining smaller types to build larger types.

One feature used for this is struct embedding.

```
package main

import "fmt"

type Person struct {
    Name  string
    Email string
}

type Employee struct {
    Person
    Position string
    Salary   float64
}

func main() {
    employee := Employee{
        Person: Person{
            Name:  "Rahim",
            Email: "rahim@example.com",
        },
        Position: "Software Engineer",
        Salary:   50000,
    }

    fmt.Println(employee.Name)
    fmt.Println(employee.Email)
    fmt.Println(employee.Position)
}
```

Output:

```
Rahim
rahim@example.com
Software Engineer
```

Notice:

```
type Employee struct {
    Person
    Position string
}
```

`Person` is embedded in `Employee`.

Its fields can be accessed directly through `employee.Name`, although the actual field belongs to the embedded `Person`.

This is called **struct embedding**.

## **11. Exported and unexported fields**

Go uses capitalization to control visibility across packages.

```
type Product struct {
    Name  string
    Price float64
    stock int
}
```

- `Name` — exported; accessible from other packages.
- `Price` — exported.
- `stock` — unexported; accessible only within the package where it is declared.

This is important when designing packages and APIs.

For example, a package might expose product information but keep internal stock-management details private.

## **12. Practical project: Product management**

Let's combine what you've learned into a small product model.

```

package main

import "fmt"

type Product struct {
    ID    int
    Name  string
    Price float64
    Stock int
}

func (p *Product) UpdatePrice(price float64) {
    if price > 0 {
        p.Price = price
    }
}

func (p *Product) AddStock(quantity int) {
    if quantity > 0 {
        p.Stock += quantity
    }
}

func (p *Product) RemoveStock(quantity int) bool {
    if quantity <= 0 || quantity > p.Stock {
        return false
    }

    p.Stock -= quantity
    return true
}

func (p Product) Display() {
    fmt.Println("ID:", p.ID)
    fmt.Println("Name:", p.Name)
    fmt.Println("Price:", p.Price)
    fmt.Println("Stock:", p.Stock)
}

func main() {
    product := Product{
        ID:    1,
        Name:  "Laptop",
        Price: 75000,
        Stock: 10,
    }

    product.UpdatePrice(80000)
    product.AddStock(5)

    success := product.RemoveStock(3)

    if success {
        fmt.Println("Stock removed successfully")
    }

    product.Display()
}

```

Expected output:

```
Stock removed successfully
ID: 1
Name: Laptop
Price: 80000
Stock: 12
```

### **What did we practice?**

- Creating a custom struct.
- Initializing struct fields.
- Defining methods.
- Using pointer receivers.
- Updating product information.
- Returning boolean values from methods.
- Validating stock operations.

This is a basic model of how you might represent a product in a Go backend application.

## **13. Important concepts to remember**

| Concept          | Key idea                           |
| ---------------- | ---------------------------------- |
| Struct           | Groups related fields              |
| Field            | A variable inside a struct         |
| Nested struct    | A struct inside another struct     |
| Struct pointer   | Holds the address of a struct      |
| Method           | Function with a receiver           |
| Value receiver   | Works with a copy                  |
| Pointer receiver | Can modify the original            |
| Embedding        | Combines types through composition |
| Exported field   | Starts with uppercase              |
| Unexported field | Starts with lowercase              |

## **14. Practice exercises**

Try these independently before moving to Phase 2.

**Your practice progress**

0/7

Create a Student struct with ID, Name, Age, and Grade fields.

Create three Student objects and print their information.

Create an Address struct and embed it inside a Customer struct.

Write a method named UpdateAge that changes a student's age using a pointer receiver.

Create a BankAccount struct with Balance and write Deposit and Withdraw methods.

Create a Product struct with AddStock, RemoveStock, and UpdatePrice methods.

Create an Employee struct by embedding a Person struct.

## **Phase 1 — Completion**

**Go Fundamentals**

10/10 Completed

1. Introduction to Go
2. Variables and Data Types
3. Operators and Expressions
4. Input and Output
5. Control Flow
6. Functions
7. Arrays and Slices
8. Maps
9. Pointers
10. Structs and Methods
