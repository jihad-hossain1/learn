# Event-Driven Architecture (EDA) - Complete Guide

Event-Driven Architecture is a software design pattern that focuses on producing, detecting, consuming, and reacting to events. Let me explain this concept in depth with JavaScript examples.

## Core Concepts

### 1. **What is an Event?**

```javascript
// Basic Event Structure
class Event {
  constructor(type, data, source, timestamp = Date.now()) {
    this.type = type;        // Event type/name
    this.data = data;        // Event payload
    this.source = source;    // Event source
    this.timestamp = timestamp;
    this.id = this.generateId();
  }
  
  generateId() {
    return `event-${this.timestamp}-${Math.random().toString(36).substr(2, 9)}`;
  }
}

// Example Events
const userEvents = {
  USER_CREATED: 'user.created',
  USER_UPDATED: 'user.updated',
  USER_DELETED: 'user.deleted'
};

const orderEvents = {
  ORDER_PLACED: 'order.placed',
  ORDER_SHIPPED: 'order.shipped',
  ORDER_CANCELLED: 'order.cancelled'
};
```

## Event-Driven Architecture Patterns

### 1. **Simple Event Emitter**

```javascript
class EventEmitter {
  constructor() {
    this.listeners = new Map();
  }
  
  // Add event listener
  on(eventType, callback) {
    if (!this.listeners.has(eventType)) {
      this.listeners.set(eventType, []);
    }
    this.listeners.get(eventType).push(callback);
    return this; // For chaining
  }
  
  // Remove event listener
  off(eventType, callback) {
    const listeners = this.listeners.get(eventType);
    if (listeners) {
      const index = listeners.indexOf(callback);
      if (index > -1) {
        listeners.splice(index, 1);
      }
    }
    return this;
  }
  
  // Emit event
  emit(eventType, data) {
    const event = new Event(eventType, data, this);
    const listeners = this.listeners.get(eventType) || [];
    
    listeners.forEach(callback => {
      try {
        callback(event);
      } catch (error) {
        console.error(`Error in event listener for ${eventType}:`, error);
      }
    });
    
    // Also emit to wildcard listeners
    const wildcardListeners = this.listeners.get('*') || [];
    wildcardListeners.forEach(callback => callback(event));
  }
  
  // One-time event listener
  once(eventType, callback) {
    const onceWrapper = (event) => {
      this.off(eventType, onceWrapper);
      callback(event);
    };
    return this.on(eventType, onceWrapper);
  }
}
```

### 2. **Advanced Event Bus**

```javascript
class EventBus {
  constructor() {
    this.channels = new Map();
    this.middlewares = [];
  }
  
  // Create or get channel
  channel(name) {
    if (!this.channels.has(name)) {
      this.channels.set(name, new EventEmitter());
    }
    return this.channels.get(name);
  }
  
  // Publish to channel
  publish(channelName, eventType, data) {
    const channel = this.channel(channelName);
    
    // Apply middlewares
    let event = new Event(eventType, data, this);
    this.middlewares.forEach(middleware => {
      event = middleware(event) || event;
    });
    
    channel.emit(eventType, event);
  }
  
  // Subscribe to channel
  subscribe(channelName, eventType, callback) {
    return this.channel(channelName).on(eventType, callback);
  }
  
  // Add middleware
  use(middleware) {
    this.middlewares.push(middleware);
    return this;
  }
}
```

## Real-World Examples

### 1. **E-commerce System**

```javascript
class ECommerceSystem {
  constructor() {
    this.eventBus = new EventBus();
    this.setupEventHandlers();
  }
  
  setupEventHandlers() {
    // Order processing pipeline
    this.eventBus.subscribe('orders', 'order.placed', this.handleOrderPlaced.bind(this));
    this.eventBus.subscribe('inventory', 'order.placed', this.handleInventoryUpdate.bind(this));
    this.eventBus.subscribe('shipping', 'order.ready', this.handleShipping.bind(this));
    this.eventBus.subscribe('notifications', 'order.shipped', this.handleCustomerNotification.bind(this));
    
    // Payment events
    this.eventBus.subscribe('payments', 'payment.processed', this.handlePaymentProcessed.bind(this));
    this.eventBus.subscribe('orders', 'payment.failed', this.handlePaymentFailed.bind(this));
  }
  
  // Event handlers
  handleOrderPlaced(event) {
    const { orderId, items, customer } = event.data;
    console.log(`Processing order ${orderId} for customer ${customer.email}`);
    
    // Check inventory
    const inventoryAvailable = this.checkInventory(items);
    
    if (inventoryAvailable) {
      this.eventBus.publish('orders', 'order.ready', { 
        orderId, 
        items, 
        customer 
      });
    } else {
      this.eventBus.publish('orders', 'order.failed', { 
        orderId, 
        reason: 'Insufficient inventory' 
      });
    }
  }
  
  handleInventoryUpdate(event) {
    const { items } = event.data;
    console.log('Updating inventory levels...');
    // Update inventory logic here
  }
  
  handleShipping(event) {
    const { orderId, customer } = event.data;
    console.log(`Shipping order ${orderId} to ${customer.address}`);
    
    // Simulate shipping process
    setTimeout(() => {
      this.eventBus.publish('shipping', 'order.shipped', { 
        orderId, 
        trackingNumber: 'TRK' + Math.random().toString(36).substr(2, 9).toUpperCase(),
        shippedAt: new Date()
      });
    }, 2000);
  }
  
  handleCustomerNotification(event) {
    const { orderId, trackingNumber } = event.data;
    console.log(`Notifying customer about shipment: Order ${orderId}, Tracking: ${trackingNumber}`);
    // Send email/SMS notification
  }
  
  handlePaymentProcessed(event) {
    const { orderId, amount } = event.data;
    console.log(`Payment processed for order ${orderId}: $${amount}`);
    this.eventBus.publish('orders', 'order.placed', event.data);
  }
  
  handlePaymentFailed(event) {
    const { orderId, reason } = event.data;
    console.log(`Payment failed for order ${orderId}: ${reason}`);
    this.eventBus.publish('orders', 'order.cancelled', { 
      orderId, 
      reason: `Payment failed: ${reason}` 
    });
  }
  
  checkInventory(items) {
    // Simulate inventory check
    return Math.random() > 0.1; // 90% chance of availability
  }
  
  // Business methods
  placeOrder(orderData) {
    console.log('Placing order...');
    this.eventBus.publish('orders', 'order.placed', orderData);
  }
  
  processPayment(paymentData) {
    console.log('Processing payment...');
    // Simulate payment processing
    const success = Math.random() > 0.2; // 80% success rate
    
    if (success) {
      this.eventBus.publish('payments', 'payment.processed', paymentData);
    } else {
      this.eventBus.publish('payments', 'payment.failed', { 
        ...paymentData, 
        reason: 'Insufficient funds' 
      });
    }
  }
}

// Usage
const ecommerce = new ECommerceSystem();

// Simulate customer interaction
ecommerce.placeOrder({
  orderId: 'ORD-12345',
  items: [{ productId: 'P1', quantity: 2 }, { productId: 'P2', quantity: 1 }],
  customer: { email: 'customer@example.com', address: '123 Main St' },
  amount: 99.99
});

ecommerce.processPayment({
  orderId: 'ORD-12345',
  amount: 99.99,
  paymentMethod: 'credit_card'
});
```

### 2. **Real-Time Dashboard System**

```javascript
class RealTimeDashboard {
  constructor() {
    this.eventBus = new EventBus();
    this.metrics = new Map();
    this.setupMetricsCollectors();
    this.setupVisualizations();
  }
  
  setupMetricsCollectors() {
    // System metrics
    setInterval(() => {
      this.eventBus.publish('metrics', 'system.cpu', {
        value: Math.random() * 100,
        timestamp: new Date()
      });
    }, 1000);
    
    setInterval(() => {
      this.eventBus.publish('metrics', 'system.memory', {
        value: Math.random() * 100,
        timestamp: new Date()
      });
    }, 1500);
    
    // Business metrics
    setInterval(() => {
      this.eventBus.publish('metrics', 'sales.revenue', {
        value: Math.random() * 10000,
        timestamp: new Date()
      });
    }, 2000);
  }
  
  setupVisualizations() {
    // CPU Monitor
    this.eventBus.subscribe('metrics', 'system.cpu', (event) => {
      this.updateMetric('cpu', event.data.value);
      this.renderCPUChart(event.data);
    });
    
    // Memory Monitor
    this.eventBus.subscribe('metrics', 'system.memory', (event) => {
      this.updateMetric('memory', event.data.value);
      this.renderMemoryChart(event.data);
    });
    
    // Revenue Counter
    this.eventBus.subscribe('metrics', 'sales.revenue', (event) => {
      this.updateMetric('revenue', event.data.value);
      this.renderRevenueCounter(event.data);
    });
    
    // Alert system
    this.eventBus.subscribe('metrics', '*', (event) => {
      this.checkAlerts(event);
    });
  }
  
  updateMetric(name, value) {
    if (!this.metrics.has(name)) {
      this.metrics.set(name, []);
    }
    
    const metrics = this.metrics.get(name);
    metrics.push({ value, timestamp: new Date() });
    
    // Keep only last 100 values
    if (metrics.length > 100) {
      metrics.shift();
    }
  }
  
  checkAlerts(event) {
    const { type, data } = event;
    
    // CPU alert
    if (type === 'system.cpu' && data.value > 80) {
      this.eventBus.publish('alerts', 'high.cpu', {
        severity: 'warning',
        message: `High CPU usage: ${data.value.toFixed(1)}%`,
        timestamp: new Date()
      });
    }
    
    // Revenue alert
    if (type === 'sales.revenue' && data.value > 8000) {
      this.eventBus.publish('alerts', 'high.revenue', {
        severity: 'info',
        message: `High revenue: $${data.value.toFixed(2)}`,
        timestamp: new Date()
      });
    }
  }
  
  // Visualization methods (simplified)
  renderCPUChart(data) {
    console.log(`CPU: ${data.value.toFixed(1)}%`);
  }
  
  renderMemoryChart(data) {
    console.log(`Memory: ${data.value.toFixed(1)}%`);
  }
  
  renderRevenueCounter(data) {
    console.log(`Revenue: $${data.value.toFixed(2)}`);
  }
  
  // Add alert handler
  onAlert(callback) {
    this.eventBus.subscribe('alerts', '*', callback);
  }
}

// Usage
const dashboard = new RealTimeDashboard();

// Handle alerts
dashboard.onAlert((event) => {
  console.log(`🚨 ALERT [${event.data.severity}]: ${event.data.message}`);
});
```

## Advanced EDA Patterns

### 1. **Event Sourcing**

```javascript
class EventStore {
  constructor() {
    this.events = [];
    this.projections = new Map();
  }
  
  appendEvent(aggregateId, eventType, data) {
    const event = {
      id: this.events.length + 1,
      aggregateId,
      type: eventType,
      data,
      timestamp: new Date(),
      version: this.getNextVersion(aggregateId)
    };
    
    this.events.push(event);
    
    // Update projections
    this.updateProjections(event);
    
    return event;
  }
  
  getNextVersion(aggregateId) {
    const aggregateEvents = this.events.filter(e => e.aggregateId === aggregateId);
    return aggregateEvents.length + 1;
  }
  
  getEvents(aggregateId) {
    return this.events.filter(event => event.aggregateId === aggregateId);
  }
  
  // Rebuild state from events
  rebuildState(aggregateId, initialState = {}) {
    return this.getEvents(aggregateId).reduce((state, event) => {
      return this.applyEvent(state, event);
    }, initialState);
  }
  
  applyEvent(state, event) {
    switch (event.type) {
      case 'user.created':
        return { ...state, ...event.data, status: 'active' };
      case 'user.updated':
        return { ...state, ...event.data };
      case 'user.deactivated':
        return { ...state, status: 'inactive' };
      default:
        return state;
    }
  }
  
  // Projections for read models
  createProjection(name, eventHandler) {
    this.projections.set(name, eventHandler);
  }
  
  updateProjections(event) {
    this.projections.forEach(handler => {
      handler(event);
    });
  }
}

// Usage
const eventStore = new EventStore();

// Create projections
eventStore.createProjection('user-projection', (event) => {
  if (event.type.startsWith('user.')) {
    console.log(`Projection updated for user ${event.aggregateId}: ${event.type}`);
  }
});

// Store events
eventStore.appendEvent('user-123', 'user.created', { 
  name: 'John Doe', 
  email: 'john@example.com' 
});

eventStore.appendEvent('user-123', 'user.updated', { 
  name: 'John Smith' 
});

// Rebuild state
const userState = eventStore.rebuildState('user-123');
console.log('User state:', userState);
```

### 2. **CQRS (Command Query Responsibility Segregation)**

```javascript
class CQRSSystem {
  constructor() {
    this.commands = new EventBus();
    this.queries = new EventBus();
    this.readModels = new Map();
    this.setupCommandHandlers();
    this.setupQueryHandlers();
  }
  
  setupCommandHandlers() {
    // Command handlers (write side)
    this.commands.subscribe('users', 'createUser', async (event) => {
      const user = await this.createUserInDatabase(event.data);
      this.commands.publish('events', 'userCreated', user);
    });
    
    this.commands.subscribe('users', 'updateUser', async (event) => {
      const user = await this.updateUserInDatabase(event.data);
      this.commands.publish('events', 'userUpdated', user);
    });
  }
  
  setupQueryHandlers() {
    // Query handlers (read side)
    this.queries.subscribe('users', 'getUser', async (event) => {
      return this.readModels.get('users')?.get(event.data.userId);
    });
    
    this.queries.subscribe('users', 'listUsers', async (event) => {
      return Array.from(this.readModels.get('users')?.values() || []);
    });
    
    // Update read models from events
    this.commands.subscribe('events', 'userCreated', (event) => {
      this.updateReadModel('users', event.data.id, event.data);
    });
    
    this.commands.subscribe('events', 'userUpdated', (event) => {
      this.updateReadModel('users', event.data.id, event.data);
    });
  }
  
  updateReadModel(modelName, id, data) {
    if (!this.readModels.has(modelName)) {
      this.readModels.set(modelName, new Map());
    }
    this.readModels.get(modelName).set(id, data);
  }
  
  async executeCommand(command, data) {
    return new Promise((resolve) => {
      this.commands.publish('users', command, data);
      resolve({ status: 'accepted' });
    });
  }
  
  async executeQuery(query, data) {
    return new Promise((resolve) => {
      this.queries.subscribe('users', query, (event) => {
        resolve(event.data);
      }).once(); // Listen only once
      
      this.queries.publish('users', query, data);
    });
  }
  
  // Simulated database operations
  async createUserInDatabase(data) {
    return { ...data, id: 'user-' + Date.now(), createdAt: new Date() };
  }
  
  async updateUserInDatabase(data) {
    return { ...data, updatedAt: new Date() };
  }
}

// Usage
const cqrs = new CQRSSystem();

async function demo() {
  // Execute command
  await cqrs.executeCommand('createUser', { 
    name: 'Alice', 
    email: 'alice@example.com' 
  });
  
  // Execute query
  const users = await cqrs.executeQuery('listUsers');
  console.log('Users:', users);
}

demo();
```

## Benefits of Event-Driven Architecture

1. **Loose Coupling**: Components don't need to know about each other
2. **Scalability**: Easy to add new event consumers
3. **Extensibility**: New features can be added without modifying existing code
4. **Resilience**: Failure in one component doesn't break the entire system
5. **Auditability**: Events provide natural audit trail
6. **Real-time processing**: Immediate reaction to changes

## Common Use Cases

- **Microservices communication**
- **Real-time applications**
- **Financial systems**
- **IoT systems**
- **Game development**
- **Logging and monitoring systems**

Event-Driven Architecture is particularly powerful for building responsive, scalable systems that need to react to changes in real-time.