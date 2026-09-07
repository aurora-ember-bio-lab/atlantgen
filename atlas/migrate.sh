#!/usr/bin/env bash
set -e

ATLAS_URL="${ATLAS_URL:-postgres://user:pass@localhost:5432/aura}"

atlas schema apply \
  -u "$ATLAS_URL" \
  --to "file://schema.hcl" \
  --auto-approve
