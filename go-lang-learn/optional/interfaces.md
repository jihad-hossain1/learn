# 12. Interfaces

Interfaces in Go define a contract by specifying a set of method signatures. They enable polymorphism and loose coupling between components. Go interfaces are implemented implicitly - any type that has the methods defined by an interface automatically implements that interface.

## Interface Basics

### Interface Declaration

```go
// Basic interface declaration
type InterfaceName interface {
    MethodName1(parameters) returnType
    MethodName2(parameters) returnType
    // ...
}

// Example
type Writer interface {
    Write([]byte) (int, error)
}

type Reader interface {
    Read([]byte) (int, error)
}
```

### Simple Interface Example

```go
package main

import "fmt"

// Define an interface
type Shape interface {
    Area() float64
    Perimeter() float64
}

// Rectangle type
type Rectangle struct {
    Width, Height float64
}

// Circle type
type Circle struct {
    Radius float64
}

// Rectangle implements Shape interface
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}

// Circle implements Shape interface
func (c Circle) Area() float64 {
    return 3.14159 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * 3.14159 * c.Radius
}

// Function that works with any Shape
func PrintShapeInfo(s Shape) {
    fmt.Printf("Area: %.2f\n", s.Area())
    fmt.Printf("Perimeter: %.2f\n", s.Perimeter())
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    circle := Circle{Radius: 3}
    
    fmt.Println("Rectangle:")
    PrintShapeInfo(rect)
    
    fmt.Println("\nCircle:")
    PrintShapeInfo(circle)
    
    // Using interface as a slice
    shapes := []Shape{rect, circle}
    
    fmt.Println("\nAll shapes:")
    for i, shape := range shapes {
        fmt.Printf("Shape %d:\n", i+1)
        PrintShapeInfo(shape)
        fmt.Println()
    }
}
```

### Empty Interface

The empty interface `interface{}` can hold values of any type.

```go
package main

import "fmt"

// Function that accepts any type
func PrintAnything(value interface{}) {
    fmt.Printf("Value: %v, Type: %T\n", value, value)
}

func main() {
    // Empty interface can hold any value
    var anything interface{}
    
    anything = 42
    fmt.Printf("Integer: %v\n", anything)
    
    anything = "Hello, World!"
    fmt.Printf("String: %v\n", anything)
    
    anything = []int{1, 2, 3}
    fmt.Printf("Slice: %v\n", anything)
    
    anything = map[string]int{"a": 1, "b": 2}
    fmt.Printf("Map: %v\n", anything)
    
    // Using with function
    PrintAnything(42)
    PrintAnything("Hello")
    PrintAnything([]string{"a", "b", "c"})
    
    // Slice of empty interfaces
    mixed := []interface{}{42, "hello", true, 3.14}
    fmt.Println("\nMixed slice:")
    for i, item := range mixed {
        fmt.Printf("Item %d: %v (type: %T)\n", i, item, item)
    }
}
```

### Type Assertions

Type assertions extract the underlying concrete value from an interface.

```go
package main

import "fmt"

func main() {
    var i interface{} = "Hello, World!"
    
    // Type assertion (unsafe - panics if wrong type)
    str := i.(string)
    fmt.Printf("String value: %s\n", str)
    
    // Safe type assertion with ok idiom
    str2, ok := i.(string)
    if ok {
        fmt.Printf("Successfully extracted string: %s\n", str2)
    } else {
        fmt.Println("Not a string")
    }
    
    // Trying wrong type
    num, ok := i.(int)
    if ok {
        fmt.Printf("Number: %d\n", num)
    } else {
        fmt.Println("Not an integer")
    }
    
    // Function to demonstrate type assertions
    processValue := func(value interface{}) {
        switch v := value.(type) {
        case string:
            fmt.Printf("String: %s (length: %d)\n", v, len(v))
        case int:
            fmt.Printf("Integer: %d (squared: %d)\n", v, v*v)
        case float64:
            fmt.Printf("Float: %.2f (rounded: %.0f)\n", v, v)
        case bool:
            fmt.Printf("Boolean: %t\n", v)
        default:
            fmt.Printf("Unknown type: %T\n", v)
        }
    }
    
    fmt.Println("\nProcessing different values:")
    processValue("Hello")
    processValue(42)
    processValue(3.14159)
    processValue(true)
    processValue([]int{1, 2, 3})
}
```

### Type Switches

Type switches provide a clean way to handle different types.

```go
package main

import "fmt"

func DescribeValue(value interface{}) {
    switch v := value.(type) {
    case nil:
        fmt.Println("Value is nil")
    case bool:
        if v {
            fmt.Println("Boolean: true")
        } else {
            fmt.Println("Boolean: false")
        }
    case int:
        fmt.Printf("Integer: %d\n", v)
        if v > 0 {
            fmt.Println("  - Positive number")
        } else if v < 0 {
            fmt.Println("  - Negative number")
        } else {
            fmt.Println("  - Zero")
        }
    case float64:
        fmt.Printf("Float: %.2f\n", v)
    case string:
        fmt.Printf("String: %q (length: %d)\n", v, len(v))
        if len(v) == 0 {
            fmt.Println("  - Empty string")
        }
    case []int:
        fmt.Printf("Integer slice: %v (length: %d)\n", v, len(v))
    case []string:
        fmt.Printf("String slice: %v (length: %d)\n", v, len(v))
    case map[string]int:
        fmt.Printf("String-to-int map: %v (size: %d)\n", v, len(v))
    default:
        fmt.Printf("Unknown type: %T with value: %v\n", v, v)
    }
}

func main() {
    values := []interface{}{
        nil,
        true,
        42,
        -10,
        0,
        3.14159,
        "Hello, World!",
        "",
        []int{1, 2, 3},
        []string{"a", "b", "c"},
        map[string]int{"x": 1, "y": 2},
        struct{ Name string }{"John"},
    }
    
    for i, value := range values {
        fmt.Printf("Value %d: ", i+1)
        DescribeValue(value)
        fmt.Println()
    }
}
```

### Interface Composition

Interfaces can be composed of other interfaces.

```go
package main

import (
    "fmt"
    "io"
    "strings"
)

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

// Custom type that implements multiple interfaces
type Buffer struct {
    data   []byte
    pos    int
    closed bool
}

func NewBuffer() *Buffer {
    return &Buffer{data: make([]byte, 0)}
}

func (b *Buffer) Read(p []byte) (int, error) {
    if b.closed {
        return 0, fmt.Errorf("buffer is closed")
    }
    
    if b.pos >= len(b.data) {
        return 0, io.EOF
    }
    
    n := copy(p, b.data[b.pos:])
    b.pos += n
    return n, nil
}

func (b *Buffer) Write(p []byte) (int, error) {
    if b.closed {
        return 0, fmt.Errorf("buffer is closed")
    }
    
    b.data = append(b.data, p...)
    return len(p), nil
}

func (b *Buffer) Close() error {
    if b.closed {
        return fmt.Errorf("buffer already closed")
    }
    b.closed = true
    return nil
}

func (b *Buffer) String() string {
    return string(b.data)
}

// Functions that work with composed interfaces
func CopyData(dst Writer, src Reader) error {
    buffer := make([]byte, 1024)
    for {
        n, err := src.Read(buffer)
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        
        _, err = dst.Write(buffer[:n])
        if err != nil {
            return err
        }
    }
    return nil
}

func ProcessReadWriteCloser(rwc ReadWriteCloser) {
    defer rwc.Close()
    
    // Write some data
    data := []byte("Hello, Interface Composition!")
    rwc.Write(data)
    
    fmt.Println("Data written successfully")
}

func main() {
    // Create buffer that implements ReadWriteCloser
    buffer := NewBuffer()
    
    // Use as Writer
    buffer.Write([]byte("Hello, "))
    buffer.Write([]byte("World!"))
    
    fmt.Printf("Buffer content: %s\n", buffer.String())
    
    // Reset position for reading
    buffer.pos = 0
    
    // Use as Reader
    readBuffer := make([]byte, 5)
    n, err := buffer.Read(readBuffer)
    if err != nil && err != io.EOF {
        fmt.Printf("Read error: %v\n", err)
    } else {
        fmt.Printf("Read %d bytes: %s\n", n, string(readBuffer[:n]))
    }
    
    // Use with composed interface function
    buffer2 := NewBuffer()
    src := strings.NewReader("Data to copy")
    
    err = CopyData(buffer2, src)
    if err != nil {
        fmt.Printf("Copy error: %v\n", err)
    } else {
        fmt.Printf("Copied data: %s\n", buffer2.String())
    }
    
    // Use as ReadWriteCloser
    buffer3 := NewBuffer()
    ProcessReadWriteCloser(buffer3)
    
    // Demonstrate interface satisfaction
    var r Reader = buffer
    var w Writer = buffer
    var c Closer = buffer
    var rw ReadWriter = buffer
    var rwc ReadWriteCloser = buffer
    
    fmt.Printf("Buffer implements Reader: %T\n", r)
    fmt.Printf("Buffer implements Writer: %T\n", w)
    fmt.Printf("Buffer implements Closer: %T\n", c)
    fmt.Printf("Buffer implements ReadWriter: %T\n", rw)
    fmt.Printf("Buffer implements ReadWriteCloser: %T\n", rwc)
}
```

### Interface Values and nil

```go
package main

import "fmt"

type Printer interface {
    Print() string
}

type Document struct {
    Content string
}

func (d *Document) Print() string {
    if d == nil {
        return "<nil document>"
    }
    return d.Content
}

type EmptyPrinter struct{}

func (e EmptyPrinter) Print() string {
    return "<empty>"
}

func main() {
    var p Printer
    
    // nil interface
    fmt.Printf("nil interface: %v, %T\n", p, p)
    fmt.Printf("Is nil: %t\n", p == nil)
    
    // Interface with nil concrete value
    var doc *Document
    p = doc
    fmt.Printf("Interface with nil concrete value: %v, %T\n", p, p)
    fmt.Printf("Is nil: %t\n", p == nil) // false! Interface is not nil
    
    // Call method on nil concrete value
    fmt.Printf("Print result: %s\n", p.Print())
    
    // Interface with non-nil concrete value
    doc = &Document{Content: "Hello, World!"}
    p = doc
    fmt.Printf("Interface with concrete value: %v, %T\n", p, p)
    fmt.Printf("Print result: %s\n", p.Print())
    
    // Different concrete type
    p = EmptyPrinter{}
    fmt.Printf("Different concrete type: %v, %T\n", p, p)
    fmt.Printf("Print result: %s\n", p.Print())
    
    // Checking for nil interface vs nil concrete value
    checkPrinter := func(printer Printer) {
        if printer == nil {
            fmt.Println("Printer interface is nil")
            return
        }
        
        // Check if concrete value is nil (for pointer types)
        if doc, ok := printer.(*Document); ok && doc == nil {
            fmt.Println("Printer has nil Document")
        } else {
            fmt.Printf("Printer result: %s\n", printer.Print())
        }
    }
    
    fmt.Println("\nChecking different printer states:")
    
    var nilInterface Printer
    checkPrinter(nilInterface)
    
    var nilDoc *Document
    checkPrinter(nilDoc)
    
    validDoc := &Document{Content: "Valid document"}
    checkPrinter(validDoc)
    
    checkPrinter(EmptyPrinter{})
}
```

## Practical Examples

### Example 1: Payment Processing System

```go
package main

import (
    "fmt"
    "time"
)

// Payment interface
type PaymentProcessor interface {
    ProcessPayment(amount float64) (*PaymentResult, error)
    GetProviderName() string
    ValidatePayment(amount float64) error
}

// Payment result
type PaymentResult struct {
    TransactionID string
    Amount        float64
    Status        string
    ProcessedAt   time.Time
    Provider      string
}

// Credit Card processor
type CreditCardProcessor struct {
    CardNumber string
    ExpiryDate string
    CVV        string
}

func (cc *CreditCardProcessor) ProcessPayment(amount float64) (*PaymentResult, error) {
    if err := cc.ValidatePayment(amount); err != nil {
        return nil, err
    }
    
    // Simulate payment processing
    time.Sleep(100 * time.Millisecond)
    
    return &PaymentResult{
        TransactionID: fmt.Sprintf("CC_%d", time.Now().Unix()),
        Amount:        amount,
        Status:        "SUCCESS",
        ProcessedAt:   time.Now(),
        Provider:      cc.GetProviderName(),
    }, nil
}

func (cc *CreditCardProcessor) GetProviderName() string {
    return "Credit Card"
}

func (cc *CreditCardProcessor) ValidatePayment(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    if amount > 10000 {
        return fmt.Errorf("amount exceeds credit card limit")
    }
    if len(cc.CardNumber) != 16 {
        return fmt.Errorf("invalid card number")
    }
    return nil
}

// PayPal processor
type PayPalProcessor struct {
    Email    string
    Password string
}

func (pp *PayPalProcessor) ProcessPayment(amount float64) (*PaymentResult, error) {
    if err := pp.ValidatePayment(amount); err != nil {
        return nil, err
    }
    
    // Simulate PayPal API call
    time.Sleep(200 * time.Millisecond)
    
    return &PaymentResult{
        TransactionID: fmt.Sprintf("PP_%d", time.Now().Unix()),
        Amount:        amount,
        Status:        "SUCCESS",
        ProcessedAt:   time.Now(),
        Provider:      pp.GetProviderName(),
    }, nil
}

func (pp *PayPalProcessor) GetProviderName() string {
    return "PayPal"
}

func (pp *PayPalProcessor) ValidatePayment(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    if amount > 5000 {
        return fmt.Errorf("amount exceeds PayPal limit")
    }
    if pp.Email == "" {
        return fmt.Errorf("email is required")
    }
    return nil
}

// Bank Transfer processor
type BankTransferProcessor struct {
    AccountNumber string
    RoutingNumber string
}

func (bt *BankTransferProcessor) ProcessPayment(amount float64) (*PaymentResult, error) {
    if err := bt.ValidatePayment(amount); err != nil {
        return nil, err
    }
    
    // Simulate bank transfer processing
    time.Sleep(500 * time.Millisecond)
    
    return &PaymentResult{
        TransactionID: fmt.Sprintf("BT_%d", time.Now().Unix()),
        Amount:        amount,
        Status:        "PENDING", // Bank transfers are usually pending
        ProcessedAt:   time.Now(),
        Provider:      bt.GetProviderName(),
    }, nil
}

func (bt *BankTransferProcessor) GetProviderName() string {
    return "Bank Transfer"
}

func (bt *BankTransferProcessor) ValidatePayment(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    if amount < 10 {
        return fmt.Errorf("minimum transfer amount is $10")
    }
    if len(bt.AccountNumber) == 0 || len(bt.RoutingNumber) == 0 {
        return fmt.Errorf("account and routing numbers are required")
    }
    return nil
}

// Payment service that uses the interface
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

func (ps *PaymentService) ProcessPaymentWithBestOption(amount float64) (*PaymentResult, error) {
    var lastError error
    
    for _, processor := range ps.processors {
        result, err := processor.ProcessPayment(amount)
        if err == nil {
            fmt.Printf("Payment processed successfully with %s\n", processor.GetProviderName())
            return result, nil
        }
        
        fmt.Printf("Failed to process with %s: %v\n", processor.GetProviderName(), err)
        lastError = err
    }
    
    return nil, fmt.Errorf("all payment processors failed, last error: %v", lastError)
}

func (ps *PaymentService) GetAvailableProcessors() []string {
    var names []string
    for _, processor := range ps.processors {
        names = append(names, processor.GetProviderName())
    }
    return names
}

func main() {
    // Create payment service
    paymentService := NewPaymentService()
    
    // Add different payment processors
    paymentService.AddProcessor(&CreditCardProcessor{
        CardNumber: "1234567890123456",
        ExpiryDate: "12/25",
        CVV:        "123",
    })
    
    paymentService.AddProcessor(&PayPalProcessor{
        Email:    "user@example.com",
        Password: "password123",
    })
    
    paymentService.AddProcessor(&BankTransferProcessor{
        AccountNumber: "1234567890",
        RoutingNumber: "987654321",
    })
    
    fmt.Printf("Available processors: %v\n\n", paymentService.GetAvailableProcessors())
    
    // Test different payment amounts
    amounts := []float64{50.0, 3000.0, 8000.0, 15000.0}
    
    for _, amount := range amounts {
        fmt.Printf("Processing payment of $%.2f:\n", amount)
        result, err := paymentService.ProcessPaymentWithBestOption(amount)
        
        if err != nil {
            fmt.Printf("Payment failed: %v\n", err)
        } else {
            fmt.Printf("Payment successful: %+v\n", result)
        }
        fmt.Println()
    }
    
    // Test individual processors
    fmt.Println("Testing individual processors:")
    
    processors := []PaymentProcessor{
        &CreditCardProcessor{CardNumber: "1234567890123456", ExpiryDate: "12/25", CVV: "123"},
        &PayPalProcessor{Email: "user@example.com", Password: "password123"},
        &BankTransferProcessor{AccountNumber: "1234567890", RoutingNumber: "987654321"},
    }
    
    testAmount := 100.0
    for _, processor := range processors {
        fmt.Printf("\nTesting %s with $%.2f:\n", processor.GetProviderName(), testAmount)
        result, err := processor.ProcessPayment(testAmount)
        
        if err != nil {
            fmt.Printf("Error: %v\n", err)
        } else {
            fmt.Printf("Success: Transaction ID: %s, Status: %s\n", 
                result.TransactionID, result.Status)
        }
    }
}
```

### Example 2: Notification System

```go
package main

import (
    "fmt"
    "time"
)

// Notification interface
type NotificationSender interface {
    Send(message string, recipient string) error
    GetProviderName() string
    IsAvailable() bool
}

// Email notification
type EmailSender struct {
    SMTPServer string
    Port       int
    Username   string
    Password   string
}

func (e *EmailSender) Send(message, recipient string) error {
    if !e.IsAvailable() {
        return fmt.Errorf("email service is not available")
    }
    
    // Simulate email sending
    fmt.Printf("📧 Sending email to %s via %s\n", recipient, e.SMTPServer)
    fmt.Printf("   Message: %s\n", message)
    time.Sleep(100 * time.Millisecond)
    
    return nil
}

func (e *EmailSender) GetProviderName() string {
    return "Email"
}

func (e *EmailSender) IsAvailable() bool {
    return e.SMTPServer != "" && e.Username != ""
}

// SMS notification
type SMSSender struct {
    APIKey    string
    APISecret string
    Provider  string
}

func (s *SMSSender) Send(message, recipient string) error {
    if !s.IsAvailable() {
        return fmt.Errorf("SMS service is not available")
    }
    
    // Simulate SMS sending
    fmt.Printf("📱 Sending SMS to %s via %s\n", recipient, s.Provider)
    fmt.Printf("   Message: %s\n", message)
    time.Sleep(200 * time.Millisecond)
    
    return nil
}

func (s *SMSSender) GetProviderName() string {
    return "SMS"
}

func (s *SMSSender) IsAvailable() bool {
    return s.APIKey != "" && s.APISecret != ""
}

// Push notification
type PushSender struct {
    AppID     string
    APIKey    string
    ServerKey string
}

func (p *PushSender) Send(message, recipient string) error {
    if !p.IsAvailable() {
        return fmt.Errorf("push notification service is not available")
    }
    
    // Simulate push notification
    fmt.Printf("🔔 Sending push notification to %s\n", recipient)
    fmt.Printf("   Message: %s\n", message)
    time.Sleep(50 * time.Millisecond)
    
    return nil
}

func (p *PushSender) GetProviderName() string {
    return "Push Notification"
}

func (p *PushSender) IsAvailable() bool {
    return p.AppID != "" && p.APIKey != ""
}

// Slack notification
type SlackSender struct {
    WebhookURL string
    Channel    string
}

func (sl *SlackSender) Send(message, recipient string) error {
    if !sl.IsAvailable() {
        return fmt.Errorf("Slack service is not available")
    }
    
    // Simulate Slack message
    channel := sl.Channel
    if recipient != "" {
        channel = "@" + recipient
    }
    
    fmt.Printf("💬 Sending Slack message to %s\n", channel)
    fmt.Printf("   Message: %s\n", message)
    time.Sleep(150 * time.Millisecond)
    
    return nil
}

func (sl *SlackSender) GetProviderName() string {
    return "Slack"
}

func (sl *SlackSender) IsAvailable() bool {
    return sl.WebhookURL != ""
}

// Notification service
type NotificationService struct {
    senders []NotificationSender
}

func NewNotificationService() *NotificationService {
    return &NotificationService{
        senders: make([]NotificationSender, 0),
    }
}

func (ns *NotificationService) AddSender(sender NotificationSender) {
    ns.senders = append(ns.senders, sender)
}

func (ns *NotificationService) SendNotification(message, recipient string) []error {
    var errors []error
    
    for _, sender := range ns.senders {
        if !sender.IsAvailable() {
            errors = append(errors, fmt.Errorf("%s is not available", sender.GetProviderName()))
            continue
        }
        
        err := sender.Send(message, recipient)
        if err != nil {
            errors = append(errors, fmt.Errorf("%s failed: %v", sender.GetProviderName(), err))
        }
    }
    
    return errors
}

func (ns *NotificationService) SendWithFallback(message, recipient string) error {
    for _, sender := range ns.senders {
        if !sender.IsAvailable() {
            continue
        }
        
        err := sender.Send(message, recipient)
        if err == nil {
            fmt.Printf("✅ Successfully sent via %s\n", sender.GetProviderName())
            return nil
        }
        
        fmt.Printf("❌ Failed to send via %s: %v\n", sender.GetProviderName(), err)
    }
    
    return fmt.Errorf("all notification methods failed")
}

func (ns *NotificationService) GetAvailableSenders() []string {
    var available []string
    for _, sender := range ns.senders {
        if sender.IsAvailable() {
            available = append(available, sender.GetProviderName())
        }
    }
    return available
}

func (ns *NotificationService) SendToSpecificProvider(providerName, message, recipient string) error {
    for _, sender := range ns.senders {
        if sender.GetProviderName() == providerName {
            if !sender.IsAvailable() {
                return fmt.Errorf("%s is not available", providerName)
            }
            return sender.Send(message, recipient)
        }
    }
    return fmt.Errorf("provider %s not found", providerName)
}

func main() {
    // Create notification service
    notificationService := NewNotificationService()
    
    // Add different notification senders
    notificationService.AddSender(&EmailSender{
        SMTPServer: "smtp.gmail.com",
        Port:       587,
        Username:   "sender@example.com",
        Password:   "password",
    })
    
    notificationService.AddSender(&SMSSender{
        APIKey:    "sms_api_key",
        APISecret: "sms_api_secret",
        Provider:  "Twilio",
    })
    
    notificationService.AddSender(&PushSender{
        AppID:     "app_123",
        APIKey:    "push_api_key",
        ServerKey: "server_key",
    })
    
    notificationService.AddSender(&SlackSender{
        WebhookURL: "https://hooks.slack.com/webhook",
        Channel:    "#general",
    })
    
    // Add an unavailable sender for testing
    notificationService.AddSender(&EmailSender{
        SMTPServer: "", // Missing configuration
        Username:   "",
    })
    
    fmt.Printf("Available senders: %v\n\n", notificationService.GetAvailableSenders())
    
    // Test sending to all available providers
    fmt.Println("=== Sending to all providers ===")
    message := "Hello! This is a test notification."
    recipient := "user@example.com"
    
    errors := notificationService.SendNotification(message, recipient)
    if len(errors) > 0 {
        fmt.Println("\nErrors encountered:")
        for _, err := range errors {
            fmt.Printf("- %v\n", err)
        }
    }
    
    fmt.Println("\n=== Testing fallback mechanism ===")
    urgentMessage := "URGENT: System maintenance in 5 minutes!"
    err := notificationService.SendWithFallback(urgentMessage, "admin@example.com")
    if err != nil {
        fmt.Printf("Fallback failed: %v\n", err)
    }
    
    fmt.Println("\n=== Testing specific provider ===")
    specificMessage := "This message is sent via SMS only"
    err = notificationService.SendToSpecificProvider("SMS", specificMessage, "+1234567890")
    if err != nil {
        fmt.Printf("Specific provider failed: %v\n", err)
    }
    
    // Test with different message types
    fmt.Println("\n=== Testing different message types ===")
    
    messages := map[string]string{
        "Welcome":     "Welcome to our service! We're glad to have you.",
        "Alert":       "⚠️ Security alert: Unusual login detected",
        "Reminder":    "📅 Don't forget about your appointment tomorrow",
        "Promotion":   "🎉 Special offer: 50% off this weekend only!",
        "System":      "🔧 System will be down for maintenance at 2 AM",
    }
    
    for msgType, msg := range messages {
        fmt.Printf("\nSending %s message:\n", msgType)
        notificationService.SendWithFallback(msg, "customer@example.com")
    }
}
```

## Exercises

### Exercise 1: File Storage System
Create a file storage system with different storage providers.

```go
package main

import "time"

type FileInfo struct {
    Name         string
    Size         int64
    LastModified time.Time
    ContentType  string
}

type StorageProvider interface {
    // Implement these methods:
    // Upload(filename string, data []byte) error
    // Download(filename string) ([]byte, error)
    // Delete(filename string) error
    // List() ([]FileInfo, error)
    // GetProviderName() string
    // IsAvailable() bool
}

// Implement these storage providers:
// LocalStorage, S3Storage, GoogleCloudStorage, DropboxStorage

type FileStorageService struct {
    // Add fields and methods
}

func main() {
    // Test the file storage system
}
```

### Exercise 2: Database Connection Pool
Create a database connection pool with different database types.

```go
package main

type QueryResult struct {
    Rows     []map[string]interface{}
    Affected int64
    Error    error
}

type Database interface {
    // Implement these methods:
    // Connect() error
    // Disconnect() error
    // Query(sql string, args ...interface{}) QueryResult
    // Execute(sql string, args ...interface{}) QueryResult
    // IsConnected() bool
    // GetDatabaseType() string
}

// Implement these database types:
// MySQLDatabase, PostgreSQLDatabase, SQLiteDatabase, MongoDatabase

type ConnectionPool struct {
    // Add fields and methods
}

func main() {
    // Test the database connection pool
}
```

### Exercise 3: Logger System
Create a flexible logging system with different output destinations.

```go
package main

import "time"

type LogLevel int

const (
    DEBUG LogLevel = iota
    INFO
    WARN
    ERROR
    FATAL
)

type LogEntry struct {
    Level     LogLevel
    Message   string
    Timestamp time.Time
    Source    string
}

type Logger interface {
    // Implement these methods:
    // Log(entry LogEntry) error
    // SetLevel(level LogLevel)
    // GetName() string
    // IsEnabled(level LogLevel) bool
}

// Implement these loggers:
// ConsoleLogger, FileLogger, DatabaseLogger, RemoteLogger

type LoggingService struct {
    // Add fields and methods
}

func main() {
    // Test the logging system
}
```

## Key Takeaways

1. **Interfaces** define contracts through method signatures
2. **Implicit implementation** - no explicit "implements" keyword needed
3. **Empty interface** `interface{}` can hold any value
4. **Type assertions** extract concrete values from interfaces
5. **Type switches** handle different types cleanly
6. **Interface composition** combines multiple interfaces
7. **Interface values** can be nil or contain nil concrete values
8. **Polymorphism** achieved through interfaces
9. **Loose coupling** between components via interfaces
10. **Duck typing** - "If it walks like a duck and quacks like a duck, it's a duck"

## Next Steps

Now that you understand interfaces, let's learn about [Error Handling](13-error-handling.md) for robust Go applications!

---

**Previous**: [← Structs](11-structs.md) | **Next**: [Error Handling →](13-error-handling.md)