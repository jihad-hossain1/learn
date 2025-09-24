# Active Learning Loops

Active learning is a machine learning approach where the algorithm selectively queries humans for labels on the most informative data points, creating an interactive loop between the model and human experts.

## Core Components of Active Learning Loops

### 1. **Query Strategy** - How to select the most valuable samples
### 2. **Human-in-the-Loop** - Expert annotation interface
### 3. **Model Retraining** - Continuous learning from new labels
### 4. **Performance Monitoring** - Tracking improvement over time

## Basic Active Learning Loop Implementation

```javascript
class ActiveLearningLoop {
  constructor(model, queryStrategy, labeledData = [], unlabeledData = []) {
    this.model = model;
    this.queryStrategy = queryStrategy;
    this.labeledData = labeledData;
    this.unlabeledData = unlabeledData;
    this.iteration = 0;
    this.performanceHistory = [];
    this.humanLabeler = new HumanLabeler();
  }

  async runIteration(budget = 10) {
    console.log(`Starting active learning iteration ${this.iteration}`);
    
    // 1. Train model on current labeled data
    await this.model.train(this.labeledData);
    
    // 2. Evaluate current performance
    const performance = await this.evaluateModel();
    this.performanceHistory.push({
      iteration: this.iteration,
      performance,
      labeledDataSize: this.labeledData.length
    });
    
    // 3. Query most informative samples
    const queries = this.queryStrategy.selectQueries(
      this.unlabeledData, 
      budget, 
      this.model
    );
    
    // 4. Get human labels for queries
    const newlyLabeled = await this.humanLabeler.getLabels(queries);
    
    // 5. Update datasets
    this.updateDatasets(newlyLabeled);
    
    this.iteration++;
    return newlyLabeled;
  }

  async evaluateModel() {
    if (this.labeledData.length < 10) {
      return { accuracy: 0, confidence: 0 };
    }
    
    // Simple holdout validation
    const splitPoint = Math.floor(this.labeledData.length * 0.8);
    const trainData = this.labeledData.slice(0, splitPoint);
    const testData = this.labeledData.slice(splitPoint);
    
    const tempModel = this.model.clone();
    await tempModel.train(trainData);
    
    let correct = 0;
    let totalConfidence = 0;
    
    testData.forEach(sample => {
      const prediction = tempModel.predict(sample.features);
      if (prediction.label === sample.label) {
        correct++;
      }
      totalConfidence += prediction.confidence;
    });
    
    return {
      accuracy: correct / testData.length,
      confidence: totalConfidence / testData.length,
      testSize: testData.length
    };
  }

  updateDatasets(newlyLabeled) {
    // Move labeled samples from unlabeled to labeled dataset
    newlyLabeled.forEach(labeledSample => {
      const index = this.unlabeledData.findIndex(
        sample => sample.id === labeledSample.id
      );
      
      if (index !== -1) {
        this.unlabeledData.splice(index, 1);
        this.labeledData.push(labeledSample);
      }
    });
  }

  async runMultipleIterations(iterations = 5, budgetPerIteration = 10) {
    const results = [];
    
    for (let i = 0; i < iterations; i++) {
      if (this.unlabeledData.length === 0) {
        console.log('No more unlabeled data available');
        break;
      }
      
      const newlyLabeled = await this.runIteration(budgetPerIteration);
      results.push({
        iteration: i,
        newlyLabeled: newlyLabeled.length,
        performance: this.performanceHistory[this.performanceHistory.length - 1]
      });
      
      console.log(`Iteration ${i} completed: ${newlyLabeled.length} samples labeled`);
    }
    
    return results;
  }

  getPerformanceTrend() {
    return this.performanceHistory.map(entry => ({
      iteration: entry.iteration,
      accuracy: entry.performance.accuracy,
      dataSize: entry.labeledDataSize
    }));
  }
}
```

## Query Strategies

### 1. **Uncertainty Sampling**
```javascript
class UncertaintySampling {
  constructor(uncertaintyMeasure = 'entropy') {
    this.uncertaintyMeasure = uncertaintyMeasure;
  }

  selectQueries(unlabeledData, budget, model) {
    // Score each sample by uncertainty
    const scoredSamples = unlabeledData.map(sample => {
      const prediction = model.predict(sample.features);
      const uncertainty = this.calculateUncertainty(prediction);
      
      return {
        sample,
        uncertainty,
        prediction
      };
    });

    // Sort by uncertainty (highest first) and select top budget
    return scoredSamples
      .sort((a, b) => b.uncertainty - a.uncertainty)
      .slice(0, budget)
      .map(item => item.sample);
  }

  calculateUncertainty(prediction) {
    switch (this.uncertaintyMeasure) {
      case 'entropy':
        return this.calculateEntropy(prediction.probabilities);
      
      case 'least_confidence':
        return 1 - Math.max(...Object.values(prediction.probabilities));
      
      case 'margin':
        const sorted = Object.values(prediction.probabilities).sort((a, b) => b - a);
        return sorted[0] - sorted[1]; // Difference between top two probabilities
      
      default:
        return 1 - prediction.confidence;
    }
  }

  calculateEntropy(probabilities) {
    let entropy = 0;
    Object.values(probabilities).forEach(prob => {
      if (prob > 0) {
        entropy -= prob * Math.log(prob);
      }
    });
    return entropy;
  }
}
```

### 2. **Diversity Sampling**
```javascript
class DiversitySampling {
  constructor(distanceMetric = 'euclidean') {
    this.distanceMetric = distanceMetric;
  }

  selectQueries(unlabeledData, budget, model, existingData = []) {
    // Select diverse samples that are different from existing labeled data
    const selectedSamples = [];
    const remainingSamples = [...unlabeledData];

    // First sample: highest uncertainty
    const firstSample = this.getMostUncertain(remainingSamples, model);
    selectedSamples.push(firstSample);
    this.removeSample(remainingSamples, firstSample);

    // Subsequent samples: maximize diversity
    while (selectedSamples.length < budget && remainingSamples.length > 0) {
      const nextSample = this.getMostDiverse(
        remainingSamples, 
        selectedSamples, 
        existingData
      );
      selectedSamples.push(nextSample);
      this.removeSample(remainingSamples, nextSample);
    }

    return selectedSamples;
  }

  getMostUncertain(samples, model) {
    let maxUncertainty = -1;
    let mostUncertainSample = null;

    samples.forEach(sample => {
      const prediction = model.predict(sample.features);
      const uncertainty = 1 - prediction.confidence;
      
      if (uncertainty > maxUncertainty) {
        maxUncertainty = uncertainty;
        mostUncertainSample = sample;
      }
    });

    return mostUncertainSample;
  }

  getMostDiverse(candidates, selectedSamples, existingData) {
    let maxMinDistance = -1;
    let mostDiverseSample = null;

    candidates.forEach(candidate => {
      // Calculate minimum distance to any already selected sample
      const minDistance = this.calculateMinDistance(
        candidate, 
        [...selectedSamples, ...existingData]
      );

      if (minDistance > maxMinDistance) {
        maxMinDistance = minDistance;
        mostDiverseSample = candidate;
      }
    });

    return mostDiverseSample;
  }

  calculateMinDistance(sample, comparisonSet) {
    if (comparisonSet.length === 0) return Infinity;
    
    let minDistance = Infinity;
    
    comparisonSet.forEach(comparisonSample => {
      const distance = this.calculateDistance(sample.features, comparisonSample.features);
      if (distance < minDistance) {
        minDistance = distance;
      }
    });
    
    return minDistance;
  }

  calculateDistance(featuresA, featuresB) {
    // Euclidean distance
    const squaredDiff = Object.keys(featuresA).map(key => {
      const valA = featuresA[key] || 0;
      const valB = featuresB[key] || 0;
      return Math.pow(valA - valB, 2);
    });
    
    return Math.sqrt(squaredDiff.reduce((sum, val) => sum + val, 0));
  }

  removeSample(samples, sampleToRemove) {
    const index = samples.findIndex(s => s.id === sampleToRemove.id);
    if (index !== -1) {
      samples.splice(index, 1);
    }
  }
}
```

### 3. **Query-by-Committee**
```javascript
class QueryByCommittee {
  constructor(committeeSize = 3) {
    this.committeeSize = committeeSize;
    this.committee = [];
  }

  initializeCommittee(modelTemplate, labeledData) {
    this.committee = [];
    
    for (let i = 0; i < this.committeeSize; i++) {
      // Create committee member with different initializations
      const member = modelTemplate.clone();
      
      // Bootstrap sample from labeled data
      const bootstrapSample = this.bootstrapSample(labeledData);
      member.train(bootstrapSample);
      
      this.committee.push(member);
    }
  }

  selectQueries(unlabeledData, budget, model) {
    if (this.committee.length === 0) {
      this.initializeCommittee(model, []);
    }

    // Score by committee disagreement
    const scoredSamples = unlabeledData.map(sample => {
      const disagreements = this.calculateDisagreement(sample);
      return { sample, disagreements };
    });

    return scoredSamples
      .sort((a, b) => b.disagreements - a.disagreements)
      .slice(0, budget)
      .map(item => item.sample);
  }

  calculateDisagreement(sample) {
    const predictions = this.committee.map(member => 
      member.predict(sample.features).label
    );

    // Calculate voting entropy
    const voteCounts = {};
    predictions.forEach(prediction => {
      voteCounts[prediction] = (voteCounts[prediction] || 0) + 1;
    });

    let entropy = 0;
    Object.values(voteCounts).forEach(count => {
      const probability = count / this.committeeSize;
      entropy -= probability * Math.log(probability);
    });

    return entropy;
  }

  bootstrapSample(data, size = null) {
    const sampleSize = size || data.length;
    const sample = [];
    
    for (let i = 0; i < sampleSize; i++) {
      const randomIndex = Math.floor(Math.random() * data.length);
      sample.push(data[randomIndex]);
    }
    
    return sample;
  }

  updateCommittee(newlyLabeledData) {
    // Retrain committee members with new data
    this.committee.forEach(member => {
      const trainingData = [...this.bootstrapSample(newlyLabeledData), ...newlyLabeledData];
      member.train(trainingData);
    });
  }
}
```

## Human Labeling Interface

```javascript
class HumanLabeler {
  constructor() {
    this.labelingQueue = [];
    this.activeLabelers = new Map();
    this.labelingHistory = [];
  }

  async getLabels(samples, options = {}) {
    const {
      priority = 'normal',
      timeout = 300000, // 5 minutes
      requiredLabelers = 1
    } = options;

    const labelingTask = {
      id: this.generateTaskId(),
      samples: [...samples],
      priority,
      createdAt: new Date(),
      status: 'pending',
      labels: [],
      requiredLabelers
    };

    this.labelingQueue.push(labelingTask);
    this.sortQueueByPriority();

    // Wait for labeling completion
    return new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        reject(new Error('Labeling timeout'));
      }, timeout);

      // Poll for completion
      const checkInterval = setInterval(() => {
        if (labelingTask.status === 'completed') {
          clearTimeout(timeoutId);
          clearInterval(checkInterval);
          resolve(labelingTask.labels);
        }
      }, 1000);
    });
  }

  async presentLabelingInterface(sample, labelerId) {
    const interface = new LabelingInterface(sample, {
      availableLabels: this.getAvailableLabels(sample),
      context: this.getLabelingContext(sample),
      previousLabels: this.getPreviousLabels(sample)
    });

    return await interface.getLabel();
  }

  getAvailableLabels(sample) {
    // Define possible labels based on sample type
    if (sample.type === 'image') {
      return ['cat', 'dog', 'bird', 'other'];
    } else if (sample.type === 'text') {
      return ['positive', 'negative', 'neutral'];
    }
    return ['label1', 'label2', 'label3'];
  }

  processLabelingResult(taskId, labelerId, results) {
    const task = this.labelingQueue.find(t => t.id === taskId);
    if (!task) return;

    results.forEach(result => {
      task.labels.push({
        ...result,
        labelerId,
        timestamp: new Date(),
        taskId
      });
    });

    // Check if task is complete
    if (task.labels.length >= task.samples.length * task.requiredLabelers) {
      task.status = 'completed';
      this.recordLabelingHistory(task);
    }
  }

  recordLabelingHistory(task) {
    this.labelingHistory.push({
      taskId: task.id,
      samplesCount: task.samples.length,
      completedAt: new Date(),
      labels: task.labels
    });
  }

  sortQueueByPriority() {
    const priorityWeights = { 'high': 3, 'normal': 2, 'low': 1 };
    this.labelingQueue.sort((a, b) => 
      priorityWeights[b.priority] - priorityWeights[a.priority]
    );
  }

  generateTaskId() {
    return `task_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }
}

class LabelingInterface {
  constructor(sample, options) {
    this.sample = sample;
    this.options = options;
    this.container = null;
  }

  async getLabel() {
    return new Promise((resolve) => {
      this.renderInterface();
      this.setupEventHandlers(resolve);
    });
  }

  renderInterface() {
    this.container = document.createElement('div');
    this.container.className = 'labeling-interface';
    
    this.container.innerHTML = `
      <div class="sample-display">
        ${this.renderSample()}
      </div>
      <div class="label-options">
        ${this.renderLabelOptions()}
      </div>
      <div class="confidence-slider">
        <label>Labeling Confidence:</label>
        <input type="range" min="1" max="5" value="3" class="confidence-input">
        <span class="confidence-value">3/5</span>
      </div>
      <div class="actions">
        <button class="btn-submit">Submit Label</button>
        <button class="btn-skip">Skip</button>
        <button class="btn-flag">Flag Issue</button>
      </div>
      <div class="context-info">
        ${this.renderContextInfo()}
      </div>
    `;

    document.body.appendChild(this.container);
  }

  renderSample() {
    if (this.sample.type === 'image') {
      return `<img src="${this.sample.data}" alt="Sample to label" style="max-width: 300px;">`;
    } else if (this.sample.type === 'text') {
      return `<div class="text-sample">${this.sample.data}</div>`;
    }
    return `<pre>${JSON.stringify(this.sample.data, null, 2)}</pre>`;
  }

  renderLabelOptions() {
    return this.options.availableLabels.map(label => `
      <label class="label-option">
        <input type="radio" name="label" value="${label}">
        ${label}
      </label>
    `).join('');
  }

  setupEventHandlers(resolve) {
    this.container.querySelector('.btn-submit').addEventListener('click', () => {
      const selectedLabel = this.container.querySelector('input[name="label"]:checked');
      const confidence = this.container.querySelector('.confidence-input').value;
      
      if (selectedLabel) {
        this.container.remove();
        resolve({
          sampleId: this.sample.id,
          label: selectedLabel.value,
          confidence: parseInt(confidence),
          timestamp: new Date()
        });
      }
    });

    this.container.querySelector('.btn-skip').addEventListener('click', () => {
      this.container.remove();
      resolve(null); // Skip this sample
    });

    // Update confidence display
    this.container.querySelector('.confidence-input').addEventListener('input', (e) => {
      this.container.querySelector('.confidence-value').textContent = `${e.target.value}/5`;
    });
  }
}
```

## Advanced Active Learning Patterns

### 1. **Multi-Armed Bandit for Query Strategy Selection**
```javascript
class AdaptiveQueryStrategy {
  constructor(strategies) {
    this.strategies = strategies;
    this.strategyPerformance = new Map();
    this.strategySelections = new Map();
    
    strategies.forEach(strategy => {
      this.strategyPerformance.set(strategy.name, { successes: 1, trials: 2 });
      this.strategySelections.set(strategy.name, 0);
    });
  }

  selectStrategyAndQueries(unlabeledData, budget, model) {
    // Use Thompson Sampling to select strategy
    const strategy = this.selectStrategyByThompsonSampling();
    this.strategySelections.set(strategy.name, this.strategySelections.get(strategy.name) + 1);
    
    // Use selected strategy to choose queries
    const queries = strategy.selectQueries(unlabeledData, budget, model);
    
    return {
      strategy: strategy.name,
      queries,
      strategyInstance: strategy
    };
  }

  selectStrategyByThompsonSampling() {
    let bestStrategy = null;
    let highestSample = -1;

    this.strategies.forEach(strategy => {
      const perf = this.strategyPerformance.get(strategy.name);
      const sample = this.sampleBeta(perf.successes, perf.trials - perf.successes);
      
      if (sample > highestSample) {
        highestSample = sample;
        bestStrategy = strategy;
      }
    });

    return bestStrategy;
  }

  updateStrategyPerformance(strategyName, utilityGain) {
    const perf = this.strategyPerformance.get(strategyName);
    if (utilityGain > 0) {
      perf.successes += 1;
    }
    perf.trials += 1;
  }

  sampleBeta(alpha, beta) {
    // Simple Beta distribution sampling
    const gamma1 = this.sampleGamma(alpha, 1);
    const gamma2 = this.sampleGamma(beta, 1);
    return gamma1 / (gamma1 + gamma2);
  }

  sampleGamma(alpha, beta) {
    // Simplified Gamma sampling
    let sum = 0;
    for (let i = 0; i < alpha; i++) {
      sum -= Math.log(Math.random());
    }
    return sum / beta;
  }
}
```

### 2. **Batch Active Learning with Diversity**
```javascript
class BatchActiveLearning {
  constructor(model, batchSize = 20) {
    this.model = model;
    this.batchSize = batchSize;
    this.diversityWeight = 0.3;
    this.uncertaintyWeight = 0.7;
  }

  selectBatch(unlabeledData, existingData = []) {
    // Greedy algorithm for batch selection
    const selectedBatch = [];
    const remainingData = [...unlabeledData];
    
    while (selectedBatch.length < this.batchSize && remainingData.length > 0) {
      const nextSample = this.selectNextSample(remainingData, selectedBatch, existingData);
      selectedBatch.push(nextSample);
      this.removeSample(remainingData, nextSample);
    }
    
    return selectedBatch;
  }

  selectNextSample(candidates, selectedBatch, existingData) {
    let bestScore = -Infinity;
    let bestSample = null;
    
    candidates.forEach(sample => {
      const uncertaintyScore = this.calculateUncertainty(sample);
      const diversityScore = this.calculateDiversity(sample, [...selectedBatch, ...existingData]);
      
      const totalScore = 
        this.uncertaintyWeight * uncertaintyScore + 
        this.diversityWeight * diversityScore;
      
      if (totalScore > bestScore) {
        bestScore = totalScore;
        bestSample = sample;
      }
    });
    
    return bestSample;
  }

  calculateUncertainty(sample) {
    const prediction = this.model.predict(sample.features);
    return 1 - prediction.confidence; // Higher uncertainty = higher score
  }

  calculateDiversity(sample, comparisonSet) {
    if (comparisonSet.length === 0) return 1;
    
    let minDistance = Infinity;
    comparisonSet.forEach(existingSample => {
      const distance = this.calculateDistance(sample.features, existingSample.features);
      if (distance < minDistance) {
        minDistance = distance;
      }
    });
    
    return minDistance; // Higher distance = more diverse
  }

  calculateDistance(featuresA, featuresB) {
    // Cosine distance
    const dotProduct = this.dotProduct(featuresA, featuresB);
    const normA = this.norm(featuresA);
    const normB = this.norm(featuresB);
    
    return 1 - (dotProduct / (normA * normB));
  }

  dotProduct(vecA, vecB) {
    return Object.keys(vecA).reduce((sum, key) => {
      return sum + ((vecA[key] || 0) * (vecB[key] || 0));
    }, 0);
  }

  norm(vec) {
    return Math.sqrt(Object.values(vec).reduce((sum, val) => sum + val * val, 0));
  }
}
```

## Performance Monitoring and Visualization

```javascript
class ActiveLearningMonitor {
  constructor(activeLearningLoop) {
    this.loop = activeLearningLoop;
    this.metrics = {
      accuracy: [],
      dataEfficiency: [],
      labelingCost: [],
      modelConfidence: []
    };
  }

  trackPerformance() {
    const currentPerformance = this.loop.performanceHistory.slice(-1)[0];
    
    this.metrics.accuracy.push({
      iteration: this.loop.iteration,
      value: currentPerformance.performance.accuracy
    });
    
    this.metrics.dataEfficiency.push({
      iteration: this.loop.iteration,
      value: currentPerformance.performance.accuracy / this.loop.labeledData.length
    });
    
    this.metrics.labelingCost.push({
      iteration: this.loop.iteration,
      value: this.loop.labeledData.length
    });
  }

  generateReport() {
    return {
      summary: {
        totalIterations: this.loop.iteration,
        totalLabeled: this.loop.labeledData.length,
        currentAccuracy: this.metrics.accuracy.slice(-1)[0]?.value || 0,
        improvement: this.calculateImprovement()
      },
      trends: this.metrics,
      recommendations: this.generateRecommendations()
    };
  }

  calculateImprovement() {
    if (this.metrics.accuracy.length < 2) return 0;
    
    const first = this.metrics.accuracy[0].value;
    const last = this.metrics.accuracy.slice(-1)[0].value;
    
    return ((last - first) / first) * 100;
  }

  generateRecommendations() {
    const recommendations = [];
    
    const recentAccuracy = this.metrics.accuracy.slice(-3).map(m => m.value);
    const isPlateauing = this.checkPlateau(recentAccuracy);
    
    if (isPlateauing) {
      recommendations.push({
        type: 'strategy_change',
        message: 'Model performance is plateauing. Consider changing query strategy.',
        priority: 'medium'
      });
    }
    
    if (this.loop.unlabeledData.length < 100) {
      recommendations.push({
        type: 'data_warning',
        message: 'Low on unlabeled data. Consider acquiring more samples.',
        priority: 'high'
      });
    }
    
    return recommendations;
  }

  checkPlateau(values, threshold = 0.01) {
    if (values.length < 3) return false;
    
    const changes = [];
    for (let i = 1; i < values.length; i++) {
      changes.push(Math.abs(values[i] - values[i-1]));
    }
    
    const avgChange = changes.reduce((sum, change) => sum + change, 0) / changes.length;
    return avgChange < threshold;
  }

  visualizePerformance() {
    // This would integrate with a charting library in practice
    return {
      accuracyChart: this.createAccuracyChart(),
      efficiencyChart: this.createEfficiencyChart(),
      costBenefitChart: this.createCostBenefitChart()
    };
  }

  createAccuracyChart() {
    return {
      type: 'line',
      data: this.metrics.accuracy,
      title: 'Model Accuracy Over Iterations',
      xAxis: 'Iteration',
      yAxis: 'Accuracy'
    };
  }
}
```

## Key Benefits of Active Learning Loops

1. **Reduced Labeling Costs**: Focus on most valuable samples
2. **Faster Model Improvement**: Strategic data selection
3. **Human-in-the-Loop**: Combines AI efficiency with human expertise
4. **Adaptive Learning**: Continuously improves based on feedback
5. **Resource Optimization**: Maximizes ROI on labeling efforts

Active learning loops are particularly valuable when labeling is expensive, time-consuming, or requires specialized expertise, making them essential for many real-world machine learning applications.