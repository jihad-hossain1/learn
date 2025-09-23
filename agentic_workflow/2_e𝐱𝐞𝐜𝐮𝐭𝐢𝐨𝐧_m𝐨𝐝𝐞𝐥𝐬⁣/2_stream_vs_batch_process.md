# Streaming vs. Batch Processing in JavaScript

Streaming and batch processing are two fundamental data processing paradigms with different characteristics and use cases.

## Core Concepts

### **Batch Processing**
- **Processes data in large chunks/groups**
- **Scheduled or triggered** at specific intervals
- **High latency** - processes data after it's collected
- **Good for analytics, reporting, ETL jobs**
- **Example**: Nightly sales report generation

### **Stream Processing**
- **Processes data continuously as it arrives**
- **Real-time or near-real-time** processing
- **Low latency** - processes data immediately
- **Good for monitoring, alerts, real-time features**
- **Example**: Live stock price updates

---

## JavaScript Examples

### Example 1: Basic Batch Processing

```javascript
class BatchProcessor {
    constructor(batchSize = 100, processInterval = 5000) {
        this.batchSize = batchSize;
        this.processInterval = processInterval;
        this.batch = [];
        this.isProcessing = false;
        
        // Process batch when size limit reached or time interval elapsed
        this.intervalId = setInterval(() => this.processBatch(), processInterval);
    }

    addData(data) {
        this.batch.push({
            data,
            timestamp: Date.now(),
            id: Math.random().toString(36).substr(2, 9)
        });

        console.log(`Added data to batch. Current size: ${this.batch.length}/${this.batchSize}`);

        // Process if batch size reached
        if (this.batch.length >= this.batchSize) {
            this.processBatch();
        }
    }

    async processBatch() {
        if (this.isProcessing || this.batch.length === 0) {
            return;
        }

        this.isProcessing = true;
        const currentBatch = [...this.batch];
        this.batch = []; // Clear batch

        console.log(`🚀 Processing batch of ${currentBatch.length} items...`);

        try {
            // Simulate batch processing (database insert, analytics, etc.)
            const results = await this.processBatchData(currentBatch);
            
            console.log(`✅ Batch processed successfully. Results:`, results);
            console.log(`📊 Batch stats: ${currentBatch.length} items, avg processing time: ${results.averageProcessingTime}ms`);
            
        } catch (error) {
            console.error('❌ Batch processing failed:', error);
            // Retry logic or move to dead letter queue
            this.batch = [...currentBatch, ...this.batch]; // Put failed batch back
        } finally {
            this.isProcessing = false;
        }
    }

    async processBatchData(batch) {
        // Simulate processing time (database operations, API calls, etc.)
        const startTime = Date.now();
        
        // Process each item in batch
        const processedItems = await Promise.all(
            batch.map(async (item, index) => {
                await this.delay(10 + Math.random() * 50); // Simulate work
                return {
                    ...item,
                    processedAt: Date.now(),
                    result: `Processed ${item.data} successfully`
                };
            })
        );

        const totalTime = Date.now() - startTime;
        
        return {
            processedItems,
            totalProcessingTime: totalTime,
            averageProcessingTime: totalTime / batch.length,
            batchSize: batch.length
        };
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    getStats() {
        return {
            currentBatchSize: this.batch.length,
            isProcessing: this.isProcessing,
            batchSizeLimit: this.batchSize,
            processInterval: this.processInterval
        };
    }

    destroy() {
        clearInterval(this.intervalId);
        // Process remaining batch
        if (this.batch.length > 0) {
            this.processBatch();
        }
    }
}
```

### Example 2: Basic Stream Processing

```javascript
class StreamProcessor {
    constructor(processingFn, options = {}) {
        this.processingFn = processingFn;
        this.bufferSize = options.bufferSize || 1000;
        this.maxConcurrent = options.maxConcurrent || 5;
        this.buffer = [];
        this.processing = new Set();
        this.stats = {
            processed: 0,
            errors: 0,
            startedAt: Date.now()
        };

        this._processBuffer = this._processBuffer.bind(this);
        // Start processing buffer periodically
        this.intervalId = setInterval(this._processBuffer, options.processInterval || 100);
    }

    async addData(data) {
        const streamItem = {
            data,
            id: Math.random().toString(36).substr(2, 9),
            timestamp: Date.now(),
            status: 'queued'
        };

        this.buffer.push(streamItem);

        // Apply backpressure if buffer is too large
        if (this.buffer.length > this.bufferSize) {
            console.warn('🚨 Buffer overflow, applying backpressure');
            // Remove oldest items or wait for buffer to clear
            this.buffer.shift();
        }

        console.log(`📥 Stream received data. Buffer size: ${this.buffer.length}`);
        
        // Try to process immediately if we have capacity
        this._processBuffer();
        
        return streamItem.id;
    }

    async _processBuffer() {
        // Don't process if we're at capacity or buffer is empty
        if (this.processing.size >= this.maxConcurrent || this.buffer.length === 0) {
            return;
        }

        const availableSlots = this.maxConcurrent - this.processing.size;
        const itemsToProcess = this.buffer.splice(0, availableSlots);

        for (const item of itemsToProcess) {
            this._processItem(item);
        }
    }

    async _processItem(item) {
        item.status = 'processing';
        this.processing.add(item);

        try {
            const startTime = Date.now();
            const result = await this.processingFn(item.data);
            const processingTime = Date.now() - startTime;

            item.status = 'completed';
            item.result = result;
            item.processingTime = processingTime;

            this.stats.processed++;
            
            console.log(`✅ Stream processed item ${item.id} in ${processingTime}ms`);

        } catch (error) {
            item.status = 'error';
            item.error = error;
            this.stats.errors++;
            
            console.error(`❌ Stream processing error for item ${item.id}:`, error);
        } finally {
            this.processing.delete(item);
            // Continue processing more items
            this._processBuffer();
        }
    }

    getStats() {
        const runningTime = Date.now() - this.stats.startedAt;
        return {
            ...this.stats,
            bufferSize: this.buffer.length,
            currentlyProcessing: this.processing.size,
            itemsPerSecond: this.stats.processed / (runningTime / 1000),
            runningTime
        };
    }

    destroy() {
        clearInterval(this.intervalId);
        console.log('Stream processor destroyed');
    }
}
```

### Example 3: Real-World Comparison - Data Analytics

```javascript
// Simulate an e-commerce analytics system
class EcommerceAnalytics {
    constructor() {
        // Batch processor for daily reports
        this.batchProcessor = new BatchProcessor(50, 30000); // 30-second intervals
        
        // Stream processor for real-time alerts
        this.streamProcessor = new StreamProcessor(
            this._processRealTimeEvent.bind(this),
            { maxConcurrent: 3, bufferSize: 500 }
        );
        
        this.salesData = [];
        this.realTimeAlerts = [];
    }

    // Batch processing: Daily sales reports
    recordSale(sale) {
        const saleRecord = {
            ...sale,
            id: Math.random().toString(36).substr(2, 9),
            timestamp: new Date().toISOString()
        };

        // Add to batch for nightly processing
        this.batchProcessor.addData(saleRecord);
        
        // Also process in real-time for alerts
        this.streamProcessor.addData(saleRecord);

        this.salesData.push(saleRecord);
        console.log(`💰 Sale recorded: ${sale.product} for $${sale.amount}`);
    }

    // Stream processing: Real-time alerting
    async _processRealTimeEvent(sale) {
        // Check for high-value purchases (fraud detection)
        if (sale.amount > 1000) {
            this.realTimeAlerts.push({
                type: 'HIGH_VALUE_SALE',
                sale,
                message: `High value sale: $${sale.amount} for ${sale.product}`,
                timestamp: new Date().toISOString()
            });
            console.log(`🚨 ALERT: High value sale detected!`);
        }

        // Check for rapid purchases (potential fraud)
        const recentSales = this.salesData.filter(s => 
            Date.now() - new Date(s.timestamp).getTime() < 60000 // Last minute
        );
        
        if (recentSales.length > 5) {
            this.realTimeAlerts.push({
                type: 'RAPID_PURCHASES',
                sales: recentSales,
                message: `Rapid purchases detected: ${recentSales.length} in 1 minute`,
                timestamp: new Date().toISOString()
            });
            console.log(`🚨 ALERT: Rapid purchases detected!`);
        }

        // Simulate real-time processing
        await this.delay(5 + Math.random() * 15);
        
        return { status: 'processed', alertGenerated: this.realTimeAlerts.length };
    }

    // Generate batch report
    async generateDailyReport() {
        console.log('📊 Generating daily sales report...');
        
        const today = new Date().toDateString();
        const todaySales = this.salesData.filter(sale => 
            new Date(sale.timestamp).toDateString() === today
        );

        const report = {
            date: today,
            totalSales: todaySales.length,
            totalRevenue: todaySales.reduce((sum, sale) => sum + sale.amount, 0),
            averageSale: todaySales.reduce((sum, sale) => sum + sale.amount, 0) / todaySales.length,
            topProducts: this._getTopProducts(todaySales),
            alerts: this.realTimeAlerts.filter(alert => 
                new Date(alert.timestamp).toDateString() === today
            )
        };

        console.log('✅ Daily report generated:', report);
        return report;
    }

    _getTopProducts(sales) {
        const productCount = {};
        sales.forEach(sale => {
            productCount[sale.product] = (productCount[sale.product] || 0) + 1;
        });
        
        return Object.entries(productCount)
            .sort(([,a], [,b]) => b - a)
            .slice(0, 5)
            .map(([product, count]) => ({ product, count }));
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    getPerformanceStats() {
        return {
            batch: this.batchProcessor.getStats(),
            stream: this.streamProcessor.getStats(),
            totalSales: this.salesData.length,
            totalAlerts: this.realTimeAlerts.length
        };
    }
}
```

### Example 4: File Processing Comparison

```javascript
const fs = require('fs');
const readline = require('readline');
const { Transform, PassThrough } = require('stream');

class FileProcessor {
    // Batch processing: Read entire file into memory
    static async processFileBatch(filePath, processFn) {
        console.log('🔄 Starting BATCH file processing...');
        const startTime = Date.now();

        try {
            // Read entire file into memory
            const data = await fs.promises.readFile(filePath, 'utf8');
            const lines = data.split('\n').filter(line => line.trim());
            
            console.log(`📁 File loaded: ${lines.length} lines, ${data.length} bytes`);

            // Process all lines in batch
            const results = await Promise.all(
                lines.map(async (line, index) => {
                    return await processFn(line, index, lines.length);
                })
            );

            const processingTime = Date.now() - startTime;
            console.log(`✅ Batch processing completed in ${processingTime}ms`);
            console.log(`📊 Processed ${lines.length} lines at ${(lines.length / (processingTime / 1000)).toFixed(2)} lines/second`);

            return results;

        } catch (error) {
            console.error('❌ Batch processing failed:', error);
            throw error;
        }
    }

    // Stream processing: Process file line by line
    static async processFileStream(filePath, processFn) {
        console.log('🌊 Starting STREAM file processing...');
        const startTime = Date.now();
        let lineCount = 0;
        const results = [];

        return new Promise((resolve, reject) => {
            const fileStream = fs.createReadStream(filePath, { encoding: 'utf8' });
            const rl = readline.createInterface({
                input: fileStream,
                crlfDelay: Infinity
            });

            // Process each line as it's read
            rl.on('line', async (line) => {
                lineCount++;
                if (line.trim()) {
                    try {
                        const result = await processFn(line, lineCount);
                        results.push(result);
                        
                        // Show progress every 1000 lines
                        if (lineCount % 1000 === 0) {
                            const elapsed = Date.now() - startTime;
                            console.log(`📈 Stream processed ${lineCount} lines... (${(lineCount / (elapsed / 1000)).toFixed(2)} lines/sec)`);
                        }
                    } catch (error) {
                        console.error(`Error processing line ${lineCount}:`, error);
                    }
                }
            });

            rl.on('close', () => {
                const processingTime = Date.now() - startTime;
                console.log(`✅ Stream processing completed in ${processingTime}ms`);
                console.log(`📊 Processed ${lineCount} lines at ${(lineCount / (processingTime / 1000)).toFixed(2)} lines/second`);
                resolve(results);
            });

            rl.on('error', reject);
            fileStream.on('error', reject);
        });
    }

    // Transform stream for real-time data processing
    static createProcessingStream(processFn) {
        return new Transform({
            objectMode: true,
            transform(chunk, encoding, callback) {
                processFn(chunk.toString())
                    .then(result => callback(null, result))
                    .catch(error => callback(error));
            }
        });
    }
}

// Usage example for file processing
async function fileProcessingDemo() {
    // Create a sample large file
    const filePath = './large_data.txt';
    await createSampleFile(filePath, 10000); // 10,000 lines

    // Process function example
    const processLine = async (line, index, total) => {
        // Simulate some processing work
        await new Promise(resolve => setTimeout(resolve, 1));
        return {
            line: index + 1,
            content: line.toUpperCase(),
            length: line.length
        };
    };

    console.log('=== BATCH PROCESSING DEMO ===');
    await FileProcessor.processFileBatch(filePath, processLine);

    console.log('\n=== STREAM PROCESSING DEMO ===');
    await FileProcessor.processFileStream(filePath, processLine);

    // Cleanup
    await fs.promises.unlink(filePath);
}

async function createSampleFile(filePath, lineCount) {
    const lines = [];
    for (let i = 0; i < lineCount; i++) {
        lines.push(`Line ${i + 1}: ${Math.random().toString(36).substr(2, 10)}`);
    }
    await fs.promises.writeFile(filePath, lines.join('\n'));
}
```

### Example 5: Real-time Dashboard with Both Approaches

```javascript
class RealTimeDashboard {
    constructor() {
        this.batchData = [];
        this.streamData = [];
        this.metrics = {
            totalUsers: 0,
            activeUsers: 0,
            eventsProcessed: 0
        };

        // Batch processor for historical data
        this.batchProcessor = new BatchProcessor(100, 10000); // 10-second batches
        
        // Stream processor for real-time metrics
        this.streamProcessor = new StreamProcessor(
            this._updateRealTimeMetrics.bind(this),
            { maxConcurrent: 10, bufferSize: 1000 }
        );

        this._setupDashboard();
    }

    // Simulate user events
    simulateUserActivity() {
        setInterval(() => {
            const event = {
                type: ['pageview', 'click', 'purchase', 'login'][Math.floor(Math.random() * 4)],
                userId: Math.floor(Math.random() * 1000),
                timestamp: new Date().toISOString(),
                data: { page: '/home', duration: Math.random() * 100 }
            };

            this.processEvent(event);
        }, 100); // 10 events per second
    }

    processEvent(event) {
        // Send to batch processing for historical analysis
        this.batchProcessor.addData(event);
        
        // Send to stream processing for real-time updates
        this.streamProcessor.addData(event);
    }

    async _updateRealTimeMetrics(event) {
        // Update real-time metrics
        this.metrics.eventsProcessed++;
        
        if (event.type === 'login') {
            this.metrics.activeUsers++;
            // Set timeout to decrement active users after 5 minutes
            setTimeout(() => {
                this.metrics.activeUsers = Math.max(0, this.metrics.activeUsers - 1);
            }, 300000);
        }

        // Simulate processing time
        await this.delay(2 + Math.random() * 8);

        return { metric: 'updated', eventType: event.type };
    }

    // Batch processing for historical data analysis
    async analyzeHistoricalData() {
        console.log('📈 Analyzing historical data...');
        
        // This would typically run on a schedule (e.g., every hour)
        const events = this.batchData; // In real implementation, this would query a database
        
        const analysis = {
            period: new Date().toISOString(),
            totalEvents: events.length,
            eventsByType: this._groupByType(events),
            uniqueUsers: new Set(events.map(e => e.userId)).size,
            averageSessionDuration: this._calculateAverageDuration(events)
        };

        console.log('Historical Analysis:', analysis);
        return analysis;
    }

    _groupByType(events) {
        return events.reduce((acc, event) => {
            acc[event.type] = (acc[event.type] || 0) + 1;
            return acc;
        }, {});
    }

    _calculateAverageDuration(events) {
        const durationEvents = events.filter(e => e.data && e.data.duration);
        return durationEvents.reduce((sum, e) => sum + e.data.duration, 0) / durationEvents.length;
    }

    _setupDashboard() {
        // Update dashboard every second
        setInterval(() => {
            const stats = this.getDashboardStats();
            this._renderDashboard(stats);
        }, 1000);
    }

    _renderDashboard(stats) {
        console.clear();
        console.log('=== REAL-TIME DASHBOARD ===');
        console.log(`Active Users: ${stats.activeUsers}`);
        console.log(`Events Processed: ${stats.eventsProcessed}`);
        console.log(`Events/Second: ${stats.eventsPerSecond}`);
        console.log(`Stream Buffer: ${stats.streamBuffer}`);
        console.log(`Batch Queue: ${stats.batchQueue}`);
        console.log('---------------------------');
    }

    getDashboardStats() {
        const streamStats = this.streamProcessor.getStats();
        const batchStats = this.batchProcessor.getStats();
        
        return {
            activeUsers: this.metrics.activeUsers,
            eventsProcessed: this.metrics.eventsProcessed,
            eventsPerSecond: streamStats.itemsPerSecond,
            streamBuffer: streamStats.bufferSize,
            batchQueue: batchStats.currentBatchSize,
            processing: streamStats.currentlyProcessing
        };
    }

    delay(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Usage
// const dashboard = new RealTimeDashboard();
// dashboard.simulateUserActivity();
```

## Key Differences Summary

| Aspect | Batch Processing | Stream Processing |
|--------|------------------|-------------------|
| **Data Scope** | Entire dataset | Individual records/events |
| **Timing** | Scheduled intervals | Continuous, real-time |
| **Latency** | High (minutes to hours) | Low (milliseconds to seconds) |
| **Use Cases** | Analytics, reporting, ETL | Monitoring, alerts, real-time features |
| **Resource Usage** | Bursty, high memory | Consistent, lower memory |
| **Error Handling** | Retry entire batch | Retry individual items |
| **Complexity** | Simpler to implement | More complex, backpressure handling |

## When to Use Each

### **Use Batch Processing When:**
- You need comprehensive analysis of historical data
- Data accuracy is more important than speed
- Processing can be scheduled during off-peak hours
- You're dealing with large datasets that need aggregation
- **Examples**: Daily sales reports, monthly analytics, data migration

### **Use Stream Processing When:**
- You need immediate insights or actions
- Real-time monitoring and alerting is required
- Data arrives continuously and needs immediate processing
- Low latency is critical
- **Examples**: Fraud detection, live dashboards, real-time recommendations

## Hybrid Approach
Many modern systems use both:
- **Stream processing** for real-time alerts and immediate actions
- **Batch processing** for comprehensive analytics and data quality checks

This combination provides both real-time responsiveness and deep historical analysis.