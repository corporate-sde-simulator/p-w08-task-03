# PR Review - Graceful shutdown orchestrator (by Suresh)

## Reviewer: Vikram Patel
---

**Overall:** Good foundation but critical bugs need fixing before merge.

### `shutdownOrchestrator.go`

> **Bug #1:** Shutdown hooks execute in registration order instead of reverse order so deps are closed too early
> This is the higher priority fix. Check the logic carefully and compare against the design doc.

### `resourceManager.go`

> **Bug #2:** Timeout handler fires even when all hooks complete successfully and always logs timeout error
> This is more subtle but will cause issues in production. Make sure to add a test case for this.

---

**Suresh**
> Acknowledged. I have documented the issues for whoever picks this up.
