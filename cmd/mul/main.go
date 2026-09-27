package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	In1 float64 `json:"in1" jsonschema:"First number to multiply"`
	In2 float64 `json:"in2" jsonschema:"Second number to multiply"`
}

type Output struct {
	Result float64 `json:"result" jsonschema:"Product of in1 and in2"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Mul",
		Description: "Multiply two numbers and return the product.",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Result: in.In1 * in.In2}, nil
	})
}
