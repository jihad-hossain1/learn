# 11. Structs

Structs are user-defined types that group together related data. They are similar to classes in other languages but without inheritance. Structs are the foundation for creating complex data types in Go.

## Struct Basics

### Struct Declaration

```go
// Basic struct declaration
type StructName struct {
    field1 dataType1
    field2 dataType2
    // ...
}

// Example
type Person struct {
    Name string
    Age  int
    City string
}
```

### Struct Initialization

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    // Method 1: Zero value initialization
    var p1 Person
    fmt.Printf("Zero value: %+v\n", p1)
    
    // Method 2: Field assignment after declaration
    var p2 Person
    p2.Name = "Alice"
    p2.Age = 25
    p2.City = "New York"
    fmt.Printf("Field assignment: %+v\n", p2)
    
    // Method 3: Struct literal with field names
    p3 := Person{
        Name: "Bob",
        Age:  30,
        City: "San Francisco",
    }
    fmt.Printf("Named fields: %+v\n", p3)
    
    // Method 4: Struct literal with positional values
    p4 := Person{"Charlie", 35, "Chicago"}
    fmt.Printf("Positional: %+v\n", p4)
    
    // Method 5: Partial initialization
    p5 := Person{
        Name: "Diana",
        Age:  28,
        // City will be zero value (empty string)
    }
    fmt.Printf("Partial: %+v\n", p5)
    
    // Method 6: Using new keyword
    p6 := new(Person)
    p6.Name = "Eve"
    p6.Age = 32
    p6.City = "Boston"
    fmt.Printf("Using new: %+v\n", *p6)
}
```

### Accessing and Modifying Struct Fields

```go
package main

import "fmt"

type Student struct {
    ID     int
    Name   string
    Grade  float64
    Active bool
}

func main() {
    student := Student{
        ID:     101,
        Name:   "John Doe",
        Grade:  85.5,
        Active: true,
    }
    
    // Access fields
    fmt.Printf("Student ID: %d\n", student.ID)
    fmt.Printf("Student Name: %s\n", student.Name)
    fmt.Printf("Student Grade: %.2f\n", student.Grade)
    fmt.Printf("Is Active: %t\n", student.Active)
    
    // Modify fields
    student.Grade = 90.0
    student.Active = false
    
    fmt.Printf("\nAfter modification:\n")
    fmt.Printf("Student Grade: %.2f\n", student.Grade)
    fmt.Printf("Is Active: %t\n", student.Active)
    
    // Print entire struct
    fmt.Printf("\nComplete student: %+v\n", student)
}
```

### Anonymous Structs

```go
package main

import "fmt"

func main() {
    // Anonymous struct declaration and initialization
    config := struct {
        Host     string
        Port     int
        Database string
        SSL      bool
    }{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        SSL:      true,
    }
    
    fmt.Printf("Config: %+v\n", config)
    
    // Anonymous struct in slice
    users := []struct {
        Name string
        Role string
    }{
        {"Alice", "Admin"},
        {"Bob", "User"},
        {"Charlie", "Moderator"},
    }
    
    fmt.Println("\nUsers:")
    for _, user := range users {
        fmt.Printf("Name: %s, Role: %s\n", user.Name, user.Role)
    }
    
    // Anonymous struct as function parameter
    processOrder := func(order struct {
        ID       int
        Product  string
        Quantity int
        Price    float64
    }) {
        total := float64(order.Quantity) * order.Price
        fmt.Printf("Order %d: %d x %s = $%.2f\n", 
            order.ID, order.Quantity, order.Product, total)
    }
    
    processOrder(struct {
        ID       int
        Product  string
        Quantity int
        Price    float64
    }{1, "Laptop", 2, 999.99})
}
```

### Nested Structs

```go
package main

import "fmt"

type Address struct {
    Street  string
    City    string
    State   string
    ZipCode string
    Country string
}

type Contact struct {
    Email string
    Phone string
}

type Person struct {
    Name    string
    Age     int
    Address Address // Nested struct
    Contact Contact // Another nested struct
}

func main() {
    person := Person{
        Name: "John Smith",
        Age:  30,
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            State:   "NY",
            ZipCode: "10001",
            Country: "USA",
        },
        Contact: Contact{
            Email: "john.smith@email.com",
            Phone: "+1-555-0123",
        },
    }
    
    // Access nested fields
    fmt.Printf("Name: %s\n", person.Name)
    fmt.Printf("City: %s\n", person.Address.City)
    fmt.Printf("Email: %s\n", person.Contact.Email)
    
    // Modify nested fields
    person.Address.City = "Boston"
    person.Contact.Phone = "+1-555-0456"
    
    fmt.Printf("\nAfter modification:\n")
    fmt.Printf("City: %s\n", person.Address.City)
    fmt.Printf("Phone: %s\n", person.Contact.Phone)
    
    // Print complete struct
    fmt.Printf("\nComplete person: %+v\n", person)
}
```

### Struct Embedding (Anonymous Fields)

```go
package main

import "fmt"

type Address struct {
    Street  string
    City    string
    ZipCode string
}

type Contact struct {
    Email string
    Phone string
}

// Struct with embedded structs
type Employee struct {
    ID       int
    Name     string
    Position string
    Salary   float64
    Address  // Embedded struct (anonymous field)
    Contact  // Embedded struct (anonymous field)
}

func main() {
    emp := Employee{
        ID:       101,
        Name:     "Alice Johnson",
        Position: "Software Engineer",
        Salary:   75000,
        Address: Address{
            Street:  "456 Tech Ave",
            City:    "San Francisco",
            ZipCode: "94105",
        },
        Contact: Contact{
            Email: "alice.johnson@company.com",
            Phone: "+1-555-0789",
        },
    }
    
    // Access embedded fields directly
    fmt.Printf("Name: %s\n", emp.Name)
    fmt.Printf("City: %s\n", emp.City)    // Direct access to embedded field
    fmt.Printf("Email: %s\n", emp.Email)  // Direct access to embedded field
    
    // Can also access through embedded struct name
    fmt.Printf("Street: %s\n", emp.Address.Street)
    fmt.Printf("Phone: %s\n", emp.Contact.Phone)
    
    // Modify embedded fields
    emp.City = "Seattle"        // Direct modification
    emp.Contact.Email = "alice.j@company.com" // Through embedded struct
    
    fmt.Printf("\nAfter modification:\n")
    fmt.Printf("City: %s\n", emp.City)
    fmt.Printf("Email: %s\n", emp.Email)
}
```

### Struct Methods

```go
package main

import (
    "fmt"
    "math"
)

type Rectangle struct {
    Width  float64
    Height float64
}

type Circle struct {
    Radius float64
}

// Method with value receiver
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

// Method with value receiver
func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Method with pointer receiver (can modify the struct)
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor
    r.Height *= factor
}

// Method with value receiver
func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

// Method with value receiver
func (c Circle) Circumference() float64 {
    return 2 * math.Pi * c.Radius
}

// Method with pointer receiver
func (c *Circle) Scale(factor float64) {
    c.Radius *= factor
}

// Method that returns a new struct
func (r Rectangle) Scaled(factor float64) Rectangle {
    return Rectangle{
        Width:  r.Width * factor,
        Height: r.Height * factor,
    }
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 3}
    
    // Call methods
    fmt.Printf("Rectangle: %+v\n", rect)
    fmt.Printf("Area: %.2f\n", rect.Area())
    fmt.Printf("Perimeter: %.2f\n", rect.Perimeter())
    
    fmt.Printf("\nCircle: %+v\n", circle)
    fmt.Printf("Area: %.2f\n", circle.Area())
    fmt.Printf("Circumference: %.2f\n", circle.Circumference())
    
    // Modify using pointer receiver methods
    fmt.Println("\nScaling rectangle by 2...")
    rect.Scale(2)
    fmt.Printf("Scaled rectangle: %+v\n", rect)
    fmt.Printf("New area: %.2f\n", rect.Area())
    
    fmt.Println("\nScaling circle by 1.5...")
    circle.Scale(1.5)
    fmt.Printf("Scaled circle: %+v\n", circle)
    fmt.Printf("New area: %.2f\n", circle.Area())
    
    // Create new scaled struct
    newRect := Rectangle{Width: 4, Height: 6}
    scaledRect := newRect.Scaled(3)
    fmt.Printf("\nOriginal: %+v\n", newRect)
    fmt.Printf("Scaled copy: %+v\n", scaledRect)
}
```

### Struct Tags

```go
package main

import (
    "encoding/json"
    "fmt"
    "reflect"
)

type User struct {
    ID       int    `json:"id" db:"user_id" validate:"required"`
    Name     string `json:"name" db:"full_name" validate:"required,min=2"`
    Email    string `json:"email" db:"email_address" validate:"required,email"`
    Age      int    `json:"age,omitempty" db:"age" validate:"min=0,max=120"`
    Password string `json:"-" db:"password_hash" validate:"required,min=8"`
    IsActive bool   `json:"is_active" db:"active" validate:""`
}

func main() {
    user := User{
        ID:       1,
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      30,
        Password: "secretpassword",
        IsActive: true,
    }
    
    // JSON marshaling (uses json tags)
    jsonData, err := json.Marshal(user)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("JSON: %s\n", jsonData)
    
    // JSON unmarshaling
    jsonStr := `{"id":2,"name":"Jane Smith","email":"jane@example.com","is_active":false}`
    var newUser User
    err = json.Unmarshal([]byte(jsonStr), &newUser)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    fmt.Printf("Unmarshaled: %+v\n", newUser)
    
    // Reading struct tags using reflection
    fmt.Println("\nStruct tags:")
    userType := reflect.TypeOf(user)
    for i := 0; i < userType.NumField(); i++ {
        field := userType.Field(i)
        jsonTag := field.Tag.Get("json")
        dbTag := field.Tag.Get("db")
        validateTag := field.Tag.Get("validate")
        
        fmt.Printf("Field: %s\n", field.Name)
        fmt.Printf("  JSON tag: %s\n", jsonTag)
        fmt.Printf("  DB tag: %s\n", dbTag)
        fmt.Printf("  Validate tag: %s\n", validateTag)
        fmt.Println()
    }
}
```

### Struct Comparison

```go
package main

import "fmt"

type Point struct {
    X, Y int
}

type Person struct {
    Name string
    Age  int
}

type PersonWithSlice struct {
    Name    string
    Age     int
    Hobbies []string // Slice makes struct non-comparable
}

func main() {
    // Comparable structs
    p1 := Point{X: 1, Y: 2}
    p2 := Point{X: 1, Y: 2}
    p3 := Point{X: 2, Y: 3}
    
    fmt.Printf("p1 == p2: %t\n", p1 == p2) // true
    fmt.Printf("p1 == p3: %t\n", p1 == p3) // false
    
    person1 := Person{Name: "Alice", Age: 25}
    person2 := Person{Name: "Alice", Age: 25}
    person3 := Person{Name: "Bob", Age: 30}
    
    fmt.Printf("person1 == person2: %t\n", person1 == person2) // true
    fmt.Printf("person1 == person3: %t\n", person1 == person3) // false
    
    // Non-comparable structs (contain slices, maps, or functions)
    // The following would cause a compile error:
    // personA := PersonWithSlice{Name: "Alice", Age: 25, Hobbies: []string{"reading"}}
    // personB := PersonWithSlice{Name: "Alice", Age: 25, Hobbies: []string{"reading"}}
    // fmt.Printf("personA == personB: %t\n", personA == personB) // Compile error!
    
    // Use custom comparison function for non-comparable structs
    personA := PersonWithSlice{Name: "Alice", Age: 25, Hobbies: []string{"reading", "swimming"}}
    personB := PersonWithSlice{Name: "Alice", Age: 25, Hobbies: []string{"reading", "swimming"}}
    
    fmt.Printf("Custom comparison: %t\n", comparePersonWithSlice(personA, personB))
}

// Custom comparison function for non-comparable structs
func comparePersonWithSlice(a, b PersonWithSlice) bool {
    if a.Name != b.Name || a.Age != b.Age {
        return false
    }
    
    if len(a.Hobbies) != len(b.Hobbies) {
        return false
    }
    
    for i, hobby := range a.Hobbies {
        if hobby != b.Hobbies[i] {
            return false
        }
    }
    
    return true
}
```

## Practical Examples

### Example 1: Employee Management System

```go
package main

import (
    "fmt"
    "time"
)

type Department struct {
    ID   int
    Name string
}

type Employee struct {
    ID         int
    FirstName  string
    LastName   string
    Email      string
    Department Department
    Salary     float64
    HireDate   time.Time
    IsActive   bool
}

// Method to get full name
func (e Employee) FullName() string {
    return e.FirstName + " " + e.LastName
}

// Method to calculate years of service
func (e Employee) YearsOfService() int {
    return int(time.Since(e.HireDate).Hours() / 24 / 365)
}

// Method to give raise
func (e *Employee) GiveRaise(percentage float64) {
    e.Salary *= (1 + percentage/100)
}

// Method to check if employee is senior (5+ years)
func (e Employee) IsSenior() bool {
    return e.YearsOfService() >= 5
}

// Method to get employee summary
func (e Employee) Summary() string {
    return fmt.Sprintf("%s (%s) - %s Department, $%.2f, %d years",
        e.FullName(), e.Email, e.Department.Name, e.Salary, e.YearsOfService())
}

type Company struct {
    Name      string
    Employees []Employee
}

// Method to add employee
func (c *Company) AddEmployee(emp Employee) {
    c.Employees = append(c.Employees, emp)
}

// Method to get employees by department
func (c Company) GetEmployeesByDepartment(deptName string) []Employee {
    var result []Employee
    for _, emp := range c.Employees {
        if emp.Department.Name == deptName {
            result = append(result, emp)
        }
    }
    return result
}

// Method to get average salary
func (c Company) AverageSalary() float64 {
    if len(c.Employees) == 0 {
        return 0
    }
    
    total := 0.0
    for _, emp := range c.Employees {
        total += emp.Salary
    }
    return total / float64(len(c.Employees))
}

// Method to get senior employees
func (c Company) GetSeniorEmployees() []Employee {
    var seniors []Employee
    for _, emp := range c.Employees {
        if emp.IsSenior() {
            seniors = append(seniors, emp)
        }
    }
    return seniors
}

func main() {
    // Create departments
    engineering := Department{ID: 1, Name: "Engineering"}
    marketing := Department{ID: 2, Name: "Marketing"}
    hr := Department{ID: 3, Name: "Human Resources"}
    
    // Create company
    company := Company{Name: "TechCorp"}
    
    // Add employees
    company.AddEmployee(Employee{
        ID:         101,
        FirstName:  "John",
        LastName:   "Doe",
        Email:      "john.doe@techcorp.com",
        Department: engineering,
        Salary:     85000,
        HireDate:   time.Date(2018, 3, 15, 0, 0, 0, 0, time.UTC),
        IsActive:   true,
    })
    
    company.AddEmployee(Employee{
        ID:         102,
        FirstName:  "Jane",
        LastName:   "Smith",
        Email:      "jane.smith@techcorp.com",
        Department: marketing,
        Salary:     75000,
        HireDate:   time.Date(2019, 7, 22, 0, 0, 0, 0, time.UTC),
        IsActive:   true,
    })
    
    company.AddEmployee(Employee{
        ID:         103,
        FirstName:  "Bob",
        LastName:   "Johnson",
        Email:      "bob.johnson@techcorp.com",
        Department: engineering,
        Salary:     95000,
        HireDate:   time.Date(2016, 1, 10, 0, 0, 0, 0, time.UTC),
        IsActive:   true,
    })
    
    company.AddEmployee(Employee{
        ID:         104,
        FirstName:  "Alice",
        LastName:   "Brown",
        Email:      "alice.brown@techcorp.com",
        Department: hr,
        Salary:     65000,
        HireDate:   time.Date(2020, 9, 5, 0, 0, 0, 0, time.UTC),
        IsActive:   true,
    })
    
    // Display company information
    fmt.Printf("Company: %s\n", company.Name)
    fmt.Printf("Total Employees: %d\n", len(company.Employees))
    fmt.Printf("Average Salary: $%.2f\n\n", company.AverageSalary())
    
    // Display all employees
    fmt.Println("All Employees:")
    for _, emp := range company.Employees {
        fmt.Printf("- %s\n", emp.Summary())
    }
    
    // Display engineering employees
    fmt.Println("\nEngineering Department:")
    engEmployees := company.GetEmployeesByDepartment("Engineering")
    for _, emp := range engEmployees {
        fmt.Printf("- %s\n", emp.Summary())
    }
    
    // Display senior employees
    fmt.Println("\nSenior Employees (5+ years):")
    seniors := company.GetSeniorEmployees()
    for _, emp := range seniors {
        fmt.Printf("- %s\n", emp.Summary())
    }
    
    // Give raise to first employee
    fmt.Println("\nGiving 10% raise to John Doe...")
    fmt.Printf("Before raise: $%.2f\n", company.Employees[0].Salary)
    company.Employees[0].GiveRaise(10)
    fmt.Printf("After raise: $%.2f\n", company.Employees[0].Salary)
}
```

### Example 2: Library Management System

```go
package main

import (
    "fmt"
    "time"
)

type Author struct {
    ID        int
    FirstName string
    LastName  string
    BirthYear int
}

func (a Author) FullName() string {
    return a.FirstName + " " + a.LastName
}

type Book struct {
    ID          int
    Title       string
    Author      Author
    ISBN        string
    Pages       int
    PublishYear int
    Genre       string
    Available   bool
}

func (b Book) String() string {
    status := "Available"
    if !b.Available {
        status = "Checked Out"
    }
    return fmt.Sprintf("%s by %s (%d) - %s", 
        b.Title, b.Author.FullName(), b.PublishYear, status)
}

type Member struct {
    ID       int
    Name     string
    Email    string
    JoinDate time.Time
    Active   bool
}

func (m Member) MembershipDuration() time.Duration {
    return time.Since(m.JoinDate)
}

type Loan struct {
    ID         int
    Book       Book
    Member     Member
    LoanDate   time.Time
    DueDate    time.Time
    ReturnDate *time.Time
}

func (l Loan) IsOverdue() bool {
    if l.ReturnDate != nil {
        return false // Already returned
    }
    return time.Now().After(l.DueDate)
}

func (l Loan) DaysOverdue() int {
    if !l.IsOverdue() {
        return 0
    }
    return int(time.Since(l.DueDate).Hours() / 24)
}

type Library struct {
    Name    string
    Books   []Book
    Members []Member
    Loans   []Loan
    nextID  int
}

func NewLibrary(name string) *Library {
    return &Library{
        Name:   name,
        nextID: 1,
    }
}

func (lib *Library) AddBook(book Book) {
    book.ID = lib.nextID
    lib.nextID++
    book.Available = true
    lib.Books = append(lib.Books, book)
}

func (lib *Library) AddMember(member Member) {
    member.ID = lib.nextID
    lib.nextID++
    member.Active = true
    lib.Members = append(lib.Members, member)
}

func (lib *Library) FindBookByID(id int) (*Book, bool) {
    for i := range lib.Books {
        if lib.Books[i].ID == id {
            return &lib.Books[i], true
        }
    }
    return nil, false
}

func (lib *Library) FindMemberByID(id int) (*Member, bool) {
    for i := range lib.Members {
        if lib.Members[i].ID == id {
            return &lib.Members[i], true
        }
    }
    return nil, false
}

func (lib *Library) CheckOutBook(bookID, memberID int) error {
    book, bookExists := lib.FindBookByID(bookID)
    if !bookExists {
        return fmt.Errorf("book with ID %d not found", bookID)
    }
    
    if !book.Available {
        return fmt.Errorf("book '%s' is not available", book.Title)
    }
    
    member, memberExists := lib.FindMemberByID(memberID)
    if !memberExists {
        return fmt.Errorf("member with ID %d not found", memberID)
    }
    
    if !member.Active {
        return fmt.Errorf("member '%s' is not active", member.Name)
    }
    
    // Create loan
    loan := Loan{
        ID:       lib.nextID,
        Book:     *book,
        Member:   *member,
        LoanDate: time.Now(),
        DueDate:  time.Now().AddDate(0, 0, 14), // 2 weeks
    }
    
    lib.nextID++
    lib.Loans = append(lib.Loans, loan)
    book.Available = false
    
    return nil
}

func (lib *Library) ReturnBook(loanID int) error {
    for i := range lib.Loans {
        if lib.Loans[i].ID == loanID && lib.Loans[i].ReturnDate == nil {
            now := time.Now()
            lib.Loans[i].ReturnDate = &now
            
            // Make book available again
            book, _ := lib.FindBookByID(lib.Loans[i].Book.ID)
            book.Available = true
            
            return nil
        }
    }
    return fmt.Errorf("active loan with ID %d not found", loanID)
}

func (lib *Library) GetOverdueLoans() []Loan {
    var overdue []Loan
    for _, loan := range lib.Loans {
        if loan.IsOverdue() {
            overdue = append(overdue, loan)
        }
    }
    return overdue
}

func (lib *Library) GetBooksByAuthor(authorName string) []Book {
    var books []Book
    for _, book := range lib.Books {
        if book.Author.FullName() == authorName {
            books = append(books, book)
        }
    }
    return books
}

func main() {
    // Create library
    library := NewLibrary("City Public Library")
    
    // Add authors and books
    tolkien := Author{FirstName: "J.R.R.", LastName: "Tolkien", BirthYear: 1892}
    orwell := Author{FirstName: "George", LastName: "Orwell", BirthYear: 1903}
    
    library.AddBook(Book{
        Title:       "The Hobbit",
        Author:      tolkien,
        ISBN:        "978-0547928227",
        Pages:       366,
        PublishYear: 1937,
        Genre:       "Fantasy",
    })
    
    library.AddBook(Book{
        Title:       "1984",
        Author:      orwell,
        ISBN:        "978-0451524935",
        Pages:       328,
        PublishYear: 1949,
        Genre:       "Dystopian Fiction",
    })
    
    library.AddBook(Book{
        Title:       "Animal Farm",
        Author:      orwell,
        ISBN:        "978-0451526342",
        Pages:       112,
        PublishYear: 1945,
        Genre:       "Political Satire",
    })
    
    // Add members
    library.AddMember(Member{
        Name:     "Alice Johnson",
        Email:    "alice@email.com",
        JoinDate: time.Now().AddDate(-1, 0, 0), // Joined 1 year ago
    })
    
    library.AddMember(Member{
        Name:     "Bob Smith",
        Email:    "bob@email.com",
        JoinDate: time.Now().AddDate(0, -6, 0), // Joined 6 months ago
    })
    
    // Display library information
    fmt.Printf("Library: %s\n", library.Name)
    fmt.Printf("Books: %d\n", len(library.Books))
    fmt.Printf("Members: %d\n\n", len(library.Members))
    
    // Display all books
    fmt.Println("Available Books:")
    for _, book := range library.Books {
        fmt.Printf("- %s\n", book.String())
    }
    
    // Check out books
    fmt.Println("\nChecking out books...")
    err := library.CheckOutBook(1, 1) // Alice checks out The Hobbit
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Successfully checked out The Hobbit to Alice")
    }
    
    err = library.CheckOutBook(2, 2) // Bob checks out 1984
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Successfully checked out 1984 to Bob")
    }
    
    // Display books after checkout
    fmt.Println("\nBooks after checkout:")
    for _, book := range library.Books {
        fmt.Printf("- %s\n", book.String())
    }
    
    // Display active loans
    fmt.Println("\nActive Loans:")
    for _, loan := range library.Loans {
        if loan.ReturnDate == nil {
            fmt.Printf("- Loan #%d: %s to %s (Due: %s)\n",
                loan.ID, loan.Book.Title, loan.Member.Name, 
                loan.DueDate.Format("2006-01-02"))
        }
    }
    
    // Return a book
    fmt.Println("\nReturning The Hobbit...")
    err = library.ReturnBook(1)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Println("Successfully returned The Hobbit")
    }
    
    // Display books by author
    fmt.Println("\nBooks by George Orwell:")
    orwellBooks := library.GetBooksByAuthor("George Orwell")
    for _, book := range orwellBooks {
        fmt.Printf("- %s\n", book.String())
    }
}
```

## Exercises

### Exercise 1: Bank Account System
Create a bank account management system using structs.

```go
package main

import (
    "fmt"
    "time"
)

type Account struct {
    // Define fields: AccountNumber, HolderName, Balance, AccountType, CreatedAt
}

type Transaction struct {
    // Define fields: ID, AccountNumber, Type, Amount, Timestamp, Description
}

type Bank struct {
    // Define fields: Name, Accounts, Transactions
}

// Implement these methods:
// (a *Account) Deposit(amount float64) error
// (a *Account) Withdraw(amount float64) error
// (a *Account) GetBalance() float64
// (b *Bank) CreateAccount(holderName, accountType string) string
// (b *Bank) GetAccount(accountNumber string) (*Account, bool)
// (b *Bank) Transfer(fromAccount, toAccount string, amount float64) error
// (b *Bank) GetTransactionHistory(accountNumber string) []Transaction

func main() {
    // Test the bank account system
}
```

### Exercise 2: E-commerce Product Catalog
Create a product catalog system for an e-commerce platform.

```go
package main

import "time"

type Category struct {
    // Define fields: ID, Name, Description
}

type Product struct {
    // Define fields: ID, Name, Description, Price, Category, Stock, CreatedAt, IsActive
}

type Review struct {
    // Define fields: ID, ProductID, CustomerName, Rating, Comment, Date
}

type Catalog struct {
    // Define fields: Products, Categories, Reviews
}

// Implement these methods:
// (p *Product) UpdatePrice(newPrice float64)
// (p *Product) AddStock(quantity int)
// (p *Product) RemoveStock(quantity int) error
// (p *Product) AverageRating() float64
// (c *Catalog) AddProduct(product Product)
// (c *Catalog) GetProductsByCategory(categoryName string) []Product
// (c *Catalog) SearchProducts(keyword string) []Product
// (c *Catalog) GetTopRatedProducts(limit int) []Product
// (c *Catalog) AddReview(review Review)

func main() {
    // Test the product catalog system
}
```

### Exercise 3: School Management System
Create a comprehensive school management system.

```go
package main

import "time"

type Subject struct {
    // Define fields: Code, Name, Credits
}

type Student struct {
    // Define fields: ID, Name, Email, EnrollmentDate, GPA
}

type Teacher struct {
    // Define fields: ID, Name, Email, Department, HireDate
}

type Course struct {
    // Define fields: ID, Subject, Teacher, Students, Schedule, Semester
}

type Grade struct {
    // Define fields: StudentID, CourseID, Score, GradeDate
}

type School struct {
    // Define fields: Name, Students, Teachers, Courses, Grades
}

// Implement these methods:
// (s *Student) CalculateGPA() float64
// (c *Course) EnrollStudent(student Student) error
// (c *Course) GetEnrolledStudents() []Student
// (sch *School) AddGrade(grade Grade)
// (sch *School) GetStudentGrades(studentID int) []Grade
// (sch *School) GetCoursesByTeacher(teacherID int) []Course
// (sch *School) GetTopStudents(limit int) []Student

func main() {
    // Test the school management system
}
```

## Key Takeaways

1. **Structs** group related data together
2. **Zero values** initialize struct fields to their type's zero value
3. **Struct literals** can initialize structs with or without field names
4. **Methods** can be defined on structs with value or pointer receivers
5. **Pointer receivers** allow methods to modify the struct
6. **Embedding** provides composition and field promotion
7. **Struct tags** provide metadata for serialization and validation
8. **Comparable structs** can be compared with == and !=
9. **Anonymous structs** are useful for temporary data structures
10. **Nested structs** organize complex data hierarchically

## Next Steps

Now that you understand structs, let's learn about [Interfaces](12-interfaces.md) for defining contracts and achieving polymorphism!

---

**Previous**: [← Maps](10-maps.md) | **Next**: [Interfaces →](12-interfaces.md)