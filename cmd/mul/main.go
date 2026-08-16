package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	In1 float64 `json:"in1"`
	In2 float64 `json:"in2"`
}

type Output struct {
	Result float64 `json:"result"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Mul",
		Description: "multiplies two numbers",
	}, func(ctx context.Context, in Input) (Output, error) {
		return Output{Result: in.In1 * in.In2}, nil
	})
}
