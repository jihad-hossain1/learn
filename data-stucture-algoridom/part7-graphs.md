# Part 7: Graphs

## Table of Contents
1. [Introduction to Graphs](#introduction-graphs)
2. [Graph Terminology](#graph-terminology)
3. [Graph Representations](#graph-representations)
4. [Graph Traversal Algorithms](#graph-traversal)
5. [Shortest Path Algorithms](#shortest-path)
6. [Minimum Spanning Tree](#minimum-spanning-tree)
7. [Topological Sorting](#topological-sorting)
8. [Graph Cycle Detection](#cycle-detection)
9. [Advanced Graph Algorithms](#advanced-algorithms)
10. [Practice Exercises](#practice-exercises)

---

## 1. Introduction to Graphs {#introduction-graphs}

### What is a Graph?
A graph is a non-linear data structure consisting of vertices (nodes) connected by edges.

### Graph Components
- **Vertices (V)**: Set of nodes
- **Edges (E)**: Set of connections between vertices
- **Graph G = (V, E)**: Mathematical representation

### Why Use Graphs?
- **Network Modeling**: Social networks, computer networks
- **Path Finding**: GPS navigation, game AI
- **Dependency Modeling**: Task scheduling, compilation
- **Optimization**: Resource allocation, circuit design

### Real-World Applications
- **Social Networks**: Facebook friends, Twitter followers
- **Transportation**: Road networks, flight routes
- **Internet**: Web pages and links
- **Biology**: Protein interactions, gene networks
- **Computer Science**: Compilers, operating systems

---

## 2. Graph Terminology {#graph-terminology}

### Basic Terms
- **Vertex/Node**: Individual point in graph
- **Edge**: Connection between two vertices
- **Adjacent**: Two vertices connected by an edge
- **Incident**: Edge connected to a vertex
- **Degree**: Number of edges connected to a vertex
- **Path**: Sequence of vertices connected by edges
- **Cycle**: Path that starts and ends at same vertex

### Graph Types

#### By Direction
- **Undirected Graph**: Edges have no direction
- **Directed Graph (Digraph)**: Edges have direction
- **Mixed Graph**: Contains both directed and undirected edges

#### By Weights
- **Unweighted Graph**: All edges have same weight
- **Weighted Graph**: Edges have associated weights/costs

#### By Connectivity
- **Connected Graph**: Path exists between any two vertices
- **Disconnected Graph**: Some vertices not reachable from others
- **Strongly Connected**: In directed graph, path exists between any two vertices in both directions

#### Special Types
- **Simple Graph**: No self-loops or multiple edges
- **Multigraph**: Multiple edges between same vertices allowed
- **Complete Graph**: Every vertex connected to every other vertex
- **Bipartite Graph**: Vertices can be divided into two disjoint sets
- **Tree**: Connected acyclic graph
- **Forest**: Collection of trees

### Graph Properties
- **Order**: Number of vertices |V|
- **Size**: Number of edges |E|
- **Density**: |E| / |V|(|V|-1)/2 for undirected graphs
- **Diameter**: Longest shortest path between any two vertices

---

## 3. Graph Representations {#graph-representations}

### Adjacency Matrix

```cpp
class GraphMatrix {
public:
    int vertices;
    vector<vector<int>> adjMatrix;
    bool isDirected;
    
    GraphMatrix(int v, bool directed = false) : vertices(v), isDirected(directed) {
        adjMatrix.resize(v, vector<int>(v, 0));
    }
    
    // Add edge
    void addEdge(int src, int dest, int weight = 1) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            adjMatrix[src][dest] = weight;
            
            if (!isDirected) {
                adjMatrix[dest][src] = weight;
            }
        }
    }
    
    // Remove edge
    void removeEdge(int src, int dest) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            adjMatrix[src][dest] = 0;
            
            if (!isDirected) {
                adjMatrix[dest][src] = 0;
            }
        }
    }
    
    // Check if edge exists
    bool hasEdge(int src, int dest) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            return adjMatrix[src][dest] != 0;
        }
        return false;
    }
    
    // Get edge weight
    int getWeight(int src, int dest) {
        if (hasEdge(src, dest)) {
            return adjMatrix[src][dest];
        }
        return 0;
    }
    
    // Get all neighbors of a vertex
    vector<int> getNeighbors(int vertex) {
        vector<int> neighbors;
        if (vertex >= 0 && vertex < vertices) {
            for (int i = 0; i < vertices; i++) {
                if (adjMatrix[vertex][i] != 0) {
                    neighbors.push_back(i);
                }
            }
        }
        return neighbors;
    }
    
    // Get degree of vertex
    int getDegree(int vertex) {
        if (vertex < 0 || vertex >= vertices) return -1;
        
        int degree = 0;
        for (int i = 0; i < vertices; i++) {
            if (adjMatrix[vertex][i] != 0) {
                degree++;
            }
        }
        
        if (!isDirected) {
            return degree;
        } else {
            // For directed graphs, return out-degree
            return degree;
        }
    }
    
    // Get in-degree for directed graphs
    int getInDegree(int vertex) {
        if (vertex < 0 || vertex >= vertices || !isDirected) return -1;
        
        int inDegree = 0;
        for (int i = 0; i < vertices; i++) {
            if (adjMatrix[i][vertex] != 0) {
                inDegree++;
            }
        }
        return inDegree;
    }
    
    // Print graph
    void printGraph() {
        cout << "Adjacency Matrix:" << endl;
        cout << "   ";
        for (int i = 0; i < vertices; i++) {
            cout << setw(3) << i;
        }
        cout << endl;
        
        for (int i = 0; i < vertices; i++) {
            cout << setw(2) << i << " ";
            for (int j = 0; j < vertices; j++) {
                cout << setw(3) << adjMatrix[i][j];
            }
            cout << endl;
        }
    }
};
```

### Adjacency List

```cpp
struct Edge {
    int dest;
    int weight;
    
    Edge(int d, int w = 1) : dest(d), weight(w) {}
};

class GraphList {
public:
    int vertices;
    vector<vector<Edge>> adjList;
    bool isDirected;
    
    GraphList(int v, bool directed = false) : vertices(v), isDirected(directed) {
        adjList.resize(v);
    }
    
    // Add edge
    void addEdge(int src, int dest, int weight = 1) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            adjList[src].push_back(Edge(dest, weight));
            
            if (!isDirected) {
                adjList[dest].push_back(Edge(src, weight));
            }
        }
    }
    
    // Remove edge
    void removeEdge(int src, int dest) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            adjList[src].erase(
                remove_if(adjList[src].begin(), adjList[src].end(),
                         [dest](const Edge& e) { return e.dest == dest; }),
                adjList[src].end());
            
            if (!isDirected) {
                adjList[dest].erase(
                    remove_if(adjList[dest].begin(), adjList[dest].end(),
                             [src](const Edge& e) { return e.dest == src; }),
                    adjList[dest].end());
            }
        }
    }
    
    // Check if edge exists
    bool hasEdge(int src, int dest) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            for (const Edge& edge : adjList[src]) {
                if (edge.dest == dest) {
                    return true;
                }
            }
        }
        return false;
    }
    
    // Get edge weight
    int getWeight(int src, int dest) {
        if (src >= 0 && src < vertices && dest >= 0 && dest < vertices) {
            for (const Edge& edge : adjList[src]) {
                if (edge.dest == dest) {
                    return edge.weight;
                }
            }
        }
        return 0;
    }
    
    // Get all neighbors of a vertex
    vector<int> getNeighbors(int vertex) {
        vector<int> neighbors;
        if (vertex >= 0 && vertex < vertices) {
            for (const Edge& edge : adjList[vertex]) {
                neighbors.push_back(edge.dest);
            }
        }
        return neighbors;
    }
    
    // Get degree of vertex
    int getDegree(int vertex) {
        if (vertex >= 0 && vertex < vertices) {
            return adjList[vertex].size();
        }
        return -1;
    }
    
    // Get in-degree for directed graphs
    int getInDegree(int vertex) {
        if (vertex < 0 || vertex >= vertices || !isDirected) return -1;
        
        int inDegree = 0;
        for (int i = 0; i < vertices; i++) {
            for (const Edge& edge : adjList[i]) {
                if (edge.dest == vertex) {
                    inDegree++;
                }
            }
        }
        return inDegree;
    }
    
    // Print graph
    void printGraph() {
        cout << "Adjacency List:" << endl;
        for (int i = 0; i < vertices; i++) {
            cout << "Vertex " << i << ": ";
            for (const Edge& edge : adjList[i]) {
                cout << "(" << edge.dest << ", " << edge.weight << ") ";
            }
            cout << endl;
        }
    }
};
```

### Edge List

```cpp
struct GraphEdge {
    int src, dest, weight;
    
    GraphEdge(int s, int d, int w = 1) : src(s), dest(d), weight(w) {}
    
    bool operator<(const GraphEdge& other) const {
        return weight < other.weight;
    }
};

class GraphEdgeList {
public:
    int vertices;
    vector<GraphEdge> edges;
    bool isDirected;
    
    GraphEdgeList(int v, bool directed = false) : vertices(v), isDirected(directed) {}
    
    void addEdge(int src, int dest, int weight = 1) {
        edges.push_back(GraphEdge(src, dest, weight));
        if (!isDirected) {
            edges.push_back(GraphEdge(dest, src, weight));
        }
    }
    
    void printGraph() {
        cout << "Edge List:" << endl;
        for (const GraphEdge& edge : edges) {
            cout << edge.src << " -> " << edge.dest << " (" << edge.weight << ")" << endl;
        }
    }
};
```

### Comparison of Representations

| Operation | Adjacency Matrix | Adjacency List | Edge List |
|-----------|------------------|----------------|-----------|
| Space | O(V²) | O(V + E) | O(E) |
| Add Edge | O(1) | O(1) | O(1) |
| Remove Edge | O(1) | O(V) | O(E) |
| Check Edge | O(1) | O(V) | O(E) |
| Get Neighbors | O(V) | O(degree) | O(E) |
| Best For | Dense graphs | Sparse graphs | Algorithms like Kruskal's |

---

## 4. Graph Traversal Algorithms {#graph-traversal}

### Depth-First Search (DFS)

```cpp
class GraphTraversal {
public:
    // DFS Recursive
    static void dfsRecursive(GraphList& graph, int start, vector<bool>& visited) {
        visited[start] = true;
        cout << start << " ";
        
        for (const Edge& edge : graph.adjList[start]) {
            if (!visited[edge.dest]) {
                dfsRecursive(graph, edge.dest, visited);
            }
        }
    }
    
    // DFS Iterative
    static void dfsIterative(GraphList& graph, int start) {
        vector<bool> visited(graph.vertices, false);
        stack<int> stk;
        
        stk.push(start);
        
        while (!stk.empty()) {
            int current = stk.top();
            stk.pop();
            
            if (!visited[current]) {
                visited[current] = true;
                cout << current << " ";
                
                // Add neighbors to stack (in reverse order for consistent traversal)
                for (int i = graph.adjList[current].size() - 1; i >= 0; i--) {
                    int neighbor = graph.adjList[current][i].dest;
                    if (!visited[neighbor]) {
                        stk.push(neighbor);
                    }
                }
            }
        }
    }
    
    // DFS to find all connected components
    static vector<vector<int>> findConnectedComponents(GraphList& graph) {
        vector<bool> visited(graph.vertices, false);
        vector<vector<int>> components;
        
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                vector<int> component;
                dfsComponent(graph, i, visited, component);
                components.push_back(component);
            }
        }
        
        return components;
    }
    
    // DFS to check if path exists between two vertices
    static bool hasPath(GraphList& graph, int start, int end) {
        if (start == end) return true;
        
        vector<bool> visited(graph.vertices, false);
        return dfsPath(graph, start, end, visited);
    }
    
    // DFS to find all paths between two vertices
    static vector<vector<int>> findAllPaths(GraphList& graph, int start, int end) {
        vector<vector<int>> allPaths;
        vector<int> currentPath;
        vector<bool> visited(graph.vertices, false);
        
        dfsAllPaths(graph, start, end, visited, currentPath, allPaths);
        return allPaths;
    }
    
private:
    static void dfsComponent(GraphList& graph, int vertex, vector<bool>& visited, vector<int>& component) {
        visited[vertex] = true;
        component.push_back(vertex);
        
        for (const Edge& edge : graph.adjList[vertex]) {
            if (!visited[edge.dest]) {
                dfsComponent(graph, edge.dest, visited, component);
            }
        }
    }
    
    static bool dfsPath(GraphList& graph, int current, int end, vector<bool>& visited) {
        if (current == end) return true;
        
        visited[current] = true;
        
        for (const Edge& edge : graph.adjList[current]) {
            if (!visited[edge.dest]) {
                if (dfsPath(graph, edge.dest, end, visited)) {
                    return true;
                }
            }
        }
        
        return false;
    }
    
    static void dfsAllPaths(GraphList& graph, int current, int end, vector<bool>& visited,
                           vector<int>& currentPath, vector<vector<int>>& allPaths) {
        visited[current] = true;
        currentPath.push_back(current);
        
        if (current == end) {
            allPaths.push_back(currentPath);
        } else {
            for (const Edge& edge : graph.adjList[current]) {
                if (!visited[edge.dest]) {
                    dfsAllPaths(graph, edge.dest, end, visited, currentPath, allPaths);
                }
            }
        }
        
        // Backtrack
        currentPath.pop_back();
        visited[current] = false;
    }
};
```

### Breadth-First Search (BFS)

```cpp
class BFSTraversal {
public:
    // Standard BFS
    static void bfs(GraphList& graph, int start) {
        vector<bool> visited(graph.vertices, false);
        queue<int> q;
        
        visited[start] = true;
        q.push(start);
        
        while (!q.empty()) {
            int current = q.front();
            q.pop();
            cout << current << " ";
            
            for (const Edge& edge : graph.adjList[current]) {
                if (!visited[edge.dest]) {
                    visited[edge.dest] = true;
                    q.push(edge.dest);
                }
            }
        }
    }
    
    // BFS with level information
    static void bfsWithLevels(GraphList& graph, int start) {
        vector<bool> visited(graph.vertices, false);
        vector<int> level(graph.vertices, -1);
        queue<int> q;
        
        visited[start] = true;
        level[start] = 0;
        q.push(start);
        
        while (!q.empty()) {
            int current = q.front();
            q.pop();
            cout << "Vertex " << current << " at level " << level[current] << endl;
            
            for (const Edge& edge : graph.adjList[current]) {
                if (!visited[edge.dest]) {
                    visited[edge.dest] = true;
                    level[edge.dest] = level[current] + 1;
                    q.push(edge.dest);
                }
            }
        }
    }
    
    // BFS to find shortest path (unweighted)
    static vector<int> shortestPath(GraphList& graph, int start, int end) {
        vector<bool> visited(graph.vertices, false);
        vector<int> parent(graph.vertices, -1);
        queue<int> q;
        
        visited[start] = true;
        q.push(start);
        
        while (!q.empty()) {
            int current = q.front();
            q.pop();
            
            if (current == end) {
                // Reconstruct path
                vector<int> path;
                int node = end;
                while (node != -1) {
                    path.push_back(node);
                    node = parent[node];
                }
                reverse(path.begin(), path.end());
                return path;
            }
            
            for (const Edge& edge : graph.adjList[current]) {
                if (!visited[edge.dest]) {
                    visited[edge.dest] = true;
                    parent[edge.dest] = current;
                    q.push(edge.dest);
                }
            }
        }
        
        return {}; // No path found
    }
    
    // BFS to find shortest distances from source
    static vector<int> shortestDistances(GraphList& graph, int start) {
        vector<int> distance(graph.vertices, -1);
        queue<int> q;
        
        distance[start] = 0;
        q.push(start);
        
        while (!q.empty()) {
            int current = q.front();
            q.pop();
            
            for (const Edge& edge : graph.adjList[current]) {
                if (distance[edge.dest] == -1) {
                    distance[edge.dest] = distance[current] + 1;
                    q.push(edge.dest);
                }
            }
        }
        
        return distance;
    }
    
    // Multi-source BFS
    static vector<int> multiSourceBFS(GraphList& graph, vector<int>& sources) {
        vector<int> distance(graph.vertices, -1);
        queue<int> q;
        
        // Initialize all sources
        for (int source : sources) {
            distance[source] = 0;
            q.push(source);
        }
        
        while (!q.empty()) {
            int current = q.front();
            q.pop();
            
            for (const Edge& edge : graph.adjList[current]) {
                if (distance[edge.dest] == -1) {
                    distance[edge.dest] = distance[current] + 1;
                    q.push(edge.dest);
                }
            }
        }
        
        return distance;
    }
};
```

---

## 5. Shortest Path Algorithms {#shortest-path}

### Dijkstra's Algorithm

```cpp
class ShortestPath {
public:
    // Dijkstra's algorithm for single source shortest path
    static vector<int> dijkstra(GraphList& graph, int source) {
        vector<int> distance(graph.vertices, INT_MAX);
        vector<bool> visited(graph.vertices, false);
        priority_queue<pair<int, int>, vector<pair<int, int>>, greater<pair<int, int>>> pq;
        
        distance[source] = 0;
        pq.push({0, source});
        
        while (!pq.empty()) {
            int u = pq.top().second;
            pq.pop();
            
            if (visited[u]) continue;
            visited[u] = true;
            
            for (const Edge& edge : graph.adjList[u]) {
                int v = edge.dest;
                int weight = edge.weight;
                
                if (!visited[v] && distance[u] + weight < distance[v]) {
                    distance[v] = distance[u] + weight;
                    pq.push({distance[v], v});
                }
            }
        }
        
        return distance;
    }
    
    // Dijkstra with path reconstruction
    static pair<vector<int>, vector<int>> dijkstraWithPath(GraphList& graph, int source, int target) {
        vector<int> distance(graph.vertices, INT_MAX);
        vector<int> parent(graph.vertices, -1);
        vector<bool> visited(graph.vertices, false);
        priority_queue<pair<int, int>, vector<pair<int, int>>, greater<pair<int, int>>> pq;
        
        distance[source] = 0;
        pq.push({0, source});
        
        while (!pq.empty()) {
            int u = pq.top().second;
            pq.pop();
            
            if (visited[u]) continue;
            visited[u] = true;
            
            if (u == target) break; // Early termination
            
            for (const Edge& edge : graph.adjList[u]) {
                int v = edge.dest;
                int weight = edge.weight;
                
                if (!visited[v] && distance[u] + weight < distance[v]) {
                    distance[v] = distance[u] + weight;
                    parent[v] = u;
                    pq.push({distance[v], v});
                }
            }
        }
        
        // Reconstruct path
        vector<int> path;
        if (distance[target] != INT_MAX) {
            int current = target;
            while (current != -1) {
                path.push_back(current);
                current = parent[current];
            }
            reverse(path.begin(), path.end());
        }
        
        return {distance, path};
    }
};
```

### Bellman-Ford Algorithm

```cpp
class BellmanFord {
public:
    // Bellman-Ford algorithm (handles negative weights)
    static pair<vector<int>, bool> bellmanFord(GraphEdgeList& graph, int source) {
        vector<int> distance(graph.vertices, INT_MAX);
        distance[source] = 0;
        
        // Relax edges V-1 times
        for (int i = 0; i < graph.vertices - 1; i++) {
            for (const GraphEdge& edge : graph.edges) {
                if (distance[edge.src] != INT_MAX && 
                    distance[edge.src] + edge.weight < distance[edge.dest]) {
                    distance[edge.dest] = distance[edge.src] + edge.weight;
                }
            }
        }
        
        // Check for negative cycles
        bool hasNegativeCycle = false;
        for (const GraphEdge& edge : graph.edges) {
            if (distance[edge.src] != INT_MAX && 
                distance[edge.src] + edge.weight < distance[edge.dest]) {
                hasNegativeCycle = true;
                break;
            }
        }
        
        return {distance, hasNegativeCycle};
    }
    
    // Bellman-Ford with path reconstruction
    static pair<vector<int>, vector<vector<int>>> bellmanFordWithPaths(GraphEdgeList& graph, int source) {
        vector<int> distance(graph.vertices, INT_MAX);
        vector<int> parent(graph.vertices, -1);
        distance[source] = 0;
        
        // Relax edges V-1 times
        for (int i = 0; i < graph.vertices - 1; i++) {
            for (const GraphEdge& edge : graph.edges) {
                if (distance[edge.src] != INT_MAX && 
                    distance[edge.src] + edge.weight < distance[edge.dest]) {
                    distance[edge.dest] = distance[edge.src] + edge.weight;
                    parent[edge.dest] = edge.src;
                }
            }
        }
        
        // Reconstruct all paths
        vector<vector<int>> paths(graph.vertices);
        for (int i = 0; i < graph.vertices; i++) {
            if (distance[i] != INT_MAX) {
                vector<int> path;
                int current = i;
                while (current != -1) {
                    path.push_back(current);
                    current = parent[current];
                }
                reverse(path.begin(), path.end());
                paths[i] = path;
            }
        }
        
        return {distance, paths};
    }
};
```

### Floyd-Warshall Algorithm

```cpp
class FloydWarshall {
public:
    // Floyd-Warshall algorithm (all pairs shortest path)
    static vector<vector<int>> floydWarshall(GraphMatrix& graph) {
        int V = graph.vertices;
        vector<vector<int>> distance(V, vector<int>(V));
        
        // Initialize distance matrix
        for (int i = 0; i < V; i++) {
            for (int j = 0; j < V; j++) {
                if (i == j) {
                    distance[i][j] = 0;
                } else if (graph.adjMatrix[i][j] != 0) {
                    distance[i][j] = graph.adjMatrix[i][j];
                } else {
                    distance[i][j] = INT_MAX;
                }
            }
        }
        
        // Floyd-Warshall algorithm
        for (int k = 0; k < V; k++) {
            for (int i = 0; i < V; i++) {
                for (int j = 0; j < V; j++) {
                    if (distance[i][k] != INT_MAX && distance[k][j] != INT_MAX &&
                        distance[i][k] + distance[k][j] < distance[i][j]) {
                        distance[i][j] = distance[i][k] + distance[k][j];
                    }
                }
            }
        }
        
        return distance;
    }
    
    // Floyd-Warshall with path reconstruction
    static pair<vector<vector<int>>, vector<vector<int>>> floydWarshallWithPaths(GraphMatrix& graph) {
        int V = graph.vertices;
        vector<vector<int>> distance(V, vector<int>(V));
        vector<vector<int>> next(V, vector<int>(V, -1));
        
        // Initialize
        for (int i = 0; i < V; i++) {
            for (int j = 0; j < V; j++) {
                if (i == j) {
                    distance[i][j] = 0;
                } else if (graph.adjMatrix[i][j] != 0) {
                    distance[i][j] = graph.adjMatrix[i][j];
                    next[i][j] = j;
                } else {
                    distance[i][j] = INT_MAX;
                }
            }
        }
        
        // Floyd-Warshall
        for (int k = 0; k < V; k++) {
            for (int i = 0; i < V; i++) {
                for (int j = 0; j < V; j++) {
                    if (distance[i][k] != INT_MAX && distance[k][j] != INT_MAX &&
                        distance[i][k] + distance[k][j] < distance[i][j]) {
                        distance[i][j] = distance[i][k] + distance[k][j];
                        next[i][j] = next[i][k];
                    }
                }
            }
        }
        
        return {distance, next};
    }
    
    // Reconstruct path using next matrix
    static vector<int> reconstructPath(const vector<vector<int>>& next, int start, int end) {
        if (next[start][end] == -1) {
            return {}; // No path
        }
        
        vector<int> path;
        int current = start;
        path.push_back(current);
        
        while (current != end) {
            current = next[current][end];
            path.push_back(current);
        }
        
        return path;
    }
};
```

---

## 6. Minimum Spanning Tree {#minimum-spanning-tree}

### Kruskal's Algorithm

```cpp
class UnionFind {
public:
    vector<int> parent, rank;
    
    UnionFind(int n) {
        parent.resize(n);
        rank.resize(n, 0);
        for (int i = 0; i < n; i++) {
            parent[i] = i;
        }
    }
    
    int find(int x) {
        if (parent[x] != x) {
            parent[x] = find(parent[x]); // Path compression
        }
        return parent[x];
    }
    
    bool unite(int x, int y) {
        int rootX = find(x);
        int rootY = find(y);
        
        if (rootX == rootY) {
            return false; // Already in same set
        }
        
        // Union by rank
        if (rank[rootX] < rank[rootY]) {
            parent[rootX] = rootY;
        } else if (rank[rootX] > rank[rootY]) {
            parent[rootY] = rootX;
        } else {
            parent[rootY] = rootX;
            rank[rootX]++;
        }
        
        return true;
    }
};

class MST {
public:
    // Kruskal's algorithm
    static pair<vector<GraphEdge>, int> kruskal(GraphEdgeList& graph) {
        vector<GraphEdge> mst;
        int totalWeight = 0;
        
        // Sort edges by weight
        vector<GraphEdge> edges = graph.edges;
        sort(edges.begin(), edges.end());
        
        UnionFind uf(graph.vertices);
        
        for (const GraphEdge& edge : edges) {
            if (uf.unite(edge.src, edge.dest)) {
                mst.push_back(edge);
                totalWeight += edge.weight;
                
                if (mst.size() == graph.vertices - 1) {
                    break; // MST complete
                }
            }
        }
        
        return {mst, totalWeight};
    }
};
```

### Prim's Algorithm

```cpp
class PrimMST {
public:
    // Prim's algorithm
    static pair<vector<GraphEdge>, int> prim(GraphList& graph, int start = 0) {
        vector<GraphEdge> mst;
        vector<bool> inMST(graph.vertices, false);
        priority_queue<pair<int, pair<int, int>>, 
                      vector<pair<int, pair<int, int>>>,
                      greater<pair<int, pair<int, int>>>> pq;
        
        int totalWeight = 0;
        inMST[start] = true;
        
        // Add all edges from start vertex
        for (const Edge& edge : graph.adjList[start]) {
            pq.push({edge.weight, {start, edge.dest}});
        }
        
        while (!pq.empty() && mst.size() < graph.vertices - 1) {
            auto [weight, edge_pair] = pq.top();
            auto [src, dest] = edge_pair;
            pq.pop();
            
            if (inMST[dest]) {
                continue; // Skip if already in MST
            }
            
            // Add edge to MST
            mst.push_back(GraphEdge(src, dest, weight));
            totalWeight += weight;
            inMST[dest] = true;
            
            // Add all edges from new vertex
            for (const Edge& edge : graph.adjList[dest]) {
                if (!inMST[edge.dest]) {
                    pq.push({edge.weight, {dest, edge.dest}});
                }
            }
        }
        
        return {mst, totalWeight};
    }
    
    // Prim's algorithm with adjacency matrix
    static pair<vector<GraphEdge>, int> primMatrix(GraphMatrix& graph) {
        vector<GraphEdge> mst;
        vector<bool> inMST(graph.vertices, false);
        vector<int> key(graph.vertices, INT_MAX);
        vector<int> parent(graph.vertices, -1);
        
        key[0] = 0;
        int totalWeight = 0;
        
        for (int count = 0; count < graph.vertices; count++) {
            // Find minimum key vertex not in MST
            int u = -1;
            for (int v = 0; v < graph.vertices; v++) {
                if (!inMST[v] && (u == -1 || key[v] < key[u])) {
                    u = v;
                }
            }
            
            inMST[u] = true;
            
            if (parent[u] != -1) {
                mst.push_back(GraphEdge(parent[u], u, key[u]));
                totalWeight += key[u];
            }
            
            // Update key values of adjacent vertices
            for (int v = 0; v < graph.vertices; v++) {
                if (graph.adjMatrix[u][v] != 0 && !inMST[v] && 
                    graph.adjMatrix[u][v] < key[v]) {
                    key[v] = graph.adjMatrix[u][v];
                    parent[v] = u;
                }
            }
        }
        
        return {mst, totalWeight};
    }
};
```

---

## 7. Topological Sorting {#topological-sorting}

### DFS-based Topological Sort

```cpp
class TopologicalSort {
public:
    // DFS-based topological sorting
    static vector<int> topologicalSortDFS(GraphList& graph) {
        if (!graph.isDirected) {
            throw invalid_argument("Topological sort only works on directed graphs");
        }
        
        vector<bool> visited(graph.vertices, false);
        stack<int> stk;
        
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                topologicalSortUtil(graph, i, visited, stk);
            }
        }
        
        vector<int> result;
        while (!stk.empty()) {
            result.push_back(stk.top());
            stk.pop();
        }
        
        return result;
    }
    
    // Kahn's algorithm (BFS-based)
    static vector<int> topologicalSortKahn(GraphList& graph) {
        if (!graph.isDirected) {
            throw invalid_argument("Topological sort only works on directed graphs");
        }
        
        vector<int> inDegree(graph.vertices, 0);
        
        // Calculate in-degrees
        for (int i = 0; i < graph.vertices; i++) {
            for (const Edge& edge : graph.adjList[i]) {
                inDegree[edge.dest]++;
            }
        }
        
        queue<int> q;
        for (int i = 0; i < graph.vertices; i++) {
            if (inDegree[i] == 0) {
                q.push(i);
            }
        }
        
        vector<int> result;
        
        while (!q.empty()) {
            int current = q.front();
            q.pop();
            result.push_back(current);
            
            for (const Edge& edge : graph.adjList[current]) {
                inDegree[edge.dest]--;
                if (inDegree[edge.dest] == 0) {
                    q.push(edge.dest);
                }
            }
        }
        
        if (result.size() != graph.vertices) {
            throw runtime_error("Graph contains a cycle - topological sort not possible");
        }
        
        return result;
    }
    
    // Check if topological sort is possible (DAG check)
    static bool isDAG(GraphList& graph) {
        try {
            topologicalSortKahn(graph);
            return true;
        } catch (const runtime_error&) {
            return false;
        }
    }
    
private:
    static void topologicalSortUtil(GraphList& graph, int vertex, 
                                   vector<bool>& visited, stack<int>& stk) {
        visited[vertex] = true;
        
        for (const Edge& edge : graph.adjList[vertex]) {
            if (!visited[edge.dest]) {
                topologicalSortUtil(graph, edge.dest, visited, stk);
            }
        }
        
        stk.push(vertex);
    }
};
```

---

## 8. Graph Cycle Detection {#cycle-detection}

### Cycle Detection in Undirected Graphs

```cpp
class CycleDetection {
public:
    // Cycle detection in undirected graph using DFS
    static bool hasCycleUndirected(GraphList& graph) {
        if (graph.isDirected) {
            throw invalid_argument("Use directed graph cycle detection method");
        }
        
        vector<bool> visited(graph.vertices, false);
        
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                if (hasCycleUndirectedUtil(graph, i, -1, visited)) {
                    return true;
                }
            }
        }
        
        return false;
    }
    
    // Cycle detection in undirected graph using Union-Find
    static bool hasCycleUndirectedUF(GraphEdgeList& graph) {
        UnionFind uf(graph.vertices);
        
        for (const GraphEdge& edge : graph.edges) {
            if (!uf.unite(edge.src, edge.dest)) {
                return true; // Cycle found
            }
        }
        
        return false;
    }
    
    // Cycle detection in directed graph using DFS
    static bool hasCycleDirected(GraphList& graph) {
        if (!graph.isDirected) {
            throw invalid_argument("Use undirected graph cycle detection method");
        }
        
        vector<int> color(graph.vertices, 0); // 0: white, 1: gray, 2: black
        
        for (int i = 0; i < graph.vertices; i++) {
            if (color[i] == 0) {
                if (hasCycleDirectedUtil(graph, i, color)) {
                    return true;
                }
            }
        }
        
        return false;
    }
    
    // Find all cycles in undirected graph
    static vector<vector<int>> findAllCyclesUndirected(GraphList& graph) {
        vector<vector<int>> cycles;
        vector<bool> visited(graph.vertices, false);
        
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                vector<int> path;
                findCyclesUndirectedUtil(graph, i, -1, visited, path, cycles);
            }
        }
        
        return cycles;
    }
    
private:
    static bool hasCycleUndirectedUtil(GraphList& graph, int vertex, int parent, vector<bool>& visited) {
        visited[vertex] = true;
        
        for (const Edge& edge : graph.adjList[vertex]) {
            int neighbor = edge.dest;
            
            if (!visited[neighbor]) {
                if (hasCycleUndirectedUtil(graph, neighbor, vertex, visited)) {
                    return true;
                }
            } else if (neighbor != parent) {
                return true; // Back edge found
            }
        }
        
        return false;
    }
    
    static bool hasCycleDirectedUtil(GraphList& graph, int vertex, vector<int>& color) {
        color[vertex] = 1; // Gray
        
        for (const Edge& edge : graph.adjList[vertex]) {
            int neighbor = edge.dest;
            
            if (color[neighbor] == 1) {
                return true; // Back edge found
            }
            
            if (color[neighbor] == 0 && hasCycleDirectedUtil(graph, neighbor, color)) {
                return true;
            }
        }
        
        color[vertex] = 2; // Black
        return false;
    }
    
    static void findCyclesUndirectedUtil(GraphList& graph, int vertex, int parent,
                                        vector<bool>& visited, vector<int>& path,
                                        vector<vector<int>>& cycles) {
        visited[vertex] = true;
        path.push_back(vertex);
        
        for (const Edge& edge : graph.adjList[vertex]) {
            int neighbor = edge.dest;
            
            if (!visited[neighbor]) {
                findCyclesUndirectedUtil(graph, neighbor, vertex, visited, path, cycles);
            } else if (neighbor != parent) {
                // Found cycle
                vector<int> cycle;
                bool startFound = false;
                for (int node : path) {
                    if (node == neighbor) {
                        startFound = true;
                    }
                    if (startFound) {
                        cycle.push_back(node);
                    }
                }
                cycle.push_back(neighbor);
                cycles.push_back(cycle);
            }
        }
        
        path.pop_back();
    }
};
```

---

## 9. Advanced Graph Algorithms {#advanced-algorithms}

### Strongly Connected Components (Kosaraju's Algorithm)

```cpp
class StronglyConnectedComponents {
public:
    // Kosaraju's algorithm
    static vector<vector<int>> kosaraju(GraphList& graph) {
        if (!graph.isDirected) {
            throw invalid_argument("SCC only applies to directed graphs");
        }
        
        vector<bool> visited(graph.vertices, false);
        stack<int> finishOrder;
        
        // Step 1: Fill vertices in stack according to their finishing times
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                fillOrder(graph, i, visited, finishOrder);
            }
        }
        
        // Step 2: Create transpose graph
        GraphList transpose = getTranspose(graph);
        
        // Step 3: Process vertices in order defined by stack
        fill(visited.begin(), visited.end(), false);
        vector<vector<int>> sccs;
        
        while (!finishOrder.empty()) {
            int vertex = finishOrder.top();
            finishOrder.pop();
            
            if (!visited[vertex]) {
                vector<int> scc;
                dfsUtil(transpose, vertex, visited, scc);
                sccs.push_back(scc);
            }
        }
        
        return sccs;
    }
    
    // Check if graph is strongly connected
    static bool isStronglyConnected(GraphList& graph) {
        vector<vector<int>> sccs = kosaraju(graph);
        return sccs.size() == 1;
    }
    
private:
    static void fillOrder(GraphList& graph, int vertex, vector<bool>& visited, stack<int>& finishOrder) {
        visited[vertex] = true;
        
        for (const Edge& edge : graph.adjList[vertex]) {
            if (!visited[edge.dest]) {
                fillOrder(graph, edge.dest, visited, finishOrder);
            }
        }
        
        finishOrder.push(vertex);
    }
    
    static GraphList getTranspose(GraphList& graph) {
        GraphList transpose(graph.vertices, true);
        
        for (int i = 0; i < graph.vertices; i++) {
            for (const Edge& edge : graph.adjList[i]) {
                transpose.addEdge(edge.dest, i, edge.weight);
            }
        }
        
        return transpose;
    }
    
    static void dfsUtil(GraphList& graph, int vertex, vector<bool>& visited, vector<int>& component) {
        visited[vertex] = true;
        component.push_back(vertex);
        
        for (const Edge& edge : graph.adjList[vertex]) {
            if (!visited[edge.dest]) {
                dfsUtil(graph, edge.dest, visited, component);
            }
        }
    }
};
```

### Articulation Points and Bridges

```cpp
class ArticulationPointsBridges {
public:
    // Find articulation points (cut vertices)
    static vector<int> findArticulationPoints(GraphList& graph) {
        vector<int> articulationPoints;
        vector<bool> visited(graph.vertices, false);
        vector<int> discovery(graph.vertices, -1);
        vector<int> low(graph.vertices, -1);
        vector<int> parent(graph.vertices, -1);
        vector<bool> isAP(graph.vertices, false);
        int time = 0;
        
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                findAPUtil(graph, i, visited, discovery, low, parent, isAP, time);
            }
        }
        
        for (int i = 0; i < graph.vertices; i++) {
            if (isAP[i]) {
                articulationPoints.push_back(i);
            }
        }
        
        return articulationPoints;
    }
    
    // Find bridges (cut edges)
    static vector<pair<int, int>> findBridges(GraphList& graph) {
        vector<pair<int, int>> bridges;
        vector<bool> visited(graph.vertices, false);
        vector<int> discovery(graph.vertices, -1);
        vector<int> low(graph.vertices, -1);
        vector<int> parent(graph.vertices, -1);
        int time = 0;
        
        for (int i = 0; i < graph.vertices; i++) {
            if (!visited[i]) {
                findBridgesUtil(graph, i, visited, discovery, low, parent, bridges, time);
            }
        }
        
        return bridges;
    }
    
private:
    static void findAPUtil(GraphList& graph, int u, vector<bool>& visited,
                          vector<int>& discovery, vector<int>& low,
                          vector<int>& parent, vector<bool>& isAP, int& time) {
        int children = 0;
        visited[u] = true;
        discovery[u] = low[u] = ++time;
        
        for (const Edge& edge : graph.adjList[u]) {
            int v = edge.dest;
            
            if (!visited[v]) {
                children++;
                parent[v] = u;
                findAPUtil(graph, v, visited, discovery, low, parent, isAP, time);
                
                low[u] = min(low[u], low[v]);
                
                // Root is AP if it has more than one child
                if (parent[u] == -1 && children > 1) {
                    isAP[u] = true;
                }
                
                // Non-root is AP if removing it disconnects the graph
                if (parent[u] != -1 && low[v] >= discovery[u]) {
                    isAP[u] = true;
                }
            } else if (v != parent[u]) {
                low[u] = min(low[u], discovery[v]);
            }
        }
    }
    
    static void findBridgesUtil(GraphList& graph, int u, vector<bool>& visited,
                               vector<int>& discovery, vector<int>& low,
                               vector<int>& parent, vector<pair<int, int>>& bridges, int& time) {
        visited[u] = true;
        discovery[u] = low[u] = ++time;
        
        for (const Edge& edge : graph.adjList[u]) {
            int v = edge.dest;
            
            if (!visited[v]) {
                parent[v] = u;
                findBridgesUtil(graph, v, visited, discovery, low, parent, bridges, time);
                
                low[u] = min(low[u], low[v]);
                
                // Bridge condition
                if (low[v] > discovery[u]) {
                    bridges.push_back({u, v});
                }
            } else if (v != parent[u]) {
                low[u] = min(low[u], discovery[v]);
            }
        }
    }
};
```

---

## 10. Practice Exercises {#practice-exercises}

### Basic Graph Exercises
1. **Graph Representation**:
   - Convert between adjacency matrix and adjacency list
   - Implement graph with different representations
   - Calculate graph properties (density, degree distribution)

2. **Graph Traversal**:
   - Implement DFS and BFS (recursive and iterative)
   - Find connected components
   - Check if path exists between two vertices
   - Find all paths between two vertices

3. **Basic Graph Problems**:
   - Count number of islands (2D grid)
   - Find shortest path in unweighted graph
   - Detect cycles in undirected graph
   - Check if graph is bipartite

### Intermediate Graph Exercises
1. **Shortest Path Problems**:
   - Implement Dijkstra's algorithm
   - Find shortest path with exactly k edges
   - Shortest path in grid with obstacles
   - Find cheapest flights with k stops

2. **Tree and MST Problems**:
   - Implement Kruskal's and Prim's algorithms
   - Find second minimum spanning tree
   - Critical connections in network
   - Minimum cost to connect all points

3. **Topological Sorting**:
   - Course scheduling problems
   - Alien dictionary
   - Parallel course scheduling
   - Build order with dependencies

### Advanced Graph Exercises
1. **Advanced Algorithms**:
   - Implement Floyd-Warshall algorithm
   - Find strongly connected components
   - Articulation points and bridges
   - Maximum flow algorithms

2. **Complex Graph Problems**:
   - Traveling salesman problem (TSP)
   - Graph coloring problems
   - Hamiltonian path and cycle
   - Network flow problems

3. **Real-world Applications**:
   - Social network analysis
   - Web crawler implementation
   - GPS navigation system
   - Dependency resolution

### Complete Example Program

```cpp
#include <iostream>
#include <vector>
#include <queue>
#include <stack>
#include <algorithm>
#include <climits>
using namespace std;

int main() {
    cout << "=== Graph Algorithms Demo ===" << endl;
    
    // Create a sample graph
    cout << "\nCreating sample graph..." << endl;
    GraphList graph(6, false); // Undirected graph with 6 vertices
    
    // Add edges
    graph.addEdge(0, 1, 4);
    graph.addEdge(0, 2, 3);
    graph.addEdge(1, 2, 1);
    graph.addEdge(1, 3, 2);
    graph.addEdge(2, 3, 4);
    graph.addEdge(3, 4, 2);
    graph.addEdge(4, 5, 6);
    
    graph.printGraph();
    
    // Graph traversals
    cout << "\nDFS from vertex 0: ";
    vector<bool> visited(graph.vertices, false);
    GraphTraversal::dfsRecursive(graph, 0, visited);
    cout << endl;
    
    cout << "BFS from vertex 0: ";
    BFSTraversal::bfs(graph, 0);
    cout << endl;
    
    // Shortest path
    cout << "\nShortest path from 0 to 5: ";
    vector<int> path = BFSTraversal::shortestPath(graph, 0, 5);
    for (int i = 0; i < path.size(); i++) {
        cout << path[i];
        if (i < path.size() - 1) cout << " -> ";
    }
    cout << endl;
    
    // Dijkstra's algorithm
    cout << "\nDijkstra's shortest distances from vertex 0:" << endl;
    vector<int> distances = ShortestPath::dijkstra(graph, 0);
    for (int i = 0; i < distances.size(); i++) {
        cout << "Distance to " << i << ": " << distances[i] << endl;
    }
    
    // Minimum Spanning Tree
    cout << "\nMinimum Spanning Tree (Prim's):" << endl;
    auto [mstEdges, totalWeight] = PrimMST::prim(graph, 0);
    cout << "Total weight: " << totalWeight << endl;
    cout << "Edges in MST:" << endl;
    for (const GraphEdge& edge : mstEdges) {
        cout << edge.src << " - " << edge.dest << " (" << edge.weight << ")" << endl;
    }
    
    // Cycle detection
    cout << "\nCycle detection: ";
    bool hasCycle = CycleDetection::hasCycleUndirected(graph);
    cout << (hasCycle ? "Cycle found" : "No cycle") << endl;
    
    // Connected components
    cout << "\nConnected components:" << endl;
    vector<vector<int>> components = GraphTraversal::findConnectedComponents(graph);
    for (int i = 0; i < components.size(); i++) {
        cout << "Component " << i + 1 << ": ";
        for (int vertex : components[i]) {
            cout << vertex << " ";
        }
        cout << endl;
    }
    
    // Directed graph example
    cout << "\n=== Directed Graph Example ===" << endl;
    GraphList directedGraph(4, true);
    directedGraph.addEdge(0, 1);
    directedGraph.addEdge(1, 2);
    directedGraph.addEdge(2, 3);
    directedGraph.addEdge(3, 1); // Creates a cycle
    
    directedGraph.printGraph();
    
    // Topological sort
    cout << "\nTopological sort: ";
    try {
        vector<int> topoSort = TopologicalSort::topologicalSortKahn(directedGraph);
        for (int vertex : topoSort) {
            cout << vertex << " ";
        }
        cout << endl;
    } catch (const runtime_error& e) {
        cout << "Error: " << e.what() << endl;
    }
    
    // Cycle detection in directed graph
    cout << "\nCycle in directed graph: ";
    bool hasDirectedCycle = CycleDetection::hasCycleDirected(directedGraph);
    cout << (hasDirectedCycle ? "Cycle found" : "No cycle") << endl;
    
    return 0;
}
```

---

## Summary

### Key Concepts Covered
1. **Graph Fundamentals**: Terminology, types, and properties
2. **Graph Representations**: Adjacency matrix, list, and edge list
3. **Graph Traversal**: DFS and BFS algorithms
4. **Shortest Path**: Dijkstra's, Bellman-Ford, Floyd-Warshall
5. **Minimum Spanning Tree**: Kruskal's and Prim's algorithms
6. **Topological Sorting**: DFS and Kahn's algorithms
7. **Cycle Detection**: For both directed and undirected graphs
8. **Advanced Algorithms**: SCC, articulation points, bridges

### Time and Space Complexities

| Algorithm | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| DFS/BFS | O(V + E) | O(V) |
| Dijkstra's | O((V + E) log V) | O(V) |
| Bellman-Ford | O(VE) | O(V) |
| Floyd-Warshall | O(V³) | O(V²) |
| Kruskal's MST | O(E log E) | O(V) |
| Prim's MST | O((V + E) log V) | O(V) |
| Topological Sort | O(V + E) | O(V) |
| Kosaraju's SCC | O(V + E) | O(V) |

### When to Use Different Algorithms
- **DFS**: Cycle detection, topological sorting, pathfinding
- **BFS**: Shortest path in unweighted graphs, level-order traversal
- **Dijkstra's**: Single-source shortest path with non-negative weights
- **Bellman-Ford**: Single-source shortest path with negative weights
- **Floyd-Warshall**: All-pairs shortest path
- **Kruskal's**: MST for sparse graphs
- **Prim's**: MST for dense graphs

### Next Steps
Proceed to **Part 8: Sorting and Searching Algorithms** to learn about various sorting techniques and advanced searching methods.

### Additional Resources
- [Graph Theory Visualization](https://visualgo.net/en/graphds)
- [Graph Algorithms Practice](https://leetcode.com/tag/graph/)
- [Advanced Graph Theory](https://www.coursera.org/learn/algorithms-graphs-data-structures)
- [Network Analysis](https://networkx.org/)