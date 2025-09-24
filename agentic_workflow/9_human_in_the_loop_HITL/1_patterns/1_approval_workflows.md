# Human-in-the-Loop (HITL) Patterns in Approval Workflows

Human-in-the-Loop (HITL) refers to systems where human judgment is integrated into automated processes. In approval workflows, HITL patterns ensure that critical decisions involve human oversight before final execution.

## Common HITL Patterns in Approval Workflows

1. **Sequential Approval** - Multiple approvers in a specific order
2. **Parallel Approval** - Multiple approvers simultaneously
3. **Conditional Routing** - Dynamic approval paths based on data
4. **Escalation** - Automatic escalation if no response
5. **Delegate/Proxy** - Approval delegation when primary is unavailable

## JavaScript Example: Simple Approval Workflow

Here's a basic implementation of a sequential approval workflow:

```javascript
class ApprovalWorkflow {
  constructor() {
    this.steps = [];
    this.currentStep = 0;
    this.status = 'pending'; // pending, approved, rejected, escalated
    this.data = {};
  }

  addStep(approver, conditions = {}) {
    this.steps.push({
      approver,
      conditions,
      status: 'pending',
      timestamp: null,
      comments: ''
    });
  }

  async startApproval(initialData) {
    this.data = { ...initialData };
    this.status = 'in_progress';
    console.log(`Approval workflow started for: ${this.data.title}`);
    
    return await this.processCurrentStep();
  }

  async processCurrentStep() {
    if (this.currentStep >= this.steps.length) {
      this.status = 'approved';
      console.log('All approval steps completed!');
      return this.finalizeApproval();
    }

    const step = this.steps[this.currentStep];
    console.log(`Waiting for approval from: ${step.approver}`);
    
    // In real implementation, this would trigger notification
    // and wait for human response via UI/email/etc.
    return this.simulateHumanDecision();
  }

  async submitDecision(decision, comments = '') {
    if (this.status !== 'in_progress') {
      throw new Error('No active approval step');
    }

    const step = this.steps[this.currentStep];
    step.status = decision;
    step.comments = comments;
    step.timestamp = new Date();

    if (decision === 'approved') {
      this.currentStep++;
      
      if (this.currentStep < this.steps.length) {
        return await this.processCurrentStep();
      } else {
        this.status = 'approved';
        return await this.finalizeApproval();
      }
    } else if (decision === 'rejected') {
      this.status = 'rejected';
      return await this.handleRejection();
    }
  }

  async finalizeApproval() {
    console.log('Approval completed successfully!');
    // Execute the approved action
    return { status: 'approved', data: this.data };
  }

  async handleRejection() {
    console.log('Approval rejected');
    // Handle rejection logic
    return { status: 'rejected', data: this.data };
  }

  // Simulate human decision-making (for demo purposes)
  async simulateHumanDecision() {
    return new Promise((resolve) => {
      setTimeout(() => {
        // Simulate 80% approval rate
        const decision = Math.random() > 0.2 ? 'approved' : 'rejected';
        this.submitDecision(decision, `Simulated ${decision} decision`);
        resolve(decision);
      }, 1000);
    });
  }
}

// Usage Example
async function runExample() {
  const workflow = new ApprovalWorkflow();
  
  // Define approval steps
  workflow.addStep('manager@company.com', { minAmount: 1000 });
  workflow.addStep('director@company.com', { minAmount: 5000 });
  workflow.addStep('finance@company.com', { minAmount: 10000 });
  
  // Start approval process
  const requestData = {
    title: 'New Software License Purchase',
    amount: 7500,
    requester: 'john.doe@company.com',
    description: 'Annual license for design software'
  };
  
  await workflow.startApproval(requestData);
}

runExample();
```

## Advanced HITL Pattern: Conditional Approval Routing

```javascript
class ConditionalApprovalWorkflow {
  constructor() {
    this.ruleEngine = new RuleEngine();
    this.approvers = new Map();
  }

  addApprover(role, email, conditions) {
    this.approvers.set(role, { email, conditions });
  }

  async determineApprovalPath(requestData) {
    const requiredApprovers = [];
    
    // Check amount-based conditions
    if (requestData.amount > 10000) {
      requiredApprovers.push('finance_director');
    } else if (requestData.amount > 5000) {
      requiredApprovers.push('department_head');
    }
    
    // Check department-specific rules
    if (requestData.department === 'IT' && requestData.amount > 2000) {
      requiredApprovers.push('it_director');
    }
    
    // Check risk level
    if (requestData.riskLevel === 'high') {
      requiredApprovers.push('compliance_officer');
    }
    
    return requiredApprovers;
  }

  async startConditionalApproval(requestData) {
    const approvalPath = await this.determineApprovalPath(requestData);
    console.log('Determined approval path:', approvalPath);
    
    const workflow = new ApprovalWorkflow();
    
    // Add steps based on conditional routing
    for (const role of approvalPath) {
      const approver = this.approvers.get(role);
      if (approver) {
        workflow.addStep(approver.email, approver.conditions);
      }
    }
    
    return await workflow.startApproval(requestData);
  }
}

// Rule Engine for complex conditions
class RuleEngine {
  evaluateRules(requestData, rules) {
    return rules.every(rule => {
      switch (rule.operator) {
        case 'greater_than':
          return requestData[rule.field] > rule.value;
        case 'equals':
          return requestData[rule.field] === rule.value;
        case 'in_array':
          return rule.value.includes(requestData[rule.field]);
        default:
          return true;
      }
    });
  }
}
```

## Real-world Implementation Considerations

1. **Notification System**: Integrate with email/Slack for approver notifications
2. **Timeout Handling**: Implement escalation for unresponsive approvers
3. **Audit Trail**: Maintain complete history of all decisions
4. **UI Integration**: Provide user-friendly interfaces for approvers
5. **Security**: Ensure proper authentication and authorization

This example demonstrates the core HITL patterns in approval workflows. In production systems, you would integrate with actual notification systems, databases, and user interfaces.