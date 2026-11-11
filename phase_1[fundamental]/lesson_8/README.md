# Phase 1 — Go Fundamentals

## Lesson 8: Maps in Go

Welcome to Lesson 8!

In the previous lesson, we learned arrays and slices. Now we'll learn Maps, one of the most useful data structures in Go.

Maps are especially important in backend development because they allow you to store and retrieve data using key-value pairs.

You'll use maps for:

- User information
- Product lookup
- Configuration settings
- API responses
- Caching
- Counting and grouping data

## 1. What is a Map?

A Map is a collection of key-value pairs.

Each key maps to a value.

For example, consider a user:

User information

name

"Jihad"

age

30

city

"Dhaka"

isActive

true

Each key identifies its corresponding value.

In Go, the syntax is:

```
map[keyType]valueType
```

Examples:

```
map[string]string
map[string]int
map[int]string
map[string]bool
```

Meaning:

| Map type            | Meaning                        |
| ------------------- | ------------------------------ |
| `map[string]string` | String keys and string values  |
| `map[string]int`    | String keys and integer values |
| `map[int]string`    | Integer keys and string values |
| `map[string]bool`   | String keys and Boolean values |

## 2. Declaring and initializing Maps

There are several ways to create maps.

### A. Using a map literal

```
package main

import "fmt"

func main() {
    user := map[string]string{
        "name": "Jihad",
        "city": "Dhaka",
        "role": "Developer",
    }

    fmt.Println(user)
}
```

Output (key order is unspecified):

```
map[city:Dhaka name:Jihad role:Developer]
```

### B. Using `make()`

The `make()` function creates an initialized map.

```
package main

import "fmt"

func main() {
    user := make(map[string]string)

    user["name"] = "Jihad"
    user["city"] = "Dhaka"
    user["role"] = "Developer"

    fmt.Println(user)
}
```

### C. Declaring a nil map

```
var user map[string]string

fmt.Println(user == nil) // true
```

Important

A nil map can be read, but you cannot insert elements into it.

```
var user map[string]string

// user["name"] = "Jihad" // Runtime panic
```

Initialize it using `make()` or a map literal before inserting values.

## 3. Adding elements to a Map

Use the key inside square brackets.

```
package main

import "fmt"

func main() {
    products := make(map[string]int)

    products["Laptop"] = 75000
    products["Mouse"] = 1200
    products["Keyboard"] = 2500

    fmt.Println(products)
}
```

Output (order may differ):

```
map[Keyboard:2500 Laptop:75000 Mouse:1200]
```

Syntax:

```
mapName[key] = value
```

## 4. Accessing Map values

Access a value using its key.

```
package main

import "fmt"

func main() {
    prices := map[string]float64{
        "Laptop": 75000,
        "Mouse": 1200,
        "Keyboard": 2500,
    }

    fmt.Println("Laptop:", prices["Laptop"])
    fmt.Println("Mouse:", prices["Mouse"])
}
```

Output:

```
Laptop: 75000
Mouse: 1200
```

### What happens if a key doesn't exist?

Go returns the zero value of the map's value type.

```
prices := map[string]int{
    "Laptop": 75000,
}

fmt.Println(prices["Mobile"]) // 0
```

But this can be confusing because `0` might be a legitimate value.

Go provides a better solution.

## 5. Checking whether a key exists (`value, ok`)

This is one of the most important map concepts.

Go supports the comma-ok idiom to check whether a key exists.

```
package main

import "fmt"

func main() {
    prices := map[string]int{
        "Laptop": 75000,
        "Mouse": 1200,
    }

    price, ok := prices["Laptop"]

    if ok {
        fmt.Println("Price:", price)
    } else {
        fmt.Println("Product not found")
    }
}
```

Output:

```
Price: 75000
```

Another example:

```
price, ok := prices["Mobile"]

fmt.Println(price) // 0
fmt.Println(ok)    // false
```

Understanding this syntax:

```
value, ok := mapName[key]
```

| Result        | Meaning                           |
| ------------- | --------------------------------- |
| `ok == true`  | Key exists                        |
| `ok == false` | Key doesn't exist                 |
| `value`       | Corresponding value or zero value |

This is very useful when checking whether a user, product, or configuration entry exists.

## 6. Updating Map values

Assigning a value to an existing key updates it.

```
package main

import "fmt"

func main() {
    user := map[string]string{
        "name": "Jihad",
        "city": "Dhaka",
    }

    user["city"] = "Chittagong"

    fmt.Println(user["city"])
}
```

Output:

```
Chittagong
```

The same syntax is used for both inserting and updating:

```
user["city"] = "Dhaka"
```

## 7. Deleting Map elements

Go provides the built-in `delete()` function.

Syntax:

```
delete(mapName, key)
```

Example:

```
package main

import "fmt"

func main() {
    user := map[string]string{
        "name": "Jihad",
        "city": "Dhaka",
        "role": "Developer",
    }

    delete(user, "role")

    fmt.Println(user)
}
```

Output (order may differ):

```
map[city:Dhaka name:Jihad]
```

If the key doesn't exist, `delete()` does nothing.

## 8. Iterating over Maps

Use `range` to iterate over map entries.

```
package main

import "fmt"

func main() {
    user := map[string]string{
        "name": "Jihad",
        "city": "Dhaka",
        "role": "Developer",
    }

    for key, value := range user {
        fmt.Println(key, ":", value)
    }
}
```

Possible output:

```
name : Jihad
city : Dhaka
role : Developer
```

Important

Go does not guarantee map iteration order. The order can vary between iterations.

If you need predictable ordering, collect the keys and sort them.

Example:

```
package main

import (
    "fmt"
    "slices"
)

func main() {
    prices := map[string]int{
        "Laptop": 75000,
        "Mouse": 1200,
        "Keyboard": 2500,
    }

    keys := make([]string, 0, len(prices))

    for key := range prices {
        keys = append(keys, key)
    }

    slices.Sort(keys)

    for _, key := range keys {
        fmt.Println(key, prices[key])
    }
}
```

This example uses the standard library `slices` package, available since Go 1.21.

## 9. Getting the number of Map elements

Use `len()`.

```
products := map[string]int{
    "Laptop": 75000,
    "Mouse": 1200,
    "Keyboard": 2500,
}

fmt.Println("Total products:", len(products))
```

Output:

```
Total products: 3
```

Unlike slices, maps do not expose a capacity using `cap()`.

## 10. Nested Maps

A nested map is a map inside another map.

This is useful when organizing related information.

Example:

```
package main

import "fmt"

func main() {
    users := map[string]map[string]string{
        "user1": {
            "name": "Jihad",
            "city": "Dhaka",
        },
        "user2": {
            "name": "Rahim",
            "city": "Chittagong",
        },
    }

    fmt.Println(users["user1"]["name"])
    fmt.Println(users["user2"]["city"])
}
```

Output:

```
Jihad
Chittagong
```

You can also initialize nested maps separately:

```
users := make(map[string]map[string]string)

users["user1"] = make(map[string]string)
users["user1"]["name"] = "Jihad"
users["user1"]["city"] = "Dhaka"
```

Remember that an inner map must be initialized before inserting values into it.

## 11. Maps with Different Value Types

Sometimes you need to store different data types in one map.

You can use `any` (an alias for `interface{}`).

```
package main

import "fmt"

func main() {
    user := map[string]any{
        "name": "Jihad",
        "age": 30,
        "salary": 50000.50,
        "isActive": true,
    }

    fmt.Println(user)
}
```

Output (order may differ):

```
map[age:30 isActive:true name:Jihad salary:50000.5]
```

Here:

- `name` stores a string.
- `age` stores an integer.
- `salary` stores a float64.
- `isActive` stores a Boolean.

Professional recommendation: Use `map[string]any` when genuinely dynamic data is needed, such as flexible JSON. For known application data structures, prefer structs. We'll learn structs in Lesson 10.

## 12. Maps vs Arrays vs Slices

| Feature        | Array              | Slice              | Map                  |
| -------------- | ------------------ | ------------------ | -------------------- |
| Structure      | Indexed collection | Indexed collection | Key-value collection |
| Access         | Index              | Index              | Key                  |
| Size           | Fixed              | Dynamic            | Dynamic              |
| Add elements   | Not directly       | `append()`         | Assign a key         |
| Delete element | Not directly       | Custom operation   | `delete()`           |
| Iteration      | `range`            | `range`            | `range`              |
| Ordering       | Index order        | Index order        | Unspecified          |

Use:

- Arrays when you need a fixed-size collection.
- Slices when you need an ordered, flexible collection.
- Maps when you need fast key-based lookup.

## 13. Practical Project: Product Inventory Management

Let's build a small inventory management program using maps and functions.

Features:

1. Add products.
2. Display products.
3. Search for products.
4. Update product prices.
5. Delete products.

```

package main

import "fmt"

func addProduct(products map[string]float64, name string, price float64) {
    products[name] = price
    fmt.Println("Product added successfully!")
}

func displayProducts(products map[string]float64) {
    fmt.Println("\n--- Product Inventory ---")

    for name, price := range products {
        fmt.Printf("%s: %.2f TK\n", name, price)
    }

    fmt.Println("Total products:", len(products))
}

func searchProduct(products map[string]float64, name string) {
    price, ok := products[name]

    if !ok {
        fmt.Println("Product not found!")
        return
    }

    fmt.Printf("Product: %s\nPrice: %.2f TK\n", name, price)
}

func updateProduct(products map[string]float64, name string, price float64) {
    _, ok := products[name]

    if !ok {
        fmt.Println("Product not found!")
        return
    }

    products[name] = price
    fmt.Println("Product updated successfully!")
}

func deleteProduct(products map[string]float64, name string) {
    _, ok := products[name]

    if !ok {
        fmt.Println("Product not found!")
        return
    }

    delete(products, name)
    fmt.Println("Product deleted successfully!")
}

func main() {
    products := map[string]float64{
        "Laptop": 75000,
        "Mouse": 1200,
        "Keyboard": 2500,
    }

    displayProducts(products)

    addProduct(products, "Monitor", 18000)

    searchProduct(products, "Laptop")

    updateProduct(products, "Mouse", 1500)

    deleteProduct(products, "Keyboard")

    displayProducts(products)
}

```

### What did we learn from this project?

- Maps can be passed to functions.
- Changes made to a map inside a function are visible to the caller.
- `value, ok` helps validate whether a product exists.
- `delete()` removes entries.
- `len()` counts map entries.

Backend connection: This is a basic in-memory inventory store. In a real ERP application, a database would normally be the source of truth, while maps can be useful for temporary lookups and aggregations.

## 14. Practice Exercises

Try these exercises before moving to the next lesson.

Exercise 1 — Basic Map

Easy

Create a map containing your:

- Name
- Age
- City
- Profession

Print every key-value pair.

Exercise 2 — Student Marks

Easy

Create a map of student names and their marks.

- Add five students.
- Print all students and marks.
- Calculate the average marks.

Exercise 3 — Product CRUD

Medium

Create a product map with these features:

- Add a product.
- Search by product name.
- Update its price.
- Delete a product.
- Display all products.

Exercise 4 — Word Counter

Medium

Create a program that counts how many times each word appears in a sentence.

Example input:

```
go is fast and go is simple
```

Expected counts (order unspecified):

```
go: 2
is: 2
fast: 1
and: 1
simple: 1
```

Exercise 5 — Nested Map

Challenge

Create a nested map containing multiple users, with:

- Name
- Email
- City
- Role

Search for a user by ID and display their information.

Exercise 6 — Inventory Stock Counter

Challenge

Create an inventory system using `map[string]int`.

Features:

- Add stock.
- Reduce stock.
- Check available stock.
- Prevent negative stock.
- Display all inventory.

## 15. Important Go Map Concepts — Quick Revision

| Concept        | Syntax                         |
| -------------- | ------------------------------ |
| Declare map    | `var m map[string]int`         |
| Initialize map | `m := make(map[string]int)`    |
| Map literal    | `m := map[string]int{"A": 10}` |
| Insert         | `m["A"] = 10`                  |
| Read           | `value := m["A"]`              |
| Check key      | `value, ok := m["A"]`          |
| Update         | `m["A"] = 20`                  |
| Delete         | `delete(m, "A")`               |
| Iterate        | `for k, v := range m`          |
| Count          | `len(m)`                       |

Remember these five rules:

1. Maps store key-value pairs.
2. Map keys must be comparable types.
3. A nil map cannot accept new entries.
4. Map iteration order is unspecified.
5. Go maps are not safe for concurrent read/write access without synchronization.
