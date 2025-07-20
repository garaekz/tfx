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
  - Multi-path examples
- Keep usage examples runnable.

---

> Consistency beats cleverness.  
> TFX should feel familiar after 5 minutes — and powerful after 5 days.
