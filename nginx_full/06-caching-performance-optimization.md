# Caching and Performance Optimization - Part 6

## Understanding Nginx Caching

Nginx provides several types of caching mechanisms:

- **Proxy Cache**: Cache responses from upstream servers
- **FastCGI Cache**: Cache responses from FastCGI applications (PHP)
- **Browser Cache**: Control client-side caching with headers
- **Static File Caching**: Optimize static content delivery

## Proxy Caching

### Basic Proxy Cache Setup

```nginx
http {
    # Define cache path and settings
    proxy_cache_path /var/cache/nginx/proxy 
                     levels=1:2 
                     keys_zone=my_cache:10m 
                     max_size=1g 
                     inactive=60m 
                     use_temp_path=off;
    
    server {
        listen 80;
        server_name example.com;
        
        location / {
            proxy_cache my_cache;
            proxy_cache_valid 200 302 10m;
            proxy_cache_valid 404 1m;
            proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
            
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            
            # Add cache status header for debugging
            add_header X-Cache-Status $upstream_cache_status;
        }
    }
}
```

### Cache Path Parameters Explained

- `levels=1:2`: Directory structure (1 char/2 char subdirectories)
- `keys_zone=name:size`: Shared memory zone for cache keys
- `max_size=size`: Maximum cache size
- `inactive=time`: Remove cached files not accessed within this time
- `use_temp_path=off`: Write directly to cache path (faster)

### Advanced Proxy Caching

```nginx
http {
    # Multiple cache zones for different content types
    proxy_cache_path /var/cache/nginx/static 
                     levels=1:2 
                     keys_zone=static_cache:10m 
                     max_size=1g 
                     inactive=1d;
    
    proxy_cache_path /var/cache/nginx/api 
                     levels=1:2 
                     keys_zone=api_cache:10m 
                     max_size=500m 
                     inactive=1h;
    
    # Custom cache key
    proxy_cache_key "$scheme$request_method$host$request_uri$http_accept_encoding";
    
    server {
        listen 80;
        server_name example.com;
        
        # Cache static API responses
        location /api/static/ {
            proxy_cache api_cache;
            proxy_cache_valid 200 1h;
            proxy_cache_valid 404 5m;
            proxy_cache_lock on;
            proxy_cache_lock_timeout 5s;
            proxy_cache_lock_age 1h;
            
            proxy_pass http://api_backend;
        }
        
        # Cache images and static files
        location ~* \.(jpg|jpeg|png|gif|css|js)$ {
            proxy_cache static_cache;
            proxy_cache_valid 200 1d;
            proxy_cache_valid 404 1h;
            
            proxy_pass http://static_backend;
        }
        
        # Don't cache dynamic content
        location /api/dynamic/ {
            proxy_cache off;
            proxy_pass http://api_backend;
        }
    }
}
```

### Cache Bypass and Purging

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Set cache bypass conditions
    set $skip_cache 0;
    
    # Skip cache for POST requests
    if ($request_method = POST) {
        set $skip_cache 1;
    }
    
    # Skip cache for URLs with query parameters
    if ($query_string != "") {
        set $skip_cache 1;
    }
    
    # Skip cache for admin pages
    if ($request_uri ~* "/admin/|/wp-admin/") {
        set $skip_cache 1;
    }
    
    # Skip cache for logged-in users (based on cookie)
    if ($http_cookie ~* "logged_in") {
        set $skip_cache 1;
    }
    
    location / {
        proxy_cache my_cache;
        proxy_cache_bypass $skip_cache;
        proxy_no_cache $skip_cache;
        
        proxy_pass http://backend;
    }
    
    # Cache purge endpoint (requires nginx-plus or third-party module)
    location ~ /purge(/.*) {
        allow 127.0.0.1;
        allow 192.168.1.0/24;
        deny all;
        
        proxy_cache_purge my_cache "$scheme$request_method$host$1";
    }
}
```

## FastCGI Caching (PHP)

### Basic FastCGI Cache

```nginx
http {
    # FastCGI cache path
    fastcgi_cache_path /var/cache/nginx/fastcgi 
                       levels=1:2 
                       keys_zone=php_cache:10m 
                       max_size=1g 
                       inactive=1h;
    
    server {
        listen 80;
        server_name example.com;
        root /var/www/html;
        index index.php index.html;
        
        # Cache bypass conditions
        set $skip_cache 0;
        
        if ($request_method = POST) {
            set $skip_cache 1;
        }
        
        if ($query_string != "") {
            set $skip_cache 1;
        }
        
        if ($request_uri ~* "/wp-admin/|/xmlrpc.php|wp-.*.php|/feed/|index.php|sitemap(_index)?.xml") {
            set $skip_cache 1;
        }
        
        if ($http_cookie ~* "comment_author|wordpress_[a-f0-9]+|wp-postpass|wordpress_no_cache|wordpress_logged_in") {
            set $skip_cache 1;
        }
        
        location ~ \.php$ {
            fastcgi_cache php_cache;
            fastcgi_cache_valid 200 1h;
            fastcgi_cache_valid 404 1m;
            fastcgi_cache_bypass $skip_cache;
            fastcgi_no_cache $skip_cache;
            
            fastcgi_pass unix:/var/run/php/php-fpm.sock;
            fastcgi_index index.php;
            fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
            include fastcgi_params;
            
            # Add cache status header
            add_header X-Cache-Status $upstream_cache_status;
        }
        
        location / {
            try_files $uri $uri/ /index.php?$args;
        }
    }
}
```

### WordPress-Specific FastCGI Cache

```nginx
server {
    listen 80;
    server_name wordpress.com;
    root /var/www/wordpress;
    index index.php;
    
    # Cache configuration
    set $skip_cache 0;
    
    # POST requests and URLs with a query string should always go to PHP
    if ($request_method = POST) {
        set $skip_cache 1;
    }
    
    if ($query_string != "") {
        set $skip_cache 1;
    }
    
    # Don't cache URIs containing the following segments
    if ($request_uri ~* "/wp-admin/|/xmlrpc.php|wp-.*.php|/feed/|index.php|sitemap(_index)?.xml") {
        set $skip_cache 1;
    }
    
    # Don't use the cache for logged-in users or recent commenters
    if ($http_cookie ~* "comment_author|wordpress_[a-f0-9]+|wp-postpass|wordpress_no_cache|wordpress_logged_in") {
        set $skip_cache 1;
    }
    
    location / {
        try_files $uri $uri/ /index.php?$args;
    }
    
    location ~ \.php$ {
        try_files $uri =404;
        fastcgi_split_path_info ^(.+\.php)(/.+)$;
        
        fastcgi_cache php_cache;
        fastcgi_cache_valid 200 1h;
        fastcgi_cache_valid 404 1m;
        fastcgi_cache_bypass $skip_cache;
        fastcgi_no_cache $skip_cache;
        fastcgi_cache_lock on;
        
        fastcgi_pass unix:/var/run/php/php-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
        
        add_header X-Cache-Status $upstream_cache_status;
    }
    
    # Cache purge endpoint
    location ~ /purge(/.*) {
        fastcgi_cache_purge php_cache "$scheme$request_method$host$1";
    }
}
```

## Browser Caching

### Static File Caching Headers

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/html;
    
    # Cache static files for 1 year
    location ~* \.(jpg|jpeg|png|gif|ico|css|js|pdf|txt)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
        add_header Vary Accept-Encoding;
        access_log off;
    }
    
    # Cache fonts
    location ~* \.(woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public";
        add_header Access-Control-Allow-Origin "*";
    }
    
    # Cache HTML files for shorter period
    location ~* \.(html|htm)$ {
        expires 1h;
        add_header Cache-Control "public";
    }
    
    # Don't cache dynamic content
    location ~* \.(php|pl|py)$ {
        add_header Cache-Control "no-cache, no-store, must-revalidate";
        add_header Pragma "no-cache";
        add_header Expires "0";
    }
}
```

### Advanced Cache Control

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Map file extensions to cache times
    map $sent_http_content_type $expires {
        default                    off;
        text/html                  1h;
        text/css                   1y;
        application/javascript     1y;
        ~image/                    1M;
        ~font/                     1y;
        application/pdf            1M;
    }
    
    expires $expires;
    
    # Conditional caching based on user agent
    map $http_user_agent $cache_control {
        default "public";
        ~*bot   "no-cache";
    }
    
    add_header Cache-Control $cache_control;
    
    location / {
        root /var/www/html;
        try_files $uri $uri/ =404;
    }
}
```

## Performance Optimization

### Gzip Compression

```nginx
http {
    # Enable gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json
        application/xml
        image/svg+xml;
    
    # Disable gzip for IE6
    gzip_disable "msie6";
}
```

### Brotli Compression (if module available)

```nginx
http {
    # Enable Brotli compression
    brotli on;
    brotli_comp_level 6;
    brotli_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/json
        application/xml
        application/rss+xml
        application/atom+xml
        image/svg+xml;
}
```

### Connection Optimization

```nginx
http {
    # Optimize connections
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    keepalive_requests 100;
    
    # Client settings
    client_max_body_size 16M;
    client_body_buffer_size 128k;
    client_header_buffer_size 1k;
    large_client_header_buffers 4 4k;
    
    # Timeouts
    client_body_timeout 12;
    client_header_timeout 12;
    send_timeout 10;
    
    # File descriptor cache
    open_file_cache max=1000 inactive=20s;
    open_file_cache_valid 30s;
    open_file_cache_min_uses 2;
    open_file_cache_errors on;
}
```

### Worker Process Optimization

```nginx
# Main context
user nginx;
worker_processes auto;  # or specific number like 4
worker_rlimit_nofile 65535;

events {
    worker_connections 1024;
    use epoll;  # Linux
    multi_accept on;
}

http {
    # Worker process settings
    worker_cpu_affinity auto;
    
    # Connection pooling for upstream
    upstream backend {
        server backend1.example.com;
        server backend2.example.com;
        keepalive 32;
    }
    
    server {
        location / {
            proxy_pass http://backend;
            proxy_http_version 1.1;
            proxy_set_header Connection "";
        }
    }
}
```

## CDN Integration

### CloudFlare Integration

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Real IP from CloudFlare
    set_real_ip_from 103.21.244.0/22;
    set_real_ip_from 103.22.200.0/22;
    set_real_ip_from 103.31.4.0/22;
    set_real_ip_from 104.16.0.0/12;
    set_real_ip_from 108.162.192.0/18;
    set_real_ip_from 131.0.72.0/22;
    set_real_ip_from 141.101.64.0/18;
    set_real_ip_from 162.158.0.0/15;
    set_real_ip_from 172.64.0.0/13;
    set_real_ip_from 173.245.48.0/20;
    set_real_ip_from 188.114.96.0/20;
    set_real_ip_from 190.93.240.0/20;
    set_real_ip_from 197.234.240.0/22;
    set_real_ip_from 198.41.128.0/17;
    real_ip_header CF-Connecting-IP;
    
    # Cache static files locally
    location ~* \.(jpg|jpeg|png|gif|css|js)$ {
        expires 1y;
        add_header Cache-Control "public";
        
        # Serve from local cache if available
        try_files $uri @cdn;
    }
    
    location @cdn {
        # Fallback to CDN
        return 301 https://cdn.example.com$request_uri;
    }
}
```

### Origin Server Configuration

```nginx
server {
    listen 80;
    server_name origin.example.com;
    
    # Only allow CDN IPs
    allow 103.21.244.0/22;
    allow 103.22.200.0/22;
    # ... other CloudFlare IPs
    deny all;
    
    # Optimize for CDN
    location ~* \.(jpg|jpeg|png|gif|css|js)$ {
        expires 1y;
        add_header Cache-Control "public, max-age=31536000";
        add_header Vary "Accept-Encoding";
        
        # Enable CORS for CDN
        add_header Access-Control-Allow-Origin "https://example.com";
    }
    
    location / {
        root /var/www/html;
    }
}
```

## Monitoring and Analytics

### Cache Hit Rate Monitoring

```nginx
http {
    # Log format with cache status
    log_format cache_log '$remote_addr - $remote_user [$time_local] '
                        '"$request" $status $body_bytes_sent '
                        '"$http_referer" "$http_user_agent" '
                        'cache_status:$upstream_cache_status '
                        'response_time:$request_time';
    
    server {
        access_log /var/log/nginx/cache.log cache_log;
        
        location / {
            proxy_cache my_cache;
            proxy_pass http://backend;
            
            # Add cache headers for monitoring
            add_header X-Cache-Status $upstream_cache_status;
            add_header X-Cache-Key $scheme$request_method$host$request_uri;
        }
    }
}
```

### Performance Monitoring

```nginx
http {
    # Detailed performance logging
    log_format performance '$remote_addr - $remote_user [$time_local] '
                          '"$request" $status $body_bytes_sent '
                          'request_time:$request_time '
                          'upstream_time:$upstream_response_time '
                          'upstream_connect_time:$upstream_connect_time '
                          'upstream_header_time:$upstream_header_time';
    
    server {
        access_log /var/log/nginx/performance.log performance;
    }
}
```

## Cache Management Scripts

### Cache Warming Script

```bash
#!/bin/bash
# cache-warm.sh

SITE_URL="https://example.com"
SITEMAP_URL="$SITE_URL/sitemap.xml"

# Extract URLs from sitemap
curl -s $SITEMAP_URL | grep -oP '(?<=<loc>)[^<]+' | while read url; do
    echo "Warming cache for: $url"
    curl -s -o /dev/null "$url"
    sleep 0.1  # Be nice to the server
done

echo "Cache warming completed"
```

### Cache Statistics Script

```bash
#!/bin/bash
# cache-stats.sh

CACHE_DIR="/var/cache/nginx"

echo "=== Nginx Cache Statistics ==="
echo "Cache directory: $CACHE_DIR"
echo "Total cache size: $(du -sh $CACHE_DIR | cut -f1)"
echo "Number of cached files: $(find $CACHE_DIR -type f | wc -l)"
echo "Cache directories: $(find $CACHE_DIR -type d | wc -l)"

echo "
=== Cache Hit Rates (last 1000 requests) ==="
tail -1000 /var/log/nginx/access.log | grep -o 'cache_status:[A-Z]*' | sort | uniq -c | sort -nr

echo "
=== Top Cached URLs ==="
tail -1000 /var/log/nginx/access.log | grep 'cache_status:HIT' | awk '{print $7}' | sort | uniq -c | sort -nr | head -10
```

### Cache Cleanup Script

```bash
#!/bin/bash
# cache-cleanup.sh

CACHE_DIR="/var/cache/nginx"
MAX_SIZE="1G"

echo "Cleaning up cache directory: $CACHE_DIR"

# Remove files older than 7 days
find $CACHE_DIR -type f -mtime +7 -delete

# Check if cache size exceeds limit
CURRENT_SIZE=$(du -sb $CACHE_DIR | cut -f1)
MAX_SIZE_BYTES=$(echo $MAX_SIZE | sed 's/G/*1024*1024*1024/' | bc)

if [ $CURRENT_SIZE -gt $MAX_SIZE_BYTES ]; then
    echo "Cache size exceeds limit, removing oldest files"
    find $CACHE_DIR -type f -printf '%T@ %p\n' | sort -n | head -n 100 | cut -d' ' -f2- | xargs rm -f
fi

echo "Cache cleanup completed"
echo "Current cache size: $(du -sh $CACHE_DIR | cut -f1)"
```

## Best Practices

1. **Monitor cache hit rates** regularly
2. **Set appropriate cache times** based on content type
3. **Use cache keys wisely** to avoid cache pollution
4. **Implement cache warming** for critical content
5. **Configure proper cache invalidation**
6. **Monitor disk space** for cache directories
7. **Use compression** to reduce bandwidth
8. **Optimize worker processes** for your hardware
9. **Enable HTTP/2** for better performance
10. **Use CDN** for global content delivery

---

**Continue to Part 7: Monitoring, Logging, and Troubleshooting**