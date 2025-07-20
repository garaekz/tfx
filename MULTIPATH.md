# TFX – Multipath API Philosophy

> “If there’s only one way to do it, it’s not the right way.”  
> — Core Principle of TFX

---

## 1. 🤔 Why Multipath?

Go encourages explicitness, but often ends up verbose — nudging library authors toward rigid **single-path APIs**: one constructor, one pattern, zero freedom.

**TFX** embraces a multi-entry design, offering **three coherent entry points** for every public-facing feature:

| Entry                 | Goal                                   | Typical Consumer               |
| --------------------- | -------------------------------------- | ------------------------------ |
| **Default / Express** | Zero-config, one-liner utility         | Scripts, demos                 |
| **Instantiated**      | Fine-grained control via config struct | Library authors, tool builders |
| **Fluent / DSL**      | Declarative chaining with readability  | Power users, prototyping       |

All styles share the same internal engine — no duplicated logic, no hidden side effects.

---

## 2. 🧰 Internal Consistency via `internal/share`

TFX provides two DX primitives to enable clean multipath APIs:

### 2.1 `Overload[T]`

```go
cfg := share.Overload[Config](args, DefaultConfig())
```

- Accepts **0 or 1** anonymous value.
- Handles both `Config` and `*Config`.
- Falls back safely to the provided default.

### 2.2 `ApplyOptions[T]` (Functional Option Set)

```go
share.ApplyOptions(&cfg,
    progress.WithText("Downloading"),
    progress.WithColor(color.Cyan),
)
```

- Avoids long parameter lists.
- Zero reflection — 100% compile-time safety.
- Uses Go generics for reusability.

### 2.3 `OverloadWithOptions[T]`

```go
cfg := share.OverloadWithOptions[Config](args, DefaultConfig(), userOpts...)
```

Combines positional + keyed options in one call.

---

## 3. 📐 Canonical Pattern per Package

Every TFX package adheres to this baseline:

```go
// 1. Quick default
progress.Start()

// 2. Config struct
cfg := progress.Config{Text: "Sync", Color: color.Green}
progress.Start(cfg)

// 3. Fluent builder
progress.Start(
    progress.WithText("Sync"),
    progress.WithColor(color.Green),
)

// 4. Full object lifecycle
bar := progress.New(progress.WithWidth(40))
bar.Update(75)
bar.Complete()
```

### ✅ Enforcement Checklist

- `Start()` must work with **zero args**.
- An instantiated form (`Start(cfg)` or `New(...)`) must exist.
- All `WithX` options must live in the **same package**.
- If `...any` is exposed, provide a **typed sibling** like `StartWith(cfg)` for IDE safety.

---

## 4. 🎯 Design Rationale

1. **DX First** — Libraries are judged by the first 30 seconds.
2. **No Reflection, No Magic** — Compile-time safety over runtime tricks.
3. **Opt-in Complexity** — Keep simple things simple. Expose power without requiring it.
4. **Uniform Mental Model** — Learn one, use everywhere.

---

## 5. 🧪 How to Build a New TFX Package

1. Define a `Config` with sensible defaults.
2. Provide `DefaultConfig()`.
3. Create `WithX`-style options:

   ```go
   func WithText(text string) share.Option[Config] {
       return func(c *Config) { c.Text = text }
   }
   ```

4. Expose all three entry points: `Start()`, `Start(cfg)`, `Start(opts...)`.
5. Use `share.OverloadWithOptions` internally.
6. Document all usage styles via GoDoc.

---

## 6. ❓ FAQ

**Q:** Isn’t `...any` evil?  
**A:** It’s scoped and type-switched. Regular users stay in safe APIs.

**Q:** Isn’t this over-engineered?  
**A:** Not if it lets beginners and power users coexist without friction.

**Q:** Why not use reflection?  
**A:** It’s untestable and runtime-only. TFX stays fast, lean, and predictable.

---

## 🧠 Summary

The **Multipath API** is not a gimmick — it's a **DX multiplier**.  
It lets TFX be:

- Beginner-friendly in a one-liner.
- Scalable for complex CLI tools.
- Consistent and safe by design.

> Build APIs that scale **horizontally**, not just vertically.
