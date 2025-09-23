# Summary Memory Architecture

Summary Memory maintains a condensed version of conversation history instead of storing raw messages. It creates intelligent summaries that preserve key information while saving space.

## How Summary Memory Works

- **Condenses** multiple messages into concise summaries
- **Preserves key information** and context
- **Dramatically reduces** storage requirements
- **Maintains long-term context** beyond buffer limits

## JavaScript Implementation Examples

### 1. Basic Summary Memory Class

```javascript
class SummaryMemory {
  constructor() {
    this.summary = "";
    this.rawMessages = [];
    this.maxRawMessages = 10; // Keep some recent raw messages
  }

  addMessage(role, content) {
    this.rawMessages.push({ role, content, timestamp: Date.now() });

    // Maintain buffer of recent messages
    if (this.rawMessages.length > this.maxRawMessages) {
      this.rawMessages.shift();
    }

    // Update summary when significant new information is added
    if (this.shouldUpdateSummary()) {
      this.updateSummary();
    }
  }

  shouldUpdateSummary() {
    // Update summary every 3 messages or when conversation topic likely changed
    return this.rawMessages.length % 3 === 0;
  }

  updateSummary() {
    if (this.rawMessages.length === 0) return;

    const recentMessages = this.rawMessages.slice(-4); // Last 4 messages
    const newSummaryContent = this.generateSummary(recentMessages);

    if (this.summary) {
      this.summary = this.mergeSummaries(this.summary, newSummaryContent);
    } else {
      this.summary = newSummaryContent;
    }

    // Keep summary concise
    this.summary = this.compressSummary(this.summary);
  }

  generateSummary(messages) {
    const content = messages
      .map((msg) => `${msg.role}: ${msg.content}`)
      .join(". ");

    return `Recent discussion: ${content}.`;
  }

  mergeSummaries(oldSummary, newSummary) {
    // Simple merging - in practice, you'd use more sophisticated NLP
    return `${oldSummary} Then, ${newSummary}`;
  }

  compressSummary(summary) {
    // Basic compression - limit length
    const maxLength = 500;
    if (summary.length > maxLength) {
      return summary.substring(0, maxLength - 3) + "...";
    }
    return summary;
  }

  getContext() {
    const recentMessages = this.rawMessages.slice(-2); // Last 2 raw messages
    const recentContext = recentMessages
      .map((msg) => `${msg.role}: ${msg.content}`)
      .join("\n");

    return `${this.summary}\n\nRecent messages:\n${recentContext}`;
  }

  clear() {
    this.summary = "";
    this.rawMessages = [];
  }
}

// Usage example
const summaryMemory = new SummaryMemory();

summaryMemory.addMessage("user", "I want to learn about machine learning");
summaryMemory.addMessage(
  "assistant",
  "Machine learning involves training algorithms on data..."
);
summaryMemory.addMessage("user", "What are neural networks?");
summaryMemory.addMessage(
  "assistant",
  "Neural networks are inspired by the human brain..."
);

console.log("Summary Context:");
console.log(summaryMemory.getContext());
```

### 2. Advanced Summary Memory with Topic Tracking

```javascript
class TopicAwareSummaryMemory {
  constructor() {
    this.summaries = []; // Array of topic summaries
    this.currentTopic = null;
    this.messageBuffer = [];
    this.topicThreshold = 3; // Messages before creating a new topic
  }

  addMessage(role, content) {
    this.messageBuffer.push({ role, content, timestamp: Date.now() });

    const topics = this.detectTopics(content);
    this.manageTopics(topics);
  }

  detectTopics(content) {
    // Simple keyword-based topic detection
    const topicKeywords = {
      programming: ["code", "program", "function", "class", "variable"],
      weather: ["weather", "temperature", "rain", "sunny"],
      food: ["food", "restaurant", "cook", "recipe", "eat"],
    };

    const contentLower = content.toLowerCase();
    const detectedTopics = [];

    for (const [topic, keywords] of Object.entries(topicKeywords)) {
      if (keywords.some((keyword) => contentLower.includes(keyword))) {
        detectedTopics.push(topic);
      }
    }

    return detectedTopics.length > 0 ? detectedTopics : ["general"];
  }

  manageTopics(currentTopics) {
    if (!this.currentTopic || !this.hasTopicOverlap(currentTopics)) {
      // Topic shift detected, summarize current conversation
      if (this.messageBuffer.length >= this.topicThreshold) {
        this.createTopicSummary();
      }
      this.currentTopic = currentTopics[0];
    }
  }

  hasTopicOverlap(newTopics) {
    return newTopics.some((topic) => topic === this.currentTopic);
  }

  createTopicSummary() {
    if (this.messageBuffer.length === 0) return;

    const summary = {
      topic: this.currentTopic,
      summary: this.generateTopicSummary(this.messageBuffer),
      messageCount: this.messageBuffer.length,
      timestamp: Date.now(),
    };

    this.summaries.push(summary);
    this.messageBuffer = []; // Clear buffer after summarizing

    // Keep only recent summaries
    if (this.summaries.length > 5) {
      this.summaries.shift();
    }
  }

  generateTopicSummary(messages) {
    const content = messages.map((msg) => msg.content).join(". ");
    return `Discussion about ${this.currentTopic}: ${content.substring(
      0,
      200
    )}...`;
  }

  getFullContext() {
    const topicContext = this.summaries
      .map((s) => `[${s.topic}] ${s.summary}`)
      .join("\n");

    const recentMessages = this.messageBuffer
      .slice(-3)
      .map((msg) => `${msg.role}: ${msg.content}`)
      .join("\n");

    return `Previous topics:\n${topicContext}\n\nCurrent conversation:\n${recentMessages}`;
  }

  getSummaryStats() {
    return {
      totalTopics: this.summaries.length,
      currentTopic: this.currentTopic,
      bufferedMessages: this.messageBuffer.length,
      topics: this.summaries.map((s) => s.topic),
    };
  }
}

// Usage example
const topicMemory = new TopicAwareSummaryMemory();

topicMemory.addMessage("user", "How do I write a function in Python?");
topicMemory.addMessage(
  "assistant",
  "You can define functions using the def keyword..."
);
topicMemory.addMessage("user", "What about classes and objects?");
topicMemory.addMessage(
  "assistant",
  "Classes are defined with the class keyword..."
);
topicMemory.addMessage("user", "Is it going to rain today?");
topicMemory.addMessage("assistant", "I don't have real-time weather data...");

console.log("Topic-based Summary:");
console.log(topicMemory.getFullContext());
console.log("\nMemory Stats:", topicMemory.getSummaryStats());
```

### 3. Summary Memory with Importance Scoring

```javascript
class ImportanceWeightedSummaryMemory {
  constructor() {
    this.summary = "";
    this.messageLog = [];
    this.importanceWeights = {
      question: 2.0,
      explanation: 1.5,
      greeting: 0.5,
      confirmation: 0.7,
    };
  }

  addMessage(role, content, type = "general") {
    const importance = this.importanceWeights[type] || 1.0;

    this.messageLog.push({
      role,
      content,
      type,
      importance,
      timestamp: Date.now(),
    });

    this.updateSummaryIfNeeded();
  }

  updateSummaryIfNeeded() {
    // Update summary based on importance accumulation
    const recentImportance = this.calculateRecentImportance();

    if (recentImportance > 2.0 || this.messageLog.length % 4 === 0) {
      this.generateNewSummary();
    }
  }

  calculateRecentImportance() {
    const recentMessages = this.messageLog.slice(-5);
    return recentMessages.reduce((sum, msg) => sum + msg.importance, 0);
  }

  generateNewSummary() {
    const significantMessages = this.messageLog.filter(
      (msg) => msg.importance >= 1.0
    );

    if (significantMessages.length === 0) return;

    const summaryPoints = significantMessages
      .map((msg) => this.formatSummaryPoint(msg))
      .join(" ");

    this.summary = this.condenseSummary(summaryPoints);

    // Archive old low-importance messages
    this.messageLog = this.messageLog.filter(
      (msg) => msg.importance >= 0.8 || msg.timestamp > Date.now() - 300000 // Keep recent messages (5 minutes)
    );
  }

  formatSummaryPoint(message) {
    switch (message.type) {
      case "question":
        return `User asked: "${message.content}".`;
      case "explanation":
        return `Explained: ${message.content.substring(0, 100)}.`;
      default:
        return `${message.role}: ${message.content.substring(0, 80)}.`;
    }
  }

  condenseSummary(text) {
    // Simple condensation - in practice, use NLP techniques
    const sentences = text.split(". ");
    if (sentences.length > 5) {
      return sentences.slice(0, 3).join(". ") + "...";
    }
    return text;
  }

  getContextWithWeights() {
    const recentMessages = this.messageLog.slice(-3);
    const recentContext = recentMessages
      .map(
        (msg) => `[${msg.importance.toFixed(1)}] ${msg.role}: ${msg.content}`
      )
      .join("\n");

    return `Summary: ${this.summary}\n\nRecent messages (with importance):\n${recentContext}`;
  }

  searchMemory(query) {
    // Simple search through summaries and recent messages
    const queryLower = query.toLowerCase();

    const relevantMessages = this.messageLog.filter((msg) =>
      msg.content.toLowerCase().includes(queryLower)
    );

    return {
      query,
      foundInSummary: this.summary.toLowerCase().includes(queryLower),
      relevantMessages: relevantMessages.slice(-3),
    };
  }
}

// Usage example
const weightedMemory = new ImportanceWeightedSummaryMemory();

weightedMemory.addMessage("user", "What is JavaScript?", "question");
weightedMemory.addMessage(
  "assistant",
  "JavaScript is a programming language for web development...",
  "explanation"
);
weightedMemory.addMessage("user", "Hello there!", "greeting");
weightedMemory.addMessage("user", "Can you give examples?", "question");
weightedMemory.addMessage(
  "assistant",
  "Examples include functions, event handlers...",
  "explanation"
);

console.log("Weighted Summary Context:");
console.log(weightedMemory.getContextWithWeights());

console.log('\nSearch results for "examples":');
console.log(weightedMemory.searchMemory("examples"));
```

### 4. Hybrid Memory System (Buffer + Summary)

```javascript
class HybridConversationMemory {
  constructor(bufferSize = 5, summaryInterval = 3) {
    this.bufferMemory = new ConversationBufferMemory(bufferSize);
    this.summaryMemory = new SummaryMemory();
    this.summaryInterval = summaryInterval;
    this.messageCounter = 0;
  }

  addMessage(role, content) {
    // Add to buffer memory (always keeps recent messages)
    this.bufferMemory.addMessage(role, content);

    // Add to summary memory periodically
    this.summaryMemory.addMessage(role, content);

    this.messageCounter++;

    // Force summary update at intervals
    if (this.messageCounter % this.summaryInterval === 0) {
      this.summaryMemory.updateSummary();
    }
  }

  getFullContext() {
    const recentContext = this.bufferMemory
      .getConversationHistory()
      .map((msg) => `${msg.role}: ${msg.content}`)
      .join("\n");

    const summaryContext = this.summaryMemory.summary;

    return `Conversation Summary:\n${summaryContext}\n\nRecent Messages:\n${recentContext}`;
  }

  getContextForAI() {
    // Format context for AI model consumption
    const recentMessages = this.bufferMemory.getConversationHistory();
    const summary = this.summaryMemory.summary;

    return {
      system: `Previous conversation summary: ${summary}`,
      messages: recentMessages,
    };
  }

  clear() {
    this.bufferMemory.clear();
    this.summaryMemory.clear();
    this.messageCounter = 0;
  }

  getMemoryStats() {
    return {
      totalMessages: this.messageCounter,
      bufferSize: this.bufferMemory.getCurrentSize(),
      summaryLength: this.summaryMemory.summary.length,
      summaryExists: this.summaryMemory.summary.length > 0,
    };
  }
}

// Usage example
const hybridMemory = new HybridConversationMemory(4, 3);

// Simulate longer conversation
const conversation = [
  ["user", "I want to plan a trip to Japan"],
  ["assistant", "Great! What cities are you interested in?"],
  ["user", "Tokyo, Kyoto, and Osaka"],
  ["assistant", "Excellent choices. How long will you stay?"],
  ["user", "About 2 weeks"],
  ["assistant", "I recommend splitting time between cities"],
  ["user", "What about transportation?"],
  ["assistant", "The Shinkansen bullet train is efficient"],
];

conversation.forEach(([role, content]) => {
  hybridMemory.addMessage(role, content);
});

console.log("Hybrid Memory Context:");
console.log(hybridMemory.getFullContext());
console.log("\nMemory Statistics:", hybridMemory.getMemoryStats());
```

## Key Benefits of Summary Memory

1. **Efficient Storage**: Dramatically reduces memory footprint
2. **Long-term Context**: Maintains context beyond buffer limits
3. **Information Preservation**: Keeps important details while discarding fluff
4. **Scalability**: Handles long conversations effectively
5. **Topic Continuity**: Maintains conversation flow across sessions

## When to Use Summary Memory

- **Long conversations** with multiple topics
- **Resource-constrained environments**
- **Historical context** requirements
- **Multi-session conversations**
- **Information-dense dialogues**

Summary Memory is particularly useful for applications like customer support chatbots, therapeutic assistants, educational tools, and any system requiring long-term conversation context without the overhead of storing every message.
