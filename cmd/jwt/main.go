package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Operation string `json:"operation"`
	Payload map[string]any `json:"payload,omitempty"`
	Token string `json:"token,omitempty"`
	Secret string `json:"secret,omitempty"`
}
type Output struct {
	Token string `json:"token,omitempty"`
	Claims map[string]any `json:"claims,omitempty"`
}

func main() {
	action.Main(action.Meta{
		Name:        "JWT",
		Description: "sign or decode jwt hs256",
	}, func(ctx context.Context, in Input) (Output, error) {
		switch strings.ToLower(in.Operation) {
		case "sign":
			if in.Secret == "" { return Output{}, fmt.Errorf("secret required") }
			header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
			payload, _ := json.Marshal(in.Payload)
			body := base64.RawURLEncoding.EncodeToString(payload)
			sig := hmacBytes("sha256", []byte(in.Secret), []byte(header+"."+body))
			return Output{Token: header+"."+body+"."+base64.RawURLEncoding.EncodeToString(sig)}, nil
		case "decode":
			parts := strings.Split(in.Token, ".")
			if len(parts) < 2 { return Output{}, fmt.Errorf("invalid token") }
			raw, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil { return Output{}, err }
			var claims map[string]any
			if err := json.Unmarshal(raw, &claims); err != nil { return Output{}, err }
			return Output{Claims: claims}, nil
		default:
			return Output{}, fmt.Errorf("unsupported operation")
		}

	})
}

func hmacBytes(alg string, key, b []byte) []byte {
	h := hmac.New(sha256.New, key)
	_ = alg
	h.Write(b)
	return h.Sum(nil)
}
