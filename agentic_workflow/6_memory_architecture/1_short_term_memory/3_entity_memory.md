# Entity Memory Architecture

Entity Memory focuses on tracking and remembering specific entities (people, places, things, concepts) mentioned in conversations. It extracts and maintains knowledge about these entities over time.

## How Entity Memory Works

- **Extracts entities** from conversation text
- **Stores entity attributes** and relationships
- **Updates entity knowledge** as new information emerges
- **Provides entity context** for future interactions

## JavaScript Implementation Examples

### 1. Basic Entity Memory Class

```javascript
class EntityMemory {
  constructor() {
    this.entities = new Map(); // entityName -> Entity object
    this.entityRelationships = new Map(); // entityName -> Set of related entities
  }

  addMessage(role, content) {
    const entities = this.extractEntities(content);
    this.updateEntityKnowledge(entities, role, content);
  }

  extractEntities(text) {
    // Simple entity extraction using patterns and keywords
    const entities = new Set();

    // Patterns for different entity types
    const patterns = {
      person: /\b([A-Z][a-z]+ [A-Z][a-z]+)\b/g, // Proper names
      place: /\b(in|at|from) ([A-Z][a-zA-Z]+)\b/g, // Locations
      date: /\b(\d{1,2}\/\d{1,2}\/\d{4}|\d{4})\b/g, // Dates
      topic: /\b(JavaScript|Python|AI|ML|web development)\b/gi, // Specific topics
    };

    // Extract using patterns
    for (const [type, pattern] of Object.entries(patterns)) {
      const matches = text.match(pattern);
      if (matches) {
        matches.forEach((match) => {
          entities.add({ name: match, type, context: text });
        });
      }
    }

    // Additional keyword-based extraction
    this.extractEntitiesByKeywords(text, entities);

    return Array.from(entities);
  }

  extractEntitiesByKeywords(text, entities) {
    const keywordMap = {
      person: ["my friend", "my brother", "sister", "colleague", "teacher"],
      place: ["restaurant", "cafe", "hotel", "airport", "school"],
      organization: ["company", "university", "school", "hospital"],
    };

    for (const [type, keywords] of Object.entries(keywordMap)) {
      keywords.forEach((keyword) => {
        if (text.toLowerCase().includes(keyword)) {
          // Extract the entity name following the keyword
          const regex = new RegExp(`${keyword} ([A-Za-z]+)`, "i");
          const match = text.match(regex);
          if (match) {
            entities.add({ name: match[1], type, context: text });
          }
        }
      });
    }
  }

  updateEntityKnowledge(entities, role, context) {
    entities.forEach((entity) => {
      if (!this.entities.has(entity.name)) {
        // Create new entity
        this.entities.set(entity.name, {
          name: entity.name,
          type: entity.type,
          firstMentioned: new Date(),
          lastMentioned: new Date(),
          mentions: 1,
          contexts: [context],
          attributes: new Map(),
          source: role,
        });
      } else {
        // Update existing entity
        const existingEntity = this.entities.get(entity.name);
        existingEntity.lastMentioned = new Date();
        existingEntity.mentions++;
        existingEntity.contexts.push(context);

        // Keep only recent contexts
        if (existingEntity.contexts.length > 5) {
          existingEntity.contexts.shift();
        }
      }

      this.extractAttributes(entity.name, context, role);
    });
  }

  extractAttributes(entityName, context, role) {
    const entity = this.entities.get(entityName);
    if (!entity) return;

    // Extract attributes based on context patterns
    const attributePatterns = {
      age: /(\d+) years? old/,
      location: /(live|located) in ([A-Z][a-zA-Z]+)/,
      profession: /(is a|works as) ([a-zA-Z]+)/,
      preference: /(like|love|enjoy) ([a-zA-Z]+)/,
    };

    for (const [attrType, pattern] of Object.entries(attributePatterns)) {
      const match = context.match(pattern);
      if (match) {
        entity.attributes.set(attrType, match[2] || match[1]);
      }
    }
  }

  getEntity(name) {
    return this.entities.get(name);
  }

  getAllEntities() {
    return Array.from(this.entities.values());
  }

  getEntitiesByType(type) {
    return this.getAllEntities().filter((entity) => entity.type === type);
  }

  findRelatedEntities(entityName) {
    const related = new Set();
    this.entityRelationships.get(entityName)?.forEach((relatedName) => {
      related.add(this.entities.get(relatedName));
    });
    return Array.from(related);
  }

  addRelationship(entity1, entity2, relationship) {
    if (!this.entityRelationships.has(entity1)) {
      this.entityRelationships.set(entity1, new Set());
    }
    if (!this.entityRelationships.has(entity2)) {
      this.entityRelationships.set(entity2, new Set());
    }

    this.entityRelationships.get(entity1).add(entity2);
    this.entityRelationships.get(entity2).add(entity1);

    // Store relationship type
    const e1 = this.entities.get(entity1);
    const e2 = this.entities.get(entity2);
    if (e1 && e2) {
      e1.attributes.set(`related_to_${entity2}`, relationship);
      e2.attributes.set(`related_to_${entity1}`, relationship);
    }
  }
}

// Usage example
const entityMemory = new EntityMemory();

entityMemory.addMessage(
  "user",
  "My friend Alice is a doctor who lives in Boston"
);
entityMemory.addMessage(
  "assistant",
  "That's great! What does Alice specialize in?"
);
entityMemory.addMessage(
  "user",
  "Alice works at Massachusetts General Hospital and she loves hiking"
);

console.log("Entity Memory Contents:");
console.log(entityMemory.getEntity("Alice"));
console.log("\nAll entities:", entityMemory.getAllEntities());
```

### 2. Advanced Entity Memory with Relationship Tracking

```javascript
class RelationshipAwareEntityMemory {
  constructor() {
    this.entities = new Map();
    this.relationships = new Map(); // entity -> { target: entity, type: string }
    this.conversationHistory = [];
  }

  addMessage(role, content) {
    this.conversationHistory.push({ role, content, timestamp: Date.now() });

    const entities = this.advancedEntityExtraction(content);
    this.updateEntitiesWithRelationships(entities, content);
  }

  advancedEntityExtraction(text) {
    const entities = new Set();

    // Use more sophisticated extraction patterns
    const patterns = [
      // People with titles
      {
        pattern: /\b(Dr\.|Mr\.|Ms\.|Mrs\.) ([A-Z][a-z]+ [A-Z][a-z]+)\b/g,
        type: "person",
      },
      // Organizations
      {
        pattern: /\b([A-Z][a-z]+ (Inc|Corp|Company|University))\b/g,
        type: "organization",
      },
      // Email addresses
      { pattern: /\b\S+@\S+\.\S+\b/g, type: "contact" },
      // Phone numbers
      { pattern: /\b\d{3}-\d{3}-\d{4}\b/g, type: "contact" },
    ];

    patterns.forEach(({ pattern, type }) => {
      const matches = text.match(pattern);
      if (matches) {
        matches.forEach((match) => {
          entities.add({ name: match, type, source: text });
        });
      }
    });

    return Array.from(entities);
  }

  updateEntitiesWithRelationships(entities, context) {
    entities.forEach((entity) => {
      if (!this.entities.has(entity.name)) {
        this.entities.set(entity.name, {
          ...entity,
          firstSeen: new Date(),
          lastSeen: new Date(),
          mentionCount: 1,
          contexts: [context],
          properties: new Map(),
        });
      } else {
        const existing = this.entities.get(entity.name);
        existing.lastSeen = new Date();
        existing.mentionCount++;
        existing.contexts.push(context);
      }

      this.extractPropertiesAndRelationships(entity.name, context);
    });
  }

  extractPropertiesAndRelationships(entityName, context) {
    const entity = this.entities.get(entityName);
    if (!entity) return;

    // Extract properties
    const propertyPatterns = {
      age: /aged? (\d+)/,
      city: /in ([A-Z][a-zA-Z]+)/,
      role: /(developer|manager|director|engineer)/i,
      department: /(engineering|sales|marketing|support)/i,
    };

    Object.entries(propertyPatterns).forEach(([prop, pattern]) => {
      const match = context.match(pattern);
      if (match) {
        entity.properties.set(prop, match[1] || match[0]);
      }
    });

    // Extract relationships with other mentioned entities
    this.extractRelationships(entityName, context);
  }

  extractRelationships(entityName, context) {
    const otherEntities = Array.from(this.entities.keys()).filter(
      (name) => name !== entityName
    );

    otherEntities.forEach((otherEntity) => {
      if (context.includes(otherEntity)) {
        this.addRelationship(entityName, otherEntity, "mentioned_together");
      }
    });

    // Specific relationship patterns
    const relationshipPatterns = [
      {
        pattern: new RegExp(`${entityName} works with ([A-Za-z ]+)`, "i"),
        type: "colleague",
      },
      {
        pattern: new RegExp(`${entityName} reports to ([A-Za-z ]+)`, "i"),
        type: "reports_to",
      },
      {
        pattern: new RegExp(`${entityName} manages ([A-Za-z ]+)`, "i"),
        type: "manages",
      },
    ];

    relationshipPatterns.forEach(({ pattern, type }) => {
      const match = context.match(pattern);
      if (match && match[1]) {
        this.addRelationship(entityName, match[1].trim(), type);
      }
    });
  }

  addRelationship(entity1, entity2, relationshipType) {
    const key = `${entity1}-${entity2}`;
    if (!this.relationships.has(key)) {
      this.relationships.set(key, {
        entity1,
        entity2,
        type: relationshipType,
        strength: 1,
        lastUpdated: new Date(),
      });
    } else {
      const rel = this.relationships.get(key);
      rel.strength++;
      rel.lastUpdated = new Date();
    }
  }

  getEntityNetwork(entityName) {
    const network = {
      entity: this.entities.get(entityName),
      relationships: [],
    };

    this.relationships.forEach((rel, key) => {
      if (rel.entity1 === entityName || rel.entity2 === entityName) {
        const otherEntity =
          rel.entity1 === entityName ? rel.entity2 : rel.entity1;
        network.relationships.push({
          entity: otherEntity,
          type: rel.type,
          strength: rel.strength,
          target: this.entities.get(otherEntity),
        });
      }
    });

    return network;
  }

  searchEntities(query) {
    const results = [];
    const queryLower = query.toLowerCase();

    this.entities.forEach((entity, name) => {
      if (name.toLowerCase().includes(queryLower)) {
        results.push(entity);
      } else {
        // Search in properties and contexts
        for (const [key, value] of entity.properties) {
          if (value.toString().toLowerCase().includes(queryLower)) {
            results.push(entity);
            break;
          }
        }
      }
    });

    return results;
  }

  getMemorySummary() {
    return {
      totalEntities: this.entities.size,
      totalRelationships: this.relationships.size,
      entityTypes: this.getEntityTypeCounts(),
      mostMentioned: this.getMostMentionedEntities(5),
    };
  }

  getEntityTypeCounts() {
    const counts = {};
    this.entities.forEach((entity) => {
      counts[entity.type] = (counts[entity.type] || 0) + 1;
    });
    return counts;
  }

  getMostMentionedEntities(limit = 10) {
    return Array.from(this.entities.values())
      .sort((a, b) => b.mentionCount - a.mentionCount)
      .slice(0, limit);
  }
}

// Usage example
const advancedMemory = new RelationshipAwareEntityMemory();

advancedMemory.addMessage("user", "John Smith is our project manager");
advancedMemory.addMessage("assistant", "What team does John Smith manage?");
advancedMemory.addMessage(
  "user",
  "John manages the engineering team including Alice Brown and Bob Wilson"
);
advancedMemory.addMessage(
  "user",
  "Alice Brown is a senior developer aged 32 in Boston"
);

console.log("John Smith Network:");
console.log(advancedMemory.getEntityNetwork("John Smith"));
console.log("\nMemory Summary:", advancedMemory.getMemorySummary());
```

### 3. Entity Memory with Temporal Context

```javascript
class TemporalEntityMemory {
  constructor() {
    this.entities = new Map();
    this.timeline = []; // Chronological entity events
    this.entityStates = new Map(); // Current state of entities
  }

  addMessage(role, content, timestamp = Date.now()) {
    const entities = this.extractEntitiesWithTemporalContext(
      content,
      timestamp
    );

    entities.forEach((entity) => {
      this.recordEntityEvent(entity, content, role, timestamp);
      this.updateEntityState(entity, timestamp);
    });
  }

  extractEntitiesWithTemporalContext(text, timestamp) {
    const entities = new Set();

    // Extract entities with temporal indicators
    const temporalPatterns = [
      { pattern: /\b(currently|now) ([A-Z][a-z]+) is/g, type: "current_state" },
      {
        pattern: /\b(formerly|previously) ([A-Z][a-z]+) was/g,
        type: "past_state",
      },
      { pattern: /\b(will|soon) ([A-Z][a-z]+) will/g, type: "future_state" },
    ];

    temporalPatterns.forEach(({ pattern, type }) => {
      const matches = text.match(pattern);
      if (matches) {
        matches.forEach((match) => {
          // Extract entity name and temporal context
          const entityMatch = match.match(/([A-Z][a-z]+)/);
          if (entityMatch) {
            entities.add({
              name: entityMatch[1],
              type: "person", // Default type
              temporalContext: type,
              timestamp,
              source: text,
            });
          }
        });
      }
    });

    return Array.from(entities);
  }

  recordEntityEvent(entity, context, role, timestamp) {
    const event = {
      entity: entity.name,
      type: "mention",
      context,
      role,
      timestamp,
      temporalContext: entity.temporalContext,
    };

    this.timeline.push(event);

    // Keep timeline manageable
    if (this.timeline.length > 1000) {
      this.timeline = this.timeline.slice(-500); // Keep last 500 events
    }

    // Update entity record
    if (!this.entities.has(entity.name)) {
      this.entities.set(entity.name, {
        name: entity.name,
        firstSeen: timestamp,
        lastSeen: timestamp,
        events: [event],
        currentState: {},
        historicalStates: [],
      });
    } else {
      const existing = this.entities.get(entity.name);
      existing.lastSeen = timestamp;
      existing.events.push(event);
    }
  }

  updateEntityState(entity, timestamp) {
    const entityRecord = this.entities.get(entity.name);
    if (!entityRecord) return;

    // Extract state information from the context
    const stateInfo = this.extractStateInformation(
      entity.temporalContext,
      entity.source
    );

    if (stateInfo) {
      if (entity.temporalContext === "current_state") {
        entityRecord.currentState = {
          ...entityRecord.currentState,
          ...stateInfo,
          lastUpdated: timestamp,
        };
      } else {
        entityRecord.historicalStates.push({
          state: stateInfo,
          period: entity.temporalContext,
          timestamp,
        });
      }
    }
  }

  extractStateInformation(temporalContext, text) {
    const state = {};

    // Extract role/position
    const roleMatch = text.match(/(is|was|will be) (a|an|the) ([a-zA-Z]+)/);
    if (roleMatch) {
      state.role = roleMatch[3];
    }

    // Extract location
    const locationMatch = text.match(/(in|at) ([A-Z][a-zA-Z]+)/);
    if (locationMatch) {
      state.location = locationMatch[2];
    }

    // Extract status
    const statusMatch = text.match(/(working|managing|developing) ([a-zA-Z]+)/);
    if (statusMatch) {
      state.status = statusMatch[1];
      state.activity = statusMatch[2];
    }

    return Object.keys(state).length > 0 ? state : null;
  }

  getEntityTimeline(entityName, startTime = 0, endTime = Date.now()) {
    const entity = this.entities.get(entityName);
    if (!entity) return [];

    return entity.events
      .filter(
        (event) => event.timestamp >= startTime && event.timestamp <= endTime
      )
      .sort((a, b) => a.timestamp - b.timestamp);
  }

  getEntityCurrentState(entityName) {
    const entity = this.entities.get(entityName);
    return entity ? entity.currentState : null;
  }

  getEntityHistory(entityName) {
    const entity = this.entities.get(entityName);
    if (!entity) return null;

    return {
      entity: entityName,
      currentState: entity.currentState,
      historicalStates: entity.historicalStates,
      timeline: entity.events.slice(-10), // Last 10 events
    };
  }

  findEntitiesByState(stateQuery) {
    const results = [];

    this.entities.forEach((entity, name) => {
      // Check current state
      const currentStateStr = JSON.stringify(entity.currentState).toLowerCase();
      if (currentStateStr.includes(stateQuery.toLowerCase())) {
        results.push({
          entity: name,
          state: "current",
          details: entity.currentState,
        });
      }

      // Check historical states
      entity.historicalStates.forEach((historicalState) => {
        const historicalStr = JSON.stringify(
          historicalState.state
        ).toLowerCase();
        if (historicalStr.includes(stateQuery.toLowerCase())) {
          results.push({
            entity: name,
            state: "historical",
            period: historicalState.period,
            details: historicalState.state,
            timestamp: historicalState.timestamp,
          });
        }
      });
    });

    return results;
  }
}

// Usage example
const temporalMemory = new TemporalEntityMemory();

const now = Date.now();
temporalMemory.addMessage("user", "Currently John is the project manager", now);
temporalMemory.addMessage(
  "user",
  "Previously John was a developer",
  now - 86400000
); // 1 day ago
temporalMemory.addMessage(
  "user",
  "Soon John will be a director",
  now + 86400000
); // 1 day from now

console.log("John's History:");
console.log(temporalMemory.getEntityHistory("John"));
console.log("\nJohn's Current State:");
console.log(temporalMemory.getEntityCurrentState("John"));
```

### 4. Practical Application: Customer Service Entity Memory

```javascript
class CustomerServiceEntityMemory {
  constructor() {
    this.customers = new Map();
    this.issues = new Map();
    this.products = new Map();
    this.interactions = [];
  }

  recordInteraction(
    customerInfo,
    issueDescription,
    productMentioned,
    timestamp = Date.now()
  ) {
    // Extract entities from interaction
    const customerEntity = this.extractCustomerEntity(customerInfo);
    const issueEntity = this.extractIssueEntity(issueDescription);
    const productEntity = this.extractProductEntity(productMentioned);

    // Record the interaction
    const interaction = {
      id: this.interactions.length + 1,
      customer: customerEntity.name,
      issue: issueEntity.name,
      product: productEntity.name,
      timestamp,
      description: issueDescription,
    };

    this.interactions.push(interaction);

    // Update entity knowledge
    this.updateCustomerHistory(customerEntity, interaction);
    this.updateIssueTracking(issueEntity, interaction);
    this.updateProductKnowledge(productEntity, interaction);

    return interaction;
  }

  extractCustomerEntity(customerInfo) {
    const patterns = [
      { pattern: /name: ([A-Za-z ]+)/, field: "name" },
      { pattern: /email: (\S+@\S+)/, field: "email" },
      { pattern: /phone: ([\d-]+)/, field: "phone" },
    ];

    const customer = { name: "Unknown", type: "customer" };

    patterns.forEach(({ pattern, field }) => {
      const match = customerInfo.match(pattern);
      if (match) {
        customer[field] = match[1];
        if (field === "name") customer.name = match[1];
      }
    });

    return customer;
  }

  extractIssueEntity(issueDescription) {
    // Categorize issues
    const categories = {
      technical: ["not working", "error", "bug", "crash"],
      billing: ["charge", "invoice", "payment", "bill"],
      feature: ["how to", "can I", "possible to"],
    };

    let issueType = "general";
    for (const [type, keywords] of Object.entries(categories)) {
      if (
        keywords.some((keyword) =>
          issueDescription.toLowerCase().includes(keyword)
        )
      ) {
        issueType = type;
        break;
      }
    }

    return {
      name: issueDescription.substring(0, 50),
      type: "issue",
      category: issueType,
      severity: this.assessSeverity(issueDescription),
    };
  }

  extractProductEntity(productMentioned) {
    const knownProducts = [
      "WebApp",
      "MobileApp",
      "API",
      "Dashboard",
      "Database",
    ];
    const mentionedProduct =
      knownProducts.find((product) =>
        productMentioned.toLowerCase().includes(product.toLowerCase())
      ) || "Unknown";

    return {
      name: mentionedProduct,
      type: "product",
    };
  }

  assessSeverity(description) {
    if (
      description.toLowerCase().includes("urgent") ||
      description.includes("!!!")
    ) {
      return "high";
    }
    if (description.toLowerCase().includes("important")) {
      return "medium";
    }
    return "low";
  }

  updateCustomerHistory(customer, interaction) {
    if (!this.customers.has(customer.name)) {
      this.customers.set(customer.name, {
        ...customer,
        firstContact: interaction.timestamp,
        lastContact: interaction.timestamp,
        totalInteractions: 1,
        issues: [interaction.issue],
        products: [interaction.product],
      });
    } else {
      const existing = this.customers.get(customer.name);
      existing.lastContact = interaction.timestamp;
      existing.totalInteractions++;
      if (!existing.issues.includes(interaction.issue)) {
        existing.issues.push(interaction.issue);
      }
      if (!existing.products.includes(interaction.product)) {
        existing.products.push(interaction.product);
      }
    }
  }

  updateIssueTracking(issue, interaction) {
    const issueKey = `${issue.category}-${issue.name}`;

    if (!this.issues.has(issueKey)) {
      this.issues.set(issueKey, {
        ...issue,
        firstReported: interaction.timestamp,
        lastReported: interaction.timestamp,
        reportCount: 1,
        affectedCustomers: [interaction.customer],
        affectedProducts: [interaction.product],
      });
    } else {
      const existing = this.issues.get(issueKey);
      existing.lastReported = interaction.timestamp;
      existing.reportCount++;
      if (!existing.affectedCustomers.includes(interaction.customer)) {
        existing.affectedCustomers.push(interaction.customer);
      }
    }
  }

  updateProductKnowledge(product, interaction) {
    if (!this.products.has(product.name)) {
      this.products.set(product.name, {
        ...product,
        issues: [interaction.issue],
        customers: [interaction.customer],
      });
    } else {
      const existing = this.products.get(product.name);
      if (!existing.issues.includes(interaction.issue)) {
        existing.issues.push(interaction.issue);
      }
      if (!existing.customers.includes(interaction.customer)) {
        existing.customers.push(interaction.customer);
      }
    }
  }

  getCustomerProfile(customerName) {
    const customer = this.customers.get(customerName);
    if (!customer) return null;

    const customerInteractions = this.interactions.filter(
      (i) => i.customer === customerName
    );
    const frequentIssues = this.getFrequentItems(
      customerInteractions.map((i) => i.issue)
    );
    const usedProducts = this.getFrequentItems(
      customerInteractions.map((i) => i.product)
    );

    return {
      profile: customer,
      interactions: customerInteractions.slice(-5), // Last 5 interactions
      frequentIssues,
      usedProducts,
      satisfactionScore: this.calculateSatisfactionScore(customerInteractions),
    };
  }

  getFrequentItems(items, limit = 3) {
    const frequency = {};
    items.forEach((item) => {
      frequency[item] = (frequency[item] || 0) + 1;
    });

    return Object.entries(frequency)
      .sort(([, a], [, b]) => b - a)
      .slice(0, limit)
      .map(([item]) => item);
  }

  calculateSatisfactionScore(interactions) {
    // Simple heuristic based on interaction patterns
    const total = interactions.length;
    if (total === 0) return 0;

    const repeatedIssues = interactions.filter((interaction, index) => {
      return interactions
        .slice(0, index)
        .some((prev) => prev.issue === interaction.issue);
    }).length;

    return Math.max(0, 10 - (repeatedIssues / total) * 10);
  }

  getProductAnalysis(productName) {
    const product = this.products.get(productName);
    if (!product) return null;

    const productInteractions = this.interactions.filter(
      (i) => i.product === productName
    );
    const commonIssues = this.getFrequentItems(
      productInteractions.map((i) => i.issue)
    );

    return {
      product: productName,
      totalReports: productInteractions.length,
      commonIssues,
      affectedCustomers: product.customers.length,
      issueTrend: this.calculateIssueTrend(productInteractions),
    };
  }

  calculateIssueTrend(interactions) {
    if (interactions.length < 2) return "stable";

    const recent = interactions.slice(-5);
    const older = interactions.slice(-10, -5);

    return recent.length > older.length ? "increasing" : "decreasing";
  }
}

// Usage example
const customerMemory = new CustomerServiceEntityMemory();

// Record some interactions
customerMemory.recordInteraction(
  "name: John Doe, email: john@email.com",
  "The dashboard is not loading properly, urgent!",
  "WebApp Dashboard"
);

customerMemory.recordInteraction(
  "name: Jane Smith, email: jane@email.com",
  "How do I export data from the MobileApp?",
  "MobileApp"
);

console.log("Customer Profiles:");
console.log(customerMemory.getCustomerProfile("John Doe"));
console.log("\nProduct Analysis:");
console.log(customerMemory.getProductAnalysis("Dashboard"));
```

## Key Benefits of Entity Memory

1. **Targeted Recall**: Remembers specific people, places, and things
2. **Relationship Mapping**: Tracks connections between entities
3. **Attribute Tracking**: Maintains detailed information about entities
4. **Temporal Awareness**: Understands how entities change over time
5. **Efficient Querying**: Quickly retrieves relevant entity information

## Common Use Cases

- **Customer relationship management** (CRM) systems
- **Personal assistants** that remember user preferences
- **Knowledge management** systems
- **Investigative tools** that track people and relationships
- **Educational systems** that remember student progress

Entity Memory is particularly powerful for applications that need to maintain detailed knowledge about specific subjects across multiple conversations or sessions.
