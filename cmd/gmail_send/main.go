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
	ID       string   `json:"id"`
	ThreadID string   `json:"threadId,omitempty"`
	LabelIDs []string `json:"labelIds,omitempty"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailSend",
		Description: "Send an email via Gmail API (composes MIME from to/subject/body).",
	}, func(ctx context.Context, in Input) (Output, error) {
		if strings.TrimSpace(in.AccessToken) == "" {
			return Output{}, fmt.Errorf("accessToken is required")
		}
		raw, err := gmail.BuildRawMessage(in.From, in.To, in.Subject, in.TextBody, in.HtmlBody)
		if err != nil {
			return Output{}, err
		}
		ref, err := gmail.NewClient(in.AccessToken).Send(ctx, in.UserID, raw)
		if err != nil {
			return Output{}, err
		}
		return Output{ID: ref.ID, ThreadID: ref.ThreadID, LabelIDs: ref.LabelIDs}, nil
	})
}
