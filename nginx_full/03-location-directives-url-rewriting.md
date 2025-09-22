# Location Directives and URL Rewriting - Part 3

## Advanced Location Directives

### Location Matching Syntax

```nginx
location [modifier] pattern {
    # Configuration
}
```

### Modifiers Explained

| Modifier | Description | Example |
|----------|-------------|----------|
| `=` | Exact match | `location = /login` |
| `^~` | Priority prefix | `location ^~ /images/` |
| `~` | Case-sensitive regex | `location ~ \.php$` |
| `~*` | Case-insensitive regex | `location ~* \.(jpg\|png)$` |
| (none) | Prefix match | `location /api/` |

### Location Processing Order

1. **Exact matches** (`=`) are checked first
2. **Priority prefix matches** (`^~`) are checked next
3. **Regular expressions** (`~` and `~*`) in order of appearance
4. **Prefix matches** (longest match wins)

### Practical Location Examples

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/html;
    
    # Exact match for homepage
    location = / {
        try_files /index.html =404;
    }
    
    # Exact match for specific file
    location = /favicon.ico {
        log_not_found off;
        access_log off;
    }
    
    # Priority prefix for admin area
    location ^~ /admin/ {
        auth_basic "Admin Area";
        auth_basic_user_file /etc/nginx/.htpasswd;
        try_files $uri $uri/ /admin/index.php?$query_string;
    }
    
    # Case-sensitive regex for PHP files
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
    
    # Case-insensitive regex for static files
    location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        add_header Vary Accept-Encoding;
    }
    
    # Prefix match for API
    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    # Default location
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
}
```

## URL Rewriting with Rewrite Module

### Basic Rewrite Syntax

```nginx
rewrite regex replacement [flag];
```

### Rewrite Flags

| Flag | Description |
|------|-------------|
| `last` | Stop processing, start new location lookup |
| `break` | Stop processing in current location |
| `redirect` | Return 302 temporary redirect |
| `permanent` | Return 301 permanent redirect |

### Common Rewrite Examples

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/html;
    
    # Redirect www to non-www
    if ($host = www.example.com) {
        return 301 http://example.com$request_uri;
    }
    
    # Redirect old URLs to new structure
    rewrite ^/old-page$ /new-page permanent;
    rewrite ^/blog/([0-9]+)/([^/]+)$ /blog/$1/$2.html last;
    
    # Remove trailing slash
    rewrite ^/(.*)/$ /$1 permanent;
    
    # Add .html extension
    rewrite ^/([^.]+)$ /$1.html last;
    
    # API versioning
    rewrite ^/api/v1/(.*)$ /api/v1.php?request=$1 last;
    
    location / {
        try_files $uri $uri/ =404;
    }
}
```

## Try Files Directive

### Basic Syntax
```nginx
try_files file1 file2 ... fallback;
```

### Common Patterns

```nginx
# Basic try_files
location / {
    try_files $uri $uri/ =404;
}

# With index fallback
location / {
    try_files $uri $uri/ /index.html;
}

# For single-page applications
location / {
    try_files $uri $uri/ /index.html;
}

# For PHP applications
location / {
    try_files $uri $uri/ /index.php?$query_string;
}

# With multiple fallbacks
location / {
    try_files $uri $uri/ @fallback;
}

location @fallback {
    proxy_pass http://backend;
}
```

## Advanced URL Manipulation

### Using Named Captures

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Named capture groups
    location ~ ^/user/(?<username>[a-zA-Z0-9_]+)/profile$ {
        try_files /profiles/$username.html =404;
    }
    
    # Using captures in proxy_pass
    location ~ ^/api/v(?<version>[0-9]+)/(?<endpoint>.*)$ {
        proxy_pass http://api-v$version.internal/$endpoint;
    }
}
```

### Complex Rewrite Rules

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/html;
    
    # WordPress-style permalinks
    location / {
        try_files $uri $uri/ /index.php?$args;
    }
    
    # Custom post type rewrite
    rewrite ^/products/([^/]+)/?$ /product.php?slug=$1 last;
    rewrite ^/category/([^/]+)/?$ /category.php?name=$1 last;
    
    # Date-based archives
    rewrite ^/([0-9]{4})/([0-9]{2})/([0-9]{2})/?$ /archive.php?year=$1&month=$2&day=$3 last;
    rewrite ^/([0-9]{4})/([0-9]{2})/?$ /archive.php?year=$1&month=$2 last;
    rewrite ^/([0-9]{4})/?$ /archive.php?year=$1 last;
}
```

## Conditional Logic

### If Directive (Use Sparingly)

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Redirect mobile users
    if ($http_user_agent ~* "(iPhone|Android)") {
        return 301 http://m.example.com$request_uri;
    }
    
    # Block specific user agents
    if ($http_user_agent ~* "(bot|crawler|spider)") {
        return 403;
    }
    
    # Maintenance mode
    if (-f $document_root/maintenance.html) {
        return 503;
    }
    
    # Force HTTPS
    if ($scheme != "https") {
        return 301 https://$server_name$request_uri;
    }
}
```

### Better Alternatives to If

```nginx
# Instead of if for static files, use location
location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public";
}

# Instead of if for redirects, use return in server block
server {
    listen 80;
    server_name www.example.com;
    return 301 http://example.com$request_uri;
}

# Use map for complex conditions
map $http_user_agent $mobile {
    default 0;
    ~*iPhone 1;
    ~*Android 1;
}

server {
    if ($mobile) {
        return 301 http://m.example.com$request_uri;
    }
}
```

## Real-World Examples

### Single Page Application (SPA)

```nginx
server {
    listen 80;
    server_name myapp.com;
    root /var/www/myapp/dist;
    index index.html;
    
    # Handle client-side routing
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    # API proxy
    location /api/ {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # Static assets with caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

### WordPress Configuration

```nginx
server {
    listen 80;
    server_name wordpress.com;
    root /var/www/wordpress;
    index index.php index.html;
    
    # WordPress permalinks
    location / {
        try_files $uri $uri/ /index.php?$args;
    }
    
    # PHP processing
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
    
    # Deny access to sensitive files
    location ~* /(wp-config\.php|readme\.html|license\.txt)$ {
        deny all;
    }
    
    # Static files
    location ~* \.(css|js|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public";
    }
}
```

### E-commerce Site

```nginx
server {
    listen 80;
    server_name shop.example.com;
    root /var/www/shop;
    
    # Product pages
    location ~ ^/product/([a-zA-Z0-9-]+)$ {
        try_files /products/$1.html /product.php?slug=$1;
    }
    
    # Category pages
    location ~ ^/category/([a-zA-Z0-9-]+)$ {
        try_files /categories/$1.html /category.php?slug=$1;
    }
    
    # Search functionality
    location /search {
        try_files $uri /search.php?$query_string;
    }
    
    # Checkout process (force HTTPS)
    location ^~ /checkout {
        return 301 https://$server_name$request_uri;
    }
    
    # Default handler
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
}
```

## Debugging Location Matches

### Add Debug Headers

```nginx
location /debug {
    add_header X-Location "debug location";
    add_header X-URI "$uri";
    add_header X-Args "$args";
    return 200 "Debug info in headers";
}
```

### Test Location Matching

```bash
# Test with curl
curl -I http://example.com/test-path

# Check error logs
sudo tail -f /var/log/nginx/error.log

# Use nginx debug mode (if compiled with debug)
nginx -V 2>&1 | grep -o with-debug
```

---

**Continue to Part 4: Reverse Proxy and Load Balancing**