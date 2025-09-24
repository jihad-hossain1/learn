# Hybrid Search (Vector + Keyword)

Hybrid search combines the semantic understanding of vector search with the precision of keyword search to deliver superior results. It leverages the strengths of both approaches while mitigating their weaknesses.

## How Hybrid Search Works

- **Vector search**: Finds semantically similar content
- **Keyword search**: Finds exact term matches
- **Fusion**: Combines results using scoring algorithms
- **Re-ranking**: Optimizes final result order

## Core Concepts

### Why Hybrid Search?
- **Vector pros**: Semantic understanding, synonym handling
- **Vector cons**: May miss exact matches, computationally expensive
- **Keyword pros**: Fast, precise for known terms
- **Keyword cons**: No semantic understanding, vocabulary mismatch

### Fusion Strategies
- **Reciprocal Rank Fusion (RRF)**: Combines rankings from both methods
- **Weighted scoring**: Assigns weights to each method's scores
- **Re-ranking**: Uses cross-encoders for final ordering

## JavaScript Implementation Examples

### 1. Base Hybrid Search Class

```javascript
class HybridSearch {
  constructor(vectorSearch, keywordSearch, options = {}) {
    this.vectorSearch = vectorSearch;
    this.keywordSearch = keywordSearch;
    this.options = {
      fusionMethod: options.fusionMethod || 'rrf', // 'rrf', 'weighted', 'concatenation'
      vectorWeight: options.vectorWeight || 0.7,
      keywordWeight: options.keywordWeight || 0.3,
      maxResults: options.maxResults || 20,
      rrfK: options.rrfK || 60,
      ...options
    };
  }

  async search(query, options = {}) {
    const searchOptions = { ...this.options, ...options };
    
    // Execute both searches in parallel
    const [vectorResults, keywordResults] = await Promise.all([
      this.vectorSearch.search(query, searchOptions),
      this.keywordSearch.search(query, searchOptions)
    ]);

    // Fuse results based on selected method
    switch (searchOptions.fusionMethod) {
      case 'rrf':
        return this.fuseWithRRF(vectorResults, keywordResults, searchOptions);
      case 'weighted':
        return this.fuseWithWeightedScores(vectorResults, keywordResults, searchOptions);
      case 'concatenation':
        return this.fuseWithConcatenation(vectorResults, keywordResults, searchOptions);
      default:
        throw new Error(`Unknown fusion method: ${searchOptions.fusionMethod}`);
    }
  }

  fuseWithRRF(vectorResults, keywordResults, options) {
    const fusedScores = new Map();
    
    // Score results from vector search
    vectorResults.forEach((result, index) => {
      const rrfScore = 1 / (options.rrfK + index + 1);
      fusedScores.set(result.id, {
        ...result,
        vectorScore: result.score || (1 - index / vectorResults.length),
        rrfVectorScore: rrfScore,
        rrfKeywordScore: 0,
        finalScore: rrfScore
      });
    });

    // Score results from keyword search
    keywordResults.forEach((result, index) => {
      const rrfScore = 1 / (options.rrfK + index + 1);
      const existing = fusedScores.get(result.id);
      
      if (existing) {
        existing.rrfKeywordScore = rrfScore;
        existing.finalScore = existing.rrfVectorScore + rrfScore;
      } else {
        fusedScores.set(result.id, {
          ...result,
          vectorScore: 0,
          rrfVectorScore: 0,
          rrfKeywordScore: rrfScore,
          finalScore: rrfScore
        });
      }
    });

    // Convert to array and sort by final score
    return Array.from(fusedScores.values())
      .sort((a, b) => b.finalScore - a.finalScore)
      .slice(0, options.maxResults);
  }

  fuseWithWeightedScores(vectorResults, keywordResults, options) {
    const fusedScores = new Map();
    
    // Normalize vector scores
    const maxVectorScore = Math.max(...vectorResults.map(r => r.score || 0), 1);
    vectorResults.forEach(result => {
      const normalizedVectorScore = (result.score || 0) / maxVectorScore;
      fusedScores.set(result.id, {
        ...result,
        vectorScore: normalizedVectorScore,
        keywordScore: 0,
        finalScore: normalizedVectorScore * options.vectorWeight
      });
    });

    // Normalize keyword scores and combine
    const maxKeywordScore = Math.max(...keywordResults.map(r => r.score || 0), 1);
    keywordResults.forEach(result => {
      const normalizedKeywordScore = (result.score || 0) / maxKeywordScore;
      const existing = fusedScores.get(result.id);
      
      if (existing) {
        existing.keywordScore = normalizedKeywordScore;
        existing.finalScore = (existing.vectorScore * options.vectorWeight) + 
                            (normalizedKeywordScore * options.keywordWeight);
      } else {
        fusedScores.set(result.id, {
          ...result,
          vectorScore: 0,
          keywordScore: normalizedKeywordScore,
          finalScore: normalizedKeywordScore * options.keywordWeight
        });
      }
    });

    return Array.from(fusedScores.values())
      .sort((a, b) => b.finalScore - a.finalScore)
      .slice(0, options.maxResults);
  }

  fuseWithConcatenation(vectorResults, keywordResults, options) {
    // Simple concatenation - remove duplicates and combine
    const allResults = new Map();
    
    vectorResults.forEach(result => {
      allResults.set(result.id, {
        ...result,
        source: 'vector',
        rank: allResults.size + 1
      });
    });

    keywordResults.forEach(result => {
      if (!allResults.has(result.id)) {
        allResults.set(result.id, {
          ...result,
          source: 'keyword',
          rank: allResults.size + 1
        });
      }
    });

    return Array.from(allResults.values())
      .slice(0, options.maxResults);
  }

  async searchWithReRanking(query, options = {}) {
    // First pass: hybrid search
    const initialResults = await this.search(query, {
      ...options,
      maxResults: options.rerankTopK || 50
    });

    // Second pass: re-ranking with cross-encoder or more sophisticated model
    return this.reRankResults(query, initialResults, options);
  }

  async reRankResults(query, results, options) {
    if (!options.reRanker || results.length <= 1) {
      return results;
    }

    // Simple re-ranking based on combined features
    const reRanked = results.map(result => ({
      ...result,
      reRankScore: this.calculateReRankScore(query, result, options)
    }));

    return reRanked
      .sort((a, b) => b.reRankScore - a.reRankScore)
      .slice(0, options.maxResults || this.options.maxResults);
  }

  calculateReRankScore(query, result, options) {
    let score = result.finalScore || 0;
    
    // Boost exact matches
    if (this.hasExactMatch(query, result.content || result.text)) {
      score *= 1.5;
    }
    
    // Boost recent content
    if (result.timestamp) {
      const recency = this.calculateRecencyBoost(result.timestamp);
      score *= recency;
    }
    
    // Boost based on popularity/quality signals
    if (result.popularity) {
      score *= (1 + result.popularity * 0.1);
    }

    return score;
  }

  hasExactMatch(query, text) {
    const queryTerms = query.toLowerCase().split(/\s+/);
    const textLower = text.toLowerCase();
    
    return queryTerms.some(term => 
      textLower.includes(term) && term.length > 2
    );
  }

  calculateRecencyBoost(timestamp) {
    const daysOld = (Date.now() - new Date(timestamp).getTime()) / (1000 * 60 * 60 * 24);
    return Math.max(0.5, 1 - (daysOld / 365)); // Linear decay over 1 year
  }
}
```

### 2. Vector Search Implementation

```javascript
class VectorSearch {
  constructor(embeddingFunction, vectorStore, options = {}) {
    this.embeddingFunction = embeddingFunction;
    this.vectorStore = vectorStore;
    this.options = {
      similarityThreshold: options.similarityThreshold || 0.7,
      maxResults: options.maxResults || 10,
      ...options
    };
  }

  async search(query, options = {}) {
    const searchOptions = { ...this.options, ...options };
    
    try {
      const queryEmbedding = await this.embeddingFunction.generateEmbedding(query);
      const results = await this.vectorStore.similaritySearch(
        queryEmbedding,
        searchOptions
      );

      return results.map((result, index) => ({
        id: result.id || `vec_${index}`,
        content: result.content || result.text,
        score: result.similarity || result.score,
        metadata: result.metadata || {},
        type: 'vector',
        timestamp: result.timestamp || new Date()
      }));
    } catch (error) {
      console.error('Vector search error:', error);
      return [];
    }
  }

  async searchWithFilters(query, filters, options = {}) {
    const queryEmbedding = await this.embeddingFunction.generateEmbedding(query);
    
    return this.vectorStore.similaritySearchWithFilters(
      queryEmbedding,
      filters,
      { ...this.options, ...options }
    );
  }
}

// Example vector store implementation
class SimpleVectorStore {
  constructor() {
    this.vectors = new Map();
    this.metadata = new Map();
  }

  addVector(id, vector, metadata = {}) {
    this.vectors.set(id, vector);
    this.metadata.set(id, {
      ...metadata,
      id,
      timestamp: metadata.timestamp || new Date()
    });
  }

  similaritySearch(queryVector, options = {}) {
    const results = [];
    
    for (const [id, vector] of this.vectors) {
      const similarity = this.cosineSimilarity(queryVector, vector);
      
      if (similarity >= (options.similarityThreshold || 0.7)) {
        const metadata = this.metadata.get(id);
        results.push({
          id,
          similarity,
          score: similarity,
          ...metadata
        });
      }
    }

    return results
      .sort((a, b) => b.similarity - a.similarity)
      .slice(0, options.maxResults || 10);
  }

  cosineSimilarity(vecA, vecB) {
    if (vecA.length !== vecB.length) {
      throw new Error('Vectors must have same dimensions');
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
}
```

### 3. Keyword Search Implementation

```javascript
class KeywordSearch {
  constructor(index, options = {}) {
    this.index = index;
    this.options = {
      boostExactMatches: options.boostExactMatches !== false,
      useStemming: options.useStemming || false,
      fuzzyMatching: options.fuzzyMatching || false,
      maxResults: options.maxResults || 10,
      ...options
    };
    
    if (this.options.useStemming) {
      this.stemmer = new PorterStemmer();
    }
  }

  async search(query, options = {}) {
    const searchOptions = { ...this.options, ...options };
    const queryTerms = this.processQuery(query, searchOptions);
    
    const results = this.executeKeywordSearch(queryTerms, searchOptions);
    return this.rankResults(results, queryTerms, searchOptions);
  }

  processQuery(query, options) {
    // Tokenize and clean query
    let terms = query.toLowerCase()
      .split(/\s+/)
      .filter(term => term.length > 1) // Remove single characters
      .map(term => term.replace(/[^\w]/g, '')); // Remove punctuation

    // Apply stemming if enabled
    if (options.useStemming && this.stemmer) {
      terms = terms.map(term => this.stemmer.stem(term));
    }

    return [...new Set(terms)]; // Remove duplicates
  }

  executeKeywordSearch(queryTerms, options) {
    const results = new Map();
    
    for (const term of queryTerms) {
      const documents = this.index.get(term) || [];
      
      for (const doc of documents) {
        const existing = results.get(doc.id);
        const termFrequency = this.calculateTermFrequency(term, doc.content);
        
        if (existing) {
          existing.score += termFrequency;
          existing.matchedTerms.add(term);
        } else {
          results.set(doc.id, {
            id: doc.id,
            content: doc.content,
            score: termFrequency,
            matchedTerms: new Set([term]),
            metadata: doc.metadata || {},
            type: 'keyword',
            timestamp: doc.timestamp
          });
        }
      }
    }

    return Array.from(results.values());
  }

  calculateTermFrequency(term, content) {
    const contentLower = content.toLowerCase();
    const termRegex = new RegExp(`\\b${term}\\b`, 'gi');
    const matches = contentLower.match(termRegex);
    
    if (!matches) return 0;
    
    // TF with logarithmic scaling to prevent dominance by single terms
    return 1 + Math.log(matches.length);
  }

  rankResults(results, queryTerms, options) {
    return results
      .map(result => {
        let score = result.score;
        
        // Boost documents that match more query terms
        const termCoverage = result.matchedTerms.size / queryTerms.length;
        score *= (1 + termCoverage);
        
        // Boost exact phrase matches
        if (options.boostExactMatches) {
          const exactMatchBoost = this.calculateExactMatchBoost(queryTerms, result.content);
          score *= exactMatchBoost;
        }
        
        // Apply document length normalization
        score = this.applyLengthNormalization(score, result.content);
        
        return {
          ...result,
          score,
          matchedTerms: Array.from(result.matchedTerms)
        };
      })
      .sort((a, b) => b.score - a.score)
      .slice(0, options.maxResults);
  }

  calculateExactMatchBoost(queryTerms, content) {
    const contentLower = content.toLowerCase();
    const queryPhrase = queryTerms.join(' ');
    
    if (contentLower.includes(queryPhrase)) {
      return 2.0; // Significant boost for exact phrase matches
    }
    
    return 1.0;
  }

  applyLengthNormalization(score, content) {
    // Penalize very long documents (BM25-like normalization)
    const avgDocumentLength = 1000; // Average document length in words
    const k1 = 1.2; // BM25 parameter
    const b = 0.75; // BM25 parameter
    
    const documentLength = content.split(/\s+/).length;
    const lengthPenalty = (1 - b) + b * (documentLength / avgDocumentLength);
    
    return score / (k1 * lengthPenalty + score);
  }
}

// Simple inverted index implementation
class InvertedIndex {
  constructor() {
    this.index = new Map();
    this.documents = new Map();
  }

  addDocument(id, content, metadata = {}) {
    this.documents.set(id, { id, content, metadata });
    
    const terms = this.tokenizeContent(content);
    for (const term of terms) {
      if (!this.index.has(term)) {
        this.index.set(term, new Set());
      }
      this.index.get(term).add(id);
    }
  }

  tokenizeContent(content) {
    return content.toLowerCase()
      .split(/\s+/)
      .filter(term => term.length > 1)
      .map(term => term.replace(/[^\w]/g, ''))
      .filter(term => term.length > 0);
  }

  get(term) {
    const documentIds = this.index.get(term);
    if (!documentIds) return [];
    
    return Array.from(documentIds).map(id => this.documents.get(id));
  }
}

// Simple stemmer implementation
class PorterStemmer {
  stem(word) {
    // Simplified Porter stemming algorithm
    // In practice, use a complete implementation
    if (word.endsWith('ing')) {
      return word.substring(0, word.length - 3);
    }
    if (word.endsWith('ed')) {
      return word.substring(0, word.length - 2);
    }
    if (word.endsWith('s')) {
      return word.substring(0, word.length - 1);
    }
    return word;
  }
}
```

### 4. Advanced Hybrid Search with Query Understanding

```javascript
class IntelligentHybridSearch extends HybridSearch {
  constructor(vectorSearch, keywordSearch, options = {}) {
    super(vectorSearch, keywordSearch, options);
    this.queryAnalyzer = new QueryAnalyzer();
  }

  async search(query, options = {}) {
    // Analyze query to determine optimal strategy
    const queryAnalysis = this.queryAnalyzer.analyze(query);
    const optimizedOptions = this.optimizeSearchOptions(queryAnalysis, options);
    
    // Execute search with optimized parameters
    const results = await super.search(query, optimizedOptions);
    
    // Apply query-specific boosting
    return this.applyQuerySpecificBoosting(results, queryAnalysis);
  }

  optimizeSearchOptions(queryAnalysis, baseOptions) {
    const optimized = { ...baseOptions };
    
    switch (queryAnalysis.type) {
      case 'factual':
        optimized.keywordWeight = Math.max(0.6, optimized.keywordWeight);
        optimized.vectorWeight = 1 - optimized.keywordWeight;
        break;
        
      case 'semantic':
        optimized.vectorWeight = Math.max(0.8, optimized.vectorWeight);
        optimized.keywordWeight = 1 - optimized.vectorWeight;
        break;
        
      case 'exploratory':
        optimized.fusionMethod = 'rrf';
        optimized.maxResults = Math.max(20, optimized.maxResults);
        break;
        
      case 'navigational':
        optimized.keywordWeight = 0.9;
        optimized.vectorWeight = 0.1;
        break;
    }
    
    // Adjust based on query complexity
    if (queryAnalysis.complexity === 'high') {
      optimized.maxResults = Math.max(optimized.maxResults, 15);
    }
    
    return optimized;
  }

  applyQuerySpecificBoosting(results, queryAnalysis) {
    return results.map(result => {
      let boost = 1.0;
      
      // Boost based on query type
      if (queryAnalysis.type === 'factual' && this.isFactualContent(result)) {
        boost *= 1.3;
      }
      
      if (queryAnalysis.type === 'semantic' && this.hasSemanticDepth(result)) {
        boost *= 1.2;
      }
      
      // Boost based on content type match
      if (queryAnalysis.contentType && result.metadata?.type === queryAnalysis.contentType) {
        boost *= 1.5;
      }
      
      return {
        ...result,
        finalScore: (result.finalScore || result.score) * boost,
        boostFactors: {
          queryType: boost,
          originalScore: result.finalScore || result.score
        }
      };
    }).sort((a, b) => b.finalScore - a.finalScore);
  }

  isFactualContent(result) {
    const factualIndicators = ['fact', 'data', 'statistic', 'study', 'research'];
    const content = result.content.toLowerCase();
    return factualIndicators.some(indicator => content.includes(indicator));
  }

  hasSemanticDepth(result) {
    // Check if content has semantic richness
    const sentences = result.content.split(/[.!?]+/).length;
    const words = result.content.split(/\s+/).length;
    const uniqueWords = new Set(result.content.toLowerCase().split(/\W+/)).size;
    
    return uniqueWords / words > 0.6 && sentences > 3;
  }
}

class QueryAnalyzer {
  analyze(query) {
    return {
      type: this.classifyQueryType(query),
      complexity: this.assessComplexity(query),
      contentType: this.predictContentType(query),
      terms: this.extractKeyTerms(query),
      length: query.length
    };
  }

  classifyQueryType(query) {
    const lowerQuery = query.toLowerCase();
    
    if (this.isFactualQuery(lowerQuery)) return 'factual';
    if (this.isSemanticQuery(lowerQuery)) return 'semantic';
    if (this.isNavigationalQuery(lowerQuery)) return 'navigational';
    if (this.isExploratoryQuery(lowerQuery)) return 'exploratory';
    
    return 'general';
  }

  isFactualQuery(query) {
    const factualPatterns = [
      /^what is/,
      /^when did/,
      /^how many/,
      /^who invented/,
      /^where is/,
      /\?$/,
      /^facts about/
    ];
    
    return factualPatterns.some(pattern => pattern.test(query));
  }

  isSemanticQuery(query) {
    const semanticPatterns = [
      /explain/,
      /describe/,
      /discuss/,
      /compare/,
      /contrast/,
      /analyze/
    ];
    
    return semanticPatterns.some(pattern => pattern.test(query));
  }

  isNavigationalQuery(query) {
    const navigationalPatterns = [
      /^go to/,
      /^open/,
      /^show me/,
      /^take me to/,
      /^find the page for/
    ];
    
    return navigationalPatterns.some(pattern => pattern.test(query));
  }

  isExploratoryQuery(query) {
    const exploratoryPatterns = [
      /^ideas for/,
      /^options for/,
      /^suggest/,
      /^recommend/,
      /^what are some/
    ];
    
    return exploratoryPatterns.some(pattern => pattern.test(query));
  }

  assessComplexity(query) {
    const wordCount = query.split(/\s+/).length;
    const uniqueWords = new Set(query.toLowerCase().split(/\W+/)).size;
    const complexityRatio = uniqueWords / wordCount;
    
    if (wordCount > 10 || complexityRatio > 0.8) return 'high';
    if (wordCount > 5 || complexityRatio > 0.6) return 'medium';
    return 'low';
  }

  predictContentType(query) {
    const contentTypes = {
      technical: ['code', 'programming', 'algorithm', 'software', 'technical'],
      academic: ['research', 'study', 'paper', 'academic', 'scholarly'],
      commercial: ['buy', 'price', 'product', 'service', 'commercial'],
      general: ['information', 'about', 'what is', 'how to']
    };
    
    const lowerQuery = query.toLowerCase();
    for (const [type, keywords] of Object.entries(contentTypes)) {
      if (keywords.some(keyword => lowerQuery.includes(keyword))) {
        return type;
      }
    }
    
    return 'general';
  }

  extractKeyTerms(query) {
    // Remove stop words and extract meaningful terms
    const stopWords = new Set(['the', 'a', 'an', 'and', 'or', 'but', 'in', 'on', 'at', 'to', 'for']);
    
    return query.toLowerCase()
      .split(/\s+/)
      .filter(term => term.length > 2 && !stopWords.has(term))
      .map(term => term.replace(/[^\w]/g, ''));
  }
}
```

### 5. Hybrid Search with Database Integration

```javascript
class DatabaseHybridSearch {
  constructor(database, embeddingService, options = {}) {
    this.database = database;
    this.embeddingService = embeddingService;
    this.options = {
      vectorColumn: options.vectorColumn || 'embedding',
      contentColumn: options.contentColumn || 'content',
      maxResults: options.maxResults || 20,
      ...options
    };
  }

  async search(query, options = {}) {
    const searchOptions = { ...this.options, ...options };
    
    // Execute vector search using database capabilities
    const vectorResults = await this.vectorSearchInDB(query, searchOptions);
    
    // Execute keyword search using database full-text search
    const keywordResults = await this.keywordSearchInDB(query, searchOptions);
    
    // Fuse results
    return this.fuseResults(vectorResults, keywordResults, searchOptions);
  }

  async vectorSearchInDB(query, options) {
    try {
      const queryEmbedding = await this.embeddingService.generateEmbedding(query);
      
      // Use database vector search (e.g., PostgreSQL with pgvector)
      const results = await this.database.query(
        `SELECT id, content, metadata, timestamp,
         ${this.options.vectorColumn} <=> $1 as similarity
         FROM documents 
         WHERE ${this.options.vectorColumn} <=> $1 < $2
         ORDER BY similarity 
         LIMIT $3`,
        [queryEmbedding, 1 - options.similarityThreshold, options.maxResults]
      );

      return results.rows.map(row => ({
        id: row.id,
        content: row.content,
        score: 1 - row.similarity, // Convert distance to similarity
        metadata: row.metadata,
        timestamp: row.timestamp,
        type: 'vector'
      }));
    } catch (error) {
      console.error('Database vector search error:', error);
      return [];
    }
  }

  async keywordSearchInDB(query, options) {
    try {
      // Use database full-text search
      const results = await this.database.query(
        `SELECT id, content, metadata, timestamp,
         ts_rank(to_tsvector(content), to_tsquery($1)) as rank
         FROM documents 
         WHERE to_tsvector(content) @@ to_tsquery($1)
         ORDER BY rank DESC
         LIMIT $2`,
        [this.queryToTSQuery(query), options.maxResults]
      );

      return results.rows.map(row => ({
        id: row.id,
        content: row.content,
        score: row.rank,
        metadata: row.metadata,
        timestamp: row.timestamp,
        type: 'keyword'
      }));
    } catch (error) {
      console.error('Database keyword search error:', error);
      return [];
    }
  }

  queryToTSQuery(query) {
    // Convert natural language query to PostgreSQL TSQuery format
    return query.split(/\s+/)
      .filter(term => term.length > 2)
      .map(term => term + ':*')
      .join(' & ');
  }

  fuseResults(vectorResults, keywordResults, options) {
    const fused = new Map();
    
    // Combine results using RRF
    vectorResults.forEach((result, index) => {
      const rrfScore = 1 / (options.rrfK + index + 1);
      fused.set(result.id, {
        ...result,
        rrfScore,
        finalScore: rrfScore
      });
    });

    keywordResults.forEach((result, index) => {
      const rrfScore = 1 / (options.rrfK + index + 1);
      const existing = fused.get(result.id);
      
      if (existing) {
        existing.finalScore += rrfScore;
      } else {
        fused.set(result.id, {
          ...result,
          rrfScore,
          finalScore: rrfScore
        });
      }
    });

    return Array.from(fused.values())
      .sort((a, b) => b.finalScore - a.finalScore)
      .slice(0, options.maxResults);
  }

  async addDocument(content, metadata = {}) {
    // Generate embedding and store in database
    const embedding = await this.embeddingService.generateEmbedding(content);
    
    await this.database.query(
      `INSERT INTO documents (content, metadata, embedding, timestamp)
       VALUES ($1, $2, $3, $4)`,
      [content, metadata, embedding, new Date()]
    );
  }

  async createSearchIndexes() {
    // Create vector index
    await this.database.query(`
      CREATE INDEX IF NOT EXISTS documents_embedding_idx 
      ON documents USING ivfflat (embedding vector_cosine_ops)
    `);
    
    // Create full-text search index
    await this.database.query(`
      CREATE INDEX IF NOT EXISTS documents_content_tsvector_idx 
      ON documents USING gin(to_tsvector('english', content))
    `);
  }
}
```

### 6. Performance-Optimized Hybrid Search

```javascript
class PerformanceOptimizedHybridSearch {
  constructor(components, options = {}) {
    this.components = components;
    this.options = {
      cacheEnabled: options.cacheEnabled !== false,
      precomputeEmbeddings: options.precomputeEmbeddings || false,
      batchSize: options.batchSize || 10,
      timeout: options.timeout || 5000,
      ...options
    };
    
    this.cache = new Map();
    this.metrics = new SearchMetrics();
  }

  async search(query, options = {}) {
    const startTime = Date.now();
    const searchId = this.generateSearchId(query, options);
    
    // Check cache
    if (this.options.cacheEnabled) {
      const cached = this.cache.get(searchId);
      if (cached && Date.now() - cached.timestamp < (options.cacheTTL || 300000)) {
        this.metrics.recordCacheHit();
        return cached.results;
      }
    }

    try {
      // Execute search with timeout
      const results = await this.executeSearchWithTimeout(query, options, startTime);
      
      // Cache results
      if (this.options.cacheEnabled) {
        this.cache.set(searchId, {
          results,
          timestamp: Date.now()
        });
        this.cleanCache();
      }
      
      this.metrics.recordSuccess(Date.now() - startTime, results.length);
      return results;
    } catch (error) {
      this.metrics.recordError(error);
      throw error;
    }
  }

  async executeSearchWithTimeout(query, options, startTime) {
    return new Promise(async (resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error(`Search timeout after ${this.options.timeout}ms`));
      }, this.options.timeout);
      
      try {
        const results = await this.executeSearch(query, options);
        clearTimeout(timeout);
        resolve(results);
      } catch (error) {
        clearTimeout(timeout);
        reject(error);
      }
    });
  }

  async executeSearch(query, options) {
    const [vectorResults, keywordResults] = await Promise.all([
      this.optimizedVectorSearch(query, options),
      this.optimizedKeywordSearch(query, options)
    ]);

    // Early return if one method returns no results
    if (vectorResults.length === 0 && keywordResults.length > 0) {
      return keywordResults.slice(0, options.maxResults || this.options.maxResults);
    }
    if (keywordResults.length === 0 && vectorResults.length > 0) {
      return vectorResults.slice(0, options.maxResults || this.options.maxResults);
    }

    return this.fuseResults(vectorResults, keywordResults, options);
  }

  async optimizedVectorSearch(query, options) {
    // Use precomputed embeddings if available
    if (this.options.precomputeEmbeddings) {
      return this.components.vectorSearch.searchWithPrecomputed(query, options);
    }
    
    // Batch embedding generation for multiple terms
    const queryTerms = this.extractSearchTerms(query);
    if (queryTerms.length > 1) {
      return this.batchVectorSearch(queryTerms, options);
    }
    
    return this.components.vectorSearch.search(query, options);
  }

  async batchVectorSearch(terms, options) {
    const batchPromises = terms.map(term => 
      this.components.vectorSearch.search(term, { ...options, maxResults: 5 })
    );
    
    const batchResults = await Promise.all(batchPromises);
    return this.mergeBatchResults(batchResults);
  }

  mergeBatchResults(batchResults) {
    const merged = new Map();
    
    batchResults.forEach((results, index) => {
      results.forEach(result => {
        const existing = merged.get(result.id);
        if (existing) {
          existing.score += result.score;
          existing.termMatches++;
        } else {
          merged.set(result.id, {
            ...result,
            termMatches: 1
          });
        }
      });
    });
    
    return Array.from(merged.values())
      .sort((a, b) => b.score - a.score)
      .map(result => ({
        ...result,
        score: result.score * (1 + result.termMatches * 0.1) // Boost multi-term matches
      }));
  }

  async optimizedKeywordSearch(query, options) {
    // Query expansion for better keyword matching
    const expandedQuery = await this.expandQuery(query);
    return this.components.keywordSearch.search(expandedQuery, options);
  }

  async expandQuery(query) {
    // Simple query expansion using synonyms
    const synonyms = await this.findSynonyms(query);
    return `${query} ${synonyms.join(' ')}`;
  }

  async findSynonyms(query) {
    // Use a synonym service or predefined dictionary
    const synonymMap = {
      'big': ['large', 'huge', 'enormous'],
      'fast': ['quick', 'rapid', 'speedy'],
      'learn': ['study', 'understand', 'master']
    };
    
    const terms = query.toLowerCase().split(/\s+/);
    const synonyms = [];
    
    for (const term of terms) {
      if (synonymMap[term]) {
        synonyms.push(...synonymMap[term].slice(0, 2));
      }
    }
    
    return synonyms;
  }

  extractSearchTerms(query) {
    return query.split(/\s+/)
      .filter(term => term.length > 2)
      .slice(0, 5); // Limit to 5 terms for batching
  }

  generateSearchId(query, options) {
    const key = `${query}-${JSON.stringify(options)}`;
    return Buffer.from(key).toString('base64').substring(0, 50);
  }

  cleanCache() {
    const maxCacheSize = this.options.maxCacheSize || 1000;
    if (this.cache.size > maxCacheSize) {
      const entries = Array.from(this.cache.entries());
      entries.sort((a, b) => a[1].timestamp - b[1].timestamp);
      
      // Remove oldest 10%
      const toRemove = entries.slice(0, Math.floor(maxCacheSize * 0.1));
      toRemove.forEach(([key]) => this.cache.delete(key));
    }
  }

  getPerformanceMetrics() {
    return this.metrics.getSummary();
  }
}

class SearchMetrics {
  constructor() {
    this.metrics = {
      totalSearches: 0,
      successfulSearches: 0,
      failedSearches: 0,
      cacheHits: 0,
      averageResponseTime: 0,
      responseTimes: []
    };
  }

  recordSuccess(responseTime, resultCount) {
    this.metrics.totalSearches++;
    this.metrics.successfulSearches++;
    this.metrics.responseTimes.push(responseTime);
    this.metrics.averageResponseTime = this.metrics.responseTimes.reduce((a, b) => a + b, 0) / this.metrics.responseTimes.length;
  }

  recordError(error) {
    this.metrics.totalSearches++;
    this.metrics.failedSearches++;
  }

  recordCacheHit() {
    this.metrics.cacheHits++;
  }

  getSummary() {
    const successRate = this.metrics.totalSearches > 0 ? 
      (this.metrics.successfulSearches / this.metrics.totalSearches) * 100 : 0;
    
    const cacheHitRate = this.metrics.totalSearches > 0 ?
      (this.metrics.cacheHits / this.metrics.totalSearches) * 100 : 0;
    
    return {
      ...this.metrics,
      successRate: `${successRate.toFixed(1)}%`,
      cacheHitRate: `${cacheHitRate.toFixed(1)}%`,
      p95ResponseTime: this.calculatePercentile(95),
      p99ResponseTime: this.calculatePercentile(99)
    };
  }

  calculatePercentile(percentile) {
    if (this.metrics.responseTimes.length === 0) return 0;
    
    const sorted = [...this.metrics.responseTimes].sort((a, b) => a - b);
    const index = Math.ceil(percentile / 100 * sorted.length) - 1;
    return sorted[index];
  }
}
```

## Key Benefits of Hybrid Search

### 1. **Improved Recall**
- Catches both semantic matches and exact term matches
- Handles vocabulary mismatch problems
- Finds relevant content that pure methods might miss

### 2. **Better Precision**
- Ranks exact matches higher
- Understands user intent through query analysis
- Provides more relevant results overall

### 3. **Robustness**
- Works well for various query types
- Handles ambiguous queries effectively
- Degrades gracefully when one method fails

### 4. **Flexibility**
- Configurable weighting based on use case
- Adaptable to different domains and content types
- Extensible with additional ranking factors

## Best Practices

1. **Start with simple RRF**: Often provides good results with minimal tuning
2. **Analyze your queries**: Understand what types of searches users perform
3. **Monitor performance**: Track what works and adjust weights accordingly
4. **Consider your domain**: Technical content may need different tuning than general content
5. **Implement caching**: Especially for expensive vector operations
6. **Test extensively**: Use A/B testing to validate improvements

Hybrid search represents the state-of-the-art in search technology, combining the best of both traditional and modern approaches to deliver superior search experiences.