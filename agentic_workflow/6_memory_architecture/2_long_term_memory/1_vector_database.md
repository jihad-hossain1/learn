# Long-Term Memory (LTM) with Vector Databases

Long-Term Memory using vector databases enables AI systems to store, retrieve, and reason over large amounts of information across multiple conversations and sessions.

## How Vector Database LTM Works

- **Stores embeddings** of conversations and knowledge
- **Semantic search** for relevant context retrieval
- **Persistent storage** across sessions
- **Scalable knowledge** management

## Core Concepts

### Embeddings
- Convert text to numerical vectors
- Capture semantic meaning
- Enable similarity search

### Vector Databases
- **Pinecone**: Fully-managed vector database
- **Weaviate**: Open-source vector search engine
- **Qdrant**: High-performance vector database

## JavaScript Implementation Examples

### 1. Base Vector Memory Class

```javascript
class VectorMemory {
  constructor(embeddingFunction, options = {}) {
    this.embeddingFunction = embeddingFunction;
    this.options = {
      similarityThreshold: options.similarityThreshold || 0.7,
      maxResults: options.maxResults || 5,
      namespace: options.namespace || 'default'
    };
    this.memories = [];
  }

  async addMemory(content, metadata = {}) {
    try {
      const embedding = await this.embeddingFunction(content);
      const memory = {
        id: this.generateId(),
        content,
        embedding,
        metadata: {
          timestamp: new Date(),
          type: 'memory',
          ...metadata
        },
        namespace: this.options.namespace
      };

      this.memories.push(memory);
      return memory;
    } catch (error) {
      console.error('Error adding memory:', error);
      throw error;
    }
  }

  async searchMemories(query, options = {}) {
    const queryEmbedding = await this.embeddingFunction(query);
    const similarities = [];

    for (const memory of this.memories) {
      const similarity = this.cosineSimilarity(queryEmbedding, memory.embedding);
      
      if (similarity >= (options.similarityThreshold || this.options.similarityThreshold)) {
        similarities.push({
          memory,
          similarity,
          score: similarity * this.calculateRecencyScore(memory.metadata.timestamp)
        });
      }
    }

    return similarities
      .sort((a, b) => b.score - a.score)
      .slice(0, options.maxResults || this.options.maxResults);
  }

  cosineSimilarity(vecA, vecB) {
    if (vecA.length !== vecB.length) {
      throw new Error('Vectors must have the same length');
    }

    let dotProduct = 0;
    let normA = 0;
    let normB = 0;

    for (let i = 0; i < vecA.length; i++) {
      dotProduct += vecA[i] * vecB[i];
      normA += vecA[i] * vecA[i];
      normB += vecB[i] * vecB[i];
    }

    return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB));
  }

  calculateRecencyScore(timestamp) {
    const now = new Date();
    const hoursDiff = (now - new Date(timestamp)) / (1000 * 60 * 60);
    return Math.max(0.1, 1 - (hoursDiff / 168)); // Decay over 1 week
  }

  generateId() {
    return `mem_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  getMemoryStats() {
    return {
      totalMemories: this.memories.length,
      namespaces: [...new Set(this.memories.map(m => m.namespace))],
      oldestMemory: this.memories.length > 0 ? 
        new Date(Math.min(...this.memories.map(m => new Date(m.metadata.timestamp)))) : null,
      newestMemory: this.memories.length > 0 ? 
        new Date(Math.max(...this.memories.map(m => new Date(m.metadata.timestamp)))) : null
    };
  }
}
```

### 2. Pinecone Integration

```javascript
class PineconeMemory extends VectorMemory {
  constructor(pineconeClient, indexName, embeddingFunction, options = {}) {
    super(embeddingFunction, options);
    this.pinecone = pineconeClient;
    this.indexName = indexName;
    this.index = pineconeClient.Index(indexName);
  }

  async addMemory(content, metadata = {}) {
    try {
      const embedding = await this.embeddingFunction(content);
      const memoryId = this.generateId();

      const upsertRequest = {
        vectors: [{
          id: memoryId,
          values: embedding,
          metadata: {
            content,
            timestamp: new Date().toISOString(),
            type: 'memory',
            namespace: this.options.namespace,
            ...metadata
          }
        }]
      };

      await this.index.upsert({ upsertRequest });
      
      const memory = {
        id: memoryId,
        content,
        embedding,
        metadata: upsertRequest.vectors[0].metadata
      };

      return memory;
    } catch (error) {
      console.error('Pinecone error:', error);
      throw error;
    }
  }

  async searchMemories(query, options = {}) {
    try {
      const queryEmbedding = await this.embeddingFunction(query);
      
      const queryRequest = {
        vector: queryEmbedding,
        topK: options.maxResults || this.options.maxResults,
        includeMetadata: true,
        includeValues: true,
        filter: options.filter || {}
      };

      if (this.options.namespace) {
        queryRequest.namespace = this.options.namespace;
      }

      const result = await this.index.query({ queryRequest });
      
      return result.matches.map(match => ({
        memory: {
          id: match.id,
          content: match.metadata.content,
          embedding: match.values,
          metadata: match.metadata
        },
        similarity: match.score,
        score: match.score * this.calculateRecencyScore(match.metadata.timestamp)
      }));
    } catch (error) {
      console.error('Pinecone search error:', error);
      throw error;
    }
  }

  async updateMemory(memoryId, updates) {
    try {
      // Pinecone doesn't support partial updates, so we need to fetch and re-upsert
      const fetchResult = await this.index.fetch({ ids: [memoryId] });
      const existing = fetchResult.vectors[memoryId];

      if (!existing) {
        throw new Error('Memory not found');
      }

      const updatedMetadata = {
        ...existing.metadata,
        ...updates,
        lastUpdated: new Date().toISOString()
      };

      const upsertRequest = {
        vectors: [{
          id: memoryId,
          values: existing.values,
          metadata: updatedMetadata
        }]
      };

      await this.index.upsert({ upsertRequest });
      return updatedMetadata;

    } catch (error) {
      console.error('Pinecone update error:', error);
      throw error;
    }
  }

  async deleteMemory(memoryId) {
    try {
      await this.index.delete1({ ids: [memoryId] });
      return true;
    } catch (error) {
      console.error('Pinecone delete error:', error);
      throw error;
    }
  }
}

// Pinecone initialization helper
class PineconeManager {
  constructor(apiKey, environment) {
    this.pinecone = new PineconeClient();
    this.initialized = false;
    this.apiKey = apiKey;
    this.environment = environment;
  }

  async initialize() {
    if (!this.initialized) {
      await this.pinecone.init({
        apiKey: this.apiKey,
        environment: this.environment
      });
      this.initialized = true;
    }
    return this.pinecone;
  }

  async getIndex(indexName, dimension = 1536) {
    await this.initialize();
    
    // Check if index exists, create if not
    const indexes = await this.pinecone.listIndexes();
    if (!indexes.includes(indexName)) {
      await this.pinecone.createIndex({
        createRequest: {
          name: indexName,
          dimension: dimension,
          metric: 'cosine'
        }
      });

      // Wait for index to be ready
      await new Promise(resolve => setTimeout(resolve, 30000));
    }

    return this.pinecone.Index(indexName);
  }
}
```

### 3. Weaviate Integration

```javascript
class WeaviateMemory extends VectorMemory {
  constructor(weaviateClient, className, embeddingFunction, options = {}) {
    super(embeddingFunction, options);
    this.weaviate = weaviateClient;
    this.className = className;
    this.initialized = false;
  }

  async initializeSchema() {
    if (this.initialized) return;

    const classObj = {
      class: this.className,
      properties: [
        {
          name: 'content',
          dataType: ['text']
        },
        {
          name: 'timestamp',
          dataType: ['date']
        },
        {
          name: 'type',
          dataType: ['text']
        },
        {
          name: 'namespace',
          dataType: ['text']
        }
      ],
      vectorizer: 'none' // We're providing our own vectors
    };

    try {
      await this.weaviate.schema.classCreator().withClass(classObj).do();
    } catch (error) {
      if (!error.message.includes('already exists')) {
        throw error;
      }
    }

    this.initialized = true;
  }

  async addMemory(content, metadata = {}) {
    await this.initializeSchema();

    const embedding = await this.embeddingFunction(content);
    const memoryId = this.generateId();

    const memoryData = {
      content,
      timestamp: new Date().toISOString(),
      type: metadata.type || 'memory',
      namespace: this.options.namespace,
      ...metadata
    };

    try {
      const result = await this.weaviate.data
        .creator()
        .withClassName(this.className)
        .withId(memoryId)
        .withProperties(memoryData)
        .withVector(embedding)
        .do();

      return {
        id: memoryId,
        content,
        embedding,
        metadata: memoryData
      };
    } catch (error) {
      console.error('Weaviate error:', error);
      throw error;
    }
  }

  async searchMemories(query, options = {}) {
    const queryEmbedding = await this.embeddingFunction(query);

    let whereFilter = {};
    if (options.filter) {
      whereFilter = this.buildWhereFilter(options.filter);
    }

    try {
      const result = await this.weaviate.graphql
        .get()
        .withClassName(this.className)
        .withFields('content timestamp type namespace _additional { id distance }')
        .withNearVector({
          vector: queryEmbedding,
          distance: options.maxDistance || 0.3
        })
        .withWhere(whereFilter)
        .withLimit(options.maxResults || this.options.maxResults)
        .do();

      return result.data.Get[this.className].map(item => ({
        memory: {
          id: item._additional.id,
          content: item.content,
          metadata: {
            timestamp: item.timestamp,
            type: item.type,
            namespace: item.namespace
          }
        },
        similarity: 1 - item._additional.distance, // Convert distance to similarity
        score: (1 - item._additional.distance) * this.calculateRecencyScore(item.timestamp)
      }));
    } catch (error) {
      console.error('Weaviate search error:', error);
      throw error;
    }
  }

  buildWhereFilter(filter) {
    const conditions = [];
    
    if (filter.namespace) {
      conditions.push({
        path: ['namespace'],
        operator: 'Equal',
        valueString: filter.namespace
      });
    }

    if (filter.type) {
      conditions.push({
        path: ['type'],
        operator: 'Equal',
        valueString: filter.type
      });
    }

    if (conditions.length === 0) {
      return undefined;
    }

    return {
      operator: 'And',
      operands: conditions
    };
  }

  async updateMemory(memoryId, updates) {
    try {
      await this.weaviate.data
        .merger()
        .withId(memoryId)
        .withClassName(this.className)
        .withProperties({
          ...updates,
          lastUpdated: new Date().toISOString()
        })
        .do();

      return true;
    } catch (error) {
      console.error('Weaviate update error:', error);
      throw error;
    }
  }

  async deleteMemory(memoryId) {
    try {
      await this.weaviate.data
        .deleter()
        .withId(memoryId)
        .withClassName(this.className)
        .do();
      return true;
    } catch (error) {
      console.error('Weaviate delete error:', error);
      throw error;
    }
  }
}

// Weaviate client initialization
class WeaviateManager {
  constructor(config) {
    this.config = config;
    this.client = weaviate.client(config);
  }

  async checkHealth() {
    try {
      const result = await this.client.misc.readyChecker().do();
      return result === true;
    } catch (error) {
      console.error('Weaviate health check failed:', error);
      return false;
    }
  }
}
```

### 4. Qdrant Integration

```javascript
class QdrantMemory extends VectorMemory {
  constructor(qdrantClient, collectionName, embeddingFunction, options = {}) {
    super(embeddingFunction, options);
    this.qdrant = qdrantClient;
    this.collectionName = collectionName;
    this.initialized = false;
  }

  async initializeCollection(dimension = 1536) {
    if (this.initialized) return;

    try {
      // Check if collection exists
      const collections = await this.qdrant.getCollections();
      const exists = collections.collections.some(c => c.name === this.collectionName);

      if (!exists) {
        await this.qdrant.createCollection(this.collectionName, {
          vectors: {
            size: dimension,
            distance: 'Cosine'
          }
        });
      }

      this.initialized = true;
    } catch (error) {
      console.error('Qdrant initialization error:', error);
      throw error;
    }
  }

  async addMemory(content, metadata = {}) {
    await this.initializeCollection();

    const embedding = await this.embeddingFunction(content);
    const memoryId = this.generateId();

    const payload = {
      content,
      timestamp: new Date().toISOString(),
      type: metadata.type || 'memory',
      namespace: this.options.namespace,
      ...metadata
    };

    try {
      await this.qdrant.upsert(this.collectionName, {
        wait: true,
        points: [
          {
            id: memoryId,
            vector: embedding,
            payload: payload
          }
        ]
      });

      return {
        id: memoryId,
        content,
        embedding,
        metadata: payload
      };
    } catch (error) {
      console.error('Qdrant error:', error);
      throw error;
    }
  }

  async searchMemories(query, options = {}) {
    await this.initializeCollection();
    const queryEmbedding = await this.embeddingFunction(query);

    const filter = this.buildFilter(options.filter);

    try {
      const result = await this.qdrant.search(this.collectionName, {
        vector: queryEmbedding,
        limit: options.maxResults || this.options.maxResults,
        score_threshold: options.similarityThreshold || this.options.similarityThreshold,
        filter: filter,
        with_payload: true,
        with_vector: false
      });

      return result.map(item => ({
        memory: {
          id: item.id,
          content: item.payload.content,
          metadata: item.payload
        },
        similarity: item.score,
        score: item.score * this.calculateRecencyScore(item.payload.timestamp)
      }));
    } catch (error) {
      console.error('Qdrant search error:', error);
      throw error;
    }
  }

  buildFilter(filterOptions = {}) {
    const must = [];

    if (filterOptions.namespace) {
      must.push({
        key: 'namespace',
        match: { value: filterOptions.namespace }
      });
    }

    if (filterOptions.type) {
      must.push({
        key: 'type',
        match: { value: filterOptions.type }
      });
    }

    return must.length > 0 ? { must } : undefined;
  }

  async updateMemory(memoryId, updates) {
    try {
      // Qdrant doesn't support partial payload updates directly
      // We need to fetch, update, and re-insert
      const points = await this.qdrant.retrieve(this.collectionName, {
        ids: [memoryId],
        with_payload: true
      });

      if (points.length === 0) {
        throw new Error('Memory not found');
      }

      const existing = points[0];
      const updatedPayload = {
        ...existing.payload,
        ...updates,
        lastUpdated: new Date().toISOString()
      };

      await this.qdrant.upsert(this.collectionName, {
        points: [
          {
            id: memoryId,
            vector: existing.vector,
            payload: updatedPayload
          }
        ]
      });

      return updatedPayload;

    } catch (error) {
      console.error('Qdrant update error:', error);
      throw error;
    }
  }

  async deleteMemory(memoryId) {
    try {
      await this.qdrant.delete(this.collectionName, {
        points: [memoryId]
      });
      return true;
    } catch (error) {
      console.error('Qdrant delete error:', error);
      throw error;
    }
  }
}

// Qdrant client setup
class QdrantManager {
  constructor(url, apiKey = null) {
    const config = {
      url: url
    };

    if (apiKey) {
      config.apiKey = apiKey;
    }

    this.client = new QdrantClient(config);
  }

  async healthCheck() {
    try {
      const result = await this.client.healthCheck();
      return result.status === 'ok';
    } catch (error) {
      console.error('Qdrant health check failed:', error);
      return false;
    }
  }
}
```

### 5. Advanced LTM System with Multiple Storage Tiers

```javascript
class AdvancedLTMSystem {
  constructor(config) {
    this.config = config;
    this.vectorDB = this.initializeVectorDB();
    this.cache = new Map(); // In-memory cache for recent memories
    this.embeddingFunction = config.embeddingFunction;
  }

  initializeVectorDB() {
    const { type, config } = this.config.vectorDB;
    
    switch (type) {
      case 'pinecone':
        return new PineconeMemory(
          config.client,
          config.indexName,
          this.embeddingFunction,
          config.options
        );
      case 'weaviate':
        return new WeaviateMemory(
          config.client,
          config.className,
          this.embeddingFunction,
          config.options
        );
      case 'qdrant':
        return new QdrantMemory(
          config.client,
          config.collectionName,
          this.embeddingFunction,
          config.options
        );
      default:
        throw new Error(`Unsupported vector DB type: ${type}`);
    }
  }

  async addConversationMemory(conversation, metadata = {}) {
    const memoryContent = this.formatConversationForMemory(conversation);
    const memory = await this.vectorDB.addMemory(memoryContent, {
      type: 'conversation',
      participants: metadata.participants || ['user', 'assistant'],
      topic: await this.extractTopic(memoryContent),
      ...metadata
    });

    // Cache recent memory
    this.cache.set(memory.id, {
      ...memory,
      lastAccessed: Date.now()
    });

    this.cleanCache();
    return memory;
  }

  async addFactMemory(fact, source, metadata = {}) {
    const memory = await this.vectorDB.addMemory(fact, {
      type: 'fact',
      source,
      verified: metadata.verified || false,
      category: metadata.category || 'general',
      ...metadata
    });

    return memory;
  }

  async searchRelevantMemories(query, context, options = {}) {
    // Combine query with context for better relevance
    const enhancedQuery = `${query} ${context}`;
    
    // Search in cache first
    const cachedResults = this.searchCache(enhancedQuery, options);
    
    // Search in vector database
    const vectorResults = await this.vectorDB.searchMemories(enhancedQuery, options);

    // Combine and deduplicate results
    const allResults = this.mergeResults(cachedResults, vectorResults);
    
    // Update cache access times
    allResults.forEach(result => {
      if (this.cache.has(result.memory.id)) {
        this.cache.get(result.memory.id).lastAccessed = Date.now();
      }
    });

    return allResults;
  }

  searchCache(query, options) {
    const results = [];
    const queryEmbedding = this.embeddingFunction(query); // Sync for cache search

    for (const [id, memory] of this.cache) {
      const similarity = this.cosineSimilarity(queryEmbedding, memory.embedding);
      
      if (similarity >= (options.similarityThreshold || 0.7)) {
        results.push({
          memory,
          similarity,
          score: similarity * this.calculateRecencyScore(memory.metadata.timestamp),
          fromCache: true
        });
      }
    }

    return results.sort((a, b) => b.score - a.score);
  }

  mergeResults(cachedResults, vectorResults) {
    const merged = [...cachedResults];
    const seenIds = new Set(cachedResults.map(r => r.memory.id));

    for (const vectorResult of vectorResults) {
      if (!seenIds.has(vectorResult.memory.id)) {
        merged.push(vectorResult);
        seenIds.add(vectorResult.memory.id);
      }
    }

    return merged.sort((a, b) => b.score - a.score);
  }

  formatConversationForMemory(conversation) {
    return conversation.map(entry => 
      `${entry.role}: ${entry.content}`
    ).join('\n');
  }

  async extractTopic(text) {
    // Simple topic extraction - in practice, use NLP
    const sentences = text.split(/[.!?]+/);
    const firstSentence = sentences[0].toLowerCase();
    
    const topics = {
      technology: ['code', 'program', 'software', 'tech', 'computer'],
      science: ['research', 'study', 'experiment', 'scientific'],
      business: ['company', 'business', 'market', 'sales']
    };

    for (const [topic, keywords] of Object.entries(topics)) {
      if (keywords.some(keyword => firstSentence.includes(keyword))) {
        return topic;
      }
    }

    return 'general';
  }

  cleanCache() {
    const maxCacheSize = this.config.cacheSize || 100;
    if (this.cache.size > maxCacheSize) {
      const entries = Array.from(this.cache.entries());
      entries.sort((a, b) => a[1].lastAccessed - b[1].lastAccessed);
      
      // Remove oldest 10% of cache
      const toRemove = entries.slice(0, Math.floor(maxCacheSize * 0.1));
      toRemove.forEach(([id]) => this.cache.delete(id));
    }
  }

  async getMemoryStatistics() {
    const stats = await this.vectorDB.getMemoryStats();
    
    return {
      ...stats,
      cacheSize: this.cache.size,
      cacheHitRate: this.calculateCacheHitRate(),
      averageSimilarityScore: await this.calculateAverageSimilarity()
    };
  }

  calculateCacheHitRate() {
    // Implementation would track cache hits/misses over time
    return 0.7; // Example value
  }

  async calculateAverageSimilarity() {
    // Sample some memories and calculate average similarity
    return 0.8; // Example value
  }
}
```

### 6. Practical Implementation with OpenAI Embeddings

```javascript
// Example using OpenAI embeddings with Pinecone
class OpenAILTM {
  constructor(openaiApiKey, pineconeConfig) {
    this.openai = new OpenAI({ apiKey: openaiApiKey });
    this.pineconeManager = new PineconeManager(
      pineconeConfig.apiKey,
      pineconeConfig.environment
    );
    this.memorySystem = null;
  }

  async initialize() {
    const index = await this.pineconeManager.getIndex('memories', 1536);
    this.memorySystem = new PineconeMemory(
      index,
      this.generateEmbedding.bind(this),
      {
        similarityThreshold: 0.75,
        maxResults: 10,
        namespace: 'conversations'
      }
    );

    console.log('LTM system initialized');
  }

  async generateEmbedding(text) {
    try {
      const response = await this.openai.embeddings.create({
        model: 'text-embedding-ada-002',
        input: text
      });

      return response.data[0].embedding;
    } catch (error) {
      console.error('OpenAI embedding error:', error);
      throw error;
    }
  }

  async rememberConversation(messages, summary, metadata = {}) {
    if (!this.memorySystem) {
      throw new Error('LTM system not initialized');
    }

    const memoryContent = this.createMemoryContent(messages, summary);
    
    return await this.memorySystem.addMemory(memoryContent, {
      type: 'conversation',
      summary,
      messageCount: messages.length,
      ...metadata
    });
  }

  async recallRelevantContext(query, recentContext = '') {
    if (!this.memorySystem) {
      throw new Error('LTM system not initialized');
    }

    const results = await this.memorySystem.searchMemories(query, {
      filter: { type: 'conversation' },
      maxResults: 5
    });

    return this.formatRecalledContext(results, recentContext);
  }

  createMemoryContent(messages, summary) {
    const recentMessages = messages.slice(-3); // Last 3 messages
    const conversationExcerpt = recentMessages.map(m => 
      `${m.role}: ${m.content}`
    ).join('\n');

    return `Summary: ${summary}\nRecent exchange:\n${conversationExcerpt}`;
  }

  formatRecalledContext(searchResults, recentContext) {
    if (searchResults.length === 0) {
      return recentContext;
    }

    const relevantMemories = searchResults
      .filter(result => result.score > 0.6)
      .map(result => result.memory.content)
      .join('\n\n');

    return `Previous conversations:\n${relevantMemories}\n\nCurrent context:\n${recentContext}`;
  }

  async manageMemoryLifecycle() {
    // Implement memory pruning, compression, etc.
    const stats = await this.memorySystem.getMemoryStats();
    
    if (stats.totalMemories > 1000) {
      console.log('Memory archive needed - implementing compression strategy');
      // Archive old memories, create summaries, etc.
    }
  }
}

// Usage example
async function demonstrateLTM() {
  const ltm = new OpenAILTM(
    'your-openai-key',
    {
      apiKey: 'your-pinecone-key',
      environment: 'us-west1-gcp'
    }
  );

  await ltm.initialize();

  // Remember a conversation
  const conversation = [
    { role: 'user', content: 'How do I learn machine learning?' },
    { role: 'assistant', content: 'Start with Python and basic statistics...' }
  ];

  await ltm.rememberConversation(
    conversation,
    'Discussion about learning machine learning',
    { topic: 'education', difficulty: 'beginner' }
  );

  // Recall relevant information
  const context = await ltm.recallRelevantContext(
    'machine learning resources for beginners',
    'User is asking about learning path'
  );

  console.log('Recalled context:', context);
}
```

## Key Benefits of Vector Database LTM

1. **Semantic Understanding**: Finds relevant content based on meaning
2. **Scalability**: Handles millions of memories efficiently
3. **Persistence**: Maintains knowledge across sessions
4. **Flexible Querying**: Supports complex semantic searches
5. **Real-time Performance**: Fast similarity searches

## Best Practices

1. **Choose appropriate embedding models** for your use case
2. **Implement proper namespacing** for different memory types
3. **Set appropriate similarity thresholds** based on your domain
4. **Combine with traditional databases** for metadata management
5. **Implement cache layers** for frequently accessed memories
6. **Regularly maintain and prune** memory collections

Vector databases provide a powerful foundation for building sophisticated long-term memory systems that can understand and recall information based on semantic meaning rather than just keyword matching.