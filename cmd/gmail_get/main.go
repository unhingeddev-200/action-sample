package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"github.com/unhingeddev-200/action-sample/internal/gmail"
)

type Input struct {
	AccessToken string `json:"accessToken" jsonschema:"Gmail access token (injected from credentialId at publish time)"`
	UserID      string `json:"userId,omitempty" jsonschema:"Gmail account; use me for the connected mailbox"`
	MessageID   string `json:"messageId" jsonschema:"Gmail message id (often bound from Gmail List / For Each item.id)"`
	Format      string `json:"format,omitempty" jsonschema:"How much detail to fetch: full (default), metadata, minimal, or raw"` // full | metadata | minimal | raw
}

type Output struct {
	ID           string   `json:"id" jsonschema:"Gmail message id"`
	ThreadID     string   `json:"threadId,omitempty" jsonschema:"Gmail thread id"`
	LabelIDs     []string `json:"labelIds,omitempty" jsonschema:"Label ids on the message"`
	Snippet      string   `json:"snippet,omitempty" jsonschema:"Short preview text of the message"`
	InternalDate string   `json:"internalDate,omitempty" jsonschema:"Gmail internal timestamp"`
	// Payload is the Gmail message.payload object (MIME tree). Use any so the
	// generated JSON Schema accepts an object (json.RawMessage is []byte → array).
	Payload      any    `json:"payload,omitempty" jsonschema:"Raw Gmail payload structure (advanced)"`
	SizeEstimate int64  `json:"sizeEstimate,omitempty" jsonschema:"Approximate message size in bytes"`
	Raw          string `json:"raw,omitempty" jsonschema:"Raw RFC822 content when format is raw"`
	// Flat fields for ProductCall / sources.create bindings.
	Subject    string `json:"subject,omitempty" jsonschema:"Email subject"`
	BodyText   string `json:"bodyText,omitempty" jsonschema:"Plain text body extracted from the message"`
	From       string `json:"from,omitempty" jsonschema:"Sender address"`
	OccurredAt string `json:"occurredAt,omitempty" jsonschema:"When the message occurred (RFC3339 when available)"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailGet",
		Description: "Fetch one Gmail message by id, including subject and body text helpers for later steps.",
	}, func(ctx context.Context, in Input) (Output, error) {
		if strings.TrimSpace(in.AccessToken) == "" {
			return Output{}, fmt.Errorf("accessToken is required")
		}
		msg, err := gmail.NewClient(in.AccessToken).Get(ctx, in.UserID, in.MessageID, in.Format)
		if err != nil {
			return Output{}, err
		}
		var payload any
		if len(msg.Payload) > 0 {
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				return Output{}, fmt.Errorf("decode payload: %w", err)
			}
		}
		fields := gmail.ParseMessageFields(msg)
		return Output{
			ID:           msg.ID,
			ThreadID:     msg.ThreadID,
			LabelIDs:     msg.LabelIDs,
			Snippet:      msg.Snippet,
			InternalDate: msg.InternalDate,
			Payload:      payload,
			SizeEstimate: msg.SizeEstimate,
			Raw:          msg.Raw,
			Subject:      fields.Subject,
			BodyText:     fields.BodyText,
			From:         fields.From,
			OccurredAt:   fields.OccurredAt,
		}, nil
	})
}
