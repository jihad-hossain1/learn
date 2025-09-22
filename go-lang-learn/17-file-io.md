# 17. File I/O and System Programming

Go provides excellent support for file operations and system programming through its standard library. This chapter covers file I/O operations, directory manipulation, and system-level programming concepts.

## Basic File Operations

### Opening and Closing Files

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Open a file for reading
    file, err := os.Open("example.txt")
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close() // Always close files
    
    fmt.Println("File opened successfully")
    
    // Get file information
    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Printf("Error getting file info: %v\n", err)
        return
    }
    
    fmt.Printf("File name: %s\n", fileInfo.Name())
    fmt.Printf("File size: %d bytes\n", fileInfo.Size())
    fmt.Printf("File mode: %v\n", fileInfo.Mode())
    fmt.Printf("Modified time: %v\n", fileInfo.ModTime())
    fmt.Printf("Is directory: %v\n", fileInfo.IsDir())
}
```

### Creating Files

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // Create a new file
    file, err := os.Create("newfile.txt")
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()
    
    fmt.Println("File created successfully")
    
    // Create file with specific permissions
    file2, err := os.OpenFile("permfile.txt", os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Error creating file with permissions: %v\n", err)
        return
    }
    defer file2.Close()
    
    fmt.Println("File with permissions created successfully")
}
```

### File Opening Modes

```go
package main

import (
    "fmt"
    "os"
)

func demonstrateFileModes() {
    // Different file opening modes
    
    // Read only
    file1, err := os.OpenFile("example.txt", os.O_RDONLY, 0)
    if err != nil {
        fmt.Printf("Error opening for read: %v\n", err)
    } else {
        defer file1.Close()
        fmt.Println("Opened for reading")
    }
    
    // Write only (truncate if exists)
    file2, err := os.OpenFile("example.txt", os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        fmt.Printf("Error opening for write: %v\n", err)
    } else {
        defer file2.Close()
        fmt.Println("Opened for writing (truncate)")
    }
    
    // Append mode
    file3, err := os.OpenFile("example.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        fmt.Printf("Error opening for append: %v\n", err)
    } else {
        defer file3.Close()
        fmt.Println("Opened for appending")
    }
    
    // Read and write
    file4, err := os.OpenFile("example.txt", os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        fmt.Printf("Error opening for read/write: %v\n", err)
    } else {
        defer file4.Close()
        fmt.Println("Opened for read/write")
    }
}
```

## Reading Files

### Reading Entire File

```go
package main

import (
    "fmt"
    "io"
    "os"
)

// Method 1: Using os.ReadFile (Go 1.16+)
func readFileMethod1(filename string) {
    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        return
    }
    
    fmt.Printf("File contents:\n%s\n", string(data))
}

// Method 2: Using io.ReadAll
func readFileMethod2(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    
    data, err := io.ReadAll(file)
    if err != nil {
        fmt.Printf("Error reading file: %v\n", err)
        return
    }
    
    fmt.Printf("File contents:\n%s\n", string(data))
}

// Method 3: Reading in chunks
func readFileInChunks(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    
    buffer := make([]byte, 1024) // 1KB buffer
    
    for {
        n, err := file.Read(buffer)
        if err != nil {
            if err == io.EOF {
                break // End of file
            }
            fmt.Printf("Error reading file: %v\n", err)
            return
        }
        
        fmt.Printf("Read %d bytes: %s\n", n, string(buffer[:n]))
    }
}

func main() {
    filename := "example.txt"
    
    fmt.Println("=== Method 1: os.ReadFile ===")
    readFileMethod1(filename)
    
    fmt.Println("\n=== Method 2: io.ReadAll ===")
    readFileMethod2(filename)
    
    fmt.Println("\n=== Method 3: Reading in chunks ===")
    readFileInChunks(filename)
}
```

### Reading Line by Line

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func readFileLineByLine(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    lineNumber := 1
    
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Printf("Line %d: %s\n", lineNumber, line)
        lineNumber++
    }
    
    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file: %v\n", err)
    }
}

// Reading with custom delimiter
func readWithCustomDelimiter(filename string, delimiter byte) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Error opening file: %v\n", err)
        return
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    
    // Custom split function
    scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
        for i := 0; i < len(data); i++ {
            if data[i] == delimiter {
                return i + 1, data[:i], nil
            }
        }
        if atEOF && len(data) > 0 {
            return len(data), data, nil
        }
        return 0, nil, nil
    })
    
    tokenNumber := 1
    for scanner.Scan() {
        token := scanner.Text()
        fmt.Printf("Token %d: %s\n", tokenNumber, token)
        tokenNumber++
    }
    
    if err := scanner.Err(); err != nil {
        fmt.Printf("Error reading file: %v\n", err)
    }
}

func main() {
    filename := "example.txt"
    
    fmt.Println("=== Reading line by line ===")
    readFileLineByLine(filename)
    
    fmt.Println("\n=== Reading with custom delimiter (comma) ===")
    readWithCustomDelimiter("data.csv", ',')
}
```

## Writing Files

### Basic Writing

```go
package main

import (
    "fmt"
    "os"
)

// Method 1: Using os.WriteFile (Go 1.16+)
func writeFileMethod1(filename string, data []byte) {
    err := os.WriteFile(filename, data, 0644)
    if err != nil {
        fmt.Printf("Error writing file: %v\n", err)
        return
    }
    
    fmt.Println("File written successfully using os.WriteFile")
}

// Method 2: Using file.Write
func writeFileMethod2(filename string, data []byte) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()
    
    n, err := file.Write(data)
    if err != nil {
        fmt.Printf("Error writing to file: %v\n", err)
        return
    }
    
    fmt.Printf("Written %d bytes using file.Write\n", n)
}

// Method 3: Using file.WriteString
func writeStringToFile(filename string, text string) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()
    
    n, err := file.WriteString(text)
    if err != nil {
        fmt.Printf("Error writing string to file: %v\n", err)
        return
    }
    
    fmt.Printf("Written %d bytes using WriteString\n", n)
}

func main() {
    data := []byte("Hello, World!\nThis is a test file.\n")
    text := "This is a string written to file.\n"
    
    writeFileMethod1("output1.txt", data)
    writeFileMethod2("output2.txt", data)
    writeStringToFile("output3.txt", text)
}
```

### Buffered Writing

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func bufferedWriting(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()
    
    writer := bufio.NewWriter(file)
    defer writer.Flush() // Important: flush the buffer
    
    // Write multiple lines
    lines := []string{
        "Line 1: Hello, World!",
        "Line 2: This is buffered writing.",
        "Line 3: Much more efficient for multiple writes.",
        "Line 4: Don't forget to flush!",
    }
    
    for i, line := range lines {
        _, err := writer.WriteString(fmt.Sprintf("%s\n", line))
        if err != nil {
            fmt.Printf("Error writing line %d: %v\n", i+1, err)
            return
        }
    }
    
    fmt.Println("Buffered writing completed")
}

// Writing large amounts of data efficiently
func writeLotsOfData(filename string) {
    file, err := os.Create(filename)
    if err != nil {
        fmt.Printf("Error creating file: %v\n", err)
        return
    }
    defer file.Close()
    
    writer := bufio.NewWriterSize(file, 8192) // 8KB buffer
    defer writer.Flush()
    
    // Write 1000 lines
    for i := 1; i <= 1000; i++ {
        line := fmt.Sprintf("This is line number %d\n", i)
        _, err := writer.WriteString(line)
        if err != nil {
            fmt.Printf("Error writing line %d: %v\n", i, err)
            return
        }
        
        // Flush every 100 lines to ensure data is written
        if i%100 == 0 {
            writer.Flush()
        }
    }
    
    fmt.Println("Large file writing completed")
}

func main() {
    bufferedWriting("buffered_output.txt")
    writeLotsOfData("large_output.txt")
}
```

### Appending to Files

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func appendToFile(filename string, text string) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Error opening file for append: %v\n", err)
        return
    }
    defer file.Close()
    
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] %s\n", timestamp, text)
    
    _, err = file.WriteString(logEntry)
    if err != nil {
        fmt.Printf("Error appending to file: %v\n", err)
        return
    }
    
    fmt.Println("Text appended successfully")
}

// Simple logging function
func logMessage(filename string, level string, message string) {
    file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        fmt.Printf("Error opening log file: %v\n", err)
        return
    }
    defer file.Close()
    
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    logEntry := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)
    
    file.WriteString(logEntry)
}

func main() {
    logFile := "app.log"
    
    // Append some log entries
    logMessage(logFile, "INFO", "Application started")
    logMessage(logFile, "DEBUG", "Processing user request")
    logMessage(logFile, "ERROR", "Database connection failed")
    logMessage(logFile, "INFO", "Application stopped")
    
    // Append custom text
    appendToFile("notes.txt", "This is a new note")
    appendToFile("notes.txt", "This is another note")
    
    fmt.Println("Logging completed")
}
```

## Working with Directories

### Creating and Removing Directories

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
)

func directoryOperations() {
    // Create a single directory
    err := os.Mkdir("testdir", 0755)
    if err != nil {
        fmt.Printf("Error creating directory: %v\n", err)
    } else {
        fmt.Println("Directory 'testdir' created")
    }
    
    // Create nested directories
    err = os.MkdirAll("path/to/nested/dir", 0755)
    if err != nil {
        fmt.Printf("Error creating nested directories: %v\n", err)
    } else {
        fmt.Println("Nested directories created")
    }
    
    // Remove a directory (must be empty)
    err = os.Remove("testdir")
    if err != nil {
        fmt.Printf("Error removing directory: %v\n", err)
    } else {
        fmt.Println("Directory 'testdir' removed")
    }
    
    // Remove directory and all contents
    err = os.RemoveAll("path")
    if err != nil {
        fmt.Printf("Error removing directory tree: %v\n", err)
    } else {
        fmt.Println("Directory tree removed")
    }
}

// Check if directory exists
func directoryExists(path string) bool {
    info, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return info.IsDir()
}

// Create directory if it doesn't exist
func ensureDirectory(path string) error {
    if !directoryExists(path) {
        return os.MkdirAll(path, 0755)
    }
    return nil
}

func main() {
    directoryOperations()
    
    // Ensure directories exist
    dirs := []string{"logs", "data", "temp"}
    
    for _, dir := range dirs {
        if err := ensureDirectory(dir); err != nil {
            fmt.Printf("Error ensuring directory %s: %v\n", dir, err)
        } else {
            fmt.Printf("Directory %s is ready\n", dir)
        }
    }
}
```

### Listing Directory Contents

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

// Method 1: Using os.ReadDir (Go 1.16+)
func listDirectoryMethod1(dirPath string) {
    entries, err := os.ReadDir(dirPath)
    if err != nil {
        fmt.Printf("Error reading directory: %v\n", err)
        return
    }
    
    fmt.Printf("Contents of %s:\n", dirPath)
    for _, entry := range entries {
        fileType := "FILE"
        if entry.IsDir() {
            fileType = "DIR "
        }
        
        info, err := entry.Info()
        if err != nil {
            fmt.Printf("Error getting info for %s: %v\n", entry.Name(), err)
            continue
        }
        
        fmt.Printf("%s %10d %s %s\n", 
            fileType, 
            info.Size(), 
            info.ModTime().Format("2006-01-02 15:04:05"), 
            entry.Name())
    }
}

// Method 2: Using filepath.Walk for recursive listing
func listDirectoryRecursive(rootPath string) {
    fmt.Printf("Recursive listing of %s:\n", rootPath)
    
    err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            fmt.Printf("Error accessing %s: %v\n", path, err)
            return nil // Continue walking
        }
        
        // Calculate relative path
        relPath, _ := filepath.Rel(rootPath, path)
        if relPath == "." {
            relPath = "."
        }
        
        fileType := "FILE"
        if info.IsDir() {
            fileType = "DIR "
        }
        
        fmt.Printf("%s %10d %s %s\n", 
            fileType, 
            info.Size(), 
            info.ModTime().Format("2006-01-02 15:04:05"), 
            relPath)
        
        return nil
    })
    
    if err != nil {
        fmt.Printf("Error walking directory: %v\n", err)
    }
}

// Method 3: Using filepath.WalkDir (Go 1.16+) - more efficient
func listDirectoryWalkDir(rootPath string) {
    fmt.Printf("WalkDir listing of %s:\n", rootPath)
    
    err := filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            fmt.Printf("Error accessing %s: %v\n", path, err)
            return nil
        }
        
        relPath, _ := filepath.Rel(rootPath, path)
        if relPath == "." {
            relPath = "."
        }
        
        fileType := "FILE"
        if d.IsDir() {
            fileType = "DIR "
        }
        
        info, err := d.Info()
        if err != nil {
            fmt.Printf("Error getting info for %s: %v\n", path, err)
            return nil
        }
        
        fmt.Printf("%s %10d %s %s\n", 
            fileType, 
            info.Size(), 
            info.ModTime().Format("2006-01-02 15:04:05"), 
            relPath)
        
        return nil
    })
    
    if err != nil {
        fmt.Printf("Error walking directory: %v\n", err)
    }
}

// Filter files by extension
func findFilesByExtension(rootPath string, extension string) []string {
    var files []string
    
    filepath.WalkDir(rootPath, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        
        if !d.IsDir() && filepath.Ext(path) == extension {
            files = append(files, path)
        }
        
        return nil
    })
    
    return files
}

func main() {
    dirPath := "."
    
    fmt.Println("=== Method 1: os.ReadDir ===")
    listDirectoryMethod1(dirPath)
    
    fmt.Println("\n=== Method 2: filepath.Walk ===")
    listDirectoryRecursive(dirPath)
    
    fmt.Println("\n=== Method 3: filepath.WalkDir ===")
    listDirectoryWalkDir(dirPath)
    
    fmt.Println("\n=== Find Go files ===")
    goFiles := findFilesByExtension(".", ".go")
    for _, file := range goFiles {
        fmt.Println(file)
    }
}
```

## File Path Operations

### Working with Paths

```go
package main

import (
    "fmt"
    "path/filepath"
    "strings"
)

func pathOperations() {
    path := "/home/user/documents/file.txt"
    
    fmt.Printf("Original path: %s\n", path)
    fmt.Printf("Directory: %s\n", filepath.Dir(path))
    fmt.Printf("Filename: %s\n", filepath.Base(path))
    fmt.Printf("Extension: %s\n", filepath.Ext(path))
    
    // Remove extension
    nameWithoutExt := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
    fmt.Printf("Name without extension: %s\n", nameWithoutExt)
    
    // Join paths
    newPath := filepath.Join("/home", "user", "documents", "newfile.txt")
    fmt.Printf("Joined path: %s\n", newPath)
    
    // Clean path
    messyPath := "/home/user/../user/./documents//file.txt"
    cleanPath := filepath.Clean(messyPath)
    fmt.Printf("Messy path: %s\n", messyPath)
    fmt.Printf("Clean path: %s\n", cleanPath)
    
    // Absolute path
    relPath := "../documents/file.txt"
    absPath, err := filepath.Abs(relPath)
    if err != nil {
        fmt.Printf("Error getting absolute path: %v\n", err)
    } else {
        fmt.Printf("Relative path: %s\n", relPath)
        fmt.Printf("Absolute path: %s\n", absPath)
    }
    
    // Relative path
    basePath := "/home/user"
    targetPath := "/home/user/documents/file.txt"
    relativeToBase, err := filepath.Rel(basePath, targetPath)
    if err != nil {
        fmt.Printf("Error getting relative path: %v\n", err)
    } else {
        fmt.Printf("Path relative to %s: %s\n", basePath, relativeToBase)
    }
}

// Path matching with patterns
func pathMatching() {
    patterns := []string{
        "*.txt",
        "test_*.go",
        "[abc]*.log",
        "data/*/output.json",
    }
    
    testPaths := []string{
        "file.txt",
        "test_main.go",
        "afile.log",
        "data/2023/output.json",
        "readme.md",
    }
    
    fmt.Println("\n=== Path Matching ===")
    for _, pattern := range patterns {
        fmt.Printf("\nPattern: %s\n", pattern)
        for _, path := range testPaths {
            matched, err := filepath.Match(pattern, path)
            if err != nil {
                fmt.Printf("  Error matching %s: %v\n", path, err)
                continue
            }
            if matched {
                fmt.Printf("  ✓ %s\n", path)
            }
        }
    }
}

// Find files matching pattern
func findFilesWithPattern(rootDir, pattern string) ([]string, error) {
    var matches []string
    
    err := filepath.WalkDir(rootDir, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return nil
        }
        
        if d.IsDir() {
            return nil
        }
        
        matched, err := filepath.Match(pattern, filepath.Base(path))
        if err != nil {
            return err
        }
        
        if matched {
            matches = append(matches, path)
        }
        
        return nil
    })
    
    return matches, err
}

func main() {
    pathOperations()
    pathMatching()
    
    // Find all .go files
    fmt.Println("\n=== Finding .go files ===")
    goFiles, err := findFilesWithPattern(".", "*.go")
    if err != nil {
        fmt.Printf("Error finding files: %v\n", err)
    } else {
        for _, file := range goFiles {
            fmt.Println(file)
        }
    }
}
```

## File Information and Permissions

### Getting File Information

```go
package main

import (
    "fmt"
    "os"
    "time"
)

func getFileInfo(filename string) {
    info, err := os.Stat(filename)
    if err != nil {
        if os.IsNotExist(err) {
            fmt.Printf("File %s does not exist\n", filename)
        } else {
            fmt.Printf("Error getting file info: %v\n", err)
        }
        return
    }
    
    fmt.Printf("=== File Information for %s ===\n", filename)
    fmt.Printf("Name: %s\n", info.Name())
    fmt.Printf("Size: %d bytes\n", info.Size())
    fmt.Printf("Mode: %v\n", info.Mode())
    fmt.Printf("ModTime: %v\n", info.ModTime())
    fmt.Printf("IsDir: %v\n", info.IsDir())
    
    // Check file permissions
    mode := info.Mode()
    fmt.Printf("\n=== Permissions ===\n")
    fmt.Printf("Is regular file: %v\n", mode.IsRegular())
    fmt.Printf("Is directory: %v\n", mode.IsDir())
    fmt.Printf("Is symlink: %v\n", mode&os.ModeSymlink != 0)
    fmt.Printf("Is device: %v\n", mode&os.ModeDevice != 0)
    fmt.Printf("Is named pipe: %v\n", mode&os.ModeNamedPipe != 0)
    fmt.Printf("Is socket: %v\n", mode&os.ModeSocket != 0)
    fmt.Printf("Is setuid: %v\n", mode&os.ModeSetuid != 0)
    fmt.Printf("Is setgid: %v\n", mode&os.ModeSetgid != 0)
    fmt.Printf("Is sticky: %v\n", mode&os.ModeSticky != 0)
    
    // Permission bits
    perm := mode.Perm()
    fmt.Printf("\n=== Permission Bits ===\n")
    fmt.Printf("Octal: %o\n", perm)
    fmt.Printf("String: %s\n", perm.String())
    
    // Check specific permissions
    fmt.Printf("\n=== Permission Checks ===\n")
    fmt.Printf("Owner read: %v\n", perm&0400 != 0)
    fmt.Printf("Owner write: %v\n", perm&0200 != 0)
    fmt.Printf("Owner execute: %v\n", perm&0100 != 0)
    fmt.Printf("Group read: %v\n", perm&0040 != 0)
    fmt.Printf("Group write: %v\n", perm&0020 != 0)
    fmt.Printf("Group execute: %v\n", perm&0010 != 0)
    fmt.Printf("Other read: %v\n", perm&0004 != 0)
    fmt.Printf("Other write: %v\n", perm&0002 != 0)
    fmt.Printf("Other execute: %v\n", perm&0001 != 0)
}

// Check if file is readable/writable
func checkFileAccess(filename string) {
    fmt.Printf("\n=== Access Checks for %s ===\n", filename)
    
    // Check if file exists
    if _, err := os.Stat(filename); os.IsNotExist(err) {
        fmt.Println("File does not exist")
        return
    }
    
    // Try to open for reading
    if file, err := os.Open(filename); err == nil {
        file.Close()
        fmt.Println("File is readable")
    } else {
        fmt.Printf("File is not readable: %v\n", err)
    }
    
    // Try to open for writing
    if file, err := os.OpenFile(filename, os.O_WRONLY, 0); err == nil {
        file.Close()
        fmt.Println("File is writable")
    } else {
        fmt.Printf("File is not writable: %v\n", err)
    }
}

// Change file permissions
func changeFilePermissions(filename string, mode os.FileMode) {
    err := os.Chmod(filename, mode)
    if err != nil {
        fmt.Printf("Error changing permissions: %v\n", err)
        return
    }
    
    fmt.Printf("Changed permissions of %s to %o\n", filename, mode)
}

// Change file timestamps
func changeFileTimestamps(filename string) {
    now := time.Now()
    yesterday := now.Add(-24 * time.Hour)
    
    err := os.Chtimes(filename, yesterday, yesterday)
    if err != nil {
        fmt.Printf("Error changing timestamps: %v\n", err)
        return
    }
    
    fmt.Printf("Changed timestamps of %s\n", filename)
}

func main() {
    filename := "example.txt"
    
    // Create a test file
    file, err := os.Create(filename)
    if err != nil {
        fmt.Printf("Error creating test file: %v\n", err)
        return
    }
    file.WriteString("This is a test file.")
    file.Close()
    
    // Get file information
    getFileInfo(filename)
    
    // Check file access
    checkFileAccess(filename)
    
    // Change permissions
    changeFilePermissions(filename, 0644)
    
    // Change timestamps
    changeFileTimestamps(filename)
    
    // Get updated information
    fmt.Println("\n=== After Changes ===")
    getFileInfo(filename)
    
    // Clean up
    os.Remove(filename)
}
```

## Advanced File Operations

### File Copying and Moving

```go
package main

import (
    "fmt"
    "io"
    "os"
    "path/filepath"
)

// Copy file using io.Copy
func copyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("error opening source file: %w", err)
    }
    defer sourceFile.Close()
    
    destFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("error creating destination file: %w", err)
    }
    defer destFile.Close()
    
    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return fmt.Errorf("error copying file: %w", err)
    }
    
    // Copy file permissions
    sourceInfo, err := sourceFile.Stat()
    if err != nil {
        return fmt.Errorf("error getting source file info: %w", err)
    }
    
    err = os.Chmod(dst, sourceInfo.Mode())
    if err != nil {
        return fmt.Errorf("error setting file permissions: %w", err)
    }
    
    return nil
}

// Copy file with progress reporting
func copyFileWithProgress(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer sourceFile.Close()
    
    sourceInfo, err := sourceFile.Stat()
    if err != nil {
        return err
    }
    
    destFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer destFile.Close()
    
    // Create a progress reader
    var copied int64
    buffer := make([]byte, 32*1024) // 32KB buffer
    
    for {
        n, err := sourceFile.Read(buffer)
        if n > 0 {
            _, writeErr := destFile.Write(buffer[:n])
            if writeErr != nil {
                return writeErr
            }
            copied += int64(n)
            
            // Report progress
            progress := float64(copied) / float64(sourceInfo.Size()) * 100
            fmt.Printf("\rCopying: %.1f%% (%d/%d bytes)", progress, copied, sourceInfo.Size())
        }
        
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
    }
    
    fmt.Println() // New line after progress
    return nil
}

// Move file (rename)
func moveFile(src, dst string) error {
    err := os.Rename(src, dst)
    if err != nil {
        // If rename fails (e.g., across filesystems), copy and delete
        if copyErr := copyFile(src, dst); copyErr != nil {
            return fmt.Errorf("error copying file: %w", copyErr)
        }
        
        if removeErr := os.Remove(src); removeErr != nil {
            return fmt.Errorf("error removing source file: %w", removeErr)
        }
    }
    
    return nil
}

// Copy directory recursively
func copyDirectory(src, dst string) error {
    // Get source directory info
    srcInfo, err := os.Stat(src)
    if err != nil {
        return err
    }
    
    // Create destination directory
    err = os.MkdirAll(dst, srcInfo.Mode())
    if err != nil {
        return err
    }
    
    // Walk through source directory
    return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        // Calculate destination path
        relPath, err := filepath.Rel(src, path)
        if err != nil {
            return err
        }
        dstPath := filepath.Join(dst, relPath)
        
        if info.IsDir() {
            // Create directory
            return os.MkdirAll(dstPath, info.Mode())
        } else {
            // Copy file
            return copyFile(path, dstPath)
        }
    })
}

func main() {
    // Create test files
    testFile := "test_source.txt"
    file, err := os.Create(testFile)
    if err != nil {
        fmt.Printf("Error creating test file: %v\n", err)
        return
    }
    
    // Write some data
    for i := 0; i < 1000; i++ {
        file.WriteString(fmt.Sprintf("Line %d: This is some test data\n", i))
    }
    file.Close()
    
    // Copy file
    fmt.Println("=== Copying file ===")
    err = copyFile(testFile, "test_copy.txt")
    if err != nil {
        fmt.Printf("Error copying file: %v\n", err)
    } else {
        fmt.Println("File copied successfully")
    }
    
    // Copy with progress
    fmt.Println("\n=== Copying with progress ===")
    err = copyFileWithProgress(testFile, "test_copy_progress.txt")
    if err != nil {
        fmt.Printf("Error copying file with progress: %v\n", err)
    } else {
        fmt.Println("File copied with progress successfully")
    }
    
    // Move file
    fmt.Println("\n=== Moving file ===")
    err = moveFile("test_copy.txt", "test_moved.txt")
    if err != nil {
        fmt.Printf("Error moving file: %v\n", err)
    } else {
        fmt.Println("File moved successfully")
    }
    
    // Create test directory structure
    os.MkdirAll("test_src/subdir", 0755)
    os.WriteFile("test_src/file1.txt", []byte("File 1 content"), 0644)
    os.WriteFile("test_src/subdir/file2.txt", []byte("File 2 content"), 0644)
    
    // Copy directory
    fmt.Println("\n=== Copying directory ===")
    err = copyDirectory("test_src", "test_dst")
    if err != nil {
        fmt.Printf("Error copying directory: %v\n", err)
    } else {
        fmt.Println("Directory copied successfully")
    }
    
    // Clean up
    os.Remove(testFile)
    os.Remove("test_copy_progress.txt")
    os.Remove("test_moved.txt")
    os.RemoveAll("test_src")
    os.RemoveAll("test_dst")
}
```

### File Watching

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "time"
)

// Simple file watcher using polling
type FileWatcher struct {
    files    map[string]os.FileInfo
    interval time.Duration
    stop     chan bool
}

func NewFileWatcher(interval time.Duration) *FileWatcher {
    return &FileWatcher{
        files:    make(map[string]os.FileInfo),
        interval: interval,
        stop:     make(chan bool),
    }
}

func (fw *FileWatcher) AddFile(filename string) error {
    info, err := os.Stat(filename)
    if err != nil {
        return err
    }
    
    fw.files[filename] = info
    return nil
}

func (fw *FileWatcher) AddDirectory(dirPath string) error {
    return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        
        if !info.IsDir() {
            fw.files[path] = info
        }
        
        return nil
    })
}

func (fw *FileWatcher) Start() {
    ticker := time.NewTicker(fw.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            fw.checkFiles()
        case <-fw.stop:
            return
        }
    }
}

func (fw *FileWatcher) Stop() {
    fw.stop <- true
}

func (fw *FileWatcher) checkFiles() {
    for filename, oldInfo := range fw.files {
        newInfo, err := os.Stat(filename)
        if err != nil {
            if os.IsNotExist(err) {
                fmt.Printf("File deleted: %s\n", filename)
                delete(fw.files, filename)
            } else {
                fmt.Printf("Error checking file %s: %v\n", filename, err)
            }
            continue
        }
        
        // Check if file was modified
        if newInfo.ModTime().After(oldInfo.ModTime()) {
            fmt.Printf("File modified: %s\n", filename)
            fw.files[filename] = newInfo
        }
        
        // Check if file size changed
        if newInfo.Size() != oldInfo.Size() {
            fmt.Printf("File size changed: %s (was %d, now %d)\n", 
                filename, oldInfo.Size(), newInfo.Size())
            fw.files[filename] = newInfo
        }
    }
}

// Directory watcher that detects new files
type DirectoryWatcher struct {
    directory string
    files     map[string]bool
    interval  time.Duration
    stop      chan bool
}

func NewDirectoryWatcher(directory string, interval time.Duration) *DirectoryWatcher {
    return &DirectoryWatcher{
        directory: directory,
        files:     make(map[string]bool),
        interval:  interval,
        stop:      make(chan bool),
    }
}

func (dw *DirectoryWatcher) Start() {
    // Initial scan
    dw.scanDirectory()
    
    ticker := time.NewTicker(dw.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            dw.scanDirectory()
        case <-dw.stop:
            return
        }
    }
}

func (dw *DirectoryWatcher) Stop() {
    dw.stop <- true
}

func (dw *DirectoryWatcher) scanDirectory() {
    entries, err := os.ReadDir(dw.directory)
    if err != nil {
        fmt.Printf("Error reading directory: %v\n", err)
        return
    }
    
    currentFiles := make(map[string]bool)
    
    for _, entry := range entries {
        if !entry.IsDir() {
            filename := entry.Name()
            fullPath := filepath.Join(dw.directory, filename)
            currentFiles[fullPath] = true
            
            // Check if this is a new file
            if !dw.files[fullPath] {
                fmt.Printf("New file detected: %s\n", fullPath)
            }
        }
    }
    
    // Check for deleted files
    for filename := range dw.files {
        if !currentFiles[filename] {
            fmt.Printf("File deleted: %s\n", filename)
        }
    }
    
    dw.files = currentFiles
}

func main() {
    // Create test directory
    testDir := "watch_test"
    os.MkdirAll(testDir, 0755)
    
    // Create initial test file
    testFile := filepath.Join(testDir, "test.txt")
    os.WriteFile(testFile, []byte("Initial content"), 0644)
    
    // Start file watcher
    fmt.Println("Starting file watcher...")
    fileWatcher := NewFileWatcher(1 * time.Second)
    fileWatcher.AddFile(testFile)
    
    go fileWatcher.Start()
    
    // Start directory watcher
    fmt.Println("Starting directory watcher...")
    dirWatcher := NewDirectoryWatcher(testDir, 1*time.Second)
    go dirWatcher.Start()
    
    // Simulate file changes
    time.Sleep(2 * time.Second)
    
    fmt.Println("Modifying file...")
    os.WriteFile(testFile, []byte("Modified content"), 0644)
    
    time.Sleep(2 * time.Second)
    
    fmt.Println("Creating new file...")
    newFile := filepath.Join(testDir, "new.txt")
    os.WriteFile(newFile, []byte("New file content"), 0644)
    
    time.Sleep(2 * time.Second)
    
    fmt.Println("Deleting file...")
    os.Remove(newFile)
    
    time.Sleep(2 * time.Second)
    
    // Stop watchers
    fileWatcher.Stop()
    dirWatcher.Stop()
    
    // Clean up
    os.RemoveAll(testDir)
    
    fmt.Println("File watching demo completed")
}
```

## Working with CSV Files

### Reading and Writing CSV

```go
package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
)

type Person struct {
    ID    int
    Name  string
    Email string
    Age   int
}

// Write CSV file
func writeCSV(filename string, people []Person) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    // Write header
    header := []string{"ID", "Name", "Email", "Age"}
    if err := writer.Write(header); err != nil {
        return err
    }
    
    // Write data
    for _, person := range people {
        record := []string{
            strconv.Itoa(person.ID),
            person.Name,
            person.Email,
            strconv.Itoa(person.Age),
        }
        
        if err := writer.Write(record); err != nil {
            return err
        }
    }
    
    return nil
}

// Read CSV file
func readCSV(filename string) ([]Person, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }
    
    var people []Person
    
    // Skip header (first record)
    for i, record := range records {
        if i == 0 {
            continue
        }
        
        if len(record) != 4 {
            fmt.Printf("Skipping invalid record: %v\n", record)
            continue
        }
        
        id, err := strconv.Atoi(record[0])
        if err != nil {
            fmt.Printf("Invalid ID in record: %v\n", record)
            continue
        }
        
        age, err := strconv.Atoi(record[3])
        if err != nil {
            fmt.Printf("Invalid age in record: %v\n", record)
            continue
        }
        
        person := Person{
            ID:    id,
            Name:  record[1],
            Email: record[2],
            Age:   age,
        }
        
        people = append(people, person)
    }
    
    return people, nil
}

// Read CSV with custom delimiter
func readCSVWithCustomDelimiter(filename string, delimiter rune) ([][]string, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    
    reader := csv.NewReader(file)
    reader.Comma = delimiter
    
    return reader.ReadAll()
}

func main() {
    // Sample data
    people := []Person{
        {1, "Alice Johnson", "alice@example.com", 30},
        {2, "Bob Smith", "bob@example.com", 25},
        {3, "Charlie Brown", "charlie@example.com", 35},
        {4, "Diana Prince", "diana@example.com", 28},
    }
    
    filename := "people.csv"
    
    // Write CSV
    fmt.Println("Writing CSV file...")
    err := writeCSV(filename, people)
    if err != nil {
        fmt.Printf("Error writing CSV: %v\n", err)
        return
    }
    
    // Read CSV
    fmt.Println("Reading CSV file...")
    readPeople, err := readCSV(filename)
    if err != nil {
        fmt.Printf("Error reading CSV: %v\n", err)
        return
    }
    
    // Display results
    fmt.Println("\nPeople from CSV:")
    for _, person := range readPeople {
        fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d\n", 
            person.ID, person.Name, person.Email, person.Age)
    }
    
    // Clean up
    os.Remove(filename)
}
```

## Practical Examples

### Log File Analyzer

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
    IP        string
}

type LogAnalyzer struct {
    entries []LogEntry
}

func NewLogAnalyzer() *LogAnalyzer {
    return &LogAnalyzer{
        entries: make([]LogEntry, 0),
    }
}

func (la *LogAnalyzer) ParseLogFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    scanner := bufio.NewScanner(file)
    
    // Regular expression for common log format
    // Example: 2023-12-01 10:30:45 [INFO] 192.168.1.1 User login successful
    logRegex := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) \[(\w+)\] (\d+\.\d+\.\d+\.\d+) (.+)`)
    
    for scanner.Scan() {
        line := scanner.Text()
        matches := logRegex.FindStringSubmatch(line)
        
        if len(matches) == 5 {
            timestamp, err := time.Parse("2006-01-02 15:04:05", matches[1])
            if err != nil {
                continue
            }
            
            entry := LogEntry{
                Timestamp: timestamp,
                Level:     matches[2],
                IP:        matches[3],
                Message:   matches[4],
            }
            
            la.entries = append(la.entries, entry)
        }
    }
    
    return scanner.Err()
}

func (la *LogAnalyzer) GetEntriesByLevel(level string) []LogEntry {
    var filtered []LogEntry
    
    for _, entry := range la.entries {
        if strings.EqualFold(entry.Level, level) {
            filtered = append(filtered, entry)
        }
    }
    
    return filtered
}

func (la *LogAnalyzer) GetEntriesByTimeRange(start, end time.Time) []LogEntry {
    var filtered []LogEntry
    
    for _, entry := range la.entries {
        if entry.Timestamp.After(start) && entry.Timestamp.Before(end) {
            filtered = append(filtered, entry)
        }
    }
    
    return filtered
}

func (la *LogAnalyzer) GetTopIPs(limit int) map[string]int {
    ipCounts := make(map[string]int)
    
    for _, entry := range la.entries {
        ipCounts[entry.IP]++
    }
    
    // Sort by count
    type ipCount struct {
        IP    string
        Count int
    }
    
    var sorted []ipCount
    for ip, count := range ipCounts {
        sorted = append(sorted, ipCount{IP: ip, Count: count})
    }
    
    sort.Slice(sorted, func(i, j int) bool {
        return sorted[i].Count > sorted[j].Count
    })
    
    result := make(map[string]int)
    for i, item := range sorted {
        if i >= limit {
            break
        }
        result[item.IP] = item.Count
    }
    
    return result
}

func (la *LogAnalyzer) GenerateReport(outputFile string) error {
    file, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer file.Close()
    
    fmt.Fprintf(file, "Log Analysis Report\n")
    fmt.Fprintf(file, "==================\n\n")
    
    fmt.Fprintf(file, "Total entries: %d\n\n", len(la.entries))
    
    // Count by level
    levelCounts := make(map[string]int)
    for _, entry := range la.entries {
        levelCounts[entry.Level]++
    }
    
    fmt.Fprintf(file, "Entries by level:\n")
    for level, count := range levelCounts {
        fmt.Fprintf(file, "  %s: %d\n", level, count)
    }
    
    fmt.Fprintf(file, "\nTop 5 IP addresses:\n")
    topIPs := la.GetTopIPs(5)
    for ip, count := range topIPs {
        fmt.Fprintf(file, "  %s: %d requests\n", ip, count)
    }
    
    // Time range
    if len(la.entries) > 0 {
        sort.Slice(la.entries, func(i, j int) bool {
            return la.entries[i].Timestamp.Before(la.entries[j].Timestamp)
        })
        
        fmt.Fprintf(file, "\nTime range:\n")
        fmt.Fprintf(file, "  From: %s\n", la.entries[0].Timestamp.Format("2006-01-02 15:04:05"))
        fmt.Fprintf(file, "  To: %s\n", la.entries[len(la.entries)-1].Timestamp.Format("2006-01-02 15:04:05"))
    }
    
    return nil
}

func createSampleLogFile(filename string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()
    
    // Generate sample log entries
    entries := []string{
        "2023-12-01 10:30:45 [INFO] 192.168.1.100 User login successful",
        "2023-12-01 10:31:12 [ERROR] 192.168.1.101 Failed login attempt",
        "2023-12-01 10:32:05 [INFO] 192.168.1.100 Page accessed: /dashboard",
        "2023-12-01 10:33:22 [WARN] 192.168.1.102 Slow query detected",
        "2023-12-01 10:34:15 [INFO] 192.168.1.103 User logout",
        "2023-12-01 10:35:08 [ERROR] 192.168.1.101 Database connection failed",
        "2023-12-01 10:36:33 [INFO] 192.168.1.100 File uploaded",
        "2023-12-01 10:37:45 [DEBUG] 192.168.1.104 Cache miss for key: user_123",
        "2023-12-01 10:38:12 [INFO] 192.168.1.105 API request processed",
        "2023-12-01 10:39:28 [ERROR] 192.168.1.106 Payment processing failed",
    }
    
    for _, entry := range entries {
        _, err := file.WriteString(entry + "\n")
        if err != nil {
            return err
        }
    }
    
    return nil
}

func main() {
    logFile := "sample.log"
    reportFile := "analysis_report.txt"
    
    // Create sample log file
    fmt.Println("Creating sample log file...")
    err := createSampleLogFile(logFile)
    if err != nil {
        fmt.Printf("Error creating sample log: %v\n", err)
        return
    }
    
    // Analyze log file
    fmt.Println("Analyzing log file...")
    analyzer := NewLogAnalyzer()
    err = analyzer.ParseLogFile(logFile)
    if err != nil {
        fmt.Printf("Error parsing log file: %v\n", err)
        return
    }
    
    // Get error entries
    errorEntries := analyzer.GetEntriesByLevel("ERROR")
    fmt.Printf("Found %d error entries\n", len(errorEntries))
    
    // Get top IPs
    topIPs := analyzer.GetTopIPs(3)
    fmt.Println("Top IP addresses:")
    for ip, count := range topIPs {
        fmt.Printf("  %s: %d requests\n", ip, count)
    }
    
    // Generate report
    fmt.Println("Generating analysis report...")
    err = analyzer.GenerateReport(reportFile)
    if err != nil {
        fmt.Printf("Error generating report: %v\n", err)
        return
    }
    
    fmt.Printf("Analysis complete. Report saved to %s\n", reportFile)
    
    // Clean up
    os.Remove(logFile)
    os.Remove(reportFile)
}
```

### Configuration File Manager

```go
package main

import (
    "encoding/json"
    "fmt"
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
    Host string `json:"host"`
    Port int    `json:"port"`
    SSL  bool   `json:"ssl"`
}

type AppConfig struct {
    AppName  string         `json:"app_name"`
    Version  string         `json:"version"`
    Debug    bool           `json:"debug"`
    Database DatabaseConfig `json:"database"`
    Server   ServerConfig   `json:"server"`
    LogLevel string         `json:"log_level"`
}

type ConfigManager struct {
    configPath string
    config     *AppConfig
}

func NewConfigManager(configPath string) *ConfigManager {
    return &ConfigManager{
        configPath: configPath,
    }
}

func (cm *ConfigManager) LoadConfig() error {
    // Check if config file exists
    if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
        // Create default config
        return cm.CreateDefaultConfig()
    }
    
    // Read config file
    data, err := os.ReadFile(cm.configPath)
    if err != nil {
        return fmt.Errorf("error reading config file: %w", err)
    }
    
    // Parse JSON
    cm.config = &AppConfig{}
    err = json.Unmarshal(data, cm.config)
    if err != nil {
        return fmt.Errorf("error parsing config file: %w", err)
    }
    
    return nil
}

func (cm *ConfigManager) CreateDefaultConfig() error {
    defaultConfig := &AppConfig{
        AppName: "MyApp",
        Version: "1.0.0",
        Debug:   false,
        Database: DatabaseConfig{
            Host:     "localhost",
            Port:     5432,
            Username: "user",
            Password: "password",
            Database: "myapp",
        },
        Server: ServerConfig{
            Host: "localhost",
            Port: 8080,
            SSL:  false,
        },
        LogLevel: "info",
    }
    
    cm.config = defaultConfig
    return cm.SaveConfig()
}

func (cm *ConfigManager) SaveConfig() error {
    // Ensure directory exists
    dir := filepath.Dir(cm.configPath)
    err := os.MkdirAll(dir, 0755)
    if err != nil {
        return fmt.Errorf("error creating config directory: %w", err)
    }
    
    // Marshal to JSON with indentation
    data, err := json.MarshalIndent(cm.config, "", "  ")
    if err != nil {
        return fmt.Errorf("error marshaling config: %w", err)
    }
    
    // Write to file
    err = os.WriteFile(cm.configPath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing config file: %w", err)
    }
    
    return nil
}

func (cm *ConfigManager) GetConfig() *AppConfig {
    return cm.config
}

func (cm *ConfigManager) UpdateDatabaseConfig(dbConfig DatabaseConfig) error {
    cm.config.Database = dbConfig
    return cm.SaveConfig()
}

func (cm *ConfigManager) UpdateServerConfig(serverConfig ServerConfig) error {
    cm.config.Server = serverConfig
    return cm.SaveConfig()
}

func (cm *ConfigManager) SetDebugMode(debug bool) error {
    cm.config.Debug = debug
    return cm.SaveConfig()
}

func (cm *ConfigManager) BackupConfig(backupPath string) error {
    data, err := os.ReadFile(cm.configPath)
    if err != nil {
        return fmt.Errorf("error reading config file: %w", err)
    }
    
    err = os.WriteFile(backupPath, data, 0644)
    if err != nil {
        return fmt.Errorf("error writing backup file: %w", err)
    }
    
    return nil
}

func (cm *ConfigManager) RestoreConfig(backupPath string) error {
    data, err := os.ReadFile(backupPath)
    if err != nil {
        return fmt.Errorf("error reading backup file: %w", err)
    }
    
    err = os.WriteFile(cm.configPath, data, 0644)
    if err != nil {
        return fmt.Errorf("error restoring config file: %w", err)
    }
    
    return cm.LoadConfig()
}

func main() {
    configPath := "config/app.json"
    backupPath := "config/app.backup.json"
    
    // Create config manager
    cm := NewConfigManager(configPath)
    
    // Load or create config
    fmt.Println("Loading configuration...")
    err := cm.LoadConfig()
    if err != nil {
        fmt.Printf("Error loading config: %v\n", err)
        return
    }
    
    // Display current config
    config := cm.GetConfig()
    fmt.Printf("App: %s v%s\n", config.AppName, config.Version)
    fmt.Printf("Debug mode: %v\n", config.Debug)
    fmt.Printf("Database: %s@%s:%d/%s\n", 
        config.Database.Username, config.Database.Host, 
        config.Database.Port, config.Database.Database)
    fmt.Printf("Server: %s:%d (SSL: %v)\n", 
        config.Server.Host, config.Server.Port, config.Server.SSL)
    
    // Create backup
    fmt.Println("\nCreating backup...")
    err = cm.BackupConfig(backupPath)
    if err != nil {
        fmt.Printf("Error creating backup: %v\n", err)
    } else {
        fmt.Println("Backup created successfully")
    }
    
    // Update configuration
    fmt.Println("\nUpdating configuration...")
    err = cm.SetDebugMode(true)
    if err != nil {
        fmt.Printf("Error updating debug mode: %v\n", err)
    }
    
    newServerConfig := ServerConfig{
        Host: "0.0.0.0",
        Port: 9090,
        SSL:  true,
    }
    
    err = cm.UpdateServerConfig(newServerConfig)
    if err != nil {
        fmt.Printf("Error updating server config: %v\n", err)
    } else {
        fmt.Println("Configuration updated successfully")
    }
    
    // Display updated config
    config = cm.GetConfig()
    fmt.Printf("\nUpdated configuration:\n")
    fmt.Printf("Debug mode: %v\n", config.Debug)
    fmt.Printf("Server: %s:%d (SSL: %v)\n", 
        config.Server.Host, config.Server.Port, config.Server.SSL)
    
    // Clean up
    os.RemoveAll("config")
}
```

## Best Practices

### 1. Always Close Files

```go
file, err := os.Open("example.txt")
if err != nil {
    return err
}
defer file.Close() // Always use defer
```

### 2. Handle Errors Properly

```go
if _, err := os.Stat(filename); os.IsNotExist(err) {
    // File doesn't exist
} else if err != nil {
    // Other error
} else {
    // File exists
}
```

### 3. Use Buffered I/O for Large Files

```go
// For reading
reader := bufio.NewReader(file)

// For writing
writer := bufio.NewWriter(file)
defer writer.Flush()
```

### 4. Check File Permissions

```go
info, err := os.Stat(filename)
if err != nil {
    return err
}

if info.Mode().Perm()&0200 == 0 {
    return errors.New("file is not writable")
}
```

### 5. Use Appropriate File Modes

```go
// Common file modes
os.O_RDONLY          // Read only
os.O_WRONLY          // Write only
os.O_RDWR            // Read and write
os.O_APPEND          // Append to file
os.O_CREATE          // Create if doesn't exist
os.O_TRUNC           // Truncate file
os.O_EXCL            // Fail if file exists
```

## Exercises

### Exercise 1: File Backup Utility

Create a file backup utility that:
- Backs up files to a specified directory
- Maintains multiple backup versions
- Compresses backups
- Provides restore functionality

### Exercise 2: Directory Synchronizer

Implement a directory synchronization tool that:
- Compares two directories
- Identifies differences (new, modified, deleted files)
- Synchronizes directories
- Provides dry-run mode

### Exercise 3: Log Rotation System

Build a log rotation system that:
- Monitors log file size
- Rotates logs when they exceed size limit
- Maintains specified number of old logs
- Compresses old log files

## Key Takeaways

1. **Always close files**: Use `defer file.Close()` to ensure files are closed.

2. **Handle errors properly**: Check for specific error types like `os.IsNotExist()`.

3. **Use appropriate file modes**: Choose the right combination of flags for your use case.

4. **Buffer I/O operations**: Use `bufio` for better performance with large files.

5. **Check file permissions**: Verify file accessibility before operations.

6. **Use filepath package**: For cross-platform path operations.

7. **Be careful with file paths**: Use absolute paths when necessary.

8. **Consider file locking**: For concurrent access scenarios.

9. **Validate file operations**: Check return values and handle partial operations.

10. **Clean up resources**: Remove temporary files and directories.

## Next Steps

Now that you understand file I/O and system programming, let's explore [Concurrency Patterns](18-concurrency-patterns.md) to learn advanced concurrent programming techniques!

---

**Previous**: [← Packages and Modules](16-packages-modules.md) | **Next**: [Concurrency Patterns →](18-concurrency-patterns.md)