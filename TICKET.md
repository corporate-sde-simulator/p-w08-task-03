# PLATFORM-2975: Fix graceful shutdown orchestrator

**Status:** In Progress · **Priority:** Critical
**Sprint:** Sprint 30 · **Story Points:** 5
**Reporter:** Vikram Patel (Infra Lead) · **Assignee:** You (Intern)
**Due:** End of sprint (Friday)
**Labels:** `backend`, `golang`, `production`, `shutdown`
**Task Type:** Bug Fix

---

## Description

The graceful shutdown orchestrator manages cleanup when a service is shutting down — draining connections, flushing caches, and closing resources in the correct order. Two bugs are causing data loss during deployments. Bugs are marked with `// BUG:` comments.

## Acceptance Criteria

- [ ] Bug #1 fixed: Shutdown tasks execute in parallel ignoring declared dependencies (should respect order)
- [ ] Bug #2 fixed: Timeout tracking starts AFTER all tasks complete instead of at shutdown signal
- [ ] All unit tests pass
