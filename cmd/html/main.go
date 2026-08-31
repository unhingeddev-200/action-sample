package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"fmt"
	"regexp"
	"strings"
)


type Input struct {
	Operation string `json:"operation"`
	Value string `json:"value"`
	Tag string `json:"tag,omitempty"`
}
type Output struct { Result string `json:"result"` }

func main() {
	action.Main(action.Meta{
		Name:        "HTML",
		Description: "html wrap/strip",
	}, func(ctx context.Context, in Input) (Output, error) {
		tag := in.Tag
		if tag == "" { tag = "div" }
		switch strings.ToLower(in.Operation) {
		case "wrap":
			return Output{Result: fmt.Sprintf("<%s>%s</%s>", tag, in.Value, tag)}, nil
		case "striptags":
			re := regexp.MustCompile(`<[^>]*>`)
			return Output{Result: re.ReplaceAllString(in.Value, "")}, nil
		default:
			return Output{}, fmt.Errorf("unsupported operation")
		}

	})
}

