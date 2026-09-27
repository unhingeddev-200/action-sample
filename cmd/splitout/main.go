package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Data any `json:"data,omitempty" jsonschema:"Object that contains the list field"`
	Field string `json:"field" jsonschema:"Name of the field that holds the list to extract"`
}
type Output struct { Items []any `json:"items" jsonschema:"The extracted list"` }

func main() {
	action.Main(action.Meta{
		Name:        "SplitOut",
		Description: "Take a list nested inside an object and promote it to a top-level list for later steps.",
	}, func(ctx context.Context, in Input) (Output, error) {
		m, _ := in.Data.(map[string]any)
		var raw any = in.Data
		if m != nil && in.Field != "" { raw = m[in.Field] }
		switch t := raw.(type) {
		case []any: return Output{Items: t}, nil
		default: return Output{Items: []any{raw}}, nil
		}

	})
}

func fieldOf(it any, field string) any {
	if field == "" { return it }
	if m, ok := it.(map[string]any); ok { return m[field] }
	return it
}
func asF(v any) (float64, bool) {
	switch t := v.(type) {
	case float64: return t, true
	case float32: return float64(t), true
	case int: return float64(t), true
	case int64: return float64(t), true
	case string:
		f, err := strconv.ParseFloat(t, 64); return f, err == nil
	default: return 0, false
	}
}
func cmpAny(a, b any) int {
	fa, oka := asF(a); fb, okb := asF(b)
	if oka && okb {
		if fa < fb { return -1 }; if fa > fb { return 1 }; return 0
	}
	sa, sb := fmt.Sprint(a), fmt.Sprint(b)
	if sa < sb { return -1 }; if sa > sb { return 1 }; return 0
}
func match(op string, left, right any) bool {
	switch op {
	case "eq": return fmt.Sprint(left) == fmt.Sprint(right)
	case "ne": return fmt.Sprint(left) != fmt.Sprint(right)
	case "contains": return strings.Contains(fmt.Sprint(left), fmt.Sprint(right))
	case "truthy":
		switch t := left.(type) {
		case nil: return false
		case bool: return t
		case string: return t != ""
		case float64: return t != 0
		default: return left != nil
		}
	case "gt": return cmpAny(left, right) > 0
	case "gte": return cmpAny(left, right) >= 0
	case "lt": return cmpAny(left, right) < 0
	case "lte": return cmpAny(left, right) <= 0
	default: return false
	}
}
