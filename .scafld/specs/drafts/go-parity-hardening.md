---
spec_version: '2.0'
task_id: go-parity-hardening
created: '2026-05-04T00:44:18Z'
updated: '2026-05-04T00:54:47Z'
status: completed
harden_status: not_run
size: medium
risk_level: medium
---

# Harden Go runtime parity

## Current State

Status: completed
Current phase: none
Next: done
Reason: task completed
Blockers: none
Allowed follow-up command: `none`
Latest runner update: 2026-05-04T00:54:47Z
Review gate: pass

## Summary

Drive the Go implementation from foundation slice to release-track parity with real scafld workflows, review packets, handoffs, provider streams, and packaging.

## Objectives

- Preserve literate Markdown spec fields across parse and render.
- Project phase completion from recorded session evidence.
- Keep the full Go quality loop green after dogfooding.

## Scope

- Markdown adapter parse/render parity for non-placeholder spec sections.
- Build/reconcile phase status evidence projection.
- Dogfood execution through the Go binary against this repo's own `.scafld` workspace.

## Dependencies

- none

## Assumptions

- This slice should deepen the foundation without trying to replace every Python runtime behavior at once.
- Review packets, provider streams, and packaging remain follow-up parity phases after this dogfood slice.

## Touchpoints

- none

## Risks

- none

## Acceptance

Profile: standard

Validation:
- none

## Phase 1: Implementation

Status: completed
Dependencies: none

Objective: Make the current Go runtime preserve richer Markdown spec fields and project phase completion from evidence.

Changes:
- Extend the Markdown parser to retain current-state labels, list sections, acceptance profile, review state, metadata, origin, and phase dependencies.
- Add a round-trip test for the richer literate spec shape.
- Append phase evidence during build when all criteria in a phase pass, then project that phase status back into the spec.

Acceptance:
- [x] `ac1` command - Markdown adapter preserves richer spec fields
  - Command: `go test ./internal/adapters/markdown -run 'Golden|RoundTrip|Literate'`
  - Expected kind: `exit_code_zero`
  - Status: pass
  - Evidence: exit code was 0
  - Source event: entry-6
- [x] `ac2` command - Phase evidence projects completion
  - Command: `go test ./internal/app/build ./internal/core/reconcile -run 'Phase|Evidence|Projection'`
  - Expected kind: `exit_code_zero`
  - Status: pass
  - Evidence: exit code was 0
  - Source event: entry-7
- [x] `ac3` command - Full repository quality loop remains green
  - Command: `go test ./... && make check`
  - Expected kind: `exit_code_zero`
  - Status: pass
  - Evidence: exit code was 0
  - Source event: entry-8

## Rollback

- none

## Review

Status: completed
Verdict: pass

## Self Eval

- none

## Deviations

- none

## Metadata

- created_by: scafld-go

## Origin

Created by: scafld-go
Source: plan

## Harden Rounds

- none

## Planning Log

- none
