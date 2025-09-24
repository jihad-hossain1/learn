# Correction Interfaces in Human-in-the-Loop Systems

Correction interfaces are specialized UI components that allow humans to review, modify, and correct automated system outputs before final processing. They're essential for maintaining data quality and system reliability in HITL workflows.

## Key Characteristics of Effective Correction Interfaces

1. **Context Preservation** - Show original input and system reasoning
2. **Minimal Friction** - Reduce cognitive load for correctors
3. **Audit Trail** - Track all changes and decisions
4. **Validation** - Ensure corrected data meets requirements
5. **Batch Operations** - Support bulk corrections when possible

## Common Correction Interface Patterns

### 1. **Inline Editing Interface**
```javascript
class InlineCorrectionInterface {
  constructor(containerId, data) {
    this.container = document.getElementById(containerId);
    this.originalData = { ...data };
    this.correctedData = { ...data };
    this.validationRules = {};
  }

  render() {
    this.container.innerHTML = this.createInterface();
    this.attachEventListeners();
  }

  createInterface() {
    return `
      <div class="correction-interface">
        <div class="header">
          <h3>Review and Correct Data</h3>
          <div class="confidence-score">
            System Confidence: ${this.originalData.confidence}%
          </div>
        </div>
        
        <div class="correction-area">
          ${this.createFieldCorrections()}
        </div>
        
        <div class="actions">
          <button class="btn-approve">Approve As Is</button>
          <button class="btn-submit-corrections">Submit Corrections</button>
          <button class="btn-reject">Reject Entire Record</button>
        </div>
        
        <div class="audit-trail" style="display: none;">
          <h4>Change History</h4>
          <div id="change-history"></div>
        </div>
      </div>
    `;
  }

  createFieldCorrections() {
    return Object.keys(this.originalData.fields).map(fieldName => {
      const field = this.originalData.fields[fieldName];
      return `
        <div class="field-correction" data-field="${fieldName}">
          <label>${field.label}</label>
          <div class="original-value">
            <span class="value">${field.value}</span>
            <span class="confidence">(${field.confidence}%)</span>
          </div>
          <div class="correction-controls">
            <input type="text" 
                   value="${field.value}" 
                   class="correction-input"
                   data-original="${field.value}">
            <button class="btn-suggest" data-field="${fieldName}">AI Suggest</button>
            <button class="btn-reset" data-field="${fieldName}">Reset</button>
          </div>
          <div class="validation-error" style="display: none;"></div>
        </div>
      `;
    }).join('');
  }

  attachEventListeners() {
    // Input change handlers
    this.container.querySelectorAll('.correction-input').forEach(input => {
      input.addEventListener('change', (e) => this.handleFieldCorrection(e));
    });

    // Action buttons
    this.container.querySelector('.btn-approve').addEventListener('click', () => this.approveAsIs());
    this.container.querySelector('.btn-submit-corrections').addEventListener('click', () => this.submitCorrections());
    this.container.querySelector('.btn-reject').addEventListener('click', () => this.rejectRecord());

    // Field-specific buttons
    this.container.querySelectorAll('.btn-suggest').forEach(btn => {
      btn.addEventListener('click', (e) => this.getAISuggestion(e.target.dataset.field));
    });

    this.container.querySelectorAll('.btn-reset').forEach(btn => {
      btn.addEventListener('click', (e) => this.resetField(e.target.dataset.field));
    });
  }

  handleFieldCorrection(event) {
    const fieldName = event.target.closest('.field-correction').dataset.field;
    const newValue = event.target.value;
    const originalValue = event.target.dataset.original;

    this.correctedData.fields[fieldName].value = newValue;
    this.validateField(fieldName, newValue);
    this.recordChange(fieldName, originalValue, newValue);
  }

  async getAISuggestion(fieldName) {
    // Simulate AI suggestion API call
    const suggestion = await this.callAISuggestionAPI(fieldName, this.correctedData);
    
    const input = this.container.querySelector(`[data-field="${fieldName}"] .correction-input`);
    input.value = suggestion;
    this.handleFieldCorrection({ target: input });
  }

  resetField(fieldName) {
    const originalValue = this.originalData.fields[fieldName].value;
    const input = this.container.querySelector(`[data-field="${fieldName}"] .correction-input`);
    input.value = originalValue;
    this.handleFieldCorrection({ target: input });
  }

  validateField(fieldName, value) {
    const rules = this.validationRules[fieldName];
    const errorElement = this.container.querySelector(`[data-field="${fieldName}"] .validation-error`);

    if (rules) {
      for (const rule of rules) {
        if (!rule.validator(value)) {
          errorElement.textContent = rule.message;
          errorElement.style.display = 'block';
          return false;
        }
      }
    }

    errorElement.style.display = 'none';
    return true;
  }

  recordChange(fieldName, fromValue, toValue) {
    const change = {
      timestamp: new Date().toISOString(),
      field: fieldName,
      from: fromValue,
      to: toValue,
      user: 'current_user' // Would be actual user in real implementation
    };

    const historyElement = this.container.querySelector('#change-history');
    historyElement.innerHTML += `
      <div class="change-record">
        [${change.timestamp}] ${fieldName}: "${fromValue}" → "${toValue}"
      </div>
    `;
  }

  approveAsIs() {
    this.submitCorrection({
      action: 'approved',
      corrections: {},
      timestamp: new Date().toISOString()
    });
  }

  submitCorrections() {
    const hasErrors = Object.keys(this.correctedData.fields).some(fieldName => 
      !this.validateField(fieldName, this.correctedData.fields[fieldName].value)
    );

    if (!hasErrors) {
      this.submitCorrection({
        action: 'corrected',
        corrections: this.correctedData.fields,
        timestamp: new Date().toISOString()
      });
    }
  }

  rejectRecord() {
    this.submitCorrection({
      action: 'rejected',
      reason: 'manual_rejection',
      timestamp: new Date().toISOString()
    });
  }

  submitCorrection(correctionResult) {
    // In real implementation, this would send to backend
    console.log('Correction submitted:', correctionResult);
    
    // Dispatch custom event for other components to listen to
    this.container.dispatchEvent(new CustomEvent('correctionComplete', {
      detail: correctionResult
    }));
  }

  async callAISuggestionAPI(fieldName, context) {
    // Simulate API call
    return new Promise(resolve => {
      setTimeout(() => {
        const suggestions = {
          name: "John Smith", // AI might suggest based on context
          email: "john.smith@example.com",
          date: new Date().toISOString().split('T')[0]
        };
        resolve(suggestions[fieldName] || "Suggested value");
      }, 500);
    });
  }
}

// Usage Example
const sampleData = {
  confidence: 85,
  fields: {
    name: { label: "Full Name", value: "Jon Smth", confidence: 78 },
    email: { label: "Email Address", value: "jon.smth@example", confidence: 92 },
    date: { label: "Event Date", value: "2024-13-45", confidence: 45 }
  }
};

const correctionUI = new InlineCorrectionInterface('correction-container', sampleData);
correctionUI.validationRules = {
  email: [
    {
      validator: (value) => /\S+@\S+\.\S+/.test(value),
      message: "Please enter a valid email address"
    }
  ],
  date: [
    {
      validator: (value) => !isNaN(Date.parse(value)),
      message: "Please enter a valid date"
    }
  ]
};

correctionUI.render();
```

### 2. **Side-by-Side Comparison Interface**
```javascript
class ComparisonCorrectionInterface {
  constructor(containerId, originalData, suggestedCorrections) {
    this.container = document.getElementById(containerId);
    this.originalData = originalData;
    this.suggestedCorrections = suggestedCorrections;
  }

  render() {
    this.container.innerHTML = this.createComparisonView();
    this.attachEventListeners();
  }

  createComparisonView() {
    return `
      <div class="comparison-interface">
        <div class="comparison-header">
          <div class="column">Original</div>
          <div class="column">Suggested Correction</div>
          <div class="column">Your Decision</div>
        </div>
        
        ${this.createComparisonRows()}
        
        <div class="bulk-actions">
          <button class="btn-accept-all">Accept All Suggestions</button>
          <button class="btn-reject-all">Reject All Suggestions</button>
        </div>
      </div>
    `;
  }

  createComparisonRows() {
    return Object.keys(this.originalData).map(field => `
      <div class="comparison-row" data-field="${field}">
        <div class="column original-value">
          ${this.originalData[field]}
        </div>
        <div class="column suggested-value">
          ${this.suggestedCorrections[field] || this.originalData[field]}
        </div>
        <div class="column decision">
          <label>
            <input type="radio" name="${field}" value="keep" checked> Keep Original
          </label>
          <label>
            <input type="radio" name="${field}" value="accept"> Accept Suggestion
          </label>
          <label>
            <input type="radio" name="${field}" value="custom"> Custom Value
          </label>
          <input type="text" class="custom-value" style="display: none;" 
                 placeholder="Enter custom value">
        </div>
      </div>
    `).join('');
  }

  attachEventListeners() {
    // Radio button handlers
    this.container.querySelectorAll('input[type="radio"]').forEach(radio => {
      radio.addEventListener('change', (e) => this.handleDecisionChange(e));
    });

    // Bulk actions
    this.container.querySelector('.btn-accept-all').addEventListener('click', () => this.acceptAll());
    this.container.querySelector('.btn-reject-all').addEventListener('click', () => this.rejectAll());
  }

  handleDecisionChange(event) {
    const row = event.target.closest('.comparison-row');
    const customInput = row.querySelector('.custom-value');

    if (event.target.value === 'custom') {
      customInput.style.display = 'block';
    } else {
      customInput.style.display = 'none';
    }
  }

  acceptAll() {
    this.container.querySelectorAll('input[value="accept"]').forEach(radio => {
      radio.checked = true;
      radio.dispatchEvent(new Event('change'));
    });
  }

  rejectAll() {
    this.container.querySelectorAll('input[value="keep"]').forEach(radio => {
      radio.checked = true;
      radio.dispatchEvent(new Event('change'));
    });
  }

  getFinalDecisions() {
    const decisions = {};
    
    this.container.querySelectorAll('.comparison-row').forEach(row => {
      const field = row.dataset.field;
      const selectedOption = row.querySelector('input[type="radio"]:checked').value;
      
      decisions[field] = {
        decision: selectedOption,
        finalValue: selectedOption === 'custom' 
          ? row.querySelector('.custom-value').value
          : selectedOption === 'accept'
          ? this.suggestedCorrections[field]
          : this.originalData[field]
      };
    });

    return decisions;
  }
}
```

## Best Practices for Correction Interfaces

### 1. **User Experience Considerations**
```javascript
class CorrectionUXOptimizations {
  constructor() {
    this.keyboardShortcuts = new Map([
      ['Tab', 'nextField'],
      ['Shift+Tab', 'previousField'],
      ['Enter', 'submitCorrection'],
      ['Escape', 'cancelCorrection']
    ]);
  }

  setupKeyboardNavigation() {
    document.addEventListener('keydown', (e) => {
      const shortcut = this.getShortcutKey(e);
      if (this.keyboardShortcuts.has(shortcut)) {
        e.preventDefault();
        this.executeShortcut(this.keyboardShortcuts.get(shortcut));
      }
    });
  }

  setupProgressiveDisclosure() {
    // Show advanced options only when needed
    const advancedToggle = document.getElementById('advanced-options-toggle');
    const advancedSection = document.getElementById('advanced-options');
    
    advancedToggle.addEventListener('click', () => {
      advancedSection.style.display = 
        advancedSection.style.display === 'none' ? 'block' : 'none';
    });
  }

  setupAutoSave() {
    setInterval(() => {
      this.saveCorrectionProgress();
    }, 30000); // Auto-save every 30 seconds
  }

  provideVisualFeedback(correctionElement, isValid) {
    correctionElement.style.borderColor = isValid ? '#28a745' : '#dc3545';
    
    if (isValid) {
      correctionElement.classList.add('correction-valid');
      correctionElement.classList.remove('correction-invalid');
    } else {
      correctionElement.classList.add('correction-invalid');
      correctionElement.classList.remove('correction-valid');
    }
  }
}
```

### 2. **Batch Correction Interface**
```javascript
class BatchCorrectionInterface {
  constructor(containerId, records) {
    this.records = records;
    this.container = document.getElementById(containerId);
    this.correctionTemplates = new Map();
  }

  setupBulkOperations() {
    // Pattern-based corrections
    this.correctionTemplates.set('date-format', {
      name: 'Fix Date Format',
      apply: (value) => value.replace(/(\d{4})-(\d{2})-(\d{2})/, '$2/$3/$1')
    });

    this.correctionTemplates.set('trim-whitespace', {
      name: 'Trim Whitespace',
      apply: (value) => value.trim()
    });

    this.correctionTemplates.set('standardize-phone', {
      name: 'Standardize Phone Format',
      apply: (value) => value.replace(/\D/g, '').replace(/(\d{3})(\d{3})(\d{4})/, '($1) $2-$3')
    });
  }

  applyBulkCorrection(templateId, fieldName) {
    const template = this.correctionTemplates.get(templateId);
    this.records.forEach(record => {
      if (record.fields[fieldName]) {
        record.fields[fieldName].value = template.apply(record.fields[fieldName].value);
      }
    });
    this.render();
  }

  createBulkActionPanel() {
    return `
      <div class="bulk-action-panel">
        <select id="bulk-action-select">
          <option value="">Select Bulk Action</option>
          ${Array.from(this.correctionTemplates.entries()).map(([id, template]) => 
            `<option value="${id}">${template.name}</option>`
          ).join('')}
        </select>
        
        <select id="bulk-field-select">
          <option value="">Select Field</option>
          ${Object.keys(this.records[0].fields).map(field => 
            `<option value="${field}">${field}</option>`
          ).join('')}
        </select>
        
        <button onclick="applyBulkAction()">Apply to All</button>
        
        <div class="bulk-preview">
          <h4>Preview Changes</h4>
          <div id="bulk-preview-content"></div>
        </div>
      </div>
    `;
  }
}
```

## Key Implementation Considerations

1. **Performance**: Handle large datasets efficiently with virtualization
2. **Accessibility**: Ensure keyboard navigation and screen reader support
3. **Undo/Redo**: Implement change history for user confidence
4. **Collaboration**: Support multiple correctors with conflict resolution
5. **Integration**: Connect with backend systems for data persistence

These correction interface patterns provide robust foundations for building effective human-in-the-loop systems that balance automation with human oversight.