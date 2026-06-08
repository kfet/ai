# ai

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

Portable, model-agnostic Go primitives for talking to large-language-model
providers: message and content types, the tool schema, token/cost
accounting, the streaming event model, and the provider/model shape.

## Why

Building an LLM client or agent runtime means re-describing the same wire
shapes every time: a conversation is a list of user / assistant / tool-result
messages; an assistant turn is a list of text / thinking / tool-call /
server-content blocks; a model has a cost table, a context window, and a set
of capability quirks. `ai` factors those types out so the runtime above them
need not redefine them, and so provider-independent helpers (rate-limit
detection, context-overflow detection, partial-JSON repair) have a shared
vocabulary.

The package is **dependency-free** at its root (stdlib only). It prescribes
no transport, no auth, no provider catalog, and no HTTP client — those belong
to the consumer.

## Layout

| Package | Purpose |
| --- | --- |
| `ai` (root) | Message/content types, `Tool`, `Usage`, `Context`, `Model`, streaming events, `StreamFunction`, retry classification. |
| `ai/ratelimit` | Detect rate-limit conditions and parse retry-after hints from assistant error messages. |
| `ai/overflow` | Detect context-window overflow errors across providers. |
| `ai/jsonparse` | Best-effort parse of partial JSON emitted during tool-call streaming. |

## Attribution

Portions are ported from [pi-mono](https://github.com/badlogic/pi-mono)
(MIT, Copyright (c) 2025 Mario Zechner). Files derived from that project
carry a `// Ported from:` header. See [LICENSE](LICENSE).

## License

MIT — see [LICENSE](LICENSE).
