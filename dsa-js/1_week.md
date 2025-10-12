# Practice Problems: JavaScript Fundamentals

Let me explain each of these fundamental problems with detailed solutions and analysis.

## 1. Reverse a String/Array

### Problem Understanding:
Reverse the order of elements in a string or array.

### Multiple Approaches:

```javascript
// Method 1: Built-in methods (Easiest)
function reverseStringBuiltIn(str) {
    return str.split('').reverse().join('');
}

function reverseArrayBuiltIn(arr) {
    return [...arr].reverse(); // Spread operator to avoid mutating original
}

// Method 2: Two Pointers (Most Efficient - O(n) time, O(1) space)
function reverseStringTwoPointers(str) {
    const chars = str.split('');
    let left = 0, right = chars.length - 1;
    
    while (left < right) {
        // Swap characters
        [chars[left], chars[right]] = [chars[right], chars[left]];
        left++;
        right--;
    }
    return chars.join('');
}

function reverseArrayTwoPointers(arr) {
    const result = [...arr]; // Create copy to avoid mutation
    let left = 0, right = result.length - 1;
    
    while (left < right) {
        [result[left], result[right]] = [result[right], result[left]];
        left++;
        right--;
    }
    return result;
}

// Method 3: Using Stack (Educational)
function reverseStringStack(str) {
    const stack = [];
    // Push all characters onto stack
    for (let char of str) {
        stack.push(char);
    }
    
    // Pop from stack (reverses order)
    let reversed = '';
    while (stack.length > 0) {
        reversed += stack.pop();
    }
    return reversed;
}

// Test cases
console.log(reverseStringBuiltIn("hello")); // "olleh"
console.log(reverseArrayBuiltIn([1, 2, 3, 4])); // [4, 3, 2, 1]
console.log(reverseStringTwoPointers("javascript")); // "tpircsavaj"
```

**Time Complexity:** O(n)  
**Space Complexity:** O(n) for new string/array, O(1) for in-place with two pointers

---

## 2. Find Max/Min in Array

### Problem Understanding:
Find the maximum and minimum values in an array.

### Solutions:

```javascript
// Method 1: Built-in methods (Simplest)
function findMaxMinBuiltIn(arr) {
    return {
        max: Math.max(...arr),
        min: Math.min(...arr)
    };
}

// Method 2: Linear Search (Most Efficient - O(n))
function findMaxMinLinear(arr) {
    if (arr.length === 0) return { max: undefined, min: undefined };
    
    let max = arr[0];
    let min = arr[0];
    
    for (let i = 1; i < arr.length; i++) {
        if (arr[i] > max) max = arr[i];
        if (arr[i] < min) min = arr[i];
    }
    
    return { max, min };
}

// Method 3: Using Reduce (Functional Programming)
function findMaxMinReduce(arr) {
    if (arr.length === 0) return { max: undefined, min: undefined };
    
    return arr.reduce((acc, current) => {
        return {
            max: current > acc.max ? current : acc.max,
            min: current < acc.min ? current : acc.min
        };
    }, { max: arr[0], min: arr[0] });
}

// Method 4: Tournament Method (Divide & Conquer - O(n) but fewer comparisons)
function findMaxMinTournament(arr, low = 0, high = arr.length - 1) {
    // Base cases
    if (low === high) {
        return { max: arr[low], min: arr[low] };
    }
    
    if (high === low + 1) {
        return {
            max: Math.max(arr[low], arr[high]),
            min: Math.min(arr[low], arr[high])
        };
    }
    
    // Divide
    const mid = Math.floor((low + high) / 2);
    const left = findMaxMinTournament(arr, low, mid);
    const right = findMaxMinTournament(arr, mid + 1, high);
    
    // Combine
    return {
        max: Math.max(left.max, right.max),
        min: Math.min(left.min, right.min)
    };
}

// Test cases
const testArray = [3, 1, 4, 1, 5, 9, 2, 6];
console.log(findMaxMinBuiltIn(testArray)); // { max: 9, min: 1 }
console.log(findMaxMinLinear(testArray));  // { max: 9, min: 1 }
console.log(findMaxMinReduce(testArray));  // { max: 9, min: 1 }
```

**Time Complexity:** O(n) for all approaches  
**Space Complexity:** O(1) for linear search, O(log n) for tournament method (recursion stack)

---

## 3. Count Characters in String

### Problem Understanding:
Count the frequency of each character in a string.

### Solutions:

```javascript
// Method 1: Using Object (Most Common)
function countCharactersObject(str) {
    const charCount = {};
    
    for (let char of str) {
        // If character exists, increment count, else initialize to 1
        charCount[char] = (charCount[char] || 0) + 1;
    }
    
    return charCount;
}

// Method 2: Using Map (Better for non-string keys)
function countCharactersMap(str) {
    const charCount = new Map();
    
    for (let char of str) {
        charCount.set(char, (charCount.get(char) || 0) + 1);
    }
    
    return charCount;
}

// Method 3: Using Reduce (Functional Approach)
function countCharactersReduce(str) {
    return str.split('').reduce((count, char) => {
        count[char] = (count[char] || 0) + 1;
        return count;
    }, {});
}

// Method 4: Only Alphabets (Case Insensitive)
function countAlphabetsOnly(str) {
    const charCount = {};
    
    for (let char of str.toLowerCase()) {
        if (char >= 'a' && char <= 'z') {
            charCount[char] = (charCount[char] || 0) + 1;
        }
    }
    
    return charCount;
}

// Method 5: Get Most Frequent Character
function getMostFrequentChar(str) {
    const charCount = countCharactersObject(str);
    let maxChar = '';
    let maxCount = 0;
    
    for (let char in charCount) {
        if (charCount[char] > maxCount) {
            maxCount = charCount[char];
            maxChar = char;
        }
    }
    
    return { character: maxChar, count: maxCount };
}

// Test cases
const testString = "hello world";
console.log(countCharactersObject(testString));
// Output: { h: 1, e: 1, l: 3, o: 2, ' ': 1, w: 1, r: 1, d: 1 }

console.log(getMostFrequentChar(testString));
// Output: { character: 'l', count: 3 }

console.log(countAlphabetsOnly("Hello World! 123"));
// Output: { h: 1, e: 1, l: 3, o: 2, w: 1, r: 1, d: 1 }
```

**Time Complexity:** O(n) where n is string length  
**Space Complexity:** O(k) where k is number of unique characters

---

## 4. Two Sum (Brute Force)

### Problem Understanding:
Given an array of numbers and a target sum, find two numbers that add up to the target. Return their indices.

### Solutions:

```javascript
// Method 1: Brute Force (Naive Approach - O(n²))
function twoSumBruteForce(nums, target) {
    for (let i = 0; i < nums.length; i++) {
        for (let j = i + 1; j < nums.length; j++) {
            if (nums[i] + nums[j] === target) {
                return [i, j];
            }
        }
    }
    return []; // No solution found
}

// Method 2: Using Hash Map (Optimized - O(n))
function twoSumHashMap(nums, target) {
    const numMap = new Map();
    
    for (let i = 0; i < nums.length; i++) {
        const complement = target - nums[i];
        
        if (numMap.has(complement)) {
            return [numMap.get(complement), i];
        }
        
        numMap.set(nums[i], i);
    }
    return [];
}

// Method 3: Using Object (Similar to Map)
function twoSumObject(nums, target) {
    const numObj = {};
    
    for (let i = 0; i < nums.length; i++) {
        const complement = target - nums[i];
        
        if (complement in numObj) {
            return [numObj[complement], i];
        }
        
        numObj[nums[i]] = i;
    }
    return [];
}

// Method 4: Two Pointers (For Sorted Array - O(n log n) with sorting)
function twoSumTwoPointers(nums, target) {
    // Create array of [value, originalIndex] pairs
    const indexedNums = nums.map((num, index) => [num, index]);
    
    // Sort by value
    indexedNums.sort((a, b) => a[0] - b[0]);
    
    let left = 0, right = indexedNums.length - 1;
    
    while (left < right) {
        const sum = indexedNums[left][0] + indexedNums[right][0];
        
        if (sum === target) {
            return [indexedNums[left][1], indexedNums[right][1]];
        } else if (sum < target) {
            left++;
        } else {
            right--;
        }
    }
    return [];
}

// Test cases
const nums = [2, 7, 11, 15];
const target = 9;

console.log(twoSumBruteForce(nums, target)); // [0, 1]
console.log(twoSumHashMap(nums, target));    // [0, 1]
console.log(twoSumObject(nums, target));     // [0, 1]
console.log(twoSumTwoPointers(nums, target)); // [0, 1]

// Edge cases
console.log(twoSumBruteForce([3, 3], 6)); // [0, 1]
console.log(twoSumBruteForce([3, 2, 4], 6)); // [1, 2]
```

### Complexity Analysis:

**Brute Force:**
- **Time:** O(n²) - Nested loops
- **Space:** O(1) - No extra space

**HashMap/Object:**
- **Time:** O(n) - Single pass
- **Space:** O(n) - Store complements

**Two Pointers:**
- **Time:** O(n log n) - Due to sorting
- **Space:** O(n) - Store indexed array

---

## 🎯 Key Learning Points

### **Reverse Problems:**
- Understand different approaches (built-in, two pointers, stack)
- Learn about in-place vs new array/string creation
- Practice swapping elements

### **Max/Min Problems:**
- Learn to traverse arrays efficiently
- Understand edge cases (empty array, single element)
- Compare different algorithmic approaches

### **Character Counting:**
- Practice with different data structures (Object, Map)
- Learn string manipulation
- Handle edge cases (case sensitivity, special characters)

### **Two Sum:**
- Understand time-space tradeoffs
- Learn hash-based optimization
- Practice multiple approaches to the same problem

### **Next Steps:**
1. Solve each problem manually first
2. Implement multiple solutions
3. Analyze time/space complexity
4. Test with edge cases
5. Practice explaining your solutions

These fundamental problems build the foundation for more complex DSA concepts! 🚀