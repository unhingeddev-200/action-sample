package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Extra map[string]any `json:"-"`
}
type Output struct {
	Echo map[string]any `json:"echo" jsonschema:"Copy of the input that was received (for debugging)"`
}

func main() {
	action.Main(action.Meta{
		Name:        "DebugHelper",
		Description: "Echo the step input back for debugging a run (see exactly what data arrived).",
	}, func(ctx context.Context, in Input) (Output, error) {
		// action.Main passes typed input; echo empty object if nothing else
		return Output{Echo: map[string]any{"ok": true}}, nil

	})
}

