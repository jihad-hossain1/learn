# HTTP/2 Server Push

HTTP/2 Server Push allows servers to **proactively send resources** to clients before they're requested, reducing latency and improving page load times.

## Key Concepts

- **Proactive Resource Delivery**: Server sends resources it knows the client will need
- **Reduced Round Trips**: Eliminates the request-response cycle for pushed resources
- **Cache Awareness**: Pushed resources can be cached by the client
- **Connection Multiplexing**: Uses HTTP/2's multiplexing capabilities

## How HTTP/2 Server Push Works

### Traditional HTTP/1.1 vs HTTP/2 Push

**HTTP/1.1:**

```
Client: GET /index.html
Server: Sends index.html
Client: Parses HTML, finds need for style.css
Client: GET /style.css
Server: Sends style.css
Client: Parses HTML, finds need for app.js
Client: GET /app.js
Server: Sends app.js
```

**HTTP/2 with Server Push:**

```
Client: GET /index.html
Server: Sends index.html + PUSH_PROMISE for style.css + app.js
Server: Sends style.css + app.js without waiting for requests
```

## Basic Implementation Examples

### Node.js with HTTP/2 Module

```javascript
const http2 = require("http2");
const fs = require("fs");

const server = http2.createSecureServer({
  key: fs.readFileSync("localhost-privkey.pem"),
  cert: fs.readFileSync("localhost-cert.pem"),
});

server.on("stream", (stream, headers) => {
  const path = headers[":path"];

  // Serve main HTML file
  if (path === "/index.html") {
    // Push critical resources
    pushResource(stream, "/style.css");
    pushResource(stream, "/app.js");
    pushResource(stream, "/header.jpg");

    // Serve main HTML
    stream.respond({
      "content-type": "text/html; charset=utf-8",
      ":status": 200,
    });
    stream.end(`
      <!DOCTYPE html>
      <html>
      <head>
          <title>HTTP/2 Push Demo</title>
          <link rel="stylesheet" href="/style.css">
      </head>
      <body>
          <h1>HTTP/2 Server Push Demo</h1>
          <img src="/header.jpg" alt="Header">
          <script src="/app.js"></script>
      </body>
      </html>
    `);
  }

  // Handle direct resource requests
  else if (path === "/style.css") {
    stream.respond({ ":status": 200, "content-type": "text/css" });
    stream.end(`
      body { font-family: Arial, sans-serif; margin: 40px; }
      h1 { color: #2c3e50; }
    `);
  } else if (path === "/app.js") {
    stream.respond({
      ":status": 200,
      "content-type": "application/javascript",
    });
    stream.end(`
      console.log('App loaded via HTTP/2 Push!');
      document.addEventListener('DOMContentLoaded', function() {
          console.log('DOM fully loaded');
      });
    `);
  }
});

function pushResource(stream, path) {
  const pushHeaders = { ":path": path };

  stream.pushStream(pushHeaders, (err, pushStream) => {
    if (err) {
      console.error("Push error:", err);
      return;
    }

    console.log("Pushing resource:", path);

    if (path === "/style.css") {
      pushStream.respond({ ":status": 200, "content-type": "text/css" });
      pushStream.end(`
        body { font-family: Arial, sans-serif; margin: 40px; }
        h1 { color: #2c3e50; }
        .pushed { color: #27ae60; font-weight: bold; }
      `);
    } else if (path === "/app.js") {
      pushStream.respond({
        ":status": 200,
        "content-type": "application/javascript",
      });
      pushStream.end(`
        console.log('This resource was pushed by the server!');
        document.addEventListener('DOMContentLoaded', function() {
            const element = document.createElement('div');
            element.className = 'pushed';
            element.textContent = 'This content was enabled by HTTP/2 Server Push!';
            document.body.appendChild(element);
        });
      `);
    }
  });
}

server.listen(8443, () => {
  console.log("HTTP/2 Server running on https://localhost:8443");
});
```

## Advanced HTTP/2 Server with Smart Push

```javascript
const http2 = require("http2");
const fs = require("fs");
const path = require("path");

class HTTP2Server {
  constructor() {
    this.server = null;
    this.resourceCache = new Map();
    this.init();
  }

  init() {
    const options = {
      key: fs.readFileSync("./certs/localhost-privkey.pem"),
      cert: fs.readFileSync("./certs/localhost-cert.pem"),
      allowHTTP1: true, // Fallback for HTTP/1.1 clients
    };

    this.server = http2.createSecureServer(options);
    this.setupRoutes();
    this.preloadResources();
  }

  preloadResources() {
    const resources = {
      "/style.css": "text/css",
      "/app.js": "application/javascript",
      "/critical.css": "text/css",
      "/api/data.json": "application/json",
    };

    Object.keys(resources).forEach((resourcePath) => {
      this.loadResource(resourcePath, resources[resourcePath]);
    });
  }

  loadResource(resourcePath, contentType) {
    // In real application, you'd read from file system
    const mockData = this.generateMockResource(resourcePath);
    this.resourceCache.set(resourcePath, {
      data: mockData,
      contentType: contentType,
      lastModified: new Date(),
    });
  }

  generateMockResource(resourcePath) {
    switch (resourcePath) {
      case "/style.css":
        return `
          body { 
            font-family: 'Arial', sans-serif; 
            margin: 0; 
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
          }
          .container { max-width: 800px; margin: 0 auto; }
          .card { 
            background: rgba(255,255,255,0.1); 
            padding: 20px; 
            border-radius: 10px;
            margin: 20px 0;
          }
        `;

      case "/critical.css":
        return `
          .critical { 
            display: block !important;
            color: #ffeb3b;
            font-size: 1.2em;
          }
        `;

      case "/app.js":
        return `
          class HTTP2Demo {
            constructor() {
              this.init();
            }
            
            init() {
              console.log('HTTP/2 Push Demo Loaded');
              this.setupMetrics();
              this.loadAdditionalResources();
            }
            
            setupMetrics() {
              performance.mark('http2-push-start');
              window.addEventListener('load', () => {
                performance.mark('http2-push-end');
                performance.measure('pageLoad', 'http2-push-start', 'http2-push-end');
                this.displayMetrics();
              });
            }
            
            displayMetrics() {
              const measure = performance.getEntriesByName('pageLoad')[0];
              const metricsDiv = document.getElementById('metrics');
              if (metricsDiv) {
                metricsDiv.innerHTML = \`
                  <h3>Performance Metrics</h3>
                  <p>Page Load Time: \${measure.duration.toFixed(2)}ms</p>
                  <p>Resources pushed: 4</p>
                  <p>Protocol: \${window.performance.protocol}</p>
                \`;
              }
            }
            
            loadAdditionalResources() {
              // Non-critical resources loaded normally
              fetch('/api/data.json')
                .then(r => r.json())
                .then(data => {
                  console.log('Additional data loaded:', data);
                });
            }
          }
          
          new HTTP2Demo();
        `;

      case "/api/data.json":
        return JSON.stringify({
          message: "This data was pushed via HTTP/2",
          timestamp: new Date().toISOString(),
          features: ["Multiplexing", "Server Push", "Header Compression"],
        });

      default:
        return "Resource not found";
    }
  }

  setupRoutes() {
    this.server.on("stream", (stream, headers) => {
      const method = headers[":method"];
      const path = headers[":path"];

      console.log(`${method} ${path} - HTTP/${stream.session.alpnProtocol}`);

      if (method === "GET") {
        this.handleGetRequest(stream, path, headers);
      } else {
        stream.respond({ ":status": 405 });
        stream.end();
      }
    });
  }

  handleGetRequest(stream, path, headers) {
    // Main page - push critical resources
    if (path === "/" || path === "/index.html") {
      this.serveHTML(stream);
    }
    // Direct resource requests
    else if (this.resourceCache.has(path)) {
      this.serveResource(stream, path);
    } else {
      stream.respond({ ":status": 404 });
      stream.end("Not found");
    }
  }

  serveHTML(stream) {
    console.log("Serving HTML with HTTP/2 Push");

    // Push critical resources first
    this.pushCriticalResources(stream);

    // Then serve HTML
    const html = `
      <!DOCTYPE html>
      <html>
      <head>
          <title>HTTP/2 Server Push Advanced Demo</title>
          <link rel="stylesheet" href="/critical.css">
          <link rel="stylesheet" href="/style.css">
      </head>
      <body>
          <div class="container">
              <h1>🚀 HTTP/2 Server Push Demo</h1>
              
              <div class="card">
                  <div class="critical">
                      This critical CSS was pushed via HTTP/2!
                  </div>
                  
                  <div id="metrics">
                      <p>Loading metrics...</p>
                  </div>
                  
                  <h3>Pushed Resources:</h3>
                  <ul id="resource-list">
                      <li>Critical CSS (pushed)</li>
                      <li>Main Stylesheet (pushed)</li>
                      <li>JavaScript Application (pushed)</li>
                      <li>API Data (pushed)</li>
                  </ul>
              </div>
              
              <div class="card">
                  <h3>HTTP/2 Benefits:</h3>
                  <ul>
                      <li>✅ Reduced Latency</li>
                      <li>✅ Multiplexing</li>
                      <li>✅ Header Compression</li>
                      <li>✅ Server Push</li>
                  </ul>
              </div>
          </div>
          
          <script src="/app.js"></script>
      </body>
      </html>
    `;

    stream.respond({
      "content-type": "text/html; charset=utf-8",
      ":status": 200,
    });
    stream.end(html);
  }

  pushCriticalResources(stream) {
    const criticalResources = [
      "/critical.css",
      "/style.css",
      "/app.js",
      "/api/data.json",
    ];

    criticalResources.forEach((resourcePath) => {
      this.pushResource(stream, resourcePath);
    });
  }

  pushResource(stream, resourcePath) {
    if (!this.resourceCache.has(resourcePath)) return;

    const resource = this.resourceCache.get(resourcePath);

    stream.pushStream({ ":path": resourcePath }, (err, pushStream) => {
      if (err) {
        console.error(`Push error for ${resourcePath}:`, err);
        return;
      }

      console.log(`Pushing: ${resourcePath}`);

      pushStream.respond({
        ":status": 200,
        "content-type": resource.contentType,
        "cache-control": "public, max-age=3600",
        "last-modified": resource.lastModified.toUTCString(),
      });

      pushStream.end(resource.data);
    });
  }

  serveResource(stream, resourcePath) {
    const resource = this.resourceCache.get(resourcePath);

    stream.respond({
      ":status": 200,
      "content-type": resource.contentType,
      "cache-control": "public, max-age=3600",
      "last-modified": resource.lastModified.toUTCString(),
    });

    stream.end(resource.data);
  }

  start(port = 8443) {
    this.server.listen(port, () => {
      console.log(`🚀 HTTP/2 Server running on https://localhost:${port}`);
      console.log("Features: Server Push, Multiplexing, Header Compression");
    });
  }
}

// Generate self-signed certificates for demo
function generateDemoCerts() {
  const certsDir = "./certs";
  if (!fs.existsSync(certsDir)) {
    fs.mkdirSync(certsDir);
    console.log("Note: Please generate SSL certificates for production use");
    console.log("For demo, you can use mkcert or openssl to create certs");
  }
}

// Start server
generateDemoCerts();
const server = new HTTP2Server();
server.start();
```

## Client-Side Detection and Monitoring

```html
<!DOCTYPE html>
<html>
  <head>
    <title>HTTP/2 Push Detection</title>
    <style>
      .metric {
        background: #f5f5f5;
        padding: 10px;
        margin: 5px 0;
        border-left: 4px solid #4caf50;
      }
      .http2 {
        border-color: #4caf50;
      }
      .http1 {
        border-color: #ff9800;
      }
    </style>
  </head>
  <body>
    <h1>HTTP/2 Capability Detection</h1>

    <div id="protocol-info" class="metric">
      <h3>Protocol Detection</h3>
      <p id="protocol-status">Detecting...</p>
    </div>

    <div id="push-info" class="metric">
      <h3>Server Push Information</h3>
      <div id="push-resources"></div>
    </div>

    <div id="performance" class="metric">
      <h3>Performance Metrics</h3>
      <div id="timing-data"></div>
    </div>

    <script>
      class HTTP2Detector {
        constructor() {
          this.supported = false;
          this.pushedResources = new Set();
          this.detectProtocol();
          this.setupPerformanceMonitoring();
        }

        detectProtocol() {
          // Check if HTTP/2 is supported
          const protocol = performance.protocol;
          const statusElement = document.getElementById("protocol-status");

          if (protocol === "h2") {
            this.supported = true;
            statusElement.innerHTML = "✅ HTTP/2 Supported";
            statusElement.parentElement.className = "metric http2";
          } else {
            statusElement.innerHTML = "⚠️ HTTP/1.1 - No Server Push";
            statusElement.parentElement.className = "metric http1";
          }
        }

        setupPerformanceMonitoring() {
          // Monitor resource loading
          const observer = new PerformanceObserver((list) => {
            list.getEntries().forEach((entry) => {
              if (
                entry.initiatorType === "script" ||
                entry.initiatorType === "link"
              ) {
                this.trackResource(entry);
              }
            });
          });

          observer.observe({ entryTypes: ["resource"] });

          // Track page load performance
          window.addEventListener("load", () => {
            this.displayPerformanceMetrics();
          });
        }

        trackResource(entry) {
          const resourcesDiv = document.getElementById("push-resources");
          const resourceInfo = document.createElement("div");

          resourceInfo.innerHTML = `
                    <strong>${entry.name}</strong><br>
                    Size: ${(entry.transferSize / 1024).toFixed(2)}KB<br>
                    Duration: ${entry.duration.toFixed(2)}ms
                `;

          resourcesDiv.appendChild(resourceInfo);
        }

        displayPerformanceMetrics() {
          const timing = performance.timing;
          const loadTime = timing.loadEventEnd - timing.navigationStart;

          document.getElementById("timing-data").innerHTML = `
                    Total Load Time: ${loadTime}ms<br>
                    DOM Ready: ${
                      timing.domContentLoadedEventEnd - timing.navigationStart
                    }ms<br>
                    Resources: ${
                      performance.getEntriesByType("resource").length
                    }
                `;
        }
      }

      // Initialize detector
      new HTTP2Detector();
    </script>
  </body>
</html>
```

## Best Practices for HTTP/2 Server Push

### 1. Push Only Critical Resources

```javascript
// Good: Push critical above-the-fold resources
const criticalResources = [
  "/critical.css",
  "/above-the-fold.js",
  "/hero-image.jpg",
];

// Avoid: Pushing everything
const allResources = [
  "/style.css",
  "/app.js",
  "/lib1.js",
  "/lib2.js",
  "/image1.jpg",
  "/image2.png",
  "/non-critical.css",
];
```

### 2. Cache Management

```javascript
function pushWithCacheHeaders(stream, resourcePath, content) {
  stream.pushStream({ ":path": resourcePath }, (err, pushStream) => {
    pushStream.respond({
      ":status": 200,
      "content-type": "text/css",
      "cache-control": "public, max-age=31536000", // 1 year
      etag: generateETag(content),
    });
    pushStream.end(content);
  });
}
```

### 3. Smart Push Based on User Agent

```javascript
function shouldPushForClient(userAgent) {
  // Only push for HTTP/2 capable clients
  const supportsHTTP2 = !userAgent.includes("HTTP/1.1");
  const isMobile = userAgent.includes("Mobile");

  // Push fewer resources for mobile
  return supportsHTTP2 && !isMobile;
}
```

## Use Cases for HTTP/2 Server Push

1. **Critical CSS**: Push above-the-fold styles
2. **Web Fonts**: Push fonts needed for initial render
3. **Hero Images**: Push important images
4. **JavaScript Frameworks**: Push core framework code
5. **API Data**: Push initial data payloads

## Limitations and Considerations

- **Over-pushing**: Can waste bandwidth if client already has resources
- **Cache Digests**: Need careful cache management
- **Browser Support**: Not all proxies handle push correctly
- **Complexity**: More complex than traditional resource loading

HTTP/2 Server Push is powerful for performance optimization when used judiciously for critical resources that benefit from reduced latency!
