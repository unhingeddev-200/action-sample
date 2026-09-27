package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"strings"
)


type Input struct {
	Mode   string         `json:"mode,omitempty" jsonschema:"How to combine inputs: combine (merge objects), append (list of items), or chooseBranch"`
	Inputs []map[string]any `json:"inputs,omitempty" jsonschema:"Objects from upstream steps to merge (often filled via bindings)"`
}
type Output struct {
	Data map[string]any `json:"data" jsonschema:"Merged object when mode is combine"`
	Items []any `json:"items,omitempty" jsonschema:"List of items when mode is append"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Merge",
		Description: "Combine results from several previous steps into one object or list (fan-in).",
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

