# Subgraphs & Nested Workflows

Subgraphs and nested workflows are powerful concepts in workflow orchestration that allow you to create modular, reusable, and complex workflow structures by breaking down large processes into smaller, manageable components.

## Core Concepts

### 1. **Subgraphs**
- Self-contained workflow components
- Can be reused across multiple parent workflows
- Have defined inputs and outputs
- Can be executed independently or as part of larger workflows

### 2. **Nested Workflows**
- Workflows that contain other workflows
- Enable hierarchical organization of complex processes
- Support recursion and complex dependency management
- Improve maintainability and reusability

## JavaScript Implementation Examples

### Basic Subgraph Implementation

```javascript
// Base Subgraph Class
class Subgraph {
  constructor(name, description = '') {
    this.name = name;
    this.description = description;
    this.nodes = new Map();
    this.edges = [];
    this.inputs = new Map();
    this.outputs = new Map();
    this.executionHistory = [];
  }

  addNode(nodeId, nodeFunction, config = {}) {
    this.nodes.set(nodeId, {
      id: nodeId,
      execute: nodeFunction,
      config,
      status: 'pending',
      dependencies: config.dependencies || []
    });
    return this;
  }

  addEdge(fromNode, toNode, condition = null) {
    this.edges.push({
      from: fromNode,
      to: toNode,
      condition,
      type: 'dependency'
    });
    return this;
  }

  defineInput(inputName, validator = null) {
    this.inputs.set(inputName, {
      name: inputName,
      validator,
      required: true
    });
    return this;
  }

  defineOutput(outputName, generator = null) {
    this.outputs.set(outputName, {
      name: outputName,
      generator,
      required: true
    });
    return this;
  }

  async validateInputs(inputData) {
    const errors = [];
    
    for (const [inputName, inputConfig] of this.inputs) {
      if (inputConfig.required && !(inputName in inputData)) {
        errors.push(`Missing required input: ${inputName}`);
      }
      
      if (inputConfig.validator && inputName in inputData) {
        try {
          await inputConfig.validator(inputData[inputName]);
        } catch (error) {
          errors.push(`Invalid input ${inputName}: ${error.message}`);
        }
      }
    }
    
    return errors;
  }

  async execute(inputData = {}, context = {}) {
    console.log(`\n=== Executing Subgraph: ${this.name} ===`);
    
    // Validate inputs
    const validationErrors = await this.validateInputs(inputData);
    if (validationErrors.length > 0) {
      throw new Error(`Input validation failed: ${validationErrors.join(', ')}`);
    }

    const executionId = `exec_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    const executionContext = {
      id: executionId,
      startTime: new Date(),
      inputData,
      nodeResults: new Map(),
      status: 'running',
      ...context
    };

    try {
      // Topological sort for execution order
      const executionOrder = this.calculateExecutionOrder();
      
      for (const nodeId of executionOrder) {
        const node = this.nodes.get(nodeId);
        if (!node) continue;

        console.log(`Executing node: ${nodeId}`);
        
        // Check dependencies
        const dependenciesMet = node.dependencies.every(depId => 
          executionContext.nodeResults.has(depId)
        );

        if (!dependenciesMet) {
          throw new Error(`Dependencies not met for node: ${nodeId}`);
        }

        // Prepare node input
        const nodeInput = this.prepareNodeInput(nodeId, executionContext);
        
        // Execute node
        node.status = 'running';
        const result = await node.execute(nodeInput, executionContext);
        
        // Store result
        executionContext.nodeResults.set(nodeId, {
          nodeId,
          result,
          timestamp: new Date(),
          status: 'completed'
        });
        
        node.status = 'completed';
      }

      // Generate outputs
      const outputs = this.generateOutputs(executionContext);
      
      executionContext.status = 'completed';
      executionContext.endTime = new Date();
      executionContext.outputs = outputs;

      this.executionHistory.push(executionContext);
      
      console.log(`✓ Subgraph ${this.name} completed successfully`);
      return outputs;

    } catch (error) {
      executionContext.status = 'failed';
      executionContext.error = error;
      executionContext.endTime = new Date();
      
      this.executionHistory.push(executionContext);
      console.error(`✗ Subgraph ${this.name} failed: ${error.message}`);
      throw error;
    }
  }

  calculateExecutionOrder() {
    // Simple topological sort implementation
    const visited = new Set();
    const order = [];
    
    const visit = (nodeId) => {
      if (visited.has(nodeId)) return;
      visited.add(nodeId);
      
      const node = this.nodes.get(nodeId);
      node.dependencies.forEach(depId => visit(depId));
      
      order.push(nodeId);
    };
    
    this.nodes.forEach((_, nodeId) => visit(nodeId));
    return order;
  }

  prepareNodeInput(nodeId, context) {
    const node = this.nodes.get(nodeId);
    const input = { ...context.inputData };
    
    // Add results from dependency nodes
    node.dependencies.forEach(depId => {
      const depResult = context.nodeResults.get(depId);
      if (depResult) {
        input[depId] = depResult.result;
      }
    });
    
    return input;
  }

  generateOutputs(context) {
    const outputs = {};
    
    for (const [outputName, outputConfig] of this.outputs) {
      if (outputConfig.generator) {
        outputs[outputName] = outputConfig.generator(context);
      } else {
        // Default output generation: use the last node's result
        const lastNodeId = Array.from(this.nodes.keys()).pop();
        outputs[outputName] = context.nodeResults.get(lastNodeId)?.result;
      }
    }
    
    return outputs;
  }
}

// Nested Workflow Orchestrator
class NestedWorkflowOrchestrator {
  constructor() {
    this.subgraphs = new Map();
    this.workflows = new Map();
    this.executionStack = [];
  }

  registerSubgraph(name, subgraph) {
    this.subgraphs.set(name, subgraph);
    return this;
  }

  async executeWorkflow(workflowName, inputData = {}) {
    const workflow = this.workflows.get(workflowName);
    if (!workflow) {
      throw new Error(`Workflow not found: ${workflowName}`);
    }

    this.executionStack.push({
      workflow: workflowName,
      startTime: new Date(),
      depth: this.executionStack.length
    });

    try {
      const result = await workflow.execute(inputData, {
        orchestrator: this,
        parentStack: [...this.executionStack]
      });

      this.executionStack.pop();
      return result;

    } catch (error) {
      this.executionStack.pop();
      throw error;
    }
  }

  async executeSubgraph(subgraphName, inputData = {}, context = {}) {
    const subgraph = this.subgraphs.get(subgraphName);
    if (!subgraph) {
      throw new Error(`Subgraph not found: ${subgraphName}`);
    }

    return await subgraph.execute(inputData, {
      orchestrator: this,
      ...context
    });
  }
}
```

### Practical Example: E-commerce Order Processing System

```javascript
// Create specialized subgraphs

// 1. Inventory Management Subgraph
const inventorySubgraph = new Subgraph('inventory_management', 'Handles inventory checks and updates');

inventorySubgraph
  .defineInput('productId', (id) => {
    if (!id || typeof id !== 'string') {
      throw new Error('Product ID must be a non-empty string');
    }
  })
  .defineInput('quantity', (qty) => {
    if (qty <= 0) {
      throw new Error('Quantity must be positive');
    }
  })
  .defineOutput('inventoryStatus')
  .defineOutput('reservationId');

// Add nodes to inventory subgraph
inventorySubgraph
  .addNode('check_availability', async (input) => {
    // Simulate database check
    await new Promise(resolve => setTimeout(resolve, 200));
    
    const available = Math.random() > 0.1; // 90% available
    return {
      available,
      productId: input.productId,
      requestedQty: input.quantity
    };
  })
  .addNode('reserve_inventory', async (input, context) => {
    const availability = context.nodeResults.get('check_availability').result;
    
    if (!availability.available) {
      throw new Error(`Product ${availability.productId} is out of stock`);
    }

    // Simulate reservation
    await new Promise(resolve => setTimeout(resolve, 300));
    
    return {
      reservationId: `res_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      productId: availability.productId,
      quantity: availability.requestedQty,
      reservedAt: new Date()
    };
  }, { dependencies: ['check_availability'] });

// 2. Payment Processing Subgraph
const paymentSubgraph = new Subgraph('payment_processing', 'Handles payment transactions');

paymentSubgraph
  .defineInput('amount', (amt) => {
    if (amt <= 0) throw new Error('Amount must be positive');
  })
  .defineInput('paymentMethod')
  .defineOutput('transactionId')
  .defineOutput('paymentStatus');

paymentSubgraph
  .addNode('validate_payment', async (input) => {
    await new Promise(resolve => setTimeout(resolve, 150));
    
    const valid = Math.random() > 0.05; // 95% valid
    return {
      valid,
      amount: input.amount,
      method: input.paymentMethod
    };
  })
  .addNode('process_payment', async (input, context) => {
    const validation = context.nodeResults.get('validate_payment').result;
    
    if (!validation.valid) {
      throw new Error('Payment method validation failed');
    }

    await new Promise(resolve => setTimeout(resolve, 500));
    
    // Simulate payment processing
    const success = Math.random() > 0.1; // 90% success rate
    
    return {
      transactionId: `txn_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
      status: success ? 'completed' : 'failed',
      amount: validation.amount,
      processedAt: new Date()
    };
  }, { dependencies: ['validate_payment'] });

// 3. Shipping Subgraph
const shippingSubgraph = new Subgraph('shipping', 'Handles order shipping and tracking');

shippingSubgraph
  .defineInput('orderId')
  .defineInput('address')
  .defineOutput('trackingNumber')
  .defineOutput('shippingStatus');

shippingSubgraph
  .addNode('calculate_shipping', async (input) => {
    await new Promise(resolve => setTimeout(resolve, 200));
    
    return {
      cost: Math.random() * 20 + 5, // $5-$25 shipping
      estimatedDays: Math.floor(Math.random() * 7) + 2, // 2-8 days
      carrier: ['UPS', 'FedEx', 'USPS'][Math.floor(Math.random() * 3)]
    };
  })
  .addNode('create_shipment', async (input, context) => {
    const shippingInfo = context.nodeResults.get('calculate_shipping').result;
    
    await new Promise(resolve => setTimeout(resolve, 400));
    
    return {
      trackingNumber: `TRK${Date.now()}${Math.random().toString(36).substr(2, 6).toUpperCase()}`,
      carrier: shippingInfo.carrier,
      estimatedDelivery: new Date(Date.now() + shippingInfo.estimatedDays * 24 * 60 * 60 * 1000),
      cost: shippingInfo.cost
    };
  }, { dependencies: ['calculate_shipping'] });

// Main Order Processing Workflow (Nested)
class OrderProcessingWorkflow extends Subgraph {
  constructor() {
    super('order_processing', 'Complete order processing workflow');
    
    this.orchestrator = new NestedWorkflowOrchestrator();
    this.orchestrator.registerSubgraph('inventory', inventorySubgraph);
    this.orchestrator.registerSubgraph('payment', paymentSubgraph);
    this.orchestrator.registerSubgraph('shipping', shippingSubgraph);
  }

  async execute(inputData = {}, context = {}) {
    console.log('\n🎯 Starting Order Processing Workflow');
    
    try {
      // Step 1: Process Inventory
      console.log('\n📦 Processing Inventory...');
      const inventoryResult = await this.orchestrator.executeSubgraph('inventory', {
        productId: inputData.productId,
        quantity: inputData.quantity
      });

      // Step 2: Process Payment
      console.log('\n💳 Processing Payment...');
      const paymentResult = await this.orchestrator.executeSubgraph('payment', {
        amount: inputData.amount,
        paymentMethod: inputData.paymentMethod
      });

      if (paymentResult.paymentStatus !== 'completed') {
        throw new Error('Payment processing failed');
      }

      // Step 3: Arrange Shipping
      console.log('\n🚚 Arranging Shipping...');
      const shippingResult = await this.orchestrator.executeSubgraph('shipping', {
        orderId: inputData.orderId,
        address: inputData.shippingAddress
      });

      // Final: Complete Order
      const orderResult = {
        orderId: inputData.orderId,
        status: 'completed',
        inventory: inventoryResult,
        payment: paymentResult,
        shipping: shippingResult,
        completedAt: new Date()
      };

      console.log('\n✅ Order Processing Completed Successfully!');
      return orderResult;

    } catch (error) {
      console.error('\n❌ Order Processing Failed:', error.message);
      
      // Implement compensation logic (rollback)
      await this.compensate(inputData);
      throw error;
    }
  }

  async compensate(inputData) {
    console.log('🔄 Executing compensation logic...');
    // Implement rollback logic for failed transactions
    // This would typically involve calling reverse APIs
  }
}

// Advanced Example: Recursive Workflow (Factorial Calculation)
class RecursiveWorkflow extends Subgraph {
  constructor() {
    super('recursive_factorial', 'Calculate factorial using recursive workflow');
  }

  async execute(inputData = {}, context = {}) {
    const n = inputData.number;
    
    // Base case
    if (n <= 1) {
      return { result: 1, steps: 1 };
    }

    // Recursive case
    const recursiveResult = await this.execute({
      number: n - 1
    }, { ...context, depth: (context.depth || 0) + 1 });

    const result = n * recursiveResult.result;
    
    return {
      result,
      steps: recursiveResult.steps + 1,
      depth: context.depth || 0
    };
  }
}

// Demonstration
async function demonstrateNestedWorkflows() {
  console.log('=== Nested Workflows & Subgraphs Demonstration ===\n');

  // Create and test the order processing workflow
  const orderWorkflow = new OrderProcessingWorkflow();

  const orderData = {
    orderId: 'ORD_12345',
    productId: 'PROD_67890',
    quantity: 2,
    amount: 99.99,
    paymentMethod: 'credit_card',
    shippingAddress: {
      street: '123 Main St',
      city: 'Anytown',
      state: 'CA',
      zipCode: '12345'
    }
  };

  try {
    const result = await orderWorkflow.execute(orderData);
    console.log('\n📊 Final Order Result:');
    console.log(JSON.stringify(result, null, 2));
  } catch (error) {
    console.error('Workflow execution failed:', error.message);
  }

  // Demonstrate recursive workflow
  console.log('\n=== Recursive Workflow Example ===');
  const recursiveWorkflow = new RecursiveWorkflow();
  
  try {
    const factorialResult = await recursiveWorkflow.execute({ number: 5 });
    console.log(`Factorial of 5: ${factorialResult.result}`);
    console.log(`Calculation steps: ${factorialResult.steps}`);
  } catch (error) {
    console.error('Recursive workflow failed:', error.message);
  }
}

// Run demonstration
demonstrateNestedWorkflows().catch(console.error);
```

## Key Benefits of Subgraphs & Nested Workflows

### 1. **Modularity**
- Break down complex processes into manageable components
- Reuse common patterns across different workflows

### 2. **Maintainability**
- Isolate changes to specific subgraphs
- Easier testing and debugging of individual components

### 3. **Scalability**
- Distribute workload across different systems
- Parallel execution of independent subgraphs

### 4. **Reusability**
- Share common subgraphs across multiple workflows
- Build libraries of standardized components

### 5. **Complexity Management**
- Handle recursive patterns
- Manage complex dependencies efficiently

## Advanced Patterns

### 1. **Dynamic Subgraph Composition**
```javascript
class DynamicWorkflow extends Subgraph {
  async execute(inputData, context) {
    // Dynamically compose subgraphs based on input
    const requiredSubgraphs = this.analyzeRequirements(inputData);
    
    const results = {};
    for (const subgraphName of requiredSubgraphs) {
      results[subgraphName] = await context.orchestrator.executeSubgraph(
        subgraphName, 
        inputData
      );
    }
    
    return this.synthesizeResults(results);
  }
}
```

### 2. **Conditional Nesting**
```javascript
class ConditionalWorkflow extends Subgraph {
  async execute(inputData, context) {
    if (inputData.complexity === 'high') {
      // Use nested workflow for complex cases
      return await context.orchestrator.executeSubgraph('complex_processing', inputData);
    } else {
      // Use simple processing for straightforward cases
      return await this.simpleProcessing(inputData);
    }
  }
}
```

This architecture is particularly valuable in enterprise systems, microservices architectures, and complex business process automation where workflows need to be modular, reusable, and maintainable.