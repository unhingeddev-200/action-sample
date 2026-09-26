package gmail

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParsedFields are flat fields extracted from a Gmail message for workflow bindings.
type ParsedFields struct {
	Subject    string
	BodyText   string
	From       string
	OccurredAt string // RFC3339 from internalDate (ms since epoch)
}

// ParseMessageFields extracts subject, plain/html body, From, and occurredAt.
func ParseMessageFields(msg Message) ParsedFields {
	out := ParsedFields{}
	if ms, err := strconv.ParseInt(strings.TrimSpace(msg.InternalDate), 10, 64); err == nil && ms > 0 {
		out.OccurredAt = time.UnixMilli(ms).UTC().Format(time.RFC3339)
	}
	if strings.TrimSpace(msg.Snippet) != "" {
		out.BodyText = strings.TrimSpace(msg.Snippet)
	}
	if len(msg.Payload) == 0 {
		return out
	}
	var payload map[string]any
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		return out
	}
	out.Subject = headerValue(payload, "Subject")
	out.From = headerValue(payload, "From")
	plain, html := collectBodies(payload)
	chosen := plain
	if chosen == "" {
		chosen = html
	}
	chosen = strings.TrimSpace(chosen)
	if chosen != "" {
		out.BodyText = chosen
	}
	return out
}

func headerValue(payload map[string]any, name string) string {
	headers, _ := payload["headers"].([]any)
	for _, h := range headers {
		hm, _ := h.(map[string]any)
		if strings.EqualFold(fmt.Sprint(hm["name"]), name) {
			return strings.TrimSpace(fmt.Sprint(hm["value"]))
		}
	}
	return ""
}

func collectBodies(payload map[string]any) (plain, html string) {
	var walk func(any)
	walk = func(node any) {
		nm, ok := node.(map[string]any)
		if !ok {
			return
		}
		mime, _ := nm["mimeType"].(string)
		body, _ := nm["body"].(map[string]any)
		if b64, _ := body["data"].(string); b64 != "" && strings.HasPrefix(mime, "text/") {
			decoded := decodeBodyData(b64)
			if decoded == "" {
				return
			}
			switch {
			case strings.HasPrefix(mime, "text/plain") && len(decoded) > len(plain):
				plain = decoded
			case strings.HasPrefix(mime, "text/html") && len(decoded) > len(html):
				html = decoded
			}
		}
		if children, _ := nm["parts"].([]any); len(children) > 0 {
			for _, c := range children {
				walk(c)
			}
		}
	}
	walk(payload)
	return plain, html
}

func decodeBodyData(b64 string) string {
	pad := strings.Repeat("=", (4-len(b64)%4)%4)
	if decoded, err := base64.URLEncoding.DecodeString(b64 + pad); err == nil {
		return string(decoded)
	}
	if decoded, err := base64.RawURLEncoding.DecodeString(b64); err == nil {
		return string(decoded)
	}
	return ""
}
