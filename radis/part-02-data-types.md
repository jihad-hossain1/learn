# Part 2: Data Types and Basic Operations

Redis supports several data types, each optimized for specific use cases. Understanding these data types is crucial for effective Redis usage.

## String Data Type

Strings are the most basic Redis data type and can store text, numbers, or binary data up to 512MB.

### Basic String Operations

```bash
# Set and get strings
SET name "John Doe"
GET name
# Output: "John Doe"

# Set multiple keys at once
MSET key1 "value1" key2 "value2" key3 "value3"
MGET key1 key2 key3
# Output: 1) "value1" 2) "value2" 3) "value3"

# Set with expiration
SETEX session:user123 3600 "active"
# or
SET session:user123 "active" EX 3600

# Set only if key doesn't exist
SETNX counter 0
# Output: (integer) 1 (success)

# Append to existing string
SET message "Hello"
APPEND message " World"
GET message
# Output: "Hello World"

# Get string length
STRLEN message
# Output: (integer) 11
```

### Numeric Operations

```bash
# Increment/Decrement
SET counter 10
INCR counter        # counter = 11
DECR counter        # counter = 10
INCRBY counter 5    # counter = 15
DECRBY counter 3    # counter = 12

# Floating point increment
SET price 19.99
INCRBYFLOAT price 5.01
GET price
# Output: "25"
```

### Bit Operations

```bash
# Set bit at position
SETBIT user:1000:login 0 1  # Set bit 0 to 1
SETBIT user:1000:login 1 0  # Set bit 1 to 0

# Get bit at position
GETBIT user:1000:login 0
# Output: (integer) 1

# Count set bits
BITCOUNT user:1000:login
# Output: (integer) 1
```

## Hash Data Type

Hashes are maps between string fields and string values, perfect for representing objects.

### Basic Hash Operations

```bash
# Set hash fields
HSET user:1000 name "John Doe" email "john@example.com" age 30

# Get single field
HGET user:1000 name
# Output: "John Doe"

# Get multiple fields
HMGET user:1000 name email age
# Output: 1) "John Doe" 2) "john@example.com" 3) "30"

# Get all fields and values
HGETALL user:1000
# Output: 1) "name" 2) "John Doe" 3) "email" 4) "john@example.com" 5) "age" 6) "30"

# Check if field exists
HEXISTS user:1000 name
# Output: (integer) 1

# Delete field
HDEL user:1000 age

# Get all field names
HKEYS user:1000
# Output: 1) "name" 2) "email"

# Get all values
HVALS user:1000
# Output: 1) "John Doe" 2) "john@example.com"

# Get number of fields
HLEN user:1000
# Output: (integer) 2
```

### Hash Numeric Operations

```bash
# Increment hash field
HSET stats:user:1000 page_views 100
HINCRBY stats:user:1000 page_views 1
HGET stats:user:1000 page_views
# Output: "101"

# Increment by float
HSET stats:user:1000 score 95.5
HINCRBYFLOAT stats:user:1000 score 2.3
HGET stats:user:1000 score
# Output: "97.8"
```

## List Data Type

Lists are ordered collections of strings, supporting operations at both ends.

### Basic List Operations

```bash
# Push elements to list
LPUSH mylist "first"     # Add to left (beginning)
RPUSH mylist "last"      # Add to right (end)
LPUSH mylist "new_first"

# Get list elements
LRANGE mylist 0 -1
# Output: 1) "new_first" 2) "first" 3) "last"

# Get specific range
LRANGE mylist 0 1
# Output: 1) "new_first" 2) "first"

# Get list length
LLEN mylist
# Output: (integer) 3

# Get element at index
LINDEX mylist 0
# Output: "new_first"

# Set element at index
LSET mylist 1 "modified"
LRANGE mylist 0 -1
# Output: 1) "new_first" 2) "modified" 3) "last"
```

### List Pop Operations

```bash
# Pop elements
LPOP mylist    # Remove and return first element
RPOP mylist    # Remove and return last element

# Blocking pop (wait for element)
BLPOP mylist 10    # Wait up to 10 seconds
BRPOP mylist 10    # Wait up to 10 seconds

# Pop from one list and push to another
RPOPLPUSH source destination
```

### List Utility Operations

```bash
# Insert element
LINSERT mylist BEFORE "first" "before_first"
LINSERT mylist AFTER "first" "after_first"

# Remove elements
LREM mylist 2 "value"    # Remove first 2 occurrences of "value"
LREM mylist -1 "value"   # Remove last occurrence of "value"
LREM mylist 0 "value"    # Remove all occurrences of "value"

# Trim list
LTRIM mylist 0 2    # Keep only elements from index 0 to 2
```

## Set Data Type

Sets are unordered collections of unique strings.

### Basic Set Operations

```bash
# Add members to set
SADD myset "apple" "banana" "orange"
SADD myset "apple"    # Won't add duplicate

# Get all members
SMEMBERS myset
# Output: 1) "apple" 2) "banana" 3) "orange"

# Check membership
SISMEMBER myset "apple"
# Output: (integer) 1

# Get set size
SCARD myset
# Output: (integer) 3

# Remove member
SREM myset "banana"

# Get random member
SRANDMEMBER myset
# Output: "apple" (or "orange")

# Get multiple random members
SRANDMEMBER myset 2
# Output: 1) "apple" 2) "orange"

# Pop random member
SPOP myset
# Output: "orange" (removes from set)
```

### Set Operations (Union, Intersection, Difference)

```bash
# Create sets
SADD set1 "a" "b" "c"
SADD set2 "b" "c" "d"
SADD set3 "c" "d" "e"

# Union
SUNION set1 set2
# Output: 1) "a" 2) "b" 3) "c" 4) "d"

# Store union result
SUNIONSTORE result_set set1 set2

# Intersection
SINTER set1 set2
# Output: 1) "b" 2) "c"

# Store intersection result
SINTERSTORE result_set set1 set2

# Difference
SDIFF set1 set2
# Output: 1) "a"

# Store difference result
SDIFFSTORE result_set set1 set2
```

## Sorted Set (ZSet) Data Type

Sorted sets are like sets but each member has an associated score for ordering.

### Basic Sorted Set Operations

```bash
# Add members with scores
ZADD leaderboard 100 "player1" 85 "player2" 92 "player3"

# Get members by rank (lowest to highest score)
ZRANGE leaderboard 0 -1
# Output: 1) "player2" 2) "player3" 3) "player1"

# Get members with scores
ZRANGE leaderboard 0 -1 WITHSCORES
# Output: 1) "player2" 2) "85" 3) "player3" 4) "92" 5) "player1" 6) "100"

# Get members by rank (highest to lowest score)
ZREVRANGE leaderboard 0 -1 WITHSCORES
# Output: 1) "player1" 2) "100" 3) "player3" 4) "92" 5) "player2" 6) "85"

# Get member score
ZSCORE leaderboard "player1"
# Output: "100"

# Get member rank
ZRANK leaderboard "player2"     # 0-based rank (lowest score = 0)
# Output: (integer) 0

ZREVRANK leaderboard "player2"  # 0-based rank (highest score = 0)
# Output: (integer) 2
```

### Sorted Set Range Operations

```bash
# Get members by score range
ZRANGEBYSCORE leaderboard 80 95
# Output: 1) "player2" 2) "player3"

# Get members by score range with scores
ZRANGEBYSCORE leaderboard 80 95 WITHSCORES

# Get members by score range (reverse)
ZREVRANGEBYSCORE leaderboard 95 80

# Count members in score range
ZCOUNT leaderboard 80 95
# Output: (integer) 2

# Get set size
ZCARD leaderboard
# Output: (integer) 3
```

### Sorted Set Modification Operations

```bash
# Increment member score
ZINCRBY leaderboard 5 "player2"
ZSCORE leaderboard "player2"
# Output: "90"

# Remove member
ZREM leaderboard "player3"

# Remove members by rank
ZREMRANGEBYRANK leaderboard 0 0    # Remove lowest ranked member

# Remove members by score
ZREMRANGEBYSCORE leaderboard 0 50  # Remove members with score 0-50
```

## Bitmap Operations

Bitmaps are not a separate data type but string operations that treat strings as bit arrays.

### Basic Bitmap Operations

```bash
# Set bits for user activity tracking
SETBIT user:1000:2023:12:01 0 1    # User active at hour 0
SETBIT user:1000:2023:12:01 8 1    # User active at hour 8
SETBIT user:1000:2023:12:01 16 1   # User active at hour 16

# Get bit value
GETBIT user:1000:2023:12:01 8
# Output: (integer) 1

# Count set bits
BITCOUNT user:1000:2023:12:01
# Output: (integer) 3

# Find first set/unset bit
BITPOS user:1000:2023:12:01 1      # First set bit
# Output: (integer) 0

BITPOS user:1000:2023:12:01 0      # First unset bit
# Output: (integer) 1
```

### Bitmap Operations Between Keys

```bash
# Create bitmaps for different days
SETBIT users:active:2023:12:01 100 1
SETBIT users:active:2023:12:01 200 1
SETBIT users:active:2023:12:02 150 1
SETBIT users:active:2023:12:02 200 1

# Bitwise AND (users active both days)
BITOP AND users:active:both users:active:2023:12:01 users:active:2023:12:02
BITCOUNT users:active:both
# Output: (integer) 1

# Bitwise OR (users active either day)
BITOP OR users:active:either users:active:2023:12:01 users:active:2023:12:02
BITCOUNT users:active:either
# Output: (integer) 3

# Bitwise XOR (users active only one day)
BITOP XOR users:active:exclusive users:active:2023:12:01 users:active:2023:12:02
BITCOUNT users:active:exclusive
# Output: (integer) 2
```

## HyperLogLog

HyperLogLog is a probabilistic data structure for counting unique elements.

### Basic HyperLogLog Operations

```bash
# Add elements
PFADD unique_visitors "user1" "user2" "user3"
PFADD unique_visitors "user1" "user4"    # user1 won't be counted twice

# Get approximate count
PFCOUNT unique_visitors
# Output: (integer) 4

# Merge HyperLogLogs
PFADD visitors:page1 "user1" "user2"
PFADD visitors:page2 "user2" "user3"
PFMERGE visitors:total visitors:page1 visitors:page2
PFCOUNT visitors:total
# Output: (integer) 3
```

## Geospatial Data

Redis supports geospatial operations for location-based applications.

### Basic Geospatial Operations

```bash
# Add locations
GEOADD locations 13.361389 38.115556 "Palermo" 15.087269 37.502669 "Catania"

# Get coordinates
GEOPOS locations "Palermo"
# Output: 1) 1) "13.36138933897018433" 2) "38.11555639549629859"

# Calculate distance
GEODIST locations "Palermo" "Catania" km
# Output: "166.2742"

# Find nearby locations
GEORADIUS locations 15 37 200 km WITHDIST
# Output: 1) 1) "Catania" 2) "56.4413" 2) 1) "Palermo" 2) "190.4424"

# Find nearby locations by member
GEORADIUSBYMEMBER locations "Palermo" 200 km
# Output: 1) "Palermo" 2) "Catania"
```

## Data Type Selection Guide

| Use Case | Recommended Data Type | Example |
|----------|----------------------|----------|
| Simple key-value | String | `SET user:1000:name "John"` |
| Counters | String (numeric) | `INCR page_views` |
| User profiles | Hash | `HSET user:1000 name "John" age 30` |
| Shopping cart | Hash | `HSET cart:user123 item1 2 item2 1` |
| Activity feed | List | `LPUSH feed:user123 "new post"` |
| Message queue | List | `BRPOP queue 10` |
| Tags/categories | Set | `SADD post:123:tags "redis" "database"` |
| Unique visitors | Set | `SADD visitors:today "user123"` |
| Leaderboard | Sorted Set | `ZADD scores 100 "player1"` |
| Time series | Sorted Set | `ZADD metrics timestamp value` |
| Feature flags | Bitmap | `SETBIT features:user123 5 1` |
| Daily active users | Bitmap | `SETBIT dau:2023:12:01 user_id 1` |
| Unique page views | HyperLogLog | `PFADD page:123:visitors "user456"` |
| Location services | Geospatial | `GEOADD stores lng lat "store1"` |

## Performance Tips

### Memory Optimization
```bash
# Use appropriate data types
# Hash is more memory efficient for objects
HSET user:1000 name "John" age 30    # Better than multiple strings

# Use shorter key names in production
SET u:1000:n "John"    # Instead of user:1000:name

# Set expiration to prevent memory leaks
SETEX session:abc123 3600 "user_data"
```

### Operation Efficiency
```bash
# Use pipeline for multiple operations
# Use MSET instead of multiple SET commands
MSET key1 val1 key2 val2 key3 val3

# Use HMSET for multiple hash fields
HMSET user:1000 name "John" email "john@example.com" age 30

# Avoid KEYS command in production (use SCAN instead)
SCAN 0 MATCH user:* COUNT 100
```

## Common Patterns and Best Practices

### Key Naming Conventions
```bash
# Use consistent hierarchy
user:1000:profile
user:1000:preferences
user:1000:sessions:abc123

# Use prefixes for different environments
dev:user:1000:profile
prod:user:1000:profile

# Include data type in key name
counter:page_views
list:recent_posts
set:user_tags
zset:leaderboard
```

### Expiration Strategies
```bash
# Session data
SETEX session:abc123 1800 "session_data"    # 30 minutes

# Cache data
SETEX cache:user:1000 3600 "cached_data"    # 1 hour

# Temporary data
SETEX temp:upload:xyz789 300 "temp_data"    # 5 minutes
```

## Next Steps

Now that you understand Redis data types and basic operations, you're ready to explore:

- [Part 3: Advanced Data Structures](./part-03-advanced-structures.md)

## Quick Reference

### String Commands
```bash
SET key value          # Set string
GET key               # Get string
MSET k1 v1 k2 v2      # Set multiple
MGET k1 k2            # Get multiple
INCR key              # Increment number
APPEND key value      # Append to string
```

### Hash Commands
```bash
HSET key field value  # Set hash field
HGET key field        # Get hash field
HGETALL key          # Get all fields
HDEL key field       # Delete field
HINCRBY key field n  # Increment field
```

### List Commands
```bash
LPUSH key value      # Push to left
RPUSH key value      # Push to right
LPOP key             # Pop from left
RPOP key             # Pop from right
LRANGE key start end # Get range
LLEN key             # Get length
```

### Set Commands
```bash
SADD key member      # Add to set
SMEMBERS key         # Get all members
SISMEMBER key member # Check membership
SREM key member      # Remove member
SUNION key1 key2     # Union sets
SINTER key1 key2     # Intersect sets
```

### Sorted Set Commands
```bash
ZADD key score member    # Add with score
ZRANGE key start end     # Get by rank
ZRANGEBYSCORE key min max # Get by score
ZSCORE key member        # Get score
ZRANK key member         # Get rank
```

---

**Previous**: [← Part 1: Redis Basics and Installation](./part-01-basics-and-setup.md)  
**Next**: [Part 3: Advanced Data Structures →](./part-03-advanced-structures.md)