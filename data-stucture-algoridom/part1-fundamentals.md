# Part 1: Fundamentals

## Table of Contents
1. [Introduction to Data Structures and Algorithms](#introduction)
2. [C/C++ Basics Review](#cpp-basics)
3. [Time and Space Complexity](#complexity)
4. [Big O Notation](#big-o)
5. [Practice Problems](#practice)

---

## 1. Introduction to Data Structures and Algorithms {#introduction}

### What are Data Structures?
Data structures are ways of organizing and storing data in a computer so that it can be accessed and modified efficiently.

### What are Algorithms?
Algorithms are step-by-step procedures or formulas for solving problems.

### Why Study DSA?
- **Efficiency**: Write faster and more memory-efficient code
- **Problem Solving**: Develop logical thinking and problem-solving skills
- **Interviews**: Essential for technical interviews
- **Foundation**: Base for advanced computer science concepts

### Types of Data Structures

#### Linear Data Structures
- Arrays
- Linked Lists
- Stacks
- Queues

#### Non-Linear Data Structures
- Trees
- Graphs
- Hash Tables

---

## 2. C/C++ Basics Review {#cpp-basics}

### Memory Management

```cpp
#include <iostream>
#include <cstdlib>
using namespace std;

// Dynamic memory allocation in C
void cStyleMemory() {
    // Allocate memory
    int* arr = (int*)malloc(5 * sizeof(int));
    
    // Use memory
    for(int i = 0; i < 5; i++) {
        arr[i] = i + 1;
    }
    
    // Free memory
    free(arr);
}

// Dynamic memory allocation in C++
void cppStyleMemory() {
    // Allocate memory
    int* arr = new int[5];
    
    // Use memory
    for(int i = 0; i < 5; i++) {
        arr[i] = i + 1;
    }
    
    // Free memory
    delete[] arr;
}
```

### Pointers and References

```cpp
#include <iostream>
using namespace std;

void pointerBasics() {
    int x = 10;
    int* ptr = &x;  // Pointer to x
    int& ref = x;   // Reference to x
    
    cout << "Value of x: " << x << endl;
    cout << "Value through pointer: " << *ptr << endl;
    cout << "Value through reference: " << ref << endl;
    
    // Modifying through pointer
    *ptr = 20;
    cout << "After modifying through pointer: " << x << endl;
    
    // Modifying through reference
    ref = 30;
    cout << "After modifying through reference: " << x << endl;
}
```

### Structures and Classes

```cpp
// C-style structure
struct Point {
    int x, y;
};

// C++ class
class Rectangle {
private:
    int width, height;
    
public:
    // Constructor
    Rectangle(int w, int h) : width(w), height(h) {}
    
    // Methods
    int area() {
        return width * height;
    }
    
    void display() {
        cout << "Rectangle: " << width << "x" << height << endl;
    }
};
```

### Function Pointers and Templates

```cpp
// Function pointer example
int add(int a, int b) { return a + b; }
int multiply(int a, int b) { return a * b; }

void functionPointerExample() {
    int (*operation)(int, int);
    
    operation = add;
    cout << "Addition: " << operation(5, 3) << endl;
    
    operation = multiply;
    cout << "Multiplication: " << operation(5, 3) << endl;
}

// Template example
template<typename T>
T maximum(T a, T b) {
    return (a > b) ? a : b;
}

void templateExample() {
    cout << "Max of 5 and 3: " << maximum(5, 3) << endl;
    cout << "Max of 5.5 and 3.3: " << maximum(5.5, 3.3) << endl;
}
```

---

## 3. Time and Space Complexity {#complexity}

### Time Complexity
Time complexity measures how the runtime of an algorithm grows with input size.

### Space Complexity
Space complexity measures how much memory an algorithm uses relative to input size.

### Examples

```cpp
// O(1) - Constant Time
int getFirstElement(int arr[], int size) {
    if (size > 0) {
        return arr[0];  // Always takes same time
    }
    return -1;
}

// O(n) - Linear Time
int findElement(int arr[], int size, int target) {
    for (int i = 0; i < size; i++) {
        if (arr[i] == target) {
            return i;
        }
    }
    return -1;
}

// O(n²) - Quadratic Time
void bubbleSort(int arr[], int n) {
    for (int i = 0; i < n-1; i++) {
        for (int j = 0; j < n-i-1; j++) {
            if (arr[j] > arr[j+1]) {
                // Swap elements
                int temp = arr[j];
                arr[j] = arr[j+1];
                arr[j+1] = temp;
            }
        }
    }
}

// O(log n) - Logarithmic Time
int binarySearch(int arr[], int left, int right, int target) {
    while (left <= right) {
        int mid = left + (right - left) / 2;
        
        if (arr[mid] == target)
            return mid;
        
        if (arr[mid] < target)
            left = mid + 1;
        else
            right = mid - 1;
    }
    return -1;
}
```

---

## 4. Big O Notation {#big-o}

### Common Time Complexities (Best to Worst)

1. **O(1)** - Constant
2. **O(log n)** - Logarithmic
3. **O(n)** - Linear
4. **O(n log n)** - Linearithmic
5. **O(n²)** - Quadratic
6. **O(n³)** - Cubic
7. **O(2ⁿ)** - Exponential
8. **O(n!)** - Factorial

### Complexity Analysis Examples

```cpp
// Example 1: O(n)
void printArray(int arr[], int n) {
    for (int i = 0; i < n; i++) {  // Loop runs n times
        cout << arr[i] << " ";
    }
}

// Example 2: O(n²)
void printPairs(int arr[], int n) {
    for (int i = 0; i < n; i++) {      // Outer loop: n times
        for (int j = 0; j < n; j++) {  // Inner loop: n times
            cout << "(" << arr[i] << ", " << arr[j] << ") ";
        }
    }
}

// Example 3: O(log n)
int power(int base, int exp) {
    if (exp == 0) return 1;
    
    int half = power(base, exp/2);
    
    if (exp % 2 == 0)
        return half * half;
    else
        return base * half * half;
}

// Example 4: O(n log n)
void mergeSort(int arr[], int left, int right) {
    if (left < right) {
        int mid = left + (right - left) / 2;
        
        mergeSort(arr, left, mid);      // T(n/2)
        mergeSort(arr, mid + 1, right); // T(n/2)
        merge(arr, left, mid, right);   // O(n)
    }
}
```

### Space Complexity Examples

```cpp
// O(1) Space - Constant space
int sum(int a, int b) {
    return a + b;  // Only uses constant extra space
}

// O(n) Space - Linear space
int* createArray(int n) {
    int* arr = new int[n];  // Allocates n integers
    return arr;
}

// O(n) Space - Recursive call stack
int factorial(int n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1);  // n recursive calls on stack
}
```

### Rules for Big O Analysis

1. **Drop Constants**: O(2n) → O(n)
2. **Drop Lower Order Terms**: O(n² + n) → O(n²)
3. **Different Inputs Use Different Variables**: O(a + b) not O(n)
4. **Worst Case**: Consider the worst possible scenario

---

## 5. Practice Problems {#practice}

### Problem 1: Array Sum
```cpp
// Calculate sum of array elements
// Time: O(n), Space: O(1)
int arraySum(int arr[], int n) {
    int sum = 0;
    for (int i = 0; i < n; i++) {
        sum += arr[i];
    }
    return sum;
}
```

### Problem 2: Find Maximum
```cpp
// Find maximum element in array
// Time: O(n), Space: O(1)
int findMax(int arr[], int n) {
    if (n == 0) return -1;
    
    int max = arr[0];
    for (int i = 1; i < n; i++) {
        if (arr[i] > max) {
            max = arr[i];
        }
    }
    return max;
}
```

### Problem 3: Reverse Array
```cpp
// Reverse array in-place
// Time: O(n), Space: O(1)
void reverseArray(int arr[], int n) {
    int start = 0, end = n - 1;
    
    while (start < end) {
        // Swap elements
        int temp = arr[start];
        arr[start] = arr[end];
        arr[end] = temp;
        
        start++;
        end--;
    }
}
```

### Problem 4: Check if Array is Sorted
```cpp
// Check if array is sorted in ascending order
// Time: O(n), Space: O(1)
bool isSorted(int arr[], int n) {
    for (int i = 1; i < n; i++) {
        if (arr[i] < arr[i-1]) {
            return false;
        }
    }
    return true;
}
```

### Exercises

1. **Analyze Complexity**: Determine time and space complexity for each function above
2. **Two Sum Problem**: Find two numbers in array that add up to target
3. **Remove Duplicates**: Remove duplicates from sorted array
4. **Rotate Array**: Rotate array by k positions
5. **Missing Number**: Find missing number in array of 1 to n

### Complete Example Program

```cpp
#include <iostream>
using namespace std;

class ArrayOperations {
public:
    static void printArray(int arr[], int n) {
        for (int i = 0; i < n; i++) {
            cout << arr[i] << " ";
        }
        cout << endl;
    }
    
    static int linearSearch(int arr[], int n, int target) {
        for (int i = 0; i < n; i++) {
            if (arr[i] == target) {
                return i;
            }
        }
        return -1;
    }
};

int main() {
    int arr[] = {64, 34, 25, 12, 22, 11, 90};
    int n = sizeof(arr) / sizeof(arr[0]);
    
    cout << "Original array: ";
    ArrayOperations::printArray(arr, n);
    
    int target = 25;
    int index = ArrayOperations::linearSearch(arr, n, target);
    
    if (index != -1) {
        cout << "Element " << target << " found at index " << index << endl;
    } else {
        cout << "Element " << target << " not found" << endl;
    }
    
    return 0;
}
```

---

## Summary

In this part, you learned:
- Fundamentals of data structures and algorithms
- Essential C/C++ concepts for DSA
- Time and space complexity analysis
- Big O notation and its applications
- Basic problem-solving techniques

**Next**: [Part 2: Arrays and Strings](./part2-arrays-strings.md)

---

## Additional Resources

- Practice complexity analysis with different algorithms
- Implement basic operations and analyze their complexity
- Review C/C++ memory management concepts
- Solve more array-based problems to strengthen fundamentals