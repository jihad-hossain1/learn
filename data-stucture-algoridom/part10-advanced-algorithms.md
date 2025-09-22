# Part 10: Advanced Algorithms

## Table of Contents
1. [Introduction](#introduction)
2. [Greedy Algorithms](#greedy-algorithms)
3. [Divide and Conquer](#divide-and-conquer)
4. [Advanced Graph Algorithms](#advanced-graph-algorithms)
5. [String Algorithms](#string-algorithms)
6. [Number Theory Algorithms](#number-theory)
7. [Computational Geometry](#computational-geometry)
8. [Advanced Data Structures](#advanced-data-structures)
9. [Optimization Techniques](#optimization-techniques)
10. [Practice Problems](#practice-problems)
11. [Complete Example Program](#example-program)
12. [Summary](#summary)

---

## Introduction {#introduction}

Advanced algorithms represent sophisticated problem-solving techniques that go beyond basic data structures and algorithms. These algorithms often combine multiple concepts and require deep understanding of algorithmic paradigms.

### Key Algorithmic Paradigms
1. **Greedy Algorithms**: Make locally optimal choices
2. **Divide and Conquer**: Break problems into smaller subproblems
3. **Dynamic Programming**: Solve overlapping subproblems (covered in Part 9)
4. **Backtracking**: Explore solution space systematically
5. **Branch and Bound**: Optimization with pruning
6. **Randomized Algorithms**: Use randomness for efficiency

### Problem Categories
- **Optimization Problems**: Find best solution
- **Decision Problems**: Yes/No answers
- **Search Problems**: Find specific solutions
- **Counting Problems**: Count number of solutions
- **Construction Problems**: Build solutions

---

## Greedy Algorithms {#greedy-algorithms}

### Characteristics
- Make locally optimal choice at each step
- Hope that local optimum leads to global optimum
- Usually efficient but not always correct
- Require proof of correctness

### 1. Activity Selection Problem
```cpp
class ActivitySelection {
public:
    struct Activity {
        int start, finish, index;
        Activity(int s, int f, int i) : start(s), finish(f), index(i) {}
    };
    
    // Select maximum number of non-overlapping activities
    static vector<int> selectActivities(vector<pair<int, int>>& activities) {
        vector<Activity> acts;
        for (int i = 0; i < activities.size(); i++) {
            acts.push_back(Activity(activities[i].first, activities[i].second, i));
        }
        
        // Sort by finish time
        sort(acts.begin(), acts.end(), [](const Activity& a, const Activity& b) {
            return a.finish < b.finish;
        });
        
        vector<int> selected;
        selected.push_back(acts[0].index);
        int lastFinish = acts[0].finish;
        
        for (int i = 1; i < acts.size(); i++) {
            if (acts[i].start >= lastFinish) {
                selected.push_back(acts[i].index);
                lastFinish = acts[i].finish;
            }
        }
        
        return selected;
    }
    
    // Weighted activity selection (DP approach)
    static int maxWeightActivities(vector<tuple<int, int, int>>& activities) {
        // Sort by finish time
        sort(activities.begin(), activities.end(), [](const auto& a, const auto& b) {
            return get<1>(a) < get<1>(b);
        });
        
        int n = activities.size();
        vector<int> dp(n);
        dp[0] = get<2>(activities[0]); // weight of first activity
        
        for (int i = 1; i < n; i++) {
            int currentWeight = get<2>(activities[i]);
            int latestNonConflicting = findLatestNonConflicting(activities, i);
            
            int includeWeight = currentWeight;
            if (latestNonConflicting != -1) {
                includeWeight += dp[latestNonConflicting];
            }
            
            dp[i] = max(dp[i-1], includeWeight);
        }
        
        return dp[n-1];
    }
    
private:
    static int findLatestNonConflicting(vector<tuple<int, int, int>>& activities, int index) {
        for (int i = index - 1; i >= 0; i--) {
            if (get<1>(activities[i]) <= get<0>(activities[index])) {
                return i;
            }
        }
        return -1;
    }
};
```

### 2. Huffman Coding
```cpp
class HuffmanCoding {
public:
    struct Node {
        char ch;
        int freq;
        Node* left;
        Node* right;
        
        Node(char c, int f) : ch(c), freq(f), left(nullptr), right(nullptr) {}
        Node(int f) : ch('\0'), freq(f), left(nullptr), right(nullptr) {}
    };
    
    struct Compare {
        bool operator()(Node* a, Node* b) {
            return a->freq > b->freq;
        }
    };
    
    static unordered_map<char, string> buildHuffmanCodes(unordered_map<char, int>& frequencies) {
        priority_queue<Node*, vector<Node*>, Compare> pq;
        
        // Create leaf nodes
        for (auto& pair : frequencies) {
            pq.push(new Node(pair.first, pair.second));
        }
        
        // Build Huffman tree
        while (pq.size() > 1) {
            Node* right = pq.top(); pq.pop();
            Node* left = pq.top(); pq.pop();
            
            Node* merged = new Node(left->freq + right->freq);
            merged->left = left;
            merged->right = right;
            
            pq.push(merged);
        }
        
        Node* root = pq.top();
        unordered_map<char, string> codes;
        generateCodes(root, "", codes);
        
        return codes;
    }
    
    static string encode(string text, unordered_map<char, string>& codes) {
        string encoded = "";
        for (char c : text) {
            encoded += codes[c];
        }
        return encoded;
    }
    
    static string decode(string encoded, Node* root) {
        string decoded = "";
        Node* current = root;
        
        for (char bit : encoded) {
            if (bit == '0') {
                current = current->left;
            } else {
                current = current->right;
            }
            
            if (!current->left && !current->right) {
                decoded += current->ch;
                current = root;
            }
        }
        
        return decoded;
    }
    
private:
    static void generateCodes(Node* node, string code, unordered_map<char, string>& codes) {
        if (!node) return;
        
        if (!node->left && !node->right) {
            codes[node->ch] = code.empty() ? "0" : code;
            return;
        }
        
        generateCodes(node->left, code + "0", codes);
        generateCodes(node->right, code + "1", codes);
    }
};
```

### 3. Fractional Knapsack
```cpp
class FractionalKnapsack {
public:
    struct Item {
        int value, weight;
        double ratio;
        int index;
        
        Item(int v, int w, int i) : value(v), weight(w), index(i) {
            ratio = (double)value / weight;
        }
    };
    
    static double fractionalKnapsack(vector<int>& values, vector<int>& weights, int capacity) {
        vector<Item> items;
        for (int i = 0; i < values.size(); i++) {
            items.push_back(Item(values[i], weights[i], i));
        }
        
        // Sort by value-to-weight ratio in descending order
        sort(items.begin(), items.end(), [](const Item& a, const Item& b) {
            return a.ratio > b.ratio;
        });
        
        double totalValue = 0.0;
        int remainingCapacity = capacity;
        
        for (const Item& item : items) {
            if (remainingCapacity >= item.weight) {
                // Take the whole item
                totalValue += item.value;
                remainingCapacity -= item.weight;
            } else {
                // Take fraction of the item
                totalValue += item.value * ((double)remainingCapacity / item.weight);
                break;
            }
        }
        
        return totalValue;
    }
    
    // Get the actual items taken
    static vector<pair<int, double>> getKnapsackItems(vector<int>& values, vector<int>& weights, int capacity) {
        vector<Item> items;
        for (int i = 0; i < values.size(); i++) {
            items.push_back(Item(values[i], weights[i], i));
        }
        
        sort(items.begin(), items.end(), [](const Item& a, const Item& b) {
            return a.ratio > b.ratio;
        });
        
        vector<pair<int, double>> result; // {item_index, fraction_taken}
        int remainingCapacity = capacity;
        
        for (const Item& item : items) {
            if (remainingCapacity >= item.weight) {
                result.push_back({item.index, 1.0});
                remainingCapacity -= item.weight;
            } else {
                double fraction = (double)remainingCapacity / item.weight;
                result.push_back({item.index, fraction});
                break;
            }
        }
        
        return result;
    }
};
```

### 4. Job Scheduling
```cpp
class JobScheduling {
public:
    struct Job {
        int id, deadline, profit;
        Job(int i, int d, int p) : id(i), deadline(d), profit(p) {}
    };
    
    // Job scheduling with deadlines
    static pair<int, int> jobScheduling(vector<Job>& jobs) {
        // Sort jobs by profit in descending order
        sort(jobs.begin(), jobs.end(), [](const Job& a, const Job& b) {
            return a.profit > b.profit;
        });
        
        int maxDeadline = 0;
        for (const Job& job : jobs) {
            maxDeadline = max(maxDeadline, job.deadline);
        }
        
        vector<int> schedule(maxDeadline + 1, -1);
        int totalProfit = 0;
        int jobsScheduled = 0;
        
        for (const Job& job : jobs) {
            // Find a free slot for this job (starting from its deadline)
            for (int slot = job.deadline; slot >= 1; slot--) {
                if (schedule[slot] == -1) {
                    schedule[slot] = job.id;
                    totalProfit += job.profit;
                    jobsScheduled++;
                    break;
                }
            }
        }
        
        return {jobsScheduled, totalProfit};
    }
    
    // Weighted job scheduling (DP approach)
    static int weightedJobScheduling(vector<tuple<int, int, int>>& jobs) {
        // Sort by finish time
        sort(jobs.begin(), jobs.end(), [](const auto& a, const auto& b) {
            return get<1>(a) < get<1>(b);
        });
        
        int n = jobs.size();
        vector<int> dp(n);
        dp[0] = get<2>(jobs[0]);
        
        for (int i = 1; i < n; i++) {
            int currentProfit = get<2>(jobs[i]);
            int latestNonConflicting = findLatestNonConflictingJob(jobs, i);
            
            int includeProfit = currentProfit;
            if (latestNonConflicting != -1) {
                includeProfit += dp[latestNonConflicting];
            }
            
            dp[i] = max(dp[i-1], includeProfit);
        }
        
        return dp[n-1];
    }
    
private:
    static int findLatestNonConflictingJob(vector<tuple<int, int, int>>& jobs, int index) {
        for (int i = index - 1; i >= 0; i--) {
            if (get<1>(jobs[i]) <= get<0>(jobs[index])) {
                return i;
            }
        }
        return -1;
    }
};
```

---

## Divide and Conquer {#divide-and-conquer}

### Characteristics
- Break problem into smaller subproblems
- Solve subproblems recursively
- Combine solutions to solve original problem
- Often leads to efficient algorithms

### 1. Maximum Subarray (Divide and Conquer)
```cpp
class MaxSubarrayDC {
public:
    static int maxSubarraySum(vector<int>& nums) {
        return maxSubarrayHelper(nums, 0, nums.size() - 1);
    }
    
    static tuple<int, int, int> maxSubarrayWithIndices(vector<int>& nums) {
        return maxSubarrayHelperWithIndices(nums, 0, nums.size() - 1);
    }
    
private:
    static int maxSubarrayHelper(vector<int>& nums, int left, int right) {
        if (left == right) {
            return nums[left];
        }
        
        int mid = left + (right - left) / 2;
        
        int leftMax = maxSubarrayHelper(nums, left, mid);
        int rightMax = maxSubarrayHelper(nums, mid + 1, right);
        int crossMax = maxCrossingSubarray(nums, left, mid, right);
        
        return max({leftMax, rightMax, crossMax});
    }
    
    static int maxCrossingSubarray(vector<int>& nums, int left, int mid, int right) {
        int leftSum = INT_MIN;
        int sum = 0;
        for (int i = mid; i >= left; i--) {
            sum += nums[i];
            leftSum = max(leftSum, sum);
        }
        
        int rightSum = INT_MIN;
        sum = 0;
        for (int i = mid + 1; i <= right; i++) {
            sum += nums[i];
            rightSum = max(rightSum, sum);
        }
        
        return leftSum + rightSum;
    }
    
    static tuple<int, int, int> maxSubarrayHelperWithIndices(vector<int>& nums, int left, int right) {
        if (left == right) {
            return {nums[left], left, right};
        }
        
        int mid = left + (right - left) / 2;
        
        auto [leftSum, leftStart, leftEnd] = maxSubarrayHelperWithIndices(nums, left, mid);
        auto [rightSum, rightStart, rightEnd] = maxSubarrayHelperWithIndices(nums, mid + 1, right);
        auto [crossSum, crossStart, crossEnd] = maxCrossingSubarrayWithIndices(nums, left, mid, right);
        
        if (leftSum >= rightSum && leftSum >= crossSum) {
            return {leftSum, leftStart, leftEnd};
        } else if (rightSum >= leftSum && rightSum >= crossSum) {
            return {rightSum, rightStart, rightEnd};
        } else {
            return {crossSum, crossStart, crossEnd};
        }
    }
    
    static tuple<int, int, int> maxCrossingSubarrayWithIndices(vector<int>& nums, int left, int mid, int right) {
        int leftSum = INT_MIN;
        int sum = 0;
        int maxLeft = mid;
        for (int i = mid; i >= left; i--) {
            sum += nums[i];
            if (sum > leftSum) {
                leftSum = sum;
                maxLeft = i;
            }
        }
        
        int rightSum = INT_MIN;
        sum = 0;
        int maxRight = mid + 1;
        for (int i = mid + 1; i <= right; i++) {
            sum += nums[i];
            if (sum > rightSum) {
                rightSum = sum;
                maxRight = i;
            }
        }
        
        return {leftSum + rightSum, maxLeft, maxRight};
    }
};
```

### 2. Closest Pair of Points
```cpp
class ClosestPair {
public:
    struct Point {
        double x, y;
        Point(double x, double y) : x(x), y(y) {}
    };
    
    static double closestPairDistance(vector<Point>& points) {
        sort(points.begin(), points.end(), [](const Point& a, const Point& b) {
            return a.x < b.x;
        });
        
        return closestPairRec(points, 0, points.size() - 1);
    }
    
    static pair<Point, Point> closestPairPoints(vector<Point>& points) {
        sort(points.begin(), points.end(), [](const Point& a, const Point& b) {
            return a.x < b.x;
        });
        
        pair<Point, Point> result;
        double minDist = closestPairRecWithPoints(points, 0, points.size() - 1, result);
        return result;
    }
    
private:
    static double distance(const Point& p1, const Point& p2) {
        return sqrt((p1.x - p2.x) * (p1.x - p2.x) + (p1.y - p2.y) * (p1.y - p2.y));
    }
    
    static double bruteForce(vector<Point>& points, int left, int right) {
        double minDist = DBL_MAX;
        for (int i = left; i <= right; i++) {
            for (int j = i + 1; j <= right; j++) {
                minDist = min(minDist, distance(points[i], points[j]));
            }
        }
        return minDist;
    }
    
    static double closestPairRec(vector<Point>& points, int left, int right) {
        if (right - left <= 3) {
            return bruteForce(points, left, right);
        }
        
        int mid = left + (right - left) / 2;
        Point midPoint = points[mid];
        
        double leftMin = closestPairRec(points, left, mid);
        double rightMin = closestPairRec(points, mid + 1, right);
        
        double minDist = min(leftMin, rightMin);
        
        // Create array of points close to the line dividing the two halves
        vector<Point> strip;
        for (int i = left; i <= right; i++) {
            if (abs(points[i].x - midPoint.x) < minDist) {
                strip.push_back(points[i]);
            }
        }
        
        // Sort strip by y coordinate
        sort(strip.begin(), strip.end(), [](const Point& a, const Point& b) {
            return a.y < b.y;
        });
        
        // Find closest points in strip
        for (int i = 0; i < strip.size(); i++) {
            for (int j = i + 1; j < strip.size() && (strip[j].y - strip[i].y) < minDist; j++) {
                minDist = min(minDist, distance(strip[i], strip[j]));
            }
        }
        
        return minDist;
    }
    
    static double closestPairRecWithPoints(vector<Point>& points, int left, int right, pair<Point, Point>& result) {
        if (right - left <= 3) {
            return bruteForceWithPoints(points, left, right, result);
        }
        
        int mid = left + (right - left) / 2;
        Point midPoint = points[mid];
        
        pair<Point, Point> leftResult, rightResult;
        double leftMin = closestPairRecWithPoints(points, left, mid, leftResult);
        double rightMin = closestPairRecWithPoints(points, mid + 1, right, rightResult);
        
        double minDist;
        if (leftMin < rightMin) {
            minDist = leftMin;
            result = leftResult;
        } else {
            minDist = rightMin;
            result = rightResult;
        }
        
        vector<Point> strip;
        for (int i = left; i <= right; i++) {
            if (abs(points[i].x - midPoint.x) < minDist) {
                strip.push_back(points[i]);
            }
        }
        
        sort(strip.begin(), strip.end(), [](const Point& a, const Point& b) {
            return a.y < b.y;
        });
        
        for (int i = 0; i < strip.size(); i++) {
            for (int j = i + 1; j < strip.size() && (strip[j].y - strip[i].y) < minDist; j++) {
                double dist = distance(strip[i], strip[j]);
                if (dist < minDist) {
                    minDist = dist;
                    result = {strip[i], strip[j]};
                }
            }
        }
        
        return minDist;
    }
    
    static double bruteForceWithPoints(vector<Point>& points, int left, int right, pair<Point, Point>& result) {
        double minDist = DBL_MAX;
        for (int i = left; i <= right; i++) {
            for (int j = i + 1; j <= right; j++) {
                double dist = distance(points[i], points[j]);
                if (dist < minDist) {
                    minDist = dist;
                    result = {points[i], points[j]};
                }
            }
        }
        return minDist;
    }
};
```

### 3. Fast Matrix Multiplication (Strassen's Algorithm)
```cpp
class MatrixMultiplication {
public:
    static vector<vector<int>> multiply(vector<vector<int>>& A, vector<vector<int>>& B) {
        int n = A.size();
        vector<vector<int>> C(n, vector<int>(n, 0));
        
        if (n <= 64) { // Use standard multiplication for small matrices
            return standardMultiply(A, B);
        }
        
        return strassenMultiply(A, B);
    }
    
private:
    static vector<vector<int>> standardMultiply(vector<vector<int>>& A, vector<vector<int>>& B) {
        int n = A.size();
        vector<vector<int>> C(n, vector<int>(n, 0));
        
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                for (int k = 0; k < n; k++) {
                    C[i][j] += A[i][k] * B[k][j];
                }
            }
        }
        
        return C;
    }
    
    static vector<vector<int>> strassenMultiply(vector<vector<int>>& A, vector<vector<int>>& B) {
        int n = A.size();
        
        if (n == 1) {
            return {{A[0][0] * B[0][0]}};
        }
        
        int half = n / 2;
        
        // Divide matrices into quadrants
        auto A11 = getSubMatrix(A, 0, 0, half);
        auto A12 = getSubMatrix(A, 0, half, half);
        auto A21 = getSubMatrix(A, half, 0, half);
        auto A22 = getSubMatrix(A, half, half, half);
        
        auto B11 = getSubMatrix(B, 0, 0, half);
        auto B12 = getSubMatrix(B, 0, half, half);
        auto B21 = getSubMatrix(B, half, 0, half);
        auto B22 = getSubMatrix(B, half, half, half);
        
        // Calculate the 7 products
        auto M1 = strassenMultiply(addMatrices(A11, A22), addMatrices(B11, B22));
        auto M2 = strassenMultiply(addMatrices(A21, A22), B11);
        auto M3 = strassenMultiply(A11, subtractMatrices(B12, B22));
        auto M4 = strassenMultiply(A22, subtractMatrices(B21, B11));
        auto M5 = strassenMultiply(addMatrices(A11, A12), B22);
        auto M6 = strassenMultiply(subtractMatrices(A21, A11), addMatrices(B11, B12));
        auto M7 = strassenMultiply(subtractMatrices(A12, A22), addMatrices(B21, B22));
        
        // Calculate result quadrants
        auto C11 = addMatrices(subtractMatrices(addMatrices(M1, M4), M5), M7);
        auto C12 = addMatrices(M3, M5);
        auto C21 = addMatrices(M2, M4);
        auto C22 = addMatrices(subtractMatrices(addMatrices(M1, M3), M2), M6);
        
        // Combine quadrants
        return combineMatrices(C11, C12, C21, C22);
    }
    
    static vector<vector<int>> getSubMatrix(vector<vector<int>>& matrix, int row, int col, int size) {
        vector<vector<int>> sub(size, vector<int>(size));
        for (int i = 0; i < size; i++) {
            for (int j = 0; j < size; j++) {
                sub[i][j] = matrix[row + i][col + j];
            }
        }
        return sub;
    }
    
    static vector<vector<int>> addMatrices(vector<vector<int>>& A, vector<vector<int>>& B) {
        int n = A.size();
        vector<vector<int>> result(n, vector<int>(n));
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                result[i][j] = A[i][j] + B[i][j];
            }
        }
        return result;
    }
    
    static vector<vector<int>> subtractMatrices(vector<vector<int>>& A, vector<vector<int>>& B) {
        int n = A.size();
        vector<vector<int>> result(n, vector<int>(n));
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n; j++) {
                result[i][j] = A[i][j] - B[i][j];
            }
        }
        return result;
    }
    
    static vector<vector<int>> combineMatrices(vector<vector<int>>& C11, vector<vector<int>>& C12,
                                              vector<vector<int>>& C21, vector<vector<int>>& C22) {
        int half = C11.size();
        int n = 2 * half;
        vector<vector<int>> result(n, vector<int>(n));
        
        for (int i = 0; i < half; i++) {
            for (int j = 0; j < half; j++) {
                result[i][j] = C11[i][j];
                result[i][j + half] = C12[i][j];
                result[i + half][j] = C21[i][j];
                result[i + half][j + half] = C22[i][j];
            }
        }
        
        return result;
    }
};
```

---

## Advanced Graph Algorithms {#advanced-graph-algorithms}

### 1. Network Flow - Ford-Fulkerson Algorithm
```cpp
class MaxFlow {
public:
    static int fordFulkerson(vector<vector<int>>& capacity, int source, int sink) {
        int n = capacity.size();
        vector<vector<int>> residualGraph = capacity;
        vector<int> parent(n);
        int maxFlow = 0;
        
        while (bfs(residualGraph, source, sink, parent)) {
            // Find minimum residual capacity along the path
            int pathFlow = INT_MAX;
            for (int v = sink; v != source; v = parent[v]) {
                int u = parent[v];
                pathFlow = min(pathFlow, residualGraph[u][v]);
            }
            
            // Update residual capacities
            for (int v = sink; v != source; v = parent[v]) {
                int u = parent[v];
                residualGraph[u][v] -= pathFlow;
                residualGraph[v][u] += pathFlow;
            }
            
            maxFlow += pathFlow;
        }
        
        return maxFlow;
    }
    
    // Edmonds-Karp (BFS-based Ford-Fulkerson)
    static int edmondsKarp(vector<vector<int>>& capacity, int source, int sink) {
        return fordFulkerson(capacity, source, sink);
    }
    
    // Dinic's Algorithm (faster for dense graphs)
    static int dinic(vector<vector<int>>& capacity, int source, int sink) {
        int n = capacity.size();
        vector<vector<int>> residualGraph = capacity;
        int maxFlow = 0;
        
        while (true) {
            vector<int> level = bfsLevel(residualGraph, source, sink);
            if (level[sink] == -1) break;
            
            vector<int> iter(n, 0);
            int flow;
            while ((flow = dfsFlow(residualGraph, level, iter, source, sink, INT_MAX)) > 0) {
                maxFlow += flow;
            }
        }
        
        return maxFlow;
    }
    
private:
    static bool bfs(vector<vector<int>>& graph, int source, int sink, vector<int>& parent) {
        int n = graph.size();
        vector<bool> visited(n, false);
        queue<int> q;
        
        q.push(source);
        visited[source] = true;
        parent[source] = -1;
        
        while (!q.empty()) {
            int u = q.front();
            q.pop();
            
            for (int v = 0; v < n; v++) {
                if (!visited[v] && graph[u][v] > 0) {
                    q.push(v);
                    parent[v] = u;
                    visited[v] = true;
                    if (v == sink) return true;
                }
            }
        }
        
        return false;
    }
    
    static vector<int> bfsLevel(vector<vector<int>>& graph, int source, int sink) {
        int n = graph.size();
        vector<int> level(n, -1);
        queue<int> q;
        
        level[source] = 0;
        q.push(source);
        
        while (!q.empty()) {
            int u = q.front();
            q.pop();
            
            for (int v = 0; v < n; v++) {
                if (level[v] == -1 && graph[u][v] > 0) {
                    level[v] = level[u] + 1;
                    q.push(v);
                }
            }
        }
        
        return level;
    }
    
    static int dfsFlow(vector<vector<int>>& graph, vector<int>& level, vector<int>& iter,
                      int u, int sink, int flow) {
        if (u == sink) return flow;
        
        int n = graph.size();
        for (int& v = iter[u]; v < n; v++) {
            if (level[v] == level[u] + 1 && graph[u][v] > 0) {
                int minFlow = min(flow, graph[u][v]);
                int pushed = dfsFlow(graph, level, iter, v, sink, minFlow);
                
                if (pushed > 0) {
                    graph[u][v] -= pushed;
                    graph[v][u] += pushed;
                    return pushed;
                }
            }
        }
        
        return 0;
    }
};
```

### 2. Bipartite Matching
```cpp
class BipartiteMatching {
public:
    // Maximum bipartite matching using DFS
    static int maxBipartiteMatching(vector<vector<int>>& graph) {
        int n = graph.size();
        vector<int> match(n, -1);
        int result = 0;
        
        for (int u = 0; u < n; u++) {
            vector<bool> visited(n, false);
            if (dfs(graph, u, visited, match)) {
                result++;
            }
        }
        
        return result;
    }
    
    // Hungarian Algorithm for assignment problem
    static int hungarianAlgorithm(vector<vector<int>>& cost) {
        int n = cost.size();
        vector<int> u(n + 1), v(n + 1), p(n + 1), way(n + 1);
        
        for (int i = 1; i <= n; i++) {
            p[0] = i;
            int j0 = 0;
            vector<int> minv(n + 1, INT_MAX);
            vector<bool> used(n + 1, false);
            
            do {
                used[j0] = true;
                int i0 = p[j0], delta = INT_MAX, j1;
                
                for (int j = 1; j <= n; j++) {
                    if (!used[j]) {
                        int cur = cost[i0 - 1][j - 1] - u[i0] - v[j];
                        if (cur < minv[j]) {
                            minv[j] = cur;
                            way[j] = j0;
                        }
                        if (minv[j] < delta) {
                            delta = minv[j];
                            j1 = j;
                        }
                    }
                }
                
                for (int j = 0; j <= n; j++) {
                    if (used[j]) {
                        u[p[j]] += delta;
                        v[j] -= delta;
                    } else {
                        minv[j] -= delta;
                    }
                }
                
                j0 = j1;
            } while (p[j0] != 0);
            
            do {
                int j1 = way[j0];
                p[j0] = p[j1];
                j0 = j1;
            } while (j0);
        }
        
        return -v[0];
    }
    
private:
    static bool dfs(vector<vector<int>>& graph, int u, vector<bool>& visited, vector<int>& match) {
        for (int v : graph[u]) {
            if (visited[v]) continue;
            visited[v] = true;
            
            if (match[v] == -1 || dfs(graph, match[v], visited, match)) {
                match[v] = u;
                return true;
            }
        }
        return false;
    }
};
```

---

## String Algorithms {#string-algorithms}

### 1. KMP (Knuth-Morris-Pratt) Algorithm
```cpp
class KMPAlgorithm {
public:
    static vector<int> findPattern(string text, string pattern) {
        vector<int> lps = computeLPS(pattern);
        vector<int> matches;
        
        int i = 0; // index for text
        int j = 0; // index for pattern
        
        while (i < text.length()) {
            if (pattern[j] == text[i]) {
                i++;
                j++;
            }
            
            if (j == pattern.length()) {
                matches.push_back(i - j);
                j = lps[j - 1];
            } else if (i < text.length() && pattern[j] != text[i]) {
                if (j != 0) {
                    j = lps[j - 1];
                } else {
                    i++;
                }
            }
        }
        
        return matches;
    }
    
    static vector<int> computeLPS(string pattern) {
        int m = pattern.length();
        vector<int> lps(m, 0);
        int len = 0;
        int i = 1;
        
        while (i < m) {
            if (pattern[i] == pattern[len]) {
                len++;
                lps[i] = len;
                i++;
            } else {
                if (len != 0) {
                    len = lps[len - 1];
                } else {
                    lps[i] = 0;
                    i++;
                }
            }
        }
        
        return lps;
    }
};
```

### 2. Rabin-Karp Algorithm
```cpp
class RabinKarp {
public:
    static vector<int> search(string text, string pattern) {
        const int PRIME = 101;
        int m = pattern.length();
        int n = text.length();
        vector<int> matches;
        
        if (m > n) return matches;
        
        long long patternHash = 0;
        long long textHash = 0;
        long long h = 1;
        
        // Calculate h = pow(256, m-1) % PRIME
        for (int i = 0; i < m - 1; i++) {
            h = (h * 256) % PRIME;
        }
        
        // Calculate hash for pattern and first window of text
        for (int i = 0; i < m; i++) {
            patternHash = (256 * patternHash + pattern[i]) % PRIME;
            textHash = (256 * textHash + text[i]) % PRIME;
        }
        
        // Slide the pattern over text
        for (int i = 0; i <= n - m; i++) {
            if (patternHash == textHash) {
                // Check characters one by one
                bool match = true;
                for (int j = 0; j < m; j++) {
                    if (text[i + j] != pattern[j]) {
                        match = false;
                        break;
                    }
                }
                if (match) {
                    matches.push_back(i);
                }
            }
            
            // Calculate hash for next window
            if (i < n - m) {
                textHash = (256 * (textHash - text[i] * h) + text[i + m]) % PRIME;
                if (textHash < 0) {
                    textHash += PRIME;
                }
            }
        }
        
        return matches;
    }
};
```

### 3. Suffix Array and LCP Array
```cpp
class SuffixArray {
public:
    static vector<int> buildSuffixArray(string s) {
        int n = s.length();
        vector<int> sa(n), rank(n), temp(n);
        
        // Initial ranking
        for (int i = 0; i < n; i++) {
            sa[i] = i;
            rank[i] = s[i];
        }
        
        for (int k = 1; k < n; k *= 2) {
            // Sort by (rank[i], rank[i+k])
            sort(sa.begin(), sa.end(), [&](int i, int j) {
                if (rank[i] != rank[j]) return rank[i] < rank[j];
                int ri = (i + k < n) ? rank[i + k] : -1;
                int rj = (j + k < n) ? rank[j + k] : -1;
                return ri < rj;
            });
            
            // Update ranks
            temp[sa[0]] = 0;
            for (int i = 1; i < n; i++) {
                temp[sa[i]] = temp[sa[i-1]];
                if (rank[sa[i-1]] != rank[sa[i]] || 
                    (sa[i-1] + k < n ? rank[sa[i-1] + k] : -1) != 
                    (sa[i] + k < n ? rank[sa[i] + k] : -1)) {
                    temp[sa[i]]++;
                }
            }
            
            rank = temp;
        }
        
        return sa;
    }
    
    static vector<int> buildLCPArray(string s, vector<int>& sa) {
        int n = s.length();
        vector<int> rank(n), lcp(n - 1);
        
        // Build rank array
        for (int i = 0; i < n; i++) {
            rank[sa[i]] = i;
        }
        
        int h = 0;
        for (int i = 0; i < n; i++) {
            if (rank[i] > 0) {
                int j = sa[rank[i] - 1];
                while (i + h < n && j + h < n && s[i + h] == s[j + h]) {
                    h++;
                }
                lcp[rank[i] - 1] = h;
                if (h > 0) h--;
            }
        }
        
        return lcp;
    }
    
    // Find longest repeated substring
    static string longestRepeatedSubstring(string s) {
        vector<int> sa = buildSuffixArray(s);
        vector<int> lcp = buildLCPArray(s, sa);
        
        int maxLen = 0;
        int maxIndex = 0;
        
        for (int i = 0; i < lcp.size(); i++) {
            if (lcp[i] > maxLen) {
                maxLen = lcp[i];
                maxIndex = sa[i];
            }
        }
        
        return s.substr(maxIndex, maxLen);
    }
};
```

---

## Number Theory Algorithms {#number-theory}

### 1. Prime Numbers and Factorization
```cpp
class NumberTheory {
public:
    // Sieve of Eratosthenes
    static vector<bool> sieveOfEratosthenes(int n) {
        vector<bool> isPrime(n + 1, true);
        isPrime[0] = isPrime[1] = false;
        
        for (int i = 2; i * i <= n; i++) {
            if (isPrime[i]) {
                for (int j = i * i; j <= n; j += i) {
                    isPrime[j] = false;
                }
            }
        }
        
        return isPrime;
    }
    
    // Segmented Sieve
    static vector<int> segmentedSieve(int low, int high) {
        int limit = sqrt(high) + 1;
        vector<bool> mark(limit + 1, true);
        vector<int> primes;
        
        // Find all primes up to sqrt(high)
        for (int i = 2; i <= limit; i++) {
            if (mark[i]) {
                primes.push_back(i);
                for (int j = i * i; j <= limit; j += i) {
                    mark[j] = false;
                }
            }
        }
        
        // Create boolean array for [low, high]
        vector<bool> isPrime(high - low + 1, true);
        
        // Mark non-primes in [low, high]
        for (int prime : primes) {
            int start = max(prime * prime, (low + prime - 1) / prime * prime);
            for (int j = start; j <= high; j += prime) {
                isPrime[j - low] = false;
            }
        }
        
        vector<int> result;
        for (int i = low; i <= high; i++) {
            if (i > 1 && isPrime[i - low]) {
                result.push_back(i);
            }
        }
        
        return result;
    }
    
    // Prime factorization
    static vector<pair<int, int>> primeFactorization(int n) {
        vector<pair<int, int>> factors;
        
        // Check for 2
        int count = 0;
        while (n % 2 == 0) {
            count++;
            n /= 2;
        }
        if (count > 0) {
            factors.push_back({2, count});
        }
        
        // Check for odd factors
        for (int i = 3; i * i <= n; i += 2) {
            count = 0;
            while (n % i == 0) {
                count++;
                n /= i;
            }
            if (count > 0) {
                factors.push_back({i, count});
            }
        }
        
        // If n is still > 1, it's a prime
        if (n > 1) {
            factors.push_back({n, 1});
        }
        
        return factors;
    }
    
    // Euler's totient function
    static int eulerTotient(int n) {
        int result = n;
        
        for (int i = 2; i * i <= n; i++) {
            if (n % i == 0) {
                while (n % i == 0) {
                    n /= i;
                }
                result -= result / i;
            }
        }
        
        if (n > 1) {
            result -= result / n;
        }
        
        return result;
    }
    
    // Extended Euclidean Algorithm
    static tuple<int, int, int> extendedGCD(int a, int b) {
        if (b == 0) {
            return {a, 1, 0};
        }
        
        auto [gcd, x1, y1] = extendedGCD(b, a % b);
        int x = y1;
        int y = x1 - (a / b) * y1;
        
        return {gcd, x, y};
    }
    
    // Modular exponentiation
    static long long modPow(long long base, long long exp, long long mod) {
        long long result = 1;
        base %= mod;
        
        while (exp > 0) {
            if (exp & 1) {
                result = (result * base) % mod;
            }
            base = (base * base) % mod;
            exp >>= 1;
        }
        
        return result;
    }
    
    // Modular inverse
    static long long modInverse(long long a, long long mod) {
        auto [gcd, x, y] = extendedGCD(a, mod);
        if (gcd != 1) {
            return -1; // Inverse doesn't exist
        }
        return (x % mod + mod) % mod;
    }
    
    // Chinese Remainder Theorem
    static long long chineseRemainder(vector<long long>& remainders, vector<long long>& moduli) {
        long long prod = 1;
        for (long long mod : moduli) {
            prod *= mod;
        }
        
        long long result = 0;
        for (int i = 0; i < remainders.size(); i++) {
            long long partialProd = prod / moduli[i];
            long long inv = modInverse(partialProd, moduli[i]);
            result = (result + remainders[i] * partialProd * inv) % prod;
        }
        
        return (result + prod) % prod;
    }
};
```

### 2. Fast Fourier Transform (FFT)
```cpp
class FFT {
public:
    using Complex = complex<double>;
    const double PI = acos(-1);
    
    static vector<Complex> fft(vector<Complex> a, bool invert = false) {
        int n = a.size();
        
        // Bit-reversal permutation
        for (int i = 1, j = 0; i < n; i++) {
            int bit = n >> 1;
            for (; j & bit; bit >>= 1) {
                j ^= bit;
            }
            j ^= bit;
            
            if (i < j) {
                swap(a[i], a[j]);
            }
        }
        
        // FFT computation
        for (int len = 2; len <= n; len <<= 1) {
            double ang = 2 * PI / len * (invert ? -1 : 1);
            Complex wlen(cos(ang), sin(ang));
            
            for (int i = 0; i < n; i += len) {
                Complex w(1);
                for (int j = 0; j < len / 2; j++) {
                    Complex u = a[i + j];
                    Complex v = a[i + j + len / 2] * w;
                    a[i + j] = u + v;
                    a[i + j + len / 2] = u - v;
                    w *= wlen;
                }
            }
        }
        
        if (invert) {
            for (Complex& x : a) {
                x /= n;
            }
        }
        
        return a;
    }
    
    // Polynomial multiplication using FFT
    static vector<int> multiply(vector<int>& a, vector<int>& b) {
        vector<Complex> fa(a.begin(), a.end());
        vector<Complex> fb(b.begin(), b.end());
        
        int n = 1;
        while (n < a.size() + b.size()) {
            n <<= 1;
        }
        
        fa.resize(n);
        fb.resize(n);
        
        fa = fft(fa, false);
        fb = fft(fb, false);
        
        for (int i = 0; i < n; i++) {
            fa[i] *= fb[i];
        }
        
        fa = fft(fa, true);
        
        vector<int> result(n);
        for (int i = 0; i < n; i++) {
            result[i] = round(fa[i].real());
        }
        
        // Remove leading zeros
        while (result.size() > 1 && result.back() == 0) {
            result.pop_back();
        }
        
        return result;
    }
};
```

---

## Computational Geometry {#computational-geometry}

### 1. Basic Geometric Operations
```cpp
class ComputationalGeometry {
public:
    struct Point {
        double x, y;
        Point(double x = 0, double y = 0) : x(x), y(y) {}
        
        Point operator+(const Point& p) const { return Point(x + p.x, y + p.y); }
        Point operator-(const Point& p) const { return Point(x - p.x, y - p.y); }
        Point operator*(double t) const { return Point(x * t, y * t); }
        double dot(const Point& p) const { return x * p.x + y * p.y; }
        double cross(const Point& p) const { return x * p.y - y * p.x; }
        double norm() const { return sqrt(x * x + y * y); }
        double norm2() const { return x * x + y * y; }
    };
    
    // Check if three points are collinear
    static bool areCollinear(Point a, Point b, Point c) {
        return abs((b - a).cross(c - a)) < 1e-9;
    }
    
    // Check orientation of three points
    static int orientation(Point a, Point b, Point c) {
        double val = (b - a).cross(c - a);
        if (abs(val) < 1e-9) return 0; // Collinear
        return (val > 0) ? 1 : 2; // Clockwise or Counterclockwise
    }
    
    // Check if point lies on line segment
    static bool onSegment(Point p, Point q, Point r) {
        return q.x <= max(p.x, r.x) && q.x >= min(p.x, r.x) &&
               q.y <= max(p.y, r.y) && q.y >= min(p.y, r.y);
    }
    
    // Check if two line segments intersect
    static bool doIntersect(Point p1, Point q1, Point p2, Point q2) {
        int o1 = orientation(p1, q1, p2);
        int o2 = orientation(p1, q1, q2);
        int o3 = orientation(p2, q2, p1);
        int o4 = orientation(p2, q2, q1);
        
        // General case
        if (o1 != o2 && o3 != o4) return true;
        
        // Special cases
        if (o1 == 0 && onSegment(p1, p2, q1)) return true;
        if (o2 == 0 && onSegment(p1, q2, q1)) return true;
        if (o3 == 0 && onSegment(p2, p1, q2)) return true;
        if (o4 == 0 && onSegment(p2, q1, q2)) return true;
        
        return false;
    }
    
    // Convex Hull using Graham Scan
    static vector<Point> convexHull(vector<Point> points) {
        int n = points.size();
        if (n < 3) return {};
        
        // Find bottom-most point (or left most in case of tie)
        int l = 0;
        for (int i = 1; i < n; i++) {
            if (points[i].y < points[l].y || 
                (points[i].y == points[l].y && points[i].x < points[l].x)) {
                l = i;
            }
        }
        
        swap(points[0], points[l]);
        Point pivot = points[0];
        
        // Sort points by polar angle with respect to pivot
        sort(points.begin() + 1, points.end(), [&](const Point& a, const Point& b) {
            int o = orientation(pivot, a, b);
            if (o == 0) {
                return (pivot - a).norm2() < (pivot - b).norm2();
            }
            return o == 2; // Counterclockwise
        });
        
        // Create convex hull
        vector<Point> hull;
        for (Point p : points) {
            while (hull.size() > 1 && 
                   orientation(hull[hull.size()-2], hull[hull.size()-1], p) != 2) {
                hull.pop_back();
            }
            hull.push_back(p);
        }
        
        return hull;
    }
    
    // Point in polygon test
    static bool pointInPolygon(Point point, vector<Point>& polygon) {
        int n = polygon.size();
        bool inside = false;
        
        for (int i = 0, j = n - 1; i < n; j = i++) {
            if (((polygon[i].y > point.y) != (polygon[j].y > point.y)) &&
                (point.x < (polygon[j].x - polygon[i].x) * (point.y - polygon[i].y) / 
                          (polygon[j].y - polygon[i].y) + polygon[i].x)) {
                inside = !inside;
            }
        }
        
        return inside;
    }
    
    // Area of polygon
    static double polygonArea(vector<Point>& polygon) {
        double area = 0;
        int n = polygon.size();
        
        for (int i = 0; i < n; i++) {
            int j = (i + 1) % n;
            area += polygon[i].cross(polygon[j]);
        }
        
        return abs(area) / 2.0;
    }
};
```

---

## Advanced Data Structures {#advanced-data-structures}

### 1. Segment Tree
```cpp
class SegmentTree {
public:
    vector<int> tree;
    int n;
    
    SegmentTree(vector<int>& arr) {
        n = arr.size();
        tree.resize(4 * n);
        build(arr, 1, 0, n - 1);
    }
    
    void update(int pos, int val) {
        update(1, 0, n - 1, pos, val);
    }
    
    int query(int l, int r) {
        return query(1, 0, n - 1, l, r);
    }
    
private:
    void build(vector<int>& arr, int v, int tl, int tr) {
        if (tl == tr) {
            tree[v] = arr[tl];
        } else {
            int tm = (tl + tr) / 2;
            build(arr, 2*v, tl, tm);
            build(arr, 2*v+1, tm+1, tr);
            tree[v] = tree[2*v] + tree[2*v+1];
        }
    }
    
    void update(int v, int tl, int tr, int pos, int val) {
        if (tl == tr) {
            tree[v] = val;
        } else {
            int tm = (tl + tr) / 2;
            if (pos <= tm)
                update(2*v, tl, tm, pos, val);
            else
                update(2*v+1, tm+1, tr, pos, val);
            tree[v] = tree[2*v] + tree[2*v+1];
        }
    }
    
    int query(int v, int tl, int tr, int l, int r) {
        if (l > r) return 0;
        if (l == tl && r == tr) {
            return tree[v];
        }
        int tm = (tl + tr) / 2;
        return query(2*v, tl, tm, l, min(r, tm)) +
               query(2*v+1, tm+1, tr, max(l, tm+1), r);
    }
};

// Lazy Propagation Segment Tree
class LazySegmentTree {
public:
    vector<long long> tree, lazy;
    int n;
    
    LazySegmentTree(vector<int>& arr) {
        n = arr.size();
        tree.resize(4 * n);
        lazy.resize(4 * n);
        build(arr, 1, 0, n - 1);
    }
    
    void updateRange(int l, int r, int val) {
        updateRange(1, 0, n - 1, l, r, val);
    }
    
    long long queryRange(int l, int r) {
        return queryRange(1, 0, n - 1, l, r);
    }
    
private:
    void build(vector<int>& arr, int v, int tl, int tr) {
        if (tl == tr) {
            tree[v] = arr[tl];
        } else {
            int tm = (tl + tr) / 2;
            build(arr, 2*v, tl, tm);
            build(arr, 2*v+1, tm+1, tr);
            tree[v] = tree[2*v] + tree[2*v+1];
        }
    }
    
    void push(int v, int tl, int tr) {
        if (lazy[v] != 0) {
            tree[v] += lazy[v] * (tr - tl + 1);
            if (tl != tr) {
                lazy[2*v] += lazy[v];
                lazy[2*v+1] += lazy[v];
            }
            lazy[v] = 0;
        }
    }
    
    void updateRange(int v, int tl, int tr, int l, int r, int val) {
        if (l > r) return;
        if (l == tl && r == tr) {
            lazy[v] += val;
            push(v, tl, tr);
            return;
        }
        push(v, tl, tr);
        int tm = (tl + tr) / 2;
        updateRange(2*v, tl, tm, l, min(r, tm), val);
        updateRange(2*v+1, tm+1, tr, max(l, tm+1), r, val);
        push(2*v, tl, tm);
        push(2*v+1, tm+1, tr);
        tree[v] = tree[2*v] + tree[2*v+1];
    }
    
    long long queryRange(int v, int tl, int tr, int l, int r) {
        if (l > r) return 0;
        push(v, tl, tr);
        if (l == tl && r == tr) {
            return tree[v];
        }
        int tm = (tl + tr) / 2;
        return queryRange(2*v, tl, tm, l, min(r, tm)) +
               queryRange(2*v+1, tm+1, tr, max(l, tm+1), r);
    }
};
```

### 2. Fenwick Tree (Binary Indexed Tree)
```cpp
class FenwickTree {
public:
    vector<int> tree;
    int n;
    
    FenwickTree(int size) {
        n = size;
        tree.assign(n + 1, 0);
    }
    
    FenwickTree(vector<int>& arr) {
        n = arr.size();
        tree.assign(n + 1, 0);
        for (int i = 0; i < n; i++) {
            update(i, arr[i]);
        }
    }
    
    void update(int idx, int delta) {
        for (++idx; idx <= n; idx += idx & -idx) {
            tree[idx] += delta;
        }
    }
    
    int query(int idx) {
        int sum = 0;
        for (++idx; idx > 0; idx -= idx & -idx) {
            sum += tree[idx];
        }
        return sum;
    }
    
    int rangeQuery(int l, int r) {
        return query(r) - (l > 0 ? query(l - 1) : 0);
    }
};

// 2D Fenwick Tree
class FenwickTree2D {
public:
    vector<vector<int>> tree;
    int n, m;
    
    FenwickTree2D(int rows, int cols) {
        n = rows;
        m = cols;
        tree.assign(n + 1, vector<int>(m + 1, 0));
    }
    
    void update(int x, int y, int delta) {
        for (int i = x + 1; i <= n; i += i & -i) {
            for (int j = y + 1; j <= m; j += j & -j) {
                tree[i][j] += delta;
            }
        }
    }
    
    int query(int x, int y) {
        int sum = 0;
        for (int i = x + 1; i > 0; i -= i & -i) {
            for (int j = y + 1; j > 0; j -= j & -j) {
                sum += tree[i][j];
            }
        }
        return sum;
    }
    
    int rangeQuery(int x1, int y1, int x2, int y2) {
        return query(x2, y2) - query(x1 - 1, y2) - query(x2, y1 - 1) + query(x1 - 1, y1 - 1);
    }
};
```

### 3. Disjoint Set Union (Union-Find)
```cpp
class DSU {
public:
    vector<int> parent, rank, size;
    int components;
    
    DSU(int n) {
        parent.resize(n);
        rank.resize(n, 0);
        size.resize(n, 1);
        components = n;
        iota(parent.begin(), parent.end(), 0);
    }
    
    int find(int x) {
        if (parent[x] != x) {
            parent[x] = find(parent[x]); // Path compression
        }
        return parent[x];
    }
    
    bool unite(int x, int y) {
        int px = find(x), py = find(y);
        if (px == py) return false;
        
        // Union by rank
        if (rank[px] < rank[py]) {
            parent[px] = py;
            size[py] += size[px];
        } else if (rank[px] > rank[py]) {
            parent[py] = px;
            size[px] += size[py];
        } else {
            parent[py] = px;
            size[px] += size[py];
            rank[px]++;
        }
        
        components--;
        return true;
    }
    
    bool connected(int x, int y) {
        return find(x) == find(y);
    }
    
    int getSize(int x) {
        return size[find(x)];
    }
    
    int getComponents() {
        return components;
    }
};
```

---

## Optimization Techniques {#optimization-techniques}

### 1. Memoization and Caching
```cpp
class Memoization {
public:
    // Fibonacci with memoization
    static unordered_map<int, long long> fibMemo;
    
    static long long fibonacci(int n) {
        if (n <= 1) return n;
        if (fibMemo.find(n) != fibMemo.end()) {
            return fibMemo[n];
        }
        return fibMemo[n] = fibonacci(n-1) + fibonacci(n-2);
    }
    
    // LCS with memoization
    static unordered_map<string, int> lcsMemo;
    
    static int lcs(string& s1, string& s2, int i, int j) {
        if (i == s1.length() || j == s2.length()) return 0;
        
        string key = to_string(i) + "," + to_string(j);
        if (lcsMemo.find(key) != lcsMemo.end()) {
            return lcsMemo[key];
        }
        
        if (s1[i] == s2[j]) {
            return lcsMemo[key] = 1 + lcs(s1, s2, i+1, j+1);
        } else {
            return lcsMemo[key] = max(lcs(s1, s2, i+1, j), lcs(s1, s2, i, j+1));
        }
    }
};

unordered_map<int, long long> Memoization::fibMemo;
unordered_map<string, int> Memoization::lcsMemo;
```

### 2. Bit Manipulation Optimizations
```cpp
class BitOptimizations {
public:
    // Check if number is power of 2
    static bool isPowerOfTwo(int n) {
        return n > 0 && (n & (n - 1)) == 0;
    }
    
    // Count set bits
    static int countSetBits(int n) {
        int count = 0;
        while (n) {
            count++;
            n &= (n - 1); // Remove rightmost set bit
        }
        return count;
    }
    
    // Find rightmost set bit
    static int rightmostSetBit(int n) {
        return n & (-n);
    }
    
    // Subset enumeration
    static vector<vector<int>> generateSubsets(vector<int>& nums) {
        vector<vector<int>> subsets;
        int n = nums.size();
        
        for (int mask = 0; mask < (1 << n); mask++) {
            vector<int> subset;
            for (int i = 0; i < n; i++) {
                if (mask & (1 << i)) {
                    subset.push_back(nums[i]);
                }
            }
            subsets.push_back(subset);
        }
        
        return subsets;
    }
    
    // Iterate through all submasks of a mask
    static vector<int> getAllSubmasks(int mask) {
        vector<int> submasks;
        for (int submask = mask; ; submask = (submask - 1) & mask) {
            submasks.push_back(submask);
            if (submask == 0) break;
        }
        return submasks;
    }
};
```

---

## Practice Problems {#practice-problems}

### Basic Level
1. **Greedy Algorithms**
   - Activity Selection Problem
   - Fractional Knapsack
   - Job Scheduling with Deadlines
   - Minimum Number of Coins
   - Huffman Coding

2. **Divide and Conquer**
   - Maximum Subarray (Kadane's vs D&C)
   - Merge Sort Implementation
   - Quick Sort with Random Pivot
   - Binary Search Variations
   - Closest Pair of Points

3. **String Algorithms**
   - Pattern Matching (KMP, Rabin-Karp)
   - Longest Common Subsequence
   - Edit Distance
   - Palindrome Problems
   - Anagram Detection

### Intermediate Level
1. **Graph Algorithms**
   - Network Flow Problems
   - Bipartite Matching
   - Strongly Connected Components
   - Articulation Points and Bridges
   - Minimum Cut Problems

2. **Advanced Data Structures**
   - Range Query Problems (Segment Tree)
   - Dynamic Range Updates (Lazy Propagation)
   - 2D Range Queries (2D BIT)
   - Union-Find Applications
   - Persistent Data Structures

3. **Number Theory**
   - Prime Factorization
   - Modular Arithmetic
   - Chinese Remainder Theorem
   - Fast Exponentiation
   - Euler's Totient Function

### Advanced Level
1. **Computational Geometry**
   - Convex Hull Algorithms
   - Line Intersection
   - Point in Polygon
   - Closest Pair of Points
   - Voronoi Diagrams

2. **Advanced Algorithms**
   - Fast Fourier Transform
   - Matrix Exponentiation
   - Heavy-Light Decomposition
   - Centroid Decomposition
   - Suffix Arrays and LCP

3. **Optimization Problems**
   - Linear Programming
   - Assignment Problems
   - Traveling Salesman Problem
   - Knapsack Variations
   - Scheduling Problems

---

## Complete Example Program {#example-program}

```cpp
#include <iostream>
#include <vector>
#include <algorithm>
#include <string>
#include <unordered_map>
#include <queue>
#include <stack>
#include <set>
#include <map>
#include <cmath>
#include <complex>
#include <chrono>
#include <random>

using namespace std;
using namespace std::chrono;

class AdvancedAlgorithmsDemo {
public:
    static void demonstrateAll() {
        cout << "=== Advanced Algorithms Demonstration ===\n\n";
        
        demonstrateGreedyAlgorithms();
        demonstrateDivideAndConquer();
        demonstrateStringAlgorithms();
        demonstrateAdvancedDataStructures();
        demonstrateNumberTheory();
        demonstrateComputationalGeometry();
        
        cout << "\n=== All demonstrations completed! ===\n";
    }
    
private:
    static void demonstrateGreedyAlgorithms() {
        cout << "--- Greedy Algorithms ---\n";
        
        // Activity Selection
        vector<pair<int, int>> activities = {{1, 4}, {3, 5}, {0, 6}, {5, 7}, {3, 9}, {5, 9}, {6, 10}, {8, 11}, {8, 12}, {2, 14}, {12, 16}};
        auto selected = ActivitySelection::selectActivities(activities);
        cout << "Selected activities: ";
        for (int idx : selected) {
            cout << "(" << activities[idx].first << "," << activities[idx].second << ") ";
        }
        cout << "\n";
        
        // Fractional Knapsack
        vector<int> values = {60, 100, 120};
        vector<int> weights = {10, 20, 30};
        int capacity = 50;
        double maxValue = FractionalKnapsack::fractionalKnapsack(values, weights, capacity);
        cout << "Maximum value in fractional knapsack: " << maxValue << "\n";
        
        // Huffman Coding
        unordered_map<char, int> frequencies = {{'a', 5}, {'b', 9}, {'c', 12}, {'d', 13}, {'e', 16}, {'f', 45}};
        auto codes = HuffmanCoding::buildHuffmanCodes(frequencies);
        cout << "Huffman codes: ";
        for (auto& pair : codes) {
            cout << pair.first << ":" << pair.second << " ";
        }
        cout << "\n\n";
    }
    
    static void demonstrateDivideAndConquer() {
        cout << "--- Divide and Conquer ---\n";
        
        // Maximum Subarray
        vector<int> arr = {-2, -3, 4, -1, -2, 1, 5, -3};
        int maxSum = MaxSubarrayDC::maxSubarraySum(arr);
        cout << "Maximum subarray sum: " << maxSum << "\n";
        
        // Closest Pair of Points
        vector<ClosestPair::Point> points = {
            {2, 3}, {12, 30}, {40, 50}, {5, 1}, {12, 10}, {3, 4}
        };
        double minDist = ClosestPair::closestPairDistance(points);
        cout << "Closest pair distance: " << minDist << "\n";
        
        // Matrix Multiplication
        vector<vector<int>> A = {{1, 2}, {3, 4}};
        vector<vector<int>> B = {{5, 6}, {7, 8}};
        auto C = MatrixMultiplication::multiply(A, B);
        cout << "Matrix multiplication result: ";
        for (auto& row : C) {
            for (int val : row) {
                cout << val << " ";
            }
        }
        cout << "\n\n";
    }
    
    static void demonstrateStringAlgorithms() {
        cout << "--- String Algorithms ---\n";
        
        // KMP Algorithm
        string text = "ABABDABACDABABCABCABCABCABC";
        string pattern = "ABABCABCABCABC";
        auto matches = KMPAlgorithm::findPattern(text, pattern);
        cout << "KMP matches at positions: ";
        for (int pos : matches) {
            cout << pos << " ";
        }
        cout << "\n";
        
        // Rabin-Karp Algorithm
        matches = RabinKarp::search(text, "ABC");
        cout << "Rabin-Karp matches for 'ABC': ";
        for (int pos : matches) {
            cout << pos << " ";
        }
        cout << "\n";
        
        // Suffix Array
        string s = "banana";
        auto sa = SuffixArray::buildSuffixArray(s);
        cout << "Suffix array for 'banana': ";
        for (int idx : sa) {
            cout << idx << " ";
        }
        cout << "\n\n";
    }
    
    static void demonstrateAdvancedDataStructures() {
        cout << "--- Advanced Data Structures ---\n";
        
        // Segment Tree
        vector<int> arr = {1, 3, 5, 7, 9, 11};
        SegmentTree st(arr);
        cout << "Sum of range [1, 3]: " << st.query(1, 3) << "\n";
        st.update(1, 10);
        cout << "After updating index 1 to 10, sum of range [1, 3]: " << st.query(1, 3) << "\n";
        
        // Fenwick Tree
        FenwickTree ft(arr);
        cout << "Fenwick tree sum [0, 2]: " << ft.rangeQuery(0, 2) << "\n";
        ft.update(0, 5);
        cout << "After adding 5 to index 0, sum [0, 2]: " << ft.rangeQuery(0, 2) << "\n";
        
        // Disjoint Set Union
        DSU dsu(6);
        dsu.unite(0, 1);
        dsu.unite(2, 3);
        dsu.unite(1, 2);
        cout << "Connected components: " << dsu.getComponents() << "\n";
        cout << "Size of component containing 0: " << dsu.getSize(0) << "\n\n";
    }
    
    static void demonstrateNumberTheory() {
        cout << "--- Number Theory ---\n";
        
        // Prime factorization
        int n = 315;
        auto factors = NumberTheory::primeFactorization(n);
        cout << "Prime factorization of " << n << ": ";
        for (auto& factor : factors) {
            cout << factor.first << "^" << factor.second << " ";
        }
        cout << "\n";
        
        // Euler's totient
        cout << "Euler's totient of " << n << ": " << NumberTheory::eulerTotient(n) << "\n";
        
        // Modular exponentiation
        long long base = 2, exp = 10, mod = 1000;
        cout << base << "^" << exp << " mod " << mod << " = " << NumberTheory::modPow(base, exp, mod) << "\n";
        
        // Extended GCD
        auto [gcd, x, y] = NumberTheory::extendedGCD(30, 18);
        cout << "Extended GCD(30, 18): gcd=" << gcd << ", x=" << x << ", y=" << y << "\n\n";
    }
    
    static void demonstrateComputationalGeometry() {
        cout << "--- Computational Geometry ---\n";
        
        // Convex Hull
        vector<ComputationalGeometry::Point> points = {
            {0, 3}, {1, 1}, {2, 2}, {4, 4}, {0, 0}, {1, 2}, {3, 1}, {3, 3}
        };
        auto hull = ComputationalGeometry::convexHull(points);
        cout << "Convex hull points: ";
        for (auto& p : hull) {
            cout << "(" << p.x << "," << p.y << ") ";
        }
        cout << "\n";
        
        // Point in polygon
        vector<ComputationalGeometry::Point> polygon = {{0, 0}, {4, 0}, {4, 4}, {0, 4}};
        ComputationalGeometry::Point testPoint(2, 2);
        bool inside = ComputationalGeometry::pointInPolygon(testPoint, polygon);
        cout << "Point (2,2) is " << (inside ? "inside" : "outside") << " the polygon\n";
        
        // Polygon area
        double area = ComputationalGeometry::polygonArea(polygon);
        cout << "Polygon area: " << area << "\n\n";
    }
};

int main() {
    AdvancedAlgorithmsDemo::demonstrateAll();
    return 0;
}
```

---

## Summary {#summary}

Part 10 covers advanced algorithms that represent the pinnacle of algorithmic problem-solving. These algorithms combine multiple concepts and require deep understanding of various paradigms.

### Key Concepts Covered

#### Algorithmic Paradigms
- **Greedy Algorithms**: Local optimization leading to global solutions
- **Divide and Conquer**: Breaking problems into smaller subproblems
- **Advanced Graph Algorithms**: Network flows, matching, connectivity
- **String Algorithms**: Pattern matching, suffix structures
- **Number Theory**: Prime numbers, modular arithmetic, GCD
- **Computational Geometry**: Geometric algorithms and data structures

#### Advanced Data Structures
- **Segment Trees**: Range queries and updates
- **Fenwick Trees**: Efficient prefix sum operations
- **Disjoint Set Union**: Union-find with path compression
- **Advanced Trees**: Heavy-light decomposition, centroid decomposition

#### Optimization Techniques
- **Memoization**: Caching for recursive algorithms
- **Bit Manipulation**: Efficient operations using bits
- **Space-Time Tradeoffs**: Balancing memory and computation
- **Lazy Propagation**: Deferred updates for efficiency

### Time and Space Complexities

| Algorithm/Data Structure | Time Complexity | Space Complexity | Use Case |
|-------------------------|-----------------|------------------|----------|
| Activity Selection | O(n log n) | O(1) | Scheduling |
| Huffman Coding | O(n log n) | O(n) | Compression |
| Ford-Fulkerson | O(E × f) | O(V²) | Max Flow |
| KMP Algorithm | O(n + m) | O(m) | Pattern Matching |
| Suffix Array | O(n log² n) | O(n) | String Processing |
| Segment Tree | O(log n) | O(n) | Range Queries |
| Fenwick Tree | O(log n) | O(n) | Prefix Sums |
| Union-Find | O(α(n)) | O(n) | Connectivity |
| Convex Hull | O(n log n) | O(n) | Geometry |
| FFT | O(n log n) | O(n) | Polynomial Mult |

### Problem-Solving Strategy

1. **Identify the Problem Type**
   - Optimization → Greedy or DP
   - Geometric → Computational Geometry
   - String → String Algorithms
   - Graph → Advanced Graph Algorithms

2. **Choose the Right Paradigm**
   - Greedy: Local optimum leads to global
   - Divide and Conquer: Subproblems are independent
   - Dynamic Programming: Overlapping subproblems
   - Network Flow: Flow/matching problems

3. **Consider Data Structure Requirements**
   - Range queries → Segment/Fenwick Tree
   - Connectivity → Union-Find
   - String operations → Suffix structures
   - Geometric queries → Spatial data structures

4. **Optimize for Constraints**
   - Large data → Efficient algorithms
   - Online queries → Preprocessing
   - Memory limits → Space optimization
   - Time limits → Faster algorithms

### When to Use Advanced Algorithms

- **Competitive Programming**: Complex problem-solving
- **System Design**: Efficient data processing
- **Research**: Novel algorithm development
- **Optimization**: Performance-critical applications
- **Specialized Domains**: Graphics, cryptography, bioinformatics

### Next Steps

Congratulations! You've completed the comprehensive Data Structures and Algorithms guide. You now have:

1. **Solid Foundation**: Basic to advanced concepts
2. **Practical Skills**: Implementation experience
3. **Problem-Solving Ability**: Multiple paradigms
4. **Optimization Knowledge**: Efficiency techniques

**Continue Learning:**
- Practice on competitive programming platforms
- Implement algorithms in different languages
- Study algorithm analysis and complexity theory
- Explore specialized algorithm domains
- Contribute to open-source algorithm libraries

### Additional Resources

- **Books**: "Introduction to Algorithms" by CLRS, "Algorithm Design" by Kleinberg & Tardos
- **Online Judges**: Codeforces, AtCoder, TopCoder, LeetCode
- **Courses**: MIT 6.006, Stanford CS161, Princeton Algorithms
- **Competitions**: ACM ICPC, Google Code Jam, Facebook Hacker Cup
- **Research**: Algorithm conferences (SODA, STOC, FOCS)

---

**End of Part 10: Advanced Algorithms**

*This concludes our comprehensive journey through Data Structures and Algorithms. You now have the knowledge and tools to tackle complex algorithmic challenges and build efficient software systems.*