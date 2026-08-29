package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct{}
type Output struct {
	OK bool `json:"ok"`
}

func main() {
	action.Main(action.Meta{
		Name:        "NoOp",
		Description: "no-op pass-through",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{OK: true}, nil
	})
}
