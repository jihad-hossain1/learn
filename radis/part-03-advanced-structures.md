# Part 3: Advanced Data Structures

This part covers advanced Redis data structures and operations that enable sophisticated use cases like streams, modules, and complex data manipulations.

## Redis Streams

Redis Streams is a data structure that acts like an append-only log, perfect for event sourcing, activity feeds, and real-time data processing.

### Basic Stream Operations

```bash
# Add entries to stream
XADD mystream * sensor-id 1234 temperature 19.8 humidity 65.2
# Output: "1640995200000-0"

XADD mystream * sensor-id 1235 temperature 20.1 humidity 63.8
# Output: "1640995201000-0"

# Add with explicit ID
XADD mystream 1640995202000-0 sensor-id 1236 temperature 18.9

# Read all entries
XRANGE mystream - +
# Output: Shows all entries with their IDs and field-value pairs

# Read entries from specific ID
XRANGE mystream 1640995200000-0 +

# Read latest entries
XREVRANGE mystream + - COUNT 5

# Get stream length
XLEN mystream
# Output: (integer) 3
```

### Stream Consumer Groups

```bash
# Create consumer group
XGROUP CREATE mystream mygroup $ MKSTREAM

# Read from group (blocking)
XREADGROUP GROUP mygroup consumer1 COUNT 1 BLOCK 2000 STREAMS mystream >

# Read pending messages
XREADGROUP GROUP mygroup consumer1 COUNT 1 STREAMS mystream 0

# Acknowledge message processing
XACK mystream mygroup 1640995200000-0

# Get pending messages info
XPENDING mystream mygroup

# Get detailed pending info
XPENDING mystream mygroup - + 10 consumer1

# Claim abandoned messages
XCLAIM mystream mygroup consumer2 3600000 1640995200000-0
```

### Stream Information Commands

```bash
# Get stream info
XINFO STREAM mystream

# Get consumer groups info
XINFO GROUPS mystream

# Get consumers info
XINFO CONSUMERS mystream mygroup

# Trim stream (keep last 1000 entries)
XTRIM mystream MAXLEN 1000

# Trim stream (approximate, more efficient)
XTRIM mystream MAXLEN ~ 1000
```

## Advanced List Operations

### List Blocking Operations

```bash
# Blocking pop with timeout
BLPOP queue1 queue2 queue3 10    # Wait up to 10 seconds
BRPOP queue1 queue2 queue3 0     # Wait indefinitely

# Blocking pop and push
BRPOPLPUSH source destination 10

# Multiple list operations
LMOVE source destination LEFT RIGHT    # Move from left of source to right of destination
BLMOVE source destination LEFT RIGHT 10  # Blocking version
```

### List as Queue Implementation

```bash
# Producer (add to queue)
LPUSH task_queue '{"id":1,"type":"email","data":{"to":"user@example.com"}}'
LPUSH task_queue '{"id":2,"type":"sms","data":{"phone":"+1234567890"}}'

# Consumer (process from queue)
BRPOP task_queue 10
# Process the task...

# Priority queue using multiple lists
BRPOP high_priority medium_priority low_priority 10
```

## Advanced Set Operations

### Set Scanning

```bash
# Scan set members
SSCAN myset 0 MATCH "prefix:*" COUNT 100

# Move member between sets
SMOVE source_set dest_set member

# Random sampling
SRANDMEMBER myset 5        # Get 5 random members (with possible duplicates)
SRANDMEMBER myset -5       # Get 5 random unique members
```

### Set-based Analytics

```bash
# Daily active users
SADD users:active:2023:12:01 user1 user2 user3
SADD users:active:2023:12:02 user2 user3 user4
SADD users:active:2023:12:03 user1 user3 user5

# Users active all three days
SINTER users:active:2023:12:01 users:active:2023:12:02 users:active:2023:12:03

# Users active on any day
SUNION users:active:2023:12:01 users:active:2023:12:02 users:active:2023:12:03

# Users active only on first day
SDIFF users:active:2023:12:01 users:active:2023:12:02 users:active:2023:12:03
```

## Advanced Sorted Set Operations

### Sorted Set Lexicographical Operations

```bash
# Add members with same score for lexicographical ordering
ZADD words 0 "apple" 0 "banana" 0 "cherry" 0 "date"

# Range by lexicographical order
ZRANGEBYLEX words [a [c
# Output: 1) "apple" 2) "banana"

ZRANGEBYLEX words [b +
# Output: 1) "banana" 2) "cherry" 3) "date"

# Count in lexicographical range
ZLEXCOUNT words [a [c
# Output: (integer) 2

# Remove by lexicographical range
ZREMRANGEBYLEX words [a [b
```

### Sorted Set Aggregation

```bash
# Create multiple sorted sets
ZADD math_scores 85 "alice" 92 "bob" 78 "charlie"
ZADD english_scores 88 "alice" 79 "bob" 95 "charlie"

# Union with sum aggregation
ZUNIONSTORE total_scores 2 math_scores english_scores
ZRANGE total_scores 0 -1 WITHSCORES
# Output: 1) "bob" 2) "171" 3) "alice" 4) "173" 5) "charlie" 6) "173"

# Union with max aggregation
ZUNIONSTORE max_scores 2 math_scores english_scores AGGREGATE MAX

# Union with min aggregation
ZUNIONSTORE min_scores 2 math_scores english_scores AGGREGATE MIN

# Weighted union
ZUNIONSTORE weighted_scores 2 math_scores english_scores WEIGHTS 0.6 0.4

# Intersection
ZINTERSTORE common_scores 2 math_scores english_scores
```

### Time-based Sorted Sets

```bash
# Store events with timestamps
ZADD events 1640995200 "user_login:123"
ZADD events 1640995260 "page_view:home"
ZADD events 1640995320 "user_logout:123"

# Get events in time range
ZRANGEBYSCORE events 1640995200 1640995300

# Get events from last hour
ZRANGEBYSCORE events ($(date -d '1 hour ago' +%s) +inf

# Remove old events (older than 24 hours)
ZREMRANGEBYSCORE events -inf ($(date -d '24 hours ago' +%s)
```

## Advanced Hash Operations

### Hash Scanning

```bash
# Scan hash fields
HSCAN user:1000 0 MATCH "pref_*" COUNT 10

# Conditional set
HSETNX user:1000 created_at "2023-12-01"

# Get field length
HSTRLEN user:1000 bio
```

### Hash-based Counters

```bash
# Page view counters
HINCRBY stats:daily:2023:12:01 page:home 1
HINCRBY stats:daily:2023:12:01 page:about 1
HINCRBY stats:daily:2023:12:01 page:contact 1

# User activity counters
HINCRBY user:1000:stats posts 1
HINCRBY user:1000:stats comments 1
HINCRBY user:1000:stats likes 5

# Get all stats
HGETALL user:1000:stats
```

## Advanced String Operations

### String Bit Operations

```bash
# User permissions (bit flags)
SETBIT user:1000:permissions 0 1    # Read permission
SETBIT user:1000:permissions 1 1    # Write permission
SETBIT user:1000:permissions 2 0    # Delete permission
SETBIT user:1000:permissions 3 1    # Admin permission

# Check specific permission
GETBIT user:1000:permissions 1      # Check write permission

# Count permissions
BITCOUNT user:1000:permissions

# Feature flags for users
SETBIT features:new_ui 1000 1        # Enable new UI for user 1000
SETBIT features:beta_features 1000 1 # Enable beta features

# Check if user has feature
GETBIT features:new_ui 1000
```

### String Range Operations

```bash
# Get substring
SET message "Hello, World!"
GETRANGE message 0 4     # "Hello"
GETRANGE message 7 -1    # "World!"

# Set substring
SETRANGE message 7 "Redis"
GET message              # "Hello, Redis"
```

## HyperLogLog Advanced Usage

### Unique Visitor Tracking

```bash
# Track unique visitors per page
PFADD page:home:visitors user1 user2 user3
PFADD page:about:visitors user2 user4 user5
PFADD page:contact:visitors user1 user5 user6

# Get unique visitors per page
PFCOUNT page:home:visitors     # Approximate count
PFCOUNT page:about:visitors
PFCOUNT page:contact:visitors

# Total unique visitors across all pages
PFMERGE total:visitors page:home:visitors page:about:visitors page:contact:visitors
PFCOUNT total:visitors

# Daily unique visitors
PFADD visitors:2023:12:01 user1 user2 user3
PFADD visitors:2023:12:02 user2 user4 user5
PFADD visitors:2023:12:03 user1 user5 user6

# Weekly unique visitors
PFMERGE visitors:week:48 visitors:2023:12:01 visitors:2023:12:02 visitors:2023:12:03
PFCOUNT visitors:week:48
```

## Geospatial Advanced Operations

### Location-based Services

```bash
# Add restaurants with coordinates
GEOADD restaurants -74.0059 40.7128 "Restaurant A"    # New York
GEOADD restaurants -118.2437 34.0522 "Restaurant B"   # Los Angeles
GEOADD restaurants -87.6298 41.8781 "Restaurant C"    # Chicago

# Find restaurants within 1000km of New York
GEORADIUS restaurants -74.0059 40.7128 1000 km WITHDIST WITHCOORD ASC

# Find restaurants near Restaurant A
GEORADIUSBYMEMBER restaurants "Restaurant A" 2000 km WITHDIST

# Get hash of location
GEOHASH restaurants "Restaurant A"
# Output: 1) "dr5regw3p6g0"

# Store search results
GEORADIUS restaurants -74.0059 40.7128 1000 km STORE nearby_restaurants
GEORADIUS restaurants -74.0059 40.7128 1000 km STOREDIST nearby_with_dist
```

### Geofencing

```bash
# Define geofence areas
GEOADD geofences -74.0059 40.7128 "downtown_ny"
GEOADD geofences -118.2437 34.0522 "downtown_la"

# Check if user is in geofence
GEORADIUSBYMEMBER geofences "downtown_ny" 5 km WITHCOORD

# Track user locations
GEOADD user_locations:1000 -74.0060 40.7129 "current"

# Check if user is near any geofence
GEORADIUS geofences -74.0060 40.7129 10 km
```

## Scanning and Iteration

### SCAN Command Family

```bash
# Scan all keys
SCAN 0 MATCH "user:*" COUNT 100

# Scan hash fields
HSCAN user:1000 0 MATCH "pref_*" COUNT 10

# Scan set members
SSCAN tags 0 MATCH "tech:*" COUNT 50

# Scan sorted set members
ZSCAN leaderboard 0 MATCH "player:*" COUNT 20
```

### Safe Key Iteration

```bash
# Instead of KEYS * (dangerous in production)
# Use SCAN for safe iteration
redis-cli --scan --pattern "user:*"

# Count keys matching pattern
redis-cli --scan --pattern "session:*" | wc -l

# Delete keys matching pattern safely
redis-cli --scan --pattern "temp:*" | xargs redis-cli DEL
```

## Memory Optimization Techniques

### Small Hash Optimization

```bash
# Configure small hash optimization
# In redis.conf:
# hash-max-ziplist-entries 512
# hash-max-ziplist-value 64

# Use hashes for small objects
HMSET user:1000 name "John" age "30" city "NYC"  # More memory efficient
# Instead of:
# SET user:1000:name "John"
# SET user:1000:age "30"
# SET user:1000:city "NYC"
```

### Integer Set Optimization

```bash
# Small integer sets are stored efficiently
SADD small_numbers 1 2 3 4 5 6 7 8 9 10

# Configure in redis.conf:
# set-max-intset-entries 512
```

### Ziplist Optimization

```bash
# Small lists use ziplist encoding
LPUSH small_list "item1" "item2" "item3"

# Configure in redis.conf:
# list-max-ziplist-size -2
# list-compress-depth 0
```

## Advanced Expiration Patterns

### Sliding Window Expiration

```bash
# Rate limiting with sliding window
SET rate_limit:user:1000:$(date +%s) 1 EX 3600
INCR rate_limit:user:1000:$(date +%s)

# Check rate limit
KEYS "rate_limit:user:1000:*" | wc -l
```

### Lazy Expiration

```bash
# Set expiration on access
SET cache:data "value"
# Later, when accessed:
EXPIRE cache:data 3600
```

### Conditional Expiration

```bash
# Extend expiration if key exists
SET session:abc123 "data" EX 1800
# Later:
EXPIRE session:abc123 1800  # Extend if still exists
```

## Data Structure Combinations

### User Session Management

```bash
# Session data in hash
HMSET session:abc123 user_id 1000 login_time 1640995200 last_activity 1640995800

# Active sessions in set
SADD active_sessions session:abc123

# Session expiration
EXPIRE session:abc123 1800

# Cleanup expired sessions from set
# (Use Lua script or application logic)
```

### Shopping Cart System

```bash
# Cart items in hash
HMSET cart:user:1000 item:123 2 item:456 1 item:789 3

# Cart metadata
HMSET cart:user:1000:meta created 1640995200 updated 1640995800 total 6

# User's active cart
SET user:1000:active_cart cart:user:1000

# Cart expiration
EXPIRE cart:user:1000 86400  # 24 hours
EXPIRE cart:user:1000:meta 86400
```

### Real-time Analytics

```bash
# Page views by minute
ZADD analytics:page_views:2023:12:01 1640995200 "page:home"
ZADD analytics:page_views:2023:12:01 1640995260 "page:about"

# Unique visitors
PFADD analytics:unique_visitors:2023:12:01 user1 user2 user3

# Top pages
ZINCRBY analytics:top_pages 1 "page:home"
ZINCRBY analytics:top_pages 1 "page:about"

# Get top 10 pages
ZREVRANGE analytics:top_pages 0 9 WITHSCORES
```

## Performance Monitoring

### Memory Usage Analysis

```bash
# Check memory usage by data type
INFO memory

# Find memory usage of specific keys
MEMORY USAGE user:1000
MEMORY USAGE session:abc123

# Find largest keys
redis-cli --bigkeys

# Sample memory usage
redis-cli --memkeys --memkeys-samples 1000
```

### Operation Monitoring

```bash
# Monitor commands in real-time
MONITOR

# Check slow queries
SLOWLOG GET 10

# Configure slow log
CONFIG SET slowlog-log-slower-than 10000  # 10ms
CONFIG SET slowlog-max-len 128
```

## Next Steps

Now that you understand advanced Redis data structures and operations, you're ready to integrate Redis with Node.js:

- [Part 4: Node.js Redis Client Setup](./part-04-nodejs-setup.md)

## Quick Reference

### Stream Commands
```bash
XADD stream * field value     # Add entry
XRANGE stream - +            # Read all entries
XGROUP CREATE stream group $ # Create consumer group
XREADGROUP GROUP group consumer COUNT 1 STREAMS stream >
XACK stream group id         # Acknowledge processing
```

### Advanced Set Operations
```bash
SSCAN set 0 MATCH pattern   # Scan set members
SMOVE source dest member     # Move between sets
SRANDMEMBER set count       # Random sampling
```

### Geospatial Commands
```bash
GEOADD key lng lat member    # Add location
GEORADIUS key lng lat radius # Find nearby
GEODIST key m1 m2           # Calculate distance
GEOHASH key member          # Get geohash
```

### Scanning Commands
```bash
SCAN cursor MATCH pattern   # Scan keys
HSCAN key cursor            # Scan hash fields
SSCAN key cursor            # Scan set members
ZSCAN key cursor            # Scan sorted set
```

---

**Previous**: [← Part 2: Data Types and Basic Operations](./part-02-data-types.md)  
**Next**: [Part 4: Node.js Redis Client Setup →](./part-04-nodejs-setup.md)