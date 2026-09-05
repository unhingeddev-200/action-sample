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
	AccessToken string `json:"accessToken"`
	UserID      string `json:"userId,omitempty"`
	MessageID   string `json:"messageId"`
	Format      string `json:"format,omitempty"` // full | metadata | minimal | raw
}

type Output struct {
	ID           string          `json:"id"`
	ThreadID     string          `json:"threadId,omitempty"`
	LabelIDs     []string        `json:"labelIds,omitempty"`
	Snippet      string          `json:"snippet,omitempty"`
	InternalDate string          `json:"internalDate,omitempty"`
	Payload      json.RawMessage `json:"payload,omitempty"`
	SizeEstimate int64           `json:"sizeEstimate,omitempty"`
	Raw          string          `json:"raw,omitempty"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailGet",
		Description: "Get a Gmail message by id.",
	}, func(ctx context.Context, in Input) (Output, error) {
		if strings.TrimSpace(in.AccessToken) == "" {
			return Output{}, fmt.Errorf("accessToken is required")
		}
		msg, err := gmail.NewClient(in.AccessToken).Get(ctx, in.UserID, in.MessageID, in.Format)
		if err != nil {
			return Output{}, err
		}
		return Output{
			ID:           msg.ID,
			ThreadID:     msg.ThreadID,
			LabelIDs:     msg.LabelIDs,
			Snippet:      msg.Snippet,
			InternalDate: msg.InternalDate,
			Payload:      msg.Payload,
			SizeEstimate: msg.SizeEstimate,
			Raw:          msg.Raw,
		}, nil
	})
}
