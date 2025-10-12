# Practice Problems: Backtracking & Recursion

Let me explain and provide practice solutions for these essential backtracking problems that form the foundation of combinatorial algorithms.

## 1. Generate Parentheses

### Problem Understanding:
Generate all combinations of well-formed parentheses given n pairs.

**Example:**
```
Input: n = 3
Output: ["((()))","(()())","(())()","()(())","()()()"]
```

### Solutions:

```javascript
// Method 1: Backtracking (Optimal)
function generateParenthesis(n) {
    const result = [];
    
    function backtrack(current, open, close) {
        // Base case: we've used all parentheses
        if (current.length === 2 * n) {
            result.push(current);
            return;
        }
        
        // Add opening parenthesis if we haven't used all
        if (open < n) {
            backtrack(current + '(', open + 1, close);
        }
        
        // Add closing parenthesis if it's valid (more opens than closes)
        if (close < open) {
            backtrack(current + ')', open, close + 1);
        }
    }
    
    backtrack('', 0, 0);
    return result;
}

// Method 2: Iterative with Stack
function generateParenthesisIterative(n) {
    const result = [];
    const stack = [{ str: '', open: 0, close: 0 }];
    
    while (stack.length > 0) {
        const { str, open, close } = stack.pop();
        
        if (str.length === 2 * n) {
            result.push(str);
            continue;
        }
        
        if (open < n) {
            stack.push({ str: str + '(', open: open + 1, close });
        }
        
        if (close < open) {
            stack.push({ str: str + ')', open, close: close + 1 });
        }
    }
    
    return result;
}

// Method 3: Closure Number (Mathematical approach)
function generateParenthesisClosure(n) {
    if (n === 0) return [''];
    
    const result = [];
    for (let i = 0; i < n; i++) {
        for (let left of generateParenthesisClosure(i)) {
            for (let right of generateParenthesisClosure(n - 1 - i)) {
                result.push(`(${left})${right}`);
            }
        }
    }
    return result;
}

// Enhanced version with step tracking
function generateParenthesisWithSteps(n) {
    const result = [];
    const steps = [];
    
    function backtrack(current, open, close, depth = 0) {
        steps.push({
            depth,
            current,
            open,
            close,
            action: `Current: "${current}"`
        });
        
        if (current.length === 2 * n) {
            steps.push({
                depth,
                action: `✓ Found valid: "${current}"`,
                isSolution: true
            });
            result.push(current);
            return;
        }
        
        if (open < n) {
            steps.push({
                depth,
                action: '→ Adding "("'
            });
            backtrack(current + '(', open + 1, close, depth + 1);
        }
        
        if (close < open) {
            steps.push({
                depth,
                action: '→ Adding ")"'
            });
            backtrack(current + ')', open, close + 1, depth + 1);
        }
    }
    
    backtrack('', 0, 0);
    return { result, steps };
}

// Test cases
console.log("Generate Parentheses Tests:");
console.log(generateParenthesis(2)); // ['(())', '()()']
console.log(generateParenthesis(3)); // 5 combinations

const withSteps = generateParenthesisWithSteps(2);
console.log("Solutions:", withSteps.result);
console.log("Steps:", withSteps.steps.slice(0, 10)); // First 10 steps

// Count valid parentheses for different n
for (let i = 1; i <= 5; i++) {
    const result = generateParenthesis(i);
    console.log(`n=${i}: ${result.length} combinations`);
}
```

**Time Complexity:** O(4ⁿ/√n) - Catalan number  
**Space Complexity:** O(n) for recursion stack

---

## 2. Subsets

### Problem Understanding:
Given distinct integers, return all possible subsets (the power set).

**Example:**
```
Input: [1,2,3]
Output: [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]
```

### Solutions:

```javascript
// Method 1: Backtracking (Cascading)
function subsets(nums) {
    const result = [];
    
    function backtrack(start, current) {
        // Add current subset to result
        result.push([...current]);
        
        // Generate all subsets that include nums[i]
        for (let i = start; i < nums.length; i++) {
            current.push(nums[i]);
            backtrack(i + 1, current);
            current.pop(); // backtrack
        }
    }
    
    backtrack(0, []);
    return result;
}

// Method 2: Bit Manipulation
function subsetsBitmask(nums) {
    const n = nums.length;
    const total = 1 << n; // 2^n subsets
    const result = [];
    
    for (let mask = 0; mask < total; mask++) {
        const subset = [];
        for (let i = 0; i < n; i++) {
            // Check if i-th bit is set in mask
            if (mask & (1 << i)) {
                subset.push(nums[i]);
            }
        }
        result.push(subset);
    }
    
    return result;
}

// Method 3: Iterative BFS-like
function subsetsIterative(nums) {
    let result = [[]];
    
    for (let num of nums) {
        const newSubsets = [];
        for (let subset of result) {
            newSubsets.push([...subset, num]);
        }
        result = result.concat(newSubsets);
    }
    
    return result;
}

// Method 4: Recursive without backtracking
function subsetsRecursive(nums) {
    if (nums.length === 0) return [[]];
    
    const first = nums[0];
    const rest = subsetsRecursive(nums.slice(1));
    
    const withFirst = rest.map(subset => [first, ...subset]);
    return [...rest, ...withFirst];
}

// Enhanced version with visualization
function subsetsWithSteps(nums) {
    const result = [];
    const steps = [];
    
    function backtrack(start, current, depth = 0) {
        const indent = '  '.repeat(depth);
        steps.push(`${indent}Subset: [${current}]`);
        
        result.push([...current]);
        
        for (let i = start; i < nums.length; i++) {
            steps.push(`${indent}→ Adding ${nums[i]}`);
            current.push(nums[i]);
            backtrack(i + 1, current, depth + 1);
            steps.push(`${indent}← Backtrack (remove ${nums[i]})`);
            current.pop();
        }
    }
    
    backtrack(0, []);
    return { result, steps };
}

// Test cases
console.log("\nSubsets Tests:");
console.log(subsets([1, 2, 3]));
// [[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]

console.log(subsetsBitmask([1, 2, 3])); // Same result
console.log(subsetsIterative([1, 2])); // [[],[1],[2],[1,2]]

const withSteps = subsetsWithSteps([1, 2]);
console.log("Subsets:", withSteps.result);
```

**Time Complexity:** O(n × 2ⁿ) - 2ⁿ subsets, each taking O(n) time  
**Space Complexity:** O(n) for recursion stack

---

## 3. Permutations

### Problem Understanding:
Given distinct integers, return all possible permutations.

**Example:**
```
Input: [1,2,3]
Output: [[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
```

### Solutions:

```javascript
// Method 1: Backtracking with Swapping
function permute(nums) {
    const result = [];
    
    function backtrack(start) {
        if (start === nums.length) {
            result.push([...nums]);
            return;
        }
        
        for (let i = start; i < nums.length; i++) {
            // Swap elements
            [nums[start], nums[i]] = [nums[i], nums[start]];
            backtrack(start + 1);
            // Backtrack (swap back)
            [nums[start], nums[i]] = [nums[i], nums[start]];
        }
    }
    
    backtrack(0);
    return result;
}

// Method 2: Backtracking with Used Array
function permuteUsed(nums) {
    const result = [];
    const used = new Array(nums.length).fill(false);
    
    function backtrack(current) {
        if (current.length === nums.length) {
            result.push([...current]);
            return;
        }
        
        for (let i = 0; i < nums.length; i++) {
            if (!used[i]) {
                used[i] = true;
                current.push(nums[i]);
                backtrack(current);
                current.pop();
                used[i] = false;
            }
        }
    }
    
    backtrack([]);
    return result;
}

// Method 3: Heap's Algorithm (Efficient)
function permuteHeap(nums) {
    const result = [];
    const n = nums.length;
    
    function generate(k, arr) {
        if (k === 1) {
            result.push([...arr]);
            return;
        }
        
        generate(k - 1, arr);
        
        for (let i = 0; i < k - 1; i++) {
            if (k % 2 === 0) {
                // Even k: swap i and k-1
                [arr[i], arr[k - 1]] = [arr[k - 1], arr[i]];
            } else {
                // Odd k: swap 0 and k-1
                [arr[0], arr[k - 1]] = [arr[k - 1], arr[0]];
            }
            generate(k - 1, arr);
        }
    }
    
    generate(n, [...nums]);
    return result;
}

// Method 4: Iterative
function permuteIterative(nums) {
    let result = [[]];
    
    for (let num of nums) {
        const newPermutations = [];
        
        for (let permutation of result) {
            for (let i = 0; i <= permutation.length; i++) {
                const newPerm = [...permutation];
                newPerm.splice(i, 0, num);
                newPermutations.push(newPerm);
            }
        }
        
        result = newPermutations;
    }
    
    return result;
}

// Enhanced version with visualization
function permuteWithSteps(nums) {
    const result = [];
    const steps = [];
    
    function backtrack(start, depth = 0) {
        const indent = '  '.repeat(depth);
        steps.push(`${indent}Current array: [${nums}]`);
        
        if (start === nums.length) {
            steps.push(`${indent}✓ Found permutation: [${nums}]`);
            result.push([...nums]);
            return;
        }
        
        for (let i = start; i < nums.length; i++) {
            steps.push(`${indent}Swap ${start} with ${i}`);
            [nums[start], nums[i]] = [nums[i], nums[start]];
            backtrack(start + 1, depth + 1);
            steps.push(`${indent}Backtrack (swap ${start} with ${i})`);
            [nums[start], nums[i]] = [nums[i], nums[start]];
        }
    }
    
    backtrack(0);
    return { result, steps };
}

// Test cases
console.log("\nPermutations Tests:");
console.log(permute([1, 2, 3])); // 6 permutations
console.log(permuteUsed([1, 2])); // [[1,2],[2,1]]
console.log(permuteHeap([1, 2, 3])); // Same as permute

// Count permutations for different sizes
for (let i = 1; i <= 4; i++) {
    const arr = Array.from({ length: i }, (_, idx) => idx + 1);
    const perms = permute(arr);
    console.log(`Size ${i}: ${perms.length} permutations`);
}
```

**Time Complexity:** O(n × n!) - n! permutations, each taking O(n) time  
**Space Complexity:** O(n) for recursion stack

---

## 4. Combination Sum

### Problem Understanding:
Given candidates and target, find all unique combinations where candidate numbers sum to target. The same number may be used unlimited times.

**Example:**
```
Input: candidates = [2,3,6,7], target = 7
Output: [[2,2,3],[7]]
```

### Solutions:

```javascript
// Method 1: Backtracking
function combinationSum(candidates, target) {
    const result = [];
    
    function backtrack(start, current, remaining) {
        if (remaining === 0) {
            result.push([...current]);
            return;
        }
        
        if (remaining < 0) return;
        
        for (let i = start; i < candidates.length; i++) {
            current.push(candidates[i]);
            backtrack(i, current, remaining - candidates[i]); // Note: i not i+1 (reuse allowed)
            current.pop();
        }
    }
    
    backtrack(0, [], target);
    return result;
}

// Method 2: Backtracking with Sorting Optimization
function combinationSumOptimized(candidates, target) {
    const result = [];
    candidates.sort((a, b) => a - b); // Sort to enable early termination
    
    function backtrack(start, current, remaining) {
        if (remaining === 0) {
            result.push([...current]);
            return;
        }
        
        for (let i = start; i < candidates.length; i++) {
            // Early termination if candidate > remaining
            if (candidates[i] > remaining) break;
            
            current.push(candidates[i]);
            backtrack(i, current, remaining - candidates[i]);
            current.pop();
        }
    }
    
    backtrack(0, [], target);
    return result;
}

// Method 3: Dynamic Programming (Bottom-up)
function combinationSumDP(candidates, target) {
    const dp = new Array(target + 1).fill().map(() => []);
    dp[0] = [[]]; // Base case: one way to make sum 0 (empty combination)
    
    for (let candidate of candidates) {
        for (let sum = candidate; sum <= target; sum++) {
            for (let combination of dp[sum - candidate]) {
                dp[sum].push([...combination, candidate]);
            }
        }
    }
    
    return dp[target];
}

// Method 4: Combination Sum II (No duplicates, use once)
function combinationSum2(candidates, target) {
    const result = [];
    candidates.sort((a, b) => a - b);
    
    function backtrack(start, current, remaining) {
        if (remaining === 0) {
            result.push([...current]);
            return;
        }
        
        for (let i = start; i < candidates.length; i++) {
            // Skip duplicates
            if (i > start && candidates[i] === candidates[i - 1]) continue;
            
            if (candidates[i] > remaining) break;
            
            current.push(candidates[i]);
            backtrack(i + 1, current, remaining - candidates[i]); // i+1: no reuse
            current.pop();
        }
    }
    
    backtrack(0, [], target);
    return result;
}

// Enhanced version with visualization
function combinationSumWithSteps(candidates, target) {
    const result = [];
    const steps = [];
    
    function backtrack(start, current, remaining, depth = 0) {
        const indent = '  '.repeat(depth);
        steps.push(`${indent}Current: [${current}], Remaining: ${remaining}`);
        
        if (remaining === 0) {
            steps.push(`${indent}✓ Found combination: [${current}]`);
            result.push([...current]);
            return;
        }
        
        if (remaining < 0) {
            steps.push(`${indent}✗ Exceeded target`);
            return;
        }
        
        for (let i = start; i < candidates.length; i++) {
            steps.push(`${indent}→ Trying ${candidates[i]}`);
            current.push(candidates[i]);
            backtrack(i, current, remaining - candidates[i], depth + 1);
            steps.push(`${indent}← Backtrack (remove ${candidates[i]})`);
            current.pop();
        }
    }
    
    backtrack(0, [], target);
    return { result, steps };
}

// Test cases
console.log("\nCombination Sum Tests:");
console.log(combinationSum([2, 3, 6, 7], 7)); // [[2,2,3],[7]]
console.log(combinationSum([2, 3, 5], 8)); // [[2,2,2,2],[2,3,3],[3,5]]

console.log("Combination Sum II (no duplicates):");
console.log(combinationSum2([10, 1, 2, 7, 6, 1, 5], 8)); 
// [[1,1,6],[1,2,5],[1,7],[2,6]]

console.log("DP approach:");
console.log(combinationSumDP([2, 3, 5], 8));

const withSteps = combinationSumWithSteps([2, 3], 5);
console.log("Solutions:", withSteps.result);
```

**Time Complexity:** O(N^(T/M + 1)) where N is candidates length, T is target, M is min candidate  
**Space Complexity:** O(T/M) for recursion stack

---

## 🎯 Practice Exercises

### Exercise 1: Parentheses Variations
```javascript
// 1. Generate parentheses with different characters
function generateParenthesesVariations(n, openChar = '(', closeChar = ')') {
    // Your code here
}

// 2. Count valid parentheses without generating
function countValidParentheses(n) {
    // Your code here (Catalan numbers)
}
```

### Exercise 2: Subsets Variations
```javascript
// 1. Subsets with duplicates
function subsetsWithDup(nums) {
    // Your code here
}

// 2. Subsets of fixed size k
function subsetsOfSizeK(nums, k) {
    // Your code here
}
```

### Exercise 3: Permutations Variations
```javascript
// 1. Permutations with duplicates
function permuteUnique(nums) {
    // Your code here
}

// 2. Next permutation
function nextPermutation(nums) {
    // Your code here
}
```

### Exercise 4: Combination Variations
```javascript
// 1. Combination sum with constraints
function combinationSum3(k, n) {
    // Your code here (use k numbers that sum to n)
}

// 2. Letter combinations of phone number
function letterCombinations(digits) {
    // Your code here
}
```

## 🚀 Key Patterns Learned:

1. **Backtracking Template**: Choose → Explore → Unchoose
2. **State Tracking**: Maintain current state and constraints
3. **Pruning**: Early termination to avoid unnecessary work
4. **Duplicate Handling**: Sorting and skipping duplicates
5. **Multiple Approaches**: Backtracking vs DP vs Iterative

These combinatorial problems are fundamental for technical interviews and teach essential recursive thinking patterns!