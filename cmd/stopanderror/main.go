package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"fmt"
	"strings"
)


type Input struct { Message string `json:"message"` }
type Output struct {}

func main() {
	action.Main(action.Meta{
		Name:        "StopAndError",
		Description: "fail the run",
	}, func(ctx context.Context, in Input) (Output, error) {
		msg := strings.TrimSpace(in.Message)
		if msg == "" { msg = "stopped" }
		return Output{}, fmt.Errorf("%s", msg)

	})
}

