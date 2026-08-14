package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct{}

type Output struct {
	Message string `json:"message"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Sample",
		Description: "sample hello-world action",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Message: "Hello From Action"}, nil
	})
}
