# Server Blocks and Virtual Hosts - Part 2

## Understanding Server Blocks

Server blocks in Nginx are similar to virtual hosts in Apache. They allow you to host multiple websites on a single server by defining different configurations for different domain names or IP addresses.

## Basic Server Block Structure

```nginx
server {
    listen 80;                    # Port to listen on
    server_name example.com;      # Domain name
    root /var/www/example.com;    # Document root
    index index.html index.php;  # Default files
    
    # Location blocks
    location / {
        try_files $uri $uri/ =404;
    }
}
```

## Multiple Server Blocks

### Example: Hosting Multiple Sites

```nginx
# Site 1: example.com
server {
    listen 80;
    server_name example.com www.example.com;
    root /var/www/example.com;
    index index.html;
    
    access_log /var/log/nginx/example.com.access.log;
    error_log /var/log/nginx/example.com.error.log;
    
    location / {
        try_files $uri $uri/ =404;
    }
}

# Site 2: blog.example.com
server {
    listen 80;
    server_name blog.example.com;
    root /var/www/blog;
    index index.html index.php;
    
    access_log /var/log/nginx/blog.access.log;
    error_log /var/log/nginx/blog.error.log;
    
    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }
    
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php7.4-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}

# Site 3: api.example.com
server {
    listen 80;
    server_name api.example.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

## Server Name Matching

### Exact Match
```nginx
server_name example.com;
```

### Wildcard Match
```nginx
server_name *.example.com;     # Matches any subdomain
server_name example.*;         # Matches any TLD
```

### Regular Expression Match
```nginx
server_name ~^(www\.)?(.+)$;   # Regex pattern
```

### Multiple Server Names
```nginx
server_name example.com www.example.com site.example.com;
```

## Location Directives

### Basic Location Types

```nginx
# Exact match
location = /exact {
    # Only matches /exact
}

# Prefix match
location /prefix {
    # Matches /prefix, /prefix/, /prefix/anything
}

# Case-sensitive regex
location ~ \.php$ {
    # Matches files ending with .php
}

# Case-insensitive regex
location ~* \.(jpg|jpeg|png|gif)$ {
    # Matches image files (case insensitive)
}

# Priority prefix match
location ^~ /priority {
    # Higher priority than regex
}
```

### Location Priority Order

1. Exact match (`=`)
2. Priority prefix match (`^~`)
3. Regular expressions (`~` and `~*`)
4. Prefix match (no modifier)

### Common Location Examples

```nginx
server {
    listen 80;
    server_name example.com;
    root /var/www/example.com;
    
    # Root location
    location / {
        try_files $uri $uri/ =404;
    }
    
    # Static files with caching
    location ~* \.(css|js|jpg|jpeg|png|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
    
    # API proxy
    location /api/ {
        proxy_pass http://backend_server;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
    
    # Deny access to hidden files
    location ~ /\. {
        deny all;
    }
    
    # PHP processing
    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php-fpm.sock;
        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }
}
```

## Setting Up Virtual Hosts

### Step 1: Create Directory Structure
```bash
sudo mkdir -p /var/www/site1.com/html
sudo mkdir -p /var/www/site2.com/html
sudo chown -R $USER:$USER /var/www/site1.com/html
sudo chown -R $USER:$USER /var/www/site2.com/html
sudo chmod -R 755 /var/www
```

### Step 2: Create Sample Pages
```bash
echo "<h1>Welcome to Site 1</h1>" > /var/www/site1.com/html/index.html
echo "<h1>Welcome to Site 2</h1>" > /var/www/site2.com/html/index.html
```

### Step 3: Create Server Block Files

**Site 1 Configuration** (`/etc/nginx/sites-available/site1.com`):
```nginx
server {
    listen 80;
    listen [::]:80;
    
    root /var/www/site1.com/html;
    index index.html index.htm index.nginx-debian.html;
    
    server_name site1.com www.site1.com;
    
    location / {
        try_files $uri $uri/ =404;
    }
}
```

**Site 2 Configuration** (`/etc/nginx/sites-available/site2.com`):
```nginx
server {
    listen 80;
    listen [::]:80;
    
    root /var/www/site2.com/html;
    index index.html index.htm index.nginx-debian.html;
    
    server_name site2.com www.site2.com;
    
    location / {
        try_files $uri $uri/ =404;
    }
}
```

### Step 4: Enable Sites
```bash
sudo ln -s /etc/nginx/sites-available/site1.com /etc/nginx/sites-enabled/
sudo ln -s /etc/nginx/sites-available/site2.com /etc/nginx/sites-enabled/
```

### Step 5: Test and Reload
```bash
sudo nginx -t
sudo systemctl reload nginx
```

## Default Server

```nginx
server {
    listen 80 default_server;
    listen [::]:80 default_server;
    
    server_name _;
    
    return 444;  # Close connection without response
    # or
    # return 301 https://example.com$request_uri;
}
```

## Common Variables

- `$host` - Host header from request
- `$server_name` - Server name from server block
- `$request_uri` - Full original request URI
- `$uri` - Current URI in request
- `$document_root` - Root directory for current request
- `$remote_addr` - Client IP address
- `$scheme` - Request scheme (http or https)

## Best Practices

1. **Use separate files** for each site in `sites-available`
2. **Enable sites** by symlinking to `sites-enabled`
3. **Use descriptive names** for configuration files
4. **Set up proper logging** for each site
5. **Test configurations** before reloading
6. **Use consistent directory structure**

## Troubleshooting

### Check which server block is being used:
```bash
curl -H "Host: example.com" http://your-server-ip/
```

### View error logs:
```bash
sudo tail -f /var/log/nginx/error.log
```

### Test specific configuration:
```bash
sudo nginx -t -c /etc/nginx/sites-available/site1.com
```

---

**Continue to Part 3: Location Directives and URL Rewriting**