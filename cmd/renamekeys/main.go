package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Data map[string]any `json:"data,omitempty" jsonschema:"Object whose keys should be renamed"`
	Mapping map[string]string `json:"mapping" jsonschema:"Map of old field name → new field name"`
}
type Output struct { Data map[string]any `json:"data" jsonschema:"Object with renamed keys"` }

func main() {
	action.Main(action.Meta{
		Name:        "RenameKeys",
		Description: "Rename fields on an object (for example change customer_name to customerName).",
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

