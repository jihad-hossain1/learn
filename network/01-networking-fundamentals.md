# Networking Fundamentals for Full Stack Development

## Table of Contents
1. [What is Networking?](#what-is-networking)
2. [OSI Model](#osi-model)
3. [TCP/IP Model](#tcpip-model)
4. [IP Addresses](#ip-addresses)
5. [Ports and Protocols](#ports-and-protocols)
6. [DNS (Domain Name System)](#dns-domain-name-system)
7. [HTTP/HTTPS Basics](#httphttps-basics)

## What is Networking?

Networking is the practice of connecting computers and other devices to share resources and communicate. In full stack development, understanding networking is crucial for:

- Building web applications that communicate over the internet
- Designing APIs and microservices
- Implementing security measures
- Optimizing performance
- Troubleshooting connectivity issues

## OSI Model

The Open Systems Interconnection (OSI) model is a conceptual framework that describes how network protocols interact in seven layers:

### Layer 7: Application Layer
- **Purpose**: User interface and network services
- **Examples**: HTTP, HTTPS, FTP, SMTP, DNS
- **Full Stack Relevance**: Where your web applications operate

### Layer 6: Presentation Layer
- **Purpose**: Data encryption, compression, and translation
- **Examples**: SSL/TLS, JPEG, GIF
- **Full Stack Relevance**: Data formatting and encryption

### Layer 5: Session Layer
- **Purpose**: Establishes, manages, and terminates connections
- **Examples**: NetBIOS, RPC
- **Full Stack Relevance**: Session management in web apps

### Layer 4: Transport Layer
- **Purpose**: Reliable data transfer
- **Examples**: TCP, UDP
- **Full Stack Relevance**: How data is transmitted reliably

### Layer 3: Network Layer
- **Purpose**: Routing and logical addressing
- **Examples**: IP, ICMP, OSPF
- **Full Stack Relevance**: How data finds its way across networks

### Layer 2: Data Link Layer
- **Purpose**: Physical addressing and error detection
- **Examples**: Ethernet, Wi-Fi
- **Full Stack Relevance**: Local network communication

### Layer 1: Physical Layer
- **Purpose**: Physical transmission of data
- **Examples**: Cables, radio waves, fiber optics
- **Full Stack Relevance**: The physical infrastructure

## TCP/IP Model

The TCP/IP model is a more practical, four-layer model used in real-world networking:

### 1. Application Layer
- Combines OSI layers 5, 6, and 7
- Where your applications run (HTTP, HTTPS, FTP, SMTP)

### 2. Transport Layer
- Same as OSI Layer 4
- TCP (reliable) vs UDP (fast but unreliable)

### 3. Internet Layer
- Same as OSI Layer 3
- IP addressing and routing

### 4. Network Access Layer
- Combines OSI layers 1 and 2
- Physical network access

## IP Addresses

### IPv4
- **Format**: 32-bit addresses (e.g., 192.168.1.1)
- **Range**: 0.0.0.0 to 255.255.255.255
- **Classes**:
  - Class A: 1.0.0.0 to 126.255.255.255 (large networks)
  - Class B: 128.0.0.0 to 191.255.255.255 (medium networks)
  - Class C: 192.0.0.0 to 223.255.255.255 (small networks)

### IPv6
- **Format**: 128-bit addresses (e.g., 2001:0db8:85a3:0000:0000:8a2e:0370:7334)
- **Purpose**: Solve IPv4 address exhaustion
- **Benefits**: More addresses, better security, improved performance

### Private vs Public IP Addresses

**Private IP Ranges** (not routable on the internet):
- 10.0.0.0 to 10.255.255.255
- 172.16.0.0 to 172.31.255.255
- 192.168.0.0 to 192.168.255.255

**Public IP Addresses**: Routable on the internet, assigned by ISPs

## Ports and Protocols

### Common Ports for Full Stack Development

| Port | Protocol | Service | Usage |
|------|----------|---------|-------|
| 80 | HTTP | Web traffic | Standard web servers |
| 443 | HTTPS | Secure web traffic | SSL/TLS encrypted web |
| 22 | SSH | Secure Shell | Remote server access |
| 21 | FTP | File Transfer | File uploads/downloads |
| 25 | SMTP | Email sending | Mail servers |
| 53 | DNS | Domain resolution | Name to IP translation |
| 3000 | Custom | Development | React dev server |
| 3306 | MySQL | Database | MySQL connections |
| 5432 | PostgreSQL | Database | PostgreSQL connections |
| 27017 | MongoDB | Database | MongoDB connections |
| 6379 | Redis | Cache/Database | Redis connections |

### TCP vs UDP

**TCP (Transmission Control Protocol)**:
- Reliable, connection-oriented
- Guarantees delivery and order
- Error checking and correction
- Used for: HTTP, HTTPS, FTP, email

**UDP (User Datagram Protocol)**:
- Fast, connectionless
- No delivery guarantee
- No error correction
- Used for: DNS, video streaming, gaming

## DNS (Domain Name System)

DNS translates human-readable domain names to IP addresses.

### DNS Hierarchy
```
. (root)
├── com
│   ├── google
│   │   └── www
│   └── github
│       └── api
├── org
│   └── wikipedia
│       └── en
└── net
    └── cloudflare
```

### DNS Record Types

| Type | Purpose | Example |
|------|---------|----------|
| A | IPv4 address | example.com → 192.168.1.1 |
| AAAA | IPv6 address | example.com → 2001:db8::1 |
| CNAME | Alias | www.example.com → example.com |
| MX | Mail server | example.com → mail.example.com |
| TXT | Text data | SPF, DKIM records |
| NS | Name server | example.com → ns1.provider.com |

### DNS Resolution Process
1. Browser checks local cache
2. Query local DNS resolver
3. Query root name servers
4. Query TLD name servers
5. Query authoritative name servers
6. Return IP address to browser

## HTTP/HTTPS Basics

### HTTP (HyperText Transfer Protocol)
- **Port**: 80
- **Security**: None (plain text)
- **Use case**: Development, non-sensitive data

### HTTPS (HTTP Secure)
- **Port**: 443
- **Security**: SSL/TLS encryption
- **Use case**: Production, sensitive data

### HTTP Methods

| Method | Purpose | Idempotent | Safe |
|--------|---------|------------|------|
| GET | Retrieve data | Yes | Yes |
| POST | Create/submit data | No | No |
| PUT | Update/replace data | Yes | No |
| PATCH | Partial update | No | No |
| DELETE | Remove data | Yes | No |
| HEAD | Get headers only | Yes | Yes |
| OPTIONS | Get allowed methods | Yes | Yes |

### HTTP Status Codes

**1xx - Informational**
- 100 Continue
- 101 Switching Protocols

**2xx - Success**
- 200 OK
- 201 Created
- 204 No Content

**3xx - Redirection**
- 301 Moved Permanently
- 302 Found (Temporary Redirect)
- 304 Not Modified

**4xx - Client Error**
- 400 Bad Request
- 401 Unauthorized
- 403 Forbidden
- 404 Not Found
- 429 Too Many Requests

**5xx - Server Error**
- 500 Internal Server Error
- 502 Bad Gateway
- 503 Service Unavailable
- 504 Gateway Timeout

## Key Takeaways

1. **Understand the layers**: Each layer has a specific purpose
2. **Know your protocols**: HTTP/HTTPS for web, TCP for reliability, UDP for speed
3. **Master IP addressing**: Public vs private, IPv4 vs IPv6
4. **DNS is crucial**: Domain resolution affects performance
5. **Ports matter**: Know common ports for development
6. **HTTP is the foundation**: Status codes and methods are essential

## Next Steps

In the next part, we'll dive deeper into:
- Advanced HTTP concepts
- RESTful API design
- WebSockets and real-time communication
- Network security fundamentals

---

*Continue to [Part 2: HTTP Deep Dive and APIs](02-http-apis.md)*