# TFX – Theme System

TFX provides a flexible color theming system designed for terminal-safe visual output, supporting ANSI, 256-color, and TrueColor modes. You can choose from built-in semantic themes or define your own.

---

## 🎨 Built-in Themes

| Theme Name | Preview                         | Notes                      |
| ---------- | ------------------------------- | -------------------------- |
| `Dracula`  | `#ff79c6`, `#bd93f9`, `#f8f8f2` | High-contrast neon palette |
| `Nord`     | `#8fbcbb`, `#88c0d0`, `#5e81ac` | Calm, desaturated colors   |
| `Material` | `#009688`, `#ff9800`, `#e91e63` | Based on Material Design   |
| `Tailwind` | `#0ea5e9`, `#10b981`, `#f43f5e` | Derived from Tailwind CSS  |
| `GitHub`   | `#24292f`, `#0366d6`, `#6f42c1` | Inspired by GitHub UI      |

> Future versions may support light/dark switching and user-defined palettes via CLI.

---

## 🧪 How to Use a Theme

You can set a global theme, or pass one directly into individual components:

### Global Theme

```go
import "github.com/garaekz/tfx/color"

color.SetTheme(color.NordTheme)
```

### Component-specific

```go
progress.Start(
    progress.WithTheme(color.DraculaTheme),
    progress.WithText("Loading"),
)
```

---

## 🧱 Defining Custom Themes

A `Theme` is just a set of named semantic colors. You can create your own:

```go
myTheme := color.Theme{
    Primary:   color.Hex("#00ffd0"),
    Secondary: color.Hex("#ff006e"),
    Accent:    color.Hex("#fefefe"),
}

color.SetTheme(myTheme)
```

Each field can be any `Color` — including RGB, ANSI, or named.

---

## 🖥️ Theme Compatibility

TFX detects your terminal's capabilities and renders as close as possible:

- 24-bit TrueColor terminals: Full fidelity.
- 256-color terminals: Best-match fallback.
- Basic ANSI-only: Graceful degradation.

You can inspect detection using:

```go
fmt.Println(terminal.Capabilities())
```

---

## 🚧 Roadmap

- [ ] Light/dark mode switching
- [ ] User CLI overrides
- [ ] Theme-based spinner/progress styles
- [ ] Runtime theme preview

---

> Themes in TFX are more than color — they’re DX affordances.  
> Code should look as intentional as it feels.
