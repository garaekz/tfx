# TFX – Design Guidelines

This document outlines the design principles, conventions, and architectural patterns that guide TFX development.
It ensures a consistent developer experience and a predictable API surface across all internal and public modules.

---

## ✅ General Principles

- **DX First** – Favor ergonomics and fast feedback over premature abstraction.
- **Zero Reflection** – Use generics and types instead of unsafe reflection or interface{}.
- **Predictable APIs** – Every feature should expose a `Default`, `Config`, and `Fluent` interface.
- **Go Native** – Stick to idiomatic Go whenever possible.
- **Minimal API Surface** – Expose only what’s useful. Internal helpers go in `internal/share`.
- **Separation of Concerns** – Multipath flexibility belongs to express functions, not constructors.

---

## 📐 API Shapes

### 1. Functional Options

Every public-facing feature must support `WithX` options:

```go
func WithText(text string) share.Option[Config] {
    return func(c *Config) { c.Text = text }
}
```

- Options must live in the same package (`progress.WithText`, not `progressopts.WithText`).
- Avoid nested Option structs.

### 2. Overload + Fallback

Always allow:

- Zero-arg express usage: `Start()`
- Config struct usage: `Start(cfg)`
- Fluent chaining usage: `Start(WithX, WithY...)`

Use `share.OverloadWithOptions` internally to keep all paths unified.

**Typed sibling required:** If you expose `Start(...any)`, you **must also** expose `StartWith(cfg)` to retain type-safe DX and IDE discoverability.

---

## 🧬 Entry Point Enforcement

| Entry Type     | Function Example              | Rules                                                                 |
| -------------- | ----------------------------- | --------------------------------------------------------------------- |
| Express        | `Start()`                     | Uses multipath (`...any`) — must support config, options, or nothing. |
| Object Creator | `New()`, `NewWithConfig(cfg)` | Must be strongly typed. **Never** use `...any` in constructors.       |
| IDE Support    | `StartWith(cfg)`              | Mandatory if `Start(...any)` exists. Must be GoDoc-visible.           |

### Example

```go
Start() // default
Start(WithTotal(50), WithLabel("Sync")) // fluent
StartWith(Config{Total: 50, Label: "Sync"}) // typed

bar := New()
bar.Set(10)
bar.Complete("Done!")

bar2 := NewWithConfig(Config{Total: 100})
bar2.Start()
```

---

## 🧱 File & Package Structure

- Avoid `utils/` or `helpers/` — name by responsibility (`color`, `logx`, `writers`, `progress`, etc).
- Internal helpers must live in `internal/share/`.
- Functional options go in the same package as the consumer.
- Keep public packages flat; no deep nesting.

---

## 🧪 Testing Standards

- Always test with race detector: `go test -race ./...`
- Use `TestCaptureWriter` or mocks for log assertions.
- Use `go tool cover -html=...` and maintain high coverage for core packages.
- Validate all usage paths: express, config, and fluent APIs.

---

## 📦 Versioning & Compatibility

- No breaking changes before `v1.0.0`.
- All public APIs must be reviewed for:

  - Overload safety
  - Zero-config usability
  - Composability

- Mark anything internal or unstable with a comment: `// EXPERIMENTAL` or `// INTERNAL ONLY`

---

## 📚 Documentation Rules

- Each package must include:

  - A `Config` example in its GoDoc
  - Functional Option examples
  - Multi-path usage examples (Start, StartWith, WithX...)

- Keep usage examples runnable.
- Reference [MULTIPATH.md](./MULTIPATH.md) for philosophical alignment.

---

> Consistency beats cleverness.
> TFX should feel familiar after 5 minutes — and powerful after 5 days.
