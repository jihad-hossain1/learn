# 2. Installation and Setup

## Installing Go

### Method 1: Official Installer (Recommended)

#### Windows
1. Visit [https://golang.org/dl/](https://golang.org/dl/)
2. Download the Windows installer (`.msi` file)
3. Run the installer and follow the setup wizard
4. Go will be installed to `C:\Program Files\Go` by default

#### macOS
1. Visit [https://golang.org/dl/](https://golang.org/dl/)
2. Download the macOS installer (`.pkg` file)
3. Run the installer and follow the instructions
4. Go will be installed to `/usr/local/go`

#### Linux
1. Download the Linux tarball from [https://golang.org/dl/](https://golang.org/dl/)
2. Extract to `/usr/local`:
   ```bash
   sudo tar -C /usr/local -xzf go1.21.x.linux-amd64.tar.gz
   ```
3. Add Go to your PATH in `~/.bashrc` or `~/.zshrc`:
   ```bash
   export PATH=$PATH:/usr/local/go/bin
   ```

### Method 2: Package Managers

#### Windows (Chocolatey)
```powershell
choco install golang
```

#### macOS (Homebrew)
```bash
brew install go
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang-go
```

## Verifying Installation

Open a terminal/command prompt and run:
```bash
go version
```

You should see output like:
```
go version go1.21.x operating_system/architecture
```

## Setting Up Go Workspace

### Understanding GOPATH vs Go Modules

- **GOPATH**: Legacy workspace model (Go < 1.11)
- **Go Modules**: Modern dependency management (Go 1.11+)

### Go Modules (Recommended)

Go modules allow you to work anywhere on your system without setting up GOPATH.

#### Creating Your First Module

1. Create a project directory:
   ```bash
   mkdir my-go-project
   cd my-go-project
   ```

2. Initialize a Go module:
   ```bash
   go mod init example.com/my-go-project
   ```

3. This creates a `go.mod` file:
   ```go
   module example.com/my-go-project
   
   go 1.21
   ```

## Environment Variables

### Important Go Environment Variables

- **GOROOT**: Go installation directory
- **GOPATH**: Workspace directory (legacy)
- **GOPROXY**: Module proxy settings
- **GOSUMDB**: Checksum database
- **GO111MODULE**: Module mode (auto/on/off)

### Checking Environment Variables
```bash
go env
```

### Setting Environment Variables

#### Windows
```cmd
set GOPROXY=https://proxy.golang.org,direct
```

#### macOS/Linux
```bash
export GOPROXY=https://proxy.golang.org,direct
```

## Development Environment Setup

### Code Editors and IDEs

#### 1. Visual Studio Code (Recommended)
- Install the Go extension by Google
- Features: IntelliSense, debugging, testing, formatting

**Setup Steps:**
1. Install VS Code
2. Install Go extension
3. Open Command Palette (`Ctrl+Shift+P`)
4. Run "Go: Install/Update Tools"

#### 2. GoLand (JetBrains)
- Professional IDE with advanced features
- Built-in debugger, profiler, and testing tools

#### 3. Vim/Neovim
- Install vim-go plugin
- Lightweight and fast

#### 4. Sublime Text
- Install GoSublime package
- Simple and efficient

### Essential Go Tools

Install these tools for better development experience:

```bash
# Code formatting
go install golang.org/x/tools/cmd/goimports@latest

# Linting
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Documentation
go install golang.org/x/tools/cmd/godoc@latest

# Debugging
go install github.com/go-delve/delve/cmd/dlv@latest
```

## Your First Go Program

Create a file named `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### Running the Program

```bash
# Method 1: Run directly
go run main.go

# Method 2: Build and run
go build main.go
./main        # Linux/macOS
main.exe      # Windows
```

## Go Commands Overview

| Command | Description |
|---------|-------------|
| `go run` | Compile and run Go program |
| `go build` | Compile packages and dependencies |
| `go install` | Compile and install packages |
| `go mod init` | Initialize new module |
| `go mod tidy` | Add missing and remove unused modules |
| `go get` | Download and install packages |
| `go test` | Run tests |
| `go fmt` | Format Go source code |
| `go vet` | Examine Go source code |
| `go doc` | Show documentation |

## Project Structure Best Practices

### Simple Project
```
my-project/
├── go.mod
├── go.sum
├── main.go
├── README.md
└── internal/
    └── helper.go
```

### Larger Project
```
my-app/
├── go.mod
├── go.sum
├── main.go
├── README.md
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/
│   └── utils/
├── api/
├── web/
├── configs/
├── scripts/
└── docs/
```

## Configuration Tips

### 1. Enable Go Modules
```bash
go env -w GO111MODULE=on
```

### 2. Set Module Proxy
```bash
go env -w GOPROXY=https://proxy.golang.org,direct
```

### 3. Private Repositories
```bash
go env -w GOPRIVATE=github.com/yourcompany/*
```

## Troubleshooting

### Common Issues

1. **"go: command not found"**
   - Check if Go is in your PATH
   - Restart terminal after installation

2. **Module issues**
   - Run `go mod tidy` to clean up dependencies
   - Check `go.mod` file syntax

3. **Permission errors**
   - Check file permissions
   - Run with appropriate privileges

### Useful Commands for Debugging

```bash
# Check Go installation
go version

# Check environment
go env

# Verify module
go mod verify

# Clean module cache
go clean -modcache
```

## Next Steps

Now that you have Go installed and configured, let's move on to [Basic Syntax and Hello World](03-basic-syntax.md) to start writing Go code!

---

## Quick Reference

### Installation Verification Checklist
- [ ] `go version` works
- [ ] Can create and run "Hello, World!" program
- [ ] Go modules initialized successfully
- [ ] Development environment configured
- [ ] Essential tools installed

---

**Previous**: [← Introduction](01-introduction.md) | **Next**: [Basic Syntax →](03-basic-syntax.md)