package ai_test

// These examples are compile-checked usage documentation for the headline
// message API (visible on pkg.go.dev). They use no provider, no transport,
// and no HTTP — just the portable types.

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kfet/ai"
)

// Example builds a tiny conversation and walks it, discriminating each
// message by its role. A Message is a tagged union — Role reports the
// variant and the As* accessors return the concrete payload (or nil).
func Example() {
	convo := []ai.Message{
		ai.NewUserMsg("What is 2+2?", 0),
		ai.NewAssistantMsg(ai.AssistantMessage{
			Content: []ai.AssistantContent{ai.NewTextContent("4")},
		}),
	}

	for _, m := range convo {
		switch m.Role() {
		case ai.RoleUser:
			fmt.Printf("user: %v\n", m.AsUser().Content)
		case ai.RoleAssistant:
			var b strings.Builder
			for _, c := range m.AsAssistant().Content {
				if c.Text != nil {
					b.WriteString(c.Text.Text)
				}
			}
			fmt.Printf("assistant: %s\n", b.String())
		}
	}
	// Output:
	// user: What is 2+2?
	// assistant: 4
}

// Example_jsonRoundTrip shows that a Message survives a JSON round-trip:
// MarshalJSON serialises the active variant and UnmarshalJSON re-selects it
// from the "role" field, so the discriminated content blocks come back
// intact.
func Example_jsonRoundTrip() {
	orig := ai.NewAssistantMsg(ai.AssistantMessage{
		Content: []ai.AssistantContent{
			ai.NewTextContent("hi"),
			ai.NewToolCallContent("call_1", "get_weather", map[string]any{"city": "SF"}),
		},
	})

	data, err := json.Marshal(orig)
	if err != nil {
		fmt.Println("marshal:", err)
		return
	}

	var back ai.Message
	if err := json.Unmarshal(data, &back); err != nil {
		fmt.Println("unmarshal:", err)
		return
	}

	am := back.AsAssistant()
	fmt.Println(back.Role())
	fmt.Println(am.Content[0].Text.Text)
	fmt.Println(am.Content[1].ToolCall.Name)
	// Output:
	// assistant
	// hi
	// get_weather
}
