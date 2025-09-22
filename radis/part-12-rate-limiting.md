# Part 12: Rate Limiting and Security with Redis

Redis is excellent for implementing rate limiting and security features due to its atomic operations, expiration capabilities, and high performance. This part covers comprehensive rate limiting strategies and security implementations.

## Table of Contents

1. [Rate Limiting Fundamentals](#rate-limiting-fundamentals)
2. [Basic Rate Limiting](#basic-rate-limiting)
3. [Advanced Rate Limiting Algorithms](#advanced-rate-limiting-algorithms)
4. [Distributed Rate Limiting](#distributed-rate-limiting)
5. [Security Features](#security-features)
6. [API Protection](#api-protection)
7. [DDoS Protection](#ddos-protection)
8. [Real-world Examples](#real-world-examples)
9. [Performance Optimization](#performance-optimization)

## Rate Limiting Fundamentals

### Why Rate Limiting?

- **Prevent Abuse**: Protect against malicious attacks
- **Resource Protection**: Ensure fair usage of resources
- **Cost Control**: Manage API usage costs
- **Quality of Service**: Maintain performance for all users
- **Compliance**: Meet SLA requirements

### Common Rate Limiting Algorithms

1. **Fixed Window**: Simple time-based windows
2. **Sliding Window**: More accurate time-based limiting
3. **Token Bucket**: Burst handling with sustained rates
4. **Leaky Bucket**: Smooth rate enforcement
5. **Sliding Window Log**: Precise request tracking

## Basic Rate Limiting

### Simple Fixed Window Rate Limiter

```typescript
// src/rate-limiting/fixed-window-limiter.ts
import { RedisClientType } from 'redis';

interface RateLimitOptions {
  windowSizeMs: number;
  maxRequests: number;
  keyPrefix?: string;
}

interface RateLimitResult {
  allowed: boolean;
  remaining: number;
  resetTime: number;
  totalRequests: number;
}

class FixedWindowRateLimiter {
  private client: RedisClientType;
  private options: RateLimitOptions;

  constructor(client: RedisClientType, options: RateLimitOptions) {
    this.client = client;
    this.options = {
      keyPrefix: 'rate_limit:',
      ...options
    };
  }

  private getKey(identifier: string): string {
    const window = Math.floor(Date.now() / this.options.windowSizeMs);
    return `${this.options.keyPrefix}${identifier}:${window}`;
  }

  async checkLimit(identifier: string): Promise<RateLimitResult> {
    const key = this.getKey(identifier);
    const windowStart = Math.floor(Date.now() / this.options.windowSizeMs) * this.options.windowSizeMs;
    const resetTime = windowStart + this.options.windowSizeMs;

    // Use Lua script for atomic operation
    const luaScript = `
      local key = KEYS[1]
      local limit = tonumber(ARGV[1])
      local window = tonumber(ARGV[2])
      
      local current = redis.call('GET', key)
      if current == false then
        current = 0
      else
        current = tonumber(current)
      end
      
      if current < limit then
        local new_val = redis.call('INCR', key)
        if new_val == 1 then
          redis.call('EXPIRE', key, window)
        end
        return {1, limit - new_val, new_val}
      else
        return {0, 0, current}
      end
    `;

    const result = await this.client.eval(luaScript, {
      keys: [key],
      arguments: [
        this.options.maxRequests.toString(),
        Math.ceil(this.options.windowSizeMs / 1000).toString()
      ]
    }) as number[];

    return {
      allowed: result[0] === 1,
      remaining: result[1],
      resetTime,
      totalRequests: result[2]
    };
  }

  async getRemainingRequests(identifier: string): Promise<number> {
    const key = this.getKey(identifier);
    const current = await this.client.get(key);
    const used = current ? parseInt(current) : 0;
    return Math.max(0, this.options.maxRequests - used);
  }

  async resetLimit(identifier: string): Promise<boolean> {
    const key = this.getKey(identifier);
    const result = await this.client.del(key);
    return result > 0;
  }
}

export { FixedWindowRateLimiter, RateLimitOptions, RateLimitResult };
```

### Sliding Window Rate Limiter

```typescript
// src/rate-limiting/sliding-window-limiter.ts
class SlidingWindowRateLimiter {
  private client: RedisClientType;
  private windowSizeMs: number;
  private maxRequests: number;
  private keyPrefix: string;

  constructor(
    client: RedisClientType,
    windowSizeMs: number,
    maxRequests: number,
    keyPrefix = 'sliding_rate_limit:'
  ) {
    this.client = client;
    this.windowSizeMs = windowSizeMs;
    this.maxRequests = maxRequests;
    this.keyPrefix = keyPrefix;
  }

  private getKey(identifier: string): string {
    return `${this.keyPrefix}${identifier}`;
  }

  async checkLimit(identifier: string): Promise<RateLimitResult> {
    const key = this.getKey(identifier);
    const now = Date.now();
    const windowStart = now - this.windowSizeMs;

    const luaScript = `
      local key = KEYS[1]
      local window_start = tonumber(ARGV[1])
      local now = tonumber(ARGV[2])
      local limit = tonumber(ARGV[3])
      
      -- Remove expired entries
      redis.call('ZREMRANGEBYSCORE', key, 0, window_start)
      
      -- Count current requests
      local current = redis.call('ZCARD', key)
      
      if current < limit then
        -- Add current request
        redis.call('ZADD', key, now, now)
        redis.call('EXPIRE', key, math.ceil((tonumber(ARGV[4])) / 1000))
        return {1, limit - current - 1, current + 1}
      else
        return {0, 0, current}
      end
    `;

    const result = await this.client.eval(luaScript, {
      keys: [key],
      arguments: [
        windowStart.toString(),
        now.toString(),
        this.maxRequests.toString(),
        this.windowSizeMs.toString()
      ]
    }) as number[];

    return {
      allowed: result[0] === 1,
      remaining: result[1],
      resetTime: now + this.windowSizeMs,
      totalRequests: result[2]
    };
  }

  async getRequestTimes(identifier: string): Promise<number[]> {
    const key = this.getKey(identifier);
    const now = Date.now();
    const windowStart = now - this.windowSizeMs;
    
    // Clean up expired entries and get remaining
    await this.client.zRemRangeByScore(key, 0, windowStart);
    const scores = await this.client.zRange(key, 0, -1, { BY: 'SCORE' });
    
    return scores.map(score => parseFloat(score.toString()));
  }

  async resetLimit(identifier: string): Promise<boolean> {
    const key = this.getKey(identifier);
    const result = await this.client.del(key);
    return result > 0;
  }
}

export { SlidingWindowRateLimiter };
```

## Advanced Rate Limiting Algorithms

### Token Bucket Rate Limiter

```typescript
// src/rate-limiting/token-bucket-limiter.ts
interface TokenBucketOptions {
  capacity: number;
  refillRate: number; // tokens per second
  refillPeriodMs?: number;
  keyPrefix?: string;
}

interface TokenBucketState {
  tokens: number;
  lastRefill: number;
}

class TokenBucketRateLimiter {
  private client: RedisClientType;
  private options: TokenBucketOptions;

  constructor(client: RedisClientType, options: TokenBucketOptions) {
    this.client = client;
    this.options = {
      refillPeriodMs: 1000,
      keyPrefix: 'token_bucket:',
      ...options
    };
  }

  private getKey(identifier: string): string {
    return `${this.options.keyPrefix}${identifier}`;
  }

  async checkLimit(identifier: string, tokensRequested = 1): Promise<RateLimitResult & {
    tokensAvailable: number;
  }> {
    const key = this.getKey(identifier);
    const now = Date.now();

    const luaScript = `
      local key = KEYS[1]
      local capacity = tonumber(ARGV[1])
      local refill_rate = tonumber(ARGV[2])
      local refill_period = tonumber(ARGV[3])
      local tokens_requested = tonumber(ARGV[4])
      local now = tonumber(ARGV[5])
      
      local bucket = redis.call('HMGET', key, 'tokens', 'last_refill')
      local tokens = tonumber(bucket[1]) or capacity
      local last_refill = tonumber(bucket[2]) or now
      
      -- Calculate tokens to add
      local time_passed = now - last_refill
      local tokens_to_add = math.floor((time_passed / refill_period) * refill_rate)
      tokens = math.min(capacity, tokens + tokens_to_add)
      
      local allowed = 0
      local remaining_tokens = tokens
      
      if tokens >= tokens_requested then
        tokens = tokens - tokens_requested
        allowed = 1
        remaining_tokens = tokens
      end
      
      -- Update bucket state
      redis.call('HMSET', key, 'tokens', tokens, 'last_refill', now)
      redis.call('EXPIRE', key, math.ceil(capacity / refill_rate * refill_period / 1000))
      
      return {allowed, remaining_tokens, tokens_requested}
    `;

    const result = await this.client.eval(luaScript, {
      keys: [key],
      arguments: [
        this.options.capacity.toString(),
        this.options.refillRate.toString(),
        this.options.refillPeriodMs!.toString(),
        tokensRequested.toString(),
        now.toString()
      ]
    }) as number[];

    return {
      allowed: result[0] === 1,
      remaining: result[1],
      resetTime: now + (this.options.capacity / this.options.refillRate * this.options.refillPeriodMs!),
      totalRequests: result[2],
      tokensAvailable: result[1]
    };
  }

  async getBucketState(identifier: string): Promise<TokenBucketState | null> {
    const key = this.getKey(identifier);
    const result = await this.client.hmGet(key, ['tokens', 'last_refill']);
    
    if (!result.tokens) {
      return null;
    }

    return {
      tokens: parseFloat(result.tokens),
      lastRefill: parseFloat(result.last_refill || '0')
    };
  }

  async addTokens(identifier: string, tokensToAdd: number): Promise<number> {
    const key = this.getKey(identifier);
    
    const luaScript = `
      local key = KEYS[1]
      local capacity = tonumber(ARGV[1])
      local tokens_to_add = tonumber(ARGV[2])
      local now = tonumber(ARGV[3])
      
      local current_tokens = tonumber(redis.call('HGET', key, 'tokens')) or capacity
      local new_tokens = math.min(capacity, current_tokens + tokens_to_add)
      
      redis.call('HMSET', key, 'tokens', new_tokens, 'last_refill', now)
      
      return new_tokens
    `;

    const result = await this.client.eval(luaScript, {
      keys: [key],
      arguments: [
        this.options.capacity.toString(),
        tokensToAdd.toString(),
        Date.now().toString()
      ]
    }) as number;

    return result;
  }
}

export { TokenBucketRateLimiter, TokenBucketOptions, TokenBucketState };
```

### Leaky Bucket Rate Limiter

```typescript
// src/rate-limiting/leaky-bucket-limiter.ts
interface LeakyBucketOptions {
  capacity: number;
  leakRate: number; // requests per second
  keyPrefix?: string;
}

class LeakyBucketRateLimiter {
  private client: RedisClientType;
  private options: LeakyBucketOptions;

  constructor(client: RedisClientType, options: LeakyBucketOptions) {
    this.client = client;
    this.options = {
      keyPrefix: 'leaky_bucket:',
      ...options
    };
  }

  private getKey(identifier: string): string {
    return `${this.options.keyPrefix}${identifier}`;
  }

  async checkLimit(identifier: string): Promise<RateLimitResult & {
    queueSize: number;
    estimatedWaitTime: number;
  }> {
    const key = this.getKey(identifier);
    const now = Date.now();

    const luaScript = `
      local key = KEYS[1]
      local capacity = tonumber(ARGV[1])
      local leak_rate = tonumber(ARGV[2])
      local now = tonumber(ARGV[3])
      
      local bucket = redis.call('HMGET', key, 'volume', 'last_leak')
      local volume = tonumber(bucket[1]) or 0
      local last_leak = tonumber(bucket[2]) or now
      
      -- Calculate leaked volume
      local time_passed = (now - last_leak) / 1000
      local leaked = time_passed * leak_rate
      volume = math.max(0, volume - leaked)
      
      local allowed = 0
      local wait_time = 0
      
      if volume < capacity then
        volume = volume + 1
        allowed = 1
      else
        -- Calculate wait time
        wait_time = ((volume - capacity + 1) / leak_rate) * 1000
      end
      
      -- Update bucket state
      redis.call('HMSET', key, 'volume', volume, 'last_leak', now)
      redis.call('EXPIRE', key, math.ceil(capacity / leak_rate))
      
      return {allowed, capacity - volume, volume, wait_time}
    `;

    const result = await this.client.eval(luaScript, {
      keys: [key],
      arguments: [
        this.options.capacity.toString(),
        this.options.leakRate.toString(),
        now.toString()
      ]
    }) as number[];

    return {
      allowed: result[0] === 1,
      remaining: Math.floor(result[1]),
      resetTime: now + (result[2] / this.options.leakRate * 1000),
      totalRequests: 1,
      queueSize: Math.floor(result[2]),
      estimatedWaitTime: result[3]
    };
  }

  async getBucketVolume(identifier: string): Promise<number> {
    const key = this.getKey(identifier);
    const now = Date.now();
    
    const luaScript = `
      local key = KEYS[1]
      local leak_rate = tonumber(ARGV[1])
      local now = tonumber(ARGV[2])
      
      local bucket = redis.call('HMGET', key, 'volume', 'last_leak')
      local volume = tonumber(bucket[1]) or 0
      local last_leak = tonumber(bucket[2]) or now
      
      local time_passed = (now - last_leak) / 1000
      local leaked = time_passed * leak_rate
      volume = math.max(0, volume - leaked)
      
      return volume
    `;

    const result = await this.client.eval(luaScript, {
      keys: [key],
      arguments: [
        this.options.leakRate.toString(),
        now.toString()
      ]
    }) as number;

    return result;
  }
}

export { LeakyBucketRateLimiter, LeakyBucketOptions };
```

## Distributed Rate Limiting

### Multi-Tier Rate Limiter

```typescript
// src/rate-limiting/multi-tier-limiter.ts
interface TierConfig {
  name: string;
  windowSizeMs: number;
  maxRequests: number;
  priority: number;
}

interface MultiTierOptions {
  tiers: TierConfig[];
  keyPrefix?: string;
  strategy?: 'strict' | 'graceful';
}

class MultiTierRateLimiter {
  private client: RedisClientType;
  private options: MultiTierOptions;
  private limiters: Map<string, FixedWindowRateLimiter>;

  constructor(client: RedisClientType, options: MultiTierOptions) {
    this.client = client;
    this.options = {
      keyPrefix: 'multi_tier:',
      strategy: 'strict',
      ...options
    };
    
    this.limiters = new Map();
    
    // Create limiters for each tier
    for (const tier of this.options.tiers) {
      this.limiters.set(tier.name, new FixedWindowRateLimiter(client, {
        windowSizeMs: tier.windowSizeMs,
        maxRequests: tier.maxRequests,
        keyPrefix: `${this.options.keyPrefix}${tier.name}:`
      }));
    }
  }

  async checkLimit(identifier: string): Promise<{
    allowed: boolean;
    limitingTier?: string;
    tierResults: Array<{
      tier: string;
      result: RateLimitResult;
    }>;
  }> {
    const tierResults: Array<{ tier: string; result: RateLimitResult }> = [];
    
    // Sort tiers by priority
    const sortedTiers = [...this.options.tiers].sort((a, b) => a.priority - b.priority);
    
    for (const tier of sortedTiers) {
      const limiter = this.limiters.get(tier.name)!;
      const result = await limiter.checkLimit(identifier);
      
      tierResults.push({ tier: tier.name, result });
      
      if (!result.allowed) {
        if (this.options.strategy === 'strict') {
          return {
            allowed: false,
            limitingTier: tier.name,
            tierResults
          };
        }
      }
    }
    
    return {
      allowed: true,
      tierResults
    };
  }

  async getTierStatus(identifier: string): Promise<Array<{
    tier: string;
    remaining: number;
    resetTime: number;
  }>> {
    const status = [];
    
    for (const tier of this.options.tiers) {
      const limiter = this.limiters.get(tier.name)!;
      const remaining = await limiter.getRemainingRequests(identifier);
      
      status.push({
        tier: tier.name,
        remaining,
        resetTime: Date.now() + tier.windowSizeMs
      });
    }
    
    return status;
  }

  async resetAllTiers(identifier: string): Promise<number> {
    let resetCount = 0;
    
    for (const tier of this.options.tiers) {
      const limiter = this.limiters.get(tier.name)!;
      const reset = await limiter.resetLimit(identifier);
      if (reset) resetCount++;
    }
    
    return resetCount;
  }
}

export { MultiTierRateLimiter, TierConfig, MultiTierOptions };
```

### Distributed Rate Limiter with Consensus

```typescript
// src/rate-limiting/distributed-consensus-limiter.ts
interface DistributedNode {
  id: string;
  weight: number;
  lastSeen: number;
}

class DistributedConsensusRateLimiter {
  private client: RedisClientType;
  private nodeId: string;
  private windowSizeMs: number;
  private maxRequests: number;
  private keyPrefix: string;
  private nodeWeight: number;

  constructor(
    client: RedisClientType,
    nodeId: string,
    windowSizeMs: number,
    maxRequests: number,
    nodeWeight = 1,
    keyPrefix = 'distributed_rate_limit:'
  ) {
    this.client = client;
    this.nodeId = nodeId;
    this.windowSizeMs = windowSizeMs;
    this.maxRequests = maxRequests;
    this.keyPrefix = keyPrefix;
    this.nodeWeight = nodeWeight;
  }

  private getKey(identifier: string, suffix = ''): string {
    return `${this.keyPrefix}${identifier}${suffix}`;
  }

  async registerNode(): Promise<void> {
    const nodeKey = this.getKey('nodes');
    const nodeData = {
      id: this.nodeId,
      weight: this.nodeWeight,
      lastSeen: Date.now()
    };
    
    await this.client.hSet(nodeKey, this.nodeId, JSON.stringify(nodeData));
    await this.client.expire(nodeKey, Math.ceil(this.windowSizeMs / 1000) * 2);
  }

  async getActiveNodes(): Promise<DistributedNode[]> {
    const nodeKey = this.getKey('nodes');
    const nodes = await this.client.hGetAll(nodeKey);
    const now = Date.now();
    const activeNodes: DistributedNode[] = [];
    
    for (const [nodeId, nodeDataStr] of Object.entries(nodes)) {
      try {
        const nodeData = JSON.parse(nodeDataStr) as DistributedNode;
        
        // Consider node active if seen within 2 window periods
        if (now - nodeData.lastSeen < this.windowSizeMs * 2) {
          activeNodes.push(nodeData);
        }
      } catch (error) {
        console.error('Error parsing node data:', error);
      }
    }
    
    return activeNodes;
  }

  async checkLimit(identifier: string): Promise<RateLimitResult & {
    nodeAllocation: number;
    totalNodes: number;
  }> {
    // Register this node
    await this.registerNode();
    
    // Get active nodes
    const activeNodes = await this.getActiveNodes();
    const totalWeight = activeNodes.reduce((sum, node) => sum + node.weight, 0);
    const nodeAllocation = Math.floor((this.nodeWeight / totalWeight) * this.maxRequests);
    
    // Use fixed window limiter with node allocation
    const nodeKey = this.getKey(identifier, `:node:${this.nodeId}`);
    const window = Math.floor(Date.now() / this.windowSizeMs);
    const windowKey = `${nodeKey}:${window}`;
    
    const luaScript = `
      local key = KEYS[1]
      local limit = tonumber(ARGV[1])
      local window_ttl = tonumber(ARGV[2])
      
      local current = redis.call('GET', key)
      if current == false then
        current = 0
      else
        current = tonumber(current)
      end
      
      if current < limit then
        local new_val = redis.call('INCR', key)
        if new_val == 1 then
          redis.call('EXPIRE', key, window_ttl)
        end
        return {1, limit - new_val, new_val}
      else
        return {0, 0, current}
      end
    `;

    const result = await this.client.eval(luaScript, {
      keys: [windowKey],
      arguments: [
        nodeAllocation.toString(),
        Math.ceil(this.windowSizeMs / 1000).toString()
      ]
    }) as number[];

    const windowStart = window * this.windowSizeMs;
    const resetTime = windowStart + this.windowSizeMs;

    return {
      allowed: result[0] === 1,
      remaining: result[1],
      resetTime,
      totalRequests: result[2],
      nodeAllocation,
      totalNodes: activeNodes.length
    };
  }

  async getGlobalUsage(identifier: string): Promise<{
    totalUsage: number;
    nodeUsages: Array<{ nodeId: string; usage: number }>;
  }> {
    const window = Math.floor(Date.now() / this.windowSizeMs);
    const pattern = this.getKey(identifier, `:node:*:${window}`);
    const keys = await this.client.keys(pattern);
    
    let totalUsage = 0;
    const nodeUsages: Array<{ nodeId: string; usage: number }> = [];
    
    for (const key of keys) {
      const usage = await this.client.get(key);
      const usageCount = usage ? parseInt(usage) : 0;
      totalUsage += usageCount;
      
      // Extract node ID from key
      const match = key.match(/:node:([^:]+):/);
      if (match) {
        nodeUsages.push({
          nodeId: match[1],
          usage: usageCount
        });
      }
    }
    
    return { totalUsage, nodeUsages };
  }
}

export { DistributedConsensusRateLimiter, DistributedNode };
```

## Security Features

### IP-based Security Manager

```typescript
// src/security/ip-security-manager.ts
interface SecurityRule {
  type: 'whitelist' | 'blacklist' | 'rate_limit' | 'geo_block';
  value: string;
  expiresAt?: number;
  reason?: string;
}

interface SecurityEvent {
  ip: string;
  event: string;
  timestamp: number;
  severity: 'low' | 'medium' | 'high' | 'critical';
  metadata?: Record<string, any>;
}

class IPSecurityManager {
  private client: RedisClientType;
  private keyPrefix: string;

  constructor(client: RedisClientType, keyPrefix = 'security:') {
    this.client = client;
    this.keyPrefix = keyPrefix;
  }

  private getKey(suffix: string): string {
    return `${this.keyPrefix}${suffix}`;
  }

  async addSecurityRule(ip: string, rule: SecurityRule): Promise<void> {
    const key = this.getKey(`rules:${ip}`);
    const ruleData = {
      ...rule,
      createdAt: Date.now()
    };
    
    await this.client.hSet(key, rule.type, JSON.stringify(ruleData));
    
    if (rule.expiresAt) {
      const ttl = Math.ceil((rule.expiresAt - Date.now()) / 1000);
      await this.client.expire(key, ttl);
    }
  }

  async removeSecurityRule(ip: string, ruleType: string): Promise<boolean> {
    const key = this.getKey(`rules:${ip}`);
    const result = await this.client.hDel(key, ruleType);
    return result > 0;
  }

  async getSecurityRules(ip: string): Promise<SecurityRule[]> {
    const key = this.getKey(`rules:${ip}`);
    const rules = await this.client.hGetAll(key);
    const now = Date.now();
    
    const activeRules: SecurityRule[] = [];
    
    for (const [type, ruleDataStr] of Object.entries(rules)) {
      try {
        const ruleData = JSON.parse(ruleDataStr);
        
        // Check if rule is still valid
        if (!ruleData.expiresAt || ruleData.expiresAt > now) {
          activeRules.push({
            type: type as SecurityRule['type'],
            value: ruleData.value,
            expiresAt: ruleData.expiresAt,
            reason: ruleData.reason
          });
        } else {
          // Remove expired rule
          await this.client.hDel(key, type);
        }
      } catch (error) {
        console.error('Error parsing security rule:', error);
      }
    }
    
    return activeRules;
  }

  async isIPBlocked(ip: string): Promise<{
    blocked: boolean;
    reason?: string;
    rule?: SecurityRule;
  }> {
    const rules = await this.getSecurityRules(ip);
    
    // Check blacklist
    const blacklistRule = rules.find(rule => rule.type === 'blacklist');
    if (blacklistRule) {
      return {
        blocked: true,
        reason: blacklistRule.reason || 'IP is blacklisted',
        rule: blacklistRule
      };
    }
    
    // Check if IP is in whitelist (if whitelist exists, only whitelisted IPs are allowed)
    const whitelistRules = await this.getWhitelistRules();
    if (whitelistRules.length > 0) {
      const isWhitelisted = rules.some(rule => rule.type === 'whitelist');
      if (!isWhitelisted) {
        return {
          blocked: true,
          reason: 'IP not in whitelist'
        };
      }
    }
    
    return { blocked: false };
  }

  async getWhitelistRules(): Promise<SecurityRule[]> {
    const pattern = this.getKey('rules:*');
    const keys = await this.client.keys(pattern);
    const whitelistRules: SecurityRule[] = [];
    
    for (const key of keys) {
      const whitelistRule = await this.client.hGet(key, 'whitelist');
      if (whitelistRule) {
        try {
          const ruleData = JSON.parse(whitelistRule);
          whitelistRules.push(ruleData);
        } catch (error) {
          console.error('Error parsing whitelist rule:', error);
        }
      }
    }
    
    return whitelistRules;
  }

  async logSecurityEvent(event: SecurityEvent): Promise<void> {
    const eventKey = this.getKey(`events:${event.ip}`);
    const globalEventKey = this.getKey('events:global');
    
    const eventData = {
      ...event,
      id: crypto.randomUUID()
    };
    
    // Store in IP-specific log
    await this.client.lPush(eventKey, JSON.stringify(eventData));
    await this.client.lTrim(eventKey, 0, 999); // Keep last 1000 events
    await this.client.expire(eventKey, 86400 * 7); // 7 days
    
    // Store in global log
    await this.client.lPush(globalEventKey, JSON.stringify(eventData));
    await this.client.lTrim(globalEventKey, 0, 9999); // Keep last 10000 events
    await this.client.expire(globalEventKey, 86400 * 7); // 7 days
    
    // Update event counters
    const counterKey = this.getKey(`counters:${event.event}:${event.severity}`);
    await this.client.incr(counterKey);
    await this.client.expire(counterKey, 86400); // 24 hours
  }

  async getSecurityEvents(
    ip?: string,
    limit = 100
  ): Promise<SecurityEvent[]> {
    const key = ip 
      ? this.getKey(`events:${ip}`)
      : this.getKey('events:global');
    
    const events = await this.client.lRange(key, 0, limit - 1);
    
    return events.map(eventStr => {
      try {
        return JSON.parse(eventStr) as SecurityEvent;
      } catch (error) {
        console.error('Error parsing security event:', error);
        return null;
      }
    }).filter(Boolean) as SecurityEvent[];
  }

  async getEventCounters(): Promise<Record<string, number>> {
    const pattern = this.getKey('counters:*');
    const keys = await this.client.keys(pattern);
    const counters: Record<string, number> = {};
    
    for (const key of keys) {
      const count = await this.client.get(key);
      const keyParts = key.replace(this.keyPrefix, '').split(':');
      const eventName = keyParts.slice(1, -1).join(':');
      const severity = keyParts[keyParts.length - 1];
      
      counters[`${eventName}:${severity}`] = parseInt(count || '0');
    }
    
    return counters;
  }

  async autoBlockSuspiciousIP(
    ip: string,
    threshold: number = 10,
    timeWindowMs: number = 300000 // 5 minutes
  ): Promise<boolean> {
    const events = await this.getSecurityEvents(ip, 1000);
    const now = Date.now();
    
    // Count high severity events in time window
    const recentHighSeverityEvents = events.filter(event => 
      (event.severity === 'high' || event.severity === 'critical') &&
      (now - event.timestamp) < timeWindowMs
    );
    
    if (recentHighSeverityEvents.length >= threshold) {
      await this.addSecurityRule(ip, {
        type: 'blacklist',
        value: ip,
        reason: `Auto-blocked: ${recentHighSeverityEvents.length} high severity events in ${timeWindowMs}ms`,
        expiresAt: now + 3600000 // 1 hour
      });
      
      await this.logSecurityEvent({
        ip,
        event: 'auto_block',
        timestamp: now,
        severity: 'critical',
        metadata: {
          eventCount: recentHighSeverityEvents.length,
          timeWindow: timeWindowMs,
          threshold
        }
      });
      
      return true;
    }
    
    return false;
  }
}

export { IPSecurityManager, SecurityRule, SecurityEvent };
```

## API Protection

### Comprehensive API Rate Limiter

```typescript
// src/api-protection/api-rate-limiter.ts
import { Request, Response, NextFunction } from 'express';

interface APIRateLimitConfig {
  global?: {
    windowMs: number;
    maxRequests: number;
  };
  perUser?: {
    windowMs: number;
    maxRequests: number;
  };
  perIP?: {
    windowMs: number;
    maxRequests: number;
  };
  perEndpoint?: Record<string, {
    windowMs: number;
    maxRequests: number;
  }>;
  skipSuccessfulRequests?: boolean;
  skipFailedRequests?: boolean;
  keyGenerator?: (req: Request) => string;
  onLimitReached?: (req: Request, res: Response) => void;
}

class APIRateLimiter {
  private client: RedisClientType;
  private config: APIRateLimitConfig;
  private limiters: Map<string, FixedWindowRateLimiter>;

  constructor(client: RedisClientType, config: APIRateLimitConfig) {
    this.client = client;
    this.config = config;
    this.limiters = new Map();
    
    this.initializeLimiters();
  }

  private initializeLimiters(): void {
    if (this.config.global) {
      this.limiters.set('global', new FixedWindowRateLimiter(this.client, {
        windowSizeMs: this.config.global.windowMs,
        maxRequests: this.config.global.maxRequests,
        keyPrefix: 'api:global:'
      }));
    }
    
    if (this.config.perUser) {
      this.limiters.set('perUser', new FixedWindowRateLimiter(this.client, {
        windowSizeMs: this.config.perUser.windowMs,
        maxRequests: this.config.perUser.maxRequests,
        keyPrefix: 'api:user:'
      }));
    }
    
    if (this.config.perIP) {
      this.limiters.set('perIP', new FixedWindowRateLimiter(this.client, {
        windowSizeMs: this.config.perIP.windowMs,
        maxRequests: this.config.perIP.maxRequests,
        keyPrefix: 'api:ip:'
      }));
    }
    
    if (this.config.perEndpoint) {
      for (const [endpoint, config] of Object.entries(this.config.perEndpoint)) {
        this.limiters.set(`endpoint:${endpoint}`, new FixedWindowRateLimiter(this.client, {
          windowSizeMs: config.windowMs,
          maxRequests: config.maxRequests,
          keyPrefix: `api:endpoint:${endpoint}:`
        }));
      }
    }
  }

  middleware() {
    return async (req: Request, res: Response, next: NextFunction) => {
      try {
        const checks = await this.performRateLimitChecks(req);
        
        // Find the most restrictive limit
        const blockedCheck = checks.find(check => !check.allowed);
        
        if (blockedCheck) {
          // Set rate limit headers
          res.set({
            'X-RateLimit-Limit': blockedCheck.limit.toString(),
            'X-RateLimit-Remaining': blockedCheck.remaining.toString(),
            'X-RateLimit-Reset': new Date(blockedCheck.resetTime).toISOString()
          });
          
          if (this.config.onLimitReached) {
            this.config.onLimitReached(req, res);
          } else {
            res.status(429).json({
              error: 'Too Many Requests',
              message: `Rate limit exceeded for ${blockedCheck.type}`,
              retryAfter: Math.ceil((blockedCheck.resetTime - Date.now()) / 1000)
            });
          }
          
          return;
        }
        
        // Set headers for successful requests
        const mostRestrictive = checks.reduce((min, check) => 
          check.remaining < min.remaining ? check : min
        );
        
        res.set({
          'X-RateLimit-Limit': mostRestrictive.limit.toString(),
          'X-RateLimit-Remaining': mostRestrictive.remaining.toString(),
          'X-RateLimit-Reset': new Date(mostRestrictive.resetTime).toISOString()
        });
        
        next();
      } catch (error) {
        console.error('Rate limiting error:', error);
        next(); // Continue on error
      }
    };
  }

  private async performRateLimitChecks(req: Request): Promise<Array<{
    type: string;
    allowed: boolean;
    remaining: number;
    resetTime: number;
    limit: number;
  }>> {
    const checks = [];
    
    // Global rate limit
    if (this.limiters.has('global')) {
      const limiter = this.limiters.get('global')!;
      const result = await limiter.checkLimit('global');
      checks.push({
        type: 'global',
        allowed: result.allowed,
        remaining: result.remaining,
        resetTime: result.resetTime,
        limit: this.config.global!.maxRequests
      });
    }
    
    // Per-IP rate limit
    if (this.limiters.has('perIP')) {
      const ip = this.getClientIP(req);
      const limiter = this.limiters.get('perIP')!;
      const result = await limiter.checkLimit(ip);
      checks.push({
        type: 'ip',
        allowed: result.allowed,
        remaining: result.remaining,
        resetTime: result.resetTime,
        limit: this.config.perIP!.maxRequests
      });
    }
    
    // Per-user rate limit
    if (this.limiters.has('perUser') && req.user?.id) {
      const limiter = this.limiters.get('perUser')!;
      const result = await limiter.checkLimit(req.user.id);
      checks.push({
        type: 'user',
        allowed: result.allowed,
        remaining: result.remaining,
        resetTime: result.resetTime,
        limit: this.config.perUser!.maxRequests
      });
    }
    
    // Per-endpoint rate limit
    const endpoint = this.getEndpointKey(req);
    if (this.limiters.has(`endpoint:${endpoint}`)) {
      const identifier = this.config.keyGenerator ? 
        this.config.keyGenerator(req) : 
        this.getClientIP(req);
      
      const limiter = this.limiters.get(`endpoint:${endpoint}`)!;
      const result = await limiter.checkLimit(identifier);
      checks.push({
        type: `endpoint:${endpoint}`,
        allowed: result.allowed,
        remaining: result.remaining,
        resetTime: result.resetTime,
        limit: this.config.perEndpoint![endpoint].maxRequests
      });
    }
    
    return checks;
  }

  private getClientIP(req: Request): string {
    return (req.headers['x-forwarded-for'] as string)?.split(',')[0] ||
           req.connection.remoteAddress ||
           req.socket.remoteAddress ||
           'unknown';
  }

  private getEndpointKey(req: Request): string {
    return `${req.method}:${req.route?.path || req.path}`;
  }

  async getUsageStats(): Promise<Record<string, any>> {
    const stats: Record<string, any> = {};
    
    for (const [key, limiter] of this.limiters) {
      // This would require additional methods in the rate limiter classes
      // to get usage statistics
      stats[key] = {
        type: key,
        // Add usage statistics here
      };
    }
    
    return stats;
  }
}

// Extend Express Request type
declare global {
  namespace Express {
    interface Request {
      user?: {
        id: string;
        [key: string]: any;
      };
    }
  }
}

export { APIRateLimiter, APIRateLimitConfig };
```

## Usage Examples

```typescript
// src/examples/rate-limiting-usage.ts
import { createClient } from 'redis';
import express from 'express';
import { FixedWindowRateLimiter } from '../rate-limiting/fixed-window-limiter';
import { SlidingWindowRateLimiter } from '../rate-limiting/sliding-window-limiter';
import { TokenBucketRateLimiter } from '../rate-limiting/token-bucket-limiter';
import { MultiTierRateLimiter } from '../rate-limiting/multi-tier-limiter';
import { IPSecurityManager } from '../security/ip-security-manager';
import { APIRateLimiter } from '../api-protection/api-rate-limiter';

async function runRateLimitingExamples() {
  const client = createClient({
    host: 'localhost',
    port: 6379
  });
  
  await client.connect();
  
  try {
    // Fixed Window Rate Limiter
    console.log('=== Fixed Window Rate Limiter ===');
    const fixedLimiter = new FixedWindowRateLimiter(client, {
      windowSizeMs: 60000, // 1 minute
      maxRequests: 100
    });
    
    const result1 = await fixedLimiter.checkLimit('user123');
    console.log('Fixed window result:', result1);
    
    // Sliding Window Rate Limiter
    console.log('\n=== Sliding Window Rate Limiter ===');
    const slidingLimiter = new SlidingWindowRateLimiter(
      client,
      60000, // 1 minute window
      50     // 50 requests
    );
    
    const result2 = await slidingLimiter.checkLimit('user456');
    console.log('Sliding window result:', result2);
    
    // Token Bucket Rate Limiter
    console.log('\n=== Token Bucket Rate Limiter ===');
    const tokenBucket = new TokenBucketRateLimiter(client, {
      capacity: 100,
      refillRate: 10, // 10 tokens per second
      refillPeriodMs: 1000
    });
    
    const result3 = await tokenBucket.checkLimit('user789', 5);
    console.log('Token bucket result:', result3);
    
    // Multi-Tier Rate Limiter
    console.log('\n=== Multi-Tier Rate Limiter ===');
    const multiTier = new MultiTierRateLimiter(client, {
      tiers: [
        { name: 'second', windowSizeMs: 1000, maxRequests: 10, priority: 1 },
        { name: 'minute', windowSizeMs: 60000, maxRequests: 100, priority: 2 },
        { name: 'hour', windowSizeMs: 3600000, maxRequests: 1000, priority: 3 }
      ]
    });
    
    const result4 = await multiTier.checkLimit('user_multi');
    console.log('Multi-tier result:', result4);
    
    // IP Security Manager
    console.log('\n=== IP Security Manager ===');
    const securityManager = new IPSecurityManager(client);
    
    // Add a blacklist rule
    await securityManager.addSecurityRule('192.168.1.100', {
      type: 'blacklist',
      value: '192.168.1.100',
      reason: 'Suspicious activity detected',
      expiresAt: Date.now() + 3600000 // 1 hour
    });
    
    // Check if IP is blocked
    const blockCheck = await securityManager.isIPBlocked('192.168.1.100');
    console.log('IP block check:', blockCheck);
    
    // Log security event
    await securityManager.logSecurityEvent({
      ip: '192.168.1.101',
      event: 'failed_login',
      timestamp: Date.now(),
      severity: 'medium',
      metadata: { attempts: 3 }
    });
    
    // Express API Rate Limiter
    console.log('\n=== Express API Rate Limiter ===');
    const app = express();
    
    const apiLimiter = new APIRateLimiter(client, {
      global: {
        windowMs: 60000,
        maxRequests: 1000
      },
      perIP: {
        windowMs: 60000,
        maxRequests: 100
      },
      perUser: {
        windowMs: 60000,
        maxRequests: 200
      },
      perEndpoint: {
        'POST:/api/login': {
          windowMs: 300000, // 5 minutes
          maxRequests: 5
        },
        'GET:/api/data': {
          windowMs: 60000,
          maxRequests: 50
        }
      }
    });
    
    app.use(apiLimiter.middleware());
    
    app.get('/api/data', (req, res) => {
      res.json({ message: 'Data retrieved successfully' });
    });
    
    app.post('/api/login', (req, res) => {
      res.json({ message: 'Login attempt' });
    });
    
    console.log('Express app configured with rate limiting');
    
  } finally {
    await client.disconnect();
  }
}

// Performance testing function
async function performanceTest() {
  const client = createClient({
    host: 'localhost',
    port: 6379
  });
  
  await client.connect();
  
  const limiter = new FixedWindowRateLimiter(client, {
    windowSizeMs: 60000,
    maxRequests: 1000
  });
  
  const startTime = Date.now();
  const promises = [];
  
  // Test 1000 concurrent requests
  for (let i = 0; i < 1000; i++) {
    promises.push(limiter.checkLimit(`user${i % 100}`));
  }
  
  const results = await Promise.all(promises);
  const endTime = Date.now();
  
  const allowed = results.filter(r => r.allowed).length;
  const denied = results.filter(r => !r.allowed).length;
  
  console.log('Performance Test Results:');
  console.log(`Total requests: ${results.length}`);
  console.log(`Allowed: ${allowed}`);
  console.log(`Denied: ${denied}`);
  console.log(`Time taken: ${endTime - startTime}ms`);
  console.log(`Requests per second: ${(results.length / (endTime - startTime) * 1000).toFixed(2)}`);
  
  await client.disconnect();
}

// Run examples
if (require.main === module) {
  runRateLimitingExamples()
    .then(() => performanceTest())
    .catch(console.error);
}

export { runRateLimitingExamples, performanceTest };
```

## Best Practices

### 1. Algorithm Selection

- **Fixed Window**: Simple, memory efficient, but allows bursts
- **Sliding Window**: More accurate, prevents burst issues
- **Token Bucket**: Good for APIs that need burst handling
- **Leaky Bucket**: Smooth rate enforcement, good for streaming

### 2. Key Design

- Use hierarchical keys for easy management
- Include time windows in keys for automatic cleanup
- Consider using hashes for complex rate limit data

### 3. Performance

- Use Lua scripts for atomic operations
- Implement proper key expiration
- Monitor Redis memory usage
- Consider using Redis clustering for high load

### 4. Security

- Implement multiple layers of protection
- Log security events for analysis
- Use automatic blocking for suspicious activity
- Regularly review and update security rules

## Next Steps

Now that you understand rate limiting and security, you're ready to explore:

- [Part 13: Real-time Applications](./part-13-realtime-applications.md)

## Quick Reference

### Basic Rate Limiting
```typescript
// Fixed window
const limiter = new FixedWindowRateLimiter(client, {
  windowSizeMs: 60000,
  maxRequests: 100
});

const result = await limiter.checkLimit('user123');
if (!result.allowed) {
  // Rate limit exceeded
}
```

### Security Rules
```typescript
// Block IP
await securityManager.addSecurityRule('192.168.1.100', {
  type: 'blacklist',
  value: '192.168.1.100',
  reason: 'Suspicious activity'
});

// Check if blocked
const blocked = await securityManager.isIPBlocked('192.168.1.100');
```

### Express Integration
```typescript
// API rate limiter middleware
app.use(new APIRateLimiter(client, {
  perIP: { windowMs: 60000, maxRequests: 100 },
  perUser: { windowMs: 60000, maxRequests: 200 }
}).middleware());
```

---

**Previous**: [← Part 11: Session Management](./part-11-session-management.md)  
**Next**: [Part 13: Real-time Applications →](./part-13-realtime-applications.md)