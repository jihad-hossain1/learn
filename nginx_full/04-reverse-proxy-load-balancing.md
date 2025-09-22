# Reverse Proxy and Load Balancing - Part 4

## Understanding Reverse Proxy

A reverse proxy sits between clients and backend servers, forwarding client requests to backend servers and returning responses back to clients. Unlike a forward proxy that acts on behalf of clients, a reverse proxy acts on behalf of servers.

### Benefits of Reverse Proxy

- **Load Distribution**: Distribute requests across multiple backend servers
- **SSL Termination**: Handle SSL encryption/decryption
- **Caching**: Cache responses to improve performance
- **Security**: Hide backend server details from clients
- **Compression**: Compress responses before sending to clients
- **Rate Limiting**: Control request rates

## Basic Reverse Proxy Configuration

### Simple Proxy Pass

```nginx
server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend-server:8080;
    }
}
```

### Essential Proxy Headers

```nginx
server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend-server:8080;
        
        # Essential headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $host;
        proxy_set_header X-Forwarded-Port $server_port;
    }
}
```

### Complete Proxy Configuration

```nginx
server {
    listen 80;
    server_name api.example.com;
    
    location / {
        proxy_pass http://backend;
        
        # Headers
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
        
        # Buffering
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;
        
        # Other settings
        proxy_redirect off;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }
}
```

## Load Balancing

### Upstream Blocks

```nginx
upstream backend {
    server 192.168.1.10:8080;
    server 192.168.1.11:8080;
    server 192.168.1.12:8080;
}

server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend;
    }
}
```

### Load Balancing Methods

#### Round Robin (Default)
```nginx
upstream backend {
    server server1.example.com;
    server server2.example.com;
    server server3.example.com;
}
```

#### Least Connections
```nginx
upstream backend {
    least_conn;
    server server1.example.com;
    server server2.example.com;
    server server3.example.com;
}
```

#### IP Hash (Session Persistence)
```nginx
upstream backend {
    ip_hash;
    server server1.example.com;
    server server2.example.com;
    server server3.example.com;
}
```

#### Weighted Round Robin
```nginx
upstream backend {
    server server1.example.com weight=3;
    server server2.example.com weight=2;
    server server3.example.com weight=1;
}
```

#### Hash (Custom Key)
```nginx
upstream backend {
    hash $request_uri consistent;
    server server1.example.com;
    server server2.example.com;
    server server3.example.com;
}
```

### Server Parameters

```nginx
upstream backend {
    server server1.example.com:8080 weight=5;
    server server2.example.com:8080 max_fails=3 fail_timeout=30s;
    server server3.example.com:8080 backup;
    server server4.example.com:8080 down;
    server unix:/tmp/backend3 weight=1;
}
```

#### Parameter Explanations

- `weight=n` - Set weight for server (default: 1)
- `max_fails=n` - Maximum failed attempts (default: 1)
- `fail_timeout=time` - Time to consider server failed (default: 10s)
- `backup` - Mark server as backup (only used when primary servers fail)
- `down` - Mark server as permanently unavailable
- `max_conns=n` - Limit maximum connections to server

## Advanced Proxy Configurations

### Microservices Architecture

```nginx
# User service
upstream user-service {
    server user1.internal:3001;
    server user2.internal:3001;
}

# Product service
upstream product-service {
    server product1.internal:3002;
    server product2.internal:3002;
}

# Order service
upstream order-service {
    server order1.internal:3003;
    server order2.internal:3003;
}

server {
    listen 80;
    server_name api.example.com;
    
    # User API
    location /api/users/ {
        proxy_pass http://user-service/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # Product API
    location /api/products/ {
        proxy_pass http://product-service/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
    
    # Order API
    location /api/orders/ {
        proxy_pass http://order-service/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### WebSocket Proxying

```nginx
upstream websocket {
    server websocket1.example.com:8080;
    server websocket2.example.com:8080;
}

server {
    listen 80;
    server_name ws.example.com;
    
    location / {
        proxy_pass http://websocket;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket specific timeouts
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }
}
```

### Health Checks (Nginx Plus)

```nginx
upstream backend {
    zone backend 64k;
    server server1.example.com:8080;
    server server2.example.com:8080;
    server server3.example.com:8080;
}

server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend;
        health_check;
    }
    
    # Health check endpoint
    location /health {
        health_check_status;
        access_log off;
    }
}
```

### Custom Health Checks (Open Source)

```nginx
upstream backend {
    server server1.example.com:8080;
    server server2.example.com:8080;
    server server3.example.com:8080;
}

server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend;
        proxy_next_upstream error timeout invalid_header http_500 http_502 http_503 http_504;
        proxy_next_upstream_tries 3;
        proxy_next_upstream_timeout 30s;
    }
}
```

## SSL/TLS Termination

### Basic SSL Termination

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto https;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name example.com;
    return 301 https://$server_name$request_uri;
}
```

### Advanced SSL Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # SSL certificates
    ssl_certificate /etc/ssl/certs/example.com.crt;
    ssl_certificate_key /etc/ssl/private/example.com.key;
    
    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Caching with Proxy

### Basic Proxy Caching

```nginx
http {
    proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=my_cache:10m max_size=10g 
                     inactive=60m use_temp_path=off;
    
    server {
        listen 80;
        server_name example.com;
        
        location / {
            proxy_cache my_cache;
            proxy_cache_valid 200 302 10m;
            proxy_cache_valid 404 1m;
            proxy_cache_use_stale error timeout updating http_500 http_502 http_503 http_504;
            proxy_cache_lock on;
            
            proxy_pass http://backend;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            
            # Add cache status header
            add_header X-Cache-Status $upstream_cache_status;
        }
    }
}
```

### Selective Caching

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Cache static API responses
    location /api/static/ {
        proxy_cache my_cache;
        proxy_cache_valid 200 1h;
        proxy_pass http://backend;
    }
    
    # Don't cache dynamic content
    location /api/dynamic/ {
        proxy_cache off;
        proxy_pass http://backend;
    }
    
    # Cache with custom key
    location /api/user/ {
        proxy_cache my_cache;
        proxy_cache_key "$scheme$request_method$host$request_uri$http_authorization";
        proxy_cache_valid 200 5m;
        proxy_pass http://backend;
    }
}
```

## Monitoring and Debugging

### Status Module

```nginx
server {
    listen 8080;
    server_name localhost;
    
    location /nginx_status {
        stub_status on;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
    
    location /upstream_status {
        upstream_status;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
}
```

### Logging

```nginx
http {
    log_format upstream_log '$remote_addr - $remote_user [$time_local] '
                           '"$request" $status $body_bytes_sent '
                           '"$http_referer" "$http_user_agent" '
                           'rt=$request_time uct="$upstream_connect_time" '
                           'uht="$upstream_header_time" urt="$upstream_response_time"';
    
    server {
        access_log /var/log/nginx/upstream.log upstream_log;
        
        location / {
            proxy_pass http://backend;
        }
    }
}
```

## Best Practices

1. **Always set proper headers** for backend applications
2. **Use health checks** to ensure backend availability
3. **Implement proper timeouts** to avoid hanging connections
4. **Monitor upstream performance** with logging
5. **Use SSL termination** at the proxy level
6. **Implement caching** where appropriate
7. **Configure proper buffer sizes** for your use case
8. **Use connection pooling** with `keepalive`

### Connection Pooling

```nginx
upstream backend {
    server server1.example.com:8080;
    server server2.example.com:8080;
    keepalive 32;
}

server {
    location / {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }
}
```

---

**Continue to Part 5: SSL/TLS Configuration and Security**