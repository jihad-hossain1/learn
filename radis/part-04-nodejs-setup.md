# Part 4: Node.js Redis Client Setup

This part covers setting up Redis clients in Node.js, connection management, and basic integration patterns.

## Redis Client Libraries for Node.js

### Popular Redis Clients

1. **node-redis** (Official, Recommended)
   - Most popular and actively maintained
   - Full Redis command support
   - TypeScript support
   - Promise-based API

2. **ioredis**
   - High performance
   - Built-in cluster support
   - Lua scripting support
   - Pipeline support

3. **redis** (Legacy)
   - Older callback-based API
   - Being phased out

## Installing Redis Client

### Project Setup

```bash
# Create new Node.js project
mkdir redis-nodejs-app
cd redis-nodejs-app
npm init -y

# Install Redis client
npm install redis

# Install development dependencies
npm install --save-dev @types/node typescript ts-node nodemon

# Install additional utilities
npm install dotenv
```

### Package.json Configuration

```json
{
  "name": "redis-nodejs-app",
  "version": "1.0.0",
  "description": "Redis with Node.js learning project",
  "main": "index.js",
  "scripts": {
    "start": "node dist/index.js",
    "dev": "nodemon --exec ts-node src/index.ts",
    "build": "tsc",
    "test": "node --test"
  },
  "dependencies": {
    "redis": "^4.6.0",
    "dotenv": "^16.3.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "typescript": "^5.0.0",
    "ts-node": "^10.9.0",
    "nodemon": "^3.0.0"
  }
}
```

### TypeScript Configuration

```json
// tsconfig.json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist"]
}
```

## Basic Redis Connection

### Simple Connection (JavaScript)

```javascript
// src/basic-connection.js
const redis = require('redis');

async function basicConnection() {
  // Create Redis client
  const client = redis.createClient({
    host: 'localhost',
    port: 6379,
    // password: 'your_password', // if authentication is required
  });

  // Handle connection events
  client.on('connect', () => {
    console.log('Connected to Redis');
  });

  client.on('error', (err) => {
    console.error('Redis connection error:', err);
  });

  client.on('ready', () => {
    console.log('Redis client ready');
  });

  client.on('end', () => {
    console.log('Redis connection closed');
  });

  try {
    // Connect to Redis
    await client.connect();

    // Test connection
    const pong = await client.ping();
    console.log('Ping response:', pong);

    // Basic operations
    await client.set('test:key', 'Hello Redis!');
    const value = await client.get('test:key');
    console.log('Retrieved value:', value);

  } catch (error) {
    console.error('Error:', error);
  } finally {
    // Close connection
    await client.quit();
  }
}

basicConnection();
```

### TypeScript Connection

```typescript
// src/redis-client.ts
import { createClient, RedisClientType } from 'redis';
import dotenv from 'dotenv';

dotenv.config();

interface RedisConfig {
  host: string;
  port: number;
  password?: string;
  database?: number;
}

class RedisManager {
  private client: RedisClientType;
  private isConnected: boolean = false;

  constructor(config: RedisConfig) {
    this.client = createClient({
      socket: {
        host: config.host,
        port: config.port,
      },
      password: config.password,
      database: config.database || 0,
    });

    this.setupEventHandlers();
  }

  private setupEventHandlers(): void {
    this.client.on('connect', () => {
      console.log('🔗 Connecting to Redis...');
    });

    this.client.on('ready', () => {
      console.log('✅ Redis client ready');
      this.isConnected = true;
    });

    this.client.on('error', (err) => {
      console.error('❌ Redis connection error:', err);
      this.isConnected = false;
    });

    this.client.on('end', () => {
      console.log('🔌 Redis connection closed');
      this.isConnected = false;
    });

    this.client.on('reconnecting', () => {
      console.log('🔄 Reconnecting to Redis...');
    });
  }

  async connect(): Promise<void> {
    try {
      await this.client.connect();
    } catch (error) {
      console.error('Failed to connect to Redis:', error);
      throw error;
    }
  }

  async disconnect(): Promise<void> {
    try {
      await this.client.quit();
    } catch (error) {
      console.error('Error disconnecting from Redis:', error);
      throw error;
    }
  }

  getClient(): RedisClientType {
    if (!this.isConnected) {
      throw new Error('Redis client is not connected');
    }
    return this.client;
  }

  isClientConnected(): boolean {
    return this.isConnected;
  }

  async ping(): Promise<string> {
    return await this.client.ping();
  }
}

export default RedisManager;
```

### Environment Configuration

```bash
# .env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DATABASE=0

# Development
NODE_ENV=development

# Production settings
# REDIS_HOST=your-redis-server.com
# REDIS_PASSWORD=your-secure-password
# REDIS_TLS=true
```

### Configuration Manager

```typescript
// src/config.ts
import dotenv from 'dotenv';

dotenv.config();

interface Config {
  redis: {
    host: string;
    port: number;
    password?: string;
    database: number;
    tls?: boolean;
  };
  app: {
    port: number;
    env: string;
  };
}

const config: Config = {
  redis: {
    host: process.env.REDIS_HOST || 'localhost',
    port: parseInt(process.env.REDIS_PORT || '6379'),
    password: process.env.REDIS_PASSWORD || undefined,
    database: parseInt(process.env.REDIS_DATABASE || '0'),
    tls: process.env.REDIS_TLS === 'true',
  },
  app: {
    port: parseInt(process.env.PORT || '3000'),
    env: process.env.NODE_ENV || 'development',
  },
};

export default config;
```

## Connection Patterns

### Singleton Pattern

```typescript
// src/redis-singleton.ts
import { createClient, RedisClientType } from 'redis';
import config from './config';

class RedisSingleton {
  private static instance: RedisSingleton;
  private client: RedisClientType;
  private connected: boolean = false;

  private constructor() {
    this.client = createClient({
      socket: {
        host: config.redis.host,
        port: config.redis.port,
        tls: config.redis.tls,
      },
      password: config.redis.password,
      database: config.redis.database,
    });

    this.setupEventHandlers();
  }

  public static getInstance(): RedisSingleton {
    if (!RedisSingleton.instance) {
      RedisSingleton.instance = new RedisSingleton();
    }
    return RedisSingleton.instance;
  }

  private setupEventHandlers(): void {
    this.client.on('ready', () => {
      this.connected = true;
      console.log('Redis singleton connected');
    });

    this.client.on('error', (err) => {
      this.connected = false;
      console.error('Redis singleton error:', err);
    });

    this.client.on('end', () => {
      this.connected = false;
      console.log('Redis singleton disconnected');
    });
  }

  public async connect(): Promise<void> {
    if (!this.connected) {
      await this.client.connect();
    }
  }

  public async disconnect(): Promise<void> {
    if (this.connected) {
      await this.client.quit();
    }
  }

  public getClient(): RedisClientType {
    if (!this.connected) {
      throw new Error('Redis client not connected');
    }
    return this.client;
  }

  public isConnected(): boolean {
    return this.connected;
  }
}

export default RedisSingleton;
```

### Factory Pattern

```typescript
// src/redis-factory.ts
import { createClient, RedisClientType } from 'redis';

interface RedisConnectionOptions {
  host: string;
  port: number;
  password?: string;
  database?: number;
  retryDelayOnFailover?: number;
  maxRetriesPerRequest?: number;
}

class RedisFactory {
  static createClient(options: RedisConnectionOptions): RedisClientType {
    return createClient({
      socket: {
        host: options.host,
        port: options.port,
        reconnectStrategy: (retries) => {
          if (retries > 10) {
            return new Error('Too many retries');
          }
          return Math.min(retries * 50, 500);
        },
      },
      password: options.password,
      database: options.database || 0,
    });
  }

  static async createConnectedClient(options: RedisConnectionOptions): Promise<RedisClientType> {
    const client = this.createClient(options);
    await client.connect();
    return client;
  }

  static createClusterClient(nodes: string[]): RedisClientType {
    // For Redis Cluster (covered in advanced topics)
    return createClient({
      socket: {
        host: 'localhost',
        port: 6379,
      },
    });
  }
}

export default RedisFactory;
```

## Connection Pool Management

### Basic Pool Implementation

```typescript
// src/redis-pool.ts
import { createClient, RedisClientType } from 'redis';

interface PoolOptions {
  min: number;
  max: number;
  host: string;
  port: number;
  password?: string;
  database?: number;
}

class RedisPool {
  private pool: RedisClientType[] = [];
  private available: RedisClientType[] = [];
  private options: PoolOptions;

  constructor(options: PoolOptions) {
    this.options = options;
  }

  async initialize(): Promise<void> {
    // Create minimum number of connections
    for (let i = 0; i < this.options.min; i++) {
      const client = await this.createClient();
      this.pool.push(client);
      this.available.push(client);
    }
  }

  private async createClient(): Promise<RedisClientType> {
    const client = createClient({
      socket: {
        host: this.options.host,
        port: this.options.port,
      },
      password: this.options.password,
      database: this.options.database,
    });

    await client.connect();
    return client;
  }

  async acquire(): Promise<RedisClientType> {
    if (this.available.length > 0) {
      return this.available.pop()!;
    }

    if (this.pool.length < this.options.max) {
      const client = await this.createClient();
      this.pool.push(client);
      return client;
    }

    // Wait for available connection
    return new Promise((resolve) => {
      const checkAvailable = () => {
        if (this.available.length > 0) {
          resolve(this.available.pop()!);
        } else {
          setTimeout(checkAvailable, 10);
        }
      };
      checkAvailable();
    });
  }

  release(client: RedisClientType): void {
    if (this.pool.includes(client)) {
      this.available.push(client);
    }
  }

  async destroy(): Promise<void> {
    await Promise.all(this.pool.map(client => client.quit()));
    this.pool = [];
    this.available = [];
  }

  getStats() {
    return {
      total: this.pool.length,
      available: this.available.length,
      inUse: this.pool.length - this.available.length,
    };
  }
}

export default RedisPool;
```

## Error Handling and Retry Logic

### Robust Error Handling

```typescript
// src/redis-error-handler.ts
import { RedisClientType } from 'redis';

class RedisErrorHandler {
  private client: RedisClientType;
  private maxRetries: number;
  private retryDelay: number;

  constructor(client: RedisClientType, maxRetries = 3, retryDelay = 1000) {
    this.client = client;
    this.maxRetries = maxRetries;
    this.retryDelay = retryDelay;
  }

  async executeWithRetry<T>(
    operation: () => Promise<T>,
    operationName: string
  ): Promise<T> {
    let lastError: Error;

    for (let attempt = 1; attempt <= this.maxRetries; attempt++) {
      try {
        return await operation();
      } catch (error) {
        lastError = error as Error;
        console.warn(
          `${operationName} failed (attempt ${attempt}/${this.maxRetries}):`,
          error
        );

        if (attempt < this.maxRetries) {
          await this.delay(this.retryDelay * attempt);
        }
      }
    }

    throw new Error(
      `${operationName} failed after ${this.maxRetries} attempts: ${lastError.message}`
    );
  }

  private delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async safeGet(key: string): Promise<string | null> {
    return this.executeWithRetry(
      () => this.client.get(key),
      `GET ${key}`
    );
  }

  async safeSet(key: string, value: string, ttl?: number): Promise<string | null> {
    return this.executeWithRetry(
      () => ttl ? this.client.setEx(key, ttl, value) : this.client.set(key, value),
      `SET ${key}`
    );
  }

  async safeDel(key: string): Promise<number> {
    return this.executeWithRetry(
      () => this.client.del(key),
      `DEL ${key}`
    );
  }
}

export default RedisErrorHandler;
```

### Circuit Breaker Pattern

```typescript
// src/redis-circuit-breaker.ts
import { RedisClientType } from 'redis';

enum CircuitState {
  CLOSED = 'CLOSED',
  OPEN = 'OPEN',
  HALF_OPEN = 'HALF_OPEN'
}

class RedisCircuitBreaker {
  private client: RedisClientType;
  private state: CircuitState = CircuitState.CLOSED;
  private failureCount: number = 0;
  private lastFailureTime: number = 0;
  private successCount: number = 0;

  private readonly failureThreshold: number;
  private readonly recoveryTimeout: number;
  private readonly monitoringPeriod: number;

  constructor(
    client: RedisClientType,
    failureThreshold = 5,
    recoveryTimeout = 60000,
    monitoringPeriod = 10000
  ) {
    this.client = client;
    this.failureThreshold = failureThreshold;
    this.recoveryTimeout = recoveryTimeout;
    this.monitoringPeriod = monitoringPeriod;
  }

  async execute<T>(operation: () => Promise<T>): Promise<T> {
    if (this.state === CircuitState.OPEN) {
      if (Date.now() - this.lastFailureTime > this.recoveryTimeout) {
        this.state = CircuitState.HALF_OPEN;
        this.successCount = 0;
      } else {
        throw new Error('Circuit breaker is OPEN');
      }
    }

    try {
      const result = await operation();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }

  private onSuccess(): void {
    this.failureCount = 0;
    
    if (this.state === CircuitState.HALF_OPEN) {
      this.successCount++;
      if (this.successCount >= 3) {
        this.state = CircuitState.CLOSED;
      }
    }
  }

  private onFailure(): void {
    this.failureCount++;
    this.lastFailureTime = Date.now();

    if (this.failureCount >= this.failureThreshold) {
      this.state = CircuitState.OPEN;
    }
  }

  getState(): CircuitState {
    return this.state;
  }

  getStats() {
    return {
      state: this.state,
      failureCount: this.failureCount,
      successCount: this.successCount,
      lastFailureTime: this.lastFailureTime,
    };
  }
}

export default RedisCircuitBreaker;
```

## Health Checks and Monitoring

### Health Check Implementation

```typescript
// src/redis-health-check.ts
import { RedisClientType } from 'redis';

interface HealthStatus {
  status: 'healthy' | 'unhealthy' | 'degraded';
  latency: number;
  memory: {
    used: string;
    peak: string;
    fragmentation: number;
  };
  connections: number;
  uptime: number;
  version: string;
}

class RedisHealthCheck {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  async checkHealth(): Promise<HealthStatus> {
    try {
      const startTime = Date.now();
      
      // Test basic connectivity
      await this.client.ping();
      
      const latency = Date.now() - startTime;
      
      // Get server info
      const info = await this.client.info();
      const parsedInfo = this.parseInfo(info);
      
      return {
        status: this.determineStatus(latency, parsedInfo),
        latency,
        memory: {
          used: parsedInfo.used_memory_human || '0B',
          peak: parsedInfo.used_memory_peak_human || '0B',
          fragmentation: parseFloat(parsedInfo.mem_fragmentation_ratio || '1'),
        },
        connections: parseInt(parsedInfo.connected_clients || '0'),
        uptime: parseInt(parsedInfo.uptime_in_seconds || '0'),
        version: parsedInfo.redis_version || 'unknown',
      };
    } catch (error) {
      return {
        status: 'unhealthy',
        latency: -1,
        memory: { used: '0B', peak: '0B', fragmentation: 0 },
        connections: 0,
        uptime: 0,
        version: 'unknown',
      };
    }
  }

  private parseInfo(info: string): Record<string, string> {
    const result: Record<string, string> = {};
    
    info.split('\r\n').forEach(line => {
      if (line.includes(':')) {
        const [key, value] = line.split(':');
        result[key] = value;
      }
    });
    
    return result;
  }

  private determineStatus(
    latency: number,
    info: Record<string, string>
  ): 'healthy' | 'unhealthy' | 'degraded' {
    if (latency > 1000) return 'unhealthy';
    if (latency > 100) return 'degraded';
    
    const fragmentation = parseFloat(info.mem_fragmentation_ratio || '1');
    if (fragmentation > 2) return 'degraded';
    
    return 'healthy';
  }
}

export default RedisHealthCheck;
```

## Usage Examples

### Basic Usage Example

```typescript
// src/examples/basic-usage.ts
import RedisManager from '../redis-client';
import config from '../config';

async function basicUsageExample() {
  const redisManager = new RedisManager(config.redis);
  
  try {
    await redisManager.connect();
    
    const client = redisManager.getClient();
    
    // String operations
    await client.set('user:1000:name', 'John Doe');
    const name = await client.get('user:1000:name');
    console.log('User name:', name);
    
    // Hash operations
    await client.hSet('user:1000:profile', {
      name: 'John Doe',
      email: 'john@example.com',
      age: '30'
    });
    
    const profile = await client.hGetAll('user:1000:profile');
    console.log('User profile:', profile);
    
    // List operations
    await client.lPush('notifications:user:1000', 'Welcome!');
    await client.lPush('notifications:user:1000', 'New message');
    
    const notifications = await client.lRange('notifications:user:1000', 0, -1);
    console.log('Notifications:', notifications);
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await redisManager.disconnect();
  }
}

basicUsageExample();
```

### Express.js Integration

```typescript
// src/examples/express-integration.ts
import express from 'express';
import RedisSingleton from '../redis-singleton';
import RedisHealthCheck from '../redis-health-check';

const app = express();
app.use(express.json());

// Initialize Redis connection
const redis = RedisSingleton.getInstance();
redis.connect();

const healthCheck = new RedisHealthCheck(redis.getClient());

// Health check endpoint
app.get('/health', async (req, res) => {
  try {
    const health = await healthCheck.checkHealth();
    res.status(health.status === 'healthy' ? 200 : 503).json(health);
  } catch (error) {
    res.status(503).json({ status: 'unhealthy', error: error.message });
  }
});

// User endpoints
app.get('/users/:id', async (req, res) => {
  try {
    const client = redis.getClient();
    const user = await client.hGetAll(`user:${req.params.id}`);
    
    if (Object.keys(user).length === 0) {
      return res.status(404).json({ error: 'User not found' });
    }
    
    res.json(user);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

app.post('/users/:id', async (req, res) => {
  try {
    const client = redis.getClient();
    await client.hSet(`user:${req.params.id}`, req.body);
    res.status(201).json({ message: 'User created' });
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Graceful shutdown
process.on('SIGTERM', async () => {
  console.log('Shutting down gracefully...');
  await redis.disconnect();
  process.exit(0);
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
```

## Next Steps

Now that you have Redis client setup and basic integration, you're ready to explore:

- [Part 5: CRUD Operations in Node.js](./part-05-crud-operations.md)

## Quick Reference

### Installation
```bash
npm install redis
npm install --save-dev @types/node typescript
```

### Basic Connection
```typescript
import { createClient } from 'redis';

const client = createClient({
  socket: { host: 'localhost', port: 6379 },
  password: 'your_password'
});

await client.connect();
```

### Error Handling
```typescript
client.on('error', (err) => console.error('Redis error:', err));
client.on('connect', () => console.log('Connected to Redis'));
client.on('ready', () => console.log('Redis client ready'));
```

### Environment Variables
```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=your_password
REDIS_DATABASE=0
```

---

**Previous**: [← Part 3: Advanced Data Structures](./part-03-advanced-structures.md)  
**Next**: [Part 5: CRUD Operations in Node.js →](./part-05-crud-operations.md)