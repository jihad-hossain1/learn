# Practice Problems: Linked List Operations

Let me explain and provide practice solutions for these essential linked list problems.

## 1. Detect Cycle in Linked List

### Problem Understanding:
Determine if a linked list has a cycle (where a node's next pointer points to a previous node).

**Example:**
```
Input: 1 → 2 → 3 → 4 → 2 (cycle back to node 2)
Output: true
```

### Solutions:

```javascript
class ListNode {
    constructor(val, next = null) {
        this.val = val;
        this.next = next;
    }
}

// Method 1: Floyd's Cycle Detection (Two Pointers - Optimal)
function hasCycle(head) {
    if (!head || !head.next) return false;
    
    let slow = head;
    let fast = head;
    
    while (fast && fast.next) {
        slow = slow.next;          // Move one step
        fast = fast.next.next;     // Move two steps
        
        // If they meet, cycle exists
        if (slow === fast) return true;
    }
    
    return false; // Fast reached null, no cycle
}

// Method 2: Using Hash Set (O(n) space)
function hasCycleHash(head) {
    const visited = new Set();
    let current = head;
    
    while (current) {
        if (visited.has(current)) return true;
        visited.add(current);
        current = current.next;
    }
    
    return false;
}

// Method 3: Find Cycle Start Node
function detectCycleStart(head) {
    if (!head || !head.next) return null;
    
    let slow = head;
    let fast = head;
    let hasCycle = false;
    
    // Detect cycle
    while (fast && fast.next) {
        slow = slow.next;
        fast = fast.next.next;
        
        if (slow === fast) {
            hasCycle = true;
            break;
        }
    }
    
    if (!hasCycle) return null;
    
    // Find cycle start
    slow = head;
    while (slow !== fast) {
        slow = slow.next;
        fast = fast.next;
    }
    
    return slow; // Cycle start node
}

// Helper function to create a cycle for testing
function createLinkedListWithCycle(values, cycleIndex) {
    if (values.length === 0) return null;
    
    const head = new ListNode(values[0]);
    let current = head;
    let cycleNode = null;
    let count = 0;
    
    // Create nodes
    for (let i = 1; i < values.length; i++) {
        current.next = new ListNode(values[i]);
        current = current.next;
        
        if (count === cycleIndex) cycleNode = current;
        count++;
    }
    
    // Create cycle
    if (cycleIndex >= 0) {
        current.next = cycleNode;
    }
    
    return head;
}

// Test cases
console.log("Detect Cycle Tests:");
const cycleList = createLinkedListWithCycle([1, 2, 3, 4, 5], 2);
console.log(hasCycle(cycleList)); // true

const noCycleList = createLinkedListWithCycle([1, 2, 3, 4, 5], -1);
console.log(hasCycle(noCycleList)); // false

const cycleStart = detectCycleStart(cycleList);
console.log(cycleStart?.val); // 3
```

**Time Complexity:** O(n)  
**Space Complexity:** O(1) for Floyd's, O(n) for hash set

---

## 2. Merge Two Sorted Lists

### Problem Understanding:
Merge two sorted linked lists into one sorted linked list.

**Example:**
```
Input: 1 → 3 → 5 and 2 → 4 → 6
Output: 1 → 2 → 3 → 4 → 5 → 6
```

### Solutions:

```javascript
// Method 1: Iterative (Optimal)
function mergeTwoLists(l1, l2) {
    // Create dummy node to simplify code
    const dummy = new ListNode(-1);
    let current = dummy;
    
    while (l1 && l2) {
        if (l1.val <= l2.val) {
            current.next = l1;
            l1 = l1.next;
        } else {
            current.next = l2;
            l2 = l2.next;
        }
        current = current.next;
    }
    
    // Attach remaining nodes
    current.next = l1 || l2;
    
    return dummy.next;
}

// Method 2: Recursive (Elegant but O(n) space)
function mergeTwoListsRecursive(l1, l2) {
    // Base cases
    if (!l1) return l2;
    if (!l2) return l1;
    
    if (l1.val <= l2.val) {
        l1.next = mergeTwoListsRecursive(l1.next, l2);
        return l1;
    } else {
        l2.next = mergeTwoListsRecursive(l1, l2.next);
        return l2;
    }
}

// Method 3: In-place (Modify original lists)
function mergeTwoListsInPlace(l1, l2) {
    if (!l1) return l2;
    if (!l2) return l1;
    
    // Ensure l1 starts with smaller value
    if (l1.val > l2.val) {
        [l1, l2] = [l2, l1];
    }
    
    let head = l1;
    
    while (l1 && l2) {
        let temp = null;
        
        // Find the point where l2 should be inserted
        while (l1 && l1.val <= l2.val) {
            temp = l1;
            l1 = l1.next;
        }
        
        // Insert l2 node
        temp.next = l2;
        
        // Swap l1 and l2
        [l1, l2] = [l2, l1];
    }
    
    return head;
}

// Helper function to create linked list from array
function createLinkedList(arr) {
    if (arr.length === 0) return null;
    
    const head = new ListNode(arr[0]);
    let current = head;
    
    for (let i = 1; i < arr.length; i++) {
        current.next = new ListNode(arr[i]);
        current = current.next;
    }
    
    return head;
}

// Helper function to convert linked list to array
function linkedListToArray(head) {
    const result = [];
    let current = head;
    
    while (current) {
        result.push(current.val);
        current = current.next;
    }
    
    return result;
}

// Test cases
console.log("\nMerge Two Sorted Lists Tests:");
const list1 = createLinkedList([1, 3, 5]);
const list2 = createLinkedList([2, 4, 6]);

const merged = mergeTwoLists(list1, list2);
console.log(linkedListToArray(merged)); // [1, 2, 3, 4, 5, 6]

const mergedRecursive = mergeTwoListsRecursive(
    createLinkedList([1, 5, 9]),
    createLinkedList([2, 3, 10])
);
console.log(linkedListToArray(mergedRecursive)); // [1, 2, 3, 5, 9, 10]
```

**Time Complexity:** O(n + m)  
**Space Complexity:** O(1) for iterative, O(n + m) for recursive

---

## 3. Remove Nth Node From End

### Problem Understanding:
Remove the nth node from the end of the list and return the head.

**Example:**
```
Input: 1 → 2 → 3 → 4 → 5, n = 2
Output: 1 → 2 → 3 → 5
```

### Solutions:

```javascript
// Method 1: Two Pass (Find length first)
function removeNthFromEndTwoPass(head, n) {
    // Calculate length of list
    let length = 0;
    let current = head;
    
    while (current) {
        length++;
        current = current.next;
    }
    
    // Find position from start
    const positionFromStart = length - n;
    
    // If removing head
    if (positionFromStart === 0) {
        return head.next;
    }
    
    // Find node before the one to remove
    current = head;
    for (let i = 0; i < positionFromStart - 1; i++) {
        current = current.next;
    }
    
    // Remove the node
    current.next = current.next?.next || null;
    
    return head;
}

// Method 2: One Pass with Two Pointers (Optimal)
function removeNthFromEnd(head, n) {
    const dummy = new ListNode(-1);
    dummy.next = head;
    
    let fast = dummy;
    let slow = dummy;
    
    // Move fast pointer n+1 steps ahead
    for (let i = 0; i <= n; i++) {
        fast = fast.next;
    }
    
    // Move both pointers until fast reaches end
    while (fast) {
        slow = slow.next;
        fast = fast.next;
    }
    
    // Remove the nth node from end
    slow.next = slow.next.next;
    
    return dummy.next;
}

// Method 3: Using Stack
function removeNthFromEndStack(head, n) {
    const stack = [];
    const dummy = new ListNode(-1);
    dummy.next = head;
    let current = dummy;
    
    // Push all nodes to stack
    while (current) {
        stack.push(current);
        current = current.next;
    }
    
    // Pop n nodes to find the node before target
    for (let i = 0; i < n; i++) {
        stack.pop();
    }
    
    // Remove the target node
    const nodeBeforeTarget = stack.pop();
    nodeBeforeTarget.next = nodeBeforeTarget.next.next;
    
    return dummy.next;
}

// Test cases
console.log("\nRemove Nth From End Tests:");
const list = createLinkedList([1, 2, 3, 4, 5]);

console.log("Original:", linkedListToArray(list));
console.log("Remove 2nd from end:", linkedListToArray(removeNthFromEnd(list, 2))); // [1, 2, 3, 5]

const list2 = createLinkedList([1]);
console.log("Remove from single node:", linkedListToArray(removeNthFromEnd(list2, 1))); // []

const list3 = createLinkedList([1, 2]);
console.log("Remove head:", linkedListToArray(removeNthFromEnd(list3, 2))); // [2]
```

**Time Complexity:** O(n)  
**Space Complexity:** O(1) for two pointers, O(n) for stack

---

## 4. Palindrome Linked List

### Problem Understanding:
Check if a linked list is a palindrome.

**Example:**
```
Input: 1 → 2 → 2 → 1
Output: true
```

### Solutions:

```javascript
// Method 1: Reverse Second Half (Optimal - O(n) time, O(1) space)
function isPalindrome(head) {
    if (!head || !head.next) return true;
    
    // Step 1: Find middle using slow/fast pointers
    let slow = head;
    let fast = head;
    
    while (fast && fast.next) {
        slow = slow.next;
        fast = fast.next.next;
    }
    
    // Step 2: Reverse second half
    let secondHalf = reverseList(slow);
    let firstHalf = head;
    let copySecondHalf = secondHalf; // Save for restoring later
    
    // Step 3: Compare both halves
    let isPal = true;
    while (secondHalf) {
        if (firstHalf.val !== secondHalf.val) {
            isPal = false;
            break;
        }
        firstHalf = firstHalf.next;
        secondHalf = secondHalf.next;
    }
    
    // Step 4: Restore the list (optional)
    reverseList(copySecondHalf);
    
    return isPal;
}

// Helper function to reverse linked list
function reverseList(head) {
    let prev = null;
    let current = head;
    
    while (current) {
        const next = current.next;
        current.next = prev;
        prev = current;
        current = next;
    }
    
    return prev;
}

// Method 2: Using Stack (O(n) space)
function isPalindromeStack(head) {
    if (!head || !head.next) return true;
    
    const stack = [];
    let slow = head;
    let fast = head;
    
    // Push first half to stack
    while (fast && fast.next) {
        stack.push(slow.val);
        slow = slow.next;
        fast = fast.next.next;
    }
    
    // If odd number of nodes, skip middle
    if (fast) {
        slow = slow.next;
    }
    
    // Compare second half with stack
    while (slow) {
        if (stack.pop() !== slow.val) return false;
        slow = slow.next;
    }
    
    return true;
}

// Method 3: Recursive (O(n) space)
function isPalindromeRecursive(head) {
    let frontPointer = head;
    
    function recursivelyCheck(current) {
        if (current) {
            if (!recursivelyCheck(current.next)) return false;
            if (current.val !== frontPointer.val) return false;
            frontPointer = frontPointer.next;
        }
        return true;
    }
    
    return recursivelyCheck(head);
}

// Method 4: Convert to Array and Check
function isPalindromeArray(head) {
    const values = [];
    let current = head;
    
    // Convert to array
    while (current) {
        values.push(current.val);
        current = current.next;
    }
    
    // Check if array is palindrome
    let left = 0, right = values.length - 1;
    while (left < right) {
        if (values[left] !== values[right]) return false;
        left++;
        right--;
    }
    
    return true;
}

// Test cases
console.log("\nPalindrome Linked List Tests:");
console.log(isPalindrome(createLinkedList([1, 2, 2, 1]))); // true
console.log(isPalindrome(createLinkedList([1, 2, 3, 2, 1]))); // true
console.log(isPalindrome(createLinkedList([1, 2]))); // false
console.log(isPalindrome(createLinkedList([1]))); // true

console.log("Stack method:", isPalindromeStack(createLinkedList([1, 2, 3, 2, 1]))); // true
console.log("Array method:", isPalindromeArray(createLinkedList([1, 2, 1]))); // true
```

**Time Complexity:** O(n) for all methods  
**Space Complexity:** O(1) for reverse method, O(n) for others

---

## 🎯 Practice Exercises

### Exercise 1: Cycle Detection Variations
```javascript
// 1. Find cycle length
function cycleLength(head) {
    // Your code here
}

// 2. Check if cycle exists and return cycle nodes
function getCycleNodes(head) {
    // Your code here
}
```

### Exercise 2: Merge Lists Variations
```javascript
// 1. Merge K sorted lists
function mergeKLists(lists) {
    // Your code here
}

// 2. Merge two lists in alternating fashion
function mergeAlternating(l1, l2) {
    // Your code here
}
```

### Exercise 3: Remove Node Variations
```javascript
// 1. Remove all occurrences of a value
function removeElements(head, val) {
    // Your code here
}

// 2. Remove duplicates from sorted list
function deleteDuplicates(head) {
    // Your code here
}
```

### Exercise 4: Palindrome Variations
```javascript
// 1. Check if linked list can be palindrome after removing one node
function validPalindromeList(head) {
    // Your code here
}

// 2. Find longest palindromic subsequence in linked list
function longestPalindromeSubsequence(head) {
    // Your code here
}
```

## 🚀 Key Patterns Learned:

1. **Fast & Slow Pointers**: Cycle detection, find middle
2. **Dummy Nodes**: Simplify edge cases in list operations
3. **List Reversal**: Palindrome checking, list manipulation
4. **Multiple Passes**: Length calculation, position finding
5. **Recursive vs Iterative**: Different tradeoffs for list problems

These linked list problems are fundamental for technical interviews. Practice until you can implement them confidently!