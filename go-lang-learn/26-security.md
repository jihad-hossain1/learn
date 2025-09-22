# 26. Security

This chapter covers security best practices for Go applications, including input validation, authentication, authorization, cryptography, and protection against common vulnerabilities.

## Input Validation and Sanitization

### Basic Input Validation

```go
// validation/validator.go
package validation

import (
    "fmt"
    "net/mail"
    "regexp"
    "strings"
    "unicode"
)

type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
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

func (e ValidationErrors) HasErrors() bool {
    return len(e) > 0
}

// String validation
func ValidateRequired(value, fieldName string) *ValidationError {
    if strings.TrimSpace(value) == "" {
        return &ValidationError{
            Field:   fieldName,
            Message: "is required",
        }
    }
    return nil
}

func ValidateMinLength(value, fieldName string, minLength int) *ValidationError {
    if len(value) < minLength {
        return &ValidationError{
            Field:   fieldName,
            Message: fmt.Sprintf("must be at least %d characters long", minLength),
        }
    }
    return nil
}

func ValidateMaxLength(value, fieldName string, maxLength int) *ValidationError {
    if len(value) > maxLength {
        return &ValidationError{
            Field:   fieldName,
            Message: fmt.Sprintf("must be no more than %d characters long", maxLength),
        }
    }
    return nil
}

func ValidateEmail(email, fieldName string) *ValidationError {
    if _, err := mail.ParseAddress(email); err != nil {
        return &ValidationError{
            Field:   fieldName,
            Message: "must be a valid email address",
        }
    }
    return nil
}

func ValidatePassword(password, fieldName string) *ValidationError {
    if len(password) < 8 {
        return &ValidationError{
            Field:   fieldName,
            Message: "must be at least 8 characters long",
        }
    }
    
    var (
        hasUpper   = false
        hasLower   = false
        hasNumber  = false
        hasSpecial = false
    )
    
    for _, char := range password {
        switch {
        case unicode.IsUpper(char):
            hasUpper = true
        case unicode.IsLower(char):
            hasLower = true
        case unicode.IsNumber(char):
            hasNumber = true
        case unicode.IsPunct(char) || unicode.IsSymbol(char):
            hasSpecial = true
        }
    }
    
    if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
        return &ValidationError{
            Field:   fieldName,
            Message: "must contain at least one uppercase letter, one lowercase letter, one number, and one special character",
        }
    }
    
    return nil
}

func ValidateAlphanumeric(value, fieldName string) *ValidationError {
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9]+$`, value)
    if !matched {
        return &ValidationError{
            Field:   fieldName,
            Message: "must contain only letters and numbers",
        }
    }
    return nil
}

func ValidatePhoneNumber(phone, fieldName string) *ValidationError {
    // Simple phone number validation (adjust regex as needed)
    matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
    if !matched {
        return &ValidationError{
            Field:   fieldName,
            Message: "must be a valid phone number",
        }
    }
    return nil
}

// Sanitization functions
func SanitizeString(input string) string {
    // Remove null bytes and control characters
    cleaned := strings.Map(func(r rune) rune {
        if r == 0 || (r < 32 && r != '\t' && r != '\n' && r != '\r') {
            return -1
        }
        return r
    }, input)
    
    return strings.TrimSpace(cleaned)
}

func SanitizeHTML(input string) string {
    // Basic HTML entity encoding
    replacer := strings.NewReplacer(
        "&", "&amp;",
        "<", "&lt;",
        ">", "&gt;",
        "\"", "&quot;",
        "'", "&#39;",
    )
    return replacer.Replace(input)
}

func SanitizeSQL(input string) string {
    // Remove or escape SQL injection attempts
    dangerous := []string{
        "'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
        "DROP", "DELETE", "INSERT", "UPDATE", "SELECT",
        "UNION", "OR", "AND", "WHERE",
    }
    
    cleaned := input
    for _, danger := range dangerous {
        cleaned = strings.ReplaceAll(cleaned, danger, "")
        cleaned = strings.ReplaceAll(cleaned, strings.ToLower(danger), "")
    }
    
    return cleaned
}
```

### Struct Validation

```go
// models/user.go
package models

import (
    "myapp/validation"
)

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"-"` // Never include in JSON
    Phone    string `json:"phone,omitempty"`
}

func (u *User) Validate() validation.ValidationErrors {
    var errors validation.ValidationErrors
    
    // Username validation
    if err := validation.ValidateRequired(u.Username, "username"); err != nil {
        errors = append(errors, *err)
    } else {
        if err := validation.ValidateMinLength(u.Username, "username", 3); err != nil {
            errors = append(errors, *err)
        }
        if err := validation.ValidateMaxLength(u.Username, "username", 50); err != nil {
            errors = append(errors, *err)
        }
        if err := validation.ValidateAlphanumeric(u.Username, "username"); err != nil {
            errors = append(errors, *err)
        }
    }
    
    // Email validation
    if err := validation.ValidateRequired(u.Email, "email"); err != nil {
        errors = append(errors, *err)
    } else {
        if err := validation.ValidateEmail(u.Email, "email"); err != nil {
            errors = append(errors, *err)
        }
    }
    
    // Password validation (only if provided)
    if u.Password != "" {
        if err := validation.ValidatePassword(u.Password, "password"); err != nil {
            errors = append(errors, *err)
        }
    }
    
    // Phone validation (optional)
    if u.Phone != "" {
        if err := validation.ValidatePhoneNumber(u.Phone, "phone"); err != nil {
            errors = append(errors, *err)
        }
    }
    
    return errors
}

func (u *User) Sanitize() {
    u.Username = validation.SanitizeString(u.Username)
    u.Email = validation.SanitizeString(u.Email)
    u.Phone = validation.SanitizeString(u.Phone)
}
```

## Authentication and Authorization

### JWT Authentication

```go
// auth/jwt.go
package auth

import (
    "errors"
    "fmt"
    "time"
    
    "github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
    secretKey     string
    tokenDuration time.Duration
}

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
    return &JWTManager{
        secretKey:     secretKey,
        tokenDuration: tokenDuration,
    }
}

func (manager *JWTManager) Generate(userID int, username, role string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.tokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "myapp",
            Subject:   fmt.Sprintf("%d", userID),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(manager.secretKey))
}

func (manager *JWTManager) Verify(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return []byte(manager.secretKey), nil
        },
    )
    
    if err != nil {
        return nil, err
    }
    
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return claims, nil
}

func (manager *JWTManager) Refresh(tokenString string) (string, error) {
    claims, err := manager.Verify(tokenString)
    if err != nil {
        return "", err
    }
    
    // Check if token is close to expiry (within 1 hour)
    if time.Until(claims.ExpiresAt.Time) > time.Hour {
        return "", errors.New("token doesn't need refresh yet")
    }
    
    return manager.Generate(claims.UserID, claims.Username, claims.Role)
}
```

### Password Hashing

```go
// auth/password.go
package auth

import (
    "crypto/rand"
    "crypto/subtle"
    "encoding/base64"
    "errors"
    "fmt"
    "strings"
    
    "golang.org/x/crypto/argon2"
)

type PasswordConfig struct {
    Memory      uint32
    Iterations  uint32
    Parallelism uint8
    SaltLength  uint32
    KeyLength   uint32
}

func DefaultPasswordConfig() *PasswordConfig {
    return &PasswordConfig{
        Memory:      64 * 1024, // 64 MB
        Iterations:  3,
        Parallelism: 2,
        SaltLength:  16,
        KeyLength:   32,
    }
}

func HashPassword(password string, config *PasswordConfig) (string, error) {
    if config == nil {
        config = DefaultPasswordConfig()
    }
    
    // Generate random salt
    salt := make([]byte, config.SaltLength)
    if _, err := rand.Read(salt); err != nil {
        return "", err
    }
    
    // Generate hash
    hash := argon2.IDKey(
        []byte(password),
        salt,
        config.Iterations,
        config.Memory,
        config.Parallelism,
        config.KeyLength,
    )
    
    // Encode to base64
    b64Salt := base64.RawStdEncoding.EncodeToString(salt)
    b64Hash := base64.RawStdEncoding.EncodeToString(hash)
    
    // Format: $argon2id$v=19$m=65536,t=3,p=2$salt$hash
    format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
    return fmt.Sprintf(format, argon2.Version, config.Memory, config.Iterations, config.Parallelism, b64Salt, b64Hash), nil
}

func VerifyPassword(password, hashedPassword string) (bool, error) {
    parts := strings.Split(hashedPassword, "$")
    if len(parts) != 6 {
        return false, errors.New("invalid hash format")
    }
    
    var version int
    if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
        return false, err
    }
    if version != argon2.Version {
        return false, errors.New("incompatible version")
    }
    
    var memory, iterations uint32
    var parallelism uint8
    if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism); err != nil {
        return false, err
    }
    
    salt, err := base64.RawStdEncoding.DecodeString(parts[4])
    if err != nil {
        return false, err
    }
    
    hash, err := base64.RawStdEncoding.DecodeString(parts[5])
    if err != nil {
        return false, err
    }
    
    // Generate hash for comparison
    comparisonHash := argon2.IDKey(
        []byte(password),
        salt,
        iterations,
        memory,
        parallelism,
        uint32(len(hash)),
    )
    
    // Use constant-time comparison
    return subtle.ConstantTimeCompare(hash, comparisonHash) == 1, nil
}
```

### Middleware for Authentication

```go
// middleware/auth.go
package middleware

import (
    "context"
    "encoding/json"
    "net/http"
    "strings"
    
    "myapp/auth"
)

type contextKey string

const (
    UserContextKey contextKey = "user"
)

type AuthMiddleware struct {
    jwtManager *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
    return &AuthMiddleware{
        jwtManager: jwtManager,
    }
}

func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get token from Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header required", http.StatusUnauthorized)
            return
        }
        
        // Check Bearer token format
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
            return
        }
        
        // Verify token
        claims, err := m.jwtManager.Verify(parts[1])
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }
        
        // Add user info to context
        ctx := context.WithValue(r.Context(), UserContextKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, ok := r.Context().Value(UserContextKey).(*auth.Claims)
            if !ok {
                http.Error(w, "User not authenticated", http.StatusUnauthorized)
                return
            }
            
            if claims.Role != role {
                http.Error(w, "Insufficient permissions", http.StatusForbidden)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader != "" {
            parts := strings.SplitN(authHeader, " ", 2)
            if len(parts) == 2 && parts[0] == "Bearer" {
                if claims, err := m.jwtManager.Verify(parts[1]); err == nil {
                    ctx := context.WithValue(r.Context(), UserContextKey, claims)
                    r = r.WithContext(ctx)
                }
            }
        }
        
        next.ServeHTTP(w, r)
    })
}

// Helper function to get user from context
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
    claims, ok := ctx.Value(UserContextKey).(*auth.Claims)
    return claims, ok
}
```

## Cryptography

### Encryption and Decryption

```go
// crypto/encryption.go
package crypto

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "errors"
    "io"
    
    "golang.org/x/crypto/pbkdf2"
)

type AESEncryption struct {
    key []byte
}

func NewAESEncryption(password string, salt []byte) *AESEncryption {
    // Derive key from password using PBKDF2
    key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
    return &AESEncryption{key: key}
}

func (e *AESEncryption) Encrypt(plaintext []byte) (string, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }
    
    // Create GCM mode
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }
    
    // Generate random nonce
    nonce := make([]byte, gcm.NonceSize())
    if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }
    
    // Encrypt and authenticate
    ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
    
    // Encode to base64
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (e *AESEncryption) Decrypt(ciphertextB64 string) ([]byte, error) {
    // Decode from base64
    ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
    if err != nil {
        return nil, err
    }
    
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return nil, err
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, err
    }
    
    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }
    
    nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
    
    // Decrypt and verify
    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }
    
    return plaintext, nil
}
```

### Digital Signatures

```go
// crypto/signature.go
package crypto

import (
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha256"
    "crypto/x509"
    "encoding/base64"
    "encoding/pem"
    "errors"
)

type RSASignature struct {
    privateKey *rsa.PrivateKey
    publicKey  *rsa.PublicKey
}

func NewRSASignature(keySize int) (*RSASignature, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
    if err != nil {
        return nil, err
    }
    
    return &RSASignature{
        privateKey: privateKey,
        publicKey:  &privateKey.PublicKey,
    }, nil
}

func (r *RSASignature) LoadPrivateKey(pemData []byte) error {
    block, _ := pem.Decode(pemData)
    if block == nil {
        return errors.New("failed to decode PEM block")
    }
    
    privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return err
    }
    
    r.privateKey = privateKey
    r.publicKey = &privateKey.PublicKey
    return nil
}

func (r *RSASignature) LoadPublicKey(pemData []byte) error {
    block, _ := pem.Decode(pemData)
    if block == nil {
        return errors.New("failed to decode PEM block")
    }
    
    publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
    if err != nil {
        return err
    }
    
    r.publicKey = publicKey
    return nil
}

func (r *RSASignature) Sign(data []byte) (string, error) {
    if r.privateKey == nil {
        return "", errors.New("private key not loaded")
    }
    
    hash := sha256.Sum256(data)
    signature, err := rsa.SignPKCS1v15(rand.Reader, r.privateKey, crypto.SHA256, hash[:])
    if err != nil {
        return "", err
    }
    
    return base64.StdEncoding.EncodeToString(signature), nil
}

func (r *RSASignature) Verify(data []byte, signatureB64 string) error {
    if r.publicKey == nil {
        return errors.New("public key not loaded")
    }
    
    signature, err := base64.StdEncoding.DecodeString(signatureB64)
    if err != nil {
        return err
    }
    
    hash := sha256.Sum256(data)
    return rsa.VerifyPKCS1v15(r.publicKey, crypto.SHA256, hash[:], signature)
}

func (r *RSASignature) ExportPrivateKey() ([]byte, error) {
    if r.privateKey == nil {
        return nil, errors.New("private key not available")
    }
    
    privateKeyBytes := x509.MarshalPKCS1PrivateKey(r.privateKey)
    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: privateKeyBytes,
    })
    
    return privateKeyPEM, nil
}

func (r *RSASignature) ExportPublicKey() ([]byte, error) {
    if r.publicKey == nil {
        return nil, errors.New("public key not available")
    }
    
    publicKeyBytes := x509.MarshalPKCS1PublicKey(r.publicKey)
    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PUBLIC KEY",
        Bytes: publicKeyBytes,
    })
    
    return publicKeyPEM, nil
}
```

## Security Headers and CORS

### Security Middleware

```go
// middleware/security.go
package middleware

import (
    "net/http"
    "strings"
)

type SecurityConfig struct {
    ContentSecurityPolicy string
    AllowedOrigins       []string
    AllowedMethods       []string
    AllowedHeaders       []string
    AllowCredentials     bool
    MaxAge               int
}

func DefaultSecurityConfig() *SecurityConfig {
    return &SecurityConfig{
        ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'",
        AllowedOrigins:       []string{"*"},
        AllowedMethods:       []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:       []string{"Origin", "Content-Type", "Accept", "Authorization"},
        AllowCredentials:     false,
        MaxAge:               86400, // 24 hours
    }
}

func SecurityHeaders(config *SecurityConfig) func(http.Handler) http.Handler {
    if config == nil {
        config = DefaultSecurityConfig()
    }
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Security headers
            w.Header().Set("X-Content-Type-Options", "nosniff")
            w.Header().Set("X-Frame-Options", "DENY")
            w.Header().Set("X-XSS-Protection", "1; mode=block")
            w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
            w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
            
            if config.ContentSecurityPolicy != "" {
                w.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func CORS(config *SecurityConfig) func(http.Handler) http.Handler {
    if config == nil {
        config = DefaultSecurityConfig()
    }
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            
            // Check if origin is allowed
            allowed := false
            for _, allowedOrigin := range config.AllowedOrigins {
                if allowedOrigin == "*" || allowedOrigin == origin {
                    allowed = true
                    break
                }
            }
            
            if allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
            }
            
            w.Header().Set("Access-Control-Allow-Methods", strings.Join(config.AllowedMethods, ", "))
            w.Header().Set("Access-Control-Allow-Headers", strings.Join(config.AllowedHeaders, ", "))
            
            if config.AllowCredentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            
            if config.MaxAge > 0 {
                w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
            }
            
            // Handle preflight requests
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## Rate Limiting

### Token Bucket Rate Limiter

```go
// middleware/ratelimit.go
package middleware

import (
    "fmt"
    "net"
    "net/http"
    "sync"
    "time"
)

type TokenBucket struct {
    capacity     int
    tokens       int
    refillRate   int
    lastRefill   time.Time
    mutex        sync.Mutex
}

func NewTokenBucket(capacity, refillRate int) *TokenBucket {
    return &TokenBucket{
        capacity:   capacity,
        tokens:     capacity,
        refillRate: refillRate,
        lastRefill: time.Now(),
    }
}

func (tb *TokenBucket) Allow() bool {
    tb.mutex.Lock()
    defer tb.mutex.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(tb.lastRefill)
    
    // Refill tokens based on elapsed time
    tokensToAdd := int(elapsed.Seconds()) * tb.refillRate
    tb.tokens += tokensToAdd
    
    if tb.tokens > tb.capacity {
        tb.tokens = tb.capacity
    }
    
    tb.lastRefill = now
    
    // Check if we can consume a token
    if tb.tokens > 0 {
        tb.tokens--
        return true
    }
    
    return false
}

type RateLimiter struct {
    buckets map[string]*TokenBucket
    mutex   sync.RWMutex
    capacity int
    refillRate int
}

func NewRateLimiter(capacity, refillRate int) *RateLimiter {
    rl := &RateLimiter{
        buckets:    make(map[string]*TokenBucket),
        capacity:   capacity,
        refillRate: refillRate,
    }
    
    // Clean up old buckets periodically
    go rl.cleanup()
    
    return rl
}

func (rl *RateLimiter) Allow(key string) bool {
    rl.mutex.RLock()
    bucket, exists := rl.buckets[key]
    rl.mutex.RUnlock()
    
    if !exists {
        rl.mutex.Lock()
        // Double-check after acquiring write lock
        bucket, exists = rl.buckets[key]
        if !exists {
            bucket = NewTokenBucket(rl.capacity, rl.refillRate)
            rl.buckets[key] = bucket
        }
        rl.mutex.Unlock()
    }
    
    return bucket.Allow()
}

func (rl *RateLimiter) cleanup() {
    ticker := time.NewTicker(time.Hour)
    defer ticker.Stop()
    
    for range ticker.C {
        rl.mutex.Lock()
        now := time.Now()
        for key, bucket := range rl.buckets {
            bucket.mutex.Lock()
            if now.Sub(bucket.lastRefill) > time.Hour {
                delete(rl.buckets, key)
            }
            bucket.mutex.Unlock()
        }
        rl.mutex.Unlock()
    }
}

func RateLimit(capacity, refillRate int) func(http.Handler) http.Handler {
    limiter := NewRateLimiter(capacity, refillRate)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Get client IP
            ip, _, err := net.SplitHostPort(r.RemoteAddr)
            if err != nil {
                ip = r.RemoteAddr
            }
            
            // Check for forwarded IP headers
            if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
                ip = forwarded
            } else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
                ip = realIP
            }
            
            if !limiter.Allow(ip) {
                w.Header().Set("Retry-After", "60")
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

// Per-user rate limiting
func UserRateLimit(capacity, refillRate int) func(http.Handler) http.Handler {
    limiter := NewRateLimiter(capacity, refillRate)
    
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, ok := GetUserFromContext(r.Context())
            if !ok {
                // If no user context, fall back to IP-based limiting
                ip, _, _ := net.SplitHostPort(r.RemoteAddr)
                if !limiter.Allow(ip) {
                    http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                    return
                }
            } else {
                userKey := fmt.Sprintf("user:%d", claims.UserID)
                if !limiter.Allow(userKey) {
                    http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                    return
                }
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## SQL Injection Prevention

### Safe Database Operations

```go
// database/safe.go
package database

import (
    "database/sql"
    "fmt"
    "strings"
)

type SafeDB struct {
    db *sql.DB
}

func NewSafeDB(db *sql.DB) *SafeDB {
    return &SafeDB{db: db}
}

// Safe query with parameters
func (sdb *SafeDB) GetUser(userID int) (*User, error) {
    query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
    
    var user User
    err := sdb.db.QueryRow(query, userID).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.CreatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// Safe search with LIKE operator
func (sdb *SafeDB) SearchUsers(searchTerm string, limit, offset int) ([]User, error) {
    // Escape special characters for LIKE operator
    escapedTerm := strings.ReplaceAll(searchTerm, "%", "\\%")
    escapedTerm = strings.ReplaceAll(escapedTerm, "_", "\\_")
    
    query := `
        SELECT id, username, email, created_at 
        FROM users 
        WHERE username ILIKE $1 OR email ILIKE $1
        ORDER BY username
        LIMIT $2 OFFSET $3
    `
    
    rows, err := sdb.db.Query(query, "%"+escapedTerm+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, rows.Err()
}

// Safe dynamic query builder
type QueryBuilder struct {
    query  strings.Builder
    args   []interface{}
    argNum int
}

func NewQueryBuilder() *QueryBuilder {
    return &QueryBuilder{
        argNum: 1,
    }
}

func (qb *QueryBuilder) Select(columns ...string) *QueryBuilder {
    qb.query.WriteString("SELECT ")
    qb.query.WriteString(strings.Join(columns, ", "))
    return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
    qb.query.WriteString(" FROM ")
    qb.query.WriteString(table)
    return qb
}

func (qb *QueryBuilder) Where(condition string, args ...interface{}) *QueryBuilder {
    if qb.argNum == 1 {
        qb.query.WriteString(" WHERE ")
    } else {
        qb.query.WriteString(" AND ")
    }
    
    // Replace ? placeholders with numbered placeholders
    for range args {
        condition = strings.Replace(condition, "?", fmt.Sprintf("$%d", qb.argNum), 1)
        qb.argNum++
    }
    
    qb.query.WriteString(condition)
    qb.args = append(qb.args, args...)
    return qb
}

func (qb *QueryBuilder) OrderBy(column, direction string) *QueryBuilder {
    qb.query.WriteString(" ORDER BY ")
    qb.query.WriteString(column)
    qb.query.WriteString(" ")
    qb.query.WriteString(direction)
    return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
    qb.query.WriteString(fmt.Sprintf(" LIMIT $%d", qb.argNum))
    qb.args = append(qb.args, limit)
    qb.argNum++
    return qb
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
    return qb.query.String(), qb.args
}

// Example usage of QueryBuilder
func (sdb *SafeDB) GetUsersWithFilters(username, email string, limit int) ([]User, error) {
    qb := NewQueryBuilder()
    qb.Select("id", "username", "email", "created_at").From("users")
    
    if username != "" {
        qb.Where("username ILIKE ?", "%"+username+"%")
    }
    
    if email != "" {
        qb.Where("email ILIKE ?", "%"+email+"%")
    }
    
    qb.OrderBy("created_at", "DESC").Limit(limit)
    
    query, args := qb.Build()
    
    rows, err := sdb.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(
            &user.ID,
            &user.Username,
            &user.Email,
            &user.CreatedAt,
        )
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    
    return users, rows.Err()
}
```

## Secure Configuration

### Environment-based Secrets Management

```go
// config/secrets.go
package config

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "os"
    "strconv"
    "strings"
)

type SecretManager struct {
    secrets map[string]string
}

func NewSecretManager() *SecretManager {
    return &SecretManager{
        secrets: make(map[string]string),
    }
}

func (sm *SecretManager) LoadFromEnv() error {
    // Load secrets from environment variables
    for _, env := range os.Environ() {
        pair := strings.SplitN(env, "=", 2)
        if len(pair) == 2 {
            key := pair[0]
            value := pair[1]
            
            // Only store variables that look like secrets
            if strings.Contains(strings.ToLower(key), "secret") ||
               strings.Contains(strings.ToLower(key), "key") ||
               strings.Contains(strings.ToLower(key), "password") ||
               strings.Contains(strings.ToLower(key), "token") {
                sm.secrets[key] = value
            }
        }
    }
    
    return nil
}

func (sm *SecretManager) LoadFromFile(filename string) error {
    // Load secrets from file (Docker secrets, Kubernetes secrets, etc.)
    content, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    lines := strings.Split(string(content), "\n")
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        
        pair := strings.SplitN(line, "=", 2)
        if len(pair) == 2 {
            sm.secrets[pair[0]] = pair[1]
        }
    }
    
    return nil
}

func (sm *SecretManager) Get(key string) (string, bool) {
    value, exists := sm.secrets[key]
    return value, exists
}

func (sm *SecretManager) GetRequired(key string) (string, error) {
    value, exists := sm.secrets[key]
    if !exists {
        return "", fmt.Errorf("required secret %s not found", key)
    }
    return value, nil
}

func (sm *SecretManager) GetInt(key string, defaultValue int) int {
    value, exists := sm.secrets[key]
    if !exists {
        return defaultValue
    }
    
    intValue, err := strconv.Atoi(value)
    if err != nil {
        return defaultValue
    }
    
    return intValue
}

func (sm *SecretManager) GetBool(key string, defaultValue bool) bool {
    value, exists := sm.secrets[key]
    if !exists {
        return defaultValue
    }
    
    boolValue, err := strconv.ParseBool(value)
    if err != nil {
        return defaultValue
    }
    
    return boolValue
}

// Generate secure random secrets
func GenerateSecret(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

func GenerateJWTSecret() (string, error) {
    return GenerateSecret(32) // 256 bits
}

func GenerateAPIKey() (string, error) {
    return GenerateSecret(24) // 192 bits
}

// Validate configuration security
func ValidateSecurityConfig(config *Config) []string {
    var warnings []string
    
    // Check for weak JWT secret
    if len(config.JWTSecret) < 32 {
        warnings = append(warnings, "JWT secret should be at least 32 characters long")
    }
    
    // Check for default passwords
    defaultPasswords := []string{"password", "admin", "123456", "secret"}
    for _, defaultPwd := range defaultPasswords {
        if strings.Contains(strings.ToLower(config.Database.Password), defaultPwd) {
            warnings = append(warnings, "Database password appears to be a default/weak password")
            break
        }
    }
    
    // Check for HTTP in production
    if config.Environment == "production" && !config.Server.TLS.Enabled {
        warnings = append(warnings, "TLS should be enabled in production")
    }
    
    // Check for debug mode in production
    if config.Environment == "production" && config.Debug {
        warnings = append(warnings, "Debug mode should be disabled in production")
    }
    
    return warnings
}
```

## Exercises

### Exercise 1: Secure User Registration
Implement a secure user registration system that:
- Validates all input fields
- Hashes passwords securely
- Prevents username/email enumeration
- Implements rate limiting
- Sends verification emails

### Exercise 2: API Security
Secure a REST API with:
- JWT authentication
- Role-based authorization
- Input validation and sanitization
- Rate limiting per user and endpoint
- Security headers
- Request/response logging

### Exercise 3: Encryption Service
Build an encryption service that:
- Encrypts sensitive data at rest
- Implements key rotation
- Provides secure key derivation
- Handles multiple encryption algorithms
- Includes integrity verification

## Key Takeaways

- Always validate and sanitize user input
- Use parameterized queries to prevent SQL injection
- Implement proper authentication and authorization
- Hash passwords using strong algorithms (Argon2, bcrypt)
- Use HTTPS/TLS for all communications
- Implement rate limiting to prevent abuse
- Add security headers to protect against common attacks
- Keep secrets secure and rotate them regularly
- Log security events for monitoring and auditing
- Follow the principle of least privilege
- Keep dependencies updated and scan for vulnerabilities
- Implement proper error handling without information leakage

## Next Steps

Next, we'll explore [Performance Optimization](27-performance-optimization.md) to learn about optimizing Go applications for better performance and efficiency.