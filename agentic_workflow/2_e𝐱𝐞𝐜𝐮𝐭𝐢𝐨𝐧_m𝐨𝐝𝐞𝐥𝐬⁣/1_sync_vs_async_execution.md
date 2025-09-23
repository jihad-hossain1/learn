# Synchronous vs. Asynchronous Execution in JavaScript

Execution models define how code runs and handles tasks. JavaScript is single-threaded but uses asynchronous execution to handle non-blocking operations efficiently.

## Core Concepts

### **Synchronous Execution**
- **Sequential**: Code runs line by line, in order
- **Blocking**: Each operation must complete before moving to the next
- **Predictable**: Easy to reason about execution flow
- **Simple**: Straightforward error handling

### **Asynchronous Execution**
- **Non-blocking**: Operations start and complete later without blocking the main thread
- **Event-driven**: Uses callbacks, promises, async/await
- **Efficient**: Handles I/O operations without waiting
- **Complex**: Requires understanding of event loop and callbacks

---

## JavaScript Examples

### Example 1: Basic Synchronous vs Asynchronous

```javascript
console.log("=== SYNCHRONOUS EXECUTION ===");

function synchronousExample() {
    console.log("1. Start synchronous function");
    
    // Blocking operation simulation
    const result = expensiveCalculation(1000000);
    console.log("2. Calculation result:", result);
    
    console.log("3. End synchronous function");
}

function expensiveCalculation(iterations) {
    let sum = 0;
    for (let i = 0; i < iterations; i++) {
        sum += Math.sqrt(i) * Math.random();
    }
    return sum;
}

// This will block the main thread
synchronousExample();
console.log("4. After synchronous function"); // This waits for the function to complete

console.log("\n=== ASYNCHRONOUS EXECUTION ===");

function asynchronousExample() {
    console.log("1. Start asynchronous function");
    
    // Non-blocking operation
    setTimeout(() => {
        console.log("3. Async operation completed after 2 seconds");
    }, 2000);
    
    console.log("2. End asynchronous function (but operation continues in background)");
}

asynchronousExample();
console.log("4. This runs immediately after async function starts"); // Doesn't wait for setTimeout
```

### Example 2: Real-World Comparison

```javascript
// Synchronous file reading (Node.js)
const fs = require('fs');

console.log("=== SYNCHRONOUS FILE READING ===");

function readFilesSynchronously() {
    console.log("1. Starting synchronous file read");
    
    try {
        const file1 = fs.readFileSync('file1.txt', 'utf8');
        console.log("2. File 1 content length:", file1.length);
        
        const file2 = fs.readFileSync('file2.txt', 'utf8');
        console.log("3. File 2 content length:", file2.length);
        
        const file3 = fs.readFileSync('file3.txt', 'utf8');
        console.log("4. File 3 content length:", file3.length);
    } catch (error) {
        console.error("Error reading files:", error);
    }
    
    console.log("5. All files read synchronously");
}

// Asynchronous file reading
function readFilesAsynchronously() {
    console.log("1. Starting asynchronous file read");
    
    fs.readFile('file1.txt', 'utf8', (err, data1) => {
        if (err) {
            console.error("Error reading file1:", err);
            return;
        }
        console.log("2. File 1 content length:", data1.length);
        
        // Chain next read
        fs.readFile('file2.txt', 'utf8', (err, data2) => {
            if (err) {
                console.error("Error reading file2:", err);
                return;
            }
            console.log("3. File 2 content length:", data2.length);
            
            fs.readFile('file3.txt', 'utf8', (err, data3) => {
                if (err) {
                    console.error("Error reading file3:", err);
                    return;
                }
                console.log("4. File 3 content length:", data3.length);
                console.log("5. All files read asynchronously");
            });
        });
    });
    
    console.log("6. Async function continues immediately");
}

// readFilesSynchronously(); // This would block
// readFilesAsynchronously(); // This doesn't block
```

### Example 3: Callback Hell vs Promises vs Async/Await

```javascript
console.log("=== EVOLUTION OF ASYNC PATTERNS ===");

// 1. Callback Hell (Pyramid of Doom)
function callbackHellExample(userId, callback) {
    console.log("1. Starting callback hell...");
    
    getUser(userId, (err, user) => {
        if (err) return callback(err);
        console.log("2. Got user:", user.name);
        
        getPosts(user.id, (err, posts) => {
            if (err) return callback(err);
            console.log("3. Got posts:", posts.length);
            
            getComments(posts[0].id, (err, comments) => {
                if (err) return callback(err);
                console.log("4. Got comments:", comments.length);
                
                // Nested callbacks create "callback hell"
                callback(null, { user, posts, comments });
            });
        });
    });
}

// 2. Promises (ES6)
function promiseExample(userId) {
    console.log("1. Starting promise chain...");
    
    return getUserPromise(userId)
        .then(user => {
            console.log("2. Got user:", user.name);
            return getPostsPromise(user.id);
        })
        .then(posts => {
            console.log("3. Got posts:", posts.length);
            return getCommentsPromise(posts[0].id);
        })
        .then(comments => {
            console.log("4. Got comments:", comments.length);
            return { comments };
        })
        .catch(error => {
            console.error("Promise chain error:", error);
            throw error;
        });
}

// 3. Async/Await (ES2017) - Most readable
async function asyncAwaitExample(userId) {
    try {
        console.log("1. Starting async/await...");
        
        const user = await getUserPromise(userId);
        console.log("2. Got user:", user.name);
        
        const posts = await getPostsPromise(user.id);
        console.log("3. Got posts:", posts.length);
        
        const comments = await getCommentsPromise(posts[0].id);
        console.log("4. Got comments:", comments.length);
        
        return { user, posts, comments };
    } catch (error) {
        console.error("Async/await error:", error);
        throw error;
    }
}

// Mock implementations
function getUser(id, callback) {
    setTimeout(() => callback(null, { id, name: "John Doe" }), 100);
}

function getPosts(userId, callback) {
    setTimeout(() => callback(null, [{ id: 1, title: "Post 1" }]), 100);
}

function getComments(postId, callback) {
    setTimeout(() => callback(null, [{ id: 1, text: "Great post!" }]), 100);
}

// Promise versions
function getUserPromise(id) {
    return new Promise((resolve) => {
        setTimeout(() => resolve({ id, name: "John Doe" }), 100);
    });
}

function getPostsPromise(userId) {
    return new Promise((resolve) => {
        setTimeout(() => resolve([{ id: 1, title: "Post 1" }]), 100);
    });
}

function getCommentsPromise(postId) {
    return new Promise((resolve) => {
        setTimeout(() => resolve([{ id: 1, text: "Great post!" }]), 100);
    });
}

// Test the examples
async function testExamples() {
    console.log("=== TESTING ASYNC PATTERNS ===\n");
    
    // Callback hell
    callbackHellExample(1, (err, result) => {
        if (err) console.error(err);
        else console.log("Callback hell result:", result);
    });
    
    // Promises
    promiseExample(1)
        .then(result => console.log("Promise result:", result))
        .catch(err => console.error(err));
    
    // Async/await
    try {
        const result = await asyncAwaitExample(1);
        console.log("Async/await result:", result);
    } catch (err) {
        console.error(err);
    }
}

// testExamples();
```

### Example 4: Understanding the Event Loop

```javascript
console.log("=== JAVASCRIPT EVENT LOOP DEMONSTRATION ===");

function eventLoopDemo() {
    console.log("1. Start of script");
    
    // Synchronous code
    const x = 10;
    console.log("2. Synchronous variable:", x);
    
    // Macro-task: setTimeout
    setTimeout(() => {
        console.log("7. setTimeout callback (macro-task)");
    }, 0);
    
    // Micro-task: Promise
    Promise.resolve().then(() => {
        console.log("5. Promise resolution (micro-task)");
    });
    
    // Another micro-task
    queueMicrotask(() => {
        console.log("6. queueMicrotask (micro-task)");
    });
    
    // Synchronous function
    function syncFunction() {
        console.log("3. Inside synchronous function");
    }
    syncFunction();
    
    console.log("4. End of synchronous code");
}

eventLoopDemo();

// Output order explanation:
// 1. Start of script
// 2. Synchronous variable: 10
// 3. Inside synchronous function
// 4. End of synchronous code
// 5. Promise resolution (micro-task)  - Micro-tasks run before macro-tasks
// 6. queueMicrotask (micro-task)      - All micro-tasks complete first
// 7. setTimeout callback (macro-task) - Then macro-tasks run
```

### Example 5: Practical Async Patterns

```javascript
class AsyncOperations {
    // 1. Sequential execution (one after another)
    static async sequentialOperations() {
        console.log("=== SEQUENTIAL EXECUTION ===");
        
        const start = Date.now();
        
        // These run one after another
        const user = await this.fetchUser(1);
        const posts = await this.fetchPosts(user.id);
        const comments = await this.fetchComments(posts[0].id);
        
        const end = Date.now();
        console.log(`Sequential completed in ${end - start}ms`);
        
        return { user, posts, comments };
    }
    
    // 2. Parallel execution (all at once)
    static async parallelOperations() {
        console.log("=== PARALLEL EXECUTION ===");
        
        const start = Date.now();
        
        // These start simultaneously
        const [user, settings, notifications] = await Promise.all([
            this.fetchUser(1),
            this.fetchUserSettings(1),
            this.fetchNotifications(1)
        ]);
        
        const end = Date.now();
        console.log(`Parallel completed in ${end - start}ms`);
        
        return { user, settings, notifications };
    }
    
    // 3. Race conditions - which completes first?
    static async raceOperations() {
        console.log("=== RACE EXECUTION ===");
        
        const winner = await Promise.race([
            this.fetchFastAPI(),
            this.fetchSlowAPI(),
            this.timeout(500) // Reject if too slow
        ]);
        
        console.log("Race winner:", winner);
        return winner;
    }
    
    // 4. Error handling in parallel
    static async parallelWithErrorHandling() {
        console.log("=== PARALLEL WITH ERROR HANDLING ===");
        
        try {
            const results = await Promise.allSettled([
                this.fetchUser(1),
                this.fetchUser(999), // Might fail
                this.fetchPosts(1)
            ]);
            
            const successful = results.filter(r => r.status === 'fulfilled');
            const failed = results.filter(r => r.status === 'rejected');
            
            console.log(`Successful: ${successful.length}, Failed: ${failed.length}`);
            
            return successful.map(r => r.value);
        } catch (error) {
            console.error("Parallel operation failed:", error);
            throw error;
        }
    }
    
    // 5. Async generators for streaming data
    static async *asyncGenerator() {
        console.log("=== ASYNC GENERATOR ===");
        
        for (let i = 1; i <= 5; i++) {
            // Simulate async data fetching
            const data = await this.fetchItem(i);
            yield data;
        }
    }
    
    // Mock API calls
    static fetchUser(id) {
        return new Promise(resolve => {
            setTimeout(() => resolve({ id, name: `User ${id}` }), 100);
        });
    }
    
    static fetchPosts(userId) {
        return new Promise(resolve => {
            setTimeout(() => resolve([{ id: 1, title: "Post 1" }]), 150);
        });
    }
    
    static fetchComments(postId) {
        return new Promise(resolve => {
            setTimeout(() => resolve([{ id: 1, text: "Comment 1" }]), 200);
        });
    }
    
    static fetchUserSettings(userId) {
        return new Promise(resolve => {
            setTimeout(() => resolve({ theme: "dark", language: "en" }), 120);
        });
    }
    
    static fetchNotifications(userId) {
        return new Promise(resolve => {
            setTimeout(() => resolve([{ id: 1, message: "Notification 1" }]), 180);
        });
    }
    
    static fetchFastAPI() {
        return new Promise(resolve => {
            setTimeout(() => resolve("Fast API response"), 100);
        });
    }
    
    static fetchSlowAPI() {
        return new Promise(resolve => {
            setTimeout(() => resolve("Slow API response"), 1000);
        });
    }
    
    static timeout(ms) {
        return new Promise((_, reject) => {
            setTimeout(() => reject(new Error("Timeout")), ms);
        });
    }
    
    static fetchItem(id) {
        return new Promise(resolve => {
            setTimeout(() => resolve(`Item ${id}`), 50);
        });
    }
}

// Usage examples
async function runAsyncExamples() {
    console.log("=== RUNNING ASYNC EXAMPLES ===\n");
    
    // Sequential vs Parallel comparison
    await AsyncOperations.sequentialOperations();
    await AsyncOperations.parallelOperations();
    
    // Race condition
    try {
        await AsyncOperations.raceOperations();
    } catch (error) {
        console.error("Race failed:", error.message);
    }
    
    // Error handling
    await AsyncOperations.parallelWithErrorHandling();
    
    // Async generator
    for await (const item of AsyncOperations.asyncGenerator()) {
        console.log("Generated:", item);
    }
}

// runAsyncExamples();
```

### Example 6: Advanced Async Patterns

```javascript
// Advanced async patterns for production applications
class AdvancedAsyncPatterns {
    // 1. Async semaphore for limiting concurrent operations
    static createSemaphore(maxConcurrent) {
        let current = 0;
        const queue = [];
        
        return async (fn) => {
            if (current >= maxConcurrent) {
                await new Promise(resolve => queue.push(resolve));
            }
            
            current++;
            try {
                return await fn();
            } finally {
                current--;
                if (queue.length > 0) {
                    queue.shift()(); // Release next waiting function
                }
            }
        };
    }
    
    // 2. Retry with exponential backoff
    static async retryWithBackoff(operation, maxRetries = 3, baseDelay = 1000) {
        for (let attempt = 1; attempt <= maxRetries; attempt++) {
            try {
                return await operation();
            } catch (error) {
                if (attempt === maxRetries) throw error;
                
                const delay = baseDelay * Math.pow(2, attempt - 1);
                console.log(`Attempt ${attempt} failed, retrying in ${delay}ms`);
                await this.delay(delay + Math.random() * 1000); // Jitter
            }
        }
    }
    
    // 3. Async queue for processing tasks sequentially
    static createAsyncQueue() {
        const queue = [];
        let isProcessing = false;
        
        const processQueue = async () => {
            if (isProcessing || queue.length === 0) return;
            
            isProcessing = true;
            const { operation, resolve, reject } = queue.shift();
            
            try {
                const result = await operation();
                resolve(result);
            } catch (error) {
                reject(error);
            } finally {
                isProcessing = false;
                processQueue();
            }
        };
        
        return (operation) => {
            return new Promise((resolve, reject) => {
                queue.push({ operation, resolve, reject });
                processQueue();
            });
        };
    }
    
    // 4. Debounced async function
    static debounceAsync(fn, delay) {
        let timeoutId;
        let pendingPromise = null;
        let latestArgs;
        
        return async (...args) => {
            latestArgs = args;
            
            if (timeoutId) {
                clearTimeout(timeoutId);
            }
            
            if (!pendingPromise) {
                pendingPromise = new Promise((resolve) => {
                    timeoutId = setTimeout(async () => {
                        try {
                            const result = await fn(...latestArgs);
                            resolve(result);
                        } finally {
                            pendingPromise = null;
                        }
                    }, delay);
                });
            }
            
            return pendingPromise;
        };
    }
    
    // 5. Async pool for limiting concurrent promises
    static async promisePool(tasks, concurrency) {
        const results = [];
        const executing = new Set();
        
        for (const [index, task] of tasks.entries()) {
            const promise = task();
            results[index] = promise;
            executing.add(promise);
            
            const clean = () => executing.delete(promise);
            promise.then(clean).catch(clean);
            
            if (executing.size >= concurrency) {
                await Promise.race(executing);
            }
        }
        
        return Promise.all(results);
    }
    
    static delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Usage examples
async function advancedPatternsDemo() {
    console.log("=== ADVANCED ASYNC PATTERNS ===\n");
    
    // Semaphore example - limit to 2 concurrent operations
    const semaphore = AdvancedAsyncPatterns.createSemaphore(2);
    
    const tasks = Array.from({ length: 5 }, (_, i) => 
        semaphore(() => {
            console.log(`Task ${i} started`);
            return AdvancedAsyncPatterns.delay(1000).then(() => {
                console.log(`Task ${i} completed`);
                return i;
            });
        })
    );
    
    await Promise.all(tasks);
    console.log("All tasks completed with semaphore\n");
    
    // Retry with backoff example
    let attempt = 0;
    const flakyOperation = () => {
        attempt++;
        if (attempt < 3) {
            return Promise.reject(new Error("Temporary failure"));
        }
        return Promise.resolve("Success!");
    };
    
    try {
        const result = await AdvancedAsyncPatterns.retryWithBackoff(flakyOperation);
        console.log("Retry success:", result);
    } catch (error) {
        console.error("Retry failed:", error.message);
    }
}

// advancedPatternsDemo();
```

## Key Differences Summary

| Aspect | Synchronous | Asynchronous |
|--------|-------------|--------------|
| **Execution** | Sequential, blocking | Concurrent, non-blocking |
| **Performance** | Can block UI | Better responsiveness |
| **Error Handling** | Try/catch | .catch() or try/catch with async/await |
| **Complexity** | Simple, linear | More complex, requires understanding of event loop |
| **Use Cases** | Simple calculations, immediate operations | I/O operations, API calls, timers |

## When to Use Each

### **Use Synchronous When:**
- Operations are fast and CPU-bound
- You need simple, predictable code flow
- Operations must happen in specific order
- Dealing with small datasets

### **Use Asynchronous When:**
- Performing I/O operations (file system, network requests)
- Dealing with user interactions and events
- Handling multiple concurrent operations
- Building responsive applications

## JavaScript Event Loop Priority:
1. **Synchronous code** (highest priority)
2. **Micro-tasks** (Promises, queueMicrotask)
3. **Macro-tasks** (setTimeout, setInterval, I/O)
4. **Animation frames** (requestAnimationFrame)

Understanding synchronous vs asynchronous execution is fundamental to writing efficient, responsive JavaScript applications.