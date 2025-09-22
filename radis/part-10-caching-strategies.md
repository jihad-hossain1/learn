# Part 10: Caching Strategies with Redis

Caching is one of Redis's most popular use cases. This part covers various caching strategies, patterns, and implementations in Node.js.

## Table of Contents

1. [Caching Fundamentals](#caching-fundamentals)
2. [Cache Patterns](#cache-patterns)
3. [Cache-Aside Pattern](#cache-aside-pattern)
4. [Write-Through Cache](#write-through-cache)
5. [Write-Behind Cache](#write-behind-cache)
6. [Read-Through Cache](#read-through-cache)
7. [Multi-Level Caching](#multi-level-caching)
8. [Cache Invalidation](#cache-invalidation)
9. [Performance Optimization](#performance-optimization)
10. [Real-world Examples](#real-world-examples)

## Caching Fundamentals

### Why Cache?

- **Performance**: Faster data access
- **Scalability**: Reduced database load
- **Cost**: Lower infrastructure costs
- **User Experience**: Improved response times

### Cache Considerations

- **Cache Hit Ratio**: Percentage of requests served from cache
- **TTL (Time To Live)**: How long data stays in cache
- **Cache Size**: Memory limitations
- **Consistency**: Data freshness vs performance

## Cache Patterns

### Overview of Common Patterns

1. **Cache-Aside (Lazy Loading)**: Application manages cache
2. **Write-Through**: Write to cache and database simultaneously
3. **Write-Behind (Write-Back)**: Write to cache first, database later
4. **Read-Through**: Cache loads data on miss
5. **Refresh-Ahead**: Proactively refresh cache before expiration

## Cache-Aside Pattern

The most common caching pattern where the application manages the cache.

```typescript
// src/caching/cache-aside.ts
import { RedisClientType } from 'redis';

interface CacheOptions {
  ttl?: number;
  prefix?: string;
}

class CacheAside {
  private client: RedisClientType;
  private defaultTTL: number;
  private keyPrefix: string;

  constructor(client: RedisClientType, options: CacheOptions = {}) {
    this.client = client;
    this.defaultTTL = options.ttl || 3600; // 1 hour default
    this.keyPrefix = options.prefix || 'cache:';
  }

  private buildKey(key: string): string {
    return `${this.keyPrefix}${key}`;
  }

  async get<T>(key: string): Promise<T | null> {
    try {
      const cached = await this.client.get(this.buildKey(key));
      if (cached) {
        return JSON.parse(cached);
      }
      return null;
    } catch (error) {
      console.error('Cache get error:', error);
      return null;
    }
  }

  async set<T>(key: string, value: T, ttl?: number): Promise<boolean> {
    try {
      const cacheKey = this.buildKey(key);
      const serialized = JSON.stringify(value);
      const expiration = ttl || this.defaultTTL;
      
      await this.client.setEx(cacheKey, expiration, serialized);
      return true;
    } catch (error) {
      console.error('Cache set error:', error);
      return false;
    }
  }

  async delete(key: string): Promise<boolean> {
    try {
      const result = await this.client.del(this.buildKey(key));
      return result > 0;
    } catch (error) {
      console.error('Cache delete error:', error);
      return false;
    }
  }

  async getOrSet<T>(
    key: string,
    fetcher: () => Promise<T>,
    ttl?: number
  ): Promise<T> {
    // Try to get from cache first
    const cached = await this.get<T>(key);
    if (cached !== null) {
      return cached;
    }

    // Cache miss - fetch from source
    const value = await fetcher();
    
    // Store in cache for next time
    await this.set(key, value, ttl);
    
    return value;
  }

  async mget<T>(keys: string[]): Promise<(T | null)[]> {
    try {
      const cacheKeys = keys.map(key => this.buildKey(key));
      const results = await this.client.mGet(cacheKeys);
      
      return results.map(result => {
        if (result) {
          try {
            return JSON.parse(result);
          } catch {
            return null;
          }
        }
        return null;
      });
    } catch (error) {
      console.error('Cache mget error:', error);
      return keys.map(() => null);
    }
  }

  async mset<T>(items: Array<{ key: string; value: T; ttl?: number }>): Promise<boolean> {
    try {
      const pipeline = this.client.multi();
      
      items.forEach(({ key, value, ttl }) => {
        const cacheKey = this.buildKey(key);
        const serialized = JSON.stringify(value);
        const expiration = ttl || this.defaultTTL;
        
        pipeline.setEx(cacheKey, expiration, serialized);
      });
      
      await pipeline.exec();
      return true;
    } catch (error) {
      console.error('Cache mset error:', error);
      return false;
    }
  }
}

export { CacheAside, CacheOptions };
```

### User Service with Cache-Aside

```typescript
// src/caching/user-service.ts
interface User {
  id: string;
  name: string;
  email: string;
  createdAt: Date;
  updatedAt: Date;
}

interface UserRepository {
  findById(id: string): Promise<User | null>;
  findByEmail(email: string): Promise<User | null>;
  create(user: Omit<User, 'id' | 'createdAt' | 'updatedAt'>): Promise<User>;
  update(id: string, updates: Partial<User>): Promise<User | null>;
  delete(id: string): Promise<boolean>;
}

class CachedUserService {
  private cache: CacheAside;
  private repository: UserRepository;
  private userTTL: number = 1800; // 30 minutes

  constructor(cache: CacheAside, repository: UserRepository) {
    this.cache = cache;
    this.repository = repository;
  }

  async getUserById(id: string): Promise<User | null> {
    const cacheKey = `user:${id}`;
    
    return await this.cache.getOrSet(
      cacheKey,
      () => this.repository.findById(id),
      this.userTTL
    );
  }

  async getUserByEmail(email: string): Promise<User | null> {
    const cacheKey = `user:email:${email}`;
    
    return await this.cache.getOrSet(
      cacheKey,
      () => this.repository.findByEmail(email),
      this.userTTL
    );
  }

  async createUser(userData: Omit<User, 'id' | 'createdAt' | 'updatedAt'>): Promise<User> {
    const user = await this.repository.create(userData);
    
    // Cache the new user
    await this.cache.set(`user:${user.id}`, user, this.userTTL);
    await this.cache.set(`user:email:${user.email}`, user, this.userTTL);
    
    return user;
  }

  async updateUser(id: string, updates: Partial<User>): Promise<User | null> {
    const user = await this.repository.update(id, updates);
    
    if (user) {
      // Update cache
      await this.cache.set(`user:${user.id}`, user, this.userTTL);
      await this.cache.set(`user:email:${user.email}`, user, this.userTTL);
      
      // If email changed, invalidate old email cache
      if (updates.email) {
        const oldUser = await this.cache.get<User>(`user:${id}`);
        if (oldUser && oldUser.email !== updates.email) {
          await this.cache.delete(`user:email:${oldUser.email}`);
        }
      }
    }
    
    return user;
  }

  async deleteUser(id: string): Promise<boolean> {
    // Get user first to invalidate email cache
    const user = await this.cache.get<User>(`user:${id}`);
    
    const deleted = await this.repository.delete(id);
    
    if (deleted) {
      // Invalidate cache
      await this.cache.delete(`user:${id}`);
      if (user) {
        await this.cache.delete(`user:email:${user.email}`);
      }
    }
    
    return deleted;
  }

  async getUsersBatch(ids: string[]): Promise<(User | null)[]> {
    const cacheKeys = ids.map(id => `user:${id}`);
    const cached = await this.cache.mget<User>(cacheKeys);
    
    const missingIds: string[] = [];
    const missingIndices: number[] = [];
    
    cached.forEach((user, index) => {
      if (user === null) {
        missingIds.push(ids[index]);
        missingIndices.push(index);
      }
    });
    
    if (missingIds.length > 0) {
      // Fetch missing users from database
      const missingUsers = await Promise.all(
        missingIds.map(id => this.repository.findById(id))
      );
      
      // Cache the fetched users
      const cacheItems = missingUsers
        .filter(user => user !== null)
        .map(user => ({
          key: `user:${user!.id}`,
          value: user!,
          ttl: this.userTTL
        }));
      
      if (cacheItems.length > 0) {
        await this.cache.mset(cacheItems);
      }
      
      // Fill in the missing users
      missingIndices.forEach((index, i) => {
        cached[index] = missingUsers[i];
      });
    }
    
    return cached;
  }
}

export { CachedUserService, User, UserRepository };
```

## Write-Through Cache

Writes go to both cache and database simultaneously.

```typescript
// src/caching/write-through.ts
class WriteThroughCache {
  private cache: CacheAside;
  private repository: UserRepository;

  constructor(cache: CacheAside, repository: UserRepository) {
    this.cache = cache;
    this.repository = repository;
  }

  async createUser(userData: Omit<User, 'id' | 'createdAt' | 'updatedAt'>): Promise<User> {
    // Write to database first
    const user = await this.repository.create(userData);
    
    // Then write to cache
    await this.cache.set(`user:${user.id}`, user);
    await this.cache.set(`user:email:${user.email}`, user);
    
    return user;
  }

  async updateUser(id: string, updates: Partial<User>): Promise<User | null> {
    // Update database first
    const user = await this.repository.update(id, updates);
    
    if (user) {
      // Then update cache
      await this.cache.set(`user:${user.id}`, user);
      await this.cache.set(`user:email:${user.email}`, user);
    }
    
    return user;
  }

  async getUserById(id: string): Promise<User | null> {
    // Try cache first
    const cached = await this.cache.get<User>(`user:${id}`);
    if (cached) {
      return cached;
    }
    
    // If not in cache, load from database and cache it
    const user = await this.repository.findById(id);
    if (user) {
      await this.cache.set(`user:${id}`, user);
    }
    
    return user;
  }
}

export { WriteThroughCache };
```

## Write-Behind Cache

Writes go to cache immediately, database writes are deferred.

```typescript
// src/caching/write-behind.ts
interface WriteOperation {
  type: 'create' | 'update' | 'delete';
  key: string;
  data?: any;
  timestamp: number;
}

class WriteBehindCache {
  private cache: CacheAside;
  private repository: UserRepository;
  private writeQueue: WriteOperation[] = [];
  private flushInterval: number;
  private maxQueueSize: number;
  private flushTimer?: NodeJS.Timeout;

  constructor(
    cache: CacheAside,
    repository: UserRepository,
    options: {
      flushInterval?: number;
      maxQueueSize?: number;
    } = {}
  ) {
    this.cache = cache;
    this.repository = repository;
    this.flushInterval = options.flushInterval || 5000; // 5 seconds
    this.maxQueueSize = options.maxQueueSize || 100;
    
    this.startFlushTimer();
  }

  private startFlushTimer(): void {
    this.flushTimer = setInterval(() => {
      this.flushWrites().catch(console.error);
    }, this.flushInterval);
  }

  private async flushWrites(): Promise<void> {
    if (this.writeQueue.length === 0) {
      return;
    }

    const operations = [...this.writeQueue];
    this.writeQueue = [];

    try {
      // Group operations by type for batch processing
      const creates = operations.filter(op => op.type === 'create');
      const updates = operations.filter(op => op.type === 'update');
      const deletes = operations.filter(op => op.type === 'delete');

      // Process operations
      await Promise.all([
        this.processCreates(creates),
        this.processUpdates(updates),
        this.processDeletes(deletes)
      ]);

      console.log(`Flushed ${operations.length} write operations`);
    } catch (error) {
      console.error('Error flushing writes:', error);
      // Re-queue failed operations
      this.writeQueue.unshift(...operations);
    }
  }

  private async processCreates(operations: WriteOperation[]): Promise<void> {
    for (const op of operations) {
      try {
        await this.repository.create(op.data);
      } catch (error) {
        console.error(`Failed to create user ${op.key}:`, error);
      }
    }
  }

  private async processUpdates(operations: WriteOperation[]): Promise<void> {
    for (const op of operations) {
      try {
        const userId = op.key.replace('user:', '');
        await this.repository.update(userId, op.data);
      } catch (error) {
        console.error(`Failed to update user ${op.key}:`, error);
      }
    }
  }

  private async processDeletes(operations: WriteOperation[]): Promise<void> {
    for (const op of operations) {
      try {
        const userId = op.key.replace('user:', '');
        await this.repository.delete(userId);
      } catch (error) {
        console.error(`Failed to delete user ${op.key}:`, error);
      }
    }
  }

  private queueWrite(operation: WriteOperation): void {
    this.writeQueue.push(operation);
    
    // Flush immediately if queue is full
    if (this.writeQueue.length >= this.maxQueueSize) {
      this.flushWrites().catch(console.error);
    }
  }

  async createUser(userData: Omit<User, 'id' | 'createdAt' | 'updatedAt'>): Promise<User> {
    // Create user object with generated ID
    const user: User = {
      id: `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      ...userData,
      createdAt: new Date(),
      updatedAt: new Date()
    };

    // Write to cache immediately
    await this.cache.set(`user:${user.id}`, user);
    await this.cache.set(`user:email:${user.email}`, user);

    // Queue database write
    this.queueWrite({
      type: 'create',
      key: `user:${user.id}`,
      data: user,
      timestamp: Date.now()
    });

    return user;
  }

  async updateUser(id: string, updates: Partial<User>): Promise<User | null> {
    // Get current user from cache
    const currentUser = await this.cache.get<User>(`user:${id}`);
    if (!currentUser) {
      return null;
    }

    // Create updated user
    const updatedUser: User = {
      ...currentUser,
      ...updates,
      updatedAt: new Date()
    };

    // Update cache immediately
    await this.cache.set(`user:${id}`, updatedUser);
    await this.cache.set(`user:email:${updatedUser.email}`, updatedUser);

    // Queue database write
    this.queueWrite({
      type: 'update',
      key: `user:${id}`,
      data: updates,
      timestamp: Date.now()
    });

    return updatedUser;
  }

  async deleteUser(id: string): Promise<boolean> {
    // Remove from cache immediately
    const user = await this.cache.get<User>(`user:${id}`);
    await this.cache.delete(`user:${id}`);
    if (user) {
      await this.cache.delete(`user:email:${user.email}`);
    }

    // Queue database delete
    this.queueWrite({
      type: 'delete',
      key: `user:${id}`,
      timestamp: Date.now()
    });

    return true;
  }

  async getUserById(id: string): Promise<User | null> {
    return await this.cache.get<User>(`user:${id}`);
  }

  async shutdown(): Promise<void> {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
    }
    await this.flushWrites();
  }
}

export { WriteBehindCache, WriteOperation };
```

## Read-Through Cache

Cache automatically loads data on cache miss.

```typescript
// src/caching/read-through.ts
class ReadThroughCache {
  private cache: CacheAside;
  private repository: UserRepository;
  private loadingKeys: Set<string> = new Set();

  constructor(cache: CacheAside, repository: UserRepository) {
    this.cache = cache;
    this.repository = repository;
  }

  async getUserById(id: string): Promise<User | null> {
    const cacheKey = `user:${id}`;
    
    // Try cache first
    const cached = await this.cache.get<User>(cacheKey);
    if (cached !== null) {
      return cached;
    }

    // Prevent multiple concurrent loads of the same key
    if (this.loadingKeys.has(cacheKey)) {
      // Wait a bit and try cache again
      await new Promise(resolve => setTimeout(resolve, 50));
      return await this.cache.get<User>(cacheKey);
    }

    this.loadingKeys.add(cacheKey);

    try {
      // Load from database
      const user = await this.repository.findById(id);
      
      if (user) {
        // Store in cache
        await this.cache.set(cacheKey, user);
        await this.cache.set(`user:email:${user.email}`, user);
      }
      
      return user;
    } finally {
      this.loadingKeys.delete(cacheKey);
    }
  }

  async getUserByEmail(email: string): Promise<User | null> {
    const cacheKey = `user:email:${email}`;
    
    // Try cache first
    const cached = await this.cache.get<User>(cacheKey);
    if (cached !== null) {
      return cached;
    }

    // Prevent concurrent loads
    if (this.loadingKeys.has(cacheKey)) {
      await new Promise(resolve => setTimeout(resolve, 50));
      return await this.cache.get<User>(cacheKey);
    }

    this.loadingKeys.add(cacheKey);

    try {
      const user = await this.repository.findByEmail(email);
      
      if (user) {
        await this.cache.set(`user:${user.id}`, user);
        await this.cache.set(cacheKey, user);
      }
      
      return user;
    } finally {
      this.loadingKeys.delete(cacheKey);
    }
  }
}

export { ReadThroughCache };
```

## Multi-Level Caching

Combining multiple cache layers for optimal performance.

```typescript
// src/caching/multi-level.ts
interface L1Cache {
  get<T>(key: string): T | undefined;
  set<T>(key: string, value: T, ttl?: number): void;
  delete(key: string): void;
  clear(): void;
}

class MemoryCache implements L1Cache {
  private cache = new Map<string, { value: any; expires: number }>();
  private defaultTTL: number;

  constructor(defaultTTL: number = 300000) { // 5 minutes
    this.defaultTTL = defaultTTL;
  }

  get<T>(key: string): T | undefined {
    const item = this.cache.get(key);
    if (!item) {
      return undefined;
    }

    if (Date.now() > item.expires) {
      this.cache.delete(key);
      return undefined;
    }

    return item.value;
  }

  set<T>(key: string, value: T, ttl?: number): void {
    const expires = Date.now() + (ttl || this.defaultTTL);
    this.cache.set(key, { value, expires });
  }

  delete(key: string): void {
    this.cache.delete(key);
  }

  clear(): void {
    this.cache.clear();
  }

  size(): number {
    return this.cache.size;
  }
}

class MultiLevelCache {
  private l1Cache: L1Cache; // Memory cache
  private l2Cache: CacheAside; // Redis cache
  private repository: UserRepository;

  constructor(
    l1Cache: L1Cache,
    l2Cache: CacheAside,
    repository: UserRepository
  ) {
    this.l1Cache = l1Cache;
    this.l2Cache = l2Cache;
    this.repository = repository;
  }

  async getUserById(id: string): Promise<User | null> {
    const cacheKey = `user:${id}`;

    // Try L1 cache (memory) first
    let user = this.l1Cache.get<User>(cacheKey);
    if (user) {
      return user;
    }

    // Try L2 cache (Redis)
    user = await this.l2Cache.get<User>(cacheKey);
    if (user) {
      // Populate L1 cache
      this.l1Cache.set(cacheKey, user, 300000); // 5 minutes
      return user;
    }

    // Load from database
    user = await this.repository.findById(id);
    if (user) {
      // Populate both cache levels
      this.l1Cache.set(cacheKey, user, 300000); // 5 minutes
      await this.l2Cache.set(cacheKey, user, 1800); // 30 minutes
    }

    return user;
  }

  async updateUser(id: string, updates: Partial<User>): Promise<User | null> {
    const user = await this.repository.update(id, updates);
    
    if (user) {
      const cacheKey = `user:${id}`;
      const emailKey = `user:email:${user.email}`;
      
      // Update both cache levels
      this.l1Cache.set(cacheKey, user);
      this.l1Cache.set(emailKey, user);
      
      await this.l2Cache.set(cacheKey, user);
      await this.l2Cache.set(emailKey, user);
    }
    
    return user;
  }

  async deleteUser(id: string): Promise<boolean> {
    const user = await this.l2Cache.get<User>(`user:${id}`);
    const deleted = await this.repository.delete(id);
    
    if (deleted) {
      // Invalidate both cache levels
      this.l1Cache.delete(`user:${id}`);
      await this.l2Cache.delete(`user:${id}`);
      
      if (user) {
        this.l1Cache.delete(`user:email:${user.email}`);
        await this.l2Cache.delete(`user:email:${user.email}`);
      }
    }
    
    return deleted;
  }

  async warmupCache(userIds: string[]): Promise<void> {
    const users = await Promise.all(
      userIds.map(id => this.repository.findById(id))
    );

    const validUsers = users.filter(user => user !== null) as User[];
    
    // Warm up both cache levels
    const l2CacheItems = validUsers.map(user => ({
      key: `user:${user.id}`,
      value: user,
      ttl: 1800
    }));
    
    await this.l2Cache.mset(l2CacheItems);
    
    validUsers.forEach(user => {
      this.l1Cache.set(`user:${user.id}`, user);
    });
  }
}

export { MultiLevelCache, MemoryCache, L1Cache };
```

## Cache Invalidation

Strategies for keeping cache data fresh.

```typescript
// src/caching/invalidation.ts
class CacheInvalidation {
  private cache: CacheAside;
  private client: RedisClientType;

  constructor(cache: CacheAside, client: RedisClientType) {
    this.cache = cache;
    this.client = client;
  }

  // Time-based invalidation
  async setWithTTL<T>(key: string, value: T, ttl: number): Promise<void> {
    await this.cache.set(key, value, ttl);
  }

  // Tag-based invalidation
  async setWithTags<T>(
    key: string,
    value: T,
    tags: string[],
    ttl?: number
  ): Promise<void> {
    await this.cache.set(key, value, ttl);
    
    // Associate key with tags
    const pipeline = this.client.multi();
    tags.forEach(tag => {
      pipeline.sAdd(`tag:${tag}`, key);
      if (ttl) {
        pipeline.expire(`tag:${tag}`, ttl + 3600); // Tag expires later
      }
    });
    await pipeline.exec();
  }

  async invalidateByTag(tag: string): Promise<number> {
    const tagKey = `tag:${tag}`;
    const keys = await this.client.sMembers(tagKey);
    
    if (keys.length === 0) {
      return 0;
    }

    // Delete all keys associated with the tag
    const pipeline = this.client.multi();
    keys.forEach(key => pipeline.del(key));
    pipeline.del(tagKey); // Remove the tag set itself
    
    const results = await pipeline.exec();
    return keys.length;
  }

  // Pattern-based invalidation
  async invalidateByPattern(pattern: string): Promise<number> {
    const script = `
      local keys = redis.call('KEYS', ARGV[1])
      local count = 0
      for i = 1, #keys do
        redis.call('DEL', keys[i])
        count = count + 1
      end
      return count
    `;
    
    const result = await this.client.eval(script, {
      keys: [],
      arguments: [pattern]
    });
    
    return result as number;
  }

  // Version-based invalidation
  async setWithVersion<T>(
    key: string,
    value: T,
    version: string,
    ttl?: number
  ): Promise<void> {
    const versionedKey = `${key}:v:${version}`;
    await this.cache.set(versionedKey, value, ttl);
    await this.cache.set(`${key}:current`, version, ttl);
  }

  async getWithVersion<T>(key: string): Promise<T | null> {
    const currentVersion = await this.cache.get<string>(`${key}:current`);
    if (!currentVersion) {
      return null;
    }
    
    const versionedKey = `${key}:v:${currentVersion}`;
    return await this.cache.get<T>(versionedKey);
  }

  async invalidateVersion(key: string, version: string): Promise<void> {
    const versionedKey = `${key}:v:${version}`;
    await this.cache.delete(versionedKey);
  }

  // Event-driven invalidation
  async setupEventInvalidation(): Promise<void> {
    // Subscribe to invalidation events
    const subscriber = this.client.duplicate();
    await subscriber.connect();
    
    await subscriber.subscribe('cache:invalidate', (message) => {
      const { type, key, pattern, tag } = JSON.parse(message);
      
      switch (type) {
        case 'key':
          this.cache.delete(key);
          break;
        case 'pattern':
          this.invalidateByPattern(pattern);
          break;
        case 'tag':
          this.invalidateByTag(tag);
          break;
      }
    });
  }

  async publishInvalidation(type: string, data: any): Promise<void> {
    await this.client.publish('cache:invalidate', JSON.stringify({
      type,
      ...data,
      timestamp: Date.now()
    }));
  }
}

export { CacheInvalidation };
```

## Performance Optimization

### Cache Warming

```typescript
// src/caching/cache-warming.ts
class CacheWarming {
  private cache: CacheAside;
  private repository: UserRepository;

  constructor(cache: CacheAside, repository: UserRepository) {
    this.cache = cache;
    this.repository = repository;
  }

  async warmPopularUsers(limit: number = 100): Promise<void> {
    // This would typically come from analytics
    const popularUserIds = await this.getPopularUserIds(limit);
    
    const batchSize = 10;
    for (let i = 0; i < popularUserIds.length; i += batchSize) {
      const batch = popularUserIds.slice(i, i + batchSize);
      
      const users = await Promise.all(
        batch.map(id => this.repository.findById(id))
      );
      
      const cacheItems = users
        .filter(user => user !== null)
        .map(user => ({
          key: `user:${user!.id}`,
          value: user!,
          ttl: 3600
        }));
      
      if (cacheItems.length > 0) {
        await this.cache.mset(cacheItems);
      }
    }
  }

  async warmUsersByActivity(): Promise<void> {
    // Warm cache based on recent activity
    const recentlyActiveUsers = await this.getRecentlyActiveUsers();
    
    const cacheItems = recentlyActiveUsers.map(user => ({
      key: `user:${user.id}`,
      value: user,
      ttl: 1800 // 30 minutes
    }));
    
    await this.cache.mset(cacheItems);
  }

  private async getPopularUserIds(limit: number): Promise<string[]> {
    // Mock implementation - would typically query analytics
    return Array.from({ length: limit }, (_, i) => `user_${i + 1}`);
  }

  private async getRecentlyActiveUsers(): Promise<User[]> {
    // Mock implementation - would typically query activity logs
    return [];
  }
}

export { CacheWarming };
```

### Cache Metrics

```typescript
// src/caching/metrics.ts
interface CacheMetrics {
  hits: number;
  misses: number;
  sets: number;
  deletes: number;
  errors: number;
}

class CacheMetricsCollector {
  private metrics: CacheMetrics = {
    hits: 0,
    misses: 0,
    sets: 0,
    deletes: 0,
    errors: 0
  };

  recordHit(): void {
    this.metrics.hits++;
  }

  recordMiss(): void {
    this.metrics.misses++;
  }

  recordSet(): void {
    this.metrics.sets++;
  }

  recordDelete(): void {
    this.metrics.deletes++;
  }

  recordError(): void {
    this.metrics.errors++;
  }

  getMetrics(): CacheMetrics & { hitRatio: number } {
    const total = this.metrics.hits + this.metrics.misses;
    const hitRatio = total > 0 ? this.metrics.hits / total : 0;
    
    return {
      ...this.metrics,
      hitRatio
    };
  }

  reset(): void {
    this.metrics = {
      hits: 0,
      misses: 0,
      sets: 0,
      deletes: 0,
      errors: 0
    };
  }
}

class InstrumentedCache extends CacheAside {
  private metricsCollector: CacheMetricsCollector;

  constructor(
    client: RedisClientType,
    options: CacheOptions = {},
    metricsCollector: CacheMetricsCollector
  ) {
    super(client, options);
    this.metricsCollector = metricsCollector;
  }

  async get<T>(key: string): Promise<T | null> {
    try {
      const result = await super.get<T>(key);
      
      if (result !== null) {
        this.metricsCollector.recordHit();
      } else {
        this.metricsCollector.recordMiss();
      }
      
      return result;
    } catch (error) {
      this.metricsCollector.recordError();
      throw error;
    }
  }

  async set<T>(key: string, value: T, ttl?: number): Promise<boolean> {
    try {
      const result = await super.set(key, value, ttl);
      this.metricsCollector.recordSet();
      return result;
    } catch (error) {
      this.metricsCollector.recordError();
      throw error;
    }
  }

  async delete(key: string): Promise<boolean> {
    try {
      const result = await super.delete(key);
      this.metricsCollector.recordDelete();
      return result;
    } catch (error) {
      this.metricsCollector.recordError();
      throw error;
    }
  }
}

export { CacheMetricsCollector, InstrumentedCache, CacheMetrics };
```

## Real-world Examples

### E-commerce Product Cache

```typescript
// src/examples/product-cache.ts
interface Product {
  id: string;
  name: string;
  price: number;
  category: string;
  inventory: number;
  description: string;
  images: string[];
  createdAt: Date;
  updatedAt: Date;
}

class ProductCacheService {
  private cache: CacheAside;
  private invalidation: CacheInvalidation;

  constructor(cache: CacheAside, invalidation: CacheInvalidation) {
    this.cache = cache;
    this.invalidation = invalidation;
  }

  async getProduct(id: string): Promise<Product | null> {
    return await this.cache.getOrSet(
      `product:${id}`,
      () => this.fetchProductFromDB(id),
      3600 // 1 hour
    );
  }

  async getProductsByCategory(category: string): Promise<Product[]> {
    return await this.cache.getOrSet(
      `products:category:${category}`,
      () => this.fetchProductsByCategoryFromDB(category),
      1800 // 30 minutes
    );
  }

  async updateProduct(id: string, updates: Partial<Product>): Promise<Product | null> {
    const product = await this.updateProductInDB(id, updates);
    
    if (product) {
      // Update product cache
      await this.cache.set(`product:${id}`, product, 3600);
      
      // Invalidate category cache
      await this.invalidation.invalidateByTag(`category:${product.category}`);
      
      // If category changed, invalidate old category too
      if (updates.category && updates.category !== product.category) {
        await this.invalidation.invalidateByTag(`category:${updates.category}`);
      }
    }
    
    return product;
  }

  async updateInventory(id: string, quantity: number): Promise<boolean> {
    // Use Lua script for atomic inventory update
    const script = `
      local product_key = KEYS[1]
      local quantity = tonumber(ARGV[1])
      
      local product = redis.call('GET', product_key)
      if product == false then
        return nil
      end
      
      local product_data = cjson.decode(product)
      product_data.inventory = quantity
      product_data.updatedAt = ARGV[2]
      
      redis.call('SET', product_key, cjson.encode(product_data))
      return product_data.inventory
    `;
    
    try {
      const result = await this.cache['client'].eval(script, {
        keys: [`product:${id}`],
        arguments: [quantity.toString(), new Date().toISOString()]
      });
      
      return result !== null;
    } catch (error) {
      // Fallback to cache invalidation
      await this.cache.delete(`product:${id}`);
      return false;
    }
  }

  private async fetchProductFromDB(id: string): Promise<Product | null> {
    // Mock database fetch
    return null;
  }

  private async fetchProductsByCategoryFromDB(category: string): Promise<Product[]> {
    // Mock database fetch
    return [];
  }

  private async updateProductInDB(id: string, updates: Partial<Product>): Promise<Product | null> {
    // Mock database update
    return null;
  }
}

export { ProductCacheService, Product };
```

### Usage Examples

```typescript
// src/examples/caching-usage.ts
import { createClient } from 'redis';
import { CacheAside } from '../caching/cache-aside';
import { CachedUserService } from '../caching/user-service';
import { MultiLevelCache, MemoryCache } from '../caching/multi-level';
import { CacheMetricsCollector, InstrumentedCache } from '../caching/metrics';

async function runCachingExamples() {
  const client = createClient({
    host: 'localhost',
    port: 6379
  });
  
  await client.connect();
  
  try {
    // Basic cache-aside example
    console.log('=== Cache-Aside Example ===');
    const cache = new CacheAside(client, { prefix: 'app:', ttl: 3600 });
    
    // Set and get
    await cache.set('user:123', { name: 'John', email: 'john@example.com' });
    const user = await cache.get('user:123');
    console.log('Cached user:', user);
    
    // Get or set pattern
    const userData = await cache.getOrSet(
      'user:456',
      async () => {
        console.log('Fetching from database...');
        return { name: 'Jane', email: 'jane@example.com' };
      }
    );
    console.log('User data:', userData);
    
    // Multi-level caching example
    console.log('\n=== Multi-Level Cache Example ===');
    const l1Cache = new MemoryCache(300000); // 5 minutes
    const multiCache = new MultiLevelCache(l1Cache, cache, {} as any);
    
    // This would hit L1 cache on subsequent calls
    const cachedUser = await multiCache.getUserById('123');
    console.log('Multi-level cached user:', cachedUser);
    
    // Metrics example
    console.log('\n=== Cache Metrics Example ===');
    const metricsCollector = new CacheMetricsCollector();
    const instrumentedCache = new InstrumentedCache(client, {}, metricsCollector);
    
    // Perform some operations
    await instrumentedCache.set('test:1', 'value1');
    await instrumentedCache.get('test:1'); // Hit
    await instrumentedCache.get('test:2'); // Miss
    
    const metrics = metricsCollector.getMetrics();
    console.log('Cache metrics:', metrics);
    
  } finally {
    await client.disconnect();
  }
}

// Run examples
if (require.main === module) {
  runCachingExamples().catch(console.error);
}

export { runCachingExamples };
```

## Best Practices

### 1. Cache Key Design

- Use consistent naming conventions
- Include version information when needed
- Use hierarchical keys for easy invalidation
- Avoid very long keys

### 2. TTL Strategy

- Set appropriate TTLs based on data volatility
- Use shorter TTLs for frequently changing data
- Consider using different TTLs for different data types

### 3. Cache Size Management

- Monitor cache memory usage
- Implement cache eviction policies
- Use compression for large objects
- Consider data structure optimization

### 4. Error Handling

- Always handle cache failures gracefully
- Implement fallback to data source
- Log cache errors for monitoring
- Use circuit breakers for cache dependencies

## Next Steps

Now that you understand caching strategies, you're ready to explore:

- [Part 11: Session Management](./part-11-session-management.md)

## Quick Reference

### Cache Patterns
```typescript
// Cache-Aside
const cached = await cache.get(key);
if (!cached) {
  const data = await database.get(key);
  await cache.set(key, data);
  return data;
}
return cached;

// Write-Through
const data = await database.update(key, value);
await cache.set(key, data);

// Write-Behind
await cache.set(key, value);
queue.push({ key, value }); // Async DB write
```

### Cache Operations
```typescript
// Basic operations
await cache.set('key', value, ttl);
const value = await cache.get('key');
await cache.delete('key');

// Batch operations
await cache.mset([{key: 'k1', value: 'v1'}, {key: 'k2', value: 'v2'}]);
const values = await cache.mget(['k1', 'k2']);

// Conditional operations
const value = await cache.getOrSet('key', () => fetchFromDB());
```

---

**Previous**: [← Part 9: Lua Scripting](./part-09-lua-scripting.md)  
**Next**: [Part 11: Session Management →](./part-11-session-management.md)