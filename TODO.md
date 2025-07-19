# ✅ TFX v0.1 Milestone Checklist

## 🥇 Core Features

- [ ] **Progress Bars**

  - [ ] `NewProgress(max int, label string)`
  - [ ] `.Update(percent int)`
  - [ ] `.Complete(msg string)`
  - [ ] `ProgressStyle` (bar, percent, etc.)

- [ ] **Spinners**

  - [ ] `NewSpinner(label string)`
  - [ ] `.Start()`, `.Stop(msg string)`
  - [ ] Ciclo animado con `\r` y goroutine

- [ ] **Markup DSL (`[bold red]Text[/]`)**

  - [ ] Parser de tags tipo `[color style]text[/]`
  - [ ] Mapeo a colores desde `internal/color`
  - [ ] Función: `lfx.PrintMarkup(...)`

- [ ] **`Dump()` para structs**

  - [ ] `lfx.Dump(obj interface{})`
  - [ ] Usa `json.MarshalIndent` o `spew`
  - [ ] Colorea tipos y claves

- [ ] **Context Enrichment**
  - [ ] `.WithSpanID()`, `.WithService()`
  - [ ] `TraceStart()` / `TraceEnd()` simulados
  - [ ] Inyección en `Fields` automáticamente

## 🧪 Testing Suite Inicial

- [ ] Test para `Success`, `Error`, `Info` con `TestCapture`
- [ ] Test para `Fields` y logger contextual
- [ ] Test para `Progress` y `Spinner` (uso de time, mock de writer)
- [ ] Test de parsing de Markup DSL

## 🎨 Bonus Visual

- [ ] `lfx.PreviewThemes()` imprime 1 línea por theme
- [ ] `lfx.SetFormat(FormatLogfmt)` + formato tipo `key=value`
