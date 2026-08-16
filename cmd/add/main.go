package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	A float64 `json:"a"`
	B float64 `json:"b"`
}

type Output struct {
	Result float64 `json:"result"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Add",
		Description: "adds two numbers",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Result: in.A + in.B}, nil
	})
}
