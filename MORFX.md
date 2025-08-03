# MORFX CLI — User Guide

> **Version:** 0.1.1 (Go provider only)
> **Binary:** `morfx` (built via `make build` → `bin/morfx`)

---

## 1. Quick Start

```bash
# Run a DSL query against the current directory (Go source)
$ echo "call:util.SHA1FileHex" | morfx -

# Same but from a .dsl file
$ morfx patterns/my_check.dsl
```

The CLI reads **one DSL query** per invocation, translates it to a Tree‑sitter query, executes it over Go files in the working tree, and prints any matches.

---

## 2. Core Commands & Flags

| Flag            | Default | Description                                                                        |
| --------------- | ------- | ---------------------------------------------------------------------------------- |
| `-`             | –       | Read DSL from **STDIN** (dash = stdin convention).                                 |
| `<file>`        | –       | Path to a text file containing the DSL query (first non‑flag arg).                 |
| `--dry-run`     | `false` | Stage changes only (writes to `.morfx/` but does **not** touch source files).      |
| `--commit`      | `false` | Apply all staged changes found in `.morfx/` and delete the staging dir on success. |
| `--summary`     | `false` | After running, print a diff/summary of staged or committed changes.                |
| `--help` / `-h` | –       | Show help.                                                                         |

> **Tip:** `morfx` always exits **non‑zero** if parsing fails or if any staged changes fail to commit (hash mismatch).

---

## 3. Staging vs Commit Workflow

1. **Dry‑run / Stage**

   ```bash
   echo "struct:* > field:Secret string" | morfx - --dry-run --summary
   ```

   _Writes JSON change files under `.morfx/` but leaves source untouched._

2. **Inspect staged diff**

   ```bash
   cat .morfx/change_internal_lang_golang_writer.go.json | jq .
   ```

3. **Commit staged changes**

   ```bash
   morfx --commit --summary
   ```

   _Applies each change atomically with hash verification; skips and reports conflicts; deletes `.morfx/` on full success._

## 4. Writing DSL Queries (Go only)

```text
# exact match
func:Init

# wildcard (prefix)
struct:User*

# negate (skip tests)
!func:Test*

# hierarchy (struct → field)
struct:* > field:Secret string
```

Supported node types: `func`, `const`, `var`, `struct`, `field`, `call`, `assign`, `if`, `import`, `block`.

Wildcards: `Foo*`, `*Foo`, `*Foo*`, `Foo*Bar`

Negation: one per level, `!type:pattern`

Identifier lists: use `any‑` predicates automatically (`var a, b T` handled)

> Full spec: `RULES.md#dsl-v0-1-1` in repo.

---

## 5. What You Can Do

- Search Go codebase with expressive, AST‑aware DSL queries.
- Stage or commit code modifications via Writer abstraction.
- Regenerate and lock golden Tree‑sitter queries (snapshots).
- Run comprehensive gate to ensure DSL compliance (tests, validators, E2E).

---

## 8. Current Limitations

| Area           | Limitation                                                                |
| -------------- | ------------------------------------------------------------------------- |
| Languages      | **Go only**. (PHP/Python providers planned.)                              |
| DSL Logic      | No `AND/OR`, no parentheses. Only single negation per level.              |
| Assign LHS     | Captures **identifiers** only (not `obj.Field` or `arr[idx]`).            |
| Scopes         | No local/global scope predicates yet.                                     |
| Concurrency    | StagingWriter is in‑memory; concurrent runs in same dir may race.         |
| CI Integration | No built‑in Git hooks; run `make gate` manually or add to your CI script. |

---

## 6. Example Session

```bash
# 1) Find insecure md5 usages and stage replacement to sha256
$ echo '!call:crypto.md5* > call:crypto.sha1*' | morfx - --dry-run --summary
Staged 2 change(s) in .morfx/:
  diff path/internal/hash.go …

# 2) Review staged JSON
$ jq '. | {path,operation,original_sha256,modified_sha256}' .morfx/*.json

# 3) Commit if happy
$ morfx --commit --summary
Applied changes to 2 file(s):
  ✓ path/internal/hash.go
```

---

## 7. FAQ (short)

> **Q:** What happens if a file changes after staging but before commit?
> **A:** `CommitWriter` compares the current SHA‑256 with the staged `original_sha256`. If they differ, that file is **skipped** and reported as a conflict.

> **Q:** Can I run morfx on a sub‑directory?
> **A:** Yes, just `cd` into that sub‑folder first; morfx walks recursively from CWD.

> **Q:** Can I combine multiple DSL patterns in one file?
> **A:** Currently morfx processes **one** DSL string per run. Use a shell loop or future playbook support to batch them.

---

_Happy querying!_
