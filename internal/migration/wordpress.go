package migration

import (
  "context"
  "log"
)

type WordPressRunner struct{}

func (w WordPressRunner) Run(ctx context.Context, job Job) error {
  log.Println("Running WordPress migration for job:", job.ID)
  return nil
}
