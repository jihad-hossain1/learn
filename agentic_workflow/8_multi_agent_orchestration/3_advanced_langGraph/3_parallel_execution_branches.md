# Parallel Execution

Parallel execution refers to running multiple tasks simultaneously rather than sequentially, significantly improving performance and efficiency in systems with multiple processors, cores, or distributed resources.

## Core Concepts

### 1. **Types of Parallelism**
- **Task Parallelism**: Different tasks running concurrently
- **Data Parallelism**: Same operation on different data sets
- **Pipeline Parallelism**: Different stages of processing running concurrently

### 2. **Parallel Execution Patterns**
- **Fork-Join**: Split work, process in parallel, combine results
- **Map-Reduce**: Transform data in parallel, aggregate results
- **Pipeline**: Sequential stages with parallel processing within stages

## JavaScript Implementation Examples

### Basic Parallel Execution Framework

```javascript
// Parallel Execution Manager
class ParallelExecutor {
  constructor(maxConcurrency = 4) {
    this.maxConcurrency = maxConcurrency;
    this.activeTasks = new Set();
    this.queue = [];
    this.taskId = 0;
    this.stats = {
      totalTasks: 0,
      completedTasks: 0,
      failedTasks: 0,
      totalExecutionTime: 0
    };
  }

  async execute(tasks, options = {}) {
    const {
      timeout = 0,
      breakOnError = false,
      progressCallback = null
    } = options;

    const taskPromises = [];
    const results = [];
    const errors = [];

    // Create task entries
    const taskEntries = tasks.map((task, index) => ({
      id: this.taskId++,
      task,
      index,
      status: 'pending',
      startTime: null,
      endTime: null
    }));

    this.stats.totalTasks += taskEntries.length;

    // Execute with concurrency control
    const executing = new Set();
    
    for (const taskEntry of taskEntries) {
      // Wait if we've reached max concurrency
      while (executing.size >= this.maxConcurrency) {
        await Promise.race(executing);
      }

      const taskPromise = this.executeSingleTask(taskEntry, timeout);
      executing.add(taskPromise);
      
      taskPromise
        .then(result => {
          results[taskEntry.index] = result;
          taskEntry.status = 'completed';
        })
        .catch(error => {
          errors[taskEntry.index] = error;
          taskEntry.status = 'failed';
          this.stats.failedTasks++;
          
          if (breakOnError) {
            // Cancel remaining tasks
            taskEntries.forEach(te => {
              if (te.status === 'pending') te.status = 'cancelled';
            });
          }
        })
        .finally(() => {
          executing.delete(taskPromise);
          this.stats.completedTasks++;
          
          if (progressCallback) {
            progressCallback({
              completed: this.stats.completedTasks,
              total: this.stats.totalTasks,
              successful: results.length,
              failed: errors.length
            });
          }
        });

      taskPromises.push(taskPromise);
    }

    // Wait for all tasks to complete
    await Promise.allSettled(taskPromises);

    return {
      results,
      errors: errors.filter(e => e), // Remove undefined entries
      stats: this.getExecutionStats()
    };
  }

  async executeSingleTask(taskEntry, timeout = 0) {
    taskEntry.startTime = Date.now();
    taskEntry.status = 'running';

    let timeoutId;
    const timeoutPromise = timeout > 0 ? new Promise((_, reject) => {
      timeoutId = setTimeout(() => {
        reject(new Error(`Task ${taskEntry.id} timed out after ${timeout}ms`));
      }, timeout);
    }) : null;

    try {
      const taskPromise = typeof taskEntry.task === 'function' 
        ? taskEntry.task() 
        : taskEntry.task;

      const result = await (timeoutPromise 
        ? Promise.race([taskPromise, timeoutPromise])
        : taskPromise
      );

      if (timeoutId) clearTimeout(timeoutId);
      
      taskEntry.endTime = Date.now();
      taskEntry.executionTime = taskEntry.endTime - taskEntry.startTime;
      this.stats.totalExecutionTime += taskEntry.executionTime;

      return result;
    } catch (error) {
      if (timeoutId) clearTimeout(timeoutId);
      taskEntry.endTime = Date.now();
      taskEntry.error = error;
      throw error;
    }
  }

  getExecutionStats() {
    return {
      ...this.stats,
      successRate: this.stats.totalTasks > 0 
        ? (this.stats.totalTasks - this.stats.failedTasks) / this.stats.totalTasks 
        : 0,
      averageExecutionTime: this.stats.completedTasks > 0
        ? this.stats.totalExecutionTime / this.stats.completedTasks
        : 0
    };
  }
}

// Specialized Parallel Patterns
class ParallelPatterns {
  // Map-Reduce Pattern
  static async mapReduce(data, mapper, reducer, concurrency = 4) {
    console.log('🚀 Starting Map-Reduce Process');
    
    // Map phase (parallel)
    const executor = new ParallelExecutor(concurrency);
    
    const mapTasks = data.map((item, index) => async () => {
      console.log(`📊 Mapping item ${index + 1}/${data.length}`);
      return await mapper(item, index);
    });

    const mapResult = await executor.execute(mapTasks);
    
    if (mapResult.errors.length > 0) {
      throw new Error(`Map phase failed: ${mapResult.errors[0].message}`);
    }

    console.log('✅ Map phase completed');

    // Reduce phase (sequential or parallel depending on reducer)
    console.log('🔗 Starting Reduce phase');
    const finalResult = await reducer(mapResult.results);
    console.log('✅ Reduce phase completed');

    return {
      finalResult,
      mapResults: mapResult.results,
      stats: mapResult.stats
    };
  }

  // Fork-Join Pattern
  static async forkJoin(tasks, concurrency = 4) {
    console.log('🎯 Starting Fork-Join Process');
    
    const executor = new ParallelExecutor(concurrency);
    const taskEntries = Object.entries(tasks).map(([key, task]) => ({
      key,
      task: typeof task === 'function' ? task : () => task
    }));

    const forkTasks = taskEntries.map(entry => async () => {
      console.log(`🔧 Executing fork task: ${entry.key}`);
      return await entry.task();
    });

    const result = await executor.execute(forkTasks);
    
    // Join results back into object structure
    const joinedResults = {};
    taskEntries.forEach((entry, index) => {
      joinedResults[entry.key] = result.results[index];
    });

    console.log('✅ Fork-Join completed');

    return {
      results: joinedResults,
      errors: result.errors,
      stats: result.stats
    };
  }

  // Pipeline Pattern with Parallel Stages
  static async pipeline(data, stages, concurrency = 4) {
    console.log('🏭 Starting Pipeline Process');
    
    let currentData = data;
    const stageResults = [];

    for (let i = 0; i < stages.length; i++) {
      const stage = stages[i];
      console.log(`\n📋 Processing Stage ${i + 1}/${stages.length}: ${stage.name}`);

      if (stage.parallel) {
        // Process stage in parallel
        const executor = new ParallelExecutor(concurrency);
        
        const tasks = currentData.map((item, index) => async () => {
          return await stage.processor(item, index, currentData);
        });

        const result = await executor.execute(tasks);
        
        if (result.errors.length > 0) {
          throw new Error(`Pipeline stage ${i + 1} failed: ${result.errors[0].message}`);
        }

        currentData = result.results;
        stageResults.push({
          stage: stage.name,
          results: result.results,
          stats: result.stats
        });
      } else {
        // Process stage sequentially
        const results = [];
        for (let j = 0; j < currentData.length; j++) {
          console.log(`Processing item ${j + 1}/${currentData.length}`);
          const result = await stage.processor(currentData[j], j, currentData);
          results.push(result);
        }
        currentData = results;
        stageResults.push({
          stage: stage.name,
          results: results
        });
      }
    }

    console.log('✅ Pipeline completed');
    return {
      finalResult: currentData,
      stageResults
    };
  }
}
```

### Practical Examples

#### Example 1: Image Processing Pipeline
```javascript
// Simulated image processing functions
class ImageProcessor {
  static async resizeImage(image) {
    console.log(`🖼️  Resizing image: ${image.name}`);
    await new Promise(resolve => setTimeout(resolve, 100 + Math.random() * 200));
    return { ...image, size: 'resized', width: 800, height: 600 };
  }

  static async applyFilter(image) {
    console.log(`🎨 Applying filter to: ${image.name}`);
    await new Promise(resolve => setTimeout(resolve, 150 + Math.random() * 300));
    return { ...image, filter: 'enhanced', quality: 'high' };
  }

  static async optimizeImage(image) {
    console.log(`⚡ Optimizing image: ${image.name}`);
    await new Promise(resolve => setTimeout(resolve, 200 + Math.random() * 250));
    return { ...image, optimized: true, sizeKB: Math.random() * 100 + 50 };
  }

  static async addWatermark(image) {
    console.log(`💧 Adding watermark to: ${image.name}`);
    await new Promise(resolve => setTimeout(resolve, 100 + Math.random() * 150));
    return { ...image, watermarked: true, copyright: '2024' };
  }
}

async function demonstrateImageProcessing() {
  console.log('=== Image Processing Pipeline ===\n');

  const images = [
    { name: 'photo1.jpg', size: 'original', width: 4000, height: 3000 },
    { name: 'photo2.png', size: 'original', width: 3000, height: 2000 },
    { name: 'photo3.jpeg', size: 'original', width: 5000, height: 4000 },
    { name: 'photo4.jpg', size: 'original', width: 2500, height: 1800 },
    { name: 'photo5.png', size: 'original', width: 3500, height: 2500 }
  ];

  // Define pipeline stages
  const pipelineStages = [
    {
      name: 'Resize',
      parallel: true,
      processor: ImageProcessor.resizeImage
    },
    {
      name: 'Filter Application',
      parallel: true,
      processor: ImageProcessor.applyFilter
    },
    {
      name: 'Optimization',
      parallel: true,
      processor: ImageProcessor.optimizeImage
    },
    {
      name: 'Watermark',
      parallel: false, // Sequential to ensure proper ordering
      processor: ImageProcessor.addWatermark
    }
  ];

  const result = await ParallelPatterns.pipeline(images, pipelineStages, 3);
  
  console.log('\n📊 Pipeline Results:');
  result.stageResults.forEach((stage, index) => {
    console.log(`\nStage ${index + 1}: ${stage.stage}`);
    console.log(`Processed ${stage.results.length} images`);
    if (stage.stats) {
      console.log(`Average time: ${stage.stats.averageExecutionTime.toFixed(2)}ms`);
    }
  });

  console.log('\n✅ Final processed images:');
  result.finalResult.forEach(image => {
    console.log(`- ${image.name}: ${image.width}x${image.height}, ${image.filter}, ${image.optimized ? 'optimized' : 'not optimized'}`);
  });
}
```

#### Example 2: Data Analysis with Map-Reduce
```javascript
async function demonstrateMapReduce() {
  console.log('\n=== Data Analysis with Map-Reduce ===\n');

  // Sample data: website analytics
  const pageViews = [
    { page: '/home', views: 1500, duration: 45 },
    { page: '/about', views: 800, duration: 60 },
    { page: '/products', views: 1200, duration: 30 },
    { page: '/contact', views: 400, duration: 120 },
    { page: '/blog', views: 2000, duration: 90 },
    { page: '/pricing', views: 600, duration: 50 }
  ];

  // Mapper: Calculate engagement score for each page
  const mapper = async (pageData, index) => {
    await new Promise(resolve => setTimeout(resolve, 100)); // Simulate processing
    const engagement = (pageData.views * pageData.duration) / 1000;
    return {
      page: pageData.page,
      views: pageData.views,
      engagement: Math.round(engagement * 100) / 100,
      category: pageData.page === '/blog' ? 'content' : 'navigation'
    };
  };

  // Reducer: Aggregate results by category
  const reducer = async (mappedResults) => {
    const categories = {};
    
    mappedResults.forEach(result => {
      if (!categories[result.category]) {
        categories[result.category] = {
          totalViews: 0,
          totalEngagement: 0,
          pages: []
        };
      }
      
      categories[result.category].totalViews += result.views;
      categories[result.category].totalEngagement += result.engagement;
      categories[result.category].pages.push(result.page);
    });

    // Calculate averages
    Object.keys(categories).forEach(category => {
      categories[category].averageEngagement = 
        categories[category].totalEngagement / categories[category].pages.length;
    });

    return categories;
  };

  const result = await ParallelPatterns.mapReduce(pageViews, mapper, reducer, 2);

  console.log('📈 Analysis Results:');
  Object.entries(result.finalResult).forEach(([category, data]) => {
    console.log(`\n${category.toUpperCase()}:`);
    console.log(`  Total Views: ${data.totalViews}`);
    console.log(`  Total Engagement: ${data.totalEngagement}`);
    console.log(`  Average Engagement: ${data.averageEngagement.toFixed(2)}`);
    console.log(`  Pages: ${data.pages.join(', ')}`);
  });

  console.log(`\n📊 Map-Reduce Statistics:`);
  console.log(`Total processing time: ${result.stats.totalExecutionTime}ms`);
  console.log(`Average task time: ${result.stats.averageExecutionTime.toFixed(2)}ms`);
  console.log(`Success rate: ${(result.stats.successRate * 100).toFixed(1)}%`);
}
```

#### Example 3: API Data Aggregation with Fork-Join
```javascript
async function demonstrateForkJoin() {
  console.log('\n=== API Data Aggregation with Fork-Join ===\n');

  // Simulated API calls
  const apiTasks = {
    users: async () => {
      console.log('👥 Fetching users...');
      await new Promise(resolve => setTimeout(resolve, 300));
      return [{ id: 1, name: 'Alice' }, { id: 2, name: 'Bob' }];
    },
    
    products: async () => {
      console.log('📦 Fetching products...');
      await new Promise(resolve => setTimeout(resolve, 400));
      return [{ id: 1, name: 'Laptop' }, { id: 2, name: 'Phone' }];
    },
    
    orders: async () => {
      console.log('📋 Fetching orders...');
      await new Promise(resolve => setTimeout(resolve, 500));
      return [{ id: 1, userId: 1, productId: 1 }, { id: 2, userId: 2, productId: 2 }];
    },
    
    analytics: async () => {
      console.log('📊 Fetching analytics...');
      await new Promise(resolve => setTimeout(resolve, 200));
      return { visits: 1000, conversions: 50 };
    }
  };

  const result = await ParallelPatterns.forkJoin(apiTasks, 2);

  console.log('\n🔗 Aggregating results...');
  
  // Combine the results
  const dashboardData = {
    summary: {
      totalUsers: result.results.users.length,
      totalProducts: result.results.products.length,
      totalOrders: result.results.orders.length,
      ...result.results.analytics
    },
    users: result.results.users,
    products: result.results.products,
    orders: result.results.orders.map(order => ({
      ...order,
      user: result.results.users.find(u => u.id === order.userId),
      product: result.results.products.find(p => p.id === order.productId)
    }))
  };

  console.log('\n📊 Dashboard Data:');
  console.log(JSON.stringify(dashboardData.summary, null, 2));

  console.log('\n📈 Fork-Join Statistics:');
  console.log(`Total execution time: ${result.stats.totalExecutionTime}ms`);
  console.log(`Parallel efficiency: ${((1 - result.stats.totalExecutionTime / 1400) * 100).toFixed(1)}%`);
}
```

### Advanced Parallel Execution Features

```javascript
// Advanced Parallel Executor with Error Handling and Retries
class AdvancedParallelExecutor extends ParallelExecutor {
  constructor(maxConcurrency = 4, maxRetries = 3) {
    super(maxConcurrency);
    this.maxRetries = maxRetries;
    this.retryDelays = [100, 500, 1000]; // Exponential backoff
  }

  async executeWithRetry(tasks, options = {}) {
    const { retryCondition = () => true } = options;
    
    const finalResults = [];
    const finalErrors = [];

    for (let i = 0; i < tasks.length; i++) {
      let lastError;
      let attempt = 0;

      while (attempt <= this.maxRetries) {
        try {
          if (attempt > 0) {
            console.log(`🔄 Retrying task ${i} (attempt ${attempt + 1})`);
            const delay = this.retryDelays[Math.min(attempt - 1, this.retryDelays.length - 1)];
            await new Promise(resolve => setTimeout(resolve, delay));
          }

          const result = await this.executeSingleTask(
            { id: i, task: tasks[i], index: i },
            options.timeout
          );

          finalResults[i] = result;
          break; // Success, break retry loop

        } catch (error) {
          lastError = error;
          attempt++;

          if (attempt > this.maxRetries || !retryCondition(error)) {
            finalErrors[i] = lastError;
            break;
          }
        }
      }
    }

    return {
      results: finalResults,
      errors: finalErrors.filter(e => e)
    };
  }

  // Batch processing for large datasets
  async executeInBatches(tasks, batchSize = 10) {
    const results = [];
    const errors = [];

    for (let i = 0; i < tasks.length; i += batchSize) {
      const batch = tasks.slice(i, i + batchSize);
      console.log(`Processing batch ${Math.floor(i / batchSize) + 1}/${Math.ceil(tasks.length / batchSize)}`);
      
      const batchResult = await this.execute(batch);
      
      results.push(...batchResult.results);
      errors.push(...batchResult.errors);
      
      // Optional delay between batches
      await new Promise(resolve => setTimeout(resolve, 100));
    }

    return { results, errors };
  }
}

// Demonstration of advanced features
async function demonstrateAdvancedFeatures() {
  console.log('\n=== Advanced Parallel Execution Features ===\n');

  const advancedExecutor = new AdvancedParallelExecutor(3, 2);

  // Tasks with simulated failures
  const unreliableTasks = Array.from({ length: 10 }, (_, i) => async () => {
    await new Promise(resolve => setTimeout(resolve, 100));
    
    // Simulate occasional failures
    if (Math.random() < 0.3) {
      throw new Error(`Task ${i} failed randomly`);
    }
    
    return `Task ${i} completed successfully`;
  });

  console.log('🔄 Testing retry mechanism...');
  const result = await advancedExecutor.executeWithRetry(unreliableTasks, {
    retryCondition: (error) => error.message.includes('failed randomly')
  });

  console.log(`\n📊 Results with retries:`);
  console.log(`Successful: ${result.results.filter(r => r).length}`);
  console.log(`Failed: ${result.errors.length}`);
  console.log(`Success rate: ${(result.results.filter(r => r).length / unreliableTasks.length * 100).toFixed(1)}%`);
}

// Run all demonstrations
async function runAllDemonstrations() {
  await demonstrateImageProcessing();
  await demonstrateMapReduce();
  await demonstrateForkJoin();
  await demonstrateAdvancedFeatures();
}

runAllDemonstrations().catch(console.error);
```

## Key Benefits of Parallel Execution

### 1. **Performance Improvement**
- Significant speedup for CPU-intensive tasks
- Better utilization of multi-core systems

### 2. **Resource Efficiency**
- Parallel I/O operations
- Reduced overall execution time

### 3. **Scalability**
- Handles large datasets efficiently
- Scales with available resources

### 4. **Fault Tolerance**
- Isolated task failures
- Retry mechanisms for failed tasks

### 5. **Flexibility**
- Adaptable concurrency levels
- Multiple parallel patterns for different use cases

This approach is essential for modern applications dealing with large-scale data processing, real-time analytics, and high-performance computing requirements.