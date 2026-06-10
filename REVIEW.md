# OSS readiness review — `github.com/kfet/ai`

Independent review of this repo as a **standalone open-source Go library**
(not as the in-tree fir package it was extracted from). Build is green:
`make all` passes and the module is at a 100% coverage gate.

**Verdict: ship as v0.1.0.** No correctness blockers. Everything below is
P2/P3 polish. Address the findings, keep `make all` green, then it's tag-ready.

---

## Structure — idiomatic, keep as-is

Root package `ai` owns the headline types; `ratelimit/`, `overflow/`,
`jsonparse/` subpackages are cohesive and each earns its place. The
discriminated unions (`types.go:Message`, `AssistantContent`) are done
properly: unexported variant pointers, `Role()`/`As*()` accessors, custom
`MarshalJSON`/`UnmarshalJSON` round-tripping through the tag field. Zero-value
usable. Tests genuinely hit 100%.

`Api = string`, `Provider = string`, `InputModality = string` as type *aliases*
(not defined types) is the **right** call for a wire-shape lib — keep them.

## Findings & required changes

### P2 — Scrub fir-isms / dangling internal references from public doc comments
These leak the extraction onto pkg.go.dev:
- `types.go:280` — `ServerContent` doc: "that **fir** does not interpret
  semantically" → reword generically (e.g. "that this package does not
  interpret semantically").
- `types.go:284` — "see **BACKLOG.md** for the motivation" → BACKLOG.md does
  not exist in this repo. Remove or inline the motivation.
- `types.go:286` — "// **fir**-internal discriminator" → "// internal
  discriminator".
- `types.go:865` — "narrower than **fir's** {minimal,…}" → reword.
- `ratelimit/ratelimit.go:51,59,67` — comments reference the internal path
  "pkg/ai/ai" → should read `github.com/kfet/ai`.
- `retryclass.go:7-8` — comments reference "fir-side", "pkg/ai", "pkg/agent" →
  reword for a standalone lib.

### P2 — Go-initialism naming (coordinate cross-repo, see below)
`Api` → `API`, `ApiKey` / `StreamOptions.ApiKey` → `APIKey`. Pre-1.0 is the
only free window to fix this. **This is a coordinated cross-repo rename** — see
the "Cross-repo coordination" section; do not ship it half-done.

### P3 — Split provider compat structs out of `types.go`
`types.go` is ~870 lines and mixes the core primitives with very
provider-specific compat structs (`OpenRouterRouting`, `VercelGatewayRouting`,
`OpenAICompletionsCompat`, z.ai `ZaiToolStream`, Cerebras quirks). Move the
compat structs into a new `compat.go` (same package, zero API change) to sharpen
the "model-agnostic primitives" story.

### P3 — Drop or document `UserContentBlock = any`
`types.go:UserContentBlock = any` is an alias to `any` that buys nothing. Drop
it or replace with a documented union.

---

## Cross-repo coordination — the `Api`/`ApiKey` rename

The initialism rename threads between `kfet/ai` and `kfet/agent`
(`StreamOptions.ApiKey` etc. are consumed by the agent runtime). For **this
repo**, the rename is self-contained — perform it fully here:
`Api`→`API`, `ApiKey`→`APIKey`, update all call sites, doc comments, and tests,
keep `make all` green. The sibling `kfet/agent` repo will pick up the renamed
symbols in a separate coordinated dependency bump (it currently pins
`kfet/ai v0.0.1` with the old names, so it stays building against the old tag
until this work is released).
