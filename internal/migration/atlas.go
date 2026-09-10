// Package migration implements the Atlas Schema Orchestrator: it drives the
// Atlas CLI (https://atlasgo.io) to keep a target Postgres database's schema
// in sync with schema.hcl, and tracks migration jobs run by the worker.
package migration

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

const defaultSchemaPath = "atlas/schema.hcl"

// ApplySchema runs `atlas schema apply`, bringing targetDBURL's schema in
// line with the HCL definition at schemaPath (defaultSchemaPath if empty).
// Requires the Atlas CLI on PATH.
func ApplySchema(ctx context.Context, targetDBURL, schemaPath string) error {
	if targetDBURL == "" {
		return fmt.Errorf("targetDBURL is required")
	}
	if schemaPath == "" {
		schemaPath = defaultSchemaPath
	}

	cmd := exec.CommandContext(ctx, "atlas", "schema", "apply",
		"--url", targetDBURL,
		"--to", fmt.Sprintf("file://%s", schemaPath),
		"--auto-approve",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("atlas schema apply failed: %w: %s", err, stderr.String())
	}
	return nil
}

// DiffSchema returns the pending SQL diff between targetDBURL's current
// state and the desired schema, without applying it.
func DiffSchema(ctx context.Context, targetDBURL, schemaPath string) (string, error) {
	if targetDBURL == "" {
		return "", fmt.Errorf("targetDBURL is required")
	}
	if schemaPath == "" {
		schemaPath = defaultSchemaPath
	}

	cmd := exec.CommandContext(ctx, "atlas", "schema", "diff",
		"--from", targetDBURL,
		"--to", fmt.Sprintf("file://%s", schemaPath),
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("atlas schema diff failed: %w: %s", err, string(out))
	}
	return string(out), nil
}
