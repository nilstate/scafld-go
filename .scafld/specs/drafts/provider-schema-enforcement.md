---
spec_version: '2.0'
task_id: provider-schema-enforcement
created: '2026-05-04T02:33:21Z'
updated: '2026-05-04T02:33:22Z'
status: completed
harden_status: not_run
size: medium
risk_level: medium
---

# Attach provider review schemas

## Current State

Status: completed
Current phase: none
Next: done
Reason: task completed
Blockers: none
Allowed follow-up command: `none`
Latest runner update: 2026-05-04T02:33:22Z
Review gate: pass

## Summary

Attach a default ReviewPacket JSON schema to Claude and Codex provider runs so real model reviews emit structured packets instead of best-effort prose.

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
  - Command: `go test ./internal/adapters/providers && make check`
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
