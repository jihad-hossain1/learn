# 28. Microservices

This chapter covers building microservices with Go, including service communication, service discovery, load balancing, and distributed system patterns.

## Microservices Architecture

### Basic Service Structure

```go
// service/user/main.go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"
    
    "github.com/gorilla/mux"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Created  time.Time `json:"created"`
}

type UserService struct {
    users map[int]*User
    nextID int
}

func NewUserService() *UserService {
    return &UserService{
        users:  make(map[int]*User),
        nextID: 1,
    }
}

func (us *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    user.ID = us.nextID
    user.Created = time.Now()
    us.nextID++
    
    us.users[user.ID] = &user
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}

func (us *UserService) GetUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    user, exists := us.users[id]
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (us *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
    users := make([]*User, 0, len(us.users))
    for _, user := range us.users {
        users = append(users, user)
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}

func (us *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    user, exists := us.users[id]
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    var updates User
    if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if updates.Username != "" {
        user.Username = updates.Username
    }
    if updates.Email != "" {
        user.Email = updates.Email
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}

func (us *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }
    
    if _, exists := us.users[id]; !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    
    delete(us.users, id)
    w.WriteHeader(http.StatusNoContent)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "status": "healthy",
        "service": "user-service",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

func main() {
    userService := NewUserService()
    
    r := mux.NewRouter()
    
    // API routes
    api := r.PathPrefix("/api/v1").Subrouter()
    api.HandleFunc("/users", userService.CreateUser).Methods("POST")
    api.HandleFunc("/users", userService.GetUsers).Methods("GET")
    api.HandleFunc("/users/{id}", userService.GetUser).Methods("GET")
    api.HandleFunc("/users/{id}", userService.UpdateUser).Methods("PUT")
    api.HandleFunc("/users/{id}", userService.DeleteUser).Methods("DELETE")
    
    // Health check
    r.HandleFunc("/health", healthCheck).Methods("GET")
    
    // Get port from environment
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    server := &http.Server{
        Addr:         ":" + port,
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    // Start server in goroutine
    go func() {
        log.Printf("User service starting on port %s", port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Server failed to start: %v", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}
```

### Service Configuration

```go
// config/config.go
package config

import (
    "os"
    "strconv"
    "time"
)

type Config struct {
    Service  ServiceConfig
    Database DatabaseConfig
    Redis    RedisConfig
    Consul   ConsulConfig
    Jaeger   JaegerConfig
}

type ServiceConfig struct {
    Name    string
    Version string
    Port    int
    Host    string
    Timeout time.Duration
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

type ConsulConfig struct {
    Host   string
    Port   int
    Scheme string
}

type JaegerConfig struct {
    Endpoint string
    Service  string
}

func Load() (*Config, error) {
    config := &Config{
        Service: ServiceConfig{
            Name:    getEnv("SERVICE_NAME", "user-service"),
            Version: getEnv("SERVICE_VERSION", "1.0.0"),
            Port:    getEnvAsInt("PORT", 8080),
            Host:    getEnv("HOST", "localhost"),
            Timeout: getEnvAsDuration("SERVICE_TIMEOUT", 30*time.Second),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnvAsInt("DB_PORT", 5432),
            Username: getEnv("DB_USERNAME", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            Database: getEnv("DB_DATABASE", "userdb"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
        Redis: RedisConfig{
            Host:     getEnv("REDIS_HOST", "localhost"),
            Port:     getEnvAsInt("REDIS_PORT", 6379),
            Password: getEnv("REDIS_PASSWORD", ""),
            Database: getEnvAsInt("REDIS_DATABASE", 0),
        },
        Consul: ConsulConfig{
            Host:   getEnv("CONSUL_HOST", "localhost"),
            Port:   getEnvAsInt("CONSUL_PORT", 8500),
            Scheme: getEnv("CONSUL_SCHEME", "http"),
        },
        Jaeger: JaegerConfig{
            Endpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
            Service:  getEnv("JAEGER_SERVICE", "user-service"),
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
```

## Service Communication

### HTTP Client

```go
// client/http.go
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type HTTPClient struct {
    client  *http.Client
    baseURL string
    timeout time.Duration
}

func NewHTTPClient(baseURL string, timeout time.Duration) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: timeout,
            Transport: &http.Transport{
                MaxIdleConns:        100,
                MaxIdleConnsPerHost: 10,
                IdleConnTimeout:     90 * time.Second,
            },
        },
        baseURL: baseURL,
        timeout: timeout,
    }
}

func (c *HTTPClient) Get(ctx context.Context, path string, result interface{}) error {
    req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+path, nil)
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
    }
    
    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }
    
    return nil
}

func (c *HTTPClient) Post(ctx context.Context, path string, body interface{}, result interface{}) error {
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return err
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, reqBody)
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        respBody, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
    }
    
    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }
    
    return nil
}

func (c *HTTPClient) Put(ctx context.Context, path string, body interface{}, result interface{}) error {
    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return err
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }
    
    req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+path, reqBody)
    if err != nil {
        return err
    }
    
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        respBody, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
    }
    
    if result != nil {
        return json.NewDecoder(resp.Body).Decode(result)
    }
    
    return nil
}

func (c *HTTPClient) Delete(ctx context.Context, path string) error {
    req, err := http.NewRequestWithContext(ctx, "DELETE", c.baseURL+path, nil)
    if err != nil {
        return err
    }
    
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        body, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
    }
    
    return nil
}

// User service client
type UserClient struct {
    client *HTTPClient
}

func NewUserClient(baseURL string) *UserClient {
    return &UserClient{
        client: NewHTTPClient(baseURL, 30*time.Second),
    }
}

func (uc *UserClient) GetUser(ctx context.Context, userID int) (*User, error) {
    var user User
    err := uc.client.Get(ctx, fmt.Sprintf("/api/v1/users/%d", userID), &user)
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (uc *UserClient) CreateUser(ctx context.Context, user *User) (*User, error) {
    var result User
    err := uc.client.Post(ctx, "/api/v1/users", user, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (uc *UserClient) UpdateUser(ctx context.Context, userID int, user *User) (*User, error) {
    var result User
    err := uc.client.Put(ctx, fmt.Sprintf("/api/v1/users/%d", userID), user, &result)
    if err != nil {
        return nil, err
    }
    return &result, nil
}

func (uc *UserClient) DeleteUser(ctx context.Context, userID int) error {
    return uc.client.Delete(ctx, fmt.Sprintf("/api/v1/users/%d", userID))
}
```

### gRPC Communication

```protobuf
// proto/user.proto
syntax = "proto3";

package user;

option go_package = "./proto";

service UserService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}

message User {
    int32 id = 1;
    string username = 2;
    string email = 3;
    int64 created = 4;
}

message GetUserRequest {
    int32 id = 1;
}

message GetUserResponse {
    User user = 1;
}

message CreateUserRequest {
    string username = 1;
    string email = 2;
}

message CreateUserResponse {
    User user = 1;
}

message UpdateUserRequest {
    int32 id = 1;
    string username = 2;
    string email = 3;
}

message UpdateUserResponse {
    User user = 1;
}

message DeleteUserRequest {
    int32 id = 1;
}

message DeleteUserResponse {
    bool success = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
}

message ListUsersResponse {
    repeated User users = 1;
    int32 total = 2;
}
```

```go
// grpc/server.go
package grpc

import (
    "context"
    "fmt"
    "net"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "./proto"
)

type UserServer struct {
    pb.UnimplementedUserServiceServer
    users  map[int32]*pb.User
    nextID int32
}

func NewUserServer() *UserServer {
    return &UserServer{
        users:  make(map[int32]*pb.User),
        nextID: 1,
    }
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
    user, exists := s.users[req.Id]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.Id)
    }
    
    return &pb.GetUserResponse{User: user}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
    user := &pb.User{
        Id:       s.nextID,
        Username: req.Username,
        Email:    req.Email,
        Created:  time.Now().Unix(),
    }
    
    s.users[s.nextID] = user
    s.nextID++
    
    return &pb.CreateUserResponse{User: user}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
    user, exists := s.users[req.Id]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.Id)
    }
    
    if req.Username != "" {
        user.Username = req.Username
    }
    if req.Email != "" {
        user.Email = req.Email
    }
    
    return &pb.UpdateUserResponse{User: user}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
    _, exists := s.users[req.Id]
    if !exists {
        return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.Id)
    }
    
    delete(s.users, req.Id)
    
    return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
    users := make([]*pb.User, 0, len(s.users))
    for _, user := range s.users {
        users = append(users, user)
    }
    
    // Simple pagination
    page := req.Page
    pageSize := req.PageSize
    if page <= 0 {
        page = 1
    }
    if pageSize <= 0 {
        pageSize = 10
    }
    
    start := (page - 1) * pageSize
    end := start + pageSize
    
    if start >= int32(len(users)) {
        return &pb.ListUsersResponse{
            Users: []*pb.User{},
            Total: int32(len(users)),
        }, nil
    }
    
    if end > int32(len(users)) {
        end = int32(len(users))
    }
    
    return &pb.ListUsersResponse{
        Users: users[start:end],
        Total: int32(len(users)),
    }, nil
}

func StartGRPCServer(port int) error {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        return err
    }
    
    s := grpc.NewServer()
    pb.RegisterUserServiceServer(s, NewUserServer())
    
    fmt.Printf("gRPC server listening on port %d\n", port)
    return s.Serve(lis)
}
```

```go
// grpc/client.go
package grpc

import (
    "context"
    "fmt"
    "time"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    
    pb "./proto"
)

type UserGRPCClient struct {
    client pb.UserServiceClient
    conn   *grpc.ClientConn
}

func NewUserGRPCClient(address string) (*UserGRPCClient, error) {
    conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        return nil, err
    }
    
    client := pb.NewUserServiceClient(conn)
    
    return &UserGRPCClient{
        client: client,
        conn:   conn,
    }, nil
}

func (c *UserGRPCClient) Close() error {
    return c.conn.Close()
}

func (c *UserGRPCClient) GetUser(ctx context.Context, userID int32) (*pb.User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    resp, err := c.client.GetUser(ctx, &pb.GetUserRequest{Id: userID})
    if err != nil {
        return nil, err
    }
    
    return resp.User, nil
}

func (c *UserGRPCClient) CreateUser(ctx context.Context, username, email string) (*pb.User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    resp, err := c.client.CreateUser(ctx, &pb.CreateUserRequest{
        Username: username,
        Email:    email,
    })
    if err != nil {
        return nil, err
    }
    
    return resp.User, nil
}

func (c *UserGRPCClient) UpdateUser(ctx context.Context, userID int32, username, email string) (*pb.User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    resp, err := c.client.UpdateUser(ctx, &pb.UpdateUserRequest{
        Id:       userID,
        Username: username,
        Email:    email,
    })
    if err != nil {
        return nil, err
    }
    
    return resp.User, nil
}

func (c *UserGRPCClient) DeleteUser(ctx context.Context, userID int32) error {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    _, err := c.client.DeleteUser(ctx, &pb.DeleteUserRequest{Id: userID})
    return err
}

func (c *UserGRPCClient) ListUsers(ctx context.Context, page, pageSize int32) ([]*pb.User, int32, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    
    resp, err := c.client.ListUsers(ctx, &pb.ListUsersRequest{
        Page:     page,
        PageSize: pageSize,
    })
    if err != nil {
        return nil, 0, err
    }
    
    return resp.Users, resp.Total, nil
}
```

## Service Discovery

### Consul Integration

```go
// discovery/consul.go
package discovery

import (
    "fmt"
    "log"
    "strconv"
    "time"
    
    "github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
    client *api.Client
    config *api.Config
}

func NewConsulRegistry(address string) (*ConsulRegistry, error) {
    config := api.DefaultConfig()
    config.Address = address
    
    client, err := api.NewClient(config)
    if err != nil {
        return nil, err
    }
    
    return &ConsulRegistry{
        client: client,
        config: config,
    }, nil
}

type ServiceInfo struct {
    ID      string
    Name    string
    Address string
    Port    int
    Tags    []string
    Meta    map[string]string
}

func (cr *ConsulRegistry) RegisterService(service ServiceInfo) error {
    registration := &api.AgentServiceRegistration{
        ID:      service.ID,
        Name:    service.Name,
        Address: service.Address,
        Port:    service.Port,
        Tags:    service.Tags,
        Meta:    service.Meta,
        Check: &api.AgentServiceCheck{
            HTTP:                           fmt.Sprintf("http://%s:%d/health", service.Address, service.Port),
            Interval:                       "10s",
            Timeout:                        "3s",
            DeregisterCriticalServiceAfter: "30s",
        },
    }
    
    return cr.client.Agent().ServiceRegister(registration)
}

func (cr *ConsulRegistry) DeregisterService(serviceID string) error {
    return cr.client.Agent().ServiceDeregister(serviceID)
}

func (cr *ConsulRegistry) DiscoverServices(serviceName string) ([]ServiceInfo, error) {
    services, _, err := cr.client.Health().Service(serviceName, "", true, nil)
    if err != nil {
        return nil, err
    }
    
    var result []ServiceInfo
    for _, service := range services {
        result = append(result, ServiceInfo{
            ID:      service.Service.ID,
            Name:    service.Service.Service,
            Address: service.Service.Address,
            Port:    service.Service.Port,
            Tags:    service.Service.Tags,
            Meta:    service.Service.Meta,
        })
    }
    
    return result, nil
}

func (cr *ConsulRegistry) WatchServices(serviceName string, callback func([]ServiceInfo)) error {
    plan, err := api.ParseWatchPlan(map[string]interface{}{
        "type":    "service",
        "service": serviceName,
    })
    if err != nil {
        return err
    }
    
    plan.Handler = func(idx uint64, data interface{}) {
        services, ok := data.([]*api.ServiceEntry)
        if !ok {
            log.Printf("Invalid data type for service watch")
            return
        }
        
        var result []ServiceInfo
        for _, service := range services {
            result = append(result, ServiceInfo{
                ID:      service.Service.ID,
                Name:    service.Service.Service,
                Address: service.Service.Address,
                Port:    service.Service.Port,
                Tags:    service.Service.Tags,
                Meta:    service.Service.Meta,
            })
        }
        
        callback(result)
    }
    
    go func() {
        if err := plan.Run(cr.config.Address); err != nil {
            log.Printf("Error running watch plan: %v", err)
        }
    }()
    
    return nil
}

// Service registry with automatic registration
type ServiceRegistry struct {
    consul  *ConsulRegistry
    service ServiceInfo
}

func NewServiceRegistry(consulAddress string, service ServiceInfo) (*ServiceRegistry, error) {
    consul, err := NewConsulRegistry(consulAddress)
    if err != nil {
        return nil, err
    }
    
    return &ServiceRegistry{
        consul:  consul,
        service: service,
    }, nil
}

func (sr *ServiceRegistry) Start() error {
    // Register service
    if err := sr.consul.RegisterService(sr.service); err != nil {
        return err
    }
    
    log.Printf("Service %s registered with Consul", sr.service.Name)
    
    // Start health check updater
    go sr.healthCheckUpdater()
    
    return nil
}

func (sr *ServiceRegistry) Stop() error {
    return sr.consul.DeregisterService(sr.service.ID)
}

func (sr *ServiceRegistry) healthCheckUpdater() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // Update TTL check
        checkID := "service:" + sr.service.ID
        if err := sr.consul.client.Agent().UpdateTTL(checkID, "Service is healthy", api.HealthPassing); err != nil {
            log.Printf("Failed to update TTL check: %v", err)
        }
    }
}

func (sr *ServiceRegistry) DiscoverService(serviceName string) ([]ServiceInfo, error) {
    return sr.consul.DiscoverServices(serviceName)
}
```

### Load Balancing

```go
// loadbalancer/roundrobin.go
package loadbalancer

import (
    "context"
    "errors"
    "sync"
    "sync/atomic"
)

type LoadBalancer interface {
    Next() (string, error)
    Add(address string)
    Remove(address string)
    Healthy(address string, healthy bool)
}

// Round Robin Load Balancer
type RoundRobinLB struct {
    addresses []string
    current   uint64
    mutex     sync.RWMutex
}

func NewRoundRobinLB() *RoundRobinLB {
    return &RoundRobinLB{
        addresses: make([]string, 0),
    }
}

func (lb *RoundRobinLB) Next() (string, error) {
    lb.mutex.RLock()
    defer lb.mutex.RUnlock()
    
    if len(lb.addresses) == 0 {
        return "", errors.New("no available addresses")
    }
    
    index := atomic.AddUint64(&lb.current, 1) % uint64(len(lb.addresses))
    return lb.addresses[index], nil
}

func (lb *RoundRobinLB) Add(address string) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    // Check if address already exists
    for _, addr := range lb.addresses {
        if addr == address {
            return
        }
    }
    
    lb.addresses = append(lb.addresses, address)
}

func (lb *RoundRobinLB) Remove(address string) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    for i, addr := range lb.addresses {
        if addr == address {
            lb.addresses = append(lb.addresses[:i], lb.addresses[i+1:]...)
            break
        }
    }
}

func (lb *RoundRobinLB) Healthy(address string, healthy bool) {
    if healthy {
        lb.Add(address)
    } else {
        lb.Remove(address)
    }
}

// Weighted Round Robin Load Balancer
type WeightedRoundRobinLB struct {
    servers []WeightedServer
    current int
    mutex   sync.RWMutex
}

type WeightedServer struct {
    Address       string
    Weight        int
    CurrentWeight int
    Healthy       bool
}

func NewWeightedRoundRobinLB() *WeightedRoundRobinLB {
    return &WeightedRoundRobinLB{
        servers: make([]WeightedServer, 0),
    }
}

func (lb *WeightedRoundRobinLB) AddWeighted(address string, weight int) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    // Check if server already exists
    for i, server := range lb.servers {
        if server.Address == address {
            lb.servers[i].Weight = weight
            lb.servers[i].Healthy = true
            return
        }
    }
    
    lb.servers = append(lb.servers, WeightedServer{
        Address:       address,
        Weight:        weight,
        CurrentWeight: 0,
        Healthy:       true,
    })
}

func (lb *WeightedRoundRobinLB) Next() (string, error) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    if len(lb.servers) == 0 {
        return "", errors.New("no available servers")
    }
    
    // Find healthy servers
    var healthyServers []int
    totalWeight := 0
    
    for i, server := range lb.servers {
        if server.Healthy {
            healthyServers = append(healthyServers, i)
            totalWeight += server.Weight
        }
    }
    
    if len(healthyServers) == 0 {
        return "", errors.New("no healthy servers")
    }
    
    // Weighted round robin algorithm
    var selected int = -1
    maxCurrentWeight := -1
    
    for _, i := range healthyServers {
        lb.servers[i].CurrentWeight += lb.servers[i].Weight
        
        if lb.servers[i].CurrentWeight > maxCurrentWeight {
            maxCurrentWeight = lb.servers[i].CurrentWeight
            selected = i
        }
    }
    
    if selected != -1 {
        lb.servers[selected].CurrentWeight -= totalWeight
        return lb.servers[selected].Address, nil
    }
    
    return "", errors.New("no server selected")
}

func (lb *WeightedRoundRobinLB) Add(address string) {
    lb.AddWeighted(address, 1) // Default weight
}

func (lb *WeightedRoundRobinLB) Remove(address string) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    for i, server := range lb.servers {
        if server.Address == address {
            lb.servers = append(lb.servers[:i], lb.servers[i+1:]...)
            break
        }
    }
}

func (lb *WeightedRoundRobinLB) Healthy(address string, healthy bool) {
    lb.mutex.Lock()
    defer lb.mutex.Unlock()
    
    for i, server := range lb.servers {
        if server.Address == address {
            lb.servers[i].Healthy = healthy
            break
        }
    }
}

// Service client with load balancing
type ServiceClient struct {
    lb         LoadBalancer
    httpClient *HTTPClient
}

func NewServiceClient(lb LoadBalancer) *ServiceClient {
    return &ServiceClient{
        lb:         lb,
        httpClient: NewHTTPClient("", 30*time.Second),
    }
}

func (sc *ServiceClient) Call(ctx context.Context, method, path string, body, result interface{}) error {
    address, err := sc.lb.Next()
    if err != nil {
        return err
    }
    
    // Update base URL for this request
    sc.httpClient.baseURL = "http://" + address
    
    switch method {
    case "GET":
        return sc.httpClient.Get(ctx, path, result)
    case "POST":
        return sc.httpClient.Post(ctx, path, body, result)
    case "PUT":
        return sc.httpClient.Put(ctx, path, body, result)
    case "DELETE":
        return sc.httpClient.Delete(ctx, path)
    default:
        return errors.New("unsupported HTTP method")
    }
}
```

## Circuit Breaker Pattern

```go
// circuitbreaker/breaker.go
package circuitbreaker

import (
    "errors"
    "sync"
    "time"
)

type State int

const (
    StateClosed State = iota
    StateHalfOpen
    StateOpen
)

type CircuitBreaker struct {
    maxRequests  uint32
    interval     time.Duration
    timeout      time.Duration
    readyToTrip  func(counts Counts) bool
    onStateChange func(name string, from State, to State)
    
    mutex      sync.Mutex
    name       string
    state      State
    generation uint64
    counts     Counts
    expiry     time.Time
}

type Counts struct {
    Requests             uint32
    TotalSuccesses       uint32
    TotalFailures        uint32
    ConsecutiveSuccesses uint32
    ConsecutiveFailures  uint32
}

func (c Counts) SuccessRate() float64 {
    if c.Requests == 0 {
        return 0
    }
    return float64(c.TotalSuccesses) / float64(c.Requests)
}

func (c Counts) FailureRate() float64 {
    if c.Requests == 0 {
        return 0
    }
    return float64(c.TotalFailures) / float64(c.Requests)
}

type Settings struct {
    Name          string
    MaxRequests   uint32
    Interval      time.Duration
    Timeout       time.Duration
    ReadyToTrip   func(counts Counts) bool
    OnStateChange func(name string, from State, to State)
}

func NewCircuitBreaker(st Settings) *CircuitBreaker {
    cb := &CircuitBreaker{
        name:         st.Name,
        maxRequests:  st.MaxRequests,
        interval:     st.Interval,
        timeout:      st.Timeout,
        readyToTrip:  st.ReadyToTrip,
        onStateChange: st.OnStateChange,
        state:        StateClosed,
        expiry:       time.Now().Add(st.Interval),
    }
    
    if cb.maxRequests == 0 {
        cb.maxRequests = 1
    }
    
    if cb.interval == 0 {
        cb.interval = time.Duration(0)
    }
    
    if cb.timeout == 0 {
        cb.timeout = 60 * time.Second
    }
    
    if cb.readyToTrip == nil {
        cb.readyToTrip = func(counts Counts) bool {
            return counts.ConsecutiveFailures > 5
        }
    }
    
    return cb
}

var (
    ErrTooManyRequests = errors.New("circuit breaker: too many requests")
    ErrOpenState       = errors.New("circuit breaker: open state")
)

func (cb *CircuitBreaker) Execute(req func() (interface{}, error)) (interface{}, error) {
    generation, err := cb.beforeRequest()
    if err != nil {
        return nil, err
    }
    
    defer func() {
        e := recover()
        if e != nil {
            cb.afterRequest(generation, false)
            panic(e)
        }
    }()
    
    result, err := req()
    cb.afterRequest(generation, err == nil)
    return result, err
}

func (cb *CircuitBreaker) beforeRequest() (uint64, error) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    now := time.Now()
    state, generation := cb.currentState(now)
    
    if state == StateOpen {
        return generation, ErrOpenState
    } else if state == StateHalfOpen && cb.counts.Requests >= cb.maxRequests {
        return generation, ErrTooManyRequests
    }
    
    cb.counts.Requests++
    return generation, nil
}

func (cb *CircuitBreaker) afterRequest(before uint64, success bool) {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    now := time.Now()
    state, generation := cb.currentState(now)
    if generation != before {
        return
    }
    
    if success {
        cb.onSuccess(state, now)
    } else {
        cb.onFailure(state, now)
    }
}

func (cb *CircuitBreaker) onSuccess(state State, now time.Time) {
    switch state {
    case StateClosed:
        cb.counts.TotalSuccesses++
        cb.counts.ConsecutiveSuccesses++
        cb.counts.ConsecutiveFailures = 0
    case StateHalfOpen:
        cb.counts.TotalSuccesses++
        cb.counts.ConsecutiveSuccesses++
        cb.counts.ConsecutiveFailures = 0
        if cb.counts.ConsecutiveSuccesses >= cb.maxRequests {
            cb.setState(StateClosed, now)
        }
    }
}

func (cb *CircuitBreaker) onFailure(state State, now time.Time) {
    switch state {
    case StateClosed:
        cb.counts.TotalFailures++
        cb.counts.ConsecutiveFailures++
        cb.counts.ConsecutiveSuccesses = 0
        if cb.readyToTrip(cb.counts) {
            cb.setState(StateOpen, now)
        }
    case StateHalfOpen:
        cb.setState(StateOpen, now)
    }
}

func (cb *CircuitBreaker) currentState(now time.Time) (State, uint64) {
    switch cb.state {
    case StateClosed:
        if !cb.expiry.IsZero() && cb.expiry.Before(now) {
            cb.toNewGeneration(now)
        }
    case StateOpen:
        if cb.expiry.Before(now) {
            cb.setState(StateHalfOpen, now)
        }
    }
    return cb.state, cb.generation
}

func (cb *CircuitBreaker) setState(state State, now time.Time) {
    if cb.state == state {
        return
    }
    
    prev := cb.state
    cb.state = state
    
    cb.toNewGeneration(now)
    
    if cb.onStateChange != nil {
        cb.onStateChange(cb.name, prev, state)
    }
}

func (cb *CircuitBreaker) toNewGeneration(now time.Time) {
    cb.generation++
    cb.counts = Counts{}
    
    var zero time.Time
    switch cb.state {
    case StateClosed:
        if cb.interval == 0 {
            cb.expiry = zero
        } else {
            cb.expiry = now.Add(cb.interval)
        }
    case StateOpen:
        cb.expiry = now.Add(cb.timeout)
    default: // StateHalfOpen
        cb.expiry = zero
    }
}

func (cb *CircuitBreaker) State() State {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    state, _ := cb.currentState(time.Now())
    return state
}

func (cb *CircuitBreaker) Counts() Counts {
    cb.mutex.Lock()
    defer cb.mutex.Unlock()
    
    return cb.counts
}
```

## Distributed Tracing

### OpenTelemetry Integration

```go
// tracing/tracer.go
package tracing

import (
    "context"
    "log"
    
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
    "go.opentelemetry.io/otel/trace"
)

func InitTracer(serviceName, jaegerEndpoint string) func() {
    // Create Jaeger exporter
    exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
    if err != nil {
        log.Fatalf("Failed to create Jaeger exporter: %v", err)
    }
    
    // Create resource
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String(serviceName),
            semconv.ServiceVersionKey.String("1.0.0"),
        ),
    )
    if err != nil {
        log.Fatalf("Failed to create resource: %v", err)
    }
    
    // Create trace provider
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exp),
        trace.WithResource(res),
        trace.WithSampler(trace.AlwaysSample()),
    )
    
    // Set global trace provider
    otel.SetTracerProvider(tp)
    
    // Set global propagator
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{},
        propagation.Baggage{},
    ))
    
    return func() {
        if err := tp.Shutdown(context.Background()); err != nil {
            log.Printf("Error shutting down tracer provider: %v", err)
        }
    }
}

// HTTP middleware for tracing
func TracingMiddleware(serviceName string) func(http.Handler) http.Handler {
    tracer := otel.Tracer(serviceName)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
            
            ctx, span := tracer.Start(ctx, r.Method+" "+r.URL.Path,
                trace.WithAttributes(
                    semconv.HTTPMethodKey.String(r.Method),
                    semconv.HTTPURLKey.String(r.URL.String()),
                    semconv.HTTPUserAgentKey.String(r.UserAgent()),
                ),
            )
            defer span.End()
            
            // Inject trace context into response headers
            otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(w.Header()))
            
            // Wrap response writer to capture status code
            wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
            
            next.ServeHTTP(wrapped, r.WithContext(ctx))
            
            span.SetAttributes(
                semconv.HTTPStatusCodeKey.Int(wrapped.statusCode),
            )
            
            if wrapped.statusCode >= 400 {
                span.RecordError(fmt.Errorf("HTTP %d", wrapped.statusCode))
            }
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

// HTTP client with tracing
type TracingHTTPClient struct {
    client *http.Client
    tracer trace.Tracer
}

func NewTracingHTTPClient(serviceName string) *TracingHTTPClient {
    return &TracingHTTPClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
        tracer: otel.Tracer(serviceName),
    }
}

func (c *TracingHTTPClient) Do(req *http.Request) (*http.Response, error) {
    ctx, span := c.tracer.Start(req.Context(), "HTTP "+req.Method,
        trace.WithAttributes(
            semconv.HTTPMethodKey.String(req.Method),
            semconv.HTTPURLKey.String(req.URL.String()),
        ),
    )
    defer span.End()
    
    // Inject trace context into request headers
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
    
    req = req.WithContext(ctx)
    resp, err := c.client.Do(req)
    
    if err != nil {
        span.RecordError(err)
        return nil, err
    }
    
    span.SetAttributes(
        semconv.HTTPStatusCodeKey.Int(resp.StatusCode),
    )
    
    if resp.StatusCode >= 400 {
        span.RecordError(fmt.Errorf("HTTP %d", resp.StatusCode))
    }
    
    return resp, nil
}

// Database tracing
func TraceDBOperation(ctx context.Context, operation, query string, fn func() error) error {
    tracer := otel.Tracer("database")
    
    ctx, span := tracer.Start(ctx, operation,
        trace.WithAttributes(
            attribute.String("db.operation", operation),
            attribute.String("db.statement", query),
        ),
    )
    defer span.End()
    
    err := fn()
    if err != nil {
        span.RecordError(err)
    }
    
    return err
}

// Custom span creation
func StartSpan(ctx context.Context, serviceName, operationName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
    tracer := otel.Tracer(serviceName)
    return tracer.Start(ctx, operationName, trace.WithAttributes(attrs...))
}

// Add attributes to current span
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
    span := trace.SpanFromContext(ctx)
    span.SetAttributes(attrs...)
}

// Record error in current span
func RecordError(ctx context.Context, err error) {
    span := trace.SpanFromContext(ctx)
    span.RecordError(err)
}
```

## Exercises

### Exercise 1: Build a Complete Microservice
Create a complete microservice with:
- RESTful API endpoints
- Database integration
- Service registration with Consul
- Health checks
- Graceful shutdown
- Docker containerization

### Exercise 2: Service Communication
Implement communication between multiple services:
- HTTP client with retry logic
- gRPC communication
- Circuit breaker pattern
- Load balancing
- Error handling and timeouts

### Exercise 3: Distributed System
Build a distributed system with:
- Multiple microservices
- Service discovery
- Distributed tracing
- Centralized logging
- Monitoring and metrics
- API gateway

## Key Takeaways

- Design services around business capabilities
- Keep services small and focused (single responsibility)
- Use service discovery for dynamic service location
- Implement circuit breakers to prevent cascade failures
- Use distributed tracing for observability
- Design for failure - services will fail
- Implement proper health checks and monitoring
- Use asynchronous communication when possible
- Maintain service independence and loose coupling
- Version your APIs and maintain backward compatibility
- Implement proper security between services
- Use containers for consistent deployment

## Next Steps

Next, we'll explore [GraphQL](29-graphql.md) to learn about building GraphQL APIs and integrating them with Go applications.