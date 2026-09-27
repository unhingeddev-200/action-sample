package main

import (
	"context"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	Assignments []Assignment `json:"assignments,omitempty" jsonschema:"List of fields to set on the output object"`
}

type Assignment struct {
	Name  string `json:"name" jsonschema:"Field name to set on the output"`
	Value any    `json:"value" jsonschema:"Value to store (text, number, true/false, object, or list)"`
}

type Output struct {
	Data map[string]any `json:"data" jsonschema:"Object containing all assigned fields"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Set",
		Description: "Build or update a data object by setting named fields. Later steps can read those fields.",
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
