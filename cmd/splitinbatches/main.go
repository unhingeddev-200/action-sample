package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Items []any `json:"items" jsonschema:"Full list to split into smaller groups"`
	BatchSize int `json:"batchSize" jsonschema:"How many items per batch (must be at least 1)"`
}
type Output struct { Batches [][]any `json:"batches" jsonschema:"List of batches; each batch is a smaller list of items"` }

func main() {
	action.Main(action.Meta{
		Name:        "SplitInBatches",
		Description: "Split a long list into smaller batches (for example, process 10 emails at a time).",
	}, func(ctx context.Context, in Input) (Output, error) {
		n := in.BatchSize
		if n < 1 { n = 10 }
		var batches [][]any
		for i := 0; i < len(in.Items); i += n {
			j := i + n
			if j > len(in.Items) { j = len(in.Items) }
			batches = append(batches, in.Items[i:j])
		}
		return Output{Batches: batches}, nil

	})
}

