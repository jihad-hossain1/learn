# Real-Time Communication with WebSockets

WebSockets provide a **bi-directional, full-duplex communication channel** over a single TCP connection, enabling real-time data exchange between clients and servers.

## Key WebSocket Concepts

- **Bi-directional**: Both client and server can send messages at any time
- **Persistent Connection**: Remains open until explicitly closed
- **Low Latency**: Minimal overhead compared to HTTP polling
- **Real-time**: Instant data transmission

## Simple WebSocket Example

Here's a basic implementation with a Node.js server and JavaScript client:

### Server (Node.js with ws library)

```javascript
const WebSocket = require("ws");
const server = new WebSocket.Server({ port: 8080 });

server.on("connection", (socket) => {
  console.log("Client connected");

  // Send welcome message to client
  socket.send("Welcome to the WebSocket server!");

  // Handle messages from client
  socket.on("message", (message) => {
    console.log(`Received: ${message}`);

    // Echo message back to client
    socket.send(`Server received: ${message}`);
  });

  // Handle connection close
  socket.on("close", () => {
    console.log("Client disconnected");
  });
});

console.log("WebSocket server running on ws://localhost:8080");
```

### Client (Browser JavaScript)

```javascript
// Create WebSocket connection
const socket = new WebSocket("ws://localhost:8080");

// Connection opened
socket.addEventListener("open", (event) => {
  console.log("Connected to server");
  socket.send("Hello Server!");
});

// Listen for messages
socket.addEventListener("message", (event) => {
  console.log("Message from server:", event.data);

  // Display message on page
  const messages = document.getElementById("messages");
  messages.innerHTML += `<p>${event.data}</p>`;
});

// Connection closed
socket.addEventListener("close", (event) => {
  console.log("Disconnected from server");
});

// Error handling
socket.addEventListener("error", (event) => {
  console.error("WebSocket error:", event);
});

// Send message function
function sendMessage() {
  const input = document.getElementById("messageInput");
  const message = input.value;
  if (message) {
    socket.send(message);
    input.value = "";
  }
}
```

### HTML for Client

```html
<!DOCTYPE html>
<html>
  <head>
    <title>WebSocket Client</title>
  </head>
  <body>
    <h1>WebSocket Chat</h1>
    <div
      id="messages"
      style="border: 1px solid #ccc; height: 300px; overflow-y: scroll;"
    ></div>
    <input type="text" id="messageInput" placeholder="Type your message" />
    <button onclick="sendMessage()">Send</button>

    <script>
      // Client JavaScript code from above
    </script>
  </body>
</html>
```

## Advanced Example: Real-Time Chat Application

### Enhanced Server with Multiple Clients

```javascript
const WebSocket = require("ws");
const server = new WebSocket.Server({ port: 8080 });

const clients = new Set();

server.on("connection", (socket) => {
  clients.add(socket);
  console.log(`New client connected. Total clients: ${clients.size}`);

  // Broadcast to all clients when someone joins
  broadcast(`${clients.size} users online`, socket);

  socket.on("message", (message) => {
    try {
      const data = JSON.parse(message);

      if (data.type === "chat") {
        // Broadcast chat message to all clients
        broadcast(
          JSON.stringify({
            type: "chat",
            user: data.user,
            message: data.message,
            timestamp: new Date().toISOString(),
          }),
          socket
        );
      }
    } catch (error) {
      console.error("Error parsing message:", error);
    }
  });

  socket.on("close", () => {
    clients.delete(socket);
    console.log(`Client disconnected. Total clients: ${clients.size}`);
    broadcast(`${clients.size} users online`);
  });
});

function broadcast(message, sender = null) {
  clients.forEach((client) => {
    if (client !== sender && client.readyState === WebSocket.OPEN) {
      client.send(message);
    }
  });
}
```

### Enhanced Client

```javascript
class ChatClient {
  constructor() {
    this.socket = new WebSocket("ws://localhost:8080");
    this.user = `User${Math.floor(Math.random() * 1000)}`;
    this.setupEventListeners();
  }

  setupEventListeners() {
    this.socket.addEventListener("open", () => {
      this.displaySystemMessage("Connected to chat");
    });

    this.socket.addEventListener("message", (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.type === "chat") {
          this.displayChatMessage(data.user, data.message, data.timestamp);
        }
      } catch {
        this.displaySystemMessage(event.data);
      }
    });

    this.socket.addEventListener("close", () => {
      this.displaySystemMessage("Disconnected from chat");
    });
  }

  sendChatMessage(message) {
    if (this.socket.readyState === WebSocket.OPEN) {
      this.socket.send(
        JSON.stringify({
          type: "chat",
          user: this.user,
          message: message,
        })
      );
    }
  }

  displayChatMessage(user, message, timestamp) {
    const messagesDiv = document.getElementById("messages");
    const time = new Date(timestamp).toLocaleTimeString();
    messagesDiv.innerHTML += `
      <div class="message">
        <span class="time">[${time}]</span>
        <span class="user">${user}:</span>
        <span class="text">${message}</span>
      </div>
    `;
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }

  displaySystemMessage(message) {
    const messagesDiv = document.getElementById("messages");
    messagesDiv.innerHTML += `
      <div class="system-message">${message}</div>
    `;
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
  }
}

// Initialize chat client
const chat = new ChatClient();

// Send message when button is clicked
document.getElementById("sendButton").addEventListener("click", () => {
  const input = document.getElementById("messageInput");
  const message = input.value.trim();
  if (message) {
    chat.sendChatMessage(message);
    input.value = "";
  }
});

// Send message on Enter key
document.getElementById("messageInput").addEventListener("keypress", (e) => {
  if (e.key === "Enter") {
    document.getElementById("sendButton").click();
  }
});
```

## WebSocket API Methods and Events

### Client-Side API

```javascript
// Creating a connection
const socket = new WebSocket("ws://example.com/socket");

// Events
socket.onopen = (event) => {
  console.log("Connection established");
};

socket.onmessage = (event) => {
  console.log("Data received:", event.data);
};

socket.onclose = (event) => {
  console.log("Connection closed", event.code, event.reason);
};

socket.onerror = (error) => {
  console.error("WebSocket error:", error);
};

// Methods
socket.send("Hello Server"); // Send string
socket.send(JSON.stringify({ data: "object" })); // Send object
socket.send(blobData); // Send binary data

socket.close(1000, "Work complete"); // Close connection
```

## Use Cases for WebSockets

1. **Real-time chat applications**
2. **Live sports scores and updates**
3. **Multiplayer online games**
4. **Collaborative editing tools**
5. **Financial trading platforms**
6. **Live customer support**
7. **IoT device monitoring**

## Installation and Setup

To run the server examples, install the ws library:

```bash
npm install ws
```

## Key Advantages over HTTP Polling

- **Lower latency**: No need to repeatedly establish connections
- **Reduced overhead**: No HTTP headers with each message
- **True real-time**: Instant push notifications
- **Bi-directional**: Server can initiate communication

This covers the fundamentals of real-time communication with WebSockets! The technology is perfect for applications requiring instant data updates and bidirectional communication.
