# Server-Sent Events (SSE) - Real-Time Server-to-Client Communication

SSE provides a **one-way real-time communication channel** from server to client over a single HTTP connection.

## Key SSE Concepts

- **Unidirectional**: Server can push data to client, but client cannot send data back through SSE
- **HTTP-based**: Uses standard HTTP protocol, no special protocol like WebSockets
- **Automatic Reconnection**: Built-in reconnection mechanism
- **Text-based**: Primarily designed for text data (can send JSON)
- **Simple API**: Easy to implement on both client and server

## Basic SSE Example

### Server (Node.js)

```javascript
const http = require("http");
const fs = require("fs");

const server = http.createServer((req, res) => {
  // Serve HTML page
  if (req.url === "/") {
    fs.readFile("./index.html", (err, data) => {
      if (err) {
        res.writeHead(404);
        res.end("Not found");
        return;
      }
      res.writeHead(200, { "Content-Type": "text/html" });
      res.end(data);
    });
    return;
  }

  // SSE endpoint
  if (req.url === "/events") {
    console.log("SSE connection established");

    // Set SSE headers
    res.writeHead(200, {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
      "Access-Control-Allow-Origin": "*",
    });

    // Send initial event
    res.write("data: Connected to SSE server\n\n");

    // Send periodic updates
    let counter = 0;
    const interval = setInterval(() => {
      counter++;
      res.write(`data: Server time: ${new Date().toISOString()}\n`);
      res.write(`id: ${counter}\n`);
      res.write(`event: update\n`);
      res.write(`retry: 5000\n\n`); // Reconnect after 5 seconds if connection lost
    }, 2000);

    // Handle client disconnect
    req.on("close", () => {
      console.log("SSE connection closed");
      clearInterval(interval);
      res.end();
    });

    return;
  }

  res.writeHead(404);
  res.end();
});

server.listen(3000, () => {
  console.log("SSE server running on http://localhost:3000");
});
```

### Client (Browser JavaScript)

```javascript
class SSEClient {
  constructor() {
    this.eventSource = null;
    this.isConnected = false;
  }

  connect() {
    try {
      this.eventSource = new EventSource("/events");

      this.eventSource.onopen = (event) => {
        this.isConnected = true;
        this.displayMessage("System: Connection established", "system");
        console.log("SSE connection opened");
      };

      this.eventSource.onmessage = (event) => {
        this.displayMessage(`Data: ${event.data}`, "message");
      };

      this.eventSource.addEventListener("update", (event) => {
        this.displayMessage(
          `Update: ${event.data} (ID: ${event.lastEventId})`,
          "update"
        );
      });

      this.eventSource.onerror = (event) => {
        this.isConnected = false;
        this.displayMessage("System: Connection error", "error");
        console.error("SSE error:", event);
      };
    } catch (error) {
      console.error("Failed to create EventSource:", error);
    }
  }

  disconnect() {
    if (this.eventSource) {
      this.eventSource.close();
      this.isConnected = false;
      this.displayMessage("System: Connection closed", "system");
    }
  }

  displayMessage(message, type) {
    const messagesDiv = document.getElementById("messages");
    const messageElement = document.createElement("div");
    messageElement.className = `message ${type}`;
    messageElement.textContent = message;
    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }
}

// Initialize SSE client
const sseClient = new SSEClient();

// Connect when page loads
window.addEventListener("load", () => {
  sseClient.connect();
});

// Clean up when leaving page
window.addEventListener("beforeunload", () => {
  sseClient.disconnect();
});
```

### HTML Client Page

```html
<!DOCTYPE html>
<html>
  <head>
    <title>SSE Client</title>
    <style>
      body {
        font-family: Arial, sans-serif;
        margin: 20px;
      }
      .container {
        max-width: 800px;
        margin: 0 auto;
      }
      #messages {
        border: 1px solid #ccc;
        height: 400px;
        overflow-y: auto;
        padding: 10px;
        margin-bottom: 20px;
      }
      .message {
        margin: 5px 0;
        padding: 5px;
      }
      .system {
        color: #666;
        font-style: italic;
      }
      .update {
        color: #2c5aa0;
        background-color: #f0f8ff;
      }
      .error {
        color: #d00;
        background-color: #ffe6e6;
      }
      .controls {
        margin: 20px 0;
      }
      button {
        padding: 10px 20px;
        margin-right: 10px;
      }
      .status {
        padding: 10px;
        margin: 10px 0;
      }
      .connected {
        background-color: #d4edda;
        color: #155724;
      }
      .disconnected {
        background-color: #f8d7da;
        color: #721c24;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <h1>Server-Sent Events Demo</h1>

      <div id="status" class="status connected">Status: Connected</div>

      <div class="controls">
        <button onclick="connectSSE()">Connect</button>
        <button onclick="disconnectSSE()">Disconnect</button>
        <button onclick="sendToServer()">Send Message to Server</button>
      </div>

      <div id="messages"></div>

      <div>
        <h3>Send data to server (via fetch):</h3>
        <input
          type="text"
          id="messageInput"
          placeholder="Type message for server"
        />
        <button onclick="sendToServer()">Send</button>
      </div>
    </div>

    <script>
      const sseClient = new SSEClient();

      function connectSSE() {
        sseClient.connect();
        updateStatus(true);
      }

      function disconnectSSE() {
        sseClient.disconnect();
        updateStatus(false);
      }

      function updateStatus(connected) {
        const statusDiv = document.getElementById("status");
        statusDiv.textContent = `Status: ${
          connected ? "Connected" : "Disconnected"
        }`;
        statusDiv.className = `status ${
          connected ? "connected" : "disconnected"
        }`;
      }

      async function sendToServer() {
        const input = document.getElementById("messageInput");
        const message = input.value.trim();

        if (message) {
          try {
            // Send data to server via regular HTTP request
            const response = await fetch("/message", {
              method: "POST",
              headers: { "Content-Type": "application/json" },
              body: JSON.stringify({ message: message }),
            });

            if (response.ok) {
              sseClient.displayMessage(`You: ${message}`, "message");
              input.value = "";
            }
          } catch (error) {
            console.error("Error sending message:", error);
          }
        }
      }

      // Auto-connect on page load
      window.addEventListener("load", connectSSE);
    </script>
  </body>
</html>
```

## Advanced SSE Server with Multiple Event Types

```javascript
const http = require("http");
const fs = require("fs");
const url = require("url");

class SSEServer {
  constructor(port) {
    this.port = port;
    this.clients = new Set();
    this.setupServer();
  }

  setupServer() {
    this.server = http.createServer((req, res) => {
      const parsedUrl = url.parse(req.url, true);

      switch (parsedUrl.pathname) {
        case "/":
          this.serveHTML(req, res);
          break;
        case "/events":
          this.handleSSE(req, res);
          break;
        case "/message":
          this.handleMessage(req, res);
          break;
        default:
          res.writeHead(404);
          res.end("Not found");
      }
    });

    this.server.listen(this.port, () => {
      console.log(`SSE server running on http://localhost:${this.port}`);
    });
  }

  serveHTML(req, res) {
    fs.readFile("./index.html", (err, data) => {
      if (err) {
        res.writeHead(404);
        res.end("Not found");
        return;
      }
      res.writeHead(200, { "Content-Type": "text/html" });
      res.end(data);
    });
  }

  handleSSE(req, res) {
    // Set SSE headers
    res.writeHead(200, {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache",
      Connection: "keep-alive",
      "Access-Control-Allow-Origin": "*",
    });

    // Add client to tracking
    const client = {
      res: res,
      id: Date.now(),
    };
    this.clients.add(client);

    console.log(`New SSE client connected. Total: ${this.clients.size}`);

    // Send welcome message
    this.sendToClient(client, {
      event: "system",
      data: "Welcome to SSE server!",
      id: 1,
    });

    // Send periodic updates
    const updateInterval = setInterval(() => {
      this.sendToClient(client, {
        event: "update",
        data: JSON.stringify({
          timestamp: new Date().toISOString(),
          clients: this.clients.size,
        }),
        id: Date.now(),
      });
    }, 3000);

    // Handle client disconnect
    req.on("close", () => {
      console.log(`SSE client disconnected. Total: ${this.clients.size - 1}`);
      clearInterval(updateInterval);
      this.clients.delete(client);
    });
  }

  handleMessage(req, res) {
    if (req.method === "POST") {
      let body = "";
      req.on("data", (chunk) => {
        body += chunk.toString();
      });

      req.on("end", () => {
        try {
          const { message } = JSON.parse(body);

          // Broadcast message to all clients
          this.broadcast({
            event: "chat",
            data: JSON.stringify({
              message: message,
              timestamp: new Date().toISOString(),
            }),
          });

          res.writeHead(200);
          res.end("Message received");
        } catch (error) {
          res.writeHead(400);
          res.end("Invalid JSON");
        }
      });
    } else {
      res.writeHead(405);
      res.end("Method not allowed");
    }
  }

  sendToClient(client, { event = "message", data, id = null, retry = null }) {
    let message = "";
    if (event !== "message") message += `event: ${event}\n`;
    if (id) message += `id: ${id}\n`;
    if (retry) message += `retry: ${retry}\n`;
    message += `data: ${data}\n\n`;

    client.res.write(message);
  }

  broadcast(message) {
    this.clients.forEach((client) => {
      this.sendToClient(client, message);
    });
  }
}

// Start server
new SSEServer(3000);
```

## SSE Event Format

SSE messages follow a specific format:

```
event: update
id: 12345
retry: 10000
data: This is a message
data: that spans multiple lines
```

## Client-Side EventSource API

```javascript
// Creating an EventSource connection
const eventSource = new EventSource("/events");

// Event handlers
eventSource.onopen = (event) => {
  console.log("Connection opened");
};

eventSource.onmessage = (event) => {
  console.log("Message:", event.data);
  console.log("Last Event ID:", event.lastEventId);
};

eventSource.onerror = (event) => {
  console.log("Error:", event);
};

// Custom event listeners
eventSource.addEventListener("update", (event) => {
  const data = JSON.parse(event.data);
  console.log("Update received:", data);
});

eventSource.addEventListener("chat", (event) => {
  const message = JSON.parse(event.data);
  console.log("Chat message:", message);
});

// Close connection
eventSource.close();
```

## Use Cases for SSE

1. **Live news feeds and stock tickers**
2. **Social media notifications**
3. **Real-time dashboard updates**
4. **Progress bars for long operations**
5. **Live sports scores**
6. **Weather updates**
7. **System monitoring dashboards**

## Advantages of SSE

- **Simple implementation**: Built into modern browsers
- **Automatic reconnection**: Handles network issues gracefully
- **HTTP compatibility**: Works with existing HTTP infrastructure
- **Efficient**: No polling overhead
- **Text-friendly**: Perfect for JSON and text data

## Limitations

- **Unidirectional**: Only server → client communication
- **Text-only**: No binary data support
- **Connection limits**: Browser limits on concurrent connections
- **No built-in authentication**: Need to handle separately

## Comparison with WebSockets

| Feature         | SSE            | WebSockets         |
| --------------- | -------------- | ------------------ |
| Direction       | Unidirectional | Bidirectional      |
| Protocol        | HTTP           | WebSocket protocol |
| Data type       | Text only      | Text and binary    |
| Reconnection    | Automatic      | Manual             |
| Complexity      | Simple         | More complex       |
| Browser support | Good           | Excellent          |

SSE is perfect for scenarios where you need simple, efficient server-to-client real-time updates without the complexity of WebSockets!
