package migration

import "context"

type Job struct {
  ID     int64
  Source string
}

type Runner interface {
  Run(ctx context.Context, job Job) error
}
