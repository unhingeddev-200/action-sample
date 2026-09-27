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
	Q           string `json:"q,omitempty" jsonschema:"Gmail search query (same syntax as the Gmail search box, for example is:unread newer_than:7d)"`
	MaxResults  int64  `json:"maxResults,omitempty" jsonschema:"Maximum number of messages to return (default around 25)"`
	PageToken   string `json:"pageToken,omitempty" jsonschema:"Token from a previous list call to fetch the next page"`
}

type Output struct {
	Messages           []gmail.MessageRef `json:"messages" jsonschema:"Matching messages (id/threadId); use Gmail Get for full content"`
	NextPageToken      string             `json:"nextPageToken,omitempty" jsonschema:"Pass to the next Gmail List call to continue paging"`
	ResultSizeEstimate int64              `json:"resultSizeEstimate,omitempty" jsonschema:"Approximate total matches reported by Gmail"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailList",
		Description: "List Gmail messages, optionally filtered with a Gmail search query. Pair with For Each + Gmail Get to process each message.",
	}, func(ctx context.Context, in Input) (Output, error) {
		if strings.TrimSpace(in.AccessToken) == "" {
			return Output{}, fmt.Errorf("accessToken is required")
		}
		res, err := gmail.NewClient(in.AccessToken).List(ctx, in.UserID, in.Q, in.MaxResults, in.PageToken)
		if err != nil {
			return Output{}, err
		}
		return Output{
			Messages:           res.Messages,
			NextPageToken:      res.NextPageToken,
			ResultSizeEstimate: res.ResultSizeEstimate,
		}, nil
	})
}
