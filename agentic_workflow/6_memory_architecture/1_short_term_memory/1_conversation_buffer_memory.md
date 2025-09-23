I'll explain Short-Term Memory (STM) architecture with a focus on conversation buffer memory, including JavaScript examples.

## Short-Term Memory (STM) Architecture

STM refers to temporary storage that holds information for brief periods (seconds to minutes). In AI systems, this translates to conversation history management.

## Conversation Buffer Memory

This maintains a rolling window of recent interactions, storing the most recent messages while discarding older ones to stay within context limits.

### Key Characteristics:

- **Limited capacity** (token/character limits)
- **First-In-First-Out (FIFO)** behavior
- **Temporary retention** of recent exchanges

## JavaScript Implementation Examples

### 1. Basic Conversation Buffer Class

```javascript
class ConversationBufferMemory {
  constructor(maxSize = 10) {
    this.maxSize = maxSize; // Maximum number of messages to retain
    this.messages = [];
  }

  addMessage(role, content) {
    const message = { role, content, timestamp: Date.now() };
    this.messages.push(message);

    // Remove oldest messages if exceeding max size
    while (this.messages.length > this.maxSize) {
      this.messages.shift();
    }
  }

  getConversationHistory() {
    return this.messages.slice(); // Return copy of array
  }

  clear() {
    this.messages = [];
  }

  getCurrentSize() {
    return this.messages.length;
  }
}

// Usage example
const chatMemory = new ConversationBufferMemory(5);

// Add messages to the conversation
chatMemory.addMessage('user', 'Hello, how are you?');
chatMemory.addMessage('assistant', 'I\\'m doing well! How can I help you?');
chatMemory.addMessage('user', 'What\\'s the weather like?');
chatMemory.addMessage('assistant', 'I don\\'t have real-time weather data.');
chatMemory.addMessage('user', 'Then what can you do?');

console.log('Current conversation:', chatMemory.getConversationHistory());
```

### 2. Advanced Buffer with Token Counting

```javascript
class TokenAwareBufferMemory {
  constructor(maxTokens = 1000) {
    this.maxTokens = maxTokens;
    this.messages = [];
    this.currentTokenCount = 0;
  }

  // Simple token counter (approximate)
  countTokens(text) {
    return Math.ceil(text.split(/\s+/).length * 1.3); // Rough estimation
  }

  addMessage(role, content) {
    const message = {
      role,
      content,
      tokens: this.countTokens(content),
      timestamp: Date.now(),
    };

    // Check if adding this message would exceed token limit
    while (
      this.currentTokenCount + message.tokens > this.maxTokens &&
      this.messages.length > 0
    ) {
      const removedMessage = this.messages.shift();
      this.currentTokenCount -= removedMessage.tokens;
    }

    this.messages.push(message);
    this.currentTokenCount += message.tokens;
  }

  getFormattedConversation() {
    return this.messages.map((msg) => `${msg.role}: ${msg.content}`).join("\n");
  }

  getStats() {
    return {
      messageCount: this.messages.length,
      tokenCount: this.currentTokenCount,
      maxTokens: this.maxTokens,
    };
  }
}

// Usage
const tokenMemory = new TokenAwareBufferMemory(500);

tokenMemory.addMessage("user", "Tell me about JavaScript closures");
tokenMemory.addMessage(
  "assistant",
  "Closures are functions that remember their lexical scope..."
);
tokenMemory.addMessage("user", "Can you give me an example?");

console.log(tokenMemory.getFormattedConversation());
console.log("Memory stats:", tokenMemory.getStats());
```

### 3. Conversation Buffer with Context Preservation

```javascript
class SmartConversationBuffer {
  constructor(maxMessages = 6, preserveSystemMessages = true) {
    this.maxMessages = maxMessages;
    this.preserveSystemMessages = preserveSystemMessages;
    this.conversation = [];
  }

  addMessage(role, content, isSystemContext = false) {
    const message = {
      role,
      content,
      isSystemContext,
      priority: isSystemContext ? 1 : 0, // System messages have higher priority
      timestamp: Date.now(),
    };

    this.conversation.push(message);
    this.manageBuffer();
  }

  manageBuffer() {
    if (this.conversation.length <= this.maxMessages) return;

    // Preserve system messages if configured to do so
    if (this.preserveSystemMessages) {
      const systemMessages = this.conversation.filter(
        (msg) => msg.isSystemContext
      );
      const regularMessages = this.conversation.filter(
        (msg) => !msg.isSystemContext
      );

      // Remove oldest regular messages first
      while (
        regularMessages.length + systemMessages.length > this.maxMessages &&
        regularMessages.length > 0
      ) {
        regularMessages.shift();
      }

      this.conversation = [...systemMessages, ...regularMessages];
    } else {
      // Simple FIFO removal
      while (this.conversation.length > this.maxMessages) {
        this.conversation.shift();
      }
    }
  }

  getConversationContext() {
    return this.conversation.map((msg) => ({
      role: msg.role,
      content: msg.content,
    }));
  }

  // Method to get summary when buffer needs to be compressed
  createSummary() {
    const recentMessages = this.conversation.slice(-3); // Last 3 messages
    return `Recent conversation context: ${recentMessages
      .map((m) => m.content.substring(0, 50))
      .join("... ")}`;
  }
}

// Usage example
const smartBuffer = new SmartConversationBuffer(5, true);

// Add system message (will be preserved)
smartBuffer.addMessage(
  "system",
  "You are a helpful assistant specialized in programming.",
  true
);

// Add conversation
smartBuffer.addMessage("user", "How do I create a class in JavaScript?");
smartBuffer.addMessage(
  "assistant",
  "You can create a class using the class syntax..."
);
smartBuffer.addMessage("user", "What about inheritance?");
smartBuffer.addMessage(
  "assistant",
  "Use the extends keyword for inheritance..."
);

console.log("Conversation context:", smartBuffer.getConversationContext());
```

### 4. Practical Chat Application Example

```javascript
// Simulating a chat application with memory management
class ChatSession {
  constructor() {
    this.memory = new ConversationBufferMemory(10);
    this.sessionStart = Date.now();
  }

  processUserInput(input) {
    // Add user message to memory
    this.memory.addMessage('user', input);

    // Generate response (simulated)
    const response = this.generateResponse(input);

    // Add assistant response to memory
    this.memory.addMessage('assistant', response);

    return response;
  }

  generateResponse(input) {
    const responses = {
      hello: 'Hello! How can I assist you today?',
      help: 'I can help with conversation memory management examples.',
      memory: `Current memory has ${this.memory.getCurrentSize()} messages stored.`,
      default: 'I understand. What would you like to know next?'
    };

    const lowerInput = input.toLowerCase();
    if (lowerInput.includes('hello')) return responses.hello;
    if (lowerInput.includes('help')) return responses.help;
    if (lowerInput.includes('memory')) return responses.memory;

    return responses.default;
  }

  getChatHistory() {
    return this.memory.getConversationHistory();
  }
}

// Demo usage
const chatSession = new ChatSession();

// Simulate conversation
console.log(chatSession.processUserInput('Hello'));
console.log(chatSession.processUserInput('Can you tell me about memory?'));
console.log(chatSession.processUserInput('That\\'s helpful, thank you!'));

console.log('\\nFull conversation history:');
console.log(chatSession.getChatHistory());
```

## Key Benefits of Conversation Buffer Memory:

1. **Context Maintenance**: Keeps recent conversation relevant
2. **Resource Management**: Prevents memory overflow
3. **Performance**: Limits processing to recent interactions
4. **Continuity**: Maintains conversation flow

## Typical Use Cases:

- Chatbots and virtual assistants
- Interactive dialogue systems
- Multi-turn conversation applications
- Context-aware AI applications

This architecture ensures that your AI system maintains relevant context while efficiently managing computational resources.
