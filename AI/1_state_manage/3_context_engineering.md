# Context Engineering: The Art of Smart Information Management

Context Engineering is the deliberate design and management of contextual information to improve AI agent performance, efficiency, and relevance. It's about being strategic about **what information to include, how to structure it, and when to use it**.

## Core Principles of Context Engineering

### 1. **Relevance Filtering**
Only include context that's actually useful for the current task.

### 2. **Information Prioritization**
Place the most important information where it's most effective.

### 3. **Context Compression**
Reduce noise while preserving signal.

### 4. **Dynamic Context Loading**
Load context just-in-time rather than all-at-once.

---

## JavaScript Examples

### Example 1: Basic Context Manager

```javascript
class ContextEngineer {
    constructor() {
        this.contexts = new Map();
        this.maxContextSize = 1000; // characters
        this.priorityWeights = {
            recent: 1.5,
            userPreference: 2.0,
            critical: 3.0,
            historical: 0.7
        };
    }

    // Add context with metadata
    addContext(key, content, metadata = {}) {
        this.contexts.set(key, {
            content,
            timestamp: Date.now(),
            accessCount: 0,
            priority: metadata.priority || 1.0,
            category: metadata.category || 'general',
            expiresAt: metadata.expiresAt || null,
            ...metadata
        });
    }

    // Smart context retrieval with relevance scoring
    getRelevantContext(query, maxTokens = 500) {
        const scoredContexts = [];
        
        for (const [key, context] of this.contexts) {
            if (this.isExpired(context)) {
                this.contexts.delete(key);
                continue;
            }

            const relevanceScore = this.calculateRelevance(context, query);
            if (relevanceScore > 0) {
                scoredContexts.push({
                    key,
                    ...context,
                    relevanceScore
                });
            }
        }

        // Sort by relevance and priority
        scoredContexts.sort((a, b) => {
            const scoreA = a.relevanceScore * a.priority;
            const scoreB = b.relevanceScore * b.priority;
            return scoreB - scoreA;
        });

        // Build context string within token limits
        return this.buildContextString(scoredContexts, maxTokens);
    }

    calculateRelevance(context, query) {
        let score = 0;
        const queryTerms = query.toLowerCase().split(/\s+/);
        const content = context.content.toLowerCase();

        // Term frequency scoring
        queryTerms.forEach(term => {
            if (content.includes(term)) {
                score += 1;
                // Bonus for exact matches
                if (content.includes(` ${term} `)) {
                    score += 2;
                }
            }
        });

        // Recency bonus
        const hoursOld = (Date.now() - context.timestamp) / (1000 * 60 * 60);
        if (hoursOld < 24) {
            score *= this.priorityWeights.recent;
        }

        // Access pattern bonus
        if (context.accessCount > 5) {
            score *= 1.2; // Frequently used context is probably important
        }

        return score;
    }

    buildContextString(contexts, maxTokens) {
        let result = "";
        let currentLength = 0;

        for (const context of contexts) {
            const contextString = `[${context.category}] ${context.content}\n`;
            
            if (currentLength + contextString.length <= maxTokens) {
                result += contextString;
                currentLength += contextString.length;
                
                // Update access count
                context.accessCount++;
                this.contexts.set(context.key, context);
            } else {
                break; // Respect token limits
            }
        }

        return result;
    }

    isExpired(context) {
        return context.expiresAt && Date.now() > context.expiresAt;
    }

    // Context pruning to prevent memory bloat
    pruneContexts(maxSize = 50) {
        if (this.contexts.size > maxSize) {
            const contextsArray = Array.from(this.contexts.entries());
            contextsArray.sort((a, b) => a[1].timestamp - b[1].timestamp); // Oldest first
            
            const toDelete = contextsArray.length - maxSize;
            for (let i = 0; i < toDelete; i++) {
                this.contexts.delete(contextsArray[i][0]);
            }
        }
    }
}
```

### Example 2: Specialized Context Strategies

```javascript
// Different context strategies for different types of tasks
class ContextStrategyManager {
    constructor() {
        this.strategies = {
            conversation: new ConversationContextStrategy(),
            research: new ResearchContextStrategy(),
            coding: new CodingContextStrategy(),
            analysis: new AnalysisContextStrategy()
        };
    }

    getContextForTask(taskType, query, additionalParams = {}) {
        const strategy = this.strategies[taskType];
        if (!strategy) {
            return this.strategies.conversation.getContext(query, additionalParams);
        }
        return strategy.getContext(query, additionalParams);
    }
}

class ConversationContextStrategy {
    getContext(query, params) {
        // Focus on recent conversation history and user preferences
        const contextPieces = [];
        
        // Recent messages (last 5 exchanges)
        if (params.conversationHistory) {
            const recentHistory = params.conversationHistory.slice(-10); // Last 5 exchanges
            contextPieces.push("Recent conversation:\n" + recentHistory.join("\n"));
        }

        // User preferences
        if (params.userPreferences) {
            contextPieces.push("User preferences: " + JSON.stringify(params.userPreferences));
        }

        // Current session goals
        if (params.sessionGoal) {
            contextPieces.push(`Session goal: ${params.sessionGoal}`);
        }

        return this.prioritizeAndTrim(contextPieces, 800);
    }

    prioritizeAndTrim(pieces, maxLength) {
        // Simple priority: conversation history > preferences > goals
        let result = "";
        for (const piece of pieces) {
            if (result.length + piece.length <= maxLength) {
                result += piece + "\n\n";
            }
        }
        return result;
    }
}

class CodingContextStrategy {
    getContext(query, params) {
        const contextPieces = [];
        
        // Code file context
        if (params.currentFile) {
            contextPieces.push(`Current file content:\n${params.currentFile}`);
        }

        // Related files
        if (params.relatedFiles) {
            contextPieces.push(`Related files:\n${params.relatedFiles.slice(0, 3).join('\n')}`);
        }

        // Project structure
        if (params.projectStructure) {
            contextPieces.push(`Project structure:\n${params.projectStructure}`);
        }

        // Recent errors
        if (params.recentErrors) {
            contextPieces.push(`Recent errors:\n${params.recentErrors.join('\n')}`);
        }

        return this.prioritizeCodeContext(contextPieces, 1000);
    }

    prioritizeCodeContext(pieces, maxLength) {
        // For coding, current file is most important, then errors, then related files
        const priorityOrder = [0, 3, 1, 2]; // indices based on importance
        let result = "";
        
        for (const index of priorityOrder) {
            if (pieces[index] && result.length + pieces[index].length <= maxLength) {
                result += pieces[index] + "\n\n";
            }
        }
        
        return result;
    }
}
```

### Example 3: Advanced Context Engineering with RAG (Retrieval-Augmented Generation)

```javascript
class AdvancedContextEngineer {
    constructor() {
        this.embeddingCache = new Map();
        this.semanticWeights = {
            conceptualSimilarity: 1.0,
            temporalRelevance: 0.8,
            userSpecific: 1.2,
            taskAlignment: 1.5
        };
    }

    async engineerContext(query, availableContexts, options = {}) {
        const {
            strategy = 'semantic',
            maxContextLength = 1500,
            diversityThreshold = 0.3
        } = options;

        switch (strategy) {
            case 'semantic':
                return await this.semanticRetrieval(query, availableContexts, maxContextLength);
            case 'temporal':
                return this.temporalRetrieval(availableContexts, maxContextLength);
            case 'hybrid':
                return await this.hybridRetrieval(query, availableContexts, maxContextLength, diversityThreshold);
            default:
                return await this.semanticRetrieval(query, availableContexts, maxContextLength);
        }
    }

    async semanticRetrieval(query, contexts, maxLength) {
        // Simple semantic similarity (in reality, you'd use proper embeddings)
        const queryEmbedding = this.simpleEmbed(query);
        const scoredContexts = contexts.map(context => {
            const contextEmbedding = this.simpleEmbed(context.content);
            const similarity = this.cosineSimilarity(queryEmbedding, contextEmbedding);
            
            return {
                ...context,
                score: similarity * this.semanticWeights.conceptualSimilarity
            };
        });

        return this.assembleContext(scoredContexts, maxLength);
    }

    temporalRetrieval(contexts, maxLength) {
        const scoredContexts = contexts.map(context => {
            const hoursOld = (Date.now() - context.timestamp) / (1000 * 60 * 60);
            const recencyScore = Math.max(0, 1 - (hoursOld / 168)); // Decay over 1 week
            
            return {
                ...context,
                score: recencyScore * this.semanticWeights.temporalRelevance
            };
        });

        return this.assembleContext(scoredContexts, maxLength);
    }

    async hybridRetrieval(query, contexts, maxLength, diversityThreshold) {
        const semanticResults = await this.semanticRetrieval(query, contexts, maxLength * 2);
        const temporalResults = this.temporalRetrieval(contexts, maxLength * 2);
        
        // Combine and diversify
        const combined = [...semanticResults.contexts, ...temporalResults.contexts];
        const diversified = this.diversifySelection(combined, diversityThreshold);
        
        return this.assembleContext(diversified, maxLength);
    }

    diversifySelection(contexts, threshold) {
        const diversified = [];
        
        for (const context of contexts) {
            let tooSimilar = false;
            
            for (const selected of diversified) {
                const similarity = this.cosineSimilarity(
                    this.simpleEmbed(context.content),
                    this.simpleEmbed(selected.content)
                );
                
                if (similarity > threshold) {
                    tooSimilar = true;
                    break;
                }
            }
            
            if (!tooSimilar) {
                diversified.push(context);
            }
        }
        
        return diversified;
    }

    assembleContext(scoredContexts, maxLength) {
        scoredContexts.sort((a, b) => b.score - a.score);
        
        let assembled = "";
        let currentLength = 0;
        const selectedContexts = [];

        for (const context of scoredContexts) {
            const contextString = this.formatContext(context);
            
            if (currentLength + contextString.length <= maxLength) {
                assembled += contextString + "\n\n";
                currentLength += contextString.length;
                selectedContexts.push(context);
            } else {
                break;
            }
        }

        return {
            context: assembled,
            sources: selectedContexts,
            totalLength: currentLength,
            effectivenessScore: this.calculateEffectiveness(selectedContexts)
        };
    }

    // Simple embedding function (replace with real embeddings in production)
    simpleEmbed(text) {
        if (this.embeddingCache.has(text)) {
            return this.embeddingCache.get(text);
        }

        const words = text.toLowerCase().split(/\W+/).filter(word => word.length > 2);
        const embedding = {};
        
        words.forEach(word => {
            embedding[word] = (embedding[word] || 0) + 1;
        });

        this.embeddingCache.set(text, embedding);
        return embedding;
    }

    cosineSimilarity(vecA, vecB) {
        const words = new Set([...Object.keys(vecA), ...Object.keys(vecB)]);
        let dotProduct = 0;
        let normA = 0;
        let normB = 0;

        for (const word of words) {
            const a = vecA[word] || 0;
            const b = vecB[word] || 0;
            dotProduct += a * b;
            normA += a * a;
            normB += b * b;
        }

        return dotProduct / (Math.sqrt(normA) * Math.sqrt(normB));
    }

    formatContext(context) {
        return `[${context.category.toUpperCase()}] ${context.content} (relevance: ${context.score.toFixed(2)})`;
    }

    calculateEffectiveness(contexts) {
        if (contexts.length === 0) return 0;
        
        const avgScore = contexts.reduce((sum, ctx) => sum + ctx.score, 0) / contexts.length;
        const diversity = this.calculateDiversity(contexts);
        
        return avgScore * diversity;
    }

    calculateDiversity(contexts) {
        if (contexts.length <= 1) return 1;
        
        let totalSimilarity = 0;
        let comparisons = 0;

        for (let i = 0; i < contexts.length; i++) {
            for (let j = i + 1; j < contexts.length; j++) {
                const similarity = this.cosineSimilarity(
                    this.simpleEmbed(contexts[i].content),
                    this.simpleEmbed(contexts[j].content)
                );
                totalSimilarity += similarity;
                comparisons++;
            }
        }

        return 1 - (totalSimilarity / comparisons);
    }
}
```

### Example 4: Practical Usage Example

```javascript
// Practical example of context engineering in action
class SmartAssistant {
    constructor() {
        this.contextEngineer = new AdvancedContextEngineer();
        this.contextStrategy = new ContextStrategyManager();
        this.conversationMemory = [];
    }

    async generateResponse(userMessage, userData = {}) {
        // Step 1: Determine the task type
        const taskType = this.classifyTask(userMessage);
        
        // Step 2: Engineer the optimal context
        const engineeredContext = await this.engineerOptimalContext(userMessage, taskType, userData);
        
        // Step 3: Generate response using the engineered context
        return await this.generateWithContext(userMessage, engineeredContext, taskType);
    }

    classifyTask(message) {
        const lowerMsg = message.toLowerCase();
        
        if (lowerMsg.includes('code') || lowerMsg.includes('program') || lowerMsg.includes('function')) {
            return 'coding';
        } else if (lowerMsg.includes('research') || lowerMsg.includes('find') || lowerMsg.includes('information')) {
            return 'research';
        } else if (lowerMsg.includes('analyze') || lowerMsg.includes('compare') || lowerMsg.includes('statistics')) {
            return 'analysis';
        } else {
            return 'conversation';
        }
    }

    async engineerOptimalContext(userMessage, taskType, userData) {
        // Get available contexts from various sources
        const availableContexts = await this.gatherAvailableContexts(userData);
        
        // Engineer the context based on task type and message
        const engineered = await this.contextEngineer.engineerContext(
            userMessage, 
            availableContexts, 
            {
                strategy: taskType === 'coding' ? 'hybrid' : 'semantic',
                maxContextLength: this.getOptimalLength(taskType)
            }
        );

        console.log(`Context Engineering Results:
- Task Type: ${taskType}
- Context Length: ${engineered.totalLength} chars
- Effectiveness Score: ${engineered.effectivenessScore.toFixed(2)}
- Sources Used: ${engineered.sources.length}`);

        return engineered.context;
    }

    async gatherAvailableContexts(userData) {
        // Simulate gathering context from different sources
        return [
            {
                content: "User prefers detailed explanations with examples",
                category: "preference",
                timestamp: Date.now() - 1000 * 60 * 30, // 30 minutes ago
                priority: 1.2
            },
            {
                content: "Current project: E-commerce platform using React and Node.js",
                category: "project",
                timestamp: Date.now() - 1000 * 60 * 60 * 2, // 2 hours ago
                priority: 1.5
            },
            {
                content: "Recent error: 'TypeError: Cannot read properties of undefined' in shopping cart",
                category: "error",
                timestamp: Date.now() - 1000 * 60 * 15, // 15 minutes ago
                priority: 2.0
            },
            // ... more contexts
        ];
    }

    getOptimalLength(taskType) {
        const lengthMap = {
            coding: 1200,
            research: 1500,
            analysis: 1000,
            conversation: 800
        };
        return lengthMap[taskType] || 1000;
    }

    async generateWithContext(userMessage, context, taskType) {
        // Simulate AI response generation
        const prompt = `
TASK TYPE: ${taskType.toUpperCase()}
RELEVANT CONTEXT:
${context}

USER MESSAGE: ${userMessage}

RESPONSE:`;
        
        // In reality, this would call an AI API
        return this.simulateAIResponse(prompt);
    }

    simulateAIResponse(prompt) {
        // Simulate different responses based on context quality
        if (prompt.includes('coding') && prompt.includes('error')) {
            return "Based on the error context you provided, this seems like a null reference issue. Let me help you debug the shopping cart component...";
        }
        return "I've analyzed your request considering the relevant context. Here's my response...";
    }
}

// Usage Example
const assistant = new SmartAssistant();

async function demo() {
    const response = await assistant.generateResponse(
        "I'm getting a TypeError in my shopping cart code", 
        { userId: "123", project: "ecommerce" }
    );
    console.log("Assistant Response:", response);
}

demo();
```

## Key Context Engineering Techniques

1. **Context Window Optimization**: Smartly fill the available token budget
2. **Semantic Chunking**: Break information into meaningful pieces
3. **Relevance Scoring**: Weight context by usefulness
4. **Temporal Decay**: Prioritize recent information
5. **Diversity Sampling**: Avoid redundant context
6. **Strategy Selection**: Use different approaches for different tasks

## Benefits of Good Context Engineering

- **Better AI Performance**: More relevant context = better responses
- **Cost Efficiency**: Fewer tokens wasted on irrelevant information
- **Faster Responses**: Reduced processing time
- **Improved Accuracy**: Less hallucination, more factual responses
- **Enhanced User Experience**: More personalized and relevant interactions

Context Engineering turns raw information into strategic advantage by applying engineering principles to context management.