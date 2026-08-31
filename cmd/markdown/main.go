package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct { Value string `json:"value"` }
type Output struct { Markdown string `json:"markdown"` }

func main() {
	action.Main(action.Meta{
		Name:        "Markdown",
		Description: "markdown passthrough",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Markdown: in.Value}, nil

	})
}

