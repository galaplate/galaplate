package queue

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/sheenazien8/goplate-core/database"
	"github.com/sheenazien8/goplate-core/models"
)

type Task func()

type Queue struct {
	tasks chan Task
	wg    sync.WaitGroup
}

func New(bufferSize int) *Queue {
	return &Queue{
		tasks: make(chan Task, bufferSize),
	}
}

func (q *Queue) Start(workerCount int) {
	if database.Connect.Migrator().HasTable(&models.Job{}) {
		for i := 0; i < workerCount; i++ {
			go func() {
				for {
					var jobRecord models.Job
					err := database.Connect.
						Where("state = ? AND available_at <= ?", models.JobPending, time.Now()).
						Order("created_at ASC").
						First(&jobRecord).Error

					if err != nil {
						time.Sleep(1 * time.Second)
						continue
					}

					job, err := ResolveJob(jobRecord.Type, jobRecord.Payload)
					if err != nil {
						failJob(&jobRecord, err)
						continue
					}

					start := time.Now()
					database.Connect.Model(&jobRecord).Updates(models.Job{
						State:     models.JobStarted,
						StartedAt: &start,
					})

					err = job.Handle(jobRecord.Payload)

					if err != nil {
						jobRecord.Attempts++
						if jobRecord.Attempts >= job.MaxAttempts() {
							failJob(&jobRecord, err)
						} else {
							database.Connect.Model(&jobRecord).Updates(models.Job{
								State:       models.JobPending,
								ErrorMsg:    err.Error(),
								Attempts:    jobRecord.Attempts,
								AvailableAt: time.Now().Add(job.RetryAfter()),
							})
						}
					} else {
						database.Connect.Model(&jobRecord).Updates(models.Job{
							State:      models.JobFinished,
							FinishedAt: ptr(time.Now()),
						})
					}
				}
			}()
		}

	}
}

func failJob(job *models.Job, err error) {
	database.Connect.Model(job).Updates(models.Job{
		State:      models.JobFailed,
		ErrorMsg:   err.Error(),
		FinishedAt: ptr(time.Now()),
	})
}

func ptr[T any](v T) *T {
	return &v
}

var registry = map[string]func() Job{}

func RegisterJob(job Job) {
	registry[job.Type()] = func() Job {
		return job
	}
}

func ResolveJob(typeName string, payload json.RawMessage) (Job, error) {
	creator, exists := registry[typeName]
	if !exists {
		return nil, fmt.Errorf("job type '%s' not registered", typeName)
	}

	job := creator()

	// // optional: decode payload into the job struct if needed
	// if err := json.Unmarshal(payload, &job); err != nil {
	// 	return nil, fmt.Errorf("failed to unmarshal payload for '%s': %w", typeName, err)
	// }

	return job, nil
}

type JobEnqueueRequest struct {
	Type    string
	Payload any
}

type Job interface {
	Type() string
	Handle(payload json.RawMessage) error
	MaxAttempts() int
	RetryAfter() time.Duration
}

func Dispatch(job Job, params ...any) error {
	_, err := SaveJobToDB(JobEnqueueRequest{
		Type:    job.Type(),
		Payload: params,
	})
	if err != nil {
		return fmt.Errorf("failed to save job to DB: %w", err)
	}

	return nil
}

func SaveJobToDB(req JobEnqueueRequest) (*models.Job, error) {
	payloadJSON, err := json.Marshal(req.Payload)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	job := models.Job{
		Type:        req.Type,
		Payload:     payloadJSON,
		State:       models.JobPending,
		Attempts:    0,
		AvailableAt: now,
		CreatedAt:   now,
	}

	if err := database.Connect.Create(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}
