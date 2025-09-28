
// 1. Define the possible STATES (as constants)
const STATES = {
  GREEN: 'GREEN',
  YELLOW: 'YELLOW',
  RED: 'RED'
};

// 2. Define the possible EVENTS (as constants)
const EVENTS = {
  TIMER_EXPIRED: 'TIMER_EXPIRED'
};

const TrafficTime = 10000 // every 10s later automatic light change

class TrafficLightFSM {
  constructor() {
    // 3. Set the initial state
    this.currentState = STATES.GREEN;
    this.startTimer(); // Start the timer for the initial state
  }

  // 4. The core transition function
  transition(event) {
    console.log(`Received event: ${event}. Current state: ${this.currentState}`);

    // 5. Define the transition logic
    switch (this.currentState) {
      case STATES.GREEN:
        if (event === EVENTS.TIMER_EXPIRED) {
          this.setState(STATES.YELLOW);
        }
        break;
      case STATES.YELLOW:
        if (event === EVENTS.TIMER_EXPIRED) {
          this.setState(STATES.RED);
        }
        break;
      case STATES.RED:
        if (event === EVENTS.TIMER_EXPIRED) {
          this.setState(STATES.GREEN);
        }
        break;
      default:
        // This should never happen if states are defined correctly!
        throw new Error(`Invalid state: ${this.currentState}`);
    }
  }

  // Helper method to change state and perform actions
  setState(newState) {
    console.log(`Transitioning from ${this.currentState} to ${newState}`);
    this.currentState = newState;

    // 6. Perform actions associated with entering the new state
    // This could be turning on/off lights, starting new timers, etc.
    this.startTimer();
    this.updateLights();
  }

  updateLights() {
    console.log(`>>> ACTIVATING THE ${this.currentState} LIGHT <<<`);
  }

  startTimer() {
    // Simulate a timer. In reality, this would be a setTimeout.
    console.log(`Starting timer for ${this.currentState} state...`);
    // After a simulated delay, trigger the TIMER_EXPIRED event.
    setTimeout(() => {
      this.transition(EVENTS.TIMER_EXPIRED);
    }, TrafficTime); // 3 seconds for simplicity
  }
}

// Usage
const light = new TrafficLightFSM();

// console.log(light())
// The light will now run automatically based on timers.
// Output:
// Starting timer for GREEN state...
// >>> ACTIVATING THE GREEN LIGHT <<<
// ... (after 3 seconds)
// Received event: TIMER_EXPIRED. Current state: GREEN
// Transitioning from GREEN to YELLOW
// >>> ACTIVATING THE YELLOW LIGHT <<<
// Starting timer for YELLOW state...

