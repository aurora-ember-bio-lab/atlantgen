package migration

import (
	"database/sql"
	"fmt"
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

// Job is a single migration run, persisted in Postgres.
type Job struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	SourceType  SourceType `json:"sourceType"`
	TargetDBURL string     `json:"-"`
	Status      Status     `json:"status"`
	WorkerJobID string     `json:"workerJobId,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
}

var db *sql.DB

// InitDB initializes the database connection for the migration package.
func InitDB(conn *sql.DB) {
	db = conn
}

func getDB() *sql.DB {
	if db == nil {
		return nil
	}
	return db
}

// EnqueueJob inserts a new migration job into Postgres.
func EnqueueJob(j Job) *Job {
	conn := getDB()
	if conn == nil {
		return fallbackEnqueue(j)
	}

	j.ID = fmt.Sprintf("mig_%04d", time.Now().UnixNano()%100000)
	j.Status = StatusQueued
	j.CreatedAt = time.Now().UTC()
	j.UpdatedAt = j.CreatedAt

	_, err := conn.Exec(`
		INSERT INTO migrations (id, name, source_type, status, worker_job_id, target_db_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, j.ID, j.Name, string(j.SourceType), string(j.Status), nil, j.TargetDBURL, j.CreatedAt, j.UpdatedAt)

	if err != nil {
		fmt.Printf("[migration] DB insert failed: %v, falling back to memory\n", err)
		return fallbackEnqueue(j)
	}
	return &j
}

// ListJobs returns all migration jobs from Postgres.
func ListJobs() []*Job {
	conn := getDB()
	if conn == nil {
		return fallbackList()
	}

	rows, err := conn.Query(`
		SELECT id, name, source_type, status, worker_job_id, created_at, updated_at
		FROM migrations ORDER BY created_at DESC
	`)
	if err != nil {
		fmt.Printf("[migration] DB query failed: %v\n", err)
		return fallbackList()
	}
	defer rows.Close()

	var list []*Job
	for rows.Next() {
		j := &Job{}
		var workerJobID sql.NullString
		var updatedAt sql.NullTime
		if err := rows.Scan(&j.ID, &j.Name, &j.SourceType, &j.Status, &workerJobID, &j.CreatedAt, &updatedAt); err != nil {
			fmt.Printf("[migration] DB scan failed: %v\n", err)
			continue
		}
		if workerJobID.Valid {
			j.WorkerJobID = workerJobID.String
		}
		if updatedAt.Valid {
			j.UpdatedAt = updatedAt.Time
		}
		list = append(list, j)
	}
	return list
}

// GetJob returns a single migration job by ID.
func GetJob(id string) *Job {
	conn := getDB()
	if conn == nil {
		return fallbackGet(id)
	}

	j := &Job{}
	var workerJobID sql.NullString
	var updatedAt sql.NullTime
	err := conn.QueryRow(`
		SELECT id, name, source_type, status, worker_job_id, created_at, updated_at
		FROM migrations WHERE id = $1
	`, id).Scan(&j.ID, &j.Name, &j.SourceType, &j.Status, &workerJobID, &j.CreatedAt, &updatedAt)

	if err != nil {
		return nil
	}
	if workerJobID.Valid {
		j.WorkerJobID = workerJobID.String
	}
	if updatedAt.Valid {
		j.UpdatedAt = updatedAt.Time
	}
	return j
}

// SetWorkerJobID records the BullMQ job ID and marks the job running.
func SetWorkerJobID(id, workerJobID string) {
	conn := getDB()
	if conn == nil {
		fallbackSetWorkerJobID(id, workerJobID)
		return
	}

	_, err := conn.Exec(`
		UPDATE migrations SET worker_job_id = $1, status = 'running', updated_at = now()
		WHERE id = $2
	`, workerJobID, id)
	if err != nil {
		fmt.Printf("[migration] DB update worker_job_id failed: %v\n", err)
	}
}

// MarkFailed records that a job failed.
func MarkFailed(id string) {
	conn := getDB()
	if conn == nil {
		fallbackMarkFailed(id)
		return
	}

	_, err := conn.Exec(`
		UPDATE migrations SET status = 'failed', updated_at = now()
		WHERE id = $1
	`, id)
	if err != nil {
		fmt.Printf("[migration] DB mark failed: %v\n", err)
	}
}

// MarkCompleted records that a job completed successfully.
func MarkCompleted(id string) {
	conn := getDB()
	if conn == nil {
		return
	}

	_, err := conn.Exec(`
		UPDATE migrations SET status = 'completed', updated_at = now()
		WHERE id = $1
	`, id)
	if err != nil {
		fmt.Printf("[migration] DB mark completed: %v\n", err)
	}
}

// DeleteJob removes a migration job by ID.
func DeleteJob(id string) bool {
	conn := getDB()
	if conn == nil {
		return false
	}

	result, err := conn.Exec(`DELETE FROM migrations WHERE id = $1`, id)
	if err != nil {
		return false
	}
	n, _ := result.RowsAffected()
	return n > 0
}

// --- Fallback in-memory store (used when DB is unavailable) ---

var (
	fallbackJobs = map[string]*Job{}
	fallbackSeq  int
)

func fallbackEnqueue(j Job) *Job {
	fallbackSeq++
	j.ID = fmt.Sprintf("mig_%04d", fallbackSeq)
	j.Status = StatusQueued
	j.CreatedAt = time.Now().UTC()
	stored := j
	fallbackJobs[j.ID] = &stored
	return &stored
}

func fallbackList() []*Job {
	list := make([]*Job, 0, len(fallbackJobs))
	for _, j := range fallbackJobs {
		list = append(list, j)
	}
	return list
}

func fallbackGet(id string) *Job {
	return fallbackJobs[id]
}

func fallbackSetWorkerJobID(id, workerJobID string) {
	if j, ok := fallbackJobs[id]; ok {
		j.WorkerJobID = workerJobID
		j.Status = StatusRunning
	}
}

func fallbackMarkFailed(id string) {
	if j, ok := fallbackJobs[id]; ok {
		j.Status = StatusFailed
	}
}
