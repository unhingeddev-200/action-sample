package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Value string `json:"value" jsonschema:"Markdown text"`
}

type Output struct {
	Markdown string `json:"markdown" jsonschema:"Markdown text passed through for later steps or the UI"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Markdown",
		Description: "Pass Markdown text through the workflow (rendering is left to the client/UI).",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Markdown: in.Value}, nil

	})
}

