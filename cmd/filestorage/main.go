package main

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	s3sdk "github.com/lunaya-dubai/lunaya-flow-runtime/sdk/s3"
	"github.com/spf13/afero"
)

type Input struct {
	Operation string `json:"operation" jsonschema:"What to do: read, write, list, delete, or exists"`
	Path      string `json:"path" jsonschema:"File or folder path inside the shared file storage (SeaweedFS genai bucket)"`
	Content   string `json:"content,omitempty" jsonschema:"File contents when writing"`
	Encoding  string `json:"encoding,omitempty" jsonschema:"How content is encoded: utf8 (default) or base64"`
	Recursive bool   `json:"recursive,omitempty" jsonschema:"When listing, include files in subfolders"`
}

type FileEntry struct {
	Name  string `json:"name" jsonschema:"File or folder name"`
	Path  string `json:"path" jsonschema:"Full path of this entry"`
	Size  int64  `json:"size,omitempty" jsonschema:"Size in bytes (files only)"`
	IsDir bool   `json:"isDir" jsonschema:"True if this entry is a folder"`
}

type Output struct {
	OK       bool        `json:"ok" jsonschema:"True when the operation succeeded"`
	Content  string      `json:"content,omitempty" jsonschema:"File contents when reading"`
	Encoding string      `json:"encoding,omitempty" jsonschema:"Encoding of content"`
	Exists   bool        `json:"exists,omitempty" jsonschema:"Whether the path exists (exists operation)"`
	Files    []FileEntry `json:"files,omitempty" jsonschema:"Entries returned by list"`
}

func main() {
	action.Main(action.Meta{
		Name:        "FileStorage",
		Description: "Read, write, list, delete, or check files in the platform file storage (shared S3-compatible storage).",
	}, func(ctx context.Context, in Input) (Output, error) {
		fs, err := s3sdk.NewFSFromEnv()
		if err != nil {
			return Output{}, fmt.Errorf("s3 filesystem: %w", err)
		}
		key := cleanKey(in.Path)
		switch strings.ToLower(strings.TrimSpace(in.Operation)) {
		case "read":
			return readFile(fs, key, in.Encoding)
		case "write":
			return writeFile(fs, key, in.Content, in.Encoding)
		case "list":
			return listFiles(fs, key, in.Recursive)
		case "delete":
			return deletePath(fs, key)
		case "exists":
			return statPath(fs, key)
		default:
			return Output{}, fmt.Errorf("unsupported operation %q", in.Operation)
		}
	})
}

func cleanKey(p string) string {
	p = strings.TrimSpace(p)
	p = strings.TrimPrefix(p, "/")
	p = path.Clean(p)
	if p == "." {
		return ""
	}
	return p
}

func encodingMode(enc string) string {
	switch strings.ToLower(strings.TrimSpace(enc)) {
	case "", "utf8", "utf-8", "text":
		return "utf8"
	case "base64":
		return "base64"
	default:
		return "utf8"
	}
}

func readFile(fs afero.Fs, key, enc string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	mode := encodingMode(enc)
	f, err := fs.Open(key)
	if err != nil {
		return Output{}, err
	}
	defer f.Close()
	raw, err := io.ReadAll(f)
	if err != nil {
		return Output{}, err
	}
	out := Output{OK: true, Encoding: mode}
	if mode == "base64" {
		out.Content = base64.StdEncoding.EncodeToString(raw)
	} else {
		out.Content = string(raw)
	}
	return out, nil
}

func writeFile(_ afero.Fs, key, content, enc string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	mode := encodingMode(enc)
	var raw []byte
	switch mode {
	case "base64":
		var err error
		raw, err = base64.StdEncoding.DecodeString(content)
		if err != nil {
			return Output{}, fmt.Errorf("decode base64 content: %w", err)
		}
	default:
		raw = []byte(content)
	}
	cfg := s3sdk.EnvFromOS()
	if err := s3sdk.PutBytes(context.Background(), cfg, key, raw); err != nil {
		return Output{}, err
	}
	return Output{OK: true, Encoding: mode}, nil
}

func listFiles(fs afero.Fs, prefix string, recursive bool) (Output, error) {
	root := prefix
	if root == "" {
		root = "."
	}
	entries := make([]FileEntry, 0, 16)
	if !recursive {
		dir, err := afero.ReadDir(fs, root)
		if err != nil {
			return Output{}, err
		}
		for _, ent := range dir {
			name := ent.Name()
			child := name
			if prefix != "" {
				child = path.Join(prefix, name)
			}
			entries = append(entries, FileEntry{
				Name:  name,
				Path:  child,
				Size:  ent.Size(),
				IsDir: ent.IsDir(),
			})
		}
		return Output{OK: true, Files: entries}, nil
	}
	err := afero.Walk(fs, root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if p == root {
			return nil
		}
		entries = append(entries, FileEntry{
			Name:  path.Base(p),
			Path:  p,
			Size:  info.Size(),
			IsDir: info.IsDir(),
		})
		return nil
	})
	if err != nil {
		return Output{}, err
	}
	return Output{OK: true, Files: entries}, nil
}

func deletePath(fs afero.Fs, key string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	info, err := fs.Stat(key)
	if err != nil {
		return Output{}, err
	}
	if info.IsDir() {
		if err := fs.RemoveAll(key); err != nil {
			return Output{}, err
		}
	} else if err := fs.Remove(key); err != nil {
		return Output{}, err
	}
	return Output{OK: true}, nil
}

func statPath(fs afero.Fs, key string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	_, err := fs.Stat(key)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "not found") {
			return Output{OK: true, Exists: false}, nil
		}
		return Output{}, err
	}
	return Output{OK: true, Exists: true}, nil
}
