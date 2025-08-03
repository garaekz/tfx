# Implementation Plan — `formfx`

> Interactive package for prompt-driven input capture and validation inside `tfx/formfx`, designed to be **fully decoupled** from `flowfx` and `progress`. Uses **TFX multipath API** and **RunFX integration** for robust interactive features.

## Guiding Principles

- **No Cross-Package Dependencies**: `formfx` does not import `flowfx` or `progress`. Integration happens in user code.
- **Multipath API**: Follows TFX patterns with Express/Instantiated/DSL paths.
- **RunFX Integration**: Uses RunFX for robust terminal management and interactive features (arrow keys, WASD navigation).
- **Go-first**: no reflection for public APIs, compile-time safety where possible.
- **Deterministic**: explicit contracts; panics only in internal helpers when contracts are violated.
- **Non-interactive safe**: graceful degradation when TTY is unavailable using `runfx.DetectTTY()`.

## ✅ COMPLETED STAGES

### ✅ Stage 0-1 — Package Scaffolding

- Created `formfx/` directory with proper Go documentation
- Integrated with `tfx/{writer,terminal,color,runfx}` only

### ✅ Stage 2-3 — IO & Terminal Adapters + Error Model

- Terminal detection via `runfx.DetectTTY()`
- Standardized `ErrCanceled` semantics
- TTY-aware fallback for non-interactive environments

### ✅ Stage 4-5 — Multipath API + Prompt Primitives (REFACTORED)

**Replaced functional options with TFX multipath pattern:**

```go
// EXPRESS: Zero-config, quick usage
result, err := formfx.Confirm("Are you sure?")
result, err := formfx.Input("Your name?")
result, err := formfx.Select("Choose:", []string{"A", "B", "C"})

// INSTANTIATED: Config struct for control
cfg := formfx.ConfirmConfig{Label: "Proceed?", Default: true}
result, err := formfx.Confirm(cfg)

// DSL: Chained builder pattern
result, err := formfx.NewConfirm().
    Label("Proceed with deployment?").
    Default(true).
    Interactive(true).
    Show()
```

**Interactive Features:**

- ✅ **Confirm**: Visual yes/no selection with RunFX
- ✅ **Input**: Enhanced text editing capabilities
- ✅ **Select**: Arrow key navigation through options (TODO: Full RunFX implementation)
- ✅ **Secret**: Echo-off password input

### ✅ Stage 6-7 — Validation + Non-TTY Behavior

- Validation helpers via DSL `.Validate()` method
- Automatic non-interactive detection via `runfx.DetectTTY()`
- `FORM_NONINTERACTIVE=1` environment override

## 🚀 NEXT STEPS (Future Implementation)

### Stage 8 — Full RunFX Interactive Integration

**TODO: Complete RunFX integration for enhanced interactive features**

```go
// Advanced Select with arrow key navigation
result, err := formfx.NewSelect().
    Label("Choose deployment target:").
    Options([]string{"Production", "Staging", "Development"}).
    Interactive(true).  // Enables full RunFX visual selection
    PageSize(5).
    Show()

// Enhanced Input with history and autocomplete
result, err := formfx.NewInput().
    Label("Enter command:").
    Interactive(true).  // Enables advanced editing features
    Show()
```

**Implementation Plan:**

1. Create RunFX Visual implementations for each prompt type
2. Add keyboard handling (arrow keys, WASD, Enter, Escape)
3. Visual highlighting and cursor management
4. Graceful fallback to simple mode in non-TTY environments

### Stage 9 — Examples & Documentation

- Comprehensive examples for all three API styles
- Integration examples with FlowFX (in user code)
- Non-interactive usage patterns

### Stage 10 — Advanced Features

- Multi-select support (`MultiSelect`)
- Input validation with real-time feedback
- Autocomplete and history for input fields
- Custom themes and styling

---

## 🏆 CURRENT STATUS

**✅ COMPLETE:**

- Multipath API implementation (Express/Instantiated/DSL)
- Basic interactive prompts (Confirm, Input, Select, Secret)
- RunFX integration foundation with TTY detection
- Non-interactive fallback behavior
- Validation support via DSL

**🚧 IN PROGRESS:**

- Full RunFX Visual implementations for enhanced interactivity
- Arrow key navigation for Select prompts
- Advanced editing features for Input prompts

**📋 REMAINING:**

- Complete RunFX Visual integration
- Advanced keyboard handling
- Comprehensive examples and documentation

---

## Global Acceptance Criteria

- ✅ All prompts: `(value, error)`; cancel → `ErrCanceled`.
- ✅ No cross-package imports with `flowfx` or `progress`.
- ✅ Works on TTY and non-TTY; Windows/\*nix parity for echo control.
- ✅ TFX multipath API compliance (Express/Instantiated/DSL).
- ✅ RunFX integration for robust terminal management.
- 🚧 Unit + golden tests; CI green; Go Report Card A+.

---

## Stage 3 — Error Model & Cancel Semantics

**Tasks**

- Define `var ErrCanceled = errors.New("form: canceled")`.
- Map ESC/Ctrl+C/EOF to `ErrCanceled`.
- Add `CancelPolicy` option.

**Deliverables**

- `errors.go`
- Cancel handling integrated in prompt loop.

---

## Stage 4 — Config & Options (Functional Options)

**Tasks**

- Introduce `Option[T]` + `ApplyOptions` (reuse from `share` if available).
- For each prompt, define a typed `Config` and `DefaultConfig()`.
- Provide helpers like `WithDefault(v)`, `WithValidate(fn)`, `WithPageSize(n)`.

**Deliverables**

- `options.go`
- `*_config.go` for each prompt.

---

## Stage 5 — Prompt Primitives

**Tasks**

1. **Confirm**: `(bool, error)`
2. **Input**: `(string, error)`
3. **Secret**: `(string, error)` with echo-off
4. **Select**: `(int, error)` index-based selection

- All prompts block, return on success/cancel/error.
- No dependencies outside `tfx/{writer,terminal,color}`.

**Deliverables**

- `confirm.go`, `input.go`, `secret.go`, `select.go`

---

## Stage 6 — Validation Helpers

**Tasks**

- Subpackage `formfx/validate` with reusable validators.

**Deliverables**

- `validate/validate.go`

---

## Stage 7 — Non‑TTY & Piped Input Behavior

**Tasks**

- Provide non-interactive fallbacks.
- Environment override: `FORM_NONINTERACTIVE=1`.

**Deliverables**

- `non_tty.go`
- `NONINTERACTIVE.md`

---

## Stage 8 — Examples & Docs

**Tasks**

- `examples/` for each primitive.
- Update package comments.

**Deliverables**

- Example programs.

---

## Stage 9 — Tests & Golden Files

**Tasks**

- Unit tests for config merging, validators, cancel paths.
- Golden tests for rendered output.

**Deliverables**

- `*_test.go` files, `testdata/` fixtures.

---

## Stage 10 — API Stability

**Tasks**

- No breaking changes to public function signatures.
- All integration with `flowfx` or `tfx/progress` must happen in **user code**.

**Deliverables**

- `API_POLICY.md`

---

## Global Acceptance Criteria

- All prompts: `(value, error)`; cancel → `ErrCanceled`.
- No cross-package imports with `flowfx` or `tfx/progress`.
- Works on TTY and non-TTY; Windows/\*nix parity for echo control.
- Unit + golden tests; CI green; Go Report Card A+.
