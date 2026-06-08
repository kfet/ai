// Package ai defines the portable, model-agnostic primitives for talking
// to large-language-model providers: the message and content types, the
// tool schema, token/cost accounting, the streaming event model, and the
// provider/model shape.
//
// The types are deliberately dependency-free — they describe the wire
// shape of a conversation without prescribing transport, authentication,
// a provider catalog, or an HTTP client. Higher layers (a provider
// registry, an agent runtime) build on these:
//
//   - [Message] is a discriminated union of [UserMessage],
//     [AssistantMessage], and [ToolResultMessage], selected by [Message.Role].
//   - [AssistantContent] is a discriminated union of [TextContent],
//     [ThinkingContent], [ToolCall], and [ServerContent].
//   - [Context] carries the system prompt, the message list, and the tool
//     set for a single request.
//   - [Model] describes a provider model's capabilities, cost, and compat
//     quirks. The concrete catalog of models lives in the consumer.
//   - [StreamFunction] / [SimpleStreamFunction] are the provider streaming
//     entry points; events arrive over an [AssistantMessageEventStream].
//
// Subpackages provide portable, provider-independent helpers:
//
//   - [github.com/kfet/ai/ratelimit] detects rate-limit conditions and
//     parses retry-after hints from assistant error messages.
//   - [github.com/kfet/ai/overflow] detects context-window overflow
//     errors across providers.
//   - [github.com/kfet/ai/jsonparse] best-effort parses partial JSON
//     emitted during tool-call streaming.
//
// Portions are ported from pi-mono (https://github.com/badlogic/pi-mono,
// MIT, Copyright (c) 2025 Mario Zechner); files carry a "// Ported from:"
// header recording the upstream source.
package ai
