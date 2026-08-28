# Background Tasks

Galaplate provides a database-backed queue with worker pools for processing jobs asynchronously.

## Creating Jobs

Generate a job:

```bash
go run main.go console make:job SendEmail
```

This creates `pkg/jobs/send_email.go`:

```go
package jobs

import (
    "context"
    "encoding/json"
    "time"
    "github.com/galaplate/core/queue"
)

type SendEmail struct{}

func (j SendEmail) Type() string {
    return "send_email"
}

func (j SendEmail) Handle(ctx context.Context, payload json.RawMessage) error {
    // Job logic here — ctx carries cancellation signals from the queue worker
    return nil
}

func (j SendEmail) MaxAttempts() int {
    return 3
}

func (j SendEmail) RetryAfter() time.Duration {
    return 2 * time.Minute
}

func init() {
    queue.RegisterJob(SendEmail{})
}
```

## Dispatching Jobs

```go
import "github.com/galaplate/core/queue"

queue.Dispatch(jobs.SendEmail{}, map[string]any{
    "to":      "user@example.com",
    "subject": "Welcome",
})
```

## Job States

Jobs are stored in the database with these states:

| State | Description |
|-------|-------------|
| `pending` | Waiting to be processed |
| `started` | Currently being processed |
| `finished` | Completed successfully |
| `failed` | Failed after max attempts |

## Job Model

```go
type Job struct {
    ID          uint
    Type        string
    Payload     json.RawMessage
    State       JobState
    ErrorMsg    string
    Attempts    int
    AvailableAt time.Time
    CreatedAt   time.Time
    StartedAt   *time.Time
    FinishedAt  *time.Time
}
```

## Scheduling

Generate a cron task:

```bash
go run main.go console make:cron DailyCleanup
```

This creates `pkg/scheduler/daily_cleanup.go`:

```go
package scheduler

import (
    "context"
    "github.com/galaplate/core/scheduler"
)

type DailyCleanup struct{}

func (c DailyCleanup) Handle() (string, func(context.Context)) {
    return "0 0 * * *", func(ctx context.Context) {
        // Cleanup logic — ctx carries the scheduler lifecycle
    }
}

func init() {
    scheduler.RegisterScheduler("daily_cleanup", DailyCleanup{})
}
```

## Context Propagation

Queue jobs and scheduler tasks receive a `context.Context` that carries cancellation signals from the application lifecycle.

Use `ctx` to make cancellable database queries or HTTP requests:

```go
func (j SendEmail) Handle(ctx context.Context, payload json.RawMessage) error {
    db := database.Connect.WithContext(ctx)
    // queries will be cancelled if the worker shuts down
    db.Find(&users)
    return nil
}
```

Cron tasks are registered automatically via `init()` and started when the application boots with background jobs enabled.

## Worker Configuration

Workers are configured via `bootstrap.AppConfig`:

```go
app := bootstrap.NewApp(func(cfg *bootstrap.AppConfig) {
    cfg.QueueSize = 100
    cfg.WorkerCount = 5
})
```

The queue polls the database every second for pending jobs.
