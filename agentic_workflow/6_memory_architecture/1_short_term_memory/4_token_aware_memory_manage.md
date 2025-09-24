# Token-Aware Memory Management

Token-aware memory management focuses on optimizing memory usage based on token counts rather than simple message counts. This is crucial for AI systems with strict token limits (like GPT models with 4K-128K context windows).

## How Token-Aware Memory Works

- **Tracks token counts** for each message/conversation piece
- **Prioritizes content** based on importance and token efficiency
- **Implements smart truncation** strategies
- **Optimizes context windows** for AI model consumption

## JavaScript Implementation Examples

### 1. Basic Token-Aware Memory Class

```javascript
class TokenAwareMemory {
  constructor(maxTokens = 4000, tokenBuffer = 100) {
    this.maxTokens = maxTokens;
    this.tokenBuffer = tokenBuffer; // Reserve tokens for response
    this.messages = [];
    this.totalTokens = 0;
  }

  // Approximate token counting (simple version)
  countTokens(text) {
    if (!text) return 0;
    
    // Basic estimation: 1 token ≈ 4 characters or 0.75 words
    const wordCount = text.split(/\s+/).length;
    const charCount = text.length;
    
    // Use average of word-based and char-based estimation
    return Math.ceil((wordCount * 1.3 + charCount / 4) / 2);
  }

  addMessage(role, content, metadata = {}) {
    const tokens = this.countTokens(content);
    const message = {
      role,
      content,
      tokens,
      timestamp: Date.now(),
      importance: metadata.importance || 1.0,
      ...metadata
    };

    // Check if we need to make space
    this.ensureSpace(tokens);

    this.messages.push(message);
    this.totalTokens += tokens;
  }

  ensureSpace(requiredTokens) {
    const availableTokens = this.maxTokens - this.tokenBuffer;
    
    while (this.totalTokens + requiredTokens > availableTokens && this.messages.length > 0) {
      // Remove least important message first
      const leastImportantIndex = this.findLeastImportantMessage();
      const removed = this.messages.splice(leastImportantIndex, 1)[0];
      this.totalTokens -= removed.tokens;
    }
  }

  findLeastImportantMessage() {
    let minImportance = Infinity;
    let minIndex = 0;

    this.messages.forEach((msg, index) => {
      // Skip system messages if possible
      if (msg.role === 'system') return;

      // Calculate score: lower importance + older = more likely to remove
      const ageFactor = (Date.now() - msg.timestamp) / (60 * 1000); // minutes
      const score = msg.importance - (ageFactor * 0.01); // Older messages get slightly penalized

      if (score < minImportance) {
        minImportance = score;
        minIndex = index;
      }
    });

    return minIndex;
  }

  getContext() {
    return this.messages.map(msg => ({
      role: msg.role,
      content: msg.content
    }));
  }

  getTokenUsage() {
    return {
      totalTokens: this.totalTokens,
      availableTokens: this.maxTokens - this.tokenBuffer - this.totalTokens,
      messageCount: this.messages.length,
      tokenLimit: this.maxTokens
    };
  }

  compressMessage(message, targetTokens) {
    if (message.tokens <= targetTokens) return message;

    // Simple compression by truncation
    const words = message.content.split(' ');
    let compressedContent = '';
    let compressedTokens = 0;

    for (const word of words) {
      const newTokens = this.countTokens(compressedContent + ' ' + word);
      if (newTokens > targetTokens) break;
      
      compressedContent += (compressedContent ? ' ' : '') + word;
      compressedTokens = newTokens;
    }

    return {
      ...message,
      content: compressedContent + '...',
      tokens: compressedTokens
    };
  }

  clear() {
    this.messages = [];
    this.totalTokens = 0;
  }
}

// Usage example
const tokenMemory = new TokenAwareMemory(1000, 50);

tokenMemory.addMessage('system', 'You are a helpful assistant.', { importance: 2.0 });
tokenMemory.addMessage('user', 'Can you explain machine learning?', { importance: 1.5 });
tokenMemory.addMessage('assistant', 'Machine learning is a subset of AI that enables computers to learn without explicit programming...', { importance: 1.0 });

console.log('Token Usage:', tokenMemory.getTokenUsage());
console.log('Context:', tokenMemory.getContext());
```

### 2. Advanced Token Manager with Multiple Strategies

```javascript
class AdvancedTokenManager {
  constructor(maxTokens = 4000, strategies = {}) {
    this.maxTokens = maxTokens;
    this.strategies = {
      compression: strategies.compression || 'truncate',
      prioritization: strategies.prioritization || 'importance',
      removal: strategies.removal || 'least_important'
    };
    
    this.conversation = [];
    this.tokenUsage = {
      total: 0,
      byRole: { system: 0, user: 0, assistant: 0 },
      byImportance: { high: 0, medium: 0, low: 0 }
    };
  }

  // More accurate token estimation
  accurateTokenCount(text) {
    if (typeof text !== 'string') return 0;

    // Improved token estimation algorithm
    const cleanText = text.replace(/[^\w\s]/g, ' ').replace(/\s+/g, ' ');
    const words = cleanText.trim().split(/\s+/);
    
    let tokenCount = 0;
    
    words.forEach(word => {
      if (word.length <= 3) {
        tokenCount += 1; // Short words ≈ 1 token
      } else if (word.length <= 6) {
        tokenCount += 1.5; // Medium words
      } else {
        tokenCount += 2; // Long words
      }
    });

    // Add tokens for punctuation and special characters
    const specialChars = text.match(/[^\w\s]/g) || [];
    tokenCount += specialChars.length * 0.3;

    return Math.ceil(tokenCount);
  }

  addMessage(role, content, options = {}) {
    const tokens = this.accurateTokenCount(content);
    const message = {
      role,
      content,
      tokens,
      timestamp: Date.now(),
      importance: options.importance || this.calculateImportance(role, content),
      id: this.generateId(),
      metadata: options.metadata || {}
    };

    this.applyMemoryStrategy(message);
    this.updateTokenUsage(message);
  }

  calculateImportance(role, content) {
    let baseImportance = 1.0;
    
    // Role-based importance
    if (role === 'system') baseImportance = 3.0;
    else if (role === 'user') baseImportance = 2.0;
    
    // Content-based importance factors
    const importanceFactors = {
      question: content.includes('?') ? 0.5 : 0,
      length: Math.min(content.length / 100, 1.0), // Longer messages slightly more important
      urgency: content.match(/urgent|important|asap/i) ? 0.5 : 0
    };

    return baseImportance + Object.values(importanceFactors).reduce((a, b) => a + b, 0);
  }

  applyMemoryStrategy(newMessage) {
    const availableTokens = this.maxTokens - newMessage.tokens;

    if (availableTokens >= 0) {
      this.conversation.push(newMessage);
      return;
    }

    // Need to make space - apply chosen strategy
    switch (this.strategies.removal) {
      case 'least_important':
        this.removeLeastImportant(newMessage);
        break;
      case 'fifo':
        this.removeFIFO(newMessage);
        break;
      case 'selective_compression':
        this.compressSelective(newMessage);
        break;
      default:
        this.removeLeastImportant(newMessage);
    }
  }

  removeLeastImportant(newMessage) {
    // Sort by importance (ascending) and remove until we have space
    const sorted = [...this.conversation].sort((a, b) => a.importance - b.importance);
    
    let removedTokens = 0;
    const toRemove = new Set();

    for (const msg of sorted) {
      if (removedTokens >= newMessage.tokens) break;
      
      toRemove.add(msg.id);
      removedTokens += msg.tokens;
    }

    this.conversation = this.conversation.filter(msg => !toRemove.has(msg.id));
    this.conversation.push(newMessage);
    this.recalculateTokenUsage();
  }

  removeFIFO(newMessage) {
    // Remove oldest messages first
    const sorted = [...this.conversation].sort((a, b) => a.timestamp - b.timestamp);
    
    let removedTokens = 0;
    const toRemove = new Set();

    for (const msg of sorted) {
      if (removedTokens >= newMessage.tokens) break;
      
      toRemove.add(msg.id);
      removedTokens += msg.tokens;
    }

    this.conversation = this.conversation.filter(msg => !toRemove.has(msg.id));
    this.conversation.push(newMessage);
    this.recalculateTokenUsage();
  }

  compressSelective(newMessage) {
    // Try compressing low-importance messages first
    const compressible = this.conversation.filter(msg => 
      msg.importance < 2.0 && msg.tokens > 20
    ).sort((a, b) => a.importance - b.importance);

    let spaceFreed = 0;

    for (const msg of compressible) {
      if (spaceFreed >= newMessage.tokens) break;

      const originalTokens = msg.tokens;
      const compressed = this.compressMessage(msg, Math.ceil(msg.tokens * 0.6)); // Compress to 60%
      const tokensFreed = originalTokens - compressed.tokens;

      if (tokensFreed > 0) {
        const index = this.conversation.findIndex(m => m.id === msg.id);
        this.conversation[index] = compressed;
        spaceFreed += tokensFreed;
      }
    }

    // If compression didn't free enough space, remove messages
    if (spaceFreed < newMessage.tokens) {
      this.removeLeastImportant(newMessage);
    } else {
      this.conversation.push(newMessage);
      this.recalculateTokenUsage();
    }
  }

  compressMessage(message, targetTokens) {
    if (message.tokens <= targetTokens) return message;

    const sentences = message.content.split(/[.!?]+/).filter(s => s.trim());
    let compressedContent = '';
    let currentTokens = 0;

    for (const sentence of sentences) {
      const sentenceTokens = this.accurateTokenCount(sentence);
      
      if (currentTokens + sentenceTokens <= targetTokens) {
        compressedContent += (compressedContent ? '. ' : '') + sentence.trim();
        currentTokens += sentenceTokens;
      } else {
        break;
      }
    }

    if (compressedContent && currentTokens > 0) {
      return {
        ...message,
        content: compressedContent + (compressedContent !== message.content ? '...' : ''),
        tokens: currentTokens
      };
    }

    // Fallback to word-level truncation
    return this.truncateMessage(message, targetTokens);
  }

  truncateMessage(message, targetTokens) {
    const words = message.content.split(' ');
    let truncated = '';
    let currentTokens = 0;

    for (const word of words) {
      const wordTokens = this.accurateTokenCount(word);
      
      if (currentTokens + wordTokens <= targetTokens) {
        truncated += (truncated ? ' ' : '') + word;
        currentTokens += wordTokens;
      } else {
        break;
      }
    }

    return {
      ...message,
      content: truncated + '...',
      tokens: currentTokens
    };
  }

  updateTokenUsage(message) {
    this.tokenUsage.total += message.tokens;
    this.tokenUsage.byRole[message.role] = (this.tokenUsage.byRole[message.role] || 0) + message.tokens;
    
    const importanceLevel = message.importance >= 2.5 ? 'high' : message.importance >= 1.5 ? 'medium' : 'low';
    this.tokenUsage.byImportance[importanceLevel] += message.tokens;
  }

  recalculateTokenUsage() {
    this.tokenUsage = {
      total: 0,
      byRole: { system: 0, user: 0, assistant: 0 },
      byImportance: { high: 0, medium: 0, low: 0 }
    };

    this.conversation.forEach(msg => this.updateTokenUsage(msg));
  }

  generateId() {
    return Math.random().toString(36).substr(2, 9);
  }

  getOptimizedContext(targetTokens = null) {
    const availableTokens = targetTokens || (this.maxTokens - 200); // Reserve for response
    
    if (this.tokenUsage.total <= availableTokens) {
      return this.conversation.map(msg => ({ role: msg.role, content: msg.content }));
    }

    // Need to optimize for target token count
    return this.optimizeForTokenLimit(availableTokens);
  }

  optimizeForTokenLimit(availableTokens) {
    // Sort by importance (descending) and take as many as fit
    const sorted = [...this.conversation].sort((a, b) => b.importance - a.importance);
    
    const optimized = [];
    let usedTokens = 0;

    for (const msg of sorted) {
      if (usedTokens + msg.tokens <= availableTokens) {
        optimized.push(msg);
        usedTokens += msg.tokens;
      } else {
        // Try to compress this message to fit
        const remainingTokens = availableTokens - usedTokens;
        if (remainingTokens > 10) { // Only compress if we have meaningful space
          const compressed = this.compressMessage(msg, remainingTokens);
          if (compressed.tokens <= remainingTokens) {
            optimized.push(compressed);
            break;
          }
        }
        break;
      }
    }

    // Return in chronological order
    return optimized.sort((a, b) => a.timestamp - b.timestamp)
                   .map(msg => ({ role: msg.role, content: msg.content }));
  }

  getTokenAnalytics() {
    return {
      ...this.tokenUsage,
      efficiency: (this.tokenUsage.total / this.maxTokens * 100).toFixed(1) + '%',
      messageCount: this.conversation.length,
      averageTokensPerMessage: this.conversation.length > 0 ? 
        (this.tokenUsage.total / this.conversation.length).toFixed(1) : 0
    };
  }
}

// Usage example
const advancedTokenManager = new AdvancedTokenManager(500, {
  removal: 'selective_compression',
  prioritization: 'importance'
});

// Add messages with different importance levels
advancedTokenManager.addMessage('system', 'You are a helpful AI assistant.', { importance: 3.0 });
advancedTokenManager.addMessage('user', 'Can you explain quantum computing?', { importance: 2.5 });
advancedTokenManager.addMessage('assistant', 'Quantum computing uses qubits instead of bits...', { importance: 2.0 });
advancedTokenManager.addMessage('user', 'That sounds interesting. Tell me more.', { importance: 1.5 });

console.log('Token Analytics:', advancedTokenManager.getTokenAnalytics());
console.log('Optimized Context:', advancedTokenManager.getOptimizedContext(300));
```

### 3. Token-Aware Memory with Chunking

```javascript
class ChunkedTokenMemory {
  constructor(maxChunkTokens = 2000, maxTotalTokens = 8000) {
    this.maxChunkTokens = maxChunkTokens;
    this.maxTotalTokens = maxTotalTokens;
    this.chunks = []; // Array of message chunks
    this.chunkSummaries = []; // Summaries of each chunk
  }

  countTokens(text) {
    // Simple token counter
    return Math.ceil(text.split(/\s+/).length * 1.3);
  }

  addMessage(role, content) {
    const tokens = this.countTokens(content);
    
    // If message is too large for a single chunk, split it
    if (tokens > this.maxChunkTokens) {
      this.handleOversizedMessage(role, content, tokens);
      return;
    }

    // Try to add to current chunk
    const currentChunk = this.getCurrentChunk();
    if (currentChunk.tokens + tokens <= this.maxChunkTokens) {
      currentChunk.messages.push({ role, content });
      currentChunk.tokens += tokens;
    } else {
      // Create new chunk
      this.createNewChunk([{ role, content }], tokens);
    }

    this.enforceTotalLimit();
  }

  getCurrentChunk() {
    if (this.chunks.length === 0) {
      return this.createNewChunk([], 0);
    }
    return this.chunks[this.chunks.length - 1];
  }

  createNewChunk(messages, initialTokens = 0) {
    const newChunk = {
      id: this.chunks.length + 1,
      messages,
      tokens: initialTokens,
      timestamp: Date.now(),
      summary: ''
    };
    
    this.chunks.push(newChunk);
    return newChunk;
  }

  handleOversizedMessage(role, content, totalTokens) {
    // Split oversized message into parts
    const words = content.split(' ');
    const targetWordCount = Math.floor((this.maxChunkTokens / 1.3) * 0.8); // Safe estimate
    
    for (let i = 0; i < words.length; i += targetWordCount) {
      const chunkWords = words.slice(i, i + targetWordCount);
      const chunkContent = chunkWords.join(' ');
      const chunkTokens = this.countTokens(chunkContent);
      
      this.addMessage(role, chunkContent);
    }
  }

  enforceTotalLimit() {
    let totalTokens = this.chunks.reduce((sum, chunk) => sum + chunk.tokens, 0);
    
    while (totalTokens > this.maxTotalTokens && this.chunks.length > 1) {
      const oldestChunk = this.chunks.shift(); // Remove oldest chunk
      totalTokens -= oldestChunk.tokens;
      
      // Preserve summary of removed chunk
      if (oldestChunk.messages.length > 0) {
        this.chunkSummaries.push(this.generateChunkSummary(oldestChunk));
      }
    }
  }

  generateChunkSummary(chunk) {
    // Simple summary generation - in practice, use NLP
    const content = chunk.messages.map(msg => msg.content).join(' ');
    const sentences = content.split(/[.!?]+/).filter(s => s.trim());
    
    return {
      chunkId: chunk.id,
      summary: sentences.slice(0, 2).join('. ') + '...',
      messageCount: chunk.messages.length,
      originalTokens: chunk.tokens
    };
  }

  getContext(maxTokens = 4000) {
    const context = {
      summaries: [...this.chunkSummaries],
      recentMessages: []
    };

    let usedTokens = context.summaries.reduce((sum, s) => sum + this.countTokens(s.summary), 0);
    
    // Add messages from recent chunks (newest first)
    for (let i = this.chunks.length - 1; i >= 0; i--) {
      const chunk = this.chunks[i];
      
      // Add messages from this chunk (in reverse order within chunk)
      for (let j = chunk.messages.length - 1; j >= 0; j--) {
        const message = chunk.messages[j];
        const messageTokens = this.countTokens(message.content);
        
        if (usedTokens + messageTokens <= maxTokens) {
          context.recentMessages.unshift(message); // Add to beginning to maintain order
          usedTokens += messageTokens;
        } else {
          break;
        }
      }
      
      if (usedTokens >= maxTokens * 0.9) break; // Leave some buffer
    }

    // Ensure chronological order
    context.recentMessages.sort((a, b) => {
      // This would need timestamp tracking for proper sorting
      return 0; // Simplified
    });

    return context;
  }

  getMemoryStats() {
    const totalTokens = this.chunks.reduce((sum, chunk) => sum + chunk.tokens, 0) +
                      this.chunkSummaries.reduce((sum, s) => sum + s.originalTokens, 0);
    
    return {
      totalChunks: this.chunks.length,
      totalSummaries: this.chunkSummaries.length,
      currentTokens: totalTokens,
      maxTokens: this.maxTotalTokens,
      utilization: ((totalTokens / this.maxTotalTokens) * 100).toFixed(1) + '%'
    };
  }
}

// Usage example
const chunkedMemory = new ChunkedTokenMemory(1000, 4000);

// Add multiple messages
for (let i = 1; i <= 10; i++) {
  chunkedMemory.addMessage('user', `This is message number ${i} about various topics including AI and machine learning.`);
  chunkedMemory.addMessage('assistant', `Response to message ${i} with detailed information.`);
}

console.log('Chunked Memory Stats:', chunkedMemory.getMemoryStats());
console.log('Context:', chunkedMemory.getContext(2000));
```

### 4. Practical Implementation: AI Conversation Manager

```javascript
class AIConversationManager {
  constructor(modelConfig = { maxTokens: 4000, tokenBuffer: 200 }) {
    this.modelConfig = modelConfig;
    this.tokenManager = new AdvancedTokenManager(
      modelConfig.maxTokens - modelConfig.tokenBuffer
    );
    this.messageHistory = [];
  }

  async processUserMessage(userMessage, systemPrompt = null) {
    // Add system prompt if provided
    if (systemPrompt && !this.hasSystemPrompt()) {
      this.tokenManager.addMessage('system', systemPrompt, { 
        importance: 3.0,
        metadata: { type: 'system_prompt' }
      });
    }

    // Add user message
    this.tokenManager.addMessage('user', userMessage, {
      importance: 2.5,
      metadata: { type: 'user_query' }
    });

    // Get optimized context for AI
    const context = this.tokenManager.getOptimizedContext();
    
    // Simulate AI response (in real implementation, call AI API)
    const aiResponse = await this.generateAIResponse(context, userMessage);
    
    // Add AI response to memory
    this.tokenManager.addMessage('assistant', aiResponse, {
      importance: 2.0,
      metadata: { type: 'ai_response' }
    });

    this.messageHistory.push({
      user: userMessage,
      assistant: aiResponse,
      timestamp: Date.now(),
      tokensUsed: this.tokenManager.getTokenAnalytics().total
    });

    return {
      response: aiResponse,
      tokenUsage: this.tokenManager.getTokenAnalytics(),
      contextSize: context.length
    };
  }

  async generateAIResponse(context, userMessage) {
    // Simulated AI response generation
    // In practice, this would call GPT-4, Claude, etc.
    await new Promise(resolve => setTimeout(resolve, 100)); // Simulate API call
    
    return `I understand you're asking about "${userMessage.substring(0, 50)}...". Based on our conversation, I can provide information about this topic.`;
  }

  hasSystemPrompt() {
    return this.tokenManager.conversation.some(msg => 
      msg.metadata.type === 'system_prompt'
    );
  }

  compressConversation() {
    const analytics = this.tokenManager.getTokenAnalytics();
    
    if (analytics.total > this.modelConfig.maxTokens * 0.8) {
      // High token usage, apply aggressive compression
      this.tokenManager.conversation.forEach((msg, index) => {
        if (msg.importance < 2.0 && msg.tokens > 50) {
          const compressed = this.tokenManager.compressMessage(msg, Math.ceil(msg.tokens * 0.4));
          this.tokenManager.conversation[index] = compressed;
        }
      });
      
      this.tokenManager.recalculateTokenUsage();
    }
  }

  getConversationSummary() {
    const userMessages = this.tokenManager.conversation.filter(msg => msg.role === 'user');
    const recentTopics = userMessages.slice(-5).map(msg => 
      msg.content.substring(0, 100)
    );

    return {
      totalMessages: this.messageHistory.length,
      recentTopics,
      tokenEfficiency: this.tokenManager.getTokenAnalytics().efficiency,
      currentContextSize: this.tokenManager.conversation.length
    };
  }

  exportConversation() {
    return {
      messages: this.messageHistory,
      tokenAnalytics: this.tokenManager.getTokenAnalytics(),
      configuration: this.modelConfig
    };
  }
}

// Usage example
const conversationManager = new AIConversationManager({
  maxTokens: 2000,
  tokenBuffer: 150
});

// Simulate conversation
async function simulateConversation() {
  const messages = [
    "Hello, how are you?",
    "Can you explain artificial intelligence?",
    "What are the different types of machine learning?",
    "How does deep learning work?",
    "What's the difference between AI and ML?"
  ];

  for (const message of messages) {
    const result = await conversationManager.processUserMessage(
      message,
      "You are a helpful AI assistant specializing in technology topics."
    );
    
    console.log(`User: ${message}`);
    console.log(`Assistant: ${result.response.substring(0, 100)}...`);
    console.log(`Tokens: ${result.tokenUsage.total}/${result.tokenUsage.tokenLimit}`);
    console.log('---');
  }

  console.log('Conversation Summary:', conversationManager.getConversationSummary());
}

// simulateConversation();
```

## Key Benefits of Token-Aware Memory Management

1. **Optimized Context Windows**: Maximizes useful information within token limits
2. **Intelligent Prioritization**: Keeps most important information accessible
3. **Efficient Compression**: Reduces token usage without losing critical context
4. **Model Compliance**: Ensures compatibility with AI model token limits
5. **Cost Optimization**: Reduces API costs by minimizing token usage

## Best Practices

- **Set appropriate buffers** for AI responses
- **Balance importance scoring** with recency
- **Implement multiple compression strategies**
- **Monitor token usage efficiency**
- **Use chunking for very long conversations**

Token-aware memory management is essential for building efficient, cost-effective AI applications that can handle long conversations while staying within model constraints.