# Monitoring, Logging, and Troubleshooting - Part 7

## Nginx Logging

### Understanding Log Formats

#### Default Log Format
```nginx
http {
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';
    
    access_log /var/log/nginx/access.log main;
    error_log /var/log/nginx/error.log warn;
}
```

#### Custom Log Formats

```nginx
http {
    # Detailed performance logging
    log_format detailed '$remote_addr - $remote_user [$time_local] '
                       '"$request" $status $body_bytes_sent '
                       '"$http_referer" "$http_user_agent" '
                       'rt=$request_time uct="$upstream_connect_time" '
                       'uht="$upstream_header_time" urt="$upstream_response_time" '
                       'cache="$upstream_cache_status" host="$host"';
    
    # JSON format for structured logging
    log_format json_combined escape=json
        '{'
            '"time_local":"$time_local",'  
            '"remote_addr":"$remote_addr",'  
            '"remote_user":"$remote_user",'  
            '"request":"$request",'  
            '"status": "$status",'  
            '"body_bytes_sent":"$body_bytes_sent",'  
            '"request_time":"$request_time",'  
            '"http_referrer":"$http_referer",'  
            '"http_user_agent":"$http_user_agent",'  
            '"upstream_addr":"$upstream_addr",'  
            '"upstream_response_time":"$upstream_response_time",'  
            '"upstream_cache_status":"$upstream_cache_status"'  
        '}';
    
    # Security-focused logging
    log_format security '$remote_addr - $remote_user [$time_local] '
                       '"$request" $status $body_bytes_sent '
                       '"$http_referer" "$http_user_agent" '
                       'country="$geoip_country_code" '
                       'ssl_protocol="$ssl_protocol" '
                       'ssl_cipher="$ssl_cipher"';
    
    server {
        access_log /var/log/nginx/access.log detailed;
        access_log /var/log/nginx/security.log security;
        access_log /var/log/nginx/json.log json_combined;
    }
}
```

### Conditional Logging

```nginx
http {
    # Map to determine if request should be logged
    map $status $loggable {
        ~^[23]  0;  # Don't log successful requests
        default 1;  # Log everything else
    }
    
    # Map based on user agent
    map $http_user_agent $log_ua {
        ~*bot     0;  # Don't log bots
        ~*crawler 0;  # Don't log crawlers
        default   1;  # Log everything else
    }
    
    server {
        # Conditional logging based on status
        access_log /var/log/nginx/access.log combined if=$loggable;
        
        # Conditional logging based on user agent
        access_log /var/log/nginx/human.log combined if=$log_ua;
        
        # Don't log static files
        location ~* \.(css|js|jpg|jpeg|png|gif|ico)$ {
            access_log off;
            expires 1y;
        }
        
        # Separate log for API endpoints
        location /api/ {
            access_log /var/log/nginx/api.log detailed;
            proxy_pass http://backend;
        }
    }
}
```

### Error Logging Levels

```nginx
# Error log levels (from least to most verbose):
# emerg, alert, crit, error, warn, notice, info, debug

http {
    error_log /var/log/nginx/error.log warn;
    
    server {
        # Override error log level for specific server
        error_log /var/log/nginx/site-error.log error;
        
        location /debug/ {
            # Debug logging for specific location
            error_log /var/log/nginx/debug.log debug;
        }
    }
}
```

## Real-time Monitoring

### Nginx Status Module

```nginx
server {
    listen 8080;
    server_name localhost;
    
    # Basic status
    location /nginx_status {
        stub_status on;
        access_log off;
        allow 127.0.0.1;
        allow 192.168.1.0/24;
        deny all;
    }
    
    # Extended status (Nginx Plus)
    location /status {
        status;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
    
    # Status in JSON format (Nginx Plus)
    location /status.json {
        status format=json;
        access_log off;
        allow 127.0.0.1;
        deny all;
    }
}
```

### Custom Monitoring Endpoints

```nginx
server {
    listen 8080;
    server_name localhost;
    
    # Health check endpoint
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
    
    # Detailed health check
    location /health/detailed {
        access_log off;
        
        # Check if upstream is available
        proxy_pass http://backend/health;
        proxy_connect_timeout 1s;
        proxy_read_timeout 1s;
        
        # Return 503 if upstream fails
        error_page 502 503 504 = @unhealthy;
    }
    
    location @unhealthy {
        return 503 "unhealthy\n";
        add_header Content-Type text/plain;
    }
    
    # Metrics endpoint
    location /metrics {
        access_log off;
        allow 127.0.0.1;
        allow 192.168.1.0/24;
        deny all;
        
        content_by_lua_block {
            -- Custom metrics using Lua (requires lua module)
            ngx.say("nginx_connections_active ", ngx.var.connections_active)
            ngx.say("nginx_connections_reading ", ngx.var.connections_reading)
            ngx.say("nginx_connections_writing ", ngx.var.connections_writing)
            ngx.say("nginx_connections_waiting ", ngx.var.connections_waiting)
        }
    }
}
```

## Log Analysis Tools

### GoAccess Real-time Analysis

```bash
# Install GoAccess
sudo apt install goaccess  # Ubuntu/Debian
sudo yum install goaccess  # CentOS/RHEL

# Real-time HTML dashboard
goaccess /var/log/nginx/access.log -o /var/www/html/report.html --log-format=COMBINED --real-time-html

# Terminal dashboard
goaccess /var/log/nginx/access.log --log-format=COMBINED

# Custom log format
goaccess /var/log/nginx/access.log --log-format='%h %^[%d:%t %^] "%r" %s %b "%R" "%u"' --date-format=%d/%b/%Y --time-format=%H:%M:%S
```

### AWK Scripts for Log Analysis

```bash
#!/bin/bash
# analyze-logs.sh

LOG_FILE="/var/log/nginx/access.log"

echo "=== Top 10 IP Addresses ==="
awk '{print $1}' $LOG_FILE | sort | uniq -c | sort -nr | head -10

echo "\n=== Top 10 Requested URLs ==="
awk '{print $7}' $LOG_FILE | sort | uniq -c | sort -nr | head -10

echo "\n=== HTTP Status Codes ==="
awk '{print $9}' $LOG_FILE | sort | uniq -c | sort -nr

echo "\n=== Top User Agents ==="
awk -F'"' '{print $6}' $LOG_FILE | sort | uniq -c | sort -nr | head -10

echo "\n=== Requests by Hour ==="
awk '{print $4}' $LOG_FILE | cut -c 14-15 | sort | uniq -c

echo "\n=== 4xx and 5xx Errors ==="
awk '$9 ~ /^[45]/ {print $9, $7}' $LOG_FILE | sort | uniq -c | sort -nr

echo "\n=== Large Requests (>1MB) ==="
awk '$10 > 1048576 {print $10, $7}' $LOG_FILE | sort -nr | head -10

echo "\n=== Slow Requests (>5s) ==="
awk '$NF > 5.000 {print $NF, $7}' $LOG_FILE | sort -nr | head -10
```

### Python Log Analyzer

```python
#!/usr/bin/env python3
# nginx-log-analyzer.py

import re
import sys
from collections import Counter, defaultdict
from datetime import datetime

def parse_nginx_log(log_file):
    # Nginx log pattern
    pattern = r'(\S+) - (\S+) \[(.*?)\] "(.*?)" (\d+) (\d+) "(.*?)" "(.*?)"'
    
    stats = {
        'total_requests': 0,
        'status_codes': Counter(),
        'ip_addresses': Counter(),
        'urls': Counter(),
        'user_agents': Counter(),
        'hourly_requests': defaultdict(int),
        'errors': [],
        'large_requests': []
    }
    
    with open(log_file, 'r') as f:
        for line in f:
            match = re.match(pattern, line.strip())
            if match:
                ip, user, timestamp, request, status, size, referer, user_agent = match.groups()
                
                stats['total_requests'] += 1
                stats['status_codes'][status] += 1
                stats['ip_addresses'][ip] += 1
                
                # Extract URL from request
                url = request.split()[1] if len(request.split()) > 1 else request
                stats['urls'][url] += 1
                stats['user_agents'][user_agent] += 1
                
                # Parse timestamp for hourly stats
                try:
                    dt = datetime.strptime(timestamp, '%d/%b/%Y:%H:%M:%S %z')
                    hour = dt.strftime('%H:00')
                    stats['hourly_requests'][hour] += 1
                except:
                    pass
                
                # Track errors
                if status.startswith(('4', '5')):
                    stats['errors'].append((ip, status, url, timestamp))
                
                # Track large requests
                if int(size) > 1048576:  # > 1MB
                    stats['large_requests'].append((ip, size, url))
    
    return stats

def print_report(stats):
    print(f"Total Requests: {stats['total_requests']}")
    print("\nTop 10 IP Addresses:")
    for ip, count in stats['ip_addresses'].most_common(10):
        print(f"  {ip}: {count}")
    
    print("\nStatus Code Distribution:")
    for status, count in stats['status_codes'].most_common():
        print(f"  {status}: {count}")
    
    print("\nTop 10 URLs:")
    for url, count in stats['urls'].most_common(10):
        print(f"  {url}: {count}")
    
    print("\nRequests by Hour:")
    for hour in sorted(stats['hourly_requests'].keys()):
        print(f"  {hour}: {stats['hourly_requests'][hour]}")
    
    print(f"\nTotal Errors: {len(stats['errors'])}")
    if stats['errors']:
        print("Recent Errors:")
        for error in stats['errors'][-5:]:
            print(f"  {error[1]} - {error[2]} from {error[0]}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 nginx-log-analyzer.py /path/to/access.log")
        sys.exit(1)
    
    log_file = sys.argv[1]
    stats = parse_nginx_log(log_file)
    print_report(stats)
```

## Monitoring with External Tools

### Prometheus + Grafana

#### Nginx Prometheus Exporter

```yaml
# docker-compose.yml for monitoring stack
version: '3.8'
services:
  nginx-exporter:
    image: nginx/nginx-prometheus-exporter:latest
    ports:
      - "9113:9113"
    command:
      - -nginx.scrape-uri=http://nginx:8080/nginx_status
    depends_on:
      - nginx
  
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
  
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

#### Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'nginx'
    static_configs:
      - targets: ['nginx-exporter:9113']
    scrape_interval: 5s
    metrics_path: /metrics
```

### ELK Stack Integration

#### Filebeat Configuration

```yaml
# filebeat.yml
filebeat.inputs:
- type: log
  enabled: true
  paths:
    - /var/log/nginx/access.log
  fields:
    logtype: nginx-access
  fields_under_root: true

- type: log
  enabled: true
  paths:
    - /var/log/nginx/error.log
  fields:
    logtype: nginx-error
  fields_under_root: true

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
  index: "nginx-logs-%{+yyyy.MM.dd}"

setup.template.name: "nginx"
setup.template.pattern: "nginx-*"
```

#### Logstash Configuration

```ruby
# logstash.conf
input {
  beats {
    port => 5044
  }
}

filter {
  if [logtype] == "nginx-access" {
    grok {
      match => { "message" => "%{NGINXACCESS}" }
    }
    
    date {
      match => [ "timestamp", "dd/MMM/yyyy:HH:mm:ss Z" ]
    }
    
    mutate {
      convert => { "response" => "integer" }
      convert => { "bytes" => "integer" }
      convert => { "responsetime" => "float" }
    }
  }
  
  if [logtype] == "nginx-error" {
    grok {
      match => { "message" => "%{NGINXERROR}" }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "nginx-logs-%{+YYYY.MM.dd}"
  }
}
```

## Troubleshooting Common Issues

### Configuration Testing

```bash
# Test configuration syntax
nginx -t

# Test specific configuration file
nginx -t -c /etc/nginx/sites-available/example.com

# Check configuration and show parsed result
nginx -T

# Show version and compiled modules
nginx -V
```

### Debug Mode

```nginx
# Enable debug logging (requires nginx compiled with --with-debug)
error_log /var/log/nginx/debug.log debug;

# Debug specific events
events {
    debug_connection 192.168.1.100;
    debug_connection 10.0.0.0/8;
}
```

### Common Error Diagnosis

#### 502 Bad Gateway

```bash
# Check if upstream is running
curl -I http://backend-server:8080/

# Check nginx error log
tail -f /var/log/nginx/error.log

# Check upstream connectivity
telnet backend-server 8080

# Check firewall rules
sudo iptables -L
sudo ufw status
```

#### 504 Gateway Timeout

```nginx
# Increase timeout values
server {
    location / {
        proxy_pass http://backend;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }
}
```

#### High Memory Usage

```bash
# Check nginx processes
ps aux | grep nginx

# Check memory usage
top -p $(pgrep nginx | tr '\n' ',')

# Check open files
lsof -p $(pgrep nginx)

# Check worker process limits
ulimit -n
```

### Performance Troubleshooting

#### Slow Response Times

```nginx
# Add timing information to logs
log_format timing '$remote_addr - $remote_user [$time_local] '
                 '"$request" $status $body_bytes_sent '
                 'rt=$request_time uct="$upstream_connect_time" '
                 'uht="$upstream_header_time" urt="$upstream_response_time"';

access_log /var/log/nginx/timing.log timing;
```

#### High CPU Usage

```bash
# Check worker process CPU usage
top -H -p $(pgrep nginx)

# Profile with perf (if available)
sudo perf top -p $(pgrep nginx)

# Check for infinite loops in configuration
nginx -t
```

### Debugging Tools and Scripts

#### Real-time Log Monitoring

```bash
#!/bin/bash
# monitor-nginx.sh

echo "Monitoring Nginx logs in real-time..."
echo "Press Ctrl+C to stop"

# Monitor access log for errors
tail -f /var/log/nginx/access.log | while read line; do
    status=$(echo $line | awk '{print $9}')
    if [[ $status =~ ^[45] ]]; then
        echo "ERROR: $line"
    fi
done &

# Monitor error log
tail -f /var/log/nginx/error.log &

wait
```

#### Connection Monitoring

```bash
#!/bin/bash
# connection-monitor.sh

while true; do
    echo "$(date): Active connections: $(curl -s http://localhost:8080/nginx_status | grep 'Active connections' | awk '{print $3}')"
    sleep 5
done
```

#### Automated Health Check

```bash
#!/bin/bash
# health-check.sh

SERVER="http://localhost"
EMAIL="admin@example.com"

response=$(curl -s -o /dev/null -w "%{http_code}" $SERVER)

if [ $response -ne 200 ]; then
    echo "ALERT: Server returned $response" | mail -s "Nginx Health Check Failed" $EMAIL
    
    # Try to restart nginx
    sudo systemctl restart nginx
    
    # Wait and check again
    sleep 10
    response=$(curl -s -o /dev/null -w "%{http_code}" $SERVER)
    
    if [ $response -eq 200 ]; then
        echo "Server recovered after restart" | mail -s "Nginx Recovered" $EMAIL
    else
        echo "Server still down after restart" | mail -s "Nginx Still Down" $EMAIL
    fi
else
    echo "$(date): Server is healthy (HTTP $response)"
fi
```

## Best Practices for Monitoring

1. **Set up proper log rotation** to prevent disk space issues
2. **Monitor key metrics**: response times, error rates, connection counts
3. **Use structured logging** (JSON) for better parsing
4. **Set up alerts** for critical issues
5. **Regular log analysis** to identify trends
6. **Monitor upstream health** in load balancer scenarios
7. **Track cache hit rates** for performance optimization
8. **Monitor SSL certificate expiration**
9. **Set up automated health checks**
10. **Use external monitoring** for comprehensive coverage

### Log Rotation Configuration

```bash
# /etc/logrotate.d/nginx
/var/log/nginx/*.log {
    daily
    missingok
    rotate 52
    compress
    delaycompress
    notifempty
    create 644 nginx adm
    sharedscripts
    prerotate
        if [ -d /etc/logrotate.d/httpd-prerotate ]; then \
            run-parts /etc/logrotate.d/httpd-prerotate; \
        fi \
    endscript
    postrotate
        invoke-rc.d nginx rotate >/dev/null 2>&1
    endscript
}
```

---

**Continue to Part 8: Advanced Topics and Best Practices**