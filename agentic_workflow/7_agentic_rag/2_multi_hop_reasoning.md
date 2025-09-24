# Multi-hop Reasoning in AI Systems

Multi-hop reasoning is the process of answering complex questions that require connecting multiple pieces of information across different documents or knowledge sources through several logical steps ("hops").

## What is Multi-hop Reasoning?

**Multi-hop reasoning** involves making inferences that require connecting multiple facts or pieces of evidence through sequential logical steps. Unlike single-hop questions that can be answered directly from one source, multi-hop questions require:

- **Multiple retrieval steps**
- **Intermediate reasoning**
- **Information synthesis**
- **Logical inference chains**

## Types of Multi-hop Reasoning

### 1. **Sequential Reasoning**
```
Q: "Where was the inventor of the telephone born?"
Hops: [Inventor of telephone → Alexander Graham Bell → Birth place → Edinburgh, Scotland]
```

### 2. **Comparative Reasoning**
```
Q: "Which is older: Python or JavaScript?"
Hops: [Python release year → 1991] → [JavaScript release year → 1995] → [Compare dates]
```

### 3. **Set Operations Reasoning**
```
Q: "What programming languages were created before 1990 and are still popular today?"
Hops: [Languages before 1990] ∩ [Popular languages today]
```

## Multi-hop Reasoning Architectures

### 1. **Retrieve-Then-Read Approach**
```javascript
class MultiHopRetrieveThenRead {
    constructor(retriever, reader) {
        this.retriever = retriever;
        this.reader = reader;
    }

    async answerQuestion(question) {
        // Step 1: Initial retrieval
        let context = await this.retriever.retrieve(question);
        let hopCount = 0;
        const maxHops = 3;

        while (hopCount < maxHops) {
            // Step 2: Analyze if we need more information
            const analysis = await this.analyzeCompleteness(question, context);
            
            if (analysis.isComplete) {
                break;
            }

            // Step 3: Generate follow-up questions for missing information
            const followUpQueries = await this.generateFollowUpQueries(question, context, analysis.missingInfo);
            
            // Step 4: Retrieve additional information
            for (const query of followUpQueries) {
                const additionalContext = await this.retriever.retrieve(query);
                context = this.mergeContexts(context, additionalContext);
            }
            
            hopCount++;
        }

        // Final answer generation
        return await this.reader.generateAnswer(question, context);
    }

    async analyzeCompleteness(question, context) {
        // Use LLM to analyze if context contains all needed information
        const prompt = `
        Question: ${question}
        Context: ${JSON.stringify(context.slice(0, 500))}
        
        Analyze if the context contains sufficient information to answer the question completely.
        Identify any missing information needed.
        
        Return JSON: {isComplete: boolean, missingInfo: string[]}
        `;
        
        // In practice, call your LLM here
        return { isComplete: context.length > 5, missingInfo: [] };
    }
}
```

### 2. **Iterative Retrieval Approach**
```javascript
class IterativeMultiHopReasoner {
    constructor(vectorStore, llm) {
        this.vectorStore = vectorStore;
        this.llm = llm;
    }

    async performMultiHopReasoning(question) {
        const reasoningChain = [];
        let currentQuery = question;
        let currentContext = [];
        const maxHops = 3;

        for (let hop = 0; hop < maxHops; hop++) {
            console.log(`Hop ${hop + 1}: ${currentQuery}`);
            
            // Retrieve relevant documents
            const documents = await this.vectorStore.similaritySearch(currentQuery, 3);
            currentContext.push(...documents);
            
            // Analyze if we can answer or need another hop
            const analysis = await this.analyzeStep(question, currentQuery, currentContext);
            reasoningChain.push({
                hop: hop + 1,
                query: currentQuery,
                documents: documents,
                intermediateAnswer: analysis.intermediateAnswer
            });

            if (analysis.canAnswer) {
                break;
            }

            if (hop < maxHops - 1) {
                currentQuery = await this.generateNextQuery(question, currentContext, reasoningChain);
            }
        }

        return await this.synthesizeFinalAnswer(question, reasoningChain);
    }

    async analyzeStep(originalQuestion, currentQuery, context) {
        const prompt = `
        Original Question: ${originalQuestion}
        Current Focus: ${currentQuery}
        Context: ${context.map(d => d.content).join('\n')}
        
        Based on the context, can we fully answer the original question?
        If not, what intermediate information do we have?
        
        Return JSON: {
            canAnswer: boolean,
            intermediateAnswer: string,
            confidence: number
        }
        `;
        
        // Mock response - in practice, call LLM
        return {
            canAnswer: context.length >= 2,
            intermediateAnswer: "Found some relevant information",
            confidence: 0.7
        };
    }

    async generateNextQuery(originalQuestion, context, reasoningChain) {
        const prompt = `
        Original Question: ${originalQuestion}
        Reasoning Chain So Far: ${JSON.stringify(reasoningChain)}
        
        What specific piece of information should we search for next to help answer the question?
        Generate a precise search query.
        `;
        
        // Mock response
        return `More specific information about ${originalQuestion}`;
    }
}
```

## Advanced Multi-hop Reasoning with Graph-based Approach

```javascript
class GraphBasedMultiHopReasoner {
    constructor() {
        this.knowledgeGraph = new Map();
    }

    async buildReasoningGraph(question) {
        const graph = {
            nodes: new Map(),
            edges: new Map(),
            questions: new Set([question])
        };

        await this.expandGraph(graph, question, 0, 3);
        return graph;
    }

    async expandGraph(graph, currentQuestion, depth, maxDepth) {
        if (depth >= maxDepth) return;

        const entities = await this.extractEntities(currentQuestion);
        const relationships = await this.findRelationships(entities);

        // Add nodes and edges to graph
        for (const entity of entities) {
            graph.nodes.set(entity, await this.getEntityInfo(entity));
        }

        for (const rel of relationships) {
            graph.edges.set(`${rel.source}-${rel.target}`, rel);
            
            // Generate new questions based on relationships
            const newQuestion = await this.generateFollowUpQuestion(rel, graph.questions);
            if (newQuestion && !graph.questions.has(newQuestion)) {
                graph.questions.add(newQuestion);
                await this.expandGraph(graph, newQuestion, depth + 1, maxDepth);
            }
        }
    }

    async extractEntities(text) {
        // Use NER or LLM to extract entities
        const mockEntities = text.split(' ').filter(word => 
            word.length > 3 && /[A-Z]/.test(word[0]));
        return [...new Set(mockEntities)];
    }
}
```

## Complete Multi-hop RAG System

```javascript
class MultiHopRAGSystem {
    constructor(vectorStore, llm) {
        this.vectorStore = vectorStore;
        this.llm = llm;
        this.reasoningStrategies = {
            'sequential': this.sequentialReasoning,
            'comparative': this.comparativeReasoning,
            'set_operations': this.setOperationsReasoning
        };
    }

    async answerComplexQuestion(question) {
        // Determine reasoning type
        const reasoningType = await this.classifyReasoningType(question);
        console.log(`Detected reasoning type: ${reasoningType}`);
        
        // Execute appropriate reasoning strategy
        const strategy = this.reasoningStrategies[reasoningType];
        return await strategy.call(this, question);
    }

    async sequentialReasoning(question) {
        const hops = await this.decomposeSequentialQuestion(question);
        const intermediateResults = [];
        
        for (let i = 0; i < hops.length; i++) {
            const hopQuestion = hops[i];
            const documents = await this.vectorStore.similaritySearch(hopQuestion, 2);
            const answer = await this.answerSingleHop(hopQuestion, documents);
            
            intermediateResults.push({
                hop: i + 1,
                question: hopQuestion,
                answer: answer,
                documents: documents
            });

            // Use answer to inform next hop if needed
            if (i < hops.length - 1) {
                hops[i + 1] = await this.contextualizeNextHop(hops[i + 1], answer);
            }
        }
        
        return await this.synthesizeSequentialAnswer(question, intermediateResults);
    }

    async comparativeReasoning(question) {
        // Extract comparison entities
        const entities = await this.extractComparisonEntities(question);
        const comparisons = [];
        
        // Retrieve information for each entity
        for (const entity of entities) {
            const entityInfo = await this.vectorStore.similaritySearch(entity, 3);
            const summary = await this.summarizeEntity(entity, entityInfo);
            comparisons.push({ entity, summary, info: entityInfo });
        }
        
        // Perform comparison
        return await this.performComparison(question, comparisons);
    }

    async decomposeSequentialQuestion(question) {
        const prompt = `
        Question: "${question}"
        
        Decompose this into sequential sub-questions that need to be answered in order.
        Return as JSON array: ["sub-question 1", "sub-question 2", ...]
        `;
        
        // Mock decomposition
        if (question.includes("inventor of")) {
            return [
                "Who invented the telephone?",
                "Where was Alexander Graham Bell born?"
            ];
        }
        return [question];
    }

    async answerSingleHop(question, documents) {
        const context = documents.map(d => d.content).join('\n');
        const prompt = `
        Context: ${context}
        Question: ${question}
        
        Answer the question based only on the context provided.
        `;
        
        // In practice: return await this.llm.generate(prompt);
        return `Answer for: ${question}`;
    }

    async classifyReasoningType(question) {
        if (question.includes("compare") || question.includes("vs")) return 'comparative';
        if (question.includes("before") && question.includes("after")) return 'set_operations';
        if (question.includes("inventor of") || question.includes("founder of")) return 'sequential';
        return 'sequential';
    }
}

// Usage Example
async function demoMultiHopReasoning() {
    const ragSystem = new MultiHopRAGSystem();
    
    const questions = [
        "What programming language was created first: Python or JavaScript?",
        "Where was the inventor of the telephone born?",
        "Which companies founded before 2000 are currently worth over $1 trillion?"
    ];
    
    for (const question of questions) {
        console.log(`\nQuestion: ${question}`);
        const answer = await ragSystem.answerComplexQuestion(question);
        console.log(`Answer: ${answer}`);
    }
}
```

## Evaluation Metrics for Multi-hop Reasoning

```javascript
class MultiHopEvaluator {
    evaluateReasoningChain(question, groundTruth, reasoningChain) {
        const metrics = {
            hopAccuracy: this.calculateHopAccuracy(reasoningChain),
            informationCompleteness: this.calculateCompleteness(question, reasoningChain, groundTruth),
            logicalCoherence: this.assessLogicalFlow(reasoningChain),
            finalAnswerAccuracy: this.assessFinalAnswer(reasoningChain[reasoningChain.length - 1].answer, groundTruth)
        };
        
        return {
            ...metrics,
            overallScore: this.calculateOverallScore(metrics)
        };
    }

    calculateHopAccuracy(chain) {
        // Evaluate if each hop was necessary and productive
        let validHops = 0;
        for (let i = 1; i < chain.length; i++) {
            if (this.isHopProductive(chain[i-1], chain[i])) {
                validHops++;
            }
        }
        return validHops / (chain.length - 1);
    }
}
```

## Key Challenges in Multi-hop Reasoning

1. **Error Propagation**: Mistakes in early hops compound
2. **Information Overload**: Too many documents to process
3. **Query Formulation**: Generating effective follow-up queries
4. **Stopping Criteria**: Knowing when reasoning is complete
5. **Context Management**: Handling large context windows

## Best Practices

1. **Start with query decomposition**
2. **Implement early stopping mechanisms**
3. **Use confidence scoring for each hop**
4. **Maintain reasoning transparency**
5. **Validate intermediate results**

Multi-hop reasoning significantly enhances RAG systems' ability to handle complex questions that require connecting information from multiple sources through logical inference chains.