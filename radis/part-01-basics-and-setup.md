# Part 1: Redis Basics and Installation

## What is Redis?

Redis (Remote Dictionary Server) is an open-source, in-memory data structure store that can be used as a database, cache, and message broker. It supports various data structures such as strings, hashes, lists, sets, and more.

### Key Features

- **In-Memory Storage**: Extremely fast read/write operations
- **Persistence**: Optional data persistence to disk
- **Data Structures**: Rich set of data types
- **Atomic Operations**: All operations are atomic
- **Pub/Sub**: Built-in publish/subscribe messaging
- **Lua Scripting**: Server-side scripting support
- **Clustering**: Horizontal scaling capabilities

### Use Cases

- **Caching**: Session storage, page caching, API response caching
- **Real-time Analytics**: Counters, metrics, leaderboards
- **Message Queues**: Task queues, pub/sub systems
- **Session Store**: User session management
- **Rate Limiting**: API rate limiting and throttling

## Installation

### Windows Installation

#### Option 1: Using WSL (Recommended)
```bash
# Install WSL and Ubuntu
wsl --install

# Update package list
sudo apt update

# Install Redis
sudo apt install redis-server

# Start Redis service
sudo service redis-server start

# Test installation
redis-cli ping
```

#### Option 2: Using Docker
```bash
# Pull Redis image
docker pull redis:latest

# Run Redis container
docker run --name redis-server -p 6379:6379 -d redis:latest

# Connect to Redis CLI
docker exec -it redis-server redis-cli
```

#### Option 3: Windows Binary (Not Recommended for Production)
```bash
# Download from: https://github.com/microsoftarchive/redis/releases
# Extract and run redis-server.exe
```

### macOS Installation
```bash
# Using Homebrew
brew install redis

# Start Redis service
brew services start redis

# Test installation
redis-cli ping
```

### Linux Installation
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install redis-server

# CentOS/RHEL
sudo yum install redis
# or
sudo dnf install redis

# Start Redis service
sudo systemctl start redis
sudo systemctl enable redis
```

## Basic Redis CLI Commands

### Starting Redis
```bash
# Start Redis server
redis-server

# Start with custom configuration
redis-server /path/to/redis.conf

# Start Redis CLI
redis-cli

# Connect to specific host/port
redis-cli -h localhost -p 6379
```

### Basic Operations
```bash
# Test connection
PING
# Response: PONG

# Set a key-value pair
SET mykey "Hello Redis"
# Response: OK

# Get a value
GET mykey
# Response: "Hello Redis"

# Check if key exists
EXISTS mykey
# Response: (integer) 1

# Delete a key
DEL mykey
# Response: (integer) 1

# List all keys
KEYS *

# Get information about Redis server
INFO

# Clear all data
FLUSHALL
```

## Redis Configuration

### Basic Configuration File (redis.conf)
```conf
# Network
bind 127.0.0.1
port 6379

# General
daemonize yes
pidfile /var/run/redis/redis-server.pid
loglevel notice
logfile /var/log/redis/redis-server.log

# Memory
maxmemory 256mb
maxmemory-policy allkeys-lru

# Persistence
save 900 1
save 300 10
save 60 10000

# Security
requirepass your_password_here
```

### Key Configuration Options

- **bind**: IP addresses to bind to
- **port**: Port number (default: 6379)
- **maxmemory**: Maximum memory usage
- **maxmemory-policy**: Eviction policy when memory limit is reached
- **save**: Automatic saving intervals
- **requirepass**: Password authentication

## Redis Data Persistence

### RDB (Redis Database Backup)
- Point-in-time snapshots
- Compact binary format
- Good for backups and disaster recovery
- Faster restart times

### AOF (Append Only File)
- Logs every write operation
- Better durability
- Larger file sizes
- Slower restart times

### Configuration Examples
```conf
# RDB Configuration
save 900 1      # Save if at least 1 key changed in 900 seconds
save 300 10     # Save if at least 10 keys changed in 300 seconds
save 60 10000   # Save if at least 10000 keys changed in 60 seconds

# AOF Configuration
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec  # Options: always, everysec, no
```

## Security Basics

### Authentication
```bash
# Set password in redis.conf
requirepass your_strong_password

# Authenticate in CLI
AUTH your_strong_password

# Or connect with password
redis-cli -a your_strong_password
```

### Network Security
```conf
# Bind to specific interfaces only
bind 127.0.0.1 192.168.1.100

# Disable dangerous commands
rename-command FLUSHDB ""
rename-command FLUSHALL ""
rename-command KEYS ""
rename-command CONFIG "CONFIG_9a90f2b4c8d7e1f3"
```

## Performance Considerations

### Memory Usage
```bash
# Check memory usage
INFO memory

# Get memory usage of specific key
MEMORY USAGE mykey

# Monitor commands in real-time
MONITOR
```

### Benchmarking
```bash
# Basic benchmark
redis-benchmark

# Custom benchmark
redis-benchmark -h localhost -p 6379 -n 100000 -c 50

# Test specific operations
redis-benchmark -t set,get -n 100000 -q
```

## Common Redis Patterns

### Key Naming Conventions
```bash
# Use colons for hierarchy
user:1000:profile
user:1000:sessions
product:electronics:laptop:123

# Use prefixes for different data types
cache:user:1000
session:abc123def456
counter:page_views:2023:12:01
```

### Expiration Patterns
```bash
# Set key with expiration (seconds)
SETEX mykey 3600 "value"

# Set expiration on existing key
EXPIRE mykey 3600

# Check time to live
TTL mykey

# Remove expiration
PERSIST mykey
```

## Troubleshooting Common Issues

### Connection Issues
```bash
# Check if Redis is running
ps aux | grep redis

# Check port availability
netstat -tlnp | grep 6379

# Test connection
telnet localhost 6379
```

### Memory Issues
```bash
# Check memory usage
redis-cli INFO memory

# Find large keys
redis-cli --bigkeys

# Analyze memory usage
redis-cli --memkeys
```

### Performance Issues
```bash
# Monitor slow queries
redis-cli CONFIG SET slowlog-log-slower-than 10000
redis-cli SLOWLOG GET 10

# Check latency
redis-cli --latency

# Monitor in real-time
redis-cli --stat
```

## Next Steps

Now that you have Redis installed and understand the basics, you're ready to move on to:

- [Part 2: Data Types and Basic Operations](./part-02-data-types.md)

## Quick Reference

### Essential Commands
```bash
PING                    # Test connection
SET key value          # Set string value
GET key                # Get string value
DEL key                # Delete key
EXISTS key             # Check if key exists
KEYS pattern           # Find keys matching pattern
TTL key                # Get time to live
EXPIRE key seconds     # Set expiration
INFO                   # Server information
FLUSHALL               # Clear all data
```

### Configuration Files
- **Linux**: `/etc/redis/redis.conf`
- **macOS**: `/usr/local/etc/redis.conf`
- **Docker**: Custom mount or environment variables

---

**Next**: [Part 2: Data Types and Basic Operations →](./part-02-data-types.md)