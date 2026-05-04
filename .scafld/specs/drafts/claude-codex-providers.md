---
spec_version: '2.0'
task_id: claude-codex-providers
created: '2026-05-04T02:13:52Z'
updated: '2026-05-04T02:13:53Z'
status: completed
harden_status: not_run
size: medium
risk_level: medium
---

# Add Claude and Codex review providers

## Current State

Status: completed
Current phase: none
Next: done
Reason: task completed
Blockers: none
Allowed follow-up command: `none`
Latest runner update: 2026-05-04T02:13:53Z
Review gate: pass

## Summary

Implement provider-specific review adapters with exact Claude stream-json and Codex read-only argv contracts, provider output extraction, and CLI selection.

## Objectives

- none

## Scope

- none

## Dependencies

- none

## Assumptions

- none

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

Objective: Complete the requested change.

Changes:
- Implement the requested behavior.

Acceptance:
- [x] `ac1` command - Primary validation command
  - Command: `go test ./internal/adapters/providers ./internal/adapters/process ./internal/core/review ./internal/adapters/cli ./test/e2e && make check`
  - Expected kind: `exit_code_zero`
  - Status: pass
  - Evidence: exit code was 0
  - Source event: entry-2

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
