# State Persistence Strategies in JavaScript

State persistence involves saving application state to durable storage and restoring it later. This is crucial for user experience, allowing users to resume where they left off, maintain preferences, and recover from crashes.

## Core Persistence Strategies

### 1. **Storage Locations**
- **LocalStorage/SessionStorage**: Browser key-value storage
- **IndexedDB**: Browser database for larger data
- **Cookies**: Small data with expiration
- **Server/Cloud**: Remote persistence
- **File System**: Node.js file-based storage

### 2. **Persistence Patterns**
- **Auto-save**: Periodic background saving
- **On-demand**: Manual save points
- **Debounced**: Save after inactivity
- **Differential**: Only save changes

---

## JavaScript Examples

### Example 1: Multi-Storage Persistence Manager

```javascript
class PersistenceManager {
    constructor(options = {}) {
        this.strategies = new Map();
        this.defaultStrategy = options.defaultStrategy || 'localStorage';
        this.autoSave = options.autoSave || false;
        this.autoSaveInterval = options.autoSaveInterval || 30000; // 30 seconds
        this.debounceDelay = options.debounceDelay || 1000;
        
        this._initStrategies();
        this._setupAutoSave();
    }

    _initStrategies() {
        // LocalStorage Strategy
        this.addStrategy('localStorage', {
            save: (key, data) => {
                try {
                    localStorage.setItem(key, JSON.stringify({
                        data,
                        timestamp: Date.now(),
                        version: '1.0'
                    }));
                    return true;
                } catch (error) {
                    console.warn('LocalStorage save failed:', error);
                    return false;
                }
            },
            load: (key) => {
                try {
                    const item = localStorage.getItem(key);
                    if (!item) return null;
                    
                    const parsed = JSON.parse(item);
                    return this._migrateData(parsed.data, parsed.version);
                } catch (error) {
                    console.warn('LocalStorage load failed:', error);
                    return null;
                }
            },
            remove: (key) => {
                localStorage.removeItem(key);
            },
            getSize: () => {
                return JSON.stringify(localStorage).length;
            }
        });

        // SessionStorage Strategy
        this.addStrategy('sessionStorage', {
            save: (key, data) => {
                try {
                    sessionStorage.setItem(key, JSON.stringify(data));
                    return true;
                } catch (error) {
                    console.warn('SessionStorage save failed:', error);
                    return false;
                }
            },
            load: (key) => {
                try {
                    const item = sessionStorage.getItem(key);
                    return item ? JSON.parse(item) : null;
                } catch (error) {
                    console.warn('SessionStorage load failed:', error);
                    return null;
                }
            },
            remove: (key) => {
                sessionStorage.removeItem(key);
            }
        });

        // IndexedDB Strategy
        this.addStrategy('indexedDB', {
            save: async (key, data) => {
                return await this._indexedDBSave(key, data);
            },
            load: async (key) => {
                return await this._indexedDBLoad(key);
            },
            remove: async (key) => {
                return await this._indexedDBRemove(key);
            }
        });

        // Memory Strategy (fallback)
        this.addStrategy('memory', {
            _storage: new Map(),
            save: (key, data) => {
                this.strategies.get('memory')._storage.set(key, data);
                return true;
            },
            load: (key) => {
                return this.strategies.get('memory')._storage.get(key) || null;
            },
            remove: (key) => {
                this.strategies.get('memory')._storage.delete(key);
            }
        });
    }

    addStrategy(name, implementation) {
        this.strategies.set(name, implementation);
    }

    async save(key, data, strategy = this.defaultStrategy) {
        const strategyImpl = this.strategies.get(strategy);
        if (!strategyImpl) {
            throw new Error(`Unknown persistence strategy: ${strategy}`);
        }

        try {
            const result = await strategyImpl.save(key, data);
            this._emit('save', { key, strategy, success: result });
            return result;
        } catch (error) {
            this._emit('error', { key, strategy, error });
            throw error;
        }
    }

    async load(key, strategy = this.defaultStrategy) {
        const strategyImpl = this.strategies.get(strategy);
        if (!strategyImpl) {
            throw new Error(`Unknown persistence strategy: ${strategy}`);
        }

        try {
            const data = await strategyImpl.load(key);
            this._emit('load', { key, strategy, success: !!data });
            return data;
        } catch (error) {
            this._emit('error', { key, strategy, error });
            return null;
        }
    }

    async remove(key, strategy = this.defaultStrategy) {
        const strategyImpl = this.strategies.get(strategy);
        if (!strategyImpl) {
            throw new Error(`Unknown persistence strategy: ${strategy}`);
        }

        await strategyImpl.remove(key);
        this._emit('remove', { key, strategy });
    }

    // Fallback strategy loading
    async loadWithFallback(key, strategies = ['localStorage', 'sessionStorage', 'memory']) {
        for (const strategy of strategies) {
            const data = await this.load(key, strategy);
            if (data) {
                console.log(`Loaded from ${strategy} fallback`);
                return data;
            }
        }
        return null;
    }

    _setupAutoSave() {
        if (this.autoSave) {
            this.autoSaveTimer = setInterval(() => {
                this._emit('autoSave');
            }, this.autoSaveInterval);
        }
    }

    _migrateData(data, version) {
        // Data migration logic for different versions
        if (version === '0.9') {
            // Migrate from v0.9 to v1.0
            return { ...data, migrated: true };
        }
        return data;
    }

    // Simple event system
    _emit(event, data) {
        if (this.listeners && this.listeners[event]) {
            this.listeners[event].forEach(listener => listener(data));
        }
    }

    on(event, listener) {
        if (!this.listeners) this.listeners = {};
        if (!this.listeners[event]) this.listeners[event] = [];
        this.listeners[event].push(listener);
    }

    destroy() {
        if (this.autoSaveTimer) {
            clearInterval(this.autoSaveTimer);
        }
    }

    // IndexedDB implementation
    async _indexedDBSave(key, data) {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open('StatePersistence', 1);
            
            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                const db = request.result;
                const transaction = db.transaction(['states'], 'readwrite');
                const store = transaction.objectStore('states');
                const putRequest = store.put(data, key);
                
                putRequest.onsuccess = () => resolve(true);
                putRequest.onerror = () => reject(putRequest.error);
            };
            
            request.onupgradeneeded = (event) => {
                const db = event.target.result;
                if (!db.objectStoreNames.contains('states')) {
                    db.createObjectStore('states');
                }
            };
        });
    }

    async _indexedDBLoad(key) {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open('StatePersistence', 1);
            
            request.onerror = () => reject(request.error);
            request.onsuccess = () => {
                const db = request.result;
                const transaction = db.transaction(['states'], 'readonly');
                const store = transaction.objectStore('states');
                const getRequest = store.get(key);
                
                getRequest.onsuccess = () => resolve(getRequest.result || null);
                getRequest.onerror = () => reject(getRequest.error);
            };
        });
    }

    async _indexedDBRemove(key) {
        // Similar implementation to save but with delete operation
    }
}
```

### Example 2: Advanced State Container with Persistence

```javascript
class PersistentStateContainer {
    constructor(config = {}) {
        this.state = config.initialState || {};
        this.persistence = config.persistence || new PersistenceManager();
        this.persistenceKey = config.persistenceKey || 'app-state';
        
        this.debounceTimers = new Map();
        this.changeListeners = new Set();
        this.saveQueue = new Set();
        
        this.config = {
            autoSave: config.autoSave !== false,
            saveDebounce: config.saveDebounce ?? 1000,
            maxSaveAttempts: config.maxSaveAttempts ?? 3,
            compression: config.compression ?? false,
            ...config
        };

        this._initialize();
    }

    async _initialize() {
        // Load existing state
        await this.loadState();
        
        // Set up auto-save listeners
        if (this.config.autoSave) {
            this.persistence.on('autoSave', () => this.saveState());
            
            // Save on page unload
            window.addEventListener('beforeunload', () => {
                this.saveStateImmediate();
            });
        }
    }

    // State update with persistence
    update(updater, options = {}) {
        const oldState = this.state;
        const newState = typeof updater === 'function' ? updater(oldState) : updater;
        
        // Check if state actually changed
        if (this._isEqual(oldState, newState)) {
            return Promise.resolve(false);
        }

        this.state = newState;
        
        // Notify listeners
        this._notifyChangeListeners(newState, oldState);
        
        // Persist changes
        if (options.immediate) {
            return this.saveStateImmediate();
        } else if (options.persist !== false) {
            return this.debouncedSave();
        }
        
        return Promise.resolve(true);
    }

    // Debounced save to prevent excessive writes
    debouncedSave() {
        return new Promise((resolve) => {
            const key = this.persistenceKey;
            
            if (this.debounceTimers.has(key)) {
                clearTimeout(this.debounceTimers.get(key));
            }
            
            const timer = setTimeout(async () => {
                try {
                    await this._saveWithRetry();
                    resolve(true);
                } catch (error) {
                    console.error('Failed to save state:', error);
                    resolve(false);
                } finally {
                    this.debounceTimers.delete(key);
                }
            }, this.config.saveDebounce);
            
            this.debounceTimers.set(key, timer);
        });
    }

    // Immediate save without debouncing
    async saveStateImmediate() {
        try {
            await this._saveWithRetry();
            return true;
        } catch (error) {
            console.error('Immediate save failed:', error);
            return false;
        }
    }

    async _saveWithRetry(attempt = 1) {
        try {
            const data = this.config.compression ? this._compress(this.state) : this.state;
            await this.persistence.save(this.persistenceKey, data);
            console.log('State saved successfully');
        } catch (error) {
            if (attempt < this.config.maxSaveAttempts) {
                console.warn(`Save attempt ${attempt} failed, retrying...`);
                await new Promise(resolve => setTimeout(resolve, 1000 * attempt));
                return this._saveWithRetry(attempt + 1);
            }
            throw error;
        }
    }

    // Load state from persistence
    async loadState() {
        try {
            let data = await this.persistence.load(this.persistenceKey);
            
            if (data && this.config.compression) {
                data = this._decompress(data);
            }
            
            if (data) {
                const oldState = this.state;
                this.state = this._migrateState(data);
                this._notifyChangeListeners(this.state, oldState);
                console.log('State loaded successfully');
                return true;
            }
        } catch (error) {
            console.error('Failed to load state:', error);
        }
        
        return false;
    }

    // State migration for version updates
    _migrateState(data) {
        if (!data.version) {
            // Migrate from unversioned state
            return { ...data, version: '1.0', migrated: true };
        }
        
        // Add more migration logic as needed
        return data;
    }

    // Selective persistence - only save certain paths
    updatePath(path, value, options = {}) {
        const newState = this._setIn(this.state, path, value);
        return this.update(newState, options);
    }

    _setIn(obj, path, value) {
        const [key, ...rest] = path;
        
        if (rest.length === 0) {
            return { ...obj, [key]: value };
        }
        
        return {
            ...obj,
            [key]: this._setIn(obj[key] || {}, rest, value)
        };
    }

    // Compression utilities
    _compress(data) {
        // Simple compression - in real apps, use libraries like lz-string
        const jsonString = JSON.stringify(data);
        return btoa(unescape(encodeURIComponent(jsonString))); // Base64 encoding
    }

    _decompress(data) {
        try {
            const jsonString = decodeURIComponent(escape(atob(data)));
            return JSON.parse(jsonString);
        } catch (error) {
            console.warn('Decompression failed, returning raw data');
            return data;
        }
    }

    // Change listeners
    subscribe(listener) {
        this.changeListeners.add(listener);
        return () => this.changeListeners.delete(listener);
    }

    _notifyChangeListeners(newState, oldState) {
        this.changeListeners.forEach(listener => {
            try {
                listener(newState, oldState);
            } catch (error) {
                console.error('Change listener error:', error);
            }
        });
    }

    _isEqual(a, b) {
        return JSON.stringify(a) === JSON.stringify(b);
    }

    // Export/import state
    exportState() {
        return {
            state: this.state,
            timestamp: Date.now(),
            version: '1.0'
        };
    }

    async importState(data) {
        if (this._validateImportedState(data)) {
            return this.update(data.state, { immediate: true });
        }
        throw new Error('Invalid state format');
    }

    _validateImportedState(data) {
        return data && data.state && data.version;
    }

    // Cleanup
    destroy() {
        this.persistence.destroy();
        this.debounceTimers.forEach(timer => clearTimeout(timer));
        this.changeListeners.clear();
    }
}
```

### Example 3: Differential Persistence Strategy

```javascript
class DifferentialPersistence {
    constructor(persistenceManager, options = {}) {
        this.persistence = persistenceManager;
        this.key = options.key || 'app-state';
        this.patchSizeThreshold = options.patchSizeThreshold || 1024; // 1KB
        
        this.baseState = null;
        this.patches = [];
        this.lastSavedState = null;
    }

    async initialize() {
        // Load base state and patches
        this.baseState = await this.persistence.load(`${this.key}-base`) || {};
        const savedPatches = await this.persistence.load(`${this-key}-patches`) || [];
        this.patches = savedPatches;
        
        // Reconstruct current state
        this.lastSavedState = this._applyPatches(this.baseState, this.patches);
        return this.lastSavedState;
    }

    async save(newState) {
        if (!this.lastSavedState) {
            // First save - save complete state
            await this._saveCompleteState(newState);
            return;
        }

        // Calculate difference
        const patch = this._createPatch(this.lastSavedState, newState);
        
        if (this._shouldSaveComplete(patch, newState)) {
            // Patch is too large or too many patches - save complete state
            await this._saveCompleteState(newState);
        } else {
            // Save differential patch
            this.patches.push(patch);
            await this.persistence.save(`${this-key}-patches`, this.patches);
            this.lastSavedState = newState;
        }
    }

    _createPatch(oldState, newState) {
        const patch = {};
        
        // Simple diff algorithm - in practice, use a library like fast-json-patch
        this._buildPatch(oldState, newState, patch, '');
        
        return {
            timestamp: Date.now(),
            changes: patch,
            size: JSON.stringify(patch).length
        };
    }

    _buildPatch(oldObj, newObj, patch, path) {
        const allKeys = new Set([
            ...Object.keys(oldObj || {}),
            ...Object.keys(newObj || {})
        ]);

        for (const key of allKeys) {
            const currentPath = path ? `${path}.${key}` : key;
            const oldVal = oldObj?.[key];
            const newVal = newObj?.[key];

            if (JSON.stringify(oldVal) !== JSON.stringify(newVal)) {
                if (typeof newVal === 'object' && newVal !== null && 
                    typeof oldVal === 'object' && oldVal !== null) {
                    this._buildPatch(oldVal, newVal, patch, currentPath);
                } else {
                    patch[currentPath] = newVal;
                }
            }
        }
    }

    _applyPatches(baseState, patches) {
        let state = { ...baseState };
        
        for (const patch of patches) {
            state = this._applyPatch(state, patch.changes);
        }
        
        return state;
    }

    _applyPatch(state, changes) {
        const newState = { ...state };
        
        for (const [path, value] of Object.entries(changes)) {
            this._setValueByPath(newState, path, value);
        }
        
        return newState;
    }

    _setValueByPath(obj, path, value) {
        const keys = path.split('.');
        let current = obj;
        
        for (let i = 0; i < keys.length - 1; i++) {
            const key = keys[i];
            if (!current[key] || typeof current[key] !== 'object') {
                current[key] = {};
            }
            current = current[key];
        }
        
        current[keys[keys.length - 1]] = value;
    }

    _shouldSaveComplete(patch, newState) {
        // Save complete state if patch is too large
        if (patch.size > this.patchSizeThreshold) {
            return true;
        }
        
        // Save complete state if we have too many patches
        if (this.patches.length >= 50) {
            return true;
        }
        
        // Save complete state periodically
        if (this.patches.length > 0 && this.patches.length % 10 === 0) {
            return true;
        }
        
        return false;
    }

    async _saveCompleteState(state) {
        this.baseState = state;
        this.patches = [];
        this.lastSavedState = state;
        
        await this.persistence.save(`${this-key}-base`, state);
        await this.persistence.save(`${this-key}-patches`, this.patches);
    }

    async getCurrentState() {
        return this.lastSavedState;
    }

    // Time travel - get state at specific time
    async getStateAtTime(timestamp) {
        const relevantPatches = this.patches.filter(p => p.timestamp <= timestamp);
        return this._applyPatches(this.baseState, relevantPatches);
    }

    // Clean up old patches
    async cleanupOlderThan(timestamp) {
        const recentPatches = this.patches.filter(p => p.timestamp >= timestamp);
        
        if (recentPatches.length < this.patches.length) {
            this.patches = recentPatches;
            await this.persistence.save(`${this-key}-patches`, this.patches);
        }
    }
}
```

### Example 4: Practical Usage Example

```javascript
// Complete example using all strategies
class Application {
    constructor() {
        this.persistence = new PersistenceManager({
            autoSave: true,
            autoSaveInterval: 30000
        });

        this.stateContainer = new PersistentStateContainer({
            persistence: this.persistence,
            persistenceKey: 'todo-app-state',
            initialState: {
                todos: [],
                filter: 'all',
                userPreferences: {
                    theme: 'light',
                    language: 'en'
                },
                version: '1.0'
            },
            autoSave: true,
            saveDebounce: 2000
        });

        this.differentialPersistence = new DifferentialPersistence(
            this.persistence, 
            { key: 'todo-app-diff' }
        );

        this._setupEventListeners();
    }

    async initialize() {
        // Try differential loading first, fallback to full state
        try {
            const differentialState = await this.differentialPersistence.initialize();
            if (differentialState) {
                await this.stateContainer.update(() => differentialState);
            }
        } catch (error) {
            console.warn('Differential loading failed, using full state');
            await this.stateContainer.loadState();
        }

        this._render();
    }

    async addTodo(text) {
        await this.stateContainer.update(state => ({
            ...state,
            todos: [
                ...state.todos,
                {
                    id: Date.now(),
                    text,
                    completed: false,
                    createdAt: new Date().toISOString()
                }
            ]
        }));

        // Also save differentially for backup
        await this.differentialPersistence.save(this.stateContainer.state);
    }

    async toggleTodo(id) {
        await this.stateContainer.update(state => ({
            ...state,
            todos: state.todos.map(todo =>
                todo.id === id ? { ...todo, completed: !todo.completed } : todo
            )
        }));
    }

    async updatePreferences(preferences) {
        // Use path-based update for better performance
        await this.stateContainer.updatePath(
            ['userPreferences'], 
            { ...this.stateContainer.state.userPreferences, ...preferences }
        );
    }

    async exportData() {
        const exported = this.stateContainer.exportState();
        const blob = new Blob([JSON.stringify(exported, null, 2)], {
            type: 'application/json'
        });
        
        // Create download link
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `todo-backup-${new Date().toISOString().split('T')[0]}.json`;
        a.click();
        URL.revokeObjectURL(url);
    }

    async importData(file) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = async (e) => {
                try {
                    const data = JSON.parse(e.target.result);
                    await this.stateContainer.importState(data);
                    this._render();
                    resolve();
                } catch (error) {
                    reject(error);
                }
            };
            reader.readAsText(file);
        });
    }

    _setupEventListeners() {
        // Save on visibility change (tab switch)
        document.addEventListener('visibilitychange', () => {
            if (document.visibilityState === 'hidden') {
                this.stateContainer.saveStateImmediate();
            }
        });

        // Save on online status change
        window.addEventListener('online', () => {
            this.stateContainer.saveStateImmediate();
        });

        // Subscribe to state changes
        this.stateContainer.subscribe((newState, oldState) => {
            console.log('State changed:', { newState, oldState });
            this._render();
        });

        // Error handling
        this.persistence.on('error', (error) => {
            console.error('Persistence error:', error);
            this._showError('Failed to save data. Please check storage space.');
        });
    }

    _render() {
        // Update UI based on current state
        const state = this.stateContainer.state;
        console.log('Current state:', state);
        // ... actual UI rendering logic
    }

    _showError(message) {
        // Show error to user
        console.error('Application error:', message);
    }

    destroy() {
        this.stateContainer.destroy();
        this.persistence.destroy();
    }
}

// Usage
const app = new Application();
app.initialize();

// Example usage
app.addTodo('Learn state persistence');
app.updatePreferences({ theme: 'dark' });
```

## Key Persistence Strategies Summary

| Strategy | Use Case | Pros | Cons |
|----------|----------|------|------|
| **LocalStorage** | Small data, simple apps | Simple API, synchronous | 5MB limit, string-only |
| **IndexedDB** | Large data, complex queries | Large storage, async | Complex API |
| **Differential** | Frequent updates, large states | Efficient storage | Complex implementation |
| **Multi-strategy** | Production apps | Fallback options | More complex |
| **Auto-save** | User experience | No data loss | Performance impact |

## Best Practices

1. **Use debouncing** for frequent updates
2. **Implement fallback strategies** for reliability
3. **Compress large state** when possible
4. **Handle storage limits** gracefully
5. **Provide export/import** functionality
6. **Monitor performance** and optimize
7. **Implement data migration** for version changes

State persistence is essential for professional applications, providing users with seamless experiences across sessions and devices.