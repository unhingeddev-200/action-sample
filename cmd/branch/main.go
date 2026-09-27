package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Case struct {
	Match any    `json:"match,omitempty" jsonschema:"Value that selects this branch in switch mode"`
	Label string `json:"label" jsonschema:"Label of the outgoing edge to follow when this case matches"`
}

type Input struct {
	Mode    string `json:"mode" jsonschema:"Branch style: if (true/false style) or switch (pick among labeled cases)"`
	Op      string `json:"op,omitempty" jsonschema:"Comparison for if mode: eq, ne, gt, gte, lt, lte, or truthy"`
	Left    any    `json:"left,omitempty" jsonschema:"Left side of the comparison (if mode)"`
	Right   any    `json:"right,omitempty" jsonschema:"Right side of the comparison (if mode)"`
	Value   any    `json:"value,omitempty" jsonschema:"Value to match against cases (switch mode)"`
	Cases   []Case `json:"cases,omitempty" jsonschema:"Named cases for switch mode; each case points to an edge label"`
	Default string `json:"default,omitempty" jsonschema:"Edge label to follow when no case matches"`
}

type Output struct {
	Selected string `json:"selected" jsonschema:"Label of the edge that was chosen"`
}

func main() {
	action.Main(action.Meta{
		Name:        "Branch",
		Description: "Choose which path the workflow takes next (if/else or switch among labeled edges).",
	}, func(ctx context.Context, in Input) (Output, error) {
		switch stringsToLower(in.Mode) {
		case "if":
			ok, err := evalIf(in)
			if err != nil {
				return Output{}, err
			}
			if ok {
				return Output{Selected: "true"}, nil
			}
			return Output{Selected: "false"}, nil
		case "switch":
			sel, err := evalSwitch(in)
			if err != nil {
				return Output{}, err
			}
			return Output{Selected: sel}, nil
		default:
			return Output{}, fmt.Errorf("mode must be if or switch")
		}
	})
}

func stringsToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func evalIf(in Input) (bool, error) {
	op := stringsToLower(in.Op)
	if op == "" {
		op = "truthy"
	}
	switch op {
	case "truthy":
		return isTruthy(in.Left), nil
	case "eq":
		return softEqual(in.Left, in.Right), nil
	case "ne":
		return !softEqual(in.Left, in.Right), nil
	case "gt", "gte", "lt", "lte":
		lf, ok1 := asFloat(in.Left)
		rf, ok2 := asFloat(in.Right)
		if !ok1 || !ok2 {
			return false, fmt.Errorf("numeric compare requires numbers")
		}
		switch op {
		case "gt":
			return lf > rf, nil
		case "gte":
			return lf >= rf, nil
		case "lt":
			return lf < rf, nil
		default:
			return lf <= rf, nil
		}
	default:
		return false, fmt.Errorf("unsupported op %q", in.Op)
	}
}

func evalSwitch(in Input) (string, error) {
	for _, c := range in.Cases {
		if softEqual(in.Value, c.Match) {
			if c.Label == "" {
				return "", fmt.Errorf("case label required")
			}
			return c.Label, nil
		}
	}
	if in.Default != "" {
		return in.Default, nil
	}
	return "", fmt.Errorf("no matching case and no default")
}

func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return t != ""
	case float64:
		return t != 0
	case json.Number:
		f, _ := t.Float64()
		return f != 0
	default:
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Map:
			return rv.Len() > 0
		default:
			return true
		}
	}
}

func softEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if fa, ok := asFloat(a); ok {
		if fb, ok2 := asFloat(b); ok2 {
			return fa == fb
		}
	}
	return fmt.Sprint(a) == fmt.Sprint(b)
}

func asFloat(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case float32:
		return float64(t), true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	case json.Number:
		f, err := t.Float64()
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(t, 64)
		return f, err == nil
	default:
		return 0, false
	}
}
