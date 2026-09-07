package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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
		_ = vm.Set("env", environMap())
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
		v, err := vm.RunString(wrapUserSource(in.Source))
		if err != nil {
			return Output{}, err
		}
		return Output{Result: v.Export()}, nil
	})
}

// wrapUserSource runs user code inside an IIFE so `return` works at the top level
// (goja rejects bare return statements in RunString scripts).
func wrapUserSource(source string) string {
	s := strings.TrimSpace(source)
	if s == "" {
		return s
	}
	lower := strings.ToLower(s)
	if strings.HasPrefix(lower, "(function") || strings.HasPrefix(lower, "function") {
		return s
	}
	return "(function() {\n" + source + "\n})()"
}

func environMap() map[string]string {
	out := make(map[string]string)
	for _, e := range os.Environ() {
		k, v, ok := strings.Cut(e, "=")
		if !ok || k == "" {
			continue
		}
		out[k] = v
	}
	return out
}
