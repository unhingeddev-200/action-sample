package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)


type Input struct {
	Secret string `json:"secret"`
	Digits int `json:"digits,omitempty"`
	Period int `json:"period,omitempty"`
}
type Output struct { Code string `json:"code"` }

func main() {
	action.Main(action.Meta{
		Name:        "Totp",
		Description: "generate totp",
	}, func(ctx context.Context, in Input) (Output, error) {
		digits := in.Digits
		if digits == 0 { digits = 6 }
		period := in.Period
		if period == 0 { period = 30 }
		key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(strings.ReplaceAll(in.Secret, " ", "")))
		if err != nil {
			// try with padding
			key, err = base32.StdEncoding.DecodeString(strings.ToUpper(strings.ReplaceAll(in.Secret, " ", "")))
			if err != nil { return Output{}, err }
		}
		counter := uint64(time.Now().Unix()) / uint64(period)
		code := hotp(key, counter, digits)
		return Output{Code: code}, nil

	})
}

func hotp(key []byte, counter uint64, digits int) string {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)
	mac := hmac.New(sha1.New, key)
	mac.Write(buf)
	sum := mac.Sum(nil)
	off := sum[len(sum)-1] & 0x0f
	val := binary.BigEndian.Uint32(sum[off:off+4]) & 0x7fffffff
	mod := uint32(1)
	for i := 0; i < digits; i++ { mod *= 10 }
	return fmt.Sprintf("%0*d", digits, val%mod)
}
