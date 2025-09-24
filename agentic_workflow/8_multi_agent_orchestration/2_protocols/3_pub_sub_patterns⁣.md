# Pub/Sub (Publish-Subscribe) Patterns - Complete Guide

Pub/Sub is a messaging pattern where senders (publishers) categorize messages into topics/classes without knowing the subscribers. Subscribers express interest in topics and receive relevant messages.

## Core Concepts

### 1. **Basic Pub/Sub Structure**

```javascript
class PubSub {
    constructor() {
        this.topics = new Map(); // topic -> Set of subscribers
        this.history = new Map(); // topic -> message history
        this.maxHistory = 100; // max messages per topic to retain
    }

    // Subscribe to a topic
    subscribe(topic, subscriber, options = {}) {
        if (!this.topics.has(topic)) {
            this.topics.set(topic, new Set());
            this.history.set(topic, []);
        }

        const subscription = {
            id: this.generateId(),
            subscriber,
            topic,
            options: {
                receiveHistory: options.receiveHistory || false,
                filter: options.filter || null,
                once: options.once || false
            }
        };

        this.topics.get(topic).add(subscription);

        // Send historical messages if requested
        if (options.receiveHistory) {
            this.history.get(topic).forEach(message => {
                this.deliverMessage(subscription, message);
            });
        }

        return subscription.id;
    }

    // Unsubscribe from a topic
    unsubscribe(topic, subscriptionId) {
        if (!this.topics.has(topic)) return false;

        const subscriptions = this.topics.get(topic);
        for (let sub of subscriptions) {
            if (sub.id === subscriptionId) {
                subscriptions.delete(sub);
                return true;
            }
        }
        return false;
    }

    // Publish to a topic
    publish(topic, message, metadata = {}) {
        if (!this.topics.has(topic)) {
            this.topics.set(topic, new Set());
            this.history.set(topic, []);
        }

        const messageObj = {
            id: this.generateId(),
            topic,
            data: message,
            metadata: {
                timestamp: new Date(),
                publisher: metadata.publisher || 'unknown',
                ...metadata
            }
        };

        // Store in history
        const history = this.history.get(topic);
        history.push(messageObj);
        if (history.length > this.maxHistory) {
            history.shift(); // Remove oldest message
        }

        // Deliver to subscribers
        this.deliverToSubscribers(topic, messageObj);

        return messageObj.id;
    }

    // Deliver message to all subscribers of a topic
    deliverToSubscribers(topic, message) {
        if (!this.topics.has(topic)) return;

        const subscriptions = this.topics.get(topic);
        const toRemove = [];

        subscriptions.forEach(subscription => {
            try {
                // Apply filter if exists
                if (subscription.options.filter && 
                    !subscription.options.filter(message)) {
                    return;
                }

                this.deliverMessage(subscription, message);

                // Remove if it's a once subscription
                if (subscription.options.once) {
                    toRemove.push(subscription);
                }
            } catch (error) {
                console.error(`Error delivering message to subscriber:`, error);
            }
        });

        // Remove once subscriptions
        toRemove.forEach(sub => subscriptions.delete(sub));
    }

    deliverMessage(subscription, message) {
        if (typeof subscription.subscriber === 'function') {
            subscription.subscriber(message);
        } else if (subscription.subscriber && 
                   typeof subscription.subscriber.onMessage === 'function') {
            subscription.subscriber.onMessage(message);
        }
    }

    generateId() {
        return Math.random().toString(36).substr(2, 9);
    }

    // Get all topics
    getTopics() {
        return Array.from(this.topics.keys());
    }

    // Get subscribers for a topic
    getSubscribers(topic) {
        return this.topics.has(topic) ? 
            Array.from(this.topics.get(topic)).map(s => s.subscriber) : [];
    }
}
```

## Advanced Pub/Sub Patterns

### 1. **Wildcard and Pattern Matching**

```javascript
class AdvancedPubSub extends PubSub {
    constructor() {
        super();
        this.patternSubscriptions = new Map(); // pattern -> subscriptions
    }

    // Support for wildcard subscriptions (topic.*, topic.#, etc.)
    subscribePattern(pattern, subscriber, options = {}) {
        const subscription = {
            id: this.generateId(),
            subscriber,
            pattern,
            options
        };

        if (!this.patternSubscriptions.has(pattern)) {
            this.patternSubscriptions.set(pattern, new Set());
        }
        this.patternSubscriptions.get(pattern).add(subscription);

        return subscription.id;
    }

    publish(topic, message, metadata = {}) {
        super.publish(topic, message, metadata);

        // Deliver to pattern subscribers
        this.deliverToPatternSubscribers(topic, message, metadata);
    }

    deliverToPatternSubscribers(topic, message, metadata) {
        const messageObj = {
            id: this.generateId(),
            topic,
            data: message,
            metadata: {
                timestamp: new Date(),
                publisher: metadata.publisher || 'unknown',
                ...metadata
            }
        };

        this.patternSubscriptions.forEach((subscriptions, pattern) => {
            if (this.topicMatchesPattern(topic, pattern)) {
                subscriptions.forEach(subscription => {
                    this.deliverMessage(subscription, messageObj);
                });
            }
        });
    }

    topicMatchesPattern(topic, pattern) {
        // Support for:
        // - exact match: "user.created"
        // - single level wildcard: "user.*" matches "user.created" but not "user.created.email"
        // - multi-level wildcard: "user.#" matches "user.created", "user.created.email", etc.
        
        if (pattern === topic) return true;
        if (pattern === '#') return true;
        
        const patternParts = pattern.split('.');
        const topicParts = topic.split('.');
        
        for (let i = 0; i < patternParts.length; i++) {
            if (patternParts[i] === '#') return true; // Multi-level wildcard
            if (patternParts[i] === '*') { // Single-level wildcard
                if (i >= topicParts.length) return false;
                continue;
            }
            if (patternParts[i] !== topicParts[i]) return false;
        }
        
        return patternParts.length === topicParts.length;
    }
}
```

### 2. **Durable Subscriptions with Message Queuing**

```javascript
class DurablePubSub extends AdvancedPubSub {
    constructor() {
        super();
        this.queues = new Map(); // subscriberId -> message queue
        this.offlineSubscribers = new Set();
    }

    subscribe(topic, subscriber, options = {}) {
        const subscriptionId = super.subscribe(topic, subscriber, options);
        
        if (options.durable) {
            this.queues.set(subscriptionId, {
                messages: [],
                maxQueueSize: options.maxQueueSize || 1000,
                subscriber
            });
        }
        
        return subscriptionId;
    }

    deliverMessage(subscription, message) {
        // If subscriber is offline and durable, queue the message
        if (this.offlineSubscribers.has(subscription.id) && this.queues.has(subscription.id)) {
            const queue = this.queues.get(subscription.id);
            if (queue.messages.length < queue.maxQueueSize) {
                queue.messages.push(message);
            }
            return;
        }

        try {
            super.deliverMessage(subscription, message);
        } catch (error) {
            // If delivery fails and it's durable, go offline
            if (this.queues.has(subscription.id)) {
                this.offlineSubscribers.add(subscription.id);
                const queue = this.queues.get(subscription.id);
                queue.messages.push(message);
            }
        }
    }

    // Mark subscriber as online and deliver queued messages
    setOnline(subscriptionId) {
        if (!this.queues.has(subscriptionId)) return;

        this.offlineSubscribers.delete(subscriptionId);
        const queue = this.queues.get(subscriptionId);
        
        // Deliver queued messages
        while (queue.messages.length > 0) {
            const message = queue.messages.shift();
            try {
                if (typeof queue.subscriber === 'function') {
                    queue.subscriber(message);
                }
            } catch (error) {
                // If delivery fails, put back in queue and go offline
                queue.messages.unshift(message);
                this.offlineSubscribers.add(subscriptionId);
                break;
            }
        }
    }

    setOffline(subscriptionId) {
        this.offlineSubscribers.add(subscriptionId);
    }

    getQueueSize(subscriptionId) {
        return this.queues.has(subscriptionId) ? 
            this.queues.get(subscriptionId).messages.length : 0;
    }
}
```

## Real-World Examples

### 1. **Stock Market Data System**

```javascript
class StockMarketSystem {
    constructor() {
        this.pubsub = new DurablePubSub();
        this.stocks = new Map();
        this.setupStockFeeds();
    }

    setupStockFeeds() {
        // Simulate stock price updates
        setInterval(() => {
            this.updateStockPrices();
        }, 1000);
    }

    updateStockPrices() {
        const stocks = ['AAPL', 'GOOGL', 'MSFT', 'AMZN', 'TSLA'];
        
        stocks.forEach(symbol => {
            const currentPrice = this.stocks.get(symbol) || 
                Math.random() * 1000 + 100; // Initial price
            const change = (Math.random() - 0.5) * 10; // Random change
            const newPrice = Math.max(1, currentPrice + change); // Prevent negative prices
            
            this.stocks.set(symbol, newPrice);
            
            this.pubsub.publish(`stocks.${symbol}`, {
                symbol,
                price: newPrice,
                change: newPrice - currentPrice,
                changePercent: ((newPrice - currentPrice) / currentPrice) * 100,
                volume: Math.floor(Math.random() * 1000000),
                timestamp: new Date()
            }, { publisher: 'stock-feed' });
            
            // Also publish to aggregate topics
            this.pubsub.publish('stocks.all', {
                symbol,
                price: newPrice,
                timestamp: new Date()
            }, { publisher: 'stock-feed' });
        });
    }

    // Subscribe to individual stocks
    subscribeToStock(symbol, callback, options = {}) {
        return this.pubsub.subscribe(`stocks.${symbol}`, callback, options);
    }

    // Subscribe to all stocks with filtering
    subscribeToAllStocks(callback, filter = null) {
        return this.pubsub.subscribePattern('stocks.*', callback, {
            filter: filter ? (msg) => filter(msg.data) : null
        });
    }

    // Alert system for price thresholds
    setPriceAlert(symbol, condition, threshold, callback) {
        return this.subscribeToStock(symbol, (message) => {
            const { price } = message.data;
            let trigger = false;
            
            switch (condition) {
                case 'above':
                    trigger = price > threshold;
                    break;
                case 'below':
                    trigger = price < threshold;
                    break;
                case 'crossesAbove':
                    trigger = price > threshold && 
                              this.previousPrices.get(symbol) <= threshold;
                    break;
                case 'crossesBelow':
                    trigger = price < threshold && 
                              this.previousPrices.get(symbol) >= threshold;
                    break;
            }
            
            if (trigger) {
                callback({
                    symbol,
                    price,
                    threshold,
                    condition,
                    timestamp: new Date()
                });
            }
            
            this.previousPrices.set(symbol, price);
        });
    }
}

// Usage example
const stockSystem = new StockMarketSystem();

// Subscribe to Apple stock
stockSystem.subscribeToStock('AAPL', (message) => {
    console.log(`AAPL: $${message.data.price.toFixed(2)} ` +
                `(${message.data.change >= 0 ? '+' : ''}${message.data.change.toFixed(2)})`);
});

// Subscribe to all tech stocks
stockSystem.subscribeToAllStocks((message) => {
    if (['AAPL', 'GOOGL', 'MSFT'].includes(message.data.symbol)) {
        console.log(`Tech Stock ${message.data.symbol}: $${message.data.price.toFixed(2)}`);
    }
});

// Price alert
stockSystem.setPriceAlert('TSLA', 'above', 800, (alert) => {
    console.log(`🚀 ALERT: ${alert.symbol} crossed $${alert.threshold}! Current: $${alert.price}`);
});
```

### 2. **Real-Time Chat Application**

```javascript
class ChatApplication {
    constructor() {
        this.pubsub = new AdvancedPubSub();
        this.rooms = new Map();
        this.users = new Map();
        this.setupDefaultHandlers();
    }

    setupDefaultHandlers() {
        // Handle user presence
        this.pubsub.subscribePattern('user.*.online', (message) => {
            this.broadcastToRoom(message.data.roomId, 'user_joined', {
                userId: message.data.userId,
                username: message.data.username,
                timestamp: new Date()
            });
        });

        this.pubsub.subscribePattern('user.*.offline', (message) => {
            this.broadcastToRoom(message.data.roomId, 'user_left', {
                userId: message.data.userId,
                username: message.data.username,
                timestamp: new Date()
            });
        });

        // Handle message moderation
        this.pubsub.subscribePattern('room.*.message', (message) => {
            if (this.containsProfanity(message.data.content)) {
                this.pubsub.publish(`moderation.flagged`, {
                    messageId: message.id,
                    content: message.data.content,
                    userId: message.data.userId,
                    reason: 'profanity'
                });
            }
        });
    }

    createRoom(roomId, name, options = {}) {
        this.rooms.set(roomId, {
            id: roomId,
            name,
            users: new Set(),
            createdAt: new Date(),
            ...options
        });
        return roomId;
    }

    joinRoom(userId, roomId, username) {
        if (!this.rooms.has(roomId)) {
            throw new Error(`Room ${roomId} does not exist`);
        }

        const room = this.rooms.get(roomId);
        room.users.add(userId);

        this.users.set(userId, {
            id: userId,
            username,
            currentRoom: roomId,
            joinedAt: new Date()
        });

        // Notify about user joining
        this.pubsub.publish(`user.${userId}.online`, {
            userId,
            username,
            roomId
        });

        this.pubsub.publish(`room.${roomId}.user.join`, {
            userId,
            username,
            roomId,
            timestamp: new Date()
        });

        return userId;
    }

    leaveRoom(userId) {
        const user = this.users.get(userId);
        if (!user) return;

        const room = this.rooms.get(user.currentRoom);
        if (room) {
            room.users.delete(userId);
            
            this.pubsub.publish(`user.${userId}.offline`, {
                userId,
                username: user.username,
                roomId: user.currentRoom
            });

            this.pubsub.publish(`room.${user.currentRoom}.user.leave`, {
                userId,
                username: user.username,
                roomId: user.currentRoom,
                timestamp: new Date()
            });
        }

        this.users.delete(userId);
    }

    sendMessage(userId, content, roomId = null) {
        const user = this.users.get(userId);
        if (!user) throw new Error('User not found');

        const targetRoom = roomId || user.currentRoom;
        const room = this.rooms.get(targetRoom);
        if (!room) throw new Error('Room not found');

        const message = {
            id: this.generateMessageId(),
            userId,
            username: user.username,
            content,
            roomId: targetRoom,
            timestamp: new Date(),
            type: 'message'
        };

        this.pubsub.publish(`room.${targetRoom}.message`, message);
        return message.id;
    }

    broadcastToRoom(roomId, eventType, data) {
        this.pubsub.publish(`room.${roomId}.${eventType}`, {
            ...data,
            roomId,
            timestamp: new Date()
        });
    }

    // Subscribe to room events
    subscribeToRoom(roomId, eventTypes, callback) {
        const subscriptions = [];
        
        eventTypes.forEach(eventType => {
            const subId = this.pubsub.subscribe(`room.${roomId}.${eventType}`, callback);
            subscriptions.push(subId);
        });
        
        return subscriptions;
    }

    // Subscribe to user events across all rooms
    subscribeToUserEvents(userId, callback) {
        return this.pubsub.subscribePattern(`user.${userId}.*`, callback);
    }

    containsProfanity(text) {
        const profanity = ['badword1', 'badword2', 'badword3'];
        return profanity.some(word => text.toLowerCase().includes(word));
    }

    generateMessageId() {
        return `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    }
}

// Usage example
const chatApp = new ChatApplication();

// Create rooms
chatApp.createRoom('general', 'General Chat');
chatApp.createRoom('tech', 'Technology Discussion');

// User joins and interacts
const aliceId = chatApp.joinRoom('user1', 'general', 'Alice');
const bobId = chatApp.joinRoom('user2', 'general', 'Bob');

// Subscribe to room messages
chatApp.subscribeToRoom('general', ['message', 'user_joined', 'user_left'], (message) => {
    switch (message.topic) {
        case 'room.general.message':
            console.log(`[${message.data.username}]: ${message.data.content}`);
            break;
        case 'room.general.user_joined':
            console.log(`→ ${message.data.username} joined the room`);
            break;
        case 'room.general.user_left':
            console.log(`← ${message.data.username} left the room`);
            break;
    }
});

// Send messages
chatApp.sendMessage(aliceId, 'Hello everyone!');
chatApp.sendMessage(bobId, 'Hi Alice! How are you?');

// User leaves
setTimeout(() => {
    chatApp.leaveRoom(bobId);
}, 3000);
```

### 3. **Microservices Communication**

```javascript
class MicroservicesOrchestrator {
    constructor() {
        this.pubsub = new DurablePubSub();
        this.services = new Map();
        this.messageCounts = new Map();
        this.setupMonitoring();
    }

    registerService(serviceName, instance) {
        this.services.set(serviceName, {
            instance,
            status: 'running',
            lastHeartbeat: new Date(),
            subscriptions: []
        });

        // Subscribe to service-specific commands
        const commandSub = this.pubsub.subscribe(`command.${serviceName}`, (message) => {
            this.handleServiceCommand(serviceName, message);
        });

        this.services.get(serviceName).subscriptions.push(commandSub);
    }

    // Service-to-service communication
    sendCommand(targetService, command, data, options = {}) {
        const message = {
            id: this.generateId(),
            from: options.from || 'orchestrator',
            command,
            data,
            timestamp: new Date(),
            replyTo: options.replyTo // For request-reply pattern
        };

        return this.pubsub.publish(`command.${targetService}`, message, {
            persistent: options.persistent || false
        });
    }

    // Publish service event
    publishEvent(serviceName, eventType, data) {
        const event = {
            id: this.generateId(),
            service: serviceName,
            type: eventType,
            data,
            timestamp: new Date(),
            version: '1.0'
        };

        this.pubsub.publish(`event.${serviceName}.${eventType}`, event);
        this.pubsub.publish(`event.${serviceName}.*`, event); // Wildcard for all events from service
        this.pubsub.publish(`event.*.${eventType}`, event); // Wildcard for event type across services

        // Update metrics
        this.incrementMessageCount(`event.${eventType}`);
    }

    // Request-reply pattern
    async request(serviceName, command, data, timeout = 5000) {
        return new Promise((resolve, reject) => {
            const replyTopic = `reply.${this.generateId()}`;
            const timeoutId = setTimeout(() => {
                this.pubsub.unsubscribe(replyTopic, replySub);
                reject(new Error(`Request timeout after ${timeout}ms`));
            }, timeout);

            const replySub = this.pubsub.subscribe(replyTopic, (message) => {
                clearTimeout(timeoutId);
                this.pubsub.unsubscribe(replyTopic, replySub);
                resolve(message.data);
            }, { once: true });

            this.sendCommand(serviceName, command, data, {
                replyTo: replyTopic
            });
        });
    }

    handleServiceCommand(serviceName, message) {
        const service = this.services.get(serviceName);
        if (!service || service.status !== 'running') {
            // Service unavailable, maybe persist and retry later
            if (message.metadata.persistent) {
                console.log(`Service ${serviceName} unavailable, message queued`);
            }
            return;
        }

        try {
            if (service.instance.handleCommand) {
                const result = service.instance.handleCommand(message.data.command, message.data.data);
                
                // Send reply if requested
                if (message.data.replyTo) {
                    this.pubsub.publish(message.data.replyTo, {
                        success: true,
                        result,
                        originalMessageId: message.id
                    });
                }
            }
        } catch (error) {
            if (message.data.replyTo) {
                this.pubsub.publish(message.data.replyTo, {
                    success: false,
                    error: error.message,
                    originalMessageId: message.id
                });
            }
        }
    }

    setupMonitoring() {
        // Monitor all messages
        this.pubsub.subscribePattern('*', (message) => {
            this.incrementMessageCount(message.topic);
        });

        // Health checks
        setInterval(() => {
            this.checkServiceHealth();
        }, 30000);
    }

    incrementMessageCount(topic) {
        const count = this.messageCounts.get(topic) || 0;
        this.messageCounts.set(topic, count + 1);
    }

    getMessageStats() {
        return Array.from(this.messageCounts.entries()).map(([topic, count]) => ({
            topic,
            count
        }));
    }

    checkServiceHealth() {
        this.services.forEach((service, name) => {
            const timeSinceHeartbeat = Date.now() - service.lastHeartbeat.getTime();
            if (timeSinceHeartbeat > 60000) { // 1 minute
                service.status = 'unhealthy';
                this.publishEvent('orchestrator', 'service.unhealthy', {
                    service: name,
                    lastHeartbeat: service.lastHeartbeat
                });
            }
        });
    }

    generateId() {
        return Math.random().toString(36).substr(2, 9);
    }
}

// Example service implementation
class OrderService {
    constructor(orchestrator) {
        this.orchestrator = orchestrator;
        this.setupEventHandlers();
    }

    setupEventHandlers() {
        // Listen for payment events
        this.orchestrator.pubsub.subscribePattern('event.payment.*', (message) => {
            if (message.data.type === 'payment.completed') {
                this.handlePaymentCompleted(message.data);
            }
        });
    }

    handleCommand(command, data) {
        switch (command) {
            case 'createOrder':
                return this.createOrder(data);
            case 'cancelOrder':
                return this.cancelOrder(data.orderId);
            default:
                throw new Error(`Unknown command: ${command}`);
        }
    }

    createOrder(orderData) {
        const order = {
            id: `order_${Date.now()}`,
            ...orderData,
            status: 'created',
            createdAt: new Date()
        };

        // Publish order created event
        this.orchestrator.publishEvent('order', 'created', order);

        return order;
    }

    handlePaymentCompleted(paymentData) {
        // Update order status when payment is completed
        this.orchestrator.publishEvent('order', 'payment_received', {
            orderId: paymentData.orderId,
            paymentId: paymentData.paymentId
        });
    }

    cancelOrder(orderId) {
        this.orchestrator.publishEvent('order', 'cancelled', { orderId });
        return { success: true, orderId };
    }
}

// Usage
const orchestrator = new MicroservicesOrchestrator();

// Register services
const orderService = new OrderService(orchestrator);
orchestrator.registerService('order', orderService);

// Use request-reply pattern
async function demo() {
    try {
        const result = await orchestrator.request('order', 'createOrder', {
            userId: 'user123',
            items: [{ productId: 'p1', quantity: 2 }],
            total: 99.99
        });
        
        console.log('Order created:', result);
    } catch (error) {
        console.error('Failed to create order:', error);
    }
}

demo();
```

## Key Benefits of Pub/Sub Pattern

1. **Decoupling**: Publishers and subscribers don't need to know about each other
2. **Scalability**: Easy to add new subscribers without modifying publishers
3. **Flexibility**: Multiple subscription patterns (exact, wildcard, filtered)
4. **Durability**: Message persistence for offline subscribers
5. **Load Distribution**: Natural load balancing across subscribers
6. **Extensibility**: Easy to add new message types and subscribers

## Common Use Cases

- **Real-time notifications**
- **Event-driven architectures**
- **Microservices communication**
- **IoT data distribution**
- **Financial market data**
- **Social media feeds**
- **Log aggregation and monitoring**

Pub/Sub is fundamental to modern distributed systems, providing the backbone for scalable, responsive applications.