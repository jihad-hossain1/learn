# 9. Arrays and Slices

Arrays and slices are fundamental data structures in Go for storing collections of elements. While arrays have a fixed size, slices are dynamic and more commonly used.

## Arrays

Arrays are fixed-size sequences of elements of the same type.

### Array Declaration

```go
// Declare array with specific size
var arrayName [size]dataType

// Examples
var numbers [5]int
var names [3]string
var flags [4]bool
```

### Array Initialization

```go
package main

import "fmt"

func main() {
    // Method 1: Declare then assign
    var numbers [5]int
    numbers[0] = 10
    numbers[1] = 20
    numbers[2] = 30
    
    // Method 2: Declare and initialize
    var fruits = [3]string{"apple", "banana", "orange"}
    
    // Method 3: Short declaration with initialization
    colors := [4]string{"red", "green", "blue", "yellow"}
    
    // Method 4: Let Go infer the size
    grades := [...]int{85, 90, 78, 92, 88}
    
    // Method 5: Initialize specific indices
    sparse := [10]int{2: 100, 5: 200, 8: 300}
    
    fmt.Printf("Numbers: %v\n", numbers)
    fmt.Printf("Fruits: %v\n", fruits)
    fmt.Printf("Colors: %v\n", colors)
    fmt.Printf("Grades: %v\n", grades)
    fmt.Printf("Sparse: %v\n", sparse)
}
```

### Array Operations

```go
package main

import "fmt"

func main() {
    numbers := [5]int{10, 20, 30, 40, 50}
    
    // Access elements
    fmt.Printf("First element: %d\n", numbers[0])
    fmt.Printf("Last element: %d\n", numbers[4])
    
    // Modify elements
    numbers[2] = 35
    fmt.Printf("Modified array: %v\n", numbers)
    
    // Array length
    fmt.Printf("Array length: %d\n", len(numbers))
    
    // Iterate through array
    fmt.Println("Array elements:")
    for i := 0; i < len(numbers); i++ {
        fmt.Printf("Index %d: %d\n", i, numbers[i])
    }
    
    // Iterate using range
    fmt.Println("\nUsing range:")
    for index, value := range numbers {
        fmt.Printf("Index %d: %d\n", index, value)
    }
    
    // Iterate values only
    fmt.Println("\nValues only:")
    for _, value := range numbers {
        fmt.Printf("%d ", value)
    }
    fmt.Println()
}
```

### Multidimensional Arrays

```go
package main

import "fmt"

func main() {
    // 2D array (matrix)
    var matrix [3][3]int
    
    // Initialize 2D array
    grid := [2][3]int{
        {1, 2, 3},
        {4, 5, 6},
    }
    
    // 3D array
    cube := [2][2][2]int{
        {
            {1, 2},
            {3, 4},
        },
        {
            {5, 6},
            {7, 8},
        },
    }
    
    // Access and modify 2D array
    matrix[0][0] = 10
    matrix[1][1] = 20
    matrix[2][2] = 30
    
    fmt.Printf("Matrix: %v\n", matrix)
    fmt.Printf("Grid: %v\n", grid)
    fmt.Printf("Cube: %v\n", cube)
    
    // Iterate through 2D array
    fmt.Println("\nGrid elements:")
    for i := 0; i < len(grid); i++ {
        for j := 0; j < len(grid[i]); j++ {
            fmt.Printf("grid[%d][%d] = %d\n", i, j, grid[i][j])
        }
    }
    
    // Using range with 2D array
    fmt.Println("\nUsing range:")
    for i, row := range grid {
        for j, value := range row {
            fmt.Printf("grid[%d][%d] = %d\n", i, j, value)
        }
    }
}
```

## Slices

Slices are dynamic arrays that can grow and shrink. They are more flexible and commonly used than arrays.

### Slice Declaration and Initialization

```go
package main

import "fmt"

func main() {
    // Method 1: Declare empty slice
    var numbers []int
    fmt.Printf("Empty slice: %v, length: %d, capacity: %d\n", numbers, len(numbers), cap(numbers))
    
    // Method 2: Initialize with values
    fruits := []string{"apple", "banana", "orange"}
    fmt.Printf("Fruits: %v, length: %d, capacity: %d\n", fruits, len(fruits), cap(fruits))
    
    // Method 3: Using make function
    grades := make([]int, 5)      // length 5, capacity 5
    scores := make([]int, 3, 10)  // length 3, capacity 10
    
    fmt.Printf("Grades: %v, length: %d, capacity: %d\n", grades, len(grades), cap(grades))
    fmt.Printf("Scores: %v, length: %d, capacity: %d\n", scores, len(scores), cap(scores))
    
    // Method 4: From array
    array := [5]int{1, 2, 3, 4, 5}
    slice := array[1:4] // Elements from index 1 to 3
    fmt.Printf("Slice from array: %v\n", slice)
}
```

### Slice Operations

```go
package main

import "fmt"

func main() {
    // Create a slice
    numbers := []int{10, 20, 30}
    fmt.Printf("Original: %v\n", numbers)
    
    // Append elements
    numbers = append(numbers, 40)
    numbers = append(numbers, 50, 60, 70)
    fmt.Printf("After append: %v\n", numbers)
    
    // Append another slice
    moreNumbers := []int{80, 90, 100}
    numbers = append(numbers, moreNumbers...)
    fmt.Printf("After appending slice: %v\n", numbers)
    
    // Access elements
    fmt.Printf("First element: %d\n", numbers[0])
    fmt.Printf("Last element: %d\n", numbers[len(numbers)-1])
    
    // Modify elements
    numbers[0] = 15
    fmt.Printf("After modification: %v\n", numbers)
    
    // Length and capacity
    fmt.Printf("Length: %d, Capacity: %d\n", len(numbers), cap(numbers))
}
```

### Slice Slicing

Slicing creates a new slice from an existing slice or array.

```go
package main

import "fmt"

func main() {
    numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
    fmt.Printf("Original: %v\n", numbers)
    
    // Basic slicing: slice[start:end]
    slice1 := numbers[2:6]  // Elements from index 2 to 5
    fmt.Printf("numbers[2:6]: %v\n", slice1)
    
    // Omit start (defaults to 0)
    slice2 := numbers[:5]   // Elements from index 0 to 4
    fmt.Printf("numbers[:5]: %v\n", slice2)
    
    // Omit end (defaults to length)
    slice3 := numbers[3:]   // Elements from index 3 to end
    fmt.Printf("numbers[3:]: %v\n", slice3)
    
    // Full slice
    slice4 := numbers[:]    // All elements
    fmt.Printf("numbers[:]: %v\n", slice4)
    
    // Three-index slicing: slice[start:end:capacity]
    slice5 := numbers[1:5:7] // Elements 1-4, capacity limited to 7-1=6
    fmt.Printf("numbers[1:5:7]: %v, len: %d, cap: %d\n", slice5, len(slice5), cap(slice5))
    
    // Slices share underlying array
    slice1[0] = 999
    fmt.Printf("After modifying slice1: %v\n", numbers)
    fmt.Printf("slice1: %v\n", slice1)
}
```

### Copy Function

```go
package main

import "fmt"

func main() {
    source := []int{1, 2, 3, 4, 5}
    
    // Create destination slice
    dest := make([]int, len(source))
    
    // Copy elements
    copied := copy(dest, source)
    fmt.Printf("Source: %v\n", source)
    fmt.Printf("Destination: %v\n", dest)
    fmt.Printf("Elements copied: %d\n", copied)
    
    // Modify destination (doesn't affect source)
    dest[0] = 999
    fmt.Printf("After modifying dest: source=%v, dest=%v\n", source, dest)
    
    // Copy with different sizes
    small := make([]int, 3)
    copied = copy(small, source)
    fmt.Printf("Copy to smaller slice: %v, copied: %d\n", small, copied)
    
    large := make([]int, 8)
    copied = copy(large, source)
    fmt.Printf("Copy to larger slice: %v, copied: %d\n", large, copied)
}
```

### Slice Internals

```go
package main

import "fmt"

func main() {
    // Understanding slice growth
    var slice []int
    fmt.Printf("Initial: len=%d, cap=%d\n", len(slice), cap(slice))
    
    for i := 0; i < 10; i++ {
        slice = append(slice, i)
        fmt.Printf("After append %d: len=%d, cap=%d, slice=%v\n", i, len(slice), cap(slice), slice)
    }
    
    // Demonstrating capacity growth
    fmt.Println("\nCapacity growth pattern:")
    slice = nil
    for i := 0; i < 20; i++ {
        oldCap := cap(slice)
        slice = append(slice, i)
        newCap := cap(slice)
        if newCap != oldCap {
            fmt.Printf("Capacity changed from %d to %d at length %d\n", oldCap, newCap, len(slice))
        }
    }
}
```

## Practical Examples

### Example 1: Dynamic Array Operations

```go
package main

import "fmt"

// Insert element at specific index
func insert(slice []int, index, value int) []int {
    if index < 0 || index > len(slice) {
        return slice // Invalid index
    }
    
    // Grow slice by one
    slice = append(slice, 0)
    
    // Shift elements to the right
    copy(slice[index+1:], slice[index:])
    
    // Insert new value
    slice[index] = value
    
    return slice
}

// Remove element at specific index
func remove(slice []int, index int) []int {
    if index < 0 || index >= len(slice) {
        return slice // Invalid index
    }
    
    return append(slice[:index], slice[index+1:]...)
}

// Find index of element
func indexOf(slice []int, value int) int {
    for i, v := range slice {
        if v == value {
            return i
        }
    }
    return -1 // Not found
}

// Check if slice contains element
func contains(slice []int, value int) bool {
    return indexOf(slice, value) != -1
}

// Remove duplicates
func removeDuplicates(slice []int) []int {
    seen := make(map[int]bool)
    result := []int{}
    
    for _, value := range slice {
        if !seen[value] {
            seen[value] = true
            result = append(result, value)
        }
    }
    
    return result
}

func main() {
    numbers := []int{1, 2, 3, 5, 6}
    fmt.Printf("Original: %v\n", numbers)
    
    // Insert element
    numbers = insert(numbers, 3, 4)
    fmt.Printf("After inserting 4 at index 3: %v\n", numbers)
    
    // Remove element
    numbers = remove(numbers, 0)
    fmt.Printf("After removing element at index 0: %v\n", numbers)
    
    // Find element
    index := indexOf(numbers, 5)
    fmt.Printf("Index of 5: %d\n", index)
    
    // Check if contains
    fmt.Printf("Contains 3: %t\n", contains(numbers, 3))
    fmt.Printf("Contains 10: %t\n", contains(numbers, 10))
    
    // Remove duplicates
    duplicates := []int{1, 2, 2, 3, 3, 3, 4, 5, 5}
    unique := removeDuplicates(duplicates)
    fmt.Printf("Original with duplicates: %v\n", duplicates)
    fmt.Printf("After removing duplicates: %v\n", unique)
}
```

### Example 2: Matrix Operations

```go
package main

import "fmt"

// Create matrix
func createMatrix(rows, cols int) [][]int {
    matrix := make([][]int, rows)
    for i := range matrix {
        matrix[i] = make([]int, cols)
    }
    return matrix
}

// Fill matrix with values
func fillMatrix(matrix [][]int, value int) {
    for i := range matrix {
        for j := range matrix[i] {
            matrix[i][j] = value
        }
    }
}

// Print matrix
func printMatrix(matrix [][]int) {
    for _, row := range matrix {
        for _, value := range row {
            fmt.Printf("%3d ", value)
        }
        fmt.Println()
    }
}

// Matrix addition
func addMatrices(a, b [][]int) [][]int {
    if len(a) != len(b) || len(a[0]) != len(b[0]) {
        return nil // Incompatible dimensions
    }
    
    result := createMatrix(len(a), len(a[0]))
    for i := range a {
        for j := range a[i] {
            result[i][j] = a[i][j] + b[i][j]
        }
    }
    return result
}

// Matrix transpose
func transpose(matrix [][]int) [][]int {
    rows, cols := len(matrix), len(matrix[0])
    result := createMatrix(cols, rows)
    
    for i := 0; i < rows; i++ {
        for j := 0; j < cols; j++ {
            result[j][i] = matrix[i][j]
        }
    }
    return result
}

func main() {
    // Create and initialize matrices
    matrix1 := [][]int{
        {1, 2, 3},
        {4, 5, 6},
    }
    
    matrix2 := [][]int{
        {7, 8, 9},
        {10, 11, 12},
    }
    
    fmt.Println("Matrix 1:")
    printMatrix(matrix1)
    
    fmt.Println("\nMatrix 2:")
    printMatrix(matrix2)
    
    // Add matrices
    sum := addMatrices(matrix1, matrix2)
    fmt.Println("\nSum:")
    printMatrix(sum)
    
    // Transpose matrix
    transposed := transpose(matrix1)
    fmt.Println("\nMatrix 1 Transposed:")
    printMatrix(transposed)
    
    // Create identity matrix
    identity := createMatrix(3, 3)
    for i := 0; i < 3; i++ {
        identity[i][i] = 1
    }
    fmt.Println("\nIdentity Matrix:")
    printMatrix(identity)
}
```

### Example 3: Sorting and Searching

```go
package main

import "fmt"

// Bubble sort
func bubbleSort(slice []int) {
    n := len(slice)
    for i := 0; i < n-1; i++ {
        for j := 0; j < n-i-1; j++ {
            if slice[j] > slice[j+1] {
                slice[j], slice[j+1] = slice[j+1], slice[j]
            }
        }
    }
}

// Selection sort
func selectionSort(slice []int) {
    n := len(slice)
    for i := 0; i < n-1; i++ {
        minIdx := i
        for j := i + 1; j < n; j++ {
            if slice[j] < slice[minIdx] {
                minIdx = j
            }
        }
        slice[i], slice[minIdx] = slice[minIdx], slice[i]
    }
}

// Binary search (requires sorted slice)
func binarySearch(slice []int, target int) int {
    left, right := 0, len(slice)-1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if slice[mid] == target {
            return mid
        } else if slice[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1 // Not found
}

// Linear search
func linearSearch(slice []int, target int) int {
    for i, value := range slice {
        if value == target {
            return i
        }
    }
    return -1 // Not found
}

// Find min and max
func findMinMax(slice []int) (int, int) {
    if len(slice) == 0 {
        return 0, 0
    }
    
    min, max := slice[0], slice[0]
    for _, value := range slice[1:] {
        if value < min {
            min = value
        }
        if value > max {
            max = value
        }
    }
    return min, max
}

func main() {
    // Test data
    numbers := []int{64, 34, 25, 12, 22, 11, 90}
    fmt.Printf("Original: %v\n", numbers)
    
    // Test bubble sort
    bubbleData := make([]int, len(numbers))
    copy(bubbleData, numbers)
    bubbleSort(bubbleData)
    fmt.Printf("Bubble sorted: %v\n", bubbleData)
    
    // Test selection sort
    selectionData := make([]int, len(numbers))
    copy(selectionData, numbers)
    selectionSort(selectionData)
    fmt.Printf("Selection sorted: %v\n", selectionData)
    
    // Test searching
    target := 25
    linearIdx := linearSearch(numbers, target)
    binaryIdx := binarySearch(bubbleData, target)
    
    fmt.Printf("\nSearching for %d:\n", target)
    fmt.Printf("Linear search result: %d\n", linearIdx)
    fmt.Printf("Binary search result: %d\n", binaryIdx)
    
    // Find min and max
    min, max := findMinMax(numbers)
    fmt.Printf("\nMin: %d, Max: %d\n", min, max)
}
```

## Exercises

### Exercise 1: Slice Statistics
Create functions to calculate statistics for a slice of numbers.

```go
package main

import "fmt"

// Implement these functions:
// sum(numbers []float64) float64
// average(numbers []float64) float64
// median(numbers []float64) float64
// mode(numbers []int) int
// standardDeviation(numbers []float64) float64

func main() {
    // Test with sample data
    data := []float64{1.5, 2.3, 3.7, 2.3, 4.1, 5.2, 2.3, 6.8, 3.7, 4.9}
    
    // Calculate and display all statistics
}
```

### Exercise 2: 2D Array Game Board
Create a tic-tac-toe game board using a 2D array.

```go
package main

import "fmt"

// Implement these functions:
// initializeBoard() [3][3]string
// printBoard(board [3][3]string)
// makeMove(board *[3][3]string, row, col int, player string) bool
// checkWinner(board [3][3]string) string
// isBoardFull(board [3][3]string) bool

func main() {
    // Implement a simple tic-tac-toe game
}
```

### Exercise 3: Dynamic Shopping Cart
Implement a shopping cart using slices.

```go
package main

import "fmt"

type Item struct {
    Name     string
    Price    float64
    Quantity int
}

// Implement these functions:
// addItem(cart []Item, item Item) []Item
// removeItem(cart []Item, name string) []Item
// updateQuantity(cart []Item, name string, quantity int) []Item
// calculateTotal(cart []Item) float64
// printCart(cart []Item)

func main() {
    // Test shopping cart functionality
}
```

## Key Takeaways

1. **Arrays** have fixed size, **slices** are dynamic
2. **Slices** are more commonly used than arrays
3. **make()** creates slices with specified length and capacity
4. **append()** adds elements to slices
5. **copy()** copies elements between slices
6. **Slicing** creates new slices from existing ones
7. **Slices share underlying arrays** - modifications affect all slices
8. **Capacity** determines how much a slice can grow before reallocation
9. **Range** provides easy iteration over arrays and slices
10. **Zero value** of slice is nil

## Next Steps

Now that you understand arrays and slices, let's learn about [Maps](10-maps.md) for key-value data storage!

---

**Previous**: [← Functions](08-functions.md) | **Next**: [Maps →](10-maps.md)