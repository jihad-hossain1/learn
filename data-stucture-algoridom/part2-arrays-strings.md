# Part 2: Arrays and Strings

## Table of Contents
1. [Arrays Fundamentals](#arrays-fundamentals)
2. [Multi-dimensional Arrays](#multi-dimensional-arrays)
3. [Dynamic Arrays](#dynamic-arrays)
4. [String Fundamentals](#string-fundamentals)
5. [String Algorithms](#string-algorithms)
6. [Common Array Problems](#array-problems)
7. [Common String Problems](#string-problems)
8. [Practice Exercises](#practice-exercises)

---

## 1. Arrays Fundamentals {#arrays-fundamentals}

### What is an Array?
An array is a collection of elements of the same data type stored in contiguous memory locations.

### Array Declaration and Initialization

```cpp
#include <iostream>
#include <vector>
using namespace std;

// Static array declaration
int staticArray[5];                    // Uninitialized
int initializedArray[5] = {1, 2, 3, 4, 5};
int partialArray[5] = {1, 2};         // {1, 2, 0, 0, 0}
int autoSizeArray[] = {1, 2, 3, 4};   // Size automatically determined

// Character arrays (C-style strings)
char str1[10] = "Hello";
char str2[] = "World";

void arrayBasics() {
    // Array access
    cout << "First element: " << initializedArray[0] << endl;
    cout << "Last element: " << initializedArray[4] << endl;
    
    // Array modification
    initializedArray[0] = 10;
    
    // Array size
    int size = sizeof(initializedArray) / sizeof(initializedArray[0]);
    cout << "Array size: " << size << endl;
}
```

### Array Operations

```cpp
class ArrayOperations {
public:
    // Print array
    static void printArray(int arr[], int size) {
        for (int i = 0; i < size; i++) {
            cout << arr[i] << " ";
        }
        cout << endl;
    }
    
    // Insert element at specific position
    static bool insertAt(int arr[], int& size, int capacity, int pos, int value) {
        if (size >= capacity || pos < 0 || pos > size) {
            return false;
        }
        
        // Shift elements to right
        for (int i = size; i > pos; i--) {
            arr[i] = arr[i-1];
        }
        
        arr[pos] = value;
        size++;
        return true;
    }
    
    // Delete element at specific position
    static bool deleteAt(int arr[], int& size, int pos) {
        if (pos < 0 || pos >= size) {
            return false;
        }
        
        // Shift elements to left
        for (int i = pos; i < size - 1; i++) {
            arr[i] = arr[i + 1];
        }
        
        size--;
        return true;
    }
    
    // Search for element (Linear Search)
    static int linearSearch(int arr[], int size, int target) {
        for (int i = 0; i < size; i++) {
            if (arr[i] == target) {
                return i;
            }
        }
        return -1;
    }
    
    // Binary Search (for sorted arrays)
    static int binarySearch(int arr[], int size, int target) {
        int left = 0, right = size - 1;
        
        while (left <= right) {
            int mid = left + (right - left) / 2;
            
            if (arr[mid] == target) {
                return mid;
            }
            
            if (arr[mid] < target) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
        
        return -1;
    }
};
```

---

## 2. Multi-dimensional Arrays {#multi-dimensional-arrays}

### 2D Arrays

```cpp
// 2D Array declaration
int matrix[3][4];                    // 3 rows, 4 columns
int initialized2D[3][3] = {
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9}
};

class Matrix {
public:
    // Print 2D array
    static void print2D(int arr[][4], int rows, int cols) {
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                cout << arr[i][j] << " ";
            }
            cout << endl;
        }
    }
    
    // Matrix addition
    static void addMatrices(int a[][3], int b[][3], int result[][3], int rows, int cols) {
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                result[i][j] = a[i][j] + b[i][j];
            }
        }
    }
    
    // Matrix multiplication
    static void multiplyMatrices(int a[][3], int b[][3], int result[][3], int r1, int c1, int c2) {
        // Initialize result matrix
        for (int i = 0; i < r1; i++) {
            for (int j = 0; j < c2; j++) {
                result[i][j] = 0;
            }
        }
        
        // Multiply matrices
        for (int i = 0; i < r1; i++) {
            for (int j = 0; j < c2; j++) {
                for (int k = 0; k < c1; k++) {
                    result[i][j] += a[i][k] * b[k][j];
                }
            }
        }
    }
    
    // Transpose matrix
    static void transpose(int matrix[][3], int result[][3], int rows, int cols) {
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                result[j][i] = matrix[i][j];
            }
        }
    }
};
```

### Dynamic 2D Arrays

```cpp
class Dynamic2DArray {
public:
    // Create dynamic 2D array
    static int** create2D(int rows, int cols) {
        int** arr = new int*[rows];
        for (int i = 0; i < rows; i++) {
            arr[i] = new int[cols];
        }
        return arr;
    }
    
    // Delete dynamic 2D array
    static void delete2D(int** arr, int rows) {
        for (int i = 0; i < rows; i++) {
            delete[] arr[i];
        }
        delete[] arr;
    }
    
    // Initialize 2D array
    static void initialize2D(int** arr, int rows, int cols, int value = 0) {
        for (int i = 0; i < rows; i++) {
            for (int j = 0; j < cols; j++) {
                arr[i][j] = value;
            }
        }
    }
};
```

---

## 3. Dynamic Arrays {#dynamic-arrays}

### C++ Vector (Dynamic Array)

```cpp
#include <vector>
#include <algorithm>

class VectorOperations {
public:
    static void vectorBasics() {
        // Vector declaration
        vector<int> v1;                    // Empty vector
        vector<int> v2(5);                 // 5 elements, default value 0
        vector<int> v3(5, 10);             // 5 elements, all set to 10
        vector<int> v4 = {1, 2, 3, 4, 5}; // Initialize with values
        
        // Basic operations
        v1.push_back(10);                  // Add element at end
        v1.push_back(20);
        v1.push_back(30);
        
        cout << "Size: " << v1.size() << endl;
        cout << "Capacity: " << v1.capacity() << endl;
        
        // Access elements
        cout << "First element: " << v1[0] << endl;
        cout << "Last element: " << v1.back() << endl;
        
        // Remove last element
        v1.pop_back();
        
        // Insert at specific position
        v1.insert(v1.begin() + 1, 15);    // Insert 15 at index 1
        
        // Erase element
        v1.erase(v1.begin() + 1);          // Remove element at index 1
    }
    
    static void vectorAlgorithms() {
        vector<int> v = {5, 2, 8, 1, 9, 3};
        
        // Sort vector
        sort(v.begin(), v.end());
        
        // Binary search (vector must be sorted)
        bool found = binary_search(v.begin(), v.end(), 5);
        
        // Find element
        auto it = find(v.begin(), v.end(), 8);
        if (it != v.end()) {
            cout << "Element found at position: " << it - v.begin() << endl;
        }
        
        // Reverse vector
        reverse(v.begin(), v.end());
        
        // Find min and max
        auto minIt = min_element(v.begin(), v.end());
        auto maxIt = max_element(v.begin(), v.end());
        
        cout << "Min: " << *minIt << ", Max: " << *maxIt << endl;
    }
};
```

---

## 4. String Fundamentals {#string-fundamentals}

### C-Style Strings vs C++ Strings

```cpp
#include <string>
#include <cstring>

class StringBasics {
public:
    static void cStyleStrings() {
        // C-style string declaration
        char str1[20] = "Hello";
        char str2[] = "World";
        char str3[20];
        
        // String operations using cstring library
        strcpy(str3, str1);                // Copy string
        strcat(str3, " ");                 // Concatenate
        strcat(str3, str2);
        
        cout << "Length: " << strlen(str3) << endl;
        cout << "String: " << str3 << endl;
        
        // String comparison
        if (strcmp(str1, "Hello") == 0) {
            cout << "Strings are equal" << endl;
        }
    }
    
    static void cppStrings() {
        // C++ string declaration
        string s1 = "Hello";
        string s2 = "World";
        string s3;
        
        // String operations
        s3 = s1 + " " + s2;               // Concatenation
        
        cout << "Length: " << s3.length() << endl;
        cout << "Size: " << s3.size() << endl;
        cout << "String: " << s3 << endl;
        
        // String comparison
        if (s1 == "Hello") {
            cout << "Strings are equal" << endl;
        }
        
        // Access characters
        cout << "First char: " << s3[0] << endl;
        cout << "Last char: " << s3.back() << endl;
        
        // Substring
        string sub = s3.substr(0, 5);     // Extract "Hello"
        cout << "Substring: " << sub << endl;
        
        // Find substring
        size_t pos = s3.find("World");
        if (pos != string::npos) {
            cout << "'World' found at position: " << pos << endl;
        }
    }
};
```

### String Operations

```cpp
class StringOperations {
public:
    // Check if string is palindrome
    static bool isPalindrome(const string& str) {
        int left = 0, right = str.length() - 1;
        
        while (left < right) {
            if (str[left] != str[right]) {
                return false;
            }
            left++;
            right--;
        }
        
        return true;
    }
    
    // Reverse string
    static string reverseString(string str) {
        int left = 0, right = str.length() - 1;
        
        while (left < right) {
            swap(str[left], str[right]);
            left++;
            right--;
        }
        
        return str;
    }
    
    // Count character frequency
    static void charFrequency(const string& str) {
        int freq[256] = {0};  // ASCII characters
        
        for (char c : str) {
            freq[c]++;
        }
        
        for (int i = 0; i < 256; i++) {
            if (freq[i] > 0) {
                cout << (char)i << ": " << freq[i] << endl;
            }
        }
    }
    
    // Remove duplicates from string
    static string removeDuplicates(const string& str) {
        bool seen[256] = {false};
        string result = "";
        
        for (char c : str) {
            if (!seen[c]) {
                seen[c] = true;
                result += c;
            }
        }
        
        return result;
    }
};
```

---

## 5. String Algorithms {#string-algorithms}

### Pattern Matching

```cpp
class PatternMatching {
public:
    // Naive pattern matching
    static vector<int> naiveSearch(const string& text, const string& pattern) {
        vector<int> matches;
        int n = text.length();
        int m = pattern.length();
        
        for (int i = 0; i <= n - m; i++) {
            int j;
            for (j = 0; j < m; j++) {
                if (text[i + j] != pattern[j]) {
                    break;
                }
            }
            
            if (j == m) {
                matches.push_back(i);
            }
        }
        
        return matches;
    }
    
    // KMP Algorithm - Compute LPS array
    static vector<int> computeLPS(const string& pattern) {
        int m = pattern.length();
        vector<int> lps(m, 0);
        int len = 0;
        int i = 1;
        
        while (i < m) {
            if (pattern[i] == pattern[len]) {
                len++;
                lps[i] = len;
                i++;
            } else {
                if (len != 0) {
                    len = lps[len - 1];
                } else {
                    lps[i] = 0;
                    i++;
                }
            }
        }
        
        return lps;
    }
    
    // KMP Algorithm - Pattern matching
    static vector<int> KMPSearch(const string& text, const string& pattern) {
        vector<int> matches;
        int n = text.length();
        int m = pattern.length();
        
        vector<int> lps = computeLPS(pattern);
        
        int i = 0; // index for text
        int j = 0; // index for pattern
        
        while (i < n) {
            if (pattern[j] == text[i]) {
                i++;
                j++;
            }
            
            if (j == m) {
                matches.push_back(i - j);
                j = lps[j - 1];
            } else if (i < n && pattern[j] != text[i]) {
                if (j != 0) {
                    j = lps[j - 1];
                } else {
                    i++;
                }
            }
        }
        
        return matches;
    }
};
```

### String Manipulation

```cpp
class StringManipulation {
public:
    // Convert to lowercase
    static string toLowerCase(string str) {
        for (char& c : str) {
            c = tolower(c);
        }
        return str;
    }
    
    // Convert to uppercase
    static string toUpperCase(string str) {
        for (char& c : str) {
            c = toupper(c);
        }
        return str;
    }
    
    // Trim whitespace
    static string trim(const string& str) {
        size_t start = str.find_first_not_of(" \t\n\r");
        if (start == string::npos) return "";
        
        size_t end = str.find_last_not_of(" \t\n\r");
        return str.substr(start, end - start + 1);
    }
    
    // Split string by delimiter
    static vector<string> split(const string& str, char delimiter) {
        vector<string> tokens;
        string token = "";
        
        for (char c : str) {
            if (c == delimiter) {
                if (!token.empty()) {
                    tokens.push_back(token);
                    token = "";
                }
            } else {
                token += c;
            }
        }
        
        if (!token.empty()) {
            tokens.push_back(token);
        }
        
        return tokens;
    }
    
    // Check if strings are anagrams
    static bool areAnagrams(string str1, string str2) {
        if (str1.length() != str2.length()) {
            return false;
        }
        
        sort(str1.begin(), str1.end());
        sort(str2.begin(), str2.end());
        
        return str1 == str2;
    }
};
```

---

## 6. Common Array Problems {#array-problems}

### Two Pointer Technique

```cpp
class TwoPointerProblems {
public:
    // Two Sum (sorted array)
    static pair<int, int> twoSum(vector<int>& arr, int target) {
        int left = 0, right = arr.size() - 1;
        
        while (left < right) {
            int sum = arr[left] + arr[right];
            
            if (sum == target) {
                return {left, right};
            } else if (sum < target) {
                left++;
            } else {
                right--;
            }
        }
        
        return {-1, -1};
    }
    
    // Remove duplicates from sorted array
    static int removeDuplicates(vector<int>& arr) {
        if (arr.empty()) return 0;
        
        int writeIndex = 1;
        
        for (int i = 1; i < arr.size(); i++) {
            if (arr[i] != arr[i-1]) {
                arr[writeIndex] = arr[i];
                writeIndex++;
            }
        }
        
        return writeIndex;
    }
    
    // Container with most water
    static int maxArea(vector<int>& height) {
        int left = 0, right = height.size() - 1;
        int maxWater = 0;
        
        while (left < right) {
            int width = right - left;
            int currentWater = min(height[left], height[right]) * width;
            maxWater = max(maxWater, currentWater);
            
            if (height[left] < height[right]) {
                left++;
            } else {
                right--;
            }
        }
        
        return maxWater;
    }
};
```

### Sliding Window Technique

```cpp
class SlidingWindowProblems {
public:
    // Maximum sum of k consecutive elements
    static int maxSumSubarray(vector<int>& arr, int k) {
        if (arr.size() < k) return -1;
        
        // Calculate sum of first window
        int windowSum = 0;
        for (int i = 0; i < k; i++) {
            windowSum += arr[i];
        }
        
        int maxSum = windowSum;
        
        // Slide the window
        for (int i = k; i < arr.size(); i++) {
            windowSum = windowSum - arr[i-k] + arr[i];
            maxSum = max(maxSum, windowSum);
        }
        
        return maxSum;
    }
    
    // Longest substring without repeating characters
    static int lengthOfLongestSubstring(string s) {
        unordered_set<char> window;
        int left = 0, maxLength = 0;
        
        for (int right = 0; right < s.length(); right++) {
            while (window.count(s[right])) {
                window.erase(s[left]);
                left++;
            }
            
            window.insert(s[right]);
            maxLength = max(maxLength, right - left + 1);
        }
        
        return maxLength;
    }
};
```

---

## 7. Common String Problems {#string-problems}

### String Matching and Validation

```cpp
class StringProblems {
public:
    // Valid parentheses
    static bool isValidParentheses(const string& s) {
        stack<char> st;
        
        for (char c : s) {
            if (c == '(' || c == '[' || c == '{') {
                st.push(c);
            } else {
                if (st.empty()) return false;
                
                char top = st.top();
                st.pop();
                
                if ((c == ')' && top != '(') ||
                    (c == ']' && top != '[') ||
                    (c == '}' && top != '{')) {
                    return false;
                }
            }
        }
        
        return st.empty();
    }
    
    // Longest common prefix
    static string longestCommonPrefix(vector<string>& strs) {
        if (strs.empty()) return "";
        
        string prefix = strs[0];
        
        for (int i = 1; i < strs.size(); i++) {
            int j = 0;
            while (j < prefix.length() && j < strs[i].length() && 
                   prefix[j] == strs[i][j]) {
                j++;
            }
            prefix = prefix.substr(0, j);
            
            if (prefix.empty()) break;
        }
        
        return prefix;
    }
    
    // Group anagrams
    static vector<vector<string>> groupAnagrams(vector<string>& strs) {
        unordered_map<string, vector<string>> groups;
        
        for (string str : strs) {
            string key = str;
            sort(key.begin(), key.end());
            groups[key].push_back(str);
        }
        
        vector<vector<string>> result;
        for (auto& group : groups) {
            result.push_back(group.second);
        }
        
        return result;
    }
};
```

---

## 8. Practice Exercises {#practice-exercises}

### Array Exercises

1. **Rotate Array**: Rotate array to the right by k steps
2. **Missing Number**: Find missing number in array containing n-1 numbers
3. **Majority Element**: Find element that appears more than n/2 times
4. **Best Time to Buy and Sell Stock**: Find maximum profit
5. **Merge Sorted Arrays**: Merge two sorted arrays in-place
6. **Product of Array Except Self**: Calculate product array without division
7. **Find Peak Element**: Find any peak element in array
8. **Search in Rotated Sorted Array**: Binary search in rotated array

### String Exercises

1. **Reverse Words in String**: Reverse order of words
2. **String to Integer (atoi)**: Implement string to integer conversion
3. **Valid Palindrome**: Check if string is palindrome (ignoring case/spaces)
4. **Implement strStr()**: Find first occurrence of substring
5. **Longest Palindromic Substring**: Find longest palindromic substring
6. **Regular Expression Matching**: Implement regex matching
7. **Edit Distance**: Find minimum operations to convert one string to another
8. **Minimum Window Substring**: Find minimum window containing all characters

### Complete Example Program

```cpp
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <unordered_map>
using namespace std;

class ArrayStringDemo {
public:
    static void demonstrateArrays() {
        cout << "=== Array Operations Demo ===" << endl;
        
        vector<int> arr = {64, 34, 25, 12, 22, 11, 90};
        
        cout << "Original array: ";
        for (int x : arr) cout << x << " ";
        cout << endl;
        
        // Sort array
        sort(arr.begin(), arr.end());
        cout << "Sorted array: ";
        for (int x : arr) cout << x << " ";
        cout << endl;
        
        // Binary search
        int target = 25;
        bool found = binary_search(arr.begin(), arr.end(), target);
        cout << "Element " << target << (found ? " found" : " not found") << endl;
    }
    
    static void demonstrateStrings() {
        cout << "\n=== String Operations Demo ===" << endl;
        
        string str = "Hello World";
        cout << "Original string: " << str << endl;
        
        // Reverse string
        string reversed = str;
        reverse(reversed.begin(), reversed.end());
        cout << "Reversed: " << reversed << endl;
        
        // Check palindrome
        string palindrome = "racecar";
        string temp = palindrome;
        reverse(temp.begin(), temp.end());
        bool isPalin = (palindrome == temp);
        cout << palindrome << " is " << (isPalin ? "" : "not ") << "a palindrome" << endl;
        
        // Character frequency
        unordered_map<char, int> freq;
        for (char c : str) {
            if (c != ' ') freq[c]++;
        }
        
        cout << "Character frequencies:" << endl;
        for (auto& p : freq) {
            cout << p.first << ": " << p.second << endl;
        }
    }
};

int main() {
    ArrayStringDemo::demonstrateArrays();
    ArrayStringDemo::demonstrateStrings();
    
    return 0;
}
```

---

## Summary

In this part, you learned:
- Array fundamentals and operations
- Multi-dimensional arrays and matrices
- Dynamic arrays (vectors) in C++
- String fundamentals and operations
- String algorithms and pattern matching
- Common array and string problems
- Two-pointer and sliding window techniques

**Key Concepts:**
- Arrays provide O(1) access but O(n) insertion/deletion
- Strings are arrays of characters with special operations
- Two-pointer technique is useful for sorted arrays
- Sliding window is efficient for subarray problems
- Pattern matching algorithms like KMP improve search efficiency

**Next**: [Part 3: Linked Lists](./part3-linked-lists.md)

---

## Additional Practice Resources

- LeetCode Array and String problems
- HackerRank Data Structures section
- GeeksforGeeks Array and String articles
- Practice implementing all algorithms from scratch
- Time complexity analysis for each solution