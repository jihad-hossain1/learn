# Comprehensive DSA Learning Guide: JavaScript Edition

## 🎯 Learning Path Overview

### **Phase 1: Fundamentals (2-3 weeks)**
**Prerequisites: Basic JavaScript syntax**

#### **Week 1: JavaScript Basics & Simple Data Structures**
```javascript
// Arrays
const arr = [1, 2, 3, 4, 5];
arr.push(6); // O(1)
arr.pop();   // O(1)
arr.unshift(0); // O(n)
arr.splice(2, 1); // O(n)

// Strings
const str = "hello";
str.length; // O(1)
str.includes("ell"); // O(n)

// Objects
const obj = { a: 1, b: 2 };
obj.c = 3; // O(1)
delete obj.a; // O(1)
```

**Practice Problems:**
- Reverse a string/array
- Find max/min in array
- Count characters in string
- Two sum (brute force)

---

### **Phase 2: Core Data Structures (4-6 weeks)**

#### **Week 2-3: Arrays & Strings Deep Dive**
```javascript
// Two Pointers Technique
function twoSumSorted(arr, target) {
    let left = 0, right = arr.length - 1;
    while (left < right) {
        const sum = arr[left] + arr[right];
        if (sum === target) return [left, right];
        if (sum < target) left++;
        else right--;
    }
    return [-1, -1];
}

// Sliding Window
function maxSubarraySum(arr, k) {
    let maxSum = 0, windowSum = 0;
    for (let i = 0; i < k; i++) windowSum += arr[i];
    
    for (let i = k; i < arr.length; i++) {
        windowSum += arr[i] - arr[i - k];
        maxSum = Math.max(maxSum, windowSum);
    }
    return maxSum;
}
```

**Practice Problems:**
- Move zeros
- Container with most water
- Longest substring without repeating characters
- Valid palindrome

#### **Week 4-5: Linked Lists**
```javascript
class ListNode {
    constructor(val, next = null) {
        this.val = val;
        this.next = next;
    }
}

class LinkedList {
    constructor() {
        this.head = null;
    }
    
    // Add to end - O(n)
    append(val) {
        if (!this.head) {
            this.head = new ListNode(val);
            return;
        }
        let current = this.head;
        while (current.next) current = current.next;
        current.next = new ListNode(val);
    }
    
    // Reverse linked list - O(n)
    reverse() {
        let prev = null, current = this.head;
        while (current) {
            const next = current.next;
            current.next = prev;
            prev = current;
            current = next;
        }
        this.head = prev;
    }
}
```

**Practice Problems:**
- Detect cycle in linked list
- Merge two sorted lists
- Remove nth node from end
- Palindrome linked list

#### **Week 6-7: Stacks & Queues**
```javascript
// Stack implementation
class Stack {
    constructor() {
        this.items = [];
    }
    
    push(item) { this.items.push(item); }
    pop() { return this.items.pop(); }
    peek() { return this.items[this.items.length - 1]; }
    isEmpty() { return this.items.length === 0; }
}

// Queue implementation
class Queue {
    constructor() {
        this.items = [];
    }
    
    enqueue(item) { this.items.push(item); }
    dequeue() { return this.items.shift(); } // O(n) - can be optimized
    front() { return this.items[0]; }
    isEmpty() { return this.items.length === 0; }
}

// Monotonic stack example
function nextGreaterElement(nums) {
    const stack = [];
    const result = new Array(nums.length).fill(-1);
    
    for (let i = 0; i < nums.length; i++) {
        while (stack.length && nums[stack[stack.length - 1]] < nums[i]) {
            const index = stack.pop();
            result[index] = nums[i];
        }
        stack.push(i);
    }
    return result;
}
```

**Practice Problems:**
- Valid parentheses
- Implement queue using stacks
- Daily temperatures
- Min stack

---

### **Phase 3: Intermediate Algorithms (4-5 weeks)**

#### **Week 8-9: Recursion & Backtracking**
```javascript
// Basic recursion patterns
function factorial(n) {
    if (n <= 1) return 1; // base case
    return n * factorial(n - 1); // recursive case
}

// Backtracking template
function backtrack(candidates, target, start = 0, path = [], result = []) {
    if (target === 0) {
        result.push([...path]);
        return;
    }
    
    for (let i = start; i < candidates.length; i++) {
        if (target - candidates[i] >= 0) {
            path.push(candidates[i]);
            backtrack(candidates, target - candidates[i], i, path, result);
            path.pop(); // backtrack
        }
    }
    return result;
}
```

**Practice Problems:**
- Generate parentheses
- Subsets
- Permutations
- Combination sum

#### **Week 10-11: Trees & Binary Search Trees**
```javascript
class TreeNode {
    constructor(val, left = null, right = null) {
        this.val = val;
        this.left = left;
        this.right = right;
    }
}

// Tree traversals
function inOrderTraversal(root, result = []) {
    if (!root) return result;
    inOrderTraversal(root.left, result);
    result.push(root.val);
    inOrderTraversal(root.right, result);
    return result;
}

// BST validation
function isValidBST(root, min = -Infinity, max = Infinity) {
    if (!root) return true;
    if (root.val <= min || root.val >= max) return false;
    return isValidBST(root.left, min, root.val) && 
           isValidBST(root.right, root.val, max);
}
```

**Practice Problems:**
- Maximum depth of binary tree
- Invert binary tree
- Level order traversal
- Kth smallest element in BST

---

### **Phase 4: Advanced Data Structures (4-5 weeks)**

#### **Week 12-13: Hash Tables & Heaps**
```javascript
// Hash table implementation
class HashTable {
    constructor(size = 53) {
        this.keyMap = new Array(size);
    }
    
    _hash(key) {
        let total = 0, WEIRD_PRIME = 31;
        for (let i = 0; i < Math.min(key.length, 100); i++) {
            const char = key[i];
            const value = char.charCodeAt(0) - 96;
            total = (total * WEIRD_PRIME + value) % this.keyMap.length;
        }
        return total;
    }
    
    set(key, value) {
        const index = this._hash(key);
        if (!this.keyMap[index]) this.keyMap[index] = [];
        this.keyMap[index].push([key, value]);
    }
}

// Min heap implementation
class MinHeap {
    constructor() {
        this.heap = [];
    }
    
    insert(val) {
        this.heap.push(val);
        this.bubbleUp();
    }
    
    bubbleUp() {
        let index = this.heap.length - 1;
        while (index > 0) {
            const parent = Math.floor((index - 1) / 2);
            if (this.heap[parent] <= this.heap[index]) break;
            [this.heap[parent], this.heap[index]] = [this.heap[index], this.heap[parent]];
            index = parent;
        }
    }
}
```

**Practice Problems:**
- First unique character
- Group anagrams
- Top K frequent elements
- Find median from data stream

#### **Week 14-15: Graphs**
```javascript
// Graph implementation
class Graph {
    constructor() {
        this.adjacencyList = {};
    }
    
    addVertex(vertex) {
        if (!this.adjacencyList[vertex]) this.adjacencyList[vertex] = [];
    }
    
    addEdge(v1, v2) {
        this.adjacencyList[v1].push(v2);
        this.adjacencyList[v2].push(v1);
    }
    
    // BFS
    bfs(start) {
        const queue = [start], result = [], visited = {[start]: true};
        
        while (queue.length) {
            const vertex = queue.shift();
            result.push(vertex);
            
            this.adjacencyList[vertex].forEach(neighbor => {
                if (!visited[neighbor]) {
                    visited[neighbor] = true;
                    queue.push(neighbor);
                }
            });
        }
        return result;
    }
    
    // DFS
    dfs(start) {
        const result = [], visited = {};
        
        const dfsHelper = (vertex) => {
            if (!vertex) return null;
            visited[vertex] = true;
            result.push(vertex);
            this.adjacencyList[vertex].forEach(neighbor => {
                if (!visited[neighbor]) dfsHelper(neighbor);
            });
        };
        
        dfsHelper(start);
        return result;
    }
}
```

**Practice Problems:**
- Number of islands
- Clone graph
- Course schedule
- Word ladder

---

### **Phase 5: Advanced Algorithms (4-6 weeks)**

#### **Week 16-18: Dynamic Programming**
```javascript
// Fibonacci with memoization
function fibonacci(n, memo = {}) {
    if (n in memo) return memo[n];
    if (n <= 2) return 1;
    memo[n] = fibonacci(n - 1, memo) + fibonacci(n - 2, memo);
    return memo[n];
}

// Tabulation approach
function fibTabulation(n) {
    if (n <= 2) return 1;
    const dp = [0, 1, 1];
    for (let i = 3; i <= n; i++) {
        dp[i] = dp[i - 1] + dp[i - 2];
    }
    return dp[n];
}

// Knapsack problem
function knapsack(weights, values, capacity) {
    const n = weights.length;
    const dp = Array(n + 1).fill().map(() => Array(capacity + 1).fill(0));
    
    for (let i = 1; i <= n; i++) {
        for (let w = 1; w <= capacity; w++) {
            if (weights[i - 1] <= w) {
                dp[i][w] = Math.max(
                    values[i - 1] + dp[i - 1][w - weights[i - 1]],
                    dp[i - 1][w]
                );
            } else {
                dp[i][w] = dp[i - 1][w];
            }
        }
    }
    return dp[n][capacity];
}
```

**Practice Problems:**
- Climbing stairs
- Coin change
- Longest common subsequence
- House robber

#### **Week 19-20: Sorting & Searching Algorithms**
```javascript
// Quick sort
function quickSort(arr) {
    if (arr.length <= 1) return arr;
    
    const pivot = arr[0];
    const left = [], right = [];
    
    for (let i = 1; i < arr.length; i++) {
        if (arr[i] < pivot) left.push(arr[i]);
        else right.push(arr[i]);
    }
    
    return [...quickSort(left), pivot, ...quickSort(right)];
}

// Binary search variations
function binarySearch(arr, target) {
    let left = 0, right = arr.length - 1;
    
    while (left <= right) {
        const mid = Math.floor((left + right) / 2);
        if (arr[mid] === target) return mid;
        if (arr[mid] < target) left = mid + 1;
        else right = mid - 1;
    }
    return -1;
}
```

**Practice Problems:**
- Merge intervals
- Search in rotated sorted array
- Find peak element
- K closest points to origin

---

## 🛠️ Practice Platforms & Resources

### **Recommended Platforms:**
- **LeetCode** (start with Easy, progress to Medium)
- **HackerRank** (good for fundamentals)
- **AlgoExpert** (JavaScript-focused)
- **NeetCode.io** (curated problem lists)

### **Learning Resources:**
- **Books:** "JavaScript Data Structures and Algorithms" by Sammie Bae
- **YouTube:** FreeCodeCamp, NeetCode, CS Dojo
- **Courses:** Udemy "JavaScript Algorithms and Data Structures Masterclass"

## 📈 Study Strategy

### **Daily Routine:**
1. **Morning (30 mins):** Review one concept
2. **Afternoon (1 hour):** Solve 2-3 problems
3. **Evening (30 mins):** Review solutions, optimize code

### **Weekly Goals:**
- Complete 10-15 problems
- Master 1-2 new patterns
- Review previous week's topics

### **Tracking Progress:**
```javascript
// Keep a learning log
const progress = {
    week1: { topics: ["Arrays", "Strings"], problemsSolved: 15 },
    week2: { topics: ["Linked Lists"], problemsSolved: 12 },
    // ... continue tracking
};
```

## 🎯 Final Tips

1. **Consistency > Intensity:** 1 hour daily is better than 7 hours once a week
2. **Understand, Don't Memorize:** Focus on patterns and reasoning
3. **Time & Space Complexity:** Always analyze your solutions
4. **Practice System Design:** Many interviews include this aspect
5. **Mock Interviews:** Practice explaining your thought process

Remember: This is a 20-25 week journey. Adjust the pace based on your background and learning style. Good luck! 🚀