# State Channels in LangGraph

State channels are a powerful concept in LangGraph that enable parallel execution, conditional branching based on multiple state values, and complex workflow orchestration. They allow you to manage multiple state paths simultaneously and merge them back together.

## Core Concepts

### **What are State Channels?**
- **Parallel state paths**: Multiple state trajectories running concurrently
- **Conditional merging**: Combine channels based on conditions
- **Isolated execution**: Each channel maintains its own state context
- **Synchronization points**: Where channels merge or interact

### **Key Benefits**
- **Parallel processing**: Handle multiple tasks simultaneously
- **Conditional workflows**: Dynamic path selection based on state
- **Error isolation**: Failures in one channel don't affect others
- **Complex orchestration**: Manage sophisticated multi-step processes

---

## JavaScript Examples

### Example 1: Basic State Channel System

```javascript
class StateChannelSystem {
    constructor() {
        this.channels = new Map();
        this.activeChannels = new Set();
        this.mergeConditions = new Map();
        this.globalState = {};
    }

    // Create a new state channel
    createChannel(name, initialState = {}) {
        this.channels.set(name, {
            name,
            state: { ...initialState, _channel: name },
            isActive: true,
            createdAt: Date.now()
        });
        this.activeChannels.add(name);
        return this;
    }

    // Execute a function within a specific channel
    async executeInChannel(channelName, operationFn, metadata = {}) {
        const channel = this.channels.get(channelName);
        if (!channel || !channel.isActive) {
            throw new Error(`Channel ${channelName} not found or inactive`);
        }

        console.log(`🔹 Executing in channel ${channelName}: ${metadata.description || 'operation'}`);
        
        try {
            const result = await operationFn(channel.state);
            channel.state = { ...channel.state, ...result };
            channel.lastOperation = metadata.description;
            channel.updatedAt = Date.now();
            
            return result;
        } catch (error) {
            console.error(`Error in channel ${channelName}:`, error);
            channel.state._error = error.message;
            channel.state._failed = true;
            throw error;
        }
    }

    // Execute the same operation in multiple channels in parallel
    async executeInChannels(channelNames, operationFn, metadata = {}) {
        const executions = channelNames.map(channelName => 
            this.executeInChannel(channelName, operationFn, metadata)
        );
        
        return Promise.allSettled(executions);
    }

    // Define merge conditions between channels
    addMergeCondition(sourceChannel, targetChannel, conditionFn, metadata = {}) {
        const key = `${sourceChannel}-${targetChannel}`;
        if (!this.mergeConditions.has(key)) {
            this.mergeConditions.set(key, []);
        }
        
        this.mergeConditions.get(key).push({
            condition: conditionFn,
            metadata: { ...metadata, source: sourceChannel, target: targetChannel }
        });
        
        return this;
    }

    // Check and execute merge conditions
    async checkMergeConditions() {
        const merges = [];
        
        for (const [key, conditions] of this.mergeConditions) {
            const [sourceChannel, targetChannel] = key.split('-');
            const source = this.channels.get(sourceChannel);
            const target = this.channels.get(targetChannel);
            
            if (!source || !target || !source.isActive || !target.isActive) {
                continue;
            }
            
            for (const mergeCondition of conditions) {
                if (mergeCondition.condition(source.state, target.state)) {
                    console.log(`🔄 Merge condition met: ${sourceChannel} → ${targetChannel}`);
                    merges.push(this.mergeChannels(sourceChannel, targetChannel, mergeCondition.metadata));
                }
            }
        }
        
        return Promise.all(merges);
    }

    // Merge two channels
    async mergeChannels(sourceChannelName, targetChannelName, mergeConfig = {}) {
        const source = this.channels.get(sourceChannelName);
        const target = this.channels.get(targetChannelName);
        
        if (!source || !target) {
            throw new Error('One or both channels not found');
        }

        console.log(`🔄 Merging ${sourceChannelName} into ${targetChannelName}`);
        
        // Merge strategy can be customized
        const mergeStrategy = mergeConfig.strategy || 'combine';
        let mergedState = {};
        
        switch (mergeStrategy) {
            case 'source_wins':
                mergedState = { ...target.state, ...source.state };
                break;
            case 'target_wins':
                mergedState = { ...source.state, ...target.state };
                break;
            case 'combine':
                mergedState = this.deepMerge(target.state, source.state);
                break;
            case 'custom':
                mergedState = mergeConfig.customMerge 
                    ? mergeConfig.customMerge(source.state, target.state)
                    : this.deepMerge(target.state, source.state);
                break;
        }
        
        // Update target channel
        target.state = mergedState;
        target.state._mergedFrom = target.state._mergedFrom || [];
        target.state._mergedFrom.push({
            source: sourceChannelName,
            timestamp: Date.now(),
            strategy: mergeStrategy
        });
        
        // Deactivate source channel
        source.isActive = false;
        this.activeChannels.delete(sourceChannelName);
        
        return {
            merged: true,
            source: sourceChannelName,
            target: targetChannelName,
            strategy: mergeStrategy
        };
    }

    // Deep merge utility
    deepMerge(target, source) {
        const result = { ...target };
        
        for (const key in source) {
            if (source[key] && typeof source[key] === 'object' && !Array.isArray(source[key])) {
                result[key] = this.deepMerge(result[key] || {}, source[key]);
            } else {
                result[key] = source[key];
            }
        }
        
        return result;
    }

    // Get active channels
    getActiveChannels() {
        return Array.from(this.activeChannels).map(name => this.channels.get(name));
    }

    // Get channel state
    getChannelState(channelName) {
        const channel = this.channels.get(channelName);
        return channel ? channel.state : null;
    }

    // Close a channel
    closeChannel(channelName) {
        const channel = this.channels.get(channelName);
        if (channel) {
            channel.isActive = false;
            this.activeChannels.delete(channelName);
            channel.closedAt = Date.now();
        }
        return this;
    }

    // Visualize channel system
    visualize() {
        console.log('\n📊 STATE CHANNELS VISUALIZATION');
        console.log(`Active channels: ${this.activeChannels.size}`);
        
        this.channels.forEach((channel, name) => {
            const status = channel.isActive ? '🟢 ACTIVE' : '🔴 INACTIVE';
            console.log(`\n${name}: ${status}`);
            console.log(`  Operations: ${channel.lastOperation || 'None'}`);
            console.log(`  State keys: ${Object.keys(channel.state).join(', ')}`);
        });
        
        console.log('\nMerge Conditions:');
        this.mergeConditions.forEach((conditions, key) => {
            console.log(`  ${key}: ${conditions.length} condition(s)`);
        });
    }
}
```

### Example 2: Parallel Research with State Channels

```javascript
class ParallelResearchSystem {
    constructor() {
        this.channelSystem = new StateChannelSystem();
        this.researchTopics = new Map();
    }

    async conductParallelResearch(mainTopic, subtopics) {
        console.log(`🔬 Starting parallel research on: ${mainTopic}`);
        
        // Create main research channel
        this.channelSystem.createChannel('main_research', {
            topic: mainTopic,
            startTime: new Date().toISOString(),
            status: 'initializing'
        });

        // Create channels for each subtopic
        subtopics.forEach((subtopic, index) => {
            this.channelSystem.createChannel(`subtopic_${index}`, {
                topic: subtopic,
                parentTopic: mainTopic,
                startTime: new Date().toISOString(),
                status: 'researching'
            });
            
            this.researchTopics.set(`subtopic_${index}`, subtopic);
        });

        // Define merge conditions - when subtopics are complete
        subtopics.forEach((_, index) => {
            this.channelSystem.addMergeCondition(
                `subtopic_${index}`,
                'main_research',
                (sourceState, targetState) => {
                    return sourceState.researchComplete === true && 
                           sourceState.confidence > 0.7;
                },
                {
                    strategy: 'combine',
                    description: 'Merge completed research into main channel'
                }
            );
        });

        // Execute parallel research in all channels
        await this.executeParallelResearch();

        // Check for merges periodically
        await this.monitorAndMerge();

        // Get final results
        const finalState = this.channelSystem.getChannelState('main_research');
        return this.compileResearchReport(finalState);
    }

    async executeParallelResearch() {
        const researchOperations = [
            'gather_basic_info',
            'analyze_sources', 
            'extract_key_points',
            'validate_information',
            'synthesize_findings'
        ];

        // Execute each operation across all active channels
        for (const operation of researchOperations) {
            console.log(`\n⚡ Executing: ${operation}`);
            
            const activeChannels = this.channelSystem.getActiveChannels()
                .filter(ch => ch.name !== 'main_research')
                .map(ch => ch.name);

            await this.channelSystem.executeInChannels(
                activeChannels,
                async (state) => {
                    return await this.executeResearchStep(state, operation);
                },
                { description: operation }
            );

            // Brief pause between operations
            await new Promise(resolve => setTimeout(resolve, 100));
        }
    }

    async executeResearchStep(state, operation) {
        // Simulate different research operations
        switch (operation) {
            case 'gather_basic_info':
                return {
                    basicInfo: `Information about ${state.topic}`,
                    sources: ['Source A', 'Source B', 'Source C'],
                    researchStep: 'basic_info_gathered'
                };
            
            case 'analyze_sources':
                return {
                    sourceAnalysis: {
                        reliability: 0.8 + Math.random() * 0.2,
                        relevance: 0.7 + Math.random() * 0.3
                    },
                    researchStep: 'sources_analyzed'
                };
            
            case 'extract_key_points':
                const keyPoints = [
                    `Key point 1 about ${state.topic}`,
                    `Key point 2 about ${state.topic}`,
                    `Key point 3 about ${state.topic}`
                ];
                return {
                    keyPoints,
                    pointsCount: keyPoints.length,
                    researchStep: 'key_points_extracted'
                };
            
            case 'validate_information':
                const confidence = 0.6 + Math.random() * 0.4;
                return {
                    validation: {
                        isValid: confidence > 0.7,
                        confidence,
                        verifiedAt: new Date().toISOString()
                    },
                    researchStep: 'information_validated'
                };
            
            case 'synthesize_findings':
                return {
                    synthesis: `Comprehensive synthesis of research on ${state.topic}`,
                    researchComplete: true,
                    completedAt: new Date().toISOString(),
                    researchStep: 'findings_synthesized'
                };
            
            default:
                return { researchStep: 'unknown_operation' };
        }
    }

    async monitorAndMerge() {
        let attempts = 0;
        const maxAttempts = 10;
        
        while (attempts < maxAttempts) {
            attempts++;
            
            const activeChannels = this.channelSystem.getActiveChannels();
            const subtopicChannels = activeChannels.filter(ch => 
                ch.name.startsWith('subtopic_')
            );
            
            console.log(`\n📈 Research progress: ${subtopicChannels.length} subtopics active`);
            
            // Check if any channels can be merged
            await this.channelSystem.checkMergeConditions();
            
            // If only main channel remains, research is complete
            if (activeChannels.length === 1 && activeChannels[0].name === 'main_research') {
                console.log('✅ All research channels merged successfully');
                break;
            }
            
            // Wait before next check
            await new Promise(resolve => setTimeout(resolve, 200));
        }
        
        if (attempts >= maxAttempts) {
            console.log('⚠️ Research timeout - forcing channel closure');
            this.forceCloseRemainingChannels();
        }
    }

    forceCloseRemainingChannels() {
        const activeChannels = this.channelSystem.getActiveChannels();
        activeChannels.forEach(channel => {
            if (channel.name !== 'main_research') {
                console.log(`🔴 Force-closing channel: ${channel.name}`);
                this.channelSystem.closeChannel(channel.name);
                
                // Merge any partial results
                const partialResults = this.extractPartialResults(channel.state);
                this.channelSystem.executeInChannel('main_research', () => ({
                    partialResults: {
                        [channel.name]: partialResults,
                        mergedAt: new Date().toISOString(),
                        forced: true
                    }
                }));
            }
        });
    }

    extractPartialResults(state) {
        return {
            topic: state.topic,
            keyPoints: state.keyPoints || [],
            confidence: state.validation?.confidence || 0.5,
            status: 'partial'
        };
    }

    compileResearchReport(finalState) {
        console.log('\n📋 COMPILING RESEARCH REPORT');
        
        const report = {
            mainTopic: finalState.topic,
            completedAt: new Date().toISOString(),
            totalSubtopics: this.researchTopics.size,
            mergedChannels: finalState._mergedFrom || [],
            findings: this.extractFindings(finalState)
        };
        
        return report;
    }

    extractFindings(state) {
        const findings = [];
        
        // Extract from main research
        if (state.keyPoints) {
            findings.push(...state.keyPoints);
        }
        
        // Extract from merged channels
        if (state.partialResults) {
            Object.values(state.partialResults).forEach(partial => {
                if (partial.keyPoints) {
                    findings.push(...partial.keyPoints);
                }
            });
        }
        
        return [...new Set(findings)]; // Remove duplicates
    }
}

// Demonstration
async function demoParallelResearch() {
    const researchSystem = new ParallelResearchSystem();
    
    const mainTopic = "Artificial Intelligence Ethics";
    const subtopics = [
        "Bias in AI systems",
        "Privacy concerns with AI", 
        "AI accountability frameworks",
        "Ethical AI development guidelines"
    ];
    
    console.log('🧪 PARALLEL RESEARCH DEMONSTRATION\n');
    const report = await researchSystem.conductParallelResearch(mainTopic, subtopics);
    
    console.log('\n📊 FINAL RESEARCH REPORT:');
    console.log(JSON.stringify(report, null, 2));
    
    researchSystem.channelSystem.visualize();
}

// demoParallelResearch();
```

### Example 3: Multi-Stage Approval Workflow with Channels

```javascript
class MultiStageApprovalWorkflow {
    constructor() {
        this.channelSystem = new StateChannelSystem();
        this.workflowStages = ['submission', 'review', 'approval', 'implementation'];
    }

    async startApprovalProcess(submissionData) {
        console.log('🚀 Starting multi-stage approval workflow');
        
        // Create main workflow channel
        this.channelSystem.createChannel('main_workflow', {
            ...submissionData,
            workflowStage: 'submission',
            createdAt: new Date().toISOString(),
            currentApprovers: [],
            approvals: []
        });

        // Create parallel review channels for different departments
        const departments = ['technical', 'legal', 'business', 'security'];
        departments.forEach(dept => {
            this.channelSystem.createChannel(`${dept}_review`, {
                department: dept,
                parentWorkflow: 'main_workflow',
                status: 'pending_review',
                assignedTo: this.assignReviewer(dept),
                reviewStarted: new Date().toISOString()
            });
        });

        // Set up merge conditions for department reviews
        departments.forEach(dept => {
            this.channelSystem.addMergeCondition(
                `${dept}_review`,
                'main_workflow',
                (sourceState, targetState) => {
                    return sourceState.reviewComplete === true;
                },
                {
                    strategy: 'combine',
                    description: `Merge ${dept} review results`
                }
            );
        });

        // Execute the workflow
        await this.executeWorkflowStages();
        
        return this.getWorkflowResult();
    }

    assignReviewer(department) {
        const reviewers = {
            technical: 'tech_lead@company.com',
            legal: 'legal_counsel@company.com', 
            business: 'business_analyst@company.com',
            security: 'security_officer@company.com'
        };
        return reviewers[department] || 'default_reviewer@company.com';
    }

    async executeWorkflowStages() {
        // Stage 1: Parallel department reviews
        console.log('\n📋 STAGE 1: Parallel Department Reviews');
        await this.executeParallelReviews();
        
        // Stage 2: Check if all reviews are complete
        console.log('\n🔍 STAGE 2: Review Completion Check');
        const allReviewsComplete = await this.checkReviewCompletion();
        
        if (allReviewsComplete) {
            // Stage 3: Final approval
            console.log('\n✅ STAGE 3: Final Approval');
            await this.executeFinalApproval();
        } else {
            console.log('\n❌ STAGE 3: Workflow Rejection');
            await this.handleWorkflowRejection();
        }
    }

    async executeParallelReviews() {
        const reviewOperations = [
            'initial_assessment',
            'detailed_analysis', 
            'risk_evaluation',
            'recommendation'
        ];

        for (const operation of reviewOperations) {
            const activeChannels = this.channelSystem.getActiveChannels()
                .filter(ch => ch.name.endsWith('_review'))
                .map(ch => ch.name);

            await this.channelSystem.executeInChannels(
                activeChannels,
                async (state) => {
                    return await this.executeReviewStep(state, operation);
                },
                { description: operation }
            );

            // Simulate processing time
            await new Promise(resolve => setTimeout(resolve, 150));
        }
    }

    async executeReviewStep(state, operation) {
        const department = state.department;
        
        switch (operation) {
            case 'initial_assessment':
                return {
                    assessment: {
                        complexity: Math.random() > 0.5 ? 'high' : 'low',
                        impact: ['low', 'medium', 'high'][Math.floor(Math.random() * 3)],
                        departmentSpecific: `${department} initial assessment complete`
                    },
                    reviewProgress: 0.25
                };
            
            case 'detailed_analysis':
                const issuesFound = Math.floor(Math.random() * 3);
                return {
                    analysis: {
                        issuesFound,
                        criticalIssues: issuesFound > 0 ? Math.floor(Math.random() * issuesFound) : 0,
                        analysisComplete: true
                    },
                    reviewProgress: 0.5
                };
            
            case 'risk_evaluation':
                const riskLevel = Math.random() > 0.7 ? 'high' : 
                                 Math.random() > 0.4 ? 'medium' : 'low';
                return {
                    riskAssessment: {
                        level: riskLevel,
                        mitigation: riskLevel === 'high' ? 'requires_mitigation' : 'acceptable',
                        departmentRisk: `${department} risk: ${riskLevel}`
                    },
                    reviewProgress: 0.75
                };
            
            case 'recommendation':
                const recommendApprove = Math.random() > 0.3; // 70% approval rate
                return {
                    recommendation: {
                        decision: recommendApprove ? 'approve' : 'reject',
                        confidence: 0.6 + Math.random() * 0.4,
                        comments: `Recommend ${recommendApprove ? 'approval' : 'rejection'} from ${department}`
                    },
                    reviewComplete: true,
                    reviewProgress: 1.0,
                    completedAt: new Date().toISOString()
                };
            
            default:
                return { reviewProgress: state.reviewProgress || 0 };
        }
    }

    async checkReviewCompletion() {
        let attempts = 0;
        const maxAttempts = 5;
        
        while (attempts < maxAttempts) {
            attempts++;
            
            // Check merge conditions
            await this.channelSystem.checkMergeConditions();
            
            const mainState = this.channelSystem.getChannelState('main_workflow');
            const departmentReviews = ['technical', 'legal', 'business', 'security']
                .map(dept => mainState[`${dept}_review`])
                .filter(Boolean);
            
            console.log(`📊 Reviews completed: ${departmentReviews.length}/4`);
            
            if (departmentReviews.length === 4) {
                // All reviews merged
                const allApproved = departmentReviews.every(review => 
                    review.recommendation?.decision === 'approve'
                );
                
                await this.channelSystem.executeInChannel('main_workflow', () => ({
                    allReviewsComplete: true,
                    allApproved: allApproved,
                    reviewSummary: {
                        totalReviews: departmentReviews.length,
                        approved: departmentReviews.filter(r => r.recommendation?.decision === 'approve').length,
                        rejected: departmentReviews.filter(r => r.recommendation?.decision === 'reject').length
                    }
                }));
                
                return allApproved;
            }
            
            await new Promise(resolve => setTimeout(resolve, 200));
        }
        
        // Timeout - force closure
        this.forceClosePendingReviews();
        return false;
    }

    async executeFinalApproval() {
        await this.channelSystem.executeInChannel('main_workflow', (state) => {
            console.log('🎉 FINAL APPROVAL GRANTED');
            return {
                workflowStage: 'approved',
                finalDecision: 'approved',
                approvedAt: new Date().toISOString(),
                approvalAuthority: 'workflow_system',
                implementationStatus: 'pending'
            };
        });
    }

    async handleWorkflowRejection() {
        await this.channelSystem.executeInChannel('main_workflow', (state) => {
            console.log('❌ WORKFLOW REJECTED');
            return {
                workflowStage: 'rejected',
                finalDecision: 'rejected',
                rejectedAt: new Date().toISOString(),
                rejectionReasons: this.compileRejectionReasons(state)
            };
        });
    }

    compileRejectionReasons(state) {
        const reasons = [];
        const departments = ['technical', 'legal', 'business', 'security'];
        
        departments.forEach(dept => {
            const review = state[`${dept}_review`];
            if (review && review.recommendation?.decision === 'reject') {
                reasons.push({
                    department: dept,
                    reason: review.recommendation.comments,
                    riskLevel: review.riskAssessment?.level
                });
            }
        });
        
        return reasons;
    }

    forceClosePendingReviews() {
        const activeChannels = this.channelSystem.getActiveChannels();
        activeChannels.forEach(channel => {
            if (channel.name.endsWith('_review')) {
                console.log(`⏰ Force-closing pending review: ${channel.name}`);
                this.channelSystem.closeChannel(channel.name);
                
                // Mark as timeout
                this.channelSystem.executeInChannel('main_workflow', (state) => ({
                    [`${channel.name}_timeout`]: true,
                    timeoutOccurred: true
                }));
            }
        });
    }

    getWorkflowResult() {
        const finalState = this.channelSystem.getChannelState('main_workflow');
        
        return {
            decision: finalState.finalDecision,
            timestamp: finalState.approvedAt || finalState.rejectedAt,
            duration: this.calculateWorkflowDuration(finalState),
            reviewSummary: finalState.reviewSummary,
            reasons: finalState.rejectionReasons
        };
    }

    calculateWorkflowDuration(state) {
        const start = new Date(state.createdAt);
        const end = new Date(state.approvedAt || state.rejectedAt || new Date());
        return end - start;
    }
}

// Demonstration
async function demoApprovalWorkflow() {
    const workflow = new MultiStageApprovalWorkflow();
    
    const submissionData = {
        projectName: "New AI Feature Implementation",
        submitter: "product_team@company.com",
        description: "Implement new ML feature for user recommendations",
        budget: 50000,
        timeline: "3 months"
    };
    
    console.log('🏢 MULTI-STAGE APPROVAL WORKFLOW DEMO\n');
    const result = await workflow.startApprovalProcess(submissionData);
    
    console.log('\n📋 WORKFLOW RESULT:');
    console.log(JSON.stringify(result, null, 2));
    
    workflow.channelSystem.visualize();
}

// demoApprovalWorkflow();
```

### Example 4: Dynamic Channel Creation and Routing

```javascript
class DynamicChannelManagement {
    constructor() {
        this.channelSystem = new StateChannelSystem();
        this.channelTemplates = new Map();
        this.routingRules = new Map();
        this._initializeTemplates();
    }

    _initializeTemplates() {
        // Pre-defined channel templates
        this.channelTemplates.set('customer_support', {
            initialState: {
                type: 'support',
                priority: 'normal',
                status: 'open',
                createdAt: new Date().toISOString()
            },
            operations: ['triage', 'research', 'resolve', 'escalate']
        });

        this.channelTemplates.set('data_processing', {
            initialState: {
                type: 'processing',
                status: 'queued',
                dataSize: 0,
                processed: 0
            },
            operations: ['validate', 'transform', 'analyze', 'export']
        });

        this.channelTemplates.set('content_moderation', {
            initialState: {
                type: 'moderation',
                severity: 'unknown',
                status: 'pending',
                confidence: 0
            },
            operations: ['analyze', 'classify', 'review', 'action']
        });
    }

    // Dynamically create channels based on conditions
    async createDynamicChannels(triggerState, conditionConfigs) {
        const createdChannels = [];
        
        for (const config of conditionConfigs) {
            if (config.condition(triggerState)) {
                const channelName = this.generateChannelName(config.template);
                const template = this.channelTemplates.get(config.template);
                
                if (template) {
                    // Create channel with template and trigger state
                    this.channelSystem.createChannel(channelName, {
                        ...template.initialState,
                        ...config.initialData,
                        triggeredBy: triggerState._channel,
                        triggerReason: config.description
                    });
                    
                    createdChannels.push(channelName);
                    console.log(`🎯 Created dynamic channel: ${channelName} for ${config.description}`);
                }
            }
        }
        
        return createdChannels;
    }

    generateChannelName(template) {
        const timestamp = Date.now().toString(36);
        return `${template}_${timestamp}`;
    }

    // Add dynamic routing rules
    addRoutingRule(sourcePattern, targetGenerator, condition) {
        this.routingRules.set(sourcePattern, {
            targetGenerator,
            condition,
            description: `Route from ${sourcePattern} to dynamic target`
        });
    }

    // Execute dynamic routing
    async executeDynamicRouting() {
        const activeChannels = this.channelSystem.getActiveChannels();
        
        for (const channel of activeChannels) {
            for (const [pattern, rule] of this.routingRules) {
                if (channel.name.match(pattern) && rule.condition(channel.state)) {
                    const targetChannel = rule.targetGenerator(channel.state);
                    await this.routeChannel(channel.name, targetChannel, rule.description);
                }
            }
        }
    }

    async routeChannel(sourceChannel, targetChannel, reason) {
        console.log(`🔄 Routing ${sourceChannel} → ${targetChannel}: ${reason}`);
        
        // Close current channel and create/merge with target
        this.channelSystem.closeChannel(sourceChannel);
        
        if (!this.channelSystem.channels.has(targetChannel)) {
            // Create new target channel
            this.channelSystem.createChannel(targetChannel, {
                routedFrom: sourceChannel,
                routingReason: reason,
                createdAt: new Date().toISOString()
            });
        } else {
            // Merge with existing channel
            await this.channelSystem.mergeChannels(sourceChannel, targetChannel, {
                strategy: 'combine',
                description: `Dynamic routing: ${reason}`
            });
        }
    }

    // Example usage: Customer support escalation system
    async handleCustomerRequest(request) {
        // Create main support channel
        const mainChannel = 'support_main';
        this.channelSystem.createChannel(mainChannel, {
            ...request,
            type: 'customer_support',
            status: 'new',
            priority: this.calculatePriority(request)
        });

        // Define dynamic channel creation rules
        const channelConditions = [
            {
                condition: (state) => state.priority === 'high' && state.complexity === 'technical',
                template: 'technical_support',
                description: 'High priority technical issue',
                initialData: { escalationLevel: 1, specializedTeam: 'tech_support' }
            },
            {
                condition: (state) => state.urgency === 'critical',
                template: 'crisis_management', 
                description: 'Critical situation requiring immediate attention',
                initialData: { crisisLevel: 'high', emergencyProtocol: true }
            },
            {
                condition: (state) => state.requiresApproval === true,
                template: 'approval_workflow',
                description: 'Requires managerial approval',
                initialData: { approvalLevel: 'manager', workflowType: 'standard' }
            }
        ];

        // Create dynamic channels based on request analysis
        const mainState = this.channelSystem.getChannelState(mainChannel);
        const dynamicChannels = await this.createDynamicChannels(mainState, channelConditions);

        // Set up routing rules
        this.addRoutingRule(
            /support_.*/,
            (state) => {
                if (state.resolutionAttempts >= 3) return 'escalation_channel';
                if (state.customerSatisfaction < 0.3) return 'retention_specialist';
                return 'standard_processing';
            },
            (state) => state.status === 'processing'
        );

        // Execute the support workflow
        await this.executeSupportWorkflow(mainChannel, dynamicChannels);

        return this.compileSupportResult(mainChannel);
    }

    calculatePriority(request) {
        if (request.urgency === 'critical') return 'high';
        if (request.impact === 'widespread') return 'high';
        if (request.complexity === 'simple') return 'low';
        return 'normal';
    }

    async executeSupportWorkflow(mainChannel, dynamicChannels) {
        // Execute operations in all channels
        const allChannels = [mainChannel, ...dynamicChannels];
        
        const workflowStages = ['initial_triage', 'investigation', 'resolution_attempt', 'follow_up'];
        
        for (const stage of workflowStages) {
            console.log(`\n⚡ Workflow stage: ${stage}`);
            
            await this.channelSystem.executeInChannels(
                allChannels,
                async (state) => await this.executeSupportStage(state, stage),
                { description: stage }
            );

            // Apply dynamic routing after each stage
            await this.executeDynamicRouting();
            
            // Check for completion conditions
            await this.checkCompletionConditions();
            
            await new Promise(resolve => setTimeout(resolve, 100));
        }
    }

    async executeSupportStage(state, stage) {
        // Simulate different support stages
        switch (stage) {
            case 'initial_triage':
                return {
                    triageComplete: true,
                    initialAssessment: `Assessed as ${state.priority} priority`,
                    assignedAgent: `agent_${Math.floor(Math.random() * 100)}`
                };
            
            case 'investigation':
                const issuesFound = Math.floor(Math.random() * 5);
                return {
                    investigationComplete: true,
                    issuesIdentified: issuesFound,
                    rootCause: issuesFound > 0 ? 'identified' : 'none'
                };
            
            case 'resolution_attempt':
                const success = Math.random() > 0.4; // 60% success rate
                const attempts = (state.resolutionAttempts || 0) + 1;
                return {
                    resolutionAttempts: attempts,
                    lastAttemptSuccess: success,
                    customerSatisfaction: success ? 0.8 : 0.3,
                    status: success ? 'resolved' : 'escalation_needed'
                };
            
            case 'follow_up':
                return {
                    followUpComplete: true,
                    customerFeedback: 'positive',
                    caseClosed: true,
                    closedAt: new Date().toISOString()
                };
            
            default:
                return { stage: 'unknown' };
        }
    }

    async checkCompletionConditions() {
        const activeChannels = this.channelSystem.getActiveChannels();
        
        for (const channel of activeChannels) {
            if (channel.state.caseClosed) {
                console.log(`✅ Channel ${channel.name} completed successfully`);
                this.channelSystem.closeChannel(channel.name);
                
                // Merge results if this was a dynamic channel
                if (channel.name !== 'support_main') {
                    await this.channelSystem.mergeChannels(
                        channel.name, 
                        'support_main', 
                        { strategy: 'combine' }
                    );
                }
            }
        }
    }

    compileSupportResult(mainChannel) {
        const finalState = this.channelSystem.getChannelState(mainChannel);
        
        return {
            success: finalState.caseClosed === true,
            resolutionTime: this.calculateResolutionTime(finalState),
            customerSatisfaction: finalState.customerSatisfaction,
            channelsInvolved: finalState._mergedFrom?.length || 0,
            finalStatus: finalState.status
        };
    }

    calculateResolutionTime(state) {
        const start = new Date(state.createdAt);
        const end = new Date(state.closedAt || new Date());
        return end - start;
    }
}
```

## Key State Channel Patterns

### **1. Parallel Execution Pattern**
```javascript
// Multiple channels process different aspects simultaneously
const channels = ['analysis', 'validation', 'enrichment'];
await executeInParallel(channels, processData);
```

### **2. Conditional Merging Pattern**
```javascript
// Merge channels when conditions are met
addMergeCondition('channel_a', 'main', (stateA, stateMain) => {
    return stateA.complete && stateA.confidence > 0.8;
});
```

### **3. Dynamic Channel Creation**
```javascript
// Create channels based on runtime conditions
if (needsSpecializedProcessing(state)) {
    createChannel('specialized_processing', specializedState);
}
```

### **4. Hierarchical Channel Structure**
```javascript
// Parent-child channel relationships
main_channel
├── technical_review
├── business_review  
└── compliance_review
```

## Benefits of State Channels

1. **Parallel Processing**: Handle multiple tasks simultaneously
2. **Fault Isolation**: Errors in one channel don't affect others
3. **Dynamic Workflows**: Adapt based on real-time conditions
4. **Resource Optimization**: Distribute work based on channel priorities
5. **Complex Orchestration**: Manage sophisticated multi-step processes

## Best Practices

1. **Clear Channel Purposes**: Each channel should have a specific responsibility
2. **Defined Merge Conditions**: Know when and how channels should combine
3. **Resource Management**: Monitor channel count and resource usage
4. **Error Handling**: Plan for channel failures and timeouts
5. **Monitoring**: Track channel states and execution paths

State channels enable complex, adaptive workflows that can handle real-world scenarios where multiple parallel processes need to be coordinated and merged based on dynamic conditions.