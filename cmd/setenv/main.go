package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"unicode/utf8"

	"github.com/lunaya-dubai/lunaya-flow-runtime/sdk/action"
)

var envNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

type Input struct {
	// Name + Value set a single env var (convenience form).
	Name  string `json:"name,omitempty" jsonschema:"Environment variable name (letters, numbers, underscore; must start with a letter or underscore). Prefer a fixed name, not a binding."`
	Value any    `json:"value,omitempty" jsonschema:"Value to export. Non-text values are stored as JSON text."`
	// Assignments sets multiple env vars (like Set).
	Assignments []Assignment `json:"assignments,omitempty" jsonschema:"Set several environment variables at once instead of a single name/value"`
}

type Assignment struct {
	Name  string `json:"name" jsonschema:"Environment variable name"`
	Value any    `json:"value" jsonschema:"Value to export for this name"`
}

type Output struct {
	Env   map[string]string `json:"env" jsonschema:"All exported names and their string values"`
	Name  string            `json:"name,omitempty" jsonschema:"Echo of the singular name when used"`
	Value string            `json:"value,omitempty" jsonschema:"Echo of the singular value as text when used"`
}

func main() {
	action.Main(action.Meta{
		Name:        "SetEnv",
		Description: "Export values as real environment variables for later steps (for example so Code can read env.API_TOKEN). Connect this step to the steps that should receive the variables.",
	}, func(ctx context.Context, in Input) (Output, error) {
		env := make(map[string]string)
		add := func(name string, value any) error {
			if name == "" {
				return fmt.Errorf("env name is required")
			}
			if !envNameRe.MatchString(name) {
				return fmt.Errorf("invalid env name %q (want [A-Za-z_][A-Za-z0-9_]*)", name)
			}
			s, err := stringifyEnvValue(value)
			if err != nil {
				return fmt.Errorf("env %s: %w", name, err)
			}
			env[name] = s
			return nil
		}

		if in.Name != "" || in.Value != nil {
			if err := add(in.Name, in.Value); err != nil {
				return Output{}, err
			}
		}
		for _, a := range in.Assignments {
			if err := add(a.Name, a.Value); err != nil {
				return Output{}, err
			}
		}
		if len(env) == 0 {
			return Output{}, fmt.Errorf("provide name/value or assignments")
		}

		out := Output{Env: env}
		if in.Name != "" {
			out.Name = in.Name
			out.Value = env[in.Name]
		}
		return out, nil
	})
}

func stringifyEnvValue(v any) (string, error) {
	switch t := v.(type) {
	case nil:
		return "", nil
	case string:
		if !utf8.ValidString(t) {
			return "", fmt.Errorf("value is not valid UTF-8")
		}
		return t, nil
	case bool:
		if t {
			return "true", nil
		}
		return "false", nil
	case float64:
		return trimFloat(t), nil
	case json.Number:
		return t.String(), nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
}

func trimFloat(f float64) string {
	s := fmt.Sprintf("%v", f)
	return s
}
