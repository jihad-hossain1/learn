# Part 7: Pub/Sub Messaging

This part covers Redis Pub/Sub (Publish/Subscribe) messaging patterns, real-time communication, and event-driven architectures.

## Understanding Redis Pub/Sub

### Basic Concepts

- **Publisher**: Sends messages to channels
- **Subscriber**: Receives messages from channels
- **Channel**: Named communication pathway
- **Pattern Subscription**: Subscribe to multiple channels using patterns
- **Message**: Data sent through channels

### Pub/Sub vs Streams

| Feature | Pub/Sub | Streams |
|---------|---------|----------|
| Message Persistence | No | Yes |
| Message History | No | Yes |
| Consumer Groups | No | Yes |
| Delivery Guarantee | At most once | At least once |
| Memory Usage | Low | Higher |
| Use Case | Real-time notifications | Event sourcing, logs |

## Basic Pub/Sub Implementation

### Simple Publisher

```typescript
// src/pubsub/publisher.ts
import { createClient, RedisClientType } from 'redis';

interface PublisherConfig {
  host: string;
  port: number;
  password?: string;
  database?: number;
}

interface MessagePayload {
  id: string;
  timestamp: Date;
  type: string;
  data: any;
  metadata?: Record<string, any>;
}

class RedisPublisher {
  private client: RedisClientType;
  private isConnected: boolean = false;
  private messageCount: number = 0;

  constructor(config: PublisherConfig) {
    this.client = createClient({
      socket: {
        host: config.host,
        port: config.port
      },
      password: config.password,
      database: config.database
    });

    this.setupEventHandlers();
  }

  private setupEventHandlers(): void {
    this.client.on('connect', () => {
      console.log('Publisher connecting to Redis...');
    });

    this.client.on('ready', () => {
      this.isConnected = true;
      console.log('Publisher connected to Redis');
    });

    this.client.on('error', (error) => {
      console.error('Publisher Redis error:', error);
      this.isConnected = false;
    });

    this.client.on('end', () => {
      this.isConnected = false;
      console.log('Publisher disconnected from Redis');
    });
  }

  async connect(): Promise<void> {
    if (!this.isConnected) {
      await this.client.connect();
    }
  }

  async disconnect(): Promise<void> {
    if (this.isConnected) {
      await this.client.quit();
    }
  }

  async publish(channel: string, message: any): Promise<number> {
    if (!this.isConnected) {
      throw new Error('Publisher not connected to Redis');
    }

    const payload: MessagePayload = {
      id: this.generateMessageId(),
      timestamp: new Date(),
      type: typeof message,
      data: message,
      metadata: {
        publisher: 'redis-publisher',
        messageNumber: ++this.messageCount
      }
    };

    const serializedMessage = JSON.stringify(payload);
    const subscriberCount = await this.client.publish(channel, serializedMessage);
    
    console.log(`Published message to ${channel}: ${subscriberCount} subscribers`);
    return subscriberCount;
  }

  async publishRaw(channel: string, message: string): Promise<number> {
    if (!this.isConnected) {
      throw new Error('Publisher not connected to Redis');
    }

    return await this.client.publish(channel, message);
  }

  async publishToMultiple(channels: string[], message: any): Promise<Record<string, number>> {
    const results: Record<string, number> = {};
    
    for (const channel of channels) {
      results[channel] = await this.publish(channel, message);
    }
    
    return results;
  }

  async publishWithTTL(channel: string, message: any, ttlSeconds: number): Promise<number> {
    // Note: Redis Pub/Sub doesn't support TTL directly
    // This is a workaround using a separate key for TTL tracking
    const messageId = this.generateMessageId();
    const ttlKey = `pubsub:ttl:${messageId}`;
    
    // Set TTL key
    await this.client.setEx(ttlKey, ttlSeconds, '1');
    
    // Add TTL info to message
    const messageWithTTL = {
      ...message,
      _ttl: {
        messageId,
        expiresAt: new Date(Date.now() + ttlSeconds * 1000)
      }
    };
    
    return await this.publish(channel, messageWithTTL);
  }

  private generateMessageId(): string {
    return `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  getMessageCount(): number {
    return this.messageCount;
  }

  isPublisherConnected(): boolean {
    return this.isConnected;
  }
}

export { RedisPublisher, PublisherConfig, MessagePayload };
```

### Simple Subscriber

```typescript
// src/pubsub/subscriber.ts
import { createClient, RedisClientType } from 'redis';
import { EventEmitter } from 'events';

interface SubscriberConfig {
  host: string;
  port: number;
  password?: string;
  database?: number;
}

interface SubscriptionInfo {
  channel: string;
  subscribedAt: Date;
  messageCount: number;
}

class RedisSubscriber extends EventEmitter {
  private client: RedisClientType;
  private isConnected: boolean = false;
  private subscriptions: Map<string, SubscriptionInfo> = new Map();
  private patternSubscriptions: Map<string, SubscriptionInfo> = new Map();

  constructor(config: SubscriberConfig) {
    super();
    this.client = createClient({
      socket: {
        host: config.host,
        port: config.port
      },
      password: config.password,
      database: config.database
    });

    this.setupEventHandlers();
  }

  private setupEventHandlers(): void {
    this.client.on('connect', () => {
      console.log('Subscriber connecting to Redis...');
    });

    this.client.on('ready', () => {
      this.isConnected = true;
      console.log('Subscriber connected to Redis');
    });

    this.client.on('error', (error) => {
      console.error('Subscriber Redis error:', error);
      this.isConnected = false;
      this.emit('error', error);
    });

    this.client.on('end', () => {
      this.isConnected = false;
      console.log('Subscriber disconnected from Redis');
      this.emit('disconnected');
    });
  }

  async connect(): Promise<void> {
    if (!this.isConnected) {
      await this.client.connect();
    }
  }

  async disconnect(): Promise<void> {
    if (this.isConnected) {
      await this.client.quit();
    }
  }

  async subscribe(channel: string, handler?: (message: string, channel: string) => void): Promise<void> {
    if (!this.isConnected) {
      throw new Error('Subscriber not connected to Redis');
    }

    // Set up message handler
    if (handler) {
      this.client.on('message', (receivedChannel, message) => {
        if (receivedChannel === channel) {
          const subscription = this.subscriptions.get(channel);
          if (subscription) {
            subscription.messageCount++;
          }
          handler(message, receivedChannel);
        }
      });
    } else {
      // Default handler that emits events
      this.client.on('message', (receivedChannel, message) => {
        if (receivedChannel === channel) {
          const subscription = this.subscriptions.get(channel);
          if (subscription) {
            subscription.messageCount++;
          }
          this.emit('message', { channel: receivedChannel, message });
        }
      });
    }

    await this.client.subscribe(channel);
    
    this.subscriptions.set(channel, {
      channel,
      subscribedAt: new Date(),
      messageCount: 0
    });

    console.log(`Subscribed to channel: ${channel}`);
    this.emit('subscribed', channel);
  }

  async subscribeToMultiple(channels: string[], handler?: (message: string, channel: string) => void): Promise<void> {
    for (const channel of channels) {
      await this.subscribe(channel, handler);
    }
  }

  async subscribeToPattern(pattern: string, handler?: (message: string, channel: string) => void): Promise<void> {
    if (!this.isConnected) {
      throw new Error('Subscriber not connected to Redis');
    }

    // Set up pattern message handler
    if (handler) {
      this.client.on('pmessage', (pattern, receivedChannel, message) => {
        const subscription = this.patternSubscriptions.get(pattern);
        if (subscription) {
          subscription.messageCount++;
        }
        handler(message, receivedChannel);
      });
    } else {
      this.client.on('pmessage', (receivedPattern, receivedChannel, message) => {
        const subscription = this.patternSubscriptions.get(receivedPattern);
        if (subscription) {
          subscription.messageCount++;
        }
        this.emit('patternMessage', { pattern: receivedPattern, channel: receivedChannel, message });
      });
    }

    await this.client.pSubscribe(pattern);
    
    this.patternSubscriptions.set(pattern, {
      channel: pattern,
      subscribedAt: new Date(),
      messageCount: 0
    });

    console.log(`Subscribed to pattern: ${pattern}`);
    this.emit('patternSubscribed', pattern);
  }

  async unsubscribe(channel: string): Promise<void> {
    if (!this.isConnected) {
      throw new Error('Subscriber not connected to Redis');
    }

    await this.client.unsubscribe(channel);
    this.subscriptions.delete(channel);
    
    console.log(`Unsubscribed from channel: ${channel}`);
    this.emit('unsubscribed', channel);
  }

  async unsubscribeFromPattern(pattern: string): Promise<void> {
    if (!this.isConnected) {
      throw new Error('Subscriber not connected to Redis');
    }

    await this.client.pUnsubscribe(pattern);
    this.patternSubscriptions.delete(pattern);
    
    console.log(`Unsubscribed from pattern: ${pattern}`);
    this.emit('patternUnsubscribed', pattern);
  }

  async unsubscribeAll(): Promise<void> {
    if (!this.isConnected) {
      throw new Error('Subscriber not connected to Redis');
    }

    // Unsubscribe from all channels
    await this.client.unsubscribe();
    this.subscriptions.clear();

    // Unsubscribe from all patterns
    await this.client.pUnsubscribe();
    this.patternSubscriptions.clear();

    console.log('Unsubscribed from all channels and patterns');
    this.emit('unsubscribedAll');
  }

  getSubscriptions(): SubscriptionInfo[] {
    return Array.from(this.subscriptions.values());
  }

  getPatternSubscriptions(): SubscriptionInfo[] {
    return Array.from(this.patternSubscriptions.values());
  }

  isSubscribedTo(channel: string): boolean {
    return this.subscriptions.has(channel);
  }

  isSubscribedToPattern(pattern: string): boolean {
    return this.patternSubscriptions.has(pattern);
  }

  getTotalMessageCount(): number {
    let total = 0;
    for (const subscription of this.subscriptions.values()) {
      total += subscription.messageCount;
    }
    for (const subscription of this.patternSubscriptions.values()) {
      total += subscription.messageCount;
    }
    return total;
  }

  isSubscriberConnected(): boolean {
    return this.isConnected;
  }
}

export { RedisSubscriber, SubscriberConfig, SubscriptionInfo };
```

## Advanced Pub/Sub Patterns

### Message Queue with Acknowledgment

```typescript
// src/pubsub/reliable-messaging.ts
import { RedisPublisher } from './publisher';
import { RedisSubscriber } from './subscriber';
import { createClient, RedisClientType } from 'redis';
import { EventEmitter } from 'events';

interface ReliableMessage {
  id: string;
  channel: string;
  payload: any;
  timestamp: Date;
  retryCount: number;
  maxRetries: number;
  ttl: number;
}

interface AckMessage {
  messageId: string;
  status: 'ack' | 'nack';
  reason?: string;
}

class ReliableMessaging extends EventEmitter {
  private publisher: RedisPublisher;
  private subscriber: RedisSubscriber;
  private ackClient: RedisClientType;
  private pendingMessages: Map<string, ReliableMessage> = new Map();
  private ackTimeout: number = 30000; // 30 seconds
  private retryInterval: NodeJS.Timeout | null = null;

  constructor(config: any) {
    super();
    this.publisher = new RedisPublisher(config);
    this.subscriber = new RedisSubscriber(config);
    this.ackClient = createClient({
      socket: {
        host: config.host,
        port: config.port
      },
      password: config.password,
      database: config.database
    });

    this.setupAckHandling();
    this.startRetryLoop();
  }

  async connect(): Promise<void> {
    await Promise.all([
      this.publisher.connect(),
      this.subscriber.connect(),
      this.ackClient.connect()
    ]);
  }

  async disconnect(): Promise<void> {
    if (this.retryInterval) {
      clearInterval(this.retryInterval);
    }
    
    await Promise.all([
      this.publisher.disconnect(),
      this.subscriber.disconnect(),
      this.ackClient.quit()
    ]);
  }

  private setupAckHandling(): void {
    // Subscribe to acknowledgment channel
    this.subscriber.subscribe('_ack_channel', (message) => {
      try {
        const ackMessage: AckMessage = JSON.parse(message);
        this.handleAcknowledgment(ackMessage);
      } catch (error) {
        console.error('Failed to parse acknowledgment message:', error);
      }
    });
  }

  private handleAcknowledgment(ackMessage: AckMessage): void {
    const pendingMessage = this.pendingMessages.get(ackMessage.messageId);
    
    if (pendingMessage) {
      if (ackMessage.status === 'ack') {
        this.pendingMessages.delete(ackMessage.messageId);
        this.emit('messageAcknowledged', {
          messageId: ackMessage.messageId,
          message: pendingMessage
        });
      } else {
        this.emit('messageNacked', {
          messageId: ackMessage.messageId,
          message: pendingMessage,
          reason: ackMessage.reason
        });
        // Message will be retried in the retry loop
      }
    }
  }

  async publishReliable(
    channel: string, 
    payload: any, 
    options: {
      maxRetries?: number;
      ttl?: number;
    } = {}
  ): Promise<string> {
    const message: ReliableMessage = {
      id: this.generateMessageId(),
      channel,
      payload,
      timestamp: new Date(),
      retryCount: 0,
      maxRetries: options.maxRetries || 3,
      ttl: options.ttl || 300000 // 5 minutes
    };

    // Store message for tracking
    this.pendingMessages.set(message.id, message);

    // Publish message with acknowledgment request
    const messageWithAck = {
      ...message,
      requiresAck: true,
      ackChannel: '_ack_channel'
    };

    await this.publisher.publish(channel, messageWithAck);
    
    this.emit('messageSent', message);
    return message.id;
  }

  async subscribeReliable(
    channel: string, 
    handler: (payload: any, messageId: string) => Promise<boolean>
  ): Promise<void> {
    await this.subscriber.subscribe(channel, async (message) => {
      try {
        const parsedMessage = JSON.parse(message);
        
        if (parsedMessage.requiresAck) {
          try {
            const success = await handler(parsedMessage.payload, parsedMessage.id);
            
            // Send acknowledgment
            const ackMessage: AckMessage = {
              messageId: parsedMessage.id,
              status: success ? 'ack' : 'nack',
              reason: success ? undefined : 'Handler returned false'
            };
            
            await this.publisher.publishRaw('_ack_channel', JSON.stringify(ackMessage));
          } catch (error) {
            // Send negative acknowledgment
            const ackMessage: AckMessage = {
              messageId: parsedMessage.id,
              status: 'nack',
              reason: error instanceof Error ? error.message : 'Handler threw error'
            };
            
            await this.publisher.publishRaw('_ack_channel', JSON.stringify(ackMessage));
          }
        } else {
          // Regular message without acknowledgment
          await handler(parsedMessage.payload || parsedMessage, parsedMessage.id);
        }
      } catch (error) {
        console.error('Failed to process message:', error);
      }
    });
  }

  private startRetryLoop(): void {
    this.retryInterval = setInterval(() => {
      const now = Date.now();
      
      for (const [messageId, message] of this.pendingMessages.entries()) {
        const messageAge = now - message.timestamp.getTime();
        
        // Check if message has expired
        if (messageAge > message.ttl) {
          this.pendingMessages.delete(messageId);
          this.emit('messageExpired', message);
          continue;
        }
        
        // Check if message needs retry
        const timeSinceLastAttempt = now - message.timestamp.getTime() - (message.retryCount * this.ackTimeout);
        
        if (timeSinceLastAttempt > this.ackTimeout && message.retryCount < message.maxRetries) {
          message.retryCount++;
          
          // Retry publishing
          const messageWithAck = {
            ...message,
            requiresAck: true,
            ackChannel: '_ack_channel',
            isRetry: true
          };
          
          this.publisher.publish(message.channel, messageWithAck)
            .then(() => {
              this.emit('messageRetried', message);
            })
            .catch((error) => {
              console.error('Failed to retry message:', error);
            });
        } else if (message.retryCount >= message.maxRetries) {
          // Max retries exceeded
          this.pendingMessages.delete(messageId);
          this.emit('messageMaxRetriesExceeded', message);
        }
      }
    }, 10000); // Check every 10 seconds
  }

  private generateMessageId(): string {
    return `reliable_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  getPendingMessages(): ReliableMessage[] {
    return Array.from(this.pendingMessages.values());
  }

  getPendingMessageCount(): number {
    return this.pendingMessages.size;
  }

  setAckTimeout(timeout: number): void {
    this.ackTimeout = timeout;
  }
}

export { ReliableMessaging, ReliableMessage, AckMessage };
```

### Event Bus Implementation

```typescript
// src/pubsub/event-bus.ts
import { RedisPublisher } from './publisher';
import { RedisSubscriber } from './subscriber';
import { EventEmitter } from 'events';

interface EventHandler {
  id: string;
  eventType: string;
  handler: (data: any, metadata: EventMetadata) => Promise<void> | void;
  options: EventHandlerOptions;
}

interface EventHandlerOptions {
  priority?: number;
  once?: boolean;
  filter?: (data: any) => boolean;
  timeout?: number;
}

interface EventMetadata {
  eventId: string;
  eventType: string;
  timestamp: Date;
  source: string;
  correlationId?: string;
  userId?: string;
}

interface Event {
  id: string;
  type: string;
  data: any;
  metadata: EventMetadata;
}

class EventBus extends EventEmitter {
  private publisher: RedisPublisher;
  private subscriber: RedisSubscriber;
  private handlers: Map<string, EventHandler[]> = new Map();
  private isConnected: boolean = false;

  constructor(config: any) {
    super();
    this.publisher = new RedisPublisher(config);
    this.subscriber = new RedisSubscriber(config);
    this.setupEventHandling();
  }

  async connect(): Promise<void> {
    await Promise.all([
      this.publisher.connect(),
      this.subscriber.connect()
    ]);
    this.isConnected = true;
  }

  async disconnect(): Promise<void> {
    await Promise.all([
      this.publisher.disconnect(),
      this.subscriber.disconnect()
    ]);
    this.isConnected = false;
  }

  private setupEventHandling(): void {
    // Subscribe to all event channels using pattern
    this.subscriber.subscribeToPattern('event:*', async (message, channel) => {
      try {
        const event: Event = JSON.parse(message);
        await this.processEvent(event);
      } catch (error) {
        console.error('Failed to process event:', error);
        this.emit('eventProcessingError', { error, message, channel });
      }
    });
  }

  async emit(eventType: string, data: any, metadata: Partial<EventMetadata> = {}): Promise<string> {
    if (!this.isConnected) {
      throw new Error('EventBus not connected');
    }

    const event: Event = {
      id: this.generateEventId(),
      type: eventType,
      data,
      metadata: {
        eventId: this.generateEventId(),
        eventType,
        timestamp: new Date(),
        source: 'event-bus',
        ...metadata
      }
    };

    const channel = `event:${eventType}`;
    await this.publisher.publish(channel, event);
    
    this.emit('eventEmitted', event);
    return event.id;
  }

  on(
    eventType: string, 
    handler: (data: any, metadata: EventMetadata) => Promise<void> | void,
    options: EventHandlerOptions = {}
  ): string {
    const handlerId = this.generateHandlerId();
    
    const eventHandler: EventHandler = {
      id: handlerId,
      eventType,
      handler,
      options: {
        priority: 0,
        once: false,
        ...options
      }
    };

    if (!this.handlers.has(eventType)) {
      this.handlers.set(eventType, []);
    }

    const handlers = this.handlers.get(eventType)!;
    handlers.push(eventHandler);
    
    // Sort by priority (higher priority first)
    handlers.sort((a, b) => (b.options.priority || 0) - (a.options.priority || 0));

    return handlerId;
  }

  once(
    eventType: string,
    handler: (data: any, metadata: EventMetadata) => Promise<void> | void,
    options: EventHandlerOptions = {}
  ): string {
    return this.on(eventType, handler, { ...options, once: true });
  }

  off(eventType: string, handlerId?: string): void {
    if (!this.handlers.has(eventType)) {
      return;
    }

    const handlers = this.handlers.get(eventType)!;
    
    if (handlerId) {
      const index = handlers.findIndex(h => h.id === handlerId);
      if (index !== -1) {
        handlers.splice(index, 1);
      }
    } else {
      // Remove all handlers for this event type
      this.handlers.delete(eventType);
    }
  }

  private async processEvent(event: Event): Promise<void> {
    const handlers = this.handlers.get(event.type) || [];
    
    for (const handler of handlers) {
      try {
        // Apply filter if specified
        if (handler.options.filter && !handler.options.filter(event.data)) {
          continue;
        }

        // Execute handler with timeout if specified
        if (handler.options.timeout) {
          await this.executeWithTimeout(
            () => handler.handler(event.data, event.metadata),
            handler.options.timeout
          );
        } else {
          await handler.handler(event.data, event.metadata);
        }

        // Remove handler if it's a one-time handler
        if (handler.options.once) {
          this.off(event.type, handler.id);
        }

        this.emit('eventHandled', {
          event,
          handler: handler.id
        });

      } catch (error) {
        console.error(`Handler ${handler.id} failed for event ${event.type}:`, error);
        this.emit('handlerError', {
          error,
          event,
          handler: handler.id
        });
      }
    }
  }

  private async executeWithTimeout<T>(
    operation: () => Promise<T> | T,
    timeoutMs: number
  ): Promise<T> {
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error(`Handler timeout after ${timeoutMs}ms`));
      }, timeoutMs);

      Promise.resolve(operation())
        .then(result => {
          clearTimeout(timeout);
          resolve(result);
        })
        .catch(error => {
          clearTimeout(timeout);
          reject(error);
        });
    });
  }

  async waitFor(
    eventType: string,
    filter?: (data: any) => boolean,
    timeoutMs: number = 30000
  ): Promise<{ data: any; metadata: EventMetadata }> {
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        this.off(eventType, handlerId);
        reject(new Error(`Timeout waiting for event ${eventType}`));
      }, timeoutMs);

      const handlerId = this.once(eventType, (data, metadata) => {
        clearTimeout(timeout);
        resolve({ data, metadata });
      }, { filter });
    });
  }

  getHandlers(eventType?: string): EventHandler[] {
    if (eventType) {
      return this.handlers.get(eventType) || [];
    }
    
    const allHandlers: EventHandler[] = [];
    for (const handlers of this.handlers.values()) {
      allHandlers.push(...handlers);
    }
    return allHandlers;
  }

  getEventTypes(): string[] {
    return Array.from(this.handlers.keys());
  }

  private generateEventId(): string {
    return `evt_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private generateHandlerId(): string {
    return `handler_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

export { EventBus, Event, EventHandler, EventHandlerOptions, EventMetadata };
```

## Real-time Applications

### Chat Application

```typescript
// src/examples/chat-application.ts
import { EventBus } from '../pubsub/event-bus';
import { RedisPublisher } from '../pubsub/publisher';
import { RedisSubscriber } from '../pubsub/subscriber';
import express from 'express';
import { createServer } from 'http';
import { Server as SocketIOServer } from 'socket.io';

interface ChatMessage {
  id: string;
  roomId: string;
  userId: string;
  username: string;
  message: string;
  timestamp: Date;
  type: 'text' | 'image' | 'file';
}

interface ChatRoom {
  id: string;
  name: string;
  participants: string[];
  createdAt: Date;
  isPrivate: boolean;
}

interface UserPresence {
  userId: string;
  username: string;
  status: 'online' | 'away' | 'offline';
  lastSeen: Date;
  currentRoom?: string;
}

class ChatApplication {
  private eventBus: EventBus;
  private publisher: RedisPublisher;
  private subscriber: RedisSubscriber;
  private app: express.Application;
  private server: any;
  private io: SocketIOServer;
  private rooms: Map<string, ChatRoom> = new Map();
  private userPresence: Map<string, UserPresence> = new Map();

  constructor(redisConfig: any) {
    this.eventBus = new EventBus(redisConfig);
    this.publisher = new RedisPublisher(redisConfig);
    this.subscriber = new RedisSubscriber(redisConfig);
    
    this.app = express();
    this.server = createServer(this.app);
    this.io = new SocketIOServer(this.server, {
      cors: {
        origin: "*",
        methods: ["GET", "POST"]
      }
    });

    this.setupEventHandlers();
    this.setupSocketHandlers();
    this.setupRoutes();
  }

  async start(port: number = 3000): Promise<void> {
    await this.eventBus.connect();
    await this.publisher.connect();
    await this.subscriber.connect();

    this.server.listen(port, () => {
      console.log(`Chat server running on port ${port}`);
    });
  }

  private setupEventHandlers(): void {
    // Handle chat messages
    this.eventBus.on('chat:message', async (data: ChatMessage) => {
      // Broadcast to room participants
      this.io.to(data.roomId).emit('newMessage', data);
      
      // Store message (in real app, you'd use a database)
      console.log(`Message in room ${data.roomId}: ${data.message}`);
    });

    // Handle user presence updates
    this.eventBus.on('user:presence', async (data: UserPresence) => {
      this.userPresence.set(data.userId, data);
      
      // Broadcast presence update to relevant rooms
      if (data.currentRoom) {
        this.io.to(data.currentRoom).emit('presenceUpdate', data);
      }
    });

    // Handle room events
    this.eventBus.on('room:created', async (data: ChatRoom) => {
      this.rooms.set(data.id, data);
      this.io.emit('roomCreated', data);
    });

    this.eventBus.on('room:joined', async (data: { roomId: string; userId: string; username: string }) => {
      const room = this.rooms.get(data.roomId);
      if (room && !room.participants.includes(data.userId)) {
        room.participants.push(data.userId);
        this.io.to(data.roomId).emit('userJoined', data);
      }
    });

    this.eventBus.on('room:left', async (data: { roomId: string; userId: string; username: string }) => {
      const room = this.rooms.get(data.roomId);
      if (room) {
        room.participants = room.participants.filter(id => id !== data.userId);
        this.io.to(data.roomId).emit('userLeft', data);
      }
    });
  }

  private setupSocketHandlers(): void {
    this.io.on('connection', (socket) => {
      console.log(`User connected: ${socket.id}`);

      // Handle user authentication
      socket.on('authenticate', async (data: { userId: string; username: string }) => {
        socket.data.userId = data.userId;
        socket.data.username = data.username;
        
        // Update user presence
        await this.eventBus.emit('user:presence', {
          userId: data.userId,
          username: data.username,
          status: 'online',
          lastSeen: new Date()
        });
        
        socket.emit('authenticated', { success: true });
      });

      // Handle joining rooms
      socket.on('joinRoom', async (data: { roomId: string }) => {
        if (!socket.data.userId) {
          socket.emit('error', { message: 'Not authenticated' });
          return;
        }

        socket.join(data.roomId);
        
        // Update user presence with current room
        await this.eventBus.emit('user:presence', {
          userId: socket.data.userId,
          username: socket.data.username,
          status: 'online',
          lastSeen: new Date(),
          currentRoom: data.roomId
        });

        // Emit room joined event
        await this.eventBus.emit('room:joined', {
          roomId: data.roomId,
          userId: socket.data.userId,
          username: socket.data.username
        });

        socket.emit('joinedRoom', { roomId: data.roomId });
      });

      // Handle leaving rooms
      socket.on('leaveRoom', async (data: { roomId: string }) => {
        if (!socket.data.userId) {
          return;
        }

        socket.leave(data.roomId);
        
        await this.eventBus.emit('room:left', {
          roomId: data.roomId,
          userId: socket.data.userId,
          username: socket.data.username
        });

        socket.emit('leftRoom', { roomId: data.roomId });
      });

      // Handle sending messages
      socket.on('sendMessage', async (data: { roomId: string; message: string; type?: string }) => {
        if (!socket.data.userId) {
          socket.emit('error', { message: 'Not authenticated' });
          return;
        }

        const chatMessage: ChatMessage = {
          id: this.generateMessageId(),
          roomId: data.roomId,
          userId: socket.data.userId,
          username: socket.data.username,
          message: data.message,
          timestamp: new Date(),
          type: (data.type as any) || 'text'
        };

        await this.eventBus.emit('chat:message', chatMessage);
      });

      // Handle typing indicators
      socket.on('typing', (data: { roomId: string; isTyping: boolean }) => {
        if (!socket.data.userId) {
          return;
        }

        socket.to(data.roomId).emit('userTyping', {
          userId: socket.data.userId,
          username: socket.data.username,
          isTyping: data.isTyping
        });
      });

      // Handle disconnection
      socket.on('disconnect', async () => {
        if (socket.data.userId) {
          await this.eventBus.emit('user:presence', {
            userId: socket.data.userId,
            username: socket.data.username,
            status: 'offline',
            lastSeen: new Date()
          });
        }
        
        console.log(`User disconnected: ${socket.id}`);
      });
    });
  }

  private setupRoutes(): void {
    this.app.use(express.json());
    this.app.use(express.static('public'));

    // Create room endpoint
    this.app.post('/api/rooms', async (req, res) => {
      const { name, isPrivate = false } = req.body;
      
      const room: ChatRoom = {
        id: this.generateRoomId(),
        name,
        participants: [],
        createdAt: new Date(),
        isPrivate
      };

      await this.eventBus.emit('room:created', room);
      
      res.json({ success: true, room });
    });

    // Get rooms endpoint
    this.app.get('/api/rooms', (req, res) => {
      const rooms = Array.from(this.rooms.values())
        .filter(room => !room.isPrivate)
        .map(room => ({
          ...room,
          participantCount: room.participants.length
        }));
      
      res.json({ rooms });
    });

    // Get room messages endpoint (simplified)
    this.app.get('/api/rooms/:roomId/messages', (req, res) => {
      // In a real app, you'd fetch from database
      res.json({ messages: [] });
    });

    // Get online users endpoint
    this.app.get('/api/users/online', (req, res) => {
      const onlineUsers = Array.from(this.userPresence.values())
        .filter(user => user.status === 'online');
      
      res.json({ users: onlineUsers });
    });
  }

  private generateMessageId(): string {
    return `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  private generateRoomId(): string {
    return `room_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  async stop(): Promise<void> {
    await this.eventBus.disconnect();
    await this.publisher.disconnect();
    await this.subscriber.disconnect();
    this.server.close();
  }
}

export { ChatApplication, ChatMessage, ChatRoom, UserPresence };
```

### Notification System

```typescript
// src/examples/notification-system.ts
import { EventBus } from '../pubsub/event-bus';
import { ReliableMessaging } from '../pubsub/reliable-messaging';

interface Notification {
  id: string;
  userId: string;
  type: 'email' | 'sms' | 'push' | 'in-app';
  title: string;
  message: string;
  data?: any;
  priority: 'low' | 'normal' | 'high' | 'urgent';
  scheduledFor?: Date;
  expiresAt?: Date;
  channels: string[];
}

interface NotificationTemplate {
  id: string;
  name: string;
  type: string;
  template: {
    title: string;
    message: string;
    variables: string[];
  };
}

interface DeliveryResult {
  notificationId: string;
  channel: string;
  status: 'sent' | 'failed' | 'pending';
  timestamp: Date;
  error?: string;
  metadata?: any;
}

class NotificationSystem {
  private eventBus: EventBus;
  private reliableMessaging: ReliableMessaging;
  private templates: Map<string, NotificationTemplate> = new Map();
  private deliveryResults: Map<string, DeliveryResult[]> = new Map();

  constructor(redisConfig: any) {
    this.eventBus = new EventBus(redisConfig);
    this.reliableMessaging = new ReliableMessaging(redisConfig);
    this.setupEventHandlers();
  }

  async start(): Promise<void> {
    await this.eventBus.connect();
    await this.reliableMessaging.connect();
    console.log('Notification system started');
  }

  async stop(): Promise<void> {
    await this.eventBus.disconnect();
    await this.reliableMessaging.disconnect();
    console.log('Notification system stopped');
  }

  private setupEventHandlers(): void {
    // Handle notification requests
    this.eventBus.on('notification:send', async (notification: Notification) => {
      await this.processNotification(notification);
    });

    // Handle template-based notifications
    this.eventBus.on('notification:template', async (data: {
      templateId: string;
      userId: string;
      variables: Record<string, any>;
      channels: string[];
      priority?: string;
    }) => {
      await this.sendTemplateNotification(data);
    });

    // Handle bulk notifications
    this.eventBus.on('notification:bulk', async (data: {
      userIds: string[];
      notification: Omit<Notification, 'id' | 'userId'>;
    }) => {
      await this.sendBulkNotifications(data);
    });

    // Handle delivery confirmations
    this.eventBus.on('notification:delivered', async (result: DeliveryResult) => {
      this.recordDeliveryResult(result);
    });

    // Handle delivery failures
    this.eventBus.on('notification:failed', async (result: DeliveryResult) => {
      this.recordDeliveryResult(result);
      await this.handleDeliveryFailure(result);
    });
  }

  async sendNotification(notification: Omit<Notification, 'id'>): Promise<string> {
    const fullNotification: Notification = {
      id: this.generateNotificationId(),
      ...notification
    };

    // Validate notification
    if (!this.validateNotification(fullNotification)) {
      throw new Error('Invalid notification data');
    }

    // Check if notification should be scheduled
    if (fullNotification.scheduledFor && fullNotification.scheduledFor > new Date()) {
      await this.scheduleNotification(fullNotification);
      return fullNotification.id;
    }

    // Send immediately
    await this.eventBus.emit('notification:send', fullNotification);
    return fullNotification.id;
  }

  async sendTemplateNotification(data: {
    templateId: string;
    userId: string;
    variables: Record<string, any>;
    channels: string[];
    priority?: string;
  }): Promise<string> {
    const template = this.templates.get(data.templateId);
    if (!template) {
      throw new Error(`Template not found: ${data.templateId}`);
    }

    // Replace variables in template
    const title = this.replaceVariables(template.template.title, data.variables);
    const message = this.replaceVariables(template.template.message, data.variables);

    const notification: Omit<Notification, 'id'> = {
      userId: data.userId,
      type: 'in-app', // Default type
      title,
      message,
      priority: (data.priority as any) || 'normal',
      channels: data.channels,
      data: {
        templateId: data.templateId,
        variables: data.variables
      }
    };

    return await this.sendNotification(notification);
  }

  async sendBulkNotifications(data: {
    userIds: string[];
    notification: Omit<Notification, 'id' | 'userId'>;
  }): Promise<string[]> {
    const notificationIds: string[] = [];

    for (const userId of data.userIds) {
      try {
        const id = await this.sendNotification({
          ...data.notification,
          userId
        });
        notificationIds.push(id);
      } catch (error) {
        console.error(`Failed to send notification to user ${userId}:`, error);
      }
    }

    return notificationIds;
  }

  private async processNotification(notification: Notification): Promise<void> {
    // Check if notification has expired
    if (notification.expiresAt && notification.expiresAt < new Date()) {
      console.log(`Notification ${notification.id} has expired`);
      return;
    }

    // Process each channel
    for (const channel of notification.channels) {
      try {
        await this.sendToChannel(notification, channel);
      } catch (error) {
        console.error(`Failed to send notification ${notification.id} to channel ${channel}:`, error);
        
        await this.eventBus.emit('notification:failed', {
          notificationId: notification.id,
          channel,
          status: 'failed',
          timestamp: new Date(),
          error: error instanceof Error ? error.message : 'Unknown error'
        });
      }
    }
  }

  private async sendToChannel(notification: Notification, channel: string): Promise<void> {
    // Use reliable messaging for critical notifications
    if (notification.priority === 'urgent' || notification.priority === 'high') {
      await this.reliableMessaging.publishReliable(
        `notification:${channel}`,
        {
          notification,
          channel
        },
        {
          maxRetries: 5,
          ttl: 300000 // 5 minutes
        }
      );
    } else {
      await this.eventBus.emit(`notification:${channel}`, {
        notification,
        channel
      });
    }

    // Record successful send
    await this.eventBus.emit('notification:delivered', {
      notificationId: notification.id,
      channel,
      status: 'sent',
      timestamp: new Date()
    });
  }

  private async scheduleNotification(notification: Notification): Promise<void> {
    const delay = notification.scheduledFor!.getTime() - Date.now();
    
    setTimeout(async () => {
      await this.eventBus.emit('notification:send', notification);
    }, delay);

    console.log(`Notification ${notification.id} scheduled for ${notification.scheduledFor}`);
  }

  private validateNotification(notification: Notification): boolean {
    return !!(notification.userId && 
             notification.title && 
             notification.message && 
             notification.channels && 
             notification.channels.length > 0);
  }

  private replaceVariables(template: string, variables: Record<string, any>): string {
    let result = template;
    
    for (const [key, value] of Object.entries(variables)) {
      const placeholder = `{{${key}}}`;
      result = result.replace(new RegExp(placeholder, 'g'), String(value));
    }
    
    return result;
  }

  private recordDeliveryResult(result: DeliveryResult): void {
    if (!this.deliveryResults.has(result.notificationId)) {
      this.deliveryResults.set(result.notificationId, []);
    }
    
    this.deliveryResults.get(result.notificationId)!.push(result);
  }

  private async handleDeliveryFailure(result: DeliveryResult): Promise<void> {
    // Implement retry logic for failed deliveries
    console.log(`Handling delivery failure for notification ${result.notificationId} on channel ${result.channel}`);
    
    // Could implement exponential backoff retry here
  }

  addTemplate(template: NotificationTemplate): void {
    this.templates.set(template.id, template);
  }

  getTemplate(templateId: string): NotificationTemplate | undefined {
    return this.templates.get(templateId);
  }

  getDeliveryResults(notificationId: string): DeliveryResult[] {
    return this.deliveryResults.get(notificationId) || [];
  }

  private generateNotificationId(): string {
    return `notif_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

export { NotificationSystem, Notification, NotificationTemplate, DeliveryResult };
```

## Usage Examples

### Complete Pub/Sub Example

```typescript
// src/examples/pubsub-complete-example.ts
import { RedisPublisher } from '../pubsub/publisher';
import { RedisSubscriber } from '../pubsub/subscriber';
import { EventBus } from '../pubsub/event-bus';
import { ReliableMessaging } from '../pubsub/reliable-messaging';
import { ChatApplication } from './chat-application';
import { NotificationSystem } from './notification-system';

async function basicPubSubExample() {
  console.log('=== Basic Pub/Sub Example ===');
  
  const config = {
    host: 'localhost',
    port: 6379
  };

  const publisher = new RedisPublisher(config);
  const subscriber = new RedisSubscriber(config);

  await publisher.connect();
  await subscriber.connect();

  // Subscribe to channels
  await subscriber.subscribe('news', (message, channel) => {
    console.log(`📰 News: ${message}`);
  });

  await subscriber.subscribe('alerts', (message, channel) => {
    console.log(`🚨 Alert: ${message}`);
  });

  // Subscribe to pattern
  await subscriber.subscribeToPattern('user:*', (message, channel) => {
    console.log(`👤 User event on ${channel}: ${message}`);
  });

  // Publish messages
  await publisher.publish('news', 'Breaking: Redis Pub/Sub is awesome!');
  await publisher.publish('alerts', 'System maintenance in 1 hour');
  await publisher.publish('user:123', 'User logged in');
  await publisher.publish('user:456', 'User updated profile');

  // Wait a bit for messages to be processed
  await new Promise(resolve => setTimeout(resolve, 1000));

  await publisher.disconnect();
  await subscriber.disconnect();
}

async function eventBusExample() {
  console.log('\n=== Event Bus Example ===');
  
  const eventBus = new EventBus({
    host: 'localhost',
    port: 6379
  });

  await eventBus.connect();

  // Register event handlers
  eventBus.on('user:registered', async (data, metadata) => {
    console.log(`🎉 New user registered: ${data.username} (${data.email})`);
    
    // Send welcome email
    await eventBus.emit('email:send', {
      to: data.email,
      subject: 'Welcome!',
      template: 'welcome'
    });
  });

  eventBus.on('email:send', (data) => {
    console.log(`📧 Sending email to ${data.to}: ${data.subject}`);
  });

  // High priority handler
  eventBus.on('user:registered', (data) => {
    console.log(`🔔 High priority: User ${data.username} needs verification`);
  }, { priority: 10 });

  // One-time handler
  eventBus.once('system:startup', () => {
    console.log('🚀 System started successfully!');
  });

  // Emit events
  await eventBus.emit('system:startup', { version: '1.0.0' });
  await eventBus.emit('user:registered', {
    username: 'john_doe',
    email: 'john@example.com',
    id: '123'
  });

  // Wait for event
  try {
    const result = await eventBus.waitFor('user:login', 
      (data) => data.userId === '123', 
      5000
    );
    console.log('User login detected:', result.data);
  } catch (error) {
    console.log('No login detected within timeout');
  }

  await eventBus.disconnect();
}

async function reliableMessagingExample() {
  console.log('\n=== Reliable Messaging Example ===');
  
  const reliableMessaging = new ReliableMessaging({
    host: 'localhost',
    port: 6379
  });

  await reliableMessaging.connect();

  // Set up event handlers
  reliableMessaging.on('messageAcknowledged', (data) => {
    console.log(`✅ Message acknowledged: ${data.messageId}`);
  });

  reliableMessaging.on('messageRetried', (message) => {
    console.log(`🔄 Message retried: ${message.id} (attempt ${message.retryCount})`);
  });

  reliableMessaging.on('messageMaxRetriesExceeded', (message) => {
    console.log(`❌ Message failed after max retries: ${message.id}`);
  });

  // Subscribe to reliable messages
  await reliableMessaging.subscribeReliable('orders', async (payload, messageId) => {
    console.log(`📦 Processing order: ${JSON.stringify(payload)}`);
    
    // Simulate processing
    await new Promise(resolve => setTimeout(resolve, 100));
    
    // Return true for successful processing
    return payload.amount > 0;
  });

  // Publish reliable messages
  const messageId1 = await reliableMessaging.publishReliable('orders', {
    orderId: 'order_123',
    amount: 99.99,
    customerId: 'customer_456'
  });

  const messageId2 = await reliableMessaging.publishReliable('orders', {
    orderId: 'order_124',
    amount: -10, // This will cause nack
    customerId: 'customer_789'
  });

  console.log(`Published messages: ${messageId1}, ${messageId2}`);

  // Wait for processing
  await new Promise(resolve => setTimeout(resolve, 2000));

  console.log('Pending messages:', reliableMessaging.getPendingMessageCount());

  await reliableMessaging.disconnect();
}

async function notificationSystemExample() {
  console.log('\n=== Notification System Example ===');
  
  const notificationSystem = new NotificationSystem({
    host: 'localhost',
    port: 6379
  });

  await notificationSystem.start();

  // Add notification template
  notificationSystem.addTemplate({
    id: 'welcome',
    name: 'Welcome Message',
    type: 'user_onboarding',
    template: {
      title: 'Welcome {{username}}!',
      message: 'Thank you for joining {{appName}}. Your account ID is {{userId}}.',
      variables: ['username', 'appName', 'userId']
    }
  });

  // Send template-based notification
  await notificationSystem.sendTemplateNotification({
    templateId: 'welcome',
    userId: 'user_123',
    variables: {
      username: 'John Doe',
      appName: 'MyApp',
      userId: 'user_123'
    },
    channels: ['email', 'push', 'in-app'],
    priority: 'normal'
  });

  // Send direct notification
  await notificationSystem.sendNotification({
    userId: 'user_456',
    type: 'push',
    title: 'Order Update',
    message: 'Your order #12345 has been shipped!',
    priority: 'high',
    channels: ['push', 'email'],
    data: {
      orderId: '12345',
      trackingNumber: 'TRK789'
    }
  });

  // Send bulk notifications
  await notificationSystem.sendBulkNotifications({
    userIds: ['user_1', 'user_2', 'user_3'],
    notification: {
      type: 'in-app',
      title: 'System Maintenance',
      message: 'Scheduled maintenance will occur tonight at 2 AM.',
      priority: 'normal',
      channels: ['in-app', 'email']
    }
  });

  await new Promise(resolve => setTimeout(resolve, 1000));
  await notificationSystem.stop();
}

// Run all examples
async function runAllExamples() {
  try {
    await basicPubSubExample();
    await eventBusExample();
    await reliableMessagingExample();
    await notificationSystemExample();
  } catch (error) {
    console.error('Example failed:', error);
  }
}

// Uncomment to run examples
// runAllExamples();

export {
  basicPubSubExample,
  eventBusExample,
  reliableMessagingExample,
  notificationSystemExample
};
```

### Chat Application Usage

```typescript
// src/examples/chat-usage.ts
import { ChatApplication } from './chat-application';

async function runChatApplication() {
  const chatApp = new ChatApplication({
    host: 'localhost',
    port: 6379
  });

  await chatApp.start(3000);
  console.log('Chat application started on http://localhost:3000');

  // The application will handle WebSocket connections
  // and Redis pub/sub messaging automatically
}

// runChatApplication();
```

## Performance Considerations

### Message Size Optimization

```typescript
// Compress large messages
import zlib from 'zlib';
import { promisify } from 'util';

const gzip = promisify(zlib.gzip);
const gunzip = promisify(zlib.gunzip);

class CompressedPublisher extends RedisPublisher {
  async publishCompressed(channel: string, message: any): Promise<number> {
    const serialized = JSON.stringify(message);
    
    if (serialized.length > 1024) { // Compress if > 1KB
      const compressed = await gzip(serialized);
      const compressedMessage = {
        _compressed: true,
        data: compressed.toString('base64')
      };
      return await this.publish(channel, compressedMessage);
    }
    
    return await this.publish(channel, message);
  }
}

class CompressedSubscriber extends RedisSubscriber {
  async subscribeCompressed(channel: string, handler: (message: any) => void): Promise<void> {
    await this.subscribe(channel, async (message) => {
      try {
        const parsed = JSON.parse(message);
        
        if (parsed._compressed) {
          const compressed = Buffer.from(parsed.data, 'base64');
          const decompressed = await gunzip(compressed);
          const originalMessage = JSON.parse(decompressed.toString());
          handler(originalMessage);
        } else {
          handler(parsed);
        }
      } catch (error) {
        console.error('Failed to decompress message:', error);
      }
    });
  }
}
```

### Connection Pooling for Pub/Sub

```typescript
// src/pubsub/pubsub-pool.ts
import { RedisPublisher } from './publisher';
import { RedisSubscriber } from './subscriber';

interface PubSubPoolConfig {
  publisherPoolSize: number;
  subscriberPoolSize: number;
  redisConfig: any;
}

class PubSubPool {
  private publishers: RedisPublisher[] = [];
  private subscribers: RedisSubscriber[] = [];
  private publisherIndex: number = 0;
  private subscriberIndex: number = 0;
  private config: PubSubPoolConfig;

  constructor(config: PubSubPoolConfig) {
    this.config = config;
  }

  async initialize(): Promise<void> {
    // Create publisher pool
    for (let i = 0; i < this.config.publisherPoolSize; i++) {
      const publisher = new RedisPublisher(this.config.redisConfig);
      await publisher.connect();
      this.publishers.push(publisher);
    }

    // Create subscriber pool
    for (let i = 0; i < this.config.subscriberPoolSize; i++) {
      const subscriber = new RedisSubscriber(this.config.redisConfig);
      await subscriber.connect();
      this.subscribers.push(subscriber);
    }
  }

  getPublisher(): RedisPublisher {
    const publisher = this.publishers[this.publisherIndex];
    this.publisherIndex = (this.publisherIndex + 1) % this.publishers.length;
    return publisher;
  }

  getSubscriber(): RedisSubscriber {
    const subscriber = this.subscribers[this.subscriberIndex];
    this.subscriberIndex = (this.subscriberIndex + 1) % this.subscribers.length;
    return subscriber;
  }

  async destroy(): Promise<void> {
    const disconnectPromises = [
      ...this.publishers.map(p => p.disconnect()),
      ...this.subscribers.map(s => s.disconnect())
    ];
    
    await Promise.all(disconnectPromises);
  }
}

export { PubSubPool, PubSubPoolConfig };
```

## Best Practices

### 1. Channel Naming Conventions

```typescript
// Good channel naming patterns
const CHANNELS = {
  // Hierarchical naming
  USER_EVENTS: 'user:events',
  USER_LOGIN: 'user:login',
  USER_LOGOUT: 'user:logout',
  
  // Feature-based naming
  ORDERS: 'orders:created',
  PAYMENTS: 'payments:processed',
  NOTIFICATIONS: 'notifications:send',
  
  // Environment-specific
  PROD_ALERTS: 'prod:alerts',
  DEV_DEBUG: 'dev:debug',
  
  // Pattern examples
  ALL_USER_EVENTS: 'user:*',
  ALL_ORDER_EVENTS: 'orders:*'
};
```

### 2. Message Structure Standards

```typescript
interface StandardMessage {
  id: string;
  type: string;
  timestamp: Date;
  version: string;
  source: string;
  data: any;
  metadata?: {
    correlationId?: string;
    userId?: string;
    sessionId?: string;
    traceId?: string;
  };
}

class MessageBuilder {
  static create(type: string, data: any, metadata?: any): StandardMessage {
    return {
      id: `msg_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      type,
      timestamp: new Date(),
      version: '1.0',
      source: 'application',
      data,
      metadata
    };
  }
}
```

## Next Steps

Now that you understand Pub/Sub messaging, you're ready to explore:

- [Part 8: Caching Strategies](./part-08-caching-strategies.md)

## Quick Reference

### Basic Pub/Sub Commands
```typescript
// Publishing
await client.publish('channel', 'message');

// Subscribing
await client.subscribe('channel');
client.on('message', (channel, message) => {
  console.log(`Received: ${message} on ${channel}`);
});

// Pattern subscription
await client.pSubscribe('user:*');
client.on('pmessage', (pattern, channel, message) => {
  console.log(`Pattern ${pattern} matched ${channel}: ${message}`);
});

// Unsubscribing
await client.unsubscribe('channel');
await client.pUnsubscribe('pattern');
```

### Event Bus Pattern
```typescript
const eventBus = new EventBus(config);

// Emit event
await eventBus.emit('user:registered', userData);

// Handle event
eventBus.on('user:registered', async (data) => {
  // Handle user registration
});

// One-time handler
eventBus.once('system:ready', () => {
  console.log('System is ready!');
});
```

### Reliable Messaging
```typescript
const reliable = new ReliableMessaging(config);

// Publish with acknowledgment
const messageId = await reliable.publishReliable('orders', orderData);

// Subscribe with acknowledgment
await reliable.subscribeReliable('orders', async (data, messageId) => {
  // Process order
  return true; // Acknowledge success
});
```

---

**Previous**: [← Part 6: Connection Management](./part-06-connection-management.md)  
**Next**: [Part 8: Caching Strategies →](./part-08-caching-strategies.md)