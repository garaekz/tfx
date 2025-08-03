# FlowFX — Comprehensive Implementation Plan

> **FlowFX** is a robust library designed to manage structured, deterministic, and composable CLI workflows. It provides sequential, parallel, branching, wizard-style, scripted, and hierarchical (tree) execution models. FlowFX focuses on clear, explicit, and readable flow definitions, predictable outcomes, and seamless integration with external systems without importing any internal packages except foundational utilities.

---

## Guiding Principles

- **Explicit Composition:** Users explicitly compose flows; no implicit dependencies.
- **Stateless Execution:** Flow state remains contained and explicit.
- **Zero Cross-Imports:** No imports from other ecosystem packages (e.g., `formfx`, `runfx`, `progress`). Integration occurs exclusively through user-defined interfaces. Only foundational utilities (`terminal`, `writer`, `color`) are permissible imports.
- **Predictable and Safe:** Structured error handling, retries, timeouts, and cancellation.
- **Non-Interactive Ready:** Safe behavior in non-TTY environments.

---

## Stage 0 — Project Setup & Documentation

**Tasks:**

- Define project scope clearly in `SCOPE.md`.
- Setup directory structure (`termfx/flowfx/`).
- Add base documentation (`doc.go`) and interface stubs.

**Deliverables:**

- Initial scaffold with documentation (`doc.go`).

---

## Stage 1 — Core Interfaces & Error Definitions

**Tasks:**

- Define core interfaces clearly:

```go
type Step interface {
    Execute(ctx context.Context) error
}

type Flow interface {
    Run(ctx context.Context) error
}
```

- Common error definitions (`ErrCanceled`, `ErrTimeout`, `ErrRetryExhausted`).

**Deliverables:**

- `interfaces.go`, `errors.go`.

---

## Stage 2 — Sequential Execution (`Sequence`)

**Tasks:**

- Implement sequential flow execution.
- Support context-aware task cancellation, retries, and timeouts.
- Hooks: `OnStart`, `OnComplete`, `OnError`.

**Deliverables:**

- `sequence.go`.

---

## Stage 3 — Parallel Execution (`Parallel`)

**Tasks:**

- Concurrent execution using `sync.WaitGroup`.
- Individual task management (timeouts, retries, error isolation).
- Collect and report errors clearly.

**Deliverables:**

- `parallel.go`.

---

## Stage 4 — Task Abstraction (`Task`)

**Tasks:**

- Struct definition with clear labels, execution functions, retries, and timeouts.
- User-injectable `ProgressReporter` interface.

```go
type Task struct {
    Label    string
    Run      func(ctx context.Context) error
    Retry    int
    Timeout  time.Duration
    Reporter ProgressReporter
}
```

**Deliverables:**

- `task.go`, `progress.go` (interface only, no imports).

---

## Stage 5 — Branching & Conditional Logic (`Branch`)

**Tasks:**

- Implement conditional (`When`) logic execution.
- Allow nested flow execution (sequence, parallel, or wizard).

**Deliverables:**

- `branch.go`.

---

## Stage 6 — Wizard Flow (`Wizard`)

**Tasks:**

- Implement linear, step-by-step execution.
- Explicit structured input/output (`map[string]any`).
- Integration points via user-supplied prompt functions (no direct imports).

**Deliverables:**

- `wizard.go`.

---

## Stage 7 — Structured Configuration (`MapFlow`)

**Tasks:**

- Implement key-value structured configuration flows.
- Serialize and deserialize flow configurations to/from JSON and YAML.

**Deliverables:**

- `mapflow.go`.

---

## Stage 8 — Scripted Execution (`Script`)

**Tasks:**

- Provide sequential task chaining similar to scripting.
- Enhanced logging and error traceability.

**Deliverables:**

- `script.go`.

---

## Stage 9 — Hierarchical Execution (`Tree`)

**Tasks:**

- Implement tree-structured execution for hierarchical tasks.
- Clear hierarchical visual/logging representation.

**Deliverables:**

- `tree.go`.

---

## Stage 10 — Flow Runner & Execution Control

**Tasks:**

- Define explicit `Runner` abstraction.
- Manage continuous execution and responsive cancellation.
- Ensure compatibility with system interrupts (`Ctrl+C`).

```go
type Runner interface {
    Start(ctx context.Context) error
    Stop() error
}
```

**Deliverables:**

- `runner.go`, interrupt handling.

---

## Stage 11 — Non-Interactive Execution & Overrides

**Tasks:**

- Clear and documented behavior for non-TTY execution.
- Environment-variable overrides (`FLOWFX_NONINTERACTIVE`).

**Deliverables:**

- Non-interactive execution logic (`noninteractive.go`).

---

## Stage 12 — Observability & Enhanced Logging

**Tasks:**

- Structured, detailed logging for debugging and auditing.
- Optional debug mode with verbose flow state tracing.

**Deliverables:**

- `logging.go`.

---

## Stage 13 — Comprehensive Documentation & Examples

**Tasks:**

- Examples demonstrating each flow type clearly.
- Package-level documentation ready for `pkg.go.dev`.

**Deliverables:**

- Example directory (`examples/sequence`, `examples/parallel`, `examples/tree`, etc.).

---

## Stage 14 — Robust Testing Strategy

**Tasks:**

- Detailed unit tests covering all flows and conditions.
- Integration tests combining multiple flows.
- CI configuration with `golangci-lint`, coverage, race detection.

**Deliverables:**

- Comprehensive test suites (`*_test.go`).
- CI pipeline configuration.

---

## Stage 15 — API Stability & Error Clarity

**Tasks:**

- Finalize public APIs for semantic versioning.
- Comprehensive error documentation and handling strategies.

**Deliverables:**

- Documented API guarantees.

---

## Stage 16 — Release & Versioning

**Tasks:**

- Prepare and tag stable release (`v0.x`).
- Update changelog, release notes, and migration paths clearly.

**Deliverables:**

- Tagged stable release with clear release documentation.

---

## Acceptance Criteria

- Reliable deterministic execution of all flow types.
- Clear, structured input/output handling.
- Explicit composition without hidden imports.
- Comprehensive documentation and runnable examples.
- Robust test coverage with detailed CI checks.
- Safe and predictable behavior in non-interactive environments.
