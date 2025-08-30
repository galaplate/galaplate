# Background Tasks

Galaplate includes a powerful background job system for processing tasks asynchronously. This system is built on top of a database-backed queue with worker pools for reliable task processing.

## Overview

The background task system consists of:

- **Job Queue**: Database-backed queue for persistent job storage
- **Worker Pool**: Configurable number of workers processing jobs concurrently
- **Job Registry**: Type-safe job registration and resolution
- **Console Commands**: Easy job generation and management via CLI
- **Retry Logic**: Automatic retry with exponential backoff
- **Job States**: Track job lifecycle from pending to completion

## Job Interface

All background jobs must implement the `Job` interface:

```go
type Job interface {
    Type() string                           // Unique job type identifier
    Handle(payload json.RawMessage) error   // Job processing logic
    MaxAttempts() int                       // Maximum retry attempts
    RetryAfter() time.Duration             // Delay between retries
}
```

## Creating Jobs

### Generating Jobs with Console Commands

Use the console command system to quickly generate new background jobs:

```bash
# Generate a new job
go run main.go console make:job EmailJob

# Generate job with specific name
go run main.go console make:job ProcessPaymentJob

# Generate multiple jobs
go run main.go console make:job ImageProcessorJob
go run main.go console make:job NotificationJob

# List all available make commands
go run main.go console list | grep "make:"
```

**Generated job example:**
```go
// pkg/queue/jobs/email_job.go
package jobs

import (
    "encoding/json"
    "time"
    "github.com/galaplate/galaplate/pkg/queue"
)

type EmailJob struct{}

func (j EmailJob) Type() string {
    return "email_job"
}

func (j EmailJob) Handle(payload json.RawMessage) error {
    // Add your job logic here
    return nil
}

func (j EmailJob) MaxAttempts() int {
    return 3
}

func (j EmailJob) RetryAfter() time.Duration {
    return 30 * time.Second
}

func init() {
    queue.RegisterJob(EmailJob{})
}
```

### Basic Job Example

```go
// pkg/queue/jobs/email_job.go
package jobs

import (
    "encoding/json"
    "fmt"
    "time"
)

type EmailJob struct{}

type EmailPayload struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (j EmailJob) Type() string {
    return "send_email"
}

func (j EmailJob) Handle(payload json.RawMessage) error {
    var data EmailPayload
    if err := json.Unmarshal(payload, &data); err != nil {
        return fmt.Errorf("failed to unmarshal email payload: %w", err)
    }

    // Send email logic here
    fmt.Printf("Sending email to %s: %s\n", data.To, data.Subject)

    // Simulate email sending
    time.Sleep(2 * time.Second)

    return nil
}

func (j EmailJob) MaxAttempts() int {
    return 3
}

func (j EmailJob) RetryAfter() time.Duration {
    return 30 * time.Second
}

func init() {
	queue.RegisterJob(EmailJob{})
}
```

**Note**: When you generate a job using `go run main.go console make:job`, the job is automatically registered via the `init()` function.
```

### Advanced Job Example

```go
// pkg/queue/jobs/image_processor.go
package jobs

import (
    "encoding/json"
    "fmt"
    "time"
    "github.com/galaplate/galaplate/logs"
)

type ImageProcessorJob struct{}

type ImagePayload struct {
    ImageURL    string `json:"image_url"`
    UserID      uint   `json:"user_id"`
    ProcessType string `json:"process_type"`
}

func (j ImageProcessorJob) Type() string {
    return "process_image"
}

func (j ImageProcessorJob) Handle(payload json.RawMessage) error {
    var data ImagePayload
    if err := json.Unmarshal(payload, &data); err != nil {
        return fmt.Errorf("failed to unmarshal image payload: %w", err)
    }

    logs.Info("Processing image", map[string]interface{}{
        "user_id": data.UserID,
        "image_url": data.ImageURL,
        "type": data.ProcessType,
    })

    switch data.ProcessType {
    case "resize":
        return j.resizeImage(data.ImageURL)
    case "compress":
        return j.compressImage(data.ImageURL)
    default:
        return fmt.Errorf("unknown process type: %s", data.ProcessType)
    }
}

func (j ImageProcessorJob) MaxAttempts() int {
    return 5
}

func (j ImageProcessorJob) RetryAfter() time.Duration {
    return 1 * time.Minute
}

func (j ImageProcessorJob) resizeImage(url string) error {
    // Image resizing logic
    time.Sleep(5 * time.Second)
    return nil
}

func (j ImageProcessorJob) compressImage(url string) error {
    // Image compression logic
    time.Sleep(3 * time.Second)
    return nil
}

func init() {
	queue.RegisterJob(ImageProcessorJob{})
}
```

### Custom Job Development Workflow

1. **Generate the job structure**:
   ```bash
   go run main.go console make:job ImageProcessorJob
   ```

2. **Implement the job logic** in the generated file:
   ```go
   func (j ImageProcessorJob) Handle(payload json.RawMessage) error {
       // Your custom business logic here
       return nil
   }
   ```

3. **Configure retry behavior**:
   ```go
   func (j ImageProcessorJob) MaxAttempts() int {
       return 5  // Customize based on your needs
   }
   
   func (j ImageProcessorJob) RetryAfter() time.Duration {
       return 1 * time.Minute  // Customize retry delay
   }
```

## Dispatching Jobs

### From Controllers

```go
// pkg/controllers/user_controller.go
package controllers

import (
    "github.com/gofiber/fiber/v2"
    "github.com/galaplate/galaplate/pkg/queue"
    "github.com/galaplate/galaplate/pkg/queue/jobs"
)

func (c *UserController) SendWelcomeEmail(ctx *fiber.Ctx) error {
    userEmail := ctx.FormValue("email")

    // Dispatch email job
    err := queue.Dispatch(jobs.EmailJob{}, map[string]interface{}{
        "to":      userEmail,
        "subject": "Welcome to Galaplate!",
        "body":    "Thank you for joining us.",
    })

    if err != nil {
        return ctx.Status(500).JSON(fiber.Map{
            "error": "Failed to queue welcome email",
        })
    }

    return ctx.JSON(fiber.Map{
        "success": true,
        "message": "Welcome email queued successfully",
    })
}
```

### Immediate Dispatch

```go
// Dispatch job immediately
err := queue.Dispatch(jobs.EmailJob{}, emailData)
```

## Error Handling and Retries

### Retry Logic

Jobs are automatically retried based on their configuration:

```go
func (j EmailJob) MaxAttempts() int {
    return 3  // Retry up to 3 times
}

func (j EmailJob) RetryAfter() time.Duration {
    return 30 * time.Second  // Wait 30 seconds between retries
}
```

### Exponential Backoff

Implement exponential backoff for retries:

```go
func (j EmailJob) RetryAfter() time.Duration {
    // Exponential backoff: 30s, 1m, 2m, 4m, etc.
    attempt := j.getCurrentAttempt() // You'll need to track this
    delay := time.Duration(30 * math.Pow(2, float64(attempt))) * time.Second

    // Cap at maximum delay
    if delay > 10*time.Minute {
        delay = 10 * time.Minute
    }

    return delay
}
```

### Error Types

Handle different error types appropriately:

```go
func (j EmailJob) Handle(payload json.RawMessage) error {
    var data EmailPayload
    if err := json.Unmarshal(payload, &data); err != nil {
        // Permanent error - don't retry
        return queue.NewPermanentError(err)
    }

    if err := j.sendEmail(data); err != nil {
        if isTemporaryError(err) {
            // Temporary error - retry
            return err
        } else {
            // Permanent error - don't retry
            return queue.NewPermanentError(err)
        }
    }

    return nil
}
```

## Testing Background Jobs

### Unit Testing Jobs

```go
// pkg/queue/jobs/email_job_test.go
package jobs

import (
    "encoding/json"
    "testing"
)

func TestEmailJob_Handle(t *testing.T) {
    job := EmailJob{}

    payload := EmailPayload{
        To:      "test@example.com",
        Subject: "Test Email",
        Body:    "Test body",
    }

    payloadJSON, _ := json.Marshal(payload)

    err := job.Handle(payloadJSON)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
}

func TestEmailJob_Type(t *testing.T) {
    job := EmailJob{}
    expected := "send_email"

    if job.Type() != expected {
        t.Errorf("Expected %s, got %s", expected, job.Type())
    }
}
```

### Integration Testing

```go
// Test job dispatch and processing
func TestJobDispatchAndProcessing(t *testing.T) {
    // Setup test database
    setupTestDB()
    defer cleanupTestDB()

    // Register test job
    queue.RegisterJob(TestJob{})

    // Start queue with 1 worker
    q := queue.New(10)
    q.Start(1)
    defer q.Stop()

    // Dispatch job
    err := queue.Dispatch(TestJob{}, testPayload)
    assert.NoError(t, err)

    // Wait for job to complete
    time.Sleep(2 * time.Second)

    // Verify job was processed
    var job models.Job
    err = db.Connect.Where("type = ?", "test_job").First(&job).Error
    assert.NoError(t, err)
    assert.Equal(t, models.JobFinished, job.State)
}
```

## Best Practices

### Job Development with Console Commands

1. **Use code generation** - Always start with `go run main.go console make:job` for consistency
2. **Follow naming conventions** - Use descriptive job names like `ProcessPaymentJob`, `SendEmailJob`
3. **One job per file** - Each job should have its own file in `pkg/queue/jobs/`

### Job Design

1. **Keep jobs idempotent** - Jobs should be safe to run multiple times
2. **Make jobs atomic** - Each job should do one thing well
3. **Handle failures gracefully** - Implement proper error handling
4. **Use appropriate timeouts** - Don't let jobs run indefinitely

### Performance

1. **Batch similar operations** - Group related work together
2. **Use appropriate worker counts** - Balance concurrency with resource usage
3. **Monitor queue depth** - Scale workers based on queue size
4. **Clean up old jobs** - Remove completed jobs to keep database lean

### Monitoring

1. **Log job progress** - Use structured logging for job events
2. **Track job metrics** - Monitor success rates, processing times
3. **Set up alerts** - Alert on job failures or queue backups
4. **Dashboard visibility** - Create admin interfaces for job monitoring

---

## Development Workflow Summary

1. **Generate job**: `go run main.go console make:job YourJobName`
2. **Implement logic**: Add your business logic to the `Handle()` method
3. **Configure retries**: Set appropriate `MaxAttempts()` and `RetryAfter()` values
4. **Test thoroughly**: Write unit tests for your job logic
5. **Deploy and monitor**: Use logging to track job performance

## Console Commands Reference

```bash
# Job Development
go run main.go console make:job <JobName>    # Generate new background job
go run main.go console list                  # List all available commands

# Related Commands
go run main.go console make:cron <CronName>  # Generate CRON job
go run main.go console db:seed               # Run database seeders
```

## Next Steps

- **[Console Commands](/console-commands)** - Master the command-line tools
- **[Database](/database)** - Understand job storage and models
- **[API Reference](/api-reference)** - Build APIs that dispatch background jobs
- **[Project Structure](/project-structure)** - Understand where jobs are organized
