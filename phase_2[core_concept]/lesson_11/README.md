# **Phase 2 — Go Core Concepts**

## **Lesson 11: Strings and Runes in Go**

Phase 2 started

Welcome to Phase 2! In Phase 1, you learned Go's fundamental syntax, functions, pointers, structs, and methods.

Now we'll explore Go's core concepts in more depth, starting with **Strings and Runes**.

Strings are everywhere in backend development: usernames, emails, product names, JSON responses, database records, and API requests.

By the end of this lesson, you'll understand:

1. How strings work in Go.
2. String length and indexing.
3. Bytes vs runes.
4. Unicode and UTF-8.
5. String iteration.
6. String manipulation.
7. `strings` and `strconv` packages.
8. Practical text-processing examples.

## **1. What is a string in Go?**

A string is an immutable sequence of bytes that usually represents UTF-8 encoded text.

Example:

```
package main

import "fmt"

func main() {
    name := "Jihad"
    message := "Hello, Go Programming!"

    fmt.Println(name)
    fmt.Println(message)
}
```

Output:

```
Jihad
Hello, Go Programming!
```

**Important:** Strings in Go are immutable. You cannot change an individual character directly.

This is invalid:

```
name := "Jihad"

name[0] = 'j' // Compilation error
```

Instead, create a new string:

```
name := "Jihad"
name = "j" + name[1:]

fmt.Println(name)
```

Output:

```
jihad
```

## **2. String declaration**

You can declare strings in several ways.

```
package main

import "fmt"

func main() {
    // Short declaration
    name := "Jihad"

    // Explicit type
    var city string = "Dhaka"

    // Zero value
    var country string

    // Constant
    const language = "Go"

    fmt.Println(name)
    fmt.Println(city)
    fmt.Println(country)
    fmt.Println(language)
}
```

Output:

```
Jihad
Dhaka

Go
```

The zero value of a string is an empty string: `""`.

### **Two types of string literals**

**Interpreted string**

Uses double quotes. Escape sequences are interpreted.

```
message := "Hello\nWorld"
fmt.Println(message)
```

Output:

```
Hello
World
```

**Raw string**

Uses backticks. Escape sequences are not interpreted.

```
message := `Hello\nWorld`

fmt.Println(message)
```

Output:

```
Hello\nWorld
```

Raw strings are useful for multiline text, SQL queries, and JSON examples.

## **3. String length —&#x20;**`len()`

The built-in `len()` function returns the number of bytes in a string, not necessarily the number of characters.

```
package main

import "fmt"

func main() {
    name := "Jihad"

    fmt.Println(len(name))
}
```

Output:

```
5
```

For English text, the byte count often matches the number of characters.

But consider Bengali:

```
package main

import "fmt"

func main() {
    text := "বাংলা"

    fmt.Println(len(text))
}
```

The output is:

```
15
```

Why 15?

Because UTF-8 characters can require multiple bytes.

This introduces two important concepts:

- **Byte:** A unit of memory containing 8 bits.
- **Rune:** A Go type representing a Unicode code point.

A Unicode character can require multiple bytes in UTF-8.

## **4. What is a rune?**

In Go, `rune` is an alias for `int32`.

It represents a Unicode code point.

Example:

```
package main

import "fmt"

func main() {
    var letter rune = 'A'
    var bangla rune = 'ক'

    fmt.Println(letter)
    fmt.Println(bangla)

    fmt.Printf("%c\n", letter)
    fmt.Printf("%c\n", bangla)
}
```

Output:

```
65
2453
A
ক
```

Notice:

- `%c` prints the Unicode character.
- `%d` prints its numeric code point.
- `%T` prints the variable's type.

```
fmt.Printf("%T\n", bangla)
```

Output:

```
int32
```

### **Byte vs Rune**

| Feature     | Byte             | Rune               |
| ----------- | ---------------- | ------------------ |
| Go type     | `byte` (`uint8`) | `rune` (`int32`)   |
| Represents  | 8-bit value      | Unicode code point |
| Typical use | UTF-8 encoding   | Unicode text       |
| Size        | 1 byte           | 4 bytes            |

## **5. String indexing**

You can access individual bytes using an index.

```
package main

import "fmt"

func main() {
    text := "Go"

    fmt.Println(text[0])
    fmt.Println(text[1])

    fmt.Printf("%c\n", text[0])
    fmt.Printf("%c\n", text[1])
}
```

Output:

```
71
111
G
o
```

But be careful with Unicode strings.

```
text := "বাংলা"

fmt.Println(text[0])
```

This prints the first UTF-8 byte, not the first Bengali character.

To work with Unicode characters, convert the string to a rune slice:

```
package main

import "fmt"

func main() {
    text := "বাংলা"

    letters := []rune(text)

    fmt.Println(len(text))
    fmt.Println(len(letters))

    fmt.Printf("%c\n", letters[0])
}
```

Output:

```
15
5
ব
```

**Remember:**

```
len(text)       // Number of bytes
len([]rune(text)) // Number of Unicode code points
```

Note: A rune is not always the same as a user-perceived character. Some visible characters consist of multiple Unicode code points.

## **6. Iterating over strings**

Go provides two common ways to iterate over strings.

### **A. Using a traditional&#x20;**`for`**&#x20;loop**

```
package main

import "fmt"

func main() {
    text := "Go"

    for i := 0; i < len(text); i++ {
        fmt.Printf("%c\n", text[i])
    }
}
```

Output:

```
G
o
```

This iterates over bytes.

### **B. Using&#x20;**`range`

```
package main

import "fmt"

func main() {
    text := "বাংলা"

    for index, value := range text {
        fmt.Println(index, string(value))
    }
}
```

Output:

```
0 ব
3 া
6 ং
9 ল
12 া
```

Why do the indexes increase by 3?

Because each Bengali character in this example occupies 3 bytes in UTF-8.

**Best practice:** Use `range` when you want to iterate over Unicode code points.

## **7. String concatenation**

You can combine strings using the `+` operator.

```
firstName := "Jihad"
lastName := "Hossain"

fullName := firstName + " " + lastName

fmt.Println(fullName)
```

Output:

```
Jihad Hossain
```

For repeated concatenation in a loop, prefer `strings.Builder`.

```
package main

import (
    "fmt"
    "strings"
)

func main() {
    var builder strings.Builder

    builder.WriteString("Hello")
    builder.WriteString(" ")
    builder.WriteString("Go")

    result := builder.String()

    fmt.Println(result)
}
```

Output:

```
Hello Go
```

`strings.Builder` is useful for efficiently building strings from many parts.

## **8. Important strings package functions**

Go provides a standard library package named `strings`.

Import it:

```
import "strings"
```

Here are some commonly used functions:

| Function               | Purpose                           |
| ---------------------- | --------------------------------- |
| `strings.Contains()`   | Check whether text exists         |
| `strings.HasPrefix()`  | Check starting text               |
| `strings.HasSuffix()`  | Check ending text                 |
| `strings.ToUpper()`    | Convert to uppercase              |
| `strings.ToLower()`    | Convert to lowercase              |
| `strings.TrimSpace()`  | Remove surrounding whitespace     |
| `strings.ReplaceAll()` | Replace all occurrences           |
| `strings.Split()`      | Split a string                    |
| `strings.Join()`       | Join a slice of strings           |
| `strings.Count()`      | Count non-overlapping occurrences |
| `strings.Index()`      | Find the first byte index         |

### **Practical examples**

```
package main

import (
    "fmt"
    "strings"
)

func main() {
    name := "  Jihad Hossain  "

    fmt.Println(strings.TrimSpace(name))
    fmt.Println(strings.ToUpper(name))
    fmt.Println(strings.ToLower(name))

    fmt.Println(strings.Contains(name, "Jihad"))
    fmt.Println(strings.HasPrefix(name, "  Jihad"))
    fmt.Println(strings.HasSuffix(name, "  "))

    fmt.Println(strings.ReplaceAll(name, "Jihad", "Rahim"))
}
```

Output:

```
Jihad Hossain
  JIHAD HOSSAIN
  jihad hossain
true
true
true
  Rahim Hossain
```

## **9. Split and Join strings**

These are particularly useful when processing API input.

### **Split**

```
package main

import (
    "fmt"
    "strings"
)

func main() {
    categories := "Electronics,Clothing,Food"

    result := strings.Split(categories, ",")

    fmt.Println(result)
}
```

Output:

```
[Electronics Clothing Food]
```

The result is a slice of strings: `[]string`.

### **Join**

```
package main

import (
    "fmt"
    "strings"
)

func main() {
    categories := []string{
        "Electronics",
        "Clothing",
        "Food",
    }

    result := strings.Join(categories, " | ")

    fmt.Println(result)
}
```

Output:

```
Electronics | Clothing | Food
```

## **10. String conversion with strconv**

Sometimes you need to convert between strings and numbers.

Go provides the `strconv` package.

```
import "strconv"
```

**String to integer:**

```
age, err := strconv.Atoi("30")

if err != nil {
    fmt.Println("Invalid number")
    return
}

fmt.Println(age)
```

Output:

```
30
```

**Integer to string:**

```
age := 30

result := strconv.Itoa(age)

fmt.Println(result)
```

Output:

```
30
```

Remember that `Itoa` returns a string, not an integer.

**String to float:**

```
price, err := strconv.ParseFloat("1250.50", 64)

if err != nil {
    fmt.Println("Invalid price")
    return
}

fmt.Println(price)
```

**Float to string:**

```
price := 1250.50

result := strconv.FormatFloat(price, 'f', 2, 64)

fmt.Println(result)
```

Output:

```
1250.50
```

Notice that Go requires explicit error handling for conversions that can fail.

## **11. Practical project: Product search**

Let's build a small product search function using strings.

Requirements:

- Store products in a slice.
- Accept a search keyword.
- Ignore uppercase/lowercase differences.
- Find products whose names contain the keyword.
- Display matching products.

```

package main

import (
    "fmt"
    "strings"
)

type Product struct {
    ID    int
    Name  string
    Price float64
}

func searchProducts(products []Product, keyword string) {
    keyword = strings.ToLower(strings.TrimSpace(keyword))

    for _, product := range products {
        name := strings.ToLower(product.Name)

        if strings.Contains(name, keyword) {
            fmt.Printf(
                "ID: %d | Name: %s | Price: %.2f\n",
                product.ID,
                product.Name,
                product.Price,
            )
        }
    }
}

func main() {
    products := []Product{
        {1, "Laptop", 75000},
        {2, "Keyboard", 1200},
        {3, "Mouse", 800},
        {4, "Monitor", 15000},
        {5, "Gaming Laptop", 95000},
    }

    searchProducts(products, "LAPTOP")
}

```

Output:

```
ID: 1 | Name: Laptop | Price: 75000.00
ID: 5 | Name: Gaming Laptop | Price: 95000.00
```

This is a simple example of text processing that you can later extend into a real product-search API.

## **12. Common mistakes to avoid**

| Mistake                                            | Correct approach                                     |
| -------------------------------------------------- | ---------------------------------------------------- |
| Assuming `len(string)` counts characters           | It counts bytes                                      |
| Modifying a string by index                        | Create a new string                                  |
| Using byte indexing for Bengali text               | Use `range` or `[]rune`                              |
| Ignoring conversion errors                         | Check `err`                                          |
| Using `+` repeatedly in large loops                | Consider `strings.Builder`                           |
| Assuming rune count equals visible character count | Understand Unicode code points and grapheme clusters |

## **13. Practice exercises**

Try these exercises in your Go environment.

**Practice progress**

0/7

Declare a string containing your full name and print its byte length.

Count Unicode code points in the Bengali string "বাংলাদেশ" using []rune.

Write a function that reverses a string by Unicode code points.

Convert a string "1250.75" into float64 and handle errors.

Split "Laptop,Keyboard,Mouse,Monitor" into a slice and print each product.

Write a case-insensitive product search function using strings.Contains.

Count how many times the letter "a" occurs in "banana" using strings.Count.

## **Phase 2 Progress**

**Go Core Concepts**

1 of 8 lessons

1.  Current lesson

    &#x20;Strings and Runes

2.  Packages and Modules

3.  Go Standard Library

4.  Working with JSON

5.  File Handling

6.  Deeper Understanding of Interfaces and Types

7.  Generics

8.  Defer, Panic and Recover

**Next lesson — Lesson 12: Packages and Modules in Go**

We'll learn how to organize Go applications into packages, create modules using `go mod`, import custom packages, manage dependencies, and structure a real backend project.

I can also give you a three-question quiz on strings, bytes, and runes to check your understanding before Lesson 12.
