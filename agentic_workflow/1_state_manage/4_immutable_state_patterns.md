# Immutable State Patterns in JavaScript

Immutable state patterns are programming techniques where state objects are never modified after creation. Instead, any "change" creates a new object with the updated values. This is crucial for predictable state management, especially in React, Redux, and functional programming.

## Core Concepts

### Why Immutability?
- **Predictability**: State changes are explicit and traceable
- **Performance Optimization**: Easy change detection (=== comparison)
- **Thread Safety**: No race conditions in concurrent environments
- **Debugging**: State history is preserved
- **Undo/Redo**: Previous states are naturally available

---

## JavaScript Examples

### Example 1: Basic Object Immutability

```javascript
// MUTABLE approach (what NOT to do)
const mutableUser = { name: "John", age: 30 };
mutableUser.age = 31; // Original object is modified

// IMMUTABLE approaches
class ImmutablePatterns {
    // 1. Object Spread Operator (ES6)
    static updateWithSpread(obj, updates) {
        return { ...obj, ...updates };
    }

    // 2. Object.assign (ES5)
    static updateWithAssign(obj, updates) {
        return Object.assign({}, obj, updates);
    }

    // 3. Deep update with nested objects
    static deepUpdate(obj, path, value) {
        const [key, ...rest] = path;
        
        if (rest.length === 0) {
            // Base case: update the final property
            return { ...obj, [key]: value };
        }
        
        // Recursive case: navigate deeper
        return {
            ...obj,
            [key]: this.deepUpdate(obj[key] || {}, rest, value)
        };
    }

    // 4. Delete property immutably
    static deleteProperty(obj, property) {
        const { [property]: deleted, ...rest } = obj;
        return rest;
    }
}

// Usage examples
const originalUser = { name: "John", age: 30, address: { city: "NYC" } };

console.log("=== SPREAD OPERATOR ===");
const updatedUser = ImmutablePatterns.updateWithSpread(originalUser, { age: 31 });
console.log("Original:", originalUser); // { name: "John", age: 30, ... }
console.log("Updated:", updatedUser);   // { name: "John", age: 31, ... }
console.log("Same object?", originalUser === updatedUser); // false

console.log("\n=== DEEP UPDATE ===");
const deepUpdated = ImmutablePatterns.deepUpdate(originalUser, ['address', 'city'], 'Boston');
console.log("Original address:", originalUser.address.city); // "NYC"
console.log("Updated address:", deepUpdated.address.city);   // "Boston"

console.log("\n=== DELETE PROPERTY ===");
const withoutAge = ImmutablePatterns.deleteProperty(originalUser, 'age');
console.log("Original has age?", 'age' in originalUser); // true
console.log("New has age?", 'age' in withoutAge);        // false
```

### Example 2: Array Immutability

```javascript
class ImmutableArrays {
    // Add item
    static addItem(array, item) {
        return [...array, item];
    }

    // Add item at beginning
    static prependItem(array, item) {
        return [item, ...array];
    }

    // Remove item by index
    static removeByIndex(array, index) {
        return [...array.slice(0, index), ...array.slice(index + 1)];
    }

    // Remove item by value (first occurrence)
    static removeByValue(array, value) {
        const index = array.indexOf(value);
        return index === -1 ? array : this.removeByIndex(array, index);
    }

    // Update item by index
    static updateByIndex(array, index, updater) {
        return array.map((item, i) => 
            i === index ? (typeof updater === 'function' ? updater(item) : updater) : item
        );
    }

    // Insert item at specific position
    static insertAt(array, index, item) {
        return [...array.slice(0, index), item, ...array.slice(index)];
    }

    // Sort immutably
    static sortImmutable(array, compareFn) {
        return [...array].sort(compareFn);
    }

    // Filter with predicate
    static filterImmutable(array, predicate) {
        return array.filter(predicate);
    }
}

// Usage examples
const originalNumbers = [1, 2, 3, 4];

console.log("=== ARRAY OPERATIONS ===");
console.log("Original:", originalNumbers);

const withFive = ImmutableArrays.addItem(originalNumbers, 5);
console.log("After adding 5:", withFive);
console.log("Original unchanged:", originalNumbers);

const withoutThird = ImmutableArrays.removeByIndex(originalNumbers, 2);
console.log("After removing index 2:", withoutThird);

const doubledThird = ImmutableArrays.updateByIndex(originalNumbers, 2, n => n * 2);
console.log("After doubling index 2:", doubledThird);

const inserted = ImmutableArrays.insertAt(originalNumbers, 1, 99);
console.log("After inserting at index 1:", inserted);
```

### Example 3: Advanced Immutable State Management

```javascript
// Immutable state container pattern
class ImmutableStore {
    constructor(initialState = {}) {
        this._state = initialState;
        this._history = [initialState];
        this._listeners = new Set();
    }

    get state() {
        return this._state;
    }

    get history() {
        return [...this._history];
    }

    // Core update method
    update(updater) {
        const oldState = this._state;
        const newState = typeof updater === 'function' 
            ? updater(oldState) 
            : { ...oldState, ...updater };

        // Only update if state actually changed
        if (!this.isEqual(oldState, newState)) {
            this._state = newState;
            this._history.push(newState);
            this._notifyListeners(oldState, newState);
        }

        return this;
    }

    // Batch multiple updates
    batch(updates) {
        let newState = this._state;
        
        for (const update of updates) {
            newState = typeof update === 'function' 
                ? update(newState) 
                : { ...newState, ...update };
        }

        return this.update(newState);
    }

    // Undo functionality
    undo() {
        if (this._history.length > 1) {
            this._history.pop(); // Remove current state
            this._state = this._history[this._history.length - 1];
            this._notifyListeners(this._history[this._history.length], this._state);
        }
        return this;
    }

    // Subscribe to changes
    subscribe(listener) {
        this._listeners.add(listener);
        return () => this._listeners.delete(listener);
    }

    _notifyListeners(oldState, newState) {
        this._listeners.forEach(listener => listener(newState, oldState));
    }

    isEqual(obj1, obj2) {
        return JSON.stringify(obj1) === JSON.stringify(obj2);
    }

    // Create a selector for derived state
    createSelector(selectorFn) {
        let lastResult = null;
        let lastState = null;

        return (state = this._state) => {
            if (state === lastState) {
                return lastResult;
            }
            lastState = state;
            lastResult = selectorFn(state);
            return lastResult;
        };
    }
}

// Usage example
console.log("=== IMMUTABLE STORE ===");
const store = new ImmutableStore({
    user: { name: "Alice", age: 25 },
    todos: [],
    filter: "all"
});

// Subscribe to changes
const unsubscribe = store.subscribe((newState, oldState) => {
    console.log("State changed!");
    console.log("Old:", oldState.user);
    console.log("New:", newState.user);
});

// Update state immutably
store.update(state => ({
    ...state,
    user: { ...state.user, age: 26 }
}));

// Add a todo
store.update(state => ({
    ...state,
    todos: [...state.todos, { id: 1, text: "Learn immutability", completed: false }]
}));

// Batch updates
store.batch([
    state => ({ ...state, filter: "active" }),
    state => ({ 
        ...state, 
        todos: state.todos.map(todo => 
            todo.id === 1 ? { ...todo, completed: true } : todo
        )
    })
]);

// Create a selector
const getUserName = store.createSelector(state => state.user.name);
console.log("User name:", getUserName());

// Undo last change
store.undo();
console.log("After undo:", store.state);

unsubscribe();
```

### Example 4: Immutable Data Structures with Performance

```javascript
// For large datasets, use structural sharing (like Immutable.js)
class ImmutableMap {
    constructor(entries = {}) {
        this._data = { ...entries };
        this._version = 0;
    }

    set(key, value) {
        if (this._data[key] === value) {
            return this; // No change needed
        }

        const newMap = Object.create(ImmutableMap.prototype);
        newMap._data = { ...this._data, [key]: value };
        newMap._version = this._version + 1;
        return newMap;
    }

    get(key) {
        return this._data[key];
    }

    delete(key) {
        if (!(key in this._data)) {
            return this; // No change needed
        }

        const { [key]: deleted, ...rest } = this._data;
        const newMap = Object.create(ImmutableMap.prototype);
        newMap._data = rest;
        newMap._version = this._version + 1;
        return newMap;
    }

    has(key) {
        return key in this._data;
    }

    map(fn) {
        const newData = {};
        for (const [key, value] of Object.entries(this._data)) {
            newData[key] = fn(value, key);
        }
        
        const newMap = Object.create(ImmutableMap.prototype);
        newMap._data = newData;
        newMap._version = this._version + 1;
        return newMap;
    }

    filter(fn) {
        const newData = {};
        for (const [key, value] of Object.entries(this._data)) {
            if (fn(value, key)) {
                newData[key] = value;
            }
        }
        
        const newMap = Object.create(ImmutableMap.prototype);
        newMap._data = newData;
        newMap._version = this._version + 1;
        return newMap;
    }

    // Structural sharing - only changed parts are new
    merge(otherMap) {
        let changed = false;
        const newData = { ...this._data };

        for (const [key, value] of Object.entries(otherMap._data)) {
            if (newData[key] !== value) {
                newData[key] = value;
                changed = true;
            }
        }

        if (!changed) {
            return this;
        }

        const newMap = Object.create(ImmutableMap.prototype);
        newMap._data = newData;
        newMap._version = this._version + 1;
        return newMap;
    }

    toJS() {
        return { ...this._data };
    }

    equals(other) {
        if (this._version === other._version && this === other) {
            return true;
        }

        const keys1 = Object.keys(this._data);
        const keys2 = Object.keys(other._data);

        if (keys1.length !== keys2.length) {
            return false;
        }

        for (const key of keys1) {
            if (this._data[key] !== other._data[key]) {
                return false;
            }
        }

        return true;
    }
}

// Usage
console.log("=== IMMUTABLE MAP ===");
const map1 = new ImmutableMap({ a: 1, b: 2, c: 3 });
console.log("Map1:", map1.toJS());

const map2 = map1.set('b', 20);
console.log("Map2 (after set):", map2.toJS());
console.log("Map1 unchanged:", map1.toJS());
console.log("Same object?", map1 === map2); // false

const map3 = map2.delete('c');
console.log("Map3 (after delete):", map3.toJS());

const map4 = map3.map(value => value * 2);
console.log("Map4 (after map):", map4.toJS());

const map5 = map4.filter(value => value > 10);
console.log("Map5 (after filter):", map5.toJS());
```

### Example 5: React-like State Management with Immutability

```javascript
// Simulating React's useState with immutability
class ImmutableReactPattern {
    static useState(initialValue) {
        let state = typeof initialValue === 'function' ? initialValue() : initialValue;
        const listeners = new Set();

        const setState = (newValue) => {
            const nextState = typeof newValue === 'function' ? newValue(state) : newValue;

            // Only update if state changed
            if (!this.isEqual(state, nextState)) {
                const previousState = state;
                state = nextState;
                listeners.forEach(listener => listener(state, previousState));
            }
        };

        const getState = () => state;
        const subscribe = (listener) => {
            listeners.add(listener);
            return () => listeners.delete(listener);
        };

        return [getState, setState, subscribe];
    }

    // Immutable reducer pattern (like Redux)
    static createReducer(initialState, handlers) {
        return (state = initialState, action) => {
            const handler = handlers[action.type];
            
            if (handler) {
                const newState = handler(state, action);
                
                // Ensure immutability
                if (newState === state) {
                    console.warn('Reducer returned same state for action:', action.type);
                    return state;
                }
                
                return newState;
            }
            
            return state;
        };
    }

    // Combine reducers immutably
    static combineReducers(reducers) {
        return (state = {}, action) => {
            let hasChanged = false;
            const nextState = {};

            for (const [key, reducer] of Object.entries(reducers)) {
                const previousStateForKey = state[key];
                const nextStateForKey = reducer(previousStateForKey, action);
                
                nextState[key] = nextStateForKey;
                hasChanged = hasChanged || nextStateForKey !== previousStateForKey;
            }

            return hasChanged ? nextState : state;
        };
    }

    static isEqual(a, b) {
        if (a === b) return true;
        if (typeof a !== 'object' || a === null || typeof b !== 'object' || b === null) {
            return false;
        }

        const keysA = Object.keys(a);
        const keysB = Object.keys(b);

        if (keysA.length !== keysB.length) return false;

        for (const key of keysA) {
            if (!keysB.includes(key) || !this.isEqual(a[key], b[key])) {
                return false;
            }
        }

        return true;
    }
}

// Usage example
console.log("=== REACT-LIKED PATTERNS ===");

// Reducer example
const todoReducer = ImmutableReactPattern.createReducer([], {
    ADD_TODO: (state, action) => [
        ...state,
        {
            id: Date.now(),
            text: action.payload,
            completed: false
        }
    ],
    
    TOGGLE_TODO: (state, action) =>
        state.map(todo =>
            todo.id === action.payload
                ? { ...todo, completed: !todo.completed }
                : todo
        ),
    
    DELETE_TODO: (state, action) =>
        state.filter(todo => todo.id !== action.payload)
});

// Test the reducer
let state = todoReducer(undefined, { type: '@@INIT' });
console.log("Initial state:", state);

state = todoReducer(state, { type: 'ADD_TODO', payload: 'Learn immutability' });
console.log("After ADD_TODO:", state);

state = todoReducer(state, { type: 'TOGGLE_TODO', payload: state[0].id });
console.log("After TOGGLE_TODO:", state);

// useState example
const [getCount, setCount, subscribe] = ImmutableReactPattern.useState(0);

subscribe((newCount, oldCount) => {
    console.log(`Count changed from ${oldCount} to ${newCount}`);
});

setCount(1);
setCount(prev => prev + 1);
```

## Key Benefits Demonstrated

1. **Predictable Updates**: Every change creates a new object
2. **Easy Change Detection**: `oldState === newState` checks work
3. **Time Travel Debugging**: History is preserved naturally
4. **Performance**: No deep copies unless necessary
5. **Safety**: No accidental mutations

## Common Patterns Summary

| Pattern | Use Case | Example |
|---------|----------|---------|
| **Spread Operator** | Shallow updates | `{ ...obj, update }` |
| **Array Methods** | List operations | `map`, `filter`, `slice` |
| **Structural Sharing** | Performance | Only copy changed parts |
| **Reducer Pattern** | State machines | Redux-style updates |
| **Lens Pattern** | Deep updates | Focus on nested properties |

Immutable patterns are essential for modern JavaScript development, especially in frameworks like React where they enable optimal rendering performance through simple reference checks.