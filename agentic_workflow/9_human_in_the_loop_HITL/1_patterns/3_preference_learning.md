# Preference Learning

Preference learning is a subfield of machine learning that deals with learning from expressed preferences rather than traditional labeled examples. Instead of learning from absolute labels, the system learns from relative comparisons between items.

## Types of Preference Learning

### 1. **Object Preference Learning**
Learning preferences over individual items
### 2. **Label Preference Learning**  
Learning preferences over multiple labels for an instance
### 3. **Instance Preference Learning**
Learning preferences between different instances

## Key Algorithms and Approaches

### 1. **Pairwise Comparison-Based Learning**
```javascript
class PairwisePreferenceLearner {
  constructor() {
    this.preferences = new Map(); // Store pairwise preferences
    this.featureWeights = new Map(); // Learned feature importance
    this.trainingHistory = [];
  }

  // Add a preference: itemA > itemB (itemA is preferred over itemB)
  addPreference(itemA, itemB, context = {}) {
    const key = this.getComparisonKey(itemA.id, itemB.id);
    this.preferences.set(key, {
      preferred: itemA,
      rejected: itemB,
      timestamp: new Date(),
      context
    });
  }

  // Learn weights using pairwise comparison data
  learnWeights(iterations = 1000, learningRate = 0.01) {
    const features = this.extractAllFeatures();
    
    // Initialize weights
    features.forEach(feature => {
      if (!this.featureWeights.has(feature)) {
        this.featureWeights.set(feature, 0);
      }
    });

    // Stochastic gradient descent
    for (let i = 0; i < iterations; i++) {
      let totalError = 0;
      
      this.preferences.forEach((preference, key) => {
        const preferredScore = this.scoreItem(preference.preferred);
        const rejectedScore = this.scoreItem(preference.rejected);
        
        // Preference should be: preferredScore > rejectedScore
        const margin = preferredScore - rejectedScore;
        const error = Math.max(0, 1 - margin); // Hinge loss
        
        if (error > 0) {
          // Update weights
          this.updateWeights(preference.preferred, preference.rejected, learningRate);
          totalError += error;
        }
      });

      this.trainingHistory.push({ iteration: i, error: totalError });
      
      if (totalError === 0) break;
    }
  }

  scoreItem(item) {
    let score = 0;
    Object.entries(item.features).forEach(([feature, value]) => {
      const weight = this.featureWeights.get(feature) || 0;
      score += weight * this.normalizeValue(value, feature);
    });
    return score;
  }

  updateWeights(preferred, rejected, learningRate) {
    Object.keys(preferred.features).forEach(feature => {
      const prefValue = this.normalizeValue(preferred.features[feature], feature);
      const rejValue = this.normalizeValue(rejected.features[feature], feature);
      const diff = prefValue - rejValue;
      
      const currentWeight = this.featureWeights.get(feature) || 0;
      this.featureWeights.set(feature, currentWeight + learningRate * diff);
    });
  }

  predictPreference(itemA, itemB) {
    const scoreA = this.scoreItem(itemA);
    const scoreB = this.scoreItem(itemB);
    
    return {
      preferred: scoreA > scoreB ? itemA : itemB,
      confidence: Math.abs(scoreA - scoreB),
      scores: { [itemA.id]: scoreA, [itemB.id]: scoreB }
    };
  }

  extractAllFeatures() {
    const features = new Set();
    this.preferences.forEach(preference => {
      Object.keys(preference.preferred.features).forEach(f => features.add(f));
      Object.keys(preference.rejected.features).forEach(f => features.add(f));
    });
    return Array.from(features);
  }

  normalizeValue(value, feature) {
    // Simple normalization - in practice, you'd want more sophisticated scaling
    if (typeof value === 'number') {
      return value / 100; // Example normalization
    }
    return value;
  }

  getComparisonKey(idA, idB) {
    return [idA, idB].sort().join('_');
  }
}

// Usage Example
const learner = new PairwisePreferenceLearner();

// Training data: user preferences
const items = {
  item1: { id: 'item1', features: { price: 100, rating: 4.5, category: 'electronics' } },
  item2: { id: 'item2', features: { price: 150, rating: 4.8, category: 'electronics' } },
  item3: { id: 'item3', features: { price: 80, rating: 4.2, category: 'electronics' } }
};

// Add preferences: item1 > item2, item3 > item1, etc.
learner.addPreference(items.item1, items.item2);
learner.addPreference(items.item3, items.item1);
learner.addPreference(items.item2, items.item3);

// Learn from preferences
learner.learnWeights(500, 0.1);

// Predict new preferences
const prediction = learner.predictPreference(items.item1, items.item3);
console.log('Preference prediction:', prediction);
```

### 2. **Bayesian Preference Learning**
```javascript
class BayesianPreferenceLearner {
  constructor() {
    this.preferences = [];
    this.featurePosteriors = new Map();
    this.hyperparameters = { alpha: 1, beta: 1 }; // Prior parameters
  }

  // Model preference as Bernoulli distribution
  addPreference(choice, alternatives, context) {
    this.preferences.push({
      choice,
      alternatives,
      context,
      timestamp: new Date()
    });
  }

  // Update posterior distributions using Bayesian inference
  updatePosteriors() {
    // For each feature, update Beta posterior based on preferences
    const featureStats = this.computeFeatureStatistics();
    
    featureStats.forEach((stats, feature) => {
      const successes = stats.preferredCount + this.hyperparameters.alpha;
      const failures = stats.rejectedCount + this.hyperparameters.beta;
      
      this.featurePosteriors.set(feature, {
        distribution: 'Beta',
        alpha: successes,
        beta: failures,
        mean: successes / (successes + failures)
      });
    });
  }

  computeFeatureStatistics() {
    const stats = new Map();
    
    this.preferences.forEach(pref => {
      const chosen = pref.choice;
      const rejected = pref.alternatives.filter(alt => alt.id !== chosen.id);
      
      // Update statistics for chosen item's features
      Object.entries(chosen.features).forEach(([feature, value]) => {
        if (!stats.has(feature)) {
          stats.set(feature, { preferredCount: 0, rejectedCount: 0, totalValue: 0 });
        }
        const featureStat = stats.get(feature);
        featureStat.preferredCount += this.getFeatureImportance(value);
        featureStat.totalValue += value;
      });
      
      // Update statistics for rejected items' features
      rejected.forEach(item => {
        Object.entries(item.features).forEach(([feature, value]) => {
          if (!stats.has(feature)) {
            stats.set(feature, { preferredCount: 0, rejectedCount: 0, totalValue: 0 });
          }
          const featureStat = stats.get(feature);
          featureStat.rejectedCount += this.getFeatureImportance(value);
          featureStat.totalValue += value;
        });
      });
    });
    
    return stats;
  }

  predictChoice(alternatives) {
    const scores = alternatives.map(alt => ({
      item: alt,
      score: this.computeChoiceProbability(alt)
    }));
    
    scores.sort((a, b) => b.score - a.score);
    
    return {
      predictedChoice: scores[0].item,
      confidence: scores[0].score,
      allScores: scores
    };
  }

  computeChoiceProbability(item) {
    let probability = 1.0;
    
    Object.entries(item.features).forEach(([feature, value]) => {
      const posterior = this.featurePosteriors.get(feature);
      if (posterior) {
        const featureProb = posterior.mean;
        const featureWeight = this.getFeatureImportance(value);
        probability *= Math.pow(featureProb, featureWeight);
      }
    });
    
    return probability;
  }

  getFeatureImportance(value) {
    // Simple importance weighting - could be more sophisticated
    if (typeof value === 'number') {
      return Math.abs(value) / 100;
    }
    return 1;
  }

  getUncertainty(feature) {
    const posterior = this.featurePosteriors.get(feature);
    if (posterior) {
      // Variance of Beta distribution
      return (posterior.alpha * posterior.beta) / 
             (Math.pow(posterior.alpha + posterior.beta, 2) * 
              (posterior.alpha + posterior.beta + 1));
    }
    return 1; // Maximum uncertainty for unknown features
  }

  // Active learning: suggest comparisons that reduce uncertainty
  suggestInformativeComparison(availableItems) {
    let maxUncertaintyReduction = -Infinity;
    let bestPair = null;
    
    for (let i = 0; i < availableItems.length; i++) {
      for (let j = i + 1; j < availableItems.length; j++) {
        const uncertaintyReduction = this.estimateUncertaintyReduction(
          availableItems[i], 
          availableItems[j]
        );
        
        if (uncertaintyReduction > maxUncertaintyReduction) {
          maxUncertaintyReduction = uncertaintyReduction;
          bestPair = [availableItems[i], availableItems[j]];
        }
      }
    }
    
    return bestPair;
  }

  estimateUncertaintyReduction(itemA, itemB) {
    // Simplified uncertainty reduction estimation
    let totalUncertainty = 0;
    
    const allFeatures = new Set([
      ...Object.keys(itemA.features),
      ...Object.keys(itemB.features)
    ]);
    
    allFeatures.forEach(feature => {
      totalUncertainty += this.getUncertainty(feature);
    });
    
    return totalUncertainty;
  }
}
```

### 3. **Neural Preference Learning with Embeddings**
```javascript
class NeuralPreferenceModel {
  constructor(inputSize, hiddenLayers = [64, 32]) {
    this.inputSize = inputSize;
    this.hiddenLayers = hiddenLayers;
    this.model = this.buildModel();
    this.preferencePairs = [];
    this.lossHistory = [];
  }

  buildModel() {
    // Simple neural network for preference learning
    const model = {
      weights: [],
      biases: []
    };

    // Initialize weights and biases
    let previousSize = this.inputSize * 2; // Concatenated features of two items
    
    this.hiddenLayers.forEach(layerSize => {
      model.weights.push(this.initializeWeights(previousSize, layerSize));
      model.biases.push(new Array(layerSize).fill(0));
      previousSize = layerSize;
    });

    // Output layer (preference probability)
    model.weights.push(this.initializeWeights(previousSize, 1));
    model.biases.push([0]);

    return model;
  }

  initializeWeights(rows, cols) {
    return Array.from({ length: rows }, () => 
      Array.from({ length: cols }, () => (Math.random() - 0.5) * 0.1)
    );
  }

  addPreferencePair(preferredItem, rejectedItem) {
    this.preferencePairs.push({
      preferred: this.extractFeatures(preferredItem),
      rejected: this.extractFeatures(rejectedItem),
      label: 1 // Preferred > Rejected
    });
  }

  extractFeatures(item) {
    // Convert item features to numerical vector
    const features = [];
    Object.values(item.features).forEach(value => {
      if (typeof value === 'number') {
        features.push(value);
      } else if (typeof value === 'string') {
        // Simple string encoding
        features.push(value.length / 10);
      }
    });
    return features;
  }

  forwardPass(input) {
    let activation = input;
    
    for (let i = 0; i < this.model.weights.length; i++) {
      activation = this.matrixMultiply(activation, this.model.weights[i]);
      activation = this.vectorAdd(activation, this.model.biases[i]);
      
      if (i < this.model.weights.length - 1) {
        activation = this.relu(activation);
      } else {
        activation = this.sigmoid(activation); // Final output probability
      }
    }
    
    return activation[0];
  }

  predictPreference(itemA, itemB) {
    const featuresA = this.extractFeatures(itemA);
    const featuresB = this.extractFeatures(itemB);
    const combinedInput = [...featuresA, ...featuresB];
    
    const scoreAB = this.forwardPass(combinedInput);
    const scoreBA = this.forwardPass([...featuresB, ...featuresA]);
    
    return {
      preference: scoreAB > scoreBA ? itemA : itemB,
      confidence: Math.abs(scoreAB - scoreBA),
      probabilities: { [itemA.id]: scoreAB, [itemB.id]: scoreBA }
    };
  }

  train(epochs = 100, learningRate = 0.01) {
    for (let epoch = 0; epoch < epochs; epoch++) {
      let totalLoss = 0;
      
      this.preferencePairs.forEach(pair => {
        const preferredFeatures = pair.preferred;
        const rejectedFeatures = pair.rejected;
        
        // Forward pass for preferred > rejected
        const inputPreferredFirst = [...preferredFeatures, ...rejectedFeatures];
        const scorePreferred = this.forwardPass(inputPreferredFirst);
        
        // Loss: we want scorePreferred to be close to 1
        const loss = -Math.log(scorePreferred);
        totalLoss += loss;
        
        // Backward pass (simplified)
        this.updateWeights(inputPreferredFirst, scorePreferred, 1, learningRate);
      });
      
      this.lossHistory.push(totalLoss / this.preferencePairs.length);
      
      if (epoch % 10 === 0) {
        console.log(`Epoch ${epoch}, Loss: ${totalLoss / this.preferencePairs.length}`);
      }
    }
  }

  updateWeights(input, predicted, target, learningRate) {
    // Simplified gradient update
    const error = predicted - target;
    
    // This is a simplified version - real implementation would use proper backpropagation
    this.model.weights.forEach((layerWeights, layerIndex) => {
      layerWeights.forEach((neuronWeights, neuronIndex) => {
        neuronWeights.forEach((weight, weightIndex) => {
          const gradient = error * input[weightIndex] || 0;
          this.model.weights[layerIndex][neuronIndex][weightIndex] -= learningRate * gradient;
        });
      });
    });
  }

  // Utility functions for neural network operations
  matrixMultiply(a, b) {
    const result = new Array(b[0].length).fill(0);
    for (let i = 0; i < a.length; i++) {
      for (let j = 0; j < b[0].length; j++) {
        result[j] += a[i] * b[i][j];
      }
    }
    return result;
  }

  vectorAdd(a, b) {
    return a.map((val, i) => val + b[i]);
  }

  relu(x) {
    return x.map(val => Math.max(0, val));
  }

  sigmoid(x) {
    return x.map(val => 1 / (1 + Math.exp(-val)));
  }
}
```

## Preference Learning Applications

### 1. **Recommendation Systems**
```javascript
class PreferenceBasedRecommender {
  constructor(preferenceLearner) {
    this.learner = preferenceLearner;
    this.userHistory = [];
    this.itemCatalog = [];
  }

  addUserInteraction(userId, preferredItem, context) {
    this.userHistory.push({
      userId,
      preferredItem,
      context,
      timestamp: new Date()
    });
  }

  generateRecommendations(userId, candidateItems, count = 10) {
    const userPreferences = this.getUserPreferences(userId);
    const scoredItems = candidateItems.map(item => ({
      item,
      score: this.calculateRecommendationScore(item, userPreferences)
    }));

    return scoredItems
      .sort((a, b) => b.score - a.score)
      .slice(0, count)
      .map(result => result.item);
  }

  calculateRecommendationScore(item, userPreferences) {
    let score = 0;
    
    // Compare against user's past preferences
    userPreferences.forEach(pref => {
      const comparison = this.learner.predictPreference(item, pref.preferredItem);
      if (comparison.preferred.id === item.id) {
        score += comparison.confidence;
      }
    });

    return score;
  }

  getUserPreferences(userId) {
    return this.userHistory.filter(interaction => interaction.userId === userId);
  }

  // Explore-Exploit strategy
  exploreExploitRecommendations(userId, candidateItems, explorationRate = 0.1) {
    if (Math.random() < explorationRate) {
      // Exploration: recommend diverse items
      return this.exploreDiverseItems(candidateItems);
    } else {
      // Exploitation: recommend based on known preferences
      return this.generateRecommendations(userId, candidateItems);
    }
  }

  exploreDiverseItems(items) {
    // Simple diversity: shuffle and take top N
    return [...items]
      .sort(() => Math.random() - 0.5)
      .slice(0, 10);
  }
}
```

### 2. **A/B Testing with Preference Learning**
```javascript
class PreferenceDrivenABTesting {
  constructor(variants) {
    this.variants = variants;
    this.preferenceData = [];
    this.variantPerformance = new Map();
  }

  recordPreference(userId, preferredVariant, context) {
    this.preferenceData.push({
      userId,
      preferredVariant,
      context,
      timestamp: new Date()
    });

    this.updateVariantPerformance(preferredVariant);
  }

  updateVariantPerformance(variantId) {
    if (!this.variantPerformance.has(variantId)) {
      this.variantPerformance.set(variantId, { preferences: 0, impressions: 0 });
    }
    
    const performance = this.variantPerformance.get(variantId);
    performance.preferences += 1;
  }

  getOptimalVariant() {
    let bestVariant = null;
    let highestPreferenceRate = -1;

    this.variantPerformance.forEach((performance, variantId) => {
      const preferenceRate = performance.preferences / performance.impressions;
      
      if (preferenceRate > highestPreferenceRate) {
        highestPreferenceRate = preferenceRate;
        bestVariant = variantId;
      }
    });

    return bestVariant;
  }

  // Thompson Sampling for variant selection
  selectVariantThompson() {
    const variantSamples = [];
    
    this.variantPerformance.forEach((performance, variantId) => {
      const alpha = performance.preferences + 1;
      const beta = performance.impressions - performance.preferences + 1;
      
      // Sample from Beta distribution
      const sample = this.sampleBeta(alpha, beta);
      variantSamples.push({ variantId, sample });
    });

    return variantSamples.sort((a, b) => b.sample - a.sample)[0].variantId;
  }

  sampleBeta(alpha, beta) {
    // Simple Beta distribution sampling
    const gammaAlpha = this.sampleGamma(alpha, 1);
    const gammaBeta = this.sampleGamma(beta, 1);
    return gammaAlpha / (gammaAlpha + gammaBeta);
  }

  sampleGamma(alpha, beta) {
    // Marsaglia and Tsang method for Gamma distribution sampling
    if (alpha > 1) {
      const d = alpha - 1/3;
      const c = 1/Math.sqrt(9*d);
      
      while (true) {
        let x, v;
        do {
          x = this.sampleGaussian();
          v = 1 + c*x;
        } while (v <= 0);
        
        v = v*v*v;
        const u = Math.random();
        
        if (u < 1 - 0.0331*x*x*x*x || Math.log(u) < 0.5*x*x + d*(1 - v + Math.log(v))) {
          return beta * d * v;
        }
      }
    }
    return 0;
  }

  sampleGaussian() {
    // Box-Muller transform
    const u1 = Math.random();
    const u2 = Math.random();
    return Math.sqrt(-2 * Math.log(u1)) * Math.cos(2 * Math.PI * u2);
  }
}
```

## Key Challenges in Preference Learning

1. **Data Sparsity**: Preferences are often scarce compared to total possible comparisons
2. **Inconsistent Preferences**: Humans may provide contradictory preferences
3. **Context Dependence**: Preferences can change based on context
4. **Scalability**: Pairwise comparisons grow quadratically with number of items
5. **Cold Start**: Limited data for new users/items

Preference learning is particularly valuable in scenarios where explicit ratings are hard to obtain, but relative comparisons are natural for users to provide.