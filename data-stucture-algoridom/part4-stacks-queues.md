# Part 4: Stacks and Queues

## Table of Contents
1. [Introduction to Stacks](#introduction-stacks)
2. [Stack Implementation](#stack-implementation)
3. [Stack Applications](#stack-applications)
4. [Introduction to Queues](#introduction-queues)
5. [Queue Implementation](#queue-implementation)
6. [Queue Variations](#queue-variations)
7. [Priority Queues](#priority-queues)
8. [Applications and Problems](#applications-problems)
9. [Practice Exercises](#practice-exercises)

---

## 1. Introduction to Stacks {#introduction-stacks}

### What is a Stack?
A stack is a linear data structure that follows the **Last In, First Out (LIFO)** principle. Elements are added and removed from the same end, called the "top" of the stack.

### Stack Operations
- **Push**: Add an element to the top
- **Pop**: Remove and return the top element
- **Top/Peek**: Return the top element without removing it
- **isEmpty**: Check if stack is empty
- **Size**: Get number of elements

### Real-world Examples
- Function call stack
- Undo operations in editors
- Browser back button
- Expression evaluation
- Backtracking algorithms

### Time Complexity
- Push: O(1)
- Pop: O(1)
- Top: O(1)
- Search: O(n)

---

## 2. Stack Implementation {#stack-implementation}

### Array-based Stack

```cpp
#include <iostream>
#include <stdexcept>
using namespace std;

class ArrayStack {
private:
    int* arr;
    int topIndex;
    int capacity;
    
public:
    // Constructor
    ArrayStack(int size) {
        capacity = size;
        arr = new int[capacity];
        topIndex = -1;
    }
    
    // Destructor
    ~ArrayStack() {
        delete[] arr;
    }
    
    // Push operation
    void push(int value) {
        if (isFull()) {
            throw overflow_error("Stack overflow");
        }
        arr[++topIndex] = value;
    }
    
    // Pop operation
    int pop() {
        if (isEmpty()) {
            throw underflow_error("Stack underflow");
        }
        return arr[topIndex--];
    }
    
    // Peek/Top operation
    int top() {
        if (isEmpty()) {
            throw underflow_error("Stack is empty");
        }
        return arr[topIndex];
    }
    
    // Check if empty
    bool isEmpty() {
        return topIndex == -1;
    }
    
    // Check if full
    bool isFull() {
        return topIndex == capacity - 1;
    }
    
    // Get size
    int size() {
        return topIndex + 1;
    }
    
    // Display stack
    void display() {
        if (isEmpty()) {
            cout << "Stack is empty" << endl;
            return;
        }
        
        cout << "Stack (top to bottom): ";
        for (int i = topIndex; i >= 0; i--) {
            cout << arr[i] << " ";
        }
        cout << endl;
    }
};
```

### Linked List-based Stack

```cpp
struct StackNode {
    int data;
    StackNode* next;
    
    StackNode(int value) : data(value), next(nullptr) {}
};

class LinkedStack {
private:
    StackNode* topNode;
    int count;
    
public:
    // Constructor
    LinkedStack() {
        topNode = nullptr;
        count = 0;
    }
    
    // Destructor
    ~LinkedStack() {
        while (!isEmpty()) {
            pop();
        }
    }
    
    // Push operation
    void push(int value) {
        StackNode* newNode = new StackNode(value);
        newNode->next = topNode;
        topNode = newNode;
        count++;
    }
    
    // Pop operation
    int pop() {
        if (isEmpty()) {
            throw underflow_error("Stack underflow");
        }
        
        StackNode* temp = topNode;
        int value = temp->data;
        topNode = topNode->next;
        delete temp;
        count--;
        
        return value;
    }
    
    // Peek/Top operation
    int top() {
        if (isEmpty()) {
            throw underflow_error("Stack is empty");
        }
        return topNode->data;
    }
    
    // Check if empty
    bool isEmpty() {
        return topNode == nullptr;
    }
    
    // Get size
    int size() {
        return count;
    }
    
    // Display stack
    void display() {
        if (isEmpty()) {
            cout << "Stack is empty" << endl;
            return;
        }
        
        cout << "Stack (top to bottom): ";
        StackNode* current = topNode;
        while (current != nullptr) {
            cout << current->data << " ";
            current = current->next;
        }
        cout << endl;
    }
};
```

### Dynamic Array Stack (Resizable)

```cpp
class DynamicStack {
private:
    int* arr;
    int topIndex;
    int capacity;
    
    void resize() {
        int newCapacity = capacity * 2;
        int* newArr = new int[newCapacity];
        
        for (int i = 0; i <= topIndex; i++) {
            newArr[i] = arr[i];
        }
        
        delete[] arr;
        arr = newArr;
        capacity = newCapacity;
    }
    
public:
    DynamicStack(int initialSize = 10) {
        capacity = initialSize;
        arr = new int[capacity];
        topIndex = -1;
    }
    
    ~DynamicStack() {
        delete[] arr;
    }
    
    void push(int value) {
        if (topIndex == capacity - 1) {
            resize();
        }
        arr[++topIndex] = value;
    }
    
    int pop() {
        if (isEmpty()) {
            throw underflow_error("Stack underflow");
        }
        return arr[topIndex--];
    }
    
    int top() {
        if (isEmpty()) {
            throw underflow_error("Stack is empty");
        }
        return arr[topIndex];
    }
    
    bool isEmpty() {
        return topIndex == -1;
    }
    
    int size() {
        return topIndex + 1;
    }
};
```

---

## 3. Stack Applications {#stack-applications}

### Expression Evaluation

```cpp
#include <string>
#include <cctype>

class ExpressionEvaluator {
public:
    // Check if parentheses are balanced
    static bool isBalanced(const string& expression) {
        LinkedStack stack;
        
        for (char c : expression) {
            if (c == '(' || c == '[' || c == '{') {
                stack.push(c);
            } else if (c == ')' || c == ']' || c == '}') {
                if (stack.isEmpty()) {
                    return false;
                }
                
                char top = stack.pop();
                if ((c == ')' && top != '(') ||
                    (c == ']' && top != '[') ||
                    (c == '}' && top != '{')) {
                    return false;
                }
            }
        }
        
        return stack.isEmpty();
    }
    
    // Convert infix to postfix
    static string infixToPostfix(const string& infix) {
        LinkedStack stack;
        string postfix = "";
        
        for (char c : infix) {
            if (isalnum(c)) {
                postfix += c;
            } else if (c == '(') {
                stack.push(c);
            } else if (c == ')') {
                while (!stack.isEmpty() && stack.top() != '(') {
                    postfix += stack.pop();
                }
                if (!stack.isEmpty()) {
                    stack.pop(); // Remove '('
                }
            } else if (isOperator(c)) {
                while (!stack.isEmpty() && 
                       precedence(stack.top()) >= precedence(c)) {
                    postfix += stack.pop();
                }
                stack.push(c);
            }
        }
        
        while (!stack.isEmpty()) {
            postfix += stack.pop();
        }
        
        return postfix;
    }
    
    // Evaluate postfix expression
    static int evaluatePostfix(const string& postfix) {
        LinkedStack stack;
        
        for (char c : postfix) {
            if (isdigit(c)) {
                stack.push(c - '0');
            } else if (isOperator(c)) {
                int operand2 = stack.pop();
                int operand1 = stack.pop();
                int result = performOperation(operand1, operand2, c);
                stack.push(result);
            }
        }
        
        return stack.pop();
    }
    
private:
    static bool isOperator(char c) {
        return c == '+' || c == '-' || c == '*' || c == '/' || c == '^';
    }
    
    static int precedence(char op) {
        switch (op) {
            case '+': case '-': return 1;
            case '*': case '/': return 2;
            case '^': return 3;
            default: return 0;
        }
    }
    
    static int performOperation(int a, int b, char op) {
        switch (op) {
            case '+': return a + b;
            case '-': return a - b;
            case '*': return a * b;
            case '/': return a / b;
            case '^': return pow(a, b);
            default: return 0;
        }
    }
};
```

### Function Call Stack Simulation

```cpp
struct FunctionCall {
    string functionName;
    int lineNumber;
    
    FunctionCall(string name, int line) : functionName(name), lineNumber(line) {}
};

class CallStack {
private:
    stack<FunctionCall> callStack;
    
public:
    void enterFunction(const string& name, int line) {
        callStack.push(FunctionCall(name, line));
        cout << "Entering " << name << " at line " << line << endl;
    }
    
    void exitFunction() {
        if (!callStack.empty()) {
            FunctionCall call = callStack.top();
            callStack.pop();
            cout << "Exiting " << call.functionName << endl;
        }
    }
    
    void printCallStack() {
        cout << "Current call stack:" << endl;
        stack<FunctionCall> temp = callStack;
        
        while (!temp.empty()) {
            FunctionCall call = temp.top();
            temp.pop();
            cout << "  " << call.functionName << " (line " << call.lineNumber << ")" << endl;
        }
    }
};
```

---

## 4. Introduction to Queues {#introduction-queues}

### What is a Queue?
A queue is a linear data structure that follows the **First In, First Out (FIFO)** principle. Elements are added at the rear (enqueue) and removed from the front (dequeue).

### Queue Operations
- **Enqueue**: Add an element to the rear
- **Dequeue**: Remove and return the front element
- **Front**: Return the front element without removing it
- **Rear**: Return the rear element without removing it
- **isEmpty**: Check if queue is empty
- **Size**: Get number of elements

### Real-world Examples
- Print job scheduling
- CPU task scheduling
- Breadth-First Search
- Handling requests in web servers
- Call center systems

### Time Complexity
- Enqueue: O(1)
- Dequeue: O(1)
- Front/Rear: O(1)
- Search: O(n)

---

## 5. Queue Implementation {#queue-implementation}

### Array-based Queue (Circular)

```cpp
class CircularQueue {
private:
    int* arr;
    int front;
    int rear;
    int capacity;
    int count;
    
public:
    CircularQueue(int size) {
        capacity = size;
        arr = new int[capacity];
        front = 0;
        rear = -1;
        count = 0;
    }
    
    ~CircularQueue() {
        delete[] arr;
    }
    
    void enqueue(int value) {
        if (isFull()) {
            throw overflow_error("Queue overflow");
        }
        
        rear = (rear + 1) % capacity;
        arr[rear] = value;
        count++;
    }
    
    int dequeue() {
        if (isEmpty()) {
            throw underflow_error("Queue underflow");
        }
        
        int value = arr[front];
        front = (front + 1) % capacity;
        count--;
        
        return value;
    }
    
    int getFront() {
        if (isEmpty()) {
            throw underflow_error("Queue is empty");
        }
        return arr[front];
    }
    
    int getRear() {
        if (isEmpty()) {
            throw underflow_error("Queue is empty");
        }
        return arr[rear];
    }
    
    bool isEmpty() {
        return count == 0;
    }
    
    bool isFull() {
        return count == capacity;
    }
    
    int size() {
        return count;
    }
    
    void display() {
        if (isEmpty()) {
            cout << "Queue is empty" << endl;
            return;
        }
        
        cout << "Queue (front to rear): ";
        for (int i = 0; i < count; i++) {
            cout << arr[(front + i) % capacity] << " ";
        }
        cout << endl;
    }
};
```

### Linked List-based Queue

```cpp
struct QueueNode {
    int data;
    QueueNode* next;
    
    QueueNode(int value) : data(value), next(nullptr) {}
};

class LinkedQueue {
private:
    QueueNode* frontNode;
    QueueNode* rearNode;
    int count;
    
public:
    LinkedQueue() {
        frontNode = nullptr;
        rearNode = nullptr;
        count = 0;
    }
    
    ~LinkedQueue() {
        while (!isEmpty()) {
            dequeue();
        }
    }
    
    void enqueue(int value) {
        QueueNode* newNode = new QueueNode(value);
        
        if (isEmpty()) {
            frontNode = rearNode = newNode;
        } else {
            rearNode->next = newNode;
            rearNode = newNode;
        }
        count++;
    }
    
    int dequeue() {
        if (isEmpty()) {
            throw underflow_error("Queue underflow");
        }
        
        QueueNode* temp = frontNode;
        int value = temp->data;
        
        frontNode = frontNode->next;
        if (frontNode == nullptr) {
            rearNode = nullptr;
        }
        
        delete temp;
        count--;
        
        return value;
    }
    
    int getFront() {
        if (isEmpty()) {
            throw underflow_error("Queue is empty");
        }
        return frontNode->data;
    }
    
    int getRear() {
        if (isEmpty()) {
            throw underflow_error("Queue is empty");
        }
        return rearNode->data;
    }
    
    bool isEmpty() {
        return frontNode == nullptr;
    }
    
    int size() {
        return count;
    }
    
    void display() {
        if (isEmpty()) {
            cout << "Queue is empty" << endl;
            return;
        }
        
        cout << "Queue (front to rear): ";
        QueueNode* current = frontNode;
        while (current != nullptr) {
            cout << current->data << " ";
            current = current->next;
        }
        cout << endl;
    }
};
```

---

## 6. Queue Variations {#queue-variations}

### Deque (Double-ended Queue)

```cpp
class Deque {
private:
    int* arr;
    int front;
    int rear;
    int capacity;
    int count;
    
public:
    Deque(int size) {
        capacity = size;
        arr = new int[capacity];
        front = 0;
        rear = capacity - 1;
        count = 0;
    }
    
    ~Deque() {
        delete[] arr;
    }
    
    // Insert at front
    void insertFront(int value) {
        if (isFull()) {
            throw overflow_error("Deque overflow");
        }
        
        front = (front - 1 + capacity) % capacity;
        arr[front] = value;
        count++;
    }
    
    // Insert at rear
    void insertRear(int value) {
        if (isFull()) {
            throw overflow_error("Deque overflow");
        }
        
        rear = (rear + 1) % capacity;
        arr[rear] = value;
        count++;
    }
    
    // Delete from front
    int deleteFront() {
        if (isEmpty()) {
            throw underflow_error("Deque underflow");
        }
        
        int value = arr[front];
        front = (front + 1) % capacity;
        count--;
        
        return value;
    }
    
    // Delete from rear
    int deleteRear() {
        if (isEmpty()) {
            throw underflow_error("Deque underflow");
        }
        
        int value = arr[rear];
        rear = (rear - 1 + capacity) % capacity;
        count--;
        
        return value;
    }
    
    int getFront() {
        if (isEmpty()) {
            throw underflow_error("Deque is empty");
        }
        return arr[front];
    }
    
    int getRear() {
        if (isEmpty()) {
            throw underflow_error("Deque is empty");
        }
        return arr[rear];
    }
    
    bool isEmpty() {
        return count == 0;
    }
    
    bool isFull() {
        return count == capacity;
    }
    
    int size() {
        return count;
    }
};
```

### Circular Queue with Dynamic Resizing

```cpp
class DynamicCircularQueue {
private:
    int* arr;
    int front;
    int rear;
    int capacity;
    int count;
    
    void resize() {
        int newCapacity = capacity * 2;
        int* newArr = new int[newCapacity];
        
        // Copy elements in order
        for (int i = 0; i < count; i++) {
            newArr[i] = arr[(front + i) % capacity];
        }
        
        delete[] arr;
        arr = newArr;
        front = 0;
        rear = count - 1;
        capacity = newCapacity;
    }
    
public:
    DynamicCircularQueue(int initialSize = 10) {
        capacity = initialSize;
        arr = new int[capacity];
        front = 0;
        rear = -1;
        count = 0;
    }
    
    ~DynamicCircularQueue() {
        delete[] arr;
    }
    
    void enqueue(int value) {
        if (count == capacity) {
            resize();
        }
        
        rear = (rear + 1) % capacity;
        arr[rear] = value;
        count++;
    }
    
    int dequeue() {
        if (isEmpty()) {
            throw underflow_error("Queue underflow");
        }
        
        int value = arr[front];
        front = (front + 1) % capacity;
        count--;
        
        return value;
    }
    
    bool isEmpty() {
        return count == 0;
    }
    
    int size() {
        return count;
    }
};
```

---

## 7. Priority Queues {#priority-queues}

### Simple Priority Queue (Array-based)

```cpp
struct PriorityElement {
    int data;
    int priority;
    
    PriorityElement(int d, int p) : data(d), priority(p) {}
};

class SimplePriorityQueue {
private:
    PriorityElement* arr;
    int size;
    int capacity;
    
public:
    SimplePriorityQueue(int cap) {
        capacity = cap;
        arr = new PriorityElement[capacity];
        size = 0;
    }
    
    ~SimplePriorityQueue() {
        delete[] arr;
    }
    
    void enqueue(int data, int priority) {
        if (size >= capacity) {
            throw overflow_error("Priority queue overflow");
        }
        
        // Insert in sorted order (higher priority first)
        int i = size - 1;
        while (i >= 0 && arr[i].priority < priority) {
            arr[i + 1] = arr[i];
            i--;
        }
        
        arr[i + 1] = PriorityElement(data, priority);
        size++;
    }
    
    int dequeue() {
        if (isEmpty()) {
            throw underflow_error("Priority queue underflow");
        }
        
        int data = arr[0].data;
        
        // Shift elements
        for (int i = 1; i < size; i++) {
            arr[i - 1] = arr[i];
        }
        
        size--;
        return data;
    }
    
    int peek() {
        if (isEmpty()) {
            throw underflow_error("Priority queue is empty");
        }
        return arr[0].data;
    }
    
    bool isEmpty() {
        return size == 0;
    }
    
    int getSize() {
        return size;
    }
    
    void display() {
        cout << "Priority Queue (priority, data): ";
        for (int i = 0; i < size; i++) {
            cout << "(" << arr[i].priority << ", " << arr[i].data << ") ";
        }
        cout << endl;
    }
};
```

### Heap-based Priority Queue

```cpp
class HeapPriorityQueue {
private:
    vector<PriorityElement> heap;
    
    void heapifyUp(int index) {
        while (index > 0) {
            int parent = (index - 1) / 2;
            if (heap[index].priority <= heap[parent].priority) {
                break;
            }
            swap(heap[index], heap[parent]);
            index = parent;
        }
    }
    
    void heapifyDown(int index) {
        int size = heap.size();
        
        while (true) {
            int largest = index;
            int left = 2 * index + 1;
            int right = 2 * index + 2;
            
            if (left < size && heap[left].priority > heap[largest].priority) {
                largest = left;
            }
            
            if (right < size && heap[right].priority > heap[largest].priority) {
                largest = right;
            }
            
            if (largest == index) {
                break;
            }
            
            swap(heap[index], heap[largest]);
            index = largest;
        }
    }
    
public:
    void enqueue(int data, int priority) {
        heap.push_back(PriorityElement(data, priority));
        heapifyUp(heap.size() - 1);
    }
    
    int dequeue() {
        if (isEmpty()) {
            throw underflow_error("Priority queue underflow");
        }
        
        int data = heap[0].data;
        heap[0] = heap.back();
        heap.pop_back();
        
        if (!heap.empty()) {
            heapifyDown(0);
        }
        
        return data;
    }
    
    int peek() {
        if (isEmpty()) {
            throw underflow_error("Priority queue is empty");
        }
        return heap[0].data;
    }
    
    bool isEmpty() {
        return heap.empty();
    }
    
    int size() {
        return heap.size();
    }
};
```

---

## 8. Applications and Problems {#applications-problems}

### Stack Problems

```cpp
class StackProblems {
public:
    // Next Greater Element
    static vector<int> nextGreaterElement(vector<int>& nums) {
        vector<int> result(nums.size(), -1);
        stack<int> st;
        
        for (int i = 0; i < nums.size(); i++) {
            while (!st.empty() && nums[st.top()] < nums[i]) {
                result[st.top()] = nums[i];
                st.pop();
            }
            st.push(i);
        }
        
        return result;
    }
    
    // Largest Rectangle in Histogram
    static int largestRectangleArea(vector<int>& heights) {
        stack<int> st;
        int maxArea = 0;
        
        for (int i = 0; i <= heights.size(); i++) {
            int h = (i == heights.size()) ? 0 : heights[i];
            
            while (!st.empty() && heights[st.top()] > h) {
                int height = heights[st.top()];
                st.pop();
                int width = st.empty() ? i : i - st.top() - 1;
                maxArea = max(maxArea, height * width);
            }
            
            st.push(i);
        }
        
        return maxArea;
    }
    
    // Valid Parentheses with multiple types
    static bool isValidParentheses(string s) {
        stack<char> st;
        unordered_map<char, char> mapping = {{')', '('}, {']', '['}, {'}', '{'}};
        
        for (char c : s) {
            if (mapping.count(c)) {
                if (st.empty() || st.top() != mapping[c]) {
                    return false;
                }
                st.pop();
            } else {
                st.push(c);
            }
        }
        
        return st.empty();
    }
};
```

### Queue Problems

```cpp
class QueueProblems {
public:
    // Sliding Window Maximum
    static vector<int> maxSlidingWindow(vector<int>& nums, int k) {
        deque<int> dq;
        vector<int> result;
        
        for (int i = 0; i < nums.size(); i++) {
            // Remove elements outside window
            while (!dq.empty() && dq.front() < i - k + 1) {
                dq.pop_front();
            }
            
            // Remove smaller elements
            while (!dq.empty() && nums[dq.back()] < nums[i]) {
                dq.pop_back();
            }
            
            dq.push_back(i);
            
            // Add to result if window is complete
            if (i >= k - 1) {
                result.push_back(nums[dq.front()]);
            }
        }
        
        return result;
    }
    
    // First non-repeating character in stream
    static vector<char> firstNonRepeating(string stream) {
        vector<char> result;
        queue<char> q;
        unordered_map<char, int> freq;
        
        for (char c : stream) {
            freq[c]++;
            q.push(c);
            
            // Remove characters with frequency > 1
            while (!q.empty() && freq[q.front()] > 1) {
                q.pop();
            }
            
            if (q.empty()) {
                result.push_back('#');
            } else {
                result.push_back(q.front());
            }
        }
        
        return result;
    }
    
    // Generate binary numbers from 1 to n
    static vector<string> generateBinary(int n) {
        vector<string> result;
        queue<string> q;
        q.push("1");
        
        for (int i = 0; i < n; i++) {
            string current = q.front();
            q.pop();
            result.push_back(current);
            
            q.push(current + "0");
            q.push(current + "1");
        }
        
        return result;
    }
};
```

---

## 9. Practice Exercises {#practice-exercises}

### Stack Exercises
1. **Min Stack**: Design a stack that supports push, pop, top, and retrieving minimum element in constant time
2. **Implement Queue using Stacks**: Use two stacks to implement queue operations
3. **Reverse Polish Notation**: Evaluate arithmetic expression in postfix notation
4. **Daily Temperatures**: Find next warmer temperature for each day
5. **Trapping Rain Water**: Calculate trapped rainwater using stack
6. **Remove K Digits**: Remove k digits to make smallest possible number
7. **Asteroid Collision**: Simulate asteroid collisions
8. **Decode String**: Decode string with nested brackets

### Queue Exercises
1. **Implement Stack using Queues**: Use two queues to implement stack operations
2. **Circular Tour**: Find starting point for circular tour
3. **Rotten Oranges**: Find minimum time to rot all oranges (BFS)
4. **Perfect Squares**: Find minimum number of perfect squares that sum to n
5. **Open the Lock**: Find minimum moves to open combination lock
6. **Walls and Gates**: Fill rooms with distance to nearest gate
7. **Snake and Ladder**: Find minimum moves to reach end
8. **Word Ladder**: Transform one word to another with minimum steps

### Priority Queue Exercises
1. **Kth Largest Element**: Find kth largest element in array
2. **Merge K Sorted Lists**: Merge k sorted linked lists
3. **Top K Frequent Elements**: Find k most frequent elements
4. **Meeting Rooms II**: Find minimum meeting rooms required
5. **Task Scheduler**: Schedule tasks with cooling period
6. **Ugly Number II**: Find nth ugly number
7. **Super Ugly Number**: Generalized ugly number problem
8. **Find Median from Data Stream**: Maintain median of streaming data

### Complete Example Program

```cpp
#include <iostream>
#include <stack>
#include <queue>
using namespace std;

int main() {
    cout << "=== Stack and Queue Demo ===" << endl;
    
    // Stack Demo
    cout << "\nStack Operations:" << endl;
    LinkedStack stack;
    
    stack.push(10);
    stack.push(20);
    stack.push(30);
    
    stack.display();
    
    cout << "Top element: " << stack.top() << endl;
    cout << "Popped: " << stack.pop() << endl;
    
    stack.display();
    
    // Queue Demo
    cout << "\nQueue Operations:" << endl;
    LinkedQueue queue;
    
    queue.enqueue(100);
    queue.enqueue(200);
    queue.enqueue(300);
    
    queue.display();
    
    cout << "Front element: " << queue.getFront() << endl;
    cout << "Dequeued: " << queue.dequeue() << endl;
    
    queue.display();
    
    // Expression Evaluation Demo
    cout << "\nExpression Evaluation:" << endl;
    string expression = "((())())";
    bool balanced = ExpressionEvaluator::isBalanced(expression);
    cout << expression << " is " << (balanced ? "balanced" : "not balanced") << endl;
    
    string infix = "a+b*c";
    string postfix = ExpressionEvaluator::infixToPostfix(infix);
    cout << "Infix: " << infix << " -> Postfix: " << postfix << endl;
    
    // Priority Queue Demo
    cout << "\nPriority Queue Operations:" << endl;
    HeapPriorityQueue pq;
    
    pq.enqueue(10, 1);
    pq.enqueue(20, 3);
    pq.enqueue(30, 2);
    
    cout << "Highest priority element: " << pq.peek() << endl;
    cout << "Dequeued: " << pq.dequeue() << endl;
    cout << "Next highest: " << pq.peek() << endl;
    
    return 0;
}
```

---

## Summary

In this part, you learned:
- Stack and Queue fundamentals and implementations
- Different implementation approaches (array vs linked list)
- Stack applications: expression evaluation, function calls, backtracking
- Queue variations: circular queue, deque, priority queue
- Common problems and algorithmic techniques
- Real-world applications and use cases

**Key Concepts:**
- **Stack (LIFO)**: Last In, First Out - useful for reversing, parsing, backtracking
- **Queue (FIFO)**: First In, First Out - useful for scheduling, BFS, buffering
- **Priority Queue**: Elements served based on priority, not insertion order
- **Deque**: Double-ended queue allowing insertion/deletion at both ends

**Time Complexities:**
- All basic operations: O(1)
- Search: O(n)
- Priority queue with heap: O(log n) for insert/delete

**Space Complexity:**
- Array implementation: O(n) fixed space
- Linked list implementation: O(n) dynamic space

**Next**: [Part 5: Recursion and Backtracking](./part5-recursion-backtracking.md)

---

## Additional Practice Resources

- LeetCode Stack and Queue problems
- HackerRank Data Structures challenges
- GeeksforGeeks Stack and Queue articles
- Practice implementing all variations from scratch
- Analyze time and space complexity for each implementation