# 13. Methods

Methods in Go are functions with a special receiver argument. They allow you to define functions that operate on specific types, providing a way to add behavior to your custom types. Methods are fundamental to Go's approach to object-oriented programming.

## What are Methods?

A method is a function with a receiver. The receiver appears between the `func` keyword and the method name. Methods allow you to call functions on instances of types.

### Basic Method Syntax

```go
func (receiver ReceiverType) MethodName(parameters) ReturnType {
    // method body
}
```

### Simple Method Example

```go
package main

import (
    "fmt"
    "math"
)

type Circle struct {
    Radius float64
}

// Method with value receiver
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Method with value receiver
func (c Circle) Circumference() float64 {
    return 2 * math.Pi * c.Radius
}

func main() {
    circle := Circle{Radius: 5.0}
    
    fmt.Printf("Circle radius: %.2f\n", circle.Radius)
    fmt.Printf("Area: %.2f\n", circle.Area())
    fmt.Printf("Circumference: %.2f\n", circle.Circumference())
}
```

## Value Receivers vs Pointer Receivers

### Value Receivers

Value receivers work with a copy of the original value. Changes made inside the method don't affect the original.

```go
package main

import "fmt"

type Counter struct {
    Value int
}

// Value receiver - works with a copy
func (c Counter) IncrementValue() {
    c.Value++
    fmt.Printf("Inside method: %d\n", c.Value)
}

// Value receiver - returns new value
func (c Counter) Add(n int) Counter {
    return Counter{Value: c.Value + n}
}

func main() {
    counter := Counter{Value: 10}
    
    fmt.Printf("Original: %d\n", counter.Value)
    
    counter.IncrementValue()
    fmt.Printf("After IncrementValue: %d\n", counter.Value) // Still 10
    
    newCounter := counter.Add(5)
    fmt.Printf("Original after Add: %d\n", counter.Value)    // Still 10
    fmt.Printf("New counter: %d\n", newCounter.Value)       // 15
}
```

### Pointer Receivers

Pointer receivers work with the original value and can modify it.

```go
package main

import "fmt"

type Counter struct {
    Value int
}

// Pointer receiver - can modify the original
func (c *Counter) Increment() {
    c.Value++
}

// Pointer receiver - can modify the original
func (c *Counter) Add(n int) {
    c.Value += n
}

// Pointer receiver - reset to zero
func (c *Counter) Reset() {
    c.Value = 0
}

func main() {
    counter := Counter{Value: 10}
    
    fmt.Printf("Original: %d\n", counter.Value)
    
    counter.Increment()
    fmt.Printf("After Increment: %d\n", counter.Value) // 11
    
    counter.Add(5)
    fmt.Printf("After Add(5): %d\n", counter.Value)   // 16
    
    counter.Reset()
    fmt.Printf("After Reset: %d\n", counter.Value)    // 0
}
```

### When to Use Value vs Pointer Receivers

```go
package main

import "fmt"

type SmallStruct struct {
    X, Y int
}

type LargeStruct struct {
    Data [1000000]int
}

// Use value receiver for small structs that don't need modification
func (s SmallStruct) String() string {
    return fmt.Sprintf("(%d, %d)", s.X, s.Y)
}

// Use pointer receiver for large structs to avoid copying
func (l *LargeStruct) Process() {
    // Process the large struct without copying
    fmt.Println("Processing large struct...")
}

// Use pointer receiver when you need to modify the struct
func (s *SmallStruct) Move(dx, dy int) {
    s.X += dx
    s.Y += dy
}

func main() {
    small := SmallStruct{X: 1, Y: 2}
    fmt.Println("Small struct:", small.String())
    
    small.Move(3, 4)
    fmt.Println("After move:", small.String())
    
    large := &LargeStruct{}
    large.Process()
}
```

## Methods on Different Types

### Methods on Built-in Types (via Type Definitions)

You can define methods on custom types based on built-in types.

```go
package main

import (
    "fmt"
    "strings"
)

// Custom type based on string
type MyString string

// Method on custom string type
func (s MyString) Reverse() MyString {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return MyString(runes)
}

// Method to count words
func (s MyString) WordCount() int {
    if len(strings.TrimSpace(string(s))) == 0 {
        return 0
    }
    return len(strings.Fields(string(s)))
}

// Method to convert to uppercase
func (s MyString) Upper() MyString {
    return MyString(strings.ToUpper(string(s)))
}

// Custom type based on slice
type IntSlice []int

// Method to find sum
func (is IntSlice) Sum() int {
    total := 0
    for _, v := range is {
        total += v
    }
    return total
}

// Method to find average
func (is IntSlice) Average() float64 {
    if len(is) == 0 {
        return 0
    }
    return float64(is.Sum()) / float64(len(is))
}

func main() {
    // String methods
    text := MyString("Hello World")
    fmt.Printf("Original: %s\n", text)
    fmt.Printf("Reversed: %s\n", text.Reverse())
    fmt.Printf("Uppercase: %s\n", text.Upper())
    fmt.Printf("Word count: %d\n", text.WordCount())
    
    // Slice methods
    numbers := IntSlice{1, 2, 3, 4, 5}
    fmt.Printf("Numbers: %v\n", numbers)
    fmt.Printf("Sum: %d\n", numbers.Sum())
    fmt.Printf("Average: %.2f\n", numbers.Average())
}
```

### Methods on Struct Types

```go
package main

import (
    "fmt"
    "math"
)

type Rectangle struct {
    Width, Height float64
}

type Triangle struct {
    Base, Height float64
}

// Rectangle methods
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func (r Rectangle) IsSquare() bool {
    return r.Width == r.Height
}

// Triangle methods
func (t Triangle) Area() float64 {
    return 0.5 * t.Base * t.Height
}

func (t Triangle) Hypotenuse() float64 {
    return math.Sqrt(t.Base*t.Base + t.Height*t.Height)
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    
    fmt.Printf("Rectangle: %.2f x %.2f\n", rect.Width, rect.Height)
    fmt.Printf("Area: %.2f\n", rect.Area())
    fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
    fmt.Printf("Is square: %t\n", rect.IsSquare())
    
    rect.Scale(2)
    fmt.Printf("After scaling: %.2f x %.2f\n", rect.Width, rect.Height)
    fmt.Printf("Is square now: %t\n", rect.IsSquare())
    
    triangle := Triangle{Base: 3, Height: 4}
    fmt.Printf("\nTriangle: base=%.2f, height=%.2f\n", triangle.Base, triangle.Height)
    fmt.Printf("Area: %.2f\n", triangle.Area())
    fmt.Printf("Hypotenuse: %.2f\n", triangle.Hypotenuse())
}
```

## Method Sets and Interface Satisfaction

### Understanding Method Sets

```go
package main

import "fmt"

type Shape interface {
    Area() float64
    Perimeter() float64
}

type Scalable interface {
    Scale(factor float64)
}

type Rectangle struct {
    Width, Height float64
}

// Value receiver methods
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Pointer receiver method
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f, Perimeter: %.2f\n", s.Area(), s.Perimeter())
}

func scaleShape(s Scalable, factor float64) {
    s.Scale(factor)
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    rectPtr := &Rectangle{Width: 10, Height: 5}
    
    // Value satisfies Shape interface (value receiver methods)
    printShapeInfo(rect)
    printShapeInfo(rectPtr) // Pointer also works
    
    // Only pointer satisfies Scalable interface (pointer receiver method)
    // scaleShape(rect, 2)    // This would cause a compile error
    scaleShape(rectPtr, 2)    // This works
    
    printShapeInfo(rectPtr)
}
```

## Advanced Method Concepts

### Method Expressions and Method Values

```go
package main

import "fmt"

type Calculator struct {
    Value float64
}

func (c *Calculator) Add(x float64) {
    c.Value += x
}

func (c *Calculator) Multiply(x float64) {
    c.Value *= x
}

func (c Calculator) GetValue() float64 {
    return c.Value
}

func main() {
    calc := &Calculator{Value: 10}
    
    // Method value - method bound to specific instance
    addMethod := calc.Add
    addMethod(5)
    fmt.Printf("After add: %.2f\n", calc.Value)
    
    // Method expression - method as function
    multiplyFunc := (*Calculator).Multiply
    multiplyFunc(calc, 2)
    fmt.Printf("After multiply: %.2f\n", calc.Value)
    
    // Method value for value receiver
    getValue := calc.GetValue
    fmt.Printf("Current value: %.2f\n", getValue())
}
```

### Methods with Multiple Return Values

```go
package main

import (
    "errors"
    "fmt"
)

type BankAccount struct {
    Balance float64
    Owner   string
}

func (ba *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return errors.New("deposit amount must be positive")
    }
    ba.Balance += amount
    return nil
}

func (ba *BankAccount) Withdraw(amount float64) (float64, error) {
    if amount <= 0 {
        return 0, errors.New("withdrawal amount must be positive")
    }
    if amount > ba.Balance {
        return 0, errors.New("insufficient funds")
    }
    ba.Balance -= amount
    return amount, nil
}

func (ba BankAccount) GetBalance() (float64, string) {
    return ba.Balance, ba.Owner
}

func (ba *BankAccount) Transfer(to *BankAccount, amount float64) error {
    withdrawn, err := ba.Withdraw(amount)
    if err != nil {
        return err
    }
    
    return to.Deposit(withdrawn)
}

func main() {
    account1 := &BankAccount{Balance: 1000, Owner: "Alice"}
    account2 := &BankAccount{Balance: 500, Owner: "Bob"}
    
    // Deposit
    if err := account1.Deposit(200); err != nil {
        fmt.Printf("Deposit error: %v\n", err)
    } else {
        fmt.Printf("Deposit successful\n")
    }
    
    // Check balance
    balance, owner := account1.GetBalance()
    fmt.Printf("%s's balance: $%.2f\n", owner, balance)
    
    // Withdraw
    if amount, err := account1.Withdraw(300); err != nil {
        fmt.Printf("Withdrawal error: %v\n", err)
    } else {
        fmt.Printf("Withdrew: $%.2f\n", amount)
    }
    
    // Transfer
    if err := account1.Transfer(account2, 200); err != nil {
        fmt.Printf("Transfer error: %v\n", err)
    } else {
        fmt.Printf("Transfer successful\n")
    }
    
    // Final balances
    balance1, _ := account1.GetBalance()
    balance2, _ := account2.GetBalance()
    fmt.Printf("Final balances - Alice: $%.2f, Bob: $%.2f\n", balance1, balance2)
}
```

## Practical Examples

### Stack Implementation

```go
package main

import (
    "errors"
    "fmt"
)

type Stack struct {
    items []int
}

// Push adds an item to the top of the stack
func (s *Stack) Push(item int) {
    s.items = append(s.items, item)
}

// Pop removes and returns the top item from the stack
func (s *Stack) Pop() (int, error) {
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    
    index := len(s.items) - 1
    item := s.items[index]
    s.items = s.items[:index]
    return item, nil
}

// Peek returns the top item without removing it
func (s Stack) Peek() (int, error) {
    if len(s.items) == 0 {
        return 0, errors.New("stack is empty")
    }
    return s.items[len(s.items)-1], nil
}

// IsEmpty checks if the stack is empty
func (s Stack) IsEmpty() bool {
    return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s Stack) Size() int {
    return len(s.items)
}

// Clear removes all items from the stack
func (s *Stack) Clear() {
    s.items = s.items[:0]
}

// String returns a string representation of the stack
func (s Stack) String() string {
    return fmt.Sprintf("Stack%v", s.items)
}

func main() {
    stack := &Stack{}
    
    // Push items
    stack.Push(10)
    stack.Push(20)
    stack.Push(30)
    
    fmt.Printf("Stack: %s\n", stack)
    fmt.Printf("Size: %d\n", stack.Size())
    
    // Peek
    if top, err := stack.Peek(); err == nil {
        fmt.Printf("Top item: %d\n", top)
    }
    
    // Pop items
    for !stack.IsEmpty() {
        if item, err := stack.Pop(); err == nil {
            fmt.Printf("Popped: %d\n", item)
        }
    }
    
    fmt.Printf("Stack after popping all: %s\n", stack)
    fmt.Printf("Is empty: %t\n", stack.IsEmpty())
}
```

### Queue Implementation

```go
package main

import (
    "errors"
    "fmt"
)

type Queue struct {
    items []string
}

// Enqueue adds an item to the rear of the queue
func (q *Queue) Enqueue(item string) {
    q.items = append(q.items, item)
}

// Dequeue removes and returns the front item from the queue
func (q *Queue) Dequeue() (string, error) {
    if len(q.items) == 0 {
        return "", errors.New("queue is empty")
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item, nil
}

// Front returns the front item without removing it
func (q Queue) Front() (string, error) {
    if len(q.items) == 0 {
        return "", errors.New("queue is empty")
    }
    return q.items[0], nil
}

// Rear returns the rear item without removing it
func (q Queue) Rear() (string, error) {
    if len(q.items) == 0 {
        return "", errors.New("queue is empty")
    }
    return q.items[len(q.items)-1], nil
}

// IsEmpty checks if the queue is empty
func (q Queue) IsEmpty() bool {
    return len(q.items) == 0
}

// Size returns the number of items in the queue
func (q Queue) Size() int {
    return len(q.items)
}

// Clear removes all items from the queue
func (q *Queue) Clear() {
    q.items = q.items[:0]
}

// String returns a string representation of the queue
func (q Queue) String() string {
    return fmt.Sprintf("Queue%v", q.items)
}

func main() {
    queue := &Queue{}
    
    // Enqueue items
    queue.Enqueue("First")
    queue.Enqueue("Second")
    queue.Enqueue("Third")
    
    fmt.Printf("Queue: %s\n", queue)
    fmt.Printf("Size: %d\n", queue.Size())
    
    // Check front and rear
    if front, err := queue.Front(); err == nil {
        fmt.Printf("Front: %s\n", front)
    }
    
    if rear, err := queue.Rear(); err == nil {
        fmt.Printf("Rear: %s\n", rear)
    }
    
    // Dequeue items
    for !queue.IsEmpty() {
        if item, err := queue.Dequeue(); err == nil {
            fmt.Printf("Dequeued: %s\n", item)
        }
    }
    
    fmt.Printf("Queue after dequeuing all: %s\n", queue)
    fmt.Printf("Is empty: %t\n", queue.IsEmpty())
}
```

### HTTP Client with Methods

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type APIClient struct {
    BaseURL    string
    HTTPClient *http.Client
    Headers    map[string]string
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL string) *APIClient {
    return &APIClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        Headers: make(map[string]string),
    }
}

// SetHeader sets a header for all requests
func (c *APIClient) SetHeader(key, value string) {
    c.Headers[key] = value
}

// SetTimeout sets the HTTP client timeout
func (c *APIClient) SetTimeout(timeout time.Duration) {
    c.HTTPClient.Timeout = timeout
}

// makeRequest is a helper method for making HTTP requests
func (c *APIClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
    var reqBody io.Reader
    
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }
    
    req, err := http.NewRequest(method, c.BaseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }
    
    // Set headers
    for key, value := range c.Headers {
        req.Header.Set(key, value)
    }
    
    if body != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    
    return c.HTTPClient.Do(req)
}

// GetUser retrieves a user by ID
func (c *APIClient) GetUser(id int) (*User, error) {
    endpoint := fmt.Sprintf("/users/%d", id)
    resp, err := c.makeRequest("GET", endpoint, nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }
    
    return &user, nil
}

// CreateUser creates a new user
func (c *APIClient) CreateUser(user User) (*User, error) {
    resp, err := c.makeRequest("POST", "/users", user)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var createdUser User
    if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
        return nil, err
    }
    
    return &createdUser, nil
}

// UpdateUser updates an existing user
func (c *APIClient) UpdateUser(user User) (*User, error) {
    endpoint := fmt.Sprintf("/users/%d", user.ID)
    resp, err := c.makeRequest("PUT", endpoint, user)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API error: %s", resp.Status)
    }
    
    var updatedUser User
    if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
        return nil, err
    }
    
    return &updatedUser, nil
}

// DeleteUser deletes a user by ID
func (c *APIClient) DeleteUser(id int) error {
    endpoint := fmt.Sprintf("/users/%d", id)
    resp, err := c.makeRequest("DELETE", endpoint, nil)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusNoContent {
        return fmt.Errorf("API error: %s", resp.Status)
    }
    
    return nil
}

func main() {
    // This is a demonstration - the API endpoints don't actually exist
    client := NewAPIClient("https://jsonplaceholder.typicode.com")
    client.SetHeader("Authorization", "Bearer token123")
    client.SetTimeout(10 * time.Second)
    
    // Example usage (these would work with a real API)
    fmt.Println("API Client created with methods for:")
    fmt.Println("- GetUser(id)")
    fmt.Println("- CreateUser(user)")
    fmt.Println("- UpdateUser(user)")
    fmt.Println("- DeleteUser(id)")
    fmt.Println("- SetHeader(key, value)")
    fmt.Println("- SetTimeout(duration)")
    
    // In a real scenario, you would use:
    // user, err := client.GetUser(1)
    // newUser, err := client.CreateUser(User{Name: "John", Email: "john@example.com"})
    // etc.
}
```

## Method Chaining

```go
package main

import "fmt"

type StringBuilder struct {
    content string
}

// NewStringBuilder creates a new StringBuilder
func NewStringBuilder() *StringBuilder {
    return &StringBuilder{}
}

// Append adds text to the builder
func (sb *StringBuilder) Append(text string) *StringBuilder {
    sb.content += text
    return sb
}

// AppendLine adds text with a newline
func (sb *StringBuilder) AppendLine(text string) *StringBuilder {
    sb.content += text + "\n"
    return sb
}

// Prepend adds text to the beginning
func (sb *StringBuilder) Prepend(text string) *StringBuilder {
    sb.content = text + sb.content
    return sb
}

// Clear removes all content
func (sb *StringBuilder) Clear() *StringBuilder {
    sb.content = ""
    return sb
}

// String returns the built string
func (sb *StringBuilder) String() string {
    return sb.content
}

// Length returns the length of the content
func (sb *StringBuilder) Length() int {
    return len(sb.content)
}

func main() {
    // Method chaining example
    result := NewStringBuilder().
        Append("Hello").
        Append(" ").
        Append("World").
        AppendLine("!").
        Append("This is ").
        Append("method chaining").
        String()
    
    fmt.Println("Result:")
    fmt.Println(result)
    
    // Another example
    builder := NewStringBuilder()
    builder.AppendLine("Line 1").
        AppendLine("Line 2").
        Prepend("Header: ").
        Append("Footer")
    
    fmt.Printf("Length: %d\n", builder.Length())
    fmt.Println("Content:")
    fmt.Println(builder.String())
}
```

## Best Practices

### 1. Choose Receiver Type Carefully

```go
// Use pointer receivers when:
// - The method needs to modify the receiver
// - The receiver is a large struct
// - You want to avoid copying

type LargeStruct struct {
    data [1000000]int
}

func (ls *LargeStruct) Process() { // Pointer receiver for large struct
    // Process without copying
}

func (ls *LargeStruct) Modify() { // Pointer receiver for modification
    ls.data[0] = 42
}

// Use value receivers when:
// - The method doesn't modify the receiver
// - The receiver is a small struct or basic type
// - You want to ensure immutability

type Point struct {
    X, Y int
}

func (p Point) String() string { // Value receiver for small, immutable operation
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}
```

### 2. Be Consistent with Receiver Types

```go
// Good: All methods use pointer receivers
type Counter struct {
    value int
}

func (c *Counter) Increment() {
    c.value++
}

func (c *Counter) Value() int {
    return c.value
}

func (c *Counter) Reset() {
    c.value = 0
}
```

### 3. Use Meaningful Method Names

```go
// Good method names
func (u *User) IsActive() bool
func (u *User) Activate()
func (u *User) GetFullName() string
func (u *User) UpdateEmail(email string) error

// Avoid generic names like:
// func (u *User) Do()
// func (u *User) Process()
```

### 4. Handle Errors Appropriately

```go
func (ba *BankAccount) Withdraw(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    if amount > ba.balance {
        return errors.New("insufficient funds")
    }
    ba.balance -= amount
    return nil
}
```

## Exercises

### Exercise 1: Temperature Converter

Create a `Temperature` type and methods to convert between Celsius, Fahrenheit, and Kelvin.

```go
type Temperature struct {
    celsius float64
}

// Implement methods:
// - NewTemperatureFromCelsius(c float64) Temperature
// - NewTemperatureFromFahrenheit(f float64) Temperature
// - NewTemperatureFromKelvin(k float64) Temperature
// - Celsius() float64
// - Fahrenheit() float64
// - Kelvin() float64
// - String() string
```

### Exercise 2: Shopping Cart

Create a shopping cart with methods to add items, remove items, calculate total, and apply discounts.

```go
type Item struct {
    Name     string
    Price    float64
    Quantity int
}

type ShoppingCart struct {
    items []Item
}

// Implement methods:
// - AddItem(item Item)
// - RemoveItem(name string)
// - UpdateQuantity(name string, quantity int)
// - Total() float64
// - ApplyDiscount(percentage float64)
// - ItemCount() int
// - Clear()
```

### Exercise 3: Matrix Operations

Create a `Matrix` type with methods for basic matrix operations.

```go
type Matrix struct {
    data [][]float64
    rows, cols int
}

// Implement methods:
// - NewMatrix(rows, cols int) *Matrix
// - Set(row, col int, value float64)
// - Get(row, col int) float64
// - Add(other *Matrix) (*Matrix, error)
// - Multiply(other *Matrix) (*Matrix, error)
// - Transpose() *Matrix
// - String() string
```

## Key Takeaways

1. **Methods add behavior to types**: They allow you to define operations that can be performed on your custom types.

2. **Receiver types matter**: Use pointer receivers for modification or large structs, value receivers for small, immutable operations.

3. **Method sets determine interface satisfaction**: The receiver type affects which interfaces a type can satisfy.

4. **Methods can be defined on any named type**: Not just structs, but also custom types based on built-in types.

5. **Consistency is important**: Be consistent with receiver types within a type's method set.

6. **Methods enable object-oriented patterns**: While Go isn't traditionally OOP, methods provide similar functionality.

7. **Method chaining is possible**: Return the receiver to enable fluent interfaces.

8. **Methods can have multiple return values**: Use this for error handling and multiple results.

9. **Method expressions and values**: Methods can be treated as functions in certain contexts.

10. **Good naming improves readability**: Use clear, descriptive method names that indicate their purpose.

## Next Steps

Now that you understand methods, let's learn about [Interfaces](14-interfaces.md) to see how methods enable polymorphism and abstraction in Go!

---

**Previous**: [← Pointers](12-pointers.md) | **Next**: [Interfaces →](14-interfaces.md)