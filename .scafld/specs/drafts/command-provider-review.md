---
spec_version: '2.0'
task_id: command-provider-review
created: '2026-05-04T01:13:01Z'
updated: '2026-05-04T01:13:02Z'
status: completed
harden_status: not_run
size: medium
risk_level: medium
---

# Add command review provider

## Current State

Status: completed
Current phase: none
Next: done
Reason: task completed
Blockers: none
Allowed follow-up command: `none`
Latest runner update: 2026-05-04T01:13:02Z
Review gate: pass

## Summary

Allow scafld review to execute an external command provider and derive pass/fail from parsed review packet frames.

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
  - Command: `go test ./internal/adapters/providers ./internal/adapters/cli ./test/e2e && make check`
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
