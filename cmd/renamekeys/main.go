package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Data map[string]any `json:"data,omitempty"`
	Mapping map[string]string `json:"mapping"`
}
type Output struct { Data map[string]any `json:"data"` }

func main() {
	action.Main(action.Meta{
		Name:        "RenameKeys",
		Description: "rename object keys",
	}, func(ctx context.Context, in Input) (Output, error) {
		src := in.Data
		if src == nil { src = map[string]any{} }
		out := map[string]any{}
		for k, v := range src {
			if nk, ok := in.Mapping[k]; ok && nk != "" { out[nk] = v } else { out[k] = v }
		}
		return Output{Data: out}, nil

	})
}

