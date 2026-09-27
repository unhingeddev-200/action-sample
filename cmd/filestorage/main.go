package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
	s3sdk "github.com/lunaya-dubai/lunaya-flow-runtime/sdk/s3"
	"github.com/spf13/afero"
)

const (
	collectionMarkerName = ".collection"
	collectionPartsDir   = "parts"
)

type Input struct {
	Operation string `json:"operation" jsonschema:"What to do: read, write, append, list, delete, or exists"`
	Path      string `json:"path" jsonschema:"File or folder path inside the shared file storage (SeaweedFS genai bucket)"`
	Content   string `json:"content,omitempty" jsonschema:"File contents when writing or appending"`
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

// putObject uploads bytes; overridden in tests.
var putObject = func(ctx context.Context, key string, data []byte) error {
	return s3sdk.PutBytes(ctx, s3sdk.EnvFromOS(), key, data)
}

func main() {
	action.Main(action.Meta{
		Name:        "FileStorage",
		Description: "Read, write, append, list, delete, or check files in the platform file storage (shared S3-compatible storage).",
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
		case "append":
			return appendFile(ctx, fs, key, in.Content, in.Encoding)
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

func decodeContent(content, enc string) ([]byte, error) {
	mode := encodingMode(enc)
	if mode == "base64" {
		raw, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return nil, fmt.Errorf("decode base64 content: %w", err)
		}
		return raw, nil
	}
	return []byte(content), nil
}

func encodeContent(raw []byte, enc string) string {
	if encodingMode(enc) == "base64" {
		return base64.StdEncoding.EncodeToString(raw)
	}
	return string(raw)
}

func collectionMarkerKey(key string) string {
	return path.Join(key, collectionMarkerName)
}

func collectionPartsPrefix(key string) string {
	return path.Join(key, collectionPartsDir)
}

func isCollection(fs afero.Fs, key string) bool {
	if key == "" {
		return false
	}
	_, err := fs.Stat(collectionMarkerKey(key))
	return err == nil
}

func pathExists(fs afero.Fs, key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("path is required")
	}
	if isCollection(fs, key) {
		return true, nil
	}
	_, err := fs.Stat(key)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "not found") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func newPartID() (string, error) {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	// Lexicographically time-sortable without a central counter.
	return fmt.Sprintf("%020d-%s", time.Now().UnixNano(), hex.EncodeToString(b[:])), nil
}

func listPartKeys(fs afero.Fs, key string) ([]string, error) {
	partsDir := collectionPartsPrefix(key)
	entries, err := afero.ReadDir(fs, partsDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "not found") {
			return nil, nil
		}
		return nil, err
	}
	keys := make([]string, 0, len(entries))
	for _, ent := range entries {
		if ent.IsDir() {
			continue
		}
		keys = append(keys, path.Join(partsDir, ent.Name()))
	}
	sort.Strings(keys)
	return keys, nil
}

func readCollection(fs afero.Fs, key, enc string) (Output, error) {
	mode := encodingMode(enc)
	partKeys, err := listPartKeys(fs, key)
	if err != nil {
		return Output{}, err
	}
	var buf []byte
	for _, pk := range partKeys {
		f, err := fs.Open(pk)
		if err != nil {
			return Output{}, err
		}
		chunk, err := io.ReadAll(f)
		_ = f.Close()
		if err != nil {
			return Output{}, err
		}
		buf = append(buf, chunk...)
	}
	return Output{OK: true, Content: encodeContent(buf, mode), Encoding: mode}, nil
}

func readFile(fs afero.Fs, key, enc string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	if isCollection(fs, key) {
		return readCollection(fs, key, enc)
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
	return Output{OK: true, Content: encodeContent(raw, mode), Encoding: mode}, nil
}

func writeFile(fs afero.Fs, key, content, enc string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	if isCollection(fs, key) {
		return Output{}, fmt.Errorf("path %q is an appendable file; delete it before write, or use append", key)
	}
	mode := encodingMode(enc)
	raw, err := decodeContent(content, mode)
	if err != nil {
		return Output{}, err
	}
	if err := putObject(context.Background(), key, raw); err != nil {
		return Output{}, err
	}
	return Output{OK: true, Encoding: mode}, nil
}

func appendFile(ctx context.Context, fs afero.Fs, key, content, enc string) (Output, error) {
	if key == "" {
		return Output{}, fmt.Errorf("path is required")
	}
	mode := encodingMode(enc)
	raw, err := decodeContent(content, mode)
	if err != nil {
		return Output{}, err
	}

	if !isCollection(fs, key) {
		// Refuse to turn a plain single-object file into a collection.
		if _, err := fs.Stat(key); err == nil {
			return Output{}, fmt.Errorf("path %q is a plain file; delete it before append, or use write", key)
		} else if !errors.Is(err, os.ErrNotExist) && !strings.Contains(err.Error(), "not found") {
			return Output{}, err
		}
		if err := putObject(ctx, collectionMarkerKey(key), []byte("collection/v1\n")); err != nil {
			return Output{}, fmt.Errorf("create collection marker: %w", err)
		}
	}

	partID, err := newPartID()
	if err != nil {
		return Output{}, err
	}
	partKey := path.Join(collectionPartsPrefix(key), partID)
	if err := putObject(ctx, partKey, raw); err != nil {
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
	if isCollection(fs, key) {
		if err := fs.RemoveAll(key); err != nil {
			return Output{}, err
		}
		return Output{OK: true}, nil
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
	exists, err := pathExists(fs, key)
	if err != nil {
		return Output{}, err
	}
	return Output{OK: true, Exists: exists}, nil
}
