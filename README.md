# TFX

> Elegant terminal effects & structured output for Go CLIs — composable, fast, and developer-first.

---

## ✨ What is TFX?

**TFX → TermFX** (short for _Terminal Effects_) is a modular Go toolkit for building expressive, styled, and structured terminal output — without the verbosity and fragmentation of typical Go terminal libraries.

It is _not_ a TUI framework. It is a low-friction, highly composable set of tools for CLIs and utilities that need more than just `fmt.Println`, but less than `bubbletea`.

---

## 🤔 When should you use TFX?

| Tool          | Use Case                               |
| ------------- | -------------------------------------- |
| `fmt.Println` | You don’t care about styling/logging   |
| `log/slog`    | You want structured logs, no styling   |
| `TFX`         | You want structure _and_ style ✨      |
| `bubbletea`   | You’re building a full interactive TUI |

---

## 🧠 TFX Principles

- 💎 _Multi-path API_: same engine, multiple ways to use it.
- 🚫 _No reflection, no runtime surprises_.
- 🧰 _Structured logs with color and context_.
- 🧪 _Built for testing_: capture logs, inject writers.
- 🎨 _Themeable by design_: ANSI + semantic palettes.
- 🧱 _Minimal dependencies_: No third-party bloat.

---

## 🔥 Features

- Unified color system with support for:

  - ANSI
  - 256-color
  - TrueColor (24-bit)
  - Semantic palettes & themes (Dracula, Nord, GitHub, Tailwind)

- Structured logging with:

  - Badge-style tags (\[INFO], \[ERR], etc)
  - Contextual fields
  - Multi-writer support (console + file rotation)
  - JSON, text, and badge formats

- Terminal capability detection:

  - Per-OS fallbacks
  - Unicode/ANSI support
  - CI-awareness, `NO_COLOR`, etc

- Progress bars and spinners with smart rendering
- Internal `share/` helpers: `OptionSet`, `Overload` (standardized pattern)

---

## 📦 Packages

| Package           | Description                                               |
| ----------------- | --------------------------------------------------------- |
| `color/`          | Core color system: hex, RGB, ANSI, themes, rendering      |
| `terminal/`       | Terminal detection and capability inference               |
| `logfx/`          | Structured, badge-style logging with writers              |
| `progress/`       | Spinners and progress bars with auto-capability rendering |
| `writers/`        | Console & file writers with rotation and theming          |
| `internal/share/` | Internal DX helpers (option sets, overloads, conventions) |

---

## 🚀 Quickstart

```go
import "github.com/garaekz/tfx/logfx"

func main() {
    logfx.Success("Server started on port %d", 8080)
    logfx.Badge("DB", "Connected to postgres", color.MaterialGreen)
}
```

Use a spinner:

```go
import "github.com/garaekz/tfx/progrefx"

// Construct a spinner with a label using the progrefx package.  The
// StartSpinner helper accepts either a config struct or functional
// options.  Once created, call Tick() to advance frames and Render()
// to obtain the current string representation.
spinner := progrefx.StartSpinner(progrefx.SpinnerConfig{
    Label: "Loading data...",
})

// Simulate some work
for i := 0; i < 10; i++ {
    spinner.Tick()
    fmt.Print(spinner.Render())
    time.Sleep(100 * time.Millisecond)
}
```

---

## 🖼️ Live Preview

See TFX in action with the interactive demo:

```bash
# Build and run the demo
make demo

# Or run specific demonstrations
./bin/demo progress    # Progress bars showcase
./bin/demo color       # Color system showcase
./bin/demo spinner     # Spinners showcase
./bin/demo multipath   # Multipath API showcase
./bin/demo all         # Run all demonstrations
```

---

## 🧱 Philosophy

> "If there's only one way to use it, it's not the right way."

TFX promotes a **multi-entry design** for each API:

```go
logfx.Success("Done!")                            // 1. Quick default
logger := logfx.New(...)
logger.WithFields(...).Info("custom")            // 2. Instantiated style
logfx.If(err).As(logfx.WARN).Msg("warn msg")       // 3. DSL / fluent
```

This consistency is achieved via internal helpers like `Overload()` and `Option[T]`.

---

## 📚 Docs

- [VISION.md](./VISION.md) – project intent and market gap
- [DESIGN_GUIDELINES.md](./DESIGN_GUIDELINES.md) – API conventions and patterns
- [ROADMAP.md](./ROADMAP.md) – current status and future plans
- [MULTIPATH.md](./MULTIPATH.md) – why TermFX APIs support multiple entry paths

---

## 🧪 Status

TFX is under active development. Use at your own risk until `v1.0.0` is tagged.

---

## 📜 License

MIT

---

[![Go Report Card](https://goreportcard.com/badge/github.com/garaekz/tfx)](https://goreportcard.com/report/github.com/garaekz/tfx)

> Built with ☕ and 💢 by [@garaekz](https://github.com/garaekz)
