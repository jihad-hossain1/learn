# System Design with JavaScript & Node.js: Complete Guide

## Table of Contents
1. [Introduction](#introduction)
2. [Part 1: Fundamentals](#part-1-fundamentals)
3. [Part 2: Basic Patterns](#part-2-basic-patterns)
4. [Part 3: Intermediate Concepts](#part-3-intermediate-concepts)
5. [Part 4: Advanced Patterns](#part-4-advanced-patterns)
6. [Part 5: Scalability & Performance](#part-5-scalability--performance)
7. [Part 6: Real-world Applications](#part-6-real-world-applications)
8. [Resources](#resources)

## Introduction

This guide provides a comprehensive journey through system design using JavaScript and the Node.js ecosystem. We'll progress from basic concepts to advanced architectural patterns, covering everything you need to build scalable, maintainable systems.

---

## Part 1: Fundamentals

### 1.1 Understanding System Design

**What is System Design?**
- The process of defining architecture, components, modules, interfaces, and data for a system
- Involves making trade-offs between different quality attributes
- Focus on scalability, reliability, availability, and performance

**Key Principles:**
- **Separation of Concerns**: Each module should have a single responsibility
- **Modularity**: Break down complex systems into smaller, manageable pieces
- **Abstraction**: Hide implementation details behind well-defined interfaces
- **Loose Coupling**: Minimize dependencies between components
- **High Cohesion**: Related functionality should be grouped together

### 1.2 JavaScript & Node.js Fundamentals

**Event Loop & Asynchronous Programming**
```javascript
// Understanding the Event Loop
console.log('Start');

setTimeout(() => {
  console.log('Timeout');
}, 0);

setImmediate(() => {
  console.log('Immediate');
});

process.nextTick(() => {
  console.log('Next Tick');
});

console.log('End');
// Output: Start, End, Next Tick, Immediate, Timeout
```

**Promises and Async/Await**
```javascript
// Promise-based approach
function fetchData() {
  return new Promise((resolve, reject) => {
    // Simulate API call
    setTimeout(() => {
      resolve({ data: 'Sample data' });
    }, 1000);
  });
}

// Async/Await approach
async function getData() {
  try {
    const result = await fetchData();
    return result;
  } catch (error) {
    console.error('Error:', error);
  }
}
```

### 1.3 Core Node.js Concepts

**Modules and CommonJS**
```javascript
// math.js
function add(a, b) {
  return a + b;
}

function subtract(a, b) {
  return a - b;
}

module.exports = { add, subtract };

// app.js
const { add, subtract } = require('./math');
console.log(add(5, 3)); // 8
```

**ES6 Modules**
```javascript
// math.mjs
export function add(a, b) {
  return a + b;
}

export function subtract(a, b) {
  return a - b;
}

// app.mjs
import { add, subtract } from './math.mjs';
console.log(add(5, 3)); // 8
```

---

## Part 2: Basic Patterns

### 2.1 Module Pattern

**Revealing Module Pattern**
```javascript
const UserModule = (function() {
  let users = [];
  
  function addUser(user) {
    users.push(user);
  }
  
  function getUsers() {
    return [...users]; // Return copy to prevent mutation
  }
  
  function findUser(id) {
    return users.find(user => user.id === id);
  }
  
  // Public API
  return {
    add: addUser,
    getAll: getUsers,
    find: findUser
  };
})();
```

### 2.2 Observer Pattern

**Event Emitter Implementation**
```javascript
const EventEmitter = require('events');

class UserService extends EventEmitter {
  constructor() {
    super();
    this.users = [];
  }
  
  addUser(user) {
    this.users.push(user);
    this.emit('userAdded', user);
  }
  
  removeUser(userId) {
    const index = this.users.findIndex(u => u.id === userId);
    if (index !== -1) {
      const user = this.users.splice(index, 1)[0];
      this.emit('userRemoved', user);
    }
  }
}

// Usage
const userService = new UserService();

userService.on('userAdded', (user) => {
  console.log(`User ${user.name} was added`);
});

userService.on('userRemoved', (user) => {
  console.log(`User ${user.name} was removed`);
});
```

### 2.3 Factory Pattern

**Database Connection Factory**
```javascript
class DatabaseFactory {
  static createConnection(type, config) {
    switch (type) {
      case 'mysql':
        return new MySQLConnection(config);
      case 'postgresql':
        return new PostgreSQLConnection(config);
      case 'mongodb':
        return new MongoDBConnection(config);
      default:
        throw new Error(`Unsupported database type: ${type}`);
    }
  }
}

class MySQLConnection {
  constructor(config) {
    this.config = config;
  }
  
  connect() {
    console.log('Connecting to MySQL...');
  }
}

class PostgreSQLConnection {
  constructor(config) {
    this.config = config;
  }
  
  connect() {
    console.log('Connecting to PostgreSQL...');
  }
}
```

### 2.4 Singleton Pattern

**Configuration Manager**
```javascript
class ConfigManager {
  constructor() {
    if (ConfigManager.instance) {
      return ConfigManager.instance;
    }
    
    this.config = {};
    ConfigManager.instance = this;
  }
  
  set(key, value) {
    this.config[key] = value;
  }
  
  get(key) {
    return this.config[key];
  }
  
  getAll() {
    return { ...this.config };
  }
}

// Usage
const config1 = new ConfigManager();
const config2 = new ConfigManager();

console.log(config1 === config2); // true
```

---

## Part 3: Intermediate Concepts

### 3.1 Layered Architecture

**Three-Layer Architecture**
```javascript
// models/User.js
class User {
  constructor(id, name, email) {
    this.id = id;
    this.name = name;
    this.email = email;
  }
  
  validate() {
    if (!this.name || !this.email) {
      throw new Error('Name and email are required');
    }
    
    if (!this.email.includes('@')) {
      throw new Error('Invalid email format');
    }
  }
}

// repositories/UserRepository.js
class UserRepository {
  constructor(database) {
    this.db = database;
  }
  
  async create(user) {
    const query = 'INSERT INTO users (name, email) VALUES (?, ?)';
    const result = await this.db.execute(query, [user.name, user.email]);
    return result.insertId;
  }
  
  async findById(id) {
    const query = 'SELECT * FROM users WHERE id = ?';
    const [rows] = await this.db.execute(query, [id]);
    return rows[0] ? new User(rows[0].id, rows[0].name, rows[0].email) : null;
  }
  
  async findAll() {
    const query = 'SELECT * FROM users';
    const [rows] = await this.db.execute(query);
    return rows.map(row => new User(row.id, row.name, row.email));
  }
}

// services/UserService.js
class UserService {
  constructor(userRepository) {
    this.userRepository = userRepository;
  }
  
  async createUser(userData) {
    const user = new User(null, userData.name, userData.email);
    user.validate();
    
    const userId = await this.userRepository.create(user);
    return await this.userRepository.findById(userId);
  }
  
  async getUserById(id) {
    if (!id || id <= 0) {
      throw new Error('Invalid user ID');
    }
    
    return await this.userRepository.findById(id);
  }
  
  async getAllUsers() {
    return await this.userRepository.findAll();
  }
}

// controllers/UserController.js
class UserController {
  constructor(userService) {
    this.userService = userService;
  }
  
  async createUser(req, res) {
    try {
      const user = await this.userService.createUser(req.body);
      res.status(201).json(user);
    } catch (error) {
      res.status(400).json({ error: error.message });
    }
  }
  
  async getUser(req, res) {
    try {
      const user = await this.userService.getUserById(req.params.id);
      if (!user) {
        return res.status(404).json({ error: 'User not found' });
      }
      res.json(user);
    } catch (error) {
      res.status(400).json({ error: error.message });
    }
  }
}
```

### 3.2 Dependency Injection

**Simple DI Container**
```javascript
class DIContainer {
  constructor() {
    this.services = new Map();
    this.singletons = new Map();
  }
  
  register(name, factory, options = {}) {
    this.services.set(name, {
      factory,
      singleton: options.singleton || false
    });
  }
  
  resolve(name) {
    const service = this.services.get(name);
    
    if (!service) {
      throw new Error(`Service ${name} not found`);
    }
    
    if (service.singleton) {
      if (!this.singletons.has(name)) {
        this.singletons.set(name, service.factory(this));
      }
      return this.singletons.get(name);
    }
    
    return service.factory(this);
  }
}

// Usage
const container = new DIContainer();

// Register services
container.register('database', () => new Database(), { singleton: true });
container.register('userRepository', (container) => 
  new UserRepository(container.resolve('database'))
);
container.register('userService', (container) => 
  new UserService(container.resolve('userRepository'))
);
container.register('userController', (container) => 
  new UserController(container.resolve('userService'))
);

// Resolve dependencies
const userController = container.resolve('userController');
```

### 3.3 Middleware Pattern

**Express-style Middleware**
```javascript
class MiddlewareStack {
  constructor() {
    this.middlewares = [];
  }
  
  use(middleware) {
    this.middlewares.push(middleware);
  }
  
  async execute(context) {
    let index = 0;
    
    const next = async () => {
      if (index < this.middlewares.length) {
        const middleware = this.middlewares[index++];
        await middleware(context, next);
      }
    };
    
    await next();
  }
}

// Middleware functions
const loggingMiddleware = async (context, next) => {
  console.log(`${new Date().toISOString()} - ${context.method} ${context.url}`);
  await next();
};

const authMiddleware = async (context, next) => {
  if (!context.headers.authorization) {
    throw new Error('Unauthorized');
  }
  context.user = { id: 1, name: 'John Doe' }; // Mock user
  await next();
};

const validationMiddleware = async (context, next) => {
  if (context.method === 'POST' && !context.body) {
    throw new Error('Request body is required');
  }
  await next();
};

// Usage
const stack = new MiddlewareStack();
stack.use(loggingMiddleware);
stack.use(authMiddleware);
stack.use(validationMiddleware);

// Execute middleware stack
const context = {
  method: 'POST',
  url: '/api/users',
  headers: { authorization: 'Bearer token' },
  body: { name: 'John', email: 'john@example.com' }
};

stack.execute(context);
```

---

## Part 4: Advanced Patterns

### 4.1 CQRS (Command Query Responsibility Segregation)

**CQRS Implementation**
```javascript
// Commands
class CreateUserCommand {
  constructor(name, email) {
    this.name = name;
    this.email = email;
  }
}

class UpdateUserCommand {
  constructor(id, name, email) {
    this.id = id;
    this.name = name;
    this.email = email;
  }
}

// Command Handlers
class CreateUserCommandHandler {
  constructor(userRepository, eventBus) {
    this.userRepository = userRepository;
    this.eventBus = eventBus;
  }
  
  async handle(command) {
    const user = new User(null, command.name, command.email);
    user.validate();
    
    const userId = await this.userRepository.create(user);
    
    // Publish event
    this.eventBus.publish(new UserCreatedEvent(userId, command.name, command.email));
    
    return userId;
  }
}

// Queries
class GetUserQuery {
  constructor(id) {
    this.id = id;
  }
}

class GetAllUsersQuery {
  constructor(filters = {}) {
    this.filters = filters;
  }
}

// Query Handlers
class GetUserQueryHandler {
  constructor(userReadRepository) {
    this.userReadRepository = userReadRepository;
  }
  
  async handle(query) {
    return await this.userReadRepository.findById(query.id);
  }
}

// Command/Query Bus
class CommandBus {
  constructor() {
    this.handlers = new Map();
  }
  
  register(commandType, handler) {
    this.handlers.set(commandType, handler);
  }
  
  async execute(command) {
    const handler = this.handlers.get(command.constructor);
    if (!handler) {
      throw new Error(`No handler registered for ${command.constructor.name}`);
    }
    return await handler.handle(command);
  }
}

class QueryBus {
  constructor() {
    this.handlers = new Map();
  }
  
  register(queryType, handler) {
    this.handlers.set(queryType, handler);
  }
  
  async execute(query) {
    const handler = this.handlers.get(query.constructor);
    if (!handler) {
      throw new Error(`No handler registered for ${query.constructor.name}`);
    }
    return await handler.handle(query);
  }
}
```

### 4.2 Event Sourcing

**Event Store Implementation**
```javascript
class Event {
  constructor(aggregateId, eventType, data, version) {
    this.aggregateId = aggregateId;
    this.eventType = eventType;
    this.data = data;
    this.version = version;
    this.timestamp = new Date();
  }
}

class EventStore {
  constructor() {
    this.events = [];
  }
  
  async saveEvents(aggregateId, events, expectedVersion) {
    const existingEvents = this.events.filter(e => e.aggregateId === aggregateId);
    const currentVersion = existingEvents.length;
    
    if (expectedVersion !== currentVersion) {
      throw new Error('Concurrency conflict');
    }
    
    events.forEach((event, index) => {
      event.version = currentVersion + index + 1;
      this.events.push(event);
    });
  }
  
  async getEvents(aggregateId, fromVersion = 0) {
    return this.events
      .filter(e => e.aggregateId === aggregateId && e.version > fromVersion)
      .sort((a, b) => a.version - b.version);
  }
}

// Aggregate Root
class UserAggregate {
  constructor(id) {
    this.id = id;
    this.version = 0;
    this.uncommittedEvents = [];
    this.name = null;
    this.email = null;
    this.isActive = false;
  }
  
  static fromHistory(events) {
    const aggregate = new UserAggregate(events[0].aggregateId);
    events.forEach(event => aggregate.apply(event));
    return aggregate;
  }
  
  createUser(name, email) {
    if (this.name) {
      throw new Error('User already exists');
    }
    
    this.raiseEvent(new UserCreatedEvent(this.id, name, email));
  }
  
  updateUser(name, email) {
    if (!this.isActive) {
      throw new Error('Cannot update inactive user');
    }
    
    this.raiseEvent(new UserUpdatedEvent(this.id, name, email));
  }
  
  deactivateUser() {
    if (!this.isActive) {
      throw new Error('User is already inactive');
    }
    
    this.raiseEvent(new UserDeactivatedEvent(this.id));
  }
  
  raiseEvent(event) {
    this.apply(event);
    this.uncommittedEvents.push(event);
  }
  
  apply(event) {
    switch (event.constructor) {
      case UserCreatedEvent:
        this.name = event.data.name;
        this.email = event.data.email;
        this.isActive = true;
        break;
      case UserUpdatedEvent:
        this.name = event.data.name;
        this.email = event.data.email;
        break;
      case UserDeactivatedEvent:
        this.isActive = false;
        break;
    }
    this.version++;
  }
  
  getUncommittedEvents() {
    return [...this.uncommittedEvents];
  }
  
  markEventsAsCommitted() {
    this.uncommittedEvents = [];
  }
}

// Events
class UserCreatedEvent extends Event {
  constructor(aggregateId, name, email) {
    super(aggregateId, 'UserCreated', { name, email });
  }
}

class UserUpdatedEvent extends Event {
  constructor(aggregateId, name, email) {
    super(aggregateId, 'UserUpdated', { name, email });
  }
}

class UserDeactivatedEvent extends Event {
  constructor(aggregateId) {
    super(aggregateId, 'UserDeactivated', {});
  }
}
```

### 4.3 Microservices Communication

**Message Queue Implementation**
```javascript
class MessageQueue {
  constructor() {
    this.queues = new Map();
    this.subscribers = new Map();
  }
  
  createQueue(queueName) {
    if (!this.queues.has(queueName)) {
      this.queues.set(queueName, []);
      this.subscribers.set(queueName, []);
    }
  }
  
  publish(queueName, message) {
    if (!this.queues.has(queueName)) {
      this.createQueue(queueName);
    }
    
    const queue = this.queues.get(queueName);
    queue.push({
      id: Date.now() + Math.random(),
      data: message,
      timestamp: new Date()
    });
    
    this.processQueue(queueName);
  }
  
  subscribe(queueName, handler) {
    if (!this.subscribers.has(queueName)) {
      this.createQueue(queueName);
    }
    
    const subscribers = this.subscribers.get(queueName);
    subscribers.push(handler);
  }
  
  async processQueue(queueName) {
    const queue = this.queues.get(queueName);
    const subscribers = this.subscribers.get(queueName);
    
    while (queue.length > 0 && subscribers.length > 0) {
      const message = queue.shift();
      
      for (const handler of subscribers) {
        try {
          await handler(message);
        } catch (error) {
          console.error('Error processing message:', error);
          // Implement retry logic or dead letter queue
        }
      }
    }
  }
}

// Service Communication
class UserService {
  constructor(messageQueue) {
    this.messageQueue = messageQueue;
    this.setupSubscriptions();
  }
  
  setupSubscriptions() {
    this.messageQueue.subscribe('user.created', this.handleUserCreated.bind(this));
    this.messageQueue.subscribe('user.updated', this.handleUserUpdated.bind(this));
  }
  
  async createUser(userData) {
    // Create user logic
    const user = await this.userRepository.create(userData);
    
    // Publish event
    this.messageQueue.publish('user.created', {
      userId: user.id,
      name: user.name,
      email: user.email
    });
    
    return user;
  }
  
  async handleUserCreated(message) {
    console.log('User created event received:', message.data);
    // Handle user creation side effects
  }
  
  async handleUserUpdated(message) {
    console.log('User updated event received:', message.data);
    // Handle user update side effects
  }
}

class NotificationService {
  constructor(messageQueue) {
    this.messageQueue = messageQueue;
    this.setupSubscriptions();
  }
  
  setupSubscriptions() {
    this.messageQueue.subscribe('user.created', this.sendWelcomeEmail.bind(this));
  }
  
  async sendWelcomeEmail(message) {
    const { email, name } = message.data;
    console.log(`Sending welcome email to ${email}`);
    // Email sending logic
  }
}
```

---

## Part 5: Scalability & Performance

### 5.1 Caching Strategies

**Multi-level Caching**
```javascript
class CacheManager {
  constructor() {
    this.memoryCache = new Map();
    this.redisClient = null; // Redis client would be initialized here
  }
  
  async get(key) {
    // Level 1: Memory cache
    if (this.memoryCache.has(key)) {
      return this.memoryCache.get(key);
    }
    
    // Level 2: Redis cache
    if (this.redisClient) {
      const value = await this.redisClient.get(key);
      if (value) {
        const parsed = JSON.parse(value);
        this.memoryCache.set(key, parsed); // Populate memory cache
        return parsed;
      }
    }
    
    return null;
  }
  
  async set(key, value, ttl = 3600) {
    // Set in memory cache
    this.memoryCache.set(key, value);
    
    // Set in Redis cache
    if (this.redisClient) {
      await this.redisClient.setex(key, ttl, JSON.stringify(value));
    }
    
    // Set TTL for memory cache
    setTimeout(() => {
      this.memoryCache.delete(key);
    }, ttl * 1000);
  }
  
  async invalidate(key) {
    this.memoryCache.delete(key);
    if (this.redisClient) {
      await this.redisClient.del(key);
    }
  }
}

// Cache-aside pattern
class UserService {
  constructor(userRepository, cacheManager) {
    this.userRepository = userRepository;
    this.cache = cacheManager;
  }
  
  async getUserById(id) {
    const cacheKey = `user:${id}`;
    
    // Try cache first
    let user = await this.cache.get(cacheKey);
    if (user) {
      return user;
    }
    
    // Fetch from database
    user = await this.userRepository.findById(id);
    if (user) {
      // Cache the result
      await this.cache.set(cacheKey, user, 1800); // 30 minutes
    }
    
    return user;
  }
  
  async updateUser(id, userData) {
    const user = await this.userRepository.update(id, userData);
    
    // Invalidate cache
    await this.cache.invalidate(`user:${id}`);
    
    return user;
  }
}
```

### 5.2 Connection Pooling

**Database Connection Pool**
```javascript
class ConnectionPool {
  constructor(config) {
    this.config = config;
    this.pool = [];
    this.activeConnections = 0;
    this.waitingQueue = [];
  }
  
  async getConnection() {
    return new Promise((resolve, reject) => {
      // Check if there's an available connection
      if (this.pool.length > 0) {
        const connection = this.pool.pop();
        this.activeConnections++;
        resolve(connection);
        return;
      }
      
      // Check if we can create a new connection
      if (this.activeConnections < this.config.maxConnections) {
        this.createConnection()
          .then(connection => {
            this.activeConnections++;
            resolve(connection);
          })
          .catch(reject);
        return;
      }
      
      // Add to waiting queue
      this.waitingQueue.push({ resolve, reject });
      
      // Set timeout
      setTimeout(() => {
        const index = this.waitingQueue.findIndex(item => item.resolve === resolve);
        if (index !== -1) {
          this.waitingQueue.splice(index, 1);
          reject(new Error('Connection timeout'));
        }
      }, this.config.acquireTimeout || 10000);
    });
  }
  
  releaseConnection(connection) {
    this.activeConnections--;
    
    // Check if there are waiting requests
    if (this.waitingQueue.length > 0) {
      const { resolve } = this.waitingQueue.shift();
      this.activeConnections++;
      resolve(connection);
      return;
    }
    
    // Return connection to pool
    if (this.pool.length < this.config.minConnections) {
      this.pool.push(connection);
    } else {
      connection.close();
    }
  }
  
  async createConnection() {
    // Mock connection creation
    return {
      id: Date.now() + Math.random(),
      query: async (sql, params) => {
        // Mock query execution
        return { rows: [], affectedRows: 0 };
      },
      close: () => {
        console.log('Connection closed');
      }
    };
  }
}
```

### 5.3 Load Balancing

**Round Robin Load Balancer**
```javascript
class LoadBalancer {
  constructor(servers) {
    this.servers = servers;
    this.currentIndex = 0;
    this.healthCheck();
  }
  
  getNextServer() {
    const availableServers = this.servers.filter(server => server.healthy);
    
    if (availableServers.length === 0) {
      throw new Error('No healthy servers available');
    }
    
    const server = availableServers[this.currentIndex % availableServers.length];
    this.currentIndex++;
    
    return server;
  }
  
  async makeRequest(path, options = {}) {
    const maxRetries = 3;
    let lastError;
    
    for (let i = 0; i < maxRetries; i++) {
      try {
        const server = this.getNextServer();
        const response = await this.sendRequest(server, path, options);
        return response;
      } catch (error) {
        lastError = error;
        console.log(`Request failed, retrying... (${i + 1}/${maxRetries})`);
      }
    }
    
    throw lastError;
  }
  
  async sendRequest(server, path, options) {
    // Mock HTTP request
    if (!server.healthy) {
      throw new Error('Server is not healthy');
    }
    
    return {
      status: 200,
      data: { message: `Response from ${server.url}${path}` }
    };
  }
  
  healthCheck() {
    setInterval(async () => {
      for (const server of this.servers) {
        try {
          await this.sendRequest(server, '/health', {});
          server.healthy = true;
        } catch (error) {
          server.healthy = false;
          console.log(`Server ${server.url} is unhealthy`);
        }
      }
    }, 30000); // Check every 30 seconds
  }
}

// Usage
const servers = [
  { url: 'http://server1:3000', healthy: true },
  { url: 'http://server2:3000', healthy: true },
  { url: 'http://server3:3000', healthy: true }
];

const loadBalancer = new LoadBalancer(servers);
```

### 5.4 Rate Limiting

**Token Bucket Rate Limiter**
```javascript
class TokenBucket {
  constructor(capacity, refillRate) {
    this.capacity = capacity;
    this.tokens = capacity;
    this.refillRate = refillRate;
    this.lastRefill = Date.now();
  }
  
  consume(tokens = 1) {
    this.refill();
    
    if (this.tokens >= tokens) {
      this.tokens -= tokens;
      return true;
    }
    
    return false;
  }
  
  refill() {
    const now = Date.now();
    const timePassed = (now - this.lastRefill) / 1000;
    const tokensToAdd = timePassed * this.refillRate;
    
    this.tokens = Math.min(this.capacity, this.tokens + tokensToAdd);
    this.lastRefill = now;
  }
}

class RateLimiter {
  constructor() {
    this.buckets = new Map();
  }
  
  isAllowed(identifier, capacity = 10, refillRate = 1) {
    if (!this.buckets.has(identifier)) {
      this.buckets.set(identifier, new TokenBucket(capacity, refillRate));
    }
    
    const bucket = this.buckets.get(identifier);
    return bucket.consume();
  }
  
  // Middleware for Express
  middleware(options = {}) {
    const { capacity = 100, refillRate = 10, keyGenerator = (req) => req.ip } = options;
    
    return (req, res, next) => {
      const key = keyGenerator(req);
      
      if (this.isAllowed(key, capacity, refillRate)) {
        next();
      } else {
        res.status(429).json({
          error: 'Too Many Requests',
          message: 'Rate limit exceeded'
        });
      }
    };
  }
}
```

---

## Part 6: Real-world Applications

### 6.1 E-commerce System Architecture

**Microservices Structure**
```javascript
// Product Service
class ProductService {
  constructor(productRepository, cacheManager, eventBus) {
    this.productRepository = productRepository;
    this.cache = cacheManager;
    this.eventBus = eventBus;
  }
  
  async getProduct(id) {
    const cacheKey = `product:${id}`;
    let product = await this.cache.get(cacheKey);
    
    if (!product) {
      product = await this.productRepository.findById(id);
      if (product) {
        await this.cache.set(cacheKey, product, 3600);
      }
    }
    
    return product;
  }
  
  async updateInventory(productId, quantity) {
    const product = await this.productRepository.findById(productId);
    
    if (product.inventory < quantity) {
      throw new Error('Insufficient inventory');
    }
    
    await this.productRepository.updateInventory(productId, -quantity);
    await this.cache.invalidate(`product:${productId}`);
    
    this.eventBus.publish('inventory.updated', {
      productId,
      newQuantity: product.inventory - quantity
    });
  }
}

// Order Service
class OrderService {
  constructor(orderRepository, productService, paymentService, eventBus) {
    this.orderRepository = orderRepository;
    this.productService = productService;
    this.paymentService = paymentService;
    this.eventBus = eventBus;
  }
  
  async createOrder(orderData) {
    const { userId, items } = orderData;
    
    // Validate inventory
    for (const item of items) {
      const product = await this.productService.getProduct(item.productId);
      if (product.inventory < item.quantity) {
        throw new Error(`Insufficient inventory for product ${item.productId}`);
      }
    }
    
    // Calculate total
    let total = 0;
    for (const item of items) {
      const product = await this.productService.getProduct(item.productId);
      total += product.price * item.quantity;
    }
    
    // Create order
    const order = await this.orderRepository.create({
      userId,
      items,
      total,
      status: 'pending'
    });
    
    // Reserve inventory
    for (const item of items) {
      await this.productService.updateInventory(item.productId, item.quantity);
    }
    
    // Process payment
    try {
      await this.paymentService.processPayment(order.id, total);
      await this.orderRepository.updateStatus(order.id, 'paid');
      
      this.eventBus.publish('order.created', {
        orderId: order.id,
        userId,
        total
      });
    } catch (error) {
      // Rollback inventory
      for (const item of items) {
        await this.productService.updateInventory(item.productId, -item.quantity);
      }
      throw error;
    }
    
    return order;
  }
}
```

### 6.2 Real-time Chat System

**WebSocket Implementation**
```javascript
const WebSocket = require('ws');

class ChatServer {
  constructor() {
    this.rooms = new Map();
    this.users = new Map();
    this.wss = new WebSocket.Server({ port: 8080 });
    this.setupWebSocketServer();
  }
  
  setupWebSocketServer() {
    this.wss.on('connection', (ws) => {
      ws.on('message', (data) => {
        try {
          const message = JSON.parse(data);
          this.handleMessage(ws, message);
        } catch (error) {
          ws.send(JSON.stringify({ error: 'Invalid message format' }));
        }
      });
      
      ws.on('close', () => {
        this.handleDisconnection(ws);
      });
    });
  }
  
  handleMessage(ws, message) {
    switch (message.type) {
      case 'join':
        this.handleJoin(ws, message);
        break;
      case 'message':
        this.handleChatMessage(ws, message);
        break;
      case 'leave':
        this.handleLeave(ws, message);
        break;
    }
  }
  
  handleJoin(ws, message) {
    const { roomId, userId, username } = message;
    
    // Add user to room
    if (!this.rooms.has(roomId)) {
      this.rooms.set(roomId, new Set());
    }
    
    this.rooms.get(roomId).add(ws);
    this.users.set(ws, { userId, username, roomId });
    
    // Notify room members
    this.broadcastToRoom(roomId, {
      type: 'user_joined',
      userId,
      username,
      timestamp: new Date().toISOString()
    }, ws);
    
    // Send confirmation to user
    ws.send(JSON.stringify({
      type: 'joined',
      roomId,
      message: `Welcome to room ${roomId}`
    }));
  }
  
  handleChatMessage(ws, message) {
    const user = this.users.get(ws);
    if (!user) {
      ws.send(JSON.stringify({ error: 'Not authenticated' }));
      return;
    }
    
    const chatMessage = {
      type: 'message',
      userId: user.userId,
      username: user.username,
      content: message.content,
      timestamp: new Date().toISOString()
    };
    
    this.broadcastToRoom(user.roomId, chatMessage);
  }
  
  handleLeave(ws, message) {
    const user = this.users.get(ws);
    if (user) {
      this.removeUserFromRoom(ws, user);
    }
  }
  
  handleDisconnection(ws) {
    const user = this.users.get(ws);
    if (user) {
      this.removeUserFromRoom(ws, user);
    }
  }
  
  removeUserFromRoom(ws, user) {
    const room = this.rooms.get(user.roomId);
    if (room) {
      room.delete(ws);
      if (room.size === 0) {
        this.rooms.delete(user.roomId);
      } else {
        this.broadcastToRoom(user.roomId, {
          type: 'user_left',
          userId: user.userId,
          username: user.username,
          timestamp: new Date().toISOString()
        });
      }
    }
    this.users.delete(ws);
  }
  
  broadcastToRoom(roomId, message, excludeWs = null) {
    const room = this.rooms.get(roomId);
    if (room) {
      const messageStr = JSON.stringify(message);
      room.forEach(ws => {
        if (ws !== excludeWs && ws.readyState === WebSocket.OPEN) {
          ws.send(messageStr);
        }
      });
    }
  }
}

// Message persistence
class MessageStore {
  constructor() {
    this.messages = new Map(); // roomId -> messages array
  }
  
  saveMessage(roomId, message) {
    if (!this.messages.has(roomId)) {
      this.messages.set(roomId, []);
    }
    
    this.messages.get(roomId).push({
      ...message,
      id: Date.now() + Math.random()
    });
  }
  
  getMessages(roomId, limit = 50, offset = 0) {
    const roomMessages = this.messages.get(roomId) || [];
    return roomMessages
      .slice(-limit - offset, -offset || undefined)
      .reverse();
  }
}
```

### 6.3 API Gateway Pattern

**Gateway Implementation**
```javascript
class APIGateway {
  constructor() {
    this.services = new Map();
    this.middleware = [];
    this.rateLimiter = new RateLimiter();
  }
  
  registerService(name, config) {
    this.services.set(name, {
      baseUrl: config.baseUrl,
      loadBalancer: new LoadBalancer(config.instances),
      timeout: config.timeout || 5000
    });
  }
  
  use(middleware) {
    this.middleware.push(middleware);
  }
  
  async handleRequest(req, res) {
    try {
      // Apply middleware
      for (const middleware of this.middleware) {
        await middleware(req, res);
        if (res.headersSent) return;
      }
      
      // Rate limiting
      if (!this.rateLimiter.isAllowed(req.ip, 100, 10)) {
        return res.status(429).json({ error: 'Rate limit exceeded' });
      }
      
      // Route to service
      const route = this.parseRoute(req.path);
      const service = this.services.get(route.serviceName);
      
      if (!service) {
        return res.status(404).json({ error: 'Service not found' });
      }
      
      const response = await this.forwardRequest(service, route.path, req);
      
      res.status(response.status).json(response.data);
    } catch (error) {
      console.error('Gateway error:', error);
      res.status(500).json({ error: 'Internal server error' });
    }
  }
  
  parseRoute(path) {
    // Parse /api/v1/users/123 -> { serviceName: 'users', path: '/123' }
    const parts = path.split('/').filter(Boolean);
    if (parts.length < 3) {
      throw new Error('Invalid route format');
    }
    
    return {
      serviceName: parts[2], // users, orders, etc.
      path: '/' + parts.slice(3).join('/')
    };
  }
  
  async forwardRequest(service, path, req) {
    const server = service.loadBalancer.getNextServer();
    const url = `${server.url}${path}`;
    
    // Mock HTTP request forwarding
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error('Request timeout'));
      }, service.timeout);
      
      // Simulate request
      setTimeout(() => {
        clearTimeout(timeout);
        resolve({
          status: 200,
          data: { message: `Response from ${url}`, method: req.method }
        });
      }, Math.random() * 1000);
    });
  }
}

// Authentication middleware
const authMiddleware = async (req, res) => {
  const token = req.headers.authorization;
  
  if (!token) {
    res.status(401).json({ error: 'Authentication required' });
    return;
  }
  
  try {
    // Verify JWT token (mock)
    const decoded = { userId: 123, role: 'user' };
    req.user = decoded;
  } catch (error) {
    res.status(401).json({ error: 'Invalid token' });
  }
};

// Logging middleware
const loggingMiddleware = async (req, res) => {
  console.log(`${new Date().toISOString()} - ${req.method} ${req.path}`);
};

// Usage
const gateway = new APIGateway();

gateway.use(loggingMiddleware);
gateway.use(authMiddleware);

gateway.registerService('users', {
  baseUrl: 'http://user-service',
  instances: [
    { url: 'http://user-service-1:3000', healthy: true },
    { url: 'http://user-service-2:3000', healthy: true }
  ]
});

gateway.registerService('orders', {
  baseUrl: 'http://order-service',
  instances: [
    { url: 'http://order-service-1:3000', healthy: true }
  ]
});
```

---

## Resources

### Books
- "Designing Data-Intensive Applications" by Martin Kleppmann
- "Building Microservices" by Sam Newman
- "Clean Architecture" by Robert C. Martin
- "Node.js Design Patterns" by Mario Casciaro

### Online Resources
- [Node.js Best Practices](https://github.com/goldbergyoni/nodebestpractices)
- [System Design Primer](https://github.com/donnemartin/system-design-primer)
- [Microservices.io](https://microservices.io/)
- [High Scalability](http://highscalability.com/)

### Tools & Frameworks
- **API Development**: Express.js, Fastify, Koa.js
- **Database**: MongoDB, PostgreSQL, Redis
- **Message Queues**: RabbitMQ, Apache Kafka, Bull Queue
- **Monitoring**: Prometheus, Grafana, ELK Stack
- **Container Orchestration**: Docker, Kubernetes
- **Testing**: Jest, Mocha, Supertest

### Practice Projects
1. **URL Shortener** (like bit.ly)
2. **Chat Application** with real-time messaging
3. **E-commerce Platform** with microservices
4. **Social Media Feed** with caching and pagination
5. **File Storage Service** with CDN integration

---

## Conclusion

This guide provides a comprehensive foundation for understanding system design with JavaScript and Node.js. Start with the fundamentals and gradually work your way through the advanced patterns. Practice implementing these concepts in real projects to solidify your understanding.

Remember that system design is about making informed trade-offs based on your specific requirements. There's no one-size-fits-all solution, so always consider your use case, scale, team size, and constraints when making architectural decisions.

Good luck on your system design journey!