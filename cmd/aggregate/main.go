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
	Field string `json:"field,omitempty"`
	Op string `json:"op"`
	Separator string `json:"separator,omitempty"`
}
type Output struct { Result any `json:"result"` }

func main() {
	action.Main(action.Meta{
		Name:        "Aggregate",
		Description: "aggregate field",
	}, func(ctx context.Context, in Input) (Output, error) {
		op := strings.ToLower(in.Op)
		switch op {
		case "count":
			return Output{Result: len(in.Items)}, nil
		case "sum", "avg", "min", "max":
			var nums []float64
			for _, it := range in.Items {
				v := fieldOf(it, in.Field)
				if f, ok := asF(v); ok { nums = append(nums, f) }
			}
			if len(nums) == 0 { return Output{Result: 0}, nil }
			switch op {
			case "sum":
				s := 0.0; for _, n := range nums { s += n }; return Output{Result: s}, nil
			case "avg":
				s := 0.0; for _, n := range nums { s += n }; return Output{Result: s/float64(len(nums))}, nil
			case "min":
				m := nums[0]; for _, n := range nums { if n < m { m = n } }; return Output{Result: m}, nil
			case "max":
				m := nums[0]; for _, n := range nums { if n > m { m = n } }; return Output{Result: m}, nil
			}
		case "join":
			sep := in.Separator
			if sep == "" { sep = "," }
			var parts []string
			for _, it := range in.Items { parts = append(parts, fmt.Sprint(fieldOf(it, in.Field))) }
			return Output{Result: strings.Join(parts, sep)}, nil
		}
		return Output{}, fmt.Errorf("unsupported op %q", in.Op)

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
