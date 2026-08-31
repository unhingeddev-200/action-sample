package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Items []any `json:"items"`
	GroupBy string `json:"groupBy"`
	Field string `json:"field,omitempty"`
	Op string `json:"op"`
}
type Output struct { Rows []map[string]any `json:"rows"` }

func main() {
	action.Main(action.Meta{
		Name:        "Summarize",
		Description: "group by summarize",
	}, func(ctx context.Context, in Input) (Output, error) {
		groups := map[string][]any{}
		order := []string{}
		for _, it := range in.Items {
			k := fmt.Sprint(fieldOf(it, in.GroupBy))
			if _, ok := groups[k]; !ok { order = append(order, k) }
			groups[k] = append(groups[k], it)
		}
		var rows []map[string]any
		for _, k := range order {
			aggIn := Input{Items: groups[k], Field: in.Field, Op: in.Op}
			// inline count/sum
			var result any
			switch strings.ToLower(in.Op) {
			case "count": result = len(groups[k])
			case "sum":
				s := 0.0
				for _, it := range groups[k] { if f, ok := asF(fieldOf(it, in.Field)); ok { s += f } }
				result = s
			case "avg":
				s := 0.0; n := 0
				for _, it := range groups[k] { if f, ok := asF(fieldOf(it, in.Field)); ok { s += f; n++ } }
				if n == 0 { result = 0 } else { result = s/float64(n) }
			default:
				result = len(groups[k])
			}
			_ = aggIn
			rows = append(rows, map[string]any{"key": k, "value": result})
		}
		return Output{Rows: rows}, nil

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
