package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"strings"
)


type Input struct {
	Mode   string         `json:"mode,omitempty"`
	Inputs []map[string]any `json:"inputs,omitempty"`
}
type Output struct {
	Data map[string]any `json:"data"`
	Items []any `json:"items,omitempty"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Merge",
		Description: "merge input objects",
	}, func(ctx context.Context, in Input) (Output, error) {
		mode := strings.ToLower(strings.TrimSpace(in.Mode))
		if mode == "" { mode = "combine" }
		switch mode {
		case "append":
			var items []any
			for _, m := range in.Inputs { items = append(items, m) }
			return Output{Data: map[string]any{"count": len(items)}, Items: items}, nil
		default:
			out := map[string]any{}
			for _, m := range in.Inputs {
				for k, v := range m { out[k] = v }
			}
			return Output{Data: out}, nil
		}

	})
}

