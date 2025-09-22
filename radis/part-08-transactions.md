# Part 8: Transactions and Pipelining

Redis transactions and pipelining are powerful features that help optimize performance and ensure data consistency in your Node.js applications.

## Table of Contents

1. [Understanding Redis Transactions](#understanding-redis-transactions)
2. [MULTI/EXEC Commands](#multiexec-commands)
3. [WATCH Command](#watch-command)
4. [Pipelining Concepts](#pipelining-concepts)
5. [Node.js Implementation](#nodejs-implementation)
6. [Advanced Patterns](#advanced-patterns)
7. [Performance Comparison](#performance-comparison)
8. [Best Practices](#best-practices)
9. [Real-world Examples](#real-world-examples)

## Understanding Redis Transactions

Redis transactions allow you to execute a group of commands atomically. All commands in a transaction are executed sequentially without interruption.

### Key Characteristics

- **Atomic**: All commands execute or none do
- **Isolated**: No other commands can execute during transaction
- **No rollback**: Redis doesn't support rollback on failure
- **Queued execution**: Commands are queued and executed together

## MULTI/EXEC Commands

### Basic Transaction Structure

```typescript
// src/transactions/basic-transaction.ts
import { createClient, RedisClientType } from 'redis';

class BasicTransaction {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async transferFunds(fromAccount: string, toAccount: string, amount: number): Promise<boolean> {
    try {
      // Start transaction
      const multi = this.client.multi();
      
      // Queue commands
      multi.decrBy(`account:${fromAccount}`, amount);
      multi.incrBy(`account:${toAccount}`, amount);
      multi.set(`transfer:${Date.now()}`, JSON.stringify({
        from: fromAccount,
        to: toAccount,
        amount,
        timestamp: new Date().toISOString()
      }));
      
      // Execute transaction
      const results = await multi.exec();
      
      // Check if transaction was successful
      return results !== null;
    } catch (error) {
      console.error('Transaction failed:', error);
      return false;
    }
  }

  async batchUserUpdate(userId: string, updates: any): Promise<boolean> {
    const multi = this.client.multi();
    
    // Update user profile
    multi.hMSet(`user:${userId}`, updates.profile || {});
    
    // Update user preferences
    if (updates.preferences) {
      multi.hMSet(`user:${userId}:preferences`, updates.preferences);
    }
    
    // Update last activity
    multi.set(`user:${userId}:last_activity`, Date.now());
    
    // Add to activity log
    multi.lPush(`user:${userId}:activity`, JSON.stringify({
      action: 'profile_update',
      timestamp: Date.now()
    }));
    
    const results = await multi.exec();
    return results !== null;
  }
}

export { BasicTransaction };
```

## WATCH Command

The WATCH command provides optimistic locking for transactions.

```typescript
// src/transactions/watched-transaction.ts
class WatchedTransaction {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async conditionalIncrement(key: string, maxValue: number): Promise<boolean> {
    let attempts = 0;
    const maxAttempts = 10;

    while (attempts < maxAttempts) {
      try {
        // Watch the key
        await this.client.watch(key);
        
        // Get current value
        const currentValue = parseInt(await this.client.get(key) || '0');
        
        // Check condition
        if (currentValue >= maxValue) {
          await this.client.unwatch();
          return false;
        }
        
        // Start transaction
        const multi = this.client.multi();
        multi.incr(key);
        
        // Execute transaction
        const results = await multi.exec();
        
        // If results is null, the key was modified
        if (results === null) {
          attempts++;
          continue;
        }
        
        return true;
      } catch (error) {
        console.error('Watched transaction failed:', error);
        attempts++;
      }
    }
    
    return false;
  }

  async safeTransfer(fromKey: string, toKey: string, amount: number): Promise<{
    success: boolean;
    reason?: string;
  }> {
    let attempts = 0;
    const maxAttempts = 5;

    while (attempts < maxAttempts) {
      try {
        // Watch both accounts
        await this.client.watch([fromKey, toKey]);
        
        // Get current balances
        const fromBalance = parseInt(await this.client.get(fromKey) || '0');
        const toBalance = parseInt(await this.client.get(toKey) || '0');
        
        // Check if sufficient funds
        if (fromBalance < amount) {
          await this.client.unwatch();
          return { success: false, reason: 'Insufficient funds' };
        }
        
        // Start transaction
        const multi = this.client.multi();
        multi.decrBy(fromKey, amount);
        multi.incrBy(toKey, amount);
        
        // Execute transaction
        const results = await multi.exec();
        
        // Check if transaction succeeded
        if (results === null) {
          attempts++;
          continue;
        }
        
        return { success: true };
      } catch (error) {
        console.error('Safe transfer failed:', error);
        attempts++;
      }
    }
    
    return { success: false, reason: 'Max attempts exceeded' };
  }
}

export { WatchedTransaction };
```

## Pipelining Concepts

Pipelining allows you to send multiple commands without waiting for replies, improving performance.

```typescript
// src/transactions/pipeline.ts
class RedisPipeline {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async batchGet(keys: string[]): Promise<(string | null)[]> {
    // Create pipeline
    const pipeline = this.client.multi();
    
    // Add all GET commands
    keys.forEach(key => pipeline.get(key));
    
    // Execute pipeline
    const results = await pipeline.exec();
    
    // Extract values from results
    return results ? results.map(result => result as string | null) : [];
  }

  async batchSet(keyValuePairs: { key: string; value: string; ttl?: number }[]): Promise<boolean> {
    const pipeline = this.client.multi();
    
    keyValuePairs.forEach(({ key, value, ttl }) => {
      if (ttl) {
        pipeline.setEx(key, ttl, value);
      } else {
        pipeline.set(key, value);
      }
    });
    
    const results = await pipeline.exec();
    return results !== null;
  }

  async batchHashOperations(operations: Array<{
    type: 'hset' | 'hget' | 'hdel' | 'hincrby';
    key: string;
    field?: string;
    value?: string | number;
    fields?: string[];
  }>): Promise<any[]> {
    const pipeline = this.client.multi();
    
    operations.forEach(op => {
      switch (op.type) {
        case 'hset':
          if (op.field && op.value !== undefined) {
            pipeline.hSet(op.key, op.field, op.value.toString());
          }
          break;
        case 'hget':
          if (op.field) {
            pipeline.hGet(op.key, op.field);
          }
          break;
        case 'hdel':
          if (op.fields) {
            pipeline.hDel(op.key, op.fields);
          }
          break;
        case 'hincrby':
          if (op.field && typeof op.value === 'number') {
            pipeline.hIncrBy(op.key, op.field, op.value);
          }
          break;
      }
    });
    
    const results = await pipeline.exec();
    return results || [];
  }
}

export { RedisPipeline };
```

## Advanced Patterns

### Transaction with Conditional Logic

```typescript
// src/transactions/conditional-transaction.ts
class ConditionalTransaction {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async processOrder(orderId: string, userId: string, items: Array<{
    productId: string;
    quantity: number;
    price: number;
  }>): Promise<{
    success: boolean;
    orderId?: string;
    error?: string;
  }> {
    let attempts = 0;
    const maxAttempts = 3;

    while (attempts < maxAttempts) {
      try {
        // Watch inventory keys
        const inventoryKeys = items.map(item => `inventory:${item.productId}`);
        await this.client.watch(inventoryKeys);
        
        // Check inventory availability
        const inventoryChecks = await Promise.all(
          items.map(async item => {
            const available = parseInt(
              await this.client.get(`inventory:${item.productId}`) || '0'
            );
            return {
              productId: item.productId,
              requested: item.quantity,
              available,
              sufficient: available >= item.quantity
            };
          })
        );
        
        // Check if all items are available
        const insufficientItems = inventoryChecks.filter(check => !check.sufficient);
        if (insufficientItems.length > 0) {
          await this.client.unwatch();
          return {
            success: false,
            error: `Insufficient inventory for: ${insufficientItems.map(i => i.productId).join(', ')}`
          };
        }
        
        // Calculate total
        const total = items.reduce((sum, item) => sum + (item.price * item.quantity), 0);
        
        // Start transaction
        const multi = this.client.multi();
        
        // Update inventory
        items.forEach(item => {
          multi.decrBy(`inventory:${item.productId}`, item.quantity);
        });
        
        // Create order
        const orderData = {
          orderId,
          userId,
          items,
          total,
          status: 'confirmed',
          createdAt: new Date().toISOString()
        };
        
        multi.hMSet(`order:${orderId}`, orderData as any);
        multi.lPush(`user:${userId}:orders`, orderId);
        multi.incrBy(`user:${userId}:total_spent`, total);
        
        // Execute transaction
        const results = await multi.exec();
        
        if (results === null) {
          attempts++;
          continue;
        }
        
        return { success: true, orderId };
      } catch (error) {
        console.error('Order processing failed:', error);
        attempts++;
      }
    }
    
    return { success: false, error: 'Max attempts exceeded' };
  }
}

export { ConditionalTransaction };
```

### Distributed Lock with Transactions

```typescript
// src/transactions/distributed-lock.ts
class DistributedLock {
  private client: RedisClientType;
  private lockTimeout: number;

  constructor(client: RedisClientType, lockTimeout: number = 10000) {
    this.client = client;
    this.lockTimeout = lockTimeout;
  }

  async acquireLock(resource: string, identifier: string): Promise<boolean> {
    const lockKey = `lock:${resource}`;
    const expireTime = Date.now() + this.lockTimeout;
    
    const result = await this.client.set(
      lockKey,
      identifier,
      {
        PX: this.lockTimeout,
        NX: true
      }
    );
    
    return result === 'OK';
  }

  async releaseLock(resource: string, identifier: string): Promise<boolean> {
    const lockKey = `lock:${resource}`;
    
    // Lua script to ensure we only release our own lock
    const luaScript = `
      if redis.call('get', KEYS[1]) == ARGV[1] then
        return redis.call('del', KEYS[1])
      else
        return 0
      end
    `;
    
    const result = await this.client.eval(luaScript, {
      keys: [lockKey],
      arguments: [identifier]
    });
    
    return result === 1;
  }

  async withLock<T>(
    resource: string,
    identifier: string,
    operation: () => Promise<T>
  ): Promise<T | null> {
    const acquired = await this.acquireLock(resource, identifier);
    
    if (!acquired) {
      return null;
    }
    
    try {
      return await operation();
    } finally {
      await this.releaseLock(resource, identifier);
    }
  }
}

export { DistributedLock };
```

## Performance Comparison

```typescript
// src/transactions/performance-test.ts
class PerformanceTest {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async testSequentialOperations(count: number): Promise<number> {
    const start = Date.now();
    
    for (let i = 0; i < count; i++) {
      await this.client.set(`test:seq:${i}`, `value${i}`);
      await this.client.get(`test:seq:${i}`);
    }
    
    return Date.now() - start;
  }

  async testPipelinedOperations(count: number): Promise<number> {
    const start = Date.now();
    
    const pipeline = this.client.multi();
    
    for (let i = 0; i < count; i++) {
      pipeline.set(`test:pipe:${i}`, `value${i}`);
      pipeline.get(`test:pipe:${i}`);
    }
    
    await pipeline.exec();
    
    return Date.now() - start;
  }

  async testTransactionOperations(count: number): Promise<number> {
    const start = Date.now();
    
    const multi = this.client.multi();
    
    for (let i = 0; i < count; i++) {
      multi.set(`test:trans:${i}`, `value${i}`);
      multi.incr(`counter:trans`);
    }
    
    await multi.exec();
    
    return Date.now() - start;
  }

  async runPerformanceComparison(operationCount: number = 1000): Promise<void> {
    console.log(`\nPerformance Comparison (${operationCount} operations):\n`);
    
    const sequentialTime = await this.testSequentialOperations(operationCount);
    console.log(`Sequential: ${sequentialTime}ms`);
    
    const pipelineTime = await this.testPipelinedOperations(operationCount);
    console.log(`Pipelined: ${pipelineTime}ms`);
    
    const transactionTime = await this.testTransactionOperations(operationCount);
    console.log(`Transaction: ${transactionTime}ms`);
    
    console.log(`\nSpeedup:`);
    console.log(`Pipeline vs Sequential: ${(sequentialTime / pipelineTime).toFixed(2)}x`);
    console.log(`Transaction vs Sequential: ${(sequentialTime / transactionTime).toFixed(2)}x`);
  }
}

export { PerformanceTest };
```

## Real-world Examples

### Shopping Cart with Transactions

```typescript
// src/examples/shopping-cart.ts
interface CartItem {
  productId: string;
  quantity: number;
  price: number;
  name: string;
}

class ShoppingCart {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async addToCart(userId: string, item: CartItem): Promise<boolean> {
    const cartKey = `cart:${userId}`;
    const inventoryKey = `inventory:${item.productId}`;
    
    let attempts = 0;
    const maxAttempts = 3;
    
    while (attempts < maxAttempts) {
      try {
        await this.client.watch(inventoryKey);
        
        const available = parseInt(await this.client.get(inventoryKey) || '0');
        
        if (available < item.quantity) {
          await this.client.unwatch();
          return false;
        }
        
        const multi = this.client.multi();
        
        // Add item to cart
        multi.hSet(cartKey, item.productId, JSON.stringify(item));
        
        // Reserve inventory
        multi.decrBy(inventoryKey, item.quantity);
        multi.incrBy(`reserved:${item.productId}`, item.quantity);
        
        // Set cart expiration
        multi.expire(cartKey, 3600); // 1 hour
        
        const results = await multi.exec();
        
        if (results === null) {
          attempts++;
          continue;
        }
        
        return true;
      } catch (error) {
        console.error('Add to cart failed:', error);
        attempts++;
      }
    }
    
    return false;
  }

  async checkout(userId: string): Promise<{
    success: boolean;
    orderId?: string;
    total?: number;
  }> {
    const cartKey = `cart:${userId}`;
    const orderId = `order_${Date.now()}_${userId}`;
    
    try {
      // Get all cart items
      const cartData = await this.client.hGetAll(cartKey);
      
      if (Object.keys(cartData).length === 0) {
        return { success: false };
      }
      
      const items: CartItem[] = Object.values(cartData).map(item => JSON.parse(item));
      const total = items.reduce((sum, item) => sum + (item.price * item.quantity), 0);
      
      const multi = this.client.multi();
      
      // Create order
      multi.hMSet(`order:${orderId}`, {
        userId,
        items: JSON.stringify(items),
        total: total.toString(),
        status: 'confirmed',
        createdAt: new Date().toISOString()
      });
      
      // Clear cart
      multi.del(cartKey);
      
      // Update user order history
      multi.lPush(`user:${userId}:orders`, orderId);
      
      // Release reserved inventory (it's already decremented)
      items.forEach(item => {
        multi.decrBy(`reserved:${item.productId}`, item.quantity);
      });
      
      const results = await multi.exec();
      
      if (results === null) {
        return { success: false };
      }
      
      return { success: true, orderId, total };
    } catch (error) {
      console.error('Checkout failed:', error);
      return { success: false };
    }
  }
}

export { ShoppingCart, CartItem };
```

### Usage Examples

```typescript
// src/examples/usage.ts
import { createClient } from 'redis';
import { BasicTransaction } from '../transactions/basic-transaction';
import { WatchedTransaction } from '../transactions/watched-transaction';
import { RedisPipeline } from '../transactions/pipeline';
import { PerformanceTest } from '../transactions/performance-test';
import { ShoppingCart } from './shopping-cart';

async function runTransactionExamples() {
  const client = createClient({
    host: 'localhost',
    port: 6379
  });
  
  await client.connect();
  
  try {
    // Basic Transaction Example
    console.log('=== Basic Transaction Example ===');
    const basicTx = new BasicTransaction(client);
    
    // Set up initial balances
    await client.set('account:alice', '1000');
    await client.set('account:bob', '500');
    
    const transferResult = await basicTx.transferFunds('alice', 'bob', 200);
    console.log('Transfer result:', transferResult);
    
    const aliceBalance = await client.get('account:alice');
    const bobBalance = await client.get('account:bob');
    console.log(`Alice: ${aliceBalance}, Bob: ${bobBalance}`);
    
    // Watched Transaction Example
    console.log('\n=== Watched Transaction Example ===');
    const watchedTx = new WatchedTransaction(client);
    
    await client.set('counter', '5');
    const incrementResult = await watchedTx.conditionalIncrement('counter', 10);
    console.log('Increment result:', incrementResult);
    
    // Pipeline Example
    console.log('\n=== Pipeline Example ===');
    const pipeline = new RedisPipeline(client);
    
    const batchSetResult = await pipeline.batchSet([
      { key: 'user:1', value: 'Alice', ttl: 3600 },
      { key: 'user:2', value: 'Bob', ttl: 3600 },
      { key: 'user:3', value: 'Charlie' }
    ]);
    console.log('Batch set result:', batchSetResult);
    
    const batchGetResult = await pipeline.batchGet(['user:1', 'user:2', 'user:3']);
    console.log('Batch get result:', batchGetResult);
    
    // Performance Test
    console.log('\n=== Performance Test ===');
    const perfTest = new PerformanceTest(client);
    await perfTest.runPerformanceComparison(100);
    
    // Shopping Cart Example
    console.log('\n=== Shopping Cart Example ===');
    const cart = new ShoppingCart(client);
    
    // Set up inventory
    await client.set('inventory:product1', '10');
    await client.set('inventory:product2', '5');
    
    const addResult = await cart.addToCart('user123', {
      productId: 'product1',
      quantity: 2,
      price: 29.99,
      name: 'Awesome Product'
    });
    console.log('Add to cart result:', addResult);
    
    const checkoutResult = await cart.checkout('user123');
    console.log('Checkout result:', checkoutResult);
    
  } finally {
    await client.disconnect();
  }
}

// Run examples
if (require.main === module) {
  runTransactionExamples().catch(console.error);
}

export { runTransactionExamples };
```

## Best Practices

### 1. Transaction Guidelines

- Keep transactions short and simple
- Avoid long-running operations in transactions
- Use WATCH for optimistic locking
- Handle transaction failures gracefully
- Implement retry logic with exponential backoff

### 2. Pipelining Guidelines

- Use pipelining for bulk operations
- Batch related commands together
- Monitor memory usage with large pipelines
- Consider network latency benefits

### 3. Error Handling

```typescript
class TransactionErrorHandler {
  static async executeWithRetry<T>(
    operation: () => Promise<T>,
    maxAttempts: number = 3,
    baseDelay: number = 100
  ): Promise<T | null> {
    let attempts = 0;
    
    while (attempts < maxAttempts) {
      try {
        return await operation();
      } catch (error) {
        attempts++;
        
        if (attempts >= maxAttempts) {
          console.error('Max attempts exceeded:', error);
          return null;
        }
        
        // Exponential backoff
        const delay = baseDelay * Math.pow(2, attempts - 1);
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }
    
    return null;
  }
}
```

## Next Steps

Now that you understand transactions and pipelining, you're ready to explore:

- [Part 9: Lua Scripting](./part-09-lua-scripting.md)

## Quick Reference

### Transaction Commands
```typescript
// Basic transaction
const multi = client.multi();
multi.set('key1', 'value1');
multi.incr('counter');
const results = await multi.exec();

// With WATCH
await client.watch('key');
const multi = client.multi();
multi.set('key', 'new_value');
const results = await multi.exec(); // null if key was modified

// Discard transaction
multi.discard();
```

### Pipeline Commands
```typescript
// Basic pipeline
const pipeline = client.multi();
pipeline.get('key1');
pipeline.get('key2');
const results = await pipeline.exec();

// Batch operations
const keys = ['key1', 'key2', 'key3'];
const pipeline = client.multi();
keys.forEach(key => pipeline.get(key));
const values = await pipeline.exec();
```

---

**Previous**: [← Part 7: Pub/Sub Messaging](./part-07-pubsub.md)  
**Next**: [Part 9: Lua Scripting →](./part-09-lua-scripting.md)