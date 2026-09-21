// ProductCall: perform a governed Work Office operation (work item create /
// transition / assign / comment) from inside a workflow run.
//
// This is the workflow-side twin of what a human does in the UI. It does NOT
// hand-roll HTTP: it uses the runtime SDK for run identity + internal auth and
// the generated Connect client for genius.productops.v1.ProductOps. The BFF
// scopes the call to the run's org, acts as the run's attributed user, and
// applies the same permission checks, history entries, and events as the
// public RPCs.
package main

import (
	"context"
	"fmt"
	"strings"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/structpb"

	productopsv1 "github.com/lunaya-dubai/genius-ai-api/gen/genius/productops/v1"
	"github.com/lunaya-dubai/genius-ai-api/gen/genius/productops/v1/productopsv1connect"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

type Input struct {
	// Operation is a catalog name: workitems.create | workitems.transition |
	// workitems.assign | workitems.comment (allowlist enforced by the BFF).
	Operation string `json:"operation"`
	// Payload holds operation arguments; upstream steps typically bind
	// fields like workItemId here.
	Payload map[string]any `json:"payload,omitempty"`
}

type Output struct {
	Operation string         `json:"operation"`
	Result    map[string]any `json:"result"`
}

func main() {
	action.Main(action.Meta{
		Name:        "ProductCall",
		Description: "Invoke a governed product operation (work items) via the BFF ProductOps service, attributed to the run.",
	}, func(ctx context.Context, in Input) (Output, error) {
		op := strings.TrimSpace(in.Operation)
		if op == "" {
			return Output{}, fmt.Errorf("operation is required")
		}
		rc, err := action.RunContextFromEnv()
		if err != nil {
			return Output{}, err
		}
		base, hc, opts, err := action.InternalRPC()
		if err != nil {
			return Output{}, err
		}
		payload, err := structpb.NewStruct(in.Payload)
		if err != nil {
			return Output{}, fmt.Errorf("payload: %w", err)
		}
		client := productopsv1connect.NewProductOpsClient(hc, base, opts...)
		resp, err := client.Invoke(ctx, &productopsv1.InvokeRequest{
			RuntimeRunId: rc.RuntimeRunID,
			NodeId:       rc.NodeID,
			Operation:    op,
			Payload:      payload,
		})
		if err != nil {
			var cerr *connect.Error
			if ok := asConnectError(err, &cerr); ok {
				return Output{}, fmt.Errorf("%s: %s: %s", op, cerr.Code(), cerr.Message())
			}
			return Output{}, fmt.Errorf("%s: %w", op, err)
		}
		return Output{Operation: op, Result: resp.GetResult().AsMap()}, nil
	})
}

func asConnectError(err error, target **connect.Error) bool {
	ce, ok := err.(*connect.Error)
	if ok {
		*target = ce
	}
	return ok
}
