# Checkpointing Basics in LangGraph

Checkpointing is a critical feature that allows you to save and restore the state of your graph execution. This enables persistence, recovery, long-running workflows, and debugging capabilities.

## Core Concepts

### **What is Checkpointing?**
- **State snapshots**: Save the complete state of graph execution at specific points
- **Persistence**: Store state to durable storage (disk, database, cloud)
- **Recovery**: Resume execution from saved points after failures or restarts
- **Long-running workflows**: Handle processes that run for hours, days, or longer

### **Key Benefits**
- **Fault tolerance**: Continue after crashes or errors
- **Debugging**: Inspect state at various execution points
- **Cost efficiency**: Don't lose progress on expensive operations
- **Time travel**: Replay execution from any checkpoint

---

## JavaScript Examples

### Example 1: Basic Checkpointing System

```javascript
class BasicCheckpointer {
    constructor(storageBackend = new MemoryStorage()) {
        this.storage = storageBackend;
        this.checkpoints = new Map();
        this.config = {
            autoSave: true,
            saveInterval: 1000, // ms
            maxCheckpoints: 100,
            compression: false
        };
    }

    // Create a checkpoint of current state
    async createCheckpoint(graphId, state, metadata = {}) {
        const checkpointId = this.generateCheckpointId();
        const timestamp = Date.now();
        
        const checkpoint = {
            id: checkpointId,
            graphId,
            state: this.deepClone(state),
            metadata: {
                timestamp,
                node: metadata.node || 'unknown',
                executionStep: metadata.step || 0,
                ...metadata
            },
            version: '1.0'
        };

        // Apply compression if enabled
        if (this.config.compression) {
            checkpoint.state = this.compressState(checkpoint.state);
        }

        // Store checkpoint
        await this.storage.save(checkpointId, checkpoint);
        this.checkpoints.set(checkpointId, checkpoint);

        // Manage checkpoint limits
        await this.cleanupOldCheckpoints(graphId);

        console.log(`💾 Checkpoint created: ${checkpointId} at node ${metadata.node}`);
        return checkpointId;
    }

    // Load a checkpoint by ID
    async loadCheckpoint(checkpointId) {
        let checkpoint = this.checkpoints.get(checkpointId);
        
        if (!checkpoint) {
            checkpoint = await this.storage.load(checkpointId);
            if (checkpoint) {
                this.checkpoints.set(checkpointId, checkpoint);
            }
        }

        if (!checkpoint) {
            throw new Error(`Checkpoint not found: ${checkpointId}`);
        }

        // Decompress if needed
        if (this.config.compression && checkpoint.state) {
            checkpoint.state = this.decompressState(checkpoint.state);
        }

        return checkpoint;
    }

    // Resume execution from checkpoint
    async resumeFromCheckpoint(checkpointId, graph) {
        const checkpoint = await this.loadCheckpoint(checkpointId);
        
        console.log(`🔁 Resuming from checkpoint: ${checkpointId}`);
        console.log(`   Node: ${checkpoint.metadata.node}`);
        console.log(`   Timestamp: ${new Date(checkpoint.metadata.timestamp).toISOString()}`);
        
        // Resume graph execution from the checkpoint state
        return await graph.executeFromState(checkpoint.state, checkpoint.metadata.node);
    }

    // List checkpoints for a graph
    async listCheckpoints(graphId, options = {}) {
        const allCheckpoints = await this.storage.list();
        const graphCheckpoints = allCheckpoints.filter(cp => cp.graphId === graphId);
        
        // Apply filters
        let filtered = graphCheckpoints;
        
        if (options.since) {
            filtered = filtered.filter(cp => cp.metadata.timestamp >= options.since);
        }
        
        if (options.node) {
            filtered = filtered.filter(cp => cp.metadata.node === options.node);
        }
        
        return filtered.sort((a, b) => b.metadata.timestamp - a.metadata.timestamp);
    }

    // Clean up old checkpoints
    async cleanupOldCheckpoints(graphId) {
        const checkpoints = await this.listCheckpoints(graphId);
        
        if (checkpoints.length > this.config.maxCheckpoints) {
            const toDelete = checkpoints.slice(this.config.maxCheckpoints);
            
            for (const checkpoint of toDelete) {
                await this.storage.delete(checkpoint.id);
                this.checkpoints.delete(checkpoint.id);
                console.log(`🗑️ Deleted old checkpoint: ${checkpoint.id}`);
            }
        }
    }

    // Utility methods
    generateCheckpointId() {
        return `cp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    }

    deepClone(obj) {
        return JSON.parse(JSON.stringify(obj));
    }

    compressState(state) {
        // Simple compression - in real implementation use proper compression
        return Buffer.from(JSON.stringify(state)).toString('base64');
    }

    decompressState(compressed) {
        return JSON.parse(Buffer.from(compressed, 'base64').toString());
    }

    // Get checkpoint statistics
    async getStats(graphId) {
        const checkpoints = await this.listCheckpoints(graphId);
        
        return {
            totalCheckpoints: checkpoints.length,
            firstCheckpoint: checkpoints[checkpoints.length - 1]?.metadata.timestamp,
            lastCheckpoint: checkpoints[0]?.metadata.timestamp,
            nodesWithCheckpoints: [...new Set(checkpoints.map(cp => cp.metadata.node))],
            totalSize: checkpoints.reduce((sum, cp) => sum + JSON.stringify(cp).length, 0)
        };
    }
}

// Storage backends
class MemoryStorage {
    constructor() {
        this.data = new Map();
    }

    async save(key, value) {
        this.data.set(key, value);
    }

    async load(key) {
        return this.data.get(key);
    }

    async delete(key) {
        this.data.delete(key);
    }

    async list() {
        return Array.from(this.data.values());
    }
}

class FileStorage {
    constructor(directory = './checkpoints') {
        this.directory = directory;
        this.fs = require('fs').promises;
    }

    async save(key, value) {
        const filename = `${this.directory}/${key}.json`;
        await this.fs.writeFile(filename, JSON.stringify(value, null, 2));
    }

    async load(key) {
        try {
            const filename = `${this.directory}/${key}.json`;
            const data = await this.fs.readFile(filename, 'utf8');
            return JSON.parse(data);
        } catch (error) {
            return null;
        }
    }

    async delete(key) {
        try {
            const filename = `${this.directory}/${key}.json`;
            await this.fs.unlink(filename);
        } catch (error) {
            // File might not exist
        }
    }

    async list() {
        try {
            const files = await this.fs.readdir(this.directory);
            const checkpoints = [];
            
            for (const file of files) {
                if (file.endsWith('.json')) {
                    const data = await this.load(file.replace('.json', ''));
                    if (data) checkpoints.push(data);
                }
            }
            
            return checkpoints;
        } catch (error) {
            return [];
        }
    }
}
```

### Example 2: Checkpoint-Aware Graph Execution

```javascript
class CheckpointAwareGraph {
    constructor(checkpointer) {
        this.checkpointer = checkpointer;
        this.nodes = new Map();
        this.edges = new Map();
        this.checkpointConfig = {
            checkpointAtNodes: [], // Specific nodes to checkpoint
            checkpointInterval: 5, // Checkpoint every N nodes
            checkpointOnError: true,
            autoResume: true
        };
    }

    addNode(name, handler, checkpointConfig = {}) {
        this.nodes.set(name, {
            name,
            handler,
            checkpoint: checkpointConfig.checkpoint || false,
            metadata: checkpointConfig.metadata || {}
        });
        return this;
    }

    // Execute graph with checkpointing support
    async execute(initialState, options = {}) {
        const executionId = options.executionId || this.generateExecutionId();
        let currentState = { ...initialState, _executionId: executionId };
        let currentNode = options.startNode || Array.from(this.nodes.keys())[0];
        
        const executionHistory = {
            startTime: Date.now(),
            nodesVisited: [],
            checkpointsCreated: []
        };

        console.log(`🚀 Starting execution: ${executionId}`);

        // Try to resume from existing checkpoint
        if (options.resumeFromCheckpoint) {
            try {
                const resumeResult = await this.resumeFromCheckpoint(options.resumeFromCheckpoint);
                currentState = resumeResult.state;
                currentNode = resumeResult.currentNode;
                executionHistory.nodesVisited = resumeResult.history || [];
                console.log(`🔁 Resumed from checkpoint: ${options.resumeFromCheckpoint}`);
            } catch (error) {
                console.warn(`Failed to resume from checkpoint: ${error.message}`);
            }
        }

        while (currentNode) {
            try {
                executionHistory.nodesVisited.push(currentNode);
                
                // Execute current node
                const node = this.nodes.get(currentNode);
                if (!node) throw new Error(`Node not found: ${currentNode}`);

                console.log(`🔹 Executing node: ${currentNode}`);
                const result = await node.handler(currentState);
                currentState = { ...currentState, ...result };
                currentState._lastNode = currentNode;
                currentState._executionStep = executionHistory.nodesVisited.length;

                // Check if we should create a checkpoint
                if (await this.shouldCheckpoint(currentNode, currentState, executionHistory)) {
                    const checkpointId = await this.createCheckpoint(
                        executionId, 
                        currentNode, 
                        currentState, 
                        executionHistory
                    );
                    executionHistory.checkpointsCreated.push(checkpointId);
                }

                // Determine next node
                currentNode = this.getNextNode(currentNode, currentState);

            } catch (error) {
                console.error(`❌ Error at node ${currentNode}:`, error);
                
                // Create error checkpoint if configured
                if (this.checkpointConfig.checkpointOnError) {
                    await this.createCheckpoint(
                        executionId,
                        currentNode,
                        { ...currentState, _error: error.message, _failed: true },
                        executionHistory,
                        { error: true, errorMessage: error.message }
                    );
                }
                
                if (options.continueOnError) {
                    console.log('⚠️ Continuing execution despite error');
                    currentNode = this.getNextNode(currentNode, currentState);
                } else {
                    throw error;
                }
            }
        }

        executionHistory.endTime = Date.now();
        executionHistory.duration = executionHistory.endTime - executionHistory.startTime;
        
        console.log(`✅ Execution completed: ${executionId}`);
        console.log(`   Duration: ${executionHistory.duration}ms`);
        console.log(`   Nodes visited: ${executionHistory.nodesVisited.length}`);
        console.log(`   Checkpoints created: ${executionHistory.checkpointsCreated.length}`);

        return {
            finalState: currentState,
            executionHistory,
            executionId
        };
    }

    // Determine if checkpoint should be created
    async shouldCheckpoint(nodeName, state, history) {
        const node = this.nodes.get(nodeName);
        
        // Checkpoint at specific nodes
        if (node && node.checkpoint) {
            return true;
        }

        // Checkpoint interval
        if (history.nodesVisited.length % this.checkpointConfig.checkpointInterval === 0) {
            return true;
        }

        // Checkpoint based on state conditions
        if (state._requiresCheckpoint) {
            return true;
        }

        // Checkpoint if state is large (crude heuristic)
        const stateSize = JSON.stringify(state).length;
        if (stateSize > 100000) { // 100KB
            return true;
        }

        return false;
    }

    // Create a checkpoint
    async createCheckpoint(executionId, nodeName, state, history, metadata = {}) {
        const checkpointMetadata = {
            executionId,
            node: nodeName,
            step: history.nodesVisited.length,
            timestamp: Date.now(),
            nodesVisited: [...history.nodesVisited],
            ...metadata
        };

        return await this.checkpointer.createCheckpoint(
            executionId,
            state,
            checkpointMetadata
        );
    }

    // Resume execution from checkpoint
    async resumeFromCheckpoint(checkpointId) {
        const checkpoint = await this.checkpointer.loadCheckpoint(checkpointId);
        
        return {
            state: checkpoint.state,
            currentNode: checkpoint.metadata.node,
            history: checkpoint.metadata.nodesVisited || [],
            checkpointData: checkpoint
        };
    }

    getNextNode(currentNode, state) {
        const edges = this.edges.get(currentNode) || [];
        
        for (const edge of edges) {
            if (edge.condition(state)) {
                return edge.target;
            }
        }
        
        return null;
    }

    generateExecutionId() {
        return `exec_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    }

    // Add edge with condition
    addEdge(source, target, condition) {
        if (!this.edges.has(source)) {
            this.edges.set(source, []);
        }
        this.edges.get(source).push({ source, target, condition });
        return this;
    }
}

// Example usage: Long-running workflow with checkpointing
class DocumentProcessingWorkflow {
    constructor() {
        this.checkpointer = new BasicCheckpointer(new MemoryStorage());
        this.graph = new CheckpointAwareGraph(this.checkpointer);
        this.setupWorkflow();
    }

    setupWorkflow() {
        // Define nodes with checkpointing configuration
        this.graph
            .addNode('start', this.startProcessing.bind(this), {
                checkpoint: true,
                metadata: { description: 'Start of document processing' }
            })
            .addNode('validate_format', this.validateFormat.bind(this), {
                checkpoint: false // Don't checkpoint here - quick operation
            })
            .addNode('extract_text', this.extractText.bind(this), {
                checkpoint: true,
                metadata: { description: 'After text extraction - large state' }
            })
            .addNode('analyze_content', this.analyzeContent.bind(this), {
                checkpoint: true,
                metadata: { description: 'After content analysis' }
            })
            .addNode('generate_summary', this.generateSummary.bind(this), {
                checkpoint: true,
                metadata: { description: 'After summary generation' }
            })
            .addNode('save_results', this.saveResults.bind(this), {
                checkpoint: true,
                metadata: { description: 'Final results checkpoint' }
            });

        // Define edges
        this.graph
            .addEdge('start', 'validate_format', () => true)
            .addEdge('validate_format', 'extract_text', state => state.isValid)
            .addEdge('extract_text', 'analyze_content', state => state.textExtracted)
            .addEdge('analyze_content', 'generate_summary', state => state.analysisComplete)
            .addEdge('generate_summary', 'save_results', state => state.summaryGenerated);

        // Configure checkpointing
        this.graph.checkpointConfig = {
            checkpointInterval: 2, // Checkpoint every 2 nodes
            checkpointOnError: true,
            autoResume: true
        };
    }

    // Node implementations
    async startProcessing(state) {
        console.log('📄 Starting document processing...');
        await this.simulateWork(100);
        return {
            documentId: state.documentId || `doc_${Date.now()}`,
            startedAt: new Date().toISOString(),
            isValid: true
        };
    }

    async validateFormat(state) {
        console.log('🔍 Validating document format...');
        await this.simulateWork(50);
        return {
            formatValid: Math.random() > 0.1, // 90% valid
            validationTime: new Date().toISOString()
        };
    }

    async extractText(state) {
        console.log('📝 Extracting text from document...');
        // Simulate long-running operation
        await this.simulateWork(2000);
        
        const largeText = "This is a simulated large document text. ".repeat(1000);
        return {
            extractedText: largeText,
            textLength: largeText.length,
            textExtracted: true,
            extractionTime: new Date().toISOString()
        };
    }

    async analyzeContent(state) {
        console.log('🧠 Analyzing document content...');
        await this.simulateWork(1500);
        return {
            analysis: {
                sentiment: 'positive',
                keyTopics: ['technology', 'innovation', 'future'],
                wordCount: state.textLength,
                complexity: 'high'
            },
            analysisComplete: true,
            analyzedAt: new Date().toISOString()
        };
    }

    async generateSummary(state) {
        console.log('📋 Generating summary...');
        await this.simulateWork(800);
        return {
            summary: `Document about ${state.analysis.keyTopics.join(', ')} with ${state.analysis.wordCount} words`,
            summaryGenerated: true,
            summaryTime: new Date().toISOString()
        };
    }

    async saveResults(state) {
        console.log('💾 Saving results...');
        await this.simulateWork(300);
        return {
            saved: true,
            savedAt: new Date().toISOString(),
            finalResult: {
                documentId: state.documentId,
                summary: state.summary,
                processingTime: Date.now() - new Date(state.startedAt).getTime()
            }
        };
    }

    simulateWork(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    async processDocument(documentId) {
        const initialState = { documentId };
        return await this.graph.execute(initialState, {
            executionId: `doc_processing_${documentId}`,
            continueOnError: false
        });
    }

    // Demonstrate checkpoint recovery
    async demonstrateRecovery() {
        console.log('🧪 DEMONSTRATING CHECKPOINT RECOVERY\n');
        
        // Start a processing job
        const documentId = 'test_doc_123';
        console.log('1. Starting document processing...');
        
        // Simulate interruption after 2 seconds
        setTimeout(async () => {
            console.log('\n⚡ Simulating system interruption...');
            
            // List available checkpoints
            const checkpoints = await this.checkpointer.listCheckpoints(`doc_processing_${documentId}`);
            if (checkpoints.length > 0) {
                const latestCheckpoint = checkpoints[0];
                console.log(`📋 Found checkpoint: ${latestCheckpoint.id} at node ${latestCheckpoint.metadata.node}`);
                
                // Resume from checkpoint
                console.log('2. Resuming from checkpoint...');
                await this.graph.execute(
                    {}, 
                    { 
                        resumeFromCheckpoint: latestCheckpoint.id,
                        executionId: `doc_processing_${documentId}_resumed`
                    }
                );
            }
        }, 2000);

        // Start initial processing
        await this.processDocument(documentId);
    }
}

// Demonstration
async function demoCheckpointing() {
    const workflow = new DocumentProcessingWorkflow();
    
    // Run the recovery demonstration
    await workflow.demonstrateRecovery();
}

// demoCheckpointing();
```

### Example 3: Advanced Checkpoint Strategies

```javascript
class AdvancedCheckpointStrategies {
    constructor() {
        this.checkpointer = new BasicCheckpointer(new MemoryStorage());
        this.strategies = new Map();
        this.setupStrategies();
    }

    setupStrategies() {
        // Time-based strategy
        this.strategies.set('time_based', {
            shouldCheckpoint: (state, history, config) => {
                const now = Date.now();
                const lastCheckpoint = history.lastCheckpointTime || history.startTime;
                return now - lastCheckpoint >= config.interval;
            },
            config: { interval: 30000 } // 30 seconds
        });

        // Size-based strategy
        this.strategies.set('size_based', {
            shouldCheckpoint: (state, history, config) => {
                const stateSize = JSON.stringify(state).length;
                return stateSize >= config.threshold;
            },
            config: { threshold: 50000 } // 50KB
        });

        // Progress-based strategy
        this.strategies.set('progress_based', {
            shouldCheckpoint: (state, history, config) => {
                const progress = history.nodesVisited.length / config.totalNodes;
                return progress >= history.lastProgress + config.progressInterval;
            },
            config: { progressInterval: 0.1, totalNodes: 10 } // Every 10% progress
        });

        // Cost-based strategy (for expensive operations)
        this.strategies.set('cost_based', {
            shouldCheckpoint: (state, history, config) => {
                const expensiveOperations = ['llm_call', 'api_call', 'large_processing'];
                const lastNode = history.nodesVisited[history.nodesVisited.length - 1];
                return expensiveOperations.some(op => lastNode.includes(op));
            },
            config: { expensiveNodes: ['llm', 'api', 'process'] }
        });

        // Adaptive strategy combining multiple factors
        this.strategies.set('adaptive', {
            shouldCheckpoint: (state, history, config) => {
                const factors = [];
                
                // Time factor
                const timeSinceLast = Date.now() - (history.lastCheckpointTime || history.startTime);
                factors.push(timeSinceLast >= config.timeThreshold ? 1 : 0);
                
                // Size factor
                const stateSize = JSON.stringify(state).length;
                factors.push(stateSize >= config.sizeThreshold ? 1 : 0);
                
                // Progress factor
                const progress = history.nodesVisited.length / config.estimatedTotalNodes;
                factors.push(progress - (history.lastProgress || 0) >= config.progressThreshold ? 1 : 0);
                
                // Weighted score
                const score = factors.reduce((sum, factor, index) => 
                    sum + (factor * config.weights[index]), 0);
                
                return score >= config.decisionThreshold;
            },
            config: {
                timeThreshold: 60000, // 1 minute
                sizeThreshold: 100000, // 100KB
                progressThreshold: 0.2, // 20%
                estimatedTotalNodes: 10,
                weights: [0.3, 0.4, 0.3], // Size is most important
                decisionThreshold: 0.6
            }
        });
    }

    // Smart checkpointing with strategy selection
    async smartCheckpoint(executionId, nodeName, state, history, strategyName = 'adaptive') {
        const strategy = this.strategies.get(strategyName);
        if (!strategy) {
            throw new Error(`Unknown strategy: ${strategyName}`);
        }

        const shouldCheckpoint = strategy.shouldCheckpoint(state, history, strategy.config);
        
        if (shouldCheckpoint) {
            console.log(`🎯 Smart checkpoint triggered by ${strategyName} strategy`);
            const checkpointId = await this.checkpointer.createCheckpoint(
                executionId,
                state,
                {
                    node: nodeName,
                    step: history.nodesVisited.length,
                    timestamp: Date.now(),
                    strategy: strategyName,
                    reason: this.getCheckpointReason(state, history, strategyName)
                }
            );

            // Update history
            history.lastCheckpointTime = Date.now();
            history.lastProgress = history.nodesVisited.length / (strategy.config.estimatedTotalNodes || 10);
            history.checkpointsCreated.push(checkpointId);

            return checkpointId;
        }

        return null;
    }

    getCheckpointReason(state, history, strategyName) {
        switch (strategyName) {
            case 'time_based':
                return `Time threshold exceeded: ${Date.now() - (history.lastCheckpointTime || history.startTime)}ms`;
            case 'size_based':
                return `State size: ${JSON.stringify(state).length} bytes`;
            case 'progress_based':
                const progress = (history.nodesVisited.length / 10) * 100;
                return `Progress: ${progress.toFixed(1)}%`;
            case 'adaptive':
                return 'Multiple factors triggered checkpoint';
            default:
                return 'Strategy-based checkpoint';
        }
    }

    // Checkpoint optimization - incremental checkpoints
    async createIncrementalCheckpoint(executionId, previousCheckpointId, state, metadata) {
        // Load previous checkpoint to compute differences
        const previous = await this.checkpointer.loadCheckpoint(previousCheckpointId);
        
        // Compute delta (simplified - in reality use proper diff algorithm)
        const delta = this.computeStateDelta(previous.state, state);
        
        const incrementalCheckpoint = {
            id: this.checkpointer.generateCheckpointId(),
            graphId: executionId,
            type: 'incremental',
            baseCheckpoint: previousCheckpointId,
            delta: delta,
            metadata: {
                ...metadata,
                fullStateSize: JSON.stringify(state).length,
                deltaSize: JSON.stringify(delta).length,
                compressionRatio: (JSON.stringify(delta).length / JSON.stringify(state).length) * 100
            }
        };

        await this.checkpointer.storage.save(incrementalCheckpoint.id, incrementalCheckpoint);
        return incrementalCheckpoint.id;
    }

    computeStateDelta(previousState, currentState) {
        const delta = {};
        
        // Simple delta computation - only changed fields
        for (const key in currentState) {
            if (JSON.stringify(previousState[key]) !== JSON.stringify(currentState[key])) {
                delta[key] = currentState[key];
            }
        }
        
        return delta;
    }

    // Restore state from incremental checkpoints
    async restoreFromIncremental(checkpointId) {
        const checkpoint = await this.checkpointer.loadCheckpoint(checkpointId);
        
        if (checkpoint.type !== 'incremental') {
            return checkpoint.state;
        }

        // Reconstruct full state by applying deltas to base
        const baseState = await this.restoreFromIncremental(checkpoint.baseCheckpoint);
        const fullState = this.applyDelta(baseState, checkpoint.delta);
        
        return fullState;
    }

    applyDelta(baseState, delta) {
        return { ...baseState, ...delta };
    }

    // Checkpoint analysis and cleanup
    async analyzeCheckpoints(executionId) {
        const checkpoints = await this.checkpointer.listCheckpoints(executionId);
        
        const analysis = {
            total: checkpoints.length,
            byStrategy: {},
            sizeStats: {},
            frequency: {},
            recommendations: []
        };

        checkpoints.forEach(cp => {
            // Group by strategy
            const strategy = cp.metadata.strategy || 'manual';
            analysis.byStrategy[strategy] = (analysis.byStrategy[strategy] || 0) + 1;
            
            // Size statistics
            const size = JSON.stringify(cp.state).length;
            analysis.sizeStats.total = (analysis.sizeStats.total || 0) + size;
            analysis.sizeStats.average = analysis.sizeStats.total / checkpoints.length;
            analysis.sizeStats.largest = Math.max(analysis.sizeStats.largest || 0, size);
            analysis.sizeStats.smallest = Math.min(analysis.sizeStats.smallest || Infinity, size);
        });

        // Generate recommendations
        if (analysis.byStrategy.manual && analysis.byStrategy.manual > analysis.total * 0.5) {
            analysis.recommendations.push('Consider using automated checkpointing strategies');
        }

        if (analysis.sizeStats.average > 100000) {
            analysis.recommendations.push('Consider implementing incremental checkpointing');
        }

        return analysis;
    }
}
```

### Example 4: Distributed Checkpointing for Scalability

```javascript
class DistributedCheckpointer {
    constructor(config = {}) {
        this.config = {
            storageBackend: config.storageBackend || 'memory',
            replicationFactor: config.replicationFactor || 3,
            consistencyLevel: config.consistencyLevel || 'quorum', // quorum, all, one
            timeout: config.timeout || 5000,
            ...config
        };

        this.storageNodes = new Map();
        this.checkpointRegistry = new Map();
        this.pendingOperations = new Map();
    }

    // Add storage node to the cluster
    addStorageNode(nodeId, storageInstance) {
        this.storageNodes.set(nodeId, {
            id: nodeId,
            storage: storageInstance,
            status: 'healthy',
            lastHeartbeat: Date.now()
        });

        console.log(`➕ Added storage node: ${nodeId}`);
    }

    // Distributed checkpoint save
    async createCheckpoint(graphId, state, metadata = {}) {
        const checkpointId = this.generateCheckpointId();
        const checkpoint = {
            id: checkpointId,
            graphId,
            state: this.compressState(state),
            metadata: { ...metadata, timestamp: Date.now() },
            version: '1.0'
        };

        // Save to multiple nodes based on replication factor
        const savePromises = [];
        const nodesToUse = this.selectStorageNodes(this.config.replicationFactor);

        for (const node of nodesToUse) {
            savePromises.push(
                this.saveToNode(node, checkpointId, checkpoint)
            );
        }

        // Wait for required consistency level
        const results = await Promise.allSettled(savePromises);
        const successfulSaves = results.filter(r => r.status === 'fulfilled').length;

        if (this.isConsistencyMet(successfulSaves)) {
            // Register checkpoint in registry
            this.checkpointRegistry.set(checkpointId, {
                graphId,
                storageNodes: nodesToUse.map(n => n.id),
                createdAt: Date.now(),
                size: JSON.stringify(checkpoint).length
            });

            console.log(`💾 Distributed checkpoint saved: ${checkpointId} to ${successfulSaves} nodes`);
            return checkpointId;
        } else {
            throw new Error(`Failed to meet consistency requirements. Successful: ${successfulSaves}, Required: ${this.getRequiredSuccessCount()}`);
        }
    }

    // Load checkpoint with fault tolerance
    async loadCheckpoint(checkpointId) {
        const registryEntry = this.checkpointRegistry.get(checkpointId);
        if (!registryEntry) {
            throw new Error(`Checkpoint not found in registry: ${checkpointId}`);
        }

        // Try to load from storage nodes
        const loadPromises = registryEntry.storageNodes.map(nodeId => 
            this.loadFromNode(nodeId, checkpointId)
        );

        const results = await Promise.allSettled(loadPromises);
        const successfulLoads = results.filter(r => r.status === 'fulfilled' && r.value !== null);

        if (successfulLoads.length > 0) {
            // Use the first successful load
            const checkpoint = successfulLoads[0].value;
            checkpoint.state = this.decompressState(checkpoint.state);
            
            // If we have multiple successful loads, verify consistency
            if (successfulLoads.length > 1) {
                const isConsistent = this.verifyConsistency(successfulLoads.map(r => r.value));
                if (!isConsistent) {
                    console.warn('⚠️ Checkpoint data inconsistency detected');
                }
            }

            return checkpoint;
        } else {
            throw new Error(`Failed to load checkpoint from any storage node: ${checkpointId}`);
        }
    }

    // Select storage nodes for replication
    selectStorageNodes(count) {
        const healthyNodes = Array.from(this.storageNodes.values())
            .filter(node => node.status === 'healthy')
            .sort((a, b) => a.load - b.load); // Simple load-based selection

        return healthyNodes.slice(0, count);
    }

    // Save to individual storage node with timeout
    async saveToNode(node, checkpointId, checkpoint) {
        const timeoutPromise = new Promise((_, reject) => 
            setTimeout(() => reject(new Error('Save timeout')), this.config.timeout)
        );

        return Promise.race([
            node.storage.save(checkpointId, checkpoint),
            timeoutPromise
        ]);
    }

    // Load from individual storage node
    async loadFromNode(nodeId, checkpointId) {
        const node = this.storageNodes.get(nodeId);
        if (!node) throw new Error(`Storage node not found: ${nodeId}`);

        const timeoutPromise = new Promise((_, reject) => 
            setTimeout(() => reject(new Error('Load timeout')), this.config.timeout)
        );

        return Promise.race([
            node.storage.load(checkpointId),
            timeoutPromise
        ]);
    }

    // Consistency checking
    isConsistencyMet(successfulCount) {
        switch (this.config.consistencyLevel) {
            case 'all':
                return successfulCount === this.config.replicationFactor;
            case 'quorum':
                return successfulCount > Math.floor(this.config.replicationFactor / 2);
            case 'one':
                return successfulCount >= 1;
            default:
                return successfulCount > Math.floor(this.config.replicationFactor / 2);
        }
    }

    getRequiredSuccessCount() {
        switch (this.config.consistencyLevel) {
            case 'all': return this.config.replicationFactor;
            case 'quorum': return Math.floor(this.config.replicationFactor / 2) + 1;
            case 'one': return 1;
            default: return Math.floor(this.config.replicationFactor / 2) + 1;
        }
    }

    verifyConsistency(checkpoints) {
        if (checkpoints.length < 2) return true;
        
        const first = JSON.stringify(checkpoints[0]);
        return checkpoints.every(cp => JSON.stringify(cp) === first);
    }

    // Node health monitoring
    startHealthMonitoring() {
        setInterval(() => {
            this.checkNodeHealth();
        }, 30000); // Check every 30 seconds
    }

    async checkNodeHealth() {
        for (const [nodeId, node] of this.storageNodes) {
            try {
                // Simple health check - try to save and load a small test object
                const testKey = `health_check_${nodeId}_${Date.now()}`;
                const testData = { timestamp: Date.now(), node: nodeId };
                
                await node.storage.save(testKey, testData);
                const loaded = await node.storage.load(testKey);
                
                if (JSON.stringify(loaded) === JSON.stringify(testData)) {
                    node.status = 'healthy';
                    node.lastHeartbeat = Date.now();
                } else {
                    node.status = 'unhealthy';
                }
                
                await node.storage.delete(testKey);
            } catch (error) {
                node.status = 'unhealthy';
                console.error(`Health check failed for node ${nodeId}:`, error.message);
            }
        }
    }

    // Garbage collection for orphaned checkpoints
    async garbageCollect() {
        const allCheckpoints = Array.from(this.checkpointRegistry.entries());
        const orphaned = allCheckpoints.filter(([id, entry]) => 
            Date.now() - entry.createdAt > 24 * 60 * 60 * 1000 // 24 hours
        );

        for (const [checkpointId, entry] of orphaned) {
            await this.deleteCheckpoint(checkpointId);
            console.log(`🗑️ Collected orphaned checkpoint: ${checkpointId}`);
        }
    }

    async deleteCheckpoint(checkpointId) {
        const registryEntry = this.checkpointRegistry.get(checkpointId);
        if (!registryEntry) return;

        // Delete from all storage nodes
        const deletePromises = registryEntry.storageNodes.map(nodeId => {
            const node = this.storageNodes.get(nodeId);
            return node ? node.storage.delete(checkpointId) : Promise.resolve();
        });

        await Promise.allSettled(deletePromises);
        this.checkpointRegistry.delete(checkpointId);
    }

    generateCheckpointId() {
        return `dist_cp_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    }

    compressState(state) {
        // In production, use proper compression
        return state;
    }

    decompressState(state) {
        return state;
    }

    // Get cluster status
    getClusterStatus() {
        const nodes = Array.from(this.storageNodes.values());
        const healthyNodes = nodes.filter(n => n.status === 'healthy');
        
        return {
            totalNodes: nodes.length,
            healthyNodes: healthyNodes.length,
            unhealthyNodes: nodes.length - healthyNodes.length,
            totalCheckpoints: this.checkpointRegistry.size,
            consistencyLevel: this.config.consistencyLevel,
            replicationFactor: this.config.replicationFactor
        };
    }
}
```

## Key Checkpointing Patterns

### **1. Automatic Checkpointing**
```javascript
// Checkpoint at regular intervals
setInterval(() => {
    createCheckpoint(state);
}, checkpointInterval);
```

### **2. Conditional Checkpointing**
```javascript
// Checkpoint based on conditions
if (state.requiresCheckpoint || isExpensiveOperation(node)) {
    createCheckpoint(state);
}
```

### **3. Incremental Checkpointing**
```javascript
// Only save changes from previous state
const delta = computeDelta(previousState, currentState);
saveIncrementalCheckpoint(delta);
```

### **4. Distributed Checkpointing**
```