# Part 9: Dynamic Programming

## Table of Contents
1. [Introduction to Dynamic Programming](#introduction)
2. [Key Concepts](#key-concepts)
3. [Types of Dynamic Programming](#types)
4. [Basic DP Problems](#basic-problems)
5. [1D Dynamic Programming](#1d-dp)
6. [2D Dynamic Programming](#2d-dp)
7. [String DP Problems](#string-dp)
8. [Tree DP](#tree-dp)
9. [Advanced DP Techniques](#advanced-techniques)
10. [DP Optimization Techniques](#optimization)
11. [Practice Problems](#practice-problems)
12. [Complete Example Program](#example-program)
13. [Summary](#summary)

---

## Introduction to Dynamic Programming {#introduction}

**Dynamic Programming (DP)** is an algorithmic paradigm that solves complex problems by breaking them down into simpler subproblems. It is applicable when the subproblems overlap and the optimal solution can be constructed from optimal solutions of subproblems.

### Key Characteristics
1. **Overlapping Subproblems**: Same subproblems are solved multiple times
2. **Optimal Substructure**: Optimal solution contains optimal solutions of subproblems
3. **Memoization**: Store results to avoid recomputation

### DP vs Divide and Conquer
- **Divide and Conquer**: Subproblems are independent
- **Dynamic Programming**: Subproblems overlap

### When to Use DP
- Optimization problems (min/max)
- Counting problems
- Decision problems (yes/no)
- Problems with overlapping subproblems

---

## Key Concepts {#key-concepts}

### 1. Memoization (Top-Down)
```cpp
// Recursive approach with memoization
int fibonacci(int n, vector<int>& memo) {
    if (n <= 1) return n;
    if (memo[n] != -1) return memo[n];
    
    memo[n] = fibonacci(n-1, memo) + fibonacci(n-2, memo);
    return memo[n];
}
```

### 2. Tabulation (Bottom-Up)
```cpp
// Iterative approach with table
int fibonacci(int n) {
    if (n <= 1) return n;
    
    vector<int> dp(n + 1);
    dp[0] = 0;
    dp[1] = 1;
    
    for (int i = 2; i <= n; i++) {
        dp[i] = dp[i-1] + dp[i-2];
    }
    
    return dp[n];
}
```

### 3. Space Optimization
```cpp
// Space optimized version
int fibonacci(int n) {
    if (n <= 1) return n;
    
    int prev2 = 0, prev1 = 1;
    
    for (int i = 2; i <= n; i++) {
        int curr = prev1 + prev2;
        prev2 = prev1;
        prev1 = curr;
    }
    
    return prev1;
}
```

---

## Types of Dynamic Programming {#types}

### 1. Linear DP
- Problems with 1D state space
- Examples: Fibonacci, Climbing Stairs, House Robber

### 2. Grid DP
- Problems with 2D state space
- Examples: Unique Paths, Minimum Path Sum

### 3. Interval DP
- Problems on intervals/ranges
- Examples: Matrix Chain Multiplication, Palindrome Partitioning

### 4. Tree DP
- Problems on trees
- Examples: Tree Diameter, Maximum Path Sum

### 5. Digit DP
- Problems involving digits
- Examples: Count numbers with specific properties

### 6. Bitmask DP
- Problems with subset states
- Examples: Traveling Salesman Problem

---

## Basic DP Problems {#basic-problems}

### 1. Fibonacci Sequence
```cpp
class FibonacciDP {
public:
    // Memoization approach
    static int fibMemo(int n) {
        vector<int> memo(n + 1, -1);
        return fibMemoHelper(n, memo);
    }
    
    // Tabulation approach
    static int fibTab(int n) {
        if (n <= 1) return n;
        
        vector<int> dp(n + 1);
        dp[0] = 0;
        dp[1] = 1;
        
        for (int i = 2; i <= n; i++) {
            dp[i] = dp[i-1] + dp[i-2];
        }
        
        return dp[n];
    }
    
    // Space optimized
    static int fibOptimized(int n) {
        if (n <= 1) return n;
        
        int prev2 = 0, prev1 = 1;
        
        for (int i = 2; i <= n; i++) {
            int curr = prev1 + prev2;
            prev2 = prev1;
            prev1 = curr;
        }
        
        return prev1;
    }
    
private:
    static int fibMemoHelper(int n, vector<int>& memo) {
        if (n <= 1) return n;
        if (memo[n] != -1) return memo[n];
        
        memo[n] = fibMemoHelper(n-1, memo) + fibMemoHelper(n-2, memo);
        return memo[n];
    }
};
```

### 2. Climbing Stairs
```cpp
class ClimbingStairs {
public:
    // You can climb 1 or 2 steps at a time
    static int climbStairs(int n) {
        if (n <= 2) return n;
        
        vector<int> dp(n + 1);
        dp[1] = 1;
        dp[2] = 2;
        
        for (int i = 3; i <= n; i++) {
            dp[i] = dp[i-1] + dp[i-2];
        }
        
        return dp[n];
    }
    
    // Space optimized
    static int climbStairsOptimized(int n) {
        if (n <= 2) return n;
        
        int prev2 = 1, prev1 = 2;
        
        for (int i = 3; i <= n; i++) {
            int curr = prev1 + prev2;
            prev2 = prev1;
            prev1 = curr;
        }
        
        return prev1;
    }
    
    // With variable steps (can climb 1, 2, or 3 steps)
    static int climbStairsVariable(int n) {
        if (n <= 0) return 0;
        if (n == 1) return 1;
        if (n == 2) return 2;
        if (n == 3) return 4;
        
        vector<int> dp(n + 1);
        dp[0] = 1;
        dp[1] = 1;
        dp[2] = 2;
        dp[3] = 4;
        
        for (int i = 4; i <= n; i++) {
            dp[i] = dp[i-1] + dp[i-2] + dp[i-3];
        }
        
        return dp[n];
    }
};
```

### 3. House Robber
```cpp
class HouseRobber {
public:
    // Cannot rob adjacent houses
    static int rob(vector<int>& nums) {
        int n = nums.size();
        if (n == 0) return 0;
        if (n == 1) return nums[0];
        
        vector<int> dp(n);
        dp[0] = nums[0];
        dp[1] = max(nums[0], nums[1]);
        
        for (int i = 2; i < n; i++) {
            dp[i] = max(dp[i-1], dp[i-2] + nums[i]);
        }
        
        return dp[n-1];
    }
    
    // Space optimized
    static int robOptimized(vector<int>& nums) {
        int n = nums.size();
        if (n == 0) return 0;
        if (n == 1) return nums[0];
        
        int prev2 = nums[0];
        int prev1 = max(nums[0], nums[1]);
        
        for (int i = 2; i < n; i++) {
            int curr = max(prev1, prev2 + nums[i]);
            prev2 = prev1;
            prev1 = curr;
        }
        
        return prev1;
    }
    
    // House Robber II (circular arrangement)
    static int robCircular(vector<int>& nums) {
        int n = nums.size();
        if (n == 0) return 0;
        if (n == 1) return nums[0];
        if (n == 2) return max(nums[0], nums[1]);
        
        // Case 1: Rob houses 0 to n-2
        int case1 = robRange(nums, 0, n-2);
        
        // Case 2: Rob houses 1 to n-1
        int case2 = robRange(nums, 1, n-1);
        
        return max(case1, case2);
    }
    
private:
    static int robRange(vector<int>& nums, int start, int end) {
        int prev2 = 0, prev1 = 0;
        
        for (int i = start; i <= end; i++) {
            int curr = max(prev1, prev2 + nums[i]);
            prev2 = prev1;
            prev1 = curr;
        }
        
        return prev1;
    }
};
```

---

## 1D Dynamic Programming {#1d-dp}

### 1. Coin Change
```cpp
class CoinChange {
public:
    // Minimum coins to make amount
    static int coinChange(vector<int>& coins, int amount) {
        vector<int> dp(amount + 1, amount + 1);
        dp[0] = 0;
        
        for (int i = 1; i <= amount; i++) {
            for (int coin : coins) {
                if (coin <= i) {
                    dp[i] = min(dp[i], dp[i - coin] + 1);
                }
            }
        }
        
        return dp[amount] > amount ? -1 : dp[amount];
    }
    
    // Number of ways to make amount
    static int coinChangeWays(vector<int>& coins, int amount) {
        vector<int> dp(amount + 1, 0);
        dp[0] = 1;
        
        for (int coin : coins) {
            for (int i = coin; i <= amount; i++) {
                dp[i] += dp[i - coin];
            }
        }
        
        return dp[amount];
    }
    
    // Coin change with limited coins
    static int coinChangeLimited(vector<int>& coins, vector<int>& counts, int amount) {
        vector<int> dp(amount + 1, amount + 1);
        dp[0] = 0;
        
        for (int i = 0; i < coins.size(); i++) {
            int coin = coins[i];
            int count = counts[i];
            
            for (int j = amount; j >= coin; j--) {
                for (int k = 1; k <= count && k * coin <= j; k++) {
                    dp[j] = min(dp[j], dp[j - k * coin] + k);
                }
            }
        }
        
        return dp[amount] > amount ? -1 : dp[amount];
    }
};
```

### 2. Longest Increasing Subsequence
```cpp
class LongestIncreasingSubsequence {
public:
    // O(n²) DP solution
    static int lengthOfLIS(vector<int>& nums) {
        int n = nums.size();
        if (n == 0) return 0;
        
        vector<int> dp(n, 1);
        int maxLength = 1;
        
        for (int i = 1; i < n; i++) {
            for (int j = 0; j < i; j++) {
                if (nums[j] < nums[i]) {
                    dp[i] = max(dp[i], dp[j] + 1);
                }
            }
            maxLength = max(maxLength, dp[i]);
        }
        
        return maxLength;
    }
    
    // O(n log n) Binary Search solution
    static int lengthOfLISOptimized(vector<int>& nums) {
        vector<int> tails;
        
        for (int num : nums) {
            auto it = lower_bound(tails.begin(), tails.end(), num);
            if (it == tails.end()) {
                tails.push_back(num);
            } else {
                *it = num;
            }
        }
        
        return tails.size();
    }
    
    // Get actual LIS
    static vector<int> getLIS(vector<int>& nums) {
        int n = nums.size();
        if (n == 0) return {};
        
        vector<int> dp(n, 1);
        vector<int> parent(n, -1);
        int maxLength = 1;
        int maxIndex = 0;
        
        for (int i = 1; i < n; i++) {
            for (int j = 0; j < i; j++) {
                if (nums[j] < nums[i] && dp[j] + 1 > dp[i]) {
                    dp[i] = dp[j] + 1;
                    parent[i] = j;
                }
            }
            if (dp[i] > maxLength) {
                maxLength = dp[i];
                maxIndex = i;
            }
        }
        
        // Reconstruct LIS
        vector<int> lis;
        int current = maxIndex;
        while (current != -1) {
            lis.push_back(nums[current]);
            current = parent[current];
        }
        
        reverse(lis.begin(), lis.end());
        return lis;
    }
};
```

### 3. Maximum Subarray (Kadane's Algorithm)
```cpp
class MaximumSubarray {
public:
    // Maximum sum subarray
    static int maxSubArray(vector<int>& nums) {
        int maxSoFar = nums[0];
        int maxEndingHere = nums[0];
        
        for (int i = 1; i < nums.size(); i++) {
            maxEndingHere = max(nums[i], maxEndingHere + nums[i]);
            maxSoFar = max(maxSoFar, maxEndingHere);
        }
        
        return maxSoFar;
    }
    
    // Get the actual subarray
    static vector<int> maxSubArrayIndices(vector<int>& nums) {
        int maxSoFar = nums[0];
        int maxEndingHere = nums[0];
        int start = 0, end = 0, tempStart = 0;
        
        for (int i = 1; i < nums.size(); i++) {
            if (maxEndingHere < 0) {
                maxEndingHere = nums[i];
                tempStart = i;
            } else {
                maxEndingHere += nums[i];
            }
            
            if (maxEndingHere > maxSoFar) {
                maxSoFar = maxEndingHere;
                start = tempStart;
                end = i;
            }
        }
        
        return {start, end};
    }
    
    // Maximum product subarray
    static int maxProduct(vector<int>& nums) {
        int maxSoFar = nums[0];
        int maxEndingHere = nums[0];
        int minEndingHere = nums[0];
        
        for (int i = 1; i < nums.size(); i++) {
            if (nums[i] < 0) {
                swap(maxEndingHere, minEndingHere);
            }
            
            maxEndingHere = max(nums[i], maxEndingHere * nums[i]);
            minEndingHere = min(nums[i], minEndingHere * nums[i]);
            
            maxSoFar = max(maxSoFar, maxEndingHere);
        }
        
        return maxSoFar;
    }
};
```

---

## 2D Dynamic Programming {#2d-dp}

### 1. Unique Paths
```cpp
class UniquePaths {
public:
    // Unique paths in grid
    static int uniquePaths(int m, int n) {
        vector<vector<int>> dp(m, vector<int>(n, 1));
        
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                dp[i][j] = dp[i-1][j] + dp[i][j-1];
            }
        }
        
        return dp[m-1][n-1];
    }
    
    // Space optimized
    static int uniquePathsOptimized(int m, int n) {
        vector<int> dp(n, 1);
        
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                dp[j] += dp[j-1];
            }
        }
        
        return dp[n-1];
    }
    
    // Unique paths with obstacles
    static int uniquePathsWithObstacles(vector<vector<int>>& obstacleGrid) {
        int m = obstacleGrid.size();
        int n = obstacleGrid[0].size();
        
        if (obstacleGrid[0][0] == 1 || obstacleGrid[m-1][n-1] == 1) {
            return 0;
        }
        
        vector<vector<int>> dp(m, vector<int>(n, 0));
        dp[0][0] = 1;
        
        // Fill first row
        for (int j = 1; j < n; j++) {
            dp[0][j] = (obstacleGrid[0][j] == 0) ? dp[0][j-1] : 0;
        }
        
        // Fill first column
        for (int i = 1; i < m; i++) {
            dp[i][0] = (obstacleGrid[i][0] == 0) ? dp[i-1][0] : 0;
        }
        
        // Fill rest of the grid
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                if (obstacleGrid[i][j] == 0) {
                    dp[i][j] = dp[i-1][j] + dp[i][j-1];
                }
            }
        }
        
        return dp[m-1][n-1];
    }
};
```

### 2. Minimum Path Sum
```cpp
class MinimumPathSum {
public:
    static int minPathSum(vector<vector<int>>& grid) {
        int m = grid.size();
        int n = grid[0].size();
        
        vector<vector<int>> dp(m, vector<int>(n));
        dp[0][0] = grid[0][0];
        
        // Fill first row
        for (int j = 1; j < n; j++) {
            dp[0][j] = dp[0][j-1] + grid[0][j];
        }
        
        // Fill first column
        for (int i = 1; i < m; i++) {
            dp[i][0] = dp[i-1][0] + grid[i][0];
        }
        
        // Fill rest of the grid
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j];
            }
        }
        
        return dp[m-1][n-1];
    }
    
    // Space optimized
    static int minPathSumOptimized(vector<vector<int>>& grid) {
        int m = grid.size();
        int n = grid[0].size();
        
        vector<int> dp(n);
        dp[0] = grid[0][0];
        
        // Fill first row
        for (int j = 1; j < n; j++) {
            dp[j] = dp[j-1] + grid[0][j];
        }
        
        // Process remaining rows
        for (int i = 1; i < m; i++) {
            dp[0] += grid[i][0];
            for (int j = 1; j < n; j++) {
                dp[j] = min(dp[j], dp[j-1]) + grid[i][j];
            }
        }
        
        return dp[n-1];
    }
};
```

### 3. Longest Common Subsequence
```cpp
class LongestCommonSubsequence {
public:
    static int longestCommonSubsequence(string text1, string text2) {
        int m = text1.length();
        int n = text2.length();
        
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
        
        for (int i = 1; i <= m; i++) {
            for (int j = 1; j <= n; j++) {
                if (text1[i-1] == text2[j-1]) {
                    dp[i][j] = dp[i-1][j-1] + 1;
                } else {
                    dp[i][j] = max(dp[i-1][j], dp[i][j-1]);
                }
            }
        }
        
        return dp[m][n];
    }
    
    // Get actual LCS
    static string getLCS(string text1, string text2) {
        int m = text1.length();
        int n = text2.length();
        
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
        
        // Fill DP table
        for (int i = 1; i <= m; i++) {
            for (int j = 1; j <= n; j++) {
                if (text1[i-1] == text2[j-1]) {
                    dp[i][j] = dp[i-1][j-1] + 1;
                } else {
                    dp[i][j] = max(dp[i-1][j], dp[i][j-1]);
                }
            }
        }
        
        // Reconstruct LCS
        string lcs = "";
        int i = m, j = n;
        
        while (i > 0 && j > 0) {
            if (text1[i-1] == text2[j-1]) {
                lcs = text1[i-1] + lcs;
                i--;
                j--;
            } else if (dp[i-1][j] > dp[i][j-1]) {
                i--;
            } else {
                j--;
            }
        }
        
        return lcs;
    }
    
    // Space optimized
    static int longestCommonSubsequenceOptimized(string text1, string text2) {
        int m = text1.length();
        int n = text2.length();
        
        vector<int> prev(n + 1, 0);
        vector<int> curr(n + 1, 0);
        
        for (int i = 1; i <= m; i++) {
            for (int j = 1; j <= n; j++) {
                if (text1[i-1] == text2[j-1]) {
                    curr[j] = prev[j-1] + 1;
                } else {
                    curr[j] = max(prev[j], curr[j-1]);
                }
            }
            prev = curr;
        }
        
        return curr[n];
    }
};
```

---

## String DP Problems {#string-dp}

### 1. Edit Distance
```cpp
class EditDistance {
public:
    static int minDistance(string word1, string word2) {
        int m = word1.length();
        int n = word2.length();
        
        vector<vector<int>> dp(m + 1, vector<int>(n + 1));
        
        // Initialize base cases
        for (int i = 0; i <= m; i++) {
            dp[i][0] = i; // Delete all characters
        }
        
        for (int j = 0; j <= n; j++) {
            dp[0][j] = j; // Insert all characters
        }
        
        // Fill DP table
        for (int i = 1; i <= m; i++) {
            for (int j = 1; j <= n; j++) {
                if (word1[i-1] == word2[j-1]) {
                    dp[i][j] = dp[i-1][j-1]; // No operation needed
                } else {
                    dp[i][j] = 1 + min({
                        dp[i-1][j],     // Delete
                        dp[i][j-1],     // Insert
                        dp[i-1][j-1]    // Replace
                    });
                }
            }
        }
        
        return dp[m][n];
    }
    
    // Space optimized
    static int minDistanceOptimized(string word1, string word2) {
        int m = word1.length();
        int n = word2.length();
        
        vector<int> prev(n + 1);
        vector<int> curr(n + 1);
        
        // Initialize first row
        for (int j = 0; j <= n; j++) {
            prev[j] = j;
        }
        
        for (int i = 1; i <= m; i++) {
            curr[0] = i;
            for (int j = 1; j <= n; j++) {
                if (word1[i-1] == word2[j-1]) {
                    curr[j] = prev[j-1];
                } else {
                    curr[j] = 1 + min({prev[j], curr[j-1], prev[j-1]});
                }
            }
            prev = curr;
        }
        
        return curr[n];
    }
};
```

### 2. Palindromic Subsequences
```cpp
class PalindromicSubsequences {
public:
    // Longest palindromic subsequence
    static int longestPalindromeSubseq(string s) {
        int n = s.length();
        vector<vector<int>> dp(n, vector<int>(n, 0));
        
        // Every single character is a palindrome
        for (int i = 0; i < n; i++) {
            dp[i][i] = 1;
        }
        
        // Fill for substrings of length 2 to n
        for (int len = 2; len <= n; len++) {
            for (int i = 0; i <= n - len; i++) {
                int j = i + len - 1;
                
                if (s[i] == s[j]) {
                    dp[i][j] = dp[i+1][j-1] + 2;
                } else {
                    dp[i][j] = max(dp[i+1][j], dp[i][j-1]);
                }
            }
        }
        
        return dp[0][n-1];
    }
    
    // Count palindromic subsequences
    static int countPalindromicSubseq(string s) {
        int n = s.length();
        vector<vector<int>> dp(n, vector<int>(n, 0));
        
        // Every single character is a palindrome
        for (int i = 0; i < n; i++) {
            dp[i][i] = 1;
        }
        
        // Fill for substrings of length 2 to n
        for (int len = 2; len <= n; len++) {
            for (int i = 0; i <= n - len; i++) {
                int j = i + len - 1;
                
                if (s[i] == s[j]) {
                    dp[i][j] = dp[i+1][j] + dp[i][j-1] + 1;
                } else {
                    dp[i][j] = dp[i+1][j] + dp[i][j-1] - dp[i+1][j-1];
                }
            }
        }
        
        return dp[0][n-1];
    }
    
    // Minimum insertions to make palindrome
    static int minInsertions(string s) {
        int n = s.length();
        int lps = longestPalindromeSubseq(s);
        return n - lps;
    }
};
```

### 3. Word Break
```cpp
class WordBreak {
public:
    static bool wordBreak(string s, vector<string>& wordDict) {
        unordered_set<string> wordSet(wordDict.begin(), wordDict.end());
        int n = s.length();
        vector<bool> dp(n + 1, false);
        dp[0] = true;
        
        for (int i = 1; i <= n; i++) {
            for (int j = 0; j < i; j++) {
                if (dp[j] && wordSet.count(s.substr(j, i - j))) {
                    dp[i] = true;
                    break;
                }
            }
        }
        
        return dp[n];
    }
    
    // Word Break II - return all possible sentences
    static vector<string> wordBreakII(string s, vector<string>& wordDict) {
        unordered_set<string> wordSet(wordDict.begin(), wordDict.end());
        unordered_map<string, vector<string>> memo;
        return wordBreakHelper(s, wordSet, memo);
    }
    
private:
    static vector<string> wordBreakHelper(string s, unordered_set<string>& wordSet, 
                                         unordered_map<string, vector<string>>& memo) {
        if (memo.count(s)) {
            return memo[s];
        }
        
        vector<string> result;
        
        if (wordSet.count(s)) {
            result.push_back(s);
        }
        
        for (int i = 1; i < s.length(); i++) {
            string prefix = s.substr(0, i);
            if (wordSet.count(prefix)) {
                string suffix = s.substr(i);
                vector<string> suffixBreaks = wordBreakHelper(suffix, wordSet, memo);
                
                for (string suffixBreak : suffixBreaks) {
                    result.push_back(prefix + " " + suffixBreak);
                }
            }
        }
        
        memo[s] = result;
        return result;
    }
};
```

---

## Tree DP {#tree-dp}

### 1. Binary Tree Maximum Path Sum
```cpp
struct TreeNode {
    int val;
    TreeNode *left;
    TreeNode *right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode *left, TreeNode *right) : val(x), left(left), right(right) {}
};

class TreeDP {
public:
    static int maxPathSum(TreeNode* root) {
        int maxSum = INT_MIN;
        maxPathSumHelper(root, maxSum);
        return maxSum;
    }
    
    // House Robber III
    static int rob(TreeNode* root) {
        auto result = robHelper(root);
        return max(result.first, result.second);
    }
    
    // Binary Tree Diameter
    static int diameterOfBinaryTree(TreeNode* root) {
        int diameter = 0;
        diameterHelper(root, diameter);
        return diameter;
    }
    
private:
    static int maxPathSumHelper(TreeNode* node, int& maxSum) {
        if (!node) return 0;
        
        int leftMax = max(0, maxPathSumHelper(node->left, maxSum));
        int rightMax = max(0, maxPathSumHelper(node->right, maxSum));
        
        int currentMax = node->val + leftMax + rightMax;
        maxSum = max(maxSum, currentMax);
        
        return node->val + max(leftMax, rightMax);
    }
    
    // Returns {rob_current, not_rob_current}
    static pair<int, int> robHelper(TreeNode* node) {
        if (!node) return {0, 0};
        
        auto left = robHelper(node->left);
        auto right = robHelper(node->right);
        
        int robCurrent = node->val + left.second + right.second;
        int notRobCurrent = max(left.first, left.second) + max(right.first, right.second);
        
        return {robCurrent, notRobCurrent};
    }
    
    static int diameterHelper(TreeNode* node, int& diameter) {
        if (!node) return 0;
        
        int leftHeight = diameterHelper(node->left, diameter);
        int rightHeight = diameterHelper(node->right, diameter);
        
        diameter = max(diameter, leftHeight + rightHeight);
        
        return 1 + max(leftHeight, rightHeight);
    }
};
```

---

## Advanced DP Techniques {#advanced-techniques}

### 1. Matrix Chain Multiplication
```cpp
class MatrixChainMultiplication {
public:
    static int matrixChainOrder(vector<int>& dimensions) {
        int n = dimensions.size() - 1; // Number of matrices
        vector<vector<int>> dp(n, vector<int>(n, 0));
        
        // len is chain length
        for (int len = 2; len <= n; len++) {
            for (int i = 0; i <= n - len; i++) {
                int j = i + len - 1;
                dp[i][j] = INT_MAX;
                
                for (int k = i; k < j; k++) {
                    int cost = dp[i][k] + dp[k+1][j] + 
                              dimensions[i] * dimensions[k+1] * dimensions[j+1];
                    dp[i][j] = min(dp[i][j], cost);
                }
            }
        }
        
        return dp[0][n-1];
    }
    
    // Get optimal parenthesization
    static string getOptimalParentheses(vector<int>& dimensions) {
        int n = dimensions.size() - 1;
        vector<vector<int>> dp(n, vector<int>(n, 0));
        vector<vector<int>> split(n, vector<int>(n, 0));
        
        for (int len = 2; len <= n; len++) {
            for (int i = 0; i <= n - len; i++) {
                int j = i + len - 1;
                dp[i][j] = INT_MAX;
                
                for (int k = i; k < j; k++) {
                    int cost = dp[i][k] + dp[k+1][j] + 
                              dimensions[i] * dimensions[k+1] * dimensions[j+1];
                    if (cost < dp[i][j]) {
                        dp[i][j] = cost;
                        split[i][j] = k;
                    }
                }
            }
        }
        
        return buildParentheses(split, 0, n-1);
    }
    
private:
    static string buildParentheses(vector<vector<int>>& split, int i, int j) {
        if (i == j) {
            return "M" + to_string(i);
        }
        
        return "(" + buildParentheses(split, i, split[i][j]) + 
               " x " + buildParentheses(split, split[i][j] + 1, j) + ")";
    }
};
```

### 2. Knapsack Problems
```cpp
class KnapsackProblems {
public:
    // 0/1 Knapsack
    static int knapsack01(vector<int>& weights, vector<int>& values, int capacity) {
        int n = weights.size();
        vector<vector<int>> dp(n + 1, vector<int>(capacity + 1, 0));
        
        for (int i = 1; i <= n; i++) {
            for (int w = 1; w <= capacity; w++) {
                if (weights[i-1] <= w) {
                    dp[i][w] = max(dp[i-1][w], 
                                  dp[i-1][w - weights[i-1]] + values[i-1]);
                } else {
                    dp[i][w] = dp[i-1][w];
                }
            }
        }
        
        return dp[n][capacity];
    }
    
    // Space optimized 0/1 Knapsack
    static int knapsack01Optimized(vector<int>& weights, vector<int>& values, int capacity) {
        vector<int> dp(capacity + 1, 0);
        
        for (int i = 0; i < weights.size(); i++) {
            for (int w = capacity; w >= weights[i]; w--) {
                dp[w] = max(dp[w], dp[w - weights[i]] + values[i]);
            }
        }
        
        return dp[capacity];
    }
    
    // Unbounded Knapsack
    static int knapsackUnbounded(vector<int>& weights, vector<int>& values, int capacity) {
        vector<int> dp(capacity + 1, 0);
        
        for (int w = 1; w <= capacity; w++) {
            for (int i = 0; i < weights.size(); i++) {
                if (weights[i] <= w) {
                    dp[w] = max(dp[w], dp[w - weights[i]] + values[i]);
                }
            }
        }
        
        return dp[capacity];
    }
    
    // Get items in 0/1 Knapsack
    static vector<int> getKnapsackItems(vector<int>& weights, vector<int>& values, int capacity) {
        int n = weights.size();
        vector<vector<int>> dp(n + 1, vector<int>(capacity + 1, 0));
        
        // Fill DP table
        for (int i = 1; i <= n; i++) {
            for (int w = 1; w <= capacity; w++) {
                if (weights[i-1] <= w) {
                    dp[i][w] = max(dp[i-1][w], 
                                  dp[i-1][w - weights[i-1]] + values[i-1]);
                } else {
                    dp[i][w] = dp[i-1][w];
                }
            }
        }
        
        // Backtrack to find items
        vector<int> items;
        int w = capacity;
        for (int i = n; i > 0 && w > 0; i--) {
            if (dp[i][w] != dp[i-1][w]) {
                items.push_back(i-1);
                w -= weights[i-1];
            }
        }
        
        reverse(items.begin(), items.end());
        return items;
    }
};
```

---

## DP Optimization Techniques {#optimization}

### 1. Space Optimization
```cpp
class SpaceOptimization {
public:
    // Example: LCS with O(min(m,n)) space
    static int lcsSpaceOptimized(string text1, string text2) {
        // Make text1 the shorter string
        if (text1.length() > text2.length()) {
            swap(text1, text2);
        }
        
        int m = text1.length();
        int n = text2.length();
        
        vector<int> prev(m + 1, 0);
        vector<int> curr(m + 1, 0);
        
        for (int j = 1; j <= n; j++) {
            for (int i = 1; i <= m; i++) {
                if (text1[i-1] == text2[j-1]) {
                    curr[i] = prev[i-1] + 1;
                } else {
                    curr[i] = max(prev[i], curr[i-1]);
                }
            }
            prev = curr;
        }
        
        return curr[m];
    }
};
```

### 2. Memoization with Custom Hash
```cpp
class MemoizationOptimization {
public:
    // Custom hash for pair<int, int>
    struct PairHash {
        size_t operator()(const pair<int, int>& p) const {
            return hash<int>()(p.first) ^ (hash<int>()(p.second) << 1);
        }
    };
    
    static int uniquePathsMemo(int m, int n) {
        unordered_map<pair<int, int>, int, PairHash> memo;
        return uniquePathsHelper(m-1, n-1, memo);
    }
    
private:
    static int uniquePathsHelper(int i, int j, 
                                unordered_map<pair<int, int>, int, PairHash>& memo) {
        if (i == 0 || j == 0) return 1;
        
        pair<int, int> key = {i, j};
        if (memo.count(key)) {
            return memo[key];
        }
        
        memo[key] = uniquePathsHelper(i-1, j, memo) + uniquePathsHelper(i, j-1, memo);
        return memo[key];
    }
};
```

---

## Practice Problems {#practice-problems}

### Basic Problems
1. **Fibonacci Numbers** - Classic DP introduction
2. **Climbing Stairs** - Simple state transition
3. **House Robber** - Linear DP with constraints
4. **Maximum Subarray** - Kadane's algorithm
5. **Coin Change** - Unbounded knapsack variant
6. **Unique Paths** - 2D grid DP
7. **Minimum Path Sum** - 2D DP with optimization
8. **Longest Increasing Subsequence** - Classic LIS problem

### Intermediate Problems
1. **Edit Distance** - String DP
2. **Longest Common Subsequence** - 2D string DP
3. **Word Break** - String segmentation
4. **Palindromic Subsequence** - Interval DP
5. **0/1 Knapsack** - Classic optimization problem
6. **Partition Equal Subset Sum** - Subset sum variant
7. **Target Sum** - Count ways DP
8. **Decode Ways** - Linear DP with conditions

### Advanced Problems
1. **Matrix Chain Multiplication** - Interval DP
2. **Burst Balloons** - Interval DP
3. **Regular Expression Matching** - 2D string DP
4. **Interleaving String** - 2D string DP
5. **Distinct Subsequences** - String counting DP
6. **Maximum Rectangle** - 2D optimization
7. **Wildcard Matching** - Pattern matching DP
8. **Russian Doll Envelopes** - 2D LIS variant

### Tree DP Problems
1. **Binary Tree Maximum Path Sum**
2. **House Robber III**
3. **Binary Tree Diameter**
4. **Longest Univalue Path**
5. **Binary Tree Cameras**

---

## Complete Example Program {#example-program}

```cpp
#include <iostream>
#include <vector>
#include <string>
#include <algorithm>
#include <unordered_map>
#include <unordered_set>
#include <climits>
#include <chrono>

using namespace std;
using namespace chrono;

// All the classes defined above would be included here...
// (FibonacciDP, ClimbingStairs, HouseRobber, etc.)

class DPBenchmark {
public:
    static void benchmarkDPAlgorithms() {
        cout << "\n=== Dynamic Programming Benchmarks ===" << endl;
        
        // Fibonacci benchmarks
        cout << "\n--- Fibonacci Sequence ---" << endl;
        int n = 40;
        
        auto start = high_resolution_clock::now();
        int result1 = FibonacciDP::fibMemo(n);
        auto end = high_resolution_clock::now();
        auto duration1 = duration_cast<microseconds>(end - start);
        cout << "Memoization: " << result1 << " (" << duration1.count() << " μs)" << endl;
        
        start = high_resolution_clock::now();
        int result2 = FibonacciDP::fibTab(n);
        end = high_resolution_clock::now();
        auto duration2 = duration_cast<microseconds>(end - start);
        cout << "Tabulation: " << result2 << " (" << duration2.count() << " μs)" << endl;
        
        start = high_resolution_clock::now();
        int result3 = FibonacciDP::fibOptimized(n);
        end = high_resolution_clock::now();
        auto duration3 = duration_cast<microseconds>(end - start);
        cout << "Optimized: " << result3 << " (" << duration3.count() << " μs)" << endl;
        
        // LIS benchmark
        cout << "\n--- Longest Increasing Subsequence ---" << endl;
        vector<int> nums = {10, 9, 2, 5, 3, 7, 101, 18};
        
        start = high_resolution_clock::now();
        int lis1 = LongestIncreasingSubsequence::lengthOfLIS(nums);
        end = high_resolution_clock::now();
        duration1 = duration_cast<microseconds>(end - start);
        cout << "O(n²) DP: " << lis1 << " (" << duration1.count() << " μs)" << endl;
        
        start = high_resolution_clock::now();
        int lis2 = LongestIncreasingSubsequence::lengthOfLISOptimized(nums);
        end = high_resolution_clock::now();
        duration2 = duration_cast<microseconds>(end - start);
        cout << "O(n log n): " << lis2 << " (" << duration2.count() << " μs)" << endl;
        
        // String DP benchmark
        cout << "\n--- String DP Problems ---" << endl;
        string text1 = "abcde";
        string text2 = "ace";
        
        start = high_resolution_clock::now();
        int lcs = LongestCommonSubsequence::longestCommonSubsequence(text1, text2);
        end = high_resolution_clock::now();
        duration1 = duration_cast<microseconds>(end - start);
        cout << "LCS: " << lcs << " (" << duration1.count() << " μs)" << endl;
        
        string word1 = "horse";
        string word2 = "ros";
        
        start = high_resolution_clock::now();
        int editDist = EditDistance::minDistance(word1, word2);
        end = high_resolution_clock::now();
        duration2 = duration_cast<microseconds>(end - start);
        cout << "Edit Distance: " << editDist << " (" << duration2.count() << " μs)" << endl;
    }
};

int main() {
    cout << "=== Dynamic Programming Comprehensive Demo ===" << endl;
    
    // Basic DP demonstrations
    cout << "\n--- Basic DP Problems ---" << endl;
    
    // Fibonacci
    int n = 10;
    cout << "Fibonacci(" << n << "): " << FibonacciDP::fibOptimized(n) << endl;
    
    // Climbing Stairs
    int stairs = 5;
    cout << "Ways to climb " << stairs << " stairs: " << ClimbingStairs::climbStairs(stairs) << endl;
    
    // House Robber
    vector<int> houses = {2, 7, 9, 3, 1};
    cout << "Max money from houses: " << HouseRobber::robOptimized(houses) << endl;
    
    // Maximum Subarray
    vector<int> nums = {-2, 1, -3, 4, -1, 2, 1, -5, 4};
    cout << "Maximum subarray sum: " << MaximumSubarray::maxSubArray(nums) << endl;
    
    // 1D DP demonstrations
    cout << "\n--- 1D DP Problems ---" << endl;
    
    // Coin Change
    vector<int> coins = {1, 3, 4};
    int amount = 6;
    cout << "Min coins for amount " << amount << ": " << CoinChange::coinChange(coins, amount) << endl;
    
    // LIS
    vector<int> sequence = {10, 9, 2, 5, 3, 7, 101, 18};
    cout << "Length of LIS: " << LongestIncreasingSubsequence::lengthOfLIS(sequence) << endl;
    
    // 2D DP demonstrations
    cout << "\n--- 2D DP Problems ---" << endl;
    
    // Unique Paths
    int m = 3, n = 7;
    cout << "Unique paths (" << m << "x" << n << "): " << UniquePaths::uniquePaths(m, n) << endl;
    
    // LCS
    string text1 = "abcde";
    string text2 = "ace";
    cout << "LCS length: " << LongestCommonSubsequence::longestCommonSubsequence(text1, text2) << endl;
    cout << "LCS string: " << LongestCommonSubsequence::getLCS(text1, text2) << endl;
    
    // String DP demonstrations
    cout << "\n--- String DP Problems ---" << endl;
    
    // Edit Distance
    string word1 = "horse";
    string word2 = "ros";
    cout << "Edit distance: " << EditDistance::minDistance(word1, word2) << endl;
    
    // Palindromic Subsequence
    string s = "bbbab";
    cout << "Longest palindromic subsequence: " << PalindromicSubsequences::longestPalindromeSubseq(s) << endl;
    
    // Word Break
    string str = "leetcode";
    vector<string> wordDict = {"leet", "code"};
    cout << "Word break possible: " << (WordBreak::wordBreak(str, wordDict) ? "Yes" : "No") << endl;
    
    // Advanced DP demonstrations
    cout << "\n--- Advanced DP Problems ---" << endl;
    
    // Matrix Chain Multiplication
    vector<int> dimensions = {1, 2, 3, 4};
    cout << "Min matrix multiplications: " << MatrixChainMultiplication::matrixChainOrder(dimensions) << endl;
    
    // Knapsack
    vector<int> weights = {10, 20, 30};
    vector<int> values = {60, 100, 120};
    int capacity = 50;
    cout << "0/1 Knapsack max value: " << KnapsackProblems::knapsack01(weights, values, capacity) << endl;
    
    // Performance benchmarking
    DPBenchmark::benchmarkDPAlgorithms();
    
    return 0;
}
```

---

## Summary

### Key Concepts Covered
1. **Memoization vs Tabulation**: Top-down vs bottom-up approaches
2. **Space Optimization**: Reducing space complexity
3. **State Definition**: Identifying DP states and transitions
4. **Base Cases**: Proper initialization
5. **Reconstruction**: Getting actual solutions, not just optimal values

### Time and Space Complexities

#### Common DP Problems
| Problem | Time Complexity | Space Complexity | Optimized Space |
|---------|----------------|------------------|------------------|
| Fibonacci | O(n) | O(n) | O(1) |
| Climbing Stairs | O(n) | O(n) | O(1) |
| House Robber | O(n) | O(n) | O(1) |
| Coin Change | O(n × m) | O(n) | O(n) |
| LIS | O(n²) / O(n log n) | O(n) | O(n) |
| Unique Paths | O(m × n) | O(m × n) | O(min(m,n)) |
| LCS | O(m × n) | O(m × n) | O(min(m,n)) |
| Edit Distance | O(m × n) | O(m × n) | O(min(m,n)) |
| 0/1 Knapsack | O(n × W) | O(n × W) | O(W) |
| Matrix Chain | O(n³) | O(n²) | O(n²) |

### DP Problem Solving Strategy
1. **Identify if it's a DP problem**:
   - Optimal substructure
   - Overlapping subproblems
   - Optimization/counting problem

2. **Define the state**:
   - What parameters uniquely identify a subproblem?
   - What's the meaning of dp[i] or dp[i][j]?

3. **Find the recurrence relation**:
   - How to compute dp[i] from previous states?
   - What are the choices at each step?

4. **Identify base cases**:
   - What are the simplest subproblems?
   - How to initialize the DP table?

5. **Determine the order of computation**:
   - Bottom-up or top-down?
   - What's the dependency between states?

6. **Optimize space if possible**:
   - Can we reduce dimensions?
   - Do we need the entire DP table?

### When to Use Different DP Approaches
- **Memoization**: When the recurrence is complex or not all subproblems are needed
- **Tabulation**: When all subproblems need to be solved and order is clear
- **Space Optimization**: When only recent states are needed for computation

### Next Steps
Proceed to **Part 10: Advanced Algorithms** to learn about more sophisticated algorithmic techniques including greedy algorithms, divide and conquer, and advanced graph algorithms.

### Additional Resources
- [Dynamic Programming Patterns](https://leetcode.com/discuss/general-discussion/458695/dynamic-programming-patterns)
- [DP Practice Problems](https://leetcode.com/tag/dynamic-programming/)
- [Competitive Programming DP](https://cp-algorithms.com/dynamic_programming/)
- [DP Visualization](https://algorithm-visualizer.org/dynamic-programming/)