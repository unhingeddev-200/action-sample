package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"encoding/xml"
	"io"
	"strings"
)


type Input struct { Value string `json:"value"` }
type Output struct { Data map[string]any `json:"data"` }

func main() {
	action.Main(action.Meta{
		Name:        "XML",
		Description: "parse xml to map best effort",
	}, func(ctx context.Context, in Input) (Output, error) {
		dec := xml.NewDecoder(strings.NewReader(in.Value))
		var root map[string]any
		var stack []map[string]any
		for {
			tok, err := dec.Token()
			if err == io.EOF { break }
			if err != nil { return Output{}, err }
			switch t := tok.(type) {
			case xml.StartElement:
				m := map[string]any{}
				if root == nil { root = m } else if len(stack)>0 {
					parent := stack[len(stack)-1]
					parent[t.Name.Local] = m
				}
				stack = append(stack, m)
			case xml.EndElement:
				if len(stack)>0 { stack = stack[:len(stack)-1] }
			case xml.CharData:
				if len(stack)==0 { continue }
				s := strings.TrimSpace(string(t))
				if s == "" { continue }
				stack[len(stack)-1]["_text"] = s
			}
		}
		if root == nil { root = map[string]any{} }
		return Output{Data: root}, nil

	})
}

