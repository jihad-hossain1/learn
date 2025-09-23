**Durable Execution Patterns** for long-running processes. This is a powerful concept where workflows can survive process failures, infrastructure issues, and continue execution seamlessly.

## What is Durable Execution?

**Durable Execution** ensures that once a workflow starts, it will complete exactly once, even in the face of:

- Process crashes
- Infrastructure failures
- Network issues
- Code deployments

## Key Patterns in Durable Execution

### 1. **Event Sourcing Pattern**

### 2. **Saga Pattern**

### 3. **Workflow Pattern**

### 4. **Compensating Actions Pattern**

Let me show you JavaScript examples for each:

## Pattern 1: Event Sourcing with Durable Execution

```javascript
// event-store.js - Event sourcing foundation
class EventStore {
  constructor() {
    this.events = [];
    this.subscribers = [];
  }

  async appendEvent(aggregateId, eventType, eventData) {
    const event = {
      id: this.events.length + 1,
      aggregateId,
      type: eventType,
      data: eventData,
      timestamp: new Date(),
      version: this.getNextVersion(aggregateId),
    };

    this.events.push(event);

    // Notify subscribers (durably)
    await this.notifySubscribers(event);

    return event;
  }

  getEvents(aggregateId) {
    return this.events.filter((e) => e.aggregateId === aggregateId);
  }

  async notifySubscribers(event) {
    for (const subscriber of this.subscribers) {
      try {
        await subscriber(event);
      } catch (error) {
        console.error(`Subscriber failed: ${error.message}`);
        // Retry logic would go here in production
      }
    }
  }

  getNextVersion(aggregateId) {
    const aggregateEvents = this.getEvents(aggregateId);
    return aggregateEvents.length + 1;
  }
}

// durable-workflow.js
class DurableWorkflow {
  constructor(eventStore) {
    this.eventStore = eventStore;
    this.handlers = new Map();
    this.workflowStates = new Map();
  }

  registerHandler(eventType, handler) {
    this.handlers.set(eventType, handler);
  }

  async startWorkflow(workflowId, initialEvent) {
    // Store initial state
    this.workflowStates.set(workflowId, {
      status: "running",
      currentStep: 0,
      data: {},
      createdAt: new Date(),
    });

    await this.eventStore.appendEvent(
      workflowId,
      "WORKFLOW_STARTED",
      initialEvent
    );
  }

  async handleEvent(event) {
    const workflowId = event.aggregateId;
    const handler = this.handlers.get(event.type);

    if (!handler) return;

    // Recover workflow state from events
    const currentState = await this.recoverState(workflowId);

    try {
      // Execute handler with current state
      const result = await handler(event, currentState);

      // Update state durably
      await this.updateWorkflowState(workflowId, result);
    } catch (error) {
      await this.eventStore.appendEvent(workflowId, "WORKFLOW_FAILED", {
        error: error.message,
        step: event.type,
      });
    }
  }

  async recoverState(workflowId) {
    const events = this.eventStore.getEvents(workflowId);
    let state = this.workflowStates.get(workflowId) || {
      status: "initializing",
      currentStep: 0,
      data: {},
    };

    // Replay events to rebuild state (this is the durable part!)
    for (const event of events) {
      const handler = this.handlers.get(event.type);
      if (handler) {
        state = await handler(event, state, true); // true indicates replay
      }
    }

    return state;
  }

  async updateWorkflowState(workflowId, updates) {
    const currentState = this.workflowStates.get(workflowId) || {};
    const newState = { ...currentState, ...updates };
    this.workflowStates.set(workflowId, newState);
  }
}

// Usage example: E-commerce Order Saga
const eventStore = new EventStore();
const workflow = new DurableWorkflow(eventStore);

// Register durable handlers
workflow.registerHandler("ORDER_CREATED", async (event, state, isReplay) => {
  if (!isReplay) {
    console.log("Processing order creation...");
    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 1000));
  }

  return {
    ...state,
    currentStep: 1,
    data: { ...state.data, order: event.data },
  };
});

workflow.registerHandler(
  "PAYMENT_PROCESSED",
  async (event, state, isReplay) => {
    if (!isReplay) {
      console.log("Processing payment...");
      await new Promise((resolve) => setTimeout(resolve, 2000));
    }

    return {
      ...state,
      currentStep: 2,
      data: { ...state.data, payment: event.data },
    };
  }
);

// Subscribe to events
eventStore.subscribers.push(async (event) => {
  await workflow.handleEvent(event);
});

// Start a workflow
await workflow.startWorkflow("order-123", {
  orderId: "order-123",
  items: [{ productId: "prod-1", quantity: 2 }],
  total: 99.99,
});
```

## Pattern 2: Saga Pattern with Compensating Actions

```javascript
// saga-orchestrator.js
class SagaOrchestrator {
  constructor() {
    this.sagas = new Map();
    this.executionLog = [];
  }

  defineSaga(sagaId, steps) {
    this.sagas.set(sagaId, steps);
  }

  async executeSaga(sagaId, initialData) {
    const saga = this.sagas.get(sagaId);
    if (!saga) throw new Error(`Saga ${sagaId} not found`);

    const executionId = `${sagaId}-${Date.now()}`;
    const context = {
      executionId,
      currentStep: 0,
      data: initialData,
      compensated: false,
    };

    // Log start (durable storage)
    await this.logExecution(executionId, "STARTED", context);

    try {
      for (let i = 0; i < saga.length; i++) {
        const step = saga[i];
        context.currentStep = i;

        // Log step start
        await this.logExecution(executionId, `STEP_${i}_STARTED`, context);

        // Execute step
        const result = await this.executeStep(step, context.data);
        context.data = { ...context.data, ...result };

        // Log step completion
        await this.logExecution(executionId, `STEP_${i}_COMPLETED`, context);
      }

      await this.logExecution(executionId, "COMPLETED", context);
      return context.data;
    } catch (error) {
      await this.logExecution(executionId, "FAILED", {
        ...context,
        error: error.message,
      });

      // Execute compensation
      await this.compensateSaga(executionId, saga, context);
      throw error;
    }
  }

  async executeStep(step, data) {
    // Simulate potential failure
    if (Math.random() < 0.3) {
      // 30% chance of failure
      throw new Error(`Step ${step.name} failed randomly`);
    }

    console.log(`Executing step: ${step.name}`);
    await new Promise((resolve) => setTimeout(resolve, 1000));

    return await step.execute(data);
  }

  async compensateSaga(executionId, saga, context) {
    if (context.compensated) return;

    console.log(`Compensating saga: ${executionId}`);

    // Execute compensation in reverse order
    for (let i = context.currentStep; i >= 0; i--) {
      const step = saga[i];
      if (step.compensate) {
        try {
          await step.compensate(context.data);
          await this.logExecution(
            executionId,
            `STEP_${i}_COMPENSATED`,
            context
          );
        } catch (compError) {
          console.error(`Compensation failed for step ${i}:`, compError);
        }
      }
    }

    context.compensated = true;
    await this.logExecution(executionId, "COMPENSATED", context);
  }

  async logExecution(executionId, status, context) {
    const logEntry = {
      executionId,
      status,
      timestamp: new Date(),
      context: JSON.parse(JSON.stringify(context)), // Deep clone
    };

    this.executionLog.push(logEntry);
    console.log(`[${executionId}] ${status}`);
  }

  async recoverSaga(executionId) {
    const logs = this.executionLog.filter(
      (log) => log.executionId === executionId
    );
    if (logs.length === 0)
      throw new Error(`Execution ${executionId} not found`);

    const lastLog = logs[logs.length - 1];

    if (lastLog.status === "FAILED" && !lastLog.context.compensated) {
      console.log(`Recovering failed saga: ${executionId}`);
      const saga = this.sagas.get(executionId.split("-")[0]);
      await this.compensateSaga(executionId, saga, lastLog.context);
    }
  }
}

// Order Processing Saga Example
const orchestrator = new SagaOrchestrator();

// Define saga steps with compensation
orchestrator.defineSaga("order-processing", [
  {
    name: "Validate Order",
    execute: async (data) => {
      console.log("Validating order...");
      if (!data.orderId) throw new Error("Order ID required");
      return { validated: true };
    },
    compensate: async (data) => {
      console.log("Compensating order validation...");
      // Nothing to compensate for validation
    },
  },
  {
    name: "Process Payment",
    execute: async (data) => {
      console.log("Processing payment...");
      // Simulate payment processing
      const paymentResult = {
        paymentId: "pay-" + Date.now(),
        status: "completed",
      };
      return { payment: paymentResult };
    },
    compensate: async (data) => {
      console.log("Refunding payment...");
      // Call refund API
      return { payment: { ...data.payment, status: "refunded" } };
    },
  },
  {
    name: "Update Inventory",
    execute: async (data) => {
      console.log("Updating inventory...");
      // Simulate inventory update
      return { inventoryUpdated: true };
    },
    compensate: async (data) => {
      console.log("Reverting inventory...");
      // Restore inventory levels
      return { inventoryUpdated: false };
    },
  },
  {
    name: "Send Confirmation",
    execute: async (data) => {
      console.log("Sending confirmation...");
      // Send email/SMS
      return { confirmationSent: true };
    },
    compensate: async (data) => {
      console.log("Sending failure notification...");
      // Notify user of failure
      return { failureNotified: true };
    },
  },
]);

// Execute the saga
async function runOrderSaga() {
  try {
    const result = await orchestrator.executeSaga("order-processing", {
      orderId: "order-123",
      amount: 99.99,
      items: [{ productId: "prod-1", quantity: 2 }],
    });

    console.log("Saga completed successfully:", result);
  } catch (error) {
    console.error("Saga failed:", error.message);
  }
}

// runOrderSaga();
```

## Pattern 3: Durable Workflow with Checkpointing

```javascript
// durable-workflow-engine.js
class DurableWorkflowEngine {
  constructor(storage) {
    this.storage = storage;
    this.workflows = new Map();
  }

  registerWorkflow(name, definition) {
    this.workflows.set(name, definition);
  }

  async startWorkflow(name, input, workflowId = null) {
    const workflowId = workflowId || `${name}-${Date.now()}`;
    const definition = this.workflows.get(name);

    if (!definition) {
      throw new Error(`Workflow ${name} not found`);
    }

    const execution = {
      workflowId,
      name,
      status: "running",
      currentStep: 0,
      input,
      output: null,
      checkpoints: [],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    // Save initial state
    await this.saveExecution(execution);

    // Start execution
    return this.executeWorkflow(execution, definition);
  }

  async executeWorkflow(execution, definition) {
    try {
      for (let i = execution.currentStep; i < definition.steps.length; i++) {
        const step = definition.steps[i];
        execution.currentStep = i;
        execution.updatedAt = new Date();

        console.log(`Executing step ${i}: ${step.name}`);

        // Create checkpoint before executing step
        await this.createCheckpoint(execution);

        // Execute step
        const result = await step.execute(execution);
        execution.data = { ...execution.data, ...result };

        // Update checkpoint after successful step
        execution.checkpoints[i] = {
          status: "completed",
          timestamp: new Date(),
          result,
        };

        await this.saveExecution(execution);
      }

      execution.status = "completed";
      execution.output = execution.data;
      await this.saveExecution(execution);

      return execution.output;
    } catch (error) {
      execution.status = "failed";
      execution.error = error.message;
      await this.saveExecution(execution);

      throw error;
    }
  }

  async createCheckpoint(execution) {
    const checkpoint = {
      step: execution.currentStep,
      data: JSON.parse(JSON.stringify(execution.data)), // Deep clone
      timestamp: new Date(),
    };

    execution.checkpoints[execution.currentStep] = checkpoint;
    await this.saveExecution(execution);
  }

  async recoverWorkflow(workflowId) {
    const execution = await this.storage.get(workflowId);

    if (!execution) {
      throw new Error(`Workflow ${workflowId} not found`);
    }

    if (execution.status === "running") {
      console.log(
        `Recovering workflow ${workflowId} from step ${execution.currentStep}`
      );
      const definition = this.workflows.get(execution.name);
      return this.executeWorkflow(execution, definition);
    }

    return execution;
  }

  async saveExecution(execution) {
    // In production, this would save to durable storage
    await this.storage.set(execution.workflowId, execution);
  }
}

// Example: Document Processing Workflow
const engine = new DurableWorkflowEngine({
  storage: new Map(), // In production, use Redis/Database
});

// Define a durable workflow
engine.registerWorkflow("document-processing", {
  steps: [
    {
      name: "Upload Document",
      execute: async (execution) => {
        await simulateProcess("Uploading document", 2000);
        return { uploaded: true, documentId: "doc-" + Date.now() };
      },
    },
    {
      name: "Validate Format",
      execute: async (execution) => {
        await simulateProcess("Validating format", 1500);
        if (!execution.data.documentId) throw new Error("No document ID");
        return { formatValid: true };
      },
    },
    {
      name: "Extract Text",
      execute: async (execution) => {
        await simulateProcess("Extracting text", 3000);
        return { text: "Extracted document content..." };
      },
    },
    {
      name: "Analyze Content",
      execute: async (execution) => {
        await simulateProcess("Analyzing content", 2500);
        return {
          analysis: { keywords: ["important", "data"], sentiment: "positive" },
        };
      },
    },
    {
      name: "Generate Summary",
      execute: async (execution) => {
        await simulateProcess("Generating summary", 2000);
        return { summary: "This document contains important information..." };
      },
    },
  ],
});

async function simulateProcess(name, duration) {
  console.log(`${name}...`);
  await new Promise((resolve) => setTimeout(resolve, duration));
  // Simulate occasional failures
  if (Math.random() < 0.2) {
    throw new Error(`${name} failed randomly`);
  }
  console.log(`${name} completed`);
}

// Usage example
async function runDurableWorkflow() {
  try {
    const result = await engine.startWorkflow("document-processing", {
      filename: "report.pdf",
      size: 1024000,
    });

    console.log("Workflow completed:", result);
  } catch (error) {
    console.error("Workflow failed:", error.message);

    // The workflow can be recovered later!
    console.log("Workflow can be recovered using engine.recoverWorkflow()");
  }
}

// runDurableWorkflow();
```

## Key Benefits of Durable Execution Patterns:

1. **Fault Tolerance**: Survives process crashes and infrastructure failures
2. **Exactly-Once Execution**: Guarantees completion even after failures
3. **Debuggability**: Full audit trail of execution
4. **Recovery**: Can resume from last known good state
5. **Scalability**: Can handle long-running processes (days/months)

## When to Use Each Pattern:

- **Event Sourcing**: When you need full auditability and time travel
- **Saga Pattern**: For distributed transactions requiring compensation
- **Workflow Pattern**: For complex business processes with multiple steps
- **Compensating Actions**: When you need to undo partial completions

These patterns are fundamental to building reliable, long-running systems that can handle real-world failures gracefully!
