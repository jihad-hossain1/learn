# Pipeline Pattern (Sequential Processing)

The **Pipeline Pattern** processes data through a series of sequential stages, where each stage transforms the output of the previous stage. This is ideal for data processing, ETL (Extract, Transform, Load) workflows, and any sequential transformation tasks.

## Key Concepts

- **Sequential Stages**: Data flows through stages in a fixed order
- **Stage Specialization**: Each stage has a specific responsibility
- **Data Transformation**: Each stage modifies the data
- **Stream Processing**: Can handle data streams efficiently
- **Quality Gates**: Validation at each stage

## JavaScript Implementation Example

### Basic Data Processing Pipeline

```javascript
// Base Pipeline Stage Class
class PipelineStage {
  constructor(name) {
    this.name = name;
    this.nextStage = null;
  }

  setNextStage(stage) {
    this.nextStage = stage;
    return stage; // Allow chaining
  }

  async process(data, context = {}) {
    console.log(`\n🔄 ${this.name} processing...`);
    console.log(`   Input: ${JSON.stringify(data).substring(0, 100)}...`);
    
    const startTime = Date.now();
    const result = await this.execute(data, context);
    const processingTime = Date.now() - startTime;
    
    console.log(`   Output: ${JSON.stringify(result).substring(0, 100)}...`);
    console.log(`   ⏱️  Processing time: ${processingTime}ms`);
    
    if (this.nextStage) {
      return await this.nextStage.process(result, context);
    }
    
    return result;
  }

  async execute(data, context) {
    // To be implemented by subclasses
    return data;
  }

  async delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

// Concrete Pipeline Stages
class DataExtractor extends PipelineStage {
  constructor() {
    super("Data Extractor");
  }

  async execute(data, context) {
    await this.delay(300);
    
    // Simulate data extraction from various sources
    const extractedData = {
      rawText: "Customer feedback: The product is great but delivery was slow. Rating: 4/5",
      metadata: {
        source: "customer_feedback",
        timestamp: new Date().toISOString(),
        format: "text"
      },
      additionalInfo: {
        customerId: "CUST-12345",
        product: "Smartphone X"
      }
    };
    
    context.originalSource = data.source || "unknown";
    return extractedData;
  }
}

class TextCleaner extends PipelineStage {
  constructor() {
    super("Text Cleaner");
  }

  async execute(data, context) {
    await this.delay(200);
    
    // Clean and normalize text
    const cleanedText = data.rawText
      .replace(/[^\w\s.,!?]/g, '') // Remove special characters
      .replace(/\s+/g, ' ')        // Normalize whitespace
      .trim()
      .toLowerCase();
    
    return {
      ...data,
      cleanedText,
      cleaningStats: {
        originalLength: data.rawText.length,
        cleanedLength: cleanedText.length,
        charactersRemoved: data.rawText.length - cleanedText.length
      }
    };
  }
}

class SentimentAnalyzer extends PipelineStage {
  constructor() {
    super("Sentiment Analyzer");
  }

  async execute(data, context) {
    await this.delay(500);
    
    // Simple sentiment analysis
    const positiveWords = ['great', 'excellent', 'good', 'awesome', 'amazing'];
    const negativeWords = ['slow', 'bad', 'terrible', 'poor', 'disappointing'];
    
    const text = data.cleanedText;
    let positiveScore = 0;
    let negativeScore = 0;
    
    positiveWords.forEach(word => {
      if (text.includes(word)) positiveScore++;
    });
    
    negativeWords.forEach(word => {
      if (text.includes(word)) negativeScore++;
    });
    
    const sentimentScore = positiveScore - negativeScore;
    const sentiment = sentimentScore > 0 ? 'positive' : sentimentScore < 0 ? 'negative' : 'neutral';
    
    return {
      ...data,
      sentimentAnalysis: {
        sentiment,
        score: sentimentScore,
        positiveWordsFound: positiveWords.filter(word => text.includes(word)),
        negativeWordsFound: negativeWords.filter(word => text.includes(word)),
        confidence: Math.min(Math.abs(sentimentScore) / 3, 1.0)
      }
    };
  }
}

class EntityExtractor extends PipelineStage {
  constructor() {
    super("Entity Extractor");
  }

  async execute(data, context) {
    await this.delay(400);
    
    // Extract entities using simple pattern matching
    const text = data.cleanedText;
    const entities = {
      products: this.extractProducts(text),
      ratings: this.extractRatings(text),
      issues: this.extractIssues(text),
      features: this.extractFeatures(text)
    };
    
    return {
      ...data,
      entities,
      extractionStats: {
        totalEntities: Object.values(entities).flat().length
      }
    };
  }

  extractProducts(text) {
    const productKeywords = ['product', 'delivery', 'service', 'phone', 'device'];
    return productKeywords.filter(keyword => text.includes(keyword));
  }

  extractRatings(text) {
    const ratingMatch = text.match(/(\d+)\/\s*(\d+)/);
    return ratingMatch ? [`${ratingMatch[1]}/${ratingMatch[2]}`] : [];
  }

  extractIssues(text) {
    const issueKeywords = ['slow', 'problem', 'issue', 'late'];
    return issueKeywords.filter(keyword => text.includes(keyword));
  }

  extractFeatures(text) {
    const featureKeywords = ['great', 'good', 'excellent', 'fast'];
    return featureKeywords.filter(keyword => text.includes(keyword));
  }
}

class DataEnricher extends PipelineStage {
  constructor() {
    super("Data Enricher");
  }

  async execute(data, context) {
    await this.delay(300);
    
    // Enrich data with additional context
    const enrichedData = {
      ...data,
      enrichedInfo: {
        processingTimestamp: new Date().toISOString(),
        pipelineVersion: "1.0",
        customerSegment: this.determineCustomerSegment(data.additionalInfo.customerId),
        priority: this.determinePriority(data.sentimentAnalysis.score),
        suggestedActions: this.generateSuggestions(data)
      }
    };
    
    return enrichedData;
  }

  determineCustomerSegment(customerId) {
    const segments = ['new', 'regular', 'vip', 'enterprise'];
    return segments[Math.floor(Math.random() * segments.length)];
  }

  determinePriority(sentimentScore) {
    if (sentimentScore < -1) return 'high';
    if (sentimentScore > 1) return 'low';
    return 'medium';
  }

  generateSuggestions(data) {
    const suggestions = [];
    
    if (data.sentimentAnalysis.sentiment === 'negative') {
      suggestions.push('Follow up with customer service');
    }
    
    if (data.entities.issues.length > 0) {
      suggestions.push('Review operational processes');
    }
    
    if (data.sentimentAnalysis.sentiment === 'positive') {
      suggestions.push('Consider for testimonial');
    }
    
    return suggestions.length > 0 ? suggestions : ['No action required'];
  }
}

class QualityValidator extends PipelineStage {
  constructor() {
    super("Quality Validator");
  }

  async execute(data, context) {
    await this.delay(250);
    
    // Validate data quality
    const validationResults = {
      hasText: !!data.cleanedText && data.cleanedText.length > 0,
      hasSentiment: !!data.sentimentAnalysis,
      hasEntities: data.entities && data.extractionStats.totalEntities > 0,
      textLengthValid: data.cleanedText.length >= 10,
      sentimentConfidence: data.sentimentAnalysis.confidence > 0.3
    };
    
    const isValid = Object.values(validationResults).every(result => result);
    
    if (!isValid) {
      console.warn(`   ⚠️  Validation failed:`, validationResults);
      throw new Error(`Data quality validation failed at ${this.name}`);
    }
    
    console.log(`   ✅ Quality validation passed`);
    
    return {
      ...data,
      validation: {
        isValid,
        results: validationResults,
        timestamp: new Date().toISOString()
      }
    };
  }
}

// Pipeline Builder
class DataProcessingPipeline {
  constructor() {
    this.stages = [];
  }

  addStage(stage) {
    this.stages.push(stage);
    return this; // Allow method chaining
  }

  build() {
    // Connect stages in sequence
    for (let i = 0; i < this.stages.length - 1; i++) {
      this.stages[i].setNextStage(this.stages[i + 1]);
    }
    
    return this.stages[0]; // Return the first stage
  }

  async process(data, context = {}) {
    const pipeline = this.build();
    console.log(`🚀 Starting pipeline with ${this.stages.length} stages`);
    console.log("=" .repeat(60));
    
    const startTime = Date.now();
    try {
      const result = await pipeline.process(data, context);
      const totalTime = Date.now() - startTime;
      
      console.log(`\n🎉 Pipeline completed successfully!`);
      console.log(`⏱️  Total processing time: ${totalTime}ms`);
      console.log(`📊 Stages executed: ${this.stages.length}`);
      
      return result;
    } catch (error) {
      console.error(`\n💥 Pipeline failed: ${error.message}`);
      throw error;
    }
  }
}

// Usage Example
async function demonstratePipeline() {
  // Build the pipeline
  const pipeline = new DataProcessingPipeline();
  
  pipeline
    .addStage(new DataExtractor())
    .addStage(new TextCleaner())
    .addStage(new SentimentAnalyzer())
    .addStage(new EntityExtractor())
    .addStage(new DataEnricher())
    .addStage(new QualityValidator());

  // Process sample data
  const sampleData = {
    source: "customer_feedback_system",
    type: "text_feedback",
    id: "FBK-2024-001"
  };

  try {
    const result = await pipeline.process(sampleData);
    
    console.log("\n📋 Final Result Summary:");
    console.log("=" .repeat(40));
    console.log(`Sentiment: ${result.sentimentAnalysis.sentiment} (score: ${result.sentimentAnalysis.score})`);
    console.log(`Entities found: ${result.extractionStats.totalEntities}`);
    console.log(`Suggested actions: ${result.enrichedInfo.suggestedActions.join(', ')}`);
    console.log(`Priority: ${result.enrichedInfo.priority}`);
    
  } catch (error) {
    console.error("Pipeline execution failed:", error);
  }
}

demonstratePipeline();
```

## Advanced Example: Configurable Pipeline with Error Handling

```javascript
class AdvancedPipelineStage extends PipelineStage {
  constructor(name, config = {}) {
    super(name);
    this.config = {
      timeout: config.timeout || 5000,
      retryAttempts: config.retryAttempts || 3,
      enabled: config.enabled !== false,
      ...config
    };
    this.metrics = {
      processedItems: 0,
      totalProcessingTime: 0,
      errors: 0
    };
  }

  async process(data, context) {
    if (!this.config.enabled) {
      console.log(`   ⏭️  ${this.name} skipped (disabled)`);
      return await this.nextStage?.process(data, context) || data;
    }

    this.metrics.processedItems++;
    const startTime = Date.now();

    try {
      const result = await this.executeWithRetry(data, context);
      const processingTime = Date.now() - startTime;
      this.metrics.totalProcessingTime += processingTime;
      
      return await this.nextStage?.process(result, context) || result;
    } catch (error) {
      this.metrics.errors++;
      throw error;
    }
  }

  async executeWithRetry(data, context) {
    for (let attempt = 1; attempt <= this.config.retryAttempts; attempt++) {
      try {
        return await Promise.race([
          this.execute(data, context),
          new Promise((_, reject) => 
            setTimeout(() => reject(new Error(`Timeout after ${this.config.timeout}ms`)), this.config.timeout)
          )
        ]);
      } catch (error) {
        if (attempt === this.config.retryAttempts) {
          throw new Error(`Failed after ${attempt} attempts: ${error.message}`);
        }
        console.log(`   🔄 ${this.name} attempt ${attempt} failed, retrying...`);
        await this.delay(1000 * attempt); // Exponential backoff
      }
    }
  }

  getMetrics() {
    const avgTime = this.metrics.processedItems > 0 
      ? this.metrics.totalProcessingTime / this.metrics.processedItems 
      : 0;
    
    return {
      stage: this.name,
      processedItems: this.metrics.processedItems,
      averageProcessingTime: avgTime,
      errorRate: this.metrics.processedItems > 0 
        ? this.metrics.errors / this.metrics.processedItems 
        : 0
    };
  }
}

// Specialized Advanced Stages
class MLModelStage extends AdvancedPipelineStage {
  constructor(modelName, config = {}) {
    super(`ML Model: ${modelName}`, config);
    this.modelName = modelName;
  }

  async execute(data, context) {
    console.log(`   🤖 Running ${this.modelName}...`);
    await this.delay(800);
    
    // Simulate ML model processing
    const confidence = 0.7 + Math.random() * 0.3;
    
    return {
      ...data,
      mlResults: {
        ...data.mlResults,
        [this.modelName]: {
          prediction: this.generatePrediction(data),
          confidence: confidence,
          featuresUsed: this.extractFeatures(data),
          timestamp: new Date().toISOString()
        }
      }
    };
  }

  generatePrediction(data) {
    const predictions = ['positive', 'negative', 'neutral', 'urgent', 'normal'];
    return predictions[Math.floor(Math.random() * predictions.length)];
  }

  extractFeatures(data) {
    return ['text_length', 'sentiment_score', 'entity_count', 'word_diversity'];
  }
}

class DatabaseWriter extends AdvancedPipelineStage {
  constructor(tableName, config = {}) {
    super(`Database Writer: ${tableName}`, config);
    this.tableName = tableName;
  }

  async execute(data, context) {
    console.log(`   💾 Writing to ${this.tableName}...`);
    await this.delay(600);
    
    // Simulate database write
    const recordId = `rec_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    return {
      ...data,
      databaseRecords: {
        ...data.databaseRecords,
        [this.tableName]: {
          id: recordId,
          writtenAt: new Date().toISOString(),
          recordSize: JSON.stringify(data).length,
          status: 'success'
        }
      }
    };
  }
}

class APIExportStage extends AdvancedPipelineStage {
  constructor(apiEndpoint, config = {}) {
    super(`API Export: ${apiEndpoint}`, config);
    this.apiEndpoint = apiEndpoint;
  }

  async execute(data, context) {
    console.log(`   🌐 Exporting to ${this.apiEndpoint}...`);
    await this.delay(700);
    
    // Simulate API call
    const success = Math.random() > 0.1; // 90% success rate
    
    if (!success) {
      throw new Error(`API export failed for ${this.apiEndpoint}`);
    }
    
    return {
      ...data,
      apiExports: {
        ...data.apiExports,
        [this.apiEndpoint]: {
          exportedAt: new Date().toISOString(),
          status: 'success',
          dataPoints: Object.keys(data).length
        }
      }
    };
  }
}

// Advanced Pipeline Manager
class AdvancedPipelineManager {
  constructor() {
    this.pipelines = new Map();
    this.metricsCollector = new MetricsCollector();
  }

  createPipeline(name, stages) {
    const pipeline = new DataProcessingPipeline();
    stages.forEach(stage => pipeline.addStage(stage));
    
    this.pipelines.set(name, pipeline);
    console.log(`📋 Created pipeline: ${name} with ${stages.length} stages`);
    
    return pipeline;
  }

  async runPipeline(name, data) {
    const pipeline = this.pipelines.get(name);
    if (!pipeline) {
      throw new Error(`Pipeline '${name}' not found`);
    }

    console.log(`\n🚀 Executing pipeline: ${name}`);
    console.log("=" .repeat(50));
    
    const result = await pipeline.process(data);
    
    // Collect metrics
    this.metricsCollector.recordExecution(name, result);
    
    return result;
  }

  getPipelineMetrics(name) {
    return this.metricsCollector.getMetrics(name);
  }

  listPipelines() {
    return Array.from(this.pipelines.keys());
  }
}

class MetricsCollector {
  constructor() {
    this.executions = new Map();
  }

  recordExecution(pipelineName, result) {
    if (!this.executions.has(pipelineName)) {
      this.executions.set(pipelineName, []);
    }
    
    this.executions.get(pipelineName).push({
      timestamp: new Date().toISOString(),
      resultSummary: this.createSummary(result)
    });
  }

  createSummary(result) {
    return {
      dataSize: JSON.stringify(result).length,
      processingStages: Object.keys(result).length,
      hasErrors: !!result.error,
      timestamp: new Date().toISOString()
    };
  }

  getMetrics(pipelineName) {
    const executions = this.executions.get(pipelineName) || [];
    return {
      totalExecutions: executions.length,
      averageDataSize: executions.reduce((sum, e) => sum + e.resultSummary.dataSize, 0) / executions.length,
      lastExecution: executions[executions.length - 1]
    };
  }
}

// Advanced Usage Example
async function advancedPipelineDemo() {
  const manager = new AdvancedPipelineManager();

  // Create different pipelines for different use cases
  const feedbackPipeline = manager.createPipeline("customer_feedback", [
    new DataExtractor(),
    new TextCleaner(),
    new MLModelStage("SentimentClassifier", { retryAttempts: 2 }),
    new MLModelStage("TopicExtractor"),
    new DatabaseWriter("feedback_analysis"),
    new APIExportStage("https://api.crm.com/feedback")
  ]);

  const dataEnrichmentPipeline = manager.createPipeline("data_enrichment", [
    new DataExtractor(),
    new EntityExtractor(),
    new DataEnricher(),
    new DatabaseWriter("enriched_data"),
    new APIExportStage("https://api.analytics.com/data")
  ]);

  // Process multiple items
  const sampleData = [
    { source: "feedback_system", id: "FBK-001", text: "Great product, fast shipping!" },
    { source: "feedback_system", id: "FBK-002", text: "Slow delivery and poor packaging" },
    { source: "feedback_system", id: "FBK-003", text: "Average product, good customer service" }
  ];

  for (const data of sampleData) {
    console.log(`\n📦 Processing: ${data.id}`);
    try {
      const result = await manager.runPipeline("customer_feedback", data);
      console.log(`✅ Completed: ${data.id}`);
    } catch (error) {
      console.error(`❌ Failed: ${data.id} - ${error.message}`);
    }
  }

  // Show metrics
  console.log("\n📊 Pipeline Metrics:");
  console.log("=" .repeat(30));
  manager.listPipelines().forEach(name => {
    const metrics = manager.getPipelineMetrics(name);
    console.log(`${name}: ${metrics.totalExecutions} executions`);
  });
}

advancedPipelineDemo();
```

## Key Benefits of Pipeline Pattern

1. **Modularity**: Each stage is independent and testable
2. **Reusability**: Stages can be reused across different pipelines
3. **Maintainability**: Easy to modify or extend individual stages
4. **Monitoring**: Each stage can be monitored independently
5. **Scalability**: Stages can be parallelized or distributed
6. **Quality Control**: Validation at each stage

## Common Use Cases

- **Data Processing and ETL Pipelines**
- **Image and Video Processing**
- **Natural Language Processing**
- **CI/CD Deployment Pipelines**
- **Financial Transaction Processing**
- **Scientific Data Analysis**

This pattern is excellent for workflows where data undergoes multiple transformations and each step has clear responsibilities and quality requirements.

```javascript
// Quick simple example for clarity
const simplePipeline = new DataProcessingPipeline()
  .addStage(new DataExtractor())
  .addStage(new TextCleaner())
  .addStage(new SentimentAnalyzer());

// Usage
simplePipeline.process({source: "test"})
  .then(result => console.log("Simple pipeline result:", result))
  .catch(error => console.error("Error:", error));
```