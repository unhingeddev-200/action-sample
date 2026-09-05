package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"github.com/unhingeddev-200/action-sample/internal/gmail"
)

type Input struct {
	AccessToken string `json:"accessToken"`
	UserID      string `json:"userId,omitempty"`
	To          string `json:"to"`
	From        string `json:"from,omitempty"`
	Subject     string `json:"subject"`
	TextBody    string `json:"textBody,omitempty"`
	HtmlBody    string `json:"htmlBody,omitempty"`
}

type Output struct {
	ID              string   `json:"id"`
	MessageID       string   `json:"messageId,omitempty"`
	MessageThreadID string   `json:"messageThreadId,omitempty"`
	LabelIDs        []string `json:"labelIds,omitempty"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailCreateDraft",
		Description: "Create a Gmail draft from to/subject/body fields.",
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
