# Practice Exercises for C and C++

This file contains additional practice exercises to reinforce your learning. Complete these exercises after studying each section of the main guide.

## 📋 Exercise Categories

- [Beginner Exercises (C Fundamentals)](#beginner-exercises-c-fundamentals)
- [Intermediate Exercises (Advanced C)](#intermediate-exercises-advanced-c)
- [Advanced Exercises (C++ and OOP)](#advanced-exercises-c-and-oop)
- [Project Ideas](#project-ideas)
- [Challenge Problems](#challenge-problems)

---

## Beginner Exercises (C Fundamentals)

### Exercise 1: Basic Input/Output
**Difficulty:** ⭐

Write a program that:
1. Asks the user for their name, age, and favorite number
2. Calculates what their age will be in 10 years
3. Displays a personalized message with all the information

**Sample Output:**
```
Enter your name: Alice
Enter your age: 25
Enter your favorite number: 7

Hello Alice! You are 25 years old.
In 10 years, you will be 35 years old.
Your favorite number 7 multiplied by your age is 175.
```

### Exercise 2: Temperature Converter
**Difficulty:** ⭐

Create a program that converts temperatures between Celsius, Fahrenheit, and Kelvin.
- Formula: F = C × 9/5 + 32
- Formula: K = C + 273.15

**Requirements:**
- Display a menu for conversion options
- Handle invalid input gracefully
- Allow multiple conversions in one run

### Exercise 3: Simple Calculator
**Difficulty:** ⭐⭐

Build a calculator that:
1. Shows a menu of operations (+, -, *, /, %)
2. Takes two numbers from the user
3. Performs the selected operation
4. Handles division by zero
5. Continues until the user chooses to exit

### Exercise 4: Number Pattern Generator
**Difficulty:** ⭐⭐

Write a program that generates various number patterns:

**Pattern 1:**
```
1
1 2
1 2 3
1 2 3 4
1 2 3 4 5
```

**Pattern 2:**
```
    1
   2 2
  3 3 3
 4 4 4 4
5 5 5 5 5
```

**Pattern 3 (Pascal's Triangle):**
```
    1
   1 1
  1 2 1
 1 3 3 1
1 4 6 4 1
```

### Exercise 5: Prime Number Checker
**Difficulty:** ⭐⭐

Create a program that:
1. Checks if a number is prime
2. Finds all prime numbers in a given range
3. Finds the largest prime number less than a given number
4. Calculates the sum of all primes up to a given number

---

## Intermediate Exercises (Advanced C)

### Exercise 6: Array Operations
**Difficulty:** ⭐⭐

Implement the following array operations:
1. **Array Statistics**: Find min, max, sum, average, median
2. **Array Sorting**: Implement bubble sort and selection sort
3. **Array Searching**: Linear search and binary search
4. **Array Rotation**: Rotate array left/right by k positions

**Bonus:** Compare performance of different sorting algorithms.

### Exercise 7: String Manipulation Library
**Difficulty:** ⭐⭐⭐

Create a library of string functions:
```c
int stringLength(char str[]);
void stringReverse(char str[]);
void stringToUpper(char str[]);
void stringToLower(char str[]);
int stringCompare(char str1[], char str2[]);
void stringCopy(char dest[], char src[]);
int countWords(char str[]);
void removeSpaces(char str[]);
int isPalindrome(char str[]);
```

### Exercise 8: Dynamic Memory Management
**Difficulty:** ⭐⭐⭐

Write programs that:
1. Create a dynamic array that grows as needed
2. Implement a simple linked list with insert/delete operations
3. Create a dynamic 2D matrix for matrix operations
4. Build a simple memory pool allocator

### Exercise 9: File Processing
**Difficulty:** ⭐⭐⭐

Create programs that:
1. **Word Counter**: Count words, lines, and characters in a text file
2. **CSV Parser**: Read and process comma-separated values
3. **Log Analyzer**: Parse log files and extract statistics
4. **File Merger**: Merge multiple sorted files into one sorted file

### Exercise 10: Data Structures
**Difficulty:** ⭐⭐⭐

Implement these data structures in C:
1. **Stack**: With push, pop, peek, isEmpty operations
2. **Queue**: With enqueue, dequeue, front, rear operations
3. **Binary Tree**: With insert, search, traversal operations
4. **Hash Table**: With insert, search, delete operations

---

## Advanced Exercises (C++ and OOP)

### Exercise 11: Class Design
**Difficulty:** ⭐⭐

Design and implement these classes:

**BankAccount Class:**
```cpp
class BankAccount {
private:
    string accountNumber;
    string holderName;
    double balance;
    vector<string> transactionHistory;
    
public:
    // Constructors, getters, setters
    void deposit(double amount);
    bool withdraw(double amount);
    void displayStatement();
    double getInterest(double rate);
};
```

**Rectangle Class:**
```cpp
class Rectangle {
private:
    double length, width;
    
public:
    // Implement constructors, area, perimeter, comparison operators
};
```

### Exercise 12: Inheritance and Polymorphism
**Difficulty:** ⭐⭐⭐

Create a shape hierarchy:
```cpp
class Shape {
public:
    virtual double area() = 0;
    virtual double perimeter() = 0;
    virtual void display() = 0;
};

class Circle : public Shape { /* implement */ };
class Rectangle : public Shape { /* implement */ };
class Triangle : public Shape { /* implement */ };
```

Create a program that:
1. Stores different shapes in a vector
2. Calculates total area of all shapes
3. Finds the largest shape by area
4. Sorts shapes by perimeter

### Exercise 13: Template Programming
**Difficulty:** ⭐⭐⭐

Implement these generic templates:

**Generic Stack:**
```cpp
template <typename T>
class Stack {
private:
    vector<T> elements;
    
public:
    void push(const T& element);
    T pop();
    T top() const;
    bool empty() const;
    size_t size() const;
};
```

**Generic Algorithms:**
```cpp
template <typename T>
T findMax(const vector<T>& vec);

template <typename T>
void bubbleSort(vector<T>& vec);

template <typename T>
int binarySearch(const vector<T>& vec, const T& target);
```

### Exercise 14: STL Mastery
**Difficulty:** ⭐⭐⭐

Solve these problems using STL:

1. **Word Frequency Counter**: Use `map` to count word frequencies in text
2. **Student Grade Manager**: Use `vector` and `algorithm` functions
3. **Phone Book**: Use `unordered_map` for fast lookups
4. **Task Scheduler**: Use `priority_queue` for task management
5. **Set Operations**: Implement union, intersection, difference using `set`

---

## Project Ideas

### Beginner Projects

1. **Tic-Tac-Toe Game**
   - 3x3 grid implementation
   - Player vs Player or Player vs Computer
   - Win condition checking

2. **Simple Text Editor**
   - Basic file operations (open, save, edit)
   - Search and replace functionality
   - Line numbering

3. **Expense Tracker**
   - Add/remove expenses
   - Categorize expenses
   - Generate monthly reports
   - File persistence

### Intermediate Projects

4. **Library Management System**
   - Book inventory management
   - Member registration
   - Book borrowing/returning
   - Fine calculation

5. **Snake Game**
   - Console-based implementation
   - Score tracking
   - Increasing difficulty
   - High score persistence

6. **Scientific Calculator**
   - Advanced mathematical functions
   - Expression parsing
   - Memory functions
   - History tracking

### Advanced Projects

7. **Mini Database System**
   - Table creation and management
   - Basic SQL-like queries
   - Indexing for fast searches
   - Transaction support

8. **Chat Application**
   - Multi-client support
   - Message broadcasting
   - Private messaging
   - File transfer

9. **Compiler/Interpreter**
   - Simple programming language
   - Lexical analysis
   - Parsing and execution
   - Error handling

---

## Challenge Problems

### Algorithm Challenges

1. **Sorting Algorithms**: Implement and compare:
   - Bubble Sort, Selection Sort, Insertion Sort
   - Merge Sort, Quick Sort, Heap Sort
   - Radix Sort, Counting Sort

2. **Graph Algorithms**: Implement:
   - Breadth-First Search (BFS)
   - Depth-First Search (DFS)
   - Dijkstra's Shortest Path
   - Minimum Spanning Tree

3. **Dynamic Programming**: Solve:
   - Fibonacci sequence (optimized)
   - Longest Common Subsequence
   - Knapsack Problem
   - Coin Change Problem

### Mathematical Problems

4. **Number Theory**:
   - Greatest Common Divisor (GCD)
   - Least Common Multiple (LCM)
   - Prime factorization
   - Modular arithmetic

5. **Matrix Operations**:
   - Matrix multiplication
   - Matrix determinant
   - Matrix inverse
   - Eigenvalue calculation

### System Programming

6. **Memory Management**:
   - Custom memory allocator
   - Garbage collector simulation
   - Memory leak detector
   - Buffer overflow protection

7. **File System Simulation**:
   - Directory structure
   - File operations
   - Permission system
   - Disk space management

---

## 💡 Tips for Success

### Before Starting an Exercise:
1. **Read the problem carefully** - Understand what's being asked
2. **Plan your approach** - Write pseudocode or draw diagrams
3. **Start simple** - Implement basic functionality first
4. **Test frequently** - Test each function as you write it

### While Coding:
1. **Use meaningful variable names**
2. **Add comments to explain complex logic**
3. **Handle edge cases and errors**
4. **Follow consistent coding style**

### After Completing:
1. **Test with various inputs**
2. **Optimize for better performance**
3. **Refactor for better readability**
4. **Document your solution**

### Getting Help:
1. **Review the main guide** for concept clarification
2. **Check example programs** for similar patterns
3. **Use online resources** like cppreference.com
4. **Practice debugging** with print statements or debugger

---

## 📚 Exercise Solutions

Solutions to these exercises are not provided intentionally. The goal is to:
- Encourage independent problem-solving
- Develop debugging skills
- Build confidence through practice
- Learn from mistakes and iterations

Remember: There are often multiple correct solutions to a problem. Focus on:
- **Correctness**: Does it work for all test cases?
- **Efficiency**: Is it reasonably fast?
- **Readability**: Can others understand your code?
- **Maintainability**: Is it easy to modify?

---

**Happy Practicing! 🎯**

Remember: The key to mastering programming is consistent practice. Start with easier exercises and gradually work your way up to more challenging problems. Don't be discouraged if you get stuck—debugging and problem-solving are essential programming skills that improve with experience!