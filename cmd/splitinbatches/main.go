package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Items []any `json:"items"`
	BatchSize int `json:"batchSize"`
}
type Output struct { Batches [][]any `json:"batches"` }

func main() {
	action.Main(action.Meta{
		Name:        "SplitInBatches",
		Description: "split array into batches",
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

