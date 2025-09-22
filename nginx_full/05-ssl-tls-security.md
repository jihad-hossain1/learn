# SSL/TLS Configuration and Security - Part 5

## SSL/TLS Fundamentals

### What is SSL/TLS?

- **SSL (Secure Sockets Layer)**: Legacy protocol for encrypting data
- **TLS (Transport Layer Security)**: Modern successor to SSL
- **Purpose**: Encrypt data in transit, authenticate servers, ensure data integrity

### TLS Versions

- **TLS 1.0/1.1**: Deprecated, security vulnerabilities
- **TLS 1.2**: Widely supported, secure
- **TLS 1.3**: Latest, fastest, most secure

## Basic SSL Configuration

### Simple HTTPS Setup

```nginx
server {
    listen 443 ssl;
    server_name example.com;
    
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    
    location / {
        root /var/www/html;
        index index.html;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name example.com;
    return 301 https://$server_name$request_uri;
}
```

### HTTP/2 Support

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;
    
    location / {
        root /var/www/html;
        index index.html;
    }
}
```

## Advanced SSL Configuration

### Modern SSL Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # Certificates
    ssl_certificate /etc/ssl/certs/example.com.crt;
    ssl_certificate_key /etc/ssl/private/example.com.key;
    
    # SSL Protocols
    ssl_protocols TLSv1.2 TLSv1.3;
    
    # Cipher Suites
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    
    # Session Settings
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;
    
    # OCSP Stapling
    ssl_stapling on;
    ssl_stapling_verify on;
    ssl_trusted_certificate /etc/ssl/certs/ca-certificates.crt;
    resolver 8.8.8.8 8.8.4.4 valid=300s;
    resolver_timeout 5s;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    location / {
        root /var/www/html;
        index index.html;
    }
}
```

### SSL Configuration for Different Use Cases

#### High Security (Government/Financial)

```nginx
server {
    listen 443 ssl http2;
    server_name secure.example.com;
    
    ssl_certificate /etc/ssl/certs/secure.example.com.crt;
    ssl_certificate_key /etc/ssl/private/secure.example.com.key;
    
    # Only TLS 1.3
    ssl_protocols TLSv1.3;
    
    # Strongest ciphers only
    ssl_ciphers TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256;
    ssl_prefer_server_ciphers off;
    
    # Enhanced session security
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    ssl_session_tickets off;
    
    # OCSP Must-Staple
    ssl_stapling on;
    ssl_stapling_verify on;
    
    # Strict security headers
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'" always;
    
    location / {
        root /var/www/secure;
        index index.html;
    }
}
```

#### Compatibility Mode (Legacy Support)

```nginx
server {
    listen 443 ssl http2;
    server_name legacy.example.com;
    
    ssl_certificate /etc/ssl/certs/legacy.example.com.crt;
    ssl_certificate_key /etc/ssl/private/legacy.example.com.key;
    
    # Support older TLS versions
    ssl_protocols TLSv1.2 TLSv1.3;
    
    # Broader cipher support
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384:DHE-RSA-CHACHA20-POLY1305:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-AES128-SHA:ECDHE-RSA-AES128-SHA:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES256-SHA:ECDHE-RSA-AES256-SHA:DHE-RSA-AES128-SHA256:DHE-RSA-AES256-SHA256:AES128-GCM-SHA256:AES256-GCM-SHA384:AES128-SHA256:AES256-SHA256:AES128-SHA:AES256-SHA:DES-CBC3-SHA;
    ssl_prefer_server_ciphers on;
    
    location / {
        root /var/www/legacy;
        index index.html;
    }
}
```

## Certificate Management

### Let's Encrypt with Certbot

#### Installation
```bash
# Ubuntu/Debian
sudo apt install certbot python3-certbot-nginx

# CentOS/RHEL
sudo yum install certbot python3-certbot-nginx
```

#### Obtain Certificate
```bash
# Automatic configuration
sudo certbot --nginx -d example.com -d www.example.com

# Manual certificate only
sudo certbot certonly --nginx -d example.com -d www.example.com
```

#### Auto-renewal
```bash
# Test renewal
sudo certbot renew --dry-run

# Add to crontab
0 12 * * * /usr/bin/certbot renew --quiet
```

### Self-Signed Certificates (Development)

```bash
# Generate private key
openssl genrsa -out server.key 2048

# Generate certificate signing request
openssl req -new -key server.key -out server.csr

# Generate self-signed certificate
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Or generate both key and certificate in one command
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
```

### Wildcard Certificates

```nginx
server {
    listen 443 ssl http2;
    server_name *.example.com;
    
    ssl_certificate /etc/ssl/certs/wildcard.example.com.crt;
    ssl_certificate_key /etc/ssl/private/wildcard.example.com.key;
    
    # Configuration based on subdomain
    set $subdomain "";
    if ($host ~ ^([^.]+)\.example\.com$) {
        set $subdomain $1;
    }
    
    location / {
        root /var/www/$subdomain;
        index index.html;
    }
}
```

## Security Headers

### Essential Security Headers

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # SSL configuration...
    
    # HTTP Strict Transport Security (HSTS)
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
    
    # Prevent clickjacking
    add_header X-Frame-Options DENY always;
    
    # Prevent MIME type sniffing
    add_header X-Content-Type-Options nosniff always;
    
    # XSS Protection
    add_header X-XSS-Protection "1; mode=block" always;
    
    # Referrer Policy
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    
    # Content Security Policy
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; frame-ancestors 'none';" always;
    
    # Permissions Policy
    add_header Permissions-Policy "geolocation=(), microphone=(), camera=()" always;
    
    location / {
        root /var/www/html;
        index index.html;
    }
}
```

### Advanced Content Security Policy

```nginx
# For a React/Vue.js application
add_header Content-Security-Policy "
    default-src 'self';
    script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.jsdelivr.net;
    style-src 'self' 'unsafe-inline' https://fonts.googleapis.com;
    font-src 'self' https://fonts.gstatic.com;
    img-src 'self' data: https:;
    connect-src 'self' https://api.example.com;
    frame-src 'none';
    object-src 'none';
    base-uri 'self';
    form-action 'self';
" always;

# For a WordPress site
add_header Content-Security-Policy "
    default-src 'self';
    script-src 'self' 'unsafe-inline' 'unsafe-eval';
    style-src 'self' 'unsafe-inline';
    img-src 'self' data: https:;
    font-src 'self' https://fonts.gstatic.com;
    connect-src 'self';
    frame-src 'self' https://www.youtube.com https://player.vimeo.com;
    object-src 'none';
" always;
```

## Rate Limiting and DDoS Protection

### Basic Rate Limiting

```nginx
http {
    # Define rate limiting zones
    limit_req_zone $binary_remote_addr zone=login:10m rate=5r/m;
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=general:10m rate=1r/s;
    
    server {
        listen 443 ssl http2;
        server_name example.com;
        
        # General rate limiting
        limit_req zone=general burst=5 nodelay;
        
        # Strict rate limiting for login
        location /login {
            limit_req zone=login burst=2 nodelay;
            proxy_pass http://backend;
        }
        
        # API rate limiting
        location /api/ {
            limit_req zone=api burst=20 nodelay;
            proxy_pass http://backend;
        }
    }
}
```

### Connection Limiting

```nginx
http {
    # Limit connections per IP
    limit_conn_zone $binary_remote_addr zone=conn_limit_per_ip:10m;
    limit_conn_zone $server_name zone=conn_limit_per_server:10m;
    
    server {
        listen 443 ssl http2;
        server_name example.com;
        
        # Limit connections
        limit_conn conn_limit_per_ip 10;
        limit_conn conn_limit_per_server 1000;
        
        location / {
            root /var/www/html;
        }
    }
}
```

### Geographic Blocking

```nginx
http {
    # GeoIP configuration (requires geoip module)
    geoip_country /usr/share/GeoIP/GeoIP.dat;
    
    # Map countries to allowed/blocked
    map $geoip_country_code $allowed_country {
        default yes;
        CN no;  # Block China
        RU no;  # Block Russia
        KP no;  # Block North Korea
    }
    
    server {
        listen 443 ssl http2;
        server_name example.com;
        
        # Block based on country
        if ($allowed_country = no) {
            return 403;
        }
        
        location / {
            root /var/www/html;
        }
    }
}
```

## Authentication and Access Control

### Basic Authentication

```nginx
server {
    listen 443 ssl http2;
    server_name admin.example.com;
    
    # Admin area protection
    location /admin {
        auth_basic "Admin Area";
        auth_basic_user_file /etc/nginx/.htpasswd;
        
        root /var/www/admin;
        index index.html;
    }
}
```

#### Create Password File
```bash
# Install htpasswd utility
sudo apt install apache2-utils

# Create password file
sudo htpasswd -c /etc/nginx/.htpasswd admin

# Add more users
sudo htpasswd /etc/nginx/.htpasswd user2
```

### IP-based Access Control

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # Allow specific IPs
    location /admin {
        allow 192.168.1.0/24;
        allow 10.0.0.0/8;
        deny all;
        
        root /var/www/admin;
    }
    
    # Block specific IPs
    location / {
        deny 192.168.1.100;
        deny 10.0.0.50;
        allow all;
        
        root /var/www/html;
    }
}
```

### Client Certificate Authentication

```nginx
server {
    listen 443 ssl http2;
    server_name secure.example.com;
    
    ssl_certificate /etc/ssl/certs/server.crt;
    ssl_certificate_key /etc/ssl/private/server.key;
    
    # Client certificate configuration
    ssl_client_certificate /etc/ssl/certs/ca.crt;
    ssl_verify_client on;
    ssl_verify_depth 2;
    
    location / {
        # Client certificate info available in variables
        proxy_set_header X-SSL-Client-Cert $ssl_client_cert;
        proxy_set_header X-SSL-Client-DN $ssl_client_s_dn;
        proxy_set_header X-SSL-Client-Verify $ssl_client_verify;
        
        proxy_pass http://backend;
    }
}
```

## Security Best Practices

### Hide Nginx Version

```nginx
http {
    server_tokens off;
    
    # Custom server header (optional)
    more_set_headers "Server: MyServer";
}
```

### Disable Unused HTTP Methods

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # Only allow specific methods
    if ($request_method !~ ^(GET|HEAD|POST|PUT|DELETE|OPTIONS)$) {
        return 405;
    }
    
    location / {
        # Limit methods per location
        limit_except GET HEAD POST {
            deny all;
        }
        
        root /var/www/html;
    }
}
```

### Secure File Upload

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # File upload location
    location /upload {
        client_max_body_size 10M;
        client_body_timeout 60s;
        
        # Prevent execution of uploaded files
        location ~* \.(php|pl|py|jsp|asp|sh|cgi)$ {
            return 403;
        }
        
        proxy_pass http://backend;
    }
}
```

### Security Configuration Template

```nginx
# /etc/nginx/conf.d/security.conf

# Hide nginx version
server_tokens off;

# Prevent access to hidden files
location ~ /\. {
    deny all;
    access_log off;
    log_not_found off;
}

# Prevent access to backup files
location ~* \.(bak|backup|old|orig|original|tmp|temp|~)$ {
    deny all;
    access_log off;
    log_not_found off;
}

# Security headers
add_header X-Frame-Options DENY always;
add_header X-Content-Type-Options nosniff always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;

# Rate limiting
limit_req_zone $binary_remote_addr zone=general:10m rate=1r/s;
limit_req zone=general burst=5 nodelay;

# Connection limiting
limit_conn_zone $binary_remote_addr zone=conn_limit:10m;
limit_conn conn_limit 10;
```

## SSL Testing and Monitoring

### Test SSL Configuration

```bash
# Test with OpenSSL
openssl s_client -connect example.com:443 -servername example.com

# Test specific TLS version
openssl s_client -connect example.com:443 -tls1_2
openssl s_client -connect example.com:443 -tls1_3

# Check certificate details
openssl x509 -in certificate.crt -text -noout

# Check certificate expiration
openssl x509 -in certificate.crt -noout -dates
```

### Online SSL Testing Tools

- **SSL Labs**: https://www.ssllabs.com/ssltest/
- **Mozilla Observatory**: https://observatory.mozilla.org/
- **Security Headers**: https://securityheaders.com/

### Certificate Monitoring Script

```bash
#!/bin/bash
# check-ssl-expiry.sh

DOMAIN="example.com"
THRESHOLD=30  # Days before expiration to alert

EXPIRY_DATE=$(openssl s_client -connect $DOMAIN:443 -servername $DOMAIN 2>/dev/null | openssl x509 -noout -dates | grep notAfter | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY_DATE" +%s)
CURRENT_EPOCH=$(date +%s)
DAYS_LEFT=$(( (EXPIRY_EPOCH - CURRENT_EPOCH) / 86400 ))

if [ $DAYS_LEFT -lt $THRESHOLD ]; then
    echo "WARNING: SSL certificate for $DOMAIN expires in $DAYS_LEFT days!"
    # Send alert (email, Slack, etc.)
else
    echo "SSL certificate for $DOMAIN is valid for $DAYS_LEFT more days."
fi
```

---

**Continue to Part 6: Caching and Performance Optimization**