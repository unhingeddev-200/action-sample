package main

import (
	"context"
	"encoding/base64"
	"path"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func withMemPut(t *testing.T, fs afero.Fs) {
	t.Helper()
	prev := putObject
	putObject = func(_ context.Context, key string, data []byte) error {
		if err := fs.MkdirAll(path.Dir(key), 0o755); err != nil {
			return err
		}
		return afero.WriteFile(fs, key, data, 0o644)
	}
	t.Cleanup(func() { putObject = prev })
}

func TestAppendReadConcatUTF8(t *testing.T) {
	fs := afero.NewMemMapFs()
	withMemPut(t, fs)
	ctx := context.Background()
	key := "debug/run-1/results"

	if _, err := appendFile(ctx, fs, key, "hello", "utf8"); err != nil {
		t.Fatal(err)
	}
	if _, err := appendFile(ctx, fs, key, " world", "utf8"); err != nil {
		t.Fatal(err)
	}

	out, err := readFile(fs, key, "utf8")
	if err != nil {
		t.Fatal(err)
	}
	if out.Content != "hello world" {
		t.Fatalf("content=%q want %q", out.Content, "hello world")
	}
	if !isCollection(fs, key) {
		t.Fatal("expected collection marker")
	}
}

func TestAppendReadConcatBase64(t *testing.T) {
	fs := afero.NewMemMapFs()
	withMemPut(t, fs)
	ctx := context.Background()
	key := "bin/blob"

	a := []byte{0x00, 0x01, 0xff}
	b := []byte{0x10, 0x20}
	if _, err := appendFile(ctx, fs, key, base64.StdEncoding.EncodeToString(a), "base64"); err != nil {
		t.Fatal(err)
	}
	if _, err := appendFile(ctx, fs, key, base64.StdEncoding.EncodeToString(b), "base64"); err != nil {
		t.Fatal(err)
	}

	out, err := readFile(fs, key, "base64")
	if err != nil {
		t.Fatal(err)
	}
	got, err := base64.StdEncoding.DecodeString(out.Content)
	if err != nil {
		t.Fatal(err)
	}
	want := append(append([]byte{}, a...), b...)
	if string(got) != string(want) {
		t.Fatalf("got %v want %v", got, want)
	}
}

func TestDeleteCollection(t *testing.T) {
	fs := afero.NewMemMapFs()
	withMemPut(t, fs)
	ctx := context.Background()
	key := "tmp/collect"

	if _, err := appendFile(ctx, fs, key, "x", "utf8"); err != nil {
		t.Fatal(err)
	}
	if _, err := deletePath(fs, key); err != nil {
		t.Fatal(err)
	}
	st, err := statPath(fs, key)
	if err != nil {
		t.Fatal(err)
	}
	if st.Exists {
		t.Fatal("expected exists=false after delete")
	}
}

func TestWriteThenAppendConflict(t *testing.T) {
	fs := afero.NewMemMapFs()
	withMemPut(t, fs)
	ctx := context.Background()
	key := "plain.txt"

	if _, err := writeFile(fs, key, "hi", "utf8"); err != nil {
		t.Fatal(err)
	}
	_, err := appendFile(ctx, fs, key, "!", "utf8")
	if err == nil || !strings.Contains(err.Error(), "plain file") {
		t.Fatalf("expected plain file conflict, got %v", err)
	}
}

func TestAppendThenWriteConflict(t *testing.T) {
	fs := afero.NewMemMapFs()
	withMemPut(t, fs)
	ctx := context.Background()
	key := "collect.txt"

	if _, err := appendFile(ctx, fs, key, "a", "utf8"); err != nil {
		t.Fatal(err)
	}
	_, err := writeFile(fs, key, "b", "utf8")
	if err == nil || !strings.Contains(err.Error(), "appendable file") {
		t.Fatalf("expected appendable conflict, got %v", err)
	}
}

func TestExistsCollection(t *testing.T) {
	fs := afero.NewMemMapFs()
	withMemPut(t, fs)
	key := "c/path"
	st, err := statPath(fs, key)
	if err != nil {
		t.Fatal(err)
	}
	if st.Exists {
		t.Fatal("expected missing")
	}
	if _, err := appendFile(context.Background(), fs, key, "z", "utf8"); err != nil {
		t.Fatal(err)
	}
	st, err = statPath(fs, key)
	if err != nil {
		t.Fatal(err)
	}
	if !st.Exists {
		t.Fatal("expected exists after append")
	}
}
