Excellent! Let's dive into **Retry and Error Handling Strategies** for long-running processes. This is crucial for building resilient systems that can handle transient failures.

## Core Retry Strategies

### 1. **Exponential Backoff**

### 2. **Jitter/Randomization**

### 3. **Circuit Breaker Pattern**

### 4. **Dead Letter Queues**

### 5. **Fallback Strategies**

Here are comprehensive JavaScript examples for each:

## Strategy 1: Exponential Backoff with Jitter

```javascript
// advanced-retry-strategy.js
class AdvancedRetry {
  constructor(options = {}) {
    this.maxAttempts = options.maxAttempts || 5;
    this.initialDelay = options.initialDelay || 1000; // 1 second
    this.maxDelay = options.maxDelay || 30000; // 30 seconds
    this.backoffFactor = options.backoffFactor || 2;
    this.jitter = options.jitter !== false; // Enable jitter by default
    this.retryableErrors = options.retryableErrors || [
      "ECONNRESET",
      "ETIMEDOUT",
      "ECONNREFUSED",
      "ENOTFOUND",
      "TIMEOUT",
    ];
  }

  // Calculate delay with exponential backoff and jitter
  calculateDelay(attempt) {
    const exponentialDelay =
      this.initialDelay * Math.pow(this.backoffFactor, attempt - 1);
    const delay = Math.min(exponentialDelay, this.maxDelay);

    if (!this.jitter) return delay;

    // Add jitter to avoid thundering herd problem
    const jitter = delay * 0.1 * Math.random(); // ±10% jitter
    return delay + jitter;
  }

  // Check if error is retryable
  isRetryableError(error) {
    if (this.retryableErrors.includes(error.code)) return true;

    // Check HTTP status codes for retry
    if (error.statusCode) {
      return [429, 500, 502, 503, 504].includes(error.statusCode);
    }

    // Check error message patterns
    const retryablePatterns = [
      /timeout/i,
      /network/i,
      /connection/i,
      /server/i,
    ];
    return retryablePatterns.some((pattern) => pattern.test(error.message));
  }

  // Main retry method
  async execute(operation, context = {}) {
    let lastError;

    for (let attempt = 1; attempt <= this.maxAttempts; attempt++) {
      try {
        console.log(`Attempt ${attempt}/${this.maxAttempts}...`);

        const result = await operation(attempt, context);

        // Success - log and return
        if (attempt > 1) {
          console.log(`Operation succeeded on attempt ${attempt}`);
        }

        return result;
      } catch (error) {
        lastError = error;

        // Check if error is retryable
        if (!this.isRetryableError(error) || attempt === this.maxAttempts) {
          console.error(
            `Non-retryable error or max attempts reached:`,
            error.message
          );
          throw this.wrapError(error, attempt, context);
        }

        // Calculate delay and wait
        const delay = this.calculateDelay(attempt);
        console.warn(
          `Attempt ${attempt} failed: ${error.message}. Retrying in ${delay}ms...`
        );

        await this.wait(delay);

        // Update context for next attempt
        context.lastAttempt = attempt;
        context.lastError = error.message;
      }
    }

    throw this.wrapError(lastError, this.maxAttempts, context);
  }

  async wait(delay) {
    return new Promise((resolve) => setTimeout(resolve, delay));
  }

  wrapError(originalError, attempt, context) {
    const retryError = new Error(
      `Operation failed after ${attempt} attempts: ${originalError.message}`
    );

    retryError.originalError = originalError;
    retryError.attempts = attempt;
    retryError.context = context;
    retryError.isRetryError = true;

    return retryError;
  }
}

// Usage Example: API Call with Retry
async function apiCallWithRetry() {
  const retry = new AdvancedRetry({
    maxAttempts: 5,
    initialDelay: 1000,
    maxDelay: 30000,
    jitter: true,
  });

  try {
    const result = await retry.execute(
      async (attempt) => {
        // Simulate API call that might fail
        console.log(`Making API call (attempt ${attempt})...`);

        // Simulate different failure scenarios
        if (attempt < 3 && Math.random() > 0.3) {
          throw simulateTransientError();
        }

        // Success case
        return { data: "API response", status: "success" };
      },
      { operation: "api-call" }
    );

    console.log("Final result:", result);
    return result;
  } catch (error) {
    console.error("All retry attempts failed:", error.message);
    throw error;
  }
}

function simulateTransientError() {
  const errors = [
    { code: "ECONNRESET", message: "Connection reset by peer" },
    { code: "ETIMEDOUT", message: "Connection timed out" },
    { code: "ESOCKETTIMEDOUT", message: "Socket timeout" },
    { statusCode: 503, message: "Service unavailable" },
    { statusCode: 429, message: "Rate limit exceeded" },
  ];

  return errors[Math.floor(Math.random() * errors.length)];
}

// apiCallWithRetry();
```

## Strategy 2: Circuit Breaker Pattern

```javascript
// circuit-breaker.js
class CircuitBreaker {
  constructor(options = {}) {
    this.failureThreshold = options.failureThreshold || 5;
    this.successThreshold = options.successThreshold || 3;
    this.timeout = options.timeout || 10000; // 10 seconds
    this.resetTimeout = options.resetTimeout || 30000; // 30 seconds

    this.state = "CLOSED"; // CLOSED, OPEN, HALF_OPEN
    this.failureCount = 0;
    this.successCount = 0;
    this.lastFailureTime = null;
    this.nextAttemptTime = null;
  }

  async call(operation, fallback = null) {
    if (this.state === "OPEN") {
      if (Date.now() < this.nextAttemptTime) {
        console.log("Circuit breaker is OPEN, using fallback");
        return fallback ? fallback() : this.getFallbackResponse();
      }

      // Time to try again - move to HALF_OPEN
      this.state = "HALF_OPEN";
      console.log("Circuit breaker moving to HALF_OPEN state");
    }

    try {
      const result = await this.executeWithTimeout(operation);
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure(error);

      if (this.state === "OPEN") {
        return fallback ? fallback(error) : this.getFallbackResponse(error);
      }

      throw error;
    }
  }

  async executeWithTimeout(operation) {
    return new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        reject(new Error("Operation timeout"));
      }, this.timeout);

      operation()
        .then(resolve)
        .catch(reject)
        .finally(() => clearTimeout(timeoutId));
    });
  }

  onSuccess() {
    this.failureCount = 0;

    if (this.state === "HALF_OPEN") {
      this.successCount++;

      if (this.successCount >= this.successThreshold) {
        this.state = "CLOSED";
        this.successCount = 0;
        console.log("Circuit breaker moving to CLOSED state");
      }
    }
  }

  onFailure(error) {
    this.failureCount++;
    this.lastFailureTime = Date.now();

    console.warn(`Circuit breaker failure count: ${this.failureCount}`);

    if (
      this.failureCount >= this.failureThreshold ||
      this.state === "HALF_OPEN"
    ) {
      this.state = "OPEN";
      this.nextAttemptTime = Date.now() + this.resetTimeout;
      console.error(
        `Circuit breaker OPEN until ${new Date(this.nextAttemptTime)}`
      );
    }
  }

  getFallbackResponse(error = null) {
    return {
      status: "fallback",
      message: "Service unavailable - using fallback response",
      originalError: error?.message,
      timestamp: new Date().toISOString(),
    };
  }

  getStatus() {
    return {
      state: this.state,
      failureCount: this.failureCount,
      successCount: this.successCount,
      lastFailureTime: this.lastFailureTime,
      nextAttemptTime: this.nextAttemptTime,
    };
  }
}

// Usage Example: Protected API Service
class APIService {
  constructor() {
    this.circuitBreaker = new CircuitBreaker({
      failureThreshold: 3,
      successThreshold: 2,
      timeout: 5000,
      resetTimeout: 30000,
    });
  }

  async callExternalAPI(url, options = {}) {
    const operation = async () => {
      console.log(`Calling external API: ${url}`);

      // Simulate API call
      if (Math.random() < 0.6) {
        // 60% chance of failure for demo
        throw new Error("API service unavailable");
      }

      return { data: "API response", status: 200 };
    };

    const fallback = (error) => ({
      data: "Fallback data from cache",
      status: "fallback",
      error: error?.message,
    });

    return this.circuitBreaker.call(operation, fallback);
  }

  async makeMultipleCalls() {
    const results = [];

    for (let i = 0; i < 10; i++) {
      try {
        const result = await this.callExternalAPI(
          `https://api.example.com/data/${i}`
        );
        results.push(result);

        // Wait between calls to see circuit breaker in action
        await new Promise((resolve) => setTimeout(resolve, 1000));

        console.log(`Call ${i + 1}:`, this.circuitBreaker.getStatus());
      } catch (error) {
        console.error(`Call ${i + 1} failed:`, error.message);
        results.push({ error: error.message, status: "failed" });
      }
    }

    return results;
  }
}

// const apiService = new APIService();
// apiService.makeMultipleCalls().then(console.log);
```

## Strategy 3: Dead Letter Queue with Retry

```javascript
// dead-letter-queue.js
class DeadLetterQueue {
  constructor(options = {}) {
    this.maxRetries = options.maxRetries || 3;
    this.retryDelay = options.retryDelay || 5000;
    this.dlqThreshold = options.dlqThreshold || 0.1; // Move to DLQ if >10% fail
    this.storage = new Map();
    this.metrics = {
      processed: 0,
      failed: 0,
      retried: 0,
      dlq: 0,
    };
  }

  async processMessage(message, processor) {
    this.metrics.processed++;

    for (let attempt = 1; attempt <= this.maxRetries + 1; attempt++) {
      try {
        console.log(`Processing message (attempt ${attempt}):`, message.id);

        const result = await processor(message);

        if (attempt > 1) {
          this.metrics.retried++;
          console.log(
            `Message ${message.id} succeeded on retry attempt ${attempt}`
          );
        }

        return result;
      } catch (error) {
        console.error(
          `Attempt ${attempt} failed for message ${message.id}:`,
          error.message
        );

        if (attempt === this.maxRetries + 1) {
          // Final attempt failed - move to DLQ
          await this.moveToDLQ(message, error, attempt);
          this.metrics.dlq++;
          break;
        }

        // Wait before retry
        await this.wait(this.retryDelay * attempt);

        // Update message with retry context
        message.retryCount = attempt;
        message.lastError = error.message;
      }
    }
  }

  async moveToDLQ(message, error, attempt) {
    const dlqMessage = {
      ...message,
      dlqReason: error.message,
      finalAttempt: attempt,
      movedToDLQ: new Date(),
      originalTimestamp: message.timestamp || new Date(),
    };

    this.storage.set(`dlq-${message.id}`, dlqMessage);
    console.error(
      `Message ${message.id} moved to DLQ after ${attempt} attempts`
    );
  }

  async wait(delay) {
    return new Promise((resolve) => setTimeout(resolve, delay));
  }

  getDLQMessages() {
    const dlqMessages = [];
    for (const [key, value] of this.storage.entries()) {
      if (key.startsWith("dlq-")) {
        dlqMessages.push(value);
      }
    }
    return dlqMessages;
  }

  getMetrics() {
    const failureRate =
      this.metrics.processed > 0
        ? this.metrics.failed / this.metrics.processed
        : 0;

    return {
      ...this.metrics,
      failureRate: Math.round(failureRate * 100) + "%",
      dlqSize: this.getDLQMessages().length,
      shouldInvestigate: failureRate > this.dlqThreshold,
    };
  }

  async reprocessDLQ(processor, batchSize = 10) {
    const dlqMessages = this.getDLQMessages();
    const toReprocess = dlqMessages.slice(0, batchSize);

    console.log(`Reprocessing ${toReprocess.length} messages from DLQ`);

    const results = [];

    for (const message of toReprocess) {
      try {
        // Remove from DLQ before reprocessing
        this.storage.delete(`dlq-${message.id}`);

        const result = await this.processMessage(
          { ...message, fromDLQ: true },
          processor
        );

        results.push({ messageId: message.id, status: "success", result });
      } catch (error) {
        results.push({
          messageId: message.id,
          status: "failed",
          error: error.message,
        });
        // Message will go back to DLQ through normal processing
      }
    }

    return results;
  }
}

// Usage Example: Message Processor with DLQ
class MessageProcessor {
  constructor() {
    this.dlq = new DeadLetterQueue({
      maxRetries: 3,
      retryDelay: 2000,
      dlqThreshold: 0.1,
    });
  }

  async processMessages(messages) {
    const results = [];

    for (const message of messages) {
      try {
        const result = await this.dlq.processMessage(
          message,
          this.processMessage.bind(this)
        );
        results.push({ messageId: message.id, status: "success", result });
      } catch (error) {
        results.push({
          messageId: message.id,
          status: "dlq",
          error: error.message,
        });
      }
    }

    console.log("Processing metrics:", this.dlq.getMetrics());
    return results;
  }

  async processMessage(message) {
    // Simulate various processing scenarios
    console.log(
      `Processing message ${message.id} with content:`,
      message.content
    );

    // Simulate different failure conditions
    if (message.content?.includes("fail-permanently")) {
      throw new Error("Permanent failure - should go to DLQ immediately");
    }

    if (message.content?.includes("fail-transient") && !message.fromDLQ) {
      throw new Error("Transient failure - should retry");
    }

    if (Math.random() < 0.3 && !message.fromDLQ) {
      // 30% failure rate for demo
      throw new Error("Random processing failure");
    }

    // Simulate processing time
    await new Promise((resolve) => setTimeout(resolve, 500));

    return { processed: true, timestamp: new Date(), messageId: message.id };
  }
}

// Example usage
const messages = [
  { id: "msg-1", content: "Normal message", timestamp: new Date() },
  {
    id: "msg-2",
    content: "Message with fail-transient error",
    timestamp: new Date(),
  },
  {
    id: "msg-3",
    content: "Message with fail-permanently error",
    timestamp: new Date(),
  },
  { id: "msg-4", content: "Another normal message", timestamp: new Date() },
];

// const processor = new MessageProcessor();
// processor.processMessages(messages).then(console.log);
```

## Strategy 4: Comprehensive Retry Orchestrator

```javascript
// retry-orchestrator.js
class RetryOrchestrator {
  constructor() {
    this.strategies = new Map();
    this.initializeDefaultStrategies();
  }

  initializeDefaultStrategies() {
    // Immediate retry strategy
    this.addStrategy("immediate", {
      maxAttempts: 3,
      calculateDelay: (attempt) => 0, // No delay
      shouldRetry: (error) => this.isNetworkError(error),
    });

    // Exponential backoff strategy
    this.addStrategy("exponential", {
      maxAttempts: 5,
      calculateDelay: (attempt) =>
        Math.min(1000 * Math.pow(2, attempt - 1), 30000),
      shouldRetry: (error) => this.isRetryableError(error),
    });

    // Progressive strategy (increasing delays)
    this.addStrategy("progressive", {
      maxAttempts: 4,
      calculateDelay: (attempt) =>
        [1000, 5000, 15000, 30000][attempt - 1] || 30000,
      shouldRetry: (error) => true, // Retry all errors
    });
  }

  addStrategy(name, config) {
    this.strategies.set(name, {
      maxAttempts: config.maxAttempts,
      calculateDelay: config.calculateDelay,
      shouldRetry: config.shouldRetry || (() => true),
      onRetry: config.onRetry || (() => {}),
      onFailure: config.onFailure || (() => {}),
    });
  }

  async executeWithRetry(
    operation,
    strategyName = "exponential",
    context = {}
  ) {
    const strategy = this.strategies.get(strategyName);
    if (!strategy) throw new Error(`Strategy ${strategyName} not found`);

    let lastError;
    const startTime = Date.now();

    for (let attempt = 1; attempt <= strategy.maxAttempts; attempt++) {
      try {
        context.attempt = attempt;
        context.strategy = strategyName;

        const result = await operation(attempt, context);

        if (attempt > 1) {
          console.log(`✅ Operation succeeded on attempt ${attempt}`);
        }

        return {
          result,
          attempts: attempt,
          duration: Date.now() - startTime,
          strategy: strategyName,
        };
      } catch (error) {
        lastError = error;

        // Check if we should retry
        const shouldRetry =
          attempt < strategy.maxAttempts && strategy.shouldRetry(error);

        if (!shouldRetry) {
          strategy.onFailure(error, attempt, context);
          throw this.wrapFinalError(error, attempt, strategyName, context);
        }

        // Calculate delay and wait
        const delay = strategy.calculateDelay(attempt);
        console.warn(
          `⚠️ Attempt ${attempt} failed: ${error.message}. Retrying in ${delay}ms...`
        );

        strategy.onRetry(error, attempt, delay, context);
        await this.wait(delay);

        // Update context for next attempt
        context.lastError = error.message;
        context.lastAttempt = attempt;
      }
    }

    throw lastError;
  }

  isNetworkError(error) {
    const networkErrorCodes = ["ECONNRESET", "ETIMEDOUT", "ECONNREFUSED"];
    return networkErrorCodes.includes(error.code);
  }

  isRetryableError(error) {
    // HTTP status codes that are typically retryable
    if (error.statusCode) {
      return [408, 429, 500, 502, 503, 504].includes(error.statusCode);
    }

    // Database connection errors
    if (error.code && error.code.startsWith("ECONN")) {
      return true;
    }

    // Timeout errors
    if (error.code === "ETIMEDOUT" || error.message?.includes("timeout")) {
      return true;
    }

    return false;
  }

  wrapFinalError(error, attempts, strategy, context) {
    const finalError = new Error(
      `Operation failed after ${attempts} attempts using ${strategy} strategy: ${error.message}`
    );

    finalError.originalError = error;
    finalError.attempts = attempts;
    finalError.strategy = strategy;
    finalError.context = context;
    finalError.timestamp = new Date();

    return finalError;
  }

  async wait(delay) {
    return new Promise((resolve) => setTimeout(resolve, delay));
  }

  // Batch processing with individual retry strategies
  async processBatch(operations, strategyName = "exponential") {
    const results = await Promise.allSettled(
      operations.map((operation, index) =>
        this.executeWithRetry(operation, strategyName, { batchIndex: index })
      )
    );

    return this.analyzeBatchResults(results);
  }

  analyzeBatchResults(results) {
    const analysis = {
      total: results.length,
      successful: 0,
      failed: 0,
      totalAttempts: 0,
      failures: [],
    };

    results.forEach((result, index) => {
      if (result.status === "fulfilled") {
        analysis.successful++;
        analysis.totalAttempts += result.value.attempts;
      } else {
        analysis.failed++;
        analysis.failures.push({
          index,
          error: result.reason.message,
          attempts: result.reason.attempts,
        });
      }
    });

    return analysis;
  }
}

// Usage Example: Comprehensive Retry System
async function demonstrateRetryStrategies() {
  const orchestrator = new RetryOrchestrator();

  // Add custom strategy
  orchestrator.addStrategy("custom-backoff", {
    maxAttempts: 4,
    calculateDelay: (attempt) =>
      Math.min(500 * Math.pow(3, attempt - 1), 20000),
    shouldRetry: (error) => !error.message.includes("permanent"),
    onRetry: (error, attempt, delay) => {
      console.log(`🔄 Custom retry: Attempt ${attempt}, Delay ${delay}ms`);
    },
  });

  // Test different scenarios
  const testScenarios = [
    {
      name: "Network Operation",
      operation: async (attempt) => {
        if (attempt < 3)
          throw { code: "ECONNRESET", message: "Connection reset" };
        return { data: "Network response" };
      },
      strategy: "immediate",
    },
    {
      name: "API Call",
      operation: async (attempt) => {
        if (attempt < 4)
          throw { statusCode: 503, message: "Service unavailable" };
        return { data: "API response" };
      },
      strategy: "exponential",
    },
    {
      name: "Permanent Failure",
      operation: async () => {
        throw new Error("Permanent failure - should not retry");
      },
      strategy: "custom-backoff",
    },
  ];

  for (const scenario of testScenarios) {
    console.log(`\n🧪 Testing: ${scenario.name}`);
    console.log(`Strategy: ${scenario.strategy}`);

    try {
      const result = await orchestrator.executeWithRetry(
        scenario.operation,
        scenario.strategy,
        { scenario: scenario.name }
      );

      console.log(`✅ Success:`, result);
    } catch (error) {
      console.error(`❌ Final failure:`, error.message);
      console.log(`Attempts: ${error.attempts}, Strategy: ${error.strategy}`);
    }
  }
}

// demonstrateRetryStrategies();
```

## Key Takeaways:

1. **Exponential Backoff**: Prevents overwhelming failing systems
2. **Jitter**: Avoids thundering herd problems
3. **Circuit Breaker**: Protects systems from cascading failures
4. **Dead Letter Queue**: Handles permanently failing messages
5. **Strategy Pattern**: Allows different retry approaches for different scenarios

These strategies ensure your long-running processes can handle failures gracefully and maintain system reliability!
