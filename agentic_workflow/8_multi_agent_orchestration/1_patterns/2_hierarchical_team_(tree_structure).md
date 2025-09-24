# Hierarchical Teams Pattern (Tree Structure)

The **Hierarchical Teams** pattern organizes agents in a tree-like structure where higher-level managers coordinate subordinate agents, creating multiple layers of specialization and responsibility.

## Key Concepts

- **Tree Structure**: Parent-child relationships between agents
- **Delegation Chain**: Tasks flow down the hierarchy
- **Specialization Layers**: Different levels handle different abstraction levels
- **Reporting Chain**: Results bubble up through the hierarchy
- **Span of Control**: Each manager handles a limited number of subordinates

## JavaScript Implementation Example

### Basic Hierarchical Structure

```javascript
// Base Agent Class
class Agent {
  constructor(name, role) {
    this.name = name;
    this.role = role;
    this.subordinates = [];
    this.supervisor = null;
  }

  addSubordinate(agent) {
    agent.supervisor = this;
    this.subordinates.push(agent);
    console.log(`${agent.name} (${agent.role}) now reports to ${this.name}`);
  }

  async delegate(task, depth = 0) {
    const indent = "  ".repeat(depth);
    console.log(`${indent}${this.name} [${this.role}] considering: ${task}`);
    
    if (this.canHandle(task)) {
      console.log(`${indent}${this.name} handling task directly`);
      return await this.handle(task, depth);
    }
    
    if (this.subordinates.length > 0) {
      console.log(`${indent}${this.name} delegating to subordinates...`);
      for (const subordinate of this.subordinates) {
        try {
          const result = await subordinate.delegate(task, depth + 1);
          if (result.handled) {
            console.log(`${indent}${this.name} received result from ${subordinate.name}`);
            return result;
          }
        } catch (error) {
          console.log(`${indent}${subordinate.name} failed, trying next...`);
        }
      }
    }
    
    throw new Error(`No one in hierarchy could handle: ${task}`);
  }

  canHandle(task) {
    // Base implementation - override in subclasses
    return false;
  }

  async handle(task, depth) {
    // Base implementation
    await this.delay(500);
    return { handled: true, by: this.name, result: `Task completed by ${this.name}` };
  }

  delay(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }
}

// Specialized Agent Classes
class CEO extends Agent {
  constructor() {
    super("Alice Chen", "CEO");
  }

  canHandle(task) {
    return task.includes("strategic") || task.includes("merger") || task.includes("investment");
  }

  async handle(task, depth) {
    console.log(`  CEO making high-level decision: ${task}`);
    await this.delay(1000);
    return { handled: true, by: this.name, result: "Strategic decision made" };
  }
}

class EngineeringManager extends Agent {
  constructor() {
    super("Bob Smith", "Engineering Manager");
  }

  canHandle(task) {
    return task.includes("technical") || task.includes("architecture") || task.includes("team");
  }

  async handle(task, depth) {
    console.log(`  Engineering Manager coordinating: ${task}`);
    await this.delay(800);
    return { handled: true, by: this.name, result: "Engineering solution designed" };
  }
}

class MarketingManager extends Agent {
  constructor() {
    super("Carol Davis", "Marketing Manager");
  }

  canHandle(task) {
    return task.includes("marketing") || task.includes("campaign") || task.includes("brand");
  }

  async handle(task, depth) {
    console.log(`  Marketing Manager planning: ${task}`);
    await this.delay(700);
    return { handled: true, by: this.name, result: "Marketing plan created" };
  }
}

class SeniorDeveloper extends Agent {
  constructor() {
    super("David Wilson", "Senior Developer");
  }

  canHandle(task) {
    return task.includes("complex") || task.includes("algorithm") || task.includes("review");
  }

  async handle(task, depth) {
    console.log(`    Senior Developer working on: ${task}`);
    await this.delay(600);
    return { handled: true, by: this.name, result: "Complex code implemented" };
  }
}

class JuniorDeveloper extends Agent {
  constructor() {
    super("Eva Garcia", "Junior Developer");
  }

  canHandle(task) {
    return task.includes("simple") || task.includes("bug") || task.includes("feature");
  }

  async handle(task, depth) {
    console.log(`      Junior Developer implementing: ${task}`);
    await this.delay(400);
    return { handled: true, by: this.name, result: "Task implemented" };
  }
}

// Building the Hierarchy
function buildOrganization() {
  // Create agents
  const ceo = new CEO();
  const engineeringManager = new EngineeringManager();
  const marketingManager = new MarketingManager();
  const seniorDev = new SeniorDeveloper();
  const juniorDev1 = new JuniorDeveloper();
  const juniorDev2 = new JuniorDeveloper();

  // Build hierarchy
  ceo.addSubordinate(engineeringManager);
  ceo.addSubordinate(marketingManager);
  engineeringManager.addSubordinate(seniorDev);
  seniorDev.addSubordinate(juniorDev1);
  seniorDev.addSubordinate(juniorDev2);

  return ceo;
}

// Usage Example
async function demonstrateHierarchy() {
  console.log("🏢 Building Organizational Hierarchy...\n");
  const ceo = buildOrganization();

  const tasks = [
    "strategic company direction",
    "complex algorithm implementation",
    "simple bug fix in login system",
    "marketing campaign planning",
    "technical architecture review"
  ];

  for (const task of tasks) {
    console.log(`\n📋 Processing: "${task}"`);
    console.log("=" .repeat(50));
    
    try {
      const result = await ceo.delegate(task);
      console.log(`✅ Final Result: ${result.result}\n`);
    } catch (error) {
      console.log(`❌ Error: ${error.message}\n`);
    }
  }
}

demonstrateHierarchy();
```

## Advanced Example: Multi-level Problem Solving

```javascript
class HierarchicalProblemSolver {
  constructor() {
    this.root = this.buildProblemSolvingHierarchy();
  }

  buildProblemSolvingHierarchy() {
    // Level 1: High-level strategists
    const chiefScientist = new StrategicSolver("Dr. Sophia", "Chief Scientist");
    const leadArchitect = new StrategicSolver("Mark", "Lead Architect");

    // Level 2: Domain experts
    const aiExpert = new DomainSolver("Dr. Chen", "AI Specialist");
    const dataExpert = new DomainSolver("Sarah", "Data Scientist");
    const cloudExpert = new DomainSolver("Mike", "Cloud Architect");

    // Level 3: Technical specialists
    const mlEngineer = new TechnicalSolver("Alex", "ML Engineer");
    const backendEngineer = new TechnicalSolver("Lisa", "Backend Engineer");
    const devopsEngineer = new TechnicalSolver("Tom", "DevOps Engineer");

    // Level 4: Implementers
    const juniorEngineer1 = new ImplementationSolver("John", "Junior Engineer");
    const juniorEngineer2 = new ImplementationSolver("Emma", "Junior Engineer");

    // Build the tree
    chiefScientist.addSubordinate(aiExpert);
    chiefScientist.addSubordinate(dataExpert);
    leadArchitect.addSubordinate(cloudExpert);
    
    aiExpert.addSubordinate(mlEngineer);
    cloudExpert.addSubordinate(backendEngineer);
    cloudExpert.addSubordinate(devopsEngineer);
    
    mlEngineer.addSubordinate(juniorEngineer1);
    backendEngineer.addSubordinate(juniorEngineer2);

    return chiefScientist;
  }

  async solveProblem(problem) {
    console.log(`\n🔍 Solving: ${problem}`);
    console.log("=" .repeat(60));
    
    return await this.root.delegate(problem);
  }
}

// Specialized Solver Classes
class StrategicSolver extends Agent {
  canHandle(problem) {
    return problem.includes("strategy") || 
           problem.includes("research") || 
           problem.includes("architecture");
  }

  async handle(problem, depth) {
    await this.delay(1000);
    return { 
      handled: true, 
      by: this.name, 
      result: `Strategic solution for: ${problem}`,
      level: "strategic"
    };
  }
}

class DomainSolver extends Agent {
  canHandle(problem) {
    return problem.includes("AI") || 
           problem.includes("data") || 
           problem.includes("analysis") ||
           problem.includes("model");
  }

  async handle(problem, depth) {
    await this.delay(800);
    return { 
      handled: true, 
      by: this.name, 
      result: `Domain-specific solution for: ${problem}`,
      level: "domain"
    };
  }
}

class TechnicalSolver extends Agent {
  canHandle(problem) {
    return problem.includes("implementation") || 
           problem.includes("technical") || 
           problem.includes("code") ||
           problem.includes("system");
  }

  async handle(problem, depth) {
    await this.delay(600);
    return { 
      handled: true, 
      by: this.name, 
      result: `Technical implementation for: ${problem}`,
      level: "technical"
    };
  }
}

class ImplementationSolver extends Agent {
  canHandle(problem) {
    return problem.includes("simple") || 
           problem.includes("routine") || 
           problem.includes("basic") ||
           problem.includes("task");
  }

  async handle(problem, depth) {
    await this.delay(400);
    return { 
      handled: true, 
      by: this.name, 
      result: `Implemented: ${problem}`,
      level: "implementation"
    };
  }
}

// Advanced Usage with Monitoring
async function advancedDemo() {
  const solver = new HierarchicalProblemSolver();
  
  const problems = [
    "research new AI strategy for company",
    "implement complex machine learning model",
    "design cloud architecture for scalability",
    "perform data analysis on customer behavior",
    "simple routine data processing task",
    "technical implementation of authentication system"
  ];

  const results = [];

  for (const problem of problems) {
    const result = await solver.solveProblem(problem);
    results.push({
      problem,
      solvedBy: result.by,
      level: result.level
    });
    
    console.log(`✅ ${result.result}\n`);
  }

  // Analysis
  console.log("📊 Performance Analysis:");
  console.log("=" .repeat(40));
  results.forEach((r, i) => {
    console.log(`${i + 1}. ${r.problem}`);
    console.log(`   Solved by: ${r.solvedBy} (${r.level} level)`);
  });
}

advancedDemo();
```

## Key Benefits of Hierarchical Teams

1. **Scalability**: Can handle organizations of any size
2. **Clear Responsibility**: Each level has well-defined responsibilities
3. **Efficient Delegation**: Tasks are routed to appropriate expertise levels
4. **Error Containment**: Problems can be handled at the lowest possible level
5. **Career Progression**: Natural growth path for agents

## Common Use Cases

- **Corporate Organizations** (CEO → VPs → Directors → Managers → Employees)
- **Military Command Structures**
- **Large-scale Software Architecture**
- **Research and Development Teams**
- **Emergency Response Systems**
- **Complex Manufacturing Processes**

This pattern is excellent for scenarios where you need to model real-world organizational structures or manage complex systems with clear chains of command and specialization.