package main

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"strings"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Operation string `json:"operation"`
	Algorithm string `json:"algorithm,omitempty"`
	Value string `json:"value"`
	Secret string `json:"secret,omitempty"`
}
type Output struct { Result string `json:"result"` }

func main() {
	action.Main(action.Meta{
		Name:        "Crypto",
		Description: "hash hmac base64",
	}, func(ctx context.Context, in Input) (Output, error) {
		alg := strings.ToLower(in.Algorithm)
		if alg == "" { alg = "sha256" }
		switch strings.ToLower(in.Operation) {
		case "base64encode":
			return Output{Result: base64.StdEncoding.EncodeToString([]byte(in.Value))}, nil
		case "base64decode":
			b, err := base64.StdEncoding.DecodeString(in.Value)
			if err != nil { return Output{}, err }
			return Output{Result: string(b)}, nil
		case "hash":
			sum := hashBytes(alg, []byte(in.Value))
			return Output{Result: hex.EncodeToString(sum)}, nil
		case "hmac":
			sum := hmacBytes(alg, []byte(in.Secret), []byte(in.Value))
			return Output{Result: hex.EncodeToString(sum)}, nil
		default:
			return Output{}, fmt.Errorf("unsupported operation")
		}

	})
}

func hashBytes(alg string, b []byte) []byte {
	switch alg {
	case "sha1": s := sha1.Sum(b); return s[:]
	case "md5": s := md5.Sum(b); return s[:]
	default: s := sha256.Sum256(b); return s[:]
	}
}
func hmacBytes(alg string, key, b []byte) []byte {
	var h hash.Hash
	switch alg {
	case "sha1": h = hmac.New(sha1.New, key)
	case "md5": h = hmac.New(md5.New, key)
	default: h = hmac.New(sha256.New, key)
	}
	h.Write(b)
	return h.Sum(nil)
}
