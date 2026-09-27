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
	Operation string `json:"operation" jsonschema:"sign (create a token) or decode (read a token)"`
	Payload map[string]any `json:"payload,omitempty" jsonschema:"Claims/data to put inside the token when signing"`
	Token string `json:"token,omitempty" jsonschema:"Existing JWT to decode (input) or signed JWT (output)"`
	Secret string `json:"secret,omitempty" jsonschema:"Shared secret used to sign or verify (HS256)"`
}
type Output struct {
	Token string `json:"token,omitempty" jsonschema:"Existing JWT to decode (input) or signed JWT (output)"`
	Claims map[string]any `json:"claims,omitempty" jsonschema:"Decoded claims when operation is decode"`
}

func main() {
	action.Main(action.Meta{
		Name:        "JWT",
		Description: "Create or read a JWT (JSON Web Token) using a shared secret (HS256).",
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
