package migration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// workerBaseURL returns the migration worker's HTTP address (its small
// /enqueue shim in worker/src/server.ts). Defaults to the docker-compose
// service name; override with WORKER_URL for local development.
func workerBaseURL() string {
	if url := os.Getenv("WORKER_URL"); url != "" {
		return url
	}
	return "http://localhost:8090"
}

type enqueueRequest struct {
	Name                 string            `json:"name"`
	SourceType           SourceType        `json:"sourceType"`
	SourceConnectionInfo map[string]string `json:"sourceConnectionInfo"`
	TargetDBURL          string            `json:"targetDbUrl"`
}

type enqueueResponse struct {
	JobID string `json:"jobId"`
	Error string `json:"error"`
}

// TriggerWorker asks the Node/BullMQ worker (worker/src/index.ts) to
// actually run a migration job. The Go side only tracks job metadata for
// the dashboard; the worker owns real progress and execution.
func TriggerWorker(ctx context.Context, name string, sourceType SourceType, sourceConnectionInfo map[string]string, targetDBURL string) (string, error) {
	body, err := json.Marshal(enqueueRequest{
		Name:                 name,
		SourceType:           sourceType,
		SourceConnectionInfo: sourceConnectionInfo,
		TargetDBURL:          targetDBURL,
	})
	if err != nil {
		return "", fmt.Errorf("marshal enqueue request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, workerBaseURL()+"/enqueue", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build enqueue request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("worker unreachable at %s: %w", workerBaseURL(), err)
	}
	defer resp.Body.Close()

	var out enqueueResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", fmt.Errorf("decode worker response: %w", err)
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("worker rejected job: %s", out.Error)
	}
	return out.JobID, nil
}
