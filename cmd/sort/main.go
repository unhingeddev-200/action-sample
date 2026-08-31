package main

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Items []any `json:"items"`
	Field string `json:"field,omitempty"`
	Order string `json:"order,omitempty"`
}
type Output struct { Items []any `json:"items"` }

func main() {
	action.Main(action.Meta{
		Name:        "Sort",
		Description: "sort items",
	}, func(ctx context.Context, in Input) (Output, error) {
		items := append([]any(nil), in.Items...)
		desc := strings.EqualFold(in.Order, "desc")
		sort.SliceStable(items, func(i, j int) bool {
			a, b := fieldOf(items[i], in.Field), fieldOf(items[j], in.Field)
			cmp := cmpAny(a, b)
			if desc { return cmp > 0 }
			return cmp < 0
		})
		return Output{Items: items}, nil

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
