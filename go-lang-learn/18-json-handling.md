# 18. JSON Handling

JSON (JavaScript Object Notation) is a lightweight data interchange format that's widely used in web APIs and configuration files. Go provides excellent built-in support for JSON through the `encoding/json` package.

## JSON Basics

### What is JSON?

JSON is a text-based data format that's easy for humans to read and write, and easy for machines to parse and generate. It's built on two structures:
- A collection of name/value pairs (similar to a map)
- An ordered list of values (similar to an array)

### JSON Data Types

```json
{
  "string": "Hello, World!",
  "number": 42,
  "float": 3.14,
  "boolean": true,
  "null": null,
  "array": [1, 2, 3],
  "object": {
    "nested": "value"
  }
}
```

## Marshaling (Go to JSON)

### Basic Marshaling

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
    IsAdmin bool   `json:"is_admin"`
}

func main() {
    person := Person{
        Name:    "John Doe",
        Age:     30,
        Email:   "john@example.com",
        IsAdmin: false,
    }
    
    // Marshal to JSON
    jsonData, err := json.Marshal(person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(jsonData))
    // Output: {"name":"John Doe","age":30,"email":"john@example.com","is_admin":false}
}
```

### Pretty Printing JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Product struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Price       float64  `json:"price"`
    Categories  []string `json:"categories"`
    InStock     bool     `json:"in_stock"`
    Description string   `json:"description,omitempty"`
}

func main() {
    product := Product{
        ID:         1,
        Name:       "Laptop",
        Price:      999.99,
        Categories: []string{"Electronics", "Computers"},
        InStock:    true,
    }
    
    // Pretty print JSON
    jsonData, err := json.MarshalIndent(product, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(jsonData))
}
```

### Marshaling Maps and Slices

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

func main() {
    // Marshaling a map
    data := map[string]interface{}{
        "name":    "Alice",
        "age":     25,
        "hobbies": []string{"reading", "swimming"},
        "address": map[string]string{
            "city":    "New York",
            "country": "USA",
        },
    }
    
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(jsonData))
    
    // Marshaling a slice of structs
    type User struct {
        ID   int    `json:"id"`
        Name string `json:"name"`
    }
    
    users := []User{
        {ID: 1, Name: "Alice"},
        {ID: 2, Name: "Bob"},
        {ID: 3, Name: "Charlie"},
    }
    
    usersJSON, err := json.MarshalIndent(users, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\nUsers:")
    fmt.Println(string(usersJSON))
}
```

## Unmarshaling (JSON to Go)

### Basic Unmarshaling

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Person struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Email   string `json:"email"`
    IsAdmin bool   `json:"is_admin"`
}

func main() {
    jsonData := `{
        "name": "Jane Smith",
        "age": 28,
        "email": "jane@example.com",
        "is_admin": true
    }`
    
    var person Person
    err := json.Unmarshal([]byte(jsonData), &person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Name: %s\n", person.Name)
    fmt.Printf("Age: %d\n", person.Age)
    fmt.Printf("Email: %s\n", person.Email)
    fmt.Printf("Is Admin: %t\n", person.IsAdmin)
}
```

### Unmarshaling into Maps

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

func main() {
    jsonData := `{
        "name": "John",
        "age": 30,
        "city": "New York",
        "hobbies": ["reading", "gaming"]
    }`
    
    var data map[string]interface{}
    err := json.Unmarshal([]byte(jsonData), &data)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Name: %s\n", data["name"])
    fmt.Printf("Age: %.0f\n", data["age"].(float64)) // JSON numbers are float64
    fmt.Printf("City: %s\n", data["city"])
    
    // Handle array
    hobbies := data["hobbies"].([]interface{})
    fmt.Print("Hobbies: ")
    for i, hobby := range hobbies {
        if i > 0 {
            fmt.Print(", ")
        }
        fmt.Print(hobby.(string))
    }
    fmt.Println()
}
```

### Unmarshaling Arrays

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func main() {
    jsonData := `[
        {"id": 1, "username": "alice", "email": "alice@example.com"},
        {"id": 2, "username": "bob", "email": "bob@example.com"},
        {"id": 3, "username": "charlie", "email": "charlie@example.com"}
    ]`
    
    var users []User
    err := json.Unmarshal([]byte(jsonData), &users)
    if err != nil {
        log.Fatal(err)
    }
    
    for _, user := range users {
        fmt.Printf("ID: %d, Username: %s, Email: %s\n", 
            user.ID, user.Username, user.Email)
    }
}
```

## JSON Struct Tags

### Common Struct Tags

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
)

type Employee struct {
    // Custom field name
    ID int `json:"employee_id"`
    
    // Omit if empty
    Name string `json:"name,omitempty"`
    
    // Ignore this field
    Password string `json:"-"`
    
    // Include field even if empty
    Department string `json:"department"`
    
    // Custom field name with omitempty
    PhoneNumber string `json:"phone,omitempty"`
    
    // Embedded struct
    Address Address `json:"address"`
    
    // Time field
    CreatedAt time.Time `json:"created_at"`
    
    // Pointer field
    Manager *Employee `json:"manager,omitempty"`
}

type Address struct {
    Street  string `json:"street"`
    City    string `json:"city"`
    Country string `json:"country"`
    ZipCode string `json:"zip_code,omitempty"`
}

func main() {
    emp := Employee{
        ID:         123,
        Name:       "John Doe",
        Password:   "secret123", // This will be ignored
        Department: "",           // This will be included even if empty
        Address: Address{
            Street:  "123 Main St",
            City:    "New York",
            Country: "USA",
        },
        CreatedAt: time.Now(),
    }
    
    jsonData, err := json.MarshalIndent(emp, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(jsonData))
}
```

### Advanced Struct Tags

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Product struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price,string"` // Convert number to string
    
    // Nested JSON as string
    Metadata string `json:"metadata,omitempty"`
    
    // Array with omitempty
    Tags []string `json:"tags,omitempty"`
    
    // Map field
    Attributes map[string]interface{} `json:"attributes,omitempty"`
}

func main() {
    product := Product{
        ID:    1,
        Name:  "Widget",
        Price: 29.99,
        Tags:  []string{"electronics", "gadget"},
        Attributes: map[string]interface{}{
            "color":  "blue",
            "weight": 1.5,
        },
    }
    
    jsonData, err := json.MarshalIndent(product, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(string(jsonData))
}
```

## Custom JSON Marshaling

### Implementing json.Marshaler Interface

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
)

type CustomTime struct {
    time.Time
}

// MarshalJSON implements the json.Marshaler interface
func (ct CustomTime) MarshalJSON() ([]byte, error) {
    // Custom format: "2006-01-02 15:04:05"
    formatted := ct.Time.Format("2006-01-02 15:04:05")
    return json.Marshal(formatted)
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ct *CustomTime) UnmarshalJSON(data []byte) error {
    var timeStr string
    if err := json.Unmarshal(data, &timeStr); err != nil {
        return err
    }
    
    parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
    if err != nil {
        return err
    }
    
    ct.Time = parsedTime
    return nil
}

type Event struct {
    ID        int        `json:"id"`
    Name      string     `json:"name"`
    Timestamp CustomTime `json:"timestamp"`
}

func main() {
    event := Event{
        ID:        1,
        Name:      "Meeting",
        Timestamp: CustomTime{time.Now()},
    }
    
    // Marshal
    jsonData, err := json.MarshalIndent(event, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Marshaled:")
    fmt.Println(string(jsonData))
    
    // Unmarshal
    var parsedEvent Event
    err = json.Unmarshal(jsonData, &parsedEvent)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\nUnmarshaled:")
    fmt.Printf("ID: %d\n", parsedEvent.ID)
    fmt.Printf("Name: %s\n", parsedEvent.Name)
    fmt.Printf("Timestamp: %s\n", parsedEvent.Timestamp.Time.Format(time.RFC3339))
}
```

### Custom Marshaling for Complex Types

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "strconv"
    "strings"
)

type Version struct {
    Major int
    Minor int
    Patch int
}

// MarshalJSON converts Version to "major.minor.patch" format
func (v Version) MarshalJSON() ([]byte, error) {
    versionStr := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
    return json.Marshal(versionStr)
}

// UnmarshalJSON parses "major.minor.patch" format to Version
func (v *Version) UnmarshalJSON(data []byte) error {
    var versionStr string
    if err := json.Unmarshal(data, &versionStr); err != nil {
        return err
    }
    
    parts := strings.Split(versionStr, ".")
    if len(parts) != 3 {
        return fmt.Errorf("invalid version format: %s", versionStr)
    }
    
    major, err := strconv.Atoi(parts[0])
    if err != nil {
        return err
    }
    
    minor, err := strconv.Atoi(parts[1])
    if err != nil {
        return err
    }
    
    patch, err := strconv.Atoi(parts[2])
    if err != nil {
        return err
    }
    
    v.Major = major
    v.Minor = minor
    v.Patch = patch
    
    return nil
}

type Software struct {
    Name    string  `json:"name"`
    Version Version `json:"version"`
}

func main() {
    software := Software{
        Name: "MyApp",
        Version: Version{
            Major: 2,
            Minor: 1,
            Patch: 3,
        },
    }
    
    // Marshal
    jsonData, err := json.MarshalIndent(software, "", "  ")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Marshaled:")
    fmt.Println(string(jsonData))
    
    // Unmarshal
    jsonInput := `{"name": "AnotherApp", "version": "1.5.2"}`
    var parsedSoftware Software
    err = json.Unmarshal([]byte(jsonInput), &parsedSoftware)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\nUnmarshaled:")
    fmt.Printf("Name: %s\n", parsedSoftware.Name)
    fmt.Printf("Version: %d.%d.%d\n", 
        parsedSoftware.Version.Major,
        parsedSoftware.Version.Minor,
        parsedSoftware.Version.Patch)
}
```

## Working with JSON Streams

### JSON Encoder and Decoder

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    "strings"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func main() {
    users := []User{
        {ID: 1, Username: "alice", Email: "alice@example.com"},
        {ID: 2, Username: "bob", Email: "bob@example.com"},
        {ID: 3, Username: "charlie", Email: "charlie@example.com"},
    }
    
    // Create a file for writing
    file, err := os.Create("users.json")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    // Create JSON encoder
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ") // Pretty print
    
    // Encode users to file
    err = encoder.Encode(users)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Users written to users.json")
    
    // Read back from file
    file, err = os.Open("users.json")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    // Create JSON decoder
    decoder := json.NewDecoder(file)
    
    var readUsers []User
    err = decoder.Decode(&readUsers)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("\nUsers read from file:")
    for _, user := range readUsers {
        fmt.Printf("ID: %d, Username: %s, Email: %s\n", 
            user.ID, user.Username, user.Email)
    }
    
    // Clean up
    os.Remove("users.json")
}
```

### Streaming JSON Processing

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "strings"
)

type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    Message   string `json:"message"`
}

func main() {
    // Simulate multiple JSON objects in a stream
    jsonStream := `
    {"timestamp": "2023-01-01T10:00:00Z", "level": "INFO", "message": "Application started"}
    {"timestamp": "2023-01-01T10:01:00Z", "level": "WARN", "message": "Low memory warning"}
    {"timestamp": "2023-01-01T10:02:00Z", "level": "ERROR", "message": "Database connection failed"}
    {"timestamp": "2023-01-01T10:03:00Z", "level": "INFO", "message": "Database connection restored"}
    `
    
    reader := strings.NewReader(jsonStream)
    decoder := json.NewDecoder(reader)
    
    fmt.Println("Processing log entries:")
    
    for decoder.More() {
        var entry LogEntry
        err := decoder.Decode(&entry)
        if err != nil {
            log.Printf("Error decoding JSON: %v", err)
            continue
        }
        
        // Process each log entry
        fmt.Printf("[%s] %s: %s\n", entry.Timestamp, entry.Level, entry.Message)
        
        // Example: Filter error messages
        if entry.Level == "ERROR" {
            fmt.Printf("  -> ERROR DETECTED: %s\n", entry.Message)
        }
    }
}
```

## Error Handling in JSON Operations

### Common JSON Errors

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func safeUnmarshal(data []byte) {
    var person Person
    err := json.Unmarshal(data, &person)
    
    if err != nil {
        switch err := err.(type) {
        case *json.SyntaxError:
            fmt.Printf("Syntax error at offset %d: %v\n", err.Offset, err)
        case *json.UnmarshalTypeError:
            fmt.Printf("Type error: cannot unmarshal %s into Go value of type %s\n", 
                err.Value, err.Type)
        default:
            fmt.Printf("Other JSON error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Successfully parsed: %+v\n", person)
}

func main() {
    // Valid JSON
    validJSON := `{"name": "John", "age": 30}`
    fmt.Println("Testing valid JSON:")
    safeUnmarshal([]byte(validJSON))
    
    // Invalid syntax
    invalidSyntax := `{"name": "John", "age": 30` // Missing closing brace
    fmt.Println("\nTesting invalid syntax:")
    safeUnmarshal([]byte(invalidSyntax))
    
    // Type mismatch
    typeMismatch := `{"name": "John", "age": "thirty"}` // Age should be number
    fmt.Println("\nTesting type mismatch:")
    safeUnmarshal([]byte(typeMismatch))
    
    // Extra fields (this is actually OK)
    extraFields := `{"name": "John", "age": 30, "city": "New York"}`
    fmt.Println("\nTesting extra fields:")
    safeUnmarshal([]byte(extraFields))
}
```

### Validating JSON Structure

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

type Config struct {
    DatabaseURL string `json:"database_url"`
    Port        int    `json:"port"`
    Debug       bool   `json:"debug"`
}

func (c *Config) Validate() error {
    if c.DatabaseURL == "" {
        return fmt.Errorf("database_url is required")
    }
    
    if c.Port <= 0 || c.Port > 65535 {
        return fmt.Errorf("port must be between 1 and 65535")
    }
    
    return nil
}

func parseConfig(jsonData []byte) (*Config, error) {
    var config Config
    
    // First, check if JSON is valid
    if !json.Valid(jsonData) {
        return nil, fmt.Errorf("invalid JSON")
    }
    
    // Unmarshal
    err := json.Unmarshal(jsonData, &config)
    if err != nil {
        return nil, fmt.Errorf("failed to parse JSON: %w", err)
    }
    
    // Validate business logic
    err = config.Validate()
    if err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }
    
    return &config, nil
}

func main() {
    // Valid config
    validConfig := `{
        "database_url": "postgres://localhost:5432/mydb",
        "port": 8080,
        "debug": true
    }`
    
    config, err := parseConfig([]byte(validConfig))
    if err != nil {
        log.Printf("Error: %v", err)
    } else {
        fmt.Printf("Valid config: %+v\n", config)
    }
    
    // Invalid config - missing database_url
    invalidConfig := `{
        "port": 8080,
        "debug": true
    }`
    
    config, err = parseConfig([]byte(invalidConfig))
    if err != nil {
        log.Printf("Error: %v", err)
    } else {
        fmt.Printf("Config: %+v\n", config)
    }
    
    // Invalid config - bad port
    badPortConfig := `{
        "database_url": "postgres://localhost:5432/mydb",
        "port": 70000,
        "debug": true
    }`
    
    config, err = parseConfig([]byte(badPortConfig))
    if err != nil {
        log.Printf("Error: %v", err)
    } else {
        fmt.Printf("Config: %+v\n", config)
    }
}
```

## Practical Examples

### REST API Response Handling

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"
)

type APIResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message,omitempty"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    
    if err := json.NewEncoder(w).Encode(response); err != nil {
        log.Printf("Error encoding JSON response: %v", err)
    }
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    // Simulate getting user data
    user := User{
        ID:        1,
        Username:  "johndoe",
        Email:     "john@example.com",
        CreatedAt: time.Now(),
    }
    
    response := APIResponse{
        Success: true,
        Message: "User retrieved successfully",
        Data:    user,
    }
    
    sendJSONResponse(w, http.StatusOK, response)
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        response := APIResponse{
            Success: false,
            Error:   "Method not allowed",
        }
        sendJSONResponse(w, http.StatusMethodNotAllowed, response)
        return
    }
    
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        response := APIResponse{
            Success: false,
            Error:   "Invalid JSON payload",
        }
        sendJSONResponse(w, http.StatusBadRequest, response)
        return
    }
    
    // Simulate user creation
    user.ID = 123
    user.CreatedAt = time.Now()
    
    response := APIResponse{
        Success: true,
        Message: "User created successfully",
        Data:    user,
    }
    
    sendJSONResponse(w, http.StatusCreated, response)
}

func main() {
    http.HandleFunc("/user", getUserHandler)
    http.HandleFunc("/users", createUserHandler)
    
    fmt.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

### Configuration File Management

```go
package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path/filepath"
)

type DatabaseConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    Database string `json:"database"`
}

type ServerConfig struct {
    Port         int  `json:"port"`
    Debug        bool `json:"debug"`
    ReadTimeout  int  `json:"read_timeout"`
    WriteTimeout int  `json:"write_timeout"`
}

type AppConfig struct {
    AppName  string         `json:"app_name"`
    Version  string         `json:"version"`
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    Features map[string]bool `json:"features"`
}

func (c *AppConfig) LoadFromFile(filename string) error {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }
    
    err = json.Unmarshal(data, c)
    if err != nil {
        return fmt.Errorf("failed to parse config file: %w", err)
    }
    
    return nil
}

func (c *AppConfig) SaveToFile(filename string) error {
    // Create directory if it doesn't exist
    dir := filepath.Dir(filename)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return fmt.Errorf("failed to create config directory: %w", err)
    }
    
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return fmt.Errorf("failed to marshal config: %w", err)
    }
    
    err = ioutil.WriteFile(filename, data, 0644)
    if err != nil {
        return fmt.Errorf("failed to write config file: %w", err)
    }
    
    return nil
}

func (c *AppConfig) SetDefaults() {
    c.AppName = "MyApp"
    c.Version = "1.0.0"
    c.Database = DatabaseConfig{
        Host:     "localhost",
        Port:     5432,
        Username: "user",
        Password: "password",
        Database: "myapp",
    }
    c.Server = ServerConfig{
        Port:         8080,
        Debug:        false,
        ReadTimeout:  30,
        WriteTimeout: 30,
    }
    c.Features = map[string]bool{
        "authentication": true,
        "logging":        true,
        "metrics":        false,
    }
}

func main() {
    configFile := "config/app.json"
    
    var config AppConfig
    
    // Try to load existing config
    err := config.LoadFromFile(configFile)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Println("Config file not found, creating with defaults...")
            config.SetDefaults()
            
            err = config.SaveToFile(configFile)
            if err != nil {
                log.Fatal("Failed to save default config:", err)
            }
            
            fmt.Printf("Default config saved to %s\n", configFile)
        } else {
            log.Fatal("Failed to load config:", err)
        }
    } else {
        fmt.Printf("Config loaded from %s\n", configFile)
    }
    
    // Display current config
    fmt.Println("\nCurrent configuration:")
    fmt.Printf("App Name: %s\n", config.AppName)
    fmt.Printf("Version: %s\n", config.Version)
    fmt.Printf("Database: %s@%s:%d/%s\n", 
        config.Database.Username, config.Database.Host, 
        config.Database.Port, config.Database.Database)
    fmt.Printf("Server Port: %d\n", config.Server.Port)
    fmt.Printf("Debug Mode: %t\n", config.Server.Debug)
    
    fmt.Println("\nFeatures:")
    for feature, enabled := range config.Features {
        status := "disabled"
        if enabled {
            status = "enabled"
        }
        fmt.Printf("  %s: %s\n", feature, status)
    }
    
    // Clean up
    os.RemoveAll("config")
}
```

## Best Practices

### 1. Use Struct Tags Appropriately

```go
// Good: Clear field mapping and omitempty for optional fields
type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email,omitempty"`
    Password string `json:"-"` // Never include in JSON
}

// Avoid: Inconsistent naming
type BadUser struct {
    ID       int    `json:"user_id"`
    Name     string `json:"userName"`  // Inconsistent case
    Email    string `json:"email_addr"` // Inconsistent naming
}
```

### 2. Handle Errors Properly

```go
func parseJSON(data []byte, v interface{}) error {
    if !json.Valid(data) {
        return fmt.Errorf("invalid JSON")
    }
    
    err := json.Unmarshal(data, v)
    if err != nil {
        return fmt.Errorf("failed to unmarshal JSON: %w", err)
    }
    
    return nil
}
```

### 3. Use Streaming for Large Data

```go
// For large JSON files, use streaming
func processLargeJSONFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    decoder := json.NewDecoder(file)
    
    for decoder.More() {
        var item MyStruct
        if err := decoder.Decode(&item); err != nil {
            return err
        }
        
        // Process item
        processItem(item)
    }
    
    return nil
}
```

### 4. Validate Input Data

```go
type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

func (r *CreateUserRequest) Validate() error {
    if r.Username == "" {
        return fmt.Errorf("username is required")
    }
    
    if r.Email == "" {
        return fmt.Errorf("email is required")
    }
    
    if r.Age < 0 || r.Age > 150 {
        return fmt.Errorf("age must be between 0 and 150")
    }
    
    return nil
}
```

## Exercises

### Exercise 1: User Management System

Create a user management system that can:
1. Load users from a JSON file
2. Add new users
3. Update existing users
4. Save users back to JSON file
5. Handle validation errors

### Exercise 2: Configuration Parser

Build a configuration parser that:
1. Supports nested configuration objects
2. Has default values
3. Validates configuration values
4. Supports environment variable overrides

### Exercise 3: API Response Handler

Implement an API response handler that:
1. Handles different response formats
2. Includes proper error handling
3. Supports pagination metadata
4. Validates response structure

## Key Takeaways

- Go's `encoding/json` package provides powerful JSON handling capabilities
- Use struct tags to control JSON marshaling/unmarshaling behavior
- Implement custom marshaling for complex data types
- Always handle JSON errors appropriately
- Use streaming for large JSON data
- Validate input data for security and correctness
- Follow consistent naming conventions in JSON APIs
- Use `omitempty` for optional fields to keep JSON clean

## Next Steps

Next, we'll explore [HTTP Client and Server](19-http-client-server.md) to learn how to build web applications and APIs that work with JSON data.