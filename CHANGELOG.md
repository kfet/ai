# Changelog

All notable changes to this project are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.2] - 2026-06-10

### Added

- Tool-result meta channel: `ToolResultMessage.Meta` (`map[string]string`,
  `json:"meta,omitempty"`) plus `RenderToolResultMeta`, which renders meta
  deterministically (sorted keys, `[key: value]` lines) for the LLM-bound
  message copy. Internal consumers that join content blocks never see Meta;
  the rendered bytes are stable across turns so provider prompt caches stay
  valid. Additive, non-breaking.

## [0.1.1] - 2026-06-10

### Added

- `Model.AdaptiveThinking` (`json:"adaptiveThinking,omitempty"`): marks models
  with always-on, effort-based adaptive thinking that cannot disable thinking
  or set an explicit token budget. Additive, non-breaking.

## [0.1.0] - 2026-06-10

### Changed

- **Breaking:** renamed `Api` → `API` and `ApiKey` → `APIKey` per Go
  initialism convention: the `API` type alias, all `API*` constants
  (`APIOpenAICompletions`, `APIAnthropicMessages`, …), `Model.API`,
  `StreamOptions.APIKey`, and `StreamOptions.RefreshAPIKey`. JSON wire
  format is unchanged.
- **Breaking:** dropped the `UserContentBlock = any` alias; the
  `UserMessage.Content` union (string, or `[]any` of `TextContent` |
  `ImageContent`) is now documented directly on the struct.
- Moved provider compatibility structs (`ThinkingFormat`,
  `MaxTokensField`, `OpenAICompletionsCompat`, `OpenAIResponsesCompat`,
  `AnthropicMessagesCompat`, `OpenRouterRouting`, `OpenRouterMaxPrice`,
  `VercelGatewayRouting`, and the `Model.Get*Compat` accessors) from
  `types.go` into a new `compat.go`. Same package; no API change.
- Scrubbed internal references from public doc comments (fir-isms,
  `BACKLOG.md`, `pkg/ai`/`pkg/agent` paths) for the standalone library.

### Added

- Initial extraction from fir (`github.com/kfet/fir`) Phase 5. Portable
  AI primitives: message/content types, `Tool`, `Usage`, `Context`,
  `Model`, streaming events, `StreamFunction`, and retry classification
  in the root `ai` package; `ratelimit/`, `overflow/`, and `jsonparse/`
  subpackages. Portions ported from pi-mono (MIT, Copyright (c) 2025
  Mario Zechner).
