package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	Assignments []Assignment `json:"assignments"`
}

type Assignment struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type Output struct {
	Data map[string]any `json:"data"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Set",
		Description: "assign fields into an output object",
	}, func(ctx context.Context, in Input) (Output, error) {
		out := make(map[string]any, len(in.Assignments))
		for _, a := range in.Assignments {
			if a.Name == "" {
				continue
			}
			out[a.Name] = a.Value
		}
		return Output{Data: out}, nil
	})
}
