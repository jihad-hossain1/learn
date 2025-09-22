# Performance and Optimization for Full Stack Development

## Table of Contents
1. [Performance Fundamentals](#performance-fundamentals)
2. [Caching Strategies](#caching-strategies)
3. [Content Delivery Networks (CDN)](#content-delivery-networks-cdn)
4. [Load Balancing](#load-balancing)
5. [Database Optimization](#database-optimization)
6. [Frontend Performance](#frontend-performance)
7. [Monitoring and Metrics](#monitoring-and-metrics)
8. [Scalability Patterns](#scalability-patterns)

## Performance Fundamentals

### Key Performance Metrics

#### Frontend Metrics
- **First Contentful Paint (FCP)**: Time to first content render
- **Largest Contentful Paint (LCP)**: Time to largest content render
- **First Input Delay (FID)**: Time to first user interaction
- **Cumulative Layout Shift (CLS)**: Visual stability measure
- **Time to Interactive (TTI)**: Time until page is fully interactive

#### Backend Metrics
- **Response Time**: Time to process and return response
- **Throughput**: Requests processed per second
- **Error Rate**: Percentage of failed requests
- **Availability**: Uptime percentage
- **Resource Utilization**: CPU, memory, disk usage

### Performance Budget
```javascript
// Example performance budget
const performanceBudget = {
    // Network
    maxResponseTime: 200, // milliseconds
    maxThroughput: 1000,  // requests per second
    
    // Frontend
    maxBundleSize: 250,   // KB
    maxImageSize: 100,    // KB
    maxFCP: 1500,         // milliseconds
    maxLCP: 2500,         // milliseconds
    maxFID: 100,          // milliseconds
    maxCLS: 0.1,          // score
    
    // Backend
    maxDatabaseQueryTime: 50, // milliseconds
    maxAPIResponseTime: 100,  // milliseconds
};
```

## Caching Strategies

### 1. Browser Caching

#### Cache-Control Headers
```javascript
// Static assets (images, CSS, JS)
app.use('/static', express.static('public', {
    maxAge: '1y', // 1 year
    etag: true,
    lastModified: true
}));

// API responses
app.get('/api/products', (req, res) => {
    res.set({
        'Cache-Control': 'public, max-age=300', // 5 minutes
        'ETag': generateETag(data),
        'Last-Modified': new Date().toUTCString()
    });
    res.json(products);
});

// No cache for sensitive data
app.get('/api/user/profile', (req, res) => {
    res.set({
        'Cache-Control': 'no-store, no-cache, must-revalidate',
        'Pragma': 'no-cache',
        'Expires': '0'
    });
    res.json(userProfile);
});
```

#### ETag Implementation
```javascript
const crypto = require('crypto');

const generateETag = (data) => {
    return crypto.createHash('md5')
        .update(JSON.stringify(data))
        .digest('hex');
};

const checkETag = (req, res, data) => {
    const etag = generateETag(data);
    const clientETag = req.headers['if-none-match'];
    
    if (clientETag === etag) {
        res.status(304).end(); // Not Modified
        return true;
    }
    
    res.set('ETag', etag);
    return false;
};

app.get('/api/data', (req, res) => {
    const data = fetchData();
    
    if (checkETag(req, res, data)) {
        return; // 304 response sent
    }
    
    res.json(data);
});
```

### 2. Server-Side Caching

#### In-Memory Caching with Node.js
```javascript
const NodeCache = require('node-cache');
const cache = new NodeCache({ 
    stdTTL: 600, // 10 minutes default
    checkperiod: 120 // Check for expired keys every 2 minutes
});

// Cache middleware
const cacheMiddleware = (duration = 300) => {
    return (req, res, next) => {
        const key = req.originalUrl;
        const cachedResponse = cache.get(key);
        
        if (cachedResponse) {
            console.log('Cache hit:', key);
            return res.json(cachedResponse);
        }
        
        // Override res.json to cache the response
        const originalJson = res.json;
        res.json = function(data) {
            cache.set(key, data, duration);
            console.log('Cache set:', key);
            originalJson.call(this, data);
        };
        
        next();
    };
};

// Usage
app.get('/api/products', cacheMiddleware(600), async (req, res) => {
    const products = await Product.findAll();
    res.json(products);
});
```

#### Redis Caching
```javascript
const redis = require('redis');
const client = redis.createClient({
    host: process.env.REDIS_HOST,
    port: process.env.REDIS_PORT,
    password: process.env.REDIS_PASSWORD
});

class RedisCache {
    async get(key) {
        try {
            const data = await client.get(key);
            return data ? JSON.parse(data) : null;
        } catch (error) {
            console.error('Redis get error:', error);
            return null;
        }
    }
    
    async set(key, data, ttl = 3600) {
        try {
            await client.setex(key, ttl, JSON.stringify(data));
        } catch (error) {
            console.error('Redis set error:', error);
        }
    }
    
    async del(key) {
        try {
            await client.del(key);
        } catch (error) {
            console.error('Redis del error:', error);
        }
    }
    
    async flush() {
        try {
            await client.flushall();
        } catch (error) {
            console.error('Redis flush error:', error);
        }
    }
}

const redisCache = new RedisCache();

// Redis cache middleware
const redisCacheMiddleware = (ttl = 3600) => {
    return async (req, res, next) => {
        const key = `cache:${req.originalUrl}`;
        
        try {
            const cachedData = await redisCache.get(key);
            
            if (cachedData) {
                console.log('Redis cache hit:', key);
                return res.json(cachedData);
            }
            
            // Override res.json to cache the response
            const originalJson = res.json;
            res.json = function(data) {
                redisCache.set(key, data, ttl);
                console.log('Redis cache set:', key);
                originalJson.call(this, data);
            };
            
            next();
        } catch (error) {
            console.error('Cache middleware error:', error);
            next();
        }
    };
};
```

### 3. Database Query Caching
```javascript
// Sequelize with Redis caching
const { Sequelize, DataTypes } = require('sequelize');

class CachedModel {
    constructor(model, cache) {
        this.model = model;
        this.cache = cache;
    }
    
    async findById(id, options = {}) {
        const cacheKey = `${this.model.name}:${id}`;
        
        // Try cache first
        let result = await this.cache.get(cacheKey);
        
        if (!result) {
            result = await this.model.findByPk(id, options);
            if (result) {
                await this.cache.set(cacheKey, result, 3600); // 1 hour
            }
        }
        
        return result;
    }
    
    async findAll(options = {}) {
        const cacheKey = `${this.model.name}:all:${JSON.stringify(options)}`;
        
        let results = await this.cache.get(cacheKey);
        
        if (!results) {
            results = await this.model.findAll(options);
            await this.cache.set(cacheKey, results, 1800); // 30 minutes
        }
        
        return results;
    }
    
    async create(data) {
        const result = await this.model.create(data);
        
        // Invalidate related caches
        await this.invalidateCache();
        
        return result;
    }
    
    async update(id, data) {
        const result = await this.model.update(data, { where: { id } });
        
        // Invalidate specific and related caches
        await this.cache.del(`${this.model.name}:${id}`);
        await this.invalidateCache();
        
        return result;
    }
    
    async invalidateCache() {
        // Invalidate all cached queries for this model
        const pattern = `${this.model.name}:*`;
        // Implementation depends on your cache system
    }
}
```

### 4. Cache Invalidation Strategies

#### Time-based Invalidation (TTL)
```javascript
// Set TTL based on data volatility
const cacheTTL = {
    userProfile: 300,      // 5 minutes (frequently updated)
    productCatalog: 3600,  // 1 hour (moderately updated)
    staticContent: 86400,  // 24 hours (rarely updated)
    configuration: 604800  // 1 week (very rarely updated)
};
```

#### Event-based Invalidation
```javascript
const EventEmitter = require('events');
const cacheInvalidator = new EventEmitter();

// Listen for data changes
cacheInvalidator.on('user:updated', async (userId) => {
    await redisCache.del(`user:${userId}`);
    await redisCache.del(`user:${userId}:profile`);
    await redisCache.del(`user:${userId}:posts`);
});

cacheInvalidator.on('product:updated', async (productId) => {
    await redisCache.del(`product:${productId}`);
    await redisCache.del('products:all');
    await redisCache.del('products:featured');
});

// Emit events when data changes
app.put('/api/users/:id', async (req, res) => {
    const userId = req.params.id;
    const updatedUser = await User.update(req.body, { where: { id: userId } });
    
    // Invalidate cache
    cacheInvalidator.emit('user:updated', userId);
    
    res.json(updatedUser);
});
```

## Content Delivery Networks (CDN)

### CDN Benefits
- **Reduced Latency**: Content served from nearest edge server
- **Improved Availability**: Distributed infrastructure
- **Reduced Bandwidth**: Offload traffic from origin server
- **DDoS Protection**: Built-in attack mitigation
- **Global Reach**: Worldwide presence

### CDN Implementation

#### Static Asset CDN
```javascript
// Environment-based CDN configuration
const config = {
    development: {
        cdnUrl: 'http://localhost:3000',
        usecdn: false
    },
    production: {
        cdnUrl: 'https://cdn.example.com',
        usecdn: true
    }
};

const getCDNUrl = (path) => {
    const env = process.env.NODE_ENV || 'development';
    const baseUrl = config[env].usecdn ? config[env].cdnUrl : '';
    return `${baseUrl}${path}`;
};

// Usage in templates
app.locals.cdnUrl = getCDNUrl;

// In HTML template
<link rel="stylesheet" href="<%= cdnUrl('/css/styles.css') %>">
<script src="<%= cdnUrl('/js/app.js') %>"></script>
<img src="<%= cdnUrl('/images/logo.png') %>" alt="Logo">
```

#### API Response CDN
```javascript
// CloudFlare CDN headers
app.get('/api/public-data', (req, res) => {
    res.set({
        'Cache-Control': 'public, max-age=3600, s-maxage=7200',
        'CDN-Cache-Control': 'max-age=7200',
        'Cloudflare-CDN-Cache-Control': 'max-age=7200'
    });
    
    res.json(publicData);
});

// Vary header for different responses
app.get('/api/content', (req, res) => {
    res.set({
        'Vary': 'Accept-Language, Accept-Encoding',
        'Cache-Control': 'public, max-age=1800'
    });
    
    const language = req.headers['accept-language'];
    const content = getLocalizedContent(language);
    
    res.json(content);
});
```

### CDN Cache Purging
```javascript
// CloudFlare API cache purging
const cloudflare = require('cloudflare')({
    email: process.env.CLOUDFLARE_EMAIL,
    key: process.env.CLOUDFLARE_API_KEY
});

const purgeCache = async (urls) => {
    try {
        await cloudflare.zones.purgeCache(process.env.CLOUDFLARE_ZONE_ID, {
            files: urls
        });
        console.log('Cache purged for:', urls);
    } catch (error) {
        console.error('Cache purge failed:', error);
    }
};

// Purge cache when content updates
app.put('/api/content/:id', async (req, res) => {
    const contentId = req.params.id;
    const updatedContent = await Content.update(req.body, { where: { id: contentId } });
    
    // Purge related CDN cache
    await purgeCache([
        `https://example.com/api/content/${contentId}`,
        'https://example.com/api/content',
        'https://example.com/'
    ]);
    
    res.json(updatedContent);
});
```

## Load Balancing

### Load Balancing Algorithms

#### 1. Round Robin
```javascript
class RoundRobinBalancer {
    constructor(servers) {
        this.servers = servers;
        this.currentIndex = 0;
    }
    
    getNextServer() {
        const server = this.servers[this.currentIndex];
        this.currentIndex = (this.currentIndex + 1) % this.servers.length;
        return server;
    }
}

const servers = [
    { host: 'server1.example.com', port: 3000 },
    { host: 'server2.example.com', port: 3000 },
    { host: 'server3.example.com', port: 3000 }
];

const balancer = new RoundRobinBalancer(servers);
```

#### 2. Weighted Round Robin
```javascript
class WeightedRoundRobinBalancer {
    constructor(servers) {
        this.servers = servers.map(server => ({
            ...server,
            currentWeight: 0
        }));
        this.totalWeight = servers.reduce((sum, server) => sum + server.weight, 0);
    }
    
    getNextServer() {
        let selected = null;
        let maxWeight = -1;
        
        for (let server of this.servers) {
            server.currentWeight += server.weight;
            
            if (server.currentWeight > maxWeight) {
                maxWeight = server.currentWeight;
                selected = server;
            }
        }
        
        if (selected) {
            selected.currentWeight -= this.totalWeight;
        }
        
        return selected;
    }
}

const weightedServers = [
    { host: 'server1.example.com', port: 3000, weight: 5 },
    { host: 'server2.example.com', port: 3000, weight: 3 },
    { host: 'server3.example.com', port: 3000, weight: 2 }
];
```

#### 3. Least Connections
```javascript
class LeastConnectionsBalancer {
    constructor(servers) {
        this.servers = servers.map(server => ({
            ...server,
            activeConnections: 0
        }));
    }
    
    getNextServer() {
        return this.servers.reduce((min, server) => 
            server.activeConnections < min.activeConnections ? server : min
        );
    }
    
    incrementConnections(server) {
        server.activeConnections++;
    }
    
    decrementConnections(server) {
        server.activeConnections = Math.max(0, server.activeConnections - 1);
    }
}
```

### Nginx Load Balancer Configuration
```nginx
# /etc/nginx/nginx.conf
upstream backend {
    # Load balancing method
    least_conn;  # or ip_hash, hash, random
    
    # Backend servers
    server 10.0.1.10:3000 weight=3 max_fails=3 fail_timeout=30s;
    server 10.0.1.11:3000 weight=2 max_fails=3 fail_timeout=30s;
    server 10.0.1.12:3000 weight=1 max_fails=3 fail_timeout=30s backup;
}

server {
    listen 80;
    server_name example.com;
    
    location / {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # Health checks
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
        
        # Connection pooling
        proxy_http_version 1.1;
        proxy_set_header Connection "";
    }
    
    # Health check endpoint
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
```

### Application-Level Load Balancing
```javascript
const http = require('http');
const httpProxy = require('http-proxy-middleware');
const express = require('express');

const app = express();

// Health check for servers
const checkServerHealth = async (server) => {
    try {
        const response = await fetch(`http://${server.host}:${server.port}/health`);
        return response.ok;
    } catch (error) {
        return false;
    }
};

// Proxy middleware with load balancing
const createLoadBalancedProxy = (servers) => {
    const balancer = new RoundRobinBalancer(servers);
    
    return httpProxy({
        target: 'http://placeholder', // Will be overridden
        changeOrigin: true,
        router: async (req) => {
            // Filter healthy servers
            const healthyServers = [];
            for (let server of servers) {
                if (await checkServerHealth(server)) {
                    healthyServers.push(server);
                }
            }
            
            if (healthyServers.length === 0) {
                throw new Error('No healthy servers available');
            }
            
            const server = balancer.getNextServer();
            return `http://${server.host}:${server.port}`;
        },
        onError: (err, req, res) => {
            console.error('Proxy error:', err);
            res.status(503).json({ error: 'Service temporarily unavailable' });
        }
    });
};

app.use('/api', createLoadBalancedProxy(servers));
```

## Database Optimization

### Query Optimization

#### Indexing Strategies
```sql
-- Single column index
CREATE INDEX idx_user_email ON users(email);

-- Composite index
CREATE INDEX idx_order_user_date ON orders(user_id, created_at);

-- Partial index
CREATE INDEX idx_active_users ON users(email) WHERE active = true;

-- Covering index
CREATE INDEX idx_user_profile ON users(id) INCLUDE (name, email, created_at);
```

#### Query Analysis
```javascript
// Sequelize query optimization
const { QueryTypes } = require('sequelize');

// Bad: N+1 query problem
const users = await User.findAll();
for (let user of users) {
    user.posts = await Post.findAll({ where: { userId: user.id } });
}

// Good: Use includes/joins
const users = await User.findAll({
    include: [{
        model: Post,
        as: 'posts'
    }]
});

// Raw query with explain
const explainQuery = async (query, replacements = {}) => {
    const explainResult = await sequelize.query(
        `EXPLAIN ANALYZE ${query}`,
        {
            replacements,
            type: QueryTypes.SELECT
        }
    );
    console.log('Query execution plan:', explainResult);
};
```

### Connection Pooling
```javascript
// Sequelize connection pool configuration
const sequelize = new Sequelize(database, username, password, {
    host: 'localhost',
    dialect: 'postgres',
    pool: {
        max: 20,          // Maximum connections
        min: 5,           // Minimum connections
        acquire: 30000,   // Maximum time to get connection (ms)
        idle: 10000,      // Maximum idle time (ms)
        evict: 1000,      // Check for idle connections interval
        handleDisconnects: true
    },
    logging: process.env.NODE_ENV === 'development' ? console.log : false
});

// MongoDB connection pool
const { MongoClient } = require('mongodb');

const client = new MongoClient(uri, {
    maxPoolSize: 20,
    minPoolSize: 5,
    maxIdleTimeMS: 30000,
    serverSelectionTimeoutMS: 5000,
    socketTimeoutMS: 45000,
});
```

### Database Sharding
```javascript
// Simple sharding implementation
class DatabaseSharding {
    constructor(shards) {
        this.shards = shards;
    }
    
    getShardKey(userId) {
        // Hash-based sharding
        const hash = this.hashFunction(userId.toString());
        return hash % this.shards.length;
    }
    
    hashFunction(str) {
        let hash = 0;
        for (let i = 0; i < str.length; i++) {
            const char = str.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash; // Convert to 32-bit integer
        }
        return Math.abs(hash);
    }
    
    getShard(userId) {
        const shardIndex = this.getShardKey(userId);
        return this.shards[shardIndex];
    }
    
    async query(userId, query, params) {
        const shard = this.getShard(userId);
        return await shard.query(query, params);
    }
}

const shards = [
    new Sequelize(/* shard 1 config */),
    new Sequelize(/* shard 2 config */),
    new Sequelize(/* shard 3 config */)
];

const sharding = new DatabaseSharding(shards);
```

## Frontend Performance

### Code Splitting
```javascript
// React lazy loading
import React, { Suspense, lazy } from 'react';

const Dashboard = lazy(() => import('./components/Dashboard'));
const Profile = lazy(() => import('./components/Profile'));
const Settings = lazy(() => import('./components/Settings'));

function App() {
    return (
        <Router>
            <Suspense fallback={<div>Loading...</div>}>
                <Routes>
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/profile" element={<Profile />} />
                    <Route path="/settings" element={<Settings />} />
                </Routes>
            </Suspense>
        </Router>
    );
}

// Webpack code splitting
const loadDashboard = () => import(/* webpackChunkName: "dashboard" */ './Dashboard');
const loadProfile = () => import(/* webpackChunkName: "profile" */ './Profile');
```

### Image Optimization
```javascript
// Next.js Image component
import Image from 'next/image';

function OptimizedImage() {
    return (
        <Image
            src="/images/hero.jpg"
            alt="Hero image"
            width={800}
            height={600}
            priority // Load immediately
            placeholder="blur"
            blurDataURL="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQ..."
        />
    );
}

// Responsive images
function ResponsiveImage() {
    return (
        <picture>
            <source 
                media="(min-width: 800px)" 
                srcSet="/images/hero-large.webp" 
                type="image/webp" 
            />
            <source 
                media="(min-width: 400px)" 
                srcSet="/images/hero-medium.webp" 
                type="image/webp" 
            />
            <img 
                src="/images/hero-small.jpg" 
                alt="Hero image"
                loading="lazy"
            />
        </picture>
    );
}
```

### Bundle Optimization
```javascript
// Webpack configuration
module.exports = {
    optimization: {
        splitChunks: {
            chunks: 'all',
            cacheGroups: {
                vendor: {
                    test: /[\\/]node_modules[\\/]/,
                    name: 'vendors',
                    chunks: 'all',
                },
                common: {
                    name: 'common',
                    minChunks: 2,
                    chunks: 'all',
                    enforce: true
                }
            }
        },
        usedExports: true,
        sideEffects: false
    },
    resolve: {
        alias: {
            'lodash': 'lodash-es' // Use ES modules version
        }
    }
};

// Tree shaking
// Instead of importing entire library
import _ from 'lodash';

// Import only what you need
import { debounce, throttle } from 'lodash';
// or
import debounce from 'lodash/debounce';
```

## Monitoring and Metrics

### Application Performance Monitoring
```javascript
// Custom performance monitoring
class PerformanceMonitor {
    constructor() {
        this.metrics = new Map();
    }
    
    startTimer(name) {
        this.metrics.set(name, { start: Date.now() });
    }
    
    endTimer(name) {
        const metric = this.metrics.get(name);
        if (metric) {
            metric.duration = Date.now() - metric.start;
            metric.end = Date.now();
        }
        return metric;
    }
    
    recordMetric(name, value, tags = {}) {
        console.log(`Metric: ${name} = ${value}`, tags);
        // Send to monitoring service (DataDog, New Relic, etc.)
    }
    
    middleware() {
        return (req, res, next) => {
            const startTime = Date.now();
            const timerName = `request_${req.method}_${req.path}`;
            
            this.startTimer(timerName);
            
            res.on('finish', () => {
                const duration = Date.now() - startTime;
                
                this.recordMetric('http_request_duration', duration, {
                    method: req.method,
                    path: req.path,
                    status: res.statusCode
                });
                
                this.recordMetric('http_request_count', 1, {
                    method: req.method,
                    path: req.path,
                    status: res.statusCode
                });
            });
            
            next();
        };
    }
}

const monitor = new PerformanceMonitor();
app.use(monitor.middleware());
```

### Database Performance Monitoring
```javascript
// Sequelize query monitoring
const originalQuery = sequelize.query;
sequelize.query = function(...args) {
    const startTime = Date.now();
    const query = args[0];
    
    return originalQuery.apply(this, args)
        .then(result => {
            const duration = Date.now() - startTime;
            
            if (duration > 100) { // Log slow queries
                console.warn('Slow query detected:', {
                    query: query.substring(0, 100),
                    duration,
                    timestamp: new Date().toISOString()
                });
            }
            
            monitor.recordMetric('db_query_duration', duration, {
                type: query.split(' ')[0].toLowerCase()
            });
            
            return result;
        })
        .catch(error => {
            monitor.recordMetric('db_query_error', 1, {
                error: error.name
            });
            throw error;
        });
};
```

### Real User Monitoring (RUM)
```javascript
// Client-side performance monitoring
class RUMMonitor {
    constructor(apiEndpoint) {
        this.apiEndpoint = apiEndpoint;
        this.init();
    }
    
    init() {
        // Monitor page load performance
        window.addEventListener('load', () => {
            setTimeout(() => this.collectPageMetrics(), 0);
        });
        
        // Monitor navigation performance
        this.observeNavigation();
        
        // Monitor resource loading
        this.observeResources();
    }
    
    collectPageMetrics() {
        const navigation = performance.getEntriesByType('navigation')[0];
        const paint = performance.getEntriesByType('paint');
        
        const metrics = {
            url: window.location.href,
            timestamp: Date.now(),
            navigation: {
                domContentLoaded: navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
                loadComplete: navigation.loadEventEnd - navigation.loadEventStart,
                firstByte: navigation.responseStart - navigation.requestStart,
                domInteractive: navigation.domInteractive - navigation.navigationStart
            },
            paint: {
                firstPaint: paint.find(p => p.name === 'first-paint')?.startTime,
                firstContentfulPaint: paint.find(p => p.name === 'first-contentful-paint')?.startTime
            }
        };
        
        this.sendMetrics(metrics);
    }
    
    observeNavigation() {
        // Monitor SPA navigation
        let lastUrl = window.location.href;
        
        const observer = new MutationObserver(() => {
            if (window.location.href !== lastUrl) {
                lastUrl = window.location.href;
                this.trackPageView();
            }
        });
        
        observer.observe(document.body, {
            childList: true,
            subtree: true
        });
    }
    
    observeResources() {
        const observer = new PerformanceObserver((list) => {
            for (const entry of list.getEntries()) {
                if (entry.duration > 1000) { // Slow resource
                    this.sendMetrics({
                        type: 'slow_resource',
                        name: entry.name,
                        duration: entry.duration,
                        size: entry.transferSize
                    });
                }
            }
        });
        
        observer.observe({ entryTypes: ['resource'] });
    }
    
    async sendMetrics(metrics) {
        try {
            await fetch(this.apiEndpoint, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(metrics)
            });
        } catch (error) {
            console.error('Failed to send metrics:', error);
        }
    }
}

// Initialize RUM monitoring
const rum = new RUMMonitor('/api/metrics');
```

## Scalability Patterns

### Microservices Architecture
```javascript
// API Gateway pattern
const express = require('express');
const httpProxy = require('http-proxy-middleware');

const app = express();

// Service registry
const services = {
    auth: 'http://auth-service:3001',
    users: 'http://user-service:3002',
    orders: 'http://order-service:3003',
    payments: 'http://payment-service:3004'
};

// Route to appropriate microservice
Object.keys(services).forEach(service => {
    app.use(`/api/${service}`, httpProxy({
        target: services[service],
        changeOrigin: true,
        pathRewrite: {
            [`^/api/${service}`]: ''
        }
    }));
});

// Health check aggregation
app.get('/health', async (req, res) => {
    const healthChecks = await Promise.allSettled(
        Object.entries(services).map(async ([name, url]) => {
            try {
                const response = await fetch(`${url}/health`);
                return { service: name, status: response.ok ? 'healthy' : 'unhealthy' };
            } catch (error) {
                return { service: name, status: 'unhealthy', error: error.message };
            }
        })
    );
    
    const results = healthChecks.map(check => check.value || check.reason);
    const allHealthy = results.every(result => result.status === 'healthy');
    
    res.status(allHealthy ? 200 : 503).json({
        status: allHealthy ? 'healthy' : 'degraded',
        services: results
    });
});
```

### Event-Driven Architecture
```javascript
// Event bus implementation
class EventBus {
    constructor() {
        this.events = new Map();
    }
    
    subscribe(eventType, handler) {
        if (!this.events.has(eventType)) {
            this.events.set(eventType, []);
        }
        this.events.get(eventType).push(handler);
    }
    
    async publish(eventType, data) {
        const handlers = this.events.get(eventType) || [];
        
        await Promise.all(
            handlers.map(handler => {
                try {
                    return handler(data);
                } catch (error) {
                    console.error(`Event handler error for ${eventType}:`, error);
                }
            })
        );
    }
}

const eventBus = new EventBus();

// Event handlers
eventBus.subscribe('user.created', async (userData) => {
    // Send welcome email
    await emailService.sendWelcomeEmail(userData.email);
    
    // Create user profile
    await profileService.createProfile(userData.id);
    
    // Track analytics
    await analyticsService.track('user_registered', userData);
});

eventBus.subscribe('order.completed', async (orderData) => {
    // Update inventory
    await inventoryService.updateStock(orderData.items);
    
    // Send confirmation email
    await emailService.sendOrderConfirmation(orderData);
    
    // Generate invoice
    await invoiceService.generateInvoice(orderData);
});

// Publish events
app.post('/api/users', async (req, res) => {
    const user = await User.create(req.body);
    
    // Publish event
    await eventBus.publish('user.created', user);
    
    res.status(201).json(user);
});
```

## Key Takeaways

1. **Measure first**: Always measure before optimizing
2. **Cache strategically**: Use appropriate caching layers
3. **Optimize databases**: Proper indexing and query optimization
4. **Use CDNs**: Distribute static content globally
5. **Load balance**: Distribute traffic across multiple servers
6. **Monitor continuously**: Track performance metrics
7. **Scale horizontally**: Add more servers rather than bigger servers
8. **Optimize frontend**: Reduce bundle sizes and improve loading

## Next Steps

In the next part, we'll cover:
- Advanced networking concepts
- Microservices communication
- Message queues and pub/sub
- Real-time communication patterns

---

*Continue to [Part 5: Advanced Networking Concepts](05-advanced-networking.md)*