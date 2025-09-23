# Conditional Edges in LangGraph

Conditional edges are the decision-making components of a graph that determine the flow based on the current state. They enable dynamic routing, branching logic, and adaptive behavior in LangGraph applications.

## Core Concepts

### **What are Conditional Edges?**
- **Runtime decisions**: Evaluate state to choose next node
- **Branching logic**: Enable multiple execution paths
- **State-dependent routing**: Flow changes based on data
- **Graph intelligence**: Make graphs responsive and adaptive

### **Key Characteristics**
- **Boolean conditions**: Return true/false for each potential path
- **Mutually exclusive**: Only one path is typically taken
- **Fallback support**: Default paths when no conditions match
- **State inspection**: Analyze current state to make decisions

---

## JavaScript Examples

### Example 1: Basic Conditional Edge System

```javascript
class ConditionalEdgeSystem {
    constructor() {
        this.nodes = new Map();
        this.conditionalEdges = new Map();
        this.defaultEdges = new Map();
    }

    addNode(name, handler) {
        this.nodes.set(name, { name, handler });
        return this;
    }

    // Add conditional edge with a condition function
    addConditionalEdge(sourceNode, conditionFn, targetNode, metadata = {}) {
        if (!this.conditionalEdges.has(sourceNode)) {
            this.conditionalEdges.set(sourceNode, []);
        }

        this.conditionalEdges.get(sourceNode).push({
            condition: conditionFn,
            target: targetNode,
            description: metadata.description || 'Conditional edge',
            priority: metadata.priority || 1
        });

        return this;
    }

    // Add default edge (when no conditions match)
    addDefaultEdge(sourceNode, targetNode) {
        this.defaultEdges.set(sourceNode, targetNode);
        return this;
    }

    async execute(startNode, initialState) {
        let currentState = { ...initialState };
        let currentNode = startNode;
        const path = [];

        while (currentNode) {
            path.push(currentNode);
            console.log(`📍 Current node: ${currentNode}`);

            // Execute current node
            const node = this.nodes.get(currentNode);
            if (!node) throw new Error(`Node not found: ${currentNode}`);

            const result = await node.handler(currentState);
            currentState = { ...currentState, ...result };

            // Determine next node using conditional edges
            const nextNode = this._evaluateConditionalEdges(currentNode, currentState);
            console.log(`➡️ Next node: ${nextNode}`);

            if (!nextNode) break; // End of graph
            currentNode = nextNode;
        }

        return { finalState: currentState, executionPath: path };
    }

    _evaluateConditionalEdges(currentNode, state) {
        const edges = this.conditionalEdges.get(currentNode) || [];

        // Sort by priority (higher priority first)
        const sortedEdges = edges.sort((a, b) => b.priority - a.priority);

        // Evaluate conditions in priority order
        for (const edge of sortedEdges) {
            try {
                const shouldFollow = edge.condition(state);
                console.log(`🔍 Condition "${edge.description}": ${shouldFollow}`);
                
                if (shouldFollow) {
                    return edge.target;
                }
            } catch (error) {
                console.error(`Error evaluating condition: ${error.message}`);
            }
        }

        // No conditions matched, use default edge
        const defaultTarget = this.defaultEdges.get(currentNode);
        if (defaultTarget) {
            console.log(`🔀 Using default edge to: ${defaultTarget}`);
        }

        return defaultTarget;
    }

    visualize() {
        console.log('\n📊 CONDITIONAL EDGES VISUALIZATION');
        this.conditionalEdges.forEach((edges, source) => {
            console.log(`\n${source}:`);
            edges.forEach(edge => {
                console.log(`  ──▷ ${edge.target} (${edge.description}) [priority: ${edge.priority}]`);
            });
            if (this.defaultEdges.has(source)) {
                console.log(`  ──▷ ${this.defaultEdges.get(source)} [DEFAULT]`);
            }
        });
    }
}
```

### Example 2: Customer Service with Conditional Routing

```javascript
class CustomerServiceWithConditionals {
    constructor() {
        this.system = new ConditionalEdgeSystem();
        this._buildGraph();
    }

    _buildGraph() {
        // Define nodes
        this.system
            .addNode('receive_message', this._receiveMessage.bind(this))
            .addNode('classify_urgency', this._classifyUrgency.bind(this))
            .addNode('handle_urgent', this._handleUrgent.bind(this))
            .addNode('handle_normal', this._handleNormal.bind(this))
            .addNode('check_resolution', this._checkResolution.bind(this))
            .addNode('escalate', this._escalate.bind(this))
            .addNode('end_conversation', this._endConversation.bind(this));

        // Define conditional edges
        this.system
            // After receiving message, always classify urgency
            .addDefaultEdge('receive_message', 'classify_urgency')
            
            // Urgency-based routing
            .addConditionalEdge('classify_urgency', 
                state => state.urgency === 'high', 'handle_urgent',
                { description: 'High urgency messages', priority: 10 })
            
            .addConditionalEdge('classify_urgency', 
                state => state.urgency === 'normal', 'handle_normal',
                { description: 'Normal urgency messages', priority: 5 })
            
            .addConditionalEdge('classify_urgency', 
                state => state.urgency === 'low', 'handle_normal',
                { description: 'Low urgency messages', priority: 1 })
            
            // After handling, check if resolved
            .addDefaultEdge('handle_urgent', 'check_resolution')
            .addDefaultEdge('handle_normal', 'check_resolution')
            
            // Resolution checking
            .addConditionalEdge('check_resolution',
                state => state.resolved === true, 'end_conversation',
                { description: 'Issue resolved', priority: 10 })
            
            .addConditionalEdge('check_resolution',
                state => state.attempts >= 3, 'escalate',
                { description: 'Max attempts reached', priority: 8 })
            
            .addConditionalEdge('check_resolution',
                state => state.requires_human === true, 'escalate',
                { description: 'Requires human intervention', priority: 9 })
            
            .addConditionalEdge('check_resolution',
                state => state.resolved === false, 'classify_urgency',
                { description: 'Retry with reclassification', priority: 1 })
            
            // From escalate to end
            .addDefaultEdge('escalate', 'end_conversation');

        return this.system;
    }

    // Node implementations
    async _receiveMessage(state) {
        console.log(`📨 Received message: "${state.message}"`);
        return {
            received_at: new Date().toISOString(),
            attempts: (state.attempts || 0) + 1
        };
    }

    async _classifyUrgency(state) {
        const message = state.message.toLowerCase();
        let urgency = 'normal';
        
        if (message.includes('urgent') || message.includes('emergency') || message.includes('help!')) {
            urgency = 'high';
        } else if (message.includes('when you have time') || message.includes('no rush')) {
            urgency = 'low';
        }

        // Simulate AI classification
        const confidence = Math.random();
        console.log(`🔍 Urgency classification: ${urgency} (confidence: ${confidence.toFixed(2)})`);

        return { 
            urgency,
            classification_confidence: confidence,
            requires_human: confidence < 0.6 // Low confidence requires human
        };
    }

    async _handleUrgent(state) {
        console.log('🚨 Handling URGENT request immediately');
        
        // Simulate urgent handling
        const resolved = Math.random() > 0.2; // 80% success rate for urgent
        const response = resolved ? 
            "URGENT: Issue has been resolved immediately." : 
            "URGENT: Escalating to senior support.";

        return {
            response,
            resolved,
            handled_as_urgent: true,
            response_time: 'immediate'
        };
    }

    async _handleNormal(state) {
        console.log('⏱️ Handling normal priority request');
        
        // Simulate normal handling
        const resolved = Math.random() > 0.5; // 50% success rate
        const response = resolved ? 
            "Your issue has been addressed." : 
            "I need more information to help you.";

        return {
            response,
            resolved,
            handled_as_urgent: false,
            response_time: 'normal'
        };
    }

    async _checkResolution(state) {
        console.log(`✓ Checking resolution: ${state.resolved ? 'RESOLVED' : 'NOT RESOLVED'}`);
        return {
            last_check: new Date().toISOString(),
            check_result: state.resolved ? 'success' : 'needs_more_work'
        };
    }

    async _escalate(state) {
        console.log('👥 Escalating to human agent');
        return {
            escalated: true,
            escalated_at: new Date().toISOString(),
            escalation_reason: state.attempts >= 3 ? 'max_attempts' : 'low_confidence'
        };
    }

    async _endConversation(state) {
        console.log('✅ Conversation ended');
        return {
            ended: true,
            end_time: new Date().toISOString(),
            success: state.resolved === true
        };
    }

    async processMessage(message) {
        const initialState = { message };
        return await this.system.execute('receive_message', initialState);
    }
}

// Demonstration
async function demoCustomerService() {
    const service = new CustomerServiceWithConditionals();
    service.system.visualize();

    console.log('\n💬 CUSTOMER SERVICE DEMONSTRATION\n');

    // Test different scenarios
    const testCases = [
        "URGENT: My account has been hacked!",
        "Hello, when you have time can you help me?",
        "I'm having trouble logging in",
        "HELP! My payment failed and I need this fixed now!",
        "Just a general question about your services"
    ];

    for (const message of testCases) {
        console.log(`\n=== Testing: "${message}" ===`);
        const result = await service.processMessage(message);
        console.log('Execution path:', result.executionPath);
        console.log('Final state:', result.finalState);
    }
}

// demoCustomerService();
```

### Example 3: Advanced Conditional Logic Patterns

```javascript
class AdvancedConditionalPatterns {
    constructor() {
        this.system = new ConditionalEdgeSystem();
        this._buildAdvancedGraph();
    }

    _buildAdvancedGraph() {
        // Add nodes for complex workflow
        this.system
            .addNode('start', this._startNode.bind(this))
            .addNode('analyze_sentiment', this._analyzeSentiment.bind(this))
            .addNode('extract_entities', this._extractEntities.bind(this))
            .addNode('classify_intent', this._classifyIntent.bind(this))
            .addNode('generate_response', this._generateResponse.bind(this))
            .addNode('validate_response', this._validateResponse.bind(this))
            .addNode('apply_safety_check', this._applySafetyCheck.bind(this))
            .addNode('send_response', this._sendResponse.bind(this))
            .addNode('log_interaction', this._logInteraction.bind(this));

        // Complex conditional edges with multiple factors
        this.system
            .addDefaultEdge('start', 'analyze_sentiment')
            
            // Sentiment-based routing with confidence threshold
            .addConditionalEdge('analyze_sentiment',
                state => state.sentiment.score < -0.7 && state.sentiment.confidence > 0.8,
                'apply_safety_check',
                { description: 'Highly negative sentiment - safety first', priority: 20 })
            
            .addConditionalEdge('analyze_sentiment',
                state => state.sentiment.score > 0.7,
                'extract_entities',
                { description: 'Positive sentiment - proceed normally', priority: 10 })
            
            .addConditionalEdge('analyze_sentiment',
                state => Math.abs(state.sentiment.score) <= 0.3,
                'classify_intent', 
                { description: 'Neutral sentiment - focus on intent', priority: 5 })
            
            // Entity extraction conditions
            .addConditionalEdge('extract_entities',
                state => state.entities?.length > 0,
                'classify_intent',
                { description: 'Entities found - proceed to intent', priority: 10 })
            
            .addConditionalEdge('extract_entities',
                state => !state.entities || state.entities.length === 0,
                'generate_response',
                { description: 'No entities - generate direct response', priority: 5 })
            
            // Intent classification with fallbacks
            .addConditionalEdge('classify_intent',
                state => state.intent.confidence > 0.9,
                'generate_response',
                { description: 'High confidence intent', priority: 15 })
            
            .addConditionalEdge('classify_intent',
                state => state.intent.confidence <= 0.9 && state.intent.confidence > 0.6,
                'apply_safety_check',
                { description: 'Medium confidence - safety check', priority: 10 })
            
            .addConditionalEdge('classify_intent',
                state => state.intent.confidence <= 0.6,
                'log_interaction',
                { description: 'Low confidence - log for review', priority: 1 })
            
            // Response validation flow
            .addConditionalEdge('validate_response',
                state => state.response_quality >= 0.8,
                'send_response',
                { description: 'High quality response', priority: 10 })
            
            .addConditionalEdge('validate_response',
                state => state.response_quality < 0.8 && state.regeneration_attempts < 3,
                'generate_response',
                { description: 'Regenerate response', priority: 5 })
            
            .addConditionalEdge('validate_response',
                state => state.response_quality < 0.8 && state.regeneration_attempts >= 3,
                'log_interaction',
                { description: 'Max regeneration attempts', priority: 1 })
            
            // Safety check results
            .addConditionalEdge('apply_safety_check',
                state => state.safety_approved === true,
                state => state.previous_node === 'analyze_sentiment' ? 'extract_entities' : 'validate_response',
                { description: 'Safety approved - continue flow', priority: 10 })
            
            .addConditionalEdge('apply_safety_check',
                state => state.safety_approved === false,
                'log_interaction',
                { description: 'Safety violation - log and stop', priority: 20 })
            
            // Default edges for completion
            .addDefaultEdge('generate_response', 'validate_response')
            .addDefaultEdge('send_response', 'log_interaction')
            .addDefaultEdge('log_interaction', null); // End of graph

        return this.system;
    }

    // Node implementations
    async _startNode(state) {
        console.log('🚀 Starting advanced workflow');
        return {
            workflow_started: new Date().toISOString(),
            regeneration_attempts: 0,
            previous_node: 'start'
        };
    }

    async _analyzeSentiment(state) {
        console.log('😊 Analyzing sentiment...');
        // Simulate sentiment analysis
        const score = Math.random() * 2 - 1; // -1 to 1
        const confidence = 0.7 + Math.random() * 0.3; // 0.7-1.0
        
        return {
            sentiment: { score, confidence },
            requires_careful_handling: score < -0.5,
            previous_node: 'analyze_sentiment'
        };
    }

    async _extractEntities(state) {
        console.log('🔍 Extracting entities...');
        const entities = ['user', 'service', 'problem'].filter(() => Math.random() > 0.3);
        
        return {
            entities,
            entity_count: entities.length,
            previous_node: 'extract_entities'
        };
    }

    async _classifyIntent(state) {
        console.log('🎯 Classifying intent...');
        const intents = ['question', 'complaint', 'request', 'feedback'];
        const intent = intents[Math.floor(Math.random() * intents.length)];
        const confidence = 0.5 + Math.random() * 0.5; // 0.5-1.0
        
        return {
            intent: { type: intent, confidence },
            requires_specific_handling: confidence > 0.8,
            previous_node: 'classify_intent'
        };
    }

    async _generateResponse(state) {
        console.log('💭 Generating response...');
        const quality = 0.6 + Math.random() * 0.4; // 0.6-1.0
        const attempts = (state.regeneration_attempts || 0) + 1;
        
        return {
            response: `Generated response for ${state.intent?.type || 'unknown'} intent`,
            response_quality: quality,
            regeneration_attempts: attempts,
            previous_node: 'generate_response'
        };
    }

    async _validateResponse(state) {
        console.log('✓ Validating response quality...');
        const passes = state.response_quality > 0.7;
        
        return {
            validation_passed: passes,
            validation_score: state.response_quality,
            previous_node: 'validate_response'
        };
    }

    async _applySafetyCheck(state) {
        console.log('🛡️ Applying safety checks...');
        const approved = Math.random() > 0.1; // 90% approval rate
        
        return {
            safety_approved: approved,
            safety_check_performed: true,
            previous_node: 'apply_safety_check'
        };
    }

    async _sendResponse(state) {
        console.log('📤 Sending response to user');
        return {
            response_sent: true,
            sent_at: new Date().toISOString(),
            previous_node: 'send_response'
        };
    }

    async _logInteraction(state) {
        console.log('📝 Logging interaction');
        return {
            logged: true,
            log_timestamp: new Date().toISOString(),
            final_state: state
        };
    }

    async processInput(input) {
        const initialState = { input };
        return await this.system.execute('start', initialState);
    }
}

// Demonstration of advanced conditionals
async function demoAdvancedConditionals() {
    const advanced = new AdvancedConditionalPatterns();
    
    console.log('🧠 ADVANCED CONDITIONAL PATTERNS DEMO\n');
    
    const inputs = [
        "I'm really angry about your service!",
        "Can you help me with my account?",
        "This is the worst experience ever!",
        "Thanks for your help yesterday"
    ];

    for (const input of inputs) {
        console.log(`\n=== Processing: "${input}" ===`);
        const result = await advanced.processInput(input);
        console.log('Execution path:', result.executionPath);
        console.log('Final outcome:', result.finalState.logged ? 'Logged' : 'Responded');
    }
}

// demoAdvancedConditionals();
```

### Example 4: Conditional Edge Composition and Reusability

```javascript
class ComposableConditionalEdges {
    constructor() {
        this.conditionLibrary = new Map();
        this._buildConditionLibrary();
    }

    _buildConditionLibrary() {
        // Reusable condition functions
        this.conditionLibrary.set('high_urgency', (state) => 
            state.urgency === 'high' || state.priority === 'critical');
        
        this.conditionLibrary.set('needs_human', (state) => 
            state.confidence < 0.6 || state.complexity === 'high');
        
        this.conditionLibrary.set('can_automate', (state) => 
            state.confidence > 0.8 && state.complexity === 'low');
        
        this.conditionLibrary.set('has_entities', (state) => 
            state.entities && state.entities.length > 0);
        
        this.conditionLibrary.set('positive_sentiment', (state) => 
            state.sentiment && state.sentiment.score > 0.3);
        
        this.conditionLibrary.set('negative_sentiment', (state) => 
            state.sentiment && state.sentiment.score < -0.3);
        
        this.conditionLibrary.set('max_attempts_reached', (state) => 
            state.attempts >= (state.max_attempts || 3));
        
        this.conditionLibrary.set('business_hours', (state) => {
            const hour = new Date().getHours();
            return hour >= 9 && hour < 17;
        });
    }

    // Compose multiple conditions
    composeConditions(conditionNames, operator = 'AND') {
        const conditions = conditionNames.map(name => this.conditionLibrary.get(name));
        
        return (state) => {
            if (operator === 'AND') {
                return conditions.every(condition => condition(state));
            } else if (operator === 'OR') {
                return conditions.some(condition => condition(state));
            }
            return false;
        };
    }

    // Create conditional edges with composed conditions
    createCompositeEdge(source, target, conditionNames, operator = 'AND', metadata = {}) {
        const compositeCondition = this.composeConditions(conditionNames, operator);
        const description = `${conditionNames.join(` ${operator} `)}`;
        
        return {
            source,
            target,
            condition: compositeCondition,
            description,
            metadata
        };
    }

    // Example workflow using composite conditions
    buildCompositeWorkflow() {
        const system = new ConditionalEdgeSystem();
        
        // Add simple nodes
        system.addNode('start', async (s) => ({ ...s, started: true }));
        system.addNode('process_auto', async (s) => ({ ...s, automated: true }));
        system.addNode('process_human', async (s) => ({ ...s, human_handled: true }));
        system.addNode('escalate', async (s) => ({ ...s, escalated: true }));
        system.addNode('end', async (s) => ({ ...s, ended: true }));

        // Create composite conditional edges
        const edges = [
            // Automated processing when conditions are ideal
            this.createCompositeEdge(
                'start', 'process_auto',
                ['can_automate', 'positive_sentiment', 'business_hours'],
                'AND',
                { priority: 20, description: 'Ideal conditions for automation' }
            ),
            
            // Human processing for complex or sensitive cases
            this.createCompositeEdge(
                'start', 'process_human',
                ['needs_human', 'high_urgency'],
                'OR',
                { priority: 15, description: 'Requires human attention' }
            ),
            
            // Escalation for negative sentiment during business hours
            this.createCompositeEdge(
                'start', 'escalate',
                ['negative_sentiment', 'business_hours', 'max_attempts_reached'],
                'AND',
                { priority: 25, description: 'Critical issue needing escalation' }
            )
        ];

        // Add edges to system
        edges.forEach(edge => {
            system.addConditionalEdge(edge.source, edge.condition, edge.target, {
                description: edge.description,
                priority: edge.metadata.priority
            });
        });

        system.addDefaultEdge('start', 'process_auto');
        system.addDefaultEdge('process_auto', 'end');
        system.addDefaultEdge('process_human', 'end');
        system.addDefaultEdge('escalate', 'end');

        return system;
    }

    // Test the composite conditions
    testCompositeConditions() {
        console.log('🧪 TESTING COMPOSITE CONDITIONS\n');
        
        const testStates = [
            { urgency: 'high', confidence: 0.9, complexity: 'low', attempts: 1 },
            { urgency: 'low', confidence: 0.4, complexity: 'high', attempts: 5 },
            { sentiment: { score: 0.8 }, confidence: 0.9, complexity: 'low', attempts: 2 }
        ];

        testStates.forEach((state, i) => {
            console.log(`\nTest Case ${i + 1}:`, state);
            
            const conditions = {
                'High urgency OR needs human': this.composeConditions(['high_urgency', 'needs_human'], 'OR'),
                'Can automate AND positive sentiment': this.composeConditions(['can_automate', 'positive_sentiment'], 'AND'),
                'Business hours AND max attempts': this.composeConditions(['business_hours', 'max_attempts_reached'], 'AND')
            };

            Object.entries(conditions).forEach(([desc, condition]) => {
                console.log(`  ${desc}: ${condition(state)}`);
            });
        });
    }
}

// Demonstration
function demoComposableConditions() {
    const composable = new ComposableConditionalEdges();
    composable.testCompositeConditions();
    
    const workflow = composable.buildCompositeWorkflow();
    console.log('\n📊 COMPOSITE WORKFLOW VISUALIZATION');
    workflow.visualize();
}

// demoComposableConditions();
```

## Key Conditional Edge Patterns

### **1. Basic Conditional Patterns**
```javascript
// Simple boolean condition
state => state.value > threshold

// String matching
state => state.category === 'urgent'

// Array checking
state => state.tags.includes('important')
```

### **2. Composite Conditions**
```javascript
// AND condition
state => condition1(state) && condition2(state)

// OR condition  
state => condition1(state) || condition2(state)

// Complex logic
state => (condition1(state) && !condition2(state)) || condition3(state)
```

### **3. State-Based Routing**
```javascript
// Route based on previous node
state => state.previous_node === 'classification'

// Route based on attempt count
state => state.attempts < max_attempts

// Route based on confidence
state => state.confidence > confidence_threshold
```

### **4. Time-Based Conditions**
```javascript
// Business hours check
state => {
    const hour = new Date().getHours();
    return hour >= 9 && hour < 17;
}

// Timeout check
state => Date.now() - state.start_time < timeout_ms
```

## Best Practices for Conditional Edges

### **1. Keep Conditions Simple**
```javascript
// ✅ Good - simple, readable
state => state.urgency === 'high'

// ❌ Avoid - too complex
state => {
    // Complex logic that should be in a node
    return someComplexCalculation(state) && anotherCheck(state);
}
```

### **2. Use Descriptive Condition Names**
```javascript
// ✅ Good - clear intent
const isHighPriority = state => state.priority === 'high' && state.impact === 'critical';

// ❌ Avoid - unclear
const condition1 = state => state.a === 'x' && state.b > 5;
```

### **3. Handle Edge Cases**
```javascript
// ✅ Good - defensive programming
state => {
    if (!state || !state.sentiment) return false;
    return state.sentiment.score < -0.7;
}

// ❌ Avoid - potential errors
state => state.sentiment.score < -0.7 // Could throw if sentiment is undefined
```

### **4. Test Conditions Independently**
```javascript
// Unit test your conditions
const testCondition = (condition, testStates) => {
    testStates.forEach((state, i) => {
        console.log(`Test ${i}:`, condition(state));
    });
};
```

## Common Use Cases

1. **Content Moderation**: Route content based on safety scores
2. **Customer Support**: Escalate based on urgency and complexity
3. **Workflow Routing**: Choose different processing paths
4. **Error Handling**: Retry or escalate based on error types
5. **A/B Testing**: Route users to different experiences
6. **Load Balancing**: Distribute work based on system load

Conditional edges are the "brain" of your LangGraph application, enabling intelligent routing and adaptive behavior based on real-time state analysis.