package migration

import (
	"fmt"
	"sync"
	"time"
)

// SourceType identifies which platform a migration job reads from.
type SourceType string

const (
	SourceWordPress   SourceType = "wordpress"
	SourceShopify     SourceType = "shopify"
	SourceWooCommerce SourceType = "woocommerce"
	SourceMagento     SourceType = "magento"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

// Job is a single migration run, tracked here for the dashboard's benefit.
// Actual execution happens in the Node/TS worker (worker/), which owns the
// real BullMQ/Redis queue; the API reaches it over HTTP via TriggerWorker
// (see worker_client.go) since Go can't write BullMQ's queue format
// directly.
type Job struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	SourceType  SourceType `json:"sourceType"`
	TargetDBURL string     `json:"-"`
	Status      Status     `json:"status"`
	WorkerJobID string     `json:"workerJobId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

var (
	mu   sync.Mutex
	jobs = map[string]*Job{}
	seq  int
)

// EnqueueJob records a new migration job and returns it with a generated ID.
func EnqueueJob(j Job) *Job {
	mu.Lock()
	defer mu.Unlock()

	seq++
	j.ID = fmt.Sprintf("mig_%04d", seq)
	j.Status = StatusQueued
	j.CreatedAt = time.Now().UTC()

	stored := j
	jobs[j.ID] = &stored
	return &stored
}

// ListJobs returns every known migration job.
func ListJobs() []*Job {
	mu.Lock()
	defer mu.Unlock()

	list := make([]*Job, 0, len(jobs))
	for _, j := range jobs {
		list = append(list, j)
	}
	return list
}

// SetWorkerJobID records which BullMQ job (on the worker side) corresponds
// to a given Aura Amber job, and marks it running once the worker has
// actually accepted it.
func SetWorkerJobID(id, workerJobID string) {
	mu.Lock()
	defer mu.Unlock()

	if j, ok := jobs[id]; ok {
		j.WorkerJobID = workerJobID
		j.Status = StatusRunning
	}
}

// MarkFailed records that handing a job off to the worker failed.
func MarkFailed(id string) {
	mu.Lock()
	defer mu.Unlock()

	if j, ok := jobs[id]; ok {
		j.Status = StatusFailed
	}
}
