package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Operation string `json:"operation"`
	Items []any `json:"items,omitempty"`
	Field string `json:"field,omitempty"`
	Limit int `json:"limit,omitempty"`
}
type Output struct { Items []any `json:"items"` }

func main() {
	action.Main(action.Meta{
		Name:        "ItemLists",
		Description: "item list helpers",
	}, func(ctx context.Context, in Input) (Output, error) {
		switch strings.ToLower(in.Operation) {
		case "limit":
			items := in.Items
			if in.Limit >= 0 && in.Limit < len(items) { items = items[:in.Limit] }
			return Output{Items: items}, nil
		case "concatenate":
			var out []any
			for _, it := range in.Items {
				if arr, ok := it.([]any); ok { out = append(out, arr...) } else { out = append(out, it) }
			}
			return Output{Items: out}, nil
		case "splitout":
			var out []any
			for _, it := range in.Items {
				v := fieldOf(it, in.Field)
				if arr, ok := v.([]any); ok { out = append(out, arr...) } else { out = append(out, v) }
			}
			return Output{Items: out}, nil
		default:
			return Output{}, fmt.Errorf("unsupported operation %q", in.Operation)
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
