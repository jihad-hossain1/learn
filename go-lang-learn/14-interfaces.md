# 14. Interfaces

Interfaces in Go are one of the most powerful features of the language. They define a contract that types must fulfill by implementing specific methods. Go's interfaces are implicitly satisfied, meaning you don't need to explicitly declare that a type implements an interface.

## What are Interfaces?

An interface is a type that specifies a method set. Any type that implements all the methods in the interface automatically satisfies that interface.

### Basic Interface Syntax

```go
type InterfaceName interface {
    MethodName(parameters) ReturnType
    AnotherMethod(parameters) ReturnType
}
```

### Simple Interface Example

```go
package main

import "fmt"

// Shape interface defines what it means to be a shape
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle implements Shape interface
type Rectangle struct {
    Width, Height float64
}

func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle implements Shape interface
type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// Function that works with any Shape
func printShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 3}
    
    fmt.Println("Rectangle:")
    printShapeInfo(rect)
    
    fmt.Println("\nCircle:")
    printShapeInfo(circle)
}
```

## Empty Interface

The empty interface `interface{}` can hold values of any type since every type implements at least zero methods.

```go
package main

import "fmt"

func printAnything(value interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", value, value)
}

func main() {
    printAnything(42)
    printAnything("hello")
    printAnything(3.14)
    printAnything([]int{1, 2, 3})
    printAnything(map[string]int{"a": 1})
    
    // Slice of empty interfaces
    var items []interface{}
    items = append(items, 1, "two", 3.0, true)
    
    for i, item := range items {
        fmt.Printf("Item %d: %v (%T)\n", i, item, item)
    }
}
```

## Type Assertions

Type assertions allow you to extract the underlying concrete value from an interface.

### Basic Type Assertion

```go
package main

import "fmt"

func main() {
    var i interface{} = "hello"
    
    // Type assertion
    s := i.(string)
    fmt.Println(s) // hello
    
    // Safe type assertion with ok idiom
    s, ok := i.(string)
    if ok {
        fmt.Printf("String value: %s\n", s)
    }
    
    // This would panic
    // n := i.(int) // panic: interface conversion
    
    // Safe version
    n, ok := i.(int)
    if ok {
        fmt.Printf("Int value: %d\n", n)
    } else {
        fmt.Println("Not an int")
    }
}
```

### Type Assertion with Custom Types

```go
package main

import "fmt"

type Writer interface {
    Write([]byte) (int, error)
}

type FileWriter struct {
    filename string
}

func (fw FileWriter) Write(data []byte) (int, error) {
    fmt.Printf("Writing to file %s: %s\n", fw.filename, string(data))
    return len(data), nil
}

type ConsoleWriter struct{}

func (cw ConsoleWriter) Write(data []byte) (int, error) {
    fmt.Printf("Console: %s\n", string(data))
    return len(data), nil
}

func processWriter(w Writer) {
    // Type assertion to check specific type
    if fw, ok := w.(FileWriter); ok {
        fmt.Printf("This is a FileWriter for: %s\n", fw.filename)
    } else if _, ok := w.(ConsoleWriter); ok {
        fmt.Println("This is a ConsoleWriter")
    }
    
    w.Write([]byte("Hello, World!"))
}

func main() {
    fileWriter := FileWriter{filename: "output.txt"}
    consoleWriter := ConsoleWriter{}
    
    processWriter(fileWriter)
    processWriter(consoleWriter)
}
```

## Type Switches

Type switches allow you to handle different types in a clean way.

```go
package main

import "fmt"

func describe(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Integer: %d\n", v)
    case string:
        fmt.Printf("String: %s (length: %d)\n", v, len(v))
    case bool:
        fmt.Printf("Boolean: %t\n", v)
    case []int:
        fmt.Printf("Slice of ints: %v (length: %d)\n", v, len(v))
    case map[string]int:
        fmt.Printf("Map: %v\n", v)
    case nil:
        fmt.Println("Nil value")
    default:
        fmt.Printf("Unknown type: %T\n", v)
    }
}

func main() {
    describe(42)
    describe("hello")
    describe(true)
    describe([]int{1, 2, 3})
    describe(map[string]int{"a": 1, "b": 2})
    describe(nil)
    describe(3.14)
}
```

### Advanced Type Switch Example

```go
package main

import "fmt"

type Processor interface {
    Process() string
}

type TextProcessor struct {
    text string
}

func (tp TextProcessor) Process() string {
    return "Processed text: " + tp.text
}

type NumberProcessor struct {
    number int
}

func (np NumberProcessor) Process() string {
    return fmt.Sprintf("Processed number: %d", np.number*2)
}

func handleValue(value interface{}) {
    switch v := value.(type) {
    case Processor:
        fmt.Println("Processor result:", v.Process())
    case string:
        fmt.Println("Raw string:", v)
    case int:
        fmt.Println("Raw number:", v)
    case []Processor:
        fmt.Println("Processing slice of processors:")
        for i, p := range v {
            fmt.Printf("  %d: %s\n", i, p.Process())
        }
    default:
        fmt.Printf("Cannot handle type: %T\n", v)
    }
}

func main() {
    tp := TextProcessor{text: "hello"}
    np := NumberProcessor{number: 21}
    
    handleValue(tp)
    handleValue(np)
    handleValue("raw string")
    handleValue(42)
    handleValue([]Processor{tp, np})
    handleValue(3.14)
}
```

## Interface Composition

Interfaces can be composed of other interfaces.

```go
package main

import "fmt"

// Basic interfaces
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type Closer interface {
    Close() error
}

// Composed interfaces
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// Implementation
type File struct {
    name string
    data []byte
    pos  int
}

func (f *File) Read(p []byte) (int, error) {
    if f.pos >= len(f.data) {
        return 0, fmt.Errorf("EOF")
    }
    
    n := copy(p, f.data[f.pos:])
    f.pos += n
    return n, nil
}

func (f *File) Write(p []byte) (int, error) {
    f.data = append(f.data, p...)
    return len(p), nil
}

func (f *File) Close() error {
    fmt.Printf("Closing file: %s\n", f.name)
    return nil
}

// Functions using composed interfaces
func copyData(src Reader, dst Writer) error {
    buffer := make([]byte, 1024)
    for {
        n, err := src.Read(buffer)
        if err != nil {
            break
        }
        if n > 0 {
            _, writeErr := dst.Write(buffer[:n])
            if writeErr != nil {
                return writeErr
            }
        }
    }
    return nil
}

func processFile(rwc ReadWriteCloser) {
    defer rwc.Close()
    
    // Write some data
    rwc.Write([]byte("Hello, World!"))
    
    // Read it back
    buffer := make([]byte, 100)
    // Reset position for reading (in a real implementation)
    n, _ := rwc.Read(buffer)
    fmt.Printf("Read: %s\n", string(buffer[:n]))
}

func main() {
    file := &File{
        name: "example.txt",
        data: []byte("Initial content"),
        pos:  0,
    }
    
    // File implements ReadWriteCloser
    processFile(file)
    
    // Create another file for copying
    src := &File{
        name: "source.txt",
        data: []byte("Source data"),
        pos:  0,
    }
    
    dst := &File{
        name: "dest.txt",
        data: []byte{},
        pos:  0,
    }
    
    copyData(src, dst)
    fmt.Printf("Destination data: %s\n", string(dst.data))
}
```

## Interface Values and Nil

### Understanding Interface Values

```go
package main

import "fmt"

type Describer interface {
    Describe() string
}

type Person struct {
    Name string
}

func (p Person) Describe() string {
    return "Person: " + p.Name
}

func (p *Person) DescribePtr() string {
    if p == nil {
        return "Person: <nil>"
    }
    return "Person pointer: " + p.Name
}

func main() {
    var d Describer
    
    // Interface is nil
    fmt.Printf("Interface: %v, Type: %T\n", d, d)
    if d == nil {
        fmt.Println("Interface is nil")
    }
    
    // Interface with nil concrete value
    var p *Person
    d = p
    fmt.Printf("Interface: %v, Type: %T\n", d, d)
    if d == nil {
        fmt.Println("Interface is nil")
    } else {
        fmt.Println("Interface is not nil (but concrete value is nil)")
        // This would panic because p is nil
        // fmt.Println(d.Describe())
    }
    
    // Interface with concrete value
    person := Person{Name: "Alice"}
    d = person
    fmt.Printf("Interface: %v, Type: %T\n", d, d)
    fmt.Println(d.Describe())
}
```

### Safe Interface Usage

```go
package main

import "fmt"

type Processor interface {
    Process() string
}

type SafeProcessor struct {
    name string
}

func (sp *SafeProcessor) Process() string {
    if sp == nil {
        return "SafeProcessor: <nil>"
    }
    return "SafeProcessor: " + sp.name
}

func safeProcess(p Processor) string {
    if p == nil {
        return "Processor interface is nil"
    }
    
    // Check if the concrete value is nil
    if sp, ok := p.(*SafeProcessor); ok && sp == nil {
        return "SafeProcessor concrete value is nil"
    }
    
    return p.Process()
}

func main() {
    var p Processor
    fmt.Println(safeProcess(p)) // Interface is nil
    
    var sp *SafeProcessor
    p = sp
    fmt.Println(safeProcess(p)) // Concrete value is nil
    
    sp = &SafeProcessor{name: "test"}
    p = sp
    fmt.Println(safeProcess(p)) // Normal processing
}
```

## Practical Examples

### Payment Processing System

```go
package main

import (
    "fmt"
    "errors"
)

// Payment interface
type PaymentProcessor interface {
    ProcessPayment(amount float64) error
    GetFee(amount float64) float64
    GetName() string
}

// Credit Card implementation
type CreditCard struct {
    Number string
    Name   string
}

func (cc CreditCard) ProcessPayment(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    fmt.Printf("Processing $%.2f via Credit Card ending in %s\n", 
        amount, cc.Number[len(cc.Number)-4:])
    return nil
}

func (cc CreditCard) GetFee(amount float64) float64 {
    return amount * 0.029 // 2.9% fee
}

func (cc CreditCard) GetName() string {
    return "Credit Card"
}

// PayPal implementation
type PayPal struct {
    Email string
}

func (pp PayPal) ProcessPayment(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    fmt.Printf("Processing $%.2f via PayPal account %s\n", amount, pp.Email)
    return nil
}

func (pp PayPal) GetFee(amount float64) float64 {
    return amount * 0.034 // 3.4% fee
}

func (pp PayPal) GetName() string {
    return "PayPal"
}

// Bank Transfer implementation
type BankTransfer struct {
    AccountNumber string
    RoutingNumber string
}

func (bt BankTransfer) ProcessPayment(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    fmt.Printf("Processing $%.2f via Bank Transfer to account %s\n", 
        amount, bt.AccountNumber)
    return nil
}

func (bt BankTransfer) GetFee(amount float64) float64 {
    return 0.50 // Flat $0.50 fee
}

func (bt BankTransfer) GetName() string {
    return "Bank Transfer"
}

// Payment service
type PaymentService struct {
    processors []PaymentProcessor
}

func NewPaymentService() *PaymentService {
    return &PaymentService{
        processors: make([]PaymentProcessor, 0),
    }
}

func (ps *PaymentService) AddProcessor(processor PaymentProcessor) {
    ps.processors = append(ps.processors, processor)
}

func (ps *PaymentService) ProcessPayment(amount float64, processorName string) error {
    for _, processor := range ps.processors {
        if processor.GetName() == processorName {
            fee := processor.GetFee(amount)
            total := amount + fee
            
            fmt.Printf("Payment Details:\n")
            fmt.Printf("  Amount: $%.2f\n", amount)
            fmt.Printf("  Fee: $%.2f\n", fee)
            fmt.Printf("  Total: $%.2f\n", total)
            
            return processor.ProcessPayment(total)
        }
    }
    return errors.New("payment processor not found")
}

func (ps *PaymentService) GetBestProcessor(amount float64) PaymentProcessor {
    if len(ps.processors) == 0 {
        return nil
    }
    
    best := ps.processors[0]
    bestFee := best.GetFee(amount)
    
    for _, processor := range ps.processors[1:] {
        fee := processor.GetFee(amount)
        if fee < bestFee {
            best = processor
            bestFee = fee
        }
    }
    
    return best
}

func main() {
    // Create payment service
    service := NewPaymentService()
    
    // Add payment processors
    service.AddProcessor(CreditCard{Number: "1234567890123456", Name: "John Doe"})
    service.AddProcessor(PayPal{Email: "john@example.com"})
    service.AddProcessor(BankTransfer{AccountNumber: "987654321", RoutingNumber: "123456789"})
    
    amount := 100.0
    
    // Process payments with different methods
    fmt.Println("=== Processing Payments ===")
    service.ProcessPayment(amount, "Credit Card")
    fmt.Println()
    service.ProcessPayment(amount, "PayPal")
    fmt.Println()
    service.ProcessPayment(amount, "Bank Transfer")
    fmt.Println()
    
    // Find best processor
    fmt.Println("=== Best Processor ===")
    best := service.GetBestProcessor(amount)
    if best != nil {
        fmt.Printf("Best processor for $%.2f: %s (fee: $%.2f)\n", 
            amount, best.GetName(), best.GetFee(amount))
    }
}
```

### Notification System

```go
package main

import (
    "fmt"
    "strings"
    "time"
)

// Notification interface
type Notifier interface {
    Send(message string) error
    GetType() string
}

// Email notifier
type EmailNotifier struct {
    SMTPServer string
    Port       int
    Username   string
    Recipients []string
}

func (en EmailNotifier) Send(message string) error {
    fmt.Printf("[EMAIL] Sending to %s via %s:%d\n", 
        strings.Join(en.Recipients, ", "), en.SMTPServer, en.Port)
    fmt.Printf("[EMAIL] Message: %s\n", message)
    return nil
}

func (en EmailNotifier) GetType() string {
    return "Email"
}

// SMS notifier
type SMSNotifier struct {
    APIKey      string
    PhoneNumbers []string
}

func (sn SMSNotifier) Send(message string) error {
    fmt.Printf("[SMS] Sending to %s\n", strings.Join(sn.PhoneNumbers, ", "))
    fmt.Printf("[SMS] Message: %s\n", message)
    return nil
}

func (sn SMSNotifier) GetType() string {
    return "SMS"
}

// Slack notifier
type SlackNotifier struct {
    WebhookURL string
    Channel    string
}

func (sl SlackNotifier) Send(message string) error {
    fmt.Printf("[SLACK] Sending to #%s\n", sl.Channel)
    fmt.Printf("[SLACK] Message: %s\n", message)
    return nil
}

func (sl SlackNotifier) GetType() string {
    return "Slack"
}

// Push notification
type PushNotifier struct {
    AppID    string
    DeviceIDs []string
}

func (pn PushNotifier) Send(message string) error {
    fmt.Printf("[PUSH] Sending to %d devices\n", len(pn.DeviceIDs))
    fmt.Printf("[PUSH] Message: %s\n", message)
    return nil
}

func (pn PushNotifier) GetType() string {
    return "Push"
}

// Notification service
type NotificationService struct {
    notifiers []Notifier
}

func NewNotificationService() *NotificationService {
    return &NotificationService{
        notifiers: make([]Notifier, 0),
    }
}

func (ns *NotificationService) AddNotifier(notifier Notifier) {
    ns.notifiers = append(ns.notifiers, notifier)
}

func (ns *NotificationService) SendToAll(message string) {
    fmt.Printf("=== Broadcasting message to all notifiers ===\n")
    for _, notifier := range ns.notifiers {
        if err := notifier.Send(message); err != nil {
            fmt.Printf("Error sending via %s: %v\n", notifier.GetType(), err)
        }
    }
    fmt.Println()
}

func (ns *NotificationService) SendToType(message string, notifierType string) {
    fmt.Printf("=== Sending message via %s ===\n", notifierType)
    for _, notifier := range ns.notifiers {
        if notifier.GetType() == notifierType {
            if err := notifier.Send(message); err != nil {
                fmt.Printf("Error sending via %s: %v\n", notifier.GetType(), err)
            }
            return
        }
    }
    fmt.Printf("No notifier found for type: %s\n", notifierType)
    fmt.Println()
}

func (ns *NotificationService) SendUrgent(message string) {
    urgentTypes := []string{"SMS", "Push"}
    fmt.Printf("=== Sending urgent message ===\n")
    
    for _, urgentType := range urgentTypes {
        for _, notifier := range ns.notifiers {
            if notifier.GetType() == urgentType {
                if err := notifier.Send("URGENT: " + message); err != nil {
                    fmt.Printf("Error sending urgent via %s: %v\n", notifier.GetType(), err)
                }
            }
        }
    }
    fmt.Println()
}

// Alert levels
type AlertLevel int

const (
    Info AlertLevel = iota
    Warning
    Error
    Critical
)

func (ns *NotificationService) SendAlert(message string, level AlertLevel) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    formattedMessage := fmt.Sprintf("[%s] %s", timestamp, message)
    
    switch level {
    case Info:
        ns.SendToType(formattedMessage, "Email")
    case Warning:
        ns.SendToType(formattedMessage, "Slack")
    case Error:
        ns.SendToAll(formattedMessage)
    case Critical:
        ns.SendUrgent(formattedMessage)
    }
}

func main() {
    // Create notification service
    service := NewNotificationService()
    
    // Add notifiers
    service.AddNotifier(EmailNotifier{
        SMTPServer: "smtp.gmail.com",
        Port:       587,
        Username:   "admin@company.com",
        Recipients: []string{"user1@company.com", "user2@company.com"},
    })
    
    service.AddNotifier(SMSNotifier{
        APIKey:       "sms-api-key",
        PhoneNumbers: []string{"+1234567890", "+0987654321"},
    })
    
    service.AddNotifier(SlackNotifier{
        WebhookURL: "https://hooks.slack.com/webhook",
        Channel:    "alerts",
    })
    
    service.AddNotifier(PushNotifier{
        AppID:     "com.company.app",
        DeviceIDs: []string{"device1", "device2", "device3"},
    })
    
    // Send different types of alerts
    service.SendAlert("System maintenance scheduled", Info)
    service.SendAlert("High CPU usage detected", Warning)
    service.SendAlert("Database connection failed", Error)
    service.SendAlert("Server is down!", Critical)
    
    // Send custom messages
    service.SendToAll("Welcome to our new notification system!")
    service.SendToType("Check your email for updates", "Email")
}
```

### Logger Interface

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"
)

// Logger interface
type Logger interface {
    Info(message string)
    Warning(message string)
    Error(message string)
    Debug(message string)
}

// Console logger
type ConsoleLogger struct {
    prefix string
}

func NewConsoleLogger(prefix string) *ConsoleLogger {
    return &ConsoleLogger{prefix: prefix}
}

func (cl *ConsoleLogger) Info(message string) {
    fmt.Printf("[%s] INFO: %s\n", cl.prefix, message)
}

func (cl *ConsoleLogger) Warning(message string) {
    fmt.Printf("[%s] WARNING: %s\n", cl.prefix, message)
}

func (cl *ConsoleLogger) Error(message string) {
    fmt.Printf("[%s] ERROR: %s\n", cl.prefix, message)
}

func (cl *ConsoleLogger) Debug(message string) {
    fmt.Printf("[%s] DEBUG: %s\n", cl.prefix, message)
}

// File logger
type FileLogger struct {
    filename string
    logger   *log.Logger
}

func NewFileLogger(filename string) (*FileLogger, error) {
    file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }
    
    return &FileLogger{
        filename: filename,
        logger:   log.New(file, "", log.LstdFlags),
    }, nil
}

func (fl *FileLogger) Info(message string) {
    fl.logger.Printf("INFO: %s", message)
}

func (fl *FileLogger) Warning(message string) {
    fl.logger.Printf("WARNING: %s", message)
}

func (fl *FileLogger) Error(message string) {
    fl.logger.Printf("ERROR: %s", message)
}

func (fl *FileLogger) Debug(message string) {
    fl.logger.Printf("DEBUG: %s", message)
}

// Multi logger (logs to multiple destinations)
type MultiLogger struct {
    loggers []Logger
}

func NewMultiLogger(loggers ...Logger) *MultiLogger {
    return &MultiLogger{loggers: loggers}
}

func (ml *MultiLogger) Info(message string) {
    for _, logger := range ml.loggers {
        logger.Info(message)
    }
}

func (ml *MultiLogger) Warning(message string) {
    for _, logger := range ml.loggers {
        logger.Warning(message)
    }
}

func (ml *MultiLogger) Error(message string) {
    for _, logger := range ml.loggers {
        logger.Error(message)
    }
}

func (ml *MultiLogger) Debug(message string) {
    for _, logger := range ml.loggers {
        logger.Debug(message)
    }
}

// Application service that uses logger
type UserService struct {
    logger Logger
    users  map[int]string
}

func NewUserService(logger Logger) *UserService {
    return &UserService{
        logger: logger,
        users:  make(map[int]string),
    }
}

func (us *UserService) CreateUser(id int, name string) error {
    us.logger.Info(fmt.Sprintf("Creating user: %d - %s", id, name))
    
    if _, exists := us.users[id]; exists {
        us.logger.Error(fmt.Sprintf("User %d already exists", id))
        return fmt.Errorf("user %d already exists", id)
    }
    
    us.users[id] = name
    us.logger.Info(fmt.Sprintf("User %d created successfully", id))
    return nil
}

func (us *UserService) GetUser(id int) (string, error) {
    us.logger.Debug(fmt.Sprintf("Looking up user: %d", id))
    
    name, exists := us.users[id]
    if !exists {
        us.logger.Warning(fmt.Sprintf("User %d not found", id))
        return "", fmt.Errorf("user %d not found", id)
    }
    
    us.logger.Debug(fmt.Sprintf("User %d found: %s", id, name))
    return name, nil
}

func (us *UserService) DeleteUser(id int) error {
    us.logger.Info(fmt.Sprintf("Deleting user: %d", id))
    
    if _, exists := us.users[id]; !exists {
        us.logger.Error(fmt.Sprintf("Cannot delete user %d: not found", id))
        return fmt.Errorf("user %d not found", id)
    }
    
    delete(us.users, id)
    us.logger.Info(fmt.Sprintf("User %d deleted successfully", id))
    return nil
}

func main() {
    // Create different loggers
    consoleLogger := NewConsoleLogger("APP")
    
    fileLogger, err := NewFileLogger("app.log")
    if err != nil {
        fmt.Printf("Error creating file logger: %v\n", err)
        return
    }
    
    // Create multi-logger that logs to both console and file
    multiLogger := NewMultiLogger(consoleLogger, fileLogger)
    
    // Create user service with multi-logger
    userService := NewUserService(multiLogger)
    
    // Demonstrate logging
    fmt.Println("=== User Service Operations ===")
    
    userService.CreateUser(1, "Alice")
    userService.CreateUser(2, "Bob")
    userService.CreateUser(1, "Charlie") // This will fail
    
    if name, err := userService.GetUser(1); err == nil {
        fmt.Printf("Retrieved user: %s\n", name)
    }
    
    userService.GetUser(999) // This will fail
    
    userService.DeleteUser(2)
    userService.DeleteUser(999) // This will fail
    
    fmt.Println("\n=== Check app.log file for file logger output ===")
}
```

## Best Practices

### 1. Keep Interfaces Small

```go
// Good: Small, focused interface
type Writer interface {
    Write([]byte) (int, error)
}

// Good: Another small interface
type Reader interface {
    Read([]byte) (int, error)
}

// Good: Compose when needed
type ReadWriter interface {
    Reader
    Writer
}

// Avoid: Large interfaces
// type FileManager interface {
//     Read([]byte) (int, error)
//     Write([]byte) (int, error)
//     Seek(int64, int) (int64, error)
//     Close() error
//     Stat() (FileInfo, error)
//     Chmod(FileMode) error
//     // ... many more methods
// }
```

### 2. Accept Interfaces, Return Concrete Types

```go
// Good: Accept interface
func ProcessData(r io.Reader) *ProcessedData {
    // Implementation
    return &ProcessedData{}
}

// Good: Return concrete type
func NewFileProcessor(filename string) *FileProcessor {
    return &FileProcessor{filename: filename}
}

// Avoid: Returning interfaces unless necessary
// func NewProcessor() Processor {
//     return &FileProcessor{}
// }
```

### 3. Use Interface Segregation

```go
// Good: Segregated interfaces
type Validator interface {
    Validate() error
}

type Serializer interface {
    Serialize() ([]byte, error)
}

type Deserializer interface {
    Deserialize([]byte) error
}

// Compose when needed
type ValidatingSerializer interface {
    Validator
    Serializer
}
```

### 4. Handle Nil Interface Values

```go
func SafeProcess(p Processor) string {
    if p == nil {
        return "No processor provided"
    }
    return p.Process()
}
```

## Exercises

### Exercise 1: File Storage System

Create a file storage system with different storage backends.

```go
type FileStorage interface {
    Save(filename string, data []byte) error
    Load(filename string) ([]byte, error)
    Delete(filename string) error
    Exists(filename string) bool
    List() ([]string, error)
}

// Implement:
// - LocalFileStorage
// - MemoryFileStorage
// - CloudFileStorage (simulated)
```

### Exercise 2: Database Connection Pool

Create a database interface with different implementations.

```go
type Database interface {
    Connect() error
    Disconnect() error
    Query(sql string, args ...interface{}) ([]map[string]interface{}, error)
    Execute(sql string, args ...interface{}) error
    IsConnected() bool
}

// Implement:
// - MySQLDatabase
// - PostgreSQLDatabase
// - MockDatabase (for testing)
```

### Exercise 3: Message Queue System

Create a message queue interface with different implementations.

```go
type MessageQueue interface {
    Publish(topic string, message []byte) error
    Subscribe(topic string) (<-chan []byte, error)
    Unsubscribe(topic string) error
    Close() error
}

// Implement:
// - InMemoryQueue
// - RedisQueue (simulated)
// - RabbitMQQueue (simulated)
```

## Key Takeaways

1. **Interfaces define contracts**: They specify what methods a type must implement.

2. **Implicit satisfaction**: Types automatically satisfy interfaces by implementing the required methods.

3. **Empty interface accepts anything**: `interface{}` can hold any value.

4. **Type assertions extract concrete values**: Use the comma ok idiom for safe assertions.

5. **Type switches handle multiple types**: Clean way to handle different types in one place.

6. **Interface composition**: Build complex interfaces from simpler ones.

7. **Interface values can be nil**: Always check for nil when necessary.

8. **Accept interfaces, return concrete types**: This provides flexibility for callers.

9. **Keep interfaces small**: Follow the Interface Segregation Principle.

10. **Interfaces enable polymorphism**: Different types can be used interchangeably if they satisfy the same interface.

## Next Steps

Now that you understand interfaces, let's learn about [Error Handling](15-error-handling.md) to see how Go approaches error management!

---

**Previous**: [← Methods](13-methods.md) | **Next**: [Error Handling →](15-error-handling.md)