# Part 6: Trees

## Table of Contents
1. [Introduction to Trees](#introduction-trees)
2. [Tree Terminology](#tree-terminology)
3. [Binary Trees](#binary-trees)
4. [Tree Traversals](#tree-traversals)
5. [Binary Search Trees (BST)](#binary-search-trees)
6. [AVL Trees](#avl-trees)
7. [Heap and Priority Queue](#heap-priority-queue)
8. [Tree Applications](#tree-applications)
9. [Advanced Tree Problems](#advanced-problems)
10. [Practice Exercises](#practice-exercises)

---

## 1. Introduction to Trees {#introduction-trees}

### What is a Tree?
A tree is a hierarchical data structure consisting of nodes connected by edges, with no cycles.

### Tree Characteristics
- **Hierarchical Structure**: Parent-child relationships
- **No Cycles**: Acyclic connected graph
- **Single Root**: One node with no parent
- **Connected**: Path exists between any two nodes

### Why Use Trees?
- **Hierarchical Data**: Natural representation of hierarchical relationships
- **Efficient Search**: O(log n) search in balanced trees
- **Dynamic Size**: Can grow and shrink during runtime
- **Sorted Data**: Maintain sorted order efficiently

### Tree vs Other Data Structures

| Operation | Array | Linked List | Tree (Balanced) |
|-----------|-------|-------------|----------------|
| Search | O(n) | O(n) | O(log n) |
| Insert | O(n) | O(1) | O(log n) |
| Delete | O(n) | O(1) | O(log n) |
| Access | O(1) | O(n) | O(log n) |

---

## 2. Tree Terminology {#tree-terminology}

### Basic Terms
- **Node**: Basic unit containing data
- **Root**: Top node with no parent
- **Parent**: Node with children
- **Child**: Node with a parent
- **Leaf**: Node with no children
- **Sibling**: Nodes with same parent
- **Ancestor**: Node on path from root to current node
- **Descendant**: Node in subtree of current node

### Tree Properties
- **Height**: Maximum distance from root to leaf
- **Depth**: Distance from root to specific node
- **Level**: All nodes at same depth
- **Degree**: Number of children of a node
- **Size**: Total number of nodes

### Tree Types
- **Binary Tree**: Each node has at most 2 children
- **Complete Tree**: All levels filled except possibly last
- **Full Tree**: Every node has 0 or maximum children
- **Perfect Tree**: All internal nodes have same degree, all leaves at same level
- **Balanced Tree**: Height difference between subtrees ≤ 1

```cpp
// Basic tree node structure
struct TreeNode {
    int data;
    TreeNode* left;
    TreeNode* right;
    
    TreeNode(int val) : data(val), left(nullptr), right(nullptr) {}
};

// General tree node (multiple children)
struct GeneralTreeNode {
    int data;
    vector<GeneralTreeNode*> children;
    
    GeneralTreeNode(int val) : data(val) {}
};
```

---

## 3. Binary Trees {#binary-trees}

### Binary Tree Implementation

```cpp
class BinaryTree {
public:
    TreeNode* root;
    
    BinaryTree() : root(nullptr) {}
    
    // Insert node (level order)
    void insert(int data) {
        TreeNode* newNode = new TreeNode(data);
        
        if (root == nullptr) {
            root = newNode;
            return;
        }
        
        queue<TreeNode*> q;
        q.push(root);
        
        while (!q.empty()) {
            TreeNode* current = q.front();
            q.pop();
            
            if (current->left == nullptr) {
                current->left = newNode;
                return;
            } else {
                q.push(current->left);
            }
            
            if (current->right == nullptr) {
                current->right = newNode;
                return;
            } else {
                q.push(current->right);
            }
        }
    }
    
    // Search for a value
    bool search(int data) {
        return searchHelper(root, data);
    }
    
    // Calculate height
    int height() {
        return heightHelper(root);
    }
    
    // Count total nodes
    int size() {
        return sizeHelper(root);
    }
    
    // Check if tree is balanced
    bool isBalanced() {
        return isBalancedHelper(root) != -1;
    }
    
    // Find maximum value
    int findMax() {
        if (root == nullptr) {
            throw runtime_error("Tree is empty");
        }
        return findMaxHelper(root);
    }
    
    // Find minimum value
    int findMin() {
        if (root == nullptr) {
            throw runtime_error("Tree is empty");
        }
        return findMinHelper(root);
    }
    
    // Mirror the tree
    void mirror() {
        mirrorHelper(root);
    }
    
    // Check if two trees are identical
    bool isIdentical(BinaryTree& other) {
        return isIdenticalHelper(root, other.root);
    }
    
    // Destructor
    ~BinaryTree() {
        destroyTree(root);
    }
    
private:
    bool searchHelper(TreeNode* node, int data) {
        if (node == nullptr) {
            return false;
        }
        
        if (node->data == data) {
            return true;
        }
        
        return searchHelper(node->left, data) || searchHelper(node->right, data);
    }
    
    int heightHelper(TreeNode* node) {
        if (node == nullptr) {
            return -1;
        }
        
        return 1 + max(heightHelper(node->left), heightHelper(node->right));
    }
    
    int sizeHelper(TreeNode* node) {
        if (node == nullptr) {
            return 0;
        }
        
        return 1 + sizeHelper(node->left) + sizeHelper(node->right);
    }
    
    int isBalancedHelper(TreeNode* node) {
        if (node == nullptr) {
            return 0;
        }
        
        int leftHeight = isBalancedHelper(node->left);
        if (leftHeight == -1) return -1;
        
        int rightHeight = isBalancedHelper(node->right);
        if (rightHeight == -1) return -1;
        
        if (abs(leftHeight - rightHeight) > 1) {
            return -1;
        }
        
        return 1 + max(leftHeight, rightHeight);
    }
    
    int findMaxHelper(TreeNode* node) {
        if (node == nullptr) {
            return INT_MIN;
        }
        
        int maxVal = node->data;
        int leftMax = findMaxHelper(node->left);
        int rightMax = findMaxHelper(node->right);
        
        return max({maxVal, leftMax, rightMax});
    }
    
    int findMinHelper(TreeNode* node) {
        if (node == nullptr) {
            return INT_MAX;
        }
        
        int minVal = node->data;
        int leftMin = findMinHelper(node->left);
        int rightMin = findMinHelper(node->right);
        
        return min({minVal, leftMin, rightMin});
    }
    
    void mirrorHelper(TreeNode* node) {
        if (node == nullptr) {
            return;
        }
        
        swap(node->left, node->right);
        mirrorHelper(node->left);
        mirrorHelper(node->right);
    }
    
    bool isIdenticalHelper(TreeNode* node1, TreeNode* node2) {
        if (node1 == nullptr && node2 == nullptr) {
            return true;
        }
        
        if (node1 == nullptr || node2 == nullptr) {
            return false;
        }
        
        return (node1->data == node2->data) &&
               isIdenticalHelper(node1->left, node2->left) &&
               isIdenticalHelper(node1->right, node2->right);
    }
    
    void destroyTree(TreeNode* node) {
        if (node != nullptr) {
            destroyTree(node->left);
            destroyTree(node->right);
            delete node;
        }
    }
};
```

---

## 4. Tree Traversals {#tree-traversals}

### Depth-First Traversals

```cpp
class TreeTraversal {
public:
    // Inorder Traversal (Left, Root, Right)
    static void inorderRecursive(TreeNode* root) {
        if (root == nullptr) {
            return;
        }
        
        inorderRecursive(root->left);
        cout << root->data << " ";
        inorderRecursive(root->right);
    }
    
    // Preorder Traversal (Root, Left, Right)
    static void preorderRecursive(TreeNode* root) {
        if (root == nullptr) {
            return;
        }
        
        cout << root->data << " ";
        preorderRecursive(root->left);
        preorderRecursive(root->right);
    }
    
    // Postorder Traversal (Left, Right, Root)
    static void postorderRecursive(TreeNode* root) {
        if (root == nullptr) {
            return;
        }
        
        postorderRecursive(root->left);
        postorderRecursive(root->right);
        cout << root->data << " ";
    }
    
    // Iterative Inorder Traversal
    static vector<int> inorderIterative(TreeNode* root) {
        vector<int> result;
        stack<TreeNode*> stk;
        TreeNode* current = root;
        
        while (current != nullptr || !stk.empty()) {
            while (current != nullptr) {
                stk.push(current);
                current = current->left;
            }
            
            current = stk.top();
            stk.pop();
            result.push_back(current->data);
            current = current->right;
        }
        
        return result;
    }
    
    // Iterative Preorder Traversal
    static vector<int> preorderIterative(TreeNode* root) {
        vector<int> result;
        if (root == nullptr) return result;
        
        stack<TreeNode*> stk;
        stk.push(root);
        
        while (!stk.empty()) {
            TreeNode* current = stk.top();
            stk.pop();
            
            result.push_back(current->data);
            
            if (current->right) {
                stk.push(current->right);
            }
            if (current->left) {
                stk.push(current->left);
            }
        }
        
        return result;
    }
    
    // Iterative Postorder Traversal
    static vector<int> postorderIterative(TreeNode* root) {
        vector<int> result;
        if (root == nullptr) return result;
        
        stack<TreeNode*> stk;
        TreeNode* lastVisited = nullptr;
        TreeNode* current = root;
        
        while (current != nullptr || !stk.empty()) {
            if (current != nullptr) {
                stk.push(current);
                current = current->left;
            } else {
                TreeNode* peekNode = stk.top();
                
                if (peekNode->right != nullptr && lastVisited != peekNode->right) {
                    current = peekNode->right;
                } else {
                    result.push_back(peekNode->data);
                    lastVisited = stk.top();
                    stk.pop();
                }
            }
        }
        
        return result;
    }
};
```

### Breadth-First Traversal

```cpp
class BreadthFirstTraversal {
public:
    // Level Order Traversal
    static vector<int> levelOrder(TreeNode* root) {
        vector<int> result;
        if (root == nullptr) return result;
        
        queue<TreeNode*> q;
        q.push(root);
        
        while (!q.empty()) {
            TreeNode* current = q.front();
            q.pop();
            
            result.push_back(current->data);
            
            if (current->left) {
                q.push(current->left);
            }
            if (current->right) {
                q.push(current->right);
            }
        }
        
        return result;
    }
    
    // Level Order with Level Information
    static vector<vector<int>> levelOrderByLevels(TreeNode* root) {
        vector<vector<int>> result;
        if (root == nullptr) return result;
        
        queue<TreeNode*> q;
        q.push(root);
        
        while (!q.empty()) {
            int levelSize = q.size();
            vector<int> currentLevel;
            
            for (int i = 0; i < levelSize; i++) {
                TreeNode* current = q.front();
                q.pop();
                
                currentLevel.push_back(current->data);
                
                if (current->left) {
                    q.push(current->left);
                }
                if (current->right) {
                    q.push(current->right);
                }
            }
            
            result.push_back(currentLevel);
        }
        
        return result;
    }
    
    // Zigzag Level Order Traversal
    static vector<vector<int>> zigzagLevelOrder(TreeNode* root) {
        vector<vector<int>> result;
        if (root == nullptr) return result;
        
        queue<TreeNode*> q;
        q.push(root);
        bool leftToRight = true;
        
        while (!q.empty()) {
            int levelSize = q.size();
            vector<int> currentLevel(levelSize);
            
            for (int i = 0; i < levelSize; i++) {
                TreeNode* current = q.front();
                q.pop();
                
                int index = leftToRight ? i : (levelSize - 1 - i);
                currentLevel[index] = current->data;
                
                if (current->left) {
                    q.push(current->left);
                }
                if (current->right) {
                    q.push(current->right);
                }
            }
            
            leftToRight = !leftToRight;
            result.push_back(currentLevel);
        }
        
        return result;
    }
    
    // Vertical Order Traversal
    static vector<vector<int>> verticalOrder(TreeNode* root) {
        vector<vector<int>> result;
        if (root == nullptr) return result;
        
        map<int, vector<int>> columnMap;
        queue<pair<TreeNode*, int>> q;
        q.push({root, 0});
        
        while (!q.empty()) {
            auto [node, column] = q.front();
            q.pop();
            
            columnMap[column].push_back(node->data);
            
            if (node->left) {
                q.push({node->left, column - 1});
            }
            if (node->right) {
                q.push({node->right, column + 1});
            }
        }
        
        for (auto& [col, nodes] : columnMap) {
            result.push_back(nodes);
        }
        
        return result;
    }
};
```

---

## 5. Binary Search Trees (BST) {#binary-search-trees}

### BST Properties
- Left subtree contains only nodes with values less than root
- Right subtree contains only nodes with values greater than root
- Both subtrees are also BSTs
- No duplicate values (in this implementation)

### BST Implementation

```cpp
class BinarySearchTree {
public:
    TreeNode* root;
    
    BinarySearchTree() : root(nullptr) {}
    
    // Insert a value
    void insert(int data) {
        root = insertHelper(root, data);
    }
    
    // Search for a value
    bool search(int data) {
        return searchHelper(root, data);
    }
    
    // Delete a value
    void remove(int data) {
        root = removeHelper(root, data);
    }
    
    // Find minimum value
    int findMin() {
        if (root == nullptr) {
            throw runtime_error("Tree is empty");
        }
        return findMinHelper(root)->data;
    }
    
    // Find maximum value
    int findMax() {
        if (root == nullptr) {
            throw runtime_error("Tree is empty");
        }
        return findMaxHelper(root)->data;
    }
    
    // Find inorder successor
    int findSuccessor(int data) {
        TreeNode* node = findNode(root, data);
        if (node == nullptr) {
            throw runtime_error("Node not found");
        }
        
        TreeNode* successor = findSuccessorHelper(node);
        if (successor == nullptr) {
            throw runtime_error("No successor found");
        }
        
        return successor->data;
    }
    
    // Find inorder predecessor
    int findPredecessor(int data) {
        TreeNode* node = findNode(root, data);
        if (node == nullptr) {
            throw runtime_error("Node not found");
        }
        
        TreeNode* predecessor = findPredecessorHelper(node);
        if (predecessor == nullptr) {
            throw runtime_error("No predecessor found");
        }
        
        return predecessor->data;
    }
    
    // Validate if tree is BST
    bool isValidBST() {
        return isValidBSTHelper(root, LLONG_MIN, LLONG_MAX);
    }
    
    // Find kth smallest element
    int kthSmallest(int k) {
        int count = 0;
        TreeNode* result = kthSmallestHelper(root, k, count);
        if (result == nullptr) {
            throw runtime_error("k is larger than tree size");
        }
        return result->data;
    }
    
    // Convert BST to sorted array
    vector<int> toSortedArray() {
        vector<int> result;
        inorderTraversal(root, result);
        return result;
    }
    
    // Print tree structure
    void printTree() {
        printTreeHelper(root, 0);
    }
    
private:
    TreeNode* insertHelper(TreeNode* node, int data) {
        if (node == nullptr) {
            return new TreeNode(data);
        }
        
        if (data < node->data) {
            node->left = insertHelper(node->left, data);
        } else if (data > node->data) {
            node->right = insertHelper(node->right, data);
        }
        // If data == node->data, do nothing (no duplicates)
        
        return node;
    }
    
    bool searchHelper(TreeNode* node, int data) {
        if (node == nullptr) {
            return false;
        }
        
        if (data == node->data) {
            return true;
        } else if (data < node->data) {
            return searchHelper(node->left, data);
        } else {
            return searchHelper(node->right, data);
        }
    }
    
    TreeNode* removeHelper(TreeNode* node, int data) {
        if (node == nullptr) {
            return node;
        }
        
        if (data < node->data) {
            node->left = removeHelper(node->left, data);
        } else if (data > node->data) {
            node->right = removeHelper(node->right, data);
        } else {
            // Node to be deleted found
            
            // Case 1: Node with only right child or no child
            if (node->left == nullptr) {
                TreeNode* temp = node->right;
                delete node;
                return temp;
            }
            // Case 2: Node with only left child
            else if (node->right == nullptr) {
                TreeNode* temp = node->left;
                delete node;
                return temp;
            }
            
            // Case 3: Node with two children
            TreeNode* temp = findMinHelper(node->right);
            node->data = temp->data;
            node->right = removeHelper(node->right, temp->data);
        }
        
        return node;
    }
    
    TreeNode* findMinHelper(TreeNode* node) {
        while (node->left != nullptr) {
            node = node->left;
        }
        return node;
    }
    
    TreeNode* findMaxHelper(TreeNode* node) {
        while (node->right != nullptr) {
            node = node->right;
        }
        return node;
    }
    
    TreeNode* findNode(TreeNode* node, int data) {
        if (node == nullptr || node->data == data) {
            return node;
        }
        
        if (data < node->data) {
            return findNode(node->left, data);
        } else {
            return findNode(node->right, data);
        }
    }
    
    TreeNode* findSuccessorHelper(TreeNode* node) {
        if (node->right != nullptr) {
            return findMinHelper(node->right);
        }
        
        // Find the deepest ancestor for which node is in left subtree
        TreeNode* successor = nullptr;
        TreeNode* current = root;
        
        while (current != nullptr) {
            if (node->data < current->data) {
                successor = current;
                current = current->left;
            } else if (node->data > current->data) {
                current = current->right;
            } else {
                break;
            }
        }
        
        return successor;
    }
    
    TreeNode* findPredecessorHelper(TreeNode* node) {
        if (node->left != nullptr) {
            return findMaxHelper(node->left);
        }
        
        // Find the deepest ancestor for which node is in right subtree
        TreeNode* predecessor = nullptr;
        TreeNode* current = root;
        
        while (current != nullptr) {
            if (node->data > current->data) {
                predecessor = current;
                current = current->right;
            } else if (node->data < current->data) {
                current = current->left;
            } else {
                break;
            }
        }
        
        return predecessor;
    }
    
    bool isValidBSTHelper(TreeNode* node, long long minVal, long long maxVal) {
        if (node == nullptr) {
            return true;
        }
        
        if (node->data <= minVal || node->data >= maxVal) {
            return false;
        }
        
        return isValidBSTHelper(node->left, minVal, node->data) &&
               isValidBSTHelper(node->right, node->data, maxVal);
    }
    
    TreeNode* kthSmallestHelper(TreeNode* node, int k, int& count) {
        if (node == nullptr) {
            return nullptr;
        }
        
        TreeNode* left = kthSmallestHelper(node->left, k, count);
        if (left != nullptr) {
            return left;
        }
        
        count++;
        if (count == k) {
            return node;
        }
        
        return kthSmallestHelper(node->right, k, count);
    }
    
    void inorderTraversal(TreeNode* node, vector<int>& result) {
        if (node != nullptr) {
            inorderTraversal(node->left, result);
            result.push_back(node->data);
            inorderTraversal(node->right, result);
        }
    }
    
    void printTreeHelper(TreeNode* node, int space) {
        const int COUNT = 10;
        if (node == nullptr) {
            return;
        }
        
        space += COUNT;
        printTreeHelper(node->right, space);
        
        cout << endl;
        for (int i = COUNT; i < space; i++) {
            cout << " ";
        }
        cout << node->data << endl;
        
        printTreeHelper(node->left, space);
    }
};
```

---

## 6. AVL Trees {#avl-trees}

### AVL Tree Properties
- Self-balancing binary search tree
- Height difference between left and right subtrees ≤ 1
- Guarantees O(log n) operations

### AVL Tree Implementation

```cpp
struct AVLNode {
    int data;
    AVLNode* left;
    AVLNode* right;
    int height;
    
    AVLNode(int val) : data(val), left(nullptr), right(nullptr), height(1) {}
};

class AVLTree {
public:
    AVLNode* root;
    
    AVLTree() : root(nullptr) {}
    
    void insert(int data) {
        root = insertHelper(root, data);
    }
    
    void remove(int data) {
        root = removeHelper(root, data);
    }
    
    bool search(int data) {
        return searchHelper(root, data);
    }
    
    void printInorder() {
        inorderHelper(root);
        cout << endl;
    }
    
private:
    int getHeight(AVLNode* node) {
        return node ? node->height : 0;
    }
    
    int getBalance(AVLNode* node) {
        return node ? getHeight(node->left) - getHeight(node->right) : 0;
    }
    
    void updateHeight(AVLNode* node) {
        if (node) {
            node->height = 1 + max(getHeight(node->left), getHeight(node->right));
        }
    }
    
    AVLNode* rotateRight(AVLNode* y) {
        AVLNode* x = y->left;
        AVLNode* T2 = x->right;
        
        // Perform rotation
        x->right = y;
        y->left = T2;
        
        // Update heights
        updateHeight(y);
        updateHeight(x);
        
        return x;
    }
    
    AVLNode* rotateLeft(AVLNode* x) {
        AVLNode* y = x->right;
        AVLNode* T2 = y->left;
        
        // Perform rotation
        y->left = x;
        x->right = T2;
        
        // Update heights
        updateHeight(x);
        updateHeight(y);
        
        return y;
    }
    
    AVLNode* insertHelper(AVLNode* node, int data) {
        // Step 1: Perform normal BST insertion
        if (node == nullptr) {
            return new AVLNode(data);
        }
        
        if (data < node->data) {
            node->left = insertHelper(node->left, data);
        } else if (data > node->data) {
            node->right = insertHelper(node->right, data);
        } else {
            return node; // No duplicates
        }
        
        // Step 2: Update height
        updateHeight(node);
        
        // Step 3: Get balance factor
        int balance = getBalance(node);
        
        // Step 4: Perform rotations if needed
        
        // Left Left Case
        if (balance > 1 && data < node->left->data) {
            return rotateRight(node);
        }
        
        // Right Right Case
        if (balance < -1 && data > node->right->data) {
            return rotateLeft(node);
        }
        
        // Left Right Case
        if (balance > 1 && data > node->left->data) {
            node->left = rotateLeft(node->left);
            return rotateRight(node);
        }
        
        // Right Left Case
        if (balance < -1 && data < node->right->data) {
            node->right = rotateRight(node->right);
            return rotateLeft(node);
        }
        
        return node;
    }
    
    AVLNode* findMin(AVLNode* node) {
        while (node->left != nullptr) {
            node = node->left;
        }
        return node;
    }
    
    AVLNode* removeHelper(AVLNode* node, int data) {
        // Step 1: Perform normal BST deletion
        if (node == nullptr) {
            return node;
        }
        
        if (data < node->data) {
            node->left = removeHelper(node->left, data);
        } else if (data > node->data) {
            node->right = removeHelper(node->right, data);
        } else {
            if (node->left == nullptr || node->right == nullptr) {
                AVLNode* temp = node->left ? node->left : node->right;
                
                if (temp == nullptr) {
                    temp = node;
                    node = nullptr;
                } else {
                    *node = *temp;
                }
                delete temp;
            } else {
                AVLNode* temp = findMin(node->right);
                node->data = temp->data;
                node->right = removeHelper(node->right, temp->data);
            }
        }
        
        if (node == nullptr) {
            return node;
        }
        
        // Step 2: Update height
        updateHeight(node);
        
        // Step 3: Get balance factor
        int balance = getBalance(node);
        
        // Step 4: Perform rotations if needed
        
        // Left Left Case
        if (balance > 1 && getBalance(node->left) >= 0) {
            return rotateRight(node);
        }
        
        // Left Right Case
        if (balance > 1 && getBalance(node->left) < 0) {
            node->left = rotateLeft(node->left);
            return rotateRight(node);
        }
        
        // Right Right Case
        if (balance < -1 && getBalance(node->right) <= 0) {
            return rotateLeft(node);
        }
        
        // Right Left Case
        if (balance < -1 && getBalance(node->right) > 0) {
            node->right = rotateRight(node->right);
            return rotateLeft(node);
        }
        
        return node;
    }
    
    bool searchHelper(AVLNode* node, int data) {
        if (node == nullptr) {
            return false;
        }
        
        if (data == node->data) {
            return true;
        } else if (data < node->data) {
            return searchHelper(node->left, data);
        } else {
            return searchHelper(node->right, data);
        }
    }
    
    void inorderHelper(AVLNode* node) {
        if (node != nullptr) {
            inorderHelper(node->left);
            cout << node->data << " ";
            inorderHelper(node->right);
        }
    }
};
```

---

## 7. Heap and Priority Queue {#heap-priority-queue}

### Heap Properties
- **Complete Binary Tree**: All levels filled except possibly last
- **Heap Property**: 
  - Max Heap: Parent ≥ children
  - Min Heap: Parent ≤ children

### Binary Heap Implementation

```cpp
class MaxHeap {
public:
    vector<int> heap;
    
    MaxHeap() {}
    
    MaxHeap(const vector<int>& arr) {
        heap = arr;
        buildHeap();
    }
    
    void insert(int data) {
        heap.push_back(data);
        heapifyUp(heap.size() - 1);
    }
    
    int extractMax() {
        if (heap.empty()) {
            throw runtime_error("Heap is empty");
        }
        
        int maxVal = heap[0];
        heap[0] = heap.back();
        heap.pop_back();
        
        if (!heap.empty()) {
            heapifyDown(0);
        }
        
        return maxVal;
    }
    
    int getMax() {
        if (heap.empty()) {
            throw runtime_error("Heap is empty");
        }
        return heap[0];
    }
    
    void increaseKey(int index, int newVal) {
        if (index >= heap.size() || newVal < heap[index]) {
            throw invalid_argument("Invalid operation");
        }
        
        heap[index] = newVal;
        heapifyUp(index);
    }
    
    void deleteKey(int index) {
        if (index >= heap.size()) {
            throw invalid_argument("Index out of bounds");
        }
        
        increaseKey(index, INT_MAX);
        extractMax();
    }
    
    bool isEmpty() {
        return heap.empty();
    }
    
    int size() {
        return heap.size();
    }
    
    void printHeap() {
        for (int val : heap) {
            cout << val << " ";
        }
        cout << endl;
    }
    
    // Heap Sort
    static void heapSort(vector<int>& arr) {
        MaxHeap heap(arr);
        
        for (int i = arr.size() - 1; i >= 0; i--) {
            arr[i] = heap.extractMax();
        }
    }
    
private:
    int parent(int i) { return (i - 1) / 2; }
    int leftChild(int i) { return 2 * i + 1; }
    int rightChild(int i) { return 2 * i + 2; }
    
    void heapifyUp(int index) {
        while (index > 0 && heap[parent(index)] < heap[index]) {
            swap(heap[index], heap[parent(index)]);
            index = parent(index);
        }
    }
    
    void heapifyDown(int index) {
        int largest = index;
        int left = leftChild(index);
        int right = rightChild(index);
        
        if (left < heap.size() && heap[left] > heap[largest]) {
            largest = left;
        }
        
        if (right < heap.size() && heap[right] > heap[largest]) {
            largest = right;
        }
        
        if (largest != index) {
            swap(heap[index], heap[largest]);
            heapifyDown(largest);
        }
    }
    
    void buildHeap() {
        for (int i = heap.size() / 2 - 1; i >= 0; i--) {
            heapifyDown(i);
        }
    }
};

// Min Heap (similar implementation with reversed comparisons)
class MinHeap {
public:
    vector<int> heap;
    
    void insert(int data) {
        heap.push_back(data);
        heapifyUp(heap.size() - 1);
    }
    
    int extractMin() {
        if (heap.empty()) {
            throw runtime_error("Heap is empty");
        }
        
        int minVal = heap[0];
        heap[0] = heap.back();
        heap.pop_back();
        
        if (!heap.empty()) {
            heapifyDown(0);
        }
        
        return minVal;
    }
    
    int getMin() {
        if (heap.empty()) {
            throw runtime_error("Heap is empty");
        }
        return heap[0];
    }
    
private:
    int parent(int i) { return (i - 1) / 2; }
    int leftChild(int i) { return 2 * i + 1; }
    int rightChild(int i) { return 2 * i + 2; }
    
    void heapifyUp(int index) {
        while (index > 0 && heap[parent(index)] > heap[index]) {
            swap(heap[index], heap[parent(index)]);
            index = parent(index);
        }
    }
    
    void heapifyDown(int index) {
        int smallest = index;
        int left = leftChild(index);
        int right = rightChild(index);
        
        if (left < heap.size() && heap[left] < heap[smallest]) {
            smallest = left;
        }
        
        if (right < heap.size() && heap[right] < heap[smallest]) {
            smallest = right;
        }
        
        if (smallest != index) {
            swap(heap[index], heap[smallest]);
            heapifyDown(smallest);
        }
    }
};
```

---

## 8. Tree Applications {#tree-applications}

### Expression Trees

```cpp
class ExpressionTree {
public:
    struct ExprNode {
        string data;
        ExprNode* left;
        ExprNode* right;
        
        ExprNode(string val) : data(val), left(nullptr), right(nullptr) {}
    };
    
    ExprNode* root;
    
    ExpressionTree() : root(nullptr) {}
    
    // Build expression tree from postfix expression
    void buildFromPostfix(const vector<string>& postfix) {
        stack<ExprNode*> stk;
        
        for (const string& token : postfix) {
            ExprNode* node = new ExprNode(token);
            
            if (isOperator(token)) {
                node->right = stk.top(); stk.pop();
                node->left = stk.top(); stk.pop();
            }
            
            stk.push(node);
        }
        
        root = stk.top();
    }
    
    // Evaluate expression tree
    double evaluate() {
        return evaluateHelper(root);
    }
    
    // Convert to infix expression
    string toInfix() {
        return toInfixHelper(root);
    }
    
    // Convert to prefix expression
    string toPrefix() {
        return toPrefixHelper(root);
    }
    
private:
    bool isOperator(const string& token) {
        return token == "+" || token == "-" || token == "*" || token == "/";
    }
    
    double evaluateHelper(ExprNode* node) {
        if (node == nullptr) {
            return 0;
        }
        
        if (!isOperator(node->data)) {
            return stod(node->data);
        }
        
        double left = evaluateHelper(node->left);
        double right = evaluateHelper(node->right);
        
        if (node->data == "+") return left + right;
        if (node->data == "-") return left - right;
        if (node->data == "*") return left * right;
        if (node->data == "/") return left / right;
        
        return 0;
    }
    
    string toInfixHelper(ExprNode* node) {
        if (node == nullptr) {
            return "";
        }
        
        if (!isOperator(node->data)) {
            return node->data;
        }
        
        return "(" + toInfixHelper(node->left) + " " + 
               node->data + " " + toInfixHelper(node->right) + ")";
    }
    
    string toPrefixHelper(ExprNode* node) {
        if (node == nullptr) {
            return "";
        }
        
        if (!isOperator(node->data)) {
            return node->data;
        }
        
        return node->data + " " + toPrefixHelper(node->left) + 
               " " + toPrefixHelper(node->right);
    }
};
```

### Trie (Prefix Tree)

```cpp
class TrieNode {
public:
    unordered_map<char, TrieNode*> children;
    bool isEndOfWord;
    
    TrieNode() : isEndOfWord(false) {}
};

class Trie {
public:
    TrieNode* root;
    
    Trie() {
        root = new TrieNode();
    }
    
    void insert(const string& word) {
        TrieNode* current = root;
        
        for (char ch : word) {
            if (current->children.find(ch) == current->children.end()) {
                current->children[ch] = new TrieNode();
            }
            current = current->children[ch];
        }
        
        current->isEndOfWord = true;
    }
    
    bool search(const string& word) {
        TrieNode* current = root;
        
        for (char ch : word) {
            if (current->children.find(ch) == current->children.end()) {
                return false;
            }
            current = current->children[ch];
        }
        
        return current->isEndOfWord;
    }
    
    bool startsWith(const string& prefix) {
        TrieNode* current = root;
        
        for (char ch : prefix) {
            if (current->children.find(ch) == current->children.end()) {
                return false;
            }
            current = current->children[ch];
        }
        
        return true;
    }
    
    void remove(const string& word) {
        removeHelper(root, word, 0);
    }
    
    vector<string> getAllWords() {
        vector<string> result;
        string current = "";
        getAllWordsHelper(root, current, result);
        return result;
    }
    
    vector<string> getWordsWithPrefix(const string& prefix) {
        vector<string> result;
        TrieNode* prefixNode = root;
        
        // Navigate to prefix node
        for (char ch : prefix) {
            if (prefixNode->children.find(ch) == prefixNode->children.end()) {
                return result; // No words with this prefix
            }
            prefixNode = prefixNode->children[ch];
        }
        
        // Get all words from prefix node
        string current = prefix;
        getAllWordsHelper(prefixNode, current, result);
        return result;
    }
    
private:
    bool removeHelper(TrieNode* node, const string& word, int index) {
        if (index == word.length()) {
            if (!node->isEndOfWord) {
                return false; // Word doesn't exist
            }
            
            node->isEndOfWord = false;
            return node->children.empty();
        }
        
        char ch = word[index];
        if (node->children.find(ch) == node->children.end()) {
            return false; // Word doesn't exist
        }
        
        bool shouldDeleteChild = removeHelper(node->children[ch], word, index + 1);
        
        if (shouldDeleteChild) {
            delete node->children[ch];
            node->children.erase(ch);
            return !node->isEndOfWord && node->children.empty();
        }
        
        return false;
    }
    
    void getAllWordsHelper(TrieNode* node, string& current, vector<string>& result) {
        if (node->isEndOfWord) {
            result.push_back(current);
        }
        
        for (auto& [ch, child] : node->children) {
            current.push_back(ch);
            getAllWordsHelper(child, current, result);
            current.pop_back();
        }
    }
};
```

---

## 9. Advanced Tree Problems {#advanced-problems}

### Tree Construction Problems

```cpp
class TreeConstruction {
public:
    // Build tree from inorder and preorder traversals
    static TreeNode* buildTreeFromInorderPreorder(vector<int>& preorder, vector<int>& inorder) {
        unordered_map<int, int> inorderMap;
        for (int i = 0; i < inorder.size(); i++) {
            inorderMap[inorder[i]] = i;
        }
        
        int preorderIndex = 0;
        return buildHelper(preorder, inorderMap, preorderIndex, 0, inorder.size() - 1);
    }
    
    // Build tree from inorder and postorder traversals
    static TreeNode* buildTreeFromInorderPostorder(vector<int>& inorder, vector<int>& postorder) {
        unordered_map<int, int> inorderMap;
        for (int i = 0; i < inorder.size(); i++) {
            inorderMap[inorder[i]] = i;
        }
        
        int postorderIndex = postorder.size() - 1;
        return buildHelperPost(postorder, inorderMap, postorderIndex, 0, inorder.size() - 1);
    }
    
    // Serialize tree to string
    static string serialize(TreeNode* root) {
        if (root == nullptr) {
            return "null,";
        }
        
        return to_string(root->data) + "," + serialize(root->left) + serialize(root->right);
    }
    
    // Deserialize string to tree
    static TreeNode* deserialize(string data) {
        queue<string> nodes;
        stringstream ss(data);
        string item;
        
        while (getline(ss, item, ',')) {
            nodes.push(item);
        }
        
        return deserializeHelper(nodes);
    }
    
private:
    static TreeNode* buildHelper(vector<int>& preorder, unordered_map<int, int>& inorderMap,
                               int& preorderIndex, int inorderStart, int inorderEnd) {
        if (inorderStart > inorderEnd) {
            return nullptr;
        }
        
        int rootVal = preorder[preorderIndex++];
        TreeNode* root = new TreeNode(rootVal);
        
        int inorderIndex = inorderMap[rootVal];
        
        root->left = buildHelper(preorder, inorderMap, preorderIndex, inorderStart, inorderIndex - 1);
        root->right = buildHelper(preorder, inorderMap, preorderIndex, inorderIndex + 1, inorderEnd);
        
        return root;
    }
    
    static TreeNode* buildHelperPost(vector<int>& postorder, unordered_map<int, int>& inorderMap,
                                   int& postorderIndex, int inorderStart, int inorderEnd) {
        if (inorderStart > inorderEnd) {
            return nullptr;
        }
        
        int rootVal = postorder[postorderIndex--];
        TreeNode* root = new TreeNode(rootVal);
        
        int inorderIndex = inorderMap[rootVal];
        
        root->right = buildHelperPost(postorder, inorderMap, postorderIndex, inorderIndex + 1, inorderEnd);
        root->left = buildHelperPost(postorder, inorderMap, postorderIndex, inorderStart, inorderIndex - 1);
        
        return root;
    }
    
    static TreeNode* deserializeHelper(queue<string>& nodes) {
        string val = nodes.front();
        nodes.pop();
        
        if (val == "null") {
            return nullptr;
        }
        
        TreeNode* root = new TreeNode(stoi(val));
        root->left = deserializeHelper(nodes);
        root->right = deserializeHelper(nodes);
        
        return root;
    }
};
```

### Tree Path Problems

```cpp
class TreePathProblems {
public:
    // Find all root-to-leaf paths
    static vector<vector<int>> allRootToLeafPaths(TreeNode* root) {
        vector<vector<int>> result;
        vector<int> currentPath;
        findPaths(root, currentPath, result);
        return result;
    }
    
    // Path sum - check if path with given sum exists
    static bool hasPathSum(TreeNode* root, int targetSum) {
        if (root == nullptr) {
            return false;
        }
        
        if (root->left == nullptr && root->right == nullptr) {
            return root->data == targetSum;
        }
        
        return hasPathSum(root->left, targetSum - root->data) ||
               hasPathSum(root->right, targetSum - root->data);
    }
    
    // Find all paths with given sum
    static vector<vector<int>> pathSum(TreeNode* root, int targetSum) {
        vector<vector<int>> result;
        vector<int> currentPath;
        findPathsWithSum(root, targetSum, currentPath, result);
        return result;
    }
    
    // Maximum path sum (any node to any node)
    static int maxPathSum(TreeNode* root) {
        int maxSum = INT_MIN;
        maxPathSumHelper(root, maxSum);
        return maxSum;
    }
    
    // Diameter of tree (longest path between any two nodes)
    static int diameter(TreeNode* root) {
        int maxDiameter = 0;
        diameterHelper(root, maxDiameter);
        return maxDiameter;
    }
    
    // Lowest Common Ancestor
    static TreeNode* lowestCommonAncestor(TreeNode* root, TreeNode* p, TreeNode* q) {
        if (root == nullptr || root == p || root == q) {
            return root;
        }
        
        TreeNode* left = lowestCommonAncestor(root->left, p, q);
        TreeNode* right = lowestCommonAncestor(root->right, p, q);
        
        if (left != nullptr && right != nullptr) {
            return root;
        }
        
        return left != nullptr ? left : right;
    }
    
private:
    static void findPaths(TreeNode* node, vector<int>& currentPath, vector<vector<int>>& result) {
        if (node == nullptr) {
            return;
        }
        
        currentPath.push_back(node->data);
        
        if (node->left == nullptr && node->right == nullptr) {
            result.push_back(currentPath);
        } else {
            findPaths(node->left, currentPath, result);
            findPaths(node->right, currentPath, result);
        }
        
        currentPath.pop_back();
    }
    
    static void findPathsWithSum(TreeNode* node, int targetSum, vector<int>& currentPath, 
                               vector<vector<int>>& result) {
        if (node == nullptr) {
            return;
        }
        
        currentPath.push_back(node->data);
        
        if (node->left == nullptr && node->right == nullptr && node->data == targetSum) {
            result.push_back(currentPath);
        } else {
            findPathsWithSum(node->left, targetSum - node->data, currentPath, result);
            findPathsWithSum(node->right, targetSum - node->data, currentPath, result);
        }
        
        currentPath.pop_back();
    }
    
    static int maxPathSumHelper(TreeNode* node, int& maxSum) {
        if (node == nullptr) {
            return 0;
        }
        
        int leftSum = max(0, maxPathSumHelper(node->left, maxSum));
        int rightSum = max(0, maxPathSumHelper(node->right, maxSum));
        
        int currentMax = node->data + leftSum + rightSum;
        maxSum = max(maxSum, currentMax);
        
        return node->data + max(leftSum, rightSum);
    }
    
    static int diameterHelper(TreeNode* node, int& maxDiameter) {
        if (node == nullptr) {
            return 0;
        }
        
        int leftHeight = diameterHelper(node->left, maxDiameter);
        int rightHeight = diameterHelper(node->right, maxDiameter);
        
        maxDiameter = max(maxDiameter, leftHeight + rightHeight);
        
        return 1 + max(leftHeight, rightHeight);
    }
};
```

---

## 10. Practice Exercises {#practice-exercises}

### Basic Tree Exercises
1. **Tree Basics**:
   - Count total nodes in tree
   - Count leaf nodes
   - Find height/depth of tree
   - Check if tree is complete/full/perfect

2. **Tree Traversals**:
   - Implement all traversals (recursive and iterative)
   - Level order traversal variations
   - Boundary traversal
   - Diagonal traversal

3. **Binary Search Tree**:
   - Insert, delete, search operations
   - Find min/max elements
   - Find inorder successor/predecessor
   - Validate BST

### Intermediate Tree Exercises
1. **Tree Construction**:
   - Build tree from traversals
   - Serialize and deserialize tree
   - Clone a tree
   - Convert tree to its mirror

2. **Tree Properties**:
   - Check if trees are identical
   - Check if tree is balanced
   - Find diameter of tree
   - Calculate tree width

3. **Path Problems**:
   - Root to leaf paths
   - Path with given sum
   - Maximum path sum
   - Lowest common ancestor

### Advanced Tree Exercises
1. **Tree Modifications**:
   - Convert BST to sorted doubly linked list
   - Flatten tree to linked list
   - Connect nodes at same level
   - Populate next right pointers

2. **Special Trees**:
   - Implement AVL tree with rotations
   - Red-Black tree operations
   - Segment tree for range queries
   - Fenwick tree (Binary Indexed Tree)

3. **Tree Applications**:
   - Expression tree evaluation
   - Trie for string operations
   - Huffman coding tree
   - Decision tree implementation

### Complete Example Program

```cpp
#include <iostream>
#include <vector>
#include <queue>
#include <stack>
#include <unordered_map>
using namespace std;

int main() {
    cout << "=== Trees Demo ===" << endl;
    
    // Binary Tree operations
    cout << "\nBinary Tree:" << endl;
    BinaryTree bt;
    bt.insert(1);
    bt.insert(2);
    bt.insert(3);
    bt.insert(4);
    bt.insert(5);
    
    cout << "Tree height: " << bt.height() << endl;
    cout << "Tree size: " << bt.size() << endl;
    cout << "Is balanced: " << (bt.isBalanced() ? "Yes" : "No") << endl;
    
    // Tree traversals
    cout << "\nTraversals:" << endl;
    cout << "Inorder: ";
    TreeTraversal::inorderRecursive(bt.root);
    cout << endl;
    
    cout << "Preorder: ";
    TreeTraversal::preorderRecursive(bt.root);
    cout << endl;
    
    cout << "Level order: ";
    vector<int> levelOrder = BreadthFirstTraversal::levelOrder(bt.root);
    for (int val : levelOrder) {
        cout << val << " ";
    }
    cout << endl;
    
    // Binary Search Tree
    cout << "\nBinary Search Tree:" << endl;
    BinarySearchTree bst;
    vector<int> values = {50, 30, 70, 20, 40, 60, 80};
    
    for (int val : values) {
        bst.insert(val);
    }
    
    cout << "BST Inorder: ";
    TreeTraversal::inorderRecursive(bst.root);
    cout << endl;
    
    cout << "Search 40: " << (bst.search(40) ? "Found" : "Not found") << endl;
    cout << "Min value: " << bst.findMin() << endl;
    cout << "Max value: " << bst.findMax() << endl;
    
    // Heap operations
    cout << "\nMax Heap:" << endl;
    MaxHeap heap;
    vector<int> heapValues = {10, 20, 15, 30, 40};
    
    for (int val : heapValues) {
        heap.insert(val);
    }
    
    cout << "Heap: ";
    heap.printHeap();
    cout << "Max element: " << heap.getMax() << endl;
    
    cout << "Extract max: " << heap.extractMax() << endl;
    cout << "Heap after extraction: ";
    heap.printHeap();
    
    // AVL Tree
    cout << "\nAVL Tree:" << endl;
    AVLTree avl;
    vector<int> avlValues = {10, 20, 30, 40, 50, 25};
    
    for (int val : avlValues) {
        avl.insert(val);
    }
    
    cout << "AVL Inorder: ";
    avl.printInorder();
    
    // Trie operations
    cout << "\nTrie:" << endl;
    Trie trie;
    vector<string> words = {"cat", "car", "card", "care", "careful", "cars", "carry"};
    
    for (const string& word : words) {
        trie.insert(word);
    }
    
    cout << "Search 'car': " << (trie.search("car") ? "Found" : "Not found") << endl;
    cout << "Search 'care': " << (trie.search("care") ? "Found" : "Not found") << endl;
    cout << "Search 'careful': " << (trie.search("careful") ? "Found" : "Not found") << endl;
    cout << "Search 'carefully': " << (trie.search("carefully") ? "Found" : "Not found") << endl;
    
    cout << "Words with prefix 'car': ";
    vector<string> carWords = trie.getWordsWithPrefix("car");
    for (const string& word : carWords) {
        cout << word << " ";
    }
    cout << endl;
    
    return 0;
}
```

---

## Summary

In this part, you learned:
- Tree fundamentals and terminology
- Binary tree implementation and operations
- Tree traversal algorithms (DFS and BFS)
- Binary Search Trees with efficient operations
- Self-balancing trees (AVL trees)
- Heap data structure and priority queues
- Advanced tree applications (Trie, Expression trees)
- Complex tree problems and algorithms

**Key Concepts:**
- **Tree Structure**: Hierarchical data organization
- **Traversals**: Different ways to visit tree nodes
- **BST Property**: Left < Root < Right for efficient searching
- **Balancing**: Maintaining tree height for optimal performance
- **Heap Property**: Parent-child relationship for priority operations
- **Trie**: Efficient string storage and retrieval

**Time Complexities:**

| Operation | Binary Tree | BST (Balanced) | BST (Skewed) | AVL Tree | Heap |
|-----------|-------------|----------------|--------------|----------|---------|
| Search | O(n) | O(log n) | O(n) | O(log n) | O(n) |
| Insert | O(n) | O(log n) | O(n) | O(log n) | O(log n) |
| Delete | O(n) | O(log n) | O(n) | O(log n) | O(log n) |
| Min/Max | O(n) | O(log n) | O(n) | O(log n) | O(1) |

**Space Complexity:**
- **Storage**: O(n) for all tree types
- **Recursion**: O(h) where h is height
- **Balanced trees**: O(log n) space for operations
- **Skewed trees**: O(n) space for operations

**When to Use Different Trees:**
- **Binary Tree**: General hierarchical data
- **BST**: Sorted data with frequent search operations
- **AVL Tree**: Guaranteed balanced operations
- **Heap**: Priority queue operations
- **Trie**: String prefix operations and autocomplete

**Next**: [Part 7: Graphs](./part7-graphs.md)

---

## Additional Practice Resources

- LeetCode Tree problems (Easy to Hard)
- HackerRank Tree challenges
- GeeksforGeeks Tree articles and problems
- Practice tree construction from different traversals
- Implement various tree algorithms from scratch
- Analyze time and space complexity for tree operations
- Study real-world applications of different tree types