package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	A float64 `json:"a" jsonschema:"First number to add"`
	B float64 `json:"b" jsonschema:"Second number to add"`
}

type Output struct {
	Result float64 `json:"result" jsonschema:"Sum of a and b"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Add",
		Description: "Add two numbers together and return the sum. Useful for simple math in a workflow.",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Result: in.A + in.B}, nil
	})
}
