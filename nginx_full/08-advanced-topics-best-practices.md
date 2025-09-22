# Advanced Topics and Best Practices - Part 8

## Advanced Nginx Modules

### Lua Module (OpenResty)

OpenResty extends Nginx with Lua scripting capabilities, enabling dynamic content generation and complex logic.

#### Installation
```bash
# Ubuntu/Debian
sudo apt install openresty

# Or compile from source
wget https://openresty.org/download/openresty-1.21.4.1.tar.gz
tar -xzf openresty-1.21.4.1.tar.gz
cd openresty-1.21.4.1
./configure --with-luajit
make && sudo make install
```

#### Basic Lua Usage

```nginx
server {
    listen 80;
    server_name example.com;
    
    # Simple Lua content
    location /hello {
        content_by_lua_block {
            ngx.say("Hello from Lua!")
            ngx.say("Current time: ", os.date())
        }
    }
    
    # Dynamic routing based on time
    location /dynamic {
        access_by_lua_block {
            local hour = tonumber(os.date("%H"))
            if hour < 9 or hour > 17 then
                ngx.status = 503
                ngx.say("Service unavailable outside business hours")
                ngx.exit(503)
            end
        }
        
        proxy_pass http://backend;
    }
    
    # Custom authentication
    location /api/ {
        access_by_lua_block {
            local token = ngx.var.http_authorization
            if not token or token ~= "Bearer secret123" then
                ngx.status = 401
                ngx.header["WWW-Authenticate"] = "Bearer"
                ngx.say("Unauthorized")
                ngx.exit(401)
            end
        }
        
        proxy_pass http://api_backend;
    }
    
    # Rate limiting with Redis
    location /limited {
        access_by_lua_block {
            local redis = require "resty.redis"
            local red = redis:new()
            red:set_timeout(1000)
            
            local ok, err = red:connect("127.0.0.1", 6379)
            if not ok then
                ngx.log(ngx.ERR, "Failed to connect to Redis: ", err)
                return
            end
            
            local key = "rate_limit:" .. ngx.var.remote_addr
            local current = red:incr(key)
            
            if current == 1 then
                red:expire(key, 60)  -- 1 minute window
            end
            
            if current > 10 then  -- 10 requests per minute
                ngx.status = 429
                ngx.say("Rate limit exceeded")
                ngx.exit(429)
            end
        }
        
        proxy_pass http://backend;
    }
}
```

### Stream Module (TCP/UDP Load Balancing)

```nginx
# /etc/nginx/nginx.conf
stream {
    # TCP load balancing for database
    upstream database {
        server db1.example.com:5432;
        server db2.example.com:5432;
        server db3.example.com:5432;
    }
    
    server {
        listen 5432;
        proxy_pass database;
        proxy_timeout 1s;
        proxy_responses 1;
    }
    
    # UDP load balancing for DNS
    upstream dns_servers {
        server 8.8.8.8:53;
        server 8.8.4.4:53;
    }
    
    server {
        listen 53 udp;
        proxy_pass dns_servers;
        proxy_timeout 1s;
        proxy_responses 1;
    }
    
    # SSL termination for TCP
    server {
        listen 443 ssl;
        ssl_certificate /path/to/cert.pem;
        ssl_certificate_key /path/to/key.pem;
        
        proxy_pass backend_ssl;
    }
}
```

### Image Processing Module

```nginx
# Requires nginx-module-image-filter
server {
    listen 80;
    server_name images.example.com;
    
    location ~* ^/resize/(\d+)x(\d+)/(.+)$ {
        set $width $1;
        set $height $2;
        set $image_path $3;
        
        # Security check
        if ($width !~ "^[0-9]+$") {
            return 400;
        }
        if ($height !~ "^[0-9]+$") {
            return 400;
        }
        if ($width > 2000) {
            return 400;
        }
        if ($height > 2000) {
            return 400;
        }
        
        image_filter resize $width $height;
        image_filter_jpeg_quality 85;
        image_filter_buffer 10M;
        
        try_files /$image_path =404;
    }
    
    location ~* ^/crop/(\d+)x(\d+)/(.+)$ {
        set $width $1;
        set $height $2;
        set $image_path $3;
        
        image_filter crop $width $height;
        image_filter_jpeg_quality 85;
        
        try_files /$image_path =404;
    }
    
    location / {
        root /var/www/images;
        expires 1y;
    }
}
```

## High Availability and Clustering

### Nginx with Keepalived

#### Master Configuration
```bash
# /etc/keepalived/keepalived.conf (Master)
vrrp_script chk_nginx {
    script "/usr/bin/curl -f http://localhost/health || exit 1"
    interval 2
    weight -2
    fall 3
    rise 2
}

vrrp_instance VI_1 {
    state MASTER
    interface eth0
    virtual_router_id 51
    priority 110
    advert_int 1
    authentication {
        auth_type PASS
        auth_pass mypassword
    }
    virtual_ipaddress {
        192.168.1.100
    }
    track_script {
        chk_nginx
    }
}
```

#### Backup Configuration
```bash
# /etc/keepalived/keepalived.conf (Backup)
vrrp_script chk_nginx {
    script "/usr/bin/curl -f http://localhost/health || exit 1"
    interval 2
    weight -2
    fall 3
    rise 2
}

vrrp_instance VI_1 {
    state BACKUP
    interface eth0
    virtual_router_id 51
    priority 100
    advert_int 1
    authentication {
        auth_type PASS
        auth_pass mypassword
    }
    virtual_ipaddress {
        192.168.1.100
    }
    track_script {
        chk_nginx
    }
}
```

### Docker Swarm with Nginx

```yaml
# docker-compose.yml
version: '3.8'
services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl
    deploy:
      replicas: 2
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
    networks:
      - frontend
      - backend
  
  app:
    image: myapp:latest
    deploy:
      replicas: 3
    networks:
      - backend

networks:
  frontend:
    external: true
  backend:
    driver: overlay
```

### Kubernetes Ingress

```yaml
# nginx-ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    nginx.ingress.kubernetes.io/rate-limit: "100"
    nginx.ingress.kubernetes.io/rate-limit-window: "1m"
spec:
  tls:
  - hosts:
    - example.com
    secretName: example-tls
  rules:
  - host: example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 80
      - path: /
        pathType: Prefix
        backend:
          service:
            name: web-service
            port:
              number: 80
```

## Security Hardening

### Complete Security Configuration

```nginx
# /etc/nginx/conf.d/security.conf

# Hide nginx version and server tokens
server_tokens off;
more_set_headers "Server: WebServer";

# Security headers
more_set_headers "X-Frame-Options: DENY";
more_set_headers "X-XSS-Protection: 1; mode=block";
more_set_headers "X-Content-Type-Options: nosniff";
more_set_headers "Referrer-Policy: strict-origin-when-cross-origin";
more_set_headers "X-Download-Options: noopen";
more_set_headers "X-Permitted-Cross-Domain-Policies: none";

# Content Security Policy
more_set_headers "Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self';";

# HSTS
more_set_headers "Strict-Transport-Security: max-age=31536000; includeSubDomains; preload";

# Disable unused HTTP methods
map $request_method $not_allowed_method {
    default 0;
    ~^(TRACE|DELETE|PUT)$ 1;
}

# Rate limiting zones
limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;
limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=general:10m rate=1r/s;

# Connection limiting
limit_conn_zone $binary_remote_addr zone=conn_limit_per_ip:10m;
limit_conn_zone $server_name zone=conn_limit_per_server:10m;

# Geo-blocking (requires GeoIP module)
map $geoip_country_code $blocked_country {
    default 0;
    CN 1;  # Block China
    RU 1;  # Block Russia
    KP 1;  # Block North Korea
}

# Bot detection
map $http_user_agent $blocked_agent {
    default 0;
    ~*malicious 1;
    ~*scanner 1;
    ~*bot 1;
}

# Common security rules
location ~ /\. {
    deny all;
    access_log off;
    log_not_found off;
}

location ~* \.(bak|backup|old|orig|tmp|~)$ {
    deny all;
    access_log off;
    log_not_found off;
}

location ~* \.(sql|log|conf)$ {
    deny all;
    access_log off;
    log_not_found off;
}
```

### WAF Integration (ModSecurity)

```nginx
# Requires nginx-module-modsecurity
server {
    listen 80;
    server_name example.com;
    
    modsecurity on;
    modsecurity_rules_file /etc/nginx/modsec/main.conf;
    
    location / {
        proxy_pass http://backend;
    }
}
```

```bash
# /etc/nginx/modsec/main.conf
Include /etc/nginx/modsec/modsecurity.conf
Include /usr/share/modsecurity-crs/crs-setup.conf
Include /usr/share/modsecurity-crs/rules/*.conf
```

## Performance Tuning

### System-Level Optimizations

```bash
# /etc/sysctl.conf
# Network optimizations
net.core.somaxconn = 65535
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_max_syn_backlog = 65535
net.ipv4.tcp_fin_timeout = 10
net.ipv4.tcp_keepalive_time = 600
net.ipv4.tcp_keepalive_intvl = 60
net.ipv4.tcp_keepalive_probes = 10
net.ipv4.tcp_tw_reuse = 1

# File descriptor limits
fs.file-max = 2097152

# Apply changes
sudo sysctl -p
```

```bash
# /etc/security/limits.conf
nginx soft nofile 65535
nginx hard nofile 65535
```

### Nginx Performance Configuration

```nginx
# /etc/nginx/nginx.conf
user nginx;
worker_processes auto;
worker_rlimit_nofile 65535;
worker_cpu_affinity auto;

events {
    worker_connections 4096;
    use epoll;
    multi_accept on;
    accept_mutex off;
}

http {
    # Basic optimizations
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 30;
    keepalive_requests 1000;
    
    # Buffer sizes
    client_body_buffer_size 128k;
    client_max_body_size 10m;
    client_header_buffer_size 1k;
    large_client_header_buffers 4 4k;
    output_buffers 1 32k;
    postpone_output 1460;
    
    # Timeouts
    client_header_timeout 3m;
    client_body_timeout 3m;
    send_timeout 3m;
    
    # File caching
    open_file_cache max=1000 inactive=20s;
    open_file_cache_valid 30s;
    open_file_cache_min_uses 5;
    open_file_cache_errors off;
    
    # Compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;
    
    # Connection pooling
    upstream backend {
        server backend1.example.com;
        server backend2.example.com;
        keepalive 32;
        keepalive_requests 100;
        keepalive_timeout 60s;
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

## Best Practices Summary

### Configuration Management

1. **Use version control** for configuration files
2. **Modularize configurations** with include files
3. **Test configurations** before deployment
4. **Document changes** and maintain changelog
5. **Use configuration templates** for consistency

### Security Best Practices

1. **Keep Nginx updated** to latest stable version
2. **Use strong SSL/TLS configuration**
3. **Implement proper access controls**
4. **Hide server information**
5. **Use security headers**
6. **Implement rate limiting**
7. **Monitor for suspicious activity**
8. **Regular security audits**

### Performance Best Practices

1. **Optimize worker processes** for your hardware
2. **Use appropriate buffer sizes**
3. **Enable compression**
4. **Implement caching strategies**
5. **Use HTTP/2**
6. **Optimize SSL/TLS**
7. **Monitor performance metrics**
8. **Use CDN for static content**

### Monitoring and Maintenance

1. **Set up comprehensive logging**
2. **Monitor key metrics**
3. **Implement health checks**
4. **Set up alerting**
5. **Regular log analysis**
6. **Capacity planning**
7. **Backup configurations**
8. **Document procedures**

### Deployment Best Practices

1. **Use blue-green deployments**
2. **Implement graceful reloads**
3. **Test in staging environment**
4. **Have rollback procedures**
5. **Monitor after deployment**
6. **Use automation tools**
7. **Maintain deployment logs**

## Advanced Configuration Examples

### Multi-tenant SaaS Application

```nginx
http {
    # Map subdomain to tenant
    map $host $tenant {
        ~^(?<subdomain>[^.]+)\.example\.com$ $subdomain;
        default "";
    }
    
    # Upstream per tenant
    upstream tenant_app {
        server app1.internal:3000;
        server app2.internal:3000;
    }
    
    server {
        listen 443 ssl http2;
        server_name *.example.com;
        
        ssl_certificate /etc/ssl/wildcard.crt;
        ssl_certificate_key /etc/ssl/wildcard.key;
        
        # Reject if no tenant
        if ($tenant = "") {
            return 404;
        }
        
        location / {
            proxy_pass http://tenant_app;
            proxy_set_header Host $host;
            proxy_set_header X-Tenant $tenant;
            proxy_set_header X-Real-IP $remote_addr;
        }
        
        # Tenant-specific static files
        location /static/ {
            root /var/www/tenants/$tenant;
            expires 1y;
        }
    }
}
```

### API Gateway

```nginx
http {
    # Rate limiting per API key
    map $http_x_api_key $api_client_id {
        "key1" "client1";
        "key2" "client2";
        default "anonymous";
    }
    
    limit_req_zone $api_client_id zone=api_limit:10m rate=100r/m;
    
    # Service discovery
    upstream user_service {
        server user1.internal:3001;
        server user2.internal:3001;
    }
    
    upstream order_service {
        server order1.internal:3002;
        server order2.internal:3002;
    }
    
    server {
        listen 443 ssl http2;
        server_name api.example.com;
        
        # API authentication
        location / {
            access_by_lua_block {
                local api_key = ngx.var.http_x_api_key
                if not api_key then
                    ngx.status = 401
                    ngx.say("API key required")
                    ngx.exit(401)
                end
                
                -- Validate API key (check database/cache)
                local redis = require "resty.redis"
                local red = redis:new()
                red:connect("127.0.0.1", 6379)
                
                local valid = red:get("api_key:" .. api_key)
                if not valid or valid == ngx.null then
                    ngx.status = 401
                    ngx.say("Invalid API key")
                    ngx.exit(401)
                end
            }
            
            limit_req zone=api_limit burst=20 nodelay;
            
            # Route based on path
            location /v1/users/ {
                proxy_pass http://user_service/;
            }
            
            location /v1/orders/ {
                proxy_pass http://order_service/;
            }
        }
    }
}
```

## Conclusion

This comprehensive guide covers Nginx from basic installation to advanced enterprise configurations. Key takeaways:

1. **Start with basics** - understand core concepts before moving to advanced features
2. **Security first** - always implement proper security measures
3. **Monitor everything** - comprehensive monitoring is crucial for production
4. **Performance matters** - optimize for your specific use case
5. **Test thoroughly** - always test configurations before production deployment
6. **Document well** - maintain good documentation for your team
7. **Stay updated** - keep Nginx and modules updated for security and features
8. **Plan for scale** - design configurations that can grow with your needs

### Next Steps

1. Practice with the examples in a test environment
2. Implement monitoring and logging for your use case
3. Explore additional modules that fit your requirements
4. Join the Nginx community for support and updates
5. Consider Nginx Plus for enterprise features
6. Implement automation for configuration management
7. Regular security audits and performance reviews

### Additional Resources

- **Official Documentation**: http://nginx.org/en/docs/
- **Nginx Blog**: https://www.nginx.com/blog/
- **Community Forums**: https://forum.nginx.org/
- **GitHub**: https://github.com/nginx/nginx
- **OpenResty**: https://openresty.org/
- **Nginx Plus**: https://www.nginx.com/products/nginx/

Remember: Nginx is a powerful tool, but with great power comes great responsibility. Always prioritize security, performance, and reliability in your configurations.

---

**End of Nginx Complete Guide**

*This guide provides a comprehensive foundation for working with Nginx. Continue learning by exploring specific use cases and advanced modules that match your requirements.*