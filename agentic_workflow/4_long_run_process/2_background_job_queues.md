**Long-Running Processes** with Background Job Queues like BullMQ and Trigger.dev.

## BullMQ (Redis-based Queue)

BullMQ is a Redis-based queue system that excels at handling long-running processes with reliability and scalability.

### Key Features for Long-Running Processes:

- **Job persistence** in Redis
- **Progress tracking**
- **Retry mechanisms**
- **Priority queues**
- **Delayed jobs**

### JavaScript Example: E-commerce Order Processing Pipeline

```javascript
// queue-setup.js
import { Queue, Worker, QueueEvents } from "bullmq";
import IORedis from "ioredis";

const connection = new IORedis(process.env.REDIS_URL);

// Create queues for different stages
export const orderQueue = new Queue("order processing", { connection });
export const paymentQueue = new Queue("payment processing", { connection });
export const notificationQueue = new Queue("notifications", { connection });

// Order Processing Worker
const orderWorker = new Worker(
  "order processing",
  async (job) => {
    try {
      const { orderId, userId, items } = job.data;

      // Update job progress
      await job.updateProgress(10);

      // Step 1: Validate order
      await job.log("Validating order...");
      await validateOrder(items);
      await job.updateProgress(30);

      // Step 2: Reserve inventory (can take time)
      await job.log("Reserving inventory...");
      await reserveInventory(items);
      await job.updateProgress(60);

      // Step 3: Create payment job
      await paymentQueue.add(
        "process-payment",
        {
          orderId,
          userId,
          amount: calculateTotal(items),
        },
        {
          delay: 1000, // 1 second delay
          jobId: `payment-${orderId}`,
        }
      );

      await job.updateProgress(90);
      await job.log("Order validation completed");

      return { orderId, status: "validated" };
    } catch (error) {
      await job.log(`Error: ${error.message}`);
      throw error; // Will trigger retry mechanism
    }
  },
  {
    connection,
    concurrency: 5, // Process 5 jobs concurrently
    limiter: {
      max: 10, // Max 10 jobs per second
      duration: 1000,
    },
  }
);

// Payment Processing Worker (Long-running example)
const paymentWorker = new Worker(
  "payment processing",
  async (job) => {
    const { orderId, userId, amount } = job.data;

    // Simulate a long-running payment process (5-30 seconds)
    await job.updateProgress(0);
    await job.log("Starting payment processing...");

    // Step 1: Validate payment method
    await simulateLongProcess("Validating payment method", 2000, job, 20);

    // Step 2: Process with payment gateway (can be slow)
    await simulateLongProcess("Processing with payment gateway", 5000, job, 50);

    // Step 3: Verify transaction
    await simulateLongProcess("Verifying transaction", 3000, job, 80);

    // Step 4: Create notification job
    await notificationQueue.add("payment-success", {
      orderId,
      userId,
      amount,
    });

    await job.updateProgress(100);
    return { success: true, transactionId: `txn_${Date.now()}` };
  },
  {
    connection,
    concurrency: 3, // Fewer concurrent payments due to external API limits
    attempts: 3, // Retry up to 3 times
    backoff: {
      type: "exponential",
      delay: 2000,
    },
  }
);

// Helper function for simulating long processes
async function simulateLongProcess(stepName, duration, job, progress) {
  await job.log(stepName);
  await new Promise((resolve) => setTimeout(resolve, duration));
  await job.updateProgress(progress);
}

// Queue event listeners for monitoring
const queueEvents = new QueueEvents("order processing", { connection });

queueEvents.on("completed", ({ jobId, returnvalue }) => {
  console.log(`Job ${jobId} completed:`, returnvalue);
});

queueEvents.on("failed", ({ jobId, failedReason }) => {
  console.error(`Job ${jobId} failed:`, failedReason);
});

queueEvents.on("progress", ({ jobId, data }) => {
  console.log(`Job ${jobId} progress: ${data}%`);
});

// Usage: Adding jobs to the queue
export async function createOrder(orderData) {
  const job = await orderQueue.add("process-order", orderData, {
    jobId: `order-${orderData.orderId}`,
    removeOnComplete: 50, // Keep only last 50 completed jobs
    removeOnFail: 10, // Keep only last 10 failed jobs
    attempts: 3,
    delay: 1000, // Start after 1 second
  });

  return job;
}

// Example usage
const orderData = {
  orderId: "12345",
  userId: "user-67890",
  items: [
    { productId: "prod-1", quantity: 2, price: 25 },
    { productId: "prod-2", quantity: 1, price: 50 },
  ],
};

// createOrder(orderData);
```

## Trigger.dev

Trigger.dev is a modern background jobs framework that's particularly good for **very long-running processes** (hours/days) and **event-driven workflows**.

### Key Features:

- **Extremely long timeouts** (up to 1 year)
- **Event-driven workflows**
- **Integrated with popular APIs**
- **Visual debugging and monitoring**

### JavaScript Example: Multi-day Video Processing Pipeline

```javascript
// triggers/video-processing.ts
import { trigger, eventTrigger, intervalTrigger } from "@trigger.dev/sdk";
import { client } from "../trigger";
import { z } from "zod";

// Define events
export const videoUploadedEvent = client.defineEvent({
  name: "video.uploaded",
  schema: z.object({
    videoId: z.string(),
    userId: z.string(),
    filePath: z.string(),
    format: z.string(),
    duration: z.number(),
  }),
});

// Long-running video processing workflow
client.defineJob({
  id: "video-processing-pipeline",
  name: "Video Processing Pipeline",
  version: "0.1.0",
  enabled: true,

  // Trigger when video is uploaded
  trigger: eventTrigger({
    name: videoUploadedEvent.name,
  }),

  run: async (payload, io, ctx) => {
    const { videoId, userId, filePath, format, duration } = payload;

    // This job can run for hours/days without timing out

    // Step 1: Validate video file (could take minutes for large files)
    const validation = await io.runTask("validate-video", async () => {
      // Simulate validation process
      await io.wait("wait-for-validation", 2 * 60); // Wait 2 minutes
      return { valid: true, format: 'mp4', size: '2.1GB' };
    });

    if (!validation.valid) {
      await io.sendEvent("video-validation-failed", {
        name: "video.processing.failed",
        payload: { videoId, reason: "Invalid format" },
      });
      return;
    }

    // Step 2: Transcode to multiple formats (could take hours)
    const transcodingPromises = [
      io.runTask("transcode-1080p", async () => {
        await io.wait("transcode-1080p-wait", 30 * 60); // 30 minutes
        return { format: '1080p', status: 'completed' };
      }),

      io.runTask("transcode-720p", async () => {
        await io.wait("transcode-720p-wait", 20 * 60); // 20 minutes
        return { format: '720p', status: 'completed' };
      }),

      io.runTask("transcode-480p", async () => {
        await io.wait("transcode-480p-wait", 15 * 60); // 15 minutes
        return { format: '480p', status: 'completed' };
      })
    ];

    const transcodingResults = await Promise.all(transcodingPromises);

    // Step 3: Generate thumbnails (parallel processing)
    const thumbnailTasks = [];
    for (let i = 0; i < 5; i++) {
      thumbnailTasks.push(
        io.runTask(`generate-thumbnail-${i}`, async () => {
          const timestamp = i * (duration / 5);
          await io.wait(`thumbnail-wait-${i}`, 2 * 60); // 2 minutes per thumbnail
          return { thumbnailId: `thumb-${i}`, timestamp, status: 'generated' };
        })
      );
    }

    const thumbnails = await Promise.all(thumbnailTasks);

    // Step 4: AI Content Analysis (could take variable time)
    const analysis = await io.runTask("ai-content-analysis", async () => {
      // This could take hours for long videos
      await io.wait("ai-analysis-wait", 60 * 60); // 1 hour wait
      return {
        objectsDetected: ['person', 'car', 'building'],
        sceneChanges: 45,
        contentRating: 'G'
      };
    });

    // Step 5: Update database and notify user
    await io.runTask("update-database", async () => {
      // Update video status to processed
      return { status: 'completed', processedAt: new Date() };
    });

    // Send completion event
    await io.sendEvent("video-processing-complete", {
      name: "video.processing.completed",
      payload: {
        videoId,
        userId,
        formats: transcodingResults,
        thumbnails,
        analysis,
        completedAt: new Date().toISOString()
      },
    });

    return {
      videoId,
      status: "fully processed",
      processingTime: "several hours", // Actually could be days
      outputs: transcodingResults.length + thumbnails.length
    };
  },
});

// Scheduled cleanup job (runs daily)
client.defineJob({
  id: "cleanup-old-videos",
  name: "Cleanup Old Videos",
  version: "0.1.0",

  trigger: intervalTrigger({
    seconds: 24 * 60 * 60, // Run daily
  }),

  run: async (payload, io, ctx) => {
    // This job can run for hours processing thousands of records

    await io.logger.info("Starting cleanup of old videos...");

    // Step 1: Find videos older than 30 days
    const oldVideos = await io.runTask("find-old-videos", async () => {
      // Database query to find old videos
      return [{ videoId: 'vid1' }, { videoId: 'vid2' }]; // Simplified
    });

    // Step 2: Delete from storage (parallel processing)
    const deletionPromises = oldVideos.map(video =>
      io.runTask(`delete-video-${video.videoId}`, async () => {
        await io.wait(`delete-wait-${video.videoId}`, 30); // 30 seconds per deletion
        return { videoId: video.videoId, deleted: true };
      })
    );

    // Process in batches of 10 to avoid overwhelming the system
    const batchSize = 10;
    for (let i = 0; i < deletionPromises.length; i += batchSize) {
      const batch = deletionPromises.slice(i, i + batchSize);
      await Promise.all(batch);
      await io.wait(`batch-wait-${i/batchSize}`, 60); // Wait 1 minute between batches
    }

    // Step 3: Update database records
    await io.runTask("update-database-records", async () => {
      // Mark videos as archived in database
      return { updated: oldVideos.length };
    });

    await io.logger.info(`Cleanup completed: ${oldVideos.length} videos processed`);

    return { videosProcessed: oldVideos.length };
  },
});

// trigger-client.js
import { TriggerClient } from "@trigger.dev/sdk";

export const client = new TriggerClient({
  id: "video-platform",
  apiKey: process.env.TRIGGER_API_KEY!,
  apiUrl: process.env.TRIGGER_API_URL,
});

// Usage: Trigger the workflow
import { videoUploadedEvent } from "./triggers/video-processing";

async function handleVideoUpload(videoData) {
  // Send event to trigger the long-running workflow
  await client.sendEvent({
    name: videoUploadedEvent.name,
    payload: {
      videoId: videoData.id,
      userId: videoData.userId,
      filePath: videoData.path,
      format: videoData.format,
      duration: videoData.duration,
    },
  });
}
```

## Comparison: BullMQ vs Trigger.dev for Long-Running Processes

### BullMQ is better for:

- **Redis-based infrastructure** you already have
- **Fine-grained control** over queue mechanics
- **High-throughput** processing needs
- **Custom retry and backoff** strategies
- **Priority-based** job processing

### Trigger.dev is better for:

- **Extremely long-running** processes (hours to months)
- **Event-driven architectures**
- **API-heavy workflows** with built-in integrations
- **Visual debugging** and monitoring
- **Serverless environments**
- **Complex workflow dependencies**

## Real-World Use Cases:

**BullMQ Examples:**

- E-commerce order processing
- Batch email sending
- PDF generation and processing
- Data synchronization between systems

**Trigger.dev Examples:**

- Multi-day video rendering pipelines
- AI model training workflows
- Complex data migration processes
- Multi-step approval workflows that take weeks

Both are excellent choices depending on your specific requirements for long-running processes!
