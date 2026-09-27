package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"github.com/unhingeddev-200/action-sample/internal/gmail"
)

type Input struct {
	AccessToken string `json:"accessToken" jsonschema:"Gmail access token (injected automatically from credentialId when publishing a product workflow)"`
	UserID      string `json:"userId,omitempty" jsonschema:"Gmail account to use; leave as me for the signed-in mailbox"`
	To          string `json:"to" jsonschema:"Recipient email address"`
	From        string `json:"from,omitempty" jsonschema:"From address (optional; defaults to the mailbox)"`
	Subject     string `json:"subject" jsonschema:"Email subject line"`
	TextBody    string `json:"textBody,omitempty" jsonschema:"Plain text body"`
	HtmlBody    string `json:"htmlBody,omitempty" jsonschema:"HTML body (optional alternative to plain text)"`
}

type Output struct {
	ID       string   `json:"id" jsonschema:"Gmail message id of the sent email"`
	ThreadID string   `json:"threadId,omitempty" jsonschema:"Gmail thread id"`
	LabelIDs []string `json:"labelIds,omitempty" jsonschema:"Gmail label ids applied to the message"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailSend",
		Description: "Send an email with Gmail using a saved Google credential (needs gmail.send permission).",
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
