# Part 8: Networking Fundamentals for DevOps

## Introduction

Networking is the backbone of modern DevOps infrastructure. Understanding networking concepts, protocols, and tools is essential for designing, deploying, and troubleshooting distributed systems. This guide covers fundamental networking concepts, practical tools, and real-world scenarios that DevOps engineers encounter daily.

## OSI Model and TCP/IP Stack

### OSI Model (7 Layers)

```
┌─────────────────┬─────────────────┬─────────────────┐
│   Layer 7       │   Application   │   HTTP, HTTPS,  │
│   Application   │   Layer         │   FTP, SMTP     │
├─────────────────┼─────────────────┼─────────────────┤
│   Layer 6       │   Presentation  │   SSL/TLS,      │
│   Presentation  │   Layer         │   Encryption    │
├─────────────────┼─────────────────┼─────────────────┤
│   Layer 5       │   Session       │   NetBIOS,      │
│   Session       │   Layer         │   RPC           │
├─────────────────┼─────────────────┼─────────────────┤
│   Layer 4       │   Transport     │   TCP, UDP      │
│   Transport     │   Layer         │                 │
├─────────────────┼─────────────────┼─────────────────┤
│   Layer 3       │   Network       │   IP, ICMP,     │
│   Network       │   Layer         │   Routing       │
├─────────────────┼─────────────────┼─────────────────┤
│   Layer 2       │   Data Link     │   Ethernet,     │
│   Data Link     │   Layer         │   WiFi, MAC     │
├─────────────────┼─────────────────┼─────────────────┤
│   Layer 1       │   Physical      │   Cables,       │
│   Physical      │   Layer         │   Radio Waves   │
└─────────────────┴─────────────────┴─────────────────┘
```

### TCP/IP Model (4 Layers)

```
┌─────────────────┬─────────────────┬─────────────────┐
│   Application   │   Layer 4       │   HTTP, FTP,    │
│   Layer         │                 │   SMTP, DNS     │
├─────────────────┼─────────────────┼─────────────────┤
│   Transport     │   Layer 3       │   TCP, UDP      │
│   Layer         │                 │                 │
├─────────────────┼─────────────────┼─────────────────┤
│   Internet      │   Layer 2       │   IP, ICMP,     │
│   Layer         │                 │   ARP           │
├─────────────────┼─────────────────┼─────────────────┤
│   Network       │   Layer 1       │   Ethernet,     │
│   Access Layer  │                 │   WiFi          │
└─────────────────┴─────────────────┴─────────────────┘
```

## IP Addressing and Subnetting

### IPv4 Addressing

**IPv4 Address Structure**:
- 32-bit address divided into 4 octets
- Each octet ranges from 0-255
- Example: 192.168.1.100

**Address Classes**:
```
Class A: 1.0.0.0    to 126.255.255.255  (/8)  - Large networks
Class B: 128.0.0.0  to 191.255.255.255  (/16) - Medium networks
Class C: 192.0.0.0  to 223.255.255.255  (/24) - Small networks
Class D: 224.0.0.0  to 239.255.255.255        - Multicast
Class E: 240.0.0.0  to 255.255.255.255        - Reserved
```

**Private IP Ranges**:
```
Class A: 10.0.0.0/8        (10.0.0.0 - 10.255.255.255)
Class B: 172.16.0.0/12     (172.16.0.0 - 172.31.255.255)
Class C: 192.168.0.0/16    (192.168.0.0 - 192.168.255.255)
```

### Subnetting

**CIDR Notation**:
```
/24 = 255.255.255.0   (256 addresses, 254 usable)
/25 = 255.255.255.128 (128 addresses, 126 usable)
/26 = 255.255.255.192 (64 addresses, 62 usable)
/27 = 255.255.255.224 (32 addresses, 30 usable)
/28 = 255.255.255.240 (16 addresses, 14 usable)
/29 = 255.255.255.248 (8 addresses, 6 usable)
/30 = 255.255.255.252 (4 addresses, 2 usable)
```

**Subnetting Example**:
```
Network: 192.168.1.0/24
Subnet 1: 192.168.1.0/26   (192.168.1.1 - 192.168.1.62)
Subnet 2: 192.168.1.64/26  (192.168.1.65 - 192.168.1.126)
Subnet 3: 192.168.1.128/26 (192.168.1.129 - 192.168.1.190)
Subnet 4: 192.168.1.192/26 (192.168.1.193 - 192.168.1.254)
```

**Subnet Calculator Script**:
```python
#!/usr/bin/env python3
# subnet_calculator.py

import ipaddress
import sys

def calculate_subnet(network_str):
    try:
        network = ipaddress.IPv4Network(network_str, strict=False)
        
        print(f"Network: {network}")
        print(f"Network Address: {network.network_address}")
        print(f"Broadcast Address: {network.broadcast_address}")
        print(f"Subnet Mask: {network.netmask}")
        print(f"Wildcard Mask: {network.hostmask}")
        print(f"Total Addresses: {network.num_addresses}")
        print(f"Usable Addresses: {network.num_addresses - 2}")
        print(f"First Usable: {list(network.hosts())[0]}")
        print(f"Last Usable: {list(network.hosts())[-1]}")
        
    except ValueError as e:
        print(f"Error: {e}")

def subnet_network(network_str, num_subnets):
    try:
        network = ipaddress.IPv4Network(network_str, strict=False)
        
        # Calculate required subnet bits
        import math
        subnet_bits = math.ceil(math.log2(num_subnets))
        new_prefix = network.prefixlen + subnet_bits
        
        if new_prefix > 30:
            print("Error: Too many subnets requested")
            return
        
        subnets = list(network.subnets(new_prefix=new_prefix))
        
        print(f"Original Network: {network}")
        print(f"Number of Subnets: {len(subnets)}")
        print(f"New Subnet Mask: /{new_prefix}")
        print("\nSubnets:")
        
        for i, subnet in enumerate(subnets[:num_subnets], 1):
            hosts = list(subnet.hosts())
            print(f"Subnet {i}: {subnet}")
            print(f"  Range: {hosts[0]} - {hosts[-1]}")
            print(f"  Usable Addresses: {len(hosts)}")
            print()
            
    except ValueError as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print("Usage: python subnet_calculator.py <network> [num_subnets]")
        print("Example: python subnet_calculator.py 192.168.1.0/24")
        print("Example: python subnet_calculator.py 192.168.1.0/24 4")
        sys.exit(1)
    
    network = sys.argv[1]
    
    if len(sys.argv) == 3:
        num_subnets = int(sys.argv[2])
        subnet_network(network, num_subnets)
    else:
        calculate_subnet(network)
```

### IPv6 Addressing

**IPv6 Address Structure**:
- 128-bit address
- Written in hexadecimal notation
- Example: 2001:0db8:85a3:0000:0000:8a2e:0370:7334
- Compressed: 2001:db8:85a3::8a2e:370:7334

**IPv6 Address Types**:
```
Unicast:   One-to-one communication
Multicast: One-to-many communication
Anycast:   One-to-nearest communication
```

## DNS (Domain Name System)

### DNS Hierarchy

```
                    Root (".") 
                       |
        ┌──────────────┼──────────────┐
       com            org            net
        |
    ┌───┼───┐
  google  amazon  microsoft
    |
  ┌─┼─┐
 www mail drive
```

### DNS Record Types

```
A Record:     Maps domain to IPv4 address
AAAA Record:  Maps domain to IPv6 address
CNAME:        Canonical name (alias)
MX:           Mail exchange server
NS:           Name server
PTR:          Reverse DNS lookup
SOA:          Start of authority
TXT:          Text records (SPF, DKIM, etc.)
SRV:          Service records
```

### DNS Tools and Commands

```bash
# DNS Lookup Tools
nslookup google.com              # Basic DNS lookup
dig google.com                   # Detailed DNS information
dig @8.8.8.8 google.com          # Query specific DNS server
dig google.com MX                # Query specific record type
dig +trace google.com            # Trace DNS resolution path

# Reverse DNS lookup
dig -x 8.8.8.8                   # Reverse lookup
nslookup 8.8.8.8                 # Reverse lookup

# DNS cache management
sudo systemctl flush-dns         # Flush DNS cache (systemd)
sudo dscacheutil -flushcache     # Flush DNS cache (macOS)
ipconfig /flushdns               # Flush DNS cache (Windows)
```

### DNS Configuration

**BIND DNS Server Configuration**:
```bash
# /etc/bind/named.conf.local
zone "example.com" {
    type master;
    file "/etc/bind/db.example.com";
};

zone "1.168.192.in-addr.arpa" {
    type master;
    file "/etc/bind/db.192.168.1";
};
```

**Zone File Example**:
```bash
# /etc/bind/db.example.com
$TTL    604800
@       IN      SOA     ns1.example.com. admin.example.com. (
                              2023010101         ; Serial
                              604800             ; Refresh
                              86400              ; Retry
                              2419200            ; Expire
                              604800 )           ; Negative Cache TTL

; Name servers
@       IN      NS      ns1.example.com.
@       IN      NS      ns2.example.com.

; A records
ns1     IN      A       192.168.1.10
ns2     IN      A       192.168.1.11
www     IN      A       192.168.1.20
mail    IN      A       192.168.1.30
ftp     IN      A       192.168.1.40

; CNAME records
webmail IN      CNAME   mail
api     IN      CNAME   www

; MX records
@       IN      MX      10      mail.example.com.

; TXT records
@       IN      TXT     "v=spf1 mx a ~all"
_dmarc  IN      TXT     "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
```

## Network Protocols

### HTTP/HTTPS

**HTTP Methods**:
```
GET:     Retrieve data
POST:    Submit data
PUT:     Update/create resource
DELETE:  Remove resource
PATCH:   Partial update
HEAD:    Get headers only
OPTIONS: Get allowed methods
```

**HTTP Status Codes**:
```
1xx: Informational
  100 Continue
  101 Switching Protocols

2xx: Success
  200 OK
  201 Created
  204 No Content

3xx: Redirection
  301 Moved Permanently
  302 Found (Temporary Redirect)
  304 Not Modified

4xx: Client Error
  400 Bad Request
  401 Unauthorized
  403 Forbidden
  404 Not Found
  429 Too Many Requests

5xx: Server Error
  500 Internal Server Error
  502 Bad Gateway
  503 Service Unavailable
  504 Gateway Timeout
```

**HTTP Headers**:
```bash
# Request Headers
Host: example.com
User-Agent: Mozilla/5.0...
Accept: text/html,application/xhtml+xml
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive
Authorization: Bearer token123

# Response Headers
Content-Type: text/html; charset=UTF-8
Content-Length: 1234
Server: nginx/1.18.0
Cache-Control: max-age=3600
Set-Cookie: sessionid=abc123; HttpOnly; Secure
X-Frame-Options: DENY
Strict-Transport-Security: max-age=31536000
```

### SSL/TLS

**TLS Handshake Process**:
```
1. Client Hello
   - Supported cipher suites
   - Random number
   - Session ID

2. Server Hello
   - Selected cipher suite
   - Random number
   - Session ID
   - Certificate

3. Client Key Exchange
   - Pre-master secret (encrypted)
   - Certificate verification

4. Finished
   - Both sides compute session keys
   - Secure communication begins
```

**SSL Certificate Management**:
```bash
# Generate private key
openssl genrsa -out private.key 2048

# Generate certificate signing request (CSR)
openssl req -new -key private.key -out certificate.csr

# Generate self-signed certificate
openssl req -x509 -new -nodes -key private.key -sha256 -days 365 -out certificate.crt

# View certificate information
openssl x509 -in certificate.crt -text -noout

# Test SSL connection
openssl s_client -connect example.com:443

# Check certificate expiration
echo | openssl s_client -connect example.com:443 2>/dev/null | openssl x509 -noout -dates
```

### TCP and UDP

**TCP (Transmission Control Protocol)**:
- Connection-oriented
- Reliable delivery
- Error checking and correction
- Flow control
- Used by: HTTP, HTTPS, FTP, SSH

**UDP (User Datagram Protocol)**:
- Connectionless
- Fast but unreliable
- No error correction
- Used by: DNS, DHCP, streaming media

**TCP Connection States**:
```bash
# View TCP connections
netstat -tuln                    # Linux/macOS
Get-NetTCPConnection            # Windows PowerShell

# TCP connection states
LISTEN:      Waiting for connections
ESTABLISHED: Active connection
TIME_WAIT:   Connection closed, waiting
CLOSE_WAIT:  Remote end closed connection
FIN_WAIT:    Local end closed connection
```

## Network Tools and Troubleshooting

### Essential Network Commands

```bash
# Connectivity Testing
ping google.com                  # Test basic connectivity
ping -c 4 google.com            # Send 4 packets
ping6 google.com                # IPv6 ping
traceroute google.com           # Trace network path
tracepath google.com            # Alternative to traceroute
mtr google.com                  # Continuous traceroute

# Port Testing
telnet google.com 80            # Test TCP port
nc -zv google.com 80            # Test port with netcat
nmap -p 80,443 google.com       # Port scanning

# Network Interface Information
ip addr show                    # Show interfaces (Linux)
ifconfig                        # Show interfaces (Linux/macOS)
ip route show                   # Show routing table
route -n                        # Show routing table
arp -a                          # Show ARP table

# Network Statistics
netstat -tuln                   # Show listening ports
netstat -rn                     # Show routing table
ss -tuln                        # Modern netstat alternative
lsof -i                         # Show network connections
lsof -i :80                     # Show processes using port 80
```

### Advanced Network Troubleshooting

**Network Diagnostic Script**:
```bash
#!/bin/bash
# network-diagnostic.sh

TARGET_HOST="${1:-google.com}"
TARGET_PORT="${2:-80}"
LOG_FILE="/tmp/network-diagnostic-$(date +%Y%m%d-%H%M%S).log"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

log "Starting network diagnostic for $TARGET_HOST:$TARGET_PORT"

# Basic connectivity test
log "=== Basic Connectivity Test ==="
if ping -c 3 "$TARGET_HOST" &>/dev/null; then
    log "✓ Ping to $TARGET_HOST successful"
else
    log "✗ Ping to $TARGET_HOST failed"
fi

# DNS resolution test
log "\n=== DNS Resolution Test ==="
if nslookup "$TARGET_HOST" &>/dev/null; then
    log "✓ DNS resolution for $TARGET_HOST successful"
    nslookup "$TARGET_HOST" | tee -a "$LOG_FILE"
else
    log "✗ DNS resolution for $TARGET_HOST failed"
fi

# Port connectivity test
log "\n=== Port Connectivity Test ==="
if nc -z -w5 "$TARGET_HOST" "$TARGET_PORT" &>/dev/null; then
    log "✓ Port $TARGET_PORT on $TARGET_HOST is open"
else
    log "✗ Port $TARGET_PORT on $TARGET_HOST is closed or filtered"
fi

# Traceroute test
log "\n=== Traceroute Test ==="
log "Tracing route to $TARGET_HOST:"
traceroute "$TARGET_HOST" 2>&1 | tee -a "$LOG_FILE"

# Local network configuration
log "\n=== Local Network Configuration ==="
log "Network interfaces:"
ip addr show 2>&1 | tee -a "$LOG_FILE"

log "\nRouting table:"
ip route show 2>&1 | tee -a "$LOG_FILE"

log "\nDNS configuration:"
cat /etc/resolv.conf 2>&1 | tee -a "$LOG_FILE"

# Network statistics
log "\n=== Network Statistics ==="
log "Active connections:"
ss -tuln 2>&1 | tee -a "$LOG_FILE"

log "\nNetwork interface statistics:"
cat /proc/net/dev 2>&1 | tee -a "$LOG_FILE"

log "\nNetwork diagnostic completed. Log saved to: $LOG_FILE"
```

### Packet Capture and Analysis

**tcpdump Examples**:
```bash
# Basic packet capture
sudo tcpdump -i eth0                    # Capture on interface eth0
sudo tcpdump -i any                     # Capture on all interfaces
sudo tcpdump host google.com            # Capture traffic to/from host
sudo tcpdump port 80                    # Capture HTTP traffic
sudo tcpdump -w capture.pcap            # Save to file
sudo tcpdump -r capture.pcap            # Read from file

# Advanced filters
sudo tcpdump 'tcp port 80 and host google.com'
sudo tcpdump 'udp port 53'              # DNS traffic
sudo tcpdump 'icmp'                     # ICMP traffic
sudo tcpdump 'tcp[tcpflags] & tcp-syn != 0'  # SYN packets

# Detailed output
sudo tcpdump -v -i eth0                 # Verbose output
sudo tcpdump -X -i eth0                 # Hex and ASCII output
sudo tcpdump -A -i eth0                 # ASCII output
```

**Wireshark Command Line (tshark)**:
```bash
# Basic capture
tshark -i eth0                          # Live capture
tshark -r capture.pcap                  # Read file
tshark -w output.pcap -i eth0           # Write to file

# Filters
tshark -f "tcp port 80"                 # Capture filter
tshark -Y "http.request.method == GET"  # Display filter
tshark -Y "dns"                         # DNS traffic
tshark -Y "tcp.flags.syn == 1"          # SYN packets

# Statistics
tshark -z conv,tcp                      # TCP conversations
tshark -z hosts                         # Host statistics
tshark -z io,phs                        # Protocol hierarchy
```

## Load Balancing and Reverse Proxies

### Nginx as Reverse Proxy

**Basic Reverse Proxy Configuration**:
```nginx
# /etc/nginx/sites-available/app
server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Timeouts
        proxy_connect_timeout 30s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;
        
        # Buffer settings
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
    }
    
    # Health check endpoint
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}

# Upstream backend servers
upstream backend {
    least_conn;  # Load balancing method
    
    server 192.168.1.10:8080 weight=3 max_fails=3 fail_timeout=30s;
    server 192.168.1.11:8080 weight=2 max_fails=3 fail_timeout=30s;
    server 192.168.1.12:8080 weight=1 max_fails=3 fail_timeout=30s backup;
    
    # Health checks (nginx plus)
    # health_check interval=5s fails=3 passes=2;
}
```

**SSL Termination**:
```nginx
server {
    listen 443 ssl http2;
    server_name example.com;
    
    # SSL configuration
    ssl_certificate /etc/ssl/certs/example.com.crt;
    ssl_certificate_key /etc/ssl/private/example.com.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    
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

### HAProxy Configuration

**HAProxy Load Balancer**:
```bash
# /etc/haproxy/haproxy.cfg
global
    daemon
    user haproxy
    group haproxy
    log stdout local0
    chroot /var/lib/haproxy
    stats socket /run/haproxy/admin.sock mode 660 level admin
    stats timeout 30s
    
    # SSL configuration
    ssl-default-bind-ciphers ECDHE+AESGCM:ECDHE+CHACHA20:RSA+AESGCM:RSA+SHA256:!aNULL:!MD5:!DSS
    ssl-default-bind-options ssl-min-ver TLSv1.2 no-tls-tickets

defaults
    mode http
    log global
    option httplog
    option dontlognull
    option log-health-checks
    option forwardfor
    option http-server-close
    
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms
    
    errorfile 400 /etc/haproxy/errors/400.http
    errorfile 403 /etc/haproxy/errors/403.http
    errorfile 408 /etc/haproxy/errors/408.http
    errorfile 500 /etc/haproxy/errors/500.http
    errorfile 502 /etc/haproxy/errors/502.http
    errorfile 503 /etc/haproxy/errors/503.http
    errorfile 504 /etc/haproxy/errors/504.http

# Frontend configuration
frontend web_frontend
    bind *:80
    bind *:443 ssl crt /etc/ssl/certs/example.com.pem
    
    # Redirect HTTP to HTTPS
    redirect scheme https if !{ ssl_fc }
    
    # ACLs
    acl is_api path_beg /api
    acl is_static path_beg /static
    
    # Routing
    use_backend api_backend if is_api
    use_backend static_backend if is_static
    default_backend web_backend

# Backend configurations
backend web_backend
    balance roundrobin
    option httpchk GET /health
    
    server web1 192.168.1.10:8080 check weight 100
    server web2 192.168.1.11:8080 check weight 100
    server web3 192.168.1.12:8080 check weight 50 backup

backend api_backend
    balance leastconn
    option httpchk GET /api/health
    
    server api1 192.168.1.20:8080 check
    server api2 192.168.1.21:8080 check
    server api3 192.168.1.22:8080 check

backend static_backend
    balance source
    option httpchk GET /static/health.txt
    
    server static1 192.168.1.30:8080 check
    server static2 192.168.1.31:8080 check

# Statistics page
listen stats
    bind *:8404
    stats enable
    stats uri /stats
    stats refresh 30s
    stats admin if TRUE
```

## Firewalls and Security

### iptables (Linux)

**Basic iptables Rules**:
```bash
# View current rules
sudo iptables -L -n -v
sudo iptables -t nat -L -n -v

# Basic firewall setup
#!/bin/bash
# firewall-setup.sh

# Flush existing rules
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X
iptables -t mangle -F
iptables -t mangle -X

# Set default policies
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# Allow loopback traffic
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Allow established and related connections
iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# Allow SSH (change port as needed)
iptables -A INPUT -p tcp --dport 22 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT

# Allow HTTP and HTTPS
iptables -A INPUT -p tcp --dport 80 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT

# Allow DNS
iptables -A INPUT -p udp --dport 53 -j ACCEPT
iptables -A INPUT -p tcp --dport 53 -j ACCEPT

# Allow ping (ICMP)
iptables -A INPUT -p icmp --icmp-type echo-request -j ACCEPT

# Rate limiting for SSH
iptables -A INPUT -p tcp --dport 22 -m recent --name ssh --set
iptables -A INPUT -p tcp --dport 22 -m recent --name ssh --rcheck --seconds 60 --hitcount 4 -j DROP

# Log dropped packets
iptables -A INPUT -m limit --limit 5/min -j LOG --log-prefix "iptables denied: " --log-level 7

# Save rules
iptables-save > /etc/iptables/rules.v4

echo "Firewall rules applied successfully"
```

**Advanced iptables Rules**:
```bash
# Port forwarding (NAT)
iptables -t nat -A PREROUTING -p tcp --dport 8080 -j DNAT --to-destination 192.168.1.10:80
iptables -t nat -A POSTROUTING -p tcp -d 192.168.1.10 --dport 80 -j SNAT --to-source 10.0.0.1

# Load balancing with iptables
iptables -t nat -A PREROUTING -p tcp --dport 80 -m statistic --mode nth --every 3 --packet 0 -j DNAT --to-destination 192.168.1.10:80
iptables -t nat -A PREROUTING -p tcp --dport 80 -m statistic --mode nth --every 2 --packet 0 -j DNAT --to-destination 192.168.1.11:80
iptables -t nat -A PREROUTING -p tcp --dport 80 -j DNAT --to-destination 192.168.1.12:80

# Block specific IP ranges
iptables -A INPUT -s 192.168.100.0/24 -j DROP
iptables -A INPUT -s 10.0.0.0/8 -j DROP

# Allow specific applications
iptables -A OUTPUT -p tcp --dport 443 -m owner --uid-owner www-data -j ACCEPT
iptables -A OUTPUT -p tcp --dport 3306 -m owner --uid-owner mysql -j ACCEPT
```

### UFW (Uncomplicated Firewall)

```bash
# Enable/disable UFW
sudo ufw enable
sudo ufw disable

# Default policies
sudo ufw default deny incoming
sudo ufw default allow outgoing

# Allow services
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow from 192.168.1.0/24 to any port 22

# Deny services
sudo ufw deny 23/tcp
sudo ufw deny from 192.168.100.0/24

# Application profiles
sudo ufw app list
sudo ufw allow 'Nginx Full'
sudo ufw allow 'OpenSSH'

# Status and rules
sudo ufw status verbose
sudo ufw status numbered
sudo ufw delete 2

# Logging
sudo ufw logging on
sudo ufw logging medium
```

### Windows Firewall

```powershell
# Windows Firewall with PowerShell

# Enable/disable firewall
Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled True
Set-NetFirewallProfile -Profile Domain,Public,Private -Enabled False

# View firewall rules
Get-NetFirewallRule | Where-Object {$_.Enabled -eq 'True'}
Get-NetFirewallRule -DisplayName "*HTTP*"

# Create firewall rules
New-NetFirewallRule -DisplayName "Allow HTTP" -Direction Inbound -Protocol TCP -LocalPort 80 -Action Allow
New-NetFirewallRule -DisplayName "Allow HTTPS" -Direction Inbound -Protocol TCP -LocalPort 443 -Action Allow
New-NetFirewallRule -DisplayName "Block Telnet" -Direction Inbound -Protocol TCP -LocalPort 23 -Action Block

# Remove firewall rules
Remove-NetFirewallRule -DisplayName "Allow HTTP"

# Advanced rules
New-NetFirewallRule -DisplayName "Allow SSH from Management" -Direction Inbound -Protocol TCP -LocalPort 22 -RemoteAddress 192.168.1.0/24 -Action Allow

# Application rules
New-NetFirewallRule -DisplayName "Allow App" -Direction Inbound -Program "C:\Program Files\MyApp\app.exe" -Action Allow
```

## Network Monitoring and Performance

### Bandwidth Monitoring

**iftop - Real-time Bandwidth Usage**:
```bash
# Install iftop
sudo apt install iftop        # Ubuntu/Debian
sudo yum install iftop        # CentOS/RHEL

# Usage
sudo iftop                    # Monitor default interface
sudo iftop -i eth0            # Monitor specific interface
sudo iftop -n                 # Don't resolve hostnames
sudo iftop -P                 # Show ports
sudo iftop -B                 # Show bandwidth in bytes
```

**vnstat - Network Statistics**:
```bash
# Install vnstat
sudo apt install vnstat

# Initialize database
sudo vnstat -u -i eth0

# View statistics
vnstat                        # Summary
vnstat -h                     # Hourly stats
vnstat -d                     # Daily stats
vnstat -m                     # Monthly stats
vnstat -w                     # Weekly stats
vnstat -t                     # Top 10 days

# Live monitoring
vnstat -l -i eth0             # Live mode
```

**Network Monitoring Script**:
```bash
#!/bin/bash
# network-monitor.sh

INTERFACE="${1:-eth0}"
INTERVAL="${2:-5}"
LOG_FILE="/var/log/network-monitor.log"

log_network_stats() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local rx_bytes=$(cat /sys/class/net/$INTERFACE/statistics/rx_bytes)
    local tx_bytes=$(cat /sys/class/net/$INTERFACE/statistics/tx_bytes)
    local rx_packets=$(cat /sys/class/net/$INTERFACE/statistics/rx_packets)
    local tx_packets=$(cat /sys/class/net/$INTERFACE/statistics/tx_packets)
    local rx_errors=$(cat /sys/class/net/$INTERFACE/statistics/rx_errors)
    local tx_errors=$(cat /sys/class/net/$INTERFACE/statistics/tx_errors)
    
    echo "$timestamp,$INTERFACE,$rx_bytes,$tx_bytes,$rx_packets,$tx_packets,$rx_errors,$tx_errors" >> "$LOG_FILE"
}

calculate_bandwidth() {
    local prev_rx=$1
    local prev_tx=$2
    local curr_rx=$3
    local curr_tx=$4
    local interval=$5
    
    local rx_rate=$(( (curr_rx - prev_rx) / interval ))
    local tx_rate=$(( (curr_tx - prev_tx) / interval ))
    
    echo "RX: $(numfmt --to=iec-i --suffix=B/s $rx_rate), TX: $(numfmt --to=iec-i --suffix=B/s $tx_rate)"
}

echo "Monitoring interface $INTERFACE every $INTERVAL seconds..."
echo "Timestamp,Interface,RX_Bytes,TX_Bytes,RX_Packets,TX_Packets,RX_Errors,TX_Errors" > "$LOG_FILE"

prev_rx=$(cat /sys/class/net/$INTERFACE/statistics/rx_bytes)
prev_tx=$(cat /sys/class/net/$INTERFACE/statistics/tx_bytes)

while true; do
    sleep $INTERVAL
    
    curr_rx=$(cat /sys/class/net/$INTERFACE/statistics/rx_bytes)
    curr_tx=$(cat /sys/class/net/$INTERFACE/statistics/tx_bytes)
    
    log_network_stats
    
    bandwidth=$(calculate_bandwidth $prev_rx $prev_tx $curr_rx $curr_tx $INTERVAL)
    echo "$(date '+%H:%M:%S') - $bandwidth"
    
    prev_rx=$curr_rx
    prev_tx=$curr_tx
done
```

### Network Performance Testing

**iperf3 - Network Performance Testing**:
```bash
# Install iperf3
sudo apt install iperf3

# Server mode
iperf3 -s                     # Start server
iperf3 -s -p 5002             # Start server on port 5002

# Client mode
iperf3 -c server_ip           # Test to server
iperf3 -c server_ip -t 30     # Test for 30 seconds
iperf3 -c server_ip -P 4      # Use 4 parallel streams
iperf3 -c server_ip -u        # UDP test
iperf3 -c server_ip -R        # Reverse test (server to client)

# Advanced testing
iperf3 -c server_ip -w 1M     # Set window size
iperf3 -c server_ip -M 1400   # Set MSS
iperf3 -c server_ip -J        # JSON output
```

**Network Latency Testing**:
```bash
#!/bin/bash
# latency-test.sh

TARGETS=("8.8.8.8" "1.1.1.1" "google.com" "github.com")
COUNT=10
RESULT_FILE="latency-results-$(date +%Y%m%d-%H%M%S).csv"

echo "Target,Min,Avg,Max,Stddev,Loss%" > "$RESULT_FILE"

for target in "${TARGETS[@]}"; do
    echo "Testing latency to $target..."
    
    result=$(ping -c $COUNT "$target" 2>/dev/null | tail -2)
    
    if [[ $? -eq 0 ]]; then
        # Extract statistics
        stats=$(echo "$result" | grep "min/avg/max/stddev" | cut -d'=' -f2 | tr '/' ',')
        loss=$(echo "$result" | grep "packet loss" | grep -o '[0-9]*%')
        
        echo "$target,$stats,$loss" >> "$RESULT_FILE"
        echo "  Result: $stats (Loss: $loss)"
    else
        echo "$target,,,,,100%" >> "$RESULT_FILE"
        echo "  Result: Failed to reach target"
    fi
done

echo "Results saved to: $RESULT_FILE"
```

## Container Networking

### Docker Networking

**Docker Network Types**:
```bash
# List networks
docker network ls

# Create networks
docker network create mynetwork                    # Bridge network
docker network create --driver host hostnetwork    # Host network
docker network create --driver none nonenetwork    # No network

# Custom bridge network
docker network create --driver bridge \
  --subnet=172.20.0.0/16 \
  --ip-range=172.20.240.0/20 \
  --gateway=172.20.0.1 \
  mybridge

# Connect containers
docker run -d --name web --network mynetwork nginx
docker run -d --name db --network mynetwork mysql

# Inspect network
docker network inspect mynetwork

# Connect/disconnect containers
docker network connect mynetwork container_name
docker network disconnect mynetwork container_name
```

**Docker Compose Networking**:
```yaml
# docker-compose.yml
version: '3.8'

services:
  web:
    image: nginx
    ports:
      - "80:80"
    networks:
      - frontend
      - backend
    depends_on:
      - api
  
  api:
    image: myapi:latest
    ports:
      - "8080:8080"
    networks:
      - backend
    depends_on:
      - database
  
  database:
    image: postgres:13
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    networks:
      - backend
    volumes:
      - db_data:/var/lib/postgresql/data

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: true  # No external access

volumes:
  db_data:
```

### Kubernetes Networking

**Pod Networking**:
```yaml
# pod-with-network-policy.yaml
apiVersion: v1
kind: Pod
metadata:
  name: web-pod
  labels:
    app: web
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: web-netpol
spec:
  podSelector:
    matchLabels:
      app: web
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 80
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: database
    ports:
    - protocol: TCP
      port: 5432
```

**Service Types**:
```yaml
# ClusterIP Service (internal)
apiVersion: v1
kind: Service
metadata:
  name: web-service
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
  type: ClusterIP
---
# NodePort Service (external access)
apiVersion: v1
kind: Service
metadata:
  name: web-nodeport
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
    nodePort: 30080
  type: NodePort
---
# LoadBalancer Service (cloud)
apiVersion: v1
kind: Service
metadata:
  name: web-loadbalancer
spec:
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 80
  type: LoadBalancer
```

## Practical Exercises

### Exercise 1: Network Troubleshooting Lab
1. Set up a multi-tier application with network issues
2. Use various tools to diagnose connectivity problems
3. Implement monitoring and alerting
4. Document troubleshooting procedures

### Exercise 2: Load Balancer Configuration
1. Configure Nginx or HAProxy as a load balancer
2. Set up SSL termination
3. Implement health checks
4. Test failover scenarios

### Exercise 3: Container Networking
1. Create a multi-container application
2. Configure custom networks
3. Implement network policies
4. Monitor network performance

## Key Takeaways

- Networking is fundamental to DevOps infrastructure
- Understanding protocols and tools is essential for troubleshooting
- Load balancing and reverse proxies improve application availability
- Security should be implemented at the network level
- Monitoring and performance testing ensure optimal network health
- Container networking requires special considerations
- Automation and scripting reduce manual network management tasks
- Documentation and procedures are crucial for incident response

## Conclusion

This completes our comprehensive DevOps guide covering all eight essential areas:

1. **Introduction to DevOps** - Foundation and principles
2. **DevOps Culture and Mindset** - Cultural transformation
3. **Version Control Systems** - Code management with Git
4. **Programming Fundamentals** - Essential coding skills
5. **Build Automation** - Automated build processes
6. **Testing Strategies** - Quality assurance practices
7. **Operating Systems** - System administration
8. **Networking Fundamentals** - Network concepts and tools

Each part builds upon the previous ones, creating a solid foundation for DevOps practices. Continue learning by:

- Practicing hands-on exercises
- Building real-world projects
- Staying updated with industry trends
- Contributing to open-source projects
- Pursuing relevant certifications
- Joining DevOps communities

Remember that DevOps is a journey of continuous learning and improvement. The tools and practices evolve rapidly, but the fundamental principles remain constant: collaboration, automation, measurement, and sharing.

---

**End of DevOps Learning Guide**