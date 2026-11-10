# Phase 1 — Go Fundamentals

## Lesson 7: Arrays and Slices in Go

Arrays and slices are essential Go data structures. They allow you to store and manage multiple values using a single variable.

As a backend developer, you'll frequently use slices when working with:

- Product lists
- Customer records
- Invoice items
- API responses
- Database query results

Let's understand them step by step.

## 1. What is an Array?

An array is a collection of elements of the same data type with a fixed length.

Example:

```
package main

import "fmt"

func main() {
    var numbers [5]int

    numbers[0] = 10
    numbers[1] = 20
    numbers[2] = 30
    numbers[3] = 40
    numbers[4] = 50

    fmt.Println(numbers)
}
```

Output:

```
[10 20 30 40 50]
```

Here:

- `[5]` means the array contains exactly five elements.
- `int` means every element must be an integer.
- Array indexing starts from `0`.

### Array syntax

```
var arrayName [size]dataType
```

Examples:

```
var ages [5]int
var prices [10]float64
var names [3]string
var active [4]bool
```

## 2. Array initialization

You can initialize an array when declaring it.

```
package main

import "fmt"

func main() {
    numbers := [5]int{10, 20, 30, 40, 50}

    fmt.Println(numbers)
}
```

You can also let Go determine the array length:

```
numbers := [...]int{10, 20, 30, 40, 50}

fmt.Println(numbers)
```

The `...` tells Go to calculate the array length automatically.

## 3. Accessing and updating array elements

Use an index to access or update an element.

```
package main

import "fmt"

func main() {
    fruits := [3]string{"Apple", "Banana", "Orange"}

    fmt.Println(fruits[0])
    fmt.Println(fruits[1])

    fruits[1] = "Mango"

    fmt.Println(fruits)
}
```

Output:

```
Apple
Banana
[Apple Mango Orange]
```

Important: Accessing an index outside the array's valid range causes a runtime panic.

For example, an array of length 3 has indexes `0`, `1`, and `2`.

## 4. Array length

Use the built-in `len()` function to determine the number of elements.

```
package main

import "fmt"

func main() {
    numbers := [5]int{10, 20, 30, 40, 50}

    fmt.Println("Array:", numbers)
    fmt.Println("Length:", len(numbers))
}
```

Output:

```
Array: [10 20 30 40 50]
Length: 5
```

## 5. Iterating over arrays

There are two common ways to iterate through an array.

A. Traditional for loop

```
package main

import "fmt"

func main() {
    numbers := [5]int{10, 20, 30, 40, 50}

    for i := 0; i < len(numbers); i++ {
        fmt.Println(numbers[i])
    }
}
```

B. Using range

```
package main

import "fmt"

func main() {
    numbers := [5]int{10, 20, 30, 40, 50}

    for index, value := range numbers {
        fmt.Println("Index:", index, "Value:", value)
    }
}
```

Output:

```
Index: 0 Value: 10
Index: 1 Value: 20
Index: 2 Value: 30
Index: 3 Value: 40
Index: 4 Value: 50
```

If you don't need the index, use `_`:

```
for _, value := range numbers {
    fmt.Println(value)
}
```

# 6. What is a Slice?

Very important

A slice is a flexible, dynamically sized view into an underlying array.

Unlike arrays, slices can grow or shrink as needed.

In real-world Go backend development, slices are used much more frequently than arrays.

Example:

```
package main

import "fmt"

func main() {
    numbers := []int{10, 20, 30, 40, 50}

    fmt.Println(numbers)
}
```

Notice the difference:

```
// Array
numbers := [5]int{10, 20, 30, 40, 50}

// Slice
numbers := []int{10, 20, 30, 40, 50}
```

Array has a fixed length. Slice does not have a fixed length.

## 7. Creating slices

There are several ways to create slices.

A. Using a slice literal

```
names := []string{"Jihad", "Rahim", "Karim"}

fmt.Println(names)
```

B. Using `make()`

The `make()` function creates an initialized slice.

```
numbers := make([]int, 5)

fmt.Println(numbers)
```

Output:

```
[0 0 0 0 0]
```

Syntax:

```
slice := make([]dataType, length, capacity)
```

Example:

```
numbers := make([]int, 3, 5)

fmt.Println("Slice:", numbers)
fmt.Println("Length:", len(numbers))
fmt.Println("Capacity:", cap(numbers))
```

Output:

```
Slice: [0 0 0]
Length: 3
Capacity: 5
```

Here:

- Length = 3 (currently accessible elements).
- Capacity = 5 (available storage before reallocation is needed).

## 8. Accessing and updating slice elements

Slices work similarly to arrays.

```
package main

import "fmt"

func main() {
    fruits := []string{"Apple", "Banana", "Orange"}

    fmt.Println(fruits[0])

    fruits[1] = "Mango"

    fmt.Println(fruits)
}
```

Output:

```
Apple
[Apple Mango Orange]
```

## 9. Adding elements using `append()`

The `append()` function adds elements to a slice.

```
package main

import "fmt"

func main() {
    numbers := []int{10, 20, 30}

    numbers = append(numbers, 40)
    numbers = append(numbers, 50)

    fmt.Println(numbers)
}
```

Output:

```
[10 20 30 40 50]
```

You can append multiple elements:

```
numbers = append(numbers, 60, 70, 80)
```

You can also append another slice using `...`:

```
first := []int{10, 20, 30}
second := []int{40, 50, 60}

result := append(first, second...)

fmt.Println(result)
```

Output:

```
[10 20 30 40 50 60]
```

Important: Always assign the result of `append()` back to a slice variable because it may return a new slice with a different underlying array.

## 10. Slice operations

You can extract a portion of an existing slice using slicing syntax.

```
package main

import "fmt"

func main() {
    numbers := []int{10, 20, 30, 40, 50}

    fmt.Println(numbers[1:4])
    fmt.Println(numbers[:3])
    fmt.Println(numbers[2:])
    fmt.Println(numbers[:])
}
```

Output:

```
[20 30 40]
[10 20 30]
[30 40 50]
[10 20 30 40 50]
```

Syntax:

```
slice[start:end]
```

Remember:

- `start` is inclusive.
- `end` is exclusive.

For example, `numbers[1:4]` returns indexes 1, 2, and 3.

## 11. Understanding length and capacity

This is one of the most important slice concepts.

```
package main

import "fmt"

func main() {
    numbers := make([]int, 3, 5)

    fmt.Println("Length:", len(numbers))
    fmt.Println("Capacity:", cap(numbers))

    numbers = append(numbers, 10)

    fmt.Println("After append:")
    fmt.Println("Length:", len(numbers))
    fmt.Println("Capacity:", cap(numbers))
}
```

Expected output:

```
Length: 3
Capacity: 5

After append:
Length: 4
Capacity: 5
```

Visual explanation:

Slice: make([]int, 3, 5)

## 0

Index 0

## 0

Index 1

## 0

Index 2

## —

Index 3

## —

Index 4

Length: 3

Capacity: 5

Green cells represent the current slice length. Grey cells represent additional capacity.

When you append an element, the length increases. If the slice exceeds its capacity, Go allocates a new underlying array and copies the elements.

Important: The exact growth strategy is an implementation detail. Don't depend on a specific capacity increase.

## 12. Copying slices using `copy()`

Go provides the built-in `copy()` function to copy elements from one slice into another.

```
package main

import "fmt"

func main() {
    source := []int{10, 20, 30, 40, 50}

    destination := make([]int, len(source))

    copied := copy(destination, source)

    fmt.Println("Source:", source)
    fmt.Println("Destination:", destination)
    fmt.Println("Copied elements:", copied)
}
```

Output:

```
Source: [10 20 30 40 50]
Destination: [10 20 30 40 50]
Copied elements: 5
```

`copy()` copies up to the length of the shorter slice.

This creates an independent copy for these integer elements. For slices of pointers, maps, or nested reference-containing values, the copied elements can still refer to shared data.

## 13. Removing an element from a slice

Go does not provide a built-in `remove()` function for slices.

You can remove an element using `append()`.

```
package main

import "fmt"

func main() {
    numbers := []int{10, 20, 30, 40, 50}

    index := 2

    numbers = append(numbers[:index], numbers[index+1:]...)

    fmt.Println(numbers)
}
```

Output:

```
[10 20 40 50]
```

We removed the element at index 2 (`30`).

## 14. Practical project: Product inventory manager

Let's combine arrays, slices, loops, and functions into a small inventory program.

```

package main

import "fmt"

func displayProducts(products []string) {
    fmt.Println("\n--- Product Inventory ---")

    for index, product := range products {
        fmt.Printf("%d. %s\n", index+1, product)
    }

    fmt.Println("Total products:", len(products))
}

func addProduct(products []string, product string) []string {
    products = append(products, product)
    return products
}

func main() {
    products := []string{
        "Keyboard",
        "Mouse",
        "Monitor",
    }

    displayProducts(products)

    products = addProduct(products, "Laptop")

    displayProducts(products)
}

```

Expected output:

```
--- Product Inventory ---
1. Keyboard
2. Mouse
3. Monitor
Total products: 3

--- Product Inventory ---
1. Keyboard
2. Mouse
3. Monitor
4. Laptop
Total products: 4
```

Notice how we passed a slice to a function and returned the updated slice.

This pattern is very common in Go applications.

## 15. Array vs Slice — Important comparison

| Feature              | Array                  | Slice                          |
| -------------------- | ---------------------- | ------------------------------ |
| Size                 | Fixed                  | Dynamic                        |
| Declaration          | `[5]int`               | `[]int`                        |
| Length               | Fixed                  | Can change                     |
| Capacity             | Same as length         | Can exceed length              |
| `append()`           | Not directly supported | Supported                      |
| `copy()`             | Not directly supported | Supported                      |
| Underlying array     | Is the array           | References an underlying array |
| Common backend usage | Less frequent          | Very frequent                  |

## 16. Practice exercises

Try these yourself before moving to the next lesson.

Exercise 1 — Array basics&#x20;

Easy

- Create an integer array containing 10 numbers.
- Print all elements.
- Calculate the sum of all elements.

Exercise 2 — Find maximum and minimum&#x20;

Easy

- Create an array of 10 integers.
- Find the maximum and minimum values.

Exercise 3 — Slice operations&#x20;

Medium

- Create a slice of five product names.
- Add two more products using `append()`.
- Print the first three products.
- Print the total number of products.

Exercise 4 — Remove an element&#x20;

Medium

- Create a slice of integers.
- Remove an element by its index.
- Print the updated slice.

Exercise 5 — Copy a slice&#x20;

Medium

- Create a source slice.
- Create a destination slice.
- Copy all elements.
- Modify the source and verify that the destination remains unchanged.

Exercise 6 — Inventory management&#x20;

Challenge

\
Build a small program with functions to:

- Add products.
- List products.
- Search for a product.
- Remove a product.
- Count total products.

## 17. Important concepts to remember

Before moving forward, make sure you understand these:

1. Arrays have fixed lengths.
2. Slices are dynamically sized views over underlying arrays.
3. `len()` returns the current number of elements.
4. `cap()` returns the capacity.
5. `append()` adds elements.
6. `copy()` copies elements between slices.
7. `range` iterates over arrays and slices.
8. Slicing uses an inclusive start and exclusive end index.
9. Multiple slices can share the same underlying array, so modifications may affect each other.
