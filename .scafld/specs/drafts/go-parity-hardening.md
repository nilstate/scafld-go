---
spec_version: '2.0'
task_id: go-parity-hardening
created: '2026-05-04T00:44:18Z'
updated: '2026-05-04T00:44:18Z'
status: draft
harden_status: not_run
size: medium
risk_level: medium
---

# Harden Go runtime parity

## Current State

Status: draft
Current phase: none
Next: approve
Reason: draft created
Blockers: none
Allowed follow-up command: `scafld approve go-parity-hardening`
Latest runner update: none
Review gate: not_started

## Summary

Drive the Go implementation from foundation slice to release-track parity with real scafld workflows, review packets, handoffs, provider streams, and packaging.

## Objectives

- none

## Scope

- none

## Dependencies

- none

## Assumptions

- none

## Acceptance

Profile: standard

Validation:
- none

## Phase 1: Implementation

Status: pending
Dependencies: none

Objective: Complete the requested change.

Changes:
- Implement the requested behavior.

Acceptance:
- [ ] `ac1` command - Primary validation command
  - Command: `go test ./... && make check`
  - Expected kind: `exit_code_zero`
  - Status: pending

## Rollback

- none

## Review

Status: not_started
Verdict: none

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
