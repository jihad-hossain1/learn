# Multi-Agent Orchestration with Message Passing

Multi-Agent Orchestration refers to coordinating multiple AI agents that work together through message passing to accomplish complex tasks. Let me explain the concepts and provide JavaScript examples.

## Core Concepts

### 1. **Message Passing Protocols**
Agents communicate through structured messages following specific protocols:

```javascript
// Basic message structure
class AgentMessage {
  constructor(sender, receiver, content, type, timestamp = Date.now()) {
    this.sender = sender;
    this.receiver = receiver;
    this.content = content;
    this.type = type; // 'request', 'response', 'notification'
    this.timestamp = timestamp;
    this.id = this.generateId();
  }
  
  generateId() {
    return `${this.sender}-${this.timestamp}-${Math.random().toString(36).substr(2, 9)}`;
  }
}
```

### 2. **Agent Base Class**

```javascript
class Agent {
  constructor(name, capabilities = []) {
    this.name = name;
    this.capabilities = capabilities;
    this.inbox = [];
    this.outbox = [];
    this.peers = new Map(); // Other agents this agent can communicate with
  }
  
  // Register with an orchestration system
  register(orchestrator) {
    this.orchestrator = orchestrator;
    orchestrator.registerAgent(this);
  }
  
  // Send message to another agent
  sendMessage(receiver, content, type = 'request') {
    const message = new AgentMessage(this.name, receiver, content, type);
    this.outbox.push(message);
    
    if (this.orchestrator) {
      this.orchestrator.routeMessage(message);
    }
    
    return message.id;
  }
  
  // Receive message
  receiveMessage(message) {
    this.inbox.push(message);
    this.processMessage(message);
  }
  
  // Process incoming message (to be implemented by specific agents)
  processMessage(message) {
    console.log(`${this.name} received message:`, message);
  }
  
  // Get unread messages
  getUnreadMessages() {
    return this.inbox.filter(msg => !msg.read);
  }
}
```

## Orchestration Patterns

### 1. **Centralized Orchestrator**

```javascript
class CentralOrchestrator {
  constructor() {
    this.agents = new Map();
    this.messageQueue = [];
    this.workflows = new Map();
  }
  
  registerAgent(agent) {
    this.agents.set(agent.name, agent);
    console.log(`Agent ${agent.name} registered`);
  }
  
  // Route messages between agents
  routeMessage(message) {
    const receiver = this.agents.get(message.receiver);
    if (receiver) {
      receiver.receiveMessage(message);
    } else {
      console.error(`Agent ${message.receiver} not found`);
    }
  }
  
  // Create workflow
  createWorkflow(name, steps) {
    this.workflows.set(name, steps);
  }
  
  // Execute workflow
  async executeWorkflow(workflowName, initialData) {
    const workflow = this.workflows.get(workflowName);
    if (!workflow) throw new Error(`Workflow ${workflowName} not found`);
    
    let currentData = initialData;
    
    for (const step of workflow) {
      const { agent, task, waitForResponse = true } = step;
      
      console.log(`Executing step: ${agent} -> ${task}`);
      
      const messageId = this.agents.get(agent).sendMessage(
        agent, // self-message for task execution
        { task, data: currentData },
        'task'
      );
      
      if (waitForResponse) {
        // Wait for response (simplified)
        await this.waitForResponse(messageId, agent);
      }
    }
  }
  
  async waitForResponse(messageId, agentName) {
    // Simplified implementation - in real scenario, use proper async handling
    return new Promise(resolve => setTimeout(resolve, 1000));
  }
}
```

### 2. **Specialized Agent Examples**

```javascript
// Research Agent
class ResearchAgent extends Agent {
  constructor(name) {
    super(name, ['web_search', 'data_analysis', 'report_generation']);
  }
  
  async processMessage(message) {
    if (message.type === 'task') {
      const { task, data } = message.content;
      
      switch (task) {
        case 'research_topic':
          const research = await this.researchTopic(data.topic);
          this.sendMessage(message.sender, { research }, 'response');
          break;
          
        case 'analyze_data':
          const analysis = await this.analyzeData(data.dataset);
          this.sendMessage(message.sender, { analysis }, 'response');
          break;
      }
    }
  }
  
  async researchTopic(topic) {
    // Simulate research
    return `Research results about ${topic}`;
  }
  
  async analyzeData(dataset) {
    // Simulate analysis
    return `Analysis of ${dataset}`;
  }
}

// Writing Agent
class WritingAgent extends Agent {
  constructor(name) {
    super(name, ['content_creation', 'editing', 'formatting']);
  }
  
  async processMessage(message) {
    if (message.type === 'task') {
      const { task, data } = message.content;
      
      switch (task) {
        case 'write_article':
          const article = await this.writeArticle(data.outline);
          this.sendMessage(message.sender, { article }, 'response');
          break;
          
        case 'format_content':
          const formatted = await this.formatContent(data.content);
          this.sendMessage(message.sender, { formatted }, 'response');
          break;
      }
    }
  }
  
  async writeArticle(outline) {
    return `Article based on: ${outline}`;
  }
  
  async formatContent(content) {
    return `Formatted: ${content}`;
  }
}
```

## Complete Example: Content Creation Workflow

```javascript
class ContentCreationOrchestration {
  constructor() {
    this.orchestrator = new CentralOrchestrator();
    this.setupAgents();
    this.setupWorkflows();
  }
  
  setupAgents() {
    // Create specialized agents
    this.researchAgent = new ResearchAgent('researcher');
    this.writingAgent = new WritingAgent('writer');
    this.qualityAgent = new QualityAgent('quality_checker');
    
    // Register agents
    this.researchAgent.register(this.orchestrator);
    this.writingAgent.register(this.orchestrator);
    this.qualityAgent.register(this.orchestrator);
  }
  
  setupWorkflows() {
    // Define content creation workflow
    this.orchestrator.createWorkflow('content_creation', [
      { agent: 'researcher', task: 'research_topic' },
      { agent: 'writer', task: 'write_article' },
      { agent: 'quality_checker', task: 'review_content' }
    ]);
  }
  
  async createContent(topic) {
    console.log(`Starting content creation for: ${topic}`);
    
    await this.orchestrator.executeWorkflow('content_creation', {
      topic: topic,
      requirements: 'comprehensive, engaging, accurate'
    });
    
    console.log('Content creation completed!');
  }
}

// Quality Check Agent
class QualityAgent extends Agent {
  constructor(name) {
    super(name, ['quality_check', 'feedback', 'validation']);
  }
  
  async processMessage(message) {
    if (message.type === 'task') {
      const { task, data } = message.content;
      
      if (task === 'review_content') {
        const review = await this.reviewContent(data.content);
        this.sendMessage(message.sender, { review, approved: true }, 'response');
      }
    }
  }
  
  async reviewContent(content) {
    return `Review completed for content: ${content.substring(0, 50)}...`;
  }
}

// Usage example
async function demo() {
  const contentSystem = new ContentCreationOrchestration();
  
  // Start content creation workflow
  await contentSystem.createContent('Artificial Intelligence in Healthcare');
  
  // Simulate direct agent communication
  contentSystem.researchAgent.sendMessage(
    'writer',
    { research: 'Latest AI healthcare trends' },
    'notification'
  );
}

demo();
```

## Advanced Message Passing Patterns

### 1. **Pub/Sub Pattern**

```javascript
class PubSubOrchestrator extends CentralOrchestrator {
  constructor() {
    super();
    this.subscriptions = new Map();
  }
  
  subscribe(agent, topic) {
    if (!this.subscriptions.has(topic)) {
      this.subscriptions.set(topic, new Set());
    }
    this.subscriptions.get(topic).add(agent);
  }
  
  publish(topic, message) {
    const subscribers = this.subscriptions.get(topic) || new Set();
    subscribers.forEach(agentName => {
      const agent = this.agents.get(agentName);
      if (agent) {
        agent.receiveMessage({
          ...message,
          topic: topic,
          type: 'pubsub'
        });
      }
    });
  }
}
```

### 2. **Request-Response with Timeout**

```javascript
class ResponsiveAgent extends Agent {
  constructor(name) {
    super(name);
    this.pendingRequests = new Map();
  }
  
  sendRequest(receiver, content, timeout = 5000) {
    const messageId = this.sendMessage(receiver, content, 'request');
    
    return new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        this.pendingRequests.delete(messageId);
        reject(new Error(`Request timeout after ${timeout}ms`));
      }, timeout);
      
      this.pendingRequests.set(messageId, { resolve, reject, timeoutId });
    });
  }
  
  processMessage(message) {
    if (message.type === 'response' && this.pendingRequests.has(message.id)) {
      const { resolve, timeoutId } = this.pendingRequests.get(message.id);
      clearTimeout(timeoutId);
      this.pendingRequests.delete(message.id);
      resolve(message.content);
    }
  }
}
```

## Key Benefits of Multi-Agent Orchestration

1. **Modularity**: Each agent has specific capabilities
2. **Scalability**: Easy to add new agents
3. **Fault Tolerance**: Failure in one agent doesn't break the system
4. **Specialization**: Agents can be optimized for specific tasks
5. **Flexibility**: Dynamic routing and workflow management

This architecture is particularly useful for complex AI systems where different models or services need to collaborate through structured communication protocols.