# Practice Problems: Array & String Techniques

Let me explain and provide practice solutions for these intermediate problems that teach important algorithmic patterns.

## 1. Move Zeros

### Problem Understanding:
Move all zeros in an array to the end while maintaining the relative order of non-zero elements.

**Example:**
```
Input: [0,1,0,3,12]
Output: [1,3,12,0,0]
```

### Solutions:

```javascript
// Method 1: Two Pointers (Optimal - O(n) time, O(1) space)
function moveZeroes(nums) {
    let nonZeroPointer = 0;
    
    // Move all non-zero elements to the front
    for (let i = 0; i < nums.length; i++) {
        if (nums[i] !== 0) {
            nums[nonZeroPointer] = nums[i];
            nonZeroPointer++;
        }
    }
    
    // Fill remaining positions with zeros
    for (let i = nonZeroPointer; i < nums.length; i++) {
        nums[i] = 0;
    }
    
    return nums;
}

// Method 2: Snowball Method (Single Pass)
function moveZeroesSnowball(nums) {
    let snowballSize = 0;
    
    for (let i = 0; i < nums.length; i++) {
        if (nums[i] === 0) {
            snowballSize++;
        } else if (snowballSize > 0) {
            // Swap current element with first zero in snowball
            [nums[i], nums[i - snowballSize]] = [nums[i - snowballSize], nums[i]];
        }
    }
    
    return nums;
}

// Method 3: Using Extra Array (Not in-place)
function moveZeroesExtraArray(nums) {
    const result = [];
    let zeroCount = 0;
    
    // Collect non-zero elements
    for (let num of nums) {
        if (num !== 0) {
            result.push(num);
        } else {
            zeroCount++;
        }
    }
    
    // Append zeros at the end
    while (zeroCount > 0) {
        result.push(0);
        zeroCount--;
    }
    
    return result;
}

// Test cases
console.log("Move Zeroes Tests:");
console.log(moveZeroes([0,1,0,3,12])); // [1,3,12,0,0]
console.log(moveZeroes([0,0,1])); // [1,0,0]
console.log(moveZeroes([1,2,3,0,0])); // [1,2,3,0,0]
console.log(moveZeroesSnowball([0,1,0,3,12])); // [1,3,12,0,0]
```

**Time Complexity:** O(n)  
**Space Complexity:** O(1) for in-place methods

---

## 2. Container With Most Water

### Problem Understanding:
Given an array of heights, find two lines that form a container with the maximum area.

**Example:**
```
Input: [1,8,6,2,5,4,8,3,7]
Output: 49 (between indices 1 and 8)
```

### Solutions:

```javascript
// Method 1: Brute Force (O(n²))
function maxAreaBruteForce(height) {
    let maxArea = 0;
    
    for (let i = 0; i < height.length; i++) {
        for (let j = i + 1; j < height.length; j++) {
            const currentArea = Math.min(height[i], height[j]) * (j - i);
            maxArea = Math.max(maxArea, currentArea);
        }
    }
    
    return maxArea;
}

// Method 2: Two Pointers (Optimal - O(n))
function maxArea(height) {
    let left = 0;
    let right = height.length - 1;
    let maxArea = 0;
    
    while (left < right) {
        // Calculate current area
        const width = right - left;
        const currentHeight = Math.min(height[left], height[right]);
        const currentArea = width * currentHeight;
        
        // Update max area
        maxArea = Math.max(maxArea, currentArea);
        
        // Move the pointer with smaller height
        if (height[left] < height[right]) {
            left++;
        } else {
            right--;
        }
    }
    
    return maxArea;
}

// Method 3: Two Pointers with Early Termination
function maxAreaOptimized(height) {
    let left = 0, right = height.length - 1;
    let maxArea = 0;
    
    while (left < right) {
        const width = right - left;
        const minHeight = Math.min(height[left], height[right]);
        maxArea = Math.max(maxArea, width * minHeight);
        
        // Move pointers until we find taller lines
        while (left < right && height[left] <= minHeight) {
            left++;
        }
        while (left < right && height[right] <= minHeight) {
            right--;
        }
    }
    
    return maxArea;
}

// Visualization helper
function visualizeContainer(height) {
    const [left, right] = findMaxContainerIndices(height);
    const area = Math.min(height[left], height[right]) * (right - left);
    
    console.log(`Max area: ${area}`);
    console.log(`Between indices ${left} (height ${height[left]}) and ${right} (height ${height[right]})`);
}

function findMaxContainerIndices(height) {
    let left = 0, right = height.length - 1;
    let maxArea = 0;
    let bestLeft = 0, bestRight = height.length - 1;
    
    while (left < right) {
        const area = Math.min(height[left], height[right]) * (right - left);
        if (area > maxArea) {
            maxArea = area;
            bestLeft = left;
            bestRight = right;
        }
        
        if (height[left] < height[right]) left++;
        else right--;
    }
    
    return [bestLeft, bestRight];
}

// Test cases
console.log("\nContainer With Most Water Tests:");
console.log(maxArea([1,8,6,2,5,4,8,3,7])); // 49
console.log(maxArea([1,1])); // 1
console.log(maxArea([4,3,2,1,4])); // 16
console.log(maxArea([1,2,1])); // 2

visualizeContainer([1,8,6,2,5,4,8,3,7]);
```

**Time Complexity:** O(n) with two pointers  
**Space Complexity:** O(1)

---

## 3. Longest Substring Without Repeating Characters

### Problem Understanding:
Find the length of the longest substring without repeating characters.

**Example:**
```
Input: "abcabcbb"
Output: 3 ("abc")
```

### Solutions:

```javascript
// Method 1: Sliding Window with Set (Optimal - O(n))
function lengthOfLongestSubstring(s) {
    const charSet = new Set();
    let left = 0;
    let maxLength = 0;
    
    for (let right = 0; right < s.length; right++) {
        // If duplicate found, shrink window from left
        while (charSet.has(s[right])) {
            charSet.delete(s[left]);
            left++;
        }
        
        // Add current character to set
        charSet.add(s[right]);
        
        // Update max length
        maxLength = Math.max(maxLength, right - left + 1);
    }
    
    return maxLength;
}

// Method 2: Sliding Window with Map (More Efficient)
function lengthOfLongestSubstringMap(s) {
    const charMap = new Map(); // Stores character -> index
    let left = 0;
    let maxLength = 0;
    
    for (let right = 0; right < s.length; right++) {
        const currentChar = s[right];
        
        // If character exists and is within current window
        if (charMap.has(currentChar) && charMap.get(currentChar) >= left) {
            // Move left pointer to position after duplicate
            left = charMap.get(currentChar) + 1;
        }
        
        // Update character's latest position
        charMap.set(currentChar, right);
        
        // Update max length
        maxLength = Math.max(maxLength, right - left + 1);
    }
    
    return maxLength;
}

// Method 3: Using Array as Character Map (ASCII optimization)
function lengthOfLongestSubstringArray(s) {
    const charIndex = new Array(128).fill(-1); // For ASCII characters
    let left = 0;
    let maxLength = 0;
    
    for (let right = 0; right < s.length; right++) {
        const charCode = s.charCodeAt(right);
        
        // If character seen and within current window
        if (charIndex[charCode] >= left) {
            left = charIndex[charCode] + 1;
        }
        
        charIndex[charCode] = right;
        maxLength = Math.max(maxLength, right - left + 1);
    }
    
    return maxLength;
}

// Helper function to get the actual substring
function getLongestUniqueSubstring(s) {
    const charMap = new Map();
    let left = 0;
    let maxLength = 0;
    let startIndex = 0;
    
    for (let right = 0; right < s.length; right++) {
        const currentChar = s[right];
        
        if (charMap.has(currentChar) && charMap.get(currentChar) >= left) {
            left = charMap.get(currentChar) + 1;
        }
        
        charMap.set(currentChar, right);
        
        if (right - left + 1 > maxLength) {
            maxLength = right - left + 1;
            startIndex = left;
        }
    }
    
    return s.substring(startIndex, startIndex + maxLength);
}

// Test cases
console.log("\nLongest Substring Tests:");
console.log(lengthOfLongestSubstring("abcabcbb")); // 3
console.log(lengthOfLongestSubstring("bbbbb")); // 1
console.log(lengthOfLongestSubstring("pwwkew")); // 3
console.log(lengthOfLongestSubstring("")); // 0
console.log(lengthOfLongestSubstring("dvdf")); // 3

console.log("Actual substrings:");
console.log(getLongestUniqueSubstring("abcabcbb")); // "abc"
console.log(getLongestUniqueSubstring("pwwkew")); // "wke"
```

**Time Complexity:** O(n)  
**Space Complexity:** O(min(m, n)) where m is character set size

---

## 4. Valid Palindrome

### Problem Understanding:
Check if a string is a palindrome, considering only alphanumeric characters and ignoring cases.

**Example:**
```
Input: "A man, a plan, a canal: Panama"
Output: true
```

### Solutions:

```javascript
// Method 1: Two Pointers with Character Validation
function isPalindrome(s) {
    // Helper function to check if character is alphanumeric
    function isAlphanumeric(char) {
        return (char >= 'a' && char <= 'z') || 
               (char >= '0' && char <= '9');
    }
    
    let left = 0;
    let right = s.length - 1;
    s = s.toLowerCase();
    
    while (left < right) {
        // Skip non-alphanumeric characters from left
        while (left < right && !isAlphanumeric(s[left])) {
            left++;
        }
        
        // Skip non-alphanumeric characters from right
        while (left < right && !isAlphanumeric(s[right])) {
            right--;
        }
        
        // Compare characters
        if (s[left] !== s[right]) {
            return false;
        }
        
        left++;
        right--;
    }
    
    return true;
}

// Method 2: Using Regular Expressions (Cleaner)
function isPalindromeRegex(s) {
    // Remove non-alphanumeric characters and convert to lowercase
    const cleanString = s.replace(/[^a-zA-Z0-9]/g, '').toLowerCase();
    
    let left = 0;
    let right = cleanString.length - 1;
    
    while (left < right) {
        if (cleanString[left] !== cleanString[right]) {
            return false;
        }
        left++;
        right--;
    }
    
    return true;
}

// Method 3: Compare with Reverse
function isPalindromeReverse(s) {
    const cleanString = s.replace(/[^a-zA-Z0-9]/g, '').toLowerCase();
    return cleanString === cleanString.split('').reverse().join('');
}

// Method 4: Recursive Approach
function isPalindromeRecursive(s) {
    const cleanString = s.replace(/[^a-zA-Z0-9]/g, '').toLowerCase();
    
    function checkPalindrome(str, left, right) {
        // Base case
        if (left >= right) return true;
        
        // Check if characters match
        if (str[left] !== str[right]) return false;
        
        // Recursive call
        return checkPalindrome(str, left + 1, right - 1);
    }
    
    return checkPalindrome(cleanString, 0, cleanString.length - 1);
}

// Test cases
console.log("\nValid Palindrome Tests:");
console.log(isPalindrome("A man, a plan, a canal: Panama")); // true
console.log(isPalindrome("race a car")); // false
console.log(isPalindrome(" ")); // true
console.log(isPalindrome("0P")); // false
console.log(isPalindrome("ab2a")); // false

console.log(isPalindromeRegex("A man, a plan, a canal: Panama")); // true
console.log(isPalindromeReverse("racecar")); // true
```

**Time Complexity:** O(n)  
**Space Complexity:** O(1) for two pointers, O(n) for reverse method

---

## 🎯 Practice Exercises

### Exercise 1: Move Zeroes Variations
```javascript
// 1. Move zeros to front instead of back
function moveZeroesToFront(nums) {
    // Your code here
}

// 2. Move all duplicates to end
function moveDuplicatesToEnd(nums) {
    // Your code here
}
```

### Exercise 2: Container Variations
```javascript
// 1. Find the container with minimum area
function minArea(height) {
    // Your code here
}

// 2. Find three lines that form maximum area
function maxAreaThreeLines(height) {
    // Your code here
}
```

### Exercise 3: Substring Variations
```javascript
// 1. Find longest substring with at most K distinct characters
function longestSubstringKDistinct(s, k) {
    // Your code here
}

// 2. Find all substrings without repeating characters
function allUniqueSubstrings(s) {
    // Your code here
}
```

### Exercise 4: Palindrome Variations
```javascript
// 1. Check if string can be palindrome after removing one character
function validPalindromeII(s) {
    // Your code here
}

// 2. Find longest palindromic substring
function longestPalindrome(s) {
    // Your code here
}
```

## 🚀 Key Patterns Learned:

1. **Two Pointers**: Move Zeroes, Container, Palindrome
2. **Sliding Window**: Longest Substring
3. **Character Validation**: Palindrome checking
4. **In-place Modification**: Move Zeroes
5. **Early Termination**: Container optimization

These problems teach essential patterns that appear in many interview questions. Practice each approach until you can implement them from memory!