# WebRTC (Web Real-Time Communication) - Advanced Peer-to-Peer

WebRTC enables **direct peer-to-peer communication** between browsers and devices without intermediaries, supporting audio, video, and data channels.

## Core WebRTC Architecture

### Key Components

- **RTCPeerConnection**: Manages peer-to-peer connection
- **RTCDataChannel**: Bi-directional data transfer
- **MediaStream**: Audio/video media handling
- **Signaling**: Coordination between peers

## Advanced WebRTC Implementation

### Complete Signaling Server (Node.js + Socket.io)

```javascript
const express = require("express");
const http = require("http");
const socketIo = require("socket.io");
const path = require("path");

class WebRTCSignalingServer {
  constructor(port = 3000) {
    this.app = express();
    this.server = http.createServer(this.app);
    this.io = socketIo(this.server, {
      cors: { origin: "*" },
    });
    this.port = port;
    this.rooms = new Map();

    this.setupMiddleware();
    this.setupSocketHandlers();
  }

  setupMiddleware() {
    this.app.use(express.static(path.join(__dirname, "public")));

    this.app.get("/api/rooms", (req, res) => {
      res.json({
        totalRooms: this.rooms.size,
        rooms: Array.from(this.rooms.entries()).map(([roomId, room]) => ({
          roomId,
          userCount: room.users.size,
          users: Array.from(room.users.keys()),
        })),
      });
    });
  }

  setupSocketHandlers() {
    this.io.on("connection", (socket) => {
      console.log("User connected:", socket.id);

      socket.on("join-room", (roomId, userId, userData) => {
        this.handleJoinRoom(socket, roomId, userId, userData);
      });

      socket.on("offer", (data) => {
        socket.to(data.roomId).emit("offer", {
          offer: data.offer,
          from: data.from,
          to: data.to,
        });
      });

      socket.on("answer", (data) => {
        socket.to(data.roomId).emit("answer", {
          answer: data.answer,
          from: data.from,
        });
      });

      socket.on("ice-candidate", (data) => {
        socket.to(data.roomId).emit("ice-candidate", {
          candidate: data.candidate,
          from: data.from,
        });
      });

      socket.on("data-channel-message", (data) => {
        socket.to(data.roomId).emit("data-channel-message", {
          message: data.message,
          type: data.type,
          from: data.from,
          timestamp: new Date().toISOString(),
        });
      });

      socket.on("disconnect", () => {
        this.handleDisconnect(socket);
      });
    });
  }

  handleJoinRoom(socket, roomId, userId, userData) {
    // Leave any existing rooms
    Array.from(socket.rooms)
      .filter((room) => room !== socket.id)
      .forEach((room) => socket.leave(room));

    // Join new room
    socket.join(roomId);

    if (!this.rooms.has(roomId)) {
      this.rooms.set(roomId, { users: new Map() });
    }

    const room = this.rooms.get(roomId);
    room.users.set(userId, { socketId: socket.id, userData });

    // Notify others in room
    socket.to(roomId).emit("user-joined", {
      userId,
      userData,
      roomSize: room.users.size,
    });

    // Send current room state to new user
    socket.emit("room-state", {
      users: Array.from(room.users.entries())
        .map(([id, data]) => ({
          userId: id,
          userData: data.userData,
        }))
        .filter((user) => user.userId !== userId),
    });

    console.log(`User ${userId} joined room ${roomId}`);
  }

  handleDisconnect(socket) {
    // Find and remove user from rooms
    this.rooms.forEach((room, roomId) => {
      room.users.forEach((userData, userId) => {
        if (userData.socketId === socket.id) {
          room.users.delete(userId);
          socket.to(roomId).emit("user-left", { userId });

          if (room.users.size === 0) {
            this.rooms.delete(roomId);
          }
        }
      });
    });

    console.log("User disconnected:", socket.id);
  }

  start() {
    this.server.listen(this.port, () => {
      console.log(`WebRTC Signaling Server running on port ${this.port}`);
    });
  }
}

module.exports = WebRTCSignalingServer;
```

### Advanced WebRTC Client Class

```javascript
class AdvancedWebRTCClient {
  constructor() {
    this.localStream = null;
    this.remoteStreams = new Map();
    this.peerConnections = new Map();
    this.dataChannels = new Map();
    this.socket = null;
    this.config = {
      iceServers: [
        { urls: "stun:stun.l.google.com:19302" },
        { urls: "stun:stun1.l.google.com:19302" },
      ],
      sdpSemantics: "unified-plan",
    };
    this.stats = {
      connections: 0,
      messagesSent: 0,
      messagesReceived: 0,
      dataChannelState: new Map(),
    };
  }

  // Initialize connection to signaling server
  async initialize(socketUrl, roomId, userId, userData = {}) {
    this.roomId = roomId;
    this.userId = userId;
    this.userData = userData;

    await this.connectSignaling(socketUrl);
    await this.initializeMedia();
  }

  // Connect to signaling server
  async connectSignaling(socketUrl) {
    return new Promise((resolve, reject) => {
      this.socket = io(socketUrl);

      this.socket.on("connect", () => {
        console.log("Connected to signaling server");
        this.setupSignalingHandlers();
        this.socket.emit("join-room", this.roomId, this.userId, this.userData);
        resolve();
      });

      this.socket.on("connect_error", (error) => {
        console.error("Signaling connection error:", error);
        reject(error);
      });
    });
  }

  // Initialize media devices
  async initializeMedia(
    constraints = {
      video: { width: 1280, height: 720, frameRate: 30 },
      audio: {
        echoCancellation: true,
        noiseSuppression: true,
        autoGainControl: true,
      },
    }
  ) {
    try {
      this.localStream = await navigator.mediaDevices.getUserMedia(constraints);
      this.onLocalStreamAcquired(this.localStream);
    } catch (error) {
      console.error("Error accessing media devices:", error);
      throw error;
    }
  }

  // Setup signaling event handlers
  setupSignalingHandlers() {
    this.socket.on("user-joined", (data) => {
      console.log("User joined:", data.userId);
      this.createPeerConnection(data.userId);
    });

    this.socket.on("user-left", (data) => {
      console.log("User left:", data.userId);
      this.cleanupPeerConnection(data.userId);
    });

    this.socket.on("offer", async (data) => {
      await this.handleOffer(data.from, data.offer);
    });

    this.socket.on("answer", async (data) => {
      await this.handleAnswer(data.from, data.answer);
    });

    this.socket.on("ice-candidate", (data) => {
      this.handleIceCandidate(data.from, data.candidate);
    });

    this.socket.on("data-channel-message", (data) => {
      this.handleDataChannelMessage(data.from, data);
    });

    this.socket.on("room-state", (data) => {
      data.users.forEach((user) => {
        this.createPeerConnection(user.userId);
      });
    });
  }

  // Create peer connection
  async createPeerConnection(remoteUserId) {
    if (this.peerConnections.has(remoteUserId)) {
      console.log("Peer connection already exists for:", remoteUserId);
      return;
    }

    try {
      const pc = new RTCPeerConnection(this.config);
      this.peerConnections.set(remoteUserId, pc);
      this.stats.connections++;

      // Add local stream tracks
      if (this.localStream) {
        this.localStream.getTracks().forEach((track) => {
          pc.addTrack(track, this.localStream);
        });
      }

      // Setup data channel
      this.setupDataChannel(pc, remoteUserId);

      // Setup event handlers
      this.setupPeerConnectionHandlers(pc, remoteUserId);

      // Create and send offer if we're the initiator
      if (this.userId < remoteUserId) {
        // Simple initiator logic
        await this.createOffer(remoteUserId);
      }
    } catch (error) {
      console.error("Error creating peer connection:", error);
    }
  }

  // Setup data channel
  setupDataChannel(pc, remoteUserId) {
    const dataChannel = pc.createDataChannel(`chat-${this.userId}`, {
      ordered: true,
      maxPacketLifeTime: 1000, // 1 second
    });

    this.dataChannels.set(remoteUserId, dataChannel);
    this.setupDataChannelHandlers(dataChannel, remoteUserId);
  }

  // Setup data channel event handlers
  setupDataChannelHandlers(dc, remoteUserId) {
    dc.onopen = () => {
      console.log("Data channel opened with:", remoteUserId);
      this.stats.dataChannelState.set(remoteUserId, "open");
      this.onDataChannelOpen(remoteUserId);
    };

    dc.onclose = () => {
      console.log("Data channel closed with:", remoteUserId);
      this.stats.dataChannelState.set(remoteUserId, "closed");
      this.onDataChannelClose(remoteUserId);
    };

    dc.onmessage = (event) => {
      this.stats.messagesReceived++;
      const message = JSON.parse(event.data);
      this.onDataChannelMessage(remoteUserId, message);
    };

    dc.onerror = (error) => {
      console.error("Data channel error:", error);
      this.onDataChannelError(remoteUserId, error);
    };
  }

  // Setup peer connection event handlers
  setupPeerConnectionHandlers(pc, remoteUserId) {
    pc.onicecandidate = (event) => {
      if (event.candidate) {
        this.socket.emit("ice-candidate", {
          roomId: this.roomId,
          to: remoteUserId,
          from: this.userId,
          candidate: event.candidate,
        });
      }
    };

    pc.ontrack = (event) => {
      const remoteStream = event.streams[0];
      this.remoteStreams.set(remoteUserId, remoteStream);
      this.onRemoteStream(remoteUserId, remoteStream);
    };

    pc.onconnectionstatechange = () => {
      console.log(`Connection state with ${remoteUserId}:`, pc.connectionState);
      this.onConnectionStateChange(remoteUserId, pc.connectionState);

      if (pc.connectionState === "connected") {
        this.onPeerConnected(remoteUserId);
      } else if (
        pc.connectionState === "disconnected" ||
        pc.connectionState === "failed"
      ) {
        this.onPeerDisconnected(remoteUserId);
      }
    };

    pc.oniceconnectionstatechange = () => {
      console.log(
        `ICE connection state with ${remoteUserId}:`,
        pc.iceConnectionState
      );
      this.onIceConnectionStateChange(remoteUserId, pc.iceConnectionState);
    };

    // Handle incoming data channel
    pc.ondatachannel = (event) => {
      const dc = event.channel;
      this.dataChannels.set(remoteUserId, dc);
      this.setupDataChannelHandlers(dc, remoteUserId);
    };
  }

  // Create and send offer
  async createOffer(remoteUserId) {
    const pc = this.peerConnections.get(remoteUserId);
    if (!pc) return;

    try {
      const offer = await pc.createOffer();
      await pc.setLocalDescription(offer);

      this.socket.emit("offer", {
        roomId: this.roomId,
        to: remoteUserId,
        from: this.userId,
        offer: offer,
      });
    } catch (error) {
      console.error("Error creating offer:", error);
    }
  }

  // Handle incoming offer
  async handleOffer(remoteUserId, offer) {
    let pc = this.peerConnections.get(remoteUserId);
    if (!pc) {
      await this.createPeerConnection(remoteUserId);
      pc = this.peerConnections.get(remoteUserId);
    }

    try {
      await pc.setRemoteDescription(offer);
      const answer = await pc.createAnswer();
      await pc.setLocalDescription(answer);

      this.socket.emit("answer", {
        roomId: this.roomId,
        to: remoteUserId,
        from: this.userId,
        answer: answer,
      });
    } catch (error) {
      console.error("Error handling offer:", error);
    }
  }

  // Handle incoming answer
  async handleAnswer(remoteUserId, answer) {
    const pc = this.peerConnections.get(remoteUserId);
    if (!pc) return;

    try {
      await pc.setRemoteDescription(answer);
    } catch (error) {
      console.error("Error handling answer:", error);
    }
  }

  // Handle ICE candidate
  async handleIceCandidate(remoteUserId, candidate) {
    const pc = this.peerConnections.get(remoteUserId);
    if (!pc || !candidate) return;

    try {
      await pc.addIceCandidate(candidate);
    } catch (error) {
      console.error("Error adding ICE candidate:", error);
    }
  }

  // Send data channel message
  sendData(remoteUserId, message, type = "text") {
    const dc = this.dataChannels.get(remoteUserId);
    if (!dc || dc.readyState !== "open") {
      console.warn("Data channel not ready");
      return false;
    }

    try {
      const data = {
        type: type,
        message: message,
        from: this.userId,
        timestamp: Date.now(),
      };

      dc.send(JSON.stringify(data));
      this.stats.messagesSent++;

      // Also send via signaling for reliability
      this.socket.emit("data-channel-message", {
        roomId: this.roomId,
        to: remoteUserId,
        from: this.userId,
        message: message,
        type: type,
      });

      return true;
    } catch (error) {
      console.error("Error sending data:", error);
      return false;
    }
  }

  // Handle incoming data channel message
  handleDataChannelMessage(remoteUserId, data) {
    this.onDataChannelMessage(remoteUserId, data);
  }

  // File sharing via data channel
  async sendFile(remoteUserId, file) {
    const dc = this.dataChannels.get(remoteUserId);
    if (!dc || dc.readyState !== "open") {
      throw new Error("Data channel not ready");
    }

    // Send file metadata first
    const metadata = {
      type: "file-metadata",
      fileName: file.name,
      fileSize: file.size,
      fileType: file.type,
      timestamp: Date.now(),
    };

    this.sendData(remoteUserId, metadata, "file-metadata");

    // Read and send file in chunks
    const chunkSize = 16384; // 16KB chunks
    const fileReader = new FileReader();
    let offset = 0;

    return new Promise((resolve, reject) => {
      fileReader.onload = (event) => {
        try {
          dc.send(event.target.result);
          offset += event.target.result.byteLength;

          this.onFileProgress(remoteUserId, file.name, offset, file.size);

          if (offset < file.size) {
            readNextChunk();
          } else {
            resolve();
          }
        } catch (error) {
          reject(error);
        }
      };

      fileReader.onerror = reject;

      const readNextChunk = () => {
        const slice = file.slice(offset, offset + chunkSize);
        fileReader.readAsArrayBuffer(slice);
      };

      readNextChunk();
    });
  }

  // Cleanup peer connection
  cleanupPeerConnection(remoteUserId) {
    const pc = this.peerConnections.get(remoteUserId);
    if (pc) {
      pc.close();
      this.peerConnections.delete(remoteUserId);
    }

    this.dataChannels.delete(remoteUserId);
    this.remoteStreams.delete(remoteUserId);
    this.stats.dataChannelState.delete(remoteUserId);
    this.stats.connections--;

    this.onPeerDisconnected(remoteUserId);
  }

  // Get connection statistics
  getStats() {
    return {
      ...this.stats,
      peerConnections: this.peerConnections.size,
      dataChannels: Array.from(this.dataChannels.entries()).map(
        ([userId, dc]) => ({
          userId,
          state: dc.readyState,
        })
      ),
    };
  }

  // Event handlers (to be overridden by application)
  onLocalStreamAcquired(stream) {
    console.log("Local stream acquired");
  }

  onRemoteStream(userId, stream) {
    console.log("Remote stream received from:", userId);
  }

  onDataChannelOpen(userId) {
    console.log("Data channel opened with:", userId);
  }

  onDataChannelClose(userId) {
    console.log("Data channel closed with:", userId);
  }

  onDataChannelMessage(userId, message) {
    console.log("Message from", userId, ":", message);
  }

  onDataChannelError(userId, error) {
    console.error("Data channel error with", userId, ":", error);
  }

  onPeerConnected(userId) {
    console.log("Peer connected:", userId);
  }

  onPeerDisconnected(userId) {
    console.log("Peer disconnected:", userId);
  }

  onConnectionStateChange(userId, state) {
    console.log("Connection state changed for", userId, ":", state);
  }

  onIceConnectionStateChange(userId, state) {
    console.log("ICE connection state changed for", userId, ":", state);
  }

  onFileProgress(userId, fileName, loaded, total) {
    const percent = ((loaded / total) * 100).toFixed(1);
    console.log(`File ${fileName}: ${percent}%`);
  }
}
```

### Advanced WebRTC Application Example

```html
<!DOCTYPE html>
<html>
  <head>
    <title>Advanced WebRTC Demo</title>
    <style>
      * {
        box-sizing: border-box;
        margin: 0;
        padding: 0;
      }
      body {
        font-family: Arial, sans-serif;
        background: #1a1a1a;
        color: white;
      }
      .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 20px;
      }
      .header {
        text-align: center;
        margin-bottom: 30px;
      }
      .video-container {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
        gap: 20px;
        margin-bottom: 30px;
      }
      .video-wrapper {
        background: #2a2a2a;
        border-radius: 10px;
        padding: 10px;
        position: relative;
      }
      .video-wrapper.local video {
        transform: scaleX(-1);
      }
      .video-title {
        position: absolute;
        top: 10px;
        left: 10px;
        background: rgba(0, 0, 0, 0.7);
        padding: 5px 10px;
        border-radius: 5px;
      }
      video {
        width: 100%;
        height: 200px;
        object-fit: cover;
        border-radius: 5px;
      }
      .controls {
        background: #2a2a2a;
        padding: 20px;
        border-radius: 10px;
        margin-bottom: 20px;
      }
      .control-group {
        margin-bottom: 15px;
      }
      button {
        background: #4caf50;
        color: white;
        border: none;
        padding: 10px 20px;
        margin: 5px;
        border-radius: 5px;
        cursor: pointer;
      }
      button:hover {
        background: #45a049;
      }
      button:disabled {
        background: #666;
        cursor: not-allowed;
      }
      .stats {
        background: #2a2a2a;
        padding: 15px;
        border-radius: 10px;
        font-family: monospace;
      }
      .chat {
        background: #2a2a2a;
        padding: 20px;
        border-radius: 10px;
      }
      #messages {
        height: 200px;
        overflow-y: auto;
        margin-bottom: 10px;
        border: 1px solid #444;
        padding: 10px;
      }
      .message {
        margin: 5px 0;
        padding: 5px;
        background: #333;
        border-radius: 3px;
      }
    </style>
  </head>
  <body>
    <div class="container">
      <div class="header">
        <h1>🚀 Advanced WebRTC Demo</h1>
        <p>
          Real-time peer-to-peer communication with video, audio, and data
          channels
        </p>
      </div>

      <div class="controls">
        <div class="control-group">
          <button id="joinBtn">Join Room</button>
          <button id="leaveBtn" disabled>Leave Room</button>
          <button id="shareScreenBtn">Share Screen</button>
          <button id="toggleVideoBtn">Toggle Video</button>
          <button id="toggleAudioBtn">Toggle Audio</button>
        </div>

        <div class="control-group">
          <input
            type="text"
            id="messageInput"
            placeholder="Type a message..."
            disabled
          />
          <button id="sendMsgBtn" disabled>Send Message</button>
          <input type="file" id="fileInput" disabled />
          <button id="sendFileBtn" disabled>Send File</button>
        </div>
      </div>

      <div class="video-container">
        <div class="video-wrapper local">
          <div class="video-title">You</div>
          <video id="localVideo" autoplay muted></video>
        </div>
        <div id="remoteVideos"></div>
      </div>

      <div class="stats">
        <h3>Connection Statistics</h3>
        <div id="statsDisplay"></div>
      </div>

      <div class="chat">
        <h3>Chat</h3>
        <div id="messages"></div>
      </div>
    </div>

    <script src="/socket.io/socket.io.js"></script>
    <script>
      class WebRTCApp {
        constructor() {
          this.webrtcClient = new AdvancedWebRTCClient();
          this.roomId = "advanced-demo-room";
          this.userId = "user-" + Math.random().toString(36).substr(2, 9);
          this.isConnected = false;

          this.setupEventHandlers();
          this.overrideClientHandlers();
        }

        setupEventHandlers() {
          document
            .getElementById("joinBtn")
            .addEventListener("click", () => this.joinRoom());
          document
            .getElementById("leaveBtn")
            .addEventListener("click", () => this.leaveRoom());
          document
            .getElementById("shareScreenBtn")
            .addEventListener("click", () => this.shareScreen());
          document
            .getElementById("toggleVideoBtn")
            .addEventListener("click", () => this.toggleVideo());
          document
            .getElementById("toggleAudioBtn")
            .addEventListener("click", () => this.toggleAudio());
          document
            .getElementById("sendMsgBtn")
            .addEventListener("click", () => this.sendMessage());
          document
            .getElementById("sendFileBtn")
            .addEventListener("click", () => this.sendFile());

          document
            .getElementById("messageInput")
            .addEventListener("keypress", (e) => {
              if (e.key === "Enter") this.sendMessage();
            });
        }

        overrideClientHandlers() {
          // Override WebRTC client event handlers
          this.webrtcClient.onLocalStreamAcquired = (stream) => {
            const localVideo = document.getElementById("localVideo");
            localVideo.srcObject = stream;
            this.updateUI();
          };

          this.webrtcClient.onRemoteStream = (userId, stream) => {
            this.addRemoteVideo(userId, stream);
          };

          this.webrtcClient.onDataChannelMessage = (userId, message) => {
            this.displayMessage(userId, message);
          };

          this.webrtcClient.onPeerConnected = (userId) => {
            this.displaySystemMessage(`User ${userId} connected`);
            this.updateStats();
          };

          this.webrtcClient.onPeerDisconnected = (userId) => {
            this.removeRemoteVideo(userId);
            this.displaySystemMessage(`User ${userId} disconnected`);
            this.updateStats();
          };
        }

        async joinRoom() {
          try {
            await this.webrtcClient.initialize("/", this.roomId, this.userId, {
              name: `User${Math.floor(Math.random() * 1000)}`,
            });

            this.isConnected = true;
            this.updateUI();
            this.displaySystemMessage("Joined room successfully");

            // Start stats update loop
            this.statsInterval = setInterval(() => this.updateStats(), 2000);
          } catch (error) {
            console.error("Error joining room:", error);
            this.displaySystemMessage("Error joining room: " + error.message);
          }
        }

        async leaveRoom() {
          if (this.webrtcClient.socket) {
            this.webrtcClient.socket.disconnect();
          }

          this.webrtcClient.peerConnections.forEach((pc, userId) => {
            this.webrtcClient.cleanupPeerConnection(userId);
          });

          if (this.localStream) {
            this.localStream.getTracks().forEach((track) => track.stop());
          }

          this.isConnected = false;
          this.updateUI();

          if (this.statsInterval) {
            clearInterval(this.statsInterval);
          }

          this.displaySystemMessage("Left room");
        }

        async shareScreen() {
          try {
            const screenStream = await navigator.mediaDevices.getDisplayMedia({
              video: true,
              audio: true,
            });

            // Replace local stream
            if (this.webrtcClient.localStream) {
              this.webrtcClient.localStream
                .getTracks()
                .forEach((track) => track.stop());
            }

            this.webrtcClient.localStream = screenStream;
            document.getElementById("localVideo").srcObject = screenStream;

            // Update all peer connections with new stream
            this.webrtcClient.peerConnections.forEach((pc, userId) => {
              const videoSender = pc
                .getSenders()
                .find(
                  (sender) => sender.track && sender.track.kind === "video"
                );

              if (videoSender) {
                const videoTrack = screenStream.getVideoTracks()[0];
                videoSender.replaceTrack(videoTrack);
              }
            });

            screenStream.getVideoTracks()[0].onended = () => {
              this.webrtcClient.initializeMedia(); // Restore camera
            };
          } catch (error) {
            console.error("Error sharing screen:", error);
          }
        }

        toggleVideo() {
          if (!this.webrtcClient.localStream) return;

          const videoTrack = this.webrtcClient.localStream.getVideoTracks()[0];
          if (videoTrack) {
            videoTrack.enabled = !videoTrack.enabled;
            document.getElementById("toggleVideoBtn").textContent =
              videoTrack.enabled ? "Disable Video" : "Enable Video";
          }
        }

        toggleAudio() {
          if (!this.webrtcClient.localStream) return;

          const audioTrack = this.webrtcClient.localStream.getAudioTracks()[0];
          if (audioTrack) {
            audioTrack.enabled = !audioTrack.enabled;
            document.getElementById("toggleAudioBtn").textContent =
              audioTrack.enabled ? "Mute Audio" : "Unmute Audio";
          }
        }

        sendMessage() {
          const input = document.getElementById("messageInput");
          const message = input.value.trim();

          if (message) {
            this.webrtcClient.peerConnections.forEach((pc, userId) => {
              this.webrtcClient.sendData(userId, message, "chat");
            });

            this.displayMessage(this.userId, { message, type: "chat" });
            input.value = "";
          }
        }

        async sendFile() {
          const fileInput = document.getElementById("fileInput");
          const file = fileInput.files[0];

          if (!file) return;

          try {
            this.webrtcClient.peerConnections.forEach((pc, userId) => {
              this.webrtcClient.sendFile(userId, file);
            });

            this.displaySystemMessage(`Sending file: ${file.name}`);
            fileInput.value = "";
          } catch (error) {
            console.error("Error sending file:", error);
            this.displaySystemMessage("Error sending file: " + error.message);
          }
        }

        addRemoteVideo(userId, stream) {
          const remoteVideos = document.getElementById("remoteVideos");

          if (document.getElementById(`video-${userId}`)) {
            return; // Video already exists
          }

          const videoWrapper = document.createElement("div");
          videoWrapper.className = "video-wrapper remote";
          videoWrapper.id = `video-wrapper-${userId}`;

          const videoTitle = document.createElement("div");
          videoTitle.className = "video-title";
          videoTitle.textContent = userId;

          const video = document.createElement("video");
          video.id = `video-${userId}`;
          video.autoplay = true;
          video.srcObject = stream;

          videoWrapper.appendChild(videoTitle);
          videoWrapper.appendChild(video);
          remoteVideos.appendChild(videoWrapper);
        }

        removeRemoteVideo(userId) {
          const videoWrapper = document.getElementById(
            `video-wrapper-${userId}`
          );
          if (videoWrapper) {
            videoWrapper.remove();
          }
        }

        displayMessage(userId, data) {
          const messages = document.getElementById("messages");
          const messageDiv = document.createElement("div");
          messageDiv.className = "message";

          const timestamp = new Date().toLocaleTimeString();
          messageDiv.innerHTML = `
                    <strong>${userId}</strong> 
                    [${timestamp}]: 
                    ${
                      data.type === "file-metadata"
                        ? `File: ${data.message.fileName}`
                        : data.message
                    }
                `;

          messages.appendChild(messageDiv);
          messages.scrollTop = messages.scrollHeight;
        }

        displaySystemMessage(message) {
          const messages = document.getElementById("messages");
          const messageDiv = document.createElement("div");
          messageDiv.className = "message system";
          messageDiv.innerHTML = `<em>${message}</em>`;
          messages.appendChild(messageDiv);
          messages.scrollTop = messages.scrollHeight;
        }

        updateStats() {
          const stats = this.webrtcClient.getStats();
          const statsDisplay = document.getElementById("statsDisplay");

          statsDisplay.innerHTML = `
                    Connections: ${stats.connections} |
                    Messages Sent: ${stats.messagesSent} |
                    Messages Received: ${stats.messagesReceived} |
                    Data Channels: ${stats.dataChannels.length}
                `;
        }

        updateUI() {
          document.getElementById("joinBtn").disabled = this.isConnected;
          document.getElementById("leaveBtn").disabled = !this.isConnected;
          document.getElementById("messageInput").disabled = !this.isConnected;
          document.getElementById("sendMsgBtn").disabled = !this.isConnected;
          document.getElementById("fileInput").disabled = !this.isConnected;
          document.getElementById("sendFileBtn").disabled = !this.isConnected;
        }
      }

      // Initialize application when page loads
      window.addEventListener("load", () => {
        window.webrtcApp = new WebRTCApp();
      });
    </script>
  </body>
</html>
```

## Advanced WebRTC Features

### 1. TURN Server Configuration

```javascript
const advancedConfig = {
  iceServers: [
    { urls: "stun:stun.l.google.com:19302" },
    {
      urls: "turn:your-turn-server.com:3478",
      username: "your-username",
      credential: "your-credential",
    },
  ],
  iceTransportPolicy: "all", // or 'relay' for firewall restrictions
  bundlePolicy: "max-bundle",
  rtcpMuxPolicy: "require",
};
```

### 2. Quality Control and Bandwidth Management

```javascript
// Adjust video quality based on network conditions
function adjustVideoQuality(pc, bandwidth) {
  const senders = pc.getSenders();
  const videoSender = senders.find((s) => s.track && s.track.kind === "video");

  if (videoSender && videoSender.track) {
    const parameters = videoSender.getParameters();
    if (!parameters.encodings) {
      parameters.encodings = [{}];
    }

    parameters.encodings[0].maxBitrate = bandwidth * 1000; // Convert to bps
    videoSender.setParameters(parameters);
  }
}
```

### 3. Recording and Media Processing

```javascript
// Record media streams
class MediaRecorder {
  constructor(stream) {
    this.recorder = new MediaRecorder(stream, {
      mimeType: "video/webm;codecs=vp9",
      videoBitsPerSecond: 2500000,
    });
    this.chunks = [];

    this.recorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        this.chunks.push(event.data);
      }
    };
  }

  start() {
    this.recorder.start(1000); // Collect data every second
  }

  stop() {
    return new Promise((resolve) => {
      this.recorder.onstop = () => {
        const blob = new Blob(this.chunks, { type: "video/webm" });
        resolve(blob);
      };
      this.recorder.stop();
    });
  }
}
```

## WebRTC Use Cases

1. **Video Conferencing**: Multi-party video calls
2. **Live Streaming**: Real-time broadcasting
3. **Gaming**: Multiplayer online games
4. **Remote Desktop**: Screen sharing and control
5. **IoT Communication**: Device-to-device data transfer
6. **File Sharing**: Peer-to-peer file transfer
7. **Collaborative Tools**: Real-time document editing

## Key Advantages

- **Low Latency**: Direct peer-to-peer communication
- **No Intermediaries**: Reduced server costs and complexity
- **Secure**: End-to-end encryption
- **Flexible**: Supports various media types and data channels
- **Standardized**: Cross-browser compatibility

WebRTC provides powerful peer-to-peer capabilities for real-time communication applications!
