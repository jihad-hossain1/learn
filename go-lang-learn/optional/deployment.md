# 25. Deployment

This chapter covers building, packaging, and deploying Go applications across different environments and platforms. We'll explore various deployment strategies, containerization, and production best practices.

## Building Go Applications

### Basic Building

```bash
# Build for current platform
go build main.go

# Build with custom output name
go build -o myapp main.go

# Build all files in current directory
go build .

# Build and install to $GOPATH/bin
go install .

# Build with optimizations
go build -ldflags="-s -w" main.go
# -s: strip symbol table
# -w: strip debug info
```

### Cross-Platform Building

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o myapp-linux main.go

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o myapp.exe main.go

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o myapp-mac main.go

# Build for ARM (Raspberry Pi)
GOOS=linux GOARCH=arm go build -o myapp-arm main.go

# Build for ARM64
GOOS=linux GOARCH=arm64 go build -o myapp-arm64 main.go

# List all supported platforms
go tool dist list
```

### Build Scripts

```bash
#!/bin/bash
# build.sh - Cross-platform build script

APP_NAME="myapp"
VERSION="1.0.0"
BUILD_DIR="build"

# Create build directory
mkdir -p $BUILD_DIR

# Build flags
LDFLAGS="-s -w -X main.version=$VERSION -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)"

echo "Building $APP_NAME v$VERSION..."

# Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="$LDFLAGS" -o $BUILD_DIR/${APP_NAME}-linux-amd64 .

# Windows AMD64
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o $BUILD_DIR/${APP_NAME}-windows-amd64.exe .

# macOS AMD64
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="$LDFLAGS" -o $BUILD_DIR/${APP_NAME}-darwin-amd64 .

# Linux ARM64
echo "Building for Linux ARM64..."
GOOS=linux GOARCH=arm64 go build -ldflags="$LDFLAGS" -o $BUILD_DIR/${APP_NAME}-linux-arm64 .

echo "Build complete! Binaries are in $BUILD_DIR/"
ls -la $BUILD_DIR/
```

### Build Information Injection

```go
// main.go
package main

import (
    "fmt"
    "runtime"
)

// These variables will be set at build time
var (
    version   = "dev"
    buildTime = "unknown"
    gitCommit = "unknown"
)

func main() {
    fmt.Printf("Application Info:\n")
    fmt.Printf("  Version: %s\n", version)
    fmt.Printf("  Build Time: %s\n", buildTime)
    fmt.Printf("  Git Commit: %s\n", gitCommit)
    fmt.Printf("  Go Version: %s\n", runtime.Version())
    fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
    
    // Your application logic here
    fmt.Println("\nApplication is running...")
}
```

```bash
# Build with version information
VERSION="1.2.3"
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)
GIT_COMMIT=$(git rev-parse --short HEAD)

go build -ldflags="\
    -X main.version=$VERSION \
    -X main.buildTime=$BUILD_TIME \
    -X main.gitCommit=$GIT_COMMIT" \
    -o myapp main.go
```

## Configuration Management

### Environment-Based Configuration

```go
// config/config.go
package config

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Logging  LoggingConfig
}

type ServerConfig struct {
    Host         string
    Port         int
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    IdleTimeout  time.Duration
}

type DatabaseConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    Database string
    SSLMode  string
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    Database int
}

type LoggingConfig struct {
    Level  string
    Format string
    Output string
}

func Load() (*Config, error) {
    config := &Config{
        Server: ServerConfig{
            Host:         getEnv("SERVER_HOST", "localhost"),
            Port:         getEnvAsInt("SERVER_PORT", 8080),
            ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", 30*time.Second),
            WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", 30*time.Second),
            IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            Username: getEnv("DB_USERNAME", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            Database: getEnv("DB_DATABASE", "myapp"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            Database: getEnvAsInt("REDIS_DATABASE", 0),
        },
        Logging: LoggingConfig{
            Level:  getEnv("LOG_LEVEL", "info"),
            Format: getEnv("LOG_FORMAT", "json"),
            Output: getEnv("LOG_OUTPUT", "stdout"),
        },
    }
    
    return config, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}

func (c *Config) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid server port: %d", c.Server.Port)
    }
    
    if c.Database.Host == "" {
        return fmt.Errorf("database host is required")
    }
    
    if c.Database.Username == "" {
        return fmt.Errorf("database username is required")
    }
    
    return nil
}

func (c *Config) Print() {
    fmt.Printf("Configuration:\n")
    fmt.Printf("  Server: %s:%d\n", c.Server.Host, c.Server.Port)
    fmt.Printf("  Database: %s:%d/%s\n", c.Database.Host, c.Database.Port, c.Database.Database)
    fmt.Printf("  Redis: %s:%d/%d\n", c.Redis.Host, c.Redis.Port, c.Redis.Database)
    fmt.Printf("  Logging: %s (%s)\n", c.Logging.Level, c.Logging.Format)
}
```

### Configuration Files

```yaml
# config.yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: "30s"
  write_timeout: "30s"
  idle_timeout: "60s"

database:
  host: "localhost"
  port: 5432
  username: "myapp"
  password: "secret"
  database: "myapp_prod"
  ssl_mode: "require"

redis:
  host: "localhost"
  port: 6379
  password: ""
  database: 0

logging:
  level: "info"
  format: "json"
  output: "stdout"
```

```go
// config/yaml.go
package config

import (
    "fmt"
    "io/ioutil"
    "time"
    
    "gopkg.in/yaml.v2"
)

type YAMLConfig struct {
    Server struct {
        Host         string `yaml:"host"`
        Port         int    `yaml:"port"`
        ReadTimeout  string `yaml:"read_timeout"`
        WriteTimeout string `yaml:"write_timeout"`
        IdleTimeout  string `yaml:"idle_timeout"`
    } `yaml:"server"`
    
    Database struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Username string `yaml:"username"`
        Password string `yaml:"password"`
        Database string `yaml:"database"`
        SSLMode  string `yaml:"ssl_mode"`
    } `yaml:"database"`
    
    Redis struct {
        Host     string `yaml:"host"`
        Port     int    `yaml:"port"`
        Password string `yaml:"password"`
        Database int    `yaml:"database"`
    } `yaml:"redis"`
    
    Logging struct {
        Level  string `yaml:"level"`
        Format string `yaml:"format"`
        Output string `yaml:"output"`
    } `yaml:"logging"`
}

func LoadFromYAML(filename string) (*Config, error) {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read config file: %w", err)
    }
    
    var yamlConfig YAMLConfig
    if err := yaml.Unmarshal(data, &yamlConfig); err != nil {
        return nil, fmt.Errorf("failed to parse YAML: %w", err)
    }
    
    // Convert to our Config struct
    config := &Config{
        Server: ServerConfig{
            Host: yamlConfig.Server.Host,
            Port: yamlConfig.Server.Port,
        },
        Database: DatabaseConfig{
            Host:     yamlConfig.Database.Host,
            Port:     yamlConfig.Database.Port,
            Username: yamlConfig.Database.Username,
            Password: yamlConfig.Database.Password,
            Database: yamlConfig.Database.Database,
            SSLMode:  yamlConfig.Database.SSLMode,
        },
        Redis: RedisConfig{
            Host:     yamlConfig.Redis.Host,
            Port:     yamlConfig.Redis.Port,
            Password: yamlConfig.Redis.Password,
            Database: yamlConfig.Redis.Database,
        },
        Logging: LoggingConfig{
            Level:  yamlConfig.Logging.Level,
            Format: yamlConfig.Logging.Format,
            Output: yamlConfig.Logging.Output,
        },
    }
    
    // Parse durations
    if config.Server.ReadTimeout, err = time.ParseDuration(yamlConfig.Server.ReadTimeout); err != nil {
        return nil, fmt.Errorf("invalid read timeout: %w", err)
    }
    
    if config.Server.WriteTimeout, err = time.ParseDuration(yamlConfig.Server.WriteTimeout); err != nil {
        return nil, fmt.Errorf("invalid write timeout: %w", err)
    }
    
    if config.Server.IdleTimeout, err = time.ParseDuration(yamlConfig.Server.IdleTimeout); err != nil {
        return nil, fmt.Errorf("invalid idle timeout: %w", err)
    }
    
    return config, nil
}
```

## Docker Deployment

### Basic Dockerfile

```dockerfile
# Dockerfile
# Multi-stage build for smaller final image

# Build stage
FROM golang:1.21-alpine AS builder

# Install git (needed for go mod download)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/main .

# Copy config files if needed
COPY --from=builder /app/config.yaml .

# Change ownership to appuser
RUN chown -R appuser:appuser /root

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
```

### Optimized Dockerfile

```dockerfile
# Dockerfile.optimized
# Optimized multi-stage build

# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /build

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o app .

# Final stage - use scratch for minimal image
FROM scratch

# Copy timezone data
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Copy SSL certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy passwd file for user
COPY --from=builder /etc/passwd /etc/passwd

# Copy binary
COPY --from=builder /build/app /app

# Use non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/app"]
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USERNAME=myapp
      - DB_PASSWORD=secret
      - DB_DATABASE=myapp
      - REDIS_HOST=redis
      - LOG_LEVEL=info
    depends_on:
      - postgres
      - redis
    restart: unless-stopped
    networks:
      - app-network

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=myapp
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    ports:
      - "5432:5432"
    restart: unless-stopped
    networks:
      - app-network

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    ports:
      - "6379:6379"
    restart: unless-stopped
    networks:
      - app-network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
    networks:
      - app-network

volumes:
  postgres_data:
  redis_data:

networks:
  app-network:
    driver: bridge
```

### Production Docker Compose

```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  app:
    image: myapp:latest
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.25'
          memory: 256M
    environment:
      - DB_HOST=postgres
      - DB_USERNAME_FILE=/run/secrets/db_username
      - DB_PASSWORD_FILE=/run/secrets/db_password
      - REDIS_HOST=redis
      - LOG_LEVEL=warn
    secrets:
      - db_username
      - db_password
    networks:
      - app-network
      - db-network
    healthcheck:
      test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER_FILE=/run/secrets/db_username
      - POSTGRES_PASSWORD_FILE=/run/secrets/db_password
      - POSTGRES_DB=myapp
    volumes:
      - postgres_data:/var/lib/postgresql/data
    secrets:
      - db_username
      - db_password
    networks:
      - db-network
    deploy:
      placement:
        constraints:
          - node.role == manager

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes --requirepass-file /run/secrets/redis_password
    volumes:
      - redis_data:/data
    secrets:
      - redis_password
    networks:
      - app-network

secrets:
  db_username:
    external: true
  db_password:
    external: true
  redis_password:
    external: true

volumes:
  postgres_data:
    external: true
  redis_data:
    external: true

networks:
  app-network:
    external: true
  db-network:
    external: true
```

## Kubernetes Deployment

### Basic Kubernetes Manifests

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: myapp
  labels:
    name: myapp
```

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: myapp-config
  namespace: myapp
data:
  SERVER_HOST: "0.0.0.0"
  SERVER_PORT: "8080"
  DB_HOST: "postgres-service"
  DB_PORT: "5432"
  DB_DATABASE: "myapp"
  DB_SSLMODE: "require"
  REDIS_HOST: "redis-service"
  REDIS_PORT: "6379"
  LOG_LEVEL: "info"
  LOG_FORMAT: "json"
```

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: myapp-secrets
  namespace: myapp
type: Opaque
data:
  # Base64 encoded values
  DB_USERNAME: bXlhcHA=  # myapp
  DB_PASSWORD: c2VjcmV0  # secret
  REDIS_PASSWORD: cmVkaXNwYXNz  # redispass
```

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  namespace: myapp
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: myapp-config
        - secretRef:
            name: myapp-secrets
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
          timeoutSeconds: 5
          failureThreshold: 3
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
          timeoutSeconds: 3
          failureThreshold: 3
        securityContext:
          runAsNonRoot: true
          runAsUser: 1000
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          capabilities:
            drop:
            - ALL
```

```yaml
# k8s/service.yaml
apiVersion: v1
kind: Service
metadata:
  name: myapp-service
  namespace: myapp
  labels:
    app: myapp
spec:
  selector:
    app: myapp
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: ClusterIP
```

```yaml
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myapp-ingress
  namespace: myapp
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/force-ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - myapp.example.com
    secretName: myapp-tls
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: myapp-service
            port:
              number: 80
```

### Horizontal Pod Autoscaler

```yaml
# k8s/hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: myapp-hpa
  namespace: myapp
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: myapp
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60
    scaleUp:
      stabilizationWindowSeconds: 60
      policies:
      - type: Percent
        value: 50
        periodSeconds: 60
```

## Health Checks and Monitoring

### Health Check Implementation

```go
// health/health.go
package health

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    
    "github.com/go-redis/redis/v8"
)

type HealthChecker struct {
    db    *sql.DB
    redis *redis.Client
}

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Checks    map[string]Check  `json:"checks"`
    Version   string            `json:"version,omitempty"`
    Uptime    string            `json:"uptime,omitempty"`
}

type Check struct {
    Status  string        `json:"status"`
    Message string        `json:"message,omitempty"`
    Latency time.Duration `json:"latency,omitempty"`
}

func NewHealthChecker(db *sql.DB, redis *redis.Client) *HealthChecker {
    return &HealthChecker{
        db:    db,
        redis: redis,
    }
}

func (h *HealthChecker) CheckDatabase(ctx context.Context) Check {
    start := time.Now()
    
    if h.db == nil {
        return Check{
            Status:  "fail",
            Message: "database connection not initialized",
            Latency: time.Since(start),
        }
    }
    
    if err := h.db.PingContext(ctx); err != nil {
        return Check{
            Status:  "fail",
            Message: fmt.Sprintf("database ping failed: %v", err),
            Latency: time.Since(start),
        }
    }
    
    return Check{
        Status:  "pass",
        Latency: time.Since(start),
    }
}

func (h *HealthChecker) CheckRedis(ctx context.Context) Check {
    start := time.Now()
    
    if h.redis == nil {
        return Check{
            Status:  "fail",
            Message: "redis connection not initialized",
            Latency: time.Since(start),
        }
    }
    
    if err := h.redis.Ping(ctx).Err(); err != nil {
        return Check{
            Status:  "fail",
            Message: fmt.Sprintf("redis ping failed: %v", err),
            Latency: time.Since(start),
        }
    }
    
    return Check{
        Status:  "pass",
        Latency: time.Since(start),
    }
}

func (h *HealthChecker) HealthHandler(version string, startTime time.Time) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()
        
        checks := map[string]Check{
            "database": h.CheckDatabase(ctx),
            "redis":    h.CheckRedis(ctx),
        }
        
        // Determine overall status
        status := "pass"
        for _, check := range checks {
            if check.Status == "fail" {
                status = "fail"
                break
            }
        }
        
        health := HealthStatus{
            Status:    status,
            Timestamp: time.Now().UTC(),
            Checks:    checks,
            Version:   version,
            Uptime:    time.Since(startTime).String(),
        }
        
        w.Header().Set("Content-Type", "application/json")
        
        if status == "fail" {
            w.WriteHeader(http.StatusServiceUnavailable)
        } else {
            w.WriteHeader(http.StatusOK)
        }
        
        json.NewEncoder(w).Encode(health)
    }
}

func (h *HealthChecker) ReadinessHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
        defer cancel()
        
        // Quick checks for readiness
        dbCheck := h.CheckDatabase(ctx)
        
        if dbCheck.Status == "fail" {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]string{
                "status": "not ready",
                "reason": "database not available",
            })
            return
        }
        
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "ready",
        })
    }
}

func (h *HealthChecker) LivenessHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Simple liveness check - just return OK if the service is running
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "alive",
        })
    }
}
```

### Metrics and Monitoring

```go
// metrics/metrics.go
package metrics

import (
    "net/http"
    "strconv"
    "time"
    
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
    requestsTotal     *prometheus.CounterVec
    requestDuration   *prometheus.HistogramVec
    activeConnections prometheus.Gauge
    dbConnections     *prometheus.GaugeVec
}

func NewMetrics() *Metrics {
    m := &Metrics{
        requestsTotal: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "http_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name:    "http_request_duration_seconds",
                Help:    "HTTP request duration in seconds",
                Buckets: prometheus.DefBuckets,
            },
            []string{"method", "endpoint"},
        ),
        activeConnections: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "http_active_connections",
                Help: "Number of active HTTP connections",
            },
        ),
        dbConnections: prometheus.NewGaugeVec(
            prometheus.GaugeOpts{
                Name: "database_connections",
                Help: "Number of database connections",
            },
            []string{"state"},
        ),
    }
    
    // Register metrics
    prometheus.MustRegister(
        m.requestsTotal,
        m.requestDuration,
        m.activeConnections,
        m.dbConnections,
    )
    
    return m
}

func (m *Metrics) MiddlewareFunc(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Increment active connections
        m.activeConnections.Inc()
        defer m.activeConnections.Dec()
        
        // Wrap response writer to capture status code
        wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
        
        // Call next handler
        next.ServeHTTP(wrapped, r)
        
        // Record metrics
        duration := time.Since(start).Seconds()
        status := strconv.Itoa(wrapped.statusCode)
        
        m.requestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        m.requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}

func (m *Metrics) UpdateDBConnections(open, inUse, idle int) {
    m.dbConnections.WithLabelValues("open").Set(float64(open))
    m.dbConnections.WithLabelValues("in_use").Set(float64(inUse))
    m.dbConnections.WithLabelValues("idle").Set(float64(idle))
}

func (m *Metrics) Handler() http.Handler {
    return promhttp.Handler()
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## Graceful Shutdown

```go
// server/server.go
package server

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Server struct {
    httpServer *http.Server
    logger     *log.Logger
}

func NewServer(addr string, handler http.Handler, logger *log.Logger) *Server {
    return &Server{
        httpServer: &http.Server{
            Addr:         addr,
            Handler:      handler,
            ReadTimeout:  30 * time.Second,
            WriteTimeout: 30 * time.Second,
            IdleTimeout:  60 * time.Second,
        },
        logger: logger,
    }
}

func (s *Server) Start() error {
    // Start server in a goroutine
    go func() {
        s.logger.Printf("Starting server on %s", s.httpServer.Addr)
        if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            s.logger.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    <-quit
    s.logger.Println("Shutting down server...")
    
    // Create shutdown context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Attempt graceful shutdown
    if err := s.httpServer.Shutdown(ctx); err != nil {
        s.logger.Printf("Server forced to shutdown: %v", err)
        return err
    }
    
    s.logger.Println("Server exited")
    return nil
}

func (s *Server) Stop(ctx context.Context) error {
    return s.httpServer.Shutdown(ctx)
}
```

## CI/CD Pipeline

### GitHub Actions

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Cache Go modules
      uses: actions/cache@v3
      with:
        path: ~/go/pkg/mod
        key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
        restore-keys: |
          ${{ runner.os }}-go-
    
    - name: Download dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v -race -coverprofile=coverage.out ./...
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: ./coverage.out
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest

  build:
    needs: test
    runs-on: ubuntu-latest
    if: github.event_name == 'push'
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Build binary
      run: |
        CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
          -ldflags="-s -w -X main.version=${{ github.sha }} -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
          -o app .
    
    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Extract metadata
      id: meta
      uses: docker/metadata-action@v4
      with:
        images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
        tags: |
          type=ref,event=branch
          type=ref,event=pr
          type=sha,prefix={{branch}}-
          type=raw,value=latest,enable={{is_default_branch}}
    
    - name: Build and push Docker image
      uses: docker/build-push-action@v4
      with:
        context: .
        push: true
        tags: ${{ steps.meta.outputs.tags }}
        labels: ${{ steps.meta.outputs.labels }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
    - uses: actions/checkout@v4
    
    - name: Deploy to Kubernetes
      run: |
        echo "Deploying to production..."
        # Add your deployment commands here
        # kubectl apply -f k8s/
        # kubectl set image deployment/myapp myapp=${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
```

## Exercises

### Exercise 1: Multi-Stage Docker Build
Create a multi-stage Dockerfile that:
- Builds a Go application with proper optimization flags
- Uses a minimal base image for the final stage
- Includes health checks and proper security practices
- Results in an image smaller than 20MB

### Exercise 2: Kubernetes Deployment
Deploy a Go web application to Kubernetes with:
- ConfigMaps and Secrets for configuration
- Horizontal Pod Autoscaler
- Ingress with TLS termination
- Health checks and monitoring

### Exercise 3: CI/CD Pipeline
Set up a complete CI/CD pipeline that:
- Runs tests and linting on every commit
- Builds and pushes Docker images
- Deploys to staging and production environments
- Includes rollback capabilities

## Key Takeaways

- Use multi-stage Docker builds for smaller, more secure images
- Implement proper configuration management with environment variables
- Always include health checks for containerized applications
- Use graceful shutdown to handle termination signals properly
- Implement comprehensive monitoring and metrics collection
- Follow security best practices (non-root users, minimal images)
- Use Kubernetes for orchestration and scaling
- Automate deployment with CI/CD pipelines
- Monitor application performance and resource usage
- Plan for rollbacks and disaster recovery

## Next Steps

Next, we'll explore [Security](26-security.md) to learn about securing Go applications and following security best practices.