# Message Passing Patterns in LangGraph

Message passing is the fundamental communication mechanism between nodes in a graph. It enables data flow, state management, and coordination between different components of your LangGraph application.

## Core Concepts

### **What is Message Passing?**

- **Inter-node communication**: How nodes exchange data and signals
- **State propagation**: Passing state from one node to another
- **Event-driven architecture**: Nodes react to incoming messages
- **Decoupled components**: Nodes don't need to know about each other directly

### **Key Message Passing Patterns**

1. **Direct Passing**: Explicit state transfer between nodes
2. **Broadcast**: Sending messages to multiple nodes
3. **Request-Response**: Two-way communication between nodes
4. **Pub/Sub**: Publish-subscribe pattern for loose coupling
5. **Channel-based**: Using intermediate channels for communication

---

## JavaScript Examples

### Example 1: Basic Message Passing System

```javascript
class MessagePassingSystem {
  constructor() {
    this.nodes = new Map();
    this.messageQueues = new Map();
    this.handlers = new Map();
    this.messageHistory = [];
    this.messageIdCounter = 0;
  }

  // Register a node in the system
  registerNode(nodeId, messageHandler) {
    this.nodes.set(nodeId, {
      id: nodeId,
      handler: messageHandler,
      status: "active",
    });

    this.messageQueues.set(nodeId, []);
    this.handlers.set(nodeId, messageHandler);

    console.log(`✅ Registered node: ${nodeId}`);
    return this;
  }

  // Send a message from one node to another
  async sendMessage(senderId, receiverId, message, options = {}) {
    const messageId = this.generateMessageId();
    const timestamp = Date.now();

    const messageEnvelope = {
      id: messageId,
      sender: senderId,
      receiver: receiverId,
      payload: message,
      timestamp,
      type: options.type || "standard",
      priority: options.priority || 1,
      metadata: options.metadata || {},
    };

    // Add to receiver's message queue
    if (this.messageQueues.has(receiverId)) {
      const queue = this.messageQueues.get(receiverId);

      // Handle priority insertion
      if (options.priority > 1) {
        const insertIndex = queue.findIndex(
          (msg) => msg.priority < options.priority
        );
        if (insertIndex === -1) {
          queue.push(messageEnvelope);
        } else {
          queue.splice(insertIndex, 0, messageEnvelope);
        }
      } else {
        queue.push(messageEnvelope);
      }

      console.log(`📨 Message ${messageId} sent: ${senderId} → ${receiverId}`);
    } else {
      throw new Error(`Receiver node not found: ${receiverId}`);
    }

    // Record in history
    this.messageHistory.push(messageEnvelope);

    // Auto-process if enabled
    if (options.autoProcess !== false) {
      await this.processMessages(receiverId);
    }

    return messageId;
  }

  // Broadcast message to multiple nodes
  async broadcastMessage(senderId, receiverIds, message, options = {}) {
    const messagePromises = receiverIds.map((receiverId) =>
      this.sendMessage(senderId, receiverId, message, {
        ...options,
        autoProcess: false,
      })
    );

    const messageIds = await Promise.all(messagePromises);

    // Process all messages
    await Promise.all(
      receiverIds.map((receiverId) => this.processMessages(receiverId))
    );

    return messageIds;
  }

  // Process all messages in a node's queue
  async processMessages(nodeId) {
    const node = this.nodes.get(nodeId);
    if (!node) {
      throw new Error(`Node not found: ${nodeId}`);
    }

    const queue = this.messageQueues.get(nodeId);
    const processedMessages = [];

    while (queue.length > 0) {
      const message = queue.shift();

      try {
        console.log(`🔹 Processing message ${message.id} at node ${nodeId}`);

        // Execute node's message handler
        const result = await node.handler(message, node);
        processedMessages.push({
          messageId: message.id,
          success: true,
          result: result,
          processedAt: Date.now(),
        });

        // Handle responses if needed
        if (result && result.responseTo) {
          await this.sendMessage(nodeId, result.responseTo, result.response, {
            type: "response",
          });
        }
      } catch (error) {
        console.error(
          `❌ Error processing message ${message.id} at node ${nodeId}:`,
          error
        );
        processedMessages.push({
          messageId: message.id,
          success: false,
          error: error.message,
          processedAt: Date.now(),
        });

        // Retry logic
        if (message.metadata.retryCount < (message.metadata.maxRetries || 3)) {
          message.metadata.retryCount = (message.metadata.retryCount || 0) + 1;
          queue.push(message); // Put back in queue for retry
          console.log(
            `🔄 Retrying message ${message.id} (attempt ${message.metadata.retryCount})`
          );
        }
      }
    }

    return processedMessages;
  }

  // Request-response pattern
  async sendRequest(senderId, receiverId, request, options = {}) {
    return new Promise(async (resolve, reject) => {
      const requestId = this.generateMessageId();
      const timeout = options.timeout || 30000;

      // Set up response handler
      const responseHandler = (message) => {
        if (message.payload._requestId === requestId) {
          clearTimeout(timeoutId);
          resolve(message.payload);
        }
      };

      // Temporary subscription for response
      const originalHandler = this.handlers.get(senderId);
      this.handlers.set(senderId, async (message) => {
        if (message.payload._requestId === requestId) {
          return responseHandler(message);
        }
        return originalHandler(message);
      });

      // Send request
      await this.sendMessage(
        senderId,
        receiverId,
        {
          ...request,
          _requestId: requestId,
          _responseAddress: senderId,
        },
        { ...options, autoProcess: true }
      );

      // Timeout handling
      const timeoutId = setTimeout(() => {
        this.handlers.set(senderId, originalHandler); // Restore original handler
        reject(new Error(`Request timeout after ${timeout}ms`));
      }, timeout);
    });
  }

  // Pub/Sub pattern implementation
  createTopic(topicName) {
    const subscribers = new Set();

    return {
      subscribe: (nodeId) => {
        subscribers.add(nodeId);
        console.log(`📢 Node ${nodeId} subscribed to topic: ${topicName}`);
      },

      unsubscribe: (nodeId) => {
        subscribers.delete(nodeId);
        console.log(`📢 Node ${nodeId} unsubscribed from topic: ${topicName}`);
      },

      publish: async (senderId, message, options = {}) => {
        const subscriberArray = Array.from(subscribers);
        console.log(
          `📢 Publishing to ${subscriberArray.length} subscribers of topic: ${topicName}`
        );

        return await this.broadcastMessage(
          senderId,
          subscriberArray,
          message,
          options
        );
      },
    };
  }

  // Utility methods
  generateMessageId() {
    return `msg_${this.messageIdCounter++}_${Date.now()}`;
  }

  getNodeQueueSize(nodeId) {
    return this.messageQueues.get(nodeId)?.length || 0;
  }

  getMessageHistory(filter = {}) {
    let history = this.messageHistory;

    if (filter.sender) {
      history = history.filter((msg) => msg.sender === filter.sender);
    }

    if (filter.receiver) {
      history = history.filter((msg) => msg.receiver === filter.receiver);
    }

    if (filter.since) {
      history = history.filter((msg) => msg.timestamp >= filter.since);
    }

    return history.sort((a, b) => a.timestamp - b.timestamp);
  }

  // Visualization of message flow
  visualizeMessageFlow() {
    console.log("\n📊 MESSAGE FLOW VISUALIZATION");

    this.nodes.forEach((node, nodeId) => {
      const sentMessages = this.messageHistory.filter(
        (msg) => msg.sender === nodeId
      );
      const receivedMessages = this.messageHistory.filter(
        (msg) => msg.receiver === nodeId
      );

      console.log(`\n${nodeId}:`);
      console.log(`  📤 Sent: ${sentMessages.length} messages`);
      console.log(`  📥 Received: ${receivedMessages.length} messages`);
      console.log(`  📬 Queue: ${this.getNodeQueueSize(nodeId)} pending`);

      // Show recent communication partners
      const recentReceivers = [
        ...new Set(sentMessages.slice(-5).map((msg) => msg.receiver)),
      ];
      const recentSenders = [
        ...new Set(receivedMessages.slice(-5).map((msg) => msg.sender)),
      ];

      if (recentReceivers.length > 0) {
        console.log(`  ➡️ Recent receivers: ${recentReceivers.join(", ")}`);
      }

      if (recentSenders.length > 0) {
        console.log(`  ⬅️ Recent senders: ${recentSenders.join(", ")}`);
      }
    });

    console.log(`\n📈 Total messages: ${this.messageHistory.length}`);
  }
}
```

### Example 2: Chat Application with Message Passing

```javascript
class ChatApplication {
  constructor() {
    this.messageSystem = new MessagePassingSystem();
    this.rooms = new Map();
    this.users = new Map();
    this.setupChatNodes();
  }

  setupChatNodes() {
    // Define different node types for chat system
    this.messageSystem
      .registerNode("user_manager", this.handleUserManagement.bind(this))
      .registerNode("room_manager", this.handleRoomManagement.bind(this))
      .registerNode(
        "message_processor",
        this.handleMessageProcessing.bind(this)
      )
      .registerNode("moderation_bot", this.handleModeration.bind(this))
      .registerNode("notification_service", this.handleNotifications.bind(this))
      .registerNode("message_logger", this.handleMessageLogging.bind(this));

    // Create some chat rooms
    this.createRoom("general", "General discussion");
    this.createRoom("tech", "Technology talk");
    this.createRoom("random", "Random chatter");
  }

  // User management
  async registerUser(userId, userData) {
    this.users.set(userId, {
      id: userId,
      ...userData,
      joinedAt: new Date().toISOString(),
      status: "online",
    });

    // Notify user manager
    await this.messageSystem.sendMessage("system", "user_manager", {
      type: "user_registered",
      userId,
      userData,
    });

    console.log(`👤 User registered: ${userId}`);
    return userId;
  }

  // Room management
  createRoom(roomId, roomName) {
    this.rooms.set(roomId, {
      id: roomId,
      name: roomName,
      users: new Set(),
      createdAt: new Date().toISOString(),
      messageCount: 0,
    });

    console.log(`💬 Room created: ${roomName} (${roomId})`);
    return roomId;
  }

  async joinRoom(userId, roomId) {
    const room = this.rooms.get(roomId);
    if (!room) throw new Error(`Room not found: ${roomId}`);

    room.users.add(userId);

    // Notify room manager and other services
    await this.messageSystem.broadcastMessage(
      "system",
      ["room_manager", "notification_service", "message_logger"],
      {
        type: "user_joined_room",
        userId,
        roomId,
        timestamp: new Date().toISOString(),
      }
    );

    console.log(`🚪 User ${userId} joined room: ${roomId}`);
  }

  // Send chat message
  async sendMessage(userId, roomId, messageText) {
    const message = {
      type: "chat_message",
      userId,
      roomId,
      text: messageText,
      timestamp: new Date().toISOString(),
      messageId: `msg_${Date.now()}_${Math.random().toString(36).substr(2, 5)}`,
    };

    // Send through processing pipeline
    await this.messageSystem.sendMessage(
      "user_input",
      "message_processor",
      message,
      {
        priority: 2, // Higher priority for user messages
        metadata: { requiresModeration: true },
      }
    );

    return message.messageId;
  }

  // Node handlers
  async handleUserManagement(message, node) {
    console.log(`👥 User Manager processing: ${message.payload.type}`);

    switch (message.payload.type) {
      case "user_registered":
        // Track user registration
        return { processed: true, action: "user_registered" };

      case "user_left":
        // Handle user leaving
        return { processed: true, action: "user_removed" };

      default:
        return { processed: false, error: "Unknown user action" };
    }
  }

  async handleRoomManagement(message, node) {
    console.log(`💬 Room Manager processing: ${message.payload.type}`);

    switch (message.payload.type) {
      case "user_joined_room":
        const { userId, roomId } = message.payload;
        const room = this.rooms.get(roomId);

        if (room) {
          room.users.add(userId);

          // Send welcome message
          await this.messageSystem.sendMessage(
            "room_manager",
            "notification_service",
            {
              type: "room_join_notification",
              userId,
              roomId,
              welcomeMessage: `Welcome to ${room.name}!`,
            }
          );
        }
        return { processed: true, usersInRoom: room.users.size };

      case "user_left_room":
        // Handle user leaving room
        return { processed: true };

      default:
        return { processed: false };
    }
  }

  async handleMessageProcessing(message, node) {
    console.log(
      `🔧 Message Processor: ${message.payload.text?.substring(0, 50)}...`
    );

    const payload = message.payload;

    // Basic message processing
    const processedMessage = {
      ...payload,
      processedAt: new Date().toISOString(),
      containsLinks: this.containsLinks(payload.text),
      sentiment: this.analyzeSentiment(payload.text),
      length: payload.text.length,
    };

    // Route to appropriate services
    const nextSteps = [];

    if (processedMessage.containsLinks) {
      nextSteps.push("moderation_bot");
    }

    nextSteps.push("notification_service", "message_logger");

    await this.messageSystem.broadcastMessage(
      "message_processor",
      nextSteps,
      processedMessage
    );

    return { processed: true, nextSteps, messageId: payload.messageId };
  }

  async handleModeration(message, node) {
    console.log(`🛡️ Moderation Bot checking message...`);

    const payload = message.payload;
    const violations = this.checkForViolations(payload.text);

    if (violations.length > 0) {
      console.log(`🚫 Message violated rules: ${violations.join(", ")}`);

      // Notify system about violation
      await this.messageSystem.sendMessage(
        "moderation_bot",
        "notification_service",
        {
          type: "message_blocked",
          messageId: payload.messageId,
          userId: payload.userId,
          roomId: payload.roomId,
          violations,
          originalMessage: payload.text,
        }
      );

      return { blocked: true, violations };
    }

    return { approved: true };
  }

  async handleNotifications(message, node) {
    console.log(`🔔 Notification Service: ${message.payload.type}`);

    const payload = message.payload;

    switch (payload.type) {
      case "room_join_notification":
        // Send welcome message to user
        console.log(
          `📨 Sending welcome to user ${payload.userId}: ${payload.welcomeMessage}`
        );
        return { notificationSent: true };

      case "chat_message":
        // Notify room participants
        const room = this.rooms.get(payload.roomId);
        if (room) {
          console.log(
            `📢 Broadcasting message to ${room.users.size} users in room ${payload.roomId}`
          );
          // In real implementation, send to each user's client
        }
        return { notified: room.users.size };

      case "message_blocked":
        // Notify user about blocked message
        console.log(
          `🚫 Notifying user ${payload.userId} about blocked message`
        );
        return { userNotified: true };

      default:
        return { processed: false };
    }
  }

  async handleMessageLogging(message, node) {
    console.log(
      `📝 Message Logger: Logging message ${message.payload.messageId}`
    );

    // Simulate logging to database
    const logEntry = {
      ...message.payload,
      loggedAt: new Date().toISOString(),
      logId: `log_${Date.now()}`,
    };

    // In real implementation, save to database
    return { logged: true, logId: logEntry.logId };
  }

  // Utility methods
  containsLinks(text) {
    return /http[s]?:\/\/\S+/.test(text);
  }

  analyzeSentiment(text) {
    // Simple sentiment analysis
    const positiveWords = ["good", "great", "awesome", "excellent", "happy"];
    const negativeWords = ["bad", "terrible", "awful", "hate", "angry"];

    const words = text.toLowerCase().split(/\s+/);
    const positiveCount = words.filter((word) =>
      positiveWords.includes(word)
    ).length;
    const negativeCount = words.filter((word) =>
      negativeWords.includes(word)
    ).length;

    if (positiveCount > negativeCount) return "positive";
    if (negativeCount > positiveCount) return "negative";
    return "neutral";
  }

  checkForViolations(text) {
    const violations = [];
    const bannedWords = ["spam", "scam", "malware", "phishing"];

    bannedWords.forEach((word) => {
      if (text.toLowerCase().includes(word)) {
        violations.push(`contains_banned_word:${word}`);
      }
    });

    if (text.length > 500) {
      violations.push("message_too_long");
    }

    return violations;
  }

  // Demo the chat system
  async demo() {
    console.log("💬 CHAT SYSTEM DEMONSTRATION\n");

    // Register users
    await this.registerUser("alice", {
      name: "Alice",
      email: "alice@example.com",
    });
    await this.registerUser("bob", { name: "Bob", email: "bob@example.com" });
    await this.registerUser("charlie", {
      name: "Charlie",
      email: "charlie@example.com",
    });

    // Join rooms
    await this.joinRoom("alice", "general");
    await this.joinRoom("bob", "general");
    await this.joinRoom("charlie", "tech");

    // Send messages
    console.log("\n--- Sending messages ---");
    await this.sendMessage(
      "alice",
      "general",
      "Hello everyone! How are you doing?"
    );
    await this.sendMessage(
      "bob",
      "general",
      "I'm doing great! This is an awesome chat system."
    );
    await this.sendMessage(
      "charlie",
      "tech",
      "Check out this new technology: http://example.com"
    );
    await this.sendMessage(
      "alice",
      "general",
      "This message contains spam and scam words which should be blocked."
    );

    // Process all messages
    await this.messageSystem.processMessages("message_processor");
    await this.messageSystem.processMessages("moderation_bot");
    await this.messageSystem.processMessages("notification_service");
    await this.messageSystem.processMessages("message_logger");

    // Show message flow
    console.log("\n--- Message Flow Analysis ---");
    this.messageSystem.visualizeMessageFlow();

    // Show message history
    console.log("\n--- Recent Message History ---");
    const recentMessages = this.messageSystem.getMessageHistory({
      since: Date.now() - 60000,
    });
    recentMessages.forEach((msg) => {
      console.log(`${msg.sender} → ${msg.receiver}: ${msg.payload.type}`);
    });
  }
}

// Run the demo
// const chatApp = new ChatApplication();
// chatApp.demo();
```

### Example 3: Workflow Orchestration with Message Passing

```javascript
class WorkflowOrchestrator {
  constructor() {
    this.messageSystem = new MessagePassingSystem();
    this.workflows = new Map();
    this.tasks = new Map();
    this.setupWorkflowNodes();
  }

  setupWorkflowNodes() {
    // Define workflow processing nodes
    this.messageSystem
      .registerNode(
        "workflow_manager",
        this.handleWorkflowManagement.bind(this)
      )
      .registerNode("task_dispatcher", this.handleTaskDispatch.bind(this))
      .registerNode("task_executor", this.handleTaskExecution.bind(this))
      .registerNode(
        "result_aggregator",
        this.handleResultAggregation.bind(this)
      )
      .registerNode("error_handler", this.handleErrors.bind(this))
      .registerNode("progress_tracker", this.handleProgressTracking.bind(this));
  }

  // Define a new workflow
  async defineWorkflow(workflowId, definition) {
    const workflow = {
      id: workflowId,
      definition,
      status: "defined",
      createdAt: new Date().toISOString(),
      currentStep: 0,
      tasks: [],
      results: {},
    };

    this.workflows.set(workflowId, workflow);

    // Notify workflow manager
    await this.messageSystem.sendMessage("system", "workflow_manager", {
      type: "workflow_defined",
      workflowId,
      definition,
    });

    console.log(`📋 Workflow defined: ${workflowId}`);
    return workflowId;
  }

  // Start workflow execution
  async startWorkflow(workflowId, inputData) {
    const workflow = this.workflows.get(workflowId);
    if (!workflow) throw new Error(`Workflow not found: ${workflowId}`);

    workflow.status = "running";
    workflow.startedAt = new Date().toISOString();
    workflow.input = inputData;
    workflow.currentStep = 0;

    // Start with first task
    await this.executeWorkflowStep(workflowId);

    console.log(`🚀 Workflow started: ${workflowId}`);
    return workflowId;
  }

  // Execute a single workflow step
  async executeWorkflowStep(workflowId) {
    const workflow = this.workflows.get(workflowId);
    const step = workflow.definition.steps[workflow.currentStep];

    if (!step) {
      // Workflow completed
      workflow.status = "completed";
      workflow.completedAt = new Date().toISOString();

      await this.messageSystem.sendMessage(
        "workflow_manager",
        "progress_tracker",
        {
          type: "workflow_completed",
          workflowId,
          results: workflow.results,
        }
      );

      return;
    }

    // Create task for this step
    const taskId = this.createTask(workflowId, step, workflow.currentStep);

    // Dispatch task
    await this.messageSystem.sendMessage(
      "workflow_manager",
      "task_dispatcher",
      {
        type: "task_ready",
        workflowId,
        taskId,
        step: workflow.currentStep,
        taskDefinition: step,
      }
    );

    console.log(
      `🔹 Executing step ${workflow.currentStep + 1} of workflow ${workflowId}`
    );
  }

  createTask(workflowId, step, stepIndex) {
    const taskId = `task_${workflowId}_${stepIndex}_${Date.now()}`;

    const task = {
      id: taskId,
      workflowId,
      stepIndex,
      type: step.type,
      parameters: step.parameters,
      status: "pending",
      createdAt: new Date().toISOString(),
      retryCount: 0,
    };

    this.tasks.set(taskId, task);
    return taskId;
  }

  // Node handlers
  async handleWorkflowManagement(message, node) {
    console.log(`📋 Workflow Manager: ${message.payload.type}`);

    const payload = message.payload;

    switch (payload.type) {
      case "workflow_defined":
        // Acknowledge workflow definition
        return { acknowledged: true, workflowId: payload.workflowId };

      case "task_completed":
        // Move to next workflow step
        const workflow = this.workflows.get(payload.workflowId);
        if (workflow) {
          workflow.currentStep++;
          workflow.results[payload.stepIndex] = payload.result;

          // Execute next step
          await this.executeWorkflowStep(payload.workflowId);
        }
        return { nextStepScheduled: true };

      case "task_failed":
        // Handle task failure
        await this.messageSystem.sendMessage(
          "workflow_manager",
          "error_handler",
          {
            type: "workflow_task_failure",
            workflowId: payload.workflowId,
            taskId: payload.taskId,
            error: payload.error,
            stepIndex: payload.stepIndex,
          }
        );
        return { failureHandled: true };

      default:
        return { processed: false };
    }
  }

  async handleTaskDispatch(message, node) {
    console.log(`🚚 Task Dispatcher: ${message.payload.taskId}`);

    const payload = message.payload;
    const task = this.tasks.get(payload.taskId);

    if (!task) {
      throw new Error(`Task not found: ${payload.taskId}`);
    }

    task.status = "dispatched";
    task.dispatchedAt = new Date().toISOString();

    // Route to appropriate executor based on task type
    let executorNode = "task_executor";

    if (task.type === "api_call") {
      executorNode = "api_executor";
    } else if (task.type === "data_processing") {
      executorNode = "data_processor";
    } else if (task.type === "llm_call") {
      executorNode = "llm_executor";
    }

    await this.messageSystem.sendMessage("task_dispatcher", executorNode, {
      type: "execute_task",
      taskId: payload.taskId,
      taskDefinition: task,
    });

    return { dispatched: true, executor: executorNode };
  }

  async handleTaskExecution(message, node) {
    console.log(`⚡ Task Executor: ${message.payload.taskId}`);

    const payload = message.payload;
    const task = this.tasks.get(payload.taskId);

    if (!task) {
      throw new Error(`Task not found: ${payload.taskId}`);
    }

    task.status = "executing";
    task.startedAt = new Date().toISOString();

    try {
      // Simulate task execution
      const result = await this.executeTask(task);

      task.status = "completed";
      task.completedAt = new Date().toISOString();
      task.result = result;

      // Notify workflow manager
      await this.messageSystem.sendMessage(
        "task_executor",
        "workflow_manager",
        {
          type: "task_completed",
          workflowId: task.workflowId,
          taskId: task.id,
          stepIndex: task.stepIndex,
          result: result,
        }
      );

      // Also notify result aggregator
      await this.messageSystem.sendMessage(
        "task_executor",
        "result_aggregator",
        {
          type: "task_result",
          taskId: task.id,
          result: result,
        }
      );

      return { executed: true, result };
    } catch (error) {
      task.status = "failed";
      task.error = error.message;

      await this.messageSystem.sendMessage(
        "task_executor",
        "workflow_manager",
        {
          type: "task_failed",
          workflowId: task.workflowId,
          taskId: task.id,
          stepIndex: task.stepIndex,
          error: error.message,
        }
      );

      throw error;
    }
  }

  async executeTask(task) {
    // Simulate different task types
    await new Promise((resolve) =>
      setTimeout(resolve, 100 + Math.random() * 200)
    );

    switch (task.type) {
      case "data_processing":
        return {
          processed: true,
          records: Math.floor(Math.random() * 1000),
          summary: `Processed data with parameters: ${JSON.stringify(
            task.parameters
          )}`,
        };

      case "api_call":
        return {
          success: true,
          statusCode: 200,
          data: { message: "API call successful" },
          duration: 150,
        };

      case "llm_call":
        return {
          response: `LLM response for: ${task.parameters.prompt}`,
          tokensUsed: 45,
          model: "gpt-3.5-turbo",
        };

      default:
        return { completed: true, taskType: task.type };
    }
  }

  async handleResultAggregation(message, node) {
    console.log(`📊 Result Aggregator: ${message.payload.type}`);

    const payload = message.payload;

    if (payload.type === "task_result") {
      // Aggregate results (in real implementation, store in database)
      console.log(
        `✅ Task ${payload.taskId} completed with result:`,
        payload.result
      );
      return { aggregated: true };
    }

    return { processed: false };
  }

  async handleErrors(message, node) {
    console.log(`🚨 Error Handler: ${message.payload.type}`);

    const payload = message.payload;

    switch (payload.type) {
      case "workflow_task_failure":
        console.log(
          `❌ Task failure in workflow ${payload.workflowId}:`,
          payload.error
        );

        // Retry logic
        const task = this.tasks.get(payload.taskId);
        if (task && task.retryCount < 3) {
          task.retryCount++;
          task.status = "pending";

          console.log(
            `🔄 Retrying task ${task.id} (attempt ${task.retryCount})`
          );

          await this.messageSystem.sendMessage(
            "error_handler",
            "task_dispatcher",
            {
              type: "retry_task",
              taskId: payload.taskId,
            }
          );
        } else {
          console.log(`💀 Task ${payload.taskId} failed after maximum retries`);
          // Escalate or take other action
        }
        return { handled: true };

      default:
        return { processed: false };
    }
  }

  async handleProgressTracking(message, node) {
    console.log(`📈 Progress Tracker: ${message.payload.type}`);

    const payload = message.payload;

    if (payload.type === "workflow_completed") {
      console.log(`🎉 Workflow ${payload.workflowId} completed successfully!`);
      console.log("Final results:", payload.results);
    }

    return { tracked: true };
  }

  // Demo workflow execution
  async demo() {
    console.log("⚙️ WORKFLOW ORCHESTRATION DEMO\n");

    // Define a sample workflow
    const workflowDefinition = {
      name: "Data Processing Pipeline",
      steps: [
        { type: "data_processing", parameters: { operation: "clean" } },
        { type: "api_call", parameters: { endpoint: "/validate" } },
        { type: "llm_call", parameters: { prompt: "Analyze this data" } },
        { type: "data_processing", parameters: { operation: "transform" } },
      ],
    };

    const workflowId = await this.defineWorkflow(
      "demo_workflow",
      workflowDefinition
    );
    await this.startWorkflow(workflowId, { sample: "input data" });

    // Process messages to simulate execution
    await this.processWorkflowMessages();

    // Show workflow status
    this.showWorkflowStatus(workflowId);
  }

  async processWorkflowMessages() {
    // Process messages through the workflow pipeline
    const nodes = [
      "workflow_manager",
      "task_dispatcher",
      "task_executor",
      "result_aggregator",
      "error_handler",
      "progress_tracker",
    ];

    for (const node of nodes) {
      await this.messageSystem.processMessages(node);
    }
  }

  showWorkflowStatus(workflowId) {
    const workflow = this.workflows.get(workflowId);
    console.log("\n📋 WORKFLOW STATUS:");
    console.log(`ID: ${workflow.id}`);
    console.log(`Status: ${workflow.status}`);
    console.log(
      `Current Step: ${workflow.currentStep}/${workflow.definition.steps.length}`
    );
    console.log(`Results:`, workflow.results);

    console.log("\n📊 TASK STATUS:");
    const workflowTasks = Array.from(this.tasks.values()).filter(
      (task) => task.workflowId === workflowId
    );

    workflowTasks.forEach((task) => {
      console.log(`  ${task.id}: ${task.status} (${task.type})`);
    });
  }
}

// Run the demo
// const orchestrator = new WorkflowOrchestrator();
// orchestrator.demo();
```

### Example 4: Advanced Message Patterns with Middleware

```javascript
class AdvancedMessageSystem extends MessagePassingSystem {
  constructor() {
    super();
    this.middleware = new Map();
    this.interceptors = new Map();
    this.messageTransformers = new Map();
    this.setupAdvancedFeatures();
  }

  setupAdvancedFeatures() {
    // Default middleware
    this.addMiddleware("logging", this.loggingMiddleware.bind(this));
    this.addMiddleware("validation", this.validationMiddleware.bind(this));
    this.addMiddleware("retry", this.retryMiddleware.bind(this));

    // Default interceptors
    this.addInterceptor("metrics", this.metricsInterceptor.bind(this));
    this.addInterceptor("caching", this.cachingInterceptor.bind(this));
  }

  // Middleware system
  addMiddleware(name, middlewareFn) {
    this.middleware.set(name, middlewareFn);
    console.log(`🔧 Middleware added: ${name}`);
  }

  async executeMiddleware(message, direction) {
    // direction: 'inbound' or 'outbound'
    let currentMessage = message;

    for (const [name, middleware] of this.middleware) {
      try {
        currentMessage = await middleware(currentMessage, direction);
      } catch (error) {
        console.error(`Middleware ${name} failed:`, error);
        throw error;
      }
    }

    return currentMessage;
  }

  // Interceptor system
  addInterceptor(name, interceptorFn) {
    this.interceptors.set(name, interceptorFn);
    console.log(`🎯 Interceptor added: ${name}`);
  }

  async executeInterceptors(message, phase) {
    // phase: 'pre_send', 'post_send', 'pre_receive', 'post_receive'
    for (const [name, interceptor] of this.interceptors) {
      try {
        await interceptor(message, phase);
      } catch (error) {
        console.error(`Interceptor ${name} failed:`, error);
        // Don't throw - interceptors shouldn't block message flow
      }
    }
  }

  // Override sendMessage to include middleware and interceptors
  async sendMessage(senderId, receiverId, message, options = {}) {
    let processedMessage = message;

    // Pre-send processing
    processedMessage = await this.executeMiddleware(
      { senderId, receiverId, payload: message, options },
      "outbound"
    );

    await this.executeInterceptors(
      { senderId, receiverId, payload: processedMessage, options },
      "pre_send"
    );

    // Original sendMessage logic
    const messageId = await super.sendMessage(
      senderId,
      receiverId,
      processedMessage,
      { ...options, autoProcess: false }
    );

    // Post-send processing
    await this.executeInterceptors(
      { senderId, receiverId, payload: processedMessage, messageId },
      "post_send"
    );

    // Auto-process if enabled
    if (options.autoProcess !== false) {
      await this.processMessages(receiverId);
    }

    return messageId;
  }

  // Middleware implementations
  async loggingMiddleware(message, direction) {
    const logEntry = {
      timestamp: Date.now(),
      direction,
      sender: message.senderId,
      receiver: message.receiverId,
      messageType: message.payload.type,
      middleware: "logging",
    };

    console.log(
      `📝 Logging middleware: ${direction} message from ${message.senderId} to ${message.receiverId}`
    );

    // Add logging metadata to message
    return {
      ...message,
      payload: {
        ...message.payload,
        _logMetadata: {
          ...(message.payload._logMetadata || {}),
          loggedAt: logEntry.timestamp,
        },
      },
    };
  }

  async validationMiddleware(message, direction) {
    if (direction === "inbound") {
      // Validate incoming messages
      if (!message.payload.type) {
        throw new Error("Validation failed: message type required");
      }
      // Add more validation rules as needed
    }
    return message;
  }
}
```
