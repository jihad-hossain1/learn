# Part 11: Session Management with Redis

Redis is an excellent choice for session storage due to its speed, persistence options, and built-in expiration capabilities. This part covers comprehensive session management strategies.

## Table of Contents

1. [Session Management Basics](#session-management-basics)
2. [Session Storage Strategies](#session-storage-strategies)
3. [Express.js Integration](#expressjs-integration)
4. [Advanced Session Features](#advanced-session-features)
5. [Session Security](#session-security)
6. [Multi-Server Sessions](#multi-server-sessions)
7. [Session Analytics](#session-analytics)
8. [Performance Optimization](#performance-optimization)
9. [Real-world Examples](#real-world-examples)

## Session Management Basics

### What are Sessions?

Sessions maintain state between HTTP requests for a specific user, storing:
- User authentication status
- User preferences and settings
- Shopping cart contents
- Temporary form data
- Application state

### Why Redis for Sessions?

- **Fast Access**: In-memory storage for quick retrieval
- **Automatic Expiration**: Built-in TTL support
- **Persistence**: Optional data durability
- **Scalability**: Shared across multiple servers
- **Data Structures**: Rich data types for complex session data

## Session Storage Strategies

### Basic Session Store

```typescript
// src/session/basic-session-store.ts
import { RedisClientType } from 'redis';
import crypto from 'crypto';

interface SessionData {
  [key: string]: any;
}

interface SessionOptions {
  ttl?: number;
  prefix?: string;
  generateId?: () => string;
}

class BasicSessionStore {
  private client: RedisClientType;
  private ttl: number;
  private prefix: string;
  private generateId: () => string;

  constructor(client: RedisClientType, options: SessionOptions = {}) {
    this.client = client;
    this.ttl = options.ttl || 1800; // 30 minutes default
    this.prefix = options.prefix || 'sess:';
    this.generateId = options.generateId || this.defaultGenerateId;
  }

  private defaultGenerateId(): string {
    return crypto.randomBytes(32).toString('hex');
  }

  private getKey(sessionId: string): string {
    return `${this.prefix}${sessionId}`;
  }

  async create(data: SessionData = {}): Promise<string> {
    const sessionId = this.generateId();
    const key = this.getKey(sessionId);
    
    const sessionData = {
      ...data,
      id: sessionId,
      createdAt: Date.now(),
      lastAccessed: Date.now()
    };

    await this.client.setEx(key, this.ttl, JSON.stringify(sessionData));
    return sessionId;
  }

  async get(sessionId: string): Promise<SessionData | null> {
    const key = this.getKey(sessionId);
    const data = await this.client.get(key);
    
    if (!data) {
      return null;
    }

    try {
      const sessionData = JSON.parse(data);
      
      // Update last accessed time
      sessionData.lastAccessed = Date.now();
      await this.client.setEx(key, this.ttl, JSON.stringify(sessionData));
      
      return sessionData;
    } catch (error) {
      console.error('Error parsing session data:', error);
      return null;
    }
  }

  async set(sessionId: string, data: SessionData): Promise<boolean> {
    const key = this.getKey(sessionId);
    const existing = await this.get(sessionId);
    
    if (!existing) {
      return false;
    }

    const sessionData = {
      ...existing,
      ...data,
      lastAccessed: Date.now()
    };

    await this.client.setEx(key, this.ttl, JSON.stringify(sessionData));
    return true;
  }

  async update(sessionId: string, updates: Partial<SessionData>): Promise<boolean> {
    const existing = await this.get(sessionId);
    if (!existing) {
      return false;
    }

    return await this.set(sessionId, { ...existing, ...updates });
  }

  async destroy(sessionId: string): Promise<boolean> {
    const key = this.getKey(sessionId);
    const result = await this.client.del(key);
    return result > 0;
  }

  async touch(sessionId: string): Promise<boolean> {
    const key = this.getKey(sessionId);
    const exists = await this.client.exists(key);
    
    if (exists) {
      await this.client.expire(key, this.ttl);
      return true;
    }
    
    return false;
  }

  async exists(sessionId: string): Promise<boolean> {
    const key = this.getKey(sessionId);
    const exists = await this.client.exists(key);
    return exists > 0;
  }

  async getTTL(sessionId: string): Promise<number> {
    const key = this.getKey(sessionId);
    return await this.client.ttl(key);
  }

  async extendTTL(sessionId: string, additionalSeconds: number): Promise<boolean> {
    const key = this.getKey(sessionId);
    const currentTTL = await this.client.ttl(key);
    
    if (currentTTL > 0) {
      await this.client.expire(key, currentTTL + additionalSeconds);
      return true;
    }
    
    return false;
  }
}

export { BasicSessionStore, SessionData, SessionOptions };
```

### Hash-Based Session Store

Using Redis hashes for more efficient storage and partial updates.

```typescript
// src/session/hash-session-store.ts
class HashSessionStore {
  private client: RedisClientType;
  private ttl: number;
  private prefix: string;
  private generateId: () => string;

  constructor(client: RedisClientType, options: SessionOptions = {}) {
    this.client = client;
    this.ttl = options.ttl || 1800;
    this.prefix = options.prefix || 'sess:';
    this.generateId = options.generateId || this.defaultGenerateId;
  }

  private defaultGenerateId(): string {
    return crypto.randomBytes(32).toString('hex');
  }

  private getKey(sessionId: string): string {
    return `${this.prefix}${sessionId}`;
  }

  async create(data: SessionData = {}): Promise<string> {
    const sessionId = this.generateId();
    const key = this.getKey(sessionId);
    
    const sessionData = {
      ...data,
      id: sessionId,
      createdAt: Date.now().toString(),
      lastAccessed: Date.now().toString()
    };

    // Convert all values to strings for Redis hash
    const hashData: Record<string, string> = {};
    for (const [k, v] of Object.entries(sessionData)) {
      hashData[k] = typeof v === 'object' ? JSON.stringify(v) : String(v);
    }

    await this.client.hMSet(key, hashData);
    await this.client.expire(key, this.ttl);
    
    return sessionId;
  }

  async get(sessionId: string): Promise<SessionData | null> {
    const key = this.getKey(sessionId);
    const data = await this.client.hGetAll(key);
    
    if (Object.keys(data).length === 0) {
      return null;
    }

    // Convert back to appropriate types
    const sessionData: SessionData = {};
    for (const [k, v] of Object.entries(data)) {
      try {
        // Try to parse as JSON first
        sessionData[k] = JSON.parse(v);
      } catch {
        // If not JSON, keep as string
        sessionData[k] = v;
      }
    }

    // Update last accessed time
    await this.client.hSet(key, 'lastAccessed', Date.now().toString());
    await this.client.expire(key, this.ttl);
    
    return sessionData;
  }

  async set(sessionId: string, field: string, value: any): Promise<boolean> {
    const key = this.getKey(sessionId);
    const exists = await this.client.exists(key);
    
    if (!exists) {
      return false;
    }

    const stringValue = typeof value === 'object' ? JSON.stringify(value) : String(value);
    
    await this.client.hSet(key, field, stringValue);
    await this.client.hSet(key, 'lastAccessed', Date.now().toString());
    await this.client.expire(key, this.ttl);
    
    return true;
  }

  async getField(sessionId: string, field: string): Promise<any> {
    const key = this.getKey(sessionId);
    const value = await this.client.hGet(key, field);
    
    if (value === null) {
      return null;
    }

    try {
      return JSON.parse(value);
    } catch {
      return value;
    }
  }

  async deleteField(sessionId: string, field: string): Promise<boolean> {
    const key = this.getKey(sessionId);
    const result = await this.client.hDel(key, field);
    
    if (result > 0) {
      await this.client.hSet(key, 'lastAccessed', Date.now().toString());
      await this.client.expire(key, this.ttl);
    }
    
    return result > 0;
  }

  async updateFields(sessionId: string, updates: Record<string, any>): Promise<boolean> {
    const key = this.getKey(sessionId);
    const exists = await this.client.exists(key);
    
    if (!exists) {
      return false;
    }

    const hashUpdates: Record<string, string> = {
      lastAccessed: Date.now().toString()
    };
    
    for (const [k, v] of Object.entries(updates)) {
      hashUpdates[k] = typeof v === 'object' ? JSON.stringify(v) : String(v);
    }

    await this.client.hMSet(key, hashUpdates);
    await this.client.expire(key, this.ttl);
    
    return true;
  }

  async destroy(sessionId: string): Promise<boolean> {
    const key = this.getKey(sessionId);
    const result = await this.client.del(key);
    return result > 0;
  }

  async getAllFields(sessionId: string): Promise<string[]> {
    const key = this.getKey(sessionId);
    return await this.client.hKeys(key);
  }
}

export { HashSessionStore };
```

## Express.js Integration

### Custom Session Middleware

```typescript
// src/session/express-middleware.ts
import { Request, Response, NextFunction } from 'express';
import { BasicSessionStore } from './basic-session-store';

interface SessionConfig {
  store: BasicSessionStore;
  cookieName?: string;
  cookieOptions?: {
    httpOnly?: boolean;
    secure?: boolean;
    sameSite?: 'strict' | 'lax' | 'none';
    maxAge?: number;
    domain?: string;
    path?: string;
  };
  genid?: () => string;
  resave?: boolean;
  saveUninitialized?: boolean;
}

declare global {
  namespace Express {
    interface Request {
      session?: SessionData & {
        id: string;
        save: () => Promise<void>;
        destroy: () => Promise<void>;
        regenerate: () => Promise<void>;
        touch: () => Promise<void>;
      };
      sessionID?: string;
    }
  }
}

function createSessionMiddleware(config: SessionConfig) {
  const {
    store,
    cookieName = 'sessionId',
    cookieOptions = {
      httpOnly: true,
      secure: false,
      sameSite: 'lax',
      maxAge: 1800000 // 30 minutes
    },
    resave = false,
    saveUninitialized = false
  } = config;

  return async (req: Request, res: Response, next: NextFunction) => {
    let sessionId = req.cookies[cookieName];
    let sessionData: SessionData | null = null;
    let isNew = false;

    // Try to load existing session
    if (sessionId) {
      sessionData = await store.get(sessionId);
      if (!sessionData) {
        // Session expired or invalid
        sessionId = undefined;
      }
    }

    // Create new session if needed
    if (!sessionId) {
      sessionId = await store.create();
      sessionData = await store.get(sessionId);
      isNew = true;
      
      // Set cookie
      res.cookie(cookieName, sessionId, cookieOptions);
    }

    if (!sessionData) {
      return next(new Error('Failed to create session'));
    }

    req.sessionID = sessionId;
    
    // Create session object with methods
    req.session = {
      ...sessionData,
      id: sessionId,
      
      save: async () => {
        if (req.session) {
          const { save, destroy, regenerate, touch, ...data } = req.session;
          await store.set(sessionId!, data);
        }
      },
      
      destroy: async () => {
        await store.destroy(sessionId!);
        res.clearCookie(cookieName);
        req.session = undefined;
        req.sessionID = undefined;
      },
      
      regenerate: async () => {
        const oldData = req.session ? { ...req.session } : {};
        await store.destroy(sessionId!);
        
        const newSessionId = await store.create(oldData);
        const newSessionData = await store.get(newSessionId);
        
        res.cookie(cookieName, newSessionId, cookieOptions);
        req.sessionID = newSessionId;
        
        if (newSessionData) {
          req.session = {
            ...newSessionData,
            id: newSessionId,
            save: req.session!.save,
            destroy: req.session!.destroy,
            regenerate: req.session!.regenerate,
            touch: req.session!.touch
          };
        }
      },
      
      touch: async () => {
        await store.touch(sessionId!);
      }
    };

    // Auto-save session on response end
    const originalEnd = res.end;
    res.end = function(this: Response, ...args: any[]) {
      res.end = originalEnd;
      
      if (req.session && (resave || req.session !== sessionData || isNew)) {
        if (saveUninitialized || !isNew || Object.keys(req.session).length > 4) {
          // Save session (excluding methods)
          const { save, destroy, regenerate, touch, ...data } = req.session;
          store.set(sessionId!, data).catch(console.error);
        }
      }
      
      return originalEnd.apply(this, args);
    };

    next();
  };
}

export { createSessionMiddleware, SessionConfig };
```

### Express App Example

```typescript
// src/session/express-app.ts
import express from 'express';
import cookieParser from 'cookie-parser';
import { createClient } from 'redis';
import { BasicSessionStore } from './basic-session-store';
import { createSessionMiddleware } from './express-middleware';

async function createApp() {
  const app = express();
  
  // Redis client setup
  const redisClient = createClient({
    host: 'localhost',
    port: 6379
  });
  await redisClient.connect();
  
  // Session store
  const sessionStore = new BasicSessionStore(redisClient, {
    ttl: 1800, // 30 minutes
    prefix: 'sess:'
  });
  
  // Middleware
  app.use(express.json());
  app.use(cookieParser());
  app.use(createSessionMiddleware({
    store: sessionStore,
    cookieName: 'sessionId',
    cookieOptions: {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      maxAge: 1800000 // 30 minutes
    },
    resave: false,
    saveUninitialized: false
  }));
  
  // Routes
  app.post('/login', async (req, res) => {
    const { username, password } = req.body;
    
    // Validate credentials (mock)
    if (username === 'admin' && password === 'password') {
      if (req.session) {
        req.session.userId = 'user123';
        req.session.username = username;
        req.session.isAuthenticated = true;
        req.session.loginTime = Date.now();
        
        await req.session.save();
      }
      
      res.json({ success: true, message: 'Logged in successfully' });
    } else {
      res.status(401).json({ success: false, message: 'Invalid credentials' });
    }
  });
  
  app.post('/logout', async (req, res) => {
    if (req.session) {
      await req.session.destroy();
    }
    res.json({ success: true, message: 'Logged out successfully' });
  });
  
  app.get('/profile', (req, res) => {
    if (req.session?.isAuthenticated) {
      res.json({
        userId: req.session.userId,
        username: req.session.username,
        loginTime: req.session.loginTime,
        sessionId: req.session.id
      });
    } else {
      res.status(401).json({ error: 'Not authenticated' });
    }
  });
  
  app.put('/profile', async (req, res) => {
    if (!req.session?.isAuthenticated) {
      return res.status(401).json({ error: 'Not authenticated' });
    }
    
    const { preferences } = req.body;
    
    if (req.session) {
      req.session.preferences = preferences;
      await req.session.save();
    }
    
    res.json({ success: true, message: 'Preferences updated' });
  });
  
  app.get('/session-info', (req, res) => {
    res.json({
      sessionId: req.sessionID,
      session: req.session,
      authenticated: !!req.session?.isAuthenticated
    });
  });
  
  return { app, redisClient };
}

export { createApp };
```

## Advanced Session Features

### Session with User Tracking

```typescript
// src/session/user-session-manager.ts
interface UserSession extends SessionData {
  userId?: string;
  username?: string;
  isAuthenticated?: boolean;
  loginTime?: number;
  lastActivity?: number;
  ipAddress?: string;
  userAgent?: string;
  permissions?: string[];
  preferences?: Record<string, any>;
}

class UserSessionManager {
  private sessionStore: BasicSessionStore;
  private client: RedisClientType;

  constructor(sessionStore: BasicSessionStore, client: RedisClientType) {
    this.sessionStore = sessionStore;
    this.client = client;
  }

  async createUserSession(
    userId: string,
    userData: Partial<UserSession>,
    options: { ipAddress?: string; userAgent?: string } = {}
  ): Promise<string> {
    const sessionData: UserSession = {
      userId,
      isAuthenticated: true,
      loginTime: Date.now(),
      lastActivity: Date.now(),
      ipAddress: options.ipAddress,
      userAgent: options.userAgent,
      ...userData
    };

    const sessionId = await this.sessionStore.create(sessionData);
    
    // Track user sessions
    await this.client.sAdd(`user:${userId}:sessions`, sessionId);
    await this.client.expire(`user:${userId}:sessions`, 86400); // 24 hours
    
    return sessionId;
  }

  async getUserSessions(userId: string): Promise<string[]> {
    return await this.client.sMembers(`user:${userId}:sessions`);
  }

  async destroyUserSession(sessionId: string): Promise<boolean> {
    const session = await this.sessionStore.get(sessionId) as UserSession;
    
    if (session?.userId) {
      await this.client.sRem(`user:${session.userId}:sessions`, sessionId);
    }
    
    return await this.sessionStore.destroy(sessionId);
  }

  async destroyAllUserSessions(userId: string): Promise<number> {
    const sessionIds = await this.getUserSessions(userId);
    
    let destroyedCount = 0;
    for (const sessionId of sessionIds) {
      const destroyed = await this.sessionStore.destroy(sessionId);
      if (destroyed) destroyedCount++;
    }
    
    await this.client.del(`user:${userId}:sessions`);
    return destroyedCount;
  }

  async updateLastActivity(sessionId: string): Promise<boolean> {
    return await this.sessionStore.update(sessionId, {
      lastActivity: Date.now()
    });
  }

  async getActiveUsers(): Promise<string[]> {
    const pattern = 'user:*:sessions';
    const keys = await this.client.keys(pattern);
    
    return keys.map(key => {
      const match = key.match(/user:(.+):sessions/);
      return match ? match[1] : '';
    }).filter(Boolean);
  }

  async getSessionsByUser(userId: string): Promise<UserSession[]> {
    const sessionIds = await this.getUserSessions(userId);
    const sessions: UserSession[] = [];
    
    for (const sessionId of sessionIds) {
      const session = await this.sessionStore.get(sessionId) as UserSession;
      if (session) {
        sessions.push(session);
      }
    }
    
    return sessions;
  }

  async cleanupExpiredSessions(): Promise<number> {
    const userSessionKeys = await this.client.keys('user:*:sessions');
    let cleanedCount = 0;
    
    for (const userKey of userSessionKeys) {
      const sessionIds = await this.client.sMembers(userKey);
      
      for (const sessionId of sessionIds) {
        const exists = await this.sessionStore.exists(sessionId);
        if (!exists) {
          await this.client.sRem(userKey, sessionId);
          cleanedCount++;
        }
      }
    }
    
    return cleanedCount;
  }
}

export { UserSessionManager, UserSession };
```

### Session Analytics

```typescript
// src/session/session-analytics.ts
interface SessionMetrics {
  totalSessions: number;
  activeSessions: number;
  averageSessionDuration: number;
  uniqueUsers: number;
  sessionsPerUser: number;
  topUserAgents: Array<{ userAgent: string; count: number }>;
  topIpAddresses: Array<{ ip: string; count: number }>;
}

class SessionAnalytics {
  private client: RedisClientType;
  private sessionStore: BasicSessionStore;

  constructor(client: RedisClientType, sessionStore: BasicSessionStore) {
    this.client = client;
    this.sessionStore = sessionStore;
  }

  async trackSessionEvent(
    sessionId: string,
    event: string,
    data?: Record<string, any>
  ): Promise<void> {
    const eventData = {
      sessionId,
      event,
      timestamp: Date.now(),
      ...data
    };

    // Store in a time-series like structure
    const key = `session:events:${sessionId}`;
    await this.client.lPush(key, JSON.stringify(eventData));
    await this.client.expire(key, 86400); // 24 hours
    
    // Update global counters
    await this.client.incr(`session:events:${event}:count`);
    await this.client.expire(`session:events:${event}:count`, 86400);
  }

  async getSessionEvents(sessionId: string): Promise<any[]> {
    const key = `session:events:${sessionId}`;
    const events = await this.client.lRange(key, 0, -1);
    
    return events.map(event => {
      try {
        return JSON.parse(event);
      } catch {
        return null;
      }
    }).filter(Boolean);
  }

  async getSessionMetrics(): Promise<SessionMetrics> {
    // Get all session keys
    const sessionKeys = await this.client.keys('sess:*');
    const totalSessions = sessionKeys.length;
    
    // Get active sessions (those accessed in last hour)
    const oneHourAgo = Date.now() - 3600000;
    let activeSessions = 0;
    let totalDuration = 0;
    const uniqueUsers = new Set<string>();
    const userAgents = new Map<string, number>();
    const ipAddresses = new Map<string, number>();
    
    for (const key of sessionKeys) {
      const sessionData = await this.client.get(key);
      if (sessionData) {
        try {
          const session = JSON.parse(sessionData) as UserSession;
          
          if (session.lastActivity && session.lastActivity > oneHourAgo) {
            activeSessions++;
          }
          
          if (session.createdAt && session.lastActivity) {
            totalDuration += session.lastActivity - session.createdAt;
          }
          
          if (session.userId) {
            uniqueUsers.add(session.userId);
          }
          
          if (session.userAgent) {
            userAgents.set(session.userAgent, (userAgents.get(session.userAgent) || 0) + 1);
          }
          
          if (session.ipAddress) {
            ipAddresses.set(session.ipAddress, (ipAddresses.get(session.ipAddress) || 0) + 1);
          }
        } catch (error) {
          console.error('Error parsing session data:', error);
        }
      }
    }
    
    const averageSessionDuration = totalSessions > 0 ? totalDuration / totalSessions : 0;
    const sessionsPerUser = uniqueUsers.size > 0 ? totalSessions / uniqueUsers.size : 0;
    
    // Sort and get top user agents and IPs
    const topUserAgents = Array.from(userAgents.entries())
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
      .map(([userAgent, count]) => ({ userAgent, count }));
    
    const topIpAddresses = Array.from(ipAddresses.entries())
      .sort((a, b) => b[1] - a[1])
      .slice(0, 10)
      .map(([ip, count]) => ({ ip, count }));
    
    return {
      totalSessions,
      activeSessions,
      averageSessionDuration,
      uniqueUsers: uniqueUsers.size,
      sessionsPerUser,
      topUserAgents,
      topIpAddresses
    };
  }

  async getEventCounts(): Promise<Record<string, number>> {
    const eventKeys = await this.client.keys('session:events:*:count');
    const counts: Record<string, number> = {};
    
    for (const key of eventKeys) {
      const match = key.match(/session:events:(.+):count/);
      if (match) {
        const event = match[1];
        const count = await this.client.get(key);
        counts[event] = parseInt(count || '0');
      }
    }
    
    return counts;
  }

  async generateDailyReport(): Promise<{
    date: string;
    metrics: SessionMetrics;
    events: Record<string, number>;
  }> {
    const metrics = await this.getSessionMetrics();
    const events = await this.getEventCounts();
    
    return {
      date: new Date().toISOString().split('T')[0],
      metrics,
      events
    };
  }
}

export { SessionAnalytics, SessionMetrics };
```

## Session Security

### Secure Session Implementation

```typescript
// src/session/secure-session.ts
import crypto from 'crypto';
import { Request } from 'express';

interface SecurityOptions {
  maxSessionsPerUser?: number;
  sessionTimeout?: number;
  ipValidation?: boolean;
  userAgentValidation?: boolean;
  regenerateOnAuth?: boolean;
  encryptionKey?: string;
}

class SecureSessionManager {
  private sessionStore: BasicSessionStore;
  private userSessionManager: UserSessionManager;
  private options: SecurityOptions;
  private encryptionKey: Buffer;

  constructor(
    sessionStore: BasicSessionStore,
    userSessionManager: UserSessionManager,
    options: SecurityOptions = {}
  ) {
    this.sessionStore = sessionStore;
    this.userSessionManager = userSessionManager;
    this.options = {
      maxSessionsPerUser: 5,
      sessionTimeout: 1800,
      ipValidation: true,
      userAgentValidation: true,
      regenerateOnAuth: true,
      ...options
    };
    
    this.encryptionKey = options.encryptionKey 
      ? Buffer.from(options.encryptionKey, 'hex')
      : crypto.randomBytes(32);
  }

  private encrypt(text: string): string {
    const iv = crypto.randomBytes(16);
    const cipher = crypto.createCipher('aes-256-cbc', this.encryptionKey);
    let encrypted = cipher.update(text, 'utf8', 'hex');
    encrypted += cipher.final('hex');
    return iv.toString('hex') + ':' + encrypted;
  }

  private decrypt(encryptedText: string): string {
    const parts = encryptedText.split(':');
    const iv = Buffer.from(parts[0], 'hex');
    const encrypted = parts[1];
    const decipher = crypto.createDecipher('aes-256-cbc', this.encryptionKey);
    let decrypted = decipher.update(encrypted, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    return decrypted;
  }

  async createSecureSession(
    userId: string,
    userData: Partial<UserSession>,
    req: Request
  ): Promise<string> {
    // Check session limits
    if (this.options.maxSessionsPerUser) {
      const existingSessions = await this.userSessionManager.getUserSessions(userId);
      if (existingSessions.length >= this.options.maxSessionsPerUser) {
        // Remove oldest session
        const oldestSessionId = existingSessions[0];
        await this.userSessionManager.destroyUserSession(oldestSessionId);
      }
    }

    // Create session with security metadata
    const secureData: UserSession = {
      ...userData,
      userId,
      ipAddress: this.getClientIP(req),
      userAgent: req.get('User-Agent'),
      fingerprint: this.generateFingerprint(req),
      isAuthenticated: true,
      securityLevel: 'high'
    };

    return await this.userSessionManager.createUserSession(userId, secureData, {
      ipAddress: secureData.ipAddress,
      userAgent: secureData.userAgent
    });
  }

  async validateSession(sessionId: string, req: Request): Promise<{
    valid: boolean;
    reason?: string;
    session?: UserSession;
  }> {
    const session = await this.sessionStore.get(sessionId) as UserSession;
    
    if (!session) {
      return { valid: false, reason: 'Session not found' };
    }

    // Check IP validation
    if (this.options.ipValidation && session.ipAddress) {
      const currentIP = this.getClientIP(req);
      if (currentIP !== session.ipAddress) {
        return { valid: false, reason: 'IP address mismatch' };
      }
    }

    // Check User-Agent validation
    if (this.options.userAgentValidation && session.userAgent) {
      const currentUserAgent = req.get('User-Agent');
      if (currentUserAgent !== session.userAgent) {
        return { valid: false, reason: 'User agent mismatch' };
      }
    }

    // Check session timeout
    if (session.lastActivity) {
      const timeSinceLastActivity = Date.now() - session.lastActivity;
      if (timeSinceLastActivity > (this.options.sessionTimeout! * 1000)) {
        await this.userSessionManager.destroyUserSession(sessionId);
        return { valid: false, reason: 'Session timeout' };
      }
    }

    // Check fingerprint
    const currentFingerprint = this.generateFingerprint(req);
    if (session.fingerprint && session.fingerprint !== currentFingerprint) {
      return { valid: false, reason: 'Browser fingerprint mismatch' };
    }

    return { valid: true, session };
  }

  private getClientIP(req: Request): string {
    return (req.headers['x-forwarded-for'] as string)?.split(',')[0] ||
           req.connection.remoteAddress ||
           req.socket.remoteAddress ||
           '';
  }

  private generateFingerprint(req: Request): string {
    const components = [
      req.get('User-Agent') || '',
      req.get('Accept-Language') || '',
      req.get('Accept-Encoding') || '',
      req.get('Accept') || ''
    ];
    
    return crypto.createHash('sha256')
      .update(components.join('|'))
      .digest('hex');
  }

  async regenerateSession(oldSessionId: string, req: Request): Promise<string | null> {
    const oldSession = await this.sessionStore.get(oldSessionId) as UserSession;
    
    if (!oldSession || !oldSession.userId) {
      return null;
    }

    // Destroy old session
    await this.userSessionManager.destroyUserSession(oldSessionId);
    
    // Create new session with same data
    return await this.createSecureSession(oldSession.userId, oldSession, req);
  }

  async detectSuspiciousActivity(userId: string): Promise<{
    suspicious: boolean;
    reasons: string[];
  }> {
    const sessions = await this.userSessionManager.getSessionsByUser(userId);
    const reasons: string[] = [];
    
    // Check for multiple IPs
    const ips = new Set(sessions.map(s => s.ipAddress).filter(Boolean));
    if (ips.size > 3) {
      reasons.push('Multiple IP addresses detected');
    }
    
    // Check for multiple user agents
    const userAgents = new Set(sessions.map(s => s.userAgent).filter(Boolean));
    if (userAgents.size > 2) {
      reasons.push('Multiple browsers/devices detected');
    }
    
    // Check for rapid session creation
    const recentSessions = sessions.filter(s => 
      s.loginTime && (Date.now() - s.loginTime) < 300000 // 5 minutes
    );
    if (recentSessions.length > 3) {
      reasons.push('Rapid session creation detected');
    }
    
    return {
      suspicious: reasons.length > 0,
      reasons
    };
  }
}

export { SecureSessionManager, SecurityOptions };
```

## Usage Examples

```typescript
// src/examples/session-usage.ts
import { createClient } from 'redis';
import { BasicSessionStore } from '../session/basic-session-store';
import { HashSessionStore } from '../session/hash-session-store';
import { UserSessionManager } from '../session/user-session-manager';
import { SessionAnalytics } from '../session/session-analytics';
import { SecureSessionManager } from '../session/secure-session';

async function runSessionExamples() {
  const client = createClient({
    host: 'localhost',
    port: 6379
  });
  
  await client.connect();
  
  try {
    // Basic session store example
    console.log('=== Basic Session Store ===');
    const basicStore = new BasicSessionStore(client, {
      ttl: 1800,
      prefix: 'sess:'
    });
    
    const sessionId = await basicStore.create({
      userId: 'user123',
      username: 'john_doe',
      isAuthenticated: true
    });
    console.log('Created session:', sessionId);
    
    const session = await basicStore.get(sessionId);
    console.log('Retrieved session:', session);
    
    await basicStore.update(sessionId, {
      lastPage: '/dashboard',
      preferences: { theme: 'dark' }
    });
    
    // Hash session store example
    console.log('\n=== Hash Session Store ===');
    const hashStore = new HashSessionStore(client, {
      ttl: 1800,
      prefix: 'hsess:'
    });
    
    const hashSessionId = await hashStore.create({
      userId: 'user456',
      username: 'jane_doe'
    });
    
    await hashStore.set(hashSessionId, 'cart', [
      { productId: 'p1', quantity: 2 },
      { productId: 'p2', quantity: 1 }
    ]);
    
    const cart = await hashStore.getField(hashSessionId, 'cart');
    console.log('Shopping cart:', cart);
    
    // User session manager example
    console.log('\n=== User Session Manager ===');
    const userManager = new UserSessionManager(basicStore, client);
    
    const userSessionId = await userManager.createUserSession('user789', {
      username: 'admin',
      permissions: ['read', 'write', 'admin']
    }, {
      ipAddress: '192.168.1.100',
      userAgent: 'Mozilla/5.0...'
    });
    
    const userSessions = await userManager.getUserSessions('user789');
    console.log('User sessions:', userSessions);
    
    // Session analytics example
    console.log('\n=== Session Analytics ===');
    const analytics = new SessionAnalytics(client, basicStore);
    
    await analytics.trackSessionEvent(sessionId, 'page_view', {
      page: '/dashboard',
      duration: 5000
    });
    
    await analytics.trackSessionEvent(sessionId, 'button_click', {
      button: 'save_profile',
      location: 'header'
    });
    
    const events = await analytics.getSessionEvents(sessionId);
    console.log('Session events:', events);
    
    const metrics = await analytics.getSessionMetrics();
    console.log('Session metrics:', metrics);
    
    // Secure session manager example
    console.log('\n=== Secure Session Manager ===');
    const secureManager = new SecureSessionManager(basicStore, userManager, {
      maxSessionsPerUser: 3,
      sessionTimeout: 1800,
      ipValidation: true,
      userAgentValidation: true
    });
    
    // Mock request object
    const mockReq = {
      headers: {
        'user-agent': 'Mozilla/5.0 (Test Browser)',
        'x-forwarded-for': '192.168.1.200'
      },
      get: (header: string) => mockReq.headers[header.toLowerCase()],
      connection: { remoteAddress: '192.168.1.200' },
      socket: { remoteAddress: '192.168.1.200' }
    } as any;
    
    const secureSessionId = await secureManager.createSecureSession(
      'secure_user',
      { username: 'secure_admin' },
      mockReq
    );
    
    const validation = await secureManager.validateSession(secureSessionId, mockReq);
    console.log('Session validation:', validation);
    
    const suspiciousActivity = await secureManager.detectSuspiciousActivity('secure_user');
    console.log('Suspicious activity check:', suspiciousActivity);
    
  } finally {
    await client.disconnect();
  }
}

// Run examples
if (require.main === module) {
  runSessionExamples().catch(console.error);
}

export { runSessionExamples };
```

## Best Practices

### 1. Security

- Use HTTPS in production
- Set secure cookie flags
- Implement session regeneration on authentication
- Validate session integrity
- Monitor for suspicious activity

### 2. Performance

- Use appropriate TTL values
- Implement session cleanup
- Consider using Redis hashes for large sessions
- Monitor memory usage

### 3. Scalability

- Use Redis clustering for high availability
- Implement session replication
- Consider session affinity vs. stateless design
- Monitor session distribution

### 4. Monitoring

- Track session metrics
- Monitor session creation/destruction rates
- Alert on unusual patterns
- Log security events

## Next Steps

Now that you understand session management, you're ready to explore:

- [Part 12: Rate Limiting and Security](./part-12-rate-limiting.md)

## Quick Reference

### Basic Session Operations
```typescript
// Create session
const sessionId = await store.create({ userId: 'user123' });

// Get session
const session = await store.get(sessionId);

// Update session
await store.update(sessionId, { lastPage: '/dashboard' });

// Destroy session
await store.destroy(sessionId);

// Touch session (extend TTL)
await store.touch(sessionId);
```

### Express Integration
```typescript
// Session middleware
app.use(createSessionMiddleware({
  store: sessionStore,
  cookieName: 'sessionId',
  cookieOptions: {
    httpOnly: true,
    secure: true,
    sameSite: 'strict'
  }
}));

// Using session in routes
app.get('/profile', (req, res) => {
  if (req.session?.isAuthenticated) {
    res.json({ user: req.session.username });
  } else {
    res.status(401).json({ error: 'Not authenticated' });
  }
});
```

---

**Previous**: [← Part 10: Caching Strategies](./part-10-caching-strategies.md)  
**Next**: [Part 12: Rate Limiting and Security →](./part-12-rate-limiting.md)