package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Items []any `json:"items" jsonschema:"List of items to trim"`
	Max int `json:"max" jsonschema:"Maximum number of items to keep (from the start of the list)"`
}
type Output struct { Items []any `json:"items" jsonschema:"First max items from the input list"` }

func main() {
	action.Main(action.Meta{
		Name:        "Limit",
		Description: "Keep only the first N items from a list.",
	}, func(ctx context.Context, in Input) (Output, error) {
		items := in.Items
		if in.Max < 0 { in.Max = 0 }
		if in.Max < len(items) { items = items[:in.Max] }
		return Output{Items: items}, nil

	})
}

