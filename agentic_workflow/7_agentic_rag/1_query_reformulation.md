# Agentic RAG in Query Reformulation

Agentic RAG (Retrieval-Augmented Generation) is an advanced approach where an AI agent actively participates in the retrieval process, including query reformulation, to improve the quality of results.

## What is Agentic RAG Query Reformulation?

Traditional RAG takes a user query and retrieves relevant documents. Agentic RAG adds an intelligent agent that can:
- **Analyze** the original query
- **Reformulate** it into multiple variations
- **Select** the best reformulation strategy
- **Execute** searches with different query versions
- **Synthesize** results from multiple searches

## Key Benefits

1. **Improved recall** - Multiple query variations increase chances of finding relevant documents
2. **Better precision** - Intelligent reformulation targets the most relevant information
3. **Context awareness** - The agent understands the conversation context
4. **Adaptive searching** - Can adjust strategy based on initial results

## JavaScript Implementation Example

Here's a simplified example of Agentic RAG with query reformulation:

```javascript
class AgenticRAG {
    constructor() {
        this.queryStrategies = {
            'synonym_expansion': this.synonymExpansion,
            'query_decomposition': this.queryDecomposition,
            'contextual_enhancement': this.contextualEnhancement,
            'technical_reformulation': this.technicalReformulation
        };
    }

    // Main agent function
    async processQuery(originalQuery, context = {}) {
        console.log(`Original query: "${originalQuery}"`);
        
        // Step 1: Analyze query intent
        const intent = await this.analyzeIntent(originalQuery, context);
        
        // Step 2: Generate query variations
        const queryVariations = await this.generateQueryVariations(originalQuery, intent);
        
        // Step 3: Execute searches with different variations
        const results = await this.executeMultiQuerySearch(queryVariations);
        
        // Step 4: Synthesize and rank results
        const finalResults = await this.synthesizeResults(results, originalQuery);
        
        return finalResults;
    }

    // Intent analysis
    async analyzeIntent(query, context) {
        // Simple intent classification - in practice, use ML model
        const intent = {
            type: this.classifyIntentType(query),
            complexity: this.assessComplexity(query),
            domain: this.identifyDomain(query, context),
            technicality: this.assessTechnicality(query)
        };
        
        console.log('Detected intent:', intent);
        return intent;
    }

    // Generate multiple query variations
    async generateQueryVariations(originalQuery, intent) {
        const variations = [originalQuery]; // Always include original
        
        // Apply different reformulation strategies based on intent
        for (const [strategyName, strategyFn] of Object.entries(this.queryStrategies)) {
            if (this.shouldApplyStrategy(strategyName, intent)) {
                const strategyVariations = await strategyFn(originalQuery, intent);
                variations.push(...strategyVariations);
            }
        }
        
        console.log('Generated variations:', variations);
        return [...new Set(variations)]; // Remove duplicates
    }

    // Query reformulation strategies
    async synonymExpansion(query, intent) {
        const synonyms = {
            'how to': ['ways to', 'methods for', 'steps to'],
            'best': ['top', 'optimal', 'most effective'],
            'javascript': ['js', 'ecmascript'],
            'error': ['issue', 'problem', 'bug']
        };
        
        const variations = [query];
        
        for (const [word, alternatives] of Object.entries(synonyms)) {
            if (query.toLowerCase().includes(word)) {
                for (const alt of alternatives) {
                    variations.push(query.replace(new RegExp(word, 'i'), alt));
                }
            }
        }
        
        return variations;
    }

    async queryDecomposition(query, intent) {
        if (intent.complexity === 'high') {
            // Split complex queries into simpler components
            const parts = query.split(/(?:and|or|but|however)/i);
            return parts.filter(part => part.trim().length > 5);
        }
        return [query];
    }

    async contextualEnhancement(query, intent) {
        const enhancements = [];
        
        if (intent.domain === 'programming') {
            enhancements.push(`${query} code example`);
            enhancements.push(`${query} implementation`);
        }
        
        if (intent.technicality === 'high') {
            enhancements.push(`${query} technical details`);
            enhancements.push(`${query} advanced tutorial`);
        }
        
        return enhancements;
    }

    async technicalReformulation(query, intent) {
        if (intent.domain === 'programming') {
            return [
                `${query} algorithm`,
                `${query} time complexity`,
                `${query} best practices`
            ];
        }
        return [query];
    }

    // Strategy selection logic
    shouldApplyStrategy(strategy, intent) {
        const strategyRules = {
            'synonym_expansion': intent.complexity !== 'very low',
            'query_decomposition': intent.complexity === 'high',
            'contextual_enhancement': intent.domain !== 'general',
            'technical_reformulation': intent.technicality === 'high'
        };
        
        return strategyRules[strategy] || false;
    }

    // Mock search execution
    async executeMultiQuerySearch(queries) {
        const results = {};
        
        for (const query of queries) {
            // In practice, this would call your vector database or search API
            results[query] = await this.mockSearch(query);
        }
        
        return results;
    }

    // Result synthesis
    async synthesizeResults(results, originalQuery) {
        // Combine, deduplicate, and rank results from all queries
        const allDocuments = [];
        
        for (const [query, docs] of Object.entries(results)) {
            docs.forEach(doc => {
                doc.querySource = query;
                doc.relevanceScore = this.calculateRelevance(doc, originalQuery);
                allDocuments.push(doc);
            });
        }
        
        // Remove duplicates and sort by relevance
        const uniqueDocs = this.removeDuplicates(allDocuments);
        return uniqueDocs.sort((a, b) => b.relevanceScore - a.relevanceScore);
    }

    // Helper methods (simplified for example)
    classifyIntentType(query) {
        if (query.toLowerCase().includes('how to')) return 'tutorial';
        if (query.toLowerCase().includes('what is')) return 'definition';
        if (query.toLowerCase().includes('best')) return 'comparison';
        return 'general';
    }

    assessComplexity(query) {
        const wordCount = query.split(' ').length;
        if (wordCount > 8) return 'high';
        if (wordCount > 4) return 'medium';
        return 'low';
    }

    identifyDomain(query, context) {
        const programmingKeywords = ['code', 'program', 'function', 'variable', 'javascript'];
        return programmingKeywords.some(keyword => 
            query.toLowerCase().includes(keyword)) ? 'programming' : 'general';
    }

    assessTechnicality(query) {
        const technicalTerms = ['algorithm', 'complexity', 'implementation', 'optimization'];
        return technicalTerms.some(term => 
            query.toLowerCase().includes(term)) ? 'high' : 'medium';
    }

    async mockSearch(query) {
        // Simulate search results
        return [
            { content: `Document about ${query}`, score: Math.random() },
            { content: `Another result for ${query}`, score: Math.random() }
        ];
    }

    calculateRelevance(doc, originalQuery) {
        // Simple relevance calculation
        const words = originalQuery.toLowerCase().split(' ');
        const matches = words.filter(word => 
            doc.content.toLowerCase().includes(word)).length;
        return matches / words.length;
    }

    removeDuplicates(documents) {
        const seen = new Set();
        return documents.filter(doc => {
            const signature = doc.content.substring(0, 50);
            if (seen.has(signature)) return false;
            seen.add(signature);
            return true;
        });
    }
}

// Usage example
async function demonstrateAgenticRAG() {
    const agenticRAG = new AgenticRAG();
    
    const query = "How to handle async errors in JavaScript";
    const context = { previousQueries: [], userExpertise: 'intermediate' };
    
    console.log('=== Agentic RAG Demo ===');
    const results = await agenticRAG.processQuery(query, context);
    
    console.log('\n=== Final Results ===');
    results.forEach((result, index) => {
        console.log(`${index + 1}. ${result.content} (Score: ${result.relevanceScore.toFixed(2)})`);
    });
}

// Run the demo
demonstrateAgenticRAG();
```

## Advanced Example with Real Search Integration

```javascript
class AdvancedAgenticRAG {
    constructor(searchEngine, llmClient) {
        this.searchEngine = searchEngine; // Your vector DB or search API
        this.llm = llmClient; // LLM for intelligent reformulation
    }

    async intelligentQueryReformulation(originalQuery, context) {
        // Use LLM to generate better reformulations
        const prompt = `
        Original query: "${originalQuery}"
        Context: ${JSON.stringify(context)}
        
        Generate 3-5 alternative query formulations that might retrieve better results.
        Consider: synonyms, technical terms, broader/narrower scopes.
        
        Return as JSON array:
        ["query1", "query2", ...]
        `;
        
        const response = await this.llm.generate(prompt);
        return JSON.parse(response);
    }

    async adaptiveSearch(queries, originalQuery) {
        let allResults = [];
        
        for (const query of queries) {
            const results = await this.searchEngine.search(query);
            
            // Early termination if we find high-quality results
            const highQualityResults = results.filter(r => 
                this.calculateQualityScore(r, originalQuery) > 0.8);
            
            if (highQualityResults.length >= 3) {
                console.log(`Found high-quality results for: "${query}"`);
                allResults.push(...highQualityResults);
                break; // Stop searching if we have good results
            }
            
            allResults.push(...results);
        }
        
        return allResults;
    }
}

// Example with real API calls (conceptual)
class ProductionAgenticRAG {
    async processQueryWithAPIs(query) {
        // 1. Analyze with NLP service
        const analysis = await this.analyzeWithNLP(query);
        
        // 2. Generate variations with LLM
        const variations = await this.generateVariationsWithLLM(query, analysis);
        
        // 3. Search with multiple queries
        const results = await Promise.all(
            variations.map(v => this.vectorDBSearch(v))
        );
        
        // 4. Rerank with cross-encoder
        const reranked = await this.rerankResults(results.flat(), query);
        
        return reranked;
    }
}
```

## Key Features of Agentic RAG Query Reformulation

1. **Multi-strategy approach**: Uses different reformulation techniques
2. **Intent-aware**: Adapts based on what the user really needs
3. **Context-sensitive**: Considers conversation history and user context
4. **Quality-aware**: Can stop early when good results are found
5. **Synthetic combination**: Merges results from multiple queries intelligently

This approach significantly improves over basic RAG by making the retrieval process more intelligent and adaptive to the specific query and context.