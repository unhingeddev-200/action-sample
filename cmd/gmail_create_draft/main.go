package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"github.com/unhingeddev-200/action-sample/internal/gmail"
)

type Input struct {
	AccessToken string `json:"accessToken" jsonschema:"Gmail access token (injected from credentialId at publish time)"`
	UserID      string `json:"userId,omitempty" jsonschema:"Gmail account; use me for the connected mailbox"`
	To          string `json:"to" jsonschema:"Recipient email address"`
	From        string `json:"from,omitempty" jsonschema:"From address (optional)"`
	Subject     string `json:"subject" jsonschema:"Draft subject line"`
	TextBody    string `json:"textBody,omitempty" jsonschema:"Plain text body"`
	HtmlBody    string `json:"htmlBody,omitempty" jsonschema:"HTML body (optional)"`
}

type Output struct {
	ID              string   `json:"id" jsonschema:"Gmail draft id"`
	MessageID       string   `json:"messageId,omitempty" jsonschema:"Message id inside the draft"`
	MessageThreadID string   `json:"messageThreadId,omitempty" jsonschema:"Thread id for the draft message"`
	LabelIDs        []string `json:"labelIds,omitempty" jsonschema:"Label ids on the draft message"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailCreateDraft",
		Description: "Create a Gmail draft (does not send). Uses a saved Google credential.",
	}, func(ctx context.Context, in Input) (Output, error) {
		if strings.TrimSpace(in.AccessToken) == "" {
			return Output{}, fmt.Errorf("accessToken is required")
		}
		raw, err := gmail.BuildRawMessage(in.From, in.To, in.Subject, in.TextBody, in.HtmlBody)
		if err != nil {
			return Output{}, err
		}
		draft, err := gmail.NewClient(in.AccessToken).CreateDraft(ctx, in.UserID, raw)
		if err != nil {
			return Output{}, err
		}
		var out Output
		if id, ok := draft["id"].(string); ok {
			out.ID = id
		}
		if msg, ok := draft["message"].(map[string]any); ok {
			if id, ok := msg["id"].(string); ok {
				out.MessageID = id
			}
			if tid, ok := msg["threadId"].(string); ok {
				out.MessageThreadID = tid
			}
			if labels, ok := msg["labelIds"].([]any); ok {
				for _, l := range labels {
					if s, ok := l.(string); ok {
						out.LabelIDs = append(out.LabelIDs, s)
					}
				}
			}
		}
		return out, nil
	})
}
