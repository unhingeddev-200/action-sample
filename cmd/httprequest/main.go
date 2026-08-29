package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	Method         string            `json:"method"`
	URL            string            `json:"url"`
	Headers        map[string]string `json:"headers,omitempty"`
	Query          map[string]string `json:"query,omitempty"`
	Body           any               `json:"body,omitempty"`
	ResponseFormat string            `json:"responseFormat,omitempty"`
}

type Output struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    any               `json:"body"`
}

func main() {
	action.Main(action.Meta{
		Name:        "HttpRequest",
		Description: "perform an HTTP request",
	}, func(ctx context.Context, in Input) (Output, error) {
		method := strings.ToUpper(strings.TrimSpace(in.Method))
		if method == "" {
			method = http.MethodGet
		}
		if strings.TrimSpace(in.URL) == "" {
			return Output{}, fmt.Errorf("url required")
		}
		u, err := url.Parse(in.URL)
		if err != nil {
			return Output{}, fmt.Errorf("url: %w", err)
		}
		q := u.Query()
		for k, v := range in.Query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()

		var bodyReader io.Reader
		if in.Body != nil && method != http.MethodGet && method != http.MethodHead {
			switch b := in.Body.(type) {
			case string:
				bodyReader = strings.NewReader(b)
			default:
				raw, err := json.Marshal(b)
				if err != nil {
					return Output{}, err
				}
				bodyReader = bytes.NewReader(raw)
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
		if err != nil {
			return Output{}, err
		}
		for k, v := range in.Headers {
			req.Header.Set(k, v)
		}
		if bodyReader != nil && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return Output{}, err
		}
		defer resp.Body.Close()
		raw, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
		if err != nil {
			return Output{}, err
		}
		outHeaders := map[string]string{}
		for k, vs := range resp.Header {
			if len(vs) > 0 {
				outHeaders[k] = vs[0]
			}
		}
		format := strings.ToLower(strings.TrimSpace(in.ResponseFormat))
		if format == "" {
			format = "json"
		}
		var body any
		if format == "text" {
			body = string(raw)
		} else {
			if len(raw) == 0 {
				body = nil
			} else if err := json.Unmarshal(raw, &body); err != nil {
				body = string(raw)
			}
		}
		return Output{Status: resp.StatusCode, Headers: outHeaders, Body: body}, nil
	})
}
