package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"fmt"
	"strings"
)


type Input struct {
	Message string `json:"message" jsonschema:"Error message shown when the run fails"`
}
type Output struct {}

func main() {
	action.Main(action.Meta{
		Name:        "StopAndError",
		Description: "Stop the workflow and mark the run as failed, with a clear error message.",
	}, func(ctx context.Context, in Input) (Output, error) {
		msg := strings.TrimSpace(in.Message)
		if msg == "" { msg = "stopped" }
		return Output{}, fmt.Errorf("%s", msg)

	})
}

