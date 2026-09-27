package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"fmt"
	"regexp"
	"strings"
)


type Input struct {
	Operation string `json:"operation" jsonschema:"wrap (put text inside an HTML tag) or stripTags (remove HTML tags)"`
	Value string `json:"value" jsonschema:"HTML or plain text to process"`
	Tag string `json:"tag,omitempty" jsonschema:"Tag name used when wrapping (default div)"`
}
type Output struct { Result string `json:"result" jsonschema:"Processed HTML or plain text"` }

func main() {
	action.Main(action.Meta{
		Name:        "HTML",
		Description: "Simple HTML helpers: wrap text in a tag or strip tags to plain text.",
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

