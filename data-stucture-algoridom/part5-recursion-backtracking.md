# Part 5: Recursion and Backtracking

## Table of Contents
1. [Introduction to Recursion](#introduction-recursion)
2. [Recursion Fundamentals](#recursion-fundamentals)
3. [Types of Recursion](#types-recursion)
4. [Recursive Problem Solving](#recursive-problem-solving)
5. [Introduction to Backtracking](#introduction-backtracking)
6. [Backtracking Algorithms](#backtracking-algorithms)
7. [Classic Recursive Problems](#classic-problems)
8. [Optimization Techniques](#optimization-techniques)
9. [Practice Exercises](#practice-exercises)

---

## 1. Introduction to Recursion {#introduction-recursion}

### What is Recursion?
Recursion is a programming technique where a function calls itself to solve a smaller instance of the same problem.

### Key Components of Recursion
1. **Base Case**: The condition that stops the recursion
2. **Recursive Case**: The function calling itself with modified parameters
3. **Progress**: Each recursive call should move toward the base case

### Why Use Recursion?
- **Natural Problem Decomposition**: Many problems have recursive structure
- **Cleaner Code**: Often more readable than iterative solutions
- **Mathematical Elegance**: Mirrors mathematical definitions
- **Tree/Graph Traversal**: Natural fit for hierarchical structures

### Recursion vs Iteration

| Aspect | Recursion | Iteration |
|--------|-----------|----------|
| Memory | O(n) stack space | O(1) space |
| Speed | Slower (function calls) | Faster |
| Code | Often cleaner | Can be complex |
| Debugging | Can be harder | Usually easier |

---

## 2. Recursion Fundamentals {#recursion-fundamentals}

### Basic Recursive Structure

```cpp
#include <iostream>
using namespace std;

// Template for recursive function
int recursiveFunction(int n) {
    // Base case
    if (n <= 0) {
        return baseValue;
    }
    
    // Recursive case
    return someOperation(recursiveFunction(n - 1));
}
```

### Simple Examples

```cpp
class BasicRecursion {
public:
    // Factorial: n! = n * (n-1)!
    static long long factorial(int n) {
        // Base case
        if (n <= 1) {
            return 1;
        }
        
        // Recursive case
        return n * factorial(n - 1);
    }
    
    // Fibonacci: F(n) = F(n-1) + F(n-2)
    static long long fibonacci(int n) {
        // Base cases
        if (n <= 1) {
            return n;
        }
        
        // Recursive case
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
    
    // Sum of first n natural numbers
    static int sum(int n) {
        if (n <= 0) {
            return 0;
        }
        return n + sum(n - 1);
    }
    
    // Power function: x^n
    static long long power(int x, int n) {
        if (n == 0) {
            return 1;
        }
        if (n == 1) {
            return x;
        }
        
        return x * power(x, n - 1);
    }
    
    // Optimized power function (divide and conquer)
    static long long fastPower(int x, int n) {
        if (n == 0) {
            return 1;
        }
        
        long long half = fastPower(x, n / 2);
        
        if (n % 2 == 0) {
            return half * half;
        } else {
            return x * half * half;
        }
    }
};
```

### Array and String Recursion

```cpp
class ArrayStringRecursion {
public:
    // Print array elements
    static void printArray(int arr[], int size, int index = 0) {
        if (index >= size) {
            return;
        }
        
        cout << arr[index] << " ";
        printArray(arr, size, index + 1);
    }
    
    // Find maximum in array
    static int findMax(int arr[], int size) {
        if (size == 1) {
            return arr[0];
        }
        
        int maxOfRest = findMax(arr + 1, size - 1);
        return max(arr[0], maxOfRest);
    }
    
    // Binary search
    static int binarySearch(int arr[], int left, int right, int target) {
        if (left > right) {
            return -1;
        }
        
        int mid = left + (right - left) / 2;
        
        if (arr[mid] == target) {
            return mid;
        } else if (arr[mid] > target) {
            return binarySearch(arr, left, mid - 1, target);
        } else {
            return binarySearch(arr, mid + 1, right, target);
        }
    }
    
    // Check if string is palindrome
    static bool isPalindrome(const string& str, int left, int right) {
        if (left >= right) {
            return true;
        }
        
        if (str[left] != str[right]) {
            return false;
        }
        
        return isPalindrome(str, left + 1, right - 1);
    }
    
    // Reverse string
    static string reverseString(const string& str, int index = 0) {
        if (index >= str.length()) {
            return "";
        }
        
        return reverseString(str, index + 1) + str[index];
    }
    
    // Count occurrences of character
    static int countChar(const string& str, char ch, int index = 0) {
        if (index >= str.length()) {
            return 0;
        }
        
        int count = (str[index] == ch) ? 1 : 0;
        return count + countChar(str, ch, index + 1);
    }
};
```

---

## 3. Types of Recursion {#types-recursion}

### Linear Recursion
Function makes only one recursive call.

```cpp
class LinearRecursion {
public:
    // Linear search
    static int linearSearch(int arr[], int size, int target, int index = 0) {
        if (index >= size) {
            return -1;
        }
        
        if (arr[index] == target) {
            return index;
        }
        
        return linearSearch(arr, size, target, index + 1);
    }
    
    // Print numbers from n to 1
    static void printDescending(int n) {
        if (n <= 0) {
            return;
        }
        
        cout << n << " ";
        printDescending(n - 1);
    }
    
    // Print numbers from 1 to n
    static void printAscending(int n) {
        if (n <= 0) {
            return;
        }
        
        printAscending(n - 1);
        cout << n << " ";
    }
};
```

### Binary Recursion
Function makes two recursive calls.

```cpp
class BinaryRecursion {
public:
    // Fibonacci (inefficient version)
    static long long fibonacci(int n) {
        if (n <= 1) {
            return n;
        }
        
        return fibonacci(n - 1) + fibonacci(n - 2);
    }
    
    // Tower of Hanoi
    static void towerOfHanoi(int n, char source, char destination, char auxiliary) {
        if (n == 1) {
            cout << "Move disk 1 from " << source << " to " << destination << endl;
            return;
        }
        
        towerOfHanoi(n - 1, source, auxiliary, destination);
        cout << "Move disk " << n << " from " << source << " to " << destination << endl;
        towerOfHanoi(n - 1, auxiliary, destination, source);
    }
    
    // Binary tree traversal (conceptual)
    struct TreeNode {
        int data;
        TreeNode* left;
        TreeNode* right;
        
        TreeNode(int val) : data(val), left(nullptr), right(nullptr) {}
    };
    
    static void inorderTraversal(TreeNode* root) {
        if (root == nullptr) {
            return;
        }
        
        inorderTraversal(root->left);
        cout << root->data << " ";
        inorderTraversal(root->right);
    }
};
```

### Tail Recursion
Recursive call is the last operation in the function.

```cpp
class TailRecursion {
public:
    // Tail recursive factorial
    static long long factorialTail(int n, long long accumulator = 1) {
        if (n <= 1) {
            return accumulator;
        }
        
        return factorialTail(n - 1, n * accumulator);
    }
    
    // Tail recursive Fibonacci
    static long long fibonacciTail(int n, long long a = 0, long long b = 1) {
        if (n == 0) {
            return a;
        }
        if (n == 1) {
            return b;
        }
        
        return fibonacciTail(n - 1, b, a + b);
    }
    
    // Tail recursive GCD
    static int gcdTail(int a, int b) {
        if (b == 0) {
            return a;
        }
        
        return gcdTail(b, a % b);
    }
};
```

### Multiple Recursion
Function makes more than two recursive calls.

```cpp
class MultipleRecursion {
public:
    // Tribonacci: T(n) = T(n-1) + T(n-2) + T(n-3)
    static long long tribonacci(int n) {
        if (n <= 1) {
            return n;
        }
        if (n == 2) {
            return 1;
        }
        
        return tribonacci(n - 1) + tribonacci(n - 2) + tribonacci(n - 3);
    }
    
    // Generate all subsets
    static void generateSubsets(vector<int>& nums, int index, vector<int>& current, vector<vector<int>>& result) {
        if (index >= nums.size()) {
            result.push_back(current);
            return;
        }
        
        // Include current element
        current.push_back(nums[index]);
        generateSubsets(nums, index + 1, current, result);
        
        // Exclude current element
        current.pop_back();
        generateSubsets(nums, index + 1, current, result);
    }
};
```

---

## 4. Recursive Problem Solving {#recursive-problem-solving}

### Problem-Solving Strategy

1. **Identify the base case(s)**
2. **Define the recursive relation**
3. **Ensure progress toward base case**
4. **Implement and test**

### Advanced Recursive Problems

```cpp
class AdvancedRecursion {
public:
    // Generate all permutations
    static void generatePermutations(vector<int>& nums, int start, vector<vector<int>>& result) {
        if (start >= nums.size()) {
            result.push_back(nums);
            return;
        }
        
        for (int i = start; i < nums.size(); i++) {
            swap(nums[start], nums[i]);
            generatePermutations(nums, start + 1, result);
            swap(nums[start], nums[i]); // backtrack
        }
    }
    
    // Generate all combinations
    static void generateCombinations(vector<int>& nums, int start, int k, 
                                   vector<int>& current, vector<vector<int>>& result) {
        if (current.size() == k) {
            result.push_back(current);
            return;
        }
        
        for (int i = start; i < nums.size(); i++) {
            current.push_back(nums[i]);
            generateCombinations(nums, i + 1, k, current, result);
            current.pop_back();
        }
    }
    
    // Merge sort
    static void mergeSort(vector<int>& arr, int left, int right) {
        if (left >= right) {
            return;
        }
        
        int mid = left + (right - left) / 2;
        
        mergeSort(arr, left, mid);
        mergeSort(arr, mid + 1, right);
        merge(arr, left, mid, right);
    }
    
    // Quick sort
    static void quickSort(vector<int>& arr, int low, int high) {
        if (low < high) {
            int pi = partition(arr, low, high);
            
            quickSort(arr, low, pi - 1);
            quickSort(arr, pi + 1, high);
        }
    }
    
private:
    static void merge(vector<int>& arr, int left, int mid, int right) {
        vector<int> temp(right - left + 1);
        int i = left, j = mid + 1, k = 0;
        
        while (i <= mid && j <= right) {
            if (arr[i] <= arr[j]) {
                temp[k++] = arr[i++];
            } else {
                temp[k++] = arr[j++];
            }
        }
        
        while (i <= mid) temp[k++] = arr[i++];
        while (j <= right) temp[k++] = arr[j++];
        
        for (i = left, k = 0; i <= right; i++, k++) {
            arr[i] = temp[k];
        }
    }
    
    static int partition(vector<int>& arr, int low, int high) {
        int pivot = arr[high];
        int i = low - 1;
        
        for (int j = low; j < high; j++) {
            if (arr[j] < pivot) {
                i++;
                swap(arr[i], arr[j]);
            }
        }
        
        swap(arr[i + 1], arr[high]);
        return i + 1;
    }
};
```

---

## 5. Introduction to Backtracking {#introduction-backtracking}

### What is Backtracking?
Backtracking is an algorithmic approach that considers searching every possible combination to solve computational problems.

### Backtracking Strategy
1. **Choose**: Make a choice from available options
2. **Explore**: Recursively explore the consequences
3. **Unchoose**: If the path doesn't lead to solution, backtrack

### Backtracking Template

```cpp
void backtrack(parameters) {
    // Base case
    if (isValidSolution(parameters)) {
        processSolution(parameters);
        return;
    }
    
    // Try all possible choices
    for (each choice in choices) {
        // Make choice
        makeChoice(choice);
        
        // Recurse
        backtrack(modifiedParameters);
        
        // Unmake choice (backtrack)
        unmakeChoice(choice);
    }
}
```

---

## 6. Backtracking Algorithms {#backtracking-algorithms}

### N-Queens Problem

```cpp
class NQueens {
public:
    static vector<vector<string>> solveNQueens(int n) {
        vector<vector<string>> solutions;
        vector<string> board(n, string(n, '.'));
        
        backtrack(board, 0, solutions);
        return solutions;
    }
    
private:
    static void backtrack(vector<string>& board, int row, vector<vector<string>>& solutions) {
        if (row == board.size()) {
            solutions.push_back(board);
            return;
        }
        
        for (int col = 0; col < board.size(); col++) {
            if (isSafe(board, row, col)) {
                board[row][col] = 'Q';
                backtrack(board, row + 1, solutions);
                board[row][col] = '.';
            }
        }
    }
    
    static bool isSafe(const vector<string>& board, int row, int col) {
        int n = board.size();
        
        // Check column
        for (int i = 0; i < row; i++) {
            if (board[i][col] == 'Q') {
                return false;
            }
        }
        
        // Check diagonal (top-left to bottom-right)
        for (int i = row - 1, j = col - 1; i >= 0 && j >= 0; i--, j--) {
            if (board[i][j] == 'Q') {
                return false;
            }
        }
        
        // Check diagonal (top-right to bottom-left)
        for (int i = row - 1, j = col + 1; i >= 0 && j < n; i--, j++) {
            if (board[i][j] == 'Q') {
                return false;
            }
        }
        
        return true;
    }
};
```

### Sudoku Solver

```cpp
class SudokuSolver {
public:
    static bool solveSudoku(vector<vector<char>>& board) {
        return backtrack(board);
    }
    
private:
    static bool backtrack(vector<vector<char>>& board) {
        for (int row = 0; row < 9; row++) {
            for (int col = 0; col < 9; col++) {
                if (board[row][col] == '.') {
                    for (char num = '1'; num <= '9'; num++) {
                        if (isValid(board, row, col, num)) {
                            board[row][col] = num;
                            
                            if (backtrack(board)) {
                                return true;
                            }
                            
                            board[row][col] = '.';
                        }
                    }
                    return false;
                }
            }
        }
        return true;
    }
    
    static bool isValid(const vector<vector<char>>& board, int row, int col, char num) {
        // Check row
        for (int j = 0; j < 9; j++) {
            if (board[row][j] == num) {
                return false;
            }
        }
        
        // Check column
        for (int i = 0; i < 9; i++) {
            if (board[i][col] == num) {
                return false;
            }
        }
        
        // Check 3x3 box
        int startRow = (row / 3) * 3;
        int startCol = (col / 3) * 3;
        
        for (int i = startRow; i < startRow + 3; i++) {
            for (int j = startCol; j < startCol + 3; j++) {
                if (board[i][j] == num) {
                    return false;
                }
            }
        }
        
        return true;
    }
};
```

### Maze Solving

```cpp
class MazeSolver {
public:
    static bool solveMaze(vector<vector<int>>& maze) {
        int n = maze.size();
        vector<vector<bool>> solution(n, vector<bool>(n, false));
        
        if (backtrack(maze, 0, 0, solution)) {
            printSolution(solution);
            return true;
        }
        
        return false;
    }
    
private:
    static bool backtrack(const vector<vector<int>>& maze, int x, int y, 
                         vector<vector<bool>>& solution) {
        int n = maze.size();
        
        // Base case: reached destination
        if (x == n - 1 && y == n - 1 && maze[x][y] == 1) {
            solution[x][y] = true;
            return true;
        }
        
        // Check if current cell is valid
        if (isSafe(maze, x, y)) {
            solution[x][y] = true;
            
            // Move right
            if (backtrack(maze, x, y + 1, solution)) {
                return true;
            }
            
            // Move down
            if (backtrack(maze, x + 1, y, solution)) {
                return true;
            }
            
            // Backtrack
            solution[x][y] = false;
        }
        
        return false;
    }
    
    static bool isSafe(const vector<vector<int>>& maze, int x, int y) {
        int n = maze.size();
        return (x >= 0 && x < n && y >= 0 && y < n && maze[x][y] == 1);
    }
    
    static void printSolution(const vector<vector<bool>>& solution) {
        for (const auto& row : solution) {
            for (bool cell : row) {
                cout << (cell ? "1 " : "0 ");
            }
            cout << endl;
        }
    }
};
```

### Word Search

```cpp
class WordSearch {
public:
    static bool exist(vector<vector<char>>& board, string word) {
        int m = board.size();
        int n = board[0].size();
        
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (backtrack(board, word, i, j, 0)) {
                    return true;
                }
            }
        }
        
        return false;
    }
    
private:
    static bool backtrack(vector<vector<char>>& board, const string& word, 
                         int row, int col, int index) {
        if (index == word.length()) {
            return true;
        }
        
        if (row < 0 || row >= board.size() || col < 0 || col >= board[0].size() ||
            board[row][col] != word[index]) {
            return false;
        }
        
        char temp = board[row][col];
        board[row][col] = '#'; // Mark as visited
        
        // Explore all 4 directions
        bool found = backtrack(board, word, row + 1, col, index + 1) ||
                    backtrack(board, word, row - 1, col, index + 1) ||
                    backtrack(board, word, row, col + 1, index + 1) ||
                    backtrack(board, word, row, col - 1, index + 1);
        
        board[row][col] = temp; // Restore
        
        return found;
    }
};
```

---

## 7. Classic Recursive Problems {#classic-problems}

### Mathematical Problems

```cpp
class MathematicalRecursion {
public:
    // Greatest Common Divisor (Euclidean Algorithm)
    static int gcd(int a, int b) {
        if (b == 0) {
            return a;
        }
        return gcd(b, a % b);
    }
    
    // Least Common Multiple
    static int lcm(int a, int b) {
        return (a * b) / gcd(a, b);
    }
    
    // Check if number is prime (recursive)
    static bool isPrime(int n, int divisor = 2) {
        if (n <= 1) return false;
        if (n == 2) return true;
        if (n % divisor == 0) return false;
        if (divisor * divisor > n) return true;
        
        return isPrime(n, divisor + 1);
    }
    
    // Sum of digits
    static int sumOfDigits(int n) {
        if (n == 0) {
            return 0;
        }
        return (n % 10) + sumOfDigits(n / 10);
    }
    
    // Reverse a number
    static int reverseNumber(int n, int reversed = 0) {
        if (n == 0) {
            return reversed;
        }
        return reverseNumber(n / 10, reversed * 10 + n % 10);
    }
    
    // Check if number is palindrome
    static bool isPalindrome(int n) {
        return n == reverseNumber(n);
    }
};
```

### String Problems

```cpp
class StringRecursion {
public:
    // Generate all subsequences
    static void generateSubsequences(const string& str, int index, 
                                   string current, vector<string>& result) {
        if (index >= str.length()) {
            result.push_back(current);
            return;
        }
        
        // Include current character
        generateSubsequences(str, index + 1, current + str[index], result);
        
        // Exclude current character
        generateSubsequences(str, index + 1, current, result);
    }
    
    // Check if string contains only digits
    static bool isAllDigits(const string& str, int index = 0) {
        if (index >= str.length()) {
            return true;
        }
        
        if (!isdigit(str[index])) {
            return false;
        }
        
        return isAllDigits(str, index + 1);
    }
    
    // Remove all occurrences of character
    static string removeChar(const string& str, char ch, int index = 0) {
        if (index >= str.length()) {
            return "";
        }
        
        if (str[index] == ch) {
            return removeChar(str, ch, index + 1);
        } else {
            return str[index] + removeChar(str, ch, index + 1);
        }
    }
    
    // Generate all permutations of string
    static void stringPermutations(string str, int start, vector<string>& result) {
        if (start >= str.length()) {
            result.push_back(str);
            return;
        }
        
        for (int i = start; i < str.length(); i++) {
            swap(str[start], str[i]);
            stringPermutations(str, start + 1, result);
            swap(str[start], str[i]); // backtrack
        }
    }
};
```

---

## 8. Optimization Techniques {#optimization-techniques}

### Memoization

```cpp
class Memoization {
public:
    // Memoized Fibonacci
    static long long fibonacciMemo(int n, unordered_map<int, long long>& memo) {
        if (n <= 1) {
            return n;
        }
        
        if (memo.find(n) != memo.end()) {
            return memo[n];
        }
        
        memo[n] = fibonacciMemo(n - 1, memo) + fibonacciMemo(n - 2, memo);
        return memo[n];
    }
    
    // Wrapper function
    static long long fibonacci(int n) {
        unordered_map<int, long long> memo;
        return fibonacciMemo(n, memo);
    }
    
    // Memoized combination (nCr)
    static long long combination(int n, int r, vector<vector<long long>>& memo) {
        if (r == 0 || r == n) {
            return 1;
        }
        
        if (memo[n][r] != -1) {
            return memo[n][r];
        }
        
        memo[n][r] = combination(n - 1, r - 1, memo) + combination(n - 1, r, memo);
        return memo[n][r];
    }
    
    // Longest Common Subsequence with memoization
    static int lcs(const string& s1, const string& s2, int i, int j, 
                  vector<vector<int>>& memo) {
        if (i >= s1.length() || j >= s2.length()) {
            return 0;
        }
        
        if (memo[i][j] != -1) {
            return memo[i][j];
        }
        
        if (s1[i] == s2[j]) {
            memo[i][j] = 1 + lcs(s1, s2, i + 1, j + 1, memo);
        } else {
            memo[i][j] = max(lcs(s1, s2, i + 1, j, memo), 
                           lcs(s1, s2, i, j + 1, memo));
        }
        
        return memo[i][j];
    }
};
```

### Tail Call Optimization

```cpp
class TailCallOptimization {
public:
    // Iterative version of tail recursive function
    static long long factorialIterative(int n) {
        long long result = 1;
        while (n > 1) {
            result *= n;
            n--;
        }
        return result;
    }
    
    // Convert tail recursion to iteration
    static long long fibonacciIterative(int n) {
        if (n <= 1) return n;
        
        long long a = 0, b = 1;
        for (int i = 2; i <= n; i++) {
            long long temp = a + b;
            a = b;
            b = temp;
        }
        return b;
    }
};
```

---

## 9. Practice Exercises {#practice-exercises}

### Basic Recursion Exercises
1. **Sum of Array**: Calculate sum of array elements recursively
2. **Product of Array**: Calculate product of array elements
3. **Count Zeros**: Count number of zeros in an integer
4. **Digital Root**: Find digital root of a number
5. **Binary Representation**: Convert decimal to binary recursively
6. **String Length**: Calculate string length without using built-in function
7. **Array Reverse**: Reverse array in-place recursively
8. **Nested Parentheses**: Check if parentheses are properly nested

### Intermediate Recursion Exercises
1. **Subset Sum**: Check if subset with given sum exists
2. **Coin Change**: Count ways to make change for given amount
3. **Staircase Problem**: Count ways to climb n stairs
4. **Josephus Problem**: Find survivor in Josephus problem
5. **Hanoi Variations**: Solve Tower of Hanoi with variations
6. **Palindrome Partitioning**: Partition string into palindromes
7. **Word Break**: Check if string can be segmented into words
8. **Expression Evaluation**: Evaluate mathematical expressions

### Advanced Backtracking Exercises
1. **Knight's Tour**: Find path for knight to visit all squares
2. **Rat in Maze**: Find all paths from source to destination
3. **Graph Coloring**: Color graph with minimum colors
4. **Hamiltonian Path**: Find Hamiltonian path in graph
5. **Subset with Equal Sum**: Partition array into equal sum subsets
6. **Cryptarithmetic**: Solve cryptarithmetic puzzles
7. **Crossword Puzzle**: Fill crossword puzzle
8. **Constraint Satisfaction**: Solve general CSP problems

### Complete Example Program

```cpp
#include <iostream>
#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

int main() {
    cout << "=== Recursion and Backtracking Demo ===" << endl;
    
    // Basic recursion examples
    cout << "\nBasic Recursion:" << endl;
    cout << "Factorial of 5: " << BasicRecursion::factorial(5) << endl;
    cout << "Fibonacci of 10: " << BasicRecursion::fibonacci(10) << endl;
    cout << "2^10 = " << BasicRecursion::fastPower(2, 10) << endl;
    
    // Array recursion
    cout << "\nArray Recursion:" << endl;
    int arr[] = {3, 7, 1, 9, 4};
    int size = sizeof(arr) / sizeof(arr[0]);
    cout << "Array: ";
    ArrayStringRecursion::printArray(arr, size);
    cout << "\nMaximum: " << ArrayStringRecursion::findMax(arr, size) << endl;
    
    // String recursion
    cout << "\nString Recursion:" << endl;
    string str = "racecar";
    bool isPalin = ArrayStringRecursion::isPalindrome(str, 0, str.length() - 1);
    cout << str << " is " << (isPalin ? "" : "not ") << "a palindrome" << endl;
    
    // Tower of Hanoi
    cout << "\nTower of Hanoi (3 disks):" << endl;
    BinaryRecursion::towerOfHanoi(3, 'A', 'C', 'B');
    
    // Backtracking - Generate permutations
    cout << "\nPermutations of [1,2,3]:" << endl;
    vector<int> nums = {1, 2, 3};
    vector<vector<int>> permutations;
    AdvancedRecursion::generatePermutations(nums, 0, permutations);
    
    for (const auto& perm : permutations) {
        cout << "[";
        for (int i = 0; i < perm.size(); i++) {
            cout << perm[i];
            if (i < perm.size() - 1) cout << ", ";
        }
        cout << "]" << endl;
    }
    
    // Memoized Fibonacci
    cout << "\nMemoized Fibonacci:" << endl;
    cout << "Fibonacci of 40: " << Memoization::fibonacci(40) << endl;
    
    return 0;
}
```

---

## Summary

In this part, you learned:
- Recursion fundamentals and different types
- Recursive problem-solving strategies
- Backtracking algorithms and applications
- Classic recursive problems and solutions
- Optimization techniques like memoization
- When to use recursion vs iteration

**Key Concepts:**
- **Base Case**: Essential to prevent infinite recursion
- **Recursive Case**: Problem decomposition into smaller subproblems
- **Backtracking**: Systematic exploration with undoing choices
- **Memoization**: Caching results to avoid redundant calculations
- **Tail Recursion**: Can be optimized to iteration

**Time and Space Complexity:**
- **Time**: Depends on number of recursive calls and work per call
- **Space**: O(depth) for call stack, additional space for memoization
- **Optimization**: Memoization can reduce time complexity significantly

**When to Use Recursion:**
- Problems with recursive structure (trees, graphs)
- Mathematical problems with recursive definitions
- Backtracking and constraint satisfaction problems
- Divide and conquer algorithms

**Next**: [Part 6: Trees](./part6-trees.md)

---

## Additional Practice Resources

- LeetCode Recursion and Backtracking problems
- HackerRank Recursion challenges
- GeeksforGeeks Recursion articles
- Practice converting recursive solutions to iterative
- Analyze time and space complexity for recursive algorithms
- Implement memoization for optimization