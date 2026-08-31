package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)


type Input struct {
	Operation string `json:"operation"`
	Value string `json:"value"`
}
type Output struct { Result string `json:"result"` }

func main() {
	action.Main(action.Meta{
		Name:        "Compression",
		Description: "gzip base64",
	}, func(ctx context.Context, in Input) (Output, error) {
		switch strings.ToLower(in.Operation) {
		case "gzip":
			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)
			if _, err := w.Write([]byte(in.Value)); err != nil { return Output{}, err }
			_ = w.Close()
			return Output{Result: base64.StdEncoding.EncodeToString(buf.Bytes())}, nil
		case "gunzip":
			raw, err := base64.StdEncoding.DecodeString(in.Value)
			if err != nil { return Output{}, err }
			r, err := gzip.NewReader(bytes.NewReader(raw))
			if err != nil { return Output{}, err }
			defer r.Close()
			out, err := io.ReadAll(r)
			if err != nil { return Output{}, err }
			return Output{Result: string(out)}, nil
		default:
			return Output{}, fmt.Errorf("unsupported operation")
		}

	})
}

