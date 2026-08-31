package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"fmt"
	"strings"
	"time"
)


type Input struct {
	Operation string `json:"operation"`
	Value string `json:"value,omitempty"`
	Layout string `json:"layout,omitempty"`
	Amount int `json:"amount,omitempty"`
	Unit string `json:"unit,omitempty"`
}
type Output struct { Result string `json:"result"` }

func main() {
	action.Main(action.Meta{
		Name:        "DateTime",
		Description: "date time helpers",
	}, func(ctx context.Context, in Input) (Output, error) {
		layout := in.Layout
		if layout == "" { layout = time.RFC3339 }
		now := time.Now().UTC()
		switch strings.ToLower(in.Operation) {
		case "now":
			return Output{Result: now.Format(layout)}, nil
		case "format":
			t, err := time.Parse(time.RFC3339, in.Value)
			if err != nil { t, err = time.Parse(layout, in.Value) }
			if err != nil { return Output{}, err }
			return Output{Result: t.Format(layout)}, nil
		case "parse":
			t, err := time.Parse(layout, in.Value)
			if err != nil { return Output{}, err }
			return Output{Result: t.UTC().Format(time.RFC3339)}, nil
		case "add":
			t := now
			if in.Value != "" {
				var err error
				t, err = time.Parse(time.RFC3339, in.Value)
				if err != nil { return Output{}, err }
			}
			d := time.Duration(in.Amount)
			switch strings.ToLower(in.Unit) {
			case "seconds": d *= time.Second
			case "minutes": d *= time.Minute
			case "hours": d *= time.Hour
			case "days": d *= 24 * time.Hour
			default: d *= time.Second
			}
			return Output{Result: t.Add(d).Format(time.RFC3339)}, nil
		default:
			return Output{}, fmt.Errorf("unsupported operation")
		}

	})
}

