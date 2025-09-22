# Part 3: Linked Lists

## Table of Contents
1. [Introduction to Linked Lists](#introduction)
2. [Singly Linked Lists](#singly-linked-lists)
3. [Doubly Linked Lists](#doubly-linked-lists)
4. [Circular Linked Lists](#circular-linked-lists)
5. [Linked List Operations](#operations)
6. [Common Linked List Problems](#problems)
7. [Advanced Techniques](#advanced-techniques)
8. [Practice Exercises](#practice-exercises)

---

## 1. Introduction to Linked Lists {#introduction}

### What is a Linked List?
A linked list is a linear data structure where elements are stored in nodes, and each node contains data and a pointer/reference to the next node.

### Advantages over Arrays
- **Dynamic size**: Can grow or shrink during runtime
- **Efficient insertion/deletion**: O(1) at any position if you have the node reference
- **Memory efficient**: Only allocates memory as needed

### Disadvantages
- **No random access**: Must traverse from head to reach any element
- **Extra memory**: Requires additional memory for storing pointers
- **Cache performance**: Poor cache locality compared to arrays

### Types of Linked Lists
1. **Singly Linked List**: Each node points to the next node
2. **Doubly Linked List**: Each node has pointers to both next and previous nodes
3. **Circular Linked List**: Last node points back to the first node

---

## 2. Singly Linked Lists {#singly-linked-lists}

### Node Structure

```cpp
#include <iostream>
using namespace std;

// Node structure for singly linked list
struct ListNode {
    int data;
    ListNode* next;
    
    // Constructor
    ListNode(int value) : data(value), next(nullptr) {}
    ListNode(int value, ListNode* nextNode) : data(value), next(nextNode) {}
};

// Alternative class-based approach
class Node {
public:
    int data;
    Node* next;
    
    Node(int value) {
        data = value;
        next = nullptr;
    }
};
```

### Singly Linked List Implementation

```cpp
class SinglyLinkedList {
private:
    ListNode* head;
    int size;
    
public:
    // Constructor
    SinglyLinkedList() {
        head = nullptr;
        size = 0;
    }
    
    // Destructor
    ~SinglyLinkedList() {
        clear();
    }
    
    // Insert at beginning
    void insertAtHead(int value) {
        ListNode* newNode = new ListNode(value);
        newNode->next = head;
        head = newNode;
        size++;
    }
    
    // Insert at end
    void insertAtTail(int value) {
        ListNode* newNode = new ListNode(value);
        
        if (head == nullptr) {
            head = newNode;
        } else {
            ListNode* current = head;
            while (current->next != nullptr) {
                current = current->next;
            }
            current->next = newNode;
        }
        size++;
    }
    
    // Insert at specific position
    void insertAt(int position, int value) {
        if (position < 0 || position > size) {
            cout << "Invalid position" << endl;
            return;
        }
        
        if (position == 0) {
            insertAtHead(value);
            return;
        }
        
        ListNode* newNode = new ListNode(value);
        ListNode* current = head;
        
        for (int i = 0; i < position - 1; i++) {
            current = current->next;
        }
        
        newNode->next = current->next;
        current->next = newNode;
        size++;
    }
    
    // Delete from beginning
    void deleteFromHead() {
        if (head == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        ListNode* temp = head;
        head = head->next;
        delete temp;
        size--;
    }
    
    // Delete from end
    void deleteFromTail() {
        if (head == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        if (head->next == nullptr) {
            delete head;
            head = nullptr;
            size--;
            return;
        }
        
        ListNode* current = head;
        while (current->next->next != nullptr) {
            current = current->next;
        }
        
        delete current->next;
        current->next = nullptr;
        size--;
    }
    
    // Delete specific value
    void deleteValue(int value) {
        if (head == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        if (head->data == value) {
            deleteFromHead();
            return;
        }
        
        ListNode* current = head;
        while (current->next != nullptr && current->next->data != value) {
            current = current->next;
        }
        
        if (current->next != nullptr) {
            ListNode* temp = current->next;
            current->next = current->next->next;
            delete temp;
            size--;
        } else {
            cout << "Value not found" << endl;
        }
    }
    
    // Search for value
    int search(int value) {
        ListNode* current = head;
        int position = 0;
        
        while (current != nullptr) {
            if (current->data == value) {
                return position;
            }
            current = current->next;
            position++;
        }
        
        return -1; // Not found
    }
    
    // Get value at position
    int get(int position) {
        if (position < 0 || position >= size) {
            throw out_of_range("Invalid position");
        }
        
        ListNode* current = head;
        for (int i = 0; i < position; i++) {
            current = current->next;
        }
        
        return current->data;
    }
    
    // Display list
    void display() {
        ListNode* current = head;
        cout << "List: ";
        while (current != nullptr) {
            cout << current->data << " -> ";
            current = current->next;
        }
        cout << "NULL" << endl;
    }
    
    // Get size
    int getSize() {
        return size;
    }
    
    // Check if empty
    bool isEmpty() {
        return head == nullptr;
    }
    
    // Clear all nodes
    void clear() {
        while (head != nullptr) {
            ListNode* temp = head;
            head = head->next;
            delete temp;
        }
        size = 0;
    }
    
    // Reverse the list
    void reverse() {
        ListNode* prev = nullptr;
        ListNode* current = head;
        ListNode* next = nullptr;
        
        while (current != nullptr) {
            next = current->next;
            current->next = prev;
            prev = current;
            current = next;
        }
        
        head = prev;
    }
};
```

---

## 3. Doubly Linked Lists {#doubly-linked-lists}

### Node Structure

```cpp
// Node structure for doubly linked list
struct DoublyListNode {
    int data;
    DoublyListNode* next;
    DoublyListNode* prev;
    
    DoublyListNode(int value) : data(value), next(nullptr), prev(nullptr) {}
};
```

### Doubly Linked List Implementation

```cpp
class DoublyLinkedList {
private:
    DoublyListNode* head;
    DoublyListNode* tail;
    int size;
    
public:
    // Constructor
    DoublyLinkedList() {
        head = nullptr;
        tail = nullptr;
        size = 0;
    }
    
    // Destructor
    ~DoublyLinkedList() {
        clear();
    }
    
    // Insert at beginning
    void insertAtHead(int value) {
        DoublyListNode* newNode = new DoublyListNode(value);
        
        if (head == nullptr) {
            head = tail = newNode;
        } else {
            newNode->next = head;
            head->prev = newNode;
            head = newNode;
        }
        size++;
    }
    
    // Insert at end
    void insertAtTail(int value) {
        DoublyListNode* newNode = new DoublyListNode(value);
        
        if (tail == nullptr) {
            head = tail = newNode;
        } else {
            tail->next = newNode;
            newNode->prev = tail;
            tail = newNode;
        }
        size++;
    }
    
    // Insert at specific position
    void insertAt(int position, int value) {
        if (position < 0 || position > size) {
            cout << "Invalid position" << endl;
            return;
        }
        
        if (position == 0) {
            insertAtHead(value);
            return;
        }
        
        if (position == size) {
            insertAtTail(value);
            return;
        }
        
        DoublyListNode* newNode = new DoublyListNode(value);
        DoublyListNode* current = head;
        
        for (int i = 0; i < position; i++) {
            current = current->next;
        }
        
        newNode->next = current;
        newNode->prev = current->prev;
        current->prev->next = newNode;
        current->prev = newNode;
        size++;
    }
    
    // Delete from beginning
    void deleteFromHead() {
        if (head == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        DoublyListNode* temp = head;
        
        if (head == tail) {
            head = tail = nullptr;
        } else {
            head = head->next;
            head->prev = nullptr;
        }
        
        delete temp;
        size--;
    }
    
    // Delete from end
    void deleteFromTail() {
        if (tail == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        DoublyListNode* temp = tail;
        
        if (head == tail) {
            head = tail = nullptr;
        } else {
            tail = tail->prev;
            tail->next = nullptr;
        }
        
        delete temp;
        size--;
    }
    
    // Display forward
    void displayForward() {
        DoublyListNode* current = head;
        cout << "Forward: ";
        while (current != nullptr) {
            cout << current->data << " <-> ";
            current = current->next;
        }
        cout << "NULL" << endl;
    }
    
    // Display backward
    void displayBackward() {
        DoublyListNode* current = tail;
        cout << "Backward: ";
        while (current != nullptr) {
            cout << current->data << " <-> ";
            current = current->prev;
        }
        cout << "NULL" << endl;
    }
    
    // Clear all nodes
    void clear() {
        while (head != nullptr) {
            DoublyListNode* temp = head;
            head = head->next;
            delete temp;
        }
        tail = nullptr;
        size = 0;
    }
    
    // Get size
    int getSize() {
        return size;
    }
};
```

---

## 4. Circular Linked Lists {#circular-linked-lists}

### Circular Singly Linked List

```cpp
class CircularLinkedList {
private:
    ListNode* tail; // Points to last node
    int size;
    
public:
    CircularLinkedList() {
        tail = nullptr;
        size = 0;
    }
    
    ~CircularLinkedList() {
        clear();
    }
    
    // Insert at beginning
    void insertAtHead(int value) {
        ListNode* newNode = new ListNode(value);
        
        if (tail == nullptr) {
            tail = newNode;
            newNode->next = newNode; // Points to itself
        } else {
            newNode->next = tail->next;
            tail->next = newNode;
        }
        size++;
    }
    
    // Insert at end
    void insertAtTail(int value) {
        ListNode* newNode = new ListNode(value);
        
        if (tail == nullptr) {
            tail = newNode;
            newNode->next = newNode;
        } else {
            newNode->next = tail->next;
            tail->next = newNode;
            tail = newNode;
        }
        size++;
    }
    
    // Delete from beginning
    void deleteFromHead() {
        if (tail == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        ListNode* head = tail->next;
        
        if (head == tail) { // Only one node
            delete tail;
            tail = nullptr;
        } else {
            tail->next = head->next;
            delete head;
        }
        size--;
    }
    
    // Display list
    void display() {
        if (tail == nullptr) {
            cout << "List is empty" << endl;
            return;
        }
        
        ListNode* current = tail->next; // Start from head
        cout << "Circular List: ";
        
        do {
            cout << current->data << " -> ";
            current = current->next;
        } while (current != tail->next);
        
        cout << "(back to start)" << endl;
    }
    
    // Clear all nodes
    void clear() {
        if (tail == nullptr) return;
        
        ListNode* current = tail->next;
        while (current != tail) {
            ListNode* temp = current;
            current = current->next;
            delete temp;
        }
        delete tail;
        tail = nullptr;
        size = 0;
    }
    
    int getSize() {
        return size;
    }
};
```

---

## 5. Linked List Operations {#operations}

### Advanced Operations

```cpp
class LinkedListOperations {
public:
    // Find middle of linked list
    static ListNode* findMiddle(ListNode* head) {
        if (head == nullptr) return nullptr;
        
        ListNode* slow = head;
        ListNode* fast = head;
        
        while (fast != nullptr && fast->next != nullptr) {
            slow = slow->next;
            fast = fast->next->next;
        }
        
        return slow;
    }
    
    // Detect cycle in linked list
    static bool hasCycle(ListNode* head) {
        if (head == nullptr || head->next == nullptr) {
            return false;
        }
        
        ListNode* slow = head;
        ListNode* fast = head;
        
        while (fast != nullptr && fast->next != nullptr) {
            slow = slow->next;
            fast = fast->next->next;
            
            if (slow == fast) {
                return true;
            }
        }
        
        return false;
    }
    
    // Find cycle start node
    static ListNode* detectCycleStart(ListNode* head) {
        if (head == nullptr || head->next == nullptr) {
            return nullptr;
        }
        
        ListNode* slow = head;
        ListNode* fast = head;
        
        // Detect if cycle exists
        while (fast != nullptr && fast->next != nullptr) {
            slow = slow->next;
            fast = fast->next->next;
            
            if (slow == fast) {
                break;
            }
        }
        
        // No cycle found
        if (fast == nullptr || fast->next == nullptr) {
            return nullptr;
        }
        
        // Find start of cycle
        slow = head;
        while (slow != fast) {
            slow = slow->next;
            fast = fast->next;
        }
        
        return slow;
    }
    
    // Merge two sorted linked lists
    static ListNode* mergeSorted(ListNode* l1, ListNode* l2) {
        ListNode dummy(0);
        ListNode* current = &dummy;
        
        while (l1 != nullptr && l2 != nullptr) {
            if (l1->data <= l2->data) {
                current->next = l1;
                l1 = l1->next;
            } else {
                current->next = l2;
                l2 = l2->next;
            }
            current = current->next;
        }
        
        // Attach remaining nodes
        current->next = (l1 != nullptr) ? l1 : l2;
        
        return dummy.next;
    }
    
    // Remove nth node from end
    static ListNode* removeNthFromEnd(ListNode* head, int n) {
        ListNode dummy(0);
        dummy.next = head;
        
        ListNode* first = &dummy;
        ListNode* second = &dummy;
        
        // Move first pointer n+1 steps ahead
        for (int i = 0; i <= n; i++) {
            first = first->next;
        }
        
        // Move both pointers until first reaches end
        while (first != nullptr) {
            first = first->next;
            second = second->next;
        }
        
        // Remove the nth node
        ListNode* nodeToDelete = second->next;
        second->next = second->next->next;
        delete nodeToDelete;
        
        return dummy.next;
    }
    
    // Check if linked list is palindrome
    static bool isPalindrome(ListNode* head) {
        if (head == nullptr || head->next == nullptr) {
            return true;
        }
        
        // Find middle
        ListNode* slow = head;
        ListNode* fast = head;
        
        while (fast->next != nullptr && fast->next->next != nullptr) {
            slow = slow->next;
            fast = fast->next->next;
        }
        
        // Reverse second half
        ListNode* secondHalf = reverseList(slow->next);
        
        // Compare first and second half
        ListNode* firstHalf = head;
        bool result = true;
        
        while (secondHalf != nullptr) {
            if (firstHalf->data != secondHalf->data) {
                result = false;
                break;
            }
            firstHalf = firstHalf->next;
            secondHalf = secondHalf->next;
        }
        
        return result;
    }
    
private:
    static ListNode* reverseList(ListNode* head) {
        ListNode* prev = nullptr;
        ListNode* current = head;
        
        while (current != nullptr) {
            ListNode* next = current->next;
            current->next = prev;
            prev = current;
            current = next;
        }
        
        return prev;
    }
};
```

---

## 6. Common Linked List Problems {#problems}

### Problem Solutions

```cpp
class LinkedListProblems {
public:
    // 1. Reverse linked list iteratively
    static ListNode* reverseIterative(ListNode* head) {
        ListNode* prev = nullptr;
        ListNode* current = head;
        
        while (current != nullptr) {
            ListNode* next = current->next;
            current->next = prev;
            prev = current;
            current = next;
        }
        
        return prev;
    }
    
    // 2. Reverse linked list recursively
    static ListNode* reverseRecursive(ListNode* head) {
        if (head == nullptr || head->next == nullptr) {
            return head;
        }
        
        ListNode* newHead = reverseRecursive(head->next);
        head->next->next = head;
        head->next = nullptr;
        
        return newHead;
    }
    
    // 3. Add two numbers represented as linked lists
    static ListNode* addTwoNumbers(ListNode* l1, ListNode* l2) {
        ListNode dummy(0);
        ListNode* current = &dummy;
        int carry = 0;
        
        while (l1 != nullptr || l2 != nullptr || carry != 0) {
            int sum = carry;
            
            if (l1 != nullptr) {
                sum += l1->data;
                l1 = l1->next;
            }
            
            if (l2 != nullptr) {
                sum += l2->data;
                l2 = l2->next;
            }
            
            carry = sum / 10;
            current->next = new ListNode(sum % 10);
            current = current->next;
        }
        
        return dummy.next;
    }
    
    // 4. Intersection of two linked lists
    static ListNode* getIntersectionNode(ListNode* headA, ListNode* headB) {
        if (headA == nullptr || headB == nullptr) {
            return nullptr;
        }
        
        ListNode* pA = headA;
        ListNode* pB = headB;
        
        while (pA != pB) {
            pA = (pA == nullptr) ? headB : pA->next;
            pB = (pB == nullptr) ? headA : pB->next;
        }
        
        return pA;
    }
    
    // 5. Remove duplicates from sorted list
    static ListNode* deleteDuplicates(ListNode* head) {
        ListNode* current = head;
        
        while (current != nullptr && current->next != nullptr) {
            if (current->data == current->next->data) {
                ListNode* duplicate = current->next;
                current->next = current->next->next;
                delete duplicate;
            } else {
                current = current->next;
            }
        }
        
        return head;
    }
    
    // 6. Rotate list to the right by k places
    static ListNode* rotateRight(ListNode* head, int k) {
        if (head == nullptr || head->next == nullptr || k == 0) {
            return head;
        }
        
        // Find length and make it circular
        ListNode* tail = head;
        int length = 1;
        
        while (tail->next != nullptr) {
            tail = tail->next;
            length++;
        }
        
        tail->next = head; // Make circular
        
        // Find new tail (length - k % length - 1 steps from head)
        k = k % length;
        int stepsToNewTail = length - k;
        
        ListNode* newTail = head;
        for (int i = 1; i < stepsToNewTail; i++) {
            newTail = newTail->next;
        }
        
        ListNode* newHead = newTail->next;
        newTail->next = nullptr; // Break the circle
        
        return newHead;
    }
    
    // 7. Partition list around value x
    static ListNode* partition(ListNode* head, int x) {
        ListNode beforeHead(0);
        ListNode afterHead(0);
        
        ListNode* before = &beforeHead;
        ListNode* after = &afterHead;
        
        while (head != nullptr) {
            if (head->data < x) {
                before->next = head;
                before = before->next;
            } else {
                after->next = head;
                after = after->next;
            }
            head = head->next;
        }
        
        after->next = nullptr;
        before->next = afterHead.next;
        
        return beforeHead.next;
    }
};
```

---

## 7. Advanced Techniques {#advanced-techniques}

### Fast and Slow Pointers (Floyd's Algorithm)

```cpp
class FloydAlgorithm {
public:
    // Find cycle length
    static int findCycleLength(ListNode* head) {
        ListNode* slow = head;
        ListNode* fast = head;
        
        // Detect cycle
        while (fast != nullptr && fast->next != nullptr) {
            slow = slow->next;
            fast = fast->next->next;
            
            if (slow == fast) {
                // Cycle detected, now find length
                int length = 1;
                ListNode* current = slow->next;
                
                while (current != slow) {
                    current = current->next;
                    length++;
                }
                
                return length;
            }
        }
        
        return 0; // No cycle
    }
    
    // Find kth node from end
    static ListNode* findKthFromEnd(ListNode* head, int k) {
        ListNode* first = head;
        ListNode* second = head;
        
        // Move first pointer k steps ahead
        for (int i = 0; i < k; i++) {
            if (first == nullptr) {
                return nullptr; // k is larger than list length
            }
            first = first->next;
        }
        
        // Move both pointers until first reaches end
        while (first != nullptr) {
            first = first->next;
            second = second->next;
        }
        
        return second;
    }
};
```

### Merge Sort for Linked Lists

```cpp
class LinkedListMergeSort {
public:
    static ListNode* mergeSort(ListNode* head) {
        if (head == nullptr || head->next == nullptr) {
            return head;
        }
        
        // Split the list into two halves
        ListNode* middle = getMiddle(head);
        ListNode* nextOfMiddle = middle->next;
        middle->next = nullptr;
        
        // Recursively sort both halves
        ListNode* left = mergeSort(head);
        ListNode* right = mergeSort(nextOfMiddle);
        
        // Merge the sorted halves
        return merge(left, right);
    }
    
private:
    static ListNode* getMiddle(ListNode* head) {
        if (head == nullptr) return head;
        
        ListNode* slow = head;
        ListNode* fast = head->next;
        
        while (fast != nullptr && fast->next != nullptr) {
            slow = slow->next;
            fast = fast->next->next;
        }
        
        return slow;
    }
    
    static ListNode* merge(ListNode* left, ListNode* right) {
        if (left == nullptr) return right;
        if (right == nullptr) return left;
        
        if (left->data <= right->data) {
            left->next = merge(left->next, right);
            return left;
        } else {
            right->next = merge(left, right->next);
            return right;
        }
    }
};
```

---

## 8. Practice Exercises {#practice-exercises}

### Basic Exercises
1. **Implement Stack using Linked List**
2. **Implement Queue using Linked List**
3. **Find Length of Linked List** (iterative and recursive)
4. **Search Element in Linked List**
5. **Delete All Occurrences of a Key**

### Intermediate Exercises
1. **Reverse Nodes in k-Group**
2. **Swap Nodes in Pairs**
3. **Remove Duplicates from Unsorted List**
4. **Flatten a Multilevel Doubly Linked List**
5. **Copy List with Random Pointer**

### Advanced Exercises
1. **LRU Cache Implementation**
2. **Merge k Sorted Lists**
3. **Reverse Alternate k Nodes**
4. **Clone a Linked List with Next and Random Pointer**
5. **Convert Binary Tree to Doubly Linked List**

### Complete Example Program

```cpp
#include <iostream>
using namespace std;

int main() {
    cout << "=== Linked List Demo ===" << endl;
    
    // Singly Linked List Demo
    SinglyLinkedList sll;
    
    cout << "\nSingly Linked List Operations:" << endl;
    sll.insertAtHead(10);
    sll.insertAtHead(20);
    sll.insertAtTail(30);
    sll.insertAt(1, 15);
    
    sll.display();
    cout << "Size: " << sll.getSize() << endl;
    
    cout << "Search 15: " << sll.search(15) << endl;
    
    sll.reverse();
    cout << "After reverse: ";
    sll.display();
    
    // Doubly Linked List Demo
    DoublyLinkedList dll;
    
    cout << "\nDoubly Linked List Operations:" << endl;
    dll.insertAtHead(100);
    dll.insertAtTail(200);
    dll.insertAtHead(50);
    
    dll.displayForward();
    dll.displayBackward();
    
    // Circular Linked List Demo
    CircularLinkedList cll;
    
    cout << "\nCircular Linked List Operations:" << endl;
    cll.insertAtHead(1);
    cll.insertAtTail(2);
    cll.insertAtTail(3);
    
    cll.display();
    
    return 0;
}
```

---

## Summary

In this part, you learned:
- Different types of linked lists and their implementations
- Basic operations: insertion, deletion, search, traversal
- Advanced techniques: cycle detection, finding middle, merging
- Common problems and their solutions
- Time and space complexity analysis

**Key Concepts:**
- **Singly Linked List**: Simple, memory efficient, one-way traversal
- **Doubly Linked List**: Bidirectional traversal, more memory overhead
- **Circular Linked List**: No null pointers, useful for round-robin algorithms
- **Two Pointers**: Powerful technique for many linked list problems
- **Recursion**: Natural fit for linked list operations

**Time Complexities:**
- Access: O(n)
- Search: O(n)
- Insertion: O(1) if position known, O(n) otherwise
- Deletion: O(1) if node reference known, O(n) otherwise

**Next**: [Part 4: Stacks and Queues](./part4-stacks-queues.md)

---

## Additional Practice Resources

- LeetCode Linked List problems
- HackerRank Linked List challenges
- GeeksforGeeks Linked List articles
- Practice implementing all operations from scratch
- Visualize pointer movements while coding