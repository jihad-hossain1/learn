# LangGraph Essentials: Graph Nodes vs. Edges

LangGraph is a library for building stateful, multi-actor applications with language models, using graphs to define workflows. Understanding nodes and edges is fundamental to designing effective LangGraph applications.

## Core Concepts

### **Graph Nodes**
- **Processing units** where work happens
- **State transformers** that modify the graph state
- **Can be**: LLM calls, tools, functions, conditional logic
- **Key characteristic**: They **do something** with the state

### **Graph Edges**
- **Connections** between nodes
- **Control flow** determiners
- **Can be**: Conditional, unconditional, dynamic
- **Key characteristic**: They **route** the flow based on state

---

## JavaScript Examples

### Example 1: Basic LangGraph Structure with Nodes & Edges

```javascript
class LangGraph {
    constructor() {
        this.nodes = new Map();
        this.edges = new Map();
        this.entryPoint = null;
    }

    // Define a node (processing unit)
    addNode(name, handler) {
        this.nodes.set(name, {
            name,
            handler,
            type: 'function' // Could be 'llm', 'tool', 'conditional'
        });
        return this;
    }

    // Define edges (connections between nodes)
    addEdge(sourceNode, targetNode, condition = null) {
        if (!this.edges.has(sourceNode)) {
            this.edges.set(sourceNode, []);
        }
        
        this.edges.get(sourceNode).push({
            target: targetNode,
            condition: condition || (() => true), // Default: always follow
            type: condition ? 'conditional' : 'unconditional'
        });
        
        return this;
    }

    // Set the starting point
    setEntryPoint(nodeName) {
        this.entryPoint = nodeName;
        return this;
    }

    // Execute the graph
    async execute(initialState) {
        if (!this.entryPoint) {
            throw new Error('No entry point defined');
        }

        let currentState = { ...initialState };
        let currentNode = this.entryPoint;
        const executionPath = [];

        while (currentNode) {
            executionPath.push(currentNode);
            
            // Execute the current node
            const node = this.nodes.get(currentNode);
            if (!node) {
                throw new Error(`Node not found: ${currentNode}`);
            }

            console.log(`🔹 Executing node: ${currentNode}`);
            const result = await node.handler(currentState);
            
            // Update state with node's result
            currentState = { ...currentState, ...result };
            currentState._lastNode = currentNode;
            currentState._executionPath = executionPath;

            // Determine next node based on edges
            const nextNode = this._getNextNode(currentNode, currentState);
            console.log(`➡️ Moving from ${currentNode} to ${nextNode || 'END'}`);
            
            currentNode = nextNode;
        }

        return currentState;
    }

    // Determine next node based on edges and conditions
    _getNextNode(currentNode, state) {
        const edges = this.edges.get(currentNode) || [];
        
        for (const edge of edges) {
            if (edge.condition(state)) {
                return edge.target;
            }
        }
        
        return null; // No valid edges - graph ends
    }

    // Visualize the graph structure
    visualize() {
        console.log('\n📊 GRAPH VISUALIZATION');
        console.log('Nodes:');
        this.nodes.forEach((node, name) => {
            console.log(`  ○ ${name} (${node.type})`);
        });
        
        console.log('\nEdges:');
        this.edges.forEach((edges, source) => {
            edges.forEach(edge => {
                console.log(`  ${source} ──${edge.type}──> ${edge.target}`);
            });
        });
        
        if (this.entryPoint) {
            console.log(`\nEntry Point: ${this.entryPoint}`);
        }
    }
}
```

### Example 2: Chatbot with Different Node Types

```javascript
// Example: Customer Service Chatbot
class CustomerServiceBot {
    constructor() {
        this.graph = new LangGraph();
        this._buildGraph();
    }

    _buildGraph() {
        // Define nodes (different types of processing)
        this.graph
            .addNode('greet', this._greetNode.bind(this))
            .addNode('classify_intent', this._classifyIntentNode.bind(this))
            .addNode('handle_billing', this._handleBillingNode.bind(this))
            .addNode('handle_technical', this._handleTechnicalNode.bind(this))
            .addNode('handle_general', this._handleGeneralNode.bind(this))
            .addNode('escalate', this._escalateNode.bind(this))
            .addNode('end_conversation', this._endConversationNode.bind(this));

        // Define edges (control flow)
        this.graph
            .setEntryPoint('greet')
            
            // From greet: always go to intent classification
            .addEdge('greet', 'classify_intent')
            
            // From classify_intent: conditional edges based on intent
            .addEdge('classify_intent', 'handle_billing', 
                state => state.intent === 'billing')
            .addEdge('classify_intent', 'handle_technical', 
                state => state.intent === 'technical')
            .addEdge('classify_intent', 'handle_general', 
                state => state.intent === 'general')
            .addEdge('classify_intent', 'escalate', 
                state => state.intent === 'escalate')
            
            // From handling nodes: check if resolution achieved
            .addEdge('handle_billing', 'end_conversation', 
                state => state.issue_resolved)
            .addEdge('handle_billing', 'escalate', 
                state => state.attempts > 2 && !state.issue_resolved)
            .addEdge('handle_billing', 'classify_intent', 
                state => !state.issue_resolved && state.attempts <= 2)
            
            .addEdge('handle_technical', 'end_conversation', 
                state => state.issue_resolved)
            .addEdge('handle_technical', 'escalate', 
                state => state.attempts > 3)
            .addEdge('handle_technical', 'classify_intent', 
                state => !state.issue_resolved)
            
            // From escalate: always end conversation
            .addEdge('escalate', 'end_conversation')
            
            // From end_conversation: no edges (graph ends)
            ;

        return this.graph;
    }

    // Node implementations
    async _greetNode(state) {
        console.log('🤖: Hello! Welcome to customer service. How can I help you today?');
        return {
            conversation: [...(state.conversation || []), { role: 'assistant', content: 'Hello! How can I help you?' }],
            attempts: 0
        };
    }

    async _classifyIntentNode(state) {
        const lastMessage = state.conversation?.[state.conversation.length - 1]?.content || '';
        
        // Simple intent classification (in reality, use LLM)
        let intent = 'general';
        if (lastMessage.toLowerCase().includes('bill') || lastMessage.includes('payment')) {
            intent = 'billing';
        } else if (lastMessage.includes('error') || lastMessage.includes('not working')) {
            intent = 'technical';
        } else if (lastMessage.includes('manager') || lastMessage.includes('human')) {
            intent = 'escalate';
        }

        console.log(`🔍 Classified intent: ${intent}`);
        return { intent, attempts: (state.attempts || 0) + 1 };
    }

    async _handleBillingNode(state) {
        console.log('🤖: I see you have a billing question. Let me check your account...');
        
        // Simulate billing resolution
        const resolved = Math.random() > 0.3; // 70% success rate
        const response = resolved 
            ? 'I\'ve found the issue with your bill and applied a credit. Is there anything else?'
            : 'I need more information about your billing issue. Can you provide your account number?';
        
        return {
            issue_resolved: resolved,
            conversation: [...(state.conversation || []), { role: 'assistant', content: response }]
        };
    }

    async _handleTechnicalNode(state) {
        console.log('🤖: I understand this is a technical issue. Let me help troubleshoot...');
        
        // Simulate technical support
        const resolved = Math.random() > 0.5; // 50% success rate
        const response = resolved 
            ? 'The technical issue has been resolved. Please try again now.'
            : 'I suggest restarting the application. If that doesn\'t work, we may need to escalate.';
        
        return {
            issue_resolved: resolved,
            conversation: [...(state.conversation || []), { role: 'assistant', content: response }]
        };
    }

    async _handleGeneralNode(state) {
        console.log('🤖: I\'d be happy to help with your question...');
        return {
            conversation: [...(state.conversation || []), { role: 'assistant', content: 'I understand. Let me look into that for you.' }]
        };
    }

    async _escalateNode(state) {
        console.log('🤖: Let me connect you with a human specialist...');
        return {
            conversation: [...(state.conversation || []), { role: 'assistant', content: 'I\'m escalating this to our specialist team. They will contact you shortly.' }],
            escalated: true
        };
    }

    async _endConversationNode(state) {
        console.log('🤖: Thank you for contacting us. Have a great day!');
        return {
            conversation: [...(state.conversation || []), { role: 'assistant', content: 'Thank you! Have a great day!' }],
            ended: true
        };
    }

    async processMessage(userMessage) {
        const initialState = {
            conversation: [{ role: 'user', content: userMessage }],
            timestamp: new Date().toISOString()
        };

        return await this.graph.execute(initialState);
    }
}

// Usage example
async function demoChatbot() {
    const bot = new CustomerServiceBot();
    bot.graph.visualize();
    
    console.log('\n💬 CHATBOT DEMO');
    const result1 = await bot.processMessage("I have a problem with my bill");
    console.log('Final state:', result1);
    
    console.log('\n--- Another conversation ---');
    const result2 = await bot.processMessage("My app is showing an error");
    console.log('Final state:', result2);
}

// demoChatbot();
```

### Example 3: Advanced Node Types with LLM Integration

```javascript
class AdvancedLangGraph {
    constructor() {
        this.nodes = new Map();
        this.edges = new Map();
        this.nodeTypes = {
            LLM: 'llm_node',
            TOOL: 'tool_node',
            CONDITIONAL: 'conditional_node',
            FUNCTION: 'function_node',
            PARALLEL: 'parallel_node'
        };
    }

    // Different types of nodes with specialized handlers
    addLLMNode(name, promptTemplate, modelConfig = {}) {
        this.nodes.set(name, {
            name,
            type: this.nodeTypes.LLM,
            handler: async (state) => await this._handleLLMNode(state, promptTemplate, modelConfig),
            config: modelConfig
        });
        return this;
    }

    addToolNode(name, toolFunction, inputMapper, outputMapper) {
        this.nodes.set(name, {
            name,
            type: this.nodeTypes.TOOL,
            handler: async (state) => await this._handleToolNode(state, toolFunction, inputMapper, outputMapper),
            tool: toolFunction
        });
        return this;
    }

    addConditionalNode(name, conditionFunction, branches) {
        this.nodes.set(name, {
            name,
            type: this.nodeTypes.CONDITIONAL,
            handler: async (state) => await this._handleConditionalNode(state, conditionFunction),
            branches: branches || []
        });
        return this;
    }

    addParallelNode(name, parallelNodes, aggregator) {
        this.nodes.set(name, {
            name,
            type: this.nodeTypes.PARALLEL,
            handler: async (state) => await this._handleParallelNode(state, parallelNodes, aggregator),
            parallelNodes: parallelNodes || []
        });
        return this;
    }

    // Edge with additional metadata
    addEdge(source, target, condition = null, metadata = {}) {
        if (!this.edges.has(source)) {
            this.edges.set(source, []);
        }

        this.edges.get(source).push({
            source,
            target,
            condition: condition || (() => true),
            metadata: {
                priority: metadata.priority || 1,
                description: metadata.description || '',
                ...metadata
            }
        });
        return this;
    }

    // LLM Node Handler
    async _handleLLMNode(state, promptTemplate, modelConfig) {
        console.log(`🧠 LLM Node: Generating response...`);
        
        // Build prompt from template and state
        const prompt = this._renderPrompt(promptTemplate, state);
        
        // Simulate LLM call (in reality, call OpenAI, Anthropic, etc.)
        const response = await this._simulateLLMCall(prompt, modelConfig);
        
        return {
            llm_response: response,
            last_prompt: prompt,
            conversation: [
                ...(state.conversation || []),
                { role: 'assistant', content: response }
            ]
        };
    }

    // Tool Node Handler
    async _handleToolNode(state, toolFunction, inputMapper, outputMapper) {
        console.log(`🛠️ Tool Node: Executing tool...`);
        
        // Map state to tool inputs
        const toolInputs = inputMapper ? inputMapper(state) : state;
        
        // Execute tool
        const toolResult = await toolFunction(toolInputs);
        
        // Map tool results back to state
        const stateUpdates = outputMapper ? outputMapper(toolResult, state) : { tool_result: toolResult };
        
        return stateUpdates;
    }

    // Conditional Node Handler
    async _handleConditionalNode(state, conditionFunction) {
        console.log(`🎯 Conditional Node: Evaluating conditions...`);
        
        const conditionResult = conditionFunction(state);
        
        return {
            condition_result: conditionResult,
            branching_decisions: [
                ...(state.branching_decisions || []),
                { node: 'conditional', result: conditionResult, timestamp: Date.now() }
            ]
        };
    }

    // Parallel Node Handler
    async _handleParallelNode(state, parallelNodes, aggregator) {
        console.log(`⚡ Parallel Node: Executing ${parallelNodes.length} nodes in parallel...`);
        
        // Execute all parallel nodes
        const parallelResults = await Promise.all(
            parallelNodes.map(async (nodeConfig) => {
                const node = this.nodes.get(nodeConfig.nodeName);
                if (!node) throw new Error(`Node not found: ${nodeConfig.nodeName}`);
                
                const result = await node.handler(state);
                return {
                    node: nodeConfig.nodeName,
                    result: nodeConfig.resultKey ? result[nodeConfig.resultKey] : result,
                    metadata: nodeConfig
                };
            })
        );

        // Aggregate results
        const aggregated = aggregator ? aggregator(parallelResults, state) : {
            parallel_results: parallelResults
        };

        return aggregated;
    }

    // Utility methods
    _renderPrompt(template, state) {
        return template.replace(/\{(\w+)\}/g, (match, key) => state[key] || match);
    }

    async _simulateLLMCall(prompt, config) {
        // Simulate LLM processing time
        await new Promise(resolve => setTimeout(resolve, config.delay || 100));
        
        // Simple simulated response based on prompt
        if (prompt.includes('greet')) return "Hello! How can I assist you today?";
        if (prompt.includes('weather')) return "The weather is sunny and 72°F.";
        if (prompt.includes('help')) return "I'd be happy to help with that!";
        
        return "I understand your request. Let me process that information.";
    }
}

// Example: Research Assistant Graph
class ResearchAssistant {
    constructor() {
        this.graph = new AdvancedLangGraph();
        this._buildResearchGraph();
    }

    _buildResearchGraph() {
        // Define various node types
        this.graph
            // LLM Nodes
            .addLLMNode('analyze_query', 
                "Analyze this research query: {user_query}. Identify key topics and search terms.",
                { model: 'gpt-4', temperature: 0.1 })
            
            .addLLMNode('synthesize_results',
                "Synthesize these research findings: {search_results}. Create a comprehensive summary.",
                { model: 'gpt-4', temperature: 0.3 })
            
            // Tool Nodes
            .addToolNode('web_search', 
                this._simulateWebSearch.bind(this),
                state => ({ query: state.search_terms }),
                (result, state) => ({ search_results: result }))
            
            .addToolNode('format_output',
                this._formatResearchReport.bind(this),
                state => ({ summary: state.synthesis, sources: state.search_results }),
                (result) => ({ final_report: result }))
            
            // Conditional Node
            .addConditionalNode('needs_deeper_research',
                state => state.search_results?.length < 3 || state.confidence_score < 0.7)
            
            // Parallel Node for multi-source validation
            .addParallelNode('validate_sources', 
                [
                    { nodeName: 'check_fact1', resultKey: 'fact_check' },
                    { nodeName: 'check_fact2', resultKey: 'secondary_check' }
                ],
                (results, state) => ({
                    validation_results: results,
                    all_valid: results.every(r => r.result.valid),
                    confidence_score: results.filter(r => r.result.valid).length / results.length
                }));

        // Define edges with conditions
        this.graph
            .addEdge('analyze_query', 'web_search')
            .addEdge('web_search', 'needs_deeper_research')
            .addEdge('needs_deeper_research', 'web_search', 
                state => state.condition_result === true,
                { description: 'Loop back for more research if needed' })
            .addEdge('needs_deeper_research', 'validate_sources', 
                state => state.condition_result === false)
            .addEdge('validate_sources', 'synthesize_results')
            .addEdge('synthesize_results', 'format_output');

        return this.graph;
    }

    async _simulateWebSearch(query) {
        console.log(`🔍 Searching web for: ${query.query}`);
        await new Promise(resolve => setTimeout(resolve, 200));
        
        // Simulated search results
        return [
            { title: "Source 1 about " + query.query, content: "Information from source 1", reliability: 0.9 },
            { title: "Source 2 about " + query.query, content: "Information from source 2", reliability: 0.8 },
            { title: "Source 3 about " + query.query, content: "Information from source 3", reliability: 0.7 }
        ];
    }

    async _formatResearchReport(data) {
        console.log('📝 Formatting research report...');
        return {
            title: "Research Report",
            summary: data.summary,
            sources: data.sources,
            generated_at: new Date().toISOString()
        };
    }

    async research(query) {
        const initialState = {
            user_query: query,
            timestamp: new Date().toISOString()
        };

        // In real implementation, this would execute the graph
        console.log('Starting research graph execution...');
        // const result = await this.graph.execute(initialState);
        // return result;
        
        return { status: 'Research completed', query };
    }
}
```

### Example 4: Edge Types and Routing Logic

```javascript
class EdgeTypesDemo {
    constructor() {
        this.graph = new LangGraph();
        this.edgeTypes = {
            UNCONDITIONAL: 'always_follow',
            CONDITIONAL: 'if_condition',
            PROBABILISTIC: 'probability_based',
            DYNAMIC: 'runtime_determined'
        };
    }

    buildEdgeDemoGraph() {
        // Add nodes
        this.graph
            .addNode('start', async (state) => {
                console.log('Start node');
                return { ...state, start: true };
            })
            .addNode('process_a', async (state) => {
                console.log('Processing path A');
                return { ...state, path: 'A', processed: true };
            })
            .addNode('process_b', async (state) => {
                console.log('Processing path B');
                return { ...state, path: 'B', processed: true };
            })
            .addNode('process_c', async (state) => {
                console.log('Processing path C');
                return { ...state, path: 'C', processed: true };
            })
            .addNode('end', async (state) => {
                console.log('End node');
                return { ...state, ended: true };
            });

        // Different edge types
        this.graph
            .setEntryPoint('start')
            
            // Unconditional edge - always follow
            .addEdge('start', 'process_a')
            
            // Conditional edge - based on state
            .addEdge('process_a', 'process_b', 
                state => state.user_type === 'premium')
            .addEdge('process_a', 'process_c', 
                state => state.user_type !== 'premium')
            
            // Probabilistic edge - random choice
            .addEdge('process_b', 'end', this._probabilisticEdge(0.7))
            .addEdge('process_b', 'process_c', this._probabilisticEdge(0.3))
            
            // Dynamic edge - determined at runtime
            .addEdge('process_c', 'end', this._dynamicEdge());

        return this.graph;
    }

    // Probabilistic edge condition
    _probabilisticEdge(probability) {
        return (state) => {
            const randomValue = Math.random();
            const shouldFollow = randomValue < probability;
            console.log(`🎲 Probabilistic edge: ${probability} chance -> ${shouldFollow ? 'FOLLOW' : 'SKIP'}`);
            return shouldFollow;
        };
    }

    // Dynamic edge condition
    _dynamicEdge() {
        return (state) => {
            // Complex logic that can change based on multiple factors
            const currentHour = new Date().getHours();
            const isBusinessHours = currentHour >= 9 && currentHour < 17;
            const hasHighPriority = state.priority === 'high';
            const systemLoad = Math.random(); // Simulated system load
            
            const shouldFollow = isBusinessHours && (hasHighPriority || systemLoad < 0.8);
            console.log(`🔄 Dynamic edge evaluation: business_hours=${isBusinessHours}, high_priority=${hasHighPriority}, load=${systemLoad.toFixed(2)} -> ${shouldFollow}`);
            
            return shouldFollow;
        };
    }

    async demonstrateEdges() {
        console.log('=== EDGE TYPE DEMONSTRATION ===\n');
        
        // Test case 1: Premium user
        console.log('🧪 Test 1: Premium user');
        await this.graph.execute({ user_type: 'premium' });
        
        console.log('\n🧪 Test 2: Regular user');
        await this.graph.execute({ user_type: 'regular' });
        
        console.log('\n🧪 Test 3: High priority during business hours');
        await this.graph.execute({ 
            user_type: 'premium', 
            priority: 'high' 
        });
    }
}
```

## Key Differences: Nodes vs Edges

| Aspect | Nodes | Edges |
|--------|-------|-------|
| **Purpose** | **Process** data | **Route** flow |
| **Function** | Transform state | Evaluate conditions |
| **Complexity** | Can be complex (LLMs, tools) | Usually simple logic |
| **State Impact** | **Modify** state | **Use** state for decisions |
| **Types** | LLM, Tool, Function, Conditional | Unconditional, Conditional, Dynamic |

## Node Types in LangGraph

1. **LLM Nodes**: Call language models, generate text
2. **Tool Nodes**: Execute functions, call APIs, process data
3. **Conditional Nodes**: Make branching decisions
4. **Function Nodes**: Custom business logic
5. **Parallel Nodes**: Execute multiple nodes concurrently

## Edge Types in LangGraph

1. **Unconditional Edges**: Always follow this path
2. **Conditional Edges**: Follow if condition is met
3. **Probabilistic Edges**: Follow based on probability
4. **Dynamic Edges**: Runtime-determined routing

## Best Practices

### **For Nodes:**
- Keep nodes focused on single responsibilities
- Design nodes to be reusable across graphs
- Handle errors gracefully within nodes
- Document node inputs/outputs clearly

### **For Edges:**
- Keep edge conditions simple and testable
- Use descriptive condition names
- Avoid complex business logic in edges
- Consider fallback edges for error handling

### **Graph Design:**
- Start with simple linear flows
- Add complexity gradually with conditionals
- Test each node and edge independently
- Monitor graph execution paths

Understanding the distinction between nodes (the "workers") and edges (the "directors") is crucial for designing effective LangGraph applications that are maintainable, debuggable, and scalable.