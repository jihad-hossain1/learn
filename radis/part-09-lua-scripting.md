# Part 9: Lua Scripting in Redis

Lua scripting in Redis allows you to execute complex operations atomically on the server side, reducing network round trips and ensuring data consistency.

## Table of Contents

1. [Introduction to Lua Scripting](#introduction-to-lua-scripting)
2. [Basic Lua Syntax](#basic-lua-syntax)
3. [Redis Lua API](#redis-lua-api)
4. [EVAL and EVALSHA Commands](#eval-and-evalsha-commands)
5. [Node.js Implementation](#nodejs-implementation)
6. [Common Patterns](#common-patterns)
7. [Advanced Examples](#advanced-examples)
8. [Performance Optimization](#performance-optimization)
9. [Best Practices](#best-practices)
10. [Real-world Use Cases](#real-world-use-cases)

## Introduction to Lua Scripting

### Why Use Lua Scripts?

- **Atomicity**: Scripts execute atomically
- **Performance**: Reduced network round trips
- **Complex Logic**: Server-side computation
- **Consistency**: Guaranteed execution order
- **Reusability**: Scripts can be cached and reused

### Key Concepts

- Scripts are executed in a sandboxed environment
- No access to file system or network
- Deterministic execution required
- Scripts block other operations during execution

## Basic Lua Syntax

### Variables and Data Types

```lua
-- Variables
local name = "Redis"
local count = 42
local active = true
local data = nil

-- Tables (arrays and objects)
local array = {1, 2, 3, 4, 5}
local object = {name = "John", age = 30}

-- Accessing table elements
local first = array[1]  -- Lua arrays are 1-indexed
local userName = object.name
local userAge = object["age"]
```

### Control Structures

```lua
-- If statements
if count > 0 then
    return "positive"
elseif count < 0 then
    return "negative"
else
    return "zero"
end

-- For loops
for i = 1, 10 do
    -- loop body
end

for i, v in ipairs(array) do
    -- iterate over array
end

for k, v in pairs(object) do
    -- iterate over object
end

-- While loops
while count > 0 do
    count = count - 1
end
```

### Functions

```lua
-- Function definition
local function add(a, b)
    return a + b
end

-- Anonymous function
local multiply = function(a, b)
    return a * b
end
```

## Redis Lua API

### Available Redis Commands

```lua
-- Basic operations
redis.call('SET', 'key', 'value')
local value = redis.call('GET', 'key')

-- Hash operations
redis.call('HSET', 'hash', 'field', 'value')
local field_value = redis.call('HGET', 'hash', 'field')

-- List operations
redis.call('LPUSH', 'list', 'item')
local items = redis.call('LRANGE', 'list', 0, -1)

-- Set operations
redis.call('SADD', 'set', 'member')
local members = redis.call('SMEMBERS', 'set')
```

### Error Handling

```lua
-- Using redis.call (throws error on Redis error)
local result = redis.call('GET', 'key')

-- Using redis.pcall (returns error as table)
local result = redis.pcall('GET', 'key')
if type(result) == 'table' and result.err then
    return {err = result.err}
end
```

## EVAL and EVALSHA Commands

### Basic Usage in Node.js

```typescript
// src/lua/basic-lua.ts
import { createClient, RedisClientType } from 'redis';

class BasicLuaScripts {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async simpleIncrement(key: string): Promise<number> {
    const script = `
      local current = redis.call('GET', KEYS[1])
      if current == false then
        current = 0
      else
        current = tonumber(current)
      end
      
      local new_value = current + 1
      redis.call('SET', KEYS[1], new_value)
      return new_value
    `;
    
    const result = await this.client.eval(script, {
      keys: [key],
      arguments: []
    });
    
    return result as number;
  }

  async conditionalSet(key: string, expectedValue: string, newValue: string): Promise<boolean> {
    const script = `
      local current = redis.call('GET', KEYS[1])
      if current == ARGV[1] then
        redis.call('SET', KEYS[1], ARGV[2])
        return 1
      else
        return 0
      end
    `;
    
    const result = await this.client.eval(script, {
      keys: [key],
      arguments: [expectedValue, newValue]
    });
    
    return result === 1;
  }

  async getAndExpire(key: string, ttl: number): Promise<string | null> {
    const script = `
      local value = redis.call('GET', KEYS[1])
      if value ~= false then
        redis.call('EXPIRE', KEYS[1], ARGV[1])
      end
      return value
    `;
    
    const result = await this.client.eval(script, {
      keys: [key],
      arguments: [ttl.toString()]
    });
    
    return result as string | null;
  }
}

export { BasicLuaScripts };
```

### Script Caching with EVALSHA

```typescript
// src/lua/script-cache.ts
class LuaScriptCache {
  private client: RedisClientType;
  private scriptHashes: Map<string, string> = new Map();

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async loadScript(name: string, script: string): Promise<string> {
    const sha = await this.client.scriptLoad(script);
    this.scriptHashes.set(name, sha);
    return sha;
  }

  async executeScript(name: string, keys: string[] = [], args: string[] = []): Promise<any> {
    const sha = this.scriptHashes.get(name);
    if (!sha) {
      throw new Error(`Script '${name}' not loaded`);
    }

    try {
      return await this.client.evalSha(sha, {
        keys,
        arguments: args
      });
    } catch (error: any) {
      // If script not in cache, reload and retry
      if (error.message.includes('NOSCRIPT')) {
        console.log(`Script '${name}' not in Redis cache, reloading...`);
        // Note: You'd need to store the original script to reload
        throw new Error('Script needs to be reloaded');
      }
      throw error;
    }
  }

  async loadAllScripts(): Promise<void> {
    const scripts = {
      'rate_limit': `
        local key = KEYS[1]
        local window = tonumber(ARGV[1])
        local limit = tonumber(ARGV[2])
        local current_time = tonumber(ARGV[3])
        
        local current = redis.call('GET', key)
        if current == false then
          redis.call('SET', key, 1)
          redis.call('EXPIRE', key, window)
          return {1, limit - 1}
        end
        
        current = tonumber(current)
        if current < limit then
          local new_count = redis.call('INCR', key)
          local ttl = redis.call('TTL', key)
          return {new_count, limit - new_count}
        else
          local ttl = redis.call('TTL', key)
          return {current, 0, ttl}
        end
      `,
      
      'sliding_window': `
        local key = KEYS[1]
        local window = tonumber(ARGV[1])
        local limit = tonumber(ARGV[2])
        local current_time = tonumber(ARGV[3])
        
        -- Remove old entries
        redis.call('ZREMRANGEBYSCORE', key, 0, current_time - window * 1000)
        
        -- Count current entries
        local current_count = redis.call('ZCARD', key)
        
        if current_count < limit then
          -- Add current request
          redis.call('ZADD', key, current_time, current_time)
          redis.call('EXPIRE', key, window + 1)
          return {current_count + 1, limit - current_count - 1}
        else
          return {current_count, 0}
        end
      `,
      
      'atomic_counter': `
        local key = KEYS[1]
        local increment = tonumber(ARGV[1]) or 1
        local max_value = tonumber(ARGV[2])
        
        local current = redis.call('GET', key)
        if current == false then
          current = 0
        else
          current = tonumber(current)
        end
        
        if max_value and current + increment > max_value then
          return {current, false}
        end
        
        local new_value = redis.call('INCRBY', key, increment)
        return {new_value, true}
      `
    };

    for (const [name, script] of Object.entries(scripts)) {
      await this.loadScript(name, script);
    }
  }
}

export { LuaScriptCache };
```

## Common Patterns

### Rate Limiting

```typescript
// src/lua/rate-limiter.ts
class LuaRateLimiter {
  private client: RedisClientType;
  private scriptCache: LuaScriptCache;

  constructor(client: RedisClientType) {
    this.client = client;
    this.scriptCache = new LuaScriptCache(client);
  }

  async initialize(): Promise<void> {
    await this.scriptCache.loadAllScripts();
  }

  async checkRateLimit(
    identifier: string,
    windowSeconds: number,
    maxRequests: number
  ): Promise<{
    allowed: boolean;
    count: number;
    remaining: number;
    resetTime?: number;
  }> {
    const key = `rate_limit:${identifier}`;
    const currentTime = Date.now();
    
    const result = await this.scriptCache.executeScript(
      'rate_limit',
      [key],
      [windowSeconds.toString(), maxRequests.toString(), currentTime.toString()]
    ) as number[];
    
    const [count, remaining, resetTime] = result;
    
    return {
      allowed: remaining > 0,
      count,
      remaining: Math.max(0, remaining),
      resetTime: resetTime ? Date.now() + (resetTime * 1000) : undefined
    };
  }

  async checkSlidingWindow(
    identifier: string,
    windowSeconds: number,
    maxRequests: number
  ): Promise<{
    allowed: boolean;
    count: number;
    remaining: number;
  }> {
    const key = `sliding_window:${identifier}`;
    const currentTime = Date.now();
    
    const result = await this.scriptCache.executeScript(
      'sliding_window',
      [key],
      [windowSeconds.toString(), maxRequests.toString(), currentTime.toString()]
    ) as number[];
    
    const [count, remaining] = result;
    
    return {
      allowed: remaining > 0,
      count,
      remaining: Math.max(0, remaining)
    };
  }
}

export { LuaRateLimiter };
```

### Distributed Lock

```typescript
// src/lua/distributed-lock.ts
class LuaDistributedLock {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async acquireLock(
    resource: string,
    identifier: string,
    ttlMs: number
  ): Promise<boolean> {
    const script = `
      local key = KEYS[1]
      local identifier = ARGV[1]
      local ttl = tonumber(ARGV[2])
      
      local current = redis.call('GET', key)
      if current == false then
        redis.call('SET', key, identifier, 'PX', ttl)
        return 1
      elseif current == identifier then
        -- Extend our own lock
        redis.call('PEXPIRE', key, ttl)
        return 1
      else
        return 0
      end
    `;
    
    const result = await this.client.eval(script, {
      keys: [`lock:${resource}`],
      arguments: [identifier, ttlMs.toString()]
    });
    
    return result === 1;
  }

  async releaseLock(resource: string, identifier: string): Promise<boolean> {
    const script = `
      local key = KEYS[1]
      local identifier = ARGV[1]
      
      local current = redis.call('GET', key)
      if current == identifier then
        redis.call('DEL', key)
        return 1
      else
        return 0
      end
    `;
    
    const result = await this.client.eval(script, {
      keys: [`lock:${resource}`],
      arguments: [identifier]
    });
    
    return result === 1;
  }

  async extendLock(
    resource: string,
    identifier: string,
    ttlMs: number
  ): Promise<boolean> {
    const script = `
      local key = KEYS[1]
      local identifier = ARGV[1]
      local ttl = tonumber(ARGV[2])
      
      local current = redis.call('GET', key)
      if current == identifier then
        redis.call('PEXPIRE', key, ttl)
        return 1
      else
        return 0
      end
    `;
    
    const result = await this.client.eval(script, {
      keys: [`lock:${resource}`],
      arguments: [identifier, ttlMs.toString()]
    });
    
    return result === 1;
  }
}

export { LuaDistributedLock };
```

### Atomic Operations

```typescript
// src/lua/atomic-operations.ts
class AtomicOperations {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async atomicTransfer(
    fromKey: string,
    toKey: string,
    amount: number
  ): Promise<{
    success: boolean;
    fromBalance?: number;
    toBalance?: number;
    error?: string;
  }> {
    const script = `
      local from_key = KEYS[1]
      local to_key = KEYS[2]
      local amount = tonumber(ARGV[1])
      
      local from_balance = redis.call('GET', from_key)
      if from_balance == false then
        return {0, 'Source account not found'}
      end
      
      from_balance = tonumber(from_balance)
      if from_balance < amount then
        return {0, 'Insufficient funds'}
      end
      
      local to_balance = redis.call('GET', to_key)
      if to_balance == false then
        to_balance = 0
      else
        to_balance = tonumber(to_balance)
      end
      
      -- Perform transfer
      local new_from = from_balance - amount
      local new_to = to_balance + amount
      
      redis.call('SET', from_key, new_from)
      redis.call('SET', to_key, new_to)
      
      return {1, new_from, new_to}
    `;
    
    const result = await this.client.eval(script, {
      keys: [fromKey, toKey],
      arguments: [amount.toString()]
    }) as any[];
    
    if (result[0] === 0) {
      return { success: false, error: result[1] };
    }
    
    return {
      success: true,
      fromBalance: result[1],
      toBalance: result[2]
    };
  }

  async atomicIncrement(
    key: string,
    increment: number = 1,
    maxValue?: number
  ): Promise<{
    success: boolean;
    value: number;
    hitLimit?: boolean;
  }> {
    const script = `
      local key = KEYS[1]
      local increment = tonumber(ARGV[1])
      local max_value = tonumber(ARGV[2])
      
      local current = redis.call('GET', key)
      if current == false then
        current = 0
      else
        current = tonumber(current)
      end
      
      if max_value and current + increment > max_value then
        return {0, current, 1}
      end
      
      local new_value = redis.call('INCRBY', key, increment)
      return {1, new_value, 0}
    `;
    
    const result = await this.client.eval(script, {
      keys: [key],
      arguments: [increment.toString(), maxValue?.toString() || 'nil']
    }) as number[];
    
    return {
      success: result[0] === 1,
      value: result[1],
      hitLimit: result[2] === 1
    };
  }

  async atomicListPushWithLimit(
    key: string,
    value: string,
    maxLength: number
  ): Promise<{
    success: boolean;
    length: number;
    removed?: string[];
  }> {
    const script = `
      local key = KEYS[1]
      local value = ARGV[1]
      local max_length = tonumber(ARGV[2])
      
      -- Push new value
      local new_length = redis.call('LPUSH', key, value)
      
      local removed = {}
      if new_length > max_length then
        -- Remove excess elements
        local excess = new_length - max_length
        for i = 1, excess do
          local removed_item = redis.call('RPOP', key)
          table.insert(removed, removed_item)
        end
        new_length = max_length
      end
      
      return {1, new_length, removed}
    `;
    
    const result = await this.client.eval(script, {
      keys: [key],
      arguments: [value, maxLength.toString()]
    }) as any[];
    
    return {
      success: result[0] === 1,
      length: result[1],
      removed: result[2] || []
    };
  }
}

export { AtomicOperations };
```

## Advanced Examples

### Leaderboard with Lua

```typescript
// src/lua/leaderboard.ts
class LuaLeaderboard {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async updateScore(
    leaderboardKey: string,
    playerId: string,
    score: number,
    operation: 'set' | 'increment' = 'set'
  ): Promise<{
    newScore: number;
    rank: number;
    totalPlayers: number;
  }> {
    const script = `
      local key = KEYS[1]
      local player_id = ARGV[1]
      local score = tonumber(ARGV[2])
      local operation = ARGV[3]
      
      local new_score
      if operation == 'increment' then
        new_score = redis.call('ZINCRBY', key, score, player_id)
      else
        redis.call('ZADD', key, score, player_id)
        new_score = score
      end
      
      -- Get rank (1-based)
      local rank = redis.call('ZREVRANK', key, player_id)
      if rank then
        rank = rank + 1
      else
        rank = -1
      end
      
      -- Get total players
      local total = redis.call('ZCARD', key)
      
      return {new_score, rank, total}
    `;
    
    const result = await this.client.eval(script, {
      keys: [leaderboardKey],
      arguments: [playerId, score.toString(), operation]
    }) as number[];
    
    return {
      newScore: result[0],
      rank: result[1],
      totalPlayers: result[2]
    };
  }

  async getPlayerStats(
    leaderboardKey: string,
    playerId: string
  ): Promise<{
    score: number | null;
    rank: number | null;
    totalPlayers: number;
  }> {
    const script = `
      local key = KEYS[1]
      local player_id = ARGV[1]
      
      local score = redis.call('ZSCORE', key, player_id)
      local rank = redis.call('ZREVRANK', key, player_id)
      local total = redis.call('ZCARD', key)
      
      if rank then
        rank = rank + 1
      end
      
      return {score, rank, total}
    `;
    
    const result = await this.client.eval(script, {
      keys: [leaderboardKey],
      arguments: [playerId]
    }) as any[];
    
    return {
      score: result[0] ? parseFloat(result[0]) : null,
      rank: result[1] ? parseInt(result[1]) : null,
      totalPlayers: result[2]
    };
  }

  async getTopPlayers(
    leaderboardKey: string,
    count: number = 10
  ): Promise<Array<{
    playerId: string;
    score: number;
    rank: number;
  }>> {
    const script = `
      local key = KEYS[1]
      local count = tonumber(ARGV[1])
      
      local top_players = redis.call('ZREVRANGE', key, 0, count - 1, 'WITHSCORES')
      
      local result = {}
      for i = 1, #top_players, 2 do
        local player_id = top_players[i]
        local score = tonumber(top_players[i + 1])
        local rank = math.floor(i / 2) + 1
        
        table.insert(result, {player_id, score, rank})
      end
      
      return result
    `;
    
    const result = await this.client.eval(script, {
      keys: [leaderboardKey],
      arguments: [count.toString()]
    }) as any[];
    
    return result.map((player: any[]) => ({
      playerId: player[0],
      score: player[1],
      rank: player[2]
    }));
  }
}

export { LuaLeaderboard };
```

### Session Management

```typescript
// src/lua/session-manager.ts
interface SessionData {
  userId: string;
  data: Record<string, any>;
  createdAt: number;
  lastAccessed: number;
}

class LuaSessionManager {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async createSession(
    sessionId: string,
    userId: string,
    data: Record<string, any>,
    ttlSeconds: number
  ): Promise<boolean> {
    const script = `
      local session_key = KEYS[1]
      local user_sessions_key = KEYS[2]
      local session_id = ARGV[1]
      local user_id = ARGV[2]
      local session_data = ARGV[3]
      local ttl = tonumber(ARGV[4])
      local current_time = tonumber(ARGV[5])
      
      -- Check if session already exists
      local existing = redis.call('EXISTS', session_key)
      if existing == 1 then
        return 0
      end
      
      -- Create session
      local session = {
        userId = user_id,
        data = session_data,
        createdAt = current_time,
        lastAccessed = current_time
      }
      
      redis.call('HMSET', session_key,
        'userId', user_id,
        'data', session_data,
        'createdAt', current_time,
        'lastAccessed', current_time
      )
      
      redis.call('EXPIRE', session_key, ttl)
      
      -- Add to user's session list
      redis.call('SADD', user_sessions_key, session_id)
      redis.call('EXPIRE', user_sessions_key, ttl + 3600)
      
      return 1
    `;
    
    const result = await this.client.eval(script, {
      keys: [`session:${sessionId}`, `user_sessions:${userId}`],
      arguments: [
        sessionId,
        userId,
        JSON.stringify(data),
        ttlSeconds.toString(),
        Date.now().toString()
      ]
    });
    
    return result === 1;
  }

  async getSession(sessionId: string): Promise<SessionData | null> {
    const script = `
      local session_key = KEYS[1]
      local current_time = tonumber(ARGV[1])
      
      local session = redis.call('HMGET', session_key,
        'userId', 'data', 'createdAt', 'lastAccessed'
      )
      
      if session[1] == false then
        return nil
      end
      
      -- Update last accessed time
      redis.call('HSET', session_key, 'lastAccessed', current_time)
      
      return session
    `;
    
    const result = await this.client.eval(script, {
      keys: [`session:${sessionId}`],
      arguments: [Date.now().toString()]
    }) as string[] | null;
    
    if (!result || !result[0]) {
      return null;
    }
    
    return {
      userId: result[0],
      data: JSON.parse(result[1] || '{}'),
      createdAt: parseInt(result[2]),
      lastAccessed: parseInt(result[3])
    };
  }

  async updateSession(
    sessionId: string,
    data: Record<string, any>,
    ttlSeconds?: number
  ): Promise<boolean> {
    const script = `
      local session_key = KEYS[1]
      local session_data = ARGV[1]
      local current_time = tonumber(ARGV[2])
      local ttl = tonumber(ARGV[3])
      
      local exists = redis.call('EXISTS', session_key)
      if exists == 0 then
        return 0
      end
      
      redis.call('HMSET', session_key,
        'data', session_data,
        'lastAccessed', current_time
      )
      
      if ttl and ttl > 0 then
        redis.call('EXPIRE', session_key, ttl)
      end
      
      return 1
    `;
    
    const result = await this.client.eval(script, {
      keys: [`session:${sessionId}`],
      arguments: [
        JSON.stringify(data),
        Date.now().toString(),
        ttlSeconds?.toString() || '0'
      ]
    });
    
    return result === 1;
  }

  async destroySession(sessionId: string): Promise<boolean> {
    const script = `
      local session_key = KEYS[1]
      
      -- Get user ID before deleting
      local user_id = redis.call('HGET', session_key, 'userId')
      
      if user_id then
        local user_sessions_key = 'user_sessions:' .. user_id
        redis.call('SREM', user_sessions_key, ARGV[1])
      end
      
      local deleted = redis.call('DEL', session_key)
      return deleted
    `;
    
    const result = await this.client.eval(script, {
      keys: [`session:${sessionId}`],
      arguments: [sessionId]
    });
    
    return result === 1;
  }

  async destroyAllUserSessions(userId: string): Promise<number> {
    const script = `
      local user_sessions_key = KEYS[1]
      local user_id = ARGV[1]
      
      local session_ids = redis.call('SMEMBERS', user_sessions_key)
      local deleted_count = 0
      
      for i = 1, #session_ids do
        local session_key = 'session:' .. session_ids[i]
        local deleted = redis.call('DEL', session_key)
        deleted_count = deleted_count + deleted
      end
      
      redis.call('DEL', user_sessions_key)
      
      return deleted_count
    `;
    
    const result = await this.client.eval(script, {
      keys: [`user_sessions:${userId}`],
      arguments: [userId]
    });
    
    return result as number;
  }
}

export { LuaSessionManager, SessionData };
```

## Usage Examples

```typescript
// src/examples/lua-usage.ts
import { createClient } from 'redis';
import { BasicLuaScripts } from '../lua/basic-lua';
import { LuaRateLimiter } from '../lua/rate-limiter';
import { LuaDistributedLock } from '../lua/distributed-lock';
import { AtomicOperations } from '../lua/atomic-operations';
import { LuaLeaderboard } from '../lua/leaderboard';
import { LuaSessionManager } from '../lua/session-manager';

async function runLuaExamples() {
  const client = createClient({
    host: 'localhost',
    port: 6379
  });
  
  await client.connect();
  
  try {
    // Basic Lua Scripts
    console.log('=== Basic Lua Scripts ===');
    const basicLua = new BasicLuaScripts(client);
    
    const incrementResult = await basicLua.simpleIncrement('counter');
    console.log('Increment result:', incrementResult);
    
    const conditionalResult = await basicLua.conditionalSet('test_key', 'old_value', 'new_value');
    console.log('Conditional set result:', conditionalResult);
    
    // Rate Limiting
    console.log('\n=== Rate Limiting ===');
    const rateLimiter = new LuaRateLimiter(client);
    await rateLimiter.initialize();
    
    const rateResult = await rateLimiter.checkRateLimit('user123', 60, 10);
    console.log('Rate limit result:', rateResult);
    
    // Distributed Lock
    console.log('\n=== Distributed Lock ===');
    const lock = new LuaDistributedLock(client);
    
    const lockAcquired = await lock.acquireLock('resource1', 'worker1', 5000);
    console.log('Lock acquired:', lockAcquired);
    
    if (lockAcquired) {
      // Do some work
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const lockReleased = await lock.releaseLock('resource1', 'worker1');
      console.log('Lock released:', lockReleased);
    }
    
    // Atomic Operations
    console.log('\n=== Atomic Operations ===');
    const atomic = new AtomicOperations(client);
    
    // Set up accounts
    await client.set('account:alice', '1000');
    await client.set('account:bob', '500');
    
    const transferResult = await atomic.atomicTransfer('account:alice', 'account:bob', 200);
    console.log('Transfer result:', transferResult);
    
    // Leaderboard
    console.log('\n=== Leaderboard ===');
    const leaderboard = new LuaLeaderboard(client);
    
    await leaderboard.updateScore('game_scores', 'player1', 1500);
    await leaderboard.updateScore('game_scores', 'player2', 1200);
    await leaderboard.updateScore('game_scores', 'player3', 1800);
    
    const topPlayers = await leaderboard.getTopPlayers('game_scores', 3);
    console.log('Top players:', topPlayers);
    
    const playerStats = await leaderboard.getPlayerStats('game_scores', 'player1');
    console.log('Player1 stats:', playerStats);
    
    // Session Management
    console.log('\n=== Session Management ===');
    const sessionManager = new LuaSessionManager(client);
    
    const sessionCreated = await sessionManager.createSession(
      'sess_123',
      'user_456',
      { theme: 'dark', language: 'en' },
      3600
    );
    console.log('Session created:', sessionCreated);
    
    const session = await sessionManager.getSession('sess_123');
    console.log('Retrieved session:', session);
    
  } finally {
    await client.disconnect();
  }
}

// Run examples
if (require.main === module) {
  runLuaExamples().catch(console.error);
}

export { runLuaExamples };
```

## Performance Optimization

### Script Optimization Tips

1. **Minimize Redis Calls**: Combine operations when possible
2. **Use Local Variables**: Store frequently accessed values
3. **Avoid Loops**: Use Redis bulk operations instead
4. **Cache Scripts**: Use EVALSHA for repeated executions
5. **Keep Scripts Short**: Long scripts block other operations

### Example: Optimized Batch Operations

```lua
-- Optimized batch hash update
local key = KEYS[1]
local fields_and_values = {}

for i = 1, #ARGV, 2 do
  table.insert(fields_and_values, ARGV[i])
  table.insert(fields_and_values, ARGV[i + 1])
end

if #fields_and_values > 0 then
  redis.call('HMSET', key, unpack(fields_and_values))
end

return redis.call('HLEN', key)
```

## Best Practices

### 1. Script Design

- Keep scripts atomic and deterministic
- Avoid time-dependent operations
- Use proper error handling
- Document script parameters clearly

### 2. Error Handling

```lua
-- Good error handling
local function safe_operation()
  local result = redis.pcall('GET', 'some_key')
  if type(result) == 'table' and result.err then
    return {err = 'Operation failed: ' .. result.err}
  end
  return result
end
```

### 3. Performance Considerations

- Use EVALSHA for frequently executed scripts
- Minimize script execution time
- Avoid blocking operations
- Use appropriate data structures

### 4. Security

- Validate input parameters
- Avoid dynamic script generation
- Use parameterized scripts
- Limit script complexity

## Next Steps

Now that you understand Lua scripting, you're ready to explore:

- [Part 10: Caching Strategies](./part-10-caching-strategies.md)

## Quick Reference

### Basic Lua Syntax
```lua
-- Variables
local value = redis.call('GET', KEYS[1])
local number = tonumber(value) or 0

-- Conditionals
if value then
  return value
else
  return 'default'
end

-- Loops
for i = 1, #ARGV do
  redis.call('SET', 'key' .. i, ARGV[i])
end

-- Tables
local result = {}
table.insert(result, value)
return result
```

### Redis Commands in Lua
```lua
-- Basic operations
redis.call('SET', 'key', 'value')
local value = redis.call('GET', 'key')

-- Error handling
local result = redis.pcall('GET', 'key')
if type(result) == 'table' and result.err then
  return {err = result.err}
end

-- Bulk operations
redis.call('HMSET', 'hash', 'field1', 'value1', 'field2', 'value2')
```

### Node.js Integration
```typescript
// Execute script
const result = await client.eval(script, {
  keys: ['key1', 'key2'],
  arguments: ['arg1', 'arg2']
});

// Execute cached script
const sha = await client.scriptLoad(script);
const result = await client.evalSha(sha, {
  keys: ['key1'],
  arguments: ['arg1']
});
```

---

**Previous**: [← Part 8: Transactions and Pipelining](./part-08-transactions.md)  
**Next**: [Part 10: Caching Strategies →](./part-10-caching-strategies.md)