# Token Streaming Fundamentals in JavaScript

Token streaming is a technique where data is processed and transmitted in small chunks (tokens) rather than waiting for the complete payload. This is essential for large language models, real-time data processing, and responsive user interfaces.

## Core Concepts

### **What are Tokens?**
- **Chunks of data** - pieces of text, JSON objects, binary data
- **Sequential delivery** - tokens arrive in order
- **Incremental processing** - each token can be processed immediately

### **Why Stream Tokens?**
- **Low latency** - First token arrives quickly
- **Memory efficiency** - Don't need to hold entire response in memory
- **User experience** - Progressive rendering/display
- **Network efficiency** - Can start processing before download completes

---

## JavaScript Examples

### Example 1: Basic Text Token Streamer

```javascript
class TextTokenStreamer {
    constructor(options = {}) {
        this.chunkSize = options.chunkSize || 10; // characters per token
        this.delay = options.delay || 50; // ms between tokens
        this.onToken = options.onToken || (() => {});
        this.onComplete = options.onComplete || (() => {});
        this.onError = options.onError || (() => {});
        
        this.isStreaming = false;
        this.currentPosition = 0;
        this.text = '';
    }

    // Stream text token by token
    async stream(text) {
        if (this.isStreaming) {
            throw new Error('Streaming already in progress');
        }

        this.isStreaming = true;
        this.text = text;
        this.currentPosition = 0;

        try {
            while (this.currentPosition < this.text.length && this.isStreaming) {
                const chunk = this.text.substring(
                    this.currentPosition,
                    this.currentPosition + this.chunkSize
                );
                
                this.currentPosition += chunk.length;
                
                // Emit token
                this.onToken(chunk, {
                    position: this.currentPosition,
                    total: this.text.length,
                    progress: (this.currentPosition / this.text.length) * 100
                });

                // Simulate network delay or processing time
                await this.sleep(this.delay);
            }

            if (this.isStreaming) {
                this.onComplete();
            }
        } catch (error) {
            this.onError(error);
        } finally {
            this.isStreaming = false;
        }
    }

    // Pause the stream
    pause() {
        this.isStreaming = false;
    }

    // Resume from current position
    resume() {
        if (!this.isStreaming && this.currentPosition < this.text.length) {
            this.stream(this.text.substring(this.currentPosition));
        }
    }

    // Stop streaming completely
    stop() {
        this.isStreaming = false;
        this.currentPosition = 0;
    }

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }

    getProgress() {
        return {
            position: this.currentPosition,
            total: this.text.length,
            progress: (this.currentPosition / this.text.length) * 100,
            isStreaming: this.isStreaming
        };
    }
}

// Usage example
function demoTextStreaming() {
    const streamer = new TextTokenStreamer({
        chunkSize: 5,
        delay: 100,
        onToken: (token, metadata) => {
            process.stdout.write(token); // Typewriter effect
            console.log(` | Progress: ${metadata.progress.toFixed(1)}%`);
        },
        onComplete: () => {
            console.log('\n✅ Streaming completed!');
        },
        onError: (error) => {
            console.error('❌ Streaming error:', error);
        }
    });

    const longText = "Token streaming is a powerful technique for processing large amounts of data incrementally. Instead of waiting for the entire payload, we process small chunks as they arrive.";
    
    streamer.stream(longText);
}

// demoTextStreaming();
```

### Example 2: AI Response Streaming (ChatGPT-like)

```javascript
class AIResponseStreamer {
    constructor() {
        this.subscribers = new Set();
        this.isStreaming = false;
        this.buffer = [];
    }

    // Simulate AI response generation
    async simulateAIResponse(prompt, options = {}) {
        if (this.isStreaming) {
            throw new Error('Already streaming');
        }

        this.isStreaming = true;
        this.buffer = [];

        // Simulate different AI responses based on prompt
        const responses = {
            greeting: "Hello! I'm an AI assistant. How can I help you today?",
            weather: "The weather today is sunny with a high of 75°F. Perfect day for outdoor activities!",
            joke: "Why don't scientists trust atoms? Because they make up everything!",
            default: "I understand your question. Let me provide a detailed response about this topic."
        };

        const responseType = prompt.toLowerCase().includes('hello') ? 'greeting' :
                           prompt.toLowerCase().includes('weather') ? 'weather' :
                           prompt.toLowerCase().includes('joke') ? 'joke' : 'default';

        const fullResponse = responses[responseType];
        const words = fullResponse.split(' ');
        
        try {
            // Stream words one by one with realistic timing
            for (let i = 0; i < words.length; i++) {
                if (!this.isStreaming) break; // Check if streaming was stopped
                
                const word = words[i];
                this.buffer.push(word);
                
                // Emit the current state
                this.emit('token', {
                    token: word,
                    fullText: this.buffer.join(' '),
                    isComplete: false,
                    tokensSoFar: i + 1,
                    totalTokens: words.length
                });

                // Simulate AI "thinking" time - longer for complex responses
                const delay = responseType === 'default' ? 120 : 80;
                await this.sleep(delay + Math.random() * 40);
            }

            if (this.isStreaming) {
                this.emit('complete', {
                    fullText: this.buffer.join(' '),
                    totalTokens: words.length,
                    responseType: responseType
                });
            }
        } catch (error) {
            this.emit('error', error);
        } finally {
            this.isStreaming = false;
        }
    }

    // Subscribe to stream events
    subscribe(callback) {
        this.subscribers.add(callback);
        return () => this.subscribers.delete(callback);
    }

    // Emit events to all subscribers
    emit(event, data) {
        this.subscribers.forEach(callback => {
            try {
                callback(event, data);
            } catch (error) {
                console.error('Subscriber error:', error);
            }
        });
    }

    // Stop the stream
    stop() {
        this.isStreaming = false;
        this.emit('stopped', { position: this.buffer.length });
    }

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Usage example
function demoAIStreaming() {
    const ai = new AIResponseStreamer();
    
    // Subscribe to stream events
    const unsubscribe = ai.subscribe((event, data) => {
        switch (event) {
            case 'token':
                process.stdout.write(data.token + ' ');
                break;
            case 'complete':
                console.log('\n\n✅ Response complete!');
                console.log('Full response:', data.fullText);
                break;
            case 'error':
                console.error('AI error:', data);
                break;
            case 'stopped':
                console.log('\n\n⏹️ Streaming stopped');
                break;
        }
    });

    // Simulate user interaction
    console.log('User: Tell me a joke');
    ai.simulateAIResponse("Tell me a joke")
        .then(() => {
            console.log('\n--- Conversation complete ---');
            unsubscribe();
        })
        .catch(error => console.error('Error:', error));

    // You could stop the stream after 2 seconds to see interruption
    // setTimeout(() => ai.stop(), 2000);
}

// demoAIStreaming();
```

### Example 3: HTTP Response Streaming with Fetch API

```javascript
class HTTPStreamProcessor {
    constructor() {
        this.readers = new Map();
        this.controller = null;
    }

    // Stream HTTP response with proper chunk handling
    async streamHTTPResponse(url, onChunk, onComplete, options = {}) {
        try {
            const response = await fetch(url, {
                method: options.method || 'GET',
                headers: options.headers,
                signal: options.signal
            });

            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }

            if (!response.body) {
                throw new Error('Response body is not available for streaming');
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let accumulatedText = '';
            let chunkCount = 0;

            const streamId = Math.random().toString(36).substr(2, 9);
            this.readers.set(streamId, reader);

            try {
                while (true) {
                    const { done, value } = await reader.read();
                    
                    if (done) {
                        onComplete?.(accumulatedText, chunkCount);
                        this.readers.delete(streamId);
                        break;
                    }

                    chunkCount++;
                    const chunk = decoder.decode(value, { stream: true });
                    accumulatedText += chunk;

                    // Process chunk (could be JSON, text, etc.)
                    onChunk(chunk, {
                        chunkNumber: chunkCount,
                        accumulatedLength: accumulatedText.length,
                        isComplete: false
                    });
                }
            } catch (error) {
                this.readers.delete(streamId);
                throw error;
            }

        } catch (error) {
            throw new Error(`Streaming failed: ${error.message}`);
        }
    }

    // Stream JSON lines (common in AI APIs)
    async streamJSONLines(url, onObject, onComplete) {
        let buffer = '';
        
        return this.streamHTTPResponse(
            url,
            (chunk) => {
                buffer += chunk;
                const lines = buffer.split('\n');
                
                // Keep last incomplete line in buffer
                buffer = lines.pop() || '';
                
                // Process complete lines
                lines.forEach(line => {
                    if (line.trim()) {
                        try {
                            const obj = JSON.parse(line);
                            onObject(obj);
                        } catch (error) {
                            console.warn('Failed to parse JSON line:', line);
                        }
                    }
                });
            },
            (finalText) => {
                // Process any remaining buffer
                if (buffer.trim()) {
                    try {
                        const obj = JSON.parse(buffer);
                        onObject(obj);
                    } catch (error) {
                        console.warn('Failed to parse final JSON:', buffer);
                    }
                }
                onComplete?.(finalText);
            }
        );
    }

    // Cancel a specific stream
    cancelStream(streamId) {
        const reader = this.readers.get(streamId);
        if (reader) {
            reader.cancel();
            this.readers.delete(streamId);
        }
    }

    // Cancel all streams
    cancelAll() {
        this.readers.forEach(reader => reader.cancel());
        this.readers.clear();
    }
}

// Mock server for demonstration
class MockStreamingServer {
    constructor() {
        this.messages = [
            "This is the first token of the response.",
            " Here's the second part of the message.",
            " The streaming continues with more content.",
            " Finally, we reach the end of the response."
        ];
    }

    async createMockResponse() {
        const stream = new ReadableStream({
            start(controller) {
                // Simulate delayed token delivery
                let index = 0;
                
                function pushNextToken() {
                    if (index < this.messages.length) {
                        const encoder = new TextEncoder();
                        const chunk = encoder.encode(this.messages[index]);
                        controller.enqueue(chunk);
                        index++;
                        
                        // Simulate network delay
                        setTimeout(pushNextToken.bind(this), 100 + Math.random() * 100);
                    } else {
                        controller.close();
                    }
                }
                
                setTimeout(pushNextToken.bind(this), 200);
            }.bind(this)
        });

        return new Response(stream, {
            headers: { 'Content-Type': 'text/plain' }
        });
    }
}
```

### Example 4: Advanced Token Stream with Backpressure

```javascript
class AdvancedTokenStream {
    constructor(options = {}) {
        this.bufferSize = options.bufferSize || 100;
        this.processingRate = options.processingRate || 10; // tokens/sec
        this.subscribers = new Set();
        this.tokenBuffer = [];
        this.isProducing = false;
        this.isConsuming = false;
        this.producedCount = 0;
        this.consumedCount = 0;
        
        this.stats = {
            tokensProduced: 0,
            tokensConsumed: 0,
            bufferHighWaterMark: 0,
            backpressureEvents: 0,
            startTime: Date.now()
        };
    }

    // Producer: Generate tokens
    async startProducing(tokenGenerator, totalTokens) {
        if (this.isProducing) {
            throw new Error('Already producing tokens');
        }

        this.isProducing = true;
        this.producedCount = 0;

        for (let i = 0; i < totalTokens; i++) {
            if (!this.isProducing) break;

            // Apply backpressure if buffer is too full
            while (this.tokenBuffer.length >= this.bufferSize) {
                this.stats.backpressureEvents++;
                console.log('⏸️ Backpressure applied - buffer full');
                await this.sleep(100);
            }

            const token = await tokenGenerator(i);
            this.tokenBuffer.push(token);
            this.producedCount++;
            this.stats.tokensProduced++;

            this.stats.bufferHighWaterMark = Math.max(
                this.stats.bufferHighWaterMark, 
                this.tokenBuffer.length
            );

            this.emit('tokenProduced', {
                token,
                position: i,
                total: totalTokens,
                bufferSize: this.tokenBuffer.length
            });

            // Rate limiting
            await this.sleep(1000 / this.processingRate);
        }

        this.isProducing = false;
        this.emit('productionComplete', { totalProduced: this.producedCount });
    }

    // Consumer: Process tokens
    async startConsuming(processorFn) {
        if (this.isConsuming) {
            throw new Error('Already consuming tokens');
        }

        this.isConsuming = true;
        this.consumedCount = 0;

        while (this.isConsuming) {
            if (this.tokenBuffer.length > 0) {
                const token = this.tokenBuffer.shift();
                this.consumedCount++;
                this.stats.tokensConsumed++;

                try {
                    await processorFn(token, this.consumedCount);
                } catch (error) {
                    this.emit('processingError', { error, token });
                }

                this.emit('tokenConsumed', {
                    token,
                    count: this.consumedCount,
                    bufferSize: this.tokenBuffer.length
                });
            } else if (!this.isProducing) {
                // No more tokens and production is complete
                break;
            } else {
                // Wait for more tokens
                await this.sleep(50);
            }
        }

        this.isConsuming = false;
        this.emit('consumptionComplete', { totalConsumed: this.consumedCount });
    }

    // Subscribe to stream events
    on(event, callback) {
        this.subscribers.add({ event, callback });
        return () => {
            this.subscribers.forEach(sub => {
                if (sub.event === event && sub.callback === callback) {
                    this.subscribers.delete(sub);
                }
            });
        };
    }

    emit(event, data) {
        this.subscribers.forEach(({ event: subEvent, callback }) => {
            if (subEvent === event) {
                try {
                    callback(data);
                } catch (error) {
                    console.error('Event handler error:', error);
                }
            }
        });
    }

    // Get current statistics
    getStats() {
        const runningTime = Date.now() - this.stats.startTime;
        return {
            ...this.stats,
            runningTime,
            tokensPerSecond: this.stats.tokensConsumed / (runningTime / 1000),
            bufferUtilization: this.tokenBuffer.length / this.bufferSize,
            productionActive: this.isProducing,
            consumptionActive: this.isConsuming
        };
    }

    // Stop both production and consumption
    stop() {
        this.isProducing = false;
        this.isConsuming = false;
        this.emit('streamStopped', this.getStats());
    }

    sleep(ms) {
        return new Promise(resolve => setTimeout(resolve, ms));
    }
}

// Usage example
async function demoAdvancedStreaming() {
    const stream = new AdvancedTokenStream({
        bufferSize: 20,
        processingRate: 5 // 5 tokens/second
    });

    // Set up monitoring
    stream.on('tokenProduced', (data) => {
        console.log(`📤 Produced: "${data.token}" | Buffer: ${data.bufferSize}/${stream.bufferSize}`);
    });

    stream.on('tokenConsumed', (data) => {
        console.log(`📥 Consumed: "${data.token}" | Total: ${data.count} | Buffer: ${data.bufferSize}`);
    });

    stream.on('productionComplete', (data) => {
        console.log(`✅ Production complete: ${data.totalProduced} tokens`);
    });

    stream.on('consumptionComplete', (data) => {
        console.log(`✅ Consumption complete: ${data.totalConsumed} tokens`);
    });

    // Token generator
    const tokenGenerator = async (index) => {
        const tokens = [
            "Hello", "World", "Streaming", "Tokens", "JavaScript",
            "Async", "Await", "Promises", "Events", "Backpressure"
        ];
        return `${tokens[index % tokens.length]}_${index}`;
    };

    // Token processor
    const tokenProcessor = async (token, count) => {
        // Simulate processing work
        await stream.sleep(200);
        return `Processed: ${token}`;
    };

    console.log('🚀 Starting advanced token streaming...\n');

    // Start production and consumption
    await Promise.all([
        stream.startProducing(tokenGenerator, 50),
        stream.startConsuming(tokenProcessor)
    ]);

    console.log('\n📊 Final Stats:', stream.getStats());
}
```

### Example 5: Real-world AI API Streaming (OpenAI-like)

```javascript
class OpenAIStyleStreamer {
    constructor(apiKey, options = {}) {
        this.apiKey = apiKey;
        this.baseURL = options.baseURL || 'https://api.openai.com/v1';
        this.defaultModel = options.defaultModel || 'gpt-3.5-turbo';
        this.activeStreams = new Map();
    }

    async streamChatCompletion(messages, options = {}) {
        const streamId = Math.random().toString(36).substr(2, 9);
        const abortController = new AbortController();
        
        this.activeStreams.set(streamId, { abortController });

        try {
            const response = await fetch(`${this.baseURL}/chat/completions`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${this.apiKey}`
                },
                body: JSON.stringify({
                    model: options.model || this.defaultModel,
                    messages: messages,
                    stream: true,
                    max_tokens: options.max_tokens,
                    temperature: options.temperature
                }),
                signal: abortController.signal
            });

            if (!response.ok) {
                throw new Error(`API error: ${response.status} ${response.statusText}`);
            }

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';
            let fullResponse = '';

            while (true) {
                const { done, value } = await reader.read();
                
                if (done) {
                    this.activeStreams.delete(streamId);
                    return { fullResponse, streamId, completed: true };
                }

                buffer += decoder.decode(value, { stream: true });
                const lines = buffer.split('\n');
                buffer = lines.pop() || ''; // Keep incomplete line

                for (const line of lines) {
                    const trimmedLine = line.trim();
                    
                    if (trimmedLine === '') continue;
                    if (trimmedLine === 'data: [DONE]') {
                        this.activeStreams.delete(streamId);
                        return { fullResponse, streamId, completed: true };
                    }

                    if (trimmedLine.startsWith('data: ')) {
                        try {
                            const data = JSON.parse(trimmedLine.slice(6));
                            const content = data.choices[0]?.delta?.content;
                            
                            if (content) {
                                fullResponse += content;
                                
                                // Emit token event
                                this.emit('token', {
                                    streamId,
                                    token: content,
                                    fullResponse,
                                    choice: data.choices[0]
                                });
                            }

                            if (data.choices[0]?.finish_reason) {
                                this.emit('complete', {
                                    streamId,
                                    fullResponse,
                                    finishReason: data.choices[0].finish_reason
                                });
                            }
                        } catch (error) {
                            console.warn('Failed to parse stream data:', trimmedLine);
                        }
                    }
                }
            }

        } catch (error) {
            this.activeStreams.delete(streamId);
            if (error.name !== 'AbortError') {
                this.emit('error', { streamId, error });
                throw error;
            }
        }
    }

    // Cancel a specific stream
    cancelStream(streamId) {
        const stream = this.activeStreams.get(streamId);
        if (stream) {
            stream.abortController.abort();
            this.activeStreams.delete(streamId);
            this.emit('cancelled', { streamId });
        }
    }

    // Event emitter methods (simplified)
    emit(event, data) {
        // In real implementation, this would notify subscribers
        console.log(`Event: ${event}`, data);
    }
}

// Usage example for AI streaming
async function demoAIChatStreaming() {
    // This is a mock example - real implementation would use actual API
    const streamer = new OpenAIStyleStreamer('mock-api-key');
    
    const messages = [
        { role: 'user', content: 'Explain token streaming in simple terms.' }
    ];

    console.log('🤖 AI Response Streaming...\n');
    
    let accumulatedResponse = '';
    
    // Mock event handler
    streamer.emit = (event, data) => {
        switch (event) {
            case 'token':
                process.stdout.write(data.token); // Typewriter effect
                accumulatedResponse += data.token;
                break;
            case 'complete':
                console.log('\n\n✅ Response complete!');
                console.log('Full response length:', accumulatedResponse.length);
                break;
            case 'error':
                console.error('Stream error:', data.error);
                break;
        }
    };

    // Simulate streaming (in real use, this would be actual API call)
    await simulateAIStream(streamer, messages);
}

async function simulateAIStream(streamer, messages) {
    const responseText = "Token streaming is a technique where data is sent in small pieces instead of all at once. This allows the recipient to start processing the data immediately as it arrives, rather than waiting for the entire response to download.";
    const words = responseText.split(' ');
    
    for (let i = 0; i < words.length; i++) {
        await new Promise(resolve => setTimeout(resolve, 100));
        streamer.emit('token', {
            token: words[i] + (i < words.length - 1 ? ' ' : ''),
            fullResponse: words.slice(0, i + 1).join(' ')
        });
    }
    
    streamer.emit('complete', { fullResponse: responseText });
}
```

## Key Token Streaming Concepts

### **1. Chunking Strategies**
- **Fixed-size chunks**: Consistent token sizes
- **Variable-size chunks**: Natural boundaries (words, sentences)
- **JSON streaming**: One object per line
- **Server-Sent Events**: Standardized streaming protocol

### **2. Flow Control**
- **Backpressure**: Preventing buffer overflow
- **Rate limiting**: Controlling token production speed
- **Cancellation**: Stopping streams gracefully

### **3. Error Handling**
- **Resilience**: Continuing after token errors
- **Retry logic**: Recovering from failures
- **Completion detection**: Knowing when stream ends

### **4. Performance Considerations**
- **Memory usage**: Streaming vs buffering
- **Network efficiency**: TCP vs HTTP/2 streaming
- **Processing overhead**: Tokenization costs

## Benefits of Token Streaming

1. **Improved Latency**: Users see first results quickly
2. **Memory Efficiency**: Process large datasets without loading everything
3. **Better UX**: Progressive loading and typewriter effects
4. **Network Resilience**: Handle intermittent connectivity
5. **Real-time Processing**: Live data processing capabilities

Token streaming is fundamental to modern web applications, especially for AI interfaces, real-time dashboards, and large data processing pipelines.