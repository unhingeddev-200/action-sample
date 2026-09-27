package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct{}
type Output struct {
	OK bool `json:"ok" jsonschema:"Always true when the step finishes successfully"`
}

func main() {
	action.Main(action.Meta{
		Name:        "NoOp",
		Description: "Do nothing and continue. Useful as a placeholder or join point on the canvas.",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{OK: true}, nil
	})
}
