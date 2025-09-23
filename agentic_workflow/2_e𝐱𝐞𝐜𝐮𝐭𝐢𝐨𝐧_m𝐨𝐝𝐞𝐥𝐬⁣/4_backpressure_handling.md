# Backpressure Handling in JavaScript

Backpressure occurs when data is produced faster than it can be consumed. Proper backpressure handling prevents memory issues, system crashes, and ensures stable data flow in streaming applications.

## Core Concepts

### **What is Backpressure?**
- **Data flow imbalance**: Producer > Consumer
- **Memory buildup**: Unprocessed data accumulates
- **System stress**: Can lead to crashes or data loss
- **Flow control**: Managing the data production rate

### **Backpressure Strategies**
- **Pausing**: Temporarily stop production
- **Buffering**: Store excess data temporarily
- **Dropping**: Discard data when necessary
- **Throttling**: Slow down production rate

---

## JavaScript Examples

### Example 1: Basic Backpressure-aware Stream

```javascript
class BackpressureAwareStream {
    constructor(options = {}) {
        this.highWaterMark = options.highWaterMark || 100; // Max buffer size
        this.lowWaterMark = options.lowWaterMark || 20; // Resume threshold
        this.buffer = [];
        this.isPaused = false;
        this.producers = new Set();
        this.consumers = new Set();
        
        this.stats = {
            totalProduced: 0,
            totalConsumed: 0,
            backpressureEvents: 0,
            droppedItems: 0,
            maxBufferSize: 0
        };
    }

    // Producer: Add data to stream
    async write(data) {
        // Apply backpressure if buffer is full
        if (this.buffer.length >= this.highWaterMark) {
            this.stats.backpressureEvents++;
            
            if (this.isPaused) {
                // Buffer is still full, consider dropping data
                this.stats.droppedItems++;
                throw new Error('Backpressure: Buffer full, data dropped');
            }
            
            // Pause production
            this.isPaused = true;
            console.log('⏸️ Backpressure applied: Pausing production');
            
            // Wait for buffer to drain
            await this.waitForBufferDrain();
        }

        this.buffer.push(data);
        this.stats.totalProduced++;
        this.stats.maxBufferSize = Math.max(this.stats.maxBufferSize, this.buffer.length);
        
        // Notify consumers
        this.notifyConsumers();
        
        return true;
    }

    // Consumer: Read data from stream
    async read() {
        if (this.buffer.length === 0) {
            // Wait for data
            await new Promise(resolve => {
                this.consumers.add(resolve);
            });
        }

        const data = this.buffer.shift();
        this.stats.totalConsumed++;
        
        // Check if we can resume production
        if (this.isPaused && this.buffer.length <= this.lowWaterMark) {
            this.isPaused = false;
            console.log('▶️ Backpressure released: Resuming production');
            this.notifyProducers();
        }
        
        return data;
    }

    // Wait for buffer to drain to low water mark
    async waitForBufferDrain() {
        return new Promise(resolve => {
            const checkBuffer = () => {
                if (this.buffer.length <= this.lowWaterMark) {
                    resolve();
                } else {
                    setTimeout(checkBuffer, 10);
                }
            };
            checkBuffer();
        });
    }

    notifyConsumers() {
        this.consumers.forEach(resolve => {
            resolve();
            this.consumers.delete(resolve);
        });
    }

    notifyProducers() {
        this.producers.forEach(resolve => {
            resolve();
            this.producers.delete(resolve);
        });
    }

    getStats() {
        return {
            ...this.stats,
            currentBufferSize: this.buffer.length,
            isPaused: this.isPaused,
            waitingProducers: this.producers.size,
            waitingConsumers: this.consumers.size
        };
    }
}
```

### Example 2: Node.js Stream with Backpressure

```javascript
const { Readable, Writable, Transform } = require('stream');
const { pipeline } = require('stream/promises');

class BackpressureHandlingDemo {
    
    // Fast producer stream
    static createFastProducer() {
        let count = 0;
        
        return new Readable({
            objectMode: true,
            highWaterMark: 10, // Small buffer to trigger backpressure quickly
            
            read(size) {
                // Simulate fast data production
                if (count >= 100) {
                    this.push(null); // End stream
                    return;
                }
                
                const canPush = this.push({
                    id: count++,
                    data: `Item ${count}`,
                    timestamp: Date.now()
                });
                
                console.log(`Produced item ${count}, canPush: ${canPush}`);
                
                // If backpressure is applied (canPush is false), slow down
                if (!canPush) {
                    console.log('⏸️ Backpressure detected: Slowing down production');
                    // In real scenario, you might pause or slow down
                }
            }
        });
    }
    
    // Slow consumer stream with backpressure handling
    static createSlowConsumer() {
        return new Writable({
            objectMode: true,
            highWaterMark: 5, // Small buffer to show backpressure
            
            async write(chunk, encoding, callback) {
                console.log(`Consuming item ${chunk.id}...`);
                
                // Simulate slow processing
                await new Promise(resolve => setTimeout(resolve, 200));
                
                console.log(`✅ Processed item ${chunk.id}`);
                
                // Signal that we're ready for more data
                callback();
            },
            
            // Handle backpressure
            final(callback) {
                console.log('Consumer finished');
                callback();
            }
        });
    }
    
    // Transform stream with explicit backpressure handling
    static createProcessingPipeline() {
        let processedCount = 0;
        
        return new Transform({
            objectMode: true,
            highWaterMark: 8,
            
            async transform(chunk, encoding, callback) {
                try {
                    // Simulate variable processing time
                    const processingTime = 100 + Math.random() * 200;
                    await new Promise(resolve => setTimeout(resolve, processingTime));
                    
                    processedCount++;
                    const processedChunk = {
                        ...chunk,
                        processedAt: Date.now(),
                        processedCount: processedCount
                    };
                    
                    // Check if we should apply backpressure to upstream
                    const shouldApplyBackpressure = processedCount % 20 === 0;
                    if (shouldApplyBackpressure) {
                        console.log('🔄 Intentional backpressure for flow control');
                        await new Promise(resolve => setTimeout(resolve, 500));
                    }
                    
                    callback(null, processedChunk);
                } catch (error) {
                    callback(error);
                }
            }
        });
    }
    
    static async demonstrateBackpressure() {
        console.log('🚀 Starting backpressure demonstration...\n');
        
        try {
            await pipeline(
                this.createFastProducer(),
                this.createProcessingPipeline(),
                this.createSlowConsumer()
            );
            
            console.log('\n✅ Pipeline completed successfully');
        } catch (error) {
            console.error('Pipeline error:', error);
        }
    }
}

// BackpressureHandlingDemo.demonstrateBackpressure();
```

### Example 3: Advanced Backpressure Controller

```javascript
class AdvancedBackpressureController {
    constructor(options = {}) {
        this.maxBufferSize = options.maxBufferSize || 1000;
        this.targetLatency = options.targetLatency || 100; // ms
        this.adaptationRate = options.adaptationRate || 0.1;
        
        this.buffer = [];
        this.consumptionRate = 0;
        this.productionRate = 0;
        this.currentPressure = 0; // 0 = no pressure, 1 = max pressure
        
        this.metrics = {
            totalProduced: 0,
            totalConsumed: 0,
            backpressureEvents: 0,
            avgConsumptionRate: 0,
            avgProductionRate: 0,
            pressureHistory: []
        };
        
        this._startMetricsCollection();
    }
    
    // Adaptive backpressure based on system metrics
    async shouldProduce() {
        const currentBufferSize = this.buffer.length;
        const bufferUtilization = currentBufferSize / this.maxBufferSize;
        
        // Calculate pressure based on multiple factors
        const pressureFromBuffer = Math.min(bufferUtilization * 2, 1); // 0-1 scale
        const pressureFromRateDifference = Math.max(0, 
            (this.productionRate - this.consumptionRate) / Math.max(this.consumptionRate, 1)
        );
        
        this.currentPressure = Math.max(pressureFromBuffer, pressureFromRateDifference);
        this.metrics.pressureHistory.push({
            timestamp: Date.now(),
            pressure: this.currentPressure,
            bufferSize: currentBufferSize
        });
        
        // Keep only recent history
        if (this.metrics.pressureHistory.length > 1000) {
            this.metrics.pressureHistory = this.metrics.pressureHistory.slice(-1000);
        }
        
        // Apply backpressure based on current pressure
        if (this.currentPressure > 0.8) {
            this.metrics.backpressureEvents++;
            console.log(`🚨 High backpressure: ${(this.currentPressure * 100).toFixed(1)}%`);
            
            // Aggressive backpressure
            await this.applyBackpressure('high');
            return false;
            
        } else if (this.currentPressure > 0.5) {
            this.metrics.backpressureEvents++;
            console.log(`⚠️ Medium backpressure: ${(this.currentPressure * 100).toFixed(1)}%`);
            
            // Moderate backpressure
            await this.applyBackpressure('medium');
            return true;
            
        } else if (this.currentPressure > 0.3) {
            // Mild backpressure - just slow down
            await this.applyBackpressure('low');
            return true;
        }
        
        // No backpressure needed
        return true;
    }
    
    async applyBackpressure(level) {
        const strategies = {
            high: {
                delay: 500,
                message: 'Aggressive backpressure: Pausing production'
            },
            medium: {
                delay: 200,
                message: 'Moderate backpressure: Slowing production'
            },
            low: {
                delay: 50,
                message: 'Mild backpressure: Minor slowdown'
            }
        };
        
        const strategy = strategies[level];
        console.log(strategy.message);
        
        // Apply delay based on pressure level
        await new Promise(resolve => setTimeout(resolve, strategy.delay));
        
        // Adaptive learning: Adjust based on effectiveness
        this.adaptStrategy();
    }
    
    adaptStrategy() {
        // Simple adaptation: if pressure remains high, be more aggressive
        const recentPressure = this.metrics.pressureHistory
            .slice(-10)
            .reduce((sum, point) => sum + point.pressure, 0) / 10;
        
        if (recentPressure > 0.7) {
            this.targetLatency = Math.min(this.targetLatency * 1.1, 1000);
        } else if (recentPressure < 0.2) {
            this.targetLatency = Math.max(this.targetLatency * 0.9, 10);
        }
    }
    
    async produce(data) {
        const canProduce = await this.shouldProduce();
        
        if (!canProduce) {
            throw new Error('Backpressure: Cannot produce at this time');
        }
        
        if (this.buffer.length >= this.maxBufferSize) {
            throw new Error('Backpressure: Buffer overflow');
        }
        
        this.buffer.push(data);
        this.metrics.totalProduced++;
        this.productionRate = this.calculateCurrentRate('production');
        
        return true;
    }
    
    async consume() {
        if (this.buffer.length === 0) {
            return null;
        }
        
        const data = this.buffer.shift();
        this.metrics.totalConsumed++;
        this.consumptionRate = this.calculateCurrentRate('consumption');
        
        return data;
    }
    
    calculateCurrentRate(type) {
        const now = Date.now();
        const history = this.metrics.pressureHistory
            .filter(point => now - point.timestamp < 5000); // Last 5 seconds
        
        if (history.length < 2) return 0;
        
        const timeWindow = (now - history[0].timestamp) / 1000; // seconds
        const count = type === 'production' ? this.metrics.totalProduced : this.metrics.totalConsumed;
        
        // Simplified rate calculation
        return count / Math.max(timeWindow, 1);
    }
    
    _startMetricsCollection() {
        setInterval(() => {
            this.metrics.avgProductionRate = this.productionRate;
            this.metrics.avgConsumptionRate = this.consumptionRate;
        }, 1000);
    }
    
    getMetrics() {
        return {
            ...this.metrics,
            currentPressure: this.currentPressure,
            bufferUtilization: this.buffer.length / this.maxBufferSize,
            rateDifference: this.productionRate - this.consumptionRate,
            targetLatency: this.targetLatency
        };
    }
    
    // Emergency pressure release
    releasePressure() {
        const itemsToRelease = Math.floor(this.buffer.length * 0.5); // Release 50%
        this.buffer = this.buffer.slice(itemsToRelease);
        console.log(`🆘 Emergency pressure release: Dropped ${itemsToRelease} items`);
    }
}
```

### Example 4: Producer-Consumer with Backpressure

```javascript
class BackpressureAwareProducerConsumer {
    constructor(options = {}) {
        this.bufferSize = options.bufferSize || 50;
        this.productionDelay = options.productionDelay || 10;
        this.consumptionDelay = options.consumptionDelay || 100;
        
        this.buffer = [];
        this.isProducing = false;
        this.isConsuming = false;
        this.backpressureCallbacks = [];
        
        this.stats = {
            produced: 0,
            consumed: 0,
            dropped: 0,
            backpressureApplied: 0,
            bufferHigh: 0
        };
    }
    
    // Producer that respects backpressure
    async startProducer() {
        this.isProducing = true;
        let itemId = 0;
        
        while (this.isProducing) {
            // Check backpressure before producing
            if (await this.checkBackpressure()) {
                const item = { id: itemId++, data: `Item ${itemId}`, timestamp: Date.now() };
                
                if (this.buffer.length < this.bufferSize) {
                    this.buffer.push(item);
                    this.stats.produced++;
                    console.log(`📤 Produced item ${item.id}, buffer: ${this.buffer.length}/${this.bufferSize}`);
                } else {
                    // Buffer full - apply backpressure
                    this.stats.dropped++;
                    console.log('❌ Buffer full, item dropped');
                    await this.applyBackpressure();
                }
            }
            
            await this.delay(this.productionDelay);
        }
    }
    
    // Consumer that processes items
    async startConsumer() {
        this.isConsuming = true;
        
        while (this.isConsuming) {
            if (this.buffer.length > 0) {
                const item = this.buffer.shift();
                this.stats.consumed++;
                
                console.log(`📥 Consuming item ${item.id}, buffer: ${this.buffer.length}/${this.bufferSize}`);
                
                // Simulate processing time
                await this.delay(this.consumptionDelay);
                
                console.log(`✅ Processed item ${item.id}`);
                
                // Notify that pressure might be released
                this.notifyBackpressureRelease();
            } else {
                // No items to process
                await this.delay(50);
            }
        }
    }
    
    // Backpressure check with adaptive strategy
    async checkBackpressure() {
        const utilization = this.buffer.length / this.bufferSize;
        this.stats.bufferHigh = Math.max(this.stats.bufferHigh, utilization);
        
        if (utilization > 0.9) {
            // Critical backpressure
            this.stats.backpressureApplied++;
            console.log('🚨 CRITICAL backpressure: 90%+ buffer utilization');
            await this.applyBackpressure();
            return false;
        } else if (utilization > 0.7) {
            // High backpressure
            this.stats.backpressureApplied++;
            console.log('⚠️ HIGH backpressure: 70%+ buffer utilization');
            await this.delay(this.productionDelay * 5); // Slow down 5x
            return true;
        } else if (utilization > 0.5) {
            // Medium backpressure
            console.log('🔶 MEDIUM backpressure: 50%+ buffer utilization');
            await this.delay(this.productionDelay * 2); // Slow down 2x
            return true;
        }
        
        // No backpressure needed
        return true;
    }
    
    async applyBackpressure() {
        console.log('⏸️ Applying backpressure...');
        
        // Notify all backpressure listeners
        this.backpressureCallbacks.forEach(callback => callback('high'));
        
        // Wait for buffer to drain
        while (this.buffer.length > this.bufferSize * 0.5) {
            console.log(`Waiting for buffer drain... ${this.buffer.length}/${this.bufferSize}`);
            await this.delay(100);
        }
        
        console.log('▶️ Backpressure released');
        this.backpressureCallbacks.forEach(callback => callback('low'));
    }
    
    notifyBackpressureRelease() {
        if (this.buffer.length < this.bufferSize * 0.3) {
            this.backpressureCallbacks.forEach(callback => callback('low'));
        }
    }
    
    onBackpressure(callback) {
        this.backpressureCallbacks.push(callback);
    }
    
    stop() {
        this.isProducing = false;
        this.isConsuming = false;
    }
    
    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
    
    getStats() {
        return {
            ...this.stats,
            currentBufferSize: this.buffer.length,
            bufferUtilization: (this.buffer.length / this.bufferSize) * 100,
            isProducing: this.isProducing,
            isConsuming: this.isConsuming
        };
    }
}

// Demonstration
async function demonstrateProducerConsumer() {
    const system = new BackpressureAwareProducerConsumer({
        bufferSize: 20,
        productionDelay: 10,  // Fast producer
        consumptionDelay: 200 // Slow consumer
    });
    
    // Monitor backpressure events
    system.onBackpressure((level) => {
        console.log(`🔔 Backpressure level: ${level}`);
    });
    
    // Start producer and consumer
    const producerPromise = system.startProducer();
    const consumerPromise = system.startConsumer();
    
    // Monitor stats
    const statsInterval = setInterval(() => {
        const stats = system.getStats();
        console.log('\n📊 Current Stats:', stats);
    }, 5000);
    
    // Run for 30 seconds
    await system.delay(30000);
    
    system.stop();
    clearInterval(statsInterval);
    
    await Promise.allSettled([producerPromise, consumerPromise]);
    
    console.log('\n🎯 Final Stats:', system.getStats());
}

// demonstrateProducerConsumer();
```

### Example 5: Real-world HTTP Stream with Backpressure

```javascript
class HTTPStreamWithBackpressure {
    constructor() {
        this.activeStreams = new Map();
        this.globalBackpressure = 0;
    }
    
    async streamWithBackpressure(url, onData, options = {}) {
        const streamId = Math.random().toString(36).substr(2, 9);
        const controller = new AbortController();
        
        const streamInfo = {
            id: streamId,
            controller,
            bytesReceived: 0,
            chunksReceived: 0,
            backpressureApplied: 0,
            startTime: Date.now()
        };
        
        this.activeStreams.set(streamId, streamInfo);
        
        try {
            const response = await fetch(url, {
                signal: controller.signal,
                headers: options.headers
            });
            
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}`);
            }
            
            if (!response.body) {
                throw new Error('Response body not available for streaming');
            }
            
            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';
            
            while (true) {
                // Check global backpressure
                if (await this.checkGlobalBackpressure()) {
                    streamInfo.backpressureApplied++;
                    console.log('🌍 Global backpressure applied, pausing stream');
                    await this.delay(100);
                    continue;
                }
                
                const { done, value } = await reader.read();
                
                if (done) break;
                
                streamInfo.bytesReceived += value.length;
                streamInfo.chunksReceived++;
                
                buffer += decoder.decode(value, { stream: true });
                const lines = buffer.split('\n');
                buffer = lines.pop() || '';
                
                // Process each line with backpressure check
                for (const line of lines) {
                    if (line.trim() === '') continue;
                    
                    // Check if consumer can handle more data
                    if (await this.shouldDelayProcessing()) {
                        streamInfo.backpressureApplied++;
                        console.log('⏸️ Processing backpressure, delaying consumption');
                        await this.delay(50);
                    }
                    
                    try {
                        await onData(line.trim(), streamInfo);
                    } catch (error) {
                        console.error('Data processing error:', error);
                        // Continue with next line despite error
                    }
                }
                
                // Update global backpressure based on system load
                this.updateGlobalBackpressure();
            }
            
            // Process remaining buffer
            if (buffer.trim()) {
                await onData(buffer.trim(), streamInfo);
            }
            
            return streamInfo;
            
        } catch (error) {
            if (error.name !== 'AbortError') {
                console.error('Stream error:', error);
                throw error;
            }
        } finally {
            this.activeStreams.delete(streamId);
        }
    }
    
    async checkGlobalBackpressure() {
        // Simulate global system load check
        const activeStreams = this.activeStreams.size;
        const totalBytes = Array.from(this.activeStreams.values())
            .reduce((sum, info) => sum + info.bytesReceived, 0);
        
        // Simple backpressure heuristic
        const pressureFromStreams = activeStreams / 10; // Max 10 streams
        const pressureFromBytes = Math.min(totalBytes / (1024 * 1024), 1); // Max 1MB
        
        this.globalBackpressure = Math.max(pressureFromStreams, pressureFromBytes);
        
        return this.globalBackpressure > 0.7;
    }
    
    async shouldDelayProcessing() {
        // Simulate processing capacity check
        const memoryUsage = process.memoryUsage ? process.memoryUsage().heapUsed / (1024 * 1024) : 0;
        const memoryPressure = memoryUsage > 100 ? 1 : memoryUsage / 100; // Max 100MB
        
        // Random processing delay simulation
        const randomLoad = Math.random() > 0.8; // 20% chance of being busy
        
        return memoryPressure > 0.8 || randomLoad;
    }
    
    updateGlobalBackpressure() {
        // Reduce backpressure over time
        this.globalBackpressure = Math.max(0, this.globalBackpressure - 0.1);
    }
    
    // Apply backpressure to a specific stream
    throttleStream(streamId, factor = 2) {
        const streamInfo = this.activeStreams.get(streamId);
        if (streamInfo) {
            console.log(`🐌 Throttling stream ${streamId} by factor ${factor}`);
            // In real implementation, this would adjust stream parameters
        }
    }
    
    // Emergency: Stop all streams
    emergencyStop() {
        console.log('🆘 EMERGENCY: Stopping all streams');
        this.activeStreams.forEach((info, id) => {
            info.controller.abort();
        });
        this.activeStreams.clear();
    }
    
    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
    
    getSystemStats() {
        return {
            activeStreams: this.activeStreams.size,
            globalBackpressure: this.globalBackpressure,
            totalBytes: Array.from(this.activeStreams.values())
                .reduce((sum, info) => sum + info.bytesReceived, 0),
            avgBackpressure: Array.from(this.activeStreams.values())
                .reduce((sum, info) => sum + info.backpressureApplied, 0) / Math.max(this.activeStreams.size, 1)
        };
    }
}
```

## Key Backpressure Strategies

### **1. Reactive Strategies**
- **Pause/Resume**: Stop production when buffer is full
- **Slow Down**: Reduce production rate gradually
- **Buffer Management**: Monitor and control buffer size

### **2. Proactive Strategies**
- **Rate Limiting**: Prevent over-production from start
- **Adaptive Speed**: Adjust production based on consumer capacity
- **Load Shedding**: Discard low-priority data when necessary

### **3. System-level Strategies**
- **Global Backpressure**: Coordinate across multiple streams
- **Resource Monitoring**: Consider memory, CPU, network
- **Circuit Breakers**: Stop entirely when system is overwhelmed

## Best Practices

1. **Monitor Buffer Levels**: Always know your buffer utilization
2. **Implement Graceful Degradation**: Slow down before stopping completely
3. **Provide Feedback**: Let producers know about backpressure
4. **Test Under Load**: Simulate high-pressure scenarios
5. **Plan for Failure**: Have strategies for when backpressure isn't enough

Backpressure handling is crucial for building robust, production-ready streaming systems that can handle real-world load conditions gracefully.