package ai

// This file collects per-provider compatibility overrides and routing
// preferences (OpenAI-compatible completions/responses quirks,
// Anthropic Messages quirks, OpenRouter and Vercel AI Gateway routing)
// so the core primitives in types.go stay model-agnostic.

// --- OpenAI compatibility settings ---

// ThinkingFormat controls how reasoning/thinking is sent to the provider.
type ThinkingFormat string

const (
	ThinkingFormatOpenAI      ThinkingFormat = "openai"
	ThinkingFormatOpenRouter  ThinkingFormat = "openrouter"
	ThinkingFormatDeepSeek    ThinkingFormat = "deepseek"
	ThinkingFormatZAI         ThinkingFormat = "zai"
	ThinkingFormatQwen        ThinkingFormat = "qwen"
	ThinkingFormatQwenChatTpl ThinkingFormat = "qwen-chat-template"
)

// MaxTokensField controls which JSON field name is used for max tokens.
type MaxTokensField string

const (
	MaxTokensFieldMaxCompletionTokens MaxTokensField = "max_completion_tokens"
	MaxTokensFieldMaxTokens           MaxTokensField = "max_tokens"
)

// OpenAICompletionsCompat holds compatibility overrides for OpenAI-compatible completions APIs.
type OpenAICompletionsCompat struct {
	SupportsStore                               *bool                 `json:"supportsStore,omitempty"`
	SupportsDeveloperRole                       *bool                 `json:"supportsDeveloperRole,omitempty"`
	SupportsReasoningEffort                     *bool                 `json:"supportsReasoningEffort,omitempty"`
	ReasoningEffortMap                          map[string]string     `json:"reasoningEffortMap,omitempty"`
	SupportsUsageInStreaming                    *bool                 `json:"supportsUsageInStreaming,omitempty"`
	MaxTokensField                              MaxTokensField        `json:"maxTokensField,omitempty"`
	RequiresToolResultName                      *bool                 `json:"requiresToolResultName,omitempty"`
	RequiresAssistantAfterToolResult            *bool                 `json:"requiresAssistantAfterToolResult,omitempty"`
	RequiresThinkingAsText                      *bool                 `json:"requiresThinkingAsText,omitempty"`
	RequiresReasoningContentOnAssistantMessages *bool                 `json:"requiresReasoningContentOnAssistantMessages,omitempty"`
	ThinkingFormat                              ThinkingFormat        `json:"thinkingFormat,omitempty"`
	OpenRouterRouting                           *OpenRouterRouting    `json:"openRouterRouting,omitempty"`
	VercelGatewayRouting                        *VercelGatewayRouting `json:"vercelGatewayRouting,omitempty"`
	// ZaiToolStream: whether z.ai supports top-level `tool_stream: true` for
	// streaming tool call deltas. Default: false.
	ZaiToolStream      *bool `json:"zaiToolStream,omitempty"`
	SupportsStrictMode *bool `json:"supportsStrictMode,omitempty"`
	// CacheControlFormat controls prompt caching convention. "anthropic" applies
	// Anthropic-style cache_control markers to system prompt, last tool definition,
	// and last user/assistant text content.
	CacheControlFormat string `json:"cacheControlFormat,omitempty"` // "anthropic" or ""
	// SendSessionAffinityHeaders: whether to send session_id, x-client-request-id,
	// x-session-affinity headers from options.sessionId when caching is enabled.
	SendSessionAffinityHeaders *bool `json:"sendSessionAffinityHeaders,omitempty"`
	// SupportsLongCacheRetention: whether the provider supports long prompt cache
	// retention (prompt_cache_retention: "24h" or Anthropic-style cache_control.ttl: "1h").
	SupportsLongCacheRetention *bool `json:"supportsLongCacheRetention,omitempty"`
}

// OpenAIResponsesCompat holds compatibility overrides for OpenAI Responses APIs.
type OpenAIResponsesCompat struct {
	// SendSessionIdHeader: whether to send the OpenAI session_id cache-affinity header
	// from options.sessionId when caching is enabled. Default: true.
	SendSessionIdHeader *bool `json:"sendSessionIdHeader,omitempty"`
	// SupportsLongCacheRetention: whether the provider supports prompt_cache_retention: "24h". Default: true.
	SupportsLongCacheRetention *bool `json:"supportsLongCacheRetention,omitempty"`
}

// AnthropicMessagesCompat holds compatibility overrides for Anthropic Messages-compatible APIs.
type AnthropicMessagesCompat struct {
	// SupportsEagerToolInputStreaming: whether the provider accepts per-tool eager_input_streaming.
	// When false, the Anthropic provider omits tools[].eager_input_streaming and sends the legacy
	// fine-grained-tool-streaming beta header for tool-enabled requests. Default: true.
	SupportsEagerToolInputStreaming *bool `json:"supportsEagerToolInputStreaming,omitempty"`
	// SupportsLongCacheRetention: whether the provider supports Anthropic long cache retention
	// (cache_control.ttl: "1h"). Default: true.
	SupportsLongCacheRetention *bool `json:"supportsLongCacheRetention,omitempty"`
}

// OpenRouterRouting configures OpenRouter provider routing preferences.
// Sent as the `provider` field in the OpenRouter API request body.
// See https://openrouter.ai/docs/guides/routing/provider-selection
type OpenRouterRouting struct {
	// AllowFallbacks: whether to allow backup providers to serve requests. Default: true.
	AllowFallbacks *bool `json:"allow_fallbacks,omitempty"`
	// RequireParameters: whether to filter providers to only those that support
	// all parameters in the request. Default: false.
	RequireParameters *bool `json:"require_parameters,omitempty"`
	// DataCollection: "allow" (default) or "deny".
	DataCollection string `json:"data_collection,omitempty"`
	// ZDR: restrict routing to only Zero Data Retention endpoints.
	ZDR *bool `json:"zdr,omitempty"`
	// EnforceDistillableText: restrict routing to only models that allow text distillation.
	EnforceDistillableText *bool `json:"enforce_distillable_text,omitempty"`
	// Order: an ordered list of provider names/slugs to try in sequence.
	Order []string `json:"order,omitempty"`
	// Only: providers to exclusively allow for this request.
	Only []string `json:"only,omitempty"`
	// Ignore: providers to skip for this request.
	Ignore []string `json:"ignore,omitempty"`
	// Quantizations: quantization levels to filter providers by.
	Quantizations []string `json:"quantizations,omitempty"`
	// Sort: sorting strategy. Either a string (e.g. "price", "throughput", "latency")
	// or an object with `by` and `partition`. Use any for the union.
	Sort any `json:"sort,omitempty"`
	// MaxPrice: maximum price per million tokens (USD).
	MaxPrice *OpenRouterMaxPrice `json:"max_price,omitempty"`
	// PreferredMinThroughput: preferred minimum throughput (tokens/second).
	// Can be a number or an object with p50/p75/p90/p99 cutoffs.
	PreferredMinThroughput any `json:"preferred_min_throughput,omitempty"`
	// PreferredMaxLatency: preferred maximum latency (seconds).
	// Can be a number or an object with p50/p75/p90/p99 cutoffs.
	PreferredMaxLatency any `json:"preferred_max_latency,omitempty"`
}

// OpenRouterMaxPrice represents maximum price limits per million tokens.
// Values can be numbers or strings (OpenRouter accepts both), so we use any.
type OpenRouterMaxPrice struct {
	Prompt     any `json:"prompt,omitempty"`
	Completion any `json:"completion,omitempty"`
	Image      any `json:"image,omitempty"`
	Audio      any `json:"audio,omitempty"`
	Request    any `json:"request,omitempty"`
}

// VercelGatewayRouting configures Vercel AI Gateway routing preferences.
type VercelGatewayRouting struct {
	Only  []string `json:"only,omitempty"`
	Order []string `json:"order,omitempty"`
}

// GetOpenAICompletionsCompat returns the OpenAI completions compat settings, or nil.
func (m *Model) GetOpenAICompletionsCompat() *OpenAICompletionsCompat {
	if c, ok := m.Compat.(*OpenAICompletionsCompat); ok {
		return c
	}
	return nil
}

// GetOpenAIResponsesCompat returns the OpenAI responses compat settings, or nil.
func (m *Model) GetOpenAIResponsesCompat() *OpenAIResponsesCompat {
	if c, ok := m.Compat.(*OpenAIResponsesCompat); ok {
		return c
	}
	return nil
}

// GetAnthropicMessagesCompat returns the Anthropic messages compat settings, or nil.
func (m *Model) GetAnthropicMessagesCompat() *AnthropicMessagesCompat {
	if c, ok := m.Compat.(*AnthropicMessagesCompat); ok {
		return c
	}
	return nil
}
