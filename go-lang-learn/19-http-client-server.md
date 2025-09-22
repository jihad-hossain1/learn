# 19. HTTP Client and Server

Go provides excellent built-in support for HTTP through the `net/http` package. This chapter covers building HTTP servers, making HTTP requests, and working with web APIs.

## HTTP Server Basics

### Simple HTTP Server

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "About page")
}

func main() {
    // Register handlers
    http.HandleFunc("/", helloHandler)
    http.HandleFunc("/about", aboutHandler)
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### HTTP Handler Interface

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// Custom handler type
type TimeHandler struct {
    format string
}

// Implement http.Handler interface
func (th *TimeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    currentTime := time.Now().Format(th.format)
    fmt.Fprintf(w, "Current time: %s", currentTime)
}

// Handler function
func greetingHandler(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Query().Get("name")
    if name == "" {
        name = "Guest"
    }
    fmt.Fprintf(w, "Hello, %s!", name)
}

func main() {
    // Using custom handler
    timeHandler := &TimeHandler{format: "2006-01-02 15:04:05"}
    http.Handle("/time", timeHandler)
    
    // Using handler function
    http.HandleFunc("/greeting", greetingHandler)
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: http://localhost:8080/time")
    fmt.Println("Try: http://localhost:8080/greeting?name=John")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### HTTP Methods and Routing

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "strings"
)

type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
}

// In-memory storage
var users = []User{
    {ID: 1, Name: "John Doe", Email: "john@example.com"},
    {ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
}
var nextID = 3

func usersHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        getUsersHandler(w, r)
    case http.MethodPost:
        createUserHandler(w, r)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func userHandler(w http.ResponseWriter, r *http.Request) {
    // Extract user ID from URL path
    path := strings.TrimPrefix(r.URL.Path, "/users/")
    if path == "" {
        http.Error(w, "User ID required", http.StatusBadRequest)
        return
    }
    
    userID, err := strconv.Atoi(path)
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    switch r.Method {
    case http.MethodGet:
        getUserHandler(w, r, userID)
    case http.MethodPut:
        updateUserHandler(w, r, userID)
    case http.MethodDelete:
        deleteUserHandler(w, r, userID)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getUsersHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user.ID = nextID
    nextID++
    users = append(users, user)
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func getUserHandler(w http.ResponseWriter, r *http.Request, userID int) {
    for _, user := range users {
        if user.ID == userID {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(user)
            return
        }
    }
    
    http.Error(w, "User not found", http.StatusNotFound)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request, userID int) {
    for i, user := range users {
        if user.ID == userID {
            var updatedUser User
            err := json.NewDecoder(r.Body).Decode(&updatedUser)
            if err != nil {
                http.Error(w, "Invalid JSON", http.StatusBadRequest)
                return
            }
            
            updatedUser.ID = userID
            users[i] = updatedUser
            
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(updatedUser)
            return
        }
    }
    
    http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request, userID int) {
    for i, user := range users {
        if user.ID == userID {
            users = append(users[:i], users[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }
    
    http.Error(w, "User not found", http.StatusNotFound)
}

func main() {
    http.HandleFunc("/users", usersHandler)
    http.HandleFunc("/users/", userHandler)
    
    fmt.Println("Server starting on :8080")
    fmt.Println("API endpoints:")
    fmt.Println("  GET    /users       - List all users")
    fmt.Println("  POST   /users       - Create user")
    fmt.Println("  GET    /users/{id}  - Get user by ID")
    fmt.Println("  PUT    /users/{id}  - Update user")
    fmt.Println("  DELETE /users/{id}  - Delete user")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## HTTP Middleware

### Basic Middleware

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

// Middleware type
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging middleware
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call the next handler
        next(w, r)
        
        // Log the request
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    }
}

// Authentication middleware
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token != "Bearer secret-token" {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        next(w, r)
    }
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

// Chain multiple middlewares
func chainMiddleware(h http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

// Handlers
func publicHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "This is a public endpoint")
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "This is a protected endpoint")
}

func main() {
    // Public endpoint with logging and CORS
    http.HandleFunc("/public", chainMiddleware(
        publicHandler,
        loggingMiddleware,
        corsMiddleware,
    ))
    
    // Protected endpoint with all middlewares
    http.HandleFunc("/protected", chainMiddleware(
        protectedHandler,
        loggingMiddleware,
        corsMiddleware,
        authMiddleware,
    ))
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: curl http://localhost:8080/public")
    fmt.Println("Try: curl -H 'Authorization: Bearer secret-token' http://localhost:8080/protected")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Advanced Middleware with Context

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "time"
)

type contextKey string

const (
    userIDKey    contextKey = "userID"
    requestIDKey contextKey = "requestID"
)

// Request ID middleware
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := fmt.Sprintf("%d", time.Now().UnixNano())
        ctx := context.WithValue(r.Context(), requestIDKey, requestID)
        
        w.Header().Set("X-Request-ID", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// User context middleware
func userContextMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // In a real app, you'd extract this from JWT token or session
        userID := r.Header.Get("X-User-ID")
        if userID == "" {
            userID = "anonymous"
        }
        
        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Timeout middleware
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx, cancel := context.WithTimeout(r.Context(), timeout)
            defer cancel()
            
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// Enhanced logging middleware
func enhancedLoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Wrap ResponseWriter to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        next.ServeHTTP(wrapped, r)
        
        // Extract context values
        requestID := r.Context().Value(requestIDKey)
        userID := r.Context().Value(userIDKey)
        
        log.Printf("[%v] %s %s %s %d %v", 
            requestID, userID, r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
    requestID := r.Context().Value(requestIDKey)
    userID := r.Context().Value(userIDKey)
    
    fmt.Fprintf(w, "Hello! Request ID: %v, User ID: %v", requestID, userID)
}

func slowHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate slow operation
    select {
    case <-time.After(2 * time.Second):
        fmt.Fprintf(w, "Slow operation completed")
    case <-r.Context().Done():
        http.Error(w, "Request timeout", http.StatusRequestTimeout)
    }
}

func main() {
    mux := http.NewServeMux()
    
    mux.HandleFunc("/", homeHandler)
    mux.HandleFunc("/slow", slowHandler)
    
    // Apply middlewares
    handler := requestIDMiddleware(
        userContextMiddleware(
            enhancedLoggingMiddleware(
                timeoutMiddleware(1 * time.Second)(mux),
            ),
        ),
    )
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Try: curl -H 'X-User-ID: john123' http://localhost:8080/")
    fmt.Println("Try: curl http://localhost:8080/slow (will timeout)")
    
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```

## HTTP Client

### Basic HTTP Client

```go
package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

func main() {
    // Simple GET request
    resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Status: %s\n", resp.Status)
    fmt.Printf("Body: %s\n", string(body))
    
    // Custom client with timeout
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    resp2, err := client.Get("https://jsonplaceholder.typicode.com/posts/2")
    if err != nil {
        log.Fatal(err)
    }
    defer resp2.Body.Close()
    
    fmt.Printf("\nSecond request status: %s\n", resp2.Status)
}
```

### Advanced HTTP Client

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "time"
)

type Post struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
    UserID int    `json:"userId"`
}

type APIClient struct {
    baseURL    string
    httpClient *http.Client
    apiKey     string
}

func NewAPIClient(baseURL, apiKey string) *APIClient {
    return &APIClient{
        baseURL: baseURL,
        apiKey:  apiKey,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                MaxIdleConns:        10,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     30 * time.Second,
            },
        },
    }
}

func (c *APIClient) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
    var reqBody *bytes.Buffer
    
    if body != nil {
        jsonData, err := json.Marshal(body)
        if err != nil {
            return nil, err
        }
        reqBody = bytes.NewBuffer(jsonData)
    } else {
        reqBody = bytes.NewBuffer(nil)
    }
    
    req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
    if err != nil {
        return nil, err
    }
    
    // Set headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "Go-HTTP-Client/1.0")
    
    if c.apiKey != "" {
        req.Header.Set("Authorization", "Bearer "+c.apiKey)
    }
    
    return c.httpClient.Do(req)
}

func (c *APIClient) GetPost(id int) (*Post, error) {
    resp, err := c.doRequest("GET", fmt.Sprintf("/posts/%d", id), nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
    }
    
    var post Post
    err = json.NewDecoder(resp.Body).Decode(&post)
    if err != nil {
        return nil, err
    }
    
    return &post, nil
}

func (c *APIClient) CreatePost(post *Post) (*Post, error) {
    resp, err := c.doRequest("POST", "/posts", post)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusCreated {
        body, _ := ioutil.ReadAll(resp.Body)
        return nil, fmt.Errorf("API request failed with status: %s, body: %s", resp.Status, string(body))
    }
    
    var createdPost Post
    err = json.NewDecoder(resp.Body).Decode(&createdPost)
    if err != nil {
        return nil, err
    }
    
    return &createdPost, nil
}

func (c *APIClient) GetPosts() ([]Post, error) {
    resp, err := c.doRequest("GET", "/posts", nil)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
    }
    
    var posts []Post
    err = json.NewDecoder(resp.Body).Decode(&posts)
    if err != nil {
        return nil, err
    }
    
    return posts, nil
}

func main() {
    client := NewAPIClient("https://jsonplaceholder.typicode.com", "")
    
    // Get a single post
    fmt.Println("Getting post 1...")
    post, err := client.GetPost(1)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Post: %+v\n\n", post)
    
    // Create a new post
    fmt.Println("Creating new post...")
    newPost := &Post{
        Title:  "My New Post",
        Body:   "This is the body of my new post",
        UserID: 1,
    }
    
    createdPost, err := client.CreatePost(newPost)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Created post: %+v\n\n", createdPost)
    
    // Get all posts (limited output)
    fmt.Println("Getting all posts (first 5)...")
    posts, err := client.GetPosts()
    if err != nil {
        log.Fatal(err)
    }
    
    for i, p := range posts {
        if i >= 5 {
            break
        }
        fmt.Printf("Post %d: %s\n", p.ID, p.Title)
    }
}
```

## File Upload and Download

### File Upload Server

```go
package main

import (
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "path/filepath"
)

const uploadDir = "uploads"

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    // Parse multipart form (32MB max)
    err := r.ParseMultipartForm(32 << 20)
    if err != nil {
        http.Error(w, "Unable to parse form", http.StatusBadRequest)
        return
    }
    
    file, header, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Unable to get file", http.StatusBadRequest)
        return
    }
    defer file.Close()
    
    // Create uploads directory if it doesn't exist
    err = os.MkdirAll(uploadDir, 0755)
    if err != nil {
        http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
        return
    }
    
    // Create destination file
    filename := filepath.Join(uploadDir, header.Filename)
    dst, err := os.Create(filename)
    if err != nil {
        http.Error(w, "Unable to create file", http.StatusInternalServerError)
        return
    }
    defer dst.Close()
    
    // Copy file content
    _, err = io.Copy(dst, file)
    if err != nil {
        http.Error(w, "Unable to save file", http.StatusInternalServerError)
        return
    }
    
    fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
    filename := r.URL.Query().Get("file")
    if filename == "" {
        http.Error(w, "File parameter required", http.StatusBadRequest)
        return
    }
    
    filePath := filepath.Join(uploadDir, filename)
    
    // Check if file exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        http.Error(w, "File not found", http.StatusNotFound)
        return
    }
    
    // Set headers for file download
    w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
    w.Header().Set("Content-Type", "application/octet-stream")
    
    // Serve file
    http.ServeFile(w, r, filePath)
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
    files, err := os.ReadDir(uploadDir)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Fprintf(w, "No files uploaded yet")
            return
        }
        http.Error(w, "Unable to read directory", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, "<h2>Uploaded Files</h2><ul>")
    
    for _, file := range files {
        if !file.IsDir() {
            fmt.Fprintf(w, `<li><a href="/download?file=%s">%s</a></li>`, file.Name(), file.Name())
        }
    }
    
    fmt.Fprintf(w, "</ul>")
}

func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>File Upload</title>
    </head>
    <body>
        <h2>Upload File</h2>
        <form action="/upload" method="post" enctype="multipart/form-data">
            <input type="file" name="file" required>
            <input type="submit" value="Upload">
        </form>
        <br>
        <a href="/files">View uploaded files</a>
    </body>
    </html>
    `
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, html)
}

func main() {
    http.HandleFunc("/", uploadFormHandler)
    http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/download", downloadHandler)
    http.HandleFunc("/files", listFilesHandler)
    
    fmt.Println("Server starting on :8080")
    fmt.Println("Visit http://localhost:8080 to upload files")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### File Upload Client

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "os"
    "path/filepath"
)

func uploadFile(url, filename string) error {
    // Open file
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Create multipart form
    var buffer bytes.Buffer
    writer := multipart.NewWriter(&buffer)
    
    // Create form file field
    part, err := writer.CreateFormFile("file", filepath.Base(filename))
    if err != nil {
        return err
    }
    
    // Copy file content to form
    _, err = io.Copy(part, file)
    if err != nil {
        return err
    }
    
    // Close writer to finalize form
    err = writer.Close()
    if err != nil {
        return err
    }
    
    // Create request
    req, err := http.NewRequest("POST", url, &buffer)
    if err != nil {
        return err
    }
    
    // Set content type
    req.Header.Set("Content-Type", writer.FormDataContentType())
    
    // Send request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    fmt.Printf("Response: %s\n", string(body))
    return nil
}

func downloadFile(url, filename string) error {
    // Create request
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("download failed with status: %s", resp.Status)
    }
    
    // Create file
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Copy response body to file
    _, err = io.Copy(file, resp.Body)
    if err != nil {
        return err
    }
    
    fmt.Printf("File downloaded: %s\n", filename)
    return nil
}

func main() {
    // Create a test file
    testFile := "test.txt"
    content := "This is a test file for upload/download demo."
    
    err := os.WriteFile(testFile, []byte(content), 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    // Upload file
    fmt.Println("Uploading file...")
    err = uploadFile("http://localhost:8080/upload", testFile)
    if err != nil {
        log.Printf("Upload failed: %v", err)
    }
    
    // Download file
    fmt.Println("Downloading file...")
    err = downloadFile("http://localhost:8080/download?file=test.txt", "downloaded_test.txt")
    if err != nil {
        log.Printf("Download failed: %v", err)
    }
    
    // Clean up
    os.Remove(testFile)
    os.Remove("downloaded_test.txt")
}
```

## WebSocket Support

### WebSocket Server

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all origins in development
    },
}

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
    id   string
}

type Message struct {
    Type    string `json:"type"`
    Content string `json:"content"`
    From    string `json:"from"`
    Time    string `json:"time"`
}

func newHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            
            log.Printf("Client %s connected. Total clients: %d", client.id, len(h.clients))
            
        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mutex.Unlock()
            
            log.Printf("Client %s disconnected. Total clients: %d", client.id, len(h.clients))
            
        case message := <-h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()
    
    c.conn.SetReadLimit(512)
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })
    
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        // Broadcast message to all clients
        c.hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()
    
    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            
            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }
            
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func wsHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket upgrade error: %v", err)
        return
    }
    
    clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())
    client := &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
        id:   clientID,
    }
    
    client.hub.register <- client
    
    go client.writePump()
    go client.readPump()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    html := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>WebSocket Chat</title>
    </head>
    <body>
        <div id="messages"></div>
        <input type="text" id="messageInput" placeholder="Type a message...">
        <button onclick="sendMessage()">Send</button>
        
        <script>
            const ws = new WebSocket('ws://localhost:8080/ws');
            const messages = document.getElementById('messages');
            const messageInput = document.getElementById('messageInput');
            
            ws.onmessage = function(event) {
                const message = JSON.parse(event.data);
                const div = document.createElement('div');
                div.textContent = message.from + ': ' + message.content + ' (' + message.time + ')';
                messages.appendChild(div);
            };
            
            function sendMessage() {
                const content = messageInput.value;
                if (content) {
                    const message = {
                        type: 'chat',
                        content: content,
                        from: 'User',
                        time: new Date().toLocaleTimeString()
                    };
                    ws.send(JSON.stringify(message));
                    messageInput.value = '';
                }
            }
            
            messageInput.addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    sendMessage();
                }
            });
        </script>
    </body>
    </html>
    `
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprintf(w, html)
}

func main() {
    hub := newHub()
    go hub.run()
    
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        wsHandler(hub, w, r)
    })
    
    fmt.Println("WebSocket server starting on :8080")
    fmt.Println("Visit http://localhost:8080 to test the chat")
    
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Note: This example requires the Gorilla WebSocket package:
```bash
go mod init websocket-example
go get github.com/gorilla/websocket
```

## Best Practices

### 1. Error Handling

```go
// Good: Proper error handling
func apiHandler(w http.ResponseWriter, r *http.Request) {
    var data RequestData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    result, err := processData(data)
    if err != nil {
        log.Printf("Processing error: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
```

### 2. Timeouts and Context

```go
// Good: Use timeouts and context
func makeAPICall(ctx context.Context, url string) (*http.Response, error) {
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    
    return client.Do(req)
}
```

### 3. Resource Management

```go
// Good: Always close response bodies
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close() // Important!

body, err := ioutil.ReadAll(resp.Body)
```

### 4. Security Headers

```go
func securityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000")
        
        next.ServeHTTP(w, r)
    })
}
```

## Exercises

### Exercise 1: REST API
Build a complete REST API for a todo application with:
- CRUD operations for todos
- User authentication
- Input validation
- Error handling
- Middleware for logging and CORS

### Exercise 2: File Server
Create a file server that:
- Supports file upload with size limits
- Provides file listing with metadata
- Implements file download with resume support
- Includes basic authentication

### Exercise 3: HTTP Client Library
Develop an HTTP client library that:
- Supports retries with exponential backoff
- Implements request/response logging
- Handles different authentication methods
- Provides rate limiting

## Key Takeaways

- Go's `net/http` package provides powerful HTTP server and client capabilities
- Use middleware for cross-cutting concerns like logging, authentication, and CORS
- Always handle errors appropriately and return proper HTTP status codes
- Use context for timeouts and cancellation
- Implement proper resource management (close response bodies)
- Follow REST conventions for API design
- Use JSON for data exchange in modern web APIs
- Implement security best practices (HTTPS, security headers, input validation)

## Next Steps

Next, we'll explore [Goroutines](20-goroutines.md) to learn about Go's powerful concurrency model for building high-performance applications.