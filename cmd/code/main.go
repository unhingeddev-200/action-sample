package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dop251/goja"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	Source string `json:"source"`
	Input  any    `json:"input,omitempty"`
}

type Output struct {
	Result any `json:"result"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Code",
		Description: "run JavaScript (goja) against input JSON",
	}, func(ctx context.Context, in Input) (Output, error) {
		if in.Source == "" {
			return Output{}, fmt.Errorf("source required")
		}
		vm := goja.New()
		_ = vm.Set("input", in.Input)
		_ = vm.Set("JSON", map[string]any{
			"stringify": func(v any) string {
				b, _ := json.Marshal(v)
				return string(b)
			},
			"parse": func(s string) any {
				var v any
				_ = json.Unmarshal([]byte(s), &v)
				return v
			},
		})
		v, err := vm.RunString(in.Source)
		if err != nil {
			return Output{}, err
		}
		return Output{Result: v.Export()}, nil
	})
}
