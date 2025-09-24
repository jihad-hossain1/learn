# Memory Consolidation Patterns

Memory consolidation refers to the process of transforming transient, short-term memories into stable, long-term representations. In AI systems, this involves strategic management of conversation history and knowledge retention.

## Core Consolidation Concepts

### Why Consolidate Memories?
- **Reduce cognitive load** on AI context windows
- **Preserve important information** across sessions
- **Prevent information overload**
- **Maintain conversation coherence** over long interactions

### Consolidation Strategies
- **Summarization**: Condense conversations into key points
- **Extraction**: Pull out important facts and entities
- **Forgetting**: Remove irrelevant or redundant information
- **Reorganization**: Structure knowledge for efficient retrieval

## JavaScript Implementation Examples

### 1. Base Memory Consolidator Class

```javascript
class MemoryConsolidator {
  constructor(options = {}) {
    this.options = {
      consolidationThreshold: options.consolidationThreshold || 10, // messages
      timeThreshold: options.timeThreshold || 3600000, // 1 hour in ms
      importanceThreshold: options.importanceThreshold || 0.7,
      maxSummaryLength: options.maxSummaryLength || 500,
      ...options
    };
    this.consolidationHistory = [];
  }

  async needsConsolidation(conversationState) {
    const messageCount = conversationState.messages.length;
    const timeSinceLastConsolidation = this.getTimeSinceLastConsolidation();
    const memoryLoad = this.calculateMemoryLoad(conversationState);

    return (
      messageCount >= this.options.consolidationThreshold ||
      timeSinceLastConsolidation >= this.options.timeThreshold ||
      memoryLoad > 0.8
    );
  }

  async consolidate(conversationState, strategy = 'summary') {
    if (!(await this.needsConsolidation(conversationState))) {
      return conversationState; // No consolidation needed
    }

    switch (strategy) {
      case 'summary':
        return await this.consolidateWithSummary(conversationState);
      case 'extraction':
        return await this.consolidateWithExtraction(conversationState);
      case 'hierarchical':
        return await this.consolidateHierarchically(conversationState);
      case 'selective':
        return await this.consolidateSelectively(conversationState);
      default:
        throw new Error(`Unknown consolidation strategy: ${strategy}`);
    }
  }

  calculateMemoryLoad(conversationState) {
    const totalTokens = conversationState.messages.reduce(
      (sum, msg) => sum + (msg.tokenCount || 0), 0
    );
    const maxTokens = this.options.maxContextTokens || 4000;
    return totalTokens / maxTokens;
  }

  getTimeSinceLastConsolidation() {
    if (this.consolidationHistory.length === 0) {
      return Infinity;
    }
    const lastConsolidation = this.consolidationHistory[this.consolidationHistory.length - 1];
    return Date.now() - lastConsolidation.timestamp;
  }

  recordConsolidation(strategy, beforeState, afterState) {
    this.consolidationHistory.push({
      timestamp: Date.now(),
      strategy,
      messagesBefore: beforeState.messages.length,
      messagesAfter: afterState.messages.length,
      tokensBefore: this.calculateTotalTokens(beforeState),
      tokensAfter: this.calculateTotalTokens(afterState),
      reductionRatio: this.calculateReductionRatio(beforeState, afterState)
    });

    // Keep history manageable
    if (this.consolidationHistory.length > 100) {
      this.consolidationHistory.shift();
    }
  }

  calculateTotalTokens(conversationState) {
    return conversationState.messages.reduce((sum, msg) => sum + (msg.tokenCount || 0), 0);
  }

  calculateReductionRatio(beforeState, afterState) {
    const beforeTokens = this.calculateTotalTokens(beforeState);
    const afterTokens = this.calculateTotalTokens(afterState);
    return beforeTokens > 0 ? (beforeTokens - afterTokens) / beforeTokens : 0;
  }

  getConsolidationMetrics() {
    const totalConsolidations = this.consolidationHistory.length;
    if (totalConsolidations === 0) return null;

    const averageReduction = this.consolidationHistory.reduce(
      (sum, entry) => sum + entry.reductionRatio, 0
    ) / totalConsolidations;

    return {
      totalConsolidations,
      averageReduction: `${(averageReduction * 100).toFixed(1)}%`,
      mostUsedStrategy: this.getMostUsedStrategy(),
      lastConsolidation: this.consolidationHistory[this.consolidationHistory.length - 1]
    };
  }

  getMostUsedStrategy() {
    const strategyCount = {};
    this.consolidationHistory.forEach(entry => {
      strategyCount[entry.strategy] = (strategyCount[entry.strategy] || 0) + 1;
    });

    return Object.entries(strategyCount)
      .sort(([,a], [,b]) => b - a)[0][0];
  }
}
```

### 2. Summary-Based Consolidation

```javascript
class SummaryConsolidator extends MemoryConsolidator {
  constructor(summarizer, options = {}) {
    super(options);
    this.summarizer = summarizer;
  }

  async consolidateWithSummary(conversationState) {
    const messagesToConsolidate = this.selectMessagesForConsolidation(conversationState);
    
    if (messagesToConsolidate.length === 0) {
      return conversationState;
    }

    // Generate summary of selected messages
    const summary = await this.generateSummary(messagesToConsolidate);
    
    // Create new conversation state with summary
    const consolidatedState = this.createConsolidatedState(
      conversationState,
      messagesToConsolidate,
      summary
    );

    this.recordConsolidation('summary', conversationState, consolidatedState);
    return consolidatedState;
  }

  selectMessagesForConsolidation(conversationState) {
    const messages = conversationState.messages;
    
    // Always keep recent messages (last 2-3 exchanges)
    const keepRecentCount = this.options.keepRecentMessages || 4;
    const messagesToKeep = messages.slice(-keepRecentCount);
    const messagesToConsolidate = messages.slice(0, -keepRecentCount);

    // Don't consolidate if too few messages
    if (messagesToConsolidate.length < this.options.minConsolidationMessages || 3) {
      return [];
    }

    return messagesToConsolidate;
  }

  async generateSummary(messages) {
    const conversationText = messages.map(msg => 
      `${msg.role}: ${msg.content}`
    ).join('\n');

    const summary = await this.summarizer.summarize(conversationText, {
      maxLength: this.options.maxSummaryLength,
      focusPoints: this.extractKeyTopics(messages)
    });

    return {
      type: 'summary',
      content: summary,
      originalMessageCount: messages.length,
      consolidatedAt: new Date(),
      sourceMessages: messages.map(msg => msg.id)
    };
  }

  extractKeyTopics(messages) {
    // Extract frequent nouns and important entities
    const content = messages.map(msg => msg.content).join(' ');
    const words = content.toLowerCase().split(/\W+/).filter(word => word.length > 3);
    
    const wordFreq = {};
    words.forEach(word => {
      wordFreq[word] = (wordFreq[word] || 0) + 1;
    });

    return Object.entries(wordFreq)
      .sort(([,a], [,b]) => b - a)
      .slice(0, 5)
      .map(([word]) => word);
  }

  createConsolidatedState(originalState, consolidatedMessages, summary) {
    const keptMessages = originalState.messages.slice(-this.options.keepRecentMessages || 4);
    
    return {
      ...originalState,
      messages: [
        summary, // Summary becomes first message
        ...keptMessages
      ],
      summaries: [
        ...(originalState.summaries || []),
        summary
      ],
      lastConsolidation: new Date()
    };
  }
}

// Example summarizer implementation
class TextSummarizer {
  async summarize(text, options = {}) {
    // In practice, use an AI model or sophisticated algorithm
    // This is a simplified implementation
    
    const sentences = text.split(/[.!?]+/).filter(s => s.trim().length > 0);
    if (sentences.length <= 3) {
      return text; // Too short to summarize meaningfully
    }

    // Simple extractive summarization - take key sentences
    const importantSentences = this.scoreSentences(sentences, options.focusPoints);
    const summarySentences = importantSentences
      .sort((a, b) => b.score - a.score)
      .slice(0, Math.min(3, Math.ceil(sentences.length * 0.3)))
      .sort((a, b) => a.originalIndex - b.originalIndex) // Maintain order
      .map(item => item.sentence.trim());

    let summary = summarySentences.join('. ') + '.';
    
    // Trim to max length if needed
    if (summary.length > options.maxLength) {
      summary = summary.substring(0, options.maxLength - 3) + '...';
    }

    return summary;
  }

  scoreSentences(sentences, focusPoints = []) {
    return sentences.map((sentence, index) => {
      let score = 0;
      
      // Position scoring (first and last sentences are often important)
      if (index === 0 || index === sentences.length - 1) {
        score += 2;
      }
      
      // Length scoring (medium-length sentences are often more informative)
      const wordCount = sentence.split(/\s+/).length;
      if (wordCount > 5 && wordCount < 25) {
        score += 1;
      }
      
      // Focus points scoring
      if (focusPoints) {
        focusPoints.forEach(point => {
          if (sentence.toLowerCase().includes(point)) {
            score += 3;
          }
        });
      }
      
      // Question scoring (questions often indicate important content)
      if (sentence.includes('?')) {
        score += 1;
      }
      
      return {
        sentence,
        originalIndex: index,
        score
      };
    });
  }
}
```

### 3. Entity-Based Consolidation

```javascript
class EntityConsolidator extends MemoryConsolidator {
  constructor(entityExtractor, options = {}) {
    super(options);
    this.entityExtractor = entityExtractor;
    this.knownEntities = new Map();
  }

  async consolidateWithExtraction(conversationState) {
    const entities = await this.extractEntitiesFromConversation(conversationState);
    const importantMessages = this.identifyImportantMessages(conversationState, entities);
    
    // Create entity-based summary
    const entitySummary = this.createEntitySummary(entities);
    
    const consolidatedState = this.createEntityConsolidatedState(
      conversationState,
      importantMessages,
      entitySummary,
      entities
    );

    this.recordConsolidation('extraction', conversationState, consolidatedState);
    return consolidatedState;
  }

  async extractEntitiesFromConversation(conversationState) {
    const allContent = conversationState.messages
      .map(msg => msg.content)
      .join('\n');

    const newEntities = await this.entityExtractor.extract(allContent);
    
    // Merge with known entities
    this.mergeEntities(newEntities);
    
    return this.knownEntities;
  }

  mergeEntities(newEntities) {
    newEntities.forEach(entity => {
      const existing = this.knownEntities.get(entity.name);
      
      if (existing) {
        // Update existing entity
        existing.mentionCount += entity.mentionCount;
        existing.lastMentioned = new Date();
        existing.contexts = [...existing.contexts, ...entity.contexts].slice(-5); // Keep recent
        existing.attributes = { ...existing.attributes, ...entity.attributes };
      } else {
        // Add new entity
        this.knownEntities.set(entity.name, {
          ...entity,
          firstMentioned: new Date(),
          lastMentioned: new Date()
        });
      }
    });
  }

  identifyImportantMessages(conversationState, entities) {
    return conversationState.messages.filter(message => {
      // Messages that mention important entities
      const entityMentions = this.countEntityMentions(message.content, entities);
      
      // Messages with high importance signals
      const importanceSignals = this.calculateImportanceSignals(message);
      
      return entityMentions > 0 || importanceSignals > this.options.importanceThreshold;
    });
  }

  countEntityMentions(content, entities) {
    let count = 0;
    for (const entityName of entities.keys()) {
      if (content.toLowerCase().includes(entityName.toLowerCase())) {
        count++;
      }
    }
    return count;
  }

  calculateImportanceSignals(message) {
    let score = 0;
    
    // User questions are important
    if (message.role === 'user' && message.content.includes('?')) {
      score += 0.3;
    }
    
    // Long messages often contain important information
    const wordCount = message.content.split(/\s+/).length;
    if (wordCount > 20) {
      score += 0.2;
    }
    
    // Messages with specific keywords
    const importantKeywords = ['important', 'remember', 'key', 'critical', 'essential'];
    if (importantKeywords.some(keyword => message.content.toLowerCase().includes(keyword))) {
      score += 0.3;
    }
    
    return Math.min(score, 1.0);
  }

  createEntitySummary(entities) {
    const importantEntities = Array.from(entities.values())
      .filter(entity => entity.mentionCount >= 2 || entity.importance > 0.7)
      .slice(0, 5); // Top 5 entities

    const summaryPoints = importantEntities.map(entity => 
      `${entity.name} (mentioned ${entity.mentionCount} times): ${this.describeEntity(entity)}`
    );

    return {
      type: 'entity_summary',
      content: `Key entities discussed: ${summaryPoints.join('; ')}`,
      entities: importantEntities,
      consolidatedAt: new Date()
    };
  }

  describeEntity(entity) {
    if (entity.attributes && Object.keys(entity.attributes).length > 0) {
      const attributes = Object.entries(entity.attributes)
        .slice(0, 2)
        .map(([key, value]) => `${key}: ${value}`)
        .join(', ');
      return attributes;
    }
    
    return `mentioned in context of ${entity.contexts.slice(-1)[0]?.substring(0, 50)}...`;
  }

  createEntityConsolidatedState(originalState, importantMessages, entitySummary, entities) {
    return {
      ...originalState,
      messages: [
        entitySummary,
        ...importantMessages.slice(-this.options.keepImportantMessages || 6)
      ],
      entities: entities,
      entitySummaries: [
        ...(originalState.entitySummaries || []),
        entitySummary
      ]
    };
  }
}

// Example entity extractor
class EntityExtractor {
  async extract(text) {
    // Simple rule-based entity extraction
    // In practice, use NER models or more sophisticated approaches
    
    const entities = new Map();
    
    // Extract proper nouns (simplified)
    const words = text.split(/\s+/);
    const properNouns = words.filter(word => 
      word.length > 2 && /^[A-Z]/.test(word) && !this.isCommonWord(word)
    );

    // Count occurrences and extract context
    properNouns.forEach((noun, index) => {
      const cleanNoun = noun.replace(/[^\w]/g, '');
      if (cleanNoun.length < 2) return;

      const existing = entities.get(cleanNoun) || {
        name: cleanNoun,
        type: 'entity',
        mentionCount: 0,
        contexts: [],
        attributes: {}
      };

      existing.mentionCount++;
      
      // Extract context (surrounding words)
      const start = Math.max(0, index - 5);
      const end = Math.min(words.length, index + 6);
      const context = words.slice(start, end).join(' ');
      existing.contexts.push(context);
      
      entities.set(cleanNoun, existing);
    });

    return entities;
  }

  isCommonWord(word) {
    const commonWords = new Set(['The', 'This', 'That', 'There', 'What', 'How', 'Why']);
    return commonWords.has(word);
  }
}
```

### 4. Hierarchical Consolidation

```javascript
class HierarchicalConsolidator extends MemoryConsolidator {
  constructor(options = {}) {
    super(options);
    this.levels = options.levels || [
      { name: 'detailed', maxMessages: 20, maxAge: 3600000 }, // 1 hour
      { name: 'summarized', maxMessages: 10, maxAge: 86400000 }, // 1 day
      { name: 'archival', maxMessages: 5, maxAge: 604800000 } // 1 week
    ];
  }

  async consolidateHierarchically(conversationState) {
    const currentTime = Date.now();
    const messagesByLevel = this.categorizeMessagesByAge(conversationState.messages, currentTime);
    
    const consolidatedLevels = await this.consolidateEachLevel(messagesByLevel);
    
    const finalState = this.mergeConsolidatedLevels(
      conversationState,
      consolidatedLevels,
      currentTime
    );

    this.recordConsolidation('hierarchical', conversationState, finalState);
    return finalState;
  }

  categorizeMessagesByAge(messages, currentTime) {
    const levels = this.levels.map(level => ({ ...level, messages: [] }));
    
    messages.forEach(message => {
      const messageAge = currentTime - new Date(message.timestamp || currentTime).getTime();
      
      for (let i = levels.length - 1; i >= 0; i--) {
        if (messageAge <= levels[i].maxAge) {
          levels[i].messages.push(message);
          break;
        }
      }
    });
    
    return levels;
  }

  async consolidateEachLevel(levels) {
    const consolidatedLevels = [];
    
    for (let i = 0; i < levels.length; i++) {
      const level = levels[i];
      
      if (level.messages.length <= level.maxMessages) {
        consolidatedLevels.push(level);
        continue;
      }
      
      // Need to consolidate this level
      const consolidationMethod = i === 0 ? 'extraction' : 'summary';
      const consolidated = await this.consolidateLevel(
        level.messages,
        consolidationMethod,
        level.name
      );
      
      consolidatedLevels.push({
        ...level,
        messages: consolidated
      });
    }
    
    return consolidatedLevels;
  }

  async consolidateLevel(messages, method, levelName) {
    switch (method) {
      case 'summary':
        const summarizer = new TextSummarizer();
        const content = messages.map(msg => msg.content).join('\n');
        const summary = await summarizer.summarize(content, {
          maxLength: this.options.maxSummaryLength
        });
        
        return [{
          type: 'level_summary',
          level: levelName,
          content: summary,
          originalCount: messages.length,
          consolidatedAt: new Date()
        }];
        
      case 'extraction':
        const importantMessages = messages.filter(msg => 
          this.isMessageImportant(msg, levelName)
        ).slice(-this.options.keepImportantMessages || 8);
        
        return importantMessages;
        
      default:
        return messages.slice(-this.options.keepRecentMessages || 10);
    }
  }

  isMessageImportant(message, levelName) {
    const importanceFactors = {
      detailed: ['?', 'important', 'key', 'critical'],
      summarized: ['decision', 'conclusion', 'agreement', 'action'],
      archival: ['summary', 'result', 'outcome', 'final']
    };
    
    const factors = importanceFactors[levelName] || [];
    return factors.some(factor => message.content.toLowerCase().includes(factor));
  }

  mergeConsolidatedLevels(originalState, consolidatedLevels, currentTime) {
    const allMessages = consolidatedLevels.flatMap(level => level.messages);
    
    // Add metadata about the consolidation
    const consolidationInfo = {
      type: 'hierarchical_consolidation',
      timestamp: currentTime,
      levels: consolidatedLevels.map(level => ({
        name: level.name,
        messageCount: level.messages.length
      }))
    };
    
    return {
      ...originalState,
      messages: allMessages,
      consolidationHistory: [
        ...(originalState.consolidationHistory || []),
        consolidationInfo
      ],
      memoryHierarchy: consolidatedLevels.reduce((acc, level) => {
        acc[level.name] = level.messages;
        return acc;
      }, {})
    };
  }

  getLevelEfficiency() {
    const metrics = {};
    this.levels.forEach(level => {
      // Calculate how well each level is managing its memory allocation
      // This would track actual usage vs allocated capacity
    });
    return metrics;
  }
}
```

### 5. Selective Consolidation with Importance Scoring

```javascript
class SelectiveConsolidator extends MemoryConsolidator {
  constructor(importanceScorer, options = {}) {
    super(options);
    this.importanceScorer = importanceScorer;
    this.learningRate = options.learningRate || 0.1;
  }

  async consolidateSelectively(conversationState) {
    // Score importance of each message
    const scoredMessages = await this.scoreMessagesImportance(conversationState.messages);
    
    // Group messages by topic/cluster
    const clusters = this.clusterMessages(scoredMessages);
    
    // Select representative messages from each cluster
    const selectedMessages = this.selectRepresentativeMessages(clusters);
    
    // Create continuity bridges for removed messages
    const continuityMessages = this.createContinuityBridges(selectedMessages, scoredMessages);
    
    const consolidatedState = this.createSelectiveConsolidatedState(
      conversationState,
      selectedMessages,
      continuityMessages
    );

    this.recordConsolidation('selective', conversationState, consolidatedState);
    return consolidatedState;
  }

  async scoreMessagesImportance(messages) {
    const scoredMessages = [];
    
    for (const message of messages) {
      const importance = await this.importanceScorer.score(message);
      scoredMessages.push({
        ...message,
        importance,
        retentionScore: this.calculateRetentionScore(message, importance)
      });
    }
    
    return scoredMessages;
  }

  calculateRetentionScore(message, importance) {
    let score = importance;
    
    // Recency boost (newer messages are more important to keep)
    const ageHours = (Date.now() - new Date(message.timestamp).getTime()) / (1000 * 60 * 60);
    const recencyBoost = Math.max(0, 1 - (ageHours / 24)); // Linear decay over 24 hours
    score *= (1 + recencyBoost * 0.5);
    
    // Question boost (questions often need to be remembered)
    if (message.content.includes('?')) {
      score *= 1.3;
    }
    
    // Length factor (very short or very long messages might be less important)
    const wordCount = message.content.split(/\s+/).length;
    if (wordCount < 3) score *= 0.7; // Very short
    if (wordCount > 100) score *= 0.8; // Very long
    
    return score;
  }

  clusterMessages(scoredMessages) {
    // Simple time-based clustering
    const clusters = [];
    let currentCluster = [];
    const clusterThreshold = this.options.clusterTimeThreshold || 300000; // 5 minutes
    
    scoredMessages.forEach((message, index) => {
      if (currentCluster.length === 0) {
        currentCluster.push(message);
        return;
      }
      
      const lastMessage = currentCluster[currentCluster.length - 1];
      const timeDiff = new Date(message.timestamp).getTime() - 
                      new Date(lastMessage.timestamp).getTime();
      
      if (timeDiff <= clusterThreshold) {
        currentCluster.push(message);
      } else {
        clusters.push([...currentCluster]);
        currentCluster = [message];
      }
    });
    
    if (currentCluster.length > 0) {
      clusters.push(currentCluster);
    }
    
    return clusters;
  }

  selectRepresentativeMessages(clusters) {
    const selectedMessages = [];
    
    clusters.forEach(cluster => {
      if (cluster.length <= 2) {
        // Keep all messages in small clusters
        selectedMessages.push(...cluster);
        return;
      }
      
      // Select top messages by retention score
      const sortedCluster = cluster.sort((a, b) => b.retentionScore - a.retentionScore);
      const keepCount = Math.max(2, Math.ceil(cluster.length * 0.4)); // Keep 40% or 2 messages
      
      selectedMessages.push(...sortedCluster.slice(0, keepCount));
    });
    
    return selectedMessages.sort((a, b) => 
      new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
    );
  }

  createContinuityBridges(selectedMessages, allMessages) {
    const bridges = [];
    const selectedIds = new Set(selectedMessages.map(msg => msg.id));
    
    // Find gaps in the timeline where messages were removed
    for (let i = 0; i < allMessages.length - 1; i++) {
      const current = allMessages[i];
      const next = allMessages[i + 1];
      
      if (!selectedIds.has(current.id) && selectedIds.has(next.id)) {
        // Gap detected - create a bridge message
        const gapMessages = this.findMessagesInGap(allMessages, i);
        if (gapMessages.length > 0) {
          bridges.push(this.createBridgeMessage(gapMessages, next));
        }
      }
    }
    
    return bridges;
  }

  findMessagesInGap(messages, startIndex) {
    const gapMessages = [];
    let i = startIndex;
    
    while (i < messages.length && gapMessages.length < 5) { // Limit gap size
      if (messages[i].importance > 0.3) { // Only include somewhat important messages
        gapMessages.push(messages[i]);
      }
      i++;
    }
    
    return gapMessages;
  }

  createBridgeMessage(gapMessages, nextMessage) {
    const gapContent = gapMessages.map(msg => 
      msg.content.substring(0, 100)
    ).join('; ');
    
    return {
      type: 'continuity_bridge',
      content: `Earlier discussion included: ${gapContent}... leading to current topic.`,
      timestamp: new Date(nextMessage.timestamp).getTime() - 1, // Just before next message
      bridgedMessages: gapMessages.map(msg => msg.id),
      importance: 0.5 // Medium importance for continuity
    };
  }

  createSelectiveConsolidatedState(originalState, selectedMessages, continuityMessages) {
    // Combine selected messages with continuity bridges
    const allMessages = [...selectedMessages, ...continuityMessages]
      .sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());
    
    return {
      ...originalState,
      messages: allMessages,
      consolidationInfo: {
        strategy: 'selective',
        originalCount: originalState.messages.length,
        retainedCount: allMessages.length,
        reductionRatio: (originalState.messages.length - allMessages.length) / originalState.messages.length,
        timestamp: new Date()
      }
    };
  }

  async updateImportanceModel(feedback) {
    // Learn from user feedback to improve importance scoring
    // This is a simplified implementation
    if (feedback.messageId && feedback.importanceFeedback) {
      const adjustment = feedback.importanceFeedback === 'important' ? this.learningRate : -this.learningRate;
      // Update scoring model based on feedback
      this.importanceScorer.adjustWeights(feedback.messageId, adjustment);
    }
  }
}

// Importance scorer implementation
class ImportanceScorer {
  constructor() {
    this.featureWeights = {
      question: 0.3,
      length: 0.2,
      recency: 0.25,
      entityMention: 0.15,
      userIntent: 0.1
    };
  }

  async score(message) {
    let score = 0;
    
    // Question detection
    if (message.content.includes('?')) {
      score += this.featureWeights.question;
    }
    
    // Length factor (medium-length messages are often most important)
    const wordCount = message.content.split(/\s+/).length;
    const lengthScore = Math.max(0, 1 - Math.abs(wordCount - 15) / 50); // Peak around 15 words
    score += lengthScore * this.featureWeights.length;
    
    // Recency (already handled in retention score)
    score += 0.5 * this.featureWeights.recency; // Base recency score
    
    // Entity mentions
    const entityCount = (message.content.match(/[A-Z][a-z]+/g) || []).length;
    score += Math.min(entityCount / 5, 1) * this.featureWeights.entityMention;
    
    // User intent signals
    const intentWords = ['please', 'help', 'important', 'remember', 'need'];
    const hasIntent = intentWords.some(word => message.content.toLowerCase().includes(word));
    if (hasIntent) {
      score += this.featureWeights.userIntent;
    }
    
    return Math.min(score, 1.0);
  }

  adjustWeights(messageId, adjustment) {
    // Simple weight adjustment based on feedback
    // In practice, use more sophisticated machine learning
    Object.keys(this.featureWeights).forEach(feature => {
      this.featureWeights[feature] = Math.max(0.05, 
        Math.min(0.9, this.featureWeights[feature] + adjustment * 0.05)
      );
    });
    
    // Normalize weights
    this.normalizeWeights();
  }

  normalizeWeights() {
    const total = Object.values(this.featureWeights).reduce((sum, weight) => sum + weight, 0);
    Object.keys(this.featureWeights).forEach(feature => {
      this.featureWeights[feature] /= total;
    });
  }
}
```

### 6. Adaptive Consolidation System

```javascript
class AdaptiveConsolidationSystem {
  constructor(consolidators, options = {}) {
    this.consolidators = consolidators; // Array of different consolidation strategies
    this.options = {
      evaluationInterval: options.evaluationInterval || 10000, // 10 seconds
      strategyWeights: options.strategyWeights || {},
      learningEnabled: options.learningEnabled !== false,
      ...options
    };
    
    this.performanceMetrics = new Map();
    this.currentStrategy = consolidators[0];
    this.evaluationTimer = null;
  }

  async initialize() {
    // Start periodic evaluation of consolidation strategies
    this.startStrategyEvaluation();
  }

  async consolidate(conversationState, context = {}) {
    // Choose best strategy based on current context and performance
    const strategy = await this.selectOptimalStrategy(conversationState, context);
    
    const startTime = Date.now();
    try {
      const result = await strategy.consolidate(conversationState);
      this.recordPerformance(strategy, 'success', Date.now() - startTime, conversationState, result);
      return result;
    } catch (error) {
      this.recordPerformance(strategy, 'error', Date.now() - startTime, conversationState, null);
      // Fallback to another strategy
      return this.fallbackConsolidation(conversationState, strategy);
    }
  }

  async selectOptimalStrategy(conversationState, context) {
    if (!this.options.learningEnabled) {
      return this.currentStrategy;
    }

    const scores = await Promise.all(
      this.consolidators.map(consolidator => 
        this.scoreStrategy(consolidator, conversationState, context)
      )
    );

    const bestStrategyIndex = scores.reduce((bestIndex, score, index) => 
      score > scores[bestIndex] ? index : bestIndex, 0
    );

    this.currentStrategy = this.consolidators[bestStrategyIndex];
    return this.currentStrategy;
  }

  async scoreStrategy(consolidator, conversationState, context) {
    let score = 0;
    
    // Historical performance
    const performance = this.performanceMetrics.get(consolidator.constructor.name) || 
                       { successRate: 0.5, averageReduction: 0.3 };
    score += performance.successRate * 0.4;
    score += performance.averageReduction * 0.3;
    
    // Context suitability
    const contextScore = this.evaluateContextSuitability(consolidator, context);
    score += contextScore * 0.3;
    
    return score;
  }

  evaluateContextSuitability(consolidator, context) {
    let score = 0;
    
    // Conversation length suitability
    if (context.conversationLength) {
      if (context.conversationLength > 20 && consolidator.constructor.name === 'SummaryConsolidator') {
        score += 0.5; // Summary works well for long conversations
      }
      if (context.conversationLength < 10 && consolidator.constructor.name === 'SelectiveConsolidator') {
        score += 0.3; // Selective works well for shorter conversations
      }
    }
    
    // Domain suitability
    if (context.domain === 'technical' && consolidator.constructor.name === 'EntityConsolidator') {
      score += 0.2; // Entity consolidation good for technical content
    }
    
    return score;
  }

  async fallbackConsolidation(conversationState, failedStrategy) {
    const otherStrategies = this.consolidators.filter(s => s !== failedStrategy);
    
    for (const strategy of otherStrategies) {
      try {
        return await strategy.consolidate(conversationState);
      } catch (error) {
        console.warn(`Fallback strategy also failed: ${error.message}`);
      }
    }
    
    // Last resort: simple truncation
    return this.truncateConversation(conversationState);
  }

  truncateConversation(conversationState) {
    // Keep only most recent messages
    const keepCount = this.options.minKeepMessages || 5;
    return {
      ...conversationState,
      messages: conversationState.messages.slice(-keepCount),
      consolidationInfo: {
        strategy: 'truncation',
        reason: 'all_strategies_failed',
        timestamp: new Date()
      }
    };
  }

  recordPerformance(strategy, outcome, duration, beforeState, afterState) {
    const key = strategy.constructor.name;
    const current = this.performanceMetrics.get(key) || {
      totalRuns: 0,
      successfulRuns: 0,
      totalDuration: 0,
      totalReduction: 0
    };
    
    current.totalRuns++;
    current.totalDuration += duration;
    
    if (outcome === 'success') {
      current.successfulRuns++;
      if (afterState) {
        const reduction = this.calculateReductionRatio(beforeState, afterState);
        current.totalReduction += reduction;
      }
    }
    
    this.performanceMetrics.set(key, current);
  }

  calculateReductionRatio(beforeState, afterState) {
    const beforeTokens = beforeState.messages.reduce((sum, msg) => sum + (msg.tokenCount || 0), 0);
    const afterTokens = afterState.messages.reduce((sum, msg) => sum + (msg.tokenCount || 0), 0);
    return beforeTokens > 0 ? (beforeTokens - afterTokens) / beforeTokens : 0;
  }

  startStrategyEvaluation() {
    this.evaluationTimer = setInterval(() => {
      this.evaluateStrategies();
    }, this.options.evaluationInterval);
  }

  evaluateStrategies() {
    // Recalculate strategy weights based on recent performance
    this.consolidators.forEach(consolidator => {
      const metrics = this.performanceMetrics.get(consolidator.constructor.name);
      if (metrics && metrics.totalRuns > 0) {
        const successRate = metrics.successfulRuns / metrics.totalRuns;
        const averageReduction = metrics.totalReduction / metrics.successfulRuns || 0;
        
        // Update strategy weights (simplified)
        const newWeight = successRate * 0.6 + averageReduction * 0.4;
        this.options.strategyWeights[consolidator.constructor.name] = newWeight;
      }
    });
  }

  getSystemMetrics() {
    const metrics = {};
    
    this.consolidators.forEach(consolidator => {
      const perf = this.performanceMetrics.get(consolidator.constructor.name) || {};
      metrics[consolidator.constructor.name] = {
        successRate: perf.totalRuns > 0 ? (perf.successfulRuns / perf.totalRuns) : 0,
        averageDuration: perf.totalRuns > 0 ? (perf.totalDuration / perf.totalRuns) : 0,
        averageReduction: perf.successfulRuns > 0 ? (perf.totalReduction / perf.successfulRuns) : 0,
        totalRuns: perf.totalRuns || 0
      };
    });
    
    return metrics;
  }

  shutdown() {
    if (this.evaluationTimer) {
      clearInterval(this.evaluationTimer);
    }
  }
}
```

## Key Memory Consolidation Patterns

### 1. **Summary-Based Consolidation**
- **Pros**: Maintains context, reduces token usage significantly
- **Cons**: May lose specific details, summary quality critical
- **Best for**: Long conversations, narrative content

### 2. **Entity-Based Consolidation**
- **Pros**: Preserves important facts and relationships
- **Cons**: May lose conversational flow, requires good entity extraction
- **Best for**: Factual discussions, technical content

### 3. **Hierarchical Consolidation**
- **Pros**: Maintains multiple resolution levels, scalable
- **Cons**: Complex to implement, requires careful level management
- **Best for**: Very long-term conversations, archival systems

### 4. **Selective Consolidation**
- **Pros**: Preserves most important content, maintains flow
- **Cons**: Requires good importance scoring, may create gaps
- **Best for**: Balanced conversations, adaptive systems

## Best Practices

1. **Monitor consolidation quality**: Track what information is being lost
2. **Use multiple strategies**: Different conversations need different approaches
3. **Consider conversation type**: Technical vs casual conversations need different handling
4. **Preserve user intent**: Don't lose the original questions and requests
5. **Implement feedback loops**: Learn from user interactions to improve consolidation
6. **Balance compression and fidelity**: Don't over-optimize for size at the cost of usefulness

Memory consolidation is essential for building AI systems that can maintain long-term conversations without hitting context limits or becoming inefficient. The right consolidation strategy depends on your specific use case, conversation types, and performance requirements.