package gmail

import (
	"encoding/base64"
	"encoding/json"
	"testing"
)

func TestParseMessageFields(t *testing.T) {
	body := base64.URLEncoding.EncodeToString([]byte("hello world"))
	payload, _ := json.Marshal(map[string]any{
		"headers": []any{
			map[string]any{"name": "Subject", "value": "Hi there"},
			map[string]any{"name": "From", "value": "Ada <ada@example.com>"},
		},
		"mimeType": "text/plain",
		"body":     map[string]any{"data": body},
	})
	got := ParseMessageFields(Message{
		InternalDate: "1700000000000",
		Snippet:      "snip",
		Payload:      payload,
	})
	if got.Subject != "Hi there" {
		t.Fatalf("subject=%q", got.Subject)
	}
	if got.From != "Ada <ada@example.com>" {
		t.Fatalf("from=%q", got.From)
	}
	if got.BodyText != "hello world" {
		t.Fatalf("body=%q", got.BodyText)
	}
	if got.OccurredAt == "" {
		t.Fatal("occurredAt empty")
	}
}
