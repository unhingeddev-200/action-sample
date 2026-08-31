package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Items []any `json:"items"`
	Max int `json:"max"`
}
type Output struct { Items []any `json:"items"` }

func main() {
	action.Main(action.Meta{
		Name:        "Limit",
		Description: "limit items",
	}, func(ctx context.Context, in Input) (Output, error) {
		items := in.Items
		if in.Max < 0 { in.Max = 0 }
		if in.Max < len(items) { items = items[:in.Max] }
		return Output{Items: items}, nil

	})
}

