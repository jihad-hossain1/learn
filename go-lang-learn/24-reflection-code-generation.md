# 20. Reflection and Code Generation

Reflection in Go allows programs to examine their own structure at runtime. The `reflect` package provides the ability to inspect types and values, manipulate them dynamically, and even generate code. While powerful, reflection should be used judiciously as it can make code harder to understand and maintain.

## Understanding Reflection

### The reflect Package

Go's reflection is built around two main types:
- `reflect.Type`: Represents a Go type
- `reflect.Value`: Represents a Go value

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var x float64 = 3.4
    
    // Get the type and value
    t := reflect.TypeOf(x)
    v := reflect.ValueOf(x)
    
    fmt.Println("Type:", t)           // Type: float64
    fmt.Println("Value:", v)          // Value: 3.4
    fmt.Println("Kind:", t.Kind())    // Kind: float64
    fmt.Println("String:", t.String()) // String: float64
}
```

### Basic Type Inspection

```go
package main

import (
    "fmt"
    "reflect"
)

func inspectType(x interface{}) {
    t := reflect.TypeOf(x)
    v := reflect.ValueOf(x)
    
    fmt.Printf("Type: %v\n", t)
    fmt.Printf("Kind: %v\n", t.Kind())
    fmt.Printf("Value: %v\n", v)
    fmt.Printf("Interface: %v\n", v.Interface())
    
    // Check if the value can be set
    fmt.Printf("CanSet: %v\n", v.CanSet())
    
    // Get the underlying value if it's a pointer
    if t.Kind() == reflect.Ptr {
        fmt.Printf("Elem Type: %v\n", t.Elem())
        fmt.Printf("Elem Value: %v\n", v.Elem())
        fmt.Printf("Elem CanSet: %v\n", v.Elem().CanSet())
    }
    
    fmt.Println("---")
}

func main() {
    var i int = 42
    var s string = "hello"
    var f float64 = 3.14
    var ptr *int = &i
    
    inspectType(i)
    inspectType(s)
    inspectType(f)
    inspectType(ptr)
}
```

### Working with Pointers and Addressability

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    var x float64 = 3.4
    
    // This won't work - value is not addressable
    v := reflect.ValueOf(x)
    fmt.Println("CanSet:", v.CanSet()) // false
    
    // This will work - pass a pointer
    p := reflect.ValueOf(&x)
    fmt.Println("Pointer CanSet:", p.CanSet()) // false
    fmt.Println("Elem CanSet:", p.Elem().CanSet()) // true
    
    // Set the value
    p.Elem().SetFloat(7.1)
    fmt.Println("New value:", x) // 7.1
}
```

## Struct Reflection

### Inspecting Struct Fields

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name    string `json:"name" validate:"required"`
    Age     int    `json:"age" validate:"min=0,max=150"`
    Email   string `json:"email" validate:"email"`
    private string // unexported field
}

func inspectStruct(s interface{}) {
    t := reflect.TypeOf(s)
    v := reflect.ValueOf(s)
    
    // Handle pointer to struct
    if t.Kind() == reflect.Ptr {
        t = t.Elem()
        v = v.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        fmt.Println("Not a struct")
        return
    }
    
    fmt.Printf("Struct: %s\n", t.Name())
    fmt.Printf("Number of fields: %d\n", t.NumField())
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        fmt.Printf("Field %d:\n", i)
        fmt.Printf("  Name: %s\n", field.Name)
        fmt.Printf("  Type: %s\n", field.Type)
        fmt.Printf("  Tag: %s\n", field.Tag)
        fmt.Printf("  JSON Tag: %s\n", field.Tag.Get("json"))
        fmt.Printf("  Validate Tag: %s\n", field.Tag.Get("validate"))
        fmt.Printf("  Exported: %v\n", field.IsExported())
        
        if value.CanInterface() {
            fmt.Printf("  Value: %v\n", value.Interface())
        } else {
            fmt.Printf("  Value: <unexported>\n")
        }
        
        fmt.Println()
    }
}

func main() {
    p := Person{
        Name:    "John Doe",
        Age:     30,
        Email:   "john@example.com",
        private: "secret",
    }
    
    inspectStruct(p)
    inspectStruct(&p)
}
```

### Modifying Struct Fields

```go
package main

import (
    "fmt"
    "reflect"
)

type Config struct {
    Host     string
    Port     int
    Debug    bool
    Features []string
}

func setField(obj interface{}, fieldName string, value interface{}) error {
    v := reflect.ValueOf(obj)
    
    // Must be a pointer to be settable
    if v.Kind() != reflect.Ptr {
        return fmt.Errorf("object must be a pointer")
    }
    
    v = v.Elem()
    if v.Kind() != reflect.Struct {
        return fmt.Errorf("object must be a struct")
    }
    
    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        return fmt.Errorf("field %s not found", fieldName)
    }
    
    if !field.CanSet() {
        return fmt.Errorf("field %s cannot be set", fieldName)
    }
    
    fieldType := field.Type()
    valueType := reflect.TypeOf(value)
    
    if !valueType.AssignableTo(fieldType) {
        return fmt.Errorf("value type %s is not assignable to field type %s", 
            valueType, fieldType)
    }
    
    field.Set(reflect.ValueOf(value))
    return nil
}

func getField(obj interface{}, fieldName string) (interface{}, error) {
    v := reflect.ValueOf(obj)
    
    // Handle pointer
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }
    
    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("object must be a struct")
    }
    
    field := v.FieldByName(fieldName)
    if !field.IsValid() {
        return nil, fmt.Errorf("field %s not found", fieldName)
    }
    
    if !field.CanInterface() {
        return nil, fmt.Errorf("field %s cannot be accessed", fieldName)
    }
    
    return field.Interface(), nil
}

func main() {
    config := &Config{
        Host:     "localhost",
        Port:     8080,
        Debug:    false,
        Features: []string{"auth", "logging"},
    }
    
    fmt.Printf("Original config: %+v\n", config)
    
    // Get field values
    if host, err := getField(config, "Host"); err == nil {
        fmt.Printf("Host: %v\n", host)
    }
    
    // Set field values
    if err := setField(config, "Host", "0.0.0.0"); err != nil {
        fmt.Printf("Error setting Host: %v\n", err)
    }
    
    if err := setField(config, "Port", 9090); err != nil {
        fmt.Printf("Error setting Port: %v\n", err)
    }
    
    if err := setField(config, "Debug", true); err != nil {
        fmt.Printf("Error setting Debug: %v\n", err)
    }
    
    if err := setField(config, "Features", []string{"auth", "logging", "metrics"}); err != nil {
        fmt.Printf("Error setting Features: %v\n", err)
    }
    
    fmt.Printf("Modified config: %+v\n", config)
}
```

## Method Reflection

### Inspecting and Calling Methods

```go
package main

import (
    "fmt"
    "reflect"
)

type Calculator struct {
    Value float64
}

func (c *Calculator) Add(x float64) float64 {
    c.Value += x
    return c.Value
}

func (c *Calculator) Subtract(x float64) float64 {
    c.Value -= x
    return c.Value
}

func (c *Calculator) Multiply(x float64) float64 {
    c.Value *= x
    return c.Value
}

func (c Calculator) GetValue() float64 {
    return c.Value
}

func (c Calculator) String() string {
    return fmt.Sprintf("Calculator{Value: %.2f}", c.Value)
}

func inspectMethods(obj interface{}) {
    t := reflect.TypeOf(obj)
    v := reflect.ValueOf(obj)
    
    fmt.Printf("Type: %s\n", t)
    fmt.Printf("Number of methods: %d\n", t.NumMethod())
    
    for i := 0; i < t.NumMethod(); i++ {
        method := t.Method(i)
        fmt.Printf("Method %d: %s\n", i, method.Name)
        fmt.Printf("  Type: %s\n", method.Type)
        fmt.Printf("  NumIn: %d\n", method.Type.NumIn())
        fmt.Printf("  NumOut: %d\n", method.Type.NumOut())
        
        // Print input types
        for j := 0; j < method.Type.NumIn(); j++ {
            fmt.Printf("  In[%d]: %s\n", j, method.Type.In(j))
        }
        
        // Print output types
        for j := 0; j < method.Type.NumOut(); j++ {
            fmt.Printf("  Out[%d]: %s\n", j, method.Type.Out(j))
        }
        
        fmt.Println()
    }
}

func callMethod(obj interface{}, methodName string, args ...interface{}) ([]reflect.Value, error) {
    v := reflect.ValueOf(obj)
    method := v.MethodByName(methodName)
    
    if !method.IsValid() {
        return nil, fmt.Errorf("method %s not found", methodName)
    }
    
    // Convert arguments to reflect.Value
    in := make([]reflect.Value, len(args))
    for i, arg := range args {
        in[i] = reflect.ValueOf(arg)
    }
    
    // Call the method
    return method.Call(in), nil
}

func main() {
    calc := &Calculator{Value: 10.0}
    
    fmt.Println("=== Inspecting methods ===")
    inspectMethods(calc)
    
    fmt.Println("=== Calling methods dynamically ===")
    
    // Call Add method
    if results, err := callMethod(calc, "Add", 5.0); err == nil {
        fmt.Printf("Add(5.0) = %v\n", results[0].Interface())
    }
    
    // Call Multiply method
    if results, err := callMethod(calc, "Multiply", 2.0); err == nil {
        fmt.Printf("Multiply(2.0) = %v\n", results[0].Interface())
    }
    
    // Call GetValue method
    if results, err := callMethod(calc, "GetValue"); err == nil {
        fmt.Printf("GetValue() = %v\n", results[0].Interface())
    }
    
    // Call String method
    if results, err := callMethod(calc, "String"); err == nil {
        fmt.Printf("String() = %v\n", results[0].Interface())
    }
}
```

## Slice and Map Reflection

### Working with Slices

```go
package main

import (
    "fmt"
    "reflect"
)

func inspectSlice(s interface{}) {
    v := reflect.ValueOf(s)
    t := reflect.TypeOf(s)
    
    if t.Kind() != reflect.Slice {
        fmt.Println("Not a slice")
        return
    }
    
    fmt.Printf("Slice type: %s\n", t)
    fmt.Printf("Element type: %s\n", t.Elem())
    fmt.Printf("Length: %d\n", v.Len())
    fmt.Printf("Capacity: %d\n", v.Cap())
    
    // Print elements
    for i := 0; i < v.Len(); i++ {
        elem := v.Index(i)
        fmt.Printf("[%d]: %v\n", i, elem.Interface())
    }
}

func createSlice(elementType reflect.Type, length, capacity int) reflect.Value {
    sliceType := reflect.SliceOf(elementType)
    return reflect.MakeSlice(sliceType, length, capacity)
}

func appendToSlice(slice reflect.Value, elements ...interface{}) reflect.Value {
    for _, elem := range elements {
        elemValue := reflect.ValueOf(elem)
        slice = reflect.Append(slice, elemValue)
    }
    return slice
}

func main() {
    // Inspect existing slice
    numbers := []int{1, 2, 3, 4, 5}
    fmt.Println("=== Inspecting existing slice ===")
    inspectSlice(numbers)
    
    // Create new slice dynamically
    fmt.Println("\n=== Creating slice dynamically ===")
    intType := reflect.TypeOf(0)
    newSlice := createSlice(intType, 0, 5)
    
    // Append elements
    newSlice = appendToSlice(newSlice, 10, 20, 30)
    
    fmt.Printf("Created slice: %v\n", newSlice.Interface())
    inspectSlice(newSlice.Interface())
    
    // Modify slice elements
    fmt.Println("\n=== Modifying slice elements ===")
    v := reflect.ValueOf(&numbers).Elem() // Get addressable value
    
    // Set element at index 0
    v.Index(0).SetInt(100)
    fmt.Printf("Modified slice: %v\n", numbers)
}
```

### Working with Maps

```go
package main

import (
    "fmt"
    "reflect"
)

func inspectMap(m interface{}) {
    v := reflect.ValueOf(m)
    t := reflect.TypeOf(m)
    
    if t.Kind() != reflect.Map {
        fmt.Println("Not a map")
        return
    }
    
    fmt.Printf("Map type: %s\n", t)
    fmt.Printf("Key type: %s\n", t.Key())
    fmt.Printf("Value type: %s\n", t.Elem())
    fmt.Printf("Length: %d\n", v.Len())
    
    // Print key-value pairs
    keys := v.MapKeys()
    for _, key := range keys {
        value := v.MapIndex(key)
        fmt.Printf("%v: %v\n", key.Interface(), value.Interface())
    }
}

func createMap(keyType, valueType reflect.Type) reflect.Value {
    mapType := reflect.MapOf(keyType, valueType)
    return reflect.MakeMap(mapType)
}

func setMapValue(mapValue reflect.Value, key, value interface{}) {
    keyValue := reflect.ValueOf(key)
    valueValue := reflect.ValueOf(value)
    mapValue.SetMapIndex(keyValue, valueValue)
}

func getMapValue(mapValue reflect.Value, key interface{}) (interface{}, bool) {
    keyValue := reflect.ValueOf(key)
    valueValue := mapValue.MapIndex(keyValue)
    
    if !valueValue.IsValid() {
        return nil, false
    }
    
    return valueValue.Interface(), true
}

func main() {
    // Inspect existing map
    scores := map[string]int{
        "Alice": 95,
        "Bob":   87,
        "Carol": 92,
    }
    
    fmt.Println("=== Inspecting existing map ===")
    inspectMap(scores)
    
    // Create new map dynamically
    fmt.Println("\n=== Creating map dynamically ===")
    stringType := reflect.TypeOf("")
    intType := reflect.TypeOf(0)
    newMap := createMap(stringType, intType)
    
    // Set values
    setMapValue(newMap, "David", 88)
    setMapValue(newMap, "Eve", 94)
    setMapValue(newMap, "Frank", 79)
    
    fmt.Printf("Created map: %v\n", newMap.Interface())
    inspectMap(newMap.Interface())
    
    // Get values
    fmt.Println("\n=== Getting map values ===")
    if value, ok := getMapValue(newMap, "David"); ok {
        fmt.Printf("David's score: %v\n", value)
    }
    
    if value, ok := getMapValue(newMap, "Unknown"); ok {
        fmt.Printf("Unknown's score: %v\n", value)
    } else {
        fmt.Println("Unknown not found")
    }
}
```

## Practical Applications

### JSON-like Serialization

```go
package main

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

type Person struct {
    Name    string `serialize:"name"`
    Age     int    `serialize:"age"`
    Email   string `serialize:"email"`
    Salary  float64 `serialize:"salary"`
    Active  bool   `serialize:"active"`
    private string // This won't be serialized
}

func serialize(obj interface{}) (map[string]string, error) {
    v := reflect.ValueOf(obj)
    t := reflect.TypeOf(obj)
    
    // Handle pointer
    if t.Kind() == reflect.Ptr {
        v = v.Elem()
        t = t.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        return nil, fmt.Errorf("object must be a struct")
    }
    
    result := make(map[string]string)
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        // Skip unexported fields
        if !field.IsExported() {
            continue
        }
        
        // Get the serialize tag
        tag := field.Tag.Get("serialize")
        if tag == "" {
            tag = strings.ToLower(field.Name)
        }
        
        // Convert value to string based on type
        var strValue string
        switch value.Kind() {
        case reflect.String:
            strValue = value.String()
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            strValue = strconv.FormatInt(value.Int(), 10)
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            strValue = strconv.FormatUint(value.Uint(), 10)
        case reflect.Float32, reflect.Float64:
            strValue = strconv.FormatFloat(value.Float(), 'f', -1, 64)
        case reflect.Bool:
            strValue = strconv.FormatBool(value.Bool())
        default:
            strValue = fmt.Sprintf("%v", value.Interface())
        }
        
        result[tag] = strValue
    }
    
    return result, nil
}

func deserialize(data map[string]string, obj interface{}) error {
    v := reflect.ValueOf(obj)
    
    // Must be a pointer to be settable
    if v.Kind() != reflect.Ptr {
        return fmt.Errorf("object must be a pointer")
    }
    
    v = v.Elem()
    t := v.Type()
    
    if t.Kind() != reflect.Struct {
        return fmt.Errorf("object must be a struct")
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        // Skip unexported fields
        if !field.IsExported() {
            continue
        }
        
        // Get the serialize tag
        tag := field.Tag.Get("serialize")
        if tag == "" {
            tag = strings.ToLower(field.Name)
        }
        
        // Get the string value from data
        strValue, ok := data[tag]
        if !ok {
            continue
        }
        
        // Convert string to appropriate type
        switch value.Kind() {
        case reflect.String:
            value.SetString(strValue)
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
                value.SetInt(intValue)
            }
        case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            if uintValue, err := strconv.ParseUint(strValue, 10, 64); err == nil {
                value.SetUint(uintValue)
            }
        case reflect.Float32, reflect.Float64:
            if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
                value.SetFloat(floatValue)
            }
        case reflect.Bool:
            if boolValue, err := strconv.ParseBool(strValue); err == nil {
                value.SetBool(boolValue)
            }
        }
    }
    
    return nil
}

func main() {
    person := Person{
        Name:    "John Doe",
        Age:     30,
        Email:   "john@example.com",
        Salary:  75000.50,
        Active:  true,
        private: "secret",
    }
    
    fmt.Printf("Original: %+v\n", person)
    
    // Serialize
    data, err := serialize(person)
    if err != nil {
        fmt.Printf("Serialization error: %v\n", err)
        return
    }
    
    fmt.Println("Serialized:")
    for k, v := range data {
        fmt.Printf("  %s: %s\n", k, v)
    }
    
    // Deserialize
    var newPerson Person
    if err := deserialize(data, &newPerson); err != nil {
        fmt.Printf("Deserialization error: %v\n", err)
        return
    }
    
    fmt.Printf("Deserialized: %+v\n", newPerson)
}
```

### Validation Framework

```go
package main

import (
    "fmt"
    "reflect"
    "strconv"
    "strings"
)

type User struct {
    Name     string `validate:"required,min=2,max=50"`
    Email    string `validate:"required,email"`
    Age      int    `validate:"min=0,max=150"`
    Password string `validate:"required,min=8"`
}

type ValidationError struct {
    Field   string
    Tag     string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
    var messages []string
    for _, err := range e {
        messages = append(messages, err.Error())
    }
    return strings.Join(messages, "; ")
}

func validate(obj interface{}) error {
    v := reflect.ValueOf(obj)
    t := reflect.TypeOf(obj)
    
    // Handle pointer
    if t.Kind() == reflect.Ptr {
        v = v.Elem()
        t = t.Elem()
    }
    
    if t.Kind() != reflect.Struct {
        return fmt.Errorf("object must be a struct")
    }
    
    var errors ValidationErrors
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        value := v.Field(i)
        
        // Skip unexported fields
        if !field.IsExported() {
            continue
        }
        
        // Get validation tags
        validateTag := field.Tag.Get("validate")
        if validateTag == "" {
            continue
        }
        
        // Parse validation rules
        rules := strings.Split(validateTag, ",")
        for _, rule := range rules {
            rule = strings.TrimSpace(rule)
            if err := validateField(field.Name, value, rule); err != nil {
                errors = append(errors, *err)
            }
        }
    }
    
    if len(errors) > 0 {
        return errors
    }
    
    return nil
}

func validateField(fieldName string, value reflect.Value, rule string) *ValidationError {
    parts := strings.Split(rule, "=")
    ruleName := parts[0]
    var ruleValue string
    if len(parts) > 1 {
        ruleValue = parts[1]
    }
    
    switch ruleName {
    case "required":
        if isZeroValue(value) {
            return &ValidationError{
                Field:   fieldName,
                Tag:     ruleName,
                Message: "field is required",
            }
        }
    
    case "min":
        minVal, _ := strconv.Atoi(ruleValue)
        switch value.Kind() {
        case reflect.String:
            if len(value.String()) < minVal {
                return &ValidationError{
                    Field:   fieldName,
                    Tag:     ruleName,
                    Message: fmt.Sprintf("minimum length is %d", minVal),
                }
            }
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            if value.Int() < int64(minVal) {
                return &ValidationError{
                    Field:   fieldName,
                    Tag:     ruleName,
                    Message: fmt.Sprintf("minimum value is %d", minVal),
                }
            }
        }
    
    case "max":
        maxVal, _ := strconv.Atoi(ruleValue)
        switch value.Kind() {
        case reflect.String:
            if len(value.String()) > maxVal {
                return &ValidationError{
                    Field:   fieldName,
                    Tag:     ruleName,
                    Message: fmt.Sprintf("maximum length is %d", maxVal),
                }
            }
        case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
            if value.Int() > int64(maxVal) {
                return &ValidationError{
                    Field:   fieldName,
                    Tag:     ruleName,
                    Message: fmt.Sprintf("maximum value is %d", maxVal),
                }
            }
        }
    
    case "email":
        if value.Kind() == reflect.String {
            email := value.String()
            if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
                return &ValidationError{
                    Field:   fieldName,
                    Tag:     ruleName,
                    Message: "invalid email format",
                }
            }
        }
    }
    
    return nil
}

func isZeroValue(value reflect.Value) bool {
    switch value.Kind() {
    case reflect.String:
        return value.String() == ""
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return value.Int() == 0
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        return value.Uint() == 0
    case reflect.Float32, reflect.Float64:
        return value.Float() == 0
    case reflect.Bool:
        return !value.Bool()
    case reflect.Slice, reflect.Map, reflect.Array:
        return value.Len() == 0
    case reflect.Ptr, reflect.Interface:
        return value.IsNil()
    }
    return false
}

func main() {
    // Valid user
    validUser := User{
        Name:     "John Doe",
        Email:    "john@example.com",
        Age:      30,
        Password: "secretpassword",
    }
    
    fmt.Println("=== Validating valid user ===")
    if err := validate(validUser); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Validation passed!")
    }
    
    // Invalid user
    invalidUser := User{
        Name:     "J",                // Too short
        Email:    "invalid-email",   // Invalid format
        Age:      -5,               // Negative age
        Password: "123",            // Too short
    }
    
    fmt.Println("\n=== Validating invalid user ===")
    if err := validate(invalidUser); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Validation passed!")
    }
    
    // Empty user
    emptyUser := User{}
    
    fmt.Println("\n=== Validating empty user ===")
    if err := validate(emptyUser); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("Validation passed!")
    }
}
```

## Code Generation

### Using go:generate

Go provides the `go:generate` directive to run code generation tools:

```go
// person.go
package main

//go:generate stringer -type=Status

type Status int

const (
    Active Status = iota
    Inactive
    Pending
    Suspended
)

type Person struct {
    Name   string
    Status Status
}

func main() {
    p := Person{
        Name:   "John",
        Status: Active,
    }
    
    // After running go generate, Status will have a String() method
    fmt.Printf("Person: %s, Status: %s\n", p.Name, p.Status)
}
```

Run with:
```bash
go generate
go run .
```

### Custom Code Generator

```go
// generator.go
package main

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "os"
    "strings"
    "text/template"
)

type StructInfo struct {
    Name   string
    Fields []FieldInfo
}

type FieldInfo struct {
    Name string
    Type string
}

const builderTemplate = `// Code generated by builder generator. DO NOT EDIT.

package {{.Package}}

{{range .Structs}}
// {{.Name}}Builder provides a builder pattern for {{.Name}}
type {{.Name}}Builder struct {
    {{range .Fields}}{{.Name}} {{.Type}}
    {{end}}
}

// New{{.Name}}Builder creates a new {{.Name}}Builder
func New{{.Name}}Builder() *{{.Name}}Builder {
    return &{{.Name}}Builder{}
}

{{range .Fields}}
// {{.Name}} sets the {{.Name}} field
func (b *{{$.Name}}Builder) {{.Name}}({{.Name}} {{.Type}}) *{{$.Name}}Builder {
    b.{{.Name}} = {{.Name}}
    return b
}
{{end}}

// Build creates a {{.Name}} from the builder
func (b *{{.Name}}Builder) Build() {{.Name}} {
    return {{.Name}}{
        {{range .Fields}}{{.Name}}: b.{{.Name}},
        {{end}}
    }
}

{{end}}
`

func generateBuilders(filename string) error {
    // Parse the Go file
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
    if err != nil {
        return err
    }
    
    var structs []StructInfo
    
    // Find all struct declarations
    ast.Inspect(node, func(n ast.Node) bool {
        switch x := n.(type) {
        case *ast.GenDecl:
            if x.Tok == token.TYPE {
                for _, spec := range x.Specs {
                    if typeSpec, ok := spec.(*ast.TypeSpec); ok {
                        if structType, ok := typeSpec.Type.(*ast.StructType); ok {
                            structInfo := StructInfo{
                                Name: typeSpec.Name.Name,
                            }
                            
                            // Extract fields
                            for _, field := range structType.Fields.List {
                                for _, name := range field.Names {
                                    if name.IsExported() {
                                        fieldInfo := FieldInfo{
                                            Name: name.Name,
                                            Type: getTypeString(field.Type),
                                        }
                                        structInfo.Fields = append(structInfo.Fields, fieldInfo)
                                    }
                                }
                            }
                            
                            structs = append(structs, structInfo)
                        }
                    }
                }
            }
        }
        return true
    })
    
    // Generate the builder code
    tmpl, err := template.New("builder").Parse(builderTemplate)
    if err != nil {
        return err
    }
    
    // Create output file
    outputFile := strings.TrimSuffix(filename, ".go") + "_builder.go"
    file, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer file.Close()
    
    data := struct {
        Package string
        Structs []StructInfo
    }{
        Package: node.Name.Name,
        Structs: structs,
    }
    
    return tmpl.Execute(file, data)
}

func getTypeString(expr ast.Expr) string {
    switch x := expr.(type) {
    case *ast.Ident:
        return x.Name
    case *ast.StarExpr:
        return "*" + getTypeString(x.X)
    case *ast.ArrayType:
        return "[]" + getTypeString(x.Elt)
    case *ast.MapType:
        return "map[" + getTypeString(x.Key) + "]" + getTypeString(x.Value)
    case *ast.SelectorExpr:
        return getTypeString(x.X) + "." + x.Sel.Name
    default:
        return "interface{}"
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: generator <file.go>")
        os.Exit(1)
    }
    
    filename := os.Args[1]
    if err := generateBuilders(filename); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Generated builder for %s\n", filename)
}
```

### Example Usage of Generated Code

```go
// models.go
package main

import "fmt"

type User struct {
    ID       int
    Name     string
    Email    string
    Age      int
    IsActive bool
}

type Product struct {
    ID          int
    Name        string
    Description string
    Price       float64
    InStock     bool
}

func main() {
    // After running the generator, we can use the builder pattern
    user := NewUserBuilder().
        ID(1).
        Name("John Doe").
        Email("john@example.com").
        Age(30).
        IsActive(true).
        Build()
    
    fmt.Printf("User: %+v\n", user)
    
    product := NewProductBuilder().
        ID(1).
        Name("Laptop").
        Description("High-performance laptop").
        Price(999.99).
        InStock(true).
        Build()
    
    fmt.Printf("Product: %+v\n", product)
}
```

## Performance Considerations

### Reflection Performance

```go
package main

import (
    "fmt"
    "reflect"
    "time"
)

type Person struct {
    Name string
    Age  int
}

// Direct access
func directAccess(p *Person, name string, age int) {
    p.Name = name
    p.Age = age
}

// Reflection access
func reflectionAccess(p *Person, name string, age int) {
    v := reflect.ValueOf(p).Elem()
    v.FieldByName("Name").SetString(name)
    v.FieldByName("Age").SetInt(int64(age))
}

// Cached reflection access
var (
    personType = reflect.TypeOf(Person{})
    nameField  = getFieldIndex(personType, "Name")
    ageField   = getFieldIndex(personType, "Age")
)

func getFieldIndex(t reflect.Type, name string) int {
    for i := 0; i < t.NumField(); i++ {
        if t.Field(i).Name == name {
            return i
        }
    }
    return -1
}

func cachedReflectionAccess(p *Person, name string, age int) {
    v := reflect.ValueOf(p).Elem()
    v.Field(nameField).SetString(name)
    v.Field(ageField).SetInt(int64(age))
}

func benchmark(name string, fn func(), iterations int) {
    start := time.Now()
    for i := 0; i < iterations; i++ {
        fn()
    }
    duration := time.Since(start)
    fmt.Printf("%s: %v (%.2f ns/op)\n", name, duration, float64(duration.Nanoseconds())/float64(iterations))
}

func main() {
    iterations := 1000000
    
    var p Person
    
    fmt.Printf("Performance comparison (%d iterations):\n", iterations)
    
    benchmark("Direct Access", func() {
        directAccess(&p, "John", 30)
    }, iterations)
    
    benchmark("Reflection Access", func() {
        reflectionAccess(&p, "John", 30)
    }, iterations)
    
    benchmark("Cached Reflection", func() {
        cachedReflectionAccess(&p, "John", 30)
    }, iterations)
}
```

## Best Practices

### When to Use Reflection

1. **Library and Framework Development**: When you need to work with types unknown at compile time
2. **Serialization/Deserialization**: Converting between different data formats
3. **Validation Frameworks**: Inspecting struct tags and validating fields
4. **Dependency Injection**: Automatically wiring dependencies
5. **Code Generation**: Analyzing code structure to generate boilerplate

### When NOT to Use Reflection

1. **Performance-Critical Code**: Reflection is slower than direct access
2. **Simple Type Conversions**: Use type assertions instead
3. **Known Types at Compile Time**: Use interfaces and type switches
4. **Basic Operations**: Don't over-engineer simple problems

### Reflection Guidelines

```go
package main

import (
    "fmt"
    "reflect"
)

// Good: Check for nil before using reflection
func safeReflection(obj interface{}) {
    if obj == nil {
        fmt.Println("Object is nil")
        return
    }
    
    v := reflect.ValueOf(obj)
    if !v.IsValid() {
        fmt.Println("Invalid value")
        return
    }
    
    // Handle pointer to nil
    if v.Kind() == reflect.Ptr && v.IsNil() {
        fmt.Println("Pointer to nil")
        return
    }
    
    fmt.Printf("Type: %v, Value: %v\n", v.Type(), v.Interface())
}

// Good: Cache reflection operations
type FieldCache struct {
    fields map[string]int
}

func NewFieldCache(t reflect.Type) *FieldCache {
    cache := &FieldCache{
        fields: make(map[string]int),
    }
    
    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        cache.fields[field.Name] = i
    }
    
    return cache
}

func (c *FieldCache) GetFieldIndex(name string) (int, bool) {
    index, ok := c.fields[name]
    return index, ok
}

// Good: Use type switches for known types
func processValue(v interface{}) {
    switch val := v.(type) {
    case string:
        fmt.Printf("String: %s\n", val)
    case int:
        fmt.Printf("Int: %d\n", val)
    case float64:
        fmt.Printf("Float: %f\n", val)
    default:
        // Only use reflection as a fallback
        rv := reflect.ValueOf(v)
        fmt.Printf("Unknown type %v: %v\n", rv.Type(), rv.Interface())
    }
}

func main() {
    // Test safe reflection
    safeReflection(nil)
    safeReflection("hello")
    safeReflection((*int)(nil))
    
    // Test field cache
    type Person struct {
        Name string
        Age  int
    }
    
    cache := NewFieldCache(reflect.TypeOf(Person{}))
    if index, ok := cache.GetFieldIndex("Name"); ok {
        fmt.Printf("Name field index: %d\n", index)
    }
    
    // Test type switches
    processValue("hello")
    processValue(42)
    processValue(3.14)
    processValue([]int{1, 2, 3})
}
```

## Key Takeaways

1. **Reflection is powerful but expensive**: Use it judiciously and consider performance implications.

2. **Type safety is reduced**: Reflection moves type checking from compile time to runtime.

3. **Cache reflection operations**: Store `reflect.Type` and field indices to improve performance.

4. **Handle edge cases**: Always check for nil values, invalid types, and unexported fields.

5. **Use interfaces when possible**: Prefer compile-time polymorphism over runtime reflection.

6. **Code generation is often better**: For known patterns, generate code at build time rather than using reflection at runtime.

7. **Reflection enables powerful libraries**: Many Go libraries (JSON, ORM, validation) rely on reflection.

8. **Document reflection usage**: Make it clear when and why reflection is being used.

9. **Test thoroughly**: Reflection code is harder to verify at compile time.

10. **Consider alternatives**: Type assertions, interfaces, and code generation might be better solutions.

## Next Steps

Now that you understand reflection and code generation in Go, let's explore [File I/O and System Programming](21-file-io-system.md) to learn about working with files, directories, and system resources!

---

**Previous**: [← Testing in Go](19-testing.md) | **Next**: [File I/O and System Programming →](21-file-io-system.md)