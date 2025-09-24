# Adaptive Retrieval Strategies

Adaptive retrieval strategies dynamically adjust retrieval parameters and approaches based on query characteristics, context, and real-time feedback to optimize RAG system performance.

## What are Adaptive Retrieval Strategies?

**Adaptive retrieval** refers to intelligent systems that modify their retrieval behavior in response to:
- Query complexity and type
- User context and history
- Document collection characteristics
- Real-time performance feedback
- Resource constraints

## Core Adaptive Strategies

### 1. **Query Complexity-Based Adaptation**
```javascript
class ComplexityAdaptiveRetriever {
    constructor(vectorStore, llm) {
        this.vectorStore = vectorStore;
        this.llm = llm;
        this.complexityThresholds = {
            low: { k: 3, similarity: 0.7 },
            medium: { k: 7, similarity: 0.6 },
            high: { k: 15, similarity: 0.5 },
            very_high: { k: 25, similarity: 0.4 }
        };
    }

    async adaptiveRetrieve(query, context = {}) {
        // Analyze query complexity
        const complexity = await this.analyzeQueryComplexity(query, context);
        console.log(`Detected complexity: ${complexity.level}`);
        
        // Adjust retrieval parameters
        const params = this.complexityThresholds[complexity.level];
        let results = await this.vectorStore.similaritySearch(query, params.k);
        
        // Apply complexity-specific strategies
        if (complexity.level === 'high' || complexity.level === 'very_high') {
            results = await this.handleComplexQuery(query, results, complexity);
        }
        
        return {
            documents: results,
            complexity: complexity,
            parameters: params
        };
    }

    async analyzeQueryComplexity(query, context) {
        const prompt = `
        Analyze this query's complexity:
        Query: "${query}"
        Context: ${JSON.stringify(context)}
        
        Consider:
        - Number of concepts involved
        - Specificity vs breadth
        - Technical depth required
        - Need for multi-step reasoning
        
        Return JSON: {level: "low"|"medium"|"high"|"very_high", reasoning: string, factors: string[]}
        `;
        
        // Mock analysis - real implementation would use LLM
        const factors = [];
        const words = query.split(' ');
        const hasTechnicalTerms = /(algorithm|complexity|implementation|optimization)/i.test(query);
        
        if (words.length > 10) factors.push('long_query');
        if (hasTechnicalTerms) factors.push('technical');
        if (query.includes('compare') || query.includes('vs')) factors.push('comparative');
        if (query.includes('how to') && query.includes('and')) factors.push('multi_step');
        
        let level = 'medium';
        if (factors.length === 0) level = 'low';
        if (factors.length >= 2) level = 'high';
        if (factors.includes('multi_step') && factors.includes('technical')) level = 'very_high';
        
        return { level, factors, reasoning: `Based on factors: ${factors.join(', ')}` };
    }

    async handleComplexQuery(query, initialResults, complexity) {
        // For complex queries, employ advanced strategies
        const strategies = [];
        
        if (complexity.factors.includes('multi_step')) {
            strategies.push('query_decomposition');
        }
        
        if (complexity.factors.includes('comparative')) {
            strategies.push('multi_perspective');
        }
        
        if (complexity.factors.includes('technical')) {
            strategies.push('expert_sources');
        }
        
        let enhancedResults = initialResults;
        
        for (const strategy of strategies) {
            enhancedResults = await this.applyRetrievalStrategy(
                strategy, query, enhancedResults
            );
        }
        
        return enhancedResults;
    }

    async applyRetrievalStrategy(strategy, query, currentResults) {
        switch (strategy) {
            case 'query_decomposition':
                return await this.decomposeAndRetrieve(query, currentResults);
            case 'multi_perspective':
                return await this.multiPerspectiveRetrieval(query, currentResults);
            case 'expert_sources':
                return await this.expertSourceFocus(query, currentResults);
            default:
                return currentResults;
        }
    }

    async decomposeAndRetrieve(query, currentResults) {
        const subQueries = await this.decomposeQuery(query);
        const allResults = [...currentResults];
        
        for (const subQuery of subQueries) {
            const subResults = await this.vectorStore.similaritySearch(subQuery, 5);
            allResults.push(...subResults);
        }
        
        return this.deduplicateAndRank(allResults, query);
    }

    async multiPerspectiveRetrieval(query, currentResults) {
        const perspectives = await this.generatePerspectives(query);
        const allResults = [...currentResults];
        
        for (const perspective of perspectives) {
            const perspectiveQuery = `${query} ${perspective}`;
            const perspectiveResults = await this.vectorStore.similaritySearch(perspectiveQuery, 3);
            allResults.push(...perspectiveResults);
        }
        
        return this.deduplicateAndRank(allResults, query);
    }
}
```

### 2. **Context-Aware Adaptive Retrieval**
```javascript
class ContextAwareRetriever {
    constructor(vectorStore, contextManager) {
        this.vectorStore = vectorStore;
        this.contextManager = contextManager;
        this.adaptationRules = {
            follow_up: this.adaptToFollowUpQuery,
            clarification: this.adaptToClarification,
            topic_shift: this.adaptToTopicShift,
            deepening: this.adaptToDeepeningQuery
        };
    }

    async adaptiveRetrieve(query, sessionContext) {
        // Analyze conversation context
        const contextAnalysis = await this.analyzeConversationContext(query, sessionContext);
        const adaptationType = contextAnalysis.adaptationNeeded;
        
        console.log(`Context analysis: ${adaptationType}`);
        
        // Apply context-specific adaptation
        const adaptationRule = this.adaptationRules[adaptationType];
        let results = await this.vectorStore.similaritySearch(query, 10);
        
        if (adaptationRule) {
            results = await adaptationRule.call(this, query, results, sessionContext, contextAnalysis);
        }
        
        return {
            documents: results,
            contextAnalysis: contextAnalysis,
            adaptationApplied: adaptationType
        };
    }

    async analyzeConversationContext(currentQuery, sessionContext) {
        const { previousQueries = [], previousResults = [], conversationHistory = [] } = sessionContext;
        
        if (previousQueries.length === 0) {
            return { adaptationNeeded: 'none', reasoning: 'First query in session' };
        }
        
        const lastQuery = previousQueries[previousQueries.length - 1];
        const lastResults = previousResults[previousResults.length - 1];
        
        // Analyze relationship with previous query
        const relationship = await this.analyzeQueryRelationship(lastQuery, currentQuery);
        
        if (relationship.type === 'follow_up') {
            return {
                adaptationNeeded: 'follow_up',
                relationship: relationship,
                focusArea: relationship.focus
            };
        }
        
        if (relationship.type === 'clarification') {
            return {
                adaptationNeeded: 'clarification',
                relationship: relationship,
                clarificationType: relationship.clarificationType
            };
        }
        
        if (relationship.type === 'topic_shift') {
            return {
                adaptationNeeded: 'topic_shift',
                relationship: relationship,
                newTopic: relationship.newTopic
            };
        }
        
        return { adaptationNeeded: 'none', relationship: relationship };
    }

    async adaptToFollowUpQuery(query, currentResults, sessionContext, analysis) {
        const { previousResults, previousQueries } = sessionContext;
        const lastResults = previousResults[previousResults.length - 1];
        
        // Build upon previous results
        const enhancedResults = [...currentResults];
        
        // Include relevant documents from previous search
        const relevantPreviousDocs = lastResults.documents.filter(doc => 
            this.calculateRelevanceToQuery(doc, query) > 0.3
        );
        
        enhancedResults.push(...relevantPreviousDocs);
        
        // Focus on the specific aspect mentioned in follow-up
        if (analysis.focusArea) {
            const focusedQuery = `${query} ${analysis.focusArea}`;
            const focusedResults = await this.vectorStore.similaritySearch(focusedQuery, 5);
            enhancedResults.push(...focusedResults);
        }
        
        return this.deduplicateAndRank(enhancedResults, query);
    }

    async adaptToTopicShift(query, currentResults, sessionContext, analysis) {
        // For topic shifts, we need broader retrieval initially
        const broadQuery = await this.broadenQuery(query);
        const broadResults = await this.vectorStore.similaritySearch(broadQuery, 15);
        
        const allResults = [...currentResults, ...broadResults];
        
        // Use diversity ranking to explore new topic
        return this.rankForDiversity(allResults, query);
    }

    async broadenQuery(query) {
        // Add broader terms to explore new topic space
        const broadeningTerms = {
            'introduction': ['overview', 'basics', 'fundamentals'],
            'technical': ['explanation', 'guide', 'tutorial'],
            'comparison': ['vs', 'differences', 'advantages']
        };
        
        let broadenedQuery = query;
        for (const [key, terms] of Object.entries(broadeningTerms)) {
            if (query.toLowerCase().includes(key)) {
                broadenedQuery += ' ' + terms[0];
                break;
            }
        }
        
        return broadenedQuery;
    }

    rankForDiversity(documents, query, maxSimilar = 0.7) {
        // Maximize diversity while maintaining relevance
        const ranked = [];
        const used = new Set();
        
        documents.forEach(doc => {
            if (!used.has(doc.id)) {
                ranked.push(doc);
                used.add(doc.id);
                
                // Mark very similar documents as used
                documents.forEach(otherDoc => {
                    if (!used.has(otherDoc.id) && 
                        this.calculateSimilarity(doc.content, otherDoc.content) > maxSimilar) {
                        used.add(otherDoc.id);
                    }
                });
            }
        });
        
        return ranked.slice(0, 10);
    }
}
```

### 3. **Performance-Adaptive Retrieval**
```javascript
class PerformanceAdaptiveRetriever {
    constructor(vectorStore, config = {}) {
        this.vectorStore = vectorStore;
        this.performanceHistory = [];
        this.config = {
            historySize: 100,
            adaptationWindow: 20,
            minConfidenceThreshold: 0.6,
            ...config
        };
    }

    async adaptiveRetrieve(query, context = {}) {
        // Check recent performance to adapt strategy
        const performanceAnalysis = this.analyzeRecentPerformance();
        const adaptation = this.determineAdaptation(performanceAnalysis);
        
        console.log(`Performance-based adaptation: ${adaptation.strategy}`);
        
        // Apply adaptive strategy
        let results;
        switch (adaptation.strategy) {
            case 'conservative':
                results = await this.conservativeRetrieval(query, adaptation.parameters);
                break;
            case 'aggressive':
                results = await this.aggressiveRetrieval(query, adaptation.parameters);
                break;
            case 'balanced':
            default:
                results = await this.balancedRetrieval(query, adaptation.parameters);
        }
        
        // Store query for future performance tracking
        this.storeQueryMetrics(query, results.length, adaptation.strategy);
        
        return {
            documents: results,
            adaptation: adaptation,
            performanceMetrics: performanceAnalysis
        };
    }

    analyzeRecentPerformance() {
        if (this.performanceHistory.length < 5) {
            return { confidence: 'high', trend: 'stable', recommendation: 'balanced' };
        }
        
        const recent = this.performanceHistory.slice(-this.config.adaptationWindow);
        const successRate = recent.filter(entry => entry.success).length / recent.length;
        const avgResults = recent.reduce((sum, entry) => sum + entry.resultCount, 0) / recent.length;
        
        let confidence, trend, recommendation;
        
        if (successRate > 0.8) {
            confidence = 'high';
            recommendation = 'aggressive';
        } else if (successRate > 0.6) {
            confidence = 'medium';
            recommendation = 'balanced';
        } else {
            confidence = 'low';
            recommendation = 'conservative';
        }
        
        // Analyze trend
        const recentSuccess = recent.slice(-10).filter(entry => entry.success).length;
        const previousSuccess = recent.slice(-20, -10).filter(entry => entry.success).length;
        trend = recentSuccess > previousSuccess ? 'improving' : 
                recentSuccess < previousSuccess ? 'declining' : 'stable';
        
        return { confidence, trend, successRate, avgResults, recommendation };
    }

    determineAdaptation(performance) {
        const baseConfig = {
            k: 10,
            similarityThreshold: 0.5,
            enableReranking: true
        };
        
        switch (performance.recommendation) {
            case 'aggressive':
                return {
                    strategy: 'aggressive',
                    parameters: {
                        ...baseConfig,
                        k: 15,
                        similarityThreshold: 0.4,
                        enableReranking: false // Faster but less precise
                    }
                };
                
            case 'conservative':
                return {
                    strategy: 'conservative',
                    parameters: {
                        ...baseConfig,
                        k: 7,
                        similarityThreshold: 0.6,
                        enableReranking: true,
                        maxRerankTime: 1000
                    }
                };
                
            default: // balanced
                return {
                    strategy: 'balanced',
                    parameters: baseConfig
                };
        }
    }

    async conservativeRetrieval(query, parameters) {
        // High precision, lower recall approach
        let results = await this.vectorStore.similaritySearch(query, parameters.k);
        
        // Apply strict filtering
        results = results.filter(doc => doc.similarityScore >= parameters.similarityThreshold);
        
        if (parameters.enableReranking) {
            results = await this.applyReranking(results, query, parameters.maxRerankTime);
        }
        
        return results;
    }

    async aggressiveRetrieval(query, parameters) {
        // High recall, lower precision approach
        let results = await this.vectorStore.similaritySearch(query, parameters.k);
        
        // Apply minimal filtering
        results = results.filter(doc => doc.similarityScore >= parameters.similarityThreshold);
        
        // Focus on coverage rather than precision
        if (results.length < parameters.k / 2) {
            // If too few results, broaden search
            const broadResults = await this.vectorStore.similaritySearch(query, parameters.k * 2);
            results = [...results, ...broadResults];
        }
        
        return this.deduplicateAndRank(results, query);
    }

    storeQueryMetrics(query, resultCount, strategy) {
        // In real implementation, this would track actual success metrics
        // For demo, we'll simulate based on query characteristics
        const success = this.simulateSuccessMetric(query, resultCount);
        
        this.performanceHistory.push({
            query: query.substring(0, 50),
            timestamp: Date.now(),
            strategy,
            resultCount,
            success
        });
        
        // Maintain history size
        if (this.performanceHistory.length > this.config.historySize) {
            this.performanceHistory.shift();
        }
    }

    simulateSuccessMetric(query, resultCount) {
        // Simulate success based on query length and results
        // Real system would use actual user feedback or downstream success metrics
        const baseProbability = 0.7;
        const lengthFactor = Math.min(1, query.length / 50); // Longer queries harder
        const resultsFactor = Math.min(1, resultCount / 5); // More results better
        
        return Math.random() < (baseProbability * (1 - lengthFactor) * resultsFactor);
    }
}
```

### 4. **Resource-Aware Adaptive Retrieval**
```javascript
class ResourceAwareRetriever {
    constructor(vectorStore, resourceMonitor) {
        this.vectorStore = vectorStore;
        this.resourceMonitor = resourceMonitor;
        this.qualityDegradationTable = {
            'high': { k: 10, reranking: true, methods: ['semantic', 'keyword'] },
            'medium': { k: 7, reranking: false, methods: ['semantic'] },
            'low': { k: 5, reranking: false, methods: ['keyword'] },
            'critical': { k: 3, reranking: false, methods: ['keyword_fast'] }
        };
    }

    async adaptiveRetrieve(query, context = {}) {
        // Check current resource levels
        const resourceStatus = await this.resourceMonitor.getStatus();
        const adaptationLevel = this.determineAdaptationLevel(resourceStatus);
        
        console.log(`Resource status: ${resourceStatus.level}, Adaptation: ${adaptationLevel}`);
        
        // Get adaptation parameters
        const params = this.qualityDegradationTable[adaptationLevel];
        
        // Execute retrieval with adapted parameters
        const results = await this.executeResourceAwareRetrieval(query, params, resourceStatus);
        
        return {
            documents: results,
            resourceStatus: resourceStatus,
            adaptationLevel: adaptationLevel,
            parameters: params
        };
    }

    determineAdaptationLevel(resourceStatus) {
        const { cpu, memory, latency, cost } = resourceStatus;
        
        if (cost > 0.9 || latency > 5000) return 'critical';
        if (cpu > 0.8 || memory > 0.85) return 'low';
        if (cpu > 0.6 || memory > 0.7) return 'medium';
        return 'high';
    }

    async executeResourceAwareRetrieval(query, params, resourceStatus) {
        let allResults = [];
        
        for (const method of params.methods) {
            if (this.shouldContinue(retrieveStatus, resourceStatus)) {
                const methodResults = await this.executeRetrievalMethod(method, query, params);
                allResults = this.mergeResults(allResults, methodResults);
            }
        }
        
        if (params.reranking && this.shouldContinue(retrieveStatus, resourceStatus)) {
            allResults = await this.applyLightReranking(allResults, query);
        }
        
        return allResults.slice(0, params.k);
    }

    async executeRetrievalMethod(method, query, params) {
        switch (method) {
            case 'semantic':
                return await this.vectorStore.similaritySearch(query, params.k);
            case 'keyword':
                return await this.keywordSearch(query, params.k);
            case 'keyword_fast':
                return await this.fastKeywordSearch(query, params.k);
            default:
                return [];
        }
    }

    shouldContinue(retrieveStatus, resourceStatus) {
        // Check if we should continue based on resource constraints
        if (resourceStatus.cost > 0.95) return false;
        if (resourceStatus.latency > 10000) return false;
        if (retrieveStatus.documentsCount >= 20) return false; // Safety limit
        
        return true;
    }

    applyLightReranking(documents, query, maxTime = 100) {
        // Fast, resource-light reranking
        return documents.sort((a, b) => {
            const scoreA = this.fastRelevanceScore(a, query);
            const scoreB = this.fastRelevanceScore(b, query);
            return scoreB - scoreA;
        });
    }

    fastRelevanceScore(document, query) {
        // Fast heuristic-based scoring
        const queryTerms = query.toLowerCase().split(' ');
        const docText = document.content.toLowerCase();
        
        let score = 0;
        for (const term of queryTerms) {
            if (docText.includes(term)) {
                score += 1;
                // Bonus for exact matches
                if (docText.includes(` ${term} `)) {
                    score += 0.5;
                }
            }
        }
        
        return score;
    }
}
```

### 5. **Hybrid Adaptive Retrieval System**
```javascript
class HybridAdaptiveRetriever {
    constructor(components) {
        this.complexityRetriever = components.complexityRetriever;
        this.contextRetriever = components.contextRetriever;
        this.performanceRetriever = components.performanceRetriever;
        this.resourceRetriever = components.resourceRetriever;
        
        this.adaptationWeights = {
            complexity: 0.3,
            context: 0.3,
            performance: 0.25,
            resources: 0.15
        };
    }

    async adaptiveRetrieve(query, context = {}) {
        // Get recommendations from all adaptive strategies
        const recommendations = await this.getAdaptationRecommendations(query, context);
        
        // Combine recommendations using weighted strategy
        const finalStrategy = this.combineStrategies(recommendations);
        
        // Execute retrieval with combined strategy
        const results = await this.executeCombinedStrategy(query, context, finalStrategy);
        
        return {
            documents: results,
            strategy: finalStrategy,
            recommendations: recommendations,
            combinedScore: this.calculateCombinedScore(recommendations)
        };
    }

    async getAdaptationRecommendations(query, context) {
        const [complexityRec, contextRec, performanceRec, resourceRec] = await Promise.all([
            this.complexityRetriever.adaptiveRetrieve(query, context),
            this.contextRetriever.adaptiveRetrieve(query, context),
            this.performanceRetriever.adaptiveRetrieve(query, context),
            this.resourceRetriever.adaptiveRetrieve(query, context)
        ]);
        
        return {
            complexity: complexityRec.complexity,
            context: contextRec.contextAnalysis,
            performance: performanceRec.performanceMetrics,
            resources: resourceRec.resourceStatus,
            individualResults: { complexityRec, contextRec, performanceRec, resourceRec }
        };
    }

    combineStrategies(recommendations) {
        const strategyScores = {
            aggressive: 0,
            conservative: 0,
            balanced: 0
        };
        
        // Score based on complexity
        if (recommendations.complexity.level === 'very_high') {
            strategyScores.aggressive += this.adaptationWeights.complexity;
        } else if (recommendations.complexity.level === 'low') {
            strategyScores.conservative += this.adaptationWeights.complexity;
        } else {
            strategyScores.balanced += this.adaptationWeights.complexity;
        }
        
        // Score based on context
        if (recommendations.context.adaptationNeeded === 'topic_shift') {
            strategyScores.aggressive += this.adaptationWeights.context;
        } else {
            strategyScores.balanced += this.adaptationWeights.context;
        }
        
        // Score based on performance
        strategyScores[recommendations.performance.recommendation] += 
            this.adaptationWeights.performance;
            
        // Score based on resources
        if (recommendations.resources.level === 'critical') {
            strategyScores.conservative += this.adaptationWeights.resources;
        } else {
            strategyScores.balanced += this.adaptationWeights.resources;
        }
        
        // Return strategy with highest score
        return Object.keys(strategyScores).reduce((a, b) => 
            strategyScores[a] > strategyScores[b] ? a : b
        );
    }

    calculateCombinedScore(recommendations) {
        // Calculate overall confidence score
        let score = 0;
        
        // Complexity confidence (higher for medium complexity)
        if (recommendations.complexity.level === 'medium') score += 0.3;
        
        // Context confidence
        if (recommendations.context.adaptationNeeded === 'none') score += 0.3;
        
        // Performance confidence
        if (recommendations.performance.confidence === 'high') score += 0.25;
        
        // Resource confidence
        if (recommendations.resources.level === 'high') score += 0.15;
        
        return Math.min(1, score);
    }
}
```

## Key Benefits of Adaptive Retrieval

1. **Improved Precision**: Right strategy for each query type
2. **Better Resource Utilization**: Adapts to system constraints
3. **Enhanced User Experience**: Context-aware responses
4. **Scalability**: Handles varying load intelligently
5. **Continuous Improvement**: Learns from performance feedback

## Implementation Considerations

- **Monitoring**: Track adaptation effectiveness
- **Fallbacks**: Always have conservative fallback strategies
- **Testing**: A/B test different adaptation strategies
- **Explainability**: Log adaptation decisions for debugging
- **User Control**: Allow users to override adaptations when needed

Adaptive retrieval strategies make RAG systems more intelligent and efficient by dynamically adjusting to the specific needs of each query and the current system context.