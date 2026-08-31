package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	A []any `json:"a"`
	B []any `json:"b"`
	Key string `json:"key"`
}
type Output struct {
	OnlyA []any `json:"onlyA"`
	OnlyB []any `json:"onlyB"`
	Both []any `json:"both"`
}

func main() {
	action.Main(action.Meta{
		Name:        "CompareDatasets",
		Description: "diff arrays by key",
	}, func(ctx context.Context, in Input) (Output, error) {
		index := func(items []any) map[string]any {
			m := map[string]any{}
			for _, it := range items { m[fmt.Sprint(fieldOf(it, in.Key))] = it }
			return m
		}
		ia, ib := index(in.A), index(in.B)
		var onlyA, onlyB, both []any
		for k, v := range ia {
			if _, ok := ib[k]; ok { both = append(both, v) } else { onlyA = append(onlyA, v) }
		}
		for k, v := range ib {
			if _, ok := ia[k]; !ok { onlyB = append(onlyB, v) }
		}
		return Output{OnlyA: onlyA, OnlyB: onlyB, Both: both}, nil

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
