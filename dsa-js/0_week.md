# Practice Problems: JavaScript Fundamentals - Detailed Explanation & Practice

Let me provide a comprehensive explanation and practice for these fundamental problems that form the building blocks of algorithmic thinking.

## 1. Reverse a String/Array

### Problem Understanding:
Reverse the order of elements in a string or array.

**Key Concepts:**
- String/Array manipulation
- Two-pointer technique
- In-place vs new collection

### Detailed Solutions:

```javascript
// Method 1: Built-in Methods (Simplest)
function reverseStringBuiltIn(str) {
    return str.split('').reverse().join('');
}

function reverseArrayBuiltIn(arr) {
    return [...arr].reverse(); // Spread to avoid mutating original
}

// Method 2: Two Pointers (Most Efficient)
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
    const result = [...arr]; // Create copy
    let left = 0, right = result.length - 1;
    
    while (left < right) {
        [result[left], result[right]] = [result[right], result[left]];
        left++;
        right--;
    }
    return result;
}

// Method 3: Using Stack (Educational Purpose)
function reverseStringStack(str) {
    const stack = [];
    // Push all characters onto stack
    for (let char of str) {
        stack.push(char);
    }
    
    // Pop from stack (reverses order due to LIFO)
    let reversed = '';
    while (stack.length > 0) {
        reversed += stack.pop();
    }
    return reversed;
}

// Method 4: Recursive Approach
function reverseStringRecursive(str) {
    if (str === '') return '';
    return reverseStringRecursive(str.substr(1)) + str[0];
}

// Method 5: Using Reduce
function reverseStringReduce(str) {
    return str.split('').reduce((reversed, char) => char + reversed, '');
}

function reverseArrayReduce(arr) {
    return arr.reduce((reversed, item) => [item, ...reversed], []);
}

// Enhanced Practice with Error Handling
function enhancedReverse(input) {
    // Handle both strings and arrays
    if (typeof input === 'string') {
        return reverseStringTwoPointers(input);
    } else if (Array.isArray(input)) {
        return reverseArrayTwoPointers(input);
    } else {
        throw new Error('Input must be string or array');
    }
}

// Test Cases
console.log("=== Reverse String/Array Tests ===");
console.log("String 'hello':", reverseStringBuiltIn("hello")); // "olleh"
console.log("Array [1,2,3]:", reverseArrayBuiltIn([1,2,3])); // [3,2,1]
console.log("Two Pointers 'javascript':", reverseStringTwoPointers("javascript")); // "tpircsavaj"
console.log("Recursive 'recursion':", reverseStringRecursive("recursion")); // "noisrucer"
console.log("Enhanced Reverse:", enhancedReverse("test")); // "tset"

// Performance Comparison
const testString = "a".repeat(1000);
console.time('Built-in Reverse');
reverseStringBuiltIn(testString);
console.timeEnd('Built-in Reverse');

console.time('Two Pointers Reverse');
reverseStringTwoPointers(testString);
console.timeEnd('Two Pointers Reverse');
```

### Practice Exercises:
```javascript
// Exercise 1: Reverse words in a sentence
function reverseWords(sentence) {
    // "hello world" → "world hello"
    // Your code here
    return sentence.split(' ').reverse().join(' ');
}

// Exercise 2: Reverse only vowels in a string
function reverseVowels(str) {
    // "hello" → "holle"
    // Your code here
    const vowels = 'aeiouAEIOU';
    const chars = str.split('');
    let left = 0, right = chars.length - 1;
    
    while (left < right) {
        if (!vowels.includes(chars[left])) {
            left++;
        } else if (!vowels.includes(chars[right])) {
            right--;
        } else {
            [chars[left], chars[right]] = [chars[right], chars[left]];
            left++;
            right--;
        }
    }
    return chars.join('');
}

// Test your solutions
console.log(reverseWords("hello world")); // "world hello"
console.log(reverseVowels("hello")); // "holle"
```

---

## 2. Find Max/Min in Array

### Problem Understanding:
Find the maximum and minimum values in an array.

**Key Concepts:**
- Array traversal
- Comparison operations
- Edge cases handling

### Detailed Solutions:

```javascript
// Method 1: Built-in Methods (Simplest)
function findMaxMinBuiltIn(arr) {
    if (arr.length === 0) return { max: undefined, min: undefined };
    
    return {
        max: Math.max(...arr),
        min: Math.min(...arr)
    };
}

// Method 2: Linear Search (Most Efficient)
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

// Method 4: Tournament Method (Divide & Conquer)
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

// Method 5: One Pass with Pair Comparison
function findMaxMinPairs(arr) {
    if (arr.length === 0) return { max: undefined, min: undefined };
    if (arr.length === 1) return { max: arr[0], min: arr[0] };
    
    let max, min;
    let i;
    
    // Initialize first two elements
    if (arr[0] > arr[1]) {
        max = arr[0];
        min = arr[1];
    } else {
        max = arr[1];
        min = arr[0];
    }
    
    // Process remaining elements in pairs
    for (i = 2; i < arr.length - 1; i += 2) {
        if (arr[i] > arr[i + 1]) {
            max = Math.max(max, arr[i]);
            min = Math.min(min, arr[i + 1]);
        } else {
            max = Math.max(max, arr[i + 1]);
            min = Math.min(min, arr[i]);
        }
    }
    
    // Handle odd-length array
    if (i === arr.length - 1) {
        max = Math.max(max, arr[i]);
        min = Math.min(min, arr[i]);
    }
    
    return { max, min };
}

// Enhanced version with statistics
function arrayStats(arr) {
    if (arr.length === 0) {
        return {
            max: undefined,
            min: undefined,
            average: undefined,
            sum: 0,
            count: 0
        };
    }
    
    let max = arr[0];
    let min = arr[0];
    let sum = 0;
    
    for (let i = 0; i < arr.length; i++) {
        if (arr[i] > max) max = arr[i];
        if (arr[i] < min) min = arr[i];
        sum += arr[i];
    }
    
    return {
        max,
        min,
        average: sum / arr.length,
        sum,
        count: arr.length
    };
}

// Test Cases
console.log("\n=== Find Max/Min Tests ===");
const testArray = [3, 1, 4, 1, 5, 9, 2, 6];
console.log("Linear:", findMaxMinLinear(testArray)); // { max: 9, min: 1 }
console.log("Reduce:", findMaxMinReduce(testArray)); // { max: 9, min: 1 }
console.log("Tournament:", findMaxMinTournament(testArray)); // { max: 9, min: 1 }
console.log("Pairs:", findMaxMinPairs(testArray)); // { max: 9, min: 1 }
console.log("Array Stats:", arrayStats(testArray));

// Edge Cases
console.log("Empty array:", findMaxMinLinear([])); // { max: undefined, min: undefined }
console.log("Single element:", findMaxMinLinear([5])); // { max: 5, min: 5 }
```

### Practice Exercises:
```javascript
// Exercise 1: Find second largest element
function secondLargest(arr) {
    if (arr.length < 2) return undefined;
    
    let largest = -Infinity;
    let secondLargest = -Infinity;
    
    for (let num of arr) {
        if (num > largest) {
            secondLargest = largest;
            largest = num;
        } else if (num > secondLargest && num < largest) {
            secondLargest = num;
        }
    }
    
    return secondLargest !== -Infinity ? secondLargest : undefined;
}

// Exercise 2: Find max and min with their indices
function findMaxMinWithIndices(arr) {
    if (arr.length === 0) return { max: undefined, min: undefined };
    
    let max = { value: arr[0], index: 0 };
    let min = { value: arr[0], index: 0 };
    
    for (let i = 1; i < arr.length; i++) {
        if (arr[i] > max.value) {
            max.value = arr[i];
            max.index = i;
        }
        if (arr[i] < min.value) {
            min.value = arr[i];
            min.index = i;
        }
    }
    
    return { max, min };
}

// Test your solutions
console.log(secondLargest([3, 1, 4, 1, 5, 9, 2, 6])); // 6
console.log(findMaxMinWithIndices([3, 1, 4, 1, 5, 9, 2, 6]));
// { max: { value: 9, index: 5 }, min: { value: 1, index: 1 } }
```

---

## 3. Count Characters in String

### Problem Understanding:
Count the frequency of each character in a string.

**Key Concepts:**
- Hash maps/objects for frequency counting
- String iteration
- Character validation

### Detailed Solutions:

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

// Method 5: Using Array for ASCII (Optimized for English)
function countCharactersASCII(str) {
    const count = new Array(128).fill(0); // For ASCII characters
    
    for (let i = 0; i < str.length; i++) {
        const charCode = str.charCodeAt(i);
        if (charCode < 128) {
            count[charCode]++;
        }
    }
    
    // Convert to readable object
    const result = {};
    for (let i = 0; i < 128; i++) {
        if (count[i] > 0) {
            result[String.fromCharCode(i)] = count[i];
        }
    }
    
    return result;
}

// Enhanced version with analysis
function analyzeString(str) {
    const charCount = countCharactersObject(str);
    let mostFrequent = { char: '', count: 0 };
    let leastFrequent = { char: '', count: Infinity };
    let totalChars = 0;
    let uniqueChars = 0;
    
    for (let char in charCount) {
        const count = charCount[char];
        totalChars += count;
        uniqueChars++;
        
        if (count > mostFrequent.count) {
            mostFrequent = { char, count };
        }
        if (count < leastFrequent.count) {
            leastFrequent = { char, count };
        }
    }
    
    return {
        charCount,
        mostFrequent,
        leastFrequent,
        totalChars,
        uniqueChars,
        averageFrequency: totalChars / uniqueChars
    };
}

// Test Cases
console.log("\n=== Count Characters Tests ===");
const testString = "hello world";
console.log("Object method:", countCharactersObject(testString));
// { h: 1, e: 1, l: 3, o: 2, ' ': 1, w: 1, r: 1, d: 1 }

console.log("Map method:", countCharactersMap(testString));
console.log("Alphabets only:", countAlphabetsOnly("Hello World! 123"));
// { h: 1, e: 1, l: 3, o: 2, w: 1, r: 1, d: 1 }

console.log("String analysis:", analyzeString("hello"));

// Performance Test
const longString = "a".repeat(1000) + "b".repeat(500) + "c".repeat(250);
console.time('Object Counting');
countCharactersObject(longString);
console.timeEnd('Object Counting');

console.time('Map Counting');
countCharactersMap(longString);
console.timeEnd('Map Counting');
```

### Practice Exercises:
```javascript
// Exercise 1: Find first non-repeating character
function firstNonRepeatingChar(str) {
    const charCount = countCharactersObject(str);
    
    for (let char of str) {
        if (charCount[char] === 1) {
            return char;
        }
    }
    return null;
}

// Exercise 2: Check if string is anagram
function isAnagram(str1, str2) {
    if (str1.length !== str2.length) return false;
    
    const count1 = countCharactersObject(str1);
    const count2 = countCharactersObject(str2);
    
    for (let char in count1) {
        if (count1[char] !== count2[char]) {
            return false;
        }
    }
    return true;
}

// Exercise 3: Count only vowels and consonants
function countVowelsConsonants(str) {
    const vowels = 'aeiouAEIOU';
    let vowelCount = 0;
    let consonantCount = 0;
    
    for (let char of str) {
        if (char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z') {
            if (vowels.includes(char)) {
                vowelCount++;
            } else {
                consonantCount++;
            }
        }
    }
    
    return { vowels: vowelCount, consonants: consonantCount };
}

// Test your solutions
console.log(firstNonRepeatingChar("swiss")); // "w"
console.log(isAnagram("listen", "silent")); // true
console.log(countVowelsConsonants("Hello World!")); // { vowels: 3, consonants: 7 }
```

---

## 4. Two Sum (Brute Force)

### Problem Understanding:
Given an array of numbers and a target sum, find two numbers that add up to the target. Return their indices.

**Key Concepts:**
- Nested loops
- Array indexing
- Early termination

### Detailed Solutions:

```javascript
// Method 1: Brute Force (Naive Approach)
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

// Method 2: Enhanced Brute Force with Early Return
function twoSumBruteForceEnhanced(nums, target) {
    for (let i = 0; i < nums.length - 1; i++) {
        for (let j = i + 1; j < nums.length; j++) {
            if (nums[i] + nums[j] === target) {
                return [i, j];
            }
        }
    }
    return [];
}

// Method 3: Brute Force with All Solutions
function twoSumAllSolutions(nums, target) {
    const solutions = [];
    
    for (let i = 0; i < nums.length; i++) {
        for (let j = i + 1; j < nums.length; j++) {
            if (nums[i] + nums[j] === target) {
                solutions.push([i, j]);
            }
        }
    }
    
    return solutions;
}

// Method 4: Using Hash Map (For Comparison - More Efficient)
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

// Enhanced Brute Force with Debugging
function twoSumWithSteps(nums, target) {
    console.log(`Looking for two numbers that sum to ${target}`);
    console.log(`Array: [${nums}]`);
    
    for (let i = 0; i < nums.length; i++) {
        for (let j = i + 1; j < nums.length; j++) {
            const sum = nums[i] + nums[j];
            console.log(`Checking nums[${i}]=${nums[i]} + nums[${j}]=${nums[j]} = ${sum}`);
            
            if (sum === target) {
                console.log(`✓ Found solution: indices ${i} and ${j}`);
                return [i, j];
            }
        }
    }
    
    console.log('✗ No solution found');
    return [];
}

// Test Cases
console.log("\n=== Two Sum Tests ===");
const nums = [2, 7, 11, 15];
const target = 9;

console.log("Brute Force:", twoSumBruteForce(nums, target)); // [0, 1]
console.log("All Solutions:", twoSumAllSolutions([1, 2, 3, 4, 5], 6)); // [[0,4], [1,3]]

// Edge Cases
console.log("Empty array:", twoSumBruteForce([], 9)); // []
console.log("Single element:", twoSumBruteForce([5], 5)); // []
console.log("No solution:", twoSumBruteForce([1, 2, 3], 10)); // []
console.log("Duplicate elements:", twoSumBruteForce([3, 3], 6)); // [0, 1]

// Performance Comparison
const largeArray = Array.from({ length: 1000 }, (_, i) => i);
console.time('Brute Force');
twoSumBruteForce(largeArray, 1997);
console.timeEnd('Brute Force');

console.time('Hash Map');
twoSumHashMap(largeArray, 1997);
console.timeEnd('Hash Map');
```

### Practice Exercises:
```javascript
// Exercise 1: Two Sum with sorted array (two pointers)
function twoSumSorted(nums, target) {
    let left = 0, right = nums.length - 1;
    
    while (left < right) {
        const sum = nums[left] + nums[right];
        
        if (sum === target) {
            return [left, right];
        } else if (sum < target) {
            left++;
        } else {
            right--;
        }
    }
    
    return [];
}

// Exercise 2: Three Sum (brute force)
function threeSumBruteForce(nums, target) {
    for (let i = 0; i < nums.length; i++) {
        for (let j = i + 1; j < nums.length; j++) {
            for (let k = j + 1; k < nums.length; k++) {
                if (nums[i] + nums[j] + nums[k] === target) {
                    return [i, j, k];
                }
            }
        }
    }
    return [];
}

// Exercise 3: Two Sum returning values instead of indices
function twoSumValues(nums, target) {
    for (let i = 0; i < nums.length; i++) {
        for (let j = i + 1; j < nums.length; j++) {
            if (nums[i] + nums[j] === target) {
                return [nums[i], nums[j]];
            }
        }
    }
    return [];
}

// Test your solutions
console.log(twoSumSorted([1, 2, 3, 4, 5], 6)); // [0, 4]
console.log(threeSumBruteForce([1, 2, 3, 4, 5], 6)); // [0, 1, 2]
console.log(twoSumValues([2, 7, 11, 15], 9)); // [2, 7]
```

## 🎯 Key Learning Summary

### **Reverse Problems:**
- **Built-in methods**: Quick but may not be optimal for interviews
- **Two pointers**: Most efficient, O(n) time, O(1) space
- **Stack/Recursive**: Good for understanding data structures

### **Max/Min Problems:**
- **Linear scan**: Most practical, O(n) time
- **Tournament method**: Fewer comparisons but more complex
- **Edge cases**: Empty arrays, single elements

### **Character Counting:**
- **Object vs Map**: Object for simple cases, Map for complex keys
- **ASCII optimization**: For limited character sets
- **Analysis**: Building on basic counting

### **Two Sum:**
- **Brute force**: Foundation for understanding the problem
- **Time complexity**: O(n²) - understand why
- **Preparation**: For more efficient solutions

### **Next Steps:**
1. Practice each problem until you can code it from memory
2. Analyze time/space complexity for each approach
3. Test with edge cases
4. Try the variations and exercises

These fundamental problems build the essential skills needed for more complex algorithmic challenges!