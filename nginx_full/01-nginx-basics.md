# Nginx Basics - Part 1

## What is Nginx?

Nginx (pronounced "engine-x") is a high-performance web server, reverse proxy server, and load balancer. Originally created by Igor Sysoev in 2004, it has become one of the most popular web servers in the world.

## Key Features

- **High Performance**: Can handle thousands of concurrent connections
- **Low Memory Usage**: Efficient memory management
- **Reverse Proxy**: Can act as a proxy server for HTTP, HTTPS, SMTP, POP3, and IMAP
- **Load Balancing**: Distribute incoming requests across multiple servers
- **SSL/TLS Termination**: Handle SSL encryption/decryption
- **Static Content Serving**: Efficiently serve static files
- **Caching**: Built-in caching capabilities

## Installation

### Ubuntu/Debian
```bash
sudo apt update
sudo apt install nginx
```

### CentOS/RHEL
```bash
sudo yum install nginx
# or for newer versions
sudo dnf install nginx
```

### Windows
1. Download from official website: http://nginx.org/en/download.html
2. Extract the zip file
3. Run nginx.exe from command prompt

### macOS
```bash
brew install nginx
```

## Basic Commands

### Start Nginx
```bash
sudo systemctl start nginx
# or
sudo nginx
```

### Stop Nginx
```bash
sudo systemctl stop nginx
# or
sudo nginx -s stop
```

### Restart Nginx
```bash
sudo systemctl restart nginx
# or
sudo nginx -s reload
```

### Check Status
```bash
sudo systemctl status nginx
```

### Test Configuration
```bash
sudo nginx -t
```

## Directory Structure

### Main Configuration Files
- `/etc/nginx/nginx.conf` - Main configuration file
- `/etc/nginx/sites-available/` - Available site configurations
- `/etc/nginx/sites-enabled/` - Enabled site configurations
- `/etc/nginx/conf.d/` - Additional configuration files

### Log Files
- `/var/log/nginx/access.log` - Access logs
- `/var/log/nginx/error.log` - Error logs

### Web Root
- `/var/www/html/` - Default web root directory

## Basic Configuration Structure

```nginx
# Main context
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;

# Events context
events {
    worker_connections 1024;
}

# HTTP context
http {
    # MIME types
    include /etc/nginx/mime.types;
    default_type application/octet-stream;
    
    # Logging
    access_log /var/log/nginx/access.log;
    
    # Server context
    server {
        listen 80;
        server_name example.com;
        root /var/www/html;
        index index.html index.htm;
        
        # Location context
        location / {
            try_files $uri $uri/ =404;
        }
    }
}
```

## Your First Nginx Server

1. Create a simple HTML file:
```bash
sudo mkdir -p /var/www/mysite
echo "<h1>Hello from Nginx!</h1>" | sudo tee /var/www/mysite/index.html
```

2. Create a server configuration:
```nginx
server {
    listen 80;
    server_name localhost;
    root /var/www/mysite;
    index index.html;
    
    location / {
        try_files $uri $uri/ =404;
    }
}
```

3. Test and reload:
```bash
sudo nginx -t
sudo systemctl reload nginx
```

## Next Steps

In the next part, we'll cover:
- Server blocks and virtual hosts
- Location directives
- Basic security configurations
- Serving static files efficiently

---

**Continue to Part 2: Server Blocks and Virtual Hosts**