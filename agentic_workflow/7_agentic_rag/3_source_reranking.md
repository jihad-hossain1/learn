# Source Reranking in RAG Systems

Source reranking is a critical component in Retrieval-Augmented Generation (RAG) systems that improves the quality of retrieved documents by reordering them based on relevance to the query before passing them to the generator.

## What is Source Reranking?

**Source reranking** is the process of taking initially retrieved documents and reordering them using a more sophisticated (but computationally expensive) model to ensure the most relevant documents appear first in the context window.

## Why Reranking Matters

### Problems with Simple Retrieval:
- **Semantic mismatch**: Vector similarity ≠ relevance
- **Position bias**: Important documents might not be top-ranked
- **Context window limits**: Only top-k documents get used
- **Quality degradation**: Irrelevant documents hurt generation quality

### Benefits of Reranking:
- **Improved precision**: More relevant documents upfront
- **Better context utilization**: Optimal use of limited context windows
- **Higher quality answers**: Generator receives better information
- **Reduced noise**: Filters out irrelevant documents

## Reranking Architectures

### 1. **Two-Stage Retrieval with Reranking**
```javascript
class TwoStageRetriever {
    constructor(initialRetriever, reranker) {
        this.initialRetriever = initialRetriever;
        this.reranker = reranker;
    }

    async retrieve(query, options = {}) {
        const {
            initialTopK = 100,
            finalTopK = 10,
            rerankThreshold = 0.1
        } = options;

        // Stage 1: Broad retrieval
        console.log(`Stage 1: Retrieving ${initialTopK} documents...`);
        const initialResults = await this.initialRetriever.retrieve(
            query, 
            initialTopK
        );

        // Stage 2: Precision reranking
        console.log(`Stage 2: Reranking top ${initialResults.length} documents...`);
        const rerankedResults = await this.reranker.rerank(
            query, 
            initialResults
        );

        // Apply threshold and return top-k
        const filteredResults = rerankedResults.filter(
            doc => doc.rerankScore >= rerankThreshold
        ).slice(0, finalTopK);

        console.log(`Final: Returning ${filteredResults.length} documents`);
        return filteredResults;
    }
}
```

### 2. **Cross-Encoder Reranking**
```javascript
class CrossEncoderReranker {
    constructor(model) {
        this.model = model; // Cross-encoder model
    }

    async rerank(query, documents) {
        // Create query-document pairs for cross-encoding
        const pairs = documents.map(doc => ({
            query: query,
            document: doc.content,
            originalDoc: doc
        }));

        // Score each pair using cross-encoder
        const scoredPairs = await Promise.all(
            pairs.map(async pair => ({
                ...pair,
                relevanceScore: await this.scoreRelevance(pair.query, pair.document)
            }))
        );

        // Sort by relevance score (descending)
        return scoredPairs
            .sort((a, b) => b.relevanceScore - a.relevanceScore)
            .map(pair => ({
                ...pair.originalDoc,
                rerankScore: pair.relevanceScore,
                finalRank: scoredPairs.indexOf(pair) + 1
            }));
    }

    async scoreRelevance(query, document) {
        // Simulate cross-encoder scoring
        // In practice: this would call a trained cross-encoder model
        const prompt = `
        Query: "${query}"
        Document: "${document.substring(0, 500)}"
        
        Score the relevance from 0.0 to 1.0 where:
        0.0 = completely irrelevant
        0.5 = somewhat relevant
        1.0 = perfectly relevant
        
        Return only the score as a number.
        `;
        
        // Mock implementation - real system would use actual model
        const baseScore = Math.random() * 0.3 + 0.6; // Simulate decent relevance
        return Math.min(1.0, baseScore + this.calculateSemanticOverlap(query, document));
    }

    calculateSemanticOverlap(query, document) {
        const queryWords = new Set(query.toLowerCase().split(/\W+/));
        const docWords = new Set(document.toLowerCase().split(/\W+/));
        
        const intersection = new Set([...queryWords].filter(x => docWords.has(x)));
        return intersection.size / queryWords.size;
    }
}
```

## Advanced Reranking Strategies

### 1. **Multi-Factor Reranking**
```javascript
class MultiFactorReranker {
    constructor() {
        this.factors = {
            semanticRelevance: { weight: 0.6 },
            freshness: { weight: 0.15 },
            authority: { weight: 0.15 },
            diversity: { weight: 0.1 }
        };
    }

    async rerank(query, documents, context = {}) {
        const scoredDocuments = await Promise.all(
            documents.map(async doc => {
                const scores = {
                    semantic: await this.scoreSemanticRelevance(query, doc),
                    freshness: this.scoreFreshness(doc, context.currentDate),
                    authority: this.scoreAuthority(doc),
                    diversity: await this.scoreDiversity(doc, documents)
                };

                const finalScore = this.calculateWeightedScore(scores);
                
                return {
                    ...doc,
                    scores: scores,
                    finalScore: finalScore,
                    factors: this.factors
                };
            })
        );

        return scoredDocuments.sort((a, b) => b.finalScore - a.finalScore);
    }

    calculateWeightedScore(scores) {
        return Object.entries(this.factors).reduce((total, [factor, config]) => {
            return total + (scores[factor] * config.weight);
        }, 0);
    }

    async scoreSemanticRelevance(query, document) {
        // Use cross-encoder or semantic similarity
        return Math.random() * 0.4 + 0.5; // Base score + randomness
    }

    scoreFreshness(document, currentDate) {
        if (!document.metadata || !document.metadata.date) return 0.5;
        
        const docDate = new Date(document.metadata.date);
        const current = currentDate ? new Date(currentDate) : new Date();
        const diffTime = Math.abs(current - docDate);
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        
        // Recent documents score higher (exponential decay)
        return Math.exp(-diffDays / 365); // 1 year half-life
    }

    scoreAuthority(document) {
        const authorityScores = {
            'wikipedia.org': 0.9,
            'arxiv.org': 0.8,
            'github.com': 0.7,
            'personal.blog': 0.4
        };

        const source = document.metadata?.source || '';
        for (const [domain, score] of Object.entries(authorityScores)) {
            if (source.includes(domain)) return score;
        }
        
        return 0.5; // Default score
    }

    async scoreDiversity(currentDoc, allDocs) {
        // Penalize documents too similar to higher-ranked ones
        let diversityPenalty = 0;
        const similarDocs = allDocs.filter(doc => 
            doc !== currentDoc && 
            this.calculateSimilarity(currentDoc.content, doc.content) > 0.8
        );
        
        return Math.max(0, 1 - (similarDocs.length * 0.2));
    }

    calculateSimilarity(text1, text2) {
        // Simple Jaccard similarity
        const words1 = new Set(text1.toLowerCase().split(/\W+/));
        const words2 = new Set(text2.toLowerCase().split(/\W+/));
        
        const intersection = new Set([...words1].filter(x => words2.has(x)));
        const union = new Set([...words1, ...words2]);
        
        return intersection.size / union.size;
    }
}
```

### 2. **Learning-to-Rank (LTR) Approach**
```javascript
class LearningToRankReranker {
    constructor(model, features) {
        this.model = model;
        this.featureExtractor = new FeatureExtractor(features);
    }

    async rerank(query, documents) {
        // Extract features for each query-document pair
        const featuresList = await Promise.all(
            documents.map(doc => 
                this.featureExtractor.extractFeatures(query, doc)
            )
        );

        // Predict relevance scores using LTR model
        const scores = await this.model.predict(featuresList);

        // Combine documents with scores and sort
        return documents
            .map((doc, index) => ({
                ...doc,
                ltrScore: scores[index],
                features: featuresList[index]
            }))
            .sort((a, b) => b.ltrScore - a.ltrScore);
    }
}

class FeatureExtractor {
    constructor(featureConfig) {
        this.featureConfig = featureConfig;
    }

    async extractFeatures(query, document) {
        const features = {};

        if (this.featureConfig.semantic) {
            features.semanticSimilarity = await this.extractSemanticFeatures(query, document);
        }

        if (this.featureConfig.lexical) {
            features.lexicalOverlap = this.extractLexicalFeatures(query, document);
        }

        if (this.featureConfig.quality) {
            features.qualityMetrics = this.extractQualityFeatures(document);
        }

        if (this.featureConfig.contextual) {
            features.contextualFeatures = await this.extractContextualFeatures(query, document);
        }

        return features;
    }

    extractLexicalFeatures(query, document) {
        const queryTerms = query.toLowerCase().split(/\W+/);
        const docText = document.content.toLowerCase();
        
        return {
            termOverlap: queryTerms.filter(term => docText.includes(term)).length / queryTerms.length,
            exactMatch: queryTerms.some(term => docText.includes(term)) ? 1 : 0,
            queryLength: query.length,
            documentLength: document.content.length
        };
    }

    extractQualityFeatures(document) {
        return {
            hasCitations: document.metadata?.citations ? 1 : 0,
            readabilityScore: this.calculateReadability(document.content),
            sourceAuthority: this.scoreSourceAuthority(document.metadata?.source)
        };
    }
}
```

## Hybrid Reranking Approaches

### 1. **Cascade Reranking**
```javascript
class CascadeReranker {
    constructor(rerankers) {
        this.rerankers = rerankers; // Array of rerankers in priority order
    }

    async rerank(query, documents, budget = 100) {
        let currentDocs = documents;
        let cost = 0;

        for (const reranker of this.rerankers) {
            if (cost >= budget) break;
            
            const startTime = Date.now();
            currentDocs = await reranker.rerank(query, currentDocs);
            cost += Date.now() - startTime;

            // Early stopping if confidence is high
            if (this.isHighConfidence(currentDocs)) {
                break;
            }
        }

        return currentDocs;
    }

    isHighConfidence(documents) {
        if (documents.length < 2) return true;
        
        const topScore = documents[0].finalScore || documents[0].rerankScore;
        const secondScore = documents[1].finalScore || documents[1].rerankScore;
        
        // High confidence if top document is significantly better
        return (topScore - secondScore) > 0.2;
    }
}
```

### 2. **Ensemble Reranking**
```javascript
class EnsembleReranker {
    constructor(rerankers, weights) {
        this.rerankers = rerankers;
        this.weights = weights;
    }

    async rerank(query, documents) {
        // Get rankings from all rerankers
        const allRankings = await Promise.all(
            this.rerankers.map(reranker => reranker.rerank(query, documents))
        );

        // Combine using Borda count or similar method
        return this.combineRankings(allRankings, documents);
    }

    combineRankings(rankings, originalDocuments) {
        const documentScores = new Map();

        originalDocuments.forEach(doc => {
            documentScores.set(doc.id, { doc, totalScore: 0, ranks: [] });
        });

        // Calculate scores based on positions in each ranking
        rankings.forEach((ranking, rerankerIndex) => {
            ranking.forEach((doc, position) => {
                const score = documentScores.get(doc.id);
                if (score) {
                    // Higher positions get more points
                    const points = originalDocuments.length - position;
                    const weightedPoints = points * (this.weights[rerankerIndex] || 1);
                    score.totalScore += weightedPoints;
                    score.ranks.push(position + 1);
                }
            });
        });

        // Sort by combined score
        return Array.from(documentScores.values())
            .sort((a, b) => b.totalScore - a.totalScore)
            .map(item => ({
                ...item.doc,
                ensembleScore: item.totalScore,
                individualRanks: item.ranks
            }));
    }
}
```

## Practical Implementation Example

```javascript
class PracticalRerankingSystem {
    constructor(vectorStore, options = {}) {
        this.vectorStore = vectorStore;
        this.options = {
            initialRetrievalK: 50,
            finalK: 7,
            rerankThreshold: 0.3,
            useMultiFactor: true,
            ...options
        };

        this.setupRerankers();
    }

    setupRerankers() {
        this.initialRetriever = new SimpleRetriever(this.vectorStore);
        
        if (this.options.useMultiFactor) {
            this.reranker = new MultiFactorReranker();
        } else {
            this.reranker = new CrossEncoderReranker();
        }

        // Setup cascade for expensive rerankers
        this.cascadeReranker = new CascadeReranker([
            new FastReranker(),    // Quick semantic check
            this.reranker,         // Full multi-factor rerank
            new ExpensiveReranker() // Only if needed
        ]);
    }

    async searchWithReranking(query, context = {}) {
        try {
            // Step 1: Initial broad retrieval
            const initialResults = await this.initialRetriever.retrieve(
                query, 
                this.options.initialRetrievalK
            );

            if (initialResults.length === 0) {
                return { documents: [], rerankingApplied: false };
            }

            // Step 2: Apply reranking
            const rerankedResults = await this.cascadeReranker.rerank(
                query, 
                initialResults, 
                context
            );

            // Step 3: Filter and return top-k
            const finalResults = rerankedResults
                .filter(doc => doc.finalScore >= this.options.rerankThreshold)
                .slice(0, this.options.finalK);

            return {
                documents: finalResults,
                rerankingApplied: true,
                stats: {
                    initialCount: initialResults.length,
                    finalCount: finalResults.length,
                    reductionRatio: (initialResults.length - finalResults.length) / initialResults.length
                }
            };

        } catch (error) {
            console.error('Reranking failed:', error);
            // Fallback to simple retrieval
            const fallbackResults = await this.initialRetriever.retrieve(query, this.options.finalK);
            return { documents: fallbackResults, rerankingApplied: false, error: error.message };
        }
    }
}

// Usage example
async function demonstrateReranking() {
    const system = new PracticalRerankingSystem(vectorStore, {
        initialRetrievalK: 100,
        finalK: 5,
        rerankThreshold: 0.4
    });

    const query = "What are the latest advancements in quantum computing?";
    const context = { currentDate: '2024-01-15' };

    const results = await system.searchWithReranking(query, context);
    
    console.log('Reranking Results:');
    console.log(`Initial documents: ${results.stats.initialCount}`);
    console.log(`Final documents: ${results.stats.finalCount}`);
    console.log(`Reduction: ${(results.stats.reductionRatio * 100).toFixed(1)}%`);
    
    results.documents.forEach((doc, index) => {
        console.log(`${index + 1}. Score: ${doc.finalScore.toFixed(3)} - ${doc.metadata.title}`);
    });
}
```

## Key Benefits and Trade-offs

### Benefits:
- **Significantly improved relevance**
- **Better context window utilization**
- **Higher answer quality**
- **Reduced hallucination**

### Trade-offs:
- **Increased latency** (additional processing)
- **Higher computational cost**
- **Complexity in implementation**
- **Model dependency**

## Best Practices

1. **Start with simple reranking** and gradually add complexity
2. **Monitor reranking performance** with A/B testing
3. **Implement fallback mechanisms** for reranking failures
4. **Cache frequent query-document scores**
5. **Balance precision vs. latency** based on use case

Source reranking is a powerful technique that can dramatically improve RAG system performance by ensuring the generator receives the most relevant information in optimal order.