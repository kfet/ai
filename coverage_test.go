package ai

import (
	"encoding/json"
	"testing"
)

// TestZeroUsage exercises the trivial constructor.
func TestZeroUsage(t *testing.T) {
	if (ZeroUsage() != Usage{}) {
		t.Fatal("ZeroUsage should be the zero value")
	}
}

// TestBoolIntPtr covers the pointer helpers.
func TestBoolIntPtr(t *testing.T) {
	if b := BoolPtr(true); b == nil || *b != true {
		t.Fatal("BoolPtr")
	}
	if i := IntPtr(7); i == nil || *i != 7 {
		t.Fatal("IntPtr")
	}
}

// TestBudgetForLevel covers all branches including nil receiver, the
// default case, and the nil-pointer case.
func TestBudgetForLevel(t *testing.T) {
	var nilTB *ThinkingBudgets
	if nilTB.BudgetForLevel(ThinkingLow) != 0 {
		t.Fatal("nil receiver should be 0")
	}
	tb := &ThinkingBudgets{Minimal: IntPtr(1), Low: IntPtr(2), Medium: IntPtr(3), High: IntPtr(4)}
	cases := map[ThinkingLevel]int{
		ThinkingMinimal: 1,
		ThinkingLow:     2,
		ThinkingMedium:  3,
		ThinkingHigh:    4,
		ThinkingMax:     0, // default branch
	}
	for level, want := range cases {
		if got := tb.BudgetForLevel(level); got != want {
			t.Fatalf("BudgetForLevel(%q) = %d, want %d", level, got, want)
		}
	}
	// nil inner pointer → 0.
	empty := &ThinkingBudgets{}
	if empty.BudgetForLevel(ThinkingLow) != 0 {
		t.Fatal("nil inner pointer should be 0")
	}
}

// TestAssistantContentHelpers covers DeepCopy, ContentType, the Is*
// predicates, the constructors, and MarshalJSON/UnmarshalJSON for every
// branch.
func TestAssistantContentHelpers(t *testing.T) {
	text := NewTextContent("hi")
	think := NewThinkingContent("hmm")
	call := NewToolCallContent("id1", "tool", map[string]any{"a": 1})
	server := NewServerContent("server_tool_use", json.RawMessage(`{"x":1}`), "display")
	var empty AssistantContent

	// ContentType.
	if text.ContentType() != ContentTypeText {
		t.Fatal("text content type")
	}
	if think.ContentType() != ContentTypeThinking {
		t.Fatal("thinking content type")
	}
	if call.ContentType() != ContentTypeToolCall {
		t.Fatal("toolcall content type")
	}
	if server.ContentType() != ContentTypeServer {
		t.Fatal("server content type")
	}
	if empty.ContentType() != "" {
		t.Fatal("empty content type")
	}

	// Is* predicates.
	if !text.IsText() || !think.IsThinking() || !call.IsToolCall() || !server.IsServerContent() {
		t.Fatal("Is* predicates")
	}
	if empty.IsText() || empty.IsThinking() || empty.IsToolCall() || empty.IsServerContent() {
		t.Fatal("empty Is* predicates should be false")
	}

	// DeepCopy for each variant + empty.
	for _, c := range []AssistantContent{text, think, call, server, empty} {
		cp := c.DeepCopy()
		if cp.ContentType() != c.ContentType() {
			t.Fatal("DeepCopy changed content type")
		}
	}
	// DeepCopy of server with nil Raw.
	srvNilRaw := NewServerContent("t", nil, "d")
	if cp := srvNilRaw.DeepCopy(); cp.Server.Raw != nil {
		t.Fatal("DeepCopy nil Raw")
	}

	// Marshal/Unmarshal round-trip for each variant.
	for _, c := range []AssistantContent{text, think, call, server} {
		b, err := json.Marshal(c)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		var back AssistantContent
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if back.ContentType() != c.ContentType() {
			t.Fatalf("round-trip type mismatch: %s vs %s", back.ContentType(), c.ContentType())
		}
	}
	// Marshal empty → null.
	if b, _ := json.Marshal(empty); string(b) != "null" {
		t.Fatalf("empty marshal = %s, want null", b)
	}
	// Unmarshal unknown type → no field set, no error.
	var unknown AssistantContent
	if err := json.Unmarshal([]byte(`{"type":"bogus"}`), &unknown); err != nil {
		t.Fatalf("unmarshal unknown: %v", err)
	}
	if unknown.ContentType() != "" {
		t.Fatal("unknown type should leave content empty")
	}
	// Calling UnmarshalJSON directly with malformed bytes hits the probe
	// decode-error branch (json.Unmarshal would reject the syntax before
	// dispatching to the Unmarshaler).
	if err := unknown.UnmarshalJSON([]byte(`not json`)); err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}

// TestSnapshotContent covers the deep-copy snapshot helper.
func TestSnapshotContent(t *testing.T) {
	m := &AssistantMessage{
		Role:    RoleAssistant,
		Content: []AssistantContent{NewTextContent("a"), NewToolCallContent("i", "n", nil)},
	}
	snap := m.SnapshotContent()
	if len(snap.Content) != 2 {
		t.Fatal("snapshot length")
	}
	// Mutating the original text block must not affect the snapshot.
	m.Content[0].Text.Text = "changed"
	if snap.Content[0].Text.Text != "a" {
		t.Fatal("snapshot not isolated")
	}
}

// TestToolResultContentPredicates covers IsText/IsImage.
func TestToolResultContentPredicates(t *testing.T) {
	txt := &ToolResultContent{Type: ContentTypeText}
	img := &ToolResultContent{Type: ContentTypeImage}
	if !txt.IsText() || txt.IsImage() {
		t.Fatal("text predicates")
	}
	if !img.IsImage() || img.IsText() {
		t.Fatal("image predicates")
	}
}

// TestMessageUnion covers Role, As*, and Marshal/Unmarshal for all
// three message variants plus the empty and error cases.
func TestMessageUnion(t *testing.T) {
	user := NewUserMsg("hello", 1)
	asst := NewAssistantMsg(AssistantMessage{Model: "m"})
	tr := NewToolResultMsg(ToolResultMessage{ToolCallID: "x"})
	var empty Message

	if user.Role() != RoleUser || asst.Role() != RoleAssistant || tr.Role() != RoleToolResult || empty.Role() != "" {
		t.Fatal("Role")
	}
	if user.AsUser() == nil || asst.AsAssistant() == nil || tr.AsToolResult() == nil {
		t.Fatal("As* accessors")
	}
	if user.AsAssistant() != nil || asst.AsUser() != nil {
		t.Fatal("As* mismatched should be nil")
	}

	for _, m := range []Message{user, asst, tr} {
		b, err := json.Marshal(m)
		if err != nil {
			t.Fatalf("marshal: %v", err)
		}
		var back Message
		if err := json.Unmarshal(b, &back); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if back.Role() != m.Role() {
			t.Fatalf("round-trip role mismatch: %s vs %s", back.Role(), m.Role())
		}
	}
	if b, _ := json.Marshal(empty); string(b) != "null" {
		t.Fatalf("empty message marshal = %s, want null", b)
	}
	// Unknown role → no field set.
	var unknown Message
	if err := json.Unmarshal([]byte(`{"role":"bogus"}`), &unknown); err != nil {
		t.Fatalf("unmarshal unknown role: %v", err)
	}
	if unknown.Role() != "" {
		t.Fatal("unknown role should leave message empty")
	}
	// Invalid JSON → error (direct call; see AssistantContent note).
	if err := unknown.UnmarshalJSON([]byte(`nope`)); err == nil {
		t.Fatal("expected error on invalid JSON")
	}
}

// TestModelHelpers covers the server-tool and compat getter helpers.
func TestModelHelpers(t *testing.T) {
	m := &Model{ServerTools: []string{"web_search"}}
	if !m.SupportsServerTool("web_search") || m.SupportsServerTool("nope") {
		t.Fatal("SupportsServerTool")
	}
	if !m.SupportsAnyServerTools() {
		t.Fatal("SupportsAnyServerTools")
	}
	if (&Model{}).SupportsAnyServerTools() {
		t.Fatal("no server tools")
	}

	// Compat getters: matching and non-matching.
	oc := &Model{Compat: &OpenAICompletionsCompat{}}
	if oc.GetOpenAICompletionsCompat() == nil || oc.GetOpenAIResponsesCompat() != nil || oc.GetAnthropicMessagesCompat() != nil {
		t.Fatal("completions compat getter")
	}
	or := &Model{Compat: &OpenAIResponsesCompat{}}
	if or.GetOpenAIResponsesCompat() == nil {
		t.Fatal("responses compat getter")
	}
	// Non-matching compat type → nil (false branch).
	if or.GetOpenAICompletionsCompat() != nil {
		t.Fatal("completions getter should be nil for responses compat")
	}
	am := &Model{Compat: &AnthropicMessagesCompat{}}
	if am.GetAnthropicMessagesCompat() == nil {
		t.Fatal("anthropic compat getter")
	}
}
