# Advanced Networking Concepts for Full Stack Development

## Table of Contents
1. [Microservices Communication](#microservices-communication)
2. [Message Queues and Pub/Sub](#message-queues-and-pubsub)
3. [Service Discovery](#service-discovery)
4. [Circuit Breakers and Resilience](#circuit-breakers-and-resilience)
5. [Real-time Communication](#real-time-communication)
6. [gRPC and Protocol Buffers](#grpc-and-protocol-buffers)
7. [Event Sourcing and CQRS](#event-sourcing-and-cqrs)
8. [Network Protocols Deep Dive](#network-protocols-deep-dive)

## Microservices Communication

### Communication Patterns

#### 1. Synchronous Communication

**HTTP/REST API Calls**
```javascript
// Service A calling Service B
class UserService {
    constructor() {
        this.orderServiceUrl = process.env.ORDER_SERVICE_URL;
    }
    
    async getUserWithOrders(userId) {
        try {
            // Get user data
            const user = await User.findById(userId);
            
            // Call order service
            const response = await fetch(`${this.orderServiceUrl}/users/${userId}/orders`, {
                headers: {
                    'Authorization': `Bearer ${this.getServiceToken()}`,
                    'Content-Type': 'application/json'
                },
                timeout: 5000 // 5 second timeout
            });
            
            if (!response.ok) {
                throw new Error(`Order service error: ${response.status}`);
            }
            
            const orders = await response.json();
            
            return {
                ...user.toJSON(),
                orders: orders.data
            };
        } catch (error) {
            console.error('Error fetching user with orders:', error);
            // Return user without orders if order service fails
            return user.toJSON();
        }
    }
    
    getServiceToken() {
        // Generate service-to-service JWT token
        return jwt.sign(
            { service: 'user-service', iat: Date.now() },
            process.env.SERVICE_SECRET,
            { expiresIn: '5m' }
        );
    }
}
```

**Service-to-Service Authentication**
```javascript
// Service authentication middleware
const authenticateService = (req, res, next) => {
    const token = req.headers.authorization?.split(' ')[1];
    
    if (!token) {
        return res.status(401).json({ error: 'Service token required' });
    }
    
    try {
        const decoded = jwt.verify(token, process.env.SERVICE_SECRET);
        
        if (!decoded.service) {
            return res.status(401).json({ error: 'Invalid service token' });
        }
        
        req.service = decoded;
        next();
    } catch (error) {
        return res.status(401).json({ error: 'Invalid service token' });
    }
};

// Protected service endpoint
app.get('/internal/users/:id', authenticateService, async (req, res) => {
    const user = await User.findById(req.params.id);
    res.json(user);
});
```

#### 2. Asynchronous Communication

**Event-Driven Communication**
```javascript
// Event publisher
class EventPublisher {
    constructor(messageQueue) {
        this.messageQueue = messageQueue;
    }
    
    async publishEvent(eventType, data, options = {}) {
        const event = {
            id: uuidv4(),
            type: eventType,
            data,
            timestamp: new Date().toISOString(),
            source: process.env.SERVICE_NAME,
            version: '1.0',
            ...options
        };
        
        await this.messageQueue.publish(eventType, event);
        console.log(`Event published: ${eventType}`, event.id);
    }
}

// Event subscriber
class EventSubscriber {
    constructor(messageQueue) {
        this.messageQueue = messageQueue;
        this.handlers = new Map();
    }
    
    subscribe(eventType, handler) {
        this.handlers.set(eventType, handler);
        this.messageQueue.subscribe(eventType, this.handleEvent.bind(this));
    }
    
    async handleEvent(event) {
        const handler = this.handlers.get(event.type);
        
        if (!handler) {
            console.warn(`No handler for event type: ${event.type}`);
            return;
        }
        
        try {
            await handler(event.data, event);
            console.log(`Event processed: ${event.type}`, event.id);
        } catch (error) {
            console.error(`Event processing failed: ${event.type}`, error);
            // Implement retry logic or dead letter queue
            await this.handleFailedEvent(event, error);
        }
    }
    
    async handleFailedEvent(event, error) {
        // Retry logic
        const retryCount = event.retryCount || 0;
        const maxRetries = 3;
        
        if (retryCount < maxRetries) {
            event.retryCount = retryCount + 1;
            event.lastError = error.message;
            
            // Exponential backoff
            const delay = Math.pow(2, retryCount) * 1000;
            
            setTimeout(() => {
                this.messageQueue.publish(event.type, event);
            }, delay);
        } else {
            // Send to dead letter queue
            await this.messageQueue.publish('dead-letter', event);
        }
    }
}
```

### API Gateway Pattern
```javascript
// API Gateway with routing and middleware
class APIGateway {
    constructor() {
        this.app = express();
        this.services = new Map();
        this.setupMiddleware();
        this.setupRoutes();
    }
    
    registerService(name, config) {
        this.services.set(name, {
            url: config.url,
            healthCheck: config.healthCheck,
            timeout: config.timeout || 5000,
            retries: config.retries || 3
        });
    }
    
    setupMiddleware() {
        // Rate limiting
        this.app.use(rateLimit({
            windowMs: 15 * 60 * 1000,
            max: 1000
        }));
        
        // Authentication
        this.app.use('/api', this.authenticateRequest.bind(this));
        
        // Request logging
        this.app.use(this.logRequest.bind(this));
        
        // CORS
        this.app.use(cors());
    }
    
    setupRoutes() {
        // Dynamic routing based on registered services
        this.services.forEach((config, serviceName) => {
            this.app.use(`/api/${serviceName}`, this.createServiceProxy(serviceName));
        });
        
        // Health check aggregation
        this.app.get('/health', this.aggregateHealth.bind(this));
    }
    
    createServiceProxy(serviceName) {
        return httpProxy({
            target: this.services.get(serviceName).url,
            changeOrigin: true,
            pathRewrite: {
                [`^/api/${serviceName}`]: ''
            },
            timeout: this.services.get(serviceName).timeout,
            onError: this.handleProxyError.bind(this),
            onProxyReq: this.modifyProxyRequest.bind(this),
            onProxyRes: this.modifyProxyResponse.bind(this)
        });
    }
    
    async authenticateRequest(req, res, next) {
        const token = req.headers.authorization?.split(' ')[1];
        
        if (!token) {
            return res.status(401).json({ error: 'Authentication required' });
        }
        
        try {
            const user = await this.validateToken(token);
            req.user = user;
            next();
        } catch (error) {
            return res.status(401).json({ error: 'Invalid token' });
        }
    }
    
    logRequest(req, res, next) {
        const start = Date.now();
        
        res.on('finish', () => {
            const duration = Date.now() - start;
            console.log(`${req.method} ${req.path} - ${res.statusCode} - ${duration}ms`);
        });
        
        next();
    }
    
    async aggregateHealth() {
        const healthChecks = await Promise.allSettled(
            Array.from(this.services.entries()).map(async ([name, config]) => {
                try {
                    const response = await fetch(`${config.url}${config.healthCheck}`, {
                        timeout: 2000
                    });
                    return {
                        service: name,
                        status: response.ok ? 'healthy' : 'unhealthy',
                        responseTime: response.headers.get('x-response-time')
                    };
                } catch (error) {
                    return {
                        service: name,
                        status: 'unhealthy',
                        error: error.message
                    };
                }
            })
        );
        
        return healthChecks.map(check => check.value || check.reason);
    }
}

// Usage
const gateway = new APIGateway();

gateway.registerService('users', {
    url: 'http://user-service:3001',
    healthCheck: '/health',
    timeout: 5000
});

gateway.registerService('orders', {
    url: 'http://order-service:3002',
    healthCheck: '/health',
    timeout: 3000
});
```

## Message Queues and Pub/Sub

### Redis Pub/Sub
```javascript
const redis = require('redis');

class RedisPubSub {
    constructor() {
        this.publisher = redis.createClient();
        this.subscriber = redis.createClient();
        this.subscribers = new Map();
    }
    
    async publish(channel, message) {
        const serialized = JSON.stringify({
            data: message,
            timestamp: Date.now(),
            id: uuidv4()
        });
        
        await this.publisher.publish(channel, serialized);
    }
    
    subscribe(channel, handler) {
        if (!this.subscribers.has(channel)) {
            this.subscribers.set(channel, []);
            this.subscriber.subscribe(channel);
        }
        
        this.subscribers.get(channel).push(handler);
    }
    
    init() {
        this.subscriber.on('message', async (channel, message) => {
            const handlers = this.subscribers.get(channel) || [];
            const parsedMessage = JSON.parse(message);
            
            await Promise.all(
                handlers.map(handler => {
                    try {
                        return handler(parsedMessage.data, parsedMessage);
                    } catch (error) {
                        console.error(`Handler error for channel ${channel}:`, error);
                    }
                })
            );
        });
    }
}

// Usage
const pubsub = new RedisPubSub();
pubsub.init();

// Publisher
app.post('/api/orders', async (req, res) => {
    const order = await Order.create(req.body);
    
    // Publish order created event
    await pubsub.publish('order.created', {
        orderId: order.id,
        userId: order.userId,
        total: order.total,
        items: order.items
    });
    
    res.status(201).json(order);
});

// Subscribers
pubsub.subscribe('order.created', async (orderData) => {
    // Send confirmation email
    await emailService.sendOrderConfirmation(orderData);
});

pubsub.subscribe('order.created', async (orderData) => {
    // Update inventory
    await inventoryService.updateStock(orderData.items);
});

pubsub.subscribe('order.created', async (orderData) => {
    // Generate invoice
    await invoiceService.generateInvoice(orderData);
});
```

### RabbitMQ Implementation
```javascript
const amqp = require('amqplib');

class RabbitMQService {
    constructor() {
        this.connection = null;
        this.channel = null;
    }
    
    async connect() {
        try {
            this.connection = await amqp.connect(process.env.RABBITMQ_URL);
            this.channel = await this.connection.createChannel();
            
            // Handle connection errors
            this.connection.on('error', (err) => {
                console.error('RabbitMQ connection error:', err);
            });
            
            this.connection.on('close', () => {
                console.log('RabbitMQ connection closed');
                // Implement reconnection logic
                setTimeout(() => this.connect(), 5000);
            });
            
            console.log('Connected to RabbitMQ');
        } catch (error) {
            console.error('Failed to connect to RabbitMQ:', error);
            setTimeout(() => this.connect(), 5000);
        }
    }
    
    async createExchange(exchangeName, type = 'topic') {
        await this.channel.assertExchange(exchangeName, type, { durable: true });
    }
    
    async createQueue(queueName, options = {}) {
        const defaultOptions = {
            durable: true,
            exclusive: false,
            autoDelete: false
        };
        
        await this.channel.assertQueue(queueName, { ...defaultOptions, ...options });
    }
    
    async publish(exchangeName, routingKey, message, options = {}) {
        const messageBuffer = Buffer.from(JSON.stringify({
            data: message,
            timestamp: Date.now(),
            id: uuidv4()
        }));
        
        const defaultOptions = {
            persistent: true,
            messageId: uuidv4(),
            timestamp: Date.now()
        };
        
        return this.channel.publish(
            exchangeName,
            routingKey,
            messageBuffer,
            { ...defaultOptions, ...options }
        );
    }
    
    async subscribe(queueName, handler, options = {}) {
        const defaultOptions = {
            noAck: false
        };
        
        await this.channel.consume(queueName, async (msg) => {
            if (msg) {
                try {
                    const content = JSON.parse(msg.content.toString());
                    await handler(content.data, content, msg);
                    
                    if (!options.noAck) {
                        this.channel.ack(msg);
                    }
                } catch (error) {
                    console.error('Message processing error:', error);
                    
                    // Reject and requeue with limit
                    const retryCount = msg.properties.headers?.retryCount || 0;
                    
                    if (retryCount < 3) {
                        // Requeue with retry count
                        const headers = {
                            ...msg.properties.headers,
                            retryCount: retryCount + 1
                        };
                        
                        await this.publish(
                            '',
                            queueName,
                            JSON.parse(msg.content.toString()).data,
                            { headers }
                        );
                    } else {
                        // Send to dead letter queue
                        await this.publish(
                            'dead-letter',
                            'failed',
                            {
                                originalMessage: msg.content.toString(),
                                error: error.message,
                                queue: queueName
                            }
                        );
                    }
                    
                    this.channel.nack(msg, false, false);
                }
            }
        }, { ...defaultOptions, ...options });
    }
    
    async bindQueue(queueName, exchangeName, routingKey) {
        await this.channel.bindQueue(queueName, exchangeName, routingKey);
    }
}

// Usage
const rabbitmq = new RabbitMQService();

// Initialize
const initializeMessageQueue = async () => {
    await rabbitmq.connect();
    
    // Create exchanges
    await rabbitmq.createExchange('orders', 'topic');
    await rabbitmq.createExchange('users', 'topic');
    await rabbitmq.createExchange('dead-letter', 'direct');
    
    // Create queues
    await rabbitmq.createQueue('order.created');
    await rabbitmq.createQueue('order.updated');
    await rabbitmq.createQueue('user.created');
    await rabbitmq.createQueue('email.notifications');
    await rabbitmq.createQueue('dead-letter-queue');
    
    // Bind queues to exchanges
    await rabbitmq.bindQueue('order.created', 'orders', 'order.created');
    await rabbitmq.bindQueue('order.updated', 'orders', 'order.updated');
    await rabbitmq.bindQueue('user.created', 'users', 'user.created');
    await rabbitmq.bindQueue('dead-letter-queue', 'dead-letter', 'failed');
};

// Publishers
const publishOrderCreated = async (orderData) => {
    await rabbitmq.publish('orders', 'order.created', orderData);
};

// Subscribers
const subscribeToOrderEvents = async () => {
    await rabbitmq.subscribe('order.created', async (orderData) => {
        console.log('Processing order created:', orderData);
        // Process order
    });
    
    await rabbitmq.subscribe('email.notifications', async (emailData) => {
        console.log('Sending email:', emailData);
        // Send email
    });
};
```

## Service Discovery

### Consul Service Discovery
```javascript
const consul = require('consul')();

class ServiceRegistry {
    constructor(serviceName, port, health) {
        this.serviceName = serviceName;
        this.serviceId = `${serviceName}-${process.pid}`;
        this.port = port;
        this.health = health;
    }
    
    async register() {
        const serviceDefinition = {
            id: this.serviceId,
            name: this.serviceName,
            port: this.port,
            address: process.env.SERVICE_HOST || 'localhost',
            tags: [
                `version-${process.env.SERVICE_VERSION || '1.0.0'}`,
                `environment-${process.env.NODE_ENV || 'development'}`
            ],
            check: {
                http: `http://${process.env.SERVICE_HOST || 'localhost'}:${this.port}${this.health}`,
                interval: '10s',
                timeout: '5s',
                deregistercriticalserviceafter: '30s'
            }
        };
        
        try {
            await consul.agent.service.register(serviceDefinition);
            console.log(`Service ${this.serviceName} registered with Consul`);
            
            // Graceful shutdown
            process.on('SIGINT', () => this.deregister());
            process.on('SIGTERM', () => this.deregister());
        } catch (error) {
            console.error('Failed to register service:', error);
        }
    }
    
    async deregister() {
        try {
            await consul.agent.service.deregister(this.serviceId);
            console.log(`Service ${this.serviceName} deregistered from Consul`);
            process.exit(0);
        } catch (error) {
            console.error('Failed to deregister service:', error);
            process.exit(1);
        }
    }
    
    async discoverService(serviceName) {
        try {
            const services = await consul.health.service({
                service: serviceName,
                passing: true // Only healthy services
            });
            
            return services.map(service => ({
                id: service.Service.ID,
                address: service.Service.Address,
                port: service.Service.Port,
                tags: service.Service.Tags
            }));
        } catch (error) {
            console.error(`Failed to discover service ${serviceName}:`, error);
            return [];
        }
    }
    
    async getServiceUrl(serviceName) {
        const services = await this.discoverService(serviceName);
        
        if (services.length === 0) {
            throw new Error(`No healthy instances of ${serviceName} found`);
        }
        
        // Simple round-robin selection
        const service = services[Math.floor(Math.random() * services.length)];
        return `http://${service.address}:${service.port}`;
    }
}

// Usage
const registry = new ServiceRegistry('user-service', 3001, '/health');

// Register service on startup
registry.register();

// Discover and call other services
const callOrderService = async (userId) => {
    try {
        const orderServiceUrl = await registry.getServiceUrl('order-service');
        const response = await fetch(`${orderServiceUrl}/users/${userId}/orders`);
        return await response.json();
    } catch (error) {
        console.error('Failed to call order service:', error);
        throw error;
    }
};
```

### DNS-based Service Discovery
```javascript
const dns = require('dns').promises;

class DNSServiceDiscovery {
    constructor() {
        this.cache = new Map();
        this.cacheTTL = 60000; // 1 minute
    }
    
    async resolveService(serviceName) {
        const cacheKey = serviceName;
        const cached = this.cache.get(cacheKey);
        
        if (cached && Date.now() - cached.timestamp < this.cacheTTL) {
            return cached.addresses;
        }
        
        try {
            // DNS SRV record lookup
            const srvRecords = await dns.resolveSrv(`_http._tcp.${serviceName}.service.consul`);
            
            const addresses = await Promise.all(
                srvRecords.map(async (srv) => {
                    try {
                        const addresses = await dns.resolve4(srv.name);
                        return addresses.map(address => ({
                            host: address,
                            port: srv.port,
                            priority: srv.priority,
                            weight: srv.weight
                        }));
                    } catch (error) {
                        console.error(`Failed to resolve ${srv.name}:`, error);
                        return [];
                    }
                })
            );
            
            const flatAddresses = addresses.flat();
            
            // Cache the result
            this.cache.set(cacheKey, {
                addresses: flatAddresses,
                timestamp: Date.now()
            });
            
            return flatAddresses;
        } catch (error) {
            console.error(`DNS service discovery failed for ${serviceName}:`, error);
            return [];
        }
    }
    
    selectAddress(addresses) {
        if (addresses.length === 0) {
            throw new Error('No addresses available');
        }
        
        // Weighted random selection
        const totalWeight = addresses.reduce((sum, addr) => sum + addr.weight, 0);
        let random = Math.random() * totalWeight;
        
        for (const address of addresses) {
            random -= address.weight;
            if (random <= 0) {
                return address;
            }
        }
        
        return addresses[0]; // Fallback
    }
    
    async getServiceEndpoint(serviceName) {
        const addresses = await this.resolveService(serviceName);
        const selected = this.selectAddress(addresses);
        return `http://${selected.host}:${selected.port}`;
    }
}
```

## Circuit Breakers and Resilience

### Circuit Breaker Pattern
```javascript
class CircuitBreaker {
    constructor(options = {}) {
        this.failureThreshold = options.failureThreshold || 5;
        this.recoveryTimeout = options.recoveryTimeout || 60000; // 1 minute
        this.monitoringPeriod = options.monitoringPeriod || 10000; // 10 seconds
        
        this.state = 'CLOSED'; // CLOSED, OPEN, HALF_OPEN
        this.failureCount = 0;
        this.lastFailureTime = null;
        this.successCount = 0;
        
        this.metrics = {
            totalRequests: 0,
            successfulRequests: 0,
            failedRequests: 0,
            timeouts: 0
        };
    }
    
    async execute(operation, fallback) {
        this.metrics.totalRequests++;
        
        if (this.state === 'OPEN') {
            if (Date.now() - this.lastFailureTime >= this.recoveryTimeout) {
                this.state = 'HALF_OPEN';
                this.successCount = 0;
                console.log('Circuit breaker transitioning to HALF_OPEN');
            } else {
                console.log('Circuit breaker is OPEN, executing fallback');
                return fallback ? await fallback() : this.getDefaultFallback();
            }
        }
        
        try {
            const result = await this.executeWithTimeout(operation);
            this.onSuccess();
            return result;
        } catch (error) {
            this.onFailure(error);
            
            if (fallback) {
                return await fallback();
            }
            
            throw error;
        }
    }
    
    async executeWithTimeout(operation, timeout = 5000) {
        return new Promise((resolve, reject) => {
            const timer = setTimeout(() => {
                this.metrics.timeouts++;
                reject(new Error('Operation timeout'));
            }, timeout);
            
            operation()
                .then(result => {
                    clearTimeout(timer);
                    resolve(result);
                })
                .catch(error => {
                    clearTimeout(timer);
                    reject(error);
                });
        });
    }
    
    onSuccess() {
        this.metrics.successfulRequests++;
        this.failureCount = 0;
        
        if (this.state === 'HALF_OPEN') {
            this.successCount++;
            if (this.successCount >= 3) { // Require 3 successes to close
                this.state = 'CLOSED';
                console.log('Circuit breaker closed after successful recovery');
            }
        }
    }
    
    onFailure(error) {
        this.metrics.failedRequests++;
        this.failureCount++;
        this.lastFailureTime = Date.now();
        
        if (this.state === 'HALF_OPEN' || this.failureCount >= this.failureThreshold) {
            this.state = 'OPEN';
            console.log(`Circuit breaker opened due to ${this.failureCount} failures`);
        }
    }
    
    getDefaultFallback() {
        return {
            error: 'Service temporarily unavailable',
            fallback: true,
            timestamp: new Date().toISOString()
        };
    }
    
    getMetrics() {
        const successRate = this.metrics.totalRequests > 0 
            ? (this.metrics.successfulRequests / this.metrics.totalRequests) * 100 
            : 0;
            
        return {
            ...this.metrics,
            successRate: successRate.toFixed(2) + '%',
            state: this.state,
            failureCount: this.failureCount
        };
    }
}

// Usage
const orderServiceBreaker = new CircuitBreaker({
    failureThreshold: 3,
    recoveryTimeout: 30000,
    monitoringPeriod: 5000
});

const callOrderService = async (userId) => {
    return await orderServiceBreaker.execute(
        // Primary operation
        async () => {
            const response = await fetch(`${orderServiceUrl}/users/${userId}/orders`, {
                timeout: 5000
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            
            return await response.json();
        },
        // Fallback operation
        async () => {
            console.log('Using cached order data');
            return await getCachedOrders(userId);
        }
    );
};
```
    try {
        const orderServiceUrl = await registry.getServiceUrl('order-service');
        const response = await fetch(`${orderServiceUrl}/users/${userId}/orders`);
        return await response.json();
    } catch (error) {
        console.error('Failed to call order service:', error);
        throw error;
    }
};
```

### DNS-based Service Discovery
```javascript
const dns = require('dns').promises;

class DNSServiceDiscovery {
    constructor() {
        this.cache = new Map();
        this.cacheTTL = 60000; // 1 minute
    }
    
    async resolveService(serviceName) {
        const cacheKey = serviceName;
        const cached = this.cache.get(cacheKey);
        
        if (cached && Date.now() - cached.timestamp < this.cacheTTL) {
            return cached.addresses;
        }
        
        try {
            // DNS SRV record lookup
            const srvRecords = await dns.resolveSrv(`_http._tcp.${serviceName}.service.consul`);
            
            const addresses = await Promise.all(
                srvRecords.map(async (srv) => {
                    try {
                        const addresses = await dns.resolve4(srv.name);
                        return addresses.map(address => ({
                            host: address,
                            port: srv.port,
                            priority: srv.priority,
                            weight: srv.weight
                        }));
                    } catch (error) {
                        console.error(`Failed to resolve ${srv.name}:`, error);
                        return [];
                    }
                })
            );
            
            const flatAddresses = addresses.flat();
            
            // Cache the result
            this.cache.set(cacheKey, {
                addresses: flatAddresses,
                timestamp: Date.now()
            });
            
            return flatAddresses;
        } catch (error) {
            console.error(`DNS service discovery failed for ${serviceName}:`, error);
            return [];
        }
    }
    
    selectAddress(addresses) {
        if (addresses.length === 0) {
            throw new Error('No addresses available');
        }
        
        // Weighted random selection
        const totalWeight = addresses.reduce((sum, addr) => sum + addr.weight, 0);
        let random = Math.random() * totalWeight;
        
        for (const address of addresses) {
            random -= address.weight;
            if (random <= 0) {
                return address;
            }
        }
        
        return addresses[0]; // Fallback
    }
    
    async getServiceEndpoint(serviceName) {
        const addresses = await this.resolveService(serviceName);
        const selected = this.selectAddress(addresses);
        return `http://${selected.host}:${selected.port}`;
    }
}
```

## Circuit Breakers and Resilience

### Circuit Breaker Pattern
```javascript
class CircuitBreaker {
    constructor(options = {}) {
        this.failureThreshold = options.failureThreshold || 5;
        this.recoveryTimeout = options.recoveryTimeout || 60000; // 1 minute
        this.monitoringPeriod = options.monitoringPeriod || 10000; // 10 seconds
        
        this.state = 'CLOSED'; // CLOSED, OPEN, HALF_OPEN
        this.failureCount = 0;
        this.lastFailureTime = null;
        this.successCount = 0;
        
        this.metrics = {
            totalRequests: 0,
            successfulRequests: 0,
            failedRequests: 0,
            timeouts: 0
        };
    }
    
    async execute(operation, fallback) {
        this.metrics.totalRequests++;
        
        if (this.state === 'OPEN') {
            if (Date.now() - this.lastFailureTime >= this.recoveryTimeout) {
                this.state = 'HALF_OPEN';
                this.successCount = 0;
                console.log('Circuit breaker transitioning to HALF_OPEN');
            } else {
                console.log('Circuit breaker is OPEN, executing fallback');
                return fallback ? await fallback() : this.getDefaultFallback();
            }
        }
        
        try {
            const result = await this.executeWithTimeout(operation);
            this.onSuccess();
            return result;
        } catch (error) {
            this.onFailure(error);
            
            if (fallback) {
                return await fallback();
            }
            
            throw error;
        }
    }
    
    async executeWithTimeout(operation, timeout = 5000) {
        return new Promise((resolve, reject) => {
            const timer = setTimeout(() => {
                this.metrics.timeouts++;
                reject(new Error('Operation timeout'));
            }, timeout);
            
            operation()
                .then(result => {
                    clearTimeout(timer);
                    resolve(result);
                })
                .catch(error => {
                    clearTimeout(timer);
                    reject(error);
                });
        });
    }
    
    onSuccess() {
        this.metrics.successfulRequests++;
        this.failureCount = 0;
        
        if (this.state === 'HALF_OPEN') {
            this.successCount++;
            if (this.successCount >= 3) { // Require 3 successes to close
                this.state = 'CLOSED';
                console.log('Circuit breaker closed after successful recovery');
            }
        }
    }
    
    onFailure(error) {
        this.metrics.failedRequests++;
        this.failureCount++;
        this.lastFailureTime = Date.now();
        
        if (this.state === 'HALF_OPEN' || this.failureCount >= this.failureThreshold) {
            this.state = 'OPEN';
            console.log(`Circuit breaker opened due to ${this.failureCount} failures`);
        }
    }
    
    getDefaultFallback() {
        return {
            error: 'Service temporarily unavailable',
            fallback: true,
            timestamp: new Date().toISOString()
        };
    }
    
    getMetrics() {
        const successRate = this.metrics.totalRequests > 0 
            ? (this.metrics.successfulRequests / this.metrics.totalRequests) * 100 
            : 0;
            
        return {
            ...this.metrics,
            successRate: successRate.toFixed(2) + '%',
            state: this.state,
            failureCount: this.failureCount
        };
    }
}

// Usage
const orderServiceBreaker = new CircuitBreaker({
    failureThreshold: 3,
    recoveryTimeout: 30000,
    monitoringPeriod: 5000
});

const callOrderService = async (userId) => {
    return await orderServiceBreaker.execute(
        // Primary operation
        async () => {
            const response = await fetch(`${orderServiceUrl}/users/${userId}/orders`, {
                timeout: 5000
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            
            return await response.json();
        },
        // Fallback operation
        async () => {
            console.log('Using cached order data');
            return await getCachedOrders(userId);
        }
    );
};
```

### Retry Pattern with Exponential Backoff
```javascript
class RetryHandler {
    constructor(options = {}) {
        this.maxRetries = options.maxRetries || 3;
        this.baseDelay = options.baseDelay || 1000;
        this.maxDelay = options.maxDelay || 30000;
        this.backoffMultiplier = options.backoffMultiplier || 2;
        this.jitter = options.jitter || true;
    }
    
    async execute(operation, retryCondition) {
        let lastError;
        
        for (let attempt = 0; attempt <= this.maxRetries; attempt++) {
            try {
                return await operation();
            } catch (error) {
                lastError = error;
                
                if (attempt === this.maxRetries) {
                    break; // No more retries
                }
                
                if (retryCondition && !retryCondition(error)) {
                    break; // Don't retry for this error
                }
                
                const delay = this.calculateDelay(attempt);
                console.log(`Attempt ${attempt + 1} failed, retrying in ${delay}ms:`, error.message);
                
                await this.sleep(delay);
            }
        }
        
        throw lastError;
    }
    
    calculateDelay(attempt) {
        let delay = this.baseDelay * Math.pow(this.backoffMultiplier, attempt);
        delay = Math.min(delay, this.maxDelay);
        
        if (this.jitter) {
            // Add random jitter (±25%)
            const jitterRange = delay * 0.25;
            delay += (Math.random() - 0.5) * 2 * jitterRange;
        }
        
        return Math.max(delay, 0);
    }
    
    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Usage
const retryHandler = new RetryHandler({
    maxRetries: 3,
    baseDelay: 1000,
    backoffMultiplier: 2
});

const callExternalAPI = async () => {
    return await retryHandler.execute(
        async () => {
            const response = await fetch('https://api.external.com/data');
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            return await response.json();
        },
        (error) => {
            // Retry condition: retry on network errors and 5xx status codes
            return error.message.includes('fetch') || 
                   error.message.includes('HTTP 5');
        }
    );
};
```

## Real-time Communication

### Server-Sent Events (SSE)
```javascript
// Server-side SSE implementation
class SSEManager {
    constructor() {
        this.clients = new Map();
    }
    
    addClient(clientId, res) {
        // Set SSE headers
        res.writeHead(200, {
            'Content-Type': 'text/event-stream',
            'Cache-Control': 'no-cache',
            'Connection': 'keep-alive',
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': 'Cache-Control'
        });
        
        // Send initial connection message
        this.sendEvent(res, 'connected', { clientId, timestamp: Date.now() });
        
        // Store client connection
        this.clients.set(clientId, {
            response: res,
            lastPing: Date.now()
        });
        
        // Handle client disconnect
        res.on('close', () => {
            this.clients.delete(clientId);
            console.log(`SSE client ${clientId} disconnected`);
        });
        
        // Send periodic heartbeat
        const heartbeat = setInterval(() => {
            if (this.clients.has(clientId)) {
                this.sendEvent(res, 'heartbeat', { timestamp: Date.now() });
            } else {
                clearInterval(heartbeat);
            }
        }, 30000); // Every 30 seconds
    }
    
    sendEvent(res, event, data) {
        const eventData = `event: ${event}\ndata: ${JSON.stringify(data)}\n\n`;
        res.write(eventData);
    }
    
    broadcast(event, data) {
        this.clients.forEach((client, clientId) => {
            try {
                this.sendEvent(client.response, event, data);
            } catch (error) {
                console.error(`Failed to send to client ${clientId}:`, error);
                this.clients.delete(clientId);
            }
        });
    }
    
    sendToClient(clientId, event, data) {
        const client = this.clients.get(clientId);
        if (client) {
            this.sendEvent(client.response, event, data);
        }
    }
    
    getConnectedClients() {
        return Array.from(this.clients.keys());
    }
}

const sseManager = new SSEManager();

// SSE endpoint
app.get('/api/events', (req, res) => {
    const clientId = req.query.clientId || uuidv4();
    sseManager.addClient(clientId, res);
});

// Broadcast events
app.post('/api/broadcast', (req, res) => {
    const { event, data } = req.body;
    sseManager.broadcast(event, data);
    res.json({ message: 'Event broadcasted' });
});
```
```javascript
class RetryHandler {
    constructor(options = {}) {
        this.maxRetries = options.maxRetries || 3;
        this.baseDelay = options.baseDelay || 1000;
        this.maxDelay = options.maxDelay || 30000;
        this.backoffMultiplier = options.backoffMultiplier || 2;
        this.jitter = options.jitter || true;
    }
    
    async execute(operation, retryCondition) {
        let lastError;
        
        for (let attempt = 0; attempt <= this.maxRetries; attempt++) {
            try {
                return await operation();
            } catch (error) {
                lastError = error;
                
                if (attempt === this.maxRetries) {
                    break; // No more retries
                }
                
                if (retryCondition && !retryCondition(error)) {
                    break; // Don't retry for this error
                }
                
                const delay = this.calculateDelay(attempt);
                console.log(`Attempt ${attempt + 1} failed, retrying in ${delay}ms:`, error.message);
                
                await this.sleep(delay);
            }
        }
        
        throw lastError;
    }
    
    calculateDelay(attempt) {
        let delay = this.baseDelay * Math.pow(this.backoffMultiplier, attempt);
        delay = Math.min(delay, this.maxDelay);
        
        if (this.jitter) {
            // Add random jitter (±25%)
            const jitterRange = delay * 0.25;
            delay += (Math.random() - 0.5) * 2 * jitterRange;
        }
        
        return Math.max(delay, 0);
    }
    
    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Usage
const retryHandler = new RetryHandler({
    maxRetries: 3,
    baseDelay: 1000,
    backoffMultiplier: 2
});

const callExternalAPI = async () => {
    return await retryHandler.execute(
        async () => {
            const response = await fetch('https://api.external.com/data');
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            
            return await response.json();
        },
        (error) => {
            // Retry condition: retry on network errors and 5xx status codes
            return error.message.includes('fetch') || 
                   error.message.includes('HTTP 5');
        }
    );
};
```

## Real-time Communication

### Server-Sent Events (SSE)
```javascript
// Server-side SSE implementation
class SSEManager {
    constructor() {
        this.clients = new Map();
    }
    
    addClient(clientId, res) {
        // Set SSE headers
        res.writeHead(200, {
            'Content-Type': 'text/event-stream',
            'Cache-Control': 'no-cache',
            'Connection': 'keep-alive',
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': 'Cache-Control'
        });
        
        // Send initial connection message
        this.sendEvent(res, 'connected', { clientId, timestamp: Date.now() });
        
        // Store client connection
        this.clients.set(clientId, {
            response: res,
            lastPing: Date.now()
        });
        
        // Handle client disconnect
        res.on('close', () => {
            this.clients.delete(clientId);
            console.log(`SSE client ${clientId} disconnected`);
        });
        
        // Send periodic heartbeat
        const heartbeat = setInterval(() => {
            if (this.clients.has(clientId)) {
                this.sendEvent(res, 'heartbeat', { timestamp: Date.now() });
            } else {
                clearInterval(heartbeat);
            }
        }, 30000); // Every 30 seconds
    }
    
    sendEvent(res, event, data) {
        const eventData = `event: ${event}\ndata: ${JSON.stringify(data)}\n\n`;
        res.write(eventData);
    }
    
    broadcast(event, data) {
        this.clients.forEach((client, clientId) => {
            try {
                this.sendEvent(client.response, event, data);
            } catch (error) {
                console.error(`Failed to send to client ${clientId}:`, error);
                this.clients.delete(clientId);
            }
        });
    }
    
    sendToClient(clientId, event, data) {
        const client = this.clients.get(clientId);
        if (client) {
            this.sendEvent(client.response, event, data);
        }
    }
    
    getConnectedClients() {
        return Array.from(this.clients.keys());
    }
}

const sseManager = new SSEManager();

// SSE endpoint
app.get('/api/events', (req, res) => {
    const clientId = req.query.clientId || uuidv4();
    sseManager.addClient(clientId, res);
});

// Broadcast events
app.post('/api/broadcast', (req, res) => {
    const { event, data } = req.body;
    sseManager.broadcast(event, data);
    res.json({ message: 'Event broadcasted' });
});

// Client-side SSE
class SSEClient {
    constructor(url) {
        this.url = url;
        this.eventSource = null;
        this.reconnectInterval = 5000;
        this.maxReconnectAttempts = 10;
        this.reconnectAttempts = 0;
        this.listeners = new Map();
    }
    
    connect() {
        this.eventSource = new EventSource(this.url);
        
        this.eventSource.onopen = () => {
            console.log('SSE connection opened');
            this.reconnectAttempts = 0;
        };
        
        this.eventSource.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleEvent('message', data);
        };
        
        this.eventSource.onerror = (error) => {
            console.error('SSE error:', error);
            this.eventSource.close();
            this.reconnect();
        };
        
        // Handle custom events
        this.listeners.forEach((handler, eventType) => {
            this.eventSource.addEventListener(eventType, (event) => {
                const data = JSON.parse(event.data);
                handler(data);
            });
        });
    }
    
    reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Reconnecting SSE (attempt ${this.reconnectAttempts})...`);
            
            setTimeout(() => {
                this.connect();
            }, this.reconnectInterval * this.reconnectAttempts);
        } else {
            console.error('Max reconnection attempts reached');
        }
    }
    
    addEventListener(eventType, handler) {
        this.listeners.set(eventType, handler);
        
        if (this.eventSource) {
            this.eventSource.addEventListener(eventType, (event) => {
                const data = JSON.parse(event.data);
                handler(data);
            });
        }
    }
    
    handleEvent(eventType, data) {
        const handler = this.listeners.get(eventType);
        if (handler) {
            handler(data);
        }
    }
    
    close() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }
}

// Usage
const sseClient = new SSEClient('/api/events?clientId=user123');

sseClient.addEventListener('notification', (data) => {
    console.log('Received notification:', data);
    showNotification(data.message);
});

sseClient.addEventListener('update', (data) => {
    console.log('Received update:', data);
    updateUI(data);
});

sseClient.connect();
```

## Network Protocols Deep Dive

### TCP vs UDP Implementation

#### TCP Server (Node.js)
```javascript
const net = require('net');

class TCPServer {
    constructor(port) {
        this.port = port;
        this.clients = new Map();
        this.server = net.createServer();
        this.setupEventHandlers();
    }
    
    setupEventHandlers() {
        this.server.on('connection', (socket) => {
            const clientId = `${socket.remoteAddress}:${socket.remotePort}`;
            console.log(`TCP client connected: ${clientId}`);
            
            this.clients.set(clientId, socket);
            
            // Handle incoming data
            socket.on('data', (data) => {
                this.handleMessage(clientId, data);
            });
            
            // Handle client disconnect
            socket.on('close', () => {
                console.log(`TCP client disconnected: ${clientId}`);
                this.clients.delete(clientId);
            });
            
            // Handle errors
            socket.on('error', (error) => {
                console.error(`TCP client error ${clientId}:`, error);
                this.clients.delete(clientId);
            });
            
            // Send welcome message
            this.sendToClient(clientId, {
                type: 'welcome',
                message: 'Connected to TCP server',
                clientId
            });
        });
        
        this.server.on('error', (error) => {
            console.error('TCP server error:', error);
        });
    }
    
    handleMessage(clientId, data) {
        try {
            const message = JSON.parse(data.toString());
            console.log(`Message from ${clientId}:`, message);
            
            // Echo message back to client
            this.sendToClient(clientId, {
                type: 'echo',
                originalMessage: message,
                timestamp: Date.now()
            });
            
            // Broadcast to other clients
            this.broadcast({
                type: 'broadcast',
                from: clientId,
                message: message
            }, clientId);
        } catch (error) {
            console.error('Error parsing message:', error);
        }
    }
    
    sendToClient(clientId, data) {
        const client = this.clients.get(clientId);
        if (client && !client.destroyed) {
            client.write(JSON.stringify(data) + '\n');
        }
    }
    
    broadcast(data, excludeClientId = null) {
        const message = JSON.stringify(data) + '\n';
        
        this.clients.forEach((client, clientId) => {
            if (clientId !== excludeClientId && !client.destroyed) {
                client.write(message);
            }
        });
    }
    
    start() {
        this.server.listen(this.port, () => {
            console.log(`TCP server listening on port ${this.port}`);
        });
    }
    
    stop() {
        this.clients.forEach((client) => {
            client.destroy();
        });
        this.server.close();
    }
}

// Usage
const tcpServer = new TCPServer(8080);
tcpServer.start();
```

#### UDP Server (Node.js)
```javascript
const dgram = require('dgram');

class UDPServer {
    constructor(port) {
        this.port = port;
        this.clients = new Map();
        this.socket = dgram.createSocket('udp4');
        this.setupEventHandlers();
    }
    
    setupEventHandlers() {
        this.socket.on('message', (data, rinfo) => {
            const clientId = `${rinfo.address}:${rinfo.port}`;
            
            // Track client
            this.clients.set(clientId, {
                address: rinfo.address,
                port: rinfo.port,
                lastSeen: Date.now()
            });
            
            this.handleMessage(clientId, data, rinfo);
        });
        
        this.socket.on('error', (error) => {
            console.error('UDP server error:', error);
        });
        
        this.socket.on('listening', () => {
            const address = this.socket.address();
            console.log(`UDP server listening on ${address.address}:${address.port}`);
        });
        
        // Clean up inactive clients
        setInterval(() => {
            const now = Date.now();
            const timeout = 60000; // 1 minute
            
            this.clients.forEach((client, clientId) => {
                if (now - client.lastSeen > timeout) {
                    console.log(`Removing inactive UDP client: ${clientId}`);
                    this.clients.delete(clientId);
                }
            });
        }, 30000); // Check every 30 seconds
    }
    
    handleMessage(clientId, data, rinfo) {
        try {
            const message = JSON.parse(data.toString());
            console.log(`UDP message from ${clientId}:`, message);
            
            // Send acknowledgment
            this.sendToClient(clientId, {
                type: 'ack',
                originalMessage: message,
                timestamp: Date.now()
            });
            
            // Handle different message types
            switch (message.type) {
                case 'ping':
                    this.sendToClient(clientId, {
                        type: 'pong',
                        timestamp: Date.now()
                    });
                    break;
                    
                case 'broadcast':
                    this.broadcast({
                        type: 'message',
                        from: clientId,
                        data: message.data
                    }, clientId);
                    break;
                    
                default:
                    console.log(`Unknown message type: ${message.type}`);
            }
        } catch (error) {
            console.error('Error parsing UDP message:', error);
        }
    }
    
    sendToClient(clientId, data) {
        const client = this.clients.get(clientId);
        if (client) {
            const message = Buffer.from(JSON.stringify(data));
            this.socket.send(message, client.port, client.address, (error) => {
                if (error) {
                    console.error(`Error sending to ${clientId}:`, error);
                }
            });
        }
    }
    
    broadcast(data, excludeClientId = null) {
        const message = Buffer.from(JSON.stringify(data));
        
        this.clients.forEach((client, clientId) => {
            if (clientId !== excludeClientId) {
                this.socket.send(message, client.port, client.address, (error) => {
                    if (error) {
                        console.error(`Error broadcasting to ${clientId}:`, error);
                    }
                });
            }
        });
    }
    
    start() {
        this.socket.bind(this.port);
    }
    
    stop() {
        this.socket.close();
    }
}

// Usage
const udpServer = new UDPServer(8081);
udpServer.start();
```

### HTTP/2 Implementation
```javascript
const http2 = require('http2');
const fs = require('fs');

class HTTP2Server {
    constructor(options = {}) {
        this.port = options.port || 3000;
        this.routes = new Map();
        
        // Create HTTP/2 server with SSL (required for HTTP/2)
        this.server = http2.createSecureServer({
            key: fs.readFileSync(options.keyPath || 'server.key'),
            cert: fs.readFileSync(options.certPath || 'server.crt')
        });
        
        this.setupEventHandlers();
    }
    
    setupEventHandlers() {
        this.server.on('stream', (stream, headers) => {
            const method = headers[':method'];
            const path = headers[':path'];
            
            console.log(`HTTP/2 ${method} ${path}`);
            
            // Handle the request
            this.handleRequest(stream, method, path, headers);
        });
        
        this.server.on('error', (error) => {
            console.error('HTTP/2 server error:', error);
        });
    }
    
    handleRequest(stream, method, path, headers) {
        const routeKey = `${method} ${path}`;
        const handler = this.routes.get(routeKey);
        
        if (handler) {
            try {
                handler(stream, headers);
            } catch (error) {
                this.sendError(stream, 500, 'Internal Server Error');
            }
        } else {
            this.sendError(stream, 404, 'Not Found');
        }
    }
    
    addRoute(method, path, handler) {
        this.routes.set(`${method} ${path}`, handler);
    }
    
    sendJSON(stream, data, statusCode = 200) {
        stream.respond({
            ':status': statusCode,
            'content-type': 'application/json'
        });
        
        stream.end(JSON.stringify(data));
    }
    
    sendError(stream, statusCode, message) {
        stream.respond({
            ':status': statusCode,
            'content-type': 'application/json'
        });
        
        stream.end(JSON.stringify({ error: message }));
    }
    
    // Server Push example
    pushResource(stream, path, contentType, data) {
        stream.pushStream({ ':path': path }, (error, pushStream) => {
            if (error) {
                console.error('Push stream error:', error);
                return;
            }
            
            pushStream.respond({
                ':status': 200,
                'content-type': contentType
            });
            
            pushStream.end(data);
        });
    }
    
    start() {
        this.server.listen(this.port, () => {
            console.log(`HTTP/2 server listening on port ${this.port}`);
        });
    }
}

// Usage
const http2Server = new HTTP2Server({
    port: 3000,
    keyPath: './certs/server.key',
    certPath: './certs/server.crt'
});

// Add routes
http2Server.addRoute('GET', '/', (stream, headers) => {
    // Push CSS and JS resources
    http2Server.pushResource(stream, '/styles.css', 'text/css', 
        fs.readFileSync('./public/styles.css'));
    http2Server.pushResource(stream, '/app.js', 'application/javascript', 
        fs.readFileSync('./public/app.js'));
    
    // Send HTML response
    stream.respond({
        ':status': 200,
        'content-type': 'text/html'
    });
    
    stream.end(fs.readFileSync('./public/index.html'));
});

http2Server.addRoute('GET', '/api/users', (stream, headers) => {
    // Simulate async operation
    setTimeout(() => {
        http2Server.sendJSON(stream, {
            users: [
                { id: 1, name: 'John Doe' },
                { id: 2, name: 'Jane Smith' }
            ]
        });
    }, 100);
});

http2Server.start();
```

## Summary

This advanced networking guide covers:

1. **Microservices Communication**: Synchronous and asynchronous patterns, API gateways
2. **Message Queues**: Redis Pub/Sub and RabbitMQ implementations
3. **Service Discovery**: Consul and DNS-based discovery
4. **Resilience Patterns**: Circuit breakers and retry mechanisms
5. **Real-time Communication**: Server-Sent Events and WebSocket patterns
6. **gRPC**: Protocol Buffers and high-performance RPC
7. **Event Sourcing & CQRS**: Advanced architectural patterns
8. **Network Protocols**: TCP, UDP, and HTTP/2 implementations

### Key Takeaways

- **Choose the right communication pattern** based on your requirements (sync vs async)
- **Implement resilience patterns** to handle failures gracefully
- **Use appropriate protocols** for different use cases (TCP for reliability, UDP for speed)
- **Consider real-time requirements** when choosing between WebSockets, SSE, or polling
- **Plan for scalability** with proper service discovery and load balancing
- **Monitor and measure** network performance and reliability

### Next Steps

1. Practice implementing these patterns in your applications
2. Learn about container orchestration (Kubernetes networking)
3. Explore service mesh technologies (Istio, Linkerd)
4. Study distributed systems patterns and CAP theorem
5. Implement monitoring and observability for your network layer
class SSEClient {
    constructor(url) {
        this.url = url;
        this.eventSource = null;
        this.reconnectInterval = 5000;
        this.maxReconnectAttempts = 10;
        this.reconnectAttempts = 0;
        this.listeners = new Map();
    }
    
    connect() {
        this.eventSource = new EventSource(this.url);
        
        this.eventSource.onopen = () => {
            console.log('SSE connection opened');
            this.reconnectAttempts = 0;
        };
        
        this.eventSource.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleEvent('message', data);
        };
        
        this.eventSource.onerror = (error) => {
            console.error('SSE error:', error);
            this.eventSource.close();
            this.reconnect();
        };
        
        // Handle custom events
        this.listeners.forEach((handler, eventType) => {
            this.eventSource.addEventListener(eventType, (event) => {
                const data = JSON.parse(event.data);
                handler(data);
            });
        });
    }
    
    reconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            console.log(`Reconnecting SSE (attempt ${this.reconnectAttempts})...`);
            
            setTimeout(() => {
                this.connect();
            }, this.reconnectInterval * this.reconnectAttempts);
        } else {
            console.error('Max reconnection attempts reached');
        }
    }
    
    addEventListener(eventType, handler) {
        this.listeners.set(eventType, handler);
        
        if (this.eventSource) {
            this.eventSource.addEventListener(eventType, (event) => {
                const data = JSON.parse(event.data);
                handler(data);
            });
        }
    }
    
    handleEvent(eventType, data) {
        const handler = this.listeners.get(eventType);
        if (handler) {
            handler(data);
        }
    }
    
    close() {
        if (this.eventSource) {
            this.eventSource.close();
        }
    }
}

// Usage
const sseClient = new SSEClient('/api/events?clientId=user123');

sseClient.addEventListener('notification', (data) => {
    console.log('Received notification:', data);
    showNotification(data.message);
});

sseClient.addEventListener('update', (data) => {
    console.log('Received update:', data);
    updateUI(data);
});

sseClient.connect();
```

### WebSocket Advanced Patterns
```javascript
// WebSocket with rooms and authentication
const WebSocket = require('ws');
const jwt = require('jsonwebtoken');

class WebSocketManager {
    constructor() {
        this.wss = new WebSocket.Server({ port: 8080 });
        this.clients = new Map();
        this.rooms = new Map();
        this.init();
    }
    
    init() {
        this.wss.on('connection', (ws, req) => {
            this.handleConnection(ws, req);
        });
    }
    
    async handleConnection(ws, req) {
        try {
            // Extract token from query string or headers
            const token = new URL(req.url, 'http://localhost').searchParams.get('token');
            
            if (!token) {
                ws.close(1008, 'Authentication required');
                return;
            }
            
            const user = jwt.verify(token, process.env.JWT_SECRET);
            const clientId = uuidv4();
            
            // Store client info
            this.clients.set(clientId, {
                ws,
                user,
                rooms: new Set(),
                lastPing: Date.now()
            });
            
            // Send welcome message
            this.sendToClient(clientId, {
                type: 'welcome',
                clientId,
                user: { id: user.id, name: user.name }
            });
            
            // Handle messages
            ws.on('message', (data) => {
                this.handleMessage(clientId, data);
            });
            
            // Handle disconnect
            ws.on('close', () => {
                this.handleDisconnect(clientId);
            });
            
            // Handle ping/pong for connection health
            ws.on('pong', () => {
                const client = this.clients.get(clientId);
                if (client) {
                    client.lastPing = Date.now();
                }
            });
            
            console.log(`Client ${clientId} connected (user: ${user.name})`);
        } catch (error) {
            console.error('WebSocket authentication failed:', error);
            ws.close(1008, 'Invalid token');
        }
    }
    
    handleMessage(clientId, data) {
        try {
            const message = JSON.parse(data);
            const client = this.clients.get(clientId);
            
            if (!client) return;
            
            switch (message.type) {
                case 'join_room':
                    this.joinRoom(clientId, message.room);
                    break;
                    
                case 'leave_room':
                    this.leaveRoom(clientId, message.room);
                    break;
                    
                case 'room_message':
                    this.sendToRoom(message.room, {
                        type: 'room_message',
                        room: message.room,
                        message: message.message,
                        user: client.user,
                        timestamp: Date.now()
                    }, clientId);
                    break;
                    
                case 'private_message':
                    this.sendToUser(message.targetUserId, {
                        type: 'private_message',
                        message: message.message,
                        from: client.user,
                        timestamp: Date.now()
                    });
                    break;
                    
                default:
                    console.warn(`Unknown message type: ${message.type}`);
            }
        } catch (error) {
            console.error('Error handling message:', error);
        }
    }
    
    joinRoom(clientId, roomName) {
        const client = this.clients.get(clientId);
        if (!client) return;
        
        // Add client to room
        if (!this.rooms.has(roomName)) {
            this.rooms.set(roomName, new Set());
        }
        
        this.rooms.get(roomName).add(clientId);
        client.rooms.add(roomName);
        
        // Notify room members
        this.sendToRoom(roomName, {
            type: 'user_joined',
            room: roomName,
            user: client.user,
            timestamp: Date.now()
        });
        
        console.log(`Client ${clientId} joined room ${roomName}`);
    }
    
    leaveRoom(clientId, roomName) {
        const client = this.clients.get(clientId);
        if (!client) return;
        
        // Remove client from room
        if (this.rooms.has(roomName)) {
            this.rooms.get(roomName).delete(clientId);
            
            // Clean up empty rooms
            if (this.rooms.get(roomName).size === 0) {
                this.rooms.delete(roomName);
            }
        }
        
        client.rooms.delete(roomName);
        
        // Notify room members
        this.sendToRoom(roomName, {
            type: 'user_left',
            room: roomName,
            user: client.user,
            timestamp: Date.now()
        });
        
        console.log(`Client ${clientId} left room ${roomName}`);
    }
    
    sendToClient(clientId, message) {
        const client = this.clients.get(clientId);
        if (client && client.ws.readyState === WebSocket.OPEN) {
            client.ws.send(JSON.stringify(message));
        }
    }
    
    sendToRoom(roomName, message, excludeClientId = null) {
        const room = this.rooms.get(roomName);
        if (!room) return;
        
        room.forEach(clientId => {
            if (clientId !== excludeClientId) {
                this.sendToClient(clientId, message);
            }
        });
    }
    
    sendToUser(userId, message) {
        // Find client by user ID
        for (const [clientId, client] of this.clients) {
            if (client.user.id === userId) {
                this.sendToClient(clientId, message);
                break;
            }
        }
    }
    
    broadcast(message) {
        this.clients.forEach((client, clientId) => {
            this.sendToClient(clientId, message);
        });
    }
    
    handleDisconnect(clientId) {
        const client = this.clients.get(clientId);
        if (!client) return;
        
        // Remove from all rooms
        client.rooms.forEach(roomName => {
            this.leaveRoom(clientId, roomName);
        });
        
        // Remove client
        this.clients.delete(clientId);
        
        console.log(`Client ${clientId} disconnected`);
    }
    
    // Health check - ping clients periodically
    startHealthCheck() {
        setInterval(() => {
            const now = Date.now();
            
            this.clients.forEach((client, clientId) => {
                if (now - client.lastPing > 60000) { // 1 minute timeout
                    console.log(`Client ${clientId} timed out`);
                    client.ws.terminate();
                    this.handleDisconnect(clientId);
                } else {
                    // Send ping
                    client.ws.ping();
                }
            });
        }, 30000); // Check every 30 seconds
    }
    
    getStats() {
        return {
            connectedClients: this.clients.size,
            activeRooms: this.rooms.size,
            roomDetails: Array.from(this.rooms.entries()).map(([name, clients]) => ({
                name,
                clientCount: clients.size
            }))
        };
    }
}

const wsManager = new WebSocketManager();
wsManager.startHealthCheck();
```

## gRPC and Protocol Buffers

### Protocol Buffer Definition
```protobuf
// user.proto
syntax = "proto3";

package user;

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (User);
    rpc UpdateUser(UpdateUserRequest) returns (User);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
    rpc StreamUsers(StreamUsersRequest) returns (stream User);
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
    int64 created_at = 4;
    int64 updated_at = 5;
    repeated string roles = 6;
    UserStatus status = 7;
}

enum UserStatus {
    UNKNOWN = 0;
    ACTIVE = 1;
    INACTIVE = 2;
    SUSPENDED = 3;
}

message GetUserRequest {
    string id = 1;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
    repeated string roles = 3;
}

message UpdateUserRequest {
    string id = 1;
    string name = 2;
    string email = 3;
    repeated string roles = 4;
    UserStatus status = 5;
}

message DeleteUserRequest {
    string id = 1;
}

message DeleteUserResponse {
    bool success = 1;
    string message = 2;
}

message ListUsersRequest {
    int32 page = 1;
    int32 limit = 2;
    string filter = 3;
}

message ListUsersResponse {
    repeated User users = 1;
    int32 total = 2;
    int32 page = 3;
    int32 limit = 4;
}

message StreamUsersRequest {
    string filter = 1;
}
```

### gRPC Server Implementation
```javascript
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

// Load proto file
const PROTO_PATH = path.join(__dirname, 'protos/user.proto');
const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
});

const userProto = grpc.loadPackageDefinition(packageDefinition).user;

class UserServiceImpl {
    async getUser(call, callback) {
        try {
            const { id } = call.request;
            const user = await User.findById(id);
            
            if (!user) {
                return callback({
                    code: grpc.status.NOT_FOUND,
                    message: 'User not found'
                });
            }
            
            callback(null, {
                id: user.id,
                name: user.name,
                email: user.email,
                created_at: user.createdAt.getTime(),
                updated_at: user.updatedAt.getTime(),
                roles: user.roles,
                status: user.status
            });
        } catch (error) {
            callback({
                code: grpc.status.INTERNAL,
                message: error.message
            });
        }
    }
    
    async createUser(call, callback) {
        try {
            const { name, email, roles } = call.request;
            
            const user = await User.create({
                name,
                email,
                roles: roles || [],
                status: 'ACTIVE'
            });
            
            callback(null, {
                id: user.id,
                name: user.name,
                email: user.email,
                created_at: user.createdAt.getTime(),
                updated_at: user.updatedAt.getTime(),
                roles: user.roles,
                status: user.status
            });
        } catch (error) {
            callback({
                code: grpc.status.INVALID_ARGUMENT,
                message: error.message
            });
        }
    }
    
    async listUsers(call, callback) {
        try {
            const { page = 1, limit = 10, filter } = call.request;
            const offset = (page - 1) * limit;
            
            const where = filter ? {
                [Op.or]: [
                    { name: { [Op.iLike]: `%${filter}%` } },
                    { email: { [Op.iLike]: `%${filter}%` } }
                ]
            } : {};
            
            const { rows: users, count: total } = await User.findAndCountAll({
                where,
                limit,
                offset,
                order: [['createdAt', 'DESC']]
            });
            
            callback(null, {
                users: users.map(user => ({
                    id: user.id,
                    name: user.name,
                    email: user.email,
                    created_at: user.createdAt.getTime(),
                    updated_at: user.updatedAt.getTime(),
                    roles: user.roles,
                    status: user.status
                })),
                total,
                page,
                limit
            });
        } catch (error) {
            callback({
                code: grpc.status.INTERNAL,
                message: error.message
            });
        }
    }
    
    async streamUsers(call) {
        try {
            const { filter } = call.request;
            
            const where = filter ? {
                [Op.or]: [
                    { name: { [Op.iLike]: `%${filter}%` } },
                    { email: { [Op.iLike]: `%${filter}%` } }
                ]
            } : {};
            
            const users = await User.findAll({ where });
            
            for (const user of users) {
                call.write({
                    id: user.id,
                    name: user.name,
                    email: user.email,
                    created_at: user.createdAt.getTime(),
                    updated_at: user.updatedAt.getTime(),
                    roles: user.roles,
                    status: user.status
                });
                
                // Simulate streaming delay
                await new Promise(resolve => setTimeout(resolve, 100));
            }
            
            call.end();
        } catch (error) {
            call.emit('error', {
                code: grpc.status.INTERNAL,
                message: error.message
            });
        }
    }
}

// Create and start gRPC server
const server = new grpc.Server();
server.addService(userProto.UserService.service, new UserServiceImpl());

const port = process.env.GRPC_PORT || 50051;
server.bindAsync(
    `0.0.0.0:${port}`,
    grpc.ServerCredentials.createInsecure(),
    (error, port) => {
        if (error) {
            console.error('Failed to start gRPC server:', error);
            return;
        }
        
        console.log(`gRPC server running on port ${port}`);
        server.start();
    }
);
```

### gRPC Client Implementation
```javascript
class UserServiceClient {
    constructor(serverAddress = 'localhost:50051') {
        this.client = new userProto.UserService(
            serverAddress,
            grpc.credentials.createInsecure()
        );
    }
    
    async getUser(id) {
        return new Promise((resolve, reject) => {
            this.client.getUser({ id }, (error, response) => {
                if (error) {
                    reject(error);
                } else {
                    resolve(response);
                }
            });
        });
    }
    
    async createUser(userData) {
        return new Promise((resolve, reject) => {
            this.client.createUser(userData, (error, response) => {
                if (error) {
                    reject(error);
                } else {
                    resolve(response);
                }
            });
        });
    }
    
    async listUsers(options = {}) {
        return new Promise((resolve, reject) => {
            this.client.listUsers(options, (error, response) => {
                if (error) {
                    reject(error);
                } else {
                    resolve(response);
                }
            });
        });
    }
    
    streamUsers(filter = '') {
        const call = this.client.streamUsers({ filter });
        const users = [];
        
        return new Promise((resolve, reject) => {
            call.on('data', (user) => {
                users.push(user);
                console.log('Received user:', user.name);
            });
            
            call.on('end', () => {
                resolve(users);
            });
            
            call.on('error', (error) => {
                reject(error);
            });
        });
    }
    
    close() {
        this.client.close();
    }
}

// Usage
const userClient = new UserServiceClient();

// Get user
const user = await userClient.getUser('123');
console.log('User:', user);

// Create user
const newUser = await userClient.createUser({
    name: 'John Doe',
    email: 'john@example.com',
    roles: ['user']
});

// List users
const userList = await userClient.listUsers({
    page: 1,
    limit: 10,
    filter: 'john'
});

// Stream users
const streamedUsers = await userClient.streamUsers('active');
```

## Event Sourcing and CQRS

### Event Store Implementation
```javascript
class EventStore {
    constructor() {
        this.events = new Map(); // In production, use a database
        this.snapshots = new Map();
    }
    
    async saveEvent(streamId, event, expectedVersion = -1) {
        if (!this.events.has(streamId)) {
            this.events.set(streamId, []);
        }
        
        const stream = this.events.get(streamId);
        
        // Optimistic concurrency check
        if (expectedVersion !== -1 && stream.length !== expectedVersion) {
            throw new Error('Concurrency conflict');
        }
        
        const eventWithMetadata = {
            ...event,
            eventId: uuidv4(),
            streamId,
            version: stream.length + 1,
            timestamp: new Date().toISOString(),
            metadata: event.metadata || {}
        };
        
        stream.push(eventWithMetadata);
        
        // Publish event for projections
        await this.publishEvent(eventWithMetadata);
        
        return eventWithMetadata;
    }
    
    async getEvents(streamId, fromVersion = 1) {
        const stream = this.events.get(streamId) || [];
        return stream.filter(event => event.version >= fromVersion);
    }
    
    async saveSnapshot(streamId, snapshot, version) {
        this.snapshots.set(streamId, {
            ...snapshot,
            version,
            timestamp: new Date().toISOString()
        });
    }
    
    async getSnapshot(streamId) {
        return this.snapshots.get(streamId);
    }
    
    async publishEvent(event) {
        // Publish to event bus for projections and other services
        await eventBus.publish('domain.event', event);
    }
    
    async getAllEvents(fromTimestamp = null) {
        const allEvents = [];
        
        for (const stream of this.events.values()) {
            for (const event of stream) {
                if (!fromTimestamp || event.timestamp >= fromTimestamp) {