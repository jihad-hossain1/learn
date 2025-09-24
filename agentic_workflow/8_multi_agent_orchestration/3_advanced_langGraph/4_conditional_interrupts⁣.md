# Conditional Interrupts

Conditional interrupts are control flow mechanisms that pause normal program execution based on specific conditions being met. They allow programs to respond dynamically to changing states without constant polling.

## Types of Conditional Interrupts

### 1. **Event-Driven Interrupts**
```javascript
class EventInterruptSystem {
    constructor() {
        this.interruptHandlers = new Map();
        this.isInterrupted = false;
        this.pausedState = null;
    }

    // Register conditional interrupt
    registerInterrupt(eventName, conditionFn, handlerFn, priority = 0) {
        if (!this.interruptHandlers.has(eventName)) {
            this.interruptHandlers.set(eventName, []);
        }
        
        this.interruptHandlers.get(eventName).push({
            condition: conditionFn,
            handler: handlerFn,
            priority: priority
        });
        
        // Sort by priority (higher first)
        this.interruptHandlers.get(eventName).sort((a, b) => b.priority - a.priority);
    }

    // Check for interrupts when event occurs
    async triggerEvent(eventName, eventData) {
        if (!this.interruptHandlers.has(eventName)) return false;

        const handlers = this.interruptHandlers.get(eventName);
        
        for (const { condition, handler } of handlers) {
            if (await condition(eventData)) {
                await this.executeInterrupt(handler, eventData);
                return true; // Interrupt handled
            }
        }
        return false; // No interrupt conditions met
    }

    async executeInterrupt(handler, eventData) {
        this.isInterrupted = true;
        console.log(`🔄 Interrupt triggered: ${handler.name}`);
        
        try {
            await handler(eventData);
        } finally {
            this.isInterrupted = false;
        }
    }
}

// Usage Example
const interruptSystem = new EventInterruptSystem();

// Define interrupt condition: CPU usage > 80%
const highCpuCondition = (metrics) => metrics.cpuUsage > 80;
const highCpuHandler = async (metrics) => {
    console.log(`🚨 High CPU detected: ${metrics.cpuUsage}%`);
    // Scale down operations or trigger cooling
    await this.scaleDownOperations();
};

interruptSystem.registerInterrupt('systemMetrics', highCpuCondition, highCpuHandler, 10);
```

### 2. **State-Based Conditional Interrupts**
```javascript
class StateInterruptController {
    constructor() {
        this.states = new Map();
        this.currentState = 'normal';
        this.stateTransitions = new Map();
        this.interruptQueue = [];
    }

    addState(stateName, entryHandler, exitHandler) {
        this.states.set(stateName, {
            entry: entryHandler,
            exit: exitHandler,
            interrupts: []
        });
    }

    // Register state-dependent interrupt
    registerStateInterrupt(stateName, conditionFn, interruptFn) {
        if (!this.states.has(stateName)) {
            this.states.set(stateName, { interrupts: [] });
        }
        
        this.states.get(stateName).interrupts.push({
            condition: conditionFn,
            interrupt: interruptFn
        });
    }

    async changeState(newState, data = {}) {
        const currentStateObj = this.states.get(this.currentState);
        if (currentStateObj && currentStateObj.exit) {
            await currentStateObj.exit(data);
        }

        this.currentState = newState;
        
        const newStateObj = this.states.get(newState);
        if (newStateObj && newStateObj.entry) {
            await newStateObj.entry(data);
        }

        console.log(`State changed: ${this.currentState}`);
    }

    // Check for state-based interrupts
    async checkStateInterrupts(context) {
        const currentState = this.states.get(this.currentState);
        if (!currentState || !currentState.interrupts) return false;

        for (const { condition, interrupt } of currentState.interrupts) {
            if (await condition(context)) {
                await interrupt(context);
                return true;
            }
        }
        return false;
    }
}

// Usage Example
const stateController = new StateInterruptController();

// Define states
stateController.addState(
    'processing',
    () => console.log('Entered processing state'),
    () => console.log('Exited processing state')
);

stateController.addState(
    'paused',
    () => console.log('Entered paused state'),
    () => console.log('Exited paused state')
);

// Register interrupt: If memory usage > 90% while processing, pause
stateController.registerStateInterrupt(
    'processing',
    (context) => context.memoryUsage > 90,
    async (context) => {
        console.log('🔄 Interrupt: High memory usage, pausing processing');
        await stateController.changeState('paused', context);
    }
);
```

### 3. **Timeout-Based Conditional Interrupts**
```javascript
class TimeoutInterruptManager {
    constructor() {
        this.timeouts = new Map();
        this.intervalChecks = new Map();
        this.executionContexts = new Map();
    }

    // Execute function with timeout interrupt
    async executeWithInterrupt(operationId, operationFn, timeoutMs, interruptHandler) {
        return new Promise(async (resolve, reject) => {
            const timeoutId = setTimeout(() => {
                this.handleTimeoutInterrupt(operationId, interruptHandler, reject);
            }, timeoutMs);

            this.timeouts.set(operationId, timeoutId);
            this.executionContexts.set(operationId, { resolve, reject });

            try {
                const result = await operationFn();
                clearTimeout(timeoutId);
                this.timeouts.delete(operationId);
                this.executionContexts.delete(operationId);
                resolve(result);
            } catch (error) {
                clearTimeout(timeoutId);
                this.timeouts.delete(operationId);
                this.executionContexts.delete(operationId);
                reject(error);
            }
        });
    }

    async handleTimeoutInterrupt(operationId, interruptHandler, reject) {
        console.log(`⏰ Timeout interrupt for operation: ${operationId}`);
        
        try {
            await interruptHandler();
            reject(new Error(`Operation ${operationId} interrupted by timeout`));
        } catch (error) {
            reject(new Error(`Interrupt handler failed: ${error.message}`));
        }
    }

    // Conditional timeout based on external factors
    async conditionalTimeout(operationId, operationFn, conditionFn, checkInterval = 100) {
        return new Promise(async (resolve, reject) => {
            const intervalId = setInterval(async () => {
                if (await conditionFn()) {
                    clearInterval(intervalId);
                    this.intervalChecks.delete(operationId);
                    reject(new Error(`Operation ${operationId} interrupted by condition`));
                }
            }, checkInterval);

            this.intervalChecks.set(operationId, intervalId);

            try {
                const result = await operationFn();
                clearInterval(intervalId);
                this.intervalChecks.delete(operationId);
                resolve(result);
            } catch (error) {
                clearInterval(intervalId);
                this.intervalChecks.delete(operationId);
                reject(error);
            }
        });
    }
}

// Usage Example
const timeoutManager = new TimeoutInterruptManager();

// Timeout interrupt example
const longRunningTask = async () => {
    await new Promise(resolve => setTimeout(resolve, 5000)); // 5-second task
    return "Task completed";
};

const timeoutHandler = async () => {
    console.log("Cleaning up interrupted task...");
    // Cleanup resources
};

// Execute with 2-second timeout
timeoutManager.executeWithInterrupt(
    'task-1', 
    longRunningTask, 
    2000, 
    timeoutHandler
).catch(error => console.error(error.message));
```

### 4. **Resource-Based Conditional Interrupts**
```javascript
class ResourceInterruptMonitor {
    constructor() {
        this.resourceThresholds = new Map();
        this.resourceMonitors = new Map();
        this.interruptCallbacks = new Map();
    }

    setResourceThreshold(resourceType, threshold, interruptCallback) {
        this.resourceThresholds.set(resourceType, threshold);
        this.interruptCallbacks.set(resourceType, interruptCallback);
        
        this.startMonitoring(resourceType);
    }

    startMonitoring(resourceType) {
        const monitorInterval = setInterval(() => {
            this.checkResourceUsage(resourceType);
        }, 1000); // Check every second

        this.resourceMonitors.set(resourceType, monitorInterval);
    }

    async checkResourceUsage(resourceType) {
        const currentUsage = await this.getCurrentResourceUsage(resourceType);
        const threshold = this.resourceThresholds.get(resourceType);
        
        if (currentUsage >= threshold) {
            const callback = this.interruptCallbacks.get(resourceType);
            if (callback) {
                console.log(`🚨 Resource interrupt: ${resourceType} at ${currentUsage}%`);
                await callback(currentUsage);
            }
        }
    }

    async getCurrentResourceUsage(resourceType) {
        // Simulate resource monitoring
        switch(resourceType) {
            case 'memory':
                return process.memoryUsage().heapUsed / process.memoryUsage().heapTotal * 100;
            case 'cpu':
                return Math.random() * 100; // Simulated CPU usage
            case 'disk':
                return Math.random() * 100; // Simulated disk usage
            default:
                return 0;
        }
    }

    stopMonitoring(resourceType) {
        const interval = this.resourceMonitors.get(resourceType);
        if (interval) {
            clearInterval(interval);
            this.resourceMonitors.delete(resourceType);
        }
    }
}

// Usage Example
const resourceMonitor = new ResourceInterruptMonitor();

// Set memory threshold interrupt
resourceMonitor.setResourceThreshold(
    'memory',
    85, // 85% threshold
    async (usage) => {
        console.log(`Memory usage critical: ${usage.toFixed(2)}%`);
        // Trigger garbage collection or scale down
        await this.forceGarbageCollection();
    }
);
```

### 5. **Priority-Based Interrupt System**
```javascript
class PriorityInterruptController {
    constructor() {
        this.interruptQueue = [];
        this.currentPriority = 0;
        this.isHandlingInterrupt = false;
    }

    // Add interrupt with priority
    requestInterrupt(priority, conditionFn, handlerFn, context = {}) {
        this.interruptQueue.push({
            priority,
            condition: conditionFn,
            handler: handlerFn,
            context,
            timestamp: Date.now()
        });

        // Sort by priority (higher numbers = higher priority)
        this.interruptQueue.sort((a, b) => b.priority - a.priority);
        
        this.processInterruptQueue();
    }

    async processInterruptQueue() {
        if (this.isHandlingInterrupt || this.interruptQueue.length === 0) {
            return;
        }

        this.isHandlingInterrupt = true;

        while (this.interruptQueue.length > 0) {
            const interrupt = this.interruptQueue[0];
            
            // Check if condition is met
            if (await interrupt.condition(interrupt.context)) {
                // Remove from queue and execute
                this.interruptQueue.shift();
                
                console.log(`🎯 Executing interrupt with priority ${interrupt.priority}`);
                await interrupt.handler(interrupt.context);
            } else {
                // Condition not met, skip for now
                this.interruptQueue.shift();
            }
        }

        this.isHandlingInterrupt = false;
    }

    // Nested interrupt handling
    async executeWithNestedInterrupts(mainTask, interruptConfigs) {
        const originalQueue = [...this.interruptQueue];
        this.interruptQueue = [];

        // Register temporary interrupts
        interruptConfigs.forEach(config => {
            this.requestInterrupt(config.priority, config.condition, config.handler);
        });

        try {
            const result = await mainTask();
            await this.processInterruptQueue(); // Handle any pending interrupts
            return result;
        } finally {
            // Restore original interrupt queue
            this.interruptQueue = originalQueue;
        }
    }
}

// Usage Example
const priorityController = new PriorityInterruptController();

// High-priority emergency interrupt
priorityController.requestInterrupt(
    100, // Highest priority
    (context) => context.systemTemperature > 90,
    async (context) => {
        console.log("🔥 EMERGENCY: System overheating!");
        await this.emergencyShutdown();
    }
);

// Medium-priority interrupt
priorityController.requestInterrupt(
    50,
    (context) => context.networkLatency > 1000,
    async (context) => {
        console.log("🌐 High latency detected, switching to fallback");
        await this.switchToFallbackServer();
    }
);
```

### 6. **Conditional Interrupt Patterns for Distributed Systems**
```javascript
class DistributedInterruptCoordinator {
    constructor(nodeId, clusterNodes) {
        this.nodeId = nodeId;
        this.clusterNodes = clusterNodes;
        this.globalConditions = new Map();
        this.consensusThreshold = Math.floor(clusterNodes.length / 2) + 1;
    }

    // Register global condition that requires consensus
    async registerGlobalCondition(conditionId, conditionFn, interruptHandler) {
        this.globalConditions.set(conditionId, {
            condition: conditionFn,
            handler: interruptHandler,
            votes: new Map(),
            isTriggered: false
        });

        // Notify other nodes about this condition
        await this.broadcastConditionRegistration(conditionId);
    }

    // Check condition across cluster
    async checkGlobalCondition(conditionId, localContext) {
        const condition = this.globalConditions.get(conditionId);
        if (!condition || condition.isTriggered) return false;

        // Check local condition
        const localResult = await condition.condition(localContext);
        
        if (localResult) {
            // Broadcast vote to other nodes
            const votes = await this.collectClusterVotes(conditionId, localContext);
            
            if (votes >= this.consensusThreshold) {
                condition.isTriggered = true;
                await condition.handler(localContext);
                return true;
            }
        }
        
        return false;
    }

    async collectClusterVotes(conditionId, context) {
        const votePromises = this.clusterNodes
            .filter(node => node !== this.nodeId)
            .map(async node => {
                try {
                    const response = await this.sendVoteRequest(node, conditionId, context);
                    return response.votesFor ? 1 : 0;
                } catch (error) {
                    console.error(`Failed to get vote from ${node}:`, error);
                    return 0; // Node failure counts as no vote
                }
            });

        const votes = await Promise.all(votePromises);
        const totalVotes = votes.reduce((sum, vote) => sum + vote, 1); // +1 for self
        
        return totalVotes;
    }
}
```

## Key Benefits of Conditional Interrupts

1. **Efficiency**: Avoid constant polling by triggering only when conditions met
2. **Responsiveness**: Immediate reaction to important events
3. **Resource Management**: Prevent resource exhaustion
4. **Error Prevention**: Proactive handling of potential issues
5. **System Stability**: Maintain operation within safe parameters

## Common Use Cases

- **Operating Systems**: Process scheduling, I/O operations
- **Embedded Systems**: Hardware monitoring, sensor triggers
- **Distributed Systems**: Consensus mechanisms, failure detection
- **Real-time Systems**: Deadline monitoring, quality of service
- **Game Development**: AI behavior, event triggers
- **Network Systems**: Congestion control, security monitoring

Conditional interrupts provide a powerful mechanism for building responsive, efficient, and robust systems that can adapt dynamically to changing conditions.