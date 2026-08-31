package main

import (
	"context"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)


type Input struct {
	WorkflowId string `json:"workflowId"`
	BffUrl string `json:"bffUrl,omitempty"`
	Token string `json:"token,omitempty"`
}
type Output struct {
	Status int `json:"status"`
	Body any `json:"body"`
}

func main() {
	action.Main(action.Meta{
		Name:        "ExecuteWorkflow",
		Description: "start another workflow via BFF",
	}, func(ctx context.Context, in Input) (Output, error) {
		base := strings.TrimSpace(in.BffUrl)
		if base == "" { base = strings.TrimSpace(os.Getenv("GENIUS_BFF_URL")) }
		tok := strings.TrimSpace(in.Token)
		if tok == "" { tok = strings.TrimSpace(os.Getenv("GENIUS_SERVICE_TOKEN")) }
		if base == "" || tok == "" {
			return Output{}, fmt.Errorf("ExecuteWorkflow requires bffUrl/token config or GENIUS_BFF_URL + GENIUS_SERVICE_TOKEN")
		}
		if strings.TrimSpace(in.WorkflowId) == "" {
			return Output{}, fmt.Errorf("workflowId required")
		}
		payload, _ := json.Marshal(map[string]any{"workflowId": in.WorkflowId})
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(base, "/")+"/genius.workflows.v1.ProductWorkflows/StartRun", bytes.NewReader(payload))
		if err != nil { return Output{}, err }
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+tok)
		resp, err := http.DefaultClient.Do(req)
		if err != nil { return Output{}, err }
		defer resp.Body.Close()
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
		var body any
		if err := json.Unmarshal(raw, &body); err != nil { body = string(raw) }
		return Output{Status: resp.StatusCode, Body: body}, nil

	})
}

