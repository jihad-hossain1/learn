# 30. WebAssembly

This chapter covers using Go with WebAssembly (WASM) to run Go code in web browsers, including compilation, JavaScript interop, and practical applications.

## Introduction to WebAssembly

WebAssembly (WASM) is a binary instruction format that allows code written in languages like Go, Rust, C++, and others to run in web browsers at near-native performance. Go has built-in support for compiling to WebAssembly, making it possible to run Go code directly in the browser.

### Benefits of WebAssembly

- **Performance**: Near-native execution speed
- **Language flexibility**: Use Go instead of JavaScript for complex logic
- **Code reuse**: Share code between server and client
- **Security**: Runs in a sandboxed environment
- **Portability**: Runs on any platform that supports WebAssembly

## Setting Up WebAssembly with Go

### Basic Setup

Go has built-in support for WebAssembly. You can compile Go code to WASM using the standard Go toolchain:

```bash
# Set environment variables for WebAssembly compilation
set GOOS=js
set GOARCH=wasm

# Compile Go code to WebAssembly
go build -o main.wasm main.go
```

### Simple "Hello World" Example

Let's start with a basic example:

```go
// main.go
package main

import (
    "fmt"
    "syscall/js"
)

func main() {
    fmt.Println("Hello, WebAssembly!")
    
    // Keep the program running
    select {}
}
```

Compile it to WebAssembly:

```bash
set GOOS=js
set GOARCH=wasm
go build -o main.wasm main.go
```

### HTML Setup

Create an HTML file to load and run the WebAssembly module:

```html
<!-- index.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Go WebAssembly Example</title>
</head>
<body>
    <h1>Go WebAssembly Example</h1>
    <div id="output"></div>
    
    <script src="wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</body>
</html>
```

You'll need the `wasm_exec.js` file from the Go installation:

```bash
# Copy the WebAssembly support file
copy "%GOROOT%\misc\wasm\wasm_exec.js" .
```

### Serving the Files

WebAssembly requires files to be served over HTTP (not file://). You can use Go's built-in HTTP server:

```go
// server.go
package main

import (
    "log"
    "net/http"
)

func main() {
    fs := http.FileServer(http.Dir("."))
    http.Handle("/", fs)
    
    log.Println("Server starting on http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## JavaScript Interoperability

### Calling JavaScript from Go

The `syscall/js` package provides the interface for calling JavaScript from Go:

```go
// js_interop.go
package main

import (
    "fmt"
    "syscall/js"
)

func main() {
    // Get the global object (window in browsers)
    global := js.Global()
    
    // Get the document object
    document := global.Get("document")
    
    // Get an element by ID
    output := document.Call("getElementById", "output")
    
    // Set the innerHTML
    output.Set("innerHTML", "<p>Hello from Go!</p>")
    
    // Call JavaScript functions
    global.Call("alert", "Hello from Go WebAssembly!")
    
    // Access JavaScript objects
    console := global.Get("console")
    console.Call("log", "This is logged from Go")
    
    // Work with JavaScript arrays
    jsArray := js.Global().Get("Array").New()
    jsArray.Call("push", "item1", "item2", "item3")
    fmt.Printf("Array length: %d\n", jsArray.Get("length").Int())
    
    // Keep the program running
    select {}
}
```

### Exposing Go Functions to JavaScript

You can expose Go functions to JavaScript using `js.FuncOf`:

```go
// exposed_functions.go
package main

import (
    "fmt"
    "strconv"
    "syscall/js"
)

// Add function that can be called from JavaScript
func add(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: add requires exactly 2 arguments")
    }
    
    a := args[0].Float()
    b := args[1].Float()
    result := a + b
    
    fmt.Printf("Go: Adding %.2f + %.2f = %.2f\n", a, b, result)
    return js.ValueOf(result)
}

// Fibonacci function
func fibonacci(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: fibonacci requires exactly 1 argument")
    }
    
    n := args[0].Int()
    result := fib(n)
    
    fmt.Printf("Go: Fibonacci(%d) = %d\n", n, result)
    return js.ValueOf(result)
}

func fib(n int) int {
    if n <= 1 {
        return n
    }
    return fib(n-1) + fib(n-2)
}

// String manipulation function
func reverseString(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: reverseString requires exactly 1 argument")
    }
    
    input := args[0].String()
    runes := []rune(input)
    
    // Reverse the string
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    
    result := string(runes)
    fmt.Printf("Go: Reversing '%s' = '%s'\n", input, result)
    return js.ValueOf(result)
}

// Process array function
func processArray(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: processArray requires exactly 1 argument")
    }
    
    jsArray := args[0]
    length := jsArray.Get("length").Int()
    
    // Convert JavaScript array to Go slice
    goSlice := make([]float64, length)
    for i := 0; i < length; i++ {
        goSlice[i] = jsArray.Index(i).Float()
    }
    
    // Process the array (calculate sum and average)
    sum := 0.0
    for _, v := range goSlice {
        sum += v
    }
    average := sum / float64(length)
    
    // Create result object
    result := js.Global().Get("Object").New()
    result.Set("sum", js.ValueOf(sum))
    result.Set("average", js.ValueOf(average))
    result.Set("count", js.ValueOf(length))
    
    fmt.Printf("Go: Processed array of %d elements, sum=%.2f, avg=%.2f\n", length, sum, average)
    return result
}

func main() {
    fmt.Println("Go WebAssembly Initialized")
    
    // Expose functions to JavaScript
    js.Global().Set("goAdd", js.FuncOf(add))
    js.Global().Set("goFibonacci", js.FuncOf(fibonacci))
    js.Global().Set("goReverseString", js.FuncOf(reverseString))
    js.Global().Set("goProcessArray", js.FuncOf(processArray))
    
    // Signal that Go is ready
    js.Global().Call("goReady")
    
    // Keep the program running
    select {}
}
```

Update the HTML to use these functions:

```html
<!-- index.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Go WebAssembly Functions</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ccc; }
        button { margin: 5px; padding: 10px; }
        input { margin: 5px; padding: 5px; }
        #output { background: #f5f5f5; padding: 10px; margin: 10px 0; }
    </style>
</head>
<body>
    <h1>Go WebAssembly Functions</h1>
    
    <div class="section">
        <h3>Addition</h3>
        <input type="number" id="num1" placeholder="Number 1" value="5">
        <input type="number" id="num2" placeholder="Number 2" value="3">
        <button onclick="testAdd()">Add Numbers</button>
    </div>
    
    <div class="section">
        <h3>Fibonacci</h3>
        <input type="number" id="fibNum" placeholder="Fibonacci number" value="10">
        <button onclick="testFibonacci()">Calculate Fibonacci</button>
    </div>
    
    <div class="section">
        <h3>String Reversal</h3>
        <input type="text" id="stringInput" placeholder="Enter text" value="Hello WebAssembly">
        <button onclick="testReverseString()">Reverse String</button>
    </div>
    
    <div class="section">
        <h3>Array Processing</h3>
        <input type="text" id="arrayInput" placeholder="Enter numbers (comma-separated)" value="1,2,3,4,5">
        <button onclick="testProcessArray()">Process Array</button>
    </div>
    
    <div id="output"></div>
    
    <script src="wasm_exec.js"></script>
    <script>
        let goReady = false;
        
        // Called when Go is ready
        function goReady() {
            goReady = true;
            document.getElementById('output').innerHTML += '<p><strong>Go WebAssembly is ready!</strong></p>';
        }
        
        function testAdd() {
            if (!goReady) {
                alert('Go WebAssembly is not ready yet');
                return;
            }
            
            const num1 = parseFloat(document.getElementById('num1').value);
            const num2 = parseFloat(document.getElementById('num2').value);
            
            const result = goAdd(num1, num2);
            document.getElementById('output').innerHTML += `<p>Add: ${num1} + ${num2} = ${result}</p>`;
        }
        
        function testFibonacci() {
            if (!goReady) {
                alert('Go WebAssembly is not ready yet');
                return;
            }
            
            const num = parseInt(document.getElementById('fibNum').value);
            const start = performance.now();
            const result = goFibonacci(num);
            const end = performance.now();
            
            document.getElementById('output').innerHTML += 
                `<p>Fibonacci(${num}) = ${result} (calculated in ${(end - start).toFixed(2)}ms)</p>`;
        }
        
        function testReverseString() {
            if (!goReady) {
                alert('Go WebAssembly is not ready yet');
                return;
            }
            
            const input = document.getElementById('stringInput').value;
            const result = goReverseString(input);
            document.getElementById('output').innerHTML += `<p>Reverse: "${input}" → "${result}"</p>`;
        }
        
        function testProcessArray() {
            if (!goReady) {
                alert('Go WebAssembly is not ready yet');
                return;
            }
            
            const input = document.getElementById('arrayInput').value;
            const numbers = input.split(',').map(s => parseFloat(s.trim())).filter(n => !isNaN(n));
            
            const result = goProcessArray(numbers);
            document.getElementById('output').innerHTML += 
                `<p>Array [${numbers.join(', ')}]: Sum=${result.sum}, Average=${result.average.toFixed(2)}, Count=${result.count}</p>`;
        }
        
        // Load and run the WebAssembly module
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        }).catch((err) => {
            console.error('Failed to load WebAssembly:', err);
            document.getElementById('output').innerHTML += '<p><strong>Error loading WebAssembly</strong></p>';
        });
    </script>
</body>
</html>
```

## Working with DOM

You can manipulate the DOM directly from Go:

```go
// dom_manipulation.go
package main

import (
    "fmt"
    "strconv"
    "syscall/js"
    "time"
)

// Create a new element
func createElement(tag, text string) js.Value {
    document := js.Global().Get("document")
    element := document.Call("createElement", tag)
    element.Set("textContent", text)
    return element
}

// Add event listener
func addEventListener(element js.Value, event string, callback js.Func) {
    element.Call("addEventListener", event, callback)
}

// Dynamic content generation
func generateContent(this js.Value, args []js.Value) interface{} {
    document := js.Global().Get("document")
    container := document.Call("getElementById", "dynamic-content")
    
    // Clear existing content
    container.Set("innerHTML", "")
    
    // Create a list of items
    ul := createElement("ul", "")
    
    for i := 1; i <= 5; i++ {
        li := createElement("li", fmt.Sprintf("Dynamic item %d - Generated at %s", i, time.Now().Format("15:04:05")))
        
        // Add click event to each item
        clickHandler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
            this.Get("style").Set("color", "red")
            this.Set("textContent", this.Get("textContent").String()+" (clicked!)")
            return nil
        })
        
        addEventListener(li, "click", clickHandler)
        ul.Call("appendChild", li)
    }
    
    container.Call("appendChild", ul)
    return nil
}

// Form handling
func handleForm(this js.Value, args []js.Value) interface{} {
    event := args[0]
    event.Call("preventDefault") // Prevent form submission
    
    document := js.Global().Get("document")
    
    // Get form values
    name := document.Call("getElementById", "name").Get("value").String()
    email := document.Call("getElementById", "email").Get("value").String()
    age := document.Call("getElementById", "age").Get("value").String()
    
    // Validate
    if name == "" || email == "" || age == "" {
        js.Global().Call("alert", "Please fill in all fields")
        return nil
    }
    
    ageInt, err := strconv.Atoi(age)
    if err != nil || ageInt < 0 || ageInt > 150 {
        js.Global().Call("alert", "Please enter a valid age")
        return nil
    }
    
    // Display result
    result := document.Call("getElementById", "form-result")
    result.Set("innerHTML", fmt.Sprintf(
        "<h3>Form Submitted Successfully!</h3><p>Name: %s<br>Email: %s<br>Age: %d</p>",
        name, email, ageInt,
    ))
    
    // Clear form
    document.Call("getElementById", "user-form").Call("reset")
    
    return nil
}

// Canvas drawing
func drawOnCanvas(this js.Value, args []js.Value) interface{} {
    document := js.Global().Get("document")
    canvas := document.Call("getElementById", "canvas")
    ctx := canvas.Call("getContext", "2d")
    
    // Clear canvas
    ctx.Call("clearRect", 0, 0, 400, 300)
    
    // Draw background
    ctx.Set("fillStyle", "#f0f0f0")
    ctx.Call("fillRect", 0, 0, 400, 300)
    
    // Draw some shapes
    ctx.Set("fillStyle", "#ff6b6b")
    ctx.Call("fillRect", 50, 50, 100, 80)
    
    ctx.Set("fillStyle", "#4ecdc4")
    ctx.Call("beginPath")
    ctx.Call("arc", 250, 90, 40, 0, 2*3.14159)
    ctx.Call("fill")
    
    // Draw text
    ctx.Set("fillStyle", "#333")
    ctx.Set("font", "20px Arial")
    ctx.Call("fillText", "Hello from Go!", 50, 200)
    
    // Draw current time
    ctx.Set("font", "16px Arial")
    ctx.Call("fillText", "Time: "+time.Now().Format("15:04:05"), 50, 250)
    
    return nil
}

// Animation loop
func startAnimation(this js.Value, args []js.Value) interface{} {
    go animationLoop()
    return nil
}

func animationLoop() {
    document := js.Global().Get("document")
    canvas := document.Call("getElementById", "animation-canvas")
    ctx := canvas.Call("getContext", "2d")
    
    x := 0.0
    direction := 1.0
    
    for {
        // Clear canvas
        ctx.Call("clearRect", 0, 0, 400, 200)
        
        // Draw background
        ctx.Set("fillStyle", "#e8f4f8")
        ctx.Call("fillRect", 0, 0, 400, 200)
        
        // Draw moving circle
        ctx.Set("fillStyle", "#ff6b6b")
        ctx.Call("beginPath")
        ctx.Call("arc", x, 100, 20, 0, 2*3.14159)
        ctx.Call("fill")
        
        // Update position
        x += direction * 2
        if x >= 380 || x <= 20 {
            direction *= -1
        }
        
        // Wait for next frame
        time.Sleep(16 * time.Millisecond) // ~60 FPS
    }
}

func main() {
    fmt.Println("DOM Manipulation WebAssembly Initialized")
    
    // Expose functions to JavaScript
    js.Global().Set("goGenerateContent", js.FuncOf(generateContent))
    js.Global().Set("goHandleForm", js.FuncOf(handleForm))
    js.Global().Set("goDrawOnCanvas", js.FuncOf(drawOnCanvas))
    js.Global().Set("goStartAnimation", js.FuncOf(startAnimation))
    
    // Set up form event listener
    document := js.Global().Get("document")
    document.Call("addEventListener", "DOMContentLoaded", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        form := document.Call("getElementById", "user-form")
        if !form.IsNull() {
            form.Call("addEventListener", "submit", js.FuncOf(handleForm))
        }
        return nil
    }))
    
    // Signal that Go is ready
    js.Global().Call("goReady")
    
    // Keep the program running
    select {}
}
```

HTML for DOM manipulation:

```html
<!-- dom.html -->
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Go WebAssembly DOM Manipulation</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .section { margin: 20px 0; padding: 15px; border: 1px solid #ccc; }
        button { margin: 5px; padding: 10px; background: #007bff; color: white; border: none; cursor: pointer; }
        button:hover { background: #0056b3; }
        input { margin: 5px; padding: 8px; border: 1px solid #ccc; }
        canvas { border: 1px solid #ccc; }
        li { cursor: pointer; padding: 5px; }
        li:hover { background: #f0f0f0; }
    </style>
</head>
<body>
    <h1>Go WebAssembly DOM Manipulation</h1>
    
    <div class="section">
        <h3>Dynamic Content Generation</h3>
        <button onclick="goGenerateContent()">Generate Dynamic Content</button>
        <div id="dynamic-content"></div>
    </div>
    
    <div class="section">
        <h3>Form Handling</h3>
        <form id="user-form">
            <input type="text" id="name" placeholder="Name" required>
            <input type="email" id="email" placeholder="Email" required>
            <input type="number" id="age" placeholder="Age" min="0" max="150" required>
            <button type="submit">Submit Form</button>
        </form>
        <div id="form-result"></div>
    </div>
    
    <div class="section">
        <h3>Canvas Drawing</h3>
        <button onclick="goDrawOnCanvas()">Draw on Canvas</button><br>
        <canvas id="canvas" width="400" height="300"></canvas>
    </div>
    
    <div class="section">
        <h3>Animation</h3>
        <button onclick="goStartAnimation()">Start Animation</button><br>
        <canvas id="animation-canvas" width="400" height="200"></canvas>
    </div>
    
    <script src="wasm_exec.js"></script>
    <script>
        function goReady() {
            console.log('Go WebAssembly is ready!');
        }
        
        // Load and run the WebAssembly module
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        }).catch((err) => {
            console.error('Failed to load WebAssembly:', err);
        });
    </script>
</body>
</html>
```

## Advanced WebAssembly Features

### Working with Binary Data

```go
// binary_data.go
package main

import (
    "bytes"
    "compress/gzip"
    "crypto/md5"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "fmt"
    "io"
    "syscall/js"
)

// Compress data using gzip
func compressData(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: compressData requires exactly 1 argument")
    }
    
    input := args[0].String()
    
    var buf bytes.Buffer
    gzipWriter := gzip.NewWriter(&buf)
    
    _, err := gzipWriter.Write([]byte(input))
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    err = gzipWriter.Close()
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    compressed := buf.Bytes()
    encoded := base64.StdEncoding.EncodeToString(compressed)
    
    result := js.Global().Get("Object").New()
    result.Set("original", js.ValueOf(input))
    result.Set("compressed", js.ValueOf(encoded))
    result.Set("originalSize", js.ValueOf(len(input)))
    result.Set("compressedSize", js.ValueOf(len(compressed)))
    result.Set("compressionRatio", js.ValueOf(float64(len(compressed))/float64(len(input))))
    
    return result
}

// Decompress gzip data
func decompressData(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: decompressData requires exactly 1 argument")
    }
    
    encoded := args[0].String()
    
    compressed, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    reader, err := gzip.NewReader(bytes.NewReader(compressed))
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    defer reader.Close()
    
    decompressed, err := io.ReadAll(reader)
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    return js.ValueOf(string(decompressed))
}

// Calculate hash
func calculateHash(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: calculateHash requires exactly 2 arguments (data, algorithm)")
    }
    
    data := args[0].String()
    algorithm := args[1].String()
    
    var hash string
    
    switch algorithm {
    case "md5":
        h := md5.Sum([]byte(data))
        hash = hex.EncodeToString(h[:])
    case "sha256":
        h := sha256.Sum256([]byte(data))
        hash = hex.EncodeToString(h[:])
    default:
        return js.ValueOf("Error: Unsupported algorithm. Use 'md5' or 'sha256'")
    }
    
    result := js.Global().Get("Object").New()
    result.Set("data", js.ValueOf(data))
    result.Set("algorithm", js.ValueOf(algorithm))
    result.Set("hash", js.ValueOf(hash))
    
    return result
}

// Base64 encode/decode
func base64Encode(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: base64Encode requires exactly 1 argument")
    }
    
    input := args[0].String()
    encoded := base64.StdEncoding.EncodeToString([]byte(input))
    
    return js.ValueOf(encoded)
}

func base64Decode(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: base64Decode requires exactly 1 argument")
    }
    
    encoded := args[0].String()
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        return js.ValueOf("Error: " + err.Error())
    }
    
    return js.ValueOf(string(decoded))
}

// Process image data (simple example)
func processImageData(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: processImageData requires exactly 1 argument")
    }
    
    // Get ImageData from canvas
    imageData := args[0]
    data := imageData.Get("data")
    width := imageData.Get("width").Int()
    height := imageData.Get("height").Int()
    
    // Convert to grayscale
    for i := 0; i < width*height*4; i += 4 {
        r := data.Index(i).Int()
        g := data.Index(i + 1).Int()
        b := data.Index(i + 2).Int()
        
        // Calculate grayscale value
        gray := int(0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b))
        
        // Set RGB to grayscale value
        data.SetIndex(i, js.ValueOf(gray))
        data.SetIndex(i+1, js.ValueOf(gray))
        data.SetIndex(i+2, js.ValueOf(gray))
        // Alpha channel (i+3) remains unchanged
    }
    
    return imageData
}

func main() {
    fmt.Println("Binary Data Processing WebAssembly Initialized")
    
    // Expose functions to JavaScript
    js.Global().Set("goCompressData", js.FuncOf(compressData))
    js.Global().Set("goDecompressData", js.FuncOf(decompressData))
    js.Global().Set("goCalculateHash", js.FuncOf(calculateHash))
    js.Global().Set("goBase64Encode", js.FuncOf(base64Encode))
    js.Global().Set("goBase64Decode", js.FuncOf(base64Decode))
    js.Global().Set("goProcessImageData", js.FuncOf(processImageData))
    
    // Signal that Go is ready
    js.Global().Call("goReady")
    
    // Keep the program running
    select {}
}
```

### Performance Optimization

```go
// performance.go
package main

import (
    "fmt"
    "math"
    "runtime"
    "syscall/js"
    "time"
)

// CPU-intensive calculation
func calculatePrimes(this js.Value, args []js.Value) interface{} {
    if len(args) != 1 {
        return js.ValueOf("Error: calculatePrimes requires exactly 1 argument")
    }
    
    limit := args[0].Int()
    start := time.Now()
    
    primes := sieveOfEratosthenes(limit)
    
    duration := time.Since(start)
    
    result := js.Global().Get("Object").New()
    result.Set("primes", js.ValueOf(len(primes)))
    result.Set("limit", js.ValueOf(limit))
    result.Set("duration", js.ValueOf(duration.Milliseconds()))
    result.Set("primesArray", js.ValueOf(primes[:min(100, len(primes))])) // Return first 100 primes
    
    return result
}

func sieveOfEratosthenes(limit int) []int {
    if limit < 2 {
        return []int{}
    }
    
    isPrime := make([]bool, limit+1)
    for i := 2; i <= limit; i++ {
        isPrime[i] = true
    }
    
    for i := 2; i*i <= limit; i++ {
        if isPrime[i] {
            for j := i * i; j <= limit; j += i {
                isPrime[j] = false
            }
        }
    }
    
    var primes []int
    for i := 2; i <= limit; i++ {
        if isPrime[i] {
            primes = append(primes, i)
        }
    }
    
    return primes
}

// Matrix multiplication
func multiplyMatrices(this js.Value, args []js.Value) interface{} {
    if len(args) != 2 {
        return js.ValueOf("Error: multiplyMatrices requires exactly 2 arguments")
    }
    
    matrixA := args[0]
    matrixB := args[1]
    
    rowsA := matrixA.Get("length").Int()
    colsA := matrixA.Index(0).Get("length").Int()
    rowsB := matrixB.Get("length").Int()
    colsB := matrixB.Index(0).Get("length").Int()
    
    if colsA != rowsB {
        return js.ValueOf("Error: Matrix dimensions don't match for multiplication")
    }
    
    start := time.Now()
    
    // Convert JavaScript arrays to Go slices
    a := make([][]float64, rowsA)
    for i := 0; i < rowsA; i++ {
        a[i] = make([]float64, colsA)
        for j := 0; j < colsA; j++ {
            a[i][j] = matrixA.Index(i).Index(j).Float()
        }
    }
    
    b := make([][]float64, rowsB)
    for i := 0; i < rowsB; i++ {
        b[i] = make([]float64, colsB)
        for j := 0; j < colsB; j++ {
            b[i][j] = matrixB.Index(i).Index(j).Float()
        }
    }
    
    // Multiply matrices
    result := make([][]float64, rowsA)
    for i := 0; i < rowsA; i++ {
        result[i] = make([]float64, colsB)
        for j := 0; j < colsB; j++ {
            for k := 0; k < colsA; k++ {
                result[i][j] += a[i][k] * b[k][j]
            }
        }
    }
    
    duration := time.Since(start)
    
    // Convert result back to JavaScript array
    jsResult := js.Global().Get("Array").New()
    for i := 0; i < rowsA; i++ {
        row := js.Global().Get("Array").New()
        for j := 0; j < colsB; j++ {
            row.Call("push", js.ValueOf(result[i][j]))
        }
        jsResult.Call("push", row)
    }
    
    response := js.Global().Get("Object").New()
    response.Set("result", jsResult)
    response.Set("duration", js.ValueOf(duration.Milliseconds()))
    response.Set("dimensions", js.ValueOf(fmt.Sprintf("%dx%d × %dx%d = %dx%d", rowsA, colsA, rowsB, colsB, rowsA, colsB)))
    
    return response
}

// Mandelbrot set calculation
func calculateMandelbrot(this js.Value, args []js.Value) interface{} {
    if len(args) != 4 {
        return js.ValueOf("Error: calculateMandelbrot requires exactly 4 arguments (width, height, maxIter, zoom)")
    }
    
    width := args[0].Int()
    height := args[1].Int()
    maxIter := args[2].Int()
    zoom := args[3].Float()
    
    start := time.Now()
    
    data := make([]int, width*height)
    
    for y := 0; y < height; y++ {
        for x := 0; x < width; x++ {
            // Map pixel to complex plane
            real := (float64(x)/float64(width)-0.5)*4.0/zoom - 0.5
            imag := (float64(y)/float64(height)-0.5)*4.0/zoom
            
            // Calculate Mandelbrot iteration
            iter := mandelbrotIteration(complex(real, imag), maxIter)
            data[y*width+x] = iter
        }
    }
    
    duration := time.Since(start)
    
    // Convert to JavaScript array
    jsData := js.Global().Get("Array").New()
    for _, value := range data {
        jsData.Call("push", js.ValueOf(value))
    }
    
    result := js.Global().Get("Object").New()
    result.Set("data", jsData)
    result.Set("width", js.ValueOf(width))
    result.Set("height", js.ValueOf(height))
    result.Set("duration", js.ValueOf(duration.Milliseconds()))
    
    return result
}

func mandelbrotIteration(c complex128, maxIter int) int {
    z := complex(0, 0)
    for i := 0; i < maxIter; i++ {
        if real(z)*real(z)+imag(z)*imag(z) > 4 {
            return i
        }
        z = z*z + c
    }
    return maxIter
}

// Memory usage information
func getMemoryStats(this js.Value, args []js.Value) interface{} {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    result := js.Global().Get("Object").New()
    result.Set("allocatedMB", js.ValueOf(float64(m.Alloc)/1024/1024))
    result.Set("totalAllocatedMB", js.ValueOf(float64(m.TotalAlloc)/1024/1024))
    result.Set("systemMB", js.ValueOf(float64(m.Sys)/1024/1024))
    result.Set("numGC", js.ValueOf(m.NumGC))
    result.Set("goroutines", js.ValueOf(runtime.NumGoroutine()))
    
    return result
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    fmt.Println("Performance WebAssembly Initialized")
    
    // Expose functions to JavaScript
    js.Global().Set("goCalculatePrimes", js.FuncOf(calculatePrimes))
    js.Global().Set("goMultiplyMatrices", js.FuncOf(multiplyMatrices))
    js.Global().Set("goCalculateMandelbrot", js.FuncOf(calculateMandelbrot))
    js.Global().Set("goGetMemoryStats", js.FuncOf(getMemoryStats))
    
    // Signal that Go is ready
    js.Global().Call("goReady")
    
    // Keep the program running
    select {}
}
```

## Best Practices

### 1. Memory Management

```go
// Always clean up JavaScript function references
func cleanupExample() {
    callback := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
        // Function implementation
        return nil
    })
    defer callback.Release() // Important: release the function
    
    // Use the callback
    js.Global().Get("someAPI").Call("addEventListener", "event", callback)
}
```

### 2. Error Handling

```go
func safeFunction(this js.Value, args []js.Value) interface{} {
    defer func() {
        if r := recover(); r != nil {
            fmt.Printf("Recovered from panic: %v\n", r)
        }
    }()
    
    // Your function implementation
    if len(args) == 0 {
        return js.ValueOf(map[string]interface{}{
            "error": "No arguments provided",
        })
    }
    
    // Process arguments safely
    return js.ValueOf(map[string]interface{}{
        "success": true,
        "result":  "Operation completed",
    })
}
```

### 3. Performance Considerations

- Minimize JavaScript ↔ Go calls
- Use batch operations when possible
- Avoid frequent memory allocations
- Use goroutines for concurrent operations
- Profile your WebAssembly code

### 4. Build Optimization

```bash
# Optimize for size
set GOOS=js
set GOARCH=wasm
go build -ldflags="-s -w" -o main.wasm main.go

# Further compression with tools like wasm-opt
# wasm-opt -Oz main.wasm -o main.optimized.wasm
```

## Exercises

### Exercise 1: Calculator App
Build a scientific calculator using Go WebAssembly with:
- Basic arithmetic operations
- Advanced functions (sin, cos, log, etc.)
- Memory functions
- History of calculations

### Exercise 2: Image Processing Tool
Create an image processing application that can:
- Load images from file input
- Apply filters (blur, sharpen, grayscale)
- Adjust brightness and contrast
- Save processed images

### Exercise 3: Data Visualization
Build a data visualization tool that:
- Parses CSV data
- Generates charts using Canvas API
- Supports different chart types
- Allows interactive data exploration

### Exercise 4: Game Development
Create a simple game using WebAssembly:
- Real-time graphics using Canvas
- User input handling
- Game state management
- Score tracking and persistence

## Key Takeaways

- WebAssembly allows Go code to run in browsers at near-native speed
- The `syscall/js` package provides JavaScript interoperability
- Go functions can be exposed to JavaScript using `js.FuncOf`
- DOM manipulation is possible directly from Go
- Memory management and error handling are crucial
- Performance optimization requires careful consideration of Go ↔ JS boundaries
- WebAssembly is ideal for CPU-intensive tasks in the browser
- Proper build optimization can significantly reduce WASM file size

## Next Steps

Next, we'll explore [Advanced Topics](31-advanced-topics.md) covering more specialized Go features and advanced programming patterns.