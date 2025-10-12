# Practice Problems: Stack-Based Challenges

Let me explain and provide practice solutions for these stack-based problems that teach important algorithmic patterns.

## 1. Valid Parentheses

### Problem Understanding:
Check if a string containing parentheses `()`, brackets `[]`, and braces `{}` is valid (properly closed and nested).

**Example:**
```
Input: "()[]{}"
Output: true

Input: "([)]"
Output: false
```

### Solutions:

```javascript
// Method 1: Stack (Optimal - O(n) time, O(n) space)
function isValid(s) {
    const stack = [];
    const mapping = {
        ')': '(',
        ']': '[',
        '}': '{'
    };
    
    for (let char of s) {
        if (char === '(' || char === '[' || char === '{') {
            // Opening bracket - push to stack
            stack.push(char);
        } else {
            // Closing bracket - check if matches top of stack
            if (stack.length === 0 || stack.pop() !== mapping[char]) {
                return false;
            }
        }
    }
    
    // Stack should be empty if all brackets are matched
    return stack.length === 0;
}

// Method 2: More Explicit Version
function isValidExplicit(s) {
    const stack = [];
    
    for (let char of s) {
        if (char === '(' || char === '[' || char === '{') {
            stack.push(char);
        } else if (char === ')' && stack.length > 0 && stack[stack.length - 1] === '(') {
            stack.pop();
        } else if (char === ']' && stack.length > 0 && stack[stack.length - 1] === '[') {
            stack.pop();
        } else if (char === '}' && stack.length > 0 && stack[stack.length - 1] === '{') {
            stack.pop();
        } else {
            return false; // Invalid closing bracket
        }
    }
    
    return stack.length === 0;
}

// Method 3: Using Switch Statement
function isValidSwitch(s) {
    const stack = [];
    
    for (let char of s) {
        switch (char) {
            case '(':
            case '[':
            case '{':
                stack.push(char);
                break;
            case ')':
                if (stack.pop() !== '(') return false;
                break;
            case ']':
                if (stack.pop() !== '[') return false;
                break;
            case '}':
                if (stack.pop() !== '{') return false;
                break;
        }
    }
    
    return stack.length === 0;
}

// Enhanced version that tracks position for error reporting
function isValidWithPosition(s) {
    const stack = []; // Store both character and position
    const mapping = { ')': '(', ']': '[', '}': '{' };
    
    for (let i = 0; i < s.length; i++) {
        const char = s[i];
        
        if (char === '(' || char === '[' || char === '{') {
            stack.push({ char, position: i });
        } else if (char in mapping) {
            if (stack.length === 0) {
                return { isValid: false, error: `Unmatched closing ${char} at position ${i}` };
            }
            
            const last = stack.pop();
            if (last.char !== mapping[char]) {
                return { isValid: false, error: `Mismatch: ${last.char} at ${last.position} with ${char} at ${i}` };
            }
        }
    }
    
    if (stack.length > 0) {
        const last = stack.pop();
        return { isValid: false, error: `Unmatched opening ${last.char} at position ${last.position}` };
    }
    
    return { isValid: true, error: null };
}

// Test cases
console.log("Valid Parentheses Tests:");
console.log(isValid("()")); // true
console.log(isValid("()[]{}")); // true
console.log(isValid("(]")); // false
console.log(isValid("([)]")); // false
console.log(isValid("{[]}")); // true
console.log(isValid("")); // true

console.log(isValidWithPosition("({)}")); // { isValid: false, error: "Mismatch: { at 1 with ) at 2" }
```

**Time Complexity:** O(n)  
**Space Complexity:** O(n) for the stack

---

## 2. Implement Queue Using Stacks

### Problem Understanding:
Implement a queue (FIFO) using only stacks (LIFO).

**Operations:**
- `push(x)`: Add element to back
- `pop()`: Remove element from front  
- `peek()`: Get front element
- `empty()`: Check if queue is empty

### Solutions:

```javascript
// Method 1: Two Stacks - Amortized O(1) operations
class MyQueue {
    constructor() {
        this.inputStack = [];  // For push operations
        this.outputStack = []; // For pop/peek operations
    }
    
    push(x) {
        this.inputStack.push(x);
    }
    
    pop() {
        if (this.outputStack.length === 0) {
            this._transfer();
        }
        return this.outputStack.pop();
    }
    
    peek() {
        if (this.outputStack.length === 0) {
            this._transfer();
        }
        return this.outputStack[this.outputStack.length - 1];
    }
    
    empty() {
        return this.inputStack.length === 0 && this.outputStack.length === 0;
    }
    
    _transfer() {
        // Move all elements from input to output stack (reversing order)
        while (this.inputStack.length > 0) {
            this.outputStack.push(this.inputStack.pop());
        }
    }
    
    // Helper to see both stacks
    getStacks() {
        return {
            input: [...this.inputStack],
            output: [...this.outputStack]
        };
    }
}

// Method 2: Two Stacks with Different Approach
class MyQueueAlternative {
    constructor() {
        this.stack1 = [];
        this.stack2 = [];
    }
    
    push(x) {
        // Move all elements from stack1 to stack2
        while (this.stack1.length > 0) {
            this.stack2.push(this.stack1.pop());
        }
        
        // Push new element to stack1
        this.stack1.push(x);
        
        // Move all elements back to stack1
        while (this.stack2.length > 0) {
            this.stack1.push(this.stack2.pop());
        }
    }
    
    pop() {
        return this.stack1.pop();
    }
    
    peek() {
        return this.stack1[this.stack1.length - 1];
    }
    
    empty() {
        return this.stack1.length === 0;
    }
}

// Method 3: Using JavaScript Array as Queue (for comparison)
class SimpleQueue {
    constructor() {
        this.items = [];
    }
    
    push(x) {
        this.items.push(x);
    }
    
    pop() {
        return this.items.shift(); // O(n) operation!
    }
    
    peek() {
        return this.items[0];
    }
    
    empty() {
        return this.items.length === 0;
    }
}

// Test cases
console.log("\nQueue Using Stacks Tests:");
const queue = new MyQueue();
queue.push(1);
queue.push(2);
queue.push(3);

console.log(queue.peek()); // 1
console.log(queue.pop());  // 1
console.log(queue.empty()); // false
console.log(queue.pop());  // 2
console.log(queue.pop());  // 3
console.log(queue.empty()); // true

// Performance comparison
function testPerformance() {
    const n = 10000;
    
    console.time('MyQueue');
    const myQueue = new MyQueue();
    for (let i = 0; i < n; i++) myQueue.push(i);
    for (let i = 0; i < n; i++) myQueue.pop();
    console.timeEnd('MyQueue');
    
    console.time('SimpleQueue');
    const simpleQueue = new SimpleQueue();
    for (let i = 0; i < n; i++) simpleQueue.push(i);
    for (let i = 0; i < n; i++) simpleQueue.pop();
    console.timeEnd('SimpleQueue');
}

// testPerformance(); // Uncomment to see performance difference
```

**Time Complexity:**
- **MyQueue:** Amortized O(1) for all operations
- **MyQueueAlternative:** O(n) for push, O(1) for pop/peek
- **SimpleQueue:** O(1) for push, O(n) for pop

---

## 3. Daily Temperatures

### Problem Understanding:
For each day, find how many days you have to wait until a warmer temperature. If no warmer day, return 0.

**Example:**
```
Input: [73, 74, 75, 71, 69, 72, 76, 73]
Output: [1, 1, 4, 2, 1, 1, 0, 0]
```

### Solutions:

```javascript
// Method 1: Monotonic Decreasing Stack (Optimal - O(n))
function dailyTemperatures(temperatures) {
    const result = new Array(temperatures.length).fill(0);
    const stack = []; // Store indices of temperatures
    
    for (let i = 0; i < temperatures.length; i++) {
        const currentTemp = temperatures[i];
        
        // While current temp is warmer than stack top
        while (stack.length > 0 && temperatures[stack[stack.length - 1]] < currentTemp) {
            const colderDayIndex = stack.pop();
            result[colderDayIndex] = i - colderDayIndex;
        }
        
        stack.push(i);
    }
    
    return result;
}

// Method 2: Brute Force (O(n²) - for comparison)
function dailyTemperaturesBrute(temperatures) {
    const result = [];
    
    for (let i = 0; i < temperatures.length; i++) {
        let days = 0;
        let found = false;
        
        for (let j = i + 1; j < temperatures.length; j++) {
            if (temperatures[j] > temperatures[i]) {
                days = j - i;
                found = true;
                break;
            }
        }
        
        result.push(found ? days : 0);
    }
    
    return result;
}

// Method 3: Array Approach (Optimized brute force)
function dailyTemperaturesArray(temperatures) {
    const n = temperatures.length;
    const result = new Array(n).fill(0);
    const hottest = Math.max(...temperatures);
    
    for (let i = n - 1; i >= 0; i--) {
        if (temperatures[i] >= hottest) {
            result[i] = 0;
            continue;
        }
        
        let days = 1;
        while (temperatures[i + days] <= temperatures[i]) {
            days += result[i + days];
        }
        result[i] = days;
    }
    
    return result;
}

// Enhanced version with visualization
function dailyTemperaturesWithExplanation(temperatures) {
    const result = new Array(temperatures.length).fill(0);
    const stack = [];
    const steps = [];
    
    for (let i = 0; i < temperatures.length; i++) {
        const currentTemp = temperatures[i];
        
        while (stack.length > 0 && temperatures[stack[stack.length - 1]] < currentTemp) {
            const colderIndex = stack.pop();
            result[colderIndex] = i - colderIndex;
            
            steps.push({
                currentIndex: i,
                currentTemp,
                poppedIndex: colderIndex,
                daysToWait: i - colderIndex,
                stack: [...stack],
                result: [...result]
            });
        }
        
        stack.push(i);
        steps.push({
            currentIndex: i,
            currentTemp,
            action: 'pushed to stack',
            stack: [...stack],
            result: [...result]
        });
    }
    
    return { result, steps };
}

// Test cases
console.log("\nDaily Temperatures Tests:");
const temps = [73, 74, 75, 71, 69, 72, 76, 73];
console.log(dailyTemperatures(temps)); // [1, 1, 4, 2, 1, 1, 0, 0]

const explanation = dailyTemperaturesWithExplanation([73, 74, 75, 71, 69, 72, 76, 73]);
console.log("Final result:", explanation.result);

// Compare with brute force
console.log("Brute force:", dailyTemperaturesBrute(temps));
```

**Time Complexity:** O(n) with stack, O(n²) brute force  
**Space Complexity:** O(n) for stack and result

---

## 4. Min Stack

### Problem Understanding:
Design a stack that supports push, pop, top, and retrieving the minimum element in constant time.

**Operations:**
- `push(x)`: Add element
- `pop()`: Remove top element  
- `top()`: Get top element
- `getMin()`: Get minimum element

### Solutions:

```javascript
// Method 1: Two Stacks (Optimal)
class MinStack {
    constructor() {
        this.stack = [];     // Main stack
        this.minStack = [];  // Stack to track minimums
    }
    
    push(x) {
        this.stack.push(x);
        
        // If minStack is empty or x <= current min, push to minStack
        if (this.minStack.length === 0 || x <= this.getMin()) {
            this.minStack.push(x);
        }
    }
    
    pop() {
        if (this.stack.length === 0) return undefined;
        
        const popped = this.stack.pop();
        
        // If popped element is current min, remove from minStack
        if (popped === this.getMin()) {
            this.minStack.pop();
        }
        
        return popped;
    }
    
    top() {
        return this.stack[this.stack.length - 1];
    }
    
    getMin() {
        return this.minStack[this.minStack.length - 1];
    }
}

// Method 2: Single Stack with Tuples
class MinStackSingle {
    constructor() {
        this.stack = []; // Each element is [value, currentMin]
    }
    
    push(x) {
        let currentMin = x;
        if (this.stack.length > 0) {
            currentMin = Math.min(x, this.getMin());
        }
        this.stack.push([x, currentMin]);
    }
    
    pop() {
        if (this.stack.length === 0) return undefined;
        return this.stack.pop()[0];
    }
    
    top() {
        return this.stack[this.stack.length - 1][0];
    }
    
    getMin() {
        return this.stack[this.stack.length - 1][1];
    }
}

// Method 3: Using Linked List
class MinNode {
    constructor(val, min, next = null) {
        this.val = val;
        this.min = min; // Minimum so far
        this.next = next;
    }
}

class MinStackLinkedList {
    constructor() {
        this.head = null;
    }
    
    push(x) {
        if (this.head === null) {
            this.head = new MinNode(x, x);
        } else {
            this.head = new MinNode(x, Math.min(x, this.getMin()), this.head);
        }
    }
    
    pop() {
        if (this.head === null) return undefined;
        const val = this.head.val;
        this.head = this.head.next;
        return val;
    }
    
    top() {
        return this.head?.val;
    }
    
    getMin() {
        return this.head?.min;
    }
}

// Enhanced version with tracking
class TrackingMinStack extends MinStack {
    constructor() {
        super();
        this.operations = [];
    }
    
    push(x) {
        super.push(x);
        this.operations.push(`Pushed ${x}, min is ${this.getMin()}`);
    }
    
    pop() {
        const val = super.pop();
        this.operations.push(`Popped ${val}, min is ${this.getMin()}`);
        return val;
    }
    
    getHistory() {
        return this.operations;
    }
}

// Test cases
console.log("\nMin Stack Tests:");
const minStack = new MinStack();
minStack.push(-2);
minStack.push(0);
minStack.push(-3);
console.log(minStack.getMin()); // -3
minStack.pop();
console.log(minStack.top());    // 0
console.log(minStack.getMin()); // -2

// Test with tracking
const trackingStack = new TrackingMinStack();
trackingStack.push(5);
trackingStack.push(3);
trackingStack.push(7);
trackingStack.pop();
trackingStack.push(1);
console.log(trackingStack.getHistory());

// Performance test
function testMinStackPerformance() {
    const n = 100000;
    
    console.time('TwoStacks');
    const stack1 = new MinStack();
    for (let i = 0; i < n; i++) stack1.push(Math.random());
    for (let i = 0; i < n; i++) stack1.getMin();
    console.timeEnd('TwoStacks');
    
    console.time('SingleStack');
    const stack2 = new MinStackSingle();
    for (let i = 0; i < n; i++) stack2.push(Math.random());
    for (let i = 0; i < n; i++) stack2.getMin();
    console.timeEnd('SingleStack');
}

// testMinStackPerformance(); // Uncomment to compare performance
```

**Time Complexity:** O(1) for all operations  
**Space Complexity:** O(n) for all methods

---

## 🎯 Practice Exercises

### Exercise 1: Parentheses Variations
```javascript
// 1. Generate all valid parentheses combinations
function generateParentheses(n) {
    // Your code here
}

// 2. Minimum removals to make parentheses valid
function minRemoveToMakeValid(s) {
    // Your code here
}
```

### Exercise 2: Queue Variations
```javascript
// 1. Implement stack using queues
class MyStack {
    // Your code here
}

// 2. Circular queue implementation
class CircularQueue {
    // Your code here
}
```

### Exercise 3: Temperature Variations
```javascript
// 1. Next greater element (general case)
function nextGreaterElement(nums) {
    // Your code here
}

// 2. Stock span problem
function stockSpan(prices) {
    // Your code here
}
```

### Exercise 4: Stack Variations
```javascript
// 1. Max stack (supports getMax like min stack)
class MaxStack {
    // Your code here
}

// 2. Stack that supports getMiddle in O(1)
class MiddleStack {
    // Your code here
}
```

## 🚀 Key Patterns Learned:

1. **Stack for LIFO Problems**: Parentheses matching, function calls
2. **Monotonic Stack**: Daily temperatures, next greater element
3. **Two Stack Technique**: Queue implementation, min stack
4. **State Tracking**: Maintaining additional information in stacks
5. **Amortized Analysis**: Understanding average case performance

These stack-based problems are fundamental for understanding how to use LIFO data structures to solve complex problems efficiently!