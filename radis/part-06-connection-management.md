# Part 6: Connection Management and Error Handling

This part covers advanced connection management, error handling strategies, and best practices for production Redis applications.

## Connection Management Strategies

### Connection Lifecycle Management

```typescript
// src/connection/connection-manager.ts
import { createClient, RedisClientType } from 'redis';
import { EventEmitter } from 'events';

interface ConnectionConfig {
  host: string;
  port: number;
  password?: string;
  database?: number;
  connectTimeout?: number;
  lazyConnect?: boolean;
  retryDelayOnFailover?: number;
  maxRetriesPerRequest?: number;
  keepAlive?: number;
}

interface ConnectionStats {
  totalConnections: number;
  activeConnections: number;
  failedConnections: number;
  reconnections: number;
  lastConnectedAt?: Date;
  lastDisconnectedAt?: Date;
  uptime: number;
}

class ConnectionManager extends EventEmitter {
  private client: RedisClientType | null = null;
  private config: ConnectionConfig;
  private stats: ConnectionStats;
  private isConnecting: boolean = false;
  private isConnected: boolean = false;
  private reconnectTimer: NodeJS.Timeout | null = null;
  private healthCheckInterval: NodeJS.Timeout | null = null;
  private startTime: Date;

  constructor(config: ConnectionConfig) {
    super();
    this.config = config;
    this.startTime = new Date();
    this.stats = {
      totalConnections: 0,
      activeConnections: 0,
      failedConnections: 0,
      reconnections: 0,
      uptime: 0
    };
  }

  async connect(): Promise<void> {
    if (this.isConnected || this.isConnecting) {
      return;
    }

    this.isConnecting = true;
    
    try {
      this.client = createClient({
        socket: {
          host: this.config.host,
          port: this.config.port,
          connectTimeout: this.config.connectTimeout || 10000,
          keepAlive: this.config.keepAlive || 30000,
          reconnectStrategy: (retries) => {
            if (retries > 10) {
              this.emit('maxRetriesReached', retries);
              return new Error('Max retries reached');
            }
            const delay = Math.min(retries * 50, 500);
            this.emit('reconnecting', { attempt: retries, delay });
            return delay;
          }
        },
        password: this.config.password,
        database: this.config.database || 0
      });

      this.setupEventHandlers();
      
      await this.client.connect();
      
      this.isConnected = true;
      this.isConnecting = false;
      this.stats.totalConnections++;
      this.stats.activeConnections++;
      this.stats.lastConnectedAt = new Date();
      
      this.startHealthCheck();
      this.emit('connected');
      
    } catch (error) {
      this.isConnecting = false;
      this.stats.failedConnections++;
      this.emit('connectionError', error);
      throw error;
    }
  }

  private setupEventHandlers(): void {
    if (!this.client) return;

    this.client.on('connect', () => {
      this.emit('connecting');
    });

    this.client.on('ready', () => {
      this.isConnected = true;
      this.emit('ready');
    });

    this.client.on('error', (error) => {
      this.isConnected = false;
      this.stats.failedConnections++;
      this.emit('error', error);
    });

    this.client.on('end', () => {
      this.isConnected = false;
      this.stats.activeConnections = Math.max(0, this.stats.activeConnections - 1);
      this.stats.lastDisconnectedAt = new Date();
      this.stopHealthCheck();
      this.emit('disconnected');
    });

    this.client.on('reconnecting', () => {
      this.stats.reconnections++;
      this.emit('reconnecting');
    });
  }

  async disconnect(): Promise<void> {
    if (!this.client || !this.isConnected) {
      return;
    }

    try {
      this.stopHealthCheck();
      await this.client.quit();
      this.isConnected = false;
      this.client = null;
      this.emit('disconnected');
    } catch (error) {
      this.emit('disconnectionError', error);
      throw error;
    }
  }

  async forceDisconnect(): Promise<void> {
    if (!this.client) {
      return;
    }

    try {
      this.stopHealthCheck();
      await this.client.disconnect();
      this.isConnected = false;
      this.client = null;
      this.emit('forceDisconnected');
    } catch (error) {
      this.emit('disconnectionError', error);
      throw error;
    }
  }

  getClient(): RedisClientType {
    if (!this.client || !this.isConnected) {
      throw new Error('Redis client is not connected');
    }
    return this.client;
  }

  isClientConnected(): boolean {
    return this.isConnected && this.client !== null;
  }

  getStats(): ConnectionStats {
    return {
      ...this.stats,
      uptime: Date.now() - this.startTime.getTime()
    };
  }

  private startHealthCheck(): void {
    this.healthCheckInterval = setInterval(async () => {
      try {
        if (this.client && this.isConnected) {
          await this.client.ping();
          this.emit('healthCheckPassed');
        }
      } catch (error) {
        this.emit('healthCheckFailed', error);
      }
    }, 30000); // Check every 30 seconds
  }

  private stopHealthCheck(): void {
    if (this.healthCheckInterval) {
      clearInterval(this.healthCheckInterval);
      this.healthCheckInterval = null;
    }
  }

  async ping(): Promise<string> {
    if (!this.client || !this.isConnected) {
      throw new Error('Redis client is not connected');
    }
    return await this.client.ping();
  }
}

export default ConnectionManager;
```

### Connection Pool Implementation

```typescript
// src/connection/connection-pool.ts
import { createClient, RedisClientType } from 'redis';
import { EventEmitter } from 'events';

interface PoolConfig {
  min: number;
  max: number;
  acquireTimeoutMillis: number;
  idleTimeoutMillis: number;
  host: string;
  port: number;
  password?: string;
  database?: number;
}

interface PoolStats {
  size: number;
  available: number;
  borrowed: number;
  pending: number;
  created: number;
  destroyed: number;
}

interface PooledConnection {
  client: RedisClientType;
  createdAt: Date;
  lastUsedAt: Date;
  isIdle: boolean;
}

class ConnectionPool extends EventEmitter {
  private config: PoolConfig;
  private connections: Map<string, PooledConnection> = new Map();
  private availableConnections: string[] = [];
  private pendingAcquisitions: Array<{
    resolve: (client: RedisClientType) => void;
    reject: (error: Error) => void;
    timeout: NodeJS.Timeout;
  }> = [];
  private idleCheckInterval: NodeJS.Timeout | null = null;

  constructor(config: PoolConfig) {
    super();
    this.config = config;
    this.startIdleCheck();
  }

  async initialize(): Promise<void> {
    // Create minimum number of connections
    const promises = [];
    for (let i = 0; i < this.config.min; i++) {
      promises.push(this.createConnection());
    }
    await Promise.all(promises);
    this.emit('initialized');
  }

  private async createConnection(): Promise<string> {
    const connectionId = `conn_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    try {
      const client = createClient({
        socket: {
          host: this.config.host,
          port: this.config.port
        },
        password: this.config.password,
        database: this.config.database
      });

      await client.connect();

      const connection: PooledConnection = {
        client,
        createdAt: new Date(),
        lastUsedAt: new Date(),
        isIdle: true
      };

      this.connections.set(connectionId, connection);
      this.availableConnections.push(connectionId);
      
      this.emit('connectionCreated', connectionId);
      return connectionId;
    } catch (error) {
      this.emit('connectionCreationFailed', error);
      throw error;
    }
  }

  async acquire(): Promise<RedisClientType> {
    // Check for available connections
    if (this.availableConnections.length > 0) {
      const connectionId = this.availableConnections.shift()!;
      const connection = this.connections.get(connectionId)!;
      connection.isIdle = false;
      connection.lastUsedAt = new Date();
      this.emit('connectionAcquired', connectionId);
      return connection.client;
    }

    // Create new connection if under max limit
    if (this.connections.size < this.config.max) {
      try {
        const connectionId = await this.createConnection();
        const connection = this.connections.get(connectionId)!;
        connection.isIdle = false;
        this.availableConnections.splice(this.availableConnections.indexOf(connectionId), 1);
        this.emit('connectionAcquired', connectionId);
        return connection.client;
      } catch (error) {
        // Fall through to waiting logic
      }
    }

    // Wait for available connection
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        const index = this.pendingAcquisitions.findIndex(p => p.resolve === resolve);
        if (index !== -1) {
          this.pendingAcquisitions.splice(index, 1);
        }
        reject(new Error('Connection acquisition timeout'));
      }, this.config.acquireTimeoutMillis);

      this.pendingAcquisitions.push({ resolve, reject, timeout });
    });
  }

  release(client: RedisClientType): void {
    // Find the connection
    let connectionId: string | null = null;
    for (const [id, connection] of this.connections.entries()) {
      if (connection.client === client) {
        connectionId = id;
        break;
      }
    }

    if (!connectionId) {
      this.emit('releaseError', new Error('Connection not found in pool'));
      return;
    }

    const connection = this.connections.get(connectionId)!;
    connection.isIdle = true;
    connection.lastUsedAt = new Date();

    // Check for pending acquisitions
    if (this.pendingAcquisitions.length > 0) {
      const pending = this.pendingAcquisitions.shift()!;
      clearTimeout(pending.timeout);
      connection.isIdle = false;
      pending.resolve(client);
      this.emit('connectionAcquired', connectionId);
    } else {
      this.availableConnections.push(connectionId);
      this.emit('connectionReleased', connectionId);
    }
  }

  async destroy(client: RedisClientType): Promise<void> {
    let connectionId: string | null = null;
    for (const [id, connection] of this.connections.entries()) {
      if (connection.client === client) {
        connectionId = id;
        break;
      }
    }

    if (!connectionId) {
      return;
    }

    try {
      await client.quit();
    } catch (error) {
      // Ignore errors during destruction
    }

    this.connections.delete(connectionId);
    const availableIndex = this.availableConnections.indexOf(connectionId);
    if (availableIndex !== -1) {
      this.availableConnections.splice(availableIndex, 1);
    }

    this.emit('connectionDestroyed', connectionId);
  }

  async destroyAll(): Promise<void> {
    const destroyPromises = Array.from(this.connections.values()).map(
      connection => this.destroy(connection.client)
    );
    
    await Promise.all(destroyPromises);
    
    // Clear pending acquisitions
    this.pendingAcquisitions.forEach(pending => {
      clearTimeout(pending.timeout);
      pending.reject(new Error('Pool is being destroyed'));
    });
    this.pendingAcquisitions = [];
    
    this.stopIdleCheck();
    this.emit('poolDestroyed');
  }

  private startIdleCheck(): void {
    this.idleCheckInterval = setInterval(() => {
      const now = new Date();
      const connectionsToDestroy: string[] = [];

      for (const [id, connection] of this.connections.entries()) {
        if (connection.isIdle && 
            now.getTime() - connection.lastUsedAt.getTime() > this.config.idleTimeoutMillis &&
            this.connections.size > this.config.min) {
          connectionsToDestroy.push(id);
        }
      }

      connectionsToDestroy.forEach(id => {
        const connection = this.connections.get(id);
        if (connection) {
          this.destroy(connection.client);
        }
      });
    }, 30000); // Check every 30 seconds
  }

  private stopIdleCheck(): void {
    if (this.idleCheckInterval) {
      clearInterval(this.idleCheckInterval);
      this.idleCheckInterval = null;
    }
  }

  getStats(): PoolStats {
    const borrowed = Array.from(this.connections.values())
      .filter(conn => !conn.isIdle).length;
    
    return {
      size: this.connections.size,
      available: this.availableConnections.length,
      borrowed,
      pending: this.pendingAcquisitions.length,
      created: this.connections.size,
      destroyed: 0 // Would need to track this separately
    };
  }
}

export default ConnectionPool;
```

## Error Handling Strategies

### Comprehensive Error Handler

```typescript
// src/error-handling/redis-error-handler.ts
import { RedisClientType } from 'redis';
import { EventEmitter } from 'events';

enum ErrorType {
  CONNECTION_ERROR = 'CONNECTION_ERROR',
  TIMEOUT_ERROR = 'TIMEOUT_ERROR',
  COMMAND_ERROR = 'COMMAND_ERROR',
  MEMORY_ERROR = 'MEMORY_ERROR',
  AUTHENTICATION_ERROR = 'AUTHENTICATION_ERROR',
  UNKNOWN_ERROR = 'UNKNOWN_ERROR'
}

interface ErrorContext {
  operation: string;
  key?: string;
  command?: string;
  timestamp: Date;
  retryCount: number;
}

interface RetryConfig {
  maxRetries: number;
  baseDelay: number;
  maxDelay: number;
  backoffMultiplier: number;
  retryableErrors: ErrorType[];
}

class RedisErrorHandler extends EventEmitter {
  private client: RedisClientType;
  private retryConfig: RetryConfig;
  private errorStats: Map<ErrorType, number> = new Map();

  constructor(client: RedisClientType, retryConfig?: Partial<RetryConfig>) {
    super();
    this.client = client;
    this.retryConfig = {
      maxRetries: 3,
      baseDelay: 1000,
      maxDelay: 10000,
      backoffMultiplier: 2,
      retryableErrors: [
        ErrorType.CONNECTION_ERROR,
        ErrorType.TIMEOUT_ERROR,
        ErrorType.COMMAND_ERROR
      ],
      ...retryConfig
    };
    
    this.initializeErrorStats();
  }

  private initializeErrorStats(): void {
    Object.values(ErrorType).forEach(type => {
      this.errorStats.set(type, 0);
    });
  }

  private classifyError(error: Error): ErrorType {
    const message = error.message.toLowerCase();
    
    if (message.includes('connection') || message.includes('connect')) {
      return ErrorType.CONNECTION_ERROR;
    }
    if (message.includes('timeout')) {
      return ErrorType.TIMEOUT_ERROR;
    }
    if (message.includes('auth') || message.includes('password')) {
      return ErrorType.AUTHENTICATION_ERROR;
    }
    if (message.includes('memory') || message.includes('oom')) {
      return ErrorType.MEMORY_ERROR;
    }
    if (message.includes('command') || message.includes('wrong number of arguments')) {
      return ErrorType.COMMAND_ERROR;
    }
    
    return ErrorType.UNKNOWN_ERROR;
  }

  private incrementErrorCount(errorType: ErrorType): void {
    const current = this.errorStats.get(errorType) || 0;
    this.errorStats.set(errorType, current + 1);
  }

  private calculateDelay(retryCount: number): number {
    const delay = this.retryConfig.baseDelay * 
      Math.pow(this.retryConfig.backoffMultiplier, retryCount);
    return Math.min(delay, this.retryConfig.maxDelay);
  }

  private async delay(ms: number): Promise<void> {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  async executeWithRetry<T>(
    operation: () => Promise<T>,
    context: Omit<ErrorContext, 'timestamp' | 'retryCount'>
  ): Promise<T> {
    let lastError: Error;
    
    for (let attempt = 0; attempt <= this.retryConfig.maxRetries; attempt++) {
      try {
        const result = await operation();
        
        if (attempt > 0) {
          this.emit('retrySuccess', {
            ...context,
            retryCount: attempt,
            timestamp: new Date()
          });
        }
        
        return result;
      } catch (error) {
        lastError = error as Error;
        const errorType = this.classifyError(lastError);
        this.incrementErrorCount(errorType);
        
        const errorContext: ErrorContext = {
          ...context,
          timestamp: new Date(),
          retryCount: attempt
        };
        
        this.emit('error', {
          error: lastError,
          errorType,
          context: errorContext
        });
        
        // Check if error is retryable
        if (!this.retryConfig.retryableErrors.includes(errorType)) {
          this.emit('nonRetryableError', {
            error: lastError,
            errorType,
            context: errorContext
          });
          throw lastError;
        }
        
        // Check if we've exhausted retries
        if (attempt >= this.retryConfig.maxRetries) {
          this.emit('maxRetriesExceeded', {
            error: lastError,
            errorType,
            context: errorContext
          });
          throw new Error(
            `Operation failed after ${this.retryConfig.maxRetries} retries: ${lastError.message}`
          );
        }
        
        // Wait before retry
        const delayMs = this.calculateDelay(attempt);
        this.emit('retryAttempt', {
          error: lastError,
          errorType,
          context: errorContext,
          delayMs,
          nextAttempt: attempt + 1
        });
        
        await this.delay(delayMs);
      }
    }
    
    throw lastError!;
  }

  // Wrapper methods for common Redis operations
  async safeGet(key: string): Promise<string | null> {
    return this.executeWithRetry(
      () => this.client.get(key),
      { operation: 'GET', key, command: 'GET' }
    );
  }

  async safeSet(key: string, value: string, ttl?: number): Promise<string | null> {
    return this.executeWithRetry(
      () => ttl ? this.client.setEx(key, ttl, value) : this.client.set(key, value),
      { operation: 'SET', key, command: ttl ? 'SETEX' : 'SET' }
    );
  }

  async safeDel(key: string): Promise<number> {
    return this.executeWithRetry(
      () => this.client.del(key),
      { operation: 'DELETE', key, command: 'DEL' }
    );
  }

  async safeHGet(key: string, field: string): Promise<string | undefined> {
    return this.executeWithRetry(
      () => this.client.hGet(key, field),
      { operation: 'HGET', key, command: 'HGET' }
    );
  }

  async safeHSet(key: string, field: string, value: string): Promise<number> {
    return this.executeWithRetry(
      () => this.client.hSet(key, field, value),
      { operation: 'HSET', key, command: 'HSET' }
    );
  }

  async safeLPush(key: string, ...values: string[]): Promise<number> {
    return this.executeWithRetry(
      () => this.client.lPush(key, values),
      { operation: 'LPUSH', key, command: 'LPUSH' }
    );
  }

  async safeSAdd(key: string, ...members: string[]): Promise<number> {
    return this.executeWithRetry(
      () => this.client.sAdd(key, members),
      { operation: 'SADD', key, command: 'SADD' }
    );
  }

  getErrorStats(): Map<ErrorType, number> {
    return new Map(this.errorStats);
  }

  resetErrorStats(): void {
    this.initializeErrorStats();
  }

  getRetryConfig(): RetryConfig {
    return { ...this.retryConfig };
  }

  updateRetryConfig(config: Partial<RetryConfig>): void {
    this.retryConfig = { ...this.retryConfig, ...config };
  }
}

export { RedisErrorHandler, ErrorType, ErrorContext, RetryConfig };
```

### Circuit Breaker Implementation

```typescript
// src/error-handling/circuit-breaker.ts
import { EventEmitter } from 'events';

enum CircuitState {
  CLOSED = 'CLOSED',
  OPEN = 'OPEN',
  HALF_OPEN = 'HALF_OPEN'
}

interface CircuitBreakerConfig {
  failureThreshold: number;
  recoveryTimeout: number;
  monitoringPeriod: number;
  halfOpenMaxCalls: number;
  halfOpenSuccessThreshold: number;
}

interface CircuitBreakerStats {
  state: CircuitState;
  failureCount: number;
  successCount: number;
  totalCalls: number;
  lastFailureTime?: Date;
  lastSuccessTime?: Date;
  stateChangedAt: Date;
}

class CircuitBreaker extends EventEmitter {
  private config: CircuitBreakerConfig;
  private state: CircuitState = CircuitState.CLOSED;
  private failureCount: number = 0;
  private successCount: number = 0;
  private totalCalls: number = 0;
  private lastFailureTime?: Date;
  private lastSuccessTime?: Date;
  private stateChangedAt: Date = new Date();
  private halfOpenCalls: number = 0;

  constructor(config: Partial<CircuitBreakerConfig> = {}) {
    super();
    this.config = {
      failureThreshold: 5,
      recoveryTimeout: 60000, // 1 minute
      monitoringPeriod: 10000, // 10 seconds
      halfOpenMaxCalls: 3,
      halfOpenSuccessThreshold: 2,
      ...config
    };
  }

  async execute<T>(operation: () => Promise<T>): Promise<T> {
    this.totalCalls++;

    if (this.state === CircuitState.OPEN) {
      if (this.shouldAttemptReset()) {
        this.setState(CircuitState.HALF_OPEN);
      } else {
        throw new Error('Circuit breaker is OPEN');
      }
    }

    if (this.state === CircuitState.HALF_OPEN) {
      if (this.halfOpenCalls >= this.config.halfOpenMaxCalls) {
        throw new Error('Circuit breaker is HALF_OPEN and max calls exceeded');
      }
      this.halfOpenCalls++;
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
    this.successCount++;
    this.lastSuccessTime = new Date();

    if (this.state === CircuitState.HALF_OPEN) {
      if (this.successCount >= this.config.halfOpenSuccessThreshold) {
        this.setState(CircuitState.CLOSED);
        this.reset();
      }
    } else if (this.state === CircuitState.CLOSED) {
      this.failureCount = 0; // Reset failure count on success
    }
  }

  private onFailure(): void {
    this.failureCount++;
    this.lastFailureTime = new Date();

    if (this.state === CircuitState.HALF_OPEN) {
      this.setState(CircuitState.OPEN);
    } else if (this.state === CircuitState.CLOSED) {
      if (this.failureCount >= this.config.failureThreshold) {
        this.setState(CircuitState.OPEN);
      }
    }
  }

  private shouldAttemptReset(): boolean {
    if (!this.lastFailureTime) {
      return false;
    }
    
    const timeSinceLastFailure = Date.now() - this.lastFailureTime.getTime();
    return timeSinceLastFailure >= this.config.recoveryTimeout;
  }

  private setState(newState: CircuitState): void {
    if (this.state !== newState) {
      const oldState = this.state;
      this.state = newState;
      this.stateChangedAt = new Date();
      
      if (newState === CircuitState.HALF_OPEN) {
        this.halfOpenCalls = 0;
        this.successCount = 0;
      }
      
      this.emit('stateChanged', {
        from: oldState,
        to: newState,
        timestamp: this.stateChangedAt
      });
    }
  }

  private reset(): void {
    this.failureCount = 0;
    this.successCount = 0;
    this.halfOpenCalls = 0;
  }

  getState(): CircuitState {
    return this.state;
  }

  getStats(): CircuitBreakerStats {
    return {
      state: this.state,
      failureCount: this.failureCount,
      successCount: this.successCount,
      totalCalls: this.totalCalls,
      lastFailureTime: this.lastFailureTime,
      lastSuccessTime: this.lastSuccessTime,
      stateChangedAt: this.stateChangedAt
    };
  }

  isCallAllowed(): boolean {
    if (this.state === CircuitState.CLOSED) {
      return true;
    }
    
    if (this.state === CircuitState.OPEN) {
      return this.shouldAttemptReset();
    }
    
    if (this.state === CircuitState.HALF_OPEN) {
      return this.halfOpenCalls < this.config.halfOpenMaxCalls;
    }
    
    return false;
  }

  forceOpen(): void {
    this.setState(CircuitState.OPEN);
  }

  forceClose(): void {
    this.setState(CircuitState.CLOSED);
    this.reset();
  }

  forceHalfOpen(): void {
    this.setState(CircuitState.HALF_OPEN);
  }
}

export { CircuitBreaker, CircuitState, CircuitBreakerConfig, CircuitBreakerStats };
```

## Monitoring and Logging

### Redis Monitor

```typescript
// src/monitoring/redis-monitor.ts
import { RedisClientType } from 'redis';
import { EventEmitter } from 'events';

interface MonitoringMetrics {
  commandsPerSecond: number;
  averageResponseTime: number;
  errorRate: number;
  memoryUsage: {
    used: number;
    peak: number;
    fragmentation: number;
  };
  connections: {
    total: number;
    clients: number;
  };
  keyspaceStats: {
    totalKeys: number;
    expires: number;
    avgTtl: number;
  };
}

interface CommandMetrics {
  command: string;
  count: number;
  totalTime: number;
  averageTime: number;
  errors: number;
}

class RedisMonitor extends EventEmitter {
  private client: RedisClientType;
  private metrics: Map<string, CommandMetrics> = new Map();
  private monitoringInterval: NodeJS.Timeout | null = null;
  private isMonitoring: boolean = false;

  constructor(client: RedisClientType) {
    super();
    this.client = client;
  }

  startMonitoring(intervalMs: number = 5000): void {
    if (this.isMonitoring) {
      return;
    }

    this.isMonitoring = true;
    this.monitoringInterval = setInterval(async () => {
      try {
        const metrics = await this.collectMetrics();
        this.emit('metrics', metrics);
      } catch (error) {
        this.emit('monitoringError', error);
      }
    }, intervalMs);

    this.emit('monitoringStarted');
  }

  stopMonitoring(): void {
    if (!this.isMonitoring) {
      return;
    }

    this.isMonitoring = false;
    if (this.monitoringInterval) {
      clearInterval(this.monitoringInterval);
      this.monitoringInterval = null;
    }

    this.emit('monitoringStopped');
  }

  async collectMetrics(): Promise<MonitoringMetrics> {
    const [info, slowLog] = await Promise.all([
      this.client.info(),
      this.client.slowLogGet(10)
    ]);

    const parsedInfo = this.parseInfo(info);
    
    return {
      commandsPerSecond: this.calculateCommandsPerSecond(parsedInfo),
      averageResponseTime: this.calculateAverageResponseTime(),
      errorRate: this.calculateErrorRate(),
      memoryUsage: {
        used: parseInt(parsedInfo.used_memory || '0'),
        peak: parseInt(parsedInfo.used_memory_peak || '0'),
        fragmentation: parseFloat(parsedInfo.mem_fragmentation_ratio || '1')
      },
      connections: {
        total: parseInt(parsedInfo.total_connections_received || '0'),
        clients: parseInt(parsedInfo.connected_clients || '0')
      },
      keyspaceStats: await this.getKeyspaceStats()
    };
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

  private calculateCommandsPerSecond(info: Record<string, string>): number {
    const totalCommands = parseInt(info.total_commands_processed || '0');
    const uptime = parseInt(info.uptime_in_seconds || '1');
    return totalCommands / uptime;
  }

  private calculateAverageResponseTime(): number {
    let totalTime = 0;
    let totalCommands = 0;

    for (const metric of this.metrics.values()) {
      totalTime += metric.totalTime;
      totalCommands += metric.count;
    }

    return totalCommands > 0 ? totalTime / totalCommands : 0;
  }

  private calculateErrorRate(): number {
    let totalErrors = 0;
    let totalCommands = 0;

    for (const metric of this.metrics.values()) {
      totalErrors += metric.errors;
      totalCommands += metric.count;
    }

    return totalCommands > 0 ? totalErrors / totalCommands : 0;
  }

  private async getKeyspaceStats(): Promise<{ totalKeys: number; expires: number; avgTtl: number }> {
    try {
      const dbSize = await this.client.dbSize();
      // Note: Getting expires and avgTtl would require scanning keys
      // This is a simplified version
      return {
        totalKeys: dbSize,
        expires: 0, // Would need to scan keys to get accurate count
        avgTtl: 0   // Would need to scan keys to calculate average
      };
    } catch (error) {
      return { totalKeys: 0, expires: 0, avgTtl: 0 };
    }
  }

  recordCommand(command: string, executionTime: number, isError: boolean = false): void {
    const existing = this.metrics.get(command) || {
      command,
      count: 0,
      totalTime: 0,
      averageTime: 0,
      errors: 0
    };

    existing.count++;
    existing.totalTime += executionTime;
    existing.averageTime = existing.totalTime / existing.count;
    
    if (isError) {
      existing.errors++;
    }

    this.metrics.set(command, existing);
  }

  getCommandMetrics(): CommandMetrics[] {
    return Array.from(this.metrics.values());
  }

  resetMetrics(): void {
    this.metrics.clear();
  }

  async getSlowQueries(count: number = 10): Promise<any[]> {
    try {
      return await this.client.slowLogGet(count);
    } catch (error) {
      return [];
    }
  }

  async getMemoryUsage(key?: string): Promise<number> {
    try {
      if (key) {
        return await this.client.memoryUsage(key);
      } else {
        const info = await this.client.info('memory');
        const parsedInfo = this.parseInfo(info);
        return parseInt(parsedInfo.used_memory || '0');
      }
    } catch (error) {
      return 0;
    }
  }
}

export { RedisMonitor, MonitoringMetrics, CommandMetrics };
```

## Usage Examples

### Complete Connection Management Example

```typescript
// src/examples/connection-management-example.ts
import ConnectionManager from '../connection/connection-manager';
import ConnectionPool from '../connection/connection-pool';
import { RedisErrorHandler } from '../error-handling/redis-error-handler';
import { CircuitBreaker } from '../error-handling/circuit-breaker';
import { RedisMonitor } from '../monitoring/redis-monitor';

async function connectionManagementExample() {
  // Connection Manager Example
  const connectionManager = new ConnectionManager({
    host: 'localhost',
    port: 6379,
    connectTimeout: 5000,
    keepAlive: 30000
  });

  // Set up event listeners
  connectionManager.on('connected', () => {
    console.log('✅ Connected to Redis');
  });

  connectionManager.on('error', (error) => {
    console.error('❌ Redis error:', error);
  });

  connectionManager.on('reconnecting', (info) => {
    console.log(`🔄 Reconnecting (attempt ${info.attempt})...`);
  });

  try {
    await connectionManager.connect();
    const client = connectionManager.getClient();

    // Error Handler Example
    const errorHandler = new RedisErrorHandler(client, {
      maxRetries: 3,
      baseDelay: 1000,
      backoffMultiplier: 2
    });

    errorHandler.on('retryAttempt', (info) => {
      console.log(`🔄 Retrying operation: ${info.context.operation} (attempt ${info.nextAttempt})`);
    });

    // Circuit Breaker Example
    const circuitBreaker = new CircuitBreaker({
      failureThreshold: 5,
      recoveryTimeout: 30000
    });

    circuitBreaker.on('stateChanged', (info) => {
      console.log(`🔀 Circuit breaker state changed: ${info.from} -> ${info.to}`);
    });

    // Monitor Example
    const monitor = new RedisMonitor(client);
    
    monitor.on('metrics', (metrics) => {
      console.log('📊 Redis Metrics:', {
        commandsPerSecond: metrics.commandsPerSecond.toFixed(2),
        memoryUsed: `${(metrics.memoryUsage.used / 1024 / 1024).toFixed(2)} MB`,
        connections: metrics.connections.clients
      });
    });

    monitor.startMonitoring(10000); // Every 10 seconds

    // Perform operations with error handling and circuit breaker
    const performOperation = async () => {
      return circuitBreaker.execute(async () => {
        return errorHandler.safeSet('test:key', 'test:value', 3600);
      });
    };

    await performOperation();
    
    const value = await circuitBreaker.execute(async () => {
      return errorHandler.safeGet('test:key');
    });
    
    console.log('Retrieved value:', value);

    // Get statistics
    console.log('Connection stats:', connectionManager.getStats());
    console.log('Error stats:', errorHandler.getErrorStats());
    console.log('Circuit breaker stats:', circuitBreaker.getStats());
    console.log('Command metrics:', monitor.getCommandMetrics());

    // Cleanup
    monitor.stopMonitoring();
    await connectionManager.disconnect();

  } catch (error) {
    console.error('Example failed:', error);
  }
}

// Connection Pool Example
async function connectionPoolExample() {
  const pool = new ConnectionPool({
    min: 2,
    max: 10,
    acquireTimeoutMillis: 5000,
    idleTimeoutMillis: 30000,
    host: 'localhost',
    port: 6379
  });

  pool.on('connectionCreated', (id) => {
    console.log(`➕ Connection created: ${id}`);
  });

  pool.on('connectionAcquired', (id) => {
    console.log(`📤 Connection acquired: ${id}`);
  });

  pool.on('connectionReleased', (id) => {
    console.log(`📥 Connection released: ${id}`);
  });

  try {
    await pool.initialize();
    console.log('Pool initialized');

    // Acquire and use connections
    const client1 = await pool.acquire();
    await client1.set('pool:test:1', 'value1');
    pool.release(client1);

    const client2 = await pool.acquire();
    const value = await client2.get('pool:test:1');
    console.log('Retrieved from pool:', value);
    pool.release(client2);

    console.log('Pool stats:', pool.getStats());

    await pool.destroyAll();
    console.log('Pool destroyed');

  } catch (error) {
    console.error('Pool example failed:', error);
  }
}

// Run examples
connectionManagementExample()
  .then(() => connectionPoolExample())
  .catch(console.error);
```

## Next Steps

Now that you understand connection management and error handling, you're ready to explore:

- [Part 7: Pub/Sub Messaging](./part-07-pubsub.md)

## Quick Reference

### Connection Events
```typescript
client.on('connect', () => console.log('Connecting...'));
client.on('ready', () => console.log('Ready'));
client.on('error', (err) => console.error('Error:', err));
client.on('end', () => console.log('Disconnected'));
client.on('reconnecting', () => console.log('Reconnecting...'));
```

### Error Handling Pattern
```typescript
try {
  const result = await client.get('key');
} catch (error) {
  if (error.message.includes('connection')) {
    // Handle connection error
  } else if (error.message.includes('timeout')) {
    // Handle timeout error
  }
}
```

### Retry Logic
```typescript
const executeWithRetry = async (operation, maxRetries = 3) => {
  for (let i = 0; i <= maxRetries; i++) {
    try {
      return await operation();
    } catch (error) {
      if (i === maxRetries) throw error;
      await new Promise(resolve => setTimeout(resolve, 1000 * i));
    }
  }
};
```

---

**Previous**: [← Part 5: CRUD Operations in Node.js](./part-05-crud-operations.md)  
**Next**: [Part 7: Pub/Sub Messaging →](./part-07-pubsub.md)