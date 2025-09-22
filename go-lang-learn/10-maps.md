# 10. Maps

Maps are Go's built-in key-value data structure, similar to hash tables, dictionaries, or associative arrays in other languages. They provide efficient lookup, insertion, and deletion operations.

## Map Basics

### Map Declaration

```go
// Declare map with specific key and value types
var mapName map[keyType]valueType

// Examples
var ages map[string]int
var scores map[int]float64
var flags map[string]bool
```

### Map Initialization

```go
package main

import "fmt"

func main() {
    // Method 1: Using make function
    ages := make(map[string]int)
    ages["Alice"] = 25
    ages["Bob"] = 30
    ages["Charlie"] = 35
    
    // Method 2: Map literal with values
    scores := map[string]int{
        "Math":    95,
        "Science": 87,
        "English": 92,
    }
    
    // Method 3: Empty map literal
    grades := map[string]string{}
    grades["Alice"] = "A"
    grades["Bob"] = "B+"
    
    // Method 4: Short declaration with initialization
    colors := map[int]string{
        1: "red",
        2: "green",
        3: "blue",
    }
    
    fmt.Printf("Ages: %v\n", ages)
    fmt.Printf("Scores: %v\n", scores)
    fmt.Printf("Grades: %v\n", grades)
    fmt.Printf("Colors: %v\n", colors)
}
```

### Basic Map Operations

```go
package main

import "fmt"

func main() {
    // Create a map
    inventory := make(map[string]int)
    
    // Add/Update elements
    inventory["apples"] = 50
    inventory["bananas"] = 30
    inventory["oranges"] = 25
    
    fmt.Printf("Inventory: %v\n", inventory)
    
    // Access elements
    apples := inventory["apples"]
    fmt.Printf("Apples: %d\n", apples)
    
    // Access non-existent key (returns zero value)
    grapes := inventory["grapes"]
    fmt.Printf("Grapes: %d\n", grapes) // 0
    
    // Check if key exists
    value, exists := inventory["bananas"]
    if exists {
        fmt.Printf("Bananas: %d\n", value)
    } else {
        fmt.Println("Bananas not found")
    }
    
    // Check non-existent key
    value, exists = inventory["grapes"]
    if exists {
        fmt.Printf("Grapes: %d\n", value)
    } else {
        fmt.Println("Grapes not found")
    }
    
    // Update existing element
    inventory["apples"] = 45
    fmt.Printf("Updated inventory: %v\n", inventory)
    
    // Delete element
    delete(inventory, "oranges")
    fmt.Printf("After deleting oranges: %v\n", inventory)
    
    // Map length
    fmt.Printf("Number of items: %d\n", len(inventory))
}
```

### Iterating Over Maps

```go
package main

import "fmt"

func main() {
    student := map[string]int{
        "Math":    95,
        "Science": 87,
        "English": 92,
        "History": 88,
        "Art":     94,
    }
    
    // Iterate over key-value pairs
    fmt.Println("Subject scores:")
    for subject, score := range student {
        fmt.Printf("%s: %d\n", subject, score)
    }
    
    // Iterate over keys only
    fmt.Println("\nSubjects:")
    for subject := range student {
        fmt.Printf("%s ", subject)
    }
    fmt.Println()
    
    // Iterate over values only
    fmt.Println("\nScores:")
    for _, score := range student {
        fmt.Printf("%d ", score)
    }
    fmt.Println()
    
    // Calculate average
    total := 0
    count := 0
    for _, score := range student {
        total += score
        count++
    }
    average := float64(total) / float64(count)
    fmt.Printf("\nAverage score: %.2f\n", average)
}
```

### Maps with Different Value Types

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
    City string
}

func main() {
    // Map with struct values
    employees := map[int]Person{
        101: {Name: "Alice", Age: 30, City: "New York"},
        102: {Name: "Bob", Age: 25, City: "San Francisco"},
        103: {Name: "Charlie", Age: 35, City: "Chicago"},
    }
    
    // Map with slice values
    hobbies := map[string][]string{
        "Alice":   {"reading", "swimming", "cooking"},
        "Bob":     {"gaming", "cycling"},
        "Charlie": {"photography", "hiking", "music", "travel"},
    }
    
    // Map with map values (nested maps)
    grades := map[string]map[string]int{
        "Alice": {
            "Math":    95,
            "Science": 87,
            "English": 92,
        },
        "Bob": {
            "Math":    78,
            "Science": 85,
            "English": 88,
        },
    }
    
    fmt.Println("Employees:")
    for id, person := range employees {
        fmt.Printf("ID %d: %+v\n", id, person)
    }
    
    fmt.Println("\nHobbies:")
    for name, hobbyList := range hobbies {
        fmt.Printf("%s: %v\n", name, hobbyList)
    }
    
    fmt.Println("\nGrades:")
    for student, subjects := range grades {
        fmt.Printf("%s:\n", student)
        for subject, grade := range subjects {
            fmt.Printf("  %s: %d\n", subject, grade)
        }
    }
}
```

### Map Functions

```go
package main

import (
    "fmt"
    "sort"
)

// Check if key exists
func keyExists(m map[string]int, key string) bool {
    _, exists := m[key]
    return exists
}

// Get all keys
func getKeys(m map[string]int) []string {
    keys := make([]string, 0, len(m))
    for key := range m {
        keys = append(keys, key)
    }
    return keys
}

// Get all values
func getValues(m map[string]int) []int {
    values := make([]int, 0, len(m))
    for _, value := range m {
        values = append(values, value)
    }
    return values
}

// Copy map
func copyMap(original map[string]int) map[string]int {
    copy := make(map[string]int)
    for key, value := range original {
        copy[key] = value
    }
    return copy
}

// Merge maps
func mergeMaps(map1, map2 map[string]int) map[string]int {
    result := make(map[string]int)
    
    // Copy first map
    for key, value := range map1 {
        result[key] = value
    }
    
    // Add/overwrite with second map
    for key, value := range map2 {
        result[key] = value
    }
    
    return result
}

// Filter map
func filterMap(m map[string]int, predicate func(string, int) bool) map[string]int {
    result := make(map[string]int)
    for key, value := range m {
        if predicate(key, value) {
            result[key] = value
        }
    }
    return result
}

// Get sorted keys
func getSortedKeys(m map[string]int) []string {
    keys := getKeys(m)
    sort.Strings(keys)
    return keys
}

func main() {
    scores := map[string]int{
        "Alice":   95,
        "Bob":     78,
        "Charlie": 87,
        "Diana":   92,
        "Eve":     85,
    }
    
    fmt.Printf("Original scores: %v\n", scores)
    
    // Test key existence
    fmt.Printf("Alice exists: %t\n", keyExists(scores, "Alice"))
    fmt.Printf("Frank exists: %t\n", keyExists(scores, "Frank"))
    
    // Get keys and values
    keys := getKeys(scores)
    values := getValues(scores)
    fmt.Printf("Keys: %v\n", keys)
    fmt.Printf("Values: %v\n", values)
    
    // Copy map
    scoresCopy := copyMap(scores)
    fmt.Printf("Copy: %v\n", scoresCopy)
    
    // Merge maps
    newScores := map[string]int{
        "Frank": 89,
        "Grace": 94,
        "Alice": 98, // This will overwrite Alice's score
    }
    merged := mergeMaps(scores, newScores)
    fmt.Printf("Merged: %v\n", merged)
    
    // Filter map (scores >= 90)
    highScores := filterMap(scores, func(name string, score int) bool {
        return score >= 90
    })
    fmt.Printf("High scores (>=90): %v\n", highScores)
    
    // Sorted keys
    sortedKeys := getSortedKeys(scores)
    fmt.Printf("Sorted keys: %v\n", sortedKeys)
    
    // Print in sorted order
    fmt.Println("\nScores in alphabetical order:")
    for _, key := range sortedKeys {
        fmt.Printf("%s: %d\n", key, scores[key])
    }
}
```

## Practical Examples

### Example 1: Word Frequency Counter

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

// Count word frequencies in text
func countWords(text string) map[string]int {
    // Convert to lowercase and split into words
    words := strings.Fields(strings.ToLower(text))
    
    // Remove punctuation and count
    wordCount := make(map[string]int)
    for _, word := range words {
        // Simple punctuation removal
        word = strings.Trim(word, ".,!?;:()[]{}\"'")
        if word != "" {
            wordCount[word]++
        }
    }
    
    return wordCount
}

// Get most frequent words
func getMostFrequent(wordCount map[string]int, n int) []struct {
    Word  string
    Count int
} {
    type wordFreq struct {
        Word  string
        Count int
    }
    
    // Convert map to slice
    frequencies := make([]wordFreq, 0, len(wordCount))
    for word, count := range wordCount {
        frequencies = append(frequencies, wordFreq{word, count})
    }
    
    // Sort by count (descending)
    sort.Slice(frequencies, func(i, j int) bool {
        return frequencies[i].Count > frequencies[j].Count
    })
    
    // Return top n
    if n > len(frequencies) {
        n = len(frequencies)
    }
    return frequencies[:n]
}

func main() {
    text := `Go is an open source programming language that makes it easy to build simple, reliable, and efficient software. Go was created at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.`
    
    fmt.Println("Text:", text)
    fmt.Println()
    
    // Count words
    wordCount := countWords(text)
    
    fmt.Printf("Total unique words: %d\n", len(wordCount))
    fmt.Println("\nWord frequencies:")
    
    // Print all word counts
    for word, count := range wordCount {
        fmt.Printf("%s: %d\n", word, count)
    }
    
    // Get most frequent words
    fmt.Println("\nTop 5 most frequent words:")
    topWords := getMostFrequent(wordCount, 5)
    for i, wf := range topWords {
        fmt.Printf("%d. %s: %d\n", i+1, wf.Word, wf.Count)
    }
}
```

### Example 2: Student Management System

```go
package main

import "fmt"

type Student struct {
    ID     int
    Name   string
    Age    int
    Grades map[string]float64
}

type StudentManager struct {
    students map[int]*Student
    nextID   int
}

// Create new student manager
func NewStudentManager() *StudentManager {
    return &StudentManager{
        students: make(map[int]*Student),
        nextID:   1,
    }
}

// Add student
func (sm *StudentManager) AddStudent(name string, age int) int {
    student := &Student{
        ID:     sm.nextID,
        Name:   name,
        Age:    age,
        Grades: make(map[string]float64),
    }
    
    sm.students[sm.nextID] = student
    sm.nextID++
    
    return student.ID
}

// Get student by ID
func (sm *StudentManager) GetStudent(id int) (*Student, bool) {
    student, exists := sm.students[id]
    return student, exists
}

// Add grade for student
func (sm *StudentManager) AddGrade(studentID int, subject string, grade float64) error {
    student, exists := sm.students[studentID]
    if !exists {
        return fmt.Errorf("student with ID %d not found", studentID)
    }
    
    student.Grades[subject] = grade
    return nil
}

// Calculate average grade for student
func (sm *StudentManager) GetAverage(studentID int) (float64, error) {
    student, exists := sm.students[studentID]
    if !exists {
        return 0, fmt.Errorf("student with ID %d not found", studentID)
    }
    
    if len(student.Grades) == 0 {
        return 0, fmt.Errorf("no grades found for student %d", studentID)
    }
    
    total := 0.0
    for _, grade := range student.Grades {
        total += grade
    }
    
    return total / float64(len(student.Grades)), nil
}

// Get all students
func (sm *StudentManager) GetAllStudents() map[int]*Student {
    return sm.students
}

// Remove student
func (sm *StudentManager) RemoveStudent(id int) error {
    if _, exists := sm.students[id]; !exists {
        return fmt.Errorf("student with ID %d not found", id)
    }
    
    delete(sm.students, id)
    return nil
}

// Get students by age range
func (sm *StudentManager) GetStudentsByAgeRange(minAge, maxAge int) []*Student {
    var result []*Student
    
    for _, student := range sm.students {
        if student.Age >= minAge && student.Age <= maxAge {
            result = append(result, student)
        }
    }
    
    return result
}

// Get top performers
func (sm *StudentManager) GetTopPerformers(n int) []*Student {
    type studentAvg struct {
        Student *Student
        Average float64
    }
    
    var studentAvgs []studentAvg
    
    for _, student := range sm.students {
        if len(student.Grades) > 0 {
            avg, _ := sm.GetAverage(student.ID)
            studentAvgs = append(studentAvgs, studentAvg{student, avg})
        }
    }
    
    // Sort by average (descending)
    for i := 0; i < len(studentAvgs)-1; i++ {
        for j := i + 1; j < len(studentAvgs); j++ {
            if studentAvgs[i].Average < studentAvgs[j].Average {
                studentAvgs[i], studentAvgs[j] = studentAvgs[j], studentAvgs[i]
            }
        }
    }
    
    // Return top n
    var result []*Student
    for i := 0; i < n && i < len(studentAvgs); i++ {
        result = append(result, studentAvgs[i].Student)
    }
    
    return result
}

func main() {
    // Create student manager
    sm := NewStudentManager()
    
    // Add students
    alice := sm.AddStudent("Alice", 20)
    bob := sm.AddStudent("Bob", 19)
    charlie := sm.AddStudent("Charlie", 21)
    diana := sm.AddStudent("Diana", 20)
    
    fmt.Printf("Added students with IDs: %d, %d, %d, %d\n", alice, bob, charlie, diana)
    
    // Add grades
    sm.AddGrade(alice, "Math", 95.0)
    sm.AddGrade(alice, "Science", 87.0)
    sm.AddGrade(alice, "English", 92.0)
    
    sm.AddGrade(bob, "Math", 78.0)
    sm.AddGrade(bob, "Science", 85.0)
    sm.AddGrade(bob, "English", 88.0)
    
    sm.AddGrade(charlie, "Math", 92.0)
    sm.AddGrade(charlie, "Science", 94.0)
    sm.AddGrade(charlie, "English", 89.0)
    
    sm.AddGrade(diana, "Math", 88.0)
    sm.AddGrade(diana, "Science", 91.0)
    sm.AddGrade(diana, "English", 94.0)
    
    // Display all students
    fmt.Println("\nAll students:")
    for id, student := range sm.GetAllStudents() {
        avg, _ := sm.GetAverage(id)
        fmt.Printf("ID: %d, Name: %s, Age: %d, Average: %.2f\n", 
            student.ID, student.Name, student.Age, avg)
        fmt.Printf("  Grades: %v\n", student.Grades)
    }
    
    // Get students by age range
    fmt.Println("\nStudents aged 19-20:")
    youngStudents := sm.GetStudentsByAgeRange(19, 20)
    for _, student := range youngStudents {
        fmt.Printf("- %s (Age: %d)\n", student.Name, student.Age)
    }
    
    // Get top performers
    fmt.Println("\nTop 3 performers:")
    topStudents := sm.GetTopPerformers(3)
    for i, student := range topStudents {
        avg, _ := sm.GetAverage(student.ID)
        fmt.Printf("%d. %s (Average: %.2f)\n", i+1, student.Name, avg)
    }
}
```

### Example 3: Cache Implementation

```go
package main

import (
    "fmt"
    "time"
)

type CacheItem struct {
    Value      interface{}
    Expiration time.Time
}

type Cache struct {
    items map[string]CacheItem
    ttl   time.Duration
}

// Create new cache
func NewCache(ttl time.Duration) *Cache {
    return &Cache{
        items: make(map[string]CacheItem),
        ttl:   ttl,
    }
}

// Set item in cache
func (c *Cache) Set(key string, value interface{}) {
    expiration := time.Now().Add(c.ttl)
    c.items[key] = CacheItem{
        Value:      value,
        Expiration: expiration,
    }
}

// Get item from cache
func (c *Cache) Get(key string) (interface{}, bool) {
    item, exists := c.items[key]
    if !exists {
        return nil, false
    }
    
    // Check if expired
    if time.Now().After(item.Expiration) {
        delete(c.items, key)
        return nil, false
    }
    
    return item.Value, true
}

// Delete item from cache
func (c *Cache) Delete(key string) {
    delete(c.items, key)
}

// Clear expired items
func (c *Cache) ClearExpired() {
    now := time.Now()
    for key, item := range c.items {
        if now.After(item.Expiration) {
            delete(c.items, key)
        }
    }
}

// Get cache size
func (c *Cache) Size() int {
    return len(c.items)
}

// Clear all items
func (c *Cache) Clear() {
    c.items = make(map[string]CacheItem)
}

// Get all keys
func (c *Cache) Keys() []string {
    keys := make([]string, 0, len(c.items))
    for key := range c.items {
        keys = append(keys, key)
    }
    return keys
}

// Check if key exists and is not expired
func (c *Cache) Exists(key string) bool {
    _, exists := c.Get(key)
    return exists
}

func main() {
    // Create cache with 5 second TTL
    cache := NewCache(5 * time.Second)
    
    // Set some values
    cache.Set("user:1", map[string]interface{}{
        "name": "Alice",
        "age":  25,
    })
    cache.Set("user:2", map[string]interface{}{
        "name": "Bob",
        "age":  30,
    })
    cache.Set("config", "production")
    cache.Set("counter", 42)
    
    fmt.Printf("Cache size: %d\n", cache.Size())
    fmt.Printf("Keys: %v\n", cache.Keys())
    
    // Get values
    if user1, exists := cache.Get("user:1"); exists {
        fmt.Printf("User 1: %v\n", user1)
    }
    
    if config, exists := cache.Get("config"); exists {
        fmt.Printf("Config: %v\n", config)
    }
    
    // Check existence
    fmt.Printf("user:1 exists: %t\n", cache.Exists("user:1"))
    fmt.Printf("user:3 exists: %t\n", cache.Exists("user:3"))
    
    // Wait for expiration
    fmt.Println("\nWaiting 6 seconds for expiration...")
    time.Sleep(6 * time.Second)
    
    // Try to get expired values
    if _, exists := cache.Get("user:1"); exists {
        fmt.Println("User 1 still exists")
    } else {
        fmt.Println("User 1 has expired")
    }
    
    fmt.Printf("Cache size after expiration: %d\n", cache.Size())
    
    // Clear expired items
    cache.ClearExpired()
    fmt.Printf("Cache size after clearing expired: %d\n", cache.Size())
    
    // Add new items
    cache.Set("new_item", "Hello World")
    fmt.Printf("Cache size after adding new item: %d\n", cache.Size())
    
    // Clear all
    cache.Clear()
    fmt.Printf("Cache size after clearing all: %d\n", cache.Size())
}
```

## Exercises

### Exercise 1: Phone Book
Implement a phone book using maps.

```go
package main

import "fmt"

type Contact struct {
    Name  string
    Phone string
    Email string
}

type PhoneBook struct {
    contacts map[string]Contact
}

// Implement these methods:
// NewPhoneBook() *PhoneBook
// AddContact(name, phone, email string)
// GetContact(name string) (Contact, bool)
// UpdateContact(name, phone, email string) error
// DeleteContact(name string) error
// SearchByPhone(phone string) (Contact, bool)
// ListAllContacts() []Contact
// GetContactCount() int

func main() {
    // Test phone book functionality
}
```

### Exercise 2: Inventory Management
Create an inventory management system using maps.

```go
package main

import "fmt"

type Product struct {
    ID       string
    Name     string
    Price    float64
    Quantity int
    Category string
}

type Inventory struct {
    products map[string]Product
}

// Implement these methods:
// NewInventory() *Inventory
// AddProduct(product Product) error
// UpdateQuantity(id string, quantity int) error
// GetProduct(id string) (Product, bool)
// GetProductsByCategory(category string) []Product
// GetLowStockProducts(threshold int) []Product
// CalculateTotalValue() float64
// RemoveProduct(id string) error

func main() {
    // Test inventory management
}
```

### Exercise 3: Grade Book
Implement a grade book system using nested maps.

```go
package main

import "fmt"

type GradeBook struct {
    grades map[string]map[string]float64 // student -> subject -> grade
}

// Implement these methods:
// NewGradeBook() *GradeBook
// AddGrade(student, subject string, grade float64)
// GetGrade(student, subject string) (float64, bool)
// GetStudentAverage(student string) (float64, error)
// GetSubjectAverage(subject string) (float64, error)
// GetAllStudents() []string
// GetAllSubjects() []string
// GetTopStudent() (string, float64)
// GetStudentGrades(student string) map[string]float64

func main() {
    // Test grade book functionality
}
```

## Key Takeaways

1. **Maps** store key-value pairs with efficient lookup
2. **Zero value** of a map is nil
3. **make()** creates an empty map
4. **Map literals** can initialize maps with values
5. **Comma ok idiom** checks if a key exists
6. **delete()** removes key-value pairs
7. **Maps are reference types** - modifications affect all references
8. **Iteration order** is not guaranteed in maps
9. **Keys must be comparable** types
10. **Maps are not thread-safe** - use sync.Map for concurrent access

## Next Steps

Now that you understand maps, let's learn about [Structs](11-structs.md) for creating custom data types!

---

**Previous**: [← Arrays and Slices](09-arrays-slices.md) | **Next**: [Structs →](11-structs.md)