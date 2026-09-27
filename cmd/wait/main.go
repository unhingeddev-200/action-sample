package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	Seconds float64 `json:"seconds,omitempty" jsonschema:"How long to wait, in seconds. Values of 60 or more use a durable wait that survives restarts; shorter waits pause inside the step."`
}

type Output struct {
	Slept float64 `json:"slept" jsonschema:"How many seconds the step actually waited"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Wait",
		Description: "Pause the workflow for a number of seconds before continuing. Long waits (60s+) are durable and survive restarts.",
	}, func(ctx context.Context, in Input) (Output, error) {
		if in.Seconds < 0 {
			return Output{}, fmt.Errorf("seconds must be >= 0")
		}
		if in.Seconds > 3600 {
			return Output{}, fmt.Errorf("seconds max 3600")
		}
		d := time.Duration(in.Seconds * float64(time.Second))
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return Output{}, ctx.Err()
		case <-timer.C:
			return Output{Slept: in.Seconds}, nil
		}
	})
}
