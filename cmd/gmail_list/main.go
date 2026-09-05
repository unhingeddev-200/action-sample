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
	Q           string `json:"q,omitempty"`
	MaxResults  int64  `json:"maxResults,omitempty"`
	PageToken   string `json:"pageToken,omitempty"`
}

type Output struct {
	Messages           []gmail.MessageRef `json:"messages"`
	NextPageToken      string             `json:"nextPageToken,omitempty"`
	ResultSizeEstimate int64              `json:"resultSizeEstimate,omitempty"`
}

func main() {
	action.Main(action.Meta{
		Name:        "GmailList",
		Description: "List Gmail messages (optional search query q).",
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
